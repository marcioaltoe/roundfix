---
schema: roundfix/archive-record/v1
spec: 0334-the-template-tooling-row-passes-and-a-wave-collision-refuses-before-the-run
title: The template's tooling row passes the Spec Check, and a Wave collision is refused before the Run
status: archived
created: "2026-10-30"
archived: "2026-10-31"
disposition: partial
source: docs/specs/0334-the-template-tooling-row-passes-and-a-wave-collision-refuses-before-the-run
source_revision: c7e1eff9ed0a5f2b9395f2e6c4129ea9fbf8b344
qa_task: task_03
qa_report: qa-report-2026-10-31.md
qa_verdict: partial
unproven:
  - Quaestor's reply through ai-memory to the release message the operator sends, recorded session evidence.
adrs:
  - ADR-0131
  - ADR-0347
sources:
  - 2026-10-27-the-authoring-templates-and-the-tooling-check-disagree-on-no-tooling-change.md
  - 2026-10-27-a-wave-collision-with-a-completed-task-refuses-after-the-run-starts.md
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - .agents/skills/roundfix/SKILL.md
      - .agents/skills/write-prd/SKILL.md
      - .agents/skills/write-techspec/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/testdata/owned-skill-versions.json
      - skills/write-prd/SKILL.md
      - skills/write-techspec/SKILL.md
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
---

# The template's tooling row passes the Spec Check, and a Wave collision is refused before the Run

The template-derived no-mutation Tooling authority row passes with an operative empty grant, while nonempty grants and explicit mutations remain bounded; completed Tasks never collide, and a real parallel collision is refused before a Run exists.
