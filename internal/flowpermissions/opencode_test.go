package flowpermissions

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadOpenCode(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, file, data string
		state            State
	}{
		{"no_file_is_not_needed", "", "", NotNeeded},
		{"no_shell_rule_is_not_needed", "opencode.json", `{"permissions":[{"action":"edit","effect":"deny"}],"permission":{"edit":"ask","bash":"allow"}}`, NotNeeded},
		{"permissions_list_asks_about_shell", "opencode.json", `{"permissions":[{"action":"shell","effect":"ask"}]}`, OwnRules},
		{"permissions_list_wildcard_denies_shell", "opencode.json", `{"permissions":[{"action":"*","effect":"deny"}]}`, OwnRules},
		{"permission_bash_ask", "opencode.json", `{"permission":{"bash":"ask"}}`, OwnRules},
		{"permission_bash_pattern_deny", "opencode.json", `{"permission":{"bash":{"*":"allow","git *":"deny"}}}`, OwnRules},
		{"permission_value_deny", "opencode.json", `{"permission":"deny"}`, OwnRules},
		{"permission_wildcard_ask", "opencode.json", `{"permission":{"*":"ask"}}`, OwnRules},
		{"jsonc_comments_and_trailing_commas", "opencode.jsonc", "{ // comment\n /* block */ \"permissions\":[{\"action\":\"shell\",\"effect\":\"deny\",},],\n}", OwnRules},
		{"comment_markers_inside_strings_kept", "opencode.jsonc", `{"description":"https://host/a/*b*/?x=\"//\"", "permission":{"bash":{"echo /* */ // , }":"ask"}}}`, OwnRules},
		{"unreadable_file", "opencode.json", `{"permission":`, Unreadable},
		{"unclosed_comment", "opencode.jsonc", `{} /*`, Unreadable},
		{"top_level_null", "opencode.json", `null`, Unreadable},
		{"top_level_array", "opencode.json", `[]`, Unreadable},
		{"invalid_permission_shape", "opencode.json", `{"permission":false}`, Unreadable},
		{"invalid_bash_shape", "opencode.json", `{"permission":{"bash":[]}}`, Unreadable},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			if tc.file != "" {
				writeFixture(t, root, tc.file, []byte(tc.data))
			}
			state, reason := ReadOpenCode(root)
			if state != tc.state {
				t.Fatalf("state %q (%s); want %q", state, reason, tc.state)
			}
			if state == OwnRules && reason != tc.file {
				t.Fatalf("file %q; want %q", reason, tc.file)
			}
			if state == Unreadable && !strings.Contains(reason, tc.file) {
				t.Fatalf("reason does not name file: %q", reason)
			}
			if state == NotNeeded && reason != "" {
				t.Fatalf("unexpected reason: %q", reason)
			}
		})
	}
	t.Run("both_files_read_before_reporting", func(t *testing.T) {
		root := t.TempDir()
		writeFixture(t, root, "opencode.json", []byte(`{"permission":"ask"}`))
		writeFixture(t, root, "opencode.jsonc", []byte(`{`))
		state, reason := ReadOpenCode(root)
		if state != Unreadable || !strings.Contains(reason, "opencode.jsonc") {
			t.Fatalf("%s: %s", state, reason)
		}
	})
	t.Run("filesystem_read_error", func(t *testing.T) {
		root := t.TempDir()
		if err := os.Mkdir(filepath.Join(root, "opencode.json"), 0755); err != nil {
			t.Fatal(err)
		}
		if state, reason := ReadOpenCode(root); state != Unreadable || reason == "" {
			t.Fatalf("%s: %s", state, reason)
		}
	})
}
