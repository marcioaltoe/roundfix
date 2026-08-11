---
source: coderabbit
pr: "87"
round: 1
round_created_at: "2026-08-03T16:19:44Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/qa-task-not-flag
head_sha: d4011aade56592440d7a542682ebee7dced31f79
file: .agents/skills/roundfix/SKILL.md
line: 1259
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WCzY4,comment:PRRC_kwDOS0qyts7c4SNY
review_hash: 520e1b41b6dc34bcf0568026a8bf1da07f533be34fb7402d0aeca489711a30a4
duplicate_of: ""
source_review_id: "4846253969"
source_review_submitted_at: "2026-08-03T16:15:59Z"
---

# Issue 001: _ Maintainability & Code Quality_ _ Minor_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟡 Minor_ | _⚡ Quick win_

**Remove the stale per-run `qa` profile instruction.**

The file still says at Lines 455-456 that `implement` adds `qa` “only when requested.” The updated contract removes per-run QA selection and derives the gate from Task Graph metadata. Update the earlier sentence so agents do not look for the deleted `--qa` request.

<details>
<summary>Proposed fix</summary>

```diff
-`fetch` remains Agent-free. `resolve` and `watch` use only `review`; `implement` uses the Task categories and adds `qa` only when requested.
+`fetch` remains Agent-free. `resolve` and `watch` use only `review`; `implement` derives the `qa` category from Task Graph metadata and does not accept a per-run QA selection.
```
</details>

<details>
<summary>🧰 Tools</summary>

<details>
<summary>🪛 SkillSpector (2.4.4)</summary>

[warning] 24: [RP1] null: npx commands without a version suffix (e.g. `@1.0.0`) create a rug-pull risk if the upstream server is compromised and publishes a malicious update.

Remediation: Pin the version: npx `@scope/server`@1.2.3

(MCP Rug Pull (RP1))

---

[warning] 950: [RA2] Session Persistence: Skill establishes unauthorized persistence across sessions via cron jobs, startup scripts, or state files. Session persistence allows an attacker to maintain access beyond the current interaction.

Remediation: Remove any persistence mechanisms (cron jobs, startup scripts, state files). Skills should not maintain state across sessions without explicit user consent.

(Rogue Agent (RA2))

</details>

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In @.agents/skills/roundfix/SKILL.md around lines 1253 - 1259, Update the
earlier implement instruction around the “only when requested” wording to remove
per-run QA selection and the deleted --qa request. State that QA is determined
by the Task Graph’s authored qa metadata, consistent with the contract described
near the Daemon gate behavior.
```

</details>

<!-- fingerprinting:phantom:poseidon:luna -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:7e8f18abc8e0833e65e3002b -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes:
  - The profile-readiness section still said `implement` added `qa` only when requested, contradicting ADR-0088 and the later Implement Command contract that removes per-run QA selection.
  - Updated the canonical and distributed Roundfix Skill to derive categories, including `qa`, from Task Graph metadata.
  - `rtk make baseline-digests`: passed and regenerated only the ADR-0081 deterministic fallout from the authorized Roundfix-owned Skill edit.
  - Verification Feedback repair: regenerated the four Baseline plan-characterization goldens through their canonical `-update-baseline-plan-characterization` writer after the Skill edit changed the embedded catalog digest; `rtk env GOCACHE=/private/tmp/roundfix-batch001-feedback-gocache go test ./internal/baseline -run '^TestBaselinePlanCharacterization$' -count=1` passed afterward.
  - `rtk make skills-sync-check`: passed.
  - `rtk /usr/bin/diff -r .agents/skills/roundfix skills/roundfix`: passed.
  - Daemon Verification `make verify` was not run by this Agent; the Daemon owns authoritative Verification after this turn.
