// Suite: Stated HTTP default
// Invariant: public statements of the HTTP suggestion agree with the decision catalog.
// Boundary IN: embedded catalog and the public guide's suggested-values row.
// Boundary OUT: decision prompting and repository mutation.

package cli

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/baseline"
)

type httpDefaultStatement struct{ source, text, format string }

func httpDefaultStatementFindings(mode string, modes []string, statements []httpDefaultStatement) []string {
	var findings []string
	for _, statement := range statements {
		text := strings.Join(strings.Fields(statement.text), " ")
		if !strings.Contains(text, fmt.Sprintf(statement.format, mode)) {
			findings = append(findings, fmt.Sprintf("%s: missing HTTP default %s", statement.source, mode))
		}
		for _, other := range modes {
			if other != mode && strings.Contains(text, fmt.Sprintf(statement.format, other)) {
				findings = append(findings, fmt.Sprintf("%s: states another HTTP default %s", statement.source, other))
			}
		}
	}
	return findings
}

func TestEveryStatementOfTheHTTPDefaultNamesTheCatalogDefault(t *testing.T) {
	t.Parallel()
	catalog, err := baseline.LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	entry, ok := catalog.Decision("http.contract")
	if !ok {
		t.Fatal("HTTP decision is absent")
	}
	var decision struct {
		Default struct {
			Mode string `json:"mode"`
		} `json:"default"`
		Modes []string `json:"modes"`
	}
	if err := json.Unmarshal(entry.Data, &decision); err != nil {
		t.Fatal(err)
	}
	if decision.Default.Mode == "" || len(decision.Modes) == 0 {
		t.Fatal("HTTP decision has no default or modes")
	}
	contract, ok := catalog.Asset("contract-v1.json")
	if !ok {
		t.Fatal("Baseline contract is absent")
	}
	guide := readBaselineDocumentation(t, filepath.Join(baselineDocumentationRepoRoot(), "docs/user-guide/context-driven-development.md"))
	statements := []httpDefaultStatement{
		{"contract-v1.json", string(contract.Data), "the interactive workflow suggests %s until"},
		{"context-driven-development.md", guide, "| HTTP contract | `%s` |"},
	}
	if findings := httpDefaultStatementFindings(decision.Default.Mode, decision.Modes, statements); len(findings) != 0 {
		t.Fatalf("HTTP statement findings: %v", findings)
	}
}

func TestAStatementThatNamesAnotherHTTPDefaultIsReported(t *testing.T) {
	t.Parallel()
	statements := []httpDefaultStatement{
		{"contract", "the interactive workflow suggests Post-only until the maintainer confirms it", "the interactive workflow suggests %s until"},
		{"guide", "| HTTP contract | `Post-only` |", "| HTTP contract | `%s` |"},
	}
	want := []string{
		"contract: missing HTTP default REST", "contract: states another HTTP default Post-only",
		"guide: missing HTTP default REST", "guide: states another HTTP default Post-only",
	}
	if findings := httpDefaultStatementFindings("REST", []string{"REST", "Post-only"}, statements); !reflect.DeepEqual(findings, want) {
		t.Fatalf("findings = %v, want %v", findings, want)
	}
}
