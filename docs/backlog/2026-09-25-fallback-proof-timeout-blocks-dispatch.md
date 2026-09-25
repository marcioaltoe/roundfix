---
type: fix
status: open
created: 2026-09-25
spec: null
reason: null
---

# An unused fallback's proof timing out under load blocks every dispatch

## Symptom

Implement preflight proves every configured fallback tuple sequentially with a 30 s setup limit and no retry; under load the unused claude fallback timed out and blocked dispatch, and the advice tells the operator to reconfigure a working profile.

## Where

`internal/cli/profile_preflight.go`, `profiles_validate.go` (`proveProfileSelectionsWithOptions`, `profileProofClassification`), `internal/agent/acpx_runner.go` (`acpxPreflightSetupTimeout`).

## Expected

Classify a proof deadline as temporary with retry advice and retry once; the lazy fallback proof stays with Spec 0123 Core Feature 4.

## Evidence

Session 2026-09-24: four dispatches refused with `adapter error: context deadline exceeded` for `claude`/`opus`/`high`.
