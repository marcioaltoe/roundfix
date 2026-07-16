---
spec: 0035-agent-selection-profiles
prd: _prd.md
created: 2026-07-16
---

# Agent Selection Profiles — Technical Spec

## Executive Summary

Replace Run-wide runtime defaults with category-resolved Agent Selection Profiles while retaining the existing exact-selection ACP boundary. Configuration becomes an atomic profile of one Preferred Selection plus a required ordered Fallback Chain. The Implement Command validates every Task Type, resolves the categories used by the requested Task Graph and optional QA phase, and probes every distinct preferred and fallback tuple through disposable Agent Sessions before it creates a Run. Review commands resolve the `review` profile through the same path.

Each Task and QA action receives its own Agent Session so mixed Task Graphs can use different runtimes safely. A selection-start failure emits a structured notification and automatically advances to the next preflight-proven fallback before the first prompt. Once Agent work starts, no fallback is attempted over potentially modified state. A new `profiles` support command owns source-aware inspection, interactive or file-driven configuration, validation, and the static top-five recommendation display; skills only author and validate Task Type.

## System Architecture

The change spans these existing boundaries:

- `internal/agent` — official Model Catalog entries, immutable Agent Selection values, disposable proof, and live session creation for a selected runtime.
- `internal/config` — built-in profiles, strict YAML decoding, atomic profile precedence, optional Task Type inheritance, and backward-compatible migration from `runtimes.*` defaults.
- `internal/spec` — the closed Task Type set and early task-file validation.
- `internal/cli` — the `profiles` support command, existing explicit override semantics, category resolution, preflight ordering, stable text/JSON output, and caller-visible fallback notifications.
- `internal/run` / `internal/taskrun` — one Agent Session per Task, a separate QA Agent Session, and a review-selected session.
- `internal/store` — durable selection-attempt records and Run Event payloads.
- `.agents/skills/write-tasks` — authoring-time Task Type rules and template validation only; it has no ranking or model configuration responsibility.

No runtime-owned configuration is read or edited. All validation goes through the same pinned acpx adapter and authentication path that the live Run will use.

## Implementation Design

### Domain types

Keep the exact selection independent from routing policy:

```go
type AgentSelection struct {
    Runtime         agent.Kind
    Model           string
    ReasoningEffort string
}

type AgentSelectionProfile struct {
    Preferred AgentSelection
    Fallbacks []AgentSelection
}

type WorkCategory string

const (
    CategoryGeneral  WorkCategory = "general"
    CategoryBackend  WorkCategory = "backend"
    CategoryFrontend WorkCategory = "frontend"
    CategoryData     WorkCategory = "data"
    CategoryInfra    WorkCategory = "infra"
    CategoryDocs     WorkCategory = "docs"
    CategoryTest     WorkCategory = "test"
    CategoryChore    WorkCategory = "chore"
    CategoryQA       WorkCategory = "qa"
    CategoryReview   WorkCategory = "review"
)
```

`TaskType` is a closed type in `internal/spec` containing the seven Task values. Its parser rejects empty and unknown input with an error that includes the task file, invalid value, allowed values, and the instruction to update task frontmatter. Task parsing may continue to decode YAML before this check, but the Implement Command MUST validate the complete graph before profile resolution or any mutating Run step.

`AgentSelectionProfile.Validate()` enforces:

- complete runtime, model, and explicit reasoning field for the Preferred Selection;
- one or more complete Fallback Selections;
- no exact duplicate tuple in the profile;
- supported runtime name;
- official built-in ids for built-in entries, while explicit user custom model strings retain the existing forward-compatible ACP behavior;
- no empty profile overlay.

An empty reasoning string remains an explicit model-managed value where a runtime supports it. A missing YAML key is different from an explicitly empty string and is invalid inside a configured selection.

### Configuration schema and precedence

Add a top-level `profiles` section to User Config and Project Config:

```yaml
profiles:
  general:
    preferred:
      runtime: codex
      model: gpt-5.6-sol
      reasoning_effort: high
    fallbacks:
      - runtime: codex
        model: gpt-5.6-terra
        reasoning_effort: max
  backend:
    preferred:
      runtime: codex
      model: gpt-5.6-sol
      reasoning_effort: high
    fallbacks:
      - runtime: codex
        model: gpt-5.6-terra
        reasoning_effort: max
  frontend:
    preferred:
      runtime: claude
      model: claude-fable-5
      reasoning_effort: medium
    fallbacks:
      - runtime: codex
        model: gpt-5.6-sol
        reasoning_effort: high
  qa:
    preferred:
      runtime: codex
      model: gpt-5.6-sol
      reasoning_effort: high
    fallbacks:
      - runtime: codex
        model: gpt-5.6-terra
        reasoning_effort: max
  review:
    preferred:
      runtime: codex
      model: gpt-5.6-sol
      reasoning_effort: high
    fallbacks:
      - runtime: codex
        model: gpt-5.6-terra
        reasoning_effort: max
```

`general`, `backend`, `frontend`, `qa`, and `review` always resolve through built-ins when absent. `data`, `infra`, `docs`, `test`, and `chore` resolve to the effective `general` profile when absent. An optional profile that is present must be complete.

Precedence is:

1. invocation override of the Preferred Selection for this Run;
2. complete Project Config profile;
3. complete User Config profile;
4. complete built-in profile.

At levels 2–4, the first present profile wins as one object. The resolver never overlays `preferred`, fallback entries, or tuple fields across scopes. Invocation overrides preserve the resolved category's Fallback Chain, producing a complete one-Run profile before validation and proof. For a mixed Task Graph, the invocation override replaces the Preferred Selection of every relevant Task category and QA category; the CLI warns in stderr and JSON metadata that one explicit override is being applied across categories.

The legacy `defaults.agent` and `runtimes.<runtime>.{model,reasoning_effort}` keys remain readable for one compatibility window. When no `profiles` section exists at that scope, the resolver converts the legacy selected runtime tuple into the Preferred Selection for all categories and obtains built-in fallbacks. `roundfix profiles configure` writes the new schema only. If both schemas occur at the same scope, strict validation fails with a migration command instead of guessing precedence. The user guide and generated config are updated to the new schema.

### Official Model Catalog and adapter proof

Normalize built-in catalog entries to official identifiers:

- Codex: `gpt-5.6-sol`, `gpt-5.6-terra`, `gpt-5.6-luna` plus retained supported official older ids.
- Claude: `claude-fable-5` plus retained supported official ids such as `claude-opus-4-8`.

The existing catalog already contains the three GPT-5.6 ids; implementation updates labels, order, documentation, and tests rather than adding duplicates. The existing Claude `fable` and `opus` aliases are not rendered as official built-ins. If the installed adapter accepts an alias supplied explicitly by a user, it remains a custom value and is recorded verbatim.

ADR-0040 documented a pinned adapter that rejected non-empty reasoning for GPT-5.6. Implementation MUST first prove the exact new built-in tuples through the currently pinned adapter. If it cannot advertise and apply them, update the pinned acpx adapter through the repository's normal dependency/version path and test both configuration keys. Do not add a static exception that skips reasoning while displaying `high` or `max`; displayed and persisted values must be the values actually applied.

### Profile CLI

Add a support command with three deterministic actions:

```text
roundfix profiles show [--category <category>] [--json]
roundfix profiles configure --scope user|project [--file <path>] [--dry-run] [--json]
roundfix profiles validate [--category <category>] [--json]
```

`show` is read-only. It renders the effective source (`built-in`, `user`, `project`, or `invocation` where applicable), Preferred Selection, ordered fallbacks, inheritance source, and exactly five recommendations for each requested category. Default text output goes to stdout; warnings go to stderr. JSON uses stable arrays and never relies on map iteration order.

`configure` without `--file` uses Interactive Input. The user selects a category, then the Preferred Selection and at least one fallback, seeing the category's top-five recommendations beside the choice. The final screen shows the complete profile and target scope before the existing confirmation boundary writes atomically. `--file` is the non-interactive path: it reads a strict profile fragment, validates it completely, and writes only after success. `--dry-run` prints the normalized result without writing. `--json` reports `changed`, `scope`, `path`, normalized profiles, and validation errors without mixing diagnostics into stdout.

`validate` is read-only and proves all required effective profiles by default, or one requested category. It deduplicates identical Agent Selections but reports the proof under every category that references it. Exit `0` means every requested tuple was applied successfully in a disposable session; invalid configuration or failed runtime proof uses the repository's existing usage/readiness exit contracts and names a concrete remediation.

Unknown category, missing value, incomplete profile, an empty fallback list, or an invalid `--scope` is a usage error. `configure` never removes runtime-owned settings or credentials.

### Recommendation data

Store a versioned built-in recommendation dataset in Go data, not skill text or downloaded runtime state. Each row contains:

```go
type ModelRecommendation struct {
    Category         WorkCategory `json:"category"`
    Rank             int          `json:"rank"`
    Selection        AgentSelection `json:"selection"`
    Benchmark        string       `json:"benchmark"`
    ResultPercent    float64      `json:"result_percent"`
    AverageCostUSD   float64      `json:"average_cost_usd"`
    SourceAsOf       string       `json:"source_as_of"`
    Rationale        string       `json:"rationale"`
    CategorySpecific bool         `json:"category_specific"`
}
```

The initial dataset comes from `references/model-ranking.md` and the supplied snapshots dated 2026-07-16. Every list contains exactly five unique Agent Models, not five effort variants of one model. `CategorySpecific` is false for the initial release because the sources measure general coding-agent performance. CLI copy states that the category order is an informed recommendation, not empirical task-specific proof. A future ranking update is a code/docs change reviewed and released with its source date; it never changes existing User Config or Project Config.

### Operational preflight

The operational sequence for an Implement Command is:

1. locate and parse the Spec Task Graph and requested task files;
2. validate graph topology, task status, and every Task Type;
3. determine the Task Type categories in the executable graph and add `qa` only when QA was requested;
4. resolve complete effective profiles and apply invocation Preferred Selection overrides;
5. produce a stable sequence ordered by category, then preferred before fallback index; deduplicate exact tuples for execution while retaining category references;
6. open a disposable Agent Session for each distinct tuple, apply runtime/model/reasoning, and close it;
7. on any failure, close every opened disposable session, print the category and failed tuple with recovery guidance, and exit without creating a Run, branch, worktree, or durable selection attempt;
8. create the Run only after every relevant tuple passes.

Review fetch-only operations remain Agent-free. Resolve and watch operations use only the `review` profile. The `profiles validate` support command can probe categories that are not part of an operational Run.

Preflight is deterministic and sequential. This avoids authentication bursts, stable-order ambiguity, and new goroutine ownership solely to reduce a setup delay. Exact duplicate tuples are probed once.

### Agent Session lifecycle and fallback activation

Spec execution changes from one Run-wide Agent Session to:

- one session per Task, opened after its Task Worktree and bootstrap are ready and closed when that Task settles;
- one QA session opened after all Tasks settle and closed with the QA action;
- one review session selected from the `review` profile for review work.

For each session owner, activation follows:

1. attempt to create and prepare the Preferred Selection;
2. if preparation fails before the first prompt is sent, append a selection-attempt failure record and publish `agent_selection_fallback`;
3. render the same notification in live stderr/TUI and make it available to the Supervisor through the Run Event Stream;
4. only after the event is committed, create and prepare the next preflight-proven fallback;
5. repeat in configured order if another selection-start failure occurs;
6. if the chain is exhausted, fail the Work Item or action with every attempted tuple and recovery action.

The notification payload includes `category`, Work Item/action identity, failed selection, fallback selection, fallback index, normalized reason code, human-readable reason, and `automatic: true`. Cross-runtime activation uses a new runner/session factory keyed by the fallback's runtime; it does not mutate configuration.

The owner marks `agent_work_started` immediately before sending the first prompt. After that marker, prompt failure, tool failure, verification failure, cancellation, rate limit, or session loss follows existing Task/Run failure semantics. No replacement session is started because state may already have changed and replay would not be idempotent.

### Persistence and Run Events

Schema version 9 adds an append-only table:

```sql
CREATE TABLE run_agent_selections (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    run_id TEXT NOT NULL,
    scope_kind TEXT NOT NULL,
    scope_id TEXT NOT NULL,
    category TEXT NOT NULL,
    profile_source TEXT NOT NULL,
    attempt INTEGER NOT NULL,
    selection_role TEXT NOT NULL,
    runtime TEXT NOT NULL,
    model TEXT NOT NULL,
    reasoning_effort TEXT NOT NULL,
    status TEXT NOT NULL,
    reason_code TEXT NOT NULL DEFAULT '',
    reason TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    FOREIGN KEY (run_id) REFERENCES runs(id)
);
```

`scope_kind` is `task`, `qa`, or `review`; `scope_id` is the Task id, stable QA action id, or review Batch/Run id. `selection_role` is `preferred` or `fallback`; `status` is `attempting`, `active`, `failed`, or `closed`. The store validates monotonic attempt numbers per scope. It stores no prompt content, credentials, or runtime-owned config.

Existing `runs.agent`, `runs.model`, and `runs.reasoning_effort` remain non-null compatibility summaries: spec Runs store the effective `general` Preferred Selection, and review Runs store the effective `review` Preferred Selection. New UI and JSON detail reads `run_agent_selections` for actual per-scope selections. This avoids rewriting old Run rows while removing the false implication that one tuple drove every Task.

Add structured Run Events for selection attempt, fallback notification, selection active, and selection exhausted. The Run Event Stream remains the Supervisor contract. Caller-visible text is derived from the same payload so TUI, stderr, replay, and machine JSON cannot disagree.

### `write-tasks` contract

Update `.agents/skills/write-tasks/SKILL.md` and `references/task-template.md` so:

- every approved breakdown row and every task frontmatter has exactly one Task Type;
- valid values are `backend`, `frontend`, `data`, `infra`, `docs`, `test`, and `chore`;
- the type is selected from the dominant implementation outcome, not from the model the author prefers;
- cross-surface work uses its primary delivered behavior, while genuinely separable outcomes are split into vertical Tasks;
- graph validation parses all task frontmatter, rejects missing/unknown types, and confirms the human-readable projection matches each task file;
- no skill contains recommendations, runtime ids, model ids, or profile configuration logic.

## Cross-Spec dependencies

- Spec 0032 — Deterministic Agent Session cancellation is a hard implementation dependency. Its deterministic lifecycle boundary must pass before this Spec changes Agent Session ownership and fallback startup.
- Spec 0033 — Console Log tool-summary deduplication is an integration-order dependency. It must merge first because both Specs change Agent Run Event presentation and `internal/agent` / `internal/cli` wiring, but profile behavior does not consume deduplication behavior.
- Spec 0034 — Release Plan is an integration-order dependency. It must merge first because both Specs extend public CLI dispatch, help, documentation, and the synchronized Roundfix Skill; Agent Selection Profiles do not consume Release Plan behavior.

Implementation of Spec 0035 starts from a base containing completed Specs 0032–0034.

## API Contracts

### JSON profile output

`roundfix profiles show --json` returns one object:

```json
{
  "schema": "roundfix/profiles/v1",
  "profiles": [
    {
      "category": "backend",
      "source": "project",
      "inherited_from": "",
      "preferred": {
        "runtime": "codex",
        "model": "gpt-5.6-sol",
        "reasoning_effort": "high"
      },
      "fallbacks": [
        {
          "runtime": "codex",
          "model": "gpt-5.6-terra",
          "reasoning_effort": "max"
        }
      ],
      "recommendations": []
    }
  ]
}
```

The requested categories are ordered `general`, Task Types in their canonical order, `qa`, `review`. Fallback order and recommendation rank are stable. `inherited_from` is `general` only for an absent optional Task Type profile.

### Error and exit behavior

- Configuration syntax, unknown Task Type, or incomplete profile: usage/config error; stderr identifies the exact path and allowed values.
- Failed disposable proof: readiness error; stderr names runtime, model, reasoning effort, affected categories, adapter error, and a profiles validation/configuration next action.
- Fallback exhaustion after Run creation: existing Unresolved or Failed outcome according to the owning work path, with selection-attempt evidence.
- JSON commands return one JSON document on stdout; diagnostics remain structured in that document or on stderr according to existing command conventions.
- Read-only `show`, `validate`, and `--dry-run` never write config, Run state, sessions beyond disposable proof, or recommendation data.

## Coverage Map

- Goals 1–3 and Stories 1–3 → atomic profile resolver, relevant-category preflight, per-scope Agent Sessions, and notification-before-activation ordering.
- Goals 4 and Story 4 → official Model Catalog entries, versioned five-entry recommendation dataset, and `profiles show/configure` surfaces.
- Goals 5–7 and Stories 5–7 → Task Type parser, write-tasks contract, `profiles validate`, durable selection attempts, and deterministic JSON.
- Core Features 8–10 → disposable session integration tests, fallback lifecycle tests, store migration, and Run Event ordering assertions.

## Testing Approach

Every suite names its invariant and stays at the lowest owning layer:

- **Configuration unit tests (`internal/config`)** — table-driven precedence over built-in/User/Project profiles; atomic replacement; explicit empty reasoning versus missing key; optional inheritance; legacy migration; duplicate and empty fallback rejection. Negative companions cover partial overlays and same-scope legacy/new conflicts.
- **Task parser tests (`internal/spec`)** — all seven accepted values plus missing, mixed-case, whitespace, and unknown values. The observable contract is the returned task/error, not parser internals.
- **CLI tests (`internal/cli`)** — buffer-captured `Run(args, &stdout, &stderr) int` tests for text/JSON ordering, exactly five recommendations, no mutation from `show`/`validate --dry-run`, and actionable usage errors.
- **Agent integration tests (`internal/agent`)** — helper-process acpx sessions prove the exact official model/reasoning command sequence, session closure on success and error, cross-runtime factory selection, and no skipped reasoning assignment while reporting a non-empty effort.
- **Operational integration tests** — a fake Agent boundary with real config/spec/store files proves invalid Task Type blocks before probe and Run creation; all distinct preferred/fallback tuples are probed; one failed fallback proof blocks the Run; duplicate tuples are probed once in stable order.
- **Fallback lifecycle tests** — event/store fakes record that the caller-visible and structured notification precede fallback session creation and first prompt. A negative companion proves no fallback after `agent_work_started`. Cross-runtime fallback is exercised.
- **Persistence tests (`internal/store`)** — schema 8→9 migration, clean schema 9 creation, append/order constraints, and round-trip of empty reasoning and failure reason.
- **Macro QA** — build the binary in a temporary repository with fake pinned adapters; configure profiles from a file; show and validate them; run a two-type Task Graph plus QA; force one pre-prompt preferred failure; assert continued fallback execution, selection records, and Run Event Stream order. A second flow supplies an invalid Task Type and proves zero Run rows and Agent invocations.

Tests use real files, SQLite, CLI parsing, and helper processes at their owned boundaries. They do not assert values written into the same fake, introduce production-only hooks, depend on wall-clock sleeps, or weaken existing failure contracts. The benchmark recommendation dataset is a product contract, so exact snapshot-like fixture comparison is permitted; session behavior is asserted through explicit events and state.

Required verification:

```text
rtk gofmt -w <changed-go-files>
rtk go test ./internal/config ./internal/spec ./internal/agent ./internal/store ./internal/cli ./internal/run ./internal/taskrun
rtk go test -race ./...
rtk go run -buildvcs=false ./cmd/roundfix profiles show --category backend --json
rtk go run -buildvcs=false ./cmd/roundfix skills check
rtk make verify
```

The implementation task files must replace package placeholders with the smallest effect-proving commands for each slice.

## Build Order

1. Canonical Task Type and Agent Selection/Profile domain values, official Model Catalog normalization, and configuration schema/resolver (no dependencies).
2. `write-tasks` authoring and template contract plus Roundfix early Task Type validation (depends on: 1).
3. Versioned top-five recommendation dataset and read-only `profiles show` text/JSON surface (depends on: 1).
4. Atomic `profiles configure` and disposable `profiles validate` CLI flows, including legacy migration diagnostics (depends on: 1, 3).
5. Relevant-category operational preflight that proves every preferred/fallback tuple before Run creation (depends on: 1, 2, 4).
6. Per-Task, QA, and review Agent Session ownership with notification-before-automatic-fallback and cross-runtime session factories (depends on: 5).
7. Run Database selection-attempt migration, Run Events, Live Run View/Run Event Stream projection, and compatibility summaries (depends on: 6).
8. User guide, generated config, Roundfix Skill synchronization, integration QA, race gate, and full verification (depends on: 2–7).

## Risks & Considerations

- **Adapter drift.** Official model ids can exist while a pinned ACP adapter advertises different configuration. Disposable proof is authoritative; the pinned adapter must be upgraded or the built-in profile cannot ship.
- **Preflight latency and cost.** Proving every fallback adds setup work. Exact tuple deduplication limits repeats; correctness and guaranteed fallback readiness take priority over parallel probes.
- **Mixed-session lifecycle.** Per-Task sessions increase session count and cleanup paths. Every session has one Task/action owner and closes on success, failure, cancellation, and panic-free early return.
- **Unsafe replay.** Automatic fallback after work begins could duplicate edits or commands. The hard `agent_work_started` boundary forbids it.
- **Benchmark overreach.** The supplied rankings measure general coding performance, not each category. CLI copy and JSON metadata keep that limitation visible, and routing never consumes the rank.
- **Configuration duplication.** Required profiles duplicate the general default. Atomic replacement makes the effective policy auditable and avoids hidden inheritance; only optional profiles inherit implicitly.
- **Legacy ambiguity.** Mixing old runtime defaults and new profiles at one scope is rejected with a migration command instead of choosing silently.
- **Run summary compatibility.** Existing Run-level columns cannot describe mixed profiles. They remain compatibility summaries while all new behavior and UI use per-scope records.
- **Skill ownership.** Only the repo-owned `write-tasks` skill is updated for Task Type. Model configuration remains CLI-owned, and the shipped Roundfix Skill must be synchronized with public CLI behavior before delivery.

## Decisions

- Profiles replace lower-precedence profiles atomically and include a non-empty ordered Fallback Chain.
- Built-in defaults use official identifiers and exact efforts that must pass the pinned adapter's disposable proof.
- Operational preflight probes relevant profiles only; the explicit validation command can probe all required profiles.
- Preflight is stable-order, sequential, and deduplicated by exact Agent Selection.
- Automatic fallback is notification-first and pre-prompt only.
- Each Task and QA action owns its Agent Session; review uses the dedicated `review` profile.
- Ranking is versioned product data shown by the CLI, limited to five entries, and never routing input.
- Skills author Task Type but never own model selection or ranking.
