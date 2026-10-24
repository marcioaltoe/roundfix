---
schema: roundfix/archive-record/v1
spec: 0306-the-follow-up-command-lists-pending-work-in-priority-order
title: The Follow-up Command lists a repository's pending work in priority order
status: archived
created: "2026-10-22"
archived: "2026-10-24"
disposition: partial
source: docs/specs/0306-the-follow-up-command-lists-pending-work-in-priority-order
source_revision: 873ddaaeaed7be401ac7566bad5aa503af140ef2
qa_task: task_05
qa_report: qa-report-2026-10-23-02.md
qa_verdict: partial
unproven:
  - the operator's report of that first run, with its judge state and model id, after this Spec's release.
adrs:
  - ADR-0319
sources:
  - 2026-10-22-follow-up-status-and-unattended-modes-for-the-roundfix-skill.md
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - skills/testdata/owned-skill-versions.json
  - command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
    outputs:
      - docs/references/behavior-surfaces.json
promoted:
  - docs/references/2026-10-22-follow-up-status-and-unattended-modes-for-the-roundfix-skill.md
relinked:
  - docs/specs/0307-the-roundfix-skill-modes-and-up-to-three-questions-per-call/_prd.md
  - docs/specs/0307-the-roundfix-skill-modes-and-up-to-three-questions-per-call/task_01.md
  - docs/specs/0307-the-roundfix-skill-modes-and-up-to-three-questions-per-call/task_02.md
  - docs/specs/0307-the-roundfix-skill-modes-and-up-to-three-questions-per-call/task_03.md
---

# The Follow-up Command lists a repository's pending work in priority order

The Follow-up Command's pending inventory, priority fallback, stop state and status view passed fresh offline CLI journeys, named test cases and eleven isolated sabotage checks at a527e2a22c6cc7405c2c3439a89b47254e380419, and the corrected guide examples match the produced reports.
