# Backend flow evidence

Build: `75161e9c3a5f7554cd1e0b9290bce6c61820b5c7`.

Fresh exact integration command:

```text
rtk env GOCACHE=<repository>/.gocache go test -count=1 -run '<named selection>' ./internal/daemon ./internal/runevent ./internal/spec
```

The selection passed:

- red repository entry starts no Agent Session and records the bounded
  precondition failure;
- green entry orders precondition, Agent work, full post-Agent Verification,
  then commit;
- non-gate Tasks skip the entry precondition;
- a passing entry does not consume the one post-Agent repair;
- owner-, group-, and other-only execute modes are rejected;
- Task, Batch, and QA commits omit executable paths and retain ordinary paths;
- repository-external and symlink-crossing behavior remains unchanged;
- an explicitly selected ignored tracked path still stages;
- ordinary Task, Batch, and QA commit contracts remain stable.

`rtk ... go test -count=1 -run '^TestProjectStreamEvent.*Verification'
./internal/runevent` also passed, independently confirming the public event
projection used by the precondition classification.

The public `roundfix implement` journey remains environment-blocked. The CLI
persists the Run Database and User Config under user-scoped Roundfix Home.
This Daemon-assigned QA session may write only the Run Worktree's Spec QA
directory and may not repurpose `HOME`. A full-access session must run a
disposable red-repository Implement Command, then read the persisted Task
status and `roundfix events <run-id>` projection. This block does not invalidate
the passing current-build integration seams, but it prevents full backend
surface completion.
