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

Only a full scan (`roundfix reconcile` without a Run ID) inspects item branches
named `roundfix/deliver-<slug>-<16 hex>`. A live item's branch is preserved with
its Delivery Queue item. `--apply` deletes a proven item branch with its clean
worktree. The `roundfix-reconcile/v1` report's `itemBranchCandidates` list
names these entries, and `debrisSummary` reports `itemBranchCandidates` and
`itemBranchesApplied` counts.

For a Run of a merged Spec, reconciliation proves the Run against the merged
head. The Delivery Queue merge record is the primary source; when no usable
record remains, the default branch carrying the archived Spec is the fallback.
The existing `evidence` field names the source and proof without changing the
text or JSON report shape.

When the default branch holds the archived `<archive root>/<slug>/_prd.md` and
the delivery commit that added it is not in the branch being proven, the Spec
carries merge evidence. With merge evidence, Task commits still need their
Task completed in the archived Spec, and every other commit is superseded by
the delivery of Spec; an unrepresented Task commit preserves the Run. Without
merge evidence, the content comparison is unchanged.

A dirty Run Worktree is `superseded` only when the commit proof at the merged
head or by the delivery is `safe` or `superseded` and every uncommitted path
lies under the Spec's own directories, equals a declared or recorded path, or
lies beneath a declared or recorded path (`<declared>/` prefix). The scope
comes from the archived Tasks' `interface:`, `creates:`, `deletes:` and
`## Recorded paths`, read at the Archive Record's `source_revision`. A declared
directory need not still exist; a sibling that only shares its prefix stays
outside the scope. The superseded reason appends
`uncommitted path(s) superseded by the archived Spec`, with the path count.

A kept dirty Run whose Spec is archived at the chosen merged head or the
default branch names why: `Run Worktree has tracked or untracked changes; <commit proof reason>`
when the commit proof fails, otherwise
`Run Worktree has tracked or untracked changes; <path> is outside the archived Spec's declared and recorded paths`.
The path is the first outside the scope in byte order. Both reasons fit 160
bytes by shortening the suffix's middle, preserving the constant prefix. A
Run whose Spec is not archived there keeps the bare
`Run Worktree has tracked or untracked changes`; an unavailable archive source
keeps its existing diagnostic. Apply re-proves the state, Run head, merged
head and exact dirty path set before forced removal; clean worktrees are
removed without force. A changed path set refuses apply as stale.
A dirty item worktree keeps its item branch.
An absent target branch also uses the merged-head proof before preserving a
terminal Run Branch, including Runs recorded from linked repository worktrees.
All proof comes from local Git; reconciliation does not read GitHub.

Run Worktree Reconciliation uses six states:

| State | Evidence and behavior |
| --- | --- |
| `safe` | The Run Worktree is clean and the Run Branch is contained in its target or merged head, or its changed content is represented at the merged head. Eligible for cleanup after revalidation. |
| `superseded` | A newer QA Report or the merged-head proof represents the Run's Task or QA Report commits, archived Spec-directory work, or scoped uncommitted leftovers. Eligible for cleanup after revalidation. |
| `unintegrated` | The worktree cleanliness and ref evidence resolve, but the Run Branch tip is not an ancestor of the target tip. Roundfix preserves the worktree and branch. |
| `dirty` | Tracked or untracked changes lack archived-Spec supersession proof, even after including paths beneath a declared or recorded path. For an archived Spec, `Run Worktree has tracked or untracked changes; <commit proof reason>` names the failed commit proof, or `Run Worktree has tracked or untracked changes; <path> is outside the archived Spec's declared and recorded paths` names the first outside path in byte order; the suffix is shortened in the middle to fit 160 bytes. An unarchived Spec keeps the bare constant; an unavailable archive source keeps its diagnostic. Roundfix preserves the worktree and branch. |
| `unknown` | Invalid or missing metadata, an unsafe or unregistered worktree, an ambiguous or missing ref, or a Git inspection failure prevents proof. Roundfix preserves every surface it can identify. |
| `released` | Both the Run Worktree and Run Branch are absent. Repeated dry-run or apply is an idempotent no-op. |

When `worktree.location` is a symbolic link, a recorded Run Worktree is
resolved within its Worktree Root. A recorded path whose repository or
worktree directory is itself a symbolic link stays `unknown`. Run
`roundfix reconcile` to inspect the retained Runs, then
`roundfix reconcile --apply` to release eligible Runs that an earlier version
kept with `recorded Run Worktree path is unsafe`.

A missing worktree alone is not `released` and never authorizes deletion. When
the Run Branch remains, Roundfix still requires an unambiguous Run Branch tip,
and either target ancestry, QA supersession, or merged-head representation
proof. Age and terminal outcome are also not cleanup evidence.

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
carrying part of it.

A settled Task can be carried over an amendment of its own Spec's `_prd.md`,
`_techspec.md` or `_tasks.md` when its Binding Text is unchanged. In the manifest,
Binding Text is the Task's graph node (`id`, `file`, `needs` in order) and its
projection table row. In the PRD and TechSpec it is every section named by the
Task's References bullet for that file, from its level 2–4 heading to the next
heading of the same or a higher level. Matching ignores case, accepts singular
or plural forms, and drops numbered suffixes, ranges, commas and `and`; a unit
ends at its first `:`. Headings inside fenced blocks do not match. A missing
node, row or References bullet, or any unresolved unit, binds the whole file
and keeps the moved-input refusal.

Both bindings use the Task file at the settlement parent and compare that
parent's planning text with the accumulating staged target. A changed Binding
Text still refuses the whole remaining set. The Task file, `AGENTS.md`,
`CONTEXT.md`, the `implement-task` skill, `docs/agents/skill-dispatch.md` and
every explicit Context path still refuse on any change. This rule also serves
Implement Preflight and Delivery Retry; their refusal reasons and recovery
commands are unchanged.

Each candidate lists tolerated planning amendments separately. Text prints these
lines in order (repository-relative paths replace the placeholders):

```text
  moved-inputs: []
  amended-inputs: [<paths>]
```

The `carryForwards` JSON entry includes `"movedInputs":[]` and
`"amendedInputs":["<path>"]`. `amendedInputs` is `[]` when there are no tolerated
amendments; those paths do not appear in `movedInputs` or set `inputsMoved`.

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

Reconcile also proves the recorded owner of each terminal Run. When the live process at the owner PID presents a start identity that differs from the one the Run recorded, the system gave that PID to another process and the recorded owner ended: the Run's owner PID and identity are a Stale Run Owner Record. The report lists it once in `processCandidates` with `kind` `runOwnerRecord`, the owner PID, an empty `processIds` list and that proof, and the text report prints `Run owner record candidate: Run <run-id> owner PID <pid>`. `--apply` releases it: it clears the Run's recorded owner PID and identity, journals what it cleared, and sends no signal; it never inspects the processes of the PID's new holder. `debrisSummary.runOwnerRecordsReleased`, `run-owner-records-released` in the text summary, counts the releases. An owner identity that cannot be read or compared, and a Run recorded without one, stay in `preservedCandidates` with the refusal reason, and the Stop Command still never signals such a PID ([ADR-0313](../../adr/0313-a-proven-owner-identity-mismatch-releases-a-terminal-runs-owner-record-without-a-signal.md)).

There is no force flag or user assertion that bypasses the proof. Apply acts
only on entries classified `safe` or `superseded` during that invocation, then
rechecks the metadata, worktree cleanliness, Run head, target or merged head,
and the applicable ancestry, content, Task, and QA Report evidence before
mutation. It removes a clean Run Worktree without force, or uses force only for the
re-proved superseded dirty path set, deletes the Run Branch,
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
Spec 0038.
Adjacent terminal-cleanup diagnostics remain traced through the
[Stop Command](stop.md#stop) to the
detached-watch finding (`docs/history/findings/2026-07-16-vortex-pr87-detached-watch-notification.md@a30dc847a037fe584812b4471aa7da1701a9ec59`).
