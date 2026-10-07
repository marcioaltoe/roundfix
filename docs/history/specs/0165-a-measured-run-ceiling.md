---
schema: roundfix/archive-record/v1
spec: 0165-a-measured-run-ceiling
title: A measured Run ceiling
status: archived
created: "2026-09-24"
archived: "2026-09-25"
disposition: pass
source: docs/history/specs/0165-a-measured-run-ceiling
source_revision: 31c30719398d5c2e43af79d5868b1268fc495f10
qa_task: task_03
qa_report: qa-report-2026-09-24-01.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "245"
delivery_commit: 85d8d7a86161eff538a72d585420588af9627f6a
---

# A measured Run ceiling

On 2026-09-24 this machine ran five Implement Runs at once while the maintainer's session ran the repository gates. Load average passed 15. The QA gate of Spec 0157 then failed its precondition for a reason unrelated to the Spec: `internal/cli` exceeded `go test`'s ten-minute default inside `make verify`, and `TestRunImplementBudgetExceededPreservesRunWorktreeAndBranch` failed three times out of three, on `main` as well. With three concurrent Runs the same gates passed. Nothing in Roundfix bounds how many Runs a machine starts: `worktree.concurrency` bounds Task Worktrees inside one Run, and the Active Run lock is per checkout.
