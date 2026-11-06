---
schema: roundfix/archive-record/v1
spec: 0361-a-session-starts-with-an-update-check-and-the-update-mode-applies-it
title: A session starts with an Update Check, and the update mode applies it
status: archived
created: "2026-11-06"
archived: "2026-11-06"
disposition: partial
source: docs/specs/0361-a-session-starts-with-an-update-check-and-the-update-mode-applies-it
source_revision: 5eed63da9a2a85ab3d3e027d8818fdd8f397b04b
qa_task: task_05
qa_report: qa-report-2026-11-06.md
qa_verdict: partial
unproven:
  - the operator's session evidence in an adopting repository after the release, recorded session evidence.
adrs:
  - ADR-0371
sources:
  - 2026-11-06-an-update-check-at-session-start-replaces-the-release-notes-in-the-inboxes.md
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

# A session starts with an Update Check, and the update mode applies it

The offline archive plan lists the Spec cut and its knowledge candidates; the judge was skipped offline. The command changed no repository file and did not publish or archive anything.
