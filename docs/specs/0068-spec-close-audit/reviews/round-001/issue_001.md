---
source: coderabbit
pr: "113"
round: 1
round_created_at: "2026-08-05T02:12:07Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/spec-0068-implementation
head_sha: c9af2617f988bd63e1bd8f22c6178758a8e5fd40
file: .agents/skills/roundfix/SKILL.md
line: 641
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WhZp8,comment:PRRC_kwDOS0qyts7dk96i
review_hash: 2fe06d188447df02395b156d877e6c303f97d96a08015482898c359c4321591b
duplicate_of: ""
source_review_id: "4860420451"
source_review_submitted_at: "2026-08-05T02:11:26Z"
---

# Issue 001: _ Maintainability & Code Quality_ _ Minor_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟡 Minor_ | _⚡ Quick win_

**Correct the `preserved` description.**

The code assigns `preserved` in several deliberate cases that are not classification failures:

- An Active Run owns the survivor (`activeRunSurvivor`, `internal/specaudit/audit.go` lines 910-921).
- A scratch worktree branch is unpushed (`internal/specaudit/audit.go` lines 833-839).
- Two Run-less worktrees share one branch (`internal/specaudit/audit.go` lines 678-684).

The current wording states only that classification failed, so an operator can read a healthy Active Run survivor as an audit defect.




<details>
<summary>📝 Proposed wording</summary>

```diff
-| `preserved` | The audit could not classify the survivor, so it must remain intact. |
+| `preserved` | The survivor must remain intact: an Active Run owns it, its state is unpushed or shared, or the audit could not classify it. |
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
| `preserved` | The survivor must remain intact: an Active Run owns it, its state is unpushed or shared, or the audit could not classify it. |
```

</details>

<!-- suggestion_end -->

<details>
<summary>🧰 Tools</summary>

<details>
<summary>🪛 SkillSpector (2.5.1)</summary>

[warning] 24: [RP1] null: npx commands without a version suffix (e.g. `@1.0.0`) create a rug-pull risk if the upstream server is compromised and publishes a malicious update.

Remediation: Pin the version: npx `@scope/server`@1.2.3

(MCP Rug Pull (RP1))

---

[error] 1472: [PE3] Credential Access: Code accesses credential files (SSH keys, AWS credentials, etc.). This could indicate credential theft attempts.

Remediation: Remove references to credential paths. Use environment variables or secrets managers. For docs, use placeholder paths (e.g., /path/to/config). Never load .env or token files in production code paths.

(Privilege Escalation (PE3))

---

[error] 1472: [PE3] Credential Access: Code accesses credential files (SSH keys, AWS credentials, etc.). This could indicate credential theft attempts.

Remediation: Remove references to credential paths. Use environment variables or secrets managers. For docs, use placeholder paths (e.g., /path/to/config). Never load .env or token files in production code paths.

(Privilege Escalation (PE3))

---

[warning] 984: [RA2] Session Persistence: Skill establishes unauthorized persistence across sessions via cron jobs, startup scripts, or state files. Session persistence allows an attacker to maintain access beyond the current interaction.

Remediation: Remove any persistence mechanisms (cron jobs, startup scripts, state files). Skills should not maintain state across sessions without explicit user consent.

(Rogue Agent (RA2))

</details>

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In @.agents/skills/roundfix/SKILL.md at line 641, Update the `preserved` status
description in the audit status table to cover deliberate preservation cases,
including an Active Run survivor, an unpushed scratch worktree branch, and a
branch shared by two Run-less worktrees; do not describe it solely as a
classification failure.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:3fccac9f5dcbbe931fda3473 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `preserved` covers deliberate Active Run, unpushed, and shared-worktree outcomes in addition to failed classification, so the prior wording was incomplete.

## Resolution

- Updated the canonical Roundfix skill and generated mirror with the complete operator-facing meaning.
- Regenerated the sanctioned baseline digest and characterization fallout.
- Focused evidence: `rtk make skills-sync-check`, `rtk go run -buildvcs=false ./cmd/roundfix skills check`, both baseline characterization update tests, and `rtk git diff --check` exited 0.
- Daemon Verification: `make verify` was not run; the Daemon owns that command.
