---
schema: roundfix/archive-record/v1
spec: 0053-qa-gate-reachability-and-verdict-semantics
title: QA gate reachability and verdict semantics
status: archived
created: "2026-07-28"
archived: "2026-07-30"
disposition: pass
source: docs/history/specs/0053-qa-gate-reachability-and-verdict-semantics
source_revision: 436b2d919ede38fc45ec654c236e3f863f304756
qa_task: ""
qa_report: qa-report-2026-07-30.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0080
  - ADR-0053
sources: []
regeneration: []
promoted: []
pull_request: "52"
delivery_commit: d930d3ec775252d93f94ea64f8f408d4b7c1b7b3
---

# QA gate reachability and verdict semantics

The QA gate runs in a Run Worktree with no commit or push authority, so any acceptance row whose journey needs a live Pull Request or Final Push is structurally unreachable: it is recorded `blocked`, the verdict caps at `partial`, and the Spec can never archive through the normal loop — Spec 0039 needed a maintainer `qa_override` for exactly this. Meanwhile every failed QA attempt strands a Run Branch holding a superseded report that Branch Integrity Preflight then counts against the next review Run, and the Daemon's QA prompt and the qa-gate Skill still disagree about the report filename. Evidence: [QA gate cannot reach Pull Request journeys](../../findings/2026-07-28-qa-gate-cannot-reach-pull-request-journeys.md), [failed QA Runs strand branches](../../findings/2026-07-28-failed-qa-runs-strand-branches-that-block-review-runs.md), and the naming defect recorded in [same-day QA reruns were ignored](../../findings/2026-07-28-same-day-qa-reruns-are-ignored-by-the-verdict-selector.md).
