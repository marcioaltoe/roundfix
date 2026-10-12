---
spec: 0257-the-baseline-respects-the-repository
prd: _prd.md
created: 2026-10-12
---

# The Baseline respects the repository's own choices — Technical Spec

## Executive Summary

Four serial Tasks. task_01 adds the optional `auth.base_path` decision to the
catalog and the two TypeScript profiles, derives the Better Auth route
exception from it during project-decision normalization, and renders the
recorded path in the backend guide. task_02 lets `roundfix baseline update`
record declared decisions before it plans. task_03 replaces the trailing
classification with a Setup Snapshot Difference that Doctor names and the
managed refresh keeps. task_04 rewords two core clauses and one Rust clause to
state the divergence order and the two missing overriding rules, and states
the order in the `roundfix` skill.

The trade-off accepted is ADR-0264's: nothing the binary holds can order two
upstream revisions, so no lock-matching copy is ever classified behind, and a
genuinely stale copy is no longer restored automatically. ADR-0263 accepts a
second trade-off: the provider scope stays stored in `auth.provider` and
`http.contract` for compatibility, and becomes a projection of the base path
only when the base path is recorded.

The source entries are adopted under this Spec:
[the Better Auth base path](references/2026-10-12-the-baseline-hardcodes-the-better-auth-base-path.md),
[the trailing report](references/2026-10-12-a-skill-newer-than-the-setup-snapshot-is-reported-as-trailing.md)
and [the divergence order](references/2026-10-12-the-baseline-does-not-say-how-to-handle-a-diverging-vendored-skill.md).

## Project Constraints

- Identifier strategy: not applicable. No identifier scheme changes; every
  reworded clause keeps its `id`, and the new decision and Doctor code are
  names. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable. Roundfix gains no HTTP route,
  credential or network call; tests use the embedded catalog, temporary
  repositories and fake restore dependencies. The catalog change renders an
  adopter's own recorded Better Auth path and invents no policy:
  `docs/agents/agent-instructions.md`: "do not invent authentication,
  authorization, database, transport, or deployment policy". Source:
  `docs/agents/agent-instructions.md`, `docs/agents/cli.md`.
- Active ADR obligations: applicable. ADR-0263 and ADR-0264 (this Spec)
  govern the design. ADR-0076: "Deterministic conflict detection is required
  so the two values can never render contradictory backend guidance"; a
  recorded base path moves the one provider exception and keeps the conflict
  check. ADR-0144: "ask only when a decision identifier is missing or
  incompatible". ADR-0221: "A skill that differs is reported as trailing the"
  snapshot, replaced by ADR-0264. ADR-0224: "here a trailing skill fails"
  `make verify`, kept for this repository. ADR-0186: "Every clause keeps its
  enforcement level when its wording changes". ADR-0257: "A clause whose
  wording changes keeps its identity". ADR-0250: "A Baseline module's version
  names one content, and the record step chooses it". ADR-0256: "Spec
  authoring requires a skills Task for a change a skill describes". ADR-0244:
  "A Spec declares the domain terms it introduces". Source:
  `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable. Express maintainer authorization: "considere
  autorizado a ajustar todas as skills se necessário", "Autorizar os dois" and
  the cycle AE "Concedo" of 2026-10-11 for the Governed Paths this Spec
  declares. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`. Spec-contained authorization record:
  `docs/specs/0257-the-baseline-respects-the-repository/_authorization.md`.
  Bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `.agents/skills/roundfix/references/setup.md`,
  `.agents/skills/setup-context-driven/SKILL.md`,
  `docs/agents/setup-context.json`, `docs/agents/skill-dispatch.md`,
  `docs/references/coverage-record.json`,
  `internal/baseline/assets/decisions.json`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`,
  `internal/baseline/assets/modules/core.json`,
  `internal/baseline/assets/modules/rust.json`,
  `internal/baseline/assets/profiles/go-cli-typescript-monorepo.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/baseline/assets/templates/guides/backend.md`,
  `internal/baseline/assets/templates/index.json`,
  `skills/roundfix/SKILL.md`, `skills/setup-context-driven/SKILL.md`.

## Current behavior

Characterized on `7a8b7f70`:

- `auth.provider` (`internal/baseline/assets/decisions.json`) suggests
  `routeException.scope` `/api/auth/*`. `normalizeProjectDecisions` keys an
  `http.contract` exception by owner and scope, so a provider value with scope
  `/auth/*` over a manifest whose `http.contract` holds `/api/auth/*` appends
  a second Better Auth exception and keeps the stale one.
- `baseline update` reads decisions only from the Setup Manifest; it has no
  decision input besides `--adopt-suggested`. `baseline plan` on an adopted
  repository asks for every decision again (measured: exit 0 with "required
  Baseline decisions are missing" for nine decisions).
- `TrailingSetupSkills` (`internal/baseline/skills_trailing.go`) returns every
  lock-matching required external skill whose portable digest differs from
  the snapshot `treeDigest`. Doctor appends `DR-SKILL-TRAILS-SNAPSHOT: trails
  the Setup Snapshot: <names>`, raises the line to `warn` and adds
  `roundfix baseline update`. The update preview lists each under
  `skills.drifted` and reports `plan_ready`; the skills stage restores each.
- `clause.core.prohibit-editing-vendored-skills` says only not to edit a
  vendored skill and to state a correction as a Repository-Specific Normative
  Rule.

## System Architecture

```text
decisions.json + profiles + templates/guides/backend.md ─┐
project_decisions.go (normalize) + project_decision_render.go ─┴─> BuildPlan -> backend guide, Setup Manifest
cli baseline update --decision ─> baseline.DeclareManifestDecisions ─> BuildPlan (Managed Refresh)
skills_trailing.go SetupSnapshotDifferences ─┬─> doctor skills line
                                             └─> baseline update preview and skills stage (skills.differing)
modules/core.json + modules/rust.json ─> record step -> make baseline-digests -> Managed Refresh
```

No new package. Every change extends an existing file except the new tests.

## Implementation Design

### Interfaces

```go
// internal/baseline
const authBasePathDecisionID = "auth.base_path"

func normalizeAuthBasePath(value any) (string, error)

// DeclaredDecision reports one decision recorded or changed by an update.
type DeclaredDecision struct {
	ID       string `json:"id"`
	Value    any    `json:"value"`
	Recorded any    `json:"recorded"` // previous recorded value, or null
}

func DeclareManifestDecisions(input ManifestInput, declared []DecisionValue, catalog *Catalog) (ManifestInput, []DeclaredDecision, error)

// SetupSnapshotDifferences replaces TrailingSetupSkills with the same inputs.
func SetupSnapshotDifferences(repoRoot, profileID string, required []string) ([]string, error)
```

### Data Models

New catalog entry, inserted after `auth.provider` in `decisions.json`:

```json
{
  "id": "auth.base_path",
  "version": 1,
  "optional": true,
  "type": "string",
  "default": "/api/auth",
  "requiresCapabilities": ["capability.stack.better-auth"],
  "summary": "The path where the repository mounts Better Auth's routes, its basePath option; the provider-owned route exception covers that path.",
  "effects": [{"when": {"present": true}, "renderBindings": [
    {"artifact": "guide.backend", "template": "template.guide.backend", "token": "auth.base_path"}]}]
}
```

`auth.base_path` follows `auth.provider` in the `entryDecisions` of
`standard-typescript-monorepo.json` and `go-cli-typescript-monorepo.json`.
`template.guide.backend` rises from version 3 to 4 and declares the token; its
last line becomes `{{auth.provider}}{{auth.base_path}}`.

Valid base path: a string equal to its trimmed form, starting with `/`, at
least two characters, not ending with `/`, with no `//`, `*`, `?`, `#`,
whitespace, backtick or control character, and no marker-shaped text.

`baseline-update-result/v1` gains two always-present arrays:
`declaredDecisions` (`DeclaredDecision` items) and `skills.differing`
(`{"skill","reason"}` items, the shape of `skills.drifted`).

### API Contracts

1. API Contract: base-path normalization. When `auth.base_path` is recorded
   and valid, normalization sets the provider's `routeException.scope` to
   `<base path>/*`, removes from `http.contract` the exception equal to the
   provider exception as recorded, and then derives `http.contract` as today,
   so exactly one Better Auth exception covers the new scope. An invalid value
   is refused with "auth.base_path must be an absolute path with at least one
   segment and no trailing slash, wildcard, query, fragment or whitespace";
   `auth.base_path` without `auth.provider` is refused with "auth.base_path
   requires auth.provider". When `auth.base_path` is absent, normalization is
   byte-for-byte today's.
2. API Contract: base-path rendering. While absent, the token renders the
   empty string. When recorded, it renders a blank line followed by
   ``The recorded Better Auth base path is `<path>`, the value of Better
   Auth's `basePath` option.`` with the path through `markdownCode`.
3. API Contract: declared decisions. `roundfix baseline update` accepts
   `--decision <id=value>`, repeatable, parsed as `baseline plan --decision`
   parses it. After the manifest resolves and before planning, each value
   replaces or adds the recorded decision and the set is normalized as a plan
   normalizes it. A malformed pair, an unknown decision, a decision the
   profile does not select, a profile-fixed decision or `preservation.mode`
   exits `2` with category `invalid` and writes nothing. A declared decision
   removes itself from `newDecisions`. The result lists each declared
   decision; text prints `Declared decisions: <n>` and one
   `- <id>=<compact JSON value> (recorded: <compact JSON value or none>)` line
   each. Without confirmation nothing is written; the Plan Digest covers the
   changed manifest.
4. API Contract: help. The synopsis lines of `roundfix --help`,
   `roundfix baseline --help` and `roundfix baseline update --help` add
   `[--decision <id=value> ...]` after `[--adopt-suggested]`, and the update
   options add `--decision` with "Record or change a decision as id=value;
   repeat for several decisions".
5. API Contract: Setup Snapshot Difference in Doctor. A lock-matching
   required external skill that differs from its snapshot appends
   `; DR-SKILL-DIFFERS-FROM-SNAPSHOT: differs from the Setup Snapshot, newer or older unknown: <sorted names>`
   to the `skills` detail. The status and next action stay what readiness
   set. Doctor's exit code is unchanged.
6. API Contract: Setup Snapshot Difference in the update. The preview and the
   skills stage list each such skill under `skills.differing` with reason
   `differs from its Setup Snapshot, newer or older unknown; kept; restore deliberately with roundfix baseline skills restore --profile <profile> --skill <skill>`,
   and text prints `Skills differing: <n>` and `- differs <skill>: <reason>`
   after the drifted list. A difference alone never makes the preview
   `plan_ready`, and the skills stage never restores it. Missing and
   lock-mismatched skills keep their `skills.drifted` restore.
   `DR-SKILL-TRAILS-SNAPSHOT`, "trails its Setup Snapshot" and "trails the
   Setup Snapshot" leave every live source, guide and skill.
7. API Contract: the Baseline clauses. The catalog carries Exact clause texts
   under their current identities and enforcements; this repository's
   rendered `docs/agents/skill-dispatch.md` carries both core texts, and a
   Source Baseline adopter's refresh records both core clauses retained.

### Surface Transcripts

None. No command gains a new complete transcript: API Contracts 3 to 6 fix
the added lines, which named tests assert inside existing outputs, and the
Behavior Surface Record carries the changed help.

### Exact clause texts

`clause.core.activate-matching-skills` (mandatory, `rule.core.skill-dispatch`):

```text
Activate every matching required skill before governed work. When one skill has distinct active-module triggers, retain and follow each trigger. When a skill's default conflicts with a Baseline rule or a Repository-Specific Normative Rule, follow the rule. Read a skill's bundled files from the directory where the skill is installed, never from a fixed absolute path the skill names.
```

`clause.core.prohibit-editing-vendored-skills` (prohibited, same rule):

```text
Do not edit a skill the repository installs from an upstream source; its upstream owns its text, and the next sync replaces the folder. When such a skill diverges from what the repository needs, choose in this order. First, state the correction as a rule that governs over the skill's default, in the Baseline or as a Repository-Specific Normative Rule; this is the answer to a convention difference such as the package manager, the test runner, or a library choice. Second, add a repository-authored skill that supplies or overrides only what is missing and leaves the upstream skill installed. Third, replace the upstream skill with a repository-authored one only when it is useless or almost useless to the repository. Send an objective upstream defect to the upstream as a pull request; a defect never justifies the third choice.
```

`clause.rust.read-current-docs` (mandatory, `rule.rust.current-docs`):

```text
Use current authoritative crate and toolchain documentation before changing Rust APIs, configuration, or dependencies. Follow the edition the crate's `Cargo.toml` declares, not the one a skill example assumes: under edition 2024, `std::env::set_var` and `std::env::remove_var` are `unsafe`, and each call needs an `unsafe` block whose `// SAFETY:` comment states why no other thread reads or writes the environment at that point.
```

`rule.core.skill-dispatch` rises from 6 to 7 and `rule.rust.current-docs`
from 3 to 4. The module versions are the record step's.

### Divergence audit

The skills maintainer's report of 2026-09-30, measured at the snapshot commit
`b3c45a4` and in `~/dev/skills` at `8660c80`, against the catalog:

| skill | divergence | overriding clause |
| --- | --- | --- |
| `bun` | bare `bun test` | `clause.bun.use-bun-owned-commands` |
| `hono` | `npx hono request` | `clause.bun.prohibit-other-package-managers` |
| `drizzle-orm` | `npm i` | `clause.bun.prohibit-other-package-managers` |
| `golang-cli` | Cobra and Viper | `clause.go.record-the-reason-for-a-module` |
| `golang-testing` | testify | `clause.go.test-through-go-test` |
| `golang-lint` | golangci-lint configuration | `clause.core.require-tooling-authorization` |
| `rust-expert` | `anyhow` as the entry-point default; `async-std` | `clause.rust.keep-type-erased-errors-at-the-entry-point`; `clause.core.review-new-dependencies` |
| `rust-cli` | `env::set_var` outside `unsafe`, `edition = "2021"` template | `clause.rust.read-current-docs` (task_04) |
| `clap-rust` | fixed `/mnt/skills/user/` path | `clause.core.activate-matching-skills` (task_04) |

The `oss-publish` reference in `rust-cli` names a skill that is not vendored;
no Baseline rule can supply it, and it stays out of scope.

### Skill and document texts

task_01, `.agents/skills/setup-context-driven/SKILL.md`, appended to the
paragraph that ends "The human must keep or change that complete typed
proposal.":

```text
`/api/auth` is Better Auth's default base path, not a Baseline rule: a repository that mounts Better Auth elsewhere records the optional `auth.base_path` decision, and the provider and HTTP exceptions then cover that path followed by `/*`.
```

task_01, `docs/user-guide/context-driven-development.md`, a paragraph after
the one that ends "refuses conflicting provider and HTTP values.":

```text
The Better Auth base path is the optional `auth.base_path` decision. While none is recorded, the recorded provider exception stands, and `/api/auth` is Better Auth's own default. A repository that mounts Better Auth elsewhere records its path; the provider exception and the matching `http.contract` exception then cover that path followed by `/*`, and the backend guide names the recorded path.
```

task_02, `.agents/skills/roundfix/references/baseline.md` and
`docs/user-guide/commands/baseline.md`, after the `--adopt-suggested`
sentence:

```text
`--decision <id=value>` records or changes a decision the repository's profile selects before the plan is built; repeat it for several decisions. The value is parsed as `baseline plan --decision` parses it, the plan and its digest cover the changed Setup Manifest, and the result lists each one under `declaredDecisions` with its previously recorded value. A decision the profile does not select or fixes, `preservation.mode`, or a malformed pair exits `2` and writes nothing. A declared decision also answers a new decision without `--adopt-suggested`. A repository that mounts Better Auth at `/auth` runs `roundfix baseline update --decision auth.base_path=/auth`.
```

task_03, replacing the trailing paragraph of
`.agents/skills/roundfix/references/baseline.md` and of
`docs/user-guide/commands/baseline.md`:

```text
Doctor reports `DR-SKILL-DIFFERS-FROM-SNAPSHOT` when a required upstream skill is present and matches its lock but differs from the embedded Setup Snapshot. Newer or older is unknown, so the `skills` line keeps its readiness status and gains no next action. The managed refresh lists such a skill under `skills.differing`, keeps it, and never restores it unasked; restore one deliberately with `roundfix baseline skills restore --profile <id> --skill <name>`. A missing skill, or one that no longer matches its lock, is still restored under `skills.drifted`.
```

task_03, `docs/user-guide/commands/doctor.md`, replacing the
`DR-SKILL-TRAILS-SNAPSHOT` paragraphs and their example with the same facts
and the API Contract 5 line; and `.agents/skills/roundfix/references/setup.md`,
after the Doctor example:

```text
When a required external skill matches its lock but differs from its Setup Snapshot, the `skills` line keeps its status and appends `DR-SKILL-DIFFERS-FROM-SNAPSHOT: differs from the Setup Snapshot, newer or older unknown: <names>`; nothing is restored.
```

task_03, `docs/adr/0221-a-required-external-skill-is-held-to-its-setup-snapshot-not-only-its-lock.md`, a dated addendum and `updated_at`:

```text
## Addendum — 2026-10-12

ADR-0264 replaces the trailing classification and the automatic restore above. A lock-matching skill that differs from its Setup Snapshot is a Setup Snapshot Difference, reported with `DR-SKILL-DIFFERS-FROM-SNAPSHOT` and never restored unasked. The lock rule and the comparison stay.
```

task_04, `.agents/skills/roundfix/references/baseline.md`, a new section
under its own heading:

```text
## Vendored skills

A vendored skill is installed from an upstream source and never edited in the repository, because the next sync replaces its folder; the restore and the managed refresh only install, compare and restore it. When one diverges from what the repository needs, correct it in this order:

1. A rule that governs over the skill's default: a Baseline clause, or a Repository-Specific Normative Rule. This answers a convention difference such as the package manager, the test runner or a library choice.
2. A repository-authored skill that supplies or overrides only what is missing, installed beside the vendored one.
3. A repository-authored skill that replaces the vendored one, only when the vendored skill is useless or almost useless to the repository.

Send an objective upstream defect to the upstream as a pull request; it never justifies the third choice.
```

`CONTEXT.md` entries, written through `domain-modeling` beside related terms:

```text
**Better Auth Base Path**:
The path under which a repository mounts Better Auth's routes, the value of Better Auth's `basePath` option, recorded as the optional `auth.base_path` decision. When it is recorded, the provider's route exception and the matching HTTP Contract Decision exception cover that path followed by `/*`; when it is absent, the recorded exception stands, and Better Auth's own default is `/api/auth` (ADR-0263).
_Avoid_: Auth prefix, API prefix, auth scope

**Setup Snapshot Difference**:
A required external skill that is installed and matches its `skills-lock.json` entry but whose tree differs from the `treeDigest` its Setup Snapshot pins. Nothing offline tells whether the copy is newer or older, so Doctor names it without a warning and the managed refresh keeps it; only a deliberate per-skill restore replaces it (ADR-0264).
_Avoid_: Trailing skill, outdated skill, stale skill

**Vendored Skill**:
A skill a repository installs from an upstream source and never edits locally, because the next sync replaces its folder. A divergence is corrected in order: a rule that governs over the skill's default, then a repository-authored skill that supplies only what is missing, then a repository-authored replacement only when the vendored skill is useless or almost useless; an objective defect goes upstream as a pull request (ADR-0264).
_Avoid_: Forked skill, local skill copy, edited upstream skill
```

Each ADR reference in `CONTEXT.md` is a Markdown link to its file, as the
existing entries write it.

### Derived files

The authoring prototype, a scratch clone of `7a8b7f70`, applied task_01 and
task_04 and ran the sanctioned commands after each:

- task_01: `internal/baseline/testdata/catalog.digest`,
  `catalog.normalized.json`, the four
  `internal/baseline/testdata/plan-characterization/*.golden.json`, and
  `docs/agents/setup-context.json`. The Standard TypeScript golden
  `backend.md` and the profile digest pins did not change.
- task_04: `internal/baseline/module-versions.json`, the Standard TypeScript
  golden `docs/agents/skill-dispatch.md`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/baseline/testdata/catalog.diagnostics.golden.json`,
  `catalog.digest`, `catalog.normalized.json`, the four plan goldens,
  `docs/agents/setup-context.json` and `docs/agents/skill-dispatch.md`.
- An owned-skill edit changed only the canonical file, its mirror, the
  `SKILL.md` pair and `skills/testdata/owned-skill-versions.json`;
  `make baseline-digests` reported no change.
- A change to a user guide page or to command help changes
  `docs/references/behavior-surfaces.json`, re-recorded only with
  `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`.

### Invariants

```text
1. A Setup Manifest without auth.base_path normalizes, renders and records exactly as before; the Standard TypeScript golden backend.md is byte-identical.
2. With auth.base_path recorded, http.contract holds exactly one Better Auth exception, and its scope is the base path followed by /*.
3. baseline update never changes the profile or the preservation mode, and writes nothing without confirmation.
4. No lock-matching required external skill is restored by the managed refresh unless a person runs the per-skill restore.
5. This repository's required skills are still held to the snapshot by TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot.
6. No clause identity is added or removed; each reworded clause keeps its enforcement, and no two clauses share text.
7. After each Task, roundfix baseline update --repo . --no-skills --no-history reports current.
8. Every module and owned-skill version is the one its record command chose; no hand-edited pin, golden, guide or record.
9. No Makefile, CI workflow, lint, formatter or test-runner configuration, go.mod, Source Baseline asset or retention transition changes.
10. A renamed test is re-recorded in the coverage record by its record command, in the Task that renames it; no other test function is deleted or renamed.
```

## Coverage Map

- Goal "base path recorded, no duplicate" → API Contracts 1 and 2, Invariants 1 and 2.
- Goal "one confirmed update" → API Contracts 3 and 4, Invariant 3.
- Goal "never trailing, never restored unasked" → API Contracts 5 and 6, Invariant 4.
- Goal "divergence order and overriding rules" → API Contract 7, Divergence audit.
- User Story 1 → API Contracts 1 to 3.
- User Story 2 → API Contract 1, Invariant 1.
- User Story 3 → API Contracts 3 and 4.
- User Story 4 → API Contracts 5 and 6.
- User Story 5 → API Contract 7.
- Core Feature 1 → API Contracts 1 and 2.
- Core Feature 2 → API Contracts 3 and 4.
- Core Feature 3 → API Contracts 5 and 6.
- Core Feature 4 → API Contract 7, Divergence audit.
- Core Feature 5 → Skill and document texts, Build Order 1 to 4.
- Success Metric 1 → API Contracts 1 to 3, Invariant 2.
- Success Metric 2 → API Contract 1, Invariant 1.
- Success Metric 3 → API Contracts 5 and 6, Invariant 4.
- Success Metric 4 → API Contract 7, Divergence audit.

## Integration Points

None external. The restore path, the lock and the embedded snapshots are read
as today; no network call is added, and Doctor stays offline.

## Testing Approach

- `internal/baseline/better_auth_base_path_test.go` (task_01, new), through
  `BuildPlan` with the Standard TypeScript profile as `plan_test.go` does:
  `TestTheBetterAuthBasePathIsAnOptionalDecision` (declaration, both
  profiles), `TestARecordedBetterAuthBasePathMovesTheProviderException`
  (provider scope, one `http.contract` exception, both guide lines and the
  API Contract 2 sentence, manifest value),
  `TestAnUnrecordedBetterAuthBasePathKeepsTheRecordedScope` (no sentence, no
  manifest entry, `/api/auth/*` lines, and a recorded `/v1/auth/*` provider
  scope kept), `TestAnInvalidBetterAuthBasePathIsRefused` (`auth`, `/`,
  `/auth/`, `/auth/*`, `/a b`, `/a//b`, `/a?x`, and a base path without the
  provider). Declared break: `profile_alignment_test.go`
  (`TestDivergenceRendersProbe` now expects `RemovedDecisions`
  `[auth.provider auth.base_path]`).
- `internal/cli/baseline_update_decision_test.go` (task_02, new):
  `TestBaselineUpdateRecordsADeclaredDecision`,
  `TestBaselineUpdateDeclaresTheBetterAuthBasePath` (a temporary Standard
  TypeScript adopter; Success Metric 1 end to end, second update `current`),
  `TestBaselineUpdateRefusesAnUndeclarableDecision`. Declared possible break:
  `internal/cli/baseline_update_test.go`, where a test pins the JSON field
  set.
- task_03 renames the tests whose meaning inverts and changes them to the
  new contract, then re-records `docs/references/coverage-record.json` with
  `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`:
  in `internal/cli/doctor_trailing_skills_test.go`,
  `TestDoctorWarnsOnASkillThatTrailsItsSnapshot` becomes
  `TestDoctorNamesASetupSnapshotDifferenceWithoutWarning`,
  `TestBaselineUpdatePreviewListsATrailingSkill` becomes
  `TestBaselineUpdatePreviewListsASetupSnapshotDifference` (preview
  `current`, `skills.differing`), and
  `TestBaselineUpdateRestoresTrailingSkillsThroughExistingRequest` becomes
  `TestBaselineUpdateKeepsASetupSnapshotDifference` (a fake restore that
  fails the test if called for the differing skill, and a missing skill still
  restored); in `internal/cli/baseline_update_repository_profile_test.go`,
  `TestBaselineUpdatePreviewListsATrailingSkillWithARepositoryProfile`
  becomes `TestBaselineUpdatePreviewListsADifferenceWithARepositoryProfile`;
  in `internal/baseline/skills_trailing_test.go` and
  `skills_snapshot_profile_test.go`, the two `TestTrailingSetupSkills…` tests
  become `TestSetupSnapshotDifferencesCompareTheInstalledTreeWithTheSnapshot`
  and `TestSetupSnapshotDifferencesCompareARepositoryProfile`.
  `TestDoctorComparesARepositoryProfileWithItsSnapshot` and
  `TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot` keep their
  names and change only the function they call and the strings they expect.
- `internal/baseline/vendored_skill_clauses_test.go` (task_04, new):
  `TestTheVendoredSkillClausesCarryTheDivergenceOrder` (both core texts and
  the Rust text in the catalog, with enforcement and no `replaces`),
  `TestTheVendoredSkillClausesRenderInTheGuides` (each core text once as a
  forced bullet in the Standard TypeScript golden `skill-dispatch.md`),
  `TestAnAdopterRetainsTheVendoredSkillClauses` (Source Baseline adopter
  through `newClauseReplacementAdopter`: both core clauses `retained`, none
  `unaccounted`). Declared break: `promoted_core_clauses_test.go` (the
  vendored-skill sentence it pins).
- The prototype ran `internal/baseline`, `internal/cli` (baseline, doctor and
  skill tests) and `skills` over task_01 and task_04; the two declared breaks
  above were the only failures.

## Build Order

1. Better Auth Base Path: catalog, profiles, template, normalization,
   rendering, setup skill, user guide and **Better Auth Base Path**.
2. Declared decisions at update (depends on: 1, because its end-to-end test
   declares `auth.base_path` and both re-record the Behavior Surface Record).
3. Setup Snapshot Difference (depends on: 2, because both change
   `internal/cli/baseline_update.go`, the update guide, the `roundfix` skill
   and its version record).
4. Divergence order (depends on: 3, because both change the `roundfix` skill
   and its record; and on 1, because both change the catalog's derived files
   and `docs/agents/setup-context.json`).
5. QA gate (depends on: 1, 2, 3, 4).

## Risks & Considerations

- Specs 0255 and 0259 change `baseline update` and Doctor beside this one.
  The later one rebases and re-runs its record commands.
- An adopter whose copies are genuinely stale is no longer prompted to restore
  them; the per-skill restore and its own skills tool remain (ADR-0264).
- A repository-owned profile copied from a TypeScript profile before this
  Spec lacks `auth.base_path`; it keeps working, and only cannot declare the
  base path until the profile lists it.
- `--decision` widens what an update writes. It is held to the plan's
  normalization and the confirmation contract, and refuses profile and
  preservation changes.

## Glossary

- adds: **Better Auth Base Path**
- adds: **Setup Snapshot Difference**
- adds: **Vendored Skill**

## Decisions

- `auth.base_path` is optional with default `/api/auth`, and the provider
  exception is derived from it only when it is recorded; see ADR-0263.
- `baseline update --decision` records declared decisions under the plan's
  normalization; see ADR-0263.
- A lock-matching difference is kept and never classified behind; the
  trailing code is retired; see ADR-0264.
- The two missing overriding rules extend existing clauses instead of adding
  identities, so no Source Baseline row or force record changes; see
  ADR-0264 and ADR-0257.
- Tests whose meaning inverts are renamed, and the coverage record is
  re-recorded by its command in the same Task.
