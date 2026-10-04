### deliver

```bash
roundfix deliver plan [--json] [<slug>...]
roundfix deliver start [--max-duration <duration>] [--max-retries <n>] [--max-tokens <n>] <slug>...
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

Each queued Spec runs in its own linked worktree under `worktree.location`.
After the prerequisite check and before authorization or readiness checks,
`deliver start` reads local item branches named
`roundfix/deliver-<slug>-<16 lowercase hex digits>`. It compares them with the
local `<remote>/<default>` ref without fetching. With exactly one branch
holding commits that ref lacks, it records the new queue and prints on stdout:

```text
Continuing item branch <branch> for <slug>
```

The owner checks again after fetching the default branch and records that
branch on the new item, reusing its worktree and completed work. The item
starts at `queued`, as any new item does; Roundfix does not merge the default
branch into it at start. With no branch holding work, it creates a new item
branch from the refreshed default branch.

With two or more item branches with commits the default branch lacks, start
exits `2` before recording a queue. The reason names every branch in sorted
order:

```text
Spec "0300-example" has 2 item branches with commits origin/main lacks: roundfix/deliver-0300-example-1111111111111111, roundfix/deliver-0300-example-2222222222222222; delete every branch but the one to continue, then run roundfix deliver start again
```

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

After checking delivery authorization, `deliver start` checks the `gh` and
`remote` readiness lines before opening the Run Database. Any `failed` finding
refuses with exit `2`: the `Preflight failed` block says `Delivery Queue
cannot publish from this machine:` and lists one reason per finding, including
its `DR-` code and next action. It creates no Delivery Queue and starts no
owner. A `warn`, such as `DR-GH-UNREACHABLE` on an offline machine, is silent
and does not refuse start.

For example, a missing GitHub login produces:

```text
Preflight failed

Reason:
  Delivery Queue cannot publish from this machine:
  gh: DR-GH-UNAUTHENTICATED: gh has no account for github.com; next: gh auth login --hostname github.com

No side effects:
  Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.

Usage:
  Run 'roundfix deliver start --help' for usage.
```

A Spec that changes no Governed Path records `paths: []`; that explicit empty
list grants the listed operations and bounds no Governed Path.

Use `--max-duration <duration>` with a positive Go duration to set the queue
deadline, and use `--max-retries <n>` with an integer of at least `1` to limit
retries per item. Omitted limits are recorded as `none`. Start and status print
the recorded values as:

```text
Limits: deadline <RFC 3339 UTC|none>, retries per item <n|none>, concurrency 1, tokens <n|none>
```

`deliver start` prints the limits line after any continuation lines.
`deliver status` prints usage immediately after it, summing every Run recorded for queue items, including
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

After the item worktree is created or continued and before the first Run,
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

A repository may declare `delivery.item_binary` in Project Config with a
`build` command and a repository-relative `path`. The queue owner reads the
declaration it loaded at start. Before each `implement`, `archive` and `review`
step, it runs the build in the item worktree and writes the build output to
`<artifact dir>/delivery/<slug>/item-binary-build.log`. It then asks the built
binary to run `migrate --check`. When that exits `0`, the step runs the item
binary and the owner's console log contains:

```text
roundfix: Delivery Queue item <slug>: <step> runs the item binary <absolute path>
```

When `migrate --check` exits with any other code, the step runs the owner's
binary and the console log contains:

```text
roundfix: notice: Delivery Queue item <slug>: <step> runs the owner's binary; the item binary's migrate --check exited <n>: <first non-empty line of its stderr, else stdout>
```

A failed build, a declared path Git does not ignore, or an item binary that
cannot start parks the item as `delivery-error`; the blocker names the command
or path and the build log where applicable. A repository without
`delivery.item_binary`, including every repository that installs Roundfix from
npm, is delivered as before with the owner's binary.

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
no finding-blocked rows, and has at least one environment-blocked row, including
a row blocked only because no Pull Request is open yet, the item parks as
`qa-environment-partial`. A missing or unreadable report, a finding-blocked
partial, or a partial with no environment-blocked row keeps `run-unresolved`.

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

A retry after a correction committed on top of an archived candidate accepts
its current head when Git proves it descends from the newest candidate. It
appends the head to the candidate commits and returns to `reviewing`, so the
correction receives a fresh review before the archive and repository gate
stages. This applies to an archived `gate-failed` item and other blockers,
with two restrictions: `qa-environment-partial` still needs the recorded QA
Archive Override, and `corrective-spec-required` still refuses a moved head.
When no candidate is recorded, only an item the operator archived with the QA
Archive Override may use the Run start head, whatever its park. The retry
records the descended head as the candidate and resumes at `reviewing`. A
non-descendant head or unavailable item history refuses with the existing
reason and leaves the item unchanged. An unchanged archived candidate keeps
its `gating` stage without a Pull Request or `checking` stage with one.

An archived retry of an operator-archived item finds
the Implement start head of the Run the queue started by the repository the Run
belongs to. When no candidate exists, the retry accepts an item head descended
from that start head, records it as the candidate and resumes at `reviewing`.
An archived item with an unchanged candidate head and a recorded Pull Request
resumes at `checking`; this includes a `delivery-error` park, so a green Pull
Request can continue to merge without manual intervention.

Delivery authorization is read from the parent of the newest first-parent
commit that deleted the active Spec's `_prd.md`. A later operator commit does
not hide that pre-archive grant. A moved archived head that does not descend
from its newest candidate remains refused.

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

During `checking`, Roundfix waits while GitHub reports the Pull Request merge
state as `BLOCKED` or `UNKNOWN`, even when the listed checks pass. The existing
checks timeout parks the item as `checks-timeout`; only a mergeable state with
passing checks proceeds to merge. If GitHub refuses that merge because
`base branch policy prohibits the merge`, Roundfix returns to `checking` once
for that head. A second refusal for the same head parks `delivery-error`.

A conflict park has class `conflict`. Its next action is to merge the default
branch into the item branch in the printed worktree, resolve the named paths,
commit, then run `roundfix deliver retry <slug>`. Retry accepts a head descended
from the candidate, records it, and resumes at `reviewing`.

`deliver status` prints each item's Spec slug, stage, blocker and worktree.
After any `Warning:` lines and before `Limits:`, it prints one line per parked
item in queue order, except `queue-token-ceiling` (see Token usage below):

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
| Archived Spec with a post-archive correction descended from its newest candidate | `reviewing` |
| Operator-archived item with a QA override and head descended from its candidate or Run start | `reviewing` |
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

An item-level refusal exits `2` and prints `Retry refused`, followed by
`Reason:`, the refused retry reason, `Item:` with its stage and blocker, and
`No side effects:` with the statement that Roundfix did not change the Delivery
Queue item, start a queue owner, commit, or push. It does not print a `Usage:`
block. For example:

```text
$ roundfix deliver retry 0300-example
stdout:
stderr:
Retry refused

Reason:
  retry Delivery Queue item "0300-example": archived item head "2222222222222222222222222222222222222222" differs from candidate head "1111111111111111111111111111111111111111"

Item:
  stage: parked; blocker: delivery-error: merge pull request: read pull request before merge: gh failed

No side effects:
  Roundfix did not change the Delivery Queue item, start a queue owner, commit, or push.
exit: 2
```

### Token usage

`deliver start --max-tokens 5000000 <slug>...` records an optional queue token
ceiling. The value must be an integer of at least 1; an invalid value exits
`2` before recording a queue or creating a Run Database. The limits line ends
with `tokens 5000000`, or `tokens none` when the flag is omitted.

`deliver status` prints a `Usage:` line with the tokens and adapter-reported
cost of every Run linked to the queue, including earlier retries. Totals say
how many prompts reported usage; unreported prompts add nothing. Cost is
reported by currency and is never calculated from token prices.

At or above the ceiling, each item still in `queued` parks as
`queue-token-ceiling` before creating a branch or worktree. Its Pending
Question answer is:

```text
record a new queue for the remaining Specs with roundfix deliver start and a higher --max-tokens
```

This blocker is presented in the item row and Pending Question without a
separate `Park:` line. Any parked item's retry is refused while the queue's
recorded tokens remain at or above its ceiling, with exit `2`, no workspace
action and no item change. For example, the retry reason is:

```text
retry Delivery Queue item "0301-example": queue token ceiling 5000000 was reached with 5639755 tokens; start a new queue with roundfix deliver start
```

Record a new queue for the remaining Specs with a higher ceiling or omit the
flag. Items already past `queued`, including a running item, continue; no Run
is signalled, stopped or cancelled. The queue can exceed its ceiling by an
item's Run. Tokens outside Runs, such as the pre-PR review, do not count.
