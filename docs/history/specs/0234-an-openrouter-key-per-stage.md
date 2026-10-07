---
schema: roundfix/archive-record/v1
spec: 0234-an-openrouter-key-per-stage
title: An OpenRouter key per stage
status: archived
created: "2026-10-05"
archived: "2026-10-05"
disposition: pass
source: docs/history/specs/0234-an-openrouter-key-per-stage
source_revision: a30dc847a037fe584812b4471aa7da1701a9ec59
qa_task: task_04
qa_report: qa-report-2026-10-05.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0239
sources: []
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "414"
delivery_commit: 3c15428d186395b38ccf0ff89a0739139a16ef7e
---

# An OpenRouter key per stage

Roundfix spends on OpenRouter in two stages: the Jev judge that advises Spec authoring, and, once Spec 0233 ships the light tier, implementation of a Task on an open model through OpenCode. Both read the one Roundfix key `ROUNDFIX_OPENROUTER_API_KEY`, so OpenRouter's activity export, which groups spend by API key, shows one total and cannot say what each stage cost. On 2026-10-05 the maintainer decided to create one key per stage, each with its own monthly limit, so the export measures each stage exactly. This Spec lets each stage read its own key first, keeps the shared key as the documented fallback until the new keys exist, and records which variable each stage used without ever recording a key.
