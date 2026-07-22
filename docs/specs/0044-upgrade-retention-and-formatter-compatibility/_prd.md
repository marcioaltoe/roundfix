---
spec: 0044-upgrade-retention-and-formatter-compatibility
status: active
created: 2026-07-22
surfaces: [cli, docs, infra]
---

# Context-Driven Baseline upgrade retention and formatter compatibility

A real 0.9.0 upgrade of a managed repository proved that a Context-Driven
Baseline version transition can replace previously managed hard rules
wholesale, compress supporting guides below their operational contract, drop
normative clauses during prose-to-rule migration, render duplicate skill
dispatch, and produce output the target repository's formatter immediately
rewrites — all while the current catalog passes its own asset contract. Two
2026-07-22 findings record the evidence: the upgrade preserved bytes outside
managed markers but had no boundary for the normative strength of what a
previous setup version had placed inside them. This Spec adds that boundary
while preserving the modular architecture Spec 0043 delivered.

## Goals

- A managed version transition cannot remove or weaken a previously managed
  mandatory clause without an explicit accounting the maintainer authorizes.
- Generated supporting guides carry complete operational contracts, proven by
  behavioral fixtures rather than category labels.
- Generated skill dispatch derives from the same catalog that installs the
  Repository Skill Set, with exactly one distinct entry per skill.
- Setup output composes with the target repository's selected formatter and
  Verification: apply, format, verify, audit, and reapply leave no delta.
- Baseline adoption states its own floor: dangling delegation from
  repository-authored documents is surfaced, and project-specific hard rules
  get an explicit preserved home.

## User Stories

1. As a repository maintainer upgrading a managed Context-Driven Baseline, I
   want every previously managed mandatory clause accounted for before I
   authorize the upgrade, so that a version transition cannot silently drop or
   weaken safety and delivery rules.
2. As a repository maintainer reviewing a Change Plan, I want to see which
   prior clauses are retained, moved, replaced, or rejected and why, so that
   my authorization covers the normative change and not only the file change.
3. As an Agent working in a configured repository, I want supporting guides
   that carry complete operational contracts — Spec routing, findings
   lifecycle, delegation authority, research protocol — so that I can act
   without consulting a deleted legacy corpus.
4. As an Agent, I want generated skill dispatch to route stack-specific work
   through the domain skills the repository actually installs, with one entry
   per skill, so that dispatch is actionable instead of noise.
5. As a repository maintainer whose Verification runs a formatter, I want
   generated output that stays clean under that formatter, so that I never
   have to choose between a clean setup audit and a passing format check.
6. As a repository author with project-specific hard rules, I want an
   explicit repository-owned home linked from the baseline and a warning when
   my documents delegate to categories the baseline does not cover, so that
   adoption never silently poses as a full replacement.

## Core Features

1. Portable mandatory rules carry stable clause-level identities and an
   enforcement strength such as mandatory, prohibited, or stop-and-ask. The
   Upgrade Retention Contract requires every managed version transition to map
   each prior mandatory clause to a current rule, a Repository-Owned
   Extension, or an explicit rejection with a recorded reason. Preview and
   apply block while any prior mandatory clause is unaccounted. See ADR-0058.
2. The Change Plan presents the retention accounting — retained, moved,
   replaced, and rejected clauses with their reasons — before authorization.
   An accepted clause must render with equivalent enforcement strength
   wherever it lands, including when it moves from root instructions to a
   supporting guide; weakened rendering is a blocking retention failure.
3. The supported legacy sample corpus gets a reviewed hard-rule ledger
   maintained as a migration fixture: genuinely portable clauses are promoted
   into setup-owned rules, stack-specific guidance renders only from confirmed
   repository capabilities or remains a Repository-Owned Extension, and every
   excluded sample clause carries an explicit recorded reason. Obsolete and
   conflicting clauses stay absent, and a fixture proves they stay absent.
4. The portable rule catalog gains the missing hard rules: warnings from the
   selected Verification block completion when the profile's toolchain treats
   warnings as errors; the repository-owned design contract must be read
   before UI work; dependent interfaces are inspected before tests are
   written; user-answerable decisions are asked, never guessed; the research
   protocol names the external web-research fallback and prohibits external
   research tools on local code; commit and delivery conventions are
   dispatched to their governing skills; and intentional
   verification-configuration changes pass an explicit authority gate.
5. Supporting guides restore the operational contracts that compression
   dropped: Spec entry-point routing (large initiative, feature, refactor or
   bug fix, trivial change) as a rule distinct from tracker hygiene; the
   findings lifecycle with its frontmatter template, per-state meanings,
   recorded-reason requirements, append-evidence rule, and update-timestamp
   rule; acceptance-evidence recording and slice discipline during general
   work; the Supervisor's hard prohibition on writing feature code or tests
   plus the named delegation channel; and the Secondbrain query contract with
   its concrete consult triggers.
6. Prior template prose is the coverage baseline for the current rule set.
   One catalog check fails when two supporting guides in a module render an
   identical rule list; another fails when a prior template's normative
   clause has neither a current rule nor a recorded retirement.
7. Technology and surface modules declare dispatch for the domain skills
   their profiles install, so the generated dispatch guide derives from the
   same catalog that pins the Repository Skill Set. A module that depends on
   another must not re-declare a dispatch entry for the same skill unless the
   trigger is genuinely distinct, and the generated guide renders one entry
   per skill.
8. During apply and audit, setup scans repository agent-instruction documents
   — root and nested — for delegation to setup-managed guidance. When a
   referenced category has no rule in the active catalog, setup emits a
   non-blocking finding stating that the generated baseline is a floor, not a
   replacement, naming each affected document.
9. Behind an explicit setup decision, setup scaffolds an unmarked
   Repository-Owned Extension document for project-specific hard rules, links
   it from the generated root instructions, and preserves it byte-for-byte
   afterward.
10. Generated managed Markdown is Formatter-Stable Output: the target
    repository's selected formatter leaves it unchanged. A fixture proves the
    full composition — confirmed apply, formatter check, selected
    Verification, fresh audit, and second apply all leave an empty diff and
    an empty Change Plan. Where formatter behavior varies by profile,
    formatter compatibility is declared as part of the generated-output
    contract. See ADR-0059.
11. An upgrade fixture derived from a real pre-0.9 managed corpus proves the
    retention boundary end-to-end: one unaccounted mandatory-clause deletion
    is rejected, accepted clauses render at equivalent strength,
    repository-authored extensions survive byte-for-byte, and every
    project-specific rule needing an adoption decision is reported.
12. The canonical setup skill, its distributed copy, generated guidance,
    asset contracts, and maintainer documentation describe the same
    retention, dispatch, delegation, and formatter behavior.

## User Experience

Before an upgrade, the maintainer sees one deterministic Change Plan that now
carries two layers: the file delta and the normative delta. The retention
section lists every prior mandatory clause as retained, moved, replaced, or
rejected, each with its reason. If any prior mandatory clause is unaccounted,
preview names it and apply refuses until the transition accounts for it.

After apply, running the repository's own formatter and Verification changes
nothing, and a fresh audit stays clean; a second apply reports an empty Change
Plan. When repository-authored documents delegate to categories the baseline
does not carry, audit reports an informational floor finding naming each
document; nothing blocks.

An Agent reading the generated corpus finds routing, lifecycle, authority,
research, and dispatch guidance complete enough to act on: the routing guide
answers "which pipeline entry point does this change take", the docs guide
carries the copyable findings template, and the dispatch guide lists each
installed skill once with a distinct trigger.

## Non-Goals / Out of Scope

- Restoring the legacy sample verbatim or using line count as a coverage
  target.
- Target-repository content restoration (such as the Fluxus
  repository-authored corpus); that work belongs to the affected repository.
- The retired knowledge-workspace sparse-checkout and separate documentation
  commit flow.
- Obsolete runtime and model defaults, and template stack or product
  assumptions carried by the sample.
- Conventional REST guidance that conflicts with a repository's explicit API
  contract.
- Generating project-specific architecture, feature layout, or technology
  dispatch the repository has not confirmed.
- Generating or rewriting nested package-level agent instructions; the
  delegation scan reads them only to report dangling references.
- Running the target repository's formatter as a setup mutation; setup proves
  compatibility, it does not format.
- Expanding Doctor Skill Readiness from Spec 0036 or the mandatory Baseline
  ADR set from Spec 0040.

## Success Metrics

- The real pre-0.9 upgrade fixture blocks preview and apply on one
  unaccounted mandatory-clause deletion; after every prior clause is mapped,
  the same transition passes and its Change Plan lists the retained, moved,
  replaced, and rejected clauses.
- One hundred percent of accepted prior mandatory clauses render with
  equivalent enforcement strength; a weakened-rendering fixture fails the
  retention check.
- Zero pairs of supporting guides render identical rule lists, and removing a
  routing, lifecycle, or authority clause from the catalog without a recorded
  retirement fails the migration check.
- The generated dispatch guide renders exactly one entry per skill across
  every bundled profile, and a near-duplicate dispatch declaration fails
  asset validation.
- The formatter fixture proves apply, formatter check, selected Verification,
  fresh audit, and second apply all leave an empty diff and an empty Change
  Plan.
- Every excluded sample clause carries a recorded reason, and a fixture
  proves obsolete clauses remain absent.
- A dangling-delegation fixture produces the non-blocking floor finding
  naming each affected document, and a fully covered repository produces
  none.
- The full repository verification gate passes, and canonical and
  distributed setup assets match.

## Decisions

- One Spec owns the setup-owned corrections from both 2026-07-22 findings
  because retention, migration fidelity, dispatch, and formatter stability
  are boundaries of the same upgrade workflow; target-repository content
  restoration routes to the affected repository.
- The corrective work preserves the modular architecture from Spec 0043 —
  versioned rule catalog, compact root pointers, generated dispatch,
  repository-owned extension boundaries — and adds the missing compatibility
  boundary; it does not restore the sample verbatim.
- Baseline upgrades fail closed on unaccounted managed-rule removal. See
  ADR-0058.
- Generated output is formatter-stable in the target repository. See
  ADR-0059.
- Dispatch deduplication is enforced at the catalog level rather than merged
  at render time, so preview, audit, and apply keep resolving one declarative
  source, extending ADR-0047.
- The delegation signal is informational, never blocking: the baseline states
  its floor without gatekeeping repository-authored documents.
- The extension scaffold is decision-gated and unmarked: setup creates it
  once when authorized and never manages its content afterward.
- Compact root instructions stay; operational detail returns to the
  supporting guides, and compact wording is acceptable only while every
  decision, prohibition, and escalation point remains testable.
- The Upgrade Retention Contract extends the declarative asset model of
  ADR-0046 and ADR-0047: clause identities and transition mappings are
  catalog data, not imperative migration code.

## Open Questions

None.
