# Roundfix — operational guide

Roundfix runs two local-first loops from the terminal: implementing a Spec's
Task Graph, and resolving CodeRabbit review feedback on a pull request. This
guide is the operational path for each — for a human at a prompt and for an
agent driving Roundfix. For flags, outputs, and boundaries per command, see the
[command reference](commands.md); for config keys and local state, see
[configuration](configuration.md); for install, see the
[README](../../README.md#install).

## Before you start

1. Install Roundfix (npm launcher or `make build`) and put it on `PATH`.
2. Make the machine Run-ready and check it:

   ```bash
   roundfix setup      # proves adapters and generated profiles before writes
   roundfix doctor     # read-only profile and Repository Skill Set readiness
   ```

   Doctor prints Agent Selection Profile Readiness first, then independently
   proves the Repository Skill Set. A ready repository prints
   `skills: ok (39 required: 14 Roundfix-owned, 25 external)`; a blocking
   mismatch prints `skills: failed` with the applicable owned or external
   update command and exits `1`. Doctor is offline and read-only: it never
   updates or deletes skills, and it ignores unrelated extra installed skills
   and lock entries. See the [Doctor Command reference](commands.md#doctor)
   for the ownership authorities and exact remediation commands.

3. Authenticate the GitHub CLI for the repository (`gh auth status`). Review
   loops need it; the implement loop does not.
4. Work on a non-default branch. `implement` refuses to run on the repository's
   default branch.

## The output contract

Both loops keep the same contract, which is what lets an agent drive them:

- **stdout** carries only the deterministic report at Run end.
- **stderr** carries diagnostics, progress, the Run ID, and Agent output.
- **Exit codes**: `0` Clean, Stopped, or an already-complete no-op; `1`
  Unresolved, Failed, or Integration Pending; `2` Preflight Validation failure;
  `3` Clean Unverified (watch only); `130` for an in-terminal Ctrl-C.
- **Outcomes are fixed strings** you can branch on — see each loop below.

Agent Selection Profiles choose the runtime, model, and reasoning effort.
Complete one-Run overrides can use `codex`, `claude`, or `opencode`. The Review
Source is `coderabbit`.

## Agent Selection Profiles

Roundfix routes Agent work through Agent Selection Profiles. Each profile is an
atomic object: one Preferred Selection plus a non-empty ordered Fallback Chain.
Project Config wins over User Config, which wins over built-ins; when a higher
scope defines a profile, Roundfix uses that whole object and never merges tuple
fields or fallback entries across scopes.

The required built-in categories are `general`, `backend`, `frontend`, `qa`,
and `review`. Optional Task Type categories `data`, `infra`, `docs`, `test`,
and `chore` inherit the effective `general` profile when absent; if you define
one, it must be complete. Built-ins use official model identifiers:

- `general`, `backend`, `qa`, and `review`: preferred
  `codex / gpt-5.6-sol / high`, fallback
  `codex / gpt-5.5 / xhigh`.
- `frontend`: preferred `claude / opus / xhigh`, fallback
  `codex / gpt-5.6-sol / high`.

Use this complete Project Config or User Config shape when you want explicit
profiles for every required category:

```yaml
profiles:
  general:
    preferred:
      runtime: codex
      model: gpt-5.6-sol
      reasoning_effort: high
    fallbacks:
      - runtime: codex
        model: gpt-5.5
        reasoning_effort: xhigh
  backend:
    preferred:
      runtime: codex
      model: gpt-5.6-sol
      reasoning_effort: high
    fallbacks:
      - runtime: codex
        model: gpt-5.5
        reasoning_effort: xhigh
  frontend:
    preferred:
      runtime: claude
      model: opus
      reasoning_effort: xhigh
    fallbacks:
      - runtime: codex
        model: gpt-5.6-sol
        reasoning_effort: high
  qa:
    preferred:
      runtime: codex
      model: gpt-5.6-sol
      reasoning_effort: high
    fallbacks:
      - runtime: codex
        model: gpt-5.5
        reasoning_effort: xhigh
  review:
    preferred:
      runtime: codex
      model: gpt-5.6-sol
      reasoning_effort: high
    fallbacks:
      - runtime: codex
        model: gpt-5.5
        reasoning_effort: xhigh
```

The Codex Model Catalog recognizes `gpt-5.6-sol`, `gpt-5.6-terra`, and
`gpt-5.6-luna` as official identifiers; GPT-5.5/xhigh remains the generated
fallback for the four Codex-led profiles. A valid identifier and an advisory
recommendation rank are not readiness claims. Roundfix proves operational
availability in the effective environment through exact proof of the complete
runtime/model/reasoning tuple. Custom model values remain forward-compatible:
Roundfix sends them verbatim for the same proof instead of treating the catalog
as an allowlist. The Claude Model Catalog recognizes the identifiers the adapter
advertises: `opus`, `claude-fable-5`, `sonnet`, `haiku`, and `default`. It
previously listed `claude-opus-5` and `claude-opus-4-8`, which no adapter
advertises, and omitted three that it does; the catalog now follows the
adapter.

When an adapter advertises an independent reasoning control, Roundfix treats
every advertised Agent Model identifier as opaque. Copy bracketed identifiers
such as `opus[1m]` exactly as printed, or use the advertised canonical prefix
`opus`, and select reasoning separately. The `[1m]` suffix is a context-window
annotation, not a reasoning effort; `reasoning_effort: 1m` is rejected. See
[ADR-0079](../adr/0079-independent-reasoning-controls-make-model-identifiers-opaque.md).

Adapter Readiness requires the effective Codex command to prove official
`@agentclientprotocol/codex-acp` lineage at version `1.1.5` or newer and the
effective Claude command to prove official
`@agentclientprotocol/claude-agent-acp` lineage at version `0.63.0` or newer.
The deterministic install actions are
`npm install -g @agentclientprotocol/codex-acp@1.1.5` and
`npm install -g @agentclientprotocol/claude-agent-acp@0.63.0`.

A bare `codex-acp` override can resolve to legacy
`@zed-industries/codex-acp`; Setup diagnoses it and, after authorization,
migrates it to `npx -y @agentclientprotocol/codex-acp@1.1.5`. Setup also
migrates earlier explicit pins such as `1.1.4`. For Claude, Setup recognizes
both legacy lineages — the former `claude-code-acp` package and wrong-scope
`@zed-industries/claude-agent-acp` — and proposes
`npx -y @agentclientprotocol/claude-agent-acp@0.63.0`.

### Inspect profiles

`profiles show` is read-only and prints the effective source, inherited source,
Preferred Selection, fallback order, and five advisory recommendations:

```bash
roundfix profiles show
roundfix profiles show --category backend
roundfix profiles show --category backend --json
```

JSON uses schema `roundfix/profiles/v1`. Recommendations come from a
2026-08-07 five-entry snapshot. Each row includes benchmark, result, average
cost, source date, rationale, and `category_specific: false`. They are advisory
only: the list is not category-specific proof, not automatic routing input, and
never mutates User Config or Project Config.

### Configure profiles

`profiles configure` prepares the candidate in memory, validates it, proves
every distinct Preferred Selection and fallback, and only then asks for
confirmation:

```bash
roundfix profiles configure --scope project
roundfix profiles configure --scope user --file profiles.yml --dry-run --json
roundfix profiles configure --scope project --file profiles.yml --json
```

Without `--file`, Roundfix collects one complete profile through Interactive
Input, including at least one fallback, then shows the normalized profile and
target scope before asking for confirmation. With `--file`, the input can be a
strict fragment such as:

```yaml
backend:
  preferred:
    runtime: codex
    model: gpt-5.6-sol
    reasoning_effort: high
  fallbacks:
    - runtime: claude
      model: opus
      reasoning_effort: xhigh
```

`--dry-run` runs the same exact proof, reports `changed: false`, and leaves
files unchanged. `--json` uses schema `roundfix/profiles-configure/v1`.
Validation, proof, cleanup, declined-confirmation, and output failures all
preserve target bytes. Successful writes preserve unrelated config keys and
never edit runtime-owned settings, credentials, or adapter configuration.

### Validate profiles

`profiles validate` is read-only exact proof through disposable ACP Sessions:

```bash
roundfix profiles validate
roundfix profiles validate --category review --json
```

The command proves each distinct runtime/model/reasoning tuple once, reports
every category reference that uses it, sends no prompt, creates no Run, and
closes every disposable session on success or error. JSON uses schema
`roundfix/profiles-validate/v1`. A failed proof names the runtime, model,
reasoning effort, affected categories, adapter cause, and the next
`roundfix profiles configure` or `roundfix profiles validate` action.

### One-Run overrides and fallback boundaries

Agent-starting commands still accept one-Run `--agent`, `--model`, and
`--reasoning-effort` flags. Omit all three to use Agent Selection Profiles, or
provide all three together as one complete override. Every partial subset exits
`2` before configuration loading, exact proof, Session creation, worktree or
artifact creation, or Run persistence. An explicit empty
`--reasoning-effort ""` counts as present and requests model-managed reasoning;
Roundfix never substitutes it for a rejected explicit `high` request. A
complete override replaces only the Preferred Selection for each relevant
category and keeps its configured Fallback Chain. If one override applies
across multiple Task or QA categories, Roundfix emits a warning in text output
and JSON metadata.

Operational Runs prove every relevant preferred and fallback tuple before Run
creation. If a session cannot start after the Run exists, Roundfix records and
renders the fallback notification before creating or preparing the next fallback
session. Automatic fallback is allowed only before the first prompt. After
`agent_work_started`, prompt, tool, verification, cancellation, rate-limit, or
session-loss failures use normal Task or Run failure semantics; Roundfix never
starts a replacement session over possibly changed state.

### Legacy migration

The legacy `defaults.agent` and `runtimes.<runtime>.model` /
`runtimes.<runtime>.reasoning_effort` keys remain readable only while a scope
has no `profiles` section. A same-scope mix fails with migration guidance. To
migrate, remove `defaults.agent` and `runtimes`, then write complete profiles:

```bash
roundfix profiles configure --scope project --file profiles.yml
roundfix profiles validate --json
```

Setup diagnoses stale Codex and Claude adapter overrides before this profile
migration. It proposes each official pinned command and asks before rewriting
the ACPX config; `--yes` authorizes the offers, while decline or `--no-input`
writes nothing. Setup also proves all generated profile tuples before any User
Config or Project Config write, so failed Sol/high proof never degrades to
model-managed reasoning and never leaves partial configuration.

## Release planning before publication

Release work starts with the read-only Release Plan Command:

```bash
roundfix release plan
```

Run it before changelog edits, version-file edits, tags, pushes, package
publication, asset uploads, or GitHub Release creation. A generic release
request authorizes only a conclusive patch plan. Minor, major, and version-zero
breaking plans require explicit human approval of the printed question. If the
plan reports `manual_classification_required`, rerun it with
`--impact <none|patch|minor|major> --reason <text>`; that classification records
the impact and reason, but it does not approve the resulting version.

## Loop 1 — context-driven implementation

Execute a Spec's Task Graph from the resolved Spec Root as one Run. The default
root is `docs/specs`, so the default Spec path is `docs/specs/<slug>/`.
Configure `specs.root` in User Config or Project Config to use a different
directory; relative values resolve against the repository root, and absolute
values are used as-is. Roundfix resolves that root once from the user's checkout
and carries it into Run and Task Worktrees, so every surface reads and writes the
same Spec artifacts. Tasks run in dependency order by Wave, each Task's
Verification commands gate one commit, and the Run never pushes unless
`implement.auto_push: true` and the outcome is Clean.

Each Task owns an Agent Session selected by its Task Type profile, so frontend
and non-frontend Tasks stay in the same mixed Task Graph. The Daemon alone
writes Task status during the Run. Agents hand back implementation-ready work
after any useful focused checks; the Daemon then runs the complete declared
Verification before it can settle the Task.

The recommended capacity split overlaps up to two Task lifecycles while
serializing the repository gate:

```yaml
worktree:
  concurrency: 2

verification:
  concurrency: 1
```

These limits apply to one Implement Run only. They do not coordinate another
Run, CI, or an external command. See
[Task and Verification capacities](configuration.md#task-and-verification-capacities)
for defaults, precedence, and validation.

### Foreground

```bash
roundfix implement --spec <slug>
```

stdout is one line per Task in Task Graph order, then one outcome line:

```text
task_01 completed — first task
task_02 completed — second task
Clean: all 2 Task(s) completed.
```

Other outcome lines:

```text
Unresolved: 1 completed, 1 failed, 0 skipped, 0 pending.
IntegrationPending: 2 completed, 0 failed, 0 skipped, 0 pending; integrate with git merge --ff-only roundfix/run-<id>
All 2 Task(s) already completed; no Run was created.
```

Add `--qa` to end the Run with the qa-gate step; only a `pass` verdict lets the
Run end Clean, and the report gains a `qa <verdict> — <report path>` line.

When Daemon Verification fails deterministically, the capacity permit is
released while the same Task Agent Session receives one Verification Feedback
repair turn. The final Daemon attempt queues for capacity again. A
project-authored Verification wrapper can instead return exit `75` to classify
a Temporary Verification Failure. Roundfix retains that diagnostic and grants
one exclusive retry across the Task lifecycle; another `75` exhausts the retry
and fails the Task. Log content never triggers this protocol.

When the resolved Spec Root is outside the repository, or a task file or QA
Report path crosses a symbolic link, Daemon commits leave that artifact out of
the code-repository commit. The progress stream prints warnings such as
`roundfix: task file <path> kept outside the repository; committed without it`,
and the Run Event Journal records the dropped path and reason. If a Task's only
change is its external task file, it still settles `completed` without a commit;
an external QA Report likewise leaves the QA step free to proceed. Remove any
temporary git shims that hid symlink pathspec failures after upgrading to a
Roundfix build with this behavior.

### Detached, then monitor

For scripts, CI, or when you must not own the Run's lifetime, detach. This is the
default loop for an agent:

```bash
roundfix implement --spec <slug> --detach
```

Detach prints exactly five stdout lines and exits `0`:

```text
Run ID: <run-id>
Console Log: <path>
Attach: roundfix attach <run-id>
Supervisor monitor: roundfix events <run-id> --follow --filter outcome
Stop: roundfix stop <run-id>
```

Monitor without owning the Run. If you have the detached report, use the
captured Run ID and the stable Supervisor outcome command:

```bash
roundfix events <run-id> --follow --filter outcome
```

Use the Verification stream when you need working, queue, retry, and capacity
evidence:

```bash
roundfix events <run-id> --follow \
  --filter task-status,verification,outcome > run-events.jsonl \
  2> run-events.diagnostics
```

The requested JSONL records stay on stdout and diagnostics stay on stderr.
Verification phases appear as `waiting`, `started`, `command-passed`, `failed`,
and `verdict`; exclusive retry records use `mode: "exclusive"` and `retry: 1`.

At an interactive terminal, browse with the Run Browser; from a script or
agent, use the bounded plain-text listing:

```bash
roundfix runs                     # Run Browser: browse, Enter attaches read-only
roundfix attach                   # same Run Browser, q or Esc quits with no side effects
roundfix runs list                # agent report: 20 newest Active Runs, newest first
roundfix runs list --state all --limit 0   # widen the filter and the bound
roundfix attach <run-id>          # direct read-only Live Run View
# or tail the console log at <artifact-dir>/runs/<run-id>/console.log
```

The Run Browser is machine-wide: every repository's Runs newest first with a
repository column, Active Runs only by default; `a` toggles active/all,
`Enter` opens the read-only Live Run View, and leaving it returns to a
refreshed browser. `runs list` defaults to
the 20 newest Active Runs; widen with `--state <active|terminal|all>` and
`--limit N` (`0` unbounded), and read the single trailing stderr note that
names hidden Runs and the widening flag. The terminal outcome line
lands in the console log, and the detached child sends the configured outcome
notification when the Run reaches its terminal outcome. Delivery is
best-effort and cannot change the Run outcome. A separate durable Run Event
receipts each attempt as `sent`, `skipped`, or `failed`, with its route and
completion time; acceptance by a local route does not prove a person saw it.
Use the Supervisor outcome command as the unattended-Run state contract and
`attach` for the complete read-only Run history.
`attach` never stops, commits, or mutates the Run; detaching leaves it running.
For spec Runs the Live Run View reports both effective capacities and labels
each Task `Agent working`, `Waiting for Verification`, or `Verifying`. After a
deterministic failure releases capacity for Verification Feedback, that Task
returns to `Agent working` until its final attempt queues.

### Reconcile retained Run work

Runs List keeps requested rows on stdout. When terminal spec Runs retain a Run
Worktree or Run Branch in the selected repository scope, it writes the exact
count and `roundfix reconcile` guidance to stderr without classifying the
work:

```text
(2 terminal Run Worktrees retained; run 'roundfix reconcile' to inspect)
```

Inspect before applying anything. A Run ID selects one terminal spec Run;
omitting it scans the current repository:

```bash
roundfix reconcile <run-id>                 # single-Run dry-run
roundfix reconcile                          # repository-wide dry-run
roundfix reconcile <run-id> --format json   # requested JSON on stdout
```

The report classifies every selected Run as `safe`, `unintegrated`, `dirty`,
`unknown`, or `released`. `safe` requires resolved Run and target heads,
positive cleanliness for any present registered Run Worktree, and proof that
the Run head is an ancestor of the target head. `unintegrated` has clean,
resolved evidence but lacks ancestry; `dirty` has tracked or untracked
worktree changes; `unknown` lacks trustworthy metadata or Git proof; and
`released` means both the Run Worktree and Run Branch are absent.

Only `safe` is eligible for cleanup. Age, terminal outcome, or one missing path
does not establish safety, and Roundfix preserves dirty, unintegrated, and
unknown work. Redirect requested output and diagnostics separately in scripts:

```bash
roundfix reconcile --format json > reconcile.json
roundfix runs list > runs.txt 2> runs.diagnostics
```

After reviewing the current dry-run, opt into mutation explicitly:

```bash
roundfix reconcile <run-id> --apply
roundfix reconcile --apply
```

`--apply` is the only mutation switch; there is no force bypass. Roundfix
revalidates cleanliness and both heads, then releases only entries that remain
`safe` or `superseded`. A safe Integration Pending Run becomes Clean with its evidence
recorded before cleanup. Other terminal outcomes remain unchanged, and a
second run reports `released` without another mutation.

See the [Reconcile Command reference](commands.md#reconcile) for the full
state table, stdout and stderr contract, JSON fields, refusal behavior, and
links to the glossary, ADR, Spec, and finding trail.

### Read the outcome and act

| Outcome line | Meaning | Next action |
| --- | --- | --- |
| `Clean: all N Task(s) completed.` | Every Task passed and integrated onto the current branch | Advance to the next Spec |
| `All N Task(s) already completed; no Run was created.` | Nothing to do | Advance |
| `IntegrationPending: … git merge --ff-only roundfix/run-<id>` | Tasks done, current branch could not fast-forward | Run the printed command from the repo root, then continue |
| `Unresolved: X completed, Y failed, …` | One or more Tasks did not settle | Recover the failed Tasks (below) |

Integration Pending usually means the working tree drifted while the Run was
Active. Do not edit files a Run is touching while it is Active.

### Recover a failed Task

Read the per-Task status lines for `task_NN failed — <title>`. For each failed
Task, inspect its kept Task Worktree or Run Worktree, then recover only that Task
once its Verification passes there:

```bash
roundfix settle --spec <slug> --task <task_id>
```

`settle` picks the first surface where the Task is actually `failed` (Task
Worktree, then Run Worktree, then the current repository), reports the choice
as `Settle surface: <path>`, re-runs the Task's Verification there, and on
pass prints one `commit <path>` line per committed path before committing and
integrating onto the Run Branch — warning when other failed Tasks share the
worktree, since their work is swept into the same commit. It creates no Run
and never pushes. Re-run `roundfix implement --spec <slug>` to pick up
still-pending Tasks; completed Tasks are skipped.

### Advance

When a Spec ends Clean and QA passes, archive it inside the resolved Spec Root
and move on:

```bash
roundfix archive <slug>
```

## Loop 2 — PR review resolution

Resolve unresolved CodeRabbit findings on an Open Pull Request. Roundfix fetches
them as local Review Issue artifacts, assigns bounded Batches to the Agent,
verifies changes, commits, resolves the source threads, and pushes only when no
Unresolved Review Issues remain. Review Runs execute in your checkout on the PR
Head Branch — no Run Worktree, so fixes are always a delta over the published
HEAD — and Preflight Validation requires a clean tracked working tree
(untracked files are allowed). Before any fetch or Agent work, the Branch
Integrity Preflight integrates or refuses pending `roundfix/run-*` work and
blocks while another Active Run holds the branch; `--skip-branch-integrity`
bypasses it only by publishing a PR audit comment. Each settled issue
propagates to GitHub individually — non-resolved outcomes get an explanatory
Outcome Comment. The Run's review artifacts ride the push in one separate docs
commit (`docs: review round NNN for pr <n>`); artifact roots outside the
repository are never staged.

### One shot: watch until clean

```bash
roundfix watch --source coderabbit --pr <number> --until-clean
```

`watch` owns the waits, fetches, Rounds, Agent lifecycle, verification, commits,
Final Push, source resolution, retries, and timeouts. With `--until-clean` it
ends Clean only after no Unresolved Review Issues remain and accepted Review
Source Evidence verifies the pushed HEAD. Accepted current-head Evidence can
be a successful CodeRabbit check or commit status, or a CodeRabbit
`APPROVED` review; every form requires zero unresolved CodeRabbit threads.
Stale signals never verify the expected head, and unresolved threads leave a
successful signal only `reviewed`.

The pre-fetch wait is `WaitingForReview`; post-push Merge-Ready confirmation is
`WaitingForReviewCheck`. Both expose a start, deadline, expected head, Evidence
state and kind, and retry status. Roundfix retries only positively typed
transient conditions — timeout outside Run cancellation, temporary DNS,
connection reset, HTTP `429`, or GitHub `5xx` — and records one
started/recovered-or-exhausted episode. Existing polling, Review Source
timeout, and Run Budget bounds apply; no retry configuration or log-text
inference is added.

An exact Daemon-created artifact-only descendant can inherit its verified
parent Evidence only when the recorded identity and sole parent match, the
current non-empty diff stays entirely under the resolved in-repository review
root without a symbolic-link crossing, and refreshed parent Evidence still
has no unresolved CodeRabbit threads. Any non-Daemon, user-authored, empty,
mixed-path, out-of-root, wrong-parent, stale, or unresolved descendant returns
to normal current-head Evidence polling.

A missing accepted signal within `watch.check_grace_period` (default `5m`)
ends Clean Unverified with exit `3`; confirm the Review Source Evidence before
merging. An explicit structured skip ends Review Skipped with exit `3`, prints
the source reason and next action, fetches no Round, and never means Clean,
Clean Unverified, or zero Review Issues. Bound the loop with
`--max-rounds <number>`.

stdout is one line per Review Issue in Round order — failed, unresolved, and
invalid lines carry a ` — reason: <terminal_reason>` suffix when the artifact
has one — then two summary lines separating this Run from the PR's cumulative
history:

```text
issue 001 resolved — major: handle test issue
This Run (Clean after 1 Round(s)): 1 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved.
Pull Request cumulative: 1 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved.
```

Review Issue statuses are `resolved`, `invalid`, `failed`, `duplicated`, or
`unresolved`. Non-clean outcomes include Clean Unverified
(`CleanUnverified`), Review Skipped (`ReviewSkipped`), `MaxRoundsReached`,
`BudgetExceeded`, `TimedOut`, `Failed`, and `Stopped`. If fetch does not
complete, Roundfix omits all zero-valued count summaries and reports:

```text
Review Issues: unknown — fetch did not complete.
```

The full evidence and refusal rules trace to
[ADR-0054](../adr/0054-review-source-evidence-determines-review-outcomes.md),
[Spec 0039](../specs/_archived/0039-review-source-evidence-and-detached-outcomes/_prd.md),
and the
[detached-watch finding](../findings/_archived/2026-07-16-vortex-pr87-detached-watch-notification.md).

### Step by step

Split the loop when you want to inspect artifacts between stages:

```bash
roundfix fetch --source coderabbit --pr <number>   # write Review Issue artifacts only
roundfix resolve --pr <number>                     # resolve downloaded issues once
```

`fetch` starts no Agent and never commits or pushes. `resolve` runs one Round
over already-downloaded issues and uses the same report shape as `watch` with
`1 Round(s)`.

### Detach a review Run

`--detach` works on `resolve` and `watch` too, with the same five-line report and
`attach`/`stop` follow-up as the implement loop.

## Driving Roundfix from an agent

Roundfix ships a bundled `roundfix` skill that encodes these loops for coding
agents. Install the skill set into a repository or an Agent's skill directory:

```bash
roundfix skills list       # bundled skills + recommended external skills
roundfix skills install    # writes to <repo>/.agents/skills
```

An agent driving Roundfix must:

- Prefer `roundfix` commands over manual GitHub scraping.
- Detach the Run, then discover it with the bounded `roundfix runs list`
  (Active Runs by default; widen with `--state all` or `--limit 0`) and
  follow its terminal outcome with
  `roundfix events <run-id> --follow --filter outcome`. The Run Browser and
  `attach` are human read-only surfaces; the Console Log is not a state API.
- Branch on the deterministic outcome line and exit code, not on log scraping.
- Report the Run ID, the PR or Spec, the Agent, and the current Run state when
  summarizing progress.
- Never commit, push, or resolve source threads by hand inside a Run — Roundfix
  owns those boundaries.

Inside a daemon-assigned Batch, the agent follows the bounded Review Issue or
Task contract in the skill: it edits only assigned issue files, updates only
assigned statuses, and never marks an issue `duplicated`.

## Stopping a Run

```bash
roundfix stop --spec <slug>              # graceful: stops after the current Work Item settles
roundfix stop --pr <number>              # graceful, by PR
roundfix stop --force <run-id>           # dead, stuck, or runaway Runs only
roundfix stop --force --owner-identity-unreadable <run-id> # unreadable identity only; last resort
```

Graceful stop records a Stop Request and lets the in-flight Work Item finish its
verification and commit boundary. During a Review Source status, retry,
quiet-period, or merge-readiness wait, the owner observes the request by the
next configured poll boundary. Once observed, the Run does not start another
fetch, check, commit, push, or Review Source mutation.

Use Force Stop only for a dead, stuck, or runaway Run. It cancels registered
active Agent Sessions, terminates the recorded owner process, and reports
Stopped only after owner exit is proven. Roundfix releases the Active Run lock
only after that proof. Registered sessions that are already absent need no
recovery; other cleanup failures appear as secondary warnings after the
primary failure.

Roundfix reads owner identity directly from the kernel and spawns no
subprocess, so the proof remains available when the host cannot fork. If
identity capture fails when a Run starts, one warning names its PID-only reuse
protection. The durable `owner_identity_unproven=true` marker appears in
`roundfix runs list`.

Force Stop refuses a proven identity mismatch: investigate PID reuse and do not
signal the process now using that PID. It also fails closed when owner identity
is unreadable; resolve the reported host resource failure and retry. Only when
that normal Force Stop specifically reports an unreadable identity may an
operator use `--owner-identity-unreadable` as a last resort. The flag authorizes
PID-only termination for that condition and never applies to a proven mismatch.
If identity is readable or proves a mismatch, Roundfix exits `2` without
signaling the process.

If Force Stop cannot prove owner exit, it fails closed: the Run remains Active
and the Active Run lock stays retained. Read the reported owner PID and failed
step, inspect the Run, resolve the process-control failure, and retry:

```bash
roundfix runs list --state active
roundfix stop --force <run-id>
```

Repeating Force Stop after the Run is already Stopped is an idempotent report
of the same outcome and does not repeat cleanup. If the Run already has a
different terminal outcome, Roundfix rejects the conflict and preserves that
outcome. Never kill Agent or acpx processes by hand while a Run is Active.

For the full failure and replay contract, see the
[Stop Command reference](commands.md#stop), which traces to
[ADR-0052](../adr/0052-run-completion-is-compare-and-set.md) and the
[terminal-outcome Spec](../specs/_archived/0037-terminal-outcome-integrity/_prd.md).

## Command reference

| Command | Purpose |
| --- | --- |
| `setup` / `doctor` | Prepare / diagnose machine readiness |
| `init` | Create User or Project Config |
| `implement` | Execute a Spec's Task Graph as one Run |
| `settle` | Recover one failed Task from its kept worktree |
| `archive` | Archive a completed, QA-passed Spec |
| `release plan` | Classify committed release changes without mutating release state |
| `profiles` | Show, configure, and validate Agent Selection Profiles |
| `fetch` | Download Review Issue artifacts for a PR |
| `resolve` | Resolve downloaded Review Issues once |
| `watch` | Fetch and resolve in a watched loop |
| `runs` | Browse Runs in the read-only Run Browser (interactive terminal) |
| `runs list` | List Runs from the Run Database, bounded and plain-text |
| `reconcile` | Classify retained terminal Run work and explicitly release proven safe or superseded entries |
| `attach` | Browse or replay a Run's timeline, read-only |
| `stop` | Request or force-stop an Active Run |
| `gc` | Prune old terminal Run journals and artifacts |
| `upgrade` | Upgrade the binary from GitHub Releases |
| `skills` | List, check, or install the bundled skills |

Run `roundfix <command> --help` for the full flag list of any command, or read
the [command reference](commands.md) for output shapes and boundaries. Local
state locations (Run Database, Worktrees, Artifact Directory) are documented
in [configuration](configuration.md#local-state).
