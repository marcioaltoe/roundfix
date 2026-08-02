---
source: coderabbit
pr: "67"
round: 1
round_created_at: "2026-08-02T11:30:39Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/profiles-configure-merge-semantics
head_sha: ffcc15ebed0a055d329cb3215ae0878b90931948
file: internal/config/config_test.go
line: 1315
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Vwh2M,comment:PRRC_kwDOS0qyts7ceBLv
review_hash: 8c9b2471dfe074abc0479c07bb73631cd99ddef8891c49f69ed3fa25a6b3fcb6
duplicate_of: ""
source_review_id: "4838273774"
source_review_submitted_at: "2026-08-02T11:29:42Z"
---

# Issue 006: _ Functional Correctness_ _ Trivial_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🔵 Trivial_ | _⚡ Quick win_

**Add a case for an unsupported removal category.**

`DeriveEffectiveChangeSet` rejects a removal that `ParseWorkCategory` cannot resolve and returns a message listing the supported values. No subtest exercises that branch. The CLI surfaces this message to users as a validation failure with exit code `2`, so the wording and the rejection both matter.

Add a subtest that passes `[]WorkCategory{"not-a-category"}` and asserts the error names the unsupported value.




<details>
<summary>💚 Proposed additional subtest</summary>

```go
t.Run("rejects an unsupported removal category", func(t *testing.T) {
	_, err := DeriveEffectiveChangeSet(nil, nil, []WorkCategory{"not-a-category"})
	if err == nil {
		t.Fatal("expected unsupported removal category to fail validation")
	}
	if !strings.Contains(err.Error(), "not-a-category") {
		t.Fatalf("expected validation error to name the unsupported value, got %q", err)
	}
})
```
</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/config/config_test.go` around lines 1290 - 1315, Add a subtest
alongside the existing DeriveEffectiveChangeSet cases that passes
[]WorkCategory{"not-a-category"} with nil profiles and asserts an error is
returned. Verify the error text includes "not-a-category", covering rejection of
unsupported removal categories and preserving the user-facing validation
message.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:60088bd1e2835d406af643fd -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `DeriveEffectiveChangeSet` rejects unsupported removal categories, but its negative branch had no direct test. `TestEffectiveChangeSet/rejects_an_unsupported_removal_category` now asserts both rejection and that the diagnostic names the unsupported value.
- Focused evidence: `rtk proxy env GOCACHE=/Users/marcio/dev/roundfix/.gocache rtk go test ./internal/config -run '^(TestProfilesConfigureRemovalPreservesSpacing|TestProfilesConfigureMergePreservesOtherCategories|TestEffectiveChangeSet)$' -count=1` passed 21 tests.
