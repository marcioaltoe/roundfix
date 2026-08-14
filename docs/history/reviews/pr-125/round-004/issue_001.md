---
source: coderabbit
pr: "125"
round: 4
round_created_at: "2026-08-05T20:44:41Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0078-roundfix-asks-for-the-review
head_sha: 535c9dd97cb583f418deeca1bc639b5030e5e728
file: docs/specs/_archived/0078-roundfix-asks-for-the-review/task_04.md
line: 69
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WyExc,comment:PRRC_kwDOS0qyts7d9QPg
review_hash: 89ead044ada0d5503829fb9aa5f6449ac159e92ebf5286342e8fa816d9c25228
duplicate_of: ""
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

<!-- This is an auto-generated reply by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Task 04's YAML probe validated structure and selected values but did
  not prove the required byte-exact transformation. The Verification contract
  now derives the expected file from the Task's `HEAD` preimage, inserts the
  manual-review block immediately after `reviews:`, and rejects any other byte
  difference while retaining the duplicate-key and nested-value checks.
- Focused evidence: the historical `bb7a2067^` to `bb7a2067` transformation
  probe exited 0, and the retained structural/value Ruby probe exited 0 against
  the current configuration files. The Daemon owns `make verify`.
