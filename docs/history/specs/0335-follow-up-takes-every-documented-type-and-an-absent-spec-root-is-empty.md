---
schema: roundfix/archive-record/v1
spec: 0335-follow-up-takes-every-documented-type-and-an-absent-spec-root-is-empty
title: The follow-up command takes every documented message type, an absent Spec Root counts as empty, and the follow-up guides name the JSON fields
status: archived
created: "2026-10-30"
archived: "2026-10-31"
disposition: partial
source: docs/specs/0335-follow-up-takes-every-documented-type-and-an-absent-spec-root-is-empty
source_revision: d23afdd85d08727fecf8ce192a227118aa4c241f
qa_task: task_03
qa_report: qa-report-2026-10-31.md
qa_verdict: partial
unproven:
  - the release notice to `default/gesttione-solutions-quaestor` and Quaestor's reply through ai-memory.
adrs:
  - ADR-0348
sources:
  - 2026-10-30-the-follow-up-command-rejects-documented-message-types-and-an-empty-spec-root.md
regeneration:
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1
    outputs:
      - .agents/skills/roundfix/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/testdata/owned-skill-versions.json
  - command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
---

# The follow-up command takes every documented message type, an absent Spec Root counts as empty, and the follow-up guides name the JSON fields

Quaestor's run of `roundfix follow-up` on v0.126.0 hit three frictions. The `--messages` file was refused with exit 2 for the types `other` and `chore`, although the reference says the rest of the types are `other` and lists no accepted value. A repository whose last Spec was archived reported its `specs` source `unavailable`, because the archive removes the Spec folder and Git keeps no empty directory, so the Spec Root is gone from every fresh checkout. And the reference named the item fields `class`, `reference` and `requester`, while the JSON carries `priority_name`, `ref` and `requested_by`.
