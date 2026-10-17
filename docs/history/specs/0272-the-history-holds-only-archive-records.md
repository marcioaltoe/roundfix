---
schema: roundfix/archive-record/v1
spec: 0272-the-history-holds-only-archive-records
title: The history holds only Archive Records
status: archived
created: "2026-10-16"
archived: "2026-10-17"
disposition: partial
source: docs/specs/0272-the-history-holds-only-archive-records
source_revision: 10a9c7590d9362307370378d17dd5fd825af4b31
qa_task: task_05
qa_report: qa-report-2026-10-17-01.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0278
  - ADR-0202
sources:
  - 2026-10-13-an-archive-record-outcome-keeps-links-relative-to-the-removed-spec-folder.md
  - 2026-10-16-the-baseline-keeps-retired-adrs-and-names-two-capture-doors.md
  - 2026-10-14-promoted-copies-keep-prose-paths-to-unpromoted-folder-files.md
regeneration:
  - command: make baseline-digests
  - command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
    outputs:
      - internal/baseline/assets/modules/context-workflow.json
      - internal/baseline/assets/modules/secondbrain.json
      - internal/baseline/module-versions.json
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - skills/testdata/owned-skill-versions.json
  - command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
    outputs:
      - docs/references/behavior-surfaces.json
  - command: go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1
    outputs:
      - docs/references/coverage-record.json
promoted: []
---

# The history holds only Archive Records

The audited tree holds exactly the delivery base's Archive Records minus the ten superseded records.
