package skills

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Suite: owned skill artifact stability
// Invariant: editing owned skill bytes leaves every derived Baseline artifact and characterization corpus byte-identical.
// Boundary IN: an isolated tracked repository's owned skill files and declared derived artifact roots.
// Boundary OUT: the repository-wide verification gate and owned-skill readiness behavior covered by their canonical suites.

func derivedDigestPaths(t *testing.T, repoRoot string) []string {
	t.Helper()

	data, err := os.ReadFile(filepath.Join(repoRoot, "Makefile"))
	if err != nil {
		t.Fatalf("read Makefile DERIVED_DIGEST_PATHS: %v", err)
	}
	const assignment = "DERIVED_DIGEST_PATHS :="
	for _, line := range strings.Split(string(data), "\n") {
		if !strings.HasPrefix(line, assignment) {
			continue
		}
		paths := strings.Fields(strings.TrimSpace(strings.TrimPrefix(line, assignment)))
		if len(paths) == 0 {
			t.Fatal("Makefile DERIVED_DIGEST_PATHS is empty")
		}
		return paths
	}
	t.Fatal("Makefile has no DERIVED_DIGEST_PATHS assignment")
	return nil
}

func artifactBytes(t *testing.T, repoRoot string, relativeRoots []string) map[string][]byte {
	t.Helper()

	artifacts := make(map[string][]byte)
	for _, relativeRoot := range relativeRoots {
		root := filepath.Join(repoRoot, filepath.FromSlash(relativeRoot))
		if err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() {
				return nil
			}
			relative, err := filepath.Rel(repoRoot, path)
			if err != nil {
				return err
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			artifacts[filepath.ToSlash(relative)] = data
			return nil
		}); err != nil {
			t.Fatalf("read artifact bytes under %s: %v", relativeRoot, err)
		}
	}
	return artifacts
}

func assertArtifactBytesEqual(t *testing.T, kind string, before map[string][]byte, after map[string][]byte) {
	t.Helper()

	for path, beforeBytes := range before {
		afterBytes, exists := after[path]
		if !exists {
			t.Fatalf("%s %s was removed by an owned skill edit", kind, path)
		}
		if !bytes.Equal(afterBytes, beforeBytes) {
			t.Fatalf("%s %s changed after an owned skill edit", kind, path)
		}
	}
	for path := range after {
		if _, existed := before[path]; !existed {
			t.Fatalf("%s %s was created by an owned skill edit", kind, path)
		}
	}
}

func ambientGoBuildCacheForSkills(t *testing.T) string {
	t.Helper()
	if dir := strings.TrimSpace(os.Getenv("GOCACHE")); dir != "" {
		return dir
	}
	output, err := exec.CommandContext(t.Context(), "go", "env", "GOCACHE").Output()
	if err != nil {
		t.Fatalf("resolve ambient go build cache: %v", err)
	}
	dir := strings.TrimSpace(string(output))
	if dir == "" {
		t.Fatal("resolve ambient go build cache: empty path")
	}
	return dir
}
