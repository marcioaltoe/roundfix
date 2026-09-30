---
spec: 0207-a-go-cli-and-a-typescript-monorepo-in-one-baseline
prd: _prd.md
created: 2026-09-30
---

# A Go CLI and a TypeScript monorepo in one Baseline — Technical Spec

## Executive Summary

Four changes, in order. The Frontend Layout Decision arrives first, as an
optional catalog decision with clauses that apply under one of its values, so
an unrecorded layout keeps every adopter's clauses byte-identical. A Setup
Snapshot may then be composed from named component snapshots; the asset sync
writes it and the catalog checks it. The built-in `go-cli-typescript-monorepo`
profile takes the composed setup `go-cli-typescript-bun`, and the Go guide gains
its scope sentence. Last, profile alignment checks that the root Make gate
reaches the Verification parts the profile declares. The trade-off accepted is
decision-specific rendering in Go code for the layout sentence and a textual,
one-Makefile gate check: both are smaller than a new structured decision type
or a Make evaluator, and both leave a gap the Risks section names.

## Project Constraints

- Identifier strategy: applicable — `go-cli-typescript-monorepo` (profile),
  `go-cli-typescript-bun` (Setup Snapshot), `frontend.layout` (decision),
  `rule.frontend.recorded-layout` and `clause.frontend.follow-recorded-layout`,
  `capability.stack.go`, `verification.go` and `verification.workspace`, the
  catalog diagnostics `catalog.decision.optional.invalid`,
  `catalog.clause.applies-when.invalid`, `catalog.setup.composition.invalid`,
  `catalog.setup.composition.conflict`, `catalog.setup.composition.drift` and
  `catalog.profile.verification.invalid`, and the alignment divergence
  `verification.gate.part.missing`, all in the existing forms. No identifier
  is renamed or removed. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — embedded assets, local files and
  Git only; the asset sync reads a local checkout of the upstream skills
  repository and clones it only when that checkout lacks the commit, and no
  Verification command uses the network. The composed profile carries the
  repository-owned HTTP Contract Decision unchanged. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0204 and ADR-0205 govern the
  design. ADR-0191 governs the composed snapshot through its components and
  keeps the asset sync the only writer. ADR-0190 governs the Go guide's scope
  sentence. ADR-0063 is the Frontend Layout Decision's model and ADR-0061
  keeps the composed profile's TypeScript stack opinionated. ADR-0058,
  ADR-0060 and ADR-0099 govern retention: the systems clauses keep their
  identity, and a clause a recorded value turns off is a reasoned rejection.
  ADR-0059 governs the composed profile's formatter declaration. Under
  ADR-0067 custom Baseline Profiles are repository-owned, and they resolve
  without the new decision.
  ADR-0072's words "recorded explicitly as a designed delta" govern the
  composed snapshot in the parity comparison. ADR-0186 governs every new
  sentence. ADR-0193 governs the Prerequisites, ADR-0192 a conflict confined to
  derived paths, ADR-0081, ADR-0149 and ADR-0130 the regeneration and the
  audit, and ADR-0073 and ADR-0103 the Managed Refresh. ADR-0080, ADR-0088,
  ADR-0091, ADR-0093, ADR-0094, ADR-0096, ADR-0104, ADR-0117, ADR-0155,
  ADR-0156, ADR-0166, ADR-0167, ADR-0178 and ADR-0179 bind the gate and the
  checker. ADR-0180, ADR-0181, ADR-0182, ADR-0183, ADR-0184, ADR-0187, ADR-0189,
  ADR-0194, ADR-0195, ADR-0196, ADR-0197, ADR-0198 and ADR-0199 are present
  and do not apply; the PRD's row gives each reason. ADR-0074 cites
  ADR-0067, ADR-0097 cites ADR-0080, ADR-0168 and ADR-0176 cite ADR-0093,
  ADR-0173 cites ADR-0073, and ADR-0177 is reached only through ADR-0173; none
  applies, because this Spec changes none of their behaviors. All hold.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30, recorded in [_authorization.md](_authorization.md); bounded
  files: `internal/baseline/assets/decisions.json`,
  `internal/baseline/assets/modules/frontend.json`,
  `internal/baseline/assets/modules/go.json`,
  `internal/baseline/assets/templates/index.json`,
  `internal/baseline/assets/templates/guides/frontend.md`,
  `internal/baseline/assets/templates/guides/go.md`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/baseline/assets/profiles/go-cli-typescript-monorepo.json`,
  `internal/baseline/assets/setups/go-cli-typescript-bun.json`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/frontend.md`,
  `internal/baseline/derived_ownership_test.go`,
  `internal/cli/baseline_human_test.go`,
  `docs/agents/setup-context.json`. Sanctioned regeneration:
  `make baseline-digests`. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`.

## System Architecture

| Concern | Source of truth | Derived by |
| --- | --- | --- |
| Frontend Layout Decision | `internal/baseline/assets/decisions.json` | hand edit under the grant |
| Clauses gated by a decision value | `internal/baseline/assets/modules/frontend.json` | hand edit under the grant |
| Clause selection, rendering, retention | `internal/baseline/plan.go` | — |
| Layout sentence | `internal/baseline/project_decision_render.go` | — |
| Decision and clause validation | `internal/baseline/catalog_load.go`, `internal/baseline/catalog_validate.go` | — |
| Composed Setup Snapshot | `internal/baseline/assets/setups/go-cli-typescript-bun.json` | `roundfix baseline assets sync --source-dir <checkout>/setups` |
| Composition rule | new `internal/baseline/setup_composition.go` | — |
| Composed profile | `internal/baseline/assets/profiles/go-cli-typescript-monorepo.json` | hand edit under the grant |
| Go guide scope sentence | `internal/baseline/assets/templates/guides/go.md` | — |
| Gate parts check | new `internal/baseline/verification_gate_parts.go`, called from `internal/baseline/profile_alignment.go` | — |
| Goldens, pins, catalog snapshots | the assets | `make baseline-digests` |
| This repository's guides and Setup Manifest | the assets | `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text` |

### Measured facts

Measured on 2026-09-30 in the worktree at `5f182757` (main), before Spec 0200
landed, and against the upstream lists at `a4e18e4` in `~/dev/skills`.

- A profile's setup is one identifier. `validateProfiles` reports
  `catalog.profile.skill.outside-setup` for every required skill of a selected
  module that the setup does not list, and Spec 0200 adds the same check for
  dispatched and bundled skills. `skills_restore.go`, reconcile, Doctor and
  `catalogIdentityBytes` decode a setup file's `skills` list directly, and
  `skills/baseline_skill_contract_test.go` reads every file under
  `internal/baseline/assets/setups/`.
- The asset sync lists every `setups/*.json` file and builds each from the
  upstream file of the same name; a snapshot with no upstream file fails with
  "Canonical setup source file is missing." Activation bundles are carried
  over from the current snapshot, not from upstream.
- Upstream at `a4e18e4`: `go-cli` lists 46 skills, `typescript-bun` 110, and
  their union 125; the 15 `go-cli` entries not in `typescript-bun` are Go
  skills. No upstream list combines them. `go-cli` omits `roundfix`, which
  Spec 0195's sync rule keeps as an owned entry. In today's snapshots the 30
  skills common to `go-cli` and `typescript-bun` are byte-identical entries,
  and `typescript-bun`'s ten activation bundles include `go-cli`'s two with the
  same skills.
- `resolveManagedArtifacts` renders an artifact named by any decision
  effect's `renderBindings` or `includeArtifacts` only when that effect
  matches a recorded value. A binding on an unrecorded decision would drop
  the whole frontend guide.
- `normalizePlanDecisions` reports every selected decision without a value as
  missing; `baseline update` then stops with `action_required` unless run with
  `--adopt-suggested`, and the interactive first adoption prompts every entry
  decision of the profile in order.
- `classifySourceClauseTransition` compares clause identifiers and
  enforcement only. `TestStandardTypeScriptStructuralClauseRetention` also
  compares the text of `clause.frontend.organize-by-system` and
  `clause.frontend.public-system-boundary` with the Source Baseline corpus.
  `validateSourceBaselineCoverage` requires every clause of a required rule of
  the Standard TypeScript Monorepo Profile to have a Source Baseline row, so a
  new clause cannot join `rule.frontend.user-visible-behavior`.
- Adopter mirrors: conexus, fiscus, fluxus and vortex render both systems
  clauses; gss (repository-owned profile `gss-apps-monorepo`) and tax-poc hold
  frontend guides from an older catalog without them. No Setup Manifest
  records a frontend layout.
- `resolveVerificationProjection` resolves `make <target>` against a Makefile
  target line and `bun run <script>` against `package.json`; no code reads a
  target's prerequisites or recipes. The profile's `verification` entries are
  not validated beyond their use.
- `catalog_test.go` pins the built-in profile list, and
  `internal/baseline/derived_ownership_test.go` pins the sanctioned outputs,
  which include every setup file.

## Implementation Design

### Interfaces

```go
// internal/baseline/plan.go
// An optional decision (decision field "optional": true) is never reported
// missing by normalizePlanDecisions and is never recorded unless answered.
func decisionOptional(declaration document) bool

// clauseApplies reports whether a clause is selected under the recorded
// decision values. A clause without "appliesWhen" always applies. With
// {"decision": id, "equals": v}: the recorded value equals v, or the decision
// is unrecorded and its "default" equals v.
func clauseApplies(catalog *Catalog, clause document, recorded map[string]any) bool

// selectedClauseEnforcement and artifactRenderValues gain the recorded values
// and skip a clause that does not apply. classifySourceClauseTransition reads
// the recorded values from the target Setup Manifest; a prior clause absent
// from the selection whose catalog clause exists in an active module but does
// not apply under a recorded value becomes:
//   RetentionEvidence{Disposition: "reasoned-rejection",
//     Targets: []string{<decision id>},
//     Reason: "The repository recorded <id> = <value>; this clause applies only when it is <equals>."}
func selectedClauseEnforcement(catalog *Catalog, activeModules []string, recorded map[string]any) map[string]string
```

```go
// internal/baseline/plan.go — resolveManagedArtifacts
// An optional decision's renderBindings do not make their artifact
// decision-controlled. While such a decision is unrecorded its token renders
// renderUnrecordedProjectDecision(id, declaration).
// internal/baseline/project_decision_render.go
const frontendLayoutDecisionID = "frontend.layout"
func renderUnrecordedProjectDecision(decisionID string, declaration document) (string, error)
// renderProjectDecision gains the frontend.layout case (Fixed texts).
```

```go
// internal/baseline/setup_composition.go
// composeSetupSnapshot returns the union of the components' skills and
// activation bundles in component order. A later entry equal to an earlier
// one of the same name (or bundle id) is dropped; an unequal one is an error.
func composeSetupSnapshot(id string, components []document) (document, error)
// internal/baseline/assets_sync.go: assetsSyncSource gains
//   Setups []string `json:"setups,omitempty"`
// syncAssets skips a snapshot whose source.type is "composed" when it builds
// from upstream files, then composes each composed snapshot from the built
// (or unchanged) component snapshots and plans it like any other snapshot.
```

```go
// internal/baseline/verification_gate_parts.go
// makeTargetReach reads one Makefile and returns every target reachable from
// target through prerequisites and through recipe lines that run
// "$(MAKE) <t>", "${MAKE} <t>" or "make <t>", and every recipe line of those
// targets. Each target is visited once; no variable other than MAKE is expanded.
func makeTargetReach(makefile []byte, target string) (targets map[string]struct{}, recipes []string)
// gatePartReached: a part "make <t>" is reached when <t> is a reached target;
// any other part is reached when a reached recipe line, with a leading "@",
// "-" and "rtk " removed, contains the part's command with "rtk " removed.
func gatePartReached(command string, targets map[string]struct{}, recipes []string) bool
```

`resolveVerificationProjection` calls them after it resolves the projections,
only when the selected `verification.gate` is `make <target>` declared in a
Makefile, once per built-in profile entry with `"partOfGate": true`, and
appends for each unreached part:

```go
ProfileDivergence{Code: "verification.gate.part.missing", ID: <entry id>,
	Requirement: CapabilityRecommended, Blocking: false,
	Message: fmt.Sprintf("the repository gate %q does not run %q", gate, part),
	NextAction: "make the repository gate run the named Verification, or map its role to a command the gate runs"}
```

### Data Models

- **Decision.** A decision may carry `"optional": true`. The catalog loader
  accepts the field. `catalog.decision.optional.invalid` is reported when it
  is not a boolean, when an optional decision has no valid `default`, or when
  one of its effects holds anything other than `renderBindings` under
  `{"present": true}`.
- **Clause.** A clause may carry `"appliesWhen": {"decision": <id>, "equals": <value>}`.
  `catalog.clause.applies-when.invalid` is reported, with the reason as its
  detail, when the object has other fields, the decision is unknown, not an
  optional `enum` decision, or the value is not one of its `values`.
- **Setup Snapshot.** `source` may be `{"type": "composed", "setups": [<ids>]}`.
  `catalog.setup.composition.invalid`: fewer than two components, a
  duplicate, an unknown component or a composed component.
  `catalog.setup.composition.conflict`: two components give one skill name or
  one bundle identifier different content.
  `catalog.setup.composition.drift`: `skills` or `activationBundles` differ
  from `composeSetupSnapshot`'s result. The digest rule is unchanged.
- **Profile verification entry.** An entry may carry `"partOfGate": true`;
  `catalog.profile.verification.invalid` is reported when the field is not a
  boolean.
- **Standard TypeScript Monorepo Profile.** `entryDecisions` gains
  `frontend.layout` after `http.contract` and before `auth.provider`.
  `architecture.frontend` becomes
  `{"layoutDecision": "frontend.layout", "suggested": {"organization": "systems", "publicBoundary": "public system boundary", "internalImports": "direct"}}`.
- **Composed profile** `go-cli-typescript-monorepo`, title
  "Go CLI with TypeScript Monorepo", `setup` `go-cli-typescript-bun`,
  `formatter` `{"kind": "none"}`:
  - `modules`: `core`, `context-workflow`, `go`, `cli-surface`, `typescript`,
    `bun`, `monorepo`, `backend`, `frontend`, `autonomous-work`,
    `spec-workflow`, `external-triage`, `secondbrain`, `repository-extension`.
  - `entryDecisions`: the Standard TypeScript Monorepo Profile's list after
    this Spec's change, in the same order.
  - `requiredRules`: the Standard TypeScript Monorepo Profile's list with
    `rule.go.stdlib-first`, `rule.go.context-errors`, `rule.go.observable-tests`,
    `rule.cli.output-contract` and `rule.cli.non-interactive` after
    `rule.context.docs-layout`.
  - `stack`: `Go` followed by the Standard TypeScript Monorepo stack.
    `workspaces`, `optionalModules`, `architecture`, `httpContract`,
    `capabilitySets` and `activationBundles`: equal to the Standard TypeScript
    Monorepo Profile's.
  - `capabilities`: the Standard TypeScript Monorepo Profile's, plus
    `{"id": "capability.stack.go", "title": "Go", "category": "stack", "strength": "required", "probe": {"kind": "declared-file", "paths": ["go.mod"], "contains": "module "}}`
    in identifier order.
  - `verification`: the Standard TypeScript Monorepo Profile's `format`,
    `lint`, `test` and `build` entries;
    `{"id": "verification.workspace", "kind": "workspace", "tool": "Bun", "command": "bun run verify", "partOfGate": true}`;
    `{"id": "verification.go", "kind": "go", "tool": "Go", "command": "make verify-go", "partOfGate": true}`;
    and the `incremental` entry.
- **Composed setup** `go-cli-typescript-bun`: `source`
  `{"type": "composed", "setups": ["go-cli", "typescript-bun"]}`; the rest is
  written by the asset sync.

### Fixed texts

Replace strings and numbers in place in module, profile and template files;
never pass them through a JSON encoder. Raise each version by one from the
value on the Task's base.

- task_01, `decisions.json`, appended after `http.contract`:

```json
{
  "id": "frontend.layout",
  "version": 1,
  "type": "enum",
  "values": ["systems", "repository-defined"],
  "default": "systems",
  "optional": true,
  "summary": "The frontend layout the repository records. The systems layout is suggested and applies until one is recorded.",
  "effects": [
    {
      "when": {"present": true},
      "renderBindings": [
        {"artifact": "guide.frontend", "template": "template.guide.frontend", "token": "frontend.layout"}
      ]
    }
  ]
}
```

- task_01, `frontend.json`: `clause.frontend.organize-by-system` and
  `clause.frontend.public-system-boundary` each gain
  `"appliesWhen": {"decision": "frontend.layout", "equals": "systems"}`;
  their `guidance` stays byte-identical. A new rule after
  `rule.frontend.user-visible-behavior`, listed second in `guide.frontend`'s
  `rules`:

```json
{
  "id": "rule.frontend.recorded-layout",
  "version": 1,
  "coverage": ["coverage.frontend"],
  "clauses": [
    {
      "id": "clause.frontend.follow-recorded-layout",
      "enforcement": "mandatory",
      "appliesWhen": {"decision": "frontend.layout", "equals": "repository-defined"},
      "guidance": "Organize frontend feature code by the layout the repository's own rules state, and update those rules in the same change that alters the layout."
    }
  ]
}
```

  Versions raised: the module, `guide.frontend`, `rule.frontend.user-visible-behavior`.
- task_01, `templates/guides/frontend.md`: a paragraph `{{frontend.layout}}`
  between the scope paragraph and `{{artifact.rules}}`;
  `template.guide.frontend` gains the token `frontend.layout` and its version
  is raised.
- task_01, the sentences the token renders:
  - unrecorded: ``No frontend layout is recorded. The suggested `systems` layout applies until the repository records one.``
    (the value is the decision's `default`);
  - `systems`: ``The repository records the `systems` frontend layout.``;
  - `repository-defined`: `The repository records its own frontend layout, stated in its repository-owned rules.`
- task_01, `docs/user-guide/context-driven-development.md`: the suggested
  values table gains `| Frontend layout | `systems` (applies while none is recorded) |`
  after the HTTP contract row, and one paragraph after the table states that
  the frontend layout is optional, that `repository-defined` binds the
  repository's own rules, and that a recorded value is kept.
- task_03, `templates/guides/go.md`: between the heading and
  `{{artifact.rules}}`:

```markdown
These rules govern the repository's Go module: its commands, packages and
tests. Code in another language follows its own guide.
```

  `template.guide.go`, `guide.go` and the `go` module's versions are raised.
- task_03, `docs/user-guide/context-driven-development.md`: the Profiles
  paragraph names `go-cli-typescript-monorepo` and states that it combines the
  Go CLI and the Standard TypeScript Monorepo modules on a composed setup.
- task_04, the same guide: one paragraph under Profiles states that the
  composed profile's root gate is expected to run `make verify-go` and
  `bun run verify`, and that alignment reports
  `verification.gate.part.missing` for a part the gate does not reach.

### The composed setup procedure (task_02)

After the code change, create the seed
`{"schemaVersion": "setup-context-driven/setup-snapshot/0.0.1", "id": "go-cli-typescript-bun", "version": "0.0.1", "source": {"type": "composed", "setups": ["go-cli", "typescript-bun"]}, "digest": "", "skills": []}`,
then run the sync against the pinned upstream commit:

```bash
sha="$(jq -r .source.ref internal/baseline/assets/setups/go-cli.json)"
src="$HOME/dev/skills"
git -C "$src" cat-file -e "$sha^{commit}" 2>/dev/null || src=https://github.com/marcioaltoe/skills.git
tmp="$(mktemp -d)"
git clone --quiet --no-local "$src" "$tmp/skills"
git -C "$tmp/skills" checkout --quiet --detach "$sha"
git -C "$tmp/skills" remote set-url origin https://github.com/marcioaltoe/skills.git
go run -buildvcs=false ./cmd/roundfix baseline assets sync --source-dir "$tmp/skills/setups" --format text
```

The run must change only the composed file. Then `make baseline-digests` and
the Managed Refresh, twice.

### API Contracts

1. API Contract: embedded catalog — `frontend.layout` is an optional enum
   decision selected by the Standard TypeScript Monorepo and composed
   profiles; the catalog refuses an invalid optional decision or clause gate
   with `catalog.decision.optional.invalid` or
   `catalog.clause.applies-when.invalid`.
2. API Contract: `roundfix baseline` and `roundfix baseline plan` — an
   unrecorded `frontend.layout` is not a missing decision; `--decision
   frontend.layout=<value>` records it; the interactive first adoption
   prompts it with `systems` preselected; the rendered frontend guide carries
   the sentence and clauses of Fixed texts for each state.
3. API Contract: `roundfix baseline update` — a Setup Manifest with no
   `frontend.layout` reports no new decision and keeps both systems clauses;
   a recorded value is kept; a plan that turns a managed clause off through a
   recorded value lists it under retention as `reasoned-rejection`.
4. API Contract: `roundfix baseline assets sync --source-dir <dir> [--check]`
   — composed snapshots are written, or reported as drift under `--check`, in
   the same run as their components; a composition conflict exits `2` with
   the invalid-assets finding and writes nothing.
5. API Contract: embedded catalog — `go-cli-typescript-monorepo` is a built-in
   profile on `go-cli-typescript-bun`; the catalog refuses a composed snapshot
   with `catalog.setup.composition.invalid`, `.conflict` or `.drift`.
6. API Contract: `roundfix baseline plan` alignment — with a built-in profile
   that declares gate parts and a selected `make <target>` gate, each part the
   target does not reach is a non-blocking `verification.gate.part.missing`
   divergence; `catalog.profile.verification.invalid` refuses a non-boolean
   `partOfGate`.

## Coverage Map

- Goal 1 → Data Models (composed profile); API Contract 5; Testing Approach 3.
- Goal 2 → Interfaces (`composeSetupSnapshot`, sync); API Contract 4; Testing Approach 2.
- Goal 3 → Interfaces (`makeTargetReach`, `gatePartReached`); API Contract 6; Testing Approach 4.
- Goal 4 → Interfaces (`clauseApplies`, optional decision); API Contracts 1 and 2; Testing Approach 1.
- Goal 5 → Interfaces (`clauseApplies`, retention); API Contract 3; Testing Approach 1.
- Story 1 → Data Models (composed profile); API Contract 5; Testing Approach 3.
- Story 2 → Interfaces (`composeSetupSnapshot`); API Contract 4; Testing Approach 2.
- Story 3 → Fixed texts (Go guide); Testing Approach 3.
- Story 4 → Interfaces (gate parts); API Contract 6; Testing Approach 4.
- Story 5 → Fixed texts (frontend); API Contract 2; Testing Approach 1.
- Story 6 → Interfaces (optional decision); API Contract 3; Testing Approach 1.
- Core Feature 1 → Interfaces (`composeSetupSnapshot`); Data Models (Setup Snapshot); API Contracts 4 and 5.
- Core Feature 2 → Data Models (composed profile); API Contract 5.
- Core Feature 3 → Fixed texts (Go guide).
- Core Feature 4 → Interfaces (gate parts); Data Models (profile verification entry); API Contract 6.
- Core Feature 5 → Interfaces (`clauseApplies`, `renderUnrecordedProjectDecision`); Fixed texts (frontend); API Contracts 1 and 2.
- Core Feature 6 → Interfaces (`selectedClauseEnforcement`, retention); API Contract 3.
- Success Metric 1 → Testing Approach 2 and 3.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → Testing Approach 3.
- Success Metric 4 → Testing Approach 4.
- Success Metric 5 → Testing Approach 1.
- Success Metric 6 → Testing Approach 1.
- Success Metric 7 → Testing Approach 5.
- API Contract 1 → Data Models (Decision, Clause).
- API Contract 2 → Interfaces (optional decision, render).
- API Contract 3 → Interfaces (retention).
- API Contract 4 → Interfaces (sync).
- API Contract 5 → Data Models (Setup Snapshot, composed profile).
- API Contract 6 → Interfaces (gate parts).

## Integration Points

- **Spec 0200.** Its refreshed `go-cli` and `typescript-bun` snapshots are the
  components; its asset sync, which validates the catalog it produces, is the
  code task_02 extends; its dispatch check must pass for the composed profile.
  Its Go CLI/TUI profile moves to `go-tui`, and `go-cli` stays for this
  composition.
- **Spec 0195.** Its sync rule keeps the owned `roundfix` entry in `go-cli`,
  which the composition carries.
- **Adopters.** A Standard TypeScript Monorepo adopter's next update adds the
  suggestion sentence to its frontend guide and changes no clause. A Go
  CLI/TUI adopter's Go guide gains its scope sentence. Repository-owned
  profiles need no change.
- **Upstream skills.** Read only, at the commit the component snapshots pin.

## Testing Approach

No test uses the network. Each negative case is its own test.

1. **Frontend layout (task_01).** New
   `internal/baseline/frontend_layout_decision_test.go`, on the embedded
   catalog and `newAlignedTypeScriptRepository`:
   `TestAnUnrecordedFrontendLayoutStatesTheSuggestionAndKeepsTheSystemsClauses`,
   `TestARecordedSystemsLayoutKeepsTheSystemsClauses`,
   `TestARecordedRepositoryDefinedLayoutRendersOnlyItsOwnClause`,
   `TestAnUnrecordedOptionalDecisionIsNotMissing`,
   `TestAnOptionalDecisionWithoutADefaultIsRefused`,
   `TestAClauseGateOnARequiredDecisionIsRefused`,
   `TestAClauseARecordedLayoutTurnsOffIsAReasonedRejection`,
   `TestAClauseMissingWithoutARecordedDecisionIsStillUnaccounted` and
   `TestTheProfileStatesTheLayoutSuggestionTheCatalogDefaults`. The refusal
   tests load an overlay of the embedded assets. Existing tests that must
   still pass unchanged: `TestStandardTypeScriptStructuralClauseRetention`,
   `TestNoTwoBaselineClausesShareText`. `TestBaselineClauseForceIsCharacterized`
   gains the new clause in `clause_characterization_test.go`.
   `TestHumanBaselineDecisionDefaults` gains the row
   `{name: "frontend layout", id: "frontend.layout", want: "systems"}`.
2. **Composition (task_02).** New
   `internal/baseline/setup_composition_test.go`:
   `TestTheComposedSetupIsTheUnionOfItsComponents`,
   `TestAComposedSetupThatDriftsFromItsComponentsIsRefused`,
   `TestAComposedSetupWithAnUnknownOrComposedComponentIsRefused`,
   `TestComponentsThatDisagreeOnASkillAreAConflict`,
   `TestAnAssetSyncRewritesTheComposedSetupWithItsComponents` and
   `TestAnAssetSyncWithNoSourceChangeLeavesTheComposedSetup`. The existing
   asset-sync tests whose counts include every setup, and the parity
   comparison, which skips composed snapshots as the designed delta, are
   updated in `assets_sync_test.go`; the parity fixture is not.
3. **Composed profile (task_03).** New
   `internal/baseline/composed_profile_test.go`:
   `TestTheComposedProfileTakesTheComposedSetup`,
   `TestTheComposedProfileRendersEveryGuideWithNoRepeatedClause`,
   `TestTheGoGuideNamesWhatItGoverns`,
   `TestARepeatedRenderedClauseIsReported` and
   `TestTheComposedProfilePlanConverges`. Spec 0200's
   `TestEveryBuiltInProfileSetupListsEverySkillItsGuidesName` must pass for
   the new profile by name. `catalog_test.go` names the new profile.
4. **Gate parts (task_04).** New
   `internal/baseline/verification_gate_parts_test.go`:
   `TestAGateThatReachesBothPartsReportsNoDivergence`,
   `TestAGateThatSkipsAPartReportsItOnce`,
   `TestAPartReachedThroughARecipeInvocationCounts`,
   `TestAGateThatIsNotAMakeTargetIsNotChecked`,
   `TestMakeTargetReachStopsAtACycle` and
   `TestANonBooleanPartOfGateIsRefused`.
5. **Convergence (every Task).** `make baseline-digests`, then the Managed
   Refresh of this repository twice; the second reports `File changes: 0`.
   The Source Baseline corpus, the retention transitions and
   `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json` stay
   byte-identical.

## Build Order

1. The Frontend Layout Decision: optional decisions, clause gates, the
   layout sentence, retention, the frontend module and the Standard
   TypeScript Monorepo Profile.
2. The composed Setup Snapshot: composition, validation, the asset sync and
   the `go-cli-typescript-bun` snapshot (depends on: 1, because both change
   `catalog_validate.go` and the catalog snapshots).
3. The composed profile and the Go guide's scope sentence (depends on: 1, 2).
4. The root gate's parts (depends on: 3).
5. QA (depends on: 1, 2, 3, 4).

## Risks & Considerations

- **The gate check is textual.** It reads one Makefile; an `include`d file, a
  shell script, a variable holding the command, or `bun --cwd x run verify`
  is not followed. The divergence stays advisory for that reason.
- **The layout sentence lives in Go code**, as the HTTP contract's does; a
  second optional decision needs its own render case.
- **Union order.** Skills and bundles follow component order, so the composed
  file's bundle order differs from `typescript-bun`'s. Nothing reads bundle
  order.
- **Spec 0200 moves the sync and its tests.** task_02 builds on the code as
  Spec 0200 leaves it; the counts in `assets_sync_test.go` are whatever that
  base asserts plus the composed snapshot.
- **`repository-defined` cannot be checked.** The Baseline binds the
  repository's own rules and checks no directory.
- **Human prompt order.** The interactive first adoption of a profile with the
  frontend module gains one question; `baseline_human_test.go` scripts that
  include such an adoption gain one answer.

## Decisions

- **Optional, not required.** See ADR-0205; it removes the migration and
  keeps every caller that passes decisions working.
- **Gate clauses in place.** The systems clauses keep their bytes and carrier.
- **The rejection names the decision.** A clause turned off by a recorded
  value is a reasoned rejection with the decision as its target.
- **Materialized composition written by the sync.** See ADR-0204.
- **Component order union with equality de-duplication.** A conflict is an
  error, never a choice.
- **`go-cli`, not `go-tui`.** The planned repository has a CLI and no TUI.
- **Formatter `none`.** No fixture set proves a formatter for the composed
  profile.
- **Parts in the profile, reach in the Makefile.** The profile states what the
  gate must run; alignment reads the one Makefile it already reads.
