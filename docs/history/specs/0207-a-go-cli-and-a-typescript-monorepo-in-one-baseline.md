---
schema: roundfix/archive-record/v1
spec: 0207-a-go-cli-and-a-typescript-monorepo-in-one-baseline
title: A Go CLI and a TypeScript monorepo in one Baseline
status: archived
created: "2026-09-30"
archived: "2026-10-02"
disposition: pass
source: docs/history/specs/0207-a-go-cli-and-a-typescript-monorepo-in-one-baseline
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_05
qa_report: qa-report-2026-10-02.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0204
  - ADR-0205
  - ADR-0219
sources: []
regeneration:
  - command: make baseline-digests
promoted: []
pull_request: "335"
delivery_commit: 18ef15eb3a42dff1387cf16214f33666b82fb684
---

# A Go CLI and a TypeScript monorepo in one Baseline

The planned argus repository holds a Go command-line tool, a Hono backend and a React frontend in one repository. No built-in Baseline Profile can express it: a profile takes one Setup Snapshot, each snapshot mirrors one upstream skill list, and no upstream list carries both the Go skills and the TypeScript workspace skills, so a profile that selected both sets of modules would fail the catalog check that its setup lists every skill its modules name. A repository-owned profile cannot help either, because it may only narrow a built-in profile. The Go guide also opens with no sentence naming what it governs, so beside a Bun workspace its rules read as the whole repository's.
