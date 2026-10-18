---
schema: roundfix/archive-record/v1
spec: 0278-a-fixed-round-2-finding-gets-one-correction-review-and-the-plan-names-only-what-it-delivers
title: A fixed round-2 finding gets one Correction Review, and the plan names only what it delivers
status: archived
created: "2026-10-18"
archived: "2026-10-18"
disposition: partial
source: docs/specs/0278-a-fixed-round-2-finding-gets-one-correction-review-and-the-plan-names-only-what-it-delivers
source_revision: c73cd612da8a6abfac4c6f14683d036f92401216
qa_task: task_05
qa_report: qa-report-2026-10-18.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0286
  - ADR-0272
  - ADR-0287
sources: []
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - skills/testdata/owned-skill-versions.json
  - command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
pull_request: "575"
---

# A fixed round-2 finding gets one Correction Review, and the plan names only what it delivers

The implemented candidate gives an all-fixed round-2 finding exactly one Correction Review of its fix and admits that fix through the retry proof; the historical broken plan raises SC-VERIFY-PATH-UNDECLARED while the repaired plan does not; and a one-byte Makefile change is caught by the Setup Manifest contract. All runnable declared verification Requirements passed, including three required sabotages. The Archive Record may carry this qualifying partial with the two authorization-prohibited literature lookups and the pre-PR/Merge-Ready control explicitly unproven.
