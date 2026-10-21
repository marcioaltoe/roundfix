---
schema: roundfix/archive-record/v1
spec: 0290-the-pull-request-record-commit-fits-commitlint
title: The Pull Request Record Commit fits commitlint
status: archived
created: "2026-10-20"
archived: "2026-10-20"
disposition: partial
source: docs/specs/0290-the-pull-request-record-commit-fits-commitlint
source_revision: 61dec34e047bb57f6cb20f11963dca812427b92e
qa_task: task_03
qa_report: qa-report-2026-10-20.md
qa_verdict: partial
unproven:
  - the oraculum operator's next delivery of such a Spec after this Spec's release.
adrs:
  - ADR-0302
sources:
  - 2026-10-20-o-commit-do-registro-do-pull-request-estoura-o-commitlint-com-slug-longo.md
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

# The Pull Request Record Commit fits commitlint

The Pull Request Record Commit now uses a bounded numeric Spec subject, accepts historical hand commits regardless of their subjects, and reports queue hook refusals with their hook, subject, exit status and recovery. Real Git, unchanged Daemon tests, documentation contracts, scope checks and all four negative controls passed; the report is a qualifying partial because external documentation refresh is prohibited offline, no Pull Request exists yet, and the single post-release oraculum delivery remains unproven for its operator to verify.
