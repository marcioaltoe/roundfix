---
status: accepted
created_at: 2026-11-04T12:00:00Z
updated_at: 2026-11-04T12:00:00Z
deprecated_at: null
superseded_by: null
---

# A precondition refusal is never imported as a prior QA pass

ADR-0194 imports the newest unintegrated QA Report commit so that a failed
pass becomes the next pass's previous report. ADR-0249 rejected not importing
a failed pass and formats the QA directory instead. On 2026-11-03 Pantheon,
which configures no Format Command, retried after a precondition refusal; the
import brought the refusal in, the precondition refused on its bytes, and the
Run recommitted it with no Agent session.

## Decision

A candidate whose added or modified reports record a Precondition Refusal is
refused, with no older candidate searched. A refusal measured no requirement,
so it is not a pass and has nothing to hand on. The refusal report keeps going
through the Format Command at its commit, and the refusal commit is never
carried forward.

## Considered Options

- Importing the refusal as today. Rejected: it can make the next precondition
  refuse on the imported bytes and recommit the same refusal without an Agent
  session.
- Searching an older candidate. Rejected: a refusal must start a pass with no
  previous report rather than selecting an earlier report whose status is no
  longer the newest candidate.
- Dropping only the refusal file from a mixed commit. Rejected: a candidate
  commit that records a refusal is not a valid prior pass, so importing part of
  it would mix unrelated report history.
- Carrying the refusal commit forward and superseding it. Rejected: Task
  Carry-Forward carries only Task settlement commits, and a QA Report commit
  never leaves its Run Branch.
- Adding a flag to discard a prior report. Rejected: the refusal is a property
  of the candidate report and must be excluded by the import contract, not by
  an operator choice that could vary between retries.

## Consequences

This narrows ADR-0194 without superseding it. A pre-fix mixed commit is
refused whole, and Pantheon's manual workaround is no longer needed.
