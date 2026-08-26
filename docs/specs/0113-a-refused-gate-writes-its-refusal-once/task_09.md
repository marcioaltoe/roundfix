---
status: pending
type: docs
---

# Task: Document the vocabulary this Spec emitted

QA finding F-2 (`Trust-Damage`): the PRD declares `surfaces: [backend, docs]`
and makes the identifier strategy applicable, assigning the closing node the
check of whether the work introduced or changed a term. The check ran and found
three tokens entering the QA Report artifact with no documentation anywhere:
`rows_blocked_precondition`, `precondition_check`, and `precondition_reason`.
The whole documentation surface received nothing.

## Work

- `CONTEXT.md`: the **QA Report** entry names `rows_blocked_precondition` beside
  the counts it already carries, and the two metadata keys. Add the precondition
  refusal itself as a term if the glossary has no home for it, following the
  shape of the entries already there.
- `docs/user-guide/context-driven-development.md`: the report contract states
  what a precondition refusal leaves behind, and that the mechanical stage reads
  the newest report only. Both are behavior a maintainer whose gate refused will
  look for, and neither is written anywhere today.
- `_techspec.md`: add the `## Vocabulary Contract` section whose absence made
  the checker skip its vocabulary detector rather than decide this. The skip is
  what let three undocumented tokens reach a passing static gate.
- Claims are read from the delivered code, not from the TechSpec's draft. Where
  the two disagree, the code is the fact and the TechSpec is corrected: section
  5 of the TechSpec states the shape detector reads any markdown table, which
  task_07 measured false — it bound the scan to the `## Results` section, and
  the blockers came from subsections inside that window.

## References

- User Story 3: Maintainer reads what caused the refusal
- Core Feature 2: Mechanical stage reads the newest report
- Core Feature 4: Only the Results table is read as results

## Verification
- `grep -q "rows_blocked_precondition" CONTEXT.md && grep -q "precondition_check" CONTEXT.md && grep -qi "precondition" docs/user-guide/context-driven-development.md && grep -q "newest" docs/user-guide/context-driven-development.md && grep -q "## Vocabulary Contract" docs/specs/0113-a-refused-gate-writes-its-refusal-once/_techspec.md`

## Result
Every token this Spec emitted is named where the repository names its
vocabulary, and the TechSpec agrees with the code that was delivered.
