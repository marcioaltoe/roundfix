// Suite: published upstream Setup Snapshot revision and tree contracts
// Invariant: both TypeScript setups pin the onioncry and ui-ux-pro-max trees and every setup pins the knowledge-workspace tree published at the shared upstream revision.
// Boundary IN: embedded catalog, setup provenance and skill tree digests.
// Boundary OUT: network, host skill checkouts, asset synchronization and skill installation.
package baseline

import (
	"encoding/json"
	"testing"
)

// upstreamSetupRevision is the upstream skills commit every upstream Setup
// Snapshot and every skill it pins follow; it moves with each sync.
const upstreamSetupRevision = "86df8608726979064a8fe46003dfaf403ffe0087"

// onioncryTreeDigest is the onioncry tree that commit publishes.
const onioncryTreeDigest = "66b054d73b38954303e64c4dc2c57916aab971dc8b320f76b930e5c3ef99f3d0"

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

// uiUxProMaxTreeDigest is the completed ui-ux-pro-max tree that commit publishes.
const uiUxProMaxTreeDigest = "34f2a8a3be83f0fff1e9af200c09ad33821d0380298ef2f4952da75652d80807"

func TestTheTypeScriptSetupsCarryTheCompletedUiUxProMaxTree(t *testing.T) {
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
				if skill.Name != "ui-ux-pro-max" {
					continue
				}
				if skill.Source.Ref != upstreamSetupRevision {
					t.Errorf("setup %s ui-ux-pro-max revision = %q, want %q", id, skill.Source.Ref, upstreamSetupRevision)
				}
				if skill.TreeDigest != uiUxProMaxTreeDigest {
					t.Errorf("setup %s ui-ux-pro-max tree = %q, want %q", id, skill.TreeDigest, uiUxProMaxTreeDigest)
				}
				return
			}
			t.Fatalf("setup %s lacks ui-ux-pro-max", id)
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
