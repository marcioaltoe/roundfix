package judge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type Stage string

type PendingJudgment struct {
	Kind, Artifact string
	Line           int
	Target, Text   string
	state          json.RawMessage
}

type SkippedJudgment struct {
	Kind, Artifact string
	Line           int
	Target, Reason string
}

type SkippedArtifact struct{ Artifact, Reason string }

type Plan struct {
	Pending          []PendingJudgment
	Skipped          []SkippedJudgment
	SkippedArtifacts []SkippedArtifact
}

var adrToken = regexp.MustCompile(`ADR-\d{4}`)
var specReference = regexp.MustCompile(`\(Spec \d{4}\)`)
var itemMarker = regexp.MustCompile(`^(?:- |\d+\. )`)

// PlanSpec constructs states exclusively from bounded Source values.
func PlanSpec(q Questions, repoRoot, specDir string, stage Stage) (Plan, error) {
	var plan Plan
	if stage != "" && stage != "prd" && stage != "techspec" {
		return plan, fmt.Errorf("unknown judge stage %q", stage)
	}
	prd, prdErr := readSpecArtifact(specDir, "_prd.md")
	readable := func(name string, src Source, err error) bool {
		if err != nil {
			plan.SkippedArtifacts = append(plan.SkippedArtifacts, SkippedArtifact{name, "not a regular file in the Spec directory"})
			return false
		}
		if !q.Language.isEnglish(src.text) {
			plan.SkippedArtifacts = append(plan.SkippedArtifacts, SkippedArtifact{name, "not English"})
			return false
		}
		return true
	}
	prdOK := readable("_prd.md", prd, prdErr)
	seen := make(map[string]bool)
	if prdOK && stage != "techspec" {
		planClaims(q, repoRoot, prd, &plan, seen)
	}
	if stage != "prd" {
		tech, err := readSpecArtifact(specDir, "_techspec.md")
		if readable("_techspec.md", tech, err) {
			planClaims(q, repoRoot, tech, &plan, seen)
			if prdOK {
				planGoals(q, prd, tech, &plan, seen)
			}
		}
	}
	anchors, skipped, err := adoptedSources(q, specDir)
	if err != nil {
		plan.SkippedArtifacts = append(plan.SkippedArtifacts, SkippedArtifact{filepath.Join(specDir, "references", "_index.md"), "not a regular file in its directory"})
	}
	plan.SkippedArtifacts = append(plan.SkippedArtifacts, skipped...)
	if len(anchors) > 0 {
		candidates, skipped, err := openSources(q, repoRoot)
		if err != nil {
			plan.SkippedArtifacts = append(plan.SkippedArtifacts, SkippedArtifact{repoRoot, "grouping sources unreadable: " + err.Error()})
		}
		plan.SkippedArtifacts = append(plan.SkippedArtifacts, skipped...)
		for _, anchor := range anchors {
			for _, candidate := range candidates {
				state := struct {
					First  string `json:"first"`
					Second string `json:"second"`
				}{prepareSource(q, anchor), prepareSource(q, candidate)}
				addPending(&plan, seen, PendingJudgment{Kind: "source-grouping", Artifact: anchor.path, Target: candidate.path}, state)
			}
		}
	}
	if prdErr != nil {
		return plan, fmt.Errorf("read Spec PRD: %w", prdErr)
	}
	return plan, nil
}

type textLine struct {
	text   string
	number int
}

func sourceLines(src Source) []textLine {
	_, body, offset, _ := splitFrontMatter(src.text)
	var lines []textLine
	for i, text := range strings.Split(body, "\n") {
		lines = append(lines, textLine{text, i + offset + 1})
	}
	return lines
}

func collapse(text string) string { return strings.Join(strings.Fields(text), " ") }
func cut(text string, limit int) string {
	runes := []rune(text)
	if len(runes) > limit {
		runes = runes[:limit]
	}
	return string(runes)
}

func sentences(text string) []string {
	runes := []rune(text)
	start := 0
	var result []string
	for i, r := range runes {
		if r != '.' && r != '!' && r != '?' {
			continue
		}
		j := i + 1
		for j < len(runes) && unicode.IsSpace(runes[j]) {
			j++
		}
		if j == i+1 || j == len(runes) {
			continue
		}
		next := runes[j]
		if (next >= 'A' && next <= 'Z') || next == '`' || next == '(' || next == '[' {
			result = append(result, string(runes[start:i+1]))
			start = j
		}
	}
	return append(result, string(runes[start:]))
}

func planClaims(q Questions, repoRoot string, src Source, plan *Plan, seen map[string]bool) {
	var paragraph []textLine
	flush := func() {
		if len(paragraph) == 0 {
			return
		}
		var texts []string
		for _, line := range paragraph {
			texts = append(texts, line.text)
		}
		for _, sentence := range sentences(strings.Join(texts, " ")) {
			tokens := adrToken.FindAllString(sentence, -1)
			if len(tokens) != 1 || !q.Citation.Attribution.MatchString(sentence) {
				continue
			}
			token := tokens[0]
			lineNumber := paragraph[0].number
			for _, line := range paragraph {
				if strings.Contains(line.text, token) {
					lineNumber = line.number
					break
				}
			}
			skip := func(reason string) {
				plan.Skipped = append(plan.Skipped, SkippedJudgment{"citation-support", src.path, lineNumber, token, reason})
			}
			adr, ok, err := readADR(repoRoot, strings.TrimPrefix(token, "ADR-"))
			if err != nil || !ok {
				skip("cited decision is not an accepted regular ADR")
				continue
			}
			claim := collapse(specReference.ReplaceAllString(strings.ReplaceAll(sentence, token, "the cited decision"), ""))
			runes := []rune(claim)
			if len(runes) < q.Citation.ClaimMinChars || len(runes) > q.Citation.ClaimMaxChars {
				skip("claim length outside the measured 60-500 characters")
				continue
			}
			runes[0] = unicode.ToUpper(runes[0])
			claim = string(runes)
			if !q.Language.isEnglish(adr.text) {
				skip("decision record is not English")
				continue
			}
			_, body, _, _ := splitFrontMatter(adr.text)
			state := struct {
				Claim          string `json:"claim"`
				DecisionRecord string `json:"decision_record"`
			}{claim, cut(strings.TrimSpace(body), q.Citation.DecisionRecordMaxChars)}
			addPending(plan, seen, PendingJudgment{Kind: "citation-support", Artifact: src.path, Line: lineNumber, Target: token, Text: claim}, state)
		}
		paragraph = nil
	}
	fenced := false
	for _, line := range sourceLines(src) {
		trimmed := strings.TrimSpace(line.text)
		if strings.HasPrefix(line.text, "```") {
			flush()
			fenced = !fenced
			continue
		}
		if fenced {
			continue
		}
		if trimmed == "" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "|") {
			flush()
			continue
		}
		if marker := itemMarker.FindString(trimmed); marker != "" {
			flush()
			line.text = strings.TrimPrefix(trimmed, marker)
		}
		paragraph = append(paragraph, line)
	}
	flush()
}

func addPending(plan *Plan, seen map[string]bool, pending PendingJudgment, state any) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	// States contain strings only, so encoding cannot fail.
	if err := encoder.Encode(state); err != nil {
		panic(err)
	}
	pending.state = bytes.TrimSuffix(buf.Bytes(), []byte("\n"))
	key := string(pending.state)
	if seen[key] {
		return
	}
	seen[key] = true
	plan.Pending = append(plan.Pending, pending)
}

type section struct {
	title, body string
	level, line int
}

func sections(src Source) []section {
	lines := sourceLines(src)
	var result []section
	for i, line := range lines {
		level := 0
		if strings.HasPrefix(line.text, "## ") {
			level = 2
		}
		if strings.HasPrefix(line.text, "### ") {
			level = 3
		}
		if level == 0 {
			continue
		}
		var body []string
		for _, next := range lines[i+1:] {
			if strings.HasPrefix(next.text, "# ") || strings.HasPrefix(next.text, "## ") || (level == 3 && strings.HasPrefix(next.text, "### ")) {
				break
			}
			body = append(body, next.text)
		}
		result = append(result, section{strings.TrimSpace(line.text[level+1:]), strings.TrimSpace(strings.Join(body, "\n")), level, line.number})
	}
	return result
}

func normalized(text string) string {
	return strings.ToLower(collapse(strings.NewReplacer("`", "", "*", "").Replace(text)))
}

func goals(prd Source) []string {
	var result []string
	for _, sec := range sections(prd) {
		if sec.level != 2 || sec.title != "Goals" {
			continue
		}
		current := ""
		flush := func() {
			if current != "" {
				result = append(result, collapse(current))
				current = ""
			}
		}
		for _, line := range strings.Split(sec.body, "\n") {
			text := strings.TrimSpace(line)
			if marker := itemMarker.FindString(text); marker != "" {
				flush()
				current = strings.TrimPrefix(text, marker)
			} else if text != "" && current != "" {
				current += " " + text
			}
		}
		flush()
		break
	}
	return result
}

func planGoals(q Questions, prd, tech Source, plan *Plan, seen map[string]bool) {
	goalList := goals(prd)
	allSections := sections(tech)
	var coverage []textLine
	for _, sec := range allSections {
		if sec.level != 2 || sec.title != "Coverage Map" {
			continue
		}
		// Reuse the untrimmed source to preserve continuation whitespace and lines.
		for _, line := range sourceLines(tech) {
			if line.number <= sec.line {
				continue
			}
			if strings.HasPrefix(line.text, "# ") || strings.HasPrefix(line.text, "## ") {
				break
			}
			if len(line.text) > 0 && unicode.IsSpace([]rune(line.text)[0]) && len(coverage) > 0 {
				coverage[len(coverage)-1].text += " " + strings.TrimSpace(line.text)
			} else {
				coverage = append(coverage, line)
			}
		}
		break
	}
	for _, line := range coverage {
		match := q.Goal.CoverageLine.FindStringSubmatch(strings.TrimSpace(line.text))
		if match == nil {
			continue
		}
		lo, err := strconv.Atoi(match[1])
		if err != nil {
			continue
		}
		hi := lo
		if match[2] != "" {
			hi, err = strconv.Atoi(match[2])
			if err != nil {
				continue
			}
		}
		var selected *section
		longest := 0
		for i := range allSections {
			sec := &allSections[i]
			title := normalized(sec.title)
			length := len([]rune(title))
			if q.Goal.GenericSectionTitle.MatchString(title) || len([]rune(sec.body)) < q.Goal.SectionMinChars || length < q.Goal.SectionTitleMinChars {
				continue
			}
			if length > longest && strings.Contains(normalized(match[3]), title) {
				selected = sec
				longest = length
			}
		}
		if selected == nil {
			plan.Skipped = append(plan.Skipped, SkippedJudgment{"goal-mechanism", tech.path, line.number, match[3], "no named section qualifies"})
			continue
		}
		for n := lo; n <= hi; n++ {
			target := fmt.Sprintf("Goal %d → %s", n, selected.title)
			if n < 1 || n > len(goalList) {
				plan.Skipped = append(plan.Skipped, SkippedJudgment{"goal-mechanism", tech.path, line.number, target, fmt.Sprintf("Goal %d is not in the PRD", n)})
				continue
			}
			state := struct {
				Goal         string `json:"goal"`
				SectionTitle string `json:"section_title"`
				Section      string `json:"section"`
			}{goalList[n-1], selected.title, cut(selected.body, q.Goal.SectionMaxChars)}
			addPending(plan, seen, PendingJudgment{Kind: "goal-mechanism", Artifact: tech.path, Line: line.number, Target: target, Text: goalList[n-1]}, state)
		}
	}
}
