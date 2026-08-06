---
task: task_03
spec: replay-0060-task-03
status: pending
type: docs
complexity: medium
---

# Task 03: Rehearse the lifecycle end to end and prove the gates

## Requirements

1. MUST rehearse on a scratch branch with throwaway documents — one finding and
   one inbox note — and MUST NOT commit the rehearsal fixtures.
2. MUST prove the move preserves history: `git log --follow` on the moved file
   reaches its pre-adoption commits.
3. MUST prove the index populates with the provenance column holding the
   pre-adoption path and the path column resolving from the index.
4. MUST prove no repository link points at the pre-adoption path afterward,
   counting occurrences before and after rather than sampling one.
5. MUST prove the archive gate passes for the self-contained Spec.
6. MUST prove the archive gate fails, naming the offending link, when a stale
   link is injected.
7. MUST prove the archive gate does not fail on prose that names
   `docs/findings/` without linking into it — the false-positive case.
8. MUST prove the archive gate passes trivially for a Spec carrying no index,
   which is the migration boundary.
9. MUST record every command and its outcome in this Task's `## Result`, and
   MUST leave the repository with the rehearsal fully discarded.

## Acceptance Criteria

- [ ] `git log --follow` on the moved document reaches its pre-adoption history.
- [ ] The occurrence count for the pre-adoption path is zero outside the index's
      provenance column and Git history.
- [ ] The archive gate passes for the self-contained Spec.
- [ ] The archive gate fails on the injected stale link and names it.
- [ ] The archive gate passes on the prose-mention case.
- [ ] The archive gate passes on a Spec with no index.
- [ ] `git status --porcelain` is empty of rehearsal artifacts at the end.

## Verification

- `make verify` — expected: exit 0.
- `git status --porcelain` — expected: no rehearsal fixture remains.

## References

`_prd.md` → Success Metrics; `_techspec.md` → Testing Approach, Risks.
