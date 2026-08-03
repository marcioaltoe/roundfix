---
status: accepted
created_at: 2026-08-03T00:00:00Z
updated_at: 2026-08-03T00:00:00Z
deprecated_at: null
superseded_by: null
---

# The consistency check is artifact-presence-aware

The Spec Consistency Check runs at every point in the authoring pipeline, from
a folder holding one PRD to a folder holding a full Task Graph and QA evidence.
It therefore never requires an artifact to exist. Each detector declares the
artifacts it reads and is skipped, not failed, when one of them is absent:
Coverage Map completeness is silent without a TechSpec, task-reference coverage
is silent without a Task Graph, and a vocabulary contract is checked only where
a Spec declares one. Presence is the pipeline's business — `spec-routing`
decides which artifacts a change needs — and agreement is this check's.

The alternative, demanding the full artifact set, was rejected on evidence: at
the moment this decision was taken, nine of ten active Specs carried a PRD
alone, all of them legitimately, because a TechSpec is authored when the Spec
reaches the front of the queue. A check that reported nine errors on the first
run would be switched off on the first day, and the check is only worth
building if it can be wired into the local gate as fail-closed.

This is the same non-regression clause the Spec's PRD states: a Spec that
passes today's gates must not be blocked by a new check that merely disagrees
with its style. Skipping is recorded, not silent — the report names each
detector it skipped and the artifact that was missing, so a reader can tell an
absent artifact from a passing one.
