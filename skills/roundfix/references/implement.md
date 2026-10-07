## User-Facing Spec Runs

The Implement Command executes a Spec's Task Graph on the current branch as
one Run. Tasks whose dependencies are completed form the current Wave; the
scheduler starts up to `worktree.concurrency` Tasks from that Wave at once.
The default is `2`; `1` keeps sequential behavior. Each Task's Verification
commands gate one commit. By default the Run never pushes; a repository can
opt in with `implement.auto_push: true`, which pushes only after a Clean
outcome and never opens pull requests (ADR-0138).

When the Run Budget is enabled, an Implement Run starts with the configured
maximum Run duration and its allowance renews at each Task settlement. The
renewed allowance bounds the next Task or QA gate and post-cycle integration,
push, and cleanup. When the allowance expires, the Run settles
`BudgetExceeded` with a reason naming both the configured maximum and the
elapsed time since the renewal point. The bounded Run preserves its Run
Worktree and Run Branch for inspection and recovery.

Before creating a Run, `implement` inspects prior terminal Runs for the same
Spec in the current repository. When a `BudgetExceeded`, `Stopped`, or
`Unresolved` Run with a present Run Worktree has a complete candidate set that
would carry, Preflight Validation refuses before creating a Run or Agent
Session. The complete set
must pass Task Carry-Forward's existing proofs, including a passing
Verification verdict, exactly one settlement commit, and unmoved declared
inputs for each candidate. Input proofs use the checkout plus the accumulating
staged carries, so each later candidate is compared with the state established
by earlier carries rather than the raw checkout. A Task already completed on
the checkout is reported as `already completed; nothing to carry`, never
refuses the set, and is omitted from the Tasks named by the refusal. A Run
whose candidates are all already completed produces no not-available note. It
exits `2`, leaves stdout
empty, and writes no Git or Run Database state. The refusal names the Run and
Tasks to recover and gives the exact next action:

```bash
roundfix reconcile <run-id> --carry-forward
```

When no complete candidate set would carry, `implement` reports the inspection
result and proceeds. An inspection failure is reported and also lets the Run
proceed. If several complete sets qualify, `implement` selects the Run with
the largest carriable Task set, breaking ties with the newest Run.

1. Start the Implement Command with:

   ```bash
   roundfix implement --spec <slug>
   ```

2. Flags:
   - `--spec` — Spec slug under the resolved Spec Root (`docs/specs/` by
     default).
   - There is no QA parameter. The gate is authored into the Task Graph:
     the manifest's `qa: task_NN` names a terminal `type: qa` Task depending
     on every leaf, or `qa: declined` records the reasoned decline. The
     Daemon runs an authored gate after the last Task settles; the one
     declared-acceptance eligibility policy decides whether the newest QA Report
     is acceptable for settlement. A `pass` with no disallowed blocked rows, or
     a `partial` whose blocked rows are declared unreachable, settles the
     terminal `qa` Task. A `fail`, an undeclared `partial`, a missing or
     unparseable report, or a `pass` carrying finding-, declared-, or
     precondition-blocked rows still refuses; an environment-blocked row remains
     acceptable under the existing policy. Passing `--qa` is an
     unknown-flag error whose remediation names this contract. A reasoned
     `qa: declined` declaration creates no QA Task and executes no gate.
   - `--agent` — Agent runtime. Supported: `codex`, `claude`, `opencode`.
   - `--model` — Agent Model override.
   - `--reasoning-effort` — Default Reasoning Effort override.
   - `--agent-command` — Agent command override.
   - `--agent-full-access` — opt into Agent runtime full-access mode.
   - `--no-agent-console` — hide Agent-source console events from non-TTY
     stderr; the Run Event Journal is not filtered.
   - `--detach` — start a Detached Run and print the five-line report with
     Console Log, Attach, Supervisor monitor, and Stop actions.
   - `--interactive` — open Interactive Input before starting.
   - `--no-input` — fail instead of opening Interactive Input.

3. stdout carries only the deterministic report; diagnostics, the run id,
   and the agent log go to stderr:
   - One line per Task in Task Graph order: `task_NN <status> — <title>`,
     with status `completed`, `failed`, `skipped`, or `pending`.
     Failed and skipped Task lines are followed by one indented reason line:
     two spaces followed by `reason: <one line>`. Verification failure reasons name the failed
     command and exit status and point to diagnostics. Completed Task lines do
     not gain an extra line.
   - When the graph authors a gate, one verdict line after the Task lines:
     `qa <verdict> — <report path>`; a missing report prints
     `qa missing — no QA Report found`.
   - One outcome line: `Clean: all N Task(s) completed.`,
     `Unresolved: X completed, Y failed, Z skipped, W pending.`,
     `IntegrationPending: X completed, Y failed, Z skipped, W pending; integrate with git merge --ff-only roundfix/run-<id>`,
     or — when every Task is already completed and the graph has no unsettled
     gate, including a `qa: declined` declaration —
     `All N Task(s) already completed; no Run was created.`
     An all-completed graph whose authored gate is still unsettled is not a
     no-op: the Run starts and executes the gate.
   - When `implement.auto_push: true` and the Run ends Clean with an upstream
     branch, one final line follows the outcome: `pushed <remote>/<branch>`.
     A tested example is `pushed origin/ma/widget-flow`.

   When the resolved Spec Root is not the default `<repo>/docs/specs`, Run
   startup prints `Spec Root: <path>` on stderr.

   The spec Run header names both effective capacities with this shape:

   ```text
   Task Capacity: N
   Verification Capacity: M
   ```

4. Exit codes: `0` Clean, Stopped, or the all-completed no-op, `1`
   BudgetExceeded, Unresolved, Failed, or Integration Pending, `2` Preflight
   Validation failure, `130` for in-terminal Ctrl-C interrupt mapping.

5. Preflight Validation exits `2` with one actionable message when the Spec
   or its Task Graph is invalid (each failure names the offending Task or
   check), the current branch is the repository default branch, another Active
   Run holds the work target or working tree (the error names the run id and
   `roundfix stop <id>` unless the owner is proven dead and reclaimed
   automatically), or the Agent runtime probe fails. A dirty user
   checkout no longer blocks `implement`; stderr prints a note shaped like
   `roundfix: note: working tree <path> has N uncommitted change(s); implement will run in a Run Worktree, and overlapping local changes end the Run Integration Pending.`

   Before creating an Agent Session for a Task whose declared Verification
   contains the configured repository Verification command exactly, the Daemon
   runs that command once through shared Verification Capacity. Failure starts
   no Agent Session, settles the Task with a `repository not green on entry`
   reason, and publishes a Verification event classified `precondition` with
   reason `repository_not_green_on_entry`. A passing precondition does not
   replace post-Agent Verification, and a failing precondition does not consume
   Verification Feedback.

   Every Daemon Batch, Task, and QA commit filters stageable paths through the
   same boundary. Repository-external paths, symlink crossings, and regular
   files carrying any execute bit are omitted while remaining paths continue
   to the commit. An executable refusal names its path and mode:

   ```text
   roundfix: refused executable file <path> (mode <mode>); build artifacts and deliberately executable repository files are not valid Work Item output
   ```

   A task file or QA Report outside the repository, or reached through a
   symlinked path, is dropped from staging with one Run Event Journal entry
   naming the path and reason. Progress prints warnings shaped like:

   ```text
   roundfix: task file <path> kept outside the repository; omitted from the commit
   roundfix: QA Report <path> kept outside the repository; omitted from the commit
   ```

   If a Task commit has no change outside the Spec Root, including an empty
   stageable set, the Daemon still settles the Task `completed` but emits one
   stderr warning and one Run Event warning for the no-op shape. An external QA
   Report is left uncommitted and the QA step proceeds. Remove temporary git
   shims that hid symlink pathspec failures after upgrading to a Roundfix
   build with this behavior; those shims can mask regressions in the real
   commit boundary.

   At the QA Report commit boundary, a configured `verification.format` runs
   over the regular files under the Spec's `qa/` directory immediately before
   the QA Report commit. When an imported pass is present, the same command
   runs over its QA files after the mechanical stage and before the repository
   Verification precondition. Failed, timed-out, or verdict-changing runs
   restore the original bytes and record the format outcome.

6. Without `--spec`, Interactive Input lists the repository's active Specs
   from the resolved Spec Root under an `Active Specs:` picker that accepts a
   number or a slug, and the Agent field suggests the remembered Agent. Agent
   selection then asks for Agent Model and Default Reasoning Effort using the
   selected runtime's effective configuration as the default. Interactive
   Input asks nothing about QA: the gate is an authoring decision recorded in
   the Task Graph, not a per-run choice. The Agent is remembered across runs;
   the Spec slug and selection overrides never are. `--no-input` fails
   instead of opening Interactive Input.

7. Discover spec Runs with the bounded `roundfix runs list` or open the Run
   Browser with `roundfix attach` when the Run ID was not captured. Follow
   `roundfix events <run-id> --follow` for unattended JSONL monitoring. Attach
   directly with `roundfix attach <run-id>` for the Live Run View; it shows the
   Spec's Tasks as Work Items in the shared cockpit.

8. `implement.auto_push` is a bool in config, default `false`. User Config can
   provide a default, and Project Config can override it:

   ```yaml
   implement:
     auto_push: true
   ```

   The push uses the branch's detected upstream. Missing upstream prints one
   stderr note and leaves a Clean Run Clean. Integration Pending, Unresolved,
   Failed, Stopped, and failing-QA Runs do not invoke the pusher. Push failure
   ends the Run Failed.

9. Task Capacity and Verification Capacity are independent config-only limits
   for one Implement Run. `worktree.concurrency` is Task Capacity, defaults to
   `2`, and limits concurrent Task Worktree lifecycles.
   `verification.concurrency` is Verification Capacity, defaults to `1`, and
   limits concurrent Task Verification attempts. Both must be positive
   integers; Project Config overrides User Config, which overrides built-ins.
   Neither setting coordinates other Runs, CI, or external processes.

   The recommended split overlaps implementation while serializing the
   complete repository gate:

   ```yaml
   worktree:
     concurrency: 2

   verification:
     concurrency: 1
   ```

   Every normal attempt journals `waiting` before it acquires shared capacity,
   then `started`, command results, and `verdict`. A deterministic failure
   releases capacity before one Verification Feedback repair turn and
   reacquires capacity for the final Daemon attempt.

   Exit `75` from a project-authored Verification wrapper is the sole Temporary
   Verification Failure signal. Roundfix retains the diagnostics and grants
   the Task one exclusive retry across its Verification lifecycle. The retry
   waits for other Verification attempts in that Run to drain, consumes the
   entire Verification Capacity, and does not consume the Agent repair. A
   repeated exit `75` exhausts the retry and fails the Task. Roundfix never
   classifies a failure from logs, timing, ports, package names, or framework
   text.

   `worktree.location` is the parent directory, default
   `~/.roundfix/worktrees`; Roundfix always appends `<repo-slug>/<run-id>` or
   `<repo-slug>/<run-id>.<task_id>`. `worktree.copy` copies
   repository-relative, gitignored files into each new worktree.
   `worktree.bootstrap` runs in each new worktree after copy and before Agent
   work; `worktree.bootstrap_timeout` defaults to `10m`.

   ```yaml
   worktree:
     location: "~/.roundfix/worktrees"
     concurrency: 2
     copy: []
     bootstrap: ""
     bootstrap_timeout: 10m
   ```

   For a stateful monorepo that uses one shared database, keep Task execution
   sequential so bootstrap runs once on the reused Run Worktree:

   ```yaml
   worktree:
     concurrency: 1
     copy: [".env", "packages/backend/.env"]
     bootstrap: "bun install && bun run db:migrate && bun run db:seed"
     bootstrap_timeout: 10m
   ```

10. Stop an Active Run for a Spec with `roundfix stop --spec <slug>` from inside
   the current repository. This resolves that repository's Spec target and
   records a Stop Request; the Run stops after the current Work Item settles.
   Use `roundfix stop --force --spec <slug>` only for a dead, stuck, or runaway
   Run. It proves the recorded owner exited before cleaning up registered
   active Agent Sessions, and reports Stopped and releases the Active Run lock
   only after that proof. A failed proof leaves the Run Active with its Agent
   Sessions unchanged and its lock retained.


### Token usage

The Implement Run's terminal summary prints a `Tokens:` line after its outcome
on every outcome. It reports recorded tokens and how many prompts reported,
plus adapter-reported cost by currency and Agent Session coverage. Entirely
unreported usage says `none reported`, never `0 tokens`; no usage rows says
`no prompts recorded`. Roundfix sums increases in cumulative Agent Session
cost readings and starts a new count when a reading decreases. It computes
no price from tokens. Prompts outside a Run, including the pre-PR review,
do not count. Use `roundfix runs show <run-id>` for per-scope details.
