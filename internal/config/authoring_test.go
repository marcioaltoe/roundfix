package config

import (
	"bytes"
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	"roundfix/internal/gittest"
)

const authoringExample = `skills:
  authoring:
    claude:
      idea:
        model: opus
        reasoning_effort: xhigh
      prd:
        model: opus
        reasoning_effort: xhigh
      techspec:
        model: opus
        reasoning_effort: xhigh
      tasks:
        model: opus
        reasoning_effort: high
      research:
        model: sonnet
        reasoning_effort: medium
        tools: [WebSearch, WebFetch]
`

func TestTheAuthoringModelDeclarationIsReadFromProjectConfig(t *testing.T) {
	t.Parallel()
	for _, text := range []string{authoringExample, "", "skills: {}", "skills: {authoring: {}}", "skills: {authoring: {claude: {}}}"} {
		t.Run(text, func(t *testing.T) {
			home, root := t.TempDir(), t.TempDir()
			gittest.MarkRepository(t, root)
			mustWrite(t, filepath.Join(root, ".roundfixrc.yml"), text)
			var warnings bytes.Buffer
			loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: root, Stderr: &warnings})
			if err != nil || warnings.Len() != 0 {
				t.Fatalf("load=%v warnings=%s", err, &warnings)
			}
			if text != authoringExample {
				if loaded.Config.Skills.Authoring.Claude != nil {
					t.Fatalf("declared: %+v", loaded.Config.Skills)
				}
				return
			}
			want := &AuthoringHarness{
				Idea:     &AuthoringEntry{Model: "opus", ReasoningEffort: "xhigh"},
				PRD:      &AuthoringEntry{Model: "opus", ReasoningEffort: "xhigh"},
				TechSpec: &AuthoringEntry{Model: "opus", ReasoningEffort: "xhigh"},
				Tasks:    &AuthoringEntry{Model: "opus", ReasoningEffort: "high"},
				Research: &ResearchEntry{AuthoringEntry: AuthoringEntry{Model: "sonnet", ReasoningEffort: "medium"}, Tools: []string{"WebSearch", "WebFetch"}},
			}
			if !reflect.DeepEqual(loaded.Config.Skills.Authoring.Claude, want) {
				t.Fatalf("got=%+v want=%+v", loaded.Config.Skills.Authoring.Claude, want)
			}
		})
	}
}

func TestTheAuthoringModelDeclarationRefusesWhatItDoesNotKnow(t *testing.T) {
	t.Parallel()
	cases := []struct{ name, text, message string }{
		{"skills_key", "skills: {other: {}}", "skills.other is not a supported key; expected authoring"},
		{"harness", "skills: {authoring: {codex: {}}}", "skills.authoring.codex is not a supported harness; this Roundfix supports claude"},
		{"phase_key", "skills: {authoring: {claude: {other: {}}}}", "skills.authoring.claude.other is not a supported key; expected idea, prd, techspec, tasks or research"},
		{"phase_field", "skills: {authoring: {claude: {prd: {effort: high}}}}", "skills.authoring.claude.prd.effort is not a supported key; expected model or reasoning_effort"},
		{"research_field", "skills: {authoring: {claude: {research: {effort: high}}}}", "skills.authoring.claude.research.effort is not a supported key; expected model, reasoning_effort or tools"},
		{"empty_phase", "skills: {authoring: {claude: {prd: {}}}}", "skills.authoring.claude.prd declares nothing; set model, reasoning_effort or both"},
		{"empty_research", "skills: {authoring: {claude: {research: {}}}}", "skills.authoring.claude.research declares nothing; set model, reasoning_effort or tools"},
		{"model_with_a_space", "skills: {authoring: {claude: {prd: {model: bad model}}}}", `skills.authoring.claude.prd.model "bad model" is not a model name; use letters, digits and . _ : / @ [ ] - without spaces`},
		{"effort_hihg", "skills: {authoring: {claude: {prd: {reasoning_effort: hihg}}}}", `skills.authoring.claude.prd.reasoning_effort "hihg" is not supported for claude; expected low, medium, high, xhigh or max`},
		{"tool_Bash", "skills: {authoring: {claude: {research: {tools: [Bash]}}}}", `skills.authoring.claude.research.tools names "Bash"; expected WebSearch, WebFetch or a name that begins with mcp__`},
		{"duplicate_tool", "skills: {authoring: {claude: {research: {tools: [WebSearch, WebSearch]}}}}", `skills.authoring.claude.research.tools names "WebSearch" twice`},
		{"scalar_phase", "skills: {authoring: {claude: {prd: opus}}}", "skills.authoring.claude.prd must be a mapping"},
		{"scalar_skills", "skills: false", "skills must be a mapping"},
		{"scalar_authoring", "skills: {authoring: []}", "skills.authoring must be a mapping"},
		{"scalar_claude", "skills: {authoring: {claude: null}}", "skills.authoring.claude must be a mapping"},
		{"non_string_model", "skills: {authoring: {claude: {prd: {model: 123}}}}", `skills.authoring.claude.prd.model "123" is not a model name; use letters, digits and . _ : / @ [ ] - without spaces`},
		{"non_string_effort", "skills: {authoring: {claude: {prd: {reasoning_effort: true}}}}", `skills.authoring.claude.prd.reasoning_effort "true" is not supported for claude; expected low, medium, high, xhigh or max`},
		{"scalar_tools", "skills: {authoring: {claude: {research: {tools: Bash}}}}", `skills.authoring.claude.research.tools names "Bash"; expected WebSearch, WebFetch or a name that begins with mcp__`},
		{"non_string_tool", "skills: {authoring: {claude: {research: {tools: [123]}}}}", `skills.authoring.claude.research.tools names "123"; expected WebSearch, WebFetch or a name that begins with mcp__`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ResolveConfigProposal(nil, []byte(tc.text))
			want := fmt.Sprintf("parse config %q: %s", "Project Config proposal", tc.message)
			if err == nil || err.Error() != want {
				t.Fatalf("error=%v want=%s", err, want)
			}
		})
	}
}

func TestUserConfigIgnoresTheAuthoringModelDeclaration(t *testing.T) {
	t.Parallel()
	for _, project := range []string{"", "skills: {authoring: {claude: {prd: {model: project}}}}"} {
		t.Run(project, func(t *testing.T) {
			home, root := t.TempDir(), t.TempDir()
			gittest.MarkRepository(t, root)
			mustMkdir(t, filepath.Join(home, ".roundfix"))
			mustWrite(t, filepath.Join(home, ".roundfix", "config.yml"), "skills: [invalid, ignored]\n")
			mustWrite(t, filepath.Join(root, ".roundfixrc.yml"), project)
			var warnings bytes.Buffer
			loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: root, Stderr: &warnings})
			if err != nil {
				t.Fatal(err)
			}
			if warnings.String() != "config: skills in User Config is ignored; set skills in Project Config\n" {
				t.Fatalf("warnings=%s", &warnings)
			}
			h := loaded.Config.Skills.Authoring.Claude
			if project == "" {
				if h != nil {
					t.Fatalf("user declaration leaked: %+v", h)
				}
			} else if h == nil || h.PRD == nil || h.PRD.Model != "project" || h.Idea != nil || h.Research != nil {
				t.Fatalf("declaration=%+v", h)
			}
		})
	}
}
