package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	roundskills "roundfix/skills"
)

func TestSkillsRemoveListsThenRemovesOnlyOwnedCopies(t *testing.T) {
	t.Parallel()
	t.Run("surface_transcript_4", func(t *testing.T) {
		home, root := authoringReportFixture(t, false)
		driftAuthoringReportFixture(t, root)
		path := filepath.Join(home, ".claude", "skills", "write-prd")
		writeStaleUserSkillFixture(t, path, "write-prd", bundledUserSkillVersion(t, "write-prd"))
		before, err := os.ReadFile(filepath.Join(path, "SKILL.md"))
		if err != nil {
			t.Fatal(err)
		}
		assertAuthoringCommand(t, []string{"skills", "remove", "--target", "claude"}, "Would remove "+path+"\nNothing was removed; run the command again with --apply.\n", "", 0)
		after, err := os.ReadFile(filepath.Join(path, "SKILL.md"))
		if err != nil || !bytes.Equal(before, after) {
			t.Fatalf("listing changed copy: %q %v", after, err)
		}
	})
	t.Run("surface_transcript_5", func(t *testing.T) {
		home, _ := withCLIWorkspace(t)
		path := filepath.Join(home, ".claude", "skills", "write-prd")
		writeStaleUserSkillFixture(t, path, "write-prd", bundledUserSkillVersion(t, "write-prd"))
		assertAuthoringCommand(t, []string{"skills", "remove", "--target", "project"}, "", "roundfix: skills remove refuses --target project: the project's copies are the Repository Skill Set\n", 2)
		if _, err := os.Stat(filepath.Join(path, "SKILL.md")); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("apply", func(t *testing.T) {
		home, repo := withCLIWorkspace(t)
		dir := filepath.Join(home, ".claude", "skills")
		folder := filepath.Join(dir, "roundfix")
		nested := filepath.Join(folder, "nested", "file")
		target := filepath.Join(home, "vault")
		targetFile := filepath.Join(target, "nested", "file")
		other := filepath.Join(dir, "not-owned", "file")
		file := filepath.Join(dir, "qa-gate")
		project := filepath.Join(repo, ".agents", "skills", "roundfix", "SKILL.md")
		for _, path := range []string{nested, targetFile, other, file, project} {
			writeRemovalFile(t, path)
		}
		link := filepath.Join(dir, "write-prd")
		// Keep the link's own relative text in output, rather than its resolved path.
		linkText := filepath.Join("..", "..", "vault")
		if err := os.Symlink(linkText, link); err != nil {
			t.Fatal(err)
		}
		assertAuthoringCommand(t, []string{"skills", "remove", "--target", "claude"}, "Would remove "+folder+"\nWould unlink "+link+" -> "+linkText+"\nNothing was removed; run the command again with --apply.\n", "", 0)
		for _, path := range []string{nested, targetFile, other, file, project} {
			assertRemovalFileKept(t, path)
		}
		if got, err := os.Readlink(link); err != nil || got != linkText {
			t.Fatalf("link=%q err=%v", got, err)
		}
		assertAuthoringCommand(t, []string{"skills", "remove", "--target", "claude", "--apply"}, "Removed "+folder+"\nUnlinked "+link+" -> "+linkText+"\n", "", 0)
		for _, path := range []string{folder, link} {
			if _, err := os.Lstat(path); !os.IsNotExist(err) {
				t.Fatalf("removed entry %s: %v", path, err)
			}
		}
		for _, path := range []string{targetFile, other, file, project} {
			assertRemovalFileKept(t, path)
		}
		assertAuthoringCommand(t, []string{"skills", "remove", "--target", "claude", "--apply"}, "No Roundfix-owned skill copy under "+dir+"\n", "", 0)
	})
	t.Run("target_all", func(t *testing.T) {
		home, _ := withCLIWorkspace(t)
		directories := roundskills.UserSkillDirectories(home, nil)
		for _, directory := range directories {
			writeRemovalFile(t, filepath.Join(directory.Path, "roundfix", "SKILL.md"))
		}
		var listing, applied strings.Builder
		for _, directory := range directories[:3] {
			listing.WriteString("Would remove " + filepath.Join(directory.Path, "roundfix") + "\n")
			applied.WriteString("Removed " + filepath.Join(directory.Path, "roundfix") + "\n")
		}
		listing.WriteString("Nothing was removed; run the command again with --apply.\n")
		assertAuthoringCommand(t, []string{"skills", "remove", "--target", "all"}, listing.String(), "", 0)
		for _, directory := range directories {
			assertRemovalFileKept(t, filepath.Join(directory.Path, "roundfix", "SKILL.md"))
		}
		assertAuthoringCommand(t, []string{"skills", "remove", "--target", "all", "--apply"}, applied.String(), "", 0)
		for _, directory := range directories[3:] {
			assertRemovalFileKept(t, filepath.Join(directory.Path, "roundfix", "SKILL.md"))
		}
	})
	t.Run("dir", func(t *testing.T) {
		home, _ := withCLIWorkspace(t)
		first := filepath.Join(home, ".codex", "skills", "roundfix", "SKILL.md")
		dir := filepath.Join(home, ".agents", "skills")
		overridden := filepath.Join(dir, "roundfix", "SKILL.md")
		writeRemovalFile(t, first)
		writeRemovalFile(t, overridden)
		args := []string{"skills", "remove", "--target", "codex", "--dir", dir}
		assertAuthoringCommand(t, args, "Would remove "+filepath.Dir(overridden)+"\nNothing was removed; run the command again with --apply.\n", "", 0)
		assertRemovalFileKept(t, overridden)
		assertAuthoringCommand(t, append(args, "--apply"), "Removed "+filepath.Dir(overridden)+"\n", "", 0)
		assertRemovalFileKept(t, first)
		if _, err := os.Stat(overridden); !os.IsNotExist(err) {
			t.Fatalf("override remains: %v", err)
		}
		assertAuthoringCommand(t, args, "No Roundfix-owned skill copy under "+dir+"\n", "", 0)
	})
	t.Run("refusals", func(t *testing.T) {
		home, _ := withCLIWorkspace(t)
		path := filepath.Join(home, ".claude", "skills", "roundfix", "SKILL.md")
		writeRemovalFile(t, path)
		for _, args := range [][]string{nil, {"--target"}, {"--target", ""}, {"--target", "unknown"}} {
			assertAuthoringCommand(t, append([]string{"skills", "remove"}, args...), "", "roundfix: skills remove requires --target <claude|codex|opencode|all>\n", 2)
		}
		assertAuthoringCommand(t, []string{"skills", "remove", "--target", "project", "--apply"}, "", "roundfix: skills remove refuses --target project: the project's copies are the Repository Skill Set\n", 2)
		assertAuthoringCommand(t, []string{"skills", "remove", "--target", "all", "--dir", home, "--apply"}, "", "roundfix: --dir requires a single --target value\n", 2)
		assertRemovalFileKept(t, path)
		t.Run("removal_fails", func(t *testing.T) {
			if runtime.GOOS == "windows" {
				t.Skip("Windows lacks Unix directory write permissions")
			}
			folder := filepath.Dir(path)
			if err := os.Chmod(folder, 0o555); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				if err := os.Chmod(folder, 0o755); err != nil {
					t.Error(err)
				}
			})
			var stdout, stderr bytes.Buffer
			code := runCLI(t, []string{"skills", "remove", "--target", "claude", "--apply"}, &stdout, &stderr)
			if code != 1 || stdout.Len() != 0 || !strings.HasPrefix(stderr.String(), "roundfix: skills remove failed: remove "+folder+": ") || !strings.Contains(stderr.String(), "permission denied") {
				t.Fatalf("exit=%d stdout=%q stderr=%q", code, &stdout, &stderr)
			}
			assertRemovalFileKept(t, path)
		})
	})
	t.Run("no_known_directory", func(t *testing.T) {
		home, _ := withCLIWorkspace(t)
		environment := commandEnvironmentForTest(t)
		path := filepath.Join(home, ".claude", "skills", "roundfix", "SKILL.md")
		writeRemovalFile(t, path)
		// An unresolved later target must refuse before an earlier target is touched.
		environment.dependencies.userSkillDirectories = func(home string) []roundskills.UserSkillDirectory {
			return roundskills.UserSkillDirectories(home, nil)[:1]
		}
		var stdout, stderr bytes.Buffer
		code := runWithContext(t.Context(), []string{"skills", "remove", "--target", "all", "--apply"}, &stdout, &stderr, environment)
		if code != 1 || stdout.Len() != 0 || stderr.String() != "roundfix: skills remove failed: no user-level skill directory is known for codex\n" {
			t.Fatalf("exit=%d stdout=%q stderr=%q", code, &stdout, &stderr)
		}
		assertRemovalFileKept(t, path)
		environment.dependencies.userSkillDirectories = nil
		stdout.Reset()
		stderr.Reset()
		code = runWithContext(t.Context(), []string{"skills", "remove", "--target", "claude", "--apply"}, &stdout, &stderr, environment)
		if code != 1 || stdout.Len() != 0 || stderr.String() != "roundfix: skills remove failed: no user-level skill directory is known for claude\n" {
			t.Fatalf("exit=%d stdout=%q stderr=%q", code, &stdout, &stderr)
		}
	})
	t.Run("dangling_link", func(t *testing.T) {
		home, _ := withCLIWorkspace(t)
		link := filepath.Join(home, ".claude", "skills", "write-prd")
		if err := os.MkdirAll(filepath.Dir(link), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink("missing", link); err != nil {
			t.Fatal(err)
		}
		assertAuthoringCommand(t, []string{"skills", "remove", "--target", "claude", "--apply"}, "Unlinked "+link+" -> missing\n", "", 0)
		if _, err := os.Lstat(link); !os.IsNotExist(err) {
			t.Fatalf("link remains: %v", err)
		}
	})
}

const skillsRemovalUsage = `Usage:
  roundfix skills list
  roundfix skills check
  roundfix skills install [--target <project|codex|claude|opencode|all>] [--dir <path>]
  roundfix skills remove --target <codex|claude|opencode|all> [--dir <path>] [--apply]

Commands:
  list       List the bundled Roundfix skills and recommended external skills
  check      Validate shipped Roundfix skill artifacts and, inside an adopted
             repository, the skills its Baseline requires
  install    Install the bundled Roundfix skills into this project or compatible Agent skill directories
  remove     List, and with --apply remove, the Roundfix skills in a user-level Agent skill directory

Options:
  --target   Install target. Supported: project, codex, claude, opencode, all
             project writes <repo>/.agents/skills and can link .claude/skills
             remove accepts codex, claude, opencode and all, never project
  --dir      Override the target skills directory for a single target
  --apply    With remove, delete what the listing names
`

func TestSkillsUsageNamesTheRemoval(t *testing.T) {
	t.Parallel()
	withCLIWorkspace(t)
	for _, args := range [][]string{{"skills", "--help"}, {"skills", "remove", "--help"}} {
		assertAuthoringCommand(t, args, skillsRemovalUsage, "", 0)
	}
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"--help"}, &stdout, &stderr)
	if code != 0 || stderr.Len() != 0 {
		t.Fatalf("exit=%d stderr=%q", code, &stderr)
	}
	lines := strings.Split(stdout.String(), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "  roundfix skills install ") {
			if i+1 >= len(lines) || lines[i+1] != "  roundfix skills remove --target <codex|claude|opencode|all> [--dir <path>] [--apply]" {
				t.Fatalf("top-level usage: %s", &stdout)
			}
			return
		}
	}
	t.Fatal("top-level usage lacks skills install")
}

func writeRemovalFile(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("kept content"), 0o644); err != nil {
		t.Fatal(err)
	}
}

func assertRemovalFileKept(t *testing.T, path string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil || string(data) != "kept content" {
		t.Fatalf("file %s changed: %q %v", path, data, err)
	}
}

func bundledUserSkillVersion(t *testing.T, name string) string {
	t.Helper()
	for _, owned := range roundskills.CheckReadiness().Owned {
		if owned.Skill == name {
			return owned.Minimum
		}
	}
	t.Fatalf("missing bundled version for %s", name)
	return ""
}
