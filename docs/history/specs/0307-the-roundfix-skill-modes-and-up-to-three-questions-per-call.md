---
schema: roundfix/archive-record/v1
spec: 0307-the-roundfix-skill-modes-and-up-to-three-questions-per-call
title: The Roundfix skill's follow-up, status and unattended modes, and up to three questions per call
status: archived
created: "2026-10-22"
archived: "2026-10-25"
disposition: partial
source: docs/specs/0307-the-roundfix-skill-modes-and-up-to-three-questions-per-call
source_revision: b163ad545f2c63bdafebfbbcea7ba571db1fd99e
qa_task: task_04
qa_report: qa-report-2026-10-25.md
qa_verdict: partial
unproven:
  - the maintainer's first follow-up session after the release, the use he asked for on 2026-10-22, recorded as this Spec's unproven step
adrs:
  - ADR-0320
sources: []
regeneration:
  - command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
    outputs:
      - internal/baseline/assets/modules/core.json
      - internal/baseline/module-versions.json
  - command: make baseline-digests
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - skills/testdata/owned-skill-versions.json
promoted: []
---

# The Roundfix skill's follow-up, status and unattended modes, and up to three questions per call

The Archive Record may carry a qualifying partial for the audited candidate: the Baseline retains its clause identity while allowing three independent questions, the skill's bare/dashed modes and answer-placement contracts pass all offline checks and sixteen negative probes, and live harness behavior remains the PRD's expressly unproven post-release step; external publications were not fetched under the offline authorization, no Pull Request exists, and the offline archive plan skipped the judge without modifying the checkout.
