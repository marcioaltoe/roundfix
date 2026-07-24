package cli

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"roundfix/internal/app"
	"roundfix/internal/baseline"
)

type baselineSkillsRestoreCommandRequest struct {
	repo         string
	profile      string
	skills       stringListFlag
	sourceDir    string
	confirmation string
	format       string
}

type baselineSkillsRestoreExecutor func(
	context.Context,
	baseline.SkillsRestoreRequest,
) (baseline.SkillsRestorePayload, error)

func runBaselineSkillsRestoreCommand(
	ctx context.Context,
	args []string,
	stdout, stderr io.Writer,
) int {
	return runBaselineSkillsRestoreCommandWith(
		ctx,
		args,
		stdout,
		stderr,
		baseline.RestoreSkills,
	)
}

func runBaselineSkillsRestoreCommandWith(
	ctx context.Context,
	args []string,
	stdout, stderr io.Writer,
	execute baselineSkillsRestoreExecutor,
) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("baseline skills restore"))
		return exitOK
	}
	request, err := parseBaselineSkillsRestoreCommand(args)
	jsonOutput := request.format == "json" || baselineSkillsRestoreJSONRequested(args)
	if err != nil {
		payload := baseline.SkillsRestorePayload{
			SchemaVersion:  baseline.SkillsRestoreSchemaVersion,
			Profile:        request.profile,
			Acquisitions:   []baseline.RestoreAcquisition{},
			Skills:         []baseline.RestoreSkill{},
			PlannedChanges: []baseline.RestorePlannedChange{},
			Finding: &baseline.RestoreFinding{
				Code:    "restore.arguments-invalid",
				Message: err.Error(),
				Action:  "Correct the command input and rerun roundfix baseline skills restore.",
			},
		}
		return writeBaselineSkillsRestoreFailure(
			payload,
			err,
			exitPreflight,
			jsonOutput,
			stdout,
			stderr,
		)
	}
	payload, restoreErr := execute(ctx, baseline.SkillsRestoreRequest{
		Repository:   request.repo,
		ProfileID:    request.profile,
		Skills:       append([]string(nil), request.skills...),
		SourceDir:    request.sourceDir,
		Confirmation: request.confirmation,
	})
	if restoreErr != nil {
		exit := baselineSkillsRestoreExit(restoreErr)
		return writeBaselineSkillsRestoreFailure(
			payload,
			restoreErr,
			exit,
			jsonOutput,
			stdout,
			stderr,
		)
	}
	if err := writeBaselineSkillsRestorePayload(payload, jsonOutput, stdout); err != nil {
		fmt.Fprintf(stderr, "%s: baseline skills restore output failed: %v\n", app.Name, err)
		return exitRunFailed
	}
	return exitOK
}

func parseBaselineSkillsRestoreCommand(args []string) (baselineSkillsRestoreCommandRequest, error) {
	request := baselineSkillsRestoreCommandRequest{repo: ".", format: "text"}
	flags := flag.NewFlagSet("baseline skills restore", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&request.repo, "repo", ".", "Git worktree or a path inside it")
	flags.StringVar(&request.profile, "profile", "", "Built-in Baseline Profile")
	flags.Var(&request.skills, "skill", "External profile skill to restore (repeatable)")
	flags.StringVar(&request.sourceDir, "source-dir", "", "Offline Git checkout or bare object store")
	flags.StringVar(&request.confirmation, "confirm-plan", "", "Exact approved restoration Plan Digest")
	flags.StringVar(&request.format, "format", "text", "Output format: text or json")
	if err := flags.Parse(args); err != nil {
		return request, validationError{
			message: fmt.Sprintf(
				"invalid baseline skills restore arguments: %v; run '%s baseline skills restore --help' for usage",
				err,
				app.Name,
			),
		}
	}
	if remaining := flags.Args(); len(remaining) != 0 {
		return request, validationError{
			message: fmt.Sprintf(
				"unexpected argument %q; run '%s baseline skills restore --help' for usage",
				remaining[0],
				app.Name,
			),
		}
	}
	request.repo = strings.TrimSpace(request.repo)
	request.profile = strings.TrimSpace(request.profile)
	request.sourceDir = strings.TrimSpace(request.sourceDir)
	request.confirmation = strings.TrimSpace(request.confirmation)
	request.format = strings.TrimSpace(request.format)
	if request.repo == "" {
		return request, validationError{message: "--repo cannot be empty"}
	}
	if request.profile == "" {
		return request, validationError{message: "--profile is required"}
	}
	for index, skill := range request.skills {
		request.skills[index] = strings.TrimSpace(skill)
		if request.skills[index] == "" {
			return request, validationError{message: "--skill cannot be empty"}
		}
	}
	if request.sourceDir != "" {
		absolute, err := filepath.Abs(request.sourceDir)
		if err != nil {
			return request, validationError{
				message: fmt.Sprintf("resolve --source-dir: %v", err),
			}
		}
		request.sourceDir = absolute
	}
	if request.format != "text" && request.format != "json" {
		return request, validationError{
			message: fmt.Sprintf("unsupported --format %q; use text or json", request.format),
		}
	}
	return request, nil
}

func writeBaselineSkillsRestoreFailure(
	payload baseline.SkillsRestorePayload,
	err error,
	exit int,
	jsonOutput bool,
	stdout, stderr io.Writer,
) int {
	fmt.Fprintf(stderr, "%s: baseline skills restore failed: %v\n", app.Name, err)
	if writeErr := writeBaselineSkillsRestorePayload(payload, jsonOutput, stdout); writeErr != nil {
		fmt.Fprintf(stderr, "%s: baseline skills restore output failed: %v\n", app.Name, writeErr)
		return exitRunFailed
	}
	return exit
}

func writeBaselineSkillsRestorePayload(
	payload baseline.SkillsRestorePayload,
	jsonOutput bool,
	stdout io.Writer,
) error {
	if jsonOutput {
		encoder := json.NewEncoder(stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(payload)
	}
	status := "blocked"
	switch {
	case payload.Applied:
		status = "applied"
	case payload.OK:
		status = "no changes"
	}
	if _, err := fmt.Fprintf(stdout, "Baseline skills restore: %s\n", status); err != nil {
		return err
	}
	if payload.Finding != nil {
		if _, err := fmt.Fprintf(
			stdout,
			"%s: %s\nNext action: %s\n",
			payload.Finding.Code,
			payload.Finding.Message,
			payload.Finding.Action,
		); err != nil {
			return err
		}
	}
	if payload.PlanDigest != nil {
		if _, err := fmt.Fprintf(stdout, "Plan Digest: %s\n", *payload.PlanDigest); err != nil {
			return err
		}
	}
	for _, change := range payload.PlannedChanges {
		if _, err := fmt.Fprintf(
			stdout,
			"- %s %s [%s]\n",
			change.Action,
			change.Path,
			change.Skill,
		); err != nil {
			return err
		}
	}
	return nil
}

func baselineSkillsRestoreExit(err error) int {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return exitSIGINT
	}
	var restoreErr *baseline.SkillsRestoreError
	if !errors.As(err, &restoreErr) {
		return exitRunFailed
	}
	switch restoreErr.Category {
	case baseline.SkillsRestoreInvalid:
		return exitPreflight
	case baseline.SkillsRestoreAction:
		return exitUnverified
	default:
		return exitRunFailed
	}
}

func baselineSkillsRestoreJSONRequested(args []string) bool {
	for index, arg := range args {
		if arg == "--format=json" {
			return true
		}
		if arg == "--format" && index+1 < len(args) && strings.TrimSpace(args[index+1]) == "json" {
			return true
		}
	}
	return false
}
