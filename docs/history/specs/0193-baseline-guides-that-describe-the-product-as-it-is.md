---
schema: roundfix/archive-record/v1
spec: 0193-baseline-guides-that-describe-the-product-as-it-is
title: Baseline guides that describe the product as it is
status: archived
created: "2026-09-30"
archived: "2026-09-30"
disposition: pass
source: docs/history/specs/0193-baseline-guides-that-describe-the-product-as-it-is
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_05
qa_report: qa-report-2026-09-30.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0186
sources:
  - 2026-09-30-baseline-guides-contradict-the-shipped-product.md
  - 2026-09-30-the-backend-guide-renders-one-clause-twice.md
regeneration:
  - command: make baseline-digests
promoted: []
pull_request: "298"
delivery_commit: 18e3d7f2cd74cd38198725c1be03728387571c65
---

# Baseline guides that describe the product as it is

The Context-Driven Baseline ships its clauses to every adopter as mandatory guidance, and an Agent obeys them literally. An audit on 2026-09-30 found clauses that describe an older product. The loop clause orders archive before the pre-PR review, while the Delivery Queue reviews first. One clause fails a tooling Task for any undeclared path, which the audit stopped doing. One describes an `execution_approvals` record that no code reads. Four clauses cite this repository's own Spec and ADR numbers, which name other records in an adopter's repository. The backend guide renders one paragraph twice. The adopted entries in [references/_index.md](references/_index.md) list each case with its evidence.
