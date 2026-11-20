package flowpermissions

import (
	"os"
	"slices"
	"strings"
	"testing"
)

func TestFlowCommandList(t *testing.T) {
	t.Parallel()
	entries := Entries()
	t.Run("roundfix_group_first_and_sorted", func(t *testing.T) {
		var group []string
		left := false
		for _, entry := range entries {
			if entry.Words[0] != "roundfix" {
				left = true
				continue
			}
			if left {
				t.Errorf("roundfix entry after another group: %v", entry.Words)
			}
			group = append(group, strings.Join(entry.Words, " "))
		}
		if len(group) == 0 || !slices.IsSorted(group) {
			t.Errorf("roundfix group is empty or unsorted: %v", group)
		}
	})
	t.Run("git_and_gh_groups_sorted", func(t *testing.T) {
		for _, group := range []string{"git", "gh"} {
			var commands []string
			for _, entry := range entries {
				if entry.Words[0] == group {
					commands = append(commands, strings.Join(entry.Words, " "))
				}
			}
			if len(commands) == 0 || !slices.IsSorted(commands) {
				t.Errorf("%s group is empty or unsorted: %v", group, commands)
			}
		}
		previous := -1
		for _, entry := range entries {
			rank := slices.Index([]string{"roundfix", "git", "gh"}, entry.Words[0])
			if rank < previous || rank < 0 {
				t.Errorf("wrong group order: %v", entry.Words)
			}
			previous = rank
		}
	})
	t.Run("no_entry_twice", func(t *testing.T) {
		seen := map[string]bool{}
		for _, entry := range PlainEntries() {
			if seen[entry] {
				t.Errorf("duplicate entry %s", entry)
			}
			seen[entry] = true
		}
	})
	t.Run("no_word_on_both_sides", func(t *testing.T) {
		for _, exclusion := range Exclusions() {
			name := strings.Join(exclusion.Words, " ")
			for _, entry := range entries {
				if wordPrefix(entry.Words, exclusion.Words) || wordPrefix(exclusion.Words, entry.Words) {
					t.Errorf("entry overlaps exclusion %s", name)
				}
			}
			if len(exclusion.Words) > 1 && slices.Contains(NotCommands(), exclusion.Words[1]) {
				t.Errorf("exclusion %s is also a noncommand", name)
			}
		}
		for _, entry := range entries {
			if slices.Contains(NotCommands(), entry.Words[1]) {
				t.Errorf("entry %v is also a noncommand", entry.Words)
			}
		}
		for _, word := range NotCommands() {
			if strings.HasPrefix(word, "-") {
				t.Errorf("option %s is not a diagnostic word", word)
			}
		}
	})
	t.Run("every_exclusion_has_a_reason", func(t *testing.T) {
		for _, exclusion := range Exclusions() {
			if len(exclusion.Words) < 2 || strings.TrimSpace(exclusion.Reason) == "" {
				t.Errorf("incomplete exclusion: %+v", exclusion)
			}
		}
	})
	t.Run("permissions_is_excluded", func(t *testing.T) {
		for _, exclusion := range Exclusions() {
			if slices.Equal(exclusion.Words, []string{"roundfix", "permissions"}) && exclusion.Reason == "the command that applies this list" {
				return
			}
		}
		t.Fatal("roundfix permissions needs its Named Exclusion")
	})
	t.Run("entries_name_no_argument", func(t *testing.T) {
		for _, entry := range entries {
			words := entry.Words
			want := 2
			if len(words) > 1 && words[0] == "gh" && words[1] == "pr" {
				want = 3
			}
			if len(words) != want {
				t.Errorf("entry names an argument or lacks a subcommand: %v", words)
			}
			for _, word := range words {
				if !codexWord.MatchString(word) {
					t.Errorf("entry names no plain command word: %v", words)
				}
			}
		}
	})
}

func wordPrefix(words, prefix []string) bool {
	return len(words) >= len(prefix) && slices.Equal(words[:len(prefix)], prefix)
}

func TestListAccessorsReturnIndependentData(t *testing.T) {
	t.Parallel()
	entries := Entries()
	entries[0].Words[0] = "changed"
	exclusions := Exclusions()
	exclusions[0].Words[0] = "changed"
	exclusions[0].Reason = "changed"
	guards := Guards()
	guards[0].Entry.Words[0] = "changed"
	guards[0].ClaudeRule = "changed"
	words := NotCommands()
	words[0] = "changed"
	plain := PlainEntries()
	plain[0] = "changed"
	if Entries()[0].Words[0] == "changed" || Exclusions()[0].Words[0] == "changed" || Exclusions()[0].Reason == "changed" || Guards()[0].Entry.Words[0] == "changed" || Guards()[0].ClaudeRule == "changed" || NotCommands()[0] == "changed" || PlainEntries()[0] == "changed" {
		t.Fatal("an accessor exposed the shipped list's backing storage")
	}
}

func currentReleaseLines() []string {
	lines := PlainEntries()
	for _, entry := range Entries() {
		lines = append(lines, ClaudeRule(entry))
	}
	for _, guard := range Guards() {
		lines = append(lines, guard.ClaudeRule)
	}
	slices.Sort(lines)
	return lines
}

func TestRetiredSetKeepsEveryReleasedEntry(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile("testdata/released.txt")
	if err != nil {
		t.Fatal(err)
	}
	recorded := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
	current := currentReleaseLines()
	entries, rules := Retired()
	retired := append(entries, rules...)
	report := func(t *testing.T, line string) {
		t.Helper()
		t.Errorf("%s: a release that drops an entry adds it to the retired set and a release that adds one appends it to the record", line)
	}
	t.Run("every_recorded_line_is_current_or_retired", func(t *testing.T) {
		if !slices.IsSorted(recorded) {
			t.Error("released.txt must be sorted")
		}
		seen := map[string]bool{}
		for _, line := range recorded {
			if seen[line] || (!slices.Contains(current, line) && !slices.Contains(retired, line)) {
				report(t, line)
			}
			seen[line] = true
		}
	})
	t.Run("every_current_entry_and_rule_is_recorded", func(t *testing.T) {
		for _, line := range current {
			if !slices.Contains(recorded, line) {
				report(t, line)
			}
		}
	})
	t.Run("retired_is_not_current", func(t *testing.T) {
		for _, line := range retired {
			if slices.Contains(current, line) {
				report(t, line)
			}
		}
	})
}
