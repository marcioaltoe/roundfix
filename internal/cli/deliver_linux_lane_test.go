// Boundary: disposable Git repositories and compiled Docker/Go stand-ins.
// No container or external service is used; loopback only signals cancellation.
package cli

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
	"roundfix/internal/testfixture"
)

const linuxLaneToolSource = `package main
import ("encoding/json"; "fmt"; "io"; "net"; "os"; "path/filepath"; "strconv"; "strings")
func main() {
 tool := strings.TrimSuffix(filepath.Base(os.Args[0]), ".exe")
 dir, err := os.Getwd(); if err != nil { panic(err) }
 record, err := os.OpenFile(os.Getenv("LANE_RECORD"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600); if err != nil { panic(err) }
 if err := json.NewEncoder(record).Encode(struct{Tool string; Args, Env []string; Dir string}{tool, os.Args[1:], os.Environ(), dir}); err != nil { panic(err) }
 if err := record.Close(); err != nil { panic(err) }
 if tool == "go" { fmt.Println(os.Getenv("LANE_CACHE")); return }
 switch os.Args[1] {
 case "info": if os.Getenv("LANE_INFO_FAIL") == "1" { fmt.Fprint(os.Stderr, "\n daemon unavailable\nsecond line\n"); os.Exit(1) }
 case "image": if os.Getenv("LANE_IMAGE_FAIL") == "1" { os.Exit(1) }
 case "run":
  data, err := io.ReadAll(os.Stdin); if err != nil { panic(err) }; if err := os.WriteFile(os.Getenv("LANE_ARCHIVE"), data, 0600); err != nil { panic(err) }
  if addr := os.Getenv("LANE_READY"); addr != "" { conn, err := net.Dial("tcp", addr); if err != nil { panic(err) }; defer conn.Close(); if _, err := io.Copy(io.Discard, conn); err != nil { panic(err) }; return }
  fmt.Fprint(os.Stdout, "\n fixture output\nsecond line\n")
  fmt.Fprint(os.Stderr, os.Getenv("LANE_OUTPUT"))
  code, err := strconv.Atoi(os.Getenv("LANE_EXIT")); if err != nil { panic(err) }; os.Exit(code)
 }
}`

type linuxLaneStart struct {
	Tool      string
	Args, Env []string
	Dir       string
}
type linuxLaneFixture struct {
	workflow                                 *commandDeliveryWorkflow
	repo, record, archive, cache, head, base string
	environment                              []string
	log                                      bytes.Buffer
}

func linuxLaneTools(t *testing.T) (string, string) {
	t.Helper()
	suffix := ""
	if runtime.GOOS == "windows" {
		suffix = ".exe"
	}
	return testfixture.FixtureBinary(t, "docker"+suffix, linuxLaneToolSource), testfixture.FixtureBinary(t, "go"+suffix, linuxLaneToolSource)
}

func newLinuxLaneFixture(t *testing.T, docker, goBinary string, files map[string]string) *linuxLaneFixture {
	t.Helper()
	root := t.TempDir()
	f := &linuxLaneFixture{repo: filepath.Join(root, "item"), record: filepath.Join(root, "starts.jsonl"), archive: filepath.Join(root, "archive.tar"), cache: filepath.Join(root, "module cache")}
	gittest.InitRepo(t, f.repo, "-b", "main")
	mustWrite(t, filepath.Join(f.repo, ".roundfixrc.yml"), "delivery:\n  linux_lane: true\n")
	for name, data := range files {
		mustMkdir(t, filepath.Dir(filepath.Join(f.repo, name)))
		mustWrite(t, filepath.Join(f.repo, name), data)
	}
	gittest.Run(t, f.repo, "add", ".")
	gittest.Run(t, f.repo, "commit", "--allow-empty", "-m", "base")
	f.base = strings.TrimSpace(gittest.Run(t, f.repo, "rev-parse", "HEAD"))
	gittest.Run(t, f.repo, "update-ref", "refs/remotes/origin/main", f.base)
	gittest.Run(t, f.repo, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/main")
	f.head = f.base
	home := filepath.Join(root, "owner home")
	mustMkdir(t, home)
	f.workflow = &commandDeliveryWorkflow{git: preflight.ExecGitRunner{}, log: &f.log, loaded: roundconfig.Loaded{GitRoot: f.repo, HomeDir: home}}
	f.workflow.loaded.Config.Defaults.ArtifactDir = filepath.Join(root, "artifacts")
	f.workflow.loaded.Config.Defaults.Verification = "printf 'repository passed'"
	f.workflow.loaded.Config.Delivery.LinuxLane = true
	f.environment = []string{"PATH=" + strings.Join([]string{filepath.Dir(docker), filepath.Dir(goBinary)}, string(os.PathListSeparator)), "HOME=" + filepath.Join(root, "operator home"), "LANE_RECORD=" + f.record, "LANE_ARCHIVE=" + f.archive, "LANE_CACHE=" + f.cache, "LANE_EXIT=0", "TYPESAFE_API_KEY=private-owner-key", "DOCKER_HOST=fixture-context"}
	return f
}

func (f *linuxLaneFixture) context(ctx context.Context) context.Context {
	deps := defaultCommandDependencies()
	deps.environ = func() []string { return append([]string(nil), f.environment...) }
	return contextWithCommandDependencies(ctx, deps)
}
func (f *linuxLaneFixture) commit(t *testing.T, name, data string) {
	t.Helper()
	mustMkdir(t, filepath.Dir(filepath.Join(f.repo, name)))
	mustWrite(t, filepath.Join(f.repo, name), data)
	gittest.Run(t, f.repo, "add", ".")
	gittest.Run(t, f.repo, "commit", "-m", "candidate")
	f.head = strings.TrimSpace(gittest.Run(t, f.repo, "rev-parse", "HEAD"))
}
func (f *linuxLaneFixture) starts(t *testing.T) []linuxLaneStart {
	t.Helper()
	data, err := os.ReadFile(f.record)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	var starts []linuxLaneStart
	for _, line := range bytes.Split(bytes.TrimSpace(data), []byte("\n")) {
		var start linuxLaneStart
		if err := json.Unmarshal(line, &start); err != nil {
			t.Fatal(err)
		}
		starts = append(starts, start)
	}
	return starts
}
func (f *linuxLaneFixture) logPath() string {
	return filepath.Join(f.workflow.loaded.Config.Defaults.ArtifactDir, "delivery", "widget", "linux-lane.log")
}

func TestLinuxLaneRunsTheChangedPackagesInAContainer(t *testing.T) {
	t.Parallel()
	docker, goBinary := linuxLaneTools(t)
	files := map[string]string{"go.mod": "module fixture\n\ngo 1.26.1\n", "z/z.go": "package z\n", "a/a.go": "package a\n", "docs/readme.md": "base"}
	for _, tc := range []struct {
		name, path, data string
		packages         []string
	}{
		{"go file", "a/new.go", "package a\n", []string{"./a"}},
		{"testdata", "z/testdata/nested/input.txt", "new data", []string{"./z"}},
		{"module", "go.mod", "module fixture\n\ngo 1.26.1\n// changed\n", []string{"./..."}},
		{"sum", "go.sum", "changed", []string{"./..."}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newLinuxLaneFixture(t, docker, goBinary, files)
			f.workflow.loaded.Config.Delivery.LinuxLane = false // The default commit enables the lane.
			f.commit(t, tc.path, tc.data)
			// Neither dirty files nor the current checkout's later commit may enter the archive.
			f.commit(t, "docs/later.md", "later commit")
			later := f.head
			f.head = strings.TrimSpace(gittest.Run(t, f.repo, "rev-parse", later+"^"))
			mustWrite(t, filepath.Join(f.repo, "go.mod"), "broken dirty module")
			mustWrite(t, filepath.Join(f.repo, ".roundfixrc.yml"), "delivery:\n  linux_lane: false\n")
			mustWrite(t, filepath.Join(f.repo, "untracked.go"), "package dirty")
			before := gittest.Run(t, f.repo, "status", "--porcelain")
			result, err := f.workflow.Gate(f.context(t.Context()), f.repo, "widget", f.head)
			if err != nil || !result.Passed {
				t.Fatalf("gate=%+v err=%v", result, err)
			}
			starts := f.starts(t)
			if len(starts) != 4 || starts[0].Tool != "docker" || !reflect.DeepEqual(starts[0].Args, []string{"info"}) || !reflect.DeepEqual(starts[1].Args, []string{"image", "inspect", "golang:1.26.1"}) || starts[2].Tool != "go" || !reflect.DeepEqual(starts[2].Args, []string{"env", "GOMODCACHE"}) {
				t.Fatalf("starts=%+v", starts)
			}
			run := starts[3]
			name := fmt.Sprintf("roundfix-linux-lane-widget-%d", os.Getpid())
			// Deliberately independent of the production script constant.
			script := `fail() { printf 'roundfix-linux-lane: container: %s\n' "$1" >&2; exit 3; }; for tool in tar git go; do command -v "$tool" >/dev/null 2>&1 || fail "missing tool $tool"; done; [ -d /repo-objects ] && [ -r /repo-objects ] && [ -x /repo-objects ] || fail "cannot read the repository history at /repo-objects"; mkdir -p /tmp/work && cd /tmp/work || fail "cannot prepare /tmp/work"; tar -xf - || fail "cannot unpack the candidate"; { git init -q && echo /repo-objects > .git/objects/info/alternates && git add -A && git -c user.name=roundfix -c user.email=roundfix@localhost commit -q -m "Linux Lane"; } || fail "cannot record the candidate"; pkgs="$(go list -e -f '{{if or .GoFiles .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}' "$@")" && if [ -z "$pkgs" ]; then echo "no package builds on linux"; exit 0; fi && exec go test -count=1 $pkgs`
			want := []string{"run", "--rm", "-i", "--init", "--pull", "never", "--network", "none", "--name", name, "--user", "65534:65534", "-e", "HOME=/tmp", "-e", "GOCACHE=/tmp/go-build", "-e", "GOMODCACHE=/go/pkg/mod", "-e", "GOFLAGS=-mod=readonly", "-e", "GOPROXY=off", "-e", "GOTOOLCHAIN=local", "-e", "GOSUMDB=off", "-v", f.cache + ":/go/pkg/mod:ro", "-v", linuxLaneFixtureObjects(t, f.repo) + ":/repo-objects:ro", "golang:1.26.1", "sh", "-ceu", script, "linux-lane"}
			want = append(want, tc.packages...)
			if run.Tool != "docker" || !reflect.DeepEqual(run.Args, want) {
				t.Fatalf("run args=%q\nwant=%q", run.Args, want)
			}
			for _, start := range starts {
				for _, entry := range f.environment {
					if !containsExact(start.Env, entry) {
						t.Fatalf("%s lost owner environment %q: %q", start.Tool, entry, start.Env)
					}
				}
				dir, err := filepath.EvalSymlinks(f.repo)
				if err != nil {
					t.Fatal(err)
				}
				if start.Dir != dir {
					t.Fatalf("dir=%q want=%q", start.Dir, dir)
				}
			}
			data, err := os.ReadFile(f.archive)
			if err != nil {
				t.Fatal(err)
			}
			archived := map[string]string{}
			tr := tar.NewReader(bytes.NewReader(data))
			for {
				header, err := tr.Next()
				if err == io.EOF {
					break
				}
				if err != nil {
					t.Fatal(err)
				}
				content, err := io.ReadAll(tr)
				if err != nil {
					t.Fatal(err)
				}
				if header.Typeflag == tar.TypeReg {
					archived[header.Name] = string(content)
				}
			}
			expected := map[string]string{}
			for k, v := range files {
				expected[k] = v
			}
			expected[".roundfixrc.yml"] = "delivery:\n  linux_lane: true\n"
			expected[tc.path] = tc.data
			if !reflect.DeepEqual(archived, expected) {
				t.Fatalf("archive=%v want=%v", archived, expected)
			}
			if after := gittest.Run(t, f.repo, "status", "--porcelain"); after != before {
				t.Fatalf("worktree changed: %q -> %q", before, after)
			}
			line := fmt.Sprintf("roundfix: Delivery Queue item widget: Linux Lane passed %d package(s) on golang:1.26.1\n", len(tc.packages))
			if f.log.String() != line {
				t.Fatalf("console=%q", f.log.String())
			}
			log := mustRead(t, f.logPath())
			if !strings.Contains(log, "docker 'run'") || !strings.Contains(log, "fixture output") || !strings.HasSuffix(log, strings.TrimPrefix(line, "roundfix: Delivery Queue item widget: ")) {
				t.Fatalf("lane log=%q", log)
			}
		})
	}
	t.Run("sorted nearest unique packages", func(t *testing.T) {
		f := newLinuxLaneFixture(t, docker, goBinary, files)
		f.commit(t, "z/testdata/input", "data")
		f.commit(t, "a/new.go", "package a")
		f.commit(t, "a/testdata/other", "other")
		f.commit(t, "docs/change.md", "docs")
		result, err := f.workflow.runLinuxLane(f.context(t.Context()), linuxLaneRequest{WorkDir: f.repo, Head: f.head, Label: "Delivery Queue item widget: ", LogPath: f.logPath(), Name: fmt.Sprintf("roundfix-linux-lane-widget-%d", os.Getpid())})
		if err != nil || !reflect.DeepEqual(result.Packages, []string{"./a", "./z"}) {
			t.Fatalf("result=%+v err=%v", result, err)
		}
	})
}

func containsExact(values []string, value string) bool {
	for _, entry := range values {
		if entry == value {
			return true
		}
	}
	return false
}

func TestLinuxLaneSkipsWithANamedReason(t *testing.T) {
	t.Parallel()
	docker, goBinary := linuxLaneTools(t)
	for _, tc := range []struct {
		name, reason string
		setup        func(*testing.T, *linuxLaneFixture)
	}{
		{"no module", "no go.mod at the repository root", func(t *testing.T, f *linuxLaneFixture) {
			if err := os.Remove(filepath.Join(f.repo, "go.mod")); err != nil {
				t.Fatal(err)
			}
			gittest.Run(t, f.repo, "add", ".")
			gittest.Run(t, f.repo, "commit", "-m", "remove module")
			f.head = strings.TrimSpace(gittest.Run(t, f.repo, "rev-parse", "HEAD"))
			f.environment = withEnvValue(f.environment, "PATH", t.TempDir())
		}},
		{"no version", "go.mod declares no go version", func(t *testing.T, f *linuxLaneFixture) {
			f.commit(t, "go.mod", "module fixture\n")
			f.environment = withEnvValue(f.environment, "PATH", t.TempDir())
		}},
		{"no package", "no changed Go package", func(t *testing.T, f *linuxLaneFixture) {
			gittest.Run(t, f.repo, "update-ref", "refs/remotes/origin/main", f.head)
			f.commit(t, "docs/readme", "docs")
			f.environment = withEnvValue(f.environment, "PATH", t.TempDir())
		}},
		{"no go", "go is not on PATH", func(t *testing.T, f *linuxLaneFixture) {
			f.environment = withEnvValue(f.environment, "PATH", t.TempDir())
		}},
		{"no docker", "docker is not on PATH", func(t *testing.T, f *linuxLaneFixture) {
			f.environment = withEnvValue(f.environment, "PATH", filepath.Dir(goBinary))
		}},
		{"daemon", "docker info failed: daemon unavailable", func(t *testing.T, f *linuxLaneFixture) {
			f.environment = withEnvValue(f.environment, "LANE_INFO_FAIL", "1")
			f.environment = withEnvValue(f.environment, "LANE_IMAGE_FAIL", "1")
		}},
		{"image", "image golang:1.26.1 is not present locally", func(t *testing.T, f *linuxLaneFixture) {
			f.environment = withEnvValue(f.environment, "LANE_IMAGE_FAIL", "1")
		}},
		{"start refused", "docker could not start the container: fixture output", func(t *testing.T, f *linuxLaneFixture) {
			f.environment = withEnvValue(f.environment, "LANE_EXIT", "125")
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newLinuxLaneFixture(t, docker, goBinary, map[string]string{"go.mod": "module fixture\n\ngo 1.26.1\n", "pkg/code.go": "package pkg\n"})
			f.commit(t, "pkg/changed.go", "package pkg\n")
			tc.setup(t, f)
			result, err := f.workflow.Gate(f.context(t.Context()), f.repo, "widget", f.head)
			if err != nil || !result.Passed {
				t.Fatalf("gate=%+v err=%v", result, err)
			}
			want := "roundfix: notice: Delivery Queue item widget: Linux Lane skipped: " + tc.reason + "\n"
			if f.log.String() != want {
				t.Fatalf("console=%q want=%q", f.log.String(), want)
			}
			if log := mustRead(t, f.logPath()); !strings.HasSuffix(log, want) {
				t.Fatalf("log=%q", log)
			}
			runs := 0
			for _, start := range f.starts(t) {
				if start.Tool == "docker" && start.Args[0] == "run" {
					runs++
				}
			}
			wantRuns := 0
			if tc.name == "start refused" {
				wantRuns = 1
			} // Exit 125 necessarily requires a run attempt, but starts no container.
			if runs != wantRuns {
				t.Fatalf("docker run attempts=%d want=%d", runs, wantRuns)
			}
		})
	}
}

func TestLinuxLaneFailureFailsTheGate(t *testing.T) {
	t.Parallel()
	docker, goBinary := linuxLaneTools(t)
	for _, code := range []string{"1", "2"} {
		t.Run(code, func(t *testing.T) {
			f := newLinuxLaneFixture(t, docker, goBinary, map[string]string{"go.mod": "module fixture\n\ngo 1.26\n", "pkg/code.go": "package pkg"})
			f.commit(t, "pkg/testdata/input", "data")
			f.environment = withEnvValue(f.environment, "LANE_EXIT", code)
			result, err := f.workflow.Gate(f.context(t.Context()), f.repo, "widget", f.head)
			want := "Linux Lane change failure: go test exited " + code + "; see " + f.logPath()
			if err != nil || result.Passed || result.Reason != want || result.Blocker != delivery.BlockerGateFailed+": "+want {
				t.Fatalf("gate=%+v err=%v want=%q", result, err, want)
			}
			if log := mustRead(t, f.logPath()); !strings.Contains(log, "fixture output") || !strings.HasSuffix(log, "Linux Lane change failure: go test exited "+code+"\n") {
				t.Fatal("missing failure output")
			}
		})
	}
}

func TestLinuxLaneCancellationRemovesTheContainer(t *testing.T) {
	t.Parallel()
	docker, goBinary := linuxLaneTools(t)
	f := newLinuxLaneFixture(t, docker, goBinary, map[string]string{"go.mod": "module fixture\n\ngo 1.26\n", "pkg/code.go": "package pkg"})
	f.commit(t, "pkg/testdata/input", "data")
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	f.environment = withEnvValue(f.environment, "LANE_READY", listener.Addr().String())
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	stop := context.AfterFunc(ctx, func() { _ = listener.Close() })
	defer stop()
	connected := make(chan net.Conn, 1)
	accepted := make(chan struct{})
	go func() {
		defer close(accepted)
		conn, err := listener.Accept()
		if err == nil {
			connected <- conn
		}
	}()
	type outcome struct {
		result delivery.GateResult
		err    error
	}
	done := make(chan outcome, 1)
	go func() {
		result, err := f.workflow.Gate(f.context(ctx), f.repo, "widget", f.head)
		done <- outcome{result, err}
	}()
	// The test owns both goroutines and cancels and joins them on every return.
	joined := false
	defer func() {
		cancel()
		if !joined {
			<-done
		}
		<-accepted
		select {
		case conn := <-connected:
			_ = conn.Close()
		default:
		}
	}()
	select {
	case conn := <-connected:
		defer conn.Close()
		cancel()
	case result := <-done:
		joined = true
		t.Fatalf("gate returned before docker run blocked: %+v", result)
	case <-t.Context().Done():
		t.Fatal(t.Context().Err())
	}
	result := <-done
	joined = true
	if !errors.Is(result.err, context.Canceled) {
		t.Fatalf("gate=%+v err=%v", result.result, result.err)
	}
	starts := f.starts(t)
	want := []string{"rm", "-f", fmt.Sprintf("roundfix-linux-lane-widget-%d", os.Getpid())}
	if len(starts) != 5 || !reflect.DeepEqual(starts[4].Args, want) {
		t.Fatalf("starts=%+v want cleanup=%v", starts, want)
	}
}

func TestGateWithoutTheLaneStartsNoDocker(t *testing.T) {
	t.Parallel()
	docker, goBinary := linuxLaneTools(t)
	for _, tc := range []struct {
		name         string
		enabled      bool
		verification string
		passed       bool
	}{
		{"off", false, "printf 'repository passed'", true},
		{"repository fails first", true, "exit 9", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newLinuxLaneFixture(t, docker, goBinary, nil)
			f.workflow.loaded.Config.Delivery.LinuxLane = !tc.enabled // The fetched default decides.
			if !tc.enabled {
				f.commit(t, ".roundfixrc.yml", "delivery:\n  linux_lane: false\n")
				gittest.Run(t, f.repo, "update-ref", "refs/remotes/origin/main", f.head)
			}
			f.workflow.loaded.Config.Defaults.Verification = tc.verification
			result, err := f.workflow.Gate(f.context(t.Context()), f.repo, "widget", f.head)
			if err != nil || result.Passed != tc.passed || result.Blocker != "" {
				t.Fatalf("gate=%+v err=%v", result, err)
			}
			if starts := f.starts(t); len(starts) != 0 {
				t.Fatalf("started tools=%+v", starts)
			}
			if f.log.Len() != 0 {
				t.Fatalf("console=%q", f.log.String())
			}
			if _, err := os.Stat(f.logPath()); !os.IsNotExist(err) {
				t.Fatalf("lane log created: %v", err)
			}
		})
	}
}

func TestGateSkipsTheLaneWhenTheFetchedDefaultCannotBeRead(t *testing.T) {
	t.Parallel()
	docker, goBinary := linuxLaneTools(t)
	f := newLinuxLaneFixture(t, docker, goBinary, nil)
	f.workflow.loaded.Config.Delivery.LinuxLane = true // Only the fetched default could enable it.
	gittest.Run(t, f.repo, "update-ref", "-d", "refs/remotes/origin/main")
	f.workflow.loaded.Config.Defaults.Verification = "printf 'repository passed'"
	result, err := f.workflow.Gate(f.context(t.Context()), f.repo, "widget", f.head)
	if err != nil || !result.Passed {
		t.Fatalf("gate=%+v err=%v", result, err)
	}
	if starts := f.starts(t); len(starts) != 0 {
		t.Fatalf("started tools=%+v", starts)
	}
	want := "roundfix: notice: Delivery Queue item widget: Linux Lane skipped: the fetched default refs/remotes/origin/main cannot be resolved\n"
	if f.log.String() != want {
		t.Fatalf("console=%q want=%q", f.log.String(), want)
	}
	logged, err := os.ReadFile(f.logPath())
	if err != nil || string(logged) != want {
		t.Fatalf("lane log=%q err=%v want=%q", logged, err, want)
	}
}

func TestLinuxLaneOptInPropagatesCancellationAndAnUnreadableConfig(t *testing.T) {
	t.Parallel()
	docker, goBinary := linuxLaneTools(t)
	t.Run("cancelled", func(t *testing.T) {
		f := newLinuxLaneFixture(t, docker, goBinary, nil)
		ctx, cancel := context.WithCancel(f.context(t.Context()))
		cancel()
		enabled, image, reason, err := f.workflow.linuxLaneOptIn(ctx, f.repo)
		if !errors.Is(err, context.Canceled) || enabled || image != "" || reason != "" {
			t.Fatalf("enabled=%t reason=%q err=%v, want context.Canceled", enabled, reason, err)
		}
	})
	t.Run("unreadable config", func(t *testing.T) {
		f := newLinuxLaneFixture(t, docker, goBinary, nil)
		f.commit(t, ".roundfixrc.yml", "delivery: [not a mapping\n")
		gittest.Run(t, f.repo, "update-ref", "refs/remotes/origin/main", f.head)
		enabled, image, reason, err := f.workflow.linuxLaneOptIn(f.context(t.Context()), f.repo)
		if err == nil || enabled || image != "" || reason != "" {
			t.Fatalf("enabled=%t reason=%q err=%v, want a configuration error", enabled, reason, err)
		}
	})
}

func TestLinuxLaneRunsTheConfiguredImage(t *testing.T) {
	t.Parallel()
	docker, goBinary := linuxLaneTools(t)
	const image = "registry.example.com/team/go-ci:1.26"
	for _, absent := range []bool{false, true} {
		t.Run(fmt.Sprintf("absent=%t", absent), func(t *testing.T) {
			f := newLinuxLaneFixture(t, docker, goBinary, map[string]string{
				".roundfixrc.yml": "delivery:\n  linux_lane: true\n  linux_lane_image: " + image + "\n",
				"go.mod":          "module fixture\n", // Configured images do not need a go version.
				"a/a.go":          "package a\n",
			})
			f.workflow.loaded.Config.Delivery.LinuxLaneImage = "owner/image:wrong"
			f.commit(t, ".roundfixrc.yml", "delivery:\n  linux_lane: false\n  linux_lane_image: item/image:wrong\n")
			f.commit(t, "a/new.go", "package a\n")
			if absent {
				f.environment = append(f.environment, "LANE_IMAGE_FAIL=1")
			}
			result, err := f.workflow.Gate(f.context(t.Context()), f.repo, "widget", f.head)
			if err != nil || !result.Passed {
				t.Fatalf("gate=%+v err=%v", result, err)
			}
			starts := f.starts(t)
			if len(starts) < 2 || !reflect.DeepEqual(starts[1].Args, []string{"image", "inspect", image}) {
				t.Fatalf("starts=%+v", starts)
			}
			if absent {
				if len(starts) != 2 || !strings.Contains(f.log.String(), "image "+image+" is not present locally") {
					t.Fatalf("starts=%+v log=%s", starts, f.log.String())
				}
				return
			}
			if len(starts) != 4 {
				t.Fatalf("starts=%+v", starts)
			}
			args := starts[3].Args
			if len(args) < 4 || !reflect.DeepEqual(args[:4], []string{"run", "--rm", "-i", "--init"}) || !containsExact(args, image) || containsExact(args, "item/image:wrong") || !strings.Contains(f.log.String(), "Linux Lane passed 1 package(s) on "+image) {
				t.Fatalf("args=%q log=%s", args, f.log.String())
			}
		})
	}
}

func TestLinuxLaneContainerFailureIsToldApartFromAChangeFailure(t *testing.T) {
	t.Parallel()
	docker, goBinary := linuxLaneTools(t)
	for _, tc := range []struct{ name, code, output, cause string }{
		{"setup marker", "3", "roundfix-linux-lane: container: missing tool tar\nroundfix-linux-lane: container: missing tool git\nother output\n", "missing tool git"},
		{"setup without trailing newline", "3", "roundfix-linux-lane: container: cannot record the candidate", "cannot record the candidate"},
		{"setup without marker", "3", "other output\n", "the lane script exited 3 before go test"},
		{"cannot invoke", "126", "", "docker run exited 126: the image could not invoke the lane command"},
		{"command absent", "127", "", "docker run exited 127: the image lacks the lane command"},
		{"signal", "137", "", "docker run exited 137: the container process ended by signal 9"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newLinuxLaneFixture(t, docker, goBinary, map[string]string{"go.mod": "module fixture\n\ngo 1.26\n", "pkg/code.go": "package pkg"})
			f.commit(t, "pkg/testdata/input", "data")
			f.environment = withEnvValue(f.environment, "LANE_EXIT", tc.code)
			f.environment = withEnvValue(f.environment, "LANE_OUTPUT", tc.output)
			result, err := f.workflow.Gate(f.context(t.Context()), f.repo, "widget", f.head)
			summary := "Linux Lane container failure: " + tc.cause
			want := summary + "; see " + f.logPath()
			blocker := delivery.BlockerLinuxLaneContainer + ": " + tc.cause + "; see " + f.logPath()
			if err != nil || result.Passed || result.Reason != want || result.Blocker != blocker {
				t.Fatalf("gate=%+v err=%v want reason=%q blocker=%q", result, err, want, blocker)
			}
			if log := mustRead(t, f.logPath()); !strings.HasSuffix(log, summary+"\n") {
				t.Fatalf("lane log=%q want last line=%q", log, summary)
			}
		})
	}
}

func TestLinuxLaneOutcomeClassification(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name                     string
		code                     int
		capped                   bool
		output, summary, blocker string
	}{
		{"pass", 0, false, "", "", ""},
		{"test failure", 1, false, "roundfix-linux-lane: container: ignored\n", "Linux Lane change failure: go test exited 1", delivery.BlockerGateFailed},
		{"usage", 2, false, "", "Linux Lane change failure: go test exited 2", delivery.BlockerGateFailed},
		{"other change exit", 4, false, "", "Linux Lane change failure: go test exited 4", delivery.BlockerGateFailed},
		{"setup", 3, false, "roundfix-linux-lane: container: missing tool go\n", "Linux Lane container failure: missing tool go", delivery.BlockerLinuxLaneContainer},
		{"last marker", 3, false, "roundfix-linux-lane: container: missing tool tar\nroundfix-linux-lane: container: cannot unpack the candidate\nother output\n", "Linux Lane container failure: cannot unpack the candidate", delivery.BlockerLinuxLaneContainer},
		{"no marker", 3, false, "unrelated output", "Linux Lane container failure: the lane script exited 3 before go test", delivery.BlockerLinuxLaneContainer},
		{"embedded marker ignored", 3, false, "unrelated roundfix-linux-lane: container: ignored\n", "Linux Lane container failure: the lane script exited 3 before go test", delivery.BlockerLinuxLaneContainer},
		{"invoke", 126, false, "", "Linux Lane container failure: docker run exited 126: the image could not invoke the lane command", delivery.BlockerLinuxLaneContainer},
		{"missing command", 127, false, "", "Linux Lane container failure: docker run exited 127: the image lacks the lane command", delivery.BlockerLinuxLaneContainer},
		{"signal zero", 128, false, "", "Linux Lane container failure: docker run exited 128: the container process ended by signal 0", delivery.BlockerLinuxLaneContainer},
		{"interrupt", 130, false, "", "Linux Lane container failure: docker run exited 130: the container process ended by signal 2", delivery.BlockerLinuxLaneContainer},
		{"killed", 137, false, "", "Linux Lane container failure: docker run exited 137: the container process ended by signal 9", delivery.BlockerLinuxLaneContainer},
		{"cap takes precedence", 3, true, "roundfix-linux-lane: container: missing tool go\n", "Linux Lane container failure: the container exceeded 20m0s", delivery.BlockerLinuxLaneContainer},
		{"cap after kill", -1, true, "", "Linux Lane container failure: the container exceeded 20m0s", delivery.BlockerLinuxLaneContainer},
	} {
		t.Run(tc.name, func(t *testing.T) {
			summary, blocker := classifyLinuxLaneOutcome(tc.code, tc.capped, []byte(tc.output))
			if summary != tc.summary || blocker != tc.blocker {
				t.Fatalf("summary=%q blocker=%q want %q %q", summary, blocker, tc.summary, tc.blocker)
			}
		})
	}
}

func linuxLaneFixtureObjects(t *testing.T, repo string) string {
	t.Helper()
	objects, err := filepath.EvalSymlinks(filepath.Join(repo, ".git", "objects"))
	if err != nil {
		t.Fatal(err)
	}
	return objects
}

func assertLinuxLaneHistoryMount(t *testing.T, f *linuxLaneFixture, objects string) {
	t.Helper()
	var runs int
	for _, start := range f.starts(t) {
		if start.Tool != "docker" || len(start.Args) == 0 || start.Args[0] != "run" {
			continue
		}
		runs++
		var mounts []string
		for i, arg := range start.Args {
			if arg == "-v" {
				if i+1 == len(start.Args) {
					t.Fatal("mount lacks value")
				}
				mounts = append(mounts, start.Args[i+1])
			}
		}
		want := []string{f.cache + ":/go/pkg/mod:ro", objects + ":/repo-objects:ro"}
		if !reflect.DeepEqual(mounts, want) {
			t.Fatalf("mounts=%q want=%q", mounts, want)
		}
		worktree, err := filepath.EvalSymlinks(f.repo)
		if err != nil {
			t.Fatal(err)
		}
		for _, mount := range mounts {
			source := strings.TrimSuffix(strings.TrimSuffix(mount, ":ro"), ":/repo-objects")
			source = strings.TrimSuffix(source, ":/go/pkg/mod")
			if source == worktree {
				t.Fatalf("working tree mounted: %s", mount)
			}
		}
		scriptAt := -1
		for i, arg := range start.Args {
			if arg == "-ceu" {
				scriptAt = i + 1
				break
			}
		}
		if scriptAt < 0 || scriptAt >= len(start.Args) {
			t.Fatal("missing lane script")
		}
		script := start.Args[scriptAt]
		for _, text := range []string{
			`[ -d /repo-objects ] && [ -r /repo-objects ] && [ -x /repo-objects ] || fail "cannot read the repository history at /repo-objects"`,
			`git init -q && echo /repo-objects > .git/objects/info/alternates && git add -A`,
		} {
			if !strings.Contains(script, text) {
				t.Fatalf("script missing %q: %s", text, script)
			}
		}
	}
	if runs != 1 {
		t.Fatalf("docker runs=%d want=1", runs)
	}
}

func TestLinuxLaneMountsTheRepositoryHistoryReadOnly(t *testing.T) {
	t.Parallel()
	docker, goBinary := linuxLaneTools(t)
	for _, mode := range []string{"linked item gate", "on demand all", "on demand changed"} {
		t.Run(mode, func(t *testing.T) {
			f := newLinuxLaneFixture(t, docker, goBinary, map[string]string{"go.mod": "module fixture\n\ngo 1.26.1\n", "pkg/code.go": "package pkg\n"})
			f.commit(t, "pkg/new.go", "package pkg\n")
			objects := linuxLaneFixtureObjects(t, f.repo)
			if mode == "linked item gate" {
				mainRepo := f.repo
				linked := filepath.Join(t.TempDir(), "linked item")
				gittest.Run(t, mainRepo, "worktree", "add", "--detach", linked, f.head)
				t.Cleanup(func() { gittest.Run(t, mainRepo, "worktree", "remove", "--force", linked) })
				f.repo = linked
				result, err := f.workflow.Gate(f.context(t.Context()), f.repo, "widget", f.head)
				if err != nil || !result.Passed {
					t.Fatalf("gate=%+v err=%v", result, err)
				}
			} else {
				var args []string
				if mode == "on demand all" {
					args = []string{"--all"}
				}
				stdout, stderr, code := invokeDeliverLane(t, f, args...)
				if code != 0 || stderr != "" {
					t.Fatalf("stdout=%q stderr=%q exit=%d", stdout, stderr, code)
				}
			}
			assertLinuxLaneHistoryMount(t, f, objects)
		})
	}
}

func TestLinuxLaneSkipsAnObjectStoreItCannotMount(t *testing.T) {
	t.Parallel()
	docker, goBinary := linuxLaneTools(t)
	for _, commentOnly := range []bool{false, true} {
		name := "borrowed store"
		if commentOnly {
			name = "comment and blank line"
		}
		t.Run(name, func(t *testing.T) {
			f := newLinuxLaneFixture(t, docker, goBinary, map[string]string{"go.mod": "module fixture\n\ngo 1.26.1\n", "pkg/code.go": "package pkg\n"})
			f.commit(t, "pkg/new.go", "package pkg\n")
			objects := linuxLaneFixtureObjects(t, f.repo)
			other := filepath.Join(t.TempDir(), "other")
			gittest.InitRepo(t, other, "-b", "main")
			data := linuxLaneFixtureObjects(t, other) + "\n"
			if commentOnly {
				data = "# no borrowed stores\n\n"
			}
			mustWrite(t, filepath.Join(objects, "info", "alternates"), data)
			result, err := f.workflow.Gate(f.context(t.Context()), f.repo, "widget", f.head)
			if err != nil || !result.Passed {
				t.Fatalf("gate=%+v err=%v", result, err)
			}
			if commentOnly {
				assertLinuxLaneHistoryMount(t, f, objects)
			} else {
				line := "roundfix: notice: Delivery Queue item widget: Linux Lane skipped: the object store " + objects + " borrows objects through info/alternates, which the Linux Lane does not mount\n"
				if f.log.String() != line || mustRead(t, f.logPath()) != line {
					t.Fatalf("console=%q log=%q want=%q", f.log.String(), mustRead(t, f.logPath()), line)
				}
				if starts := f.starts(t); len(starts) != 0 {
					t.Fatalf("skip started tools: %+v", starts)
				}
			}
		})
	}
	for _, kind := range []string{"colon", "file", "missing"} {
		t.Run(kind, func(t *testing.T) {
			objects := filepath.Join(t.TempDir(), "objects")
			cause := "it is not a directory"
			if kind == "colon" {
				objects += ":store"
				cause = "its path holds a colon"
			}
			if kind == "file" {
				mustWrite(t, objects, "file")
			}
			mount, skip, err := linuxLaneHistoryMount(objects)
			want := "the object store " + objects + " cannot be mounted: " + cause
			if mount != "" || skip != want || err != nil {
				t.Fatalf("mount=%q skip=%q err=%v want=%q", mount, skip, err, want)
			}
			if kind != "file" {
				if _, err := os.Stat(objects); !os.IsNotExist(err) {
					t.Fatalf("skip created source: %v", err)
				}
			}
		})
	}
}
