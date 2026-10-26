---
spec: 0320-the-baseline-setups-follow-the-consolidated-skills-catalog
prd: _prd.md
created: 2026-10-26
---

# The Baseline setups follow the consolidated skills catalog — Technical Spec

## Executive Summary

The five bundled Setup Snapshots are data the asset sync writes from one
upstream checkout (ADR-0191). The import is therefore a run of existing
commands: fetch the published commit read-only from the local upstream
checkout, run the asset sync and its `--check`, regenerate the derived
Baseline files with `make baseline-digests`, and refresh this repository's
Setup Manifest twice. The only source edit is in tests: a new test file
declares the upstream revision once and pins the republished `onioncry`
tree, and the two existing tests that each carried the old revision as a
literal read that constant. No Go production file, module, profile, template,
guide or installed skill changes.

The trade-off accepted: the new test pins one upstream tree digest, which the
next upstream change to `onioncry` will move; that is the point of the pin,
since it proves the sync, not a hand edit, wrote the tree.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential is added; the import reads a local checkout and no Verification
  or QA row reaches the network. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0191, ADR-0204, ADR-0241 and
  ADR-0322 govern the import, and ADR-0206, ADR-0219 and ADR-0246 are
  recorded as not applying, as `_prd.md` → Project Constraints states and
  ADR Alternatives below relates. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the five Setup Snapshots and this
  repository's Setup Manifest are Governed Paths, written only by the
  sanctioned commands of the authorization record. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0320-the-baseline-setups-follow-the-consolidated-skills-catalog/_authorization.md`;
  bounded files:
  `internal/baseline/assets/setups/go-cli-typescript-bun.json`,
  `internal/baseline/assets/setups/go.json`,
  `internal/baseline/assets/setups/python.json`,
  `internal/baseline/assets/setups/rust.json`,
  `internal/baseline/assets/setups/typescript.json`,
  `docs/agents/setup-context.json`.

## Current behavior

Read on 2026-10-26 at `86cb6b73` (v0.123.0).

1. `internal/baseline/assets/setups/{go,go-cli-typescript-bun,python,rust,typescript}.json`
   name `65d5f365cace27c4b468db1723213723b10c870b` as their own `source.ref`
   and as the `ref` of each of their 293 `github` skills; their 65 `repo`
   skills (Roundfix's own) carry no ref or tree. They hold 47, 115, 52, 40
   and 104 skills.
2. `roundfix baseline assets sync --source-dir <clone>/setups --check` exits
   `0` against a clean clone at `65d5f365` and `1` against `8d71fade`, with
   `skills.setup-snapshot.drift` for all five files. The command requires a
   clean Git checkout whose origin is
   `https://github.com/marcioaltoe/skills.git` (`internal/baseline/assets_sync.go`).
3. `internal/baseline/ai_memory_setup_test.go`
   `TestEveryUpstreamSetupFollowsTheAiMemoryRevision` declares a local
   `const revision = "65d5f365…"`; `internal/baseline/python_setup_test.go`
   declares the package constant `pythonPresetRevision` with the same value,
   read by `TestThePythonSetupIsThePublishedUpstreamPreset` and
   `TestEveryUpstreamSetupFollowsThePythonPresetRevision`. No other test,
   golden, guide, changelog line or production file names `65d5f365`; the
   ADRs that name it record their own decision dates and stay.
4. `make baseline-digests` owns `internal/baseline/testdata/catalog.digest`,
   `catalog.normalized.json`, `catalog.diagnostics.golden.json` and the plan
   characterization goldens; `TestCatalogCompatibility` fails when
   `catalog.digest` trails the catalog.
5. `docs/agents/setup-context.json` records the catalog's `catalogDigest`;
   `TestTheSetupManifestIsCurrent` (`internal/docscontract`, tag
   `docscontract`) fails while a managed refresh would change it.
6. This repository adopts `go-cli-tui`, whose setup is `go`; its
   `skills-lock.json` records mixed install revisions (`65d5f365`,
   `b3c45a45`, `a4e18e4f`, `14fdf46b`, `main`), and
   `TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot` compares the
   installed trees with the `go` snapshot by digest, not by ref.
7. `onioncry` appears only in `typescript.json` and
   `go-cli-typescript-bun.json`, required and dispatched by the `typescript`
   module. No Go file names it or reads a skill's `metadata.source`.

## System Architecture

- `internal/baseline/assets/setups/*.json` (five files): rewritten by the
  asset sync.
- `internal/baseline/testdata/` (seven files): rewritten by
  `make baseline-digests`.
- `docs/agents/setup-context.json`: rewritten by the managed refresh.
- `internal/baseline/upstream_setup_revision_test.go` (new): the revision
  constant, the `onioncry` tree constant and the tree test.
- `internal/baseline/ai_memory_setup_test.go`,
  `internal/baseline/python_setup_test.go`: read the constant.

## Implementation Design

### Interfaces

```go
// upstreamSetupRevision is the upstream skills commit every upstream Setup
// Snapshot and every skill it pins follow; it moves with each sync.
const upstreamSetupRevision = "8d71fade68c54693bedbd917371f0514311af157"

// onioncryTreeDigest is the onioncry tree that commit publishes.
const onioncryTreeDigest = "be502b02544be90059c126e9b07bfeb6a4af115aada86455845b1407a5eab8b3"

func TestTheTypeScriptSetupsCarryTheRepublishedOnioncryTree(t *testing.T)
```

### Data Models

None. The snapshot schema (`setup-context-driven/setup-snapshot/0.0.1`) is
unchanged; only values move.

### API Contracts

1. API Contract: the import. In an empty temporary repository, `git fetch`
   commit `8d71fade68c54693bedbd917371f0514311af157` from the local checkout
   `/Users/marcio/dev/skills`, check it out detached, and add the remote
   `origin` `https://github.com/marcioaltoe/skills.git` without contacting
   it. With a binary built under the Run's environment and a fresh `HOME`,
   run `roundfix baseline assets sync --source-dir <fetched>/setups --format
   text` (exit `0`, five `skills.setup-snapshot.updated` findings), then the
   same command with `--check` (exit `0`, `setup-context-driven audit: ok`).
   When the checkout lacks the commit, stop before any write.
2. API Contract: the derived files. `make baseline-digests` exits `0` and
   rewrites exactly `internal/baseline/testdata/catalog.diagnostics.golden.json`,
   `catalog.digest`, `catalog.normalized.json` and the four goldens
   `plan-characterization/{advisory-only-divergences,clean-adoption,idempotent-replan-after-verified-apply,same-baseline-changed-profile-and-catalog-digests}.golden.json`.
   Then, with a fresh `HOME`,
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`
   rewrites one managed entry of `docs/agents/setup-context.json`, and a
   second run reports `File changes: 0`.
3. API Contract: one revision constant. `upstreamSetupRevision` and
   `onioncryTreeDigest` live in
   `internal/baseline/upstream_setup_revision_test.go`. The local `revision`
   constant and its comment leave `ai_memory_setup_test.go`, and
   `pythonPresetRevision` leaves `python_setup_test.go`; each former use reads
   `upstreamSetupRevision`. Every other assertion of the two files keeps its
   bytes. This is the declared break: the expected revision moves from
   `65d5f365` to `8d71fade`.
4. API Contract: the tree test.
   `TestTheTypeScriptSetupsCarryTheRepublishedOnioncryTree` calls
   `t.Parallel()` first and, in subtests `typescript` and
   `go-cli-typescript-bun` (each calling `t.Parallel()`), loads the embedded
   catalog, finds the `onioncry` entry and asserts its `source.ref` equals
   `upstreamSetupRevision` and its `treeDigest` equals `onioncryTreeDigest`;
   a snapshot without `onioncry` fails.

### Surface Transcripts

None. No command, flag, exit code, message or stream placement changes.

### Invariants

```text
1. Every upstream Setup Snapshot and every github skill it pins name one upstream commit, the one upstreamSetupRevision holds.
2. A Setup Snapshot changes only through the asset sync; its derived files only through make baseline-digests and the managed refresh.
3. The import adds and removes no skill: each snapshot's ordered names equal the Delivery Base's.
4. This repository's installed skills, lock, modules, profiles, templates and rendered guides other than the Setup Manifest keep their bytes.
5. The typescript and go-cli-typescript-bun snapshots carry the onioncry tree the upstream commit publishes, and no other skill tree changes.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, Core Feature 1, Success Metric 1 → API Contract 1, Invariants 1-3, task_01
- Goal 2, Core Feature 3, Success Metrics 1 and 2 → API Contracts 1 and 4, Invariants 3 and 5, task_01
- Goal 3, Core Feature 2, Success Metrics 2 and 3 → API Contract 2, Invariants 2, task_01
- Goal 4, Core Feature 3, Success Metric 4 → API Contract 3, Interfaces, Invariants 1, task_01
- Goal 5, Success Metric 4 → Current behavior 6, Invariants 4, task_01 and task_02
- Success Metric 5 → the PRD's Unreachable Acceptance

## Absence Evidence

1. The shared revision constant and the onioncry tree test — command: `rg -n "upstreamSetupRevision|onioncryTreeDigest|onioncry" internal/baseline --type go`; result: no match at `86cb6b73`; the revision lives only in the two literals of Current behavior 3, and no test reads the `onioncry` tree.

## ADR Alternatives

- ADR-0191: the design follows the decision that the asset sync is the only writer of a Setup Snapshot; it takes neither the hand edit of refs nor a partial sync of one snapshot, so all five move together from one checkout.
- ADR-0204: the composed `go-cli-typescript-bun` snapshot is rewritten by the same sync as the union of `go` and `typescript`; the design never edits it on its own, the alternative that ADR rejected.
- ADR-0206: not applying; no upstream list is renamed and no skill leaves a preset, so no module, trigger or snapshot drops a skill.
- ADR-0219: not applying; no profile or profile draft changes, so no binding rule is exercised.
- ADR-0241: both snapshots that list `onioncry` carry the same new tree, so a repository-owned profile still resolves one contract and `restore.snapshot-conflict` stays unreachable; `TestEmbeddedSetupSnapshotsAgreeOnEveryExternalSkill` proves it.
- ADR-0246: not applying; the Retired Skill set is unchanged, and the two skills upstream retired never reached a snapshot, so nothing leaves the Baseline.
- ADR-0322: the design keeps the rule that every upstream snapshot names one commit, rejecting a `python` snapshot at a commit the others do not name; the shared test constant states that commit once.

## Integration Points

- The local checkout `/Users/marcio/dev/skills`, read only through `git
  fetch` into an empty temporary repository by task_01's import. Nothing is
  written there, and no remote is contacted. No Verification command or QA
  row reads it.

## Testing Approach

- `internal/baseline/upstream_setup_revision_test.go` (new), `package
  baseline`, with a `Suite:`/`Invariant:`/`Boundary IN:`/`Boundary OUT:`
  header, the two constants and the test of API Contract 4. A helper that
  loads one embedded snapshot must not reuse a name the package already
  declares.
- Changed: `ai_memory_setup_test.go` and `python_setup_test.go` as API
  Contract 3 states.
- Unchanged and run: the whole `./internal/baseline` package
  (`TestCatalogCompatibility`, `TestEmbeddedSetupSnapshotsAgreeOnEveryExternalSkill`,
  `TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot`, the
  derived-ownership and plan characterization tests), `./skills`, and
  `TestTheSetupManifestIsCurrent` with `-tags docscontract`.
- Prototype at authoring time: before the test edit only the three
  revision tests failed; after it every named package passed. Sabotages run
  on that prototype: restoring the snapshots from the Delivery Base fails both
  subtests of the new test on the `onioncry` ref and tree and the three
  revision tests; restoring `docs/agents/setup-context.json` fails
  `TestTheSetupManifestIsCurrent` ("Setup Manifest is stale"); restoring
  `catalog.digest` fails `TestCatalogCompatibility`; a hand edit of one
  `onioncry` tree line fails the new test at catalog load
  (`catalog.setup.digest.mismatch`).

## Build Order

1. The import, the derived files, the Setup Manifest refresh and the test
   edits (task_01; no dependency).
2. QA gate (task_02) (depends on: 1).

No Spec is listed in `requires`.

## Risks & Considerations

- The local checkout `~/dev/skills` belongs to another agent and carries an
  unpublished commit on its branch; the import only fetches the published
  object `8d71fade` from it, which the fetch proves by hash. A checkout that
  lacks the object stops the Task before any write.
- Upstream may publish again before delivery; the Spec still imports
  `8d71fade`, so the snapshots name a commit that the
  `marcioaltoe-skills` project announced, and a later sync is a later Spec.
- TypeScript adopters see `onioncry` as differing until they restore it; the
  refresh never overwrites it, so a locally adapted copy is not lost.

## Glossary

None.

## Decisions

- Move all five snapshots in one sync from `8d71fade`, the newest published
  commit, rather than `d12d7452`.
- Keep this repository's installed skills and lock: the `go` snapshot's
  trees are unchanged.
- One test constant for the revision, and a pinned `onioncry` tree as proof
  that the sync wrote the snapshots.
- No ADR and no guide, module or skill edit.
