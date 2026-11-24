---
schema: roundfix/archive-record/v1
spec: 0419-round-2-of-a-pre-pr-review-opens-a-new-session-when-the-configured-model-changed
title: Round 2 of a pre-PR review opens a new session when the configured model changed
status: archived
created: "2026-11-23"
archived: "2026-11-24"
disposition: partial
source: docs/specs/0419-round-2-of-a-pre-pr-review-opens-a-new-session-when-the-configured-model-changed
source_revision: e82357e8b179595e075f56f75544aa652e8d3f66
qa_task: task_02
qa_report: qa-report-2026-11-24-01.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0197
  - ADR-0451
sources: []
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

# Round 2 of a pre-PR review opens a new session when the configured model changed

The offline archive plan exited 0, skipped the judge and listed references/sources.md as its candidate; archive-plan (`docs/specs/0419-round-2-of-a-pre-pr-review-opens-a-new-session-when-the-configured-model-changed/qa/evidence/qa-report-2026-11-24-01/archive-plan.txt` (not reachable from any branch or tag)). No source was promoted because the reusable model/session rule already lives in accepted ADR-0451, the glossary, guide and skill reference. This report records a qualifying partial: every runnable check passed and only the pre-PR Pull Request row remains blocked, with no Merge-Ready acceptance claimed; real-adapter confirmation remains the operator's post-release action.
