---
spec: 0378-the-catalog-skill-bump-refreshes-with-a-binary-built-after-the-digests
prd: _prd.md
created: 2026-11-10
---

# The Catalog Skill Bump refreshes with a binary built after the digests — Technical Spec

## Executive Summary

The `skills-third-party-bump` recipe in the `Makefile` gains one line,
`$(GO) build -o "$$tmp/roundfix" ./cmd/roundfix`, right after its
`make baseline-digests` step and before its two Managed Refresh runs. The
binary those runs use then embeds every Baseline asset the bump wrote, so the
catalog digest the Setup Manifest records is the digest of the tree the
repository holds. The Catalog Skill Bump Makefile contract in
`internal/verifyselect` proves the order through its recording stubs: the
`go` stub records, at each build, the bytes of a fixture asset the probe
digests step rewrites, and the `roundfix` stub copies what its binary
recorded when it refreshes. The trade-off accepted is one more compile of
`./cmd/roundfix` in a command the operator runs about once a cycle, instead
of a currency check or a printed instruction that would leave the author to
run the refresh by hand.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted or changed. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the network steps of the bump are
  unchanged, and every test uses a local `file://` catalog. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0360 (the Catalog Skill Bump is
  the only writer and runs the asset sync, a build, the restore,
  `make baseline-digests` and the Managed Refresh in that order). ADR-0364
  and ADR-0380 are not applicable: they cite ADR-0360, and their rules are
  unchanged. Source: `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle BB batch 2 answer "Mesmas do BB, teto minor (Recommended)" ("Makefile
  na 0378 só nos arquivos que elas exigirem"); bounded files: `Makefile`,
  `docs/agents/setup-context.json`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0378-the-catalog-skill-bump-refreshes-with-a-binary-built-after-the-digests/_authorization.md`;
  bounded files: `Makefile`, `docs/agents/setup-context.json`.

## System Architecture

The bump recipe at `b233761e`, in order:

1. preflight refusals (lock, folders, profile, catalog, commit);
2. `git fetch --depth=1` of the one catalog commit into `$tmp/catalog`;
3. `go build` of `$tmp/roundfix`, the asset sync of the five Setup
   Snapshots into `internal/baseline/assets/setups` and its check;
4. a second `go build` of `$tmp/roundfix`, which embeds the new snapshots;
5. the restore preview and its confirmed apply (`skills-lock.json`,
   `.agents/skills/<name>`), and the folder and ref refusals;
6. `make baseline-digests`, which rewrites derived files under
   `internal/baseline/assets/{setups,source-baselines,formatter-fixtures,profiles}`
   and `internal/baseline/testdata`;
7. two runs of `roundfix baseline update --repo . --no-skills --no-history --yes --format text`
   with the binary of step 4, the second of which must report
   `File changes: 0`.

`internal/baseline/catalog.go` embeds `assets` (`//go:embed assets`), and
`Catalog.Digest` hashes every embedded file, so step 6 can change the catalog
digest after step 4 compiled the old bytes in. Step 7 then writes the old
digest into `docs/agents/setup-context.json`, and its convergence check
passes because the same binary planned both runs. The change inserts step
6a: a third `go build` of `$tmp/roundfix`, between steps 6 and 7. Nothing
after step 6a writes an embedded asset, so step 7 judges the tree on disk.

The step 4 build stays: the restore in step 5 needs a binary that embeds the
new snapshots, because the restore reads the profile's Setup Snapshot.

## Implementation Design

### Interfaces

The recipe change, in place in the `skills-third-party-bump` recipe:

```make
	$(MAKE) --no-print-directory -f "$(abspath $(firstword $(MAKEFILE_LIST)))" baseline-digests >&2; \
	$(GO) build -o "$$tmp/roundfix" ./cmd/roundfix; \
	roundfix baseline update --repo . --no-skills --no-history --yes --format text >&2; \
```

The fixture change in `internal/verifyselect/skills_bump_makefile_test.go`
(`newSkillsBumpFixture`): the `go` stub, at `build`, also copies the
working directory's `embedded-asset` to `<output>.embedded`, or writes
`absent` there when the file is missing; at `test probe`, it writes
`regenerated` to `embedded-asset`. The `roundfix` stub, at `baseline update`,
copies `"$0.embedded"` to `refreshed-from` before printing `File changes: 0`.

### Data Models

No changes. The Setup Manifest keeps its fields; the Task regenerates it once
because its verification entries record the `Makefile` digest.

### API Contracts

1. API Contract: refresh binary. Every `roundfix baseline update` the bump
   runs uses a `$tmp/roundfix` built after the bump's `make baseline-digests`
   step; the recorded call sequence is the one of today with
   `go build -o <tmp>/roundfix ./cmd/roundfix` inserted between
   `go test ./internal/baseline -run TestCatalogCompatibility -count=1` and
   the first `roundfix baseline update`. Every refusal, diagnostic and the
   final line keep their bytes.

### Surface Transcripts

None. The bump's standard output keeps its one final line and its refusals
keep their text; the only new effect is one more `go build`, which prints
nothing on success.

## Coverage Map

- Goal 1, User Story 1, Core Feature 1, Success Metric 2 → API Contract 1;
  the recipe line; Build Order 1.
- Goal 2, User Story 2, Core Feature 2, Success Metric 1 → Testing Approach;
  `TestTheSkillsBumpRefreshesWithABinaryBuiltAfterTheDigests`;
  `TestTheSkillsBumpPinsTheCatalogMainCommit`; Build Order 1.

Deviations from the PRD: None.

## Absence Evidence

1. A build of the refresh binary after the digests step — command: `awk 'index($0, "baseline-digests >&2;") { seen = 1; next } seen && index($0, "build -o") { built = 1 } seen && index($0, "baseline update") { exit } END { exit !built }' Makefile`; result: exit 1 at `b233761e`: the line after `baseline-digests >&2;` is the first `roundfix baseline update`.
2. A test that relates the refresh binary to the digests step — command: `grep -n 'embedded\|refreshed-from\|BuiltAfter' internal/verifyselect/skills_bump_makefile_test.go`; result: no match at `b233761e`; the only order proof is the pinned call list of `TestTheSkillsBumpPinsTheCatalogMainCommit`.

## ADR Alternatives

- ADR-0360: kept; the bump stays the only writer with the same order and
  refusals, and takes none of its rejected options (bumping only the lock
  and folders, `ref: main`, re-pinning test constants); the added build only
  makes its last step read what its earlier steps wrote.
- ADR-0364: not applicable; the presence-only check of a required external
  skill is unchanged.
- ADR-0380: not applicable; following and retiring a replaced catalog skill
  is unchanged.

## Integration Points

None outside the repository. The bump's `git ls-remote` and `git fetch` of
the catalog are unchanged; the tests serve the catalog from a local
repository the fixture creates.

## Testing Approach

The seam is the existing Catalog Skill Bump Makefile contract: the real
`Makefile` run by `make` with a fixture-owned `PATH`, a local Git catalog and
recording `go` and `roundfix` stubs (`skillsBumpCommand`, which passes
`DERIVED_DIGEST_PATHS` and `BASELINE_DIGEST_STEPS=probe:Probe`, so the
digests step runs `go test probe -run Probe -update -count=1` in the
repository directory).

- `TestTheSkillsBumpPinsTheCatalogMainCommit` (existing, declared break): its
  `want` call list gains `build` between the `TestCatalogCompatibility` call
  and the first `baseline update`.
- `TestTheSkillsBumpRefreshesWithABinaryBuiltAfterTheDigests` (new): seeds
  `repo/embedded-asset` with `before digests`, runs the bump, and asserts
  that `repo/refreshed-from` holds `regenerated`; on any other value it fails
  with `the managed refresh ran a binary that embeds <value>; want the asset
  make baseline-digests regenerated, "regenerated\n"`.

The other three bump tests keep their assertions; the stub changes write
only fixture files they never read.

Proven during authoring on a `git archive` copy: dropping the new build, or
moving it before `make baseline-digests`, makes both named tests fail with
the refresh reading `before digests`. The real-binary mechanism (Success
Metric 2) is reproduced by the QA gate with the module edit of the PRD's
Validation, since a real bump needs the whole catalog from another
repository.

## Build Order

1. The recipe line, the two test changes, then the Setup Manifest
   regenerated with
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`
   run twice, the second reporting `File changes: 0` (task_01). It must run
   after Specs 0373 and 0377 merged, which the queue order ensures; it reads
   no artifact of another Spec.
2. The QA gate (depends on: 1).

## Risks & Considerations

- A later edit that writes an embedded asset after the new build would
  bring the defect back; the new test fails for any build that precedes the
  digests step, and the Absence Evidence awk names where the build belongs.
- The extra compile costs one `go build` of `./cmd/roundfix` per bump,
  served by the shared build cache.
- The Setup Manifest is a shared file: the regeneration is the Task's last
  step and runs on whatever the merge base holds, so the Specs delivered
  first need no rebase of this plan.

## Glossary

None.

## Decisions

- Rebuild after the digests step rather than run the refresh with
  `go run`: the recipe already builds into `$tmp/roundfix` and runs it
  through the `roundfix` function with an isolated `HOME`; a third build
  keeps that.
- No currency check inside the bump: with the binary built after every
  embedded write, the existing "File changes: 0" check is sufficient.
- No ADR; ADR-0360 is kept as written.
