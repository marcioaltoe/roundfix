package spec

import "testing"

func TestADeferredBacklogEntryWithAReasonIsRetired(t *testing.T) {
	content := []byte("---\nstatus: deferred\nreason: no longer needed\n---\n\n# Backlog Entry\n")
	want := Retirement{Retired: true, Reason: "deferred"}
	if got := ClassifyBacklogEntry(content); got != want {
		t.Fatalf("retirement = %+v, want %+v", got, want)
	}
}

func TestADeferredBacklogEntryWithoutAReasonStaysLive(t *testing.T) {
	if got := ClassifyBacklogEntry([]byte(retirementDocument("deferred"))); got != (Retirement{}) {
		t.Fatalf("retirement = %+v, want live", got)
	}
}

func TestAConsumingSpecDoesNotRetireADeferredEntryWithoutAReason(t *testing.T) {
	content := []byte("---\nstatus: deferred\nspec: 0193-baseline-guides-that-describe-the-product-as-it-is\n---\n\n# Backlog Entry\n")
	if got := ClassifyBacklogEntry(content); got != (Retirement{}) {
		t.Fatalf("retirement = %+v, want live", got)
	}
}

func TestADeferredBacklogEntryWithAnEmptyReasonStaysLive(t *testing.T) {
	for _, reason := range []string{"null", "'  '", "'null'", "''"} {
		t.Run(reason, func(t *testing.T) {
			content := []byte("---\nstatus: deferred\nspec: consuming-spec\nreason: " + reason + "\n---\n")
			if got := ClassifyBacklogEntry(content); got != (Retirement{}) {
				t.Fatalf("retirement = %+v, want live", got)
			}
		})
	}
}
