package cli

import (
	"context"

	roundskills "roundfix/skills"
)

// repositorySkillDiagnostics reports only named missing required skills.
// Resolution failures leave the bundle-only check unchanged.
func repositorySkillDiagnostics(ctx context.Context, environment commandEnvironment) []roundskills.Diagnostic {
	workDir, err := environment.resolveWorkDir("resolve skills check working directory")
	if err != nil {
		return nil
	}
	root, err := environment.dependencies.resolveSkillsProjectRoot(ctx, workDir)
	if err != nil {
		return nil
	}
	external, manifestOK, err := environment.dependencies.doctor.resolveExternal(root)
	if err != nil || !manifestOK || len(external) == 0 {
		return nil
	}
	profile, err := requiredSkillsManifestProfile(root)
	if err != nil {
		return nil
	}
	readiness, checkErr := environment.dependencies.doctor.checkSkills(ctx, root, external)
	var diagnostics []roundskills.Diagnostic
	for _, missing := range missingRequiredSkills(profile, readiness, checkErr) {
		diagnostics = append(diagnostics, roundskills.Diagnostic{Path: ".agents/skills/" + missing.Skill, Message: missing.Reason})
	}
	return diagnostics
}
