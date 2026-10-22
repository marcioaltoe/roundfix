---
schema: roundfix/archive-record/v1
spec: 0296-the-queues-archive-subject-and-pull-request-title-fit-100-characters
title: The queue's archive subject and pull request title fit 100 characters
status: archived
created: "2026-10-21"
archived: "2026-10-21"
disposition: partial
source: docs/specs/0296-the-queues-archive-subject-and-pull-request-title-fit-100-characters
source_revision: 3dc0b8dd17eeb4cebc861e71ef4cb16d8a079438
qa_task: task_03
qa_report: qa-report-2026-10-21.md
qa_verdict: partial
unproven:
  - the operator's `git log -1 --format=%s` of that squash commit on the default branch after the first such delivery following this Spec's release.
adrs:
  - ADR-0309
sources:
  - 2026-10-20-the-queue-archive-subject-and-pull-request-title-grow-past-100-characters-with-the-slug.md
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

# The queue's archive subject and pull request title fit 100 characters

The gate records a qualifying declared partial: eight runnable rows pass, the outside remote lookup and Pull Request rows remain environment-blocked, and the single post-release squash declaration remains unproven. The Run-built `./bin/roundfix archive 0296-the-queues-archive-subject-and-pull-request-title-fit-100-characters --plan --offline` exited 0 under a disposable HOME with every Judge Key Variable and TYPESAFE_API_KEY unset; it proposed the Archive Record, identified core and evidence files, and reported `judge: skipped (offline)`. Evidence is `evidence/2026-10-21-gate/archive-plan.txt`. This was a plan only; no archive, Task settlement, commit, push or Pull Request was performed.
