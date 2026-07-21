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
   roundfix doctor     # read-only adapter and profile readiness; mutates nothing
   ```

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
- `frontend`: preferred `claude / claude-fable-5 / medium`, fallback
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
      model: claude-fable-5
      reasoning_effort: medium
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
as an allowlist.

Codex Adapter Readiness requires the effective command to prove the official
`@agentclientprotocol/codex-acp` lineage at version `1.1.4` or newer. A bare
`codex-acp` override can resolve to legacy `@zed-industries/codex-acp`; use
Setup to diagnose it and, after authorization, migrate to
`npx -y @agentclientprotocol/codex-acp@1.1.4`. The deterministic install action
is `npm install -g @agentclientprotocol/codex-acp@1.1.4`.

### Inspect profiles

`profiles show` is read-only and prints the effective source, inherited source,
Preferred Selection, fallback order, and five advisory recommendations:

```bash
roundfix profiles show
roundfix profiles show --category backend
roundfix profiles show --category backend --json
```

JSON uses schema `roundfix/profiles/v1`. Recommendations come from a
2026-07-16 five-entry snapshot. Each row includes benchmark, result, average
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
      model: claude-fable-5
      reasoning_effort: high
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

Setup diagnoses a stale Codex adapter override before this profile migration.
It proposes the official pinned command and asks before rewriting the ACPX
config; `--yes` authorizes the offer, while decline or `--no-input` writes
nothing. Setup also proves all generated profile tuples before any User Config
or Project Config write, so failed Sol/high proof never degrades to
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

Detach prints exactly four stdout lines and exits `0`:

```text
Run detached: <run-id>
Console log: <path>
Follow: roundfix attach <run-id>
Stop: roundfix stop <run-id>
```

Monitor without owning the Run. If you have the detached report, use the
captured Run ID. At an interactive terminal, browse with the Run Browser;
from a script or agent, use the bounded plain-text listing:

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
notification when the Run reaches its terminal outcome. Treat that notification
as the unattended-Run signal; use `attach` or the console log for details.
`attach` never stops, commits, or mutates the Run; detaching leaves it running.

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
ends Clean only after no Unresolved Review Issues remain and the Review Source
check on the pushed HEAD succeeds; a check that never appears within
`watch.check_grace_period` (default `5m`) ends `CleanUnverified` with exit
code `3` — confirm the check yourself before merging. Bound the loop with
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
`unresolved`. Non-clean outcomes include `CleanUnverified`,
`MaxRoundsReached`, `BudgetExceeded`, `TimedOut`, `Failed`, and `Stopped`.

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

`--detach` works on `resolve` and `watch` too, with the same four-line report and
`attach`/`stop` follow-up as the implement loop.

## Driving Roundfix from an agent

Roundfix ships a bundled `roundfix` skill that encodes these loops for coding
agents. Install the skill set into a repository or an Agent's skill directory:

```bash
roundfix skills list       # bundled skills + recommended external skills
roundfix skills install    # writes to <repo>/.agents/skills
```

An agent driving Roundfix should:

- Prefer `roundfix` commands over manual GitHub scraping.
- Detach the Run, then discover it with the bounded `roundfix runs list`
  (Active Runs by default; widen with `--state all` or `--limit 0`) and
  follow it with `roundfix attach <run-id>` or the console log rather than
  blocking a foreground process. The Run Browser is the human surface; agents
  stay on the plain-text listing.
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
```

Graceful stop records a Stop Request and lets the in-flight Work Item finish its
verification and commit boundary. `--force` cancels the Agent Session
best-effort, completes the Run Stopped immediately, releases its locks, and reaps
empty terminal Worktree debris. Never kill Agent or acpx processes by hand while
a Run is Active.

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
| `attach` | Browse or replay a Run's timeline, read-only |
| `stop` | Request or force-stop an Active Run |
| `gc` | Prune old terminal Run journals and artifacts |
| `upgrade` | Upgrade the binary from GitHub Releases |
| `skills` | List, check, or install the bundled skills |

Run `roundfix <command> --help` for the full flag list of any command, or read
the [command reference](commands.md) for output shapes and boundaries. Local
state locations (Run Database, Worktrees, Artifact Directory) are documented
in [configuration](configuration.md#local-state).
