---
schema: roundfix/archive-record/v1
spec: 0338-the-release-decision-travels-with-the-spec
title: The release decision travels with the Spec
status: archived
created: "2026-11-01"
archived: "2026-11-02"
disposition: partial
source: docs/specs/0338-the-release-decision-travels-with-the-spec
source_revision: e5a86f33fe51bee98e08ecdc425e33d697d52f22
qa_task: task_05
qa_report: qa-report-2026-11-02.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0143
  - ADR-0309
  - ADR-0349
sources:
  - 2026-11-01-the-delivery-pull-request-type-follows-the-spec.md
  - 2026-11-01-a-spec-authorization-carries-the-release-decision.md
regeneration:
  - command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
    outputs:
      - internal/baseline/assets/modules/core.json
      - internal/baseline/assets/modules/spec-workflow.json
      - internal/baseline/module-versions.json
  - command: make baseline-digests
  - command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
    outputs:
      - docs/agents/agent-instructions.md
      - docs/agents/setup-context.json
      - docs/agents/spec-routing.md
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - .agents/skills/roundfix/SKILL.md
      - .agents/skills/write-prd/SKILL.md
      - .agents/skills/write-tasks/SKILL.md
      - .agents/skills/write-techspec/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/testdata/owned-skill-versions.json
      - skills/write-prd/SKILL.md
      - skills/write-tasks/SKILL.md
      - skills/write-techspec/SKILL.md
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
type: feat
---

# The release decision travels with the Spec

The PRD's Delivery Type now survives archivi…
