---
status: accepted
created_at: 2026-10-27T15:00:00Z
updated_at: 2026-10-27T15:00:00Z
deprecated_at: null
superseded_by: null
---

# Settle follows the branch's plan, and a declared row matches the criterion it names

On 2026-10-27 Quaestor reported two refusals of one recovery. Its QA gate
blocked two Success Metrics in separate rows that one `## Unreachable
Acceptance` declaration covered, and `roundfix qa-report accept` refused the
report because `rows_blocked_declared` (2) exceeded the number of
declarations (1): ADR-0240 says the declared-blocked rows must be "covered by
the Spec's declarations", and the code measured coverage by count. After
Quaestor split the declaration on its branch, `roundfix settle` read the old
PRD from the kept Run Worktree and refused again until the worktree was moved
to the new commit by hand. The maintainer chose that acceptance matches row
by row: "Uma declaração cobre cada criterion que nomeia; o accept confere
linha a linha, não a contagem."

## Decision

**A declared row matches the criterion it names.** A declaration names one or
more criteria, one `criterion:` line each, and shares its `reason` and
`satisfied-by` among them. A Results row recorded
`blocked (declared: <criterion>)` matches when a declaration names that
criterion: the same text once inline-code backticks are dropped, whitespace
is collapsed and case and a trailing period are ignored, or the criterion's
leading label, the words before its first comma, colon, semicolon, period,
` (`, ` —` or ` - `, such as `Success Metric 5`. Matching only widens
acceptance; it never refuses a report today's rule accepts. A report is
covered when the Results table lists declared-blocked rows, every one
matches, and `rows_blocked_declared` does not exceed them; otherwise it is
covered by today's rule, `rows_blocked_declared` not exceeding the number of
declarations. When both fail, the refusal keeps today's shortfall words and
names each listed row no declaration names. The rule lives in
`internal/spec` beside the other partial rules, so
`roundfix qa-report accept`, `roundfix settle`, Daemon settlement and
`roundfix archive` read it alike, and the mechanical stage classifies a
declared row with the same function.

**Settle follows the branch's plan.** When settle selects a kept Run Worktree
or Task Worktree and the Run's target branch changed files under the Spec
folder after that worktree forked from it, settle reads the plan the branch
carries. A Run Worktree whose HEAD is an ancestor of the branch is
fast-forwarded to it before the Spec is loaded, and settle says so on stderr.
Any other such drift, a branch that diverged or a Task Worktree, is refused
before anything is written, naming the changed files. A fast-forward never
rewrites a Run commit, so the settlement commits Task Carry-Forward reads keep
their bytes, and the Task settled after it records the amended plan.

**Settle runs from the user checkout.** Settle refuses to run inside a Run
Worktree or a Task Worktree, whose branch is a `roundfix/run-` branch, and
names the user checkout the Run recorded.

**Settle commits what a Task commit keeps.** Settle refuses an untracked
executable file that the authorization it reads does not list under
`executables`, as a Task commit does (ADR-0340). A tracked file whose mode a
Task changes stays work output in both, as ADR-0157 decided.

## Considered Options

- Matching alone, refusing every row no declaration names: rejected. A
  report accepted today, such as Spec 0324's whose declared row paraphrases
  its criterion, would be refused after an upgrade; Specs evolve and never
  regress, so the count stays as the fallback.
- Compare counts of criteria instead of declarations: rejected. A count
  cannot tell which criterion a row stands for, and it would widen the
  fallback to rows the Spec never declared.
- Exact criterion text only: rejected. Gates write the label a declaration
  starts with, such as `Success Metric 5` (Spec 0326's report).
- Read the plan from the user checkout while settling the Run Worktree's work:
  rejected. The settle commit would carry a QA Report judged against a plan
  its own branch does not hold, and the integration would still be refused.
- Always refuse a stale plan and leave the fast-forward to the user: rejected.
  The fast-forward is the step Quaestor took by hand, Git performs it only
  when it overwrites no local change, and refusing would keep the silent
  stale read as an extra command.
- Resolve the user checkout from the Run record and settle anyway from inside
  the worktree: rejected. Settle then removes the directory the caller stands
  in, and the caller's next Git command exits 128 there.
- Require the executables grant for a tracked file that gains an execute bit:
  rejected. ADR-0157 holds a tracked path as work output whatever its mode,
  and a tracked file is not the build artifact ADR-0340's list guards
  against; the mode change stays visible in the diff.

## Consequences

- ADR-0240's coverage of declared rows widens. A Spec whose declaration
  names several criteria is accepted when the gate records one row per
  criterion; every report the count accepted is still accepted.
- A refused report names the declared rows no declaration names; the
  `qa-gate` skill says to write the criterion or its leading label.
- Settle may move a kept Run Worktree forward before its Verification; a
  failed Verification leaves it at the branch's commit with the Task status
  and the uncommitted work unchanged.
- Settle run from inside a Run Worktree now exits 2 instead of committing and
  reporting `integration pending`.
