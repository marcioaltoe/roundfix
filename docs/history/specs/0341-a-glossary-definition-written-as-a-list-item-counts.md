---
schema: roundfix/archive-record/v1
spec: 0341-a-glossary-definition-written-as-a-list-item-counts
title: A glossary definition written as a list item counts
status: archived
created: "2026-11-02"
archived: "2026-11-02"
disposition: partial
source: docs/specs/0341-a-glossary-definition-written-as-a-list-item-counts
source_revision: 1b170d5682050abb9a7636decd170e88ce6e1130
qa_task: task_03
qa_report: qa-report-2026-11-02.md
qa_verdict: partial
unproven: []
adrs: []
sources:
  - 2026-11-02-a-glossary-definition-written-as-a-list-item-is-not-read.md
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - .agents/skills/roundfix/SKILL.md
      - .agents/skills/write-prd/SKILL.md
      - .agents/skills/write-tasks/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/testdata/owned-skill-versions.json
      - skills/write-prd/SKILL.md
      - skills/write-tasks/SKILL.md
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
pull_request: "729"
type: fix
---

# A glossary definition written as a list item counts

The glossary check now recognizes the specified top-level bullet and spaced-em-dash definitions, preserves rejection of prose and nested items, and gives authors the accepted forms and a verified anchored command. The worktree build, required suites, nineteen public fixture cases, eight sabotage probes, documentation/history audits and both offline external-source comparisons support this result.
