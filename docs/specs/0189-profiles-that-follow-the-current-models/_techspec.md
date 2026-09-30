---
spec: 0189-profiles-that-follow-the-current-models
prd: _prd.md
created: 2026-09-30
---

# Profiles that follow the current models — Technical Spec

## Executive Summary

Four places name models, and they drifted apart: the Model Catalog, the
recommendation ranking, the built-in profiles with the generated config, and
the Baseline semantic analysis. This design gives them one source. A dated
Recommended Profile per Agent Work Category lives in
`internal/config/recommendations.go`. The built-ins, the generated config and
the legacy Codex default are derived from it, `profiles show` prints it, and
tests tie the catalog and the reference document to it. The catalog, the
picker's efforts and the adapter floors are refreshed to what the installed
adapters advertised on 2026-09-30.

The trade-off this design accepts is the recommendation's evidence. The ranking
carried a DeepSWE result and an average cost per row. The new rows carry a
rationale and a date only, because no current leading model has both figures
published. Figures and sources move to the reference document, where a missing
figure is written as missing. That costs a JSON schema change, from
`roundfix/profiles/v1` to `v2`.

## Project Constraints

- Identifier strategy: not applicable — no new identifier; models keep the
  identifiers the adapters advertise. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — static data, local configuration
  and documentation only. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0180 (this Spec), ADR-0037,
  ADR-0049, ADR-0050, ADR-0069, ADR-0079, ADR-0107, ADR-0140, ADR-0147,
  ADR-0151 and ADR-0153 hold as the PRD states. The gate is bound by ADR-0080,
  ADR-0088, ADR-0091, ADR-0096, ADR-0104, ADR-0117, ADR-0155 and ADR-0156, and
  by ADR-0093 and ADR-0094. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30 for the skill files and `.roundfixrc.yml`, and the standing grant
  of 2026-09-21 for the two governed test files, recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/agents/openai.yaml`,
  `skills/roundfix/agents/openai.yaml`, `internal/cli/cli_test.go`,
  `internal/docscontract/publicdocs_test.go`, `.roundfixrc.yml`. Sanctioned
  regeneration: `make skills-sync`, `make baseline-digests`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## System Architecture

No new package. Each reader is changed where it lives:

| Concern | Owner | File |
| --- | --- | --- |
| Model Catalog | `ModelCatalog` | `internal/agent/catalog.go` |
| Picker efforts | `reasoningEffortChoices` | `internal/cli/cli.go` |
| Adapter floors | `PinnedCodexAdapterVersion`, `PinnedClaudeAdapterVersion` | `internal/agent/acpx_runner.go` |
| Recommended Profile | `RecommendedProfile`, `ModelRecommendations` | `internal/config/recommendations.go` |
| Show | `buildProfilesShowResponse`, `printProfilesShowText` | `internal/cli/profiles.go` |
| Configure preview | `printProfilesConfigureRecommendations` | `internal/cli/profiles_configure.go` |
| Built-ins | `builtinProfiles` | `internal/config/profiles.go` |
| Generated config, legacy default | `defaultConfigYAML`, `Builtin` | `internal/config/config.go` |
| Baseline analysis | `PreferredModel`, `FallbackModel` | `internal/baselineacp/analyzer.go` |

`internal/config` does not import `internal/agent`, and that stays so. The
tests that tie the two together live in `internal/config`'s test package,
which already imports `agent`.

## Reference data

Every Task takes its data from this section and adds none.

### Model Catalog

| Runtime | Value, in order | Description |
| --- | --- | --- |
| codex | `gpt-6.1-sol` | current Codex workhorse; built-in implementation default |
| codex | `gpt-6-astra` | frontier model for the most demanding work; drains quota fastest |
| codex | `gpt-6-sol` | previous workhorse |
| codex | `gpt-6-luna` | fast and affordable |
| codex | `gpt-5.6-sol` | older workhorse |
| codex | `gpt-5.6-terra` | older balanced model |
| codex | `gpt-5.6-luna` | older fast model; built-in review default |
| claude | `opus` | Opus 5.5; design and frontend default |
| claude | `sonnet` | Sonnet 5.5; efficient for routine tasks |
| claude | `claude-fable-5-1` | Fable 5.1; most capable for the hardest work, at the highest latency and quota cost |
| claude | `haiku` | Haiku 4.5; fastest, with no reasoning control |
| claude | `default` | adapter default; currently Opus 5.5 |

`Label` equals `Value`. Picker reasoning efforts: Codex `low`, `medium`,
`high`, `xhigh`, `max`; Claude `default`, `low`, `medium`, `high`, `xhigh`,
`max`.

Two entries stay until task_03, because the shipped ranking and the built-in
profiles name them and both adapters still advertise them:

| Runtime | Value | Position | Description |
| --- | --- | --- | --- |
| codex | `gpt-5.5` | last | leaves Codex on 2026-10-14 |
| claude | `claude-fable-5` | after `claude-fable-5-1` | replaced by `claude-fable-5-1` |

task_01 writes the catalog with these two entries. task_03 removes them, and
the table above is then the whole catalog.

### Recommended Profile, snapshot 2026-09-30

| Category | Preferred Selection | Fallback Chain |
| --- | --- | --- |
| `general`, `backend`, `data`, `infra`, `test`, `qa` | `codex / gpt-6.1-sol / high` | `claude / opus / high` |
| `frontend` | `claude / opus / high` | `codex / gpt-6.1-sol / xhigh` |
| `review` | `codex / gpt-5.6-luna / max` | `codex / gpt-6.1-sol / high` |
| `docs`, `chore` | `codex / gpt-5.6-luna / max` | `claude / sonnet / high` |

Rationales, one per selection:

- Implementation preferred: the Codex workhorse since 2026-09-29, in the quota
  band of the model it replaces; adopted directly by the maintainer.
- Implementation fallback: changes runtime; `opus` resolves to Opus 5.5.
- `frontend` preferred: design judgment; `opus` resolves to Opus 5.5.
- `frontend` fallback: changes runtime; the provider places connected visual
  work at extra-high effort.
- `review`, `docs` and `chore` preferred: bounded work that blocks no other
  Task; about five minutes per session at a fraction of the price.
- `review` fallback: stays on Codex, because every review selection must use
  the pre-PR review provider's runtime.
- `docs` and `chore` fallback: changes runtime; Sonnet 5.5 is the efficient
  Claude model.

### Adapter floors and Baseline analysis

- `PinnedCodexAdapterVersion`: `2.0.1`. `PinnedClaudeAdapterVersion`: `0.84.0`.
- Baseline semantic analysis: `gpt-6.1-sol`, then `gpt-5.6-sol`, both `xhigh`.

### Published evidence for the reference document

Read on 2026-09-30 by the model research; not re-measured here. "Not found"
means no published figure was found, and the document must say so.

| Model | Released | Published evidence | Source |
| --- | --- | --- | --- |
| `gpt-6.1-sol` | 2026-09-29 | Artificial Analysis Index 52; LiveBench agentic 54.5; DeepSWE and Terminal-Bench 4.0 not found | <https://openai.com/index/introducing-gpt-6-1-sol/>, <https://livebench.ai/> |
| `gpt-6-astra` | 2026-09-03 | DeepSWE 74.1% at `xhigh`; Terminal-Bench 4.0 59.6 | <https://deepswe.datacurve.ai/>, <https://artificialanalysis.ai/evaluations/terminalbench-4-0> |
| `gpt-6-sol` | 2026-09-22 | DeepSWE 68.8% at `max`; Terminal-Bench 4.0 43.9 | same |
| `gpt-6-luna` | 2026-09-22 | DeepSWE 66.6% at `max`; Terminal-Bench 4.0 12.6 | same |
| `gpt-5.6-sol` | 2026-07-09 | DeepSWE 72.7% at `max`, 69% at `high`; Terminal-Bench 4.0 39.9 | same |
| `gpt-5.6-luna` | 2026-07-09 | DeepSWE 67.2% at `max` | <https://deepswe.datacurve.ai/> |
| Opus 5.5 (`opus`) | 2026-09-22 | Terminal-Bench 4.0 59.6; LiveBench 71.7; DeepSWE not found | <https://www.anthropic.com/claude-opus-5-5>, <https://livebench.ai/> |
| Sonnet 5.5 (`sonnet`) | 2026-09-28 | LiveBench 39.3; Terminal-Bench 4.0 63.6, unconfirmed on the primary page | <https://www.anthropic.com/claude-sonnet-5-5> |
| Fable 5.1 (`claude-fable-5-1`) | 2026-09-01 | Terminal-Bench 4.0 55.1; LiveBench 66.1 | <https://artificialanalysis.ai/evaluations/terminalbench-4-0> |

Retirements, from <https://learn.chatgpt.com/docs/models>: `gpt-5.5` leaves
Codex with ChatGPT sign-in on 2026-10-14; `gpt-5.4` and `gpt-5.4-mini` left on
2026-08-31. The alias `opus` resolves to Opus 5.5 since Claude Code 2.1.280
(<https://code.claude.com/docs/en/model-config>). Local measurement,
2026-09-15 to 2026-09-30: `gpt-5.6-sol`/`high` averaged 18.8 minutes per
`backend` session and 15.1 per `qa` session; `gpt-5.6-luna`/`max` averaged
about 5 minutes per `docs` session.

## Implementation Design

### Interfaces

```go
// internal/config/recommendations.go
const ModelRecommendationSnapshotVersion = "2026-09-30"

type RecommendationRole string

const (
	RecommendationPreferred RecommendationRole = "preferred"
	RecommendationFallback  RecommendationRole = "fallback"
)

type ModelRecommendation struct {
	Category   WorkCategory       `json:"category"`
	Rank       int                `json:"rank"`
	Role       RecommendationRole `json:"role"`
	Selection  AgentSelection     `json:"selection"`
	SourceAsOf string             `json:"source_as_of"`
	Rationale  string             `json:"rationale"`
}

// RecommendedProfile returns the dated profile Roundfix recommends for
// category. Every Agent Work Category has one.
func RecommendedProfile(category WorkCategory) (AgentSelectionProfile, bool)

// ModelRecommendations lists the Recommended Profile of category in order:
// the Preferred Selection at rank 1, then the Fallback Chain.
func ModelRecommendations(category WorkCategory) ([]ModelRecommendation, bool)
```

`ModelRecommendations` loses its middle result, the source category, because
every category now has its own profile. That is a declared signature change
inside `internal/`.

### Data Models

No persisted schema changes. `roundfix/profiles/v2` is the only changed
output shape, described under API Contract 1.

### The Recommended Profile

`recommendations.go` holds one table keyed by all ten Agent Work Categories.
`RecommendedProfile` returns a copy. `ModelRecommendations` derives its rows
from the same table, so the two cannot disagree. The numeric fields and
`CategorySpecific` leave `ModelRecommendation`.

### Show and configure

`profiles show` keeps the effective profile lines and replaces the ranking
block:

```text
Recommended profile (snapshot 2026-09-30):
  1. preferred codex / gpt-5.6-luna / max
     rationale: <rationale>
  2. fallback claude / sonnet / high
     rationale: <rationale>
```

The `unavailable:` line under a row is kept. The lines `Recommendation source`
and `Recommendations snapshot` are removed. `profilesShowSchema` becomes
`roundfix/profiles/v2`, and `recommendation_source` leaves the profile object.
`printProfilesConfigureRecommendations` prints the same rows in the
interactive flow.

### Built-ins, generated config and the legacy default

- `builtinProfiles` builds each required category from `RecommendedProfile`.
  It defines no optional category, so `ConfiguredWorkCategories` and Profile
  Readiness still cover five categories on a bare configuration.
- `defaultConfigYAML` renders its `profiles:` block from `builtinProfiles`
  instead of a literal, in `requiredWorkCategories` order.
- The legacy Codex runtime default is the Codex selection of the `general`
  Recommended Profile: its model and its effort. The constants
  `defaultCodexModel` and `defaultCodexReasoningEffort` are removed. The
  Claude runtime default does not change.
- `applyLegacyRuntimeProfiles` is unchanged. It already replaces a fallback
  that equals the legacy selection with the built-in Preferred Selection.

### Baseline semantic analysis

`PreferredModel` becomes `gpt-6.1-sol` and `FallbackModel` becomes
`gpt-5.6-sol`. `RequiredReasoningEffort` stays `xhigh`. The Task also edits
ADR-0069's text so its two model names match, and points it at ADR-0180.

### API Contracts

1. API Contract: `roundfix profiles show [--category <category>] [--json]` —
   JSON schema `roundfix/profiles/v2`. Each profile carries `category`,
   `source`, `inherited_from`, `preferred`, `fallbacks` and `recommendations`.
   Each recommendation carries `category`, `rank`, `role`, `selection`,
   `source_as_of`, `rationale` and, when set, `unavailable_reason`. Text
   prints the block above. Exit codes are unchanged.
2. API Contract: Adapter Readiness in `roundfix doctor`, `roundfix setup` and
   Preflight — an official adapter below `codex-acp` 2.0.1 or
   `claude-agent-acp` 0.84.0 is refused with the existing message shape and
   the install action `npm install -g <package>@<floor>`.
3. API Contract: built-in effective profiles — with no configured profile,
   each required category resolves to its Recommended Profile with source
   `built-in`. An optional category resolves to `general` by inheritance.
4. API Contract: generated config — `roundfix setup` and `roundfix init`
   write a `profiles:` block equal to the five built-in profiles.
5. API Contract: Interactive Input — the picker lists the Model Catalog and
   the reasoning efforts under Reference data.

## Coverage Map

- Goal 1 → Built-ins, generated config and the legacy default; Baseline
  semantic analysis; API Contracts 3 and 4.
- Goal 2 → Reference data (Model Catalog); API Contracts 2 and 5.
- Goal 3 → The Recommended Profile; Show and configure; API Contract 1.
- Goal 4 → Testing Approach 6.
- Core Feature 1 → Reference data (Model Catalog); API Contract 5.
- Core Feature 2 → Reference data (Adapter floors); API Contract 2.
- Core Feature 3 → The Recommended Profile; Show and configure; API Contract 1.
- Core Feature 4 → Built-ins, generated config and the legacy default;
  API Contracts 3 and 4.
- Core Feature 5 → Baseline semantic analysis.
- Core Feature 6 → Testing Approach 6.
- Success Metric 1 → Testing Approach 4.
- Success Metric 2 → Testing Approach 1, 3 and 7.
- Success Metric 3 → Testing Approach 3.
- Success Metric 4 → Testing Approach 2.
- Success Metric 5 → Testing Approach 6.

## Integration Points

- **ACP adapters.** No Task calls an adapter. The catalog, the efforts and the
  floors are the advertisement recorded in the PRD. Tests use the existing
  fake adapters.
- **Setup and Doctor.** They read the floors through the two constants, so
  their messages follow without other changes.
- **Pre-PR review.** `validateReviewProfileProvider` requires every selection
  of `profiles.review` to use the provider's runtime. The built-in `review`
  profile satisfies it because the recommended one stays on Codex.

## Testing Approach

1. **Catalog and efforts.** New `internal/agent/catalog_current_test.go`
   requires each catalog to open with the models of Reference data in order,
   and requires the absence of `gpt-5.4`, `gpt-5.4-mini` and
   `gpt-5.3-codex-spark`. New `internal/cli/effort_choices_test.go` requires
   the two effort lists. The existing exact-list test in
   `internal/agent/agent_test.go` is updated, first with the two transitional
   entries and then, in task_03, without them.
2. **Adapter floors.** New `internal/cli/adapter_floor_test.go` drives Doctor
   with the existing fake adapter reporting `codex-acp` 1.1.5 and
   `claude-agent-acp` 0.63.0. Each is refused and the install action names the
   floor constant. The same adapters reporting the floor pass.
3. **Recommended Profile.** New `internal/config/recommended_profile_test.go`
   requires:
   - a profile for each of the ten categories, equal to Reference data;
   - a first fallback on another runtime in every category except `review`,
     and one runtime throughout `review`;
   - every recommended Codex or Claude model present in `agent.ModelCatalog`;
   - `ModelRecommendations` rows equal to the profile in order, with roles.

   New `internal/cli/profiles_show_recommended_test.go` requires the text
   block and the `v2` JSON without the removed fields.
4. **Built-ins.** New `internal/config/built_in_profiles_test.go` requires
   the five built-in profiles to equal `RecommendedProfile`, the generated
   YAML to parse back to the same profiles, the legacy Codex default to be the
   `general` Codex selection, and an optional category to stay undefined. New
   `internal/cli/built_in_review_provider_test.go` requires the built-in
   `review` profile to pass `validateReviewProfileProvider` for `codex`.
5. **Baseline analysis.** New `internal/baselineacp/analyzer_models_test.go`
   requires both analysis models to be in the Codex Model Catalog.
6. **Reference document.** New `internal/docscontract/model_selection_test.go`
   (build tag `docscontract`) requires `docs/references/model-selection.md` to
   state `Updated <snapshot>` and every recommended selection of every
   category in `runtime / model / effort` form.
7. **No replaced model in production code.** New
   `internal/agent/catalog_replaced_test.go` requires the catalog to offer
   neither `gpt-5.5` nor `claude-fable-5`. Task Verification greps non-test Go
   under `internal` and `cmd` for `gpt-5.5`.

Each new gate is proved to fail: the Task that adds it records, in its Result,
the sabotage it applied (for example removing `gpt-6.1-sol` from the catalog)
and the failing test name, then restores the code.

## Build Order

1. Model Catalog, picker efforts and adapter floors, with their guide and
   skill text, task_01 (depends on: none).
2. The Recommended Profile and `profiles show`, with their guide, skill and
   glossary text, task_02 (depends on: 1). The catalog invariant needs the
   refreshed catalog, and both Tasks edit the guides and the skill.
3. Built-ins, generated config, legacy default and Baseline analysis, and the
   removal of the two transitional catalog entries, with their guide and
   skill text, task_03 (depends on: 2).
4. Reference document, its contract test and the `.roundfixrc.yml` comment,
   task_04 (depends on: 2). It shares no file with task_03.
5. Terminal QA, task_05 (depends on: 1, 2, 3, 4).

## Risks & Considerations

- **Test fallout in task_03.** Twenty existing tests pin the old built-ins
  (measured on a scratch copy): seven in `internal/config/config_test.go`,
  seven in `internal/cli/cli_test.go`, two each in
  `internal/cli/doctor_test.go` and
  `internal/cli/doctor_characterization_test.go`, and one each in
  `internal/cli/implement_test.go` and `internal/cli/selection_test.go`. The
  characterization tests change because the break is declared. Where a test
  means "the built-in value", it must read `RecommendedProfile` instead of
  repeating a literal.
- **The order of removal.** The shipped ranking names `gpt-5.5` and
  `claude-fable-5`, and an existing test requires every ranked model to be in
  the catalog. Removing both in task_01 fails that test (measured on a scratch
  copy), so they leave in task_03, after the ranking is gone.
- **A `review` fallback on Claude.** Deriving `review` from `general` would
  fail every pre-PR review. The single-runtime test in Testing Approach 3
  guards it.
- **Raising the floors.** A user on an older adapter is refused by Doctor with
  the install command. That is the intended, declared break; the older
  versions were never shown to advertise `gpt-6.1-sol`.
- **Figures in the reference document.** They were read by the model research
  and not re-measured. The document says so, and the binary carries none.

## Decisions

- One dated Recommended Profile is the source of every built-in selection.
  See ADR-0180.
- Recommendation rows carry a rationale and a date, and `profiles show` moves
  to schema `v2`.
- Optional categories are recommended but not built in.
- `gpt-5.5` leaves the picker before it leaves Codex.
