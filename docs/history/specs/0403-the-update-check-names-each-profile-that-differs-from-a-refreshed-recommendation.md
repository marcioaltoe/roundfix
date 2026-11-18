---
schema: roundfix/archive-record/v1
spec: 0403-the-update-check-names-each-profile-that-differs-from-a-refreshed-recommendation
title: The Update Check names each profile that differs from a refreshed recommendation
status: archived
created: "2026-11-17"
archived: "2026-11-18"
disposition: partial
source: docs/specs/0403-the-update-check-names-each-profile-that-differs-from-a-refreshed-recommendation
source_revision: aaa2e4d7d1dbcdff4d124337f97bf0fac244f7a8
qa_task: task_05
qa_report: qa-report-2026-11-18-01.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0181
  - ADR-0416
sources:
  - 2026-11-17-maintainer-answers.md
  - 2026-11-17-model-evidence.md
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

# The Update Check names each profile that differs from a refreshed recommendation

The Update Check names each differing configured category and identifies an older pin with its re-pin key and adoption command, while preserving quiet output and legacy profile behavior. The refreshed snapshot, catalog, repository configuration, documentation and offline model evidence satisfy the declared checks. QA is partial solely because no Pull Request exists yet; all runnable acceptance rows passed.
