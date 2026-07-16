---
source: coderabbit
pr: "29"
round: 1
round_created_at: "2026-07-16T20:45:26Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/setup-context-driven-validator
head_sha: 49cdc07dcdf5b8fcb40eb459f27383b00995c0e3
file: skills/setup-context-driven/assets/templates/guides/autonomous-work.md
line: 5
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6RlYe9,comment:PRRC_kwDOS0qyts7WgS4U
review_hash: aacc89a67387039ccddb5003bcee887a39618cb78e17b0633d2e7da67953c544
duplicate_of: ""
source_review_id: "4717533168"
source_review_submitted_at: "2026-07-16T20:44:20Z"
---

# Issue 001: _ Functional Correctness_ _ Major_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Make ACP Runtime delegation unconditional.**

The phrase “while a runtime can execute the Task” permits the Supervisor to implement feature code or tests when no runtime is available. Autonomous sessions require the Supervisor to orchestrate only; remove that fallback.
As per coding guidelines, `For autonomous sessions, the Supervisor orchestrates only; implementation must be delegated to an ACP Runtime.`




<details>
<summary>Proposed wording</summary>

```diff
-The Supervisor does not write
-feature code or tests while a runtime can execute the Task.
+The Supervisor never writes feature code or tests. Implementation is always
+delegated to the selected ACP Runtime.
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
The Supervisor orchestrates and authors Specs. Implementation goes to the
selected ACP Runtime, normally through Roundfix. The Supervisor never writes
feature code or tests. Implementation is always delegated to the selected ACP
Runtime.
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@skills/setup-context-driven/assets/templates/guides/autonomous-work.md`
around lines 3 - 5, Update the Supervisor responsibility statement in
autonomous-work.md to make ACP Runtime delegation unconditional: remove the
qualification “while a runtime can execute the Task” and state that autonomous
sessions only allow orchestration, with all implementation delegated to an ACP
Runtime.
```

</details>

<!-- fingerprinting:phantom:poseidon:luna -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:2f333e0b0f52541e6e1ba581 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The template allowed Supervisor implementation fallback. Updated the autonomous-work guide template to require ACP Runtime delegation unconditionally.

## Resolution

- Status: `resolved`
- Evidence: `rtk make setup-context-check` passed; `rtk make skills-sync-check` passed; `rtk make verify` passed.
