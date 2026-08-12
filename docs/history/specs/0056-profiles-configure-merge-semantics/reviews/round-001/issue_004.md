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
line: 836
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Vwh2I,comment:PRRC_kwDOS0qyts7ceBLp
review_hash: c6f50790ea5603a6a9b7e7ffea253e5599e69046eb59c8c089349de4627cd98c
duplicate_of: ""
source_review_id: "4838273774"
source_review_submitted_at: "2026-08-02T11:29:42Z"
---

# Issue 004: _ Maintainability & Code Quality_ _ Trivial_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚡ Quick win_

**Do not assert on yaml.v3's internal error text.**

`"unknown anchor"` is produced by `gopkg.in/yaml.v3`, not by this package. A dependency upgrade that rewords the message breaks this test without any behavior change. The same assertion appears at line 1214.

Assert the observable contract instead: the write returns a non-nil error and the file bytes are unchanged. If the specific cause must be checked, wrap the validation failure in a package sentinel error in `mergeProfilesConfigContent` and match it with `errors.Is`.





As per coding guidelines: "NEVER test implementation details; test observable behavior and public API contracts."

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/config/config_test.go` around lines 834 - 836, Replace the
yaml.v3-specific “unknown anchor” string assertions in the
remove-anchored-backend tests, including the matching assertion near line 1214,
with checks that the write returns a non-nil error and the file bytes remain
unchanged. Update the relevant test setup and comparisons around
mergeProfilesConfigContent without asserting dependency-generated error text;
only introduce a package sentinel and errors.Is matching if validating the
specific cause is required.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:607f0b6d44084b9bb8c1efcb -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Both alias-stranding tests asserted yaml.v3's current `unknown anchor` wording even though the owned contract is rejection before persistence. They now require a non-nil error and byte-identical config content without coupling to dependency text.
- Focused evidence: `rtk proxy env GOCACHE=/Users/marcio/dev/roundfix/.gocache rtk go test ./internal/config -run '^(TestProfilesConfigureRemovalPreservesSpacing|TestProfilesConfigureMergePreservesOtherCategories|TestEffectiveChangeSet)$' -count=1` passed 21 tests.
