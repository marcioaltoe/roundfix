---
schema: roundfix/archive-record/v1
spec: 0313-internal-cli-tests-stop-reading-the-live-documentation
title: internal/cli tests stop reading the live documentation
status: archived
created: "2026-10-24"
archived: "2026-10-25"
disposition: partial
source: docs/specs/0313-internal-cli-tests-stop-reading-the-live-documentation
source_revision: 9be06b2521f7d4f961f2c4f450687011e6f36769
qa_task: task_03
qa_report: qa-report-2026-10-25-01.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0326
sources:
  - 2026-10-24-internal-cli-tests-read-the-live-documentation-trees.md
regeneration: []
promoted: []
pull_request: "675"
---

# internal/cli tests stop reading the live documentation

The thirteen documentation-dependent tests pass under docscontract with unchanged bodies, remain selected as always-relevant contracts, and are absent from the untagged suite. The audited-head result stays cached after guide and skill edits while the current Delivery Base reruns; the full test binary logs no forbidden documentation reads and the suite passes with all five documentation trees removed. All nine deliberate regressions are detected, all 22 leak cases and every guarded package pass, and the scope and glossary audits pass. The read-only offline archive plan skips the judge and leaves Spec files unchanged. This gate records a qualifying partial solely because no Pull Request is open; settlement remains Daemon-owned.
