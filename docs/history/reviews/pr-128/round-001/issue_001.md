---
source: coderabbit
pr: "128"
round: 1
round_created_at: "2026-08-06T03:35:45Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0069-review-run-targets-its-pull-request
head_sha: 62cd2ea6f84aa181570ef18f0e05225c6e4ebb88
file: .agents/skills/roundfix/SKILL.md
line: 1054
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6W1GuR,comment:PRRC_kwDOS0qyts7eBurO
review_hash: d75db9fe951286a22e59a2687eccb782d5e73e45499f3d1baff877ef8ca0e804
duplicate_of: ""
source_review_id: "4869925235"
source_review_submitted_at: "2026-08-06T00:16:28Z"
---

# Issue 001: _ Functional Correctness_ _ Minor_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**Handle revision-only mismatches in the recovery text.**

The contract checks both the checkout branch and the Pull Request head revision, but this section says refusal occurs only when branches differ. It also gives only `git switch` to the branch. If the branch matches but the checkout is at a different revision, the text does not explain the refusal or how to restore the expected revision. State both mismatch conditions and document the recovery action that restores both branch and revision.

<details>
<summary>🧰 Tools</summary>

<details>
<summary>🪛 SkillSpector (2.5.1)</summary>

[warning] 24: [RP1] null: npx commands without a version suffix (e.g. `@1.0.0`) create a rug-pull risk if the upstream server is compromised and publishes a malicious update.

Remediation: Pin the version: npx `@scope/server`@1.2.3

(MCP Rug Pull (RP1))

---

[error] 1561: [PE3] Credential Access: Code accesses credential files (SSH keys, AWS credentials, etc.). This could indicate credential theft attempts.

Remediation: Remove references to credential paths. Use environment variables or secrets managers. For docs, use placeholder paths (e.g., /path/to/config). Never load .env or token files in production code paths.

(Privilege Escalation (PE3))

---

[error] 1561: [PE3] Credential Access: Code accesses credential files (SSH keys, AWS credentials, etc.). This could indicate credential theft attempts.

Remediation: Remove references to credential paths. Use environment variables or secrets managers. For docs, use placeholder paths (e.g., /path/to/config). Never load .env or token files in production code paths.

(Privilege Escalation (PE3))

---

[warning] 1037: [RA2] Session Persistence: Skill establishes unauthorized persistence across sessions via cron jobs, startup scripts, or state files. Session persistence allows an attacker to maintain access beyond the current interaction.

Remediation: Remove any persistence mechanisms (cron jobs, startup scripts, state files). Skills should not maintain state across sessions without explicit user consent.

(Rogue Agent (RA2))

</details>

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In @.agents/skills/roundfix/SKILL.md around lines 1042 - 1054, Update the
Preflight Validation refusal text to cover mismatches in either the checkout
branch or the Pull Request head revision, not only differing branches. Expand
the recovery guidance after the diagnostic to restore both the named PR Head
Branch and its expected revision, including the appropriate command shape for a
revision-only mismatch, while preserving the no-side-effects behavior.
```

</details>

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:4fbe07b5c26ef1938b26c20c -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Preflight now treats a known PR-head revision mismatch as a target
  mismatch even when the branch names match. The canonical Roundfix Skill and
  its embedded mirror document branch-only, revision-only, and combined
  recovery commands; `make skills-sync-check` and `roundfix skills check`
  passed. Authoritative `make verify` remains Daemon-owned.
