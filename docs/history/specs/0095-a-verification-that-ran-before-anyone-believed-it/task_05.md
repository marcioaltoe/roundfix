---
task: task_05
spec: 0095-a-verification-that-ran-before-anyone-believed-it
status: completed # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: medium
---

# Task 05: Refuse a Verification that reaches outside the repository

## Overview

A Verification that depends on state the repository does not carry produces a
false red: the work is right and the gate fails. Three forms were measured in one
fleet session, and all three cost Tasks that had finished correctly — an
undeclared environment variable, a `test -n "$VAR" &&` guard in front, and a
dependency on a temporary directory or tree snapshot outside the repository.

## Requirements

1. MUST refuse an authored Verification command that references an environment
   variable the repository does not declare.
2. MUST refuse a command that guards itself on an environment variable's
   presence, since an absent variable fails the whole chain rather than the work.
3. MUST refuse a command that depends on a path outside the repository, unless
   the Task itself creates that path.
4. MUST NOT refuse a command that writes to and reads from a temporary path it
   creates within the same command, which is the working redirect form the
   authoring contract teaches.
5. MUST skip rather than fail when the artifact it reads is absent.
6. MUST register at the tasks stage.

## Subtasks

- [ ] Add the refusal code and its detector.
- [ ] Cover each measured non-hermetic form and the working redirect form.
- [ ] Register it at the tasks stage.

## Acceptance Criteria

- [ ] Each of the three measured forms is refused and named.
- [ ] The working redirect form the authoring contract teaches is not refused.
- [ ] A Spec without a task graph skips rather than fails.
- [ ] The code appears in the staged registry at the tasks stage.

## Verification

- `grep -q 'SC-VERIFY-NON-HERMETIC' internal/speccheck/*.go && grep -q 'CodeVerifyNonHermetic' internal/speccheck/coherence.go` — expected: exits 0, proving the code exists and is registered. Fails today on both clauses.
- `go test -count=1 ./internal/speccheck -run 'TestVerifyNonHermetic' -v > /tmp/0095-t05.log 2>&1; s=$?; grep -q '^--- PASS: TestVerifyNonHermetic' /tmp/0095-t05.log || { cat /tmp/0095-t05.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing test; fails today.
- `grep -c 'PASS' /tmp/0095-t05.log > /tmp/0095-t05-n.txt; test "$(cat /tmp/0095-t05-n.txt)" -ge 4 || { echo "expected at least four passing cases, got $(cat /tmp/0095-t05-n.txt)"; cat /tmp/0095-t05.log; exit 1; }` — expected: exits 0, proving the three refused forms and the permitted redirect form each run as their own case.

## References

`_techspec.md` → Build Order 5; Testing Approach, the static detectors.
`_prd.md` → Core Feature 3; User Story 4. ADR-0093, ADR-0094.

## Result

Implemented `SC-VERIFY-NON-HERMETIC` as a tasks-stage Spec Consistency Check
detector. It names undeclared environment-variable references, prioritizes the
measured `test -n "$VAR" &&` presence-guard diagnosis, and refuses reads from
paths outside the repository unless an earlier command-local action or Task
Verification command creates the path. Shell assignments, `read`, and `for`
variables remain command-local rather than being mistaken for environment
dependencies. The stable finding code is documented as Non-Hermetic
Verification in `CONTEXT.md`.

Focused checks:

- Before implementation,
  `rtk env GOCACHE=/tmp/roundfix-0095-task05-go-cache go test ./internal/speccheck -run '^TestVerifyNonHermetic$/undeclared_environment_variable$' -count=1`
  exited 1 because `CodeVerifyNonHermetic` and
  `NonHermeticVerification` did not exist.
- After the final edits,
  `rtk env GOCACHE=/tmp/roundfix-0095-task05-go-cache go test ./internal/speccheck -run '^TestVerifyNonHermetic' -count=1`
  exited 0: `ok roundfix/internal/speccheck 0.419s`.
- `rtk env GOCACHE=/tmp/roundfix-0095-task05-go-cache go test ./internal/speccheck -count=1`
  exited 0: `ok roundfix/internal/speccheck 0.879s`.
- `rtk env GOCACHE=/tmp/roundfix-0095-task05-go-cache go run -buildvcs=false ./cmd/roundfix spec check 0095-a-verification-that-ran-before-anyone-believed-it`
  exited 0 with `No findings.` The only recorded skip was the existing absent
  `references/_index.md`; the new detector did not reject this Spec's valid
  temporary-output chains.
- `rtk git diff --check` exited 0 with no diagnostics.

Acceptance evidence:

- `TestVerifyNonHermetic` has independent refusal cases for an undeclared
  environment variable, the environment-presence guard, and an external tree
  snapshot. Each case asserts `CodeVerifyNonHermetic`, error severity, the
  declaring Task location, and the named matched form.
- The same test permits the working redirect that writes and then reads its
  log in one command. Companion cases permit an output created in an earlier
  Verification command, an explicit temporary-directory creation, and a
  binary created with `-o`; a read before a later write remains refused.
- `TestVerifyNonHermeticSkipsWithoutTaskGraph` checks that the existing
  `no-taskgraph` fixture produces no finding and records the detector as skipped
  for missing `_tasks.md`.
- `TestVerifyNonHermeticRegistersAtTasksStage` loads a complete temporary Task
  Graph through `CheckStage(StageTasks)` and observes the named refusal. The
  stage-scope tests also require the code to stay absent from earlier-stage
  findings and present in their named skip sets.

The Daemon-owned commands under `## Verification` were not run in this Agent
turn.
