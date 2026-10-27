// Suite: published upstream Setup Snapshot revision and tree contracts
// Invariant: both TypeScript setups pin the onioncry tree and every setup pins the knowledge-workspace tree published at the shared upstream revision.
// Boundary IN: embedded catalog, setup provenance and skill tree digests.
// Boundary OUT: network, host skill checkouts, asset synchronization and skill installation.
package baseline

import (
	"encoding/json"
	"testing"
)

// upstreamSetupRevision is the upstream skills commit every upstream Setup
// Snapshot and every skill it pins follow; it moves with each sync.
const upstreamSetupRevision = "b56fc5ffcad1ab366fbc03839a6ce8c662401178"

// onioncryTreeDigest is the onioncry tree that commit publishes.
const onioncryTreeDigest = "be502b02544be90059c126e9b07bfeb6a4af115aada86455845b1407a5eab8b3"

func TestTheTypeScriptSetupsCarryTheRepublishedOnioncryTree(t *testing.T) {
	t.Parallel()

	for _, id := range []string{"typescript", "go-cli-typescript-bun"} {
		t.Run(id, func(t *testing.T) {
			t.Parallel()

			catalog, err := LoadEmbeddedCatalog()
			if err != nil {
				t.Fatal(err)
			}
			entry, ok := catalog.Setup(id)
			if !ok {
				t.Fatalf("embedded setup %q is missing", id)
			}
			var setup assetsSyncSnapshot
			if err := json.Unmarshal(entry.Data, &setup); err != nil {
				t.Fatal(err)
			}
			for _, skill := range setup.Skills {
				if skill.Name != "onioncry" {
					continue
				}
				if skill.Source.Ref != upstreamSetupRevision {
					t.Errorf("setup %s onioncry revision = %q, want %q", id, skill.Source.Ref, upstreamSetupRevision)
				}
				if skill.TreeDigest != onioncryTreeDigest {
					t.Errorf("setup %s onioncry tree = %q, want %q", id, skill.TreeDigest, onioncryTreeDigest)
				}
				return
			}
			t.Fatalf("setup %s lacks onioncry", id)
		})
	}
}

// knowledgeWorkspaceTreeDigest is the knowledge-workspace tree that commit
// publishes, the one that sends research through ai-memory.
const knowledgeWorkspaceTreeDigest = "490af37df4faba388e560380df0863b88881716d96e6185858918681182f2a96"

func TestEverySetupCarriesTheKnowledgeWorkspaceThatRetiredTheSecondbrainInbox(t *testing.T) {
	t.Parallel()

	for _, id := range []string{"go", "go-cli-typescript-bun", "python", "rust", "typescript"} {
		t.Run(id, func(t *testing.T) {
			t.Parallel()

			catalog, err := LoadEmbeddedCatalog()
			if err != nil {
				t.Fatal(err)
			}
			entry, ok := catalog.Setup(id)
			if !ok {
				t.Fatalf("embedded setup %q is missing", id)
			}
			var setup assetsSyncSnapshot
			if err := json.Unmarshal(entry.Data, &setup); err != nil {
				t.Fatal(err)
			}
			for _, skill := range setup.Skills {
				if skill.Name != "knowledge-workspace" {
					continue
				}
				if skill.Source.Ref != upstreamSetupRevision {
					t.Errorf("setup %s knowledge-workspace revision = %q, want %q", id, skill.Source.Ref, upstreamSetupRevision)
				}
				if skill.TreeDigest != knowledgeWorkspaceTreeDigest {
					t.Errorf("setup %s knowledge-workspace tree = %q, want %q", id, skill.TreeDigest, knowledgeWorkspaceTreeDigest)
				}
				return
			}
			t.Fatalf("setup %s lacks knowledge-workspace", id)
		})
	}
}
