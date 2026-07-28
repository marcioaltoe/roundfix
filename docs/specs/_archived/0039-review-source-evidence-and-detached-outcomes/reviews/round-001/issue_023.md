---
source: coderabbit
pr: "40"
round: 1
round_created_at: "2026-07-28T18:11:57Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/review-evidence-and-verification-capacity
head_sha: 2f9f75357d07ae37bd94ddc54579d8e8b8d6ef0b
file: internal/tui/cockpit.go
line: 414
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Uf6Oz,comment:PRRC_kwDOS0qyts7aoLR3
review_hash: 7e7175d646e90745ab2d8847102ee49098d73178d084f29fc8be3e17438afcec
duplicate_of: ""
source_review_id: "4800337236"
source_review_submitted_at: "2026-07-28T17:53:09Z"
---

# Issue 023: _ Functional Correctness_ _ Minor_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**Terminal-to-terminal transitions contradict the documented stickiness invariant.**

The comment on `terminal()` (lines 337-339) states that a stale, duplicated, or replayed event after settlement is "never a regression", but the guard only blocks terminal → non-terminal. A replayed `skipped` event arriving after `settled completed` passes the guard and rewrites `status` to `skipped`, which then flows into `taskStatuses` via `applyTaskJournalStatuses` and renders `[skip] Skipped` for a completed Task. `TestCockpitSpecRunTaskSettlementResistsStaleAndReplayedEvents` replays `started`, so this path is uncovered.

Either make the guard match the comment, or relax the comment to say the last terminal event wins.

<details>
<summary>🛡️ Proposed guard matching the stated invariant</summary>

```diff
 func (model *cockpitModel) advanceTaskPhase(index int, next taskJournalState) {
 	if index < 0 || index >= len(model.taskJournalStates) {
 		return
 	}
-	if model.taskJournalStates[index].phase.terminal() && !next.phase.terminal() {
+	if model.taskJournalStates[index].phase.terminal() {
 		return
 	}
 	model.taskJournalStates[index] = next
 }
```

</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
func (model *cockpitModel) advanceTaskPhase(index int, next taskJournalState) {
	if index < 0 || index >= len(model.taskJournalStates) {
		return
	}
	if model.taskJournalStates[index].phase.terminal() {
		return
	}
	model.taskJournalStates[index] = next
}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/tui/cockpit.go` around lines 406 - 414, Update advanceTaskPhase so
any already-terminal phase ignores subsequent events, including
terminal-to-terminal transitions; only non-terminal entries may advance.
Preserve the existing bounds check and allow normal non-terminal progression,
ensuring replayed terminal events cannot overwrite a settled state.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:06f476828960f8f99430fdd0 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The existing guard allowed a later terminal event to overwrite an
  already-settled or skipped Task. Made every terminal phase sticky and
  extended
  `TestCockpitSpecRunTaskSettlementResistsStaleAndReplayedEvents` with both
  settled-to-skipped and skipped-to-settled replays. Focused evidence:
  `GOCACHE=/private/tmp/roundfix-batch-002-gocache rtk go test ./internal/tui -run '^TestCockpitSpecRunTaskSettlementResistsStaleAndReplayedEvents$' -count=1`
  passed; the combined daemon, Run Event, and TUI package check passed 395
  tests. The Daemon owns the configured `make verify` run after this Agent
  turn.
