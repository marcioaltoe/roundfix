---
status: completed
type: backend
---

# Task: Write Terminal Row on Precondition Refusal

Implement gate refusal path with terminal row.

## Work
- Gate refusal → single terminal row
- Row: `| 0 | blocked | precondition |`
- Frontmatter: `rows_blocked_precondition: 1`
- Store check name and reason
- Set verdict: `fail`

## Verification
- `grep -q "rows_blocked_precondition" internal/spec/qa.go && grep -q "precondition" internal/spec/qa.go && go test -count=1 ./internal/spec 2>&1 | grep -q "^ok"`


## References

- User Story 1: Gate writes valid report
- Core Feature 1: Terminal Row Writing

## Result

A gate that refuses at a precondition check can now write its refusal as a
structurally valid QA Report instead of an empty matrix. `internal/spec` owns
the shape: `WritePreconditionRefusalReport` renders the whole report — verdict
`fail`, `rows_blocked_precondition: 1` beside the three existing typed counts at
zero, the refusing check and its reason in the frontmatter, and a Results table
carrying exactly one terminal row, `| 0 | blocked | precondition |`. The reader
side learned the same field: `QAReport` now carries
`RowsBlockedPrecondition`, validated as a non-negative integer like every other
typed count, and a report claiming `pass` while declaring a precondition-blocked
row is rejected — a gate that never built its matrix measured nothing that could
pass.

Two details are deliberate. The check and reason are collapsed onto one line and
YAML-quoted, so a multi-line refusal reason cannot break the frontmatter it is
written into; and a refusal whose cause the caller cannot name still writes,
degrading to a recorded placeholder rather than to no report at all — a refusal
that cannot be written is the deadlock this Spec exists to end. The prose that
justifies the row is a list, never a second markdown table, so it cannot be read
as further result rows.

### Rendered artifact

Emitted by `WritePreconditionRefusalReport` for check `spec check --strict`,
reason `SC-VOCABULARY-UNDOCUMENTED: term "Run Ledger" is undocumented`:

```markdown
---
verdict: fail
rows_blocked_precondition: 1
rows_blocked_environment: 0
rows_blocked_finding: 0
rows_blocked_declared: 0
precondition_check: "spec check --strict"
precondition_reason: "SC-VOCABULARY-UNDOCUMENTED: term \"Run Ledger\" is undocumented"
---

# QA Report

## Results

| # | Status | Provenance |
| - | --- | --- |
| 0 | blocked | precondition |

## Precondition refusal

- check: spec check --strict
- reason: SC-VOCABULARY-UNDOCUMENTED: term "Run Ledger" is undocumented

The gate stopped at this check before it built its QA matrix, so no requirement was measured and the row above records the refusal itself.
```

### Evidence per acceptance criterion

Red first: with the tests written and the contract absent,
`go test -run 'TestWritePreconditionRefusalReport|TestReadQAReportRejectsA' -count=1 ./internal/spec`
failed to build — `undefined: WritePreconditionRefusalReport`, `undefined:
PreconditionRefusal`, `undefined: QAPreconditionRowID`, and `QAReport has no
field or method RowsBlockedPrecondition`.

Focused check after the last edit — `go test -run
'TestWritePreconditionRefusalReport|TestReadQAReportRejectsA|TestQAVerdict|TestReadQAReport|TestNewestQAReport|TestArchivedQAReportCorpus'
-count=1 ./internal/spec` → `ok roundfix/internal/spec 0.668s`.

| Criterion (TechSpec Acceptance 1) | Evidence |
| --- | --- |
| Results table has one terminal row, not empty | `TestWritePreconditionRefusalReportWritesOneTerminalRow` reads the rows back out of the rendered `## Results` section and fails on any count other than one |
| Row status is `blocked`, provenance is `precondition` | the same test compares the row against `QAPreconditionRowID`/`QAPreconditionRowStatus`/`QAPreconditionRowProvenance` rather than a copied literal |
| Frontmatter has `rows_blocked_precondition: 1` | asserted in the same test, and read back through the package's own reader by `TestWritePreconditionRefusalReportIsReadableAsARefusal`, which gets verdict `fail` and `RowsBlockedPrecondition = 1` with every other cause zero |
| Precondition check name and reason recorded | `precondition_check` and `precondition_reason` asserted in the same test; `TestWritePreconditionRefusalReportKeepsTheRefusalOnOneLine` proves a reason carrying a newline, a pipe, and quotes still parses back through `ReadQAReport`; `TestWritePreconditionRefusalReportRecordsAnUnnamedRefusal` proves an unnamed refusal is still recorded |
| Verdict is `fail` | asserted at render and again after the round-trip; `TestReadQAReportRejectsAPreconditionBlockedPass` proves a precondition-blocked `pass` is unreadable, and `TestReadQAReportRejectsANegativePreconditionCount` holds the count to the same non-negative contract as its siblings |

Supporting checks, all after the last edit: `gofmt -l internal/spec/` → no
output; `go vet ./internal/spec` → clean; `go build ./...` → clean. `go vet
./...` reports only pre-existing `passes lock by value` findings in
`internal/agent`, untouched by this Task.

Changed paths: `internal/spec/qa.go`, `internal/spec/qa_test.go`, and this Task
file. No tooling configuration, derived artifact, or other Spec file was
touched.

### Follow-ups for later Tasks

- Nothing calls `WritePreconditionRefusalReport` yet. The daemon's
  `writeMechanicalQAReport` (`internal/daemon/task_engine.go:2172`) still writes
  its own `| QA-PRECONDITION | fail | mechanical refusal: … |` row; wiring the
  refusal path to this contract belongs to task_02.
- `detectMechanicalReportShape` (`internal/speccheck/mechanical.go:1353`) would
  today reject the new row as "a blocked cause outside environment, finding, or
  declared", and its `mechanicalBlockedCounts` does not track
  `rows_blocked_precondition`. Teaching the detector the status and provenance
  pair is task_04, which the graph already sequences after this Task.
- `precondition_check` and `precondition_reason` are written here but not yet
  read into `QAReport`; task_03 owns the read side and their preservation.
- Identifier strategy: the `QA Report` glossary entry in `CONTEXT.md` names only
  `rows_blocked_environment` and `rows_blocked_finding`, so it already omits
  `rows_blocked_declared` and now `rows_blocked_precondition`. The domain rule
  places that check at the close of the Spec, so it is left for the closing node
  rather than done piecemeal here.
