---
schema: roundfix/archive-record/v1
spec: 0404-a-correction-records-its-defect-class-and-a-phase-its-authoring-model
title: A Correction records its Defect Class and a phase entry its authoring model
status: archived
created: "2026-11-17"
archived: "2026-11-18"
disposition: partial
source: docs/specs/0404-a-correction-records-its-defect-class-and-a-phase-its-authoring-model
source_revision: e586e004fe8c8ae1bda3590e2d36d6a0e105bba4
qa_task: task_05
qa_report: qa-report-2026-11-18-01.md
qa_verdict: partial
unproven:
  - the operator turns the setting on after the release, reads the report after 20 delivered Specs and records the counts, the maintainer's declarations and the adopter's first delivery in the session; the Verdict of the Idea says what a miss changes.
adrs:
  - ADR-0417
  - ADR-0422
  - ADR-0423
  - ADR-0424
  - ADR-0418
  - ADR-0419
  - ADR-0407
sources:
  - 2026-11-17-maintainer-request-and-answers.md
  - 2026-11-17-fable-and-opus-authoring-study.md
  - 2026-11-17-phase-record-and-causes-view-measurement.md
  - 2026-11-17-item-branch-commits-and-author-session-measurement.md
  - 2026-11-17-techspec-design-measurements.md
  - 2026-11-17-tasks-measurements.md
regeneration:
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
pull_request: "909"
type: feat
---

# A Correction records its Defect Class and a phase entry its authoring model

The Correction and authoring-model features pass the runnable local gate, including all eight Surface Transcripts, the eight negative mutation probes, persistence across worktrees and clones, delivery-base compatibility, glossary contracts and history scope. The qualifying partial retains the seven PRD post-release satisfied-by actions as unproven and the pre-PR Pull Request row as environment-blocked; the Supervisor must obtain current-candidate review and complete those future observations under the declared policy.
