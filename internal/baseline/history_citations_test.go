package baseline

// Suite: History Relocation citation discovery
// Invariant: only tracked citations broken by applied History Relocations become deterministic findings.
// Boundary IN: a real Git index and the working-tree paths selected from it
// Boundary OUT: plan wiring, rendering, digest binding, and apply behavior

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

const citationSource = "docs/adr/old.md"

func TestRelocationCitationsReportEachCitationForm(t *testing.T) {
	repo := newCitationTestRepository(t)
	trackCitationTestFile(t, repo, citationSource, "retired\n")
	files := map[string]string{
		"a-path.txt":     "docs/adr/old.md\n",
		"b-inline.md":    "[old](docs/adr/old.md)\n",
		"c-image.md":     "![old](docs/adr/old.md)\n",
		"d-reference.md": "[old]: docs/adr/old.md\n",
		"e-root-link.md": "[old](/docs/adr/old.md)\n",
	}
	for path, content := range files {
		trackCitationTestFile(t, repo, path, content)
	}

	findings := citationTestFindings(t, repo, citationTestMove(citationSource), nil)
	if len(findings) != len(files) {
		t.Fatalf("relocationCitationFindings() returned %d findings, want %d: %#v", len(findings), len(files), findings)
	}
	for _, finding := range findings {
		if finding.Code != "baseline.history.citation" || !strings.Contains(finding.Message, "line 1 cites ") || !strings.Contains(finding.Message, "which resolves to docs/adr/old.md") {
			t.Errorf("finding = %#v, want one line-1 citation of %s", finding, citationSource)
		}
	}
}

func TestRelocationCitationsReportARelocatedFilesOutwardLink(t *testing.T) {
	repo := newCitationTestRepository(t)
	trackCitationTestFile(t, repo, citationSource, "[current](current.md)\n")
	trackCitationTestFile(t, repo, "docs/adr/current.md", "current\n")

	findings := citationTestFindings(t, repo, citationTestMove(citationSource), nil)
	if len(findings) != 1 {
		t.Fatalf("relocationCitationFindings() = %#v, want one finding", findings)
	}
	want := "line 1 cites current.md, which resolves to docs/adr/current.md; after this plan it resolves to docs/history/adr/current.md, where nothing exists"
	if findings[0].Path != citationSource || findings[0].Message != want {
		t.Errorf("finding = %#v, want path %q and message %q", findings[0], citationSource, want)
	}
}

func TestRelocationCitationsSkipCoRelocatedLinks(t *testing.T) {
	repo := newCitationTestRepository(t)
	trackCitationTestFile(t, repo, citationSource, "[peer](peer.md)\n")
	trackCitationTestFile(t, repo, "docs/adr/peer.md", "peer\n")
	moves := []HistoryMove{
		{Ordinal: 0, From: citationSource, To: "docs/history/adr/old.md"},
		{Ordinal: 1, From: "docs/adr/peer.md", To: "docs/history/adr/peer.md"},
	}

	if findings := citationTestFindings(t, repo, moves, nil); len(findings) != 0 {
		t.Fatalf("relocationCitationFindings() = %#v, want no finding", findings)
	}
}

func TestRelocationCitationsSkipAlreadyBrokenCitations(t *testing.T) {
	repo := newCitationTestRepository(t)
	trackCitationTestFile(t, repo, citationSource, "retired\n")
	trackCitationTestFile(t, repo, "guide.md", "[missing](docs/adr/missing.md)\n")

	if findings := citationTestFindings(t, repo, citationTestMove(citationSource), nil); len(findings) != 0 {
		t.Fatalf("relocationCitationFindings() = %#v, want no finding", findings)
	}
}

func TestRelocationCitationsSkipCodeFencesAndURLs(t *testing.T) {
	repo := newCitationTestRepository(t)
	trackCitationTestFile(t, repo, citationSource, "retired\n")
	trackCitationTestFile(t, repo, "docs/guide.md", "```markdown\n[old](../adr/old.md)\n```\n[url](https://example.test/docs/adr/old.md)\n")

	if findings := citationTestFindings(t, repo, citationTestMove(citationSource), nil); len(findings) != 0 {
		t.Fatalf("relocationCitationFindings() = %#v, want no finding", findings)
	}
}

func TestRelocationCitationsNeverOpenUntrackedOrIgnoredFiles(t *testing.T) {
	repo := newCitationTestRepository(t)
	trackCitationTestFile(t, repo, citationSource, "retired\n")
	trackCitationTestFile(t, repo, ".gitignore", "ignored.md\n")
	writeCitationTestFile(t, repo, "untracked.md", "docs/adr/old.md\n")
	writeCitationTestFile(t, repo, "ignored.md", "docs/adr/old.md\n")
	for _, path := range []string{"untracked.md", "ignored.md"} {
		absolute := filepath.Join(repo, path)
		if err := os.Chmod(absolute, 0); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = os.Chmod(absolute, 0o600) })
	}

	if findings := citationTestFindings(t, repo, citationTestMove(citationSource), nil); len(findings) != 0 {
		t.Fatalf("relocationCitationFindings() = %#v, want no finding or unscanned summary", findings)
	}
}

func TestRelocationCitationsNeverFollowSymbolicLinks(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symbolic-link replacement requires Unix permissions")
	}
	repo := newCitationTestRepository(t)
	trackCitationTestFile(t, repo, citationSource, "retired\n")
	trackCitationTestFile(t, repo, "linked.md", "placeholder\n")
	trackCitationTestFile(t, repo, "linked-dir/citation.md", "placeholder\n")
	writeCitationTestFile(t, repo, "outside.md", "docs/adr/old.md\n")
	if err := os.Remove(filepath.Join(repo, "linked.md")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("outside.md", filepath.Join(repo, "linked.md")); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(filepath.Join(repo, "linked-dir"), filepath.Join(repo, "real-dir")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("real-dir", filepath.Join(repo, "linked-dir")); err != nil {
		t.Fatal(err)
	}
	writeCitationTestFile(t, repo, "real-dir/citation.md", "docs/adr/old.md\n")

	if findings := citationTestFindings(t, repo, citationTestMove(citationSource), nil); len(findings) != 0 {
		t.Fatalf("relocationCitationFindings() = %#v, want no finding or unscanned summary", findings)
	}
}

func TestRelocationCitationsSummarizeUnscannedFiles(t *testing.T) {
	repo := newCitationTestRepository(t)
	trackCitationTestFile(t, repo, citationSource, "retired\n")
	trackCitationTestFile(t, repo, "binary.dat", "\x00docs/adr/old.md\n")
	trackCitationTestFile(t, repo, "large.txt", strings.Repeat("x", 4*1024*1024+1))

	findings := citationTestFindings(t, repo, citationTestMove(citationSource), nil)
	want := []Finding{{
		Code:    "baseline.history.citation.unscanned",
		Path:    ".",
		Message: "1 tracked files were not scanned for citations (larger than 4 MiB or unreadable): large.txt",
	}}
	if !reflect.DeepEqual(findings, want) {
		t.Fatalf("relocationCitationFindings() = %#v, want %#v", findings, want)
	}
}

func TestRelocationCitationsCapEachFileAndThePlan(t *testing.T) {
	repo := newCitationTestRepository(t)
	trackCitationTestFile(t, repo, citationSource, "retired\n")
	for index := 0; index < 201; index++ {
		content := "docs/adr/old.md\n"
		if index == 0 {
			content = strings.Repeat("docs/adr/old.md\n", 4)
		}
		trackCitationTestFile(t, repo, fmt.Sprintf("citations/%03d.txt", index), content)
	}

	findings := citationTestFindings(t, repo, citationTestMove(citationSource), nil)
	if len(findings) != 201 {
		t.Fatalf("relocationCitationFindings() returned %d findings, want 201", len(findings))
	}
	if !strings.HasSuffix(findings[0].Message, "; and 1 more") {
		t.Errorf("first finding message = %q, want three-citation cap", findings[0].Message)
	}
	wantOmitted := Finding{
		Code:    "baseline.history.citation.omitted",
		Path:    ".",
		Message: "1 more tracked files cite paths this plan relocates and are not listed",
	}
	if findings[len(findings)-1] != wantOmitted {
		t.Errorf("last finding = %#v, want %#v", findings[len(findings)-1], wantOmitted)
	}
}

func TestRelocationCitationsCountUnitDirectoriesNotFamilyRoots(t *testing.T) {
	repo := newCitationTestRepository(t)
	from := "docs/specs/_archived/0012-old/_prd.md"
	reviewFrom := "docs/specs/_reviews/pr-12/report.md"
	moves := []HistoryMove{
		{From: from, To: "docs/history/specs/0012-old/_prd.md"},
		{Ordinal: 1, From: reviewFrom, To: "docs/history/reviews/pr-12/report.md"},
	}
	trackCitationTestFile(t, repo, from, "retired\n")
	trackCitationTestFile(t, repo, reviewFrom, "review\n")
	trackCitationTestFile(t, repo, "unit.txt", "docs/specs/_archived/0012-old\n")
	trackCitationTestFile(t, repo, "family.txt", "docs/specs/_archived\n")
	trackCitationTestFile(t, repo, "review-unit.txt", "docs/specs/_reviews/pr-12\n")
	trackCitationTestFile(t, repo, "review-family.txt", "docs/specs/_reviews\n")

	findings := citationTestFindings(t, repo, moves, nil)
	if len(findings) != 2 || findings[0].Path != "review-unit.txt" || findings[1].Path != "unit.txt" ||
		!strings.Contains(findings[0].Message, "(it moves to docs/history/reviews/pr-12)") ||
		!strings.Contains(findings[1].Message, "(it moves to docs/history/specs/0012-old)") {
		t.Fatalf("relocationCitationFindings() = %#v, want only the relocated unit directories", findings)
	}
}

func TestRelocationCitationsDoNothingWithoutMoves(t *testing.T) {
	findings, err := relocationCitationFindings(t.Context(), filepath.Join(t.TempDir(), "does-not-exist"), nil, nil)
	if err != nil || findings != nil {
		t.Fatalf("relocationCitationFindings() = (%#v, %v), want (nil, nil)", findings, err)
	}
}

func TestRelocationCitationsAreDeterministic(t *testing.T) {
	repo := newCitationTestRepository(t)
	trackCitationTestFile(t, repo, citationSource, "retired\n")
	trackCitationTestFile(t, repo, "b.txt", "docs/adr/old.md\n")
	trackCitationTestFile(t, repo, "a.txt", "docs/adr/old.md\n")
	first := citationTestFindings(t, repo, citationTestMove(citationSource), nil)
	second := citationTestFindings(t, repo, citationTestMove(citationSource), nil)
	firstJSON, err := json.Marshal(first)
	if err != nil {
		t.Fatal(err)
	}
	secondJSON, err := json.Marshal(second)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(firstJSON, secondJSON) {
		t.Fatalf("two scans differ:\n%s\n%s", firstJSON, secondJSON)
	}
}

func TestRelocationCitationsNeverPrintAControlCharacter(t *testing.T) {
	repo := newCitationTestRepository(t)
	trackCitationTestFile(t, repo, citationSource, "retired\n")
	trackCitationTestFile(t, repo, "control\nname.md", "docs/adr/old.md\n")
	trackCitationTestFile(t, repo, "encoded.md", "[bad](docs/adr/old.md%0Asecret)\n")

	findings := citationTestFindings(t, repo, citationTestMove(citationSource), nil)
	if len(findings) != 0 {
		t.Fatalf("relocationCitationFindings() = %#v, want no control-bearing finding", findings)
	}
}

func TestRelocationCitationsIgnoreAMoveApplyWouldRefuse(t *testing.T) {
	repo := newCitationTestRepository(t)
	trackCitationTestFile(t, repo, citationSource, "retired\n")
	trackCitationTestFile(t, repo, "guide.md", "docs/adr/old.md\n")
	writeCitationTestFile(t, repo, "docs/history/adr/old.md", "occupied and untracked\n")
	refused := map[string]bool{citationSource: true}

	if findings := citationTestFindings(t, repo, citationTestMove(citationSource), refused); len(findings) != 0 {
		t.Fatalf("relocationCitationFindings() = %#v, want no finding", findings)
	}
}

func TestRelocationCitationsNeverBlockOnAFIFO(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("FIFO replacement requires Unix")
	}
	repo := newCitationTestRepository(t)
	trackCitationTestFile(t, repo, citationSource, "retired\n")
	trackCitationTestFile(t, repo, "citation.pipe", "docs/adr/old.md\n")
	absolute := filepath.Join(repo, "citation.pipe")
	if err := os.Remove(absolute); err != nil {
		t.Fatal(err)
	}
	command := exec.Command("mkfifo", absolute)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("mkfifo: %v: %s", err, output)
	}

	if findings := citationTestFindings(t, repo, citationTestMove(citationSource), nil); len(findings) != 0 {
		t.Fatalf("relocationCitationFindings() = %#v, want no finding", findings)
	}
}

func citationTestMove(from string) []HistoryMove {
	return []HistoryMove{{Ordinal: 0, From: from, To: "docs/history/adr/old.md"}}
}

func citationTestFindings(t *testing.T, repo string, moves []HistoryMove, refused map[string]bool) []Finding {
	t.Helper()
	findings, err := relocationCitationFindings(context.Background(), repo, moves, refused)
	if err != nil {
		t.Fatalf("relocationCitationFindings() error = %v", err)
	}
	return findings
}

func newCitationTestRepository(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	runCitationTestGit(t, repo, "init", "--quiet")
	runCitationTestGit(t, repo, "config", "user.email", "tests@example.test")
	runCitationTestGit(t, repo, "config", "user.name", "Roundfix Tests")
	return repo
}

func trackCitationTestFile(t *testing.T, repo string, relative string, content string) {
	t.Helper()
	writeCitationTestFile(t, repo, relative, content)
	runCitationTestGit(t, repo, "add", "--", relative)
}

func writeCitationTestFile(t *testing.T, repo string, relative string, content string) {
	t.Helper()
	absolute := filepath.Join(repo, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(absolute), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(absolute, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func runCitationTestGit(t *testing.T, repo string, args ...string) {
	t.Helper()
	gitArgs := append([]string{"-C", repo, "-c", "core.fsmonitor=false"}, args...)
	command := exec.Command("git", gitArgs...)
	command.Env = append(os.Environ(), "GIT_CONFIG_NOSYSTEM=1", "GIT_OPTIONAL_LOCKS=0", "GIT_TERMINAL_PROMPT=0")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v: %s", strings.Join(args, " "), err, output)
	}
}
