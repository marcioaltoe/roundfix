---
source: coderabbit
pr: "110"
round: 1
round_created_at: "2026-08-04T22:55:35Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/spec-0070-implementation
head_sha: a588c6ca3ab9d977284ba1f9e80a89b0e6336786
file: .agents/skills/roundfix/SKILL.md
line: 1524
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WeYXx,comment:PRRC_kwDOS0qyts7dggqR
review_hash: 55900fbe607591c123da816c4de8de6aa07f82fd65b05c283ed54c82a65d8c54
duplicate_of: ""
source_review_id: "4859094834"
source_review_submitted_at: "2026-08-04T21:23:48Z"
---

# Issue 001: _ Functional Correctness_ _ Major_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Do not require `Clean` for the partial archive path.**

Lines 1257-1260 state that only `pass` ends a Run Clean and `partial` ends Unresolved. The `partial` case at Lines 1519-1524 can therefore never meet the stated advance condition. Users will enter recovery instead of archive.

Base both workflow transitions on `Archive` preconditions: every Task is completed and the QA Report is archive-eligible. State explicitly that an archive-eligible `partial` report can still have an Unresolved Run outcome.






Also applies to: 1581-1585

<details>
<summary>🧰 Tools</summary>

<details>
<summary>🪛 SkillSpector (2.5.1)</summary>

[warning] 24: [RP1] null: npx commands without a version suffix (e.g. `@1.0.0`) create a rug-pull risk if the upstream server is compromised and publishes a malicious update.

Remediation: Pin the version: npx `@scope/server`@1.2.3

(MCP Rug Pull (RP1))

---

[error] 1439: [PE3] Credential Access: Code accesses credential files (SSH keys, AWS credentials, etc.). This could indicate credential theft attempts.

Remediation: Remove references to credential paths. Use environment variables or secrets managers. For docs, use placeholder paths (e.g., /path/to/config). Never load .env or token files in production code paths.

(Privilege Escalation (PE3))

---

[error] 1439: [PE3] Credential Access: Code accesses credential files (SSH keys, AWS credentials, etc.). This could indicate credential theft attempts.

Remediation: Remove references to credential paths. Use environment variables or secrets managers. For docs, use placeholder paths (e.g., /path/to/config). Never load .env or token files in production code paths.

(Privilege Escalation (PE3))

---

[warning] 951: [RA2] Session Persistence: Skill establishes unauthorized persistence across sessions via cron jobs, startup scripts, or state files. Session persistence allows an attacker to maintain access beyond the current interaction.

Remediation: Remove any persistence mechanisms (cron jobs, startup scripts, state files). Skills should not maintain state across sessions without explicit user consent.

(Rogue Agent (RA2))

</details>

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In @.agents/skills/roundfix/SKILL.md around lines 1519 - 1524, Update the
Advance transition and its corresponding section around the second referenced
occurrence to base archiving on Archive preconditions: all Tasks completed and
the QA Report archive-eligible, without requiring a Clean Run. Explicitly allow
an archive-eligible partial report when the Run outcome remains Unresolved,
while preserving the declared-unreachable coverage and unproven-recording rules.
```

</details>

<!-- fingerprinting:phantom:medusa:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:d1e1a57d67fca99035cd62ed -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The user-facing loop required `Clean`, but the same skill states that an archive-eligible `partial` QA Report leaves the Run Unresolved. Both Advance transitions now require completed Tasks plus an archive-eligible QA Report, explicitly permit the Unresolved partial outcome, and preserve the `unproven` recording rule.
- Evidence: `rtk make skills-sync` regenerated the embedded skill; `rtk make baseline-digests` passed and regenerated the sanctioned derived pins; `rtk make skills-sync-check` passed.
- Verification feedback: `rtk go test ./internal/baseline -run '^TestBaselinePlanCharacterization$' -update-baseline-plan-characterization -count=1` regenerated the four dependent plan-characterization goldens, and the same command without the update flag passed with 6 tests.
