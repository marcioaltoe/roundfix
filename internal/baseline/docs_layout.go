package baseline

import (
	"regexp"
	"slices"
	"strings"
)

// DocumentationDeparture is one advisory departure from a documentation rule.
type DocumentationDeparture struct {
	Code, Path string
	Line       int
	Detail     string
}

// DocsLayoutFacts contains only repository-relative paths and observed content.
type DocsLayoutFacts struct {
	RootDocuments      []string
	DocsEntries        []string
	SpecRoot, History  string
	GlossaryFiles      []string
	ReadmePresent      bool
	SpecificRepository []byte
}

var rootReferenceDocument = regexp.MustCompile(`^[A-Z0-9_-]+\.md$`)

// DocsLayoutDepartures compares observed paths with the Baseline's named layout.
func DocsLayoutDepartures(facts DocsLayoutFacts) []DocumentationDeparture {
	var result []DocumentationDeparture
	if !facts.ReadmePresent {
		result = append(result, DocumentationDeparture{Code: "DR-DOCS-REQUIRED-MISSING", Path: "README.md", Detail: "README.md is missing"})
	}
	if len(facts.GlossaryFiles) == 0 {
		result = append(result, DocumentationDeparture{Code: "DR-DOCS-REQUIRED-MISSING", Detail: "the glossary is missing"})
	}
	folders := []string{"_inbox", "adr", "agents", "backlog", "design", "findings", "handoffs", "history", "lessons", "references", "user-guide", "workflow"}
	for _, root := range []string{facts.SpecRoot, facts.History} {
		if rest, ok := strings.CutPrefix(root, "docs/"); ok {
			first, _, _ := strings.Cut(rest, "/")
			folders = append(folders, first)
		}
	}
	own := string(facts.SpecificRepository)
	outside := func(path, detail string) {
		if !strings.Contains(own, path) {
			result = append(result, DocumentationDeparture{Code: "DR-DOCS-OUTSIDE-LAYOUT", Path: path, Detail: detail})
		}
	}
	for _, path := range facts.DocsEntries {
		if path == "docs/README.md" {
			continue
		}
		if strings.HasSuffix(path, "/") {
			if !slices.Contains(folders, strings.TrimSuffix(strings.TrimPrefix(path, "docs/"), "/")) {
				outside(path, "a folder")
			}
		} else {
			outside(path, "a file directly under docs/")
		}
	}
	named := []string{"README", "AGENTS", "CLAUDE", "CONTEXT", "CONTEXT-MAP", "GLOSSARY", "GLOSSARY-MAP", "CHANGELOG", "CONTRIBUTING", "LICENSE", "LICENCE", "SECURITY", "CODE_OF_CONDUCT", "SUPPORT", "DESIGN"}
	for _, path := range facts.RootDocuments {
		if rootReferenceDocument.MatchString(path) && !slices.Contains(named, strings.TrimSuffix(path, ".md")) {
			outside(path, "a reference document at the root")
		}
	}
	return result
}
