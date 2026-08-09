---
source: coderabbit
pr: "143"
round: 1
round_created_at: "2026-08-09T01:14:50Z"
status: failed
head_repository: marcioaltoe/roundfix
head_branch: ma/specs-0082-0083
head_sha: 49ad4407ca050772b13a812b9625cf45d940e19e
file: internal/cli/baseline_human_test.go
line: 1743
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XbA2w,comment:PRRC_kwDOS0qyts7e5EBN
review_hash: d3dc8212ac39f525a6012a5836e4c9940281697c1243c94bf61b4fdda6169e6b
duplicate_of: ""
terminal_reason: 'Verification failed: command "make verify" exited with exit status 2; diagnostics: /Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20260809T011807Z_77466ac6469b46dc/verification/batch-001-attempt-2.log'
source_review_id: "4887493512"
source_review_submitted_at: "2026-08-08T00:23:24Z"
---


# Issue 013: _ Stability & Availability_ _ Trivial_ _ Quick win_

## Review Comment

_🩺 Stability & Availability_ | _🔵 Trivial_ | _⚡ Quick win_

**Record the forbidden analyzer call instead of calling `t.Fatal` inside the double.**

`t.Fatal` stops only the goroutine that calls it. If the Baseline workflow ever invokes `Segment` or `Classify` from a worker goroutine, the failure is reported without stopping the test, and the surrounding assertions continue against undefined state. The file already has `countingBaselineSemanticAnalyzer` (Lines 1717-1721), which records calls. Use `t.Errorf` plus a recorded flag and assert the flag in each test.

<details>
<summary>♻️ Proposed goroutine-safe double</summary>

```diff
 type forbiddenBaselineSemanticAnalyzer struct {
 	t *testing.T
+	called bool
 }
 
 func (analyzer *forbiddenBaselineSemanticAnalyzer) Segment(
 	context.Context,
 	baseline.RuleSegmentationSnapshot,
 ) (baseline.RuleSegmentationProposal, error) {
 	analyzer.t.Helper()
-	analyzer.t.Fatal("semantic analyzer Segment called on resolved manifest path")
+	analyzer.called = true
+	analyzer.t.Error("semantic analyzer Segment called on resolved manifest path")
 	return baseline.RuleSegmentationProposal{}, nil
 }
```
</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/baseline_human_test.go` around lines 1723 - 1743, Update
forbiddenBaselineSemanticAnalyzer.Segment and Classify to record each unexpected
invocation with a shared flag and report it using t.Errorf rather than t.Fatal,
preserving goroutine-safe test execution. Add assertions for that recorded flag
in each test using this double, alongside the existing workflow assertions, so
any forbidden call fails the test without continuing on undefined state.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:327f81312be0ad585d4a753a -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Changed t.Fatal to t.Errorf + called flag in forbiddenBaselineSemanticAnalyzer. Added assertion checks for called flag in all three tests using this double (TestBaselineHumanResolvedManifestSkipsPromptsAndAnalyzer, TestHumanBaselineProfileDigestDriftRemainsUpdate, TestHumanBaselinePromptsOnlyForManifestMissingDecision). `rtk go build ./...` passes.
