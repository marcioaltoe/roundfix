---
status: pending
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
