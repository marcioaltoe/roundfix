---
schema: roundfix/archive-record/v1
spec: 0292-runs-causes-reads-the-diagnostic-where-the-run-wrote-it
title: Runs causes reads the diagnostic where the Run wrote it
status: archived
created: "2026-10-20"
archived: "2026-10-21"
disposition: partial
source: docs/specs/0292-runs-causes-reads-the-diagnostic-where-the-run-wrote-it
source_revision: b642fb65a3d97782cd44f83772ad60687277c17b
qa_task: task_03
qa_report: qa-report-2026-10-21.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0301
sources:
  - 2026-10-20-runs-causes-reads-no-diagnostic-of-a-run-that-recorded-no-artifact-root.md
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - skills/testdata/owned-skill-versions.json
  - command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
pull_request: "619"
---

# Runs causes reads the diagnostic where the Run wrote it

The observed feature meets every runnable criterion at e14d09a163f49b4cb7372de8e27742730509ec98, including the exact blank-root transcript, relative-root classification, all diagnostic guards, sabotage sensitivity, privacy, unchanged signatures/help and bounded Task scope. The report records a qualifying partial because only the pre-PR Pull Request row and the outside-evidence lookup prohibited by Spec authorization remain environment-blocked. The required `./bin/roundfix archive 0292-runs-causes-reads-the-diagnostic-where-the-run-wrote-it --plan --offline` ran in the disposable archived head with a disposable HOME, every Judge Key Variable and TYPESAFE_API_KEY unset, exited 0 and reported `judge: skipped (offline)`; it planned an Archive Record and removal of nine Spec files, and performed no archive.
