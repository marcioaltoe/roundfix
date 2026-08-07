package speccheck

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// CodeBacklogUnmoved identifies a Backlog Entry that declares itself promoted
// to a Spec while still living in docs/backlog/. The Backlog Operational
// Contract requires a promoted entry to move into that Spec's references/,
// where the adoption index records it.
const CodeBacklogUnmoved = "SC-BACKLOG-UNMOVED"

const backlogPath = "docs/backlog"

// backlogFrontmatter carries only the two declared values this rule reads.
// Nothing here infers intent: an entry is a defect when it says it is promoted
// and its own location contradicts that.
type backlogFrontmatter struct {
	status    findingFrontmatterValue
	hasStatus bool
	spec      findingFrontmatterValue
	hasSpec   bool
}

type backlogDocument struct {
	displayPath string
	frontmatter backlogFrontmatter
}

// detectBacklogPromotion reports promoted Backlog Entries that never moved.
// It is presence-aware per ADR-0094: a repository with no backlog directory,
// or one holding no entries, skips instead of failing.
func detectBacklogPromotion(result *Result, repoRoot string) error {
	documents, present, err := readBacklogDocuments(repoRoot, backlogPath)
	if err != nil {
		return err
	}
	if !present || len(documents) == 0 {
		addSkip(result, CodeBacklogUnmoved, backlogPath)
		return nil
	}

	activeSpecs, err := repositoryDirectoryNames(filepath.Join(filepath.Clean(repoRoot), "docs", "specs"), true)
	if err != nil {
		return err
	}
	archivedSpecs, err := repositoryDirectoryNames(filepath.Join(filepath.Clean(repoRoot), "docs", "specs", "_archived"), false)
	if err != nil {
		return err
	}

	for _, document := range documents {
		if !document.frontmatter.hasStatus || document.frontmatter.status.value != "promoted" {
			continue
		}
		if !document.frontmatter.hasSpec || strings.TrimSpace(document.frontmatter.spec.value) == "" ||
			document.frontmatter.spec.value == "null" {
			result.Findings = append(result.Findings, Finding{
				Code:     CodeBacklogUnmoved,
				Severity: SeverityError,
				Summary:  document.displayPath + " is promoted without naming its Spec",
				Where:    []Location{{Path: document.displayPath, Line: document.frontmatter.status.line}},
				Fix:      "Set spec to the owning Spec slug in " + document.displayPath + ", then move the entry into that Spec's references/ and index it.",
			})
			continue
		}
		slug := document.frontmatter.spec.value
		destination := "docs/specs/" + slug + "/references/"
		if !activeSpecs[slug] && !archivedSpecs[slug] {
			result.Findings = append(result.Findings, Finding{
				Code:     CodeBacklogUnmoved,
				Severity: SeverityError,
				Summary:  document.displayPath + " names unresolvable Spec " + strconv.Quote(slug),
				Where:    []Location{{Path: document.displayPath, Line: document.frontmatter.spec.line}},
				Fix:      "Correct the spec slug in " + document.displayPath + " to a Spec directory under docs/specs/ or docs/specs/_archived/.",
			})
			continue
		}
		if archivedSpecs[slug] {
			destination = "docs/specs/_archived/" + slug + "/references/"
		}
		result.Findings = append(result.Findings, Finding{
			Code:     CodeBacklogUnmoved,
			Severity: SeverityError,
			Summary:  document.displayPath + " is promoted to " + slug + " but still lives in " + backlogPath,
			Where:    []Location{{Path: document.displayPath, Line: document.frontmatter.status.line}},
			Fix:      "Move " + document.displayPath + " to " + destination + " and add its row to that Spec's references/_index.md with type backlog.",
		})
	}
	return nil
}

func readBacklogDocuments(repoRoot, relativeDir string) ([]backlogDocument, bool, error) {
	directory := filepath.Join(filepath.Clean(repoRoot), filepath.FromSlash(relativeDir))
	entries, err := os.ReadDir(directory)
	if errors.Is(err, os.ErrNotExist) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("read backlog directory %q: %w", directory, err)
	}

	documents := make([]backlogDocument, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") || strings.HasPrefix(entry.Name(), "_") {
			continue
		}
		path := filepath.Join(directory, entry.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, false, fmt.Errorf("read Backlog Entry %q: %w", path, err)
		}
		frontmatter, err := parseBacklogFrontmatter(content)
		if err != nil {
			return nil, false, fmt.Errorf("parse Backlog Entry %q: %w", path, err)
		}
		documents = append(documents, backlogDocument{
			displayPath: relativeDir + "/" + entry.Name(),
			frontmatter: frontmatter,
		})
	}
	return documents, true, nil
}

func parseBacklogFrontmatter(content []byte) (backlogFrontmatter, error) {
	const opening = "---\n"
	text := string(content)
	if !strings.HasPrefix(text, opening) {
		return backlogFrontmatter{}, nil
	}
	rest := text[len(opening):]
	end := strings.Index(rest, "\n---")
	if end < 0 {
		return backlogFrontmatter{}, nil
	}

	var document yaml.Node
	if err := yaml.Unmarshal([]byte(rest[:end]), &document); err != nil {
		return backlogFrontmatter{}, fmt.Errorf("parse YAML frontmatter: %w", err)
	}
	if len(document.Content) == 0 || len(document.Content[0].Content) == 0 {
		return backlogFrontmatter{}, nil
	}
	mapping := document.Content[0]
	if mapping.Kind != yaml.MappingNode {
		return backlogFrontmatter{}, errors.New("YAML frontmatter must be a mapping")
	}

	var frontmatter backlogFrontmatter
	for index := 0; index+1 < len(mapping.Content); index += 2 {
		key := mapping.Content[index]
		value := mapping.Content[index+1]
		switch key.Value {
		case "status":
			decoded, err := findingScalar(value)
			if err != nil {
				return backlogFrontmatter{}, fmt.Errorf("read status: %w", err)
			}
			frontmatter.status = findingFrontmatterValue{value: decoded, line: value.Line}
			frontmatter.hasStatus = true
		case "spec":
			decoded, err := findingScalar(value)
			if err != nil {
				return backlogFrontmatter{}, fmt.Errorf("read spec: %w", err)
			}
			frontmatter.spec = findingFrontmatterValue{value: decoded, line: value.Line}
			frontmatter.hasSpec = decoded != ""
		}
	}
	return frontmatter, nil
}
