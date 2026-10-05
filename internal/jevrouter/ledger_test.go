package jevrouter

// Suite: router records read through the shared judge reader and private modes.
import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"roundfix/internal/judge"
)

func TestRouterLineIsReadByTheJudgeLog(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	now := time.Date(2026, 9, 30, 22, 0, 0, 0, time.FixedZone("UTC-3", -3*3600))
	ledger := Ledger{HomeDir: home}
	record := PromptRecord{RunID: "run-fixture", Spec: "spec-fixture", ScopeKind: "task", ScopeID: "task_02", Category: "docs", Repository: "/fixture/repo", Attempt: 2, UsageBefore: 0.25, UsageAfter: 0.75, Latency: 1500 * time.Millisecond, InputTokens: 100, OutputTokens: 20}
	if err := ledger.Append(record, now); err != nil {
		t.Fatal(err)
	}
	rows, err := judge.ReadMonth(context.Background(), home, now)
	if err != nil || len(rows) != 1 {
		t.Fatalf("rows=%+v err=%v", rows, err)
	}
	want := judge.LogLine{Schema: "roundfix/judge-log/v1", Time: now.UTC(), Repository: record.Repository, Spec: record.Spec, Judgment: "router-prompt", Target: "run-fixture task task_02", QuestionID: "docs", Transport: "openrouter", RequestedModel: "roundfix-openrouter/typesafe/jev-router", LatencyMS: 1500, InputTokens: 100, OutputTokens: 20, CostUSD: 0.5, CostSource: "reported", Attempts: 2, Outcome: "clear"}
	if !reflect.DeepEqual(rows[0], want) {
		t.Fatalf("row=%+v want=%+v", rows[0], want)
	}
	path := filepath.Join(home, ".roundfix", "judge", "2026-10.jsonl")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		t.Fatal(err)
	}
	if len(fields) != 28 {
		t.Fatalf("fields=%v", fields)
	}
	for _, name := range []string{"answer", "probabilities", "confidence", "noul"} {
		if string(fields[name]) != "null" {
			t.Fatalf("%s=%s", name, fields[name])
		}
	}
	for _, tc := range []struct {
		path string
		mode os.FileMode
	}{{path, 0600}, {filepath.Dir(path), 0700}} {
		info, err := os.Stat(tc.path)
		if err != nil || info.Mode().Perm() != tc.mode {
			t.Fatalf("mode for %s: %v err=%v", tc.path, info, err)
		}
	}
	// Append, rather than replace: a failed prompt still contributes its cost.
	record.Failed = true
	if err := ledger.Append(record, now); err != nil {
		t.Fatal(err)
	}
	rows, err = judge.ReadMonth(context.Background(), home, now)
	if err != nil || len(rows) != 2 || rows[0].CostUSD+rows[1].CostUSD != 1 || rows[1].Outcome != "skipped" {
		t.Fatalf("appended rows=%+v err=%v", rows, err)
	}
}

func TestLedgerRecordsUnreadableUsageAndClampsNegativeDelta(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	now := time.Now()
	ledger := Ledger{HomeDir: home}
	for _, record := range []PromptRecord{{UsageBefore: 1, UsageAfter: 0.5}, {UsageBefore: 1, UsageAfter: 2, UsageAfterErr: errors.New("HTTP 500"), Failed: true}} {
		if err := ledger.Append(record, now); err != nil {
			t.Fatal(err)
		}
	}
	rows, err := judge.ReadMonth(context.Background(), home, now)
	if err != nil || len(rows) != 2 || rows[0].CostUSD != 0 || rows[0].Error != "" || rows[1].CostUSD != 0 || rows[1].Error != "key usage unreadable after prompt: HTTP 500" || rows[1].Outcome != "skipped" {
		t.Fatalf("rows=%+v err=%v", rows, err)
	}
}

func TestRouterLineRecordsTheReportedModels(t *testing.T) {
	home := t.TempDir()
	now := time.Now()
	reported := Observation{Models: []string{"a/one", "b/two"}, Providers: []string{"p", "q"}, ResponseID: "last"}
	if err := (Ledger{HomeDir: home}).Append(PromptRecord{Reported: reported}, now); err != nil {
		t.Fatal(err)
	}
	rows, err := judge.ReadMonth(t.Context(), home, now)
	if err != nil || len(rows) != 1 {
		t.Fatalf("rows=%+v err=%v", rows, err)
	}
	if rows[0].Model != "a/one, b/two" || rows[0].Provider != "p, q" || rows[0].ResponseID != "last" {
		t.Fatalf("row=%+v", rows[0])
	}
}

func TestRouterLineRecordsACreditRefusal(t *testing.T) {
	for _, source := range []string{"openrouter_credits", ""} {
		for _, usageErr := range []error{nil, errors.New("HTTP 500")} {
			home := t.TempDir()
			now := time.Now()
			record := PromptRecord{Reported: Observation{Refusal: &Refusal{Status: 402, LimitSource: source}}, Failed: true, UsageBefore: 1, UsageAfter: 2, UsageAfterErr: usageErr}
			if err := (Ledger{HomeDir: home}).Append(record, now); err != nil {
				t.Fatal(err)
			}
			rows, err := judge.ReadMonth(t.Context(), home, now)
			if err != nil || len(rows) != 1 {
				t.Fatalf("rows=%+v err=%v", rows, err)
			}
			limit := source
			if limit == "" {
				limit = "unspecified"
			}
			want := "openrouter_credit_refused: " + limit
			cost := 1.0
			if usageErr != nil {
				want += "; key usage unreadable after prompt: HTTP 500"
				cost = 0
			}
			if rows[0].Judgment != "router-prompt" || rows[0].Error != want || rows[0].Outcome != "skipped" || rows[0].CostUSD != cost {
				t.Fatalf("row=%+v want error=%q", rows[0], want)
			}
		}
	}
}
