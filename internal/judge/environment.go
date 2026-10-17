package judge

import (
	"slices"
	"strings"
)

// RunMarkerVariable names the Run whose Agent and Verification processes must not use the judge.
const RunMarkerVariable = "ROUNDFIX_RUN_ID"

// WithoutKeyVariables returns a copy of environ without any entry whose name
// is one of Questions.KeyVariables of the embedded questions or TYPESAFE_API_KEY.
func WithoutKeyVariables(environ []string) []string {
	questions, err := Load()
	if err != nil {
		// Invalid embedded questions must never let credentials reach a child.
		panic(err)
	}
	variables := append(questions.KeyVariables(), "TYPESAFE_API_KEY")
	filtered := make([]string, 0, len(environ))
	for _, entry := range environ {
		name, _, _ := strings.Cut(entry, "=")
		if !slices.Contains(variables, name) {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}
