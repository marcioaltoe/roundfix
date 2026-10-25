// Suite: repository-profile Baseline commands
// Invariant: repository profiles use their modules' embedded skill contracts.
// Boundary IN: real profile adoption, update preview, owned install, Doctor skills check, and restore refusal streams.
// Boundary OUT: source acquisition during refresh is replaced by a restore fake; machine Doctor checks are not exercised.
package cli

import (
	"bytes"
	"path/filepath"
	"testing"

	"roundfix/internal/baseline"
)

const baselineRepositoryProfileTestID = "repository-go-cli-tui"

func TestBaselineSkillsRestoreNamesTheProfilePath(t *testing.T) {
	t.Parallel()
	root := newBaselineRepositoryProfileFixture(t, false)
	before := baselinePlanTestTree(t, root)
	out, stderr, code := runRepositoryProfileCommand(t, "baseline", "skills", "restore", "--profile", "missing-profile", "--repo", root, "--format", "text")
	const message = `Baseline Profile "missing-profile" is neither built-in nor resolvable at .roundfix/baseline/profiles/missing-profile.json.`
	wantOut := "Baseline skills restore: blocked\nrestore.profile-unresolved: " + message + "\nNext action: Run roundfix baseline profile validate missing-profile to see why, restore that file, or choose a built-in profile id.\n"
	wantErr := "roundfix: baseline skills restore failed: " + message + "\n"
	if code != exitPreflight || out != wantOut || stderr != wantErr {
		t.Fatalf("Surface Transcript 2: exit=%d stdout=%q stderr=%q", code, out, stderr)
	}
	if before != baselinePlanTestTree(t, root) {
		t.Fatal("blocked restore changed repository bytes")
	}
}

func TestBaselineSkillsRestoreAcceptsARepositoryProfile(t *testing.T) {
	t.Parallel()
	root := newBaselineRepositoryProfileFixture(t, false)
	source := filepath.Join(root, "missing-source")
	before := baselinePlanTestTree(t, root)
	out, stderr, code := runRepositoryProfileCommand(t, "baseline", "skills", "restore", "--profile", baselineRepositoryProfileTestID,
		"--skill", "context7-cli", "--source-dir", source, "--repo", root, "--format", "text")
	message := "Offline Git object store is not a directory: " + source + "."
	wantOut := "Baseline skills restore: blocked\nrestore.source-dir-invalid: " + message + "\nNext action: Pass an existing Git checkout or bare object store to --source-dir.\n"
	wantErr := "roundfix: baseline skills restore failed: " + message + "\n"
	if code != exitPreflight || out != wantOut || stderr != wantErr {
		t.Fatalf("Surface Transcript 3: exit=%d stdout=%q stderr=%q", code, out, stderr)
	}
	if before != baselinePlanTestTree(t, root) {
		t.Fatal("blocked restore changed repository bytes")
	}
}

func newBaselineRepositoryProfileFixture(t *testing.T, copySkills bool) string {
	t.Helper()
	source, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	root := newBaselineApplyTestRepository(t)
	setCommandWorkDirForTest(t, root)
	out, stderr, code := runRepositoryProfileCommand(t, "baseline", "profile", "init", "--id", baselineRepositoryProfileTestID, "--from", "go-cli-tui")
	if code != exitOK || stderr != "" {
		t.Fatalf("profile init exit=%d stdout=%s stderr=%s", code, out, stderr)
	}
	commitBaselinePlanTestRepository(t, root)
	args := []string{"baseline", "plan", "--repo", root, "--profile", baselineRepositoryProfileTestID,
		"--decision", "preservation.mode=greenfield",
		"--decision", "language.generated=English",
		"--decision", "verification.gate=make verify",
		"--decision", "verification.incremental=make verify-incremental",
		"--decision", "branch.prefix=ma/",
		"--decision", "spec.scaffold=true",
		"--decision", "domain.layout=single-context",
		"--decision", "triage.external=false",
		"--decision", "autonomous.enabled=true",
		"--decision", "runtime.backend=codex gpt-5.5 xhigh",
		"--decision", "runtime.design=claude opus xhigh",
		"--decision", "secondbrain.enabled=false",
		"--decision", "repository.extension.enabled=false", "--format", "json"}
	out, stderr, code = runRepositoryProfileCommand(t, args...)
	if code != exitOK || stderr != "" {
		t.Fatalf("plan exit=%d stdout=%s stderr=%s", code, out, stderr)
	}
	plan, err := baseline.ParsePlanDocument([]byte(out))
	if err != nil {
		t.Fatal(err)
	}
	planPath := filepath.Join(t.TempDir(), "plan.json")
	mustWrite(t, planPath, out)
	out, stderr, code = runRepositoryProfileCommand(t, "baseline", "apply", "--repo", root, "--plan", planPath, "--confirm-plan", plan.PlanDigest, "--format", "json")
	if code != exitOK || stderr != "" {
		t.Fatalf("apply exit=%d stdout=%s stderr=%s", code, out, stderr)
	}
	if copySkills {
		for _, relative := range []string{".agents/skills", "skills-lock.json"} {
			copyThisRepositorySkillSetPath(t, filepath.Join(source, relative), filepath.Join(root, relative))
		}
	}
	commitBaselinePlanTestRepository(t, root)
	return root
}

func runRepositoryProfileCommand(t *testing.T, args ...string) (string, string, int) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := runCLIContext(t, t.Context(), args, &stdout, &stderr)
	return stdout.String(), stderr.String(), code
}
