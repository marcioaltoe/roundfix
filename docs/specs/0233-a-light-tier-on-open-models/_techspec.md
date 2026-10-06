---
spec: 0233-a-light-tier-on-open-models
prd: _prd.md
created: 2026-10-05
---

# A light tier on open models — Technical Spec

## Executive Summary

Dispatch today gives every Task the Agent Selection Profile of its Agent Work
Category, through the Agent Session owner in the daemon. This design adds one
decision before the owner is built: a pure tier rule over the Task file
(`complexity`, `type`, the declared files) in a new small package,
`internal/lighttier`, which also owns the Light Spend Log. For a `light` Task
the daemon derives a profile whose Preferred Selection is the first light model
on `opencode`, with no reasoning effort, and whose Fallback Chain is the other
light models followed by the category's own profile; the existing
fallback-before-work machinery then covers start failures, and escalation is
the owner jumping past the light candidates for the one Verification Feedback
turn. The ACPX runner gives a light session an inline OpenCode provider option
that names the key variable, and the CLI builds the tier's plan from User
Config and the Run's environment through one key helper. The judge gains a
`tasks` stage that suggests a tier. The primary trade-off is cost accounting:
Roundfix trusts OpenCode's reported session cost, already parsed into each
prompt's usage, instead of reading OpenRouter's key usage, accepting that an
unreported cost counts as zero (ADR-0238).

Measured on the starting tree (`9085540d`):

1. **An empty effort already sends no warm-up.** `PlanSelectionAssignment`
   gives an OpenCode selection with an empty effort the `runtime_managed`
   encoding, and `warmSessionForDeferredEffort` returns before the
   `Session setup.` prompt for every encoding except `runtime_deferred`. The
   light tier forces the empty effort; no warm-up change is needed.
2. **The warm-up's restriction does not reach OpenCode.** The inert prompt
   passes `--deny-all --allowed-tools ""`; acpx 0.19.4 forwards allowed tools
   only into Claude Code session metadata and the Qoder command line, which
   matches the measurement's 38-step `Session setup.` answer.
3. **OpenCode's cost is already captured.** The prompt stream parses
   `usage_update` cost readings into `TurnUsage.Cost`, the session's last
   cumulative reading, and the daemon stores it per prompt in the Run
   Database's token usage.
4. **`complexity` is parsed but unused.** The Task front matter is decoded
   with `complexity`, but the value never reaches `spec.Task`.
5. **The Run Database allows only `preferred` and `fallback` selection roles**
   and leaves `profile_source` free, so a derived profile with source
   `light-tier` needs no schema change.

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes; the new
  keys follow the dotted configuration names and the profile source
  `light-tier` follows the existing source words. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: applicable — a light session reaches OpenRouter
  through OpenCode with the operator's key. Roundfix sets
  `OPENCODE_CONFIG_CONTENT` for that session only, with
  `provider.openrouter.options.apiKey` holding `{env:<variable>}`, never the
  value, and removes `OPENROUTER_API_KEY` from the session's environment.
  Roundfix itself sends no request. The judge's recipients and keys are
  unchanged. No test or Verification reaches the network. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0238 governs the design: "A Task
  is `light` when its `complexity` is `low`, its `type` is not `qa`, and it
  declares no Governed Path to create, change or delete", and "An empty
  effort sends no warm-up". ADR-0235: "Other OpenRouter models stay
  selectable through OpenCode"; each light model id is checked by the
  existing subscription predicate. ADR-0108: "Roundfix therefore warms an
  OpenCode session"; unchanged for selections with an effort. ADR-0050: "once
  Agent work begins, Roundfix fails the Work Item instead of switching models
  over potentially modified state", kept for the Fallback Chain; ADR-0114: "A
  Fallback Selection may switch ACP Runtime automatically only while Agent
  work has not begun", which the derived profile relies on. ADR-0049: "each
  present profile replaces the complete lower-precedence profile"; configured
  profiles are untouched and the derived one exists only at dispatch.
  ADR-0231: "It must be a finite number greater than zero"; the new ceiling
  follows it. ADR-0002: "Roundfix uses YAML for User Config at", and ADR-0027:
  "truly unknown keys keep failing strict validation". ADR-0200: "Spec
  authoring gets an advisory judge that never gates", and ADR-0201: "Each
  request appends one line to a Judge Log in Roundfix Home"; the tier
  suggestion follows both. ADR-0187 and ADR-0189 govern the skill edits and
  their versions, and ADR-0233 regenerates the raised versions at merge.
  ADR-0184: "A TechSpec now declares numbered Surface Transcripts"; three are
  declared below. The authored QA gate follows ADR-0080: "QA verdicts
  distinguish environment-blocked rows", and ADR-0091: "required to be
  terminal and to depend on every leaf"; ADR-0096, ADR-0097 and ADR-0167 bind
  its machine stage, its row carry and its pre-PR Pull Request row, ADR-0104:
  "Every Spec therefore rests at least one named" acceptance row on outside
  evidence, ADR-0194, ADR-0195 and ADR-0210 bind what a QA row records, when
  it is observed again and its evidence snapshot, and ADR-0182 settles each
  Task on the facts its gate checks. ADR-0093, ADR-0117, ADR-0156, ADR-0168,
  ADR-0176 and ADR-0183 check this Spec's consistency by citation and receipt. ADR-0069 cites ADR-0050 but decides the Baseline semantic analysis, ADR-0180 cites ADR-0069 but decides the dated Recommended Profile, which this Spec leaves unchanged, ADR-0181 cites ADR-0180 but decides where a configuration is compared with that profile, ADR-0209 cites ADR-0200 but decides how the judge suggests sources and ADR-0208 cites ADR-0209 but decides how sources are grouped, ADR-0217 cites ADR-0049 but decides how a Cursor selection is named, and ADR-0229 cites ADR-0167 but decides how an operator archive resumes a park and ADR-0237 cites ADR-0229 but decides how a Delivery Retry records a merge made outside the queue; this Spec changes none of them, so none applies.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — task_01 edits the Roundfix Skill's
  canonical `SKILL.md`, `runtime` and `spec` references and the `SKILL.md`
  mirror, and the write-tasks skill and its mirror, all Governed Paths;
  express maintainer authorization: "considere autorizado a ajustar todas as
  skills se necessário" (2026-09-30) and "Concedo" (2026-10-05); bounded
  files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/runtime.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/write-tasks/SKILL.md`, `skills/roundfix/SKILL.md`,
  `skills/write-tasks/SKILL.md`. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0233-a-light-tier-on-open-models/_authorization.md`.

## System Architecture

| Component | Where | Change |
| --- | --- | --- |
| Task complexity | `Task` in `internal/spec/spec.go`; `ReloadTask` and the graph loader in `internal/spec/task.go` | `Complexity` carries the front-matter value, trimmed |
| Light tier keys | `internal/config/config.go`; new `internal/config/light_tier.go` | `OpenRouter` config with the two User Config keys, defaults, validation, Project Config warning; `OpenRouterImplementKey` |
| Tier rule and Light Spend Log | new `internal/lighttier` | `TierFor`, `Plan`, `SpendLog` |
| Light session | `RuntimeSpec` in `internal/agent/agent.go`; session environment in `internal/agent/acpx_runner.go` | `OpenRouterKeyVariable`; inline provider option and removed generic variable for every acpx command of that session |
| Dispatch | `taskAgentSessionOwner`, `agentSessionOwner` in `internal/daemon/agent_session_owner.go`; `TaskPlan` and the repair path in `internal/daemon/task_engine.go` | Derived profile, skip warnings, spend lines, escalation |
| Run wiring | `executeImplementCycle` and its caller in `internal/cli/implement.go` | Builds the `lighttier.Plan` from User Config, the Run's environment and Roundfix Home; off under a one-Run override |
| Advisory tier | `internal/judge` question catalog and planner; `internal/cli/spec_judge.go` | `tasks` stage, `model-tier` judgment |
| Record | configuration guide, `spec` command reference, model reference, Roundfix Skill `runtime` and `spec` references, write-tasks skill | Describe the tier, keys, warnings, log, escalation and stage |

`internal/lighttier` earns its place against two alternatives: the daemon
package would mix a file ledger into the engine, and the config package must
not write files. It imports `internal/spec` and `internal/speccheck` for the
Governed Path predicate and nothing from the daemon.

## Implementation Design

### Interfaces

```go
// internal/config/light_tier.go
const DefaultLightModel = "deepseek/deepseek-v4.1-flash"
const DefaultImplementMonthlyCeilingUSD = 10.0
type OpenRouter struct {
	LightModels                 []string // nil means the default list; empty means off
	ImplementMonthlyCeilingUSD  float64  // 0 means the default
}
// OpenRouterImplementKey is the only place that names the implementation
// key's variable. It returns that name and whether it holds a non-empty value.
func OpenRouterImplementKey(getenv func(string) string) (variable string, present bool)
```

```go
// internal/lighttier
type Tier string // "light" or "standard"
func TierFor(task spec.Task) Tier
type Plan struct {
	Models     []string // OpenRouter ids; empty means off
	CeilingUSD float64
	KeyVariable string; KeyPresent bool
	HomeDir, Repository string
}
func (plan Plan) Enabled() bool
func ReadMonth(home string, now time.Time) (float64, error)
func AppendSpend(home string, now time.Time, line SpendLine) error
```

### Invariants

```text
1. TierFor is light only when Complexity == "low", Type != qa, and no Context ref of kind creates, interface or deletes names a GovernedPath.
2. Plan.Enabled is false for a zero Plan, an empty model list, or a one-Run override; a zero Plan in TaskPlan keeps every existing test's dispatch.
3. A light model id is "<author>/<slug>" with no whitespace and passes CheckSubscriptionRule with runtime "opencode" and model "openrouter/" + id.
4. The derived profile is Preferred = models[0], Fallbacks = models[1:] then the category's Preferred and Fallbacks; every light selection has an empty effort; Source = "light-tier".
5. The skip order is: key absent, then log unreadable, then month at or above the ceiling; each skip uses the category's profile unchanged.
6. A light RuntimeSpec carries OpenRouterKeyVariable; only then does the runner set OPENCODE_CONFIG_CONTENT, merging provider.openrouter.options.apiKey = "{env:<variable>}" into an inherited JSON object, and drop OPENROUTER_API_KEY.
7. An inherited OPENCODE_CONFIG_CONTENT that is not a JSON object fails the light session before work, so the next candidate takes it.
8. After each prompt of a light session one Light Spend Log line records the increase of the session's reported USD cost, or 0 with cost_source "unreported".
9. Escalation happens at most once per Task, only after the first Verification of a Task whose active candidate is a light model, and starts the first non-light candidate in a new Agent Session.
10. No Light Spend Log line, Run Event, progress line or log holds the key value.
```

### Data Models

User Config gains:

```yaml
openrouter:
  light_models: [deepseek/deepseek-v4.1-flash]
  implement_monthly_ceiling_usd: 10
```

The Light Spend Log is `<home>/.roundfix/openrouter/implement/<YYYY-MM>.jsonl`
(UTC month), directory mode 0700 and file mode 0600 as the Judge Log, one JSON
object per line: `schema` (`roundfix/light-spend/v1`), `time`, `repository`,
`run_id`, `spec`, `task`, `session`, `model`, `cost_usd`, `cost_source`
(`opencode` or `unreported`). The month's spend is the sum of `cost_usd`; a
line that does not parse, or a negative cost, makes the month unreadable.
`TaskPlan` gains `LightTier lighttier.Plan`. No Run Database change.

### API Contracts

1. API Contract: User Config `openrouter.light_models` is a sequence of
   OpenRouter model ids; unset means `[deepseek/deepseek-v4.1-flash]` and `[]`
   turns the tier off. An entry that is not `<author>/<slug>` fails with
   `openrouter.light_models[<i>] "<id>" must be an OpenRouter model id <author>/<slug>`,
   and an entry the subscription rule refuses fails with that rule's message
   for field `openrouter.light_models[<i>]`.
2. API Contract: User Config `openrouter.implement_monthly_ceiling_usd` must
   be a finite number greater than zero, else loading fails with
   `openrouter.implement_monthly_ceiling_usd must be a finite number greater than 0`;
   unset means 10.
3. API Contract: either key in Project Config is removed before validation
   with the existing User Config-only warning naming the key.
4. API Contract: a skipped light Task prints to standard error
   `roundfix: warning: light tier skipped for Task <id>: <reason>; it runs on its <category> profile`,
   where `<reason>` is `<variable> is not set`,
   `the Light Spend Log could not be read: <error>` or
   `this month's light spend US$<spent, 4 decimals> reached the ceiling of US$<ceiling, 2 decimals>`,
   and publishes a Task Run Event with phase `light_tier_skipped` and
   `reason_code` `key_missing`, `spend_unreadable` or `ceiling_reached`.
5. API Contract: an escalation prints
   `Task <id> escalates from the light tier to its <category> profile after Verification failed.`
   and publishes a Task Run Event with phase `light_tier_escalated`, the
   light model and the failed commands.
6. API Contract: a light Task's Agent Selection attempts carry
   `profile_source` `light-tier`; its token usage rows carry the light
   runtime and model as today.
7. API Contract: `roundfix spec judge <slug> --stage tasks` requires
   `_tasks.md`, asks one `model-tier` Choice per non-QA Task file (`light`,
   `standard`, `heavy`), prints
   `suggested model-tier <task file>: <answer> at confidence <c>` for each,
   records a Judge Log line per request with judgment `model-tier`, and
   writes no Spec file. Without `--stage` the judge still judges the PRD and
   the TechSpec only. An unknown stage fails with
   `unsupported --stage "<value>"; use prd, techspec or tasks`.

### Surface Transcripts

1. Surface Transcript: `roundfix spec judge` refuses an unknown stage, in a
   disposable repository with a disposable Home.

   ```transcript
   $ roundfix spec judge 0300-example --stage qa
   stdout:
   stderr:
   roundfix: spec judge failed: unsupported --stage "qa"; use prd, techspec or tasks
   Run 'roundfix spec judge --help' for usage.
   exit: 2
   ```

2. Surface Transcript: `roundfix spec judge --stage tasks` with no key, on a
   Spec with two non-QA Tasks and a QA gate.

   ```transcript
   $ roundfix spec judge 0300-example --stage tasks
   stdout:
   Judge: skipped: ROUNDFIX_OPENROUTER_API_KEY is not set (nor ROUNDFIX_TYPESAFE_API_KEY); 2 judgment(s) not asked
   stderr:
   exit: 0
   ```

3. Surface Transcript: `roundfix implement` with the fake ACP adapter, a
   disposable Home with no OpenRouter key, and one `low` `docs` Task.

   ```transcript
   $ roundfix implement --spec 0300-example --no-input
   stdout:
   ...
   stderr:
   ...
   roundfix: warning: light tier skipped for Task task_01: ROUNDFIX_OPENROUTER_API_KEY is not set; it runs on its docs profile
   ...
   exit: 0
   ```

## Coverage Map

- Goal 1 → `TierFor`, derived profile, light session environment, Run wiring
- Goal 2 → `openrouter.implement_monthly_ceiling_usd`, `ReadMonth`, `AppendSpend`
- Goal 3 → skip order (Invariant 5), API Contract 4
- Goal 4 → escalation (Invariant 9), API Contract 5
- Goal 5 → `tasks` stage, `model-tier` judgment
- Story 1 → `TierFor`, derived profile, default light model
- Story 2 → ceiling key, Light Spend Log
- Story 3 → `openrouter.light_models`
- Story 4 → key helper, skip warning
- Story 5 → escalation
- Story 6 → `tasks` stage
- Core Feature 1 → `TierFor`, `Plan.Enabled`
- Core Feature 2 → derived profile, light model validation, empty effort
- Core Feature 3 → `OpenRouterImplementKey`, light session environment
- Core Feature 4 → ceiling key, `ReadMonth`, skip warning
- Core Feature 5 → `AppendSpend` after each light prompt
- Core Feature 6 → escalation, start-failure fallback
- Core Feature 7 → `tasks` stage, `model-tier` judgment
- Core Feature 8 → the record
- Success Metric 1 → `TierFor`, derived profile, light session test
- Success Metric 2 → skip order, start-failure fallback
- Success Metric 3 → light session environment, `AppendSpend`
- Success Metric 4 → escalation
- Success Metric 5 → User Config keys
- Success Metric 6 → `tasks` stage
- Success Metric 7 → QA live light Task

## Integration Points

- OpenCode 1.18.34 through acpx 0.19.4: the light session receives
  `--model openrouter/<id>` and the inline provider option; inline
  configuration overrides the global and project files.
- OpenRouter: reached only by OpenCode in a light session, on the operator's
  key; Roundfix sends no request.
- The judge's transports are unchanged.

## Testing Approach

- Config: table tests through `Load` with disposable Homes and repositories
  for both keys, defaults, refusals, Project Config warnings and the key
  helper over a fake `getenv`.
- `lighttier`: `TierFor` over Tasks loaded from disposable Spec files
  (governed `creates`, governed `instruction` only, `qa`, `medium`, missing
  complexity); Light Spend Log append, sum, missing file, malformed line,
  UTC month and file modes in a temporary Home.
- Runner: the fake acpx harness records the child environment and prompts: a
  light runtime gets the merged inline option with the variable reference and
  no generic variable, never the key value, and sends only the work prompt; a
  non-light runtime keeps its inherited configuration.
- Daemon: the existing Task engine harness with a fake runner and a temporary
  Home: dispatch order, each skip reason with its warning and Run Event, a
  start failure taken by the category profile, spend lines from reported
  costs, and one escalation with a new session.
- CLI: `implement` with the fake runner builds the plan from User Config and
  the Run's environment, and an `--agent` override turns it off; `spec judge
  --stage tasks` through the existing fake transport.

## Build Order

1. Describe the light tier, its keys, warnings, log, escalation and the
   judge's `tasks` stage in the guides, the Roundfix Skill and the write-tasks
   skill, raising both skill versions, so the CLI change ships with its guide.
2. Add the Task's complexity, the User Config keys, the key helper and the
   `lighttier` package with the tier rule and the Light Spend Log.
3. Give the runner its light session and the daemon its dispatch, skips,
   spend lines and escalation, and wire the plan from `implement` (depends
   on: 2).
4. Add the judge's `tasks` stage and `model-tier` suggestion (depends on: 1).
5. Final QA gate (depends on: 1, 2, 3, 4).

## Risks & Considerations

- A light session that OpenCode does not bill on Roundfix's key (an inline
  `{env:}` reference left unresolved, issue 13219) would spend elsewhere or
  fail; the measurement saw it resolve on 1.18.34 and the QA gate compares
  the log with the key's usage.
- Concurrent light Tasks can pass the ceiling by their own cost; OpenRouter's
  monthly limit on the key is the backstop.
- An unreported cost counts as zero; the warning-free path depends on
  OpenCode reporting cost, which the measurement observed for every session.
- Every `low` Task's prompt, the files its agent reads and its diagnostics
  reach OpenRouter and the model's provider; the guides say so.
- Spec 0234 will change `OpenRouterImplementKey`; every caller goes through
  it, so the change is one function.

## Vocabulary Contract

- emits: `internal/daemon/agent_session_owner.go`
  pattern: `light tier skipped for Task`
  documented-in: `docs/user-guide/configuration.md`
- emits: `internal/daemon/task_engine.go`
  pattern: `light_tier_escalated`
  documented-in: `.agents/skills/roundfix/references/runtime.md`
- emits: `internal/cli/spec_judge.go`
  pattern: `suggested model-tier`
  documented-in: `.agents/skills/roundfix/references/spec.md`

Two candidate glossary terms are used and left to the close-of-Spec glossary
check: **Light Tier** (the dispatch of a `light` Task on an open model) and
**Light Spend Log** (its monthly spend record in Roundfix Home), beside
**Agent Selection Profile**, **Preferred Selection**, **Fallback Chain**,
**Verification Feedback** and **Judge Log**.

## Decisions

- Tier from `complexity: low`, QA and Governed Paths excluded; judge advisory.
  See ADR-0238.
- A derived Agent Selection Profile with source `light-tier`, rather than a
  new selection role and a Run Database migration. See ADR-0238.
- Empty effort, so no warm-up; ADR-0108 unchanged. See ADR-0238.
- OpenCode's reported cost in a Light Spend Log, rather than OpenRouter key
  reads. See ADR-0238.
- Escalation keeps the light attempt's changes. See ADR-0238.
- Light models are not proved by Preflight; a start failure falls through
  the derived Fallback Chain.
