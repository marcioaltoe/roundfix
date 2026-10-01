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
| `superseded` | A newer QA Report or the merged-head proof represents the Run's Task or QA Report commits. Preserve it during dry-run; `--apply` can release it after revalidation. |
| `unintegrated` | Clean, resolved evidence proves that the Run Branch tip is not an ancestor of the target tip. Preserve the Run Worktree and Run Branch. |
| `dirty` | A present registered Run Worktree has tracked or untracked changes. Preserve the Run Worktree and Run Branch. |
| `unknown` | Metadata or Git evidence cannot prove another state. Preserve every identified Run Worktree and Run Branch. |
| `released` | Both the Run Worktree and Run Branch are absent. No cleanup is needed. |

The same report can add three debris candidate kinds beside those legacy Run
Worktree classifications:

For a Run of a merged Spec, reconciliation proves the Run against the merged
head. The Delivery Queue merge record is the primary source; when no usable
record remains, the default branch carrying the archived Spec is the fallback.
A Run is released only when each Run commit is represented by completed Task
status, a superseding QA Report, or matching changed content. If neither source
is usable or a commit is unrepresented, the Run remains preserved and its
reason names the missing proof; an absent target alone is not proof of release.

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
entries classified `safe` or `superseded`, or proven process, Run Branch, and
staging candidates, after rechecking the applicable metadata, ownership,
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
Carry-forward staging commits run without repository hooks because the carried commits already passed Daemon Verification and the repository hooks when the Daemon settled them. The checkout receives those commits only through a fast-forward merge, and carry-forward leaves its Git configuration unchanged.
There is no force bypass.

Process termination succeeds only when Roundfix proves every reported process
absent. An unprovable termination is reported with its reason and is never
treated as success; silence from the host does not mean the process stopped.

Never substitute manual Git deletion for this supported workflow. Do not run
`git worktree remove`, delete the Run Branch, or remove a recorded worktree
directory by hand. The Reconcile Command owns safety proof, evidence recording,
the guarded Integration Pending transition, and cleanup.

