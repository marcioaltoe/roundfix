---
task: task_01
spec: 0378-the-catalog-skill-bump-refreshes-with-a-binary-built-after-the-digests
status: completed
type: infra
complexity: low
---

# Task 01: The Catalog Skill Bump refreshes with a binary built after the digests

## Overview

On 2026-11-10 the operator recorded that `make skills-third-party-bump`
leaves a stale catalog digest in the Setup Manifest
([`references/2026-11-10-the-catalog-skill-bump-leaves-a-stale-catalog-digest-in-the-setup-manifest.md`](references/2026-11-10-the-catalog-skill-bump-leaves-a-stale-catalog-digest-in-the-setup-manifest.md),
recorded by the operator): the recipe runs its Managed Refresh with a binary
it built before its own `make baseline-digests` step rewrote embedded
Baseline assets. This Task builds the binary again after that step, proves
the order in the Catalog Skill Bump Makefile contract, and regenerates the
Setup Manifest that records the `Makefile` digest.

## Requirements

1. MUST, in the `skills-third-party-bump` recipe of the `Makefile`, add the
   line `$(GO) build -o "$$tmp/roundfix" ./cmd/roundfix; \` immediately
   after the line that runs `baseline-digests` and before the first
   `roundfix baseline update` line, as `_techspec.md` → Interfaces shows.
   Nothing else in the `Makefile` changes, and every refusal, diagnostic and
   the final line of the bump keep their bytes.
2. MUST, in `internal/verifyselect/skills_bump_makefile_test.go`
   `TestTheSkillsBumpPinsTheCatalogMainCommit`, insert `build` in the `want`
   call list between
   `"go test ./internal/baseline -run TestCatalogCompatibility -count=1"` and
   the first `"roundfix baseline update --repo . --no-skills --no-history --yes --format text"`.
   This is the declared break of the pinned call sequence (API Contract 1).
3. MUST, in `newSkillsBumpFixture`, extend the stubs as `_techspec.md` →
   Interfaces describes:
   - the `go` stub, when its first argument is `build`, after copying the
     `roundfix` stub to `$3`, copies the file `embedded-asset` of its working
     directory to `$3.embedded`, or writes `absent` and a newline there when
     that file does not exist;
   - the `go` stub, when its first two arguments are `test probe`, writes
     `regenerated` and a newline to `embedded-asset` in its working
     directory;
   - the `roundfix` stub's `baseline update` branch copies `"$0.embedded"`
     to `refreshed-from` in its working directory, then prints
     `File changes: 0` as today.
4. MUST add `TestTheSkillsBumpRefreshesWithABinaryBuiltAfterTheDigests` to
   the same file, calling `t.Parallel()` first: build the fixture, write
   `repo/embedded-asset` with `before digests` and a newline, run
   `skillsBumpCommand` and require exit 0, then read `repo/refreshed-from`
   and fail unless it is `regenerated` and a newline, with the message
   `the managed refresh ran a binary that embeds %q; want the asset make baseline-digests regenerated, %q`
   (the read bytes, then `"regenerated\n"`); end with
   `assertSkillsBumpTemporaryDirectoryRemoved`. The other three bump tests
   keep their assertions.
5. MUST, as the Task's last step, regenerate `docs/agents/setup-context.json`
   with
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`,
   run twice, the second reporting `File changes: 0`; never hand-edit a
   digest. Only the `declarationDigest` values of the two verification
   entries are expected to change.
6. MUST keep every test hermetic: the fixture's local catalog, `PATH`,
   `HOME` and `TMPDIR` as today; no network, no real catalog, no real bump.

## Subtasks

- [ ] Add the build after the digests step in the bump recipe.
- [ ] Extend the recording stubs and the pinned call sequence.
- [ ] Add the test that proves the refresh embeds what the digests step wrote.
- [ ] Regenerate the Setup Manifest as the last step.

## Acceptance Criteria

- [ ] `TestTheSkillsBumpRefreshesWithABinaryBuiltAfterTheDigests` and `TestTheSkillsBumpPinsTheCatalogMainCommit` pass, and both fail when the new build is dropped or moved before the digests step.
- [ ] `TestTheSetupManifestIsCurrent` passes after the `Makefile` edit.
- [ ] The `Makefile` diff is the one recipe line; the Setup Manifest diff is the two `declarationDigest` values.

## Context

- interface: `Makefile`
- interface: `internal/verifyselect/skills_bump_makefile_test.go`
- interface: `docs/agents/setup-context.json`
- instruction: `docs/adr/0360-the-vendored-catalog-skills-follow-the-catalog-through-one-pinned-commit.md`

## Consumers

- The `skills-third-party-bump` recipe — consumers: `./internal/verifyselect`
- The Setup Manifest's `Makefile` digest — consumers: `./internal/docscontract`

## Verification

- `out="$(GOMAXPROCS=2 go test -count=1 -timeout 10m ./internal/verifyselect -run '^TestTheSkillsBump(PinsTheCatalogMainCommit|RefreshesWithABinaryBuiltAfterTheDigests)$' -v 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestTheSkillsBumpRefreshesWithABinaryBuiltAfterTheDigests'` — expected: both tests pass; fails before the work because the new test does not exist.
- `GOMAXPROCS=2 go test -count=1 -timeout 10m -tags docscontract ./internal/docscontract -run '^TestTheSetupManifestIsCurrent$' && awk 'index($0, "baseline-digests >&2;") { seen = 1; next } seen && index($0, "build -o") { built = 1 } seen && index($0, "baseline update") { exit } END { exit !built }' Makefile` — expected: the Setup Manifest is current and a build follows the digests step in the recipe; fails before the work because no build follows it.

## References

- `_prd.md` → Goals; User Stories 1-2; Core Features 1-2; Success Metric 1; Success Metric 2
- `_techspec.md` → Interfaces; API Contract 1; Testing Approach; Build Order 1
- ADR-0360

## Result

Skills loaded: implement-task, coding-guidelines, golang-testing, testing-boss, no-workarounds, systematic-debugging, setup-context-driven, roundfix, evidence-gate, tech-writer.

The Catalog Skill Bump now rebuilds its temporary Roundfix binary immediately
following `baseline-digests`, before either Managed Refresh. The recording Go
stub snapshots `embedded-asset` at build time (or `absent\n`), the probe digests
step writes `regenerated\n`, and the refresh stub copies the embedded snapshot
to `refreshed-from`. The pinned call list includes the new build, and the new
parallel test checks the refreshed asset and temporary-directory cleanup.

Focused evidence by acceptance criterion:

1. Bump ordering and freshness: compiled the affected package with
   `go test -c -o <probe-dir>/verifyselect.test ./internal/verifyselect`, then
   ran the binary from `internal/verifyselect` with
   `-test.run '^TestTheSkillsBump' -test.v -test.timeout 3m -test.parallel 2`.
   Before the recipe edit, both changed contract tests failed: the pinned
   sequence lacked the new build, and the refresh embedded `"before digests\n"`
   instead of `"regenerated\n"`. With the recipe edit, all five bump tests and
   all ten refusal subtests passed (exit 0). Separate negative probes dropped
   the added build and moved it before `baseline-digests`: each exited 1 with
   both named tests failing, including the exact stale-asset diagnostic. Each
   probe restored the fixed Makefile bytes in `finally`; all five bump tests
   passed again after restoration. These tests use only their fixture-owned
   catalog, PATH, HOME and TMPDIR; no real bump or network was used.
2. Setup Manifest currency: compiled the focused contract with
   `go test -c -tags docscontract -o <probe-dir>/docscontract.test ./internal/docscontract`
   and ran that binary from `internal/docscontract` with
   `-test.run '^TestTheSetupManifestIsCurrent$' -test.v -test.timeout 3m`.
   After the Makefile edit and before regeneration it failed with
   `Setup Manifest is stale`; after regeneration it passed (exit 0).
3. Exact diff scope: byte comparison against HEAD proves the Makefile differs
   only by the required recipe line. The three other bump-test bodies and
   `skillsBumpCommand` remain byte-identical, preserving their assertions and
   environment isolation. `gofmt -l` reports no files. Structural JSON comparison
   and the text diff show exactly `/verification/0/declarationDigest` and
   `/verification/1/declarationDigest` changed; the catalog digest is unchanged.

As the final implementation step, ran the sanctioned generator twice:
`go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`.
Both runs exited 0. The first reported `File changes: 1`, updating only the
Setup Manifest; the second reported `File changes: 0` and `Idempotence: verified`.
No digest was hand-edited. The subsequent manifest check and diff inspection
were read-only; this Result records their evidence.

The probe binaries live in a task-owned temporary directory. An initial direct
binary invocation from the repository root was refused by suiteguard because
its relative repository path requires the package directory; the focused runs
above used the package directory.

Declared Task Verification and repository Verification were not run; the Daemon
owns them and Task settlement. No Task status, Task Graph, other Task file,
commit, push or pull request was changed by this Agent. No follow-up work was
added to this slice.
