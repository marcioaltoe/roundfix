package skills

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckValidatesRoundfixSkillArtifacts(t *testing.T) {
	if diagnostics := Check(); len(diagnostics) > 0 {
		var messages []string
		for _, diagnostic := range diagnostics {
			messages = append(messages, diagnostic.Path+": "+diagnostic.Message)
		}
		t.Fatalf("expected no skill diagnostics, got %s", strings.Join(messages, "\n"))
	}
}

func TestCheckOpenAIManifestRequiresEntrypointAndRuntimeCommand(t *testing.T) {
	diagnostics := checkOpenAIManifest("roundfix/agents/openai.yaml", []byte(`
name: roundfix
runtime_hints: {}
`))

	messages := make([]string, 0, len(diagnostics))
	for _, diagnostic := range diagnostics {
		messages = append(messages, diagnostic.Message)
	}
	joined := strings.Join(messages, "\n")
	for _, expected := range []string{
		"manifest field entrypoint is required",
		"manifest runtime command is required",
	} {
		if !strings.Contains(joined, expected) {
			t.Fatalf("expected diagnostic %q, got %s", expected, joined)
		}
	}
}

func TestCheckOpenAIManifestAcceptsNestedRuntimeHints(t *testing.T) {
	diagnostics := checkOpenAIManifest("roundfix/agents/openai.yaml", []byte(`
name: roundfix
entrypoint: SKILL.md
runtime:
  hints:
    command: roundfix watch --source coderabbit --pr <number> --agent <agent> --until-clean
`))

	if len(diagnostics) > 0 {
		t.Fatalf("expected nested runtime hints to pass, got %#v", diagnostics)
	}
}

func TestFilesIncludeEveryOwnedSkill(t *testing.T) {
	files, err := Files()
	if err != nil {
		t.Fatalf("read skill files: %v", err)
	}

	present := make(map[string]bool, len(files))
	for _, file := range files {
		present[file.Path] = true
		if len(file.Data) == 0 {
			t.Fatalf("expected embedded data for %s", file.Path)
		}
	}

	// The operational roundfix skill ships its SKILL.md and OpenAI manifest,
	// and every owned skill ships a SKILL.md.
	required := []string{"roundfix/SKILL.md", "roundfix/agents/openai.yaml"}
	for _, skill := range Names() {
		required = append(required, skill+"/SKILL.md")
	}
	for _, path := range required {
		if !present[path] {
			t.Fatalf("expected embedded file %q, got %d files without it", path, len(files))
		}
	}
}

func TestFrontmatterName(t *testing.T) {
	cases := []struct {
		name string
		text string
		want string
	}{
		{"valid", "---\nname: write-prd\ndescription: x\n---\nbody", "write-prd"},
		{"no frontmatter", "# Heading\nbody", ""},
		{"unterminated", "---\nname: write-prd\nbody", ""},
		{"missing name", "---\ndescription: x\n---\nbody", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := frontmatterName(tc.text); got != tc.want {
				t.Fatalf("frontmatterName(%q) = %q, want %q", tc.text, got, tc.want)
			}
		})
	}
}

func TestRecommendedListsExternallyManagedSkills(t *testing.T) {
	recommended := Recommended()
	if len(recommended) == 0 {
		t.Fatal("expected recommended skills from skills-lock.json, got none")
	}
	owned := make(map[string]bool, len(Names()))
	for _, skill := range Names() {
		owned[skill] = true
	}
	for _, skill := range recommended {
		if owned[skill] {
			t.Fatalf("recommended skill %q must not be an owned bundle skill", skill)
		}
	}
}

func TestInstallCopiesSkillsToSupportedTargetDirectories(t *testing.T) {
	root := t.TempDir()
	targetDirs := map[string]string{
		"codex":    filepath.Join(root, "codex"),
		"claude":   filepath.Join(root, "claude"),
		"opencode": filepath.Join(root, "opencode"),
	}

	result, err := Install(context.Background(), InstallRequest{
		Target:     "all",
		TargetDirs: targetDirs,
	})
	if err != nil {
		t.Fatalf("install skills: %v", err)
	}
	if len(result.Targets) != 3 {
		t.Fatalf("expected three install targets, got %#v", result.Targets)
	}
	wantFiles := embeddedFileCount(t)
	for _, target := range result.Targets {
		if target.Files != wantFiles {
			t.Fatalf("expected %d files for %s, got %d", wantFiles, target.Target, target.Files)
		}
		for _, path := range []string{
			"roundfix/SKILL.md",
			"roundfix/agents/openai.yaml",
			"write-prd/SKILL.md",
		} {
			if _, err := os.Stat(filepath.Join(target.Dir, path)); err != nil {
				t.Fatalf("expected installed file %s for %s: %v", path, target.Target, err)
			}
		}
	}
}

func embeddedFileCount(t *testing.T) int {
	t.Helper()
	files, err := Files()
	if err != nil {
		t.Fatalf("read skill files: %v", err)
	}
	return len(files)
}

func TestInstallCopiesSkillsToProjectDirectoryByDefault(t *testing.T) {
	projectDir := t.TempDir()

	result, err := Install(context.Background(), InstallRequest{ProjectDir: projectDir})
	if err != nil {
		t.Fatalf("install project skill: %v", err)
	}
	if len(result.Targets) != 1 {
		t.Fatalf("expected one install target, got %#v", result.Targets)
	}
	target := result.Targets[0]
	if target.Target != "project" {
		t.Fatalf("expected project target, got %q", target.Target)
	}
	expectedDir := filepath.Join(projectDir, ".agents", "skills")
	if target.Dir != expectedDir {
		t.Fatalf("expected project target dir %q, got %q", expectedDir, target.Dir)
	}
	if want := embeddedFileCount(t); target.Files != want {
		t.Fatalf("expected %d installed files, got %d", want, target.Files)
	}
	for _, path := range []string{
		"roundfix/SKILL.md",
		"roundfix/agents/openai.yaml",
		"write-prd/SKILL.md",
	} {
		if _, err := os.Stat(filepath.Join(expectedDir, path)); err != nil {
			t.Fatalf("expected installed file %s: %v", path, err)
		}
	}
}

func TestInstallRejectsUnsupportedTarget(t *testing.T) {
	_, err := Install(context.Background(), InstallRequest{Target: "other"})
	if err == nil {
		t.Fatal("expected unsupported target error")
	}
	if !strings.Contains(err.Error(), "unsupported skill install target") {
		t.Fatalf("expected unsupported target error, got %v", err)
	}
}
