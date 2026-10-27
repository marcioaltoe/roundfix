package baseline

import (
	"os"
	"strings"
	"testing"
)

func TestTheGuidesExemptANetworkProhibitedOutsideEvidenceRow(t *testing.T) {
	t.Parallel()
	const specWorkflowSentence = "So is an outside-evidence row the QA gate did not look up because the Spec's own authorization forbids that lookup, recorded as `blocked (environment: <subject> prohibited by Spec authorization)`, where `<subject>` names what the authorization forbids, such as `network` or `another repository`."
	const autonomousWorkSentence = "Nor does an outside-evidence row the QA gate did not look up because the Spec's own authorization forbids that lookup, recorded as `blocked (environment: <subject> prohibited by Spec authorization)`, where `<subject>` names what the authorization forbids, such as `network` or `another repository`."

	for _, want := range []struct {
		id       string
		sentence string
		paths    []string
	}{
		{
			id:       "clause.spec.project-constraints-06-outside-evidence",
			sentence: specWorkflowSentence,
			paths: []string{
				"assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md",
				"../../docs/agents/spec-routing.md",
			},
		},
		{
			id:       "clause.autonomous.loop-01-qa-once",
			sentence: autonomousWorkSentence,
			paths: []string{
				"assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md",
				"../../docs/agents/autonomous-work.md",
			},
		},
	} {
		t.Run(want.id, func(t *testing.T) {
			normalizedSentence := strings.Join(strings.Fields(want.sentence), " ")
			found := false
			for _, clause := range embeddedBaselineClauses(t) {
				if clause.ID != want.id {
					continue
				}
				found = true
				if !strings.Contains(strings.Join(strings.Fields(clause.Guidance), " "), normalizedSentence) {
					t.Errorf("embedded clause lacks %q", want.sentence)
				}
			}
			if !found {
				t.Fatalf("missing embedded clause %s", want.id)
			}

			for _, path := range want.paths {
				content, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				if !strings.Contains(strings.Join(strings.Fields(string(content)), " "), normalizedSentence) {
					t.Errorf("%s lacks %q", path, want.sentence)
				}
			}
		})
	}
}
