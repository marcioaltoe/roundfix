// Suite: contract horizon
// Invariant: only provably older Specs escape the receipt gap after guide adoption.
// Boundary IN: disposable Git histories and CheckStage reports.
// Boundary OUT: the repository's own history and Daemon settlement.
package speccheck

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
)

func writeReceiptFixture(t *testing.T, root, path, text string) {
	t.Helper()
	full := filepath.Join(root, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(text), 0644); err != nil {
		t.Fatal(err)
	}
}

const receiptFixturePRD = "docs/specs/0200-example/_prd.md"
const receiptFixtureADR = "docs/adr/0116-a-citation-is-checked-against-what-it-cites.md"
const receiptFixtureClaim = "ADR-0116 requires the check to read the cited record."

func receiptHorizonRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	gittest.InitRepo(t, root, "-b", "main")
	writeReceiptFixture(t, root, receiptFixtureADR, "# A citation is checked against what it cites\n\nRoundfix reads the cited record when an artifact makes a claim.\n")
	receiptCommit(t, root)
	return root
}

func receiptCommit(t *testing.T, root string) {
	t.Helper()
	gittest.Run(t, root, "add", "--all")
	gittest.Run(t, root, "commit", "-m", "fixture")
}

func receiptStage(t *testing.T, root string, stage Stage) Result {
	t.Helper()
	result, err := CheckStage(filepath.Join(root, "docs/specs"), root, "0200-example", stage)
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func receiptCodeFindings(result Result, code string) []Finding {
	var found []Finding
	for _, finding := range result.Findings {
		if finding.Code == code {
			found = append(found, finding)
		}
	}
	return found
}

func requireReceiptHorizon(t *testing.T, root string, held bool, missing string) {
	t.Helper()
	horizon := newContractHorizon(root, filepath.Join(root, receiptFixturePRD))
	if horizon.held != held || horizon.missing != missing {
		t.Fatalf("horizon = %#v, want held %v missing %q", horizon, held, missing)
	}
	result := receiptStage(t, root, StagePRD)
	want := 0
	if held {
		want = 1
	}
	if got := receiptCodeFindings(result, CodeReceiptMissing); len(got) != want {
		t.Fatalf("receipt gaps = %#v, want %d", got, want)
	}
	skips := 0
	for _, skip := range result.Skipped {
		if skip.Code == CodeReceiptMissing {
			skips++
			if !strings.Contains(RenderText(result, VerificationCoverage{}), CodeReceiptMissing+": missing "+missing) {
				t.Fatalf("missing skip reason: %#v", skip)
			}
		}
	}
	if held && skips != 0 || !held && skips != 1 {
		t.Fatalf("skips = %d, held %v", skips, held)
	}
}

func TestASpecCommittedBeforeTheGuideIsNotHeld(t *testing.T) {
	t.Parallel()
	root := receiptHorizonRepo(t)
	writeReceiptFixture(t, root, receiptFixturePRD, receiptFixtureClaim)
	receiptCommit(t, root)
	writeReceiptFixture(t, root, ConcreteContractGuidePath, "guide")
	receiptCommit(t, root)
	requireReceiptHorizon(t, root, false, "a PRD committed at or after "+ConcreteContractGuidePath)
}

func TestASpecCommittedWithOrAfterTheGuideIsHeld(t *testing.T) {
	t.Parallel()
	for _, when := range []string{"with", "after"} {
		t.Run(when, func(t *testing.T) {
			root := receiptHorizonRepo(t)
			writeReceiptFixture(t, root, ConcreteContractGuidePath, "guide")
			if when == "after" {
				receiptCommit(t, root)
			}
			writeReceiptFixture(t, root, receiptFixturePRD, receiptFixtureClaim)
			receiptCommit(t, root)
			requireReceiptHorizon(t, root, true, "")
		})
	}
}

func TestAnUncommittedSpecIsHeldOnceTheGuideIsCommitted(t *testing.T) {
	t.Parallel()
	root := receiptHorizonRepo(t)
	writeReceiptFixture(t, root, ConcreteContractGuidePath, "guide")
	receiptCommit(t, root)
	writeReceiptFixture(t, root, receiptFixturePRD, receiptFixtureClaim)
	requireReceiptHorizon(t, root, true, "")
}

func TestARepositoryWithoutTheGuideHoldsNoSpec(t *testing.T) {
	t.Parallel()
	root := receiptHorizonRepo(t)
	writeReceiptFixture(t, root, receiptFixturePRD, receiptFixtureClaim)
	requireReceiptHorizon(t, root, false, ConcreteContractGuidePath)
}

func TestAnUncommittedGuideHoldsNoSpec(t *testing.T) {
	t.Parallel()
	root := receiptHorizonRepo(t)
	writeReceiptFixture(t, root, receiptFixturePRD, receiptFixtureClaim)
	receiptCommit(t, root)
	writeReceiptFixture(t, root, ConcreteContractGuidePath, "guide")
	requireReceiptHorizon(t, root, false, "a committed "+ConcreteContractGuidePath)
}

func TestAGuideWithUnreadableHistoryHoldsEverySpec(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	writeReceiptFixture(t, root, ConcreteContractGuidePath, "guide")
	writeReceiptFixture(t, root, receiptFixtureADR, "# A citation\n\nreads the cited record\n")
	writeReceiptFixture(t, root, receiptFixturePRD, receiptFixtureClaim)
	requireReceiptHorizon(t, root, true, "")
}

func TestARevisedOldPRDKeepsItsContractHorizon(t *testing.T) {
	t.Parallel()
	root := receiptHorizonRepo(t)
	writeReceiptFixture(t, root, receiptFixturePRD, receiptFixtureClaim)
	receiptCommit(t, root)
	writeReceiptFixture(t, root, ConcreteContractGuidePath, "guide")
	receiptCommit(t, root)
	writeReceiptFixture(t, root, receiptFixturePRD, receiptFixtureClaim+"\nRevised.")
	receiptCommit(t, root)
	requireReceiptHorizon(t, root, false, "a PRD committed at or after "+ConcreteContractGuidePath)
}

func TestAReadoptedGuideKeepsItsOldestAddingCommit(t *testing.T) {
	t.Parallel()
	root := receiptHorizonRepo(t)
	writeReceiptFixture(t, root, ConcreteContractGuidePath, "guide")
	receiptCommit(t, root)
	writeReceiptFixture(t, root, receiptFixturePRD, receiptFixtureClaim)
	receiptCommit(t, root)
	if err := os.Remove(filepath.Join(root, ConcreteContractGuidePath)); err != nil {
		t.Fatal(err)
	}
	receiptCommit(t, root)
	writeReceiptFixture(t, root, ConcreteContractGuidePath, "guide again")
	receiptCommit(t, root)
	requireReceiptHorizon(t, root, true, "")
}

func TestContractHorizonFailsClosedForShallowOrMismatchedRoots(t *testing.T) {
	t.Parallel()
	source := receiptHorizonRepo(t)
	writeReceiptFixture(t, source, ConcreteContractGuidePath, "guide")
	writeReceiptFixture(t, source, receiptFixturePRD, receiptFixtureClaim)
	receiptCommit(t, source)
	t.Run("shallow", func(t *testing.T) {
		clone := filepath.Join(t.TempDir(), "clone")
		gittest.Run(t, "", "clone", "--depth=1", "file://"+filepath.ToSlash(source), clone)
		gittest.Harden(t, clone)
		requireReceiptHorizon(t, clone, true, "")
	})
	t.Run("nested root", func(t *testing.T) {
		root := filepath.Join(source, "nested")
		writeReceiptFixture(t, root, ConcreteContractGuidePath, "guide")
		writeReceiptFixture(t, root, receiptFixturePRD, receiptFixtureClaim)
		if got := newContractHorizon(root, filepath.Join(root, receiptFixturePRD)); !got.held {
			t.Fatalf("horizon = %#v", got)
		}
	})
	t.Run("external PRD", func(t *testing.T) {
		external := receiptHorizonRepo(t)
		writeReceiptFixture(t, external, receiptFixturePRD, receiptFixtureClaim)
		if got := newContractHorizon(source, filepath.Join(external, receiptFixturePRD)); !got.held {
			t.Fatalf("horizon = %#v", got)
		}
	})
}

func TestANonFileGuideHoldsNoSpec(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"directory", "symlink"} {
		t.Run(kind, func(t *testing.T) {
			root := receiptHorizonRepo(t)
			writeReceiptFixture(t, root, receiptFixturePRD, receiptFixtureClaim)
			path := filepath.Join(root, ConcreteContractGuidePath)
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				t.Fatal(err)
			}
			if kind == "directory" {
				if err := os.Mkdir(path, 0755); err != nil {
					t.Fatal(err)
				}
			} else {
				writeReceiptFixture(t, root, "guide.md", "guide")
				if err := os.Symlink(filepath.Join(root, "guide.md"), path); err != nil {
					t.Fatal(err)
				}
			}
			requireReceiptHorizon(t, root, false, ConcreteContractGuidePath)
		})
	}
}

func TestAGuideOnADivergentBranchDoesNotHoldTheSpec(t *testing.T) {
	t.Parallel()
	root := receiptHorizonRepo(t)
	gittest.Run(t, root, "branch", "spec")
	writeReceiptFixture(t, root, ConcreteContractGuidePath, "guide")
	receiptCommit(t, root)
	gittest.Run(t, root, "checkout", "spec")
	writeReceiptFixture(t, root, receiptFixturePRD, receiptFixtureClaim)
	receiptCommit(t, root)
	gittest.Run(t, root, "merge", "--no-edit", "main")
	requireReceiptHorizon(t, root, false, "a PRD committed at or after "+ConcreteContractGuidePath)
}

func TestAnUnreadableGuidePathDoesNotProveTheContractAbsent(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	writeReceiptFixture(t, root, ConcreteContractGuidePath, "guide")
	writeReceiptFixture(t, root, receiptFixturePRD, receiptFixtureClaim)
	directory := filepath.Dir(filepath.Join(root, ConcreteContractGuidePath))
	if err := os.Chmod(directory, 0); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chmod(directory, 0755); err != nil {
			t.Error(err)
		}
	})
	if got := newContractHorizon(root, filepath.Join(root, receiptFixturePRD)); !got.held {
		t.Fatalf("horizon = %#v, want held", got)
	}
}
