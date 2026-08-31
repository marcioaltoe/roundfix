---
status: pending
type: backend
---

# Task: Roundfix owns the QA Task's Verification

The gate Task's Verification is authored by hand over a verdict the Daemon
already derived. In one measured case that produced a gate which passed itself
having failed — a hand-written predicate accepted a verdict outside the domain,
and it read as ordinary Task authoring to every reviewer.

## Work

- For a Task of type `qa`, Roundfix supplies the Verification instead of reading
  the author's. The derived command requires the newest QA Report to record a
  passing verdict, so a predicate can no longer accept a verdict outside the
  domain or select an older report.
- Render the derived command into the Task file rather than leaving it implicit.
  Removing the author's control must not also remove the reader's view of what
  will run.
- A `qa` Task that carries its own authored command is refused by name through a
  Spec check finding. Refuse rather than overwrite: silently replacing an
  author's text is how a contract becomes invisible.
- Change nothing for any other Task Type. This is bounded to the one node
  ADR-0091 makes terminal.

## References

- `_prd.md` → Goal 3, User Story 4, Core Feature 3
- `_techspec.md` → Build Order 1; Interfaces: `DerivedQAVerification`,
  `AuthoredQAVerification`
- ADR-0091 makes the gate one terminal Task node of type `qa`

## Verification
- `grep -q "DerivedQAVerification" internal/spec/task.go || exit 1; grep -q "AuthoredQAVerification" internal/spec/task.go || exit 1; grep -q "TestDerivedQAVerification" internal/spec/task_test.go || exit 1; go test -count=1 ./internal/spec ./internal/speccheck`
