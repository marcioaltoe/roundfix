// Suite: Stated HTTP default
// Invariant: public statements of the HTTP suggestion agree with the decision catalog.
// Boundary IN: embedded catalog and the public guide's suggested-values row.
// Boundary OUT: decision prompting and repository mutation.

package cli

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
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
