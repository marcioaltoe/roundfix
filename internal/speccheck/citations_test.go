// Suite: Spec citation, coverage, reference, loop-order, and findings consistency
// Invariant: declared ADR, coverage, path, loop-order, and findings sources report only their written inconsistencies.
// Boundary IN: public speccheck API, fixture artifacts, accepted ADR corpus, and internal/spec loader
// Boundary OUT: CLI exit-code policy and characterization corpus assigned to later Spec Tasks
package speccheck_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/speccheck"
)

const (
	loopOrderShippedClausePath   = "internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md"
	loopOrderRepositoryGuidePath = "docs/agents/autonomous-work.md"
	loopOrderBaselineModulePath  = "internal/baseline/assets/modules/autonomous-work.json"
	loopOrderCarrierSlug         = "loop-order"
	findingsCarrierSlug          = "findings"
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
			name:      "repository guide",
			path:      loopOrderRepositoryGuidePath,
			current:   "archive, open the Pull\nRequest, watch until Clean, and merge",
			divergent: "merge, open the Pull\nRequest, watch until Clean, and archive",
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

func TestCheckFindingLifecycle(t *testing.T) {
	t.Parallel()

	t.Run("red missing lifecycle frontmatter", func(t *testing.T) {
		t.Parallel()

		repoRoot := writeFindingsCarrier(t)
		const findingPath = "docs/findings/2026-08-06-missing-lifecycle.md"
		writeFindingsArtifact(t, repoRoot, findingPath, "# Missing lifecycle\n")

		result := checkFindingsCarrier(t, repoRoot)
		finding := requireRenderedFinding(t, result, speccheck.CodeFindingLifecycle, findingPath, 1)
		if !strings.Contains(finding.Summary, "no lifecycle status") {
			t.Fatalf("summary = %q, want missing lifecycle status", finding.Summary)
		}
	})

	t.Run("red unknown lifecycle value", func(t *testing.T) {
		t.Parallel()

		repoRoot := writeFindingsCarrier(t)
		const findingPath = "docs/findings/2026-08-06-unknown-lifecycle.md"
		writeFindingsArtifact(t, repoRoot, findingPath, "---\nstatus: blocked\ncreated_at: 2026-08-06\nupdated_at: 2026-08-06\n---\n\n# Unknown lifecycle\n")

		result := checkFindingsCarrier(t, repoRoot)
		finding := requireRenderedFinding(t, result, speccheck.CodeFindingLifecycle, findingPath, 2)
		if !strings.Contains(finding.Summary, `"blocked"`) {
			t.Fatalf("summary = %q, want offending value %q", finding.Summary, "blocked")
		}
	})

	t.Run("green exact lifecycle values", func(t *testing.T) {
		t.Parallel()

		repoRoot := writeFindingsCarrier(t)
		for _, status := range []string{"pending", "partial", "deferred", "done"} {
			path := "docs/findings/2026-08-06-" + status + ".md"
			writeFindingsArtifact(t, repoRoot, path, "---\nstatus: "+status+"\ncreated_at: 2026-08-06\nupdated_at: 2026-08-06\n---\n\n# "+status+"\n")
		}

		result := checkFindingsCarrier(t, repoRoot)
		if findings := findingsWithCode(result, speccheck.CodeFindingLifecycle); len(findings) != 0 {
			t.Fatalf("%s findings = %#v, want exact lifecycle values accepted", speccheck.CodeFindingLifecycle, findings)
		}
	})

	t.Run("presence aware without findings directory", func(t *testing.T) {
		t.Parallel()

		repoRoot := writeFindingsCarrier(t)
		result := checkFindingsCarrier(t, repoRoot)
		if findings := findingsWithCode(result, speccheck.CodeFindingLifecycle); len(findings) != 0 {
			t.Fatalf("%s findings = %#v, want silent skip", speccheck.CodeFindingLifecycle, findings)
		}
		if !hasSkip(result, speccheck.CodeFindingLifecycle, "docs/findings") {
			t.Fatalf("Skipped = %#v, want %s missing docs/findings", result.Skipped, speccheck.CodeFindingLifecycle)
		}
	})
}

func TestCheckRollupMember(t *testing.T) {
	t.Parallel()

	t.Run("red unresolved declared member", func(t *testing.T) {
		t.Parallel()

		repoRoot := writeFindingsCarrier(t)
		const rollupPath = "docs/findings/2026-08-06-rollup.md"
		writeFindingsArtifact(t, repoRoot, rollupPath, "---\nstatus: pending\nkind: rollup\nmembers:\n  - missing-member.md\ncreated_at: 2026-08-06\nupdated_at: 2026-08-06\n---\n\n# Rollup\n")

		result := checkFindingsCarrier(t, repoRoot)
		finding := requireRenderedFinding(t, result, speccheck.CodeRollupMember, rollupPath, 5)
		if !strings.Contains(finding.Summary, "missing-member.md") {
			t.Fatalf("summary = %q, want unresolved member", finding.Summary)
		}
	})

	t.Run("green active and archived members", func(t *testing.T) {
		t.Parallel()

		repoRoot := writeFindingsCarrier(t)
		const rollup = "2026-08-06-rollup.md"
		writeFindingsArtifact(t, repoRoot, "docs/findings/"+rollup, "---\nstatus: pending\nkind: rollup\nmembers:\n  - 2026-08-06-active.md\n  - 2026-08-06-archived.md\ncreated_at: 2026-08-06\nupdated_at: 2026-08-06\n---\n\n# Rollup\n")
		writeFindingsArtifact(t, repoRoot, "docs/findings/2026-08-06-active.md", "---\nstatus: deferred\ncreated_at: 2026-08-06\nupdated_at: 2026-08-06\n---\n\n# Active\n")
		writeFindingsArtifact(t, repoRoot, "docs/findings/_archived/2026-08-06-archived.md", "---\nstatus: done\ncreated_at: 2026-08-06\nupdated_at: 2026-08-06\nabsorbed_by: "+rollup+"\n---\n\n# Archived\n")

		result := checkFindingsCarrier(t, repoRoot)
		if findings := findingsWithCode(result, speccheck.CodeRollupMember); len(findings) != 0 {
			t.Fatalf("%s findings = %#v, want active and archived members resolved", speccheck.CodeRollupMember, findings)
		}
	})

	t.Run("presence aware without rollups", func(t *testing.T) {
		t.Parallel()

		repoRoot := writeFindingsCarrier(t)
		writeFindingsArtifact(t, repoRoot, "docs/findings/2026-08-06-ordinary.md", "---\nstatus: pending\ncreated_at: 2026-08-06\nupdated_at: 2026-08-06\n---\n\n# Ordinary finding\n")

		result := checkFindingsCarrier(t, repoRoot)
		if findings := findingsWithCode(result, speccheck.CodeRollupMember); len(findings) != 0 {
			t.Fatalf("%s findings = %#v, want silent skip", speccheck.CodeRollupMember, findings)
		}
		if !hasSkip(result, speccheck.CodeRollupMember, "rollup") {
			t.Fatalf("Skipped = %#v, want %s without rollups", result.Skipped, speccheck.CodeRollupMember)
		}
	})
}

func TestCheckArchiveLicense(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name        string
		frontmatter string
		wantSummary string
		wantLine    int
	}{
		{
			name:        "red missing license",
			frontmatter: "---\nstatus: done\ncreated_at: 2026-08-06\nupdated_at: 2026-08-06\n---\n",
			wantSummary: "no absorbed_by license",
			wantLine:    1,
		},
		{
			name:        "red unresolved license",
			frontmatter: "---\nstatus: done\ncreated_at: 2026-08-06\nupdated_at: 2026-08-06\nabsorbed_by: missing-owner\n---\n",
			wantSummary: `"missing-owner"`,
			wantLine:    5,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repoRoot := writeFindingsCarrier(t)
			const archivedPath = "docs/findings/_archived/2026-08-06-archived.md"
			writeFindingsArtifact(t, repoRoot, archivedPath, tt.frontmatter+"\n# Archived\n")

			result := checkFindingsCarrier(t, repoRoot)
			finding := requireRenderedFinding(t, result, speccheck.CodeArchiveLicense, archivedPath, tt.wantLine)
			if !strings.Contains(finding.Summary, tt.wantSummary) {
				t.Fatalf("summary = %q, want %q", finding.Summary, tt.wantSummary)
			}
		})
	}

	t.Run("green rollup active Spec and archived Spec licenses", func(t *testing.T) {
		t.Parallel()

		repoRoot := writeFindingsCarrier(t)
		const rollup = "2026-08-06-rollup.md"
		writeFindingsArtifact(t, repoRoot, "docs/findings/"+rollup, "---\nstatus: pending\nkind: rollup\nmembers:\n  - 2026-08-06-rollup-owned.md\n  - 2026-08-06-active-spec-owned.md\n  - 2026-08-06-archived-spec-owned.md\ncreated_at: 2026-08-06\nupdated_at: 2026-08-06\n---\n\n# Rollup\n")
		writeFindingsArtifact(t, repoRoot, "docs/specs/_archived/archived-spec/_prd.md", "---\nspec: archived-spec\nstatus: archived\n---\n\n# Archived Spec\n")
		for name, owner := range map[string]string{
			"2026-08-06-rollup-owned.md":        rollup,
			"2026-08-06-active-spec-owned.md":   findingsCarrierSlug,
			"2026-08-06-archived-spec-owned.md": "archived-spec",
		} {
			writeFindingsArtifact(t, repoRoot, "docs/findings/_archived/"+name, "---\nstatus: done\ncreated_at: 2026-08-06\nupdated_at: 2026-08-06\nabsorbed_by: "+owner+"\n---\n\n# Archived\n")
		}

		result := checkFindingsCarrier(t, repoRoot)
		if findings := findingsWithCode(result, speccheck.CodeArchiveLicense); len(findings) != 0 {
			t.Fatalf("%s findings = %#v, want all declared license classes resolved", speccheck.CodeArchiveLicense, findings)
		}
	})

	t.Run("presence aware without archive", func(t *testing.T) {
		t.Parallel()

		repoRoot := writeFindingsCarrier(t)
		writeFindingsArtifact(t, repoRoot, "docs/findings/2026-08-06-ordinary.md", "---\nstatus: pending\ncreated_at: 2026-08-06\nupdated_at: 2026-08-06\n---\n\n# Ordinary finding\n")

		result := checkFindingsCarrier(t, repoRoot)
		if findings := findingsWithCode(result, speccheck.CodeArchiveLicense); len(findings) != 0 {
			t.Fatalf("%s findings = %#v, want silent skip", speccheck.CodeArchiveLicense, findings)
		}
		if !hasSkip(result, speccheck.CodeArchiveLicense, "docs/findings/_archived") {
			t.Fatalf("Skipped = %#v, want %s missing archive", result.Skipped, speccheck.CodeArchiveLicense)
		}
	})
}

func writeFindingsCarrier(t *testing.T) string {
	t.Helper()

	repoRoot := t.TempDir()
	specDir := filepath.Join(repoRoot, "docs", "specs", findingsCarrierSlug)
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatalf("create findings fixture carrier: %v", err)
	}
	return repoRoot
}

func writeFindingsArtifact(t *testing.T, repoRoot, relative, content string) {
	t.Helper()

	path := filepath.Join(repoRoot, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create findings fixture directory for %q: %v", relative, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write findings fixture %q: %v", relative, err)
	}
}

func checkFindingsCarrier(t *testing.T, repoRoot string) speccheck.Result {
	t.Helper()

	result, err := speccheck.Check(filepath.Join(repoRoot, "docs", "specs"), repoRoot, findingsCarrierSlug)
	if err != nil {
		t.Fatalf("Check(findings carrier) error = %v", err)
	}
	return result
}

func requireRenderedFinding(t *testing.T, result speccheck.Result, code, path string, line int) speccheck.Finding {
	t.Helper()

	finding := requireFinding(t, result, code)
	if finding.Severity != speccheck.SeverityError {
		t.Fatalf("severity = %q, want %q", finding.Severity, speccheck.SeverityError)
	}
	if !hasExactLocation(finding, path, line) {
		t.Fatalf("locations = %#v, want %s:%d", finding.Where, path, line)
	}
	if strings.TrimSpace(finding.Fix) == "" {
		t.Fatal("finding has no concrete fix")
	}
	report := speccheck.RenderText(result)
	if !strings.Contains(report, finding.Code+": "+finding.Summary) {
		t.Fatalf("rendered report does not contain finding summary:\n%s", report)
	}
	if !strings.Contains(report, "  fix: "+finding.Fix) {
		t.Fatalf("rendered report does not contain fix line:\n%s", report)
	}
	return finding
}

func TestCheckADRUnlisted(t *testing.T) {
	t.Parallel()

	result := checkFixture(t, "citation-dirty")
	findings := findingsWithCode(result, speccheck.CodeADRUnlisted)
	if len(findings) != 1 {
		t.Fatalf("%s findings = %#v, want exactly one", speccheck.CodeADRUnlisted, findings)
	}
	finding := findings[0]
	if finding.Severity != speccheck.SeverityError {
		t.Fatalf("severity = %q, want %q", finding.Severity, speccheck.SeverityError)
	}
	for _, path := range []string{
		"docs/specs/citation-dirty/_prd.md",
		"docs/specs/citation-dirty/_techspec.md",
	} {
		if !hasLocation(finding, path) {
			t.Errorf("locations = %#v, want %q", finding.Where, path)
		}
	}
}

func TestCheckADRClosureDepthOne(t *testing.T) {
	t.Parallel()

	result := checkFixture(t, "citation-dirty")
	findings := findingsWithCode(result, speccheck.CodeADRRelated)
	if len(findings) != 1 {
		t.Fatalf("%s findings = %#v, want related ADR once", speccheck.CodeADRRelated, findings)
	}
	if findings[0].Severity != speccheck.SeverityGap {
		t.Fatalf("severity = %q, want %q", findings[0].Severity, speccheck.SeverityGap)
	}
	if !strings.Contains(findings[0].Summary, "ADR-0102") {
		t.Fatalf("summary = %q, want ADR-0102", findings[0].Summary)
	}
	for _, listed := range []string{"ADR-0100", "ADR-0101"} {
		for _, finding := range findings {
			if strings.HasPrefix(finding.Summary, listed+" ") {
				t.Errorf("listed %s reported as related: %#v", listed, finding)
			}
		}
	}
}

func TestCitationReportsAnUnsupportedClaim(t *testing.T) {
	t.Parallel()

	repoRoot := writeSemanticCitationFixture(t, "_prd.md", "ADR-0083 makes `make verify` the only authoritative gate.")
	result := checkSemanticCitationFixture(t, repoRoot)
	finding := requireFinding(t, result, speccheck.CodeCitationUnsupported)

	const claimingSentence = "ADR-0083 makes `make verify` the only authoritative gate."
	for _, text := range []string{claimingSentence, "Adopted sources move to their owning Spec"} {
		if !strings.Contains(finding.Summary, text) {
			t.Errorf("summary = %q, want quoted evidence %q", finding.Summary, text)
		}
	}
	for _, path := range []string{
		"docs/specs/semantic-citation/_prd.md",
		"docs/adr/0083-adopted-sources-move-to-their-owning-spec.md",
	} {
		if !hasLocation(finding, path) {
			t.Errorf("locations = %#v, want %q", finding.Where, path)
		}
	}
}

func TestCitationAcceptsASupportedClaim(t *testing.T) {
	t.Parallel()

	repoRoot := writeSemanticCitationFixture(t, "_techspec.md", "ADR-0096 has the QA gate prove machine facts before it spends an Agent turn.")
	result := checkSemanticCitationFixture(t, repoRoot)
	if findings := findingsWithCode(result, speccheck.CodeCitationUnsupported); len(findings) != 0 {
		t.Fatalf("%s findings = %#v, want the cited record to support the claim", speccheck.CodeCitationUnsupported, findings)
	}
}

func TestCitationBareListingIsNotAClaim(t *testing.T) {
	t.Parallel()

	repoRoot := writeSemanticCitationFixture(t, "_prd.md", "Active ADRs: ADR-0083, ADR-0096.")
	result := checkSemanticCitationFixture(t, repoRoot)
	if findings := findingsWithCode(result, speccheck.CodeCitationUnsupported); len(findings) != 0 {
		t.Fatalf("%s findings = %#v, want a bare ADR listing ignored", speccheck.CodeCitationUnsupported, findings)
	}
}

func TestCitationReportsHowManyClaimsItResolved(t *testing.T) {
	t.Parallel()

	repoRoot := writeSemanticCitationFixture(t, "_prd.md", "ADR-0083 makes `make verify` the only authoritative gate.")
	artifact := "docs/specs/semantic-citation/_prd.md"
	content, err := os.ReadFile(filepath.Join(repoRoot, filepath.FromSlash(artifact)))
	if err != nil {
		t.Fatalf("read semantic citation fixture: %v", err)
	}
	claims := speccheck.CitationClaims(artifact, content)
	if len(claims) != 1 {
		t.Fatalf("parsed claims = %#v, want exactly one attribution", claims)
	}
	resolved, err := speccheck.ResolvedCitationClaimCount(repoRoot, claims)
	if err != nil {
		t.Fatalf("ResolvedCitationClaimCount() error = %v", err)
	}
	if resolved != 1 {
		t.Fatalf("resolved claims = %d, want 1", resolved)
	}

	bareClaims := speccheck.CitationClaims(artifact, []byte("Active ADRs: ADR-0083, ADR-0096.\n"))
	if len(bareClaims) != 0 {
		t.Fatalf("bare listing claims = %#v, want none", bareClaims)
	}
	resolved, err = speccheck.ResolvedCitationClaimCount(repoRoot, bareClaims)
	if err != nil {
		t.Fatalf("ResolvedCitationClaimCount(bare listing) error = %v", err)
	}
	if resolved != 0 {
		t.Fatalf("bare listing resolved claims = %d, want 0", resolved)
	}
}

func writeSemanticCitationFixture(t *testing.T, artifactName, sentence string) string {
	t.Helper()

	repoRoot := t.TempDir()
	for _, source := range []string{
		"docs/agents/agent-instructions.md",
		"docs/agents/cli.md",
		"docs/agents/domain.md",
	} {
		writeCitationFixtureFile(t, repoRoot, source, "# Fixture source\n")
	}
	writeCitationFixtureFile(t, repoRoot, "docs/adr/0083-adopted-sources-move-to-their-owning-spec.md", `---
status: accepted
---

# Adopted sources move to their owning Spec

An adopted source moves into the owning Spec so one authoritative home cannot drift.
`)
	const adr0096Path = "docs/adr/0096-the-qa-gate-proves-machine-facts-before-it-spends-an-agent-turn.md"
	adr0096, err := os.ReadFile(filepath.Join(characterizationRepositoryRoot(t), filepath.FromSlash(adr0096Path)))
	if err != nil {
		t.Fatalf("read real ADR-0096: %v", err)
	}
	writeCitationFixtureFile(t, repoRoot, adr0096Path, string(adr0096))

	const prd = `---
spec: semantic-citation
status: active
---

# Semantic citation

## Project Constraints

- Identifier strategy: not applicable — no persisted identity. Source: ` + "`docs/agents/domain.md`" + `.
- Authentication and HTTP: not applicable — file reads only. Source: ` + "`docs/agents/cli.md`" + `.
- Active ADR obligations: applicable — ADR-0083 and ADR-0096 are fixture decisions. Source: ` + "`docs/agents/domain.md`" + `.
- Tooling authority: not applicable — ordinary source only. Source: ` + "`docs/agents/agent-instructions.md`" + `.
`
	writeCitationFixtureFile(t, repoRoot, "docs/specs/semantic-citation/_prd.md", prd)
	writeCitationFixtureFile(t, repoRoot, "docs/specs/semantic-citation/"+artifactName, prd+"\n"+sentence+"\n")
	return repoRoot
}

func checkSemanticCitationFixture(t *testing.T, repoRoot string) speccheck.Result {
	t.Helper()

	result, err := speccheck.Check(filepath.Join(repoRoot, "docs", "specs"), repoRoot, "semantic-citation")
	if err != nil {
		t.Fatalf("Check(semantic-citation) error = %v", err)
	}
	return result
}

func TestCheckCoverageUnmapped(t *testing.T) {
	t.Parallel()

	result := checkFixture(t, "coverage-unmapped")
	findings := findingsWithCode(result, speccheck.CodeCoverageUnmapped)
	if len(findings) != 4 {
		t.Fatalf("%s findings = %#v, want one per four Core Features", speccheck.CodeCoverageUnmapped, findings)
	}
	for _, finding := range findings {
		if !strings.Contains(finding.Summary, "Core Feature") {
			t.Errorf("unexpected uncovered unit: %q", finding.Summary)
		}
	}
}

func TestCheckCoverageRange(t *testing.T) {
	t.Parallel()

	result := checkFixture(t, "coverage-range")
	if findings := findingsWithCode(result, speccheck.CodeCoverageUnmapped); len(findings) != 0 {
		t.Fatalf("range produced findings = %#v, want none", findings)
	}
	if findings := findingsWithCode(result, speccheck.CodeCoverageUntasked); len(findings) != 0 {
		t.Fatalf("Task reference range produced findings = %#v, want none", findings)
	}
}

func TestCheckCoverageUntasked(t *testing.T) {
	t.Parallel()

	result := checkFixture(t, "coverage-untasked")
	findings := findingsWithCode(result, speccheck.CodeCoverageUntasked)
	if len(findings) != 1 || !strings.Contains(findings[0].Summary, "Core Feature 4") {
		t.Fatalf("%s findings = %#v, want only Core Feature 4", speccheck.CodeCoverageUntasked, findings)
	}
}

func TestCheckReferenceUnresolved(t *testing.T) {
	t.Parallel()

	result := checkFixture(t, "reference-unresolved")
	finding := requireFinding(t, result, speccheck.CodeReferenceUnresolved)
	const missingPath = "missing/guide.md"
	if !strings.Contains(finding.Summary, missingPath) {
		t.Fatalf("summary = %q, want missing path %q", finding.Summary, missingPath)
	}
	if !hasLocation(finding, "docs/specs/reference-unresolved/task_01.md") {
		t.Errorf("locations = %#v, want declaring Task line", finding.Where)
	}
	if !hasLocation(finding, missingPath) {
		t.Errorf("locations = %#v, want unresolved path", finding.Where)
	}
}

func TestCheckReferenceIndexUnresolved(t *testing.T) {
	t.Parallel()

	result := checkFixture(t, "reference-index-unresolved")
	finding := requireFinding(t, result, speccheck.CodeReferenceUnresolved)
	if !strings.Contains(finding.Summary, "missing-note.md") {
		t.Fatalf("summary = %q, want missing-note.md", finding.Summary)
	}
	if !hasLocation(finding, "docs/specs/reference-index-unresolved/references/_index.md") {
		t.Errorf("locations = %#v, want declaring reference index line", finding.Where)
	}
}

func TestCheckCoverageUntaskedSkippedWithoutTaskGraph(t *testing.T) {
	t.Parallel()

	result := checkFixture(t, "no-taskgraph")
	if findings := findingsWithCode(result, speccheck.CodeCoverageUntasked); len(findings) != 0 {
		t.Fatalf("%s findings = %#v, want none", speccheck.CodeCoverageUntasked, findings)
	}
	if !hasSkip(result, speccheck.CodeCoverageUntasked, "_tasks.md") {
		t.Fatalf("Skipped = %#v, want %s missing _tasks.md", result.Skipped, speccheck.CodeCoverageUntasked)
	}
}

func TestCheckCoverageUnmappedSkippedWithoutTechSpec(t *testing.T) {
	t.Parallel()

	result := checkFixture(t, "no-techspec")
	if findings := findingsWithCode(result, speccheck.CodeCoverageUnmapped); len(findings) != 0 {
		t.Fatalf("%s findings = %#v, want none", speccheck.CodeCoverageUnmapped, findings)
	}
	if !hasSkip(result, speccheck.CodeCoverageUnmapped, "_techspec.md") {
		t.Fatalf("Skipped = %#v, want %s missing _techspec.md", result.Skipped, speccheck.CodeCoverageUnmapped)
	}
}

func TestCheckADRCorpusAbsentSkipsADRDetectors(t *testing.T) {
	t.Parallel()

	repoRoot, err := filepath.Abs("testdata/no-adr-repo")
	if err != nil {
		t.Fatalf("resolve fixture repository: %v", err)
	}
	result, err := speccheck.Check(filepath.Join(repoRoot, "docs", "specs"), repoRoot, "no-adr")
	if err != nil {
		t.Fatalf("Check(no-adr) error = %v", err)
	}
	for _, code := range []string{speccheck.CodeADRUnlisted, speccheck.CodeADRRelated} {
		if !hasSkip(result, code, "docs/adr") {
			t.Errorf("Skipped = %#v, want %s missing docs/adr", result.Skipped, code)
		}
	}
}

func TestCheckCitationCoverageErrorLocations(t *testing.T) {
	t.Parallel()

	for _, slug := range []string{
		"citation-dirty",
		"coverage-unmapped",
		"coverage-untasked",
		"reference-unresolved",
		"reference-index-unresolved",
	} {
		slug := slug
		t.Run(slug, func(t *testing.T) {
			t.Parallel()

			result := checkFixture(t, slug)
			for _, finding := range result.Findings {
				if finding.Code == speccheck.CodeADRRelated {
					continue
				}
				if finding.Severity != speccheck.SeverityError {
					t.Errorf("%s severity = %q, want %q", finding.Code, finding.Severity, speccheck.SeverityError)
				}
				if len(finding.Where) < 2 {
					t.Errorf("%s locations = %#v, want both sides", finding.Code, finding.Where)
				}
			}
		})
	}
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
