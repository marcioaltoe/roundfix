---
task: task_04
spec: 0172-a-qa-gate-that-tells-the-truth
status: completed
type: backend
complexity: medium
---

# Task 04: An unobserved Verification is classified

## Overview

`publishUnknownFailure` in `internal/daemon/engine.go` publishes an unobserved Verification with its cause in the `error` prose and an optional `diagnostic_path`, and publishes the verdict with empty failure metadata, so neither event carries a classification. `projectVerificationRecord` in `internal/runevent/stream.go` projects an event as unknown only when `classification` is `verification_unknown`, and then requires `command`, `reason` and `diagnostic_path`. A Supervisor reading `roundfix events` therefore reads "we did not find out" as "the work is wrong", which ADR-0111 forbids. The payload is written by the Daemon into the Run Event Journal and read by every Supervisor following the stream.

## Requirements

1. MUST make `publishUnknownFailure` write `classification: verification_unknown`, `command`, `reason` and `diagnostic_path` on both the `failed` and the `verdict` event of an unobserved outcome; `reason` is the cause's error text, or `reason unavailable` when it has none, and `diagnostic_path` is the retained path, or `unavailable` when none was retained.
2. MUST let this classification override a request's preset `FailureClassification` and `FailureReason`: a precondition Verification the runner could not observe is unknown, not a repository that was red on entry.
3. MUST leave `publishFailedCommand` and `publishVerdict` for a command verdict unchanged: a deterministic failure stays unclassified, and a temporary or precondition failure keeps its classification.
4. MUST prove the contract through the projection: each new test captures the published payloads, wraps them as `daemon.verification` Run Events and asserts the fields `runevent.ProjectStreamEvent` returns, value by value.
5. MUST state in `docs/user-guide/commands.md` (events) and in the Supervisor Run Event Stream section of `.agents/skills/roundfix/SKILL.md` that an unobserved Verification adds `classification` `verification_unknown` with `command`, `reason` and `diagnostic_path`, correcting the skill's claim that command strings and diagnostic paths are never projected, and regenerate `skills/roundfix/SKILL.md` with `make skills-sync`.
6. MUST put the new tests in `internal/daemon/unknown_verification_projection_test.go`.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] Both events of an unobserved Verification project with classification `verification_unknown` and the cause's command, reason and diagnostic path.
- [ ] Without a retained diagnostic the projection carries `unavailable`, and an unobserved precondition Verification projects as unknown.
- [ ] A deterministic command failure projects with no classification.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/daemon/engine.go`
- interface: `internal/runevent/stream.go`
- interface: `docs/user-guide/commands.md`
- interface: `.agents/skills/roundfix/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestUnobservedVerificationProjectsAsUnknown|TestUnobservedVerificationWithoutDiagnosticProjectsUnavailable|TestUnobservedPreconditionVerificationProjectsAsUnknown|TestDeterministicVerificationFailureProjectsUnclassified)$" ./internal/daemon 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestUnobservedVerificationProjectsAsUnknown TestUnobservedVerificationWithoutDiagnosticProjectsUnavailable TestUnobservedPreconditionVerificationProjectsAsUnknown TestDeterministicVerificationFailureProjectsUnclassified; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && grep -q "verification_unknown" docs/user-guide/commands.md && grep -q "verification_unknown" .agents/skills/roundfix/SKILL.md && diff -r .agents/skills/roundfix skills/roundfix >/dev/null` — expected: exit 0; before this Task none of the named cases exists and neither document names `verification_unknown`, so the command fails.

## References

- [_techspec.md](_techspec.md) — The unknown classification

## Result

### Implementation

- `publishUnknownFailure` now writes `classification: verification_unknown`, the command, the runner cause or `reason unavailable`, and the retained diagnostic path or `unavailable` on both the `failed` and `verdict` events. This explicit metadata bypasses a request's preset precondition classification and reason.
- The existing `publishFailedCommand` and `publishVerdict` command-verdict path remains unchanged. The Task-owned projection suite in `internal/daemon/unknown_verification_projection_test.go` covers the ordinary unknown, missing-evidence, precondition-override, and deterministic negative cases; the existing QA-gate integration assertion now expects the classified unknown evidence.
- The user guide and the Roundfix skill now document the projected unknown fields. `make skills-sync` regenerated the shipped `skills/roundfix/SKILL.md` mirror from the canonical skill.

### Focused checks

- Before the production change, `rtk env GOCACHE=/private/tmp/roundfix-task04-gocache go test -count=1 -run '^TestUnobservedVerificationProjectsAsUnknown$' ./internal/daemon` failed with `record 0 classification = "", want "verification_unknown"`.
- `rtk env GOCACHE=/private/tmp/roundfix-task04-gocache go test -count=1 ./internal/daemon` passed after the implementation and integration-assertion update.
- `rtk make skills-sync` exited `0`; `rtk shasum .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` reported the same SHA-1 (`aa4e4aa6de44033397d0ad157e744ff3daa3690d`) for both copies.
- `rtk git diff --check` exited `0`.
- The Task's declared `## Verification` command was not run; the Daemon owns that gate.

### Acceptance evidence

- `TestUnobservedVerificationProjectsAsUnknown` captures both published payloads, projects them through `runevent.ProjectStreamEvent`, and proves the failed and verdict records carry `verification_unknown` with the exact command, cause text, and diagnostic path.
- `TestUnobservedVerificationWithoutDiagnosticProjectsUnavailable` proves missing cause text and diagnostics project as `reason unavailable` and `unavailable`; `TestUnobservedPreconditionVerificationProjectsAsUnknown` proves unknown evidence overrides preset precondition metadata on both records.
- `TestDeterministicVerificationFailureProjectsUnclassified` proves a deterministic command failure and its verdict still project with no classification or reason.
