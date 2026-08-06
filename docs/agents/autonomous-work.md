# Autonomous work: Supervisor and implementation runtimes

One Supervisor session orchestrates autonomous work. Each Task owns a Task
Type-selected Agent Session inside its Roundfix Run, and an authored terminal
`qa` Task owns a separate `qa` Agent Session. This split binds interactive and
unattended sessions without forcing a mixed Task Graph through one Run-wide
Agent.

## Roles

| Role | Agent Model and reasoning | Work |
| --- | --- | --- |
| Supervisor | Supervising Claude Code session | Author Specs and their gate decisions, launch and monitor Runs, integrate outcomes, archive, make boundary commits, and route work |
| Non-frontend implementer | Codex `gpt-5.6-sol` with `high` Default Reasoning Effort; `gpt-5.5`/`xhigh` fallback | CLI, backend, infrastructure, documentation, and other non-frontend Tasks and Review Issue Batches |
| Frontend implementer | Claude `opus` (Opus 5, 1M context) with `xhigh` Default Reasoning Effort | Design, UI, UX, Bubble Tea/Lip Gloss TUI, and web frontend Tasks |

## The Supervisor does not implement

The Supervisor:

- routes work through `docs/agents/spec-routing.md`;
- authors planning artifacts;
- launches detached Runs and monitors them to a terminal outcome;
- integrates results, monitors the authored `qa` Task, and archives passing
  Specs;
- maintains the Agent Selection Profiles that select each Task and QA Agent
  Session;
- limits its direct edits to Spec and documentation fixes, Run recovery, and
  boundary commits.

Feature code, tests, and operational implementation belong to an ACP Runtime.

`Fable` is an Agent Model label, not the Supervisor role. Never pass a
Supervisor-backed runtime through `--agent`.

## Authoring Specs

The Supervisor authors Specs using `write-idea`, `write-prd`,
`write-techspec`, and `write-tasks`.

- Directing autonomous work authorizes the whole authoring chain. For every
  Spec already in the approved queue, the Supervisor writes the missing
  `_techspec.md` and Task Graph and starts the Run without asking again. A
  missing TechSpec is authored, never a reason to stop and request
  permission, and a Task Graph the Spec already authorizes is never held for
  approval — granularity, dependency order, and Task Type are derivations.
- Ask only before adding a Spec that is not in the approved queue, or when an
  escalation trigger fires: absent tooling authority with bounded files, an
  acceptance that no hermetic Verification can reach, a Spec whose artifacts
  the decomposition would have to contradict, or an irreversible action inside
  one slice. Technical clarification during authoring is not a permission
  request and stays allowed.
- Record accepted product and technical decisions in ADRs.
- Add canonical terms to `CONTEXT.md` as they are resolved.
- Edit shared files only while no Run is Active.
- Write self-contained Task Verification commands that run in a bare worktree.
  Use real paths, no placeholders, no assumed install state, and Roundfix's
  required `-buildvcs=false` build flags.
- Use `## Context` only for Task-specific instruction and interface paths:
  `- instruction: <path>` and `- interface: <path>`.
- Maintain one dependency-and-risk-ordered queue of approved Specs.

## Profile-led Task routing

```bash
roundfix implement --spec <slug> --detach
```

The Implement Command resolves a profile for each Task Type. Review Runs use
the same profile-led selection for their review Agent Session. When the Task
Graph declares a terminal `qa` Task, that Task resolves the separate `qa`
profile after every dependency settles completed.

1. `.roundfixrc.yml` pins `profiles.general`, `profiles.backend`,
   `profiles.qa`, and `profiles.review` to the Sol/high Preferred Selection.
2. Each Codex-led profile keeps GPT-5.5/xhigh as its ordered Fallback Chain.
3. Roundfix proves every exact tuple before creating a Run.
4. The effective Codex adapter must prove the official
   `@agentclientprotocol/codex-acp` lineage and supported version.

For a one-Run exception, provide `--agent`, `--model`, and
`--reasoning-effort` together. A partial override is rejected.

## Frontend Tasks: Claude Opus 5 xhigh

Keep frontend Tasks in the same Task Graph and give them `type: frontend` when
they are dominated by:

- visual or interaction design;
- UI, UX, navigation, information hierarchy, feedback states, or user-facing
  copy;
- the Bubble Tea/Lip Gloss TUI;
- a future web frontend.

The built-in `frontend` profile selects `claude / opus / xhigh`, and the
`claude` runtime launches through the official
`@agentclientprotocol/claude-agent-acp` adapter. The current configuration
defaults remain authoritative; model recommendations do not change routing.

## Routing rules

- Every Task owns one Agent Session selected from its Task Type profile.
- A mixed Task Graph keeps frontend and non-frontend Tasks together; Task Type
  selects the applicable profile without changing dependency order or Waves.
- `general` handles work without a more specific Task Type profile. Optional
  `data`, `infra`, `docs`, `test`, and `chore` profiles inherit `general` when
  absent.
- `backend` selects backend work, and `frontend` selects frontend/design work
  where interaction and visual quality are part of the outcome.
- The authored terminal `qa` Task runs in its own `qa` Agent Session after all
  of its dependencies complete.
- Review work uses a review-selected Agent Session.
- One-off delegation outside a Run follows the same routing.

## Verification

Runtime selection does not change completion requirements:

- The Implement Agent hands back implementation-ready work after focused
  checks and does not run the Task's declared `## Verification` commands,
  edit Task status, or claim a terminal verdict.
- The Daemon alone writes Implement Task status and runs every Task's
  `## Verification` commands verbatim before settlement.
- A deterministic failure releases Verification Capacity before the same
  Agent Session receives one Verification Feedback repair turn; the final
  Daemon attempt queues for capacity again.
- Exit `75` is a project-authored Temporary Verification Failure signal for
  one exclusive retry. Roundfix never infers it from logs.
- `make verify` gates every completion claim.
- `qa-gate` runs from the authored terminal `qa` Task regardless of runtime;
  a Spec that declines the gate records that decision and its reason during
  decomposition.
- Agent instructions and Spec conventions bind both runtimes.

### Author the QA gate once per Spec

`write-tasks` declares the gate during decomposition. Include one terminal
`qa` Task that depends on every non-QA leaf, or declare `qa: declined` with a
non-empty reason. The Implement Command never makes that decision. An
all-completed implementation graph with an unsettled authored gate remains
runnable because the gate is still a pending Task.

Rebuild any binary the gate exercises before the gate becomes runnable, or the
gate reports defects the running artifact does not contain.

If the Task Graph grows after its gate reports, load-time validation
invalidates that gate and names the inserted Tasks. Batch corrective work:
appending one corrective Task after each finding turns discovery into a serial
chain of full gate cycles. When a gate returns several findings, close them
together before the graph reaches the gate again. If a Spec needs more than
two corrective Tasks generated by QA findings, stop and re-examine the Spec's
decomposition rather than appending a third.

## The loop that implements a Spec

The Supervisor runs Specs end to end without pausing between them. After
branching from the synced default branch, authoring any missing artifacts,
including the gate decision, and committing and pushing them, follow one
order: implement the graph including its authored gate, archive, open the Pull
Request, watch until Clean, and merge. Reconcile the terminal Run afterward.

ADR-0091 keeps the authored QA gate as the graph's terminal Task, so the gate
runs before any Pull Request exists. ADR-0080 lets an environment-blocked row
reach `pass` when the report carries equivalent evidence. Spec 0078 proved
that path on 2026-08-05: its gate reached `pass` with eleven of eighteen rows
blocked, nine of those eleven on no open Pull Request and each of those nine
backed by recorded payload, command-runner, and event-stream evidence.

Invocation is the authorization. Stop only for a decision that is genuinely
the maintainer's: a tooling authorization the Spec does not carry, an
irreversible action outside the Spec's declared scope, or a Spec that turns
out to be wrong rather than merely thin.

### Specify the hazard, not only the objective

The costliest rework comes from requirements whose objective is clear and
whose hazard is implicit. Before closing a requirement, ask where the data
comes from and who reads its destination:

| Written | Delivered | Should have been written |
| --- | --- | --- |
| "key derived from the caller's IP" | IP read from a forgeable header | "IP from the **unforgeable** source address" |
| "per-dependency time cap" | a race that abandons the query | "a cap that **cancels the underlying operation**" |
| "technical detail in the log only" | the credential moved into the log | "no secret in the response **or the log**" |

A risk written down in a Spec and dismissed as hypothetical is the same
defect one step earlier. If a probe would settle it, run the probe while
writing the Spec.

### Verify the class, not the case

When a finding is "nothing here may do X", the Verification must sweep every
place that could, not the one place that did. A redaction applied to two of
four call sites passed its scoped check, passed the repository gate, and
passed review — only a manual read caught it.

Write absence assertions in the portable form, because a bare negative
`grep` and `grep -L` both exit non-zero when the pattern is absent and the
Daemon reads that as failure:

```bash
if grep -rn "<forbidden>" <dirs> --include="*.go" | grep -q .; then exit 1; fi
```

Put focused checks first and close every Verification with the repository
gate: scoped suites prove the Task's effect and are blind to the regression
it caused next door.

### Close the test seam before changing the layer

If a Spec touches a layer the tests cannot reach, the first Task creates the
reach. The warning sign is a branch by client, driver, or environment where
tests always take one path and production takes the other — a guard living on
a line no test executes will only ever be found by the slowest detector.

Put the detector where the risk is: a Task that changes a contract with an
external surface must probe that surface in its Verification, not merely
compile and unit-test.

### `Clean` is a status, not evidence

Read the diff after a Run reports `Clean`. Two minutes of reading has both
confirmed a correct fix and refuted a partial one in the same session.

The same applies to review: a Review Source rate-limit notice can be recorded
as an approval, so a `Clean` watch with zero issues may mean no review
happened. Confirm a review occurred before treating its silence as evidence.

### Never commit to the branch while a Run is active

Roundfix integrates a Run Branch with `merge --ff-only`; a commit of your own
on the target branch forces Integration Pending. Parallel tooling work goes to
its own branch. Before opening a Pull Request, confirm which branch is
checked out — `gh pr create` uses the current one.

## Monitoring Runs

Supervisors use the Run Event Stream:

```bash
roundfix events <run-id> --follow
roundfix events <run-id> --follow --filter task-status,verification,outcome
```

Each stdout line is a `roundfix-events/v1` JSON object. Diagnostics go to
stderr. Verification records expose `waiting`, `started`, command, verdict,
retry, and capacity evidence. Use `roundfix attach <run-id>` only when a human
needs the Live Run View, whose spec rows show `Agent working`,
`Waiting for Verification`, and `Verifying`; the Console Log is not a state
API.

<!-- setup-context-driven:begin id=guide.autonomous-work version=0.0.1 -->

# Autonomous work

Default backend work uses `codex gpt-5.6-sol`. Design, UI, UX, and
frontend-dominant work uses `claude opus 5 xhigh` when the Task Graph routes that
surface.

- **mandatory**: The Supervisor authors Specs, starts and monitors Runs, and orchestrates outcomes. Delegate implementation to the selected ACP Runtime through a Roundfix Run.

- **prohibited**: The Supervisor must not write feature code or tests.

- **mandatory**: Each Task owns a Task Type-selected Agent Session; mixed frontend and non-frontend Tasks remain in one Task Graph, and requested QA owns a separate `qa` Agent Session.

- **mandatory**: The Daemon runs each Task's declared Verification verbatim. A Task can settle `completed` only after that Verification passes; failed diagnostics return to the same Agent Session for the bounded retry policy.

<!-- setup-context-driven:end id=guide.autonomous-work -->

<!-- roundfix:repository-rule:begin id=rule.0d49c2a353155f1dc777e6513b646c37287d41769fd717c0463a98c9567d295e -->
- **HARD RULE — autonomous work model**: binding for every autonomous
  session — the Supervisor orchestrates only; implementation is delegated to an
  ACP Runtime per `docs/agents/autonomous-work.md`.


<!-- roundfix:repository-rule:end id=rule.0d49c2a353155f1dc777e6513b646c37287d41769fd717c0463a98c9567d295e -->

<!-- roundfix:repository-rule:begin id=rule.6b28f6bdd57f8a34a73143025f571d9c12bd6d571192ec9a27ac12cc21e74bce -->
### Autonomous work

Supervisor orchestrates and authors Specs; each Task is delegated to its
Task Type-selected Agent Session. The current Agent Selection Profiles define
the effective runtime, model, reasoning effort, and fallback order for backend,
frontend, QA, review, and other categories. Mixed frontend and non-frontend
Tasks remain in one Task Graph. Binding for every autonomous session. See
`docs/agents/autonomous-work.md`.


<!-- roundfix:repository-rule:end id=rule.6b28f6bdd57f8a34a73143025f571d9c12bd6d571192ec9a27ac12cc21e74bce -->
