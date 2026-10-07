---
schema: roundfix/archive-record/v1
spec: 0020-run-browser
title: Run Browser
status: archived
created: "2026-07-07"
archived: "2026-07-07"
disposition: pass
source: docs/history/specs/0020-run-browser
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: qa-report-2026-07-07.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
---

# Run Browser

Run discovery shipped as a flat, unbounded listing: the Attach picker dumps every Run the Run Database has ever recorded — 43 lines and growing — as a numbered prompt, each row showing only id, state, kind, and target. A user cannot tell which repository a Run belongs to, which Agent ran it, when it started, or how long it took; the numbers are ephemeral picker positions that look like stable arguments (`roundfix attach 41` fails confusingly); and the history drowns the one question that matters at a terminal: what is running now. The Run Browser replaces the prompt with a navigable TUI over the same data, defaults every surface to Active Runs, and enriches each row with the context a human needs to pick.
