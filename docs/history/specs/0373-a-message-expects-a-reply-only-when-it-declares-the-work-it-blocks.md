---
schema: roundfix/archive-record/v1
spec: 0373-a-message-expects-a-reply-only-when-it-declares-the-work-it-blocks
title: A message expects a reply only when it declares the work it blocks
status: archived
created: "2026-11-10"
archived: "2026-11-10"
disposition: partial
source: docs/specs/0373-a-message-expects-a-reply-only-when-it-declares-the-work-it-blocks
source_revision: 9647393dce770c560d3fa42531a39eff764a8f35
qa_task: task_04
qa_report: qa-report-2026-11-10.md
qa_verdict: partial
unproven:
  - the operator's follow-up session evidence after the release, recorded session evidence.
adrs:
  - ADR-0353
  - ADR-0383
  - ADR-0319
  - ADR-0330
sources:
  - 2026-11-10-a-message-expects-a-reply-only-when-the-sender-is-blocked.md
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
pull_request: "817"
type: feat
---

# A message expects a reply only when it declares the work it blocks

The local Follow-up Command now lists sender-declared blocked work first without judging it, refuses malformed declarations and preserves reports without that field.
