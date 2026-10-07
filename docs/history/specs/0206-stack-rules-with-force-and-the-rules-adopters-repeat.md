---
schema: roundfix/archive-record/v1
spec: 0206-stack-rules-with-force-and-the-rules-adopters-repeat
title: Stack rules with force, and the rules adopters repeat
status: archived
created: "2026-09-30"
archived: "2026-10-02"
disposition: pass
source: docs/history/specs/0206-stack-rules-with-force-and-the-rules-adopters-repeat
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_05
qa_report: qa-report-2026-10-02.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0202
  - ADR-0203
sources:
  - 2026-09-30-stack-modules-need-rules-with-force.md
regeneration:
  - command: make baseline-digests
promoted: []
pull_request: "339"
delivery_commit: f9caffee78b524df602448acf7881afd6dd44207
---

# Stack rules with force, and the rules adopters repeat

The Context-Driven Baseline ships Normative Clauses with a stated force to every adopter, but its Go, Rust, CLI and TUI modules still ship paragraphs with no force at all. An Agent in a Go repository cannot tell whether "prefer the standard library" is advice or a ban, and a Rust adopter receives no error policy. Meanwhile adopters keep writing the same rules by hand: five repositories require express authorization for every database mutation, five require Verification that fails before the change, three say that a lint warning blocks. This Spec gives every rule of the four stack and surface modules a clause with force, adds the Rust error policy the maintainer decided, and promotes into the Baseline the hand-written rules that recur, so adopters stop re-deriving them and Agents read them with force.
