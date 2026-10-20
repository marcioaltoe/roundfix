---
schema: roundfix/archive-record/v1
spec: 0282-the-linux-lane-runs-like-ci
title: The Linux Lane runs like CI
status: archived
created: "2026-10-19"
archived: "2026-10-19"
disposition: partial
source: docs/specs/0282-the-linux-lane-runs-like-ci
source_revision: b026a989f1fb3fc600aa71dd5c060b27eb022fdc
qa_task: task_05
qa_report: qa-report-2026-10-19-01.md
qa_verdict: partial
unproven:
  - 'the operator runs `roundfix deliver lane` on `main` after this Spec merges and before the change that sets `delivery.linux_lane: true`.'
adrs:
  - ADR-0294
  - ADR-0288
sources:
  - 2026-10-18-the-linux-lane-container-lacks-ci-tools-and-an-init.md
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

# The Linux Lane runs like CI

On 2026-10-18 the Linux Lane ran for the first time, on Spec 0280's candidate, and parked the item `gate-failed` although every package the candidate changed passed in the container. Two environment gaps failed it:
