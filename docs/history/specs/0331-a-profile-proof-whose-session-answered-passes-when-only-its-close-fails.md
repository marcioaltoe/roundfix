---
schema: roundfix/archive-record/v1
spec: 0331-a-profile-proof-whose-session-answered-passes-when-only-its-close-fails
title: A profile proof whose session answered passes when only its close fails
status: archived
created: "2026-10-27"
archived: "2026-10-27"
disposition: partial
source: docs/specs/0331-a-profile-proof-whose-session-answered-passes-when-only-its-close-fails
source_revision: 119d8e6544fc8665b793257d723d2cd943ce3b4e
qa_task: task_04
qa_report: qa-report-2026-10-27.md
qa_verdict: partial
unproven:
  - the release notice to `default/vortex` and the next loaded delivery of this repository's queue, reported through ai-memory or the queue log.
adrs:
  - ADR-0344
  - ADR-0140
sources:
  - 2026-10-27-a-profile-proof-fails-when-closing-its-disposable-session-under-load.md
regeneration:
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1
    outputs:
      - .agents/skills/roundfix/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/testdata/owned-skill-versions.json
  - command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
    outputs:
      - docs/references/behavior-surfaces.json
  - command: go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1
    outputs:
      - docs/references/coverage-record.json
promoted: []
---

# A profile proof whose session answered passes when only its close fails

The offline archive plan exited 0 and skipped the judge, proposed the Archive Record under docs/history/specs/, and reported no advice for the adopted reference; it changed no archive files. The final verdict is a qualifying declared partial, with Success Metric 6 and its release notice and next loaded delivery evidence remaining unproven, and the pre-PR Pull Request row carrying its equivalent controls. No release, archive, review, merge or publication was performed.
