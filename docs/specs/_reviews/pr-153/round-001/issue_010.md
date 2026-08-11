---
source: coderabbit
pr: "153"
round: 1
round_created_at: "2026-08-10T21:31:46Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/a-proof-that-can-refuse
head_sha: d68e5a2a65875cf1f7a5e9976514c7be60ee5d5d
file: internal/agent/selection_catalogue_characterization_test.go
line: 132
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6YBWjD,comment:PRRC_kwDOS0qyts7fswOn
review_hash: 9ef695a1dd402c027c72a7d42d9235b300234575885bef5b9bfe4eef5dacf7b6
duplicate_of: ""
source_review_id: "4900622643"
source_review_submitted_at: "2026-08-10T20:19:14Z"
---

# Issue 010: _ Functional Correctness_ _ Minor_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**This fixture fails the catalogue ensure, so it does not prove the adapter fast path.**

`fakeACPXExitBy` keys on the bare command, so `"sessions ensure": 2` fails every `sessions ensure` call. Preflight now issues the override-free catalogue ensure first (`readRuntimeCatalogueWithEvidence`), so that call fails and the model-bearing ensure is never reached. The comment at lines 121-122 claims this test preserves the adapter refusal for the requested `--model`, but the recorded refusal comes from the catalogue ensure instead.

`TestProofKeepsTheAdapterRefusalFastPath` in `internal/agent/selection_assignment_test.go` (line 83) uses `fakeACPXExitByCall` with the key `"sessions ensure model=" + unofferedCodexModel`. Use the same call-scoped key here so only the model-bearing ensure fails.



<details>
<summary>🐛 Proposed fix</summary>

```diff
 	harness := newFakeACPXHarness(t)
-	harness.setEnv(fakeACPXExitBy, mustJSONForTest(t, map[string]int{
-		"sessions ensure": 2,
-	}))
+	harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
+		"sessions show": sessionCapabilitySnapshotFixture(t, "gpt-5.6-sol", []string{"gpt-5.6-sol", "gpt-5.5"}, "reasoning_effort", "medium", []string{"low", "medium", "high"}),
+	}))
+	harness.setEnv(fakeACPXExitByCall, mustJSONForTest(t, map[string]int{
+		"sessions ensure model=" + unofferedCodexModel: 2,
+	}))
 	harness.setEnv(fakeACPXStderrBy, mustJSONForTest(t, map[string]string{
 		"sessions ensure": "Cannot apply --model " + unofferedCodexModel + ": the ACP agent did not advertise that model.\nAvailable models: gpt-5.6-sol, gpt-5.5\n",
 	}))
```

Also assert `modelErr.Err != nil` so the test distinguishes the adapter-owned refusal from the Roundfix-owned membership verdict, which sets `Err` to nil.
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
func TestSelectionCatalogueCharacterizationCodexRefusesAnUnofferedModel(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
		"sessions show": sessionCapabilitySnapshotFixture(t, "gpt-5.6-sol", []string{"gpt-5.6-sol", "gpt-5.5"}, "reasoning_effort", "medium", []string{"low", "medium", "high"}),
	}))
	harness.setEnv(fakeACPXExitByCall, mustJSONForTest(t, map[string]int{
		"sessions ensure model=" + unofferedCodexModel: 2,
	}))
	harness.setEnv(fakeACPXStderrBy, mustJSONForTest(t, map[string]string{
		"sessions ensure": "Cannot apply --model " + unofferedCodexModel + ": the ACP agent did not advertise that model.\nAvailable models: gpt-5.6-sol, gpt-5.5\n",
	}))
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/agent/selection_catalogue_characterization_test.go` around lines 123
- 132, Update TestSelectionCatalogueCharacterizationCodexRefusesAnUnofferedModel
to use fakeACPXExitByCall with the call-scoped key "sessions ensure model=" +
unofferedCodexModel, allowing the override-free catalogue ensure to succeed
while only the model-bearing ensure fails. Also assert that modelErr.Err is
non-nil to verify the adapter-owned refusal rather than the Roundfix-owned
membership verdict.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:8ec297d5ca382aba40a4cccb -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `RESOLVED`
- Notes: Fixed `TestSelectionCatalogueCharacterizationCodexRefusesAnUnofferedModel` to prove the adapter fast path. The override-free catalogue ensure now succeeds via a `sessions show` fixture, and only the model-bearing ensure fails, driven by `fakeACPXExitByCall` keyed `"sessions ensure model=" + unofferedCodexModel`. Added `modelErr.Err != nil` so the test asserts the adapter-owned refusal rather than the Roundfix membership verdict. Focused: `go test ./internal/agent -run 'TestSelectionCatalogueCharacterizationCodexRefusesAnUnofferedModel|TestProofKeepsTheAdapterRefusalFastPath'` ok.
