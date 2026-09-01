---
status: completed
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

## Result

Implementation:

- The Active ADR obligations row now reads four-digit decision numbers with or
  without an `ADR-` prefix. The exact-row boundary keeps bare-number recognition
  out of prose, while token-based parsing accepts comma lists and both `and` and
  `e` conjunctions.
- An unlisted-citation finding now tells the author that `ADR-NNNN` is the
  recognised form.

Acceptance evidence:

- `TestCitationAcceptsWrittenForms/English_conjunction_separates_decision_numbers`
  and `Portuguese_conjunction_separates_decision_numbers` exercise `0026, 0029
  and 0031` and `0026, 0029 e 0031` as three listed decisions.
- `TestCitationAcceptsWrittenForms/bare_decision_number_outside_obligations_is_not_a_citation`
  keeps a bare prose number outside the citation graph.
- `TestCitationFailureNamesTheRecognisedForm` asserts that an unrecognised
  obligations-row form reports `ADR-NNNN` in its fix.

Focused checks:

- Before the production change,
  `GOCACHE=/tmp/roundfix-task03-gocache go test -count=1 ./internal/speccheck -run 'TestCitation(AcceptsWrittenForms|FailureNamesTheRecognisedForm)$'`
  failed with three `SC-ADR-UNLISTED` findings for each conjunction case and a
  fix that omitted `ADR-NNNN`.
- After the production change, the same focused command passed.
- A final focused regression run,
  `GOCACHE=/tmp/roundfix-task03-gocache go test -count=1 ./internal/speccheck -run 'Test(CheckADRUnlisted|CitationAcceptsWrittenForms|CitationFailureNamesTheRecognisedForm)$'`,
  passed, including the existing prefixed-citation behavior.
- The Daemon-owned `## Verification` command was not run during implementation.

## Carry-forward provenance

- Source Run: `run_20260831T171256Z_df3674a688059467`
- Source commit: `5ba017a2cc9f7324b7f2d01de34aed973d6c36df`
