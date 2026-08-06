# Every Run that failed tonight failed on a contract

**Date:** 2026-08-06
**Session:** five Specs delivered — 0077, 0078, 0069, 0075, 0065 — plus the
Run Database `busy_timeout` fix, over roughly twelve hours of supervised
autonomous work.

Across those five Specs, Runs failed fifteen times. **Not one failure was a
defect in the code an Agent wrote.** Every one was a contract: between the
Task's declared Verification and what it meant to prove, between the Daemon's
guards and a legitimate decomposition, or between two artifacts that were
supposed to agree.

That distinction is the finding. The queue moves; what costs time is the
contract between whoever authors and whoever executes.

## The seven classes, each with its measurement

### 1. The entry precondition is incompatible with a Task that repairs red

The Daemon runs the configured repository Verification command as a
precondition when a Task declares it verbatim. Spec 0075's `task_02` existed to
re-record the corpus that `task_01` legitimately invalidated — so the
precondition failed on exactly the state the Task repairs, and the Task settled
`failed` **without ever starting**.

```text
repository not green on entry: command "make verify" exited with exit status 2
```

Cost: one Run. Workaround: name the same targets individually
(`make fmt-check test spec-budget skills-sync-check skills-check build
spec-check`), which is the same gate under a different string.

Suggested behaviour: the precondition should be skippable by a Task that
declares it repairs the gate, or should not apply to a Task whose graph
position follows a Task that declares it leaves the repository red.

### 2. A literal assertion does not just fail — it masks the diagnostic

`TestReadoptionCompatibilityMaintainedFixture` asserted `104` entries.
`task_01` added two, and it failed for no defect. Two Runs were then spent on
wrong hypotheses — "a corpus is missing", "the regeneration steps are ordered
wrong" — because **the real error was behind the stale literal**. Once the
literal was fixed, the true message appeared in one line and named its own
remedy:

```text
clause.context.backlog-01-operational-contract; the regenerator maintains
manifest rows but never creates them, so add this row first
```

Cost: three Runs. The repository's existing rule — "an assertion reads the
constant it means" — already names the first-order defect. What it does not yet
say is the second-order one: **a frail assertion hides every diagnostic behind
it**, which is what made this expensive rather than merely annoying.

### 3. Authored requirements that forbid the only action that unblocks them

Spec 0075's `task_02` required `MUST NOT hand-edit any derived artifact`. The
regenerator states plainly that it maintains Source Baseline manifest rows and
never creates them. The prohibition therefore forbade the single step the Task
existed to perform.

Cost: one Run. Spec 0065's `SC-REQUIREMENT-CONTRADICTORY` now refuses this
shape at authoring time — it shipped the same night, from the same queue.

### 4. Verifications that assert something other than the contract

Three separate instances in one night:

- `^type: (feat|fix|perf|refactor)$` rejected `type: perf # feat | fix | perf |
  refactor` — the template the Spec itself defines, with its documented inline
  enum comment.
- An absence search over five carriers omitted Task files, so it could not
  prove its own acceptance criterion.
- Widening that same search to the Spec folder would then have failed on the
  Spec's own prose *explaining* the correction. **An absence check that cannot
  distinguish a claim from a citation of that claim is unusable in both
  directions.**

Cost: two Runs and two gate cycles.

### 5. A test that pins a path the workflow is guaranteed to move

`TestCheckLoopOrderRepositoryAgrees` named `docs/specs/0065-.../` to run the
check against. Archiving moves that folder, and archiving is the last step of
every Spec — so the test was guaranteed to break at the moment the Spec that
created it archived itself. It did, after the Run had already reported Clean.

The rule this repository allows makes it worse: an archived Spec may be deleted
at any time, so even the archived path is not a durable anchor. The fix is a
dedicated fixture carrier, independent of both active and archived Specs.

### 6. Shared git worktrees leak Run Branches across Specs

Four `git worktree` checkouts of the same repository share one object store and
one set of refs. Branch Integrity Preflight for Spec 0075 therefore refused,
naming a `roundfix/run-*` branch belonging to Spec 0078's Run — and prescribed
integrating it, which would have overwritten newer artifacts with superseded
ones.

Worse, the documented bypass has a second-order cost: `--skip-branch-integrity`
also skips the Active-Run guard, which let two Runs operate on the same
worktree and branch simultaneously. That was recoverable only because the
conflict was noticed by hand.

Suggested behaviour: Branch Integrity Preflight should scope pending Run Branch
work to the Run's own repository checkout, not to every ref visible in the
shared object store. And the bypass should be separable — one flag per
guardrail, since they protect different things.

### 7. Typed transient failures are a list, and the list is incomplete

`gh` exiting non-zero from a network failure ends a Run `Failed`. The same
underlying failure, when it surfaces as a context deadline, DNS failure,
connection reset, HTTP 429, or a GitHub 5xx, is retried. Two Runs died this way,
one with a bare `gh command failed`, one with `HTTP 422`.

Classifying by presentation rather than by cause means the retry contract holds
only for the shapes someone happened to enumerate.

## What already shipped against this

Spec 0065, delivered the same night, turns three of these classes into
mechanical `SC-*` rules: `SC-VERIFY-WORK-INDEPENDENT`,
`SC-REQUIREMENT-CONTRADICTORY`, and `SC-REHEARSAL-UNDECLARED`. The Run Database
`busy_timeout` moved from 5s to 30s after two Agent Batches died on
`SQLITE_BUSY` writing their completion event — an infrastructure failure with no
work defect, which failed the Task and cost a recovery Run each time.

## What has not

Classes 1, 5, 6 and 7 have no owner. Class 2's second-order effect — masking —
is not stated anywhere. These are the candidates for the Spec that follows.

## The number worth carrying

Fifteen Run failures, zero implementation defects. Every Spec that reached its
gate passed it, most on the first or second cycle, and every finding the gates
produced was real — including two that were errors in artifacts this session
authored, caught by a gate reading the archived source rather than the carrier
that quoted it.
