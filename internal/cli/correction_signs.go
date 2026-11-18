package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"roundfix/internal/baseline"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/correction"
	"roundfix/internal/delivery"
	"roundfix/internal/preflight"
	"roundfix/internal/runcause"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

var _ delivery.CorrectionSignGate = (*commandDeliveryWorkflow)(nil)

func (workflow *commandDeliveryWorkflow) correctionDefaultHead(ctx context.Context, root string) (preflight.DefaultBranchHead, bool, error) {
	branch, err := workflow.git.RunGit(ctx, root, "symbolic-ref", "--quiet", "--short", "HEAD")
	if err != nil {
		branch = ""
	} // Detached HEAD can still resolve origin/HEAD.
	return preflight.ResolveDefaultBranchHead(ctx, root, strings.TrimSpace(branch), workflow.git)
}

// CorrectionSigns is the single local-object reader used by the owner, status,
// and the per-Spec report. It neither fetches nor records its result.
func (workflow *commandDeliveryWorkflow) CorrectionSigns(ctx context.Context, gitRoot, workDir, slug, head, defaultHead string) (correction.Snapshot, error) {
	snapshot := correction.Snapshot{Head: head, Signs: []correction.Sign{}}
	all, err := workflow.git.RunGit(ctx, gitRoot, "rev-list", defaultHead+".."+head)
	if err != nil {
		return snapshot, fmt.Errorf("read item commits: %w", err)
	}
	members := map[string]bool{}
	for _, commit := range strings.Fields(all) {
		members[commit] = true
	}
	artifacts, err := roundconfig.ResolveArtifactDirectory(workflow.loaded.Config.Defaults.ArtifactDir, gitRoot, workflow.loaded.HomeDir)
	if err != nil {
		return snapshot, err
	}
	ledger, err := readReviewFindingDispositions(filepath.Join(artifacts, reviewDispositionLedgerFileName))
	if err != nil {
		return snapshot, err
	}
	fixed := map[string]bool{}
	seen := map[string]bool{}
	for _, disposition := range ledger {
		// The artifact directory scopes this ledger to the repository. Its
		// originating worktree may be gone while the reviewed head survives.
		if !members[disposition.HeadCommit] {
			continue
		}
		if disposition.Disposition == "dismissed" {
			snapshot.Dismissed++
			continue
		}
		if disposition.Disposition != "fixed" {
			continue
		}
		id := disposition.Finding + "@" + disposition.HeadCommit
		fixed[disposition.FixedBy] = true
		if !seen[id] {
			snapshot.Signs = append(snapshot.Signs, correction.Sign{Kind: "finding", ID: id})
			seen[id] = true
		}
	}
	commits, err := workflow.git.RunGit(ctx, gitRoot, "rev-list", "--first-parent", defaultHead+".."+head)
	if err != nil {
		return snapshot, fmt.Errorf("read first-parent item commits: %w", err)
	}
	source, destination, err := workflow.archivePaths(gitRoot, slug)
	if err != nil {
		return snapshot, err
	}
	for _, commit := range strings.Fields(commits) {
		if fixed[commit] {
			continue
		}
		metadata, err := workflow.git.RunGit(ctx, gitRoot, "show", "-s", "--format=%P%x00%B%x00%(trailers:key=Roundfix-Spec,keyonly)%x00%(trailers:key=Roundfix-Delivery,valueonly)", commit)
		if err != nil {
			return snapshot, fmt.Errorf("read commit metadata: %w", err)
		}
		fields := strings.SplitN(metadata, "\x00", 4)
		if len(fields) != 4 {
			return snapshot, fmt.Errorf("read commit metadata: expected four fields")
		}
		if strings.TrimSpace(fields[2]) != "" || containsDerivedMergeTrailer(fields[3]) {
			continue
		}
		parents := strings.Fields(fields[0])
		defaultMerge := len(parents) > 1
		for _, parent := range parents[min(1, len(parents)):] {
			if _, err := workflow.git.RunGit(ctx, gitRoot, "merge-base", "--is-ancestor", parent, defaultHead); err != nil {
				var exit interface{ ExitCode() int }
				if !errors.As(err, &exit) || exit.ExitCode() != 1 {
					return snapshot, fmt.Errorf("inspect merge parent: %w", err)
				}
				defaultMerge = false
			}
		}
		if defaultMerge {
			continue
		}
		if len(parents) == 1 {
			exact, err := workflow.correctionQueueCommitIsExact(ctx, gitRoot, parents[0], commit, source, destination)
			if err != nil {
				return snapshot, err
			}
			if exact {
				continue
			}
		}
		snapshot.Signs = append(snapshot.Signs, correction.Sign{Kind: "commit", ID: commit, Subject: correctionSubject(fields[1])})
	}
	graph, err := workflow.correctionTaskGraph(ctx, gitRoot, workDir, slug, head, source, destination)
	if err != nil {
		return snapshot, err
	}
	for _, task := range runcause.CorrectiveTasks(graph) {
		snapshot.Signs = append(snapshot.Signs, correction.Sign{Kind: "task", ID: task})
	}
	return snapshot, nil
}

func containsDerivedMergeTrailer(value string) bool {
	for _, line := range strings.Split(value, "\n") {
		if strings.TrimSpace(line) == "derived-merge" {
			return true
		}
	}
	return false
}

func correctionSubject(message string) string {
	subject, _, _ := strings.Cut(message, "\n")
	subject = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) || r == '\u2028' || r == '\u2029' {
			return ' '
		}
		if r == '"' {
			return '\''
		}
		return r
	}, subject)
	runes := []rune(subject)
	if len(runes) > 72 {
		return string(runes[:71]) + "…"
	}
	return subject
}

func (workflow *commandDeliveryWorkflow) correctionQueueCommitIsExact(ctx context.Context, root, parent, head, source, destination string) (bool, error) {
	changed, err := workflow.git.RunGit(ctx, root, "diff", "--no-renames", "--name-only", "-z", parent, head, "--")
	if err != nil {
		return false, fmt.Errorf("read queue commit paths: %w", err)
	}
	relevant := false
	for _, path := range nulPaths(changed) {
		if path == destination+".md" || strings.HasPrefix(path, destination+"/") {
			relevant = true
		}
	}
	if !relevant {
		return false, nil
	}
	exact, err := workflow.archiveCommitIsExact(ctx, root, parent, head, source, destination)
	if err != nil || exact {
		return exact, err
	}
	paths := nulPaths(changed)
	if len(paths) != 1 || paths[0] != destination+".md" {
		return false, nil
	}
	content, err := workflow.git.RunGit(ctx, root, "show", head+":"+destination+".md")
	if err != nil {
		return false, fmt.Errorf("read Pull Request Record: %w", err)
	}
	record, err := spec.ParseArchiveRecord([]byte(content))
	if err != nil || record.PullRequest == "" {
		return false, nil
	}
	return workflow.pullRequestRecordCommitIsExact(ctx, root, parent, head, destination+".md", record.PullRequest)
}

func (workflow *commandDeliveryWorkflow) correctionTaskGraph(ctx context.Context, root, workDir, slug, head, source, destination string) (spec.CauseGraph, error) {
	if workDir != "" {
		dir := filepath.Join(workDir, filepath.FromSlash(source))
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			return spec.ReadCauseGraph(filepath.Dir(dir), slug)
		} else if err != nil && !errors.Is(err, os.ErrNotExist) {
			return spec.CauseGraph{}, fmt.Errorf("inspect item Spec: %w", err)
		}
	}
	content, err := workflow.git.RunGit(ctx, root, "show", head+":"+destination+".md")
	if err != nil {
		return spec.CauseGraph{}, fmt.Errorf("read signs Archive Record: %w", err)
	}
	record, err := spec.ParseArchiveRecord([]byte(content))
	if err != nil {
		return spec.CauseGraph{}, err
	}
	return spec.ReadCauseGraphAt(record.Source, slug, func(path string) ([]byte, error) {
		value, err := workflow.git.RunGit(ctx, root, "show", record.SourceRevision+":"+path)
		return []byte(value), err
	})
}

func correctionSignLabels(signs []correction.Sign) []string {
	labels := make([]string, 0, len(signs))
	for _, sign := range signs {
		label := correctionSignNames([]correction.Sign{sign})
		if sign.Kind == "commit" {
			label += fmt.Sprintf(" %q", correctionSubject(sign.Subject))
		}
		labels = append(labels, label)
	}
	return labels
}

func (workflow *commandDeliveryWorkflow) correctionRecordPath(ctx context.Context, root, slug string) (string, error) {
	identity, err := baseline.InspectRepositoryIdentity(ctx, root, workflow.git)
	if err != nil {
		return "", err
	}
	return correction.Path(workflow.loaded.HomeDir, identity.Digest, slug), nil
}

func (workflow *commandDeliveryWorkflow) CheckCorrectionSigns(ctx context.Context, gitRoot, workDir, specSlug, head string) (delivery.CorrectionSignCheck, error) {
	check := delivery.CorrectionSignCheck{}
	defaultHead, found, err := workflow.correctionDefaultHead(ctx, gitRoot)
	if err != nil || !found {
		return check, nil
	}
	config, err := roundconfig.DeliveryConfigAtCommit(ctx, workflow.git, gitRoot, workflow.loaded.UserConfigPath, defaultHead.Commit)
	if err != nil {
		return delivery.CorrectionSignCheck{Required: true}, err
	}
	check.Required = config.Delivery.RequireDefectClass
	snapshot, readErr := workflow.CorrectionSigns(ctx, gitRoot, workDir, specSlug, head, defaultHead.Commit)
	if readErr != nil {
		if check.Required {
			return check, readErr
		}
		snapshot.Unavailable = readErr.Error()
	}
	path, err := workflow.correctionRecordPath(ctx, gitRoot, specSlug)
	if err == nil {
		snapshot.Time = time.Now().UTC()
		var file correction.File
		file, err = appendCorrectionSnapshot(ctx, path, specSlug, snapshot, check.Required)
		if err == nil && check.Required {
			check.Undeclared = correctionSignLabels(correction.Undeclared(snapshot.Signs, correction.Standing(file)))
		}
	}
	if err != nil && check.Required {
		return check, err
	}
	return check, nil
}

// The snapshot and declaration share a lock, so a declaration made while the
// owner runs is visible in the very merge check that appends this snapshot.
func appendCorrectionSnapshot(ctx context.Context, path, slug string, snapshot correction.Snapshot, readStanding bool) (file correction.File, resultErr error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return file, fmt.Errorf("create corrections directory: %w", err)
	}
	lock, err := os.OpenFile(filepath.Dir(path)+".lock", os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return file, fmt.Errorf("open corrections lock: %w", err)
	}
	defer func() { resultErr = errors.Join(resultErr, lock.Close()) }()
	if err := lockReviewDispositionFile(ctx, lock); err != nil {
		return file, err
	}
	defer func() { resultErr = errors.Join(resultErr, unlockReviewDispositionFile(lock)) }()
	if readStanding {
		file, err = correction.Read(path)
		if err != nil {
			return file, err
		}
	}
	line, err := correction.EncodeSnapshot(slug, snapshot)
	if err != nil {
		return file, err
	}
	output, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return file, fmt.Errorf("open corrections record: %w", err)
	}
	defer func() { resultErr = errors.Join(resultErr, output.Close()) }()
	n, err := output.Write(line)
	if err == nil && n != len(line) {
		err = io.ErrShortWrite
	}
	if err != nil {
		return file, fmt.Errorf("append signs snapshot: %w", err)
	}
	return file, nil
}

func printDeliveryCorrectionSigns(ctx context.Context, loaded roundconfig.Loaded, queue store.DeliveryQueue, out io.Writer) {
	git := commandDependenciesForContext(ctx).reviewSpecGitRunner
	if git == nil {
		git = preflight.ExecGitRunner{}
	}
	workflow := &commandDeliveryWorkflow{loaded: loaded, git: git}
	defaultHead, found, err := workflow.correctionDefaultHead(ctx, loaded.GitRoot)
	if err != nil || !found {
		return
	}
	config, err := roundconfig.DeliveryConfigAtCommit(ctx, git, loaded.GitRoot, loaded.UserConfigPath, defaultHead.Commit)
	if err != nil || !config.Delivery.RequireDefectClass {
		return
	}
	for _, item := range queue.Items {
		if item.Stage == store.DeliveryStageMerged || item.Branch == "" {
			continue
		}
		head, err := git.RunGit(ctx, loaded.GitRoot, "rev-parse", "--verify", "--end-of-options", "refs/heads/"+item.Branch+"^{commit}")
		if err != nil {
			continue
		}
		snapshot, err := workflow.CorrectionSigns(ctx, loaded.GitRoot, item.Worktree, item.SpecSlug, strings.TrimSpace(head), defaultHead.Commit)
		if err != nil {
			continue
		}
		path, err := workflow.correctionRecordPath(ctx, loaded.GitRoot, item.SpecSlug)
		if err != nil {
			continue
		}
		file, err := correction.Read(path)
		if err != nil {
			continue
		}
		undeclared := correction.Undeclared(snapshot.Signs, correction.Standing(file))
		if len(undeclared) > 0 {
			fmt.Fprintf(out, "Corrections: %s undeclared: %s; declare with roundfix corrections declare --spec %s\n", item.SpecSlug, strings.Join(correctionSignLabels(undeclared), ", "), item.SpecSlug)
		}
	}
}
