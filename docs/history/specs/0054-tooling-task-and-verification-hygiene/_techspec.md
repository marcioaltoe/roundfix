---
spec: 0054-tooling-task-and-verification-hygiene
prd: _prd.md
created: 2026-07-28
---

# Tooling task and verification environment hygiene — Technical Spec

## Executive Summary

The regeneration path is the load-bearing decision. `roundfix baseline
assets sync` cannot serve: it deliberately carries repo-owned skills'
existing `contentDigest` forward (propagating staleness) and requires a
clean committed checkout, which is hostile to a target run mid-edit. Instead
the three validating test suites gain `-update` golden-refresh flags, and
`make baseline-digests` runs them in dependency order — setups, then
normalized catalog, then parity corpus. That adds zero new packages, keeps
the validators and regenerators the same code, and makes drift impossible
between what the gate checks and what the command writes. The second
decision is unifying the digest algorithm: the contract test's local hash
helper sorts by raw bytes while production `skillhash.Sum` sorts by
collation; the regenerator standardizes on the production
`skills.SkillFolderHash`, and the test helper is retired in its favor. The
Daemon-side guards (executable staging refusal, green-on-entry precondition)
attach to existing seams — `FilterStageablePaths` and the pre-Agent phase of
the task worker — so every commit path and every Task type is covered
without new orchestration.

## Project Constraints

- Identifier strategy: not applicable — Make target names, digest values,
  and Skill names keep their existing identities; no project-owned Internal
  Identifier is created. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local build tooling, digest
  computation, Daemon staging, and documentation only. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0045 keeps review Runs requiring
  a clean tracked working tree; ADR-0038 keeps the single Verification
  repair; ADR-0057 keeps Task status Daemon-owned; ADR-0081 makes
  sanctioned regeneration fallout of the authorized Skill edit. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — on 2026-07-28, the maintainer expressly
  authorizes changes to exactly `Makefile` and `.gitignore`, to exactly
  `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md`, and to
  the deterministic Skill-digest fallout in exactly
  `internal/baseline/assets/setups/go-cli.json`,
  `internal/baseline/assets/setups/rust-cli.json`,
  `internal/baseline/assets/setups/typescript-bun.json`,
  `internal/baseline/testdata/catalog.digest`,
  `internal/baseline/testdata/catalog.normalized.json`,
  `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`,
  and `internal/baseline/testdata/parity-corpus/v1/manifest.json`. No other
  protected tooling mutation is authorized. Source:
  `docs/agents/agent-instructions.md`.

## System Architecture

- `skills/` and `internal/baseline` test suites — the three validators gain
  regeneration modes: the authorial-sync contract test rewrites each setup
  snapshot's repo-sourced `contentDigest` (via the production
  `skills.SkillFolderHash`) and its top-level digest (via the existing
  canonical-JSON hash); the catalog compatibility test rewrites the
  normalized snapshot and digest from `Catalog.Normalized()` /
  `Catalog.Digest()`; the compatibility-corpus test rewrites the
  parity-corpus fixture's digests and the manifest's sha256/byte rows while
  leaving the frozen literals (inventory digest, frozen date, test counts)
  untouched.
- `Makefile` — the `baseline-digests` target sequencing the three
  regenerations, the `GOCACHE` default, and `.PHONY` registration.
- `internal/daemon` — the executable-file drop in `FilterStageablePaths`
  (covering Batch, Task, and QA commits, which all route through it) and
  the green-on-entry precondition in the task worker's pre-Agent phase.
- `.gitignore` — `/roundfix` and the repository-local Go cache directory.
- `docs/agents/specific-repository.md` (fully repository-authored) — the
  commit choreography, the ADR-0081 policy, and the cache guidance;
  `docs/agents/agent-instructions.md` is setup-owned and is not touched.

## Implementation Design

### Interfaces

Regeneration flags (per owning test package, standard golden-file idiom):

```go
var update = flag.Bool("update", false, "regenerate derived digest artifacts")
// TestAuthorialSkillSync:        rewrites assets/setups/*.json when *update
// TestCatalogCompatibility:      rewrites testdata/catalog.{normalized.json,digest}
// TestBaselineCompatibilityCorpus: rewrites parity-corpus fixture + manifest rows
```

Make targets:

```make
GOCACHE ?= $(CURDIR)/.gocache
export GOCACHE

baseline-digests:
	$(GO) test ./skills -run TestAuthorialSkillSync -update -count=1
	$(GO) test ./internal/baseline -run TestCatalogCompatibility -update -count=1
	$(GO) test ./internal/baseline -run TestBaselineCompatibilityCorpus -update -count=1
	$(GO) test ./internal/baseline -run TestFormatterComposition -update -count=1
	$(GO) test ./internal/baseline -run TestReadoptionCompatibilityMaintainedFixture -update -count=1
```

The last two steps cover the module/asset chain proven manually on
2026-07-28 (commit `66d9b63`): the formatter-composition update mode writes
the golden fixtures from the plan's own generated postimages and re-pins the
profile's `goldenDigest`; the maintained-fixture update mode regenerates the
source-baseline manifest's byte-range entries marker-based — each entry's
span is recomputed from its `source-baseline-entry` delimiters in the corpus,
never by offset arithmetic — then rewrites the identity and index digests
(manifest sha256, length-prefixed corpus walk). The frozen source-baseline
corpus bytes themselves are edited only when a clause's canonical guidance
changes, and the update mode self-validates every entry digest before
writing.

Failure messages: when run without `-update`, each validator's mismatch
message ends with `run 'make baseline-digests'` — the hint lives in this
repository's tests, not in the product's catalog diagnostics, which serve
target repositories too.

Executable staging drop (in `FilterStageablePaths`, beside the existing
external-path and symlink drops):

```go
// mode check on the worktree file: a regular file with any execute bit is
// dropped with Reason "executable file"; symlinks are already handled.
```

Green-on-entry precondition (task worker, before Agent Session creation):

```go
// When task.Verification contains the configured repository Verification
// command, run that command once in the Task's worktree first. On failure,
// settle the Task failed with reason
// "repository not green on entry: <bounded output tail>" and create no
// Agent Session.
```

### Data Models

None. The dropped-path report reuses `DroppedStagePath`; the precondition
failure reuses the existing Task settlement reason field. `.gocache/` is a
build byproduct, ignored and never staged (ignored paths never enter the
snapshot diff).

### API Contracts

- `make baseline-digests` is idempotent: a second run reports no changes.
  It regenerates all three setup snapshots — a non-roundfix Skill edit
  cascades into `go-cli.json` and `rust-cli.json` as well; those two files
  are not part of this Spec's own mutation surface because implementing the
  target changes no Skill, but the command must handle them for the Specs
  that will (0053, 0060, per ADR-0081).
- `make verify` behavior changes only in environment determinism: with no
  `GOCACHE` exported, the repository-local cache is used identically inside
  and outside the ACP sandbox; an exported `GOCACHE` wins unchanged.
- Daemon staging: an executable file in the snapshot diff is dropped and
  reported with its path and mode; the commit proceeds with the remaining
  paths exactly as the symlink and external drops do today. The existing
  contract that explicitly selected ignore-matched paths stage with
  `add -f` is preserved — the ignore entry protects the snapshot layer, the
  mode guard protects the staging layer, and the two are deliberately
  redundant.
- Precondition: only Tasks whose Verification includes the repository-wide
  gate pay the pre-run; its failure is a Task failure with a precondition
  reason, not a Run stop, and ADR-0038's single Verification repair does
  not apply to a Task that never started Agent work.

## Coverage Map

- Goal 1 / Stories 1–2 → `-update` regeneration modes, `baseline-digests`
  target, stale-pin failure hint, ADR-0081 policy text (Features 1–3).
- Goal 2 / Story 3 → `GOCACHE` default in the Makefile (Feature 5).
- Goal 3 / Story 6 → `/roundfix` ignore entry + executable staging drop
  (Feature 6).
- Goal 4 / Stories 4–5 → choreography documentation in
  `specific-repository.md` and the green-on-entry precondition (Features 4,
  7). The one-pass QA authorization-audit half of Feature 4 is delivered
  through Spec 0053's authorized qa-gate Skill edit and is referenced, not
  duplicated, here.

## Integration Points

- Git — staging and snapshot behavior via the Daemon's existing runner.
- The ACP sandbox — only through the cache default; no sandbox-specific
  code.

## Testing Approach

Existing seams: the `skills` package contract tests
(`TestAuthorialSkillSync` gains its update mode and a
regenerated-then-green round-trip test), `internal/baseline`
(`TestCatalogCompatibility`, `TestBaselineCompatibilityCorpus`,
`TestEmbeddedCatalog` for the cascade message), `internal/daemon`
(`TestTaskCommitDropsSymlinkCrossingTaskFileAndCommitsRepositoryPaths` is
the direct analogue for the executable drop;
`TestGitCommitterStagesSelectedTrackedPathMatchedByGlobalIgnore` pins the
`add -f` contract the guard must not break; task-worker tests cover the
precondition settle path and that a green repository adds exactly one gate
run). Makefile behavior is proven by the round-trip: edit a fixture skill in
a temp copy, regenerate, assert green — hermetic inside the update-mode
tests themselves rather than shelling to `make`.

## Build Order

1. Regeneration modes in the validator suites for both derived chains —
   the Skill-digest chain on the unified `skills.SkillFolderHash`
   algorithm, and the module/asset chain (formatter goldens plus
   marker-based source-baseline regeneration) — with stale-pin messages
   naming the command.
2. `Makefile`: `baseline-digests` target, `GOCACHE` default, `.PHONY`
   (depends on: 1).
3. `.gitignore`: `/roundfix` and the repository-local cache directory
   (no dependency).
4. Daemon executable staging drop (no dependency).
5. Green-on-entry precondition in the task worker (no dependency).
6. Documentation and policy: choreography and ADR-0081 in
   `specific-repository.md`, roundfix Skill pair alignment, and the
   authorized derived digest pins (depends on: 1–5).

## Risks & Considerations

- **Algorithm unification.** Switching the contract test's byte-order hash
  helper to the collated production hash could change a pin value if the
  two orders ever disagree on a real tree; the regeneration command makes
  that a one-command repair, and the round-trip test proves agreement.
- **`-update` writes into `assets/`.** Unusual for a test flag but
  deliberate: the validator and regenerator share one code path, so they
  cannot drift. The flag never runs in CI or `make verify`.
- **Precondition cost.** One extra repository-gate run per tooling Task; it
  is paid only by Tasks that already declare the gate, and it replaces a
  full wasted Agent turn plus recovery when the repository is red.
- **Frozen corpus literals.** The corpus regeneration must not touch the
  frozen inventory digest, date, and count fields — asserted by a test that
  regenerates against an unchanged tree and diffs nothing.

## Decisions

- Regeneration lives in the validators via `-update`, not in
  `baseline assets sync`, which by design carries digests forward and
  demands a clean checkout.
- The digest algorithm standardizes on production `skills.SkillFolderHash`.
- The cache default changes only the unset case; explicit `GOCACHE` always
  wins.
- Policy and choreography text land only in repository-authored carriers;
  setup-owned guides are untouched. See ADR-0081.
