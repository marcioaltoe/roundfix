---
source: coderabbit
pr: "129"
round: 1
round_created_at: "2026-08-06T03:33:20Z"
status: invalid
terminal_reason: "The originating Spec requires the Backlog Entry glossary definition to include its lifecycle, type vocabulary, evidence boundary, and write-idea distinction."
head_repository: marcioaltoe/roundfix
head_branch: ma/0075-typed-docs-backlog
head_sha: 04b156c2a36969a67c06958bbc366fc47a6db816
file: CONTEXT.md
line: 46
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6W1gkH,comment:PRRC_kwDOS0qyts7eCVqA
review_hash: 0719ab69cd3d79e86e00d97bcd337b689fe33f936e27555db0568dc64b75f750
duplicate_of: ""
source_review_id: "4870101613"
source_review_submitted_at: "2026-08-06T00:57:04Z"
---

# Issue 001: _ Maintainability & Code Quality_ _ Major_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟠 Major_ | _⚡ Quick win_

**Keep `CONTEXT.md` limited to terminology.**

Line 44 includes lifecycle states, type rules, promotion routing, and the `write-idea` boundary. These are workflow rules, not glossary content. The same contract is already defined in `internal/baseline/assets/modules/context-workflow.json`.

Keep a short term definition in `CONTEXT.md`. Keep lifecycle and routing rules in the workflow documentation.

As per coding guidelines, `**/CONTEXT.md` is a domain glossary only and must exclude implementation details, specifications, and implementation decisions.




<details>
<summary>Suggested scope</summary>

```diff
 **Backlog Entry**:
-Typed intent for what to do next, with a lifecycle of `open`, `promoted` to a named Spec, or `declined` with a reason; it is never evidence, just as a finding is never a commitment. Its types are the Conventional Commits intent vocabulary — `feat`, `fix`, `perf`, `refactor` — so one word carries intent from entry to Spec to commit; a `feat` entry is upstream raw material, never the `write-idea` artifact.
+A typed record of proposed work, distinct from evidence-backed findings.
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
**Backlog Entry**:
A typed record of proposed work, distinct from evidence-backed findings.
_Avoid_: Finding, idea artifact, untyped suggestion
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@CONTEXT.md` around lines 43 - 46, Reduce the “Typed intent” entry in
CONTEXT.md to a concise terminology definition only. Remove lifecycle states,
type vocabulary, promotion routing, and the write-idea boundary from this
glossary entry; retain those rules in
internal/baseline/assets/modules/context-workflow.json and keep CONTEXT.md free
of workflow or implementation details.
```

</details>

<!-- fingerprinting:phantom:poseidon:caracal -->

<!-- cr-indicator-types:refactor_suggestion -->

<!-- cr-comment:v1:5a84493eef90d9fd555e67e7 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: The finding conflicts with the accepted originating contract. Spec
  0075 Task 04 Requirements 1-4 require this glossary entry to name the
  lifecycle, the four Conventional Commits type tokens, the finding/evidence
  boundary, and the `write-idea` distinction. The PRD and TechSpec also require
  those concepts in `CONTEXT.md`, while the full operational templates and
  promotion procedure remain owned by `docs/agents/docs-layout.md` and the
  Baseline module. No production or documentation change is warranted.
