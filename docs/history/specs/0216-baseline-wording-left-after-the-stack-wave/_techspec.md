---
spec: 0216-baseline-wording-left-after-the-stack-wave
prd: _prd.md
created: 2026-10-02
---

# Baseline wording left after the stack wave — Technical Spec

## Executive Summary

Two Baseline modules, two profiles, two guide templates and one render
function change; no command does. The backend bucket clause takes a new
identity that declares `replaces`, the two built-in TypeScript profiles bind
each declared workspace to its guide and a new `workspace.location` token
renders the path, and core gains one skill-dispatch clause for skills only a
person can start. Each new clause of a Standard TypeScript Monorepo module
gains a Source Baseline row by hand, the one thing the sanctioned regeneration
cannot create. The trade-off accepted is a general clause instead of a
per-trigger marker for person-only skills: it covers every such skill without
catalog data the Baseline cannot verify offline, at the cost of naming none of
them. Every change below was applied in a disposable clone of `6364d3c9` on
2026-10-02, one commit per Task; after each commit `make baseline-digests`
converged on its second run, the second Managed Refresh reported
`File changes: 0`, and the Baseline, CLI and skills test packages passed.

## Project Constraints

- Identifier strategy: not applicable — new Baseline clause identifiers in
  the existing `clause.<module>.<name>` form and one render token; no domain
  identifier. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — Baseline assets, one render
  function and local tests only; no credential and no network call; the
  `auth.provider` decision is unchanged. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — the PRD's row lists every active ADR
  with its receipt or the reason it does not apply. ADR-0222 (this Spec)
  governs the design, ADR-0222: "One core clause covers every such skill,
  present or future". ADR-0190's deferred alternative is the workspace
  rendering below, ADR-0190: "Rendering each workspace path into the guide
  needs a profile field and a render token that do not exist". Retention is a
  comparison of identities, ADR-0099: "retention accounting compares prior
  managed clauses against the current catalog's clauses". Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30 ("Autorizar os dois") and 2026-10-01 ("Tudo, de A a F"). Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0216-baseline-wording-left-after-the-stack-wave/_authorization.md`;
  bounded files: `internal/baseline/assets/modules/backend.json`,
  `internal/baseline/assets/modules/core.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/baseline/assets/profiles/go-cli-typescript-monorepo.json`,
  `internal/baseline/assets/templates/index.json`,
  `internal/baseline/assets/templates/guides/backend.md`,
  `internal/baseline/assets/templates/guides/frontend.md`,
  `internal/baseline/assets/source-baselines/index.json`,
  `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/baseline.json`,
  `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/manifest.json`,
  `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/backend.md`,
  `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/agent-instructions.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/backend.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/frontend.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`,
  `internal/baseline/plan_test.go`, `docs/agents/skill-dispatch.md`,
  `docs/agents/setup-context.json`.

## System Architecture

| Concern | Source of truth | Derived by |
| --- | --- | --- |
| Clause text, force, versions | `internal/baseline/assets/modules/backend.json`, `core.json` | hand edit under the grant |
| Workspace binding | `profiles/standard-typescript-monorepo.json`, `profiles/go-cli-typescript-monorepo.json` `workspaces[].guide` | hand edit under the grant |
| Scope sentence | `templates/guides/backend.md`, `frontend.md`, `templates/index.json` | hand edit under the grant |
| Workspace rendering | `internal/baseline/plan.go` | code |
| Source Baseline rows | corpus files, `manifest.json` row, `index.json` `entryIds` | hand edit: text, identity, force, carrier; the regeneration fills offsets, digests and `baseline.json` |
| Goldens, digest pin, catalog snapshots, plan goldens | the sources above | `make baseline-digests` |
| This repository's guides and Setup Manifest | the sources above | `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text` |

A module or profile edit keeps the file's formatting: change objects in place
and raise version numbers, never re-encode the file. Each Task edits its
sources and tests, runs `make baseline-digests` twice (the second reports
`"changed":false`), then the Managed Refresh twice (the second reports
`File changes: 0`).

### Outputs measured in the rehearsal

- task_01 rewrites `catalog.digest`, `catalog.normalized.json`,
  `catalog.diagnostics.golden.json`, the four plan goldens under
  `internal/baseline/testdata/plan-characterization/`
  (`advisory-only-divergences`, `clean-adoption`,
  `idempotent-replan-after-verified-apply`,
  `same-baseline-changed-profile-and-catalog-digests`), the Standard
  TypeScript Monorepo digest pin in its profile, the golden `backend.md`, the
  Source Baseline `baseline.json`, `manifest.json` and `index.json`, and this
  repository's `setup-context.json`.
- task_02 rewrites the same catalog files and plan goldens, the digest pin,
  the goldens `backend.md` and `frontend.md`, and `setup-context.json`.
- task_03 rewrites what task_01 rewrites, except that its golden is
  `skill-dispatch.md`, and the Managed Refresh also rewrites this repository's
  `docs/agents/skill-dispatch.md`.

`make verify` exited 0 on the clone with all three Tasks applied.

## Implementation Design

### Exact texts

- task_01, `clause.backend.prohibit-generic-buckets` — prohibited —
  `"replaces": ["clause.backend.prohibit-generic-layers"]` — "Do not organize
  backend code into generic `modules` or `services` buckets in place of the
  domain, application, and infrastructure layers. A domain service that lives
  in the domain layer is not such a bucket." It takes the place of
  `clause.backend.prohibit-generic-layers` in `rule.backend.boundary-contracts`.
- task_02, backend template: the sentence becomes "These rules govern the
  repository's TypeScript backend workspace{{workspace.location}}. A service
  or" with the line break where it is today; frontend: "These rules govern the
  repository's web frontend workspace{{workspace.location}}. A terminal". Each
  workspace entry of both built-in TypeScript profiles gains
  `"guide": "guide.backend"` or `"guide": "guide.frontend"`.
- task_03, `clause.core.ask-the-person-to-start-a-person-only-skill` —
  mandatory — "When a matching skill can be started only by a person, because
  its metadata turns off model invocation, ask the person to run it instead of
  activating it or carrying out its workflow yourself." It is the last clause
  of `rule.core.skill-dispatch`.

### Interfaces

```go
// plan.go — default for the two guides, overwritten when a workspace binds them.
case "guide.frontend":
	values["workspace.location"] = ""

// profileWorkspaceLocations returns, per guide ID, " at `p`" or
// " at `a` and `b`" (paths sorted), from the built-in Profile's
// workspaces whose guide is set and whose path is safe and relative.
func profileWorkspaceLocations(profile document) map[string]string
```

`resolveManagedArtifacts` computes the map once from
`catalog.profiles[profile.ID]` and sets the token for the artifact it renders.
A repository-owned profile is not in `catalog.profiles`, so the map is empty.

### Data Models

- Profile `workspaces[]` entries gain an optional `guide` string, the managed
  identifier of a supporting guide one of the profile's modules ships.
- `templates/index.json`: `template.guide.backend` 2 → 3 and
  `template.guide.frontend` 3 → 4, each listing `workspace.location`.

### Source Baseline rows

Each row is a corpus entry `<!-- source-baseline-entry: <id> -->`, one line
`- <text>` and the closing marker, after the named anchor with a blank line
between entries; a `manifest.json` row after the anchor's row with `kind`
`normative-clause`, the force, the carrier, `structure` null and placeholder
offsets and digest; and the identifier after the anchor in `index.json`
`entryIds`. The regeneration fills the rest.

| Task | Clause | After | Corpus file and carrier | Text |
| --- | --- | --- | --- | --- |
| 01 | `clause.backend.prohibit-generic-buckets` | `clause.backend.prohibit-generic-layers` | `backend.md` | MUST NOT organize backend code into generic `modules` or `services` buckets in place of the domain, application, and infrastructure layers. A domain service that lives in the domain layer is not such a bucket. |
| 03 | `clause.core.ask-the-person-to-start-a-person-only-skill` | `clause.core.prohibit-editing-vendored-skills` | `agent-instructions.md` | MUST ask the person to run a matching skill that only a person can start, instead of activating it or carrying out its workflow yourself. |

The old bucket row stays: it records what adopters hold, and the classifier
reads it to report the replacement. Its retention disposition is `replaced`.
The task_03 clause sits beside its rule's siblings, whose rows also use the
`agent-instructions.md` carrier.

### Version changes

Raise each by one from its value on the Task's starting main. task_01:
`backend` (5), `guide.backend` (5), `rule.backend.boundary-contracts` (4).
task_02: `template.guide.backend` (2), `template.guide.frontend` (3). task_03:
`core` (17), `rule.core.skill-dispatch` (5), `guide.skill-dispatch` (4).

### Existing tests that change

- `internal/baseline/clause_characterization_test.go` (task_01, task_03): the
  force record swaps the bucket clause's key (task_01) and gains the
  person-only clause (task_03). Neither key is the longest, so `gofmt`
  realigns nothing.
- `internal/baseline/preservation_test.go` (task_01, task_03): the maintained
  Source Baseline entry count rises by one in each, from 161.
- `internal/baseline/plan_test.go` (task_01):
  `TestStandardTypeScriptStructuralClauseRetention` expects the bucket clause
  `replaced` by its successor and compares the successor's guidance with the
  successor's own Source Baseline row. The two replacements become one map.
- `internal/cli/baseline_update_test.go` (task_01): the fleet structural
  fixture names the new identity and line.
- `internal/baseline/stack_scope_and_http_default_test.go` (task_02):
  `TestTheBackendAndFrontendGuidesSayWhatTheyGovern` expects the path in both
  sentences.

### API Contracts

1. API Contract: rendered `docs/agents/backend.md` of a Standard TypeScript
   Monorepo or composed-profile plan — the scope sentence names
   `packages/backend` and the bucket clause carries the task_01 text with
   `prohibited`; a repository-owned profile's backend guide renders the scope
   sentence without a path.
2. API Contract: rendered `docs/agents/frontend.md` — the scope sentence names
   `packages/frontend` under the same rule.
3. API Contract: rendered `docs/agents/skill-dispatch.md` of every profile —
   states the task_03 clause with `mandatory`.
4. API Contract: `roundfix baseline update` on a Standard TypeScript Monorepo
   adopter — retention lists `clause.backend.prohibit-generic-layers` as
   `replaced` by `clause.backend.prohibit-generic-buckets`.

### Surface Transcripts

None. No command, flag, stream or exit code changes; API Contract 4 is
exercised through the classifier the command calls.

## Coverage Map

- Goal 1, Core Feature 2, Success Metric 1 → profile `workspaces[].guide`,
  `workspace.location`, `profileWorkspaceLocations` (task_02).
- Goal 2, Core Feature 1, Success Metric 2 → `backend.json`, Source Baseline
  row, `classifySourceClauseTransition` (task_01).
- Goal 3, Core Feature 3, Success Metric 3 → `core.json`, Source Baseline row,
  this repository's skill guide (task_03).
- Goal 4, Core Feature 4 → ADR-0222 and the existing
  `catalog.profile.skill.dispatch-outside-setup` check (QA).
- Success Metric 4 → every Task's regeneration and refresh.

## Integration Points

None external. The Managed Refresh writes this repository's own guides.

## Testing Approach

1. task_01 creates `internal/baseline/backend_bucket_clause_test.go`:
   `TestTheBackendGuideScopesTheBucketProhibition` renders a Standard
   TypeScript Monorepo plan and checks the new text present and the old absent;
   `TestTheBucketClauseReplacesTheClauseAdoptersHold` runs the Source Baseline
   transition and expects `replaced` with the successor as target, and the
   successor `retained`. Removing `replaces` makes it report `unaccounted`
   (measured).
2. task_02 creates `internal/baseline/workspace_location_test.go`:
   `TestTheBackendAndFrontendGuidesNameTheirDeclaredWorkspace` (plan render),
   `TestAGuideWithoutADeclaredWorkspaceKeepsItsGenericScope` (no binding,
   two paths to one guide, an unbound entry and an unsafe path), and
   `TestEveryBuiltInWorkspaceBindsAGuideItsProfileRenders` (class check over
   every built-in profile).
3. task_03 creates `internal/baseline/person_only_skill_test.go`:
   `TestEveryBuiltInProfileTellsTheAgentToAskForAPersonOnlySkill` (force and
   core selection for every built-in profile) and
   `TestTheSkillGuidesStateThePersonOnlySkillClause` (Rust CLI, Go CLI/TUI and
   Standard TypeScript Monorepo plans, plus the Rust `cut-release` trigger).
4. Every Task's Verification also runs `TestCatalogCompatibility`,
   `TestBaselinePlanCharacterization` and the tests it changes, and checks
   this repository's refresh is current.

## Build Order

1. Scoped bucket clause, its Source Baseline row and its tests (task_01).
2. Declared workspace paths (task_02) (depends on: 1 — both rewrite the
   digest pin, catalog snapshots and plan goldens).
3. Person-only skill clause, its Source Baseline row and its tests (task_03)
   (depends on: 2 — it shares the derived files and, with step 1, the Source
   Baseline files, force record and entry count).
4. QA gate (task_04) (depends on: 1, 2, 3).

## Risks & Considerations

- **The Source Baseline is shared history.** A new row claims that a Source
  Baseline adopter holds the clause, which is true after its next update. The
  old bucket row is never removed.
- **Counts other Specs move.** Each Task raises versions and the entry count
  from its starting main, not from the numbers above.
- **Byte stability for repository-owned profiles.** The templates keep their
  line breaks so an adopter without a bound workspace sees no change; the
  plan goldens prove the generic render path.
- **Hazard: data flow.** Workspace paths come only from embedded built-in
  profile assets and are rendered into setup-owned guides the adopter reads;
  an unsafe path is skipped, never rendered.

## Decisions

- One general person-only clause, no per-trigger marker. See ADR-0222.
- `typescript` keeps the production-code and debugging activations. See
  ADR-0222.
- A reworded clause takes a new identity with `replaces`. See ADR-0222.
- The binding lives on the profile's workspace entry, not on the module,
  because modules are profile-independent. See ADR-0222.
