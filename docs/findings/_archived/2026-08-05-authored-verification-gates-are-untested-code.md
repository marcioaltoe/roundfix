---
status: done
created_at: 2026-08-05
updated_at: 2026-08-06
absorbed_by: 2026-08-06-rollup-qa-gates-and-verification-evidence.md
---

# 2026-08-05 — Authored verification gates are untested code

status: pending

Source: the second real Supervisor session in the `fiscus` repository, spec
`0002-auth-staff-e-directory` (runs `run_20260805T131149Z_18c1483e9d4abee4` Stopped and
`run_20260805T134406Z_41a3bfa9659f0917` in flight), one day after the spec-0001 night that
produced `2026-08-05-five-frictions-from-a-full-autonomous-spec-night.md`. The unifying
observation: **a Task's `## Verification` commands are the only code in the pipeline that is
never executed before it starts deciding outcomes.** Implementation passes through the gates;
the gates themselves ship untested. One run burned three defective gates authored by a careful
Supervisor following the current `write-tasks` contract — the contract itself is missing a
step. This finding proposes improvements to the `write-*` skills (whose Baseline roundfix
owns) and to Roundfix behavior.

## 1. Three authored-gate defect classes, one run

What was observed, each with the exact failing command:

- **Structurally impossible** — task_04 declared
  `grep -rq "CREATE TABLE" … && grep -rql "dir" … | grep -q .`. The `-q` flag mutes `-l`'s
  filename output, so the trailing `grep -q .` never receives a byte: the command has no
  reachable success state. Attempt 2 failed on it while the work was demonstrably complete
  (the migration existed in the Task Worktree; a manual rerun of the first grep exited 0).
  Diagnostics log: **0 bytes**, because a muted grep fails silently.
- **Self-contradictory negative** — task_07 declared
  `! grep -rq "secret" …/schema/` while its own Requirement 1 mandates persisting a
  `secret_hash` column. Correct work makes the gate fail forever. Caught by a Supervisor sweep
  before it burned a cycle, only because the task_04 failure prompted suspicion.
- **Non-idempotent task-owned setup under shared state** — every DB task declared
  `docker compose up -d postgres-test --wait`. The repository compose pins
  `container_name: fiscus-postgres-test`; each Task Worktree is a distinct compose project, so
  the first project to create the container owns it and every later `up` — from any other
  worktree, forever — dies with `Conflict. The container name … is already in use`
  (batch-002-attempt-1.log, 725 bytes). The container was running and healthy; only `up`
  refused. task_02 burned its single repair turn removing the stale container instead of
  repairing code.

The existing `write-tasks` guard covers only the inverse class ("every Verification command
must be able to fail when no work was done" — the vacuous always-pass). Nothing checks that a
command **can pass against correct work**.

## 2. Suggested `write-tasks` improvement: a mandatory Verification lint

Add a gate-lint step to the graph-validation phase (the step that already parses nodes, needs,
and types mechanically):

1. **Dry-run every command in the pre-work tree** and classify: *red for the right reason*
   (the artifact the task will create is missing), *green as precondition*, or *anything else =
   defect*. The task_04 grep fails this immediately: no author can explain what future state
   turns it green.
2. **Adversarial negative check**: for every `! grep <pattern>` (or equivalent negative
   assertion), search the same task file's Requirements/Subtasks for identifiers matching the
   pattern. A hit is a self-contradiction (catches `secret` vs `secret_hash` by reading the
   task file alone).
3. **Idempotence rule for setup commands**: any command that provisions shared state
   (containers, databases, ports, files outside the worktree) must tolerate the state already
   existing. `docker compose up` against a name-pinned singleton fails this; the repaired form
   (`docker inspect … | grep -q true || docker compose up -d … --wait`) passes.
4. **Banned-construct list** (cheap static screen): `-q` combined with `-l`/`-c` feeding a
   pipe; any pipeline whose producer is output-silenced feeding `grep -q`; count comparisons
   over silenced commands.
5. **Grep is presence, tests are semantics**: semantic assertions belong in a test file the
   suite runs (executed, typed, reviewed code); `## Verification` keeps cheap presence checks
   plus the suite invocation. Most of the broken gates were semantic claims encoded in shell.

## 3. Suggested `write-techspec` improvement: CI/tooling impact row

What was observed (spec 0001): the techspec declared "no protected tooling mutation" while
introducing the repository's **first database-backed test** — which the CI workflow could not
execute (no database service). The merge froze overnight on a red check whose fix was an
unauthorized `.github/workflows` edit. The constraint resolution asked "does this spec *touch*
tooling?" but never "can CI *execute* what this spec creates?".

Suggested behavior: the techspec's Project Constraints resolution asks the second question
whenever a spec changes the verification class (first DB test, first E2E, first image build),
and requests the pre-authorization **with bounded files at authoring time** — when the
maintainer is awake. Applied manually in spec 0002 (drizzle.config pre-grant in the PRD/
TechSpec Tooling row): the authorized task ran first and unblocked everything downstream, zero
freezes.

## 4. Suggested `write-tasks` decomposition rule: serialize ordinal-artifact generators

What was observed: tasks 02 and 04 ran in parallel waves and each generated a Drizzle
migration — both minted **index 0001** with its snapshot and journal entry. task_02 integrated
first; task_04's worktree then held `0001_directory_schema_tables.sql` + a conflicting
`meta/_journal.json`/`meta/0001_snapshot.json`, guaranteeing a cherry-pick conflict for both
the integration queue and any later `settle`. The stranded work had to be discarded and redone.

Suggested behavior: decomposition adds an explicit `needs` edge between any two tasks that
generate sequentially-numbered artifacts from shared metadata (migration indexes, snapshot
chains), even when logically independent. A one-line rule in the skill prevents a guaranteed
conflict class that no amount of agent skill avoids.

## 5. Roundfix: empty-diagnostics failures should implicate the gate, not the work

What was observed: task_04's two failed attempts produced 0-byte diagnostics logs, and the
terminal reason named the command — but the operator guidance (`Inspect the diagnostics`)
points at an empty file. The signature *same command failed twice + zero output captured* is
strong evidence the gate itself is defective (silent greps being the canonical case).

Suggested behavior: when a Verification attempt fails with an empty captured output, append to
the failure reason a hint shaped like
`command produced no output; a structurally failing check (e.g. silenced grep) may be the
defect — review the Task's ## Verification`. Optionally, a static
`roundfix verify-lint --spec <slug>` (or an implement-preflight warning) applying section 2's
banned-construct screen would catch the whole class before a Run starts.

## 6. Roundfix: the preflight sweep destroys work that reconcile preserves

What was observed: the Implement preflight sweep reaps terminal Worktrees "whose branch has no
commits beyond its base" — and a **failed Task's Worktree is exactly that shape**: the Daemon
withheld the commit, so the branch has zero commits while the working tree holds the only copy
of the Agent's work. Run #2's preflight printed

```
roundfix: reaped terminal Worktree path=…run_20260805T131149Z_18c1483e9d4abee4.task_04 …
```

destroying the uncommitted migration and schema edits (intentional discard in this case — but
only because the Supervisor had already decided to regenerate). Reconcile classifies the same
worktree `dirty → preserve`; the sweep and reconcile disagree about the same surface, and the
sweep wins silently at the worst moment: right when the operator re-runs implement to recover
from a failure.

Suggested behavior: the sweep's eligibility check should exclude worktrees with tracked or
untracked changes (reconcile's `dirty` test), or at minimum print the dirty-file count in the
reap line so the destruction is visible. `settle`'s recovery contract depends on that surface
surviving until the operator chooses.

## Companion

`2026-08-05-five-frictions-from-a-full-autonomous-spec-night.md` records the spec-0001 night
(gate-invalidation recovery, doctor vs hookful worktrees, post-settlement commit failure,
Rounds burned on evidence waits, nitpick deadlock). Together the two findings sketch the same
theme from both sides: the pipeline's own control surfaces — gates, sweeps, waits — need the
scrutiny the pipeline applies to product code.

## Addendum — 2026-08-05 — two archive-path contract gaps (spec 0002 endgame)

7. **The Daemon settles a declared-partial gate `failed`; Archive requires `completed`.** The
   final QA report closed `verdict: partial`, zero findings, with the unmet rows covered by the
   Spec's `## Unreachable Acceptance` declarations — the exact shape Archive documents as
   eligible ("an archive-eligible partial Report can still leave the Run Unresolved; that
   outcome does not prevent Archive"). But the Daemon settles any non-`pass` gate as `failed`,
   and Archive's first preflight check demands every Task `completed` — so the documented
   eligible state is unreachable without the operator hand-editing the Daemon-owned gate status.
   Suggested behavior: the Daemon settles an archive-eligible declared-partial gate `completed`
   (journaling the verdict), or Archive's task-status check accepts the gate Task in `failed`
   when the newest report is archive-eligible.
8. **Declaration matching is left to gate-agent judgment, and Archive counts the label.** With a
   declaration covering "rows requiring loopback PostgreSQL inside the QA sandbox", the gate
   still labeled the repository-verification row `blocked (environment: …)` — after itself
   re-running the same command green outside the sandbox on the same build — and Archive then
   refused on `rows_blocked_environment is 1; expected 0`. The operator's only recourse was
   sharpening the declaration text to name the row class and prescribe the label, then
   re-running the whole gate for accounting. Suggested behavior: when a blocked row's blocker
   matches a Spec declaration, the gate classifies it `declared` (environment is the fallback
   for undeclared blockers), or Archive treats environment-blocked rows whose blocker text
   matches a declaration as declared.

Related SQLITE_BUSY note: one Task batch in the same spec died with
`publish Run Events: begin Run Event append: database is locked (5) (SQLITE_BUSY)` — an Agent
batch failure with no work defect; a bounded retry on Run-Event appends would absorb it.

## Addendum — 2026-08-05 — closing a failed gate needs a hand-edit the docs do not name (spec 0077)

9. **Appending the corrective Task leaves the graph unloadable until the operator rewrites the
   Daemon-owned gate status.** The documented recovery from a `fail` verdict is to append a
   corrective Task, and `StaleGateError` is the mechanism that makes the insertion loud rather
   than silent. But it is raised from `loadTaskGraph`, so between authoring `task_06` and
   resetting `task_05`, every command that loads the Spec — `spec check`, `implement`
   preflight — refuses with `QA gate result is invalidated`. The only way through is editing
   `status: failed` back to `status: pending` in a Task file whose status the Daemon otherwise
   owns exclusively, and no document names that edit as the operator's next step. The error
   text states the problem and not the action. Suggested behavior: `StaleGateError` carries a
   `next:` line naming the reset, or the loader treats a settled gate above uncompleted
   dependencies as `pending` and journals the invalidation instead of refusing the load.

## Addendum — 2026-08-05 — two more classes from the spec-0003 run

1. **Always-green commands defeat a red-only lint.** `bun --cwd <dir> run test` (space form)
   prints the script list and exits 0 without running anything — the correct form is
   `--cwd=<dir>`. Every authored suite gate across three specs used the vacuous form; the
   red-pre-work lint missed it because the class is *always green*, not never-green. The lint
   needs a **fail-ability proof** for suite commands (run them and require evidence of real
   execution — test counts in output — or intentionally break an input and require red), not
   only red-for-the-right-reason proofs for artifact greps. The roundfix skill itself documents
   this exact bun pitfall for agents; gate authoring needs the same rule.
2. **The Daemon's staging list can name files deleted during the turn.** task_04's agent
    created `components/.gitkeep`, later replaced it with the real component and removed it;
    the Daemon's `git add -f -- <paths>` still listed the .gitkeep and exited 128
    (`pathspec did not match`), failing the whole Run *after* the task settled completed — the
    0001-night post-settlement failure shape again, via a different door. Suggested behavior:
    build the staging list from the tree at commit time (or pass `--ignore-errors` / filter
    non-existent paths), and settle the task `failed` (recoverable) instead of failing the Run
    when the commit cannot be created.
