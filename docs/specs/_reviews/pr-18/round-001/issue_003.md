---
source: coderabbit
pr: "18"
round: 1
round_created_at: "2026-07-07T00:22:01Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/docs-and-skill-bundle
head_sha: 83c26c480d5fa9c22da307daeaa48d6acbe59c85
file: skills/setup-workflow/spec-routing.md
line: 9
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6OwE1T,comment:PRRC_kwDOS0qyts7SkHEB
review_hash: 9014b71976d8366f42f75569664cdc84d26fb03d49928f9673e30f19c9cda669
duplicate_of: ""
source_review_id: "4640508670"
source_review_submitted_at: "2026-07-06T23:41:28Z"
---

# Issue 003: _ Maintainability & Code Quality_ _ Minor_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟡 Minor_ | _⚡ Quick win_

**Tag the pipeline fence.**

This bare fence is tripping MD040. Add a language such as `text` so markdownlint passes.

<details>
<summary>🛠️ Proposed fix</summary>

```diff
-```
+```text
 write-idea → write-prd → write-techspec → write-tasks → implement-spec / implement-task → qa-gate → archive-spec
 ```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion

```

</details>

<!-- suggestion_end -->

<details>
<summary>🧰 Tools</summary>

<details>
<summary>🪛 markdownlint-cli2 (0.22.1)</summary>

[warning] 7-7: Fenced code blocks should have a language specified

(MD040, fenced-code-language)

</details>

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@skills/setup-workflow/spec-routing.md` around lines 7 - 9, The spec-routing
pipeline fence is missing a language tag, causing the markdownlint MD040
failure. Update the fenced block in the markdown content to use a tagged fence
such as text while keeping the existing pipeline sequence unchanged. Locate the
bare fenced line in the spec-routing document and add the language identifier to
the opening fence only.
```

</details>

<!-- fingerprinting:phantom:triton:quartz -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:bbe07d2595c154d90a98aee0 -->

_Source: Linters/SAST tools_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `valid`
- Notes: The pipeline example used a bare opening fence, which violates MD040.

## Resolution

- Tagged the pipeline example fence as `text`.
- Verification: `rtk make verify` — passed. Go tests passed (`821` tests across `18` packages), `roundfix skills check` passed, and the CLI build succeeded.
