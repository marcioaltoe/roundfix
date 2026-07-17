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

Verifies Node.js, the pinned acpx version, the configured Agent probe, acpx
local adapter overrides, User Config, and Project Config. Each check prints one
deterministic report line such as `node: ok`, `acpx: installed`, or
`User Config: skipped`. `--yes` accepts every offered install or file change;
`--no-input` skips offers instead of prompting.

### doctor

```bash
roundfix doctor
```

Read-only readiness report; mutates nothing and exits nonzero when any check
fails. One stdout line per check with `ok`, `failed`, or `skipped`; failure
lines include `next: <action>` when a remediation is known. The checks:

- `node:` — Node.js meets the minimum version.
- `acpx:` — the pinned acpx version is installed.
- `adapter:` — the configured runtime's ACP adapter binary resolves on `PATH`;
  a missing binary fails with its install command.
- `agent:` — the configured Agent probe succeeds.
- `model:` — the runtime accepted the effective Agent Model; on failure the
  line carries the runtime's currently advertised models and a `next:` action.
- `codex:` — macOS-only runtime hygiene: inspects `com.apple.quarantine` (the
  real XProtect trigger) and code-signature validity, resolving `CODEX_PATH`
  first and then `codex` on `PATH`. It does not use `spctl --assess`, which
  rejects any signed CLI that is not a notarized app. A quarantined or
  improperly-signed codex fails with the next action to reinstall codex with
  the official curl installer into `~/.local/bin` and set `CODEX_PATH`.
  Skipped on non-Darwin platforms.

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
roundfix resolve --pr <number> --agent codex [--spec <slug>]
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
roundfix watch --source coderabbit --pr <number> --agent codex --until-clean [--max-rounds N]
```

Waits for CodeRabbit status on the current PR HEAD, observes the configured
quiet period, fetches, resolves Batches, and repeats. After the Final Push it
polls for the Review Source check on the pushed HEAD through
`watch.check_grace_period` (default `5m`): a successful check with no new
Review Issues ends `Clean`; new Review Issues start the next Round; a check
that never appears ends `CleanUnverified` with exit code `3` and a report
naming the next action. Other terminal outcomes are `MaxRoundsReached`,
`BudgetExceeded`, `TimedOut`, `Failed`, and `Stopped`.

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

## Spec loop: implement, settle, archive

### implement

```bash
roundfix implement --spec <slug> --agent codex [--qa] [--detach]
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

## stop

```bash
roundfix stop <run-id>
roundfix stop --force <run-id>
```

Selectors: positional `<run-id>`, `--run-id`, `--pr`, `--spec`, or
`--head-repo` plus `--head-branch`. Graceful stop records a Stop Request and
reports `Stop Request recorded; the Run stops after the current Work Item
settles.` `--force` is for dead, stuck, or runaway Runs: it cancels the Agent
Session best-effort, completes the Run Stopped immediately, releases its
locks, and reaps kept terminal Worktrees whose branch has no commits beyond
its base.

Orphaned locks rarely need `--force` anymore: Runs record their owner process
id, and any command blocked by a lock whose owner is provably dead reclaims it
automatically — the Run completes Failed with the reason journaled and one
stderr warning names the reclaimed run id. A live owner, a PID-less legacy
Run, or any liveness result short of proof still blocks.

## Detached Runs

`--detach` is available on `resolve`, `watch`, and `implement`. The foreground
command prints exactly four stdout lines and exits `0`:

```text
Run detached: <run-id>
Console log: <path>
Follow: roundfix attach <run-id>
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
owns the terminal outcome and fires the outcome notification; monitor it with
`roundfix events <run-id> --follow` or `roundfix attach <run-id>`. Detach
implies non-interactive mode: `--interactive` is rejected and `--no-input` is
implied.

## Agent boundaries

Inside a Run, Agents own only assigned issue or task files, triage, code
edits, tests, verification commands, and assigned status updates. They must
not commit, push, resolve Review Source threads, edit unassigned files, or
mark issues `duplicated` — the Daemon owns every one of those boundaries.
