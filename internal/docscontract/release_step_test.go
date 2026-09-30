//go:build docscontract

package docscontract

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

const releaseSkillsAndGuidesHeading = "## Checking skills and guides before the release"

func TestTheReleaseRunbookRequiresTheSkillsAndGuidesCheck(t *testing.T) {
	runbook := readReleaseRunbook(t)
	section := releaseSkillsAndGuidesSection(runbook)
	normalizedSection := strings.Join(strings.Fields(section), " ")
	for _, want := range []string{
		"This is a mandatory release step.",
		"before opening the release Pull Request",
		"roundfix doctor",
		"TestEveryOwnedSkillVersionIsRecorded",
		"TestTheOwnedSkillMinimumIsTheEmbeddedVersion",
		"TestNoTwoBaselineClausesShareText",
		"TestBaselineClauseForceIsCharacterized",
		"TestShippedGuidanceCitesNoRepositoryRecord",
		"TestTheLoopClauseOrderMatchesTheDeliveryQueue",
		"TestEveryCommandIsNamedInTheRoundfixSkill",
		"TestEveryCommandIsNamedInTheUserGuide",
		"reading pass",
	} {
		if !strings.Contains(normalizedSection, want) {
			t.Fatalf("release check section is missing %q", want)
		}
	}

	cutting := strings.Index(runbook, "## Cutting a release")
	plan := strings.Index(runbook[cutting:], "Run `roundfix release plan`")
	check := strings.Index(runbook[cutting:], "skills and guides check")
	tag := strings.Index(runbook[cutting:], "git tag v<version>")
	if cutting < 0 || plan < 0 || check < 0 || tag < 0 || !(plan < check && check < tag) {
		t.Fatalf("cutting a release does not point to the check between planning and tagging")
	}
}

func TestEveryCheckTheReleaseStepNamesExists(t *testing.T) {
	runbook := readReleaseRunbook(t)
	checks := releaseStepCheckNames(releaseSkillsAndGuidesSection(runbook))
	if len(checks) < 7 {
		t.Fatalf("release step names %d tests, want at least 7", len(checks))
	}
	available := repositoryTestNames(t)
	if missing := missingReleaseStepChecks(checks, available); len(missing) != 0 {
		t.Fatalf("release step names missing repository tests: %v", missing)
	}
}

func TestAReleaseStepThatNamesAMissingCheckIsReported(t *testing.T) {
	checks := releaseStepCheckNames(releaseSkillsAndGuidesSection(releaseSkillsAndGuidesHeading + "\n\n1. `TestAReleaseCheckThatDoesNotExist`."))
	want := []string{"TestAReleaseCheckThatDoesNotExist"}
	if got := missingReleaseStepChecks(checks, map[string]bool{}); strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("missing checks = %v, want %v", got, want)
	}
}

func readReleaseRunbook(t *testing.T) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join("..", "..", "docs", "user-guide", "release-runbook.md"))
	if err != nil {
		t.Fatalf("read release runbook: %v", err)
	}
	return string(content)
}

func releaseSkillsAndGuidesSection(runbook string) string {
	start := strings.Index(runbook, releaseSkillsAndGuidesHeading)
	if start < 0 {
		return ""
	}
	end := strings.Index(runbook[start+len(releaseSkillsAndGuidesHeading):], "\n## ")
	if end < 0 {
		return runbook[start:]
	}
	return runbook[start : start+len(releaseSkillsAndGuidesHeading)+end]
}

func releaseStepCheckNames(section string) []string {
	return regexp.MustCompile(`\bTest[A-Z][A-Za-z0-9_]*\b`).FindAllString(section, -1)
}

func missingReleaseStepChecks(checks []string, available map[string]bool) []string {
	seen := make(map[string]bool)
	var missing []string
	for _, check := range checks {
		if !available[check] && !seen[check] {
			missing = append(missing, check)
			seen[check] = true
		}
	}
	return missing
}

func repositoryTestNames(t *testing.T) map[string]bool {
	t.Helper()
	root := filepath.Join("..", "..")
	decl := regexp.MustCompile(`^func (Test[A-Z][A-Za-z0-9_]*)\(`)
	names := make(map[string]bool)
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() && (entry.Name() == ".git" || entry.Name() == "vendor") {
			return filepath.SkipDir
		}
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), "_test.go") {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for _, line := range strings.Split(string(content), "\n") {
			if match := decl.FindStringSubmatch(line); len(match) == 2 {
				names[match[1]] = true
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk repository tests: %v", err)
	}
	return names
}
