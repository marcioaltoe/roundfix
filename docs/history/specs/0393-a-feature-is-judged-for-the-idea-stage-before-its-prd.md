---
schema: roundfix/archive-record/v1
spec: 0393-a-feature-is-judged-for-the-idea-stage-before-its-prd
title: A feature is judged for the Idea stage before its PRD
status: archived
created: "2026-11-14"
archived: "2026-11-14"
disposition: partial
source: docs/specs/0393-a-feature-is-judged-for-the-idea-stage-before-its-prd
source_revision: 9ce96a5a977c751064cef59fcd96e95dbefc666c
qa_task: task_04
qa_report: qa-report-2026-11-14-01.md
qa_verdict: partial
unproven:
  - the operator's session evidence after the release, recorded session evidence.
  - the calibration report of Spec 0395, recorded session evidence.
adrs:
  - ADR-0403
  - ADR-0404
sources:
  - 2026-11-14-maintainer-request-and-answers.md
  - 2026-11-14-idea-stage-judgment-pilot.md
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
      - .agents/skills/brainstorming/SKILL.md
      - .agents/skills/roundfix/SKILL.md
      - .agents/skills/write-idea/SKILL.md
      - .agents/skills/write-prd/SKILL.md
      - skills/brainstorming/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/testdata/owned-skill-versions.json
      - skills/write-idea/SKILL.md
      - skills/write-prd/SKILL.md
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
type: feat
---

# A feature is judged for the Idea stage before its PRD

T…
