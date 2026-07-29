# Autonomous work: Supervisor and implementation runtimes

One Supervisor session orchestrates autonomous work. Each Task owns a Task
Type-selected Agent Session inside its Roundfix Run, and a requested QA phase
owns a separate `qa` Agent Session. This split binds interactive and unattended
sessions without forcing a mixed Task Graph through one Run-wide Agent.

## Roles

| Role | Agent Model and reasoning | Work |
| --- | --- | --- |
| Supervisor | Supervising Claude Code session | Author Specs, launch and monitor Runs, integrate outcomes, run `qa-gate`, archive, make boundary commits, and route work |
| Non-frontend implementer | Codex `gpt-5.6-sol` with `high` Default Reasoning Effort; `gpt-5.5`/`xhigh` fallback | CLI, backend, infrastructure, documentation, and other non-frontend Tasks and Review Issue Batches |
| Frontend implementer | Claude `opus` (Opus 5, 1M context) with `xhigh` Default Reasoning Effort | Design, UI, UX, Bubble Tea/Lip Gloss TUI, and web frontend Tasks |

## The Supervisor does not implement

The Supervisor:

- routes work through `docs/agents/spec-routing.md`;
- authors planning artifacts;
- launches detached Runs and monitors them to a terminal outcome;
- integrates results, runs `qa-gate`, and archives passing Specs;
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

- Ask before creating a new Spec. Invoking an approved Spec authorizes its
  Tasks without another confirmation.
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
roundfix implement --spec <slug> --qa --detach
```

The Implement Command resolves a profile for each Task Type. Review Runs use
the same profile-led selection for their review Agent Session, and `--qa`
resolves the separate `qa` profile after every Task settles completed.

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
- Requested QA runs in its own `qa` Agent Session after all Tasks complete.
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
- `qa-gate` runs after the final Task regardless of runtime.
- Agent instructions and Spec conventions bind both runtimes.

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
