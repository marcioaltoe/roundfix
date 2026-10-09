package baseline

import (
	"os"
	"strings"
	"testing"
)

func TestTheGuidesExemptANetworkDeniedOutsideEvidenceRow(t *testing.T) {
	t.Parallel()
	const outsideEvidenceSentence = "An outside-evidence row blocked only because the Run sandbox denied network access, recorded as `blocked (environment: network denied: <host>)`, is the exception: the report records that the source was not reached, the row never decides a qualifying `partial`, and whoever needs that proof declares it under Unreachable Acceptance."
	const deliveryOrderSentence = "Neither does an outside-evidence row blocked only because the Run sandbox denied network access, recorded as `blocked (environment: network denied: <host>)`."

	for _, want := range []struct {
		id       string
		sentence string
		paths    []string
	}{
		{
			id:       "clause.spec.project-constraints-06-outside-evidence",
			sentence: outsideEvidenceSentence,
			paths: []string{
				"assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md",
				"../../docs/agents/spec-routing.md",
			},
		},
		{
			id:       "clause.autonomous.loop-01-qa-once",
			sentence: deliveryOrderSentence,
			paths: []string{
				"assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md",
				"../../docs/agents/autonomous-work.md",
			},
		},
	} {
		t.Run(want.id, func(t *testing.T) {
			found := false
			for _, clause := range embeddedBaselineClauses(t) {
				if clause.ID != want.id {
					continue
				}
				found = true
				if !strings.Contains(clause.Guidance, want.sentence) {
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
				if !strings.Contains(strings.Join(strings.Fields(string(content)), " "), want.sentence) {
					t.Errorf("%s lacks %q", path, want.sentence)
				}
			}
		})
	}
}
