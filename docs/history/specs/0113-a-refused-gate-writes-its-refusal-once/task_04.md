---
status: completed
type: backend
---

# Task: Update Mechanical Stage Validation

Modify mechanical stage to accept terminal row.

## Work
- Empty table → refuse SC-REPORT-SHAPE
- Terminal blocked row → accept
- Validate status and provenance

## Verification
- `grep -q "provenance.*precondition\|ProvenancePrecondition" internal/speccheck/mechanical.go && go test -count=1 ./internal/speccheck 2>&1 | grep -q "^ok"`


## References

- Core Feature 2: Mechanical Stage Update

## Result

The mechanical stage now accepts the report a refused gate writes. Before this
Task, `detectMechanicalReportShape` read the row `| 0 | blocked | precondition |`
as a bare `blocked` no typed cause covered and raised `QA-REPORT-SHAPE` — "row 0
has a blocked cause outside environment, finding, or declared" — so the artifact
task_01 taught the gate to write was still a blocker on the next run. It is not
one now: the whole report `spec.WritePreconditionRefusalReport` produces passes
the stage with zero findings and `Blocking = false`.

Three decisions are deliberate:

- **Status and provenance are read together.** The refusal is recognised by the
  pair, not by either cell alone. `blocked` alone still says only that nothing
  was measured — it keeps its existing refusal, because the three typed causes
  remain the only way to say *why* a measurement blocked. Provenance alone would
  let any status claim to be a refusal, so a row with provenance `precondition`
  and a status other than `blocked` is refused by name: `row 0 has provenance
  precondition with status "pending"`. Both cells come from the constants
  task_01 defined (`spec.QAPreconditionRowStatus`,
  `spec.QAPreconditionRowProvenance`), so the writer and this detector cannot
  drift apart.
- **`rows_blocked_precondition` is required only of a report that records a
  refusal.** The other three counts are required of every closed report. Adding
  a fourth to that list would have raised a fresh `QA-REPORT-SHAPE` on every
  report already in the tree — none of them declares the field, and all of them
  predate it — which is the same deadlock, moved. `detectPreconditionCount`
  therefore refuses only two states: a declared count that disagrees with the
  refusal rows present, and a refusal row whose count was never declared.
- **The Provenance column is optional.** A gate's ordinary Results table
  (`# | Story / criterion / sweep | Actor and surface | Status | Evidence`)
  declares no Provenance column, so `row.provenance` is empty there and nothing
  in this Task can fire on it. The Daemon's own mechanical Results table does
  carry the column, and its refusal row's provenance is `mechanical refusal: …`
  (`internal/daemon/task_engine.go:2184`), which is not `precondition` — so that
  path is unchanged too.

The empty-table refusal is untouched and now has a test that pins it: a Results
table with a header and separator but no rows still yields exactly one
`QA-REPORT-SHAPE`, "Results table has no report rows". A refusal report is valid
because it carries a row, not because the check was relaxed.

`mechanicalBlockedCounts` now tracks `rows_blocked_precondition` alongside the
other three, so a non-integer value under that key surfaces as the same
frontmatter parse refusal any other malformed count does. One small extraction:
the count-mismatch finding both the existing loop and the new precondition check
raise is now built in `addMechanicalCountMismatch`, so the two cannot word the
same defect differently.

### Evidence per Work item

Red first: with `internal/speccheck/mechanical.go` stashed and the new tests in
place, `go test -count=1 -run TestMechanicalStageAcceptsThePreconditionRefusalRow
./internal/speccheck` failed on five subtests — `the refusal a gate writes passes
the stage that reads it`, `blocked beside precondition provenance is terminal`,
`precondition provenance with a non-blocked status refuses`, `a declared count
that no refusal row matches refuses`, and `a refusal row without its declared
count refuses`. The three that passed without the change are the regression
guards, which is what they are for.

Focused check after the last edit: `go test -count=1 -v -run
TestMechanicalStageAcceptsThePreconditionRefusalRow ./internal/speccheck` → `ok
roundfix/internal/speccheck`, all eight subtests running and passing.

| Work item | Evidence |
| --- | --- |
| Empty table → refuse SC-REPORT-SHAPE | Subtest `an empty Results table still refuses` writes a report whose Results table has a header and separator and no rows, and requires exactly one `QA-REPORT-SHAPE` whose detail is `Results table has no report rows`. `TestMechanicalFindingsWithoutRowHintsBlockTheirRefusalCode` and the `report-red.md` fixture still hold their existing counts |
| Terminal blocked row → accept | Subtest `the refusal a gate writes passes the stage that reads it` writes the exact bytes of `spec.WritePreconditionRefusalReport` to a report path and requires no `QA-REPORT-SHAPE` finding and `Blocking = false` — the artifact task_01 writes is read by the stage that used to refuse it. Subtest `blocked beside precondition provenance is terminal` pins the same acceptance on a hand-written row |
| Validate status and provenance | Subtest `precondition provenance with a non-blocked status refuses` requires one finding naming both cells for a row whose status is `pending` beside provenance `precondition`; `bare blocked without precondition provenance still refuses` requires the pre-existing `blocked cause outside environment, finding, or declared` for a row whose status is `blocked` beside provenance `measured`, so provenance is load-bearing in both directions. `a declared count that no refusal row matches refuses` and `a refusal row without its declared count refuses` cover the frontmatter half; `a report that records no refusal need not declare the count` proves the field stays optional for every other report |

Supporting checks, all after the last edit: `make fmt-check` → no output; `go
build ./...` → clean; `go vet ./internal/speccheck` → clean; `go test -count=1
./internal/speccheck ./internal/spec` → `ok`, `ok`. Verification-shaped probe:
`grep -q "provenance.*precondition\|ProvenancePrecondition"
internal/speccheck/mechanical.go` → match.

Whole suite: `go test -count=1 ./...` → every package `ok` except one flake in
`internal/agent` (`acpx_runner_test.go:2919`, the fake `acpx` fixture exiting
before its prompt-start event under full-suite load). That package is untouched
by this Task and passes on its own: `go test -count=1 ./internal/agent` → `ok`,
twice.

Changed paths: `internal/speccheck/mechanical.go`,
`internal/speccheck/mechanical_test.go`, and this Task file. No tooling
configuration, derived artifact, or other Spec file was touched.

### Follow-ups for later Tasks

- The stage still reads one report path chosen by its caller. The Daemon passes
  `previousReportPath` (`internal/daemon/task_engine.go:2073`), so a superseded
  report is still the one validated. That is task_05.
- `parseMechanicalReport` already bounds row collection to the table under
  `## Results` (`internal/speccheck/mechanical.go:895-940`, unchanged since
  2026-08-11), which is the property Core Feature 4 asks for. The provenance
  column this Task added is read inside that same bound. task_07 owns confirming
  whether the 2026-08-26 measurement on Spec 0098 came from a different reader
  or from an older installed binary, and pinning the property either way.
- Still nothing calls `spec.WritePreconditionRefusalReport` in production: the
  Daemon writes its own `QA-PRECONDITION` row rather than the refusal report
  whose shape this Task now accepts, as task_02 and task_03 both recorded. No
  Task in `_tasks.md` claims that wiring.
- Identifier strategy: this Task introduced no glossary term. It reads the
  existing `precondition` provenance value and the `rows_blocked_precondition`
  count that task_01 minted. Whether `CONTEXT.md`'s `QA Report` entry — which
  still names only `rows_blocked_environment` and `rows_blocked_finding` —
  should name the precondition count and its two metadata keys remains the
  closing node's check, as task_01 and task_03 recorded.
