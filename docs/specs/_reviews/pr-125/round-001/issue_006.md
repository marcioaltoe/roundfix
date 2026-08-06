---
source: coderabbit
pr: "125"
round: 1
round_created_at: "2026-08-05T19:38:23Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0078-roundfix-asks-for-the-review
head_sha: 1384d928c1d0af4d0bca06e506eb8b2953f9f341
file: docs/specs/_archived/0078-roundfix-asks-for-the-review/task_01.md
line: 77
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WxMJT,comment:PRRC_kwDOS0qyts7d790h
review_hash: aa644a28e4d669c7babd1b182818dd2d7ea9c827a1abdf0afe77270337af0648
duplicate_of: ""
source_review_id: "4868102567"
source_review_submitted_at: "2026-08-05T19:37:35Z"
---

# Issue 006: _ Functional Correctness_ _ Major_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Preserve failure exit codes in task verification commands.**

The pipelines use `grep` as the pipeline result or use `&& exit 1 || exit 0`. A failed test or forbidden changed path can therefore report exit 0.

- `docs/specs/_archived/0078-roundfix-asks-for-the-review/task_01.md#L71-L77`: run `go test` directly, or use `set -o pipefail` before filtering verbose output.
- `docs/specs/_archived/0078-roundfix-asks-for-the-review/task_02.md#L68-L75`: preserve `go test` failure status and replace the unconditional scope-check expression.
- `docs/specs/_archived/0078-roundfix-asks-for-the-review/task_03.md#L71-L79`: preserve `go test` failure status and replace the unconditional scope-check expression.
- `docs/specs/_archived/0078-roundfix-asks-for-the-review/task_04.md#L66-L76`: preserve `go test` failure status and replace the unconditional scope-check expression.
- `docs/specs/_archived/0078-roundfix-asks-for-the-review/task_05.md#L71-L74`: use an explicit `if` condition so an out-of-scope path exits non-zero.

As per coding guidelines, task verification commands must be effect-proving and able to fail when the required work is absent.

<details>
<summary>📍 Affects 5 files</summary>

- `docs/specs/_archived/0078-roundfix-asks-for-the-review/task_01.md#L71-L77` (this comment)
- `docs/specs/_archived/0078-roundfix-asks-for-the-review/task_02.md#L68-L75`
- `docs/specs/_archived/0078-roundfix-asks-for-the-review/task_03.md#L71-L79`
- `docs/specs/_archived/0078-roundfix-asks-for-the-review/task_04.md#L66-L76`
- `docs/specs/_archived/0078-roundfix-asks-for-the-review/task_05.md#L71-L74`

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/specs/_archived/0078-roundfix-asks-for-the-review/task_01.md` around
lines 71 - 77, Preserve verification failure exit codes across the task
documents: in docs/specs/_archived/0078-roundfix-asks-for-the-review/task_01.md
lines 71-77, retain go test failures by running directly or enabling pipefail;
in task_02.md lines 68-75, task_03.md lines 71-79, and task_04.md lines 66-76,
preserve go test status and replace the unconditional scope-check expressions;
in task_05.md lines 71-74, use an explicit conditional that exits non-zero for
out-of-scope paths.
```

</details>

<!-- consolidated_sites_start -->
<!--
<consolidated_sites>
<site>
<role>anchor</role>
<file>docs/specs/_archived/0078-roundfix-asks-for-the-review/task_01.md</file>
<line_range>71-77</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/specs/_archived/0078-roundfix-asks-for-the-review/task_02.md</file>
<line_range>68-75</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/specs/_archived/0078-roundfix-asks-for-the-review/task_03.md</file>
<line_range>71-79</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/specs/_archived/0078-roundfix-asks-for-the-review/task_04.md</file>
<line_range>66-76</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/specs/_archived/0078-roundfix-asks-for-the-review/task_05.md</file>
<line_range>71-74</line_range>
</site>
</consolidated_sites>
-->
<!-- consolidated_sites_end -->

<!-- fingerprinting:phantom:poseidon:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:7fa5004e1e33af582527f931 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Replaced authored verification pipelines whose final filter could hide the primary command's status with direct `go test` invocations and `git diff --quiet` scope checks across tasks 01-05.
- Evidence: Static inspection confirms the affected task gates now expose their own exit status; `rtk git diff --check` exited 0. The Daemon retains ownership of the authored full Verification command.
