---
spec: 0200-a-skill-snapshot-that-matches-its-upstream
prd: _prd.md
created: 2026-09-30
---

# A skill snapshot that matches its upstream — Technical Spec

## Executive Summary

Four changes, in order. The asset sync validates the catalog it would write
instead of the one it replaces, and the parity fixture's synthetic digests
become regenerated. The sync then refreshes four setup snapshots to
`a4e18e4` after the modules name the renamed skills and the Go CLI/TUI profile
names a seeded `go-tui` snapshot. The catalog gains the dispatch check, and the
Context7 capability probes the new skill while accepting the old one. Last,
this repository runs the adopter's update, reconcile and cleanup. The
trade-off accepted is that upstream lists arrive whole: the snapshots gain
every skill upstream added, and only the skills a module requires are
installed. That costs a larger catalog and buys a snapshot the next sync
leaves byte-identical.

## Project Constraints

- Identifier strategy: applicable — three renamed dispatch triggers
  (`trigger.core.context7-cli`, `trigger.rust.rust-expert`,
  `trigger.typescript.app-renderer-systems`), one new trigger
  (`trigger.core.exa-web-search`), one new diagnostic code
  (`catalog.profile.skill.dispatch-outside-setup`), one new setup snapshot
  (`go-tui`) and one new probe field (`priorSkills`), all in the existing
  forms. `capability.context7` keeps its identifier. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — embedded assets, local files and
  Git only; production steps read a local checkout of the upstream skills
  repository and clone it only when that checkout lacks the commit; no
  Verification command uses the network. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0191 governs the whole design.
  ADR-0072 keeps the Python parity contract: its words "recorded explicitly as
  a designed delta" govern the sync change and the fixture's new row. ADR-0073
  and ADR-0103 govern the managed refresh and its convergence. ADR-0081,
  ADR-0149 and ADR-0130 govern the regeneration and the audit. ADR-0058,
  ADR-0060 and ADR-0099 hold because no clause changes. ADR-0067 keeps derived
  repository-owned profiles on the built-in profile's setup. ADR-0080, ADR-0088,
  ADR-0091, ADR-0093, ADR-0094, ADR-0096, ADR-0104, ADR-0117, ADR-0155,
  ADR-0156, ADR-0166, ADR-0178 and ADR-0179 bind the gate and the checker.
  ADR-0180 to ADR-0184 are present and do not apply; the PRD's row gives the
  reason. ADR-0074 (hybrid semantic ownership) cites ADR-0067, ADR-0097 (QA row
  carry-forward) and ADR-0167 (the pre-PR Pull Request row) cite ADR-0080,
  ADR-0168 (the related-ADR gap) and ADR-0176 (citation checks read authored
  text) cite ADR-0093, ADR-0173 (citations a History Relocation breaks) cites
  ADR-0073, ADR-0181 (comparison with the Recommended Profile) cites ADR-0180,
  and ADR-0182 (Settlement Checks) cites ADR-0096; none applies, because this
  Spec changes none of their behaviors. ADR-0177 is reached only through
  ADR-0173 and is equally untouched. All hold. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30, recorded in [_authorization.md](_authorization.md); bounded
  files: `internal/baseline/assets/setups/go-cli.json`,
  `internal/baseline/assets/setups/go-tui.json`,
  `internal/baseline/assets/setups/rust-cli.json`,
  `internal/baseline/assets/setups/typescript-bun.json`,
  `internal/baseline/assets/profiles/go-cli-tui.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/baseline/assets/modules/core.json`,
  `internal/baseline/assets/modules/go.json`,
  `internal/baseline/assets/modules/rust.json`,
  `internal/baseline/assets/modules/typescript.json`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`,
  `internal/baseline/derived_ownership_test.go`,
  `docs/agents/skill-dispatch.md`, `docs/agents/setup-context.json`,
  `skills/baseline_skill_contract_test.go`,
  `.agents/skills/context7/SKILL.md`,
  `.agents/skills/context7-cli/SKILL.md`,
  `.agents/skills/context7-cli/references/docs.md`,
  `.agents/skills/context7-cli/references/setup.md`,
  `.agents/skills/context7-cli/references/skills.md`,
  `.agents/skills/golang-dependency-management/SKILL.md`,
  `.agents/skills/golang-safety/SKILL.md`,
  `.agents/skills/golang-safety/references/nil-safety.md`,
  `.agents/skills/golang-safety/references/slice-map-safety.md`,
  `.agents/skills/golang-structs-interfaces/SKILL.md`,
  `.agents/skills/golang-structs-interfaces/references/struct-fields.md`,
  `.agents/skills/golang-structs-interfaces/references/type-assertions.md`.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## System Architecture

| Concern | Source of truth | Derived by |
| --- | --- | --- |
| Upstream skill lists | `marcioaltoe/skills` at `a4e18e4fa223196b51d0fd8224e5a33b84f97717`, `setups/<id>.txt` | read only |
| Setup snapshots | `internal/baseline/assets/setups/*.json` | `roundfix baseline assets sync --source-dir <checkout>/setups` |
| Profile → setup, module skills and triggers | `internal/baseline/assets/profiles/*.json`, `modules/*.json` | hand edit under the grant |
| Dispatch check | `internal/baseline/catalog_validate.go` | — |
| Context7 capability | `internal/baseline/profile_alignment.go` (`universalCapabilities`) | — |
| Parity fixture rows | `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json` | the fixture transform below, then `make baseline-digests` |
| Goldens, pins, catalog snapshots | the assets | `make baseline-digests` |
| This repository's guides and Setup Manifest | the assets | `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text` |
| This repository's upstream skills and lock | the snapshot | `baseline update --yes`, then `baseline skills reconcile` |

### Measured facts

Measured on 2026-09-30 in the worktree at `9e439dbb` and in disposable clones.
The asset sync itself was run only with `--check`.

- `baseline assets sync --check` against a clean detached clone of upstream at
  `a4e18e4`, origin set to `https://github.com/marcioaltoe/skills.git`, reaches
  catalog validation and exits `2` with the five `outside-setup` diagnostics.
  After the module and profile edits and a seeded `go-tui` snapshot, the same
  check exits `2` with "Go-owned canonical Baseline assets are invalid", now
  for the new names. The sync validates the current catalog before it builds
  anything, so no edit order passes both validations. Core Feature 1 removes
  that deadlock.
- Upstream list sizes: `go-cli` 46, `go-tui` 48, `rust-cli` 35,
  `typescript-bun` 110 (it already lists `roundfix`). Against the current
  snapshots, `go-tui` adds `context7-cli`, `typesafe-ai` and eight Go skills and
  drops `context7`; `rust-cli` swaps `context7` and `rust` for `context7-cli`
  and `rust-expert` and adds `typesafe-ai`; `typescript-bun` swaps
  `context7` and `feature-systems-pattern` and adds ten skills.
- With simulated snapshots of that membership and the module edits,
  `make baseline-digests` rewrote the Standard TypeScript Monorepo
  `skill-dispatch.md` golden, its profile pin, the three catalog snapshots and
  the four plan goldens `advisory-only-divergences`, `clean-adoption`,
  `idempotent-replan-after-verified-apply` and
  `same-baseline-changed-profile-and-catalog-digests`. The Source Baseline
  corpus, `retention/`, `lock-hash-compatibility-v1.json` and the parity
  fixture were untouched. The managed refresh rewrote only
  `docs/agents/skill-dispatch.md` and `docs/agents/setup-context.json`, and a
  second refresh reported `current`.
- The same simulation failed five existing tests, each named in task_02:
  three asset-sync tests count three setups, `TestOutputsForCommand` compares
  the sanctioned outputs with a frozen 2026-08-06 enumeration, and
  `TestRunDoctorDerivesExternalSkillRequirementFromSetupManifest` lists the Go
  profile's external skills. The parity comparison, one of the three, also
  needs the new fixture rows.
- The parity fixture's synthetic `treeDigest` for a skill path `P` is the
  portable digest of one file, `SKILL.md`, holding `# <base(P)>\n`; the fixture
  values for `knowledge-workspace` and `exa-web-search` equal it.
- `baseline skills restore` writes a lock entry with the snapshot's full
  commit as `ref`. `baseline update --yes --skills-source-dir <dir>` with no
  guidance change still runs the skills stage and restores a required skill
  whose lock entry is missing; a second update reports `current`.
  `baseline skills reconcile` at `a4e18e4` today refuses with
  `reconcile.required-removed` for `context7`, because `core` still requires
  it.
- Against `a4e18e4`, this repository's `context7-cli` is absent, and
  `golang-dependency-management`, `golang-safety` and
  `golang-structs-interfaces` are installed without lock entries and differ in
  1, 3 and 3 files. Ten further required skills match their locks and differ
  from the snapshot; the PRD leaves them to a Backlog Entry.
- This repository's `docs/agents/setup-context.json` already records a stale
  Makefile digest on `9e439dbb`. Any managed refresh rewrites it.

## Implementation Design

### Interfaces

```go
// internal/baseline/assets_sync.go — syncAssets
// The LoadCatalog(os.DirFS(assetRoot)) call before the source is read is
// removed. After the snapshots are built, the overlay catalog
// (assets + overrides) is always validated, overrides possibly empty:
//   len(overrides) == 0 and invalid → "Go-owned canonical Baseline assets are invalid: …"
//   len(overrides) > 0 and invalid  → "Generated setup snapshots are incompatible with the Baseline catalog: …"
// Both keep AssetsSyncInvalid, their actions and their codes.

// internal/baseline/catalog_validate.go — validateProfiles, when the setup exists:
// profileNamedSkills returns, sorted and unique, the skill of every
// skillDispatch entry (its "skill", else its "id") of the selected modules and
// every skill of an activation bundle whose owner is a selected module.
func profileNamedSkills(catalog *Catalog, activations document, selected []string) []string
// each name absent from the setup adds
// l.add("catalog.profile.skill.dispatch-outside-setup", profileID, skill)

// internal/baseline/profile_alignment.go
// capability.context7: Probe {"skill": "context7-cli", "priorSkills": ["context7"]}
// collectInstalledSkillEvidence checks "skill", then each prior skill in order;
// the first present SKILL.md is the evidence and its path is SourcePath.
var universalCapabilityRestoreSkills = map[string]string{
	"capability.context7": "context7-cli",
	"capability.exa":      "exa-web-search",
}
```

Test helpers, all in test files:

```go
// internal/baseline/assets_sync_synthetic_test.go (task_01)
func assetsSyncSyntheticSkillFile(skillPath string) []byte // "# " + path.Base(skillPath) + "\n"
func assetsSyncSyntheticTreeDigest(skillPath string) string
func parityFixtureDigestFindings(fixture map[string]any) []string
```

`buildAssetsSyncSource` writes `assetsSyncSyntheticSkillFile`, and
`regenerateBaselineCompatibilitySetups` sets every `github` skill's
`treeDigest` to `assetsSyncSyntheticTreeDigest` before it computes the setup
digest.

### Data Models

- `go-tui.json` is seeded as a byte copy of `go-cli.json` with
  `"id": "go-cli"` replaced by `"id": "go-tui"`. The seed gives the sync the
  owned entries' minimum versions, the Roundfix entry and the activation
  bundles; the sync rewrites everything else.
- `go-cli-tui.json`: `"setup": "go-cli"` becomes `"setup": "go-tui"`; nothing
  else changes.
- No schema changes. The probe map gains `priorSkills`, a list of skill names.

### Fixed texts

Replace strings and numbers in place; never pass a module through a JSON
encoder. Raise each module version by one per Task that changes it, from the
value on the Task's base.

- task_02, `core.json`: `"context7"` becomes `"context7-cli"` in
  `requiredSkills` and `skillDispatch`, and `trigger.core.context7` becomes
  `trigger.core.context7-cli`. The `when` text stays.
- task_02, `rust.json`: `"rust"` becomes `"rust-expert"` in `requiredSkills`
  and `skillDispatch`, and `trigger.rust.rust` becomes
  `trigger.rust.rust-expert`.
- task_02, `typescript.json`: every `feature-systems-pattern` becomes
  `app-renderer-systems`, including the trigger identifier.
- task_02, `go.json`: `requiredSkills` gains `golang-dependency-management`,
  `golang-safety` and `golang-structs-interfaces` in alphabetical position.
- task_03, `core.json`: `requiredSkills` gains `exa-web-search` after
  `evidence-gate`, and `skillDispatch` gains, after the `evidence-gate` entry:

```json
{
  "skill": "exa-web-search",
  "triggers": [
    {
      "id": "trigger.core.exa-web-search",
      "when": "Searching the web broadly when local sources and current documentation do not answer the question."
    }
  ]
}
```

- task_03, `capability.context7` NextAction: "Add the context7-cli skill to
  the Repository Skill Set, then rerun capability evaluation." Title and
  Explanation stay.

### The refresh procedure (task_02)

```bash
sha=a4e18e4fa223196b51d0fd8224e5a33b84f97717
src="$HOME/dev/skills"
git -C "$src" cat-file -e "$sha^{commit}" 2>/dev/null || src=https://github.com/marcioaltoe/skills.git
tmp="$(mktemp -d)"
git clone --quiet --no-local "$src" "$tmp/skills"
git -C "$tmp/skills" checkout --quiet --detach "$sha"
git -C "$tmp/skills" remote set-url origin https://github.com/marcioaltoe/skills.git
go run -buildvcs=false ./cmd/roundfix baseline assets sync --source-dir "$tmp/skills/setups" --format text
```

It runs after the module, profile and seed edits. Then the fixture transform,
the test edits, `make baseline-digests` and the managed refresh, twice.

The fixture transform rebuilds `manifest.setups`, `managedEntryLedger` and
`normalizedOutput.result` from the refreshed snapshots; the regeneration
fills every digest:

```jq
.input.sourceRevision as $rev
| def row($s): {activationBundles: $s.activationBundles, digest: "", id: $s.id,
    schemaVersion: $s.schemaVersion, version: $s.version,
    source: {path: $s.source.path, ref: $rev, repository: "example/skills", type: "github"},
    skills: [$s.skills[] | if .source.type == "github"
      then {name, path, source: {path: .source.path, ref: $rev, repository: "example/skills", type: "github"}, treeDigest: ""}
      else {minimumVersion, name, path, source: {name: .source.name, type: "repo"}} end]}
    | if .activationBundles == null then del(.activationBundles) else . end;
  .manifest.setups = [$s[] | row(.)]
| .managedEntryLedger = [.manifest.setups[] | {digest: "", id}]
| .normalizedOutput.result.findings = [.manifest.setups[] | {
    action: "Review the snapshot diff and run asset validation.",
    code: "skills.setup-snapshot.updated", managedId: ("setup." + .id),
    message: "Setup snapshot was synchronized from the canonical source.",
    path: ("assets/setups/" + .id + ".json"), severity: "info"}]
| .normalizedOutput.result.summary.info = (.manifest.setups | length)
```

Run as
`jq --slurpfile s <(jq -s 'sort_by(.id)[]' internal/baseline/assets/setups/*.json) -f <program> <fixture>`.
`plannedByteSequence` and `fileIdentities` are the frozen Python run's byte
identities; no test reads them and they stay.

### The repository procedure (task_04)

1. `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --yes --skills-source-dir "$src" --format text`,
   with `$src` as above. The skills stage reinstalls the owned skills
   byte-identically and restores `context7-cli`, `exa-web-search` only if its
   lock hash differs, and the three Go skills whose lock entries are missing,
   each at the snapshot's commit.
2. `go run -buildvcs=false ./cmd/roundfix baseline skills reconcile --repo . --profile go-cli-tui --source marcioaltoe/skills --revision a4e18e4fa223196b51d0fd8224e5a33b84f97717 --source-dir "$src" --format json`,
   then the same command with `--confirm-plan <planDigest>`. It removes the
   `context7` lock entry and keeps the tree.
3. `git rm -r .agents/skills/context7`.
4. `skills/recommended.txt` becomes the sorted lock keys, and
   `upstreamManagedSkillTreeDigest` takes the value the failing test prints.

### API Contracts

1. API Contract: `roundfix baseline assets sync --source-dir <dir> [--check]`
   — a catalog invalid only until the refresh is refreshed (exit `0`), or
   reported as drift under `--check` (exit `1`); a catalog the refresh leaves
   invalid exits `2` with "Generated setup snapshots are incompatible with the
   Baseline catalog"; an invalid catalog with no drift exits `2` with
   "Go-owned canonical Baseline assets are invalid". Nothing is written on
   exit `1` or `2`.
2. API Contract: embedded catalog — every built-in profile's setup lists every
   skill its selected modules dispatch or its owned bundles name, else
   `catalog.profile.skill.dispatch-outside-setup`.
3. API Contract: rendered guide `docs/agents/skill-dispatch.md` — names
   `context7-cli`, `rust-expert` or `app-renderer-systems` where it named the
   old skill, carries `trigger.core.exa-web-search`, and for `go-cli-tui` the
   three Go skills stay.
4. API Contract: `roundfix baseline plan` and `baseline update` — profile
   alignment satisfies `capability.context7` with `context7-cli` or
   `context7`, and otherwise reports it blocking with a next action naming
   `--skill context7-cli`.
5. API Contract: `roundfix doctor` in this repository — `skills: ok (42
   required: 14 Roundfix-owned, 28 external)`.

## Coverage Map

- Goal 1 → The refresh procedure; API Contract 1; Testing Approach 1 and 2.
- Goal 2 → Fixed texts; API Contracts 3 and 4; Testing Approach 2 and 3.
- Goal 3 → Interfaces (`profileNamedSkills`); API Contract 2; Testing
  Approach 3.
- Goal 4 → The repository procedure; API Contract 5; Testing Approach 4.
- Story 1 → Interfaces (`syncAssets`); Testing Approach 1.
- Story 2 → API Contracts 2 and 3; Testing Approach 2 and 3.
- Story 3 → API Contract 4; Testing Approach 3.
- Story 4 → API Contract 5; Testing Approach 4.
- Core Feature 1 → Interfaces (`syncAssets`); API Contract 1.
- Core Feature 2 → The refresh procedure; Fixed texts.
- Core Feature 3 → Interfaces (`profileNamedSkills`); API Contract 2.
- Core Feature 4 → Interfaces (`universalCapabilityRestoreSkills`); API
  Contract 4.
- Core Feature 5 → The repository procedure; API Contract 5.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → Testing Approach 2 and 3.
- Success Metric 4 → Testing Approach 3.
- Success Metric 5 → Testing Approach 4.
- Success Metric 6 → Testing Approach 5.
- API Contract 1 → Interfaces (`syncAssets`).
- API Contract 2 → Interfaces (`profileNamedSkills`).
- API Contract 3 → Fixed texts.
- API Contract 4 → Interfaces (`universalCapabilityRestoreSkills`).
- API Contract 5 → The repository procedure.

## Integration Points

- **Spec 0195.** Its sync rule keeps the Roundfix entry the upstream lists of
  `go-cli`, `go-tui` and `rust-cli` omit. Its release step names checks this
  Spec does not rename. Its `autonomous-work` module requires `roundfix`, so a
  snapshot without the entry fails task_02's Verification.
- **Spec 0199.** It edits `core.json` and `typescript.json` clauses and the
  activation file; this Spec edits other lines of the two modules and not the
  activation file.
- **Adopters.** The next update rewrites `skill-dispatch.md` and the Setup
  Manifest, then restores the new required skills. The old directory and lock
  entry stay until the adopter reconciles and deletes it; a second update
  reports `current`. Repository-owned profiles derived from `go-cli-tui` take
  `go-tui` with it.
- **Upstream.** Read only, at one commit. The Tasks never write
  `~/dev/skills`.

## Testing Approach

No test uses the network. Each negative case is its own test.

1. **Sync (task_01).** New `internal/baseline/assets_sync_upstream_rename_test.go`,
   on the existing temporary target and source helpers, renaming `handoff` to
   `handoff-next` (a `core` skill no bundle names):
   `TestAssetSyncFollowsASkillRenamedUpstream`,
   `TestAssetSyncCheckReportsDriftForARenameTheRefreshRepairs`,
   `TestAssetSyncRefusesARenameNoSetupProvides` and
   `TestAssetSyncStillRefusesAnInvalidCatalogWithoutDrift`. New
   `internal/baseline/assets_sync_synthetic_test.go`:
   `TestTheParityFixtureDigestsFollowTheSyntheticSource` and
   `TestASyntheticDigestThatDiffersIsReported`.
2. **Refresh (task_02).** New `internal/baseline/upstream_skill_names_test.go`:
   `TestNoCatalogEntryNamesASkillRenamedUpstream`,
   `TestARenamedSkillNameInTheCatalogIsReported`,
   `TestTheGoCLITUIProfileTakesTheGoTUISetup`,
   `TestEveryModuleRequiresEverySkillItDispatches` and
   `TestADispatchedSkillNoModuleRequiresIsReported`. The five existing tests of
   Measured facts pass after their literal updates.
3. **Check and capability (task_03).** New
   `internal/baseline/catalog_dispatch_outside_setup_test.go`:
   `TestEveryBuiltInProfileSetupListsEverySkillItsGuidesName`,
   `TestADispatchedSkillOutsideTheProfileSetupIsReported`,
   `TestABundledSkillOutsideTheProfileSetupIsReported` and
   `TestAModuleTheProfileDoesNotSelectIsNotChecked`. New
   `internal/baseline/context7_capability_test.go`:
   `TestTheContext7CapabilityIsSatisfiedByTheCurrentSkill`,
   `TestTheContext7CapabilityIsSatisfiedByItsPriorSkill`,
   `TestAMissingContext7SkillIsRemediatedWithTheCurrentName` and
   `TestEveryRestorableCapabilitySkillIsRequiredAndInEverySetup`.
4. **This repository (task_04).** New
   `internal/cli/this_repository_skill_set_test.go`:
   `TestThisRepositoryHoldsEveryRequiredExternalSkill` resolves the external
   set from this repository's Setup Manifest, requires a complete readiness
   and no lock entry or directory for `context7`, `feature-systems-pattern` or
   `rust`; `TestARepositoryMissingARequiredExternalSkillIsReported` runs the
   same helper on a temporary repository.
5. **Regeneration.** Each Task that edits assets runs `TestFormatterComposition`,
   `TestCatalogCompatibility`, `TestBaselinePlanCharacterization` and
   `TestBaselineCompatibilityCorpus`, then the read-only managed refresh.

## Build Order

1. The sync validates what it produces, and the fixture's synthetic digests
   regenerate, task_01 (depends on: none).
2. The snapshot refresh, the renames, the `go-tui` profile setup and the Go
   requirements, task_02 (depends on: 1).
3. The dispatch check, `exa-web-search` in `core`, and the Context7
   capability, task_03 (depends on: 2).
4. This repository takes its own update, and the user guide explains it,
   task_04 (depends on: 3).
5. Terminal QA, task_05 (depends on: 1, 2, 3, 4).

task_02 needs task_01's sync to run at all. task_03's check fails on the
current snapshots, and its capability test needs `context7-cli` in every
setup. task_04 restores what task_02 and task_03 require, and edits the user
guide after task_01 and task_03 do.

## Risks & Considerations

- **Upstream moves again.** Every step names the commit, and the clone checks
  it out detached; a moved `main` changes nothing.
- **The sync is a write.** It runs only through the recoverable transaction
  and only in task_02. A refused run writes nothing.
- **Parity fixture by transform.** The transform copies the refreshed
  snapshots; the comparison test still proves the sync reproduces them from a
  synthetic source.
- **Hazard: generated files.** No pin, golden, catalog snapshot, fixture
  digest or rendered guide is hand-edited. The two test literals and the
  frozen-enumeration addition are ordinary edits the Tasks name.
- **Hazard: this repository's skills.** Only the files the grant bounds may
  change. A restore that would touch another file stops the Task.

## Decisions

- The sync validates its result. See ADR-0191.
- `go-tui` is seeded from `go-cli` so the sync inherits owned minimums and
  bundles.
- The frozen-enumeration test names `go-tui.json` as an addition instead of
  changing its oracle.
- The dispatch check covers owned activation bundles, because the rendered
  guide names their skills too.
- `exa-web-search` joins `core`, the module every profile selects, because a
  universal required capability already demands it.
- The installed-but-trailing skills wait for their Backlog Entry.
