package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strings"

	"roundfix/internal/app"
	"roundfix/internal/store"
)

func runRunsShowCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("runs show"))
		return exitOK
	}
	fail := func(err error) int {
		fmt.Fprintf(stderr, "%s: runs show failed: %v\n", app.Name, err)
		fmt.Fprintf(stderr, "Run '%s runs show --help' for usage.\n", app.Name)
		return exitPreflight
	}
	// The public syntax places the ID before its optional flags. Accept flags
	// before the ID too, without letting the stdlib parser ignore trailing flags.
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		args = append(append([]string{}, args[1:]...), args[0])
	}
	fs := flag.NewFlagSet("runs show", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	asJSON := fs.Bool("json", false, "Print schema roundfix/runs-show/v1")
	if err := fs.Parse(args); err != nil {
		return fail(err)
	}
	if fs.NArg() != 1 {
		return fail(fmt.Errorf("runs show requires exactly one Run ID"))
	}
	id := fs.Arg(0)
	loaded, err := loadCommandConfig(environment, stderr)
	if err != nil {
		return fail(err)
	}
	reader, found, err := openRunsListReader(ctx, loaded.HomeDir)
	if err != nil {
		return fail(err)
	}
	if !found {
		return fail(fmt.Errorf("%s", unknownRunMessage(ctx, id, loaded.Config.Store.RunRetentionDays)))
	}
	defer func() { _ = reader.Close() }()
	run, found, err := reader.Run(ctx, id)
	if err != nil {
		return fail(err)
	}
	if !found {
		return fail(fmt.Errorf("%s", unknownRunMessage(ctx, id, loaded.Config.Store.RunRetentionDays)))
	}
	report, err := reader.RunTokenUsage(ctx, id)
	if err != nil {
		return fail(err)
	}
	if *asJSON {
		result := struct {
			Schema   string `json:"schema"`
			RunID    string `json:"run_id"`
			Kind     string `json:"kind"`
			SpecSlug string `json:"spec_slug"`
			State    string `json:"state"`
			store.TokenUsageReport
		}{"roundfix/runs-show/v1", run.ID, run.Kind, run.SpecSlug, run.State, report}
		if err := json.NewEncoder(stdout).Encode(result); err != nil {
			return fail(err)
		}
		return exitOK
	}
	fmt.Fprintf(stdout, "Run %s: %s %s %s\n", run.ID, run.Kind, run.SpecSlug, run.State)
	for _, scope := range report.Scopes {
		selections := make([]string, 0, len(scope.Selections))
		for _, selection := range scope.Selections {
			selections = append(selections, selection.Runtime+"/"+selection.Model+"/"+selection.ReasoningEffort)
		}
		tokens, cost := formatTokenTotals(scope.TokenTotals)
		if scope.Tokens != nil {
			number := fmt.Sprintf("%d", *scope.Tokens)
			tokens = strings.Replace(tokens, number, number+" tokens ("+strings.Join(scope.Bases, ", ")+")", 1)
		}
		fmt.Fprintf(stdout, "%s\t%s\t%s\t%s\n", scope.ScopeID, strings.Join(selections, ", "), tokens, cost)
	}
	printTokenTotals(stdout, report.Total)
	return exitOK
}
