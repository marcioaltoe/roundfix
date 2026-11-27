---
schema: roundfix/archive-record/v1
spec: 0417-the-baseline-carries-the-language-rule
title: The Baseline carries the language rule for code, data and products
status: archived
created: "2026-11-27"
archived: "2026-11-27"
disposition: partial
source: docs/specs/0417-the-baseline-carries-the-language-rule
source_revision: 8c11a97a6b02f3a605f85c5a00794bc94b1e92b8
qa_task: task_04
qa_report: qa-report-2026-11-27-01.md
qa_verdict: partial
unproven:
  - The sweep of the adopters' checkouts and Secondbrain mirrors that a session of this repository repeats 30 days after the release, recorded as a Backlog Entry when a target is missed.
  - The outbox of this project in ai-memory after the release, which the session that sends the messages lists in its report.
adrs:
  - ADR-0141
  - ADR-0118
  - ADR-0186
sources:
  - 2026-11-24-the-baseline-carries-the-language-rule-for-code-data-and-products.md
  - 2026-11-27-characterization.md
regeneration:
  - command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
    outputs:
      - internal/baseline/assets/modules/context-workflow.json
      - internal/baseline/assets/modules/core.json
      - internal/baseline/module-versions.json
  - command: make baseline-digests
    outputs:
      - internal/baseline/assets/profiles/standard-typescript-monorepo.json
      - internal/baseline/testdata/catalog.diagnostics.golden.json
      - internal/baseline/testdata/catalog.digest
      - internal/baseline/testdata/catalog.normalized.json
  - command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text --decision language.reader=English
    outputs:
      - docs/agents/agent-instructions.md
      - docs/agents/domain.md
      - docs/agents/setup-context.json
  - command: make skills-sync
    outputs:
      - .agents/skills/qa-gate/SKILL.md
      - .agents/skills/roundfix/SKILL.md
      - .agents/skills/roundfix/references/baseline.md
      - .agents/skills/roundfix/references/review.md
      - .agents/skills/setup-context-driven/SKILL.md
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - .agents/skills/qa-gate/SKILL.md
      - .agents/skills/roundfix/SKILL.md
      - .agents/skills/setup-context-driven/SKILL.md
      - skills/qa-gate/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/setup-context-driven/SKILL.md
      - skills/testdata/owned-skill-versions.json
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
  - command: go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1
    outputs:
      - docs/references/coverage-record.json
promoted: []
pull_request: "968"
type: feat
---

# The Baseline carries the language rule for code, data and products

The Language Rule, required Reader Language, source-bearing Code Language Exception, review instruction, QA instruction and glossary passed the offline assembled gate. This is a qualifying declared partial because the three post-release adopter criteria remain unproven and no Pull Request exists. The offline archive plan exited 0 and recorded its cut and candidates without calling the judge or archiving anything; ADR-0454, the glossary and guides already hold the reusable decision, so no additional promotion was made. archive-plan (`docs/specs/0417-the-baseline-carries-the-language-rule/qa/evidence/qa-report-2026-11-27-01/archive-plan.txt` (not reachable from any branch or tag))
