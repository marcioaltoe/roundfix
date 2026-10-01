---
task: task_06
spec: 0202-a-qa-gate-that-reruns-only-stale-rows
status: completed
type: test
complexity: medium
---

# Task 06: Two real gate passes carry a row into the committed report

## Overview

The QA gate of Run `run_20261001T133454Z_7983d7d312c8ec5c` blocked row 5 with finding F1. `TestQAGateCarriesARowFromAnUnintegratedFailedPass` builds the first pass's Evidence Snapshot by hand, and its fake Agent replaces the second pass's seeded report with a hollow passing report. So no test proves that a first real gate records a row and that the next real gate keeps it in its committed report. QA found no production defect: this Task adds the missing proof.

## Requirements

1. MUST add `internal/daemon/qa_two_pass_carry_test.go` with `TestTwoGatePassesCarryAnUnmovedRowIntoTheCommittedReport`. Through the existing task-cycle fixture it MUST:
   - drive two `TaskCycle` QA passes over a disposable repository;
   - in the first pass, use a fake Agent that writes a report with one passing carriable row declaring its inputs and one failing row, and let the Daemon record the Evidence Snapshot itself, never by hand;
   - leave the first row's inputs unmoved and change only a file the failing row depends on;
   - in the second pass, use a fake Agent that honors the carried rows the Daemon seeds, keeping them unchanged, and re-runs only the rest.
2. MUST assert, for both passes, the `prior_report`, `mechanical` and `evidence_snapshots` `daemon.qa` events with their payloads. The second pass's `prior_report` must name the first pass, and its `evidence_snapshots` must cover the re-run row.
3. MUST read the second pass's committed QA Report independently from Git (`git show <commit>:<path>`). It MUST assert that the carried row keeps the first pass's provenance cell and is listed as carried in `## Row carry-forward`, and that the changed row is re-run with its disposition.
4. MUST add `TestTwoGatePassesReRunARowWhoseInputMoved`, the mirror case: when the carried row's input moves between the passes, the second pass re-runs it as `input moved:` and does not keep it.
5. MUST NOT change production code or any existing test.

## Subtasks

- [ ] Build the two-pass fixture with honoring fake Agents.
- [ ] Assert the events and read the committed report from Git.
- [ ] Add the moved-input mirror case.

## Acceptance Criteria

- [ ] Both new tests pass. Sabotaging the carry, so that the seed drops the carried row, makes the first test fail.

## Context

- interface: `internal/daemon/task_engine_test.go`
- creates: `internal/daemon/qa_two_pass_carry_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestTwoGatePassesCarryAnUnmovedRowIntoTheCommittedReport|TestTwoGatePassesReRunARowWhoseInputMoved)$' ./internal/daemon 2>&1)" || { printf "%s\n" "$out"; exit 1; }; for name in TestTwoGatePassesCarryAnUnmovedRowIntoTheCommittedReport TestTwoGatePassesReRunARowWhoseInputMoved; do printf "%s\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task neither test exists.

## References

- task_01, task_02, task_03
- QA Report of Run `run_20261001T133454Z_7983d7d312c8ec5c`, finding F1

## Result

Added both requested tests in `internal/daemon/qa_two_pass_carry_test.go`.
Each drives two real `TaskCycle` QA passes over the existing disposable Git
fixture, using the real snapshot recorder, carry resolver, Git committer and
Run Event Journal. The first Agent declares inputs for a passing row and a
failing row without writing a snapshot. The second Agent preserves seeded
carried rows and records fresh results only for the remaining rows.

- Unmoved-row evidence: the first failed QA commit remains on its Run Branch;
  the correction changes only `repair-input.txt`. The second committed report
  is read with `git show <commit>:<path>` and must retain row 1's establishing
  report, head and provenance, list it as carried, and record row 2's fresh
  result and `re-run: not pass` disposition. The imported first report must
  remain byte-identical.
- Moved-input evidence: the mirror also changes `carry-input.txt`. The second
  report must replace row 1's old provenance with a fresh observation and
  record `re-run: input moved: carry-input.txt`; both rows are re-run.
- Both passes assert the ordered `prior_report`, `mechanical`, `verdict` and
  `evidence_snapshots` events, including the prior commit/report, carry/re-run
  counts, audited heads and recorded-row counts. Independent YAML decoding
  checks the committed snapshots' row IDs, input paths and SHA-256 digests,
  including the corrected row in the second pass.

Focused checks:

- `GOCACHE=/private/tmp/roundfix-task06-gocache rtk proxy go test ./internal/daemon -run '^TestTwoGatePasses' -count=1`
  exited 0; both new tests ran.
- Mutation check:
  `GOCACHE=/private/tmp/roundfix-task06-gocache rtk proxy go test -overlay=/private/tmp/roundfix-task06-overlay.json ./internal/daemon -run '^TestTwoGatePassesCarry' -count=1`
  exited 1 with `committed report lost carried row or provenance`. The
  temporary overlay omits carried Results rows from the seed renderer while
  retaining the carry disposition. It changes no repository production file
  and demonstrates that a dropped carried row is detected in the committed
  report, even when the Agent re-runs it to a passing result.
- `rtk proxy git diff --check` exited 0.
- `GOCACHE=/private/tmp/roundfix-task06-gocache rtk make verify-incremental`
  exited 0 with expanded sandbox permissions: formatting, vet, package tests
  (including `internal/daemon`), skill checks and build passed. The initial
  sandboxed attempt was blocked by network access to `cafe.github.com`; the
  expanded-permission rerun supplied the conclusive incremental evidence.

No production code or existing test was edited. Task status and declared
Verification remain Daemon-owned; the declared Verification command was not
run. No follow-up work was identified.
