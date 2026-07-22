---
spec: 0043-context-driven-baseline-coverage-and-skill-restoration
status: active
created: 2026-07-21
surfaces: [cli, docs, infra]
---

# Context-Driven Baseline coverage and Repository Skill Set restoration

The 2026-07-21 Fluxus setup investigation proved that a clean setup audit can
still leave a repository without portable instruction coverage, retain a
generated reference to a missing guide, omit a removal from the pre-apply
change plan, and offer no reproducible remediation for required-skill drift.
The Context-Driven Baseline must remain compact while making its coverage,
managed references, transition plan, and external skill restoration complete
and auditable across every supported profile.

## Goals

- Make every supported profile declare and prove the portable instruction
  coverage its generated repository guidance provides.
- Keep root agent instructions compact by routing detailed requirements to
  setup-owned supporting guides and preserving repository-authored content.
- Make a successful audit prove that every setup-owned reference resolves for
  the selected Decision Plan.
- Make the pre-apply change plan account for every file-tree mutation before
  the user authorizes apply, including conditional removals.
- Make every external required-skill drift finding reproducibly actionable
  from portable, immutable provenance.

## User Stories

1. As a repository maintainer adopting a Context-Driven profile, I want the
   generated guidance to cover every portable workflow invariant the profile
   promises, so that replacing legacy instructions does not silently remove
   safety or delivery rules.
2. As an Agent working in a configured repository, I want generated skill
   dispatch to match the selected modules' required skills, so that I know
   which skill governs each kind of work before I act.
3. As a repository maintainer auditing a selected Decision Plan, I want every
   setup-owned reference checked against the artifacts that plan generates, so
   that a successful audit never accepts a pointer to missing guidance.
4. As a repository maintainer approving setup changes, I want every creation,
   refresh, and removal named before apply, so that authorization covers the
   complete observed change.
5. As a repository maintainer facing required-skill drift, I want an immutable
   and portable restoration action, so that I can reproduce the expected
   directory without guessing a branch, commit, or machine-local source.
6. As a repository author, I want setup to state which guidance it owns and
   which project-specific rules remain mine, so that baseline adoption never
   claims to replace repository architecture or policy that it does not
   manage.

## Core Features

1. The portable asset catalog defines a semantic coverage contract made of
   stable rule identifiers. Every supported profile must cover the applicable
   categories for universal safety, selected Verification enforcement,
   verification-configuration integrity, skill dispatch, language, research
   sources, dependency changes, Git and delivery, security and configuration,
   and enabled application surfaces.
2. Generated root instructions remain mandatory pointers rather than a full
   manual. Setup-owned supporting guides carry the detailed portable rules,
   identify the selected repository Verification, derive skill dispatch from
   module requirements, and direct frontend-enabled profiles to the
   repository-owned design contract without generating project-specific
   architecture.
3. Generated guidance explicitly distinguishes setup-owned baseline rules from
   repository-authored extensions. Audit proves managed content only and never
   treats the portable baseline as evidence that project-specific rules are
   complete.
4. Audit resolves generated references against the exact artifact set selected
   by the Decision Plan. A reference to an excluded, absent, or
   repository-external target is a blocking finding, including references
   written as agent-facing path pointers rather than Markdown links.
5. The machine-readable pre-apply plan names every creation, refresh, removal,
   rename, and reference edit with its path, managed identity, state,
   condition when applicable, and reason. The applied file-tree delta must
   equal the authorized plan.
6. Setup preserves every unmarked repository-authored file unless an explicit
   adoption or removal decision includes it in the presented plan. Legacy
   markers and profile transitions do not grant implicit removal authority.
7. Every external required skill in a bundled setup snapshot carries a
   portable source identity, immutable revision, source-relative location, and
   expected complete-directory digest. Drift findings expose the same values
   in machine-readable remediation data.
8. The supported restoration path reproduces the complete expected skill
   directory, previews files that will disappear, verifies the resulting
   digest, and leaves portable provenance in the repository skill lock. If no
   supported mechanism can honor the immutable revision, setup blocks with an
   explicit upstream dependency instead of suggesting a generic refresh.
9. Profile, decision-transition, reference, and restoration fixtures prove the
   same contracts across every bundled profile and every decision combination
   that includes or excludes a referenced artifact.
10. The canonical setup skill, its distributed copy, generated guidance, asset
    contracts, and maintainer documentation describe the same coverage,
    preview, audit, and restoration behavior.

## User Experience

Audit remains read-only and local. A clean result means the selected profile
satisfies its semantic coverage contract, all generated references resolve,
and every required installed skill matches its expected digest. A blocking
result names the stable finding code, affected managed identity, path, reason,
and next action.

Before apply, the maintainer sees one deterministic plan containing every
managed file or block that will be created, refreshed, removed, renamed, or
edited. Apply requires confirmation for that complete plan, performs no
unlisted mutation, preserves unmarked content without explicit authority, and
leaves a second apply with no changes.

Required-skill drift reports an exact restoration action backed by immutable
portable provenance. Restoration remains an explicit maintainer action rather
than a side effect of audit or apply; after it runs, the repository skill lock
and complete skill directory pass the same audit proof on another machine.

## Non-Goals / Out of Scope

- Copying the legacy TypeScript/Bun sample verbatim or using line count as a
  coverage target.
- Generating project-specific architecture, authentication, database,
  transport, runtime-selection, or Secondbrain-write policy.
- Restoring retired knowledge-workspace commit flows or overriding a
  repository's explicit API and security contracts.
- Rewriting repository-authored instruction bodies or adopting unmarked files
  without explicit authorization.
- Generating or validating nested package-level agent instructions.
- Turning setup into a general-purpose Agent Skill package manager or removing
  extra installed skills.
- Restoring required skills automatically during audit or documentation apply.
- Modifying upstream-managed skill content inside Roundfix; this feature owns
  only portable provenance, restoration behavior, and setup-generated
  guidance.
- Expanding Doctor Skill Readiness from Spec 0036 or the mandatory Baseline ADR
  set from Spec 0040.

## Success Metrics

- One hundred percent of bundled profiles satisfy every mandatory and
  applicable semantic coverage category; removing one required rule or one
  required-skill dispatch mapping makes asset validation fail.
- Every supported Decision Plan combination produces zero unresolved
  setup-owned references after apply; removing one referenced target makes
  audit fail instead of reporting a clean result.
- For every marked and unmarked transition fixture, the pre-apply plan equals
  the complete observed file-tree delta, with zero unannounced removals and
  zero changes to unmarked content without explicit authorization.
- One hundred percent of external required skills carry immutable portable
  provenance; a drift fixture restores the exact complete-directory digest and
  records zero machine-local absolute source paths.
- A second apply produces zero changes, the final audit exits successfully,
  canonical and distributed setup assets match, and the full repository
  verification gate passes.

## Decisions

- One Spec owns all four reported gaps because coverage, managed references,
  preview authority, and skill restoration are correctness boundaries of the
  same setup workflow.
- Specs 0036 and 0040 retain their existing scopes: Doctor proves the installed
  Repository Skill Set, while setup owns generated repository documents and
  the Context-Driven Baseline.
- Preserve compact root instructions and move portable detail into setup-owned
  guides; semantic rule coverage replaces sample size as the completeness
  contract.
- Extend the declarative asset and Decision Plan model accepted by ADR-0046 and
  ADR-0047 rather than introducing imperative profile-specific behavior.
- Require plan/apply file-tree parity and explicit authority for every removal.
- Treat immutable, portable restoration provenance as mandatory for external
  required skills; a generic refresh instruction is not actionable
  remediation.
- No new ADR is required because the existing declarative ownership and
  Decision Plan decisions already govern these additions.

## Open Questions

None.
