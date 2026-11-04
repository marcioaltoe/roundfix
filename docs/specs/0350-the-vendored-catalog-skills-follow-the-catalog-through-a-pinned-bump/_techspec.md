---
spec: 0350-the-vendored-catalog-skills-follow-the-catalog-through-a-pinned-bump
prd: _prd.md
created: 2026-11-04
---

# The vendored catalog skills follow the catalog's main through a pinned bump — Technical Spec

## Executive Summary

The Catalog Skill Bump is a Makefile target that composes writers this
repository already trusts: the asset sync for the five upstream Setup
Snapshots (ADR-0191), the confirmed `roundfix baseline skills restore` for
the vendored folders and their lock entries (ADR-0264), `make
baseline-digests` and the managed refresh. Its own code only resolves and
fetches one catalog commit, checks what it must refuse, and verifies the
result. The primary trade-off: the bump moves Baseline content adopters
receive on every run, because ADR-0224 holds this repository's skills to the
`go` snapshot; a bump of the lock alone would leave this repository failing
its own snapshot test and adopters on the old trees. The tests stop
hand-pinning what the bump moves: the expected revision comes from the lock,
and the regression guards refuse the tree each earlier fix replaced.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint or credential; only
  the new target may use the network, and no Verification or QA row lets it.
  Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0360 (this Spec's decision),
  ADR-0191, ADR-0204, ADR-0221, ADR-0224, ADR-0241, ADR-0264, ADR-0322 and
  ADR-0326, as the PRD's Project Constraints state; ADR-0206, ADR-0219 and
  ADR-0350 do not apply. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the Cycle
  AX Rules answer "Mesmas do AW (Recommended)", which releases the Makefile
  when the design needs it; Spec-contained authorization record:
  `docs/specs/0350-the-vendored-catalog-skills-follow-the-catalog-through-a-pinned-bump/_authorization.md`;
  bounded files: `.agents/skills/systematic-debugging/CREATION-LOG.md`,
  `.agents/skills/systematic-debugging/test-academic.md`,
  `.agents/skills/systematic-debugging/test-pressure-1.md`,
  `.agents/skills/systematic-debugging/test-pressure-2.md`,
  `.agents/skills/systematic-debugging/test-pressure-3.md`, `Makefile`,
  `docs/agents/setup-context.json`, `docs/agents/specific-repository.md`,
  `internal/baseline/assets/setups/go-cli-typescript-bun.json`,
  `internal/baseline/assets/setups/go.json`,
  `internal/baseline/assets/setups/python.json`,
  `internal/baseline/assets/setups/rust.json`,
  `internal/baseline/assets/setups/typescript.json`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## System Architecture

Current behavior, measured at `6a0ed904`:

1. `skills-lock.json` holds 33 entries with `"source": "marcioaltoe/skills"`
   at seven refs; the 12 owned skills (`OWNED_SKILLS` in the Makefile,
   `roundskills.Names()` in Go) are absent. `.agents/skills` holds 45 folders:
   the 33 plus the 12 owned.
2. `make skills-update` runs `bunx skills experimental_install`, `bunx skills
   update -p -y` and `bun run fmt`; the repository has no `package.json`, and
   the published installer re-installs an entry at its own `ref`.
3. The five upstream Setup Snapshots name `86df860` and pin every catalog
   skill's `treeDigest`. `roundfix baseline skills restore` reads the
   snapshots embedded in the binary, so a binary built before an asset sync
   restores the old commit.
4. `TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot` compares the
   29 lock entries `go-cli-tui` requires with the `go` snapshot; the four
   unrequired ai-memory skills are held only to their `computedHash`
   (`assertUpstreamManagedSkillLocksMatch` in `skills/`), which a re-hash
   satisfies.
5. `upstreamSetupRevision` is a constant in
   `internal/baseline/upstream_setup_revision_test.go`, read by five tests in
   that file, `ai_memory_setup_test.go` and `python_setup_test.go`;
   `onioncryTreeDigest`, `uiUxProMaxTreeDigest` and
   `knowledgeWorkspaceTreeDigest` pin the trees `86df860` publishes. Specs
   0326 and 0333 re-pinned them by hand.

Components after this Spec:

- **Catalog Skill Bump** (`Makefile`, target `skills-third-party-bump`), new.
  It replaces `skills-update`.
- **Catalog Revision helper** (`upstreamSetupRevision(t)` in
  `internal/baseline/upstream_setup_revision_test.go`), replacing the
  constant.
- **Vendored skill contract**
  (`internal/baseline/vendored_catalog_skills_docscontract_test.go`), new.
- The bump's output at `c8b2e85`: lock, five snapshots, five deletions,
  derived files and the Setup Manifest.

## Implementation Design

### Interfaces

```make
SKILLS_CATALOG ?= https://github.com/marcioaltoe/skills.git
SKILLS_CATALOG_ORIGIN := https://github.com/marcioaltoe/skills.git
SKILLS_REF ?=            # empty: resolve the catalog's refs/heads/main
skills-third-party-bump: ## Pin the vendored catalog skills and Setup Snapshots to one catalog commit (default: its main)
```

```go
// upstreamSetupRevision returns the one full commit every catalog entry of
// ../../skills-lock.json names; it fails the test on a branch, a short or
// mixed ref, or no catalog entry.
func upstreamSetupRevision(t *testing.T) string

const replacedOnioncryTree = "f918508da464fcc37b436a5a26fc243ffa159ead9e07b19e991451ed0ca53b6c"
const incompleteUiUxProMaxTree = "b13db828c6a7fbebf8612185dfde2e4cf20414f4c8fdaa33cc50b3b59f427d24"
const inboxKnowledgeWorkspaceTree = "19cb0bd73c1a2429e902b5e5c0e951fdd38d396352e7b1dc007748fc11a7aa7c"

func TestEveryVendoredCatalogSkillIsTheTreeItsCatalogCommitPublishes(t *testing.T)
```

The three tree constants are the trees the snapshots pinned before Specs
0320, 0333 and 0326 replaced them (`git show f34147ae^`, `30590e30^`,
`2b4b37d3^` of `internal/baseline/assets/setups/typescript.json`).

### Data Models

No schema changes. After the bump, every catalog entry of `skills-lock.json`
keeps its fields and order; `ref` becomes the Catalog Revision and
`computedHash` the restored folder's hash.

### API Contracts

1. API Contract: the bump's inputs and refusals. Run from the repository
   root. `SKILLS_REF` empty → `git ls-remote "$(SKILLS_CATALOG)"
   refs/heads/main` (the only network read besides the fetch). The commit must
   match `^[0-9a-f]{40}$`. The catalog set is every lock entry whose `source`
   line is `marcioaltoe/skills`; it must be non-empty, contain no
   `OWNED_SKILLS` name, and equal the `.agents/skills` folders that are not
   owned. Any failure prints `skills-third-party-bump: <cause>` on stderr,
   exits non-zero and runs no build, sync, restore or regeneration. The
   messages name: the owned skill (`Roundfix-owned skill <name>`), the bad ref
   (`full 40-character`), the set difference (`differ: <names>`), the
   unserved commit (`does not serve commit <sha>`).
2. API Contract: the bump's steps, in order, in a temporary directory under
   `$TMPDIR` removed on exit: `git init`; `git fetch --depth=1
   "$(SKILLS_CATALOG)" <sha>`; detached checkout; remote `origin` set to
   `SKILLS_CATALOG_ORIGIN` without contacting it; `$(GO) build` of
   `./cmd/roundfix`; `baseline assets sync --source-dir <catalog>/setups` and
   its `--check`; a second build (the restore reads the embedded snapshots);
   `baseline skills restore --profile <profile of docs/agents/setup-context.json>`
   with one `--skill` per catalog name in sorted order, `--source-dir
   <catalog> --repo . --format text`, and on exit `3` the same command with
   `--confirm-plan <Plan Digest>`; every other preview exit fails. Every
   `roundfix` call runs with `HOME` set to an empty directory inside the
   temporary directory.
3. API Contract: the bump's checks and regeneration. After the restore each
   catalog name's folder under `<catalog>/skills/*/` must equal
   `.agents/skills/<name>` (`diff -r`), else `.agents/skills/<name> differs
   from the catalog at <sha>`; every catalog entry's `ref` must equal the
   commit, else `not pinned to <sha>: <names>`. Then `$(MAKE) -f <this
   Makefile> baseline-digests` and two managed refreshes (`baseline update
   --repo . --no-skills --no-history --yes --format text`) with the rebuilt
   binary; the second must print `File changes: 0`. On success stdout carries
   one line: `skills-third-party-bump: <N> catalog skills and every upstream
   Setup Snapshot pinned to <sha>`.
4. API Contract: the bump at `c8b2e85`. task_02 runs `GIT_ALLOW_PROTOCOL=file
   make skills-third-party-bump
   SKILLS_CATALOG=/Users/marcio/.roundfix-operator/scratch/spec-0350/skills
   SKILLS_REF=c8b2e8511a7ff89dbe9e87a56cc3c29a6acbcb84`. The prototype changed
   exactly: the five setups, five deletions under
   `.agents/skills/systematic-debugging/`, `skills-lock.json` (33 refs, one
   `computedHash`), the seven derived files of
   [task_02](task_02.md) and `docs/agents/setup-context.json`.
5. API Contract: the test changes. `upstreamSetupRevision` becomes the
   function of Interfaces at every call site (`upstream_setup_revision_test.go`,
   `ai_memory_setup_test.go`, `python_setup_test.go`); each tree test keeps
   its name and its revision check, and replaces its equality with the
   replaced tree by an inequality whose message names the Spec that replaced
   it. The new contract, under `//go:build docscontract` and `//verify:always`
   (ADR-0326), asserts for each lock entry, in sorted order: it is not a
   `roundskills.Names()` skill; an embedded snapshot pins a `treeDigest` for
   it; the `portableRestoreDigest` of `.agents/skills/<name>` equals it. It
   also asserts that each `github` snapshot names `upstreamSetupRevision(t)`.
   `TestRepositoryContractTestsDeclareTheirRelevance`
   (`internal/verifyselect/contracts_repository_test.go`) accepts a
   `docscontract` file of `internal/baseline` as `always`, as it accepts the
   two of `internal/cli`, and counts three such files.

### Surface Transcripts

None. No `roundfix` command, flag, exit code or message changes. The new
Makefile target is a maintainer tool whose stdout line and refusal messages
API Contracts 1 and 3 state and its own test asserts (Success Metric 1).

### Invariants

```text
1. Every catalog entry of skills-lock.json and every github Setup Snapshot and skill name one full commit.
2. No Roundfix-owned skill is a lock entry, and the bump refuses one before any write.
3. Every vendored folder equals the catalog tree at that commit; a re-hashed local edit fails the contract.
4. The snapshots change only through the asset sync, the vendored folders and lock entries only through the confirmed restore, the derived files only through make baseline-digests and the managed refresh.
5. A second bump at the same commit changes no file.
6. No test names a catalog revision or a tree digest the bump writes.
```

## Coverage Map

- Goal 1, Core Feature 1, Success Metrics 1 and 5 → Catalog Skill Bump, API Contracts 1-3, Invariants 2, 4 and 5, task_01
- Goal 2, Core Feature 2, Success Metrics 2 and 3 → API Contract 4, Invariants 1 and 3, task_02
- Goal 3, Success Metric 2 → API Contract 1, Vendored skill contract, Invariant 2, task_01, task_02
- Goals 4 and 5, Core Feature 3, Success Metrics 4 and 6 → Catalog Revision helper, Vendored skill contract, API Contract 5, Invariants 3 and 6, task_02
- Core Feature 4 → `CONTEXT.md` and `docs/agents/specific-repository.md`, task_03
- Every Success Metric at the audited head → task_04

Deviations from the PRD: None.

## Absence Evidence

1. Catalog Skill Bump — command: `rg -n "ls-remote|skills restore|assets sync" Makefile`; result: no match at `6a0ed904`; the only skills targets are `skills-update`, `skills-check`, `skills-sync`, `skills-version-check`, `skills-sync-check`, `skills-install` and `skills-link`, none of which pins a commit or runs the restore.
2. Catalog Revision helper — command: `rg -n "skills-lock.json" internal/baseline/upstream_setup_revision_test.go internal/baseline/ai_memory_setup_test.go internal/baseline/python_setup_test.go`; result: no match; the revision is the constant `upstreamSetupRevision`.
3. Vendored skill contract — command: `rg -l "go:build docscontract|portableRestoreDigest|roundskills.Names" internal/baseline --glob "*_test.go"`; result: six restore, trailing, snapshot-profile and asset-sync test files call `portableRestoreDigest` on temporary fixtures; none names `roundskills.Names` or carries the `docscontract` tag; no test compares every lock entry's folder with a snapshot tree regardless of the profile's required set.

## ADR Alternatives

- ADR-0360: the design is that ADR's decision; it takes none of the three alternatives it rejected (a lock-only bump, `ref: main`, re-pinning constants per bump).
- ADR-0191: the asset sync stays the only writer of the five snapshots; the bump calls it instead of editing a snapshot or syncing one setup.
- ADR-0204: the composed snapshot is written by the same sync as the union of `go` and `typescript`, never by the bump directly.
- ADR-0221: this repository stays held to its snapshot; the new contract widens that hold to the four unrequired catalog skills instead of relaxing it.
- ADR-0224: amended by ADR-0360 only in where the expected revision comes from; the design keeps holding installed trees to the records the writers produce and adds no hand-pinned installed-tree digest; the three tree constants name trees that must never return, not trees the snapshots carry.
- ADR-0241: all five snapshots move from one checkout, so a repository-owned profile still resolves one contract per skill.
- ADR-0264: the bump restores through the preview and its confirmed Plan Digest; it never makes the managed refresh replace a differing copy, the alternative that ADR rejected.
- ADR-0322: every upstream snapshot, `python` included, names the one commit the lock names.
- ADR-0206: not applying; no upstream list is renamed and no skill leaves a preset, so its rename path is not taken.
- ADR-0219: not applying; no profile draft or profile binding changes, and the composed snapshot it relies on is written by the same sync.
- ADR-0350: not applying; the Retired Skill set is unchanged and the bump deletes only files the catalog removed from a kept skill, through the restore.
- ADR-0326: the test that reads `.agents/skills` is a `docscontract` file declared `//verify:always`, not a plain code test.

## Integration Points

- The skills catalog `https://github.com/marcioaltoe/skills.git`: read by the
  target's `git ls-remote` and `git fetch` only, when the maintainer runs it.
- The operator's clone `/Users/marcio/.roundfix-operator/scratch/spec-0350/skills`:
  read only by task_02's bump, through a local `git fetch`. No Verification
  or QA row other than the QA's outside-evidence row reads it.

## Testing Approach

- `internal/verifyselect/skills_bump_makefile_test.go` (new, the package
  already holds the Makefile contract tests and installs the suite guard):
  runs the repository Makefile's target in a temporary repository fixture
  with `make -f`, a local fake catalog Git repository (two skills under
  `skills/<group>/`, `setups/go.txt`, branch `main`), `GO=` a shell stub that
  records its arguments and copies a recording `roundfix` stub on `build -o`,
  `DERIVED_DIGEST_PATHS=` an empty directory and `BASELINE_DIGEST_STEPS=probe:Probe`
  (the existing seams of `baseline-digests`), and `GIT_ALLOW_PROTOCOL=file`.
  `TestTheSkillsBumpPinsTheCatalogMainCommit` asserts API Contracts 2-3;
  `TestTheSkillsBumpRefusesWhatItMustNotPin` has one subtest per refusal of
  API Contracts 1 and 3. The stub writes only under the test's temporary
  directory.
- `internal/baseline`: the five revision tests and three tree tests run
  unchanged by name; the new contract with `-tags docscontract`.
- Sabotages, each run at authoring time on a prototype: dropping the owned
  check, the post-restore `diff`, the second build or the ref check each
  fails its named subtest; a hand edit of `.agents/skills/ai-memory-retrieval/SKILL.md`
  with a re-hashed lock entry passes `./skills` and
  `TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot` and fails the
  new contract.

## Build Order

1. The Catalog Skill Bump target, the retirement of `skills-update`, its
   fixture test and the managed refresh the Makefile edit requires (task_01).
2. The bump at `c8b2e85`, the Catalog Revision helper, the tree tests and the
   vendored skill contract (task_02) (depends on: 1).
3. The glossary terms and the repository rule (task_03) (depends on: 2).
4. The QA gate (task_04) (depends on: 1, 2, 3).

## Risks & Considerations

- The catalog can change a setup's membership: the sync or the catalog
  validation then refuses, and the bump stops before the restore; that move
  becomes a Spec (ADR-0360).
- The Run sandbox has no network; task_02 reads the operator's clone, and
  `GIT_ALLOW_PROTOCOL=file` makes any other transport fail.
- An adopter's next refresh lists moved skills under `skills.differing`; the
  release note tells adopters to restore them.

## Glossary

- adds: **Catalog Skill Bump**
- adds: **Catalog Revision**
- not a term: **Catalog Revision helper** — the name of a test function in this document, not a domain concept
- not a term: **Vendored skill contract** — the name of a test in this document, not a domain concept

## Decisions

- Compose the existing writers instead of a new Go command: they already
  prove provenance, digests and atomic writes. See ADR-0360.
- Move the Setup Snapshots in the same bump: ADR-0224's test and Doctor
  require it. See ADR-0360.
- Retire `make skills-update`: it floated the entries and its `bun run fmt`
  has no `package.json` to run.
- Put the target's test in `internal/verifyselect`, beside the other Makefile
  contract tests, instead of `skills/`, which installs no suite guard.
