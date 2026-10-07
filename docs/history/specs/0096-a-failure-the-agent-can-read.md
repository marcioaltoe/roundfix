---
schema: roundfix/archive-record/v1
spec: 0096-a-failure-the-agent-can-read
title: A failure the Agent can read
status: archived
created: "2026-08-12"
archived: "2026-08-16"
disposition: pass
source: docs/history/specs/0096-a-failure-the-agent-can-read
source_revision: 1d1726f5d227731aa229fbc336dccbae9aca1398
qa_task: task_08
qa_report: qa-report-2026-08-16-01.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "166"
delivery_commit: 101c3c9215d751bf2e93c5063cc2a17539fd903a
---

# A failure the Agent can read

A failing Verification can hand the Agent Session nothing at all. The command pattern adopted across this repository redirects output to a file to keep the exit status honest, and the check that then fails prints nothing — so the Daemon captures empty output and the feedback prompt carries the command, the exit status, and no cause. In one measured Task the Agent spent its repair turn rewriting its own task file with a diagnosis it had deduced, which is reasonable behavior given zero information. The same blindness repeats across Runs: a Task that fails twice with an identical assertion is reported as new both times, and a whole Run was spent reproducing a known diagnostic. And when a gate returns more corrective work than the contract's ceiling allows, the contract says the decomposition was wrong and does not say what to do, so the loop stops and asks a human.
