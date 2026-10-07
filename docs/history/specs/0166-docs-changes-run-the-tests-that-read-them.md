---
schema: roundfix/archive-record/v1
spec: 0166-docs-changes-run-the-tests-that-read-them
title: Docs changes run the tests that read them
status: archived
created: "2026-09-24"
archived: "2026-09-25"
disposition: pass
source: docs/history/specs/0166-docs-changes-run-the-tests-that-read-them
source_revision: 31c30719398d5c2e43af79d5868b1268fc495f10
qa_task: task_02
qa_report: qa-report-2026-09-25.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "248"
delivery_commit: 7a9b6ec635482440a980509454b5d98d4b1ace5a
---

# Docs changes run the tests that read them

Spec 0155 gave pull requests and Runs a selective gate, `make verify-changed`. Its classifier sends a path under `docs/**`, and Markdown at the repository root, to no test set. Tests read those files: `TestDurableTableLifecyclePolicyCoversEveryTable` reads `docs/user-guide/run-database-lifecycle.md`, `TestBaselineExamplesParse` reads `README.md`, and the archive corpus tests read `docs/history/specs`. A change to one of them passes the selective gate and fails only in the complete `make verify` that pushes to `main` still run. Spec 0155 recorded this limit and carried it here.
