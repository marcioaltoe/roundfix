package spec

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// PlanningBinding returns the Binding Text a Task names in one planning file.
// A whole binding cannot prove an amendment safe; callers retain the usual
// moved-input refusal in that case. Explicit Context paths always bind whole.
func PlanningBinding(file string, planning []byte, taskID string, task []byte) (binding string, whole bool, err error) {
	refs, err := parseTaskContextRefs(task)
	if err != nil {
		return "", true, err
	}
	for _, ref := range refs {
		if filepath.Clean(ref.Path) == filepath.Clean(file) {
			return string(planning), true, nil
		}
	}
	// Git text inputs can be checked out with CRLF while their blobs use LF.
	text := strings.ReplaceAll(string(planning), "\r\n", "\n")
	switch filepath.Base(file) {
	case "_tasks.md":
		return planningTaskBinding([]byte(text), taskID)
	case "_prd.md", "_techspec.md":
		return planningSectionBinding(filepath.Base(file), text, task)
	default:
		return text, true, nil
	}
}

func planningTaskBinding(planning []byte, taskID string) (string, bool, error) {
	front, body, err := splitFrontmatter(planning)
	if err != nil {
		return "", true, fmt.Errorf("read planning graph: %w", err)
	}
	var manifest struct {
		Graph struct {
			Nodes []manifestNode `yaml:"nodes"`
		} `yaml:"graph"`
	}
	if err := yaml.Unmarshal(front, &manifest); err != nil {
		return "", true, fmt.Errorf("read planning graph: %w", err)
	}
	var node *manifestNode
	for index := range manifest.Graph.Nodes {
		if manifest.Graph.Nodes[index].ID == taskID {
			node = &manifest.Graph.Nodes[index]
			break
		}
	}
	if node == nil {
		return string(planning), true, nil
	}
	code := archiveMarkdownCode(body)
	inTable := false
	offset := 0
	for _, line := range strings.SplitAfter(string(body), "\n") {
		start := offset
		offset += len(line)
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			inTable = false
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " \t"))
		if start+indent < len(code) && code[start+indent] {
			continue
		}
		cells := markdownTableCells(line)
		if taskTypeProjectionHeader(cells) {
			inTable = true
			continue
		}
		if len(cells) == 0 {
			inTable = false
			continue
		}
		if inTable && cells[0] == taskID {
			encoded, err := yaml.Marshal(node)
			if err != nil {
				return "", true, fmt.Errorf("encode planning Task node: %w", err)
			}
			return string(encoded) + line, false, nil
		}
	}
	return string(planning), true, nil
}

type planningHeading struct {
	start, depth int
	name         string
}

// planningHeadings keeps offsets into the original text, including headings
// deeper than four as section boundaries, but ignores fenced examples.
func planningHeadings(text string) []planningHeading {
	code := archiveMarkdownCode([]byte(text))
	var headings []planningHeading
	offset := 0
	for _, line := range strings.SplitAfter(text, "\n") {
		start := offset
		offset += len(line)
		trimmed := strings.TrimLeft(line, " \t")
		indent := len(line) - len(trimmed)
		if indent > 3 || start+indent >= len(code) || code[start+indent] {
			continue
		}
		if depth, name, ok := markdownHeading(strings.TrimSpace(trimmed)); ok {
			headings = append(headings, planningHeading{start, depth, name})
		}
	}
	return headings
}

var planningUnitNumbers = regexp.MustCompile(`(?i)\s+[0-9][0-9\s,–—-]*(?:and\s+[0-9][0-9\s,–—-]*)*$`)

func planningSectionBinding(file, text string, task []byte) (string, bool, error) {
	taskText := string(task)
	var units []string
	taskHeadings := planningHeadings(taskText)
	taskCode := archiveMarkdownCode(task)
	for i, heading := range taskHeadings {
		if heading.depth != 2 || heading.name != "References" {
			continue
		}
		end := len(taskText)
		for _, next := range taskHeadings[i+1:] {
			if next.depth <= 2 {
				end = next.start
				break
			}
		}
		prefix := "- `" + file + "` →"
		var bullet string
		offset := heading.start
		for _, line := range strings.SplitAfter(taskText[heading.start:end], "\n") {
			start := offset
			offset += len(line)
			trimmed := strings.TrimSpace(line)
			indent := len(line) - len(strings.TrimLeft(line, " \t"))
			if start+indent < len(taskCode) && taskCode[start+indent] {
				continue
			}
			if strings.HasPrefix(trimmed, "- ") {
				if bullet != "" {
					units = append(units, strings.Split(bullet, ";")...)
				}
				bullet = ""
				if strings.HasPrefix(trimmed, prefix) {
					bullet = strings.TrimSpace(strings.TrimPrefix(trimmed, prefix))
				}
			} else if bullet != "" && trimmed != "" {
				bullet += " " + trimmed
			}
		}
		if bullet != "" {
			units = append(units, strings.Split(bullet, ";")...)
		}
	}
	if len(units) == 0 {
		return text, true, nil
	}
	headings := planningHeadings(text)
	selected := make(map[int]bool)
	for _, unit := range units {
		name, _, _ := strings.Cut(unit, ":")
		name = strings.ToLower(strings.TrimSpace(planningUnitNumbers.ReplaceAllString(strings.TrimSpace(name), "")))
		found := false
		for i, heading := range headings {
			if heading.depth < 2 || heading.depth > 4 {
				continue
			}
			other := strings.ToLower(heading.name)
			if name == other || planningPlural(name) == other || name == planningPlural(other) {
				selected[i], found = true, true
			}
		}
		if !found {
			return text, true, nil
		}
	}
	var binding strings.Builder
	for i, heading := range headings {
		if !selected[i] {
			continue
		}
		end := len(text)
		for _, next := range headings[i+1:] {
			if next.depth <= heading.depth {
				end = next.start
				break
			}
		}
		binding.WriteString(text[heading.start:end])
	}
	return binding.String(), false, nil
}

func planningPlural(name string) string {
	if strings.HasSuffix(name, "y") {
		return strings.TrimSuffix(name, "y") + "ies"
	}
	return name + "s"
}
