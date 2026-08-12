---
spec: 0035-agent-selection-profiles
status: active
created: 2026-07-16
surfaces: [cli, backend, data, docs, frontend]
---

# Agent Selection Profiles

Roundfix currently selects one ACP Runtime, Agent Model, and reasoning effort for an entire Run. Its fallback is discovered dynamically from one runtime's Model Catalog and requires confirmation after failure. That contract cannot guarantee that backend, frontend, QA, and review work use the user's intended models, cannot safely route mixed Task Graphs, and leaves the chosen fallback unknown until the Run has already been delayed. Agent Selection Profiles make the complete preferred-and-fallback policy explicit, proven before a Run, observable per Work Item, and configurable from the Roundfix CLI.

The supplied benchmark snapshots and current vendor documentation support a recommendation ranking, but they do not provide category-specific proof for every Task Type. Rankings therefore help users configure profiles and never become an automatic router.

## Goals

- Make `general`, `backend`, `frontend`, `qa`, and `review` explicit Agent Selection Profiles with one Preferred Selection and at least one configured Fallback Selection.
- Let optional `data`, `infra`, `docs`, `test`, and `chore` profiles inherit `general` when they are absent.
- Keep `gpt-5.6-sol` at `high` as the built-in general, backend, QA, and review Preferred Selection, and use `claude-fable-5` at `medium` as the built-in frontend Preferred Selection.
- Use official model identifiers throughout the built-in Model Catalog, profile configuration, CLI output, persistence, and Agent Session preparation.
- Prove every relevant Preferred Selection and configured Fallback Selection through disposable Agent Sessions before creating a Run.
- Notify the user and Supervisor before automatically activating the next proven fallback, including a cross-runtime fallback when configured.
- Require and validate Task Type before an Implement Command can create a Run or start Agent work.
- Provide deterministic CLI surfaces to inspect, configure, and validate profiles and to show only the five best advisory recommendations for each Agent Work Category.
- Persist the effective Agent Selection and fallback history for every Task, QA action, and review action.

## User Stories

1. As a developer starting mixed backend and frontend Tasks, I want each Task routed by its declared Task Type, so the intended runtime and model are used without answering the same questions again.
2. As a developer relying on a preferred model, I want its fallback configured and proven before the Run, so a selection failure does not stop work or surprise me with an unknown replacement.
3. As a Supervisor monitoring a detached Run, I want a structured fallback notification before the replacement Agent Session starts, so I can tell the user exactly what changed and why while work continues automatically.
4. As a developer configuring Roundfix, I want five result-and-cost recommendations per category, so I can make an informed choice without Roundfix silently routing from benchmark scores.
5. As a maintainer authoring a Spec Task Graph, I want `write-tasks` to require a valid Task Type and Roundfix to enforce it, so every Task has deterministic profile routing.
6. As a developer running QA or pull-request review, I want dedicated `qa` and `review` profiles, so verification and review do not accidentally inherit the model used by the last implementation Task.
7. As an automation author, I want JSON inspection and validation output with stable fields and exit codes, so profile readiness can be checked without Interactive Input.

## Core Features

1. **Atomic profile configuration.** An Agent Selection Profile contains one complete Preferred Selection and a non-empty ordered Fallback Chain. A profile defined at a higher-precedence scope replaces the whole lower-precedence profile; partial profile merging is invalid.
2. **Required and optional categories.** Effective configuration always resolves `general`, `backend`, `frontend`, `qa`, and `review`. Missing optional Task Type profiles (`data`, `infra`, `docs`, `test`, and `chore`) resolve to `general` without creating duplicate stored configuration.
3. **Built-in profile defaults.** Built-ins use:
   - `general`, `backend`, `qa`, and `review`: preferred `codex / gpt-5.6-sol / high`, fallback `codex / gpt-5.6-terra / max`;
   - `frontend`: preferred `claude / claude-fable-5 / medium`, fallback `codex / gpt-5.6-sol / high`.
4. **Official Model Catalog entries.** Codex exposes `gpt-5.6-sol`, `gpt-5.6-terra`, and `gpt-5.6-luna`; Claude exposes `claude-fable-5`. Existing short Claude aliases are migrated or retained only as documented custom values, never rendered as the official built-in identifier.
5. **Profile CLI.** A support command shows effective source-aware profiles, configures complete profiles at User Config or Project Config scope, validates them through the installed runtimes, and renders deterministic JSON when requested.
6. **Advisory rankings.** The CLI shows exactly five versioned recommendations for a selected Agent Work Category. Each recommendation identifies runtime, official model id, reasoning effort, benchmark result, cost signal, source date, and caveat. Ranking data is shipped with Roundfix and never fetched or applied automatically.
7. **Task Type contract.** Every task file requires one of `backend`, `frontend`, `data`, `infra`, `docs`, `test`, or `chore`. The Implement Command rejects missing or unknown values before model probing, Run persistence, worktree creation, or Agent work. `write-tasks` and its template produce and validate that contract.
8. **Selection preflight.** After loading and validating the requested Work Items, Roundfix resolves the relevant categories, deduplicates exact Agent Selections, and probes every preferred and fallback tuple in stable order through disposable sessions. Any failed proof blocks Run creation and names the category, tuple, and next action.
9. **Automatic fallback boundary.** A live selection failure before the first Agent prompt triggers a structured notification followed by the next proven fallback. Roundfix never changes selection after prompt or tool work has begun and never treats general Agent failure as a model-selection fallback condition.
10. **Per-work observability.** Run persistence and Run Events record category, profile source, preferred or fallback role, attempt order, exact selection, activation result, and failure reason for each Task, QA action, or review action.
11. **Backward-compatible explicit overrides.** Existing invocation-level agent/model/reasoning flags remain explicit one-Run Preferred Selection overrides. They do not remove the resolved category's required Fallback Chain, and the resulting complete effective profile is subject to the same preflight.

## Non-Goals / Out of Scope

- Automatically routing work from benchmark rank, price, latency, or availability.
- Fetching live benchmark data or changing a profile when a leaderboard changes.
- Letting `write-tasks`, `qa-gate`, `review`, or another skill configure or display models.
- Inferring Task Type from task content at Run time.
- Switching models after an Agent prompt or tool action has begun.
- Falling back for test failures, verification failures, low-quality output, token usage, rate limits without a selection-start failure, or general Agent errors.
- Editing ACP Runtime-owned configuration, authentication, or credentials.
- Adding unsupported runtimes solely because a benchmark lists them.
- Claiming that general coding benchmarks are frontend-, QA-, or review-specific evaluations.
- Making `data`, `infra`, `docs`, `test`, or `chore` profiles mandatory.

## Success Metrics

- 100% of Implement Commands with a missing or invalid Task Type fail before any disposable model probe, Run row, Run Worktree, or Agent Session is created.
- 100% of relevant Preferred Selections and configured Fallback Selections are proven successfully before a Run is persisted; one failed proof blocks the Run.
- Every automatic fallback emits its notification Run Event and caller-visible message before the fallback Agent Session receives its first prompt.
- No test or production path activates a fallback after Agent work has begun.
- Optional Task Types resolve byte-for-byte to the effective `general` profile when no category-specific profile exists.
- CLI recommendation output contains exactly five deterministic entries per requested category in text and JSON modes and never mutates configuration.
- Every Task, QA action, and review action has a durable record of the effective selection and any fallback attempt.
- `write-tasks` rejects an approved breakdown or generated task file whose type is missing or outside the seven allowed Task Types.
- The full repository verification gate and race-enabled selection/session tests pass without retries.

## Decisions

- Profiles are explicit routing policy; rankings are advisory input only.
- A higher-precedence profile is an atomic replacement, not a field overlay.
- Project Config wins over User Config, which wins over built-ins; invocation flags override the Preferred Selection for the current Run only.
- Operational preflight proves only profiles relevant to the requested Run, while the profile validation command can prove all configured required profiles.
- All configured fallbacks are proven before Run creation; activation follows list order and may change ACP Runtime.
- Automatic activation is safe only before Agent work begins.
- Task Type is author-declared and validated, never guessed by Roundfix.

## Open Questions

None.
