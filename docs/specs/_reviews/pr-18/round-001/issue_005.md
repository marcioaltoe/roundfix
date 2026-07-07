---
source: coderabbit
pr: "18"
round: 1
round_created_at: "2026-07-07T00:22:01Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/docs-and-skill-bundle
head_sha: 83c26c480d5fa9c22da307daeaa48d6acbe59c85
file: skills/write-techspec/SKILL.md
line: 23
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6OwE1Z,comment:PRRC_kwDOS0qyts7SkHEG
review_hash: 4721c6463f61c45dbc0173bbd2e6a4b864255784c3e0566ace21020d3b078210
duplicate_of: ""
source_review_id: "4640508670"
source_review_submitted_at: "2026-07-06T23:41:29Z"
---

# Issue 005: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Align the minimal PRD contract with the template.**

`skills/write-prd/references/prd-template.md` doesn't define an `Acceptance Criteria` section, so refactor/bug-fix PRDs written from this path can end up out of sync with the template. Either add that section to the template or replace this instruction with the existing PRD sections.





<details>
<summary>📦 Proposed fix</summary>

```diff
- write a minimal _prd.md carrying only the contract downstream skills parse: the frontmatter (`spec`, `status: active`, `surfaces`), a problem statement, acceptance criteria, and non-goals
+ write a minimal _prd.md carrying only the contract downstream skills parse: the frontmatter (`spec`, `status: active`, `surfaces`), a problem statement, goals, core features, and non-goals
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
- **Refactor or bug fix** (no product behavior change) — this skill is the pipeline entry point; no PRD interview happens. When no spec folder exists yet, mint one following `write-prd`'s numbering rule (`docs/specs/NNNN-<kebab-slug>/`, scanning both `docs/specs/` and `docs/specs/_archived/` for the highest prefix), and write a **minimal `_prd.md`** carrying only the contract downstream skills parse: the frontmatter (`spec`, `status: active`, `surfaces`), a problem statement, goals, core features, and non-goals — engineering-framed, a few lines each. It exists so `write-tasks`, `qa-gate`, and `archive-spec` keep a single artifact contract; it is not a product document. If the "refactor" turns out to change product behavior, stop and route to `write-prd`.
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@skills/write-techspec/SKILL.md` around lines 22 - 23, The minimal PRD
contract in write-techspec is out of sync with the PRD template because it
պահանջs an Acceptance Criteria section that prd-template.md does not define.
Update the write-prd template or adjust the write-techspec guidance so both use
the same section set; make sure the instruction tied to the write-techspec entry
point still matches the sections emitted by the PRD path and references the
existing write-prd template and write-techspec contract.
```

</details>

<!-- fingerprinting:phantom:triton:quartz -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:a5bb7e30c64d47847c07b8a9 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `valid`
- Notes: The write-techspec minimal PRD guidance named an `acceptance criteria` section that the PRD template does not define.

## Resolution

- Updated write-techspec guidance to use the PRD template sections: problem statement, goals, core features, and non-goals.
- Updated the setup-workflow routing copy of the same minimal PRD contract to keep generated workflow guidance consistent.
- Verification: `rtk make verify` — passed. Go tests passed (`821` tests across `18` packages), `roundfix skills check` passed, and the CLI build succeeded.
