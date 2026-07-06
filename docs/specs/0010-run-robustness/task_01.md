---
task: task_01
spec: 0010-run-robustness
status: pending
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

- [ ] Registry and pre-decode stripping
- [ ] resolve.concurrent entry replacing the hard-fail
- [ ] Warning line shape and once-per-load semantics
- [ ] Table tests: deprecated-only, mixed, unknown-key still fails, both
      config files

## Acceptance Criteria

- [ ] The exact 0009 corpse-config (a User Config carrying
      `resolve.concurrent: 1`) loads successfully with the single warning
      line asserted verbatim, and a Run proceeds past Preflight.
- [ ] An unknown key (`resolve.concurent` typo) still fails with today's
      strict error.
- [ ] Warnings appear once per load, on stderr only; stdout untouched.
- [ ] Full suite passes with only the deliberate hard-fail-test replacement.

## Verification

- `rtk go test ./internal/config/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 1; Core Feature 1; Success Metrics. `_techspec.md` →
Config deprecation, Build Order 1. ADR-0027. Work-plan finding R3-5.
