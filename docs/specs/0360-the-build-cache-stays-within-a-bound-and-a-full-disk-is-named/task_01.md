---
task: task_01
spec: 0360-the-build-cache-stays-within-a-bound-and-a-full-disk-is-named
status: pending
type: infra
complexity: medium
---

# Task 01: The Makefile holds its Go build cache to the Build Cache Bound

## Overview

Every top-level `make` target that builds or tests first runs a new
`gocache-bound` command over the Go build cache the Makefile chose, which
removes, above 10,240 MiB, only entries no build used for two hours, oldest
first, with the Shared Build Cache rules of `internal/sharedcache`. The bound
never repeats inside a sub-make, never touches a caller's own `GOCACHE`
unless asked, and never fails a build. The Project Config declaration of the
same cache takes the same values, the README and the repository rules guide
describe the bound, `CONTEXT.md` defines it, ADR-0327 carries its amendment
notice, and the setup manifest is regenerated after the Makefile change.

## Requirements

1. MUST add `internal/sharedcache/bound.go` with `BoundBudget` and
   `BoundMain(args []string, stderr io.Writer, now func() time.Time) int`
   exactly as `_techspec.md` → Interfaces and API Contract 2 describe: flags
   `-dir`, `-max-mb` and `-idle`, the four standard-error lines word for
   word, exit `2` for a usage error and `0` otherwise, and no second trim
   implementation: it plans and applies only through `MeasureContext` and
   `ApplyContext` under a `BoundBudget` context.
2. MUST add `cmd/gocache-bound/main.go`, a thin `main` that exits with
   `sharedcache.BoundMain(os.Args[1:], os.Stderr, time.Now)`, with a package
   comment naming the Build Cache Bound.
3. MUST change the Makefile as API Contract 1 says, and nothing else in it:
   `GOCACHE_BOUND_MB ?= 10240` inside the `ifeq ($(origin GOCACHE),undefined)`
   block after the cache is chosen; `GOCACHE_BOUND_MB ?= 0` and
   `GOCACHE_BOUND_IDLE ?= 2h` after it with a comment naming ADR-0370; the
   phony `gocache-bound` target with its help text, its `$(MAKELEVEL)` and
   `GOCACHE_BOUND_MB` conditions, the `$(GO) run $(RUN_FLAGS)
   ./cmd/gocache-bound -dir "$(GOCACHE)" -max-mb "$(GOCACHE_BOUND_MB)" -idle
   "$(GOCACHE_BOUND_IDLE)"` invocation and the fallback that prints
   `gocache-bound: warning: the Build Cache Bound did not run (exit <n>); the build continues`
   to standard error and succeeds; and `gocache-bound` as the first
   prerequisite of each target API Contract 1 lists.
4. MUST change the `{tmpdir}/roundfix-go-build-{uid}` entry of
   `.roundfixrc.yml` to `max_size_mb: 10240` and `idle: 2h`, with its comment
   saying the Makefile holds the same cache to the same bound, as API
   Contract 3 says.
5. MUST add `internal/sharedcache/bound_test.go` with
   `TestBoundMainTrimsOnlyIdleEntriesAboveTheBound` and
   `TestBoundMainRefusesAnInvalidRequest`, and
   `internal/verifyselect/build_cache_bound_makefile_test.go` with
   `TestMakeBoundsItsOwnGoBuildCache`, with the rows `_techspec.md` →
   Testing Approach names, each asserting the exact recorded invocations or
   the exact lines. Every new top-level test calls `t.Parallel()` first.
   `TestMakeBoundsItsOwnGoBuildCache` uses a Go stub built with
   `testfixture.FixtureBinary`, as `TestMakeSharesOneGoBuildCache` does, and
   removes `GOCACHE`, `TMPDIR`, `MAKEFLAGS`, `MFLAGS`, `MAKEOVERRIDES`,
   `MAKELEVEL`, `GOCACHE_BOUND_MB` and `GOCACHE_BOUND_IDLE` from the child's
   environment before setting its own.
6. MUST change `TestVerifyDocsRunsTheDocumentationContractsOnceUnlessSkipped`
   (`internal/verifyselect/verify_docs_makefile_test.go`) to count the
   bound's `run -buildvcs=false ./cmd/gocache-bound ` invocation and require
   exactly one, as `_techspec.md` → Declared breaks says, and update every
   other test that pins the Makefile's recorded Go invocations on the default
   branch when the Task runs.
7. MUST write, word for word as `_techspec.md` → API Contract 8 gives them,
   the README's Development paragraph and the sentence added to the
   `docs/agents/specific-repository.md` rule starting "`make` shares
   `roundfix-go-build-<uid>`"; add through `domain-modeling` the `CONTEXT.md`
   entry **Build Cache Bound** with its text and `_Avoid_` line, beside
   **Shared Build Cache**, with `ADR-0370` written as the glossary's Markdown
   link to
   `docs/adr/0370-the-makefile-bounds-its-build-cache-and-the-queue-warns-of-low-disk-and-names-a-full-one.md`,
   adding no other capitalized phrase to `CONTEXT.md`; and write the ADR-0327
   amendment notice after the Decision paragraph API Contract 8 names, with
   `ADR-0370` in the bold opening as a Markdown link to
   `0370-the-makefile-bounds-its-build-cache-and-the-queue-warns-of-low-disk-and-names-a-full-one.md`,
   and `updated_at: 2026-11-06T12:00:00Z`; no other line of ADR-0327 changes.
8. MUST NOT run `go clean -cache`, and no test or step may read, trim or fill
   a real build cache: tests use `t.TempDir()` caches and the Go stub.
9. MUST regenerate `docs/agents/setup-context.json`, whose recorded Makefile
   digest the change makes stale, with
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`
   as the Task's last step, after the last Makefile edit, and confirm that a
   second run reports no change; never hand-edit it. The Task may change only
   its declared files, and of the Governed Paths only `Makefile`,
   `.roundfixrc.yml`, `docs/agents/setup-context.json` and
   `docs/agents/specific-repository.md`.

## Subtasks

- [ ] Add `BoundMain` and its tests.
- [ ] Add `cmd/gocache-bound`.
- [ ] Add the make variables, the target and its prerequisites.
- [ ] Align the Project Config declaration.
- [ ] Add the Makefile test and adjust the verify-docs count.
- [ ] Write the README and rules text, the glossary entry and the ADR-0327 notice.
- [ ] Regenerate the setup manifest last.

## Acceptance Criteria

- [ ] Over a cache above its bound, `gocache-bound` removes only entries idle for at least two hours, oldest first, until within the bound, and warns when it cannot reach it.
- [ ] Each top-level make that builds or tests runs the bound once over the cache the Makefile chose; a sub-make, an explicit `GOCACHE` and `GOCACHE_BOUND_MB=0` do not.
- [ ] A failing bound prints its warning and `make` still exits `0`.
- [ ] The declaration says 10,240 MiB and two hours, and the setup manifest is current.
- [ ] The README, the repository rules guide, `CONTEXT.md` and ADR-0327 describe the bound.

## Context

- instruction: `docs/adr/0370-the-makefile-bounds-its-build-cache-and-the-queue-warns-of-low-disk-and-names-a-full-one.md`
- instruction: `docs/adr/0327-gc-trims-a-declared-shared-build-cache.md`
- instruction: `internal/sharedcache/sharedcache.go`
- instruction: `internal/verifyselect/build_cache_makefile_test.go`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- creates: `internal/sharedcache/bound.go`
- creates: `internal/sharedcache/bound_test.go`
- creates: `cmd/gocache-bound/main.go`
- creates: `internal/verifyselect/build_cache_bound_makefile_test.go`
- interface: `Makefile`
- interface: `.roundfixrc.yml`
- interface: `internal/verifyselect/verify_docs_makefile_test.go`
- interface: `README.md`
- interface: `docs/agents/specific-repository.md`
- interface: `CONTEXT.md`
- interface: `docs/adr/0327-gc-trims-a-declared-shared-build-cache.md`
- interface: `docs/agents/setup-context.json`

## Consumers

- The Makefile's recorded Go invocations — consumers: `./internal/verifyselect`
- `sharedcache.BoundMain` — consumers: `./internal/sharedcache`
- The setup manifest's Makefile digest and the glossary — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 10m -v -run '^(TestBoundMainTrimsOnlyIdleEntriesAboveTheBound|TestBoundMainRefusesAnInvalidRequest)$' ./internal/sharedcache 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestBoundMainTrimsOnlyIdleEntriesAboveTheBound TestBoundMainRefusesAnInvalidRequest 'TestBoundMainTrimsOnlyIdleEntriesAboveTheBound/within_its_bound' 'TestBoundMainTrimsOnlyIdleEntriesAboveTheBound/above_its_bound_with_no_idle_entry'; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 10m ./internal/sharedcache && go build -buildvcs=false -o /dev/null ./cmd/gocache-bound && GOOS=windows go build -buildvcs=false -o /dev/null ./cmd/gocache-bound` — expected: exit 0. Before this Task neither test exists, so the first pass line is missing; afterwards the bound trims only idle entries above its bound, refuses an invalid request, and the command builds for the host and for Windows.

- `out="$(go test -count=1 -timeout 10m -v -run '^(TestMakeBoundsItsOwnGoBuildCache|TestVerifyDocsRunsTheDocumentationContractsOnceUnlessSkipped|TestMakeSharesOneGoBuildCache)$' ./internal/verifyselect 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestMakeBoundsItsOwnGoBuildCache 'TestMakeBoundsItsOwnGoBuildCache/default' 'TestMakeBoundsItsOwnGoBuildCache/once_per_invocation' 'TestMakeBoundsItsOwnGoBuildCache/sub-make' 'TestMakeBoundsItsOwnGoBuildCache/explicit_cache' 'TestMakeBoundsItsOwnGoBuildCache/failing_bound' TestVerifyDocsRunsTheDocumentationContractsOnceUnlessSkipped TestMakeSharesOneGoBuildCache; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./internal/verifyselect` — expected: exit 0. Before this Task the Makefile test does not exist, so its pass line is missing; afterwards the bound runs once per top-level make over the cache the Makefile chose, skips a sub-make and an explicit cache, fails open, and the existing Makefile tests pass.

- `grep -q '^      max_size_mb: 10240$' .roundfixrc.yml && grep -q '^      idle: 2h$' .roundfixrc.yml && ! grep -q 'max_size_mb: 20480' .roundfixrc.yml && grep -q '^gocache-bound: ##' Makefile || { echo 'declaration or target missing' >&2; exit 1; }; tr -s '[:space:]' ' ' < README.md | grep -q -- 'holds the Makefile.s Go build cache to its Build Cache Bound' || { echo 'missing README bound' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/agents/specific-repository.md | grep -q -- 'holds the cache it chose to the Build Cache Bound' || { echo 'missing repository rule' >&2; exit 1; }; grep -F -- '**Build Cache Bound**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Build Cache Bound\*\*(:| —( |$))' || { echo 'missing term: Build Cache Bound' >&2; exit 1; }; f=docs/adr/0327-gc-trims-a-declared-shared-build-cache.md; grep -q '^updated_at: 2026-11-06T12:00:00Z$' "$f" && tr -s '[:space:]' ' ' < "$f" | grep -qF -- '**Amended by [ADR-0370](0370-the-makefile-bounds-its-build-cache-and-the-queue-warns-of-low-disk-and-names-a-full-one.md) (2026-11-06).**' || { echo "missing amendment: $f" >&2; exit 1; }; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the declaration says 20480 and the target is absent, so the first checks fail; afterwards the declaration matches the bound, the README, the rules guide, the glossary and ADR-0327 describe it, and the documentation contracts, the regenerated setup manifest's among them, pass.

## References

- `_prd.md` → Goals 1 and 2; Core Feature 1; Core Feature 5; Success Metric 1; Success Metric 2; Success Metric 6; Glossary
- `_techspec.md` → Current behavior 1 and 2; Interfaces; API Contracts 1-3; API Contract 8; Invariants 1-3; Declared breaks; Testing Approach; Build Order 1
- ADR-0370; ADR-0327
- Source: [`references/2026-11-06-the-repository-build-cache-stays-within-a-size-bound.md`](references/2026-11-06-the-repository-build-cache-stays-within-a-size-bound.md) (2026-11-06) and [`references/2026-11-05-the-roundfix-build-cache-fills-the-mac-mini-internal-disk.md`](references/2026-11-05-the-roundfix-build-cache-fills-the-mac-mini-internal-disk.md) (2026-11-05)
