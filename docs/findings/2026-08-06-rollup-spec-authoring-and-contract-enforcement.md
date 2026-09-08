---
status: pending
created_at: 2026-08-06
updated_at: 2026-09-08
kind: rollup
members:
  - 2026-08-06-a-promoted-backlog-entry-has-nowhere-valid-to-go.md
  - 2026-08-06-minting-an-adr-opens-gaps-no-one-can-ever-close.md
  - 2026-08-07-a-live-contract-lives-inside-an-archived-spec.md
  - 2026-08-04-a-finding-does-not-prevent-the-recurrence-it-describes.md
  - 2026-08-05-authoring-has-no-procedure-for-a-disproven-premise.md
  - 2026-08-06-authoring-rules-a-release-night-made-checkable.md
  - 2026-08-06-every-run-that-failed-tonight-failed-on-a-contract.md
---

# Spec authoring and contract enforcement — prose matters only when execution reads it (2026-08-06)

These findings track defects that were already described in a finding, ADR, or
Spec but recurred because no Task or check consumed the statement. Authoring is
complete only when each promise reaches executable work and each invalidated
premise has an explicit route back to the artifact that asserted it.

## Consolidated learning

- A finding records evidence; it does not prevent recurrence until its defect
  class reaches a Task, test, or consistency check.
- A PRD or TechSpec promise that no Task owns does not ship, and same-wave Tasks
  must be file-disjoint when the graph allows them to run together.
- Citation checks prove that an obligation was named, not obeyed. Behavioral
  obligations need evidence from the surface they govern.
- A disproven premise needs a supported authoring outcome instead of repeated
  gate cycles against a Spec that is no longer true.

## Live edge

Spec 0065 made several requirement contradictions checkable. The rollup remains
`pending` for the broader traceability contract from finding and ADR through
Task ownership to the evidence that proves the authored consequence shipped.

---

Triage 2026-08-26: deferred out of the active queue. See docs/workflow/2026-08-26-triage-the-queue-earns-its-tokens.md.

## Addendum — 2026-09-08 — Current triage

The promoted Backlog adoption seam was explicitly closed in its 2026-08-07
member addendum. Existing deliveries include consistency checks (0064, 0065),
stage-scoped citation checks (0093), hermetic/effect-proving Verification
(0095), wave collision detection (0097), and real-boundary characterization
(0105). Spec 0083 moved the coverage reference out of archived Spec content.

The upstream-only guide/glossary citation rule still has no matching detector,
as recorded in the open 2026-08-26 Backlog Entry. Traceability from each promise
to assigned work and behavior evidence is broader than citation checking.
This triage did not establish a complete supported path for a falsified premise
back through authoring and corrective work. Route those remaining contracts to
provisional P1, authorization and knowledge lifecycle, and P8, durable
unattended Spec workflow. Keep this Rollup active for its seven members; no
Spec number is assigned by these provisional labels.
