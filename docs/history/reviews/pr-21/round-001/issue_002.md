---
source: coderabbit
pr: "21"
round: 1
round_created_at: "2026-07-11T15:14:27Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/context-efficient-runs
head_sha: d002674f43f98559f0029af16feeaa54173f85ad
file: skills/roundfix/SKILL.md
line: 457
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6QGNz3,comment:PRRC_kwDOS0qyts7UbgwM
review_hash: cba10b23c5c46367bb3aabf4a80b1778de1939ec25aea18cf7cedbcf41bd034d
duplicate_of: ""
source_review_id: "4677550318"
source_review_submitted_at: "2026-07-11T10:43:11Z"
---

# Issue 002: _ Functional Correctness_ _ Major_ _ Heavy lift_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _🏗️ Heavy lift_

**Restore legacy journal replay compatibility.**

Pre-0024 verification events lack `attempt`, but the new projector requires it, so `roundfix events <run-id>` exits `1` instead of replaying affected historical Runs. Normalize/migrate legacy verification events (or otherwise preserve replay) before documenting unconditional replay behavior.

<details>
<summary>🧰 Tools</summary>

<details>
<summary>🪛 SkillSpector (2.3.7)</summary>

[warning] 24: [RP1] null: npx commands without a version suffix (e.g. `@1.0.0`) create a rug-pull risk if the upstream server is compromised and publishes a malicious update.

Remediation: Pin the version: npx `@scope/server`@1.2.3

(MCP Rug Pull (RP1))

---

[warning] 545: [RA2] Session Persistence: Skill establishes unauthorized persistence across sessions via cron jobs, startup scripts, or state files. Session persistence allows an attacker to maintain access beyond the current interaction.

Remediation: Remove any persistence mechanisms (cron jobs, startup scripts, state files). Skills should not maintain state across sessions without explicit user consent.

(Rogue Agent (RA2))

</details>

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@skills/roundfix/SKILL.md` around lines 440 - 457, Update the events
replay/projector path so legacy pre-0024 verification events missing attempt are
normalized or migrated before projection, allowing affected historical Runs to
replay successfully. Preserve current validation for malformed newer payloads,
then document unconditional replay only after this compatibility behavior is
implemented around the Run projection logic.
```

</details>

<!-- fingerprinting:phantom:poseidon:terra -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:a51649ef25f3e1981827d58f -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Pre-repair `daemon.verification` journal payloads used the `started`, `failed`, and `passed` phases without an `attempt` field. The new projector required `attempt` unconditionally, so replay failed before emitting those historical records.

## Resolution

- Updated `internal/runevent/stream.go` to map the exact legacy verification payload and summary shapes to their single historical attempt while continuing to reject malformed attempt-aware events.
- Added projector-level coverage for all three legacy phases, a negative case for a malformed current payload, and CLI replay coverage for a historical event.
- Verification: `make verify` passed in this Batch: 1,103 tests passed in 19 packages, the Roundfix skill check passed, and the build completed.
