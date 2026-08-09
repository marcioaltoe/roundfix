---
task: task_07
spec: 0085-what-an-agent-reads-before-it-decides
status: pending
type: infra
complexity: low
---

# Task 07: Exclude the single archive root from review

## Overview

With every retired artifact under one root, the review configuration stops
carrying a filter per archived tree and excludes the one root instead. Authorized
tooling work with exactly one bounded file.

## Requirements

1. MUST exclude the single archive root from review.
2. MUST remove the per-tree filters the single root replaces.
3. MUST leave every non-archive path filter unchanged.
4. MUST NOT touch any path outside the bounded list below.

## Subtasks

- [ ] Exclude the single root.
- [ ] Remove the superseded per-tree filters.

## Acceptance Criteria

- [ ] The archive root is excluded.
- [ ] No per-tree archive filter remains.
- [ ] Non-archive filters are unchanged.

## Bounded scope

Authorized by
`docs/workflow/authorizations/2026-08-09-what-an-agent-reads-before-it-decides.md`,
which bounds `.coderabbit.yaml` to the path filter that excludes the archive.
This Task may create or modify only:

- `.coderabbit.yaml`
- `docs/specs/0085-what-an-agent-reads-before-it-decides/task_07.md`

Any other path is out of scope; stop and fail the Task rather than widen it.

## Verification

- `grep -q '_archived' .coderabbit.yaml` — expected: exits 0, proving the archive root is named in the configuration.
- `test "$(grep -c 'docs/specs/_archived\|docs/findings/_archived' .coderabbit.yaml)" -le 1` — expected: exits 0, proving the per-tree filters were replaced by one entry rather than added to.
- `test -z "$(git diff --name-only -- . ':!.coderabbit.yaml' ':!docs/specs/0085-what-an-agent-reads-before-it-decides/task_07.md')"` — expected: exits 0, proving no path outside the bounded list moved.

## Context

- instruction: `docs/workflow/authorizations/2026-08-09-what-an-agent-reads-before-it-decides.md`

## References

- `_prd.md` → the archive read path.
- `_techspec.md` → Build Order 7.
