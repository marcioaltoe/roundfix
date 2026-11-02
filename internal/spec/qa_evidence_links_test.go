package spec

// Suite: shared QA evidence links and history extraction.
// Invariant: only existing tracked or addable local targets resolve in Git.
// Boundary IN: Markdown, disposable Git indexes and disk targets.
// Boundary OUT: QA verdict policy and staging, owned by the callers.

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"roundfix/internal/gittest"
)

func TestUnresolvedEvidenceLinksNameTheReason(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	root := t.TempDir()
	gittest.InitRepo(t, root, "-b", "main")
	root, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	report := filepath.Join(root, "docs/specs/x/qa/report.md")
	writeFile(t, filepath.Join(root, ".gitignore"), "*.log\n")
	writeFile(t, filepath.Join(filepath.Dir(report), "evidence/a.log"), "ignored\n")
	writeFile(t, filepath.Join(filepath.Dir(report), "evidence/c.log"), "tracked\n")
	writeFile(t, filepath.Join(filepath.Dir(report), "evidence/d.txt"), "untracked\n")
	writeFile(t, filepath.Join(filepath.Dir(report), "evidence/only-ignored/a.log"), "ignored\n")
	writeFile(t, filepath.Join(filepath.Dir(report), "evidence/literal1.txt"), "glob decoy\n")
	writeFile(t, filepath.Join(filepath.Dir(report), "evidence/literal[1].txt"), "literal evidence\n")
	writeFile(t, filepath.Join(root, "../outside.txt"), "outside repository\n")
	gittest.Run(t, root, "add", ".gitignore")
	gittest.Run(t, root, "add", "-f", "docs/specs/x/qa/evidence/c.log")
	gittest.Run(t, root, "commit", "-m", "track ignore rules and forced log")
	tests := []struct {
		name, content string
		want          []UnresolvedEvidenceLink
	}{
		{"document order and lines", "[ignored](evidence/a.log)\n[missing](evidence/b.txt)\n[tracked](evidence/c.log)\n[addable](evidence/d.txt)\n[directory](evidence/)\n[outside](../../../../../outside.txt)\n`[code](evidence/code.txt)`\n[web](https://example.com/a)\n[anchor](#anchor)\n[self](report.md)\n[absolute](/absent.txt)\n[ref]: evidence/reference.txt\n", []UnresolvedEvidenceLink{
			{1, "evidence/a.log", filepath.ToSlash(filepath.Join(filepath.Dir(report), "evidence/a.log")), EvidenceLinkIgnored},
			{2, "evidence/b.txt", filepath.ToSlash(filepath.Join(filepath.Dir(report), "evidence/b.txt")), EvidenceLinkMissing},
			{6, "../../../../../outside.txt", filepath.ToSlash(filepath.Join(root, "../outside.txt")), EvidenceLinkMissing},
			{12, "evidence/reference.txt", filepath.ToSlash(filepath.Join(filepath.Dir(report), "evidence/reference.txt")), EvidenceLinkMissing},
		}},
		{"single addable target", "[addable](evidence/d.txt)", nil},
		{"single tracked ignored target", "[tracked](evidence/c.log)", nil},
		{"literal pathspec", "[literal](evidence/literal%5B1%5D.txt)", nil},
		{"ignored directory", "[directory](evidence/only-ignored/)", []UnresolvedEvidenceLink{{1, "evidence/only-ignored/", filepath.ToSlash(filepath.Join(filepath.Dir(report), "evidence/only-ignored")), EvidenceLinkIgnored}}},
		{"directory and repository root", "[directory](evidence/)\n[root](../../../..)", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeFile(t, report, tt.content)
			got, err := UnresolvedEvidenceLinks(t.Context(), report)
			if err != nil || !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("links = %+v, err %v; want %+v", got, err, tt.want)
			}
			err = RequireResolvedEvidenceLinks(t.Context(), report)
			var refusal EvidenceLinksError
			if len(tt.want) == 0 {
				if err != nil {
					t.Fatalf("resolved report refused: %v", err)
				}
			} else if !errors.As(err, &refusal) || !reflect.DeepEqual(refusal.Links, tt.want) {
				t.Fatalf("refusal = %v, want %+v", err, tt.want)
			}
		})
	}
	// A failed listing and a corrupt worktree must stay infrastructure errors.
	writeFile(t, report, "[addable](evidence/d.txt)")
	writeFile(t, filepath.Join(root, ".git/index"), "invalid index\n")
	if links, err := UnresolvedEvidenceLinks(t.Context(), report); err == nil || links != nil {
		t.Fatalf("broken index = %+v, %v; want infrastructure error", links, err)
	}
	writeFile(t, filepath.Join(root, ".git/HEAD"), "invalid HEAD\n")
	if links, err := UnresolvedEvidenceLinks(t.Context(), report); err == nil || links != nil {
		t.Fatalf("broken Git = %+v, %v; want infrastructure error", links, err)
	}
}

func TestUnresolvedEvidenceLinksOutsideGitReadTheDisk(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	root := t.TempDir()
	root, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	report := filepath.Join(root, "report.md")
	writeFile(t, filepath.Join(root, "present.log"), "present\n")
	writeFile(t, report, "[present](present.log)\n[absent](absent.txt)\n[self](report.md)\n")
	want := []UnresolvedEvidenceLink{{2, "absent.txt", filepath.ToSlash(filepath.Join(root, "absent.txt")), EvidenceLinkMissing}}
	got, err := UnresolvedEvidenceLinks(t.Context(), report)
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("outside Git = %+v, %v; want %+v", got, err, want)
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if links, err := UnresolvedEvidenceLinks(ctx, report); err == nil || links != nil {
		t.Fatalf("cancelled read = %+v, %v; want infrastructure error", links, err)
	}
	if err := os.Remove(report); err != nil {
		t.Fatal(err)
	}
	if links, err := UnresolvedEvidenceLinks(t.Context(), report); !errors.Is(err, os.ErrNotExist) || links != nil {
		t.Fatalf("missing report = %+v, %v; want read error", links, err)
	}
}

func TestEvidenceLinksErrorNamesEachLink(t *testing.T) {
	t.Parallel()
	err := EvidenceLinksError{Links: []UnresolvedEvidenceLink{
		{Line: 9, Destination: "evidence/a.log", Reason: EvidenceLinkIgnored},
		{Line: 12, Destination: "evidence/b.txt", Reason: EvidenceLinkMissing},
	}}
	want := "QA Report evidence links do not resolve in the committed tree: evidence/a.log (line 9, ignored); evidence/b.txt (line 12, missing)"
	if got := err.Error(); got != want {
		t.Fatalf("refusal = %q, want %q", got, want)
	}
}

func TestMissingRelativeLinksReadsTheSharedExtraction(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	path := "docs/report.md"
	writeFile(t, filepath.Join(root, "docs/present.txt"), "present\n")
	content := []byte("[present](present.txt)\n[absent](absent.txt)\n[self](report.md)\n[outside](../../outside.txt)\n`[code](code.txt)`\n[reference]: absent-ref.txt\n")
	var filtered []MissingLink
	for _, link := range RelativeLinks(path, content) {
		if !historyRepositoryPath(link.Target) {
			continue
		}
		self := link.Target == path
		_, err := os.Stat(filepath.Join(root, link.Target))
		if self || os.IsNotExist(err) {
			filtered = append(filtered, MissingLink{Path: path, Line: link.Line, Destination: link.Destination, Target: link.Target, Self: self})
		}
	}
	want := []MissingLink{
		{path, 2, "absent.txt", "docs/absent.txt", false},
		{path, 3, "report.md", path, true},
		{path, 6, "absent-ref.txt", "docs/absent-ref.txt", false},
	}
	if got := MissingRelativeLinks(root, path, content); !reflect.DeepEqual(got, filtered) || !reflect.DeepEqual(got, want) {
		t.Fatalf("history links = %+v, filter %+v; want %+v", got, filtered, want)
	}
}
