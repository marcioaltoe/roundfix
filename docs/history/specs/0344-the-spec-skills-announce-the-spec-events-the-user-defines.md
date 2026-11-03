---
schema: roundfix/archive-record/v1
spec: 0344-the-spec-skills-announce-the-spec-events-the-user-defines
title: The Spec skills announce the Spec Events the user defines
status: archived
created: "2026-11-03"
archived: "2026-11-03"
disposition: partial
source: docs/specs/0344-the-spec-skills-announce-the-spec-events-the-user-defines
source_revision: d0905859bb07ffa4d3314ce33cb48ee129adc389
qa_task: task_05
qa_report: qa-report-2026-11-03.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0354
sources:
  - 2026-11-03-the-spec-skills-announce-each-spec-event.md
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - .agents/skills/archive-spec/SKILL.md
      - .agents/skills/roundfix/SKILL.md
      - .agents/skills/write-prd/SKILL.md
      - .agents/skills/write-tasks/SKILL.md
      - .agents/skills/write-techspec/SKILL.md
      - skills/archive-spec/SKILL.md
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

# The Spec skills announce the Spec Events the user defines

The five Spec skill surfaces now describe announcements only when the user's own instructions define their destination and format, with all seven events assigned to their authoring or supervising session and Run agents expressly excluded. The contract tests, all four sabotage probes, independent wording and glossary checks, outside-source comparisons and history audit passed; the offline archive preview succeeded without writing an Archive Record or contacting the judge.
