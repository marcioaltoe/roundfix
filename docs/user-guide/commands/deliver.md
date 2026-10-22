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
repository gate, push, pull request, Pull Request Record Commit, current-head
checks, and squash merge.
When the configured pre-PR review is `none`, the queue records that omission
and archives after the Run's QA gate.

After it opens the pull request and before it judges checks, the Delivery Queue
records the pull request number in the Archive Record in a Pull Request Record
Commit. The commit changes the record's `pull_request` field and the queue
pushes it before judging the pull request checks.

Its subject is `docs: record pull request #<number> in the <spec number>
Archive Record`, where the Spec number is the run of digits that starts the
slug, such as `0088`; a slug that does not start with digits followed by `-`
gives `docs: record pull request #<number> in the Archive Record`. The
subject stays within 100 characters for every slug. The queue accepts an
existing record commit by its change, whatever its subject.

The archive commit's subject is `docs: archive <slug>` while that is at most
100 characters, and otherwise `docs: archive the <spec number> Spec`, or
`docs: archive the Spec` for a slug without a Spec number. The pull request
title is `feat: deliver <slug>` while that is at most 91 characters, which
leaves room for the ` (#<number>)` a squash merge appends for any pull request
number below 100000; otherwise it is `feat: deliver <spec number>: <Spec
title>`, or `feat: deliver <Spec title>` for a slug without a Spec number,
where the Spec title is the Archive Record's title cut at a word to fit. Lengths
count UTF-16 code units, as commitlint does. The pull request body still names
the full slug, and the queue accepts an existing archive commit by its change,
whatever its subject.

When a commit hook refuses a commit the queue makes (the archive commit, the
Pull Request Record Commit or the derived merge commit), the item parks as
`delivery-error` and the reason names the hook, the refused subject, the
hook's exit status and output, and the recovery. The queue never bypasses the
hook. To recover, commit the staged change in the item worktree with a subject
the hook accepts, then run `roundfix deliver retry <slug>`. A refused derived
merge is aborted, so merge the default branch into the item branch by hand as
for a conflict park, commit it with a message the hook accepts, then retry. A
hook that does not name itself is named from the installed `pre-commit` and
`commit-msg` hooks.

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

`deliver start` refuses a queue whose Spec requires a cancelled Spec, and
`deliver retry` refuses a cancelled Spec.

Before creating an item's worktree, the owner fetches the delivery remote's
default branch and reads `requires` from that refreshed branch. A prerequisite
is met only when its Archive Record or, for a legacy archive, its archived `_prd.md` exists there. An archive present only
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
warning. This cleanup proves a squash merge onto a default branch that moved
during the item by comparing the merge commit with Git's merge of the candidate
into its first parent.

The archive commit also carries the promoted copies and every relinked file
named by the Archive Record. The exact-archive proof recomputes each promoted
copy from the single removed file with its basename, applying the archive's
Git Locator rewrite with the parent as the preferred commit. The copy must
match those rewritten bytes; an ambiguous basename or any extra edit refuses
the proof. The queue accepts a relinked file only when it is
outside the archived folder and History Root, existed at the parent, and its
bytes equal the archive's rewrite of the parent's bytes under the same citation
moves; an extra or altered relinked file fails the archive proof.

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
`<artifact dir>/delivery/<slug>/item-binary-build.log`. After a successful build,
it copies the built regular file to
`<artifact dir>/delivery/<slug>/item-binary/<first 16 hex of SHA-256>/<base name>`.
It writes a temporary file in that directory, sets mode `0755`, closes it and
renames it over the final name, including when identical bytes reuse the same
hash path. A later build can write its output in place without rewriting a
running copy. Copies accumulate under the item's artifacts, one per distinct
build. After releasing a merged item's Runs, the Delivery Queue removes its
`item-binary` directory and keeps its other artifacts, including
`item-binary-build.log` and `repository-gate.log`. A removal failure joins the
item's cleanup warning. The GC Command also reclaims Item Binary directories
for items that are merged or no longer queued; every other stage keeps them.
Use `roundfix gc --dry-run` to preview that sweep and its reclaimable paths.

The queue probes the copy with `migrate --check`. When that exits `0`, the step
runs that copy and the owner's console log names its absolute path:

```text
roundfix: Delivery Queue item <slug>: <step> runs the item binary <absolute copy path>
```

When `migrate --check` exits with any other code or is killed by a signal, the step runs the owner's
binary and the console log contains:

```text
roundfix: notice: Delivery Queue item <slug>: <step> runs the owner's binary; the item binary's migrate --check <outcome>: <first non-empty line of its stderr, else stdout>
```

`<outcome>` is `exited <n>` for an exited probe, or, for a Signal Exit,
`was killed by signal <name> (<number>), core dumped` or
`was killed by signal <name> (<number>), no core dump`. When an item or owner
step child (`implement`, `archive`, or `review`) ends by a signal, the step
fails with `<operation> was killed by signal <name> (<number>)` and its core
state, followed by `: <detail>` only when the child wrote a non-empty line.

A failed build, a declared path Git does not ignore, a non-regular build
output, a copy failure, or an item binary that cannot start parks the item as
`delivery-error`; the blocker names the command or path and the build log where
applicable. A read or copy failure is reported as `copy item binary: <error>`. A repository without
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
Archive Override. A moved `corrective-spec-required` head returns to `reviewing`
when API Contract 5 and API Contract 6 prove the correction answers only the
review recorded in the Reviewed Head Record of the parked candidate in the
item worktree's resolved checkout. It descends from the parked candidate,
every standing finding has exactly one valid disposition, and every changed
path lies under an archived Spec the blocker names or among the paths its Tasks
declared (`creates:`, `interface:` and recorded paths read at the recorded
revision). One exception admits the fix whatever paths it touched: the parked
record is round 2, is not itself a **Correction Review**, and every standing
finding is disposed as `fixed` by a commit between the parked candidate and the
item head. The Correction Review then reviews the whole delta from the round-2
head. A round-1 park, a Correction Review park, or any dismissed round-2
finding keeps the path bound and its refusal text. Review and repository
checks still follow. Outside that exception, any other path still requires a
corrective Spec, and the refusal names that path.
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
interval up to `delivery.check_timeout`.

Line-scoped derived paths may resolve a conflicted path when every conflict
hunk contains only matching lines, and regeneration may change that path only
on matching lines. The two version fields in each owned `SKILL.md` are
line-scoped for the record command, so a derived merge can take the default
branch's version lines before the command raises and records the next version.

During `checking`, Roundfix waits while GitHub reports the Pull Request merge
state as `BLOCKED` or `UNKNOWN`, even when the listed checks pass. At the end of
the Check Wait the item parks as `checks-timeout: current-head checks not
settled within delivery.check_timeout <duration>`; only a mergeable state with
passing checks proceeds to merge. If GitHub refuses that merge because
`base branch policy prohibits the merge`, Roundfix returns to `checking` once
for that head. A second refusal for the same head parks `delivery-error`.

Each Check Wait lasts User Config `delivery.check_timeout`, a Go duration
greater than `0` and at most `24h`, default `60m`, from its own start; a re-run
of a failed check restarts it once. Project Config cannot set it. A Queue Owner
or `roundfix deliver retry` reads it when it starts, so raise it in User
Config, then run `roundfix deliver retry <slug>` for an item parked
`checks-timeout`.

A conflict park has class `conflict`. Its next action is to merge the default
branch into the item branch in the printed worktree, resolve the named paths,
commit, then run `roundfix deliver retry <slug>`. Retry accepts a head descended
from the candidate, records it, and resumes at `reviewing`.

`deliver status` prints each item's Spec slug, stage, blocker and worktree.
When the queue records an owner PID, at least one item is neither `merged` nor
`parked`, and that PID is proven absent, it adds this line after any `Warning:`
lines and before any `Park:` lines:

```text
Owner: Delivery Queue owner PID <pid> is not running; run roundfix deliver resume
```

If the PID is alive but its start-time identity differs from the recorded
owner identity, the state is `has a different process identity` instead of
`is not running`. A proven live owner, no recorded owner, or a queue whose
items are all merged or parked adds no `Owner:` line. If the identity cannot
be checked, status adds no `Owner:` line and writes
`roundfix: Delivery Queue owner PID <pid> could not be checked: <error>` to
stderr; it still exits `0`. Status is read-only and reclaims nothing, leaving
the recorded owner PID and identity for `deliver resume` to recover.

After any `Owner:` line and before `Limits:`, it prints one line per parked
item in queue order, except `queue-token-ceiling` (see Token usage below):

```text
Park: <slug> <class>: <next command>
```

The Park Class identifies the reason and the next action:

| Class | Blockers |
| --- | --- |
| `dependency` | `prerequisite-unmerged` |
| `conflict` | `pull-request-conflict` |
| `environment` | `linux-lane-container`, `qa-environment-partial`, `checks-timeout`, `item-worktree-missing`, `delivery-error`, `runtime-infrastructure` |
| `flaky-check` | a check that failed again outside the item's changed packages |
| `finding` | `run-unresolved`, `review-findings`, `corrective-spec-required`, `gate-failed`, `checks-failed`, `revalidation-failed`, `secret-found`: in `<worktree>`, rewrite the item branch from the earliest commit the blocker names so that each named line holds a value no Secret Pattern matches and nothing else in the final tree changes, then run `roundfix deliver retry <slug>` |
| `budget` | `run-budget-exceeded`, `queue-deadline` |
| `review` | `review-blocked`, `review-stale` |
| `authorization` | `unauthorized` |
| `unclassified` | an unknown blocker |

An unrecovered Provider System Error parks the item as
`runtime-infrastructure: <scope_id> hit a provider system error: <detail>`.
Each existing blocker keeps its previous next action. A queue without parked
items adds no `Park:` line.

A Delivery Retry from `runtime-infrastructure` is not counted against the retry
limit and re-enters as a `run-unresolved` retry does.

Before classifying a failed GitHub Actions check, the owner reads the failing
job's `tested-base` notice annotation and fetches the delivery remote's default
branch. The annotation records the full commit id of the default branch tip
merged by that attempt. If that tested base does not contain the current tip,
or is not a known commit after the fetch, the failure is stale. The owner
re-runs the failed jobs once per run and default tip and keeps polling, logging:

```text
roundfix: check stale: Delivery Queue item <slug>: <check> tested <sha>, default branch is at <sha>; re-run (run <id>)
```

A stale failure is not classified as `flaky-check` or `checks-failed`. A pass
after a stale re-run adds no Warning. The check timeout restarts at most once
across stale and outside-change re-runs; further default branch movement does
not extend it again. Stale re-runs are remembered only during this check wait.
An absent, malformed or ambiguous annotation keeps the previous rules.

When a required GitHub Actions check fails on its first attempt, or on the
first newer attempt after a stale re-run against the current tip, the owner
reads its failed log and compares the failing Go package directories with the
item's changed paths against the refreshed remote default branch. If every
failing package lies outside that change, it re-runs the failed jobs once,
logs the action, restarts the check timeout once and keeps polling. A pass adds
`flaky-check: <check> passed on re-run` to the item's Warning. A second failure
outside the change parks as `flaky-check: <package>, …`, with an action to fix
or re-run the check and then run `roundfix deliver retry <slug>`.

A failure in a changed package, a build or setup failure, a log without Go
packages, a non-Actions check, or a run already past attempt one without
eligibility from a stale re-run parks as
`checks-failed` without a re-run. An inspection or re-run error also parks as
`checks-failed` and is recorded in the owner log.

When one or more items are parked, status prints exactly one
`Pending question:` for the lowest-position parked item, the action that
answers it from the same Park Class table, and the count waiting behind it.
A dependency park can clear after its prerequisites merge. Other parks need
`deliver retry` or a new queue; elapsed time alone does not resolve them.
`deliver stop` ends the detached owner; `deliver resume` restarts it from the
persisted queue.

`roundfix deliver retry <slug>` first records a parked item `merged` when its
recorded Pull Request was merged from the item branch, or when the local
default branch already archives its Spec through a delivery commit outside
the item branch. Recording the merge does not count a retry and applies
regardless of the retry limit, queue deadline or token ceiling. It touches no
item workspace; the owner then runs the normal post-merge cleanup.

The retry first reads the recorded Pull Request through `gh pr view` in the
repository checkout. A merged Pull Request must name the recorded item branch
and have both a head and a merge commit. Otherwise it reads the Spec's Archive Record, or a legacy archive's `_prd.md` on the local default branch and proves that the delivery commit that added that file is absent from the item head: the local item branch's tip, or the
newest candidate when the branch is gone. This fallback does not fetch; the
local default branch must already hold the delivery.

A recorded merge prints this line on stdout before `Retried <slug>: <blocker>
-> merged`:

```text
Merged outside the queue: <evidence>; merge commit <sha>
```

`<evidence>` is `pull request #<n>` or `Spec archived on default branch
"<branch>"`; `<sha>` is the full merge commit from the Pull Request or the
delivery commit from the archive proof. A failed Pull Request read refuses
the retry and leaves the item unchanged.

Otherwise the retry returns the parked item to the queue. For an
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

A Secret Rewrite is a head that does not descend from the parked candidate
but does descend from the item's Run start. Only paths named by the
`secret-found` blocker may differ from the candidate, with the same modes
and line counts; every changed line must hold a Secret Match at the candidate,
and no named file may hold a Secret Match at the rewritten head. The retry
records the head as the newest candidate and returns the item to `reviewing`.
A refusal keeps the archived-head reason and appends the first failing rule
in parentheses, naming paths and line numbers without values.

| Recorded evidence | Re-entry stage |
| --- | --- |
| Recorded Pull Request merged from the item branch, or Spec archived on the local default branch through a delivery commit outside the item branch | `merged`, without counting a retry; owner runs post-merge cleanup |
| Active Spec with any unfinished Task | `running` |
| Active Spec with every Task completed | `reviewing` |
| Archived Spec with a post-archive correction descended from its newest candidate | `reviewing` |
| Archived Spec parked `secret-found` with a Secret Rewrite of its parked candidate | `reviewing` |
| Operator-archived item with a QA override and head descended from its candidate or Run start | `reviewing` |
| Resolved `pull-request-conflict` with a head descended from the candidate | `reviewing` |
| Archived Spec with unchanged candidate and no recorded pull request | `gating` |
| Archived Spec with unchanged candidate and a recorded pull request | `checking` |
| `corrective-spec-required` with the parked candidate head unchanged | `reviewing`, without Task Carry-Forward |
| `corrective-spec-required` with a moved head | `reviewing` when API Contract 5 and API Contract 6 prove the correction answers only the review recorded in the Reviewed Head Record of the parked candidate, descends from the parked candidate, every standing finding has exactly one valid disposition, and every changed path lies under an archived Spec the blocker names or among the paths its Tasks declared (`creates:`, `interface:` and recorded paths read at the recorded revision), except that a round-2 park that is not a Correction Review and has only `fixed` dispositions accepts any changed path for its one Correction Review; review and repository checks still follow; otherwise refused with exit `2`, the item stays unchanged, and the refusal names the path and requires a corrective Spec with its own authorization and QA gate |

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
a refused carry-forward, a recorded Pull Request closed without merging, a failed
merge observation, or an owner hand-off failure exits `2` and starts no owner.
An item-level refusal leaves the stored item unchanged.

A recorded Pull Request closed without merging, with no merge evidence on the
local default branch, refuses the retry with this reason:

```text
retry Delivery Queue item "<slug>": pull request #<n> was closed without merging; reopen it or merge the Spec into the default branch, then run roundfix deliver retry <slug>
```

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

When `deliver retry`, `status`, `resume` or `stop` runs in a recorded item
worktree that has no queue of its own, it exits `2` and names the item and the
checkout that owns its queue. The reason adds `this checkout is the item
worktree of <slug> in the Delivery Queue of "<checkout>"; run the command in
"<checkout>"`. Run the command from that owning checkout. The refusal does not
act on that queue: it changes no item or owner record and starts no owner.
`retry` keeps its `Retry refused` layout; the other three commands keep their
`Preflight failed` layout. Outside a recorded item worktree, the missing-queue
reason stays unchanged.

#### Linux Lane

With `delivery.linux_lane: true` (default `false`), the Delivery Queue runs the
Linux Lane once per candidate, after repository Verification passes and before
publication. The setting comes from the fetched default-branch Project Config,
with User Config beneath it. An item branch's opt-in takes effect after merge.

The lane tests the packages changed since the candidate's merge base with the
delivery remote's default branch. Each changed path selects its nearest parent
directory containing a Go file at the candidate head, including changes under
`testdata`; the sorted package list has no duplicates. A changed `go.mod` or
`go.sum` selects `./...`. Uncommitted files and later checkout changes are
excluded: `git archive` of the candidate supplies standard input, and the lane
writes nothing into the item worktree.

The container uses `delivery.linux_lane_image` from the same fetched default
configuration as the opt-in. An item branch or the owner's checkout cannot
choose its image. Empty by default, the key selects `golang:<version>` using
the root `go.mod`'s `go` version verbatim. The container runs under an init
with `--init`, `--pull never`, `--network none`, `--rm`, and user
`65534:65534`. Beside the host module cache, it mounts the repository object
store read-only as `-v <objects>:/repo-objects:ro`. The store is resolved with
`git rev-parse --path-format=absolute --git-path objects`, so a linked
worktree uses its common directory's store. The working tree and refs are
never mounted. After `git init`, the container repository registers
`/repo-objects` in its own `.git/objects/info/alternates` before recording the
candidate; the host store is never written. Its only explicitly
supplied environment variables are `HOME`, `GOCACHE`, `GOMODCACHE`, `GOFLAGS`,
`GOPROXY`, `GOTOOLCHAIN`, and `GOSUMDB`; no Queue Owner key or token is forwarded.
The Docker CLI keeps the operator's own environment and Docker context. Tests
use a temporary Git repository inside the container and cannot download modules
or toolchains; missing cached modules fail the lane. The operator must make the
image available locally before the lane can run; the lane never pulls it.

The following skip reasons are checked in order. A skip passes the gate and
writes this exact notice to the Queue Owner's console log and
`<artifact dir>/delivery/<slug>/linux-lane.log`:

```text
roundfix: notice: Delivery Queue item <slug>: Linux Lane skipped: <reason>
```

- `no go.mod at the repository root`
- `go.mod declares no go version` (only when no image is configured)
- `no changed Go package`
- `the object store <objects> cannot be mounted: its path holds a colon`
  (outside a Windows volume name)
- `the object store <objects> cannot be mounted: it is not a directory`
- `the object store <objects> borrows objects through info/alternates, which the Linux Lane does not mount`
- `go is not on PATH`
- `docker is not on PATH`
- `docker info failed: <first non-empty output line>`
- `image <image> is not present locally`
- `docker could not start the container: <first non-empty output line>` (exit `125`)

`no go.mod at the repository root` and `no changed Go package` are expected
skips that leave only the notice. Every other skip of an enabled lane, and a
fetched default whose configuration cannot be read, is a Lane Environment Skip.
The item records the warning `linux-lane-skipped: <reason>`, which
`roundfix deliver status` prints on its `Warning:` line. The Pull Request body
adds this paragraph:

```text
Linux Lane: not run. The Delivery Queue gate skipped it: <reason>. This candidate's Go packages were not tested on Linux before this pull request.
```

A later gate replaces the earlier lane warning with its own outcome. A skip
still passes the gate and never parks the item.

The lane log contains the command line and combined output. Exit `0` passes
with this console line. `<n>` is the number of packages `go list` selected in
the container, which the lane script reports with
`roundfix-linux-lane: packages: <n>` before `go test`:

```text
roundfix: Delivery Queue item <slug>: Linux Lane passed <n> package(s) on <image>
roundfix: Delivery Queue item <slug>: Linux Lane passed on <image>; the lane script reported no package count
```

Without a count line, exit `0` stays a pass and uses the no-count line above.

A **Lane Container Failure** fails the gate and parks the item in the
`environment` Park Class. The lane script checks for `tar`, `git` and `go`,
prepares `/tmp/work`, unpacks the candidate and records it before testing.
Setup failures print `roundfix-linux-lane: container: <cause>` and exit `3`.
The cause is selected in this order:

- The 20-minute cap: `the container exceeded 20m0s`. The cap cancels the run
  and removes its named container with `docker rm -f <name>`.
- Exit `3`: the text after the last `roundfix-linux-lane: container: ` output
  line, or `the lane script exited 3 before go test` when none exists. Setup
  causes name `missing tool <tool>`,
  `cannot read the repository history at /repo-objects`, `cannot prepare /tmp/work`,
  `cannot unpack the candidate` or `cannot record the candidate`.
- Exit `126`: `docker run exited 126: the image could not invoke the lane command`.
- Exit `127`: `docker run exited 127: the image lacks the lane command`.
- Exit `128` or more: `docker run exited <n>: the container process ended by signal <n-128>`.

The gate reason and parked blocker are respectively:

```text
Linux Lane container failure: <cause>; see <lane log>
linux-lane-container: <cause>; see <lane log>
```

A **Lane Change Failure** is any other non-zero exit of a run that started;
exit `125` stays a skip. It parks in `finding`, with cause `go test exited <n>`.
The gate reason and parked blocker are respectively:

```text
Linux Lane change failure: go test exited <n>; see <lane log>
gate-failed: Linux Lane change failure: go test exited <n>; see <lane log>
```

The lane log ends with exactly one outcome line:

```text
Linux Lane container failure: <cause>
Linux Lane change failure: go test exited <n>
Linux Lane passed <n> package(s) on <image>
Linux Lane passed on <image>; the lane script reported no package count
```

For `linux-lane-container`, the next action is `repair the Linux Lane container
cause in the blocker (Docker, the image, or delivery.linux_lane_image on the
default branch), then run roundfix deliver retry <slug>`. A repository
Verification failure still parks the bare `gate-failed` with its generic next
action.

Cancelling the gate also removes the container and returns the context's error.
With the setting off, the gate starts no Docker or Go process and emits no lane
notice.


### Run the Linux Lane on demand

```bash
roundfix deliver lane [--all] [<revision>]
```

Run `roundfix deliver lane --all` and obtain a pass before turning
`delivery.linux_lane` on. The command runs independently of that setting,
the Delivery Queue, the Run Database and any Spec. It reads the checkout's
Project Config over User Config, including `delivery.linux_lane_image` and
`defaults.artifact_dir`; it does not read configuration from the fetched
default branch. An empty image uses `golang:<go.mod go version>` from the
selected commit.

The revision defaults to `HEAD` and must resolve to a commit. The command
archives that commit, without including uncommitted changes. By default it
selects the changed Go packages as the queue gate does. `--all` selects
`./...` without reading a default branch or merge base. It uses the same
network-free, never-pull container with `--init` and the same read-only history
mount `-v <objects>:/repo-objects:ro`, named `roundfix-linux-lane-<pid>`, and
writes only `<artifact dir>/delivery/linux-lane.log`.

A checkout that borrows objects through `info/alternates` skips until
`git repack -a` and removal of `objects/info/alternates` make its store
self-contained. Resolve the store with the Git command above to locate that
file; the lane never copies borrowed objects or changes the host store.

Exactly one outcome line goes to standard output, with empty standard error:

```text
roundfix: notice: Linux Lane skipped: <reason>
roundfix: Linux Lane passed <n> package(s) on <image>
roundfix: Linux Lane passed on <image>; the lane script reported no package count
roundfix: Linux Lane change failure: go test exited <n>; see <lane log>
roundfix: Linux Lane container failure: <cause>; see <lane log>
```

Exit codes are `0` for a pass, `1` for a Lane Change Failure, and `3` for a
skip or Lane Container Failure. An interrupt removes the container before
exiting `130`. An extra argument, unknown flag, revision that is not a
commit, invalid configuration, or unknown default branch without `--all`
exits `2` with the `deliver lane` **Preflight failed** layout on standard
error and no outcome on standard output.


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

## Run slot wait

When an item's `roundfix implement` is refused only because Active Implement
Runs hold every `runs.max_active` slot, the item waits in stage `running`.
`roundfix deliver status` shows its blocker text on the item's status line:

```text
run-slot-wait: until <deadline>; <n> Run(s) hold <max> slot(s): <run id> in <repository> for <spec>
```

Multiple holders are separated by `, `. The deadline is UTC RFC3339. The Queue
Owner prints these console lines when it begins waiting and when a slot frees:

```text
roundfix: Delivery Queue item <spec>: waiting for a Run slot until <deadline>; <n> Run(s) hold <max> slot(s)
roundfix: Delivery Queue item <spec>: a Run slot is free; starting the Implement Run
```

Every 30 seconds it polls the Run Database without starting a child, reclaiming
a Run or signalling a process. Only a holder whose recorded owner PID is no
longer alive is excluded; a holder without a recorded PID still counts. When a
slot looks free, it clears the blocker and retries Implement. The child's
preflight and the store's CreateRun transaction still decide admission. Losing
the race keeps the original wait deadline.

User Config `runs.slot_wait` bounds the wait from the first refusal; its default
is `2h`, and `0` parks immediately. Project Config cannot set it. The earlier
queue deadline wins and parks the item as `queue-deadline`. At the wait's own
deadline, the item parks as `delivery-error` with the same Implement failure
reason it would have reported without waiting. Other refusals park as before.
The blocker is cleared when the wait ends; stopping the owner clears it, and
resume or retry starts a new wait.
