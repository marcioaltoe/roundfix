---
source: coderabbit
pr: "125"
round: 4
round_created_at: "2026-08-05T20:44:41Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0078-roundfix-asks-for-the-review
head_sha: 535c9dd97cb583f418deeca1bc639b5030e5e728
file: docs/specs/_archived/0078-roundfix-asks-for-the-review/task_05.md
line: 74
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WyExe,comment:PRRC_kwDOS0qyts7d9QPj
review_hash: 8570490970d23dea79162786fd4ddc0e8b70e3c27e8a29fb138f12f3cab525fe
duplicate_of: ""
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

<!-- This is an auto-generated reply by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The tracked-only `git diff --quiet` scope check could not detect an
  unauthorized untracked path. Task 05 now checks the complete short status
  with `--untracked-files=all`, retains every bounded-path exclusion, disables
  the unreliable fsmonitor for the inspection, and requires empty output. The
  separate Go-source check remains unchanged.
- Focused evidence: the new scoped status command exited 0 and emitted the
  known untracked review artifacts plus modified out-of-scope Task files;
  therefore its documented empty-output assertion rejects those paths. The
  Daemon owns `make verify`.
