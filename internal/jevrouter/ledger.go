package jevrouter

import (
	"errors"
	"fmt"
	"time"

	"roundfix/internal/judge"
)

type PromptRecord struct {
	RunID, Spec, ScopeKind, ScopeID, Category, Repository string
	Attempt                                               int
	UsageBefore, UsageAfter                               float64
	UsageAfterErr                                         error
	Latency                                               time.Duration
	InputTokens, OutputTokens                             int64
	Failed                                                bool
}

type Ledger struct {
	HomeDir string
}

// Append records a prompt using the judge's schema and append-only writer.
func (ledger Ledger) Append(record PromptRecord, now time.Time) error {
	if ledger.HomeDir == "" {
		return errors.New("append router prompt: Home is empty")
	}
	row := judge.LogLine{
		Schema: "roundfix/judge-log/v1", Time: now.UTC(), Repository: record.Repository,
		Spec: record.Spec, Judgment: "router-prompt",
		Target:     fmt.Sprintf("%s %s %s", record.RunID, record.ScopeKind, record.ScopeID),
		QuestionID: record.Category, Transport: "openrouter",
		RequestedModel: "roundfix-openrouter/typesafe/jev-router",
		LatencyMS:      record.Latency.Milliseconds(), InputTokens: int(record.InputTokens), OutputTokens: int(record.OutputTokens),
		CostUSD: max(0, record.UsageAfter-record.UsageBefore), CostSource: "reported",
		Attempts: record.Attempt, Outcome: "clear",
	}
	if record.Failed {
		row.Outcome = "skipped"
	}
	if record.UsageAfterErr != nil {
		row.Error = "key usage unreadable after prompt: " + record.UsageAfterErr.Error()
		row.CostUSD = 0
	}
	if err := judge.AppendLogLine(ledger.HomeDir, now, row); err != nil {
		return fmt.Errorf("append router prompt: %w", err)
	}
	return nil
}
