---
schema: roundfix/archive-record/v1
spec: 0291-the-queue-shows-a-skipped-linux-lane-and-waits-for-a-run-slot
title: The queue shows a skipped Linux Lane and waits for a Run slot
status: archived
created: "2026-10-20"
archived: "2026-10-21"
disposition: partial
source: docs/specs/0291-the-queue-shows-a-skipped-linux-lane-and-waits-for-a-run-slot
source_revision: 4923930ca997abf9f5b77bcd36b7d942d7fced32
qa_task: task_05
qa_report: qa-report-2026-10-21.md
qa_verdict: partial
unproven:
  - the maintainer runs `roundfix doctor` in this repository after the release that carries this Spec, with the image absent or after `docker image rm golang:1.26` followed by `docker pull golang:1.26`.
adrs:
  - ADR-0303
  - ADR-0304
sources:
  - 2026-10-20-the-queue-skips-the-linux-lane-without-a-trace.md
  - 2026-10-20-the-queue-parks-an-item-at-the-run-ceiling.md
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - skills/testdata/owned-skill-versions.json
  - command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
pull_request: "616"
---

# The queue shows a skipped Linux Lane and waits for a Run slot

The gate records a qualifying declared partial at audited head `29944e748dcb701f63301175144532f16d180fef`: all runnable criteria passed, both public status transcripts matched on repeated processes and all eight sabotages were detected. The only unmet acceptance is the declared real-Docker doctor run after release, with two network-prohibited external-source rows and the pre-PR Pull Request row also recorded. The offline archive plan passed with Judge Key Variables and TYPESAFE_API_KEY unset; the judge reported skipped (offline), no external request was made and no archive was performed.
