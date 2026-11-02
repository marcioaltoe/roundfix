// Suite: Delivery Queue publication content.
// Invariant: publication never requests CodeRabbit review under any pre-PR review policy.
// Boundary IN: publication title and body generated for a delivery candidate.
// Boundary OUT: pull request creation and provider-specific review execution.
package cli

import (
	"path/filepath"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/preflight"
	"roundfix/internal/spec"
)

func TestDeliveryPublicationRequestsNoCodeRabbitReview(t *testing.T) {
	t.Parallel()
	for _, provider := range []string{"codex", "claude", "coderabbit", "none"} {
		provider := provider
		t.Run(provider, func(t *testing.T) {
			_, checkout := newDeliveryBranchRepository(t)
			const slug = "0179-review-scope"
			record, err := spec.RenderArchiveRecord(spec.ArchiveRecord{Spec: slug, Title: "Review scope", Disposition: spec.ArchiveNoQA})
			if err != nil {
				t.Fatal(err)
			}
			mustMkdir(t, filepath.Join(checkout, "docs/history/specs"))
			mustWrite(t, filepath.Join(checkout, "docs/history/specs", slug+".md"), string(record))
			workflow := &commandDeliveryWorkflow{
				loaded: roundconfig.Loaded{
					Config: roundconfig.Config{
						Specs:       roundconfig.Specs{Root: "docs/specs"},
						PrePRReview: roundconfig.PrePRReview{Provider: provider, Source: "project"},
					},
				},
				git: preflight.ExecGitRunner{},
			}

			publication, err := workflow.Publication(t.Context(), checkout, slug, "main")
			if err != nil {
				t.Fatalf("plan publication under %s policy: %v", provider, err)
			}
			if publication.Title != "feat: deliver "+slug {
				t.Fatalf("title without Delivery Type=%q", publication.Title)
			}
			content := strings.ToLower(publication.Title + "\n" + publication.Body)
			for _, forbidden := range []string{"@coderabbitai", "coderabbit:review"} {
				if strings.Contains(content, forbidden) {
					t.Fatalf("publication under %s policy contains CodeRabbit request %q: title=%q body=%q", provider, forbidden, publication.Title, publication.Body)
				}
			}
		})
	}
}
