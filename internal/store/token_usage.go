package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"roundfix/internal/runevent"
)

type TokenUsageRecord struct {
	RunID, ScopeKind, ScopeID, Session string
	Attempt                            int
	Runtime, Model, ReasoningEffort    string
	Basis                              string
	TotalTokens, InputTokens           *int64
	OutputTokens, CachedReadTokens     *int64
	CachedWriteTokens, ThoughtTokens   *int64
	Readings                           int
	CostAmount                         *float64
	CostCurrency                       string
	Time                               time.Time
}

type TokenCost struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}

type TokenSelection struct {
	Runtime         string `json:"runtime"`
	Model           string `json:"model"`
	ReasoningEffort string `json:"reasoning_effort"`
}

type TokenTotals struct {
	Prompts           int         `json:"prompts"`
	ReportedPrompts   int         `json:"reported_prompts"`
	Tokens            *int64      `json:"tokens"`
	Bases             []string    `json:"bases"`
	InputTokens       *int64      `json:"input_tokens"`
	OutputTokens      *int64      `json:"output_tokens"`
	CachedReadTokens  *int64      `json:"cached_read_tokens"`
	CachedWriteTokens *int64      `json:"cached_write_tokens"`
	ThoughtTokens     *int64      `json:"thought_tokens"`
	Costs             []TokenCost `json:"costs"`
	Sessions          int         `json:"sessions"`
	CostSessions      int         `json:"cost_sessions"`
	Runs              int         `json:"runs,omitempty"`
}

type ScopeTokenUsage struct {
	ScopeKind  string           `json:"scope_kind"`
	ScopeID    string           `json:"scope_id"`
	Selections []TokenSelection `json:"selections"`
	TokenTotals
}

type TokenUsageReport struct {
	Scopes []ScopeTokenUsage `json:"scopes"`
	Total  TokenTotals       `json:"total"`
}

func tokenUsageSchemaStatements(maxTokensExists bool) []string {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS run_token_usage (
 id INTEGER PRIMARY KEY AUTOINCREMENT,
 run_id TEXT NOT NULL REFERENCES runs(id) ON DELETE CASCADE,
 scope_kind TEXT NOT NULL CHECK (scope_kind IN ('task','qa','review','session')),
 scope_id TEXT NOT NULL, session TEXT NOT NULL, attempt INTEGER NOT NULL DEFAULT 0,
 runtime TEXT NOT NULL DEFAULT '', model TEXT NOT NULL DEFAULT '', reasoning_effort TEXT NOT NULL DEFAULT '',
 basis TEXT NOT NULL CHECK (basis IN ('unreported','turn','request-sum')),
 total_tokens INTEGER, input_tokens INTEGER, output_tokens INTEGER,
 cached_read_tokens INTEGER, cached_write_tokens INTEGER, thought_tokens INTEGER,
 readings INTEGER NOT NULL DEFAULT 0, cost_amount REAL, cost_currency TEXT NOT NULL DEFAULT '',
 created_at TEXT NOT NULL,
 CHECK ((basis = 'unreported' AND total_tokens IS NULL) OR (basis != 'unreported' AND total_tokens IS NOT NULL)))`,
		`CREATE INDEX IF NOT EXISTS idx_run_token_usage_scope ON run_token_usage (run_id, scope_kind, scope_id)`,
		`CREATE TABLE IF NOT EXISTS delivery_queue_runs (
 git_root TEXT NOT NULL REFERENCES delivery_queues(git_root) ON DELETE CASCADE,
 spec_slug TEXT NOT NULL, run_id TEXT NOT NULL, PRIMARY KEY (git_root, run_id))`,
	}
	if !maxTokensExists {
		statements = append(statements, `ALTER TABLE delivery_queues ADD COLUMN max_tokens INTEGER NOT NULL DEFAULT 0`)
	}
	return statements
}

func linkDeliveryQueueRun(ctx context.Context, tx *sql.Tx, gitRoot string, item DeliveryQueueItem) error {
	if item.RunID == "" {
		return nil
	}
	if _, err := tx.ExecContext(ctx, `INSERT OR IGNORE INTO delivery_queue_runs (git_root, spec_slug, run_id) VALUES (?, ?, ?)`, gitRoot, item.SpecSlug, item.RunID); err != nil {
		return fmt.Errorf("link Delivery Queue Run %q: %w", item.RunID, err)
	}
	return nil
}

// AppendTokenUsage commits the usage row and its event atomically, after any
// pending journal events, just like an Agent Selection attempt.
func (store *Store) AppendTokenUsage(ctx context.Context, record TokenUsageRecord) error {
	if err := store.FlushJournal(ctx); err != nil {
		return fmt.Errorf("flush journal before token usage: %w", err)
	}
	if record.Time.IsZero() {
		record.Time = store.now()
	}
	payload := map[string]any{
		"scope_kind": record.ScopeKind, "scope_id": record.ScopeID, "session": record.Session,
		"attempt": record.Attempt, "runtime": record.Runtime, "model": record.Model,
		"reasoning_effort": record.ReasoningEffort, "basis": record.Basis, "readings": record.Readings,
	}
	for key, value := range map[string]*int64{"total_tokens": record.TotalTokens, "input_tokens": record.InputTokens, "output_tokens": record.OutputTokens, "cached_read_tokens": record.CachedReadTokens, "cached_write_tokens": record.CachedWriteTokens, "thought_tokens": record.ThoughtTokens} {
		if value != nil {
			payload[key] = *value
		}
	}
	if record.CostAmount != nil {
		payload["cost"] = TokenCost{Currency: record.CostCurrency, Amount: *record.CostAmount}
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("encode token usage: %w", err)
	}
	return store.withWriteTx(ctx, "token usage append", func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `INSERT INTO run_token_usage (
 run_id,scope_kind,scope_id,session,attempt,runtime,model,reasoning_effort,basis,
 total_tokens,input_tokens,output_tokens,cached_read_tokens,cached_write_tokens,thought_tokens,
 readings,cost_amount,cost_currency,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
			record.RunID, record.ScopeKind, record.ScopeID, record.Session, record.Attempt, record.Runtime, record.Model, record.ReasoningEffort, record.Basis,
			record.TotalTokens, record.InputTokens, record.OutputTokens, record.CachedReadTokens, record.CachedWriteTokens, record.ThoughtTokens,
			record.Readings, record.CostAmount, record.CostCurrency, formatTime(record.Time)); err != nil {
			return fmt.Errorf("insert token usage: %w", err)
		}
		_, err := appendRunEvent(ctx, tx, runevent.RunEvent{RunID: record.RunID, Source: runevent.SourceDaemon, Kind: runevent.KindDaemonTokenUsage, ReviewIssue: record.ScopeID, Time: record.Time, Payload: raw})
		return err
	})
}

func (store *Store) RunTokenUsage(ctx context.Context, runID string) (TokenUsageReport, error) {
	return store.readTokenUsage(ctx, `WHERE usage.run_id = ?`, runID, false)
}

func (store *Store) DeliveryQueueTokenUsage(ctx context.Context, gitRoot string) (TokenUsageReport, error) {
	return store.readTokenUsage(ctx, `JOIN delivery_queue_runs linked ON linked.run_id = usage.run_id WHERE linked.git_root = ?`, gitRoot, true)
}

func (store *Store) readTokenUsage(ctx context.Context, clause, identity string, queue bool) (TokenUsageReport, error) {
	report := TokenUsageReport{Scopes: []ScopeTokenUsage{}}
	rows, err := store.db.QueryContext(ctx, `SELECT usage.run_id,scope_kind,scope_id,session,runtime,model,reasoning_effort,basis,total_tokens,input_tokens,output_tokens,cached_read_tokens,cached_write_tokens,thought_tokens,cost_amount,cost_currency FROM run_token_usage usage `+clause+` ORDER BY usage.id`, identity)
	if err != nil {
		return report, fmt.Errorf("read token usage: %w", err)
	}
	defer func() { _ = rows.Close() }()
	total := newTokenAccumulator()
	scopes := map[[2]string]int{}
	accumulators := []*tokenAccumulator{}
	for rows.Next() {
		var r TokenUsageRecord
		if err := rows.Scan(&r.RunID, &r.ScopeKind, &r.ScopeID, &r.Session, &r.Runtime, &r.Model, &r.ReasoningEffort, &r.Basis, &r.TotalTokens, &r.InputTokens, &r.OutputTokens, &r.CachedReadTokens, &r.CachedWriteTokens, &r.ThoughtTokens, &r.CostAmount, &r.CostCurrency); err != nil {
			return report, fmt.Errorf("scan token usage: %w", err)
		}
		key := [2]string{r.ScopeKind, r.ScopeID}
		index, ok := scopes[key]
		if !ok {
			index = len(report.Scopes)
			scopes[key] = index
			report.Scopes = append(report.Scopes, ScopeTokenUsage{ScopeKind: r.ScopeKind, ScopeID: r.ScopeID, Selections: []TokenSelection{}})
			accumulators = append(accumulators, newTokenAccumulator())
		}
		selection := TokenSelection{r.Runtime, r.Model, r.ReasoningEffort}
		found := false
		for _, old := range report.Scopes[index].Selections {
			if old == selection {
				found = true
				break
			}
		}
		if !found {
			report.Scopes[index].Selections = append(report.Scopes[index].Selections, selection)
		}
		total.add(r)
		accumulators[index].add(r)
	}
	if err := rows.Err(); err != nil {
		return report, fmt.Errorf("iterate token usage: %w", err)
	}
	for i, a := range accumulators {
		report.Scopes[i].TokenTotals = a.totals
	}
	report.Total = total.totals
	if queue {
		if err := store.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM delivery_queue_runs WHERE git_root = ?`, identity).Scan(&report.Total.Runs); err != nil {
			return report, fmt.Errorf("count Delivery Queue Runs: %w", err)
		}
	}
	return report, nil
}

type tokenAccumulator struct {
	totals        TokenTotals
	sessions      map[[2]string]bool
	previousCosts map[[3]string]float64
}

func newTokenAccumulator() *tokenAccumulator {
	return &tokenAccumulator{totals: TokenTotals{Bases: []string{}, Costs: []TokenCost{}}, sessions: map[[2]string]bool{}, previousCosts: map[[3]string]float64{}}
}

func (a *tokenAccumulator) add(r TokenUsageRecord) {
	t := &a.totals
	t.Prompts++
	session := [2]string{r.RunID, r.Session}
	if _, ok := a.sessions[session]; !ok {
		a.sessions[session] = false
		t.Sessions++
	}
	if r.Basis != "unreported" {
		t.ReportedPrompts++
		if t.Tokens == nil {
			t.Tokens = new(int64)
		}
		*t.Tokens += *r.TotalTokens
		found := false
		for _, basis := range t.Bases {
			if basis == r.Basis {
				found = true
			}
		}
		if !found {
			t.Bases = append(t.Bases, r.Basis)
		}
		for _, pair := range [][2]**int64{{&t.InputTokens, &r.InputTokens}, {&t.OutputTokens, &r.OutputTokens}, {&t.CachedReadTokens, &r.CachedReadTokens}, {&t.CachedWriteTokens, &r.CachedWriteTokens}, {&t.ThoughtTokens, &r.ThoughtTokens}} {
			dest, value := pair[0], *pair[1]
			if t.ReportedPrompts == 1 && value != nil {
				*dest = new(int64)
			}
			if value == nil {
				*dest = nil
			} else if *dest != nil {
				**dest += *value
			}
		}
	}
	if r.CostAmount != nil {
		if !a.sessions[session] {
			a.sessions[session] = true
			t.CostSessions++
		}
		key := [3]string{r.RunID, r.Session, r.CostCurrency}
		amount := *r.CostAmount
		if previous, ok := a.previousCosts[key]; ok && amount >= previous {
			amount -= previous
		}
		a.previousCosts[key] = *r.CostAmount
		index := -1
		for i, cost := range t.Costs {
			if cost.Currency == r.CostCurrency {
				index = i
				break
			}
		}
		if index < 0 {
			t.Costs = append(t.Costs, TokenCost{Currency: r.CostCurrency})
			index = len(t.Costs) - 1
		}
		t.Costs[index].Amount += amount
	}
}

func tokenUsageMigrationStatements(ctx context.Context, query migrationQuerier) ([]string, error) {
	var maxTokensExists int
	if err := query.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM pragma_table_info('delivery_queues') WHERE name = 'max_tokens')`).Scan(&maxTokensExists); err != nil {
		return nil, fmt.Errorf("inspect Delivery Queue token limit column: %w", err)
	}
	return tokenUsageSchemaStatements(maxTokensExists != 0), nil
}
