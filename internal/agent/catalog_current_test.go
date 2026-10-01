// Suite: current picker catalog.
// Invariant: current models lead each catalog and retired models are never offered.
// Boundary IN: static ModelCatalog data; OUT: live adapter availability.
package agent

import "testing"

func TestModelCatalogOpensWithTheCurrentModels(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		runtime string
		want    []string
	}{
		{"codex", []string{"gpt-6.1-sol", "gpt-6-astra", "gpt-6-sol", "gpt-6-luna", "gpt-5.6-sol", "gpt-5.6-terra", "gpt-5.6-luna"}},
		{"claude", []string{"opus", "sonnet", "claude-fable-5-1", "haiku", "default"}},
	} {
		t.Run(tt.runtime, func(t *testing.T) {
			// Ignore transitional entries so task_03 can remove them independently.
			var current []ModelChoice
			for _, choice := range ModelCatalog(tt.runtime) {
				if choice.Value != "gpt-5.5" && choice.Value != "claude-fable-5" {
					current = append(current, choice)
				}
			}
			if len(current) < len(tt.want) {
				t.Fatalf("catalog = %#v, want prefix %v", current, tt.want)
			}
			for i, value := range tt.want {
				if current[i].Value != value || current[i].Label != value {
					t.Errorf("choice %d = %#v, want label and value %q", i, current[i], value)
				}
			}
		})
	}
}

func TestModelCatalogOffersNoRetiredModel(t *testing.T) {
	t.Parallel()
	for _, runtime := range []string{"codex", "claude"} {
		for _, retired := range []string{"gpt-5.4", "gpt-5.4-mini", "gpt-5.3-codex-spark"} {
			t.Run(runtime+"/"+retired, func(t *testing.T) {
				for _, choice := range ModelCatalog(runtime) {
					if choice.Value == retired || choice.Label == retired {
						t.Fatalf("retired model offered: %#v", choice)
					}
				}
			})
		}
	}
}
