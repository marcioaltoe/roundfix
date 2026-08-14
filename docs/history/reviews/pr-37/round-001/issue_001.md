---
source: coderabbit
pr: "37"
round: 1
round_created_at: "2026-07-27T01:53:02Z"
status: invalid
head_repository: marcioaltoe/roundfix
head_branch: ma/0036-doctor-skill-readiness
head_sha: 9a6b7f9433b9779afe75f38d833b780ceb2555ed
file: .agents/skills/golang-testing/SKILL.md
line: 88
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6T6oNH,comment:PRRC_kwDOS0qyts7Zy0Nf
review_hash: 26882155352a19baadcc93e050cd5e07df4c8e39682273b026790f6a893d9691
duplicate_of: ""
source_review_id: "4783144632"
source_review_submitted_at: "2026-07-27T01:52:01Z"
---

# Issue 001: _ Maintainability & Code Quality_ _ Minor_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟡 Minor_ | _⚡ Quick win_

**Specify a language for the fenced block.**

The fence at Line 77 has no language identifier, triggering MD040. Use `text` because this block documents filename mappings.

<details>
<summary>Proposed fix</summary>

```diff
-```
+```text
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
Name the test file after the source file it tests, not after the function or method under test. Go's convention is one test file per source file (`foo.go` -> `foo_test.go`), because tools (`go test`, coverage reports, IDE "jump to test" navigation, `gotests`) and reviewers all resolve tests by source file, not by symbol. A source file usually declares several functions/methods; splitting its tests by symbol name scatters them across many files and breaks that file-to-file mapping.

```

</details>

<!-- suggestion_end -->

<details>
<summary>🧰 Tools</summary>

<details>
<summary>🪛 LanguageTool</summary>

[style] ~85-~85: As an alternative to the over-used intensifier ‘very’, consider replacing this phrase.
Context: ...d be helloworld_test.go ```  Exception: very large source files MAY be split into multiple...

(EN_WEAK_ADJECTIVE)

</details>
<details>
<summary>🪛 markdownlint-cli2 (0.23.0)</summary>

[warning] 77-77: Fenced code blocks should have a language specified

(MD040, fenced-code-language)

</details>

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

````
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In @.agents/skills/golang-testing/SKILL.md around lines 75 - 88, Specify the
fenced code block in the documentation around the test-file naming examples with
the text language identifier, changing the opening fence to ```text while
preserving the filename mapping content.
````

</details>

<!-- fingerprinting:phantom:triton:luna -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:fb71716638f8ae28f87aa108 -->

_Source: Linters/SAST tools_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: The missing fence language is present in the fetched upstream
  artifact, but `.agents/skills/golang-testing/SKILL.md` is an externally
  managed Repository Skill Set member recorded in `skills-lock.json`, not a
  repository-owned skill. `docs/agents/skill-dispatch.md` prohibits adapting
  upstream-managed skills locally and requires proposed changes to go to their
  upstream source. Changing this snapshot would also invalidate its recorded
  immutable hash, so this Batch must not apply the suggested local edit.

## Verification

- Inspected the `golang-testing` entry in `skills-lock.json`; it records
  `marcioaltoe/skills` as the source and the current snapshot hash.
- Inspected commit `672eb5c`, which refreshed this file and its lock hash
  together from the external Repository Skill Set source.
