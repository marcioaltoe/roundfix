---
task: task_03
spec: 0244-a-formatted-qa-report-and-a-reopen-for-late-dependencies
status: completed
type: backend
complexity: high
---

# Task 03: The QA step formats its QA directory with the Format Command and reverts a failed run

## Overview

Makes the QA step run the configured Format Command over the files under the
Spec's `qa/` directory, at the two points of `_techspec.md` → Invariants 5 and
6, and wires `verification.format` from the configuration into the Run. This
Task answers the Backlog Entry "An unformatted QA report breaks the next Run's
precondition" of 2026-10-06: a failed pass committed an unformatted report, the
next Run imported it, and the repository Verification precondition refused the
gate before any row ran. It is verifiable on its own through the Daemon and
CLI tests below.

## Requirements

1. MUST add `FormatCommand string` to `daemon.TaskPlan` and pass the loaded
   `verification.format` into it from `executeImplementCycle`. When the value
   is empty, the QA step's bytes, commits, Run Events and stderr stay exactly
   as today (`_techspec.md` → Invariant 2).
2. MUST implement the format run in a new `internal/daemon/qa_format.go` per
   `_techspec.md` → Invariants 3, 4, 7 and 8. It keeps only regular files
   under the Spec's `qa/` directory, sorts them, and passes them as separate
   arguments through `sh -c '<command> "$@"' roundfix-format <path>...` from
   the Run Worktree root. It keeps every file's bytes and mode first and caps
   the run at 5 minutes, cancelling the process. On a non-zero exit, a start
   failure, the cap, a Stop Request, or a QA Report whose verdict no longer
   reads the same, it restores every file. It publishes one `daemon.qa` Run
   Event with phase `format` and writes the `roundfix: QA format` stderr line
   for `failed` and `reverted`. A path is never interpolated into the command
   string.
3. MUST call the run with stage `imported` in `runQAGate` after the mechanical
   stage and before `runQARepositoryVerification`, over the imported pass's
   files, only when a prior pass was imported; the import and the carry proof
   keep reading the committed bytes (`_techspec.md` → Invariant 5).
4. MUST call the run with stage `commit` in `commitQAReport` after the
   stageable set is final and before the commit, over the stageable paths,
   checking the report against the settled verdict, and MUST NOT recompute the
   stageable set (`_techspec.md` → Invariant 6).
5. MUST add `internal/daemon/qa_format_test.go`, built on the prior-pass
   fixture with `GitCommitter` and a formatter shell script written to a
   temporary directory, with the tests:
   `TestQAFormatCommitsTheFormattedReport` (the committed report and evidence
   equal the script's output);
   `TestQAFormatFormatsAnImportedPassBeforeThePrecondition` (an imported
   unformatted report is formatted before a repository Verification that
   fails on any unformatted file under `qa/`, and the gate reaches its Agent);
   `TestQAFormatRevertsAFailingFormatter` (a script that rewrites then exits 1
   leaves the original bytes committed, with outcome `failed`);
   `TestQAFormatRevertsAVerdictChange` (a script that rewrites the verdict
   leaves the original bytes committed, with outcome `reverted`);
   `TestQAFormatLeavesFilesOutsideTheQADirectory` (the QA Task file and a path
   outside `qa/` are never passed to the script); and
   `TestQAFormatEmptyCommandRunsNothing` (no `format` event and unchanged
   bytes).
6. MUST add `internal/cli/implement_qa_format_test.go` with
   `TestImplementPassesVerificationFormatToTheQAStep`: `roundfix implement`
   with `verification:` / `format:` naming a script in User Config runs that
   script with the QA Report path among its arguments.
7. MUST NOT format Task commits, Task files or any path outside the Spec's
   `qa/` directory, MUST NOT change the import (`copyPriorQAPass`) or the carry
   proof, and MUST NOT change any existing test expectation.

## Subtasks

- [ ] Add the plan field and wire it from the configuration.
- [ ] Implement the format run with its revert, event and stderr line.
- [ ] Call it at the import and at the QA Report commit.
- [ ] Write the Daemon and CLI tests.

## Acceptance Criteria

- [ ] A configured formatter's output is what the QA Report commit records.
- [ ] An imported unformatted pass no longer refuses the precondition.
- [ ] A failing or verdict-changing formatter costs no verdict and is recorded.
- [ ] An empty command changes nothing.
- [ ] Every existing Daemon prior-pass, carry and implement test passes
      unchanged.

## Context

- interface: `internal/daemon/task_engine.go`
- interface: `internal/cli/implement.go`
- creates: `internal/daemon/qa_format.go`
- creates: `internal/daemon/qa_format_test.go`
- creates: `internal/cli/implement_qa_format_test.go`
- instruction: `internal/daemon/qa_prior_pass.go`
- instruction: `internal/daemon/qa_prior_pass_test.go`
- instruction: `internal/cli/implement_settlement_checks_test.go`
- instruction: `docs/adr/0249-a-qa-step-formats-its-qa-directory-and-reopen-sees-a-late-dependency.md`

## Verification

- `out="$(go test -count=1 -v -run '^TestQAFormat' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for t in TestQAFormatCommitsTheFormattedReport TestQAFormatFormatsAnImportedPassBeforeThePrecondition TestQAFormatRevertsAFailingFormatter TestQAFormatRevertsAVerdictChange TestQAFormatLeavesFilesOutsideTheQADirectory TestQAFormatEmptyCommandRunsNothing; do printf '%s\n' "$out" | grep -q -- "--- PASS: $t " || { printf 'missing passing test %s\n' "$t" >&2; exit 1; }; done` — expected: exit 0; before this Task no such test exists, and after it all six pass.
- `out="$(go test -count=1 -v -run '^TestImplementPassesVerificationFormatToTheQAStep$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestImplementPassesVerificationFormatToTheQAStep ' || { printf 'missing passing test\n' >&2; exit 1; }` — expected: exit 0; before this Task the test does not exist, and after it the wired command receives the QA Report path.

## References

- `_prd.md` → Core Feature 2; Goals; Success Metric 1; Success Metric 2; Success Metric 3
- `_techspec.md` → API Contract 2; Invariants 2 to 8; Build Order 3
- ADR-0249
- ADR-0194

## Result

Implemented the Task 03 slice for Daemon Verification. The existing
`status: in_progress` is preserved; the authored Verification commands were
not run, and no commit, push or Pull Request was made.

### Implementation

- `TaskPlan.FormatCommand` receives the loaded `verification.format` from
  `executeImplementCycle`.
- The imported-pass format runs after the mechanical stage and its carry
  proof, before repository Verification. The commit format runs after the
  stageable set is final, immediately before the QA Report commit, without
  recomputing that set. Import and carry code remain unchanged.
- The format run keeps sorted regular files under the Spec's `qa/` directory,
  rejects symlinks and symlink ancestors, snapshots bytes and modes, and
  supplies repository-relative paths as separate shell arguments. It has a
  five-minute cancellation cap and observes durable Stop Requests.
- A failed command or changed/unreadable report verdict restores every kept
  file's bytes and mode. The newest committed report is checked against the
  settled verdict; older included reports retain their own verdicts.
- Each invocation records its format stage, outcome, command, path count and
  sorted changed paths. Failed/reverted runs record the reason and last 2048
  diagnostic bytes and write the documented stderr line. An empty command
  returns before any formatting work, event or diagnostic.
- `qa_format_unix.go` isolates process-group cancellation so shell children
  stop before restoration; `qa_format_windows.go` retains CommandContext
  process cancellation. These additional ordinary source paths belong to
  this Task's cancellation implementation.
- Test formatter scripts are private non-executable temporary files invoked
  through `sh`, preserving the frozen executable-fixture inventory.

### Acceptance evidence

| Acceptance criterion | Implementation and focused-check evidence |
| --- | --- |
| Configured output is committed | `TestQAFormatCommitsTheFormattedReport` reads real Git blobs for both the report and evidence and matches the script output. `TestImplementPassesVerificationFormatToTheQAStep` exercises User Config through implement and checks every recorded formatter argument is under `qa/`, including the report. |
| Imported unformatted pass reaches the Agent | `TestQAFormatFormatsAnImportedPassBeforeThePrecondition` imports a journal-proven prior pass, runs a real repository command rejecting any unformatted QA file, reaches the QA Agent, and records both format stages. |
| Failures and verdict changes preserve verdict and are recorded | `TestQAFormatRevertsAFailingFormatter` and `TestQAFormatRevertsAVerdictChange` inspect the original committed report/evidence bytes and failure events; the failing case also checks restored modes and diagnostics. Additional tests cover context cancellation, durable Stop Request, failure to start the shell, and the diagnostic-tail limit. |
| Empty command changes nothing | `TestQAFormatEmptyCommandRunsNothing` checks unchanged committed bytes, no format event and no format stderr line. `TestQAFormatLeavesFilesOutsideTheQADirectory` proves the QA Task, an outside path and a symlink are not supplied to the script. |
| Existing tests remain unchanged | The focused QA/prior-pass/implement selection passed, and the stable incremental rerun passed the entire Go suite, including daemon, CLI, spec and speccheck. No existing test expectation was edited. |

### Checks run

All Go commands below used `GOCACHE=/private/tmp/roundfix-task03-gocache`.

- `rtk proxy go test ./internal/daemon ./internal/cli -run 'QA|Prior|Implement' -count=1`
  — passed.
- `rtk proxy go test -race ./internal/daemon ./internal/cli ./internal/testfixture -run 'TestQAFormat|TestImplementPassesVerificationFormatToTheQAStep|TestNoTestWritesAnExecutableOutsideTheResidue' -count=1`
  — passed after the fixture correction, including all six required daemon
  tests, the CLI wiring test and the additional failure/cancellation tests.
- `GOOS=windows rtk proxy go build ./internal/daemon` and
  `GOOS=linux rtk proxy go build ./internal/daemon` — each passed; the host
  build and tests also exercise Darwin.
- `rtk make verify-incremental` — the initial sandboxed run exited 2: the new
  executable script writes violated the frozen fixture inventory, two
  existing force-stop tests lacked process-table access, and an edit to the
  new CLI test during that run triggered the repository fingerprint guard.
  The scripts were corrected to private files invoked through `sh`. The
  rerun with host process-table access and no concurrent source edits exited
  0: formatting, vet, the full Go suite, skill checks and binary build passed.
- `rtk proxy git -c core.fsmonitor=false diff --check` — passed.

The listed Task Verification commands remain for the Daemon. No follow-up
implementation outside Task 03 was included.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/daemon/qa_format_unix.go`
- `internal/daemon/qa_format_windows.go`

## Carry-forward provenance

- Source Run: `run_20261007T095139Z_d22b1ce63abf526a`
- Source commit: `3f0d8fdcdb260f8e122206ead6e1126aa05c4ce7`
