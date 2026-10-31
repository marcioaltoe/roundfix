// Suite: published Python Setup Snapshot
// Invariant: Python mirrors the curated upstream preset and shares every upstream skill revision with the other snapshots.
// Boundary IN: embedded setup provenance, ordered skill membership, and upstream revisions.
// Boundary OUT: profile composition, asset synchronization, and skill installation.

package baseline

import (
	"encoding/json"
	"reflect"
	"slices"
	"testing"
)

func TestThePythonSetupIsThePublishedUpstreamPreset(t *testing.T) {
	t.Parallel()

	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	entry, ok := catalog.Setup("python")
	if !ok {
		t.Fatal("embedded Python Setup Snapshot is missing")
	}
	var setup assetsSyncSnapshot
	if err := json.Unmarshal(entry.Data, &setup); err != nil {
		t.Fatal(err)
	}
	wantSource := assetsSyncSource{
		Type: "github", Repository: "marcioaltoe/skills",
		Ref: upstreamSetupRevision, Path: "setups/python.txt",
	}
	if !reflect.DeepEqual(setup.Source, wantSource) {
		t.Errorf("Python setup source = %+v, want %+v", setup.Source, wantSource)
	}
	wantNames := []string{
		"find-skills",
		"knowledge-workspace",
		"setup-context-driven",
		"brainstorming",
		"business-analyst",
		"context7-cli",
		"exa-web-search",
		"grill-with-docs",
		"grilling",
		"write-idea",
		"write-prd",
		"agentic-cli-design",
		"domain-modeling",
		"write-techspec",
		"write-tasks",
		"coding-guidelines",
		"implement-task",
		"no-workarounds",
		"testing-boss",
		"typesafe-ai",
		"evidence-gate",
		"qa-gate",
		"systematic-debugging",
		"roundfix",
		"archive-spec",
		"conventional-commits",
		"crafting-effective-readmes",
		"github-pr-workflow",
		"ai-memory-durable-pages",
		"ai-memory-handoff",
		"ai-memory-learning-maintenance",
		"ai-memory-messaging",
		"ai-memory-retrieval",
		"ai-memory-routing-install",
		"tech-writer",
		"uv",
		"ruff",
		"ty",
		"pytest",
		"typer",
		"fastapi",
		"python-services",
		"python-file-pipelines",
		"backend-docs",
		"external-api-adapters",
		"integration-contract-testing",
		"security-best-practices",
		"pdf",
		"pdftk-server",
		"docker",
		"developing-with-streamlit",
	}
	var gotNames []string
	for _, skill := range setup.Skills {
		gotNames = append(gotNames, skill.Name)
	}
	if !slices.Equal(gotNames, wantNames) {
		t.Errorf("Python setup skills = %v, want ordered preset %v", gotNames, wantNames)
	}
}

func TestEveryUpstreamSetupFollowsThePythonPresetRevision(t *testing.T) {
	t.Parallel()

	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"go", "rust", "typescript", "python"} {
		t.Run(id, func(t *testing.T) {
			t.Parallel()

			entry, ok := catalog.Setup(id)
			if !ok {
				t.Fatalf("embedded setup %q is missing", id)
			}
			var setup assetsSyncSnapshot
			if err := json.Unmarshal(entry.Data, &setup); err != nil {
				t.Fatal(err)
			}
			if setup.Source.Ref != upstreamSetupRevision {
				t.Errorf("setup %s revision = %q, want %q", id, setup.Source.Ref, upstreamSetupRevision)
			}
			var hasMemoryHandoff bool
			for _, skill := range setup.Skills {
				if skill.Source.Type == "github" && skill.Source.Ref != upstreamSetupRevision {
					t.Errorf("setup %s skill %s revision = %q, want %q", id, skill.Name, skill.Source.Ref, upstreamSetupRevision)
				}
				if id == "python" && skill.Name == "handoff" {
					t.Error("Python setup still contains the retired handoff skill")
				}
				hasMemoryHandoff = hasMemoryHandoff || skill.Name == "ai-memory-handoff"
			}
			if id == "python" && !hasMemoryHandoff {
				t.Error("Python setup lacks ai-memory-handoff")
			}
		})
	}
}
