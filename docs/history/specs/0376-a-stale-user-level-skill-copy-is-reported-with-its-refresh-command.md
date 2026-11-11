---
schema: roundfix/archive-record/v1
spec: 0376-a-stale-user-level-skill-copy-is-reported-with-its-refresh-command
title: A stale user-level skill copy is reported with its refresh command
status: archived
created: "2026-11-10"
archived: "2026-11-11"
disposition: partial
source: docs/specs/0376-a-stale-user-level-skill-copy-is-reported-with-its-refresh-command
source_revision: 36f10edcbd89b6dbe214f74396479f702be77072
qa_task: task_03
qa_report: qa-report-2026-11-11.md
qa_verdict: partial
unproven:
  - the operator's run on the maintainer's MacBook after the release, recorded session evidence.
adrs:
  - ADR-0189
  - ADR-0386
sources:
  - 2026-11-10-a-user-level-skill-copy-older-than-the-project-copy-is-loaded-unnoticed.md
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

# A stale user-level skill copy is reported with its refresh command

Roundfix now reports stale user-level owned skill copies with their versions, paths, linked destinations and refresh commands while keeping command exits and Update Check next behavior. The offline QA gate passed all 40 runnable rows, including public CLI transcripts, home immutability and all twelve negative probes; the post-release MacBook validation remains unproven and no Pull Request exists yet. The offline archive preview succeeded without contacting the judge; its one adopted backlog reference adds no reusable knowledge beyond ADR-0386, the glossary and the updated command guides, so no promotion is proposed.
