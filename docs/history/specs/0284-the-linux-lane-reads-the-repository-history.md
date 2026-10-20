---
schema: roundfix/archive-record/v1
spec: 0284-the-linux-lane-reads-the-repository-history
title: The Linux Lane reads the repository history
status: archived
created: "2026-10-19"
archived: "2026-10-20"
disposition: partial
source: docs/specs/0284-the-linux-lane-reads-the-repository-history
source_revision: 35ede566f7d1ebcfd7299ddaf2a8892d243d62e0
qa_task: task_02
qa_report: qa-report-2026-10-20.md
qa_verdict: partial
unproven:
  - 'the operator runs `roundfix deliver lane --all` on `main` after this Spec merges and before the change that sets `delivery.linux_lane: true`.'
adrs:
  - ADR-0297
sources:
  - 2026-10-19-the-linux-lane-cannot-read-the-repository-history.md
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - skills/testdata/owned-skill-versions.json
  - command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
---

# The Linux Lane reads the repository history

The Archive Record can carry a qualifying `partial` for audited head `28f925214daca49e847f9e8e9cfc5691c802ed76`, with the three network-prohibited outside-source rows and the pre-PR row unproven. The real Linux Lane passed. The offline archive plan exited 0, listed the Spec-owned adopted reference as its only promotion candidate and explicitly skipped the judge. That source's durable decision is already in accepted ADR-0297 and the revised glossary and guides; no additional promotion or second repository write is needed. Archive, review and Task settlement remain Daemon-owned.
