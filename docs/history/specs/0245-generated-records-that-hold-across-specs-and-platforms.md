---
schema: roundfix/archive-record/v1
spec: 0245-generated-records-that-hold-across-specs-and-platforms
title: Generated records that hold across Specs and platforms
status: archived
created: "2026-10-07"
archived: "2026-10-07"
disposition: qa-override
source: docs/specs/0245-generated-records-that-hold-across-specs-and-platforms
source_revision: dfb28c42d127373e7eef2a706091a810f62ee45b
qa_task: task_04
qa_report: qa-report-2026-10-07.md
qa_verdict: partial
unproven: []
qa_override: true
qa_override_approval: 'maintainer standing authorization 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: 'Environment-only partial after corrective task_05 and the sanctioned record step: rows Q07 and Q12 need the Pull Request (Linux CI and pre-PR controls) and row Q08c needs the original authoring measurement. Every behavior row passed, including the two-branch module merge that regenerates to the next free version.'
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: dfb28c42d127373e7eef2a706091a810f62ee45b
adrs:
  - ADR-0250
sources:
  - 2026-10-07-generated-records-that-parallel-work-or-another-platform-breaks.md
regeneration:
  - command: make baseline-digests
  - command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
    outputs:
      - internal/baseline/assets/modules/autonomous-work.json
      - internal/baseline/assets/modules/backend.json
      - internal/baseline/assets/modules/bun.json
      - internal/baseline/assets/modules/cli-surface.json
      - internal/baseline/assets/modules/context-workflow.json
      - internal/baseline/assets/modules/core.json
      - internal/baseline/assets/modules/external-triage.json
      - internal/baseline/assets/modules/frontend.json
      - internal/baseline/assets/modules/go.json
      - internal/baseline/assets/modules/monorepo.json
      - internal/baseline/assets/modules/repository-extension.json
      - internal/baseline/assets/modules/rust.json
      - internal/baseline/assets/modules/secondbrain.json
      - internal/baseline/assets/modules/spec-workflow.json
      - internal/baseline/assets/modules/tui-surface.json
      - internal/baseline/assets/modules/typescript.json
      - internal/baseline/module-versions.json
promoted: []
---

# Generated records that hold across Specs and platforms

At audited head 92de34df983515675a690237f73eeffccab00cc3, QA closed partial with 11 passing rows and 3 environment-blocked rows, no product findings, and no skipped or pending row. Live recording refused unrecorded content, selected the next free version and remained idempotent; a two-branch merge preserved both edits at version 25 and passed full verification; the platform-neutral Coverage Record regenerated identical bytes and detected Linux-only test removal on macOS. The independent operator log and upstream Go documentation support the design. The original authoring measurement and actual Linux CI/Pull Request controls remain unproven, so this partial is not eligible to settle the gate under the qualifying-partial rule.
