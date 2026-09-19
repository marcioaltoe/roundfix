---
task: task_05
spec: 0150-a-reopen-that-cannot-be-raced
status: completed
type: backend
complexity: medium
---

# Task 05: One resolution, not two

## Overview

Pre-PR review found that Task 02's symlink repair introduced the same class of
defect this Spec exists to close.

`ensureReopenTaskInsideSpecsRoot` resolves the Task path with `EvalSymlinks` and
proves the result is inside the configured Spec Root. `replaceTaskFile` then
calls `EvalSymlinks` on the *unresolved* path again and renames onto whatever
that second call returns. The two resolutions are independent: a symlink swapped
between them sends the write to a target that was never validated, outside the
Spec Root.

Task 01 closed a decide-then-write window. This closes a validate-then-resolve
window opened by the fix beside it.

## Requirements

1. MUST resolve the Task path exactly once, and MUST write to the path that was
   validated.
2. MUST carry the validated target from the confinement check through to the
   rename rather than re-deriving it.
3. MUST keep a symlinked Task path a symlink, with the target rewritten.
4. MUST keep the write atomic.
5. MUST keep both confinement checks refusing an escaping path.
6. MUST keep every behavior Tasks 01 through 03 delivered.

## Subtasks

- [ ] Thread the validated target through to the replacement.
- [ ] Stop resolving inside the replacement helper.
- [ ] Add a test that the replacement writes to the validated target.

## Acceptance Criteria

- [ ] The replacement never calls `EvalSymlinks` on its own; it writes to the
      path it is given.
- [ ] A symlinked Task path is still a symlink afterwards, with its target
      rewritten.
- [ ] The escaping-path refusals still refuse.
- [ ] The reopen tests from Specs 0149 and 0150 pass unchanged.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/spec/task.go`
- interface: `internal/cli/reopen.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestReopenGateWritesToTheValidatedTarget" ./internal/spec 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestReopenGateWritesToTheValidatedTarget"` — expected: exit 0; before this Task the case does not exist, so the command fails.
- `test "$(grep -c "EvalSymlinks" internal/spec/task.go)" = "0"` — expected: exit 0; the replacement no longer resolves paths itself. Before this Task it calls EvalSymlinks, so the command fails.

## References

- [_techspec.md](_techspec.md) — The symlinked path

## Result

### Implementation

- `ensureReopenTaskInsideSpecsRoot` now returns the resolved Task target after
  both confinement checks accept it, and `deriveReopenPlan` carries that exact
  path into `spec.ReopenGate`.
- `replaceTaskFile` now creates its temporary file beside the path it receives
  and atomically renames onto that path without resolving it again.
- The regression coverage redirects the Task symlink after validation and
  proves that reopen changes only the validated target. CLI coverage proves the
  validated target is stored in the plan and an ordinary symlinked Task path
  remains a symlink while its target is rewritten.

### Focused-check evidence

- Pre-change signal: `GOCACHE=/private/tmp/roundfix-task05-gocache go test
  -count=1 -run '^TestDeriveReopenPlanCarriesTheValidatedTaskTarget$'
  ./internal/cli` failed because the plan contained the unresolved symlink path
  instead of its validated target.
- Validated target and prior reopen behavior: `GOCACHE=/private/tmp/roundfix-task05-gocache
  go test -count=1 -run 'Reopen|DeriveReopenPlan' ./internal/spec
  ./internal/cli` passed.
- Replacement resolution: `rg -n 'EvalSymlinks' internal/spec/task.go` found no
  matches (exit 1), so the replacement helper does not resolve its input.
- Confinement: the focused reopen run passed
  `TestReopenRefusesQATaskResolvedOutsideSpecRoot` and
  `TestReopenRejectsAManifestPathOutsideTheSpecDirectory`.
- Symlink preservation: the focused reopen run passed
  `TestReopenStaleGatePreservesSymlinkedTaskPath`; the target became pending
  and `Lstat` still reported the Task path as a symlink.
- Repository incremental gate: `GOCACHE=/private/tmp/roundfix-task05-gocache
  rtk make verify-incremental` passed outside the sandbox. Its first sandboxed
  attempt reached the tests but could not read the process table in two force-stop
  integration cases (`operation not permitted`).
- The two commands under `## Verification` were not run; the Daemon owns them.
