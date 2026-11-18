package config

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

// Skills contains declarations applied to this repository's installed skills.
type Skills struct {
	Authoring AuthoringModels `yaml:"authoring"`
}
type AuthoringModels struct {
	Claude *AuthoringHarness `yaml:"claude"`
}
type AuthoringHarness struct {
	Idea     *AuthoringEntry `yaml:"idea"`
	PRD      *AuthoringEntry `yaml:"prd"`
	TechSpec *AuthoringEntry `yaml:"techspec"`
	Tasks    *AuthoringEntry `yaml:"tasks"`
	Research *ResearchEntry  `yaml:"research"`
}
type AuthoringFallback struct {
	Model           string `yaml:"model"`
	ReasoningEffort string `yaml:"reasoning_effort"`
}
type AuthoringEntry struct {
	Model           string             `yaml:"model"`
	ReasoningEffort string             `yaml:"reasoning_effort"`
	Fallback        *AuthoringFallback `yaml:"fallback"`
}
type ResearchEntry struct {
	AuthoringEntry `yaml:",inline"`
	Tools          []string `yaml:"tools"`
}

var authoringModelName = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:/@\[\]-]{0,199}$`)
var authoringMCPTool = regexp.MustCompile(`^mcp__[A-Za-z0-9_.*-]+$`)

func prepareAuthoringConfig(document *yaml.Node, warnings *configWarnings, source ProfileSource) error {
	if source == ProfileSourceUser && removeYAMLPath(document, []string{"delivery", "require_defect_class"}) {
		const key = "delivery.require_defect_class"
		if warnings != nil && !warnings.emitted[key] {
			warnings.emitted[key] = true
			fmt.Fprintln(warnings.stderr, "config: delivery.require_defect_class in User Config is ignored; set delivery.require_defect_class in Project Config")
		}
	}
	node, found := yamlValueAtPath(document, []string{"skills"})
	if !found {
		return nil
	}
	if source == ProfileSourceUser {
		removeYAMLPath(document, []string{"skills"})
		if warnings != nil && !warnings.emitted["skills"] {
			warnings.emitted["skills"] = true
			fmt.Fprintln(warnings.stderr, "config: skills in User Config is ignored; set skills in Project Config")
		}
		return nil
	}
	return validateAuthoringMapping(node, "skills")
}

func validateAuthoringMapping(node *yaml.Node, path string) error {
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("%s must be a mapping", path)
	}
	for i := 0; i < len(node.Content); i += 2 {
		key, value := node.Content[i].Value, node.Content[i+1]
		child := path + "." + key
		switch path {
		case "skills":
			if key != "authoring" {
				return fmt.Errorf("%s is not a supported key; expected authoring", child)
			}
		case "skills.authoring":
			if key != "claude" {
				return fmt.Errorf("%s is not a supported harness; this Roundfix supports claude", child)
			}
		case "skills.authoring.claude":
			if !slices.Contains([]string{"idea", "prd", "techspec", "tasks", "research"}, key) {
				return fmt.Errorf("%s is not a supported key; expected idea, prd, techspec, tasks or research", child)
			}
			if err := validateAuthoringEntry(value, child, key == "research"); err != nil {
				return err
			}
			continue
		}
		if err := validateAuthoringMapping(value, child); err != nil {
			return err
		}
	}
	return nil
}

func authoringNodeValue(node *yaml.Node) string {
	if node.Kind == yaml.ScalarNode {
		return node.Value
	}
	content, err := yaml.Marshal(node)
	if err != nil {
		return node.Value
	}
	return strings.TrimSpace(string(content))
}

func validateAuthoringEntry(node *yaml.Node, path string, research bool) error {
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("%s must be a mapping", path)
	}
	declared := false
	var fallback *yaml.Node
	model := ""
	for i := 0; i < len(node.Content); i += 2 {
		key, value := node.Content[i].Value, node.Content[i+1]
		raw := authoringNodeValue(value)
		switch key {
		case "model":
			if value.Kind != yaml.ScalarNode || value.Tag != "!!str" || !authoringModelName.MatchString(raw) {
				return fmt.Errorf("%s.model %q is not a model name; use letters, digits and . _ : / @ [ ] - without spaces", path, raw)
			}
			declared = true
			model = raw
		case "reasoning_effort":
			if value.Kind != yaml.ScalarNode || value.Tag != "!!str" || !slices.Contains([]string{"low", "medium", "high", "xhigh", "max"}, raw) {
				return fmt.Errorf("%s.reasoning_effort %q is not supported for claude; expected low, medium, high, xhigh or max", path, raw)
			}
			declared = true
		case "fallback":
			if research {
				return fmt.Errorf("%s.fallback is not a supported key; expected model, reasoning_effort or tools", path)
			}
			fallback = value
		case "tools":
			if !research {
				return fmt.Errorf("%s.%s is not a supported key; expected model, reasoning_effort or fallback", path, key)
			}
			if value.Kind != yaml.SequenceNode {
				return invalidAuthoringTool(raw)
			}
			seen := map[string]bool{}
			for _, tool := range value.Content {
				name := authoringNodeValue(tool)
				if tool.Kind != yaml.ScalarNode || tool.Tag != "!!str" || (name != "WebSearch" && name != "WebFetch" && !authoringMCPTool.MatchString(name)) {
					return invalidAuthoringTool(name)
				}
				if seen[name] {
					return fmt.Errorf("skills.authoring.claude.research.tools names %q twice", name)
				}
				seen[name] = true
			}
			declared = declared || len(value.Content) > 0
		default:
			expected := "model, reasoning_effort or fallback"
			if research {
				expected = "model, reasoning_effort or tools"
			}
			return fmt.Errorf("%s.%s is not a supported key; expected %s", path, key, expected)
		}
	}
	if fallback != nil {
		if err := validateAuthoringFallback(fallback, path, model); err != nil {
			return err
		}
	}
	if !declared {
		expected := "model, reasoning_effort or both"
		if research {
			expected = "model, reasoning_effort or tools"
		}
		return fmt.Errorf("%s declares nothing; set %s", path, expected)
	}
	return nil
}

func validateAuthoringFallback(node *yaml.Node, entryPath, model string) error {
	path := entryPath + ".fallback"
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("%s must be a mapping", path)
	}
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i].Value
		if key != "model" && key != "reasoning_effort" {
			return fmt.Errorf("%s.%s is not a supported key; expected model or reasoning_effort", path, key)
		}
	}
	fallbackModel := ""
	for i := 0; i < len(node.Content); i += 2 {
		key, value := node.Content[i].Value, node.Content[i+1]
		raw := authoringNodeValue(value)
		switch key {
		case "model":
			if value.Kind != yaml.ScalarNode || value.Tag != "!!str" || !authoringModelName.MatchString(raw) {
				return fmt.Errorf("%s.model %q is not a model name; use letters, digits and . _ : / @ [ ] - without spaces", path, raw)
			}
			fallbackModel = raw
		case "reasoning_effort":
			if value.Kind != yaml.ScalarNode || value.Tag != "!!str" || !slices.Contains([]string{"low", "medium", "high", "xhigh", "max"}, raw) {
				return fmt.Errorf("%s.reasoning_effort %q is not supported for claude; expected low, medium, high, xhigh or max", path, raw)
			}
		}
	}
	if fallbackModel == "" {
		return fmt.Errorf("%s declares no model; set model", path)
	}
	if model == "" {
		return fmt.Errorf("%s needs %s.model", path, entryPath)
	}
	if strings.EqualFold(fallbackModel, model) {
		return fmt.Errorf("%s.model %q is the declared model; a fallback must name another model", path, fallbackModel)
	}
	return nil
}

func invalidAuthoringTool(value string) error {
	return fmt.Errorf("skills.authoring.claude.research.tools names %q; expected WebSearch, WebFetch or a name that begins with mcp__", value)
}

func (h *AuthoringHarness) empty() bool {
	return h != nil && h.Idea == nil && h.PRD == nil && h.TechSpec == nil && h.Tasks == nil && h.Research == nil
}
