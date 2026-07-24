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

type baselineAssetsSyncCommandRequest struct {
	sourceDir string
	check     bool
	format    string
}

type baselineAssetsSyncExecutor func(
	context.Context,
	baseline.AssetsSyncRequest,
) (baseline.AssetsSyncPayload, error)

func runBaselineAssetsSyncCommand(
	ctx context.Context,
	args []string,
	stdout, stderr io.Writer,
) int {
	return runBaselineAssetsSyncCommandWith(
		ctx,
		args,
		stdout,
		stderr,
		baseline.SyncAssets,
	)
}

func runBaselineAssetsSyncCommandWith(
	ctx context.Context,
	args []string,
	stdout, stderr io.Writer,
	execute baselineAssetsSyncExecutor,
) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("baseline assets sync"))
		return exitOK
	}
	request, err := parseBaselineAssetsSyncCommand(args)
	jsonOutput := request.format == "json" || baselineAssetsSyncJSONRequested(args)
	if err != nil {
		finding := baseline.AssetsSyncFinding{
			Code:      "assets.arguments-invalid",
			Severity:  "error",
			Path:      "assets/setups",
			ManagedID: "setups",
			Message:   err.Error(),
			Action:    "Correct the command input and rerun roundfix baseline assets sync.",
		}
		payload := baseline.AssetsSyncPayload{
			SchemaVersion:  baseline.AssetsSyncSchemaVersion,
			OK:             false,
			Summary:        baseline.AssetsSyncSummary{Errors: 1},
			Findings:       []baseline.AssetsSyncFinding{finding},
			PlannedChanges: []baseline.AssetsSyncChange{},
		}
		return writeBaselineAssetsSyncFailure(
			payload,
			err,
			exitPreflight,
			jsonOutput,
			stdout,
			stderr,
		)
	}
	payload, syncErr := execute(ctx, baseline.AssetsSyncRequest{
		SourceDir: request.sourceDir,
		Check:     request.check,
	})
	if syncErr != nil {
		return writeBaselineAssetsSyncFailure(
			payload,
			syncErr,
			baselineAssetsSyncExit(syncErr),
			jsonOutput,
			stdout,
			stderr,
		)
	}
	if err := writeBaselineAssetsSyncPayload(payload, jsonOutput, stdout); err != nil {
		fmt.Fprintf(stderr, "%s: baseline assets sync output failed: %v\n", app.Name, err)
		return exitRunFailed
	}
	return exitOK
}

func parseBaselineAssetsSyncCommand(args []string) (baselineAssetsSyncCommandRequest, error) {
	request := baselineAssetsSyncCommandRequest{format: "text"}
	flags := flag.NewFlagSet("baseline assets sync", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&request.sourceDir, "source-dir", "", "Canonical setups directory")
	flags.BoolVar(&request.check, "check", false, "Report snapshot drift without writing canonical assets")
	flags.StringVar(&request.format, "format", "text", "Output format: text or json")
	if err := flags.Parse(args); err != nil {
		return request, validationError{
			message: fmt.Sprintf(
				"invalid baseline assets sync arguments: %v; run '%s baseline assets sync --help' for usage",
				err,
				app.Name,
			),
		}
	}
	if remaining := flags.Args(); len(remaining) != 0 {
		return request, validationError{
			message: fmt.Sprintf(
				"unexpected argument %q; run '%s baseline assets sync --help' for usage",
				remaining[0],
				app.Name,
			),
		}
	}
	request.sourceDir = strings.TrimSpace(request.sourceDir)
	request.format = strings.TrimSpace(request.format)
	if request.sourceDir == "" {
		return request, validationError{message: "--source-dir is required"}
	}
	absolute, err := filepath.Abs(request.sourceDir)
	if err != nil {
		return request, validationError{
			message: fmt.Sprintf("resolve --source-dir: %v", err),
		}
	}
	request.sourceDir = filepath.Clean(absolute)
	if request.format != "text" && request.format != "json" {
		return request, validationError{
			message: fmt.Sprintf("unsupported --format %q; use text or json", request.format),
		}
	}
	return request, nil
}

func writeBaselineAssetsSyncFailure(
	payload baseline.AssetsSyncPayload,
	err error,
	exit int,
	jsonOutput bool,
	stdout, stderr io.Writer,
) int {
	fmt.Fprintf(stderr, "%s: baseline assets sync failed: %v\n", app.Name, err)
	if writeErr := writeBaselineAssetsSyncPayload(payload, jsonOutput, stdout); writeErr != nil {
		fmt.Fprintf(stderr, "%s: baseline assets sync output failed: %v\n", app.Name, writeErr)
		return exitRunFailed
	}
	return exit
}

func writeBaselineAssetsSyncPayload(
	payload baseline.AssetsSyncPayload,
	jsonOutput bool,
	stdout io.Writer,
) error {
	if jsonOutput {
		encoder := json.NewEncoder(stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(payload)
	}
	if len(payload.Findings) == 0 && len(payload.PlannedChanges) == 0 {
		_, err := fmt.Fprintln(stdout, "setup-context-driven audit: ok")
		return err
	}
	if _, err := fmt.Fprintln(stdout, "setup-context-driven audit: findings"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(
		stdout,
		"errors=%d decisions=%d warnings=%d info=%d\n",
		payload.Summary.Errors,
		payload.Summary.Decisions,
		payload.Summary.Warnings,
		payload.Summary.Info,
	); err != nil {
		return err
	}
	for _, severity := range []string{"error", "decision", "warning", "info"} {
		wroteHeading := false
		for _, finding := range payload.Findings {
			if finding.Severity != severity {
				continue
			}
			if !wroteHeading {
				if _, err := fmt.Fprintf(stdout, "%s:\n", severity); err != nil {
					return err
				}
				wroteHeading = true
			}
			location := finding.Path
			if finding.ManagedID != "" {
				location += " [" + finding.ManagedID + "]"
			}
			if _, err := fmt.Fprintf(
				stdout,
				"- %s %s: %s\n  action: %s\n",
				finding.Code,
				location,
				finding.Message,
				finding.Action,
			); err != nil {
				return err
			}
		}
	}
	return nil
}

func baselineAssetsSyncExit(err error) int {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return exitSIGINT
	}
	var syncErr *baseline.AssetsSyncError
	if !errors.As(err, &syncErr) {
		return exitRunFailed
	}
	if syncErr.Category == baseline.AssetsSyncInvalid {
		return exitPreflight
	}
	return exitRunFailed
}

func baselineAssetsSyncJSONRequested(args []string) bool {
	for index, arg := range args {
		if arg == "--format=json" {
			return true
		}
		if arg == "--format" && index+1 < len(args) &&
			strings.TrimSpace(args[index+1]) == "json" {
			return true
		}
	}
	return false
}
