---
schema: roundfix/archive-record/v1
spec: 0359-the-delivery-queue-waits-for-every-check-the-head-triggered
title: The Delivery Queue waits for every check the head triggered
status: archived
created: "2026-11-06"
archived: "2026-11-06"
disposition: partial
source: docs/specs/0359-the-delivery-queue-waits-for-every-check-the-head-triggered
source_revision: 5a2ee15ac58c02b6097a058e605354c0d37fb4a2
qa_task: task_05
qa_report: qa-report-2026-11-06.md
qa_verdict: partial
unproven:
  - the Fluxus console log and pull request timeline after the release that carries this Spec, reported back through ai-memory
adrs:
  - ADR-0369
sources:
  - 2026-11-05-fluxus-delivery-queue-merged-before-checks.md
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
pull_request: "771"
type: fix
---

# The Delivery Queue waits for every check the head triggered

The local checking-stage replay now waits for unregistered workflow jobs and missing Expected Workflows, names failed members and their links, and exposes pending members in the console log and status. All runnable rows passed, including the seven sabotage probes and the exact public CLI transcript. The released Fluxus journey remains unproven under its pre-run declaration, and the Pull Request journey awaits publication, so this report records a qualifying declared partial. Durable knowledge is already recorded in accepted ADR-0369, the glossary and the guides.
