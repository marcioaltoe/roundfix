---
schema: roundfix/archive-record/v1
spec: 0367-the-strict-check-reads-unreachable-acceptance-and-related-adrs-as-written
title: The strict check reads Unreachable Acceptance and related ADRs as written
status: archived
created: "2026-11-09"
archived: "2026-11-09"
disposition: partial
source: docs/specs/0367-the-strict-check-reads-unreachable-acceptance-and-related-adrs-as-written
source_revision: 01e6c073df386984a33aebb3673be14bedc840c6
qa_task: task_04
qa_report: qa-report-2026-11-09.md
qa_verdict: partial
unproven:
  - the fiscus operator runs `roundfix spec check 0081-a-exclusao-retirada-sempre-reenvia --strict` after the patch release and reports the result through ai-memory, recorded session evidence.
adrs:
  - ADR-0377
  - ADR-0093
  - ADR-0168
sources:
  - 2026-11-07-the-strict-check-names-an-unreachable-acceptance-it-cannot-read.md
  - 2026-11-07-fiscus-sc-adr-related-grows-through-not-applicable-adrs.md
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - .agents/skills/roundfix/SKILL.md
      - .agents/skills/write-tasks/SKILL.md
      - .agents/skills/write-techspec/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/testdata/owned-skill-versions.json
      - skills/write-tasks/SKILL.md
      - skills/write-techspec/SKILL.md
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
pull_request: "797"
type: fix
---

# The strict check reads Unreachable Acceptance and related ADRs as written

The delivered strict check names the historical unread acceptance entry at its first line and stops ADRs marked not applicable from drawing in citing ADRs, while documented fielded entries and applied-ADR behavior retain their contracts.
