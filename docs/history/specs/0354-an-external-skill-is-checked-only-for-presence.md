---
schema: roundfix/archive-record/v1
spec: 0354-an-external-skill-is-checked-only-for-presence
title: An external skill is checked only for presence
status: archived
created: "2026-11-05"
archived: "2026-11-05"
disposition: partial
source: docs/specs/0354-an-external-skill-is-checked-only-for-presence
source_revision: bfb75f4229f7c111c6678256de4abe3196802e18
qa_task: task_04
qa_report: qa-report-2026-11-05.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0364
  - ADR-0221
  - ADR-0224
  - ADR-0241
  - ADR-0264
sources:
  - 2026-11-05-only-roundfix-owned-skills-are-checked-by-digest-or-version.md
regeneration:
  - command: go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1
    outputs:
      - docs/references/coverage-record.json
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - .agents/skills/roundfix/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/testdata/owned-skill-versions.json
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
pull_request: "763"
type: fix
---

# An external skill is checked only for presence

Required external skills are accepted by folder and lock-entry presence, and an edited present copy stays ready and is preserved by the managed refresh. Missing skills, deliberate restoration, the owned-skill version boundary and this repository's stricter lock/snapshot contracts remain verified. The report records a qualifying partial solely because no Pull Request and no Merge-Ready acceptance exist yet; the offline archive plan lists the adopted Fluxus report as its reusable-knowledge candidate, with judge advice explicitly skipped.
