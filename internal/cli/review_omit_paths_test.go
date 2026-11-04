// Boundary: base-commit review omissions, provider prompts, and persisted scope.
package cli

import (
	"crypto/sha256"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/preflight"
)

func TestReviewOmitsConfiguredPathsFromTheBaseCommit(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name     string
		atBase   bool
		derived  bool
		roundTwo bool
	}{
		{name: "base glob", atBase: true},
		{name: "head-only key"},
		{name: "derived reason wins", atBase: true, derived: true},
		{name: "round two keeps declaration base", atBase: true, roundTwo: true},
		{name: "round two ignores head-only key", roundTwo: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			answers := []string{"No findings."}
			if tc.roundTwo {
				answers = []string{"Findings:\n- src/a.go:1 Failure: candidate defect", "No findings."}
			}
			runner := newLineageRunner(answers...)
			fixture := newReviewCommandFixture(t, "codex", runner)
			config := fmt.Sprintf("defaults:\n  artifact_dir: %q\npre_pr_review:\n  provider: codex\n  omit_paths: ['gen/*_snapshot.json']\n", fixture.artifactDir)
			if tc.derived {
				config += "delivery:\n  derived_paths:\n    - paths: ['gen/*_snapshot.json']\n      regenerate: make generate\n"
			}
			commit := commitScopeFile(t, fixture, ".roundfixrc.yml", config)
			if tc.atBase {
				fixture.baseCommit = commit
			}
			commitScopeFile(t, fixture, "gen/a_snapshot.json", "round one snapshot payload\n")
			commitScopeFile(t, fixture, "src/a.go", "package reviewed\n")
			code, record, stderr := fixture.run(t)
			if tc.roundTwo {
				if code != exitRunFailed || record.Lineage.Round != 1 {
					t.Fatalf("round one exit=%d record=%+v stderr=%q", code, record, stderr)
				}
			} else if code != exitOK {
				t.Fatalf("exit=%d record=%+v stderr=%q", code, record, stderr)
			}
			want := []reviewOmittedPath{}
			if tc.atBase {
				reason := "omit-path"
				if tc.derived {
					reason = "generated"
				}
				want = append(want, reviewOmittedPath{Path: "gen/a_snapshot.json", Reason: reason})
			}
			if !reflect.DeepEqual(record.OmittedPaths, want) {
				t.Fatalf("round one omissions=%+v, want %+v", record.OmittedPaths, want)
			}
			if tc.roundTwo {
				// The head clears a base declaration, or retains a head-only one.
				// Neither may change the immutable declaration base for the delta.
				if tc.atBase {
					commitScopeFile(t, fixture, ".roundfixrc.yml", strings.Replace(config, "['gen/*_snapshot.json']", "[]", 1))
				}
				commitScopeFile(t, fixture, "gen/a_snapshot.json", "round two snapshot payload\n")
				commitScopeFile(t, fixture, "src/a.go", "package reviewed\n// round two source\n")
				code, record, stderr = fixture.run(t)
				if code != exitOK || record.Lineage.Round != 2 {
					t.Fatalf("round two exit=%d record=%+v stderr=%q", code, record, stderr)
				}
			}
			if !reflect.DeepEqual(record.OmittedPaths, want) {
				t.Fatalf("omissions=%+v, want %+v", record.OmittedPaths, want)
			}
			start, end := "--- BEGIN CANDIDATE DIFF ---\n", "--- END CANDIDATE DIFF ---\n"
			if tc.roundTwo {
				start, end = "--- BEGIN ROUND 2 DELTA DIFF ---\n", "--- END ROUND 2 DELTA DIFF ---\n"
			}
			_, body, found := strings.Cut(runner.request.Prompt, start)
			if !found {
				t.Fatalf("missing diff marker %q", start)
			}
			body, after, found := strings.Cut(body, end)
			if !found {
				t.Fatalf("missing diff marker %q", end)
			}
			if strings.Contains(body, "snapshot payload") == tc.atBase || !strings.Contains(body, "diff --git a/src/a.go b/src/a.go") {
				t.Fatalf("unexpected reviewed diff: %s", body)
			}
			if record.DiffBytes != len(strings.TrimSuffix(body, "\n")) {
				t.Fatalf("diffBytes=%d, prompt diff=%d", record.DiffBytes, len(strings.TrimSuffix(body, "\n")))
			}
			wantSummary := "Omitted from this diff (not reviewed): 0 path(s) of QA evidence, 0 of upstream-managed skills, 0 of archived Spec sources, 0 of retired history, 0 of generated files.\n"
			if tc.derived {
				wantSummary = strings.Replace(wantSummary, "0 of generated files", "1 of generated files", 1)
			} else if tc.atBase {
				wantSummary = strings.TrimSuffix(wantSummary, ".\n") + ", 1 of pre_pr_review.omit_paths.\n"
			}
			wantSummary = strings.TrimSuffix(wantSummary, "\n")
			summary, _, _ := strings.Cut(after, "\n")
			if summary != wantSummary {
				t.Fatalf("summary=%q, want %q", summary, wantSummary)
			}
			for _, omitted := range want {
				if !strings.Contains(after, omitted.Path+" ("+omitted.Reason+")\n") {
					t.Fatalf("omission not recorded in prompt: %q", after)
				}
			}
		})
	}
}

func TestReviewOmissionSummaryKeepsItsBytesWithoutOmitPaths(t *testing.T) {
	t.Parallel()

	const summary = "Omitted from this diff (not reviewed): 0 path(s) of QA evidence, 0 of upstream-managed skills, 0 of archived Spec sources, 0 of retired history, 0 of generated files."
	if got := appendReviewOmissions("existing prompt\n", nil); got != "existing prompt\n"+summary+"\n" {
		t.Fatalf("summary changed: %q", got)
	}
	prompt := buildReviewPrompt("base", "head", "unchanged diff")
	if !strings.HasSuffix(prompt, summary+"\n") || prompt != buildReviewPrompt("base", "head", "unchanged diff", []reviewOmittedPath{}) {
		t.Fatalf("prompt without omissions changed: %q", prompt)
	}
	t.Logf("prompt SHA-256: %x", sha256.Sum256([]byte(prompt)))

	runner := newLineageRunner("No findings.")
	fixture := newReviewCommandFixture(t, "codex", runner)
	head := commitScopeFile(t, fixture, "src/a.go", "package reviewed\n")
	code, record, stderr := fixture.run(t)
	if code != exitOK || len(record.OmittedPaths) != 0 {
		t.Fatalf("exit=%d record=%+v stderr=%q", code, record, stderr)
	}
	diff, err := (preflight.ExecGitRunner{}).RunGit(t.Context(), fixture.repository, "diff", "--no-ext-diff", "--no-textconv", "--no-color", fixture.baseCommit, head, "--")
	if err != nil {
		t.Fatal(err)
	}
	if want := buildReviewPrompt(fixture.baseCommit, head, diff); runner.request.Prompt != want {
		t.Fatalf("prompt without key changed: got %q, want %q", runner.request.Prompt, want)
	}
	persisted, err := readReviewRecord(filepath.Join(reviewCheckoutDir(fixture.artifactDir, fixture.repository), reviewRecordFileName))
	if err != nil || len(persisted.OmittedPaths) != 0 {
		t.Fatalf("persisted=%+v error=%v", persisted, err)
	}
}
