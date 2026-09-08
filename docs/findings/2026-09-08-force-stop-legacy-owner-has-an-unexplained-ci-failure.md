---
status: pending
created_at: 2026-09-08
updated_at: 2026-09-08
kind: finding
---

# Force Stop — the legacy owner test has an unexplained CI failure (2026-09-08)

A documentation-only PR failed the legacy owner Force Stop test once in CI.
The originating session lost the assertion detail during reruns and could not
reproduce locally. This remains an investigation, not a diagnosed timeout bug.

Source: Secondbrain `inbox/roundfix/_triaged/2026-09-02-flake-em-force-stop-legacy-owner-so-em-ci.md`.
The original capture and its dated evidence remain intact there.

## 1. Observed behavior

- Symptom / evidence: `TestRunForceStopLegacyRunWithoutOwnerIdentityStillStopsOwner` in `internal/cli/orphan_unix_test.go` remains parallel and waits up to two seconds for owner exit. The 2026-09-02 capture reports 25 focused runs and four package runs passing locally; those are historical results, not this triage's test runs.
- Root cause: `unknown`. The original assertion was not preserved, so host saturation and the wait budget remain hypotheses.
- Action / suggestion: Route to provisional P5, Verification capacity and measured economics. Preserve the next failing CI log before rerunning, characterize the failing boundary, and change behavior only when evidence identifies a cause.

The implementation group is provisional; no Spec has been assigned and no
implementation or terminal verification is claimed by this triage.
