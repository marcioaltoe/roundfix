---
spec: 0326-the-baseline-setups-stop-sending-agents-to-the-secondbrain-inbox
prd: _prd.md
created: 2026-10-26
---

# The Baseline setups stop sending agents to the Secondbrain inbox — Technical Spec

## Executive Summary

The five bundled Setup Snapshots are data the asset sync writes from one
upstream checkout (ADR-0191), and this repository's installed external skills
are held to its `go` snapshot (ADR-0221, ADR-0264). The fix is therefore a run
of existing commands, the same shape Spec 0320 used: fetch the published
commit `b56fc5ff` read-only into an empty repository, run the asset sync and
its `--check`, restore this repository's `knowledge-workspace` from the same
checkout, regenerate the derived Baseline files with
`make baseline-digests`, and refresh the Setup Manifest twice. The only
source edit is in Spec 0320's test file: its shared revision constant moves,
and a new test pins the corrected `knowledge-workspace` tree in all five
snapshots. No Go production file, module, profile, template or guide changes.

The trade-off accepted: the new test pins one upstream tree digest, which the
next upstream change to `knowledge-workspace` will move; the pin is what
proves the sync, not a hand edit, wrote every snapshot.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential is added; the import reads a local clone and no Verification or
  QA row reaches the network. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0060, ADR-0191, ADR-0204,
  ADR-0221, ADR-0224, ADR-0241, ADR-0264, ADR-0322 and ADR-0330 govern the
  import, and ADR-0206, ADR-0219, ADR-0246 and ADR-0338 are recorded as not
  applying, as `_prd.md` → Project Constraints states and ADR Alternatives
  below relates.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the five Setup Snapshots, this
  repository's Setup Manifest and its installed
  `.agents/skills/knowledge-workspace/SKILL.md` are Governed Paths, written
  only by the sanctioned commands of the authorization record. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0326-the-baseline-setups-stop-sending-agents-to-the-secondbrain-inbox/_authorization.md`;
  bounded files:
  `internal/baseline/assets/setups/go-cli-typescript-bun.json`,
  `internal/baseline/assets/setups/go.json`,
  `internal/baseline/assets/setups/python.json`,
  `internal/baseline/assets/setups/rust.json`,
  `internal/baseline/assets/setups/typescript.json`,
  `docs/agents/setup-context.json`,
  `.agents/skills/knowledge-workspace/SKILL.md`.

## Current behavior

Read on 2026-10-26 at `f06afdd7`.

1. `internal/baseline/assets/setups/{go,go-cli-typescript-bun,python,rust,typescript}.json`
   name `8d71fade68c54693bedbd917371f0514311af157` as the `ref` of every
   `github` skill; the four with a `github` source name it as `source.ref`;
   the composed `go-cli-typescript-bun.json` has
   `{"type":"composed","setups":["go","typescript"]}` and no ref. Each lists
   `knowledge-workspace` with `treeDigest`
   `19cb0bd73c1a2429e902b5e5c0e951fdd38d396352e7b1dc007748fc11a7aa7c`.
2. `roundfix baseline assets sync --source-dir <clone>/setups --check` exits
   `1` against a clean clone at `b56fc5ff`, with
   `skills.setup-snapshot.drift` for all five files. The command requires a
   clean Git checkout whose origin is
   `https://github.com/marcioaltoe/skills.git`.
3. `internal/baseline/upstream_setup_revision_test.go` (Spec 0320) declares
   `upstreamSetupRevision = "8d71fade…"` and `onioncryTreeDigest`, read by
   `TestTheTypeScriptSetupsCarryTheRepublishedOnioncryTree`,
   `TestEveryUpstreamSetupFollowsTheAiMemoryRevision`,
   `TestEveryUpstreamSetupFollowsThePythonPresetRevision` and
   `TestThePythonSetupIsThePublishedUpstreamPreset`. No test reads the
   `knowledge-workspace` tree.
4. This repository adopts `go-cli-tui` (setup `go`) and installs
   `.agents/skills/knowledge-workspace/SKILL.md` (one file, equal to the
   upstream file at `8d71fade`), locked in `skills-lock.json` at ref
   `65d5f365…` with `computedHash` `d108432a…`.
   `TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot` compares each
   locked skill's installed tree with the `go` snapshot by digest, and
   `TestCheckRepositoryMatchesRealRepository` (`./skills`) requires each
   installed tree to match its lock.
5. `make baseline-digests` owns `internal/baseline/testdata/catalog.digest`,
   `catalog.normalized.json`, `catalog.diagnostics.golden.json` and four plan
   characterization goldens; `docs/agents/setup-context.json` records the
   `catalogDigest`, held by `TestTheSetupManifestIsCurrent` (tag
   `docscontract`). Spec 0322, which delivers first, also regenerates these.
6. Every built-in profile selects `context-workflow`, which requires
   `knowledge-workspace` and dispatches it on two triggers about consulting
   the knowledge workspace and the Secondbrain. No module clause, rendered
   guide or Go file reads the skill's body.
7. The legacy Source Baseline
   `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/`
   still says sessions may write `inbox/**`. `Catalog.SourceBaseline`
   (`internal/baseline/source_contracts.go`) loads it only to verify its
   corpus digest and to account an old adoption's clauses; the managed
   refresh never renders it, and the rendered `secondbrain.md` (and its
   formatter golden for `standard-typescript-monorepo`) forbids any
   Secondbrain write.

## System Architecture

- `internal/baseline/assets/setups/*.json` (five files): rewritten by the
  asset sync.
- `.agents/skills/knowledge-workspace/SKILL.md` and `skills-lock.json`:
  rewritten by `roundfix baseline skills restore`.
- `internal/baseline/testdata/` (seven files): rewritten by
  `make baseline-digests`.
- `docs/agents/setup-context.json`: rewritten by the managed refresh.
- `internal/baseline/upstream_setup_revision_test.go`: the moved constant,
  a new tree constant and the new test.

## Implementation Design

### Interfaces

```go
// upstreamSetupRevision is the upstream skills commit every upstream Setup
// Snapshot and every skill it pins follow; it moves with each sync.
const upstreamSetupRevision = "b56fc5ffcad1ab366fbc03839a6ce8c662401178"

// knowledgeWorkspaceTreeDigest is the knowledge-workspace tree that commit
// publishes, the one that sends research through ai-memory.
const knowledgeWorkspaceTreeDigest = "490af37df4faba388e560380df0863b88881716d96e6185858918681182f2a96"

func TestEverySetupCarriesTheKnowledgeWorkspaceThatRetiredTheSecondbrainInbox(t *testing.T)
```

### Data Models

None. The snapshot schema and the lock schema are unchanged; only values
move.

### API Contracts

1. API Contract: the import. In an empty temporary repository, `git fetch`
   commit `b56fc5ffcad1ab366fbc03839a6ce8c662401178` from the local clone
   `/Users/marcio/.roundfix-operator/scratch/spec-0326/skills` with network
   protocols disabled, check it out detached, and add the remote `origin`
   `https://github.com/marcioaltoe/skills.git` without contacting it. With a
   binary built under the Run's environment and a fresh `HOME`, run
   `roundfix baseline assets sync --source-dir <fetched>/setups --format
   text` (exit `0`, five `skills.setup-snapshot.updated` findings), then the
   same command with `--check` (exit `0`, `setup-context-driven audit: ok`).
   When the clone lacks the commit, stop before any write.
2. API Contract: this repository's skill. With the same binary and a fresh
   `HOME`, run `roundfix baseline skills restore --profile go-cli-tui --skill
   knowledge-workspace --source-dir <fetched> --repo . --format json`; the
   preview exits `3` and lists exactly two planned changes, `refresh
   .agents/skills/knowledge-workspace/SKILL.md` and `update-lock-entry
   skills-lock.json` (ref `b56fc5ff…`). Rerun with `--confirm-plan
   <planDigest>` (exit `0`, `restore.completed`). A preview that lists any
   other change is a stop-and-report condition.
3. API Contract: the derived files. `make baseline-digests` exits `0` and
   rewrites only `internal/baseline/testdata/catalog.diagnostics.golden.json`,
   `catalog.digest`, `catalog.normalized.json` and the four goldens
   `plan-characterization/{advisory-only-divergences,clean-adoption,idempotent-replan-after-verified-apply,same-baseline-changed-profile-and-catalog-digests}.golden.json`.
   Then, with a fresh `HOME`,
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`
   rewrites one managed entry of `docs/agents/setup-context.json`, and a
   second run reports `File changes: 0`.
4. API Contract: the revision constant. `upstreamSetupRevision` moves from
   `8d71fade…` to `b56fc5ff…`; this is the declared break: the four tests of
   Current behavior 3 now expect the new commit. `onioncryTreeDigest` and
   every other line of the file keep their bytes, apart from the
   `Invariant:` header line, which adds the `knowledge-workspace` tree.
5. API Contract: the tree test.
   `TestEverySetupCarriesTheKnowledgeWorkspaceThatRetiredTheSecondbrainInbox`
   calls `t.Parallel()` first and, in one subtest per setup `go`,
   `go-cli-typescript-bun`, `python`, `rust` and `typescript` (each calling
   `t.Parallel()`), loads the embedded catalog, finds the
   `knowledge-workspace` entry and asserts its `source.ref` equals
   `upstreamSetupRevision` and its `treeDigest` equals
   `knowledgeWorkspaceTreeDigest`; a snapshot without the skill fails.

### Surface Transcripts

None. No command, flag, exit code, message or stream placement changes.

### Invariants

```text
1. Every upstream Setup Snapshot and every github skill it pins name one upstream commit, the one upstreamSetupRevision holds.
2. A Setup Snapshot changes only through the asset sync; an installed external skill and its lock entry only through baseline skills restore; the derived files only through make baseline-digests and the managed refresh.
3. The import adds and removes no skill: each snapshot's ordered names equal their values before the Task.
4. Only knowledge-workspace's tree changes, in all five snapshots, and the composed snapshot keeps its composed source.
5. This repository's installed knowledge-workspace matches both its lock entry and the go snapshot; no other installed skill, lock entry, module, profile, template, Source Baseline or rendered guide other than the Setup Manifest changes.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, Core Feature 1, Success Metric 1 → API Contract 1, Invariants 1-3, task_01
- Goal 2, Core Features 1 and 4, Success Metrics 1 and 2 → API Contracts 1 and 5, Invariants 3 and 4, task_01
- Goal 3, Core Feature 2, Success Metrics 2 and 4 → API Contract 2, Invariants 2 and 5, task_01
- Goal 4, Core Feature 3, Success Metrics 2 and 3 → API Contract 3, Invariant 2, task_01
- Goal 5, Core Feature 4, Success Metric 4 → API Contract 4, Interfaces, Invariant 1, task_01
- Success Metric 4 (scope) → Invariant 5, task_02
- Success Metric 5 → the PRD's Unreachable Acceptance

## Absence Evidence

1. The knowledge-workspace tree constant and its test — command: `rg -n "knowledgeWorkspaceTreeDigest|knowledge-workspace" internal/baseline --type go`; result: at `f06afdd7` only `python_setup_test.go` (the Python preset's member list) and `setup_composition_test.go` (a composition fixture) name the skill; no test reads its tree or its ref.

## ADR Alternatives

- ADR-0060: the design keeps the Source Baseline frozen and exhaustive; editing its corpus to drop the old inbox wording, the alternative that would rewrite a versioned input, would break its recorded corpus digest and change how an old adoption is accounted.
- ADR-0191: the asset sync stays the only writer of a Setup Snapshot; the design takes neither a hand edit of the `knowledge-workspace` lines nor a sync of one snapshot, so all five move together from one checkout.
- ADR-0204: the composed `go-cli-typescript-bun` snapshot is rewritten by the same sync as the union of `go` and `typescript`, never on its own.
- ADR-0221: this repository's required external skill is held to its Setup Snapshot, not only its lock, which is why the Task restores it rather than leaving it locked at the old tree.
- ADR-0224: the design restores this repository's skill in the same change that moves the `go` snapshot, and keeps holding installed trees to the lock and snapshot records the restore writes; the new `knowledgeWorkspaceTreeDigest` constant pins the tree the snapshots publish, the value the asset sync writes, not an installed tree, so it reintroduces none of the hand-pinned installed-tree digests that ADR rejected, as Spec 0320's `onioncryTreeDigest` did not.
- ADR-0206: not applying; no upstream list is renamed and no skill leaves a preset.
- ADR-0219: not applying; no profile or profile draft changes.
- ADR-0246: not applying; the Retired Skill set is unchanged.
- ADR-0338: not applying; the follow-up mode's handoff and message handling is untouched, and the corrected skill only names the ai-memory channel that ADR builds on.
- ADR-0241: all five snapshots carry the same new tree, so a repository-owned profile still resolves one contract and `restore.snapshot-conflict` stays unreachable; `TestEmbeddedSetupSnapshotsAgreeOnEveryExternalSkill` proves it.
- ADR-0264: the design restores deliberately, through the confirmed restore command, and leaves adopters' copies to their own deliberate restore; it does not make the managed refresh replace a differing copy, the alternative that ADR rejected.
- ADR-0322: every upstream snapshot, `python` included, names one commit; the shared test constant states it once.
- ADR-0330: the design changes no clause; it brings the vendored skill into line with the decision that the Secondbrain is read-only and ai-memory carries cross-project messages, rather than adding a Baseline clause that overrides the skill.

## Integration Points

- The operator's local clone
  `/Users/marcio/.roundfix-operator/scratch/spec-0326/skills`, read only
  through `git fetch` into an empty temporary repository by task_01's
  import. Nothing is written there, no remote is contacted, and no
  Verification command or QA row reads it.

## Testing Approach

- `internal/baseline/upstream_setup_revision_test.go`: API Contracts 4 and 5.
  The new test reuses `assetsSyncSnapshot` and `LoadEmbeddedCatalog` as the
  `onioncry` test does and declares no helper.
- Unchanged and run: the whole `./internal/baseline` package
  (`TestCatalogCompatibility`, `TestEmbeddedSetupSnapshotsAgreeOnEveryExternalSkill`,
  `TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot`, the
  derived-ownership and plan characterization tests), `./skills`
  (`TestCheckRepositoryMatchesRealRepository` and the lock-hash checks of
  `baseline_skill_contract_test.go`), and `TestTheSetupManifestIsCurrent`
  with `-tags docscontract`.
- Sabotages run at authoring time on a prototype of the whole change, each
  on its own worktree:
  - restore the five snapshots from before the Task: the five subtests of
    the new test fail naming ref `8d71fade…` and tree `19cb0bd7…`, and the
    `onioncry`, ai-memory and Python revision tests and
    `TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot` fail;
  - restore only `.agents/skills/knowledge-workspace/SKILL.md`:
    `TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot` fails,
    naming the skill and the restore command;
  - restore only `skills-lock.json`: `TestCheckRepositoryMatchesRealRepository`
    fails with `OutdatedExternal:[knowledge-workspace]`;
  - set `upstreamSetupRevision` back to `8d71fade…`: the new test and the
    four revision tests fail;
  - set `knowledgeWorkspaceTreeDigest` to `19cb0bd7…`: all five subtests of
    the new test fail;
  - restore `docs/agents/setup-context.json`: `TestTheSetupManifestIsCurrent`
    fails with "Setup Manifest is stale";
  - restore `internal/baseline/testdata/catalog.digest`:
    `TestCatalogCompatibility` fails.

## Build Order

1. The import, the skill restore, the derived files, the Setup Manifest
   refresh and the test edit (task_01; no dependency).
2. QA gate (task_02) (depends on: 1).

No Spec is listed in `requires`; the queue orders this Spec after 0322.

## Risks & Considerations

- Upstream may publish again before delivery; the Spec still imports
  `b56fc5ff`, the commit the skills project announced, and a later sync is
  a later Spec.
- Spec 0322 regenerates the same derived files first; this Spec reruns the
  regeneration on its own base, so neither pins the other's values.
- Adopters see `knowledge-workspace` as differing until they restore it; the
  refresh never overwrites it (ADR-0264), so a locally adapted copy is not
  lost, and the release note names the restore command.

## Glossary

None.

## Decisions

- Move all five snapshots in one sync from `b56fc5ff`, upstream `main` at
  authoring time.
- Restore this repository's `knowledge-workspace` and its lock entry through
  the confirmed restore command, from the same fetched checkout.
- Extend Spec 0320's test file: move the constant and pin the corrected
  tree in all five snapshots.
- Leave the legacy Source Baseline and the `.roundfixrc.yml` comment as they
  are; no ADR and no guide, module or owned-skill edit.
