---
source: coderabbit
pr: "33"
round: 1
round_created_at: "2026-07-21T17:02:02Z"
status: invalid
terminal_reason: "An explicitly empty --reasoning-effort is the documented model-managed selection and is covered by passing CLI regression tests."
head_repository: marcioaltoe/roundfix
head_branch: ma/0041-agent-selection-runtime-readiness
head_sha: 6b48b67ab2154bb40d396befd673ad645a528214
file: internal/cli/cli.go
line: 3008
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6SqIbI,comment:PRRC_kwDOS0qyts7YA3RX
review_hash: df3db2fd78653e766eaf5943ad162cbcc2d264a36ad16427b224b68483b6c20e
duplicate_of: ""
source_review_id: "4747041240"
source_review_submitted_at: "2026-07-21T17:01:00Z"
---

# Issue 001: _ Functional Correctness_ _ Major_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Missing `--reasoning-effort` empty-value check in the "complete override" validation.**

`validateExplicitSelectionFlags` checks `--agent` and `--model` for emptiness but never checks `req.reasoningEffort`. Worse, `req.modelSet &&` on line 3004 is dead code — once `presence.validate()` passes with a non-empty presence, `presence.validate()` already guarantees `agent && model && reasoningEffort` are all true, so `req.modelSet` is always true here. This strongly suggests a missing third check was dropped in favor of a stray duplicate guard.

Concretely, `--agent codex --model gpt-5.6-sol --reasoning-effort= --no-input` passes this validation silently today, letting an empty reasoning effort slip into what's supposed to be a fully-proved one-Run override tuple — contradicting the "fail closed" / exact-tuple guarantee the PR describes.

<details>
<summary>🐛 Proposed fix</summary>

```diff
 	if strings.TrimSpace(req.agent) == "" {
 		return validationError{message: "--agent cannot be empty in a complete one-Run Agent Selection override"}
 	}
-	if req.modelSet && strings.TrimSpace(req.model) == "" {
+	if strings.TrimSpace(req.model) == "" {
 		return validationError{message: "--model cannot be empty in a complete one-Run Agent Selection override"}
 	}
+	if strings.TrimSpace(req.reasoningEffort) == "" {
+		return validationError{message: "--reasoning-effort cannot be empty in a complete one-Run Agent Selection override"}
+	}
 	return nil
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
func validateExplicitSelectionFlags(req commandRequest) error {
	presence := invocationSelectionFlagPresence{
		agent:           req.agentSet,
		model:           req.modelSet,
		reasoningEffort: req.reasoningEffortSet,
	}
	if err := presence.validate(); err != nil {
		return err
	}
	if presence.empty() {
		return nil
	}
	if strings.TrimSpace(req.agent) == "" {
		return validationError{message: "--agent cannot be empty in a complete one-Run Agent Selection override"}
	}
	if strings.TrimSpace(req.model) == "" {
		return validationError{message: "--model cannot be empty in a complete one-Run Agent Selection override"}
	}
	if strings.TrimSpace(req.reasoningEffort) == "" {
		return validationError{message: "--reasoning-effort cannot be empty in a complete one-Run Agent Selection override"}
	}
	return nil
}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/cli.go` around lines 2989 - 3008, Update
validateExplicitSelectionFlags to validate req.reasoningEffort is non-empty
after presence.validate succeeds, replacing the redundant req.modelSet guard
with the missing reasoning-effort check. Preserve the existing validationError
behavior and complete one-Run override requirements for agent, model, and
reasoning effort.
```

</details>

<!-- fingerprinting:phantom:medusa:beignet -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:24ff1a364cc39f0d2c8fd798 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: The complete override requires all three flags to be present, but the
  Default Reasoning Effort value may intentionally be empty to request
  model-managed reasoning. This is the canonical contract in `CONTEXT.md` and
  ADR-0055, and `TestInvocationProfileOverrideRequiresCompleteTuple` plus
  `TestInvocationProfileOverrideParsingPreservesExplicitEmptyReasoning` pass
  with the explicit-empty case. Rejecting it would introduce a regression.
- Focused check: `rtk go test ./internal/cli -run 'Test(InvocationProfileOverrideRequiresCompleteTuple|InvocationProfileOverrideParsingPreservesExplicitEmptyReasoning)$'` — passed (11 tests in 1 package).
