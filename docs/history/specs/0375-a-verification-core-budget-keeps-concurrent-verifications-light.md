---
schema: roundfix/archive-record/v1
spec: 0375-a-verification-core-budget-keeps-concurrent-verifications-light
title: A Verification Core Budget keeps concurrent Verifications light
status: archived
created: "2026-11-10"
archived: "2026-11-11"
disposition: partial
source: docs/specs/0375-a-verification-core-budget-keeps-concurrent-verifications-light
source_revision: 65a87310e071ac1de5b252020b0a2024a28a7180
qa_task: task_05
qa_report: qa-report-2026-11-11-01.md
qa_verdict: partial
unproven:
  - the operator's doctor line and a Verification diagnostic from the first cycle after the release, recorded session evidence.
adrs:
  - ADR-0385
sources: []
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
type: feat
---

# A Verification Core Budget keeps concurrent Verifications light

The Verification Core Budget passes its runnable gate at 84d94eabb6b3c78ddc0be860c7a4e08a1cbf2779, including the repaired Linux Lane single-source names, affected regressions, documentation contracts, scope audit and all 14 sabotage checks. The gate closes as a qualifying declared partial because the operator’s post-release doctor/Verification diagnostics remain unproven and no Pull Request exists. The offline archive plan exited 0, skipped the judge and listed no promotion candidates; plan evidence (`docs/specs/0375-a-verification-core-budget-keeps-concurrent-verifications-light/qa/evidence/qa-report-2026-11-11-01/archive-plan.txt` (not reachable from any branch or tag)). No archive or publication was performed.
