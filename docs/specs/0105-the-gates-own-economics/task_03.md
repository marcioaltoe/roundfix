---
status: pending
type: backend
---

# Task: The citation parser reads the written forms

The parser recognises a decision only as `ADR-NNNN`. Specs are written in this
repository's two languages, and two measured forms are refused: a conjunction as
a list separator, and a number without its prefix on an obligations line.
Replacing the conjunction with a comma cleared both errors without changing a
word of meaning, which is a correction round spent on punctuation.

## Work

- Accept a conjunction as a list separator in both repository languages, so a
  list written `0026, 0029 and 0031` is read as three citations rather than two.
- Accept a decision number without its prefix when the context is an obligations
  line, where the form is unambiguous. Nowhere else: ADR-0093 bounds this check
  to recognising written forms, never to inferring intent, and a bare number in
  prose is not a citation.
- Where a citation still cannot be recognised, the failure names the form that
  is recognised. Today it does not, so an author is told the citation is missing
  without being told what would count.
- Cover both measured forms, a bare number outside an obligations line asserting
  it is still not a citation, and the message naming the recognised form.

## References

- `_prd.md` → User Story 5, Core Feature 5
- `_techspec.md` → Build Order 3
- ADR-0093 checks consistency by citation rather than inference, which is the
  boundary this Task works inside

## Verification
- `grep -q "TestCitationAcceptsWrittenForms" internal/speccheck/citations_test.go || exit 1; grep -q "TestCitationFailureNamesTheRecognisedForm" internal/speccheck/citations_test.go || exit 1; go test -count=1 ./internal/speccheck`
