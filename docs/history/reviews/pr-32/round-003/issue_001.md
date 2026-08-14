---
source: coderabbit
pr: "32"
round: 3
round_created_at: "2026-07-17T14:20:26Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/release-plan-and-agent-selection-profiles
head_sha: 204bbd00fbc648be0df0b8bf2f883b9e2dc490c8
file: internal/cli/implement.go
line: 379
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Ry2aG,comment:PRRC_kwDOS0qyts7Wzc-U
review_hash: 1afdca951226c1c660c2051d827228ed1a8e52b0c1f1e9077868580393352089
duplicate_of: ""
source_review_id: "4723452116"
source_review_submitted_at: "2026-07-17T14:16:02Z"
---

# Issue 001: _ Data Integrity & Integration_ _ Major_ _ Heavy lift_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _🏗️ Heavy lift_

**Do not silently route external Specs from the dirty checkout.**

When `resolvedSpecsRoot` is outside `gitRoot`, this fallback reloads the live path, so profile selection is no longer based on the committed graph. Either archive from the external Specs repository/revision or fail preflight with an actionable error.

As per coding guidelines, “Do not use production or test workarounds; fix the root cause.”

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/implement.go` around lines 375 - 379, Update the fallback in
implement.go around repositoryRelativePath so an external checkoutSpecsRoot is
never reloaded directly through spec.Load. Preserve committed-graph/profile
selection by archiving from the external Specs repository and revision, or fail
preflight with an actionable error; do not silently route to the dirty checkout.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:f52e467391cfc788e7864cdb -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `defaultLoadCommittedSpecGraph` fell back to `spec.Load` for external Spec Roots, which let dirty external files drive committed-graph/profile selection. Fixed by loading external Specs from their own Git `HEAD` via `git archive`; external roots outside any Git repository now fail preflight with an actionable error.

## Resolution

- Changed `internal/cli/implement.go` to archive external Spec Roots from the external Git repository instead of reading live files.
- Added regression coverage for dirty external Spec metadata and for non-Git external Spec Roots.
- Evidence:
  - `rtk go test ./internal/cli ./internal/store -run 'TestLoadCommitted|TestRunImplementUsesConfiguredExternalSpecRootEndToEnd|TestRunImplementInteractiveInputListsConfiguredExternalSpecRoot|TestAgentSelectionAttemptLifecycleUpdatesSameAttempt'` — passed.
  - `rtk go test ./internal/cli ./internal/store` — passed.
  - `rtk make verify` — passed.
