---
schema: roundfix/archive-record/v1
spec: 0283-a-git-locator-is-reachable-in-ci-and-the-storage-report-counts-every-page
title: A Git Locator is reachable in CI, and the storage report counts every page
status: archived
created: "2026-10-19"
archived: "2026-10-19"
disposition: partial
source: docs/specs/0283-a-git-locator-is-reachable-in-ci-and-the-storage-report-counts-every-page
source_revision: 81792eec2c9180001869ac8aac1f9c867caf5204
qa_task: task_04
qa_report: qa-report-2026-10-19.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0295
  - ADR-0296
sources:
  - 2026-10-18-a-git-locator-in-a-fresh-record-fails-the-check-in-ci.md
  - 2026-10-18-storage-report-fails-while-the-run-database-wal-is-live.md
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - skills/testdata/owned-skill-versions.json
  - command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
pull_request: "593"
---

# A Git Locator is reachable in CI, and the storage report counts every page

The Archive Record may state that fetched-default Git Locator reachability is accepted in a CI-style clone and over stale local history, all three consumers share the fetched head, and storage snapshots reconcile exactly including pointer-map/lock-byte pages and live WAL frames. The historical 1,363,968-byte gap is accounted for exactly. This is a qualifying partial: published-source live reads and the pre-PR controls remain unproven for their recorded environment causes. Offline archive planning lists both adopted Backlog Entries; their reusable conclusions already live in ADR-0295, ADR-0296, CONTEXT.md and the updated guides, so no additional promotion is proposed. No archive was performed.
