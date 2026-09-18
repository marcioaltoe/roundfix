---
spec: 0145-a-runner-that-owns-its-lock
prd: _prd.md
created: 2026-09-18
---

# A runner that owns its lock — Technical Spec

## Executive Summary

Sixteen methods take the acpx runner by value while the only holder keeps a
pointer to it. Each call copies a struct whose maps survive the copy and whose
mutex does not, so the guard is per-call and the maps are shared. Converting
those receivers to pointers gives the state one owner and empties the analyzer's
report for this package.

The trade-off this design accepts is that it changes receivers across a large
file without changing behavior, so its proof is the package's own tests passing
unedited plus the analyzer's count falling to zero — not a new test for a
behavior that is meant to stay identical.

## Project Constraints

- Identifier strategy: applicable — Agent Session identifiers, runtime identifiers and selection assignments keep their meaning; this Spec changes how the runner's methods receive their receiver and coins no identifier. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the change is method receivers and state ownership inside one package; no credential, transport or provider call changes. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — Agent Session lifecycle and cancellation are governed by accepted decisions this Spec must preserve. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0020 applies: a parsed prompt result stays usable despite a later teardown exit, so the cancellation path keeps the behavior it has today.
  ADR-0155 applies: the `qa` Task declares the gate's matrix, so this Spec's terminal gate covers what its Requirements name.
  ADR-0156 applies: this Spec declares its Success Metrics and API Contracts as numbered units and names each in a Task.
  ADR-0104 applies: acceptance rests on evidence this Spec did not author, which here is the analyzer's own report over the unchanged package.
- Tooling authority: applicable — no protected tooling mutation proposed or authorized. The repair is ordinary source in `internal/agent`, which is not a Governed Path, and its test files are ordinary too. This Spec adds no linter, no analyzer configuration and no repository gate; it runs the analyzer as a Task Verification over the package it changes. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## System Architecture

| Component | Existing package or proposed file | Responsibility |
| --- | --- | --- |
| acpx runner | `internal/agent` | Own its guarded session state through pointer receivers. |
| Package tests | `internal/agent` | Prove the behavior did not move, unedited wherever they do not name a receiver. |

No new package, file, type or field is proposed.

## Implementation Design

### One receiver kind

Every method of the runner that reads or writes its guarded state takes a
pointer receiver, joining the twenty-five that already do. A method that touches
no state may keep its receiver, but uniformity is cheaper to read than a rule
about which methods are exempt, so all sixteen move.

Call sites do not change: the only field holding a runner is already a pointer,
and Go takes the address automatically for an addressable value. Test helpers
that construct a runner as a value gain an address-of where the compiler asks
for one.

### What must not move

Session ensure, warm, work-start, selection assignment, probing, the codex spawn
resolution and cancellation keep their behavior. The cancellation path keeps
reporting a teardown exit after a parsed result, which ADR-0020 requires, and no
exported identifier changes.

### Interfaces

No exported signature changes. The receiver kind is not part of the package's
public contract here, because the type is consumed through a pointer.

### Data Models

No entity, schema, stored record or event payload changes.

## API Contracts

1. The analyzer, run over this package, reports no diagnostic about a lock
   passed by value.
2. No exported identifier of `internal/agent` is renamed, added or removed.
3. Every observable behavior of the runner — session lifecycle, selection
   assignment, probing and cancellation — is unchanged.

## Coverage Map

- Goal 1 → One receiver kind.
- Goal 2 → One receiver kind; API Contract 1.
- Goal 3 → What must not move; API Contract 3.
- User Story 1 → API Contract 1.
- User Story 2 → One receiver kind.
- User Story 3 → What must not move.
- Core Feature 1 → One receiver kind.
- Core Feature 2 → API Contract 1; Testing Approach 1.
- Core Feature 3 → Testing Approach 2; API Contracts 2-3.
- Success Metric 1 → Testing Approach 1 and 3.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → API Contract 2.
- API Contracts 1-3 → One receiver kind, What must not move.

## Integration Points

- **Agent Session lifecycle.** The runner is reached through one pointer field;
  nothing else holds it, so no caller's semantics change.
- **Spec 0123.** Adapter capability discovery, identifier rules and
  access-policy readiness stay there, and it owns the finding this Spec repairs.
- **Repository gates.** The analyzer is not added to any gate; this Spec runs it
  over its own package as Task Verification.

## Testing Approach

1. **Analyzer count.** The Task's Verification runs the analyzer over this
   package and requires an empty report. It fails on the tree as it stands
   today, where it names thirty-one diagnostics.
2. **Existing tests unedited.** The package's tests run as they are. Only a test
   helper that constructs a runner as a value may gain an address-of; no
   assertion changes.
3. **Outside evidence.** The finding Spec 0123 owns recorded thirty-one
   diagnostics on 2026-09-08, before this Spec existed. The replay runs the same
   analyzer and records the count before and after. Where the analyzer cannot be
   run, the row records that reason and does not block.
4. **Repository gate.** The terminal QA Task records the Daemon's `make verify`
   result as a fact.

## Build Order

1. Convert the sixteen value receivers to pointer receivers and fix whatever the
   compiler names (depends on: none).
2. Terminal QA (depends on: 1).

## Risks & Considerations

- **A wide diff with no behavior change.** That is what makes it reviewable: the
  package's tests are the control, and any behavior change shows up as a failing
  one rather than as a judgement call.
- **A copy someone wanted.** If any code depended on a per-call copy of the
  session maps, it depended on a guard that was already broken; no such caller
  exists, and the Spec says so where a reader meets it.

## Decisions

- **All sixteen, not the subset that touches state.** A rule about which methods
  may keep a value receiver is a rule someone has to apply correctly next time.
- **No repository analyzer gate.** Adding one is a tooling decision with its own
  authority; this Spec proves its own package.

## Vocabulary Contract

No token is coined. Runner, Agent Session and selection assignment are existing
terms used with their existing meanings.
