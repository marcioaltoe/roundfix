---
schema: roundfix/archive-record/v1
spec: 0170-authoring-that-fails-before-dispatch
title: Authoring that fails before dispatch
status: archived
created: "2026-09-25"
archived: "2026-09-25"
disposition: pass
source: docs/history/specs/0170-authoring-that-fails-before-dispatch
source_revision: 31c30719398d5c2e43af79d5868b1268fc495f10
qa_task: task_05
qa_report: qa-report-2026-09-25-01.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0093
sources:
  - 2026-09-25-spec-check-refuses-undeclared-governed-paths.md
  - 2026-09-25-skill-and-guide-ship-with-cli-behavior.md
  - 2026-09-25-templates-ask-for-a-label-the-checker-refuses.md
  - 2026-09-25-verification-pipe-to-grep-hides-tool-status.md
  - 2026-09-25-declared-none-must-fit-one-line.md
  - 2026-09-25-adopted-source-left-behind-passes.md
  - 2026-09-25-closure-fields-lack-one-field-tests.md
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "257"
delivery_commit: 44523ea585020e599f096f512161e493f843340a
---

# Authoring that fails before dispatch

Sixteen Runs of Specs 0155 to 0169 ended Unresolved. Seven of them failed for an authoring mistake that a file read could have caught before any Agent turn was spent:
