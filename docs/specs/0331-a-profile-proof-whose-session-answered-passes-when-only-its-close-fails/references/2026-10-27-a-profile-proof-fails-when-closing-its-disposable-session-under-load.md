---
type: fix
status: promoted
created: 2026-10-27
spec: 0331-a-profile-proof-whose-session-answered-passes-when-only-its-close-fails
reason: null
---

# A profile proof fails when closing its disposable session under load

## Symptom

Under high host load, Agent Selection readiness fails with
`profile proof setup timed out on 3 attempts of 2m0s each` even when the
disposable session answered. In one case the failing step was
`close disposable Agent Session roundfix-preflight-09678fa0be60a4a1: close acpx Agent Session ...: context deadline exceeded`.
The CLI classifies the failure as temporary and asks for a readiness retry, so
a pre-PR review or a delivery item parks until someone retries it.

## Where

`internal/agent/selection_assignment.go` and `internal/agent/fallback.go`:
after setup, the close runs under its own `disposableCleanupTimeout()`, and a
close error fails the proof even when setup succeeded
(`AgentSessionCleanupError`, classification `SessionCleanupFailed`).

## Expected

Decide whether a proof whose session answered and only failed to close should
pass with a recorded warning and a reaped session, or keep failing, and
whether the close deadline should follow `defaults.profile_proof_setup_timeout`.
The cause is host load, so the fix may also be to say so: the refusal names
the load average or the step that timed out.

## Evidence

- Vortex (ai-memory message `01a17b18-94f6-7950-b747-7811c9f8d3cc`, 2026-10-26
  23:02 BRT): Roundfix 0.124.0, second `roundfix review --base origin/main`
  at candidate `4f943f6f`, load average 54.35 at 22:55 BRT; the failing step
  was the close.
- Roundfix's own Delivery Queue: 0324 parked `review-blocked` (22:55 BRT) and
  0322 parked `delivery-error` at preflight (23:53 BRT) with the same 3×2m
  profile proof timeout while four Runs ran on the same host. Both passed on a
  plain retry.
