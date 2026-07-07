---
spec: 0017-run-discovery
prd: _prd.md
created: 2026-07-06
---

# Run discovery — Technical Spec

## Executive Summary

Run discovery adds one read-only query to the Run Database and exposes it on
two surfaces: a `runs list` command and an Attach picker. The primary trade-off
is contract weight versus reach: the listing's plain-text columns become public
API that agents will parse, so the design keeps the column set minimal (id,
state, kind, target) and defers machine-readable modes entirely, accepting that
a future JSON mode will be an additive change rather than designing it now.
Everything else reuses existing seams — the store package owns the SQL, the
attach command owns the picker, and the Interactive Input pattern from the
Implement Command carries over unchanged.

## System Architecture

No new packages. Three existing modules extend:

- `internal/store` — gains one listing query over the existing `runs` table.
  The store already owns every Run read (`Run`, `ActiveRunInGitRoot`,
  `RunIDs`); the new query joins that family.
- `internal/cli` — gains a `runs` top-level command with a `list` subcommand,
  following the `skills` subcommand dispatch shape, plus the no-argument
  Attach path that opens the picker.
- `internal/tui` — the Attach picker reuses the Interactive Input rendering
  pattern (`Active Specs:` numbered list accepting a number or a value) with a
  Runs variant.

Flow: `runs list` → store listing query → formatted lines on stdout.
`attach` (no args, TTY) → same store query scoped to the repository →
Interactive Input picker → selected run id → existing attach path
(`runAttachCockpit`), unchanged.

## Implementation Design

### Interfaces

Store listing query (in `internal/store`):

```go
// ListRunsQuery scopes a Run listing.
type ListRunsQuery struct {
    GitRoot    string // empty lists every repository
    ActiveOnly bool
}

// ListRuns returns matching Runs, newest first (created_at DESC).
func (store *Store) ListRuns(ctx context.Context, q ListRunsQuery) ([]Run, error)
```

CLI dispatch (in `internal/cli`):

```go
// runRunsCommand dispatches `runs list`; unknown subcommands exit 2.
func runRunsCommand(ctx context.Context, args []string, stdout, stderr io.Writer) int

// pickAttachRun opens Interactive Input over the repository's Runs and
// returns the chosen run id. Non-interactive contexts fail before this.
func pickAttachRun(ctx context.Context, reader *store.Store, gitRoot string, stdin io.Reader, stderr io.Writer) (string, error)
```

### Data Models

No schema change. The `runs` table already carries every listed field: `id`,
`state`, `kind`, `pr_number`, `spec_slug`, `git_root`, `created_at`. The
listing derives two presentation values:

- **target** — `pr:<pr_number>` for review-kind Runs, `spec:<spec_slug>` for
  spec Runs, empty for Runs with neither.
- **active marker** — from the existing `IsTerminalState(state)` predicate;
  Active Runs are marked, terminal Runs are not.

### API Contracts

`roundfix runs list [--all] [--active]`

- stdout: one Run per line, newest first, stable column order:
  `<run-id>  <state>  <kind>  <target>`. Active Runs carry a visible marker on
  the state column. With no matching Runs, prints exactly one line
  (`No Runs found.` shape) and exits `0`.
- `--all` widens scope from the current repository to every repository in the
  Run Database; each line then also names the repository.
- `--active` filters to Active Runs only. Flags compose.
- stderr: diagnostics only. Exit codes: `0` success (including empty), `2`
  invalid usage or store open failure.
- Outside a git repository without `--all`, fails with exit `2` and an error
  naming `--all` as the alternative.

`roundfix attach` (no run id)

- Interactive terminal: opens the picker listing the repository's Runs
  (newest first, Active first) as a numbered list accepting a number or a run
  id. Selection opens the existing Live Run View. Cancel exits `0` with no
  side effects.
- Non-interactive (`--no-input` or no TTY): exits `2` with an error naming
  `roundfix runs list` as the discovery command. The current
  "a Run ID is required to attach" failure becomes this richer message.
- `attach <run-id>` behavior is byte-for-byte unchanged.

## Coverage Map

- Story 1 (list to attach/stop) → `Store.ListRuns`, `runRunsCommand`
- Story 2 (agent deterministic listing) → `runRunsCommand` stdout contract
- Story 3 (no-arg attach picker) → `pickAttachRun`, existing attach cockpit
- Core Feature 2 (repo scope + `--all`) → `ListRunsQuery.GitRoot`
- Core Feature 3 (`--active` filter) → `ListRunsQuery.ActiveOnly`
- Core Feature 4 (empty is not an error) → `runRunsCommand` empty branch
- Core Feature 6 (non-interactive failure) → attach parse/validation path

## Integration Points

None external. Both surfaces read the local Run Database through the existing
store; no GitHub, Review Source, or Agent involvement.

## Testing Approach

Existing seams only:

- `internal/store` table tests for `ListRuns`: scope by git root, `--all`
  equivalent (empty GitRoot), ActiveOnly filter, newest-first ordering, empty
  result.
- `internal/cli` buffer-captured `Run(args, &stdout, &stderr)` tests:
  `runs list` line format and ordering against a seeded temp store, empty
  listing exit `0`, unknown subcommand exit `2`, `--active` filtering,
  attach-without-id non-interactive error naming `runs list`.
- Picker logic tested by driving the input reader with scripted lines
  (number, run id, cancel), matching the existing Interactive Input tests —
  no terminal emulation.

## Build Order

1. `Store.ListRuns` with `ListRunsQuery` and store tests.
2. `runs list` command: dispatch, flags, formatting, usage text, CLI tests
   (depends on: 1).
3. Attach no-arg picker: validation split (interactive vs non-interactive),
   `pickAttachRun` over the same query, richer non-interactive error, CLI
   tests (depends on: 1).
4. Docs and skill sync: README Commands/Boundaries, usage guide, roundfix
   SKILL.md (`runs list` + attach picker), `make skills-sync` (depends on:
   2, 3).

## Risks & Considerations

- The listing columns are public API from the first release; renaming or
  reordering them later is a breaking change. Mitigation: the minimal column
  set and a CLI test pinning the byte shape.
- Run ids are long; the picker accepts a number precisely so nobody retypes
  them. The listing deliberately does not truncate ids — agents copy them
  verbatim.
- `--all` reads rows for other repositories; output must name the repository
  per line or the ids are ambiguous. Covered in the contract above.

## Decisions

- One new store query, no view/materialization — the `runs` table is small and
  indexed by head; listing is a simple scan with filters.
- `runs` becomes a command namespace with `list` as its first subcommand,
  mirroring `skills`; future Run operations (if any) join it rather than
  minting new top-level commands.
- No JSON output mode now; the stable text columns are the agent contract
  (PRD Non-Goals).
- The picker lists terminal Runs too — Attach replays them from the Run Event
  Journal, so hiding them would remove a working capability.
