---
schema: roundfix/archive-record/v1
spec: 0274-authoring-guards-hold-the-plan-the-glossary-and-the-prd
title: Authoring guards hold the plan, the glossary and the PRD to their rules
status: archived
created: "2026-10-16"
archived: "2026-10-17"
disposition: pass
source: docs/specs/0274-authoring-guards-hold-the-plan-the-glossary-and-the-prd
source_revision: a1f1f56814eed9b26817690bfb4949183d96e9ec
qa_task: task_05
qa_report: qa-report-2026-10-17-01.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0280
sources:
  - 2026-10-14-spec-check-does-not-flag-a-guide-task-without-the-coverage-record.md
  - 2026-10-14-the-implement-reference-states-the-qa-secret-refusal-in-another-form.md
  - 2026-10-16-the-glossary-is-not-held-to-the-rule-it-applies-to-specs.md
regeneration:
  - command: make skills-sync
promoted: []
relinked:
  - docs/backlog/2026-10-16-spec-check-does-not-flag-an-owned-skill-edit-without-its-version-record.md
---

# Authoring guards hold the plan, the glossary and the PRD to their rules

Four authoring rules had no check, and each cost work. On 2026-10-14 a Task of Spec 0268 edited three user guides while its authorization's Limits forbade the Behavior Surface Record; its Verification failed twice, the Run ended Unresolved, and `roundfix spec check --strict` had been clean on the plan (the first Backlog Entry (`docs/specs/0274-authoring-guards-hold-the-plan-the-glossary-and-the-prd/references/2026-10-14-spec-check-does-not-flag-a-guide-task-without-the-coverage-record.md@123f6524f2ea8b9a49655ba5fe0972f3d409a6b3`)). The `roundfix` skill and the `implement` user guide state the QA secret refusal in a form the Daemon does not print (the second Backlog Entry (`docs/specs/0274-authoring-guards-hold-the-plan-the-glossary-and-the-prd/references/2026-10-14-the-implement-reference-states-the-qa-secret-refusal-in-another-form.md@123f6524f2ea8b9a49655ba5fe0972f3d409a6b3`)).
