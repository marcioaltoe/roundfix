package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/app"
	"roundfix/internal/flowpermissions"
	"roundfix/internal/gittest"
)

func permissionsWorkspace(t *testing.T) (string, string) {
	t.Helper()
	home, root := withCLIWorkspace(t)
	gittest.InitRepo(t, root)
	gittest.Run(t, root, "-c", "user.name=Permissions Test", "-c", "user.email=permissions@example.com", "commit", "--allow-empty", "-m", "test: permissions fixture")
	return home, root
}

func permissionsRun(t *testing.T, args ...string) (string, string, int) {
	t.Helper()
	var out, diagnostics bytes.Buffer
	code := runCLI(t, append([]string{"permissions"}, args...), &out, &diagnostics)
	return out.String(), diagnostics.String(), code
}

func permissionsAssert(t *testing.T, wantOut, wantErr string, wantCode int, args ...string) {
	t.Helper()
	out, diagnostics, code := permissionsRun(t, args...)
	if out != wantOut || diagnostics != wantErr || code != wantCode {
		t.Fatalf("%v: exit=%d (want %d)\nstdout=%q\nwant stdout=%q\nstderr=%q\nwant stderr=%q", args, code, wantCode, out, wantOut, diagnostics, wantErr)
	}
}

// The complete oracle is independent of the production text renderer.
func permissionsClaudeTranscript(p flowpermissions.Plan) string {
	var out strings.Builder
	fmt.Fprintf(&out, "Flow Command List for Claude Code, from Roundfix %s\nFile: .claude/settings.json (absent, it would be created)\nState: not applied\n\n", app.Version)
	fmt.Fprintf(&out, "Would add under permissions.allow (%d):\n", len(flowpermissions.Entries()))
	for _, entry := range flowpermissions.Entries() {
		fmt.Fprintf(&out, "  Bash(%s *)\n", strings.Join(entry.Words, " "))
	}
	fmt.Fprintf(&out, "\nWould add under permissions.ask, to keep these options asked about (%d):\n", len(flowpermissions.Guards()))
	for _, guard := range flowpermissions.Guards() {
		fmt.Fprintf(&out, "  %s\n", guard.ClaudeRule)
	}
	out.WriteString("  These ask in every permission mode, a session that skips prompts included.\n\nLeft out, and still asked about:\n")
	out.WriteString("  roundfix permissions: the command that applies this list\n")
	out.WriteString("  roundfix baseline, profiles, release, skills: the Roundfix skill tells a session to ask you first\n")
	out.WriteString("  git -C, config, init, merge, mv, push, remote, repack, switch, worktree: they change the repository or its configuration\n")
	out.WriteString("  gh auth: it changes your credentials\n")
	out.WriteString("  A commit, a push, a merge and a Pull Request stay asked about, and a delivery's push and merge stay behind the Spec's authorization.\n")
	out.WriteString("\nWhat this does not do:\n  It is not a security boundary: a roundfix command runs what this repository declares.\n  A listed read accepts any option; only the options listed above stay asked about.\n  The rules apply only after you trust this folder in Claude Code.\n")
	fmt.Fprintf(&out, "\nNothing was written. To write exactly this: roundfix permissions --harness claude --apply %s\n", p.Digest)
	return out.String()
}

func TestPermissionsSurfaceTranscripts(t *testing.T) {
	t.Parallel()
	t.Run("help", func(t *testing.T) {
		permissionsWorkspace(t)
		permissionsAssert(t, `Usage:
  roundfix permissions [--harness <claude|codex|opencode>] [--remove] [--apply <digest>] [--format <text|json>]

Shows the Flow Command List in the rule syntax of an Agent harness and writes
it into this repository's permission file for that harness only when it is
given the digest a preview printed. The list holds the roundfix subcommands
the owned skills name and the Git and GitHub reads they name. It is not a
security boundary. With no --harness it previews every harness whose files
this repository holds. Nothing under your home directory is read or written.

Options:
  --harness  The harness: claude, codex or opencode (reported, never written)
  --remove   Preview, and with --apply perform, the removal of what Roundfix wrote
  --apply    Write exactly what the preview with this digest showed
  --format   Output format: text (default) or json

Exit codes:
  0  a preview was printed, or the plan was written, or nothing needed writing
  1  the permission file could not be read as the harness documents it, or not written
  2  invalid input
  3  the preview changed since the digest was printed; nothing was written
`, "", 0, "--help")
	})
	t.Run("claude_preview", func(t *testing.T) {
		_, root := permissionsWorkspace(t)
		p := permissionsPlan(t, root, flowpermissions.Claude, false)
		before := snapshotDirectoryFiles(t, root)
		permissionsAssert(t, permissionsClaudeTranscript(p), "", 0, "--harness", "claude")
		if !reflect.DeepEqual(before, snapshotDirectoryFiles(t, root)) {
			t.Fatal("preview changed repository")
		}
	})
	t.Run("stale_digest", func(t *testing.T) {
		permissionsWorkspace(t)
		permissionsAssert(t, "", "The preview changed since digest 000000000000 was printed. Nothing was written.\nRun the preview again: roundfix permissions --harness claude\n", 3, "--harness", "claude", "--apply", "000000000000")
	})
	t.Run("opencode_report", func(t *testing.T) {
		permissionsWorkspace(t)
		permissionsAssert(t, "Flow Command List for OpenCode, from Roundfix "+app.Version+"\nState: not needed\nOpenCode's base policy allows every shell command, and this repository holds no OpenCode configuration that asks about or denies one.\nRoundfix writes nothing for OpenCode.\n", "", 0, "--harness", "opencode")
	})
}

func permissionsPlan(t *testing.T, root string, h flowpermissions.Harness, remove bool) flowpermissions.Plan {
	t.Helper()
	action := flowpermissions.ActionApply
	if remove {
		action = flowpermissions.ActionRemove
	}
	p, err := flowpermissions.Preview(root, h, action)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func permissionsWrite(t *testing.T, path, content string) {
	t.Helper()
	mustMkdir(t, filepath.Dir(path))
	mustWrite(t, path, content)
}

func permissionsRead(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func permissionsJSON(t *testing.T, args ...string) map[string]any {
	t.Helper()
	out, diagnostics, code := permissionsRun(t, append(args, "--format", "json")...)
	if code != 0 || diagnostics != "" {
		t.Fatalf("exit=%d stderr=%s stdout=%s", code, diagnostics, out)
	}
	var report map[string]any
	if err := json.Unmarshal([]byte(out), &report); err != nil {
		t.Fatal(err)
	}
	if report["schemaVersion"] != "roundfix/flow-permissions/v1" || report["roundfix"] != app.Version || len(report) != 3 {
		t.Fatalf("unexpected report: %#v", report)
	}
	return report["harnesses"].([]any)[0].(map[string]any)
}

func TestPermissionsCommand(t *testing.T) {
	t.Parallel()
	t.Run("no_harness_files", func(t *testing.T) {
		permissionsWorkspace(t)
		permissionsAssert(t, "No harness files in this repository. Name one with --harness.\n", "", 0)
	})
	t.Run("preview_every_present_harness", func(t *testing.T) {
		home, root := permissionsWorkspace(t)
		permissionsWrite(t, filepath.Join(root, ".claude/settings.json"), "{}")
		permissionsWrite(t, filepath.Join(root, ".codex/config.toml"), "")
		permissionsWrite(t, filepath.Join(root, "opencode.json"), "{}")
		permissionsWrite(t, filepath.Join(home, ".roundfixrc.yml"), "[invalid")
		permissionsWrite(t, filepath.Join(root, ".roundfixrc.yml"), "[invalid")
		mustMkdir(t, filepath.Join(root, "nested"))
		setCommandWorkDirForTest(t, filepath.Join(root, "nested"))
		before := snapshotDirectoryFiles(t, root)
		out, diagnostics, code := permissionsRun(t)
		if code != 0 || diagnostics != "" || strings.Count(out, "Flow Command List for ") != 3 {
			t.Fatalf("exit=%d stderr=%s stdout=%s", code, diagnostics, out)
		}
		if !(strings.Index(out, "for Claude Code") < strings.Index(out, "for Codex") && strings.Index(out, "for Codex") < strings.Index(out, "for OpenCode")) {
			t.Fatal("wrong harness order")
		}
		if !reflect.DeepEqual(before, snapshotDirectoryFiles(t, root)) {
			t.Fatal("preview wrote files")
		}
	})
	t.Run("json_preview_schema", func(t *testing.T) {
		_, root := permissionsWorkspace(t)
		p := permissionsPlan(t, root, flowpermissions.Claude, false)
		r := permissionsJSON(t, "--harness", "claude")
		want := map[string]any{"harness": "claude", "action": "apply", "file": ".claude/settings.json", "fileExists": false, "state": "not_applied", "add": p.Add, "addGuards": p.AddGuards, "remove": []string{}, "digest": p.Digest, "next": "roundfix permissions --harness claude --apply " + p.Digest}
		var exclusions []map[string]string
		for _, ex := range flowpermissions.Exclusions() {
			exclusions = append(exclusions, map[string]string{"command": strings.Join(ex.Words, " "), "reason": ex.Reason})
		}
		want["exclusions"] = exclusions
		data, err := json.Marshal(want)
		if err != nil {
			t.Fatal(err)
		}
		var normalized map[string]any
		if err := json.Unmarshal(data, &normalized); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(r, normalized) {
			t.Fatalf("got %#v, want %#v", r, normalized)
		}
	})
	t.Run("apply_writes_the_preview", func(t *testing.T) {
		_, root := permissionsWorkspace(t)
		p := permissionsPlan(t, root, flowpermissions.Claude, false)
		permissionsAssert(t, permissionsClaudeTranscript(p), "", 0, "--harness", "claude")
		want := fmt.Sprintf("Wrote %d entries and %d Guard Entries to .claude/settings.json.\nThe rules apply only after you trust this folder in Claude Code.\nTo remove them: roundfix permissions --harness claude --remove\n", len(p.Add), len(p.AddGuards))
		permissionsAssert(t, want, "", 0, "--harness", "claude", "--apply", p.Digest)
		allow, ask, err := flowpermissions.ReadClaudeSettings(permissionsRead(t, filepath.Join(root, p.File)))
		if err != nil || !reflect.DeepEqual(allow, p.Add) || !reflect.DeepEqual(ask, p.AddGuards) {
			t.Fatalf("allow=%v ask=%v err=%v", allow, ask, err)
		}
	})
	t.Run("apply_json_written", func(t *testing.T) {
		_, root := permissionsWorkspace(t)
		before := permissionsJSON(t, "--harness", "claude")
		p := permissionsPlan(t, root, flowpermissions.Claude, false)
		after := permissionsJSON(t, "--harness", "claude", "--apply", p.Digest)
		before["written"] = true
		if !reflect.DeepEqual(before, after) {
			t.Fatalf("before=%v after=%v", before, after)
		}
		assertPathExists(t, filepath.Join(root, p.File))
	})
	t.Run("second_preview_is_current", func(t *testing.T) {
		_, root := permissionsWorkspace(t)
		p := permissionsPlan(t, root, flowpermissions.Claude, false)
		permissionsJSON(t, "--harness", "claude", "--apply", p.Digest)
		before := permissionsRead(t, filepath.Join(root, p.File))
		r := permissionsJSON(t, "--harness", "claude")
		if r["state"] != "current" || r["digest"] != "" || r["next"] != "" {
			t.Fatalf("%v", r)
		}
		out, diagnostics, code := permissionsRun(t, "--harness", "claude")
		if code != 0 || diagnostics != "" || !strings.HasSuffix(out, "Nothing to write: the file holds the installed list.\n") {
			t.Fatalf("%d %s %s", code, diagnostics, out)
		}
		if !bytes.Equal(before, permissionsRead(t, filepath.Join(root, p.File))) {
			t.Fatal("current preview wrote file")
		}
	})
	for _, tc := range []struct {
		name   string
		args   []string
		reason string
	}{
		{"apply_without_harness_is_invalid", []string{"--apply", "anything"}, "--apply requires --harness"},
		{"apply_opencode_is_invalid", []string{"--harness", "opencode", "--apply", "anything"}, "Roundfix writes nothing for OpenCode"},
		{"unknown_harness_is_invalid", []string{"--harness", "other"}, "--harness must be claude, codex or opencode"},
		{"unknown_format_is_invalid", []string{"--format", "yaml"}, "--format must be text or json"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			permissionsWorkspace(t)
			permissionsAssert(t, "", "roundfix: permissions: "+tc.reason+"\n", 2, tc.args...)
		})
	}
	t.Run("outside_a_git_repository", func(t *testing.T) {
		setCommandEnvironmentForTest(t, t.TempDir(), t.TempDir())
		permissionsAssert(t, "", "roundfix: permissions: not inside a Git repository\n", 2, "--harness", "claude")
	})
	t.Run("stale_digest_after_file_change", func(t *testing.T) {
		_, root := permissionsWorkspace(t)
		p := permissionsPlan(t, root, flowpermissions.Claude, false)
		permissionsWrite(t, filepath.Join(root, p.File), "{}\n")
		permissionsAssert(t, "", "The preview changed since digest "+p.Digest+" was printed. Nothing was written.\nRun the preview again: roundfix permissions --harness claude\n", 3, "--harness", "claude", "--apply", p.Digest)
		if string(permissionsRead(t, filepath.Join(root, p.File))) != "{}\n" {
			t.Fatal("stale apply wrote file")
		}
	})
	t.Run("unreadable_target_on_write", func(t *testing.T) {
		_, root := permissionsWorkspace(t)
		permissionsWrite(t, filepath.Join(root, ".claude/settings.json"), "{broken")
		p := permissionsPlan(t, root, flowpermissions.Claude, false)
		for _, digest := range []string{"", "000000000000"} {
			permissionsAssert(t, "", ".claude/settings.json cannot be read as Claude Code documents it: "+p.Reason+". Nothing was written.\n", 1, "--harness", "claude", "--apply", digest)
		}
		if string(permissionsRead(t, filepath.Join(root, p.File))) != "{broken" {
			t.Fatal("unreadable apply wrote file")
		}
	})
	t.Run("link_outside_the_repository", func(t *testing.T) {
		_, root := permissionsWorkspace(t)
		outside := t.TempDir()
		permissionsWrite(t, filepath.Join(outside, "settings.json"), "{}")
		if err := os.Symlink(outside, filepath.Join(root, ".claude")); err != nil {
			t.Skipf("platform refuses link: %v", err)
		}
		for _, args := range [][]string{{"--harness", "claude"}, {"--harness", "claude", "--apply", "000000000000"}} {
			permissionsAssert(t, "", ".claude/settings.json leaves the repository through a link. Nothing was written.\n", 1, args...)
		}
		if string(permissionsRead(t, filepath.Join(outside, "settings.json"))) != "{}" {
			t.Fatal("link target changed")
		}
	})
	t.Run("remove_preview_and_apply", func(t *testing.T) {
		_, root := permissionsWorkspace(t)
		const original = "{\n  \"permissions\": {\"allow\": [\"Bash(custom *)\"], \"ask\": [\"Bash(other *)\"]},\n  \"other\": true\n}\n"
		permissionsWrite(t, filepath.Join(root, ".claude/settings.json"), original)
		p := permissionsPlan(t, root, flowpermissions.Claude, false)
		permissionsJSON(t, "--harness", "claude", "--apply", p.Digest)
		p = permissionsPlan(t, root, flowpermissions.Claude, true)
		var want strings.Builder
		fmt.Fprintf(&want, "Flow Command List for Claude Code, from Roundfix %s\nFile: .claude/settings.json\n\nWould remove (%d):\n", app.Version, len(p.Remove))
		for _, rule := range p.Remove {
			fmt.Fprintf(&want, "  %s\n", rule)
		}
		fmt.Fprintf(&want, "An entry with the same text that you wrote yourself is removed with them.\n\nNothing was written. To remove exactly this: roundfix permissions --harness claude --remove --apply %s\n", p.Digest)
		permissionsAssert(t, want.String(), "", 0, "--harness", "claude", "--remove")
		permissionsAssert(t, fmt.Sprintf("Removed %d entries from .claude/settings.json.\n", len(p.Remove)), "", 0, "--harness", "claude", "--remove", "--apply", p.Digest)
		if string(permissionsRead(t, filepath.Join(root, p.File))) != original {
			t.Fatal("removal did not restore original bytes")
		}
	})
	t.Run("nothing_to_remove", func(t *testing.T) {
		permissionsWorkspace(t)
		permissionsAssert(t, "Flow Command List for Claude Code, from Roundfix "+app.Version+"\nFile: .claude/settings.json (absent, it would be created)\nNothing to remove: the file holds no entry of the list.\n", "", 0, "--harness", "claude", "--remove")
		r := permissionsJSON(t, "--harness", "claude", "--remove")
		if r["digest"] != "" || r["next"] != "" {
			t.Fatalf("%v", r)
		}
	})
	t.Run("codex_apply_and_remove", func(t *testing.T) {
		_, root := permissionsWorkspace(t)
		permissionsWrite(t, filepath.Join(root, ".codex/rules/user.rules"), "user rule\n")
		p := permissionsPlan(t, root, flowpermissions.Codex, false)
		permissionsAssert(t, fmt.Sprintf("Wrote .codex/rules/roundfix-flow.rules with %d entries.\nThe rules apply only after you trust this project in Codex.\nTo remove them: roundfix permissions --harness codex --remove\n", len(p.Add)), "", 0, "--harness", "codex", "--apply", p.Digest)
		if !bytes.Equal(permissionsRead(t, filepath.Join(root, p.File)), flowpermissions.RenderCodexRules()) {
			t.Fatal("wrong Codex bytes")
		}
		p = permissionsPlan(t, root, flowpermissions.Codex, true)
		out, diagnostics, code := permissionsRun(t, "--harness", "codex", "--remove")
		if code != 0 || diagnostics != "" || !strings.Contains(out, "Would delete .codex/rules/roundfix-flow.rules.\n") {
			t.Fatalf("%d %s %s", code, diagnostics, out)
		}
		permissionsAssert(t, "Deleted .codex/rules/roundfix-flow.rules.\n", "", 0, "--harness", "codex", "--remove", "--apply", p.Digest)
		assertPathMissing(t, filepath.Join(root, p.File))
		if string(permissionsRead(t, filepath.Join(root, ".codex/rules/user.rules"))) != "user rule\n" {
			t.Fatal("user rules changed")
		}
	})
	t.Run("codex_preview_says_no_guard", func(t *testing.T) {
		permissionsWorkspace(t)
		out, diagnostics, code := permissionsRun(t, "--harness", "codex")
		for _, phrase := range []string{"No option of a listed read can be kept asked about in Codex.", "In Codex an allowed command runs outside the sandbox with nobody asked.", "git show --output=<file> would write that file."} {
			if code != 0 || diagnostics != "" || !strings.Contains(out, phrase) {
				t.Fatalf("missing %s: %d %s %s", phrase, code, diagnostics, out)
			}
		}
	})
	t.Run("exclusions_printed", func(t *testing.T) {
		_, root := permissionsWorkspace(t)
		permissionsAssert(t, permissionsClaudeTranscript(permissionsPlan(t, root, flowpermissions.Claude, false)), "", 0, "--harness", "claude")
	})
	t.Run("home_files_untouched", func(t *testing.T) {
		home, root := permissionsWorkspace(t)
		permissionsWrite(t, filepath.Join(home, ".claude/settings.json"), "{personal settings}\n")
		permissionsWrite(t, filepath.Join(home, ".codex/rules/default.rules"), "personal rules\n")
		before := snapshotDirectoryFiles(t, home)
		for _, h := range []flowpermissions.Harness{flowpermissions.Claude, flowpermissions.Codex} {
			p := permissionsPlan(t, root, h, false)
			permissionsJSON(t, "--harness", string(h))
			permissionsJSON(t, "--harness", string(h), "--apply", p.Digest)
			p = permissionsPlan(t, root, h, true)
			permissionsJSON(t, "--harness", string(h), "--remove", "--apply", p.Digest)
		}
		if !reflect.DeepEqual(before, snapshotDirectoryFiles(t, home)) {
			t.Fatal("home files changed")
		}
	})
}

func TestPermissionsTextPreviewNamesEveryExclusion(t *testing.T) {
	t.Parallel()
	for _, harness := range []string{"claude", "codex"} {
		t.Run(harness, func(t *testing.T) {
			permissionsWorkspace(t)
			text, diagnostics, code := permissionsRun(t, "--harness", harness)
			if code != 0 || diagnostics != "" {
				t.Fatalf("text preview: exit=%d stderr=%q", code, diagnostics)
			}
			textExclusions := permissionsTextExclusions(t, text)
			jsonExclusions := permissionsJSONExclusions(t, permissionsJSON(t, "--harness", harness))
			if len(textExclusions) == 0 {
				t.Fatal("text preview names no Named Exclusion")
			}
			for command, reason := range textExclusions {
				jsonReason, ok := jsonExclusions[command]
				if !ok {
					t.Errorf("text preview names %q, JSON does not", command)
					continue
				}
				if jsonReason != reason {
					t.Errorf("%q: text reason %q, JSON reason %q", command, reason, jsonReason)
				}
			}
			for command := range jsonExclusions {
				if _, ok := textExclusions[command]; !ok {
					t.Errorf("JSON names %q, text preview does not", command)
				}
			}
		})
	}
}

func permissionsTextExclusions(t *testing.T, preview string) map[string]string {
	t.Helper()
	const heading = "Left out, and still asked about:"
	const fixed = "  A commit, a push, a merge and a Pull Request stay asked about"
	lines := strings.Split(preview, "\n")
	start := -1
	for i, line := range lines {
		if line == heading {
			start = i
			break
		}
	}
	if start < 0 {
		t.Fatalf("text preview has no %q block", heading)
	}
	exclusions := map[string]string{}
	for _, line := range lines[start+1:] {
		if strings.HasPrefix(line, fixed) {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		reasonAt := strings.LastIndex(line, ": ")
		if reasonAt < 0 {
			t.Fatalf("exclusion line %q has no reason", line)
		}
		commands := strings.Split(line[:reasonAt], ", ")
		reason := line[reasonAt+2:]
		firstWord := strings.Fields(commands[0])[0]
		for i, command := range commands {
			if i > 0 && strings.Fields(command)[0] != firstWord {
				command = firstWord + " " + command
			}
			exclusions[command] = reason
		}
	}
	return exclusions
}

func permissionsJSONExclusions(t *testing.T, report map[string]any) map[string]string {
	t.Helper()
	raw, ok := report["exclusions"].([]any)
	if !ok {
		t.Fatalf("JSON report has no exclusions: %v", report)
	}
	exclusions := map[string]string{}
	for _, item := range raw {
		exclusion, ok := item.(map[string]any)
		if !ok {
			t.Fatalf("unexpected exclusion: %v", item)
		}
		command, _ := exclusion["command"].(string)
		reason, _ := exclusion["reason"].(string)
		if command == "" {
			t.Fatalf("exclusion has no command: %v", exclusion)
		}
		exclusions[command] = reason
	}
	return exclusions
}

func TestRootHelpNamesThePermissionsCommand(t *testing.T) {
	t.Parallel()
	var out, diagnostics bytes.Buffer
	code := runCLI(t, []string{"--help"}, &out, &diagnostics)
	for _, line := range []string{
		"  roundfix update check [--refresh] [--format <text|json>]\n  roundfix permissions [--harness <claude|codex|opencode>] [--remove] [--apply <digest>] [--format <text|json>]\n",
		"  update     Report a newer release and this repository's Baseline plan\n  permissions  Show or apply the Flow Command List for an Agent harness\n",
	} {
		if code != 0 || diagnostics.Len() != 0 || !strings.Contains(out.String(), line) {
			t.Fatalf("missing %q: exit=%d stdout=%s stderr=%s", line, code, &out, &diagnostics)
		}
	}
}

func TestPermissionsRefusalBoundaries(t *testing.T) {
	t.Parallel()
	t.Run("empty_digest_writes_nothing", func(t *testing.T) {
		_, root := permissionsWorkspace(t)
		for _, remove := range []bool{false, true} {
			args := []string{"--harness", "claude", "--apply="}
			next := "roundfix permissions --harness claude"
			if remove {
				args = append(args, "--remove")
				next += " --remove"
			}
			permissionsAssert(t, "", "The preview changed since digest  was printed. Nothing was written.\nRun the preview again: "+next+"\n", 3, args...)
		}
		assertPathMissing(t, filepath.Join(root, ".claude"))
	})
	t.Run("opencode_own_rules_and_json", func(t *testing.T) {
		_, root := permissionsWorkspace(t)
		permissionsWrite(t, filepath.Join(root, "opencode.json"), `{"permission":{"bash":"ask"}}`)
		var want strings.Builder
		fmt.Fprintf(&want, "Flow Command List for OpenCode, from Roundfix %s\nState: own rules\nThis repository's OpenCode configuration (opencode.json) holds rules that ask about or deny shell commands. Compare them with the list:\n", app.Version)
		for _, entry := range flowpermissions.Entries() {
			fmt.Fprintf(&want, "  %s\n", strings.Join(entry.Words, " "))
		}
		want.WriteString("Roundfix writes nothing for OpenCode.\n")
		permissionsAssert(t, want.String(), "", 0, "--harness", "opencode")
		r := permissionsJSON(t, "--harness", "opencode")
		if r["state"] != "own_rules" || r["digest"] != "" || r["next"] != "" || r["file"] != "" {
			t.Fatalf("unexpected report: %v", r)
		}
		for _, key := range []string{"add", "addGuards", "remove"} {
			if values, ok := r[key].([]any); !ok || len(values) != 0 {
				t.Errorf("%s is not an empty array: %v", key, r[key])
			}
		}
	})
	t.Run("codex_rewrite_reports_all_entries", func(t *testing.T) {
		_, root := permissionsWorkspace(t)
		entry := flowpermissions.Entries()[0].Words[1]
		older := strings.Replace(string(flowpermissions.RenderCodexRules()), fmt.Sprintf("%q, ", entry), "", 1)
		permissionsWrite(t, filepath.Join(root, ".codex/rules/roundfix-flow.rules"), older)
		p := permissionsPlan(t, root, flowpermissions.Codex, false)
		if p.State != flowpermissions.Behind || len(p.Add) != 1 {
			t.Fatalf("unexpected behind plan: %+v", p)
		}
		want := fmt.Sprintf("Wrote .codex/rules/roundfix-flow.rules with %d entries.\nThe rules apply only after you trust this project in Codex.\nTo remove them: roundfix permissions --harness codex --remove\n", len(flowpermissions.Entries()))
		permissionsAssert(t, want, "", 0, "--harness", "codex", "--apply", p.Digest)
		if !bytes.Equal(permissionsRead(t, filepath.Join(root, p.File)), flowpermissions.RenderCodexRules()) {
			t.Fatal("Codex rewrite is incomplete")
		}
	})
}
