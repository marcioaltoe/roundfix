---
source: coderabbit
pr: "143"
round: 1
round_created_at: "2026-08-09T01:14:50Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/specs-0082-0083
head_sha: 49ad4407ca050772b13a812b9625cf45d940e19e
file: docs/specs/0082-the-manifest-already-answered-that/task_01.md
line: 67
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XiAoc,comment:PRRC_kwDOS0qyts7fC8Qs
review_hash: f2273a685b990314824c5d95cc82d3ef969f07b70cdd2e9f5ddcbdc4cf266fb8
duplicate_of: ""
source_review_id: "4890173306"
source_review_submitted_at: "2026-08-09T00:28:48Z"
---

# Issue 031: _ Functional Correctness_ _ Major_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Make every Task boundary check inspect the complete worktree.**

`git diff --name-only HEAD` does not list untracked files. The 0082 command also limits the scan to `internal/`. An unauthorized path can therefore evade the declared allowlist.

- `docs/specs/0082-the-manifest-already-answered-that/task_01.md#L67-L67`: remove the `internal/` scope and compare all tracked and untracked paths with the declared test and test-data boundary.
- `docs/specs/0083-a-gate-that-can-say-no/task_01.md#L82-L82`: add `git ls-files --others --exclude-standard` to the existing allowlist check.
- `docs/specs/0083-a-gate-that-can-say-no/task_02.md#L70-L70`: add untracked paths to the allowlist check.
- `docs/specs/0083-a-gate-that-can-say-no/task_03.md#L63-L63`: add untracked paths to the allowlist check.
- `docs/specs/0083-a-gate-that-can-say-no/task_04.md#L65-L65`: add untracked paths to the allowlist check.
- `docs/specs/0083-a-gate-that-can-say-no/task_05.md#L62-L62`: add untracked paths to the allowlist check.
- `docs/specs/0083-a-gate-that-can-say-no/task_06.md#L59-L59`: add untracked paths to the allowlist check.

As per coding guidelines, Task Verification must use a complete, reproducible changed-file check.

<details>
<summary>📍 Affects 7 files</summary>

- `docs/specs/0082-the-manifest-already-answered-that/task_01.md#L67-L67` (this comment)
- `docs/specs/0083-a-gate-that-can-say-no/task_01.md#L82-L82`
- `docs/specs/0083-a-gate-that-can-say-no/task_02.md#L70-L70`
- `docs/specs/0083-a-gate-that-can-say-no/task_03.md#L63-L63`
- `docs/specs/0083-a-gate-that-can-say-no/task_04.md#L65-L65`
- `docs/specs/0083-a-gate-that-can-say-no/task_05.md#L62-L62`
- `docs/specs/0083-a-gate-that-can-say-no/task_06.md#L59-L59`

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/specs/0082-the-manifest-already-answered-that/task_01.md` at line 67,
Update the Task Verification changed-file checks to inspect the complete
worktree, including tracked changes and untracked paths while allowing only the
declared test and test-data exclusions. In
docs/specs/0082-the-manifest-already-answered-that/task_01.md:67-67, remove the
internal/ restriction; in
docs/specs/0083-a-gate-that-can-say-no/task_01.md:82-82, task_02.md:70-70,
task_03.md:63-63, task_04.md:65-65, task_05.md:62-62, and task_06.md:59-59,
incorporate git ls-files --others --exclude-standard into each existing
allowlist check.
```

</details>

<!-- consolidated_sites_start -->
<!--
<consolidated_sites>
<site>
<role>anchor</role>
<file>docs/specs/0082-the-manifest-already-answered-that/task_01.md</file>
<line_range>67-67</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/specs/0083-a-gate-that-can-say-no/task_01.md</file>
<line_range>82-82</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/specs/0083-a-gate-that-can-say-no/task_02.md</file>
<line_range>70-70</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/specs/0083-a-gate-that-can-say-no/task_03.md</file>
<line_range>63-63</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/specs/0083-a-gate-that-can-say-no/task_04.md</file>
<line_range>65-65</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/specs/0083-a-gate-that-can-say-no/task_05.md</file>
<line_range>62-62</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/specs/0083-a-gate-that-can-say-no/task_06.md</file>
<line_range>59-59</line_range>
</site>
</consolidated_sites>
-->
<!-- consolidated_sites_end -->

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:b3db31d2f4b2ec7bc49b6545 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `resolved`
- Notes: Updated all 7 task files. `0082/task_01.md` line 67: removed `internal/` scope restriction and added `git ls-files --others --exclude-standard` for untracked files. All six `0083/task_0[1-6].md` files: replaced `git diff --name-only HEAD |` with `(git diff --name-only HEAD; git ls-files --others --exclude-standard) |` to include untracked non-gitignored files in the boundary check. The `--exclude-standard` flag respects `.gitignore` so build artifacts and gitignored files remain excluded.
