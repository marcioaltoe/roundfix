---
schema: roundfix/archive-record/v1
spec: 0277-a-run-never-reaches-the-judge-and-an-amended-plan-keeps-its-settled-tasks
title: A Run never reaches the judge, and an amended plan keeps its settled Tasks
status: archived
created: "2026-10-17"
archived: "2026-10-18"
disposition: partial
source: docs/specs/0277-a-run-never-reaches-the-judge-and-an-amended-plan-keeps-its-settled-tasks
source_revision: ff9bd3ed20a2a0354d36d045ec4538d275bb9b6a
qa_task: task_05
qa_report: qa-report-2026-10-17.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0283
  - ADR-0273
  - ADR-0284
  - ADR-0285
sources: []
regeneration:
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

# A Run never reaches the judge, and an amended plan keeps its settled Tasks

The Archive Record may carry a qualifying partial: the Run Marker and credential boundaries, bootstrap isolation, Windows compilation, conservative Binding Text carry-forward and adopter ADR references were verified, with remote acceptance evidence explicitly unverified and the pre-PR row retained. `./bin/roundfix archive 0277-a-run-never-reaches-the-judge-and-an-amended-plan-keeps-its-settled-tasks --plan --offline` exited 0 with the judge skipped offline, listing ten core files and QA evidence, with no promotion candidates. The exact cut is captured at invocation in outcome-archive-plan-offline.
