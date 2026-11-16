---
schema: roundfix/archive-record/v1
spec: 0397-each-phase-of-a-spec-is-recorded-and-starts-in-a-fresh-session
title: Each phase of a Spec is recorded and starts in a fresh session
status: archived
created: "2026-11-16"
archived: "2026-11-16"
disposition: partial
source: docs/specs/0397-each-phase-of-a-spec-is-recorded-and-starts-in-a-fresh-session
source_revision: c47d7802d3452763f076d771c0f8d09172922dc1
qa_task: task_05
qa_report: qa-report-2026-11-16.md
qa_verdict: partial
unproven:
  - the operator's usage reports and Run listings after the release, recorded session evidence.
adrs:
  - ADR-0408
  - ADR-0407
  - ADR-0401
sources:
  - 2026-11-14-jarvis-long-orchestration-consumes-the-weekly-claude-limit.md
  - 2026-11-15-maintainer-answers.md
  - 2026-11-15-request-read-by-the-idea-stage-judge.md
  - 2026-11-15-usage-by-context-and-session-identity-measurement.md
  - 2026-11-16-delivery-speed-and-plan-quality-baseline.md
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - .agents/skills/roundfix/SKILL.md
      - .agents/skills/write-idea/SKILL.md
      - .agents/skills/write-prd/SKILL.md
      - .agents/skills/write-tasks/SKILL.md
      - .agents/skills/write-techspec/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/testdata/owned-skill-versions.json
      - skills/write-idea/SKILL.md
      - skills/write-prd/SKILL.md
      - skills/write-tasks/SKILL.md
      - skills/write-techspec/SKILL.md
  - command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
    outputs:
      - internal/baseline/assets/modules/autonomous-work.json
      - internal/baseline/module-versions.json
  - command: make baseline-digests
  - command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
    outputs:
      - docs/agents/autonomous-work.md
      - docs/agents/setup-context.json
  - command: go run -buildvcs=false ./cmd/roundfix skills install --target project
    outputs:
      - .claude/agents/roundfix-idea-author.md
      - .claude/agents/roundfix-prd-author.md
      - .claude/agents/roundfix-researcher.md
      - .claude/agents/roundfix-tasks-author.md
      - .claude/agents/roundfix-techspec-author.md
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
promoted:
  - docs/references/2026-11-15-usage-by-context-and-session-identity-measurement.md
relinked:
  - docs/backlog/2026-11-16-a-new-session-and-each-authoring-phase-start-with-a-large-context.md
pull_request: "883"
type: feat
---

# Each phase of a Spec is recorded and starts in a fresh session

The audited candidate records Supervisor phases and returned Run/review prompts, preserves advisory failure behavior, retains only eligible day files, withholds headers for published skills, and carries the phase-boundary and fix-route rules through synchronized skills and the managed Baseline. Local CLI, backend and document evidence is recorded; the six post-release/live criteria remain unproven and the Pull Request does not yet exist. The offline archive plan exited 0 and names the adopted planning/measurement references for later retention or promotion, including automatic reference relinks; QA made no promotion and contacted no judge. Its output is preserved in offline archive plan (`docs/specs/0397-each-phase-of-a-spec-is-recorded-and-starts-in-a-fresh-session/qa/evidence/qa-report-2026-11-16/archive-plan.txt` (not reachable from any branch or tag)).
