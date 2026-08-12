---
task: task_07
spec: 0042-verification-capacity-and-daemon-task-settlement
status: completed
type: docs
complexity: medium
---

# Task 07: Align operator docs and prepare Agent Skill wording

## Overview

Publish one operator and Agent contract for independent capacities,
Daemon-owned Task status, Waiting for Verification, and the explicit exit-75
protocol. Configuration, command guidance, autonomous-work instructions, the
repository-owned authorial workflow, and canonical/embedded Roundfix Skills
must match shipped behavior and preserve skill ownership boundaries.

## Requirements

1. MUST document `verification.concurrency`, default `1`, strict positive
   validation, precedence, and its independence from `worktree.concurrency`.
2. MUST document the recommended Task Capacity `2` / Verification Capacity `1`
   flow without claiming machine-wide or external-process coordination.
3. MUST document Daemon-only Implement Task status, implementation-ready Agent
   handoff, focused checks, authoritative full Verification, one Agent repair,
   and capacity release around repair.
4. MUST document exit code `75` as a project-authored Temporary Verification
   Failure, its one exclusive retry, retained evidence, exhaustion behavior,
   and the prohibition on log heuristics.
5. MUST document Run Event Stream and Live Run View working, waiting,
   verifying, retry, and capacity evidence with profile-led command examples.
6. MUST leave every protected tooling and upstream-managed Skill path
   untouched; Task 08 owns the isolated authorial Skill update.
7. MUST preserve ADR, finding, PRD, Tech Spec, glossary, and skill ownership
   and synchronization traceability and keep all repository content in
   English.
8. MUST treat the complete Daemon-owned repository verification gate as
   blocking; any format, test, Skill synchronization, or build failure prevents
   settlement.
9. MUST align `docs/agents/autonomous-work.md` and the `CONTEXT.md` Agent
   Session definition with ADR-0051: each Task owns a Task Type-selected Agent
   Session, frontend Tasks remain in the same mixed Task Graph, and current
   Agent Selection Profile defaults remain authoritative.
10. MUST prepare exact canonical wording for Task 08 without editing its four
    protected Skill targets.

## Subtasks

- [x] Update configuration, command, usage, and autonomous-work guidance.
- [x] Verify implementation-ready handoff and Daemon Verification instructions.
- [x] Document exit-75 project ownership, exclusive retry, and exhaustion.
- [x] Document event and Live Run View capacity/phase evidence.
- [x] Align autonomous-work routing and the Agent Session glossary with
      ADR-0051.
- [x] Prepare the authorial Skill wording and hand it off to Task 08.
- [x] Resolve links, terminology, and examples.

## Acceptance Criteria

- [x] Users can configure both capacities and predict default, precedence,
      validation, and per-Run scope from one canonical page.
- [x] No supported Agent instruction tells an Implement Agent to run declared
      Task Verification, edit Task status, or settle its own terminal verdict.
- [x] Full-gate requirements remain explicit: the Agent hands off after focused
      checks and the Daemon runs the mandatory Task Verification before
      settlement.
- [x] Exit `75` documentation clearly assigns classification to the project
      wrapper and promises only one observable exclusive retry.
- [x] Events/Attach examples use canonical phase names and keep requested
      command output separate from diagnostics.
- [x] Autonomous-work guidance and the Agent Session glossary route every Task
      through its Task Type-selected Agent Session, keep frontend Task 05 in
      this graph, and contain no superseded one-Agent-per-Run rule.
- [x] No protected tooling path changes in this Task.
- [x] No upstream-managed skill changes and every ADR/finding/Spec link resolves.
- [x] `make verify` passes completely after all documentation and generated
      skill changes.

## Context

- instruction: `docs/agents/skill-dispatch.md`
- instruction: `docs/agents/autonomous-work.md`
- instruction: `.agents/skills/tech-writer/SKILL.md`
- instruction: `.agents/skills/agentic-cli-design/SKILL.md`
- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `CONTEXT.md`
- interface: `README.md`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/user-guide/commands.md`
- interface: `docs/user-guide/usage.md`

## Verification

- `rtk go test ./internal/cli -run 'Test(CommandUsage|DocumentationContract)' -count=1` — expected: user guidance and command examples match implemented behavior.
- `rtk git diff --check` — expected: Spec, ADR, finding, documentation, and generated skill files contain no whitespace errors.
- `rtk make verify` — expected: formatting, Go tests, setup-context checks,
  current skill synchronization, shipped skill validation, and build all pass
  before the protected wording handoff.

## References

- `_prd.md` → all Goals; User Stories 2–7; Core Features 1 and 3–10; User Experience; Decisions.
- `_techspec.md` → System Architecture; API Contracts; Integration Points; Risks & Considerations; Build Order 7.
- `../../adr/0056-spec-runs-separate-task-and-verification-capacity.md` → canonical capacity and temporary-failure contract.
- `../../adr/0057-daemon-exclusively-owns-implement-task-status.md` → canonical Task-status and Agent-handoff contract.

## Task 08 authorial wording handoff

Task 08 owns the four protected targets. It must update each canonical file
under `.agents/skills/` and synchronize its matching embedded copy under
`skills/`; Task 07 does not edit any of them.

Use this exact execution-mode wording in the `implement-task` Skill:

```markdown
## Execution modes

Standalone execution and a Roundfix Daemon-assigned turn share the same Task
slice, Result, and evidence requirements, but they have different settlement
owners.

- In standalone execution, the Agent updates Task status, runs every command
  in `## Verification`, settles the Task from fresh evidence, and commits only
  when the user has authorized a commit.
- In a Roundfix Daemon-assigned turn, the Daemon is the sole Task-status
  writer. The Agent must not edit status, run the declared `## Verification`
  commands, claim a terminal verdict, or commit. It may run focused checks,
  records implementation and focused-check evidence in `## Result`, and hands
  back implementation-ready work. The Daemon then runs the complete declared
  Verification verbatim before settlement.
- A Daemon Verification command failure releases Verification Capacity before
  one Verification Feedback repair turn in the same Agent Session. The Agent
  repairs the implementation, updates `## Result`, runs focused checks when
  useful, and hands back again without running declared Verification or
  settling status.

For a Daemon-assigned Task that arrives `in_progress`, start the Task fresh.
The work-target lock proves that no live Agent owns it. Do not ask how to
resume, and do not normalize or settle its status.
```

Replace Daemon-assigned status, Verification, record, and commit instructions
with:

```markdown
### Daemon-assigned handoff

1. Read the assigned Task, PRD, TechSpec, `CONTEXT.md`, active ADRs, and bounded
   context paths before editing.
2. Implement only the assigned slice. Never edit `_tasks.md`, another Task
   file, or a path outside an authorized tooling allowlist.
3. Run focused implementation checks while working when useful. Do not run any
   command from the Task's `## Verification` section.
4. Append or update `## Result` with the implementation and focused-check
   evidence for every acceptance criterion. Do not use Daemon Verification
   evidence that has not run yet.
5. Hand back implementation-ready work without editing Task status or claiming
   `completed`, `failed`, passing Verification, commit readiness, or delivery.
6. Never commit, push, or open a pull request. The Daemon runs declared
   Verification, writes terminal Task status, and owns the Task commit.
```

Use this exact capacity and retry wording in the Roundfix Skill's user-facing
Spec Run section:

````markdown
Task Capacity and Verification Capacity are independent config-only limits for
one Implement Run. `worktree.concurrency` is Task Capacity, defaults to `2`,
and limits concurrent Task Worktree lifecycles. `verification.concurrency` is
Verification Capacity, defaults to `1`, and limits concurrent Task
Verification attempts. Both must be positive integers; Project Config
overrides User Config, which overrides built-ins. Neither setting coordinates
other Runs, CI, or external processes.

The recommended split overlaps implementation while serializing the complete
repository gate:

```yaml
worktree:
  concurrency: 2

verification:
  concurrency: 1
```

Every normal attempt journals `waiting` before it acquires shared capacity,
then `started`, command results, and `verdict`. A deterministic failure
releases capacity before one Verification Feedback repair turn and reacquires
capacity for the final Daemon attempt.

Exit `75` from a project-authored Verification wrapper is the sole Temporary
Verification Failure signal. Roundfix retains the diagnostics and grants the
Task one exclusive retry across its Verification lifecycle. The retry waits
for other Verification attempts in that Run to drain, consumes the entire
Verification Capacity, and does not consume the Agent repair. A repeated exit
`75` exhausts the retry and fails the Task. Roundfix never classifies a failure
from logs, timing, ports, package names, or framework text.
````

Use this exact ownership wording in the Roundfix Skill's assigned-Task section:

```markdown
## Assigned Task Batches

Inside an Implement Run, each Task owns a Task Type-selected Agent Session and
the Daemon is the sole writer of `in_progress`, `completed`, and `failed`.
Agent-authored status is never a verdict. Frontend and non-frontend Tasks stay
in the same mixed Task Graph; their Task Types select the applicable Agent
Selection Profiles.

The Agent:

1. Reads the assigned Task and bounded context completely.
2. Implements only that Task's slice.
3. Runs focused checks while working when useful, but never runs commands from
   the Task's `## Verification` section.
4. Appends or updates `## Result` with implementation and focused-check
   evidence for every acceptance criterion.
5. Hands back implementation-ready work without editing Task status, claiming
   a terminal verdict, committing, pushing, opening a pull request, editing
   `_tasks.md`, or editing another Task file.

The Daemon writes `in_progress` before Agent work, normalizes any Agent-authored
status after handoff, runs the complete Task Verification verbatim, and alone
settles status and creates the Task commit. A deterministic first failure
releases Verification Capacity before one Verification Feedback repair turn
in the same Agent Session. Exit `75` uses the one exclusive retry protocol and
does not create Agent feedback. Any declared formatter, test, Skill
synchronization, or build failure blocks settlement.
```

Use this exact assigned-Task completion report:

```markdown
For an assigned Task Batch, report:

- The assigned Task id.
- The implementation-ready behavior handed back.
- Focused checks run and their outcomes.
- Files changed in the assigned working tree.
- The `## Result` evidence recorded in the Task file.
- Any blocker that prevented an implementation-ready handoff.

Do not report a terminal Task status, declared Verification result, commit, or
delivery claim; those are Daemon-owned and occur after the Agent turn.
```

## Result

Aligned the operator and repository Agent contract around independent
per-Implement-Run capacities, Daemon-only status settlement, focused-check
handoff, complete Daemon Verification, one bounded Agent repair, and the
project-authored exit-75 protocol. The configuration guide is the canonical
source for defaults, precedence, strict positive validation, scope, and the
recommended Task Capacity `2` / Verification Capacity `1` split. Command and
usage guidance now exposes the event and Live Run View evidence without mixing
requested stdout with diagnostics.

ADR-0051 now governs the autonomous workflow and glossary: every Task owns a
Task Type-selected Agent Session, requested QA owns a separate `qa` Session,
and frontend Tasks remain in the mixed Task Graph. The authorial wording above
is the exact handoff for Task 08; this Task did not mutate its protected Skill
targets.

### Verification

- `rtk go test ./internal/cli -run 'Test(CommandUsage|DocumentationContract)' -count=1`
  passed: 4 tests in `internal/cli`.
- `rtk git diff --check` passed with no whitespace errors.
- The first sandboxed `rtk make verify` reached 2,722 passing tests but could
  not fork `/bin/ps` in five process-identity tests. Re-running the exact gate
  with local process access passed: 2,727 Go tests, 4 Skill policy tests,
  Roundfix Skill synchronization/validation, and the build.
- Read-only path inspection resolved ADR-0051, ADR-0056, ADR-0057, the PRD,
  the TechSpec, and the configuration target used by the new links.
- `rtk git -c core.fsmonitor=false diff --name-only` listed only operator docs,
  `CONTEXT.md`, and this Task file. None of the four protected Task 08 targets,
  any upstream-managed Skill, `_tasks.md`, another Task file, an ADR, or a
  finding changed.

### Acceptance evidence

1. The configuration capacity section defines both defaults, the built-in →
   User Config → Project Config precedence, strict positive-integer
   validation, absence of capacity flags, independent semantics, and per-Run
   boundary.
2. Autonomous guidance directs Implement Agents to focused checks and
   implementation-ready handoff while reserving declared Verification,
   status, settlement, and commits for the Daemon. The protected Skill wording
   needed to complete repository-wide synchronization is prepared verbatim for
   Task 08 above.
3. Command, usage, autonomous, README, and glossary wording preserves the
   complete declared gate as the blocking settlement authority; formatter,
   test, Skill synchronization, or build failure cannot be treated as success.
4. Command and usage guidance assigns exit `75` only to a project-authored
   wrapper, retains separate diagnostics, permits one exclusive retry, records
   exhaustion, and prohibits log heuristics.
5. Events examples use `task-status,verification,outcome`, canonical
   `waiting`, `started`, `command-passed`, `failed`, and `verdict` phases,
   shared/exclusive retry evidence, and separate stdout/stderr redirection.
   Attach guidance names both effective capacities and the `Agent working`,
   `Waiting for Verification`, and `Verifying` labels.
6. Autonomous routing and the Agent Session definition match ADR-0051 and
   contain no one-Agent-per-Run or separate-frontend-Spec rule.
7. Changed-path inspection proves the protected tooling boundary remained
   intact.
8. Existing ADR, finding, PRD, TechSpec, glossary, and Skill ownership
   traceability remains in English; the referenced local targets resolve.
9. The complete repository gate passed after the operator documentation
   changes, including current Skill synchronization and shipped Skill
   validation.

### Follow-up

- Task 08 must apply the exact authorial wording above to its four expressly
  authorized Skill targets and run its own post-change synchronization and
  complete repository gate.
