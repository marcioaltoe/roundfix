---
spec: 0319-the-spec-check-probe-leaves-completed-tasks-alone
status: active
created: 2026-10-26
surfaces: [backend, cli]
---

# The Spec Check probe leaves completed Tasks alone

Minimal PRD for a bug fix, minted by `write-techspec`. `roundfix spec check
--run-verification` probes the Verification of every Task in the graph, so on
a resumed Spec a completed Task's command passes because its work exists and
is called `vacuous`, or cannot run without the preparation its work had and is
called `unknown`, and the check exits `1` before a Run that would never
execute that Task again. After this Spec the check probes exactly the Tasks
the Daemon would still execute, chosen by one Daemon function both callers
share, names every completed Task's command as `not probed` with its reason in
text and JSON, and keeps failing on a vacuous, unknown or malformed command of
any Task still to run.

## Sources

One Backlog Entry, written and adopted at birth by this Spec's planning change
(see `references/_index.md`):

1. [`references/2026-10-25-the-spec-check-probe-re-runs-completed-tasks.md`](references/2026-10-25-the-spec-check-probe-re-runs-completed-tasks.md),
   `type: fix`, from Vortex through the Secondbrain project (ai-memory message
   `01a17436-c629-7622-a37d-12528b1da9f2`, 2026-10-25 14:58 BRT): resuming
   Vortex Spec 0162 at `gesttione-solutions/vortex@a961a29` with Roundfix
   0.119.0, `spec check --strict --run-verification` probed completed Task 01,
   whose Bun step exited 127 without `node_modules`, verdict `unknown`, exit
   `1`, while the three pending commands were `honest`.

No open Backlog Entry or unresolved Finding shares this context: at the
authoring base `docs/backlog/` and `docs/findings/` hold no entry.

## Prerequisites

None in `requires`. Authored against `86cb6b73` (v0.123.0) with no active
Spec under `docs/specs/`. Spec 0320 is authored in the same cycle and delivers
after this one; nothing here reads its code.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the `--run-verification` JSON gains one optional list. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; the probe runs local commands in a
  disposable checkout as today. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0148 (one Verification prober
  shared by the Daemon's pre-work probe and the authoring check) gains a
  dated amendment: both callers probe the same Tasks, chosen by one Daemon
  function. ADR-0315 keeps the prober's four-verdict classification
  unchanged, and `not probed` is a selection outcome, not a fifth verdict;
  ADR-0315: "stays the one shared classification of the Daemon's pre-work probe"
  Source:
  `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — protected paths limited to the Roundfix and
  `qa-gate` skills. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0319-the-spec-check-probe-leaves-completed-tasks-alone/_authorization.md`;
  bounded files: `.agents/skills/qa-gate/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/spec.md`, `skills/qa-gate/SKILL.md`,
  `skills/roundfix/SKILL.md`.

## Goals

- `spec check --run-verification` never probes a Task whose status is
  `completed`, so a resumed Spec's delivered work cannot turn the check red.
- The Tasks it probes are exactly the Tasks the Daemon schedules: one exported
  Daemon function decides for the Daemon's scheduler, its QA routing and the
  check.
- Every unprobed command is named, with the Task and the reason, in text and
  JSON, and the verdict line counts it apart from executed commands; it never
  changes the exit code.
- A vacuous, unknown or malformed command of a pending, in-progress or failed
  Task, or of a pending QA gate, still fails the check exactly as today; the
  uncommitted-source report keeps naming every Task file.
- The guide, the Roundfix and `qa-gate` skills, ADR-0148 and `CONTEXT.md`
  describe the rule.

## Core Features

1. One selection rule. `internal/daemon` exports `TaskAwaitsWork(task)`, true
   for every status but `completed`. The Daemon's initial scheduling status and
   its QA gate routing read it instead of comparing the status inline, and
   `probeSpecVerifications` probes a Task only when it returns true.
2. The not-probed report. Each command of a Task the rule leaves out is
   reported in Task Graph order as
   `- <task>: not probed — "<command>" (Task completed)` in text and under
   `verification.notProbed` (`task`, `command`, `reason: "completed"`) in
   JSON, a key present only when the list is not empty. The clean verdict line
   adds `; not probed (completed Tasks): <n>` when `n` is not zero.
3. Documentation. `docs/user-guide/commands/spec.md`, the Roundfix skill
   reference `references/spec.md`, the `qa-gate` precondition paragraph whose
   reason for the non-probing form becomes false, a dated amendment to
   ADR-0148, and the revised **Vacuous Verification** entry in `CONTEXT.md`.

## User Experience

The Vortex maintainer resumes Spec 0162 with Task 01 completed and runs
`roundfix spec check 0162-… --strict --run-verification`. The report reads
`- task_01: not probed — "…" (Task completed)` for each of Task 01's commands,
`honest` for the three pending ones, and a clean verdict line ending
`not probed (completed Tasks): 2.`, and the check exits `0`. A JSON consumer
reads the same commands under `verification.notProbed` and finds
`verification.commands` holding only what ran.

## Non-Goals / Out of Scope

- An audit flag that probes completed Tasks too. No caller depends on
  all-Task probing (`probeSpecVerifications` is the only caller of the probe
  outside the Daemon, measured at `86cb6b73`), and a completed command that
  passes is its proven postcondition, not a defect to report.
- Changing the prober's classification, its parser step, the disposable
  checkout at `HEAD`, the Daemon's pre-work probe, or which Tasks the Daemon
  runs.
- Changing the QA gate's precondition: it keeps the non-probing
  `spec check --strict`.
- Other readers that compare a Task's status with `completed` for their own
  rule (the Spec Consistency detectors, precondition repair, archive); they
  are not the pre-work selection.
- Installing tools or packages in the disposable checkout.

## Success Metrics

1. Success Metric: on a committed graph with one completed Task whose commands
   would pass and would exit 127, and pending, in-progress and failed Tasks
   whose commands fail, `spec check --run-verification` exits `0`, never runs
   the completed Task's commands, lists them `not probed` in text and under
   `verification.notProbed` in JSON, and probes the other three Tasks.
2. Success Metric: a pending Task whose command is vacuous, unknown or
   malformed still makes the check exit `1` with that verdict when a completed
   Task sits beside it, and a pending QA gate's derived command is still
   probed.
3. Success Metric: with every Task completed, the check exits `0` with
   `executed: 0; not probed (completed Tasks): <n>` and no "No authored
   Verification commands." line; without a completed Task, the text and JSON
   bytes are those of today (no `notProbed` key).
4. Success Metric: a completed Task file modified in the working tree is still
   named on an `Uncommitted Verification source` line.
5. Success Metric: replacing `TaskAwaitsWork`'s rule with "always true" makes
   both the new CLI test and the Daemon's
   `TestTaskCycleRerunsStaleTasksAndSkipsCompletedTasks` fail, proving the two
   callers share it.
6. Success Metric: the existing `--run-verification`, provenance and verdict
   line tests pass unchanged, the documentation contracts pass with the
   Behavior Surface record re-recorded last, and both owned-skill versions are
   recorded.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The Vortex report in the adopted Backlog Entry, written before this Spec
  existed: completed Task 01 probed at `gesttione-solutions/vortex@a961a29`,
  exit 127, verdict `unknown`, final exit `1`, "A execução normal por
  implement continuou somente pelas Tasks pendentes."
- `docs/references/2026-08-14-a-clean-checker-left-eight-vacuous-commands-standing.md`,
  a finding measured in Fluxus and Roundfix before this Spec: "a check that
  cannot be complete should say what it did not cover." The not-probed lines
  and verdict-line count follow it.
- The Test Anything Protocol 14 specification
  (https://testanything.org/tap-version-14-specification.html), read on
  2026-10-26: SKIP test points "indicate that a test was not run", "Harnesses
  must not treat failing `SKIP` test points as a test failure", and "Harnesses
  should report `SKIP` test points found as a list of items that were not
  tested". The not-probed list follows the same contract: named, with a
  reason, never a failure.

## Unreachable Acceptance

- criterion: after release, Vortex's
  `roundfix spec check <slug> --strict --run-verification` on a resumed Spec
  with a completed Task reports that Task's commands `not probed` and exits
  `0` when its pending commands are honest.
  reason: it needs the released binary in another repository, which no
  Verification or QA row may read.
  satisfied-by: the operator's release note to `default/gesttione-solutions-vortex`
  through ai-memory and Vortex's next resumed Spec Check, recorded session
  evidence.

## Validation

The decisions were checked on 2026-10-26 against these sources.

- Reproduction at `86cb6b73` with `bin/roundfix` built in the authoring clone,
  `HOME` and `TMPDIR` at fresh temporary directories: a temporary repository
  held Spec `demo` with Task 01 `completed` (`test -f done.txt`, file
  committed), Task 02 `pending` (`test -f todo.txt`) and Task 03 `completed`
  (`no-such-tool-xyz --run`). `spec check demo --run-verification` printed
  `- task_01: vacuous`, `- task_02: honest` and `- task_03: unknown … (exit
  127)`, and JSON (`schema: roundfix-speccheck/v1`) carried all three under
  `verification.commands`.
- Root cause, read at `86cb6b73`: `internal/cli/spec_check.go`
  `probeSpecVerifications` loops over every `graph.Tasks` entry. The Daemon's
  `initialTaskRunStatuses` (`internal/daemon/task_engine.go`) marks a
  `completed` Task completed and every other status pending, `nextReadyTask`
  starts only pending Tasks, `executeTaskWorker` runs `verifyTaskPreWork`, and
  the QA gate is routed only while its status is not `completed`. So the
  Daemon never probes a completed Task, and the check does.
- A prototype of both Tasks in a scratch clone passed the full
  `./internal/cli`, `./internal/daemon`, `./internal/speccheck`, `./skills`
  and `-tags docscontract ./internal/docscontract` packages, both Tasks'
  Verification commands exited `0` on it and non-zero on `86cb6b73`, the
  documentation Task changed exactly its eleven declared files, and each
  sabotage in `_techspec.md` → Testing Approach failed its named subtest.
- Secondbrain, read only: `qmd query "roundfix spec check run-verification
  vacuous completed task probe" --all --files --min-score 0.3` returned the
  Roundfix mirror's
  `docs/references/2026-08-14-a-clean-checker-left-eight-vacuous-commands-standing.md`
  (score 1.00, cited above), the `write-tasks` and `qa-gate` skills and
  `CONTEXT.md`; none records a decision to probe completed Tasks.
- Exa: the TAP 14 and TAP 13 specifications on the SKIP directive, quoted
  under Acceptance evidence.

## Glossary

- changes: **Vacuous Verification**

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AQ batch 1 briefing for Spec 0319 of 2026-10-26, the maintainer's cycle AQ answers quoted in `_authorization.md`, and the characterization above.

Taken from those sources, with no further interview:

- Selection: every Task whose status is not `completed` is probed; pending,
  in-progress and failed Tasks are probed because the Daemon runs them again,
  and a reopened QA gate is pending and stays probed.
- One exported Daemon function decides, used by the Daemon's scheduler, its QA
  routing and the check, instead of a second comparison in `internal/cli`.
- Reporting: a separate, optional `verification.notProbed` list and a
  `not probed` text line, instead of a new value in `verdict`, so the four
  verdicts of ADR-0148 and their consumers stay unchanged and every output
  without a completed Task keeps its bytes.
- No audit flag; no new ADR: ADR-0148 gets a dated amendment.

## Open Questions

None.
