### deliver

```bash
roundfix deliver plan [--json] [<slug>...]
roundfix deliver start [--max-duration <duration>] [--max-retries <n>] <slug>...
roundfix deliver status
roundfix deliver resume
roundfix deliver retry <slug>
roundfix deliver stop
```

Run `roundfix deliver plan` before recording a queue. With explicit slugs it
keeps their order; without slugs it reports every active Spec. Text output is
tab-separated:

- `spec` rows give the slug, `approved` or `blocked` verdict, unfinished and
  total Task counts, and authorization or strict Spec-check reasons.
- `shared` rows name production Go `interface:` paths a later Spec also shares
  with an earlier Spec in the plan.
- `backlog`, `finding`, and `inbox` rows list repository intent that is not
  approved to run. Backlog and finding rows carry their frontmatter status;
  inbox rows use `-`.

`--json` emits one `roundfix-deliver-plan/v1` document with the same Specs,
verdicts, reasons, shared premises, Task counts, and intent. Exit `0` means
every reported Spec is approved, exit `1` means at least one is blocked, and
exit `2` means usage or preflight failed. The plan opens no Run Database,
creates no worktree, and writes no file. It reports authority but never grants
implementation or delivery authority. A `shared` row predicts the later
item's `premise-changed` warning after the earlier item merges; it never blocks
the Spec or stops a queue.

Creates and operates a durable, ordered queue of Specs. Each item advances
from its Run to merge in this order: Run, pre-PR review, archive on the branch,
repository gate, push, pull request, current-head checks, and squash merge.
When the configured pre-PR review is `none`, the queue records that omission
and archives after the Run's QA gate.

A Spec can name prerequisite Specs in its Task Graph manifest (`_tasks.md`):

```yaml
requires:
  - 0192-first-prerequisite
  - 0193-second-prerequisite
```

`requires` is optional. It must be a list of distinct, non-empty strings and
cannot name the Spec itself. `deliver start` loads every queued Spec, then
refuses unknown prerequisites or cycles among queued Specs with exit `2`
before checking delivery authorization or recording a queue. A prerequisite
must exist in the Specs Root or its resolved archive root; it need not be in
the queue.

Before creating an item's worktree, the owner fetches the delivery remote's
default branch and reads `requires` from that refreshed branch. A prerequisite
is met only when its archived `_prd.md` exists there. An archive present only
in your checkout or an unmerged branch does not count. With no unmet
prerequisites, the item starts as usual.

When every unmet prerequisite is an item in the queue that is neither parked
nor merged, the dependent item waits at `queued`; the owner logs the wait and
moves on. Otherwise it parks as `prerequisite-unmerged: <slug>, …` without
creating a worktree. `deliver status` classifies this park as `dependency` and
points to the prerequisite's retry when that prerequisite is parked, or asks
you to deliver or merge it when it is outside the queue.

At the start of every owner pass and after any item merges, the owner returns
a dependency park to `queued` when all its prerequisites are met, clears the
blocker and leaves the retry count unchanged. When every item is parked, the
owner exits; retrying a prerequisite restarts it. You can also run
`roundfix deliver retry <slug>` on the dependent item to return it to `queued`
without a worktree, revalidation or carry-forward. The owner checks its
prerequisites again before starting it.

Each queued Spec runs in its own linked worktree under `worktree.location`,
created from the refreshed default branch.
`roundfix deliver` never switches, resets or cleans your checkout, and it does not need the checkout to be clean.
A parked item keeps its worktree, and `deliver status` prints that path. On
resume, Roundfix recreates a missing worktree from its recorded branch; if the
branch is missing too, it parks the item as `item-worktree-missing` instead of
replaying the stage. After an item merges, its worktree and local item branch
are removed. Before that removal, cleanup refreshes the default branch from the delivery remote.
Roundfix then releases every terminal Run of the merged Spec that it can prove
is represented at the recorded merged head. A Run it cannot prove stays in
place, and `deliver status` names the Run and its reason in the item's cleanup
warning.

A start requires every named Spec's committed authorization to grant
`implement`, `commit`, `push`, `pull_request`, and `merge`. If any Spec lacks
one of those operations, `deliver start` exits `2`, names every refused Spec
and its reasons, points to `roundfix deliver plan`, and records no queue. A
strict Spec-check finding appears in the plan but does not refuse start; the
queue revalidates that Spec against its own starting main.

A Spec that changes no Governed Path records `paths: []`; that explicit empty
list grants the listed operations and bounds no Governed Path.

Use `--max-duration <duration>` with a positive Go duration to set the queue
deadline, and use `--max-retries <n>` with an integer of at least `1` to limit
retries per item. Omitted limits are recorded as `none`. Start and status print
the recorded values as:

```text
Limits: deadline <RFC 3339 UTC|none>, retries per item <n|none>, concurrency 1, tokens <n|none>
```

`deliver start` prints only the limits line. `deliver status` prints usage
immediately after it, summing every Run recorded for queue items, including
Runs from earlier retries:

```text
Limits: deadline none, retries per item none, concurrency 1, tokens none
Usage: 5639755 tokens from 1 of 1 prompt(s) across 1 Run(s); cost not reported
```

With unreported prompts the line is `Usage: none reported by <n> prompt(s)
across <k> Run(s); cost not reported`. A queue with no linked Runs prints
`Usage: no Runs recorded`. Linked Runs with no usage rows print `no prompts
recorded` across their Run count. A Run in progress contributes once the queue
records it at Run end. Reported cost is grouped by currency; Roundfix computes
no prices. See [Token usage](../usage.md#token-usage).

At or after the deadline, the owner parks each item that has not started as
`queue-deadline`; an item that has started continues. A `queue-deadline` item
cannot be retried. Record a new queue for the remaining Specs instead. When an
item reaches its retry limit, `deliver retry` refuses the next retry and leaves
the item unchanged.

After the item worktree is created from that main and before the first Run,
Roundfix runs the strict Spec Consistency Check in the worktree. A finding
parks the item as `revalidation-failed: <code>, <code>` before any Run starts.

Revalidation also compares the Spec's declared production Go `interface:`
paths with the merge commits of earlier queue items. An overlap records
`premise-changed: <path>, <path> (merge <sha>, <sha>)` as a warning and the item
continues to its Run. `deliver status` prints `Warning: <slug> <warning>` after
the item rows, and the delivery console log prints `roundfix: warning: Delivery
Queue item <slug>: <warning>`. No overlap adds no warning or log line.

At that item-start boundary, Revalidation also compares the queue owner's build
commit with the starting main. When the owner predates a commit that changed
Roundfix source under `cmd/`, `internal/`, `go.mod`, or `go.sum`, the item adds
`owner-older-than-main: owner build <commit> predates starting main <commit>`
after any `premise-changed` warning. `deliver status` and the delivery console
log print the combined warning, and the item continues to its Run. A docs-only
change, a current owner, or a build commit absent from the repository adds no
owner warning. A retry keeps the recorded warning and does not recompute it.

A blocker parks its item with a reason and the queue continues with later
items. On resume, the owner reconciles every recorded action without a receipt
against observed state before retrying it, so a lost acknowledgement cannot
create a duplicate pull request or merge. Publication requires the Spec's
authorization record to grant `push`, `pull_request`, and `merge`.

A findings verdict with archived Specs parks as
`corrective-spec-required: <slug>[, <slug>]`; findings without archived Specs
still park as `review-findings`. `deliver status` prints either blocker. No Run
budget, corrective-Task ceiling, or queue grant authorizes the new corrective
Spec, and Roundfix never authors or starts it.

When an Implement Run ends `BudgetExceeded`, the queue parks its item as
`run-budget-exceeded` with that Run's ID. `roundfix deliver retry <slug>`
considers every terminal Implement Run of the item's Spec on the item branch,
newest first, together with the recorded Run, and carries each Run's remaining
settled Tasks before resuming the item.

When an unresolved Run's newest QA Report on its Run Branch is `partial`, has
no finding-blocked rows, and has environment-blocked rows beyond those blocked
only because no Pull Request is open yet, the item parks as
`qa-environment-partial`. A missing or unreadable report, a finding-blocked
partial, or a partial blocked only by the pre-PR row keeps `run-unresolved`.

For this environment park, status gives the recovery sequence:

```text
(cd <worktree> && roundfix reconcile <run-id> --carry-forward), satisfy the environment-blocked QA rows, run roundfix archive <slug> --qa-override --approval <source> --reason <text>, then run roundfix deliver retry <slug>
```

The archive override needs explicit user authority and preserves the actual QA
verdict and Task state. After an operator archives with `qa_override: true`,
retry accepts the moved item head only if it descends from the last candidate
commit, or from the recorded Run's starting head when no candidate exists. It
records that head as the candidate and resumes at `reviewing` without Task
Carry-Forward. The archive stage recognizes the reviewed Spec already in the
archive and proceeds to `gating` without another commit. Review, repository
gating, delivery authorization and required checks still apply.

Delivery authorization is read from the parent of the newest first-parent
commit that deleted the active Spec's `_prd.md`. A later operator commit does
not hide that pre-archive grant. Other moved archived heads remain refused
unless the item is a resolved
`pull-request-conflict` as described below.

GitHub's `CONFLICTING` mergeable state stops the check wait on its first read;
`UNKNOWN` stays pending. The owner starts from a clean item worktree at the
candidate head, fetches the default branch, and merges it without rebasing or
force-pushing. If a conflicted path is outside `delivery.derived_paths`, it
aborts the merge and parks `pull-request-conflict: <path>, …`, naming only the
undeclared conflict paths.

For a conflict confined to declared derived paths, the owner takes the default
branch's version of each conflicted file, then runs each matched regeneration
command once in declaration order. Declarations come from Project Config at
the fetched default-branch commit, with User Config beneath it; item-only
command changes are never run. A command that changes an undeclared path
aborts the merge and parks with `regenerated <path> outside
delivery.derived_paths`. Otherwise the owner commits the merge with the
`Roundfix-Delivery: derived-merge` trailer and returns the item to `gating`.
The repository gate, push and current-head checks run again. An existing Pull
Request still reporting an earlier candidate is read again at each check
interval up to the check timeout.

A conflict park has class `conflict`. Its next action is to merge the default
branch into the item branch in the printed worktree, resolve the named paths,
commit, then run `roundfix deliver retry <slug>`. Retry accepts a head descended
from the candidate, records it, and resumes at `reviewing`.

`deliver status` prints each item's Spec slug, stage, blocker and worktree.
After any `Warning:` lines and before `Limits:`, it prints one line per parked
item in queue order:

```text
Park: <slug> <class>: <next command>
```

The Park Class identifies the reason and the next action:

| Class | Blockers |
| --- | --- |
| `dependency` | `prerequisite-unmerged` |
| `conflict` | `pull-request-conflict` |
| `environment` | `qa-environment-partial`, `checks-timeout`, `item-worktree-missing`, `delivery-error` |
| `flaky-check` | a check that failed again outside the item's changed packages |
| `finding` | `run-unresolved`, `review-findings`, `corrective-spec-required`, `gate-failed`, `checks-failed`, `revalidation-failed` |
| `budget` | `run-budget-exceeded`, `queue-deadline` |
| `review` | `review-blocked`, `review-stale` |
| `authorization` | `unauthorized` |
| `unclassified` | an unknown blocker |

Each existing blocker keeps its previous next action. A queue without parked
items adds no `Park:` line.

When a required GitHub Actions check fails on its first attempt, the owner
reads its failed log and compares the failing Go package directories with the
item's changed paths against the refreshed remote default branch. If every
failing package lies outside that change, it re-runs the failed jobs once,
logs the action, restarts the check timeout once and keeps polling. A pass adds
`flaky-check: <check> passed on re-run` to the item's Warning. A second failure
outside the change parks as `flaky-check: <package>, …`, with an action to fix
or re-run the check and then run `roundfix deliver retry <slug>`.

A failure in a changed package, a build or setup failure, a log without Go
packages, a non-Actions check, or a run already past attempt one parks as
`checks-failed` without a re-run. An inspection or re-run error also parks as
`checks-failed` and is recorded in the owner log.

When one or more items are parked, status prints exactly one
`Pending question:` for the lowest-position parked item, the action that
answers it from the same Park Class table, and the count waiting behind it.
A dependency park can clear after its prerequisites merge. Other parks need
`deliver retry` or a new queue; elapsed time alone does not resolve them.
`deliver stop` ends the detached owner; `deliver resume` restarts it from the
persisted queue.

`roundfix deliver retry <slug>` returns one parked item to the queue. For an
active Spec that has not run, it first repeats the strict check in the item
worktree and refuses while findings remain, leaving the item unchanged. It
then carries the remaining settled Tasks from every terminal Implement Run of
the item's Spec on the item branch, newest first, so a later Run executes only
unfinished Tasks. A retry does not change a recorded `premise-changed` warning.
It also does not change or recompute a recorded `owner-older-than-main` warning.
It selects the re-entry stage from the evidence on that branch:

When every refused Task has moved inputs and only non-Task commits after the
Run started changed those inputs, the refusal adds `amended by <sha>, ...` to
the reason. Its next action prints these five commands with the item worktree,
Run ID, amendment commits, and Spec slug filled in and POSIX-quoted:

```bash
git -C '<worktree>' branch 'roundfix-amended-<run-id>' HEAD
git -C '<worktree>' reset --hard '<first-amendment>^'
(cd '<worktree>' && roundfix reconcile '<run-id>' --carry-forward)
git -C '<worktree>' cherry-pick '<amendment>' ...
roundfix deliver retry '<slug>'
```

Roundfix prints these commands and never runs them. If a Task commit changed a
moved input, or the refusal includes another cause, the existing single
`roundfix reconcile <run-id> --carry-forward` next action remains unchanged.

A retried `review-findings` item at an unchanged head advances once every
finding is dismissed with evidence. Standing findings park it again without
asking the reviewer; Roundfix asks the reviewer again only after the head
changes.

| Recorded evidence | Re-entry stage |
| --- | --- |
| Active Spec with any unfinished Task | `running` |
| Active Spec with every Task completed | `reviewing` |
| Operator-archived `qa-environment-partial` with a QA override and head descended from its candidate or Run start | `reviewing` |
| Resolved `pull-request-conflict` with a head descended from the candidate | `reviewing` |
| Archived Spec with unchanged candidate and no recorded pull request | `gating` |
| Archived Spec with unchanged candidate and a recorded pull request | `checking` |
| `corrective-spec-required` with the parked candidate head unchanged | `reviewing`, without Task Carry-Forward |
| `corrective-spec-required` after the item head moved | Refused with exit `2`; the item stays unchanged and the operator must author a corrective Spec with its own authorization and QA gate |

After the retry, a live owner whose identity Roundfix proves keeps the queue
and stdout reports `Handed <slug> to Delivery Queue owner PID <pid>.`. If the
recorded owner is dead or its identity is unproven, Roundfix reclaims the owner
record, writes the same stderr notice as `deliver resume`, and starts a new
detached owner. With no recorded owner, it starts one directly.

A successful retry exits `0`. When Tasks moved, stdout prints one `Carried
forward from Run <run-id>: <task>, <task>` line per Run carried from, in
newest-first order, before `Retried <slug>: <blocker> -> <stage>`. A missing or
empty slug, an extra argument, an unknown flag, an item that is not parked, a
missing item branch, a moved archived head without accepted recovery evidence,
a refused carry-forward, or an owner hand-off failure exits `2` and starts no owner. An item-level refusal
leaves the stored item unchanged.
