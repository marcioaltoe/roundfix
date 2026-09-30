// Suite: pre-PR review finding identities and dispositions.
// Invariant: each finding is addressable and accepts at most one head-bound recorded disposition.
// Boundary IN: review result classification, public CLI parsing, local Git checks, and Artifact Directory files.
// Boundary OUT: reviewer execution, delivery reuse of dispositions, and pull-request publication.
package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/gittest"
)

func TestReviewRecordListsEachFindingWithAnIdentity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		answer string
		want   []reviewFinding
	}{
		{
			name:   "marked answer",
			answer: "Findings:\n- internal/cli/review.go:10: first\n* internal/cli/review.go:20: second\n3) internal/cli/review.go:30: third",
			want: []reviewFinding{
				{ID: "F1", Text: "internal/cli/review.go:10: first"},
				{ID: "F2", Text: "internal/cli/review.go:20: second"},
				{ID: "F3", Text: "internal/cli/review.go:30: third"},
			},
		},
		{
			name:   "unmarked answer",
			answer: "Findings:\ninternal/cli/review.go:10: one unmarked finding",
			want: []reviewFinding{
				{ID: "F1", Text: "internal/cli/review.go:10: one unmarked finding"},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			record := dispositionReviewRecord("/tmp/repository", "head", test.answer)
			got, code := classifyReviewCommandResult(record, agent.ExecuteResult{Message: test.answer}, nil)
			if code != exitRunFailed {
				t.Fatalf("classify findings exit = %d, want %d", code, exitRunFailed)
			}
			if !reflect.DeepEqual(got.FindingItems, test.want) {
				t.Fatalf("finding items = %+v, want %+v", got.FindingItems, test.want)
			}
		})
	}

	t.Run("legacy disk record derives items", func(t *testing.T) {
		record := dispositionReviewRecord("/tmp/repository", "head", "internal/cli/review.go:10: legacy finding")
		record.Findings = "internal/cli/review.go:10: legacy finding"
		record.FindingItems = nil
		directory := reviewCheckoutDir(t.TempDir(), record.Repository)
		if err := os.MkdirAll(directory, 0o755); err != nil {
			t.Fatalf("create review checkout directory: %v", err)
		}
		path := filepath.Join(directory, reviewRecordFileName)
		writeDispositionRecordFile(t, path, record)

		got, err := readReviewRecord(path)
		if err != nil {
			t.Fatalf("read legacy review record: %v", err)
		}
		want := []reviewFinding{{ID: "F1", Text: record.Findings}}
		if !reflect.DeepEqual(got.FindingItems, want) {
			t.Fatalf("legacy finding items = %+v, want %+v", got.FindingItems, want)
		}
	})
}

func TestReviewPromptAsksForOneListItemPerFinding(t *testing.T) {
	t.Parallel()

	prompt := buildReviewPrompt("base", "head", "diff")
	if !strings.Contains(prompt, "each finding as a list item that starts with `- ` and names its file and line") {
		t.Fatalf("review prompt does not require one marked item per finding:\n%s", prompt)
	}
}

func TestSplitReviewFindingsKeepsPreambleAndIndentedLines(t *testing.T) {
	t.Parallel()

	findings := "preamble\ncontinues\n- first\n  - indented continuation\n4. second\n5) third\n*   "
	want := []reviewFinding{
		{ID: "F1", Text: "preamble\ncontinues"},
		{ID: "F2", Text: "first\n  - indented continuation"},
		{ID: "F3", Text: "second"},
		{ID: "F4", Text: "third"},
	}
	if got := splitReviewFindings(findings); !reflect.DeepEqual(got, want) {
		t.Fatalf("finding items = %+v, want %+v", got, want)
	}
}

func TestReviewDisposeHelpPrintsUsage(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"review", "dispose", "--help"}, &stdout, &stderr)

	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("review dispose help exit=%d stderr=%q", code, stderr.String())
	}
	for _, usageLine := range []string{
		"roundfix review dispose <finding-id> --dismiss --evidence <text>",
		"roundfix review dispose <finding-id> --fixed-by <commit>",
	} {
		if !strings.Contains(stdout.String(), usageLine) {
			t.Fatalf("review dispose help does not contain %q:\n%s", usageLine, stdout.String())
		}
	}
}

func TestReviewDisposeDismissesAFindingWithEvidence(t *testing.T) {
	fixture := newReviewDispositionFixture(t)
	const findings = "- internal/cli/review.go:10: first\n- internal/cli/review.go:20: second"
	writeDispositionFindingsRecord(t, fixture, findings, false)

	code, stdout, stderr := runReviewDispose(t, "F2", "--dismiss", "--evidence", "not reachable on the supported path")

	if code != exitOK || stderr != "" {
		t.Fatalf("review dispose exit=%d stderr=%q, want exit=0 and empty stderr", code, stderr)
	}
	line := readSingleDispositionLine(t, fixture.artifactDir)
	if stdout != line {
		t.Fatalf("stdout = %q, ledger line = %q", stdout, line)
	}
	entry := decodeDispositionLine(t, line)
	assertDispositionEntry(t, entry, fixture, "F2", "internal/cli/review.go:20: second", "dismissed")
	if entry.Evidence != "not reachable on the supported path" || entry.FixedBy != "" {
		t.Fatalf("dismissal details = %+v", entry)
	}
}

func TestReviewDisposeRecordsAFixingCommit(t *testing.T) {
	fixture := newReviewDispositionFixture(t)
	const findings = "- internal/cli/review.go:10: first"
	writeDispositionFindingsRecord(t, fixture, findings, true)
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
	assertDispositionEntry(t, entry, fixture, "F1", "internal/cli/review.go:10: first", "fixed")
	if entry.FixedBy != fixedBy || entry.Evidence != "" {
		t.Fatalf("fix details = %+v", entry)
	}
}

func TestReviewDisposeRefusesWithoutAFindingsRecord(t *testing.T) {
	tests := []struct {
		name  string
		write func(*testing.T, reviewCommandFixture)
	}{
		{name: "record is absent"},
		{
			name: "record belongs to another repository",
			write: func(t *testing.T, fixture reviewCommandFixture) {
				record := dispositionReviewRecord(filepath.Join(fixture.repository, "other"), fixture.headCommit, "Findings:\n- internal/cli/review.go:10: finding")
				writeDispositionRecord(t, fixture, record)
			},
		},
		{
			name: "record has another outcome",
			write: func(t *testing.T, fixture reviewCommandFixture) {
				record := newReviewRecord(fixture.repository, fixture.baseCommit, fixture.headCommit, roundconfig.PrePRReview{Provider: "codex", Source: "project"}, reviewOutcomeReviewed)
				writeDispositionRecord(t, fixture, record)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fixture := newReviewDispositionFixture(t)
			if test.write != nil {
				test.write(t, fixture)
			}
			assertDisposeRefusedWithoutLedgerChange(t, fixture.artifactDir, "F1", "--dismiss", "--evidence", "evidence")
		})
	}
}

func TestReviewDisposeRefusesAnUnknownFinding(t *testing.T) {
	fixture := newReviewDispositionFixture(t)
	writeDispositionFindingsRecord(t, fixture, "- internal/cli/review.go:10: finding", true)
	assertDisposeRefusedWithoutLedgerChange(t, fixture.artifactDir, "F9", "--dismiss", "--evidence", "evidence")
}

func TestReviewDisposeRefusesBothForms(t *testing.T) {
	fixture := newReviewDispositionFixture(t)
	writeDispositionFindingsRecord(t, fixture, "- internal/cli/review.go:10: finding", true)
	assertDisposeRefusedWithoutLedgerChange(t, fixture.artifactDir, "F1", "--dismiss", "--evidence", "evidence", "--fixed-by", fixture.baseCommit)
}

func TestReviewDisposeRefusesNeitherForm(t *testing.T) {
	fixture := newReviewDispositionFixture(t)
	writeDispositionFindingsRecord(t, fixture, "- internal/cli/review.go:10: finding", true)
	assertDisposeRefusedWithoutLedgerChange(t, fixture.artifactDir, "F1")
}

func TestReviewDisposeRefusesBlankEvidence(t *testing.T) {
	fixture := newReviewDispositionFixture(t)
	writeDispositionFindingsRecord(t, fixture, "- internal/cli/review.go:10: finding", true)
	assertDisposeRefusedWithoutLedgerChange(t, fixture.artifactDir, "F1", "--dismiss", "--evidence", "  \t")
}

func TestReviewDisposeRefusesADismissalAtAMovedHead(t *testing.T) {
	fixture := newReviewDispositionFixture(t)
	writeDispositionFindingsRecord(t, fixture, "- internal/cli/review.go:10: finding", true)
	mustWrite(t, filepath.Join(fixture.repository, "moved.txt"), "moved\n")
	gittest.Run(t, fixture.repository, "add", "moved.txt")
	gittest.Run(t, fixture.repository, "commit", "-m", "move head")
	assertDisposeRefusedWithoutLedgerChange(t, fixture.artifactDir, "F1", "--dismiss", "--evidence", "evidence")
}

func TestReviewDisposeRefusesAFixThatDoesNotDescendFromTheReviewedHead(t *testing.T) {
	tests := []struct {
		name    string
		fixedBy func(*testing.T, reviewCommandFixture) string
	}{
		{name: "commit does not resolve", fixedBy: func(*testing.T, reviewCommandFixture) string { return "missing-commit" }},
		{name: "commit is the reviewed head", fixedBy: func(_ *testing.T, f reviewCommandFixture) string { return f.headCommit }},
		{name: "commit predates the reviewed head", fixedBy: func(_ *testing.T, f reviewCommandFixture) string { return f.baseCommit }},
		{
			name: "commit is not reachable from current head",
			fixedBy: func(t *testing.T, f reviewCommandFixture) string {
				gittest.Run(t, f.repository, "checkout", "-b", "unreachable-fix")
				mustWrite(t, filepath.Join(f.repository, "unreachable.txt"), "fix\n")
				gittest.Run(t, f.repository, "add", "unreachable.txt")
				gittest.Run(t, f.repository, "commit", "-m", "unreachable fix")
				fixedBy := strings.TrimSpace(gittest.Run(t, f.repository, "rev-parse", "HEAD"))
				gittest.Run(t, f.repository, "checkout", "feature/review")
				return fixedBy
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fixture := newReviewDispositionFixture(t)
			writeDispositionFindingsRecord(t, fixture, "- internal/cli/review.go:10: finding", true)
			assertDisposeRefusedWithoutLedgerChange(t, fixture.artifactDir, "F1", "--fixed-by", test.fixedBy(t, fixture))
		})
	}
}

func TestReviewDisposeRefusesASecondDisposition(t *testing.T) {
	fixture := newReviewDispositionFixture(t)
	writeDispositionFindingsRecord(t, fixture, "- internal/cli/review.go:10: finding", true)
	code, _, stderr := runReviewDispose(t, "F1", "--dismiss", "--evidence", "evidence")
	if code != exitOK || stderr != "" {
		t.Fatalf("first disposition exit=%d stderr=%q", code, stderr)
	}
	mustWrite(t, filepath.Join(fixture.repository, "fix.txt"), "fixed\n")
	gittest.Run(t, fixture.repository, "add", "fix.txt")
	gittest.Run(t, fixture.repository, "commit", "-m", "fix finding")
	fixedBy := strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))

	assertDisposeRefusedWithoutLedgerChange(t, fixture.artifactDir, "F1", "--fixed-by", fixedBy)
}

func newReviewDispositionFixture(t *testing.T) reviewCommandFixture {
	t.Helper()
	return newReviewCommandFixture(t, "codex", &reviewCommandRunner{})
}

func dispositionReviewRecord(repository string, head string, answer string) reviewRecord {
	record := newReviewRecord(repository, "base", head, roundconfig.PrePRReview{Provider: "codex", Source: "project"}, reviewOutcomeFindings)
	findings := strings.TrimSpace(strings.TrimPrefix(answer, "Findings:"))
	record.Findings = findings
	record.FindingItems = splitReviewFindings(findings)
	return record
}

func writeDispositionFindingsRecord(t *testing.T, fixture reviewCommandFixture, findings string, includeItems bool) {
	t.Helper()
	record := newReviewRecord(fixture.repository, fixture.baseCommit, fixture.headCommit, roundconfig.PrePRReview{Provider: "codex", Source: "project"}, reviewOutcomeFindings)
	record.Findings = findings
	if includeItems {
		record.FindingItems = splitReviewFindings(findings)
	}
	writeDispositionRecord(t, fixture, record)
}

func writeDispositionRecord(t *testing.T, fixture reviewCommandFixture, record reviewRecord) {
	t.Helper()
	directory := reviewCheckoutDir(fixture.artifactDir, fixture.repository)
	if err := os.MkdirAll(directory, 0o755); err != nil {
		t.Fatalf("create review checkout directory: %v", err)
	}
	writeDispositionRecordFile(t, filepath.Join(directory, reviewRecordFileName), record)
}

func writeDispositionRecordFile(t *testing.T, path string, record reviewRecord) {
	t.Helper()
	var encoded bytes.Buffer
	if err := writeReviewRecord(&encoded, record); err != nil {
		t.Fatalf("encode review record: %v", err)
	}
	if err := os.WriteFile(path, encoded.Bytes(), 0o644); err != nil {
		t.Fatalf("write review record: %v", err)
	}
}

func runReviewDispose(t *testing.T, args ...string) (int, string, string) {
	t.Helper()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	commandArgs := append([]string{"review", "dispose"}, args...)
	code := runCLI(t, commandArgs, &stdout, &stderr)
	return code, stdout.String(), stderr.String()
}

func assertDisposeRefusedWithoutLedgerChange(t *testing.T, artifactDir string, args ...string) {
	t.Helper()
	path := filepath.Join(artifactDir, reviewDispositionLedgerFileName)
	before, beforeErr := os.ReadFile(path)
	if beforeErr != nil && !errors.Is(beforeErr, os.ErrNotExist) {
		t.Fatalf("read ledger before refusal: %v", beforeErr)
	}
	code, stdout, stderr := runReviewDispose(t, args...)
	if code != exitPreflight {
		t.Fatalf("review dispose exit=%d stdout=%q stderr=%q, want %d", code, stdout, stderr, exitPreflight)
	}
	if stdout != "" || !strings.HasPrefix(stderr, "roundfix: review dispose refused: ") {
		t.Fatalf("review dispose stdout=%q stderr=%q", stdout, stderr)
	}
	after, afterErr := os.ReadFile(path)
	if !sameDispositionLedgerState(before, beforeErr, after, afterErr) {
		t.Fatalf("refusal changed ledger: before=%q err=%v after=%q err=%v", before, beforeErr, after, afterErr)
	}
}

func sameDispositionLedgerState(before []byte, beforeErr error, after []byte, afterErr error) bool {
	beforeAbsent := errors.Is(beforeErr, os.ErrNotExist)
	afterAbsent := errors.Is(afterErr, os.ErrNotExist)
	if beforeAbsent || afterAbsent {
		return beforeAbsent && afterAbsent
	}
	return beforeErr == nil && afterErr == nil && bytes.Equal(before, after)
}

func readSingleDispositionLine(t *testing.T, artifactDir string) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(artifactDir, reviewDispositionLedgerFileName))
	if err != nil {
		t.Fatalf("read disposition ledger: %v", err)
	}
	if bytes.Count(content, []byte("\n")) != 1 || !bytes.HasSuffix(content, []byte("\n")) {
		t.Fatalf("disposition ledger = %q, want one JSON line", content)
	}
	return string(content)
}

func decodeDispositionLine(t *testing.T, line string) reviewFindingDisposition {
	t.Helper()
	var entry reviewFindingDisposition
	if err := json.Unmarshal([]byte(line), &entry); err != nil {
		t.Fatalf("decode disposition line %q: %v", line, err)
	}
	return entry
}

func assertDispositionEntry(t *testing.T, entry reviewFindingDisposition, fixture reviewCommandFixture, finding string, text string, disposition string) {
	t.Helper()
	if entry.Repository != fixture.repository || entry.HeadCommit != fixture.headCommit || entry.Finding != finding || entry.Text != text || entry.Disposition != disposition {
		t.Fatalf("disposition entry = %+v", entry)
	}
	recordedAt, err := time.Parse(time.RFC3339, entry.RecordedAt)
	if err != nil || recordedAt.Location() != time.UTC {
		t.Fatalf("recordedAt = %q, want RFC 3339 UTC: %v", entry.RecordedAt, err)
	}
}
