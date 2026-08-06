//go:build integration

package skills

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// Suite: owned skill edit verification
// Invariant: editing owned skill bytes leaves the exact repository verification gate green without regeneration.
// Boundary IN: an isolated tracked repository and its real make verify target.
// Boundary OUT: ordinary package tests, which run without the integration build tag.

func TestOwnedSkillEditLeavesMakeVerifyGreen(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join(".."))
	verificationRoot := copyTrackedRepository(t, repoRoot)
	initializeTrackedRepository(t, verificationRoot)
	for _, relative := range []string{
		filepath.Join(".agents", "skills", "roundfix", "SKILL.md"),
		filepath.Join("skills", "roundfix", "SKILL.md"),
	} {
		path := filepath.Join(verificationRoot, relative)
		file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0)
		if err != nil {
			t.Fatalf("open owned skill edit target %s: %v", relative, err)
		}
		if _, err := file.WriteString("\n<!-- compatibility-preserving owned skill edit -->\n"); err != nil {
			_ = file.Close()
			t.Fatalf("edit owned skill %s: %v", relative, err)
		}
		if err := file.Close(); err != nil {
			t.Fatalf("close owned skill edit target %s: %v", relative, err)
		}
	}

	derivedPaths := []string{
		"internal/baseline/assets/setups",
		"internal/baseline/testdata",
		"internal/baseline/assets/source-baselines",
		"internal/baseline/assets/formatter-fixtures",
		"internal/baseline/assets/profiles",
	}
	derivedBefore := trackedPathsDigest(t, verificationRoot, derivedPaths)
	archivedBefore := trackedPathsDigest(t, verificationRoot, []string{"docs/specs/_archived"})

	command := exec.Command("make", "verify", "RTK=")
	command.Dir = verificationRoot
	command.Env = append(os.Environ(), "GOCACHE="+t.TempDir())
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("make verify after owned skill edit: %v\n%s", err, tailBytes(output, 32*1024))
	}
	if bytes.Contains(output, []byte("baseline-digests")) {
		t.Fatalf("make verify invoked digest regeneration after an owned skill edit:\n%s", tailBytes(output, 32*1024))
	}
	if got := trackedPathsDigest(t, verificationRoot, derivedPaths); got != derivedBefore {
		t.Fatal("make verify changed derived digest artifacts after an owned skill edit")
	}
	if got := trackedPathsDigest(t, verificationRoot, []string{"docs/specs/_archived"}); got != archivedBefore {
		t.Fatal("make verify changed archived Spec artifacts after an owned skill edit")
	}
}
