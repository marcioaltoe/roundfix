---
task: task_01
spec: 0384-the-makefile-follows-the-go-build-cache-set-with-go-env
status: completed
type: infra
complexity: medium
---

# Task 01: The Makefile uses the Go build cache set with go env and falls back small with a warning

## Overview

On 2026-11-12 gesttione-drive reported that this repository's Go build cache
held 10 GB of the Mac mini's internal disk in `$TMPDIR/roundfix-go-build-501`
although the mini had set its cache on an external disk with
`go env -w GOCACHE`
([`references/2026-11-12-gesttione-drive-the-roundfix-go-build-cache-fills-the-mini-internal-disk.md`](references/2026-11-12-gesttione-drive-the-roundfix-go-build-cache-fills-the-mini-internal-disk.md),
reported by gesttione-drive). This Task makes the `Makefile` use a cache set
with `go env` when it is an existing writable directory, under the 10,240 MiB
default bound, and fall back to the cache it chose before, under a 2,048 MiB
default bound and with one warning line, when that cache is unusable. It adds
the test that proves both, changes the `go` stubs of the existing `Makefile`
tests so they answer the new `env` queries, and regenerates the setup
manifest. The authorized tooling paths are exactly `Makefile` and the derived
`docs/agents/setup-context.json` through its sanctioned command.

## Requirements

1. MUST replace the `Makefile` block that starts at
   `ifeq ($(origin GOCACHE),undefined)` and ends at its `endif` (the line
   before `# ADR-0370: bound only the cache make chose unless the caller opts
   in.`) with exactly these lines, per `_techspec.md` → API Contracts 1 to 3;
   continuation lines start with a tab, and the lines after the block keep
   their bytes:

   ```make
   ifeq ($(origin GOCACHE),undefined)
   # ADR-0394: a Go build cache set with `go env -w GOCACHE` comes first when it is
   # an existing writable directory, under the same bound. make never creates it:
   # when the disk that holds it is not mounted, make warns once and falls back to
   # the shared temporary cache under a 2048 MiB bound. With nothing set, make
   # chooses as ADR-0370 decided.
   GOCACHE_FROM_GO_ENV := $(shell [ -n "$$($(GO) env -changed GOCACHE 2>/dev/null)" ] || exit 0; \
   	cache=$$($(GO) env GOCACHE 2>/dev/null) || exit 0; \
   	case "$$cache" in (/*) ;; (*) exit 0;; esac; \
   	if [ -d "$$cache" ] && probe=$$(mktemp "$$cache/.write-test.XXXXXX" 2>/dev/null) && rm -f "$$probe" 2>/dev/null; then \
   		printf '%s' "$$cache"; \
   	fi)
   ifneq ($(GOCACHE_FROM_GO_ENV),)
   GOCACHE := $(GOCACHE_FROM_GO_ENV)
   GOCACHE_BOUND_MB ?= 10240
   else
   GOCACHE_FALLBACK := $(shell root="$${TMPDIR:-/tmp}"; cache="$${root%/}/roundfix-go-build-$$(id -u)"; \
   	if mkdir -p "$$cache" 2>/dev/null && probe=$$(mktemp "$$cache/.write-test.XXXXXX" 2>/dev/null) && rm -f "$$probe" 2>/dev/null; then \
   		printf '%s' "$$cache"; \
   	else \
   		printf '%s' "$(CURDIR)/.gocache"; \
   	fi)
   GOCACHE_DEFAULT_BOUND_MB := $(shell [ -n "$$($(GO) env -changed GOCACHE 2>/dev/null)" ] || { printf 10240; exit 0; }; \
   	printf 'gocache: warning: GOCACHE set with go env (%s) is not a writable directory; using %s with a 2048 MiB default bound\n' "$$($(GO) env GOCACHE 2>/dev/null)" '$(GOCACHE_FALLBACK)' >&2; \
   	printf 2048)
   GOCACHE := $(GOCACHE_FALLBACK)
   GOCACHE_BOUND_MB ?= $(GOCACHE_DEFAULT_BOUND_MB)
   endif
   endif
   ```

   The `case` patterns are written `(/*)` and `(*)` so that make's
   parenthesis count stays balanced inside `$(shell ...)`, and no `#` appears
   inside a `$(shell ...)`, because make ends the line there.
2. MUST hold these hazards, from `_techspec.md` → Invariants: `make` never runs
   `mkdir` on the path `go env` prints or on any parent of it, so a disk that
   is not mounted never becomes a directory on the internal disk; the probe
   file is removed at once; no `$(shell ...)` of the block runs after
   `GOCACHE` is assigned; with nothing set, the exported `GOCACHE` and the
   10240 default are byte-identical to today's; a `GOCACHE` in the environment
   or on the command line skips the block, as today, and stays unbounded unless
   `GOCACHE_BOUND_MB` is set.
3. MUST add `internal/verifyselect/build_cache_go_env_makefile_test.go` with the
   suite header comment the package uses (`// Suite:`, `// Invariant:`,
   `// Boundary IN:`, `// Boundary OUT:`) and one test,
   `TestMakeFollowsTheGoBuildCacheSetWithGoEnv`, whose top level and every
   subtest call `t.Parallel()` first, per `_techspec.md` → Testing Approach:
   - It compiles a `go` stub with `testfixture.FixtureBinary`. The stub runs
     the real go command (its path found with `exec.LookPath("go")` by the
     test and passed in an environment variable) with the same arguments when
     its first argument is `env`, copying its standard output and standard
     error and exiting 1 when it fails; for any other first argument it
     appends its arguments followed by `GOCACHE=<the GOCACHE it received>` as
     one JSON array to a record file named by another environment variable.
   - Each subtest runs
     `make --no-print-directory -f <repository Makefile> test GO=<stub> GO_TEST_PARALLEL=16`
     in its own `t.TempDir()`, with the caller's environment minus `GOCACHE`,
     `GOENV`, `TMPDIR`, `MAKEFLAGS`, `MFLAGS`, `MAKEOVERRIDES`, `MAKELEVEL`,
     `GOCACHE_BOUND_MB`, `GOCACHE_BOUND_IDLE` and the stub's two variables,
     plus `TMPDIR=<a second t.TempDir()>/`, `GOENV=<a file in the subtest
     directory>` holding `GOCACHE=<set path>` (empty in `nothing set`), and
     the stub's two variables. The set path is
     `<subtest dir>/external disk/cache/go-build` unless the subtest says
     otherwise; the space proves that make keeps the path whole.
   - It asserts the recorded invocations exactly: when a bound applies,
     `["run","-buildvcs=false","./cmd/gocache-bound","-dir",<cache>,"-max-mb",<bound>,"-idle","2h","GOCACHE="+<cache>]`
     first, then `["test","-parallel","16","./...","GOCACHE="+<cache>]`; and
     the whole standard error exactly. `<fallback>` is
     `<TMPDIR root>/roundfix-go-build-<id -u>` and `<warning>` is
     `gocache: warning: GOCACHE set with go env (<set path>) is not a writable directory; using <fallback> with a 2048 MiB default bound`
     followed by a newline.
   - Its seven subtests: `set and writable` (the set directory exists: the set
     path, bound 10240, empty standard error, and the set directory is empty
     afterwards, so no probe file is left); `set and not mounted` (a
     `Volumes` directory exists in the subtest directory and the set path is
     `<subtest dir>/Volumes/missing disk/cache/go-build`: `<fallback>`, bound
     2048, standard error exactly `<warning>`, and
     `<subtest dir>/Volumes/missing disk` does not exist afterwards);
     `set and read-only` (the set directory exists with mode `0555`, restored
     to `0700` in `t.Cleanup`: `<fallback>`, 2048, `<warning>`); `set to a file`
     (a regular file at the set path, which stays not writable as a directory
     even for root: `<fallback>`, 2048, `<warning>`); `set with a bound` (the
     set directory exists and `GOCACHE_BOUND_MB=4096` is in the environment:
     the set path, bound 4096, empty standard error); `environment wins` (the
     set directory exists and `GOCACHE=<subtest dir>/explicit cache` is in the
     environment: only the `test` invocation, with that path, and empty
     standard error); `nothing set` (an empty environment file: `<fallback>`,
     bound 10240, empty standard error).
4. MUST make each of these existing `go` stubs return at once, with no output
   and no record, when its first argument is `env`, and change nothing else in
   these tests (`_techspec.md` → Declared breaks):
   - `internal/verifyselect/build_cache_bound_makefile_test.go`: the line
     ` if len(os.Args) > 1 && os.Args[1] == "env" { return }` as the first
     statement of the stub's `main`;
   - `internal/verifyselect/full_contract_run_test.go` and both `goStub`
     literals of `internal/verifyselect/contracts_repository_test.go`: the
     literal starts `"#!/bin/sh\n[ \"$1\" = env ] && exit 0\n` before the
     `printf` line;
   - `internal/verifyselect/verify_docs_makefile_test.go`: the stub written
     for `go` and `roundfix-stub` starts the same way;
   - `internal/verifyselect/skills_bump_makefile_test.go`: the `bin/go` stub
     gets the line `if [ "$1" = env ]; then exit 0; fi` right before its
     `printf 'go %s\n'` line.
   `internal/verifyselect/build_cache_makefile_test.go` stays unchanged: its
   stub already exits 2 with no output for anything but `test`.
5. MUST, as the last step of the Task, after every `Makefile` edit, regenerate
   the setup manifest with
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`
   and confirm that a second run reports `File changes: 0`; never hand-edit
   `docs/agents/setup-context.json`. The command changes the two `Makefile`
   `declarationDigest` values and, when the manifest has no `roundfixVersion`
   yet, adds it (ADR-0392); keep what it writes.
6. MUST NOT change any other `Makefile` line or target, `.roundfixrc.yml`, any
   workflow, `go.mod`, `go.sum`, Go production code, any guide, skill,
   `CONTEXT.md`, `README.md`, any ADR, or any other test. No test reaches the
   network, reads the host's Go environment file or real caches, creates
   anything outside its temporary directories, or runs the real
   `cmd/gocache-bound`; the test file compiles for Windows.

## Subtasks

- [ ] Replace the `GOCACHE` block of the `Makefile`.
- [ ] Add `TestMakeFollowsTheGoBuildCacheSetWithGoEnv` and make the five existing stubs answer `env` queries.
- [ ] Regenerate the setup manifest last.

## Acceptance Criteria

- [ ] A usable cache set with `go env` is exported and bounded at 10240 by default; an unusable one is never created, and `make` warns once and falls back to the temporary cache at 2048.
- [ ] With nothing set or with `GOCACHE` exported, `make` exports and bounds exactly what it does on the Delivery Base.
- [ ] Every test of `./internal/verifyselect` passes, the setup manifest is current, and the package vets for Windows.

## Context

- interface: `Makefile`
- interface: `docs/agents/setup-context.json`
- interface: `internal/verifyselect/build_cache_bound_makefile_test.go`
- interface: `internal/verifyselect/full_contract_run_test.go`
- interface: `internal/verifyselect/contracts_repository_test.go`
- interface: `internal/verifyselect/verify_docs_makefile_test.go`
- interface: `internal/verifyselect/skills_bump_makefile_test.go`
- creates: `internal/verifyselect/build_cache_go_env_makefile_test.go`
- instruction: `internal/verifyselect/build_cache_makefile_test.go`
- instruction: `internal/verifyselect/verifyselect_test.go`
- instruction: `docs/adr/0394-the-makefile-follows-the-go-build-cache-set-with-go-env-and-falls-back-small.md`
- instruction: `docs/adr/0370-the-makefile-bounds-its-build-cache-and-the-queue-warns-of-low-disk-and-names-a-full-one.md`
- instruction: `docs/agents/agent-instructions.md`
- instruction: `docs/agents/go.md`

## Consumers

- The `GOCACHE` block of the `Makefile`, which every `make` target reads — consumers: `./internal/verifyselect`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^TestMakeFollowsTheGoBuildCacheSetWithGoEnv$' ./internal/verifyselect 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestMakeFollowsTheGoBuildCacheSetWithGoEnv TestMakeFollowsTheGoBuildCacheSetWithGoEnv/set_and_writable TestMakeFollowsTheGoBuildCacheSetWithGoEnv/set_and_not_mounted TestMakeFollowsTheGoBuildCacheSetWithGoEnv/set_and_read-only TestMakeFollowsTheGoBuildCacheSetWithGoEnv/set_to_a_file TestMakeFollowsTheGoBuildCacheSetWithGoEnv/set_with_a_bound TestMakeFollowsTheGoBuildCacheSetWithGoEnv/environment_wins TestMakeFollowsTheGoBuildCacheSetWithGoEnv/nothing_set; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the test does not exist, so its pass lines are missing; afterwards the test and its seven subtests pass.
- `grep -q 'env -changed GOCACHE' Makefile || { echo 'the Makefile does not ask go env for GOCACHE' >&2; exit 1; }; go test -count=1 -timeout 20m ./internal/verifyselect && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^TestTheSetupManifestIsCurrent$' && GOOS=windows go vet ./internal/verifyselect` — expected: exit 0. Before this Task the `Makefile` asks nothing of `go env`, so the first check fails; afterwards every `Makefile` test of the package passes with the changed stubs, the setup manifest is current after the `Makefile` edit, and the package vets for Windows.

## References

- `_prd.md` → Goals; User Story 1; User Story 2; User Story 3; Core Feature 1; Core Feature 2; Core Feature 3; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; Invariants; Declared breaks; Testing Approach; Build Order 1
- ADR-0394; ADR-0370; ADR-0327; ADR-0392

## Result

Skills loaded: implement-task, coding-guidelines, golang-testing, testing-boss, no-workarounds, tech-writer, setup-context-driven, roundfix, evidence-gate.

Replaced only the authorized `GOCACHE` block with the exact authored bytes.
The surrounding Makefile bytes match `HEAD`. Added the seven parallel,
isolated Go-environment cases and the six `env` guards in the five declared
existing test files; their assertions and remaining content stay unchanged.
The pre-existing task-file change was the Daemon's `pending` to `in_progress`
transition, which this turn preserved.

Acceptance evidence:

- Usable and unusable configured caches: all seven new subtests passed in
  the focused run below. They assert exact ordered invocations and stderr,
  including 10240/2048/4096 bounds, paths with spaces, the absent mount
  staying absent, and immediate probe cleanup. Read-only and regular-file
  cache cases both exercised the warning and fallback.
- Unchanged default and explicit-cache behavior: `nothing set` and
  `environment wins` passed, together with the unchanged
  `TestMakeSharesOneGoBuildCache` and the existing bound cases (including
  explicit cache, sub-make, idle override and disabled bound).
- Package compatibility and manifest: the focused affected-test selection
  passed, and the package compiled for Windows. The sanctioned manifest
  regeneration ran twice after the final Makefile edit: first exit 0,
  `File changes: 1`; second exit 0, `File changes: 0`,
  `Idempotence: verified`. Only the two Makefile declaration digests and
  `roundfixVersion: 0.142.1` changed in the generated manifest. Full-package
  tests, the docscontract manifest check and Windows vet remain for Daemon
  Verification; no terminal verdict is claimed.

Focused commands and outcomes:

- Before the Makefile edit,
  `go test -count=1 -timeout 5m -run '^TestMakeFollowsTheGoBuildCacheSetWithGoEnv$/^(set_and_writable|nothing_set|environment_wins)$' ./internal/verifyselect`
  exited 1: `set_and_writable` recorded the temporary cache instead of the
  configured writable cache, demonstrating the original regression.
- After the edit,
  `go test -count=1 -timeout 5m -v -run '^(TestMakeFollowsTheGoBuildCacheSetWithGoEnv|TestMakeBoundsItsOwnGoBuildCache|TestMakeSharesOneGoBuildCache|TestVerifyContractsRunsEveryContract|TestVerifyChangedRunsTheSelectedContracts|TestVerifyDocsRunsTheDocumentationContractsOnceUnlessSkipped|TestTheSkillsBumpPinsTheCatalogMainCommit|TestTheSkillsBumpRefusesWhatItMustNotPin)$' ./internal/verifyselect`
  exited 0: all eight selected top-level tests and their subtests passed.
- `GOOS=windows GOPROXY=off GOSUMDB=off go test -c -o "$task_compile_dir/verifyselect.test.exe" ./internal/verifyselect`
  exited 0; `task_compile_dir` was a fresh `mktemp -d` directory.
- `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`
  exited 0 on both runs. Both reported the same two nested-carrier warnings
  for existing Baseline fixture AGENTS.md files; those files were unchanged.
- Exact-byte comparison against the authored block and against the Makefile
  prefix/suffix in `HEAD` passed. `gofmt -l` on the new test emitted nothing;
  `git diff --check` exited 0.

No declared Verification command, repository gate, commit, push, pull request,
Task status edit, Task Graph edit, or other Task edit was performed. No
follow-up implementation was added to this slice.
