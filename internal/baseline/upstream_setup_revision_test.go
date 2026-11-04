// Suite: published upstream Setup Snapshot revision and tree contracts
// Invariant: setup revisions follow the catalog lock and onioncry, ui-ux-pro-max and knowledge-workspace never return to the trees replaced by Specs 0320, 0333 and 0326.
// Boundary IN: repository catalog lock, embedded catalog, setup provenance and skill tree digests.
// Boundary OUT: network, host skill checkouts, asset synchronization and skill installation.
package baseline

import (
	"encoding/json"
	"path/filepath"
	"testing"
)

// upstreamSetupRevision returns the one full commit every catalog entry of
// skills-lock.json names, refusing branches, short or mixed refs and no entry.
func upstreamSetupRevision(t *testing.T) string {
	t.Helper()
	lock, err := loadSkillsLock(filepath.Join("..", "..", skillsLockPath))
	if err != nil {
		t.Fatalf("read catalog revision from skills-lock.json: %v", err)
	}
	var revision string
	for _, entry := range lock.skillEntries() {
		source, _ := entry.value.stringField("source")
		if source != "marcioaltoe/skills" {
			continue
		}
		ref, _ := entry.value.stringField("ref")
		if len(ref) != 40 || !immutableGitID.MatchString(ref) {
			t.Fatalf("catalog skill %s ref = %q; skills-lock.json must name a full 40-character lowercase commit", entry.name, ref)
		}
		if revision != "" && ref != revision {
			t.Fatalf("catalog skill %s ref = %q, want shared catalog revision %q in skills-lock.json", entry.name, ref, revision)
		}
		revision = ref
	}
	if revision == "" {
		t.Fatal("skills-lock.json has no marcioaltoe/skills catalog entry")
	}
	return revision
}

// replacedOnioncryTree is the tree Spec 0320 replaced.
const replacedOnioncryTree = "f918508da464fcc37b436a5a26fc243ffa159ead9e07b19e991451ed0ca53b6c"

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
				if revision := upstreamSetupRevision(t); skill.Source.Ref != revision {
					t.Errorf("setup %s onioncry revision = %q, want %q", id, skill.Source.Ref, revision)
				}
				if skill.TreeDigest == replacedOnioncryTree {
					t.Errorf("setup %s onioncry carries tree %q replaced by Spec 0320", id, skill.TreeDigest)
				}
				return
			}
			t.Fatalf("setup %s lacks onioncry", id)
		})
	}
}

// incompleteUiUxProMaxTree is the tree Spec 0333 replaced.
const incompleteUiUxProMaxTree = "b13db828c6a7fbebf8612185dfde2e4cf20414f4c8fdaa33cc50b3b59f427d24"

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
				if revision := upstreamSetupRevision(t); skill.Source.Ref != revision {
					t.Errorf("setup %s ui-ux-pro-max revision = %q, want %q", id, skill.Source.Ref, revision)
				}
				if skill.TreeDigest == incompleteUiUxProMaxTree {
					t.Errorf("setup %s ui-ux-pro-max carries tree %q replaced by Spec 0333", id, skill.TreeDigest)
				}
				return
			}
			t.Fatalf("setup %s lacks ui-ux-pro-max", id)
		})
	}
}

// inboxKnowledgeWorkspaceTree is the tree Spec 0326 replaced.
const inboxKnowledgeWorkspaceTree = "19cb0bd73c1a2429e902b5e5c0e951fdd38d396352e7b1dc007748fc11a7aa7c"

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
				if revision := upstreamSetupRevision(t); skill.Source.Ref != revision {
					t.Errorf("setup %s knowledge-workspace revision = %q, want %q", id, skill.Source.Ref, revision)
				}
				if skill.TreeDigest == inboxKnowledgeWorkspaceTree {
					t.Errorf("setup %s knowledge-workspace carries tree %q replaced by Spec 0326", id, skill.TreeDigest)
				}
				return
			}
			t.Fatalf("setup %s lacks knowledge-workspace", id)
		})
	}
}
