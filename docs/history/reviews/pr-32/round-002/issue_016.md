---
source: coderabbit
pr: "32"
round: 2
round_created_at: "2026-07-17T13:23:47Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/release-plan-and-agent-selection-profiles
head_sha: d7ab1933ac9fdcf0c94d73e2f417d99d38e43fe7
file: skills/write-tasks/SKILL.md
line: 20
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Ru5tk,comment:PRRC_kwDOS0qyts7Wt96N
review_hash: e845c673dbc52760671df497c85463c65b17d362b2be510fc1efa27b5b7a2906
duplicate_of: ""
source_review_id: "4721765481"
source_review_submitted_at: "2026-07-17T10:25:31Z"
---

# Issue 016: _ Functional Correctness_ _ Major_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Use `AskUserQuestion` for the tech-spec bypass decision.**

This new prerequisite requires explicit user acceptance, but does not specify the repository-required interaction tool. State that the agent must use `AskUserQuestion` before proceeding without `_techspec.md`.

As per coding guidelines, “Use the AskUserQuestion tool for confirmations, clarifications, decision points, and needed user interaction; do not guess when an inexpensive user decision is required.”

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@skills/write-tasks/SKILL.md` around lines 19 - 20, Update the prerequisite
guidance around the missing _techspec.md decision to require using
AskUserQuestion to obtain explicit user acceptance before proceeding without the
tech spec. Only continue after acceptance, while preserving the existing
requirement for deeper codebase exploration.
```

</details>

<!-- fingerprinting:phantom:triton:luna -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:063a74b0d6ff7edcf8efe2c3 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Updated canonical and embedded `write-tasks` skill guidance to require `AskUserQuestion` for explicit user acceptance before proceeding without `_techspec.md`, then refreshed setup-context skill snapshot digests affected by the skill audit. Evidence: `GOCACHE=/private/tmp/roundfix-go-build rtk go test ./internal/agent ./internal/cli ./internal/config ./internal/daemon ./internal/releaseplan ./internal/spec ./internal/store ./internal/tui` passed; verification feedback repair evidence: `rtk make setup-context-check && rtk make skills-sync-check && rtk make skills-check` passed.
