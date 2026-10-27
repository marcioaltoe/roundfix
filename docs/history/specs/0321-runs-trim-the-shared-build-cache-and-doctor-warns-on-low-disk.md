---
schema: roundfix/archive-record/v1
spec: 0321-runs-trim-the-shared-build-cache-and-doctor-warns-on-low-disk
title: Runs trim the shared build cache, and doctor warns on low disk
status: archived
created: "2026-10-26"
archived: "2026-10-26"
disposition: partial
source: docs/specs/0321-runs-trim-the-shared-build-cache-and-doctor-warns-on-low-disk
source_revision: cf916856757cddfad02dfe46cd702a63da1af6ce
qa_task: task_04
qa_report: qa-report-2026-10-26.md
qa_verdict: partial
unproven:
  - the operator's queue records after the release that carries this Spec, recorded session evidence.
adrs:
  - ADR-0334
  - ADR-0327
sources:
  - 2026-10-26-the-shared-build-cache-fills-the-disk-between-manual-gc-runs.md
regeneration: []
promoted: []
---

# Runs trim the shared build cache, and doctor warns on low disk

The Archive Record can carry a qualifying declared partial for automatic Implement Run and repository-gate Shared Build Cache trims and Doctor disk-space warnings: focused tests, both cross-platform builds, GC byte comparisons, documentation and scope audits passed, and all nine sabotage probes failed as intended. The post-release operator-day criterion remains unproven, the Go and Bazel sources were kept offline by authorization, and no Pull Request or Merge-Ready acceptance exists. The offline archive plan succeeded with the judge skipped; durable guidance already lives in CONTEXT.md, ADR-0334 and the user guides, so no additional promotion is proposed.
