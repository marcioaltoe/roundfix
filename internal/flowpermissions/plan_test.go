package flowpermissions

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func writeFixture(t *testing.T, root, name string, data []byte) {
	t.Helper()
	target := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, data, 0644); err != nil {
		t.Fatal(err)
	}
}

func readFixture(t *testing.T, root, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(name)))
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func currentClaude(t *testing.T) []byte {
	t.Helper()
	allow, ask := []string{}, []string{}
	for _, entry := range Entries() {
		allow = append(allow, ClaudeRule(entry))
	}
	for _, guard := range Guards() {
		ask = append(ask, guard.ClaudeRule)
	}
	return settingsFixture(t, allow, ask)
}

func settingsFixture(t *testing.T, allow, ask []string) []byte {
	t.Helper()
	data, err := json.Marshal(map[string]any{"permissions": map[string]any{"allow": allow, "ask": ask}})
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func previewFixture(t *testing.T, root string, harness Harness, action Action) Plan {
	t.Helper()
	p, err := Preview(root, harness, action)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

// Sequential: home_settings_ignored changes the process-wide HOME variable.
func TestPlan(t *testing.T) {
	t.Run("present_harnesses_in_order", func(t *testing.T) {
		root := t.TempDir()
		if got := Present(root); len(got) != 0 {
			t.Fatalf("empty repository: %v", got)
		}
		for _, name := range []string{".opencode", ".codex", ".claude"} {
			if err := os.Mkdir(filepath.Join(root, name), 0755); err != nil {
				t.Fatal(err)
			}
		}
		writeFixture(t, root, "opencode.json", []byte(`{}`))
		if got := Present(root); !slices.Equal(got, []Harness{Claude, Codex, OpenCode}) {
			t.Fatalf("%v", got)
		}
		for _, name := range []string{"opencode.json", "opencode.jsonc", ".opencode"} {
			isolated := t.TempDir()
			writeFixture(t, isolated, name, []byte(`{}`))
			if got := Present(isolated); !slices.Equal(got, []Harness{OpenCode}) {
				t.Fatalf("%s: %v", name, got)
			}
		}
	})
	for _, tc := range []struct {
		name    string
		harness Harness
		data    []byte
		state   State
	}{
		{"claude_not_applied", Claude, nil, NotApplied},
		{"claude_current", Claude, currentClaude(t), Current},
		{"claude_behind_missing_rule", Claude, settingsFixture(t, []string{ClaudeRule(Entries()[0])}, []string{}), Behind},
		{"claude_git_rule_only_is_not_applied", Claude, settingsFixture(t, []string{ClaudeRule(Entry{Words: []string{"git", "show"}})}, []string{}), NotApplied},
		{"claude_unreadable", Claude, []byte(`{"permissions":{"ask":null}}`), Unreadable},
		{"codex_not_applied", Codex, nil, NotApplied},
		{"codex_current", Codex, RenderCodexRules(), Current},
		{"codex_behind", Codex, (list{entries: []Entry{Entries()[0]}}).renderCodexRules(), Behind},
		{"codex_unreadable", Codex, []byte("arbitrary program\n"), Unreadable},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			file := ".claude/settings.json"
			if tc.harness == Codex {
				file = ".codex/rules/roundfix-flow.rules"
			}
			if tc.data != nil {
				writeFixture(t, root, file, tc.data)
			}
			p := previewFixture(t, root, tc.harness, ActionApply)
			if p.State != tc.state || p.File != file || p.FileExists != (tc.data != nil) {
				t.Fatalf("plan: %+v", p)
			}
			if tc.state == Unreadable && (p.Digest != "" || p.Reason == "" || len(p.Add)+len(p.AddGuards)+len(p.Remove) != 0) {
				t.Fatalf("unreadable plan: %+v", p)
			}
			if tc.state != Current && tc.state != Unreadable && p.Digest == "" {
				t.Fatalf("missing digest: %+v", p)
			}
			if tc.data != nil && !bytes.Equal(readFixture(t, root, file), tc.data) {
				t.Fatal("preview changed target")
			}
			if tc.data == nil {
				if _, err := os.Stat(filepath.Join(root, strings.Split(file, "/")[0])); !errors.Is(err, os.ErrNotExist) {
					t.Fatalf("preview created directory: %v", err)
				}
			}
		})
	}
	t.Run("claude_behind_retired_rule", func(t *testing.T) {
		root := t.TempDir()
		rules := shipped
		rules.retiredEntries = []string{"roundfix retired"}
		rules.retiredClaudeRules = []string{"Bash(git old *--output*)"}
		old := "Bash(roundfix retired *)"
		writeFixture(t, root, ".claude/settings.json", settingsFixture(t, []string{old}, rules.retiredClaudeRules))
		p, err := plan(root, Claude, ActionApply, rules)
		if err != nil || p.State != Behind || !slices.Equal(p.Remove, []string{old, rules.retiredClaudeRules[0]}) {
			t.Fatalf("%+v, %v", p, err)
		}
		allow, ask, err := ReadClaudeSettings(p.after)
		if err != nil || slices.Contains(allow, old) || slices.Contains(ask, rules.retiredClaudeRules[0]) || len(allow) != len(Entries()) || len(ask) != len(Guards()) {
			t.Fatalf("retired plan result: %v, %v, %v", allow, ask, err)
		}
		removal, err := plan(root, Claude, ActionRemove, rules)
		if err != nil || !slices.Equal(removal.Remove, p.Remove) {
			t.Fatalf("retired removal: %+v, %v", removal, err)
		}
	})
	t.Run("digest_covers_the_file_bytes", func(t *testing.T) {
		root := t.TempDir()
		writeFixture(t, root, ".claude/settings.json", []byte(`{}`))
		first := previewFixture(t, root, Claude, ActionApply)
		writeFixture(t, root, ".claude/settings.json", []byte(`{} `))
		second := previewFixture(t, root, Claude, ActionApply)
		if first.Digest == second.Digest || !slices.Equal(first.Add, second.Add) {
			t.Fatalf("digest did not bind bytes: %s, %s", first.Digest, second.Digest)
		}
	})
	t.Run("digest_is_stable", func(t *testing.T) {
		root := t.TempDir()
		a := previewFixture(t, root, Claude, ActionApply)
		b := previewFixture(t, root, Claude, ActionApply)
		if a.Digest != b.Digest || len(a.Digest) != 12 {
			t.Fatalf("%q, %q", a.Digest, b.Digest)
		}
		// Independent digest oracle, including group ordering and terminal newline.
		lines := []string{"claude", "apply", ".claude/settings.json", "absent"}
		added, guards := slices.Clone(a.Add), slices.Clone(a.AddGuards)
		slices.Sort(added)
		slices.Sort(guards)
		for _, rule := range added {
			lines = append(lines, "+"+rule)
		}
		for _, rule := range guards {
			lines = append(lines, "?"+rule)
		}
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(strings.Join(lines, "\n")+"\n")))
		if a.Digest != hash[:12] {
			t.Fatalf("digest %s; want %s", a.Digest, hash[:12])
		}
		rules := shipped
		rules.entries = slices.Clone(rules.entries)
		slices.Reverse(rules.entries)
		rules.guards = slices.Clone(rules.guards)
		slices.Reverse(rules.guards)
		reversed, err := plan(root, Claude, ActionApply, rules)
		if err != nil || a.Digest != reversed.Digest {
			t.Fatalf("order changed digest: %+v, %v", reversed, err)
		}
	})
	t.Run("nothing_to_write_has_no_digest", func(t *testing.T) {
		root := t.TempDir()
		for _, h := range []Harness{Claude, Codex} {
			p := previewFixture(t, root, h, ActionRemove)
			if p.Digest != "" {
				t.Fatalf("absent removal: %+v", p)
			}
			data := currentClaude(t)
			if h == Codex {
				data = RenderCodexRules()
			}
			writeFixture(t, root, p.File, data)
			p = previewFixture(t, root, h, ActionApply)
			if p.Digest != "" {
				t.Fatalf("current apply: %+v", p)
			}
		}
		if p := previewFixture(t, root, OpenCode, ActionApply); p.Digest != "" || p.State != NotNeeded || p.File != "" {
			t.Fatalf("OpenCode: %+v", p)
		}
	})
	t.Run("home_settings_ignored", func(t *testing.T) {
		home, root := t.TempDir(), t.TempDir()
		writeFixture(t, home, ".claude/settings.json", currentClaude(t))
		t.Setenv("HOME", home)
		p := previewFixture(t, root, Claude, ActionApply)
		if p.State != NotApplied || p.FileExists {
			t.Fatalf("home affected repository: %+v", p)
		}
	})
	t.Run("opencode_states", func(t *testing.T) {
		for _, tc := range []struct {
			data  string
			state State
		}{{`{}`, NotNeeded}, {`{"permission":"ask"}`, OwnRules}, {`{`, Unreadable}} {
			root := t.TempDir()
			writeFixture(t, root, "opencode.json", []byte(tc.data))
			p := previewFixture(t, root, OpenCode, ActionApply)
			if p.State != tc.state || p.Digest != "" || p.File != "" {
				t.Fatalf("%+v", p)
			}
		}
	})
}

func TestApply(t *testing.T) {
	t.Parallel()
	t.Run("writes_exactly_the_preview", func(t *testing.T) {
		for _, h := range []Harness{Claude, Codex} {
			root := t.TempDir()
			p := previewFixture(t, root, h, ActionApply)
			if _, err := Apply(root, h, ActionApply, p.Digest); err != nil {
				t.Fatal(err)
			}
			data := readFixture(t, root, p.File)
			if h == Codex {
				if !bytes.Equal(data, RenderCodexRules()) {
					t.Fatal("Codex bytes differ from rendering")
				}
			} else {
				allow, ask, err := ReadClaudeSettings(data)
				if err != nil || !slices.Equal(allow, p.Add) || !slices.Equal(ask, p.AddGuards) {
					t.Fatalf("written rules differ: %v, %v, %v", allow, ask, err)
				}
			}
			after := previewFixture(t, root, h, ActionApply)
			if after.State != Current || after.Digest != "" {
				t.Fatalf("after write: %+v", after)
			}
		}
	})
	t.Run("second_apply_writes_nothing", func(t *testing.T) {
		for _, h := range []Harness{Claude, Codex} {
			root := t.TempDir()
			p := previewFixture(t, root, h, ActionApply)
			if _, err := Apply(root, h, ActionApply, p.Digest); err != nil {
				t.Fatal(err)
			}
			target := filepath.Join(root, filepath.FromSlash(p.File))
			before, err := os.Stat(target)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := Apply(root, h, ActionApply, p.Digest); !errors.Is(err, ErrPreviewChanged) {
				t.Fatalf("second apply: %v", err)
			}
			after, err := os.Stat(target)
			if err != nil {
				t.Fatal(err)
			}
			if !os.SameFile(before, after) || !before.ModTime().Equal(after.ModTime()) {
				t.Fatal("second apply replaced file")
			}
		}
	})
	t.Run("stale_digest_writes_nothing", func(t *testing.T) {
		for _, h := range []Harness{Claude, Codex} {
			root := t.TempDir()
			file := ".claude/settings.json"
			data := []byte(`{}`)
			if h == Codex {
				file = ".codex/rules/roundfix-flow.rules"
				data = (list{entries: []Entry{Entries()[0]}}).renderCodexRules()
			}
			writeFixture(t, root, file, data)
			p := previewFixture(t, root, h, ActionApply)
			changed := append(slices.Clone(data), '\n')
			writeFixture(t, root, file, changed)
			if _, err := Apply(root, h, ActionApply, p.Digest); !errors.Is(err, ErrPreviewChanged) {
				t.Fatalf("stale apply: %v", err)
			}
			if !bytes.Equal(readFixture(t, root, file), changed) {
				t.Fatal("stale apply wrote")
			}
		}
	})
	t.Run("empty_digest_is_refused", func(t *testing.T) {
		for _, h := range []Harness{Claude, Codex} {
			root := t.TempDir()
			if _, err := Apply(root, h, ActionApply, ""); !errors.Is(err, ErrNotWritable) {
				t.Fatal(err)
			}
			entries, err := os.ReadDir(root)
			if err != nil || len(entries) != 0 {
				t.Fatalf("empty digest created files: %v, %v", entries, err)
			}
		}
	})
	t.Run("opencode_is_refused", func(t *testing.T) {
		root := t.TempDir()
		writeFixture(t, root, "opencode.json", []byte(`{"permission":"ask"}`))
		for _, action := range []Action{ActionApply, ActionRemove} {
			if _, err := Apply(root, OpenCode, action, "any-digest"); !errors.Is(err, ErrNotWritable) {
				t.Fatal(err)
			}
		}
		if string(readFixture(t, root, "opencode.json")) != `{"permission":"ask"}` {
			t.Fatal("OpenCode changed")
		}
	})
	t.Run("unreadable_target_is_not_written", func(t *testing.T) {
		for _, tc := range []struct {
			h          Harness
			file, data string
			cause      error
		}{
			{Claude, ".claude/settings.json", `{ // comment`, ErrNotStrictJSON},
			{Claude, ".claude/settings.json", `{"permissions":{"allow":false}}`, ErrUnexpectedShape},
			{Codex, ".codex/rules/roundfix-flow.rules", "other program", nil},
		} {
			root := t.TempDir()
			writeFixture(t, root, tc.file, []byte(tc.data))
			for _, action := range []Action{ActionApply, ActionRemove} {
				p, err := Apply(root, tc.h, action, "stale-digest")
				if !errors.Is(err, ErrUnreadable) || (tc.cause != nil && !errors.Is(err, tc.cause)) || p.Reason == "" {
					t.Fatalf("%+v, %v", p, err)
				}
				if string(readFixture(t, root, tc.file)) != tc.data {
					t.Fatal("unreadable target changed")
				}
			}
		}
	})
	t.Run("remove_claude_entries_only", func(t *testing.T) {
		root := t.TempDir()
		before := []byte("{\n  \"env\":{\"X\":\"Y\"}, \"permissions\":{\"allow\":[\"user\"],\"ask\":[\"user-ask\"],\"deny\":[\"deny\"]}\n}")
		writeFixture(t, root, ".claude/settings.json", before)
		writeFixture(t, root, ".claude/settings.local.json", []byte("local file"))
		p := previewFixture(t, root, Claude, ActionApply)
		if _, err := Apply(root, Claude, ActionApply, p.Digest); err != nil {
			t.Fatal(err)
		}
		removal := previewFixture(t, root, Claude, ActionRemove)
		if len(removal.Remove) != len(Entries())+len(Guards()) {
			t.Fatalf("remove list: %+v", removal)
		}
		if _, err := Apply(root, Claude, ActionRemove, removal.Digest); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(readFixture(t, root, p.File), before) || string(readFixture(t, root, ".claude/settings.local.json")) != "local file" {
			t.Fatal("user bytes changed")
		}
	})
	t.Run("remove_codex_file_only", func(t *testing.T) {
		root := t.TempDir()
		writeFixture(t, root, ".codex/rules/user.rules", []byte("user rules"))
		writeFixture(t, root, ".codex/config.toml", []byte("user config"))
		p := previewFixture(t, root, Codex, ActionApply)
		if _, err := Apply(root, Codex, ActionApply, p.Digest); err != nil {
			t.Fatal(err)
		}
		p = previewFixture(t, root, Codex, ActionRemove)
		if _, err := Apply(root, Codex, ActionRemove, p.Digest); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(p.File))); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("target still exists: %v", err)
		}
		if string(readFixture(t, root, ".codex/rules/user.rules")) != "user rules" || string(readFixture(t, root, ".codex/config.toml")) != "user config" {
			t.Fatal("other files changed")
		}
	})
	t.Run("link_outside_the_repository_is_refused", func(t *testing.T) {
		root, outside := t.TempDir(), t.TempDir()
		writeFixture(t, outside, "settings.json", []byte(`{}`))
		if err := os.Mkdir(filepath.Join(root, ".claude"), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(filepath.Join(outside, "settings.json"), filepath.Join(root, ".claude", "settings.json")); err != nil {
			t.Skipf("platform refuses symlink: %v", err)
		}
		if _, err := Apply(root, Claude, ActionApply, "any"); !errors.Is(err, ErrOutsideRepository) {
			t.Fatalf("%v", err)
		}
		if string(readFixture(t, outside, "settings.json")) != `{}` {
			t.Fatal("outside file changed")
		}
	})
	t.Run("directory_link_outside_the_repository_is_refused", func(t *testing.T) {
		root, outside := t.TempDir(), t.TempDir()
		writeFixture(t, outside, "settings.json", []byte(`{}`))
		if err := os.Symlink(outside, filepath.Join(root, ".claude")); err != nil {
			t.Skipf("platform refuses symlink: %v", err)
		}
		if _, err := Apply(root, Claude, ActionApply, "any"); !errors.Is(err, ErrOutsideRepository) {
			t.Fatalf("%v", err)
		}
		if string(readFixture(t, outside, "settings.json")) != `{}` {
			t.Fatal("outside file changed")
		}
	})
	t.Run("directories_created_only_on_write", func(t *testing.T) {
		for _, h := range []Harness{Claude, Codex} {
			root := t.TempDir()
			p := previewFixture(t, root, h, ActionApply)
			entries, err := os.ReadDir(root)
			if err != nil || len(entries) != 0 {
				t.Fatalf("preview created directories: %v, %v", entries, err)
			}
			if _, err := Apply(root, h, ActionApply, "wrong"); !errors.Is(err, ErrPreviewChanged) {
				t.Fatal(err)
			}
			entries, err = os.ReadDir(root)
			if err != nil || len(entries) != 0 {
				t.Fatalf("refused write created directories: %v, %v", entries, err)
			}
			if _, err := Apply(root, h, ActionApply, p.Digest); err != nil {
				t.Fatal(err)
			}
			if _, err := os.Stat(filepath.Dir(filepath.Join(root, filepath.FromSlash(p.File)))); err != nil {
				t.Fatal(err)
			}
		}
	})
}

func TestPermissionPlanEdges(t *testing.T) {
	t.Parallel()
	t.Run("only_retired_rule_still_proves_applied", func(t *testing.T) {
		root := t.TempDir()
		rules := shipped
		rules.retiredEntries = []string{"roundfix old"}
		writeFixture(t, root, ".claude/settings.json", settingsFixture(t, []string{"Bash(roundfix old *)"}, []string{}))
		p, err := plan(root, Claude, ActionApply, rules)
		if err != nil || p.State != Behind || !slices.Equal(p.Remove, []string{"Bash(roundfix old *)"}) {
			t.Fatalf("%+v, %v", p, err)
		}
	})
	t.Run("retired_rule_makes_complete_list_behind", func(t *testing.T) {
		root := t.TempDir()
		rules := shipped
		rules.retiredEntries = []string{"roundfix old"}
		rules.retiredClaudeRules = []string{"Bash(git old *--output*)"}
		var doc struct {
			Permissions struct {
				Allow []string `json:"allow"`
				Ask   []string `json:"ask"`
			} `json:"permissions"`
		}
		if err := json.Unmarshal(currentClaude(t), &doc); err != nil {
			t.Fatal(err)
		}
		doc.Permissions.Allow = append(doc.Permissions.Allow, "Bash(roundfix old *)")
		doc.Permissions.Ask = append(doc.Permissions.Ask, rules.retiredClaudeRules...)
		writeFixture(t, root, ".claude/settings.json", settingsFixture(t, doc.Permissions.Allow, doc.Permissions.Ask))
		p, err := plan(root, Claude, ActionApply, rules)
		if err != nil || p.State != Behind || len(p.Add)+len(p.AddGuards) != 0 || len(p.Remove) != 2 {
			t.Fatalf("%+v, %v", p, err)
		}
	})
	t.Run("removal_respects_list_roles", func(t *testing.T) {
		root := t.TempDir()
		allowRule, askRule := ClaudeRule(Entries()[0]), Guards()[0].ClaudeRule
		writeFixture(t, root, ".claude/settings.json", settingsFixture(t, []string{allowRule, askRule}, []string{askRule, allowRule}))
		p := previewFixture(t, root, Claude, ActionRemove)
		if _, err := Apply(root, Claude, ActionRemove, p.Digest); err != nil {
			t.Fatal(err)
		}
		allow, ask, err := ReadClaudeSettings(readFixture(t, root, p.File))
		if err != nil || !slices.Equal(allow, []string{askRule}) || !slices.Equal(ask, []string{allowRule}) {
			t.Fatalf("%v, %v, %v", allow, ask, err)
		}
	})
	t.Run("target_directory_is_unreadable", func(t *testing.T) {
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, ".claude", "settings.json"), 0755); err != nil {
			t.Fatal(err)
		}
		p, err := Apply(root, Claude, ActionApply, "any")
		if !errors.Is(err, ErrUnreadable) || p.State != Unreadable || p.Digest != "" {
			t.Fatalf("%+v, %v", p, err)
		}
	})
	t.Run("non_directory_parent_is_unreadable", func(t *testing.T) {
		root := t.TempDir()
		writeFixture(t, root, ".claude", []byte("not a directory"))
		p, err := Apply(root, Claude, ActionApply, "any")
		if !errors.Is(err, ErrUnreadable) || p.State != Unreadable || p.Digest != "" {
			t.Fatalf("%+v, %v", p, err)
		}
	})
	t.Run("codex_layout_behind_has_write_digest", func(t *testing.T) {
		root := t.TempDir()
		writeFixture(t, root, ".codex/rules/roundfix-flow.rules", append(RenderCodexRules(), '\n'))
		p := previewFixture(t, root, Codex, ActionApply)
		if p.State != Behind || p.Digest == "" || len(p.Add)+len(p.Remove) != 0 {
			t.Fatalf("%+v", p)
		}
		if _, err := Apply(root, Codex, ActionApply, p.Digest); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(readFixture(t, root, p.File), RenderCodexRules()) {
			t.Fatal("did not replace full rendering")
		}
	})
	t.Run("changed_list_changes_digest", func(t *testing.T) {
		root := t.TempDir()
		a, err := plan(root, Claude, ActionApply, shipped)
		if err != nil {
			t.Fatal(err)
		}
		changed := shipped
		changed.entries = changed.entries[1:]
		b, err := plan(root, Claude, ActionApply, changed)
		if err != nil || a.Digest == b.Digest {
			t.Fatalf("%s, %s, %v", a.Digest, b.Digest, err)
		}
	})
	t.Run("stale_removal_digest_keeps_user_edit", func(t *testing.T) {
		root := t.TempDir()
		writeFixture(t, root, ".claude/settings.json", currentClaude(t))
		p := previewFixture(t, root, Claude, ActionRemove)
		changed := append(currentClaude(t), ' ')
		writeFixture(t, root, p.File, changed)
		if _, err := Apply(root, Claude, ActionRemove, p.Digest); !errors.Is(err, ErrPreviewChanged) {
			t.Fatal(err)
		}
		if !bytes.Equal(readFixture(t, root, p.File), changed) {
			t.Fatal("stale removal wrote")
		}
	})
	t.Run("file_mode_preserved_and_no_temporary_residue", func(t *testing.T) {
		root := t.TempDir()
		writeFixture(t, root, ".claude/settings.json", []byte(`{}`))
		target := filepath.Join(root, ".claude", "settings.json")
		if err := os.Chmod(target, 0600); err != nil {
			t.Fatal(err)
		}
		before, err := os.Stat(target)
		if err != nil {
			t.Fatal(err)
		}
		p := previewFixture(t, root, Claude, ActionApply)
		if _, err := Apply(root, Claude, ActionApply, p.Digest); err != nil {
			t.Fatal(err)
		}
		after, err := os.Stat(target)
		if err != nil {
			t.Fatal(err)
		}
		if before.Mode().Perm() != after.Mode().Perm() {
			t.Fatalf("mode changed from %v to %v", before.Mode(), after.Mode())
		}
		entries, err := os.ReadDir(filepath.Dir(target))
		if err != nil || len(entries) != 1 || entries[0].Name() != "settings.json" {
			t.Fatalf("residue: %v, %v", entries, err)
		}
	})
}

func TestPermissionPathBoundary(t *testing.T) {
	t.Parallel()
	t.Run("valid_preview_cannot_follow_replaced_target", func(t *testing.T) {
		root, outside := t.TempDir(), t.TempDir()
		p := previewFixture(t, root, Claude, ActionApply)
		writeFixture(t, outside, "settings.json", []byte(`{}`))
		if err := os.Mkdir(filepath.Join(root, ".claude"), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(filepath.Join(outside, "settings.json"), filepath.Join(root, ".claude", "settings.json")); err != nil {
			t.Skipf("platform refuses symlink: %v", err)
		}
		if _, err := Apply(root, Claude, ActionApply, p.Digest); !errors.Is(err, ErrOutsideRepository) {
			t.Fatalf("%v", err)
		}
		if string(readFixture(t, outside, "settings.json")) != `{}` {
			t.Fatal("outside file changed")
		}
	})
	t.Run("codex_nested_directory_link_refused", func(t *testing.T) {
		root, outside := t.TempDir(), t.TempDir()
		p := previewFixture(t, root, Codex, ActionApply)
		if err := os.Mkdir(filepath.Join(root, ".codex"), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(outside, filepath.Join(root, ".codex", "rules")); err != nil {
			t.Skipf("platform refuses symlink: %v", err)
		}
		for _, action := range []Action{ActionApply, ActionRemove} {
			if _, err := Apply(root, Codex, action, p.Digest); !errors.Is(err, ErrOutsideRepository) {
				t.Fatalf("%s: %v", action, err)
			}
		}
		entries, err := os.ReadDir(outside)
		if err != nil || len(entries) != 0 {
			t.Fatalf("outside directory changed: %v, %v", entries, err)
		}
	})
	t.Run("in_repository_directory_link_allowed", func(t *testing.T) {
		root := t.TempDir()
		if err := os.Mkdir(filepath.Join(root, "config"), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(filepath.Join(root, "config"), filepath.Join(root, ".claude")); err != nil {
			t.Skipf("platform refuses symlink: %v", err)
		}
		p := previewFixture(t, root, Claude, ActionApply)
		if _, err := Apply(root, Claude, ActionApply, p.Digest); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(filepath.Join(root, "config", "settings.json")); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("in_repository_target_link_refused", func(t *testing.T) {
		root := t.TempDir()
		writeFixture(t, root, "config.json", []byte(`{}`))
		if err := os.Mkdir(filepath.Join(root, ".claude"), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(filepath.Join(root, "config.json"), filepath.Join(root, ".claude", "settings.json")); err != nil {
			t.Skipf("platform refuses symlink: %v", err)
		}
		if _, err := Apply(root, Claude, ActionApply, "any"); !errors.Is(err, ErrOutsideRepository) {
			t.Fatalf("%v", err)
		}
		if string(readFixture(t, root, "config.json")) != `{}` {
			t.Fatal("linked file changed")
		}
	})
	t.Run("opencode_outside_link_unreadable", func(t *testing.T) {
		root, outside := t.TempDir(), t.TempDir()
		writeFixture(t, outside, "config.json", []byte(`{"permission":"ask"}`))
		if err := os.Symlink(filepath.Join(outside, "config.json"), filepath.Join(root, "opencode.json")); err != nil {
			t.Skipf("platform refuses symlink: %v", err)
		}
		state, reason := ReadOpenCode(root)
		if state != Unreadable || !strings.Contains(reason, ErrOutsideRepository.Error()) {
			t.Fatalf("%s: %s", state, reason)
		}
	})
}
