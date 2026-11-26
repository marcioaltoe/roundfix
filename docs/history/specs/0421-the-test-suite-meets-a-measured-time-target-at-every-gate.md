---
schema: roundfix/archive-record/v1
spec: 0421-the-test-suite-meets-a-measured-time-target-at-every-gate
title: The test suite meets a measured time target at every gate
status: archived
created: "2026-11-25"
archived: "2026-11-26"
disposition: partial
source: docs/specs/0421-the-test-suite-meets-a-measured-time-target-at-every-gate
source_revision: b6c5fc7b4b1ea3d1a12a91f588aa8e278d65d695
qa_task: task_05
qa_report: qa-report-2026-11-26-01.md
qa_verdict: partial
unproven:
  - the Delivery Queue waits for the required check of the delivering pull request before it merges; the operator reads the job times of that run and records them in the reference document before the release
  - the operator runs the release workflow's manual dispatch, which publishes nothing, after the release's pull request merges and before pushing the tag
  - the operator reads the windows with the command the reference document records and writes each result into that document; a miss goes to the maintainer with its measurement
  - the last delivery of this Spec, the closing pull request the maintainer decided; the contract over the reference document proves its numbers
adrs:
  - ADR-0452
  - ADR-0387
  - ADR-0252
  - ADR-0253
  - ADR-0453
sources: []
regeneration:
  - command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
    outputs:
      - docs/agents/setup-context.json
  - command: go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1
    outputs:
      - docs/references/coverage-record.json
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
pull_request: "962"
type: perf
---

# The test suite meets a measured time target at every gate

The gate exercises the lightweight Task settlement, the archive candidate’s Complete Suite Run, dependency selection, CI and release fixtu…
