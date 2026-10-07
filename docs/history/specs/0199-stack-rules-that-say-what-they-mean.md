---
schema: roundfix/archive-record/v1
spec: 0199-stack-rules-that-say-what-they-mean
title: Stack rules that say what they mean
status: archived
created: "2026-09-30"
archived: "2026-09-30"
disposition: pass
source: docs/history/specs/0199-stack-rules-that-say-what-they-mean
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_04
qa_report: qa-report-2026-09-30.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0190
sources: []
regeneration:
  - command: make baseline-digests
promoted: []
pull_request: "297"
delivery_commit: 5cba1dd5e299dd8900b5d432f75ffa935c0c4edc
---

# Stack rules that say what they mean

The Context-Driven Baseline ships its TypeScript, Bun, backend and frontend rules to every adopter as Normative Clauses, and an Agent obeys them literally. An audit on 2026-09-30 found rules whose literal reading is wrong. One tells the Agent to use "Bun-owned commands" for tests, which reads as the bare Bun runner, while the profile's own Verification runs the package's `test` script. Three adopters each wrote that correction by hand. Several rules bind "the repository" when they mean one language, so a repository that composes the TypeScript modules with Go or Rust reads a rule against its own toolchain, and a Go test change dispatches a Vitest skill. One clause states two obligations, one names a profile setting that does not exist, and three places still say the suggested HTTP mode is Post-only while the product suggests REST.
