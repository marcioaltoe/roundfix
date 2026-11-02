// Suite: installed Retired Skill inspection and approved deletion.
// Invariant: only planned, proven unedited Retired Skills can leave the repository.
// Boundary IN: skill trees, links, and the local skills lock.
// Boundary OUT: acquisition and Baseline planning.

package baseline

import (
	"context"
	"encoding/json"
	"errors"
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
				{Skill: "council", Paths: []string{".agents/skills/council", ".claude/skills/council"}, LockEntry: true, Action: RetiredSkillKeep, Reason: RetiredSkillUnsafe},
				{Skill: "the-fool", Paths: []string{".agents/skills/the-fool"}, LockEntry: true, Action: RetiredSkillKeep, Reason: RetiredSkillUnrecorded},
			},
		},
		{
			name: "lock only",
			lock: `{"version":1,"skills":{"the-fool":{}}}`,
			want: []InstalledRetiredSkill{{Skill: "the-fool", Paths: []string{}, LockEntry: true, Action: RetiredSkillDelete}},
		},
		{
			name:  "missing lock",
			paths: []string{".agents/skills/council"},
			want:  []InstalledRetiredSkill{{Skill: "council", Paths: []string{".agents/skills/council"}, Action: RetiredSkillKeep, Reason: RetiredSkillUnrecorded}},
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

func TestRetiredSkillDispositionsProveACopyUnedited(t *testing.T) {
	t.Parallel()
	repo := t.TempDir()
	hash := externalSkillsLockDigest([]restoreFile{{Path: "SKILL.md", Content: []byte("# installed\n")}})
	for _, relative := range []string{".agents/skills/tdd", ".agents/skills/design-patterns", ".agents/skills/council", ".claude/skills/fixing-metadata"} {
		writeRetiredSkillTestFile(t, repo, relative+"/SKILL.md", "# installed\n")
	}
	writeRetiredSkillTestFile(t, repo, ".agents/skills/design-patterns/SKILL.md", "# edited\n")
	writeRetiredSkillTestLink(t, repo, ".claude/skills/tdd", "../../.agents/skills/tdd")
	writeRetiredSkillTestLink(t, repo, ".claude/skills/to-prompt", "../../elsewhere/to-prompt")
	writeRetiredSkillTestLock(t, repo, map[string]any{
		"tdd":             map[string]any{"computedHash": hash},
		"design-patterns": map[string]any{"computedHash": hash},
		"fixing-metadata": map[string]any{"computedHash": hash},
		"the-fool":        map[string]any{},
	})
	want := []InstalledRetiredSkill{
		{Skill: "council", Paths: []string{".agents/skills/council"}, Action: RetiredSkillKeep, Reason: RetiredSkillUnrecorded},
		{Skill: "design-patterns", Paths: []string{".agents/skills/design-patterns"}, LockEntry: true, Action: RetiredSkillKeep, Reason: RetiredSkillEdited, LockHash: hash},
		{Skill: "fixing-metadata", Paths: []string{".claude/skills/fixing-metadata"}, LockEntry: true, Action: RetiredSkillDelete, LockHash: hash},
		{Skill: "tdd", Paths: []string{".agents/skills/tdd", ".claude/skills/tdd"}, LockEntry: true, Action: RetiredSkillDelete, LockHash: hash},
		{Skill: "the-fool", Paths: []string{}, LockEntry: true, Action: RetiredSkillDelete},
		{Skill: "to-prompt", Paths: []string{".claude/skills/to-prompt"}, Action: RetiredSkillKeep, Reason: RetiredSkillUnsafe},
	}
	got, err := InstalledRetiredSkills(repo)
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("dispositions = %+v, %v; want %+v", got, err, want)
	}
	for _, shape := range []string{"safe link only", "agents link", "plain file", "link inside tree", "linked parent", "edited claude copy"} {
		t.Run(shape, func(t *testing.T) {
			root := t.TempDir()
			writeRetiredSkillTestLock(t, root, map[string]any{"tdd": map[string]any{"computedHash": hash}})
			switch shape {
			case "safe link only":
				writeRetiredSkillTestLink(t, root, ".claude/skills/tdd", "../../.agents/skills/tdd")
			case "agents link":
				writeRetiredSkillTestLink(t, root, ".agents/skills/tdd", "../../elsewhere/tdd")
			case "plain file":
				writeRetiredSkillTestFile(t, root, ".agents/skills/tdd", "not a folder")
			case "link inside tree":
				writeRetiredSkillTestLink(t, root, ".agents/skills/tdd/SKILL.md", "missing")
			case "linked parent":
				outside := t.TempDir()
				writeRetiredSkillTestLink(t, outside, "skills/tdd", "../../.agents/skills/tdd")
				writeRetiredSkillTestLink(t, root, ".claude", outside)
			case "edited claude copy":
				writeRetiredSkillTestFile(t, root, ".agents/skills/tdd/SKILL.md", "# installed\n")
				writeRetiredSkillTestFile(t, root, ".claude/skills/tdd/SKILL.md", "# edited\n")
			}
			got, err := InstalledRetiredSkills(root)
			action, reason := RetiredSkillKeep, RetiredSkillUnsafe
			if shape == "safe link only" {
				action, reason = RetiredSkillDelete, ""
			}
			if shape == "edited claude copy" {
				reason = RetiredSkillEdited
			}
			if err != nil || len(got) != 1 || got[0].Action != action || got[0].Reason != reason {
				t.Fatalf("disposition = %+v, %v; want %s/%s", got, err, action, reason)
			}
		})
	}
}

func TestDeleteRetiredSkillsDeletesOnlyThePlannedCopies(t *testing.T) {
	t.Parallel()
	fixture := func(t *testing.T) (string, []InstalledRetiredSkill) {
		t.Helper()
		repo := t.TempDir()
		hash := externalSkillsLockDigest([]restoreFile{{Path: "SKILL.md", Content: []byte("# installed\n")}})
		entries := map[string]any{}
		for _, name := range []string{"tdd", "design-patterns", "context7-cli", "exa-web-search", "local-only"} {
			writeRetiredSkillTestFile(t, repo, ".agents/skills/"+name+"/SKILL.md", "# installed\n")
			writeRetiredSkillTestLink(t, repo, ".claude/skills/"+name, "../../.agents/skills/"+name)
			entries[name] = map[string]any{"computedHash": hash, "extra": "preserved"}
		}
		writeRetiredSkillTestFile(t, repo, ".agents/skills/design-patterns/SKILL.md", "# edited\n")
		writeRetiredSkillTestLock(t, repo, entries)
		planned, err := InstalledRetiredSkills(repo)
		if err != nil {
			t.Fatal(err)
		}
		return repo, planned
	}
	t.Run("deletes the unedited retired copy and keeps every other skill", func(t *testing.T) {
		repo, planned := fixture(t)
		lockPath := filepath.Join(repo, "skills-lock.json")
		before := readSkillsRestoreLock(t, repo)
		if err := os.Chmod(lockPath, 0o600); err != nil {
			t.Fatal(err)
		}
		mode, err := os.Stat(lockPath)
		if err != nil {
			t.Fatal(err)
		}
		deleted, err := DeleteRetiredSkills(t.Context(), repo, planned)
		if err != nil || !reflect.DeepEqual(deleted, planned[1:]) {
			t.Fatalf("deleted = %+v, %v; want %+v", deleted, err, planned[1:])
		}
		for _, root := range []string{".agents", ".claude"} {
			if _, err := os.Lstat(filepath.Join(repo, root, "skills", "tdd")); !os.IsNotExist(err) {
				t.Fatalf("tdd %s remains: %v", root, err)
			}
		}
		for _, name := range []string{"design-patterns", "context7-cli", "exa-web-search", "local-only"} {
			want := "# installed\n"
			if name == "design-patterns" {
				want = "# edited\n"
			}
			assertRestoreFile(t, repo, ".agents/skills/"+name+"/SKILL.md", want)
			link, err := os.Readlink(filepath.Join(repo, ".claude", "skills", name))
			if err != nil || link != "../../.agents/skills/"+name {
				t.Fatalf("survivor link %s = %q, %v", name, link, err)
			}
		}
		delete(before["skills"].(map[string]any), "tdd")
		if after := readSkillsRestoreLock(t, repo); !reflect.DeepEqual(after, before) {
			t.Fatalf("lock changed beyond tdd: %+v", after)
		}
		afterMode, err := os.Stat(lockPath)
		if err != nil || afterMode.Mode() != mode.Mode() {
			t.Fatalf("lock mode changed: %v, %v", afterMode, err)
		}
	})
	for _, shape := range []string{"refuses a copy edited after the plan", "refuses a name outside the Retired Skills", "refuses a later planned entry before deleting anything", "honors cancellation", "deletes only the selected copy"} {
		t.Run(shape, func(t *testing.T) {
			repo, planned := fixture(t)
			ctx := t.Context()
			wantError := ""
			switch shape {
			case "refuses a copy edited after the plan":
				writeRetiredSkillTestFile(t, repo, ".agents/skills/tdd/SKILL.md", "# changed after plan\n")
				wantError = `retired skill "tdd" changed after the plan; rerun roundfix baseline update`
			case "refuses a name outside the Retired Skills":
				planned = append(planned, InstalledRetiredSkill{Skill: "context7-cli", Action: RetiredSkillDelete})
				wantError = `refuse to delete "context7-cli": it is not a Retired Skill`
			case "refuses a later planned entry before deleting anything":
				planned = append(planned, InstalledRetiredSkill{Skill: "the-fool", Action: RetiredSkillDelete})
				wantError = `retired skill "the-fool" changed after the plan; rerun roundfix baseline update`
			case "honors cancellation":
				canceled, cancel := context.WithCancel(ctx)
				cancel()
				ctx = canceled
			case "deletes only the selected copy":
				writeRetiredSkillTestFile(t, repo, ".agents/skills/the-fool/SKILL.md", "# installed\n")
			}
			lockBefore, err := os.ReadFile(filepath.Join(repo, "skills-lock.json"))
			if err != nil {
				t.Fatal(err)
			}
			deleted, err := DeleteRetiredSkills(ctx, repo, planned)
			if shape == "deletes only the selected copy" {
				if err != nil || len(deleted) != 1 || deleted[0].Skill != "tdd" {
					t.Fatalf("selected deletion = %+v, %v", deleted, err)
				}
				assertRestoreFile(t, repo, ".agents/skills/the-fool/SKILL.md", "# installed\n")
				return
			}
			if err == nil || len(deleted) != 0 || (wantError != "" && err.Error() != wantError) {
				t.Fatalf("refusal = %+v, %v; want %q", deleted, err, wantError)
			}
			if shape == "honors cancellation" && !errors.Is(err, context.Canceled) {
				t.Fatalf("cancellation = %v", err)
			}
			lockAfter, readErr := os.ReadFile(filepath.Join(repo, "skills-lock.json"))
			if readErr != nil || string(lockBefore) != string(lockAfter) {
				t.Fatalf("refusal changed lock: %v", readErr)
			}
			for _, entry := range planned[:2] {
				for _, relative := range entry.Paths {
					if _, err := os.Lstat(filepath.Join(repo, filepath.FromSlash(relative))); err != nil {
						t.Fatalf("refusal deleted %s: %v", relative, err)
					}
				}
			}
		})
	}
	t.Run("deletes a claude tree and orphaned links and lock entries", func(t *testing.T) {
		repo := t.TempDir()
		content := "# installed\n"
		hash := externalSkillsLockDigest([]restoreFile{{Path: "SKILL.md", Content: []byte(content)}})
		writeRetiredSkillTestFile(t, repo, ".claude/skills/fixing-metadata/SKILL.md", content)
		writeRetiredSkillTestLink(t, repo, ".claude/skills/tdd", "../../.agents/skills/tdd")
		writeRetiredSkillTestLock(t, repo, map[string]any{"fixing-metadata": map[string]any{"computedHash": hash}, "the-fool": map[string]any{}})
		planned, err := InstalledRetiredSkills(repo)
		if err != nil {
			t.Fatal(err)
		}
		deleted, err := DeleteRetiredSkills(t.Context(), repo, planned)
		if err != nil || !reflect.DeepEqual(deleted, planned) || len(deleted) != 3 {
			t.Fatalf("orphaned deletion = %+v, %v", deleted, err)
		}
		remaining, err := InstalledRetiredSkills(repo)
		if err != nil || len(remaining) != 0 {
			t.Fatalf("remaining retired skills = %+v, %v", remaining, err)
		}
	})

}

func writeRetiredSkillTestFile(t *testing.T, repo, relative, content string) {
	t.Helper()
	filename := filepath.Join(repo, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filename, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeRetiredSkillTestLink(t *testing.T, repo, relative, target string) {
	t.Helper()
	filename := filepath.Join(repo, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, filename); err != nil {
		t.Fatal(err)
	}
}

func writeRetiredSkillTestLock(t *testing.T, repo string, entries map[string]any) {
	t.Helper()
	data, err := json.Marshal(map[string]any{"version": 1, "skills": entries})
	if err != nil {
		t.Fatal(err)
	}
	writeRetiredSkillTestFile(t, repo, "skills-lock.json", string(data))
}
