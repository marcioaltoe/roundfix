// Suite: Baseline release-gate skill cutover
// Invariant: the shipped setup skill is parser-aligned guidance and cannot contain an executable Baseline engine.
// Boundary IN: canonical, distributed, and embedded setup skills plus the executable-artifact rejection guard.
// Boundary OUT: public CLI journeys and repository adoption.

package skills

import "testing"

func TestBaselineDocumentationContractThinSkill(t *testing.T) {
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{name: "public CLI recipes", run: TestBaselineSkillContract},
		{name: "no Python fallback", run: TestNoPythonBaselineRuntime},
		{name: "guidance-only skill trees", run: TestThinSetupSkill},
		{name: "executable engine rejection", run: TestCheckRejectsExecutableSetupEngineArtifacts},
	}
	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}
