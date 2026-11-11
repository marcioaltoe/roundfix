package skills

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestStaleUserSkillCopiesReportsAnOlderCopyInEveryUserDirectory(t *testing.T) {
	t.Parallel()
	home := t.TempDir()
	directories := UserSkillDirectories(home, nil)
	if len(directories) != 5 {
		t.Fatalf("directories = %v, want five", directories)
	}
	relative := []string{".claude/skills", ".codex/skills", ".opencode/skills", ".agents/skills", ".config/opencode/skills"}
	refreshes := []string{
		"roundfix skills install --target claude", "roundfix skills install --target codex", "roundfix skills install --target opencode",
		"roundfix skills install --target codex --dir " + filepath.Join(home, ".agents", "skills"),
		"roundfix skills install --target opencode --dir " + filepath.Join(home, ".config", "opencode", "skills"),
	}
	var want []StaleUserSkillCopy
	for i, rel := range relative {
		path := filepath.Join(home, filepath.FromSlash(rel), "roundfix")
		writeUserSkillCopy(t, path, "roundfix", "0.0.1")
		want = append(want, StaleUserSkillCopy{Skill: "roundfix", Path: path, Found: "0.0.1", Embedded: ownedSkillMinimumVersions["roundfix"], Refresh: refreshes[i]})
	}
	claude := directories[0].Path
	writeUserSkillCopy(t, filepath.Join(claude, "qa-gate"), "qa-gate", ownedSkillMinimumVersions["qa-gate"])
	newer, ok := parseSkillVersion(ownedSkillMinimumVersions["write-tasks"])
	if !ok {
		t.Fatal("embedded write-tasks version is invalid")
	}
	writeUserSkillCopy(t, filepath.Join(claude, "write-tasks"), "write-tasks", fmt.Sprintf("%d.%d.%d", newer[0], newer[1], newer[2]+1))
	writeUserSkillCopy(t, filepath.Join(claude, "write-prd"), "write-prd", "")
	writeUserSkillCopy(t, filepath.Join(claude, "not-owned"), "not-owned", "0.0.1")
	// Malformed declarations, dangling links, non-directories, and missing files are skipped.
	writeUserSkillCopy(t, filepath.Join(claude, "archive-spec"), "archive-spec", "invalid")
	malformed := filepath.Join(claude, "business-analyst")
	writeUserSkillCopy(t, malformed, "business-analyst", "0.0.1")
	if err := os.WriteFile(filepath.Join(malformed, "SKILL.md"), []byte("---\nversion: [\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(home, "missing"), filepath.Join(claude, "brainstorming")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(claude, "evidence-gate"), []byte("not a directory"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(claude, "implement-task"), 0o755); err != nil {
		t.Fatal(err)
	}
	if got := StaleUserSkillCopies(context.Background(), directories); !reflect.DeepEqual(got, want) {
		t.Fatalf("copies = %#v, want %#v", got, want)
	}
	if got := StaleUserSkillRefreshes(append(want, want[0])); !reflect.DeepEqual(got, refreshes) {
		t.Fatalf("refreshes = %v, want %v", got, refreshes)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if got := StaleUserSkillCopies(ctx, directories); len(got) != 0 {
		t.Fatalf("cancelled scan = %v", got)
	}
}

func TestStaleUserSkillCopiesFollowsLinksAndReportsAFolderOnce(t *testing.T) {
	t.Parallel()
	home := t.TempDir()
	directories := UserSkillDirectories(home, nil)
	vault := filepath.Join(home, "vault", "roundfix")
	writeUserSkillCopy(t, vault, "roundfix", "0.0.1")
	for _, directory := range directories[:2] {
		if err := os.MkdirAll(directory.Path, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(vault, filepath.Join(directory.Path, "roundfix")); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.MkdirAll(filepath.Dir(directories[3].Path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(directories[0].Path, directories[3].Path); err != nil {
		t.Fatal(err)
	}
	resolved, err := filepath.EvalSymlinks(vault)
	if err != nil {
		t.Fatal(err)
	}
	want := []StaleUserSkillCopy{
		{Skill: "roundfix", Path: filepath.Join(directories[0].Path, "roundfix"), Resolved: resolved, Found: "0.0.1", Embedded: ownedSkillMinimumVersions["roundfix"], Refresh: directories[0].Refresh},
	}
	got := StaleUserSkillCopies(context.Background(), directories)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("copies = %#v, want %#v", got, want)
	}
	description := fmt.Sprintf("roundfix 0.0.1 (binary %s) at %s -> %s", ownedSkillMinimumVersions["roundfix"], want[0].Path, resolved)
	if got[0].Describe() != description {
		t.Fatalf("description = %q, want %q", got[0].Describe(), description)
	}
	// A plain skill inside a linked parent is followed without a Resolved suffix.
	writeUserSkillCopy(t, filepath.Join(directories[0].Path, "qa-gate"), "qa-gate", "0.0.1")
	// Read through just the linked parent, so deduplication cannot hide it.
	got = StaleUserSkillCopies(context.Background(), directories[3:4])
	if len(got) != 2 || got[0].Resolved != "" || got[0].Path != filepath.Join(directories[3].Path, "qa-gate") {
		t.Fatalf("linked parent copies = %#v", got)
	}
}

func TestStaleUserSkillCopiesWritesNothing(t *testing.T) {
	t.Parallel()
	home := t.TempDir()
	directories := UserSkillDirectories(home, nil)
	vault := filepath.Join(home, "vault", "roundfix")
	writeUserSkillCopy(t, vault, "roundfix", "0.0.1")
	if err := os.MkdirAll(directories[0].Path, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(vault, filepath.Join(directories[0].Path, "roundfix")); err != nil {
		t.Fatal(err)
	}
	writeUserSkillCopy(t, filepath.Join(directories[2].Path, "qa-gate"), "qa-gate", "0.0.1")
	before := snapshotUserSkillHome(t, home)
	if got := StaleUserSkillCopies(context.Background(), directories); len(got) != 2 {
		t.Fatalf("copies = %v, want two", got)
	}
	after := snapshotUserSkillHome(t, home)
	if !reflect.DeepEqual(after, before) {
		t.Fatalf("home changed:\nbefore: %#v\nafter: %#v", before, after)
	}
}

func TestUserSkillDirectoriesFollowTheInstallTargetVariables(t *testing.T) {
	t.Parallel()
	home := filepath.Join(t.TempDir(), "user's home")
	variables := map[string]string{"CLAUDE_HOME": filepath.Join(home, "claude"), "CODEX_HOME": filepath.Join(home, "codex"), "OPENCODE_HOME": filepath.Join(home, "opencode")}
	want := []UserSkillDirectory{
		{filepath.Join(home, "claude", "skills"), "roundfix skills install --target claude"},
		{filepath.Join(home, "codex", "skills"), "roundfix skills install --target codex"},
		{filepath.Join(home, "opencode", "skills"), "roundfix skills install --target opencode"},
		{filepath.Join(home, ".agents", "skills"), "roundfix skills install --target codex --dir '" + strings.ReplaceAll(filepath.Join(home, ".agents", "skills"), "'", "'\\''") + "'"},
		{filepath.Join(home, ".config", "opencode", "skills"), "roundfix skills install --target opencode --dir '" + strings.ReplaceAll(filepath.Join(home, ".config", "opencode", "skills"), "'", "'\\''") + "'"},
	}
	getenv := func(name string) string {
		value, ok := variables[name]
		if !ok {
			t.Fatalf("unexpected environment read %q", name)
		}
		return value
	}
	if got := UserSkillDirectories(home, getenv); !reflect.DeepEqual(got, want) {
		t.Fatalf("directories = %#v, want %#v", got, want)
	}
	if got := UserSkillDirectories(" ", func(string) string { t.Fatal("blank home read environment"); return "" }); got != nil {
		t.Fatalf("blank home = %v", got)
	}
	defaults := UserSkillDirectories(home, nil)
	if defaults[0].Path != filepath.Join(home, ".claude", "skills") || defaults[1].Path != filepath.Join(home, ".codex", "skills") || defaults[2].Path != filepath.Join(home, ".opencode", "skills") {
		t.Fatalf("defaults = %v", defaults)
	}
}

func writeUserSkillCopy(t *testing.T, path, name, version string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(path, "SKILL.md"), []byte(fmt.Sprintf("---\nname: %s\nversion: %s\n---\n", name, version)), 0o644); err != nil {
		t.Fatal(err)
	}
}

type userSkillFileSnapshot struct {
	Path     string
	Mode     fs.FileMode
	Modified time.Time
	Content  string
}

func snapshotUserSkillHome(t *testing.T, home string) []userSkillFileSnapshot {
	t.Helper()
	var snapshot []userSkillFileSnapshot
	err := filepath.WalkDir(home, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		item := userSkillFileSnapshot{Path: path, Mode: info.Mode(), Modified: info.ModTime()}
		if entry.Type()&os.ModeSymlink != 0 {
			item.Content, err = os.Readlink(path)
		} else if !entry.IsDir() {
			var data []byte
			data, err = os.ReadFile(path)
			item.Content = string(data)
		}
		if err != nil {
			return err
		}
		snapshot = append(snapshot, item)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return snapshot
}
