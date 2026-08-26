---
status: completed
type: backend
---

# Task: Store Precondition Metadata

Update QA Report for precondition metadata.

## Work
- Add: `precondition_check`, `precondition_reason`
- Fields optional (not required for passing)
- Preserved during read/write

## Verification
- `grep -q "CheckName\|checkName" internal/spec/qa.go && grep -q "Reason\|reason" internal/spec/qa.go && go test -count=1 ./internal/spec 2>&1 | grep -q "^ok"`


## References

- Core Feature 1: Terminal Row Writing

## Result

A refusal the gate wrote can now be read back out of the report it wrote it
into. `QAReport` carries the refusal itself — `Precondition PreconditionRefusal`
— read from `precondition_check` and `precondition_reason` in the frontmatter,
beside the typed counts task_01 added. The same type crosses both directions, so
a reader hands `WritePreconditionRefusalReport` exactly what that writer needs to
record the refusal again; nothing is summarized on the way out.

Three decisions are deliberate:

- **One type for both directions.** The read side returns
  `PreconditionRefusal`, not a second pair of flat fields, because the round trip
  is the property this Task owes: a metadata shape that only a writer can build
  would let the two ends drift.
- **The read is faithful, not repaired.** An absent check or reason stays empty
  rather than becoming the writer's `unnamed precondition check` placeholder, so
  a reader can tell a refusal whose cause the gate could not name — which is
  recorded, with that placeholder — from one that was never written down at all.
- **A refusal value is one line in both directions.** `qaRefusalLine` now backs
  both the writer's collapse and the reader's, so a hand-authored multi-line
  reason reads back as the single line the writer would have produced for it,
  and write → read → write stays byte-identical.

`PreconditionRefusal.Check` was renamed to `CheckName`. The field holds the name
of the check, not the check, which is the TechSpec's own wording (§2, "Check
name: the name of the Spec check that refused") and this Task's Verification
vocabulary; the one production caller
(`internal/speccheck/mechanical.go:148`) and four assertions in
`internal/speccheck/mechanical_test.go` moved with it.

Both fields stay optional. Nothing new refuses: a report that names no
precondition reads as one that refused for none, and a `pass` with no
precondition metadata is as readable as it was before this Task. The only value
that cannot be read is one that is not a scalar at all — a list or a mapping
under either key — which surfaces as the same `QAReportError` any other
malformed frontmatter value does.

### Evidence per Work item

Red first: with the tests written and the contract absent, `go test -run
'TestReadQAReportRecordsThePreconditionRefusal|TestPreconditionRefusalRoundTripsThroughTheQAReport|TestReadQAReportLeavesThePreconditionUnrecordedWhenNoneRefused|TestReadQAReportRejectsANonScalarPrecondition'
-count=1 ./internal/spec` failed to build — `unknown field CheckName in struct
literal of type PreconditionRefusal` (×5) and `report.Precondition undefined
(type QAReport has no field or method Precondition)`.

Focused check after the last edit — the same `-run` selection with `-v` → `ok
roundfix/internal/spec 0.455s`, all four tests and their nine subtests passing.

| Work item | Evidence |
| --- | --- |
| Add `precondition_check`, `precondition_reason` | `TestReadQAReportRecordsThePreconditionRefusal` writes a report through `WritePreconditionRefusalReport`, reads it back through `ReadQAReport`, and compares `report.Precondition` against the refusal that was written: a named refusal reads back verbatim, a check and reason spread over lines and tabs read back on one, and an unnamed refusal reads back as the placeholder pair the writer recorded |
| Fields optional (not required for passing) | `TestReadQAReportLeavesThePreconditionUnrecordedWhenNoneRefused` reads a `pass` report that names no precondition — no error, zero `Precondition` — and proves each field independently optional: a report carrying only `precondition_check` records the check alone, only `precondition_reason` records the reason alone. `TestArchivedQAReportCorpusRemainsReadable` re-reads every archived Spec's newest report, none of which carries either field |
| Preserved during read/write | `TestPreconditionRefusalRoundTripsThroughTheQAReport` writes a two-code refusal, reads it back, requires `report.Precondition` to equal the refusal that was written, then writes a second report from the value that was read and requires it byte-identical to the first — a read that dropped or reshaped either field cannot produce the same bytes |
| Unreadable metadata stays unreadable | `TestReadQAReportRejectsANonScalarPrecondition` requires `QAReportError` naming the unreadable value for a `precondition_check` written as a list and a `precondition_reason` written as a mapping |

Supporting checks, all after the last edit: `gofmt -l internal/spec/
internal/speccheck/` → no output; `go build ./...` → clean; `go vet
./internal/spec ./internal/speccheck` → clean; `go test -count=1 ./internal/spec
./internal/speccheck` → `ok`, `ok`; `go test -count=1 -tags repocontract
./internal/speccheck` → `ok`. Dependents of the renamed field and the widened
`QAReport`: `go test -count=1 ./internal/daemon ./internal/cli` → `ok`, `ok`.

Changed paths: `internal/spec/qa.go`, `internal/spec/qa_test.go`,
`internal/speccheck/mechanical.go`, `internal/speccheck/mechanical_test.go`, and
this Task file. No tooling configuration, derived artifact, or other Spec file
was touched.

### Follow-ups for later Tasks

- Nothing reads `QAReport.Precondition` yet. The Daemon still writes its own
  `| QA-PRECONDITION | fail | mechanical refusal: … |` row in
  `writeMechanicalQAReport` (`internal/daemon/task_engine.go:2172`) rather than
  calling `spec.WritePreconditionRefusalReport`, so no production path yet stores
  or recovers this metadata end to end. task_02 already recorded that the wiring
  has no Task claiming it in `_tasks.md`; this Task closes the read half of the
  contract that wiring will use, and does not change who calls it.
- `detectMechanicalReportShape` (`internal/speccheck/mechanical.go:1419`) still
  rejects the `blocked`/`precondition` row, and its `mechanicalBlockedCounts`
  does not track `rows_blocked_precondition`. That is task_04, next in the graph.
- Identifier strategy: this Task introduced no glossary term. It renamed one Go
  field inside `internal/spec` and added no vocabulary; the `QA Report` entry in
  `CONTEXT.md` still names only `rows_blocked_environment` and
  `rows_blocked_finding`, so whether it should also name
  `rows_blocked_precondition`, `precondition_check`, and `precondition_reason`
  remains the closing node's check, as task_01 recorded.
