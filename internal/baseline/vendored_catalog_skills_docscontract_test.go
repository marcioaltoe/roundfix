// Suite: vendored catalog skill trees and Setup Snapshot provenance.
// Invariant: every vendored catalog skill matches its snapshot tree and every upstream snapshot follows the catalog lock revision.
// Boundary IN: repository skills lock, live vendored skill folders and embedded Setup Snapshots.
// Boundary OUT: network, host skill checkouts, asset synchronization and skill installation.
// These tests read the live documentation trees and run under the docscontract tag per ADR-0326.
//go:build docscontract

//verify:always

package baseline

import (
	"encoding/json"
	"path"
	"path/filepath"
	"slices"
	"sort"
	"testing"

	roundskills "roundfix/skills"
)

func TestEveryVendoredCatalogSkillIsTheTreeItsCatalogCommitPublishes(t *testing.T) {
	t.Parallel()
	revision := upstreamSetupRevision(t)
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	trees := make(map[string]string)
	for _, id := range catalog.SetupIDs() {
		entry, ok := catalog.Setup(id)
		if !ok {
			t.Fatalf("embedded setup %s is missing", id)
		}
		var snapshot assetsSyncSnapshot
		if err := json.Unmarshal(entry.Data, &snapshot); err != nil {
			t.Fatalf("decode setup %s: %v", id, err)
		}
		if snapshot.Source.Type == "github" && snapshot.Source.Ref != revision {
			t.Errorf("setup %s revision = %q, want catalog revision %q; run make skills-third-party-bump", id, snapshot.Source.Ref, revision)
		}
		for _, skill := range snapshot.Skills {
			if skill.Source.Type != "github" {
				continue
			}
			if skill.Source.Ref != revision {
				t.Errorf("setup %s skill %s revision = %q, want catalog revision %q; run make skills-third-party-bump", id, skill.Name, skill.Source.Ref, revision)
			}
			if skill.Source.Repository != "marcioaltoe/skills" || skill.TreeDigest == "" {
				continue
			}
			if previous, exists := trees[skill.Name]; exists && previous != skill.TreeDigest {
				t.Errorf("setup %s skill %s tree = %q, differs from snapshot tree %q; run make skills-third-party-bump", id, skill.Name, skill.TreeDigest, previous)
			}
			trees[skill.Name] = skill.TreeDigest
		}
	}

	repoRoot := filepath.Join("..", "..")
	lock, err := loadSkillsLock(filepath.Join(repoRoot, skillsLockPath))
	if err != nil {
		t.Fatal(err)
	}
	entries := lock.skillEntries()
	sort.Slice(entries, func(i, j int) bool { return entries[i].name < entries[j].name })
	owned := roundskills.Names()
	for _, entry := range entries {
		t.Run(entry.name, func(t *testing.T) {
			if slices.Contains(owned, entry.name) {
				t.Fatalf("Roundfix-owned skill %s must stay out of skills-lock.json", entry.name)
			}
			digest, ok := trees[entry.name]
			if !ok {
				t.Fatalf("catalog skill %s has no embedded Setup Snapshot treeDigest; run make skills-third-party-bump", entry.name)
			}
			files, exists, err := inspectRestoreTarget(repoRoot, path.Join(".agents/skills", entry.name))
			if err != nil {
				t.Fatalf("inspect catalog skill %s: %v; run make skills-third-party-bump", entry.name, err)
			}
			if !exists {
				t.Fatalf("catalog skill %s is missing; run make skills-third-party-bump", entry.name)
			}
			if got := portableRestoreDigest(files); got != digest {
				t.Errorf("catalog skill %s tree = %q, want Setup Snapshot tree %q; run make skills-third-party-bump", entry.name, got, digest)
			}
		})
	}
}
