package speccheck

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const (
	// CodeVocabularyUndocumented identifies emitted vocabulary absent from its declared documentation.
	CodeVocabularyUndocumented = "SC-VOCABULARY-UNDOCUMENTED"

	vocabularyDeclarationShape = "## Vocabulary Contract\n\n" +
		"- emits: `<repository-relative path>`\n" +
		"  pattern: `<RE2>`\n" +
		"  documented-in: `<repository-relative path>`"
)

type vocabularyContract struct {
	EmitsPath      string
	Pattern        string
	DocumentedIn   string
	Line           int
	PatternLine    int
	DocumentedLine int
}

func detectVocabularyContract(result *Result, repoRoot, techSpecPath string, techSpecPresent bool) error {
	techSpecDisplayPath := artifactDisplayPath(repoRoot, techSpecPath)
	if !techSpecPresent {
		addSkip(result, CodeVocabularyUndocumented, techSpecDisplayPath)
		return nil
	}

	content, err := os.ReadFile(techSpecPath)
	if err != nil {
		return fmt.Errorf("read vocabulary contract %q: %w", techSpecPath, err)
	}
	contracts, declared := parseVocabularyContracts(content)
	if !declared || len(contracts) == 0 {
		addSkip(result, CodeVocabularyUndocumented, techSpecDisplayPath+" Vocabulary Contract")
		return nil
	}

	for _, contract := range contracts {
		detectVocabularyEntry(result, repoRoot, techSpecDisplayPath, contract)
	}
	return nil
}

func parseVocabularyContracts(content []byte) ([]vocabularyContract, bool) {
	lines := strings.Split(string(content), "\n")
	inSection := false
	declared := false
	var contracts []vocabularyContract
	var current *vocabularyContract

	flush := func() {
		if current == nil {
			return
		}
		contracts = append(contracts, *current)
		current = nil
	}

	for index, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## ") {
			if inSection {
				break
			}
			inSection = strings.TrimSpace(strings.TrimPrefix(trimmed, "## ")) == "Vocabulary Contract"
			declared = declared || inSection
			continue
		}
		if !inSection {
			continue
		}

		if value, ok := vocabularyField(trimmed, "- emits:"); ok {
			flush()
			current = &vocabularyContract{EmitsPath: value, Line: index + 1}
			continue
		}
		if current == nil {
			continue
		}
		if value, ok := vocabularyField(trimmed, "pattern:"); ok {
			current.Pattern = value
			current.PatternLine = index + 1
			continue
		}
		if value, ok := vocabularyField(trimmed, "documented-in:"); ok {
			current.DocumentedIn = value
			current.DocumentedLine = index + 1
		}
	}
	flush()
	return contracts, declared
}

func vocabularyField(line, prefix string) (string, bool) {
	if !strings.HasPrefix(line, prefix) {
		return "", false
	}
	value := strings.TrimSpace(strings.TrimPrefix(line, prefix))
	if len(value) >= 2 && strings.HasPrefix(value, "`") && strings.HasSuffix(value, "`") {
		value = strings.TrimSuffix(strings.TrimPrefix(value, "`"), "`")
	}
	return strings.TrimSpace(value), true
}

func detectVocabularyEntry(result *Result, repoRoot, techSpecDisplayPath string, contract vocabularyContract) {
	declarationLocation := Location{Path: techSpecDisplayPath, Line: contract.Line}
	if contract.EmitsPath == "" || contract.Pattern == "" || contract.DocumentedIn == "" {
		result.Findings = append(result.Findings, Finding{
			Code:     CodeVocabularyUndocumented,
			Severity: SeverityError,
			Summary:  techSpecDisplayPath + " has an incomplete Vocabulary Contract declaration",
			Where:    []Location{declarationLocation, {Path: techSpecDisplayPath, Line: contract.Line}},
			Fix:      vocabularyFix("Complete the vocabulary declaration.", techSpecDisplayPath),
		})
		return
	}

	pattern, err := regexp.Compile(contract.Pattern)
	if err != nil {
		patternLine := contract.PatternLine
		if patternLine == 0 {
			patternLine = contract.Line
		}
		result.Findings = append(result.Findings, Finding{
			Code:     CodeVocabularyUndocumented,
			Severity: SeverityError,
			Summary:  techSpecDisplayPath + " declares invalid RE2 pattern " + fmt.Sprintf("%q", contract.Pattern),
			Where: []Location{
				{Path: techSpecDisplayPath, Line: patternLine},
				{Path: contract.EmitsPath, Line: 1},
			},
			Fix: vocabularyFix("Replace the pattern with a valid RE2 expression.", techSpecDisplayPath),
		})
		return
	}

	emitting, ok := readVocabularyPath(repoRoot, contract.EmitsPath)
	if !ok {
		result.Findings = append(result.Findings, Finding{
			Code:     CodeVocabularyUndocumented,
			Severity: SeverityError,
			Summary:  techSpecDisplayPath + " declares unreadable emitting path " + contract.EmitsPath,
			Where:    []Location{declarationLocation, {Path: contract.EmitsPath, Line: 1}},
			Fix:      vocabularyFix("Point emits to a readable repository file.", techSpecDisplayPath),
		})
		return
	}

	documentation, ok := readVocabularyPath(repoRoot, contract.DocumentedIn)
	if !ok {
		documentedLine := contract.DocumentedLine
		if documentedLine == 0 {
			documentedLine = contract.Line
		}
		result.Findings = append(result.Findings, Finding{
			Code:     CodeVocabularyUndocumented,
			Severity: SeverityError,
			Summary:  techSpecDisplayPath + " declares unreadable documenting path " + contract.DocumentedIn,
			Where: []Location{
				{Path: techSpecDisplayPath, Line: documentedLine},
				{Path: contract.DocumentedIn, Line: 1},
			},
			Fix: vocabularyFix("Point documented-in to a readable repository file.", techSpecDisplayPath),
		})
		return
	}

	seen := make(map[string]bool)
	for _, match := range pattern.FindAllIndex(emitting, -1) {
		token := string(emitting[match[0]:match[1]])
		if seen[token] {
			continue
		}
		seen[token] = true
		if bytes.Contains(documentation, []byte(token)) {
			continue
		}
		line := 1 + bytes.Count(emitting[:match[0]], []byte("\n"))
		result.Findings = append(result.Findings, Finding{
			Code:     CodeVocabularyUndocumented,
			Severity: SeverityError,
			Summary:  contract.EmitsPath + " emits undocumented token " + fmt.Sprintf("%q", token) + " absent from " + contract.DocumentedIn,
			Where: []Location{
				{Path: contract.EmitsPath, Line: line},
				{Path: contract.DocumentedIn, Line: 1},
			},
			Fix: vocabularyFix("Document "+fmt.Sprintf("%q", token)+" in "+contract.DocumentedIn+" or narrow the declared pattern.", techSpecDisplayPath),
		})
	}
}

func readVocabularyPath(repoRoot, relative string) ([]byte, bool) {
	path, ok := resolveRepositoryPath(repoRoot, relative)
	if !ok {
		return nil, false
	}
	content, err := os.ReadFile(path)
	return content, err == nil
}

func vocabularyFix(action, techSpecDisplayPath string) string {
	return action + " Use this exact declaration shape in " + techSpecDisplayPath + ":\n" + vocabularyDeclarationShape
}
