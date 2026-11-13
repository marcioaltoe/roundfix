---
schema: roundfix/archive-record/v1
spec: 0385-a-queue-item-stops-only-for-its-own-work
title: A queue item stops only for its own work
status: archived
created: "2026-11-12"
archived: "2026-11-12"
disposition: partial
source: docs/specs/0385-a-queue-item-stops-only-for-its-own-work
source_revision: 5c3b370c9be60071708a8af2d3c409b3e4032b9a
qa_task: task_05
qa_report: qa-report-2026-11-12.md
qa_verdict: partial
unproven:
  - the operator, after the release, lists `ai-memory list-projects` before and after the first delivery of this repository whose pre-PR review runs the validator, with the `.ai-memory.toml` the planning change commits at the root
adrs:
  - ADR-0279
  - ADR-0395
sources:
  - 2026-11-12-the-delivery-owner-crashes-when-the-repository-gate-waits-for-machine-verification.md
  - 2026-11-11-a-review-finding-anchored-as-path-line-and-a-period-blocks-the-item.md
  - 2026-11-12-quaestor-the-delivery-owner-panics-at-the-gate.md
  - 2026-11-12-pantheon-the-review-validator-becomes-a-stray-ai-memory-project.md
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

# A queue item stops only for its own work

The gate records a qualifying declared partial: all runnable rows passed, all ten sabotages exposed the expected regression, and the only unmet rows are the declared post-release ai-memory attribution and the pre-PR Pull Request row. The Archive Record must retain that attribution as unproven. The offline archive plan listed the Spec's core, QA evidence and four adopted source candidates with judge advice skipped; no archive or promotion was performed.
