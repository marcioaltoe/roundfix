---
schema: roundfix/archive-record/v1
spec: 0293-archive-and-cancel-refuse-a-link-that-leaves-the-repository
title: Archive and cancel refuse a link that leaves the repository
status: archived
created: "2026-10-20"
archived: "2026-10-21"
disposition: partial
source: docs/specs/0293-archive-and-cancel-refuse-a-link-that-leaves-the-repository
source_revision: 1e57814805a3d510c2f562ea31befeaf7299c616
qa_task: task_04
qa_report: qa-report-2026-10-21.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0306
  - ADR-0269
sources:
  - 2026-10-20-archive-and-cancel-accept-a-link-that-leaves-the-repository.md
regeneration:
  - command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
    outputs:
      - docs/references/behavior-surfaces.json
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - skills/testdata/owned-skill-versions.json
promoted: []
---

# Archive and cancel refuse a link that leaves the repository

The audited head refuses all six outside layouts before writing and preserves accepted-path output, records and relinked bytes. Regression tests, four guard-sabotage proofs, documentation contracts and protected-file comparisons support that result; the Delivery Base reproduces the original outside write. The authorized `./bin/roundfix archive 0293-archive-and-cancel-refuse-a-link-that-leaves-the-repository --plan --offline` exited 0, planned the Archive Record under docs/history/specs, reported the judge skipped offline and left Git status unchanged, with all Judge Key Variables and TYPESAFE_API_KEY absent. The report records a qualifying partial solely for the prohibited Go blog lookup and absent Pull Request; archive and QA settlement remain Daemon-owned.
