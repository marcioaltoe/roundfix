---
source: coderabbit
pr: "125"
round: 3
round_created_at: "2026-08-05T20:34:12Z"
status: duplicated
head_repository: marcioaltoe/roundfix
head_branch: ma/0078-roundfix-asks-for-the-review
head_sha: a89da452f019b880472c798f58529ea8aebefb1b
file: docs/specs/_archived/0078-roundfix-asks-for-the-review/task_05.md
line: 74
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WyExe,comment:PRRC_kwDOS0qyts7d9QPj
review_hash: c4d612b6378afeed8abccb49547cb908924f8cb510c42482579304e163b90316
duplicate_of: /Users/marcio/dev/roundfix/docs/specs/_reviews/pr-125/round-004/issue_002.md
source_review_id: "4868508392"
source_review_submitted_at: "2026-08-05T20:33:32Z"
---



# Issue 002: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Include untracked paths in the scope check.**

`git diff --quiet` ignores untracked files. An unauthorized generated or temporary file can therefore pass the “only bounded paths changed” check. Use `git status --short --untracked-files=all` with matching path exclusions and fail when any out-of-scope path remains.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/specs/_archived/0078-roundfix-asks-for-the-review/task_05.md` around
lines 71 - 74, Update the scope-check instructions around the second git diff
command to include untracked files by using git status --short
--untracked-files=all. Apply equivalent exclusions for all allowed paths,
including the task file and bounded directories, and make the check fail when
any remaining path is outside that scope; leave the Go-source check unchanged.
```

</details>

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:1f19dd2edd02ce263b415a26 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `UNREVIEWED`
- Notes:
