package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"roundfix/internal/app"
	"roundfix/internal/baseline"
)

type baselineSkillsReconcileCommandRequest struct {
	repo             string
	profile          string
	sourceRepository string
	revision         string
	sourceDir        string
	confirmation     string
	format           string
}

func runBaselineSkillsReconcileCommand(
	ctx context.Context,
	args []string,
	stdout, stderr io.Writer,
) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("baseline skills reconcile"))
		return exitOK
	}
	request, err := parseBaselineSkillsReconcileCommand(args)
	jsonOutput := request.format == "json" || baselineSkillsRestoreJSONRequested(args)
	if err != nil {
		payload := baseline.SkillsReconcilePayload{
			SchemaVersion:  baseline.SkillsReconcileSchemaVersion,
			Profile:        request.profile,
			Acquisitions:   []baseline.RestoreAcquisition{},
			Skills:         []baseline.SkillsReconcileEntry{},
			PlannedChanges: []baseline.RestorePlannedChange{},
			Finding: &baseline.RestoreFinding{
				Code:    "reconcile.arguments-invalid",
				Message: err.Error(),
				Action:  "Correct the command input and rerun roundfix baseline skills reconcile.",
			},
		}
		return writeBaselineSkillsReconcileFailure(
			payload,
			err,
			exitPreflight,
			jsonOutput,
			stdout,
			stderr,
		)
	}
	payload, reconcileErr := baseline.ReconcileSkillsLock(ctx, baseline.SkillsReconcileRequest{
		Repository:       request.repo,
		ProfileID:        request.profile,
		SourceRepository: request.sourceRepository,
		Commit:           request.revision,
		SourceDir:        request.sourceDir,
		Confirmation:     request.confirmation,
	})
	if reconcileErr != nil {
		return writeBaselineSkillsReconcileFailure(
			payload,
			reconcileErr,
			baselineSkillsRestoreExit(reconcileErr),
			jsonOutput,
			stdout,
			stderr,
		)
	}
	if err := writeBaselineSkillsReconcilePayload(payload, jsonOutput, stdout); err != nil {
		fmt.Fprintf(stderr, "%s: baseline skills reconcile output failed: %v\n", app.Name, err)
		return exitRunFailed
	}
	return exitOK
}

func parseBaselineSkillsReconcileCommand(
	args []string,
) (baselineSkillsReconcileCommandRequest, error) {
	request := baselineSkillsReconcileCommandRequest{repo: ".", format: "text"}
	flags := flag.NewFlagSet("baseline skills reconcile", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&request.repo, "repo", ".", "Git worktree or a path inside it")
	flags.StringVar(&request.profile, "profile", "", "Built-in Baseline Profile")
	flags.StringVar(&request.sourceRepository, "source", "", "Source repository as owner/repo")
	flags.StringVar(&request.revision, "revision", "", "Exact immutable 40-hex source commit")
	flags.StringVar(&request.sourceDir, "source-dir", "", "Offline Git checkout or bare object store")
	flags.StringVar(&request.confirmation, "confirm-plan", "", "Exact approved reconciliation Plan Digest")
	flags.StringVar(&request.format, "format", "text", "Output format: text or json")
	if err := flags.Parse(args); err != nil {
		return request, validationError{
			message: fmt.Sprintf(
				"invalid baseline skills reconcile arguments: %v; run '%s baseline skills reconcile --help' for usage",
				err,
				app.Name,
			),
		}
	}
	if remaining := flags.Args(); len(remaining) != 0 {
		return request, validationError{
			message: fmt.Sprintf(
				"unexpected argument %q; run '%s baseline skills reconcile --help' for usage",
				remaining[0],
				app.Name,
			),
		}
	}
	request.repo = strings.TrimSpace(request.repo)
	request.profile = strings.TrimSpace(request.profile)
	request.sourceRepository = strings.TrimSpace(request.sourceRepository)
	request.revision = strings.TrimSpace(request.revision)
	request.sourceDir = strings.TrimSpace(request.sourceDir)
	request.confirmation = strings.TrimSpace(request.confirmation)
	request.format = strings.TrimSpace(request.format)
	if request.repo == "" {
		return request, validationError{message: "--repo cannot be empty"}
	}
	if request.profile == "" {
		return request, validationError{message: "--profile is required"}
	}
	if request.sourceRepository == "" {
		return request, validationError{message: "--source is required"}
	}
	if request.revision == "" {
		return request, validationError{message: "--revision is required"}
	}
	if request.sourceDir != "" {
		absolute, err := filepath.Abs(request.sourceDir)
		if err != nil {
			return request, validationError{message: fmt.Sprintf("resolve --source-dir: %v", err)}
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

func writeBaselineSkillsReconcileFailure(
	payload baseline.SkillsReconcilePayload,
	err error,
	exit int,
	jsonOutput bool,
	stdout, stderr io.Writer,
) int {
	fmt.Fprintf(stderr, "%s: baseline skills reconcile failed: %v\n", app.Name, err)
	if writeErr := writeBaselineSkillsReconcilePayload(payload, jsonOutput, stdout); writeErr != nil {
		fmt.Fprintf(stderr, "%s: baseline skills reconcile output failed: %v\n", app.Name, writeErr)
		return exitRunFailed
	}
	return exit
}

func writeBaselineSkillsReconcilePayload(
	payload baseline.SkillsReconcilePayload,
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
	if _, err := fmt.Fprintf(stdout, "Baseline skills reconcile: %s\n", status); err != nil {
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
