---
schema: roundfix/archive-record/v1
spec: 0372-a-review-whose-reviewed-head-is-gone-opens-a-new-lineage
title: A review whose reviewed head is gone opens a new lineage
status: archived
created: "2026-11-10"
archived: "2026-11-10"
disposition: partial
source: docs/specs/0372-a-review-whose-reviewed-head-is-gone-opens-a-new-lineage
source_revision: 3813dd11142b8d02f521305d539d1cf6f43c9385
qa_task: task_03
qa_report: qa-report-2026-11-10.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0197
sources:
  - 2026-11-10-conexus-review-blocks-when-the-reviewed-head-is-gone.md
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
promoted: []
pull_request: "808"
type: fix
---

# A review whose reviewed head is gone opens a new lineage

A review whose recorded head or merge base a complete repository proves absent now opens a new Reviewer Lineage at round 1; disposal and correction proof name missing commits, while incomplete history and other Git failures still block. The current-head tests, historical-base comparison, five independent sabotages and actual disposal CLI processes support that result. The offline archive plan listed the adopted Conexus source and skipped the judge; the behavior is already documented in the durable review/deliver guides and skill references. QA is partial only because the pre-PR Pull Request row lacks Merge-Ready acceptance.
