package cli

import (
	"bytes"
	"context"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/baseline"
)

func TestAnExceptionWithoutASourceIsRefused(t *testing.T) {
	t.Parallel()
	repository := languageRulePriorRepository(t)
	before := baselinePlanTestTree(t, repository)
	_, stdout, stderr, code := runBaselineUpdateTestCommand(t, context.Background(), "baseline", "update", "--repo", repository, "--no-skills", "--no-history", "--format", "text", "--decision", "language.reader=pt-BR", "--decision", `language.code_exception={"language":"Portuguese"}`)
	for _, label := range []string{"Prior catalog:", "Current catalog:"} {
		digest := baselineUpdateTextValue(stdout, label)
		if len(digest) != len("sha256:")+64 || !strings.HasPrefix(digest, "sha256:") {
			t.Fatalf("invalid digest %q", digest)
		}
		stdout = strings.ReplaceAll(stdout, digest, "<digest>")
	}
	const reason = `declare Setup Manifest decisions: resolve Setup Manifest decisions: normalize Baseline decision "language.code_exception": requires a source, the recorded decision of the owner that keeps the codebase's language: no source, no exception`
	const want = `Baseline update: failed
Category: invalid
Result: ` + reason + `
Prior catalog: <digest>
Current catalog: <digest>
File changes: 0
Misplaced repository rules: 0
Retention evidence: 0
Declared decisions: 0
Skills: not run
Skills installed: 0
Skills restored: 0
Skills drifted: 0
Next action: correct the declared decisions and rerun roundfix baseline update
`
	if stdout != want || stderr != "roundfix: baseline update failed: "+reason+"\n" || code != 2 {
		t.Fatalf("Surface Transcript 2: exit=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	if baselinePlanTestTree(t, repository) != before {
		t.Fatal("invalid declaration changed repository")
	}
}

func TestTheInteractiveFlowAsksNothingForTheCodeLanguageException(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"nothing_recorded", "a_recorded_exception_is_carried"} {
		t.Run(name, func(t *testing.T) {
			repository := newCLIProjectDecisionRepository(t)
			catalog, err := baseline.LoadEmbeddedCatalog()
			if err != nil {
				t.Fatal(err)
			}
			profile, err := baseline.ResolveProfile(repository, "standard-typescript-monorepo", catalog)
			if err != nil {
				t.Fatal(err)
			}
			value := map[string]any{"language": "Portuguese", "source": "docs/adr/0141-example.md"}
			state := baselineHumanState{currentDecisions: map[string]any{}}
			if name == "a_recorded_exception_is_carried" {
				state.currentDecisions["language.code_exception"] = value
			}
			var prompts bytes.Buffer
			answers, err := promptBaselineDecisions(context.Background(), &baselineHumanPrompt{reader: bufioReader(projectDecisionHumanAnswers()), writer: &prompts}, catalog, profile, state)
			if err != nil {
				t.Fatal(err)
			}
			if strings.Contains(prompts.String(), "language.code_exception") || strings.Contains(prompts.String(), "Code Language Exception") {
				t.Fatalf("exception prompted: %s", prompts.String())
			}
			found := false
			for _, answer := range answers {
				if answer.ID == "language.code_exception" {
					found = true
					if !reflect.DeepEqual(answer.Value, value) {
						t.Fatalf("carried = %#v", answer.Value)
					}
				}
			}
			if found != (name == "a_recorded_exception_is_carried") {
				t.Fatalf("exception recorded = %t", found)
			}
		})
	}
}
