---
spec: 0146-a-gate-that-runs-the-analyzer
status: archived
created: 2026-09-18
surfaces: [backend, infra]
archived: "2026-09-18"
source_slug: 0146-a-gate-that-runs-the-analyzer
---


# A gate that runs the analyzer

The Go toolchain ships an analyzer that the repository's Verification gate does
not run. For a long time it could not: the agent package reported thirty-one
diagnostics about a mutex copied by value, so composing the analyzer into the
gate would have meant composing a red gate.

Spec 0145 removed those, and the module now reports none. The prerequisite
Spec 0124 named for this work is therefore satisfied, and the question left is
the one it asked: does the gate run the analyzer, or does each Spec remember to
run it over the package it touched?

Remembering does not scale. Spec 0145 proved its own package and stopped there,
which is correct for a slice and useless against a regression arriving through a
package nobody edited. A gate that never fails is decoration, so this Spec also
demands the opposite proof: that the step fails when a diagnostic exists.

This Spec is the first slice carved from Spec 0124, which keeps the capacity
measurements, the economics and the remaining proposed controls.

## Project Constraints

- Identifier strategy: applicable — gate target names and their composition are the repository's own vocabulary; this Spec adds a step and renames none. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the analyzer reads local source and opens no credential store or network transport. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the gate's two-tier composition and its ownership are governed by accepted decisions this Spec must preserve. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0014 applies: the Daemon runs the repository Verification, so this Spec changes what that command does and never who invokes it.
  ADR-0020 applies: a parsed prompt result stays usable despite a later teardown exit, and this Spec changes no Agent result boundary.
  ADR-0038 applies: the Daemon allows one Verification repair per Agent Session, and adding a gate step does not enlarge that budget.
  ADR-0056 applies: Task Capacity and Verification Capacity stay distinct, and composing a step into the gate changes neither.
  ADR-0057 applies: the Daemon owns Task status, which this Spec leaves untouched.
  ADR-0096 applies: the QA gate proves machine facts before it spends an Agent turn, and the analyzer runs as one of those facts rather than as Agent work.
  ADR-0117 applies: a defect is checked by the stage that can produce it, which for an analyzer diagnostic is the repository gate.
  ADR-0127 applies: process residue stays a readiness observation, and no gate step turns it into authority to settle work.
  ADR-0155 applies: the `qa` Task declares the gate's matrix, so this Spec's terminal gate covers what its Requirements name.
  ADR-0156 applies: this Spec declares its Success Metrics and API Contracts as numbered units and names each in a Task.
  ADR-0104 applies: acceptance rests on evidence this Spec did not author, which here is the analyzer's own report over the module as it stands.
- Tooling authority: applicable — the repository Verification is composed in build-tool configuration, which is a Governed Path. Express maintainer authorization: granted 2026-09-18, recorded in [_authorization.md](_authorization.md); bounded files: `Makefile`, `.github/workflows/ci-verify.yml`. The negative control, its fixture and the contract test are ordinary source that no authorization has bounded. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## Goals

- The repository Verification runs the Go analyzer over the module.
- A diagnostic anywhere in the module fails that gate, proved rather than
  assumed.
- No third-party tool, configuration file or suppression mechanism enters the
  repository with it.

## User Stories

1. As a maintainer, I want a copied lock or a misused format verb to fail the
   gate wherever it appears, so that a regression does not wait for someone to
   run the analyzer by hand over that package.
2. As a maintainer, I want proof that the new step can fail, so that a green
   gate means the analyzer ran and found nothing rather than that it silently
   did nothing.
3. As a maintainer, I want the gate to keep its current shape otherwise, so that
   adding this step costs no other behavior.

## Core Features

1. **The gate runs the analyzer.** The repository Verification and its
   incremental sibling run the Go toolchain's analyzer over the module, and a
   diagnostic fails them.
2. **A negative control proves the step can fail.** A contract test runs the
   analyzer over a fixture that carries a known diagnostic and requires it to
   fail, and reads the gate's composition to require the step is present.
3. **Nothing else is added.** No third-party linter, no analyzer configuration
   file, no dependency, no suppression mechanism and no per-package exemption.
4. **The rest of the gate is unchanged.** Every existing step keeps its place,
   its command and its meaning, and the command CI invokes does not change.

## User Experience

A maintainer runs the repository Verification and sees the analyzer run with the
other steps. A change that copies a lock, or misuses a format verb, fails there
instead of in a later reader's analysis.

## Non-Goals / Out of Scope

- Capacity measurements, the economics of the gate, cold-cache convenience
  commands and the remaining proposed controls. Spec 0124 keeps those.
- Adopting a third-party linter or any analyzer beyond the Go toolchain's own.
- A suppression, allowlist or exemption mechanism. A diagnostic is fixed or the
  gate stays red.
- Changing which command CI invokes, or adding a second place where the gate is
  composed.
- Repairing diagnostics in other repositories that consume the Baseline.

## Declared intentional breaks

- A change that introduces an analyzer diagnostic now fails the repository
  Verification where it previously passed. The module reports none today, so the
  first failure will belong to the change that causes it.

## Regression locks

- The gate's existing steps keep their order, their commands and their names.
- The incremental gate stays the fast tier and the full gate stays the complete
  one.
- No configuration file, dependency or tool is added to the repository.

## Acceptance evidence

At least one acceptance row rests on evidence this Spec did not author: the
analyzer's report over the module as it stands, which is empty because Spec 0145
emptied it, and the finding Spec 0123 owns, which recorded thirty-one
diagnostics before that. The replay runs the analyzer over the module and over
the negative control's fixture, and records that the first is silent and the
second fails. Where the analyzer cannot be run, the row records that reason and
does not block.

## Success Metrics

1. The repository Verification fails when the module carries an analyzer
   diagnostic, proved by the negative control rather than by argument.
2. The repository Verification passes on the module as it stands, with the
   analyzer step running.
3. No dependency, configuration file or third-party tool is added.

## Decisions

- **The Go toolchain's analyzer, nothing else.** A third-party linter brings a
  configuration file, a version to pin and a taste to argue about; this Spec
  answers a narrower question.
- **No suppression mechanism.** An exemption is how a gate stops meaning
  anything; a diagnostic is repaired or the gate stays red.
- **One place composes the gate.** CI invokes the repository Verification, so the
  workflow is bounded by the grant and left unchanged unless the analyzer cannot
  run there otherwise.
- **Slice, not portfolio.** Spec 0124 carries six Core Features; this Spec takes
  the analyzer composition and leaves the measurements there.

## Research basis

**Secondbrain.** Consulted before authoring, index first, then a query on
composing static analysis into a build gate and on proving a gate can fail. The
results were dominated by mirrors of this repository, which are references
rather than independent knowledge, and no source changed the design.

The repository's pending Inbox Entries were read before authoring. One is
adjacent and deliberately not consumed here: an authored Verification whose
pipeline hides a tool's exit status, captured while delivering Spec 0145. It
argues for the negative control this Spec requires, and its own repair belongs to
the checker rather than to the gate.

**Exa MCP.** Consultation was attempted, and no Exa MCP tool was available in
this session. No external source was read, so no external validation is claimed.

**Local measurement.** The analyzer's silence over the module, the gate's current
composition, the CI job's single invocation of it, and the exclusion of testdata
directories from module-wide package matching were each read in this repository
before this PRD was written.

## Open Questions

None.
