---
schema: roundfix/archive-record/v1
spec: 0223-adjustments-the-adopters-asked-for
title: Adjustments the adopters asked for
status: archived
created: "2026-10-04"
archived: "2026-10-04"
disposition: pass
source: docs/history/specs/0223-adjustments-the-adopters-asked-for
source_revision: a30dc847a037fe584812b4471aa7da1701a9ec59
qa_task: task_04
qa_report: qa-report-2026-10-04.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0228
sources:
  - 2026-10-03-branch-prefix-is-a-required-decision.md
  - 2026-10-03-a-reopened-qa-gate-has-no-path-for-a-report-without-inputs.md
  - 2026-09-30-release-plan-does-not-report-skill-and-guide-checks.md
regeneration:
  - command: make skills-sync
  - command: make baseline-digests
promoted: []
pull_request: "377"
delivery_commit: ffe1ec1d081ffded76cb0be78c6da7883f5a966f
---

# Adjustments the adopters asked for

Two repositories that adopted Roundfix, and Roundfix's own release step, hit three small walls in the same week. The oraculum maintainer could not leave the branch prefix to the commit-type rule, because the Baseline made the decision mandatory and always rendered a prefix sentence. fluxus found that a QA gate reopened over a report written before row inputs existed had no valid move under the `qa-gate` skill. And `roundfix release plan`, the first command of every release, says nothing about the skills and guides check that the release runbook makes mandatory. This Spec removes each wall without changing what any adopter has already recorded. The fourth request of the same week, relative links that the archive breaks, needs its own grant and is owned by Spec 0226.
