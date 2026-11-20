package flowpermissions

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type Harness string
type Action string
type State string

const (
	Claude       Harness = "claude"
	Codex        Harness = "codex"
	OpenCode     Harness = "opencode"
	ActionApply  Action  = "apply"
	ActionRemove Action  = "remove"
	Current      State   = "current"
	Behind       State   = "behind"
	NotApplied   State   = "not_applied"
	NotNeeded    State   = "not_needed"
	OwnRules     State   = "own_rules"
	Unreadable   State   = "unreadable"
)

var (
	ErrPreviewChanged    = errors.New("permission preview changed")
	ErrNotWritable       = errors.New("permission plan is not writable")
	ErrUnreadable        = errors.New("permission target is unreadable")
	ErrOutsideRepository = errors.New("permission target is outside the repository or is a link")
)

// Plan describes one repository-relative target and the exact rule changes.
// Digest is empty when there is nothing to write.
type Plan struct {
	Harness    Harness
	Action     Action
	File       string
	FileExists bool
	State      State
	Add        []string
	AddGuards  []string
	Remove     []string
	Reason     string
	Digest     string
}

type planned struct {
	Plan
	after   []byte
	readErr error
}

// Present returns detected harnesses in Claude, Codex, OpenCode order.
func Present(repoRoot string) []Harness {
	var result []Harness
	for _, item := range []struct {
		harness Harness
		paths   []string
	}{
		{Claude, []string{".claude"}}, {Codex, []string{".codex"}}, {OpenCode, []string{"opencode.json", "opencode.jsonc", ".opencode"}},
	} {
		for _, name := range item.paths {
			if _, err := os.Lstat(filepath.Join(repoRoot, name)); err == nil {
				result = append(result, item.harness)
				break
			}
		}
	}
	return result
}

// Preview reads only repository files and never creates directories or files.
func Preview(repoRoot string, harness Harness, action Action) (Plan, error) {
	result, err := plan(repoRoot, harness, action, shipped)
	return result.Plan, err
}

// plan is the single planner, parameterized by a list so retired entries can be
// exercised without changing the binary's shipped data.
func plan(repoRoot string, harness Harness, action Action, rules list) (planned, error) {
	result := planned{Plan: Plan{Harness: harness, Action: action}}
	if action != ActionApply && action != ActionRemove {
		return result, fmt.Errorf("unknown permission action %q", action)
	}
	switch harness {
	case Claude:
		result.File = ".claude/settings.json"
	case Codex:
		result.File = ".codex/rules/roundfix-flow.rules"
	case OpenCode:
		result.State, result.Reason = ReadOpenCode(repoRoot)
		return result, nil
	default:
		return result, fmt.Errorf("unknown permission harness %q", harness)
	}
	if err := checkTarget(repoRoot, result.File); err != nil {
		if errors.Is(err, ErrOutsideRepository) {
			return result, err
		}
		result.State, result.Reason, result.readErr = Unreadable, err.Error(), err
		return result, nil
	}
	data, err := os.ReadFile(filepath.Join(repoRoot, filepath.FromSlash(result.File)))
	result.FileExists = !errors.Is(err, os.ErrNotExist)
	if err != nil && result.FileExists {
		result.State, result.Reason, result.readErr = Unreadable, err.Error(), err
		return result, nil
	}
	if !result.FileExists {
		data = nil
	}
	if harness == Claude {
		allow, ask, err := ReadClaudeSettings(data)
		if err != nil {
			result.State, result.Reason, result.readErr = Unreadable, err.Error(), err
			return result, nil
		}
		current := make([]string, 0, len(rules.entries))
		var retired []string
		hasRoundfix := false
		for _, entry := range rules.entries {
			rule := ClaudeRule(entry)
			current = append(current, rule)
			if len(entry.Words) > 0 && entry.Words[0] == "roundfix" && slices.Contains(allow, rule) {
				hasRoundfix = true
			}
		}
		for _, entry := range rules.retiredEntries {
			rule := ClaudeRule(Entry{Words: strings.Fields(entry)})
			retired = append(retired, rule)
			if strings.HasPrefix(entry, "roundfix ") && slices.Contains(allow, rule) {
				hasRoundfix = true
			}
		}
		var guards []string
		for _, guard := range rules.guards {
			guards = append(guards, guard.ClaudeRule)
		}
		missing, missingGuards := difference(current, allow), difference(guards, ask)
		oldAllow, oldAsk := intersection(allow, retired), intersection(ask, rules.retiredClaudeRules)
		result.State = NotApplied
		if hasRoundfix {
			result.State = Current
			if len(missing)+len(missingGuards)+len(oldAllow)+len(oldAsk) > 0 {
				result.State = Behind
			}
		}
		var removeAllow, removeAsk []string
		if action == ActionApply {
			result.Add, result.AddGuards = missing, missingGuards
			removeAllow, removeAsk = oldAllow, oldAsk
		} else {
			removeAllow = intersection(allow, append(slices.Clone(current), retired...))
			removeAsk = intersection(ask, append(slices.Clone(guards), rules.retiredClaudeRules...))
		}
		result.Remove = append(slices.Clone(removeAllow), removeAsk...)
		if len(result.Add)+len(result.AddGuards)+len(result.Remove) == 0 {
			return result, nil
		}
		result.after, err = editClaudeSettings(data, result.Add, result.AddGuards, removeAllow, removeAsk)
		if err != nil {
			return result, fmt.Errorf("plan Claude settings: %w", err)
		}
	} else {
		var entries []string
		result.State = NotApplied
		if result.FileExists {
			var ok bool
			entries, ok = ReadCodexRules(data)
			if !ok {
				result.State, result.Reason = Unreadable, "Codex rules do not have Roundfix's rendered form"
				result.readErr = errors.New(result.Reason)
				return result, nil
			}
			result.State = Behind
			if bytes.Equal(data, rules.renderCodexRules()) {
				result.State = Current
			}
		}
		if action == ActionRemove {
			if !result.FileExists {
				return result, nil
			}
			result.Remove = entries
		} else {
			result.after = rules.renderCodexRules()
			if result.FileExists && bytes.Equal(data, result.after) {
				return result, nil
			}
			result.Add = difference(rules.plainEntries(), entries)
			result.Remove = difference(entries, rules.plainEntries())
		}
	}
	result.Digest = planDigest(result.Plan, data)
	return result, nil
}

func difference(values, other []string) []string {
	var result []string
	for _, value := range values {
		if !slices.Contains(other, value) && !slices.Contains(result, value) {
			result = append(result, value)
		}
	}
	return result
}

func intersection(values, other []string) []string {
	var result []string
	for _, value := range values {
		if slices.Contains(other, value) && !slices.Contains(result, value) {
			result = append(result, value)
		}
	}
	return result
}

func planDigest(p Plan, data []byte) string {
	fileHash := "absent"
	if p.FileExists {
		fileHash = fmt.Sprintf("%x", sha256.Sum256(data))
	}
	lines := []string{string(p.Harness), string(p.Action), p.File, fileHash}
	for _, group := range []struct {
		prefix string
		rules  []string
	}{{"+", p.Add}, {"?", p.AddGuards}, {"-", p.Remove}} {
		sorted := slices.Clone(group.rules)
		slices.Sort(sorted)
		for _, rule := range sorted {
			lines = append(lines, group.prefix+rule)
		}
	}
	sum := sha256.Sum256([]byte(strings.Join(lines, "\n") + "\n"))
	return fmt.Sprintf("%x", sum)[:12]
}

// Apply recomputes the preview before touching a directory or target. An
// unreadable target retains its underlying error for errors.Is inspection.
func Apply(repoRoot string, harness Harness, action Action, digest string) (Plan, error) {
	if harness == OpenCode {
		return Plan{Harness: harness, Action: action}, ErrNotWritable
	}
	result, err := plan(repoRoot, harness, action, shipped)
	if err != nil {
		return result.Plan, err
	}
	if result.State == Unreadable {
		return result.Plan, fmt.Errorf("%w: %s: %w", ErrUnreadable, result.File, result.readErr)
	}
	if digest == "" {
		return result.Plan, ErrNotWritable
	}
	if result.Digest != digest {
		return result.Plan, ErrPreviewChanged
	}
	if err := checkTarget(repoRoot, result.File); err != nil {
		return result.Plan, err
	}
	target := filepath.Join(repoRoot, filepath.FromSlash(result.File))
	if harness == Codex && action == ActionRemove {
		if err := os.Remove(target); err != nil {
			return result.Plan, fmt.Errorf("remove Codex rules: %w", err)
		}
		return result.Plan, nil
	}
	directory := filepath.Dir(target)
	if err := os.MkdirAll(directory, 0755); err != nil {
		return result.Plan, fmt.Errorf("create permission directory: %w", err)
	}
	temp, err := os.CreateTemp(directory, ".roundfix-flow-*")
	if err != nil {
		return result.Plan, fmt.Errorf("create permission temporary file: %w", err)
	}
	defer os.Remove(temp.Name())
	mode := os.FileMode(0644)
	if result.FileExists {
		info, err := os.Stat(target)
		if err != nil {
			temp.Close()
			return result.Plan, fmt.Errorf("stat permission target: %w", err)
		}
		mode = info.Mode().Perm()
	}
	if err := temp.Chmod(mode); err != nil {
		temp.Close()
		return result.Plan, fmt.Errorf("set permission file mode: %w", err)
	}
	if _, err := temp.Write(result.after); err != nil {
		temp.Close()
		return result.Plan, fmt.Errorf("write permission file: %w", err)
	}
	if err := temp.Close(); err != nil {
		return result.Plan, fmt.Errorf("close permission file: %w", err)
	}
	if err := checkTarget(repoRoot, result.File); err != nil {
		return result.Plan, err
	}
	if err := os.Rename(temp.Name(), target); err != nil {
		return result.Plan, fmt.Errorf("replace permission file: %w", err)
	}
	return result.Plan, nil
}

// Resolve the repository root itself (worktrees can have logical roots), then
// inspect every existing component before any file is opened. Target links are
// refused outright; directory links are allowed only within this repository.
func checkTarget(repoRoot, relative string) error {
	root, err := filepath.EvalSymlinks(repoRoot)
	if err != nil {
		return fmt.Errorf("resolve repository: %w", err)
	}
	root, err = filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("resolve absolute repository: %w", err)
	}
	path := root
	parts := strings.Split(filepath.FromSlash(relative), string(filepath.Separator))
	for i, part := range parts {
		path = filepath.Join(path, part)
		info, err := os.Lstat(path)
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("inspect permission target: %w", err)
		}
		if info.Mode()&os.ModeSymlink == 0 {
			continue
		}
		if i == len(parts)-1 {
			return fmt.Errorf("%w: %s", ErrOutsideRepository, relative)
		}
		resolved, err := filepath.EvalSymlinks(path)
		if err != nil {
			return fmt.Errorf("resolve permission directory: %w", err)
		}
		rel, err := filepath.Rel(root, resolved)
		if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.IsAbs(rel) {
			return fmt.Errorf("%w: %s", ErrOutsideRepository, relative)
		}
		path = resolved
	}
	return nil
}
