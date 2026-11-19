---
schema: roundfix/archive-record/v1
spec: 0409-the-idea-and-prd-question-rules-that-block-the-goal
title: The Idea and the PRD question the rules that block the goal
status: archived
created: "2026-11-18"
archived: "2026-11-19"
disposition: partial
source: docs/specs/0409-the-idea-and-prd-question-rules-that-block-the-goal
source_revision: 061eef233cb3b1a59a00934ee1e8ae8f6ce9943c
qa_task: task_04
qa_report: qa-report-2026-11-19.md
qa_verdict: partial
unproven:
  - the operator's session evidence after the release, recorded session evidence.
  - the count of the measurement in `references/2026-11-18-rules-named-in-recent-plans.md` repeated by the operator at a cycle close, recorded session evidence.
adrs:
  - ADR-0437
  - ADR-0436
sources:
  - 2026-11-18-idea-and-prd-question-rules-that-block-the-goal.md
  - 2026-11-18-maintainer-answers.md
  - 2026-11-18-rules-named-in-recent-plans.md
  - 2026-11-19-prototype-and-governed-paths.md
regeneration:
  - command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
    outputs:
      - internal/baseline/assets/modules/spec-workflow.json
      - internal/baseline/module-versions.json
  - command: make baseline-digests
  - command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
    outputs:
      - docs/agents/setup-context.json
      - docs/agents/spec-routing.md
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - .agents/skills/brainstorming/SKILL.md
      - .agents/skills/roundfix/SKILL.md
      - .agents/skills/write-idea/SKILL.md
      - .agents/skills/write-prd/SKILL.md
      - .agents/skills/write-techspec/SKILL.md
      - skills/brainstorming/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/testdata/owned-skill-versions.json
      - skills/write-idea/SKILL.md
      - skills/write-prd/SKILL.md
      - skills/write-techspec/SKILL.md
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
type: feat
---

# The Idea and the PRD question the rules that block the goal

An author of an Idea or a PRD that finds a recorded decision in the way of the best solution has two routes today, and both are silent: drop the solution, or keep it and leave the rule for later. The maintainer learns that a rule decided the scope one, two or three Specs afterwards. After this Spec, in brainstorming, the Idea and the PRD the author asks the user one question about each rule that blocks the goal, with three answers: keep, change or retire. The rule binds until the user answers, the artifact records the question and the answer in a fixed form, the Spec that asked carries the change, and a rule the user keeps says why it stays. The Baseline delivers the rule to every repository that adopts it, and the TechSpec author asks the same question before a design contradicts an ADR.
