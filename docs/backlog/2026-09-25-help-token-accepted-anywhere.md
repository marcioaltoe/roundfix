---
type: fix
status: open
created: 2026-09-25
spec: null
reason: null
---

# A bare `help` anywhere in the arguments prints usage and exits 0

## Symptom

`roundfix archive --bogus help` exits 0 (without `help` it exits 2); `implement --spec help` prints usage; no flag can take the value `help`.

## Where

`internal/cli/cli.go` `commandWantsHelp` and its 46 call sites.

## Expected

Only a leading or sole help token triggers usage, or flag parsing owns `-h`.

## Evidence

secondbrain `inbox/roundfix/2026-09-19-help-token-accepted-anywhere.md`.
