---
type: fix
status: promoted
created: 2026-10-16
spec: 0276-the-history-cleanup-reaches-every-adopter-and-the-review-sees-a-large-spec
reason: null
---

# Review convention C2 asks a QA Task for settlement sections

## Symptom

The pre-PR review raised a finding that a completed QA Task has no
`## Result`, `## Recorded paths` or `## Carry-forward provenance` section, and
its validation kept the finding. The Daemon writes those sections only for
implementation Tasks; a QA Task settles on its committed QA Report.

## Where

Delivery convention C2 in `internal/cli/review_conventions.go`, which states
that the Daemon writes those sections and does not exempt the QA Task.

## Expected

C2 says that a QA Task settles on its QA Report and carries no settlement
sections, so the reviewer does not raise the finding.

## Evidence

Pre-PR review of Spec 0273 on 2026-10-17, finding F1 on `task_05.md` line 5,
dismissed by evidence. Spec 0271's delivery (pull request #542, commit
c2dcb593a8) shows the same QA Task shape and merged.
