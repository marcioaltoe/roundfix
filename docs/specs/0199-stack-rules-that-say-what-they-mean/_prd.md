---
spec: 0199-stack-rules-that-say-what-they-mean
status: active
created: 2026-09-30
surfaces: [backend, docs]
---

# Stack rules that say what they mean

The Context-Driven Baseline ships its TypeScript, Bun, backend and frontend
rules to every adopter as Normative Clauses, and an Agent obeys them
literally. An audit on 2026-09-30 found rules whose literal reading is wrong.
One tells the Agent to use "Bun-owned commands" for tests, which reads as the
bare Bun runner, while the profile's own Verification runs the package's `test`
script. Three adopters each wrote that correction by hand. Several rules bind
"the repository" when they mean one language, so a repository that composes
the TypeScript modules with Go or Rust reads a rule against its own toolchain,
and a Go test change dispatches a Vitest skill. One clause states two
obligations, one names a profile setting that does not exist, and three places
still say the suggested HTTP mode is Post-only while the product suggests REST.

This Spec corrects those statements and adds the checks that keep them
corrected. It adds no obligation, and it changes no clause identifier and no
enforcement level.

## Project Constraints

- Identifier strategy: not applicable — no new identifier. Clause, rule, guide,
  module, template and trigger identifiers are kept. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — Baseline assets and local files
  only; no credential is read and no network call is added. The Spec changes
  which HTTP mode the Baseline states as its suggestion to a repository that
  has recorded none. ADR-0063 keeps the HTTP Contract Decision the
  repository's own, and a recorded decision is not touched. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0190 (this Spec) makes a stack rule
  name the language or workspace it governs and states that a rule governs
  over a dispatched skill's default. ADR-0186 requires adopter-neutral wording
  that cites no Spec or ADR number, and every rewritten sentence follows it.
  ADR-0058 and ADR-0060 require every Source Baseline clause to stay accounted
  for, so clause identifiers are kept and the clauses retained by their bytes
  are not reworded. ADR-0059 keeps the generated guides formatter-stable.
  ADR-0061 keeps the Standard TypeScript Monorepo Profile opinionated, and
  this Spec removes none of its requirements. ADR-0063 keeps the HTTP Contract
  Decision repository-owned. ADR-0067 lets a repository-owned Baseline Profile
  compose the embedded modules, which is the composition the scope wording
  serves. ADR-0081 and ADR-0149 keep sanctioned regeneration outputs inside a
  grant, and ADR-0073 and ADR-0103 require this repository's Managed Refresh
  to apply in one transaction and converge to `current`. ADR-0130 has the
  audit judge only Governed Paths. This Spec's gate is bound by ADR-0080,
  ADR-0088, ADR-0091, ADR-0096, ADR-0104, ADR-0117, ADR-0155 and ADR-0156, and
  ADR-0093 and ADR-0094 check its consistency by citation and artifact
  presence. ADR-0099 makes retention accounting a mechanical comparison of
  clause identities, which the kept identifiers satisfy. ADR-0166 has the
  Daemon record a path a Task changed without declaring it, and ADR-0167 keeps
  the pre-PR Pull Request row from deciding a qualifying partial; both bind
  this Spec's gate. ADR-0187 (the Roundfix skill and the command reference) and
  ADR-0189 (an owned skill's version names its content) are active in this
  tree and do not apply: this Spec edits no skill, no skill version and no
  command reference. Six ADRs cite a listed ADR and do not apply, because this
  Spec touches none of their behaviors: ADR-0074 (hybrid ownership of
  repository-authored rules) cites ADR-0067; ADR-0097 (carrying a QA row
  forward) cites ADR-0080; ADR-0168 (the related-ADR gap) and ADR-0176
  (citation checks read authored text) cite ADR-0093; ADR-0173 (citations a
  History Relocation breaks) cites ADR-0073; ADR-0179 (an explicit empty
  `paths` list) cites ADR-0130, and this Spec's grant lists its paths.
  Reached only through those citations, and equally untouched: ADR-0177 (the
  Relocation Citation scan reads through a repository root). All hold. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30 for the Baseline source and its generated guides ("Autorizar os
  dois"), recorded in [_authorization.md](_authorization.md); bounded files:
  `internal/baseline/assets/modules/bun.json`,
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
  `docs/agents/setup-context.json`. Sanctioned regeneration:
  `make baseline-digests`. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`.

## Goals

- An Agent in a Bun workspace runs tests the way the profile's Verification
  runs them, without a repository-written correction.
- A repository that composes the TypeScript modules with another language's
  modules reads no rule that binds the other toolchain, and a test change in
  the other language dispatches no TypeScript testing skill.
- Every statement of the suggested HTTP mode says REST, and an adopter's
  recorded HTTP Contract Decision is unchanged by the next update.
- An Agent reads once, in the clause that requires the skills, that a rule
  governs when a dispatched skill's default disagrees with it.
- A later change that brings any corrected wording back fails a test before it
  ships.

## Core Features

1. **The Bun and TypeScript rules say what the profile runs.**
   - Tests run through the package's `test` script, `bun run test`, and never
     through a bare runner such as `bun test`.
   - The prohibition names the package managers and runner it means (npm,
     pnpm, yarn, npx), applies inside the Bun workspace, and says that another
     language's toolchain keeps its own package manager.
   - The type-error clause states one obligation: TypeScript type errors stay
     visible. Its lockfile half is dropped, because the dependency clause and
     the Bun clauses already state it.
   - The warnings clause names the repository's Verification as what treats
     warnings as errors. No profile has a setting by that name.
   - The TypeScript and Bun guide opens each of its two parts with one
     sentence that names what the part governs.
2. **The core rules hold in a repository with several languages.**
   - The dependency clause binds each language's declared package manager and
     lockfile workflow.
   - The clause that requires the skills says that a Baseline rule or a
     Repository-Specific Normative Rule governs over a dispatched skill's
     default.
   - The testing Skill Activation names TypeScript tests. A check refuses any
     Skill Activation whose bundle holds a skill for one technology and whose
     trigger does not name that technology.
3. **The backend and frontend guides name their workspace.** Each opens with
   one sentence of scope. Every backend and frontend clause stays
   byte-identical, because the Source Baseline retains those clauses by their
   text.
4. **The suggested HTTP mode is REST wherever it is stated.** The profile's
   declared default, the Baseline contract's sentence and the public guide's
   table agree with the decision catalog. A check compares them. A repository
   with a recorded HTTP Contract Decision keeps it, in either mode, and a
   repository with none is offered REST and nothing is adopted for it.

## Non-Goals / Out of Scope

- Any new obligation: clause force for the Go, Rust, CLI or TUI modules, an
  unconditional "warnings block" rule, and a built-in profile that composes a
  Go CLI with a Hono backend and a React frontend. The maintainer decided
  "Defeitos agora, regras novas depois"; the open Backlog Entry
  `docs/backlog/2026-09-30-stack-modules-need-rules-with-force.md` holds them.
- Rendering a workspace path into a guide. Workspaces are profile metadata
  that only capability checks read, and no render token carries them.
- Rewording any backend or frontend clause, including the one that prohibits
  generic layers.
- The duplicate backend clause, the duplicate-text check and the clause force
  record. Spec 0193 owns them.
- Setup membership of skills, the skill version floor and the release step.
  Spec 0195 owns them.
- The text of any skill. Spec 0192 owns the owned skills, and a separate Spec
  refreshes the upstream skill snapshot and follows upstream renames.
- Changing the owner of the production-code and debugging Skill Activations.
- Changing the Better Auth route exception or the name of its owner.
- Rewriting the Source Baseline corpus or the parity corpus.

## Success Metrics

1. Success Metric: the rendered TypeScript and Bun guide states the test
   script rule, the named package managers, the single type-error obligation
   and the Verification-based warnings condition, and none of the four
   replaced sentences; a guide that still holds a replaced sentence fails the
   wording check.
2. Success Metric: the rendered instruction and skill-dispatch guides state
   the per-language dependency rule and the rule-over-skill-default sentence.
3. Success Metric: a composition of the core, Go and TypeScript modules
   renders the testing trigger for TypeScript tests only, a composition
   without the TypeScript module renders no Vitest skill, and a trigger that
   dispatches a technology skill without naming the technology fails the
   trigger check.
4. Success Metric: the rendered backend and frontend guides each open with
   their scope sentence, and the existing structural-clause retention test
   still passes.
5. Success Metric: the profile default, the contract sentence and the public
   guide's table name the decision catalog's default mode, REST; a statement
   that names the other mode fails the statement check.
6. Success Metric: a Setup Manifest that records `Post-only` resolves to
   `Post-only` with its exceptions unchanged and no new decision, and a Setup
   Manifest with no HTTP Contract Decision is offered REST without adopting
   it.
7. Success Metric: each of the six reworded clauses keeps its identifier and
   its enforcement level, every backend and frontend clause keeps its text,
   and the five touched modules hold the same clause identifiers before and
   after this Spec.
8. Success Metric: after the sanctioned regeneration and this repository's
   Managed Refresh, a second refresh reports no file change.

## Declared breaks

Adopters receive these changes on their next `roundfix baseline update`. None
changes a command or a recorded decision.

- **Break — test command.** An Agent that read "Bun-owned commands for tests"
  as `bun test` now runs `bun run test`. This is the command the profile's
  Verification already declares. A repository whose `test` script is
  `bun test` runs the same runner as before.
- **Break — testing trigger.** In a repository with the TypeScript module, a
  test change in another language no longer dispatches the testing bundle. The
  core module's own trigger still dispatches the general testing skill for
  every test change.
- **Break — stated HTTP suggestion.** Three statements change from Post-only
  to REST. The prompt and the update suggestion already offered REST, because
  they read the decision catalog. No recorded decision changes.
- **Clarification — type errors.** The clause loses its lockfile half. The
  obligation stays in the dependency clause, which every profile carries, and
  in the Bun clauses.
- **Clarification — warnings.** The condition names the repository's
  Verification in place of a profile setting that never existed. The
  obligation is the same.
- **Clarification — scope.** "Inside the Bun workspace", "each language's" and
  the four scope sentences change nothing in a repository with one language.
- **Clarification — rule over skill default.** A skill never had authority
  over a rule. The sentence states it where the skills are required.

## Recorded limits

- The scope is stated in words. No check reads a workspace path, so a rule
  cannot yet be bound to a directory.
- The wording checks prove that a sentence is present or absent. They cannot
  prove that a sentence is true.
- The production-code and debugging Skill Activations are owned by the
  TypeScript module and dispatch language-neutral skills. In a composed
  repository they fire for every language, which is correct. A profile without
  the TypeScript module does not receive them; giving them to it is a new
  obligation.
- The Go context-and-errors rule claims backend coverage and the TUI
  synchronous-model-tests rule claims frontend coverage. These are catalog
  records that no guide renders, and they are unchanged.
- The clause that prohibits generic layers can be read against a
  domain-services folder. It is retained by its bytes and stays as it is.
- The Better Auth route scope is a recorded repository decision, and
  `/api/auth/*` is only its suggestion. A repository changes it by replanning
  with its own decision, which the public guide already says. The owner name
  is fixed to the provider.
- Vendored skills still teach the bare Bun runner and other package managers.
  Their text belongs to their upstream; the rule-over-skill-default sentence
  is what governs until they change.

The first, third, fifth and sixth limits belong to the wave that the open
Backlog Entry named under Non-Goals describes. That entry already names stack
rules with force and a repository with several stacks, and this Spec does not
edit it.

## Decisions

- **Defects now, new rules later.** Maintainer decision of 2026-09-30.
- **Scope in words, in two places.** A clause that collides across languages
  carries its scope; a guide sentence scopes the rest. See ADR-0190.
- **No new clause identifier.** The rule-over-skill-default sentence joins the
  clause that requires the skills. A new mandatory clause needs a Source
  Baseline row that the regeneration never creates.
- **The smallest trigger change.** Only the testing Skill Activation dispatches
  a skill for one technology under a trigger that names none.
- **REST as the stated suggestion.** Maintainer decision of 2026-09-30 ("REST
  como padrão"). The decision catalog already suggested REST; the three stale
  statements follow it.
- **Guard in code, not in prose.** Each corrected class gets a check and a
  test that feeds the check the wording it replaced.

## Acceptance evidence

Each Core Feature requires positive and negative evidence in the Task Graph.
The outside-evidence row rests on sources this Spec did not produce:

- Adopter repositories, read through their Secondbrain mirrors on 2026-09-30.
  The Repository-Specific Normative Rules of conexus, fluxus and vortex each
  carry the section "Tests run through the package script, never through
  `bun test`". The backend guides of conexus, fiscus, fluxus, gss, oraculum,
  tax-poc and vortex each record "Application HTTP mode: **REST**", and none
  records Post-only.
- Bun's documentation: a built-in command takes precedence over a package
  script of the same name, and `bun run <script>` runs the script
  (<https://bun.com/docs/runtime>). The Bun maintainers' discussion states the
  same for `bun test` (<https://github.com/oven-sh/bun/discussions/26312>).
- The AGENTS.md convention for a repository with several packages: the
  instructions nearest the edited file govern
  (<https://developers.openai.com/codex/guides/agents-md>).

## Research basis

The Secondbrain was consulted through `wiki/index.md`, the adopter mirrors
under `projects/*/mirror/docs/agents/`, the inbox entries
`inbox/skills/2026-08-24-o-guia-de-backend-fixa-api-auth-como-escopo-do-better-auth.md`
and
`inbox/skills/2026-09-30-skills-vendorizadas-contradizem-o-baseline-ou-estao-quebradas.md`,
the two research entries of 2026-09-30 under `inbox/secondbrain/`, and the
query
`qmd query "stack rules scope language polyglot repository bun test package script REST default skill precedence"`.
The mirrors showed three adopters repeating the test-script rule and seven
recording REST. The first inbox entry asks that the Better Auth scope come
from the Setup Manifest, which it already does. The second lists vendored
skills that teach the bare Bun runner, `npx` and `npm`, and records the
maintainer's decision that repository rules prevail over skill defaults. The
two research entries cover models and orchestrators and hold nothing on stack
rules. Exa located the Bun runtime documentation, the Bun discussion and issue
on `bun test` against the `test` script, and the AGENTS.md precedence
convention. They support naming the script in the rule and stating scope and
precedence in the instruction text.

## Technical candidate

The corrections were applied in a disposable clone of this repository on
2026-09-30. The sanctioned regeneration, this repository's Managed Refresh, a
second refresh with no file change, the Baseline and CLI test packages and the
documentation gate all passed. Reverting each corrected statement in turn
failed the matching new test. Two approaches were measured and rejected there:
a new clause for the rule-over-skill-default sentence fails the Source
Baseline's required-clause check, and rewording a backend clause fails the
structural-clause retention test. The TechSpec records the exact texts.
