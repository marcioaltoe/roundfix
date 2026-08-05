---
source: coderabbit
pr: "113"
round: 1
round_created_at: "2026-08-05T02:12:07Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/spec-0068-implementation
head_sha: c9af2617f988bd63e1bd8f22c6178758a8e5fd40
file: internal/worktree/worktree.go
line: 676
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WhZqX,comment:PRRC_kwDOS0qyts7dk97D
review_hash: f6bae48644588abf323061b43444fa4903902d49b254577c572a0e84e5cbdcea
duplicate_of: ""
source_review_id: "4860420451"
source_review_submitted_at: "2026-08-05T02:11:27Z"
---

# Issue 006: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/usr/bin/env bash
set -euo pipefail

repo="$(mktemp -d)"
trap 'rm -rf "$repo"' EXIT

git init -q "$repo"
git -C "$repo" config user.email reviewer@example.test
git -C "$repo" config user.name reviewer
git -C "$repo" checkout -qb main

printf 'base\n' >"$repo/removed.txt"
git -C "$repo" add removed.txt
git -C "$repo" commit -qm base

git -C "$repo" checkout -qb run
git -C "$repo" rm -q removed.txt
git -C "$repo" commit -qm 'delete tracked file'
run_head="$(git -C "$repo" rev-parse HEAD)"

git -C "$repo" checkout -q main
git -C "$repo" branch -D run

# Both current filters find no paths.
git -C "$repo" diff --quiet --no-renames --diff-filter=D "$run_head" main
git -C "$repo" diff --quiet --no-renames --diff-filter=MT "$run_head" main

# The default branch still has a file that the Run deleted.
test "$(git -C "$repo" diff --name-status "$run_head" main)" = $'A\tremoved.txt'
```

Repository: marcioaltoe/roundfix

Length of output: 192

---

🏁 Script executed:

```shell
sed -n '620,735p' internal/worktree/worktree.go
printf '\n-- related tests --\n'
rg -n -C 4 'compareRunContentToDefault|ReconciliationSafe|ReconciliationUnintegrated|differing shared|Run-only' --glob '*.go' .
```

Repository: marcioaltoe/roundfix

Length of output: 22442

---



</details>

**Detect Run-side deletions before setting `ReconciliationSafe`.**

When the Run Branch deletes a file that the default branch retains, `git diff runHead defaultHead` reports status `A`. The current filters select only `D`, `M`, and `T`, so the comparison can report no differences and incorrectly classify the Run as safe. Compare Run-side deletions with the merge base and classify retained default-branch paths as `ReconciliationUnintegrated`. Add a regression test.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/worktree/worktree.go` around lines 651 - 676, Update the
reconciliation logic before assigning ReconciliationSafe to detect files deleted
on the Run Branch but retained by the default branch, comparing Run-side
deletions against the merge base and treating matching default-branch paths as
differences. Include these paths in the existing unintegrated evidence and add a
regression test covering this deletion scenario.
```

</details>

<!-- fingerprinting:phantom:medusa:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:af8fa79b819976e3db1fba4c -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: A Run-side deletion retained by the default branch appears as `A` in the direct Run-to-default diff and bypassed both existing filters, allowing a false `safe` classification.

## Resolution

- Compared the Run Branch to its merge base for Run-side deletions, intersected those paths with the default tree, and treated retained deletions as explicit unintegrated evidence.
- Added a real Git regression fixture where the Run deletes a tracked file and the default branch retains it.
- Focused evidence: the new regression first failed with state `safe`; `rtk env GOCACHE=/private/tmp/roundfix-review0068-worktree-cache go test ./internal/worktree -run '^TestInspectTerminalRun(UnintegratedWhenDeletedTargetRetainsRunDeletedFile|SafeWhenTargetDeletedAfterSquashMerge|UnintegratedWhenDeletedTargetHasRunOnlyFile|UnintegratedWhenDeletedTargetHasDifferentSharedFile|UnintegratedWhenDeletedTargetContentComparisonFails)$' -count=1` exited 0 after the fix, and the full affected-package command exited 0.
- Daemon Verification: `make verify` was not run; the Daemon owns that command.
