---
schema: roundfix/archive-record/v1
spec: 0405-authoring-model-fallback
title: An Authoring Phase falls back to its declared fallback when the declared model's usage limit runs out
status: archived
created: "2026-11-17"
archived: "2026-11-18"
disposition: partial
source: docs/specs/0405-authoring-model-fallback
source_revision: c8208c77f527c18120f7b367c47851649a1d24bc
qa_task: task_05
qa_report: qa-report-2026-11-18-01.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0427
  - ADR-0428
  - ADR-0430
  - ADR-0431
sources:
  - 2026-11-17-maintainer-answers.md
  - 2026-11-17-harness-and-test-measurements.md
  - 2026-11-18-root-statement-decision-and-measurement.md
regeneration:
  - command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
    outputs:
      - internal/baseline/assets/modules/spec-workflow.json
      - internal/baseline/module-versions.json
  - command: make baseline-digests
  - command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
    outputs:
      - AGENTS.md
      - docs/agents/setup-context.json
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - .agents/skills/roundfix/SKILL.md
      - .agents/skills/write-idea/SKILL.md
      - .agents/skills/write-prd/SKILL.md
      - .agents/skills/write-tasks/SKILL.md
      - .agents/skills/write-techspec/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/testdata/owned-skill-versions.json
      - skills/write-idea/SKILL.md
      - skills/write-prd/SKILL.md
      - skills/write-tasks/SKILL.md
      - skills/write-techspec/SKILL.md
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
pull_request: "910"
type: feat
---

# An Authoring Phase falls back to its declared fallback when the declared model's usage limit runs out

At audited head `35cad72e18787b97af0b0bb327bfd703160d7d05`, all runnable declared acceptance checks pass, including the amended no-fallback compatibility requirement and fresh Baseline adoption. The report is a qualifying partial solely because no Pull Request is open; there is no unresolved product finding or declared unproven acceptance. The Daemon owns QA Task settlement, and the Supervisor owns the later review and delivery sequence.
