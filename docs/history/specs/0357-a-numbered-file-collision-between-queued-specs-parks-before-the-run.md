---
schema: roundfix/archive-record/v1
spec: 0357-a-numbered-file-collision-between-queued-specs-parks-before-the-run
title: A numbered file collision between queued Specs parks before the Run
status: archived
created: "2026-11-05"
archived: "2026-11-05"
disposition: partial
source: docs/specs/0357-a-numbered-file-collision-between-queued-specs-parks-before-the-run
source_revision: 1a2c59745bd3dbfce63a345f50c6549269cb74e5
qa_task: task_05
qa_report: qa-report-2026-11-05.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0367
sources:
  - 2026-11-05-fiscus-delivery-queue-report.md
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
pull_request: "765"
type: feat
---

# A numbered file collision between queued Specs parks before the Run

The declared numbered-file checks, pre-Run park and retry refusal, status action and supporting documentation passed the authored QA checks, including all 13 independent sabotage probes. QA is a qualifying partial solely because no Pull Request exists yet and Merge-Ready review acceptance remains unavailable; the Daemon repository Verification passed. The offline archive plan completed with the judge skipped and identified the adopted fiscus incident report as its single knowledge candidate; the rule and its limitations already live in CONTEXT.md, ADR-0367 and the user guides, so no additional promotion is proposed. No archive, review request or delivery mutation occurred.
