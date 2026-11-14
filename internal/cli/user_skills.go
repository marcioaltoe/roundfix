package cli

import (
	"context"
	"fmt"
	"io"
	"strings"

	roundskills "roundfix/skills"
)

func reportedUserSkillCopies(ctx context.Context, homeDir string, declaredSkills []string) []roundskills.StaleUserSkillCopy {
	directories := commandDependenciesForContext(ctx).userSkillDirectories
	if directories == nil || strings.TrimSpace(homeDir) == "" {
		return nil
	}
	return roundskills.ReportedUserSkillCopies(ctx, directories(homeDir), declaredSkills)
}

func printUserSkillWarning(stderr io.Writer, copies []roundskills.StaleUserSkillCopy) {
	if len(copies) == 0 {
		return
	}
	fmt.Fprintln(stderr, "roundfix: warning: user-level Roundfix skill copies that an Agent harness may load in place of the project copy:")
	for _, copy := range copies {
		fmt.Fprintln(stderr, "  "+copy.Describe())
	}
	for _, removal := range roundskills.UserSkillRemovals(copies) {
		fmt.Fprintln(stderr, "  remove: "+removal)
	}
	for _, refresh := range roundskills.StaleUserSkillRefreshes(copies) {
		fmt.Fprintln(stderr, "  refresh: "+refresh)
	}
}

func doctorUserSkillsResult(copies []roundskills.StaleUserSkillCopy) CheckResult {
	detail := "user-level Roundfix skill copies that an Agent harness may load in place of the project copy: " + describeUserSkillCopies(copies)
	if refreshes := roundskills.StaleUserSkillRefreshes(copies); len(refreshes) > 0 {
		detail += "; or refresh with " + strings.Join(refreshes, " && ")
	}
	return CheckResult{Name: HealthCheckUserSkills, Status: CheckStatusWarn,
		Detail: detail, NextAction: strings.Join(roundskills.UserSkillRemovals(copies), " && ")}
}

func printUpdateCheckUserSkills(stdout io.Writer, copies []roundskills.StaleUserSkillCopy) {
	if len(copies) == 0 {
		return
	}
	fmt.Fprintf(stdout, "User skills: copies that an Agent harness may load in place of the project copy: %s; remove with %s", describeUserSkillCopies(copies), strings.Join(roundskills.UserSkillRemovals(copies), " && "))
	if refreshes := roundskills.StaleUserSkillRefreshes(copies); len(refreshes) > 0 {
		fmt.Fprintf(stdout, ", or refresh with %s", strings.Join(refreshes, " && "))
	}
	fmt.Fprintln(stdout, ".")
}

func describeUserSkillCopies(copies []roundskills.StaleUserSkillCopy) string {
	descriptions := make([]string, 0, len(copies))
	for _, copy := range copies {
		descriptions = append(descriptions, copy.Describe())
	}
	return strings.Join(descriptions, ", ")
}
