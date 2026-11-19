package cli

import (
	"context"
	"strings"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/preflight"
)

type followUpAuthoringQuestions struct {
	Preference string `json:"preference"`
	Source     string `json:"source"`
	Ref        string `json:"ref"`
	Checkout   string `json:"checkout"`
	Detail     string `json:"detail"`
}

// readFollowUpAuthoringQuestions reads only local default-branch objects. An
// unavailable head or config leaves the rest of the status report usable.
func readFollowUpAuthoringQuestions(ctx context.Context, loaded roundconfig.Loaded, runner preflight.GitRunner) followUpAuthoringQuestions {
	report := followUpAuthoringQuestions{Preference: roundconfig.OpenQuestionsWait, Source: "default"}
	branch, err := runner.RunGit(ctx, loaded.GitRoot, "symbolic-ref", "--quiet", "--short", "HEAD")
	if err != nil {
		branch = "" // Detached HEAD can still resolve the default branch.
	}
	head, found, err := preflight.ResolveDefaultBranchHead(ctx, loaded.GitRoot, strings.TrimSpace(branch), runner)
	if err != nil || !found {
		report.Detail = "no default branch head was found"
	} else {
		report.Ref = head.Ref
		config, err := roundconfig.DeliveryConfigAtCommit(ctx, runner, loaded.GitRoot, loaded.UserConfigPath, head.Commit)
		if err != nil {
			report.Detail = "the default branch's Project Config was not read: " + err.Error()
		} else {
			report.Preference = config.Authoring.OpenQuestions
			if config.Authoring.OpenQuestionsDeclared {
				report.Source = "project-config"
			}
		}
	}
	if checkout := loaded.Config.Authoring.OpenQuestions; checkout != report.Preference {
		report.Checkout = checkout
	}
	return report
}
