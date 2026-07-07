# Roundfix — operational guide

Roundfix runs two local-first loops from the terminal: implementing a Spec's
Task Graph, and resolving CodeRabbit review feedback on a pull request. This
guide is the operational path for each — for a human at a prompt and for an
agent driving Roundfix. For flags and outcomes per command, see the
[README Commands](../README.md#commands) and
[Command Boundaries](../README.md#command-boundaries) sections; for install, see
[README Install](../README.md#install).

## Before you start

1. Install Roundfix (npm launcher or `make build`) and put it on `PATH`.
2. Make the machine Run-ready and check it:

   ```bash
   roundfix setup      # verifies Node, acpx 0.12.0, the Agent probe, config
   roundfix doctor     # read-only readiness check; mutates nothing
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
  `130` for an in-terminal Ctrl-C.
- **Outcomes are fixed strings** you can branch on — see each loop below.

An Agent runtime is selected with `--agent`; supported values are `codex`,
`claude`, and `opencode`. The Review Source is `coderabbit`.

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
roundfix implement --spec <slug> --agent codex
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
roundfix implement --spec <slug> --agent codex --detach
```

Detach prints exactly four stdout lines and exits `0`:

```text
Run detached: <run-id>
Console log: <path>
Follow: roundfix attach <run-id>
Stop: roundfix stop <run-id>
```

Monitor without owning the Run. If you have the detached report, use the
captured Run ID. From a fresh terminal, discover the repository's Runs first or
open the Attach picker:

```bash
roundfix runs list --active       # stable report: id, state, kind, target
roundfix attach                   # interactive picker; q or Ctrl-C detaches after attach
roundfix attach <run-id>          # direct read-only Live Run View
# or tail the console log at <artifact-dir>/runs/<run-id>/console.log
```

`runs list` prints this repository's Runs newest first, and `--active` filters
out terminal Runs. `attach` without a Run ID lists the repository's Runs in an
interactive terminal and accepts a number or Run ID. The terminal outcome line
lands in the console log. `attach` never stops, commits, or mutates the Run;
detaching leaves it running.

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

`settle` re-runs the Task's Verification in the kept surface, commits on pass,
and integrates onto the Run Branch. It creates no Run and never pushes. Re-run
`roundfix implement --spec <slug>` to pick up still-pending Tasks; completed
Tasks are skipped.

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
Unresolved Review Issues remain.

### One shot: watch until clean

```bash
roundfix watch --source coderabbit --pr <number> --agent codex --until-clean
```

`watch` owns the waits, fetches, Rounds, Agent lifecycle, verification, commits,
Final Push, source resolution, retries, and timeouts. With `--until-clean` it
ends Clean only after no Unresolved Review Issues remain and the Review Source
check on the pushed HEAD succeeds. Bound it with `--max-rounds <number>`.

stdout is one line per Review Issue in Round order, then one outcome line:

```text
issue 001 resolved — major: handle test issue
Clean after 1 Round(s): 1 resolved, 0 invalid, 0 failed, 0 unresolved.
```

Review Issue statuses are `resolved`, `invalid`, `failed`, `duplicated`, or
`unresolved`. Non-clean outcomes include `MaxRoundsReached`, `BudgetExceeded`,
`TimedOut`, `Failed`, and `Stopped`.

### Step by step

Split the loop when you want to inspect artifacts between stages:

```bash
roundfix fetch --source coderabbit --pr <number>   # write Review Issue artifacts only
roundfix resolve --pr <number> --agent codex       # resolve downloaded issues once
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
- Detach the Run, then discover it with `roundfix runs list --active` and
  follow it with `roundfix attach <run-id>` or the console log rather than
  blocking a foreground process.
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
| `fetch` | Download Review Issue artifacts for a PR |
| `resolve` | Resolve downloaded Review Issues once |
| `watch` | Fetch and resolve in a watched loop |
| `runs list` | List Runs from the Run Database |
| `attach` | Pick or replay a Run's timeline, read-only |
| `stop` | Request or force-stop an Active Run |
| `gc` | Prune old terminal Run journals and artifacts |
| `upgrade` | Upgrade the binary from GitHub Releases |
| `skills` | List, check, or install the bundled skills |

Run `roundfix <command> --help` for the full flag list of any command. Local
state locations (Run Database, Worktrees, Artifact Directory) are documented in
[README Local State](../README.md#local-state).
