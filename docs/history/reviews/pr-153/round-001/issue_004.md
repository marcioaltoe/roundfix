---
source: coderabbit
pr: "153"
round: 1
round_created_at: "2026-08-10T21:31:46Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/a-proof-that-can-refuse
head_sha: d68e5a2a65875cf1f7a5e9976514c7be60ee5d5d
file: internal/agent/acpx_runner_test.go
line: 1164
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6YBWiW,comment:PRRC_kwDOS0qyts7fswNx
review_hash: 04ce1a9317f48b6e30e8cfd25789028da91d00cf91368b08524e93a55475abe9
duplicate_of: ""
source_review_id: "4900622643"
source_review_submitted_at: "2026-08-10T20:19:13Z"
---

# Issue 004: _ Maintainability & Code Quality_ _ Major_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟠 Major_ | _⚡ Quick win_

**Remove the duplicated exit-code-4 case and the duplicated diagnosis constant.**

`TestDisposableSessionCloseIsNotAppendedWhenTheSessionNeverOpened` (line 1076) and the "missing session exit four" row of `TestMissingSessionIsRecognisedFromBothExitShapes` (line 1118) use the same fixture (`"sessions ensure": 2`, `"sessions close": 4`), assert the same diagnosis, and assert the same single warning. The `wantDiagnosis` constant is repeated verbatim at lines 1096 and 1151.

Keep the table-driven test, delete the standalone test, and lift the diagnosis into one shared constant. This satisfies the `dupl` linter and the duplication threshold in the coding guidelines.

As per coding guidelines: "avoid duplication exceeding the configured 100-token threshold".




<details>
<summary>♻️ Proposed consolidation</summary>

```diff
-func TestDisposableSessionCloseIsNotAppendedWhenTheSessionNeverOpened(t *testing.T) {
-	t.Parallel()
-
-	harness := newFakeACPXHarness(t)
-	harness.setEnv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{
-		"sessions ensure": 2,
-		"sessions close":  4,
-	}))
-	var warnings []string
-	harness.runner.warnf = func(format string, args ...any) {
-		warnings = append(warnings, fmt.Sprintf(format, args...))
-	}
-
-	_, err := harness.runner.ProveExactSelection(context.Background(), ProbeRequest{
-		Runtime: RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.6-sol", ReasoningEffort: "high"},
-		WorkDir: harness.gitRoot,
-	})
-	if err == nil {
-		t.Fatal("expected selection failure")
-	}
-	const wantDiagnosis = `apply Agent Selection "codex"/"gpt-5.6-sol"/"high" during ensure disposable Agent Session without model override: adapter rejected selection: acpx infrastructure error after exit code 2: acpx command failed; recovery: update the ACP Runtime or adapter and retry the exact Agent Selection`
-	if got := err.Error(); got != wantDiagnosis {
-		t.Fatalf("selection diagnosis changed or gained a close error\nwant: %q\ngot:  %q", wantDiagnosis, got)
-	}
-	var cleanupErr *AgentSessionCleanupError
-	if errors.As(err, &cleanupErr) {
-		t.Fatalf("missing disposable session close was appended to the diagnosis: %v", err)
-	}
-	if len(warnings) != 1 || !strings.Contains(warnings[0], "close disposable Agent Session") || !strings.Contains(warnings[0], acpxExitReasonMissingSession) {
-		t.Fatalf("missing disposable session close was not recorded: %#v", warnings)
-	}
-}
+const wantCatalogueEnsureRejectionDiagnosis = `apply Agent Selection "codex"/"gpt-5.6-sol"/"high" during ensure disposable Agent Session without model override: adapter rejected selection: acpx infrastructure error after exit code 2: acpx command failed; recovery: update the ACP Runtime or adapter and retry the exact Agent Selection`
```

Then reference `wantCatalogueEnsureRejectionDiagnosis` inside the table-driven subtest at line 1151.
</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/agent/acpx_runner_test.go` around lines 1076 - 1164, Consolidate the
duplicate missing-session coverage by deleting
TestDisposableSessionCloseIsNotAppendedWhenTheSessionNeverOpened and retaining
the equivalent “missing session exit four” case in
TestMissingSessionIsRecognisedFromBothExitShapes. Define one shared
wantCatalogueEnsureRejectionDiagnosis constant and use it in the table-driven
subtest instead of repeating the diagnosis literal, preserving the existing
warning and cleanup assertions.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:refactor_suggestion -->

<!-- cr-comment:v1:8705eb71b4b2c243c49428ed -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `RESOLVED`
- Notes: Deleted the standalone `TestDisposableSessionCloseIsNotAppendedWhenTheSessionNeverOpened` and lifted the diagnosis into one shared package-level `wantCatalogueEnsureRejectionDiagnosis` constant, referenced by the "missing session exit four" row of `TestMissingSessionIsRecognisedFromBothExitShapes`. The table-driven test still asserts the warning and the absence of a cleanup error. Focused: `go test ./internal/agent -run 'TestMissingSession|TestDisposable|TestUnrelated'` ok.
