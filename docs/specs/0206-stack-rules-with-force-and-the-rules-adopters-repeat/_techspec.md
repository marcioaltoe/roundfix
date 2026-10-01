---
spec: 0206-stack-rules-with-force-and-the-rules-adopters-repeat
prd: _prd.md
created: 2026-09-30
---

# Stack rules with force, and the rules adopters repeat — Technical Spec

## Executive Summary

Eight Baseline modules change and no production Go code does. The Go, CLI, TUI
and Rust modules move from rule-level guidance to clauses with force; the Rust
module gains one rule; core, Spec-workflow and TypeScript gain the recurring
adopter rules the candidate table below promotes; one Bun clause is removed and
declared replaced. Twelve of the new clauses belong to modules the Standard
TypeScript Monorepo selects, so each gains a Source Baseline row by hand, the
only thing the sanctioned regeneration cannot create. The trade-off accepted is
that every promoted rule is a new obligation for adopters that never wrote it,
bought with the evidence that at least two adopters, or one adopter's measured
rework, already needed it. Every text below was applied in a disposable clone
of `5f182757` on 2026-09-30, one commit per Task, and the Baseline, CLI,
skills and Spec-check test packages passed after each commit, with a second
Managed Refresh reporting `File changes: 0`.

## Project Constraints

- Identifier strategy: not applicable — new Baseline clause and rule
  identifiers follow the existing `clause.<module>.<name>` and
  `rule.<module>.<name>` forms; no domain identifier. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — Baseline assets, tests and local
  files only; no credential and no network call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0202 and ADR-0203 (this Spec)
  govern the clause force and the promotion rule. ADR-0190 has a rule govern
  over a dispatched skill's default and a stack guide name what it governs;
  the Rust guide gains that sentence here, and the Go guide's comes from
  Spec 0207. ADR-0186 requires adopter-neutral wording. ADR-0058 and ADR-0060
  keep every Source Baseline clause accounted for. ADR-0099 makes retention
  accounting a mechanical comparison of clause identities, which the declared
  replacement satisfies. ADR-0059 keeps generated guides formatter-stable.
  ADR-0061, ADR-0063 and ADR-0067 are the profile, HTTP and composition
  decisions this Spec leaves as they are. ADR-0073, ADR-0081, ADR-0103 and
  ADR-0149 govern the regeneration and the Managed Refresh. ADR-0130 has the
  audit judge only Governed Paths. ADR-0192 resolves a conflict confined to
  derived paths by regeneration, and ADR-0193's prerequisite field is not
  shipped, so the Build Order states the prerequisites. ADR-0080, ADR-0088,
  ADR-0091, ADR-0093, ADR-0094, ADR-0096, ADR-0104, ADR-0117, ADR-0155,
  ADR-0156, ADR-0166, ADR-0167, ADR-0182, ADR-0194 and ADR-0195 bind the gate
  and the checker. ADR-0187, ADR-0189, ADR-0191, ADR-0204 and ADR-0205 are
  active and do not apply: this Spec edits no skill, no setup snapshot, no
  composed profile and no frontend clause. ADR-0180, ADR-0181, ADR-0183,
  ADR-0184, ADR-0196, ADR-0197, ADR-0198, ADR-0199, ADR-0200 and ADR-0201 do
  not apply; the PRD's row gives the reason. All hold. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30, recorded in [_authorization.md](_authorization.md); bounded
  files: `internal/baseline/assets/modules/core.json`,
  `internal/baseline/assets/modules/bun.json`,
  `internal/baseline/assets/modules/go.json`,
  `internal/baseline/assets/modules/cli-surface.json`,
  `internal/baseline/assets/modules/tui-surface.json`,
  `internal/baseline/assets/modules/rust.json`,
  `internal/baseline/assets/modules/spec-workflow.json`,
  `internal/baseline/assets/modules/typescript.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
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
  `docs/agents/docs-layout.md`, `docs/agents/setup-context.json`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Candidate rules

Every hand-written rule found on 2026-09-30, read-only, in the Repository-
Specific Normative Rules and repository-authored guide text of the adopters,
plus onioncry and argus. `SB/<p>` is
`~/dev/secondbrain/projects/<p>/mirror`; line numbers are those files' lines
when read. The rule of ADR-0203 decides: at least two repositories, or one
with measured rework, and not machine- or domain-bound, not owned by a
formatter, linter or recorded decision, not already a clause.

| # | Candidate | Where it appears | Decision |
| --- | --- | --- | --- |
| 1 | Express authorization for each database mutation | `SB/oraculum/docs/agents/specific-repository.md:162-170`, `SB/conexus/docs/agents/specific-repository.md:140-150`, `SB/vortex/docs/agents/specific-repository.md:356-365`, `SB/fluxus/docs/agents/specific-repository.md:65-69`, `SB/gss/docs/agents/specific-repository.md:88-93` | Promote: `clause.core.ask-before-database-mutation` (five repositories) |
| 2 | Prove the predicate with a read; write inside an explicit transaction | oraculum `:172-174`, conexus `:152-154,166-167`, vortex `:367-370`, fluxus `:71,77` (same files) | Promote: `clause.core.prove-the-mutation-predicate` (four) |
| 3 | No mutation disguised as a migration, script, seed or test | oraculum `:176-177`, conexus `:156-164`, vortex `:372-380`, fluxus `:73-75` | Promote: `clause.core.prohibit-disguised-database-mutation` (four) |
| 4 | Database only through its command-line client, no MCP bridge | oraculum `:8-11`, conexus `:26`, vortex `:243-246`, fluxus `:31`, gss `:68-70` | Decline: a tool choice each repository owns; rows 1–3 carry the risk it guards |
| 5 | Credentials from the environment, never echoed | oraculum `:51-54`, conexus `:47-49`, vortex `:272-274`, fluxus `:41`, gss `:74-76`, `SB/pantheon/docs/agents/secrets.md:6` | Decline: `clause.core.prohibit-secret-exposure` already says it |
| 6 | Vault sourcing, Homebrew `mysql-client`, `MYSQL_PWD`, MySQL 5.5 dialect | oraculum `:13-49,76-81`, vortex `:248-303`, fluxus `:33-43,111-113`, gss `:71-73`, pantheon `specific-repository.md` | Decline: bound to one workstation or one domain |
| 7 | Generated files, such as migrations, are never hand-edited | vortex `:238-239`, conexus `:122-124`, `SB/fluxus/docs/agents/autonomy-charter.md:43-46`, `docs/agents/specific-repository.md` (this repository, digest pins) | Promote: `clause.core.regenerate-generated-files` (four) |
| 8 | Lint warnings block | `SB/conexus/docs/agents/typescript-bun.md:34-35`, `SB/gss/docs/agents/specific-repository.md:97-98`, `~/dev/onioncry/AGENTS.md:21-23,175-176`; measured in `SB/vortex/docs/history/specs/0075-a-classificacao-de-produtos-sai-do-espelho/task_04.md:41` | Promote: `clause.core.lint-warnings-block`, replacing the conditional Bun clause (maintainer decision; three repositories) |
| 9 | Flaky tests are blocking failures | vortex `:94`, gss `:128`, `~/dev/onioncry/AGENTS.md:149` | Promote: `clause.core.flaky-tests-block` (three) |
| 10 | No test-only hooks in production code | vortex `:92`, `~/dev/onioncry/AGENTS.md:146`, the Go module's own rule | Promote to core: `clause.core.prohibit-test-only-production-hooks` (three; the Go wording moves) |
| 11 | Upstream-managed skills are not edited locally | vortex `:457-460`, `docs/agents/specific-repository.md` (skill ownership), `~/dev/onioncry/AGENTS.md:253-258` | Promote: `clause.core.prohibit-editing-vendored-skills` (three) |
| 12 | A skill that describes the CLI ships with the behavior change | `docs/agents/specific-repository.md` (roundfix skill sync), `~/dev/onioncry/AGENTS.md:34-39` | Promote: `clause.cli.ship-the-skill-with-the-behavior` (two) |
| 13 | Verification probed against the untouched tree fails first | oraculum `:272-317`, conexus `:108-116`, fluxus `:11`, `SB/pantheon/docs/agents/specific-repository.md:100-102`, vortex `:434-447` | Promote: `clause.spec.verification-fails-before-the-change` (five) |
| 14 | A prior-cycle finding is named in the next Task's requirements | oraculum `:260-270`, conexus `:105-107`, fluxus `:9` | Promote: `clause.spec.name-the-finding-a-task-answers` (three) |
| 15 | Tests never read Spec artifacts | oraculum `:355-359` (nine tests red in CI, recurred); this repository's archive breaking a test pinned to an active Spec | Promote: `clause.spec.prohibit-tests-that-read-specs` (two, measured) |
| 16 | Fixtures typed from the schema's row type | fluxus `:59-63` (four wrong columns found once typed) | Promote: `clause.typescript.type-fixtures-from-the-schema` (measured) |
| 17 | Thin binary over the library | `docs/agents/specific-repository.md` (project map), `~/dev/onioncry/AGENTS.md:120-121` | Promote per stack: `clause.go.keep-entry-points-thin`, `clause.rust.keep-the-binary-thin` (two) |
| 18 | Stdlib first; no casual dependency | `docs/agents/specific-repository.md` (anti-pattern 1), `~/dev/onioncry/AGENTS.md:134` | Promote as a preference: `clause.go.record-the-reason-for-a-module` (maintainer decision); the Cobra and testify ban stays in this repository only |
| 19 | Manifests change through the package manager | `~/dev/onioncry/AGENTS.md:26-27`, Go and Bun modules | Promote for Rust: `clause.rust.change-manifests-through-cargo` |
| 20 | Typed errors, no panic on user paths | `~/dev/onioncry/AGENTS.md:123-130,225` | Promote: the Rust error rule (maintainer decision) |
| 21 | Build constrained files for every named platform | this repository: `docs/history/findings/2026-09-29-the-first-full-queue-trial-needed-manual-recovery.md:53` (a Windows defect found only in review cost a corrective Task) | Promote: `clause.go.build-every-constrained-platform` (measured) |
| 22 | Tests through the package script, never bare `bun test` | conexus `:84-92`, vortex `:384-390`, fluxus `:47-49` | Decline: `clause.bun.use-bun-owned-commands` already says it |
| 23 | Test behavior, not mocks | `~/dev/onioncry/AGENTS.md:147` | Decline: `clause.typescript.prohibit-incidental-test-oracles` covers TypeScript; Rust has one source |
| 24 | Arrange-Act-Assert, test file layout | vortex `:88-93`, gss `:118-120`, `~/dev/argus/CLAUDE.md:74` | Decline: style, and the layouts differ |
| 25 | No `any`; strict TypeScript | gss `:114,136`, `~/dev/argus/CLAUDE.md:76` (argus's linter turns the rule off) | Decline: the linter owns it, and with row 8 a `warn` rule now blocks |
| 26 | Function length cap | gss `:137`, conexus `typescript-bun.md:38-39` | Decline: the two numbers conflict; the linter owns it |
| 27 | Formatting style, naming, quotes | vortex `:464-470`, gss `:132-143`, argus `CLAUDE.md:73-75` | Decline: the formatter owns it |
| 28 | Architecture checker in the gate | stated by no guide; run by six repositories' Makefiles | Decline: Verification composition is repository-owned |
| 29 | Lean root manifest; self-contained `tsconfig` | vortex `:77-84`, gss `:103-112` | Decline: `clause.bun.add-from-owning-workspace` covers the first; the second is tooling configuration |
| 30 | Contract-first HTTP with Zod and OpenAPI | `SB/oraculum/docs/agents/backend.md:76-80`, vortex `:219-223` | Decline: the HTTP Contract Decision and `clause.backend.boundary-contracts` own it |
| 31 | Layering, thin handlers, feature systems with barrels | oraculum `backend.md:38-49`, vortex `:32-66`, argus `CLAUDE.md:72` | Decline: already backend and frontend clauses |
| 32 | Commit scope, destructive Git, one PR per Spec, review override, rerun gate after archive, two corrective Tasks | vortex, gss, oraculum, fiscus, pantheon, onioncry | Decline: already core, autonomous-work or loop clauses |
| 33 | Branch prefix `ma/` | `SB/vortex/docs/agents/agent-instructions.md:20-26`, `~/dev/onioncry/AGENTS.md:33` | Decline: contradicts the Baseline's purpose-based prefix |
| 34 | Version manifests agree before a tag | oraculum `:229-245` (measured) | Decline: the repository's release runbook owns release choreography |
| 35 | Shared surface becomes a `needs` edge; repair legacy QA labels | oraculum `:319-337`, conexus `:117-124`, fluxus `:13-15` | Decline: Roundfix enforces declared-file collisions, and the labels are a tool defect |
| 36 | `.env.example` mirrors the keys; English identifiers; Actions hygiene | gss `:162`, onioncry `:214-218`, pantheon `deploy.md:17,153`, vortex `:474-477` | Decline: covered by the secret and English clauses, or the rules differ |
| 37 | Race detector policy, gate composition, test isolation | this repository only | Decline: one repository, no measured rework |

## System Architecture

| Concern | Source of truth | Derived by |
| --- | --- | --- |
| Clause text, force, versions | `internal/baseline/assets/modules/*.json` | hand edit under the grant |
| Required rules | `internal/baseline/assets/profiles/*.json` `requiredRules` | hand edit under the grant |
| Legacy retention target | `internal/baseline/assets/retention/transition.legacy-typescript-bun-to-portable-v3.json` | hand edit under the grant |
| Rust guide scope sentence | `internal/baseline/assets/templates/guides/rust.md`, `templates/index.json` | hand edit under the grant |
| Source Baseline rows | the Source Baseline corpus files, `manifest.json` row, `index.json` `entryIds` | hand edit: text, identity, force, carrier; the regeneration fills offsets, digests and the identity record |
| Goldens, digest pin, catalog snapshots, plan goldens | the sources above | `make baseline-digests` |
| This repository's guides and Setup Manifest | the sources above | `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text` |

A module edit keeps the file's formatting: insert clause objects and raise
version numbers in place, never re-encode a module. Each Task edits its
sources and tests so the package compiles, runs `make baseline-digests`, then
the Managed Refresh; a second refresh reports `File changes: 0`.

Outputs measured in the rehearsal:

- task_01 rewrites the Standard TypeScript Monorepo digest pin, its goldens
  `agent-instructions.md`, `skill-dispatch.md` and `typescript-bun.md`, the
  Source Baseline `baseline.json`, `manifest.json` and `index.json`,
  `catalog.diagnostics.golden.json`, `catalog.digest`,
  `catalog.normalized.json`, the four plan goldens under
  `internal/baseline/testdata/plan-characterization/`
  (`advisory-only-divergences`, `clean-adoption`,
  `idempotent-replan-after-verified-apply`,
  `same-baseline-changed-profile-and-catalog-digests`), and this repository's
  `agent-instructions.md`, `skill-dispatch.md` and `setup-context.json`.
- task_02 rewrites `catalog.digest`, `catalog.normalized.json`, the four plan
  goldens and this repository's `go.md`, `cli.md`, `tui.md` and
  `setup-context.json`.
- task_03 rewrites `catalog.digest`, `catalog.normalized.json`, the four plan
  goldens and `setup-context.json`. This repository has no Rust guide.
- task_04 rewrites what task_01 rewrites, except that its goldens are
  `docs-layout.md`, `spec-routing.md` and `typescript-bun.md` and its
  repository guides are `docs-layout.md`, `spec-routing.md` and
  `setup-context.json`.

## Implementation Design

### Interfaces

No production file changes. Each Task adds one test file; helpers are
test-only pure functions so a negative test can feed them a literal.

```go
// internal/baseline/promoted_core_clauses_test.go (task_01)
type expectedClause struct{ module, id, force string }
func promotedCoreClauses() []expectedClause
func clauseForceFindings(catalog *Catalog, want []expectedClause) []string

// internal/baseline/stack_force_go_cli_tui_test.go (task_02)
func ruleLevelGuidanceFindings(moduleID string, module map[string]any) []string
func namedLibraryBanFindings(guide string, libraries []string) []string

// internal/baseline/stack_force_rust_test.go (task_03): uses the helpers above.
// internal/baseline/promoted_spec_and_typescript_clauses_test.go (task_04)
func sourceBaselineRowFindings(baseline SourceBaseline, want []expectedClause) []string
```

The files reuse `stackWordingFindings`, `clauseForce`, `planPostimage`,
`buildTestPlan`, `buildProjectDecisionPlan`, `newPlanRepository`,
`planTestDecisions`, `standardTypeScriptDecisions`,
`newClauseReplacementAdopter` and `mutateCatalogClause` from the package.

### Data Models

No schema change. `setup-context-driven/module-v3` already carries clauses;
the Go, CLI, TUI and Rust modules move to it and gain
`"repositoryExtensions": []`, as the Bun module has.

### Exact texts

Each clause is `id` — force — guidance, exact. The guidance of the replaced
rule-level text is split without other change unless shown.

**task_01, `core`.** `rule.core.root-cause-only`:
`clause.core.regenerate-generated-files` — prohibited — "Do not hand-edit a
generated file, such as a schema migration, a client generated from an API
description, or a derived digest; change its source and run the generator that
owns it." `rule.core.verification-integrity`:
`clause.core.flaky-tests-block` — mandatory — "A flaky test is a blocking
failure: find and fix the cause of its nondeterminism before completion, and
never rerun Verification until it happens to pass.";
`clause.core.lint-warnings-block` — mandatory, with
`"replaces": ["clause.bun.block-warnings-when-profile-treats-them-as-errors"]`
— "A lint warning fails Verification. Every linter the repository
Verification runs uses its option that turns a warning into a failure, such as
`oxlint --deny-warnings`, `biome check --error-on-warnings`,
`eslint --max-warnings 0`, or `cargo clippy -- -D warnings`, and every warning
Verification reports blocks completion.";
`clause.core.prohibit-test-only-production-hooks` — prohibited — "Do not add a
production hook, flag, branch, or exported symbol that exists only for tests;
test through the public entry points." `rule.core.skill-dispatch`:
`clause.core.prohibit-editing-vendored-skills` — prohibited — "Do not edit a
skill the repository installs from an upstream source; its upstream owns its
text. State a correction as a Repository-Specific Normative Rule, which
governs over the skill's default." `rule.core.security-configuration`:
`clause.core.ask-before-database-mutation` — stop-and-ask — "Stop and ask for
express authorization before any statement that changes data or schema in a
database other than a disposable local one. Name the exact statement and
database; an authorization covers only that statement.";
`clause.core.prove-the-mutation-predicate` — mandatory — "Before an authorized
database write, run a read with the same predicate and report its row count,
then run the write inside an explicit transaction and confirm the affected row
count before committing."; `clause.core.prohibit-disguised-database-mutation`
— prohibited — "Do not route a database change through a migration, script,
seed, or test to avoid that authorization. The repository's committed migrate
and seed commands run against a disposable local database are exempt."

**task_01, `bun` and retention.** Remove `rule.bun.warning-free-verification`
and its one clause, remove the rule from `guide.bun` and from the Standard
TypeScript Monorepo profile's `requiredRules`. In the legacy transition, the
mapping of `clause.legacy.block-warnings` keeps `replaced` and changes its
target to `clause.core.lint-warnings-block` and its reason to "The portable
core clause makes every lint warning fail Verification."

**task_02, `go`.** `rule.go.stdlib-first`:
`clause.go.change-module-files-through-go` — prohibited — "Do not hand-edit
`go.mod` or `go.sum`; change them through the `go` command.";
`clause.go.keep-entry-points-thin` — mandatory — "Keep each Go `main` package
thin: it parses input, wires dependencies, and calls behavior that lives in
cohesive packages."; `clause.go.record-the-reason-for-a-module` — mandatory —
"Prefer the Go standard library. Add a third-party module only for a named job
that the standard library and the modules already required cannot do, and
record that reason in an ADR or another recorded repository decision in the
change that adds it." `rule.go.context-errors`:
`clause.go.own-every-goroutine` — mandatory — "Give every goroutine an owner
that waits for it and a cancellation path that stops it.";
`clause.go.pass-context-first` — mandatory — "Pass a `context.Context` as the
first parameter of blocking and IO work, and stop that work when the context
is cancelled."; `clause.go.wrap-errors-with-the-operation` — mandatory — "Wrap
a returned error with the operation that failed using `%w`, and keep
`errors.Is` and `errors.As` matching wherever a caller branches on the
error." `rule.go.observable-tests`: `clause.go.build-every-constrained-platform`
— mandatory — "When a change touches a file with a build constraint, build
the affected non-test packages for every operating system its constraints
name, for example with `GOOS=windows go build ./...`; a build or test run on
the host compiles only the host's files."; `clause.go.test-observable-behavior`
— mandatory — "Test observable package and command behavior through public
entry points: stdout, stderr, files, exit codes, cancellation, and failure
paths."; `clause.go.test-through-go-test` — mandatory — "Run Go tests through
`go test` with the standard `testing` package as the harness. An assertion or
mocking library is a third-party module and needs its recorded reason."

**task_02, `cli-surface`.** `rule.cli.output-contract`:
`clause.cli.public-command-contract` — mandatory — "Treat command names,
flags, stdout and stderr placement, machine-readable fields, and exit codes as
public API."; `clause.cli.separate-output-streams` — mandatory — "Keep stdout
for requested output and stderr for diagnostics, progress, and warnings.";
`clause.cli.ship-the-skill-with-the-behavior` — mandatory — "When the
repository ships a skill or a command reference that describes its commands,
a change to command behavior updates that skill or reference in the same pull
request." `rule.cli.non-interactive`: `clause.cli.deterministic-non-interactive`
— mandatory — "Make automation deterministic and non-interactive.";
`clause.cli.explicit-safe-writes` — mandatory — "Make write operations
explicit, replayable, safe by default, and observable; use dry-run,
confirmation, or idempotency contracts where the repository requires them."

**task_02, `tui-surface`.** `rule.tui.synchronous-model-tests`:
`clause.tui.drive-models-synchronously` — mandatory — "Drive TUI model updates
synchronously and assert rendered state, messages, and transitions.";
`clause.tui.emulate-the-terminal-last` — mandatory — "Use terminal emulation
only when model-level tests cannot prove the behavior.";
`clause.tui.keep-design-policy-in-repository-guidance` — mandatory — "Keep
layout and interaction policy in repository-owned design guidance."

**task_03, `rust`.** `rule.rust.current-docs`:
`clause.rust.change-manifests-through-cargo` — mandatory — "Change
dependencies through Cargo commands such as `cargo add`, `cargo remove`, and
`cargo update`, and commit `Cargo.toml` and `Cargo.lock` together.";
`clause.rust.read-current-docs` — mandatory — "Use current authoritative crate
and toolchain documentation before changing Rust APIs, configuration, or
dependencies." `rule.rust.cargo-gates`:
`clause.rust.iterate-with-focused-cargo-commands` — mandatory — "Use focused
`cargo check`, `cargo test`, and lint commands while iterating, then run the
selected repository Verification."; `clause.rust.test-the-command-through-execution`
— mandatory — "Treat CLI flags, streams, JSON fields, and exit codes as public
behavior and test them through observable command execution." New rule
`rule.rust.library-and-entry-point`, version 1, coverage `coverage.rust` and
`coverage.cli`, listed in `guide.rust` and the Rust CLI profile's
`requiredRules`: `clause.rust.keep-the-binary-thin` — mandatory — "Keep each
binary target's `main` thin: it parses arguments and calls behavior that lives
in the library crate."; `clause.rust.keep-type-erased-errors-at-the-entry-point`
— prohibited — "Do not use `anyhow` or another type-erased error type outside
a binary's entry point. Only `main` and the command dispatch it calls may turn
typed errors into a type-erased report."; `clause.rust.prohibit-panics-on-user-paths`
— prohibited — "Do not call `unwrap`, `expect`, or `panic!` on a path that user
input, the environment, a file, or an IO result can reach; return a typed
error instead. Test code is not such a path."; `clause.rust.return-typed-errors`
— mandatory — "Return typed errors from library and domain code: an error enum
per boundary that implements `std::error::Error`, for example through
`thiserror`, so a caller can match its variants." The Rust guide template
gains, between its heading and the rules token: "These rules govern the
repository's Rust crates: the library and binary targets of its Cargo package
or workspace, and their tests. Code in another language follows its own
guide."

**task_04, `spec-workflow`.** `rule.spec.routing`:
`clause.spec.name-the-finding-a-task-answers` — mandatory — "When a Task
answers a recorded finding, name that finding and its date in one of the
Task's requirements."; `clause.spec.verification-fails-before-the-change` —
mandatory — "Author each Task's Verification so every command fails on the
tree before the Task's change, and run
`roundfix spec check <slug> --run-verification` before a Run starts: the
Daemon refuses a Task whose command already exits zero on the unchanged
tree." `rule.spec.docs-layout`: `clause.spec.prohibit-tests-that-read-specs` —
prohibited — "Do not make a test, fixture, or build step read a file under the
Spec root or assume that a Spec is still active; Specs archive and may be
deleted. A check whose subject is the Spec artifacts themselves is exempt."

**task_04, `typescript`.** `rule.typescript.behavior-tests`:
`clause.typescript.type-fixtures-from-the-schema` — mandatory — "Type a test
fixture that stands for a stored row from the schema's inferred row type, with
a complete base row and partial overrides, never from an untyped record."

### Source Baseline rows

Each row is a corpus entry `<!-- source-baseline-entry: <id> -->`, one line
`- <text>`, and the closing marker, placed after the named anchor entry with a
blank line between entries; a `manifest.json` row after the anchor's row with
`kind` `normative-clause`, the force, `carrier` `docs/agents/<file>`,
`structure` null and placeholder offsets and digest; and the identifier after
the anchor in `index.json` `entryIds`. The regeneration fills the rest.

| Task | Clause | After | Corpus file | Text |
| --- | --- | --- | --- | --- |
| 01 | `clause.core.regenerate-generated-files` | `clause.core.prohibit-verification-workarounds` | `agent-instructions.md` | MUST NOT hand-edit a generated file; change its source and run the generator that owns it. |
| 01 | `clause.core.flaky-tests-block` | `clause.core.assertion-reads-the-constant` | `agent-instructions.md` | MUST treat a flaky test as a blocking failure and fix the cause of its nondeterminism before completion. |
| 01 | `clause.core.lint-warnings-block` | `clause.core.flaky-tests-block` | `agent-instructions.md` | MUST fail Verification on a lint warning, and every warning Verification reports blocks completion. |
| 01 | `clause.core.prohibit-test-only-production-hooks` | `clause.core.lint-warnings-block` | `agent-instructions.md` | MUST NOT add a production hook, flag, branch, or exported symbol that exists only for tests. |
| 01 | `clause.core.prohibit-editing-vendored-skills` | `clause.core.use-github-pr-workflow` | `agent-instructions.md` | MUST NOT edit a skill installed from an upstream source; a correction is stated as a Repository-Specific Normative Rule. |
| 01 | `clause.core.ask-before-database-mutation` | `clause.core.prohibit-secret-exposure` | `agent-instructions.md` | MUST stop and ask for express authorization before a statement that changes data or schema in a database other than a disposable local one. |
| 01 | `clause.core.prove-the-mutation-predicate` | `clause.core.ask-before-database-mutation` | `agent-instructions.md` | MUST run a read with the same predicate and report its row count before an authorized database write, and run the write in an explicit transaction. |
| 01 | `clause.core.prohibit-disguised-database-mutation` | `clause.core.prove-the-mutation-predicate` | `agent-instructions.md` | MUST NOT route a database change through a migration, script, seed, or test to avoid that authorization. |
| 04 | `clause.spec.name-the-finding-a-task-answers` | `clause.spec.routing-05-task-graph` | `spec-routing.md` | MUST name, in one of a Task's requirements, the recorded finding and its date when the Task answers that finding. |
| 04 | `clause.spec.verification-fails-before-the-change` | `clause.spec.name-the-finding-a-task-answers` | `spec-routing.md` | MUST author each Task's Verification so every command fails on the tree before the Task's change. |
| 04 | `clause.spec.prohibit-tests-that-read-specs` | `clause.spec.specs-are-downstream-artifacts` | `docs-layout.md` | MUST NOT make a test, fixture, or build step read a file under the Spec root or assume that a Spec is still active. |
| 04 | `clause.typescript.type-fixtures-from-the-schema` | `clause.typescript.prohibit-incidental-test-oracles` | `typescript-bun.md` | MUST type a test fixture that stands for a stored row from the schema's inferred row type, never from an untyped record. |

The Bun clause's Source Baseline row stays: it is the record of what adopters
of that Source Baseline hold, and the classifier reads it to report the
replacement.

### Version changes

Raise each by one from its value on the Task's starting main; a new rule
starts at 1. task_01: `core`, `rule.core.root-cause-only`,
`rule.core.verification-integrity`, `rule.core.skill-dispatch`,
`rule.core.security-configuration`, `guide.agent-instructions`,
`guide.skill-dispatch`; `bun`, `guide.bun`. task_02: `go`, `guide.go` and its
three rules; `cli-surface`, `guide.cli-surface` and its two rules;
`tui-surface`, `guide.tui-surface` and its rule; the three modules'
`schemaVersion` becomes `setup-context-driven/module-v3`. task_03: `rust`,
`guide.rust`, `rule.rust.current-docs`, `rule.rust.cargo-gates`,
`template.guide.rust`; `schemaVersion` becomes v3. task_04: `spec-workflow`,
`rule.spec.routing`, `rule.spec.docs-layout`, `guide.spec-routing`,
`guide.spec-docs-layout`; `typescript`, `rule.typescript.behavior-tests`,
`guide.typescript-bun`.

### Existing tests that change

- `internal/baseline/clause_characterization_test.go` (every Task): the force
  record gains each new clause with its force; task_01 removes the Bun
  clause.
- `internal/baseline/preservation_test.go` (task_01, task_04): the maintained
  Source Baseline entry count rises by eight, then by four, from its value on
  the starting main.
- `internal/baseline/stack_rule_wording_test.go` (task_01): the conditional
  warnings sentence leaves the TypeScript and Bun guide's required list, and
  the Bun warnings row leaves the force table.
- `internal/baseline/plan_test.go` (task_02): `docs/agents/go.md`,
  `docs/agents/cli.md` and `docs/agents/tui.md` join the guides that grew past
  the frozen parity record. No other line changes.

### API Contracts

1. API Contract: rendered `docs/agents/agent-instructions.md` and
   `docs/agents/skill-dispatch.md` of every built-in profile — the eight core
   clauses with their force.
2. API Contract: rendered `docs/agents/go.md`, `cli.md` and `tui.md` of the Go
   CLI/TUI profile — every line labelled with force, the recorded-reason
   exception, the cross-platform build and the skill-ships-with-behavior
   clause.
3. API Contract: rendered `docs/agents/rust.md` of the Rust CLI profile — the
   scope sentence and the error rule with force.
4. API Contract: rendered `docs/agents/spec-routing.md`, `docs-layout.md` and
   `typescript-bun.md` — the four Spec-workflow and TypeScript clauses.
5. API Contract: `roundfix baseline update` for a Standard TypeScript Monorepo
   adopter of the Source Baseline — a ready plan whose retention records the
   Bun warnings clause `replaced` by `clause.core.lint-warnings-block`.

## Coverage Map

- Goal 1 → Exact texts task_02, task_03; Testing Approach 2, 3.
- Goal 2 → `clause.go.record-the-reason-for-a-module`,
  `clause.go.test-through-go-test`; Testing Approach 2.
- Goal 3 → `rule.rust.library-and-entry-point`; Testing Approach 3.
- Goal 4 → Candidate rules; Exact texts task_01, task_04; Testing
  Approach 1, 4.
- Goal 5 → the `replaces` declaration; Testing Approach 1.
- Story 1 → Exact texts task_02; Testing Approach 2.
- Story 2 → `clause.go.record-the-reason-for-a-module`; Testing Approach 2.
- Story 3 → Exact texts task_03; Testing Approach 3.
- Story 4 → Candidate rules; Testing Approach 1, 4.
- Story 5 → Testing Approach 1.
- Story 6 → the PRD's adopter notice; the QA gate's onioncry row.
- Core Feature 1 → task_02; API Contract 2.
- Core Feature 2 → task_03; API Contract 3.
- Core Feature 3 → task_01; API Contract 1, API Contract 5.
- Core Feature 4 → task_04; API Contract 4.
- Core Feature 5 → Candidate rules.
- Success Metric 1 → Testing Approach 2, 3.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → Testing Approach 3.
- Success Metric 4 → Testing Approach 1.
- Success Metric 5 → Testing Approach 1.
- Success Metric 6 → Testing Approach 4.
- Success Metric 7 → Testing Approach 1–4 and the existing characterization,
  duplicate-text and catalog tests.
- Success Metric 8 → Testing Approach 5.
- API Contract 1 → Exact texts task_01.
- API Contract 2 → Exact texts task_02.
- API Contract 3 → Exact texts task_03.
- API Contract 4 → Exact texts task_04.
- API Contract 5 → Exact texts task_01, `bun` and retention.

## Integration Points

- **Specs 0200 and 0207.** Both edit `go.json`, and 0200 also edits
  `core.json`, `rust.json`, `typescript.json` and the profiles; 0207 adds the
  Go guide's scope sentence and raises `template.guide.go`. This Spec touches
  neither the skills lists nor the Go template. Delivered after both, each
  Task raises versions from what they leave; a conflict confined to derived
  files is resolved by regeneration.
- **Spec 0208.** It renames the setup snapshots, edits the `core` and
  `typescript` skill lists, and adds eight external-triage clauses to the force
  record, the Source Baseline and its entry count. This Spec's counts and
  versions rise from the values 0208 leaves.
- **Adopters.** A Managed Refresh reports each changed guide. A Standard
  TypeScript Monorepo adopter's retention lists the Bun clause as `replaced`.
  Go and Rust adopters have no Source Baseline, so their update renders the
  new clauses without retention accounting.
- **Vendored skills.** Untouched. The rule-over-skill-default sentence and the
  new vendored-skill clause govern where they disagree.

## Testing Approach

Every test reads the embedded catalog or a plan built in a temporary
repository; none needs the network. Each negative case is its own test.

1. **Core (task_01),** new `internal/baseline/promoted_core_clauses_test.go`:
   `TestTheCoreGuidesStateThePromotedRules` builds a Standard TypeScript
   Monorepo plan and a Go CLI/TUI plan and requires the eight texts in the
   instruction and skill-dispatch postimages, and the absence of "When the
   repository's Verification treats warnings as errors" from the TypeScript
   and Bun postimage; `TestThePromotedCoreClausesCarryTheirForce` compares
   `clauseForceFindings` with a literal table and
   `TestAClauseForceCheckReportsAMissingOrChangedClause` feeds it a table with
   one absent identifier and one wrong force;
   `TestTheWarningsClauseIsReplacedByTheLintClause` runs
   `newClauseReplacementAdopter` and expects a ready plan whose delta and
   retention record the Bun clause `replaced` with the one target
   `clause.core.lint-warnings-block`;
   `TestAnUndeclaredWarningsReplacementIsUnaccounted` deletes `replaces` with
   `mutateCatalogClause` and expects `action_required` and `unaccounted`.
2. **Go, CLI, TUI (task_02),** new
   `internal/baseline/stack_force_go_cli_tui_test.go`:
   `TestEveryGoCLIAndTUIRuleCarriesForce` expects no finding from
   `ruleLevelGuidanceFindings` for the three modules and
   `TestARuleWithoutForceIsReported` feeds it a literal v2 rule;
   `TestTheGoCLIAndTUIGuidesStateTheirClauses` checks the `buildTestPlan`
   postimages for every new text, the absence of "Use stdlib `testing`." and
   of "Prefer the standard library; add a dependency only for a named job it
   cannot perform", and no unlabelled rule line;
   `TestTheGoGuideBansNoLibraryByName` expects no finding from
   `namedLibraryBanFindings` for Cobra, testify and Viper, and
   `TestALibraryBanInTheGoGuideIsReported` feeds it "Do not use Cobra.";
   `TestTheGoCLIAndTUIClausesCarryTheirForce` uses `clauseForceFindings`.
3. **Rust (task_03),** new `internal/baseline/stack_force_rust_test.go`:
   `TestEveryRustRuleCarriesForce`; `TestTheRustGuideStatesTheErrorPolicy`
   builds a Rust CLI plan with `newPlanRepository` and `planTestDecisions` and
   requires the scope sentence and the nine texts and the absence of "Preserve
   Cargo manifest and lockfile discipline.";
   `TestTheRustClausesCarryTheirForce`;
   `TestTheRustProfileRequiresTheErrorRule`.
4. **Spec workflow and TypeScript (task_04),** new
   `internal/baseline/promoted_spec_and_typescript_clauses_test.go`:
   `TestTheSpecAndTypeScriptGuidesStateThePromotedRules` (Standard TypeScript
   Monorepo and Go CLI/TUI postimages);
   `TestThePromotedSpecAndTypeScriptClausesCarryTheirForce`;
   `TestEveryPromotedClauseHasItsSourceBaselineRow` runs
   `sourceBaselineRowFindings` over the twelve promoted clauses and
   `TestAMissingSourceBaselineRowIsReported` feeds it a baseline without one.
5. **Regeneration and retention.** Each Task's Verification also runs
   `TestFormatterComposition`, `TestCatalogCompatibility`,
   `TestBaselinePlanCharacterization`, `TestBaselineCompatibilityCorpus`,
   `TestReadoptionCompatibilityMaintainedFixture`,
   `TestPlanDeterminismMatchesMaintainedManagedEntryFixture`,
   `TestStandardTypeScriptStructuralClauseRetention`,
   `TestBaselineClauseForceIsCharacterized` and its three negative tests,
   `TestNoTwoBaselineClausesShareText` and
   `TestShippedGuidanceCitesNoRepositoryRecord`, then the read-only Managed
   Refresh, which exits non-zero while a plan is pending.

In the rehearsal, each Task's edits turned the existing force record and, for
task_01 and task_04, the maintained Source Baseline count red until the test
edits above were made, and the replacement test failed as `unaccounted`
without `replaces`.

## Build Order

1. Core clauses and the lint replacement of the Bun clause, task_01 (depends
   on: none).
2. Go, CLI and TUI clauses with force, task_02 (depends on: 1).
3. Rust clauses with force and the error rule, task_03 (depends on: 2).
4. Spec-workflow and TypeScript clauses, task_04 (depends on: 3).
5. Terminal QA, task_05 (depends on: 1, 2, 3, 4).

The chain is serial: every Task edits the clause force record and rewrites the
catalog snapshots, plan goldens and Setup Manifest, and task_01 and task_04
both edit the Source Baseline files and the maintained count. task_02 must
follow task_01 so the Go test-hook ban is never absent. Prerequisites outside
the graph: deliver after Specs 0200 and 0207.

## Risks & Considerations

- **Hazard: the Source Baseline is shared history.** A new row claims that a
  Source Baseline adopter holds the clause, which is true after its next
  update. The rows are appended, never reworded, and the Bun row stays.
- **Hazard: generated files.** No golden, pin, snapshot or rendered guide is
  hand-edited; only the regeneration and the Managed Refresh write them.
- **The lint clause breaks adopters.** oraculum, vortex, fiscus and gss do
  not meet it until they change their Verification; the PRD's adopter notice
  names each.
- **Wording checks prove presence, not truth.** An Agent can still read a
  clause wrongly; the force labels remove the ambiguity the paragraphs had.
- **Counts that other Specs move.** Spec 0207 may add Source Baseline rows or
  clauses; each Task raises the count and the force record from the value on
  its starting main.

## Decisions

- **Promotion by recurrence.** See ADR-0203.
- **Force for every stack rule; a preference licenses its exception.** See
  ADR-0202.
- **Lint warnings in core, replacing the Bun clause.** One clause for every
  stack, declared as a replacement so retention holds.
- **The test-hook ban moves to core.** Three repositories in two languages
  wrote it.
- **No Go scope sentence here.** Spec 0207 adds it.
- **Tests in new files; existing test edits limited to the four listed.**
