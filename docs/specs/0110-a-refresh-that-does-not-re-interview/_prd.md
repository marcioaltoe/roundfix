---
spec: 0110-a-refresh-that-does-not-re-interview
status: active
created: 2026-08-12
surfaces: [backend, cli]
---

# A refresh that does not re-interview, and a prompt that ends

Spec 0082 made the Baseline refresh manifest-driven and explicitly left three
observations from the same investigation live and unowned. A sealed prompt was
observed alive past ten minutes against a declared five-minute bound, which
affects first adoption and every other caller of that path. The readiness
diagnostic and the shell disagree about which binary a runtime resolves to, and
the original finding is explicit that nobody should act on the shim it noticed
without proving it is the cause. And selecting a replacement profile partway
through the interview re-asks decisions already answered, including ones both
profiles share. Two of the three have unknown root causes on the record, so this
Spec reproduces before it designs.

## Project Constraints

- Identifier strategy: applicable — Setup Manifest, Baseline Profile, Sealed
  Prompt, and ACP Runtime are glossary terms this Spec investigates the behavior
  of, and a reproduction may show a term is missing for a state nobody named. The
  closing node checks whether the work introduced or changed a term. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no network boundary, credential, or
  request is created or read. The work is subprocess lifetime, executable
  resolution, and interactive prompt sequencing. Source:
  `docs/agents/specific-repository.md`.
- Active ADR obligations: applicable — ADR-0111 makes an unobserved Verification
  unknown rather than a verdict, which is the discipline this Spec applies to its
  own two unproven observations: neither is treated as a defect until reproduced.
  No accepted ADR governs sealed-prompt lifetime or interview sequencing, which is
  why the rules this Spec may add are new. Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — no protected tooling mutation is proposed
  or authorized. The work is production Go in the agent and CLI baseline packages
  plus their tests. Source: `docs/agents/agent-instructions.md`.

## Goals

1. A sealed prompt cannot outlive its declared bound.
2. Which executable a runtime resolves to is one answer, reported consistently.
3. Changing a profile mid-interview does not re-ask what was already answered.
4. Each of the three observations is either reproduced and fixed, or recorded as
   no longer reproducible with the evidence that settled it.

## User Stories

1. As a Supervisor running an adoption or refresh, I want a sealed prompt to end
   at its declared bound, so that a stalled run does not consume a night.
2. As a maintainer reading readiness output, I want the binary it names to be the
   one a run actually launches, so that a healthy diagnostic is not compatible
   with a hung run.
3. As a maintainer changing a profile during a refresh, I want previously answered
   decisions carried forward, so that a twenty-one-prompt interview does not
   restart at prompt eighteen.

## Core Features

1. **Reproduction before design.** Each of the three observations is exercised
   against the current tree first. An observation that no longer reproduces is
   recorded as such with the evidence, and its fix is not designed.
2. **A bounded sealed prompt.** A sealed prompt ends at its declared bound
   including when a descendant outlives the direct child, so no caller of that
   path can hang past its own timeout.
3. **One executable answer.** Readiness reporting and the sealed prompt runner
   resolve a runtime the same way, so a maintainer reading one can rely on the
   other.
4. **An interview that carries answers forward.** Selecting a replacement profile
   preserves the decisions already answered that the new profile also declares,
   and asks only what the change actually opened.

## User Experience

A refresh that changes profile asks only the decisions the new profile adds,
naming how many it carried forward. A sealed prompt that reaches its bound ends
and says so, rather than leaving a run with no output and no end.

## Non-Goals / Out of Scope

- Re-opening the manifest-driven update that Spec 0082 delivered.
- Changing the interview's decision set, defaults, or prompt wording beyond what
  carrying answers forward requires.
- Changing runtime configuration, adapter selection, or the capability proof.
- Removing or replacing the sealed prompt mechanism.

## Success Metrics

- A sealed prompt whose child spawns a descendant that outlives it still ends at
  its declared bound.
- Readiness output and the sealed prompt runner name the same executable for one
  runtime on a machine where the shell resolves more than one candidate. Source:
  a developer machine observed on 2026-08-07, outside this Spec's artifacts, where
  the diagnostic and the shell disagreed.
- A profile change mid-interview re-asks no decision both profiles declare,
  measured against the twenty-one-prompt run recorded on 2026-08-07.
- Each of the three observations ends this Spec either fixed with a test or
  recorded as not reproducible with its evidence.

## Decisions

- Reproduction is the first work, not the last, because two of the three
  observations carry `unknown` root causes and designing a fix for an unproven
  cause is how a Spec ships a change that fixes nothing.

## Open Questions

- Whether the sealed-prompt timeout observation still reproduces. The path now
  passes an explicit bound to the subprocess in addition to the context deadline,
  so it may already be closed; the default is to prove it before assuming either
  way.
- Whether resolving one executable answer means the runner adopts the diagnostic's
  resolution or both adopt a third shared one. The default is a shared resolution
  used by both, so neither surface can drift from the other again.
