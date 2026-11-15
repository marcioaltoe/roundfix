---
schema: roundfix/archive-record/v1
spec: 0396-the-verification-probe-prepares-its-tree-as-the-daemon-does
title: The Verification probe prepares its tree as the Daemon does
status: archived
created: "2026-11-14"
archived: "2026-11-15"
disposition: partial
source: docs/specs/0396-the-verification-probe-prepares-its-tree-as-the-daemon-does
source_revision: 48d1475e33ee57ea38631de93d45ae7e29422758
qa_task: task_02
qa_report: qa-report-2026-11-15.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0406
sources:
  - 2026-11-14-indexus-the-verification-probe-called-honest-a-command-the-daemon-refused-as-vacuous.md
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
pull_request: "880"
type: fix
---

# The Verification probe prepares its tree as the Daemon does

The authoring Verification probe now copies configured files and runs the Worktree Bootstrap before probing, so dependency-dependent commands receive the same vacuous verdict the Daemon gives them. Public text and JSON transcripts, the delivery-base comparison, all positive suites, eight isolated sabotages and the history audit support the change, with unconfigured repositories retaining byte-identical output. QA is a qualifying partial solely because no Pull Request or Merge-Ready acceptance exists yet.
