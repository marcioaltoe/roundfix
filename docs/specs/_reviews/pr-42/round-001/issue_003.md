---
source: coderabbit
pr: "42"
round: 1
round_created_at: "2026-07-29T02:33:36Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/claude-adapter-standardization
head_sha: 7155ba4d2ef353257a1bacf697027202d4750492
file: internal/cli/cli_test.go
line: 9989
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6UmtXo,comment:PRRC_kwDOS0qyts7ayC0B
review_hash: 49ec7914f72ece9332d17af9281a568576314201cb80c2a0d5d42ff53dbf9c9a
duplicate_of: ""
source_review_id: "4803488138"
source_review_submitted_at: "2026-07-29T02:32:46Z"
---

# Issue 003: _ Maintainability & Code Quality_ _ Trivial_ _ Low value_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _💤 Low value_

**Two overlapping adapter-error injection paths in the fake.**

`fake.adapterErr` (Codex-only, line 10070) and `fake.adapterErrors[runtimeID]` (line 10067) now express the same thing with different reach, so new tests must guess which one applies. Consider folding the legacy `adapterErr` into `adapterErrors["codex"]` in `newSetupFakeDeps` and dropping the special case.






Also applies to: 10063-10070

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/cli_test.go` around lines 9963 - 9989, Unify adapter error
injection in setupFakeDeps by removing the separate adapterErr field and routing
Codex failures through adapterErrors["codex"]. Update newSetupFakeDeps and the
fake adapter handling to use adapterErrors consistently, preserving existing
behavior for Codex and other runtimes.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:2ab127477d820badc95aafd8 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes:
  - Removed the Codex-only `adapterErr` path from `setupFakeDeps`.
  - Initialized `adapterErrors` in `newSetupFakeDeps` and routed Codex setup failures through `adapterErrors["codex"]`, matching the existing multi-runtime injection path.
  - `rtk env GOCACHE=/private/tmp/roundfix-review-001-gocache.QR9F0C go test ./internal/cli -run 'TestRunSetup(ReportsAdapterFailuresWithoutWrites|AdapterMigrationDeclinePreservesAllTargets|AdapterMigrationPersistsSupportedCommand|ClaudeAdapterMigrationAcceptAndDecline|MigratesBothStaleAdapterOverrides)$'` passed.
