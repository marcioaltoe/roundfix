package daemon

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
	"roundfix/internal/store"
)

type priorQAPass struct {
	Commit string
	Head   string
	Report string
	Files  []string
}

// findPriorQAPass reads the installed Run Database without creating or migrating it.
func findPriorQAPass(ctx context.Context, workDir, specRelDir, slug string) (priorQAPass, bool, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return priorQAPass{}, false, err
	}
	journal, err := store.OpenReader(ctx, home)
	if errors.Is(err, os.ErrNotExist) {
		return priorQAPass{}, false, nil
	}
	if err != nil {
		return priorQAPass{}, false, err
	}
	defer journal.Close()
	return findRecordedPriorQAPass(ctx, journal, workDir, specRelDir, slug)
}

func priorQAGit(ctx context.Context, root string, args ...string) ([]byte, error) {
	output, err := exec.CommandContext(ctx, "git", append([]string{"-C", root}, args...)...).Output()
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, fmt.Errorf("read prior QA Git %s: %w", strings.Join(args, " "), gitExecStderr(err))
	}
	return output, nil
}

func priorQAAncestor(ctx context.Context, root, commit, head string) (bool, error) {
	_, err := priorQAGit(ctx, root, "merge-base", "--is-ancestor", commit, head)
	if err == nil {
		return true, nil
	}
	var exit *exec.ExitError
	if errors.As(err, &exit) && exit.ExitCode() == 1 && ctx.Err() == nil {
		return false, nil
	}
	return false, err
}

func findRecordedPriorQAPass(ctx context.Context, journal runEventJournal, workDir, specRelDir, slug string) (priorQAPass, bool, error) {
	runs, err := journal.ListRuns(ctx, store.ListRunsQuery{GitRoot: workDir, States: store.StatesAll})
	if err != nil {
		return priorQAPass{}, false, err
	}
	recorded := map[string]bool{}
	qaDir := strings.TrimSuffix(filepath.ToSlash(specRelDir), "/") + "/qa/"
	for _, run := range runs {
		if run.Kind != store.KindImplement || run.SpecSlug != slug {
			continue
		}
		branch := store.RunBranchPrefix + run.ID
		refs, err := priorQAGit(ctx, workDir, "for-each-ref", "--format=%(refname)", "refs/heads/"+branch)
		if err != nil {
			return priorQAPass{}, false, err
		}
		present := false
		for _, ref := range strings.Fields(string(refs)) {
			if ref == "refs/heads/"+branch {
				present = true
				break
			}
		}
		if !present {
			continue
		}
		var cursor int64
		for {
			events, err := journal.RunEventsAfter(ctx, run.ID, cursor, repeatedFailureJournalPageSize)
			if err != nil {
				return priorQAPass{}, false, err
			}
			for _, entry := range events {
				cursor = entry.Cursor
				if entry.Event.Kind != runevent.KindDaemonCommit || entry.Event.Source != runevent.SourceDaemon {
					continue
				}
				var payload struct {
					Decision string `json:"decision"`
					Commit   string `json:"commit"`
					Report   string `json:"report"`
					Task     string `json:"task"`
				}
				if err := json.Unmarshal(entry.Event.Payload, &payload); err != nil {
					return priorQAPass{}, false, fmt.Errorf("read QA settlement journal for Run %s: %w", run.ID, err)
				}
				if payload.Decision != "created" || payload.Commit == "" || payload.Task == "" || !strings.HasPrefix(payload.Report, qaDir) {
					continue
				}
				// Read commits only through their recorded Run Branch. A pruned commit
				// will not appear in this branch's history even if another ref holds it.
				history, err := priorQAGit(ctx, workDir, "rev-list", branch)
				if err != nil {
					return priorQAPass{}, false, err
				}
				for _, commit := range strings.Fields(string(history)) {
					if commit == payload.Commit {
						recorded[commit] = true
						break
					}
				}
			}
			if len(events) < repeatedFailureJournalPageSize {
				break
			}
		}
	}
	// Date order across all refs, restricted to journal-proven candidates.
	history, err := priorQAGit(ctx, workDir, "log", "--all", "--format=%ct %H")
	if err != nil {
		return priorQAPass{}, false, err
	}
	var commits []struct {
		sha  string
		date int64
	}
	for _, line := range strings.Split(strings.TrimSpace(string(history)), "\n") {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}
		date, err := strconv.ParseInt(fields[0], 10, 64)
		if err != nil {
			return priorQAPass{}, false, fmt.Errorf("read QA commit date: %w", err)
		}
		commits = append(commits, struct {
			sha  string
			date int64
		}{fields[1], date})
	}
	sort.SliceStable(commits, func(i, j int) bool { return commits[i].date > commits[j].date })
	for _, candidate := range commits {
		commit := candidate.sha
		if !recorded[commit] {
			continue
		}
		message, err := priorQAGit(ctx, workDir, "show", "-s", "--format=%s%n%(trailers)", commit)
		if err != nil {
			return priorQAPass{}, false, err
		}
		lines := strings.Split(string(message), "\n")
		matchesSpec, hasTask := false, false
		for _, line := range lines[1:] {
			key, value, ok := strings.Cut(line, ":")
			if !ok {
				continue
			}
			if strings.EqualFold(key, "Roundfix-Spec") && strings.TrimSpace(value) == slug {
				matchesSpec = true
			}
			if strings.EqualFold(key, "Roundfix-Task") {
				hasTask = true
			}
		}
		if !strings.HasPrefix(lines[0], "docs: qa report for "+slug+" (") || !matchesSpec || hasTask {
			continue
		}
		ancestor, err := priorQAAncestor(ctx, workDir, commit, "HEAD")
		if err != nil {
			return priorQAPass{}, false, err
		}
		if ancestor {
			continue
		}
		parent, err := priorQAGit(ctx, workDir, "rev-parse", commit+"^1")
		if err != nil {
			return priorQAPass{}, false, err
		}
		paths, err := priorQAGit(ctx, workDir, "diff-tree", "--no-renames", "--diff-filter=AM", "--name-only", "-r", "-z", strings.TrimSpace(string(parent)), commit, "--", qaDir)
		if err != nil {
			return priorQAPass{}, false, err
		}
		pass := priorQAPass{Commit: commit, Head: strings.TrimSpace(string(parent))}
		var reports []string
		for _, path := range strings.Split(string(paths), "\x00") {
			if !strings.HasPrefix(path, qaDir) {
				continue
			}
			pass.Files = append(pass.Files, path)
			if strings.HasPrefix(filepath.Base(path), "qa-report-") && strings.HasSuffix(path, ".md") {
				reports = append(reports, path)
			}
		}
		if len(reports) == 0 {
			continue
		}
		pass.Report, err = spec.NewestQAReportFromPaths(reports)
		if err != nil {
			return priorQAPass{}, false, err
		}
		return pass, true, nil
	}
	return priorQAPass{}, false, nil
}

func (engine *Engine) importPriorQAPass(ctx context.Context, plan TaskPlan, ordinal int) (pass priorQAPass, imported bool, err error) {
	defer func() {
		if ctx.Err() != nil {
			err = errors.Join(err, ctx.Err(), engine.publishStop(ctx, plan.RunID, ordinal))
		}
	}()
	if err := ctx.Err(); err != nil {
		return priorQAPass{}, false, err
	}
	outcome, reason := "none", ""
	rel, err := filepath.Rel(plan.WorkDir, plan.Spec.Dir)
	if err != nil {
		return pass, false, err
	}
	if rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		journal, ok := engine.deps.Runs.(runEventJournal)
		if ok {
			var found bool
			pass, found, err = findRecordedPriorQAPass(ctx, journal, plan.WorkDir, rel, plan.Spec.Slug)
			if err != nil {
				return pass, false, err
			}
			if found {
				reason, err = engine.copyPriorQAPass(ctx, plan, pass)
				if err != nil {
					return pass, false, err
				}
				outcome = "refused"
				if reason == "" {
					outcome, imported = "imported", true
				}
			}
		}
	}
	err = engine.publishDaemonEvent(ctx, plan.RunID, ordinal, runevent.KindDaemonQA,
		fmt.Sprintf("Prior QA Report %s for Spec %s.", outcome, plan.Spec.Slug),
		map[string]any{"phase": "prior_report", "outcome": outcome, "commit": pass.Commit, "report": pass.Report, "files": pass.Files, "reason": reason})
	if !imported {
		pass = priorQAPass{}
	}
	return pass, imported, err
}

func (engine *Engine) copyPriorQAPass(ctx context.Context, plan TaskPlan, pass priorQAPass) (reason string, err error) {
	root := filepath.Clean(plan.WorkDir)
	newest, err := spec.NewestQAReport(plan.Spec.Dir)
	if err != nil && !errors.Is(err, spec.ErrNoQAReport) {
		return "", err
	}
	if err == nil {
		newer, err := spec.NewestQAReportFromPaths([]string{newest, pass.Report})
		if err != nil {
			return "", err
		}
		if filepath.Base(newer) != filepath.Base(pass.Report) || filepath.Base(newest) == filepath.Base(pass.Report) {
			return "not newer than " + taskPromptPath(plan, newest), nil
		}
	}
	name := strings.TrimPrefix(filepath.Base(pass.Report), "qa-report-")
	if len(name) >= 10 && name[:10] > engine.deps.Now().Format("2006-01-02") {
		return "dated after today", nil
	}
	blobs := map[string][]byte{}
	var created []string
	// Refusals never overwrite a pre-existing file, including symlink targets.
	for _, path := range pass.Files {
		blob, err := priorQAGit(ctx, plan.WorkDir, "cat-file", "blob", pass.Commit+":"+path)
		if err != nil {
			return "", err
		}
		full := filepath.Join(root, filepath.FromSlash(path))
		// Reject symlink ancestors so Git paths cannot write outside the Worktree.
		for parent := filepath.Dir(full); parent != root; parent = filepath.Dir(parent) {
			info, err := os.Lstat(parent)
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				return "", err
			}
			if err == nil && info.Mode()&os.ModeSymlink != 0 {
				return "path differs: " + path, nil
			}
		}
		info, statErr := os.Lstat(full)
		if statErr == nil {
			if !info.Mode().IsRegular() {
				return "path differs: " + path, nil
			}
			existing, err := os.ReadFile(full)
			if err != nil {
				return "", err
			}
			if !bytes.Equal(existing, blob) {
				return "path differs: " + path, nil
			}
		} else if errors.Is(statErr, os.ErrNotExist) {
			blobs[path] = blob
		} else {
			return "", statErr
		}
	}
	var directories []string
	defer func() {
		if reason != "" || err != nil {
			for _, path := range created {
				err = errors.Join(err, os.Remove(path))
			}
			for i := len(directories) - 1; i >= 0; i-- {
				err = errors.Join(err, os.Remove(directories[i]))
			}
		}
	}()
	for _, path := range pass.Files {
		blob, needed := blobs[path]
		if !needed {
			continue
		}
		if err := ctx.Err(); err != nil {
			return "", err
		}
		full := filepath.Join(root, filepath.FromSlash(path))
		var missing []string
		for dir := filepath.Dir(full); dir != root; dir = filepath.Dir(dir) {
			if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
				missing = append(missing, dir)
			} else if err != nil {
				return "", err
			} else {
				break
			}
		}
		for i := len(missing) - 1; i >= 0; i-- {
			if err := os.Mkdir(missing[i], 0o755); err != nil {
				return "", err
			}
			directories = append(directories, missing[i])
		}
		file, err := os.OpenFile(full, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if err != nil {
			return "", err
		}
		created = append(created, full)
		_, writeErr := file.Write(blob)
		modeErr := file.Chmod(0o644)
		closeErr := file.Close()
		if err := errors.Join(writeErr, modeErr, closeErr); err != nil {
			return "", err
		}
	}
	findings, err := speccheck.ReportShapeFindings(plan.WorkDir, pass.Report)
	if err != nil {
		return "", err
	}
	if len(findings) > 0 {
		return "report shape: " + findings[0].Code, nil
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}
	return "", nil
}
