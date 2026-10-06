package baseline

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// InstalledRetiredSkill identifies the retired entries an adopter can remove.
type InstalledRetiredSkill struct {
	Skill     string   `json:"skill"`
	Paths     []string `json:"paths"`
	LockEntry bool     `json:"lockEntry"`
}

// InstalledRetiredSkills lists, in lexical order of skill, each Retired Skill
// the repository still holds. It only reads.
func InstalledRetiredSkills(repoRoot string) ([]InstalledRetiredSkill, error) {
	lockPath := filepath.Join(repoRoot, skillsLockPath)
	lock, err := loadSkillsLock(lockPath)
	if err != nil {
		return nil, fmt.Errorf("inspect %s: %w", lockPath, err)
	}
	locked, _ := lock.root.field("skills")
	var installed []InstalledRetiredSkill
	for _, name := range RetiredSkills() {
		entry := InstalledRetiredSkill{Skill: name, Paths: []string{}}
		_, entry.LockEntry = locked.field(name)
		for _, root := range []string{".agents/skills", ".claude/skills"} {
			path := root + "/" + name
			_, err := os.Lstat(filepath.Join(repoRoot, filepath.FromSlash(path)))
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			if err != nil {
				return nil, fmt.Errorf("inspect %s: %w", path, err)
			}
			entry.Paths = append(entry.Paths, path)
		}
		if len(entry.Paths) != 0 || entry.LockEntry {
			installed = append(installed, entry)
		}
	}
	return installed, nil
}
