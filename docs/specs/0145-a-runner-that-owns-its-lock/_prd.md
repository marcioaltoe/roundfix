---
spec: 0145-a-runner-that-owns-its-lock
status: active
created: 2026-09-18
surfaces: [backend]
---

# A runner that owns its lock

The acpx runner keeps the state of every Agent Session it has ensured, warmed,
started work on, or assigned a selection to, in maps guarded by a mutex it holds
as a field. Sixteen of its methods take that struct by value.

Callers already hold the runner through a pointer, so each of those calls copies
the struct out of the pointer. The maps are references and survive the copy; the
mutex does not. Every copy locks its own mutex while writing the same maps, so
the guard protects nothing across calls, and the analyzer says so thirty-one
times over.

Nothing has been observed corrupting a session map yet. What is observed is the
noise: every analysis of this package reports the same thirty-one diagnostics,
which is how a real one would be missed.

This Spec is the first slice carved from Spec 0123, which keeps the adapter
capability discovery, the identifier rules and the access-policy readiness.

## Project Constraints

- Identifier strategy: applicable — Agent Session identifiers, runtime identifiers and selection assignments keep their meaning; this Spec changes how the runner's methods receive their receiver and coins no identifier. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the change is method receivers and state ownership inside one package; no credential, transport or provider call changes. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — Agent Session lifecycle and cancellation are governed by accepted decisions this Spec must preserve. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0020 applies: a parsed prompt result stays usable despite a later teardown exit, so the cancellation path keeps the behavior it has today.
  ADR-0155 applies: the `qa` Task declares the gate's matrix, so this Spec's terminal gate covers what its Requirements name.
  ADR-0156 applies: this Spec declares its Success Metrics and API Contracts as numbered units and names each in a Task.
  ADR-0104 applies: acceptance rests on evidence this Spec did not author, which here is the analyzer's own report over the unchanged package.
- Tooling authority: applicable — no protected tooling mutation proposed or authorized. The repair is ordinary source in `internal/agent`, which is not a Governed Path, and its test files are ordinary too. This Spec adds no linter, no analyzer configuration and no repository gate; it runs the analyzer as a Task Verification over the package it changes. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## Goals

- The runner's state has one owner, so a lock taken by one call is seen by the
  next.
- The analyzer reports nothing about copied locks in this package.
- Session lifecycle, selection assignment and cancellation behave exactly as
  they do today.

## User Stories

1. As a maintainer reading an analyzer report for this package, I want it empty,
   so that the next real diagnostic is visible instead of buried under
   thirty-one repeats.
2. As a maintainer changing the runner, I want its guarded state to have one
   owner, so that adding a field does not silently add a race.
3. As a Supervisor running Specs, I want session lifecycle and cancellation
   unchanged, so that this repair costs nothing in behavior.

## Core Features

1. **One owner for the runner's guarded state.** Every method that reads or
   writes the runner's session maps receives the runner through a pointer, so
   the mutex and the maps it guards travel together.
2. **The analyzer reports nothing here.** Running the analyzer over this package
   produces no diagnostic about a lock passed by value.
3. **Behavior is unchanged.** Session ensure, warm, work-start, selection
   assignment, probing and cancellation keep their current observable behavior,
   proved by the package's existing tests passing unedited wherever they do not
   name a receiver.

## User Experience

A maintainer runs the analyzer over the package and reads nothing. Runs behave
as before: the same sessions are ensured, the same selections assigned, the same
cancellations honored.

## Non-Goals / Out of Scope

- Adapter capability discovery for specific models, the public identifier rules,
  access-policy readiness, and proving a selection before Run creation. Spec
  0123 keeps those.
- Adding a linter, an analyzer configuration, or a repository-wide gate that
  runs one. This Spec proves its own package and stops there.
- Changing the runner's public surface, its constructor, or how callers obtain
  it — they already hold a pointer.
- Redesigning session state, its keys, or its lifetime.
- Touching any other package's analyzer diagnostics.

## Declared intentional breaks

- A caller that copied the runner by value and relied on the copy having its own
  state would now share state. No such caller exists: the only holder keeps a
  pointer, which is what makes the current copying pointless.

## Regression locks

- The package's existing tests pass without changes to their assertions.
- Cancellation still reports its teardown the way ADR-0020 requires.
- No exported identifier is renamed, added or removed.

## Acceptance evidence

At least one acceptance row rests on evidence this Spec did not author: the
analyzer's report over the unchanged package, which names thirty-one
diagnostics, each a method passing the runner's mutex by value. The replay runs
the same analyzer after the change and shows the count at zero. Where the
analyzer cannot be run, the row records that reason and does not block.

## Success Metrics

1. The analyzer reports zero copied-lock diagnostics for this package, against
   thirty-one on the unchanged tree.
2. The package's tests pass with no assertion edited.
3. No exported identifier of the package changes.

## Decisions

- **Pointer receivers, not a state struct.** The only holder already keeps a
  pointer, so moving the guarded state behind another indirection would add a
  layer to fix a copy nobody wanted.
- **The analyzer runs as this Spec's own Verification.** Making it a repository
  gate is a tooling decision with its own authority, and this slice does not
  take it.
- **Slice, not portfolio.** Spec 0123 carries five Core Features; this Spec takes
  the state-ownership repair and leaves capability discovery there.

## Research basis

**Secondbrain.** Consulted before authoring, index first, then a query on mutex
copying and state ownership in Go. The results were dominated by mirrors of this
repository, which are references rather than independent knowledge, and no
source changed the design.

The repository's pending Inbox Entries were read before authoring; none is a
source for this Spec. The defect was already recorded as a finding adopted by
Spec 0123, which owns it, and this Spec adopts nothing: it links the owner's
copy at
[references/2026-09-08-acpx-runner-value-receivers-copy-its-mutex.md](../0123-runtime-readiness-and-model-capabilities/references/2026-09-08-acpx-runner-value-receivers-copy-its-mutex.md).
That finding recorded thirty-one diagnostics on 2026-09-08; the same count was
measured again for this PRD.

**Exa MCP.** Consultation was attempted, and no Exa MCP tool was available in
this session. No external source was read, so no external validation is claimed.

**Local measurement.** The thirty-one diagnostics, the sixteen value receivers,
the twenty-five pointer receivers already present, and the single pointer-held
field through which callers reach the runner were each read in this repository
before this PRD was written.

## Open Questions

None.
