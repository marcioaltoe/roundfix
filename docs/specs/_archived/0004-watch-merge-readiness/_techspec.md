---
spec: 0004-watch-merge-readiness
prd: _prd.md
created: 2026-07-05
---

# Watch Merge Readiness — Technical Spec

## Executive Summary

Four changes to the review pipeline's edges, no new packages. The primary
trade-off: watch's terminal condition gains one more external status source
(the GitHub check on the head commit), accepted because the existing
`StatusSource` seam and bounds (Max Rounds, review timeout) already contain
exactly this kind of polling — the merge-readiness confirm phase is one more
state in the loop, not a new loop. Everything else is contract honesty:
loop ordering, stdout reports, and a display-only stderr filter.

## System Architecture

- `internal/watch` — loop reordering (immediate first status check, quiet
  skip when pre-settled) and the post-push confirm phase driven by a new
  `CheckFunc` seam beside the existing `StatusFunc`/`FetchFunc`/`ResolveFunc`
  adapters.
- `internal/cli` — wires the check adapter through the existing GitHub CLI
  integration (`gh` — same auth surface preflight already uses), prints the
  new watch/resolve stdout reports from data the cycle already returns, and
  adds `--no-agent-console` to the operational commands.
- `internal/runevent` — a filtering `Sink` decorator that drops agent-source
  events; used only for the stderr writer in non-TTY mode. The Journal sink
  is never filtered (ADR-0008/0009 untouched).
- No store, journal, config, or TUI changes.

## Implementation Design

### Interfaces

```go
// internal/watch — one new seam beside StatusFunc/FetchFunc/ResolveFunc
type HeadCheckState string // "pending" | "success" | "failure" | "missing"
type CheckFunc func(ctx context.Context, headSHA string) (HeadCheckState, error)

// Request gains: Check CheckFunc (nil => legacy behavior for tests only;
// the CLI always wires it).

// internal/runevent — display-only filter
func NewSourceFilterSink(next Sink, drop Source) Sink
```

### Watch loop changes

1. **Poll-first:** the settled-wait performs its first `StatusFunc` call
   before any sleep; `poll_interval` separates subsequent calls. If the first
   call reports settled, the quiet period is skipped entirely (the review
   settled before the Run began — there is nothing to let quiesce).
2. **Confirm phase (ADR-0019):** after a Final Push inside `--until-clean`,
   instead of ending Clean immediately, the loop polls `CheckFunc` with the
   pushed head SHA on `poll_interval`, bounded by the review timeout and the
   remaining Max Rounds budget:
   - `success` → end Clean.
   - `failure` or `pending`-then-settled with new Review Issues → next Round
     (existing fetch path picks up the new issues).
   - `missing` (no check reported for the Review Source) → end Clean with a
     stderr note — repos without the check keep the old semantics.
   - timeout/rounds exhausted → existing `TimedOut` / `MaxRoundsReached`.
3. Non-`--until-clean` watch and plain resolve are unchanged.

### Check adapter

`gh api repos/<owner>/<repo>/commits/<sha>/check-runs` (or `gh pr checks`
when simpler) filtered to the Review Source's app slug (CodeRabbit), mapped
to `HeadCheckState`; transient `gh` failures return an error the loop treats
like a status-poll failure (retry on next interval within the timeout).
Adapter lives beside the existing gh usage in the CLI wiring layer; the watch
package sees only `CheckFunc`.

### API Contracts

- **stdout (watch and resolve):** one line per Review Issue in Round/fetch
  order — `issue <id> <status> — <title>` with the final local status
  (`resolved|invalid|failed|duplicated|unresolved`) — then one outcome line:
  `<Outcome> after <N> Round(s): <X> resolved, <Y> invalid, <Z> failed, <W> unresolved.`
  Resolve prints the same shape with `1 Round(s)`. Nothing else on stdout.
- **`--no-agent-console`** on resolve, watch, and implement: stderr keeps
  Daemon-source events (and the existing header/progress lines) and drops
  Agent-source events. Display-only; Journal and cockpit unaffected. Rejected
  in interactive TTY mode (the cockpit already separates surfaces).
- Exit codes, flags otherwise, Round semantics, and Final Push mechanics are
  unchanged.

## Coverage Map

- Story 1 → watch loop poll-first + quiet skip (finding 18)
- Story 2 → CheckFunc seam + confirm phase + gh adapter (finding 20, ADR-0019)
- Story 3 → stdout reports for watch/resolve (finding 19)
- Story 4 → NewSourceFilterSink + `--no-agent-console` (finding 4)

## Integration Points

- **GitHub CLI (`gh`)** — the one new read (check runs for a commit), same
  binary and auth preflight already requires for review Runs.
- No new external systems; CodeRabbit interaction is unchanged.

## Testing Approach

Existing seams: watch tests drive the loop with fake `StatusFunc`/`CheckFunc`
clock-stepped adapters (poll-first asserted by call-ordering against the fake
clock; quiet-skip by elapsed-time assertions; confirm phase by check-state
scripts covering success, failure→new-Round, missing, and timeout). CLI tests
assert the stdout reports byte-exactly for terminal outcomes and the
`--no-agent-console` stderr stream containing zero agent-source lines. The
filter sink gets unit tests in `internal/runevent`. Full suite is the
regression net; existing watch outcome tests update only where ADR-0019
deliberately changes the Clean condition.

## Build Order

1. Poll-first ordering and pre-settled quiet skip (no deps)
2. Deterministic stdout reports for watch and resolve (no deps)
3. Agent-console suppression flag and filter sink (no deps)
4. Merge-readiness confirm phase: CheckFunc seam, gh adapter, outcome wiring
   (depends on: 1)
5. Docs and skill sync — flag, stdout contract, ADR-0019 semantics
   (depends on: 2, 3, 4)

## Risks & Considerations

- The check-run filter must identify the Review Source's app robustly
  (CodeRabbit's check name/app slug); a rename upstream degrades to
  `missing` → Clean-with-note, never a hang.
- The confirm phase extends wall-clock after the Final Push; bounded by the
  existing review timeout — document the expectation.
- ADR-0019 changes when Clean is declared: any script timing on watch exit
  observes longer Runs; exit codes and states are unchanged.
- `gh` rate limits: one call per poll interval is well inside limits.

## Decisions

- Merge-readiness is default `--until-clean` semantics. See ADR-0019.
- `missing` check ends Clean with a note (repos without the Review Source
  check keep old behavior).
- Check polling reuses `poll_interval` + review timeout; no new config.
- The stderr filter is a Sink decorator; the Journal is never filtered.
