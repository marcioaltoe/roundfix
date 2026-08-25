package daemon

import (
	"fmt"
	"strings"
)

// hookRefusedClassification is the bounded Run Event classification for a
// Task commit that a repository commit hook refused. The Daemon already
// passed the authoritative Verification (ADR 0014), so the refusal is a
// second gate the loop cannot satisfy rather than a Task failure.
const hookRefusedClassification = "hook_refused"

// hookRefusalOutputMax bounds the hook output carried into the Run Event and
// the Progress line, so a verbose hook cannot flood the journal.
const hookRefusalOutputMax = 1024

// commitHookNames are the Git hooks that can abort `git commit`. They are
// matched longest-first because the shorter names are substrings of the
// longer ones.
var commitHookNames = []string{
	"prepare-commit-msg",
	"pre-merge-commit",
	"pre-commit",
	"commit-msg",
}

// commitHookMarkers name a hook refusal that never says which hook ran:
// the hooks directory itself, or the banner a hook runner prints.
var commitHookMarkers = []string{
	".git/hooks/",
	"husky",
	"lefthook",
	"hook declined",
	"hook refused",
	"hook failed",
	"hook exited",
}

// HookRefusalError reports a commit a repository commit hook refused. The
// work behind it already passed the authoritative Verification, so the
// Daemon keeps it staged and names the recovery instead of discarding it.
type HookRefusalError struct {
	// Hook is the refusing hook when the output named one, else empty.
	Hook     string
	ExitCode int
	// Output is the hook's own output, stderr first.
	Output string
	Err    error
}

func (refusal *HookRefusalError) Error() string {
	detail := ""
	if trimmed := strings.TrimSpace(refusal.Output); trimmed != "" {
		detail = ": " + trimmed
	}
	return fmt.Sprintf("repository %s hook refused the commit (exit status %d)%s", refusal.HookName(), refusal.ExitCode, detail)
}

func (refusal *HookRefusalError) Unwrap() error { return refusal.Err }

// HookName reports the refusing hook, falling back to a generic label when
// the hook's output did not name itself.
func (refusal *HookRefusalError) HookName() string {
	if trimmed := strings.TrimSpace(refusal.Hook); trimmed != "" {
		return trimmed
	}
	return "commit"
}

// ClassifyCommitHookRefusal reports whether failed `git commit` output names
// a repository commit hook refusing the commit, and which hook it named. The
// match is case-insensitive because every hook runner formats its own
// banner. A plain Git failure — nothing staged, an empty message, unmerged
// files — carries no hook marker and stays unclassified.
func ClassifyCommitHookRefusal(output string) (string, bool) {
	lowered := strings.ToLower(strings.TrimSpace(output))
	if lowered == "" {
		return "", false
	}
	for _, name := range commitHookNames {
		if strings.Contains(lowered, name) {
			return name, true
		}
	}
	for _, marker := range commitHookMarkers {
		if strings.Contains(lowered, marker) {
			return "", true
		}
	}
	return "", false
}

// hookOutputExcerpt bounds hook output to its head, where a hook states the
// finding it objected to, and marks the cut so a reader knows it continues.
func hookOutputExcerpt(output string, limit int) string {
	excerpt := strings.TrimSpace(strings.ToValidUTF8(output, ""))
	if excerpt == "" {
		return ""
	}
	if limit < 1 || len(excerpt) <= limit {
		return excerpt
	}
	// The cut can land inside a rune, so the head is revalidated after it.
	return strings.TrimSpace(strings.ToValidUTF8(excerpt[:limit], "")) + " ..."
}
