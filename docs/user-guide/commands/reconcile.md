### reconcile

```bash
roundfix reconcile [run-id] [--apply | --discard-superseded | --carry-forward] [--format <text|json>]
```

The Reconcile Command classifies one terminal spec Run in the current
repository when a Run ID is supplied. Without a Run ID, it scans every
terminal spec Run for the current repository. Active Runs, review Runs, Runs
from another repository, and missing Run IDs fail Preflight Validation without
Git mutation.

The default is a dry-run. These examples inspect one Run, inspect the current
repository, and request the versioned JSON report:

```bash
roundfix reconcile <run-id>
roundfix reconcile
roundfix reconcile <run-id> --format json
```

Text and JSON reports are requested output and go to stdout. Diagnostics,
including validation and operational failures, go to stderr. Redirect them
independently when automating:

```bash
roundfix reconcile --format json > reconcile.json
roundfix runs list > runs.txt 2> runs.diagnostics
```

The text report includes the repository, mode, each Run's outcome,
classification, Run Worktree, Run Branch, target branch, both resolved heads,
evidence, action, refusal reason, summary counts, and the exact apply command.
JSON uses the `roundfix-reconcile/v1` envelope with `mode`, `repository`,
`applyCommand`, `results`, and `summary`. Its `stagingCandidates` array reports
each registered carry-forward staging worktree, its owner PID, proof, action,
and refusal reason. `debrisSummary.stagingCandidates` counts those entries and
`debrisSummary.stagingApplied` counts the staging worktrees released in that
invocation; every existing report field keeps its meaning.

For a Run of a merged Spec, reconciliation proves the Run against the merged
head. The Delivery Queue merge record is the primary source; when no usable
record remains, the default branch carrying the archived Spec is the fallback.
The existing `evidence` field names the source and proof without changing the
text or JSON report shape.

Run Worktree Reconciliation uses six states:

| State | Evidence and behavior |
| --- | --- |
| `safe` | The Run Worktree is clean and the Run Branch is contained in its target or merged head, or its changed content is represented at the merged head. Eligible for cleanup after revalidation. |
| `superseded` | A newer QA Report or the merged-head proof represents the Run's Task or QA Report commits. Eligible for cleanup after revalidation. |
| `unintegrated` | The worktree cleanliness and ref evidence resolve, but the Run Branch tip is not an ancestor of the target tip. Roundfix preserves the worktree and branch. |
| `dirty` | A present registered Run Worktree has tracked or untracked changes. Dirty evidence takes precedence and Roundfix preserves the worktree and branch. |
| `unknown` | Invalid or missing metadata, an unsafe or unregistered worktree, an ambiguous or missing ref, or a Git inspection failure prevents proof. Roundfix preserves every surface it can identify. |
| `released` | Both the Run Worktree and Run Branch are absent. Repeated dry-run or apply is an idempotent no-op. |

A missing worktree alone is not `released` and never authorizes deletion. When
the Run Branch remains, Roundfix still requires an unambiguous Run Branch tip,
the recorded target tip, and ancestry proof. Age and terminal outcome are also
not cleanup evidence.

Three switches mutate; every other invocation is a report. Each acts on a
different disposition, and none of them bypasses the proof above:

```bash
roundfix reconcile <run-id> --apply               # release safe and superseded surfaces
roundfix reconcile --apply
roundfix reconcile <run-id> --discard-superseded  # discard a Run Branch proven superseded
roundfix reconcile <run-id> --carry-forward       # hand a Stopped or Unresolved Run's settled Tasks back
```

`--discard-superseded` writes the branch record before removing anything, so a
discard that cannot be recorded does not happen. `--carry-forward` accepts a
terminal Run whose outcome is Stopped or Unresolved. Every other terminal
outcome is refused by name; the refusal names the actual outcome and says
which outcomes carry-forward accepts. It proves the candidates
in the order the Run integrated them, and the `carryForwards` JSON array
follows that order. It compares each
candidate's declared inputs against the accumulating staged carries—the
checkout plus the carries staged ahead of it—rather than against the raw
checkout. If one candidate refuses, the whole set is refused rather than
carrying part of it, so a Task whose Spec, instruction, or Context moved since
settlement is never silently replayed.

Carry-forward records each carried Task as completed and appends its source Run
and settlement commit to the Task file. A Task already completed on the
checkout is reported as `already completed; nothing to carry`; it is not
staged or proved and never refuses the remaining set. Repeating carry-forward
when every candidate is already completed succeeds without moving `HEAD`.
When several Runs qualify, choose the Run that `implement` names, because it
has the largest carriable set and uses the newest Run to break ties.

Carry-forward staging commits run without repository hooks because the carried commits already passed Daemon Verification and the repository hooks when the Daemon settled them. The checkout still receives the staged commits only through a fast-forward merge, and carry-forward does not change its Git configuration.

`--apply` remains the only switch that releases Run Worktrees.

Reconcile also sweeps registered carry-forward staging worktrees whose path is
`roundfix-carry-forward-*/worktree`. Dry-run reports every candidate without
removing it. A staging worktree is stale only when the recorded PID's
`OwnerProcessIdentity` lookup fails or differs from `owner.json`, or when a
legacy registration has no owner record and is still `locked initializing`. A
live matching owner, an unlocked legacy registration, or an unreadable or
malformed owner record stays in `preservedCandidates` with the refusal reason.
`--apply` releases stale staging, and `--carry-forward` releases it before
creating that mode's own staging worktree.

There is no force flag or user assertion that bypasses the proof. Apply acts
only on entries classified `safe` or `superseded` during that invocation, then
rechecks the metadata, worktree cleanliness, Run head, target or merged head,
and the applicable ancestry, content, Task, and QA Report evidence before
mutation. It removes the Run Worktree without force, deletes the Run Branch,
and reports failures while preserving any remaining path or ref. Dirty,
unintegrated, unknown, and released entries remain successful preserved
results unless an operational inspection fails.

Before cleanup, Roundfix records the reconciliation evidence. A safe
Integration Pending Run moves to Clean through the guarded terminal transition;
every other terminal outcome remains unchanged. The Reconcile Command never
integrates unique commits, repairs dirty work, chooses another target branch,
or treats a missing path as proof.

The contract uses the [Roundfix glossary](../../../CONTEXT.md#language) and follows
[ADR-0053](../../adr/0053-terminal-run-worktree-reconciliation-is-proof-based.md)
and
[Spec 0038](../../history/specs/0038-terminal-run-worktree-reconciliation/_prd.md).
Adjacent terminal-cleanup diagnostics remain traced through the
[Stop Command](stop.md#stop) to the
[detached-watch finding](../../history/findings/2026-07-16-vortex-pr87-detached-watch-notification.md#4-cleanup-noise-appeared-before-the-actionable-failure).

