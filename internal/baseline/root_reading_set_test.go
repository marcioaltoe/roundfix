// Suite: root guide reading set
// Invariant: the root keeps domain rules mandatory and scopes the two costly guides.
// Boundary IN: the embedded Standard TypeScript golden root.
// Boundary OUT: Secondbrain access and repository Verification.

package baseline

import (
	"strings"
	"testing"
)

func TestTheRootNamesWhenToReadTheDocsLayoutGuide(t *testing.T) {
	t.Parallel()
	assertRootReadingSentence(t, "- Domain rules are mandatory: `docs/agents/domain.md`. Read `docs/agents/docs-layout.md`, whose rules then apply, before creating, changing, moving, or retiring `CONTEXT.md` or a document under `docs/`, and before a test or build step reads one.")
}

func TestTheRootNamesWhenToReadTheSecondbrainGuide(t *testing.T) {
	t.Parallel()
	assertRootReadingSentence(t, "- Read `docs/agents/secondbrain.md` before consulting or writing the Secondbrain or authoring an Idea, PRD, or TechSpec; a Run session, which cannot reach it, skips it.")
}

func assertRootReadingSentence(t *testing.T, sentence string) {
	t.Helper()
	catalog := mustEmbeddedCatalog(t)
	asset, ok := catalog.Asset("formatter-fixtures/standard-typescript-monorepo/golden/AGENTS.md")
	if !ok {
		t.Fatal("missing AGENTS.md golden root")
	}
	text := string(asset.Data)
	if count := strings.Count(text, sentence); count != 1 {
		t.Errorf("root sentence occurs %d times, want exactly once: %q", count, sentence)
	}
	for _, old := range []string{"Domain and documentation rules are mandatory", "Optional cross-project knowledge follows"} {
		if strings.Contains(text, old) {
			t.Errorf("old root sentence remains: %q", old)
		}
	}
}
