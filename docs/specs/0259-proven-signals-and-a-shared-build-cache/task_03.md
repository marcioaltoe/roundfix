---
task: task_03
spec: 0259-proven-signals-and-a-shared-build-cache
status: pending
type: infra
complexity: medium
---

# Task 03: make shares one Go build cache that the sandbox can write

## Overview

The Makefile defaults `GOCACHE` to the worktree's `.gocache`, so every Run,
Task and item worktree builds its own cache of up to 8 GB. A probe under the
Codex workspace sandbox showed `$TMPDIR` and `/tmp` writable and the Go
default cache denied (`_techspec.md` → Current behavior). This Task makes
`make` use one cache under the temporary directory and keeps the worktree
cache only as the fallback. It is verifiable on its own: the Makefile's `test`
target, run with a stub Go, exports the expected path in each case.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-09,
   [each worktree keeps its own multi-gigabyte Go build cache](references/2026-10-09-each-worktree-keeps-its-own-go-build-cache.md),
   by sharing one cache between worktrees.
2. MUST change only the `GOCACHE` default of the Makefile, following
   `_techspec.md` → API Contract 7: computed once, `roundfix-go-build-<uid>`
   under `$TMPDIR` without its trailing slash or under `/tmp`, used when the
   directory can be created and a file created and removed in it, otherwise
   `$(CURDIR)/.gocache`; an explicit `GOCACHE` from the environment or the
   command line is exported unchanged. The bounded file of the cycle grant for
   this Task is `Makefile`, and `docs/agents/specific-repository.md` for the
   guide line; no other Makefile target, no CI workflow and no ignore file
   changes.
3. MUST add `TestMakeSharesOneGoBuildCache` with subtests `shared`,
   `fallback` and `explicit`. Each runs the repository Makefile's `test`
   target with `GO` set to a stub that records the `GOCACHE` it receives, in a
   temporary directory, with `GOCACHE` removed from the child environment
   unless the case sets it: `shared` sets `TMPDIR` to a writable temporary
   directory and expects `<TMPDIR>/roundfix-go-build-<uid>`, which exists
   afterwards; `fallback` sets `TMPDIR` to a directory without write
   permission and expects `<directory of the run>/.gocache`; `explicit` sets
   `GOCACHE` and expects it unchanged. The test lives beside the existing
   Makefile wiring tests, in a package that already installs the suite guard.
4. MUST add one line to `docs/agents/specific-repository.md`, outside every
   setup marker, naming the shared `roundfix-go-build-<uid>` cache, its
   fallback, that an explicit `GOCACHE` wins, and that a sandboxed Agent sets
   nothing.
5. MUST NOT remove or rewrite the `.gocache` entry of `.gitignore`, which the
   fallback still needs.

## Subtasks

- [ ] Change the `GOCACHE` default with its writability fallback.
- [ ] Add the Makefile test with its three cases.
- [ ] Add the repository guide line.

## Acceptance Criteria

- [ ] With `GOCACHE` unset and a writable `TMPDIR`, `make test` passes the shared cache to Go.
- [ ] With an unwritable `TMPDIR`, it passes the worktree's `.gocache`.
- [ ] An explicit `GOCACHE` reaches Go unchanged.

## Context

- instruction: `docs/adr/0266-a-signal-goes-only-to-a-proven-process-and-worktrees-share-one-build-cache.md`
- instruction: `internal/verifyselect/full_contract_run_test.go`
- interface: `Makefile`
- creates: `internal/verifyselect/build_cache_makefile_test.go`
- interface: `docs/agents/specific-repository.md`

## Verification

- `tr -s '[:space:]' ' ' < docs/agents/specific-repository.md | grep -qF -- 'roundfix-go-build-' || { printf 'the repository guide does not name the shared cache\n' >&2; exit 1; }; out="$(go test -count=1 -v -run '^TestMakeSharesOneGoBuildCache$' ./internal/verifyselect 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestMakeSharesOneGoBuildCache TestMakeSharesOneGoBuildCache/shared TestMakeSharesOneGoBuildCache/fallback TestMakeSharesOneGoBuildCache/explicit; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 ./internal/verifyselect || exit 1` — expected: exit 0. Before this Task the guide line and the test are absent, so the command fails; afterwards the three cases and the existing Makefile wiring tests pass.

## References

- `_prd.md` → Goals; User Story 4; Core Feature 6; Success Metric 4; Acceptance evidence
- `_techspec.md` → Current behavior; API Contract 7; Invariant 7; Testing Approach; Build Order 3
- ADR-0266
