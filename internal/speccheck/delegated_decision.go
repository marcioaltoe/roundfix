package speccheck

import (
	"regexp"
	"strings"
	"time"
)

// DelegatedDecision is an adopted authoring option, never a person's answer.
// Problem names the first defect when the entry cannot be read.
type DelegatedDecision struct {
	Line                                                   int
	Date, Phase, Question, Options, Adopted, Under, Answer string
	Problem                                                string
}

var (
	delegatedHead   = regexp.MustCompile(`^- delegated: (\d{4}-\d{2}-\d{2}), (prd|techspec|tasks)$`)
	delegatedField  = regexp.MustCompile(`^  - (question|options|adopted|under|answer): (\S.*)$`)
	delegatedGrant  = regexp.MustCompile(`^the user's answer of (\d{4}-\d{2}-\d{2})$`)
	delegatedAnswer = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}), (kept|changed) — (\S.*)$`)
)

// DelegatedDecisions reads every entry beginning "- delegated:" in a PRD's
// Decisions section, preserving head line numbers and ignoring fenced examples.
func DelegatedDecisions(prd []byte) []DelegatedDecision {
	lines := strings.Split(string(proofOutsideFences(prd)), "\n")
	var entries []DelegatedDecision
	inSection := false
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if strings.HasPrefix(line, "## ") {
			inSection = line == "## Decisions"
			continue
		}
		if !inSection || !strings.HasPrefix(line, "- delegated:") {
			continue
		}
		entry := DelegatedDecision{Line: i + 1}
		head := delegatedHead.FindStringSubmatch(line)
		if head != nil {
			entry.Date, entry.Phase = head[1], head[2]
		}
		fields := make(map[string]string)
		counts := make(map[string]int)
		for i+1 < len(lines) {
			field := delegatedField.FindStringSubmatch(lines[i+1])
			if field == nil {
				break
			}
			i++
			counts[field[1]]++
			if counts[field[1]] == 1 {
				fields[field[1]] = field[2]
			}
		}
		entry.Question, entry.Options, entry.Adopted = fields["question"], fields["options"], fields["adopted"]
		entry.Under, entry.Answer = fields["under"], fields["answer"]
		entry.Problem = delegatedProblem(entry, head != nil, counts)
		entries = append(entries, entry)
	}
	return entries
}

func delegatedDateValid(date string) bool {
	_, err := time.Parse(time.DateOnly, date)
	return err == nil
}

func delegatedProblem(entry DelegatedDecision, headOK bool, counts map[string]int) string {
	if !headOK || !delegatedDateValid(entry.Date) {
		return "the head needs a YYYY-MM-DD date and the phase prd, techspec or tasks"
	}
	fields := []string{"question", "options", "adopted", "under", "answer"}
	for _, field := range fields {
		if counts[field] == 0 {
			return field + " is missing"
		}
	}
	for _, field := range fields {
		if counts[field] > 1 {
			return field + " appears twice"
		}
	}
	option, reason, ok := strings.Cut(entry.Adopted, " — ")
	if !ok || strings.TrimSpace(option) == "" || strings.TrimSpace(reason) == "" {
		return "adopted needs the option, an em dash and the author's reason"
	}
	grant := delegatedGrant.FindStringSubmatch(entry.Under)
	if entry.Under != "Project Config" && (grant == nil || !delegatedDateValid(grant[1])) {
		return "under must be Project Config or the user's answer of YYYY-MM-DD"
	}
	answer := delegatedAnswer.FindStringSubmatch(entry.Answer)
	if entry.Answer != "none yet" && (answer == nil || !delegatedDateValid(answer[1]) || strings.TrimSpace(answer[3]) == "") {
		return "answer must be none yet, or YYYY-MM-DD, kept or changed, an em dash and the user's words"
	}
	return ""
}
