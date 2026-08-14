---
source: coderabbit
pr: "18"
round: 1
round_created_at: "2026-07-07T00:22:01Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/docs-and-skill-bundle
head_sha: 83c26c480d5fa9c22da307daeaa48d6acbe59c85
file: skills/write-prd/references/prd-template.md
line: 3
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6OwE1X,comment:PRRC_kwDOS0qyts7SkHEE
review_hash: bb9e580c65a32c2628915408a157f6b8c0b65627f46aad90339a43bdf9ea579c
duplicate_of: ""
source_review_id: "4640508670"
source_review_submitted_at: "2026-07-06T23:41:28Z"
---

# Issue 004: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Use the numbered spec slug here too.**

The path placeholder should match the numbered folder convention used everywhere else in this PR. Right now this still says `<feature-slug>`, which can lead authors to create an unnumbered path even though `spec` is numbered. Update it to `<NNNN-feature-slug>`.





<details>
<summary>📦 Proposed fix</summary>

```diff
-Write `docs/specs/<feature-slug>/_prd.md` with this exact structure.
+Write `docs/specs/<NNNN-feature-slug>/_prd.md` with this exact structure.
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
Write `docs/specs/<NNNN-feature-slug>/_prd.md` with this exact structure. Guidance appears as `<!-- comments -->`; delete every comment from the final file. Omit a section only when it is genuinely empty — write "None." rather than deleting Non-Goals or Open Questions, so readers know it was considered.
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@skills/write-prd/references/prd-template.md` at line 3, The PRD template
still uses the unnumbered placeholder <feature-slug>, which conflicts with the
numbered spec folder convention. Update the path example in prd-template.md to
use <NNNN-feature-slug> so authors create docs/specs entries that match the
numbered spec structure. Keep the rest of the template guidance intact and
ensure the placeholder appears consistently wherever the path is referenced.
```

</details>

<!-- fingerprinting:phantom:triton:quartz -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:99b0d77e531e21c4994fc3c4 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `valid`
- Notes: The PRD template header still used the unnumbered `<feature-slug>` path while the template frontmatter uses the numbered spec slug convention.

## Resolution

- Updated the PRD template path to `docs/specs/<NNNN-feature-slug>/_prd.md`.
- Verification: `rtk make verify` — passed. Go tests passed (`821` tests across `18` packages), `roundfix skills check` passed, and the CLI build succeeded.
