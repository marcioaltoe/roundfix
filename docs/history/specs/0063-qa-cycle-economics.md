---
schema: roundfix/archive-record/v1
spec: 0063-qa-cycle-economics
title: QA cycle economics
status: archived
created: "2026-08-01"
archived: "2026-08-03"
disposition: no-qa
source: docs/history/specs/0063-qa-cycle-economics
source_revision: 436b2d919ede38fc45ec654c236e3f863f304756
qa_task: ""
qa_report: ""
qa_verdict: ""
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "96"
delivery_commit: 0a420490bad320299a0dbd5a78d93c832b6beaf3
---

# QA cycle economics

> **Archived without implementation on 2026-08-03, superseded by Spec 0070.** > > Half its premise died to Spec 0071's measurements. The `148s cold against 5s > warm` figure described Go's *test-result* cache, which `-count=1` disables by > design — cold versus warm compilation was worth about nine seconds, not two > minutes. The Run Worktree it wanted to warm already inherits the ambient > `GOCACHE`; no code forces a cold one. And the twenty-minute cycle it costed > now measures eight to thirteen. > > What stayed true got sharper, not weaker: Core Features 1 and 4 — a static > failure must not suppress the rows it does not implicate, and a report must > distinguish wasted discovery from genuine blocking. Spec 0072's close proved > it at scale, with one governance defect blocking fifteen of twenty-four rows > across four gate executions. Those two features moved into Spec 0070, which > already owned the other half of the same question: what a blocked row means. > > Running cheap detectors ahead of the gate is the one live idea neither Spec > claims. It is left unclaimed here rather than carried as scope nobody owns. > > The text below is the original, unedited.
