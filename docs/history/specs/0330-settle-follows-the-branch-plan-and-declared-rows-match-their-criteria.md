---
schema: roundfix/archive-record/v1
spec: 0330-settle-follows-the-branch-plan-and-declared-rows-match-their-criteria
title: Settle follows the branch's plan and the Task commit's rules, and a declared row matches the criterion it names
status: archived
created: "2026-10-27"
archived: "2026-10-27"
disposition: partial
source: docs/specs/0330-settle-follows-the-branch-plan-and-declared-rows-match-their-criteria
source_revision: 7e3eae01ee21b3e9b06c2caea9ebfc0fbc5a7b2b
qa_task: task_05
qa_report: qa-report-2026-10-27.md
qa_verdict: partial
unproven:
  - Quaestor's reply through ai-memory to the release message the operator sends, recorded session evidence.
adrs:
  - ADR-0343
sources:
  - 2026-10-26-settle-run-inside-the-run-worktree-takes-it-for-the-user-checkout.md
  - 2026-10-26-settle-and-a-tracked-mode-change-commit-an-executable-the-task-commit-refuses.md
  - 2026-10-27-a-declaration-covering-two-criteria-is-refused-and-settle-reads-the-run-worktree-plan.md
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - .agents/skills/qa-gate/SKILL.md
      - .agents/skills/roundfix/SKILL.md
      - .agents/skills/write-prd/SKILL.md
      - skills/qa-gate/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/testdata/owned-skill-versions.json
      - skills/write-prd/SKILL.md
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
---

# Settle follows the branch's plan and the Task commit's rules, and a declared row matches the criterion it names

The built product accepts one declaration covering two criteria, preserves count compatibility, refuses settlement from Run and Task Worktrees, follows an amended branch plan and applies executable authorization. All prescribed positive checks and ten behavioral sabotage probes passed.
