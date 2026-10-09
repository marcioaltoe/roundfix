// Suite: installed Retired Skill inspection.
// Invariant: inspection reports repository entries without changing them.
// Boundary IN: skill trees, links, and the local skills lock.
// Boundary OUT: acquisition, deletion, and Baseline planning.

package baseline

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestInstalledRetiredSkillsListsTreesLinksAndLockEntries(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		paths []string
		links []string
		lock  string
		want  []InstalledRetiredSkill
	}{
		{
			name:  "trees links and lock",
			paths: []string{".agents/skills/council", ".agents/skills/the-fool", ".agents/skills/autoresearch"},
			links: []string{".claude/skills/council"},
			lock:  `{"version":1,"skills":{"the-fool":{},"council":{},"autoresearch":{}}}`,
			want: []InstalledRetiredSkill{
				{Skill: "council", Paths: []string{".agents/skills/council", ".claude/skills/council"}, LockEntry: true},
				{Skill: "the-fool", Paths: []string{".agents/skills/the-fool"}, LockEntry: true},
			},
		},
		{
			name: "lock only",
			lock: `{"version":1,"skills":{"the-fool":{}}}`,
			want: []InstalledRetiredSkill{{Skill: "the-fool", Paths: []string{}, LockEntry: true}},
		},
		{
			name:  "missing lock",
			paths: []string{".agents/skills/council"},
			want:  []InstalledRetiredSkill{{Skill: "council", Paths: []string{".agents/skills/council"}}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := t.TempDir()
			for _, path := range test.paths {
				if err := os.MkdirAll(filepath.Join(repo, filepath.FromSlash(path)), 0o755); err != nil {
					t.Fatal(err)
				}
			}
			for _, path := range test.links {
				absolute := filepath.Join(repo, filepath.FromSlash(path))
				if err := os.MkdirAll(filepath.Dir(absolute), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink("missing-target", absolute); err != nil {
					t.Fatal(err)
				}
			}
			lockPath := filepath.Join(repo, "skills-lock.json")
			if test.lock != "" {
				if err := os.WriteFile(lockPath, []byte(test.lock), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			got, err := InstalledRetiredSkills(repo)
			if err != nil || !reflect.DeepEqual(got, test.want) {
				t.Fatalf("inspection = %+v, %v; want %+v", got, err, test.want)
			}
			for _, path := range append(test.paths, test.links...) {
				if _, err := os.Lstat(filepath.Join(repo, filepath.FromSlash(path))); err != nil {
					t.Fatalf("inspection changed %s: %v", path, err)
				}
			}
			if test.lock != "" {
				data, err := os.ReadFile(lockPath)
				if err != nil || string(data) != test.lock {
					t.Fatalf("inspection changed lock: %s, %v", data, err)
				}
			}
		})
	}
}

func TestInstalledRetiredSkillsIsEmptyWithoutRetiredCopies(t *testing.T) {
	t.Parallel()
	for _, lock := range []string{"", `{"version":1,"skills":{"autoresearch":{}}}`} {
		repo := t.TempDir()
		if lock != "" {
			if err := os.WriteFile(filepath.Join(repo, "skills-lock.json"), []byte(lock), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		got, err := InstalledRetiredSkills(repo)
		if err != nil || len(got) != 0 {
			t.Fatalf("inspection = %+v, %v; want empty", got, err)
		}
	}
}

func TestInstalledRetiredSkillsReportsFilesystemErrors(t *testing.T) {
	t.Parallel()
	t.Run("lock is not a regular file", func(t *testing.T) {
		repo := t.TempDir()
		if err := os.Mkdir(filepath.Join(repo, "skills-lock.json"), 0o755); err != nil {
			t.Fatal(err)
		}
		if _, err := InstalledRetiredSkills(repo); err == nil || !strings.Contains(err.Error(), "skills-lock.json") {
			t.Fatalf("lock error = %v", err)
		}
	})
	t.Run("skill parent is not a directory", func(t *testing.T) {
		repo := t.TempDir()
		if err := os.WriteFile(filepath.Join(repo, ".agents"), []byte("not a directory"), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := InstalledRetiredSkills(repo); err == nil || !strings.Contains(err.Error(), ".agents/skills/council") {
			t.Fatalf("skill path error = %v", err)
		}
	})
}

func TestInstalledRetiredSkillsRefusesAMalformedLock(t *testing.T) {
	t.Parallel()
	for _, lock := range []string{`{`, `{"version":1,"skills":[]}`} {
		repo := t.TempDir()
		if err := os.WriteFile(filepath.Join(repo, "skills-lock.json"), []byte(lock), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := InstalledRetiredSkills(repo); err == nil || !strings.Contains(err.Error(), "skills-lock.json") {
			t.Fatalf("malformed lock error = %v", err)
		}
	}
}
