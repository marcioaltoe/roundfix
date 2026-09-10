# Public Implement dispatch and settlement

The assembled binary ran the real Implement Command in disposable Git
repository `/private/tmp/roundfix-0131-qa.4U8iM5/repo` on branch
`test/failed-gate-recovery`. Its committed graph began with `task_01` and
`task_02` completed, corrective `task_04` pending, and terminal `task_03`
failed above `task_04`.

The command used the controlled ACP adapter stored under `runtime/` so the
public CLI, Run Worktree, scheduler, Verification, commit, integration, QA and
Run Event paths remained real while Agent output stayed deterministic:

```text
HOME=/private/tmp/roundfix-0131-qa.4U8iM5/home \
PATH=/private/tmp/roundfix-0131-qa.4U8iM5/fake-bin:/opt/homebrew/bin:/usr/bin:/bin \
rtk <audited-bin>/roundfix implement --spec failed-recovery \
  --agent-command 'codex-acp --stdio' --no-input --no-agent-console
```

The command exited 0 and reported:

```text
Task: task_04 (Batch 001) Repair the finding
Verification passed (attempt 1).
Task commit created: feat: repair the finding
Task task_04 completed.
QA step (Batch 002) for Spec failed-recovery
QA verdict: pass
QA Report commit created: docs: qa report for failed-recovery (pass)
Implement Run run_20260910T122715Z_cb843bfda9ea4002 reached Clean.
task_04 completed — Repair the finding
task_03 completed — Run the final QA gate
qa pass — docs/specs/failed-recovery/qa/qa-report-2026-09-10.md
Clean: all 4 Task(s) completed.
```

The read-only Run Event Stream independently recorded `task_04` as
`in_progress`, its command and verdict as passed, `task_04` settled
`completed`, `task_03` settled `completed`, and the Run outcome `Clean`.

A fresh durable-state read found `status: completed` in both `task_04.md` and
`task_03.md`, plus `verdict: pass` in the seeded QA report. Git history showed
the corrective Task commit before the QA report commit, `git status` was clean,
and a new strict checker process exited 0 with no findings. These rereads prove
dispatch, settlement, ordering and persistence beyond the original command's
optimistic output.
