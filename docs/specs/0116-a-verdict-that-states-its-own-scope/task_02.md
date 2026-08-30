---
status: completed
type: backend
---

# Task: A QA Report names its auditor

The report records the tree it audited. It does not record the Roundfix that
produced the verdict, so a reader cannot attribute a mechanical finding to a
named auditor.

## Work

- Record the auditing binary and its staleness in the QA Report frontmatter,
  beside the existing audited-commit field and clearly distinct from it.
- The precondition-refusal writer records them too. A gate that refused before
  building its matrix still has an auditor, and its findings are read the same
  way.
- Keep the change backward compatible: a report written without the new keys
  must still validate, so reports already in the repository stay readable.
- Change no verdict rule, no row contract, and no blocked-cause count. The
  report states a condition; nothing branches on it.
- Cover: a released-build identity with empty commit and time produces a
  complete record rather than an empty field; a refusal report carries the
  auditor; an older report without the keys still validates.

## References

- `_prd.md` → Goal 3, User Story 3, Core Feature 3
- `_techspec.md` → Build Order 2; Data Models
- ADR-0080 keeps a report from crediting a row without evidence, which this
  Task does not touch

## Verification
- `grep -q "auditing_binary" internal/spec/qa.go && grep -q "auditor_staleness" internal/spec/qa.go && grep -q "TestPreconditionRefusalReportNamesItsAuditor" internal/spec/qa_test.go && go test -count=1 ./internal/spec`

## Result

The QA Report reader now exposes optional `auditing_binary` and
`auditor_staleness` strings. Their absence remains valid and reads back as
empty metadata, preserving reports written before this contract existed.

The precondition-refusal writer records the running `app.Auditor()` identity
and the explicit staleness answer returned by `CompareToTree`. A released
binary with empty commit and build-time stamps therefore records its version,
not an empty identity. Because this writer receives no audited-tree version or
ancestry signal, it records `unknown` with the missing-signal reason and never
infers `current`. The Daemon's refusal-report shape test now requires both
frontmatter keys.

Acceptance evidence:

- Auditor fields: `TestPreconditionRefusalReportNamesItsAuditor` reads the
  serialized refusal back through `ReadQAReport` and asserts both named fields.
- Released identity: the same test sets version `1.2.3` with empty commit and
  build time and observes `auditing_binary: "1.2.3"`.
- Refusal path: the same test invokes `WritePreconditionRefusalReport`; the
  Daemon's focused refusal-shape subtest observes the fields through the
  production report-materialization path.
- Backward compatibility:
  `TestReadQAReportAcceptsOlderReportWithoutAuditorMetadata` reads a report
  carrying neither key and preserves its `pass` verdict.
- Verdict, row, and blocked-cause contracts:
  `TestWritePreconditionRefusalReportWritesOneTerminalRow`,
  `TestWritePreconditionRefusalReportIsReadableAsARefusal`,
  `TestQAVerdictValidatesBlockedCounts`, and
  `TestReadQAReportBlockedCounts` retain the existing behavior.

Focused checks:

- Red starting signal: `rtk env GOCACHE=/private/tmp/roundfix-task02-gocache
  go test ./internal/spec -run
  '^(TestPreconditionRefusalReportNamesItsAuditor|TestReadQAReportAcceptsOlderReportWithoutAuditorMetadata)$'
  -count=1` failed to compile before implementation because `QAReport` had no
  auditor fields.
- The same focused command with `-v` passed after implementation.
- `rtk env GOCACHE=/private/tmp/roundfix-task02-gocache go test
  ./internal/daemon -run
  '^TestWriteMechanicalQAReportWritesThePreconditionRefusal$/^the_refusal_is_recorded_as_one_terminal_row_and_its_frontmatter$'
  -count=1 -v` passed after updating the intentional report-shape expectation.
- `rtk env GOCACHE=/private/tmp/roundfix-task02-gocache go test
  ./internal/spec -run
  '^(TestPreconditionRefusalReportNamesItsAuditor|TestReadQAReportAcceptsOlderReportWithoutAuditorMetadata|TestWritePreconditionRefusalReportWritesOneTerminalRow|TestWritePreconditionRefusalReportIsReadableAsARefusal|TestQAVerdictValidatesBlockedCounts|TestReadQAReportBlockedCounts)$'
  -count=1 -v` passed every named test and subtest.
- `rtk env GOCACHE=/private/tmp/roundfix-task02-gocache go vet
  ./internal/spec ./internal/daemon` passed.
- The first sandboxed `rtk make verify-incremental` run exposed the stale
  Daemon expectation and two unrelated process-table permission failures in
  `internal/cli`. After the expectation update, the permitted rerun passed all
  packages, skill checks, and the build.

The Task's declared Verification remains unrun for the Daemon.

## Carry-forward provenance

- Source Run: `run_20260830T161359Z_31aaee7e42ecc4e4`
- Source commit: `1374af5d694835ba0b3d4d1a537c1b3bfc037a78`
