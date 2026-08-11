// Suite: owned skill edit regeneration gate
// Invariant: an owned skill edit regenerates its derived artifacts and leaves
// every frozen artifact byte-identical.
// Boundary IN: a copy of the tracked repository and one nested make run.
// Boundary OUT: skill contract assertions that read only .agents/skills.
//
// The repocontract tag keeps this out of go test ./...: it copies every
// tracked file, which makes the whole repository its input, so any change
// anywhere re-runs it. make verify-docs runs it at the pull request boundary.

//go:build repocontract

package skills

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

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
	archivedBefore := artifactBytes(t, verificationRoot, []string{"_archived/specs"})

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
	// Reuse the suite's own warm build cache: the contract proven here is
	// byte-identical regeneration output after an owned skill edit, not a
	// cold compile. Result caching stays content-governed — the skill edit
	// below changes the files those nested tests read, so the steps that
	// consume the edit re-execute, and only checks over identical bytes may
	// replay their verdict.
	command.Env = append(os.Environ(),
		"GOCACHE="+ambientGoBuildCacheForSkills(t),
		"GOFLAGS=-buildvcs=false",
	)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("regenerate derived artifacts after owned skill edit: %v\n%s", err, tailBytes(output, 32*1024))
	}

	assertArtifactBytesEqual(t, "derived Baseline artifact", derivedBefore, artifactBytes(t, verificationRoot, derivedPaths))
	assertArtifactBytesEqual(t, "characterization corpus artifact", characterizationBefore, artifactBytes(t, verificationRoot, characterizationPaths))
	assertArtifactBytesEqual(t, "archived Spec artifact", archivedBefore, artifactBytes(t, verificationRoot, []string{"_archived/specs"}))
}
