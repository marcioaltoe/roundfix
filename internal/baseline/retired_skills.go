package baseline

import "sort"

var retiredSkills = map[string]struct{}{
	"council":                {},
	"design-patterns":        {},
	"diagnosing-bugs":        {},
	"fixing-metadata":        {},
	"game-changing-features": {},
	"implement-spec":         {},
	"refactoring-analysis":   {},
	"tdd":                    {},
	"the-fool":               {},
	"to-prompt":              {},
}

// RetiredSkills returns the Retired Skills in lexical order.
func RetiredSkills() []string {
	names := make([]string, 0, len(retiredSkills))
	for name := range retiredSkills {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
