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
	"time"

	"gopkg.in/yaml.v3"
	"roundfix/internal/suiteguardcontract"
)

const (
	legacyAuthorizationRecordsDir = "docs/workflow/authorizations"
	activeSpecAuthorizationRoot   = "docs/specs"
	archivedSpecAuthorizationRoot = "docs/history/specs"
)

const cleanupHistoricalGrantAncestor = "81a6afb48f4a3683d0e5fad52f3919cf1bdfbbf4"

func TestGovernedSetCoversOwnedShippedTemplates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path string
	}{
		{
			name: "shipped PRD template is governed",
			path: "skills/write-prd/references/prd-template.md",
		},
		{
			name: "shipped TechSpec template is governed",
			path: "skills/write-techspec/references/techspec-template.md",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if !GovernedPath(tt.path) {
				t.Fatalf("GovernedPath(%q) = false, want true", tt.path)
			}
			if !governedPathHasClause(tt.path, historicalToolingAuthorityClause) {
				t.Fatalf("GovernedPath(%q) does not carry clause %q", tt.path, historicalToolingAuthorityClause)
			}
		})
	}
}

func TestGovernedSetOnlyGrows(t *testing.T) {
	t.Parallel()

	previouslyGoverned := []string{
		".golangci.yml",
		".prettierrc",
		"tsconfig.json",
		"vitest.config.ts",
		".dependency-cruiser.js",
		"Makefile",
		"go.mod",
		"scripts/generate.sh",
		".gitignore",
		".codex-plugin/plugin.json",
		".tool-versions",
	}
	for _, governed := range previouslyGoverned {
		governed := governed
		t.Run("keeps "+governed+" governed", func(t *testing.T) {
			t.Parallel()

			if !GovernedPath(governed) {
				t.Fatalf("GovernedPath(%q) = false, want true", governed)
			}
		})
	}

	newlyGoverned := []string{
		"internal/speccheck/governed.go",
		"internal/speccheck/governed_repocontract_test.go",
		"internal/speccheck/mechanical_test.go",
		"internal/suiteguardcontract/regeneration.go",
		"internal/suiteguardcontract/regeneration_test.go",
		"skills/write-prd/references/prd-template.md",
		"skills/write-techspec/references/techspec-template.md",
	}
	for _, governed := range newlyGoverned {
		governed := governed
		t.Run("adds "+governed+" under the historical clause", func(t *testing.T) {
			t.Parallel()

			if !GovernedPath(governed) {
				t.Fatalf("GovernedPath(%q) = false, want true", governed)
			}
			if !governedPathHasClause(governed, historicalToolingAuthorityClause) {
				t.Fatalf("GovernedPath(%q) does not carry clause %q", governed, historicalToolingAuthorityClause)
			}
		})
	}

	ordinaryPaths := []string{
		"internal/app/metadata.go",
		"docs/specs/example/_prd.md",
	}
	for _, ordinary := range ordinaryPaths {
		ordinary := ordinary
		t.Run("keeps "+ordinary+" ungoverned", func(t *testing.T) {
			t.Parallel()

			if GovernedPath(ordinary) {
				t.Fatalf("GovernedPath(%q) = true, want false", ordinary)
			}
		})
	}
}

func governedPathHasClause(path, clause string) bool {
	clean := cleanMechanicalPath(path)
	for _, entry := range governedPathSet {
		if entry.clause == clause && entry.matches(clean) {
			return true
		}
	}
	return false
}

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
		legacyAuthorizationRecordsDir,
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

	findings, auditedRecords, err := auditBoundedPathsAreGoverned(historicalRoot)
	if err != nil {
		t.Fatal(err)
	}
	if len(auditedRecords) == 0 {
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
	target := commitGovernedFiles(
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
		RepoRoot:               replayRoot,
		AuthorizationPath:      proofCostRecord,
		DeliveryTargetRevision: target,
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
		RepoRoot:               replayRoot,
		AuthorizationPath:      proofCostRecord,
		DeliveryTargetRevision: target,
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

		findings, auditedRecords, err := auditBoundedPathsAreGoverned(repoRoot)
		if err != nil {
			t.Fatal(err)
		}
		for _, requiredSlug := range []string{
			"0119-spec-contained-authorization",
			"0130-documentation-cleanup-compatibility",
		} {
			if _, found := governedSpecAuthorizationRecord(auditedRecords, requiredSlug); !found {
				t.Fatalf(
					"audited authorization records = %q, want Spec %q included under %s or %s",
					auditedRecords,
					requiredSlug,
					activeSpecAuthorizationRoot,
					archivedSpecAuthorizationRoot,
				)
			}
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

		findings, auditedRecords, err := auditBoundedPathsAreGoverned(repoRoot)
		if err != nil {
			t.Fatal(err)
		}
		if len(auditedRecords) != 1 || auditedRecords[0] != record {
			t.Fatal("authorization record was not read")
		}
		if len(findings) != 1 || !strings.Contains(findings[0], "README.md") || !strings.Contains(findings[0], record) {
			t.Fatalf("findings = %q, want one finding naming path %q and record %q", findings, "README.md", record)
		}
	})

	t.Run("no authorization records skips", func(t *testing.T) {
		repoRoot := t.TempDir()
		writeGovernedAuthorizationRecord(t, repoRoot, "docs/specs/proposed/_authorization.md", `---
status: proposed
granted: null
action: synthetic proposed contract probe
consuming: proposed
paths:
  - README.md
---
`)

		findings, auditedRecords, err := auditBoundedPathsAreGoverned(repoRoot)
		if err != nil {
			t.Fatal(err)
		}
		if len(auditedRecords) != 0 || len(findings) != 0 {
			t.Fatalf("auditedRecords = %q, findings = %q, want no operative records and no findings", auditedRecords, findings)
		}
		t.Skip("no operative authorization records")
	})
}

func governedSpecAuthorizationRecord(records []string, slug string) (string, bool) {
	for _, root := range []string{activeSpecAuthorizationRoot, archivedSpecAuthorizationRoot} {
		candidate := filepath.ToSlash(filepath.Join(root, slug, "_authorization.md"))
		index := sort.SearchStrings(records, candidate)
		if index < len(records) && records[index] == candidate {
			return candidate, true
		}
	}
	return "", false
}

type governedAuthorizationRecord struct {
	path    string
	content []byte
}

type governedAuthorizationRoot struct {
	path          string
	specContained bool
}

func auditBoundedPathsAreGoverned(repoRoot string) ([]string, []string, error) {
	records, err := discoverGovernedAuthorizationRecords(repoRoot)
	if err != nil {
		return nil, nil, err
	}

	var findings []string
	auditedRecords := make([]string, 0, len(records))
	for _, record := range records {
		auditedRecords = append(auditedRecords, record.path)
		bounded := parseMechanicalAuthorizationPaths(record.content)
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
				findings = append(findings, fmt.Sprintf("%s bounds %s, which is not governed", record.path, path))
			}
		}
	}
	return findings, auditedRecords, nil
}

func discoverGovernedAuthorizationRecords(repoRoot string) ([]governedAuthorizationRecord, error) {
	roots := []governedAuthorizationRoot{
		{path: legacyAuthorizationRecordsDir},
		{path: activeSpecAuthorizationRoot, specContained: true},
		{path: archivedSpecAuthorizationRoot, specContained: true},
	}

	var records []governedAuthorizationRecord
	for _, root := range roots {
		directory := filepath.Join(repoRoot, filepath.FromSlash(root.path))
		info, err := os.Lstat(directory)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("inspect authorization root %q: %w", directory, err)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			continue
		}
		if !info.IsDir() {
			return nil, fmt.Errorf("authorization root %q is not a directory", directory)
		}

		err = filepath.WalkDir(directory, func(filePath string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return fmt.Errorf("inspect authorization record %q: %w", filePath, walkErr)
			}
			if entry.IsDir() || entry.Type()&os.ModeSymlink != 0 {
				return nil
			}
			if root.specContained {
				if entry.Name() != "_authorization.md" {
					return nil
				}
			} else if filepath.Ext(entry.Name()) != ".md" {
				return nil
			}
			entryInfo, err := entry.Info()
			if err != nil {
				return fmt.Errorf("stat authorization record %q: %w", filePath, err)
			}
			if !entryInfo.Mode().IsRegular() {
				return nil
			}

			content, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("read authorization record %q: %w", filePath, err)
			}
			if root.specContained {
				relative, err := filepath.Rel(directory, filePath)
				if err != nil {
					return fmt.Errorf("make authorization record %q relative to %q: %w", filePath, directory, err)
				}
				parts := strings.Split(filepath.ToSlash(relative), "/")
				if len(parts) != 2 || !operativeGovernedSpecAuthorization(content, parts[0]) {
					return nil
				}
			}

			relative, err := filepath.Rel(repoRoot, filePath)
			if err != nil {
				return fmt.Errorf("make authorization record %q relative to %q: %w", filePath, repoRoot, err)
			}
			records = append(records, governedAuthorizationRecord{
				path:    filepath.ToSlash(relative),
				content: content,
			})
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("walk authorization records %q: %w", directory, err)
		}
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].path < records[j].path
	})
	return records, nil
}

func operativeGovernedSpecAuthorization(content []byte, specSlug string) bool {
	text := strings.ReplaceAll(string(content), "\r\n", "\n")
	if !strings.HasPrefix(text, "---\n") {
		return false
	}
	frontmatter, _, found := strings.Cut(text[len("---\n"):], "\n---")
	if !found {
		return false
	}

	var grant struct {
		Status    string   `yaml:"status"`
		Granted   string   `yaml:"granted"`
		Action    string   `yaml:"action"`
		Consuming string   `yaml:"consuming"`
		Paths     []string `yaml:"paths"`
	}
	if err := yaml.Unmarshal([]byte(frontmatter), &grant); err != nil {
		return false
	}
	if strings.TrimSpace(grant.Status) != "approved" ||
		strings.TrimSpace(grant.Action) == "" ||
		strings.TrimSpace(grant.Consuming) != specSlug {
		return false
	}
	if _, err := time.Parse("2006-01-02", strings.TrimSpace(grant.Granted)); err != nil {
		return false
	}
	if len(grant.Paths) == 0 {
		return false
	}

	seen := make(map[string]struct{}, len(grant.Paths))
	for _, declared := range grant.Paths {
		clean := cleanMechanicalPath(declared)
		if clean == "" || clean != declared || strings.ContainsAny(clean, "*?") {
			return false
		}
		if _, duplicate := seen[clean]; duplicate {
			return false
		}
		seen[clean] = struct{}{}
	}
	return true
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
