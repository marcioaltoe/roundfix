// Suite: derived-artifact regeneration gates
// Invariant: the declared regeneration commands rewrite exactly the artifacts
// their ownership records sanction, and leave frozen ones byte-identical.
// Boundary IN: nested go and make invocations inside a copy of the repository.
// Boundary OUT: ownership record parsing and validation, which stay untagged
// in derived_ownership_test.go.
//
// The repocontract tag keeps these out of go test ./...: each copies the whole
// repository tree, which makes every file under internal/ one of their inputs,
// so any code change anywhere re-runs them. Their verdict changes only when
// Baseline assets or owned skills change, so make verify-docs runs them at the
// pull request boundary instead.

//go:build repocontract

package baseline

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestMeasuredSanctionedOwnershipMatchesRecords(t *testing.T) {
	// Parallel with the package: the fixture below is this test's own copy.
	// Subtests stay sequential, because child commands rewrite that one
	// shared fixture. Serialising the two regeneration giants against each
	// other was measured slower (91s vs 71s package wall on 2026-08-10): the
	// second giant then runs as a solo tail after the package drains.
	t.Parallel()
	repository := newDerivedRegenerationFixture(t)
	baselineRoot := filepath.Join(repository, "internal", "baseline")
	fileSystem := os.DirFS(baselineRoot)
	roots := derivedDigestScanRoots(t)
	records, err := readDerivedOwnershipRecords(fileSystem, roots)
	if err != nil {
		t.Fatal(err)
	}
	resolved, err := validateDerivedOwnership(fileSystem, roots)
	if err != nil {
		t.Fatal(err)
	}
	cacheRoot := ambientGoBuildCache(t)
	clean := snapshotDerivedTreeAt(t, baselineRoot, roots)
	sanctionedProbes := declaredSanctionedProbes(t, fileSystem, roots, records, resolved)
	if err := exerciseDeclaredRegenerationStep(
		t.Context(), repository, baselineRoot, cacheRoot, roots, clean,
		"make baseline-digests", sanctionedProbes, nil,
	); err != nil {
		t.Fatal(err)
	}

	measured := make(map[string]struct{})
	for _, probe := range sanctionedProbes {
		measured[probe.path] = struct{}{}
	}

	skillPath := filepath.Join(repository, ".agents", "skills", "qa-gate", "SKILL.md")
	skill, err := os.ReadFile(skillPath)
	if err != nil {
		t.Fatalf("read owned Skill fixture: %v", err)
	}
	skill = append(skill, []byte("\n<!-- derived ownership measurement fixture -->\n")...)
	if err := os.WriteFile(skillPath, skill, 0o644); err != nil {
		t.Fatalf("edit owned Skill fixture: %v", err)
	}

	if err := runDeclaredRegenerationStep(t.Context(), repository, cacheRoot, "make skills-sync"); err != nil {
		t.Fatal(err)
	}
	before := snapshotDerivedTreeAt(t, baselineRoot, roots)
	if err := runDeclaredRegenerationStep(t.Context(), repository, cacheRoot, "make baseline-digests"); err != nil {
		t.Fatal(err)
	}
	after := snapshotDerivedTreeAt(t, baselineRoot, roots)

	for artifactPath, afterArtifact := range after {
		beforeArtifact, existed := before[artifactPath]
		if existed && reflect.DeepEqual(afterArtifact, beforeArtifact) {
			continue
		}
		measured[artifactPath] = struct{}{}
	}
	for artifactPath := range before {
		if _, exists := after[artifactPath]; exists {
			continue
		}
		measured[artifactPath] = struct{}{}
	}

	for artifactPath := range measured {
		record, ok := resolved[artifactPath]
		if !ok {
			t.Errorf("measured sanctioned rewrite %q has no resolved ownership", artifactPath)
			continue
		}
		if record.Owner != derivedOwnerSanctioned {
			t.Errorf("measured sanctioned rewrite %q resolves to owner %q", artifactPath, record.Owner)
		}
	}

	secondOutput, err := runDeclaredRegenerationStepOutput(
		t.Context(), repository, cacheRoot, "make baseline-digests",
	)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(secondOutput, []byte(`"changed":false`)) {
		t.Fatalf("second sanctioned run output = %q, want changed:false", secondOutput)
	}
	second := snapshotDerivedTreeAt(t, baselineRoot, roots)
	if !reflect.DeepEqual(second, after) {
		t.Fatal("second sanctioned run changed the derived tree")
	}
}

func TestDeclaredStepRegenerationAndFrozenBoundaries(t *testing.T) {
	// Parallel with the package: the fixture below is this test's own copy.
	// Subtests stay sequential, because child commands deliberately rewrite
	// that one shared fixture.
	t.Parallel()
	repository := newDerivedRegenerationFixture(t)
	baselineRoot := filepath.Join(repository, "internal", "baseline")
	fileSystem := os.DirFS(baselineRoot)
	roots := derivedDigestScanRoots(t)
	records, err := readDerivedOwnershipRecords(fileSystem, roots)
	if err != nil {
		t.Fatal(err)
	}
	resolved, err := validateDerivedOwnership(fileSystem, roots)
	if err != nil {
		t.Fatal(err)
	}
	artifactsByRecord := derivedArtifactsByRecord(t, fileSystem, roots, records, resolved)
	clean := snapshotDerivedTreeAt(t, baselineRoot, roots)
	cacheRoot := ambientGoBuildCache(t)

	sanctioned := declaredSanctionedProbes(t, fileSystem, roots, records, resolved)
	frozen := declaredFrozenProbes(t, fileSystem, roots, records, resolved)

	recordPaths := make([]string, 0, len(records))
	for recordPath := range records {
		recordPaths = append(recordPaths, recordPath)
	}
	sort.Strings(recordPaths)
	var dedicatedRecordPath string
	var dedicatedCommand string
	for _, recordPath := range recordPaths {
		record := records[recordPath]
		artifacts := filterDerivedArtifactsByOwner(
			artifactsByRecord[recordPath],
			resolved,
			derivedOwnerDedicated,
		)
		if len(artifacts) == 0 {
			continue
		}
		if dedicatedRecordPath == "" {
			dedicatedRecordPath = recordPath
			dedicatedCommand = record.Command
		}
		t.Run("dedicated/"+strings.ReplaceAll(recordPath, "/", "_"), func(t *testing.T) {
			rewrite := probesForPaths(artifacts, "dedicated")
			err := exerciseDeclaredRegenerationStep(
				t.Context(), repository, baselineRoot, cacheRoot, roots, clean,
				record.Command, rewrite, frozen,
			)
			if err != nil {
				t.Fatal(err)
			}
			assertDerivedFixtureRestored(t, baselineRoot, roots, clean)
		})
	}
	if dedicatedRecordPath == "" {
		const planCommand = "go test ./internal/baseline -count=1 -run TestBaselinePlanCharacterization -update-baseline-plan-characterization"
		dedicatedRecordPath = "testdata/plan-characterization/_ownership.yml"
		dedicatedCommand = planCommand
		artifacts := artifactsByRecord[dedicatedRecordPath]
		if len(artifacts) == 0 {
			t.Fatalf("synthetic dedicated ownership record %q governs no artifacts", dedicatedRecordPath)
		}
		t.Run("dedicated/synthetic_plan_characterization", func(t *testing.T) {
			writeDedicatedCommandFixture(t, baselineRoot, dedicatedRecordPath, dedicatedCommand)
			err := exerciseDeclaredRegenerationStep(
				t.Context(), repository, baselineRoot, cacheRoot, roots, clean,
				dedicatedCommand, probesForPaths(artifacts, "dedicated"), frozen,
			)
			if err != nil {
				t.Fatal(err)
			}
			assertDerivedFixtureRestored(t, baselineRoot, roots, clean)
		})
	}

	negativeCases := []struct {
		name      string
		command   string
		wantError string
	}{
		{
			name:      "command does not exist",
			command:   "roundfix-declared-step-command-does-not-exist",
			wantError: "run declared command",
		},
		{
			name:      "declared flag is wrong",
			command:   dedicatedCommand + " -roundfix-deliberately-wrong-flag",
			wantError: "run declared command",
		},
		{
			name:      "command leaves artifacts unchanged",
			command:   "go version",
			wantError: "unchanged after deliberate perturbation",
		},
	}
	for _, test := range negativeCases {
		t.Run("failure/"+test.name, func(t *testing.T) {
			writeDedicatedCommandFixture(t, baselineRoot, dedicatedRecordPath, test.command)
			record, err := readDerivedOwnershipRecord(fileSystem, dedicatedRecordPath)
			if err != nil {
				t.Fatal(err)
			}
			err = exerciseDeclaredRegenerationStep(
				t.Context(), repository, baselineRoot, cacheRoot, roots, clean,
				record.Command,
				probesForPaths(artifactsByRecord[dedicatedRecordPath], "dedicated"),
				nil,
			)
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Fatalf("declared-step error = %v, want %q", err, test.wantError)
			}
			assertDerivedFixtureRestored(t, baselineRoot, roots, clean)
		})
	}

	t.Run("sanctioned command rewrites sanctioned and preserves frozen", func(t *testing.T) {
		err := exerciseDeclaredRegenerationStep(
			t.Context(), repository, baselineRoot, cacheRoot, roots, clean,
			"make baseline-digests", sanctioned, frozen,
		)
		if err != nil {
			t.Fatal(err)
		}
		assertDerivedFixtureRestored(t, baselineRoot, roots, clean)
	})

	t.Run("frozen resolved path rejects rewrite", func(t *testing.T) {
		const (
			recordPath   = "testdata/catalog.diagnostics.golden.json_ownership.yml"
			artifactPath = "testdata/catalog.diagnostics.golden.json"
		)
		content := []byte("owner: frozen\nreason: deliberately frozen rewrite fixture\n")
		if err := os.WriteFile(
			filepath.Join(baselineRoot, filepath.FromSlash(recordPath)),
			content,
			0o644,
		); err != nil {
			t.Fatalf("write frozen ownership fixture %q: %v", recordPath, err)
		}
		fixtureResolved, err := validateDerivedOwnership(fileSystem, roots)
		if err != nil {
			t.Fatal(err)
		}
		record, ok := fixtureResolved[artifactPath]
		if !ok {
			t.Fatalf("frozen fixture artifact %q has no resolved ownership", artifactPath)
		}
		if record.Owner != derivedOwnerFrozen {
			t.Fatalf("frozen fixture artifact %q owner = %q", artifactPath, record.Owner)
		}
		err = exerciseDeclaredRegenerationStep(
			t.Context(), repository, baselineRoot, cacheRoot, roots, clean,
			"make baseline-digests", nil,
			[]derivedArtifactProbe{{path: artifactPath, owner: string(record.Owner)}},
		)
		if err == nil || !strings.Contains(err.Error(), "rewrote frozen artifact") {
			t.Fatalf("frozen rewrite error = %v, want rewritten frozen artifact failure", err)
		}
		assertDerivedFixtureRestored(t, baselineRoot, roots, clean)
	})
}
