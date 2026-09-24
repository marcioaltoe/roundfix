---
task: task_07
spec: 0158-daemon-verification-and-access-readiness
status: completed
type: backend
complexity: low
---

# Task 07: Keep every failure and name a degraded access

## Overview

Corrective Task from the pre-PR review of 2026-09-24. In independent mode a temporary failure or unknown cause discards the deterministic failures already collected, so they never reach the repair turn or the reason. And when the Codex full-access sandbox preset is unavailable, readiness records an unqualified `full-access` effective policy.

## Requirements

1. MUST keep the failures already collected in independent mode when a later command ends in a temporary failure or unknown cause, and report all of them in the outcome and its reason.
2. MUST record a degraded effective policy, with the unavailable sandbox preset named, when full access is granted without it, and surface it in the readiness output.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] A deterministic failure beside a later temporary failure is kept and reported.
- [ ] A granted mode without its sandbox preset is recorded as degraded.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/daemon/engine.go`
- interface: `internal/agent/acpx_runner.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestIndependentVerificationKeepsCollectedFailuresBesideATemporaryOne|TestProofRecordsADegradedFullAccess|TestIndependentVerificationHandsEveryFailureToRepair)$" ./internal/daemon ./internal/agent 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestIndependentVerificationKeepsCollectedFailuresBesideATemporaryOne TestProofRecordsADegradedFullAccess TestIndependentVerificationHandsEveryFailureToRepair; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task two of the named cases do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Independent Verification

## Result

### Implementation

- Independent Verification now appends each observed command failure before a
  later temporary failure stops the attempt. The exclusive retry carries any
  earlier deterministic failures forward, so a passing retry still returns
  them to the one Verification Feedback turn. An unobserved later command also
  retains the collected failures, and the outcome reason names every retained
  command plus the unknown cause.
- Full-access application now returns the effective policy it established.
  When Codex accepts `full-access` but does not advertise the
  `danger-full-access` sandbox preset, the disposable proof records
  `full-access (degraded: danger-full-access sandbox preset unavailable)`.
  The existing profile proof report copies that value into readiness output.

### Focused checks

- Before the production changes,
  `rtk env GOCACHE=/private/tmp/roundfix-task07-gocache go test -count=1 -run '^TestIndependentVerificationKeepsCollectedFailuresBesideATemporaryOne$' ./internal/daemon`
  exited 1 because the passing exclusive retry discarded the earlier
  deterministic failure and started no repair turn.
- Before the production changes,
  `rtk env GOCACHE=/private/tmp/roundfix-task07-gocache go test -count=1 -run '^TestProofRecordsADegradedFullAccess$' ./internal/agent`
  exited 1 because the proof reported the unqualified policy `full-access`.
- After the production changes, both focused commands above exited 0.
- `rtk env GOCACHE=/private/tmp/roundfix-task07-gocache go test -count=1 -run '^TestIndependentVerificationKeepsCollectedFailuresBesideUnknownCause$' ./internal/daemon`
  exited 0 and proved the outcome reason retains the deterministic command,
  the unobserved command and its runner cause.
- `rtk env GOCACHE=/private/tmp/roundfix-task07-gocache go test -count=1 -run '^TestIndependentVerificationHandsEveryFailureToRepair$' ./internal/daemon`
  exited 0.
- `rtk env GOCACHE=/private/tmp/roundfix-task07-gocache go test -count=1 ./internal/daemon ./internal/agent`
  exited 0.
- `rtk env GOCACHE=/private/tmp/roundfix-task07-gocache make verify-incremental`
  exited 0 after rerunning with the network permission required by the
  repository's GitHub integration checks. The first restricted run was blocked
  by sandbox denial of `api.github.com` and produced no product verdict.
- The Task's declared `## Verification` command was not run; the Daemon owns
  that check.

### Acceptance evidence

- Deterministic failure beside a later temporary failure:
  `TestIndependentVerificationKeepsCollectedFailuresBesideATemporaryOne`
  passed after observing the initial attempt, exclusive retry and repaired
  attempt, and confirmed that the repair prompt retained the deterministic
  command and its diagnostic artifact.
- Granted mode without its sandbox preset:
  `TestProofRecordsADegradedFullAccess` passed after the disposable Codex
  session accepted `full-access`, refused the unavailable sandbox option and
  returned a degraded effective policy naming `danger-full-access`.
