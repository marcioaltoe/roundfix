package delivery

import (
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

// WorkflowRun is the latest Actions run for one workflow file and event of a head.
type WorkflowRun struct {
	ID, Name, Path, Event, Status, Conclusion string
	Attempt                                   int
	URL                                       string
}

// ExpectedWorkflows lists the head's workflow files that trigger for every pull
// request into base. Filtered and unparsable files create no expectation.
func ExpectedWorkflows(files map[string][]byte, base string) []string {
	var expected []string
	for path, data := range files {
		var document yaml.Node
		if err := yaml.Unmarshal(data, &document); err != nil || len(document.Content) != 1 {
			continue
		}
		root := document.Content[0]
		if root.Kind != yaml.MappingNode {
			continue
		}
		for i := 0; i < len(root.Content); i += 2 {
			if root.Content[i].Value == "on" && workflowTriggersForBase(root.Content[i+1], base) {
				expected = append(expected, path)
				break
			}
		}
	}
	slices.Sort(expected)
	return expected
}

func workflowTriggersForBase(node *yaml.Node, base string) bool {
	switch node.Kind {
	case yaml.ScalarNode:
		return node.Value == "pull_request" || node.Value == "pull_request_target" || node.Value == "push"
	case yaml.SequenceNode:
		for _, event := range node.Content {
			if event.Kind == yaml.ScalarNode && workflowTriggersForBase(event, base) {
				return true
			}
		}
	case yaml.MappingNode:
		for i := 0; i < len(node.Content); i += 2 {
			event, filters := node.Content[i].Value, node.Content[i+1]
			switch event {
			case "push":
				if emptyWorkflowTrigger(filters) {
					return true
				}
			case "pull_request", "pull_request_target":
				if emptyWorkflowTrigger(filters) || pullRequestTriggerForBase(filters, base) {
					return true
				}
			}
		}
	}
	return false
}

func emptyWorkflowTrigger(node *yaml.Node) bool {
	return (node.Kind == yaml.ScalarNode && (node.Tag == "!!null" || node.Value == "")) ||
		(node.Kind == yaml.MappingNode && len(node.Content) == 0)
}

func pullRequestTriggerForBase(node *yaml.Node, base string) bool {
	if node.Kind != yaml.MappingNode {
		return false
	}
	for i := 0; i < len(node.Content); i += 2 {
		key, value := node.Content[i].Value, node.Content[i+1]
		entries := workflowTriggerStrings(value)
		switch key {
		case "branches":
			if !slices.Contains(entries, base) {
				return false
			}
			for _, branch := range entries {
				if strings.HasPrefix(branch, "!") {
					return false
				}
			}
		case "types":
			if !slices.Contains(entries, "opened") || !slices.Contains(entries, "synchronize") {
				return false
			}
		default:
			return false
		}
	}
	return true
}

func workflowTriggerStrings(node *yaml.Node) []string {
	if node.Kind == yaml.ScalarNode && node.Tag == "!!str" {
		return []string{node.Value}
	}
	if node.Kind != yaml.SequenceNode {
		return nil
	}
	entries := make([]string, 0, len(node.Content))
	for _, entry := range node.Content {
		if entry.Kind != yaml.ScalarNode || entry.Tag != "!!str" {
			return nil
		}
		entries = append(entries, entry.Value)
	}
	return entries
}
