---
source: coderabbit
pr: "20"
round: 1
round_created_at: "2026-07-10T19:46:39Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/model-catalog-and-stream-efficiency
head_sha: f73fd40d026660e67999ceb7cbb016d7b1c039ad
file: README.md
line: 528
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6P_tl-,comment:PRRC_kwDOS0qyts7USTx7
review_hash: ed68ea047ce46def2a4415747dd64a0f105abad9031ca9252358ee903236a65d
duplicate_of: ""
source_review_id: "4674496937"
source_review_submitted_at: "2026-07-10T19:45:35Z"
---

# Issue 003: _ Functional Correctness_ _ Minor_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**Document mixed OpenCode configuration and flag overrides.**

Selection resolution is independent per value, so users can configure one OpenCode value and provide the other with a one-Run flag. The current “configure both or pass both” wording incorrectly implies that mixed configuration is unsupported. Update this wording here and in `skills/roundfix/SKILL.md`.

<details>
<summary>Suggested wording</summary>

```diff
-- OpenCode: no built-in model or reasoning value; configure both values or
-  pass both one-Run flags.
+- OpenCode: no built-in model or reasoning value; configure both values,
+  pass both one-Run flags, or combine configured values with one-Run flags.
```
</details>

   


Also applies to: 546-547

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@README.md` around lines 527 - 528, Update the OpenCode configuration guidance
in README.md and skills/roundfix/SKILL.md to state that model and reasoning
values resolve independently: users may configure either value and supply the
other via its one-Run flag, or provide both through configuration or flags.
Remove wording that requires both values to come from the same source.
```

</details>

<!-- fingerprinting:phantom:triton:luna -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:4d1b8fc0a2679ac3e21e0459 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `valid`
- Notes:
  - Updated README OpenCode guidance to match independent model and reasoning resolution: each value can come from User Config, Project Config, its matching one-Run flag, or a mix.
  - Updated both Roundfix skill copies so the embedded skill stays synced with the canonical `.agents/skills/roundfix` copy.
  - Verification: `make verify` passed in this session. Result: `go test ./...` reported 1046 passed in 19 packages, `roundfix skills check` passed, and `go build` completed.
