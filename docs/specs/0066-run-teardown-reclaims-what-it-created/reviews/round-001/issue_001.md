---
source: coderabbit
pr: "116"
round: 1
round_created_at: "2026-08-05T05:03:08Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/spec-0066-implementation
head_sha: 76774c31db4c686cac700bc83af0bcc68521b6c6
file: .agents/skills/roundfix/SKILL.md
line: 1019
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Wi9Xn,comment:PRRC_kwDOS0qyts7dnSa-
review_hash: b9648f86550292b06ce58cd049f1770cd9713caab2199f11bd9c257195cfdbf2
duplicate_of: ""
source_review_id: "4861144519"
source_review_submitted_at: "2026-08-05T05:01:43Z"
---

# Issue 001: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Keep the current evidence branch in normal preflight handling.**

The phrase “leaves their Git refs unchanged, and disregards them” applies to both the current evidence branch and the superseded branches. Only proven-superseded branches should leave the actionable set. The current evidence branch must remain subject to the existing fast-forward integration or refusal path. Otherwise, the review Run can proceed with current QA or code work left unintegrated.

<details>
<summary>Proposed wording</summary>

```diff
-  A failed-cycle set proven from its QA Reports is never integrated automatically: preflight
-  reports the current evidence branch and each branch proven superseded by it,
-  leaves their Git refs unchanged, and disregards them so the review Run can
-  proceed. Reclaim them separately with `roundfix reconcile --apply`.
+  A failed-cycle set proven from its QA Reports filters only branches proven
+  superseded by the current evidence branch: preflight reports the current
+  evidence branch and each superseded branch, leaves superseded Git refs
+  unchanged, and removes only those superseded branches from actionable work.
+  The current evidence branch remains subject to normal integration or
+  refusal handling. Reclaim superseded branches separately with
+  `roundfix reconcile --apply`.
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
  automatically and journaled before the review Run continues. A failed-cycle
  set proven from its QA Reports filters only branches proven
  superseded by the current evidence branch: preflight reports the current
  evidence branch and each superseded branch, leaves superseded Git refs
  unchanged, and removes only those superseded branches from actionable work.
  The current evidence branch remains subject to normal integration or
  refusal handling. Reclaim superseded branches separately with
  `roundfix reconcile --apply`.
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

[error] 1499: [PE3] Credential Access: Code accesses credential files (SSH keys, AWS credentials, etc.). This could indicate credential theft attempts.

Remediation: Remove references to credential paths. Use environment variables or secrets managers. For docs, use placeholder paths (e.g., /path/to/config). Never load .env or token files in production code paths.

(Privilege Escalation (PE3))

---

[error] 1499: [PE3] Credential Access: Code accesses credential files (SSH keys, AWS credentials, etc.). This could indicate credential theft attempts.

Remediation: Remove references to credential paths. Use environment variables or secrets managers. For docs, use placeholder paths (e.g., /path/to/config). Never load .env or token files in production code paths.

(Privilege Escalation (PE3))

---

[warning] 1005: [RA2] Session Persistence: Skill establishes unauthorized persistence across sessions via cron jobs, startup scripts, or state files. Session persistence allows an attacker to maintain access beyond the current interaction.

Remediation: Remove any persistence mechanisms (cron jobs, startup scripts, state files). Skills should not maintain state across sessions without explicit user consent.

(Rogue Agent (RA2))

</details>

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In @.agents/skills/roundfix/SKILL.md around lines 1015 - 1019, Update the
preflight behavior description so the current evidence branch remains in the
normal actionable set and continues through the existing fast-forward
integration or refusal path. Apply “leaves their Git refs unchanged and
disregards them” only to branches proven superseded by the current evidence
branch, which alone should be removed from actionable handling.
```

</details>

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:204de3c5ea0b3941de86c8bb -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes:
  - Confirmed `filterPendingRunWorkByTarget` incorrectly added the current evidence branch to `Disregarded`, which removed it from normal Branch Integrity integration/refusal handling.
  - Removed that current-branch disregard path, updated both Roundfix Skill copies, and regenerated the required deterministic Baseline digest fallout.
  - `TestBranchIntegrityPreflightWatchDisregardsOnlySupersededFailedQACycles` proves that only the three older branches are disregarded, the current branch is integrated, and every Run Branch ref remains unchanged.
  - Focused evidence: `rtk go test ./internal/store ./internal/worktree ./internal/cli` passed (1,247 tests); `rtk make baseline-digests` reported no remaining changes; `rtk make skills-sync-check` and `rtk /usr/bin/diff -r .agents/skills/roundfix skills/roundfix` passed.
  - Verification Feedback repair: the Skill-driven catalog digest change also invalidated four Baseline plan-characterization goldens. Regenerated them through the canonical `-update-baseline-plan-characterization` writer; `rtk go test ./internal/baseline -run '^TestBaselinePlanCharacterization$' -count=1` then passed (6 tests).
  - The Daemon owns authoritative `make verify` after this Agent turn.
