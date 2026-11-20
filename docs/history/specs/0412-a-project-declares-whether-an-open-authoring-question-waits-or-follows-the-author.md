---
schema: roundfix/archive-record/v1
spec: 0412-a-project-declares-whether-an-open-authoring-question-waits-or-follows-the-author
title: A project declares whether an open authoring question waits or follows the author
status: archived
created: "2026-11-19"
archived: "2026-11-20"
disposition: partial
source: docs/specs/0412-a-project-declares-whether-an-open-authoring-question-waits-or-follows-the-author
source_revision: 96b26cb438cf0e2ffb305a43cbbdb81adff5e933
qa_task: task_05
qa_report: qa-report-2026-11-20.md
qa_verdict: partial
unproven:
  - the operator's session evidence after the release, recorded session evidence.
  - the operator's count, at a cycle close, of the Delegated Decisions and their answers in the PRDs and of the Waiting Declarations in the open Backlog Entries, recorded session evidence.
adrs:
  - ADR-0442
  - ADR-0280
sources:
  - 2026-11-19-maintainer-answers.md
  - 2026-11-19-how-an-open-authoring-question-is-handled-today.md
  - 2026-11-19-prototype-and-governed-paths.md
regeneration:
  - command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
    outputs:
      - internal/baseline/assets/modules/spec-workflow.json
      - internal/baseline/module-versions.json
  - command: make baseline-digests
  - command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
    outputs:
      - .secondbrain-export
      - docs/adr/0280-a-guide-task-records-its-own-surface-fingerprint-a-prd-records-who-decided-and-the-glossary-is-held-to-its-own-rule.md
      - docs/agents/setup-context.json
      - docs/agents/spec-routing.md
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - .agents/skills/roundfix/SKILL.md
      - .agents/skills/write-prd/SKILL.md
      - .agents/skills/write-tasks/SKILL.md
      - .agents/skills/write-techspec/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/testdata/owned-skill-versions.json
      - skills/write-prd/SKILL.md
      - skills/write-tasks/SKILL.md
      - skills/write-techspec/SKILL.md
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
pull_request: "935"
type: feat
---

# A project declares whether an open authoring question waits or follows the author

This gate exercises the local preference, decision record and status surfaces at 820aad094c589cee8f308dbe2eb379c6ea8b7556. The archive would retain the report and its actual verdict, with the two post-release declarations unproven. No real unattended authoring flow or later-cycle metric was observed, and no Pull Request or Merge-Ready acceptance exists. The offline archive plan is recorded in archive-plan.txt (`docs/specs/0412-a-project-declares-whether-an-open-authoring-question-waits-or-follows-the-author/qa/evidence/qa-report-2026-11-20/archive-plan.txt` (not reachable from any branch or tag)); no archive was applied.
