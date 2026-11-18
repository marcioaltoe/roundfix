package config

import (
	"fmt"
	"reflect"
	"testing"
)

func TestAnAuthoringEntryAcceptsAFallback(t *testing.T) {
	t.Parallel()
	for _, phase := range []string{"idea", "prd", "techspec", "tasks"} {
		for _, effort := range []string{"", "high"} {
			t.Run(phase+"/effort_"+effort, func(t *testing.T) {
				fallback := "model: opus"
				if effort != "" {
					fallback += ", reasoning_effort: " + effort
				}
				// Put fallback first to prove validation does not depend on key order.
				text := fmt.Sprintf("skills: {authoring: {claude: {%s: {fallback: {%s}, model: fable, reasoning_effort: xhigh}}}}", phase, fallback)
				loaded, err := ResolveConfigProposal(nil, []byte(text))
				if err != nil {
					t.Fatal(err)
				}
				h := loaded.Skills.Authoring.Claude
				entry := map[string]*AuthoringEntry{"idea": h.Idea, "prd": h.PRD, "techspec": h.TechSpec, "tasks": h.Tasks}[phase]
				want := &AuthoringEntry{Model: "fable", ReasoningEffort: "xhigh", Fallback: &AuthoringFallback{Model: "opus", ReasoningEffort: effort}}
				if !reflect.DeepEqual(entry, want) {
					t.Fatalf("entry=%+v want=%+v", entry, want)
				}
			})
		}
	}
}

func TestAFallbackIsRefusedWhereItEscapesNoLimit(t *testing.T) {
	t.Parallel()
	cases := []struct{ name, phase, entry, message string }{
		{"research", "research", "model: fable, fallback: {model: opus}", "skills.authoring.claude.research.fallback is not a supported key; expected model, reasoning_effort or tools"},
		{"scalar", "prd", "model: fable, fallback: opus", "skills.authoring.claude.prd.fallback must be a mapping"},
		{"null", "prd", "model: fable, fallback: null", "skills.authoring.claude.prd.fallback must be a mapping"},
		{"sequence", "prd", "model: fable, fallback: [opus]", "skills.authoring.claude.prd.fallback must be a mapping"},
		{"unknown_key", "prd", "model: fable, fallback: {effort: high}", "skills.authoring.claude.prd.fallback.effort is not a supported key; expected model or reasoning_effort"},
		{"no_chain", "prd", "model: fable, fallback: {model: opus, fallback: {model: sonnet}}", "skills.authoring.claude.prd.fallback.fallback is not a supported key; expected model or reasoning_effort"},
		{"malformed_model", "prd", "model: fable, fallback: {model: bad model}", `skills.authoring.claude.prd.fallback.model "bad model" is not a model name; use letters, digits and . _ : / @ [ ] - without spaces`},
		{"non_string_model", "prd", "model: fable, fallback: {model: 123}", `skills.authoring.claude.prd.fallback.model "123" is not a model name; use letters, digits and . _ : / @ [ ] - without spaces`},
		{"malformed_effort", "prd", "model: fable, fallback: {model: opus, reasoning_effort: hihg}", `skills.authoring.claude.prd.fallback.reasoning_effort "hihg" is not supported for claude; expected low, medium, high, xhigh or max`},
		{"non_string_effort", "prd", "model: fable, fallback: {model: opus, reasoning_effort: true}", `skills.authoring.claude.prd.fallback.reasoning_effort "true" is not supported for claude; expected low, medium, high, xhigh or max`},
		{"no_model", "prd", "model: fable, fallback: {reasoning_effort: high}", "skills.authoring.claude.prd.fallback declares no model; set model"},
		{"empty", "prd", "model: fable, fallback: {}", "skills.authoring.claude.prd.fallback declares no model; set model"},
		{"no_entry_model", "prd", "reasoning_effort: high, fallback: {model: opus}", "skills.authoring.claude.prd.fallback needs skills.authoring.claude.prd.model"},
		{"only_fallback", "prd", "fallback: {model: opus}", "skills.authoring.claude.prd.fallback needs skills.authoring.claude.prd.model"},
		{"same_model", "prd", "model: fable, fallback: {model: fable}", `skills.authoring.claude.prd.fallback.model "fable" is the declared model; a fallback must name another model`},
		{"same_model_case", "prd", "fallback: {model: FaBlE}, model: FABLE", `skills.authoring.claude.prd.fallback.model "FaBlE" is the declared model; a fallback must name another model`},
		{"unknown_before_malformed", "prd", "fallback: {model: bad model, other: true}", "skills.authoring.claude.prd.fallback.other is not a supported key; expected model or reasoning_effort"},
		{"malformed_before_missing", "prd", "fallback: {reasoning_effort: hihg}", `skills.authoring.claude.prd.fallback.reasoning_effort "hihg" is not supported for claude; expected low, medium, high, xhigh or max`},
		{"fallback_model_before_entry_model", "prd", "fallback: {}", "skills.authoring.claude.prd.fallback declares no model; set model"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			text := fmt.Sprintf("skills: {authoring: {claude: {%s: {%s}}}}", tc.phase, tc.entry)
			_, err := ResolveConfigProposal(nil, []byte(text))
			want := fmt.Sprintf("parse config %q: %s", "Project Config proposal", tc.message)
			if err == nil || err.Error() != want {
				t.Fatalf("error=%v want=%s", err, want)
			}
		})
	}
	// Model comparisons are case insensitive for every Authoring Phase.
	for _, phase := range []string{"idea", "techspec", "tasks"} {
		t.Run(phase+"_same_model", func(t *testing.T) {
			_, err := ResolveConfigProposal(nil, []byte(fmt.Sprintf("skills: {authoring: {claude: {%s: {model: fable, fallback: {model: FABLE}}}}}", phase)))
			want := fmt.Sprintf(`parse config "Project Config proposal": skills.authoring.claude.%s.fallback.model "FABLE" is the declared model; a fallback must name another model`, phase)
			if err == nil || err.Error() != want {
				t.Fatalf("error=%v want=%s", err, want)
			}
		})
	}
}

func TestRecommendedAuthoringFallbacksAreIndependentCopies(t *testing.T) {
	t.Parallel()
	want := RecommendedAuthoringDeclaration()
	changed := RecommendedAuthoringDeclaration()
	for _, entry := range []*AuthoringEntry{changed.Idea, changed.PRD, changed.TechSpec, changed.Tasks} {
		entry.Fallback.Model = "changed"
		entry.Fallback.ReasoningEffort = "low"
	}
	if got := RecommendedAuthoringDeclaration(); !reflect.DeepEqual(got, want) {
		t.Fatalf("fallbacks alias the recommendation: got=%+v want=%+v", got, want)
	}
}
