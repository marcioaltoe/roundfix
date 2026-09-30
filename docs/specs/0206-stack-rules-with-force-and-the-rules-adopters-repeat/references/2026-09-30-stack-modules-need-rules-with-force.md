---
type: feat
status: promoted
created: 2026-09-30
spec: 0206-stack-rules-with-force-and-the-rules-adopters-repeat
reason: null
---

# Stack modules need rules with force, and a repository with several stacks

## Opportunity

The Baseline's Go, Rust, CLI and TUI modules carry guidance without clause force, and the knowledge each project needs lives in hand-written repository extensions. The 2026-09-30 audits found the same rules written by hand in several repositories:

- Rust (onioncry): typed errors and no panics on user paths, a thin binary over a library, tests that run the built binary, dependencies through Cargo with `--locked`, MSRV and edition as a contract.
- Go (Roundfix): gate composition from the module's toolchain, test isolation, build-constrained files checked across operating systems, a race-detector policy.
- TypeScript (conexus, fluxus, vortex, oraculum): tests through the package script, contract-first HTTP, database mutation authority, fixtures typed from the schema, the architecture checker in the gate, lint warnings that block.

No profile can express a Go CLI with a Hono backend and a React frontend in one repository, which is the planned shape of argus: a profile takes one setup, and the repository Verification is one command.

## Value

Adopters stop re-deriving the same rules, agents get them with force, and a repository with several stacks can adopt the Baseline at all.

## Shape

A dedicated wave after v0.22.0, with the maintainer's design decisions first:

- whether stdlib-only is a Baseline rule for every Go adopter;
- the default Rust error policy;
- a built-in composed profile against module additions in repository-owned profiles;
- whether onioncry adopts the Baseline first, so a real Rust adopter validates the rules;
- whether the frontend layout is one allowed structure or a decision.

Also covered: the `cut-release` skill is dispatched by the Rust module but can only be invoked by the user, and unused skills (`zustand`, `ai-sdk`, `motion`, `tanstack-table`) are required by the TypeScript setup.

Left out of the v0.22.0 corrections on purpose, because each one needs a design decision or changes what adopters must do:

- rendering the declared workspace paths in the backend and frontend guides;
- which module owns the production-code and debugging skill triggers, which today are TypeScript-owned with language-neutral wording;
- the wording of the generic-layers prohibition, which has no scope and reads against domain services;
- the auth owner name fixed in the backend guide;
- an unconditional rule that lint warnings block completion.

Evidence: the three stack audits of 2026-09-30, summarized in secondbrain `inbox/skills/2026-09-30-skills-vendorizadas-contradizem-o-baseline-ou-estao-quebradas.md` and in the agent's session reports.
