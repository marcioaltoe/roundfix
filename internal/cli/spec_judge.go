package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/judge"
)

const specJudgeUsage = `Usage:
  roundfix spec judge <slug> [--stage <prd|techspec>] [--format <text|json>]

Raises advisory judgments for one active Spec. Never fails for a judgment.
Also asks whether each open Finding or Backlog Entry belongs with a source
the Spec adopted. A suggestion never gates.

Options:
  --stage   Judge prd or techspec (default: both)
  --format  Output format: text or json (default: text)

Set ROUNDFIX_OPENROUTER_API_KEY for OpenRouter, or ROUNDFIX_TYPESAFE_API_KEY
as the direct TypeSafe alternative. The generic OPENROUTER_API_KEY is not read.
The monthly ceiling is jev.monthly_ceiling_usd in User Config, US$5 by default,
across both transports. Every request is
recorded in the Judge Log: <home>/.roundfix/judge/<YYYY-MM>.jsonl (UTC month).

Exit codes:
  0  ran, including advisory, skipped, or stopped results
  2  usage error, unknown active Spec, or missing required artifact
`

type specJudgeRequest struct {
	slug, format string
	stage        judge.Stage
}

func parseSpecJudgeCommand(args []string) (specJudgeRequest, error) {
	req := specJudgeRequest{format: "text"}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		name, value, equal := strings.Cut(arg, "=")
		switch name {
		case "--stage", "--format":
			if !equal {
				i++
				if i >= len(args) {
					return req, fmt.Errorf("%s requires a value", name)
				}
				value = args[i]
			}
			if name == "--stage" {
				if value != "prd" && value != "techspec" {
					return req, fmt.Errorf("unsupported --stage %q; use prd or techspec", value)
				}
				req.stage = judge.Stage(value)
			} else {
				if value != "text" && value != "json" {
					return req, fmt.Errorf("unsupported --format %q; use text or json", value)
				}
				req.format = value
			}
		default:
			if strings.HasPrefix(arg, "-") {
				return req, fmt.Errorf("unknown flag %q", arg)
			}
			if req.slug != "" {
				return req, fmt.Errorf("unexpected argument %q", arg)
			}
			req.slug = arg
		}
	}
	if req.slug == "" {
		return req, fmt.Errorf("spec judge requires one Spec slug")
	}
	return req, nil
}

func runSpecJudgeCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, specJudgeUsage)
		return exitOK
	}
	fail := func(err error) int {
		fmt.Fprintf(stderr, "%s: spec judge failed: %v\nRun '%s spec judge --help' for usage.\n", app.Name, err, app.Name)
		return exitPreflight
	}
	req, err := parseSpecJudgeCommand(args)
	if err != nil {
		return fail(err)
	}
	loaded, err := loadCommandConfig(environment, stderr)
	if err != nil {
		return fail(err)
	}
	if loaded.GitRoot == "" {
		return fail(fmt.Errorf("spec judge requires a git repository working tree"))
	}
	root, err := roundconfig.ResolveSpecsRoot(loaded, loaded.GitRoot)
	if err != nil {
		return fail(err)
	}
	if req.slug == "." || req.slug == ".." || req.slug == "_archived" || filepath.Base(req.slug) != req.slug {
		return fail(fmt.Errorf("invalid Spec slug %q", req.slug))
	}
	specDir := filepath.Join(root.Path, req.slug)
	info, err := os.Stat(specDir)
	if os.IsNotExist(err) || (err == nil && !info.IsDir()) {
		return fail(fmt.Errorf("unknown active Spec slug %q", req.slug))
	}
	if err != nil {
		return fail(fmt.Errorf("read Spec: %w", err))
	}
	if req.stage == "techspec" {
		if _, err := os.Lstat(filepath.Join(specDir, "_techspec.md")); err != nil {
			return fail(fmt.Errorf("read Spec TechSpec: %w", err))
		}
	}
	questions, err := judge.Load()
	if err != nil {
		return fail(err)
	}
	keys := make(map[string]string)
	for _, entry := range environment.environ {
		name, value, _ := strings.Cut(entry, "=")
		if name == "ROUNDFIX_OPENROUTER_API_KEY" || name == "ROUNDFIX_TYPESAFE_API_KEY" {
			keys[name] = value
		}
	}
	report, err := judge.Run(ctx, questions.WithMonthlyCeiling(loaded.Config.Jev.MonthlyCeilingUSD), judge.Request{
		RepoRoot: loaded.GitRoot, SpecDir: specDir, Spec: req.slug, Stage: req.stage,
		Keys: keys, HomeDir: environment.homeDir,
		Transport: environment.dependencies.judgeTransport, Now: environment.dependencies.judgeNow,
	})
	if err != nil {
		return fail(err)
	}
	if req.format == "json" {
		if err := json.NewEncoder(stdout).Encode(report); err != nil {
			return fail(fmt.Errorf("write judge JSON: %w", err))
		}
	} else {
		fmt.Fprint(stdout, renderSpecJudgeText(report))
	}
	return exitOK
}

func renderSpecJudgeText(report judge.Report) string {
	var out strings.Builder
	for _, artifact := range report.ArtifactsSkipped {
		fmt.Fprintf(&out, "skipped %s: %s\n", artifact.Artifact, artifact.Reason)
	}
	advisory, suggested, clear, skipped := 0, 0, 0, 0
	for _, j := range report.Judgments {
		switch j.Outcome {
		case "advisory":
			advisory++
			fmt.Fprintf(&out, "advisory %s %s:%d %s", j.Kind, j.Artifact, j.Line, j.Target)
			if j.Kind == "citation-support" {
				fmt.Fprintf(&out, " %s at confidence %.2f: %s\n", *j.Answer, *j.Confidence, j.Text)
			} else {
				fmt.Fprintf(&out, ": P(delivers) %.2f\n", *j.Noul)
			}
		case "suggested":
			suggested++
			fmt.Fprintf(&out, "suggested %s %s → %s: P(same Spec) %.2f\n", j.Kind, j.Artifact, j.Target, *j.Noul)
		case "clear":
			clear++
		case "skipped":
			skipped++
			// Run-level reasons belong only on the summary, as in Transcripts 2, 3 and 5.
			if report.Skipped == nil && !(report.Stopped != nil && j.Reason != nil && *j.Reason == *report.Stopped) {
				if j.Kind == "source-grouping" {
					fmt.Fprintf(&out, "skipped %s %s → %s: %s\n", j.Kind, j.Artifact, j.Target, *j.Reason)
				} else {
					fmt.Fprintf(&out, "skipped %s %s:%d %s: %s\n", j.Kind, j.Artifact, j.Line, j.Target, *j.Reason)
				}
			}
		}
	}
	if report.Skipped != nil {
		fmt.Fprintf(&out, "Judge: skipped: %s; %d judgment(s) not asked\n", *report.Skipped, skipped)
		return out.String()
	}
	fmt.Fprintf(&out, "Judge: %d advisory, %d suggested, %d clear, %d skipped; %d call(s), %d input tokens, US$%.4f; month US$%.4f of US$%.2f; model %s via %s", advisory, suggested, clear, skipped, report.Calls, report.InputTokens, report.CostUSD, report.MonthCostUSD, report.MonthCeilingUSD, report.Model, *report.Transport)
	if report.Stopped != nil {
		fmt.Fprintf(&out, "; stopped: %s", *report.Stopped)
	}
	out.WriteByte('\n')
	return out.String()
}
