// Suite: Skills commands on a version 2 skills-lock.json and the upgrade next action.
// Invariant: a version 2 lock is read and rewritten in place; a version Roundfix
// still refuses names roundfix upgrade and no skills command.
// Boundary IN: temporary repositories, the bundle install, real Doctor readiness
// and the Baseline update command output and exit codes.
// Boundary OUT: network, the skills CLI, Agent sessions, and every path outside
// a temporary directory.

package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/skillhash"
	"roundfix/skills"
)

const versionTwoStaleHash = "0000000000000000000000000000000000000000000000000000000000000000"

func versionTwoSkillsLock() string {
	return "{\n" +
		"  \"version\": 2,\n" +
		"  \"skills\": {\n" +
		"    \"roundfix\": {\n" +
		"      \"source\": \"marcioaltoe/skills\",\n" +
		"      \"ref\": \"old\",\n" +
		"      \"sourceType\": \"github\",\n" +
		"      \"skillPath\": \"skills/roundfix\",\n" +
		"      \"computedHash\": \"" + versionTwoStaleHash + "\"\n" +
		"    },\n" +
		"    \"testing-boss\": {\n" +
		"      \"source\": \"marcioaltoe/skills\",\n" +
		"      \"sourceType\": \"github\",\n" +
		"      \"skillPath\": \"skills/testing-boss\",\n" +
		"      \"computedHash\": \"" + versionTwoStaleHash + "\"\n" +
		"    }\n" +
		"  }\n" +
		"}\n"
}

// installedSkillHashV2 computes the version 2 digest of an installed skill
// folder exactly as the lock writer does, over the folder's own relative paths.
func installedSkillHashV2(t *testing.T, root string) string {
	t.Helper()
	var files []skillhash.File
	err := filepath.WalkDir(root, func(name string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if name == root {
			return nil
		}
		if entry.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(root, name)
		if err != nil {
			return err
		}
		content, err := os.ReadFile(name)
		if err != nil {
			return err
		}
		files = append(files, skillhash.File{Path: relative, Content: content})
		return nil
	})
	if err != nil {
		t.Fatalf("inspect installed skill %q: %v", root, err)
	}
	digest, err := skillhash.SumForLockVersion(2, files)
	if err != nil {
		t.Fatalf("hash installed skill %q at version 2: %v", root, err)
	}
	return digest
}

func versionTwoDoctorDependencies() doctorDependencies {
	return doctorDependencies{
		resolveExternal: func(string) ([]string, bool, error) { return []string{"testing-boss"}, true, nil },
		checkSkills:     skills.CheckRepositoryWithExternal,
	}
}

func TestTheSkillsCommandsWorkOnAVersionTwoSkillsLock(t *testing.T) {
	t.Parallel()
	_, root := withCLIWorkspace(t)
	if _, err := skills.Install(t.Context(), skills.InstallRequest{Target: "project", ProjectDir: root}); err != nil {
		t.Fatal(err)
	}
	external := filepath.Join(root, ".agents", "skills", "testing-boss")
	mustMkdir(t, external)
	mustWrite(t, filepath.Join(external, "SKILL.md"), "external skill\n")
	lockPath := filepath.Join(root, "skills-lock.json")
	before := versionTwoSkillsLock()
	mustWrite(t, lockPath, before)

	counts := fmt.Sprintf("%d required: %d Roundfix-owned, 1 external", len(skills.Names())+1, len(skills.Names()))
	wantWarn := "skills: warn (" + counts + "; DR-SKILL-LOCK-STALE: skills update would replace these Roundfix-owned skills with the copy their skills-lock.json entry pins: roundfix; next: roundfix skills install --target project)\n"
	wantOK := "skills: ok (" + counts + ")\n"

	dependencies := versionTwoDoctorDependencies()
	var output bytes.Buffer
	printDoctorResult(&output, repositorySkillsCheck(t.Context(), dependencies, root))
	if output.String() != wantWarn {
		t.Fatalf("skills line = %q, want %q", output.String(), wantWarn)
	}

	var stdout, stderr bytes.Buffer
	if code := runCLI(t, []string{"skills", "install", "--target", "project"}, &stdout, &stderr); code != 0 {
		t.Fatalf("skills install exit=%d stderr=%s", code, &stderr)
	}

	hash := installedSkillHashV2(t, filepath.Join(root, ".agents", "skills", "roundfix"))
	wantLock := strings.Replace(before, "      \"ref\": \"old\",\n", "", 1)
	wantLock = strings.Replace(wantLock, versionTwoStaleHash, hash, 1)
	after := string(mustReadBytes(t, lockPath))
	if after != wantLock {
		t.Fatalf("lock=%s\nwant=%s", after, wantLock)
	}
	var decoded struct {
		Version int `json:"version"`
	}
	if err := json.Unmarshal(mustReadBytes(t, lockPath), &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Version != 2 {
		t.Fatalf("lock version = %d, want 2", decoded.Version)
	}

	output.Reset()
	printDoctorResult(&output, repositorySkillsCheck(t.Context(), dependencies, root))
	if output.String() != wantOK {
		t.Fatalf("skills line after install = %q, want %q", output.String(), wantOK)
	}
}

func TestDoctorNamesRoundfixUpgradeForAnUnsupportedSkillsLockVersion(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	lockPath := filepath.Join(root, "skills-lock.json")
	mustWrite(t, lockPath, `{"version":3,"skills":{"testing-boss":{}}}`)
	dependencies := doctorDependencies{
		resolveExternal: func(string) ([]string, bool, error) { return []string{"testing-boss"}, true, nil },
		checkSkills:     skills.CheckRepositoryWithExternal,
	}
	want := "skills: failed (validate skills lock version \"" + lockPath + "\": skills-lock.json version 3 is not supported; this Roundfix reads versions 1 and 2; skills-lock.json comparison unavailable: skills-lock.json version 3 is not supported; this Roundfix reads versions 1 and 2.; next: roundfix upgrade)\n"
	var output bytes.Buffer
	printDoctorResult(&output, repositorySkillsCheck(t.Context(), dependencies, root))
	if output.String() != want {
		t.Fatalf("skills line = %q, want %q", output.String(), want)
	}
}

func TestBaselineUpdateNamesRoundfixUpgradeForAnUnsupportedSkillsLockVersion(t *testing.T) {
	t.Parallel()
	t.Run("version 2", func(t *testing.T) {
		repo := newBaselineUpdateRepository(t)
		installRequiredExternalSkills(t, repo)
		setLockVersion(t, repo, 2)
		before := baselinePlanTestTree(t, repo)
		result, out, errOut, code := runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", repo, "--format=json")
		if code != exitOK || errOut != "" || result.NextAction != "" {
			t.Fatalf("version 2 update exit=%d result=%+v stdout=%s stderr=%s", code, result, out, errOut)
		}
		if before != baselinePlanTestTree(t, repo) {
			t.Fatal("version 2 update changed repository bytes")
		}
	})
	t.Run("version 3", func(t *testing.T) {
		repo := newBaselineUpdateRepository(t)
		installRequiredExternalSkills(t, repo)
		setLockVersion(t, repo, 3)
		result, out, errOut, code := runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", repo, "--format=json")
		if code != exitRunFailed || result.Category != "execution" ||
			result.NextAction != "run roundfix upgrade, then rerun roundfix baseline update" ||
			!strings.Contains(errOut, "skills-lock.json version 3 is not supported; this Roundfix reads versions 1 and 2.") {
			t.Fatalf("version 3 update exit=%d result=%+v stdout=%s stderr=%s", code, result, out, errOut)
		}
	})
}

// setLockVersion rewrites only the lock's version number, keeping every entry.
func setLockVersion(t *testing.T, repo string, version int) {
	t.Helper()
	data := mustReadBytes(t, filepath.Join(repo, "skills-lock.json"))
	var lock map[string]any
	if err := json.Unmarshal(data, &lock); err != nil {
		t.Fatal(err)
	}
	lock["version"] = version
	encoded, err := json.Marshal(lock)
	if err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(repo, "skills-lock.json"), string(encoded))
	commitBaselinePlanTestRepository(t, repo)
}
