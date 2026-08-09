---
source: coderabbit
pr: "143"
round: 1
round_created_at: "2026-08-09T01:14:50Z"
status: failed
head_repository: marcioaltoe/roundfix
head_branch: ma/specs-0082-0083
head_sha: 49ad4407ca050772b13a812b9625cf45d940e19e
file: internal/baseline/preservation.go
line: 721
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XbA2p,comment:PRRC_kwDOS0qyts7e5EBF
review_hash: 6a9e52e7c909589f2ecf3296919de21e012abd3134fc484f7a7d609ac87ef531
duplicate_of: ""
terminal_reason: 'Verification failed: command "make verify" exited with exit status 2; diagnostics: /Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20260809T011807Z_77466ac6469b46dc/verification/batch-001-attempt-2.log'
source_review_id: "4887493512"
source_review_submitted_at: "2026-08-08T00:23:24Z"
---


# Issue 009: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Block managed refresh when the adopted manifest or carrier is absent.**

An absent or invalid Setup Manifest returns no finding. An expected managed carrier with `content == nil` also returns no finding. `BuildPlan` can then accept `managed-refresh`, which bypasses root backups and classification without proving that the repository was adopted.

Return a blocking finding when the manifest is absent or invalid, and when any manifest-managed carrier is missing. The task contract requires missing managed markers to block refresh.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/baseline/preservation.go` around lines 687 - 704, Update the
managed-refresh inspection flow around parseManagedSetupManifest and the content
== nil check to return a blocking Finding when the Setup Manifest is absent or
invalid, and when any manifest-managed carrier is missing. Ensure BuildPlan
cannot accept managed-refresh unless the repository’s adoption markers are
present, while preserving existing findings for valid manifests and available
carriers.
```

</details>

<!-- fingerprinting:phantom:medusa:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:61e534e7a368ca66077e7d3d -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Added blocking Finding when manifest is invalid in classifyManagedRegions (preservation.go:712). Added blocking Finding when manifest-managed carrier is absent (nil content). Updated TestManagedRefreshPlanNeedsNoClassificationInputOrBackup to create a valid manifest since managed-refresh now requires one. `rtk go build ./...` passes.
