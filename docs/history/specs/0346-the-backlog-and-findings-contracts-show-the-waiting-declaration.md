---
schema: roundfix/archive-record/v1
spec: 0346-the-backlog-and-findings-contracts-show-the-waiting-declaration
title: The Backlog and Findings contracts show the Waiting Declaration
status: archived
created: "2026-11-04"
archived: "2026-11-04"
disposition: partial
source: docs/specs/0346-the-backlog-and-findings-contracts-show-the-waiting-declaration
source_revision: 86c972311037f4e3209b50fae4dfa9612ddf65bc
qa_task: task_02
qa_report: qa-report-2026-11-04-01.md
qa_verdict: partial
unproven: []
adrs: []
sources:
  - 2026-11-03-the-backlog-and-findings-contracts-name-the-waiting-declaration.md
regeneration:
  - command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
    outputs:
      - internal/baseline/assets/modules/context-workflow.json
      - internal/baseline/module-versions.json
  - command: make baseline-digests
  - command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
    outputs:
      - docs/agents/docs-layout.md
      - docs/agents/setup-context.json
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
type: feat
---

# The Backlog and Findings contracts show the Waiting Declaration

The copyable Backlog and Findings contracts now show both optional waiting keys as comments and explain when to uncomment them; unchanged copies declare no wait, while deliberately uncommented copies carry the expected values. Adopter previews expose only the intended guide and Setup Manifest refresh, all positive and sabotage checks pass, and prior F-01 is resolved.
