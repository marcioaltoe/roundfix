---
source: coderabbit
pr: "56"
round: 1
round_created_at: "2026-07-31T14:02:35Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0060-spec-owned-reference-lifecycle
head_sha: 05752e266533235d41a554f01dd42584bd24d46d
file: .agents/skills/write-idea/SKILL.md
line: 30
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Vb8yS,comment:PRRC_kwDOS0qyts7b_18I
review_hash: 766a9e6864a6993fc4f00e27b920dcf3b703e6381a307aebb85fef95d4275d04
duplicate_of: ""
source_review_id: "4829144282"
source_review_submitted_at: "2026-07-31T14:01:53Z"
---

# Issue 005: _ Functional Correctness_ _ Trivial_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🔵 Trivial_ | _⚡ Quick win_

**Make adopted-source links file-relative in every authorial skill.**

These rules require the owner's post-adoption copy but do not state that Markdown destinations must resolve relative to the file containing the link. Align all three rules with the explicit relative-link requirement in `.agents/skills/write-prd/SKILL.md`.

- `.agents/skills/write-idea/SKILL.md#L27-L30`: require `_idea.md` links to resolve relative to `_idea.md`.
- `.agents/skills/write-techspec/SKILL.md#L30-L33`: require `_techspec.md` links to resolve relative to `_techspec.md`.
- `.agents/skills/write-tasks/SKILL.md#L60-L64`: require task-file links to resolve relative to each task file.

Based on coding guidelines, Markdown destinations must resolve relative to the linking file.

<details>
<summary>📍 Affects 3 files</summary>

- `.agents/skills/write-idea/SKILL.md#L27-L30` (this comment)
- `.agents/skills/write-techspec/SKILL.md#L30-L33`
- `.agents/skills/write-tasks/SKILL.md#L60-L64`

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In @.agents/skills/write-idea/SKILL.md around lines 27 - 30, Update the
adopted-source rules in .agents/skills/write-idea/SKILL.md:27-30,
.agents/skills/write-techspec/SKILL.md:30-33, and
.agents/skills/write-tasks/SKILL.md:60-64 so each Markdown link targets the
owner's post-adoption copy and resolves relative to its linking file—_idea.md,
_techspec.md, or the individual task file respectively.
```

</details>

<!-- consolidated_sites_start -->
<!--
<consolidated_sites>
<site>
<role>anchor</role>
<file>.agents/skills/write-idea/SKILL.md</file>
<line_range>27-30</line_range>
</site>
<site>
<role>sibling</role>
<file>.agents/skills/write-techspec/SKILL.md</file>
<line_range>30-33</line_range>
</site>
<site>
<role>sibling</role>
<file>.agents/skills/write-tasks/SKILL.md</file>
<line_range>60-64</line_range>
</site>
</consolidated_sites>
-->
<!-- consolidated_sites_end -->

<!-- fingerprinting:phantom:triton:luna -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:e6b3cd54ab81e5c463b229dc -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: A post-adoption repository path alone is insufficient for Markdown;
  the destination must be computed from the file that contains the link.
- Fix: `write-idea`, `write-techspec`, and `write-tasks` now require the
  owner's post-adoption copy and make the destination relative to `_idea.md`,
  `_techspec.md`, or the individual Task file, respectively. `rtk make
  skills-sync` propagated all three canonical edits to their shipped copies.
- Focused evidence: `TestSpecReferenceLifecycleSkillContracts`, full
  `./skills` tests, and `skills-sync-check` passed; canonical and shipped Skill
  trees are synchronized.
- Daemon Verification: `make verify` not run; Daemon-owned.
