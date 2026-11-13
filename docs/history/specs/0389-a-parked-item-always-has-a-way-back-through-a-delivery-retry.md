---
schema: roundfix/archive-record/v1
spec: 0389-a-parked-item-always-has-a-way-back-through-a-delivery-retry
title: A parked item always has a way back through a Delivery Retry
status: archived
created: "2026-11-13"
archived: "2026-11-13"
disposition: partial
source: docs/specs/0389-a-parked-item-always-has-a-way-back-through-a-delivery-retry
source_revision: 520fada70b2591f8b6e8f4b5c94a0de04696d047
qa_task: task_05
qa_report: qa-report-2026-11-13.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0398
sources:
  - 2026-11-13-a-retry-refuses-an-item-parked-before-it-has-a-branch.md
  - 2026-11-12-a-retry-after-resolving-a-pull-request-conflict-finds-no-specs-root-and-no-pushed-head.md
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

# A parked item always has a way back through a Delivery Retry

A Delivery Retry requeues a branchless parked item, including a half-recorded workspace, and an archived Spec's authorization remains readable after its Spec Root disappears. An earlier Pull Request candidate returns through the repository gate before the owner pushes, while an unknown head stays parked with the specific operator action. All runnable checks and eight negative controls passed at the audited head; live Pull Request Merge-Ready acceptance remains unproven because no Pull Request exists. The offline plan named two adopted backlog references; their reusable behavior is already recorded in the glossary, accepted ADR-0398, guide and synchronized skill, so neither needs another promotion.
