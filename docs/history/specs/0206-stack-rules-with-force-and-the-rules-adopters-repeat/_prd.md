---
spec: 0206-stack-rules-with-force-and-the-rules-adopters-repeat
status: archived
created: 2026-09-30
surfaces: [backend, docs]
archived: "2026-10-02"
source_slug: 0206-stack-rules-with-force-and-the-rules-adopters-repeat
---


# Stack rules with force, and the rules adopters repeat

The Context-Driven Baseline ships Normative Clauses with a stated force to
every adopter, but its Go, Rust, CLI and TUI modules still ship paragraphs with
no force at all. An Agent in a Go repository cannot tell whether "prefer the
standard library" is advice or a ban, and a Rust adopter receives no error
policy. Meanwhile adopters keep writing the same rules by hand: five
repositories require express authorization for every database mutation, five
require Verification that fails before the change, three say that a lint
warning blocks. This Spec gives every rule of the four stack and surface
modules a clause with force, adds the Rust error policy the maintainer
decided, and promotes into the Baseline the hand-written rules that recur, so
adopters stop re-deriving them and Agents read them with force.

## Project Constraints

- Identifier strategy: not applicable — the new identifiers are Baseline
  clause and rule identifiers in the existing `clause.<module>.<name>` and
  `rule.<module>.<name>` forms; no domain identifier is introduced. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — Baseline assets, tests and local
  files only; no credential is read and no network call is added. The
  promoted database clauses govern what an Agent may run against an adopter's
  database; they read and store no credential. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0202 (this Spec) makes every stack
  rule a clause with force and writes a preference as an obligation about its
  exception, and ADR-0203 (this Spec) admits a hand-written rule into the
  Baseline only when it recurs. ADR-0190 makes a stack guide name what it
  governs and has a rule govern over a dispatched skill's default, which the
  Go and Rust clauses rely on for the vendored skills that recommend Cobra,
  testify and `anyhow`. ADR-0186 requires adopter-neutral wording with no Spec
  or ADR number, and every new clause follows it. ADR-0058 and ADR-0060
  keep every Source Baseline clause accounted for, so the one removed clause
  is declared `replaced` and each new clause of a Standard TypeScript Monorepo
  module gains its Source Baseline row. ADR-0099 makes retention accounting a
  mechanical comparison of clause identities, which the declared replacement
  satisfies.
  ADR-0059 keeps the generated guides formatter-stable. ADR-0061 keeps the
  Standard TypeScript Monorepo opinionated, and ADR-0067 lets a repository
  profile compose the modules this Spec changes. ADR-0063 keeps the HTTP
  Contract Decision repository-owned; no clause here touches it. ADR-0081 and
  ADR-0149 keep the sanctioned regeneration outputs inside the grant, and
  ADR-0073 and ADR-0103 require this repository's Managed Refresh to apply in
  one transaction and converge. ADR-0130 has the audit judge only Governed
  Paths, and ADR-0179 is not needed because the grant lists its paths.
  ADR-0192 lets the delivery resolve a conflict confined to derived paths by
  regeneration, which the shared digest pin, snapshots and goldens rely on;
  ADR-0193 names prerequisites through a Task Graph field this repository has
  not shipped, so this Spec names them in its PRD and TechSpec instead. This
  Spec's gate is bound by ADR-0080, ADR-0088, ADR-0091, ADR-0096, ADR-0104,
  ADR-0117, ADR-0155, ADR-0156, ADR-0166, ADR-0167, ADR-0182, ADR-0194 and
  ADR-0195, and ADR-0093 and ADR-0094 check its consistency. ADR-0187 and
  ADR-0189 are active and do not apply: this Spec edits no skill, no skill
  version and no command reference. ADR-0191 (a setup snapshot follows its
  upstream by name) governs setup snapshots this Spec does not touch.
  ADR-0204 (a composed profile takes a composed setup) holds: this Spec
  removes the Bun rule from the composed profile's `requiredRules`, as from
  the Standard TypeScript Monorepo profile, and changes neither its setup nor
  its modules. ADR-0219 (a profile draft binds to the closest built-in profile)
  holds, because both profiles drop the same rule and their module lists, by
  which a draft's distance is measured, do not change. ADR-0205 (a recorded frontend layout)
  governs frontend clauses this Spec does not touch. ADR-0180, ADR-0181,
  ADR-0183, ADR-0184, ADR-0196, ADR-0197, ADR-0198, ADR-0199, ADR-0200 and
  ADR-0201 govern model selection, receipts, command transcripts, review,
  token accounting and the authoring judge, none of which this Spec changes.
  Five ADRs cite a listed ADR and do not apply, because this Spec touches none
  of their behaviors: ADR-0074 (hybrid ownership of repository-authored rules)
  cites ADR-0067, and this Spec leaves every adopter's own rules where they
  are; ADR-0097 (carrying a QA row forward) cites ADR-0080; ADR-0168 (the
  related-ADR gap) and ADR-0176 (citation checks read authored text) cite
  ADR-0093; ADR-0173 (citations a History Relocation breaks) cites ADR-0073.
  Reached only through those citations, and equally untouched: ADR-0177 (the
  Relocation Citation scan reads through a repository root). All hold.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30 for the Baseline source and its generated guides ("Autorizar os
  dois"), recorded in [_authorization.md](_authorization.md); bounded files:
  `internal/baseline/assets/modules/core.json`,
  `internal/baseline/assets/modules/bun.json`,
  `internal/baseline/assets/modules/go.json`,
  `internal/baseline/assets/modules/cli-surface.json`,
  `internal/baseline/assets/modules/tui-surface.json`,
  `internal/baseline/assets/modules/rust.json`,
  `internal/baseline/assets/modules/spec-workflow.json`,
  `internal/baseline/assets/modules/typescript.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/baseline/assets/profiles/go-cli-typescript-monorepo.json`,
  `internal/baseline/assets/profiles/rust-cli.json`,
  `internal/baseline/assets/retention/transition.legacy-typescript-bun-to-portable-v3.json`,
  `internal/baseline/assets/templates/index.json`,
  `internal/baseline/assets/templates/guides/rust.md`,
  `internal/baseline/assets/source-baselines/index.json`,
  `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/baseline.json`,
  `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/manifest.json`,
  `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/agent-instructions.md`,
  `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/spec-routing.md`,
  `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/docs-layout.md`,
  `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/typescript-bun.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/typescript-bun.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`,
  `internal/baseline/plan_test.go`, `docs/agents/agent-instructions.md`,
  `docs/agents/skill-dispatch.md`, `docs/agents/spec-routing.md`,
  `docs/agents/docs-layout.md`, `docs/agents/setup-context.json`. Sanctioned
  regeneration: `make baseline-digests`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- Every rule of the Go, Rust, CLI and TUI modules reaches an adopter as a
  clause whose force is stated, one obligation per clause.
- A Go adopter may add a third-party module when it records why, and no
  shipped Go clause bans a library by name.
- A Rust adopter reads a typed-error policy with force: typed errors in
  library and domain code, a type-erased error only at a binary's entry point,
  and no panic on a path a user can reach.
- A hand-written rule found in two or more adopter repositories, or tied to
  measured rework, reaches every adopter as a clause with force; every other
  candidate is declined with a recorded reason.
- An adopter whose Source Baseline carries the clause this Spec removes takes
  its next Baseline update without an unaccounted clause.

## User Stories

1. As an Agent working in a Go repository, I want each Go rule to state
   whether it is mandatory or prohibited, so that I do not read a preference
   as a ban or a ban as a preference.
2. As a maintainer of a Go repository, I want to add a third-party module when
   I record the reason, so that the Baseline does not force me to rewrite a
   library the standard library lacks.
3. As an Agent working in a Rust repository, I want the error policy stated
   with force, so that I return typed errors from library code and never
   panic on user input.
4. As a maintainer of an adopter repository, I want the rules my peers keep
   writing by hand to arrive with the Baseline, so that I stop copying them
   between repositories.
5. As a maintainer of a Standard TypeScript Monorepo repository, I want my
   next Baseline update to account for the warnings clause it replaces, so
   that the update is not refused.
6. As a maintainer of a Rust repository that has not adopted the Baseline, I
   want to know what adopting it would change, so that I can plan the
   adoption.

## Core Features

1. **Go, CLI and TUI rules carry force.**
   - Every Go rule is split into clauses with force: a thin `main` package, a
     standard-library preference whose exception is a recorded reason, module
     files changed only through the `go` command, context-first IO, an owner
     for every goroutine, errors wrapped with the failed operation, tests run
     by `go test` over observable behavior, and a cross-platform build for
     every change to a file with a build constraint.
   - An assertion or mocking library is a third-party module and follows the
     recorded-reason clause. No Go clause names a library it forbids.
   - Every CLI and TUI rule is split into clauses with force, and a new CLI
     clause requires a skill or command reference that describes the commands
     to change in the same pull request as the behavior it describes.
2. **Rust rules carry force, and the Rust error policy is stated.**
   - Every existing Rust rule is split into clauses with force, and the Cargo
     manifest and lockfile change through Cargo commands.
   - A new Rust rule states: a thin binary over the library crate, typed
     errors from library and domain code, a type-erased error only at a
     binary's entry point, and no `unwrap`, `expect` or `panic!` on a path
     user input, the environment, a file or an IO result can reach. Test code
     is exempt.
   - The Rust guide opens with one sentence that names what it governs.
3. **Recurring rules become core clauses.**
   - A lint warning fails Verification in every repository, and every warning
     Verification reports blocks completion. This clause replaces the Bun
     clause that blocked warnings only when Verification already treated them
     as errors, and it declares the replacement.
   - A flaky test is a blocking failure; a production hook that exists only
     for tests is prohibited; a generated file is regenerated, never
     hand-edited; a vendored skill is never edited in place.
   - A statement that changes data or schema in a database other than a
     disposable local one needs express authorization for that statement; an
     authorized write proves its predicate with a read and runs in an explicit
     transaction; a migration, script, seed or test never carries a change
     around that authorization.
4. **Recurring rules become Spec-workflow and TypeScript clauses.**
   - Every Task Verification command fails on the tree before the Task's
     change, and the author runs the probing Spec check before a Run starts.
   - A Task that answers a recorded finding names the finding and its date.
   - No test, fixture or build step reads a file under the Spec root or
     assumes a Spec is still active; a check of the Spec artifacts themselves
     is exempt.
   - A TypeScript fixture that stands for a stored row is typed from the
     schema's inferred row type.
5. **Every candidate rule has a recorded disposition.** The TechSpec lists
   each hand-written rule found in the adopter repositories, where it
   appears, and whether it is promoted or declined, with the reason.

## Non-Goals / Out of Scope

- A built-in profile that composes several stacks, the frontend layout
  decision, and the Go guide's scope sentence. Spec 0207 owns them; this Spec
  changes only the composed profile's `requiredRules` entry for the removed
  Bun rule.
- The Setup snapshots, upstream skill renames and skill membership. Spec 0200
  owns them. Editing any vendored skill, including the Go skills that
  recommend Cobra and testify and the Rust skill that prescribes `anyhow`.
- The dispatch trigger of `cut-release`, which its own Backlog Entry keeps.
- Scope sentences for the CLI and TUI guides. They govern a command or
  terminal surface in any language, so a language scope would be wrong.
- Clause force for the external-triage module, which the adopted entry does
  not name. Spec 0208 gives that module its clauses.
- Rendering workspace paths into guides, the owner of the production-code and
  debugging Skill Activations, the wording of the generic-layers prohibition,
  the auth owner name in the backend guide, and the TypeScript setup's unused
  skills. The adopted entry lists them; they remain open, and the report of
  this Spec names them for a new Backlog Entry.
- Editing any adopter repository or onioncry, and adopting the Baseline in
  onioncry. The adopter notice below records what each must change.
- Removing Roundfix's own prohibition of Cobra and testify, which stays a
  Repository-Specific Normative Rule of this repository.
- Any Makefile, lint, formatter, CI or `go.mod` change in this repository.

## Success Metrics

1. Success Metric: every rule of the Go, CLI, TUI and Rust modules carries
   clauses with an enforcement level, and none carries rule-level guidance;
   the rendered guides of the Go CLI/TUI and Rust CLI profiles label every
   rule line with its force.
2. Success Metric: the rendered Go guide states the recorded-reason exception
   and the cross-platform build, and names no library it forbids; a Go
   guidance text that forbids Cobra or testify fails the wording check.
3. Success Metric: the rendered Rust guide states the typed-error, entry-point
   and no-panic clauses with their force and opens with its scope sentence.
4. Success Metric: the rendered instruction guide of every built-in profile
   states the lint, flaky-test, test-hook, generated-file and three database
   clauses, the skill guide states the vendored-skill clause, and the Bun
   guide no longer states the conditional warnings clause.
5. Success Metric: a Standard TypeScript Monorepo adopter whose Source
   Baseline carries the conditional warnings clause gets a ready update plan
   that records it `replaced` by the core lint clause; without the
   replacement declaration the same plan is refused as `unaccounted`.
6. Success Metric: the rendered Spec-routing, docs-layout and TypeScript
   guides state the four Spec-workflow and TypeScript clauses.
7. Success Metric: the clause force record lists every new clause with its
   force, the Source Baseline carries one row for each new clause of a
   Standard TypeScript Monorepo module, and no two clauses share text.
8. Success Metric: after the sanctioned regeneration and this repository's
   Managed Refresh, a second refresh reports no file change.

## Declared breaks

Adopters receive these on their next `roundfix baseline update`. Each is a
Baseline change; none changes a command or a recorded decision.

- **Break — lint warnings block everywhere.** A repository whose Verification
  passes with lint warnings now reads a mandatory clause it does not meet. It
  must run its linters with their deny-warnings option. The conditional Bun
  clause `clause.bun.block-warnings-when-profile-treats-them-as-errors` is
  removed, with its rule, and the core lint clause declares that it replaces
  it, so its Source Baseline retention disposition is `replaced`. The
  Standard TypeScript Monorepo and the composed Go CLI with TypeScript
  Monorepo profiles both stop requiring the rule; the composed profile has no
  Source Baseline, so its adopters' update renders the core clause without
  retention accounting.
- **Break — new obligations.** The promoted core, Spec-workflow and
  TypeScript clauses are new obligations for every adopter whose profile
  selects their module. Each is already written by hand in at least two
  adopters or tied to measured rework.
- **Break — Go and Rust force.** Go, CLI, TUI and Rust rules that were advice
  become mandatory or prohibited clauses, and the Rust CLI profile requires a
  new rule. A Rust library that returns `anyhow` errors, or panics on user
  input, now violates a prohibited clause.
- **Clarification — Go dependencies.** "Use stdlib `testing`" becomes "run
  tests through `go test`"; an assertion library is allowed with a recorded
  reason. This loosens the shipped text; Roundfix's own ban stays in its
  Repository-Specific Normative Rules.
- **Clarification — test hooks.** The Go ban on production-only test hooks
  moves to core and now binds every language.

## Adopter notice

What each repository read for this Spec must do after its next Baseline
update, recorded here and not applied:

- oraculum and vortex run `oxlint` without `--deny-warnings` in Verification,
  with rules at `warn`; both must add `--deny-warnings` to the lint their
  Verification runs. vortex's commit hook already runs it, so its hook is
  stricter than its Verification today, which the hook-strictness clause
  forbids; the change cures both. fiscus has the same mismatch.
- gss states that warnings block but its Verification does not deny them.
- conexus, fluxus and tax-poc already deny lint warnings in Verification.
- The database clauses restate what oraculum, conexus, vortex, fluxus and gss
  already require; their hand-written copies stay until their maintainers
  remove them.
- **onioncry**, which does not adopt the Baseline now, was read without
  writing. Adopting the Rust CLI profile would change little in its source:
  its library already returns one typed error enum, it has no `anyhow`, and
  its non-test source calls no `unwrap`, `expect` or `panic!`. Four
  `debug_assert!(false, …)` fallbacks panic in debug and test builds only;
  they are not named by the clause and would be a review question. Its tests
  call `expect` freely, which the clause exempts. Its lint target already
  denies Clippy warnings. It declares a diagnostic crate its source never
  uses, which the dependency review clause would surface. The operator sends
  that note.

## Prerequisites

- Deliver after Specs 0200, 0207 and 0208. They edit Baseline modules,
  profiles and derived files this Spec also edits: 0200 the Go, Rust, core and
  TypeScript module skills, 0207 the Go module versions and the Go guide's
  scope sentence, 0208 the `core` and `typescript` skills, the setup names, the
  external-triage clauses, the force record and the Source Baseline entry
  count. This Spec's Verification reads none of their artifacts; the order
  avoids source conflicts in module files. Each Task raises versions from the
  value on its starting main. All three are on main since 2026-10-02
  (`18ef15eb`). Spec 0207 also shipped the composed profile
  `go-cli-typescript-monorepo`, which requires the Bun rule this Spec removes
  and selects the core, Go, CLI, TypeScript, Bun and Spec-workflow modules
  this Spec changes; task_01 owns its one-line `requiredRules` change, and
  Spec 0207 already added the Go guide to the parity list task_02 extends.

## Recorded limits

- The checks prove that a clause is present with its force and wording. They
  cannot prove that an adopter obeys it.
- No mechanical check reads an adopter's lint flags; the lint clause is
  enforced by the adopter's Verification and review.
- The database clauses bind what an Agent does; no Roundfix command can see a
  database statement.
- The promotion rule is applied by reading adopter repositories on
  2026-09-30. A rule an adopter adds later waits for the next promotion.
- This Spec adds no external-triage clause; Spec 0208 owns them.

## Decisions

- **Stdlib first is a preference with a recorded reason.** Maintainer
  decision of 2026-09-30. See ADR-0202.
- **Cobra and testify stay banned only in Roundfix.** Maintainer decision of
  2026-09-30; the vendored Go skills are governed by the rule, not edited.
- **Rust error policy with force.** Maintainer decision of 2026-09-30. See
  ADR-0202.
- **Promotion only on recurrence or measured rework.** Maintainer decision of
  2026-09-30. See ADR-0203.
- **Lint warnings block.** Maintainer decision of 2026-09-30. The clause
  lives in core so every stack reads it, and it replaces the conditional Bun
  clause.
- **onioncry is validated by reading only.** Maintainer decision of
  2026-09-30.
- **One Spec, four Tasks.** The work fits four implementation Tasks and a QA
  gate, so no split into a second Spec is proposed.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- Adopter repositories read on 2026-09-30 through their Secondbrain mirrors,
  and onioncry and argus read in place. The Repository-Specific Normative
  Rules of oraculum, conexus, vortex, fluxus and gss each require express
  authorization for a database mutation; those of oraculum, conexus, fluxus
  and pantheon require a Verification probe that fails first; conexus, gss
  and onioncry state that warnings block; oraculum and vortex lint without
  denying warnings. The TechSpec's candidate table names each file.
- Published sources: the Rust best-practice guides that restrict type-erased
  errors to binaries and forbid panics on user input
  (<https://canonical.github.io/rust-best-practices/error-and-panic-discipline.html>,
  <https://github.com/apollographql/rust-best-practices/blob/main/book/chapter_04.md>);
  thiserror's own guidance that it serves library-like code
  (<https://github.com/dtolnay/thiserror>); the linter options that turn a
  warning into a failing exit
  (<https://eslint.org/docs/latest/use/command-line-interface>,
  <https://doc.rust-lang.org/stable/clippy/configuration.html>, oxlint's
  `--deny-warnings` in <https://github.com/oxc-project/oxc/issues/1958>); and
  the Go proverb that prefers a little copying to a little dependency
  (<https://go-proverbs.github.io/>).

## Research basis

The Secondbrain was consulted through `wiki/index.md`, the adopter mirrors
under `projects/*/mirror/`, the inbox entry
`inbox/skills/2026-09-30-skills-vendorizadas-contradizem-o-baseline-ou-estao-quebradas.md`,
and the query
`qmd query "lint warnings block deny-warnings rust typed errors thiserror anyhow unwrap stdlib first go dependency recorded reason"`.
The mirrors gave the candidate inventory. The inbox entry records that the
vendored Go skills recommend Cobra and testify and that the Rust skill
prescribes `anyhow` at entry points, which the Go and Rust clauses govern
without editing them. The query surfaced a vortex Task that adds
`lint --deny-warnings` to its Verification to match its commit hook, which is
measured evidence for the lint clause. No pending Inbox Entry for this
repository was open. Exa located the Rust error guides, thiserror's README,
the ESLint, Clippy and oxlint warning options, and the Go proverbs page.

## Open Questions

None.
