---
schema: roundfix/archive-record/v1
spec: 0402-doctor-and-the-update-check-recommend-an-authoring-model-declaration
title: Doctor and the Update Check recommend an Authoring Model Declaration
status: archived
created: "2026-11-17"
archived: "2026-11-17"
disposition: partial
source: docs/specs/0402-doctor-and-the-update-check-recommend-an-authoring-model-declaration
source_revision: bf17cb383bb261740e653be5bf6ce16cb20120b0
qa_task: task_03
qa_report: qa-report-2026-11-17-01.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0415
  - ADR-0402
sources:
  - 2026-11-17-maintainer-answers.md
  - 2026-11-17-installed-authoring-skills-without-a-declaration-measurement.md
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - .agents/skills/roundfix/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/testdata/owned-skill-versions.json
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
  - command: go run -buildvcs=false ./cmd/roundfix skills install --target project
    outputs:
      - .claude/agents/roundfix-prd-author.md
      - .claude/agents/roundfix-tasks-author.md
      - .claude/agents/roundfix-techspec-author.md
promoted: []
pull_request: "903"
type: fix
---

# Doctor and the Update Check recommend an Authoring Model Declaration

The audited build prints the Recommended Authoring Declaration for installed unpublished authoring skills without a declaration, preserves read-only behavior and existing exemption cases, and holds its table, repository declaration, generated subagents and guide equal. All six sabotage probes detect their specified defect.
