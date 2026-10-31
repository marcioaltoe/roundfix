---
schema: roundfix/archive-record/v1
spec: 0332-the-delivery-queue-regenerates-the-derived-records-of-an-archived-spec
title: The Delivery Queue regenerates the derived records of an archived Spec
status: archived
created: "2026-10-30"
archived: "2026-10-31"
disposition: partial
source: docs/specs/0332-the-delivery-queue-regenerates-the-derived-records-of-an-archived-spec
source_revision: c6f68420e1e3c5a3f3e4c8a28dc20f3f8c41728e
qa_task: task_03
qa_report: qa-report-2026-10-31.md
qa_verdict: partial
unproven:
  - the queue log of the first such item in cycle AS (Specs 0333 to 0335), recorded by the operator in the queue interventions log.
adrs:
  - ADR-0345
  - ADR-0149
sources:
  - 2026-10-27-the-queue-cannot-regenerate-the-coverage-record-of-an-archived-spec.md
regeneration:
  - command: go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1
    outputs:
      - docs/references/coverage-record.json
promoted: []
pull_request: "712"
---

# The Delivery Queue regenerates the derived records of an archived Spec

At fe13e8fd953580d3b77b190330da63d195de6fce, all seven Project Config regenerations passed on an archived item’s Archive Record, and real local-Git queue recovery committed a sanctioned coverage-record conflict while refusing and rolling back the same conflict without its sanction. Four independent sabotage copies produced all six expected test failures, and the history audit found only declared, authorized edits. The gate closes as a qualifying partial with no findings: Success Metric 5 remains unproven until the first matching real conflict after release is recorded by the operator, and no Pull Request exists yet. The offline archive plan succeeded without contacting a judge; reusable knowledge already lives in ADR-0345, the glossary and the repository rule, so the adopted incident needs no promotion.
