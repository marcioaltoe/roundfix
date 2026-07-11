---
spec: 0025-optional-reasoning-effort
prd: _prd.md
created: 2026-07-11
---

# Optional Reasoning Effort — Technical Spec

## Executive Summary

Agent selection currently hard-requires a non-empty Default Reasoning Effort
at four independent chokepoints and always issues the acpx reasoning set
call, which makes every model that manages reasoning itself (codex gpt-5.6
family) fail Preflight Validation. The fix relaxes the emptiness checks and
skips the set call when the effort is empty, keeping the explicit-value
contract untouched. The primary trade-off: an empty effort no longer fails
fast, so a user who *forgot* to configure an effort now runs model-managed
instead of getting an error. We accept this because the effective selection
is persisted and printed on every Run header, and because ADR-0040 makes
"empty means model-managed" the documented meaning, mirroring the existing
empty-means-derive keys (`defaults.artifact_dir`, `watch.push_remote`,
`worktree.bootstrap`). No new packages, seams, or config keys.

## System Architecture

All changes extend existing modules; nothing new is created.

- `internal/cli/selection.go` — `ResolveSelection` stops treating an empty
  resolved effort as a missing selection (model stays required).
- `internal/agent/acpx_runner.go` — `validateRuntimeSelection` drops the
  non-empty effort requirement; `applyDisposableSelection` and
  `applySelection` skip the `set <key>` acpx call when the trimmed effort is
  empty. The config-key resolvability check runs only for non-empty values.
- `internal/cli/cli.go` — the explicit-empty flag guard
  (`validateExplicitSelectionFlags`) accepts `--reasoning-effort ""` as the
  model-managed override; the three Run headers render the model-managed
  state.
- `internal/tui/tui.go` — Interactive Input represents "model-managed" for
  the Default Reasoning Effort field instead of raising
  `requiredSelectionError`; `validateCollectedSelections` drops the effort
  emptiness check.
- `internal/store` — unchanged: `runs.reasoning_effort` already defaults to
  `''`; an empty persisted value now means model-managed (legacy rows render
  the same way, which is truthful for both).
- Config layer — unchanged mechanics: the `*string` overlay already
  distinguishes an absent key (inherit) from a present-empty key (blank the
  inherited value), so `.roundfixrc.yml` can override the User Config's
  `xhigh` with `reasoning_effort: ""`.

## Implementation Design

### Interfaces

No signature changes. Behavior deltas on existing contracts:

```go
// internal/cli/selection.go
// ResolveSelection: effort no longer appended to `missing`;
// an empty resolved ReasoningEffort flows through.

// internal/agent/acpx_runner.go
func validateRuntimeSelection(runtime RuntimeSpec) error {
    // Model required (unchanged). Empty ReasoningEffort is valid.
    // acpxReasoningEffortConfigKey checked only when effort != "".
}

func (runner ACPXRunner) applyDisposableSelection(...) error {
    // after model set succeeds:
    if strings.TrimSpace(runtime.ReasoningEffort) == "" {
        return nil // model-managed: no reasoning option is assigned
    }
    // existing set path unchanged
}
// applySelection gets the same guard.
```

### Data Models

None. `runs.reasoning_effort TEXT NOT NULL DEFAULT ''` already stores the
empty state; no migration.

### API Contracts

- `--reasoning-effort ""` (explicit empty) becomes a valid override meaning
  model-managed; today's exit-2 `emptySelectionFlagError` for this flag is
  removed. `--model ""` stays invalid.
- Run header line renders `Default Reasoning Effort: model-managed` when the
  effective effort is empty (resolve `cli.go:1576`, watch `cli.go:1726`,
  implement `implement.go:433`). Attach and the Live Run View keep their
  existing legacy-empty rendering.
- `SelectionPreflightError` copy for a rejected non-empty value gains the
  remediation naming the model-managed option (for example
  `recovery: ... or set runtimes.<runtime>.reasoning_effort "" when the
  model manages reasoning`).
- Stdout/stderr split, exit codes, and every other command contract are
  unchanged.

## Coverage Map

- Goal "empty effort selects without a reasoning option" → ResolveSelection,
  validateRuntimeSelection, applyDisposableSelection/applySelection guards.
- Goal "explicit values keep failing loudly" → unchanged set path plus the
  extended SelectionPreflightError recovery copy.
- Goal "gpt-5.6-sol as this repo's default" → `.roundfixrc.yml` flip.
- Core Feature 3 (honest surfaces) → Run headers, flag guard, Interactive
  Input field, persisted Run row.
- Core Feature 4 (guidance) → README, CONTEXT.md glossary, ADR-0040, both
  Roundfix skill copies (`.agents/skills/roundfix/SKILL.md` +
  `skills/roundfix/SKILL.md`, validated by `roundfix skills check`).

## Integration Points

acpx/codex-acp only, through the existing runner: the reasoning
`set <key> <value>` call becomes conditional. No new adapter behavior is
assumed; a model-managed model is driven by simply not sending the option.

## Testing Approach

Existing seams only: the fake acpx runner in `internal/agent`
(argument-recording command runner), buffer-captured CLI runs, and synchronous
`model.Update` TUI tests.

- Unit (`internal/agent`): empty-effort selection skips the reasoning set on
  both disposable and live paths (assert recorded args contain no
  `set reasoning_effort`); non-empty rejection still yields
  `SelectionPreflightError` with the new recovery copy; update the
  empty-effort case in `TestACPXRunRequiresConcreteSelection` to expect
  success-shaped behavior.
- CLI (`internal/cli`): doctor/setup accept a configured empty effort
  (rework `TestRunDoctorRejectsMissingConfiguredAgentSelection` and the setup
  twin); header renders `model-managed`; explicit `--reasoning-effort ""`
  passes validation and persists `''` on the Run row.
- TUI (`internal/tui`): Interactive Input accepts the model-managed state
  without `requiredSelectionError`; existing catalog tests keep passing.
- Docs: `roundfix skills check` inside `make verify` gates the two skill
  copies.

## Build Order

1. Core selection change: relax `ResolveSelection` and
   `validateRuntimeSelection`, add the empty-effort skip to
   `applyDisposableSelection` and `applySelection`, extend the
   `SelectionPreflightError` recovery copy, and update the
   `internal/agent` + selection unit tests.
2. CLI surfaces (depends on: 1): explicit-empty `--reasoning-effort` flag
   acceptance, `model-managed` header rendering across resolve/watch/
   implement, doctor/setup acceptance of empty configured effort, and the
   affected `internal/cli` tests.
3. Interactive Input (depends on: 1): model-managed representation for the
   Default Reasoning Effort field, drop the effort emptiness checks in
   `CollectInput`/`validateCollectedSelections`, and the affected
   `internal/tui` tests.
4. Configuration and guidance (depends on: 1, 2, 3): flip `.roundfixrc.yml`
   to `model: gpt-5.6-sol` + `reasoning_effort: ""` with an explanatory
   comment, update CONTEXT.md's Default Reasoning Effort entry, README
   selection/config sections, and both Roundfix skill copies.

## Risks & Considerations

- A forgotten effort now runs model-managed instead of failing; mitigated by
  the persisted effective selection, the header line, and documentation.
- `.roundfixrc.yml` diverges from the 0024 branch's committed
  `gpt-5.6-sol`/`high` value; whichever PR merges second resolves a trivial
  two-line conflict, and `high` on gpt-5.6-sol fails preflight with the new
  actionable recovery text until amended.
- The TUI codex effort choices (`low/medium/high/xhigh`) remain static and
  model-agnostic; picking one for a gpt-5.6 model fails preflight with the
  actionable message. Per-model choice filtering stays out of scope
  (Non-Goals; ADR-0039 rejects static catalog claims).

## Decisions

- Empty Default Reasoning Effort means the Agent Model manages reasoning;
  selection skips the set call. See ADR-0040.
- A rejected non-empty value keeps failing Preflight Validation without
  fallback. See ADR-0037, ADR-0039, ADR-0040.
- The empty string is the sentinel (no named `model-managed` config value):
  it matches the existing empty-means-derive config keys and the `*string`
  overlay's absent-vs-present-empty semantics.
- Built-in defaults are unchanged; only this repository's Project Config
  moves to gpt-5.6-sol.
