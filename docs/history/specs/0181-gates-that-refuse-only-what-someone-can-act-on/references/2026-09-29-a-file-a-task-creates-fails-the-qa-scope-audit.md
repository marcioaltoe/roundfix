---
type: fix
status: promoted
created: 2026-09-29
spec: 0181-gates-that-refuse-only-what-someone-can-act-on
reason: null
---

# A file a Task creates without declaring it fails the QA scope audit

## Symptom

The authored QA gate checks that every changed file stays inside the paths the Tasks declare. A Task that needs a file it did not plan, often a new test file, changes it and settles `completed`. The gate then reports the path as undeclared, and the Spec needs a corrective edit to the Task's `## Context` and a second QA Run. This cost two extra QA Runs across Specs 0179 and 0180. Spec 0180's `task_02` changed `internal/store/delivery_test.go` and its `task_04` changed `internal/cli/deliver_revalidate_test.go`. Spec 0179's Task 06 changed three production lock files.

## Where

- The Daemon's commit preparation, `prepareTaskCommit` in `internal/daemon/task_engine.go`. It stages every changed path and records only a count in the commit event.
- The Task prompt, `internal/agent/spec_prompt.go`, and the `implement-task` skill. Both forbid the Agent from editing anything in its Task file except `## Result`.
- The scope rule in `skills/write-tasks/SKILL.md`: "Every path a Task edits is declared under `interface:` or `creates:`".

## Expected

The Daemon records each path a Task changed that it did not declare. It writes the path into the Task file it already commits, as `creates:` for a new file or `interface:` for an existing one, and lists it in the commit event. The QA gate counts a recorded path as declared. A governed path still requires its `_authorization.md` bound, and the Daemon never adds one.

## Evidence

- `docs/history/specs/0180-a-prepared-queue-that-revalidates-before-each-spec/qa/qa-report-2026-09-28.md`, the scope row.
- `docs/history/specs/0179-review-findings-with-evidence-and-no-unselected-providers/qa/qa-report-2026-09-29-02.md`, the scope row.
- Recorded in `docs/handoffs/2026-09-29-after-v0-20-0.md`, Wave 4.
