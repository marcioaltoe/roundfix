---
spec: 0218-the-jev-router-on-docs-and-chore-tasks
prd: _prd.md
created: 2026-10-01
---

# The Jev Router on docs and chore Tasks — Technical Spec

## Executive Summary

The Jev Router becomes one OpenCode Agent Selection,
`roundfix-openrouter/typesafe/jev-router` with an empty effort. Configuration
accepts it only from Project Config. Every session of that selection gets a
fixed inline OpenCode provider whose key is the placeholder
`{env:ROUNDFIX_OPENROUTER_API_KEY}`, so the key reaches OpenCode from the
environment and never through Roundfix. A new package, `internal/jevrouter`,
reads the month's Jev spend from the Judge Log and the key's own
`usage_monthly`, and appends one `router-prompt` line per routed prompt. The
Daemon asks it before and after each routed prompt: a refusal before work
moves the Fallback Chain on, and after work fails the Work Item (ADR-0218).
A final live Task replays four archived `docs` and `chore` Tasks on the router
and on the current default.

The trade-off is a ceiling that reads a remote number before each prompt
over one that trusts only local records. The key's `usage_monthly` counts
every OpenRouter call on Roundfix's key, including a call whose line was
never written, at the cost of two small HTTPS reads per routed prompt and a
fail-closed refusal when OpenRouter cannot answer.

## Project Constraints

- Identifier strategy: not applicable — the routed selection is a fixed
  OpenCode model value, and a router line reuses the Judge Log's fields with
  the fixed judgment kind `router-prompt`. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — `internal/jevrouter` sends
  `GET https://openrouter.ai/api/v1/key` with
  `Authorization: Bearer <ROUNDFIX_OPENROUTER_API_KEY>` and no other header,
  10 seconds per attempt, no retry, honoring the standard proxy variables.
  It reads `data.usage_monthly`, a non-negative number. The key is read only
  from the process environment, checked for presence by `internal/agent`, and
  sent only in that header; it is never written to an argument, a file, the
  inline configuration, a log, a Run Event or a Judge Log line. No production
  file reads `OPENROUTER_API_KEY`. OpenCode sends the prompts to
  `https://openrouter.ai/api/v1` with the same key, through the placeholder.
  Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0218 binds the selection, the
  Project Config rule, the inline provider, the gate and the router line;
  ADR-0201's ceiling, read from the judge's question file, and ADR-0200's
  judge stay as Spec 0205 builds them; ADR-0050 and ADR-0114 decide refusal
  before and after work; ADR-0017, ADR-0037, ADR-0049, ADR-0105, ADR-0125,
  ADR-0187, ADR-0189, ADR-0193, ADR-0198 and ADR-0211 bind the runtime, the
  proof, the fixtures, the skill and the prerequisites; ADR-0108, ADR-0180,
  ADR-0181 and ADR-0199 stay unchanged; ADR-0208 explains the split; ADR-0184
  states the surface as a transcript. The gate is bound by ADR-0080,
  ADR-0088, ADR-0091, ADR-0104, ADR-0155, ADR-0156 and ADR-0167, and
  ADR-0093, ADR-0117, ADR-0168, ADR-0176 and ADR-0183 check consistency.
  ADR-0166, ADR-0178 and ADR-0182 govern each Task commit. ADR-0069,
  ADR-0096, ADR-0097, ADR-0192, ADR-0194, ADR-0195, ADR-0209 and ADR-0210 do
  not apply, as the PRD records. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30 ("considere autorizado a ajustar todas as skills se
  necessário") and the cycle authorization of 2026-10-01 ("Tudo, de A a F"),
  recorded in [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/runtime.md`. Sanctioned regeneration:
  `make skills-sync`. The Makefile, CI workflows, lint configuration,
  `.roundfixrc.yml`, `internal/cli/cli_test.go` and `go.mod` stay untouched.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## System Architecture

| Concern | Owner | File |
| --- | --- | --- |
| The routed model, its provider and its key name | `JevRouterModel`, `JevRouterKeyEnv`, `jevRouterProviderConfig`, `IsJevRouterSelection` | new `internal/agent/jev_router.go` |
| The session environment and the key check | `codexEnvForSession`, the one source of per-session overrides | `internal/agent/acpx_runner.go` |
| The Project Config rule and the effort rule | `validateProfiles`, `normalizeSelection`, `ResolveProfile` | `internal/config/profiles.go` |
| The month's Jev spend, the key reader, the router line | `jevrouter.Spend`, `jevrouter.KeyUsage`, `jevrouter.Ledger` | new `internal/jevrouter/spend.go`, `internal/jevrouter/key.go`, `internal/jevrouter/ledger.go` |
| The gate before and the record after a routed prompt | `JevRouterGate` in `Dependencies`, `agentSessionOwner.runPrepared` | `internal/daemon/engine.go`, `internal/daemon/agent_session_owner.go` |
| Documentation | the configuration guide and the Roundfix Skill | `docs/user-guide/configuration.md`, `.agents/skills/roundfix/references/runtime.md` |
| The measurement | the protocol below | `docs/specs/0218-the-jev-router-on-docs-and-chore-tasks/measurement/jev-router.md` (new) |

## Implementation Design

### Interfaces

```go
// internal/agent/jev_router.go
const (
	JevRouterModel     = "roundfix-openrouter/typesafe/jev-router"
	JevRouterKeyEnv    = "ROUNDFIX_OPENROUTER_API_KEY"
	JevRouterKeyMissing = "jev_router_key_missing"
)

// IsJevRouterSelection reports whether runtime and model name the routed
// selection: runtime opencode (with or without -custom) and JevRouterModel.
func IsJevRouterSelection(runtime, model string) bool

// internal/jevrouter/spend.go
type Spend struct {
	TypeSafeLogged, OpenRouterLogged, KeyUsageMonthly, Total, Ceiling float64
}
// MonthSpend reads the current UTC month's Judge Log and the key's usage.
// Total = TypeSafeLogged + max(OpenRouterLogged, KeyUsageMonthly).
func MonthSpend(ctx context.Context, deps Deps, now time.Time) (Spend, error)

// internal/jevrouter/key.go
// KeyUsage sends GET <endpoint>/key and returns data.usage_monthly.
func KeyUsage(ctx context.Context, client *http.Client, endpoint, key string) (float64, error)

// internal/jevrouter/ledger.go
type PromptRecord struct {
	RunID, Spec, ScopeKind, ScopeID, Category, Repository string
	Attempt                                             int
	UsageBefore, UsageAfter                             float64
	UsageAfterErr                                       error
	Latency                                             time.Duration
	InputTokens, OutputTokens                           int64
	Failed                                              bool
}
// Append writes one router-prompt line to the month's Judge Log.
func (ledger Ledger) Append(record PromptRecord, now time.Time) error

// internal/daemon/engine.go
type JevRouterGate interface {
	// Before returns the key usage it read, or a classified refusal.
	Before(ctx context.Context) (usageBefore float64, err error)
	After(ctx context.Context, record jevrouter.PromptRecord) error
}
```

`jevrouter.Deps` carries the environment, Roundfix Home, an `*http.Client`,
the key endpoint base (`https://openrouter.ai/api/v1` in production), and the
ceiling, which production takes from the judge package's loaded questions
(`monthly_ceiling_usd`, Spec 0205), never from a second literal.

### The routed selection

`internal/config` reads `agent.IsJevRouterSelection`, so the model value is
written once. `internal/agent` imports nothing from `internal/config`, so the
import adds no cycle.

`normalizeSelection` refuses a routed selection whose trimmed effort is not
empty with `<path>.reasoning_effort must be "" for the Jev Router; the router
chooses the effort`. `validateProfiles` refuses any profile entry whose
`Source` is not `project` and whose preferred selection or a fallback is
routed, with `<path> names the Jev Router (roundfix-openrouter/typesafe/jev-router),
which only Project Config may select`. `ResolveProfile` refuses a routed
`preferredOverride` with `the Jev Router cannot be a one-Run override; name
it in Project Config`. Built-in profiles and the legacy `runtimes:` section
cannot name it, because neither is Project Config.

### The inline provider and the key

`jevRouterProviderConfig` is this fixed JSON, compact, in this key order:

```json
{"$schema":"https://opencode.ai/config.json","provider":{"roundfix-openrouter":{"npm":"@ai-sdk/openai-compatible","name":"Roundfix OpenRouter","options":{"baseURL":"https://openrouter.ai/api/v1","apiKey":"{env:ROUNDFIX_OPENROUTER_API_KEY}"},"models":{"typesafe/jev-router":{"name":"Jev Router"}}}}}
```

`codexEnvForSession`, which every proof, fallback probe, sealed session,
session setup and prompt already calls for its per-session overrides, gains a
branch for a routed runtime before its Codex logic. When `JevRouterKeyEnv` is
empty or absent in the runner's base environment it returns
`&SelectionFailureError{Runtime: "opencode", Reason: JevRouterKeyMissing + ": ROUNDFIX_OPENROUTER_API_KEY is not set"}`.
Otherwise it returns `OPENCODE_CONFIG_CONTENT=<jevRouterProviderConfig>`,
which replaces any inherited `OPENCODE_CONFIG_CONTENT` for that session only.
Proof therefore sets the routed model through the same ACP config option as
any OpenCode model, and refuses without a key before any session opens.

### The month's Jev spend

`MonthSpend` reads `<home>/.roundfix/judge/<YYYY-MM>.jsonl` for the UTC month
of `now` through the judge package's Judge Log reader (Spec 0205), and sums
`cost_usd` by `transport`: `typesafe` into `TypeSafeLogged`, and `openrouter`
into `OpenRouterLogged`, router lines included. A missing file is `0`. It
then calls `KeyUsage` with `JevRouterKeyEnv`. An unreadable log, a key
endpoint that fails, times out, answers a non-`200`, or a body without a
non-negative `data.usage_monthly`, is an error. `Total` is
`TypeSafeLogged + max(OpenRouterLogged, KeyUsageMonthly)`.

### The gate and the record

`Dependencies.JevRouter` defaults, when nil, to a gate built from the process
environment, Roundfix Home and the judge's ceiling. In
`agentSessionOwner.runPrepared`, when the active runtime is routed:

1. **Before.** `Before` computes `MonthSpend`. On an error it returns
   `&SelectionFailureError{Runtime: "opencode", Reason: "jev_spend_unreadable: <error>"}`;
   at `Total >= Ceiling` it returns one with
   `Reason: "jev_ceiling_reached: month's Jev spend US$<total> of US$<ceiling>"`,
   both with four decimals. The prompt is not sent. The existing `Run` loop
   then falls back after notification when work has not begun, and returns
   the error, failing the Work Item, when it has.
2. **The prompt** runs as today, and its usage is recorded as today
   (ADR-0198).
3. **After.** The owner calls `After` with a `PromptRecord` carrying the
   usage `Before` returned; `After` reads the key's usage again, filling
   `UsageAfter` or `UsageAfterErr`, and `Append` writes one Judge Log line:
   `"schema":"roundfix/judge-log/v1"`, `"judgment":"router-prompt"`,
   `"spec"` the Run's Spec slug or `""`, `"artifact":""`, `"line":0`,
   `"target":"<run-id> <scope-kind> <scope-id>"`, `"question_id"` the Agent
   Work Category, `"transport":"openrouter"`, `"response_id":""`,
   `"provider":""`, `"requested_model"` `JevRouterModel`, `"model":""`,
   `"answer"`, `"probabilities"`, `"confidence"` and `"noul"` null,
   `"latency_ms"`, `"input_tokens"` and `"output_tokens"` from the prompt,
   `"cost_usd"` `max(0, UsageAfter - UsageBefore)`, `"cost_source":"reported"`,
   `"status":0`, `"attempts"` the attempt, `"error"` `""` or
   `key usage unreadable after prompt: <error>` with `cost_usd` `0`, and
   `"outcome":"skipped"` when the prompt failed or `"clear"` when it
   returned. The outcome values are the Judge Log's own, so its reader
   accepts the line unchanged. A failed append is reported on the Run's
   progress stream and does not change the prompt's result; the next
   `Before` still reads the key's usage.

No line holds the key, the prompt or the answer.

### The measurement protocol

The measurement Task runs only on the maintainer's machine, with
`NODE_OPTIONS` unset, `bin/roundfix` rebuilt, and
`ROUNDFIX_OPENROUTER_API_KEY` in its environment. Without the key it records
`jev_router_key_missing` in its result and stops.

1. **Replays.** Four archived Tasks, each on the router and on the current
   built-in default of its category, in a fresh scratch clone of this
   repository per replay: `0210-evidence-snapshots-that-stay-small/task_02`
   and `0194-a-skill-and-a-command-guide-read-one-command-at-a-time/task_04`
   (docs), `0200-a-skill-snapshot-that-matches-its-upstream/task_04` and
   `0195-owned-skills-and-a-release-step-that-follow-the-bundle/task_06`
   (chore); reserves, in order,
   `0202-a-qa-gate-that-reruns-only-stale-rows/task_04` and
   `0191-claims-with-receipts-and-contracts-as-they-ship/task_03`.
2. **Pre-state.** A replay starts from the Spec's squash merge commit with
   the Task's declared files restored to that commit's first parent, and the
   Spec directory restored under `docs/specs/` with that Task `pending`,
   every other Task removed from the graph, and no QA Task. A Task whose
   Verification passes on that state is replaced by the next reserve.
3. **Selections.** The routed replay's scratch `.roundfixrc.yml` names the
   router as the category's preferred selection and the current default as
   its fallback, and runs `bin/roundfix implement --spec <slug> --no-input`.
   The default replay runs the same command with `--agent`, `--model` and
   `--reasoning-effort` naming the current default. The routed replay runs
   first; when its gate refuses, the replay is recorded as stopped by the
   ceiling and no further routed replay starts.
4. **Record.** Per replay: selection, wall time, prompts and Verification
   repairs from `bin/roundfix runs show <run>`, outcome, whether a fallback
   activated and why, tokens, and for a routed replay the sum of its Run's
   `router-prompt` costs. Write
   `docs/specs/0218-the-jev-router-on-docs-and-chore-tasks/measurement/jev-router.md`
   with a `## Measured` table and a `## Reading` section that names cost per
   settled Task on the router and whether tool calls worked through it.
5. **Follow-up.** When the router settled every Task the default settled,
   with no more Verification repairs, the Task also mints a dated Backlog
   Entry under `docs/backlog/` proposing the router for `docs` and `chore`,
   which the Daemon records as an undeclared path; otherwise the reading says
   why not. It changes no profile and no Project Config of this repository.

Scratch clones live under the system temporary directory, are clones of this
repository only, and are removed after the record is written.

### Data Models

- The routed selection is an `AgentSelection` with `runtime: opencode`,
  `model: roundfix-openrouter/typesafe/jev-router`, `reasoning_effort: ""`.
- The router line is a Judge Log line, schema `roundfix/judge-log/v1`, with
  the field values of "The gate and the record". No new file or schema.

### API Contracts

1. API Contract: Project Config may name the routed selection as a preferred
   or fallback selection; User Config naming it, or a non-empty effort on it,
   is a configuration load error, exit `2` on every command that loads
   configuration; a one-Run override naming it is refused with exit `2`.
2. API Contract: `GET https://openrouter.ai/api/v1/key` with the bearer key;
   `200` with a non-negative `data.usage_monthly` is read, and anything else
   is `jev_spend_unreadable`.
3. API Contract: a refused routed prompt is a failed Agent Selection with
   reason `jev_router_key_missing`, `jev_spend_unreadable` or
   `jev_ceiling_reached`, published through the existing fallback
   notification before work, and as the Work Item's failure after it.
4. API Contract: the `router-prompt` Judge Log line of "The gate and the
   record", one per routed prompt.

### Surface Transcripts

The fixture for Transcript 1 is an empty Git repository and a Home whose
`/home/maintainer/.roundfix/config.yml` names the router as the `docs`
preferred selection and `codex gpt-6.1-sol high` as its fallback.

1. Surface Transcript: the router in User Config.

   ```transcript
   $ roundfix profiles show --category docs
   stdout:
   stderr:
   roundfix: profiles failed: parse config "/home/maintainer/.roundfix/config.yml": profiles.docs.preferred names the Jev Router (roundfix-openrouter/typesafe/jev-router), which only Project Config may select
   Run 'roundfix profiles --help' for usage.
   exit: 2
   ```

## Coverage Map

- Goal 1 → The routed selection; The inline provider and the key; API Contract 1.
- Goal 2 → The routed selection; The inline provider and the key; Surface Transcript 1.
- Goal 3 → The month's Jev spend; The gate and the record; API Contracts 2-4.
- Goal 4 → The measurement protocol.
- User Story 1 → The routed selection; API Contract 1.
- User Story 2 → The routed selection; Surface Transcript 1.
- User Story 3 → The month's Jev spend; The gate and the record.
- User Story 4 → The measurement protocol.
- Core Feature 1 → The routed selection.
- Core Feature 2 → The routed selection; Surface Transcript 1.
- Core Feature 3 → The inline provider and the key.
- Core Feature 4 → The month's Jev spend; The gate and the record; API Contract 3.
- Core Feature 5 → The gate and the record; API Contract 4.
- Core Feature 6 → The measurement protocol.
- Core Feature 7 → Testing Approach 4.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 1.
- Success Metric 3 → Testing Approach 3.
- Success Metric 4 → Testing Approach 2.
- Success Metric 5 → Testing Approach 5.

## Integration Points

- **OpenCode** `1.18.34` or newer through acpx's `opencode` agent; it reads
  `OPENCODE_CONFIG_CONTENT` and substitutes `{env:...}`.
- **OpenRouter** chat completions at `https://openrouter.ai/api/v1`, reached
  by OpenCode, and the key endpoint `GET /api/v1/key`, reached by Roundfix,
  both on `ROUNDFIX_OPENROUTER_API_KEY`.
- **The judge (Spec 0205)**: its Judge Log reader, its ceiling, and the log
  file this Spec appends to.

## Testing Approach

No test reaches OpenRouter or the network: the key endpoint is an
`httptest.Server`, every fake OpenCode and acpx is the compiled test binary
(ADR-0125), and every Home is a temporary directory.

1. **The selection and the key.** `TestJevRouterLoadsFromProjectConfig`,
   `TestJevRouterRefusedFromUserConfig` (Transcript 1's message),
   `TestJevRouterRefusesAReasoningEffort`,
   `TestJevRouterRefusedAsOneRunOverride` in `internal/config`.
   `TestJevRouterSessionCarriesThePlaceholderConfig` sets a sentinel key and
   asserts that the acpx environment holds `OPENCODE_CONFIG_CONTENT` equal to
   the fixed JSON and that the sentinel is in no argument and no
   configuration value; `TestJevRouterRefusedWithoutTheKey` returns
   `jev_router_key_missing` before any acpx call;
   `TestOtherOpenCodeModelsGetNoInlineConfig` keeps today's environment. A
   sweep finds no `"OPENROUTER_API_KEY"` literal in production Go files.
2. **The spend.** `TestMonthSpendAddsTypeSafeToTheLargerOpenRouterFigure`,
   `TestMonthSpendFailsOnAnUnreadableKeyEndpoint`,
   `TestKeyUsageSendsOnlyTheBearerHeader`,
   `TestKeyUsageRefusesAMalformedAnswer`, and
   `TestRouterLineIsReadByTheJudgeLog`, which appends a line and reads it
   back with the judge package's reader, counting its cost.
3. **The gate.** In `internal/daemon` with a fake gate and a fake runner:
   `TestJevRouterCeilingFallsBackBeforeWork` (the fallback starts, the
   routed prompt never ran), `TestJevRouterCeilingFailsTheTaskAfterWork`,
   `TestJevRouterUnreadableSpendFallsBack`, and
   `TestJevRouterPromptAppendsOneLine` with the usage change as its cost and
   no sentinel key in the line, and `TestNonRoutedPromptCallsNoGate`.
4. **Documentation.** Phrase checks over the configuration guide and the
   skill reference, the skill's recorded version and `make skills-sync-check`.
5. **The measurement.** The record exists with `## Measured` and
   `## Reading` and names the router.

## Build Order

1. The routed selection, the inline provider, the key check and the
   configuration guide, task_01 (depends on: none).
2. `internal/jevrouter`: the spend, the key reader and the router line,
   task_02 (depends on: none; it requires Spec 0205 in the tree).
3. The Daemon gate and record, and the Roundfix Skill, task_03 (depends on:
   1 and 2).
4. The measurement, task_04 (depends on: 1, 2 and 3).
5. Terminal QA, task_05 (depends on: 1, 2, 3, 4).

## Risks & Considerations

- **Prerequisites.** Task 02 and Task 03 need Spec 0205's judge package; run
  before it lands they fail to compile, which is a queue-order error, not a
  defect.
- **Usage lag.** OpenRouter may report a prompt's cost some seconds late, so
  one prompt's line can undercount and the next one's overcount; the sum
  stays right, and `max(OpenRouterLogged, KeyUsageMonthly)` keeps the gate
  conservative.
- **A prompt in flight.** The gate runs between prompts, so one long prompt
  can carry the month past US$5 by its own cost.
- **Prompt cache.** The router may change the model between turns of one
  session, which discards the provider's prompt cache; a cheaper model can
  still cost more per Task. The measurement compares whole Tasks for that
  reason.
- **Tool calls through an OpenAI-compatible provider.** The router lists
  `tools` and `tool_choice`, but a routed model may handle them differently
  turn to turn; the measurement records it.
- **A user's own inline configuration.** A maintainer who exports
  `OPENCODE_CONFIG_CONTENT` loses it for routed sessions only; the guide says
  so.

## Decisions

- Reach the router through OpenCode's inline configuration, not a file
  Roundfix writes, so no file ever needs protecting.
- Gate on `max(logged, key usage)` plus the TypeSafe log, not the log alone,
  so a missing line cannot hide spend.
- Record router prompts in the Judge Log with its own outcome vocabulary, so
  Spec 0205's reader and ceiling count them without change.
- Default the Daemon's gate inside `internal/daemon`, so no CLI wiring
  changes.
