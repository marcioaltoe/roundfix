---
schema: roundfix/archive-record/v1
spec: 0298-the-delivery-queues-check-wait-is-configurable-default-60-minutes
title: The Delivery Queue's check wait is configurable, default 60 minutes
status: archived
created: "2026-10-22"
archived: "2026-10-22"
disposition: partial
source: docs/specs/0298-the-delivery-queues-check-wait-is-configurable-default-60-minutes
source_revision: 19b955a583969ba937c86500ee247124d681139b
qa_task: task_03
qa_report: qa-report-2026-10-22.md
qa_verdict: partial
unproven:
  - the operator's `roundfix deliver status` line for that item after the first such delivery following this Spec's release.
adrs:
  - ADR-0311
sources:
  - 2026-10-22-the-delivery-queue-parks-checks-timeout-while-ci-is-still-running.md
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - skills/testdata/owned-skill-versions.json
  - command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
pull_request: "634"
---

# The Delivery Queue's check wait is configurable, default 60 minutes

The Delivery Queue Check Wait is configurable only through User Config, defaults to 60 minutes, refuses values outside the positive-through-24h bound, ignores Project Config with its warning, starts each deadline independently, retains one rerun restart, and parks with the configured key and duration. All runnable checks pass; the verdict is a qualifying declared partial, with the post-release real-CI criterion unproven, the external documentation lookup prohibited by authorization and no open Pull Request. The offline archive plan exited 0, printed `judge: skipped (offline)` and planned an Archive Record, without changing Task status or archiving the Spec.
