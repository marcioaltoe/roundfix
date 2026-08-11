---
source: coderabbit
pr: "40"
round: 1
round_created_at: "2026-07-28T18:11:57Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/review-evidence-and-verification-capacity
head_sha: 2f9f75357d07ae37bd94ddc54579d8e8b8d6ef0b
file: internal/runevent/stream.go
line: 193
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Uf6Or,comment:PRRC_kwDOS0qyts7aoLRv
review_hash: 488e31e6f2d012e94c724fda852732d6fdc6bca9e151bb7e90ec544e0f39d78f
duplicate_of: ""
source_review_id: "4800337236"
source_review_submitted_at: "2026-07-28T17:53:09Z"
---

# Issue 022: _ Maintainability & Code Quality_ _ Trivial_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚡ Quick win_

**Extract the per-category projection arms; `ProjectStreamEvent` now exceeds the configured function-length limit.**

The function spans lines 115-284 (~170 lines) after these additions, and the verification and outcome arms alone add 13 near-identical read-then-wrap blocks. Splitting each arm into `projectVerificationRecord` / `projectOutcomeRecord` (and optionally a small `readOptionalString(fields, event, key) (string, error)` shim) keeps the behavior identical while bringing the function back under the limit.

As per coding guidelines: "Keep cyclomatic complexity at or below 13, avoid deeply nested conditionals, limit functions to 120 lines and 80 statements, and avoid duplication exceeding the configured 100-token threshold."

<details>
<summary>♻️ Sketch of the extraction</summary>

```diff
-	case StreamCategoryOutcome:
-		record.Outcome, err = firstPayloadString(fields, "outcome", "state")
-		if err != nil {
-			return StreamRecord{}, false, err
-		}
-		record.Reason, err = optionalPayloadString(fields, "reason")
-		if err != nil {
-			return StreamRecord{}, false, streamPayloadFieldError(event, "reason", err)
-		}
-		// ... six more identical blocks ...
-		record.Summary = outcomeStreamSummary(record.Outcome)
+	case StreamCategoryOutcome:
+		if err := projectOutcomeRecord(&record, fields, event); err != nil {
+			return StreamRecord{}, false, err
+		}
```

</details>






Also applies to: 223-254

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/runevent/stream.go` around lines 174 - 193, Extract the per-category
projection logic from ProjectStreamEvent into focused helpers, including
projectVerificationRecord and projectOutcomeRecord for their respective arms,
while preserving all existing parsing and streamPayloadFieldError behavior.
Reuse a small readOptionalString helper where appropriate to remove repeated
optionalPayloadString read-and-wrap blocks, keeping ProjectStreamEvent within
the configured function-length and complexity limits.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:ad51ffd2d9911d28af1ffce6 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `ProjectStreamEvent` was 170 lines in the current HEAD. Extracted the
  verification and outcome arms into focused helpers and centralized the
  repeated contextual optional-string parsing; the dispatcher is now 87
  lines. Added malformed optional-field coverage to preserve
  `streamPayloadFieldError` context. Focused evidence:
  `GOCACHE=/private/tmp/roundfix-batch-002-gocache rtk go test ./internal/runevent -run '^TestProjectStreamEvent' -count=1`
  passed 13 tests; the combined daemon, Run Event, and TUI package check passed
  395 tests. The Daemon owns the configured `make verify` run after this Agent
  turn.
