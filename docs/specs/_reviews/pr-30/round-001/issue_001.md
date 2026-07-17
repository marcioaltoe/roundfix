---
source: coderabbit
pr: "30"
round: 1
round_created_at: "2026-07-17T02:15:32Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/deterministic-agent-session-cancellation
head_sha: f0fe093de636e05e33e2dd115caba1c923259e50
file: internal/tui/cockpit_test.go
line: 1702
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6RodpN,comment:PRRC_kwDOS0qyts7WkvbJ
review_hash: a6acff1930796370dda727359bb189603b376447eecc9a37e6aacdc20df41584
duplicate_of: ""
source_review_id: "4718882944"
source_review_submitted_at: "2026-07-17T02:15:17Z"
---

# Issue 001: _ Functional Correctness_ _ Minor_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**Wait for the writer before asserting rendered completion.**

The polling loop can exhaust before the released goroutine writes the remaining events. The later completion wait and final tick do not recompute `sawAll`, causing a scheduler-dependent false failure. Move the `<-written` receive immediately after `close(resumeWriter)` and before the render-poll loop.

As per coding guidelines, “Tests must assert observable behavior.”

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/tui/cockpit_test.go` around lines 1700 - 1702, In the test flow
around model.Update and pressKey, receive from the writer completion channel
immediately after close(resumeWriter), before entering the render-poll loop.
Then let the existing completion wait and final tick evaluate the newly written
events so the assertion reflects observable rendered behavior.
```

</details>

<!-- fingerprinting:phantom:poseidon:terra -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:444b4b438df043af43ee4766 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The test released the writer goroutine, then polled rendered output before waiting for the writer to finish. If the polling loop exhausted before the remaining events were committed, the later wait and tick did not recompute `sawAll`, so the test could fail despite the UI rendering correctly after completion.
- Fix: Moved the `written` receive and post-completion tick immediately after `close(resumeWriter)`, before the render-poll loop, so the assertion waits on the real writer completion signal instead of scheduler timing.
- Evidence: `rtk go test ./internal/tui -run TestOwningCockpitPollsJournalWhileOwnProcessWrites -count=20` passed.
