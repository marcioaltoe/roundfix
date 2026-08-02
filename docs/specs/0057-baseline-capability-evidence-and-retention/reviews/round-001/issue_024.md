---
source: coderabbit
pr: "72"
round: 1
round_created_at: "2026-08-02T21:04:16Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/baseline-capability-evidence-and-retention
head_sha: cb7c719e649ac8558f3b4d10a993516b29ececb5
file: internal/cli/baseline_human_test.go
line: 1416
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6V0Ym4,comment:PRRC_kwDOS0qyts7cjgFY
review_hash: 0644ef00315c1462eda86c29eb2b2e78b9ceb500b77b7d53266437350ba6877b
duplicate_of: ""
source_review_id: "4839703297"
source_review_submitted_at: "2026-08-02T21:03:30Z"
---

# Issue 024: _ Maintainability & Code Quality_ _ Major_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟠 Major_ | _⚡ Quick win_

**Reuse `standardTypeScriptDivergenceDecisions` in the existing decline subtest.**

The new helper at lines 1391-1416 duplicates the inline decision slice at lines 1122-1145 in the "decline writes nothing" subtest. Both lists carry the same twelve decisions with the same values. Two copies will drift when the catalog adds or renames a decision, and only one of the two call sites would be updated.

Replace the inline slice with a call to the helper.





<details>
<summary>♻️ Proposed change at lines 1122-1145</summary>

```diff
 			baselineHumanState{},
 			source,
-			[]baseline.DecisionValue{
-				{ID: "language.generated", Value: "English"},
-				{ID: "verification.gate", Value: "make verify"},
-				{ID: "identifier.strategy", Value: map[string]any{"kind": "uuid-v7"}},
-				{ID: "http.contract", Value: map[string]any{"mode": "Post-only"}},
-				{
-					ID: "auth.provider",
-					Value: map[string]any{
-						"kind": "better-auth",
-						"routeException": map[string]any{
-							"scope":   "/api/auth/*",
-							"methods": []any{"GET", "POST"},
-							"owner":   "Better Auth",
-							"reason":  "Provider protocol routes require GET and POST semantics.",
-						},
-					},
-				},
-				{ID: "spec.scaffold", Value: true},
-				{ID: "domain.layout", Value: "single-context"},
-				{ID: "triage.external", Value: false},
-				{ID: "autonomous.enabled", Value: false},
-				{ID: "secondbrain.enabled", Value: false},
-				{ID: "repository.extension.enabled", Value: false},
-			},
+			standardTypeScriptDivergenceDecisions(),
 		)
```
</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/baseline_human_test.go` around lines 1391 - 1416, In the
existing “decline writes nothing” subtest, replace the duplicated inline
decision slice with a call to standardTypeScriptDivergenceDecisions(). Leave the
helper’s decision values and the surrounding test behavior unchanged.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:refactor_suggestion -->

<!-- cr-comment:v1:f66ed56398ef8ecf5e81654a -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The decline subtest now reuses standardTypeScriptDivergenceDecisions, removing the duplicate decision fixture without changing behavior. `go test ./internal/baseline ./internal/cli` passed with a writable GOCACHE.
