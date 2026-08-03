---
source: coderabbit
pr: "87"
round: 1
round_created_at: "2026-08-03T16:19:44Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/qa-task-not-flag
head_sha: d4011aade56592440d7a542682ebee7dced31f79
file: internal/spec/spec.go
line: 145
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WCzZR,comment:PRRC_kwDOS0qyts7c4SN1
review_hash: 56e681d443c0e722e1c357c7b12ab762b922dc72d62018eab5d519c8ac215d1d
duplicate_of: ""
source_review_id: "4846253969"
source_review_submitted_at: "2026-08-03T16:15:59Z"
---

# Issue 004: _ Maintainability & Code Quality_ _ Minor_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟡 Minor_ | _⚡ Quick win_

**Wrap errors at the new YAML and QA-validation boundaries.**

At Line 144, `UnmarshalYAML` returns the `node.Decode` error without context. At Line 290, `Load` returns the `validateQAGate` error without context. Add lowercase context with `%w` so callers can identify the failing phase while preserving `errors.Is` and `errors.As`.

<details>
<summary>Proposed fix</summary>

```diff
 	if err := node.Decode(&decoded); err != nil {
-		return err
+		return fmt.Errorf("decode manifest frontmatter: %w", err)
 	}
...
 	if err := validateQAGate(manifestPath, nodes, tasks, qa); err != nil {
-		return nil, err
+		return nil, fmt.Errorf("validate qa gate: %w", err)
 	}
```
</details>

As per coding guidelines, Go code must wrap propagated errors with `%w` and use lowercase error strings.






Also applies to: 289-291

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/spec/spec.go` around lines 143 - 145, Wrap the error returned by
node.Decode in UnmarshalYAML with a lowercase phase-specific message using %w,
and wrap the error returned by validateQAGate in Load the same way. Preserve
error unwrapping so errors.Is and errors.As continue to work.
```

</details>

<!-- fingerprinting:phantom:poseidon:luna -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:8181c88a8a5d043fe092465b -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes:
  - The new manifest decode and QA gate validation boundaries propagated errors without identifying the failed phase, contrary to the repository's `%w` wrapping rule.
  - Added regression coverage for manifest decode context and extended the QA gate table to require validation context while preserving `errors.As` access to `ManifestError` and `QAGateError`.
  - Before the production fix, the focused tests failed because neither `decode manifest frontmatter` nor `validate qa gate` appeared in the error chain.
  - Both boundaries now wrap with lowercase operation context and `%w`.
  - `rtk env GOCACHE=/private/tmp/roundfix-batch001-spec-gocache go test ./internal/spec -run '^(TestLoadWrapsManifestDecodeErrors|TestLoadRejectsInvalidQAGateShape|TestLoadInvalidatesSettledQAGateAfterTaskAppend)$' -count=1`: passed.
  - `rtk env GOCACHE=/private/tmp/roundfix-batch001-packages-gocache go test ./internal/daemon ./internal/spec -count=1`: passed.
  - Daemon Verification `make verify` was not run by this Agent; the Daemon owns authoritative Verification after this turn.
