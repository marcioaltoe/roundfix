---
spec: 0208-a-baseline-that-follows-the-reshaped-skills-catalog
prd: _prd.md
created: 2026-10-01
---

# A Baseline that follows the reshaped skills catalog — Technical Spec

## Executive Summary

The work has four steps. First, the three setup snapshots take their upstream
names, `go-cli` retires, and `review` and `triage` leave the modules, all in
one asset-sync refresh to `b3c45a4`. Second, the `external-triage` module moves
from one unlabelled rule to eight clauses with force, with Source Baseline rows
and a declared replacement. Third, `core` requires `typesafe-ai` and
`crafting-effective-readmes`. Last, this repository runs the adopter's update,
reconcile and cleanup. No production Go code changes: every step is an asset
edit, a sanctioned regeneration and tests.

The trade-off: removing the skills and renaming the setups must land in one
Task, because the sync validates the catalog it produces, and a catalog that
still requires `review` is refused against a snapshot that no longer lists it.
That costs one larger Task and buys a refresh that is never half-applied.

## Project Constraints

- Identifier strategy: applicable — as the PRD states: the setup identifiers
  `go`, `rust` and `typescript`; the triggers `trigger.core.typesafe-ai` and
  `trigger.core.crafting-effective-readmes`; eight
  `clause.external-triage.<name>` clauses. Built-in profile identifiers do not
  change. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — embedded assets, local files and
  Git only. The refresh clones the local checkout `~/dev/skills` and stops,
  instead of reaching a network remote, when it lacks the pinned commit. No
  Verification command uses the network. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — the PRD's row is authoritative.
  ADR-0206 governs the whole design. ADR-0191 governs the rename and the tree
  that is never deleted by Roundfix. ADR-0072 governs the frozen parity corpus.
  ADR-0058, ADR-0060 and ADR-0099 govern the `replaced` disposition and the
  Source Baseline rows, and ADR-0202 and ADR-0186 govern the clause force and
  wording. ADR-0073, ADR-0103, ADR-0081, ADR-0149, ADR-0130 and ADR-0192
  govern the regeneration and the refresh. ADR-0204 is edited only where it
  names component setups. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30, reaffirmed on 2026-10-01, recorded in
  [_authorization.md](_authorization.md); bounded files:
  `internal/baseline/assets/setups/go-cli.json`,
  `internal/baseline/assets/setups/go-tui.json`,
  `internal/baseline/assets/setups/rust-cli.json`,
  `internal/baseline/assets/setups/typescript-bun.json`,
  `internal/baseline/assets/setups/go.json`,
  `internal/baseline/assets/setups/rust.json`,
  `internal/baseline/assets/setups/typescript.json`,
  `internal/baseline/assets/profiles/go-cli-tui.json`,
  `internal/baseline/assets/profiles/rust-cli.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/baseline/assets/modules/core.json`,
  `internal/baseline/assets/modules/typescript.json`,
  `internal/baseline/assets/modules/external-triage.json`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/external-triage.md`,
  `internal/baseline/assets/source-baselines/index.json`,
  `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/baseline.json`,
  `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/manifest.json`,
  `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/external-triage.md`,
  `internal/baseline/derived_ownership_test.go`,
  `docs/agents/skill-dispatch.md`, `docs/agents/setup-context.json`,
  `skills/baseline_skill_contract_test.go`,
  `.agents/skills/review/SKILL.md`,
  `.agents/skills/typesafe-ai/LICENSE`,
  `.agents/skills/typesafe-ai/SKILL.md`,
  `.agents/skills/crafting-effective-readmes/README.md`,
  `.agents/skills/crafting-effective-readmes/SKILL.md`,
  `.agents/skills/crafting-effective-readmes/references/art-of-readme.md`,
  `.agents/skills/crafting-effective-readmes/references/make-a-readme.md`,
  `.agents/skills/crafting-effective-readmes/references/standard-readme-example-maximal.md`,
  `.agents/skills/crafting-effective-readmes/references/standard-readme-example-minimal.md`,
  `.agents/skills/crafting-effective-readmes/references/standard-readme-spec.md`,
  `.agents/skills/crafting-effective-readmes/section-checklist.md`,
  `.agents/skills/crafting-effective-readmes/style-guide.md`,
  `.agents/skills/crafting-effective-readmes/templates/internal.md`,
  `.agents/skills/crafting-effective-readmes/templates/oss.md`,
  `.agents/skills/crafting-effective-readmes/templates/personal.md`,
  `.agents/skills/crafting-effective-readmes/templates/xdg-config.md`,
  `.agents/skills/crafting-effective-readmes/using-references.md`.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## System Architecture

| Concern | Source of truth | Derived by |
| --- | --- | --- |
| Upstream skill lists | `marcioaltoe/skills` at `b3c45a45f1bccd3b33aaecaaa22947d942f2fc02`, `setups/{go,rust,typescript}.txt` | read only |
| Setup snapshots | `internal/baseline/assets/setups/*.json` | seed by hand, then `roundfix baseline assets sync --source-dir <clone>/setups` |
| Profile → setup; module skills, triggers, clauses | `profiles/*.json`, `modules/*.json` | hand edit under the grant |
| Source Baseline rows | corpus file, `manifest.json` row, `index.json` `entryIds` | hand edit of text, identity, force and carrier; the regeneration fills the rest |
| Parity fixture `asset-sync.json` | the refreshed snapshots | the jq transform below, then the regeneration |
| Goldens, profile pin, catalog snapshots, plan goldens | the assets | `make baseline-digests` |
| This repository's guides and Setup Manifest | the assets | `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text` |
| This repository's upstream skills and lock | the snapshots | `baseline update --yes --skills-source-dir`, then `baseline skills reconcile` |

### Measured facts

Measured on 2026-10-01 at `c3be3bc9` (main) and in a disposable clone, against
a local clone of upstream at `b3c45a4`. The asset sync command ran only with
`--check`. The rehearsal wrote its simulated snapshots by calling the sync's
own snapshot builder in that clone.

- Today `baseline assets sync --check` exits `2` with four errors: "Canonical
  setup source file is missing." for `go-cli`, `go-tui`, `rust-cli` and
  `typescript-bun`.
- Upstream sizes: `go` 43 skills, `rust` 36, `typescript` 101. All three list
  `roundfix`, `typesafe-ai` and `crafting-effective-readmes`, and none lists
  `review` or `triage`. `go` lists `bubbletea` and `tui-design`. `go` and
  `typescript` share 32 entries with equal paths, and their union is 112.
- The sync lists every `setups/*.json` file and rebuilds each from the
  upstream file of the same name, keeping owned entries and activation bundles
  from the current file. Seeding `go.json`, `rust.json` and `typescript.json`
  as copies of their predecessors, with `id` and `source.path` changed, is how
  a snapshot is renamed. After the task_01 edits, `--check` exits `1` with
  three drift findings. After the refresh and the regeneration, it exits `0`.
- The catalog requires a Source Baseline row for every clause of a rule the
  Standard TypeScript Monorepo profile requires: without the rows, the sync
  refuses with `catalog.sourceBaseline.required-clause.missing` for each of the
  eight clauses.
- Today a Standard TypeScript Monorepo adopter with `triage.external` enabled
  and drifted artifact digests is refused on refresh: "retention transition has
  1 unaccounted clause(s): rule.external-triage". The rule has no clauses, so
  the classifier finds no successor. task_02's declared replacement repairs it.
- With all four Tasks applied, `make baseline-digests` rewrote both goldens
  (`skill-dispatch.md`, `external-triage.md`), the profile pin, the three
  catalog snapshots, the four plan goldens (`advisory-only-divergences`,
  `clean-adoption`, `idempotent-replan-after-verified-apply`,
  `same-baseline-changed-profile-and-catalog-digests`), the Source Baseline
  `baseline.json`, `manifest.json` and `index.json`, and the parity
  `v1/manifest.json`. The managed refresh rewrote only
  `docs/agents/skill-dispatch.md` and `docs/agents/setup-context.json`.
- The same rehearsal failed these existing tests until their edits:
  `TestCatalogDiagnosticCharacterization` (`catalog_test.go` names
  `setups/rust-cli.json` and an anchor `"rust",` that now occurs twice in
  `profiles/rust-cli.json`); `TestTheGoCLITUIProfileTakesTheGoTUISetup`;
  `TestBaselineAssetsSyncRefreshProducesCanonicalTreeAndIsIdempotent`,
  `TestAssetsSyncCompatibilityMatchesMaintainedPythonContract`,
  `TestAssetsSyncCheckIsReadOnlyAndReportsDrift` and
  `TestAssetsSyncProvenanceAndPreMutationRefusals` (setup counts and names);
  `TestOutputsForCommand` (frozen enumeration); the four force tests of
  `clause_characterization_test.go`; `TestReadoptionCompatibilityMaintainedFixture`
  (entry count); `TestThisRepositoryHoldsEveryRequiredExternalSkill` and
  `TestARepositoryMissingARequiredExternalSkillIsReported` (this repository
  lacks the two skills); and three `./skills` tests on the upstream tree digest.
  Every other test passed.
- `baseline update --yes --skills-source-dir <clone>` restored
  `crafting-effective-readmes` (14 files) and `typesafe-ai` (2 files) at
  `b3c45a4`. A second update reported `File changes: 0`. `baseline skills
  reconcile --profile go-cli-tui` previewed one edit, `remove-lock-entry
  skills-lock.json [review]`, with the tree retained. The confirmed run
  applied it, and a third run reported no changes. The other lock entries are
  present at `b3c45a4`.
- The frozen parity fixtures `greenfield-*`, `update-rust-cli`,
  `stale-plan-refusal`, `skill-restoration*`, `profile-change-*`,
  `readoption-preservation`, `atomic-rollback` and `missing-capability-refusal`
  still record the Python run's setup names. No test compares those fields,
  and `testdata/parity-corpus/_ownership.yml` declares them frozen: repointing
  them "was tried on 2026-07-30 and reverted".

## Implementation Design

### Interfaces

No production file changes. Test helpers are pure functions in test files, so
a negative test can feed them a literal.

```go
// internal/baseline/upstream_removed_skills_test.go (task_01)
var upstreamRemovedSkills = []string{ /* the seventeen names below */ }
func removedSkillFindings(modules, bundles, setups []document) []string
// "<owner> names removed skill <name>" for requiredSkills, a dispatch skill,
// a trigger id ending in ".<name>", a bundle skill, or a setup skill name.
var renamedSetups = map[string]string{"go-cli-tui": "go", "rust-cli": "rust",
	"standard-typescript-monorepo": "typescript"}
func setupNameFindings(profiles []document, setups []document) []string
// a listed profile on another setup; a setup named go-cli, go-tui, rust-cli or
// typescript-bun; a github-sourced setup whose source.path is not
// "setups/<id>.txt"; a setup without the Roundfix-owned "roundfix" entry.
```

```go
// internal/baseline/external_triage_clauses_test.go (task_02)
var externalTriageClauses = []struct{ id, force, guidance string }{ /* Fixed texts */ }
func externalTriageClauseFindings(module document, guide string) []string
// the rule's clauses differ in identifier, order, force or guidance; the first
// clause does not declare replaces ["rule.external-triage"]; the guide lacks
// "- **<force>**: <guidance>" for a clause, or holds a bullet with no force.
func newExternalTriageAdopter(t *testing.T) (PlanRequest, *Catalog)
// newClauseReplacementAdopter with triage.external = true.
```

```go
// internal/baseline/core_skill_requirements_test.go (task_03)
var coreSkillTriggers = map[string]struct{ id, when string }{ /* Fixed texts */ }
func coreSkillFindings(core, typescript document, guides map[string]string) []string
// a core skill not required or not dispatched with its exact trigger; the
// typescript module still names crafting-effective-readmes; a rendered
// skill-dispatch guide that lacks a "- `<id>`: <when>" line.
```

### Data Models

- Each new snapshot is seeded as a byte copy of its predecessor with only
  `"id"` and `"source.path"` changed (`go-tui` to `go`, `rust-cli` to `rust`,
  `typescript-bun` to `typescript`). `go-cli.json` and the three predecessors
  are deleted. The sync rewrites everything else.
- Each profile's `"setup"` value changes in place. Nothing else in the file is
  hand-edited.
- `external-triage.json` becomes `setup-context-driven/module-v3` with
  `"repositoryExtensions": []`, as `bun.json` is.

### Fixed texts

Replace strings and numbers in place. Never pass a module through a JSON
encoder. Raise every version by one from its value on the Task's base, and
start a new rule at 1.

task_01, `core.json`: remove `"review"` from `requiredSkills`, and remove the
`review` dispatch with `trigger.core.review`. `typescript.json`: remove
`"triage"` and its dispatch with `trigger.typescript.triage`. Raise both module
versions.

task_02, `external-triage.json`: raise the module, `guide.external-triage` and
`rule.external-triage` versions, and replace the rule's `guidance` with these
clauses, in this order:

| Clause | Force | Guidance |
| --- | --- | --- |
| `clause.external-triage.classify-before-labelling` (`"replaces": ["rule.external-triage"]`) | mandatory | Use this workflow only for issues and pull requests managed in an external forge. Classify each item as a bug or an enhancement, and state its user-visible problem and next action in English, before changing its labels or status. |
| `clause.external-triage.move-through-mapped-states` | mandatory | Move each item through the triage states needs-triage, needs-info, ready, and wontfix, applying the forge label the repository maps to each state. The repository records that mapping in this guide, outside its setup markers, when it enables external triage. |
| `clause.external-triage.ask-for-an-unmapped-label` | stop-and-ask | Stop and ask the maintainer for the forge label of a triage state the repository has not mapped; never invent a label. |
| `clause.external-triage.needs-info-asks-and-stops` | mandatory | When an item lacks what triage needs, mark it needs-info, ask the reporter specific questions that name each missing fact, and stop triaging it until the reporter answers. |
| `clause.external-triage.disclose-ai-authorship` | mandatory | Start every comment posted on the forge with a sentence stating that an AI agent generated it. |
| `clause.external-triage.record-wontfix-decisions` | mandatory | Record each wontfix decision as a declined Backlog Entry whose reason cites the forge item, and answer a repeated request with that recorded decision instead of deciding it again. |
| `clause.external-triage.pull-request-is-an-issue-with-code` | mandatory | Triage an external pull request as an issue with code attached: classify it and move it through the same states before reviewing or merging its code. |
| `clause.external-triage.route-accepted-work-to-specs` | mandatory | Route accepted work into the repository's local Spec workflow; a forge label or status never stands in for a Task's status. |

task_02, Source Baseline rows. Each corpus entry is
`<!-- source-baseline-entry: <id> -->`, one line `- <text>`, the closing
marker and a blank line. All eight go before the `rule.external-triage` entry
of `corpus/docs/agents/external-triage.md`, which stays. Each `manifest.json`
row goes before the `rule.external-triage` row, with `kind`
`normative-clause`, its force, `carrier` `docs/agents/external-triage.md`,
`structure` null, `start` and `end` `0` and `digest` `""`. The identifiers go
before `"rule.external-triage"` in `index.json` `entryIds`, in the same order.

| Clause | Text |
| --- | --- |
| `classify-before-labelling` | MUST use this workflow only for items in an external forge, and classify each item as a bug or an enhancement, in English, before changing its labels or status. |
| `move-through-mapped-states` | MUST move each item through the states needs-triage, needs-info, ready, and wontfix with the forge labels the repository maps to them. |
| `ask-for-an-unmapped-label` | MUST stop and ask the maintainer for the forge label of an unmapped triage state; never invent a label. |
| `needs-info-asks-and-stops` | MUST mark an item needs-info, ask the reporter specific questions, and stop until the reporter answers. |
| `disclose-ai-authorship` | MUST start every forge comment with a sentence stating that an AI agent generated it. |
| `record-wontfix-decisions` | MUST record each wontfix decision as a declined Backlog Entry and answer a repeated request with it. |
| `pull-request-is-an-issue-with-code` | MUST triage an external pull request as an issue with code attached, through the same states. |
| `route-accepted-work-to-specs` | MUST route accepted work into the local Spec workflow; a forge label or status never stands in for a Task's status. |

The retention disposition of the one removed Source Baseline entry,
`rule.external-triage`, is `replaced` by
`clause.external-triage.classify-before-labelling`, with the same force
(`mandatory`). The eight new entries have no prior entry. No other Source
Baseline entry changes its disposition. No fleet fixture in
`internal/cli/baseline_update_test.go` pins the rule.

task_03, `core.json`: add `"crafting-effective-readmes"` after
`"conventional-commits"` and `"typesafe-ai"` after `"testing-boss"` in
`requiredSkills`. Add these dispatch entries in alphabetical position:

```json
{"skill": "crafting-effective-readmes", "triggers": [{"id": "trigger.core.crafting-effective-readmes", "when": "Writing or revising the repository README."}]}
{"skill": "typesafe-ai", "triggers": [{"id": "trigger.core.typesafe-ai", "when": "Adding programmable semantic judgment (routing, ranking, extraction, verification) or touching TypeSafe/Jev."}]}
```

`typescript.json`: remove `"crafting-effective-readmes"` from
`requiredSkills`, and remove its dispatch with
`trigger.typescript.crafting-effective-readmes`. Raise both module versions.

### The refresh procedure (task_01)

```bash
sha=b3c45a45f1bccd3b33aaecaaa22947d942f2fc02
src="$HOME/dev/skills"
git -C "$src" cat-file -e "$sha^{commit}" || { echo "stop: $src lacks $sha" >&2; exit 1; }
tmp="$(mktemp -d)"
git clone --quiet --no-local "$src" "$tmp/skills"
git -C "$tmp/skills" checkout --quiet --detach "$sha"
git -C "$tmp/skills" remote set-url origin https://github.com/marcioaltoe/skills.git
go run -buildvcs=false ./cmd/roundfix baseline assets sync --source-dir "$tmp/skills/setups" --format text
```

It runs after the seeds, the profile edits and the module edits. Then come the
fixture transform, the test edits and `make baseline-digests`. The managed
refresh runs twice, and the second must report `File changes: 0`. A rerun of
the sync with `--check` must exit `0`.

The fixture transform is Spec 0200's, unchanged. It rebuilds
`manifest.setups`, `managedEntryLedger` and `normalizedOutput.result` of
`internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json` from
the refreshed snapshots, and the regeneration fills every digest:

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
The fixture's `plannedByteSequence` and `fileIdentities` are frozen byte
identities that no test reads, and they stay.

### Existing tests that change

- task_01, `internal/baseline/catalog_test.go`: `setups/go-cli.json` becomes
  `setups/go.json` and `setups/rust-cli.json` becomes `setups/rust.json`; the
  `"setup": "rust-cli"` edit literal becomes `"setup": "rust"`; and the two
  `"rust",` anchors become `"context-workflow", "rust",`.
- task_01, `internal/baseline/assets_sync_test.go`: every `go-cli.json`
  becomes `go.json`, `rust-cli.json` and `rust-cli.txt` become `rust.json` and
  `rust.txt`, `typescript-bun.json` becomes `typescript.json`, and each count
  of `4` setups becomes `3`.
- task_01, `internal/baseline/upstream_skill_names_test.go`:
  `TestTheGoCLITUIProfileTakesTheGoTUISetup` becomes
  `TestTheGoCLITUIProfileTakesTheGoSetup`, wanting `go` and reading
  `setups/go.json`.
- task_01, `internal/baseline/derived_ownership_test.go`: the frozen
  enumeration drops `setups/go-cli.json`, maps `rust-cli.json` and
  `typescript-bun.json` to `rust.json` and `typescript.json`, and adds
  `setups/go.json` in place of the earlier `go-tui.json` addition. Its two
  sample artifacts `assets/setups/go-cli.json` become `assets/setups/go.json`.
- task_01, `internal/baseline/catalog_dispatch_outside_setup_test.go` and
  `internal/cli/baseline_assets_sync_test.go`: each `setups/<retired>.json`
  path literal takes the new name.
- task_02, `internal/baseline/clause_characterization_test.go`: the force
  record gains the eight clauses. `internal/baseline/preservation_test.go`:
  `maintainedSourceBaselineEntries` rises by eight from its value on the
  Task's base.
- task_03 and task_04, `skills/baseline_skill_contract_test.go`:
  `upstreamManagedSkillTreeDigest` takes the value the failing test reports.

### The repository procedure (task_03, task_04)

The repository Verification runs when each Task settles, and
`TestThisRepositoryHoldsEveryRequiredExternalSkill` fails as soon as `core`
requires a skill this repository lacks. task_03 therefore restores the two
skills in the same Task that requires them (steps 1 and 4 to 5). task_04
removes `review` (steps 2 to 4).

1. `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --yes --skills-source-dir "$tmp/skills" --format text`,
   with a clone made as in the refresh procedure. It restores
   `crafting-effective-readmes` and `typesafe-ai`.
2. Build `bin/roundfix`, then run
   `bin/roundfix baseline skills reconcile --repo . --profile go-cli-tui --source marcioaltoe/skills --revision b3c45a45f1bccd3b33aaecaaa22947d942f2fc02 --source-dir "$tmp/skills" --format json`,
   and run it again with `--confirm-plan <planDigest>`. The only edit is
   `remove-lock-entry skills-lock.json [review]`.
3. `git rm -r .agents/skills/review`. No command deletes an installed tree.
4. Make `skills/recommended.txt` the sorted lock keys, and set
   `upstreamManagedSkillTreeDigest` to the value the failing test prints. Both
   Tasks do this, because both change the lock.
5. In `README.md` and `docs/user-guide/usage.md`, the Doctor line becomes
   `skills: ok (43 required: 14 Roundfix-owned, 29 external)`, once
   `go run -buildvcs=false ./cmd/roundfix doctor` prints it.

### API Contracts

1. API Contract: `roundfix baseline assets sync --source-dir <b3c45a4>/setups --check`
   exits `0` and reports `setup-context-driven audit: ok`.
2. API Contract: embedded catalog — the setups are `go`, `rust` and
   `typescript`; built-in profiles keep their identifiers; no entry names a
   removed skill; every setup lists the Roundfix-owned skill.
3. API Contract: rendered `docs/agents/external-triage.md` — eight force-labelled
   clauses. A refresh from the Standard TypeScript Monorepo Source Baseline
   reports `rule.external-triage` as `replaced`.
4. API Contract: rendered `docs/agents/skill-dispatch.md` — no `review` or
   `triage` trigger. `trigger.core.typesafe-ai` and
   `trigger.core.crafting-effective-readmes` are present for every built-in
   profile.
5. API Contract: `roundfix doctor` in this repository — `skills: ok (43
   required: 14 Roundfix-owned, 29 external)`.

### Surface Transcripts

None. No command, flag, output field or exit code changes. The commands above
are exercised, not changed.

## Coverage Map

- Goal 1, Story 1, Core Feature 1, Success Metric 1 → The refresh procedure;
  `setupNameFindings`; API Contract 1 and 2.
- Goal 2, Story 2, Core Feature 2, Success Metric 2 → Fixed texts (task_01);
  `removedSkillFindings`; API Contract 2 and 4.
- Goal 3, Story 3, Core Feature 3, Success Metric 3 → Fixed texts (task_02);
  `externalTriageClauseFindings`, `newExternalTriageAdopter`; API Contract 3.
- Goal 4, Story 4, Core Feature 4, Success Metric 4 → Fixed texts (task_03);
  `coreSkillFindings`; API Contract 4.
- Goal 4 (this repository), Goal 5, Core Feature 5, Success Metric 5 → The
  repository procedure; API Contract 5.
- API Contracts 1–2 → Testing Approach 1; 3 → 2; 4 → 3; 5 → 4.

## Integration Points

- **Spec 0207 (active).** Its composed snapshot is now composed from `go` and
  `typescript`, and its sync procedure reads the pin from `setups/go.json`.
  This Spec edits Spec 0207's artifacts and ADR-0204 only where they name the
  retired setups. Spec 0207 must be delivered after this Spec.
- **Spec 0206 (active).** It edits `core.json`, `typescript.json`, the force
  record and the Source Baseline count, which this Spec also edits. Its
  Verification does not depend on this Spec. Either order works, provided
  each Task raises versions and counts from its own base, as both Specs state.
- **Spec 0195.** Its sync rule keeps the Roundfix-owned entry. Upstream now
  lists `roundfix` in all three setups.
- **Adopters.** The next update records the new setup name, rewrites
  `skill-dispatch.md`, restores the two skills, and for external triage renders
  the clauses. `review` and `triage` trees stay until the adopter deletes them.
- **Upstream.** Read only, at one commit, from the local clone.

## Testing Approach

No test uses the network. Each negative case is its own test.

1. **Renames and removals (task_01).** New
   `internal/baseline/upstream_removed_skills_test.go`:
   `TestNoCatalogEntryNamesASkillRemovedUpstream`,
   `TestARemovedSkillNameInTheCatalogIsReported` (one subtest per field:
   required, dispatch, trigger, bundle, setup),
   `TestEveryBuiltInProfileTakesItsRenamedUpstreamSetup` and
   `TestARetiredOrUnownedSetupIsReported`. The seventeen removed names are
   `review`, `triage`, `resolving-merge-conflicts`, `lesson-learned`,
   `accessibility`, `best-practices`, `performance`, `web-quality-audit`,
   `firecrawl-developer-index`, `firecrawl-monitor`,
   `firecrawl-research-index`, `golang-benchmark`,
   `golang-continuous-integration`, `golang-database`, `golang-modernize`,
   `golang-observability` and `golang-performance`. The profile check covers
   only the three listed profiles and the retired names, so a later composed
   snapshot passes it.
2. **External triage (task_02).** New
   `internal/baseline/external_triage_clauses_test.go`:
   `TestTheExternalTriageGuideStatesItsClausesWithForce` reads the module and
   the formatter golden; `TestAMissingOrUnlabelledExternalTriageClauseIsReported`
   feeds literals; `TestTheExternalTriageRuleIsReplacedForSourceBaselineAdopters`
   expects a ready plan with `replaced` and the target named in the evidence;
   `TestAnUndeclaredExternalTriageReplacementIsRefused` removes `replaces` and
   expects `action_required`, category `classification`, and a message naming
   `rule.external-triage`.
3. **Core skills (task_03).** New
   `internal/baseline/core_skill_requirements_test.go`:
   `TestCoreRequiresAndDispatchesTheTypeSafeAndReadmeSkills` reads the modules,
   the Standard TypeScript Monorepo golden and this repository's
   `docs/agents/skill-dispatch.md`; `TestAMissingCoreSkillTriggerIsReported`
   feeds literals.
4. **This repository (task_03, task_04).** In task_03 the two existing tests
   of `internal/cli/this_repository_skill_set_test.go` fail until the restore,
   and pass after it. In task_04 the file gains `review` and `triage` in the
   names whose lock entry or directory is reported, so the same tests fail
   until the reconcile and the deletion.
5. **Regeneration.** Each Task that edits assets runs `TestFormatterComposition`,
   `TestCatalogCompatibility`, `TestBaselinePlanCharacterization`,
   `TestBaselineCompatibilityCorpus`, the read-only managed refresh, and, for
   task_01 to task_03, `TestCatalogDiagnosticCharacterization`.

## Build Order

1. Renames, the retired setup, `review` and `triage` dropped, refresh to
   `b3c45a4`, task_01 (depends on: none).
2. The external-triage clauses and Source Baseline rows, task_02 (depends on:
   1).
3. `core` requires `typesafe-ai` and `crafting-effective-readmes`, and this
   repository restores them, task_03 (depends on: 2).
4. This repository drops its `review` lock entry and tree, task_04 (depends
   on: 3).
5. Terminal QA, task_05 (depends on: 1, 2, 3, 4).

task_02 and task_03 regenerate the files task_01 regenerates, so the three run
in series. task_04 changes the lock, the recommended list and the digest pin
that task_03 changes, and removes what task_01 dropped.

## Risks & Considerations

- **Upstream moves again.** Every step names the commit and checks it out
  detached. When the local checkout lacks the commit, the Task stops.
- **The sync is a write.** It runs only in task_01, through its recoverable
  transaction. A refused run writes nothing.
- **Hazard: generated files.** No pin, golden, catalog snapshot, Source
  Baseline offset or digest, or parity digest is hand-edited.
- **Hazard: this repository's skills.** Only the files the grant bounds may
  change. A restore that would touch another `.agents/skills/` file stops the
  Task.
- **Hazard: forge comments.** The disclosure clause binds what an Agent
  posts. No Roundfix command posts to a forge.

## Decisions

- Keep profile identifiers and take upstream setup names. See ADR-0206.
- Drop `review` without following it to `code-review`. See ADR-0206.
- Triage functions become clauses. The label map lives in the repository's
  guide, and `wontfix` becomes a declined Backlog Entry. See ADR-0206.
- `crafting-effective-readmes` moves to `core`, beside `typesafe-ai`. See
  ADR-0206.
- Leave the frozen parity fixtures alone, under ADR-0072, and transform only
  the sanctioned `asset-sync.json`.
- Remove this repository's `review` lock entry with `baseline skills
  reconcile`, and delete the tree with `git rm`.
