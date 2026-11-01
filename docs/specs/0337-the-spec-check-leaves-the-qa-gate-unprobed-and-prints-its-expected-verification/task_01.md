---
task: task_01
spec: 0337-the-spec-check-leaves-the-qa-gate-unprobed-and-prints-its-expected-verification
status: completed
type: backend
complexity: medium
---

# Task 01: The Spec Check leaves the QA gate unprobed through the Daemon's QA predicate

## Overview

The Daemon package exports the rule that a Task of type `qa` is the QA gate,
and both the Daemon's routing and `roundfix spec check --run-verification`
read it. The check stops running the gate's derived Verification, which a QA
Report kept from an earlier Run makes pass before any work, and reports it
`not probed` with reason `qa_gate`. The guide, the Roundfix skill's spec
reference, `CONTEXT.md` and ADR-0148 state the rule. This answers the adopted
reference of 2026-11-01, "The Spec Check probe judges the QA gate's derived
Verification vacuous"
([references](references/2026-11-01-the-spec-check-probe-judges-the-qa-gates-derived-verification-vacuous.md)),
the Fluxus report of 2026-11-01.

## Requirements

1. MUST add `TaskIsQAGate(task spec.Task) bool` to
   `internal/daemon/verification_probe.go`, beside `TaskAwaitsWork`, returning
   `task.Type == spec.TaskTypeQA`, with a doc comment that says the Daemon
   routes the gate to its QA step and never probes it (`_techspec.md` → API
   Contract 1).
2. MUST make `taskPlanWithoutQAGate` in `internal/daemon/task_engine.go` decide
   with `TaskIsQAGate(task)` instead of its inline type test, changing no
   other line of that file.
3. MUST change `probeSpecVerifications` in `internal/cli/spec_check.go` as
   API Contract 2 states: after the completed branch, a Task that
   `daemon.TaskIsQAGate` accepts is never run, and each of its commands is
   appended to `NotProbed` with reason `qa_gate` and to the text as
   `- <task>: not probed — "<command>" (QA gate)`. Update the reason comment
   of `specCheckVerificationNotProbed`.
4. MUST add the field `QAGate int` to `speccheck.VerificationCoverage`
   (`internal/speccheck/report.go`) and make `RenderText` append
   `; not probed (QA gate): <n>` after the completed count when it is above
   zero; `renderSpecCheckText` sets `NotProbed` to the count of reason
   `completed` and `QAGate` to the count of reason `qa_gate` (API Contract 3).
   Every report without a QA gate entry keeps its bytes.
5. MUST add `TestSpecCheckNeverProbesTheQAGate` to
   `internal/cli/spec_check_completed_tasks_test.go`, calling `t.Parallel()`
   first and using `newSpecCheckStatusWorkspace` and `probeProvenanceCLI`,
   with the three subtests `_techspec.md` → Testing Approach lists: the kept
   `verdict: fail` report committed at
   `docs/specs/clean/qa/qa-report-2026-10-31.md`, no report, and a vacuous
   pending non-QA Task beside the gate. Each asserts the exit code, the exact
   verdict line and ordered text lines, and the JSON `commands` and
   `notProbed` entries, reason `qa_gate` included.
6. MUST remove the subtest "a pending QA gate is still probed" from
   `TestSpecCheckLeavesCompletedTasksAlone` (a declared break,
   `_techspec.md` → Declared breaks) and leave every other subtest of that
   test unchanged.
7. MUST replace, in `docs/user-guide/commands/spec.md` and the canonical
   `.agents/skills/roundfix/references/spec.md`, the paragraph that begins
   "The probe runs only the commands of Tasks a Run would still execute" with
   one that keeps its completed-Task text and states: the QA gate is never
   probed, whatever its status, because the Daemon routes the Task of type
   `qa` to its QA step and Roundfix derives its Verification; a pending gate's
   command is reported as `- <task>: not probed — "<command>" (QA gate)`; the
   clean verdict line adds `; not probed (QA gate): <n>`; JSON lists it with
   reason `qa_gate`; a completed gate keeps reason `completed`. The sentence
   "The QA gate is never probed, whatever its status" appears verbatim, and
   the paragraph no longer says that a vacuous command of a pending QA gate
   makes the check exit 1. Change no other paragraph.
8. MUST revise **Vacuous Verification** in `CONTEXT.md` through
   `domain-modeling`, adding the sentence "The QA gate is never probed
   either, whatever its status:" followed by the reason (the Daemon routes it
   to the QA step and its derived Verification would pass on a kept QA
   Report) and that the check reports it `not probed` with reason `qa_gate`.
9. MUST append to
   `docs/adr/0148-one-verification-prober-refuses-vacuous-tasks-before-agent-and-at-authoring.md`
   an amendment beginning "Amendment (2026-11-01): neither caller probes the
   QA gate." that names `daemon.TaskIsQAGate` and states why, and set its
   `updated_at` to `2026-11-01T00:00:00Z`. Change no other line.
10. MUST copy the canonical skill with `make skills-sync`, then raise both
    version fields of `.agents/skills/roundfix/SKILL.md` and its mirror with
    `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
    -record-skill-versions -count=1`, recording
    `skills/testdata/owned-skill-versions.json`. Never hand-edit a mirror or
    a version field.
11. MUST re-record the Behavior Surface Record as the LAST step, after the
    last edit of the guide, with `go test -count=1 -tags docscontract
    ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$'
    -record-skill-coverage`, and commit `docs/references/behavior-surfaces.json`.
12. MUST NOT touch `spec.DerivedQAVerification`, `spec.QAReportEligibility`,
    `runQAGate`, `docs/references/skill-coverage.json`, any `### QA
    settlement` section, any Baseline module, `docs/agents/`, the Makefile,
    CI workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] Characterize: the Fluxus shape in `newSpecCheckStatusWorkspace` reports the QA gate `vacuous` and exits 1 today.
- [ ] Add the predicate, route the Daemon through it and skip the gate in the probe with the new count.
- [ ] Add the regression test and remove the old subtest.
- [ ] Update the guide, the reference, `CONTEXT.md` and ADR-0148, sync the mirror, record the skill version, re-record the fingerprints last.

## Acceptance Criteria

- [ ] An honest corrective Task beside a pending QA gate and a kept failed QA Report passes `spec check --run-verification`, with the gate `not probed` and reason `qa_gate` in text and JSON.
- [ ] A vacuous, unknown or malformed non-QA command still exits 1, and the Daemon's QA routing reads the same predicate.
- [ ] The guide, both reference copies, `CONTEXT.md` and ADR-0148 carry the rule, the mirrors equal their canonical copies, and the owned-skill version and Behavior Surface records pass.

## Context

- instruction: `docs/adr/0148-one-verification-prober-refuses-vacuous-tasks-before-agent-and-at-authoring.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- interface: `internal/daemon/verification_probe.go`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/cli/spec_check.go`
- interface: `internal/cli/spec_check_completed_tasks_test.go`
- interface: `internal/speccheck/report.go`
- interface: `docs/user-guide/commands/spec.md`
- interface: `.agents/skills/roundfix/references/spec.md`
- interface: `skills/roundfix/references/spec.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`
- interface: `docs/adr/0148-one-verification-prober-refuses-vacuous-tasks-before-agent-and-at-authoring.md`

## Consumers

- The QA gate predicate `daemon.TaskIsQAGate` and the probe's `notProbed` reasons — consumers: `./internal/daemon`, `./internal/cli`
- The verdict line of `speccheck.RenderText` — consumers: `./internal/speccheck`, `./internal/cli`
- The owned-skill version record and the Behavior Surface Record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^TestSpecCheckNeverProbesTheQAGate$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestSpecCheckNeverProbesTheQAGate ' || { echo 'missing pass: TestSpecCheckNeverProbesTheQAGate' >&2; exit 1; }; go test -count=1 -timeout 20m -run 'SpecCheck' ./internal/cli && go test -count=1 -timeout 20m ./internal/daemon ./internal/speccheck` — expected: exit 0. Before this Task the test does not exist, so the pass line is missing; afterwards it passes with every Spec Check test of the CLI and the Daemon and Spec Check suites.
- `grep -qF 'TaskIsQAGate(task)' internal/daemon/task_engine.go || { echo 'the Daemon does not route the QA gate through the shared predicate' >&2; exit 1; }; grep -qF 'daemon.TaskIsQAGate(task)' internal/cli/spec_check.go || { echo 'the Spec Check does not read the shared predicate' >&2; exit 1; }; for f in docs/user-guide/commands/spec.md .agents/skills/roundfix/references/spec.md skills/roundfix/references/spec.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'The QA gate is never probed, whatever its status' || { printf 'missing rule in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Vacuous Verification**:' || { echo 'missing term in CONTEXT.md' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'The QA gate is never probed either, whatever its status:' || { echo 'missing QA gate sentence in CONTEXT.md' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/adr/0148-one-verification-prober-refuses-vacuous-tasks-before-agent-and-at-authoring.md | grep -qF -- 'Amendment (2026-11-01): neither caller probes the QA gate.' || { echo 'missing ADR-0148 amendment' >&2; exit 1; }; cmp .agents/skills/roundfix/references/spec.md skills/roundfix/references/spec.md && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task neither caller reads the predicate, so the first check fails. Afterwards both read it, the guide, both reference copies, the glossary entry and ADR-0148 carry the rule, the mirrors equal their canonical copies, and the owned-skill version and documentation contracts, the Behavior Surface fingerprint among them, pass.

## References

- `_prd.md` → Goals; Core Features 1, 2 and 4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 6; Glossary
- `_techspec.md` → Current behavior 1-4; API Contracts 1-3 and 6; Invariants 1-3; Declared breaks; Testing Approach; Build Order 1


## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-structs-interfaces, golang-testing, testing-boss, systematic-debugging, no-workarounds, domain-modeling, tech-writer, evidence-gate.

Implemented the assigned QA-gate selection and reporting slice. The Daemon
exports `TaskIsQAGate`; its routing reads that predicate with only one line
changed in `task_engine.go`. Spec Check checks completed status first, then
reports every unfinished QA gate command `not probed` with reason `qa_gate`
instead of executing it. The text renderer counts completed and QA gate
commands separately and preserves the old output when no QA gate entry exists.

Acceptance evidence from focused checks:

- Honest corrective Task beside a pending gate: the new regression uses
  committed temporary repositories, with and without a kept failed QA Report,
  and asserts exit 0, the exact verdict line, ordered command lines, and exact
  JSON `commands` and `notProbed` entries. The vacuous corrective case asserts
  exit 1 while the gate remains `qa_gate`. Before the implementation,
  `go test -count=1 -timeout 5m ./internal/cli -run
  '^TestSpecCheckNeverProbesTheQAGate$/a_preserved_failed_QA_report'` reproduced
  exit 1 with `task_03: vacuous` and two executed commands. The first fixture
  attempt lacked the QA directory; adding its explicit creation allowed the
  reproduction to reach the product behavior.
- Focused CLI check: `go test -count=1 -timeout 5m ./internal/cli -run
  '^TestSpecCheck(NeverProbesTheQAGate|LeavesCompletedTasksAlone)$'` exited 0.
  This covers all three new subtests and retained non-QA vacuous, unknown,
  malformed, completed-status and source-provenance cases. Only the explicitly
  superseded pending-QA subtest was removed; a source comparison confirmed every
  other subtest remains byte-identical.
- Shared Daemon routing: `go test -count=1 -timeout 5m ./internal/daemon -run
  '^(TestAQATaskRecordsNoPaths|TestTaskCycleQAStepSkippedUnlessEveryTaskCompleted|TestTaskSchedulerRefusesQATaskAsAgentWork)$'`
  exited 0. Source inspection confirmed that both routing and Spec Check read
  the exported predicate.
- Existing renderer output: `go test -count=1 -timeout 5m ./internal/speccheck
  -run '^(TestVerdictLineStatesProbeCoverage|TestRenderResultTextAndJSON)$'`
  exited 0; the CLI regression checks the new QA gate count and its position
  after the completed count.
- Rule documentation: changed only the probe paragraph of the guide and
  canonical Roundfix spec reference; revised **Vacuous Verification** through
  domain-modeling and appended the ADR-0148 amendment with its requested
  `updated_at`. `go test -count=1 -timeout 5m -tags docscontract
  ./internal/docscontract -run
  '^TestGlossary(PhrasesAreDefinedTerms|LinksNameAcceptedADRs|StatesEachRuleOnce)$'`
  exited 0. Source comparisons confirmed the bounded paragraph edits and
  equality of both canonical/mirror pairs.
- Generated artifacts: `make skills-sync` exited 0, followed by
  `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
  -record-skill-versions -count=1` (exit 0). The sanctioned recorder raised both
  Roundfix version fields in both copies from `0.1.165` to `0.1.166` and
  recorded the new digest. No mirror or version field was hand-edited.
- Final implementation regeneration, after the last guide edit:
  `go test -count=1 -tags docscontract ./internal/docscontract -run
  '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` exited 0 and
  changed only the spec guide fingerprint in the Behavior Surface Record.
  `git diff --check` also exited 0. Commands above were invoked through
  `rtk proxy` to preserve exact output.

The pre-existing Daemon-written `status: in_progress` is preserved. No Task
Graph, other Task, forbidden QA implementation, Baseline, tooling configuration,
or skill-coverage map was changed. Declared Verification and repository gates
remain for the Daemon; no terminal Task verdict is claimed. No follow-up work
was added to this slice.
