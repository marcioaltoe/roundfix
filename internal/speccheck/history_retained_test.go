package speccheck

// Suite: gated history retirement and Git readers.
// Invariant: a recorded rule enables diagnostics; removed sources stay readable.
// Boundary IN: planted files and private local Git repositories.
// Boundary OUT: public Check findings, receipt proofs, and skips.

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/spec"
)

// Isolate production Git's process environment without changing parallel tests.
func retainedTestHome(t *testing.T) bool {
	t.Helper()
	const key = "ROUNDFIX_RETAINED_TEST"
	if os.Getenv(key) == t.Name() {
		return true
	}
	cmd := exec.CommandContext(t.Context(), os.Args[0], "-test.run=^"+regexp.QuoteMeta(t.Name())+"$", "-test.count=1")
	cmd.Env = append(gittest.IsolatedEnv(), "HOME="+t.TempDir(), key+"="+t.Name())
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("isolated check: %v\n%s", err, out)
	}
	return false
}
func retainedFixture(t *testing.T, git bool) string {
	t.Helper()
	root := t.TempDir()
	if git {
		gittest.InitRepo(t, root, "-b", "main")
	}
	writeReceiptFixture(t, root, "docs/specs/demo/_prd.md", "# Demo\n")
	retainedGuide(t, root, spec.RetirementRuleSentinel)
	return root
}
func retainedGuide(t *testing.T, root, rule string) {
	t.Helper()
	writeReceiptFixture(t, root, "docs/agents/setup-context.json", `{"managedArtifacts":[{"id":"guide.docs-layout","path":"docs/agents/layout.md"}]}`)
	writeReceiptFixture(t, root, "docs/agents/layout.md", "<!-- setup-context-driven:begin id=guide.docs-layout version=0.0.1 -->\n"+rule+"\n<!-- setup-context-driven:end id=guide.docs-layout -->\n")
}
func retainedCommit(t *testing.T, root string) string {
	t.Helper()
	gittest.Run(t, root, "add", "--all")
	gittest.Run(t, root, "commit", "-m", "docs: fixture")
	return strings.TrimSpace(gittest.Run(t, root, "rev-parse", "HEAD"))
}
func retainedCheck(t *testing.T, root string) Result {
	t.Helper()
	r, err := CheckStage(filepath.Join(root, "docs/specs"), root, "demo", StagePRD)
	if err != nil {
		t.Fatal(err)
	}
	return r
}
func retainedFindings(r Result) []Finding {
	var f []Finding
	for _, v := range r.Findings {
		if v.Code == CodeHistoryRetained {
			f = append(f, v)
		}
	}
	return f
}
func retainedRecord(t *testing.T, root, revision, outcome string) {
	t.Helper()
	content, err := spec.RenderArchiveRecord(spec.ArchiveRecord{Spec: "0001-old", Title: "Old", Disposition: spec.ArchiveNoQA, Source: "docs/specs/0001-old", SourceRevision: revision, PullRequest: "42", Outcome: outcome})
	if err != nil {
		t.Fatal(err)
	}
	writeReceiptFixture(t, root, "docs/history/specs/0001-old.md", string(content))
}

func TestHistoryRetainedIsSilentBeforeTheBaselineCarriesTheRule(t *testing.T) {
	t.Parallel()
	if !retainedTestHome(t) {
		return
	}
	root := retainedFixture(t, false)
	writeReceiptFixture(t, root, "docs/history/adr/old.md", "# Retired\n")
	writeReceiptFixture(t, root, "docs/findings/2026-10-16-done.md", "---\nstatus: done\n---\n")
	retainedRecord(t, root, strings.Repeat("a", 40), "[missing](gone.md) [self](0001-old.md).")
	retainedGuide(t, root, "Old contract.")
	// An unrecorded default path and text outside the region cannot enable it.
	writeReceiptFixture(t, root, "docs/agents/docs-layout.md", spec.RetirementRuleSentinel)
	f := filepath.Join(root, "docs/agents/layout.md")
	content, _ := os.ReadFile(f)
	writeReceiptFixture(t, root, "docs/agents/layout.md", string(content)+spec.RetirementRuleSentinel)
	r := retainedCheck(t, root)
	if len(retainedFindings(r)) != 0 {
		t.Fatalf("findings=%+v", r)
	}
	for _, skip := range r.Skipped {
		if skip.Code == CodeHistoryRetained && strings.Contains(skip.Missing, "guide.docs-layout") {
			return
		}
	}
	t.Fatalf("missing rule skip: %+v", r.Skipped)
}
func TestHistoryRetainedReportsRetainedFamiliesAndTerminalFindings(t *testing.T) {
	t.Parallel()
	if !retainedTestHome(t) {
		return
	}
	root := retainedFixture(t, false)
	for _, family := range []string{"findings", "backlog", "adr", "reviews", "handoffs"} {
		writeReceiptFixture(t, root, "docs/history/"+family+"/retired.md", "# Retired\n")
	}
	writeReceiptFixture(t, root, "docs/history/findings/reduced.md", "---\nstatus: done\n---\n# Reduced\n\nSummary.\n\nFull text in Git at `"+strings.Repeat("a", 40)+"`: `docs/findings/original.md`.\n")
	writeReceiptFixture(t, root, "docs/findings/2026-10-16-done.md", "---\nstatus: done\n---\n")
	writeReceiptFixture(t, root, "docs/findings/2026-10-16-pending.md", "---\nstatus: pending\n---\n")
	r := retainedCheck(t, root)
	f := retainedFindings(r)
	if len(f) != 7 {
		t.Fatalf("findings=%+v", f)
	}
	for _, v := range f {
		if v.Severity != SeverityError || !strings.Contains(v.Fix, "Delete "+v.Where[0].Path) || v.Where[0].Line < 1 {
			t.Fatalf("finding=%+v", v)
		}
	}
	// Whitespace collapse enables the sentence even across a wrap.
	retainedGuide(t, root, "A retirement\n deletes the file")
	if len(retainedFindings(retainedCheck(t, root))) != 7 {
		t.Fatal("wrapped sentinel disabled rule")
	}
}
func TestHistoryRetainedReportsBrokenAndSelfLinksWithTheirGitLocator(t *testing.T) {
	t.Parallel()
	if !retainedTestHome(t) {
		return
	}
	root := retainedFixture(t, true)
	writeReceiptFixture(t, root, "docs/specs/0001-old/file.md", "# Original\n")
	revision := retainedCommit(t, root)
	gittest.Run(t, root, "rm", "-r", "docs/specs/0001-old")
	retainedRecord(t, root, revision, "[missing](file.md) [self](0001-old.md).")
	writeReceiptFixture(t, root, "docs/user-guide/example.md", "[missing](../specs/0001-old/file.md)\n`[code](none.md)`\n")
	writeReceiptFixture(t, root, "docs/references/example.md", "[missing](../specs/0001-old/file.md)\n")
	retainedCommit(t, root)
	f := retainedFindings(retainedCheck(t, root))
	if len(f) != 4 {
		t.Fatalf("findings=%+v", f)
	}
	for _, v := range f {
		if !strings.Contains(v.Fix, "@"+revision+"`") || !strings.Contains(v.Fix, "`docs/specs/0001-old") {
			t.Fatalf("locator=%+v", v)
		}
	}
}
func TestHistoryRetainedReportsAGitLocatorNoReachableCommitHolds(t *testing.T) {
	t.Parallel()
	if !retainedTestHome(t) {
		return
	}
	root := retainedFixture(t, true)
	reachable := retainedCommit(t, root)
	gittest.Run(t, root, "checkout", "-b", "feat/delivery")
	writeReceiptFixture(t, root, "docs/specs/0001-old/file.md", "# Branch only\n")
	unreachable := retainedCommit(t, root)
	gittest.Run(t, root, "checkout", "main")
	gittest.Run(t, root, "branch", "-D", "feat/delivery")
	retainedRecord(t, root, unreachable, "`docs/specs/0001-old/file.md@"+unreachable+"` and `docs/no.md@"+strings.Repeat("a", 40)+"` and `docs/missing.md@"+reachable+"`.")
	f := retainedFindings(retainedCheck(t, root))
	if len(f) != 3 {
		t.Fatalf("findings=%+v", f)
	}
	for _, v := range f {
		if !strings.Contains(v.Fix, "(pull request #42)") || strings.Contains(v.Fix, "@"+unreachable) {
			t.Fatalf("replacement=%+v", v)
		}
	}
	gittest.Run(t, root, "tag", "retained-source", unreachable)
	if got := retainedFindings(retainedCheck(t, root)); len(got) != 2 {
		t.Fatalf("tagged locator not accepted: %+v", got)
	}
}
func TestReceiptsResolveARemovedADRThroughGit(t *testing.T) {
	t.Parallel()
	if !retainedTestHome(t) {
		return
	}
	root := retainedFixture(t, true)
	path := "docs/adr/0001-removed.md"
	writeReceiptFixture(t, root, path, "---\nstatus: superseded\n---\nthree exact words\n")
	revision := retainedCommit(t, root)
	gittest.Run(t, root, "rm", path)
	retainedCommit(t, root)
	source, proven, reason, err := ProveReceipt(root, Receipt{Source: "ADR-0001", Quote: "three exact words"})
	if err != nil || !proven || reason != "" || source != path+"@"+revision {
		t.Fatalf("receipt=%s %v %s %v", source, proven, reason, err)
	}
	_, proven, _, err = ProveReceipt(root, Receipt{Source: "ADR-0001", Quote: "three wrong words"})
	if err != nil || proven {
		t.Fatalf("wrong quote=%v %v", proven, err)
	}
}
func TestRollupMembersResolveThroughGit(t *testing.T) {
	t.Parallel()
	if !retainedTestHome(t) {
		return
	}
	root := retainedFixture(t, true)
	for _, family := range []string{"docs/findings", "docs/history/findings"} {
		path := family + "/2026-10-16-" + filepath.Base(filepath.Dir(family)) + ".md"
		writeReceiptFixture(t, root, path, "---\nstatus: pending\n---\n")
		retainedCommit(t, root)
		gittest.Run(t, root, "rm", path)
	}
	writeReceiptFixture(t, root, "docs/findings/2026-10-16-rollup.md", "---\nstatus: pending\nkind: rollup\nmembers:\n  - 2026-10-16-docs.md\n  - 2026-10-16-history.md\n  - never-existed.md\n---\n")
	retainedCommit(t, root)
	r := retainedCheck(t, root)
	var f []Finding
	for _, v := range r.Findings {
		if v.Code == CodeRollupMember {
			f = append(f, v)
		}
	}
	if len(f) != 1 || !strings.Contains(f[0].Summary, "never-existed.md") {
		t.Fatalf("members=%+v", f)
	}
}
func TestTerminalBacklogEntryFixSaysDelete(t *testing.T) {
	t.Parallel()
	root := retainedFixture(t, false)
	writeReceiptFixture(t, root, "docs/backlog/2026-10-16-closed.md", "---\nstatus: done\nreason: no work remains\n---\n")
	r := retainedCheck(t, root)
	for _, v := range r.Findings {
		if v.Code == CodeBacklogUnmoved {
			if !strings.HasPrefix(v.Fix, "Delete ") || !strings.Contains(v.Fix, "Git keeps its text") {
				t.Fatalf("fix=%s", v.Fix)
			}
			return
		}
	}
	t.Fatal("missing backlog finding")
}

func TestHistoryRetainedAcceptsALocatorTheFetchedDefaultReachesInACIClone(t *testing.T) {
	t.Parallel()
	if !retainedTestHome(t) {
		return
	}
	origin := retainedFixture(t, true)
	older := retainedCommit(t, origin)
	gittest.Run(t, origin, "tag", "older", older)
	path := "docs/specs/0001-old/file.md"
	writeReceiptFixture(t, origin, path, "# Between releases\n")
	source := retainedCommit(t, origin)
	gittest.Run(t, origin, "rm", "--", path)
	removed := retainedCommit(t, origin)
	// A later release tag on a sibling must not prove source reachability.
	gittest.Run(t, origin, "checkout", "--detach", older)
	writeReceiptFixture(t, origin, "release.txt", "Later tagged release\n")
	retainedCommit(t, origin)
	gittest.Run(t, origin, "tag", "later")
	gittest.Run(t, origin, "checkout", "main")
	root := t.TempDir()
	gittest.InitRepo(t, root)
	gittest.Run(t, root, "remote", "add", "origin", origin)
	gittest.Run(t, root, "fetch", "origin", "+refs/heads/*:refs/remotes/origin/*", "+refs/tags/*:refs/tags/*")
	gittest.Run(t, root, "checkout", "--detach", removed)
	retainedRecord(t, root, source, "`"+path+"@"+source+"`")
	if got := retainedFindings(retainedCheck(t, root)); len(got) != 0 {
		t.Fatalf("fetched locator rejected: %+v", got)
	}
	gittest.Run(t, root, "branch", "main", older)
	gittest.Run(t, root, "checkout", "main")
	if got := retainedFindings(retainedCheck(t, root)); len(got) != 0 {
		t.Fatalf("stale local default hid fetched locator: %+v", got)
	}
	gittest.Run(t, root, "checkout", "--detach", removed)
	gittest.Run(t, root, "branch", "-D", "main")
	gittest.Run(t, root, "update-ref", "-d", "refs/remotes/origin/main")
	got := retainedFindings(retainedCheck(t, root))
	if len(got) != 1 || !strings.Contains(got[0].Summary, "invalid Git Locator "+path+"@"+source) {
		t.Fatalf("missing unreachable locator finding: %+v", got)
	}
}
