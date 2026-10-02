---
type: fix
status: promoted
created: 2026-10-01
spec: 0212-gates-and-run-storage-that-let-a-correct-delivery-finish
reason: null
---

# Task Context has no kind for a deleted path

## Opportunity

A Task's `## Context` accepts only `instruction`, `interface` and `creates`. Spec 0200's task_04 deletes `.agents/skills/context7/SKILL.md` and declared it as `interface`. Once the Task ran, the strict check refused the Spec with `SC-REF-UNRESOLVED`, and the QA gate stopped at its precondition (Run `run_20261001T013631Z_c59b100cbb9f7af1`). The operator relabelled the entry as `creates`, which the detector does not resolve. That passes, but it says the opposite of what the Task does.

## Value

A Task that removes a file can declare it truthfully, the scope audit sees the deletion as declared, and the QA precondition stops refusing a correct Task.

## Shape

Add a `deletes` Context kind. It is skipped by the unresolved-path detector, refused when the path still exists after the Task's commit, and counted as declared by the scope audit. Teach write-tasks to use it.
