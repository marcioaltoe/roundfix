package agent

import "testing"

func TestModelCatalogOffersNoReplacedModel(t *testing.T) {
	for _, runtime := range []string{"codex", "claude"} {
		for _, model := range []string{"gpt-5.5", "claude-fable-5"} {
			t.Run(runtime+"/"+model, func(t *testing.T) {
				for _, choice := range ModelCatalog(runtime) {
					if choice.Value == model {
						t.Fatalf("%s offers replaced model %q", runtime, model)
					}
				}
			})
		}
	}
}
