---
status: pending
type: docs
---

# Task: Document both command contracts

A reader meets these two behaviors at the command line. The reconcile page
still describes carry-forward as a stopped-Run act, and the implement page
knows nothing about a Preflight that can refuse a Run over work a prior Run
proved.

## Work

- The reconcile contract records which Run outcomes carry-forward accepts and
  that every other terminal outcome is refused by name.
- The implement contract describes the new Preflight refusal: what triggers it,
  that no Run is created, the command that clears it, and that a Run whose
  stranded work is not carriable proceeds with a report instead.
- Say plainly that the caller effectively gets one carry-forward, because a
  carried Task's own file becomes a moved input afterwards. A reader who
  carries from the wrong Run first should learn why before doing it, not after.
- Recut after QA finding F-02. The first cut of this Task documented the
  decision rule as originally specified, and task_08 then changed it. Document
  the delivered rule: Preflight refuses only when the **complete** candidate
  set would carry, and a candidate's declared inputs are compared against the
  accumulating staged carries rather than against the raw checkout.
- Remove the superseded claims: that one passing Task triggers a refusal, and
  that inputs are compared byte-for-byte with the checkout.
- Claims are read from the delivered code, not from the TechSpec draft. Where
  the two disagree, the code is the fact and the TechSpec is corrected.

## References

- `_prd.md` → User Stories 1, 2 and 4; Core Features 2, 3, 4 and 5
- `_techspec.md` → Build Order 5; API Contracts

## Verification
- `grep -q "staged carries" docs/user-guide/commands.md && ! grep -q "byte-for-byte with the checkout" docs/user-guide/commands.md && ! grep -q "passes Task Carry-Forward.s proofs, Preflight Validation refuses" docs/user-guide/commands.md && grep -q "Stopped or Unresolved" docs/user-guide/commands.md && go test -count=1 -tags docscontract ./internal/docscontract`
