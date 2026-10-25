// Suite: Missing Required Skills in skills check.
// Invariant: adopted repositories name missing required skills; other directories keep bundle output.
// Boundary IN: temporary repositories, real skill inspection, CLI output and exit codes.
// Boundary OUT: network and host-installed skills.

package cli

import (
	"bytes"
	"strings"
	"testing"

	"roundfix/skills"
)

func TestSkillsCheckNamesAMissingRequiredSkill(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"adopted repository with a missing required skill", "adopted repository with every required skill", "repository without a Setup Manifest", "outside a Git repository"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var repository, missing string
			switch name {
			case "adopted repository with a missing required skill", "adopted repository with every required skill":
				repository = newBaselineUpdateRepository(t)
				required := installRequiredExternalSkills(t, repository)
				if name == "adopted repository with a missing required skill" {
					missing = required[0]
					removeInstalledSkill(t, repository, missing, true)
				}
			case "repository without a Setup Manifest":
				repository = newBaselinePlanTestRepository(t)
			default:
				repository = t.TempDir()
			}
			setCommandWorkDirForTest(t, repository)
			var stdout, stderr bytes.Buffer
			code := runCLI(t, []string{"skills", "check"}, &stdout, &stderr)
			if missing != "" {
				want := "roundfix: Roundfix skill check failed:\n  .agents/skills/" + missing + ": required by the Baseline and not installed; restore with roundfix baseline skills restore --profile go-cli-tui --skill " + missing + "\n"
				if code != exitRunFailed || stdout.Len() != 0 || stderr.String() != want {
					t.Fatalf("exit=%d stdout=%q stderr=%q; want exit 1, empty stdout, stderr=%q", code, stdout.String(), stderr.String(), want)
				}
			} else {
				want := "Roundfix skill check passed: " + strings.Join(skills.Names(), ", ") + "\n"
				if code != exitOK || stdout.String() != want || stderr.Len() != 0 {
					t.Fatalf("exit=%d stdout=%q stderr=%q; want exit 0, stdout=%q, empty stderr", code, stdout.String(), stderr.String(), want)
				}
			}
		})
	}
}
