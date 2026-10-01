package judge

// Suite: embedded measured settings.
// Invariant: loading preserves the authored settings and compiles all patterns.
// Boundary IN: local question file and Load.
// Boundary OUT: requests and model answers.

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"
)

func loadQuestions(t *testing.T) Questions {
	t.Helper()
	q, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	return q
}

func TestQuestionFileLoads(t *testing.T) {
	// The Spec is active until it is archived into docs/history/specs.
	tech, err := os.ReadFile("../../docs/specs/0205-an-advisory-judge-for-spec-authoring/_techspec.md")
	if errors.Is(err, os.ErrNotExist) {
		tech, err = os.ReadFile("../../docs/history/specs/0205-an-advisory-judge-for-spec-authoring/_techspec.md")
	}
	if err != nil {
		t.Fatal(err)
	}
	block := strings.SplitN(strings.SplitN(strings.SplitN(string(tech), "### Questions and thresholds", 2)[1], "```json\n", 2)[1], "```", 2)[0]
	// Remove only the exact new member and its separating comma.
	original := strings.Replace(string(questionFile), ",\n"+strings.TrimSuffix(groupingBlock(t), "\n"), "", 1)
	if !bytes.Equal([]byte(original), []byte(block)) {
		t.Fatal("original embedded settings differ from the TechSpec block")
	}

	q := loadQuestions(t)
	var file struct {
		Pinned     string       `json:"pinned_model"`
		Transports []Transport  `json:"transports"`
		Language   LanguageGate `json:"language"`
		Price      float64      `json:"usd_per_million_input_tokens"`
		Ceiling    float64      `json:"monthly_ceiling_usd"`
	}
	if err := json.Unmarshal([]byte(block), &file); err != nil {
		t.Fatal(err)
	}
	if q.PinnedModel != file.Pinned || len(q.Transports) != len(file.Transports) || q.USDPerMillionInputTokens != file.Price || q.MonthlyCeilingUSD != file.Ceiling {
		t.Fatalf("settings lost: %+v", q)
	}
	for i, transport := range file.Transports {
		if q.Transports[i] != transport {
			t.Fatalf("transport %d differs", i)
		}
	}
	if !q.AcceptedModel.MatchString(q.Transports[1].RequestModel) || !q.AcceptedModel.MatchString(q.Transports[1].Name+"/"+q.PinnedModel+"-20260917") || q.AcceptedModel.MatchString("unmeasured") {
		t.Fatal("accepted model pattern not compiled faithfully")
	}
	if q.Grouping.SourceScrub == nil {
		t.Fatal("source scrub not compiled")
	}
	if !q.Citation.Attribution.MatchString("ADR-0123 keeps the gate") || !q.Goal.CoverageLine.MatchString("- Goals 1-2 → Reliable reader.") || !q.Goal.GenericSectionTitle.MatchString("Coverage Map") {
		t.Fatal("extraction patterns not compiled")
	}
	if q.Citation.Question.Instructions == "" || len(q.Citation.Question.Criteria) != 3 || q.Goal.Question.Instructions == "" || len(q.Goal.Question.Criteria) != 2 || q.Citation.RaiseMinConfidence == 0 || q.Goal.RaiseWhenNoulBelow == 0 || len(q.Language.EnglishWords) != len(file.Language.EnglishWords) {
		t.Fatal("questions or thresholds lost")
	}
}
