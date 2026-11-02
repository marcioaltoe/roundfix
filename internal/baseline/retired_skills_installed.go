package baseline

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"slices"
)

const (
	RetiredSkillDelete = "delete"
	RetiredSkillKeep   = "keep"

	RetiredSkillEdited     = "edited"
	RetiredSkillUnrecorded = "unrecorded"
	RetiredSkillUnsafe     = "unsafe"
)

// InstalledRetiredSkill identifies a held Retired Skill and its disposition.
type InstalledRetiredSkill struct {
	Skill     string   `json:"skill"`
	Paths     []string `json:"paths"`
	LockEntry bool     `json:"lockEntry"`
	Action    string   `json:"action"`
	Reason    string   `json:"reason,omitempty"`
	LockHash  string   `json:"lockHash,omitempty"`
}

// InstalledRetiredSkills lists each held Retired Skill with its disposition,
// in lexical order of skill. It only reads.
func InstalledRetiredSkills(repoRoot string) ([]InstalledRetiredSkill, error) {
	lockPath := filepath.Join(repoRoot, skillsLockPath)
	lock, err := loadSkillsLock(lockPath)
	if err != nil {
		return nil, fmt.Errorf("inspect %s: %w", lockPath, err)
	}
	locked, _ := lock.root.field("skills")
	var installed []InstalledRetiredSkill
	for _, name := range RetiredSkills() {
		entry := InstalledRetiredSkill{Skill: name, Paths: []string{}, Action: RetiredSkillDelete}
		recorded, exists := locked.field(name)
		entry.LockEntry = exists
		entry.LockHash, _ = recorded.stringField("computedHash")
		unsafe, unrecorded, edited := false, false, false
		for _, root := range []string{".agents/skills", ".claude/skills"} {
			path := root + "/" + name
			absolute := filepath.Join(repoRoot, filepath.FromSlash(path))
			info, err := os.Lstat(absolute)
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			if err != nil {
				return nil, fmt.Errorf("inspect %s: %w", path, err)
			}
			entry.Paths = append(entry.Paths, path)
			if err := validateRestoreParents(repoRoot, path); err != nil {
				unsafe = true
				continue
			}
			if root == ".claude/skills" && info.Mode()&fs.ModeSymlink != 0 {
				target, err := os.Readlink(absolute)
				if err != nil {
					return nil, fmt.Errorf("inspect link %s: %w", path, err)
				}
				if target != "../../.agents/skills/"+name {
					unsafe = true
				}
				continue
			}
			if !info.IsDir() || info.Mode()&fs.ModeSymlink != 0 {
				unsafe = true
				continue
			}
			files, _, err := inspectRestoreTarget(repoRoot, path)
			if err != nil {
				unsafe = true
				continue
			}
			if entry.LockHash == "" {
				unrecorded = true
			} else if externalSkillsLockDigest(files) != entry.LockHash {
				edited = true
			}
		}
		switch {
		case unsafe:
			entry.Action, entry.Reason = RetiredSkillKeep, RetiredSkillUnsafe
		case unrecorded:
			entry.Action, entry.Reason = RetiredSkillKeep, RetiredSkillUnrecorded
		case edited:
			entry.Action, entry.Reason = RetiredSkillKeep, RetiredSkillEdited
		}
		if len(entry.Paths) != 0 || entry.LockEntry {
			installed = append(installed, entry)
		}
	}
	return installed, nil
}

// DeleteRetiredSkills deletes exactly the planned delete entries and returns
// them; it deletes nothing when one changed since the plan.
func DeleteRetiredSkills(ctx context.Context, repoRoot string, planned []InstalledRetiredSkill) ([]InstalledRetiredSkill, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("delete retired skills: %w", err)
	}
	current, err := InstalledRetiredSkills(repoRoot)
	if err != nil {
		return nil, err
	}
	var deleted []InstalledRetiredSkill
	for _, entry := range planned {
		if !slices.Contains(RetiredSkills(), entry.Skill) {
			return nil, fmt.Errorf("refuse to delete %q: it is not a Retired Skill", entry.Skill)
		}
		if entry.Action != RetiredSkillDelete {
			continue
		}
		index := slices.IndexFunc(current, func(candidate InstalledRetiredSkill) bool { return candidate.Skill == entry.Skill })
		if index < 0 || !reflect.DeepEqual(current[index], entry) {
			return nil, fmt.Errorf("retired skill %q changed after the plan; rerun roundfix baseline update", entry.Skill)
		}
		deleted = append(deleted, entry)
	}
	if len(deleted) == 0 {
		return nil, nil
	}
	// Every write goes through an os.Root anchored at the repository, so a
	// symlinked lock file or a parent swapped for a symlink after the plan
	// cannot redirect the write or the recursive delete outside it.
	anchored, err := os.OpenRoot(repoRoot)
	if err != nil {
		return nil, fmt.Errorf("open repository for retired skills: %w", err)
	}
	defer anchored.Close()
	lockPath := filepath.Join(repoRoot, skillsLockPath)
	lock, err := loadSkillsLock(lockPath)
	if err != nil {
		return nil, fmt.Errorf("load retired skills lock: %w", err)
	}
	lockChanged := false
	for _, entry := range deleted {
		if lock.removeSkillsEntry(entry.Skill) {
			lockChanged = true
		}
	}
	var data []byte
	var mode fs.FileMode
	if lockChanged {
		data, err = lock.marshalIndent()
		if err != nil {
			return nil, fmt.Errorf("encode retired skills lock: %w", err)
		}
		info, err := anchored.Lstat(filepath.FromSlash(skillsLockPath))
		if err != nil {
			return nil, fmt.Errorf("inspect retired skills lock mode: %w", err)
		}
		if !info.Mode().IsRegular() {
			return nil, fmt.Errorf("refuse to write %s: it is not a regular file", skillsLockPath)
		}
		mode = info.Mode().Perm()
	}
	for _, entry := range deleted {
		for index := len(entry.Paths) - 1; index >= 0; index-- {
			if err := ctx.Err(); err != nil {
				return nil, fmt.Errorf("delete retired skills: %w", err)
			}
			path := entry.Paths[index]
			if err := anchored.RemoveAll(filepath.FromSlash(path)); err != nil {
				return nil, fmt.Errorf("delete retired skill %s: %w", path, err)
			}
		}
	}
	if lockChanged {
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("write retired skills lock: %w", err)
		}
		if err := anchored.WriteFile(filepath.FromSlash(skillsLockPath), data, mode); err != nil {
			return nil, fmt.Errorf("write retired skills lock: %w", err)
		}
	}
	return deleted, nil
}
