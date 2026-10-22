package cli

// Suite: bounded Delivery Queue archive subjects and pull request titles.
// Invariant: fitting slugs keep their bytes; squash readers and archive adoption survive.
// Boundary IN: repository-authored titles, disposable repositories, homes and hooks.
// Boundary OUT: no forge, network or host Git configuration.

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf16"
	"unicode/utf8"

	"roundfix/internal/gittest"
	"roundfix/internal/releaseplan"
	"roundfix/internal/spec"
)

func queueArchiveSubjectCases() []struct{ name, slug, want string } {
	return []struct{ name, slug, want string }{
		{"86 units", "0296-" + strings.Repeat("x", 81), "docs: archive 0296-" + strings.Repeat("x", 81)},
		{"87 units", "0296-" + strings.Repeat("x", 82), "docs: archive the 0296 Spec"},
		{"230 units", "0296-" + strings.Repeat("x", 225), "docs: archive the 0296 Spec"},
		{"ten digits", "1234567890-" + strings.Repeat("x", 90), "docs: archive the 1234567890 Spec"},
		{"eleven digits", "12345678901-" + strings.Repeat("x", 90), "docs: archive the Spec"},
		{"no digits", strings.Repeat("x", 100), "docs: archive the Spec"},
		{"non ASCII digits", "０２９６-" + strings.Repeat("x", 100), "docs: archive the Spec"},
		{"astral boundary", "0296-" + strings.Repeat("x", 79) + "😀", "docs: archive 0296-" + strings.Repeat("x", 79) + "😀"},
		{"astral overflow", "0296-" + strings.Repeat("x", 80) + "😀", "docs: archive the 0296 Spec"},
	}
}

func TestQueueArchiveSubjectKeepsTheSlugWhileItFits(t *testing.T) {
	t.Parallel()
	if got := queueHeaderLength("a😀é"); got != 4 {
		t.Fatalf("UTF-16 length=%d, want 4", got)
	}
	for _, tc := range queueArchiveSubjectCases() {
		t.Run(tc.name, func(t *testing.T) {
			got := archiveCommitSubject(tc.slug)
			if got != tc.want || len(utf16.Encode([]rune(got))) > 100 {
				t.Fatalf("subject=%q, want %q within 100 units", got, tc.want)
			}
		})
	}
}

func queuePullRequestTitleCases() []struct{ name, slug, title, want string } {
	long := "0296-" + strings.Repeat("x", 95)
	return []struct{ name, slug, title, want string }{
		{"77 units", "0296-" + strings.Repeat("x", 72), "", "feat: deliver 0296-" + strings.Repeat("x", 72)},
		{"78 units", "0296-" + strings.Repeat("x", 73), "A useful title", "feat: deliver 0296: A useful title"},
		{"0278", "0278-a-fixed-round-2-finding-gets-one-correction-review-and-the-plan-names-only-what-it-delivers", "A fixed round-2 finding gets one Correction Review, and the plan names only what it delivers", "feat: deliver 0278: A fixed round-2 finding gets one Correction Review, and the plan names"},
		{"full stop", long, "A useful title.", "feat: deliver 0296: A useful title"},
		{"whitespace", long, "  A\t useful\n title\r\nwith\u2003spaces  ", "feat: deliver 0296: A useful title with spaces"},
		{"long first word", long, strings.Repeat("x", 80), "feat: deliver 0296: " + strings.Repeat("x", 71)},
		{"astral rune", long, strings.Repeat("x", 69) + "😀yz", "feat: deliver 0296: " + strings.Repeat("x", 69) + "😀"},
		{"astral rune at cut", long, strings.Repeat("x", 71) + "😀", "feat: deliver 0296: " + strings.Repeat("x", 71)},
		{"whole word at cut", long, strings.Repeat("x", 71) + " more", "feat: deliver 0296: " + strings.Repeat("x", 71)},
		{"partial word at cut", long, "First " + strings.Repeat("x", 80), "feat: deliver 0296: First"},
		{"punctuation", long, "A title.,;: ", "feat: deliver 0296: A title"},
		{"empty numbered", long, " \n .,:; ", "feat: deliver 0296"},
		{"empty unnumbered", strings.Repeat("x", 100), "", "feat: deliver the Spec"},
		{"unnumbered", strings.Repeat("x", 100), "A useful title", "feat: deliver A useful title"},
		{"ten digits", "1234567890-" + strings.Repeat("x", 90), strings.Repeat("x", 80), "feat: deliver 1234567890: " + strings.Repeat("x", 65)},
		{"eleven digits", "12345678901-" + strings.Repeat("x", 90), "A useful title", "feat: deliver A useful title"},
	}
}

func TestQueuePullRequestTitleKeepsTheSlugWhileItFits(t *testing.T) {
	t.Parallel()
	for _, tc := range queuePullRequestTitleCases() {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			got, err := pullRequestTitle(tc.slug, func() (string, error) {
				calls++
				if tc.name == "77 units" {
					t.Fatal("fitting slug called title reader")
				}
				return tc.title, nil
			})
			wantCalls := 1
			if tc.name == "77 units" {
				wantCalls = 0
			}
			if err != nil || got != tc.want || calls != wantCalls || !utf8.ValidString(got) || len(utf16.Encode([]rune(got+" (#99999)"))) > 100 {
				t.Fatalf("title=%q err=%v reads=%d, want %q reads=%d within squash budget", got, err, calls, tc.want, wantCalls)
			}
		})
	}
	t.Run("reader error", func(t *testing.T) {
		failure := errors.New("missing archive")
		got, err := pullRequestTitle(strings.Repeat("x", 100), func() (string, error) { return "", failure })
		if got != "" || !errors.Is(err, failure) || err.Error() != "read Spec title: missing archive" {
			t.Fatalf("title=%q err=%v", got, err)
		}
	})
}

func TestQueueMessagesStayReadableBySquashReaders(t *testing.T) {
	t.Parallel()
	for _, tc := range queuePullRequestTitleCases() {
		t.Run("title/"+tc.name, func(t *testing.T) {
			title, err := pullRequestTitle(tc.slug, func() (string, error) { return tc.title, nil })
			if err != nil {
				t.Fatal(err)
			}
			subject := title + " (#575)"
			if number, ok := spec.PullRequestFromSubject(subject); !ok || number != "575" {
				t.Fatalf("number=%q ok=%v subject=%q", number, ok, subject)
			}
			evidence := releaseplan.ClassifyCommit(releaseplan.Commit{Subject: subject})
			if evidence.ConventionalType != "feat" || evidence.AutomaticImpact != releaseplan.ImpactMinor || evidence.Breaking {
				t.Fatalf("classification=%+v", evidence)
			}
		})
	}
	for _, tc := range queueArchiveSubjectCases() {
		t.Run("archive/"+tc.name, func(t *testing.T) {
			evidence := releaseplan.ClassifyCommit(releaseplan.Commit{Subject: archiveCommitSubject(tc.slug)})
			if evidence.ConventionalType != "docs" || evidence.Breaking {
				t.Fatalf("classification=%+v", evidence)
			}
		})
	}
}

func queueArchiveSubjectFixture(t *testing.T, slug string) (*commandDeliveryWorkflow, string, string) {
	t.Helper()
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	writeArchiveQAReport(t, repo, spec.VerdictPass)
	old := filepath.Join(repo, "docs/specs", implementTestSlug)
	next := filepath.Join(repo, "docs/specs", slug)
	if err := os.Rename(old, next); err != nil {
		t.Fatal(err)
	}
	if err := filepath.WalkDir(next, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() {
			mustWrite(t, path, strings.ReplaceAll(mustRead(t, path), implementTestSlug, slug))
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	commitArchiveFixture(t)
	gittest.PersistIdentity(t, repo)
	setCommandEnvForTest(t, cliTestHelperEnv, "1")
	return retirementWorkflow(repo, home), repo, strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
}

func TestDeliveryWorkflowArchiveCommitPassesAHeaderLimitHook(t *testing.T) {
	t.Parallel()
	if !queueCommitTestHome(t) {
		return
	}
	slug := "0296-" + strings.Repeat("x", 95)
	w, repo, candidate := queueArchiveSubjectFixture(t, slug)
	writeScriptFixture(t, filepath.Join(repo, ".git", "hooks", "commit-msg"), `IFS= read -r header < "$1"
if [ "${#header}" -gt 100 ]; then
 echo 'header must not be longer than 100 characters' >&2
 exit 1
fi
`)
	result, err := w.Archive(commandContextForTest(t, t.Context()), repo, slug, candidate)
	if err != nil || !result.ExactSpecMove || result.Parent != candidate || result.Head == "" || result.Head == candidate {
		t.Fatalf("archive=%+v err=%v", result, err)
	}
	if got := strings.TrimSpace(gittest.Run(t, repo, "show", "-s", "--format=%s", result.Head)); got != "docs: archive the 0296 Spec" {
		t.Fatalf("subject=%q", got)
	}
	if dirty := gittest.Run(t, repo, "status", "--porcelain"); dirty != "" {
		t.Fatalf("dirty=%s", dirty)
	}
}

func TestDeliveryWorkflowAdoptsAnArchiveCommitWhateverItsSubject(t *testing.T) {
	t.Parallel()
	if !queueCommitTestHome(t) {
		return
	}
	slug := "0296-" + strings.Repeat("x", 95)
	w, repo, candidate := queueArchiveSubjectFixture(t, slug)
	hook := filepath.Join(repo, ".git", "hooks", "pre-commit")
	writeScriptFixture(t, hook, "exit 1\n")
	ctx := commandContextForTest(t, t.Context())
	_, err := w.Archive(ctx, repo, slug, candidate)
	if err == nil || !strings.Contains(err.Error(), `refused the commit "docs: archive the 0296 Spec"`) {
		t.Fatalf("refusal=%v", err)
	}
	if staged := strings.TrimSpace(gittest.Run(t, repo, "diff", "--cached", "--name-only")); staged == "" {
		t.Fatal("refusal lost staged archive")
	}
	mustWrite(t, hook+scriptFixtureSuffix, "exit 0\n")
	gittest.Run(t, repo, "commit", "-m", "docs: archive "+slug)
	hand := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
	result, err := w.Archive(ctx, repo, slug, candidate)
	if err != nil || !result.ExactSpecMove || result.Head != hand || result.Parent != candidate {
		t.Fatalf("adopt=%+v err=%v", result, err)
	}
	if head := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD")); head != hand {
		t.Fatalf("head=%s want=%s", head, hand)
	}
	if count := strings.TrimSpace(gittest.Run(t, repo, "rev-list", "--count", candidate+"..HEAD")); count != "1" {
		t.Fatalf("commits=%s", count)
	}
	if dirty := gittest.Run(t, repo, "status", "--porcelain"); dirty != "" {
		t.Fatalf("dirty=%s", dirty)
	}
}

func TestDeliveryPublicationTitleNamesTheSpecNumberForALongSlug(t *testing.T) {
	t.Parallel()
	if !queueCommitTestHome(t) {
		return
	}
	slug := "0296-" + strings.Repeat("x", 95)
	w, repo, _, path := queueRecordHookFixture(t, slug)
	record := archiveRecordAt(t, filepath.Join(repo, path))
	record.Title = "A known Archive Record title."
	data, err := spec.RenderArchiveRecord(record)
	if err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(repo, path), string(data))
	branch := strings.TrimSpace(gittest.Run(t, repo, "branch", "--show-current"))
	publication, err := w.Publication(t.Context(), repo, slug, branch)
	if err != nil || publication.Title != "feat: deliver 0296: A known Archive Record title" || publication.Body != "Delivers Spec "+slug+"." {
		t.Fatalf("publication=%+v err=%v", publication, err)
	}
	// A configured Spec root stores its archive under _archived, not docs/history.
	customRoot := filepath.Join(repo, "planning")
	customArchive := filepath.Join(customRoot, "_archived")
	mustMkdir(t, customArchive)
	customPath := filepath.Join(customArchive, slug+".md")
	if err := os.Rename(filepath.Join(repo, path), customPath); err != nil {
		t.Fatal(err)
	}
	w.loaded.Config.Specs.Root = "planning"
	configured, err := w.Publication(t.Context(), repo, slug, branch)
	if err != nil || configured != publication {
		t.Fatalf("configured archive publication=%+v err=%v, want %+v", configured, err, publication)
	}
	if err := os.Remove(customPath); err != nil {
		t.Fatal(err)
	}
	if _, err := w.Publication(t.Context(), repo, slug, branch); err == nil || !strings.HasPrefix(err.Error(), "plan publication: read Spec title:") {
		t.Fatalf("missing record error=%v", err)
	}
}
