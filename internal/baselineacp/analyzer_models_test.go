package baselineacp

import (
	"roundfix/internal/agent"
	"testing"
)

func TestAnalysisModelsAreInTheCodexModelCatalog(t *testing.T) {
	for _, model := range []string{PreferredModel, FallbackModel} {
		t.Run(model, func(t *testing.T) {
			for _, choice := range agent.ModelCatalog("codex") {
				if choice.Value == model {
					return
				}
			}
			t.Fatalf("analysis model %q is absent from Codex catalog", model)
		})
	}
}
