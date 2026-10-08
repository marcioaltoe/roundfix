---
spec: 0252-one-gate-per-run-and-a-smaller-baseline
prd: _prd.md
created: 2026-10-08
---

# One gate per Run and a smaller Baseline — Technical Spec

## Executive Summary

This Spec changes Baseline content, not Roundfix code. Four serial Tasks edit
the `core`, `spec-workflow`, `autonomous-work`, `secondbrain` and
`context-workflow` modules, the decision catalog, one guide template and two
root templates. Each Task then runs the Module Version Record step,
`make baseline-digests` and this repository's Managed Refresh. task_02 also
re-records four of this repository's decisions through `roundfix baseline plan`
and `roundfix baseline apply`. Two owned skills, `implement-task` and
`setup-context-driven`, change with their versions, and the user guide and the
glossary follow.

The trade-off accepted is ADR-0257's narrowing of ADR-0222: a clause whose
wording changes keeps its identity. An authoring prototype gave the reworded
clauses new identities, and `make baseline-digests` then refused at
`TestReadoptionCompatibilityMaintainedFixture`. It reported eight
`catalog.sourceBaseline.required-clause.missing` rows ("the regenerator
maintains manifest rows but never creates them") and one
`catalog.transition.target.unknown` for `clause.core.run-selected-verification`.
In-place rewording needs neither. A removed clause is named in the `replaces`
list of the clause that absorbs it, which has the same enforcement.

## Project Constraints

- Identifier strategy: not applicable. No identifier scheme changes. Every
  reworded clause keeps its `id`, and every decision keeps its `id`. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable. No request, credential or network
  call is added. Tests use the embedded catalog and temporary repositories.
  The Managed Refresh and `baseline plan`/`apply` run locally against this
  repository. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0257 (this Spec) governs every
  decision below. ADR-0222: "A clause whose wording changes takes a new
  identity that declares `replaces` for the old one"; ADR-0257 narrows it to
  removed clauses, with the measured cost above. ADR-0186: "Every clause
  keeps its enforcement level when its wording changes", and every merge
  joins clauses of equal enforcement. ADR-0186 also says "It never cites a
  Spec number or an ADR number", and no new clause text does. ADR-0250: "A
  Baseline module's version names one content, and the record step chooses
  it". ADR-0238: "A Task is `light` when its `complexity` is
  `low`, its `type` is not `qa`, and it declares no Governed Path".
  ADR-0252: "`make verify-changed` is the Verification of Runs and of the QA
  gate". ADR-0253: "CI runs the Full Contract Run on every push to main and in
  the release workflow before any build or publication". ADR-0244: "A Spec declares the domain terms it introduces". Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable. Express maintainer authorization: "considere
  autorizado a ajustar todas as skills se necessário", "Autorizar os dois",
  the standing "Concedo", and the 2026-10-08 grants for the governed paths this
  Spec declares. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`. Spec-contained authorization record:
  `docs/specs/0252-one-gate-per-run-and-a-smaller-baseline/_authorization.md`.
  Bounded files: `.agents/skills/implement-task/SKILL.md`,
  `.agents/skills/setup-context-driven/SKILL.md`,
  `docs/agents/agent-instructions.md`, `docs/agents/autonomous-work.md`,
  `docs/agents/secondbrain.md`, `docs/agents/setup-context.json`,
  `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`,
  `internal/baseline/assets/decisions.json`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/AGENTS.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/issue-tracker.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/secondbrain.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`,
  `internal/baseline/assets/modules/autonomous-work.json`,
  `internal/baseline/assets/modules/context-workflow.json`,
  `internal/baseline/assets/modules/core.json`,
  `internal/baseline/assets/modules/secondbrain.json`,
  `internal/baseline/assets/modules/spec-workflow.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/baseline/assets/templates/guides/autonomous-work.md`,
  `internal/baseline/assets/templates/index.json`,
  `internal/baseline/assets/templates/root/context-workflow.md`,
  `internal/baseline/assets/templates/root/secondbrain.md`,
  `internal/baseline/plan_test.go`, `internal/cli/baseline_human_test.go`,
  `skills/baseline_skill_contract_test.go`,
  `skills/implement-task/SKILL.md`, `skills/setup-context-driven/SKILL.md`.

## Current behavior

Measured on 2026-10-08 at `55de2a73` with the worktree's `bin/roundfix`:

- `docs/agents/agent-instructions.md` names
  "The selected repository Verification is `rtk make verify`." and
  "The selected incremental Verification is `rtk make verify-incremental`.".
  `docs/agents/spec-routing.md` says "For each Task, run the selected
  incremental Verification ... before handoff".
- `make verify-incremental` is `fmt-check vet test skills-sync-check
  skills-check build`, the full suite. CI runs `make verify-changed` on a
  Pull Request and `make verify`, `make verify-docs` and
  `make verify-contracts` on a push to main
  (`.github/workflows/ci-verify.yml`). `.roundfixrc.yml` sets
  `defaults.verification: make verify-changed`.
- `implement-task` §7 step 3 says: "Run focused implementation checks while
  working when useful. Do not run any command from the Task's
  `## Verification` section." It does not name the selected commands.
- `docs/agents/autonomous-work.md` opens with "Default backend work uses
  `codex gpt-5.6-sol`" and "`claude opus 5 xhigh`". The Project Config
  profiles use `gpt-6.1-sol`/high, `opus`/high and `gpt-5.6-luna`/max.
- The rendered guides plus `AGENTS.md` hold 69,330 bytes. `docs-layout.md`
  holds 16,642 and `secondbrain.md` 7,036. The root block says
  "Domain and documentation rules are mandatory" for both the domain and
  docs-layout guides, and "Optional cross-project knowledge follows" the
  Secondbrain guide.
- `roundfix baseline plan` with every recorded decision answered and
  `preservation.mode=preservation`, followed by `roundfix baseline apply`,
  changes exactly the decisions it is given. A prototype changed only
  `docs/agents/agent-instructions.md` and `docs/agents/setup-context.json`
  when only the two Verification values differed, and a following
  `roundfix baseline update --repo . --no-skills` reported `current`.

## System Architecture

```text
internal/baseline/assets/modules/*.json ─┐
internal/baseline/assets/decisions.json ─┼─> Module Version Record step -> module-versions.json
internal/baseline/assets/templates/**   ─┘   make baseline-digests     -> formatter goldens, profile digest, testdata
                                              baseline update --repo .  -> AGENTS.md, docs/agents/*, setup-context.json
                                              baseline plan + apply     -> recorded decisions (task_02 only)
.agents/skills/<owned>/SKILL.md -> skill version record -> make skills-sync -> skills/<owned>/SKILL.md
```

No Go source outside tests changes.

## Implementation Design

### Interfaces

No exported or unexported production symbol changes. Each new test is a
package `baseline` test that reads the embedded catalog through
`mustEmbeddedCatalog`, `catalog.Module` and `catalog.Asset`, and builds an
adopter through `newClauseReplacementAdopter`, as
`internal/baseline/glossary_clauses_test.go` does.

### Data Models

The decision catalog keeps its schema. Four entries change their value and
rise one version:

| decision | field | before | after |
| --- | --- | --- | --- |
| `verification.gate` | `default` | `rtk make verify` | `make verify` |
| `verification.incremental` | `suggestion` | `rtk make verify-incremental` | `make verify-incremental` |
| `runtime.backend` | `default` | `codex gpt-5.6-sol` | `codex gpt-6.1-sol high` |
| `runtime.design` | `default` | `claude opus 5 xhigh` | `claude opus high` |

The `summary` of `runtime.backend` becomes "The backend ACP Runtime, model,
and reasoning effort the repository prefers, named by generated guidance."
The summary of `runtime.design` becomes "The design, UI, UX, or frontend ACP
Runtime, model, and reasoning effort the repository prefers, named by
generated guidance."

This repository's Setup Manifest records:

| decision | before | after |
| --- | --- | --- |
| `verification.gate` | `rtk make verify` | `make verify` |
| `verification.incremental` | `rtk make verify-incremental` | `make verify-changed` |
| `runtime.backend` | `codex gpt-5.6-sol` | `codex gpt-6.1-sol high` |
| `runtime.design` | `claude opus 5 xhigh` | `claude opus high` |

### API Contracts

1. API Contract: the Baseline clauses. The embedded catalog carries the texts
   of Exact clause texts under the identities named there. The three removed
   identities appear in no module, and each is named in exactly one `replaces`
   list of a clause with the same enforcement. Each rendered guide states each
   clause once as a forced bullet.
2. API Contract: the decision catalog. `decisions.json` carries the values of
   Data Models. `roundfix baseline update --adopt-suggested` over a Setup
   Manifest without `verification.incremental` reports `make verify-incremental`
   as the suggested value and records it. The interactive defaults for
   `verification.gate`, `runtime.backend` and `runtime.design` are the new
   defaults. A recorded value always wins over a default, as today.
3. API Contract: the rendered guidance. A Managed Refresh renders the
   runtime header and the two root sentences of Template texts. The root
   blocks keep their markers, identities and references.
4. API Contract: the Managed Refresh of an adopter. For a Source Baseline
   adopter it is ready and records the dispositions of Retention, with no
   `unaccounted` clause.
5. API Contract: the `implement-task` handoff. In a Daemon-assigned turn, §7
   allows focused checks of the changed packages only. It names the Task's
   Verification, the repository's selected Verification, its incremental
   Verification and the full suite as commands the Agent does not run.

### Surface Transcripts

None. No command gains, loses or reorders an output line, a flag or an exit
code. The only observable changes are catalog values that existing result
fields already carry (API Contract 2) and the bytes of managed guidance (API
Contract 3). Tests assert both on those fields and on rendered files.

### Exact clause texts

Each text below is the clause's whole `guidance`. Every other byte of the
module stays as it is, except the version lines in Version changes and the
`replaces` lists named here. Each clause keeps its `id` and its enforcement.

task_01, `core`, `rule.core.verification-selected`:

`clause.core.run-selected-verification` (mandatory):

```text
Outside a Run, run the selected repository Verification before a completion claim, treat every failure as blocking, and report the command plus its actionable diagnostic. In a Daemon-assigned turn the Daemon owns that claim: it runs the Task's declared Verification and the repository Verification at settlement.
```

`clause.core.verification-two-tiers` (mandatory). The text from "CI must run"
to the end is today's text, unchanged:

```text
Outside a Run, use the selected incremental Verification named at the top of this guide for fast local checks; it answers whether the current change remains valid while reusing safe local state. In a Daemon-assigned turn, run focused tests of the changed packages and neither selected command, because the Daemon verifies the Task after handoff. CI must run the selected repository Verification from a fresh run; it answers whether the complete tree satisfies the repository contract. Baseline planning refuses until the repository selects and declares both commands, so neither tier is ever satisfied by omission. Execute authored Verification only on committed provenance: the Spec artifacts that carry its commands must be tracked in the repository and byte-identical to their committed bytes at the resolved revision, comparing the authored projection while excluding Daemon-owned `status` and `## Result` fields. A Run satisfies this by construction, because its Run Worktree is created from a commit. No other command checks it: `roundfix spec check --run-verification` executes the commands of the working-tree Spec in a checkout of `HEAD`, and `roundfix settle` reads the Task file of the directory it settles, so whoever runs them on an uncommitted or modified Spec first reads the commands they will execute. No approval record makes an untrusted source executable; commit the artifacts or do not execute them. Read-only checking remains available for any source.
```

task_01, `spec-workflow`, `rule.spec.routing`:

`clause.spec.verification-two-tiers` (mandatory):

```text
Outside a Run, run the selected incremental Verification named in `docs/agents/agent-instructions.md` for each Task to answer whether the current slice remains valid before handoff; a Daemon-assigned Task hands back after focused tests. CI must run the selected repository Verification from a fresh run to answer whether the assembled tree satisfies the repository contract. A missing incremental selection is a Baseline decision to answer, never a license to skip the local tier or a waiver to repeat in each Spec.
```

`clause.spec.verification-fails-before-the-change` (mandatory):

```text
Author each Task's Verification so every command fails on the tree before the Task's change, and run `roundfix spec check <slug> --strict --run-verification` before a Run starts: strict mode also fails on gap findings, and the Daemon refuses a Task whose command already exits zero on the unchanged tree.
```

task_03, `core`: `clause.core.request-pull-request-review` is removed from
`rule.core.evidence-before-complete`. `clause.core.request-review-explicitly`
in `rule.core.git-delivery` (mandatory) gains
`"replaces": ["clause.core.request-pull-request-review"]` and this text:

```text
Resolve and follow the repository's pre-PR review policy before publication: `codex`, `claude`, `coderabbit`, or explicit `none`, with Codex as the built-in default, User Config over built-in defaults, and Project Config over User Config. Never re-enable an opted-out provider; CodeRabbit is optional, never a mandatory dependency. With a provider enabled, obtain its review of the current candidate, with its configured model and effort, before opening the Pull Request, and block on failed, unavailable, incomplete, absent, or stale review evidence, which never selects `none` automatically. With explicit `none`, request no review, record that review was disabled by configuration, and keep QA and required checks; disabled review is not a passing review. Neither choice waives repository or GitHub requirements, and unimplemented configuration or runtime support is never presented as a working command.
```

`clause.core.prohibit-external-research-for-local-code` in
`rule.core.research-authority` (prohibited) gains
`"replaces": ["clause.secondbrain.prohibit-external-local-discovery"]` and
this text:

```text
Do not use external research tools to discover or infer local repository code or behavior, and never put credentials, private client records, or proprietary source code in an external search query.
```

task_03, `autonomous-work`, `rule.autonomous.loop`:
`clause.autonomous.loop-05-clean-is-not-evidence` (mandatory):

```text
Treat a terminal Clean and a resolved status as claims, not evidence: read the diff after a Run reports Clean and, when review is enabled, confirm that the selected provider reviewed the current candidate, because silence or a provider-reported skip is not approval. Preserve historical Review Skipped and Clean Unverified meanings.
```

task_03, `spec-workflow`: `clause.spec.status-only-in-task` is removed from
`rule.spec.artifacts`, and `clause.spec.tracker-artifacts` gains
`"replaces": ["clause.spec.status-only-in-task"]` with its text unchanged.
`clause.spec.project-constraints-05-legacy-and-ownership` (mandatory):

```text
Keep completed or archived legacy Specs byte-identical.
```

task_03, `secondbrain`: `clause.secondbrain.prohibit-external-local-discovery`
is removed from `rule.secondbrain.secret-safety`.
`clause.secondbrain.03-decision-consultation` (mandatory):

```text
Use Secondbrain for prior decisions, related project experience, and existing research.
```

A `replaces` list follows `enforcement` and precedes `guidance`, one target
per line, as in `clause.core.lint-warnings-block`.

### Template texts

task_02, `internal/baseline/assets/templates/guides/autonomous-work.md`, whole
file, ending in one newline:

```text
# Autonomous work

Agent Selection Profiles in Project Config choose each Agent Session's ACP
Runtime, model, and reasoning effort (`roundfix profiles show`). The
repository prefers {{runtime.backend}} for backend work and {{runtime.design}}
for design, UI, UX, and frontend-dominant work. A `complexity: low` Task that
is not `qa` and changes no Governed Path runs on the Light Tier when it is
available.

{{artifact.rules}}
```

task_04, `internal/baseline/assets/templates/root/context-workflow.md`, whole
file, ending in one newline:

```text
### CONTEXT-driven workflow

- Domain rules are mandatory: {{reference.domain}}. Read {{reference.docs-layout}}, whose rules then apply, before creating, changing, moving, or retiring `CONTEXT.md` or a document under `docs/`, and before a test or build step reads one.
```

task_04, `internal/baseline/assets/templates/root/secondbrain.md`, whole
file, ending in one newline:

```text
### Secondbrain

- Read {{reference.secondbrain}} before consulting or writing the Secondbrain or authoring an Idea, PRD, or TechSpec; a Run session, which cannot reach it, skips it.
```

Each template keeps its token list in `templates/index.json`.

### Skill and document texts

task_01, `.agents/skills/implement-task/SKILL.md` §7, step 3 becomes:

```text
3. Run focused implementation checks while working when useful, limited to the
   packages the change touches. Do not run any command from the Task's
   `## Verification` section, the repository's selected Verification, its
   incremental Verification, or the full test suite: the Daemon runs the
   declared Verification and the repository Verification at settlement.
```

No other line of the skill changes, and no `### QA settlement` section of any
skill changes.

task_01, `CONTEXT.md`, the **Incremental Verification** entry becomes:

```text
The selected fast local check recorded by the `verification.incremental` Baseline decision for validating the current change outside a Run while reusing safe local state. In a Daemon-assigned turn the Agent runs focused tests instead, because the Daemon runs the Task's Verification and the repository Verification at settlement (ADR-0257). It is distinct from the complete repository Verification selected by `verification.gate`.
```

Its `_Avoid_` line is unchanged.

task_02 changes two sentences in `.agents/skills/setup-context-driven/SKILL.md`.
The suggested incremental command changes from `rtk make verify-incremental`
to `make verify-incremental`. The backend and design suggestions change from
`codex gpt-5.6-sol` and `claude opus 5 xhigh` to `codex gpt-6.1-sol high`
and `claude opus high`, respectively.

task_02, `docs/user-guide/context-driven-development.md`: the suggestion table
rows become `make verify`, `make verify-incremental`,
`codex gpt-6.1-sol high` and `claude opus high`, and the sentence "to record
`rtk make verify-incremental`" becomes "to record `make verify-incremental`".
The decision-document example keeps its values, because it shows an explicit
document, not the catalog suggestions.

task_03, `docs/agents/specific-repository.md`: the bullet that begins
"HARD RULE — roundfix skill sync" is deleted whole. The CLI surface
clause "When the repository ships a skill or a command reference that
describes its commands, a change to command behavior updates that skill or
reference in the same pull request" already states the rule. No other byte
changes.

### Version changes

Each Task raises by one, from the value on its starting commit, every version
it touches in this list. It then runs the Module Version Record step, which
chooses the module's own version (ADR-0250), so no Task names a module
version number.

- task_01: `rule.core.verification-selected` and `guide.agent-instructions`
  in `core`; `rule.spec.routing` and `guide.spec-routing` in
  `spec-workflow`; the `implement-task` skill through the skill version
  record.
- task_02: decisions `verification.gate`, `verification.incremental`,
  `runtime.backend` and `runtime.design`; `template.guide.autonomous-work` in
  `templates/index.json`; `guide.autonomous-work` in `autonomous-work`; the
  `setup-context-driven` skill through the skill version record.
- task_03: `rule.core.evidence-before-complete`, `rule.core.git-delivery`,
  `rule.core.research-authority` and `guide.agent-instructions` in `core`;
  `rule.autonomous.loop` and `guide.autonomous-work` in `autonomous-work`;
  `rule.spec.artifacts`, `guide.issue-tracker`,
  `rule.spec.project-constraints` and `guide.spec-routing` in
  `spec-workflow`; `rule.secondbrain.index-first`,
  `rule.secondbrain.secret-safety` and `guide.secondbrain` in `secondbrain`.
- task_04: `template.root.context-workflow` and `template.root.secondbrain`
  in `templates/index.json`; `root.context-workflow` in `context-workflow`
  and `root.secondbrain` in `secondbrain`.

### Retention

For an adopter whose Setup Manifest declares the Standard TypeScript Source
Baseline, the Managed Refresh records each reworded clause `retained`, because
its identity and enforcement are unchanged. It records
`clause.core.request-pull-request-review`, `clause.spec.status-only-in-task`
and `clause.secondbrain.prohibit-external-local-discovery` `replaced`, each
with exactly one successor of the same enforcement. No clause is
`unaccounted`. The Source Baseline corpus, manifest, accounting and index, and
`retention/transition.managed-v2-to-portable-v3.json`, stay byte-identical.

### Derived files

The authoring prototype measured, Task by Task in order, the files the
sanctioned commands rewrite:

- every module Task: `internal/baseline/module-versions.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/baseline/testdata/catalog.diagnostics.golden.json`,
  `catalog.digest`, `catalog.normalized.json` and the four
  `testdata/plan-characterization/*.golden.json`, plus
  `docs/agents/setup-context.json`;
- task_01: the Standard TypeScript golden `docs/agents/agent-instructions.md`
  and `docs/agents/spec-routing.md`, and this repository's same two guides;
- task_02: the golden `docs/agents/autonomous-work.md`, and this repository's
  `docs/agents/agent-instructions.md` and `docs/agents/autonomous-work.md`;
- task_03: the goldens `agent-instructions.md`, `autonomous-work.md`,
  `issue-tracker.md`, `secondbrain.md` and `spec-routing.md`, and this
  repository's same five guides;
- task_04: the golden `AGENTS.md` and this repository's `AGENTS.md`.

`CLAUDE.md` is a symbolic link to `AGENTS.md` and is not written.

### Invariants

```text
1. No clause identity is added; three are removed, each named in exactly one replaces list of a clause with the same enforcement.
2. No two Baseline clauses share text, and the force record lists exactly the catalog's clauses.
3. A Source Baseline adopter's Managed Refresh is ready and has no unaccounted clause.
4. The Source Baseline assets and the retention transition are byte-identical.
5. After each Task, roundfix baseline update --repo . --no-skills reports current, and a second refresh changes no file.
6. The runtime decisions stay required decisions of autonomous-work and keep their render bindings.
7. No catalog default or suggestion, and no rendered guide, names an rtk-prefixed command.
8. Every module version is the one the record step chose; no hand-written pin, golden or guide.
9. Every owned-skill change raises both version fields through the record command, and skills/ mirrors .agents/skills/.
10. No Makefile, .roundfixrc.yml, CI workflow, go.mod, Source Baseline asset or production Go file changes.
```

## Coverage Map

- Goal "Daemon-assigned Task runs focused tests" → API Contracts 1 and 5,
  Invariant 5.
- Goal "recorded values match CI" → API Contract 2, Invariant 7.
- Goal "runtime header defers to profiles" → API Contracts 2 and 3,
  Invariant 6.
- Goal "each rule once" → API Contracts 1 and 4, Invariants 1 to 4.
- Goal "root reading set" → API Contract 3.
- User Story 1 → API Contracts 1 and 5.
- User Story 2 → API Contract 2, Invariant 7.
- User Story 3 → API Contract 3, Invariant 6.
- User Story 4 → API Contract 1, Invariant 2.
- User Story 5 → API Contract 3.
- Core Feature 1 → API Contracts 1 and 5.
- Core Feature 2 → API Contract 2, Invariant 7.
- Core Feature 3 → API Contracts 2 and 3, Invariant 6.
- Core Feature 4 → API Contracts 1 and 4, Invariants 1 to 4.
- Core Feature 5 → API Contract 3.
- Core Feature 6 → Build Order 1.
- Success Metric 1 → API Contracts 1 and 5.
- Success Metric 2 → API Contract 2, Invariants 5 and 7.
- Success Metric 3 → API Contract 3, Invariant 6.
- Success Metric 4 → API Contract 4, Invariants 1 to 3.
- Success Metric 5 → API Contract 3.

## Testing Approach

- `internal/baseline/verification_tier_clauses_test.go` (task_01, new) holds
  the four task_01 texts as literals. It checks each text, enforcement and
  absent `replaces` in the catalog. It requires each text exactly once as a
  `mandatory` bullet in the Standard TypeScript golden guide that renders it.
  It requires an adopter's refresh to record all four `retained`.
- `internal/baseline/runtime_and_verification_decisions_test.go` (task_02,
  new) reads `decisions.json` from the embedded catalog. It requires the four
  values in Data Models, no default or suggestion beginning with `rtk `, and
  `runtime.backend` and `runtime.design` in the `autonomous-work` module's
  `requiredDecisions`. On whitespace-normalized text, it requires the golden
  `docs/agents/autonomous-work.md` to carry "Agent Selection Profiles in
  Project Config choose each Agent Session's ACP Runtime, model, and reasoning
  effort" and "runs on the Light Tier when it is available", and no longer
  "Default backend work uses". The golden renders its fixture's recorded
  runtimes, not the defaults. A plan never fills an absent runtime decision
  from its default: the prototype's go-cli-tui plan without them returned
  "required Baseline decisions are missing: runtime.backend, runtime.design".
  The interactive defaults are covered by the updated
  `TestHumanBaselineDecisionDefaults`.
- `internal/baseline/deduplicated_clauses_test.go` (task_03, new) holds the
  task_03 texts as literals. It requires the three removed identities to be
  absent from every module, each named in exactly one `replaces` list of a
  clause with the same enforcement, and an adopter's refresh to record each
  removed identity `replaced` with that successor and no clause
  `unaccounted`.
- `internal/baseline/root_reading_set_test.go` (task_04, new) requires the
  golden `AGENTS.md` to carry each root sentence once with its references
  rendered, and no longer to carry "Domain and documentation rules are
  mandatory" or "Optional cross-project knowledge follows".
- Declared breaks, each changed only to the new contract:
  `internal/baseline/promoted_spec_and_typescript_clauses_test.go` (task_01,
  the authoring clause sentence it pins);
  `internal/baseline/incremental_verification_test.go`,
  `internal/cli/baseline_human_test.go` and
  `internal/cli/baseline_incremental_verification_test.go` (task_02, the
  suggestion and the defaults); `internal/baseline/clause_characterization_test.go`,
  `internal/baseline/plan_test.go` (`TestStandardTypeScriptStructuralClauseRetention`),
  `internal/cli/baseline_update_test.go` (the structural-clause fleet
  fixture) and `skills/baseline_skill_contract_test.go`
  (`TestLegacySpecConstraintExemption`) (task_03). The prototype ran the
  `internal/baseline`, `internal/cli`, `internal/config`, `internal/speccheck`
  and `skills` packages, with and without the `docscontract` and
  `repocontract` tags, over all four Tasks. These were the only failures.

## Build Order

1. Tiers scoped by session: the four task_01 clauses, `implement-task` and the
   **Incremental Verification** entry.
2. Decision values and the runtime header (depends on: 1, because both
   change `docs/agents/agent-instructions.md`, `setup-context.json` and the
   module record).
3. One statement per rule (depends on: 2, because both change the
   `autonomous-work` module and the same derived files).
4. Root reading set (depends on: 3, because both change the `secondbrain`
   module and the same derived files).
5. QA gate (depends on: 1, 2, 3, 4).

## Risks & Considerations

- A Run's Task agent reads the guides of its Task's starting commit. task_02
  to task_04 already read task_01's scoped tiers, and the QA gate reads all
  four.
- Spec 0253 changes `spec-workflow`, `autonomous-work` and
  `rule.autonomous.loop` after this Spec. It rebases on this Spec's merge, and
  its record step takes the next version.
- An adopter that recorded `rtk make verify` or an old runtime keeps its value.
  Only absent decisions take the new defaults.
- The task_02 plan answers every recorded decision with its current value.
  A missed decision makes the plan exit 3 and write nothing.

## Glossary

- changes: **Incremental Verification**

## Decisions

- Reworded clauses keep their identity, and removed clauses are named in
  `replaces`; see ADR-0257.
- The secondbrain local-discovery prohibition merges into the core
  prohibition, because both are `prohibited`; the core `mandatory`
  local-search clause stays separate; see ADR-0257.
- This repository's decisions are re-recorded through `baseline plan` and
  `baseline apply`, the automation path, because `baseline update` never
  changes a recorded decision; see ADR-0257.
- The user guide's decision-document example keeps its explicit values; see
  ADR-0257.
