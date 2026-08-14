---
source: coderabbit
pr: "130"
round: 1
round_created_at: "2026-08-06T03:34:01Z"
status: invalid
terminal_reason: "Roundfix verifies the dirty Task worktree before the Daemon commits, so the proposed post-commit premise is incorrect."
head_repository: marcioaltoe/roundfix
head_branch: ma/0065-loop-order-and-verification-honesty
head_sha: 7d35358ba9f77ceeda86ec5c34d7c4485a7eb8f9
file: docs/specs/_archived/0065-loop-order-and-verification-honesty/task_07.md
line: 88
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6W2v6J,comment:PRRC_kwDOS0qyts7eEK7e
review_hash: 9e47c9c30c2af248da5e029f2fbc8e93d97e3a02f454aecb8d5c20fa37ab7e3c
duplicate_of: ""
source_review_id: "4870656161"
source_review_submitted_at: "2026-08-06T03:12:28Z"
---

# Issue 017: _ Functional Correctness_ _ Major_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Verify the committed Task change, not only uncommitted worktree changes.**

`git diff --name-only HEAD` compares the worktree with `HEAD`. After the Task is committed, it returns no paths and can pass even if the Task commit changed a `.go` file.

Compare the immutable Task commit or its base-to-Task range with `git diff-tree --no-commit-id --name-only -r` or an equivalent command, then reject `.go` paths.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/specs/_archived/0065-loop-order-and-verification-honesty/task_07.md`
around lines 87 - 88, Update the verification command in the Task documentation
to inspect the committed Task change rather than only the current worktree. Use
git diff-tree or an equivalent base-to-Task commit range to list changed paths,
and reject the result when any .go file is present.
```

</details>

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:6070f3ff7f6ed6ae83557719 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: ADR-0057 and the Roundfix Batch contract place declared Task
  Verification before terminal settlement and commit. At that point
  `git diff --name-only HEAD` inspects the Task's uncommitted worktree exactly
  as intended. The reviewer assumes the command runs after commit; it does not.
  The completed archived Task also cannot be rewritten.
- Daemon Verification: `make verify` not run; Daemon-owned.
