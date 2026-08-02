---
source: coderabbit
pr: "72"
round: 1
round_created_at: "2026-08-02T21:04:16Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/baseline-capability-evidence-and-retention
head_sha: cb7c719e649ac8558f3b4d10a993516b29ececb5
file: internal/baseline/plan_test.go
line: 2511
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6V0Ymf,comment:PRRC_kwDOS0qyts7cjgE7
review_hash: 6c94ec00d9878fd530d12bd174de4aa61abf08553b43377132ab36709930f342
duplicate_of: ""
source_review_id: "4839703297"
source_review_submitted_at: "2026-08-02T21:03:30Z"
---

# Issue 010: _ Maintainability & Code Quality_ _ Trivial_ _ Low value_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _💤 Low value_

**Use a comma-ok type assertion in `cloneCatalogForRetentionDrift`.**

Line 2504 uses a bare type assertion on the result of `cloneJSONValue`. If a module clone is not a `map[string]any`, the test panics without a diagnostic. The configured linter set includes `forcetypeassert`, which flags this form. Convert the helper to take `*testing.T`, or assert with the comma-ok form and fail the test.

As per coding guidelines: "Check all error returns, use comma-ok type assertions".






<details>
<summary>♻️ Proposed comma-ok assertion</summary>

```diff
-func cloneCatalogForRetentionDrift(source *Catalog) *Catalog {
+func cloneCatalogForRetentionDrift(t *testing.T, source *Catalog) *Catalog {
+	t.Helper()
 	target := *source
 	target.modules = make(map[string]document, len(source.modules))
 	for id, module := range source.modules {
-		target.modules[id] = document(cloneJSONValue(module).(map[string]any))
+		cloned, ok := cloneJSONValue(module).(map[string]any)
+		if !ok {
+			t.Fatalf("cloned catalog module %q is not an object", id)
+		}
+		target.modules[id] = document(cloned)
 	}
```
</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/baseline/plan_test.go` around lines 2500 - 2511, Update
cloneCatalogForRetentionDrift to replace the bare assertion on
cloneJSONValue(module) with a comma-ok assertion, and fail the test with a
useful diagnostic when the cloned value is not a map[string]any; thread
*testing.T into the helper if needed and update its callers.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:87c279cc0639ad9118f9d843 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `cloneCatalogForRetentionDrift` now receives `*testing.T`, marks itself as a helper, and reports a typed clone failure instead of panicking. The full Baseline package test passed.
