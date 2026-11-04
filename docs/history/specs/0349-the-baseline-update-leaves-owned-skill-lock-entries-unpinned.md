---
schema: roundfix/archive-record/v1
spec: 0349-the-baseline-update-leaves-owned-skill-lock-entries-unpinned
title: The baseline update leaves owned skill lock entries unpinned
status: archived
created: "2026-11-04"
archived: "2026-11-04"
disposition: partial
source: docs/specs/0349-the-baseline-update-leaves-owned-skill-lock-entries-unpinned
source_revision: d97ebf46c08e756d129bdc3f7592f34b12772c63
qa_task: task_04
qa_report: qa-report-2026-11-04.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0359
sources:
  - 2026-11-04-the-baseline-update-keeps-a-stale-ref-in-owned-skill-lock-entries.md
regeneration:
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
type: fix
---

# The baseline update leaves owned skill lock entries unpinned

The bundle install now removes the stale catalog pin from existing installed owned-skill lock entries while preserving external entries and all other owned metadata, reports the change, and leaves the lock untouched on replay. Doctor names pinned or mis-hashed owned entries and one bundle install clears the warning. Required tests, eight semantic sabotage probes and offline external-source comparisons passed; the pre-PR row remains unproven because no Pull Request or Merge-Ready acceptance exists. The offline archive plan (`docs/specs/0349-the-baseline-update-leaves-owned-skill-lock-entries-unpinned/qa/evidence/2026-11-04-owned-lock/archive-plan.txt` (not reachable from any branch or tag)) ran with the Run-built binary and all Judge Key Variables plus TYPESAFE_API_KEY unset, reported judge skipped offline and listed its candidates without contacting a judge.
