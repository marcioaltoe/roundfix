---
source: coderabbit
pr: "27"
round: 1
round_created_at: "2026-07-15T15:54:46Z"
status: duplicated
head_repository: marcioaltoe/roundfix
head_branch: ma/dogfood-findings-remediation
head_sha: 43147cdff5f36ec1ac2bf276c3747400474d3fab
file: internal/agent/acpx_runner.go
line: 547
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6RJ942,comment:PRRC_kwDOS0qyts7V5tZU
review_hash: d850d8eda8a619f9bff33d9a9455720bcfd02797f54d47ec27ab4b354a61d379
duplicate_of: /Users/marcio/dev/roundfix/docs/specs/_reviews/pr-27/round-002/issue_001.md
terminal_reason: 'Agent failed: Agent Batch failed after acpx exited with code 1: agent/protocol error'
source_review_id: "4705718496"
source_review_submitted_at: "2026-07-15T15:38:25Z"
---




# Issue 001: _ Functional Correctness_ _ Minor_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**Only suggest clearing reasoning effort for reasoning-related failures.**

`ModelNotAdvertisedError.RecoveryAction()` passes `includeReasoning=false`, but this block still suggests changing `reasoning_effort`, which cannot resolve a rejected model.

<details>
<summary>Proposed correction</summary>

```diff
-	if runtime != "" {
+	if includeReasoning && runtime != "" {
 		message += fmt.Sprintf(`, or set runtimes.%s.reasoning_effort "" when the model manages reasoning`, runtime)
 	}
```
</details>





As per coding guidelines, “Errors must name the failed operation and the next useful action when known.”

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
func selectionRecoveryAction(runtime string, includeReasoning bool) string {
	runtime = strings.TrimSpace(runtime)
	message := "update the ACP Runtime or adapter"
	if includeReasoning {
		message += ", choose supported Agent Model and Default Reasoning Effort values, choose an advertised Agent Model"
	} else {
		message += ", choose an advertised Agent Model"
	}
	message += ", or pass a one-Run --model override"
	if includeReasoning {
		message += " with --reasoning-effort when needed"
	}
	if includeReasoning && runtime != "" {
		message += fmt.Sprintf(`, or set runtimes.%s.reasoning_effort "" when the model manages reasoning`, runtime)
	}
	return message
}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/agent/acpx_runner.go` around lines 532 - 547, Update
selectionRecoveryAction so the runtime-specific reasoning_effort suggestion is
appended only when includeReasoning is true. Keep the existing
model-advertisement recovery guidance for includeReasoning=false, including
ModelNotAdvertisedError.RecoveryAction(), without suggesting reasoning changes.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:129b67ac44f192fc8010455f -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `UNREVIEWED`
- Notes:
