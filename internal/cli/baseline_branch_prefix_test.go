package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestBaselineUpdateWithoutABranchPrefixAsksNothing(t *testing.T) {
	t.Parallel()
	repo := newBaselineUpdateRepository(t)
	manifest := ReadBaselineSetupManifest(t, repo)
	delete(manifest.Decisions, "branch.prefix")
	WriteBaselineSetupManifest(t, repo, manifest)
	result, stdout, stderr, code := runBaselineUpdateTestCommand(t, context.Background(), "baseline", "update", "--repo", repo, "--no-skills", "--format=json")
	if code != exitUnverified || result.State != "plan_ready" || result.Category == "decision" || len(result.NewDecisions) != 0 || stderr != "" {
		t.Fatalf("plan exit=%d result=%+v stdout=%s stderr=%s", code, result, stdout, stderr)
	}
	result, stdout, stderr, code = runBaselineUpdateTestCommand(t, context.Background(), "baseline", "update", "--repo", repo, "--no-skills", "--yes", "--format=json")
	if code != exitOK || result.State != "verified" || result.Category == "decision" || len(result.NewDecisions) != 0 || stderr != "" {
		t.Fatalf("apply exit=%d result=%+v stdout=%s stderr=%s", code, result, stdout, stderr)
	}
	guide, err := os.ReadFile(filepath.Join(repo, "docs/agents/agent-instructions.md"))
	if err != nil {
		t.Fatal(err)
	}
	want := "No branch prefix is recorded. Name new work branches `<type>/<description>`,\nwhere `<type>` is the work's Conventional Commit type, as the branch rule\nbelow states. Tool-owned Run and Task branches follow their tool's documented\nnamespace."
	if !bytes.Contains(guide, []byte(want+"\n\n")) || bytes.Contains(guide, []byte("The branch-prefix pattern is")) {
		t.Fatalf("guide = %s", guide)
	}
	if _, found := ReadBaselineSetupManifest(t, repo).Decisions["branch.prefix"]; found {
		t.Fatal("update recorded an unanswered prefix")
	}
}

func TestBaselineUpdateKeepsARecordedBranchPrefix(t *testing.T) {
	t.Parallel()
	repo := newBaselineUpdateRepository(t)
	path := filepath.Join(repo, "docs/agents/agent-instructions.md")
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	manifest := ReadBaselineSetupManifest(t, repo)
	manifest.CatalogDigest = "sha256:" + string(bytes.Repeat([]byte("0"), 64))
	WriteBaselineSetupManifest(t, repo, manifest)
	result, stdout, stderr, code := runBaselineUpdateTestCommand(t, context.Background(), "baseline", "update", "--repo", repo, "--no-skills", "--yes", "--format=json")
	if code != exitOK || result.State != "verified" || len(result.NewDecisions) != 0 || stderr != "" {
		t.Fatalf("update exit=%d result=%+v stdout=%s stderr=%s", code, result, stdout, stderr)
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Fatal("update changed recorded-prefix guide bytes")
	}
	if got := ReadBaselineSetupManifest(t, repo).Decisions["branch.prefix"].Value; got != "ma/" {
		t.Fatalf("recorded prefix = %v", got)
	}
}
