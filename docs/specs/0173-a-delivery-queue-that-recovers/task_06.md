---
task: task_06
spec: 0173-a-delivery-queue-that-recovers
status: pending
type: backend
complexity: medium
---

# Task 06: A timed-out profile proof is retried once and called temporary

## Overview

`proveProfileSelectionsWithOptions` in `internal/cli/profiles_validate.go` proves every configured tuple, fallbacks included, once each, and each proof runs under the 30-second `acpxPreflightSetupTimeout` in `internal/agent`. Under load on 2026-09-24 the unused `claude`/`opus`/`high` fallback's proof expired four times and refused four dispatches with `adapter error: context deadline exceeded`; `profileProofClassification` returns no classification for a deadline and `profileProofNextAction` advises reconfiguring a profile that works. The proof talks to the ACP Runtime adapter through disposable sessions; its result is read by `roundfix profiles validate`, the Doctor Command's profile check and every operational Preflight, and reaches the operator on stderr and in `roundfix/profiles-validate/v1` JSON.

## Requirements

1. MUST prove each tuple through a helper that retries once: when `proveProfileSelection` returns an error for which `errors.Is(err, context.DeadlineExceeded)` holds while the command context's `Err()` is nil, the helper proves the tuple again with a fresh disposable session.
2. MUST return, when the retry also ends with such a timeout, a `profileProofTimeoutError` wrapping it, which `profileProofClassification` classifies `temporary` before it treats a deadline as unclassified, and for which `profileProofNextAction` advises rerunning the command when load drops, states that the configured profile was not shown to be wrong, and does not name `roundfix profiles configure`.
3. MUST return after one attempt, unchanged, a cancelled or expired command context, a runtime construction error, an access-policy error and every other proof error.
4. MUST keep `TestProfilesValidateFailedProofNamesTupleAffectedCategoriesAndRecovery` and `TestProfileOperationalPreflightMatchesProfilesValidateClassifiedFailure` green without edits.
5. MUST state in the profiles section of `docs/user-guide/commands.md` and in the profile Preflight paragraph of `.agents/skills/roundfix/SKILL.md` that a proof whose setup times out is retried once (the phrase `setup times out is retried once`) and that a second timeout is classified `temporary`, and regenerate `skills/roundfix/SKILL.md` with `make skills-sync`.
6. MUST put the new tests in `internal/cli/profile_proof_retry_test.go`, driving `roundfix profiles validate --json` and the operational Preflight with a fake Agent runner that counts proof attempts per tuple.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A proof that times out once and then passes leaves `profiles validate` at exit `0` after two attempts on that tuple.
- [ ] A proof that times out twice exits `2` with classification `temporary` and rerun advice that names no reconfiguration.
- [ ] A rejected selection and a cancelled command are each attempted once and keep their current classification.
- [ ] An operational Preflight whose fallback proof times out once still passes.

## Context

- interface: `internal/cli/profiles_validate.go`
- creates: `internal/cli/profile_proof_retry_test.go`
- interface: `docs/user-guide/commands.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestProfileProofRetriesATimedOutSetupOnce|TestProfileProofReportsASecondTimeoutAsTemporary|TestProfileProofDoesNotRetryARejectedSelection|TestProfileProofDoesNotRetryACancelledCommand|TestOperationalPreflightPassesAfterAFallbackProofTimesOutOnce|TestProfilesValidateFailedProofNamesTupleAffectedCategoriesAndRecovery|TestProfileOperationalPreflightMatchesProfilesValidateClassifiedFailure)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestProfileProofRetriesATimedOutSetupOnce TestProfileProofReportsASecondTimeoutAsTemporary TestProfileProofDoesNotRetryARejectedSelection TestProfileProofDoesNotRetryACancelledCommand TestOperationalPreflightPassesAfterAFallbackProofTimesOutOnce TestProfilesValidateFailedProofNamesTupleAffectedCategoriesAndRecovery TestProfileOperationalPreflightMatchesProfilesValidateClassifiedFailure; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && grep -q "setup times out is retried once" docs/user-guide/commands.md && grep -q "setup times out is retried once" .agents/skills/roundfix/SKILL.md && diff -r .agents/skills/roundfix skills/roundfix >/dev/null` — expected: exit 0; before this Task none of the five new named tests exists and the phrase is documented nowhere, so the command fails.

## References

- [_techspec.md](_techspec.md) — Proof timeouts
