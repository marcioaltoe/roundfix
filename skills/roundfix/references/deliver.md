## Delivery queue

Use the durable delivery queue when a sequence of Specs must outlive the
terminal session:

```bash
roundfix deliver plan [--json] [<slug>...]
roundfix deliver start [--max-duration <duration>] [--max-retries <n>] <slug>...
roundfix deliver status
roundfix deliver resume
roundfix deliver retry <slug>
roundfix deliver stop
```

Run `roundfix deliver plan` before recording a queue. With explicit slugs it
keeps their order; without slugs it reports every active Spec. Its
tab-separated output uses `spec` rows for the verdict, unfinished and total
Task counts, and authorization or strict Spec-check reasons; `shared` rows for
production Go `interface:` paths a later Spec shares with an earlier one; and
`backlog`, `finding`, and `inbox` rows for repository intent that is not
approved to run. Backlog and finding rows carry frontmatter status; inbox rows
use `-`.

`--json` emits one `roundfix-deliver-plan/v1` document with the same Specs,
verdicts, reasons, shared premises, Task counts, and intent. Exit `0` means
every reported Spec is approved, exit `1` means at least one is blocked, and
exit `2` means usage or preflight failed. The plan opens no Run Database,
creates no worktree, and writes no file. It reports authority but never grants
implementation or delivery authority. A `shared` row predicts the later
item's `premise-changed` warning after the earlier item merges; it never blocks
the Spec or stops a queue.

For each queued Spec, `roundfix deliver` advances from the Run to merge in
this order: Run, pre-PR review, archive on the branch, repository gate, push,
pull request, current-head checks, and squash merge. With pre-PR review set to
`none`, the queue records the configured omission and proceeds to archive after
the Run's QA gate.

A Spec's optional `_tasks.md` frontmatter `requires` lists prerequisite Specs.
Start refuses unknown slugs, self-references, duplicates and queued cycles
before recording a queue. Before creating a worktree, the owner reads the
manifest from the fetched default branch and requires each prerequisite's
archived `_prd.md` there. Unmerged local archives do not count. If every unmet
prerequisite is queued and neither parked nor merged, the item waits at
`queued`; otherwise it parks `prerequisite-unmerged: <slug>, …` without a
worktree. Each owner pass and each merge returns satisfied dependency parks to
`queued` without increasing retries. Retrying a dependency park also targets
`queued`, without worktree recovery or Task Carry-Forward.

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
Limits: deadline <RFC 3339 UTC|none>, retries per item <n|none>, concurrency 1, spend not measured
```

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

A blocker parks its item with a reason; it does not stop later queued items.
When a queue resumes, it reconciles every recorded action that lacks a receipt
against the observed remote state before retrying that action. This prevents a
lost acknowledgement from creating a second pull request or merge. Before
publication, the Spec authorization record must grant all three operations:
`push`, `pull_request`, and `merge`.

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


When a required Actions check fails on attempt one, the owner reads its Go
package summaries and compares them with the item's change against the
refreshed default branch. A failure entirely outside that change re-runs the
failed jobs once and restarts the timeout once. A pass records
`flaky-check: <check> passed on re-run` in Warning; a second failure outside the
change parks `flaky-check: <package>, …`. Failures in changed packages,
build/setup failures, unattributable logs, and attempts past one park
`checks-failed` without a re-run.

An unresolved Run parks `qa-environment-partial` when its newest Run-Branch QA
Report is partial, has no finding-blocked rows, and has at least one
environment-blocked row, including a row blocked only because no Pull Request
is open yet. Otherwise it keeps `run-unresolved`.
Status gives the carry-forward, environment repair, authorized archive override
and retry sequence. Retry of an operator-archived item requires
`qa_override: true` and a head descended from the last candidate, or the Run's
starting head when no candidate exists. It records that head and targets
`reviewing`. Archive recognizes an already archived reviewed Spec and advances
to `gating` without a new commit. The override preserves QA and Task state;
review, gate, authorization and checks still apply.

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

`deliver status` prints item rows and Warning lines, then one
`Park: <slug> <class>: <next command>` line per parked item in queue order,
before `Limits:`. Status and the Pending Question use the same Park Classes:

| Class | Blockers or action |
| --- | --- |
| `dependency` | `prerequisite-unmerged`; retry a parked prerequisite, or deliver/merge an external prerequisite and retry the dependent item |
| `conflict` | `pull-request-conflict`; merge, resolve, commit, then retry |
| `environment` | `qa-environment-partial`, `checks-timeout`, `item-worktree-missing`, `delivery-error`; the QA partial prints carry-forward, repair, authorized override and retry |
| `flaky-check` | `flaky-check`; fix or re-run the failing packages, then retry |
| `finding` | `run-unresolved`, `review-findings`, `corrective-spec-required`, `gate-failed`, `checks-failed`, `revalidation-failed` |
| `budget` | `run-budget-exceeded`, `queue-deadline` |
| `review` | `review-blocked`, `review-stale` |
| `authorization` | `unauthorized` |
| `unclassified` | unknown blockers; resolve the blocker, then retry |

Existing blockers keep their next actions. A queue without parked items adds
no Park line.
When one or more items are parked, it prints exactly one `Pending question:`
for the lowest-position parked item, the action that answers it, and the count
waiting behind it. A dependency park can clear after its prerequisites merge. Other parks need
`deliver retry` or a new queue; elapsed time alone does not resolve them.

Use `roundfix deliver retry <slug>` to return one parked item to the queue.
For an active Spec that has not run, Roundfix first repeats the strict check in
the item worktree and refuses while findings remain, leaving the item
unchanged. It then carries the remaining settled Tasks from every terminal
Implement Run of the item's Spec on the item branch, newest first. A retry does
not change a recorded `premise-changed` or `owner-older-than-main` warning. The
item re-enters at `running` when any Task is unfinished or at `reviewing` when
every Task is completed. An archived Spec with an unchanged candidate
re-enters at `gating` without a recorded pull request or at `checking` with one.

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

When every refused Task has moved inputs and only non-Task commits after the
Run started changed those inputs, the reason adds `amended by <sha>, ...` and
the next action prints these five POSIX-quoted commands in order:

```bash
git -C '<worktree>' branch 'roundfix-amended-<run-id>' HEAD
git -C '<worktree>' reset --hard '<first-amendment>^'
(cd '<worktree>' && roundfix reconcile '<run-id>' --carry-forward)
git -C '<worktree>' cherry-pick '<amendment>' ...
roundfix deliver retry '<slug>'
```

Roundfix prints the recovery and never runs it. A Task commit that changed a
moved input, or a refusal with another cause, keeps the existing single
`roundfix reconcile <run-id> --carry-forward` next action.

| Recorded evidence | Re-entry stage |
| --- | --- |
| Active Spec with any unfinished Task | `running` |
| Active Spec with every Task completed | `reviewing` |
| Archived Spec with a post-archive correction descended from its newest candidate | `reviewing` |
| Operator-archived item with override and proven ancestry | `reviewing` |
| Resolved `pull-request-conflict` with proven candidate ancestry | `reviewing` |
| Archived Spec with unchanged candidate and no recorded pull request | `gating` |
| Archived Spec with unchanged candidate and a recorded pull request | `checking` |
| `corrective-spec-required` with the parked candidate head unchanged | `reviewing`, without Task Carry-Forward |
| `corrective-spec-required` after the item head moved | Refused with exit `2`; the item stays unchanged and the operator must author a corrective Spec with its own authorization and QA gate |

A retried `review-findings` item at an unchanged head advances once every
finding is dismissed with evidence. Standing findings park it again without
asking the reviewer; Roundfix asks the reviewer again only after the head
changes.

The retry hands the item to the recorded owner only when Roundfix proves that
process is alive and has the recorded identity. A dead or unproven owner record
is reclaimed with a stderr notice and replaced; with no owner, Roundfix starts
a detached owner. Success exits `0`, prints one `Carried forward from Run
<run-id>: <task>, <task>` line per Run carried from, then prints `Retried
<slug>: <blocker> -> <stage>` and the live-owner hand-off or detached-owner
report. Invalid arguments, an item that is not parked, a missing item branch,
a moved archived head without accepted recovery evidence, refused carry-forward,
or owner hand-off failure exits
`2` and starts no owner; an item-level refusal leaves the item unchanged.

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

Use `roundfix upgrade [--check]` to resolve the latest Roundfix release through
the GitHub CLI. Without `--check`, it downloads the platform asset, verifies
size and checksum when present, and atomically replaces the current executable.
`--check` reports without installing. stdout outcomes are deterministic:

```text
upgraded 1.0.0 → 1.1.0
already current 1.0.0
no releases published
upgrade available 1.0.0 → 1.1.0
```

Failures exit `1`, leave the current binary untouched, and print a manual
fallback on stderr. Operational commands (`fetch`, `resolve`, `watch`, and
`implement`) run a best-effort daily freshness check. When the installed
version is behind, stderr contains exactly one line shaped like:

```text
roundfix 1.0.0 is behind latest 1.1.0; run roundfix upgrade
```

Freshness failures and offline checks stay silent and do not change the Run
outcome.

### Recommendation check

Doctor prints `recommendations:` immediately after `profiles:` from the
configuration it already loaded. It reports `ok` when no category differs,
`found` with counts and `run roundfix profiles check` otherwise, or `skipped`
with the reason when comparison cannot run. This line opens no Agent Session
and never fails Doctor.

`roundfix upgrade [--check]` writes a recommendation notice to standard error
after every successful release outcome, leaving standard output and the exit
code unchanged. After an install, the installed executable runs `profiles
check` in the process working directory under a ten-second timeout; otherwise
the running executable compares in process. A failed comparison prints only
`roundfix: recommendations not checked: <reason>`. Help, usage errors and
failed upgrades print no notice. An upgrade performed by an older executable
prints none; the notice starts with a subsequent upgrade.

`roundfix profiles check [--json]` compares every configured Agent Work
Category with the Recommended Profile in the shipped snapshot, including the
Preferred Selection and the full ordered Fallback Chain. It is read-only and
offline, opens no Agent Session, and writes nothing. Undefined optional
categories are omitted. The three statuses are `current` (the profile equals
the recommendation, even with a deviation), `differs` (a difference without a
deviation for the shipped snapshot), and `pinned` (a difference whose deviation
names that snapshot). An older deviation remains visible as `differs`, with
its original date. Text prints differences and pins, then the counts; a fully
current configuration prints only the summary. `--json` uses
`roundfix/profiles-check/v1`, with `snapshot`, `current`, `differ`, `pinned`
and every configured category in order. Each row carries `category`, `status`,
`source`, `configured`, `recommended`, and an optional `deviation`.
The command exits `0` even with differences and `2` for usage or configuration
load errors. `profiles show` adds `Recommendation status` and any declared
`Deviation`, plus JSON fields `recommendation_status` and `deviation`, under
its existing `roundfix/profiles/v2` schema; an undefined optional category's
status is `inherited`.

`roundfix profiles check --apply --scope user|project [--dry-run] [--yes]
[--json]` adopts only differing categories, writing each complete Recommended
Profile without a deviation. Current and pinned categories stay unchanged.
`--scope` is required; `--scope`, `--dry-run` and `--yes` require `--apply`.
With `--scope user`, categories supplied by Project Config are skipped and
named on standard error with advice to use `--scope project`. Adoption opens
disposable Agent Sessions to prove every exact tuple before confirmation or
writing, using the same preview, output and exit codes as `profiles configure`.
`--dry-run` proves and previews without writing; `--yes` skips confirmation.
Failed proof and declined confirmation leave all configuration bytes unchanged.
`--json` uses `roundfix/profiles-configure/v1`. With nothing to adopt, nothing
is prepared, proved or written: text prints `Profile configuration unchanged:
nothing to adopt`, JSON has `changed: false` and empty `profiles` and `changes`,
and the command exits `0`.

A configured Agent Selection Profile can carry a Profile Deviation under
`deviation`, with `from` (the snapshot calendar date, `YYYY-MM-DD`) and `reason`
(a non-empty string after trimming). It records that the profile's difference
from the Recommended Profile is deliberate for that snapshot. `profiles
configure` writes a deviation its fragment carries and removes an old deviation
when the replacement fragment omits it. A Roundfix older than this release
refuses a configuration that uses the key.


## Run Window

```bash
roundfix window set <HH:MM|YYYY-MM-DDTHH:MM> [--force]
roundfix window show
roundfix window clear
```

These commands read the current repository's durable Run Window in the Run
Database. The window bounds when an Implement Run may start;
`budget.max_run_duration` bounds how long it may run after starting, with its
allowance renewed at each Task settlement. The window does not apply to
`fetch`, `resolve`, or `watch`.

`roundfix window set` stores the next occurrence of a local `HH:MM` (tomorrow
if it has passed today), or an absolute local `YYYY-MM-DDTHH:MM` cutoff. A
past absolute instant is refused with exit `2`. An existing window is reported
and preserved unless `--force` replaces it. `roundfix window show` writes
nothing and prints the repository, cutoff, current time, and remaining duration,
or reports that no window is set; either state exits `0`.
`roundfix window clear` removes the stored window and reports whether one was
set.


### Token usage

`roundfix deliver status` prints `Usage:` across every Run linked to the queue,
including earlier retries: tokens, reporting prompt count, Run count and
adapter-reported cost by currency. Unreported prompts add nothing and remain
visible in the coverage count. Roundfix computes no spend from token prices.

`roundfix deliver start --max-tokens <n> <slug>...` accepts an integer of at
least 1 and stores the queue ceiling. Invalid values exit `2` before any
queue or Run Database is created. The limits line prints `tokens <n>` or
`tokens none`. At or above the recorded total, each queued item parks as
`queue-token-ceiling` before branch or worktree creation; any parked item's
retry is refused before workspace actions. Its Pending Question answer is
`record a new queue for the remaining Specs with roundfix deliver start and a
higher --max-tokens`. Status presents this blocker without a separate `Park:`
line. A retry refusal names the ceiling and recorded tokens and tells the
operator to start a new queue with `roundfix deliver start`.

Record a new queue for the remaining Specs with a higher ceiling or no flag.
Items already past `queued` continue, including running items; no Run is
signalled, stopped or cancelled. A queue may exceed its ceiling by an item's
Run. Unreported prompts and tokens outside Runs do not count.
