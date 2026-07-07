---
source: coderabbit
pr: "18"
round: 2
round_created_at: "2026-07-07T14:05:56Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/docs-and-skill-bundle
head_sha: 4237143afdd7097e755e14b962156aaf6c6e6654
file: internal/cli/settle.go
line: 307
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6O6PJ_,comment:PRRC_kwDOS0qyts7SyYRy
review_hash: 0578884f13bcd13a9ad2e8f55c9674b5ac92ea392a464bfdeb56db0f9b1e8343
duplicate_of: ""
source_review_id: "4645087962"
source_review_submitted_at: "2026-07-07T12:31:07Z"
---

# Issue 004: _ Stability & Availability_ _ Major_ _ Heavy lift_

## Review Comment

_🩺 Stability & Availability_ | _🟠 Major_ | _🏗️ Heavy lift_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
# Inspect how the shared committer handles paths outside WorkDir, to confirm
# it safely drops (rather than fails on) an external Spec Root artifact path.
rg -n "kept outside the repository" internal
rg -n -A20 "func.*Commit\(" internal/daemon/task_engine.go
```

Repository: marcioaltoe/roundfix

Length of output: 828

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '\n== settle.go (around referenced lines) ==\n'
sed -n '280,330p' internal/cli/settle.go

printf '\n== settleArtifactCommitPath and nearby ==\n'
sed -n '380,430p' internal/cli/settle.go

printf '\n== committer implementation references ==\n'
rg -n -A40 -B10 "kept outside the repository|type .*Committer|func .*Commit\(" internal
```

Repository: marcioaltoe/roundfix

Length of output: 50376

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '\n== internal/cli/settle.go relevant context ==\n'
nl -ba internal/cli/settle.go | sed -n '290,320p'

printf '\n== settleArtifactCommitPath ==\n'
nl -ba internal/cli/settle.go | sed -n '395,415p'

printf '\n== task_engine commit path handling ==\n'
nl -ba internal/daemon/task_engine.go | sed -n '760,810p'
```

Repository: marcioaltoe/roundfix

Length of output: 245

---

🏁 Script executed:

```shell
#!/bin/bash
set -euo pipefail

printf '\n== settleArtifactCommitPath and ensureSettleCommitPath ==\n'
sed -n '395,415p' internal/cli/settle.go

printf '\n== GitCommitter Commit path staging ==\n'
sed -n '60,82p' internal/daemon/daemon.go
```

Repository: marcioaltoe/roundfix

Length of output: 1414

---



</details>

**Handle external spec-root paths before calling `GitCommitter`.** `settleArtifactCommitPath` returns an absolute path when `plan.specsRoot` is outside `plan.workDir`, and `GitCommitter` passes `req.Paths` straight to `git add -- ...`; that pathspec will fail from `WorkDir`, so settle breaks for external spec roots. The new `TestRunSettleUsesConfiguredExternalSpecRoot` only covers `fakeCommitter` and misses this path.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/settle.go` around lines 301 - 307, `settleArtifactCommitPath`
and the `settle` flow need to handle specs roots outside `plan.workDir` before
reaching `GitCommitter`, because absolute paths will be passed straight into
`git add` and fail from the worktree. Update the logic around
`settleArtifactCommitPath`, `ensureSettleCommitPath`, and the commit request
construction so external spec-root paths are converted into a pathspec Git can
add from `WorkDir`, and add coverage in
`TestRunSettleUsesConfiguredExternalSpecRoot` for the real `GitCommitter` path,
not just `fakeCommitter`.
```

</details>

<!-- fingerprinting:phantom:poseidon:beignet -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:074c00a6b1196f5a56309d29 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes:
  - Settle now filters commit paths through the Daemon stageable-path rules before calling `GitCommitter`, so external Spec Root task files are omitted while repository changes still commit; the regression now uses the real Git committer.
  - Verification: configured command `make verify` was run as `rtk make verify` and passed: Go tests, skills check, and build completed.
