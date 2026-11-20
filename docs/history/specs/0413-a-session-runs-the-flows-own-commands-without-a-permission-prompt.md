---
schema: roundfix/archive-record/v1
spec: 0413-a-session-runs-the-flows-own-commands-without-a-permission-prompt
title: A session runs the flow's own commands without a permission prompt
status: archived
created: "2026-11-19"
archived: "2026-11-20"
disposition: partial
source: docs/specs/0413-a-session-runs-the-flows-own-commands-without-a-permission-prompt
source_revision: 53ef15f5ac816f90b1c9985653099bd5ad26f031
qa_task: task_05
qa_report: qa-report-2026-11-20-01.md
qa_verdict: partial
unproven:
  - the operator's walk in Claude Code after the release, with the count of prompts before and after applying, kept as a reference document, recorded session evidence.
  - the operator's walk in Codex after the release, kept in the same reference document, recorded session evidence.
  - the operator's count, at a cycle close, of the Backlog Entries, Findings and messages that name that cause, recorded session evidence.
adrs:
  - ADR-0444
  - ADR-0445
  - ADR-0446
sources:
  - 2026-11-19-maintainer-answers.md
  - 2026-11-19-the-flows-commands-and-the-harness-permission-rules.md
  - 2026-11-19-harness-rule-measurements.md
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

# A session runs the flow's own commands without a permission prompt

The gate supports the Flow Command List's preview, digest-bound apply, formatting-preserving removal, safe refusal paths, advisory reports and governed-file contract at this candidate, and resolves the previous omission of `git -C` in text previews. All nine sabotages are detected by their named assertions.
