---
spec: 0319-the-spec-check-probe-leaves-completed-tasks-alone
prd: _prd.md
created: 2026-10-26
---

# The Spec Check probe leaves completed Tasks alone — Technical Spec

## Executive Summary

`roundfix spec check --run-verification` and the Daemon's pre-work probe
share one prober (ADR-0148) but not one selection: the Daemon probes only the
Tasks it schedules, every Task whose status is not `completed`, while the
check probes every Task in the graph. One exported Daemon function,
`TaskAwaitsWork`, becomes the rule for the Daemon's scheduler, its QA gate
routing and the check. The check reports each command it leaves out as
`not probed` with the reason `completed`, in Task Graph order in text, in a
new optional `verification.notProbed` list in JSON, and as a count on the
clean verdict line; the four verdicts, the exit rule for probed commands and
the uncommitted-source report keep their bytes.

The trade-off accepted: the check no longer tells anyone that a completed
Task's command would fail in a fresh checkout. That question belongs to the
Task's settlement, which already proved the command, and to the repository
Verification; it is not one the pre-work probe can answer, because the work
exists.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local probe only. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0148 and ADR-0315, as
  `_prd.md` → Project Constraints states; ADR-0148 gains a dated amendment in
  task_02. Source: `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — protected paths limited to the Roundfix and
  `qa-gate` skills. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0319-the-spec-check-probe-leaves-completed-tasks-alone/_authorization.md`;
  bounded files: `.agents/skills/qa-gate/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/spec.md`, `skills/qa-gate/SKILL.md`,
  `skills/roundfix/SKILL.md`.

## Current behavior

Read on 2026-10-26 at `86cb6b73` (v0.123.0).

1. `internal/cli/spec_check.go` `probeSpecVerifications` loads the
   working-tree graph, collects `_tasks.md` and every Task file as
   Verification sources for `uncommittedVerificationSources`, then calls
   `daemon.ProbeCommands` for every `graph.Tasks` entry in a disposable
   checkout of `HEAD`, whatever its status. Each verdict becomes a
   `specCheckVerificationCommandReport` (`task`, `command`, `verdict`,
   `cause`) in `verification.commands`.
2. `runSpecCheckCommand` exits `1` when any reported verdict is `vacuous`,
   `unknown` or `malformed`. `renderSpecCheckText` counts `unknown` and
   `malformed` as could-not-be-executed and every other verdict as executed,
   passes both to `speccheck.RenderText` through `speccheck.VerificationCoverage`
   (`Ran`, `Executed`, `Unexecuted`), prints `Verification tree: HEAD`, the
   uncommitted-source lines, then one `- <task>: <verdict> — "<command>"` line
   per command, or `No authored Verification commands.` when there is none.
   `renderSpecCheckJSON` writes `schema: roundfix-speccheck/v1` with
   `verification` holding `executed`, `tree`, `commands` and `uncommitted`.
3. `internal/speccheck/report.go` `RenderText` prints the clean verdict line
   `No findings. Authored Verification commands executed: <n>.` or
   `…executed: <n>; could not be executed: <u>.`;
   `TestVerdictLineStatesProbeCoverage` pins both.
4. The Daemon (`internal/daemon/task_engine.go`): `initialTaskRunStatuses`
   maps a `completed` Task to `taskRunCompleted` and every other status
   (`pending`, `in_progress`, `failed`) to `taskRunPending`; `nextReadyTask`
   starts only pending Tasks; `executeTaskWorker` runs `verifyTaskPreWork`,
   the only pre-work caller of `probeCommandsForRun`. The QA gate is taken out
   of the scheduler (`taskPlanWithoutQAGate`) and routed only while
   `qaTask.Status != spec.StatusCompleted` (the two conditions after
   `runTaskScheduler` in `TaskCycle`). `TestTaskCycleRerunsStaleTasksAndSkipsCompletedTasks`
   pins that a completed Task never runs while in-progress and failed ones
   run again.
5. Reproduced in a temporary repository (see `_prd.md` → Validation): a
   completed Task's passing command is reported `vacuous` and a completed
   Task's missing tool `unknown` (exit 127), and the check exits `1`.
6. Docs: `docs/user-guide/commands/spec.md` and
   `.agents/skills/roundfix/references/spec.md` describe `--run-verification`
   without a selection rule. `.agents/skills/qa-gate/SKILL.md`, in its
   precondition bullet under `## Surface Transcripts at the gate`, says
   "Probing completed Tasks would therefore report their already-passing
   commands as vacuous and refuse the gate", which this Spec makes false.
   `CONTEXT.md` **Vacuous Verification** describes the Daemon's refusal only.
   ADR-0148 shares the classification between both callers, and its text
   names no rule for which Tasks each one probes;
   ADR-0148: "lives in one extracted prober that both callers use".

## System Architecture

- `internal/daemon/verification_probe.go`: `TaskAwaitsWork`, next to the
  exported `ProbeCommands` it selects for.
- `internal/daemon/task_engine.go`: `initialTaskRunStatuses` and the QA gate
  routing read `TaskAwaitsWork`.
- `internal/cli/spec_check.go`: selection, the `notProbed` list and the
  ordered text lines.
- `internal/speccheck/report.go`: `VerificationCoverage.NotProbed` and the
  verdict-line suffix.
- Tests: `internal/cli/spec_check_completed_tasks_test.go` (new),
  `internal/daemon/task_awaits_work_test.go` (new),
  `internal/speccheck/report_test.go`.
- Docs (task_02): the `spec` guide, the Roundfix `spec` reference and its
  mirror, the `qa-gate` skill and its mirror, both skills' version pairs,
  `skills/testdata/owned-skill-versions.json`, ADR-0148, `CONTEXT.md`,
  `docs/references/behavior-surfaces.json`.

## Implementation Design

### Interfaces

In `internal/daemon/verification_probe.go`:

```go
// TaskAwaitsWork reports whether a Run still has work to do on task: every
// Task whose status is not completed. The Daemon schedules, probes and gates
// exactly these Tasks, and the authoring check probes the same set.
func TaskAwaitsWork(task spec.Task) bool
```

`initialTaskRunStatuses` maps a Task to `taskRunCompleted` when
`!TaskAwaitsWork(task)` and to `taskRunPending` otherwise; both QA routing
conditions read `TaskAwaitsWork(*qaTask)` in place of
`qaTask.Status != spec.StatusCompleted`. No other Daemon status comparison
changes.

Text, one line per authored command in Task Graph order, probed and
not-probed interleaved:

```text
- <task>: not probed — "<command>" (Task completed)
```

The clean verdict line, with `n` not-probed commands (`n > 0`); the
could-not-be-executed part appears only when `u > 0`, as today:

```text
No findings. Authored Verification commands executed: <e>; could not be executed: <u>; not probed (completed Tasks): <n>.
```

`No authored Verification commands.` is printed only when the graph has no
command at all, probed or not.

The guide paragraph, appended to `docs/user-guide/commands/spec.md` and the
canonical `.agents/skills/roundfix/references/spec.md` directly after the
paragraph ending "If the parser cannot start, the existing command probe still
runs.":

```markdown
The probe runs only the commands of Tasks a Run would still execute: every
Task whose status is not `completed`, the rule the Daemon uses to schedule
Tasks and route its QA gate. A completed Task's commands are not run, because
they pass once its work exists and may need preparation the disposable
checkout lacks. Each is reported as
`- <task>: not probed — "<command>" (Task completed)`, the clean verdict line
adds `; not probed (completed Tasks): <n>`, and JSON lists them under
`verification.notProbed` as `{task, command, reason}` objects with reason
`completed`, a key present only when the list is not empty. A not-probed
command never changes the exit code. A vacuous, unknown or malformed command
of a pending, in-progress or failed Task, or of a pending QA gate, still makes
the check exit `1`, and the uncommitted-source lines still name every Task
file, completed ones included.
```

The `qa-gate` precondition bullet ("A clean authoring check is a precondition
of the gate"), under `## Surface Transcripts at the gate` and outside
`### QA settlement`: its paragraph after the `roundfix spec check <slug>
--strict` block is replaced, keeping the first and last sentences:

```markdown
  The gate keeps this non-probing form because the probe asks whether a
  command already passes before its work exists, while the gate runs after
  every implementation Task is complete. The probe leaves completed Tasks
  unprobed and reports their commands `not probed`, so it has nothing left to
  ask of them. The three authoring skills use `--run-verification` while the
  work is still being authored, when that question is answerable.
```

The `CONTEXT.md` **Vacuous Verification** definition gains, after its second
sentence:

```markdown
A Task whose status is already `completed` is never probed, by the Daemon or by `roundfix spec check --run-verification`: its command passing is the postcondition its work proved, and the check reports it `not probed`.
```

ADR-0148 gains, after its last paragraph, with `updated_at` set to
`2026-10-26T00:00:00Z`:

```markdown
Amendment (2026-10-26): both callers probe the same Tasks. The Daemon's
`TaskAwaitsWork`, true for every status but `completed`, chooses the Tasks the
Daemon schedules and the Tasks `roundfix spec check --run-verification`
probes. A completed Task's command passes because its work exists, so the
authoring check reports it `not probed` with its reason and never fails on it;
the vacuous, honest, unknown and malformed classification is unchanged.
```

### Data Models

`specCheckVerificationReport` gains, after `Uncommitted`:

```go
NotProbed []specCheckVerificationNotProbed `json:"notProbed,omitempty"`
```

with

```go
type specCheckVerificationNotProbed struct {
	Task    string `json:"task"`
	Command string `json:"command"`
	Reason  string `json:"reason"` // "completed"
}
```

and an unexported, unmarshalled ordered line list the text renderer reads.
`speccheck.VerificationCoverage` gains `NotProbed int`. The JSON schema
string stays `roundfix-speccheck/v1`: the change adds an optional key and no
`verdict` value.

### API Contracts

1. API Contract: `probeSpecVerifications` probes a Task only when
   `daemon.TaskAwaitsWork(task)`; for any other Task it runs no command and
   appends one `notProbed` entry and one text line per authored command, in
   graph order.
2. API Contract: the exit code reads only `verification.commands`, as today;
   a `notProbed` entry never sets it.
3. API Contract: `uncommittedVerificationSources` still receives `_tasks.md`
   and every Task file, completed ones included.
4. API Contract: with no completed Task, text and JSON bytes are unchanged,
   with no `notProbed` key.
5. API Contract: the Daemon's scheduler and QA routing read `TaskAwaitsWork`;
   which Tasks the Daemon runs, probes and gates is unchanged.

### Surface Transcripts

None. The new lines are fixed in Interfaces and asserted by named tests.

### Invariants

```text
1. spec check --run-verification never runs a command of a Task whose status is completed.
2. Every command it does not run is named in text and JSON with its Task and the reason completed, and never changes the exit code.
3. A vacuous, unknown or malformed command of any Task whose status is not completed still makes the check exit 1.
4. The Daemon and the check select Tasks through the same exported function.
5. Without a completed Task, the report keeps its bytes; the uncommitted-source report names every Task file.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, Core Feature 1, Success Metric 1 → API Contract 1, Invariant 1, task_01
- Goal 2, Core Feature 1, Success Metric 5 → API Contract 5, Invariant 4, task_01
- Goal 3, Core Feature 2, Success Metrics 1 and 3 → API Contracts 1, 2 and 4, Invariant 2, task_01
- Goal 4, Core Feature 1, Success Metrics 2 and 4 → API Contracts 2 and 3, Invariants 3 and 5, task_01
- Goal 5, Core Feature 3, Success Metric 6 → Interfaces (guide, references, skills, ADR, glossary), task_02
- Success Metric 6 → Invariant 5, task_01 and task_02

## Absence Evidence

1. A shared Task selection function — command: `grep -rn "func TaskAwaitsWork\|TaskAwaitsWork(" internal`; result: no match at `86cb6b73`.
2. A not-probed report — command: `grep -rn -i "notProbed\|not probed" internal/cli internal/speccheck docs/user-guide .agents/skills`; result: no match at `86cb6b73`.

## ADR Alternatives

- ADR-0148: extends it; the rejected alternative of two implementations that
  can disagree now also covers the selection, so this Spec rejects a second
  status comparison in `internal/cli` for the same reason, and amends the ADR
  instead of minting a new one.
- ADR-0315: unchanged; its decision keeps the four-verdict classification,
  and `not probed` is reported beside the verdicts, not as a fifth one.

## Integration Points

None. No network, provider, forge or other repository is reached; the probe
runs in a disposable checkout of the repository under test, as today.

## Testing Approach

Every test builds temporary repositories; none reads this repository's
Specs.

- `internal/cli/spec_check_completed_tasks_test.go` (new), with a
  `Suite:`/`Invariant:`/`Boundary` header, helper
  `newSpecCheckStatusWorkspace(t, seeds)` on the `clean` fixture of
  `newSpecCheckWorkspace` (a `qa` seed gets `qa:` in the manifest and
  `spec.DerivedQAVerification("clean")` as its commands), and
  `TestSpecCheckLeavesCompletedTasksAlone` (calls `t.Parallel()` first) with
  subtests:
  - `a resumed Spec probes only the Tasks still to run`: Task 01 `completed`
    with `touch "<marker>"` and `roundfix-verification-binary-that-does-not-exist`,
    Tasks 02–04 `pending`, `in_progress`, `failed` with failing `test -f`
    commands. Text: exit `0`, verdict line
    `No findings. Authored Verification commands executed: 3; not probed (completed Tasks): 2.`,
    and the exact five command lines in graph order after
    `Verification tree: HEAD`. JSON: exit `0`, `notProbed` equal to the two
    entries with reason `completed`, `commands` honest for Tasks 02–04 only.
    The marker file does not exist.
  - `a pending vacuous command still fails`, `a pending unknown command still
    fails`, `a pending malformed command still fails`: Task 01 `completed`
    (`true`), Task 02 `pending` with `true`, the missing tool or `if`: exit
    `1`, one command entry with that verdict, one `notProbed` entry.
  - `a pending QA gate is still probed`: Task 01 `completed`, a pending `qa`
    Task: exit `0`, its derived command `honest` in `commands`.
  - `every Task completed`: a completed Task and a completed `qa` Task: exit
    `0`, verdict line `…executed: 0; not probed (completed Tasks): 2.` and no
    `No authored Verification commands.` line.
  - `an uncommitted completed Task source is still named`: the completed Task
    file edited in the working tree is named on
    `Uncommitted Verification source: docs/specs/clean/task_01.md (modified)`.
  - `without a completed Task the JSON has no notProbed key`.
- `internal/daemon/task_awaits_work_test.go` (new):
  `TestTaskAwaitsWorkSelectsEveryStatusButCompleted`, every canonical status
  for a `backend` and a `qa` Task.
- `internal/speccheck/report_test.go`: `TestVerdictLineStatesProbeCoverage`
  gains `probe left completed Tasks unprobed` and `probe left completed Tasks
  unprobed and could not execute a command`; its three existing cases keep
  their expected lines (declared change: two added cases).
- Unchanged and run: `TestSpecCheckRunVerification`,
  `TestSpecCheckReportsAMalformedCommand`,
  `TestSpecCheckNamesAnUncommittedVerificationSource`,
  `TestSpecCheckNamesNoSourceWhenCommitted`,
  `TestTaskCycleRerunsStaleTasksAndSkipsCompletedTasks`,
  `TestTaskCycleQAStepSkippedUnlessEveryTaskCompleted`,
  `TestTaskCycleQAOnlyRunWhenEveryTaskAlreadyCompleted`,
  `TestTaskCycleQAStepRequiresEveryGraphTaskCompleted`.
- A prototype in a scratch clone passed all of them, and each mutation failed
  the named case, whose input makes it observable:
  - `TaskAwaitsWork` returns `true`: `a resumed Spec probes only the Tasks
    still to run` (the marker command runs and is `vacuous`),
    `TestTaskAwaitsWorkSelectsEveryStatusButCompleted` and
    `TestTaskCycleRerunsStaleTasksAndSkipsCompletedTasks` (three Agent runs
    instead of two), so both callers read it;
  - `TaskAwaitsWork` returns true only for `pending`: `a resumed Spec probes
    only the Tasks still to run` (Tasks 03 and 04 are `in_progress` and
    `failed`) and `TestTaskCycleRerunsStaleTasksAndSkipsCompletedTasks`;
  - the check uses its own rule that also skips `qa` Tasks: `a pending QA
    gate is still probed`;
  - the check drops completed Tasks without reporting them: `a resumed Spec
    probes only the Tasks still to run` and `every Task completed`;
  - the verdict line omits the not-probed count: `a resumed Spec probes only
    the Tasks still to run` and `every Task completed`;
  - completed Task files are dropped from the uncommitted sources: `an
    uncommitted completed Task source is still named`.

## Build Order

1. The guide paragraph, the Roundfix reference, the `qa-gate` paragraph, the
   ADR-0148 amendment, the glossary revision, both skills' versions and the
   Behavior Surface record re-recorded last (task_02; no dependency: its
   Verification reads only documents).
2. `TaskAwaitsWork`, the Daemon's use of it, the check's selection and
   report, the verdict line and their tests (task_01; depends on 1, which
   writes the guide of the CLI surface it changes; they share no file).
3. QA gate (task_03) (depends on: 1, 2).

No Spec is listed in `requires`. Spec 0320 delivers after this one and
nothing here reads its code.

## Risks & Considerations

- A JSON consumer that treated `verification.commands` as every authored
  command now misses completed Tasks' commands unless it reads `notProbed`.
  No Go package parses this output; the guide states the split.
- A working-tree edit that sets a Task's status to `completed` takes it out of
  the probe. The edit is not hidden: the uncommitted-source line names the
  file, and the Daemon reads the committed status in its Run Worktree.
- `TaskAwaitsWork` now gates the Daemon's scheduling; a wrong rule would
  re-run or skip work. The existing Task-cycle tests and the new unit test pin
  every status.

## Glossary

- changes: **Vacuous Verification**

## Decisions

- Selection: not `completed` → probed, for every Task type, through one
  exported Daemon function used by both callers.
- Reporting: an optional `notProbed` list and a `not probed` text line, not a
  fifth `verdict` value; the exit rule is unchanged.
- No audit flag; ADR-0148 amended, no new ADR.
