---
source: coderabbit
pr: "125"
round: 3
round_created_at: "2026-08-05T20:34:12Z"
status: duplicated
head_repository: marcioaltoe/roundfix
head_branch: ma/0078-roundfix-asks-for-the-review
head_sha: a89da452f019b880472c798f58529ea8aebefb1b
file: docs/specs/_archived/0078-roundfix-asks-for-the-review/task_04.md
line: 69
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WyExc,comment:PRRC_kwDOS0qyts7d9QPg
review_hash: fd95e74eda2a7d454a08ec7d2b58e22bb5908021b7b1c42c83592aa3829b4822
duplicate_of: /Users/marcio/dev/roundfix/docs/specs/_reviews/pr-125/round-004/issue_001.md
source_review_id: "4868508392"
source_review_submitted_at: "2026-08-05T20:33:32Z"
---



# Issue 001: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Make the YAML check prove the full file contract.**

The Ruby command checks parseability, duplicate keys, and selected values. It does not verify that `auto_review` is the first entry under `reviews:` or that every other byte in `.coderabbit.yaml` remains unchanged. Both are explicit requirements. Add an exact preimage/diff check for `.coderabbit.yaml` in addition to the structural validation.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/specs/_archived/0078-roundfix-asks-for-the-review/task_04.md` around
lines 67 - 69, Extend the validation command around load_unique and the
.coderabbit.yaml checks to compare the file against the exact required preimage,
proving auto_review remains the first reviews entry and all other bytes are
unchanged. Retain the existing single-document, duplicate-key, and nested-value
validations, and make the command fail when the exact content differs.
```

</details>

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:a92055a0d902166fce86ba6e -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `UNREVIEWED`
- Notes:
