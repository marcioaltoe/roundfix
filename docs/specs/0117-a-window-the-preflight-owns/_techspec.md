---
spec: 0117-a-window-the-preflight-owns
status: active
created: 2026-08-26
---

# TechSpec: A window the Preflight owns

## Project Constraints

- Identifier strategy: applicable — this Spec coins **Run Window** and emits the
  token `run_window` across a CLI command, a Run Database column, and a refusal
  message. The Vocabulary Contract below binds each emitting surface to
  `CONTEXT.md`. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no network boundary, credential, or
  request is created or read. Source: `docs/agents/specific-repository.md`.
- Active ADR obligations: applicable — ADR-0004 (central Run Database) decides
  the storage; ADR-0139 (locks keyed by work target in that database) supplies
  the scoping precedent this Spec follows; ADR-0137 (the Run budget is
  explained where it is configured) governs the second time bound; ADR-0133 (a
  diagnostic names the literal it requires) governs the refusal text; ADR-0022
  (Stop Requests travel through the Run Database) is the shape this Spec
  declines, and its Non-Goal records why. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — no protected tooling mutation is proposed or
  authorized. Production Go in `internal/store`, `internal/cli`, plus tests and
  one user-guide page. Source: `docs/agents/agent-instructions.md`.

## Vocabulary Contract

This Spec coins one term and emits one token for it. The obvious name was
unavailable: `Session` denotes an Agent Session in the glossary, so **Run
Window** names what the bound actually governs — the window during which Runs
may be created.

- emits: `internal/store/store.go`
  pattern: `run_window|RunWindow`
  documented-in: `CONTEXT.md`
- emits: `internal/cli/cli.go`
  pattern: `run_window|RunWindow`
  documented-in: `CONTEXT.md`

The two files are the durable record of the window and the command surface that
registers `window` and renders its refusal — every surface where the term
reaches a reader. The command's own implementation file does not appear here
because it does not exist yet; the contract may only name readable paths, so a
Spec authored before its code names the files that will carry the token and are
present today.

## Coverage Map

| PRD Item | TechSpec Section |
| --- | --- |
| User Story 1 | The Window Command |
| User Story 2 | The Preflight Refusal |
| User Story 3 | Starting, Never Finishing |
| Core Feature 1 | Storage And Scope |
| Core Feature 2 | Next Occurrence |
| Core Feature 3 | Idempotent Set |
| Core Feature 4 | Starting, Never Finishing |
| Core Feature 5 | Where The Two Bounds Meet |

## Context

`roundfix implement` runs a preflight block delimited in
`internal/cli/implement.go:105-196` by the comment stating that nothing is
written to the Run Database until every check has passed. Twelve checks live
there — git inspection, artifact directory, Specs root, Spec load, committed
graph, default-branch guard, three agent-selection proofs, store open, worktree
debris pruning, and the active-run lock. Each failure calls
`printPreflightFailure` (`internal/cli/cli.go:5801`), which prints the reason,
an optional next action, and the "No side effects" block, then exits `2`.

Run creation follows immediately at `implement.go:198-211` through
`store.CreateRun`. After that line the failure path changes: refusals become
Run failures with exit `1`.

Two facts from the survey shape the design. First, `budget.MaxRunDuration` is
consumed only by `internal/watch` today (`watch.go:406`, `watch.go:1067`); the
implement path never reads it, so this Spec introduces its first use there and
must not change what it bounds elsewhere. Second, the repository holds no
per-session state file: under `~/.roundfix/` there is the database, its advisory
lock, and a version-check cache, and under `.git/` only the Baseline
transaction directory. There is no existing home for a window, which is why
ADR-0004 decides the question rather than convenience.

## Solution Design

### 1. Storage And Scope

The Run Window is a row in the Run Database, not a file. Schema version rises
from 12 to 13 (`internal/store/store.go:1282`), adding:

```sql
CREATE TABLE run_windows (
  git_root   TEXT PRIMARY KEY,
  cutoff_at  INTEGER NOT NULL,   -- unix seconds
  created_at INTEGER NOT NULL
);
```

**Scope is the repository**, keyed by `git_root` — the PRD's Open Question,
settled here. The database is machine-wide, so an unscoped window would bound
Runs in a repository the Supervisor never touched; `git_root` is the same key
the single-working-tree guard already uses (`implement.go:183`,
`ActiveRunInGitRoot` at `store.go:1058`), so the scope is not a new concept.
The migration is additive: `case 12` creates the table and no existing row
changes, so a database that never sets a window behaves exactly as before.

Store API, beside the existing Run methods:

- `SetRunWindow(ctx, gitRoot string, cutoff time.Time, replace bool) (RunWindow, bool, error)` — the bool reports whether it wrote.
- `RunWindowFor(ctx, gitRoot string) (RunWindow, bool, error)`
- `ClearRunWindow(ctx, gitRoot string) (bool, error)`

### 2. Next Occurrence

The command accepts `HH:MM` and resolves it to the **next occurrence** in local
time: today at that time when it is still ahead, tomorrow otherwise. The naive
form — comparing the current hour against the target — reports a closed window
at 23:00 for an 07:00 cutoff, which kills a night session at the moment it
starts. Resolution happens once, at set time, and the stored value is an
instant: a window survives a restart, a compaction, and the clock crossing
midnight, because nothing re-derives it afterwards.

The command also accepts an absolute `YYYY-MM-DDTHH:MM` for a window that is
not on the next-day cycle. A cutoff resolved into the past is refused at set
time rather than stored.

### 3. Idempotent Set

`roundfix window set HH:MM` against a repository that already has a window
prints the standing cutoff and writes nothing, exiting `0`. `--force` replaces
it. The asymmetry is deliberate and follows the measured failure mode of the
prior art: a re-set that silently moved the cutoff forward would re-open a
window that had closed, and an unattended loop would spend the rest of the
night opening Runs. Failing to refresh a window is recoverable in one command;
re-opening one is not recoverable at all.

### 4. The Preflight Refusal

A thirteenth check joins the block, positioned **after** `store.Open`
(`implement.go:170`) because it needs the store, and **before**
`pruneTerminalRunWorktreeDebris` (`implement.go:179`) because a refused Run
should not first mutate worktree state. It reads the window for the resolved
git root; absent window means no check.

The error type implements `NextAction() string`, the extension point at
`cli.go:5808-5814` whose precedent is `preflight.TargetMismatch.NextAction`
(`internal/preflight/preflight.go:96`), so the refusal renders with the same
shape as every other:

```text
Preflight failed

Reason:
  the Run Window for /Users/x/dev/repo closed at 2026-08-27 07:00; the time is
  2026-08-27 07:14

No side effects:
  Roundfix did not create a Run, fetch Review Source issues, start an Agent,
  commit, or push.

Next action:
  move the window with `roundfix window set <HH:MM> --force`, or remove it with
  `roundfix window clear`
```

Per ADR-0133 the message names both instants literally rather than describing
them, so a reader never computes the delta to understand the refusal.

### 5. Starting, Never Finishing

The check governs Run creation only. Nothing reads the window after
`implement.go:198`, so a Run created inside the window reaches its own terminal
outcome regardless of when that happens — which is the whole point: a bound
that ended a running Run would leave implemented work with no Pull Request, and
that shape is a Stop Request on a timer, which ADR-0022 already owns and this
Spec declines.

**The crossing case reports rather than refuses.** When the window is open but
`now + budget.MaxRunDuration` falls past the cutoff, the Run is created and one
line is printed alongside the existing startup report:

```text
Run Window: closes 2026-08-27 07:00, in 12m; max_run_duration is 45m, so this
Run may run past it.
```

This is the implement path's first read of `budget.MaxRunDuration`; it is read
for reporting only and bounds nothing here, so `internal/watch`'s use of it is
untouched.

### 6. Where The Two Bounds Meet

Three surfaces state the relationship, per ADR-0137:

- `roundfix window --help` and the `window show` output say the window bounds
  when a Run may **start**, and name `budget.max_run_duration` as the bound on
  how long one may **run**.
- The crossing line above states it at the moment it matters.
- The configuration surface that prints `budget.max_run_duration`
  (`config.go:873`) gains the converse pointer.

## Specification

### The Window Command

`roundfix window <set|show|clear>`, resolving the git root the same way
`implement` does.

- `set <HH:MM|YYYY-MM-DDTHH:MM> [--force]` — resolves, refuses a past instant,
  writes unless one stands, prints the effective cutoff.
- `show` — prints the cutoff, the current time, and the remaining duration; or
  states that no window is set. Exit `0` either way: absence is a state, not a
  failure.
- `clear` — removes the window; reports whether one was there.

Exit codes follow the CLI contract at `cli.go:114-116`: `0` success, `2` for an
invalid argument or a past instant.

### The Preflight Check

1. Resolve the git root (already available at the insertion point).
2. Read the window. Absent → proceed.
3. Cutoff in the future → proceed, and evaluate the crossing report.
4. Cutoff passed → refuse with the message above, create no Run, exit `2`.

### Scope Of Application

`implement` only. `fetch`, `resolve`, and `watch` answer a Pull Request that is
already open, and refusing those by clock strands a review Round rather than
bounding an authoring session.

## Acceptance Criteria

1. A window set and passed refuses `implement` with no Run created, proven by
   the Run Database holding no new row and the exit code being `2`.
2. A window of `07:00` set at 23:00 local resolves to 07:00 the next day, and
   `implement` proceeds.
3. A second `window set` without `--force` leaves the stored cutoff
   byte-identical and exits `0`; with `--force` it replaces it.
4. A Run created before the cutoff reaches its terminal outcome after the
   cutoff passes, with no interruption attributable to the window.
5. A window whose remaining time is under `budget.max_run_duration` creates the
   Run and prints the crossing line naming both durations.
6. A database at schema 12 migrates to 13 and a repository with no window
   behaves exactly as before the change.
7. `CONTEXT.md` carries **Run Window**, and the vocabulary detector runs rather
   than skipping.

## Build Order

1. Store: schema 13, the table, and the three methods with their tests.
2. CLI: the `window` command over those methods.
3. Preflight: the check, the `NextAction` error, and the refusal.
4. Reporting: the crossing line and the two pointer sentences.
5. Docs: `CONTEXT.md` term and the user-guide page.

Steps 3 and 4 both touch `implement.go` and must not be authored as sibling
Tasks in one wave — the collision measured on 2026-08-26 between two Tasks
rewriting one file cost a passing Task at integration.
