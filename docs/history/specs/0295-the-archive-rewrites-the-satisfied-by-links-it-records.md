---
schema: roundfix/archive-record/v1
spec: 0295-the-archive-rewrites-the-satisfied-by-links-it-records
title: The archive rewrites the satisfied-by links it records as unproven
status: archived
created: "2026-10-21"
archived: "2026-10-21"
disposition: partial
source: docs/specs/0295-the-archive-rewrites-the-satisfied-by-links-it-records
source_revision: e8e58e87ed29c3507d3f1312624b554ff9b2c534
qa_task: task_04
qa_report: qa-report-2026-10-21.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0308
sources:
  - 2026-10-20-o-archive-copia-links-relativos-das-declaracoes-para-o-archive-record.md
regeneration:
  - command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
    outputs:
      - docs/references/behavior-surfaces.json
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - skills/testdata/owned-skill-versions.json
promoted: []
pull_request: "626"
---

# The archive rewrites the satisfied-by links it records as unproven

The archive now rebases each unproven entry through the same Legacy Link Rebase as its outcome, and the Record Link Repair and strict checker name an existing promoted copy. Both archive shapes, Delivery Base comparison, exact sanitize transcript, applied repair persistence, dry-run byte preservation, 26 named tests and all three sabotage proofs passed; all three Task Verification commands passed and the Daemon's repository Verification passed. The gate is a qualifying partial solely for the network-prohibited GitHub Docs row and the pre-PR row, with no finding, declared block or skipped row, and no additional knowledge promotion is needed.
