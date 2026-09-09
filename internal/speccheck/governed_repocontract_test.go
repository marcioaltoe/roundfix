//go:build repocontract

// Suite: governed-path authorization history contract.
// Invariant: every path bounded by an authorization record remains governed.
// Boundary IN: authorization records and the public governed-path predicate.
// Boundary OUT: changed-path audit integration and sanctioned regeneration.
package speccheck

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"roundfix/internal/suiteguardcontract"
)

const authorizationRecordsDir = "docs/workflow/authorizations"

const cleanupHistoricalGrantAncestor = "81a6afb48f4a3683d0e5fad52f3919cf1bdfbbf4"

func TestCleanupHistoricalGrantEvidence(t *testing.T) {
	repository, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("resolve repository root: %v", err)
	}
	if !governedCommitAvailable(repository, cleanupHistoricalGrantAncestor) {
		t.Skipf("historical grant ancestor %s is unavailable", cleanupHistoricalGrantAncestor)
	}

	listing := strings.TrimSpace(runGovernedGit(
		t,
		repository,
		"ls-tree",
		"-r",
		"--name-only",
		cleanupHistoricalGrantAncestor,
		authorizationRecordsDir,
	))
	records := strings.Split(listing, "\n")
	if len(records) != 42 {
		t.Fatalf("historical authorization records = %d, want 42", len(records))
	}

	historicalRoot := t.TempDir()
	recovered := make(map[string][]byte, len(records))
	boundedPathCount := 0
	for _, record := range records {
		if record == "" {
			t.Fatal("historical authorization listing contains an empty path")
		}
		content := []byte(runGovernedGit(
			t,
			repository,
			"show",
			cleanupHistoricalGrantAncestor+":"+record,
		))
		if len(content) == 0 {
			t.Fatalf("historical authorization record %q is empty", record)
		}
		boundedPathCount += len(parseMechanicalAuthorizationPaths(content))
		recovered[record] = append([]byte(nil), content...)
		writeGovernedAuthorizationRecord(t, historicalRoot, record, string(content))
	}
	if boundedPathCount == 0 {
		t.Fatal("historical authorization records expose no bounded paths")
	}

	findings, recordsExist, err := auditBoundedPathsAreGoverned(historicalRoot)
	if err != nil {
		t.Fatal(err)
	}
	if !recordsExist {
		t.Fatal("recovered historical authorization records were not read")
	}
	if len(findings) != 0 {
		t.Fatalf("historical bounded-path contract failed:\n%s", strings.Join(findings, "\n"))
	}

	const proofCostRecord = "docs/workflow/authorizations/2026-08-06-proof-cost.md"
	proofCost := recovered[proofCostRecord]
	if len(proofCost) == 0 {
		t.Fatalf("historical corpus has no %s bytes", proofCostRecord)
	}
	var proofCostOutputs []string
	for _, declaration := range suiteguardcontract.ParseSanctionedRegenerations(proofCost) {
		if declaration.Command == "make baseline-digests" {
			proofCostOutputs = declaration.Outputs
			break
		}
	}
	if len(proofCostOutputs) == 0 {
		t.Fatal("historical proof-cost grant has no baseline-digests output enumeration")
	}

	replayRoot := newGovernedGitRepo(t)
	writeGovernedResolverFixture(t, replayRoot)
	writeGovernedAuthorizationRecord(t, replayRoot, proofCostRecord, string(proofCost))
	commitGovernedFiles(
		t,
		replayRoot,
		"materialize recovered grant",
		proofCostRecord,
		"Makefile",
		"internal/baseline/derived/_ownership.yml",
		"internal/baseline/derived/frozen.txt",
	)

	acceptedPath := proofCostOutputs[0]
	writeGovernedAuthorizationRecord(t, replayRoot, acceptedPath, "regenerated\n")
	acceptedCommit := commitGovernedFiles(t, replayRoot, "replay accepted regeneration", acceptedPath)
	accepted := runGovernedMechanical(t, MechanicalRequest{
		RepoRoot:          replayRoot,
		AuthorizationPath: proofCostRecord,
		TaskCommits: []MechanicalTaskCommit{{
			TaskID: "task_accepted",
			SHA:    acceptedCommit,
		}},
	})
	assertNoGovernedMechanicalCode(t, accepted, CodeMechanicalAuthPaths)

	const refusedPath = ".golangci.yml"
	writeGovernedAuthorizationRecord(t, replayRoot, refusedPath, "linters: {}\n")
	refusedCommit := commitGovernedFiles(t, replayRoot, "replay refused change", refusedPath)
	refused := runGovernedMechanical(t, MechanicalRequest{
		RepoRoot:          replayRoot,
		AuthorizationPath: proofCostRecord,
		TaskCommits: []MechanicalTaskCommit{{
			TaskID: "task_refused",
			SHA:    refusedCommit,
		}},
	})
	assertGovernedPathEscapedGrant(t, refused, refusedPath, proofCostRecord)
}

func TestEveryBoundedPathIsGoverned(t *testing.T) {
	t.Run("repository records are governed", func(t *testing.T) {
		repoRoot, err := filepath.Abs(filepath.Join("..", ".."))
		if err != nil {
			t.Fatalf("resolve repository root: %v", err)
		}

		findings, recordsExist, err := auditBoundedPathsAreGoverned(repoRoot)
		if err != nil {
			t.Fatal(err)
		}
		if !recordsExist {
			t.Skipf("no authorization records under %s", authorizationRecordsDir)
		}
		if len(findings) != 0 {
			t.Fatalf("bounded-path contract failed:\n%s", strings.Join(findings, "\n"))
		}
	})

	t.Run("unmatched path names path and record", func(t *testing.T) {
		repoRoot := t.TempDir()
		const record = "docs/workflow/authorizations/synthetic.md"
		writeGovernedAuthorizationRecord(t, repoRoot, record, `---
granted: 2026-08-15
action: synthetic contract probe
paths:
  - README.md
consuming: synthetic-spec
---
`)

		findings, recordsExist, err := auditBoundedPathsAreGoverned(repoRoot)
		if err != nil {
			t.Fatal(err)
		}
		if !recordsExist {
			t.Fatal("authorization record was not read")
		}
		if len(findings) != 1 || !strings.Contains(findings[0], "README.md") || !strings.Contains(findings[0], record) {
			t.Fatalf("findings = %q, want one finding naming path %q and record %q", findings, "README.md", record)
		}
	})

	t.Run("no authorization records skips", func(t *testing.T) {
		findings, recordsExist, err := auditBoundedPathsAreGoverned(t.TempDir())
		if err != nil {
			t.Fatal(err)
		}
		if recordsExist || len(findings) != 0 {
			t.Fatalf("recordsExist = %t, findings = %q, want no records and no findings", recordsExist, findings)
		}
		t.Skipf("no authorization records under %s", authorizationRecordsDir)
	})
}

func auditBoundedPathsAreGoverned(repoRoot string) ([]string, bool, error) {
	directory := filepath.Join(repoRoot, filepath.FromSlash(authorizationRecordsDir))
	entries, err := os.ReadDir(directory)
	if errors.Is(err, os.ErrNotExist) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("read authorization records %q: %w", directory, err)
	}

	var findings []string
	recordsExist := false
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}
		recordsExist = true
		record := filepath.ToSlash(filepath.Join(authorizationRecordsDir, entry.Name()))
		content, err := os.ReadFile(filepath.Join(directory, entry.Name()))
		if err != nil {
			return nil, true, fmt.Errorf("read authorization record %q: %w", record, err)
		}

		bounded := parseMechanicalAuthorizationPaths(content)
		paths := make([]string, 0, len(bounded))
		for path := range bounded {
			// Legacy prose can cite a regeneration command in the same bullet as
			// a bounded file. Commands are not repository-relative paths.
			if strings.ContainsAny(path, " \t") {
				continue
			}
			paths = append(paths, path)
		}
		sort.Strings(paths)
		for _, path := range paths {
			if !GovernedPath(path) {
				findings = append(findings, fmt.Sprintf("%s bounds %s, which is not governed", record, path))
			}
		}
	}
	return findings, recordsExist, nil
}

func writeGovernedAuthorizationRecord(t *testing.T, repoRoot, relativePath, content string) {
	t.Helper()

	path := filepath.Join(repoRoot, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func governedCommitAvailable(repoRoot, commit string) bool {
	return exec.Command("git", "-C", repoRoot, "cat-file", "-e", commit+"^{commit}").Run() == nil
}

func runGovernedGit(t *testing.T, repoRoot string, args ...string) string {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", repoRoot}, args...)...)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v: %s", strings.Join(args, " "), err, output)
	}
	return string(output)
}

func newGovernedGitRepo(t *testing.T) string {
	t.Helper()
	repoRoot := t.TempDir()
	runGovernedGit(t, repoRoot, "init", "--quiet")
	runGovernedGit(t, repoRoot, "config", "user.name", "Roundfix Test")
	runGovernedGit(t, repoRoot, "config", "user.email", "roundfix@example.invalid")
	runGovernedGit(t, repoRoot, "config", "commit.gpgsign", "false")
	writeGovernedAuthorizationRecord(t, repoRoot, ".keep", "fixture\n")
	commitGovernedFiles(t, repoRoot, "initial", ".keep")
	return repoRoot
}

func writeGovernedResolverFixture(t *testing.T, repoRoot string) {
	t.Helper()
	writeGovernedAuthorizationRecord(t, repoRoot, "Makefile", "DERIVED_DIGEST_PATHS := internal/baseline/derived\n")
	writeGovernedAuthorizationRecord(t, repoRoot, "internal/baseline/derived/_ownership.yml", "owner: frozen\nreason: fixture\n")
	writeGovernedAuthorizationRecord(t, repoRoot, "internal/baseline/derived/frozen.txt", "frozen\n")
}

func commitGovernedFiles(t *testing.T, repoRoot, message string, paths ...string) string {
	t.Helper()
	args := append([]string{"add", "--"}, paths...)
	runGovernedGit(t, repoRoot, args...)
	runGovernedGit(t, repoRoot, "commit", "--quiet", "-m", message)
	return strings.TrimSpace(runGovernedGit(t, repoRoot, "rev-parse", "HEAD"))
}

func runGovernedMechanical(t *testing.T, request MechanicalRequest) MechanicalResult {
	t.Helper()
	result, err := RunMechanicalStage(context.Background(), request)
	if err != nil {
		t.Fatalf("RunMechanicalStage() error = %v", err)
	}
	return result
}

func assertNoGovernedMechanicalCode(t *testing.T, result MechanicalResult, code string) {
	t.Helper()
	for _, finding := range result.Findings {
		if finding.Code == code {
			t.Fatalf("unexpected %s finding: %#v", code, finding)
		}
	}
}

func assertGovernedPathEscapedGrant(
	t *testing.T,
	result MechanicalResult,
	path string,
	grant string,
) {
	t.Helper()
	var matching []MechanicalFinding
	for _, finding := range result.Findings {
		if finding.Code == CodeMechanicalAuthPaths {
			matching = append(matching, finding)
		}
	}
	if len(matching) != 1 {
		t.Fatalf("%s findings = %#v, want one finding for %s", CodeMechanicalAuthPaths, matching, path)
	}
	if !strings.Contains(matching[0].Detail, path) || matching[0].File != grant {
		t.Fatalf("%s finding = %#v, want path %s escaping grant %s", CodeMechanicalAuthPaths, matching[0], path, grant)
	}
}
