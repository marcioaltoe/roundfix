---
spec: 0204-a-run-that-reports-the-tokens-and-spend-it-used
prd: _prd.md
created: 2026-09-30
---

# A Run that reports the tokens and spend it used — Technical Spec

## Executive Summary

The acpx runner already reads every JSON-RPC line an ACP adapter prints for a
prompt, and drops two of them: the `usage_update` notifications and the
`usage` object of the prompt response. This design keeps both, reduces them to
one Turn Usage per prompt with the counting rule of ADR-0198, and returns it
on the prompt's result. The Daemon writes one row per prompt to a new Run
Database table and appends a `daemon.token_usage` Run Event in the same
write. Three readers sum those rows: `roundfix runs show`, the Implement Run
summary and `deliver status`. The Delivery Engine compares the queue's sum with
an optional `--max-tokens` ceiling before it starts an item, as it compares the
clock with the deadline.

The trade-off this design accepts is adapter knowledge in Roundfix. The
counting rule depends on which adapter lineage produced a report, because the
protocol does not yet say what a prompt response's `usage` covers. The
alternative, one rule for every adapter, would either undercount Codex turns
thirty- to fifty-fold or overcount Claude turns by summing readings taken
while a response streams. One boolean on the existing adapter lineage
contract, and a runtime guard, keep that knowledge small.

## Project Constraints

- Identifier strategy: not applicable — no new identifier; rows are keyed by
  the existing Run ID, Work Item scope and Agent Session name. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local JSON-RPC lines and the Run
  Database only; no credential, no proxy and no network call. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0198 and ADR-0199 (this Spec),
  ADR-0017, ADR-0020, ADR-0051, ADR-0008, ADR-0033, ADR-0098, ADR-0158,
  ADR-0164, ADR-0166, ADR-0167, ADR-0176 and ADR-0182 hold as the PRD states.
  ADR-0184 has this TechSpec state each changed command as a Surface
  Transcript. The gate is bound by ADR-0080, ADR-0088, ADR-0091, ADR-0096,
  ADR-0104, ADR-0117, ADR-0155 and ADR-0156, and by ADR-0093 and ADR-0094.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30 for the skill files, recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/events.md`,
  `.agents/skills/roundfix/references/implement.md`,
  `.agents/skills/roundfix/references/runs.md`. Sanctioned regeneration:
  `make skills-sync`, `make baseline-digests`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## System Architecture

No new package.

| Concern | Owner | File |
| --- | --- | --- |
| Turn Usage and the counting rule | `TurnUsage`, `countTurnUsage` | `internal/agent/usage.go` (new) |
| Reading the reports from the stream | `readPromptStream`, `handleStdoutLine` | `internal/agent/acpx_runner.go`, `internal/agent/acp_stream.go` |
| Which lineage reports its last request only | `adapterLineageContract` | `internal/agent/acpx_runner.go` |
| Turn Usage on the prompt result | `ExecuteResult.Usage` | `internal/agent/agent.go` |
| Usage rows, queue links, sums | `AppendTokenUsage`, `RunTokenUsage`, `DeliveryQueueTokenUsage` | `internal/store/token_usage.go` (new), `internal/store/store.go`, `internal/store/delivery.go` |
| Recording each prompt | `recordPromptUsage` | `internal/daemon/agent_session_owner.go` |
| Run Event kind and stream category | `KindDaemonTokenUsage`, `StreamCategoryUsage` | `internal/runevent/event.go`, `internal/runevent/stream.go` |
| One renderer for every surface | `formatTokenTotals` | `internal/cli/token_usage_render.go` (new) |
| `runs show` | `runRunsShowCommand` | `internal/cli/runs_show.go` (new), `internal/cli/runs.go` |
| Implement Run summary line | `printImplementTokensLine` | `internal/cli/implement.go` |
| `deliver status`, `deliver start` | `printDeliveryLimits`, `printDeliveryUsage`, `parseDeliverStart` | `internal/cli/deliver.go` |
| The token ceiling | `Engine.Run`, `Engine.Retry`, `PendingQuestionFor` | `internal/delivery/engine.go`, `internal/delivery/question.go` |
| Help text | `commandUsage` | `internal/cli/cli.go` |

## Implementation Design

### Interfaces

```go
// internal/agent/usage.go
type UsageBasis string

const (
	UsageBasisTurn       UsageBasis = "turn"
	UsageBasisRequestSum UsageBasis = "request-sum"
)

// TurnUsage is what one prompt's adapter reported. The zero value is an
// unreported prompt, so a result without usage compares equal to before.
type TurnUsage struct {
	Basis                                  UsageBasis // "" = unreported
	TotalTokens                            int64
	InputTokens, OutputTokens              *int64
	CachedReadTokens, CachedWriteTokens    *int64
	ThoughtTokens                          *int64
	Readings                               int
	Cost                                   *ReportedCost // last cumulative reading
}

type ReportedCost struct {
	Amount   float64
	Currency string
}
```

```go
// internal/store/token_usage.go
type TokenUsageRecord struct {
	RunID, ScopeKind, ScopeID, Session string
	Attempt                            int
	Runtime, Model, ReasoningEffort    string
	Basis                              string // "unreported", "turn", "request-sum"
	TotalTokens, InputTokens           *int64
	OutputTokens, CachedReadTokens     *int64
	CachedWriteTokens, ThoughtTokens   *int64
	Readings                           int
	CostAmount                         *float64
	CostCurrency                       string
	Time                               time.Time
}

func (store *Store) AppendTokenUsage(ctx context.Context, record TokenUsageRecord) error
func (store *Store) RunTokenUsage(ctx context.Context, runID string) (TokenUsageReport, error)
func (store *Store) DeliveryQueueTokenUsage(ctx context.Context, gitRoot string) (TokenUsageReport, error)
```

`internal/store` does not import `internal/agent`, so the record is flat and
the Daemon copies a `TurnUsage` into it. `TokenUsageReport` holds one `ScopeTokenUsage` per Work
Item scope, in the order the scope first recorded usage, and a `TokenTotals`
for the whole report: prompts, reported prompts, tokens, the bases seen, the
cost per currency, Agent Sessions and Agent Sessions that reported a cost,
and, for a queue, the number of Runs.

`DeliveryQueueLimits` gains `MaxTokens int64`. The delivery engine's store
interface gains `DeliveryQueueTokenUsage`.

### Reading the reports

`handleStdoutLine` gains one collector argument, like `stopReason`:

- A `session/update` whose `sessionUpdate` is `usage_update` adds its `used`
  to the prompt's readings and keeps its `cost` when present. It publishes no
  Agent Run Event and does not count as Agent output.
- A `result` that carries a non-null `usage` object keeps it. A later result
  without `usage`, which acpx prints after the adapter's response, does not
  erase it.
- A usage payload that does not parse is ignored and never fails the prompt.

`RunPrompt` sets `result.Usage` on every return path after the stream was
read, including a stopped prompt and a nonzero acpx exit after a parsed
result.

### The counting rule

`countTurnUsage(lastRequestOnly bool, reported *promptUsage, readings []int64, cost *ReportedCost)`
follows ADR-0198:

1. `lastRequestOnly` is true only for the `codex` lineage, a new field of
   `adapterLineageContract`.
2. When `lastRequestOnly` and at least one reading exists, and the reported
   total is absent or not greater than the last reading: basis `request-sum`,
   `TotalTokens` is the sum of the readings, and every split field is nil.
3. Otherwise, when a reported `usage` exists: basis `turn`, `TotalTokens` is
   its `totalTokens`, and each split field is set exactly when the report
   carried it.
4. Otherwise: unreported, `Basis` empty, whatever the readings.

`Readings` is the number of readings, and `Cost` is the last `cost` seen, in
every case.

### Data Models

Schema version 22 adds, in one migration:

```sql
CREATE TABLE run_token_usage (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	run_id TEXT NOT NULL REFERENCES runs(id) ON DELETE CASCADE,
	scope_kind TEXT NOT NULL CHECK (scope_kind IN ('task','qa','review','session')),
	scope_id TEXT NOT NULL, session TEXT NOT NULL, attempt INTEGER NOT NULL DEFAULT 0,
	runtime TEXT NOT NULL DEFAULT '', model TEXT NOT NULL DEFAULT '',
	reasoning_effort TEXT NOT NULL DEFAULT '',
	basis TEXT NOT NULL CHECK (basis IN ('unreported','turn','request-sum')),
	total_tokens INTEGER, input_tokens INTEGER, output_tokens INTEGER,
	cached_read_tokens INTEGER, cached_write_tokens INTEGER, thought_tokens INTEGER,
	readings INTEGER NOT NULL DEFAULT 0,
	cost_amount REAL, cost_currency TEXT NOT NULL DEFAULT '',
	created_at TEXT NOT NULL);
CREATE INDEX idx_run_token_usage_scope ON run_token_usage (run_id, scope_kind, scope_id);
CREATE TABLE delivery_queue_runs (
	git_root TEXT NOT NULL REFERENCES delivery_queues(git_root) ON DELETE CASCADE,
	spec_slug TEXT NOT NULL, run_id TEXT NOT NULL, PRIMARY KEY (git_root, run_id));
ALTER TABLE delivery_queues ADD COLUMN max_tokens INTEGER NOT NULL DEFAULT 0;
```

Rules:

- `total_tokens` is NULL exactly when `basis` is `unreported`.
- `UpdateDeliveryQueueItem` and `RetryDeliveryQueueItem` insert
  `(git_root, spec_slug, run_id)` with `INSERT OR IGNORE` in the same write
  whenever the item's Run ID is not empty. Every Run ID ever recorded on an
  item is therefore linked, and replacing the queue removes its links.
- Journal Retention and the GC Command never delete either table. Both are
  added to the durable-table lifecycle policy in
  `docs/user-guide/run-database-lifecycle.md`.
- An Agent Session's spend is the sum, over its rows in ID order, of each
  cost reading's increase over the previous reading of the same currency;
  the first reading and a reading lower than the previous one count whole.

### Recording

`recordPromptUsage` runs after each prompt: in `agentSessionOwner.runPrepared`
for the owned path, with the owner's scope, attempt and active candidate, and
in `Engine.runAgentSession` for a prompt without an owner, with scope kind
`session`, the Agent Session name as scope ID and the request's runtime and
model. It writes through an optional `tokenUsageAppender` interface on
`deps.Runs`, as `persistSelectionAttempt` does, under
`context.WithoutCancel`. A write error prints `roundfix: warning: token usage
not recorded for <scope_kind> <scope_id>: <error>` to `deps.Progress` and is
otherwise ignored.

`AppendTokenUsage` inserts the row and appends a Run Event of kind
`daemon.token_usage`, source `daemon`, Work Item the scope ID, with this
payload:

```json
{"scope_kind":"task","scope_id":"task_01","session":"roundfix-run_1-task_01","attempt":1,"runtime":"codex","model":"gpt-6.1-sol","reasoning_effort":"high","basis":"request-sum","total_tokens":5639755,"readings":46}
```

An unreported prompt has `"basis":"unreported"` and no `total_tokens`. Split
fields and `cost` (`{"amount":3.55,"currency":"USD"}`) appear only when
recorded. `IsDaemonKind` does not include the new kind, so the Live Run View
skips it as it skips every kind it does not know.

### Surfaces

`formatTokenTotals` renders a `TokenTotals` for every surface as two
phrases:

- the token phrase: `no prompts recorded` when there are no rows,
  `none reported by <n> prompt(s)` when no prompt reported, and otherwise
  `<tokens> from <reported> of <n> prompt(s)`;
- the cost phrase: `cost not reported`, or `cost <amount> <currency> from
  <c> of <s> Agent Session(s)`, with the amount rounded to two decimals and
  several currencies joined by ` + `.

`runs show` prints a header, one tab-separated line per scope, and
`Tokens: <token phrase>; <cost phrase>`. A scope line holds the scope ID, its
Agent Selections as `runtime/model/effort` joined by `, `, its token phrase
with `tokens (<bases>)` after the number, and its cost phrase. The Implement
Run summary prints the same `Tokens:` line right after its outcome line, on
every outcome; with no rows it prints `Tokens: no prompts recorded`.
`deliver status` prints, right after the limits line,
`Usage: <tokens> tokens from <reported> of <n> prompt(s) across <k> Run(s); <cost phrase>`,
`Usage: none reported by <n> prompt(s) across <k> Run(s); <cost phrase>`, or
`Usage: no Runs recorded`.

The `usage` stream category projects a `daemon.token_usage` event into the
existing record fields `scope_kind`, `scope_id`, `runtime`, `model` and
`reasoning_effort`, and new fields `token_basis`, `tokens`, `input_tokens`,
`output_tokens`, `cached_read_tokens`, `cached_write_tokens`,
`thought_tokens`, `cost_amount` and `cost_currency`, each omitted when not
recorded. Its summary is `<scope_id> used <n> tokens (<basis>)` or
`<scope_id> reported no usage`.

### The ceiling

In `Engine.Run`, right after the deadline check and only for an item in stage
`queued`, a queue with `MaxTokens > 0` reads `DeliveryQueueTokenUsage`. When
its tokens are at or above `MaxTokens`, the item parks with blocker
`queue-token-ceiling`, before any branch or worktree exists. `Engine.Retry`
refuses any retry under the same condition. `PendingQuestionFor` answers that
blocker with `record a new queue for the remaining Specs with roundfix deliver
start and a higher --max-tokens`. Nothing else in the loop reads the ceiling.

### Surface Transcripts

1. Surface Transcript: `runs show` for an Implement Run with a Codex Task, a
   Claude Task and an unreported QA gate.

   ```transcript
   $ roundfix runs show run_20261001T120000Z_0123456789abcdef
   stdout:
   Run run_20261001T120000Z_0123456789abcdef: implement 0300-example Clean
   task_01	codex/gpt-6.1-sol/high	5639755 tokens (request-sum) from 1 of 1 prompt(s)	cost not reported
   task_02	claude/opus/high	4019995 tokens (turn) from 1 of 1 prompt(s)	cost 3.55 USD from 1 of 1 Agent Session(s)
   qa	codex/gpt-6.1-sol/high	none reported by 1 prompt(s)	cost not reported
   Tokens: 9659750 from 2 of 3 prompt(s); cost 3.55 USD from 1 of 3 Agent Session(s)
   stderr:
   exit: 0
   ```

2. Surface Transcript: `runs show` for a Run that does not exist.

   ```transcript
   $ roundfix runs show run_missing
   stdout:
   stderr:
   roundfix: runs show failed: Run "run_missing" does not exist
   Run 'roundfix runs show --help' for usage.
   exit: 2
   ```

3. Surface Transcript: the end of the Implement Run summary, for a Run whose
   adapter reported nothing.

   ```transcript
   $ roundfix implement --spec 0300-example --no-input
   stdout:
   task_01 completed — Write the widget guide
   task_02 completed — Build the widget backend
   Clean: all 2 Task(s) completed.
   Tokens: none reported by 2 prompt(s); cost not reported
   stderr:
   ...
   exit: 0
   ```

4. Surface Transcript: `deliver status` with a ceiling reached.

   ```transcript
   $ roundfix deliver status
   stdout:
   0300-example	merged	-	-
   0301-example	parked	queue-token-ceiling	-
   Limits: deadline none, retries per item none, concurrency 1, tokens 5000000
   Usage: 5639755 tokens from 1 of 1 prompt(s) across 1 Run(s); cost not reported
   Pending question: 0301-example parked queue-token-ceiling
   Answer: record a new queue for the remaining Specs with roundfix deliver start and a higher --max-tokens
   stderr:
   exit: 0
   ```

5. Surface Transcript: `deliver start` without a ceiling.

   ```transcript
   $ roundfix deliver start 0300-example
   stdout:
   Limits: deadline none, retries per item none, concurrency 1, tokens none
   stderr:
   exit: 0
   ```

6. Surface Transcript: a retry refused at the ceiling.

   ```transcript
   $ roundfix deliver retry 0301-example
   stdout:
   stderr:
   Preflight failed

   Reason:
     retry Delivery Queue item "0301-example": queue token ceiling 5000000 was reached with 5639755 tokens; start a new queue with roundfix deliver start
   ...
   exit: 2
   ```

7. Surface Transcript: one `usage` record of the Run Event Stream.

   ```transcript
   $ roundfix events run_20261001T120000Z_0123456789abcdef --filter usage
   stdout:
   {"schema":"roundfix-events/v1","run_id":"run_20261001T120000Z_0123456789abcdef","category":"usage","time":"2026-10-01T12:04:00Z","cursor":41,"work_item":"task_01","summary":"task_01 used 5639755 tokens (request-sum)","scope_kind":"task","scope_id":"task_01","runtime":"codex","model":"gpt-6.1-sol","reasoning_effort":"high","token_basis":"request-sum","tokens":5639755}
   stderr:
   exit: 0
   ```

Surface Transcript 6 uses the existing failure path of `deliver retry`,
whose standard error continues after the reason exactly as it does today.

### API Contracts

1. API Contract: `roundfix runs show <run-id> [--json]` — read-only; opens the
   Run Database with the reader. Text is Surface Transcripts 1 and 2. JSON has
   schema `roundfix/runs-show/v1` and the fields `run_id`, `kind`,
   `spec_slug`, `state`, `scopes` and `total`. Each scope carries
   `scope_kind`, `scope_id`, `selections`, `prompts`, `reported_prompts`,
   `tokens` (null when none reported), `bases`, the five split fields (each
   the sum over the scope's prompts when every reported prompt carried it,
   else null), `costs` (a list of `currency` and `amount`), `sessions` and
   `cost_sessions`; `total` carries the same fields without `scope_kind`,
   `scope_id` and `selections`. Exit `0` when the Run exists, `2` for a usage
   error or an unknown Run.
2. API Contract: Implement Run summary — one `Tokens:` line after the outcome
   line on every outcome, Surface Transcript 3. Exit codes are unchanged.
3. API Contract: `roundfix deliver status` and `roundfix deliver start` — the
   limits line ends `tokens <n>` or `tokens none`; `deliver status` adds the
   `Usage:` line, Surface Transcripts 4 and 5.
4. API Contract: `roundfix deliver start --max-tokens <n>` — `n` is an integer
   of at least 1; otherwise a usage error with exit `2` that records nothing.
   The value is stored with the queue and printed on the limits line.
5. API Contract: blocker `queue-token-ceiling` and its retry refusal, Surface
   Transcripts 4 and 6.
6. API Contract: `roundfix events` — category `usage`, on by default and
   accepted by `--filter`, Surface Transcript 7. Schema stays
   `roundfix-events/v1`, because the category and fields are additions.
7. API Contract: Run Database schema 22 — the two tables and the column in
   Data Models. An older binary refuses the database as for every schema
   change.

## Coverage Map

- Goal 1 → Reading the reports; The counting rule; Recording; API Contract 7.
- Goal 2 → Surfaces; API Contracts 1, 2, 3 and 6.
- Goal 3 → The ceiling; API Contracts 4 and 5.
- Goal 4 → The counting rule; Surfaces (`formatTokenTotals`).
- User Story 1 → `runs show`; API Contracts 1 and 2.
- User Story 2 → Recording; API Contract 6.
- User Story 3 → Surfaces; API Contract 3.
- User Story 4 → The ceiling; API Contracts 4 and 5.
- User Story 5 → Testing Approach 6.
- Core Feature 1 → Reading the reports.
- Core Feature 2 → The counting rule.
- Core Feature 3 → The counting rule; Surfaces.
- Core Feature 4 → Data Models (spend rule); Surfaces.
- Core Feature 5 → Data Models; Recording.
- Core Feature 6 → Recording; API Contract 6.
- Core Feature 7 → API Contract 1.
- Core Feature 8 → API Contract 2.
- Core Feature 9 → API Contract 3.
- Core Feature 10 → The ceiling; API Contracts 4 and 5.
- Core Feature 11 → Testing Approach 6.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 4.
- Success Metric 3 → Testing Approach 4.
- Success Metric 4 → Testing Approaches 2 and 4.
- Success Metric 5 → Testing Approach 5.
- Success Metric 6 → Testing Approach 3.

## Integration Points

- **acpx and the ACP adapters.** Roundfix reads the lines `acpx --format json`
  prints; nothing is sent to an adapter. Tests feed recorded lines through the
  existing fake acpx seam and never start an adapter.
- **The Run Database.** One migration, two tables and one column. Readers use
  `OpenReader`; only the Daemon and the Delivery Queue write.
- **The Delivery Engine.** The ceiling reads the store; no Run is signalled.
- **A metering gateway.** None. The user guide describes it as an operator
  option.

## Testing Approach

1. **Counting.** New `internal/agent/usage_test.go` requires the three
   outcomes of the counting rule on two recorded fixtures, new
   `internal/agent/testdata/usage/codex-turn.ndjson` (three `usage_update`
   lines with `used` 26830, 36330 and 150849, then a result whose `usage`
   total is 150849 and a later acpx result without `usage`) and
   `internal/agent/testdata/usage/claude-turn.ndjson` (two streaming readings,
   one reading with `cost` `{"amount":3.5526034999999996,"currency":"USD"}`,
   then a result with `usage` total 4019995 and its split). It also requires
   the runtime guard, an unreported prompt, a malformed payload ignored, and
   usage returned on a stopped prompt. Session IDs in the fixtures are
   replaced with `sess_fixture`.
2. **Storage.** New `internal/store/token_usage_test.go` requires the
   migration from schema 21, the row and its Run Event in one write, the
   sums per scope and per Run, the spend rule including a reset, the queue
   links written by both item updates and removed with the queue, and that
   retention leaves both tables untouched. `TestDurableTableLifecyclePolicyCoversEveryTable`
   covers the policy rows.
3. **Recording.** New `internal/daemon/token_usage_test.go` requires one row
   per prompt with the owner's scope and candidate, a row for a failed and a
   stopped prompt, the `session` scope without an owner, and a failing
   appender leaving the prompt's result, the Task status and the Run outcome
   unchanged, with the warning line.
4. **Surfaces.** New `internal/cli/runs_show_test.go`,
   `internal/cli/implement_tokens_test.go` and
   `internal/cli/deliver_usage_test.go` reproduce Surface Transcripts 1, 2, 3,
   4 and 5 and the JSON contract, and require that a Run with only
   unreported prompts never prints `0 tokens` on any of the three surfaces.
   New `internal/runevent/stream_usage_test.go` reproduces Surface
   Transcript 7 and the default filter. Existing tests that compare the
   whole Implement stdout or the limits line are updated.
5. **Ceiling.** New `internal/delivery/token_ceiling_test.go` requires a park
   without a worktree at and above the ceiling, a start below it, an item past
   `queued` continuing, a refused retry, no effect with no ceiling, and the
   Pending Question. New `internal/cli/deliver_token_ceiling_test.go` requires
   the flag, its refusal below 1 and the limits line.
6. **Documents.** Each Task's Verification requires its phrases in the user
   guide, and task_04's requires them in the Roundfix Skill and its mirror.

Each new gate is proved to fail: the Task that adds it records, in its
Result, the sabotage it applied and the failing test name, then restores the
code.

## Build Order

1. Turn Usage, the counting rule and the stream reading, task_01 (depends on:
   none).
2. Schema 22, the store API, the queue links, recording and the `usage`
   stream category, task_02 (depends on: 1). Recording reads
   `ExecuteResult.Usage`.
3. `runs show`, the Implement Run summary line, the `deliver` limits and
   usage lines and the user guide, task_03 (depends on: 2). The surfaces
   read the store API.
4. `--max-tokens`, the ceiling in the Delivery Engine and the Roundfix Skill,
   task_04 (depends on: 3). It shares `internal/cli/deliver.go`,
   `internal/cli/cli.go` and the deliver guide with task_03.
5. Terminal QA, task_05 (depends on: 1, 2, 3, 4).

This Spec is delivered after Spec 0194. task_03 and task_04 edit the
per-command guide and skill files it creates.

## Risks & Considerations

- **Schema number.** Another Spec of the same wave may take schema 22 first.
  task_02 then takes the next number; the rule is "one more than the version
  on its starting commit".
- **An adapter changes its reporting.** A Codex report that grows to the turn
  is caught by the guard. Any other change shows up as a changed number, not
  as an error; the basis on every row lets a reader find it.
- **acpx output.** Roundfix reads what `acpx --format json` prints. If a later
  acpx stopped printing `usage_update`, Codex prompts would be recorded as
  `turn` with the last request's total. The recorded basis makes that visible.
- **Write volume.** One row and one event per prompt, a few per Task. The
  event takes the ordinary batched path of ADR-0098.
- **Changed output.** Scripts that read the limits line or the whole Implement
  stdout see the declared breaks.

## Decisions

- Tokens are counted by the report's scope, and spend is only what the
  adapter reported. See ADR-0198.
- The ceiling acts at item start and retry, never inside a Run. See ADR-0199.
- Usage lives in its own table, and totals are summed on read.
- The lineage knowledge is one boolean on the existing adapter contract.
