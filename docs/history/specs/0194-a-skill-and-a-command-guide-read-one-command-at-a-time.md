---
schema: roundfix/archive-record/v1
spec: 0194-a-skill-and-a-command-guide-read-one-command-at-a-time
title: A skill and a command guide read one command at a time
status: archived
created: "2026-09-30"
archived: "2026-10-01"
disposition: pass
source: docs/history/specs/0194-a-skill-and-a-command-guide-read-one-command-at-a-time
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_05
qa_report: qa-report-2026-10-01.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0187
sources: []
regeneration:
  - command: make skills-sync
  - command: make baseline-digests
promoted: []
pull_request: "307"
delivery_commit: 9456ee24d6b173903635f38c0598a56e91994d52
---

# A skill and a command guide read one command at a time

The Roundfix Skill is one file, `.agents/skills/roundfix/SKILL.md`: 141,699 bytes, 2,670 lines and 33 sections on 2026-09-30. The command reference, `docs/user-guide/commands.md`, is 79,467 bytes. Two costs were measured in the two days before this Spec:
