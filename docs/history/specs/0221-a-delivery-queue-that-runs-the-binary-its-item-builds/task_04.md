---
task: task_04
spec: 0221-a-delivery-queue-that-runs-the-binary-its-item-builds
status: completed
type: backend
complexity: medium
---

# Task 04: The queue owner runs an item's steps with the binary the item builds

## Overview

The queue owner starts `implement`, `archive` and `review` in the item
worktree with its own executable, which refuses a Project Config key the item
added (Spec 0215, operator log entry 105 of 2026-10-02). This Task makes the
owner, when its configuration declares `delivery.item_binary`, build that
binary in the item worktree before each step, ask it `migrate --check`, and
run the step with it unless it would change the Run Database schema
(ADR-0225). A repository without the declaration is delivered exactly as
before.

## Requirements

1. MUST add `stepExecutable` and the `executable` parameter of `runRoundfix`
   with the shapes of the TechSpec's Interfaces, and the `log` writer field of
   `commandDeliveryWorkflow`, which writes to standard error when nil.
2. MUST implement "The step executable" in order: read only
   `workflow.loaded.Config.Delivery.ItemBinary`; without a declaration return
   `os.Executable()` and run no Git command, build or probe; with one, run
   `git check-ignore -q` on the declared path in the item worktree, run the
   build through `daemon.ExecVerifier` with the build log at
   `<artifact dir>/delivery/<slug>/item-binary-build.log`, start
   `<item worktree>/<path> migrate --check` with the environment
   `deliveryCommandEnvironment` gives every child, and select the item binary
   on exit `0` and the owner's executable on any other exit.
3. MUST write, through the `log` writer, exactly one line per step: API
   Contract 2's line when the item binary is selected, or API Contract 3's
   notice with the probe's exit code and the first non-empty line of its
   stderr, else of its stdout.
4. MUST return the errors of API Contract 5 — an unignored path, a failed
   build naming the verifier's error and log path, and a binary that cannot
   start — so the engine parks the item `delivery-error`, and MUST start no
   step after any of them.
5. MUST make `RunSpec`, `Archive` and `Review` resolve the executable with
   the step names `implement`, `archive` and `review` and pass it to
   `runRoundfix`, keeping each step's arguments, working directory,
   environment and result handling unchanged.
6. MUST NOT read the item's Project Config or the default branch for the
   declaration, change the repository gate, the checks, the merge, Park Class
   names, the Run Database schema, `.roundfixrc.yml`, or any existing test.
7. MUST add the tests named in Verification to the new file
   `internal/cli/deliver_item_binary_test.go`, over a disposable repository
   whose `.gitignore` ignores `bin/`, a disposable Roundfix Home, and a fake
   item binary compiled once with `testfixture.FixtureBinary` (ADR-0125) that
   records its arguments and working directory to a file the test names and
   answers `migrate --check` with an exit code the test chooses; the declared
   build copies it to `bin/roundfix`. No test writes an executable with a
   literal file mode, and every fixture process ends with the test (ADR-0213).
   The tests MUST cover Testing Approach 3: the item binary selected with API
   Contract 2's line; `RunSpec`, `Archive` and `Review` each starting the
   fixture with `implement --spec <slug>`, `archive <slug>` and `review` in
   the item worktree (a test may let the fixture fail after it records its
   start); a probe exit of `2` selecting the owner's executable with API
   Contract 3's notice while the fixture records only the probe; a build
   exiting `7` and an unignored path each returning API Contract 5's error
   with nothing started; and no declaration selecting the owner's executable
   with no build, no probe and no line.

## Subtasks

- [ ] Thread the executable through `runRoundfix` and its three callers.
- [ ] Build, check the ignore rule and probe the item binary.
- [ ] Write the selection line or the fallback notice.
- [ ] Return the park errors.
- [ ] Add the step executable tests with the compiled fixture.

## Acceptance Criteria

- [ ] With a declaration and a probe exit of `0`, each of the three steps
      starts the built binary in the item worktree with the owner's former
      arguments, and the log holds one API Contract 2 line per step.
- [ ] A probe exit of `2` starts the owner's executable and logs one API
      Contract 3 notice; the fixture never runs a step.
- [ ] A failed build or an unignored path returns API Contract 5's error and
      starts nothing; without a declaration nothing new runs or is logged.

## Context

- interface: `internal/cli/deliver_workflow.go`
- creates: `internal/cli/deliver_item_binary_test.go`
- instruction: `internal/cli/deliver_conflict_test.go`
- instruction: `internal/cli/implement_test.go`
- instruction: `internal/testfixture/fixture.go`
- instruction: `internal/daemon/daemon.go`
- instruction: `internal/config/config.go`
- instruction: `docs/user-guide/commands/deliver.md`
- instruction: `.agents/skills/roundfix/references/deliver.md`
- instruction: `docs/adr/0225-a-queue-item-runs-the-roundfix-binary-its-branch-builds.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestDeliveryStepRunsTheItemBinary|TestDeliveryStepsStartTheItemBinaryWithTheOwnersArguments|TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate|TestDeliveryStepParksWhenTheItemBuildFails|TestDeliveryStepParksWhenTheItemBinaryPathIsNotIgnored|TestDeliveryStepWithoutADeclarationRunsTheOwnerExecutable)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestDeliveryStepRunsTheItemBinary TestDeliveryStepsStartTheItemBinaryWithTheOwnersArguments TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate TestDeliveryStepParksWhenTheItemBuildFails TestDeliveryStepParksWhenTheItemBinaryPathIsNotIgnored TestDeliveryStepWithoutADeclarationRunsTheOwnerExecutable; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the six tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestNoTestWritesAnExecutableOutsideTheResidue)$" ./internal/testfixture 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; printf "%s\\n" "$out" | grep -q -- "--- PASS: TestNoTestWritesAnExecutableOutsideTheResidue" || exit 1; test -f internal/cli/deliver_item_binary_test.go || { printf 'missing %s\n' internal/cli/deliver_item_binary_test.go >&2; exit 1; }` — expected: exit 0; the unchanged executable-residue guard passes with the new test file present, which does not exist before this Task, so the command fails.

## References

- `_prd.md` → User Stories 1-3; Core Features 2-6; Success Metrics 1-3; Goals 1-3
- `_techspec.md` → The step executable; Interfaces; API Contract 2; API Contract 3; API Contract 5; Testing Approach 3; Build Order 4; Risks & Considerations
- ADR-0225; ADR-0192; ADR-0125; ADR-0213

## Result

Implemented the Task 04 slice in `internal/cli/deliver_workflow.go` and added
`internal/cli/deliver_item_binary_test.go`. The owner uses only its loaded
`delivery.item_binary` declaration. Before each child step it checks the Git
ignore rule, builds through `daemon.ExecVerifier`, probes with the existing
child environment, and selects the item or owner executable with the specified
console line. Selection errors return before a step starts. The shared review
helper also applies this selection when recording an omitted review.

Focused evidence:

- Before implementation,
  `GOCACHE=/private/tmp/roundfix-task04-go-cache rtk proxy go test ./internal/cli -run '^TestDeliveryItemBinaryThatCannotStartReturnsAParkError$' -count=1`
  failed to compile because `log` and `stepExecutable` did not exist.
- After implementation and correction of the disposable fixture's Specs Root,
  `GOCACHE=/private/tmp/roundfix-task04-go-cache rtk proxy go test ./internal/cli -run '^TestDelivery(Step|Steps|ItemBinary)' -count=1 -v`
  exited 0; all seven new tests passed.
- `GOCACHE=/private/tmp/roundfix-task04-go-cache rtk proxy make verify-incremental`
  exited 0 with the access needed by existing integration tests. Formatting,
  `go vet`, package tests, skill checks and the build passed. Output was captured
  in `/private/tmp/roundfix-task04-incremental.log`.
- The initial focused attempt could not access the host Go cache; subsequent
  checks used the task cache above. The initial sandboxed incremental attempt
  was interrupted by blocked access to `cafe.github.com`; the successful rerun
  used an approved sandbox escalation.

| Acceptance criterion | Implementation and focused-check evidence |
| --- | --- |
| Declaration and probe exit 0 select the built binary for all three steps, preserving arguments and worktree, with one API Contract 2 line per step | `TestDeliveryStepRunsTheItemBinary` checks the selected absolute path and exact line. `TestDeliveryStepsStartTheItemBinaryWithTheOwnersArguments` passed its implement, archive and review subtests, recording the probe followed by the specified step arguments, item directory and disposable HOME. Removing the output before each step proves each step rebuilds it. The fixture intentionally exits after recording the step. |
| Probe exit 2 selects and starts the owner's executable, with one API Contract 3 notice and no fixture step | `TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate` passed stderr and stdout subtests. Each asserts the exact notice and first non-empty diagnostic line, starts the selected owner through the existing CLI test helper, and records only the fixture probe. |
| Failed build or unignored path returns API Contract 5's error and starts nothing; an absent declaration runs and logs nothing new | `TestDeliveryStepParksWhenTheItemBuildFails` checks exit 7, the verifier error, retained build log path and contents, and no fixture start or selection line. `TestDeliveryStepParksWhenTheItemBinaryPathIsNotIgnored` checks the exact refusal, no build marker, no fixture start and no line. `TestDeliveryStepWithoutADeclarationRunsTheOwnerExecutable` uses a nil Git runner, checks the owner selection and execution, and checks no build output, fixture start or line. |

`TestDeliveryItemBinaryThatCannotStartReturnsAParkError` additionally checks
that a successful build which leaves the binary absent returns the specified
absolute-path start error through Archive, with no step or selection line.
Each disposable fixture compiles its fake once through
`testfixture.FixtureBinary` and reuses it for its probes and steps. The fake
starts no detached or waiting process; each invocation is synchronously waited
through the test context. No new test writes an executable with a literal mode.

No follow-up implementation was added. Task status and authored checkboxes
remain unchanged by this Agent. The declared `## Verification` commands were
not run; the Daemon owns Verification and settlement. No commit, push or Pull
Request was created.

## Carry-forward provenance

- Source Run: `run_20261003T234640Z_3e854d86e9475f64`
- Source commit: `42f8e884c26f79f643a07d438f13cdacbeb56761`
