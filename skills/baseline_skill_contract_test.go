package skills

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"roundfix/internal/gittest"
	"sort"
	"strings"
	"testing"
)

var updateDerivedDigests = flag.Bool("update", false, "regenerate derived digest artifacts")

const baselineDigestRegenerationHint = "run 'make baseline-digests'"

// upstreamManagedSkillTreeDigest pins the tree of every upstream-managed skill
// declared in skills-lock.json. It moves only when that declared set changes on
// purpose; three tests read this one constant so a legitimate change edits one
// line instead of three.
const upstreamManagedSkillTreeDigest = "c6d857ac719caf3d1f334f378bdada4e4e89cbc0795b7771c781d0511cb7c9cd"

type baselineDigestTargetResult struct {
	SchemaVersion *int    `json:"schemaVersion"`
	Type          *string `json:"type"`
	OK            *bool   `json:"ok"`
	Changed       *bool   `json:"changed"`
}

// Suite: Baseline skill distribution contract
// Invariant: canonical and shipped Baseline guidance is identical and invokes only the public CLI.
// Boundary IN: setup-context-driven and Roundfix owned skill guidance.
// Boundary OUT: executable setup-runtime removal and Baseline CLI behavior.

func TestBaselineSkillContract(t *testing.T) {
	t.Parallel()
	repoRoot := filepath.Clean(filepath.Join(".."))
	skillNames := []string{"setup-context-driven", "roundfix"}
	embeddedFiles, err := Files()
	if err != nil {
		t.Fatalf("read embedded skills: %v", err)
	}
	embeddedByPath := make(map[string][]byte, len(embeddedFiles))
	for _, file := range embeddedFiles {
		embeddedByPath[file.Path] = file.Data
	}

	for _, name := range skillNames {
		t.Run(name, func(t *testing.T) {
			canonicalPath := filepath.Join(repoRoot, ".agents", "skills", name, "SKILL.md")
			distributedPath := filepath.Join(repoRoot, "skills", name, "SKILL.md")
			canonical := readBaselineSkillContractFile(t, canonicalPath)
			distributed := readBaselineSkillContractFile(t, distributedPath)
			if !bytes.Equal(canonical, distributed) {
				t.Fatalf("%s canonical and distributed guidance differ; run make skills-sync", name)
			}
			embedded := embeddedByPath[filepath.ToSlash(filepath.Join(name, "SKILL.md"))]
			if !bytes.Equal(distributed, embedded) {
				t.Fatalf("%s distributed and embedded guidance differ", name)
			}
			body := baselineSkillBody(string(canonical))
			if strings.Contains(strings.ToLower(body), "python") ||
				strings.Contains(body, "context_setup.py") ||
				strings.Contains(body, "scripts/context") {
				t.Fatalf("%s skill body invokes a retired independent setup runtime", name)
			}
			for _, forbidden := range []string{"/Users/", "/home/", `C:\`} {
				if strings.Contains(body, forbidden) {
					t.Fatalf("%s skill body contains environment-specific path %q", name, forbidden)
				}
			}
		})
	}

	setup := string(readBaselineSkillContractFile(
		t,
		filepath.Join(repoRoot, ".agents", "skills", "setup-context-driven", "SKILL.md"),
	))
	for _, required := range []string{
		"only runtime authority",
		"roundfix baseline --repo . --format text",
		"roundfix baseline plan",
		"roundfix baseline apply",
		"roundfix/baseline-plan/v1",
		"roundfix/baseline-result/v1",
		"Repository Skill Set restoration",
		"Canonical asset synchronization",
		"behavioral fallback",
	} {
		if !strings.Contains(setup, required) {
			t.Fatalf("setup-context-driven skill missing %q", required)
		}
	}
}

func TestNoPythonBaselineRuntime(t *testing.T) {
	t.Parallel()
	repoRoot := filepath.Clean(filepath.Join(".."))
	for _, root := range []string{
		filepath.Join(repoRoot, ".agents", "skills", "setup-context-driven"),
		filepath.Join(repoRoot, "skills", "setup-context-driven"),
	} {
		var files []string
		if err := filepath.WalkDir(root, func(filePath string, entry os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() {
				return nil
			}
			relative, err := filepath.Rel(root, filePath)
			if err != nil {
				return err
			}
			files = append(files, filepath.ToSlash(relative))
			return nil
		}); err != nil {
			t.Fatalf("walk setup skill %s: %v", root, err)
		}
		sort.Strings(files)
		if !reflect.DeepEqual(files, []string{"SKILL.md"}) {
			t.Fatalf("setup skill %s ships non-guidance files: %v", root, files)
		}
	}

	makefile := string(readBaselineSkillContractFile(t, filepath.Join(repoRoot, "Makefile")))
	for _, forbidden := range []string{
		"setup-context-check",
		"PYTHONDONTWRITEBYTECODE",
		"python3",
	} {
		if strings.Contains(makefile, forbidden) {
			t.Fatalf("post-cutover Makefile invokes retired runtime marker %q", forbidden)
		}
	}

	for _, relative := range []string{
		"README.md",
		filepath.Join("docs", "user-guide", "commands.md"),
		filepath.Join("docs", "user-guide", "context-driven-development.md"),
		filepath.Join(".agents", "skills", "setup-context-driven", "SKILL.md"),
		filepath.Join(".agents", "skills", "roundfix", "SKILL.md"),
	} {
		content := string(readBaselineSkillContractFile(t, filepath.Join(repoRoot, relative)))
		for _, forbidden := range []string{
			"context_" + "setup.py",
			"context_" + "baseline.py",
			"python3",
			"Python fallback",
		} {
			if strings.Contains(content, forbidden) {
				t.Fatalf("%s references retired Baseline runtime marker %q", relative, forbidden)
			}
		}
	}
}

func TestBaselineDigestTargetReportsMachineReadableOutcomes(t *testing.T) {
	t.Parallel()
	t.Run("reports regenerated artifacts", func(t *testing.T) {
		derivedRoot := t.TempDir()
		writeBaselineDigestTargetFile(t, filepath.Join(derivedRoot, "before"), nil, 0o644)
		regenerator := filepath.Join(t.TempDir(), "regenerate")
		writeBaselineDigestTargetFile(t, regenerator, []byte("#!/bin/sh\nset -eu\ntouch \"$DERIVED_DIGEST_PROBE/after\"\n"), 0o700)

		result, stderr, exitCode := runBaselineDigestTarget(t, derivedRoot, regenerator)
		if exitCode != 0 {
			t.Fatalf("baseline-digests exit code = %d, want 0; stderr:\n%s", exitCode, stderr)
		}
		assertBaselineDigestTargetResult(t, result, true, true)
		if !strings.Contains(stderr, "baseline-digests: regenerated") ||
			!strings.Contains(stderr, filepath.Join(derivedRoot, "after")) {
			t.Fatalf("baseline-digests stderr does not report regenerated artifact:\n%s", stderr)
		}
	})

	t.Run("reports unchanged artifacts", func(t *testing.T) {
		derivedRoot := t.TempDir()
		writeBaselineDigestTargetFile(t, filepath.Join(derivedRoot, "before"), nil, 0o644)

		result, stderr, exitCode := runBaselineDigestTarget(t, derivedRoot, "true")
		if exitCode != 0 {
			t.Fatalf("baseline-digests exit code = %d, want 0; stderr:\n%s", exitCode, stderr)
		}
		assertBaselineDigestTargetResult(t, result, true, false)
		if !strings.Contains(stderr, "baseline-digests: no changes") {
			t.Fatalf("baseline-digests stderr does not report unchanged artifacts:\n%s", stderr)
		}
	})

	t.Run("reports regeneration errors", func(t *testing.T) {
		derivedRoot := t.TempDir()
		writeBaselineDigestTargetFile(t, filepath.Join(derivedRoot, "before"), nil, 0o644)

		result, stderr, exitCode := runBaselineDigestTarget(t, derivedRoot, "false")
		if exitCode == 0 {
			t.Fatalf("baseline-digests exit code = 0, want non-zero; stderr:\n%s", stderr)
		}
		assertBaselineDigestTargetResult(t, result, false, false)
		if !strings.Contains(stderr, "baseline-digests: regeneration failed at probe:Probe") {
			t.Fatalf("baseline-digests stderr does not identify the failed step:\n%s", stderr)
		}
	})
}

func TestThinSetupSkill(t *testing.T) {
	t.Parallel()
	files, err := Files()
	if err != nil {
		t.Fatalf("read embedded skills: %v", err)
	}
	var setupFiles []string
	for _, file := range files {
		if file.Skill == "setup-context-driven" {
			setupFiles = append(setupFiles, file.Path)
		}
	}
	sort.Strings(setupFiles)
	if !reflect.DeepEqual(setupFiles, []string{"setup-context-driven/SKILL.md"}) {
		t.Fatalf("embedded setup skill files = %v", setupFiles)
	}
	for _, diagnostic := range Check() {
		if strings.HasPrefix(diagnostic.Path, "setup-context-driven/") {
			t.Fatalf("thin setup skill diagnostic: %s: %s", diagnostic.Path, diagnostic.Message)
		}
	}

	setup := string(readBaselineSkillContractFile(
		t,
		filepath.Join("..", ".agents", "skills", "setup-context-driven", "SKILL.md"),
	))
	for _, required := range []string{
		"recipes and interpretation only",
		"no independent",
		"behavioral fallback",
		"Instruction Hierarchy",
		"exact source bytes",
		"Repository-Specific Normative Rules",
		"Only `accepted` is active",
		"`pending`, `partial`, `deferred`, and `done`",
		"Profile alignment",
		"Change Baseline",
		"repository-owned Profile adaptation",
		"Decline",
		"without writing",
		"--profile-file",
		"mutually exclusive",
		"Generate a fresh plan",
	} {
		if !strings.Contains(setup, required) {
			t.Errorf("thin setup skill missing %q", required)
		}
	}
	for _, forbidden := range []string{
		"classify rules itself",
		"render managed guidance itself",
		"write Profile files directly",
		"mutate repository files itself",
	} {
		if strings.Contains(setup, forbidden) {
			t.Errorf("thin setup skill claims independent behavior %q", forbidden)
		}
	}
}

func runBaselineDigestTarget(t *testing.T, derivedRoot, goCommand string) (baselineDigestTargetResult, string, int) {
	t.Helper()
	command := exec.Command(
		"make",
		"-s",
		"baseline-digests",
		"DERIVED_DIGEST_PATHS="+derivedRoot,
		"BASELINE_DIGEST_STEPS=probe:Probe",
		"GO="+goCommand,
	)
	command.Dir = filepath.Clean(filepath.Join(".."))
	command.Env = append(os.Environ(), "DERIVED_DIGEST_PROBE="+derivedRoot)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	exitCode := 0
	if err := command.Run(); err != nil {
		exitError, ok := err.(*exec.ExitError)
		if !ok {
			t.Fatalf("run baseline-digests: %v", err)
		}
		exitCode = exitError.ExitCode()
	}

	var result baselineDigestTargetResult
	if err := json.Unmarshal(bytes.TrimSpace(stdout.Bytes()), &result); err != nil {
		t.Fatalf("decode baseline-digests stdout %q: %v", stdout.String(), err)
	}
	if result.SchemaVersion == nil || result.Type == nil || result.OK == nil || result.Changed == nil {
		t.Fatalf("baseline-digests result is missing stable fields: %s", stdout.String())
	}
	return result, stderr.String(), exitCode
}

func writeBaselineDigestTargetFile(t *testing.T, path string, data []byte, mode os.FileMode) {
	t.Helper()
	if err := os.WriteFile(path, data, mode); err != nil {
		t.Fatalf("write baseline-digests fixture %s: %v", path, err)
	}
}

func assertBaselineDigestTargetResult(t *testing.T, result baselineDigestTargetResult, wantOK, wantChanged bool) {
	t.Helper()
	if *result.SchemaVersion != 1 {
		t.Errorf("baseline-digests schemaVersion = %d, want 1", *result.SchemaVersion)
	}
	if *result.Type != "baseline-digests" {
		t.Errorf("baseline-digests type = %q, want %q", *result.Type, "baseline-digests")
	}
	if *result.OK != wantOK {
		t.Errorf("baseline-digests ok = %t, want %t", *result.OK, wantOK)
	}
	if *result.Changed != wantChanged {
		t.Errorf("baseline-digests changed = %t, want %t", *result.Changed, wantChanged)
	}
}

func TestAuthorialSkillSync(t *testing.T) {
	t.Parallel()
	repoRoot := filepath.Clean(filepath.Join(".."))
	embeddedFiles, err := Files()
	if err != nil {
		t.Fatalf("read embedded authorial skills: %v", err)
	}
	embeddedByPath := make(map[string][]byte, len(embeddedFiles))
	for _, file := range embeddedFiles {
		embeddedByPath[file.Path] = file.Data
	}

	for _, skillName := range Names() {
		t.Run(skillName, func(t *testing.T) {
			canonicalRoot := filepath.Join(repoRoot, ".agents", "skills", skillName)
			distributedRoot := filepath.Join(repoRoot, "skills", skillName)
			assertSkillTreesEqual(t, canonicalRoot, distributedRoot)

			canonicalDigest, err := SkillFolderHash(t.Context(), canonicalRoot)
			if err != nil {
				t.Fatalf("hash canonical skill: %v", err)
			}
			distributedDigest, err := SkillFolderHash(t.Context(), distributedRoot)
			if err != nil {
				t.Fatalf("hash distributed skill: %v", err)
			}
			if distributedDigest != canonicalDigest {
				t.Fatalf(
					"canonical digest = %q, distributed digest = %q",
					canonicalDigest,
					distributedDigest,
				)
			}

			canonicalFiles := 0
			if err := filepath.WalkDir(
				canonicalRoot,
				func(path string, entry os.DirEntry, walkErr error) error {
					if walkErr != nil {
						return walkErr
					}
					if entry.IsDir() {
						return nil
					}
					relative, err := filepath.Rel(canonicalRoot, path)
					if err != nil {
						return err
					}
					canonical := readBaselineSkillContractFile(t, path)
					embeddedPath := filepath.ToSlash(filepath.Join(skillName, relative))
					if !bytes.Equal(canonical, embeddedByPath[embeddedPath]) {
						t.Errorf("%s canonical and embedded bytes differ", embeddedPath)
					}
					canonicalFiles++
					return nil
				},
			); err != nil {
				t.Fatalf("walk canonical skill: %v", err)
			}
			embeddedCount := 0
			for embeddedPath := range embeddedByPath {
				if strings.HasPrefix(embeddedPath, skillName+"/") {
					embeddedCount++
				}
			}
			if embeddedCount != canonicalFiles {
				t.Fatalf(
					"canonical files = %d, embedded files = %d",
					canonicalFiles,
					embeddedCount,
				)
			}
		})
	}

	setupPaths, err := filepath.Glob(
		filepath.Join(repoRoot, "internal", "baseline", "assets", "setups", "*.json"),
	)
	if err != nil {
		t.Fatalf("list canonical setup snapshots: %v", err)
	}
	if len(setupPaths) == 0 {
		t.Fatal("canonical setup snapshots are missing")
	}
	knownOwned := make(map[string]struct{}, len(Names()))
	for _, skillName := range Names() {
		knownOwned[skillName] = struct{}{}
	}
	for _, setupPath := range setupPaths {
		if *updateDerivedDigests {
			regenerateBaselineSetupSnapshot(t, setupPath, knownOwned)
		}
		t.Run(filepath.Base(setupPath), func(t *testing.T) {
			var setup baselineSetupSnapshot
			if err := json.Unmarshal(readBaselineSkillContractFile(t, setupPath), &setup); err != nil {
				t.Fatalf("decode setup snapshot: %v", err)
			}
			for _, skill := range setup.Skills {
				if skill.Source.Type != "repo" {
					continue
				}
				if skill.Source.Name != "roundfix" {
					t.Errorf("%s repository source = %q, want roundfix", skill.Name, skill.Source.Name)
					continue
				}
				if _, ok := knownOwned[skill.Name]; !ok {
					t.Errorf("%s is repository-sourced but not owned by Roundfix", skill.Name)
					continue
				}
			}
			assertBaselineSetupHasNoOwnedContentPins(t, setupPath, knownOwned)
		})
	}

	lockBytes := readBaselineSkillContractFile(t, filepath.Join(repoRoot, "skills-lock.json"))
	var lock struct {
		Skills map[string]json.RawMessage `json:"skills"`
	}
	if err := json.Unmarshal(lockBytes, &lock); err != nil {
		t.Fatalf("decode skills lock: %v", err)
	}
	if got := upstreamManagedSkillDigest(t, repoRoot, lock.Skills); got != upstreamManagedSkillTreeDigest {
		t.Fatalf("upstream-managed skill tree digest = %q, want %q", got, upstreamManagedSkillTreeDigest)
	}
}

func TestAuthorialSkillSyncUpdateModeRoundTrip(t *testing.T) {
	t.Parallel()
	const minimum = "1.2.3"

	repoRoot := t.TempDir()
	skillRoot := filepath.Join(repoRoot, ".agents", "skills", "roundfix")
	if err := os.MkdirAll(skillRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	skillPath := filepath.Join(skillRoot, "SKILL.md")
	if err := os.WriteFile(skillPath, []byte("before\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	setupPath := filepath.Join(repoRoot, "setup.json")
	setup := baselineSetupSnapshot{
		SchemaVersion: "setup-context-driven/setup-snapshot/0.0.1",
		ID:            "fixture",
		Version:       "0.0.1",
		Source: baselineSetupSource{
			Type:       "github",
			Repository: "example/skills",
			Ref:        strings.Repeat("a", 40),
			Path:       "setups/fixture.txt",
		},
		Digest: "stale",
		Skills: []baselineSetupSkill{{
			Name: "roundfix",
			Path: "skills/roundfix",
			Source: baselineSetupSource{
				Type: "repo",
				Name: "roundfix",
			},
		}},
	}
	writeBaselineSetupFixture(t, setupPath, setup)
	var setupDocument map[string]any
	if err := json.Unmarshal(readBaselineSkillContractFile(t, setupPath), &setupDocument); err != nil {
		t.Fatal(err)
	}
	setupSkills, ok := setupDocument["skills"].([]any)
	if !ok || len(setupSkills) != 1 {
		t.Fatalf("fixture skills = %#v, want one skill", setupDocument["skills"])
	}
	setupSkill, ok := setupSkills[0].(map[string]any)
	if !ok {
		t.Fatalf("fixture skill = %#v, want object", setupSkills[0])
	}
	setupSkill["minimumVersion"] = minimum
	updatedSetup, err := json.MarshalIndent(setupDocument, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(setupPath, append(updatedSetup, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	knownOwned := map[string]struct{}{"roundfix": {}}

	regenerateBaselineSetupSnapshot(t, setupPath, knownOwned)
	before := readBaselineSetupFixture(t, setupPath)
	beforeBytes := readBaselineSkillContractFile(t, setupPath)
	if before.Digest != baselineSetupDigest(t, before.Skills, before.ActivationBundles) {
		t.Fatalf("first regenerated setup has a stale snapshot digest: %+v", before)
	}
	if got := baselineSetupMinimum(t, setupPath); got != minimum {
		t.Fatalf("first regenerated minimum = %q, want preserved %q", got, minimum)
	}

	if err := os.WriteFile(skillPath, []byte("after\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	regenerateBaselineSetupSnapshot(t, setupPath, knownOwned)
	after := readBaselineSetupFixture(t, setupPath)
	if after.Digest != baselineSetupDigest(t, after.Skills, after.ActivationBundles) {
		t.Fatalf("second regenerated setup has a stale snapshot digest: %+v", after)
	}
	if afterBytes := readBaselineSkillContractFile(t, setupPath); !bytes.Equal(afterBytes, beforeBytes) {
		t.Fatal("canonical Skill edit changed the setup snapshot")
	}
	if got := baselineSetupMinimum(t, setupPath); got != minimum {
		t.Fatalf("second regenerated minimum = %q, want preserved %q", got, minimum)
	}
}

func TestCharacterizationCorporaDoNotRecordOwnedSkillDigests(t *testing.T) {
	t.Parallel()

	repoRoot := filepath.Clean(filepath.Join(".."))
	corpusPaths := []string{
		filepath.Join(repoRoot, "internal", "baseline", "testdata", "catalog.diagnostics.golden.json"),
	}
	planPaths, err := filepath.Glob(
		filepath.Join(repoRoot, "internal", "baseline", "testdata", "plan-characterization", "*.golden.json"),
	)
	if err != nil {
		t.Fatalf("list Baseline plan characterization corpus: %v", err)
	}
	if len(planPaths) == 0 {
		t.Fatal("Baseline plan characterization corpus is missing")
	}
	corpusPaths = append(corpusPaths, planPaths...)

	knownOwned := make(map[string]struct{}, len(Names()))
	for _, skillName := range Names() {
		knownOwned[skillName] = struct{}{}
	}
	for _, corpusPath := range corpusPaths {
		var corpus any
		if err := json.Unmarshal(readBaselineSkillContractFile(t, corpusPath), &corpus); err != nil {
			t.Fatalf("decode characterization corpus %s: %v", corpusPath, err)
		}
		for _, violation := range ownedSkillDigestFields(corpus, knownOwned, "$") {
			t.Errorf("%s %s", corpusPath, violation)
		}
	}

	staleDigestFixture := map[string]any{
		"skills": []any{map[string]any{
			"name":          Names()[0],
			"treeDigest":    strings.Repeat("a", 64),
			"contentDigest": strings.Repeat("b", 64),
		}},
	}
	if violations := ownedSkillDigestFields(staleDigestFixture, knownOwned, "$"); len(violations) != 2 {
		t.Fatalf("stale owned-skill digest fixture produced %d violations, want 2: %v", len(violations), violations)
	}
}

func ownedSkillDigestFields(value any, knownOwned map[string]struct{}, jsonPath string) []string {
	var violations []string
	switch value := value.(type) {
	case []any:
		for index, item := range value {
			violations = append(violations, ownedSkillDigestFields(item, knownOwned, fmt.Sprintf("%s[%d]", jsonPath, index))...)
		}
	case map[string]any:
		name, _ := value["name"].(string)
		if _, owned := knownOwned[name]; owned {
			for _, field := range []string{"treeDigest", "contentDigest"} {
				if digest, exists := value[field]; exists && nonEmptyJSONValue(digest) {
					violations = append(violations, fmt.Sprintf("records %s for owned skill %s at %s", field, name, jsonPath))
				}
			}
		}
		for key, item := range value {
			violations = append(violations, ownedSkillDigestFields(item, knownOwned, jsonPath+"."+key)...)
			if key != "content" {
				continue
			}
			encoded, ok := item.(string)
			if !ok {
				continue
			}
			decoded, err := base64.StdEncoding.DecodeString(encoded)
			if err != nil || !json.Valid(decoded) {
				continue
			}
			var embedded any
			if err := json.Unmarshal(decoded, &embedded); err != nil {
				continue
			}
			violations = append(violations, ownedSkillDigestFields(embedded, knownOwned, jsonPath+".content(decoded)")...)
		}
	}
	return violations
}

func nonEmptyJSONValue(value any) bool {
	if value == nil {
		return false
	}
	if text, ok := value.(string); ok {
		return strings.TrimSpace(text) != ""
	}
	return true
}

func baselineSetupMinimum(t *testing.T, setupPath string) string {
	t.Helper()
	var setup struct {
		Skills []struct {
			MinimumVersion string `json:"minimumVersion"`
		} `json:"skills"`
	}
	if err := json.Unmarshal(readBaselineSkillContractFile(t, setupPath), &setup); err != nil {
		t.Fatal(err)
	}
	return setup.Skills[0].MinimumVersion
}

func assertBaselineSetupHasNoOwnedContentPins(
	t *testing.T,
	setupPath string,
	knownOwned map[string]struct{},
) {
	t.Helper()

	var setup struct {
		Skills []map[string]any `json:"skills"`
	}
	if err := json.Unmarshal(readBaselineSkillContractFile(t, setupPath), &setup); err != nil {
		t.Fatal(err)
	}
	for _, skill := range setup.Skills {
		source, ok := skill["source"].(map[string]any)
		if !ok || source["type"] != "repo" {
			continue
		}
		name, _ := skill["name"].(string)
		if _, ok := knownOwned[name]; !ok {
			continue
		}
		for _, field := range []string{"treeDigest", "contentDigest"} {
			if _, exists := skill[field]; exists {
				t.Errorf("%s repository-owned skill %s retains compatibility content pin %s", setupPath, name, field)
			}
		}
	}
}

func TestWritePRDProjectConstraints(t *testing.T) {
	t.Parallel()
	testAuthoringProjectConstraints(t, "write-prd", "prd-template.md")
}

func TestWriteTechSpecProjectConstraints(t *testing.T) {
	t.Parallel()
	testAuthoringProjectConstraints(t, "write-techspec", "techspec-template.md")
}

func TestSpecReferenceLifecycleSkillContracts(t *testing.T) {
	t.Parallel()
	repoRoot := filepath.Clean(filepath.Join(".."))
	tests := []struct {
		name     string
		path     string
		required []string
		ordered  []string
	}{
		{
			name: "archive validation",
			path: filepath.Join(repoRoot, ".agents", "skills", "archive-spec", "SKILL.md"),
			required: []string{
				"`internal/spec.QAVerdict`",
				"expected_owner",
				"duplicate source",
				"type must be `inbox`, `finding`, or `backlog`",
				"exists or is a symbolic link",
				"path must be one basename relative to `_index.md`",
				"test -L \"$index\" || test ! -f \"$index\"",
				"path must equal source basename",
				"test -L \"$current\"",
			},
			ordered: []string{
				"test -L \"$index\" || test ! -f \"$index\"",
				"awk -F '|'",
				"source_basename=${source##*/}",
				"current=\"$(dirname \"$index\")/$path\"",
			},
		},
		{
			name: "PRD adoption",
			path: filepath.Join(repoRoot, ".agents", "skills", "write-prd", "SKILL.md"),
			required: []string{
				"mkdir -p docs/specs/<slug>/references",
				"duplicate source basenames",
				"`_index.md` is a reserved basename",
				"destination path already exists",
				"before changing any finding status",
				"complete `_index.md`",
				"Exclude `_archived/specs/` from automatic link rewrites",
				"Report links from archived Specs separately",
				"only Markdown link destinations",
			},
			ordered: []string{
				"4. **Preflight.**",
				"5. **Index.**",
				"6. **Flip then link.**",
				"7. **Move.**",
				"8. **Rewrite and gate.**",
			},
		},
		{
			name: "idea links",
			path: filepath.Join(repoRoot, ".agents", "skills", "write-idea", "SKILL.md"),
			required: []string{
				"relative to `_idea.md`",
			},
		},
		{
			name: "TechSpec links",
			path: filepath.Join(repoRoot, ".agents", "skills", "write-techspec", "SKILL.md"),
			required: []string{
				"relative to `_techspec.md`",
			},
		},
		{
			name: "Task links",
			path: filepath.Join(repoRoot, ".agents", "skills", "write-tasks", "SKILL.md"),
			required: []string{
				"relative to that Task file",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := strings.Join(strings.Fields(string(readBaselineSkillContractFile(t, tt.path))), " ")
			for _, required := range tt.required {
				if !strings.Contains(content, strings.Join(strings.Fields(required), " ")) {
					t.Errorf("%s missing reference-lifecycle contract %q", tt.path, required)
				}
			}
			previous := -1
			for _, required := range tt.ordered {
				normalized := strings.Join(strings.Fields(required), " ")
				position := strings.Index(content, normalized)
				if position == -1 {
					t.Errorf("%s missing ordered reference-lifecycle contract %q", tt.path, required)
					continue
				}
				if position <= previous {
					t.Errorf("%s reference-lifecycle contract %q is out of order", tt.path, required)
				}
				previous = position
			}
		})
	}
}

func TestProjectConstraintTaskGate(t *testing.T) {
	t.Parallel()
	testWorkflowProjectConstraintContract(t, "write-tasks", []string{
		"roundfix spec check",
		"Project Constraint preflight",
		"active, non-archived, and not already completed",
		"complete `Project Constraints` sections",
		"`Identifier strategy`",
		"`Authentication and HTTP`",
		"`Active ADR obligations`",
		"`Tooling authority`",
		"`docs/agents/` source path",
		"refuse to create or update `_tasks.md`",
		"express maintainer authorization",
		"exact bounded repository-relative files",
		"Dependencies remain owned only by `_tasks.md`",
		"Task status remains owned only by each Task file",
	})
}

func TestProjectConstraintImplementationGate(t *testing.T) {
	t.Parallel()
	testWorkflowProjectConstraintContract(t, "implement-task", []string{
		"Project Constraint preflight",
		"before changing the Task status to `in_progress`",
		"Tooling authorization is not implied by Task assignment",
		"exact bounded repository-relative file list",
		"assigned Task file itself",
		"capture the pre-existing changed paths",
		"before every mutation",
		"changed-file postflight",
		"`git status --short` and `git diff --name-only`",
		"set `status: failed`",
		"Never edit `_tasks.md` or any other Task file",
	})
}

// The authoring skills carry the clauses that produce Project Constraints in the
// first place, and the checker invocation that catches a defect while the
// artifact is still open. Neither may be edited away: Spec 0093's own Task 06
// removed the QA gate's audit block under a standing grant, and only a contract
// test noticed.
// The review-request rule is the one that keeps getting forgotten, and forgetting
// it is how a pull request reaches main unreviewed: automatic review is off by
// configuration, and the rate-limit check passes by design.
func TestReviewRequestContract(t *testing.T) {
	t.Parallel()
	testWorkflowProjectConstraintContract(t, "roundfix", []string{
		"a pull request gets no review unless someone requests one",
		"coderabbit:review",
		"@coderabbitai review",
		"@coderabbitai full review",
		"@coderabbitai rate limit",
		"A green check is not evidence of a review",
		"Review rate limited",
	})
}

func TestProjectConstraintPRDGate(t *testing.T) {
	t.Parallel()
	testWorkflowProjectConstraintContract(t, "write-prd", []string{
		"Project Constraints",
		"express maintainer authorization",
		"bounded files",
		"MUST NOT report",
		"roundfix spec check",
	})
}

func TestProjectConstraintTechSpecGate(t *testing.T) {
	t.Parallel()
	testWorkflowProjectConstraintContract(t, "write-techspec", []string{
		"Project Constraints",
		"express maintainer authorization",
		"bounded files",
		"MUST NOT report",
		"roundfix spec check",
	})
}

// Archiving is where a Spec stops being auditable, so its preconditions are the
// last chance to catch an unfinished one.
func TestArchiveSpecContract(t *testing.T) {
	t.Parallel()
	testWorkflowProjectConstraintContract(t, "archive-spec", []string{
		"self-contained",
		"qa_override",
		"status: completed",
	})
}

func TestProjectConstraintQAGate(t *testing.T) {
	t.Parallel()
	// The gate keeps what a commit is needed to answer and delegates the rest.
	// Clauses that moved are asserted where their rule now lives: the authoring
	// obligations in TestProjectConstraintPRDGate and its TechSpec sibling, the
	// mechanical rules in internal/speccheck. Nothing was dropped; each is
	// required somewhere, and the delegation itself is required here so it
	// cannot be quietly removed.
	testWorkflowProjectConstraintContract(t, "qa-gate", []string{
		"roundfix spec check",
		"Authoring rule removed from the QA matrix",
		"did not run. It is not an equivalent",
		"actual paths from the Daemon-owned Task commit",
		"git diff-tree --no-commit-id --name-only -r",
		"missing, late, or untraceable authorization",
		"outside the exact bounded list",
		"Task status or Task Graph dependencies",
	})
}

func TestProjectConstraintJourney(t *testing.T) {
	t.Parallel()
	fixtureRoot := filepath.Join("testdata", "project-constraint-journey")
	fixtures := []string{"_prd.md", "_techspec.md"}
	for _, fixture := range fixtures {
		t.Run(fixture, func(t *testing.T) {
			content := string(readBaselineSkillContractFile(
				t,
				filepath.Join(fixtureRoot, fixture),
			))
			if missing := missingProjectConstraintFixtureRows(content); len(missing) != 0 {
				t.Fatalf("Project Constraint fixture missing %v", missing)
			}
			for _, row := range []string{
				"Identifier strategy",
				"Authentication and HTTP",
				"Active ADR obligations",
				"Tooling authority",
			} {
				mutated := strings.Replace(content, "- "+row+":", "- Removed row:", 1)
				if missing := missingProjectConstraintFixtureRows(mutated); !slicesContain(missing, row) {
					t.Fatalf("removing %q did not fail the fixture contract: %v", row, missing)
				}
			}
		})
	}

	matrix := string(readBaselineSkillContractFile(
		t,
		filepath.Join(fixtureRoot, "qa-matrix.md"),
	))
	for _, required := range []string{
		"fresh disposable Fluxus greenfield clone",
		"separate fresh disposable Fluxus update clone",
		"Keep-defaults reuses the persisted Better Auth HTTP reason",
		"formatter and repository Verification",
		"audit and empty reapply",
		"final `qa-gate`",
	} {
		if !strings.Contains(matrix, required) {
			t.Errorf("final QA matrix missing %q", required)
		}
	}
	if rows := strings.Count(matrix, "| pending |"); rows != 2 {
		t.Errorf("final QA matrix pending rows = %d, want 2", rows)
	}
}

func TestToolingAuthorizationJourney(t *testing.T) {
	t.Parallel()
	fixtureRoot := filepath.Join("testdata", "project-constraint-journey")
	prd := string(readBaselineSkillContractFile(t, filepath.Join(fixtureRoot, "_prd.md")))
	techspec := string(readBaselineSkillContractFile(t, filepath.Join(fixtureRoot, "_techspec.md")))
	prdPaths := projectConstraintAuthorizedPaths(prd)
	techspecPaths := projectConstraintAuthorizedPaths(techspec)
	wantPaths := []string{".golangci.yml", "scripts/verify.sh"}
	if !reflect.DeepEqual(prdPaths, wantPaths) ||
		!reflect.DeepEqual(techspecPaths, wantPaths) {
		t.Fatalf(
			"bounded tooling files differ: PRD=%v TechSpec=%v want=%v",
			prdPaths,
			techspecPaths,
			wantPaths,
		)
	}

	tests := []struct {
		name    string
		content string
		changed []string
		want    bool
	}{
		{
			name:    "exact bounded files are authorized",
			content: prd,
			changed: wantPaths,
			want:    true,
		},
		{
			name: "missing express authorization refuses decomposition",
			content: strings.Replace(
				prd,
				"expressly authorizes changes to exactly",
				"proposes changes to",
				1,
			),
			changed: wantPaths,
			want:    false,
		},
		{
			name:    "changed path outside bounds refuses settlement",
			content: prd,
			changed: []string{".golangci.yml", "scripts/release.sh"},
			want:    false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := projectConstraintAllowsToolingChanges(test.content, test.changed); got != test.want {
				t.Fatalf("tooling authorization = %v, want %v", got, test.want)
			}
		})
	}

	t.Run("decomposition contract", func(t *testing.T) {
		testWorkflowProjectConstraintContract(t, "write-tasks", []string{
			"refuse to create or update `_tasks.md`",
			"express maintainer authorization",
			"exact bounded repository-relative files",
		})
	})
	t.Run("execution contract", func(t *testing.T) {
		testWorkflowProjectConstraintContract(t, "implement-task", []string{
			"mutation allowlist",
			"before every mutation",
			"Every newly changed path",
			"set `status: failed`",
		})
	})
	t.Run("QA contract", func(t *testing.T) {
		// The gate audits what a commit is needed to answer. The declaration
		// half — express authorization and its exact bounded files — moved to
		// the authoring skills, where `spec check` now decides it before any
		// Agent turn; it is asserted there, not dropped.
		testWorkflowProjectConstraintContract(t, "qa-gate", []string{
			"outside the exact bounded list",
			"actual paths from the Daemon-owned Task commit",
			"missing, late, or untraceable authorization",
			"git diff-tree --no-commit-id --name-only -r",
		})
	})
}

func TestLegacySpecConstraintExemption(t *testing.T) {
	t.Parallel()
	repoRoot := filepath.Clean(filepath.Join(".."))
	for _, skillName := range []string{"write-tasks", "implement-task", "qa-gate"} {
		t.Run(skillName, func(t *testing.T) {
			testWorkflowProjectConstraintContract(t, skillName, []string{
				"completed or archived legacy Spec",
				"byte-identical",
			})
		})
	}

	module := string(readBaselineSkillContractFile(
		t,
		filepath.Join(repoRoot, "internal", "baseline", "assets", "modules", "spec-workflow.json"),
	))
	golden := string(readBaselineSkillContractFile(
		t,
		filepath.Join(
			repoRoot,
			"internal",
			"baseline",
			"assets",
			"formatter-fixtures",
			"standard-typescript-monorepo",
			"golden",
			"docs",
			"agents",
			"spec-routing.md",
		),
	))
	for _, required := range []string{
		"Project Constraints",
		"express maintainer authorization",
		"bounded repository-relative files",
		"actual changed-file scope",
		"completed or archived legacy Specs byte-identical",
		"Dependencies remain owned only by the Task Graph",
		"status remains owned only by each Task file",
	} {
		if !strings.Contains(module, required) {
			t.Errorf("spec-workflow module missing %q", required)
		}
		if !strings.Contains(golden, required) {
			t.Errorf("generated spec-routing fixture missing %q", required)
		}
	}
}

func TestAuthoringConstraintOwnership(t *testing.T) {
	t.Parallel()
	repoRoot := filepath.Clean(filepath.Join(".."))
	lockBytes := readBaselineSkillContractFile(t, filepath.Join(repoRoot, "skills-lock.json"))
	var lock struct {
		Skills map[string]json.RawMessage `json:"skills"`
	}
	if err := json.Unmarshal(lockBytes, &lock); err != nil {
		t.Fatalf("decode skills lock: %v", err)
	}

	for _, name := range []string{"write-prd", "write-techspec"} {
		t.Run(name, func(t *testing.T) {
			if _, managedUpstream := lock.Skills[name]; managedUpstream {
				t.Fatalf("%s must remain repository-owned and absent from skills-lock.json", name)
			}
			assertSkillTreesEqual(
				t,
				filepath.Join(repoRoot, ".agents", "skills", name),
				filepath.Join(repoRoot, "skills", name),
			)
		})
	}

	if got := upstreamManagedSkillDigest(t, repoRoot, lock.Skills); got != upstreamManagedSkillTreeDigest {
		t.Fatalf("upstream-managed skill tree digest = %q, want %q", got, upstreamManagedSkillTreeDigest)
	}
}

func TestUpstreamADRFormatUnchanged(t *testing.T) {
	t.Parallel()
	repoRoot := filepath.Clean(filepath.Join(".."))
	lockBytes := readBaselineSkillContractFile(
		t,
		filepath.Join(repoRoot, "skills-lock.json"),
	)
	var lock struct {
		Version int `json:"version"`
		Skills  map[string]struct {
			ComputedHash string `json:"computedHash"`
		} `json:"skills"`
	}
	if err := json.Unmarshal(lockBytes, &lock); err != nil {
		t.Fatalf("decode skills lock: %v", err)
	}
	if lock.Version != 1 || len(lock.Skills) == 0 {
		t.Fatalf("unexpected skills lock contract: version=%d skills=%d", lock.Version, len(lock.Skills))
	}

	names := make([]string, 0, len(lock.Skills))
	for name := range lock.Skills {
		names = append(names, name)
	}
	sort.Strings(names)
	upstreamDigest := sha256.New()
	for _, name := range names {
		folderDigest, err := SkillFolderHash(t.Context(), filepath.Join(repoRoot, ".agents", "skills", name))
		if err != nil {
			t.Fatalf("hash upstream-managed skill %q: %v", name, err)
		}
		_, _ = upstreamDigest.Write([]byte(name))
		_, _ = upstreamDigest.Write([]byte(folderDigest))
	}
	if got := hex.EncodeToString(upstreamDigest.Sum(nil)); got != upstreamManagedSkillTreeDigest {
		t.Fatalf("upstream-managed skill tree digest = %q, want %q", got, upstreamManagedSkillTreeDigest)
	}

	const wantADRFormatSHA256 = "f1f36cd3f8d3b6474ddd5855da4e233bfc4ae1a1c5024909ccf11871819a41b2"
	adrFormat := readBaselineSkillContractFile(
		t,
		filepath.Join(repoRoot, ".agents", "skills", "domain-modeling", "ADR-FORMAT.md"),
	)
	sum := sha256.Sum256(adrFormat)
	if got := hex.EncodeToString(sum[:]); got != wantADRFormatSHA256 {
		t.Fatalf("domain-modeling/ADR-FORMAT.md digest = %q, want %q", got, wantADRFormatSHA256)
	}
}

func testAuthoringProjectConstraints(t *testing.T, skillName, templateName string) {
	t.Helper()
	repoRoot := filepath.Clean(filepath.Join(".."))
	skill := string(readBaselineSkillContractFile(
		t,
		filepath.Join(repoRoot, ".agents", "skills", skillName, "SKILL.md"),
	))
	template := string(readBaselineSkillContractFile(
		t,
		filepath.Join(repoRoot, ".agents", "skills", skillName, "references", templateName),
	))

	for _, required := range []string{
		"Project Constraints",
		"docs/agents/agent-instructions.md",
		"docs/agents/domain.md",
		"docs/agents/backend.md",
		"docs/agents/spec-routing.md",
		"applicable or not applicable with a reason",
		"express maintainer authorization",
		"bounded files",
		"MUST NOT report completion",
	} {
		if !strings.Contains(skill, required) {
			t.Errorf("%s skill missing %q", skillName, required)
		}
	}

	for _, area := range []string{
		"Identifier strategy",
		"Authentication and HTTP",
		"Active ADR obligations",
		"Tooling authority",
	} {
		if !strings.Contains(template, area) {
			t.Errorf("%s template missing %q", skillName, area)
		}
	}
	if count := strings.Count(template, "<applicable | not applicable>"); count != 4 {
		t.Errorf("%s template applicability placeholders = %d, want 4", skillName, count)
	}
	if count := strings.Count(template, "Source: `docs/agents/"); count != 4 {
		t.Errorf("%s template docs/agents source placeholders = %d, want 4", skillName, count)
	}
	for _, required := range []string{
		"express maintainer authorization",
		"bounded files",
		"no protected tooling mutation",
	} {
		if !strings.Contains(template, required) {
			t.Errorf("%s template missing tooling authorization guidance %q", skillName, required)
		}
	}

	frontmatter := authoringTemplateFrontmatter(t, template)
	if strings.Contains(strings.ToLower(frontmatter), "authoriz") {
		t.Errorf("%s template adds authorization frontmatter: %q", skillName, frontmatter)
	}
}

func testWorkflowProjectConstraintContract(t *testing.T, skillName string, required []string) {
	t.Helper()
	repoRoot := filepath.Clean(filepath.Join(".."))
	canonicalPath := filepath.Join(repoRoot, ".agents", "skills", skillName, "SKILL.md")
	distributedPath := filepath.Join(repoRoot, "skills", skillName, "SKILL.md")
	canonical := readBaselineSkillContractFile(t, canonicalPath)
	distributed := readBaselineSkillContractFile(t, distributedPath)
	if !bytes.Equal(canonical, distributed) {
		t.Fatalf("%s canonical and distributed guidance differ; run make skills-sync", skillName)
	}

	for _, missing := range missingWorkflowContractClauses(string(canonical), required) {
		t.Errorf("%s skill missing %q", skillName, missing)
	}
	normalized := normalizeWorkflowContract(string(canonical))
	for _, clause := range required {
		mutated := strings.ReplaceAll(normalized, normalizeWorkflowContract(clause), "")
		if missing := missingWorkflowContractClauses(mutated, required); !slicesContain(missing, clause) {
			t.Errorf("%s contract still accepts removal of %q; missing = %v", skillName, clause, missing)
		}
	}
}

func missingWorkflowContractClauses(content string, required []string) []string {
	content = normalizeWorkflowContract(content)
	var missing []string
	for _, clause := range required {
		if !strings.Contains(content, normalizeWorkflowContract(clause)) {
			missing = append(missing, clause)
		}
	}
	return missing
}

func normalizeWorkflowContract(content string) string {
	return strings.Join(strings.Fields(content), " ")
}

func slicesContain(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func missingProjectConstraintFixtureRows(content string) []string {
	required := []struct {
		row    string
		source string
	}{
		{row: "Identifier strategy", source: "docs/agents/domain.md"},
		{row: "Authentication and HTTP", source: "docs/agents/backend.md"},
		{row: "Active ADR obligations", source: "docs/agents/domain.md"},
		{row: "Tooling authority", source: "docs/agents/agent-instructions.md"},
	}
	var missing []string
	for _, requirement := range required {
		bullet := projectConstraintFixtureBullet(content, requirement.row)
		if bullet == "" ||
			!strings.Contains(bullet, "applicable") ||
			!strings.Contains(bullet, "Source: `"+requirement.source+"`") {
			missing = append(missing, requirement.row)
		}
	}
	return missing
}

func projectConstraintFixtureBullet(content, row string) string {
	sectionStart := strings.Index(content, "## Project Constraints")
	if sectionStart == -1 {
		return ""
	}
	section := content[sectionStart:]
	if next := strings.Index(section[len("## Project Constraints"):], "\n## "); next != -1 {
		section = section[:len("## Project Constraints")+next]
	}
	marker := "- " + row + ":"
	start := strings.Index(section, marker)
	if start == -1 {
		return ""
	}
	bullet := section[start:]
	if next := strings.Index(bullet[len(marker):], "\n- "); next != -1 {
		bullet = bullet[:len(marker)+next]
	}
	return strings.Join(strings.Fields(bullet), " ")
}

func projectConstraintAuthorizedPaths(content string) []string {
	bullet := projectConstraintFixtureBullet(content, "Tooling authority")
	if !strings.Contains(bullet, "expressly authorizes changes to exactly") {
		return nil
	}
	var paths []string
	for {
		start := strings.Index(bullet, "`")
		if start == -1 {
			break
		}
		bullet = bullet[start+1:]
		end := strings.Index(bullet, "`")
		if end == -1 {
			return nil
		}
		value := bullet[:end]
		bullet = bullet[end+1:]
		if !strings.HasPrefix(value, "docs/agents/") {
			paths = append(paths, value)
		}
	}
	return paths
}

func projectConstraintAllowsToolingChanges(content string, changed []string) bool {
	authorized := projectConstraintAuthorizedPaths(content)
	if len(authorized) == 0 {
		return false
	}
	allowlist := make(map[string]struct{}, len(authorized))
	for _, path := range authorized {
		allowlist[path] = struct{}{}
	}
	for _, path := range changed {
		if _, ok := allowlist[path]; !ok {
			return false
		}
	}
	return true
}

func authoringTemplateFrontmatter(t *testing.T, template string) string {
	t.Helper()
	const fence = "---"
	start := strings.Index(template, fence)
	if start == -1 {
		t.Fatal("template is missing frontmatter")
	}
	rest := template[start+len(fence):]
	end := strings.Index(rest, fence)
	if end == -1 {
		t.Fatal("template has unterminated frontmatter")
	}
	return rest[:end]
}

func assertSkillTreesEqual(t *testing.T, canonicalRoot, distributedRoot string) {
	t.Helper()
	if err := filepath.WalkDir(canonicalRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(canonicalRoot, path)
		if err != nil {
			return err
		}
		canonical := readBaselineSkillContractFile(t, path)
		distributed := readBaselineSkillContractFile(t, filepath.Join(distributedRoot, relative))
		if !bytes.Equal(canonical, distributed) {
			t.Errorf("%s canonical and distributed bytes differ", filepath.ToSlash(relative))
		}
		return nil
	}); err != nil {
		t.Fatalf("walk canonical skill tree %s: %v", canonicalRoot, err)
	}
}

func upstreamManagedSkillDigest(
	t *testing.T,
	repoRoot string,
	managed map[string]json.RawMessage,
) string {
	t.Helper()
	names := make([]string, 0, len(managed))
	for name := range managed {
		names = append(names, name)
	}
	sort.Strings(names)
	digest := sha256.New()
	for _, name := range names {
		folderDigest, err := SkillFolderHash(t.Context(), filepath.Join(repoRoot, ".agents", "skills", name))
		if err != nil {
			t.Fatalf("hash upstream-managed skill %q: %v", name, err)
		}
		_, _ = digest.Write([]byte(name))
		_, _ = digest.Write([]byte(folderDigest))
	}
	return hex.EncodeToString(digest.Sum(nil))
}

type baselineSetupSnapshot struct {
	SchemaVersion     string                `json:"schemaVersion"`
	ID                string                `json:"id"`
	Version           string                `json:"version"`
	Source            baselineSetupSource   `json:"source"`
	Digest            string                `json:"digest"`
	Skills            []baselineSetupSkill  `json:"skills"`
	ActivationBundles []baselineSetupBundle `json:"activationBundles,omitempty"`
}

type baselineSetupSource struct {
	Type       string `json:"type"`
	Repository string `json:"repository,omitempty"`
	Ref        string `json:"ref,omitempty"`
	Path       string `json:"path,omitempty"`
	Name       string `json:"name,omitempty"`
}

type baselineSetupSkill struct {
	Name           string              `json:"name"`
	Path           string              `json:"path"`
	Source         baselineSetupSource `json:"source"`
	MinimumVersion string              `json:"minimumVersion,omitempty"`
	TreeDigest     string              `json:"treeDigest,omitempty"`
}

type baselineSetupBundle struct {
	ID     string   `json:"id"`
	Skills []string `json:"skills"`
}

func regenerateBaselineSetupSnapshot(
	t *testing.T,
	setupPath string,
	knownOwned map[string]struct{},
) {
	t.Helper()
	original := readBaselineSkillContractFile(t, setupPath)
	var setup baselineSetupSnapshot
	if err := json.Unmarshal(original, &setup); err != nil {
		t.Fatalf("decode setup snapshot for regeneration: %v", err)
	}
	for index := range setup.Skills {
		skill := &setup.Skills[index]
		if skill.Source.Type != "repo" {
			continue
		}
		if skill.Source.Name != "roundfix" {
			t.Fatalf(
				"%s repository source = %q, want roundfix",
				skill.Name,
				skill.Source.Name,
			)
		}
		if _, ok := knownOwned[skill.Name]; !ok {
			t.Fatalf("%s is repository-sourced but not owned by Roundfix", skill.Name)
		}
	}
	setup.Digest = baselineSetupDigest(t, setup.Skills, setup.ActivationBundles)
	updated, err := json.MarshalIndent(setup, "", "  ")
	if err != nil {
		t.Fatalf("encode regenerated setup snapshot: %v", err)
	}
	updated = append(updated, '\n')
	if bytes.Equal(original, updated) {
		return
	}
	if err := os.WriteFile(setupPath, updated, 0o644); err != nil {
		t.Fatalf("write regenerated setup snapshot: %v", err)
	}
	t.Logf("regenerated %s", filepath.ToSlash(setupPath))
}

func writeBaselineSetupFixture(t *testing.T, setupPath string, setup baselineSetupSnapshot) {
	t.Helper()
	data, err := json.MarshalIndent(setup, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(setupPath, append(data, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readBaselineSetupFixture(t *testing.T, setupPath string) baselineSetupSnapshot {
	t.Helper()
	var setup baselineSetupSnapshot
	if err := json.Unmarshal(readBaselineSkillContractFile(t, setupPath), &setup); err != nil {
		t.Fatal(err)
	}
	return setup
}

func baselineSetupDigest(
	t *testing.T,
	skills []baselineSetupSkill,
	bundles []baselineSetupBundle,
) string {
	t.Helper()
	var normalizedSkills []map[string]any
	encodedSkills, err := json.Marshal(skills)
	if err != nil {
		t.Fatalf("encode setup skills for digest: %v", err)
	}
	if err := json.Unmarshal(encodedSkills, &normalizedSkills); err != nil {
		t.Fatalf("normalize setup skills for digest: %v", err)
	}
	var payload any = normalizedSkills
	if bundles != nil {
		var normalizedBundles []map[string]any
		encodedBundles, err := json.Marshal(bundles)
		if err != nil {
			t.Fatalf("encode setup bundles for digest: %v", err)
		}
		if err := json.Unmarshal(encodedBundles, &normalizedBundles); err != nil {
			t.Fatalf("normalize setup bundles for digest: %v", err)
		}
		payload = map[string]any{
			"activationBundles": normalizedBundles,
			"skills":            normalizedSkills,
		}
	}
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("serialize canonical setup digest input: %v", err)
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func copyTrackedRepository(t *testing.T, repoRoot string) string {
	t.Helper()

	command := exec.Command("git", append(gittest.ConfigArgs(), "ls-files", "-z", "--cached")...)
	command.Dir = repoRoot
	command.Env = gittest.IsolatedEnv()
	output, err := command.Output()
	if err != nil {
		t.Fatalf("list tracked repository files: %v", err)
	}
	targetRoot := filepath.Join(t.TempDir(), "repository")
	for _, rawRelative := range bytes.Split(output, []byte{0}) {
		if len(rawRelative) == 0 {
			continue
		}
		relative := filepath.FromSlash(string(rawRelative))
		source := filepath.Join(repoRoot, relative)
		target := filepath.Join(targetRoot, relative)
		info, err := os.Lstat(source)
		if err != nil {
			t.Fatalf("inspect tracked file %s: %v", relative, err)
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			t.Fatalf("create tracked file parent %s: %v", relative, err)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			destination, err := os.Readlink(source)
			if err != nil {
				t.Fatalf("read tracked symlink %s: %v", relative, err)
			}
			if err := os.Symlink(destination, target); err != nil {
				t.Fatalf("copy tracked symlink %s: %v", relative, err)
			}
			continue
		}
		data, err := os.ReadFile(source)
		if err != nil {
			t.Fatalf("read tracked file %s: %v", relative, err)
		}
		if err := os.WriteFile(target, data, info.Mode().Perm()); err != nil {
			t.Fatalf("copy tracked file %s: %v", relative, err)
		}
	}
	return targetRoot
}

func tailBytes(data []byte, limit int) []byte {
	if len(data) <= limit {
		return data
	}
	return data[len(data)-limit:]
}

func readBaselineSkillContractFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}

func baselineSkillBody(content string) string {
	const delimiter = "---"
	parts := strings.SplitN(content, delimiter, 3)
	if len(parts) != 3 {
		return content
	}
	return parts[2]
}
