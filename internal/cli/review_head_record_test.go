// Boundary: public review/dispose commands, symlinked Git worktrees, and
// persisted head records. Reviewer execution uses a local fake.
package cli

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestASecondReviewRoundKeepsTheFirstRoundsRecord(t *testing.T) {
	t.Parallel()
	fixture, first, second := reviewSecondRoundUnderSymlink(t)
	dir := reviewCheckoutDir(fixture.artifactDir, fixture.repository)
	firstPath := filepath.Join(dir, "heads", first.HeadCommit+".json")
	kept, err := readReviewRecord(firstPath)
	if err != nil {
		t.Fatal(err)
	}
	if kept.Outcome != reviewOutcomeFindings || !reflect.DeepEqual(kept, first) {
		t.Fatalf("candidate record=%+v, want original round-1 findings %+v", kept, first)
	}
	secondPath := filepath.Join(dir, "heads", second.HeadCommit+".json")
	fixed, err := readReviewRecord(secondPath)
	if err != nil {
		t.Fatal(err)
	}
	current, err := readReviewRecord(filepath.Join(dir, reviewRecordFileName))
	if err != nil {
		t.Fatal(err)
	}
	if fixed.Outcome != reviewOutcomeReviewed || !reflect.DeepEqual(fixed, second) || !reflect.DeepEqual(current, second) {
		t.Fatalf("fix record=%+v current=%+v, want round 2 %+v", fixed, current, second)
	}
	if !bytes.Equal(mustReadReviewCheckoutFile(t, dir, reviewRecordFileName), mustReadReviewCheckoutFile(t, filepath.Dir(secondPath), filepath.Base(secondPath))) {
		t.Fatal("current and head record bytes differ")
	}
	for _, path := range []string{firstPath, secondPath} {
		info, err := os.Stat(path)
		if err != nil || info.Mode().Perm() != 0o644 {
			t.Fatalf("head record permissions: info=%v err=%v, want 0644", info, err)
		}
	}
	info, err := os.Stat(filepath.Join(dir, "heads"))
	if err != nil || info.Mode().Perm() != 0o755 {
		t.Fatalf("head directory permissions: info=%v err=%v, want 0755", info, err)
	}
}

func TestTheCorrectionProofAcceptsAfterASecondRoundUnderASymlinkedWorktreeRoot(t *testing.T) {
	t.Parallel()
	fixture, first, second := reviewSecondRoundUnderSymlink(t)
	result := proveSymlinkedReviewCorrection(t, fixture, first.HeadCommit, second.HeadCommit, true)
	if !result.Accepted || result.Reason != "" {
		t.Fatalf("proof=%+v, want accepted correction after round 2", result)
	}
}

func reviewSecondRoundUnderSymlink(t *testing.T) (reviewCommandFixture, reviewRecord, reviewRecord) {
	t.Helper()
	fixture, _, _ := newSymlinkedReviewCheckoutFixture(t, true)
	runner := newLineageRunner("Findings:\n- review.txt:2: failure: the changed line loses the review evidence", "No findings.")
	withAgentRunner(t, runner)
	first := reviewSymlinkedCheckout(t, fixture)
	head := commitReviewCorrectionFile(t, fixture, correctionArchivePath, "# Corrected archived Spec\n")
	disposeSymlinkedReviewFinding(t, fixture, head)
	code, second, stderr := runReviewAtCheckout(t, reviewCheckoutFixture{baseCommit: fixture.baseCommit}, fixture.repository)
	if code != exitOK || stderr != "" || second.Outcome != reviewOutcomeReviewed || second.HeadCommit != head || second.Reused || second.Lineage == nil || second.Lineage.Round != 2 {
		t.Fatalf("round 2 exit=%d record=%+v stderr=%q", code, second, stderr)
	}
	if runner.preparedCalls != 2 {
		t.Fatalf("reviewer calls=%d, want two rounds", runner.preparedCalls)
	}
	return fixture, first, second
}

func TestAReusedBlockedOrOmittedReviewWritesNoReviewedHeadRecord(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"reused", "blocked", "omitted"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			artifactDir, checkout := t.TempDir(), t.TempDir()
			existingHead, newHead := strings.Repeat("a", 40), strings.Repeat("b", 64)
			record := newReviewRecord(checkout, "base", existingHead, reviewPolicyForLineage(), reviewOutcomeReviewed)
			if err := persistReviewRecord(artifactDir, record); err != nil {
				t.Fatal(err)
			}
			dir := reviewCheckoutDir(artifactDir, checkout)
			before := mustReadReviewCheckoutFile(t, filepath.Join(dir, "heads"), existingHead+".json")
			switch name {
			case "reused":
				record.Reused = true
				record.Outcome = reviewOutcomeFindings
				record.Findings = "- review.txt:2: standing finding"
				record.FindingItems = splitReviewFindings(record.Findings)
			case "blocked":
				record.Outcome = reviewOutcomeBlocked
				record.Reason = "reviewer unavailable"
			case "omitted":
				record.Outcome = reviewOutcomeOmitted
			}
			for _, head := range []string{existingHead, newHead} {
				record.HeadCommit = head
				if err := persistReviewRecord(artifactDir, record); err != nil {
					t.Fatal(err)
				}
				after := mustReadReviewCheckoutFile(t, filepath.Join(dir, "heads"), existingHead+".json")
				if !bytes.Equal(before, after) {
					t.Fatal("excluded record replaced an existing head record")
				}
				if _, err := os.Stat(filepath.Join(dir, "heads", newHead+".json")); !errors.Is(err, os.ErrNotExist) {
					t.Fatalf("new head record stat=%v, want absent", err)
				}
				current, err := readReviewRecord(filepath.Join(dir, reviewRecordFileName))
				if err != nil || current.HeadCommit != head || current.Outcome != record.Outcome || current.Reused != record.Reused {
					t.Fatalf("current record=%+v err=%v, want excluded record still persisted", current, err)
				}
			}
		})
	}
}

func TestTheCorrectionProofReadsTheCurrentRecordWithoutAHeadRecord(t *testing.T) {
	t.Parallel()
	fixture, _, _ := newSymlinkedReviewCheckoutFixture(t, true)
	record := reviewSymlinkedCheckout(t, fixture)
	dir := reviewCheckoutDir(fixture.artifactDir, fixture.repository)
	if err := os.RemoveAll(filepath.Join(dir, "heads")); err != nil {
		t.Fatal(err)
	}
	head := commitReviewCorrectionFile(t, fixture, correctionArchivePath, "# Corrected archived Spec\n")
	disposeSymlinkedReviewFinding(t, fixture, head)
	result := proveSymlinkedReviewCorrection(t, fixture, record.HeadCommit, head, true)
	if !result.Accepted || result.Reason != "" {
		t.Fatalf("legacy proof=%+v, want accepted", result)
	}
	if _, err := os.Stat(filepath.Join(dir, "heads")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("heads directory stat=%v, want absent", err)
	}
}

func TestTheCorrectionProofRefusesWithoutARecordOfTheCandidate(t *testing.T) {
	t.Parallel()
	fixture, first, second := reviewSecondRoundUnderSymlink(t)
	dir := reviewCheckoutDir(fixture.artifactDir, fixture.repository)
	if err := os.RemoveAll(filepath.Join(dir, "heads")); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"another-current-head", "no-current-record"} {
		t.Run(name, func(t *testing.T) {
			if name == "no-current-record" {
				if err := os.Remove(filepath.Join(dir, reviewRecordFileName)); err != nil {
					t.Fatal(err)
				}
			}
			result := proveSymlinkedReviewCorrection(t, fixture, first.HeadCommit, second.HeadCommit, true)
			if result.Accepted || result.Reason != "review record must name the item worktree and parked candidate with outcome findings" {
				t.Fatalf("proof=%+v, want existing record refusal", result)
			}
		})
	}
}

func TestReviewHeadRecordPathRefusesANonCommitName(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, head string }{
		{"empty", ""},
		{"traversal", "../" + strings.Repeat("a", 37)},
		{"uppercase", strings.Repeat("A", 40)},
		{"short", strings.Repeat("a", 39)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			artifactDir, checkout := t.TempDir(), t.TempDir()
			if path, ok := reviewHeadRecordPath(artifactDir, checkout, tc.head); ok || path != "" {
				t.Fatalf("path=%q ok=%v, want rejection of %q", path, ok, tc.head)
			}
			// Empty heads are already rejected by record validation; other invalid
			// names remain valid current records but must never name a head file.
			if tc.head == "" {
				return
			}
			validHead := strings.Repeat("a", 40)
			record := newReviewRecord(checkout, "base", validHead, reviewPolicyForLineage(), reviewOutcomeReviewed)
			if err := persistReviewRecord(artifactDir, record); err != nil {
				t.Fatal(err)
			}
			dir := reviewCheckoutDir(artifactDir, checkout)
			before := mustReadReviewCheckoutFile(t, filepath.Join(dir, "heads"), validHead+".json")
			record.HeadCommit = tc.head
			if err := persistReviewRecord(artifactDir, record); err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(before, mustReadReviewCheckoutFile(t, filepath.Join(dir, "heads"), validHead+".json")) {
				t.Fatal("invalid head changed existing head record")
			}
			entries, err := os.ReadDir(filepath.Join(dir, "heads"))
			if err != nil || len(entries) != 1 || entries[0].Name() != validHead+".json" {
				t.Fatalf("head files=%v err=%v, want only original head", entries, err)
			}
		})
	}
	for _, length := range []int{40, 64} {
		head := strings.Repeat("0123456789abcdef", 4)[:length]
		artifactDir, checkout := t.TempDir(), t.TempDir()
		want := filepath.Join(reviewCheckoutDir(artifactDir, checkout), "heads", head+".json")
		if path, ok := reviewHeadRecordPath(artifactDir, checkout, head); !ok || path != want {
			t.Fatalf("valid %d-character commit: path=%q ok=%v, want %q", length, path, ok, want)
		}
	}
}
