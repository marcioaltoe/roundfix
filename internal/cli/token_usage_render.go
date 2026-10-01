package cli

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"roundfix/internal/app"
	"roundfix/internal/store"
)

// formatTokenTotals keeps absent reports distinct from reported zero usage.
func formatTokenTotals(totals store.TokenTotals) (tokens, cost string) {
	switch {
	case totals.Prompts == 0:
		tokens = "no prompts recorded"
	case totals.Tokens == nil:
		tokens = fmt.Sprintf("none reported by %d prompt(s)", totals.Prompts)
	default:
		tokens = fmt.Sprintf("%d from %d of %d prompt(s)", *totals.Tokens, totals.ReportedPrompts, totals.Prompts)
	}
	cost = "cost not reported"
	if len(totals.Costs) > 0 {
		amounts := make([]string, 0, len(totals.Costs))
		for _, reported := range totals.Costs {
			amounts = append(amounts, fmt.Sprintf("%.2f %s", reported.Amount, reported.Currency))
		}
		cost = fmt.Sprintf("cost %s from %d of %d Agent Session(s)", strings.Join(amounts, " + "), totals.CostSessions, totals.Sessions)
	}
	return tokens, cost
}

func printTokenTotals(output io.Writer, totals store.TokenTotals) {
	tokens, cost := formatTokenTotals(totals)
	fmt.Fprintf(output, "Tokens: %s; %s\n", tokens, cost)
}

func printImplementTokenUsage(ctx context.Context, runStore *store.Store, runID string, stdout, stderr io.Writer) {
	report, err := runStore.RunTokenUsage(context.WithoutCancel(ctx), runID)
	if err != nil {
		fmt.Fprintf(stderr, "%s: warning: token usage unavailable for Run %s: %v\n", app.Name, runID, err)
		fmt.Fprintln(stdout, "Tokens: unavailable")
		return
	}
	if report.Total.Prompts == 0 {
		tokens, _ := formatTokenTotals(report.Total)
		fmt.Fprintf(stdout, "Tokens: %s\n", tokens)
		return
	}
	printTokenTotals(stdout, report.Total)
}

func printDeliveryUsage(output io.Writer, totals store.TokenTotals) {
	tokens, cost := formatTokenTotals(totals)
	if totals.Runs == 0 {
		fmt.Fprintln(output, "Usage: no Runs recorded")
		return
	}
	if totals.Tokens != nil {
		tokens = strings.Replace(tokens, strconv.FormatInt(*totals.Tokens, 10), strconv.FormatInt(*totals.Tokens, 10)+" tokens", 1)
	}
	fmt.Fprintf(output, "Usage: %s across %d Run(s); %s\n", tokens, totals.Runs, cost)
}
