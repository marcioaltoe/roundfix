---
schema: roundfix/archive-record/v1
spec: 0369-reconcile-discards-confirmed-dirty-leftovers-and-ignored-only-worktrees
title: Reconcile discards confirmed dirty leftovers and releases ignored-only worktrees
status: archived
created: "2026-11-09"
archived: "2026-11-09"
disposition: partial
source: docs/specs/0369-reconcile-discards-confirmed-dirty-leftovers-and-ignored-only-worktrees
source_revision: 03043ebf3cc1f40956f99f9696a54bc962bbf7e9
qa_task: task_04
qa_report: qa-report-2026-11-09.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0379
  - ADR-0289
sources:
  - 2026-11-09-fiscus-reconcile-leaves-abandoned-run-leftovers-without-a-path.md
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

# Reconcile discards confirmed dirty leftovers and releases ignored-only worktrees

An operator can inspect and confirm an archived Spec's ended Run leftovers, obtain a Branch Disposition naming the discarded paths, and release the worktree and branch through the supported command; a changed confirmation preserves the work. Ignored-only unregistered Run directories and empty repository directories also leave through the supported apply flow, while mixed files and unrelated names stay. Offline archive planning listed the adopted Fiscus report as the sole promotion candidate and skipped the Judge; its observations are already captured by accepted ADR-0379 and the maintained reconcile guide, so no additional durable promotion is proposed. Local QA is complete; Pull Request Merge-Ready acceptance remains unproven until the Supervisor publishes the candidate.
