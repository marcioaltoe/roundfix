package spec

// Suite: removed-file provenance and shared citation rewriting.
// Invariant: locators name reachable objects, and unrelated Markdown bytes survive.
// Boundary IN: temporary Git repositories in child test processes with private HOME.
// Boundary OUT: rendered citations and Archive Record bytes; no remote access.

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"roundfix/internal/gittest"
)

// A subprocess gives each parallel test a private process environment. t.Setenv
// is forbidden after t.Parallel; production Git reads the process environment.
func gitLocatorTestHome(t *testing.T) bool {
	t.Helper()
	const key = "ROUNDFIX_GIT_LOCATOR_TEST"
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
		t.Fatalf("isolated test: %v\n%s", err, out)
	}
	return false
}

func gitLocatorRepo(t *testing.T) (string, string) {
	t.Helper()
	root := t.TempDir()
	gittest.InitRepo(t, root, "-b", "main")
	writeFile(t, filepath.Join(root, "docs/specs/demo/_prd.md"), "# Planning\n")
	return root, legacyCommit(t, root, "docs: plan demo")
}

func requireLocator(t *testing.T, root, path, preferred, pr, want string) GitLocator {
	t.Helper()
	got, err := ResolveGitLocator(t.Context(), root, path, preferred, pr)
	if err != nil || got.Commit != want || got.Path != path || got.PullRequest != pr {
		t.Fatalf("locator=%+v err=%v, want commit %s", got, err, want)
	}
	if want != "" {
		holds, err := GitLocatorHolds(t.Context(), root, path, want)
		if err != nil || !holds {
			t.Fatalf("locator does not hold path: %v %v", holds, err)
		}
	}
	return got
}

func TestResolveGitLocatorPrefersAReachableProvenanceCommit(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	root, planning := gitLocatorRepo(t)
	path := "docs/specs/demo/_prd.md"
	writeFile(t, filepath.Join(root, path), "# Revised planning\n")
	legacyCommit(t, root, "docs: revise plan")
	requireLocator(t, root, path, planning, "", planning)
	requireLocator(t, root, "docs/specs/demo", planning, "", planning)
	// An annotated tag can retain provenance outside the default branch.
	gittest.Run(t, root, "checkout", "-b", "feat/tagged")
	writeFile(t, filepath.Join(root, path), "# Tagged planning\n")
	tagged := legacyCommit(t, root, "docs: tagged plan")
	gittest.Run(t, root, "tag", "-a", "retained-plan", "-m", "retained")
	gittest.Run(t, root, "checkout", "main")
	gittest.Run(t, root, "branch", "-D", "feat/tagged")
	requireLocator(t, root, path, tagged, "", tagged)
	gittest.Run(t, root, "checkout", "--detach")
	// With no detectable default branch, tags alone can retain the preference.
	requireLocator(t, root, path, tagged, "", tagged)
}

func TestResolveGitLocatorFallsBackToTheNewestDefaultBranchCommit(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	root, _ := gitLocatorRepo(t)
	path := "docs/specs/demo/_prd.md"
	writeFile(t, filepath.Join(root, path), "# Newest\n")
	latest := legacyCommit(t, root, "docs: newest plan")
	requireLocator(t, root, path, strings.Repeat("f", 40), "", latest)
	gittest.Run(t, root, "rm", "--", path)
	legacyCommit(t, root, "docs: remove plan")
	requireLocator(t, root, path, "", "", latest)
}

func TestResolveGitLocatorNamesNoUnreachableCommit(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	root, planning := gitLocatorRepo(t)
	path := "docs/specs/demo/qa/evidence/output.txt"
	gittest.Run(t, root, "checkout", "-b", "feat/delivery")
	writeFile(t, filepath.Join(root, path), "Evidence\n")
	unretained := legacyCommit(t, root, "test: delivery evidence")
	gittest.Run(t, root, "checkout", "main")
	// A squash delivery retains the result, not the branch's original commit.
	writeFile(t, filepath.Join(root, "delivery.txt"), "Delivered\n")
	legacyCommit(t, root, "feat: squash delivery")
	gittest.Run(t, root, "branch", "-D", "feat/delivery")
	got := requireLocator(t, root, path, unretained, "42", "")
	if got.Text(false) != "`"+path+"` (pull request #42)" {
		t.Fatal(got.Text(false))
	}
	requireLocator(t, root, "docs/specs/demo/_prd.md", unretained, "", planning)
	// Other live work branches are not accepted as provenance reachability.
	gittest.Run(t, root, "branch", "feat/other", unretained)
	requireLocator(t, root, path, unretained, "", "")
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if _, err := ResolveGitLocator(ctx, root, path, planning, ""); err == nil {
		t.Fatal("cancelled lookup succeeded")
	}
	if _, err := ResolveGitLocator(t.Context(), t.TempDir(), path, planning, ""); err == nil {
		t.Fatal("non-repository lookup succeeded")
	} else {
		var exit *exec.ExitError
		if !errors.As(err, &exit) {
			t.Fatalf("Git error is not wrapped: %v", err)
		}
	}
}

func TestResolveGitLocatorReadsTheRemoteDefaultBranch(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	root, planning := gitLocatorRepo(t)
	gittest.Run(t, root, "update-ref", "refs/remotes/origin/trunk", planning)
	gittest.Run(t, root, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/trunk")
	gittest.Run(t, root, "checkout", "-b", "feat/work")
	gittest.Run(t, root, "branch", "-D", "main")
	requireLocator(t, root, "docs/specs/demo/_prd.md", "", "", planning)
	// An unpushed local commit cannot override the fetched default head.
	gittest.Run(t, root, "checkout", "-b", "trunk")
	writeFile(t, filepath.Join(root, "docs/specs/demo/_prd.md"), "# Local trunk\n")
	latest := legacyCommit(t, root, "docs: local trunk plan")
	requireLocator(t, root, "docs/specs/demo/_prd.md", "", "", planning)
	requireLocator(t, root, "docs/specs/demo/_prd.md", latest, "", planning)
}

func TestGitLocatorTextForms(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	for _, tc := range []struct {
		name    string
		locator GitLocator
		want    string
	}{
		{"commit", GitLocator{Path: "docs/plan.md", Commit: strings.Repeat("a", 40)}, "docs/plan.md@" + strings.Repeat("a", 40)},
		{"pull request", GitLocator{Path: "docs/plan.md", PullRequest: "42"}, "docs/plan.md (pull request #42)"},
		{"unreachable", GitLocator{Path: "docs/plan.md"}, "docs/plan.md (not reachable from any branch or tag)"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.locator.Text(true); got != tc.want {
				t.Fatalf("in code=%q want=%q", got, tc.want)
			}
			path, suffix, _ := strings.Cut(tc.want, " ")
			want := "`" + path + "`"
			if suffix != "" {
				want += " " + suffix
			}
			if got := tc.locator.Text(false); got != want {
				t.Fatalf("text=%q want=%q", got, want)
			}
		})
	}
}

func TestCitationRewriteReplacesLinksAndPathsWithGitLocators(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	root, commit := gitLocatorRepo(t)
	folder := "docs/specs/demo"
	rw := CitationRewrite{Moves: map[string]string{folder + "/references/promoted.md": "docs/references/promoted.md"}, Locate: func(path string) (GitLocator, bool, error) {
		if path != folder && !strings.HasPrefix(path, folder+"/") {
			return GitLocator{}, false, nil
		}
		locator, err := ResolveGitLocator(t.Context(), root, path, commit, "")
		return locator, true, err
	}}
	locator := GitLocator{Path: folder + "/_prd.md", Commit: commit}
	for _, tc := range []struct{ name, input, want string }{
		{"link title", `[**plan**](_prd.md?q=1#goals "Title")`, "**plan** (" + locator.Text(false) + ")"},
		{"image", `![plan](<_prd.md> 'Title')`, "plan (" + locator.Text(false) + ")"},
		{"reference", "  [plan]: <_prd.md#goals> \"Title\"\r\n", "  \\[plan]: " + locator.Text(false) + "\r\n"},
		{"path", folder + "/_prd.md.", locator.Text(false) + "."},
		{"code path", "`" + folder + "/_prd.md`", "`" + locator.Text(true) + "`"},
		{"fence path", "~~~md\n" + folder + "/_prd.md\n~~~\n", "~~~md\n" + locator.Text(true) + "\n~~~\n"},
		{"promoted", `[kept](references/promoted.md?q=1#goals "Title")`, `[kept](../../references/promoted.md?q=1#goals "Title")`},
		{"outside", `[ADR](../../adr/rule.md#rule)`, `[ADR](../../adr/rule.md#rule)`},
		{"unreachable", "[evidence](qa/evidence/out.txt)", "evidence (`" + folder + "/qa/evidence/out.txt` (not reachable from any branch or tag))"},
		{"preserve", "\x00[web](https://example.com)\r\n`[plan](_prd.md)`", "\x00[web](https://example.com)\r\n`[plan](_prd.md)`"},
		{"nested label", `[![icon](https://example.com/icon)](_prd.md "Title")`, "![icon](https://example.com/icon) (" + locator.Text(false) + ")"},
		{"existing moved locator", (GitLocator{Path: folder + "/references/promoted.md", PullRequest: "42"}).Text(false), (GitLocator{Path: folder + "/references/promoted.md", PullRequest: "42"}).Text(false)},
		{"existing locator", locator.Text(false), locator.Text(false)},
		{"existing unreachable locator", (GitLocator{Path: folder + "/missing.md"}).Text(false), (GitLocator{Path: folder + "/missing.md"}).Text(false)},
		{"existing PR locator", (GitLocator{Path: folder + "/missing.md", PullRequest: "42"}).Text(false), (GitLocator{Path: folder + "/missing.md", PullRequest: "42"}).Text(false)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, _, err := rw.Apply([]byte(tc.input), folder+"/_prd.md", "docs/history/specs/demo.md")
			if err != nil || string(got) != tc.want {
				t.Fatalf("rewrite=%q err=%v want=%q", got, err, tc.want)
			}
		})
	}
	sentinel := errors.New("locator unavailable")
	rw.Locate = func(string) (GitLocator, bool, error) { return GitLocator{}, false, sentinel }
	if _, _, err := rw.Apply([]byte("[plan](_prd.md)"), folder+"/_prd.md", "record.md"); !errors.Is(err, sentinel) {
		t.Fatalf("lookup failure=%v", err)
	}
}

func TestCitationRewriteKeepsRelinkContentBytes(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	folder := "docs/specs/demo"
	moves := map[string]string{folder: "docs/history/specs/demo.md", folder + "/_prd.md": "docs/history/specs/demo.md", folder + "/probe data.md": "docs/references/probe data.md"}
	input := []byte("\x00[plan](../specs/demo/_prd.md?q=1#goals \"title\")\r\n![probe](<../specs/demo/probe%20data.md>)\n[label]: ../specs/demo/_prd.md 'Title'\n`" + folder + "/_prd.md`\n" + folder + "/\n" + folder + "/not-mapped.md\n")
	want := []byte("\x00[plan](../history/specs/demo.md?q=1#goals \"title\")\r\n![probe](<probe%20data.md>)\n[label]: ../history/specs/demo.md 'Title'\n`docs/history/specs/demo.md`\ndocs/history/specs/demo.md\n" + folder + "/not-mapped.md\n")
	got, count, err := (CitationRewrite{Moves: moves}).Apply(input, "docs/references/citing.md", "docs/references/citing.md")
	old, oldCount := RelinkContent(input, "docs/references/citing.md", moves)
	if err != nil || !bytes.Equal(got, want) || !bytes.Equal(old, got) || count != 5 || oldCount != count {
		t.Fatalf("rewrite=%q count=%d oldCount=%d err=%v want=%q", got, count, oldCount, err, want)
	}
	for _, tc := range []struct {
		name, path, input, want string
		moves                   map[string]string
	}{
		{"outside repository", "docs/citing.md", "[text](../../removed.md)", "[text](kept.md)", map[string]string{"../removed.md": "docs/kept.md"}},
		{"incompatible path roots", "/absolute/citing.md", "[text](removed.md)", "[text](removed.md)", map[string]string{"/absolute/removed.md": "relative/kept.md"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, n, err := (CitationRewrite{Moves: tc.moves}).Apply([]byte(tc.input), tc.path, tc.path)
			old, oldN := RelinkContent([]byte(tc.input), tc.path, tc.moves)
			if err != nil || string(got) != tc.want || !bytes.Equal(got, old) || n != oldN {
				t.Fatalf("rewrite=%q old=%q err=%v want=%q", got, old, err, tc.want)
			}
		})
	}

}

func TestLegacyLinkRebaseNeverLinksTheRecordToItself(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	root, commit := gitLocatorRepo(t)
	folder := "docs/specs/demo"
	rebase := LegacyLinkRebase{Folder: folder, RecordPath: "docs/history/specs/demo.md", Locate: func(path string) (GitLocator, bool, error) {
		l, err := ResolveGitLocator(t.Context(), root, path, commit, "")
		return l, true, err
	}}
	got, err := rebase.Rebase([]byte("[folder](./) [plan](_prd.md#goals) [self](../../history/specs/demo.md)"), folder+"/_prd.md", rebase.RecordPath)
	want := "folder (`" + folder + "@" + commit + "`) plan (`" + folder + "/_prd.md@" + commit + "`) self (`" + folder + "@" + commit + "`)"
	if err != nil || string(got) != want || len(scanArchiveMarkdown(got)) != 0 {
		t.Fatalf("outcome=%q err=%v want=%q", got, err, want)
	}
}

func TestArchiveRecordShorteningKeepsGitLocatorsWhole(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	for _, kind := range []string{"cut inside locator", "sentence inside locator", "multiple backticks", "whole locator before cut"} {
		t.Run(kind, func(t *testing.T) {
			r := ArchiveRecord{Spec: "demo", Title: "Demo", Source: "docs/specs/demo", SourceRevision: strings.Repeat("a", 40), Disposition: ArchivePass}
			base, err := RenderArchiveRecord(r)
			if err != nil {
				t.Fatal(err)
			}
			budget := ArchiveRecordTargetBytes - len(base) - len("…")
			locator := GitLocator{Path: "docs/specs/demo/_techspec.md", Commit: strings.Repeat("b", 40)}.Text(false)
			prefix := "Before " + strings.Repeat("x", budget-len("Before ")-12)
			if kind == "sentence inside locator" {
				prefix = "Sentence. " + strings.Repeat("x", budget-len("Sentence. ")-40)
			}
			if kind == "multiple backticks" {
				locator = "`" + locator + "`"
			}
			if kind == "whole locator before cut" {
				prefix = strings.Repeat("x", budget-len(locator)-20) + " "
			}
			r.Outcome = prefix + locator + strings.Repeat(" after", 400)
			got, err := RenderArchiveRecord(r)
			if err != nil {
				t.Fatal(err)
			}
			body := string(got[len(base)-1:])
			if len(got) > ArchiveRecordTargetBytes || strings.Count(body, "`")%2 != 0 || (strings.Contains(body, "docs/specs/demo") && !strings.Contains(body, locator)) {
				t.Fatalf("cut code span: %q (%d bytes)", body, len(got))
			}
			if kind == "sentence inside locator" {
				if body != "Sentence.\n" {
					t.Fatalf("sentence cut=%q", body)
				}
			} else if kind != "whole locator before cut" {
				if body != prefix+"…\n" {
					t.Fatalf("fallback cut=%q want=%q", body, prefix+"…\n")
				}
			} else if !strings.Contains(body, locator) || !strings.HasSuffix(body, "…\n") {
				t.Fatalf("whole locator lost: %q", body)
			}
		})
	}
}

// Build the same refs as an Actions checkout, without clone's origin/HEAD.
func gitLocatorCIClone(t *testing.T) (string, string, string) {
	t.Helper()
	origin, older := gitLocatorRepo(t)
	gittest.Run(t, origin, "tag", "older", older)
	path := "docs/specs/demo/_prd.md"
	writeFile(t, filepath.Join(origin, path), "# Fetched planning\n")
	newer := legacyCommit(t, origin, "docs: fetched plan")
	gittest.Run(t, origin, "rm", "--", path)
	head := legacyCommit(t, origin, "docs: remove fetched plan")
	root := t.TempDir()
	gittest.InitRepo(t, root)
	gittest.Run(t, root, "remote", "add", "origin", origin)
	gittest.Run(t, root, "fetch", "origin", "+refs/heads/*:refs/remotes/origin/*", "+refs/tags/*:refs/tags/*")
	gittest.Run(t, root, "checkout", "--detach", head)
	refs := gittest.Run(t, root, "for-each-ref", "--format=%(refname)", "refs/heads/main", "refs/remotes/origin/HEAD")
	if strings.TrimSpace(refs) != "" {
		t.Fatalf("unexpected CI refs: %s", refs)
	}
	return root, older, newer
}

func TestResolveGitLocatorReachesTheFetchedDefaultInACIClone(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	root, _, newer := gitLocatorCIClone(t)
	requireLocator(t, root, "docs/specs/demo/_prd.md", newer, "", newer)
}

func TestResolveGitLocatorIgnoresAStaleLocalDefaultBranch(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	root, older, newer := gitLocatorCIClone(t)
	gittest.Run(t, root, "branch", "main", older)
	gittest.Run(t, root, "checkout", "main")
	requireLocator(t, root, "docs/specs/demo/_prd.md", newer, "", newer)
	requireLocator(t, root, "docs/specs/demo/_prd.md", "", "", newer)
}
