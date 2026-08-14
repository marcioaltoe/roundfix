---
source: coderabbit
pr: "157"
round: 2
round_created_at: "2026-08-12T03:19:20Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/what-an-agent-reads-before-it-decides
head_sha: 25099b8b7b5eb496a2871de5912b0515bcc5c27f
file: internal/spec/archive.go
line: 123
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6YcUUP,comment:PRRC_kwDOS0qyts7gUB2P
review_hash: 673c473a5a7568a9509a2da38533525a7f228e2a727e0ed830ed59a2307665df
duplicate_of: ""
source_review_id: "4912735267"
source_review_submitted_at: "2026-08-12T03:17:51Z"
---

# Issue 001: _ Data Integrity & Integration_ _ Major_ _ Heavy lift_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _🏗️ Heavy lift_

**Preserve explicit built-in-root classification through every archive consumer.**

`External` does not identify the built-in Spec Root. A configured non-default root such as `nested/docs/specs` has `External == false`, but `ArchiveSpecRoot` maps it to `nested/_archived/specs`. The documented contract requires `<spec-root>/_archived` for configured non-default roots. The audit path also always passes `false`, so an external root located at `docs/specs` is inspected at the wrong archive location.

- `internal/spec/archive.go#L119-L123`: replace path-suffix inference with an explicit built-in-root classification from Spec Root resolution.
- `internal/spec/archive_test.go#L55-L120`: add regression cases for configured non-default and external roots that both end in `docs/specs`.
- `internal/cli/archive.go#L59-L63`: pass the explicit built-in-root classification into `spec.Archive`.
- `internal/cli/spec_check.go#L369-L375`: use the same classification when checking archived Spec paths.
- `internal/specaudit/audit.go#L335-L339`: retain the configured root classification instead of unconditionally passing `false`.

<details>
<summary>📍 Affects 5 files</summary>

- `internal/spec/archive.go#L119-L123` (this comment)
- `internal/spec/archive_test.go#L55-L120`
- `internal/cli/archive.go#L59-L63`
- `internal/cli/spec_check.go#L369-L375`
- `internal/specaudit/audit.go#L335-L339`

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/spec/archive.go` around lines 119 - 123, Preserve explicit
built-in-root classification through all archive consumers instead of inferring
it from path suffixes. Update ArchiveSpecRoot in internal/spec/archive.go, add
regression coverage in internal/spec/archive_test.go for configured non-default
and external roots ending in docs/specs, pass the resolved classification from
internal/cli/archive.go and internal/cli/spec_check.go, and retain it in
internal/specaudit/audit.go rather than passing false; apply changes at
internal/spec/archive.go:119-123, internal/spec/archive_test.go:55-120,
internal/cli/archive.go:59-63, internal/cli/spec_check.go:369-375, and
internal/specaudit/audit.go:335-339.
```

</details>

<!-- consolidated_sites_start -->
<!--
<consolidated_sites>
<site>
<role>anchor</role>
<file>internal/spec/archive.go</file>
<line_range>119-123</line_range>
</site>
<site>
<role>sibling</role>
<file>internal/spec/archive_test.go</file>
<line_range>55-120</line_range>
</site>
<site>
<role>sibling</role>
<file>internal/cli/archive.go</file>
<line_range>59-63</line_range>
</site>
<site>
<role>sibling</role>
<file>internal/cli/spec_check.go</file>
<line_range>369-375</line_range>
</site>
<site>
<role>sibling</role>
<file>internal/specaudit/audit.go</file>
<line_range>335-339</line_range>
</site>
</consolidated_sites>
-->
<!-- consolidated_sites_end -->

<!-- fingerprinting:phantom:poseidon:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:cea99bf4d80f7640c0ea0cdc -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `RESOLVED`
- Notes: Replaced the path-suffix inference in `spec.ArchiveSpecRoot` with an explicit built-in-root classification computed during Spec Root resolution. `config.SpecsRoot` now carries a `BuiltInRoot` field set by `ResolveSpecsRoot` (true only for the repository's built-in `docs/specs` root). `ArchiveSpecRoot(specsRoot, builtInRoot)` returns the repository default `_archived/specs` only for the built-in root; external and configured non-default roots (including one whose path ends in `docs/specs`, e.g. `nested/docs/specs`) keep their archive beside the active root. Threaded the classification through `spec.Archive` (`ArchiveRequest.BuiltInRoot`), `internal/cli/archive.go`, `internal/cli/spec_check.go`, and `internal/specaudit/audit.go` (which previously hard-coded `false`). Added regression coverage in `internal/spec/archive_test.go` (configured non-default root ending in `docs/specs` plus external roots) and `internal/config/config_test.go` (configured non-default root not classified built-in). Verified: `gofmt`, `go vet`, the full `make verify-incremental` and `make verify-docs` gates pass.
