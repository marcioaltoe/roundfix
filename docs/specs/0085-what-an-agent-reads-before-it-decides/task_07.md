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

- `grep -q 'archive root: ' .coderabbit.yaml` — expected: exits 0. The configuration carries no line naming the root it excludes today, so this fails before the work.

A guard that no per-tree filter was introduced is not declared either: there are
zero today, so it passes before any work. Requirement 2 carries it.

The two commands this Task first declared — that `_archived` appears and that
per-tree filters number at most one — both passed before any work. The
configuration already excludes archives through the `!**/_archived/**` glob and
carries zero per-tree filters, work that landed after this Spec was authored. So
the Task's premise moved: what it must now produce is a configuration whose
exclusion is tied to the root `internal/spec.ArchiveDir` resolves, stated in a
line a reader and this check can both find. If Task 04 leaves the root where the
glob already matches, say so in the Result and record why nothing else changed —
the marker line is still the deliverable, because a configuration and a resolver
that must agree with nothing linking them is the defect this Spec exists to end.

## Context

- instruction: `docs/workflow/authorizations/2026-08-09-what-an-agent-reads-before-it-decides.md`

## References

- `_prd.md` → the archive read path.
- `_techspec.md` → Build Order 7.
