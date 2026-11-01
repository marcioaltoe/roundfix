---
spec: 0337-the-spec-check-leaves-the-qa-gate-unprobed-and-prints-its-expected-verification
prd: _prd.md
created: 2026-11-01
---

# The Spec Check leaves the QA gate unprobed and prints its expected Verification — Technical Spec

## Executive Summary

`probeSpecVerifications` runs the shared prober over every Task that is not
completed, the QA gate included, while the Daemon removes the gate from its
scheduler in `taskPlanWithoutQAGate` and never probes it. The fix exports one
predicate, `daemon.TaskIsQAGate`, makes `taskPlanWithoutQAGate` read it, and
makes the Spec Check report the gate's commands `not probed` with reason
`qa_gate` instead of running them. Separately, `speccheck.Finding` gains an
optional `Expected` text, filled by `AuthoredQAVerification` with the exact
`## Verification` section the loader accepts; `RenderText` prints it as a
fenced block after `fix:` and JSON carries it as `expected`. The
`write-tasks` template and skill, the Roundfix spec reference, the spec guide,
`CONTEXT.md` and ADR-0148 state the rules.

The trade-off accepted: the gate is not probed even when it has no QA Report
yet, so the check no longer shows `honest` for it. That verdict proved
nothing the author controls, because the command is derived (`_prd.md` →
Decisions).

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the JSON gains one optional field within `roundfix-speccheck/v1`.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential is added. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0148 shares one prober and one Task
  selection between the Daemon and the check, ADR-0148: "both callers probe the same Tasks"; the design keeps both and
  amends ADR-0148, as `_prd.md` → Project Constraints states and ADR
  Alternatives below relates. ADR-0315 cites ADR-0148 for the shared prober
  and governs pytest selectors at authoring; this Spec does not touch it.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AT answers of 2026-11-01 quoted in the authorization record; the Governed
  Paths are the Roundfix and `write-tasks` skill files the Tasks change and
  `internal/speccheck/coherence.go`; no Makefile, workflow, module or lint
  configuration changes. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0337-the-spec-check-leaves-the-qa-gate-unprobed-and-prints-its-expected-verification/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/spec.md`, `.agents/skills/write-tasks/SKILL.md`, `.agents/skills/write-tasks/references/task-template.md`, `internal/speccheck/coherence.go`, `skills/roundfix/SKILL.md`, `skills/write-tasks/SKILL.md`.

## Current behavior

Read on 2026-11-01 at `340c9004`; `_prd.md` → Validation records the runs.

1. `probeSpecVerifications` (`internal/cli/spec_check.go`) loads the graph and,
   for each Task, reports every command `not probed` with reason `completed`
   when `daemon.TaskAwaitsWork` is false, and otherwise runs
   `daemon.ProbeCommands` over `task.Verification`. A QA Task's
   `Verification` is always `spec.DerivedQAVerification(slug)`
   (`taskVerification`, `internal/spec/task.go`), so the probe runs the
   derived command.
2. The derived command finds the newest dated `qa-report-*.md` under
   `docs/specs/<slug>/qa` and exits 0 when its frontmatter has one non-empty
   `verdict:`. A kept `verdict: fail` report makes it exit 0 before any work,
   which the probe classifies `vacuous`; the check exits 1.
3. `Engine.TaskCycle` (`internal/daemon/task_engine.go`) calls
   `taskPlanWithoutQAGate`, which moves every Task with
   `task.Type == spec.TaskTypeQA` out of the scheduled plan. The pre-work
   probe `verifyTaskPreWork` is called only from `executeTaskWorker`, so the
   gate is never probed; `runQAGate` judges it through the QA Report and
   `spec.QAReportEligibility`.
4. `speccheck.RenderText` prints a clean verdict line
   `No findings. Authored Verification commands executed: <n>` with
   `; not probed (completed Tasks): <n>` from `VerificationCoverage.NotProbed`,
   which `renderSpecCheckText` sets to `len(verification.NotProbed)`.
5. `AuthoredQAVerification` (`internal/speccheck/coherence.go`) returns a
   finding with code, severity, summary, location and the fix "Regenerate the
   qa Task so its ## Verification renders Roundfix's derived command
   unchanged."; `Finding` has no field for the expected text, and the CLI's
   JSON marshals `[]speccheck.Finding` directly.
6. The loader reads a Verification command from the first code span of each
   `- ` bullet under `## Verification`, so the section
   `## Verification`, a blank line and ``- `<derived command>` `` loads to
   exactly the derived command. The derived command holds no backtick.
7. The `write-tasks` template's `## Verification` guidance and `SKILL.md`
   address authored commands only; neither mentions the QA gate's derived
   Verification.

## System Architecture

No package, command or flag is added. Changes by component:

- `internal/daemon/verification_probe.go`: the exported predicate.
- `internal/daemon/task_engine.go`: `taskPlanWithoutQAGate` reads it.
- `internal/cli/spec_check.go`: the probe skips the gate and counts reasons.
- `internal/speccheck/report.go`: the QA gate count in the verdict line, the
  `Expected` field and its text rendering.
- `internal/speccheck/coherence.go`: the finding fills `Expected`.
- Tests in `internal/cli` and `internal/speccheck`; the guide, the two
  skills, `CONTEXT.md` and ADR-0148.

## Implementation Design

### Interfaces

```go
// internal/daemon/verification_probe.go

// TaskIsQAGate reports whether task is the Spec's authored QA gate, the Task
// of type qa, which the Daemon routes to its QA step and never probes.
func TaskIsQAGate(task spec.Task) bool

// internal/speccheck/report.go
type Finding struct {
	// ...existing fields unchanged...
	Expected string `json:"expected,omitempty"`
}

type VerificationCoverage struct {
	Ran        bool
	Executed   int
	Unexecuted int
	NotProbed  int // completed Tasks
	QAGate     int // pending, in-progress or failed QA gate commands
}

// internal/speccheck/coherence.go

// ExpectedQAVerificationBlock renders the ## Verification section a qa Task
// must carry for the derived commands.
func ExpectedQAVerificationBlock(derived []string) string
```

### Data Models

`roundfix-speccheck/v1` gains one optional key: `findings[].expected`, a
string, present only when non-empty. `verification.notProbed[].reason` gains
the value `qa_gate` beside `completed`. No schema version changes: both are
additive.

### API Contracts

1. API Contract: the predicate. `TaskIsQAGate(task)` is
   `task.Type == spec.TaskTypeQA`. `taskPlanWithoutQAGate` moves a Task out of
   the scheduled plan exactly when it returns true.
2. API Contract: the probe. In `probeSpecVerifications`, a Task that
   `TaskAwaitsWork` refuses keeps reason `completed`. Otherwise a Task that
   `TaskIsQAGate` accepts is never run: each command is appended to
   `verification.notProbed` as `{task, command, reason: "qa_gate"}` and to the
   text as `- <task>: not probed — "<command>" (QA gate)`, in graph order with
   the other lines. Every other Task is probed as today. A not-probed command
   never changes the exit code.
3. API Contract: the verdict line. `renderSpecCheckText` sets
   `NotProbed` to the count of reason `completed` and `QAGate` to the count of
   reason `qa_gate`; `RenderText` appends `; not probed (QA gate): <n>` after
   the completed count when `QAGate > 0`. With no QA gate entry every line
   keeps its bytes.
4. API Contract: the expected block. `ExpectedQAVerificationBlock` returns
   `## Verification\n\n` followed by ``- `<command>`\n`` per command.
   `AuthoredQAVerification` sets `Expected` to the block of
   `task.Verification`, which is the derived command, and its fix to "Replace
   the qa Task's ## Verification section with the expected block this finding
   prints, unchanged."; code, severity, summary and location keep their bytes.
5. API Contract: the rendering. When `Expected` is non-empty, `RenderText`
   writes, right after the `  fix: ` line, the line `  expected block:`, a
   fence line of three backticks followed by `markdown`, the block, and a
   closing fence line of three backticks. JSON writes the
   block under `expected`; every finding with an empty `Expected` renders as
   today in both formats.
6. API Contract: the rules in writing.
   - `docs/user-guide/commands/spec.md` and
     `.agents/skills/roundfix/references/spec.md`: the paragraph beginning
     "The probe runs only the commands" states that the QA gate is never
     probed, with its text line, verdict count and reason `qa_gate`, and no
     longer says a pending QA gate's vacuous command exits 1. The guide's
     code table and the reference's code list gain
     `SC-QA-VERIFICATION-AUTHORED` with the `expected block:` text and the
     `expected` JSON field.
   - `.agents/skills/write-tasks/references/task-template.md`: a comment in
     the `## Verification` section tells the author of a `type: qa` Task to
     write a placeholder bullet, run `roundfix spec check <slug> --stage
     tasks` and replace the section with the finding's expected block.
   - `.agents/skills/write-tasks/SKILL.md`: a section `## QA Task
     Verification` says the same and that the check never probes the command.
   - `CONTEXT.md`: **Vacuous Verification** adds that the QA gate is never
     probed and is reported with reason `qa_gate`.
   - ADR-0148: an amendment dated 2026-11-01 states that neither caller
     probes the QA gate and names `daemon.TaskIsQAGate`.

### Surface Transcripts

None. The changed lines are asserted by the named tests; no whole-command
transcript is declared.

### Invariants

```text
1. No QA gate command is executed by spec check --run-verification, whatever its status and reports.
2. Every non-QA Task that is not completed is probed and classified as before.
3. The Daemon and the check decide the QA gate through TaskIsQAGate.
4. SC-QA-VERIFICATION-AUTHORED still fires for every differing QA Verification, and the section its expected block holds loads to the derived command.
5. A finding without an expected block renders the same bytes in text and JSON as before.
```

### Declared breaks

- `spec check --run-verification` no longer executes a pending QA gate's
  command. The subtest "a pending QA gate is still probed" of
  `TestSpecCheckLeavesCompletedTasksAlone`
  (`internal/cli/spec_check_completed_tasks_test.go`) asserts the old
  behavior and is removed; its top-level test stays, so the coverage record
  is unchanged. The new `TestSpecCheckNeverProbesTheQAGate` covers the case.
- The fix text of `SC-QA-VERIFICATION-AUTHORED` changes. No test asserts it
  (`TestAuthoredQAVerificationIsRefusedByTaskName` reads the summary only).
- The guide and the reference stop saying that a vacuous pending QA gate
  makes the check exit 1.

## Coverage Map

- Goal 1, Core Feature 2, Success Metrics 1 and 2 → `probeSpecVerifications`
  and the verdict line (API Contracts 2 and 3), Invariants 1 and 2.
- Goal 2, Core Feature 1, Success Metric 3 → `TaskIsQAGate` and
  `taskPlanWithoutQAGate` (API Contract 1), Invariants 3.
- Goal 3, Success Metric 2 → `probeSpecVerifications` (API Contract 2),
  Invariants 2.
- Goal 4, Core Feature 3, Success Metrics 4 and 5 →
  `ExpectedQAVerificationBlock`, `AuthoredQAVerification` and `RenderText`
  (API Contracts 4 and 5), Invariants 4 and 5.
- Goal 5, Core Feature 4, Success Metric 6 → the guide, the references, the
  template, `SKILL.md`, `CONTEXT.md` and ADR-0148 (API Contract 6).

Deviations from the PRD: None.

## Absence Evidence

1. A shared QA gate predicate — command: `grep -rn 'func TaskIsQAGate\|func IsQAGate' --include='*.go' internal`; result: no match; `taskPlanWithoutQAGate` tests `task.Type != spec.TaskTypeQA` inline.
2. A QA branch in the Spec Check probe — command: `grep -n 'TaskTypeQA\|qa_gate' internal/cli/spec_check.go`; result: no match; the only skip is `daemon.TaskAwaitsWork`.
3. An expected text on a finding — command: `grep -n 'Expected\|expected' internal/speccheck/report.go`; result: no match; `Finding` holds code, severity, summary, where and fix.
4. QA Verification guidance in the authoring skill — command: `grep -rn 'SC-QA-VERIFICATION-AUTHORED\|QA Task Verification' .agents/skills/write-tasks`; result: no match; the only `derived` there is the derived Task Type.

## ADR Alternatives

- ADR-0315: unaffected; its rejected alternatives (explaining the vacuous
  verdict from runner output, a rule for every test runner, a gap instead of
  an error) stay rejected: the prober still gains no cause text, and this
  Spec only stops running the one command an author cannot change.
- ADR-0148: followed and amended; it keeps one prober and one Task selection
  for both callers, and its rejected alternative, a second implementation in
  the check free to disagree with the Daemon, stays rejected: the check now
  reads the Daemon's own QA predicate instead of keeping a type test of its
  own.

## Integration Points

None external. The check reads local files and Git; the finding is built in
process.

## Testing Approach

All tests are hermetic: temporary Git repositories through the existing
`newSpecCheckStatusWorkspace` and `probeProvenanceCLI` helpers, in-memory
findings, fixture files under `t.TempDir()`, and no network.

- `internal/cli/spec_check_completed_tasks_test.go`:
  `TestSpecCheckNeverProbesTheQAGate`, with a completed `task_01`
  (`true`), a pending `task_02` (`test -f corrective-output`) and a pending QA
  `task_03`:
  - subtest "a preserved failed QA report beside an honest corrective Task":
    a committed `docs/specs/clean/qa/qa-report-2026-10-31.md` with
    `verdict: fail`; text exit 0 with the verdict line
    `No findings. Authored Verification commands executed: 1; not probed (completed Tasks): 1; not probed (QA gate): 1.`
    and the three ordered lines ending in `(QA gate)`; JSON exit 0, one
    `honest` command for `task_02`, and `notProbed[1]` naming `task_03`, the
    derived command and reason `qa_gate`;
  - subtest "no QA report yet": the same graph without the report, the same
    assertions;
  - subtest "a vacuous corrective Task beside the QA gate still fails": a
    pending `true` beside a pending QA gate gives exit 1 with one `vacuous`
    command for `task_01`.
  The subtest "a pending QA gate is still probed" of
  `TestSpecCheckLeavesCompletedTasksAlone` is removed; the rest of that test
  stays unchanged and green.
- `internal/speccheck/coherence_test.go`:
  `TestAuthoredQAVerificationPrintsTheExpectedBlock`: a QA Task whose
  Verification is ``- `<command>` — expected: exit 0`` gives one finding whose
  `Expected` equals `## Verification\n\n- ` + code span of
  `spec.DerivedQAVerification(slug)[0]`; the same Task file with its section
  replaced by `Expected` gives no `SC-QA-VERIFICATION-AUTHORED`.
- `internal/speccheck/report_test.go`: `TestAFindingRendersItsExpectedBlock`:
  `RenderText` of one finding with `Expected` and one without equals the exact
  text with the fenced block after the first `fix:` line only; `RenderJSON`
  carries `expected` on the first finding and no `expected` key on the second.
- Existing suites: `./internal/daemon` (QA routing),
  `TestAuthoredQAVerificationIsRefusedByTaskName`,
  `TestDerivedQAVerificationIsAcceptedByTaskStageChecker`, the
  `VerificationCoverage` cases in `internal/speccheck/report_test.go`,
  `./skills`, and `-tags docscontract ./internal/docscontract`.
- Sabotages, each run on the prototype:
  - `TaskIsQAGate` always false: `TestSpecCheckNeverProbesTheQAGate` fails on
    the preserved-report subtest (the gate is `vacuous`) and
    `./internal/daemon` fails, among others `TestAQATaskRecordsNoPaths`;
  - the QA branch removed from `probeSpecVerifications`:
    `TestSpecCheckNeverProbesTheQAGate` fails on all three subtests;
  - reason `completed` written for the QA gate:
    `TestSpecCheckNeverProbesTheQAGate` fails on the verdict line;
  - `Expected` left empty in `AuthoredQAVerification`:
    `TestAuthoredQAVerificationPrintsTheExpectedBlock` fails;
  - the bullet written without its code span:
    `TestAuthoredQAVerificationPrintsTheExpectedBlock` fails;
  - `omitempty` dropped from the `expected` tag:
    `TestAFindingRendersItsExpectedBlock` fails on the second finding.

## Build Order

1. The QA gate predicate, the probe change, the verdict count, the test, the
   probe paragraph of the guide and the reference, `CONTEXT.md`, ADR-0148,
   with the Roundfix skill version and the fingerprints re-recorded last
   (task_01).
2. The expected block, its rendering, the tests, the code row and bullet in
   the guide and the reference, the `write-tasks` template and `SKILL.md`,
   with both skill versions and the fingerprints re-recorded last (task_02)
   (depends on: 1).
3. QA gate (task_03) (depends on: 2).

task_02 follows task_01 because both edit
`docs/user-guide/commands/spec.md`, the Roundfix spec reference, its
`SKILL.md` pair, `internal/speccheck/report.go`, the owned-skill version
record and the Behavior Surface Record.

## Risks & Considerations

- An adopter script that parsed the verdict line for a QA gate's `honest`
  line loses it; the JSON `notProbed` entry with reason `qa_gate` replaces
  it, and the guide states the change.
- `DerivedQAVerification` hardcodes `docs/specs/<slug>/qa` even under a
  non-default Spec Root. The expected block copies the same command the loader
  compares against, so copying it always clears the finding; the path itself
  is out of scope.
- A text consumer that splits findings on lines starting with `[` still
  works: the fenced block's lines start with `#`, `-`, a blank or a fence.

## Glossary

- changes: **Vacuous Verification** — the QA gate is never probed either

## Decisions

- One exported predicate, `daemon.TaskIsQAGate`, read by the Daemon's routing
  and the Spec Check.
- The gate is never probed, whatever its reports; a completed gate keeps
  reason `completed`.
- `Expected` is an optional field on every finding, filled only by
  `SC-QA-VERIFICATION-AUTHORED`; the text prints it as a fenced block after
  `fix:`.
- Secondbrain: `wiki/index.md` read; the `qmd query` named in `_prd.md` →
  Validation found no page deciding either rule. Exa: the rustc JSON
  diagnostics page, which carries an exact suggested replacement as its own
  optional field, the shape `expected` follows.
- Two implementation Tasks in series and a QA gate.
