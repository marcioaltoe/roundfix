// Suite: review disposition revisions.
// Invariant: operator-supplied revisions are always passed to Git as revisions, never as options.
// Boundary IN: public review dispose command, real local Git, and Artifact Directory files.
// Boundary OUT: reviewer execution and disposition reuse by later review runs.
package cli

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
)

func TestDisposeFixedByOptionLikeValueIsRefusedAndWritesNothing(t *testing.T) {
	fixture := newReviewDispositionFixture(t)
	writeDispositionFindingsRecord(t, fixture, "- internal/cli/review.go:10: finding", true)
	outputPath := filepath.Join(t.TempDir(), "git-output")
	optionLikeRevision := "--output=" + outputPath
	withReviewSpecGitRunner(t, optionWritingGitRunner{delegate: preflight.ExecGitRunner{}})

	code, stdout, stderr := runReviewDispose(t, "F1", "--fixed-by", optionLikeRevision)

	if code != exitPreflight {
		t.Fatalf("review dispose exit=%d stdout=%q stderr=%q, want %d", code, stdout, stderr, exitPreflight)
	}
	wantStderr := "roundfix: review dispose refused: fixed-by commit \"" + optionLikeRevision + "\" does not resolve\n"
	if stdout != "" || stderr != wantStderr {
		t.Fatalf("review dispose stdout=%q stderr=%q, want empty stdout and stderr %q", stdout, stderr, wantStderr)
	}
	if _, err := os.Stat(outputPath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("option-like revision wrote %q: %v", outputPath, err)
	}
	ledgerPath := filepath.Join(fixture.artifactDir, reviewDispositionLedgerFileName)
	if _, err := os.Stat(ledgerPath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("refused disposition changed ledger %q: %v", ledgerPath, err)
	}
}

func TestReviewBaseOptionLikeValueIsRefusedAndWritesNothing(t *testing.T) {
	newReviewCommandFixture(t, "codex", &reviewCommandRunner{})
	outputPath := filepath.Join(t.TempDir(), "git-output")
	optionLikeRevision := "--output=" + outputPath
	withReviewSpecGitRunner(t, optionWritingGitRunner{delegate: preflight.ExecGitRunner{}})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"review", "--base", optionLikeRevision}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("review exit=%d stdout=%q stderr=%q, want %d", code, stdout.String(), stderr.String(), exitPreflight)
	}
	if stdout.Len() != 0 || !strings.HasPrefix(stderr.String(), "roundfix: review blocked: resolve review base ") {
		t.Fatalf("review stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
	if _, err := os.Stat(outputPath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("option-like base wrote %q: %v", outputPath, err)
	}
}

func TestDisposeFixedByValidCommitIsStillRecorded(t *testing.T) {
	fixture := newReviewDispositionFixture(t)
	writeDispositionFindingsRecord(t, fixture, "- internal/cli/review.go:10: finding", true)
	mustWrite(t, filepath.Join(fixture.repository, "fix.txt"), "fixed\n")
	gittest.Run(t, fixture.repository, "add", "fix.txt")
	gittest.Run(t, fixture.repository, "commit", "-m", "fix finding")
	fixedBy := strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))

	code, stdout, stderr := runReviewDispose(t, "F1", "--fixed-by", fixedBy)

	if code != exitOK || stderr != "" {
		t.Fatalf("review dispose exit=%d stderr=%q, want exit=0 and empty stderr", code, stderr)
	}
	line := readSingleDispositionLine(t, fixture.artifactDir)
	if stdout != line {
		t.Fatalf("stdout = %q, ledger line = %q", stdout, line)
	}
	entry := decodeDispositionLine(t, line)
	assertDispositionEntry(t, entry, fixture, "F1", "internal/cli/review.go:10: finding", "fixed")
	if entry.FixedBy != fixedBy || entry.Evidence != "" {
		t.Fatalf("fix details = %+v", entry)
	}
}

type optionWritingGitRunner struct {
	delegate preflight.GitRunner
}

func (runner optionWritingGitRunner) RunGit(ctx context.Context, workDir string, args ...string) (string, error) {
	if len(args) > 0 && args[0] == "rev-parse" {
		endOfOptions := false
		for _, arg := range args[1:] {
			if arg == "--end-of-options" {
				endOfOptions = true
				continue
			}
			if !endOfOptions && strings.HasPrefix(arg, "--output=") {
				path := strings.TrimSuffix(strings.TrimPrefix(arg, "--output="), "^{commit}")
				if err := os.WriteFile(path, []byte("unexpected Git output\n"), 0o644); err != nil {
					return "", err
				}
				return "", errors.New("option-like revision was parsed as an option")
			}
		}
	}
	return runner.delegate.RunGit(ctx, workDir, args...)
}
