---
spec: 0146-a-gate-that-runs-the-analyzer
prd: _prd.md
created: 2026-09-18
---

# A gate that runs the analyzer — Technical Spec

## Executive Summary

One step joins the gate, and one contract test keeps that step honest. The
repository Verification runs the Go toolchain's analyzer over the module, and a
test proves both that the composition contains the step and that the analyzer
fails on a fixture carrying a known diagnostic.

The trade-off this design accepts is that a diagnostic anywhere in the module
now blocks every delivery, with no exemption to fall back on. That is the point:
the alternative, an allowlist, is how a gate stops meaning anything.

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

## System Architecture

| Component | Existing package or proposed file | Responsibility |
| --- | --- | --- |
| Gate composition | `Makefile` | Run the analyzer over the module as part of the full and incremental gates. |
| Negative control | `internal/` contract test plus a fixture under `testdata` | Prove the analyzer fails on a known diagnostic and that the gate composes the step. |
| CI | `.github/workflows/ci-verify.yml` | Unchanged: it invokes the gate, which now includes the step. |

No new package, tool, dependency or configuration file is proposed.

## Implementation Design

### The step

The gate gains a target that runs the Go toolchain's analyzer over the module.
It joins the full gate and the incremental gate, so the fast tier and the
complete tier answer the same question about analysis. The analyzer's own exit
status is the step's result: no output is piped through a filter that could
decide the status in its place.

### The negative control

A repository-contract test, in the tag the other repository contracts already
use, does two things.

- It reads the gate's composition and requires the analyzer target to be part of
  the full and incremental gates. A step that is quietly dropped fails here.
- It runs the analyzer over a fixture package that carries one known diagnostic
  and requires a non-zero status naming that diagnostic. An analyzer that stops
  reporting fails here.

The fixture lives under a `testdata` directory, which module-wide package
matching excludes, so the deliberate defect never reaches the gate it exists to
test.

### What stays

Every existing gate step keeps its name, its command and its order. CI keeps
invoking the gate through the same command, so the workflow bounded by this
Spec's grant is left unchanged; the Spec records that rather than editing a file
because it may.

### Interfaces

No Go interface, exported identifier or command-line surface changes.

### Data Models

No entity, schema, stored record or event payload changes.

## API Contracts

1. The repository Verification fails when the module carries an analyzer
   diagnostic, and passes when it does not.
2. The repository Verification's other steps keep their names, commands and
   order, and the command CI invokes is unchanged.
3. No dependency, configuration file or third-party tool is added.

## Coverage Map

- Goal 1 → Gate composition.
- Goal 2 → Negative control.
- Goal 3 → What stays; API Contract 3.
- User Story 1 → Gate composition; API Contract 1.
- User Story 2 → Negative control.
- User Story 3 → What stays; API Contract 2.
- Core Feature 1 → Gate composition.
- Core Feature 2 → Negative control.
- Core Feature 3 → API Contract 3.
- Core Feature 4 → What stays; API Contract 2.
- Success Metric 1 → Testing Approach 2.
- Success Metric 2 → Testing Approach 1 and 3.
- Success Metric 3 → Testing Approach 4.
- API Contracts 1-3 → Gate composition, Negative control, What stays.

## Integration Points

- **CI.** It invokes the repository Verification and therefore the new step; the
  workflow file needs no edit.
- **The Daemon's QA step.** It runs the configured repository Verification, so
  the analyzer now runs there too, as a gate fact rather than an Agent task.
- **Spec 0124.** Capacity measurements, economics and the remaining proposed
  controls stay there.
- **Baseline consumers.** Other repositories compose their own gates; nothing
  here changes what they ship.

## Testing Approach

1. **Composition, at the repository-contract seam.** The test reads the gate's
   definition and requires the analyzer target in both tiers. It fails on the
   tree as it stands today, where neither contains it.
2. **Negative control, at the same seam.** The analyzer runs over a fixture
   carrying one known diagnostic and must exit non-zero naming it. This proves
   the step can fail without making the module fail.
3. **The module is silent.** The gate itself, run over the repository, passes —
   which is only meaningful beside the negative control.
4. **Nothing added.** The test requires no new module dependency, and the
   repository carries no analyzer configuration file.
5. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact, which now includes the analyzer.

## Build Order

1. The negative control and its fixture, proving the analyzer fails on a known
   diagnostic and asserting the composition the next step creates (depends on:
   none).
2. The analyzer step in the gate's full and incremental tiers (depends on: 1).
3. Terminal QA (depends on: 1, 2).

Step 1 before step 2 is deliberate: the control is written against a gate that
does not yet contain the step, so it fails first for the reason it exists.

## Risks & Considerations

- **A gate that blocks everything on one diagnostic.** With no exemption, a
  diagnostic in any package stops delivery until repaired. The module is silent
  today, so the first failure belongs to the change that causes it.
- **A fixture that leaks into the gate.** Module-wide matching excludes
  `testdata`, and the negative control asserts the module stays silent, so a
  leak fails the control rather than hiding.
- **A third-party linter later.** Nothing here forecloses one, and nothing here
  invites one; that decision keeps its own authority.

## Decisions

- **The toolchain's analyzer only.** It ships with Go, needs no version pin and
  no configuration, and answers the question this Spec asks.
- **No suppression mechanism.** An exemption path is how a gate stops meaning
  anything.
- **Control first, step second.** A control written after the step would pass on
  its first run and prove nothing about its own ability to fail.
- **The workflow stays unchanged.** It is bounded by the grant, and the gate is
  composed in one place; editing a second place because permission exists would
  create the duplicate this design avoids.

## Vocabulary Contract

No token is coined. Repository Verification, gate and analyzer are existing
terms used with their existing meanings.
