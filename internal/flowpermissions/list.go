// Package flowpermissions holds the closed Flow Command List shipped with Roundfix.
package flowpermissions

import "slices"

// Entry names a command down to its subcommand, without arguments.
type Entry struct{ Words []string }

// Exclusion names a command deliberately left out of the list and says why.
type Exclusion struct {
	Words  []string
	Reason string
}

// Guard keeps a dangerous option of a listed read behind a Claude Code ask rule.
type Guard struct {
	Entry              Entry
	Option, ClaudeRule string
}

type list struct {
	entries                            []Entry
	exclusions                         []Exclusion
	guards                             []Guard
	notCommands                        []string
	retiredEntries, retiredClaudeRules []string
}

// shipped is data, never derived from the repository or installed tools at runtime.
var shipped = list{
	entries: []Entry{
		{Words: []string{"roundfix", "archive"}},
		{Words: []string{"roundfix", "attach"}},
		{Words: []string{"roundfix", "cancel"}},
		{Words: []string{"roundfix", "corrections"}},
		{Words: []string{"roundfix", "deliver"}},
		{Words: []string{"roundfix", "doctor"}},
		{Words: []string{"roundfix", "events"}},
		{Words: []string{"roundfix", "fetch"}},
		{Words: []string{"roundfix", "follow-up"}},
		{Words: []string{"roundfix", "gc"}},
		{Words: []string{"roundfix", "history"}},
		{Words: []string{"roundfix", "implement"}},
		{Words: []string{"roundfix", "init"}},
		{Words: []string{"roundfix", "migrate"}},
		{Words: []string{"roundfix", "phase"}},
		{Words: []string{"roundfix", "qa-report"}},
		{Words: []string{"roundfix", "reconcile"}},
		{Words: []string{"roundfix", "reopen"}},
		{Words: []string{"roundfix", "resolve"}},
		{Words: []string{"roundfix", "review"}},
		{Words: []string{"roundfix", "runs"}},
		{Words: []string{"roundfix", "settle"}},
		{Words: []string{"roundfix", "setup"}},
		{Words: []string{"roundfix", "spec"}},
		{Words: []string{"roundfix", "stop"}},
		{Words: []string{"roundfix", "storage"}},
		{Words: []string{"roundfix", "supersede"}},
		{Words: []string{"roundfix", "update"}},
		{Words: []string{"roundfix", "upgrade"}},
		{Words: []string{"roundfix", "watch"}},
		{Words: []string{"roundfix", "window"}},
		{Words: []string{"git", "archive"}},
		{Words: []string{"git", "cat-file"}},
		{Words: []string{"git", "check-ignore"}},
		{Words: []string{"git", "diff"}},
		{Words: []string{"git", "diff-tree"}},
		{Words: []string{"git", "grep"}},
		{Words: []string{"git", "ls-files"}},
		{Words: []string{"git", "merge-tree"}},
		{Words: []string{"git", "rev-list"}},
		{Words: []string{"git", "rev-parse"}},
		{Words: []string{"git", "show"}},
		{Words: []string{"git", "status"}},
		{Words: []string{"gh", "api"}},
		{Words: []string{"gh", "pr", "view"}},
	},
	exclusions: []Exclusion{
		{Words: []string{"roundfix", "baseline"}, Reason: "the Roundfix skill tells a session to ask you first"},
		{Words: []string{"roundfix", "profiles"}, Reason: "the Roundfix skill tells a session to ask you first"},
		{Words: []string{"roundfix", "release"}, Reason: "the Roundfix skill tells a session to ask you first"},
		{Words: []string{"roundfix", "skills"}, Reason: "the Roundfix skill tells a session to ask you first"},
		{Words: []string{"roundfix", "permissions"}, Reason: "the command that applies this list"},
		{Words: []string{"git", "-C"}, Reason: "they change the repository or its configuration"},
		{Words: []string{"git", "config"}, Reason: "they change the repository or its configuration"},
		{Words: []string{"git", "init"}, Reason: "they change the repository or its configuration"},
		{Words: []string{"git", "merge"}, Reason: "they change the repository or its configuration"},
		{Words: []string{"git", "mv"}, Reason: "they change the repository or its configuration"},
		{Words: []string{"git", "push"}, Reason: "they change the repository or its configuration"},
		{Words: []string{"git", "remote"}, Reason: "they change the repository or its configuration"},
		{Words: []string{"git", "repack"}, Reason: "they change the repository or its configuration"},
		{Words: []string{"git", "switch"}, Reason: "they change the repository or its configuration"},
		{Words: []string{"git", "worktree"}, Reason: "they change the repository or its configuration"},
		{Words: []string{"gh", "auth"}, Reason: "it changes your credentials"},
	},
	guards: []Guard{
		{Entry: Entry{Words: []string{"git", "archive"}}, Option: "--output", ClaudeRule: "Bash(git archive *--output*)"},
		{Entry: Entry{Words: []string{"git", "archive"}}, Option: "-o", ClaudeRule: "Bash(git archive* -o*)"},
		{Entry: Entry{Words: []string{"git", "archive"}}, Option: "--remote", ClaudeRule: "Bash(git archive *--remote*)"},
		{Entry: Entry{Words: []string{"git", "archive"}}, Option: "--exec", ClaudeRule: "Bash(git archive *--exec*)"},
		{Entry: Entry{Words: []string{"git", "diff"}}, Option: "--output", ClaudeRule: "Bash(git diff *--output*)"},
		{Entry: Entry{Words: []string{"git", "diff"}}, Option: "--ext-diff", ClaudeRule: "Bash(git diff *--ext-diff*)"},
		{Entry: Entry{Words: []string{"git", "diff-tree"}}, Option: "--output", ClaudeRule: "Bash(git diff-tree *--output*)"},
		{Entry: Entry{Words: []string{"git", "diff-tree"}}, Option: "--ext-diff", ClaudeRule: "Bash(git diff-tree *--ext-diff*)"},
		{Entry: Entry{Words: []string{"git", "grep"}}, Option: "--open-files-in-pager", ClaudeRule: "Bash(git grep *--open-files-in-pager*)"},
		{Entry: Entry{Words: []string{"git", "grep"}}, Option: "-O", ClaudeRule: "Bash(git grep* -O*)"},
		{Entry: Entry{Words: []string{"git", "show"}}, Option: "--output", ClaudeRule: "Bash(git show *--output*)"},
		{Entry: Entry{Words: []string{"git", "show"}}, Option: "--ext-diff", ClaudeRule: "Bash(git show *--ext-diff*)"},
		{Entry: Entry{Words: []string{"gh", "api"}}, Option: "--method", ClaudeRule: "Bash(gh api *--method*)"},
		{Entry: Entry{Words: []string{"gh", "api"}}, Option: "-X", ClaudeRule: "Bash(gh api* -X*)"},
		{Entry: Entry{Words: []string{"gh", "api"}}, Option: "--field", ClaudeRule: "Bash(gh api *--field*)"},
		{Entry: Entry{Words: []string{"gh", "api"}}, Option: "-F", ClaudeRule: "Bash(gh api* -F*)"},
		{Entry: Entry{Words: []string{"gh", "api"}}, Option: "--raw-field", ClaudeRule: "Bash(gh api *--raw-field*)"},
		{Entry: Entry{Words: []string{"gh", "api"}}, Option: "-f", ClaudeRule: "Bash(gh api* -f*)"},
		{Entry: Entry{Words: []string{"gh", "api"}}, Option: "--input", ClaudeRule: "Bash(gh api *--input*)"},
	},
	notCommands: []string{"failed", "has"},
	// No entries or rules have been retired in the first release.
}

// Entries returns the roundfix, Git, then GitHub groups, sorted within each group.
func Entries() []Entry { return shipped.copyEntries() }

func (l list) copyEntries() []Entry {
	result := slices.Clone(l.entries)
	for i := range result {
		result[i].Words = slices.Clone(result[i].Words)
	}
	return result
}

// Exclusions returns the Named Exclusions, with their reasons.
func Exclusions() []Exclusion {
	result := slices.Clone(shipped.exclusions)
	for i := range result {
		result[i].Words = slices.Clone(result[i].Words)
	}
	return result
}

// Guards returns the individual Claude Code ask rules for listed reads.
func Guards() []Guard {
	result := slices.Clone(shipped.guards)
	for i := range result {
		result[i].Entry.Words = slices.Clone(result[i].Entry.Words)
	}
	return result
}

// NotCommands returns diagnostic words that follow git or gh but name no command.
func NotCommands() []string { return slices.Clone(shipped.notCommands) }

// Retired returns all formerly shipped entries and Claude rules, in plain text.
func Retired() (entries, claudeRules []string) {
	return slices.Clone(shipped.retiredEntries), slices.Clone(shipped.retiredClaudeRules)
}
