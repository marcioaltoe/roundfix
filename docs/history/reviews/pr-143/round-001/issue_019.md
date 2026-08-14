---
source: coderabbit
pr: "143"
round: 1
round_created_at: "2026-08-09T01:14:50Z"
status: failed
head_repository: marcioaltoe/roundfix
head_branch: ma/specs-0082-0083
head_sha: 49ad4407ca050772b13a812b9625cf45d940e19e
file: docs/agents/docs-layout.md
line: 18
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XeKK7,comment:PRRC_kwDOS0qyts7e9kiv
review_hash: 6b71e83fbc1c46e193a7e8eb7be8eaa32c19b6d1c840838e2e13ccab16e473fc
duplicate_of: ""
terminal_reason: 'Verification failed: command "make verify" exited with exit status 2; diagnostics: /Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20260809T011807Z_77466ac6469b46dc/verification/batch-001-attempt-2.log'
source_review_id: "4888818931"
source_review_submitted_at: "2026-08-08T12:40:11Z"
---


# Issue 019: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Route observations by owner namespace.**

Line 10 sends every observation to `inbox/roundfix/`. The capture contract routes each entry to the project namespace that owns the fix, and the origin and destination can differ. Limit this statement to Roundfix-owned observations or use `inbox/<destination>/`; otherwise triage can misassign another project's observation.

Based on learnings, capture entries must use the namespace of the project that owns the fix, not necessarily the project where the session runs.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/agents/docs-layout.md` around lines 10 - 18, Update the “Capture does
not start here” guidance to limit the inbox/roundfix/ routing claim to
Roundfix-owned observations, or generalize it to the owning project’s
inbox/<destination>/ namespace. Preserve the requirement that each observation
is captured under the project namespace responsible for the fix, even when the
session runs elsewhere.
```

</details>

<!-- fingerprinting:phantom:poseidon:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:0a09dc5212b05738302321ae -->

_Source: Learnings_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Updated docs-layout.md to generalize observation routing to the owning project's namespace (inbox/<destination>/) instead of hardcoding inbox/roundfix/.
