---
spec: 0104-a-gate-that-cannot-certify-its-own-cache
status: active
created: 2026-08-12
surfaces: [infra, docs]
---

# A gate that cannot certify its own cache

The repository's guidance tells a reader to clear the test cache before trusting
the verification gate, and the instruction does not work: the gate exports its own
cache directory, so the documented command clears a different one and leaves the
gate's exactly as stale. Measured during one Spec, a real regression reported the
gate green with the offending package marked cached. This is the same false-green
class the repository already carries a hard rule about — a gate reporting success
it did not earn — except here the documented workaround is part of the defect
rather than the repair. Anyone following the current instruction believes they
have a cold cache and does not.

## Project Constraints

- Identifier strategy: applicable — Verification and the repository's named gate
  commands are glossary terms this Spec adds to, and a cold-run target is
  vocabulary the guidance must name consistently. The closing node checks whether
  the work introduced or changed a term. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no network boundary, credential, or
  request is created or read. The work is build-tool configuration and the
  guidance that describes it. Source: `docs/agents/specific-repository.md`.
- Active ADR obligations: applicable — ADR-0111 makes an unobserved Verification
  unknown rather than a verdict, which is precisely what a fully cached suite is:
  it observed nothing about the current tree. No accepted ADR governs the gate's
  cache, which is why the rule this Spec may add is new. ADR-0135 makes an absent diagnostic a reported state in the repair prompt an Agent reads; this Spec changes which cache the gate clears and writes no Verification feedback, so it does not apply. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — the verification Makefile gains a target that
  clears the cache the gate uses, and the guidance stops naming a command that
  clears a different one. Express maintainer authorization:
  `docs/workflow/authorizations/2026-08-12-the-authoring-and-baseline-corrections.md`,
  granted 2026-08-12. Bounded files: `Makefile` and
  `docs/agents/specific-repository.md`. Source:
  `docs/agents/agent-instructions.md`.

## Goals

1. Clearing the gate's cache is one command that clears the right one.
2. A reader does not have to know which cache is in effect to trust the gate.
3. Whether a fully cached pass may be reported as a pass is settled deliberately.

## User Stories

1. As a maintainer about to trust the gate, I want one command that clears the
   cache the gate actually uses, so that a cold run is cold.
2. As a maintainer following the repository's guidance, I want it to name a
   command that works, so that I am not told I have a cold cache when I do not.
3. As a maintainer reading a green gate, I want to know whether it observed the
   tree or only its cache, so that a pass means what I take it to mean.

## Core Features

1. **A target that clears the gate's own cache.** One named command clears the
   cache the gate exports, so no reader has to know which cache is in effect.
2. **Guidance that names the working command.** The repository guidance points at
   that target instead of the one that clears an unrelated cache.
3. **A settled position on a fully cached pass.** Whether the gate may report
   success when its entire suite was served from cache is decided and implemented,
   since a cached pass proves the cache rather than the tree.

## User Experience

The named target reports what it cleared. If the gate refuses a fully cached pass,
the refusal says so and names the target that makes the run cold, rather than
failing with a result the reader must interpret.

## Non-Goals / Out of Scope

- Changing what the gate runs, its parallelism, or its package set.
- Changing the cache location the gate exports.
- Investigating the wrapper that once masked a package-level failure in a
  multi-package run; the gate no longer routes through it and that question is
  recorded rather than answered here.
- Any change to CI's own verification invocation.

## Success Metrics

- The regression that reported the gate green with a cached package is caught when
  the named target runs first, reproduced from the measured case on 2026-08-08.
- The guidance names only commands that clear the cache the gate uses.
- A fully cached suite produces the behavior this Spec decides on, proven by
  exercising it.

## Decisions

- The fix is a target rather than a documented environment variable, because a
  reader who has to know which cache is in effect is the reader this defect
  already caught.

## Open Questions

- Whether the gate refuses a fully cached pass, warns on one, or accepts it. A
  refusal is the honest reading of an unobserved verification and also the most
  disruptive to a fast local loop, so this is a product decision rather than a
  derivation.
- Whether the incremental gate and the complete gate answer that question
  differently. The default until settled is that the incremental gate may accept a
  cached pass, since its purpose is reuse, while the complete gate may not.
