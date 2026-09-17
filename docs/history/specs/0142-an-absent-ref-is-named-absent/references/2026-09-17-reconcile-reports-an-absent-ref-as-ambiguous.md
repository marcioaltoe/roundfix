---
status: done
created_at: 2026-09-17
updated_at: 2026-09-17
spec: docs/specs/0142-an-absent-ref-is-named-absent/_prd.md
---

# Reconciliation — An absent ref is reported as ambiguous (2026-09-17)

Reconciling this machine's Run debris produced fourteen refusals naming an
ambiguity that does not exist.

## 1. A ref that matches nothing is reported as matching too much

- Symptom / evidence:
  - `roundfix reconcile` reported, for fourteen candidates:
    `Run Branch set cannot be proven: classify Run Branch set: resolve target
    branch "ma/0118-a-task-proved-once-does-not-run-twice": resolve local branch
    "ma/0118-a-task-proved-once-does-not-run-twice": short ref is ambiguous`.
  - `git for-each-ref` lists no ref whose name contains that branch, so the
    short name matches nothing at all.
  - The summary read `total=15 released=14 unintegrated=1 preserved=1`, and the
    supported cleanup path was unavailable for every refused candidate.
- Root cause:
  - The resolver asks how many candidate refs share the short name and refuses
    when the count is not one, so a count of zero produces the message written
    for a count above one.
  - The validity check that runs first asks only whether the name is well
    formed, and the existence check that would tell absence from ambiguity is
    defined beside the resolver and not used on this path.
  - The failure is raised while classifying the Run Branch set, so one absent
    target branch refuses every candidate in that set.
- Action / suggestion:
  - Adopted by Spec 0142, which separates the two answers and preserves an
    absent-target candidate with a reason that names the missing ref.
  - Proving delivered content when the target branch is absent stays with Spec
    0125.
