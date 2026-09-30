---
spec: 0199-stack-rules-that-say-what-they-mean
prd: _prd.md
created: 2026-09-30
---

# Stack rules that say what they mean — Technical Spec

## Executive Summary

Six clauses in three Baseline modules change their guidance, one Skill
Activation changes its trigger, four guide templates gain one scope sentence
each, and three stale statements of the suggested HTTP mode are aligned with
the decision catalog. No production Go code changes: the modules, templates
and activation file are the only sources, and the goldens, digest pins, catalog
snapshots and this repository's guides are regenerated from them by the
sanctioned commands. Four new test files hold the checks, one or two per Task.
The trade-off
accepted is that scope is stated in words. Backend and frontend clauses are
retained by the Source Baseline by their bytes, so they are scoped by a guide
sentence and not reworded, and no guide renders a workspace path. That costs
precision a directory-bound rule would have and buys a change that every
adopter can take on its next update with no new obligation. Every text below
was applied in a disposable clone on 2026-09-30 and passed the repository
Verification there.

## Project Constraints

- Identifier strategy: not applicable — no new identifier; clause, rule, guide,
  module, template and trigger identifiers are kept. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — Baseline assets and local files
  only; no credential and no network call. The Baseline's stated suggestion
  for the HTTP Contract Decision changes; a recorded decision does not.
  Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0190 governs the scope wording and
  the rule-over-skill-default sentence, and ADR-0186 governs adopter-neutral
  wording. ADR-0058 and ADR-0060 keep every Source Baseline clause accounted
  for. ADR-0059 keeps generated guides formatter-stable. ADR-0061, ADR-0063 and
  ADR-0067 are the profile, HTTP and composition decisions the corrected
  statements restate, and none changes. ADR-0073, ADR-0081, ADR-0103 and
  ADR-0149 govern the regeneration and this repository's Managed Refresh.
  ADR-0130 has the audit judge only Governed Paths. ADR-0080, ADR-0088,
  ADR-0091, ADR-0093, ADR-0094, ADR-0096, ADR-0104, ADR-0117, ADR-0155,
  ADR-0156, ADR-0166 and ADR-0167 bind the gate and the checker. ADR-0099 makes
  retention accounting a comparison of clause identities, which the kept
  identifiers satisfy. ADR-0187 and ADR-0189 are active in this tree and do not
  apply, because this Spec edits no skill, no skill version and no command
  reference. ADR-0074, ADR-0097, ADR-0168, ADR-0173, ADR-0176 and ADR-0179 cite
  a listed ADR and do not apply; the PRD's row gives the reason for each.
  ADR-0177 is reached only through ADR-0173 and is equally untouched. All hold.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30, recorded in [_authorization.md](_authorization.md); bounded
  files: `internal/baseline/assets/modules/bun.json`,
  `internal/baseline/assets/modules/typescript.json`,
  `internal/baseline/assets/modules/core.json`,
  `internal/baseline/assets/modules/backend.json`,
  `internal/baseline/assets/modules/frontend.json`,
  `internal/baseline/assets/skill-activations.json`,
  `internal/baseline/assets/contract-v1.json`,
  `internal/baseline/assets/templates/index.json`,
  `internal/baseline/assets/templates/guides/typescript-bun.md`,
  `internal/baseline/assets/templates/guides/bun.md`,
  `internal/baseline/assets/templates/guides/backend.md`,
  `internal/baseline/assets/templates/guides/frontend.md`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/typescript-bun.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/backend.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/frontend.md`,
  `docs/agents/agent-instructions.md`, `docs/agents/skill-dispatch.md`,
  `docs/agents/setup-context.json`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## System Architecture

| Concern | Source of truth | Derived by |
| --- | --- | --- |
| Clause text and versions | `internal/baseline/assets/modules/*.json` | hand edit under the grant |
| Testing trigger | `internal/baseline/assets/skill-activations.json` | hand edit under the grant |
| Guide scope sentences and template versions | `internal/baseline/assets/templates/guides/*.md`, `internal/baseline/assets/templates/index.json` | hand edit under the grant |
| Stated HTTP suggestion | the profile's `httpContract.default`, `internal/baseline/assets/contract-v1.json`, `docs/user-guide/context-driven-development.md` | hand edit; the first two under the grant |
| Suggested HTTP mode in effect | `internal/baseline/assets/decisions.json` (`http.contract`, default `REST`) | unchanged |
| Formatter goldens, the profile's digest pin, catalog snapshots, plan characterization goldens | the sources above | `make baseline-digests` |
| This repository's guides and Setup Manifest | the sources above | `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text` |

Every Task edits its sources and then runs the two regeneration commands, in
that order. A second Managed Refresh must report `File changes: 0`. No pin,
golden or generated guide is edited by hand.

`make baseline-digests` compiles the Baseline test package. A new test file
that does not compile stops the regeneration, so a Task adds its tests in a
compiling state before it regenerates.

The outputs were measured on 2026-09-30 in a fresh clone of `9e439dbb`, one
commit per Task:

- every Task rewrites
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json` (the
  digest pin), `internal/baseline/testdata/catalog.diagnostics.golden.json`,
  `internal/baseline/testdata/catalog.digest`,
  `internal/baseline/testdata/catalog.normalized.json`, the four plan goldens
  under `internal/baseline/testdata/plan-characterization/`
  (`advisory-only-divergences`, `clean-adoption`,
  `idempotent-replan-after-verified-apply`,
  `same-baseline-changed-profile-and-catalog-digests`) and
  `docs/agents/setup-context.json`;
- task_01 also rewrites the formatter golden `typescript-bun.md`. This
  repository has no TypeScript guide, so no repository guide changes;
- task_02 also rewrites the formatter goldens `agent-instructions.md` and
  `skill-dispatch.md` and this repository's guides of the same names;
- task_03 also rewrites the formatter goldens `backend.md` and `frontend.md`.
  This repository has neither guide.

The Source Baseline corpus, the parity corpus, the setups and the decision
declarations stay byte-identical.

A module edit must preserve the file's existing formatting: replace the
guidance string and the version numbers in place. `core.json` and
`backend.json` are hand-formatted, and a rewrite through a JSON encoder
changes lines the catalog diagnostics test anchors on.

## Implementation Design

### Interfaces

No exported Go signature changes and no production file changes. The checks
are test-only functions, each a pure function over rendered text or decoded
catalog documents so a negative test can feed it the replaced wording.

```go
// internal/baseline/stack_rule_wording_test.go (task_01)
type stackWordingRule struct {
    guide   string   // rendered guide path
    must    []string // sentences the guide must carry
    mustNot []string // replaced sentences it must not carry
}
func stackWordingFindings(rendered map[string]string, rules []stackWordingRule) []string
func clauseForce(module any, clauseID string) (enforcement string, found bool)

// internal/baseline/stack_core_wording_test.go (task_02)
func languageTriggerFindings(activations document, words map[string]string) []string

// internal/baseline/stack_scope_and_http_default_test.go (task_03)
func profileHTTPDefaultFindings(profile, decision document) []string

// internal/cli/baseline_http_default_statement_test.go (task_03)
type httpDefaultStatement struct{ source, text, format string }
func httpDefaultStatementFindings(mode string, modes []string, statements []httpDefaultStatement) []string
```

Each Task creates its own test file, so no test file is edited by two Tasks.
The later files call `stackWordingFindings` from the same package.

`stackWordingFindings` collapses whitespace before it compares, because a
template sentence wraps across lines. It reports a guide that is not rendered,
a missing sentence and a replaced sentence still present.

### Data Models

None. No schema, record or payload changes.

### Version changes

Raise each version by one from its value on the starting main. Another Spec
may already have raised one of them.

- task_01: the `bun` module, `rule.bun.package-manager`,
  `rule.bun.warning-free-verification` and `guide.bun`; the `typescript`
  module, `rule.typescript.current-docs` and `guide.typescript-bun`;
  `template.guide.typescript-bun` and `template.guide.bun`.
- task_02: the `core` module, `rule.core.skill-dispatch`,
  `rule.core.dependency-discipline`, `guide.agent-instructions` and
  `guide.skill-dispatch`; the `version` of the activation file.
- task_03: the `backend` module and `guide.backend`; the `frontend` module and
  `guide.frontend`; `template.guide.backend` and `template.guide.frontend`.

### Exact texts

Each replacement is exact. Text outside the quoted sentences stays
byte-identical, and every clause keeps its `id` and its `enforcement`.

**task_01, `bun` module.**

- `clause.bun.use-bun-owned-commands` becomes: "Inside the Bun workspace, use
  Bun-owned commands for dependency installation, scripts, and lockfile
  updates, and run tests through the package's `test` script (`bun run test`),
  never through a bare runner such as `bun test`."
- `clause.bun.prohibit-other-package-managers` becomes: "Inside the Bun
  workspace, do not substitute another JavaScript package manager or runner
  (npm, pnpm, yarn, npx) or hand-edit the lockfile. A toolchain for another
  language in the same repository keeps its own package manager."
- `clause.bun.block-warnings-when-profile-treats-them-as-errors` becomes: "When
  the repository's Verification treats warnings as errors, every warning it
  reports blocks completion."

**task_01, `typescript` module.**

- `clause.typescript.keep-type-errors-visible` becomes: "Keep TypeScript type
  errors visible; never hide one to make Verification pass."

**task_01, templates.** Insert one paragraph between the heading and
`{{artifact.rules}}`, with a blank line on each side.

- `guides/typescript-bun.md`: "These rules govern the repository's TypeScript
  sources and tests. Code in another language follows its own guide."
- `guides/bun.md`: "These rules govern the Bun workspace: the packages under
  the root `package.json`. A toolchain for another language in the same
  repository keeps its own commands."

**task_02, `core` module.**

- `clause.core.follow-dependency-workflow`: replace "Use the repository's
  declared package manager and lockfile workflow." with "Use each language's
  declared package manager and lockfile workflow." The rest of the clause
  stays.
- `clause.core.activate-matching-skills`: append, after its last sentence,
  "When a skill's default conflicts with a Baseline rule or a
  Repository-Specific Normative Rule, follow the rule."

**task_02, activation file.** `trigger.testing` changes its `when` from
"Writing or changing tests." to "Writing or changing TypeScript tests."

**task_03, templates.** Insert one paragraph immediately before
`{{artifact.rules}}`, with a blank line on each side.

- `guides/backend.md`: "These rules govern the repository's TypeScript backend
  workspace. A service or command written in another language follows its own
  guide."
- `guides/frontend.md`: "These rules govern the repository's web frontend
  workspace. A terminal interface follows its own guide."

**task_03, stated HTTP suggestion.**

- The profile's `httpContract.default` changes from `"Post-only"` to `"REST"`.
- In `contract-v1.json`, "the interactive workflow suggests Post-only until"
  becomes "the interactive workflow suggests REST until".
- In the public guide's table of suggested values, the HTTP contract row
  changes from `Post-only` to `REST`.

The template lines wrap at the width the neighbouring paragraphs use. The
formatter composition test reports a line the target formatter would rewrap.

### Left as they are, deliberately

- Every clause of the `backend` and `frontend` modules. The structural-clause
  retention test compares them with the Source Baseline text. A reworded
  `clause.backend.prohibit-generic-layers` failed it in the disposable clone.
- `trigger.production-code` and `trigger.debugging`. Their bundles hold
  language-neutral skills.
- Every per-skill trigger of the stack modules. Each already names its
  technology.
- The `coverage` claims of `rule.go.context-errors` and
  `rule.tui.synchronous-model-tests`. They are catalog records that no guide
  renders.
- The `auth.provider` decision and its `/api/auth/*` suggestion. The scope a
  guide renders is the value the Setup Manifest records.
- The `http.contract` decision declaration. Its default is already `REST`.
- No new clause. A mandatory clause needs a Source Baseline manifest row, and
  the regeneration maintains rows without creating them.

### API Contracts

1. API Contract: rendered guide `docs/agents/typescript-bun.md` — both parts
   open with their scope sentence; the Bun part carries the test script rule,
   the named package managers and the Verification-based warnings condition;
   the TypeScript part carries the single type-error obligation.
2. API Contract: rendered guides `docs/agents/agent-instructions.md` and
   `docs/agents/skill-dispatch.md` — the dependency clause binds each
   language, the skill clause ends with the rule-over-skill-default sentence,
   and the `trigger.testing` line reads "Writing or changing TypeScript
   tests."
3. API Contract: rendered guides `docs/agents/backend.md` and
   `docs/agents/frontend.md` — each carries its scope sentence before its
   rules, and every rule line is unchanged.
4. API Contract: `roundfix baseline plan` and `roundfix baseline update` — a
   Setup Manifest with a recorded HTTP Contract Decision resolves to that
   decision unchanged; a Setup Manifest without one reports `http.contract` as
   a new decision whose suggested mode is `REST` and adopts nothing.

## Coverage Map

- Goal 1 → Exact texts (task_01, `bun` module); Testing Approach 1.
- Goal 2 → Exact texts (task_01 templates, task_02, task_03 templates);
  Testing Approach 1, 2 and 3.
- Goal 3 → Exact texts (task_03, stated HTTP suggestion); Testing Approach 4.
- Goal 4 → Exact texts (task_02, `core` module); Testing Approach 2.
- Goal 5 → Testing Approach 1–4.
- Core Feature 1 → Exact texts (task_01); API Contract 1; Testing Approach 1.
- Core Feature 2 → Exact texts (task_02); API Contract 2; Testing Approach 2.
- Core Feature 3 → Exact texts (task_03, templates); API Contract 3; Testing
  Approach 3.
- Core Feature 4 → Exact texts (task_03, stated HTTP suggestion); API
  Contract 4; Testing Approach 4.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → Testing Approach 2.
- Success Metric 4 → Testing Approach 3.
- Success Metric 5 → Testing Approach 4.
- Success Metric 6 → Testing Approach 4.
- Success Metric 7 → Testing Approach 1, 2 and 3; the QA gate compares the
  clause identifiers of the five modules before and after.
- Success Metric 8 → Testing Approach 5.
- API Contract 1 → Exact texts (task_01).
- API Contract 2 → Exact texts (task_02).
- API Contract 3 → Exact texts (task_03, templates).
- API Contract 4 → Left as they are, deliberately (the `http.contract`
  declaration); Testing Approach 4.

## Integration Points

- **Spec 0193.** It edits `core.json`, `backend.json` and the profile's digest
  pin, and it adds a record of every clause's identifier and enforcement
  level and a duplicate-text check. This Spec adds no clause, changes no
  level and writes no sentence another clause holds, so both checks stay
  green. The two Specs regenerate the same derived files and are delivered
  one after the other; this Spec's Tasks raise versions from whatever the
  starting main holds.
- **Adopters.** A Managed Refresh reports each changed guide as an update of
  its managed entry. Clause identifiers are kept, so no retention mapping
  names a changed clause. A recorded HTTP Contract Decision is read from the
  Setup Manifest and is not re-suggested.
- **Repository-Specific Normative Rules.** Three adopters carry the test
  script rule by hand. The Baseline rule now says the same, and their rule
  stays where it is: a managed update never edits that carrier.
- **Vendored skills.** Some still teach the bare Bun runner, `npx` and `npm`.
  This Spec edits no skill. The rule-over-skill-default sentence is the
  Baseline's answer until their upstream changes.

## Testing Approach

All tests read the embedded catalog, a plan built in a temporary repository,
or repository files, and need no network. Each negative case is its own test.
The new tests call existing helpers of the Baseline and CLI test packages from
new files and edit no existing test file.

1. **TypeScript and Bun wording.** New
   `internal/baseline/stack_rule_wording_test.go`:
   - `TestTheTypeScriptAndBunGuideSaysWhatItGoverns` builds a Standard
     TypeScript Monorepo plan in a temporary repository, takes the postimage
     of `docs/agents/typescript-bun.md` and expects no finding from
     `stackWordingFindings` for the seven sentences of task_01 and the four
     replaced sentences.
   - `TestStackWordingCheckReportsTheWordingItReplaced` feeds the same
     function the four replaced sentences as a literal and expects one finding
     per required and per replaced sentence, and a finding for a guide that is
     absent.
   - `TestTheRewordedBunAndTypeScriptClausesKeepTheirForce` finds each of the
     four reworded clauses by identifier in its embedded module and compares
     its enforcement level with a literal: `prohibited` for the package
     manager prohibition and `mandatory` for the other three. It also expects
     an identifier that does not exist to be reported as absent.
2. **Core wording and the testing trigger.** New
   `internal/baseline/stack_core_wording_test.go`:
   - `TestTheCoreGuidesNameEachLanguageAndTheGoverningRule` checks the
     postimages of `docs/agents/agent-instructions.md` and
     `docs/agents/skill-dispatch.md` for the task_02 sentences, and the
     absence of the two replaced ones.
   - `TestTheCoreWordingCheckReportsTheWordingItReplaced` feeds the wording
     check the two replaced sentences as a literal and expects one finding per
     sentence.
   - `TestTheRewordedCoreClausesKeepTheirForce` expects both reworded core
     clauses present and `mandatory`.
   - `TestEveryTriggerThatDispatchesALanguageSkillNamesItsLanguage` runs
     `languageTriggerFindings` on the embedded activation file with the table
     `vitest` → TypeScript, `react` → React, `hono` → Hono, and expects no
     finding and at least three covered skills.
   - `TestATestingTriggerThatNamesNoLanguageIsReported` runs it on a literal
     activation with the replaced trigger and expects `trigger.testing` back.
   - `TestAComposedProfileDispatchesVitestOnlyForTypeScriptTests` renders the
     skill dispatch for the modules `core`, `go` and `typescript` and expects
     the scoped testing trigger and the Go testing trigger; it renders `core`
     and `go` and expects neither `vitest` nor `trigger.testing`.
3. **Backend and frontend scope.** New
   `internal/baseline/stack_scope_and_http_default_test.go`:
   `TestTheBackendAndFrontendGuidesSayWhatTheyGovern` checks the postimages of
   `docs/agents/backend.md` and `docs/agents/frontend.md` for their scope
   sentences. The existing `TestStandardTypeScriptStructuralClauseRetention`
   proves that no backend or frontend clause changed.
4. **HTTP suggestion.** The same new file:
   - `TestTheProfileHTTPDefaultMatchesTheDecisionCatalog` and
     `TestAProfileHTTPDefaultThatDisagreesIsReported`.
   - `TestARecordedHTTPContractDecisionSurvivesAnUpdate` writes a Setup
     Manifest for the Standard TypeScript Monorepo Profile with the mode
     `Post-only`, then `REST`, resolves the manifest input and expects the
     resolved state, the recorded mode, the recorded exceptions and no new
     decision.
   - `TestAnAbsentHTTPContractDecisionIsSuggestedAsREST` removes the decision
     from the manifest and expects `http.contract` among the new decisions
     with the mode `REST`, and absent from the resolved decisions.
   New `internal/cli/baseline_http_default_statement_test.go`:
   - `TestEveryStatementOfTheHTTPDefaultNamesTheCatalogDefault` reads the
     catalog default and compares the contract sentence and the public guide's
     table row with it.
   - `TestAStatementThatNamesAnotherHTTPDefaultIsReported` feeds the check
     both statements with `Post-only`.
5. **Regeneration and retention.** Each Task's Verification runs the existing
   `TestFormatterComposition`, `TestCatalogCompatibility`,
   `TestBaselinePlanCharacterization`, `TestBaselineCompatibilityCorpus` and
   `TestStandardTypeScriptStructuralClauseRetention`, then the read-only
   Managed Refresh, which exits non-zero while a plan is pending and writes
   nothing.

Reverting each corrected statement in the disposable clone, thirteen in all,
failed the matching test of items 1 to 4, and reverting each of the two
documentation statements failed the statement test.

## Build Order

1. The Bun and TypeScript clauses, the two scope sentences of their guide, and
   the wording check, task_01 (depends on: none).
2. The core dependency and skill clauses, the testing trigger, and the trigger
   check, task_02 (depends on: 1).
3. The backend and frontend scope sentences, the three statements of the HTTP
   suggestion, and the preservation tests, task_03 (depends on: 2).
4. Terminal QA, task_04 (depends on: 1, 2, 3).

The chain is serial because every Task rewrites the same digest pin, catalog
snapshots, plan goldens and Setup Manifest, task_02 and task_03 call the
wording check task_01 creates, and task_01 and task_03 both edit the template
index.

## Risks & Considerations

- **Delivery order with Spec 0193.** Both Specs edit `core.json`,
  `backend.json` and the profile. This Spec is delivered after Spec 0193 and
  regenerates from the state it leaves. A different order needs only the same
  regeneration.
- **A scope sentence is not a gate.** An Agent can still apply a backend rule
  to a Go service. The sentence removes the contradiction; binding a rule to a
  directory waits for a workspace render token.
- **The warnings condition.** The clause still binds only a repository whose
  Verification treats warnings as errors. Making it unconditional is a new
  obligation and is out of scope.
- **The trigger check's table is short.** It names three technology skills.
  A new bundle with another technology skill is caught only when the table
  gains its row. The table stands next to the check for that reason.
- **Hazard: generated files.** A Task must never hand-edit a golden, a pin or
  a rendered guide. The regeneration is the only writer, and the second
  refresh proves convergence.
- **Hazard: governed test files.** The helpers the new tests call live in
  governed test files. A Task adds new files and leaves those files alone.

## Decisions

- **Scope by clause where the clause collides, by guide sentence elsewhere.**
  See ADR-0190.
- **The rule-over-skill-default sentence joins an existing clause.** A new
  mandatory clause fails the catalog's Source Baseline check.
- **Only the testing trigger changes.** The other two TypeScript-owned bundles
  are language-neutral.
- **Three statements follow the decision catalog.** The catalog is the value
  the prompt and the update suggestion read, so it is the side that stays.
- **Rendered output is what the wording check reads.** A plan postimage is
  what an adopter receives, and it covers the module, the template and the
  renderer in one assertion.
- **Tests live in new files.** The existing Baseline plan and CLI
  documentation test files are Governed Paths and are not edited.
