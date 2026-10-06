---
status: accepted
created_at: 2026-10-06T00:00:00Z
updated_at: 2026-10-06T00:00:00Z
deprecated_at: null
superseded_by: null
---

# One QA partial policy, and rows a Run sandbox cannot reach

Daemon settlement, `roundfix settle`, `roundfix archive` and `roundfix
qa-report accept` already call one function, `spec.QAReportEligibility`. They
still disagreed in practice, because that function refused what ADR-0167
meant to allow. A `partial` whose only unmet row was the pre-PR Pull Request
row was refused with "expected pass", since the code also required at least
one declared row. A Pull Request row whose provenance carried a note, such as
`Pull Request row (Requirement 10 constrains execution)`, was not recognized
at all. Outside-evidence rows that the Run sandbox could not reach, because
it denied network access, counted as ordinary environment rows. From
2026-10-04 to 2026-10-06 the operator archived nearly every Spec with the QA
Archive Override for these rows alone, and two adopters reported the same
refusal.

On 2026-10-06 the maintainer decided, for rows a Run sandbox cannot reach:
"Igual à linha do PR". An outside-evidence row that is blocked only because
the sandbox denied network access is treated like the Pull Request row: it
never decides a qualifying partial, the report records that it was not
reached, and whoever needs that proof declares Unreachable Acceptance.

The policy has one rule, in `internal/spec`, which every caller reads:

- A `partial` qualifies when it records no finding-blocked row and no
  `skipped` row, every environment-blocked row is a pre-PR Pull Request row or
  a network-denied outside-evidence row, the declared-blocked rows are covered
  by the Spec's `## Unreachable Acceptance` declarations, and at least one row
  of these kinds is unmet. A `partial` that records none of them, for example
  one that only skipped rows, still refuses.
- A pre-PR Pull Request row has the status
  `blocked (environment: no open Pull Request)`, and one provenance item that is
  `Pull Request row` or starts with it followed by a space, `:` or `(`.
- A network-denied outside-evidence row has the status
  `blocked (environment: network denied: <host>)` with a non-empty host, and
  one provenance item that is `outside-evidence row` or starts with it
  followed by a space, `:` or `(`. Any other environment row, including a
  network-denied row with no outside-evidence provenance, still needs
  declaring or the override.
- The Delivery Queue parks `qa-environment-partial` only when the newest
  report has an environment-blocked row outside those two kinds. Only such a
  row is one the QA Archive Override exists for.

The QA gate's verdict rule does not change. A Pull Request row with
equivalent evidence for every control still lets the report reach `pass`.
Without that evidence, the report closes `partial`, and that partial now
qualifies. We rejected requiring the gate to write `pass` whenever only the
Pull Request row is blocked: a `pass` must not claim evidence the gate did
not record, and a reader then cannot tell a `pass` with equivalent evidence
from one without it. The rule stays one rule because the verdict follows the
evidence and eligibility follows the row kinds. Neither guesses the other.

Trusting a network-denied row only when the Daemon observed the denial was
also rejected. The Daemon does not observe the agent's network calls, and
the row names its host, so a reader can repeat the lookup outside the
sandbox.

## Consequences

- ADR-0167 is extended. A partial whose only unmet row is the pre-PR Pull
  Request row qualifies without a declaration, and a provenance note after
  the source name no longer hides the row.
- ADR-0104 is superseded in part. A network-denied outside-evidence row no
  longer holds Pull Request preparation. Every other blocked outside-evidence
  row still holds it until the row is satisfied or carried forward.
- ADR-0229 is superseded in part. A partial blocked only by Pull Request rows
  or network-denied outside-evidence rows settles `completed` and never parks
  `qa-environment-partial`.
- Archived reports keep reading as before. The frontmatter counts keep their
  meaning, and each exempt row is derived from the Results rows that the
  mechanical stage already validates.
- The check rests on the text the gate writes. A gate that labels a row it
  could have reached as network-denied hides that row from the override.
  The named host and the archived report make the claim checkable after
  the fact.
