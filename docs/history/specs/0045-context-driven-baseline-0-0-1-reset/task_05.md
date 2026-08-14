---
task: task_05
spec: 0045-context-driven-baseline-0-0-1-reset
status: completed
type: backend
complexity: high
---

# Task 05: Deliver the Standard TypeScript Monorepo Profile

## Overview

Deliver the opinionated TypeScript monorepo baseline as a complete declarative
profile. The profile must produce exact topology, stack, architecture,
capability, decision, and Verification contracts rather than a loose menu of
suggestions.

## Requirements

1. MUST add the `standard-typescript-monorepo` profile with the exact required
   stack named by the PRD: TypeScript, Bun, Turborepo, Vite, React, Hono,
   Drizzle, Zod, Tailwind, shadcn, TanStack Query, TanStack Router, Better Auth,
   PostgreSQL, LogTape, Oxlint, Oxfmt, and Vitest.
2. MUST define exactly `packages/frontend` and `packages/backend` as required
   workspaces.
3. MUST model Inngest and Docker as optional modules whose absence does not
   block the baseline.
4. MUST require frontend organization by systems and backend organization by
   domain, application, and infrastructure while rejecting generic
   `modules`/`services` guidance as the normative structure.
5. MUST require a typed HTTP Contract Decision of `REST` or `Post-only`, with
   explicit repository-owned exceptions and no implicit default when the
   repository has not decided.
6. MUST bind the universal Repository Capability requirements and exact Skill
   Activation bundles to the profile.
7. MUST declare reproducible formatting, linting, testing, and build commands
   as persisted Verification entries.
8. MUST render only strict 0.0.1 owned schemas, markers, and profile versions.

## Subtasks

- [x] Define the exact profile topology, stack, and optional modules.
- [x] Add frontend systems and backend layered architecture clauses.
- [x] Add the typed HTTP Contract Decision and exception model.
- [x] Bind capability rules and activation bundles to the profile.
- [x] Declare Oxlint, Oxfmt, Vitest, build, and workspace Verification entries.
- [x] Add profile snapshots and focused mutation cases.
- [x] Synchronize the canonical and distributed setup skill trees.

## Acceptance Criteria

- [x] Profile inspection reports exactly the required stack and workspace
      topology, with Inngest and Docker identified as optional.
- [x] The generated architecture contract uses frontend systems and backend
      domain/application/infrastructure terminology without generic normative
      `modules` or `services` buckets.
- [x] An undecided repository is asked for `REST` or `Post-only`; the selected
      value and any explicit exceptions appear in the plan and snapshot.
- [x] PostgreSQL and LogTape are required and blocking when their declared
      evidence is absent.
- [x] The profile emits exact activation bundles and persisted Verification
      commands in deterministic order.
- [x] Profile snapshots contain only 0.0.1 owned versions and are byte-stable.

## Context

- instruction: `CONTEXT.md`
- instruction: `docs/adr/0061-standard-typescript-monorepo-is-opinionated.md`
- instruction: `docs/adr/0063-repositories-own-their-http-contract.md`
- interface: `.agents/skills/setup-context-driven/assets/profiles`
- interface: `.agents/skills/setup-context-driven/tests/test_macro_profiles.py`
- interface: `.agents/skills/setup-context-driven/tests/test_decision_rendering.py`

## Verification

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_standard_typescript_monorepo.py'` — expected: topology, stack, architecture, HTTP decision, capabilities, activations, and Verification entries match the profile contract exactly.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_macro_profiles.py'` — expected: existing profile snapshots remain deterministic alongside the new profile.
- `rtk make skills-sync-check` — expected: canonical and distributed setup skill trees are byte-identical.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goals 2 and 3; Core Features 2, 3, 7, 8, and 11; User Story 2.
- `_techspec.md` → Coverage Map; Integration Points; Build Order 5.
- ADR-0061 → exact Standard TypeScript Monorepo Profile contract.
- ADR-0063 → repository-owned typed HTTP Contract Decision.

## Result

Added the strict `standard-typescript-monorepo` 0.0.1 profile as a declarative
contract. It carries the exact required stack, the two required workspaces,
optional Inngest and Docker modules, frontend systems and backend layered
architecture, universal and profile-specific Repository Capabilities, exact
Skill Activation bundles, and ordered persisted Verification commands.

The asset boundary now validates the profile before rendering. An absent HTTP
Contract Decision remains unresolved with only `REST` and `Post-only` choices;
a resolved decision persists its ordered typed exceptions and source evidence
in byte-stable 0.0.1 plan and snapshot documents. Invalid modes, incomplete
exceptions, implicit defaults, topology/capability drift, activation drift,
and owned-version drift fail at catalog loading or decision normalization.

Verification:

- The pre-change focused run failed to import
  `build_standard_profile_plan`, establishing that the strict profile contract
  and renderer did not exist.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s
  .agents/skills/setup-context-driven/tests -p
  'test_standard_typescript_monorepo.py'` — passed 9 tests after final sync.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s
  .agents/skills/setup-context-driven/tests -p 'test_macro_profiles.py'` —
  passed 9 tests after final sync.
- The full canonical setup suite passed 205 tests before the final task gate.
- `rtk make skills-sync-check` and `rtk git diff --check` — passed.
- `rtk make verify` — passed on the unchanged authorized rerun after the
  sandboxed attempt was blocked by access to the host Go build cache. Both
  canonical and distributed setup suites passed 205 tests, 1,694 Go tests
  passed, both asset catalogs loaded, the Roundfix skill check passed, and the
  CLI built.

Acceptance evidence:

- `test_profile_inspection_reports_exact_stack_topology_and_optional_modules`
  proves the exact 18-member stack, `packages/frontend` and
  `packages/backend`, and optional Inngest/Docker classification.
- `test_architecture_contract_uses_systems_and_layers_not_generic_buckets`
  proves frontend systems, a public system boundary, direct internal imports,
  backend domain/application/infrastructure layers, thin HTTP handlers,
  HTTP-independent use cases, Drizzle-owned persistence, and explicit
  rejection of generic normative `modules` and `services` buckets.
- The two HTTP Contract Decision tests prove unresolved repositories receive
  only `REST` or `Post-only`, resolved plans and snapshots retain typed ordered
  exceptions and source evidence, and invalid modes or incomplete exceptions
  fail closed.
- `test_required_postgresql_and_logtape_evidence_blocks_but_optional_modules_do_not`
  proves missing PostgreSQL and LogTape evidence blocks while absent Inngest
  and Docker do not.
- `test_profile_binds_universal_capabilities_exact_bundles_and_verification`
  proves the universal capability binding, exact ten activation bundles, and
  ordered Oxfmt, Oxlint, Vitest, build, and workspace Verification entries.
- The byte-stability and mutation tests prove repeated plan/snapshot rendering
  is identical, all owned profile/plan/snapshot/marker versions are `0.0.1`,
  and contract drift fails during local asset loading without repository
  writes. The macro suite proves the new snapshot remains deterministic beside
  the three existing profile flows.

Follow-up: Task 08 remains the owner of integrating this resolved profile data
into the repository audit/apply Change Plan and atomic write transaction.
