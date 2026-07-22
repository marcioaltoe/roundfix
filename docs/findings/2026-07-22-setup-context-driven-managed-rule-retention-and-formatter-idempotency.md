---
status: done
created_at: 2026-07-22
updated_at: 2026-07-22
---

# Setup Context-Driven — managed hard-rule retention and formatter idempotency

This is a post-implementation finding for
[Spec 0043 — Context-Driven Baseline coverage and Repository Skill Set restoration](../specs/_archived/0043-context-driven-baseline-coverage-and-skill-restoration/_prd.md).
The Spec correctly introduced a modular, versioned rule catalog, compact root
pointers, generated skill dispatch, and repository-owned extension boundaries.
A real 0.9.0 upgrade in `/Users/marcio/dev/fluxus` nevertheless showed that the
new catalog proves its own internal completeness without proving that an
upgrade preserves the normative strength of previously managed instructions.
The concurrently authored
[v2 refresh content-regression finding](2026-07-22-setup-context-driven-v2-refresh-content-regression.md)
records the corpus-level rule-mapping and rendering defects in more detail;
this finding owns the cross-version retention contract and formatter
compatibility boundaries that prevent the same class of regression on later
upgrades.

The comparison used:

- `git diff HEAD -- AGENTS.md docs/agents` in Fluxus;
- `/Users/marcio/dev/skills/samples/AGENTS.typescript-bun.md` as evidence of
  useful legacy rules, not as a template to copy verbatim;
- the generated `AGENTS.md` and `docs/agents/` corpus after the confirmed 0.9.0
  apply;
- the 0.9.0 profile, module, asset-validation, and apply tests in Roundfix.

## 1. An upgrade can replace managed hard rules without accounting for them

- **Symptom / evidence**: the Fluxus diff changes every root block from its
  earlier managed version to a compact versioned pointer. That architecture is
  desirable, but the previous setup-owned bodies are replaced wholesale. The
  promise to preserve repository-authored bytes outside markers does not cover
  hard rules that a previous setup version placed inside managed markers.
- **Symptom / evidence**: the TypeScript/Bun profile lists every rule currently
  selected by its modules in `requiredRules`, and
  `_validate_profile_rule_contracts` proves that each listed rule has an owner,
  carrier, and render binding. If a rule is removed from both a module and the
  profile, however, the catalog remains self-consistent. The contract has no
  previous-version rule ledger, migration map, or required justification for a
  deleted or weakened rule.
- **Symptom / evidence**: the coverage taxonomy identifies broad categories
  such as verification, language, dependencies, Git delivery, and security.
  A rule can still satisfy one of those category labels with guidance weaker
  than the prior hard rule. Versioning the new rule does not state which old
  clauses it preserves, replaces, or intentionally retires.
- **Root cause**: semantic coverage is validated only inside the selected
  current catalog. Managed-block migration is file- and identity-aware, but not
  normative-rule-aware. The upgrade therefore has no fail-closed boundary for
  unaccounted hard-rule removal.
- **Action / suggestion**: add a versioned upgrade-retention contract. Give
  portable hard rules stable, clause-level identities and an enforcement mode
  such as `must`, `must-not`, or `stop-and-ask`. For every managed version
  transition, map each prior mandatory clause to a current rule, a
  repository-owned extension, or an explicit rejection with a reason. Block
  preview/apply when a prior mandatory clause has no mapping, and show the
  retained, moved, replaced, and rejected mappings in the Change Plan.

## 2. The TypeScript/Bun result loses or weakens useful legacy contracts

The generated baseline preserves useful generic rules: the selected
Verification blocks completion, root-cause fixes remain mandatory, fresh
evidence is required, local search is separated from external documentation,
the repository package manager and lockfile remain authoritative, delivery
requires authority, and secrets remain outside source control. It does not yet
preserve the following normative details from the sample:

| Legacy contract | Generated result | Assessment |
| --- | --- | --- |
| Read `DESIGN.md` before any UI implementation | `docs/agents/frontend.md` only says to follow the design contract | Weakened: the precondition is gone |
| Zero lint warnings; every required check blocks completion | Verification blocks failures, but the TypeScript/Bun warning contract is absent | Partial |
| Never edit verification configuration to silence failures; stop and ask before a justified config change | `agent-instructions.md` prohibits making failures disappear but permits intentional contract changes without an explicit authority gate | Weakened |
| Inspect dependent APIs before writing tests | No equivalent rule | Lost |
| Use `coding-guidelines`, `clean-code`, and `solid` together for production code | Generated dispatch includes only `coding-guidelines` from that trio | Lost |
| Use `no-workarounds` and `systematic-debugging` for fixes, plus `testing-boss` for tests | The generic dispatch retains these triggers | Preserved |
| Use local search for local code and current external sources for libraries, with explicit Context7/Exa routing | Generic source authority and Context7 remain; Exa and its research trigger are absent | Partial |
| Verify a dependency and its current version, then run `bun add` from the owning workspace | Bun commands are required, but package verification, current version, `bun add`, and workspace ownership are absent | Weakened |
| Ask the user instead of guessing cheap decisions | No equivalent rule | Lost |
| Write code, comments, migrations, API contracts, and structured data in English while retaining canonical domain language | Only generated guidance, identifiers, headings, and examples are explicitly English | Weakened |
| Require explicit permission before destructive Git commands | The guide says only to avoid destructive operations | Weakened |
| Conventional Commit and PR-title/body contracts | No equivalent rule | Lost |
| Keep `.env` secret, mirror safe keys in `.env.example`, and audit new dependencies | Secret safety remains; example-key parity and dependency audit are absent | Partial |

The generated `docs/agents/skill-dispatch.md` also exposes only the skills
declared by the current modules. It omits useful TypeScript/Bun stack triggers
from the sample, including `clean-code`, `solid`, `bun`, `turborepo`, `vite`,
React and TanStack skills, feature-system and UI skills, Hono, Drizzle, Zod,
integration-contract testing, observability, and security review. Some of those
are portable for a selected technology surface; others depend on the actual
repository stack. The setup must classify them instead of either generating
all of them blindly or dropping them silently.

- **Root cause**: the current catalog expresses only the newly selected rule
  and skill set. It has no adoption ledger connecting the sample or the prior
  managed corpus to portable modules and repository-owned extensions.
- **Action / suggestion**: maintain a reviewed TypeScript/Bun hard-rule ledger
  as a migration fixture. Promote genuinely portable clauses into setup-owned
  modules. Render stack-specific dispatch only from confirmed repository
  capabilities, or preserve it as a typed repository-authored extension
  outside setup markers. Require every excluded sample clause to carry an
  explicit reason.

## 3. Existing managed guides were compressed past their operational contract

- **Symptom / evidence**: `docs/agents/spec-routing.md` now contains Task
  ownership and status rules that duplicate `issue-tracker.md`; it no longer
  routes initiatives, features, bug fixes/refactors, and trivial changes to
  their correct Spec entry points. The root still promises to “Route and
  execute local Specs” through that guide.
- **Symptom / evidence**: `docs/agents/docs-layout.md` retains directory jobs
  and lifecycle state names, but removes the copyable finding frontmatter,
  per-state meanings, and the instruction to append evidence and routing links
  instead of rewriting observations.
- **Symptom / evidence**: the Secondbrain guide retains the broad read-only,
  query-first, citation, secret, and Hermes boundaries, but removes much of the
  actionable query order and path-specific safety detail. The autonomous guide
  says implementation is delegated, but no longer states the previous hard
  prohibition that the Supervisor never writes feature code or tests.
- **Root cause**: compactness was applied to supporting rule bodies rather than
  only to the root index. There is no regression fixture for the operational
  behavior carried by a previous managed guide.
- **Action / suggestion**: keep `AGENTS.md` compact, but restore operational
  detail in the supporting guides. Add guide-specific behavioral fixtures for
  Spec entry routing, finding lifecycle/evidence semantics, Secondbrain safety,
  and Supervisor authority. Compact wording is acceptable only when all
  decisions, prohibitions, and escalation points remain testable.

## 4. Setup output is not formatter-idempotent in the target repository

- **Symptom / evidence**: immediately after the confirmed 0.9.0 apply and clean
  setup audit in Fluxus, `rtk bunx oxfmt --check AGENTS.md` exited `1` and named
  `AGENTS.md`. Oxfmt inserts blank lines around the managed-marker content.
- **Symptom / evidence**: the Fluxus gate runs the formatter in mutating mode.
  In the setup session, formatting the generated root changed the managed
  bytes, after which setup audit planned managed-reference edits. Re-applying
  restored audit cleanliness but made Oxfmt dirty again. File-level apply
  idempotency therefore does not compose with the repository Verification.
- **Root cause**: generated Markdown is canonical only according to the setup
  renderer, not according to the formatter selected by the target repository.
  The test suite has apply/reapply and byte-idempotency coverage but no
  apply → repository formatter → audit compatibility check.
- **Action / suggestion**: make the Markdown templates formatter-stable and add
  an Oxfmt TypeScript/Bun fixture that proves: confirmed apply, formatter check,
  selected Verification, fresh setup audit, and second apply all leave an
  empty diff and empty Change Plan. If formatter behavior varies by profile,
  make formatter compatibility an explicit generated-output contract rather
  than requiring repositories to choose between setup audit and formatting.

## Required correction boundary

The corrective work should preserve the modular 0.9 architecture and add the
missing compatibility boundary; it should not restore the sample verbatim.
Explicitly reject or keep repository-owned:

- the retired `.knowledge` sparse-checkout and separate documentation commit
  flow;
- the obsolete Fable and model/runtime defaults;
- template PostgreSQL, Better Auth, and product-name assumptions;
- conventional REST guidance, which conflicts with Fluxus's package-level
  POST-only action contract and current internal unauthenticated backend;
- project-specific feature layout or technology dispatch not confirmed by the
  repository.

A follow-up Spec should include at least these proof cases:

1. upgrade a real pre-0.9 managed Fluxus fixture and reject one unaccounted
   mandatory-rule deletion;
2. prove each accepted hard rule is rendered with equivalent enforcement
   strength, even when moved from root to a supporting guide;
3. preserve repository-authored extensions byte-for-byte and report every
   project-specific rule that needs an adoption decision;
4. verify restored Spec routing and findings lifecycle behavior;
5. run apply → Oxfmt → selected Verification → audit → reapply with no delta;
6. prove excluded obsolete or conflicting sample clauses stay absent with a
   recorded reason.

## Routing

Keep this finding `pending` until a follow-up Spec explicitly owns managed-rule
retention and formatter compatibility. Spec 0043 is already archived with a
passing QA report; this real-repository upgrade is new post-QA evidence and
must not be treated as covered merely because the current catalog passes its
own v2 asset contract. Route the companion content-regression finding to the
same follow-up only where its rule-level corrections share this contract;
Fluxus-owned content restoration remains a separate repository concern.

Routed 2026-07-22:
[Spec 0044 — Context-Driven Baseline upgrade retention and formatter compatibility](../specs/0044-upgrade-retention-and-formatter-compatibility/_prd.md)
owns the Upgrade Retention Contract, the legacy hard-rule ledger, the
supporting-guide restoration, and Formatter-Stable Output, recorded in
ADR-0058 and ADR-0059.
