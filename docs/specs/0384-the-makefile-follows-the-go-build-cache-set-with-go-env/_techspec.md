---
spec: 0384-the-makefile-follows-the-go-build-cache-set-with-go-env
prd: _prd.md
created: 2026-11-12
---

# The Makefile follows the Go build cache set with go env — Technical Spec

## Executive Summary

The `GOCACHE` block of the `Makefile` gains one step before its current
choice. When the environment does not define `GOCACHE`, a `$(shell)` asks
`$(GO) env -changed GOCACHE`; when that prints a value, it reads the path
with `$(GO) env GOCACHE` and keeps it only if it is absolute, an existing
directory, and writable through a probe file it removes at once. A kept path
becomes `GOCACHE` with the 10,240 MiB default bound. Otherwise the block
chooses as before (`$TMPDIR/roundfix-go-build-<uid>` or `.gocache`), and a
second `$(shell)` prints one warning line and sets the default bound to
2,048 MiB when a cache was set but unusable, or keeps 10,240 MiB when nothing
was set. The trade-off accepted is two to four `go env` calls, a few
milliseconds each, at every top-level `make`, and six `go` stubs in existing
tests that must now answer `env` queries (ADR-0394).

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, record or schema
  is minted or changed. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — nothing reaches the network.
  Source: `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0394 (new) records the decision
  and amends ADR-0370, whose Build Cache Bound now holds the cache set with
  `go env` too. ADR-0327 applies unchanged: its trim removes only entries
  idle for at least two hours, oldest first, which is what keeps a cache
  shared with other Go projects safe. ADR-0392 applies unchanged to the manifest regeneration.
  ADR-0334 does not apply: Run-start and gate trims read
  `worktree.shared_caches`, which is unchanged. Source:
  `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-11-12 ("Autorizo; sem SSD, TMPDIR com 2 GB") for the `Makefile`, with
  the derived `docs/agents/setup-context.json` through its sanctioned command
  and the repository-owned guide `docs/agents/specific-repository.md`.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0384-the-makefile-follows-the-go-build-cache-set-with-go-env/_authorization.md`;
  bounded files: `Makefile`, `docs/agents/setup-context.json`,
  `docs/agents/specific-repository.md`.

## Current behavior

Characterized at `d79b30b4` with a printing makefile passed beside the real
one (`make -f Makefile -f print.mk print-gocache`) in temporary directories:

| environment | `GOCACHE` | `GOCACHE_BOUND_MB` |
| --- | --- | --- |
| `GOENV` file sets an existing directory | `$TMPDIR/roundfix-go-build-<uid>` | 10240 |
| nothing set | `$TMPDIR/roundfix-go-build-<uid>` | 10240 |
| `GOCACHE` exported | the exported path | 0 |

The `gocache-bound` target runs `$(GO) run -buildvcs=false ./cmd/gocache-bound
-dir "$(GOCACHE)" -max-mb "$(GOCACHE_BOUND_MB)" -idle "$(GOCACHE_BOUND_IDLE)"`
at `MAKELEVEL` 0 when the bound is not 0, and warns without failing when it
fails. Sub-makes inherit the exported `GOCACHE`. `go env -changed GOCACHE`
prints `GOCACHE='<path>'` only for a value set in the environment file (or
the environment); a relative value reads as `off`.

## System Architecture

- `Makefile` (task_01): the `ifeq ($(origin GOCACHE),undefined)` block only.
  The `GOCACHE_BOUND_MB ?= 0`, `GOCACHE_BOUND_IDLE ?= 2h` and `export GOCACHE`
  lines after it, and every target, keep their bytes.
- `internal/verifyselect/build_cache_go_env_makefile_test.go` (new, task_01):
  the new behavior through the real `Makefile` with a compiled `go` stub that
  forwards `env` to the real go command.
- Five existing test files of `internal/verifyselect` (task_01): their `go`
  stubs answer `env` with no output and no record (Declared breaks).
- `docs/agents/setup-context.json` (task_01): regenerated, because it records
  the `Makefile` digest.
- `docs/agents/specific-repository.md`, `README.md`, `CONTEXT.md` and
  ADR-0370 (task_02): the rule, the README paragraph, the **Build Cache
  Bound** entry and the amendment notice.

## Implementation Design

### Interfaces

```make
ifeq ($(origin GOCACHE),undefined)
GOCACHE_FROM_GO_ENV := $(shell ...)   # the set cache when usable, else empty
ifneq ($(GOCACHE_FROM_GO_ENV),)
GOCACHE := $(GOCACHE_FROM_GO_ENV)
GOCACHE_BOUND_MB ?= 10240
else
GOCACHE_FALLBACK := $(shell ...)      # unchanged temporary or .gocache choice
GOCACHE_DEFAULT_BOUND_MB := $(shell ...)  # 10240, or 2048 plus the warning
GOCACHE := $(GOCACHE_FALLBACK)
GOCACHE_BOUND_MB ?= $(GOCACHE_DEFAULT_BOUND_MB)
endif
endif
```

### Data Models

None. The only datum is the warning line of API Contract 2.

### API Contracts

1. API Contract: set and usable. With `GOCACHE` neither in the environment nor
   on the command line, when `$(GO) env -changed GOCACHE` prints a non-empty
   line and `$(GO) env GOCACHE` prints a path that starts with `/`, is a
   directory, and accepts `mktemp "<path>/.write-test.XXXXXX"` followed by
   its removal, `make` exports that path as `GOCACHE` and `GOCACHE_BOUND_MB`
   defaults to 10240. The probe file is gone afterwards. `make` never runs
   `mkdir` on that path or a parent.
2. API Contract: set and unusable. When `$(GO) env -changed GOCACHE` prints a
   value but the path fails any of those checks, `make` exports the cache the
   `d79b30b4` block chose, `GOCACHE_BOUND_MB` defaults to 2048, and `make`
   prints exactly one line to standard error, once per top-level `make`:
   `gocache: warning: GOCACHE set with go env (<path go env prints>) is not a writable directory; using <fallback> with a 2048 MiB default bound`.
3. API Contract: unchanged paths. With nothing set (empty output from
   `go env -changed GOCACHE`, including `GOENV=off`, a failing `go`, or a go
   without `-changed`), `GOCACHE` and the 10240 default are those of
   `d79b30b4` and nothing is printed. A `GOCACHE` in the environment or on the
   command line skips the whole block, as before. `GOCACHE_BOUND_MB` and
   `GOCACHE_BOUND_IDLE` from the caller still win.

### Surface Transcripts

None. The one new output is the single warning line API Contract 2 names
verbatim; the new test asserts the whole standard error of each case.

### Invariants

- `make` writes nothing at a set path except its probe file, and never at a
  missing one.
- No `$(shell)` of the block runs while `GOCACHE` is a defined make variable,
  so a make that exports variables to `$(shell)` (GNU make 4.4) cannot make
  `go env -changed` report the fallback as set.
- With nothing set, the exported `GOCACHE` and the bound equal those of
  `d79b30b4` on every machine.

### Declared breaks

- `TestMakeBoundsItsOwnGoBuildCache`
  (`internal/verifyselect/build_cache_bound_makefile_test.go`),
  `TestVerifyContractsRunsEveryContract`
  (`internal/verifyselect/full_contract_run_test.go`),
  `TestVerifyChangedRunsTheSelectedContracts`
  (`internal/verifyselect/contracts_repository_test.go`, two stubs),
  `TestVerifyDocsRunsTheDocumentationContractsOnceUnlessSkipped`
  (`internal/verifyselect/verify_docs_makefile_test.go`) and
  `TestTheSkillsBumpPinsTheCatalogMainCommit` with
  `TestTheSkillsBumpRefusesWhatItMustNotPin`
  (`internal/verifyselect/skills_bump_makefile_test.go`) record every
  invocation of their `go` stub and compare the list. The new block adds
  `env -changed GOCACHE` calls, so each stub returns at once, with no output
  and no record, when its first argument is `env`; empty output is "nothing
  set", so every expected list and the rest of each test keep their bytes.
  `TestMakeSharesOneGoBuildCache` needs no change: its stub exits 2 with no
  output for anything but `test`.
- `docs/agents/setup-context.json`: the two `declarationDigest` values of the
  `Makefile` change, and the regeneration adds the Recorded Roundfix Version
  (`roundfixVersion`) when the delivery base does not hold it yet (ADR-0392);
  from then on a Roundfix older than that release refuses this repository.

## Coverage Map

- Goal 1, User Story 1, Core Feature 1, Success Metric 1 → API Contract 1 and
  the `set and writable` and `set with a bound` cases; Build Order 1.
- Goal 2, User Story 2, Core Feature 2, Success Metric 2 → API Contract 2,
  the Invariants and the `set and not mounted`, `set and read-only` and
  `set to a file` cases; Build Order 1.
- Goal 3, User Story 3, Core Feature 3, Success Metrics 3 and 4 → API
  Contract 3, the `nothing set` and `environment wins` cases and the Declared
  breaks; Build Order 1.
- Core Feature 4 → task_02's documentation and ADR-0394; Build Order 2.

Deviations from the PRD: None.

## Absence Evidence

1. A `go env` query in the `Makefile` — command: `grep -n 'go env\|env -changed\|GOENV' Makefile`; result: no match at `d79b30b4`.
2. A test that sets `GOENV` for `make` — command: `grep -rln 'GOENV' internal/verifyselect`; result: no match at `d79b30b4`.
3. Production code that sets `GOCACHE` or `GOENV` for a Verification or an Agent — command: `grep -rn 'GOENV\|"GOCACHE' --include='*.go' internal cmd | grep -v _test.go`; result: only the Linux Lane's container argument in `internal/cli/deliver_linux_lane.go` at `d79b30b4`.

## ADR Alternatives

- ADR-0394: new; it rejects naming the external disk in the Makefile,
  failing make, exporting `GOCACHE` from shell startup files, creating the
  set path, reading the environment file directly and leaving the set cache
  unbounded, and defers a `worktree.shared_caches` placeholder; this design
  takes none of them.
- ADR-0370: amended; the bound, its idle rule, the explicit-`GOCACHE` rule,
  the Free-Space Floor and the Disk-Full Park stand; its rejected
  `go clean -cache`, young-entry removal, `-trimpath`, shell trim and refusal
  below the floor are not taken; its option "Moving the cache to another
  volume" is taken for the cache set with `go env`.
- ADR-0327: kept; the bound still removes only idle entries through its
  rules, and none of its rejected options is taken.
- ADR-0392: kept; the regeneration writes what the manifest command decides.
- ADR-0334: not applicable; the Run-start and gate trims are unchanged.

## Integration Points

The go command only: `go env -changed GOCACHE` (Go 1.23 or later; go.mod
requires 1.26) and `go env GOCACHE`, run through `$(GO)` so a caller's `GO`
is respected. GNU make 3.81 (macOS) and 4.x (Linux, CI): the block uses only
`ifeq`/`ifneq`/`else`, `:=`, `?=` and `$(shell)`, and balances its
parentheses inside `$(shell)` (`case` patterns `(/*)` and `(*)`), because an
unbalanced `)` or a `#` ends a make function or line.

## Testing Approach

New file `internal/verifyselect/build_cache_go_env_makefile_test.go`, top-level
test and every subtest calling `t.Parallel()` first:
`TestMakeFollowsTheGoBuildCacheSetWithGoEnv` runs
`make --no-print-directory -f <Makefile> test GO=<stub> GO_TEST_PARALLEL=16`
in a `t.TempDir()`, with the caller's environment minus `GOCACHE`, `GOENV`,
`TMPDIR`, the make variables and the bound variables, plus `TMPDIR`, `GOENV`
(a temporary file), the real go's path and the record path. The stub forwards
`env` to the real go command and records every other invocation with the
`GOCACHE` it received. Subtests and what they assert (the invocation list,
the exact standard error, and file-system effects):

- `set and writable`: the set directory (its path holds a space) is used and
  bounded at 10240, no warning, no probe left;
- `set and not mounted`: the set path lies under a missing directory below an
  existing `Volumes` directory; fallback, 2048, the warning, and the missing
  directory still absent;
- `set and read-only`: a `0555` directory; fallback, 2048, the warning;
- `set to a file`: a regular file at the set path; fallback, 2048, the
  warning;
- `set with a bound`: `GOCACHE_BOUND_MB=4096`; the set cache with 4096;
- `environment wins`: `GOCACHE` exported and a usable set cache; the exported
  path, no bound call, no warning;
- `nothing set`: an empty environment file; the temporary cache, 10240, no
  warning.

On `d79b30b4` the test fails: only `nothing set` and `environment wins` pass,
the kept paths. Sabotage proofs, each run on the prototype: the block
reverted fails the five other subtests; `mkdir -p` in place of the directory
check fails `set and not mounted`; a fallback default of 10240 fails the
three unusable cases; the warning removed fails the same three; the set cache
defaulting to bound 0 fails `set and writable`; the probe removed (directory
check only) fails `set and read-only`; the block entered when `GOCACHE` is in
the environment fails `environment wins`. The Declared breaks are proved by
the whole package: before the stub changes, the six tests fail on the new
`env` records; after them they pass.

## Build Order

1. The `Makefile` block, the new test, the five stub changes and, last, the
   manifest regeneration (task_01). One Task, because the tests read the
   `Makefile` and the manifest records its final bytes.
2. The guide rule, README paragraph, **Build Cache Bound** entry and ADR-0370
   notice (task_02, depends on: 1).
3. Final QA gate (depends on: 2).

No step needs another Spec's code. Spec 0383 merges first and may change
`CONTEXT.md` and `docs/agents/setup-context.json`; task_01 regenerates the
manifest with the command, and task_02 edits only its own entry.

## Risks & Considerations

- The bound now trims a cache other Go projects share on the mini; only
  entries idle two hours or more go, oldest first (ADR-0327).
- A sandboxed Agent whose sandbox cannot write the set cache gets the
  fallback and the warning on each top-level `make`; the Run Agents of
  2026-11-12 ran `make` and `go test` escalated, outside the sandbox.
- The regeneration adds `roundfixVersion` to this repository's manifest when
  the delivery base lacks it; an operator binary older than that release
  then refuses this repository (ADR-0392). The operator installs every
  release right after its merge.
- The real mini is proved only after the merge, by the PRD's one
  Unreachable Acceptance entry.

## Glossary

- changes: **Build Cache Bound**

## Decisions

- `make` asks the go command (`go env -changed GOCACHE`, then
  `go env GOCACHE`) and never names a machine's path.
- Usable means absolute, an existing directory and writable by probe; `make`
  never creates the set path.
- Default bounds 10,240 MiB for the set cache, 2,048 MiB for the fallback
  with one warning; an exported `GOCACHE` unchanged.
- ADR-0394, amending ADR-0370; two implementation Tasks and the QA gate.
