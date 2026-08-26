---
status: pending
type: docs
---

# Task: Document the term and the two bounds

This Spec coins **Run Window** and adds a second time bound beside
`budget.max_run_duration`. ADR-0137 requires the Run budget to be explained
where it is configured; a reader who meets either bound must learn how the two
relate without reading source.

## Work

- `CONTEXT.md`: add **Run Window** in the shape of its neighbours — the window
  during which Runs may be created, repository-scoped, durable, governing the
  start and never the finish. Record what to avoid: `Session` is taken by
  **Agent Session**, and `curfew` and `deadline` invite the wrong reading.
- The configuration surface that prints `budget.max_run_duration` gains the
  converse pointer: that key bounds how long a Run may run, the Run Window
  bounds when one may start.
- The user guide documents the command, the refusal, and the crossing report.
- Claims are read from the delivered code, not from the TechSpec draft. Where
  the two disagree, the code is the fact and the TechSpec is corrected.

## References

- User Story 2: The refusal names the cutoff and how to act
- Core Feature 5: The two time bounds explain each other

## Verification
- `grep -q "Run Window" CONTEXT.md && grep -q "Agent Session" CONTEXT.md && grep -qi "run window" docs/user-guide/commands.md && go test -count=1 -tags docscontract ./internal/docscontract 2>&1 | grep -q "^ok"`

## Result
The coined term is in the glossary and both time bounds explain each other
where a reader meets them.
