package daemon

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
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

// installedCommitHookNames are the hooks whose non-zero exit aborts a plain
// `git commit`. Finding one installed is the structural signal, so the list is
// deliberately narrower than commitHookNames: pre-merge-commit runs only for a
// merge, and neither it nor a prepare-commit-msg hook can refuse the commit
// this Daemon makes.
var installedCommitHookNames = []string{
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

// ClassifyCommitHookRefusal reports whether a failed `git commit` in workDir
// was refused by a repository commit hook, and which hook it named. It reads
// two signals, because either one alone leaves a real refusal unclassified.
//
// The first is the output itself: a hook runner's banner or a hook that names
// itself. That match is case-insensitive because every runner formats its own
// banner, and it is tried first so a hook that identifies itself keeps the
// diagnostics it already had — the named hook and its banner.
//
// The second is structural: the repository holds an executable commit hook.
// Git prints nothing identifying when a hook fails, so a hook that emits only
// its finding — the shape this Spec's three cases were measured in — names
// nothing the first signal can see, and only the repository can answer whether
// a hook ran at all.
//
// The structural signal accepts one bounded false positive: a genuine Git
// failure — an index lock, a full disk — in a repository that does have a
// commit hook is reported as a refusal. That misnaming is harmless, because
// the recovery it names re-runs the authoritative Verification before it
// commits, so a real Git failure surfaces there instead of being committed
// around. A repository with no executable commit hook keeps returning the raw
// Git error, so the same failure is never relabelled where no hook exists.
func ClassifyCommitHookRefusal(ctx context.Context, workDir string, output string) (string, bool) {
	if hook, named := classifyCommitHookOutput(output); named {
		return hook, true
	}
	return "", commitHookInstalled(ctx, workDir)
}

// classifyCommitHookOutput reports which commit hook the failed commit's own
// output names, and whether it names a refusal at all.
func classifyCommitHookOutput(output string) (string, bool) {
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

// commitHookInstalled reports whether workDir's repository holds an executable
// hook that can refuse a commit. Git resolves the hooks directory itself, so
// the answer honours core.hooksPath — the husky layout the measured repository
// used — and stays correct from a subdirectory or a linked worktree. When Git
// answers nothing, because workDir is outside any repository or the running Git
// predates --path-format, the signal reports no hook: the classification falls
// back to the output alone, which is where it stood before this second signal.
func commitHookInstalled(ctx context.Context, workDir string) bool {
	if strings.TrimSpace(workDir) == "" {
		return false
	}
	resolved, err := runGitOutput(ctx, workDir, "rev-parse", "--path-format=absolute", "--git-path", "hooks")
	if err != nil {
		return false
	}
	hooksDir := strings.TrimSpace(resolved)
	if hooksDir == "" {
		return false
	}
	for _, name := range installedCommitHookNames {
		if executableFile(filepath.Join(hooksDir, name)) {
			return true
		}
	}
	return false
}

// executableFile reports whether path is a regular file Git would execute as a
// hook. A hook without the executable bit is skipped by Git, so it cannot be
// the reason a commit failed.
func executableFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular() && info.Mode().Perm()&0o111 != 0
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
