---
source: coderabbit
pr: "120"
round: 1
round_created_at: "2026-08-05T14:15:03Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/spec-0067-implementation
head_sha: beca5c076ccfc951eaffc3aeaf7c6a06ed7f6c97
file: internal/baseline/derived_ownership.go
line: 219
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WpdPO,comment:PRRC_kwDOS0qyts7dws5B
review_hash: 3fe71af0edda6b92b1df6644e23e36530810a821159ded5cfc88d5edb89aa09f
duplicate_of: ""
source_review_id: "4864308938"
source_review_submitted_at: "2026-08-05T12:27:49Z"
---

# Issue 005: _ Maintainability & Code Quality_ _ Major_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟠 Major_ | _⚡ Quick win_

**Use `errors.New` for the static error messages.**

Lines 215 and 218 call `fmt.Errorf` with no format verbs. The coding guidelines require `errors.New` for static error messages, and `perfsprint` flags constant `fmt.Errorf` calls. Line 285 has the same pattern.

As per coding guidelines: "Use `errors.New` for static error messages and package-level sentinel errors when callers need to match a specific condition."





<details>
<summary>♻️ Proposed fix</summary>

```diff
 	if record.Reason == "" {
-		return fmt.Errorf("reason is required")
+		return errors.New("reason is required")
 	}
 	if record.Owner == derivedOwnerDedicated && record.Command == "" {
-		return fmt.Errorf("command is required for dedicated ownership")
+		return errors.New("command is required for dedicated ownership")
 	}
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
	if record.Reason == "" {
		return errors.New("reason is required")
	}
	if record.Owner == derivedOwnerDedicated && record.Command == "" {
		return errors.New("command is required for dedicated ownership")
	}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/baseline/derived_ownership.go` around lines 214 - 219, Replace the
constant fmt.Errorf calls in the validation logic for record.Reason and
dedicated record.Command with errors.New, and apply the same change to the
static error at the corresponding line-285 validation path. Add or reuse the
errors import as needed without changing error messages or behavior.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:refactor_suggestion -->

<!-- cr-comment:v1:244b069b80a985f616a4dc10 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: All three reported static `fmt.Errorf` calls were present and carried
  neither formatted values nor wrapped errors.

## Resolution

- Replaced the record reason, dedicated command, and exception-path static
  errors with `errors.New` without changing their messages.

## Focused evidence

- `rtk env GOCACHE=/private/tmp/roundfix-review-c8087f92-gocache GOFLAGS=-buildvcs=false go test ./internal/baseline -count=1`
  — exit 0 in 109.338s.
- `rtk git diff --check` — exit 0.
- `make verify` was not run; authoritative Verification is Daemon-owned.
