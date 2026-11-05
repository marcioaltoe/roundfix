---
schema: roundfix/archive-record/v1
spec: 0356-a-main-merge-is-a-review-delta-and-a-segment-is-measured-from-its-effective-base
title: A main merge is a review delta, and a segment is measured from its Effective Base
status: archived
created: "2026-11-05"
archived: "2026-11-05"
disposition: partial
source: docs/specs/0356-a-main-merge-is-a-review-delta-and-a-segment-is-measured-from-its-effective-base
source_revision: 9e850f74a5dc4760270ddb5b157b135da8988912
qa_task: task_05
qa_report: qa-report-2026-11-05.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0366
  - ADR-0169
  - ADR-0197
  - ADR-0282
  - ADR-0286
  - ADR-0352
sources:
  - 2026-11-05-fiscus-the-delivery-queue-takes-days-and-a-main-merge-restarts-the-review.md
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

# A main merge is a review delta, and a segment is measured from its Effective Base

The local QA gate confirms that a conflict-free merge of the default branch preserves the next review delta, both providers segment only the candidate’s own work, and the delivery commands name a review restart before it runs. All required regressions, package suites, six isolated sabotages, outside-source comparisons, wording, glossary and scope checks passed. Pull Request Merge-Ready acceptance remains unproven because no Pull Request is open.
