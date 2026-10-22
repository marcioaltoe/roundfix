---
schema: roundfix/archive-record/v1
spec: 0299-a-new-delivery-queue-resumes-an-item-whose-spec-is-archived-on-its-branch
title: A new Delivery Queue resumes an item whose Spec is already archived on its branch
status: archived
created: "2026-10-22"
archived: "2026-10-22"
disposition: partial
source: docs/specs/0299-a-new-delivery-queue-resumes-an-item-whose-spec-is-archived-on-its-branch
source_revision: 339a6368b20ec78dc43225fc6f937ba008f23028
qa_task: task_04
qa_report: qa-report-2026-10-22-01.md
qa_verdict: partial
unproven:
  - the operator's `roundfix deliver status` showing that item `merged` after the first such queue following this Spec's release.
adrs:
  - ADR-0312
sources:
  - 2026-10-22-a-new-delivery-queue-parks-an-item-whose-spec-is-archived-on-its-branch.md
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

# A new Delivery Queue resumes an item whose Spec is already archived on its branch

The report records a qualifying declared partial with 13 passed rows, two environment-blocked rows and one declared-blocked row, no failed, finding-blocked, skipped or pending rows; real post-release queue completion remains unproven until the operator supplies the public status evidence named by Success Metric 6, and the offline Archive Command plan was inspected without applying an archive.
