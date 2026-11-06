// Suite: Data Review Segment planning from real Git diffs.
// Invariant: cuts preserve every changed line and its old and new file coordinates.
// Boundary IN: file packing, hunk and line cuts, byte fits, scope reasons and prompts.
// Boundary OUT: provider sessions, review records and answer classification.
package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
)

func TestReviewDataSegmentsPackCutAndBound(t *testing.T) {
	t.Parallel()
	const bound = 4096
	fit := func(promptBytes, diffBytes int) bool { return diffBytes <= bound }
	t.Run("small files share one segment", func(t *testing.T) {
		t.Parallel()
		root, base := newReviewDataRepo(t, nil)
		head, diff := commitReviewDataFiles(t, root, base, map[string]string{"a.json": "alpha\n", "b.json": "beta\n"})
		segments, diffs, prompts, err := planReviewDataSegments(diff, base, head, fit)
		if err != nil {
			t.Fatal(err)
		}
		if len(segments) != 1 || !reflect.DeepEqual(segments[0].Paths, []string{"a.json", "b.json"}) || diffs[0] != diff || segments[0].DiffBytes != len(diff) {
			t.Fatalf("segments=%+v diffs=%q", segments, diffs)
		}
		want := fmt.Sprintf("Data Review Segment 1 of 1: pre_pr_review.data_paths changes %s..%s of the candidate; answer for this segment only.\n", base, head) +
			"These files are data the project writes by hand and declares in pre_pr_review.data_paths, such as rule catalogs. Check each changed value against the Spec contexts, the code that reads it and the rest of the file; report a wrong, missing, duplicated or contradictory value.\n" + buildReviewPrompt(base, head, diff)
		if len(prompts) != 1 || prompts[0] != want {
			t.Fatalf("prompts=%q, want %q", prompts, want)
		}
	})
	t.Run("a new file is cut at lines", func(t *testing.T) {
		t.Parallel()
		root, base := newReviewDataRepo(t, nil)
		body := strings.TrimSuffix(reviewDataNumberedLines(400), "\n")
		head, diff := commitReviewDataFiles(t, root, base, map[string]string{"rules.json": body})
		assertReviewDataCuts(t, root, base, head, "rules.json", diff, fit)
	})
	t.Run("a modified file is cut at hunks and lines", func(t *testing.T) {
		t.Parallel()
		original := reviewDataNumberedLines(600)
		root, base := newReviewDataRepo(t, map[string]string{"rules.json": original})
		lines := strings.SplitAfter(original, "\n")
		lines[19] = "changed rule twenty\n"
		// Deletions and a large insertion force both nonzero and zero hunk counts.
		body := strings.Join(lines[:99], "") + strings.Join(lines[109:299], "") + strings.ReplaceAll(reviewDataNumberedLines(200), "rule", "inserted") + strings.Join(lines[299:], "")
		head, diff := commitReviewDataFiles(t, root, base, map[string]string{"rules.json": body})
		if strings.Count(diff, "\n@@ ") < 3 {
			t.Fatalf("fixture needs distinct hunks: %s", diff)
		}
		assertReviewDataCuts(t, root, base, head, "rules.json", diff, fit)
	})
	t.Run("a line that fits no segment is refused", func(t *testing.T) {
		t.Parallel()
		root, base := newReviewDataRepo(t, nil)
		head, diff := commitReviewDataFiles(t, root, base, map[string]string{"rules.json": strings.Repeat("x", bound) + "\n"})
		segments, diffs, prompts, err := planReviewDataSegments(diff, base, head, fit)
		want := "review data too large: line 1 of rules.json does not fit one data review segment"
		if err == nil || err.Error() != want || segments != nil || diffs != nil || prompts != nil {
			t.Fatalf("segments=%+v error=%v, want %q", segments, err, want)
		}
	})
	t.Run("a fifth segment is refused", func(t *testing.T) {
		t.Parallel()
		root, base := newReviewDataRepo(t, nil)
		files := map[string]string{}
		for i := range 5 {
			files[fmt.Sprintf("%d.json", i)] = strings.Repeat("x", 3000) + "\n"
		}
		head, diff := commitReviewDataFiles(t, root, base, files)
		segments, diffs, prompts, err := planReviewDataSegments(diff, base, head, fit)
		want := fmt.Sprintf("review data too large: %d bytes of pre_pr_review.data_paths needs 5 data review segments; a review allows at most %d", len(diff), reviewSegmentLimit)
		if err == nil || err.Error() != want || segments != nil || diffs != nil || prompts != nil {
			t.Fatalf("segments=%+v error=%v, want %q", segments, err, want)
		}
	})
	t.Run("prompt byte fit includes both headers", func(t *testing.T) {
		t.Parallel()
		root, base := newReviewDataRepo(t, nil)
		head, diff := commitReviewDataFiles(t, root, base, map[string]string{"rules.json": reviewDataNumberedLines(400)})
		promptFit := func(promptBytes, diffBytes int) bool { return promptBytes <= 11000 }
		assertReviewDataCuts(t, root, base, head, "rules.json", diff, promptFit)
	})
}

func TestReviewDataPathScopeAndPrompt(t *testing.T) {
	t.Parallel()
	root, base := newReviewDataRepo(t, map[string]string{".roundfixrc.yml": "pre_pr_review:\n  data_paths: ['rules/']\n  omit_paths: ['rules/omit.json']\n"})
	head, _ := commitReviewDataFiles(t, root, base, map[string]string{"code.go": "package code\n", "rules/a.json": "authored alpha\n", "rules/z.json": "authored zeta\n", "rules/omit.json": "omitted payload\n", ".roundfixrc.yml": "pre_pr_review:\n  data_paths: []\n"})
	diff, omitted, err := reviewScopedDiff(t.Context(), root, base, head, nil, preflight.ExecGitRunner{})
	if err != nil {
		t.Fatal(err)
	}
	want := []reviewOmittedPath{{Path: "rules/a.json", Reason: "data-path"}, {Path: "rules/omit.json", Reason: "omit-path"}, {Path: "rules/z.json", Reason: "data-path"}}
	if !reflect.DeepEqual(omitted, want) || strings.Contains(diff, "authored") || strings.Contains(diff, "omitted payload") || !strings.Contains(diff, "+package code") {
		t.Fatalf("omitted=%+v diff=%s", omitted, diff)
	}
	// Deliberately unsorted: prompt ordering is its own contract.
	got := appendReviewOmissions("", []reviewOmittedPath{want[2], want[1], want[0]})
	expected := "Omitted from this diff (not reviewed): 0 path(s) of QA evidence, 0 of upstream-managed skills, 0 of archived Spec sources, 0 of retired history, 0 of generated files, 1 of pre_pr_review.omit_paths.\nrules/omit.json (omit-path)\nReviewed apart, after this diff, in Data Review Segments: 2 path(s) of pre_pr_review.data_paths.\nrules/a.json (data-path)\nrules/z.json (data-path)\n"
	if got != expected {
		t.Fatalf("prompt=%q, want %q", got, expected)
	}
	// A later round keeps declarations at the original base even though the
	// earlier candidate cleared both lists in Project Config.
	next, _ := commitReviewDataFiles(t, root, head, map[string]string{"rules/a.json": "updated data\n", "rules/omit.json": "updated omitted\n"})
	diff, omitted, err = reviewScopedDiff(t.Context(), root, head, next, nil, preflight.ExecGitRunner{}, base)
	if err != nil || !reflect.DeepEqual(omitted, want[:2]) || diff != "" {
		t.Fatalf("later round diff=%q omitted=%+v error=%v", diff, omitted, err)
	}
	// Without the pinned declaration base, the cleared key no longer separates
	// these files; this makes the preceding assertion sensitive to base drift.
	diff, omitted, err = reviewScopedDiff(t.Context(), root, head, next, nil, preflight.ExecGitRunner{})
	if err != nil || len(omitted) != 0 || !strings.Contains(diff, "+updated data") {
		t.Fatalf("cleared declarations diff=%q omitted=%+v error=%v", diff, omitted, err)
	}
}

func newReviewDataRepo(t *testing.T, files map[string]string) (string, string) {
	t.Helper()
	root := t.TempDir()
	gittest.InitRepo(t, root)
	for name, body := range files {
		validatorWrite(t, root, name, body)
	}
	gittest.Run(t, root, "add", ".")
	gittest.Run(t, root, "commit", "--allow-empty", "-m", "base data")
	return root, strings.TrimSpace(gittest.Run(t, root, "rev-parse", "HEAD"))
}

func commitReviewDataFiles(t *testing.T, root, base string, files map[string]string) (string, string) {
	t.Helper()
	for name, body := range files {
		file := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(file, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	gittest.Run(t, root, "add", ".")
	gittest.Run(t, root, "commit", "-m", "candidate data")
	head := strings.TrimSpace(gittest.Run(t, root, "rev-parse", "HEAD"))
	return head, gittest.Run(t, root, "diff", "--no-ext-diff", "--no-textconv", "--no-color", base, head, "--")
}

func reviewDataNumberedLines(count int) string {
	var out strings.Builder
	for i := range count {
		fmt.Fprintf(&out, "rule %04d value %04d\n", i+1, i+1)
	}
	return out.String()
}

func assertReviewDataCuts(t *testing.T, root, base, head, path, diff string, fit reviewDataFit) {
	t.Helper()
	segments, diffs, prompts, err := planReviewDataSegments(diff, base, head, fit)
	if err != nil {
		t.Fatal(err)
	}
	if len(segments) < 2 || len(segments) > reviewSegmentLimit || len(diffs) != len(segments) || len(prompts) != len(segments) {
		t.Fatalf("segments=%+v diffs=%d prompts=%d", segments, len(diffs), len(prompts))
	}
	oldText := gittest.Run(t, root, "ls-tree", "--name-only", base, "--", path)
	var oldLines []string
	if strings.TrimSpace(oldText) != "" {
		oldLines = strings.Split(gittest.Run(t, root, "show", base+":"+path), "\n")
	}
	newLines := strings.Split(gittest.Run(t, root, "show", head+":"+path), "\n")
	// Preserve the sequence, multiplicity and no-newline markers as well as anchors.
	payload := func(diff string) []string {
		var result []string
		inHunk := false
		for _, line := range strings.Split(diff, "\n") {
			if strings.HasPrefix(line, "diff --git ") {
				inHunk = false
			}
			if strings.HasPrefix(line, "@@ ") {
				inHunk = true
				continue
			}
			if inHunk && line != "" {
				result = append(result, line)
			}
		}
		return result
	}
	var gotPayload []string
	for i, piece := range diffs {
		if !fit(len(prompts[i]), len(piece)) || segments[i].DiffBytes != len(piece) || !reflect.DeepEqual(segments[i].Paths, []string{path}) || !strings.HasPrefix(piece, "diff --git a/"+path+" b/"+path+"\n") {
			t.Fatalf("piece %d does not fit or has wrong header/metadata: %+v, bytes=%d prompt=%d", i, segments[i], len(piece), len(prompts[i]))
		}
		gotPayload = append(gotPayload, payload(piece)...)
		old, newLine, oldCount, newCount := 0, 0, 0, 0
		inHunk := false
		checkCounts := func() {
			t.Helper()
			if oldCount != 0 || newCount != 0 {
				t.Fatalf("piece %d hunk counts left: old=%d new=%d", i, oldCount, newCount)
			}
		}
		for _, line := range strings.Split(piece, "\n") {
			if strings.HasPrefix(line, "@@ ") {
				if inHunk {
					checkCounts()
				}
				if _, err := fmt.Sscanf(line, "@@ -%d,%d +%d,%d @@", &old, &oldCount, &newLine, &newCount); err != nil {
					t.Fatalf("recomputed header %q: %v", line, err)
				}
				if oldCount == 0 {
					old++
				}
				if newCount == 0 {
					newLine++
				}
				inHunk = true
				continue
			}
			if !inHunk || line == "" || line[0] == '\\' {
				continue
			}
			if line[0] != '+' {
				if old < 1 || old > len(oldLines) || oldLines[old-1] != line[1:] {
					t.Fatalf("piece %d old line %d=%q differs from Git", i, old, line)
				}
				old++
				oldCount--
			}
			if line[0] != '-' {
				if newLine < 1 || newLine > len(newLines) || newLines[newLine-1] != line[1:] {
					t.Fatalf("piece %d new line %d=%q differs from Git", i, newLine, line)
				}
				newLine++
				newCount--
			}
		}
		checkCounts()
	}
	if want := payload(diff); !reflect.DeepEqual(gotPayload, want) {
		t.Fatalf("cut loses or duplicates payload: got %d lines, want %d", len(gotPayload), len(want))
	}
}
