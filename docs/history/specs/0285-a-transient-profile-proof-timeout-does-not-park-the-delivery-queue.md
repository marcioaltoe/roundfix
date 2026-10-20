---
schema: roundfix/archive-record/v1
spec: 0285-a-transient-profile-proof-timeout-does-not-park-the-delivery-queue
title: A transient profile-proof timeout does not park the Delivery Queue
status: archived
created: "2026-10-20"
archived: "2026-10-20"
disposition: partial
source: docs/specs/0285-a-transient-profile-proof-timeout-does-not-park-the-delivery-queue
source_revision: dbb4fe18913c4d5190e58f5c8a47812a20b1f862
qa_task: task_05
qa_report: qa-report-2026-10-20.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0298
sources:
  - 2026-10-20-a-prova-de-perfil-do-preflight-estoura-o-tempo-em-serie.md
regeneration:
  - command: go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1
    outputs:
      - docs/references/coverage-record.json
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - skills/testdata/owned-skill-versions.json
  - command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
pull_request: "599"
---

# A transient profile-proof timeout does not park the Delivery Queue

A profile proof whose disposable Agent Session setup times out now gets a third attempt after a backoff of 5 and then 10 seconds, its setup deadline becomes the config key `defaults.profile_proof_setup_timeout` (default `30s`), and the cleanup of a Session that never opened stops logging a misleading warning; every configured tuple is still proven.
