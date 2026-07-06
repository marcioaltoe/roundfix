---
task: task_01
spec: 0010-run-robustness
status: completed
type: backend
complexity: medium
---

# Task 01: Warn and ignore removed config keys

## Overview

Implement ADR-0027: a deprecated-keys registry in the config package strips
recognized removed keys from the document before strict decoding, emitting
one stderr warning naming the replacement — starting with the
`resolve.concurrent` entry that hard-broke the dogfood machine's config.
Unknown keys keep failing. Verifiable through config table tests replaying
the exact 0009 failure.

## Requirements

1. MUST add a deprecated-keys registry (path → replacement hint) consulted
   on every config load: registered paths present in a document are removed
   from the parsed node tree before the strict `KnownFields` decode, each
   emitting exactly one stderr warning per load in the shape
   `config: <key> is deprecated and ignored; use <replacement>`.
2. MUST register `resolve.concurrent` → `worktree.concurrency` and delete
   the current hard-fail rejection.
3. MUST keep truly unknown keys failing strict validation exactly as today
   (the registry is the only bypass), covered for both User Config and
   Project Config, including documents mixing deprecated and valid keys.
4. MUST leave every other config behavior, default, and generated output
   byte-stable.

## Subtasks

- [x] Registry and pre-decode stripping
- [x] resolve.concurrent entry replacing the hard-fail
- [x] Warning line shape and once-per-load semantics
- [x] Table tests: deprecated-only, mixed, unknown-key still fails, both
      config files

## Acceptance Criteria

- [x] The exact 0009 corpse-config (a User Config carrying
      `resolve.concurrent: 1`) loads successfully with the single warning
      line asserted verbatim, and a Run proceeds past Preflight.
- [x] An unknown key (`resolve.concurent` typo) still fails with today's
      strict error.
- [x] Warnings appear once per load, on stderr only; stdout untouched.
- [x] Full suite passes with only the deliberate hard-fail-test replacement.

## Verification

- `rtk go test ./internal/config/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 1; Core Feature 1; Success Metrics. `_techspec.md` →
Config deprecation, Build Order 1. ADR-0027. Work-plan finding R3-5.

## Result

Implemented the config deprecation registry for `resolve.concurrent` →
`worktree.concurrency`. Config loads now parse to a YAML node tree, remove
registered deprecated paths before strict `KnownFields` decoding, and emit
`config: resolve.concurrent is deprecated and ignored; use worktree.concurrency`
once per `Load` call on stderr.

Evidence by acceptance criterion:

- Exact 0009 corpse-config: `TestLoadWarnsAndIgnoresDeprecatedConfigKeys/user config exact 0009 corpse key` asserts a User Config with `resolve.concurrent: 1` loads and emits the warning verbatim. `TestRunFetchWarnsAndIgnoresDeprecatedUserConfig` asserts `roundfix fetch` exits 0, creates the normal fetch output, and does not print `Preflight failed`.
- Unknown typo still strict: `TestLoadRejectsUnknownConfigKeys` covers User Config, Project Config, and a deprecated-plus-typo mix; each fails with `resolve.concurent is not a supported config key`.
- Warning stream and once semantics: `TestLoadWarnsAndIgnoresDeprecatedConfigKeys/user and project config warn once per load` asserts one warning when both files contain the removed key. The CLI regression asserts stdout does not contain the warning and stderr contains it exactly once.
- Full suite: the former hard-fail expectation for `resolve.concurrent` was replaced with positive deprecation coverage plus the strict unknown-key typo checks.

Verification:

- `rtk go test ./internal/config/`: passed, 33 tests.
- `rtk go test ./...`: passed, 721 tests across 17 packages.
- `rtk make verify`: passed; it ran `rtk go test ./...`, `rtk go run -buildvcs=false ./cmd/roundfix skills check`, and `rtk go build -buildvcs=false -o bin/roundfix ./cmd/roundfix`.
