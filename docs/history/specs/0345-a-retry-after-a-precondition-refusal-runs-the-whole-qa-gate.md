---
schema: roundfix/archive-record/v1
spec: 0345-a-retry-after-a-precondition-refusal-runs-the-whole-qa-gate
title: A retry after a QA precondition refusal runs the whole QA gate
status: archived
created: "2026-11-04"
archived: "2026-11-04"
disposition: partial
source: docs/specs/0345-a-retry-after-a-precondition-refusal-runs-the-whole-qa-gate
source_revision: 7cd58990ff50ee99f07e4cdd39ccda36b3590872
qa_task: task_03
qa_report: qa-report-2026-11-04.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0194
  - ADR-0249
sources:
  - 2026-11-03-pantheon-retry-reimports-a-precondition-refusal.md
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
type: fix
---

# A retry after a QA precondition refusal runs the whole QA gate

A retry after a precondition refusal now rejects the whole prior refusal candidate and runs its gate again, while failed matrix passes retain their byte-for-byte import and carry behavior. Refusal formatting and failed-formatter restoration passed, all four sabotages exposed their intended regressions, and the adopter documentation, glossary and accepted ADR describe the same behavior. The only unobserved delivery control is the later Pull Request review and Merge-Ready acceptance, recorded as the pre-PR environment exception; durable knowledge already resides outside the Spec and the offline archive plan requires no promotion.
