---
schema: roundfix/archive-record/v1
spec: 0300-reconcile-releases-a-stale-run-owner-record
title: Reconcile releases a stale Run owner record
status: archived
created: "2026-10-22"
archived: "2026-10-23"
disposition: partial
source: docs/specs/0300-reconcile-releases-a-stale-run-owner-record
source_revision: 2f885128efd0ce7d6a66ae44fe0457038d1a27b8
qa_task: task_03
qa_report: qa-report-2026-10-23.md
qa_verdict: partial
unproven:
  - the operator's `roundfix reconcile`, `roundfix reconcile --apply` and second `roundfix reconcile` in this repository's checkout after this Spec's release.
adrs:
  - ADR-0313
sources:
  - 2026-10-21-reconcile-preserves-forever-a-process-candidate-whose-owner-pid-now-belongs-to-another-process.md
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - skills/testdata/owned-skill-versions.json
  - command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
pull_request: "641"
---

# Reconcile releases a stale Run owner record

Qualifying declared partial: nine rows passed, two are environment-blocked (the authorization-prohibited outside lookup and the pre-PR Pull Request row), and one is declared-blocked (Success Metric 7 on the maintainer Home). There are no failed, skipped, pending or finding-blocked rows. Runnable journeys, all named Verification tests, Windows vet, documentation contracts and all ten guard sabotages completed. The offline archive plan exited 0 with the judge skipped, and no archive was performed; the maintainer production criterion remains unproven for the future Archive Record. Public `./bin/roundfix qa-report accept` accepted this partial with exit 0 and empty stdout/stderr (report-acceptance.
