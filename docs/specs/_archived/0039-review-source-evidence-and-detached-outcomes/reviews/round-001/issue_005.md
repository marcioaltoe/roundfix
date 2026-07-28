---
source: coderabbit
pr: "40"
round: 1
round_created_at: "2026-07-28T18:11:57Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/review-evidence-and-verification-capacity
head_sha: 2f9f75357d07ae37bd94ddc54579d8e8b8d6ef0b
file: internal/cli/cli.go
line: 4300
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6USIpd,comment:PRRC_kwDOS0qyts7aUVDs
review_hash: 4986be8bde3893852931a5753ef976a0460090a2913326c7d9ec840d6f5ef7dc
duplicate_of: ""
source_review_id: "4793702781"
source_review_submitted_at: "2026-07-28T04:25:19Z"
---

# Issue 005: _ Data Integrity & Integration_ _ Minor_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟡 Minor_ | _⚡ Quick win_

**Default the receipt route/status before journaling.**

If a `Notifier` returns a zero-value receipt (with or without an error), `Route` stays empty and `Status` may stay empty, producing the event name `outcome_notification_` and the summary `Outcome notification  via .`. Since `Notifier` is an exported interface with a replaceable factory seam, normalize both fields here so the journal keys stay well-formed.



<details>
<summary>🛡️ Proposed fix</summary>

```diff
 	if err != nil {
 		receipt.Status = roundnotify.StatusFailed
 	}
+	if receipt.Status == "" {
+		receipt.Status = roundnotify.StatusSkipped
+	}
+	if receipt.Route == "" {
+		receipt.Route = roundnotify.RouteDisabled
+	}
 	if receipt.CompletedAt.IsZero() {
 		receipt.CompletedAt = time.Now().UTC()
 	}
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
	notifyCtx, cancel := context.WithTimeout(withoutCancelOrBackground(ctx), outcomeNotificationTimeout)
	defer cancel()
	receipt, err := notifier.Notify(notifyCtx, outcomeFromRun(run, terminal))
	if err != nil {
		receipt.Status = roundnotify.StatusFailed
	}
	if receipt.Status == "" {
		receipt.Status = roundnotify.StatusSkipped
	}
	if receipt.Route == "" {
		receipt.Route = roundnotify.RouteDisabled
	}
	if receipt.CompletedAt.IsZero() {
		receipt.CompletedAt = time.Now().UTC()
	}
	journalOutcomeNotificationReceipt(notifyCtx, runStore, run.ID, receipt, err, stderr)
	if err != nil {
		reportOutcomeNotificationFailure(err, stderr)
	}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/cli.go` around lines 4288 - 4300, Normalize the receipt returned
by notifier.Notify before journalOutcomeNotificationReceipt: when Route or
Status is empty, assign the canonical outcome-notification route and default
status used by the roundnotify package. Apply these defaults for zero-value
receipts whether Notify returns an error or not, while preserving the existing
failed-status override for errors.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:41ac395c83f020a555b242d8 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Zero-value notifier receipts are normalized before journaling: notifier errors become failed receipts, empty successful statuses become skipped, empty routes become disabled, and missing completion times are populated. The new zero-value regression and existing receipt tests passed.
