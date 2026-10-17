package cli

// Suite: Archive Command locators and Delivery Queue recomputation.
// Invariant: removed-folder citations survive without dangling or self links.
// Boundary IN: public archive dispatch and disposable Git histories with private HOME.
// Boundary OUT: record/promotion bytes and exact committed and working-tree proofs.

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/spec"
)

func archiveLocatorTestHome(t *testing.T) bool {
	t.Helper()
	const key = "ROUNDFIX_ARCHIVE_LOCATOR_TEST"
	if os.Getenv(key) == t.Name() {
		return true
	}
	cmd := exec.CommandContext(t.Context(), os.Args[0], "-test.run=^"+regexp.QuoteMeta(t.Name())+"$", "-test.count=1")
	for _, entry := range gittest.IsolatedEnv() {
		if !strings.HasPrefix(entry, "HOME=") && !strings.HasPrefix(entry, key+"=") {
			cmd.Env = append(cmd.Env, entry)
		}
	}
	cmd.Env = append(cmd.Env, "HOME="+t.TempDir(), key+"="+t.Name())
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("isolated archive test: %v\n%s", err, out)
	}
	return false
}

func TestArchiveOutcomeLinksBecomeGitLocators(t *testing.T) {
	t.Parallel()
	if !archiveLocatorTestHome(t) {
		return
	}
	for _, pr := range []string{"", "42"} {
		t.Run("pull-request-"+pr, func(t *testing.T) {
			repo, env, _ := archiveRelinkFixture(t, false)
			folder := "docs/specs/" + archiveRelinkOwner
			adopted := folder + "/references/adopted.md"
			evidence := folder + "/qa/evidence/output.txt"
			mustWrite(t, filepath.Join(repo, adopted), "---\ntype: fix\nstatus: promoted\ncreated: 2026-01-01\n---\n\n# Adopted intent\n")
			// Select the actual newest report written by the shared command fixture.
			report, err := spec.NewestQAReport(filepath.Join(repo, folder))
			if err != nil {
				t.Fatal(err)
			}
			qa := report
			mustWrite(t, qa, mustRead(t, qa)+"\n## Outcome\n\nSee [adopted](../references/adopted.md) and [evidence](evidence/output.txt).\n")
			commitArchiveFixture(t)
			planning := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
			gittest.Run(t, repo, "update-ref", "refs/remotes/origin/main", planning)
			gittest.Run(t, repo, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/main")
			gittest.Run(t, repo, "checkout", "-b", "feat/delivery")
			mustMkdir(t, filepath.Dir(filepath.Join(repo, evidence)))
			mustWrite(t, filepath.Join(repo, evidence), "Delivery-only QA evidence\n")
			gittest.Run(t, repo, "add", "-A")
			gittest.Run(t, repo, "commit", "-m", "test: delivery-only QA evidence")
			delivery := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
			// Keep the original delivery commit locally but remove its branch, as after
			// a squash merge. Only the planning commit is on the default branch.
			branch := strings.TrimSpace(gittest.Run(t, repo, "branch", "--show-current"))
			gittest.Run(t, repo, "checkout", "--detach", delivery)
			gittest.Run(t, repo, "branch", "-D", branch)
			var out, stderr bytes.Buffer
			args := []string{"archive", archiveRelinkOwner}
			if pr != "" {
				args = append(args, "--pull-request", pr)
			}
			code := runWithContext(t.Context(), args, &out, &stderr, env)
			if code != exitOK {
				t.Fatalf("exit=%d stdout=%s stderr=%s", code, &out, &stderr)
			}
			record := archiveRecordAt(t, filepath.Join(repo, archiveRelinkRecord))
			want := "See adopted (`" + adopted + "@" + planning + "`) and evidence (" + (spec.GitLocator{Path: evidence, PullRequest: pr}).Text(false) + ")."
			if record.Outcome != want || record.PullRequest != pr {
				t.Fatalf("outcome=%q pull request=%s want=%q", record.Outcome, record.PullRequest, want)
			}
			if strings.Contains(record.Outcome, delivery) || strings.Contains(record.Outcome, "](") {
				t.Fatalf("unsafe locator: %s", record.Outcome)
			}
		})
	}
}

func TestArchivePromotedCopyKeepsItsTargetsThroughGitLocators(t *testing.T) {
	t.Parallel()
	if !archiveLocatorTestHome(t) {
		return
	}
	repo, env, _ := archiveRelinkFixture(t, false)
	folder := "docs/specs/" + archiveRelinkOwner
	plan := folder + "/_techspec.md"
	mustWrite(t, filepath.Join(repo, plan), "# Design\n")
	mustWrite(t, filepath.Join(repo, archiveRelinkSource), "# Source\r\n\r\n[design](../_techspec.md#rule), [peer](peer.md?q=1#data), [ADR](../../../adr/rule.md).\nRead `"+plan+"`.\n")
	peer := folder + "/references/peer.md"
	mustWrite(t, filepath.Join(repo, peer), "# Peer\n")
	mustMkdir(t, filepath.Join(repo, "docs/adr"))
	mustWrite(t, filepath.Join(repo, "docs/adr/rule.md"), "# Rule\n")
	commitArchiveFixture(t)
	planning := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
	gittest.Run(t, repo, "update-ref", "refs/heads/main", planning)
	gittest.Run(t, repo, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/main")
	var out, stderr bytes.Buffer
	code := runWithContext(t.Context(), []string{"archive", archiveRelinkOwner, "--promote", "references/2026-01-01-source.md", "--promote", "references/peer.md"}, &out, &stderr, env)
	want := "# Source\r\n\r\ndesign (`" + plan + "@" + planning + "`), [peer](peer.md?q=1#data), [ADR](../adr/rule.md).\nRead `" + plan + "@" + planning + "`.\n"
	got := mustRead(t, filepath.Join(repo, archiveRelinkPromotion))
	if code != exitOK || got != want {
		t.Fatalf("exit=%d copy=%q stderr=%s want=%q", code, got, &stderr, want)
	}
	if mustRead(t, filepath.Join(repo, "docs/references/peer.md")) != "# Peer\n" {
		t.Fatal("link-free copy changed")
	}
	assertPathMissing(t, filepath.Join(repo, folder))
}

func TestDeliveryAcceptsAnArchiveWithRewrittenPromotedCopies(t *testing.T) {
	t.Parallel()
	if !archiveLocatorTestHome(t) {
		return
	}
	for _, tamper := range []string{"", "copy byte", "locator commit", "relative link", "ambiguous basename"} {
		t.Run("tamper-"+tamper, func(t *testing.T) {
			w, repo, _, source, destination := deliveryRelinkFixture(t, false)
			original := "Read [plan](../_prd.md) and `" + source + "/_prd.md`.\r\n"
			mustWrite(t, filepath.Join(repo, source, "references/knowledge.md"), original)
			if tamper == "ambiguous basename" {
				mustMkdir(t, filepath.Join(repo, source, "other"))
				mustWrite(t, filepath.Join(repo, source, "other/knowledge.md"), "different bytes\n")
			}
			commitArchiveFixture(t)
			parent := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
			gittest.Run(t, repo, "update-ref", "refs/heads/main", parent)
			gittest.Run(t, repo, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/main")
			archiveDeliveryRelinkFixture(t, repo, parent)
			copyPath := filepath.Join(repo, deliveryRelinkPromotion)
			want := "Read plan (`" + source + "/_prd.md@" + parent + "`) and `" + source + "/_prd.md@" + parent + "`.\r\n"
			if got := mustRead(t, copyPath); got != want {
				t.Fatalf("copy=%q want=%q", got, want)
			}
			switch tamper {
			case "copy byte":
				mustWrite(t, copyPath, want+"x")
			case "locator commit":
				mustWrite(t, copyPath, strings.ReplaceAll(want, parent, strings.Repeat("f", 40)))
			case "relative link":
				mustWrite(t, filepath.Join(repo, deliveryRelinkPath), mustRead(t, filepath.Join(repo, deliveryRelinkPath))+"x")
			}
			assertDeliveryRelinkProof(t, w, repo, parent, source, destination, tamper == "")
		})
	}
}
