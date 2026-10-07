package daemon

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"roundfix/internal/runevent"
	"roundfix/internal/spec"
)

const qaFormatTimeout = 5 * time.Minute

// qaFormatDiagnostics retains only the diagnostic tail, even for noisy commands.
type qaFormatDiagnostics struct{ tail []byte }

func (d *qaFormatDiagnostics) Write(p []byte) (int, error) {
	n := len(p)
	if len(p) >= 2048 {
		d.tail = append(d.tail[:0], p[len(p)-2048:]...)
	} else {
		d.tail = append(d.tail, p...)
		if len(d.tail) > 2048 {
			d.tail = d.tail[len(d.tail)-2048:]
		}
	}
	return n, nil
}

// formatQADirectory formats the supplied QA files without widening the commit set.
func (engine *Engine) formatQADirectory(ctx context.Context, plan TaskPlan, ordinal int, stage string, paths []string, verdict string) error {
	if strings.TrimSpace(plan.FormatCommand) == "" {
		return nil
	}
	root, err := filepath.Abs(plan.WorkDir)
	if err != nil {
		return fmt.Errorf("resolve QA format root: %w", err)
	}
	qaDir := filepath.Join(plan.Spec.Dir, "qa")
	if !filepath.IsAbs(qaDir) {
		qaDir = filepath.Join(root, qaDir)
	}
	type snapshot struct {
		full    string
		content []byte
		mode    os.FileMode
		verdict string
	}
	files := map[string]snapshot{}
	for _, path := range paths {
		full := path
		if !filepath.IsAbs(full) {
			full = filepath.Join(root, filepath.FromSlash(path))
		}
		full = filepath.Clean(full)
		rel, err := filepath.Rel(qaDir, full)
		if err != nil || rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			continue
		}
		// A regular leaf below a symlink directory is not a QA file in this Worktree.
		safe := true
		for parent := filepath.Dir(full); parent != root; parent = filepath.Dir(parent) {
			info, err := os.Lstat(parent)
			if err != nil || info.Mode()&os.ModeSymlink != 0 || parent == filepath.Dir(parent) {
				safe = false
				break
			}
		}
		if !safe {
			continue
		}
		info, err := os.Lstat(full)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return fmt.Errorf("inspect QA format file %q: %w", path, err)
		}
		if !info.Mode().IsRegular() {
			continue
		}
		content, err := os.ReadFile(full)
		if err != nil {
			return fmt.Errorf("snapshot QA format file %q: %w", path, err)
		}
		relative, err := filepath.Rel(root, full)
		if err != nil {
			return fmt.Errorf("resolve QA format path: %w", err)
		}
		file := snapshot{full: full, content: content, mode: info.Mode()}
		if strings.HasPrefix(filepath.Base(full), "qa-report-") && strings.HasSuffix(full, ".md") {
			report, err := spec.ReadQAReportFile(full)
			if err != nil {
				return fmt.Errorf("read QA format verdict for %q: %w", path, err)
			}
			file.verdict = report.Verdict
			// Older reports in the commit retain their own verdict; the newest report
			// is the report whose verdict the QA step just settled.
		}
		files[filepath.ToSlash(relative)] = file
	}
	if len(files) == 0 {
		return nil
	}
	names := make([]string, 0, len(files))
	for path := range files {
		names = append(names, path)
	}
	sort.Strings(names)
	if stage == "commit" && verdict != "" {
		var reports []string
		for _, name := range names {
			if files[name].verdict != "" {
				reports = append(reports, name)
			}
		}
		if len(reports) > 0 {
			newest, err := spec.NewestQAReportFromPaths(reports)
			if err != nil {
				return fmt.Errorf("resolve QA format report: %w", err)
			}
			file := files[newest]
			file.verdict = verdict
			files[newest] = file
		}
	}
	watchCtx, finish := engine.watchStopRequest(ctx, plan.RunID)
	runCtx, cancel := context.WithTimeout(watchCtx, qaFormatTimeout)
	diagnostics := &qaFormatDiagnostics{}
	args := append([]string{"-c", plan.FormatCommand + ` "$@"`, "roundfix-format"}, names...)
	command := exec.CommandContext(runCtx, "sh", args...)
	configureQAFormatProcess(command)
	command.Dir = root
	command.Stdout, command.Stderr = diagnostics, diagnostics
	command.WaitDelay = time.Second
	runErr := command.Run()
	runContextErr := runCtx.Err()
	cancel()
	stopped := finish()
	outcome, reason := "unchanged", ""
	changed := []string{}
	for _, name := range names {
		file := files[name]
		content, err := os.ReadFile(file.full)
		if err != nil || !bytes.Equal(content, file.content) {
			changed = append(changed, name)
		}
		if err != nil {
			outcome, reason = "reverted", "cannot read formatted file: "+name
		}
		if file.verdict != "" {
			report, err := spec.ReadQAReportFile(file.full)
			if err != nil || report.Verdict != file.verdict {
				outcome, reason = "reverted", "QA Report verdict changed: "+name
			}
		}
	}
	if outcome == "unchanged" && len(changed) > 0 {
		outcome = "formatted"
	}
	if runErr != nil {
		outcome, reason = "failed", terminalReasonLine(runErr.Error())
	}
	if runContextErr != nil {
		outcome, reason = "failed", runContextErr.Error()
	}
	if stopped {
		outcome, reason = "reverted", ErrStopRequested.Error()
	}
	if ctx.Err() != nil {
		outcome, reason = "reverted", ctx.Err().Error()
	}
	var restoreErr error
	if outcome == "failed" || outcome == "reverted" {
		for _, name := range names {
			file := files[name]
			// Remove a replacement symlink rather than writing through it.
			err := os.Remove(file.full)
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				restoreErr = errors.Join(restoreErr, fmt.Errorf("remove formatted file %q: %w", name, err))
				continue
			}
			if err := os.WriteFile(file.full, file.content, file.mode.Perm()); err != nil {
				restoreErr = errors.Join(restoreErr, fmt.Errorf("restore QA format bytes %q: %w", name, err))
				continue
			}
			if err := os.Chmod(file.full, file.mode); err != nil {
				restoreErr = errors.Join(restoreErr, fmt.Errorf("restore QA format mode %q: %w", name, err))
			}
		}
		fmt.Fprintf(engine.deps.Progress, "roundfix: QA format %s %s: %s\n", stage, outcome, reason)
	}
	payload := map[string]any{"phase": "format", "stage": stage, "outcome": outcome, "command": plan.FormatCommand, "paths": len(names), "changed": changed}
	if outcome == "failed" || outcome == "reverted" {
		payload["reason"], payload["diagnostics"] = reason, string(diagnostics.tail)
	}
	eventErr := engine.publishDaemonEvent(context.WithoutCancel(ctx), plan.RunID, ordinal, runevent.KindDaemonQA, fmt.Sprintf("QA format %s %s.", stage, outcome), payload)
	var stopErr error
	if stopped || ctx.Err() != nil {
		stopErr = ctx.Err()
		if stopped {
			stopErr = errors.Join(stopErr, ErrStopRequested)
		}
		stopErr = errors.Join(stopErr, engine.publishStop(context.WithoutCancel(ctx), plan.RunID, ordinal))
	}
	return errors.Join(restoreErr, eventErr, stopErr)
}
