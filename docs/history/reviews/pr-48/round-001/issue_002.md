---
source: coderabbit
pr: "48"
round: 1
round_created_at: "2026-07-29T21:58:48Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/repository-derived-skill-requirements
head_sha: 3ef6a563f8be4a4e72a2a063463d904fd0e0a9a1
file: internal/cli/doctor.go
line: 360
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6U6ur2,comment:PRRC_kwDOS0qyts7bPmQG
review_hash: 59ab30958ceb76a2d5608e76e8017651b72e0c09533cdfbc2062584760ae5948
duplicate_of: ""
source_review_id: "4813239038"
source_review_submitted_at: "2026-07-29T21:57:47Z"
---

# Issue 002: _ Maintainability & Code Quality_ _ Trivial_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚡ Quick win_

**Fragile substring check to detect "0 external" already rendered.**

`strings.Contains(result.Detail, "0 external")` (Line 349) infers, from rendered text, whether the base `doctorSkillReadinessResult` call already reported zero external skills. This couples correctness to the exact wording of an unrelated format string — if that format string ever changes (e.g. wording, i18n, or a skill legitimately named to contain the substring "0 external"), the check silently breaks. Since the condition it emulates is exactly `checkErr == nil && readiness.Ready()` (the same predicate `doctorSkillReadinessResult` uses internally), check that directly instead of parsing the rendered string.

<details>
<summary>♻️ Proposed fix</summary>

```diff
-	if !strings.Contains(result.Detail, "0 external") {
+	if !(checkErr == nil && readiness.Ready()) {
 		result.Detail += "; 0 external required"
 	}
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
func doctorMissingSetupManifestResult(readiness skills.RepositoryReadiness, checkErr error) CheckResult {
	result := doctorSkillReadinessResult(readiness, checkErr)
	result.Status = CheckStatusFailed
	if result.Detail == "" {
		result.Detail = missingSetupManifestDetail
	} else {
		result.Detail += "; " + missingSetupManifestDetail
	}
	if !(checkErr == nil && readiness.Ready()) {
		result.Detail += "; 0 external required"
	}

	next := make([]string, 0, 2)
	if doctorOwnedSkillRequirementFailed(readiness, checkErr) {
		next = append(next, ownedSkillsNextAction)
	}
	next = append(next, baselineAdoptionNextAction)
	result.NextAction = strings.Join(next, " && ")
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

In `@internal/cli/doctor.go` around lines 341 - 360, In
doctorMissingSetupManifestResult, replace the fragile
strings.Contains(result.Detail, "0 external") check with the underlying
predicate checkErr == nil && readiness.Ready(). Append "0 external required"
only when that predicate is false, preserving the existing detail construction
and next-action behavior.
```

</details>

<!-- fingerprinting:phantom:medusa:beignet -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:267735f5689771501753d30a -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes:
  - `doctorMissingSetupManifestResult` inferred readiness from rendered detail text even though `doctorSkillReadinessResult` uses the exact `checkErr == nil && readiness.Ready()` predicate.
  - Replaced the substring check with that predicate and added a public Doctor regression case proving checker error text containing `0 external` cannot suppress the explicit `0 external required` detail.
  - `rtk proxy env GOCACHE=/Users/marcio/dev/roundfix/.gocache go test ./internal/cli -run 'TestRunDoctor(DerivesExternalSkillRequirementFromSetupManifest|RepositorySkillReadiness)$' -count=1`: passed.
  - The initial equivalent focused test attempt without an explicit writable `GOCACHE` was blocked before compilation by the sandbox-denied user Go cache; the unchanged check passed with the repository-local cache above.
  - Daemon Verification `make verify` was not run by this Agent; the Daemon owns authoritative Verification after this turn.
