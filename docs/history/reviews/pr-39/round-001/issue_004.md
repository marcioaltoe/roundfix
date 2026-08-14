---
source: coderabbit
pr: "39"
round: 1
round_created_at: "2026-07-27T21:22:30Z"
status: invalid
terminal_reason: "The repository requires stdlib-only tests, and this test already joins all eight bounded goroutines through an exact-size result channel."
head_repository: marcioaltoe/roundfix
head_branch: ma/terminal-run-worktree-reconciliation
head_sha: 44fa422cea404a2d8c47e4b8011f065c4c0481ba
file: internal/worktree/worktree_test.go
line: 895
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6UNfMk,comment:PRRC_kwDOS0qyts7aNkLL
review_hash: 41b8a5ddf5baca38e1872ee841e71b9dc8aab340c0a378d3d5655774257839ce
duplicate_of: ""
source_review_id: "4791610618"
source_review_submitted_at: "2026-07-27T21:21:25Z"
---

# Issue 004: _ Maintainability & Code Quality_ _ Minor_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟡 Minor_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
# Check for TestMain/goleak usage in the worktree package and repo-wide
rg -nP --type=go 'goleak|func TestMain\(' -g '!vendor/**' -C2
```

Repository: marcioaltoe/roundfix

Length of output: 158

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf 'Files in internal/worktree:\n'
git ls-files 'internal/worktree/*' | sed 's#^`#-` #'

printf '\nSearch for TestMain and goleak in internal/worktree:\n'
rg -n --hidden --glob '!vendor/**' 'func TestMain\(|goleak' internal/worktree || true

printf '\nSearch for goroutine-spawning patterns in internal/worktree tests:\n'
rg -n --hidden --glob 'internal/worktree/*_test.go' '\bgo func\(|\bgo\s+[A-Za-z_]' internal/worktree || true

printf '\nRelevant slice of internal/worktree/worktree_test.go:\n'
sed -n '830,930p' internal/worktree/worktree_test.go
```

Repository: marcioaltoe/roundfix

Length of output: 4244

---



</details>

**Add goroutine-leak checks for `internal/worktree`**  
`internal/worktree/worktree_test.go` now spawns goroutines, and the package has no `TestMain`/`goleak` coverage. Add `goleak.VerifyTestMain(m)` for the package or `goleak.VerifyNone(t)` around the concurrent tests so leaks fail fast.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/worktree/worktree_test.go` around lines 867 - 895, Add goleak
coverage for the goroutines introduced by TestInspectTerminalRunConcurrentSafe,
using goleak.VerifyNone around the concurrent test or goleak.VerifyTestMain via
a package TestMain. Ensure leaked goroutines cause the internal/worktree tests
to fail while preserving the existing inspection assertions.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:5dd8b8672efece083086a47c -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: `docs/agents/go.md` prohibits test dependencies and requires stdlib
  `testing` only, so adding `go.uber.org/goleak` would violate a repository
  Normative Clause. The test starts exactly eight goroutines, gives each one
  result slot in an eight-element buffered channel, and receives exactly eight
  results before returning. There is no fire-and-forget package goroutine or
  lifecycle gap that warrants a new dependency.
