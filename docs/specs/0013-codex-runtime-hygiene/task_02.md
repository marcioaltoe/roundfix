---
task: task_02
spec: 0013-codex-runtime-hygiene
status: pending
type: backend
complexity: medium
---

# Task 02: Codex hygiene inspector (quarantine and notarization)

## Overview

Add the macOS-only inspector that decides whether the codex a Run would use is
safe: it resolves the codex on PATH (and any configured codex path), inspects
the `com.apple.quarantine` attribute and notarization/signing acceptance, and
reports a failure carrying the curl-reinstall fix. On other platforms it reports
not-applicable. This is the shared engine for both the Doctor check and the
verified-clean spawn.

## Requirements

1. MUST resolve the codex the same way a Run would: a configured codex path when
   set, else the PATH entry.
2. MUST, on macOS, detect the `com.apple.quarantine` extended attribute and
   assess notarization/signing acceptance with strictly read-only probes (no
   mutation, no de-quarantine).
3. MUST report a failure with the curl-to-`~/.local/bin` reinstall command as the
   next action when the binary is quarantined or not accepted; report clean
   otherwise.
4. MUST report not-applicable on non-darwin platforms and never fail there; the
   quarantine/assessment probes MUST be injectable for tests.

## Subtasks

- [ ] Codex resolution (configured path then PATH)
- [ ] macOS quarantine-attribute and acceptance probes behind interfaces
- [ ] Result mapping: quarantined/not-accepted → failure + curl next action
- [ ] macOS gate returning not-applicable elsewhere
- [ ] Table tests: {quarantined, not-accepted, clean, non-darwin}

## Acceptance Criteria

- [ ] A quarantined codex yields a failure result naming the curl-reinstall fix.
- [ ] An un-notarized/not-accepted codex yields a failure result.
- [ ] A clean codex yields a passing result.
- [ ] On non-darwin, the inspector reports not-applicable and never fails.
- [ ] The `com.apple.quarantine` and acceptance probes are injected in tests (no real `xattr`/`spctl`).

## Verification

- `rtk go test ./internal/... -run Codex` — expected: the hygiene inspector table tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 1-2, 5; Core Feature 2. `_techspec.md` → Codex hygiene
check, Build Order 2, Interfaces: `CodexHygiene`. ADR-0032. Work-plan finding
R3-8.
