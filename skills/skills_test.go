package skills

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
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

func TestFilesIncludesBothSkillsAndAgentMetadata(t *testing.T) {
	files, err := Files()
	if err != nil {
		t.Fatalf("read skill files: %v", err)
	}

	var paths []string
	for _, file := range files {
		paths = append(paths, file.Path)
		if len(file.Data) == 0 {
			t.Fatalf("expected embedded data for %s", file.Path)
		}
	}

	expected := []string{
		"roundfix-resolve-round/SKILL.md",
		"roundfix-resolve-round/agents/openai.yaml",
		"roundfix-watch/SKILL.md",
		"roundfix-watch/agents/openai.yaml",
	}
	if !reflect.DeepEqual(paths, expected) {
		t.Fatalf("expected embedded paths %#v, got %#v", expected, paths)
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
	for _, target := range result.Targets {
		if target.Files != 4 {
			t.Fatalf("expected four files for %s, got %d", target.Target, target.Files)
		}
		for _, path := range []string{
			"roundfix-watch/SKILL.md",
			"roundfix-watch/agents/openai.yaml",
			"roundfix-resolve-round/SKILL.md",
			"roundfix-resolve-round/agents/openai.yaml",
		} {
			if _, err := os.Stat(filepath.Join(target.Dir, path)); err != nil {
				t.Fatalf("expected installed file %s for %s: %v", path, target.Target, err)
			}
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
