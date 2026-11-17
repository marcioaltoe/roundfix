---
schema: roundfix/archive-record/v1
spec: 0399-adopters-keep-their-own-rules-in-the-specific-repository-guide
title: Adopters keep their own rules in the specific-repository guide
status: archived
created: "2026-11-16"
archived: "2026-11-17"
disposition: partial
source: docs/specs/0399-adopters-keep-their-own-rules-in-the-specific-repository-guide
source_revision: a8b702c03a3b960356efee1bfd29a1453f873c83
qa_task: task_04
qa_report: qa-report-2026-11-17-01.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0411
sources:
  - 2026-11-16-adopters-put-their-own-rules-in-baseline-owned-files.md
  - 2026-11-16-maintainer-answers.md
  - 2026-11-16-adopter-rule-placement-inventory.md
regeneration:
  - command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
    outputs:
      - internal/baseline/assets/modules/context-workflow.json
      - internal/baseline/assets/modules/core.json
      - internal/baseline/assets/modules/repository-extension.json
      - internal/baseline/module-versions.json
  - command: make baseline-digests
  - command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
    outputs:
      - AGENTS.md
      - docs/agents/agent-instructions.md
      - docs/agents/docs-layout.md
      - docs/agents/setup-context.json
      - docs/agents/skill-dispatch.md
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - .agents/skills/roundfix/SKILL.md
      - .agents/skills/setup-context-driven/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/setup-context-driven/SKILL.md
      - skills/testdata/owned-skill-versions.json
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
pull_request: "892"
type: fix
---

# Adopters keep their own rules in the specific-repository guide

The Baseline guides and starter teach the specific-repository destination. Fresh public update and Doctor flows name misplaced passages while preserving their bytes and exit behavior, all seven adopter checkouts match independent read-only inventories, and all eight sabotage tests reject their intended mutations. QA is a qualifying partial solely because no Pull Request exists yet; the eventual target-branch PR still requires its live controls.
