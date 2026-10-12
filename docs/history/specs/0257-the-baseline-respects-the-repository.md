---
schema: roundfix/archive-record/v1
spec: 0257-the-baseline-respects-the-repository
title: The Baseline respects the repository's own choices
status: archived
created: "2026-10-12"
archived: "2026-10-12"
disposition: pass
source: docs/specs/0257-the-baseline-respects-the-repository
source_revision: c85dda9a403e24307fdd51c1da4c811a08663349
qa_task: task_05
qa_report: qa-report-2026-10-12.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0263
  - ADR-0264
sources:
  - 2026-10-12-the-baseline-hardcodes-the-better-auth-base-path.md
  - 2026-10-12-a-skill-newer-than-the-setup-snapshot-is-reported-as-trailing.md
  - 2026-10-12-the-baseline-does-not-say-how-to-handle-a-diverging-vendored-skill.md
regeneration:
  - command: make baseline-digests
  - command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
    outputs:
      - internal/baseline/assets/modules/core.json
      - internal/baseline/assets/modules/rust.json
      - internal/baseline/module-versions.json
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
pull_request: "501"
---

# The Baseline respects the repository's own choices

Three adopter reports show the Baseline overriding choices that belong to the repository.
