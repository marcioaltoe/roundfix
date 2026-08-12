---
source: coderabbit
pr: "72"
round: 1
round_created_at: "2026-08-02T21:04:16Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/baseline-capability-evidence-and-retention
head_sha: cb7c719e649ac8558f3b4d10a993516b29ececb5
file: internal/baseline/profile_alignment_test.go
line: 167
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6V0Ymn,comment:PRRC_kwDOS0qyts7cjgFD
review_hash: af68cd91926483b9da7b2e911b3915a57186e4ede7cd806694ab50e29e05d406
duplicate_of: ""
source_review_id: "4839703297"
source_review_submitted_at: "2026-08-02T21:03:30Z"
---

# Issue 014: _ Maintainability & Code Quality_ _ Trivial_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚡ Quick win_

**Cover the remaining `ErrNoResolvableProfile` branches.**

`TestCapabilityRecheck` covers the missing-manifest branch and the successful manifest branch. Two branches of `resolveCapabilityRecheckProfile` stay untested:

- A present but invalid Setup Manifest, which returns `"%w: current Setup Manifest is invalid"`.
- An explicit `ProfileID` that `ResolveProfile` rejects, which returns `fmt.Errorf("%w %q: %w", ...)`.

The second branch uses two `%w` verbs in one `fmt.Errorf`. A test that asserts `errors.Is(err, ErrNoResolvableProfile)` and `errors.Is(err, <resolve error>)` would pin that multi-error wrapping.

Add two subtests.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/baseline/profile_alignment_test.go` around lines 125 - 167, Extend
TestCapabilityRecheck with subtests for a present but invalid Setup Manifest and
for an explicit ProfileID rejected by ResolveProfile. Assert both errors wrap
ErrNoResolvableProfile; in the explicit-profile case, also assert errors.Is
matches the underlying ResolveProfile error, preserving the multi-error wrapping
contract.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:4eff26db8645fc78c6e3413d -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Added re-check coverage for an invalid current Setup Manifest and an explicit missing Profile, including `errors.Is` assertions for `ErrNoResolvableProfile` and `fs.ErrNotExist`. The full Baseline package test passed.
