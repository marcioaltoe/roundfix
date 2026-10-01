package spec

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// ParseRequiredSpecs reads only the prerequisite declaration from a manifest,
// including manifests read from a Git ref rather than an active worktree.
func ParseRequiredSpecs(content []byte, slug string) ([]string, error) {
	frontmatter, _, err := splitFrontmatter(content)
	if err != nil {
		return nil, fmt.Errorf("_tasks.md: %w", err)
	}
	var declaration struct {
		Requires yaml.Node `yaml:"requires"`
	}
	if err := yaml.Unmarshal(frontmatter, &declaration); err != nil {
		return nil, fmt.Errorf("_tasks.md: %w", err)
	}
	return parseRequiredSpecs(declaration.Requires, slug)
}

func parseRequiredSpecs(node yaml.Node, slug string) ([]string, error) {
	if node.Kind == 0 {
		return nil, nil
	}
	if node.Kind != yaml.SequenceNode {
		return nil, fmt.Errorf("_tasks.md: requires must be a list of strings")
	}
	var requires []string
	seen := make(map[string]bool)
	for _, entry := range node.Content {
		if entry.Kind != yaml.ScalarNode || entry.Tag != "!!str" {
			return nil, fmt.Errorf("_tasks.md: requires must be a list of strings")
		}
		name := strings.TrimSpace(entry.Value)
		if name == "" || seen[name] || name == slug {
			return nil, fmt.Errorf("_tasks.md: requires entry %q is empty, duplicated, or names its own Spec", name)
		}
		seen[name] = true
		requires = append(requires, name)
	}
	return requires, nil
}
