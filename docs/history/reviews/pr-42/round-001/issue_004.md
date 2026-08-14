---
source: coderabbit
pr: "42"
round: 1
round_created_at: "2026-07-29T02:33:36Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/claude-adapter-standardization
head_sha: 7155ba4d2ef353257a1bacf697027202d4750492
file: internal/cli/doctor.go
line: 187
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6UmtXq,comment:PRRC_kwDOS0qyts7ayC0D
review_hash: 38df9c91d9f0a26d61d0e8a3dd8bc93a7e7e6819106652a195bf6e67a40d91ec
duplicate_of: ""
source_review_id: "4803488138"
source_review_submitted_at: "2026-07-29T02:32:46Z"
---

# Issue 004: _ Functional Correctness_ _ Major_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Only the first failing runtime's `NextAction`/`Err` survive aggregation.**

`result.NextAction`/`result.Err` are set once, guarded by `if result.NextAction == ""` (Line 178). If two runtimes fail (e.g., both Claude and Codex adapters are stale after this PR's version bumps), the `Detail` string correctly shows both classifications, but only the first failing runtime's remediation command is surfaced — the second failure's `NextAction`/`Err` are silently dropped. A user following the printed `next:` guidance would fix one adapter, rerun `doctor`, and only then discover the second failure. This scenario isn't covered by the tests in `doctor_test.go`, which only exercise single-runtime failures.

As per coding guidelines, use `errors.Join` to combine independent errors when appropriate; the same aggregation pattern used for `Detail` should be applied to `NextAction`.

<details>
<summary>🐛 Proposed fix to preserve remediation guidance for every failing runtime</summary>

```diff
 	result := CheckResult{
 		Name:   HealthCheckAdapter,
 		Status: CheckStatusOK,
 	}
 	details := make([]string, 0, len(runtimes))
+	nextActions := make([]string, 0, len(runtimes))
+	var failures error
 	for _, runtime := range runtimes {
 		runtimeResult := checker.Adapter(ctx, runtime)
 		detail := strings.TrimSpace(runtimeResult.Detail)
 		if detail == "" {
 			detail = string(runtimeResult.Status)
 		}
 		if runtimeResult.Status == CheckStatusFailed {
 			result.Status = CheckStatusFailed
 			if classification := doctorAdapterFailureClassification(runtimeResult.Err); classification != "" {
 				detail += "; classification: " + classification
 			}
-			if result.NextAction == "" {
-				result.NextAction = strings.TrimSpace(runtimeResult.NextAction)
-				result.Err = runtimeResult.Err
-			}
+			if action := strings.TrimSpace(runtimeResult.NextAction); action != "" {
+				nextActions = append(nextActions, action)
+			}
+			failures = errors.Join(failures, runtimeResult.Err)
 		}
 		details = append(details, runtime.ID+": "+detail)
 	}
 	result.Detail = strings.Join(details, " | ")
+	result.NextAction = strings.Join(nextActions, " && ")
+	result.Err = failures
 	return result
 }
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
func doctorAdapterCheck(ctx context.Context, checker HealthChecker, runtimes []agent.RuntimeSpec, runtimeErr error) CheckResult {
	if runtimeErr != nil {
		return CheckResult{
			Name:   HealthCheckAdapter,
			Status: CheckStatusFailed,
			Detail: runtimeErr.Error(),
			Err:    runtimeErr,
		}
	}
	if len(runtimes) == 0 {
		err := errors.New("effective required Agent Selection Profiles reference no ACP Runtime")
		return CheckResult{
			Name:   HealthCheckAdapter,
			Status: CheckStatusFailed,
			Detail: err.Error(),
			Err:    err,
		}
	}

	result := CheckResult{
		Name:   HealthCheckAdapter,
		Status: CheckStatusOK,
	}
	details := make([]string, 0, len(runtimes))
	nextActions := make([]string, 0, len(runtimes))
	var failures error
	for _, runtime := range runtimes {
		runtimeResult := checker.Adapter(ctx, runtime)
		detail := strings.TrimSpace(runtimeResult.Detail)
		if detail == "" {
			detail = string(runtimeResult.Status)
		}
		if runtimeResult.Status == CheckStatusFailed {
			result.Status = CheckStatusFailed
			if classification := doctorAdapterFailureClassification(runtimeResult.Err); classification != "" {
				detail += "; classification: " + classification
			}
			if action := strings.TrimSpace(runtimeResult.NextAction); action != "" {
				nextActions = append(nextActions, action)
			}
			failures = errors.Join(failures, runtimeResult.Err)
		}
		details = append(details, runtime.ID+": "+detail)
	}
	result.Detail = strings.Join(details, " | ")
	result.NextAction = strings.Join(nextActions, " && ")
	result.Err = failures
	return result
}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/doctor.go` around lines 143 - 187, Update doctorAdapterCheck to
aggregate every failing runtime’s NextAction and Err instead of retaining only
the first failure. Append each non-empty runtimeResult.NextAction to
result.NextAction using the same separator style as Detail, and combine failures
with errors.Join while preserving the existing classification and
successful-runtime behavior.
```

</details>

<!-- fingerprinting:phantom:medusa:beignet -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:82af2233fe479bcf46f3dfd2 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes:
  - A regression test reproduced that two failing runtimes retained only the first remediation and first error.
  - `doctorAdapterCheck` now joins every non-empty remediation with ` && ` and combines independent non-nil failures with `errors.Join`, while preserving the existing single-failure error value.
  - Added CLI-output coverage for two stale adapters and unit coverage proving `errors.Is` reaches both underlying failures.
  - `rtk env GOCACHE=/private/tmp/roundfix-review-001-gocache.QR9F0C go test ./internal/cli -run 'Test(RunDoctorAdapterReadinessReportsRequiredProfileRuntimes|DoctorAdapterCheckAggregatesEveryFailure)$'` passed.
  - A broader `internal/cli` package run reached an unrelated process-identity test and failed because the sandbox denied `/bin/ps`; the Daemon owns authoritative `make verify`.
