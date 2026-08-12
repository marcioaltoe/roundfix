---
source: coderabbit
pr: "157"
round: 1
round_created_at: "2026-08-12T01:25:35Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/what-an-agent-reads-before-it-decides
head_sha: bdc831f8de829f09257a71a04adca1b5219c6381
file: .agents/skills/write-idea/SKILL.md
line: 50
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6YbQdE,comment:PRRC_kwDOS0qyts7gSdxg
review_hash: 1ec993cfd9e3042d41b0abd7ab83fe3739e88416d831168e9fafbac7c94cf431
duplicate_of: ""
source_review_id: "4912178363"
source_review_submitted_at: "2026-08-12T01:24:11Z"
---

# Issue 002: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Do not enter update mode for an archived Spec.**

When the scan finds a matching folder under `_archived/specs/`, line 50 says to operate in update mode. That permits rewriting archived `_idea.md` or reusing an archived Spec number. Treat archived matches as read-only overlap findings. Enter update mode only for a matching active Spec, and assign a new number for new work.

Based on learnings: archived Specs are historical artifacts and must not be rewritten.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In @.agents/skills/write-idea/SKILL.md around lines 48 - 50, Update the
spec-folder rules to treat matching folders under _archived/specs/ as read-only
overlap findings: do not enter update mode, rewrite their _idea.md, or reuse
their numbers. Enter update mode only for matching active folders under
docs/specs/; otherwise assign a new number for the work while still surfacing
archived overlaps.
```

</details>

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:7f39c1aa27bca7dc9b20ff4f -->

_Source: Learnings_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `RESOLVED`
- Notes: Updated `.agents/skills/write-idea/SKILL.md` (and its `skills/` mirror via `make skills-sync`). Step 2 now enters update mode only for a matching **active** spec folder under the Spec Root; a matching folder under the archive root is treated as a read-only overlap finding — no update mode, no rewriting its `_idea.md`, no reusing its number — and new work gets a new number while the archived overlap is surfaced. Verified: `make skills-sync-check` and the full `make verify` gate pass.
