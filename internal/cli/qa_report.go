package cli

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"roundfix/internal/app"
	"roundfix/internal/spec"
)

const qaReportUsage = `Usage:
  roundfix qa-report accept <path>

Applies the declared-acceptance policy to the selected QA Report. The report's
Spec directory supplies its declarations. Produces no stdout; exit status is
the machine-readable result.

Exit codes:
  0  the report is acceptable
  1  the report is missing, unreadable, or unacceptable
  2  usage error
`

func runQAReportCommand(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) == 1 && commandWantsHelp(args) {
		fmt.Fprint(stdout, qaReportUsage)
		return exitOK
	}
	if len(args) == 0 || args[0] != "accept" {
		command := ""
		if len(args) > 0 {
			command = args[0]
		}
		fmt.Fprintf(stderr, "%s: unknown qa-report command %q\n", app.Name, command)
		fmt.Fprintf(stderr, "Run '%s qa-report --help' for usage.\n", app.Name)
		return exitPreflight
	}
	if len(args) != 2 || strings.TrimSpace(args[1]) == "" {
		fmt.Fprintf(stderr, "%s: qa-report accept requires exactly one report path\n", app.Name)
		fmt.Fprintf(stderr, "Run '%s qa-report --help' for usage.\n", app.Name)
		return exitPreflight
	}
	if err := ctx.Err(); err != nil {
		fmt.Fprintf(stderr, "%s: QA Report eligibility refused: %v\n", app.Name, err)
		return exitRunFailed
	}

	reportPath := filepath.Clean(args[1])
	if filepath.Base(filepath.Dir(reportPath)) != "qa" {
		fmt.Fprintf(stderr, "%s: QA Report eligibility refused: report path %q is not in a Spec qa directory\n", app.Name, args[1])
		return exitRunFailed
	}
	specDir := filepath.Dir(filepath.Dir(reportPath))
	report, err := spec.ReadQAReportFile(reportPath)
	if err != nil {
		fmt.Fprintf(stderr, "%s: QA Report eligibility refused: %v\n", app.Name, err)
		return exitRunFailed
	}
	if err := spec.QAReportEligibility(specDir, report); err != nil {
		fmt.Fprintf(stderr, "%s: QA Report eligibility refused: %v\n", app.Name, err)
		return exitRunFailed
	}
	return exitOK
}
