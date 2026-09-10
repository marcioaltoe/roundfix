# Public checker flows

The public entry point was `roundfix spec check <slug> --strict` against the
filesystem-backed fixture at `fixtures/failed-recovery/`. Every invocation ran
in a fresh process. The fixture adds pending corrective `task_04` beneath QA
`task_03`, whose recorded status is `failed`.

The reference binary was rebuilt from pre-change commit `62bb7874`:

```text
roundfix 0.12.0 (62bb7874, built 2026-09-10 09:24:01 -0300)
```

On the failed-gate fixture, that reference binary exited 2:

```text
validate qa gate: QA gate result is invalidated for Task "task_03" because
these dependencies are not completed: task_04
```

The assembled `2d9ab9b2` binary exited 0 on the same bytes:

```text
Spec failed-recovery
No findings. Authored Verification commands were not executed.
```

A fresh `git status --short` in the fixture repository emitted no paths after
the successful check. This confirms the load did not write Task status or any
other fixture byte.

For the completed-gate probe, only `task_03` changed from `status: failed` to
`status: completed`. Both reference and assembled binaries exited 2 with the
same diagnostic and named dependency:

```text
validate qa gate: QA gate result is invalidated for Task "task_03" because
these dependencies are not completed: task_04
```

This is the public representation of the unchanged `StaleGateError`. A fresh
Git comparison also found no change to `internal/spec/errors.go` across
`62bb7874..2d9ab9b2`, so the exported error definition and its public rendering
remain the same.

Fail-closed probes through the same CLI produced these exit-2 diagnostics:

- missing status: `unsupported status "" (allowed: pending, in_progress, completed, failed)`;
- unsupported status: `unsupported status "suspended" (allowed: pending, in_progress, completed, failed)`;
- malformed status: `parse frontmatter: yaml: line 2: did not find expected ',' or ']'`.

The uncovered-leaf fixture removed `task_04` from the QA Task's dependencies.
Both the failed and completed QA verdicts exited 2 with:

```text
QA Task "task_03" does not depend on every leaf; uncovered Tasks: task_04
```

These probes cover invalid input, out-of-order recovery, repeated fresh loads,
and the exact positive boundary: only a recorded failed gate gains acceptance.
