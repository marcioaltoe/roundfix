---
source: coderabbit
pr: "40"
round: 1
round_created_at: "2026-07-28T18:11:57Z"
status: invalid
terminal_reason: "RunEventsAfter guarantees ordered cursors strictly greater than the supplied cursor, and the proposed assignment guard would not terminate a hypothetical non-advancing full page."
head_repository: marcioaltoe/roundfix
head_branch: ma/review-evidence-and-verification-capacity
head_sha: 2f9f75357d07ae37bd94ddc54579d8e8b8d6ef0b
file: internal/cli/attach.go
line: 407
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Uf6OH,comment:PRRC_kwDOS0qyts7aoLQ5
review_hash: fcf13e7f7a41481c9811fa27bbbdabc1eacfec6b578b324b43e75eddcdbcfd87
duplicate_of: ""
source_review_id: "4800337236"
source_review_submitted_at: "2026-07-28T17:53:09Z"
---

# Issue 018: _ Stability & Availability_ _ Minor_ _ Quick win_

## Review Comment

_🩺 Stability & Availability_ | _🟡 Minor_ | _⚡ Quick win_

**Guard the replay loop against a non-advancing cursor.**

`cursor` is assigned unconditionally from `entry.Cursor`. If a full page ever comes back whose cursors do not exceed the cursor that was passed in, the `len(page) < attachReplayPageSize` exit never fires and this spins forever on a CLI foreground path. The sibling replay in `internal/tui/cockpit.go` (lines 290-293) already guards this with `if entry.Cursor > cursor`; matching it here costs one line.

<details>
<summary>🛡️ Proposed guard</summary>

```diff
 		for _, entry := range page {
-			cursor = entry.Cursor
+			if entry.Cursor > cursor {
+				cursor = entry.Cursor
+			}
 			if entry.Event.Kind != runevent.KindDaemonStatus || len(entry.Event.Payload) == 0 {
 				continue
 			}
 			var payload attachCapacityPayload
 			if json.Unmarshal(entry.Event.Payload, &payload) == nil && payload.recorded() {
 				return payload, true
 			}
 		}
```

</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
func attachRecordedCapacities(ctx context.Context, reader *store.Store, runID string) (attachCapacityPayload, bool) {
	cursor := int64(0)
	for {
		page, err := reader.RunEventsAfter(ctx, runID, cursor, attachReplayPageSize)
		if err != nil {
			return attachCapacityPayload{}, false
		}
		for _, entry := range page {
			if entry.Cursor > cursor {
				cursor = entry.Cursor
			}
			if entry.Event.Kind != runevent.KindDaemonStatus || len(entry.Event.Payload) == 0 {
				continue
			}
			var payload attachCapacityPayload
			if json.Unmarshal(entry.Event.Payload, &payload) == nil && payload.recorded() {
				return payload, true
			}
		}
		if len(page) < attachReplayPageSize {
			return attachCapacityPayload{}, false
		}
	}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/attach.go` around lines 387 - 407, Update
attachRecordedCapacities so cursor advances only when entry.Cursor is greater
than the current cursor, matching the guard in the sibling cockpit replay;
preserve event processing while preventing non-advancing full pages from causing
an infinite loop.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:5e986e6181abe8b37277d4a4 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: `attachRecordedCapacities` consumes the concrete `store.RunEventsAfter` contract, which returns ordered entries with cursors strictly greater than the supplied cursor. The suggested `if entry.Cursor > cursor` assignment alone would still loop forever if that contract were violated by a full non-advancing page, so it neither fixes the claimed failure mode nor adds protection under the real API contract.
