---
schema: roundfix/archive-record/v1
spec: 0280-qa-evidence-keeps-its-images-and-a-secret-is-refused-where-it-is-born
title: QA evidence keeps its images, and a secret is refused where it is born
status: archived
created: "2026-10-18"
archived: "2026-10-18"
disposition: partial
source: docs/specs/0280-qa-evidence-keeps-its-images-and-a-secret-is-refused-where-it-is-born
source_revision: e377b1a8818f9a4b7710aa579092ea3d0a098b82
qa_task: task_04
qa_report: qa-report-2026-10-18.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0291
  - ADR-0292
sources:
  - 2026-10-18-the-qa-gate-commits-screenshots-the-mechanical-audit-then-refuses.md
  - 2026-10-18-secret-found-asks-for-a-rewrite-the-retry-of-an-archived-item-refuses.md
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - skills/testdata/owned-skill-versions.json
  - command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
pull_request: "582"
---

# QA evidence keeps its images, and a secret is refused where it is born

The Spec accepts image evidence by content, refuses non-image evidence at its own QA Report commit, refuses newly added Secret Matches at Task settlement, and accepts only the bounded Secret Rewrite before returning an archived delivery item to review. The Archive Record may carry a qualifying partial: seven rows passed, three published-source rows were kept offline by authorization, and the pre-PR row remains unproven. Historical PNG/JPEG/SQLite/PDF inputs independently support the classifier. Reusable rules already live in ADR-0291, ADR-0292, CONTEXT and the shipped guides/skills. The two adopted finding references are incident history rather than additional reusable knowledge; neither needs promotion.
