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
		want = append(want, StaleUserSkillCopy{Skill: "roundfix", Path: path, Found: "0.0.1", Embedded: ownedSkillMinimumVersions["roundfix"], Refresh: refreshes[i], Remove: directories[i].Remove, Stale: true})
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
		{Skill: "roundfix", Path: filepath.Join(directories[0].Path, "roundfix"), Resolved: resolved, Found: "0.0.1", Embedded: ownedSkillMinimumVersions["roundfix"], Refresh: directories[0].Refresh, Remove: directories[0].Remove, Stale: true},
	}
	got := StaleUserSkillCopies(context.Background(), directories)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("copies = %#v, want %#v", got, want)
	}
	description := fmt.Sprintf("roundfix 0.0.1 (binary %s) at %s -> %s", ownedSkillMinimumVersions["roundfix"], want[0].Path, resolved) + "; older than this binary"
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
		{Path: filepath.Join(home, "claude", "skills"), Refresh: "roundfix skills install --target claude", Target: "claude", Remove: "roundfix skills remove --target claude --apply", HidesHeader: true},
		{Path: filepath.Join(home, "codex", "skills"), Refresh: "roundfix skills install --target codex", Target: "codex", Remove: "roundfix skills remove --target codex --apply"},
		{Path: filepath.Join(home, "opencode", "skills"), Refresh: "roundfix skills install --target opencode", Target: "opencode", Remove: "roundfix skills remove --target opencode --apply"},
	}
	for _, extra := range []struct{ target, path string }{
		{"codex", filepath.Join(home, ".agents", "skills")},
		{"opencode", filepath.Join(home, ".config", "opencode", "skills")},
	} {
		quoted := "'" + strings.ReplaceAll(extra.path, "'", "'\\''") + "'"
		want = append(want, UserSkillDirectory{Path: extra.path, Target: extra.target,
			Refresh: "roundfix skills install --target " + extra.target + " --dir " + quoted,
			Remove:  "roundfix skills remove --target " + extra.target + " --dir " + quoted + " --apply"})
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

func TestReportedUserSkillCopiesNameACurrentCopyThatHidesAHeader(t *testing.T) {
	t.Parallel()
	home := t.TempDir()
	directories := UserSkillDirectories(home, nil)
	claude, codex := directories[0].Path, directories[1].Path
	current := ownedSkillMinimumVersions["write-prd"]
	writeUserSkillCopy(t, filepath.Join(claude, "write-prd"), "write-prd", current)
	writeUserSkillCopy(t, filepath.Join(codex, "write-prd"), "write-prd", current)
	writeUserSkillCopy(t, filepath.Join(claude, "write-idea"), "write-idea", "0.0.1")
	writeUserSkillCopy(t, filepath.Join(claude, "write-techspec"), "write-techspec", "custom")
	writeUserSkillCopy(t, filepath.Join(claude, "qa-gate"), "qa-gate", ownedSkillMinimumVersions["qa-gate"])
	// Parseable frontmatter is still required even for a declared skill.
	writeUserSkillCopy(t, filepath.Join(claude, "write-tasks"), "write-tasks", "0.0.1")
	if err := os.WriteFile(filepath.Join(claude, "write-tasks", "SKILL.md"), []byte("---\nversion: [\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	shared := filepath.Join(home, "vault", "brainstorming")
	writeUserSkillCopy(t, shared, "brainstorming", ownedSkillMinimumVersions["brainstorming"])
	for _, dir := range []string{claude, codex} {
		if err := os.Symlink(shared, filepath.Join(dir, "brainstorming")); err != nil {
			t.Fatal(err)
		}
	}
	resolved, err := filepath.EvalSymlinks(shared)
	if err != nil {
		t.Fatal(err)
	}
	want := []StaleUserSkillCopy{
		{Skill: "brainstorming", Found: ownedSkillMinimumVersions["brainstorming"], Embedded: ownedSkillMinimumVersions["brainstorming"], Path: filepath.Join(claude, "brainstorming"), Resolved: resolved, Refresh: directories[0].Refresh, Remove: directories[0].Remove, HidesDeclaredHeader: true},
		{Skill: "write-idea", Found: "0.0.1", Embedded: ownedSkillMinimumVersions["write-idea"], Path: filepath.Join(claude, "write-idea"), Refresh: directories[0].Refresh, Remove: directories[0].Remove, Stale: true, HidesDeclaredHeader: true},
		{Skill: "write-prd", Found: current, Embedded: current, Path: filepath.Join(claude, "write-prd"), Refresh: directories[0].Refresh, Remove: directories[0].Remove, HidesDeclaredHeader: true},
		{Skill: "write-techspec", Found: "custom", Embedded: ownedSkillMinimumVersions["write-techspec"], Path: filepath.Join(claude, "write-techspec"), Refresh: directories[0].Refresh, Remove: directories[0].Remove, HidesDeclaredHeader: true},
	}
	before := snapshotUserSkillHome(t, home)
	copies := ReportedUserSkillCopies(t.Context(), directories, []string{"brainstorming", "write-idea", "write-prd", "write-techspec", "write-tasks"})
	if !reflect.DeepEqual(copies, want) {
		t.Fatalf("copies=%+v want=%+v", copies, want)
	}
	if !strings.HasSuffix(copies[1].Describe(), "; older than this binary and hides the declared header") || !strings.HasSuffix(copies[2].Describe(), "; hides the declared header") {
		t.Fatalf("reasons=%v", copies)
	}
	if got := UserSkillRemovals(copies); !reflect.DeepEqual(got, []string{directories[0].Remove}) {
		t.Fatalf("removals=%v", got)
	}
	if got := StaleUserSkillRefreshes(copies); !reflect.DeepEqual(got, []string{directories[0].Refresh}) {
		t.Fatalf("refreshes=%v", got)
	}
	want[1].HidesDeclaredHeader = false
	if got := ReportedUserSkillCopies(t.Context(), directories, nil); !reflect.DeepEqual(got, want[1:2]) {
		t.Fatalf("undeclared=%+v want=%+v", got, want[1:2])
	}
	if !reflect.DeepEqual(before, snapshotUserSkillHome(t, home)) {
		t.Fatal("reader changed the home")
	}
}
