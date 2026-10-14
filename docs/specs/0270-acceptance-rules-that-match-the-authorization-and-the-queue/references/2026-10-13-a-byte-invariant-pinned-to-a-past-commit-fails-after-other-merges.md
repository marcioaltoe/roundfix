---
type: fix
status: promoted
created: 2026-10-13
spec: 0270-acceptance-rules-that-match-the-authorization-and-the-queue
reason: null
---

# A byte invariant pinned to a past commit fails once other Specs merge first

## Symptom

task_03 of 0264 required `docs/references/behavior-surfaces.json` to be
byte-identical to its bytes at `3820bdd4`, the commit where the plan measured
it. Four Specs merged before 0264 ran and changed the file, so the QA failed a
row although the Spec's own diff left the file unchanged against its delivery
base.

## Where

The `write-tasks` QA requirements and `roundfix spec check`, which accept a
byte invariant pinned to any commit.

## Expected

A byte invariant names the delivery base (the merge base of the audited head
with the default branch), or `spec check` reports a QA requirement that pins
bytes to a commit older than the Spec's plan merge.

## Evidence

Delivery Queue item 0264 on 2026-10-13, QA finding F01 (row 8); fixed in the
item by 4af8655e.
