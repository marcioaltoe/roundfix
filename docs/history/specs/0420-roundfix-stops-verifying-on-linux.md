---
schema: roundfix/archive-record/v1
spec: 0420-roundfix-stops-verifying-on-linux
title: Roundfix stops verifying on Linux
status: archived
created: "2026-11-24"
archived: "2026-11-24"
disposition: partial
source: docs/specs/0420-roundfix-stops-verifying-on-linux
source_revision: e65f87d8440aaacd13a56a326f2dca4671cae848
qa_task: task_05
qa_report: qa-report-2026-11-24-01.md
qa_verdict: partial
unproven:
  - the Delivery Queue waits for the required checks of the delivering pull request before it merges, and an item whose check fails parks without merging
  - the operator reads that run before cutting the release, and a budget overrun is corrected with the measured time before the release
  - the operator runs the release workflow's manual dispatch for the next version, which publishes nothing, before pushing the tag
adrs:
  - ADR-0398
  - ADR-0365
sources: []
regeneration:
  - command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
    outputs:
      - docs/agents/setup-context.json
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - .agents/skills/roundfix/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/testdata/owned-skill-versions.json
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
  - command: go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1
    outputs:
      - docs/references/coverage-record.json
promoted: []
pull_request: "958"
type: refactor
---

# Roundfix stops verifying on Linux

The maintainer runs Roundfix only on macOS, and Roundfix verifies on Linux in three places: the Delivery Queue's optional Linux Lane, which tests a candidate's changed Go packages in a Docker container before publication, the job that verifies every pull request and e…
