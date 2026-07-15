---
source: coderabbit
pr: "27"
round: 2
round_created_at: "2026-07-15T16:07:28Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/dogfood-findings-remediation
head_sha: 81ff89aabce7ee1f748504ad6e4bcf1ac52ea200
file: internal/cli/cli.go
line: 2469
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6RJ95K,comment:PRRC_kwDOS0qyts7V5tZ3
review_hash: 30ca3b07b68feef965a1cf3d0fcf0537a1a34d28db041e4dceacdd9a2ad20ad4
duplicate_of: ""
source_review_id: "4705718496"
source_review_submitted_at: "2026-07-15T15:38:26Z"
---

# Issue 004: _ Stability & Availability_ _ Major_ _ Quick win_

## Review Comment

_🩺 Stability & Availability_ | _🟠 Major_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
set -euo pipefail

# Locate the relevant symbols and nearby code.
rg -n --context 3 --hidden --glob '!**/.git/**' \
  'func reclaimOrphanedActiveRun|func forceStopRun|bestEffortForceStopAgentSessions|activeOwnerPID|orphanedActiveRunReason|ReclaimOrphanedRun|ProcessAlive|Run\(ctx, active.ID\)' \
  internal/cli/cli.go internal -g '*.go'

# Show the target function with line numbers.
sed -n '2430,2495p' internal/cli/cli.go

# Find all callers of the session-stop helper and the reclaim helper.
rg -n --hidden --glob '!**/.git/**' \
  'bestEffortForceStopAgentSessions\(|reclaimOrphanedActiveRun\(' \
  internal -g '*.go'
```

Repository: marcioaltoe/roundfix

Length of output: 15887

---

🏁 Script executed:

```shell
set -euo pipefail

# Inspect reclaim semantics in the store layer.
sed -n '375,470p' internal/store/store.go

# Inspect force-stop behavior for contrast.
sed -n '533,565p' internal/cli/cli.go

# Search for any session cleanup triggered by orphan reclaim.
rg -n --hidden --glob '!**/.git/**' \
  'Agent Session|cancel.*Agent|ForceStopAgent|bestEffortForceStopAgentSessions|ReclaimOrphanedRun|reclaim orphaned Active Run|reclaimed orphaned Active Run' \
  internal -g '*.go'
```

Repository: marcioaltoe/roundfix

Length of output: 13332

---



</details>

**Stop Agent Sessions before reclaiming the orphaned run lock.** `ReclaimOrphanedRun` only marks the run failed and deletes the lock; it doesn’t cancel the lingering Agent Session like `forceStopRun` does, so a new run can start while the old session is still active. `internal/cli/cli.go:2449`

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/cli.go` around lines 2449 - 2469, Update
reclaimOrphanedActiveRun to stop the active run’s lingering Agent Session before
calling runStore.ReclaimOrphanedRun, reusing the same session-stop behavior and
symbols as forceStopRun. Preserve the existing reclaim, verification, error
handling, and reporting flow after the session is stopped.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:89e5a6876902bad02b433ee1 -->

<!-- This is an auto-generated reply by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Orphaned Active Run reclaim now best-effort cancels and closes lingering Agent Sessions before releasing the lock; `make verify` passed.
