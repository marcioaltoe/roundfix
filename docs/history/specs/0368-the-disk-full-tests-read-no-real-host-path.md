---
schema: roundfix/archive-record/v1
spec: 0368-the-disk-full-tests-read-no-real-host-path
title: The disk-full tests read no real host path
status: archived
created: "2026-11-09"
archived: "2026-11-09"
disposition: partial
source: docs/specs/0368-the-disk-full-tests-read-no-real-host-path
source_revision: eb0043c294afe31ad8053e142f17729e8781c9ed
qa_task: task_02
qa_report: qa-report-2026-11-09.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0370
sources:
  - 2026-11-09-the-disk-full-tests-read-a-real-host-path.md
regeneration: []
promoted: []
type: fix
---

# The disk-full tests read no real host path

The disk-full paths now remain inside each test's temporary directory, and the free-space lookup walks past a file in a path prefix while preserving permission errors. Baseline reproduction, both sabotages, original assertions and the scope audit confirm the change; QA is partial solely because no Pull Request or accepted current-candidate review exists yet. The fresh binary's offline archive plan (`docs/specs/0368-the-disk-full-tests-read-no-real-host-path/qa/evidence/qa-report-2026-11-09/archive-plan.txt` (not reachable from any branch or tag)) exited 0, listed the Spec cut and its single adopted Backlog candidate, and skipped the judge. That candidate records the resolved incident rather than new reusable guidance, so no promotion is proposed; the existing guides already describe the restored behavior. No archive mutation was performed.
