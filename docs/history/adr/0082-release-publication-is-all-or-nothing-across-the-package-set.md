---
status: superseded # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-07-28T22:02:41Z
updated_at: 2026-08-26T00:00:00Z
deprecated_at: null
superseded_by: ADR-0146
---

# Release publication is all-or-nothing across the package set

The `v0.0.1` reset showed the release workflow can begin publishing platform
packages while the launcher's registry coordinate is ineligible (a used
version, a post-unpublish cooldown, or missing ownership), which would leave
installable platform packages with no launcher — a partial release npm policy
makes irreversible. A read-only publication preflight therefore evaluates the
launcher and every platform package as one release set after Verification and
before cross-compilation, and no package publishes unless every coordinate is
eligible for the exact target version. Publishing the eligible subset was
rejected: a stopped release is recoverable, a partial one is not.
