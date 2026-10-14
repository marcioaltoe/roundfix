---
schema: roundfix/archive-record/v1
spec: 0270-acceptance-rules-that-match-the-authorization-and-the-queue
title: Acceptance rules that match the authorization and the queue
status: archived
created: "2026-10-13"
archived: "2026-10-14"
disposition: partial
source: docs/specs/0270-acceptance-rules-that-match-the-authorization-and-the-queue
source_revision: 1b4a19fc3b827c95ebe4dc250f44a5b2ecdd3f56
qa_task: task_05
qa_report: qa-report-2026-10-14.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0276
sources:
  - 2026-10-13-outside-evidence-blocked-by-the-spec-authorization-does-not-qualify.md
  - 2026-10-13-a-byte-invariant-pinned-to-a-past-commit-fails-after-other-merges.md
regeneration:
  - command: make baseline-digests
  - command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
    outputs:
      - internal/baseline/assets/modules/autonomous-work.json
      - internal/baseline/assets/modules/spec-workflow.json
      - internal/baseline/module-versions.json
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - skills/testdata/owned-skill-versions.json
  - command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
---

# Acceptance rules that match the authorization and the queue

The shared QA policy accepts an exact Network-Prohibited Row with outside-evidence provenance across public acceptance, archive, settlement, Daemon settlement and Delivery Queue parking. Malformed forms still refuse. The historical 0258 replay now qualifies, the original 0264 pin produces the new strict gap, and the repaired/active corpus is clear. Guides, skills, mirrors, Baseline output and glossary agree, with unchanged settlement bytes and idempotent regeneration.
