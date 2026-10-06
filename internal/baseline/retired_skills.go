package baseline

import "sort"

var retiredSkills = map[string]struct{}{"council": {}, "the-fool": {}}

// RetiredSkills returns the Retired Skills in lexical order.
func RetiredSkills() []string {
	names := make([]string, 0, len(retiredSkills))
	for name := range retiredSkills {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
