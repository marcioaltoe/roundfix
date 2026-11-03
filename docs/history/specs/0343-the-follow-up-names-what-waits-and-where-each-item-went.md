---
schema: roundfix/archive-record/v1
spec: 0343-the-follow-up-names-what-waits-and-where-each-item-went
title: The follow-up names what waits on someone else and where each item went
status: archived
created: "2026-11-03"
archived: "2026-11-03"
disposition: partial
source: docs/specs/0343-the-follow-up-names-what-waits-and-where-each-item-went
source_revision: 2d44b24913c083d3293453ea9f96fc7bf750b86a
qa_task: task_04
qa_report: qa-report-2026-11-03.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0319
  - ADR-0338
  - ADR-0353
sources:
  - 2026-11-02-the-follow-up-names-what-waits-and-where-each-item-went.md
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
pull_request: "733"
type: feat
---

# The follow-up names what waits on someone else and where each item went

The gate records a qualifying partial: every runnable criterion passed and only the pre-PR Pull Request row remains environment-blocked. Hidden inbox entries are excluded, Waiting Items are counted and held back unless --all is requested, malformed declarations stay visible, original list/status bytes are preserved without declarations, and the skill documents each destination. The offline archive plan listed the cut and evidence candidates with the judge skipped; durable knowledge already resides in CONTEXT.md, accepted ADR-0353 and the command/skill guides, so no additional promotion was identified. The Daemon owns settlement and any later archive; this QA session performed no archive mutation.
