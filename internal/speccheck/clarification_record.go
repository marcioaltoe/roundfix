package speccheck

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

const (
	// CodeDecisionsUnrecorded identifies a PRD without a Clarification Record.
	CodeDecisionsUnrecorded = "SC-DECISIONS-UNRECORDED"
	// ClarificationGuidePath is the guide whose oldest adding commit starts the horizon.
	ClarificationGuidePath = ".agents/skills/write-prd/references/clarification-record.md"
)

var clarificationDatedAnswer = regexp.MustCompile(`^- (\d{4}-\d{2}-\d{2}): \S`)
var clarificationAutonomousRoute = regexp.MustCompile(`^Autonomous route: no person present; decisions taken from \S.*\.$`)

func clarificationHorizon(repoRoot, prdPath string) contractHorizon {
	missing := "a PRD committed at or after " + ClarificationGuidePath
	if _, ok, err := receiptFilePath(repoRoot, ClarificationGuidePath); err != nil {
		return contractHorizon{held: true}
	} else if !ok {
		return contractHorizon{missing: missing}
	}
	commit, readable := prdAddingCommit(repoRoot, prdPath)
	if !readable {
		return contractHorizon{held: true}
	}
	output, err := adrHorizonGitOutput(repoRoot, "log", "--diff-filter=A", "--format=%H", "--", ClarificationGuidePath)
	if err != nil {
		return contractHorizon{held: true}
	}
	commits := strings.Fields(string(output))
	if len(commits) == 0 {
		return contractHorizon{missing: missing}
	}
	if commit == "" {
		return contractHorizon{held: true}
	}
	_, err = adrHorizonGitOutput(repoRoot, "merge-base", "--is-ancestor", commits[len(commits)-1], commit)
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
		return contractHorizon{missing: missing}
	}
	return contractHorizon{held: true}
}

func detectClarificationRecord(result *Result, repoRoot, prdPath string) error {
	content, err := os.ReadFile(prdPath)
	display := artifactDisplayPath(repoRoot, prdPath)
	if errors.Is(err, os.ErrNotExist) {
		addSkip(result, CodeDecisionsUnrecorded, display)
		return nil
	}
	if err != nil {
		return fmt.Errorf("read Clarification Record %q: %w", prdPath, err)
	}
	horizon := clarificationHorizon(repoRoot, prdPath)
	if !horizon.held {
		addSkip(result, CodeDecisionsUnrecorded, horizon.missing)
		return nil
	}
	inSection := false
	heading := 1
	for index, line := range strings.Split(string(content), "\n") {
		if strings.HasPrefix(line, "## ") {
			inSection = line == "## Decisions"
			if inSection {
				heading = index + 1
			}
			continue
		}
		if !inSection {
			continue
		}
		if match := clarificationDatedAnswer.FindStringSubmatch(line); match != nil {
			if _, err := time.Parse(time.DateOnly, match[1]); err == nil {
				return nil
			}
		}
		if clarificationAutonomousRoute.MatchString(line) {
			return nil
		}
	}
	result.Findings = append(result.Findings, Finding{
		Code: CodeDecisionsUnrecorded, Severity: SeverityGap,
		Summary: display + " has no Clarification Record",
		Where:   []Location{{Path: display, Line: heading}},
		Fix:     "Under ## Decisions, write one - YYYY-MM-DD: <question> — <answer> line for each question a person answered, or the line Autonomous route: no person present; decisions taken from <source>., following the write-prd rule \"Clarification Record\".",
	})
	return nil
}
