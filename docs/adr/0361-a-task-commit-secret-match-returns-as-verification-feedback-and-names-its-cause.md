---
status: accepted
created_at: 2026-11-05T00:00:00Z
updated_at: 2026-11-05T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A Task commit's Secret Match returns as Verification Feedback and names its cause

ADR-0292 scanned after the last attempt and deferred a repair turn; on
2026-11-04 Fiscus lost a Run to a `credential-assignment` fixture that no
report, park line or Agent turn named.

The scan is the last Settlement Check of every attempt, with or without a QA
gate. A first-attempt match returns as Verification Feedback with locations and
a `Fix:` line, never the value; the second attempt rescans; a final match keeps
the reason; and a scan error halts the Run. The Causes report lists a
`secret-scan` item classed `repository-convention` by kind, and a
`run-unresolved` park names failed Tasks, reasons and skipped Tasks.

## Considered Options

- A second repair turn only for the scan — rejected; the existing one repair
  turn handles all failed checks in the attempt.
- Scanning only on attempt 1 — rejected; the final attempt must prove that the
  Task commit no longer carries a Secret Match.
- Keeping the post-attempt scan as a second pass — rejected; the scan is a
  Settlement Check and must return its first-attempt failure as Verification
  Feedback.
- Accepting a value tracked files already hold — rejected; the Task commit is
  still an outbound boundary.
- A new signature in the cause table — rejected; the `secret-scan` item is
  classed by kind.
- A new Park Class or blocker name — rejected; `run-unresolved` carries the
  failed and skipped Task detail.
- A structured blocker column — rejected; the existing blocker line remains
  the durable user-facing contract.

## Consequences

This amends ADR-0292 by taking its rejected repair turn and amends ADR-0214 by
classing one item kind by kind. ADR-0273's value rule extends to Verification
Feedback, the Causes report and the park. Each attempt prepares the Task commit
once more, including one worktree snapshot and one `git diff` of the prepared
paths, and a Task that fails for another reason is scanned too so its repair
turn receives every problem at once.
