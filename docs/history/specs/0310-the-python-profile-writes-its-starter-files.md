---
schema: roundfix/archive-record/v1
spec: 0310-the-python-profile-writes-its-starter-files
title: The Python profile writes its Starter Files
status: archived
created: "2026-10-23"
archived: "2026-10-26"
disposition: partial
source: docs/specs/0310-the-python-profile-writes-its-starter-files
source_revision: 3701a1afa72b68e74a394705dd17ef7f4fbe1a40
qa_task: task_04
qa_report: qa-report-2026-10-26.md
qa_verdict: partial
unproven:
  - the first adopter that writes the Starter Files after the release, which the operator announces to `default/marcioaltoe-quaestor`, recorded as this Spec's unproven step
adrs:
  - ADR-0323
sources: []
regeneration:
  - command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
    outputs:
      - internal/baseline/assets/modules/python.json
      - internal/baseline/module-versions.json
  - command: make baseline-digests
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - skills/testdata/owned-skill-versions.json
  - command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
---

# The Python profile writes its Starter Files

The offline archive plan completed with judge skipped (offline), and no archive was applied. The runnable Starter File journeys and declared checks pass at 8f2f3e9cbb829f0a5de4edc370028068776fbf4a; the report is a qualifying declared partial under the QA settlement policy, with the single first-adopter acceptance left unproven and the forbidden live documentation lookup and absent Pull Request recorded separately. The Daemon owns settlement and the QA commit.
