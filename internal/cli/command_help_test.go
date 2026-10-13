package cli

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestHelpSurfaceTranscripts(t *testing.T) {
	t.Parallel()
	setCommandEnvironmentForTest(t, t.TempDir(), t.TempDir())
	const qaHelp = `Usage:
  roundfix qa-report accept <path>

Applies the declared-acceptance policy to the selected QA Report. The report's
Spec directory supplies its declarations. Produces no stdout; exit status is
the machine-readable result.

Exit codes:
  0  the report is acceptable
  1  the report is missing, unreadable, or unacceptable
  2  usage error
`
	tests := []struct {
		name   string
		args   []string
		stdout string
	}{
		{"transcript 1", []string{"qa-report", "accept", "--help"}, qaHelp},
		{"transcript 2", []string{"qa-report", "accept", "-h"}, qaHelp},
		{"transcript 3", []string{"qa-report", "accept", "docs/specs/absent/qa/qa-report-2026-10-12.md", "--help"}, qaHelp},
		{"transcript 4", []string{"baseline", "capabilities", "--help"}, commandUsage("baseline")},
		{"transcript 5", []string{"baseline", "skills", "-h"}, commandUsage("baseline")},
		{"transcript 6", []string{"baseline", "assets", "--help"}, commandUsage("baseline")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := runCLI(t, tt.args, &stdout, &stderr)
			if code != exitOK || stdout.String() != tt.stdout || stderr.Len() != 0 {
				t.Errorf("roundfix %s: exit = %d, want 0; stdout = %q, want %q; stderr = %q, want empty", strings.Join(tt.args, " "), code, stdout.String(), tt.stdout, stderr.String())
			}
		})
	}
}

func TestEveryCommandPathPrintsItsHelp(t *testing.T) {
	t.Parallel()
	setCommandEnvironmentForTest(t, t.TempDir(), t.TempDir())
	var stdout, stderr bytes.Buffer
	if code := runCLI(t, []string{"--help"}, &stdout, &stderr); code != exitOK || stderr.Len() != 0 {
		t.Fatalf("roundfix --help: exit = %d; stderr = %q", code, stderr.String())
	}
	paths := helpCommandPaths(stdout.String())
	if len(paths) < 60 {
		t.Fatalf("root help yielded %d command paths, want at least 60: %v", len(paths), paths)
	}
	t.Logf("checking %d command paths with both help flags (%d calls)", len(paths), 2*len(paths))
	for _, path := range paths {
		for _, help := range []string{"--help", "-h"} {
			args := append(strings.Fields(path), help)
			t.Run(strings.Join(args, " "), func(t *testing.T) {
				var stdout, stderr bytes.Buffer
				code := runCLI(t, args, &stdout, &stderr)
				if code != exitOK || !strings.HasPrefix(stdout.String(), "Usage:") || stderr.Len() != 0 {
					firstLine, _, _ := strings.Cut(stdout.String(), "\n")
					t.Errorf("roundfix %s: help exited %d; first stdout line = %q; stderr = %q", strings.Join(args, " "), code, firstLine, stderr.String())
				}
			})
		}
	}
}

// helpCommandPaths returns every command path the root usage names and every
// group prefix of one, in first-seen order without duplicates.
func helpCommandPaths(rootUsage string) []string {
	var paths []string
	seen := make(map[string]bool)
	add := func(words []string) {
		for count := 1; count <= len(words); count++ {
			path := strings.Join(words[:count], " ")
			if !seen[path] {
				seen[path] = true
				paths = append(paths, path)
			}
		}
	}
	for _, line := range strings.Split(rootUsage, "\n") {
		command, found := strings.CutPrefix(line, "  roundfix ")
		if !found {
			continue
		}
		var words []string
		for _, token := range strings.Fields(command) {
			if strings.HasPrefix(token, "-") || strings.HasPrefix(token, "[") || strings.HasPrefix(token, "(") || strings.HasPrefix(token, "<") {
				if len(words) > 0 && strings.HasPrefix(token, "<") && strings.HasSuffix(token, ">") && strings.Contains(token, "|") {
					for _, alternative := range strings.Split(token[1:len(token)-1], "|") {
						add(append(append([]string(nil), words...), alternative))
					}
					words = nil
				}
				break
			}
			words = append(words, token)
		}
		add(words)
	}
	return paths
}

func TestHelpCommandPathsKeepsPrefixesAndAlternatives(t *testing.T) {
	t.Parallel()
	rootUsage := `Usage:
  roundfix window <set|show|clear>
  roundfix archive <slug>
  roundfix baseline capabilities check [--profile <id>]
  roundfix archive <slug> [--qa-override]
  roundfix --help
roundfix ignored
  roundfix events --run <id>
  roundfix doctor (read-only)
  roundfix short -x
`
	want := []string{"window", "window set", "window show", "window clear", "archive", "baseline", "baseline capabilities", "baseline capabilities check", "events", "doctor", "short"}
	if got := helpCommandPaths(rootUsage); !reflect.DeepEqual(got, want) {
		t.Fatalf("helpCommandPaths() = %v, want %v", got, want)
	}
}

func TestCommandHelpArgumentBoundaries(t *testing.T) {
	t.Parallel()
	setCommandEnvironmentForTest(t, t.TempDir(), t.TempDir())
	const qaUsageError = "roundfix: qa-report accept requires exactly one report path\nRun 'roundfix qa-report --help' for usage.\n"
	type testCase struct {
		name   string
		args   []string
		code   int
		stdout string
		stderr string
	}
	tests := []testCase{
		{"report trailing short help", []string{"qa-report", "accept", "missing.md", "-h"}, exitOK, qaReportUsage, ""},
		{"report help before operands", []string{"qa-report", "accept", "--help", "one", "two"}, exitOK, qaReportUsage, ""},
		{"report help before terminator", []string{"qa-report", "accept", "-h", "--", "one"}, exitOK, qaReportUsage, ""},
		{"report terminator before help", []string{"qa-report", "accept", "--", "--help"}, exitPreflight, "", qaUsageError},
		{"report terminator after path", []string{"qa-report", "accept", "missing.md", "--", "-h"}, exitPreflight, "", qaUsageError},
		{"report no path", []string{"qa-report", "accept"}, exitPreflight, "", qaUsageError},
		{"report blank path", []string{"qa-report", "accept", " "}, exitPreflight, "", qaUsageError},
		{"report two paths", []string{"qa-report", "accept", "one", "two"}, exitPreflight, "", qaUsageError},
	}
	for _, group := range []struct {
		name    string
		command string
	}{
		{"capabilities", "check"},
		{"skills", "restore or reconcile"},
		{"assets", "sync"},
	} {
		for _, help := range []string{"--help", "-h", "help"} {
			tests = append(tests, testCase{group.name + " " + help, []string{"baseline", group.name, help}, exitOK, commandUsage("baseline"), ""})
		}
		for _, suffix := range [][]string{nil, {"unknown"}, {"unknown", "--help"}, {"--", "--help"}} {
			args := append([]string{"baseline", group.name}, suffix...)
			tests = append(tests, testCase{strings.Join(args, " "), args, exitPreflight, "", "roundfix: baseline baseline " + group.name + " failed: baseline " + group.name + " requires the " + group.command + " command\nRun 'roundfix baseline profile --help' for usage.\n"})
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := runCLI(t, tt.args, &stdout, &stderr)
			if code != tt.code || stdout.String() != tt.stdout || stderr.String() != tt.stderr {
				t.Errorf("roundfix %s: exit = %d, want %d; stdout = %q, want %q; stderr = %q, want %q", strings.Join(tt.args, " "), code, tt.code, stdout.String(), tt.stdout, stderr.String(), tt.stderr)
			}
		})
	}
}
