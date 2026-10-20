package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
	"roundfix/internal/preflight"
)

const linuxLaneHistory = "/repo-objects"

const linuxLaneCap = 20 * time.Minute
const linuxLaneScript = `fail() { printf 'roundfix-linux-lane: container: %s\n' "$1" >&2; exit 3; }; for tool in tar git go; do command -v "$tool" >/dev/null 2>&1 || fail "missing tool $tool"; done; [ -d /repo-objects ] && [ -r /repo-objects ] && [ -x /repo-objects ] || fail "cannot read the repository history at /repo-objects"; mkdir -p /tmp/work && cd /tmp/work || fail "cannot prepare /tmp/work"; tar -xf - || fail "cannot unpack the candidate"; { git init -q && echo /repo-objects > .git/objects/info/alternates && git add -A && git -c user.name=roundfix -c user.email=roundfix@localhost commit -q -m "Linux Lane"; } || fail "cannot record the candidate"; pkgs="$(go list -e -f '{{if or .GoFiles .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}' "$@")" && if [ -z "$pkgs" ]; then echo "no package builds on linux"; exit 0; fi && exec go test -count=1 $pkgs`

type linuxLaneResult struct {
	Ran, Passed                    bool
	Skip, Reason, LogPath, Blocker string
	Packages                       []string
}

type linuxLaneRequest struct {
	WorkDir, Head, Image string
	Label, LogPath, Name string
	All                  bool
}

func (workflow *commandDeliveryWorkflow) runLinuxLane(ctx context.Context, request linuxLaneRequest) (result linuxLaneResult, resultErr error) {
	workDir, head := request.WorkDir, request.Head
	result.LogPath = request.LogPath
	if err := os.MkdirAll(filepath.Dir(result.LogPath), 0755); err != nil {
		return result, fmt.Errorf("prepare Linux Lane log: %w", err)
	}
	log, err := os.Create(result.LogPath)
	if err != nil {
		return result, fmt.Errorf("create Linux Lane log: %w", err)
	}
	defer func() { resultErr = errors.Join(resultErr, log.Close()) }()
	console := workflow.log
	if console == nil {
		console = os.Stderr
	}
	skip := func(reason string) (linuxLaneResult, error) {
		result.Skip = reason
		result.Passed = true
		line := fmt.Sprintf("roundfix: notice: %sLinux Lane skipped: %s\n", request.Label, reason)
		_, err := io.WriteString(io.MultiWriter(console, log), line)
		return result, err
	}
	runner := workflow.git
	if runner == nil {
		runner = preflight.ExecGitRunner{}
	}
	listing, err := runner.RunGit(ctx, workDir, "ls-tree", "--name-only", head, "--", "go.mod")
	if err != nil {
		return result, fmt.Errorf("inspect Linux Lane go.mod: %w", err)
	}
	if strings.TrimSpace(listing) == "" {
		return skip("no go.mod at the repository root")
	}
	image := request.Image
	if image == "" {
		module, err := runner.RunGit(ctx, workDir, "show", head+":go.mod")
		if err != nil {
			return result, fmt.Errorf("read Linux Lane go.mod: %w", err)
		}
		version := ""
		for _, line := range strings.Split(module, "\n") {
			fields := strings.Fields(line)
			if len(fields) >= 2 && fields[0] == "go" {
				version = fields[1]
				break
			}
		}
		if version == "" {
			return skip("go.mod declares no go version")
		}
		image = "golang:" + version
	}
	if request.All {
		result.Packages = []string{"./..."}
	} else {
		branch := preflight.DetectDefaultBranch(ctx, workDir, "", runner)
		if branch.Source == preflight.DefaultBranchUndetermined {
			return result, errors.New("Linux Lane: repository default branch is unknown")
		}
		remote := strings.TrimSpace(workflow.loaded.Config.Watch.PushRemote)
		if remote == "" {
			remote = "origin"
		}
		base, err := runner.RunGit(ctx, workDir, "merge-base", remote+"/"+branch.Name, head)
		if err != nil {
			return result, fmt.Errorf("read Linux Lane merge base: %w", err)
		}
		changed, err := runner.RunGit(ctx, workDir, "diff", "--name-only", "-z", "--no-renames", strings.TrimSpace(base), head)
		if err != nil {
			return result, fmt.Errorf("read Linux Lane changed paths: %w", err)
		}
		tree, err := runner.RunGit(ctx, workDir, "ls-tree", "-r", "--name-only", "-z", head)
		if err != nil {
			return result, fmt.Errorf("read Linux Lane tree: %w", err)
		}
		result.Packages = linuxLanePackages(changed, tree)
	}
	if len(result.Packages) == 0 {
		return skip("no changed Go package")
	}
	objects, err := linuxLaneObjectStore(ctx, runner, workDir)
	if err != nil {
		return result, err
	}
	historyMount, reason, err := linuxLaneHistoryMount(objects)
	if err != nil {
		return result, err
	}
	if reason != "" {
		return skip(reason)
	}

	environment := commandDependenciesForContext(ctx).environ()
	goBinary := linuxLaneExecutable("go", environment, workDir)
	if goBinary == "" {
		return skip("go is not on PATH")
	}
	docker := linuxLaneExecutable("docker", environment, workDir)
	if docker == "" {
		return skip("docker is not on PATH")
	}
	command := func(commandCtx context.Context, binary string, args ...string) *exec.Cmd {
		cmd := exec.CommandContext(commandCtx, binary, args...)
		cmd.Dir, cmd.Env = workDir, environment
		return cmd
	}
	output, err := command(ctx, docker, "info").CombinedOutput()
	if ctx.Err() != nil {
		return result, ctx.Err()
	}
	if err != nil {
		return skip("docker info failed: " + linuxLaneFirstLine(output))
	}
	if _, err := command(ctx, docker, "image", "inspect", image).CombinedOutput(); err != nil {
		if ctx.Err() != nil {
			return result, ctx.Err()
		}
		return skip("image " + image + " is not present locally")
	}
	cache, err := command(ctx, goBinary, "env", "GOMODCACHE").Output()
	if err != nil {
		return result, fmt.Errorf("read Linux Lane host module cache: %w", err)
	}
	name := request.Name
	args := []string{"run", "--rm", "-i", "--init", "--pull", "never", "--network", "none", "--name", name, "--user", "65534:65534",
		"-e", "HOME=/tmp", "-e", "GOCACHE=/tmp/go-build", "-e", "GOMODCACHE=/go/pkg/mod", "-e", "GOFLAGS=-mod=readonly", "-e", "GOPROXY=off", "-e", "GOTOOLCHAIN=local", "-e", "GOSUMDB=off",
		"-v", strings.TrimSpace(string(cache)) + ":/go/pkg/mod:ro", "-v", historyMount, image, "sh", "-ceu", linuxLaneScript, "linux-lane"}
	args = append(args, result.Packages...)
	quoted := []string{"docker"}
	for _, arg := range args {
		quoted = append(quoted, posixSingleQuote(arg))
	}
	if _, err := fmt.Fprintln(log, strings.Join(quoted, " ")); err != nil {
		return result, err
	}
	laneCtx, cancel := context.WithTimeout(ctx, linuxLaneCap)
	defer cancel()
	archive := exec.CommandContext(laneCtx, "git", "archive", "--format=tar", head)
	archive.Dir, archive.Env = workDir, environment
	input, err := archive.StdoutPipe()
	if err != nil {
		return result, fmt.Errorf("open Linux Lane archive: %w", err)
	}
	var archiveErrors bytes.Buffer
	archive.Stderr = &archiveErrors
	if err := archive.Start(); err != nil {
		return result, fmt.Errorf("start Linux Lane archive: %w", err)
	}
	run := command(laneCtx, docker, args...)
	run.Stdin = input
	var runOutput bytes.Buffer
	combinedOutput := io.MultiWriter(log, &runOutput)
	run.Stdout, run.Stderr = combinedOutput, combinedOutput
	result.Ran = true
	runErr := run.Run()
	// Closing the pipe also releases git when Docker refuses the input early.
	closeErr := input.Close()
	archiveErr := archive.Wait()
	failure := func(code int, capped bool) (linuxLaneResult, error) {
		summary, blocker := classifyLinuxLaneOutcome(code, capped, runOutput.Bytes())
		result.Reason = summary + "; see " + result.LogPath
		detail := summary
		if blocker == delivery.BlockerLinuxLaneContainer {
			detail = strings.TrimPrefix(summary, "Linux Lane container failure: ")
		}
		result.Blocker = blocker + ": " + detail + "; see " + result.LogPath
		_, err := fmt.Fprintln(log, "\n"+summary)
		return result, err
	}
	if laneCtx.Err() != nil {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.WithoutCancel(ctx), time.Minute)
		defer cleanupCancel()
		cleanup := command(cleanupCtx, docker, "rm", "-f", name)
		cleanup.Stdout, cleanup.Stderr = log, log
		cleanupErr := cleanup.Run()
		if cleanupErr != nil {
			_, _ = fmt.Fprintf(log, "remove Linux Lane container: %v\n", cleanupErr)
		}
		if ctx.Err() != nil {
			return result, errors.Join(ctx.Err(), cleanupErr)
		}
		return failure(-1, true)
	}
	if runErr != nil {
		var exit *exec.ExitError
		if !errors.As(runErr, &exit) {
			return result, fmt.Errorf("start Linux Lane container: %w", runErr)
		}
		if exit.ExitCode() == 125 {
			return skip("docker could not start the container: " + linuxLaneFirstLine(runOutput.Bytes()))
		}
		return failure(exit.ExitCode(), false)
	}
	if archiveErr != nil {
		return result, fmt.Errorf("archive Linux Lane candidate: %w: %s", archiveErr, strings.TrimSpace(archiveErrors.String()))
	}
	if closeErr != nil {
		return result, fmt.Errorf("close Linux Lane archive: %w", closeErr)
	}
	if _, err := fmt.Fprintf(log, "\nLinux Lane passed %d package(s) on %s\n", len(result.Packages), image); err != nil {
		return result, err
	}
	result.Passed = true
	_, err = fmt.Fprintf(console, "roundfix: %sLinux Lane passed %d package(s) on %s\n", request.Label, len(result.Packages), image)
	return result, err
}

// classifyLinuxLaneOutcome handles runs that started; exit 125 is a skip at the caller.
// Only the lane script's own marker contributes output to the operator's cause.
func classifyLinuxLaneOutcome(code int, capped bool, output []byte) (summary, blocker string) {
	cause := ""
	switch {
	case capped:
		cause = "the container exceeded " + linuxLaneCap.String()
	case code == 0:
		return "", ""
	case code == 3:
		cause = "the lane script exited 3 before go test"
		for _, line := range strings.Split(string(output), "\n") {
			if text, ok := strings.CutPrefix(strings.TrimSuffix(line, "\r"), "roundfix-linux-lane: container: "); ok {
				cause = text
			}
		}
	case code == 126:
		cause = "docker run exited 126: the image could not invoke the lane command"
	case code == 127:
		cause = "docker run exited 127: the image lacks the lane command"
	case code >= 128:
		cause = fmt.Sprintf("docker run exited %d: the container process ended by signal %d", code, code-128)
	default:
		return fmt.Sprintf("Linux Lane change failure: go test exited %d", code), delivery.BlockerGateFailed
	}
	return "Linux Lane container failure: " + cause, delivery.BlockerLinuxLaneContainer
}

func linuxLanePackages(changed, tree string) []string {
	directories := make(map[string]bool)
	for name := range strings.SplitSeq(tree, "\x00") {
		if strings.HasSuffix(name, ".go") {
			directories[path.Dir(name)] = true
		}
	}
	selected := make(map[string]bool)
	for name := range strings.SplitSeq(changed, "\x00") {
		if name == "" {
			continue
		}
		if name == "go.mod" || name == "go.sum" {
			return []string{"./..."}
		}
		for dir := path.Dir(name); ; dir = path.Dir(dir) {
			if directories[dir] {
				pkg := "."
				if dir != "." {
					pkg = "./" + dir
				}
				selected[pkg] = true
				break
			}
			if dir == "." {
				break
			}
		}
	}
	packages := make([]string, 0, len(selected))
	for pkg := range selected {
		packages = append(packages, pkg)
	}
	slices.Sort(packages)
	return packages
}

func linuxLaneFirstLine(output []byte) string {
	for _, line := range strings.Split(string(output), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			return line
		}
	}
	return ""
}

// Resolve against the injected owner's PATH without reading or changing the
// process environment. Relative entries are relative to the item worktree.
func linuxLaneExecutable(name string, environment []string, workDir string) string {
	searchPath := ""
	for _, entry := range environment {
		if key, value, ok := strings.Cut(entry, "="); ok && strings.EqualFold(key, "PATH") {
			searchPath = value
		}
	}
	names := []string{name}
	if runtime.GOOS == "windows" {
		names = []string{name + ".exe", name + ".com", name + ".bat", name + ".cmd"}
	}
	for _, dir := range filepath.SplitList(searchPath) {
		if !filepath.IsAbs(dir) {
			dir = filepath.Join(workDir, dir)
		}
		for _, candidate := range names {
			full := filepath.Join(dir, candidate)
			info, err := os.Stat(full)
			if err == nil && !info.IsDir() && (runtime.GOOS == "windows" || info.Mode()&0111 != 0) {
				return full
			}
		}
	}
	return ""
}

// linuxLaneOptIn reads delivery.linux_lane and its image from the fetched default branch.
// A non-empty reason means no fetched default exists to read it from; a
// cancelled context or an unreadable Project Config at that default is an
// error, never a skip.
func (workflow *commandDeliveryWorkflow) linuxLaneOptIn(ctx context.Context, gitRoot string) (bool, string, string, error) {
	if err := ctx.Err(); err != nil {
		return false, "", "", err
	}
	runner := workflow.git
	if runner == nil {
		runner = preflight.ExecGitRunner{}
	}
	branch := preflight.DetectDefaultBranch(ctx, gitRoot, "", runner)
	if err := ctx.Err(); err != nil {
		return false, "", "", err
	}
	if branch.Source == preflight.DefaultBranchUndetermined {
		return false, "", "the repository default branch is unknown", nil
	}
	remote := strings.TrimSpace(workflow.loaded.Config.Watch.PushRemote)
	if remote == "" {
		remote = "origin"
	}
	ref := "refs/remotes/" + remote + "/" + branch.Name
	defaultHead, err := runner.RunGit(ctx, gitRoot, "rev-parse", "--verify", ref+"^{commit}")
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return false, "", "", ctxErr
		}
		return false, "", "the fetched default " + ref + " cannot be resolved", nil
	}
	config, err := roundconfig.DeliveryConfigAtCommit(ctx, runner, gitRoot, workflow.loaded.UserConfigPath, strings.TrimSpace(defaultHead))
	if err != nil {
		return false, "", "", fmt.Errorf("read Linux Lane configuration: %w", err)
	}
	return config.Delivery.LinuxLane, config.Delivery.LinuxLaneImage, "", nil
}

// recordLinuxLaneSkip writes a skip that happens before the lane starts to
// the console and to the lane log, like every other Linux Lane skip.
func (workflow *commandDeliveryWorkflow) recordLinuxLaneSkip(specSlug, reason string) (resultErr error) {
	artifactDir, err := roundconfig.ValidateArtifactDirectory(workflow.loaded.Config.Defaults.ArtifactDir, workflow.loaded.GitRoot, workflow.loaded.HomeDir)
	if err != nil {
		return err
	}
	logPath := filepath.Join(artifactDir, "delivery", specSlug, "linux-lane.log")
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return fmt.Errorf("prepare Linux Lane log: %w", err)
	}
	log, err := os.Create(logPath)
	if err != nil {
		return fmt.Errorf("create Linux Lane log: %w", err)
	}
	defer func() { resultErr = errors.Join(resultErr, log.Close()) }()
	console := workflow.log
	if console == nil {
		console = os.Stderr
	}
	line := fmt.Sprintf("roundfix: notice: Delivery Queue item %s: Linux Lane skipped: %s\n", specSlug, reason)
	_, err = io.WriteString(io.MultiWriter(console, log), line)
	return err
}

// linuxLaneObjectStore resolves the host store Git reads, including the common
// directory's store when workDir is a linked worktree.
func linuxLaneObjectStore(ctx context.Context, runner preflight.GitRunner, workDir string) (string, error) {
	objects, err := runner.RunGit(ctx, workDir, "rev-parse", "--path-format=absolute", "--git-path", "objects")
	if err != nil {
		return "", fmt.Errorf("read Linux Lane object store: %w", err)
	}
	return filepath.Clean(strings.TrimSpace(objects)), nil
}

// linuxLaneHistoryMount exposes all host repository objects to the lane container
// read-only. Neither the working tree nor refs are mounted. The script writes
// its alternate only into the container's own /tmp/work/.git, never this store.
func linuxLaneHistoryMount(objects string) (mount, skip string, err error) {
	if strings.Contains(strings.TrimPrefix(objects, filepath.VolumeName(objects)), ":") {
		return "", "the object store " + objects + " cannot be mounted: its path holds a colon", nil
	}
	info, err := os.Stat(objects)
	if err != nil || !info.IsDir() {
		// In particular, do not let Docker create a missing source directory.
		return "", "the object store " + objects + " cannot be mounted: it is not a directory", nil
	}
	alternates, err := os.ReadFile(filepath.Join(objects, "info", "alternates"))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return "", "", fmt.Errorf("read Linux Lane object store alternates: %w", err)
	}
	for line := range strings.SplitSeq(string(alternates), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			return "", "the object store " + objects + " borrows objects through info/alternates, which the Linux Lane does not mount", nil
		}
	}
	return objects + ":" + linuxLaneHistory + ":ro", "", nil
}
