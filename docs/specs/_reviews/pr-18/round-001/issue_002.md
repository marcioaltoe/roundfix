---
source: coderabbit
pr: "18"
round: 1
round_created_at: "2026-07-07T00:22:01Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/docs-and-skill-bundle
head_sha: 83c26c480d5fa9c22da307daeaa48d6acbe59c85
file: skills/setup-workflow/domain.md
line: 24
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6OwE1Q,comment:PRRC_kwDOS0qyts7SkHD9
review_hash: 8c97ff9a02d6d2f741ff7edbb07d3ae2f36f469c55ced0bb46131107488f848c
duplicate_of: ""
source_review_id: "4640508670"
source_review_submitted_at: "2026-07-06T23:41:28Z"
---

# Issue 002: _ Maintainability & Code Quality_ _ Minor_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟡 Minor_ | _⚡ Quick win_

**Tag the fenced examples.**

Both bare fences are tripping MD040. Add a language such as `text` so markdownlint stays green.

<details>
<summary>🛠️ Proposed fix</summary>

```diff
-```
+```text
 ...
-```
+```text
```
</details>







Also applies to: 28-39

<details>
<summary>🧰 Tools</summary>

<details>
<summary>🪛 markdownlint-cli2 (0.22.1)</summary>

[warning] 17-17: Fenced code blocks should have a language specified

(MD040, fenced-code-language)

</details>

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@skills/setup-workflow/domain.md` around lines 17 - 24, The fenced examples in
domain.md are missing language tags, causing MD040 failures. Update the bare
fenced blocks in the affected examples to use a language like text so
markdownlint passes, and make the same change in both occurrences referenced by
the diff. Use the existing fenced example sections in the document as the target
for this cleanup.
```

</details>

<!-- fingerprinting:phantom:triton:quartz -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:cd73d29fc6093821070dfe30 -->

_Source: Linters/SAST tools_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `valid`
- Notes: The file structure examples used bare opening fences, which violates MD040.

## Resolution

- Tagged both file-structure example fences as `text`.
- Verification: `rtk make verify` — passed. Go tests passed (`821` tests across `18` packages), `roundfix skills check` passed, and the CLI build succeeded.
