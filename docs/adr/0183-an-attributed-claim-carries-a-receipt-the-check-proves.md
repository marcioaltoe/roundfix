---
status: accepted
created_at: 2026-09-30T00:00:00Z
updated_at: 2026-09-30T00:00:00Z
deprecated_at: null
superseded_by: null
---

# An attributed claim carries a receipt the check proves

`SC-CITATION-UNSUPPORTED` compares the words of a claim with the words of the
decision record it names (ADR-0116). That catches a claim assigned to a record
about another subject. It cannot catch a claim the record never makes when the
claim happens to reuse the record's vocabulary. On 2026-09-29 and 2026-09-30
the planning changes of Waves 4, 5 and 6 needed two to seven review rounds
each. Many findings were attributions that a reviewer could settle only by
opening the cited text.

A claim that attributes behavior to a decision record now carries a Claim
Receipt in the same paragraph: the source, a colon and a verbatim quote, as in
`ADR-0116: "reads the cited record"`. The Spec Consistency Check proves that
the quote is a contiguous part of the source once whitespace is normalized. A
written receipt whose quote is absent from its source is an error. An
attribution with no receipt is a gap. A receipt may also quote a repository
file, and it is proved the same way.

The check proves presence, not support. Whether the quoted words establish the
claim stays with the reader, who now has both texts side by side. A reviewer,
or a later advisory judgment, starts from a quote that exists instead of a
paraphrase.

The missing-receipt gap binds a Spec from the commit that added the
concrete-contract guide to the repository's write-techspec skill. A Spec whose
PRD was committed before that guide was written under the earlier contract and
is not held to it. A repository without the guide has not adopted the rule, and
the check says so in a skip line. An uncommitted PRD and an unreadable history
keep the full check, the fail-closed direction ADR-0168 already takes.

## Consequences

- The attribution grammar does not change. A receipt is required for exactly
  the claims ADR-0116's check already reads.
- A written receipt is proved whatever the Spec's age, because writing one
  asserts that the quote exists.
- A planning branch cut before the guide landed and squash-merged after it is
  held once merged. Its author rebases and adds the receipts before merging.
- ADR-0093 still holds: the check compares written text and infers nothing.
