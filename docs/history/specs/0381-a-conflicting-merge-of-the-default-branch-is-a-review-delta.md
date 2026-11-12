---
schema: roundfix/archive-record/v1
spec: 0381-a-conflicting-merge-of-the-default-branch-is-a-review-delta
title: A conflicting merge of the default branch is a review delta
status: archived
created: "2026-11-11"
archived: "2026-11-12"
disposition: partial
source: docs/specs/0381-a-conflicting-merge-of-the-default-branch-is-a-review-delta
source_revision: b1a3a2760036ac6d30063ca99488d99f6fc9864e
qa_task: task_03
qa_report: qa-report-2026-11-12.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0366
  - ADR-0197
  - ADR-0391
sources:
  - 2026-11-11-conexus-review-false-compile-findings-after-a-main-merge.md
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

# A conflicting merge of the default branch is a review delta

The offline archive plan (`docs/specs/0381-a-conflicting-merge-of-the-default-branch-is-a-review-delta/qa/evidence/qa-report-2026-11-12/archive-plan.txt` (not reachable from any branch or tag)) exited 0, identified the Spec's core files, fresh QA evidence and adopted Conexus source, and reported `judge: skipped (offline)` with no advice requested. The Archive Record can carry a qualifying partial whose sole unmet row is the pre-PR Pull Request row, with no declared-unreachable acceptance. The adopted incident's confirmed rule already lives in accepted ADR-0391, the glossary and the updated command references, so no additional knowledge promotion is needed. This QA turn did not archive the Spec or change Task status.
