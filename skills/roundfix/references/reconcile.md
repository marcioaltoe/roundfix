## Run Worktree reconciliation

Use the Reconcile Command to inspect retained terminal spec Run Worktrees and
Run Branches, plus live process trees still owned by terminal Runs. Dry-run
remains the default; always inspect it before applying cleanup:

```bash
roundfix reconcile <run-id>
roundfix reconcile
roundfix reconcile <run-id> --format json
```

A Run ID selects one terminal spec Run; omitting it scans the current
repository. The report classifies every selected Run into one of six states:

| State | Agent action |
| --- | --- |
| `safe` | The Run Worktree is clean and the Run Branch is contained in its target or merged head, or its changed content is represented at the merged head. Eligible for cleanup after revalidation. |
| `superseded` | A newer QA Report or the merged-head proof represents the Run's Task or QA Report commits, archived Spec-directory work, or scoped uncommitted leftovers. Preserve it during dry-run; `--apply` can release it after revalidation. |
| `unintegrated` | Clean, resolved evidence proves that the Run Branch tip is not an ancestor of the target tip. Preserve the Run Worktree and Run Branch. |
| `dirty` | Tracked or untracked changes lack archived-Spec supersession proof, even after including paths beneath a declared or recorded path. For an archived Spec, `Run Worktree has tracked or untracked changes; <commit proof reason>` names the failed commit proof, or `Run Worktree has tracked or untracked changes; <path> is outside the archived Spec's declared and recorded paths` names the first outside path in byte order; the suffix is shortened in the middle to fit 160 bytes. An unarchived Spec keeps the bare constant; an unavailable archive source keeps its diagnostic. Preserve the Run Worktree and Run Branch. |
| `unknown` | Metadata or Git evidence cannot prove another state. Preserve every identified Run Worktree and Run Branch. |
| `released` | Both the Run Worktree and Run Branch are absent. No cleanup is needed. |

The same report can add debris candidates beside those legacy Run Worktree
classifications. Only a full scan (`roundfix reconcile` without a Run ID)
inspects item branches named `roundfix/deliver-<slug>-<16 hex>`. A live item's
branch is preserved with its Delivery Queue item. `--apply` deletes a proven
item branch with its clean worktree. The `roundfix-reconcile/v1` report's
`itemBranchCandidates` list names these entries, and `debrisSummary` reports
`itemBranchCandidates` and `itemBranchesApplied` counts.

For a Run of a merged Spec, reconciliation proves the Run against the merged
head. The Delivery Queue merge record is the primary source; when no usable
record remains, the default branch carrying the archived Spec is the fallback.
A Run is released only when each Run commit is represented by completed Task
status, a superseding QA Report, or matching changed content. If neither source
is usable or a commit is unrepresented, the Run remains preserved and its
reason names the missing proof; an absent target alone is not proof of release.

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

- A `process` candidate is proven when a terminal Run with a proven recorded
  owner identity still owns an inspected live process tree. Its report names
  the Run outcome, owner PID, every inspected process ID, and that ownership
  proof.
- A `runBranch` candidate is proven when set classification for one target
  branch and Spec shows that the Run Branch is superseded by a named current
  or target QA Report, and its registered Run Worktree was inspected clean.
- A `staging` candidate is any registered
  `roundfix-carry-forward-*/worktree`. Its `owner.json` PID and process identity
  prove it stale when the PID's `OwnerProcessIdentity` lookup fails or differs
  from the record. A legacy registration is stale only while it is
  `locked initializing`; an unlocked legacy staging worktree is preserved.
- A `runOwnerRecord` candidate is a Stale Run Owner Record: the live
  process at a terminal Run's recorded owner PID presents a start identity
  that differs from the recorded one, which proves the recorded owner ended.
  Its report names the Run outcome, owner PID, an empty process list and that
  proof. `--apply` clears the Run's recorded owner PID and identity, journals
  them, sends no signal, and counts the release in
  `debrisSummary.runOwnerRecordsReleased`; an identity that cannot be read or
  compared stays in `preservedCandidates`.

Ambiguous ownership, identity, Git, active-Run, or cleanliness evidence goes
to `preservedCandidates` with a refusal reason instead of becoming a cleanup
candidate.

The `roundfix-reconcile/v1` JSON report exposes every staging entry in
`stagingCandidates`. `debrisSummary.stagingCandidates` counts those entries and
`debrisSummary.stagingApplied` counts releases. Dry-run reports staging without
mutation. `--apply` releases stale staging, and `--carry-forward` releases it
before creating its own staging worktree. A live matching owner and every
case without positive stale proof remain in `preservedCandidates`.

After reviewing the dry-run, apply cleanup explicitly:

```bash
roundfix reconcile <run-id> --apply
roundfix reconcile --apply
roundfix reconcile <run-id> --discard-superseded
roundfix reconcile <run-id> --carry-forward
```

These three mutation switches are mutually exclusive. `--apply` releases
entries classified `safe` or `superseded`, or proven process, Run owner record,
Run Branch, and staging candidates, after rechecking the applicable metadata, ownership,
cleanliness, heads, ancestry, merged-head records, Task status, content, and
superseding QA Report evidence. `--discard-superseded`
records a Branch Disposition before removing a Run Branch proven superseded.
`--carry-forward` hands settled Tasks from one terminal spec Run back to the
checkout; it accepts Runs whose outcome is `BudgetExceeded`, `Stopped`, or
`Unresolved` and refuses every other terminal outcome. Carry-forward keeps its
existing proof requirements and refuses the whole Task set when any member
cannot be proved. A Task already completed on the checkout is reported as
`already completed; nothing to carry`; it is not staged or proved and never
refuses the remaining set. When every candidate is already completed, the
command succeeds without moving `HEAD`.
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
`CONTEXT.md`, the `implement-task` skill and every explicit Context path still
refuse on any change. This rule also serves Implement Preflight and Delivery
Retry; their refusal reasons and recovery commands are unchanged.

Each candidate lists tolerated planning amendments separately. Text prints these
lines in order (repository-relative paths replace the placeholders):

```text
  moved-inputs: []
  amended-inputs: [<paths>]
```

The `carryForwards` JSON entry includes `"movedInputs":[]` and
`"amendedInputs":["<path>"]`. `amendedInputs` is `[]` when there are no tolerated
amendments; those paths do not appear in `movedInputs` or set `inputsMoved`.

Carry-forward staging commits run without repository hooks because the carried commits already passed Daemon Verification and the repository hooks when the Daemon settled them. The checkout receives those commits only through a fast-forward merge, and carry-forward leaves its Git configuration unchanged.
There is no force bypass.

Process termination succeeds only when Roundfix proves every reported process
absent. An unprovable termination is reported with its reason and is never
treated as success; silence from the host does not mean the process stopped.

Never substitute manual Git deletion for this supported workflow. Do not run
`git worktree remove`, delete the Run Branch, or remove a recorded worktree
directory by hand. The Reconcile Command owns safety proof, evidence recording,
the guarded Integration Pending transition, and cleanup.
