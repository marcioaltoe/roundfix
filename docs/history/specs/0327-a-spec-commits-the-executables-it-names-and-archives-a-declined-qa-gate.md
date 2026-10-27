---
schema: roundfix/archive-record/v1
spec: 0327-a-spec-commits-the-executables-it-names-and-archives-a-declined-qa-gate
title: A Spec commits the executables it names, settle says what it removes, and archive honours a declined QA gate
status: archived
created: "2026-10-26"
archived: "2026-10-27"
disposition: partial
source: docs/specs/0327-a-spec-commits-the-executables-it-names-and-archives-a-declined-qa-gate
source_revision: 520f7ceffb7c8d9a4688daa692e51ada8244b55b
qa_task: task_04
qa_report: qa-report-2026-10-27.md
qa_verdict: partial
unproven:
  - Quaestor's reply through ai-memory to the release message the operator sends, recorded session evidence.
adrs:
  - ADR-0340
sources:
  - 2026-10-26-a-spec-cannot-commit-its-executables-settle-exits-128-and-archive-refuses-a-declined-qa.md
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - .agents/skills/archive-spec/SKILL.md
      - .agents/skills/roundfix/SKILL.md
      - .agents/skills/write-tasks/SKILL.md
      - skills/archive-spec/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/testdata/owned-skill-versions.json
      - skills/write-tasks/SKILL.md
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
pull_request: "700"
---

# A Spec commits the executables it names, settle says what it removes, and archive honours a declined QA gate

The audited candidate preserves authorized executable Task output and refuses unlisted executables, including inside new directories; declined-QA Specs archive without an override, qualifying declines refuse an override, and settle reports worktrees it actually removed. Every runnable row passed and all eight isolated sabotages were detected.
