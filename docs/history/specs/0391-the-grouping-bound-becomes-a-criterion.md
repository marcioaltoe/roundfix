---
schema: roundfix/archive-record/v1
spec: 0391-the-grouping-bound-becomes-a-criterion
title: The grouping bound becomes a criterion
status: archived
created: "2026-11-14"
archived: "2026-11-14"
disposition: partial
source: docs/specs/0391-the-grouping-bound-becomes-a-criterion
source_revision: 7642d3692746e38ea5af79aa082266a96a79298a
qa_task: task_03
qa_report: qa-report-2026-11-14.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0208
sources:
  - 2026-11-14-maintainer-the-grouping-bound-becomes-a-criterion.md
regeneration:
  - command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
    outputs:
      - internal/baseline/assets/modules/spec-workflow.json
      - internal/baseline/module-versions.json
  - command: make baseline-digests
  - command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
    outputs:
      - docs/agents/setup-context.json
      - docs/agents/spec-routing.md
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - .agents/skills/roundfix/SKILL.md
      - .agents/skills/write-idea/SKILL.md
      - .agents/skills/write-prd/SKILL.md
      - .agents/skills/write-tasks/SKILL.md
      - .agents/skills/write-techspec/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/testdata/owned-skill-versions.json
      - skills/write-idea/SKILL.md
      - skills/write-prd/SKILL.md
      - skills/write-tasks/SKILL.md
      - skills/write-techspec/SKILL.md
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
type: fix
---

# The grouping bound becomes a criterion

The Archive Record can carry a qualifying partial: every runnable criterion passed, and only the pre-PR Pull Request row is blocked because no PR is open.
