---
source: coderabbit
pr: "18"
round: 1
round_created_at: "2026-07-07T00:22:01Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/docs-and-skill-bundle
head_sha: 83c26c480d5fa9c22da307daeaa48d6acbe59c85
file: skills/archive-spec/SKILL.md
line: 38
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6OwE1L,comment:PRRC_kwDOS0qyts7SkHDy
review_hash: 1cd696bbf30ed07bed967894e94614e403fcb7b52ffebe09336d7c0dce04ebc8
duplicate_of: ""
source_review_id: "4640508670"
source_review_submitted_at: "2026-07-06T23:41:28Z"
---

# Issue 001: _ Data Integrity & Integration_ _ Minor_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟡 Minor_ | _⚡ Quick win_

**Preserve the QA override in the archive stamp.**

The override path is only mentioned in Preconditions; the frontmatter template never records `qa_override: true`. A literal run of the steps would archive a QA-bypassed spec without marking that exception.

<details>
<summary>Suggested fix</summary>

```diff
   status: archived
   archived: YYYY-MM-DD
+  qa_override: true # only when archiving despite failed/missing QA
   release: <tag or PR URL> # only when known — from --release or an already-merged PR/tag
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
2. **QA passed** — the newest report in `qa/` has `verdict: pass`. A missing `qa/` directory or a failing latest report blocks the archive; proceed only if the user explicitly says "archive anyway", and record that override in the stamped frontmatter (`qa_override: true`).

A merged PR or release tag is **not** a precondition. If the user passes `--release`, or a merged PR/tag is already known, stamp it as metadata — but never block the archive waiting for one.

If any check fails, stop and report exactly which — the spec stays active.

## Steps

1. **Stamp** `_prd.md` frontmatter:
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@skills/archive-spec/SKILL.md` around lines 24 - 38, The archive stamping flow
in SKILL.md is missing persistence for the QA-bypass path, so update the
frontmatter stamp step to include qa_override: true whenever the Preconditions
allow archive anyway after a failing or missing qa/ check. Keep the logic
aligned with the existing archive rules in the Preconditions and Steps sections
so _prd.md records the override alongside status, archived, and release
metadata.
```

</details>

<!-- fingerprinting:phantom:triton:quartz -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:83fd11575494694996ef3227 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `valid`
- Notes: The archive preconditions mention `qa_override: true`, but the stamp template did not include the field, so a literal archive flow could omit the recorded override.

## Resolution

- Added `qa_override: true` to the archive stamp template, scoped to archive operations that proceed despite failed or missing QA.
- Verification: `rtk make verify` — passed. Go tests passed (`821` tests across `18` packages), `roundfix skills check` passed, and the CLI build succeeded.
