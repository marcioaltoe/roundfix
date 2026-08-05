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
line: 58
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WpdPF,comment:PRRC_kwDOS0qyts7dws4x
review_hash: 241189fe431ad6bdeec79c4e082860e352faf2c7125b046943d983e0166ddb7d
duplicate_of: ""
source_review_id: "4864308938"
source_review_submitted_at: "2026-08-05T12:27:49Z"
---

# Issue 003: _ Maintainability & Code Quality_ _ Trivial_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚡ Quick win_

**Remove the duplicate exception-claim computation.**

`readDerivedOwnershipRecords` already calls `derivedOwnershipExceptionClaims` at Line 173 and discards the result. Line 55 computes the same claims again. Keep the validation in one place: drop the call at Line 173 and let callers build the claim map, or return the claims from `readDerivedOwnershipRecords`. The same duplication also runs in `derivedArtifactRemediation` at Lines 437-444.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/baseline/derived_ownership.go` around lines 51 - 58, Remove the
duplicate derivedOwnershipExceptionClaims computation by keeping claim
validation in one place only. Update readDerivedOwnershipRecords so it no longer
calls derivedOwnershipExceptionClaims and discards the result, then let the
caller path in derived_ownership.go build the claim map once after
readDerivedOwnershipRecords returns; apply the same cleanup to the
derivedArtifactRemediation flow so it does not recompute the same claims.
Preserve the existing error handling around records loading and claim
validation, but avoid calling derivedOwnershipExceptionClaims twice for the same
records.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:93fb480709cce34714c7ef97 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `readDerivedOwnershipRecords` validated exception claims and discarded
  them before both production callers immediately rebuilt the same claim map.

## Resolution

- Removed the discarded claim-map construction from
  `readDerivedOwnershipRecords`; validation and remediation now each build and
  validate the map once at the point where they consume it.

## Focused evidence

- `rtk env GOCACHE=/private/tmp/roundfix-review-c8087f92-gocache GOFLAGS=-buildvcs=false go test ./internal/baseline -count=1 -run '^TestDerivedOwnership'`
  — exit 0.
- The complete `internal/baseline` package check exited 0 in 109.338s.
- `make verify` was not run; authoritative Verification is Daemon-owned.
