---
schema: roundfix/archive-record/v1
spec: 0406-claude-md-is-always-a-link-to-agents-md
title: CLAUDE.md is always a symbolic link to AGENTS.md
status: archived
created: "2026-11-18"
archived: "2026-11-18"
disposition: partial
source: docs/specs/0406-claude-md-is-always-a-link-to-agents-md
source_revision: 658106f5af6ca7bba3aee632e18daa836c17f34c
qa_task: task_04
qa_report: qa-report-2026-11-18.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0432
  - ADR-0433
sources:
  - 2026-11-18-maintainer-decision.md
  - 2026-11-18-claude-link-inventory-and-measurements.md
regeneration:
  - command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
    outputs:
      - internal/baseline/assets/modules/core.json
      - internal/baseline/module-versions.json
  - command: make baseline-digests
  - command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
    outputs:
      - docs/agents/agent-instructions.md
      - docs/agents/setup-context.json
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
type: fix
---

# CLAUDE.md is always a symbolic link to AGENTS.md

The Archive Record can carry a qualifying partial: fourteen rows passed and only the pre-PR Pull Request row is environment-blocked because no Pull Request is open.
