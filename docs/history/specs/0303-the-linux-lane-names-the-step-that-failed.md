---
schema: roundfix/archive-record/v1
spec: 0303-the-linux-lane-names-the-step-that-failed
title: The Linux Lane names the step that failed
status: archived
created: "2026-10-22"
archived: "2026-10-23"
disposition: partial
source: docs/specs/0303-the-linux-lane-names-the-step-that-failed
source_revision: 0acd6d7e90e8d119b2786f10f3774f73271e40fb
qa_task: task_03
qa_report: qa-report-2026-10-23.md
qa_verdict: partial
unproven:
  - the operator's `roundfix deliver lane --all` in a scratch clone of this repository, on a throwaway commit that breaks `go.mod`, after this Spec's release.
adrs:
  - ADR-0316
  - ADR-0294
sources:
  - 2026-10-22-the-linux-lane-reports-a-failed-go-list-as-go-test.md
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - skills/testdata/owned-skill-versions.json
  - command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
---

# The Linux Lane names the step that failed

The Linux Lane now names a failed package-selecting go list and its exit code through the gate, parked blocker, lane log and public CLI; the Park Class is finding below 128 and environment from 128, while other outcomes and container arguments retain their contracts. All runnable checks passed, including both Task verification commands and all ten targeted sabotages. This is a qualifying declared partial with no findings: the published documentation lookup is prohibited by the Spec, no Pull Request exists yet, and the single post-release real-container acceptance remains unproven for the operator's declared scratch-clone check.
