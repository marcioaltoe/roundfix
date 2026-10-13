package cli

import (
	"bytes"
	"regexp"
	"sort"
	"strings"
	"testing"
)

func TestBaselineGroupUsageTranscripts(t *testing.T) {
	t.Parallel()
	setCommandEnvironmentForTest(t, t.TempDir(), t.TempDir())
	tests := []struct {
		name     string
		args     []string
		stderr   string
		helpPath string
	}{
		{
			name: "transcript_1",
			args: []string{"baseline", "capabilities"},
			stderr: "roundfix: baseline capabilities requires the check command\n" +
				"Run 'roundfix baseline capabilities --help' for usage.\n",
			helpPath: "baseline capabilities",
		},
		{
			name: "transcript_2",
			args: []string{"baseline", "skills"},
			stderr: "roundfix: baseline skills requires the restore or reconcile command\n" +
				"Run 'roundfix baseline skills --help' for usage.\n",
			helpPath: "baseline skills",
		},
		{
			name: "transcript_3",
			args: []string{"baseline", "skills", "prune"},
			stderr: "roundfix: baseline skills requires the restore or reconcile command\n" +
				"Run 'roundfix baseline skills --help' for usage.\n",
			helpPath: "baseline skills",
		},
		{
			name: "transcript_4",
			args: []string{"baseline", "assets"},
			stderr: "roundfix: baseline assets requires the sync command\n" +
				"Run 'roundfix baseline assets --help' for usage.\n",
			helpPath: "baseline assets",
		},
		{
			name: "transcript_5",
			args: []string{"baseline", "capabilities", "--format", "json"},
			stderr: "roundfix: baseline capabilities requires the check command\n" +
				"Run 'roundfix baseline capabilities --help' for usage.\n",
			helpPath: "baseline capabilities",
		},
		{
			name: "transcript_6",
			args: []string{"baseline", "prune"},
			stderr: "roundfix: unknown baseline command \"prune\"\n" +
				"Run 'roundfix baseline --help' for usage.\n",
			helpPath: "baseline",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := runCLI(t, tt.args, &stdout, &stderr)
			if code != exitPreflight || stdout.Len() != 0 || stderr.String() != tt.stderr {
				t.Errorf("roundfix %s: exit = %d, want 2; stdout = %q, want empty; stderr = %q, want %q", strings.Join(tt.args, " "), code, stdout.String(), stderr.String(), tt.stderr)
			}
			stdout.Reset()
			stderr.Reset()
			helpArgs := append(strings.Fields(tt.helpPath), "--help")
			if code := runCLI(t, helpArgs, &stdout, &stderr); code != exitOK || stderr.Len() != 0 {
				t.Errorf("pointer 'roundfix %s': exit = %d, want 0; stderr = %q, want empty", strings.Join(helpArgs, " "), code, stderr.String())
			}
		})
	}
}

func TestEveryGroupUsageErrorPointsAtItsCommands(t *testing.T) {
	t.Parallel()
	setCommandEnvironmentForTest(t, t.TempDir(), t.TempDir())
	var stdout, stderr bytes.Buffer
	if code := runCLI(t, []string{"--help"}, &stdout, &stderr); code != exitOK || stderr.Len() != 0 {
		t.Fatalf("roundfix --help: exit = %d, want 0; stderr = %q, want empty", code, stderr.String())
	}
	groups := usageErrorGroups(stdout.String())
	if len(groups) < 12 {
		t.Fatalf("root help yielded %d groups, want at least 12: %v", len(groups), groups)
	}
	t.Logf("checking usage errors for %d groups", len(groups))
	var paths []string
	for path := range groups {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	pointerPattern := regexp.MustCompile(`'roundfix ([^']+) --help'`)
	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			args := append(strings.Fields(path), "no-such-command")
			command := "roundfix " + strings.Join(args, " ")
			var stdout, stderr bytes.Buffer
			if code := runCLI(t, args, &stdout, &stderr); code == exitOK {
				t.Errorf("%s: exit = 0, want non-zero; stdout = %q; stderr = %q", command, stdout.String(), stderr.String())
			}
			words := strings.Fields(stderr.String())
			for index := 1; index < len(words); index++ {
				if words[index] == words[index-1] {
					t.Errorf("%s: stderr repeats word %q in a row: %q", command, words[index], stderr.String())
				}
			}
			for _, pointer := range pointerPattern.FindAllStringSubmatch(stderr.String(), -1) {
				var helpStdout, helpStderr bytes.Buffer
				helpArgs := append(strings.Fields(pointer[1]), "--help")
				if code := runCLI(t, helpArgs, &helpStdout, &helpStderr); code != exitOK {
					t.Errorf("%s: pointer %s exited %d, want 0; stderr = %q", command, pointer[0], code, helpStderr.String())
				}
				for _, child := range groups[path] {
					if !strings.Contains(helpStdout.String(), "roundfix "+child) {
						t.Errorf("%s: pointer %s help is missing command %q", command, pointer[0], "roundfix "+child)
					}
				}
			}
		})
	}
}

// usageErrorGroups returns each command path of helpCommandPaths(rootUsage)
// that prefixes another path, mapped to the command paths it prefixes.
func usageErrorGroups(rootUsage string) map[string][]string {
	paths := helpCommandPaths(rootUsage)
	groups := make(map[string][]string)
	for _, path := range paths {
		for _, child := range paths {
			if strings.HasPrefix(child, path+" ") {
				groups[path] = append(groups[path], child)
			}
		}
	}
	return groups
}
