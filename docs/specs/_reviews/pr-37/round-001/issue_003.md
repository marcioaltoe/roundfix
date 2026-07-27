---
source: coderabbit
pr: "37"
round: 1
round_created_at: "2026-07-27T01:53:02Z"
status: invalid
head_repository: marcioaltoe/roundfix
head_branch: ma/0036-doctor-skill-readiness
head_sha: 9a6b7f9433b9779afe75f38d833b780ceb2555ed
file: .roundfixrc.yml
line: 20
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6T6oNW,comment:PRRC_kwDOS0qyts7Zy0Nv
review_hash: d67c146ea165d91e2d59a01280b64b3f03a4fe346fe43dc4d8e82629e2f6e743
duplicate_of: ""
source_review_id: "4783144632"
source_review_submitted_at: "2026-07-27T01:52:02Z"
---

# Issue 003: _ Maintainability & Code Quality_ _ Major_ _ Heavy lift_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟠 Major_ | _🏗️ Heavy lift_

**Move concrete model fallbacks out of the project config.**

These entries bind Roundfix project profiles to the runtime-owned `gpt-5.6-terra` model. Keep profiles capability-oriented and resolve provider/model fallback selection from runtime-owned configuration instead.

As per coding guidelines, “Do not store runtime-owned model configuration, credentials, or adapter settings in Roundfix Project or User Config.”






Also applies to: 29-29, 47-47, 56-56

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In @.roundfixrc.yml at line 20, Remove the concrete gpt-5.6-terra model
assignments from all affected Roundfix profiles in .roundfixrc.yml, including
the entries referenced at lines 20, 29, 47, and 56. Keep the profiles
capability-oriented and rely on runtime-owned configuration for provider/model
fallback resolution.
```

</details>

<!-- fingerprinting:phantom:poseidon:terra -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:3031787b79eb761bfd064c1a -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: The requested removal contradicts Roundfix's accepted Agent Selection
  architecture. `CONTEXT.md` defines every Agent Selection as an exact ACP
  Runtime, Agent Model, and reasoning-effort tuple, and defines each Agent
  Selection Profile as one Preferred Selection plus a non-empty Fallback
  Chain. ADR-0049 explicitly states that exact runtime, official model
  identifier, and reasoning effort remain Roundfix-owned and that Project
  Config may provide complete higher-precedence profiles. The concrete Terra
  fallback entries therefore implement the repository's authoritative,
  reproducible Project Config policy; removing only their model fields would
  make the profiles incomplete rather than capability-oriented.

## Verification

- Inspected `CONTEXT.md` definitions for Agent Selection, Agent Selection
  Profile, Fallback Chain, and Default Agent Model.
- Inspected accepted ADR-0049 and ADR-0055; both require exact configured
  tuples and reject inheritance from runtime-local model configuration.
- Inspected commit `8564e67`, which deliberately updated the four project
  fallback selections to `gpt-5.6-terra`.
