# Command reference

Every Roundfix command: what it does, what it prints, and the boundary it never
crosses. This is the contract page — for the task-oriented walkthrough, read the
[operational guide](usage.md); for config keys, read
[configuration](configuration.md).

Examples call the installed `roundfix` binary. From a source checkout without
installing, substitute `go run ./cmd/roundfix`.

## Global contract

- **stdout** carries only the deterministic report of the requested command.
  Diagnostics, progress, the Run ID, and Agent output go to **stderr**.
- **Exit codes**: `0` Clean, Stopped, Fetched, or an already-complete no-op;
  `1` Unresolved, Failed, or Integration Pending; `2` Preflight Validation
  failure; `3` Clean Unverified (watch only); `130` in-terminal Ctrl-C.
- Color is automatic in interactive terminals. `ROUNDFIX_COLOR=always` forces
  it, `ROUNDFIX_COLOR=never` or `NO_COLOR` disables it.
- Supported Agent names are `codex`, `claude`, and `opencode`. The supported
  Review Source is `coderabbit`.

## Setup and maintenance

### setup

```bash
roundfix setup [--yes] [--no-input]
```

Verifies Node.js, the minimum supported acpx version, the effective adapter, generated
Agent Selection Profiles, acpx local adapter overrides, User Config, and
Project Config. Codex Adapter Readiness requires official
`@agentclientprotocol/codex-acp` lineage at version `1.1.4` or newer. A stale
bare override that resolves to legacy `@zed-industries/codex-acp` produces one
migration offer to `npx -y @agentclientprotocol/codex-acp@1.1.4`.

Setup builds every proposed file in memory and runs exact Agent Selection proof
before writing. It never changes explicit Sol/high to model-managed reasoning
when proof fails. Each check prints one deterministic report line such as
`node: ok`, `adapter: migration proposed`, `profile readiness: passed`, or
`User Config: skipped`. `--yes` accepts every offered install or file change;
`--no-input` performs diagnosis and skips offers without writing.
When acpx is missing or older than `0.12.0`, Setup offers
`npm install -g acpx@0.12.0`. It accepts `0.12.0` and newer versions without
offering a downgrade.

### doctor

```bash
roundfix doctor
```

Read-only readiness report; mutates nothing and exits nonzero when any check
fails. One stdout line per check with `ok`, `failed`, or `skipped`; failure
lines include `next: <action>` when a remediation is known. The checks:

- `node:` — Node.js meets the minimum version.
- `acpx:` — the installed acpx version is at least the minimum supported
  version. Newer versions are accepted and are not downgraded.
- `adapter:` — the effective adapter command proves the required package
  lineage and supported version; legacy, unknown, old, and missing adapters
  fail with the official install action.
- `profiles:` — the required Agent Selection Profiles pass exact proof through
  disposable ACP Sessions. Success names distinct tuples and category
  references. Failure names the exact tuple, affected categories,
  classification, bounded adapter evidence, and the next
  `roundfix profiles configure` or `roundfix profiles validate` action. A
  rejected explicit `high` does not recommend model-managed reasoning.
- `skills:` — the required Repository Skill Set matches its local
  authorities. The running binary's embedded artifacts are authoritative for
  the 14 Roundfix-owned skills, including the Roundfix Skill. Each of the 25
  required external skills must hash to its `computedHash` in
  `skills-lock.json`.
- `codex:` — macOS-only runtime hygiene: inspects `com.apple.quarantine` (the
  real XProtect trigger) and code-signature validity, resolving `CODEX_PATH`
  first and then `codex` on `PATH`. It does not use `spctl --assess`, which
  rejects any signed CLI that is not a notarized app. A quarantined or
  improperly-signed codex fails with the next action to reinstall codex with
  the official curl installer into `~/.local/bin` and set `CODEX_PATH`.
  Skipped on non-Darwin platforms.

Doctor has no separate `agent:` or `model:` authority. The aggregate
`profiles:` result is the Agent Selection Profile Readiness contract.
Repository Skill Set readiness runs after it as an independent check, even
when profile proof fails, and appears before `codex:`.
Outside Git, profile proof uses the process working directory while Repository
Skill Set inspection does not run. The `skills:` line reports
`Repository Skill Set readiness requires a Git repository` and keeps the
run-from-Git next action.

```text
node: ok
acpx: ok
adapter: ok (...)
profiles: ok (3 distinct tuples; 10 category references)
skills: ok (39 required: 14 Roundfix-owned, 25 external)
codex: ok
```

A missing or outdated required skill, or an invalid required lock declaration,
prints one sorted blocking line and makes Doctor exit `1`. Doctor still prints
every other readiness result:

```text
skills: failed (missing: handoff; outdated: roundfix; next: roundfix skills install --target project && bunx skills experimental_install && bunx skills update -p -y)
```

For mixed ownership, Doctor joins the Roundfix-owned restore and external
update actions with `&&`, so external remediation runs only after the owned
restore succeeds. A failure owned only by the external lock or skill set prints
only the external action.
Doctor never runs either command, never deletes skills, and never updates
`skills-lock.json`. The check is offline and read-only: it reads only local
embedded artifacts, `.agents/skills`, and `skills-lock.json`. Unrelated extra
installed skills and lock entries are ignored and are not removed or flagged.

### upgrade

```bash
roundfix upgrade [--check]
```

Resolves the latest release through the GitHub CLI. Successful stdout outcomes
are `upgraded 1.0.0 → 1.1.0`, `already current 1.0.0`, `no releases published`,
and, with `--check`, `upgrade available 1.0.0 → 1.1.0`. Failures leave the
current binary untouched and print a manual fallback on stderr. Operational
commands run a best-effort daily freshness check that prints one stderr line
when the binary is behind.

### init

```bash
roundfix init [--scope user] [--force]
```

Creates Project Config (`<repo>/.roundfixrc.yml`) or, with `--scope user`,
User Config (`~/.roundfix/config.yml`). When `--scope` is omitted, Roundfix
asks and defaults to Project Config. `--force` overwrites an existing file.

### gc

```bash
roundfix gc --dry-run
roundfix gc
```

Non-interactive storage reclamation. Resolves `store.journal_retention`,
prunes eligible terminal Runs' Run Event Journal rows and
`<artifact-dir>/runs/<run-id>` directories, removes orphaned `runs/<id>`
directories, and reports Runs, journal rows, and artifact bytes reclaimed.
`--dry-run` lists the same eligible set without deleting. Retention never
deletes Active Runs, `runs` rows, active-run locks, or Review artifacts under
the Spec Root.

### skills

```bash
roundfix skills list
roundfix skills check
roundfix skills install [--target codex|claude|opencode|all]
```

The binary ships 14 Roundfix-owned skills: the operational `roundfix` skill
plus the authorial workflow skills (`write-idea`, `write-prd`,
`write-techspec`, `write-tasks`, `setup-context-driven`, `implement-task`,
`implement-spec`, `brainstorming`, `council`, `business-analyst`,
`archive-spec`, `qa-gate`, `evidence-gate`). `skills list` also prints
recommended external skills, which install through your own tooling and are
never shipped. By default `skills install` writes to `<repo>/.agents/skills`;
`--target` selects user-scoped Agent skill directories instead.

### profiles

```bash
roundfix profiles show [--category <category>] [--json]
roundfix profiles configure --scope user|project [--file <path>] [--dry-run] [--yes] [--json]
roundfix profiles validate [--category <category>] [--json]
```

`profiles show` renders the effective Preferred Selection, Fallback Chain, and
dated advisory recommendations. Official model identifiers and advisory rank
do not prove that a tuple works in the current environment.

`profiles configure` validates a complete candidate and runs exact Agent
Selection proof for every distinct preferred and fallback tuple before
confirmation. `--dry-run` performs the proof without writing; proof failure,
decline, or output failure preserves the target bytes. `profiles validate` is
read-only, deduplicates exact tuples across category references, proves them
through disposable ACP Runtime Sessions, sends no Agent prompt, and closes
every Session on success or error. JSON schemas are
`roundfix/profiles/v1`, `roundfix/profiles-configure/v1`, and
`roundfix/profiles-validate/v1`.

### baseline

```text
roundfix baseline [--repo <path>] [--format <text|json>]
roundfix baseline plan --profile <id> [--decision <id=value> ...] [--decision-file <path> ...] [--repo <path>] [--format <text|json>]
roundfix baseline apply --plan <file> --confirm-plan <digest> [--repo <path>] [--format <text|json>]
roundfix baseline profile init --id <id> [--from <built-in-id>]
roundfix baseline profile show <id> [--format <text|json>]
roundfix baseline profile validate [<id>|<path>] [--format <text|json>]
roundfix baseline skills restore --profile <id> [--skill <name> ...] [--source-dir <path>] [--confirm-plan <digest>] [--repo <path>] [--format <text|json>]
roundfix baseline assets sync --source-dir <path> [--check] [--format <text|json>]
```

The root command is the terminal-only human workflow for first adoption and
later updates. It detects existing state, collects instruction preservation,
one Baseline Profile, and repository decisions, then presents one consolidated
Change Plan. Rejecting the plan returns to a selected decision area and
recalculates the complete plan. Mutation requires one explicit confirmation of
the current displayed Plan Digest.

`baseline plan` is read-only, non-interactive, local, and network-free. JSON
exit `0` emits one complete `roundfix/baseline-plan/v1` document. Missing
decisions, manual classification, or unresolved alignment exit `3` with a
`roundfix/baseline-result/v1` next action and no partial plan.

`baseline apply` accepts only a strict portable JSON plan and its exact
`planDigest`. It validates clone-stable Git lineage, the embedded catalog,
profile identity, every bounded preimage, the complete managed-entry ledger,
and its derived file projection. It never substitutes a newer plan. Matching
postimages make an exact reapply an idempotent success.

`baseline profile init` creates
`.roundfix/baseline/profiles/<id>.json` from one embedded built-in profile.
`show` resolves one built-in or repository-owned profile; `validate` checks one
ID, one direct profile path, or every repository-owned profile. Repository
profiles can use embedded entry IDs only and cannot compose profiles or load
remote executable content.

`baseline skills restore` is a separate confirmation-gated Repository Skill
Set operation. Its non-empty preview exits `3` with a current Plan Digest;
`--confirm-plan` applies only that exact preview. `--source-dir` selects an
offline Git checkout or bare object store containing the declared immutable
source commit.

`baseline assets sync` is a maintainer operation over an explicit canonical
setups directory. `--check` is read-only. Refresh validates the generated
catalog before updating only Go-owned canonical setup snapshots.

Baseline requested output goes to stdout; diagnostics and progress go to
stderr. Exit categories are `0` success or current no-op, `1` execution,
verification, output, recovery, or incomplete-rollback failure, `2` invalid
input/schema or unsafe repository, `3` another owner action or renewed
approval is required, and `130` cancellation. Baseline reports repository
formatter and Verification commands as recommendations but never runs them.

For the adoption, automation, Decision Document, cross-clone, migration,
recovery, and security procedures, read
[CONTEXT-driven development](context-driven-development.md#adopt-or-update-the-context-driven-baseline).

## Review loop: fetch, resolve, watch

Review Runs execute in the user's checkout on the PR Head Branch — they create
no Run Worktree and no Run Branch, so review fixes are always a delta over the
published HEAD and Integration Pending does not exist as a review outcome.
Preflight Validation requires a clean tracked working tree (untracked files are
allowed); after a failed batch, every dirty path in the checkout is Agent work
from that Run, because the tree was clean at start.

**Branch Integrity Preflight** runs on all three commands before any Review
Source fetch, Agent Session, comment, or code change. It enumerates
`roundfix/run-*` branches with commits based on the PR Head Branch:
fast-forwardable pending work is integrated automatically and reported;
anything else refuses with exit `2`, naming each branch, its ahead count, and
the exact integration command. It also refuses while another Active Run is
bound to the target, naming the run id and the stop commands.
`--skip-branch-integrity` bypasses both guardrails only after publishing an
audit comment on the pull request recording the run id, the skipped
guardrails, and the ignored state; a failed publish fails the command.

`resolve`, `watch`, and `implement` accept exactly two Agent Selection forms:
omit `--agent`, `--model`, and `--reasoning-effort` to use profiles, or provide
all three together for a complete one-Run override. Any partial subset exits
`2` before config load, adapter or profile proof, Session creation, or Run
mutation. For example:

```bash
roundfix watch --source coderabbit --pr <number> --until-clean
roundfix watch --source coderabbit --pr <number> --agent codex --model gpt-5.6-sol --reasoning-effort high --until-clean
```

### fetch

```bash
roundfix fetch --source coderabbit --pr <number> [--spec <slug>]
```

Validates local state, creates a Fetch Run, fetches unresolved CodeRabbit
review threads, writes markdown Round artifacts, and stops at the `Fetched`
outcome. It never starts an Agent, commits, pushes, or resolves Review Source
threads. With automatic Round selection it reuses an existing matching Round
when the same HEAD already has the same Review Issue fingerprints and never
overwrites existing Round artifacts.

### resolve

```bash
roundfix resolve --pr <number> [--spec <slug>]
```

Works only over downloaded Compatible Artifacts — it does not fetch. It
assigns bounded Batches to the Agent, verifies each Batch, and commits
successful Batches directly on the PR Head Branch when auto-commit is enabled.
At Batch settlement each Review Issue propagates to the Review Source
individually: `resolved` threads resolve; `invalid` and `duplicated` threads
get an explanatory Outcome Comment and then resolve; `failed` threads get the
failure reason and stay open; issues still unresolved at Run end receive a
closing comment. Comments carry an idempotency marker, so retries never
duplicate. The Run's review artifacts are committed in one separate docs
commit (`docs: review round NNN for pr <n>`, ADR-0036) and Final Push runs
only when no Unresolved Review Issues remain. Artifact roots outside the
repository, or reached through a symbolic link, are reported and never staged.

### watch

```bash
roundfix watch --source coderabbit --pr <number> --until-clean [--max-rounds N]
```

Waits for Review Source Evidence on the current PR HEAD, observes the
configured quiet period, fetches, resolves Batches, and repeats. Evidence is
always bound to the expected head:

- `pending` means no usable expected-head signal exists. A stale check or
  review remains visible as detail but cannot verify the expected head.
- `reviewing` means a current-head CodeRabbit check or status is pending or in
  progress.
- `reviewed` means CodeRabbit produced a current-head result that does not
  prove Merge-Ready. A successful check, status, or approval stays `reviewed`
  while an unresolved CodeRabbit thread exists; a non-approved review is also
  only `reviewed`.
- `verified` accepts a successful current-head CodeRabbit check or commit
  status, or a current-head CodeRabbit `APPROVED` review, only with zero
  unresolved CodeRabbit threads.
- `skipped` requires an explicit structured CodeRabbit skip for the expected
  head. It ends Review Skipped with exit `3`, prints the Review Source reason
  and next action, fetches no Review Issues, and cannot mean Clean, Clean
  Unverified, or a zero-issue Round.
- `failed` records an explicit current-head Review Source failure.

`WaitingForReview` is the pre-fetch phase.
`WaitingForReviewCheck` is the Merge-Ready phase after Final Push. Each wait
records its expected head, start time, deadline, Evidence state and kind, and
retry status. Non-TTY progress prints on phase entry and when Evidence or retry
state changes; the Live Run View derives remaining time from the deadline.

Roundfix retries only a typed transient Review Source failure: a context
deadline not caused by Run cancellation, a temporary DNS failure, a connection
reset, HTTP `429`, or a GitHub `5xx` response. One episode records `started`,
then `recovered` or `exhausted`. Retry sleeps use the existing poll interval
and remain bounded by the existing Review Source timeout and Run Budget; there
is no new retry setting. The watch loop does not infer retryability from its
Console Log, progress lines, or Run Event summaries.

After Final Push, a proven Daemon-created review-artifact commit can inherit
its parent's verified Evidence. Roundfix requires the recorded commit to be
the current head, to have exactly the recorded verified parent, and to retain
the exact Daemon-generated artifact-commit subject. Its stageable review root
must be inside the repository without a symbolic-link crossing, its diff must
be non-empty, and every changed path must be below that root. Roundfix then
refreshes the parent's Evidence, which must still be verified with no
unresolved CodeRabbit threads. Missing or changed identity, a wrong or multiple
parent, a non-current head, a changed subject, an external or symbolic-link
root, an empty diff, any mixed or out-of-root path, stale parent Evidence, or
an unresolved thread refuses inheritance and returns to ordinary current-head
Evidence polling. User-authored documentation commits and all other
non-Daemon descendants never inherit.

If accepted Evidence never appears within `watch.check_grace_period` (default
`5m`), watch ends Clean Unverified with exit `3` and names the next action.
Other terminal outcomes are `MaxRoundsReached`, `BudgetExceeded`, `TimedOut`,
`Failed`, and `Stopped`.

### Review report shape

One line per Review Issue in Round order — with a ` — reason:
<terminal_reason>` suffix on failed, unresolved, and invalid lines when the
artifact carries one — then two labeled summary lines separating this Run's
counts from the pull request's cumulative counts:

```text
issue 001 resolved — major: handle test issue
This Run (Clean after 1 Round(s)): 1 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved.
Pull Request cumulative: 1 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved.
```

Before a Review Source fetch completes, counts are not known. The report omits
all zero-valued status summaries and prints only:

```text
Review Issues: unknown — fetch did not complete.
```

Review Skipped uses its own two-line source reason and next-action report; it
does not use either count shape.

## Spec loop: implement, settle, archive

### implement

```bash
roundfix implement --spec <slug> [--qa] [--detach]
```

Executes a Spec's Task Graph in a Run Worktree as one Run — spec Runs keep
worktree isolation. The scheduler executes the current Wave up to
`worktree.concurrency` at a time (default `2`; `1` keeps sequential behavior),
with concurrently running Tasks in Task Worktrees and one commit per completed
Task on the Run Branch. It resolves `specs.root` once from the user's checkout.
`--qa` ends the Run with the qa-gate step; only a `pass` verdict lets the Run
end Clean. `implement.auto_push: true` makes a Clean Run push its branch
upstream. Integration Pending, Unresolved, Failed, Stopped, and failing-QA
Runs never push.

stdout is one line per Task in Task Graph order — failed and skipped Tasks are
followed by one indented `  reason: <one line>` naming the failed step (for
Verification failures: the command, exit status, and diagnostics path) — then
the QA line when requested and one outcome line:

```text
task_01 completed — first task
task_02 failed — second task
  reason: verification failed: make verify (exit status 2); see <path>
Unresolved: 1 completed, 1 failed, 0 skipped, 0 pending.
```

Task status vocabulary is normalized on reload: `done` and hyphen/space
variants of the canonical statuses map to canonical form and the task file is
rewritten; anything else still fails validation. A Task whose commit contains
no change outside the Spec Root still settles `completed`, with one stderr
warning and one Run Event marking the no-op.

Daemon Task and QA commits stage only repository paths that do not cross a
symbolic link; dropped paths are journaled and warned
(`roundfix: task file <path> kept outside the repository; committed without it`).

### settle

```bash
roundfix settle --spec <slug> --task <task_id>
```

Local recovery for one failed Task whose completed work already exists in a
kept Task Worktree, kept Run Worktree, or the current repository. Surface
resolution picks the first candidate, in that order, **where the task file is
actually `failed`**, and always reports the choice on stderr as
`Settle surface: <path>`; when no candidate qualifies, the refusal names each
candidate path and the status found there. It re-runs the Task's Verification
commands in the selected tree, changes nothing on failure, and on pass prints
one sorted `commit <path>` line per committed path before the settled line:

```text
verify go test ./... — ok
commit internal/example/example.go
commit docs/specs/<slug>/task_01.md
settled task_01 completed — <short sha>
```

When other Tasks of the Spec are `failed` at settle time, one stderr warning
names them: their work may be swept into this commit. Settle creates no Run,
writes no Run Event Journal entries, and never pushes; Task Worktree
settlements integrate onto the Run Branch before the Run-level integration.

### archive

```bash
roundfix archive <slug>
```

Non-interactive; creates no Run and never pushes. Verifies every Task is
`completed` and the newest QA Report has `verdict: pass`, then stamps
`_prd.md` and moves `<specs.root>/<slug>/` to `<specs.root>/_archived/<slug>/`.
Refusals exit `2` naming the first unmet condition.

### reconcile

```bash
roundfix reconcile [run-id] [--apply] [--format <text|json>]
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
`applyCommand`, `results`, and `summary`.

Run Worktree Reconciliation uses five states:

| State | Evidence and behavior |
| --- | --- |
| `safe` | The Run Branch and recorded target branch resolve, any present Run Worktree is registered and clean including untracked files, and the Run Branch tip is an ancestor of the current target tip. This is the only state eligible for cleanup. |
| `unintegrated` | The worktree cleanliness and ref evidence resolve, but the Run Branch tip is not an ancestor of the target tip. Roundfix preserves the worktree and branch. |
| `dirty` | A present registered Run Worktree has tracked or untracked changes. Dirty evidence takes precedence and Roundfix preserves the worktree and branch. |
| `unknown` | Invalid or missing metadata, an unsafe or unregistered worktree, an ambiguous or missing ref, or a Git inspection failure prevents proof. Roundfix preserves every surface it can identify. |
| `released` | Both the Run Worktree and Run Branch are absent. Repeated dry-run or apply is an idempotent no-op. |

A missing worktree alone is not `released` and never authorizes deletion. When
the Run Branch remains, Roundfix still requires an unambiguous Run Branch tip,
the recorded target tip, and ancestry proof. Age and terminal outcome are also
not cleanup evidence.

`--apply` is the only mutation switch:

```bash
roundfix reconcile <run-id> --apply
roundfix reconcile --apply
```

There is no force flag or user assertion that bypasses the proof. Apply acts
only on entries classified `safe` during that invocation, then rechecks the
metadata, worktree cleanliness, Run head, target head, and ancestry before
mutation. It removes the Run Worktree without force, deletes the Run Branch,
and reports failures while preserving any remaining path or ref. Dirty,
unintegrated, unknown, and released entries remain successful preserved
results unless an operational inspection fails.

Before safe cleanup, Roundfix records the reconciliation evidence. A safe
Integration Pending Run moves to Clean through the guarded terminal transition;
every other terminal outcome remains unchanged. The Reconcile Command never
integrates unique commits, repairs dirty work, chooses another target branch,
or treats a missing path as proof.

The contract uses the [Roundfix glossary](../../CONTEXT.md#language) and follows
[ADR-0053](../adr/0053-terminal-run-worktree-reconciliation-is-proof-based.md)
and
[Spec 0038](../specs/_archived/0038-terminal-run-worktree-reconciliation/_prd.md).
Adjacent terminal-cleanup diagnostics remain traced through the
[Stop Command](#stop) to the
[detached-watch finding](../findings/2026-07-16-vortex-pr87-detached-watch-notification.md#4-cleanup-noise-appeared-before-the-actionable-failure).

## Run discovery and monitoring

### runs and runs list

```bash
roundfix runs                     # Run Browser at an interactive terminal
roundfix runs list [--state active|terminal|all] [--limit N] [--all]
```

Bare `runs` opens the machine-wide read-only Run Browser: every repository's
Runs newest first, Active only by default, `a` toggles active/all, `Enter`
attaches read-only, `q`/`Esc`/`Ctrl-C` quits with exit `0` and no side
effects. In a non-interactive context it exits `2` naming `runs list`.

`runs list` prints one Run per line, newest first:
`<run-id>  <state>  <kind>  <target>  <agent>  <started-utc>  <duration>
<branch>`. Targets are `pr:<number>` or `spec:<slug>`. Default scope is the
current repository's 20 newest Active Runs; `--state` widens the state filter,
`--limit N` changes the bound (`0` unbounded), `--all` lists every repository
and adds a repository column. When Runs are hidden, exactly one trailing
stderr note names the hidden count and the widening flag. Empty results print
`No Runs found.` and exit `0`.

When the selected repository scope contains retained terminal Run Worktrees or
Run Branches, Runs List keeps the stdout row shape unchanged and writes one
diagnostic to stderr:

```text
(2 terminal Run Worktrees retained; run 'roundfix reconcile' to inspect)
```

The count can appear even when the default Active view has no rows. It is
discovery guidance, not a `safe` or unsafe classification; use the Reconcile
Command to classify each retained Run. When present, this retained-worktree
diagnostic is the one trailing stderr note.

### attach

```bash
roundfix attach [<run-id>]
```

Read-only. With a Run ID it replays the Run Event Journal and follows live
events — never creating Runs, fetching, starting Agents, committing, pushing,
stopping, or resolving threads. Without a Run ID at an interactive terminal it
opens the Run Browser; in non-interactive mode it exits `2` naming
`roundfix runs list`. The Live Run View shows a `WORK QUEUE` pane next to a
`SESSION.TIMELINE` pane grouping Run Events by Batch; raw payloads never
render inline and full content stays in the Detail Modal.

### events

```bash
roundfix events <run-id> [--follow] [--filter task-status,batch,verification,outcome,agent-selection]
```

Writes only `roundfix-events/v1` JSONL records to stdout; diagnostics go to
stderr. The public Supervisor categories, in journal cursor order, are
`task-status`, `batch`, `verification`, `outcome`, and `agent-selection`;
`--filter` accepts a comma-separated subset of those names only. Missing or
unknown Run IDs and invalid filters exit `2`; stream/store errors exit `1`;
interrupting `--follow` exits `130`. A terminal Run replays and exits `0`. Use
`events` for automation, `attach` for the human view, and the Detached Run
Console Log as a compact text record — not a state API.

For a Detached Run's stable terminal subscription, use:

```bash
roundfix events <run-id> --follow --filter outcome
```

The outcome record carries the terminal state plus bounded reason and next
action when non-Clean. When available, it also carries Review Issue knowledge,
Console Log, Attach command, accepted Evidence kind and head, and the verified
parent head used by artifact-only inheritance.

## stop

```bash
roundfix stop <run-id>
roundfix stop --force <run-id>
```

Selectors: positional `<run-id>`, `--run-id`, `--pr`, `--spec`, or
`--head-repo` plus `--head-branch`. Graceful stop records a Stop Request and
reports `Stop Request recorded; the Run stops after the current Work Item
settles.`

During a watch Run's Review Source status, retry, quiet-period, or
merge-readiness wait, the owner checks for the Stop Request before the next
status access and after each interruptible sleep. It reaches Stopped by the
next configured poll boundary. After detecting the request, it does not run
another fetch, check, commit, push, or Review Source mutation. A Work Item
already in flight still settles before the graceful stop completes.

Force Stop is for dead, stuck, or runaway Runs. It cancels only registered
Agent Sessions whose current lifecycle is active, then terminates the recorded
owner process and proves that process exited. Only then does Roundfix report
Stopped, release the Active Run lock, and reap kept terminal Worktrees whose
branch has no commits beyond its base. An already-absent registered Agent
Session is an idempotent cleanup result.

If owner exit cannot be proven, Force Stop fails with no stdout success report.
The diagnostic names the Run ID, owner PID, and failed process-control step;
the Run remains Active and its Active Run lock stays retained. Inspect the Run
with `roundfix runs list --state active`, resolve the reported owner-process
failure, then retry `roundfix stop --force <run-id>`. Agent Session cleanup
failures remain visible as secondary warnings after the primary failure. They
do not replace that failure or authorize terminal completion while the owner
is still alive.

Exit codes: `0` for a recorded Stop Request, a completed Force Stop, and the
idempotent already-Stopped report; `1` when Force Stop fails operationally
because owner exit cannot be proven; `2` for Preflight Validation failures
such as an invalid selector, no matching Active Run, or stopping a Run that
already holds a different terminal outcome.

Terminal results are stable. Repeating Force Stop for an already Stopped Run
reports the existing outcome without repeating process or Agent Session
actions. Force Stop against a different terminal outcome is rejected and
leaves that outcome unchanged.

Orphaned locks rarely need `--force` anymore: Runs record their owner process
id, and any command blocked by a lock whose owner is provably dead reclaims it
automatically — the Run completes Failed with the reason journaled and one
stderr warning names the reclaimed run id. A live owner, a PID-less legacy
Run, or any liveness result short of proof still blocks; a warning alone never
authorizes owner reclamation.

The terminology and behavior trace to the
[Roundfix glossary](../../CONTEXT.md#language),
[ADR-0052](../adr/0052-run-completion-is-compare-and-set.md),
[Spec 0037](../specs/_archived/0037-terminal-outcome-integrity/_prd.md), and the
[detached-watch finding](../findings/2026-07-16-vortex-pr87-detached-watch-notification.md#4-cleanup-noise-appeared-before-the-actionable-failure).

## Detached Runs

`--detach` is available on `resolve`, `watch`, and `implement`. The foreground
command prints exactly five stdout lines and exits `0`:

```text
Run ID: <run-id>
Console Log: <path>
Attach: roundfix attach <run-id>
Supervisor monitor: roundfix events <run-id> --follow --filter outcome
Stop: roundfix stop <run-id>
```

The handshake is two-phase: the child writes a liveness marker immediately on
entering child mode, before configuration load and Preflight Validation, and
the run-id line once the Run exists. The parent waits 10 seconds for liveness
and up to 5 minutes for Run creation, so a slow but healthy preflight (a real
Agent probe takes seconds) never fails a detach start. Every failure branch
prints an explicit stderr diagnostic — for example:

```text
roundfix: Detached Run child produced no liveness signal within 10s; killed (exit: <exit or signal>)
```

— followed by the child's console output when any exists. The detached child
owns the terminal outcome. Supervisors use the printed outcome command;
humans use `roundfix attach <run-id>` for the read-only Live Run View. Detach
implies non-interactive mode: `--interactive` is rejected and `--no-input` is
implied.

Run Outcome Notification delivery is best-effort and never changes the Run
outcome or exit code. Each attempt appends a separate durable Run Event with
route, completion time, and receipt status `sent`, `skipped`, or `failed`;
receipt success means the local route accepted the request, not that a person
saw it. The original `notify.command` variables remain available:
`ROUNDFIX_RUN_ID`, `ROUNDFIX_OUTCOME`, `ROUNDFIX_KIND`, and
`ROUNDFIX_TARGET`. Terminal context adds `ROUNDFIX_REASON`,
`ROUNDFIX_CONSOLE_LOG`, `ROUNDFIX_ATTACH_COMMAND`,
`ROUNDFIX_REVIEW_ISSUES_KNOWN`, and `ROUNDFIX_NEXT_ACTION`.

The review Evidence, artifact inheritance, Detached outcome, and notification
contracts trace to the [Roundfix glossary](../../CONTEXT.md#language),
[ADR-0054](../adr/0054-review-source-evidence-determines-review-outcomes.md),
[Spec 0039](../specs/0039-review-source-evidence-and-detached-outcomes/_prd.md),
and the
[detached-watch finding](../findings/2026-07-16-vortex-pr87-detached-watch-notification.md).

## Agent boundaries

Inside a Run, Agents own only assigned issue or task files, triage, code
edits, tests, verification commands, and assigned status updates. They must
not commit, push, resolve Review Source threads, edit unassigned files, or
mark issues `duplicated` — the Daemon owns every one of those boundaries.
