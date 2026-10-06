package speccheck

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"roundfix/internal/spec"
)

const (
	// CodeGlossaryUndeclared identifies a missing declaration or uncovered bold term.
	CodeGlossaryUndeclared = "SC-GLOSSARY-UNDECLARED"
	// CodeGlossaryUnplanned identifies a declared term with no glossary-writing Task.
	CodeGlossaryUnplanned = "SC-GLOSSARY-UNPLANNED"
	// CodeGlossaryMissing identifies a term absent after all binding Tasks completed.
	CodeGlossaryMissing = "SC-GLOSSARY-MISSING"
	// GlossaryGuidePath is the guide whose adding commit starts the glossary horizon.
	GlossaryGuidePath = ".agents/skills/write-prd/references/glossary.md"
)

type glossaryEntry struct {
	Kind   string
	Term   string
	Reason string
	Where  Location
}

var glossaryDeclarationPattern = regexp.MustCompile(`^- (adds|changes|not a term): \*\*([^*]+)\*\*(?: — (.+))?$`)
var glossaryBoldPattern = regexp.MustCompile(`\*\*([^*]+)\*\*`)
var glossaryMapLinkPattern = regexp.MustCompile(`\[[^\]]*\]\(\s*(?:<([^>]+)>|([^\s)]+))(?:\s+"[^"\n]*")?\s*\)|(?m)^\s*\[[^\]]+\]:\s*(?:<([^>]+)>|([^\s]+))`)
var glossaryDefinitionPattern = regexp.MustCompile(`^\*\*([^*]+)\*\*:`)

// GlossaryFindings reads only the glossary detector and the optional Task Graph.
func GlossaryFindings(specsRoot, repoRoot, slug string) ([]Finding, error) {
	if strings.TrimSpace(slug) == "" || filepath.Base(slug) != slug || slug == "." {
		return nil, fmt.Errorf("invalid Spec slug %q", slug)
	}
	dir := filepath.Join(specsRoot, slug)
	graph, _, err := loadOptionalTaskGraph(specsRoot, slug, dir)
	if err != nil {
		return nil, err
	}
	result := Result{Slug: slug}
	err = detectGlossary(&result, repoRoot, []string{filepath.Join(dir, "_prd.md"), filepath.Join(dir, "_techspec.md")}, StageAll, graph)
	return result.Findings, err
}

func glossaryNormalize(term string) string {
	return strings.ToLower(strings.Join(strings.Fields(term), " "))
}

func glossaryHorizon(repoRoot, prdPath string) contractHorizon {
	if _, ok, err := receiptFilePath(repoRoot, GlossaryGuidePath); err != nil {
		return contractHorizon{held: true}
	} else if !ok {
		return contractHorizon{missing: "a PRD committed at or after " + GlossaryGuidePath}
	}
	commit, readable := prdAddingCommit(repoRoot, prdPath)
	if !readable {
		return contractHorizon{held: true}
	}
	output, err := adrHorizonGitOutput(repoRoot, "log", "--diff-filter=A", "--format=%H", "--", GlossaryGuidePath)
	if err != nil {
		return contractHorizon{held: true}
	}
	commits := strings.Fields(string(output))
	if len(commits) == 0 {
		return contractHorizon{missing: "a PRD committed at or after " + GlossaryGuidePath}
	}
	if commit == "" {
		return contractHorizon{held: true}
	}
	_, err = adrHorizonGitOutput(repoRoot, "merge-base", "--is-ancestor", commits[len(commits)-1], commit)
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
		return contractHorizon{missing: "a PRD committed at or after " + GlossaryGuidePath}
	}
	return contractHorizon{held: true}
}

func glossaryFiles(repoRoot string) (map[string]bool, map[string]bool, error) {
	files := map[string]bool{}
	terms := map[string]bool{}
	for _, name := range []string{"CONTEXT.md", "GLOSSARY.md", "CONTEXT-MAP.md", "GLOSSARY-MAP.md"} {
		_, ok, err := receiptFilePath(repoRoot, name)
		if err != nil {
			return nil, nil, err
		}
		if !ok {
			continue
		}
		if name == "CONTEXT.md" || name == "GLOSSARY.md" {
			files[name] = true
			continue
		}
		content, err := os.ReadFile(filepath.Join(repoRoot, name))
		if err != nil {
			return nil, nil, fmt.Errorf("read glossary map: %w", err)
		}
		for _, match := range glossaryMapLinkPattern.FindAllSubmatch(content, -1) {
			target := ""
			for _, group := range match[1:] {
				if len(group) > 0 {
					target = string(group)
					break
				}
			}
			target = strings.Split(strings.Split(target, "#")[0], "?")[0]
			if strings.Contains(target, ":") || filepath.IsAbs(target) {
				continue
			}
			target = filepath.ToSlash(filepath.Clean(filepath.FromSlash(target)))
			base := filepath.Base(target)
			if base != "CONTEXT.md" && base != "GLOSSARY.md" {
				continue
			}
			_, ok, err := receiptFilePath(repoRoot, target)
			if err != nil {
				return nil, nil, err
			}
			if ok {
				files[target] = true
			}
		}
	}
	for file := range files {
		content, err := os.ReadFile(filepath.Join(repoRoot, filepath.FromSlash(file)))
		if err != nil {
			return nil, nil, fmt.Errorf("read glossary %q: %w", file, err)
		}
		for _, line := range strings.Split(string(content), "\n") {
			if match := glossaryDefinitionPattern.FindStringSubmatch(line); match != nil {
				terms[glossaryNormalize(match[1])] = true
			}
		}
	}
	return files, terms, nil
}

func glossaryFinding(result *Result, code, summary, fix string, where Location) {
	result.Findings = append(result.Findings, Finding{Code: code, Severity: SeverityError, Summary: summary, Where: []Location{where}, Fix: fix})
}

func parseGlossary(result *Result, content, display string) ([]glossaryEntry, bool, map[int]bool) {
	var entries []glossaryEntry
	declared, section := false, false
	excluded := map[int]bool{}
	var none []int
	count := 0
	for i, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, "## ") {
			section = line == "## Glossary"
			declared = declared || section
			continue
		}
		if !section {
			continue
		}
		excluded[i] = true
		if strings.TrimSpace(line) == "" {
			continue
		}
		count++
		if line == "None." {
			none = append(none, i+1)
			continue
		}
		match := glossaryDeclarationPattern.FindStringSubmatch(line)
		if match == nil || strings.TrimSpace(match[2]) == "" || (match[1] == "not a term" && strings.TrimSpace(match[3]) == "") {
			glossaryFinding(result, CodeGlossaryUndeclared, display+" has a malformed Glossary Declaration", "Use None. alone, or - adds: **term**, - changes: **term**, or - not a term: **phrase** — reason.", Location{Path: display, Line: i + 1})
			continue
		}
		entries = append(entries, glossaryEntry{Kind: match[1], Term: match[2], Reason: match[3], Where: Location{Path: display, Line: i + 1}})
	}
	if count > 1 {
		for _, line := range none {
			glossaryFinding(result, CodeGlossaryUndeclared, display+" mixes None. with other Glossary entries", "Use None. only as the section's single entry.", Location{Path: display, Line: line})
		}
	}
	return entries, declared, excluded
}

func glossaryBoldCandidate(phrase string) bool {
	words := strings.Fields(phrase)
	if len(words) < 2 || len(words) > 5 || strings.ContainsAny(phrase, "`()") || strings.ContainsAny(phrase[len(phrase)-1:], ".:;,?!") {
		return false
	}
	for _, r := range phrase {
		if unicode.IsDigit(r) {
			return false
		}
	}
	connectors := " a an and at by for from in of on or per the to with "
	for i, word := range words {
		first := []rune(word)[0]
		if unicode.IsUpper(first) {
			continue
		}
		if i > 0 && strings.Contains(connectors, " "+word+" ") {
			continue
		}
		return false
	}
	return true
}

func detectGlossary(result *Result, repoRoot string, artifacts []string, stage Stage, graph *spec.Graph) error {
	files, terms, err := glossaryFiles(repoRoot)
	if err != nil {
		return err
	}
	skipAll := func(missing string) {
		for _, code := range []string{CodeGlossaryUndeclared, CodeGlossaryUnplanned, CodeGlossaryMissing} {
			addSkip(result, code, missing)
		}
	}
	if len(files) == 0 {
		skipAll("CONTEXT.md or GLOSSARY.md")
		return nil
	}
	type artifact struct {
		content, display string
		excluded         map[int]bool
	}
	var documents []artifact
	var entries []glossaryEntry
	declared := false
	for _, path := range artifacts {
		content, err := os.ReadFile(path)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return fmt.Errorf("read glossary declaration %q: %w", path, err)
		}
		display := artifactDisplayPath(repoRoot, path)
		parsed, present, excluded := parseGlossary(result, string(content), display)
		entries = append(entries, parsed...)
		declared = declared || present
		documents = append(documents, artifact{string(content), display, excluded})
	}
	if len(documents) == 0 {
		skipAll(artifactDisplayPath(repoRoot, artifacts[0]))
		return nil
	}
	if !declared {
		horizon := glossaryHorizon(repoRoot, artifacts[0])
		if !horizon.held {
			skipAll(horizon.missing)
			return nil
		}
		glossaryFinding(result, CodeGlossaryUndeclared, documents[0].display+" omits a Glossary Declaration", "Add ## Glossary with declared terms or None.", Location{Path: documents[0].display, Line: 1})
	}
	covered := map[string]bool{}
	for term := range terms {
		covered[term] = true
		covered[term+"s"] = true
		covered[term+"es"] = true
	}
	for _, entry := range entries {
		covered[glossaryNormalize(entry.Term)] = true
	}
	for _, doc := range documents {
		seen := map[string]bool{}
		fence := ""
		fenceLength := 0
		for i, line := range strings.Split(doc.content, "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
				marker := trimmed[:1]
				length := len(trimmed) - len(strings.TrimLeft(trimmed, marker))
				if fence == "" {
					fence = marker
					fenceLength = length
				} else if marker == fence && length >= fenceLength && strings.TrimSpace(trimmed[length:]) == "" {
					fence = ""
				}
				continue
			}
			if fence != "" || doc.excluded[i] {
				continue
			}
			for _, match := range glossaryBoldPattern.FindAllStringSubmatch(line, -1) {
				phrase := match[1]
				key := glossaryNormalize(phrase)
				if !glossaryBoldCandidate(phrase) || covered[key] || seen[key] {
					continue
				}
				seen[key] = true
				glossaryFinding(result, CodeGlossaryUndeclared, doc.display+" bolds undeclared term "+fmt.Sprintf("%q", phrase), "Define the term in the glossary or cover it in ## Glossary, including - not a term: **phrase** — reason.", Location{Path: doc.display, Line: i + 1})
			}
		}
	}
	if stage == StagePRD || stage == StageTechSpec {
		return nil
	}
	if graph == nil {
		for _, code := range []string{CodeGlossaryUnplanned, CodeGlossaryMissing} {
			addSkip(result, code, artifactDisplayPath(repoRoot, filepath.Join(filepath.Dir(artifacts[0]), "_tasks.md")))
		}
		return nil
	}
	for _, entry := range entries {
		if entry.Kind == "not a term" {
			continue
		}
		defined := terms[glossaryNormalize(entry.Term)]
		if entry.Kind == "changes" && !defined {
			glossaryFinding(result, CodeGlossaryUnplanned, "Changed glossary term "+fmt.Sprintf("%q", entry.Term)+" is not defined", "Declare the term as added instead of changed.", entry.Where)
		}
		bindings := 0
		completed := true
		for _, task := range graph.Tasks {
			if task.ID == graph.QATaskID || task.Type == spec.TaskTypeQA {
				continue
			}
			writes := false
			for _, ref := range task.Context {
				if (ref.Kind == spec.ContextKindInterface || ref.Kind == spec.ContextKindCreates) && files[ref.Path] {
					writes = true
				}
			}
			if !writes {
				continue
			}
			taskPath := filepath.Join(filepath.Dir(filepath.Dir(artifacts[0])), filepath.FromSlash(task.File))
			content, err := os.ReadFile(taskPath)
			if err != nil {
				return fmt.Errorf("read glossary Task %q: %w", taskPath, err)
			}
			verification := strings.Join(markdownSectionLines(content, "Verification"), "\n")
			if !strings.Contains(verification, "**"+entry.Term+"**") {
				continue
			}
			bindings++
			completed = completed && task.Status == spec.StatusCompleted
		}
		if bindings == 0 {
			glossaryFinding(result, CodeGlossaryUnplanned, "Glossary term "+fmt.Sprintf("%q", entry.Term)+" has no binding Task", "Add a non-QA Task declaring the glossary under interface: or creates: and naming **"+entry.Term+"** in ## Verification.", entry.Where)
		} else if completed && !defined {
			glossaryFinding(result, CodeGlossaryMissing, "Glossary term "+fmt.Sprintf("%q", entry.Term)+" is missing after its binding Tasks completed", "Define the declared term in the glossary.", entry.Where)
		}
	}
	return nil
}
