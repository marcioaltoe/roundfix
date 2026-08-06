package speccheck

import (
	"regexp"
	"strings"

	"roundfix/internal/spec"
)

const (
	// CodeRequirementContradictory identifies declared MUST and MUST NOT clauses over one named subject.
	CodeRequirementContradictory = "SC-REQUIREMENT-CONTRADICTORY"
	// CodeRehearsalUndeclared identifies a gate rehearsal with no complete case declarations.
	CodeRehearsalUndeclared = "SC-REHEARSAL-UNDECLARED"
)

var (
	modalMentionPattern     = regexp.MustCompile(`(?i)\bMUST\s+and\s+MUST\s+NOT\b`)
	requirementModalPattern = regexp.MustCompile(`(?i)(?:^|[,;]\s+|\band\s+)MUST(\s+NOT)?\s+`)
	requirementWordPattern  = regexp.MustCompile(`[a-z0-9]+`)
	rehearsalWordPattern    = regexp.MustCompile(`(?i)\brehears[a-z]*\b`)
	proveWordPattern        = regexp.MustCompile(`(?i)\bprov(?:e|es|ed|ing)\b`)
	gateWordPattern         = regexp.MustCompile(`(?i)\bgates?\b`)
)

type declaredRequirementClause struct {
	Requirement int
	Forbidden   bool
	Words       []string
}

// ContradictoryRequirements reports one requirement forbidding a state that
// another requirement needs, decided from declared MUST/MUST NOT pairs over
// the same named subject.
func ContradictoryRequirements(task spec.Task) (Finding, bool) {
	clauses := declaredRequirementClauses(task.Requirements)
	for _, needed := range clauses {
		if needed.Forbidden {
			continue
		}
		for _, forbidden := range clauses {
			if !forbidden.Forbidden {
				continue
			}
			subject, ok := commonNamedSubject(needed.Words, forbidden.Words)
			if !ok {
				continue
			}
			return Finding{
				Code:     CodeRequirementContradictory,
				Severity: SeverityError,
				Summary:  task.File + " declares MUST and MUST NOT requirements over named subject " + subject,
				Where: []Location{
					{Path: task.File, Line: needed.Requirement + 1},
					{Path: task.File, Line: forbidden.Requirement + 1},
				},
				Fix: "Rewrite one requirement so the declared state of " + subject + " is consistent.",
			}, true
		}
	}
	return Finding{}, false
}

// UndeclaredRehearsal reports a Task whose title declares a gate rehearsal but
// whose Rehearsal Cases section has no complete case and observation entries.
func UndeclaredRehearsal(task spec.Task) (Finding, bool) {
	if !declaresGateRehearsal(task.Title) {
		return Finding{}, false
	}
	for _, entry := range task.RehearsalCases {
		if !completeRehearsalCase(entry) {
			return undeclaredRehearsalFinding(task), true
		}
	}
	if len(task.RehearsalCases) == 0 {
		return undeclaredRehearsalFinding(task), true
	}
	return Finding{}, false
}

func declaredRequirementClauses(requirements []string) []declaredRequirementClause {
	var clauses []declaredRequirementClause
	for requirementIndex, requirement := range requirements {
		text := modalMentionPattern.ReplaceAllString(requirement, "declared modal pair")
		matches := requirementModalPattern.FindAllStringSubmatchIndex(text, -1)
		for matchIndex, match := range matches {
			end := len(text)
			if matchIndex+1 < len(matches) {
				end = matches[matchIndex+1][0]
			}
			clauseText := strings.Trim(strings.TrimSpace(text[match[1]:end]), " .,:;—–-")
			if clauseText == "" {
				continue
			}
			clauses = append(clauses, declaredRequirementClause{
				Requirement: requirementIndex,
				Forbidden:   match[2] >= 0,
				Words:       significantRequirementWords(clauseText),
			})
		}
	}
	return clauses
}

func declaresGateRehearsal(title string) bool {
	if !gateWordPattern.MatchString(title) {
		return false
	}
	return rehearsalWordPattern.MatchString(title) || proveWordPattern.MatchString(title)
}

func significantRequirementWords(clause string) []string {
	var words []string
	for _, word := range requirementWordPattern.FindAllString(strings.ToLower(clause), -1) {
		word = requirementWordStem(word)
		if word == "" || requirementStopWords[word] {
			continue
		}
		words = append(words, word)
	}
	return words
}

func requirementWordStem(word string) string {
	switch {
	case strings.HasPrefix(word, "committ") || word == "commits":
		return "commit"
	case len(word) > 4 && strings.HasSuffix(word, "ies"):
		return strings.TrimSuffix(word, "ies") + "y"
	case len(word) > 4 && strings.HasSuffix(word, "s"):
		return strings.TrimSuffix(word, "s")
	default:
		return word
	}
}

func commonNamedSubject(needed, forbidden []string) (string, bool) {
	forbiddenWords := make(map[string]bool, len(forbidden))
	for _, word := range forbidden {
		forbiddenWords[word] = true
	}
	neededWords := make(map[string]bool, len(needed))
	for _, word := range needed {
		neededWords[word] = true
	}
	union := make(map[string]bool, len(neededWords)+len(forbiddenWords))
	for word := range neededWords {
		union[word] = true
	}
	for word := range forbiddenWords {
		union[word] = true
	}

	shared := 0
	subject := ""
	for _, word := range needed {
		if forbiddenWords[word] && subject == "" {
			subject = word
		}
	}
	for word := range neededWords {
		if forbiddenWords[word] {
			shared++
		}
	}
	if shared == 1 && subject == "commit" {
		return subject, true
	}
	if shared < 2 || shared*2 < len(union) {
		return "", false
	}
	return subject, true
}

func completeRehearsalCase(entry string) bool {
	lower := strings.ToLower(strings.TrimSpace(entry))
	if !strings.HasPrefix(lower, "case:") {
		return false
	}
	separator := strings.Index(lower, "; observation:")
	if separator < 0 {
		return false
	}
	caseText := strings.TrimSpace(entry[len("case:"):separator])
	observationText := strings.TrimSpace(entry[separator+len("; observation:"):])
	return caseText != "" && observationText != ""
}

func undeclaredRehearsalFinding(task spec.Task) Finding {
	return Finding{
		Code:     CodeRehearsalUndeclared,
		Severity: SeverityError,
		Summary:  task.File + " declares a gate rehearsal but has no complete Rehearsal Cases declaration",
		Where:    []Location{{Path: task.File, Line: 1}},
		Fix:      "Add a ## Rehearsal Cases section with one `- Case: <case>; Observation: <observation>` entry for each case.",
	}
}

var requirementStopWords = map[string]bool{
	"a": true, "an": true, "and": true, "another": true, "as": true,
	"at": true, "be": true, "by": true, "declared": true, "each": true,
	"for": true, "from": true, "in": true, "into": true, "is": true,
	"it": true, "its": true, "no": true, "of": true, "on": true,
	"named": true, "one": true, "or": true, "over": true, "same": true, "that": true,
	"the": true, "their": true, "them": true, "this": true, "to": true,
	"when": true, "whose": true, "with": true,
	"add": true, "change": true, "create": true, "decide": true,
	"declare": true, "define": true, "edit": true, "implement": true,
	"keep": true, "leave": true, "prove": true, "report": true,
	"requirement": true, "reuse": true, "run": true, "task": true,
	"update": true, "use": true,
}
