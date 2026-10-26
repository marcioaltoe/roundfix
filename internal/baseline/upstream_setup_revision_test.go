// Suite: published upstream Setup Snapshot revision and tree contracts
// Invariant: both TypeScript setups pin the onioncry tree published at the shared upstream revision.
// Boundary IN: embedded catalog, setup provenance and skill tree digests.
// Boundary OUT: network, host skill checkouts, asset synchronization and skill installation.
package baseline

import (
	"encoding/json"
	"testing"
)

// upstreamSetupRevision is the upstream skills commit every upstream Setup
// Snapshot and every skill it pins follow; it moves with each sync.
const upstreamSetupRevision = "8d71fade68c54693bedbd917371f0514311af157"

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
