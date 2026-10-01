package judge

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"regexp"
)

//go:embed questions.json
var questionFile []byte

type Transport struct {
	Name         string `json:"name"`
	KeyVariable  string `json:"key_variable"`
	Endpoint     string `json:"endpoint"`
	RequestModel string `json:"request_model"`
}

type Question struct {
	Type         string            `json:"type"`
	Instructions string            `json:"instructions"`
	Criteria     map[string]string `json:"criteria"`
}

type LanguageGate struct {
	EnglishWords       []string `json:"english_words"`
	PortugueseWords    []string `json:"portuguese_words"`
	MinEnglishShare    float64  `json:"min_english_share"`
	MaxPortugueseShare float64  `json:"max_portuguese_share"`
}

type CitationJudgment struct {
	QuestionID             string         `json:"question_id"`
	Question               Question       `json:"question"`
	RaiseWhenChoiceIsNot   string         `json:"raise_when_choice_is_not"`
	RaiseMinConfidence     float64        `json:"raise_min_confidence"`
	AttributionPattern     string         `json:"attribution_pattern"`
	Attribution            *regexp.Regexp `json:"-"`
	ClaimMinChars          int            `json:"claim_min_chars"`
	ClaimMaxChars          int            `json:"claim_max_chars"`
	DecisionRecordMaxChars int            `json:"decision_record_max_chars"`
}

type GoalJudgment struct {
	QuestionID                        string         `json:"question_id"`
	Question                          Question       `json:"question"`
	RaiseWhenNoulBelow                float64        `json:"raise_when_noul_below"`
	CoverageLinePattern               string         `json:"coverage_line_pattern"`
	GenericSectionTitlePattern        string         `json:"generic_section_title_pattern"`
	CoverageLine, GenericSectionTitle *regexp.Regexp `json:"-"`
	SectionTitleMinChars              int            `json:"section_title_min_chars"`
	SectionMinChars                   int            `json:"section_min_chars"`
	SectionMaxChars                   int            `json:"section_max_chars"`
}

type Questions struct {
	PinnedModel                                 string
	AcceptedModel                               *regexp.Regexp
	Transports                                  []Transport
	USDPerMillionInputTokens, MonthlyCeilingUSD float64
	Language                                    LanguageGate
	Citation                                    CitationJudgment
	Goal                                        GoalJudgment
}

// Load is the sole source of measured settings and compiles all four patterns.
func Load() (Questions, error) {
	var file struct {
		PinnedModel              string       `json:"pinned_model"`
		AcceptedModelPattern     string       `json:"accepted_model_pattern"`
		Transports               []Transport  `json:"transports"`
		USDPerMillionInputTokens float64      `json:"usd_per_million_input_tokens"`
		MonthlyCeilingUSD        float64      `json:"monthly_ceiling_usd"`
		Language                 LanguageGate `json:"language"`
		Judgments                struct {
			Citation CitationJudgment `json:"citation-support"`
			Goal     GoalJudgment     `json:"goal-mechanism"`
		} `json:"judgments"`
	}
	if err := json.Unmarshal(questionFile, &file); err != nil {
		return Questions{}, fmt.Errorf("load judge questions: %w", err)
	}
	q := Questions{PinnedModel: file.PinnedModel, Transports: file.Transports, USDPerMillionInputTokens: file.USDPerMillionInputTokens, MonthlyCeilingUSD: file.MonthlyCeilingUSD, Language: file.Language, Citation: file.Judgments.Citation, Goal: file.Judgments.Goal}
	for _, pattern := range []struct {
		name, text string
		dst        **regexp.Regexp
	}{
		{"accepted model", file.AcceptedModelPattern, &q.AcceptedModel},
		{"attribution", q.Citation.AttributionPattern, &q.Citation.Attribution},
		{"coverage line", q.Goal.CoverageLinePattern, &q.Goal.CoverageLine},
		{"generic section title", q.Goal.GenericSectionTitlePattern, &q.Goal.GenericSectionTitle},
	} {
		compiled, err := regexp.Compile(pattern.text)
		if err != nil {
			return Questions{}, fmt.Errorf("compile %s: %w", pattern.name, err)
		}
		*pattern.dst = compiled
	}
	return q, nil
}
