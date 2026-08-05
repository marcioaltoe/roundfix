---
source: coderabbit
pr: "120"
round: 1
round_created_at: "2026-08-05T14:15:03Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/spec-0067-implementation
head_sha: beca5c076ccfc951eaffc3aeaf7c6a06ed7f6c97
file: internal/baseline/derived_ownership_test.go
line: 1232
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WpdPB,comment:PRRC_kwDOS0qyts7dws4s
review_hash: 120f3aa89265673f5c708f1e458b8ed5ee59b7606ee0aaa544c380c9bc7af11b
duplicate_of: ""
source_review_id: "4864308938"
source_review_submitted_at: "2026-08-05T12:27:49Z"
---

# Issue 002: _ Maintainability & Code Quality_ _ Trivial_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚡ Quick win_

**Use `errors.New` for the two static error messages.**

`fmt.Errorf` is called at Line 1216 and Line 1221 with no format arguments. The configured `perfsprint` linter flags this pattern.






<details>
<summary>♻️ Proposed fix</summary>

```diff
 	start := bytes.Index(content, []byte(prefix))
 	if start == -1 {
-		return nil, fmt.Errorf("goldenDigest field is absent")
+		return nil, errors.New("goldenDigest field is absent")
 	}
 	start += len(prefix)
 	endOffset := bytes.IndexByte(content[start:], '"')
 	if endOffset <= 0 {
-		return nil, fmt.Errorf("goldenDigest value is malformed")
+		return nil, errors.New("goldenDigest value is malformed")
 	}
```

Add `"errors"` to the import block.
</details>

As per coding guidelines: "Use `errors.New` for static error messages" and "Follow the recommended performance linters: `perfsprint` ...".

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
func perturbFormatterGoldenDigest(content []byte) ([]byte, error) {
	const prefix = `"goldenDigest": "`
	start := bytes.Index(content, []byte(prefix))
	if start == -1 {
		return nil, errors.New("goldenDigest field is absent")
	}
	start += len(prefix)
	endOffset := bytes.IndexByte(content[start:], '"')
	if endOffset <= 0 {
		return nil, errors.New("goldenDigest value is malformed")
	}
	result := append([]byte(nil), content...)
	replacement := byte('0')
	if bytes.Count(result[start:start+endOffset], []byte{'0'}) == endOffset {
		replacement = '1'
	}
	for index := start; index < start+endOffset; index++ {
		result[index] = replacement
	}
	return result, nil
}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/baseline/derived_ownership_test.go` around lines 1212 - 1232, Update
perturbFormatterGoldenDigest to use errors.New instead of fmt.Errorf for the
static “goldenDigest field is absent” and “goldenDigest value is malformed”
errors, and add the errors import while preserving all existing behavior.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:60bd06fbc770c1fb7f8cdf0c -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Both static `fmt.Errorf` calls were present and did not format or wrap
  values.

## Resolution

- Replaced both static golden-digest perturbation errors with `errors.New` and
  added the standard-library import.

## Focused evidence

- `rtk env GOCACHE=/private/tmp/roundfix-review-c8087f92-gocache GOFLAGS=-buildvcs=false go test ./internal/baseline -count=1`
  — exit 0 in 109.338s.
- `rtk git diff --check` — exit 0.
- `make verify` was not run; authoritative Verification is Daemon-owned.
