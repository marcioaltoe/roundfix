---
task: task_02
spec: 0013-codex-runtime-hygiene
status: completed
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

- [x] Codex resolution (configured path then PATH)
- [x] macOS quarantine-attribute and acceptance probes behind interfaces
- [x] Result mapping: quarantined/not-accepted → failure + curl next action
- [x] macOS gate returning not-applicable elsewhere
- [x] Table tests: {quarantined, not-accepted, clean, non-darwin}

## Acceptance Criteria

- [x] A quarantined codex yields a failure result naming the curl-reinstall fix.
- [x] An un-notarized/not-accepted codex yields a failure result.
- [x] A clean codex yields a passing result.
- [x] On non-darwin, the inspector reports not-applicable and never fails.
- [x] The `com.apple.quarantine` and acceptance probes are injected in tests (no real `xattr`/`spctl`).

## Verification

- `rtk go test ./internal/... -run Codex` — expected: the hygiene inspector table tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 1-2, 5; Core Feature 2. `_techspec.md` → Codex hygiene
check, Build Order 2, Interfaces: `CodexHygiene`. ADR-0032. Work-plan finding
R3-8.

## Result

Implemented the reusable codex hygiene inspector with configured-path-then-PATH
resolution, macOS-only quarantine and Gatekeeper assessment probes, and
structured results carrying `CodexHygiene` details plus the curl-to-`~/.local/bin`
next action for unsafe codex binaries. The shared health checker now exposes a
`codex` check wrapper, but no Doctor command or spawn-path behavior was added in
this task.

Evidence:

- Quarantined codex: `TestCodexHygieneInspector/quarantined_configured_codex_fails_with_curl_reinstall`
  passed under `rtk go test ./internal/... -run Codex`.
- Not-accepted codex: `TestCodexHygieneInspector/not_accepted_path_codex_fails_with_curl_reinstall`
  passed under `rtk go test ./internal/... -run Codex`.
- Clean codex: `TestCodexHygieneInspector/clean_path_codex_passes` passed under
  `rtk go test ./internal/... -run Codex`.
- Non-darwin behavior: `TestCodexHygieneInspector/non_darwin_is_not_applicable_and_does_not_inspect`
  passed and asserts no PATH resolution or probe calls occur.
- Probe injection: the table tests use fake `QuarantineProbe` and
  `AcceptanceProbe` implementations; no real `xattr` or `spctl` is invoked.
- Required focused gate: `rtk go test ./internal/... -run Codex` passed with 10
  tests in 16 packages.
- Required full suite: `rtk go test ./...` passed with 774 tests in 18 packages.
- Repository gate: `rtk make verify` passed, including fmt-check, full tests,
  skills-sync-check, skills check, and build.

Follow-ups: later tasks still need to wire this inspector into the Doctor
Command and the codex-acp spawn path.
