---
schema: roundfix/archive-record/v1
spec: 0095-a-verification-that-ran-before-anyone-believed-it
title: A Verification that ran before anyone believed it
status: archived
created: "2026-08-12"
archived: "2026-08-14"
disposition: pass
source: docs/history/specs/0095-a-verification-that-ran-before-anyone-believed-it
source_revision: 1d1726f5d227731aa229fbc336dccbae9aca1398
qa_task: task_09
qa_report: qa-report-2026-08-14-01.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "162"
delivery_commit: 60036423f2f15851ccc3f76a9e2ed755f69c2c27
---

# A Verification that ran before anyone believed it

A Task's Verification is checked for form and never for execution, and the gap is the most expensive one measured across the fleet. Six defects in one night in a repository this Spec did not build were pure shell semantics — a count that prints a filename, a pipeline that exits with the status of its last stage, a platform that pads a number with spaces — and four of them failed work that was correct. A twelve-Task graph elsewhere came within one runner's exit code of being entirely vacuous, saved by the tool's luck rather than by design. Meanwhile the rule that a Verification passes only by exiting zero is enforced by the Daemon and written nowhere, so a natural-looking assertion whose success is an empty result fails its Task for doing the right thing.
