// Boundary: project-only Authoring Question Preference loading and diagnostics.
package config

import (
	"bytes"
	"strings"
	"testing"
)

func TestTheAuthoringQuestionPreferenceIsAProjectConfigValue(t *testing.T) {
	t.Parallel()
	t.Run("builtin_waits", func(t *testing.T) {
		t.Parallel()
		if got := Builtin().Authoring; got != (Authoring{OpenQuestions: OpenQuestionsWait}) {
			t.Fatalf("builtin authoring = %+v", got)
		}
	})
	for _, tc := range []struct {
		name, project, want string
		declared            bool
	}{
		{name: "absent_key_waits", project: "defaults:\n  agent: codex\n", want: OpenQuestionsWait},
		{name: "project_config_wait", project: "authoring:\n  open_questions: wait\n", want: OpenQuestionsWait, declared: true},
		{name: "project_config_follow", project: "authoring:\n  open_questions: follow\n", want: OpenQuestionsFollow, declared: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			proposal, err := ResolveConfigProposal(nil, []byte(tc.project))
			if err != nil {
				t.Fatal(err)
			}
			loaded, err := Load(itemBinaryLoadOptions(t, nil, []byte(tc.project)))
			if err != nil {
				t.Fatal(err)
			}
			want := Authoring{OpenQuestions: tc.want, OpenQuestionsDeclared: tc.declared}
			for name, config := range map[string]Config{"proposal": proposal, "load": loaded.Config} {
				if config.Authoring != want {
					t.Errorf("%s authoring = %+v, want %+v", name, config.Authoring, want)
				}
			}
		})
	}
	t.Run("user_config_ignored_once", func(t *testing.T) {
		t.Parallel()
		const warning = "config: authoring.open_questions in User Config is ignored; set authoring.open_questions in Project Config\n"
		for _, value := range []string{"follow", "maybe"} {
			user := []byte("authoring:\n  open_questions: " + value + "\n")
			proposal, err := ResolveConfigProposal(user, nil)
			if err != nil {
				t.Fatal(err)
			}
			options := itemBinaryLoadOptions(t, user, nil)
			var diagnostics bytes.Buffer
			options.Stderr = &diagnostics
			loaded, err := Load(options)
			if err != nil {
				t.Fatal(err)
			}
			for name, config := range map[string]Config{"proposal": proposal, "load": loaded.Config} {
				if got := config.Authoring; got != (Authoring{OpenQuestions: OpenQuestionsWait}) {
					t.Errorf("%s with User Config %s: %+v", name, value, got)
				}
			}
			if diagnostics.String() != warning {
				t.Errorf("User Config %s warning = %q, want %q", value, diagnostics.String(), warning)
			}
		}
	})
	for _, tc := range []struct{ name, value string }{
		{"refuses_maybe", "maybe"}, {"refuses_a_number", "1"},
		{"refuses_a_list", "[follow]"}, {"refuses_a_mapping", "{value: follow}"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			project := []byte("authoring:\n  open_questions: " + tc.value + "\n")
			_, proposalErr := ResolveConfigProposal(nil, project)
			_, loadErr := Load(itemBinaryLoadOptions(t, nil, project))
			for name, err := range map[string]error{"proposal": proposalErr, "load": loadErr} {
				if err == nil || !strings.HasPrefix(err.Error(), "parse config \"") || !strings.HasSuffix(err.Error(), "authoring.open_questions must be wait or follow") {
					t.Errorf("%s error = %v", name, err)
				}
			}
		})
	}
	t.Run("refuses_an_unknown_key", func(t *testing.T) {
		t.Parallel()
		project := []byte("authoring:\n  other: true\n")
		_, proposalErr := ResolveConfigProposal(nil, project)
		_, loadErr := Load(itemBinaryLoadOptions(t, nil, project))
		for name, err := range map[string]error{"proposal": proposalErr, "load": loadErr} {
			if err == nil || !strings.Contains(err.Error(), "field other not found") {
				t.Errorf("%s error = %v", name, err)
			}
		}
	})
	t.Run("default_configuration_omits_the_key", func(t *testing.T) {
		t.Parallel()
		for name, content := range map[string]string{"user": DefaultConfigYAML(), "project": DefaultProjectConfigYAML()} {
			if strings.Contains(content, "open_questions") {
				t.Errorf("%s default holds open_questions", name)
			}
		}
	})
}
