---
schema: roundfix/archive-record/v1
spec: 0276-the-history-cleanup-reaches-every-adopter-and-the-review-sees-a-large-spec
title: The history cleanup reaches every adopter, and the review sees a large Spec
status: archived
created: "2026-10-17"
archived: "2026-10-17"
disposition: partial
source: docs/specs/0276-the-history-cleanup-reaches-every-adopter-and-the-review-sees-a-large-spec
source_revision: 3d01553521cecc2e387b33ea9811fbf956757a2b
qa_task: task_05
qa_report: qa-report-2026-10-17.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0282
sources:
  - 2026-10-17-one-unnamed-superseded-adr-keeps-an-adopters-whole-history.md
  - 2026-10-17-the-pre-pr-review-refuses-a-spec-that-deletes-retired-history.md
  - 2026-10-16-review-convention-c2-asks-a-qa-task-for-settlement-sections.md
  - 2026-10-16-spec-check-does-not-flag-an-owned-skill-edit-without-its-version-record.md
  - 2026-10-16-a-qa-gate-archive-plan-contacts-the-judge-an-offline-qa-forbids.md
  - 2026-10-16-an-item-binary-built-in-place-can-be-killed-by-code-signing.md
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

# The history cleanup reaches every adopter, and the review sees a large Spec

The offline archive plan exited 0, printed `judge: skipped (offline)`, and listed all six adopted reference candidates, in plan order, with no advice: offline QA archive planning, item-binary replacement, QA C2, owned-skill version enforcement, successor lineage, and large review candidates.
