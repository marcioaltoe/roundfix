---
source: coderabbit
pr: "87"
round: 1
round_created_at: "2026-08-03T15:34:03Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/qa-task-not-flag
head_sha: a12c1a665c5970773e04c4a145c6b9b0c5a0e686
file: internal/daemon/task_engine_test.go
line: 227
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WBeNV,comment:PRRC_kwDOS0qyts7c2Vzr
review_hash: c15629be4ecb5898dcd3fb98673d10050aed50f59405016219c40df0b1011547
duplicate_of: ""
source_review_id: "4845660382"
source_review_submitted_at: "2026-08-03T15:14:34Z"
---

# Issue 002: _ Maintainability & Code Quality_ _ Minor_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟡 Minor_ | _⚡ Quick win_

**Call `t.Helper()` in `qaPlan` and `declinedQAPlan`.**

Both methods can fail the test through `writeSpecDirAtRootForTestWithQA` and `reloadGraph`. Without `fixture.t.Helper()`, a failure reports a line inside these helpers instead of the calling test. `reloadGraph` already calls it at Line 230.

<details>
<summary>♻️ Proposed fix</summary>

```diff
 func (fixture *taskCycleFixture) qaPlan() TaskPlan {
+	fixture.t.Helper()
 	seeds, gateID := taskSeedsWithQAGate(fixture.seeds)
```

```diff
 func (fixture *taskCycleFixture) declinedQAPlan() TaskPlan {
+	fixture.t.Helper()
 	writeSpecDirAtRootForTestWithQA(fixture.t, fixture.specsRoot, taskCycleSlug, fixture.seeds, qaDeclarationForTest{
```
</details>

As per coding guidelines: "Test helpers must call `t.Helper()`".

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
func (fixture *taskCycleFixture) qaPlan() TaskPlan {
	fixture.t.Helper()
	seeds, gateID := taskSeedsWithQAGate(fixture.seeds)
	writeSpecDirAtRootForTestWithQA(fixture.t, fixture.specsRoot, taskCycleSlug, seeds, qaDeclarationForTest{taskID: gateID})
	fixture.seeds = seeds
	fixture.reloadGraph()
	return fixture.plan()
}

func (fixture *taskCycleFixture) declinedQAPlan() TaskPlan {
	fixture.t.Helper()
	writeSpecDirAtRootForTestWithQA(fixture.t, fixture.specsRoot, taskCycleSlug, fixture.seeds, qaDeclarationForTest{
		declined: true,
		reason:   "this fixture has no behavioral surface",
	})
	fixture.reloadGraph()
	return fixture.plan()
}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/daemon/task_engine_test.go` around lines 212 - 227, Mark both
taskCycleFixture methods qaPlan and declinedQAPlan as test helpers by calling
fixture.t.Helper() at the start of each method, before
writeSpecDirAtRootForTestWithQA or reloadGraph can fail.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:4d0dbbbc088e2bc0a67feb00 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Both `qaPlan` and `declinedQAPlan` now mark themselves with `fixture.t.Helper()` before calling failure-capable helpers. `go test ./internal/daemon -count=1` passed with 171 tests.
