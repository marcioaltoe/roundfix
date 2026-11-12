---
schema: roundfix/archive-record/v1
spec: 0379-an-adr-marked-as-not-applying-in-portuguese-draws-in-no-citing-adr
title: An ADR marked as not applying in Portuguese draws in no citing ADR
status: archived
created: "2026-11-10"
archived: "2026-11-11"
disposition: partial
source: docs/specs/0379-an-adr-marked-as-not-applying-in-portuguese-draws-in-no-citing-adr
source_revision: c41a8fce14abae696255f4fd45dedbc2c50438f7
qa_task: task_03
qa_report: qa-report-2026-11-11.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0389
sources:
  - 2026-11-10-an-adr-marked-as-not-applying-in-portuguese-still-seeds-the-related-adrs.md
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
pull_request: "833"
type: feat
---

# An ADR marked as not applying in Portuguese draws in no citing ADR

The strict check now recognizes the closed Portuguese marker forms and suppresses ADRs citing only a not-applying ADR; unsupported wording retains the existing closure and the fix names both languages. Fresh public CLI comparisons, positive suites and six assertion-level sabotage failures support this result. The report is a qualifying partial solely because no Pull Request exists and Merge-Ready acceptance is unavailable. The offline archive plan lists the adopted Backlog Entry as its sole promotion candidate; it adds no reusable knowledge beyond ADR-0389 and the amended glossary/guides, so no promotion is proposed.
