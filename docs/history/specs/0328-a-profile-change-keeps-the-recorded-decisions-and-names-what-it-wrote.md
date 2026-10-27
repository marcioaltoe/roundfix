---
schema: roundfix/archive-record/v1
spec: 0328-a-profile-change-keeps-the-recorded-decisions-and-names-what-it-wrote
title: A profile change keeps the recorded decisions and names what it wrote
status: archived
created: "2026-10-26"
archived: "2026-10-27"
disposition: partial
source: docs/specs/0328-a-profile-change-keeps-the-recorded-decisions-and-names-what-it-wrote
source_revision: 781f3fb0f26de048c7dab1ff83cb9fefab12667b
qa_task: task_05
qa_report: qa-report-2026-10-27.md
qa_verdict: partial
unproven:
  - the release notice to `default/gesttione-solutions-quaestor` and that project's next profile change, reported to this project through ai-memory.
adrs:
  - ADR-0263
  - ADR-0341
sources:
  - 2026-10-26-a-baseline-profile-change-loses-the-recorded-decisions-and-writes-without-saying-so.md
regeneration:
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1
    outputs:
      - .agents/skills/roundfix/SKILL.md
      - .agents/skills/setup-context-driven/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/setup-context-driven/SKILL.md
      - skills/testdata/owned-skill-versions.json
  - command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
---

# A profile change keeps the recorded decisions and names what it wrote

The audited candidate makes a profile change with `baseline plan` carry every recorded Setup Manifest decision the new profile still selects. It asks for every missing answer in one refusal, and `baseline update --decision profile=…` refuses and names the plan command without writing. An apply names each Root Backup it writes, while a managed refresh names none. A no-change reconcile names the Setup Snapshot revision the binary carries. A restore creates Claude Skill Links only where the repository already keeps them and rolls them back with the tree.
