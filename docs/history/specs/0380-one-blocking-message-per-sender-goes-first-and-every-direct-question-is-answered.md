---
schema: roundfix/archive-record/v1
spec: 0380-one-blocking-message-per-sender-goes-first-and-every-direct-question-is-answered
title: One blocking message per sender goes first, and every direct question is answered
status: archived
created: "2026-11-11"
archived: "2026-11-12"
disposition: partial
source: docs/specs/0380-one-blocking-message-per-sender-goes-first-and-every-direct-question-is-answered
source_revision: 802102356fce66bb1e23a6d0fe2e7dc322813d99
qa_task: task_04
qa_report: qa-report-2026-11-12.md
qa_verdict: partial
unproven:
  - the operator's follow-up session evidence after the release, recorded session evidence.
adrs:
  - ADR-0390
  - ADR-0383
  - ADR-0319
  - ADR-0330
sources:
  - 2026-11-11-conexus-a-question-without-declared-blocked-work-has-no-reply-rule.md
  - 2026-11-11-quaestor-a-declared-blocked-message-jumps-every-queue-without-bound.md
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
      - skills/roundfix/SKILL.md
      - skills/testdata/owned-skill-versions.json
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
type: feat
---

# One blocking message per sender goes first, and every direct question is answered

The audited implementation bounds declared precedence to the oldest blocking message of e…
