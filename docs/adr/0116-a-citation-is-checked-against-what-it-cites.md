---
status: accepted
created_at: 2026-08-09T00:00:00Z
updated_at: 2026-08-09T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A citation is checked against what it cites

Spec 0090's PRD and TechSpec both stated that ADR-0083 makes `make verify` the
authoritative gate. ADR-0083 is "Adopted sources move to their owning Spec" and
contains no verification decision. The rule the Spec described is real, but it
lives in `docs/agents/specific-repository.md`, and the Spec assigned it to a
decision record that never made it.

The repository's existing citation checks could not catch this. `SC-ADR-UNLISTED`
verifies that a cited ADR appears in the Spec's Active ADR obligations;
`SC-ADR-RELATED` verifies that an ADR citing a listed one is itself accounted
for. Both ask whether an obligation was *named*. Neither asks whether it was
*obeyed*, which is exactly what a finding filed on 2026-08-06 said would happen:
"Citation checks prove that an obligation was named, not obeyed."

Roundfix therefore reads the cited record when an artifact makes a claim about
what that record establishes, and reports a finding when the claim is not
supported by the record's text. The check is cheap by construction — both
documents are files in this repository, and the existing checker reads one Spec
in 0.04 seconds.

The finding names both texts rather than asserting intent. A semantic check's
failure mode is a false accusation, and a maintainer settles that in seconds
when the claim and the cited passage sit side by side; a check that only says
"mismatch" would spend more of their attention than the defect does.

Rejecting the check because natural language is hard was considered and refused.
The defect it catches is not subtle disagreement about meaning: it is a Spec
attributing a subject to a record about a different subject, which surface
matching resolves. Cases where the claim is arguable are reported for a human,
not decided by the checker.
