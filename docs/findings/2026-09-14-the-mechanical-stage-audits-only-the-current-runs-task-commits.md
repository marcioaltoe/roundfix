---
status: pending
created_at: 2026-09-14
updated_at: 2026-09-14
---

# QA gate — The mechanical stage audits only the current Run's Task commits (2026-09-14)

This came up while writing Spec 0138, which narrows the QA gate's matrix. That
Spec meant to have the gate take the changed-path audit from the mechanical
stage instead of rebuilding it. Reading the code showed that the stage audits
only part of what a Spec delivered. Spec 0138 therefore leaves the gate auditing
the commits the stage never received, and routes the gap here.

## 1. A Spec delivered across several Runs leaves its earlier Task commits unaudited by the mechanical stage

- Symptom / evidence:
  - The QA mechanical request collects Task commits by reading
    `git log <Run start head>..HEAD` for the Spec's trailers (`mechanicalTaskCommits`
    in `internal/daemon/task_engine.go`). The same request passes the Run's start
    head as `DeliveryTargetRevision`.
  - Task commits from earlier Runs of the same Spec are already ancestors of that
    head, so the stage never receives them.
  - When no governed commit falls inside the Run, the changed-path detector
    records the skip `Task commits`.
  - Spec 0119's report `qa-report-2026-09-10-01.md` records the result: "The
    gate therefore reconstructed all thirteen non-QA Task commits" and chose the
    merge-base rule itself. That reconstruction found real out-of-grant changes
    (QA-05), so the gate's manual audit is load-bearing today, not redundant.
- Root cause:
  - The commit range starts at the Run start head rather than at the delivery
    target, and only the newest commit per Task is kept.
  - The authorizing revision is the merge base of that head and the commit's
    parent. For a commit already on the Spec branch, that merge base can include
    a grant committed on the same branch. That is weaker than the documented rule
    that a grant must land in the delivery target's ancestry before the consuming
    delivery.
- Action / suggestion:
  - Route to a Spec that decides which revision authorizes a Task commit from an
    earlier Run and which delivery target the mechanical stage reads. That is an
    authority decision, so Spec 0138 does not make it.
  - Candidate owners are Spec 0125 (repository identity and Run Branch policy)
    and Spec 0129 (gate recovery).
  - Until then, the qa-gate skill keeps the gate auditing, by command, the Task
    commits the stage did not receive.
