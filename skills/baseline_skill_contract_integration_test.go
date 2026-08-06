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

func TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join(".."))
	verificationRoot := copyTrackedRepository(t, repoRoot)
	derivedPaths := derivedDigestPaths(t, verificationRoot)
	characterizationPaths := []string{
		"internal/baseline/testdata/catalog.diagnostics.golden.json",
		"internal/baseline/testdata/plan-characterization",
	}
	derivedBefore := artifactBytes(t, verificationRoot, derivedPaths)
	characterizationBefore := artifactBytes(t, verificationRoot, characterizationPaths)
	archivedBefore := artifactBytes(t, verificationRoot, []string{"docs/specs/_archived"})

	for _, relative := range []string{
		filepath.Join(".agents", "skills", "roundfix", "SKILL.md"),
		filepath.Join("skills", "roundfix", "SKILL.md"),
	} {
		path := filepath.Join(verificationRoot, relative)
		file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0)
		if err != nil {
			t.Fatalf("open owned skill edit target %s: %v", relative, err)
		}
		_, writeErr := file.WriteString("\n<!-- compatibility-preserving owned skill edit -->\n")
		closeErr := file.Close()
		if writeErr != nil || closeErr != nil {
			t.Fatalf("edit owned skill %s: write error = %v, close error = %v", relative, writeErr, closeErr)
		}
	}

	command := exec.CommandContext(t.Context(), "make", "baseline-digests")
	command.Dir = verificationRoot
	command.Env = append(os.Environ(),
		"GOCACHE="+filepath.Join(verificationRoot, ".gocache"),
		"GOFLAGS=-buildvcs=false",
	)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("regenerate derived artifacts after owned skill edit: %v\n%s", err, tailBytes(output, 32*1024))
	}

	assertArtifactBytesEqual(t, "derived Baseline artifact", derivedBefore, artifactBytes(t, verificationRoot, derivedPaths))
	assertArtifactBytesEqual(t, "characterization corpus artifact", characterizationBefore, artifactBytes(t, verificationRoot, characterizationPaths))
	assertArtifactBytesEqual(t, "archived Spec artifact", archivedBefore, artifactBytes(t, verificationRoot, []string{"docs/specs/_archived"}))
}

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
