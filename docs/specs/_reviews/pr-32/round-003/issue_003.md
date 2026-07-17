---
source: coderabbit
pr: "32"
round: 3
round_created_at: "2026-07-17T14:20:26Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/release-plan-and-agent-selection-profiles
head_sha: 204bbd00fbc648be0df0b8bf2f883b9e2dc490c8
file: internal/store/agent_selection_test.go
line: 271
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Ry2aS,comment:PRRC_kwDOS0qyts7Wzc-k
review_hash: 52319484bd991cd50224394439d64001fe9b6852a4ae568424aff22613c6f436
duplicate_of: ""
source_review_id: "4723452116"
source_review_submitted_at: "2026-07-17T14:16:02Z"
---

# Issue 003: _ Maintainability & Code Quality_ _ Major_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟠 Major_ | _⚡ Quick win_

**Assert the failed lifecycle transition’s persisted result.**

These calls only verify that no error is returned, so regressions that omit the failed status, reason, or fallback event would pass.

<details>
<summary>Proposed assertions</summary>

```diff
 	if _, err := runStore.AppendAgentSelectionAttempt(ctx, failed); err != nil {
 		t.Fatalf("append failed update: %v", err)
 	}
+	qaAttempts, err := runStore.AgentSelectionAttemptsForScope(
+		ctx, run.ID, AgentSelectionScopeQA, "qa:default",
+	)
+	if err != nil {
+		t.Fatalf("read QA attempts: %v", err)
+	}
+	if len(qaAttempts) != 1 ||
+		qaAttempts[0].Status != AgentSelectionStatusFailed ||
+		qaAttempts[0].ReasonCode != "runtime_unavailable" {
+		t.Fatalf("expected one persisted failed QA attempt, got %#v", qaAttempts)
+	}
+
+	events, err = runStore.RunEventsAfter(ctx, run.ID, 0, 10)
+	if err != nil {
+		t.Fatalf("read Run Events: %v", err)
+	}
+	if events[len(events)-1].Event.Kind != runevent.KindDaemonAgentSelectionFallback {
+		t.Fatalf("expected fallback event, got %#v", events)
+	}
```
</details>

As per coding guidelines, “Tests must assert observable behavior and must not rely on production-only hooks.”

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
	failed := AgentSelectionAttemptRequest{
		RunID: run.ID, ScopeKind: AgentSelectionScopeQA, ScopeID: "qa:default",
		Category: "qa", ProfileSource: "built-in", Attempt: 1,
		SelectionRole: AgentSelectionRolePreferred, Runtime: "codex", Model: "gpt-5.6-terra",
		Status: AgentSelectionStatusAttempting,
	}
	if _, err := runStore.AppendAgentSelectionAttempt(ctx, failed); err != nil {
		t.Fatalf("append qa attempting: %v", err)
	}
	failed = withAttemptOverride(failed, func(req *AgentSelectionAttemptRequest) {
		req.Status = AgentSelectionStatusFailed
		req.ReasonCode = "runtime_unavailable"
		req.Reason = "runtime failed before start"
	})
	if _, err := runStore.AppendAgentSelectionAttempt(ctx, failed); err != nil {
		t.Fatalf("append failed update: %v", err)
	}
	qaAttempts, err := runStore.AgentSelectionAttemptsForScope(
		ctx, run.ID, AgentSelectionScopeQA, "qa:default",
	)
	if err != nil {
		t.Fatalf("read QA attempts: %v", err)
	}
	if len(qaAttempts) != 1 ||
		qaAttempts[0].Status != AgentSelectionStatusFailed ||
		qaAttempts[0].ReasonCode != "runtime_unavailable" {
		t.Fatalf("expected one persisted failed QA attempt, got %#v", qaAttempts)
	}

	events, err = runStore.RunEventsAfter(ctx, run.ID, 0, 10)
	if err != nil {
		t.Fatalf("read Run Events: %v", err)
	}
	if events[len(events)-1].Event.Kind != runevent.KindDaemonAgentSelectionFallback {
		t.Fatalf("expected fallback event, got %#v", events)
	}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/store/agent_selection_test.go` around lines 255 - 271, Strengthen
the failed lifecycle assertions in the test around AgentSelectionAttemptRequest
and AppendAgentSelectionAttempt: after the failed update, retrieve the persisted
attempt and assert its failed status, reason code, and reason, and verify the
expected fallback event was recorded. Keep the assertions based on observable
store results rather than production-only hooks.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:refactor_suggestion -->

<!-- cr-comment:v1:6f6d859658c73fadc3470c16 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The failed lifecycle update only asserted that `AppendAgentSelectionAttempt` returned no error. Strengthened the test to read persisted QA attempt state and verify failed status, reason code, reason text, and the fallback Run Event.

## Resolution

- Updated `TestAgentSelectionAttemptLifecycleUpdatesSameAttempt` to assert observable store state after the failed QA attempt update.
- Evidence:
  - `rtk go test ./internal/cli ./internal/store -run 'TestLoadCommitted|TestRunImplementUsesConfiguredExternalSpecRootEndToEnd|TestRunImplementInteractiveInputListsConfiguredExternalSpecRoot|TestAgentSelectionAttemptLifecycleUpdatesSameAttempt'` — passed.
  - `rtk go test ./internal/cli ./internal/store` — passed.
  - `rtk make verify` — passed.
