package judge

import (
	"slices"
	"strings"
)

// WithoutKeyVariables returns a copy of environ without any entry whose name
// is one of Questions.KeyVariables of the embedded questions.
func WithoutKeyVariables(environ []string) []string {
	questions, err := Load()
	if err != nil {
		// Invalid embedded questions must never let credentials reach a child.
		panic(err)
	}
	variables := questions.KeyVariables()
	filtered := make([]string, 0, len(environ))
	for _, entry := range environ {
		name, _, _ := strings.Cut(entry, "=")
		if !slices.Contains(variables, name) {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}
