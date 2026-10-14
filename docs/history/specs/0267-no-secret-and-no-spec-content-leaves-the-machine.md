---
schema: roundfix/archive-record/v1
spec: 0267-no-secret-and-no-spec-content-leaves-the-machine
title: No secret and no Spec content leaves the machine
status: archived
created: "2026-10-13"
archived: "2026-10-14"
disposition: partial
source: docs/specs/0267-no-secret-and-no-spec-content-leaves-the-machine
source_revision: d4f32a9b7913d461c85e529191432d7d563ca208
qa_task: task_05
qa_report: qa-report-2026-10-14.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0273
sources:
  - 2026-10-13-qa-evidence-records-environment-values.md
  - 2026-10-13-a-qa-run-can-send-spec-content-to-the-judge-provider.md
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "530"
---

# No secret and no Spec content leaves the machine

The Spec withholds all three Judge Key Variables from Agent and Verification environments, refuses sample-bearing QA commits and both Roundfix push paths, and documents location-only refusals. At this head all runnable declared checks passed, including the temporary sabotage tests, the offline public archive plan, and zero-match scans of tracked content and Task additions. The remaining no-open-PR row makes this a qualifying partial; current-candidate review and PR controls remain for the delivery owner.
