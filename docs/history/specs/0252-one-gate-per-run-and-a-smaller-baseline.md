---
schema: roundfix/archive-record/v1
spec: 0252-one-gate-per-run-and-a-smaller-baseline
title: One gate per Run and a smaller Baseline
status: archived
created: "2026-10-08"
archived: "2026-10-08"
disposition: qa-override
source: docs/specs/0252-one-gate-per-run-and-a-smaller-baseline
source_revision: 4d3e9059fa3d6704a8c5c0628cca369adbd7c2f7
qa_task: task_05
qa_report: qa-report-2026-10-08.md
qa_verdict: partial
unproven: []
qa_override: true
qa_override_approval: 'maintainer standing authorization 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: 'Environment-only partial: rows 08b and 08c need arxiv.org, which the authorized QA scope forbids; row 08d needs Secondbrain access outside the Run scope; row 12 has no Pull Request yet; row 11 is the PRD''s declared Unreachable Acceptance (fewer agent full-suite runs, measured in Runs after release). Every behavior row passed.'
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: 4d3e9059fa3d6704a8c5c0628cca369adbd7c2f7
adrs:
  - ADR-0257
sources: []
regeneration:
  - command: make baseline-digests
  - command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
    outputs:
      - internal/baseline/assets/modules/autonomous-work.json
      - internal/baseline/assets/modules/context-workflow.json
      - internal/baseline/assets/modules/core.json
      - internal/baseline/assets/modules/secondbrain.json
      - internal/baseline/assets/modules/spec-workflow.json
      - internal/baseline/module-versions.json
  - command: make skills-sync
promoted: []
---

# One gate per Run and a smaller Baseline

Read-only `./bin/roundfix archive 0252-one-gate-per-run-and-a-smaller-baseline --plan` exited 0; [plan output](evidence/2026-10-08-gate/archive-plan.txt) inventories the Spec core and QA evidence, and identifies no reusable-knowledge candidate. This is an inventory, not an eligibility verdict or archive execution.
