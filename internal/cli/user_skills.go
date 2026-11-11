package cli

import (
	"context"
	"fmt"
	"io"
	"strings"

	roundskills "roundfix/skills"
)

func staleUserSkillCopies(ctx context.Context, homeDir string) []roundskills.StaleUserSkillCopy {
	directories := commandDependenciesForContext(ctx).userSkillDirectories
	if directories == nil || strings.TrimSpace(homeDir) == "" {
		return nil
	}
	return roundskills.StaleUserSkillCopies(ctx, directories(homeDir))
}

func printUserSkillWarning(stderr io.Writer, copies []roundskills.StaleUserSkillCopy) {
	if len(copies) == 0 {
		return
	}
	fmt.Fprintln(stderr, "roundfix: warning: user-level Roundfix skill copies older than this binary; an Agent harness may load them in place of the project copy:")
	for _, copy := range copies {
		fmt.Fprintln(stderr, "  "+copy.Describe())
	}
	for _, refresh := range roundskills.StaleUserSkillRefreshes(copies) {
		fmt.Fprintln(stderr, "  refresh: "+refresh)
	}
}

func doctorUserSkillsResult(copies []roundskills.StaleUserSkillCopy) CheckResult {
	return CheckResult{Name: HealthCheckUserSkills, Status: CheckStatusWarn,
		Detail:     "user-level Roundfix skill copies older than this binary: " + describeUserSkillCopies(copies),
		NextAction: strings.Join(roundskills.StaleUserSkillRefreshes(copies), " && ")}
}

func printUpdateCheckUserSkills(stdout io.Writer, copies []roundskills.StaleUserSkillCopy) {
	if len(copies) == 0 {
		return
	}
	fmt.Fprintf(stdout, "User skills: copies older than this binary: %s; refresh with %s.\n", describeUserSkillCopies(copies), strings.Join(roundskills.StaleUserSkillRefreshes(copies), " && "))
}

func describeUserSkillCopies(copies []roundskills.StaleUserSkillCopy) string {
	descriptions := make([]string, 0, len(copies))
	for _, copy := range copies {
		descriptions = append(descriptions, copy.Describe())
	}
	return strings.Join(descriptions, ", ")
}
