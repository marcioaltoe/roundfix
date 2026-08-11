---
source: coderabbit
pr: "153"
round: 2
round_created_at: "2026-08-10T22:09:34Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/a-proof-that-can-refuse
head_sha: d68e5a2a65875cf1f7a5e9976514c7be60ee5d5d
file: internal/speccheck/verification.go
line: 29
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6YDU5l,comment:PRRC_kwDOS0qyts7fvgDG
review_hash: 75f57b8d0f4e5ab79ddabdd2a190ad430de029208e7cac08b856519cf09abd43
duplicate_of: ""
source_review_id: "4901334183"
source_review_submitted_at: "2026-08-10T22:08:33Z"
---

# Issue 001: _ Functional Correctness_ _ Major_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Anchor each success form to the terminal command.**

Line 29 leaves the `exit 0`, `test -z`, `[-z`, and equality alternatives unanchored. `MatchString` can therefore match text inside a quoted argument. For example, `git diff --name-only HEAD | grep -q 'exit 0'` is classified as vacuous, although `grep` exits 1 when the unchanged tree produces no input. Add this case as a non-vacuous regression and anchor or parse the success forms. `internal/speccheck/citations.go`, Lines 1288-1298, promotes every match to `SeverityError`, so this false positive can block a valid Task.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/speccheck/verification.go` at line 29, Anchor the exit 0, test -z,
[-z, and variable-equality alternatives in emptyOutputSucceedsPattern to the
terminal command rather than allowing matches inside quoted arguments. Add a
regression case for git diff --name-only HEAD | grep -q 'exit 0' and ensure it
is classified as non-vacuous, while preserving recognition of genuine terminal
success commands.
```

</details>

<!-- fingerprinting:phantom:poseidon:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:38976336f2fc413b85569ebe -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `resolved`
- Notes: An unanchored `\bexit\s+0\b` allowed a quote `<success form>` sitting
  inside a consumer command's argument to be misread as a terminal predicate.
  `internal/speccheck/verification.go` now anchors every success form to the
  whole terminal command: `emptyOutputSucceedsPattern` wraps each present form
  in `^(?:rtk\s+)?...`, and a new `emptyOutputFailsPattern` classifies a
  grep terminal as failing on empty input regardless of its arguments. The old
  `terminalSegment`/`commandSeparatorPattern` are gone, replaced by a
  channel-aware evaluator. Focused evidence: `go test ./internal/speccheck/`
  passes; the new regression "a quoted success form inside grep still fails
  when empty" (`git diff --name-only HEAD | grep -q 'exit 0'`) asserts
  non-vacuous, and "a quoted empty-test inside grep still fails when empty"
  pins the `test -z` variant. Daemon owns authoritative `make verify`.
