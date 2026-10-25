package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func copyThisRepositorySkillSetPath(t *testing.T, source, destination string) {
	t.Helper()
	info, err := os.Stat(source)
	if err != nil {
		t.Fatalf("stat fixture source %s: %v", source, err)
	}
	if info.IsDir() {
		if err := filepath.Walk(source, func(path string, info os.FileInfo, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			relative, err := filepath.Rel(source, path)
			if err != nil {
				return err
			}
			target := filepath.Join(destination, relative)
			if info.IsDir() {
				return os.MkdirAll(target, info.Mode().Perm())
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			return os.WriteFile(target, data, info.Mode().Perm())
		}); err != nil {
			t.Fatalf("copy fixture directory %s: %v", source, err)
		}
		return
	}
	data, err := os.ReadFile(source)
	if err != nil {
		t.Fatalf("read fixture source %s: %v", source, err)
	}
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(destination, data, info.Mode().Perm()); err != nil {
		t.Fatalf("write fixture destination %s: %v", destination, err)
	}
}
