---
schema: roundfix/archive-record/v1
spec: 0216-baseline-wording-left-after-the-stack-wave
title: Baseline wording left after the stack wave
status: archived
created: "2026-10-02"
archived: "2026-10-02"
disposition: pass
source: docs/history/specs/0216-baseline-wording-left-after-the-stack-wave
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_04
qa_report: qa-report-2026-10-02.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0222
sources:
  - 2026-09-30-baseline-wording-left-after-the-stack-wave.md
  - 2026-09-30-a-dispatch-trigger-names-a-skill-the-model-cannot-invoke.md
regeneration:
  - command: make baseline-digests
promoted: []
pull_request: "349"
delivery_commit: 513b22ebab2d62a23b4222073c937c2862907ec2
---

# Baseline wording left after the stack wave

The stack wave (Specs 0206 to 0208) gave the Baseline's stack rules force but left sentences that are untrue for some repository that reads them. The backend and frontend guides say "the workspace" without naming it, although the built-in profiles declare its path. The clause that forbids generic `modules` or `services` buckets has no scope, so it reads as a ban on domain services. Nine dispatch triggers name a skill that only a person can start, so an Agent is told to activate something it cannot load. This refactor makes the three sentences true, records who owns the production-code and debugging activations, and states with evidence which items of the adopted entries were already settled on main. It changes Baseline assets, one render token and their tests; no command, flag or recorded decision changes.
