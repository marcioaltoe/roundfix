## Driving a Spec implementation loop

Use this loop to carry one Spec — or a queue of Specs — from pending Tasks to an
archived Spec without owning the Run's terminal in the foreground. It composes
the Implement, Attach, Settle, Stop, and Archive commands documented above.

Follow one order per Spec: implement the graph including its authored gate,
run the configured pre-PR review, archive on the branch, pass the repository
gate, open the Pull Request, verify current-head checks, and merge. With
pre-PR review set to `none`, the review is recorded as omitted.

ADR-0091 keeps the authored QA gate before any Pull Request exists, while
ADR-0080 lets environment-blocked rows pass with equivalent evidence. Spec
0078 proved that path: eleven of eighteen rows were blocked, nine of those
eleven on no open Pull Request and each of those nine backed by recorded
payload, command-runner, and event-stream evidence.

1. **Prepare.** Work on a non-default branch and confirm readiness with
   `roundfix doctor`. Pick the Spec slug under the resolved Spec Root
   (`docs/specs/` by default). Do not edit files the Run will touch once it is
   Active; overlapping local edits end the Run Integration Pending.

2. **Start detached.** Launch the Run without owning its lifetime:

   ```bash
   roundfix implement --spec <slug> --detach
   ```

   Capture `<run-id>` from the five-line report. Detach implies non-interactive
   mode. If startup fails before the handshake, the foreground command relays
   the child's stderr and exit code with no stdout report — fix the reported
   Preflight Validation failure and start again.

3. **Monitor without owning.** Use the printed
   `roundfix events <run-id> --follow --filter outcome` command for the stable
   terminal subscription. Open the read-only Live Run View with
   `roundfix attach <run-id>` when a human needs progress detail. From a fresh
   terminal, discover the Run with the bounded `roundfix runs list` or open
   the Run Browser with `roundfix attach`. Attach replays the Run Event
   Journal and follows new events; `q` or `Ctrl-C` detaches and never stops
   the Run. The Console Log is a compact text record, not a state API.

4. **Detect the terminal outcome.** The Run ends with exactly one stdout outcome
   line in the console log:
   - `Clean: all N Task(s) completed.` — every Task passed and integrated onto
     the current branch; the Run Worktree and Run Branch are removed.
   - `All N Task(s) already completed; no Run was created.` — nothing to do;
     advance to the next Spec.
   - `IntegrationPending: … integrate with git merge --ff-only roundfix/run-<id>`
     — Tasks completed but the current branch could not fast-forward. Run the
     printed command from the repository root, then continue. This usually means
     the user checkout drifted while the Run was Active.
   - `Unresolved: X completed, Y failed, Z skipped, W pending.` — one or more
     Tasks did not settle. Go to recovery.

5. **Recover failed Tasks.** Read the per-Task status lines
   (`task_NN failed — <title>`) and the following indented `reason:` line when
   present. For each failed Task, inspect its kept Task Worktree or the kept Run
   Worktree, then recover only that Task once its Verification passes there:

   ```bash
   roundfix settle --spec <slug> --task <task_id>
   ```

   Settle re-runs the Task's Verification in the kept surface, commits on pass,
   and integrates onto the Run Branch. Re-run `roundfix implement --spec <slug>`
   to pick up any still-pending Tasks; completed Tasks are skipped.

6. **Stop when needed.** Prefer graceful `roundfix stop --spec <slug>`; the Run
   ends after the current Work Item settles. If a later command finds an
   Active-Run lock whose recorded owner PID is provably dead, Roundfix reclaims
   that orphan automatically with a stderr warning and proceeds. Use
   `roundfix stop --force --spec <slug>` only when the owner still appears
   live, the Run has no recorded PID, or the Run is otherwise stuck or runaway.
   Never kill Agent or acpx processes directly while a Run is Active — force
   stop reaps sessions and terminal Worktree debris for you.

7. **Advance.** When every Task is completed and the newest QA Report is
   acceptable under the one declared-acceptance eligibility policy — either
   `verdict: pass` with no disallowed blocked rows, or `partial` with only
   declared-unreachable unmet rows fully covered by the Spec's declarations —
   archive it with `roundfix archive <slug>`, then start the loop again on the
   next Spec. An archive-eligible `partial` Report now settles the
   terminal `qa` Task completed, because settlement applies the same
   declared-acceptance policy archive applies. The declared-only case
   records the declarations' satisfying actions under `unproven`.

Failure recovery stays clean when you keep two invariants: never edit
Run-touched files while a Run is Active, and never reap sessions or Worktrees by
hand — let automatic orphan reclamation, `roundfix stop --force`, and the
Implement Command preflight sweep close terminal sessions and Worktrees.

## Autonomous Spec delivery

Use this when the user asks to implement one Spec — or a queue of Specs —
autonomously. It wraps the loop above with the delivery half and the rules that
decide when to keep going and when to stop. Invocation is the authorization for
the Supervisor to carry each Spec to a merged Pull Request without pausing
between Specs, reporting at milestones rather than per Agent event.

**The delivery half is Supervisor authority and nothing else grants it.** An
Agent resolving an assigned Batch or Task never pushes, opens a Pull Request, or
merges, whatever this section says: the Daemon owns Task commits and any
configured push, and the Agent's own instructions prohibit publication. If you
are executing an assigned Work Item, stop at a verified handoff and report.

Per Spec, in order:

1. **Branch.** Sync the default branch, then create the Spec's own branch. Never
   commit to that branch while a Run is Active: Roundfix integrates with
   `merge --ff-only`, so your own commit forces Integration Pending. Parallel
   work goes to its own branch.

2. **Author what is missing.** A Spec needs a PRD, a TechSpec when it has
   architectural surface, and a Task Graph. Author them and commit before
   starting the Run. Task decomposition is a derivation from artifacts that
   already passed their own gates — decide it and report it, do not hold the
   loop for approval.

3. **Implement.** Run detached while any Task is still pending or a corrective
   Task is planned:

   ```bash
   roundfix implement --spec <slug> --detach
   ```

4. **The gate runs itself, once, at the very end.** It is authored into the
   graph at decomposition as the terminal `qa` Task, so no invocation
   requests it: the Daemon executes it when the last Task settles, and an
   all-completed graph whose gate is unsettled starts a Run that is just the
   gate:

   ```bash
   roundfix implement --spec <slug> --detach
   ```

   When a corrective Task is added as a dependency of a QA gate that already
   settled `completed`, reopen the settled gate first with
   `roundfix reopen --spec <slug>` before the next Run. A failed or pending
   gate needs no reopen. When a gate returns several findings, close them
   together and let the gate re-run once. More
   than two corrective Tasks generated by QA findings means the
   decomposition needs re-examining, not a third patch.

5. **Review the candidate.** Run `roundfix review` under the configured
   pre-PR review policy. With review enabled, require the selected provider's
   complete current-candidate review and use `roundfix review dispose` to
   record finding dispositions. With explicit `none`, record that review was
   omitted; this does not waive QA or required checks.

6. **Archive and publish.** When every Task is completed and the newest QA
   Report is acceptable under the one declared-acceptance eligibility policy —
   either `verdict: pass` with no disallowed blocked rows, or `partial` with
   only declared-unreachable unmet rows fully covered by the Spec's
   declarations —
   run `roundfix archive <slug>` on the branch, pass the repository gate,
   then push and open the Pull Request. An archive-eligible `partial` Report
   now settles the terminal `qa` Task completed, because settlement applies the same
   declared-acceptance policy archive applies. The declared-only case
   records the declarations' satisfying actions under `unproven`.

7. **Resolve Pull Request feedback.** This applies only when the repository's
   Review Source gives feedback on the Pull Request. Run
   `roundfix watch --source <review-source> --pr <n> --until-clean`.
   Only a terminal Clean with no unresolved Review Issues
   clears the Pull Request for merge. Reaching the configured round cap is not
   that signal: the cap can be exhausted while findings remain, so treat it as
   a stop and escalate with what is still open.

8. **Merge and continue.** Before merging, confirm the exact commit that will
   land: required checks pass on the current head and, when review is enabled,
   its result covers that same commit rather than an earlier one. Require Clean
   Pull Request feedback when the Review Source applies. Review fixes move the
   head, so a result from before them proves nothing about what merges. Then
   squash merge, delete the branch, sync the default branch,
   `roundfix reconcile`, rebuild any local binary, and start the next Spec.

### When to stop and ask

Stop only when the next decision would change authority, architecture, or blast
radius:

- a change would touch protected tooling the Spec does not authorize with exact
  bounded files;
- a Task's acceptance genuinely needs a person or a credential the repository
  cannot hold, so its Verification cannot be hermetic;
- the work only proceeds by contradicting the TechSpec, so the Spec is wrong
  rather than merely thin;
- one step carries irreversible blast radius — a data migration, a destructive
  deletion, a release, a credential rotation;
- a QA verdict cannot be reached without a supervised judgement, such as rows
  blocked by an environment the gate cannot provision;
- a review round cap is exhausted with findings still unresolved, or the
  current head lacks a passing check or a review result covering it.

Escalate with the blocking fact and a proposed resolution, never with an open
"is this ok?". Everything else — granularity, ordering, recovery from a failed
Task, re-running a gate — is inside the loop's authority.

### Recovering without breaking the loop

- Read every QA finding before fixing any of them. One commit that closes four
  findings costs one gate cycle; four commits cost four.
- Keep supervisor fixes inside the Spec's authorized paths. The gate audits the
  Supervisor exactly as it audits Agents, and a repair that touches an
  unauthorized path fails the gate on its own.
- Prefer fixing the source of a brittle assertion over bumping its expected
  value: an assertion that reads a constant survives the next change; a literal
  buys one green cycle.
- Verify a stranded Run Branch before discarding it. Failed gate cycles leave
  branches that Branch Integrity Preflight then refuses, and the suggested
  fast-forward can be wrong when integrating a superseded artifact would
  overwrite a newer one.
- `Clean` is a status, not evidence. Read the diff after a Run reports Clean,
  and confirm a review actually happened before treating its silence as proof.

