package cli

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"roundfix/internal/baseline"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/gittest"
	"roundfix/skills"
)

func TestDoctorKeepsAnExternalSkillThatDiffersFromItsLockReady(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	gittest.MarkRepository(t, root)
	writeDoctorReadyRepositoryFixture(t, root)
	setDoctorTestManifestProfile(t, root, "go-cli-tui")
	mustWrite(t, filepath.Join(root, ".agents", "skills", "testing-boss", "SKILL.md"), "changed after recording its lock entry\n")
	checker := newDoctorFakeHealthChecker(CheckResult{Name: HealthCheckNode, Status: CheckStatusOK}, CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK}, CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK})
	withDoctorFakeLoadedAndReadiness(t, checker, roundconfig.Loaded{Config: roundconfig.Builtin(), GitRoot: root}, func(context.Context, roundconfig.Config, []roundconfig.WorkCategory, string) profileProofResult {
		return profileProofResult{}
	})
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		fixture := dependencies.doctor
		dependencies.doctor = defaultDoctorDependencies()
		dependencies.doctor.loadConfig = fixture.loadConfig
		dependencies.doctor.readiness = fixture.readiness
		dependencies.doctor.healthChecker = fixture.healthChecker
		dependencies.doctor.profileReadiness = fixture.profileReadiness
		dependencies.doctor.residue = fixture.residue
		dependencies.doctor.storage = fixture.storage
		dependencies.doctor.diskSpace = nil
		dependencies.doctor.linuxLane = nil
		dependencies.doctor.resolveExternal = func(string) ([]string, bool, error) { return []string{"testing-boss"}, true, nil }
		dependencies.doctor.checkSkills = skills.CheckRepositoryWithExternal
	})
	before := snapshotDoctorPath(t, root)
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"doctor"}, &stdout, &stderr)
	owned := len(skills.Names())
	want := fmt.Sprintf("skills: ok (%d required: %d Roundfix-owned, 1 external)", owned+1, owned)
	if code != exitOK || stderr.Len() != 0 || !slices.Contains(strings.Split(stdout.String(), "\n"), want) || strings.Contains(stdout.String(), "outdated:") {
		t.Fatalf("exit=%d stdout=%s stderr=%s; want exit=0 line=%s", code, stdout.String(), stderr.String(), want)
	}
	if !reflect.DeepEqual(before, snapshotDoctorPath(t, root)) {
		t.Fatal("Doctor changed repository bytes")
	}
}

func TestBaselineUpdateNeverRestoresAnExternalSkillThatDiffersFromItsLock(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	gittest.MarkRepository(t, root)
	writeDoctorReadyRepositoryFixture(t, root)
	setDoctorTestManifestProfile(t, root, "go-cli-tui")
	path := filepath.Join(root, ".agents", "skills", "testing-boss", "SKILL.md")
	mustWrite(t, path, "changed after recording its lock entry\n")
	before := mustReadBytes(t, path)
	result, err := runBaselineUpdateSkillsStageWith(t.Context(), baselineUpdateSkillsRequest{
		Repository: root, ProfileID: "go-cli-tui",
	}, baselineUpdateSkillsDependencies{
		recordLock:         baseline.RecordInstalledOwnedSkills,
		resolveProjectRoot: defaultResolveSkillsProjectRoot,
		install:            skills.Install,
		ownedNames:         skills.Names,
		resolveExternal:    func(string) ([]string, bool, error) { return []string{"testing-boss"}, true, nil },
		checkRepository:    skills.CheckRepositoryWithExternal,
		restore: func(context.Context, baseline.SkillsRestoreRequest) (baseline.SkillsRestorePayload, error) {
			t.Fatal("restore called for a present external skill")
			return baseline.SkillsRestorePayload{}, nil
		},
	})
	if err != nil || result.Status != baselineUpdateSkillsVerified || len(result.Drifted) != 0 || len(result.Restored) != 0 {
		t.Fatalf("result=%+v err=%v", result, err)
	}
	if !bytes.Equal(before, mustReadBytes(t, path)) {
		t.Fatal("refresh changed the present skill's bytes")
	}
}
