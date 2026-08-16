---
task: task_03
spec: 0096-a-failure-the-agent-can-read
status: completed # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: high
---

# Task 03: Name a failure the loop has already seen

## Overview

A Task that fails twice with an identical assertion is reported as new both
times, and a whole Run was spent on 2026-08-08 reproducing a diagnostic the loop
had already recorded. This slice compares the new failure's signature against the
Work Item's earlier failures in the Run Event Journal and reports the repetition
where both a Supervisor and the Agent will see it.

## Requirements

1. MUST compare a failing Verification's signature against earlier failures of the
   same Work Item, read from the Run Event Journal.
2. MUST report a match in the Verification failure event, naming the earlier Run
   and attempt.
3. MUST report the same match in the repair prompt the Agent receives.
4. MUST report nothing on the first failure of a signature.
5. MUST NOT persist a second store of signatures; the journal already records
   each failure. This is audited on the Task's commit by the QA gate's changed-path
   audit, not by a Verification command — a command runs before the commit exists
   and can only see the working tree.
6. MUST report a failure older than the journal's retention as new rather than
   erroring.

## Subtasks

- [ ] Read the Work Item's earlier failures from the journal.
- [ ] Report the match in the event and in the prompt.
- [ ] Cover the first, the repeat, and the beyond-retention cases.

## Acceptance Criteria

- [ ] A second failure with a matching signature is reported as repeated, naming
      the earlier Run and attempt.
- [ ] The first failure of a signature reports no repetition.
- [ ] Both the event payload and the repair prompt carry the repetition.
- [ ] A failure whose predecessor is outside retention reports as new, with no
      error.
- [ ] No new store or column was added.

## Verification

- `go test -count=1 ./internal/daemon -run 'TestRepeatedFailure' -v > /tmp/0096-t03.log 2>&1; s=$?; grep -q '^--- PASS: TestRepeatedFailure' /tmp/0096-t03.log || { cat /tmp/0096-t03.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing test; fails today, where no such test exists.
- `test -s /tmp/0096-t03.log || { echo 'the suite produced no output'; exit 1; }; grep -qi 'no tests to run' /tmp/0096-t03.log && { echo 'the suite selected no cases'; exit 1; }; grep -c '^    --- PASS\|^--- PASS' /tmp/0096-t03.log > /tmp/0096-t03-n.txt; test "$(cat /tmp/0096-t03-n.txt)" -ge 4 || { echo "expected the first, repeat, prompt, and beyond-retention cases, got $(cat /tmp/0096-t03-n.txt)"; cat /tmp/0096-t03.log; exit 1; }` — expected: exits 0, refusing a vacuous run and proving each direction runs on its own.
- `grep -rq 'Repeated' internal/agent --include='*.go' || { echo 'the repair prompt carries no repetition'; exit 1; }; n=$(grep -rn 'DiagnosticSignature(' internal/daemon --include='*.go' | grep -v '_test.go' | wc -l | tr -d ' '); test "$n" -ge 1 || { echo 'the signature has no production caller'; exit 1; }` — expected: exits 0, proving the repetition reaches the Agent's prompt and that the signature is consumed in production rather than only declared. Fails today on both clauses.

## Context

- interface: `internal/daemon/task_engine.go`
- interface: `internal/runevent/event.go`

## References

`_techspec.md` → Build Order 3; Risks & Considerations, retention. `_prd.md` →
Core Feature 2; Goal 2; User Story 2; Success Metrics. ADR-0136.

## Result

Implemented retained-history matching at the Task Verification failure boundary.
Each captured failure now records its normalized diagnostic signature in the Run
Event payload. Before publishing that event, the Daemon reads earlier implement
Runs for the same repository, Spec, and Work Item from the Run Event Journal and
looks for the same signature. A match adds a typed `repeated_failure` value and a
summary naming the earlier Run and attempt; the same value reaches the repair
prompt. A missing earlier journal event, including a Run whose events were pruned
by Journal Retention, produces no repetition and no error.

Focused-check evidence:

- Before implementation,
  `rtk rg -n 'RepeatedFailure|Repeated Failure' internal/agent internal/daemon/task_engine.go internal/runevent/event.go`
  returned no matches, establishing that neither the event nor prompt could name
  a repetition.
- `rtk env GOCACHE=/tmp/roundfix-task03-go-cache go test ./internal/agent -run '^TestBuildVerificationRepairPrompt' -count=1`
  reported 9 passing tests. The new prompt case asserts the exact earlier Run,
  attempt, and signature; the existing cases prove a first failure still renders
  without repetition.
- `rtk env GOCACHE=/tmp/roundfix-task03-go-cache go test ./internal/daemon -run '^TestDiagnosticSignature$' -count=1`
  reported 8 passing tests and compiled the journal lookup, event metadata, and
  four-case `TestRepeatedFailure` integration coverage without executing the
  Daemon-owned Task Verification case.
- `rtk env GOCACHE=/tmp/roundfix-task03-go-cache go test ./internal/daemon -run '^(TestTaskCycleRepairReacquiresVerificationCapacityAfterFeedback|TestTaskCycleVerificationFailureRepairsSameSessionAndRerunsFullSequence)$' -count=1`
  reported 2 passing tests, preserving failure-event ordering and the repair
  prompt flow after adding classification.
- `rtk rg -n 'DiagnosticSignature\(|Repeated Failure|repeated_failure' internal/agent/agent.go internal/daemon/engine.go internal/daemon/task_engine.go internal/runevent/event.go`
  found the production signature caller, event payload, event summary, and prompt
  rendering.
- `rtk git diff --name-only` listed only this Task file plus Agent, Daemon, and
  Run Event source/test files. It listed no store, migration, schema, or tooling
  path, proving no second store or column was added.
- `rtk git diff --check` exited 0.

Acceptance evidence:

- `TestRepeatedFailure/matching_failure_event_names_earlier_Run_and_attempt`
  drives a real Run Event Journal entry through a Task cycle and asserts both the
  typed payload and bounded summary name the retained Run and attempt.
- `TestRepeatedFailure/first_failure_reports_no_repetition` asserts the first
  signature has neither a `repeated_failure` payload nor a repeated summary.
- `TestRepeatedFailure/repair_prompt_names_the_same_earlier_Run_and_attempt`
  asserts the repair Agent receives the same match as the event.
- `TestRepeatedFailure/failure_beyond_journal_retention_reports_as_new` models
  the retained Run row after its journal events are gone and asserts the new
  failure remains unannotated without an infrastructure error.
- The implementation writes signatures only into the existing Run Event payload
  and reads them through existing Run/Run Event Journal APIs; no persistence
  schema changed.

The Daemon-owned commands under `## Verification` were not run during this Agent
turn. `Repeated Failure` was already introduced in `CONTEXT.md`, so this slice
does not require a glossary update.
