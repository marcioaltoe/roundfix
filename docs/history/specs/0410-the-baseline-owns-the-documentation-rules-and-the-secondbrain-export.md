---
schema: roundfix/archive-record/v1
spec: 0410-the-baseline-owns-the-documentation-rules-and-the-secondbrain-export
title: The Baseline owns the documentation rules and the Secondbrain export
status: archived
created: "2026-11-18"
archived: "2026-11-19"
disposition: partial
source: docs/specs/0410-the-baseline-owns-the-documentation-rules-and-the-secondbrain-export
source_revision: 12df82472c073aeedfbf54789b74c91a8202e0ca
qa_task: task_06
qa_report: qa-report-2026-11-19-01.md
qa_verdict: partial
unproven:
  - the operator's count of the mirror after the release, recorded session evidence.
  - section 1 of the measurement in `references/2026-11-18-mirrored-documentation-measurement.md` repeated by the operator at a cycle close, recorded session evidence.
  - the first adopter's update after the release, reported by its project, recorded session evidence.
  - the Git history of `docs/lessons/` in this repository read by the operator at a cycle close, recorded session evidence.
  - the content of this repository's mirror after the release, read by the operator, recorded session evidence.
adrs:
  - ADR-0438
  - ADR-0439
  - ADR-0440
  - ADR-0248
  - ADR-0278
  - ADR-0412
sources:
  - 2026-11-18-the-baseline-owns-the-documentation-rules-and-the-secondbrain-export.md
  - 2026-11-16-the-docs-folders-of-adopters-have-unclear-jobs.md
  - 2026-11-18-maintainer-answers.md
  - 2026-11-18-mirrored-documentation-measurement.md
  - 2026-11-19-prototype-and-governed-paths.md
regeneration:
  - command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
    outputs:
      - internal/baseline/assets/modules/context-workflow.json
      - internal/baseline/assets/modules/core.json
      - internal/baseline/assets/modules/secondbrain.json
      - internal/baseline/module-versions.json
  - command: make baseline-digests
  - command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
    outputs:
      - docs/agents/agent-instructions.md
      - docs/agents/docs-layout.md
      - docs/agents/secondbrain.md
      - docs/agents/setup-context.json
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - .agents/skills/archive-spec/SKILL.md
      - .agents/skills/roundfix/SKILL.md
      - skills/archive-spec/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/testdata/owned-skill-versions.json
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
  - command: go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1
    outputs:
      - docs/references/coverage-record.json
promoted: []
pull_request: "933"
type: feat
---

# The Baseline owns the documentation rules and the Secondbrain export

Nothing in the Baseline says what a repository publishes to the Secondbrain, so each project exports all of `docs/`, and 1,509 of the 3,534 documentation files in the project mirrors are process. What a project learned the hard way is mostly absent, because the Finding that held the cause is deleted when it closes, and an ADR that a later one changed still reads as current. After this Spec the Baseline states the job of each documentation folder and keeps, in each repository that has the Secondbrain turned on, a standard list of what is published: knowledge, never process. An ADR says in its header whether the knowledge base wants it and which ADRs amend it, a Lesson Learned is written only when there is a lesson, `roundfix baseline update` applies what a command can apply, and `roundfix doctor` reports what departs without failing.
