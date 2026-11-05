---
schema: roundfix/archive-record/v1
spec: 0351-a-secret-refusal-of-a-task-commit-is-repairable-and-visible
title: A secret refusal of a Task commit is repairable and visible
status: archived
created: "2026-11-04"
archived: "2026-11-05"
disposition: partial
source: docs/specs/0351-a-secret-refusal-of-a-task-commit-is-repairable-and-visible
source_revision: c3626798c98ecd5a610186a709c5056971ad68c2
qa_task: task_05
qa_report: qa-report-2026-11-05-01.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0292
  - ADR-0342
  - ADR-0361
sources:
  - 2026-11-04-fiscus-a-task-commit-secret-refusal-is-invisible.md
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

# A secret refusal of a Task commit is repairable and visible

The Archive Record can retain qualifying QA verdict `partial`, 19 passed rows including four carried observations, the resolved ADR-0292 wording finding and the single pre-PR environment block, with zero failed, finding-blocked, declared-blocked or skipped rows. The offline archive preview (`docs/specs/0351-a-secret-refusal-of-a-task-commit-is-repairable-and-visible/qa/evidence/0351-rerun/archive-plan.txt` (not reachable from any branch or tag)) records the proposed cut and promotion candidates without executing an archive or contacting the judge. ADR-0361, ADR-0292, CONTEXT.md and the command guides already own the durable rules; the QA reports and adopted Fiscus report remain provenance. The Daemon owns QA Task settlement and the later delivery and archive steps.
