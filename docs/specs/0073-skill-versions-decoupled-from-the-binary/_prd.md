---
spec: 0073-skill-versions-decoupled-from-the-binary
status: active
created: 2026-08-03
surfaces: [backend, cli, docs]
---

# Skill versions decoupled from the binary

Roundfix pins skill *content* rather than skill *compatibility*. Each setup
snapshot carries a `treeDigest` per skill, the catalog digest is computed over
those snapshots, and both characterization corpora embed observed digests
inside recorded diagnostics. Any skill edit therefore changes what the binary
claims, and a Roundfix release asserts a fact about content it does not
version.

Three failures in one session, all from that coupling:

1. Editing one owned skill to record an authoring rule broke `make verify`
   twice — once through each corpus — and left the repository not green, which
   failed the next Task on entry. `make baseline-digests` reported no changes,
   because neither corpus is a member of its steps.
2. Three externally published Go skills could not be locked. Locking requires a
   setup snapshot entry whose `treeDigest` is computed at a pinned source ref,
   and the maintainer's checkout does not carry it. The skills work — they are
   on disk and the dispatch routes to them — but their provenance is
   unrecorded.
3. `baseline assets sync` refuses outright: the upstream `setups/go-cli.txt`
   dropped `bubbletea` and `tui-design`, which this repository's `go-cli-tui`
   profile requires, so a sync would invalidate the profile.

A byte-exact pin answers "is this the same content?" when the question Roundfix
needs answered is "is this content new enough to work with me?". Those are
different questions, and only the second survives the skill evolving on its own
schedule.

## Project Constraints

- Identifier strategy: applicable — this Spec introduces a skill version as a
  compared value. It is not a project-owned Internal Identifier: the version is
  declared by the skill and read by Roundfix, so ownership stays with the
  skill. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — the governing clause prohibits reading,
  printing, committing, or generating secrets and forbids inventing
  authentication or transport policy. Resolving a skill's declared version
  reads local installed files and introduces no credential and no new endpoint.
  Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0081 keeps sanctioned digest
  regeneration a fallout of the authorized edit, which whatever replaces the
  content pin must preserve; ADR-0085 keeps a regeneration run ungated by the
  pins it rewrites while every other load stays strict. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: on
  2026-08-02 the maintainer authorized tooling adjustment naming the `Makefile`
  and the owned skills, recorded at
  `docs/workflow/authorizations/2026-08-02-queued-spec-tooling.md`; bounded
  files: `Makefile` and the owned skill pair. Setup snapshot and profile assets
  under `internal/baseline/assets/**` are product assets rather than protected
  tooling. Source: `docs/agents/agent-instructions.md`.

## The prerequisite this Spec depends on

**No skill declares a version today.** Skill frontmatter carries `name` and
`description` and nothing else, so there is currently no value to compare a
minimum against. The compatibility contract requires skills to declare a
version, which is a change in the skills source rather than in Roundfix.

This Spec is written so Roundfix is ready for that and degrades honestly until
it arrives: a skill with no declared version is reported as unversioned, which
is distinct from satisfying the minimum and distinct from falling below it.
Roundfix never invents a version for a skill that does not state one.

## Goals

- A skill's version and Roundfix's version move independently.
- Roundfix states the minimum version of each skill it needs, and refuses to
  operate on anything below it with a message naming the skill, the minimum,
  and what was found.
- A skill above the minimum is accepted without Roundfix changing, so a skill
  can ship on its own schedule.
- Nothing in the binary has to change when skill content changes.

## Core Features

1. Roundfix declares a minimum version per required skill. That declaration —
   not a content digest — is what a profile carries.
2. A skill's installed version is read from what the skill itself declares.
   Roundfix compares; it never derives or assumes a version.
3. Readiness is a comparison, not an equality: at or above the minimum
   satisfies, and a newer skill needs no change in Roundfix.
4. Below the minimum is a blocking failure that names the skill, the required
   minimum, the version found, and how to upgrade.
5. Three states stay distinct and are never collapsed: satisfies the minimum,
   below the minimum, and unversioned or unresolvable. An unreachable source is
   never reported as a missing skill.
6. The Doctor Command reports skill readiness under this contract, and every
   command that gates on skills uses the same comparison, so two surfaces
   cannot disagree.
7. Content digests stop gating compatibility. Where a digest still protects
   something Roundfix genuinely owns — the guides it generates — it stays.
8. The characterization corpora stop embedding volatile skill digests in
   recorded diagnostics, so a skill edit cannot invalidate them.
9. A repository whose Baseline was applied before this Spec keeps validating,
   and archived Spec artifacts stay byte-identical.

## Non-Goals / Out of Scope

- Changing what any skill contains, or which skills a profile needs.
- Adding version declarations to the skills themselves; that is upstream work
  this Spec depends on but does not perform.
- Removing digest verification from artifacts Roundfix does own.
- Vendoring skills into the binary.
- The derived-artifact regeneration boundary, owned by Spec 0067. This Spec
  removes a cause; that one fixes the ownership map.

## Success Metrics

- Editing an owned skill leaves `make verify` green with no regeneration step.
- A skill one minor above its declared minimum satisfies readiness with no
  change to Roundfix.
- A skill below the minimum blocks with a message naming skill, minimum, found
  version, and upgrade path.
- A skill declaring no version is reported unversioned, distinctly from both
  satisfying and failing.
- The three Go skills added on 2026-08-02 carry recorded provenance without a
  setup-snapshot entry.
- A Baseline applied before this Spec still validates unchanged.

## Decisions

- Compatibility is a floor, not an equality. Roundfix asks whether a skill is
  new enough to work with, which is the only question that survives the skill
  evolving on its own schedule.
- The skill owns its version; Roundfix owns the minimum. Neither derives the
  other, and Roundfix never invents a version for a skill that does not declare
  one.
- Unversioned is its own state. Collapsing it into either pass or fail would
  make the contract unreadable during the transition.
- This Spec evolves the skill contract and never regresses it: a Baseline
  applied today still validates, no archived artifact is rewritten, and every
  digest protecting something Roundfix owns stays.

## Open Questions

None.
