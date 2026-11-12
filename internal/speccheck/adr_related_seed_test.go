package speccheck_test

import (
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/speccheck"
)

func TestANotApplicableADRDrawsInNoCitingADR(t *testing.T) {
	t.Parallel()

	for _, relation := range []string{"not applicable; unchanged.", "Not applicable — unchanged."} {
		t.Run(relation, func(t *testing.T) {
			t.Parallel()

			repoRoot := writeADRRelatedSeedFixture(t, relation)
			result := checkADRRelatedSeeds(t, repoRoot)
			for _, code := range []string{speccheck.CodeADRRelated, speccheck.CodeADRUnlisted, speccheck.CodeADRAlternativesUnrelated} {
				if findings := findingsWithCode(result, code); len(findings) != 0 {
					t.Errorf("%s findings = %#v, want none", code, findings)
				}
			}
		})
	}
}

func TestAnAppliedADRStillDrawsInItsCitingADR(t *testing.T) {
	t.Parallel()

	for _, relation := range []string{
		"applied; kept.",
		"kept as recorded, though its alternatives are not applicable here.",
	} {
		t.Run(relation, func(t *testing.T) {
			t.Parallel()

			repoRoot := writeADRRelatedSeedFixture(t, relation)
			result := checkADRRelatedSeeds(t, repoRoot)
			finding := requireADRRelatedSeedFinding(t, result, "0003", "0002")
			const wantFix = `Account for ADR-0003 in the Active ADR obligations row in docs/specs/adr-related-seeds/_prd.md by listing it or recording why it does not apply; when it does not apply, begin its ADR Alternatives line in the TechSpec with "not applicable" or "não se aplica", so the ADRs that cite only it are not drawn in.`
			if finding.Fix != wantFix {
				t.Errorf("fix = %q, want %q", finding.Fix, wantFix)
			}
		})
	}
}

func TestWithoutATechSpecEveryListedADRSeeds(t *testing.T) {
	t.Parallel()

	repoRoot := writeADRRelatedSeedFixture(t, "")
	result := checkADRRelatedSeeds(t, repoRoot)
	requireADRRelatedSeedFinding(t, result, "0003", "0002")
}

func TestAnADRCitingAnAppliedADRIsStillRelated(t *testing.T) {
	t.Parallel()

	repoRoot := writeADRRelatedSeedFixture(t, "not applicable; unchanged.")
	writeCitationFixtureFile(t, repoRoot, "docs/adr/0004-fixture-decision.md", "---\nstatus: accepted\n---\n\n# Fixture decision 0004\n\nADR-0001 and ADR-0002.\n")
	result := checkADRRelatedSeeds(t, repoRoot)
	requireADRRelatedSeedFinding(t, result, "0004", "0001")
}

func writeADRRelatedSeedFixture(t *testing.T, relation string) string {
	t.Helper()

	repoRoot := t.TempDir()
	for _, source := range []string{"agent-instructions", "cli", "domain"} {
		writeCitationFixtureFile(t, repoRoot, "docs/agents/"+source+".md", "# Fixture source\n")
	}
	for _, decision := range []struct {
		number   string
		citation string
	}{
		{number: "0001"},
		{number: "0002", citation: "ADR-0001."},
		{number: "0003", citation: "ADR-0002."},
	} {
		writeCitationFixtureFile(t, repoRoot, "docs/adr/"+decision.number+"-fixture-decision.md",
			"---\nstatus: accepted\n---\n\n# Fixture decision "+decision.number+"\n\n"+decision.citation+"\n")
	}
	writeCitationFixtureFile(t, repoRoot, "docs/specs/adr-related-seeds/_prd.md", `---
spec: adr-related-seeds
status: active
---

# ADR related seeds

## Project Constraints

- Identifier strategy: not applicable — fixture only. Source: `+"`docs/agents/domain.md`"+`.
- Authentication and HTTP: not applicable — fixture only. Source: `+"`docs/agents/cli.md`"+`.
- Active ADR obligations: applicable — ADR-0001 and ADR-0002. Source: `+"`docs/agents/domain.md`"+`.
- Tooling authority: not applicable — ordinary source only. Source: `+"`docs/agents/agent-instructions.md`"+`.
`)
	if relation != "" {
		writeCitationFixtureFile(t, repoRoot, "docs/specs/adr-related-seeds/_techspec.md",
			"# ADR related seeds\n\n## ADR Alternatives\n\n- ADR-0001: applied; kept.\n- ADR-0002: "+relation+"\n")
	}
	return repoRoot
}

func checkADRRelatedSeeds(t *testing.T, repoRoot string) speccheck.Result {
	t.Helper()

	result, err := speccheck.Check(filepath.Join(repoRoot, "docs", "specs"), repoRoot, "adr-related-seeds")
	if err != nil {
		t.Fatalf("Check(adr-related-seeds) error = %v", err)
	}
	return result
}

func requireADRRelatedSeedFinding(t *testing.T, result speccheck.Result, number, seed string) speccheck.Finding {
	t.Helper()

	findings := findingsWithCode(result, speccheck.CodeADRRelated)
	if len(findings) != 1 {
		t.Fatalf("%s findings = %#v, want one for ADR-%s", speccheck.CodeADRRelated, findings, number)
	}
	finding := findings[0]
	if !strings.HasPrefix(finding.Summary, "ADR-"+number+" ") || !strings.Contains(finding.Summary, "cites listed ADR-"+seed+" ") {
		t.Errorf("summary = %q, want ADR-%s citing applied ADR-%s", finding.Summary, number, seed)
	}
	if finding.Severity != speccheck.SeverityGap {
		t.Errorf("severity = %q, want %q", finding.Severity, speccheck.SeverityGap)
	}
	if len(finding.Where) != 2 || !hasLocation(finding, "docs/adr/"+number+"-fixture-decision.md") || !hasLocation(finding, "docs/specs/adr-related-seeds/_prd.md") {
		t.Errorf("locations = %#v, want citing ADR and PRD row", finding.Where)
	}
	return finding
}
