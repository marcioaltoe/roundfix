---
source: coderabbit
pr: "27"
round: 2
round_created_at: "2026-07-15T16:07:28Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/dogfood-findings-remediation
head_sha: 81ff89aabce7ee1f748504ad6e4bcf1ac52ea200
file: Makefile
line: 21
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6RKhzt,comment:PRRC_kwDOS0qyts7V6fip
review_hash: 0fdaf942aeb6d856e3df18ee4cb7bc55914eb6442acd0584e84c4264f20b89f5
duplicate_of: ""
source_review_id: "4705960653"
source_review_submitted_at: "2026-07-15T16:06:26Z"
---

# Issue 025: _ Functional Correctness_ _ Minor_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**Gate `-dirty` on valid Git metadata and include untracked files.**

Line [18] reports `-dirty` when Git is unavailable or `HEAD` does not exist, and it misses untracked files. This makes `--version` build identity inaccurate for source archives and local worktrees.






<details>
<summary>Proposed fix</summary>

```diff
-BUILD_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null)$(shell git diff --quiet HEAD 2>/dev/null || echo -dirty)
+BUILD_COMMIT := $(shell commit=$$(git rev-parse --short HEAD 2>/dev/null) || exit 0; dirty=$$(git status --porcelain --untracked-files=all 2>/dev/null); if test -n "$$dirty"; then dirty=-dirty; else dirty=; fi; printf '%s%s' "$$commit" "$$dirty")
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
# Local build identity for `roundfix --version`: short commit (plus -dirty
# when the tree has changes) and local build time. The release workflow
# stamps only app.Version from the tag and leaves these empty.
BUILD_COMMIT := $(shell commit=$$(git rev-parse --short HEAD 2>/dev/null) || exit 0; dirty=$$(git status --porcelain --untracked-files=all 2>/dev/null); if test -n "$$dirty"; then dirty=-dirty; else dirty=; fi; printf '%s%s' "$$commit" "$$dirty")
BUILD_TIME := $(shell date '+%Y-%m-%d %H:%M:%S %z')
STAMP_LDFLAGS := -X 'roundfix/internal/app.BuildCommit=$(BUILD_COMMIT)' -X 'roundfix/internal/app.BuildTime=$(BUILD_TIME)'
BUILD_FLAGS ?= -buildvcs=false -ldflags "$(STAMP_LDFLAGS)"
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@Makefile` around lines 15 - 21, Update the BUILD_COMMIT dirty-state logic in
the Makefile so “-dirty” is added only when valid Git metadata is available and
the working tree contains tracked or untracked changes. Keep BUILD_COMMIT empty
for source archives or repositories without a valid HEAD, while preserving the
short commit hash for valid clean worktrees.
```

</details>

<!-- fingerprinting:phantom:poseidon:luna -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:ec6501cba1a106e0434d57ad -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Updated `BUILD_COMMIT` to require a valid short `HEAD` before appending `-dirty` and to detect tracked or untracked changes with `git status --porcelain --untracked-files=all`; focused Makefile checks and `make verify` passed.
