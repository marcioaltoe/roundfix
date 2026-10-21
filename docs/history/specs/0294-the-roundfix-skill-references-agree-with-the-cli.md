---
schema: roundfix/archive-record/v1
spec: 0294-the-roundfix-skill-references-agree-with-the-cli
title: The Roundfix skill references agree with the CLI
status: archived
created: "2026-10-21"
archived: "2026-10-21"
disposition: partial
source: docs/specs/0294-the-roundfix-skill-references-agree-with-the-cli
source_revision: 3d96ff89ad69c46c94fc9e01bd389f947d88f5b6
qa_task: task_03
qa_report: qa-report-2026-10-21.md
qa_verdict: partial
unproven:
  - the quaestor operator's next Baseline refresh and pre-PR review after this Spec's release.
adrs: []
sources:
  - 2026-10-20-duas-contradicoes-nas-referencias-review-e-events-da-skill.md
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - skills/testdata/owned-skill-versions.json
  - command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
pull_request: "624"
---

# The Roundfix skill references agree with the CLI

The Roundfix review, events and deliver references now match the existing CLI, including Archive Record semantics, all six event categories and the revised Run Event Stream glossary entry. The source-derived guards and their six independent sabotage checks passed, both implementation Task Verification checks passed, and the history audit found no production behavior or unauthorized path changes. All 16 runnable rows passed; the two prohibited literature lookups and the unavailable Pull Request are environment-blocked, and the future quaestor review remains the single declared unproven criterion. The offline archive plan (`docs/specs/0294-the-roundfix-skill-references-agree-with-the-cli/qa/evidence/2026-10-21-gate/archive-plan.txt` (not reachable from any branch or tag)) exited 0; this report does not claim post-release adopter verification, current-candidate review acceptance or Task settlement.
