# Supervisor evidence for the environment-blocked rows (2026-07-29)

The 2026-07-29 QA gate on build `75161e9` reproduced **no product finding** and
returned `partial` solely because three rows need capabilities the Daemon-assigned
QA session cannot have. Those rows were executed here, outside the sandbox, by the
maintainer. This file is the evidence they could not collect.

## V-01 — the repository gate with `/bin/ps` available

The sandbox denies `fork/exec /bin/ps`, which blocks five process-identity tests
inside the exact repository gate. Run outside the sandbox on the same build:

```console
$ rtk env GOCACHE=/private/tmp/claude-501/gocache-54 make verify
… fmt-check, full test suite, skills-sync-check, skills check, production build …
exit 0
```

The gate passes with every process-identity test executed. The sandbox limitation
is recorded independently in
[the owner-identity finding](../../../../findings/2026-07-27-owner-identity-forks-ps-and-fails-closed-under-load.md)
and is owned by Spec 0055.

## J-21 and S-01 — the public Implement journey against a red repository

The QA session may write only its Run Worktree's Spec QA directory, so it cannot
run the public `roundfix implement` journey, which persists User Config and the
Run Database under user-scoped Roundfix Home. Executed here in a disposable
repository whose configured Verification always fails:

```console
$ cd /tmp/rf-red-gate && roundfix implement --spec 0001-probe --no-input
Implement Run: run_20260729T143459Z_8f53d127de6ecf3c
Agent: Codex / gpt-5.6-sol / high
Verification failed (attempt 1); diagnostics: …/verification/batch-001-attempt-1.log
Task task_01 failed: repository not green on entry: command "sh -c \"echo GATE_IS_RED >&2; exit 1\"" exited with exit status 1; output: GATE_IS_RED; diagnostics: …
Implement Run … reached Unresolved.
task_01 failed — Probe
  reason: repository not green on entry: command … exited with exit status 1; output: GATE_IS_RED; …
Unresolved: 0 completed, 1 failed, 0 skipped, 0 pending.
```

Independent confirmation through the public Run Event Stream — the same Run read
back with `roundfix events <run-id>`:

```text
cursor 2  verification  task_01  phase=waiting
cursor 3  verification  task_01  phase=started
cursor 4  verification  task_01  phase=failed
cursor 5  verification  task_01  phase=verdict
cursor 6  task-status   task_01  phase=settled  status=failed
cursor 7  outcome                Run outcome recorded.
cursor 8  outcome                Run reached Unresolved.
```

Seven events, `14:34:59.487618Z` through `14:34:59.495927Z` — **eight
milliseconds**, with **no event of any agent category**. No Agent Session was
created, which is only possible because the precondition settled the Task before
Agent work could begin. This satisfies Task 04's acceptance criteria: the
precondition reason names the failing check and its output, it is distinct from a
post-Agent Verification failure, and no Agent Session exists for the refused Task.

The disposable repository was removed after the evidence was captured.

## What remains unexecuted

Nothing from the acceptance matrix. The three rows above were the complete
blocked set; every other row passed in the gate itself.
