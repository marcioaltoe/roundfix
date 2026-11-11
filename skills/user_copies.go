package skills

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// UserSkillDirectory names a harness's user skill directory and its refresh command.
type UserSkillDirectory struct{ Path, Refresh string }

// StaleUserSkillCopy is a user-level owned copy below the binary's declared version.
type StaleUserSkillCopy struct {
	Skill    string `json:"skill"`
	Path     string `json:"path"`
	Resolved string `json:"resolved,omitempty"`
	Found    string `json:"found"`
	Embedded string `json:"embedded"`
	Refresh  string `json:"refresh"`
}

// UserSkillDirectories mirrors the install targets, then adds the harnesses'
// shared and global directories. A nil getenv reads no environment variables.
func UserSkillDirectories(home string, getenv func(string) string) []UserSkillDirectory {
	if strings.TrimSpace(home) == "" {
		return nil
	}
	directories := make([]UserSkillDirectory, 0, 5)
	for _, target := range []string{"claude", "codex", "opencode"} {
		root := ""
		if getenv != nil {
			root = getenv(strings.ToUpper(target) + "_HOME")
		}
		root = firstNonEmpty(root, filepath.Join(home, "."+target))
		directories = append(directories, UserSkillDirectory{Path: filepath.Join(root, "skills"), Refresh: "roundfix skills install --target " + target})
	}
	for _, extra := range []struct{ target, path string }{
		{"codex", filepath.Join(home, ".agents", "skills")},
		{"opencode", filepath.Join(home, ".config", "opencode", "skills")},
	} {
		directories = append(directories, UserSkillDirectory{Path: extra.path, Refresh: "roundfix skills install --target " + extra.target + " --dir " + quoteUserSkillPath(extra.path)})
	}
	return directories
}

func quoteUserSkillPath(path string) string {
	for _, char := range path {
		if char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char >= '0' && char <= '9' || strings.ContainsRune("@%_+=:,./-", char) {
			continue
		}
		return "'" + strings.ReplaceAll(path, "'", "'\\''") + "'"
	}
	return path
}

// StaleUserSkillCopies follows links and reads each reached owned folder once.
// Missing, unreadable, and non-comparable copies are omitted. It never writes.
func StaleUserSkillCopies(ctx context.Context, directories []UserSkillDirectory) []StaleUserSkillCopy {
	names := Names()
	sort.Strings(names)
	reached := make(map[string]bool)
	var copies []StaleUserSkillCopy
	for _, directory := range directories {
		for _, name := range names {
			if ctx.Err() != nil {
				return copies
			}
			path := filepath.Join(directory.Path, name)
			resolved, err := filepath.EvalSymlinks(path)
			if err != nil || reached[resolved] {
				continue
			}
			info, err := os.Stat(resolved)
			if err != nil || !info.IsDir() {
				continue
			}
			reached[resolved] = true
			if ctx.Err() != nil {
				return copies
			}
			source := filepath.Join(resolved, "SKILL.md")
			data, err := os.ReadFile(source)
			if err != nil {
				continue
			}
			if ctx.Err() != nil {
				return copies
			}
			metadata, ok := parseSkillFrontmatter(string(data))
			if !ok {
				continue
			}
			found := strings.TrimSpace(metadata.Version)
			minimum := ownedSkillMinimumVersions[name]
			state, err := Readiness(SkillVersion{Declared: found, Source: source}, minimum)
			if err != nil || state != ReadinessBelow {
				continue
			}
			copy := StaleUserSkillCopy{Skill: name, Path: path, Found: found, Embedded: minimum, Refresh: directory.Refresh}
			if entry, err := os.Lstat(path); err == nil && entry.Mode()&os.ModeSymlink != 0 {
				copy.Resolved = resolved
			}
			copies = append(copies, copy)
		}
	}
	return copies
}

// Describe names the copy, both versions, and a direct skill link's destination.
func (copy StaleUserSkillCopy) Describe() string {
	description := fmt.Sprintf("%s %s (binary %s) at %s", copy.Skill, copy.Found, copy.Embedded, copy.Path)
	if copy.Resolved != "" {
		description += " -> " + copy.Resolved
	}
	return description
}

// StaleUserSkillRefreshes returns refresh commands once, in first-seen order.
func StaleUserSkillRefreshes(copies []StaleUserSkillCopy) []string {
	seen := make(map[string]bool)
	var commands []string
	for _, copy := range copies {
		if !seen[copy.Refresh] {
			seen[copy.Refresh] = true
			commands = append(commands, copy.Refresh)
		}
	}
	return commands
}
