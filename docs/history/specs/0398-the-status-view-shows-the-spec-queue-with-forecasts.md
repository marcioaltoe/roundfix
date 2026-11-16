---
schema: roundfix/archive-record/v1
spec: 0398-the-status-view-shows-the-spec-queue-with-forecasts
title: The status view shows the Spec queue with a forecast start and end
status: archived
created: "2026-11-16"
archived: "2026-11-16"
disposition: partial
source: docs/specs/0398-the-status-view-shows-the-spec-queue-with-forecasts
source_revision: 320217fe678bf871b0793182cd1105826f4981d0
qa_task: task_05
qa_report: qa-report-2026-11-16.md
qa_verdict: partial
unproven:
  - the operator's repeat of the backtest count on the Run listing and the Git history after 20 deliveries, and the timing and the live report recorded in the session.
adrs:
  - ADR-0409
  - ADR-0410
  - ADR-0353
  - ADR-0407
sources:
  - 2026-11-15-the-status-mode-shows-the-spec-queue-with-forecasts.md
  - 2026-11-16-maintainer-answers.md
  - 2026-11-16-status-view-time-baseline.md
  - 2026-11-16-forecast-backtest.md
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

# The status view shows the Spec queue with a forecast start and end

The Archive Record will carry a qualifying declared partial: every runnable row passed, the four production and release criteria remain unproven under their pre-run declarations, and the pre-PR Pull Request row remains environment-blocked. The offline archive plan (`docs/specs/0398-the-status-view-shows-the-spec-queue-with-forecasts/qa/evidence/qa-report-2026-11-16/OUT-archive-plan.txt` (not reachable from any branch or tag)) exited 0 and listed the Spec cut and five knowledge candidates while skipping the judge.
