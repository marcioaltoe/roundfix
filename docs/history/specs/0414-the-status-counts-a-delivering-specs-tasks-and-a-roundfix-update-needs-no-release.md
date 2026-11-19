---
schema: roundfix/archive-record/v1
spec: 0414-the-status-counts-a-delivering-specs-tasks-and-a-roundfix-update-needs-no-release
title: The status counts a delivering Spec's Tasks, and a Roundfix update needs no release
status: archived
created: "2026-11-19"
archived: "2026-11-19"
disposition: partial
source: docs/specs/0414-the-status-counts-a-delivering-specs-tasks-and-a-roundfix-update-needs-no-release
source_revision: 58db3102ddaaa9cb92abe15cd97df027eb2bc8b7
qa_task: task_03
qa_report: qa-report-2026-11-19-01.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0057
  - ADR-0447
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

# The status counts a delivering Spec's Tasks, and a Roundfix update needs no release

The Item Branch count and adopter maintenance boundary satisfy their declared checks, all eleven sabotage probes fail at the named assertions, and the Task history stays in its authorized scope. The qualifying partial leaves only the pre-PR Pull Request row blocked; no archive, promotion, commit, push or Pull Request was performed. The offline archive plan (`docs/specs/0414-the-status-counts-a-delivering-specs-tasks-and-a-roundfix-update-needs-no-release/qa/evidence/qa-report-2026-11-19-01/archive-plan-offline.txt` (not reachable from any branch or tag)) records the proposed cut and reference handling, with judge keys unset and no network request.
