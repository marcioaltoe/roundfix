// Suite: repository markdown corpus contracts
// Invariant: the active Spec corpus, and the repository guides it cites, stay
// consistent under the shipped checker; nothing under _archived is asserted.
// Boundary IN: real repository markdown under docs/ and the public
// speccheck.Check API.
// Boundary OUT: fixture-based detector unit tests, which stay in
// internal/speccheck.

// The docscontract tag keeps this invalidation domain out of go test ./...;
// make verify-docs runs it at the pull request boundary.
//go:build docscontract

package docscontract

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

// maxCorpusCheckOperationsPerSpec bounds the sweep to one Check per Spec.
const maxCorpusCheckOperationsPerSpec = 1

func TestCheckCorpusGolden(t *testing.T) {
	repoRoot := characterizationRepositoryRoot(t)
	activeRoot := filepath.Join(repoRoot, "docs", "specs")
	want := readCorpusGolden(t)

	// Archived Specs are historical and are never validated: only the active
	// corpus is swept, where a changed count is a consistency regression.
	active := sweepCorpus(t, activeRoot, repoRoot, nil)

	if !reflect.DeepEqual(active, want.Active) {
		actual, err := json.MarshalIndent(active, "", "  ")
		if err != nil {
			t.Fatalf("render active corpus counts: %v", err)
		}
		t.Errorf("active Spec corpus finding counts changed; inspect the consistency regression or intentional detector change, then deliberately update testdata/corpus-golden.json:\n%s", actual)
	}
}

func TestCheckCorpusBudget(t *testing.T) {
	// Operation counts are load-independent, so this check no longer needs a dedicated-run guard.
	repoRoot := characterizationRepositoryRoot(t)
	activeRoot := filepath.Join(repoRoot, "docs", "specs")

	var work corpusSweepWork
	started := time.Now()
	sweepCorpus(t, activeRoot, repoRoot, &work)
	elapsed := time.Since(started)
	t.Logf(
		"active Spec corpus sweep completed in %s; work: %d Check operations across %d Specs (budget: at most %d Check operation per Spec)",
		elapsed,
		work.checkOperations,
		work.specs,
		maxCorpusCheckOperationsPerSpec,
	)

	if work.specs == 0 {
		// Zero measured Specs means either a broken sweep or an empty active
		// set, and only the first is a defect. The repository reached the
		// second state on 2026-08-12, when the last active Spec archived.
		entries, err := os.ReadDir(activeRoot)
		if err != nil {
			t.Fatalf("read active Spec Root: %v", err)
		}
		for _, entry := range entries {
			if !entry.IsDir() || strings.HasPrefix(entry.Name(), "_") {
				continue
			}
			if _, err := os.Stat(filepath.Join(activeRoot, entry.Name(), "_prd.md")); err == nil {
				t.Fatalf("active Spec corpus sweep measured no Specs while %s carries one", entry.Name())
			}
		}
		t.Log("active Spec Root holds no Spec; the sweep is correctly empty")
		return
	}
	maxCheckOperations := work.specs * maxCorpusCheckOperationsPerSpec
	if work.checkOperations > maxCheckOperations {
		t.Errorf(
			"full Spec corpus sweep performed %d Check operations across %d Specs, want at most %d",
			work.checkOperations,
			work.specs,
			maxCheckOperations,
		)
	}
}

func TestCheckActiveCorpusHasNoErrors(t *testing.T) {
	repoRoot := characterizationRepositoryRoot(t)
	specsRoot := filepath.Join(repoRoot, "docs", "specs")
	active, err := spec.ListActive(specsRoot)
	if err != nil {
		t.Fatalf("list active Specs: %v", err)
	}

	for _, activeSpec := range active {
		result, err := speccheck.Check(specsRoot, repoRoot, activeSpec.Slug)
		if err != nil {
			t.Errorf("%s: check active Spec: %v", activeSpec.Slug, err)
			continue
		}
		for _, finding := range result.Findings {
			switch finding.Severity {
			case speccheck.SeverityError:
				t.Errorf("%s: %s: %s", activeSpec.Slug, finding.Code, finding.Summary)
			case speccheck.SeverityGap:
				t.Logf("%s: %s gap remains visible: %s", activeSpec.Slug, finding.Code, finding.Summary)
			}
		}
	}
}

type corpusGolden struct {
	Schema string         `json:"schema"`
	Update string         `json:"update"`
	Active map[string]int `json:"active"`
}

type corpusSweepWork struct {
	specs           int
	checkOperations int
}

func characterizationRepositoryRoot(t *testing.T) string {
	t.Helper()

	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller could not locate the repository")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(testFile), "..", ".."))
}

func readCorpusGolden(t *testing.T) corpusGolden {
	t.Helper()

	path := filepath.Join("testdata", "corpus-golden.json")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read corpus golden %q: %v", path, err)
	}
	var golden corpusGolden
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&golden); err != nil {
		t.Fatalf("parse corpus golden %q: %v", path, err)
	}
	if golden.Schema != "roundfix-speccheck-corpus/v2" || strings.TrimSpace(golden.Update) == "" {
		t.Fatalf("corpus golden %q must declare its schema and update path", path)
	}
	return golden
}

func sweepCorpus(t *testing.T, specsRoot, repoRoot string, work *corpusSweepWork) map[string]int {
	t.Helper()

	counts := make(map[string]int, len(corpusFindingCodes))
	for _, code := range corpusFindingCodes {
		counts[code] = 0
	}
	entries, err := os.ReadDir(specsRoot)
	if err != nil {
		t.Fatalf("read Spec corpus %q: %v", specsRoot, err)
	}
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), "_") {
			continue
		}
		if _, err := os.Stat(filepath.Join(specsRoot, entry.Name(), "_prd.md")); errors.Is(err, os.ErrNotExist) {
			continue
		} else if err != nil {
			t.Fatalf("inspect Spec %q: %v", entry.Name(), err)
		}
		if work != nil {
			work.specs++
		}
		result, err := checkCorpusSpec(specsRoot, repoRoot, entry.Name(), work)
		if err != nil {
			t.Fatalf("Check(%q) in corpus %q: %v", entry.Name(), specsRoot, err)
		}
		for _, finding := range result.Findings {
			if _, known := counts[finding.Code]; !known {
				t.Fatalf("Check(%q) returned uncharacterized code %q", entry.Name(), finding.Code)
			}
			counts[finding.Code]++
		}
	}
	return counts
}

func TestCheckLoopOrderRepositoryAgrees(t *testing.T) {
	t.Parallel()

	repoRoot := characterizationRepositoryRoot(t)
	specsRoot := writeLoopOrderCarrier(t, t.TempDir())
	// SC-LOOP-ORDER-DIVERGENT reads the repository's own order statements, so
	// a dedicated carrier keeps the check independent from active and archived
	// repository Specs.
	result, err := speccheck.Check(specsRoot, repoRoot, loopOrderCarrierSlug)
	if err != nil {
		t.Fatalf("Check(repository) error = %v", err)
	}
	if findings := findingsWithCode(result, speccheck.CodeLoopOrderDivergent); len(findings) != 0 {
		t.Fatalf("%s findings = %#v, want corrected repository sources to agree", speccheck.CodeLoopOrderDivergent, findings)
	}
}

func writeLoopOrderFixture(t *testing.T) string {
	t.Helper()

	sourceRoot := characterizationRepositoryRoot(t)
	targetRoot := t.TempDir()
	for _, relative := range []string{
		loopOrderShippedClausePath,
		loopOrderRepositoryGuidePath,
		loopOrderBaselineModulePath,
	} {
		sourcePath := filepath.Join(sourceRoot, filepath.FromSlash(relative))
		content, err := os.ReadFile(sourcePath)
		if err != nil {
			t.Fatalf("read loop-order source %q: %v", sourcePath, err)
		}
		targetPath := filepath.Join(targetRoot, filepath.FromSlash(relative))
		if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
			t.Fatalf("create loop-order fixture directory: %v", err)
		}
		if err := os.WriteFile(targetPath, content, 0o644); err != nil {
			t.Fatalf("write loop-order fixture %q: %v", targetPath, err)
		}
	}

	writeLoopOrderCarrier(t, targetRoot)
	return targetRoot
}

func TestStageScopeDefaultSweepIsUnchanged(t *testing.T) {
	repoRoot := characterizationRepositoryRoot(t)
	activeRoot := filepath.Join(repoRoot, "docs", "specs")
	assertDefaultStageSweepUnchanged(t, activeRoot, repoRoot)
}

func assertDefaultStageSweepUnchanged(t *testing.T, specsRoot, repoRoot string) {
	t.Helper()

	entries, err := os.ReadDir(specsRoot)
	if err != nil {
		t.Fatalf("read Spec corpus %q: %v", specsRoot, err)
	}
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), "_") {
			continue
		}
		if _, err := os.Stat(filepath.Join(specsRoot, entry.Name(), "_prd.md")); errors.Is(err, os.ErrNotExist) {
			continue
		} else if err != nil {
			t.Fatalf("inspect Spec %q: %v", entry.Name(), err)
		}
		t.Run(filepath.Base(specsRoot)+"/"+entry.Name(), func(t *testing.T) {
			fullResult, err := speccheck.Check(specsRoot, repoRoot, entry.Name())
			if err != nil {
				t.Fatalf("Check(): %v", err)
			}
			defaultResult, err := speccheck.CheckStage(specsRoot, repoRoot, entry.Name(), speccheck.StageAll)
			if err != nil {
				t.Fatalf("CheckStage(StageAll): %v", err)
			}
			if !reflect.DeepEqual(defaultResult, fullResult) {
				t.Errorf("default stage result = %#v, want unchanged result %#v", defaultResult, fullResult)
			}
		})
	}
}

func checkCorpusSpec(specsRoot, repoRoot, slug string, work *corpusSweepWork) (speccheck.Result, error) {
	if work != nil {
		work.checkOperations++
	}
	return speccheck.Check(specsRoot, repoRoot, slug)
}

var corpusFindingCodes = []string{
	speccheck.CodeCitationUnsupported,
	speccheck.CodeConstraintMissing,
	speccheck.CodeConstraintUnreasoned,
	speccheck.CodeConstraintSource,
	speccheck.CodeToolingUnauthorized,
	speccheck.CodeToolingUnbounded,
	speccheck.CodeADRUnlisted,
	speccheck.CodeADRRelated,
	speccheck.CodeCoverageUnmapped,
	speccheck.CodeCoverageUntasked,
	speccheck.CodeLoopOrderDivergent,
	speccheck.CodeFindingLifecycle,
	speccheck.CodeRollupMember,
	speccheck.CodeArchiveLicense,
	speccheck.CodeReferenceUnresolved,
	speccheck.CodeVocabularyUndocumented,
	speccheck.CodeVerifyWorkIndependent,
	speccheck.CodeVerifyInvertedExit,
	speccheck.CodeVerifyNonHermetic,
	speccheck.CodeVerifyVacuousCommand,
	speccheck.CodeRequirementContradictory,
	speccheck.CodeRehearsalUndeclared,
}

const (
	loopOrderShippedClausePath   = "internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md"
	loopOrderRepositoryGuidePath = "docs/agents/autonomous-work.md"
	loopOrderBaselineModulePath  = "internal/baseline/assets/modules/autonomous-work.json"
	loopOrderCarrierSlug         = "loop-order"
)

func TestCheckLoopOrderDivergent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		path      string
		current   string
		divergent string
	}{
		{
			name:      "shipped clause",
			path:      loopOrderShippedClausePath,
			current:   "archive, open the Pull Request, watch until Clean, and merge",
			divergent: "merge, open the Pull Request, watch until Clean, and archive",
		},
		{
			name: "repository guide",
			path: loopOrderRepositoryGuidePath,
			// The guide is the rendered catalog clause after greenfield adoption,
			// so it carries the shipped wording without the paraphrase's wrap.
			current:   "archive, open the Pull Request, watch until Clean, and merge",
			divergent: "merge, open the Pull Request, watch until Clean, and archive",
		},
		{
			name:      "Baseline module asset",
			path:      loopOrderBaselineModulePath,
			current:   "archive, open the Pull Request, watch until Clean, and merge",
			divergent: "merge, open the Pull Request, watch until Clean, and archive",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repoRoot := writeLoopOrderFixture(t)
			path := filepath.Join(repoRoot, filepath.FromSlash(tt.path))
			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read %s fixture: %v", tt.name, err)
			}
			if !strings.Contains(string(content), tt.current) {
				t.Fatalf("%s fixture does not contain current order %q", tt.name, tt.current)
			}
			content = []byte(strings.Replace(string(content), tt.current, tt.divergent, 1))
			if err := os.WriteFile(path, content, 0o644); err != nil {
				t.Fatalf("write divergent %s fixture: %v", tt.name, err)
			}

			result, err := speccheck.Check(filepath.Join(repoRoot, "docs", "specs"), repoRoot, loopOrderCarrierSlug)
			if err != nil {
				t.Fatalf("Check(loop-order) error = %v", err)
			}
			finding := requireFinding(t, result, speccheck.CodeLoopOrderDivergent)
			if finding.Severity != speccheck.SeverityError {
				t.Fatalf("severity = %q, want %q", finding.Severity, speccheck.SeverityError)
			}
			for _, sourceLabel := range []string{"shipped clause", "repository guide", "Baseline module asset"} {
				if !strings.Contains(finding.Summary, sourceLabel) {
					t.Errorf("summary = %q, want source %q", finding.Summary, sourceLabel)
				}
			}
			if !strings.Contains(finding.Summary, strings.ReplaceAll(tt.divergent, "\n", " ")) {
				t.Errorf("summary = %q, want divergent order %q", finding.Summary, tt.divergent)
			}
			for _, sourcePath := range []string{
				loopOrderShippedClausePath,
				loopOrderRepositoryGuidePath,
				loopOrderBaselineModulePath,
			} {
				if !hasLocation(finding, sourcePath) {
					t.Errorf("locations = %#v, want source %q", finding.Where, sourcePath)
				}
			}
		})
	}
}

func writeLoopOrderCarrier(t *testing.T, root string) string {
	t.Helper()
	specsRoot := filepath.Join(root, "docs", "specs")
	specDir := filepath.Join(specsRoot, loopOrderCarrierSlug)
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatalf("create fixture Spec: %v", err)
	}
	const prd = "---\nspec: " + loopOrderCarrierSlug + "\nstatus: active\n---\n\n# Loop order\n"
	if err := os.WriteFile(filepath.Join(specDir, "_prd.md"), []byte(prd), 0o644); err != nil {
		t.Fatalf("write fixture PRD: %v", err)
	}
	return specsRoot
}

func findingsWithCode(result speccheck.Result, code string) []speccheck.Finding {
	var findings []speccheck.Finding
	for _, finding := range result.Findings {
		if finding.Code == code {
			findings = append(findings, finding)
		}
	}
	return findings
}

func requireFinding(t *testing.T, result speccheck.Result, code string) speccheck.Finding {
	t.Helper()

	for _, finding := range result.Findings {
		if finding.Code == code {
			return finding
		}
	}
	t.Fatalf("Findings = %#v, want code %s", result.Findings, code)
	return speccheck.Finding{}
}

func hasLocation(finding speccheck.Finding, path string) bool {
	for _, location := range finding.Where {
		if location.Path == path && location.Line > 0 {
			return true
		}
	}
	return false
}
