---
task: task_01
spec: 0215-a-setup-and-doctor-that-recognize-what-the-machine-and-the-repository-need
status: pending
type: backend
complexity: high
---

# Task 01: Doctor reports the forge and Git with coded findings

## Overview

Doctor proves the ACP side of a machine and says nothing about the GitHub
CLI, Git or the remote the Delivery Queue publishes through. This Task adds
the readiness findings, the `warn` status and the `gh`, `git` and `remote`
lines of the TechSpec's "The forge lines", printed after `pre-pr-review`, with
the codes of "Finding codes" and the bounded forge reads of ADR-0220. It is
verifiable alone through Doctor with scripted runners and fake executables.

## Requirements

1. MUST add `CheckStatusWarn` and the five line names to the existing status
   and line-name constants, and `readinessFinding` with `readinessResult` as
   the TechSpec's Interfaces state, folding statuses by Invariant 1.
2. MUST print a `warn` line's next action exactly as a `failed` line's, and
   MUST keep `warn` out of Doctor's exit code.
3. MUST implement the `gh`, `git` and `remote` lines and every one of their
   codes in "Finding codes" as "The forge lines" states: the delivery remote
   resolved as the Delivery Queue resolves it, the three URL forms, the gh
   floor 2.81.0, `gh auth status --active --hostname <host> --json hosts`, the
   permission read with `gh repo view <owner/repo> --json viewerPermission`,
   the Git floor 2.23.0, and the identity keys named without their values.
4. MUST run every `gh` and `git` child through `readinessRunner` with no
   standard input, `GIT_TERMINAL_PROMPT=0` in its environment, and a
   ten-second limit that cancels the child, and MUST NOT pass `--show-token`
   or print a login token, a scope list or a Git email address.
5. MUST report `warn`, never `failed`, for a read that timed out or could not
   connect, by Invariant 2: `DR-GH-UNREACHABLE`,
   `DR-GH-PERMISSION-UNVERIFIED` and `DR-REMOTE-UNREACHABLE`.
6. MUST add the `readiness` field to `doctorDependencies`, defaulting to the
   real lines, and MUST give `withDoctorFakeDeps` and its siblings ready fakes
   so the existing Doctor tests reach no network; their expected standard
   output gains the new `ok` lines and nothing else changes.
7. MUST change Doctor's help text, which says Doctor is offline, to say that
   only the `gh` and `remote` lines contact the repository's forge, each read
   bounded, while keeping the words "read-only" and "mutates nothing".
8. MUST describe the three lines, every code they print, the `warn` status and
   the bounded reads in `docs/user-guide/commands/doctor.md`.
9. MUST NOT change `internal/cli/cli_test.go`, the Run Database, or any check
   line other than the new ones.
10. MUST add the tests named in Verification to the new file
    `internal/cli/readiness_forge_test.go`: each code through a scripted
    runner; the `warn` fold and its printed next action; Surface Transcript 2
    through `runCLI`; and a probe test that puts fake `gh` and `git` scripts
    first on `PATH`, records their argv, environment and standard input, and
    proves a shortened timeout cancels a sleeping child.

## Subtasks

- [ ] Add the findings, the fold and the `warn` status.
- [ ] Implement the `gh`, `git` and `remote` lines over the bounded runner.
- [ ] Wire them into Doctor after `pre-pr-review` and update the help text.
- [ ] Give the existing Doctor tests ready fakes.
- [ ] Describe the lines and codes in the `doctor` guide.
- [ ] Add the forge readiness tests.

## Acceptance Criteria

- [ ] Each of `DR-GH-MISSING`, `DR-GH-VERSION`, `DR-GH-UNAUTHENTICATED`,
      `DR-GH-TOKEN-REJECTED`, `DR-GH-PERMISSION`, `DR-GIT-MISSING`,
      `DR-GIT-VERSION`, `DR-GIT-IDENTITY`, `DR-REMOTE-MISSING` and
      `DR-REMOTE-FORGE` makes its line `failed` with its next action and
      Doctor exit `1`.
- [ ] A timed-out or refused read yields `DR-GH-UNREACHABLE`,
      `DR-GH-PERMISSION-UNVERIFIED` or `DR-REMOTE-UNREACHABLE` as `warn`, with
      its next action printed, and Doctor exits `0` when nothing else fails,
      as Surface Transcript 2 states.
- [ ] The fake `gh` never receives `--show-token`, every child gets empty
      standard input and `GIT_TERMINAL_PROMPT=0`, and a sleeping child is
      cancelled at the bound.
- [ ] The existing Doctor tests pass with five more `ok` lines.

## Context

- creates: `internal/cli/readiness.go`
- creates: `internal/cli/readiness_forge.go`
- creates: `internal/cli/readiness_forge_test.go`
- interface: `internal/cli/health.go`
- interface: `internal/cli/doctor.go`
- interface: `internal/cli/doctor_test.go`
- interface: `internal/cli/doctor_characterization_test.go`
- interface: `internal/cli/doctor_recommendations_test.go`
- interface: `internal/cli/cli.go`
- interface: `docs/user-guide/commands/doctor.md`
- instruction: `internal/cli/deliver_workflow.go`
- instruction: `internal/preflight/preflight.go`
- instruction: `internal/delivery/github.go`
- instruction: `docs/adr/0220-doctor-checks-what-delivery-needs-and-only-a-definite-answer-fails.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestForgeReadinessReportsEachDefiniteFailure|TestForgeReadinessWarnsWhenTheForgeDoesNotAnswer|TestGitReadinessReportsVersionAndIdentity|TestForgeProbesAreBoundedAndNeverPromptOrPrintAToken|TestDoctorPrintsForgeWarningsAndExitsZero)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestForgeReadinessReportsEachDefiniteFailure TestForgeReadinessWarnsWhenTheForgeDoesNotAnswer TestGitReadinessReportsVersionAndIdentity TestForgeProbesAreBoundedAndNeverPromptOrPrintAToken TestDoctorPrintsForgeWarningsAndExitsZero; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the five tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestRunDoctorProfileReadinessProvesEffectiveCategoriesAndReportsCounts|TestCharacterizationInvariantDoctorCountsAreUnchangedWithoutOptionalCategories|TestForgeReadinessReportsEachDefiniteFailure)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRunDoctorProfileReadinessProvesEffectiveCategoriesAndReportsCounts TestCharacterizationInvariantDoctorCountsAreUnchangedWithoutOptionalCategories TestForgeReadinessReportsEachDefiniteFailure; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; the existing Doctor tests pass beside a new one, which does not exist before this Task.
- `for phrase in "DR-GH-UNAUTHENTICATED" "DR-GH-TOKEN-REJECTED" "DR-GIT-IDENTITY" "DR-REMOTE-UNREACHABLE" "gh auth status"; do tr -s '[:space:]' ' ' < docs/user-guide/commands/doctor.md | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/doctor.md "$phrase" >&2; exit 1; }; done` — expected: exit 0; before this Task the guide names none of these, so the command fails.

## References

- `_prd.md` → User Stories 1-2; Core Features 1-3; Success Metrics 1-2
- `_techspec.md` → Interfaces; Invariants 1-4; The forge lines; Finding codes; API Contract 1; Surface Transcript 2; Testing Approach 1-3; Build Order 1
- ADR-0220
