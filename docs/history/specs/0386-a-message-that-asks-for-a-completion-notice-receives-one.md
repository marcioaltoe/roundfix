---
schema: roundfix/archive-record/v1
spec: 0386-a-message-that-asks-for-a-completion-notice-receives-one
title: A message that asks for a completion notice receives one
status: archived
created: "2026-11-12"
archived: "2026-11-13"
disposition: partial
source: docs/specs/0386-a-message-that-asks-for-a-completion-notice-receives-one
source_revision: 69537f0e1f874d7e35aaa7eacb3e558436a17a4c
qa_task: task_03
qa_report: qa-report-2026-11-13.md
qa_verdict: partial
unproven:
  - the operator's release session evidence after the release, recorded session evidence.
adrs:
  - ADR-0383
  - ADR-0396
  - ADR-0390
  - ADR-0330
sources:
  - 2026-11-12-maintainer-an-exception-for-a-message-that-asks-for-a-notice.md
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
promoted: []
type: feat
---

# A message that asks for a completion notice receives one

The refreshed Baseline gives a requesting message one completion notice after its work is delivered, preserves the request for the later session and forbids notices from Daemon turns; the follow-up and release references, glossary and amendment notices agree. The local refresh, regression suites, scope audit and all eight negative probes passed at the audited head.
