---
spec: 0207-a-go-cli-and-a-typescript-monorepo-in-one-baseline
status: active
created: 2026-09-30
surfaces: [backend, cli, docs]
---

# A Go CLI and a TypeScript monorepo in one Baseline

The planned argus repository holds a Go command-line tool, a Hono backend and
a React frontend in one repository. No built-in Baseline Profile can express
it: a profile takes one Setup Snapshot, each snapshot mirrors one upstream
skill list, and no upstream list carries both the Go skills and the TypeScript
workspace skills, so a profile that selected both sets of modules would fail
the catalog check that its setup lists every skill its modules name. A
repository-owned profile cannot help either, because it may only narrow a
built-in profile. The Go guide also opens with no sentence naming what it
governs, so beside a Bun workspace its rules read as the whole repository's.

The same repository organizes its renderer by feature with shared stores and
components, while the frontend guide makes the systems layout mandatory. React
itself takes no position on folder layout. A repository that uses another
layout can today only contradict the Baseline.

This Spec ships a built-in composed profile for a Go CLI beside a TypeScript
monorepo, on a setup composed from two upstream lists by name, with one root
gate that runs both toolchains' Verification. It makes the frontend layout a
Frontend Layout Decision the repository records, with the systems layout as
the stated suggestion. See ADR-0204 and ADR-0205.

## Project Constraints

- Identifier strategy: applicable — new identifiers follow the catalog's
  existing forms: one built-in Baseline Profile identifier
  (`go-cli-typescript-monorepo`, the language and surface first, as in
  `go-cli-tui`), one Setup Snapshot identifier (`go-cli-typescript-bun`, named
  for the Go CLI and the TypeScript Bun workspace it serves), one decision identifier (`frontend.layout`, as
  `domain.layout`), one clause and rule identifier in the `frontend` module
  (`clause.frontend.follow-recorded-layout`, `rule.frontend.recorded-layout`),
  one Repository Capability (`capability.stack.go`), two Verification
  expectations (`verification.go`, `verification.workspace`), and six catalog
  diagnostic codes and one alignment divergence code in the existing dotted
  forms, which the TechSpec lists. No clause identifier is renamed or removed. No project-owned Internal
  Identifier is generated. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — embedded Baseline assets, local
  files and Git only; no credential is read and no Verification command opens
  a network connection. Under ADR-0063 repositories own their HTTP contract: the
  composed profile selects the backend module and carries the same
  repository-owned HTTP Contract Decision as the Standard TypeScript Monorepo
  Profile, and this Spec does not change it. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0204 (this Spec) governs the
  composed Setup Snapshot, the composed profile and the root gate's parts, and
  ADR-0205 (this Spec) governs the Frontend Layout Decision, the clauses that
  apply under a decision value and their retention. ADR-0191 makes a snapshot
  follow its upstream by name and the asset sync its only writer; a composed
  snapshot follows its components by name and the sync writes it. ADR-0190
  requires a stack rule to name the language or workspace it governs, which
  the Go guide's new sentence satisfies. ADR-0063 is the model the Frontend
  Layout Decision follows, and ADR-0061 keeps the composed profile's TypeScript
  stack as opinionated as the Standard TypeScript Monorepo Profile's. Under
  ADR-0058 Baseline upgrades fail closed on unaccounted rule removal, so a
  clause a recorded decision turns off is accounted for as a reasoned
  rejection. Under ADR-0060 Source Baselines are exhaustive and
  project-agnostic, so the systems clauses keep their identifiers, enforcement
  and bytes. Under ADR-0099 retention accounting is mechanical.
  ADR-0059 is why the composed profile declares no formatter instead of an
  unproven one. Under ADR-0067 custom Baseline Profiles are
  repository-owned, and a repository-owned profile that does not select the
  new decision keeps the systems clauses. Under ADR-0072 the Baseline Go cutover
  preserves Python contracts, so the composed snapshot, which the Python run
  never produced, is recorded as a designed delta. ADR-0186 requires adopter-neutral wording in every new
  sentence. ADR-0193 has this Spec name its prerequisites. ADR-0192 resolves a
  conflict confined to declared derived paths by regeneration. ADR-0081 and
  ADR-0149 keep the sanctioned regeneration outputs inside the grant,
  ADR-0073 and ADR-0103 require this repository's Managed Refresh to apply in
  one transaction and converge to `current`, and ADR-0130 has the audit judge
  only Governed Paths. The gate is bound by ADR-0080, ADR-0088, ADR-0091,
  ADR-0096, ADR-0104, ADR-0117, ADR-0155, ADR-0156, ADR-0166 and ADR-0167,
  ADR-0093 and ADR-0094 check the Spec by citation and artifact presence, and
  ADR-0178 and ADR-0179 decide how the audit reads the grant. Present in
  this tree and not applicable, because this Spec touches none of their
  subjects and edits no skill: ADR-0180 (the dated Recommended Profile),
  ADR-0181 (comparison with the recommendation, which cites ADR-0180),
  ADR-0182 (Settlement Checks, which cites ADR-0096), ADR-0183 (Claim
  Receipts, which cites ADR-0093), ADR-0184 (Surface Transcripts; the guide
  that would bind ADR-0183 and ADR-0184 is not in this tree), ADR-0187 (the
  Roundfix skill and the command reference), ADR-0189 (owned skill versions),
  ADR-0194 (QA row observation), ADR-0195 (rows observed on every pass, which
  cites ADR-0097), ADR-0196 and ADR-0197 (pre-PR review), ADR-0198 and
  ADR-0199 (token accounting and ceilings), ADR-0200 (an advisory
  judge for Spec authoring that never gates) and ADR-0201 (what that judge
  sends and spends, which cites ADR-0200), and ADR-0202 and ADR-0203 when
  present; all four belong to other Specs of this wave. ADR-0074
  (hybrid semantic ownership) cites ADR-0067, ADR-0097 (QA row carry-forward)
  cites ADR-0080, ADR-0168 (the related-ADR gap) and ADR-0176 (citation checks
  read authored text) cite ADR-0093, ADR-0173 (citations a History Relocation
  breaks) cites ADR-0073, and ADR-0177 is reached only through ADR-0173; none
  applies, because this Spec changes none of their behaviors. All hold.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30 for the Baseline source, its generated guides and this
  repository's Setup Manifest ("Autorizar os dois"), and the program decision
  that this Spec ships a built-in composed profile for a Go CLI with a
  TypeScript monorepo and records the frontend layout as a decision, recorded
  in [_authorization.md](_authorization.md); bounded files:
  `skills/baseline_skill_contract_test.go`,
  `internal/baseline/assets/decisions.json`,
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

## Goals

- A repository that holds a Go command-line tool beside a Bun workspace with a
  Hono backend and a React frontend can adopt one built-in Baseline Profile
  that passes every catalog check and renders each guide with no clause
  repeated.
- Every skill the composed profile's guides name comes from a setup that
  follows its upstream lists by name, so the next upstream refresh updates it
  with no hand edit.
- The composed profile names the Go and the Bun Verification its root gate
  runs, and profile alignment tells a repository whose root Make gate does not
  reach one of them.
- A repository records its frontend layout; the systems layout is stated as
  the suggestion only while none is recorded, and a recorded layout is never
  overwritten.
- Every existing adopter, including repository-owned profiles derived from the
  Standard TypeScript Monorepo Profile, keeps its current frontend guidance on
  its next update with no new question.

## User Stories

1. As the maintainer of a repository with a Go CLI and a TypeScript monorepo,
   I want to select one built-in Baseline Profile, so that my Agents get the
   Go, CLI, TypeScript, Bun, backend and frontend rules together without a
   rule that contradicts the other toolchain.
2. As a Roundfix maintainer, I want the composed setup to be derived from the
   upstream lists it combines, so that a snapshot refresh keeps it current in
   the same sync run.
3. As an adopter's Agent, I want the Go guide to say which code its rules
   govern, so that I do not apply the Go rules to the TypeScript workspace.
4. As the maintainer of that repository, I want Baseline to tell me when my
   root gate does not run the Go or the Bun Verification, so that a green gate
   means both toolchains passed.
5. As the maintainer of a React frontend organized by feature, I want to
   record that my repository defines its own frontend layout, so that the
   Baseline stops telling my Agents to organize code by domain system.
6. As an existing adopter, I want my next update to leave my frontend
   guidance's rules as they are and ask me nothing, so that fleet updates stay
   mechanical.

## Core Features

1. **A composed Setup Snapshot.** A Setup Snapshot may be composed from other
   snapshots named by identifier. It carries the union of their skills and
   activation bundles in component order and nothing else. The asset sync
   writes it after it refreshes its components, and the catalog refuses a
   composed snapshot that differs from that union, names an unknown or a
   composed component, or would merge two different entries under one skill
   name or one bundle identifier. The first composed snapshot combines the Go
   CLI list and the TypeScript/Bun list.
2. **A built-in composed profile.** A built-in Baseline Profile for a Go CLI
   beside a TypeScript monorepo selects the Go, CLI, TypeScript, Bun,
   monorepo, backend and frontend modules with the workflow modules every
   profile carries, takes the composed setup, requires the Standard TypeScript
   Monorepo stack and workspaces plus a Go module, and carries the same
   repository-owned HTTP Contract Decision and Frontend Layout Decision. It
   passes every catalog check, including the check that its setup lists every
   skill its guides name, and its rendered guides repeat no clause.
3. **The Go guide names what it governs.** The Go guide opens with one
   sentence that names the Go code its rules govern and sends code in another
   language to its own guide.
4. **A root gate that runs both toolchains.** The composed profile declares
   the Go Verification and the Bun workspace Verification as parts of the
   repository gate. When the selected gate is a Make target, profile alignment
   reports, as a non-blocking divergence, each part the target does not reach
   through its prerequisites or its recipes.
5. **The Frontend Layout Decision.** A repository records its frontend layout
   as `systems` or `repository-defined`. The decision is optional: while none
   is recorded, the frontend guide states that the suggested systems layout
   applies and renders the systems clauses as before; with `systems` recorded
   it states the recorded layout and renders the same clauses; with
   `repository-defined` recorded it renders one clause that binds the layout
   the repository's own rules state, and not the systems clauses. The
   interactive first adoption asks the question with `systems` preselected.
6. **Retention follows the recorded layout.** The systems clauses keep their
   identifiers, enforcement and text. An adopter that records
   `repository-defined` after it held the systems clauses updates with those
   clauses accounted for as reasoned rejections that name the recorded
   decision, not as unaccounted clauses that stop the update.

## User Experience

A maintainer runs `roundfix baseline` in a Go and TypeScript repository and
selects `go-cli-typescript-monorepo`. The interview asks the Standard
TypeScript Monorepo decisions, including the new frontend layout question with
`systems` preselected. The plan's alignment lists `verification.gate.part.missing`
for a part the root gate does not reach, with the part's command and a next
action. An existing Standard TypeScript Monorepo adopter runs
`roundfix baseline update --yes`: no decision is requested, its frontend guide
gains one sentence stating that the suggested systems layout applies while none
is recorded, and every clause stays. An adopter that later records
`repository-defined` through `roundfix baseline` sees the two systems clauses
replaced by the clause that binds its own layout, and the update's retention
lists both as reasoned rejections.

## Non-Goals / Out of Scope

- Clause force for the Go, Rust, CLI or TUI modules, stdlib-only as a Baseline
  rule, a Rust error policy, and scope sentences for the Rust, CLI and TUI
  guides. The open Backlog Entry "Stack modules need rules with force, and a
  repository with several stacks" keeps them; this Spec does not adopt it,
  because it implements only two of its questions.
- A profile that composes Go with a TUI surface, Rust with TypeScript, or any
  other combination. The composition mechanism allows them; no second
  composed profile ships.
- Rendering declared workspace paths into guides or checking any directory
  against a recorded layout.
- Moving the `app-renderer-systems` skill trigger from the TypeScript module to
  the systems layout; its trigger already names a feature system.
- A formatter fixture set for the composed profile; it declares no formatter.
- Following an included Makefile, a shell script, or a differently spelled
  command when checking the root gate's parts.
- Editing any skill, any upstream skill, `~/dev/skills`, or another repository.
- Any Makefile, lint, formatter, CI or `go.mod` change.

## Success Metrics

1. Success Metric: the embedded catalog loads with the composed profile; the
   profile's setup is the composed snapshot; the snapshot equals the union of
   its components in component order; and a composed snapshot that differs
   from that union, names an unknown or a composed component, or merges two
   different entries under one name is refused with its own diagnostic.
2. Success Metric: an asset sync from a source whose component list changed
   rewrites the composed snapshot in the same run, and a sync with no source
   change leaves it byte-identical.
3. Success Metric: a plan for the composed profile in a repository with a Go
   module and the Standard TypeScript Monorepo workspaces renders every guide
   of its modules, no rendered clause line occurs twice across them, and the
   Go guide carries its scope sentence.
4. Success Metric: with a root Make gate whose prerequisites or recipes reach
   both parts, alignment reports no part divergence; with a gate that reaches
   only one, it reports exactly one non-blocking `verification.gate.part.missing`
   naming the other.
5. Success Metric: a Standard TypeScript Monorepo plan with no recorded
   frontend layout renders the six frontend clauses it renders today and the
   suggestion sentence; with `systems` recorded it renders the same clauses and
   the recorded-layout sentence; with `repository-defined` recorded it renders
   the repository-layout clause and neither systems clause.
6. Success Metric: an update of a Setup Manifest with no recorded frontend
   layout reports no new decision; the structural-clause retention test still
   retains both systems clauses; a manifest that held the systems clauses and
   now records `repository-defined` plans with both accounted for as reasoned
   rejections and no unaccounted clause.
7. Success Metric: after the sanctioned regeneration and this repository's
   Managed Refresh, a second refresh reports no file change, and the Source
   Baseline corpus, the retention transitions and the parity fixture are
   byte-identical.

## Declared breaks

- **Break — frontend guide.** Every adopter's frontend guide gains one
  sentence on its next update: the suggestion sentence while no layout is
  recorded. No clause changes.
- **Break — first adoption.** The interactive first adoption of a profile
  with the frontend module asks one more question. An automation that passes
  decisions without it is not refused; the suggestion applies.
- **Break — retention.** An update that turns off a managed clause through a
  recorded decision value now proceeds with a reasoned rejection; the same
  change had no way to proceed before, because the clause could not be turned
  off.
- **Break — Go guide.** Every Go CLI/TUI adopter's Go guide gains its scope
  sentence on its next update. No rule changes.
- **Break — catalog.** The catalog lists one more built-in profile and one
  more Setup Snapshot, and the interactive profile list gains one entry.

## Prerequisites

- Spec 0200 is delivered first. It refreshes the `go-cli` and
  `typescript-bun` snapshots to one upstream commit, keeps the `go-cli`
  snapshot while giving the Go CLI/TUI profile `go-tui`, adds the catalog check
  that a profile's setup lists every skill its guides name, and changes the
  asset sync this Spec extends. The composed snapshot is built from its
  refreshed components, and the composed profile must pass its check.
- Spec 0195 is delivered first. Its sync rule keeps the Roundfix-owned entry
  in every setup, so the composed snapshot inherits it.
- Spec 0208 is delivered first. It follows the upstream setup renames of
  2026-10-01: `go-cli` retires into `go` and `typescript-bun` becomes
  `typescript`. The composed snapshot is built from the `go` and `typescript`
  snapshots it creates.

The evidence owner for each is the default branch. task_03's Verification
requires Spec 0200's check that every built-in profile's setup lists every
skill its guides name to pass by name for the composed profile; on a tree
without Spec 0200 that test does not exist and the Verification fails.

## Decisions

- **Compose setups, not profiles.** A composed Setup Snapshot keeps every
  reader of a profile's one setup unchanged. See ADR-0204.
- **Materialize the composition.** The composed snapshot is written in full by
  the asset sync, so raw readers of setup files keep working, and the catalog
  checks it against its components.
- **Name it as the others are named.** `go-cli-typescript-monorepo` puts the
  language and surface first, as `go-cli-tui` does, and ends with the
  TypeScript profile's shape.
- **No formatter claim.** The composed profile declares no formatter rather
  than claim fixtures no formatter has checked.
- **Advisory gate check.** The root gate's parts are checked textually in one
  Makefile, so a miss is a non-blocking divergence.
- **An optional decision with a stated suggestion.** Maintainer decision of
  2026-09-30: the layout is recorded, `systems` is the suggestion, a recorded
  value is never overwritten, and existing adopters keep their behavior. An
  optional decision satisfies all four with no migration. See ADR-0205.
- **Gate clauses, do not move them.** The systems clauses stay in the frontend
  module and apply under a decision value, so their bytes, identifiers and
  carrier stay put.

## Acceptance evidence

Each Core Feature requires positive and negative evidence in the Task Graph.
The outside-evidence row rests on sources this Spec did not produce:

- React's documentation on file structure: "React doesn't have opinions on how
  you put files into folders", followed by grouping by feature or route and by
  file type as common approaches
  (<https://legacy.reactjs.org/docs/faq-structure.html>).
- The argus repository at `~/dev/argus`, read only: its renderer is organized
  under feature directories with shared stores and components, which the
  systems layout's mandatory clauses would contradict.
- Adopter Setup Manifests and frontend guides in the Secondbrain mirrors, read
  on 2026-09-30: four Standard TypeScript Monorepo adopters (conexus, fiscus,
  fluxus, vortex) render both systems clauses, and vortex's Repository-Specific
  Normative Rules restate the systems layout by hand; one adopter (gss) uses a
  repository-owned profile, and two (gss, tax-poc) hold frontend guides from
  an older catalog without the systems clauses. None records a frontend
  layout.
- Published polyglot repositories whose root Make gate aggregates the Go and
  the JavaScript toolchains through prerequisite targets and recipes, such as
  `pedronauck/go-devstack`'s `make verify` and `namuh-eng/opensend`'s
  aggregate targets, which the gate check's prerequisite-and-recipe reading
  covers.

## Research basis

The Secondbrain was consulted through `wiki/index.md`, the adopter mirrors
under `projects/*/mirror/docs/agents/`, the upstream skills mirror
`projects/skills/mirror/setups/`, the inbox entry
`inbox/skills/2026-09-30-skills-vendorizadas-contradizem-o-baseline-ou-estao-quebradas.md`,
and `qmd query "perfil composto Go CLI TypeScript monorepo frontend layout systems argus"`,
whose top results were Roundfix's archived Specs on Baseline adoption and
manifest-backed updates and ADR-0061. The mirrors settled which adopters hold
the systems clauses and which use a repository-owned profile; the upstream
mirror settled that the `go-cli` and `typescript-bun` lists share every common
entry and that no upstream list combines them. The `roundfix` inbox namespace
held no pending entry. Exa located React's file-structure documentation and
the polyglot Makefile examples above. They support making the layout a
repository decision and reading the root gate's prerequisites and recipes.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the measured facts, the interfaces,
the exact texts and the build order.
