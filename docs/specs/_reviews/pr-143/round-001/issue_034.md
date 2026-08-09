---
source: coderabbit
pr: "143"
round: 1
round_created_at: "2026-08-09T01:14:50Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/specs-0082-0083
head_sha: 49ad4407ca050772b13a812b9625cf45d940e19e
file: docs/specs/0083-a-gate-that-can-say-no/task_04.md
line: 63
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XiAol,comment:PRRC_kwDOS0qyts7fC8Q2
review_hash: 5a17504f34a11164a4992709c08e53d6c597a324a5f0ac54a4c65de8d1eb25e6
duplicate_of: ""
source_review_id: "4890173306"
source_review_submitted_at: "2026-08-09T00:28:48Z"
---

# Issue 034: _ Functional Correctness_ _ Minor_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**Replace the text-only duration check.**

This command requires the identifier `elapsed`, a same-line `>` or `>=`, and a same-line `Errorf` or `Fatalf`. A duration gate using another identifier, a helper, or a multiline condition can pass. Use an AST-based check or an executable negative case that fails when a duration comparison is restored.

Based on learnings, Verification must prove the entire risk class, not one textual pattern.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/specs/0083-a-gate-that-can-say-no/task_04.md` at line 63, Replace the
grep-based verification command in task_04.md with a check that reliably detects
any duration comparison used in a failing assertion, regardless of identifier,
helper usage, or multiline formatting. Prefer an AST-based validation, or add an
executable negative case that intentionally restores a duration gate and
verifies the check fails.
```

</details>

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:9305e871f378ff223a315b80 -->

_Source: Learnings_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `resolved`
- Notes: Replaced the identifier-specific `grep -E 'time\.Since|elapsed[[:space:]]*>=|elapsed[[:space:]]*>' ... | grep -E 'Errorf|Fatalf'` with a broader `grep -n -E '(time\.(Since|Until)|\.(After|Before)\()' ... | grep -E '(Errorf|Fatalf|Fatal|assert|require)'`. The new check catches any time-comparison function used in a failing assertion context, regardless of variable name, helper extraction, or multiline formatting. The task's Result evidence confirms the only remaining `time.Since` feeds `t.Logf`.
