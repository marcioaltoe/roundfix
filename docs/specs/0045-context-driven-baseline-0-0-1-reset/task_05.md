---
task: task_05
spec: 0045-context-driven-baseline-0-0-1-reset
status: pending
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

- [ ] Define the exact profile topology, stack, and optional modules.
- [ ] Add frontend systems and backend layered architecture clauses.
- [ ] Add the typed HTTP Contract Decision and exception model.
- [ ] Bind capability rules and activation bundles to the profile.
- [ ] Declare Oxlint, Oxfmt, Vitest, build, and workspace Verification entries.
- [ ] Add profile snapshots and focused mutation cases.
- [ ] Synchronize the canonical and distributed setup skill trees.

## Acceptance Criteria

- [ ] Profile inspection reports exactly the required stack and workspace
      topology, with Inngest and Docker identified as optional.
- [ ] The generated architecture contract uses frontend systems and backend
      domain/application/infrastructure terminology without generic normative
      `modules` or `services` buckets.
- [ ] An undecided repository is asked for `REST` or `Post-only`; the selected
      value and any explicit exceptions appear in the plan and snapshot.
- [ ] PostgreSQL and LogTape are required and blocking when their declared
      evidence is absent.
- [ ] The profile emits exact activation bundles and persisted Verification
      commands in deterministic order.
- [ ] Profile snapshots contain only 0.0.1 owned versions and are byte-stable.

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
