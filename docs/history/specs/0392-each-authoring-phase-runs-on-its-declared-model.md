---
schema: roundfix/archive-record/v1
spec: 0392-each-authoring-phase-runs-on-its-declared-model
title: Each authoring phase runs on its declared model
status: archived
created: "2026-11-14"
archived: "2026-11-14"
disposition: partial
source: docs/specs/0392-each-authoring-phase-runs-on-its-declared-model
source_revision: 7ab8a591eb47a0daa687306f4e237be956738c94
qa_task: task_05
qa_report: qa-report-2026-11-14.md
qa_verdict: partial
unproven:
  - the operator's pilot evidence after the release, recorded session evidence.
adrs:
  - ADR-0392
  - ADR-0401
  - ADR-0402
  - ADR-0405
  - ADR-0386
sources:
  - 2026-11-13-indexus-model-and-effort-per-skill-and-a-research-subagent.md
  - 2026-11-13-pantheon-roundfix-skills-left-the-user-level.md
  - 2026-11-14-skill-model-override-measurement.md
  - 2026-11-14-maintainer-answers.md
  - 2026-11-14-subagent-model-and-effort-measurement.md
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
type: feat
---

# Each authoring phase runs on its declared model

Repositories can now install declared authoring headers and generated subagents, diagnose mismatches, and explicitly remove user-level skill copies while preserving symlink targets.
