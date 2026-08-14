---
source: coderabbit
pr: "136"
round: 3
round_created_at: "2026-08-06T20:20:19Z"
status: failed
terminal_reason: "The finding is valid, but its five targets are unassigned prior-round Review Issue files that this Batch is expressly forbidden to edit."
head_repository: marcioaltoe/roundfix
head_branch: ma/0079-one-door-for-fleet-knowledge
head_sha: fba018672a8f31a3a4f59e6afd21d2c03c6a220f
file: docs/specs/_reviews/pr-136/round-002/issue_001.md
line: 36
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XGgW7,comment:PRRC_kwDOS0qyts7ebLL3
review_hash: 73b5419ac5bd8929dfcb856daa059a1ecaad55791749d7df5ad319de5147a05a
duplicate_of: ""
source_review_id: "4877969817"
source_review_submitted_at: "2026-08-06T20:19:25Z"
---

# Issue 001: _ Maintainability & Code Quality_ _ Minor_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟡 Minor_ | _⚡ Quick win_

**Fix the shared Review Issue Markdown template.**

The five files repeat MD037 and MD040. `docs/specs/_reviews/pr-136/round-002/issue_004.md` also reports MD031. Normalize heading emphasis, label every prompt fence, and add blank lines around the affected fence.

- `docs/specs/_reviews/pr-136/round-002/issue_001.md#L21-L36`: remove heading spaces and label the prompt fence.
- `docs/specs/_reviews/pr-136/round-002/issue_002.md#L21-L38`: remove heading spaces and label the prompt fence.
- `docs/specs/_reviews/pr-136/round-002/issue_003.md#L21-L36`: remove heading spaces and label the prompt fence.
- `docs/specs/_reviews/pr-136/round-002/issue_004.md#L21-L37`: remove heading spaces and add blank lines around the affected fence.
- `docs/specs/_reviews/pr-136/round-002/issue_004.md#L73-L80`: label the prompt fence.
- `docs/specs/_reviews/pr-136/round-002/issue_005.md#L21-L36`: remove heading spaces and label the prompt fence.

<details>
<summary>🧰 Tools</summary>

<details>
<summary>🪛 markdownlint-cli2 (0.23.2)</summary>

[warning] 21-21: Spaces inside emphasis markers

(MD037, no-space-in-emphasis)

---

[warning] 21-21: Spaces inside emphasis markers

(MD037, no-space-in-emphasis)

---

[warning] 21-21: Spaces inside emphasis markers

(MD037, no-space-in-emphasis)

---

[warning] 36-36: Fenced code blocks should have a language specified

(MD040, fenced-code-language)

</details>

</details>

<details>
<summary>📍 Affects 5 files</summary>

- `docs/specs/_reviews/pr-136/round-002/issue_001.md#L21-L36` (this comment)
- `docs/specs/_reviews/pr-136/round-002/issue_002.md#L21-L38`
- `docs/specs/_reviews/pr-136/round-002/issue_003.md#L21-L36`
- `docs/specs/_reviews/pr-136/round-002/issue_004.md#L21-L37`
- `docs/specs/_reviews/pr-136/round-002/issue_004.md#L73-L80`
- `docs/specs/_reviews/pr-136/round-002/issue_005.md#L21-L36`

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/specs/_reviews/pr-136/round-002/issue_001.md` around lines 21 - 36,
Normalize the shared Review Issue Markdown formatting: in
docs/specs/_reviews/pr-136/round-002/issue_001.md (21-36), issue_002.md (21-38),
issue_003.md (21-36), and issue_005.md (21-36), remove heading spaces and label
each prompt fence; in issue_004.md (21-37), remove heading spaces and add blank
lines around the affected fence; in issue_004.md (73-80), label the prompt
fence.
```

</details>

<!-- consolidated_sites_start -->
<!--
<consolidated_sites>
<site>
<role>anchor</role>
<file>docs/specs/_reviews/pr-136/round-002/issue_001.md</file>
<line_range>21-36</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/specs/_reviews/pr-136/round-002/issue_002.md</file>
<line_range>21-38</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/specs/_reviews/pr-136/round-002/issue_003.md</file>
<line_range>21-36</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/specs/_reviews/pr-136/round-002/issue_004.md</file>
<line_range>21-37</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/specs/_reviews/pr-136/round-002/issue_004.md</file>
<line_range>73-80</line_range>
</site>
<site>
<role>sibling</role>
<file>docs/specs/_reviews/pr-136/round-002/issue_005.md</file>
<line_range>21-36</line_range>
</site>
</consolidated_sites>
-->
<!-- consolidated_sites_end -->

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:5cda19e256727fd3521ca5a2 -->

_Source: Linters/SAST tools_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Direct inspection confirms the affected Round 002 artifacts retain
  spaces inside emphasis markers and unlabeled prompt fences. The requested
  fix cannot be applied in this Batch because every target is an unassigned
  Review Issue file, and the Batch contract expressly forbids editing one.
- Focused evidence: a bounded `rtk rg` inspection over the nine Round 002
  issue files reported the cited headings and bare fences. No target file was
  changed.
- Daemon Verification: `make verify` not run; Daemon-owned.
