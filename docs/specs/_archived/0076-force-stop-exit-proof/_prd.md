---
spec: 0076-force-stop-exit-proof
status: archived
created: 2026-08-04
surfaces: [backend]
archived: "2026-08-04"
source_slug: 0076-force-stop-exit-proof
---


# Force Stop exit proof

Force Stop is one of the strongest guarantees in the glossary: it "proves owner
identity, cancels registered Agent Sessions, terminates the recorded owning
process, and completes the Run as Stopped **only after owner exit is proven**."
The test that proves the hardest half of that — that a process ignoring
`SIGTERM` is still terminated — has never proved it.

`TestOwnerProcessControllerForceKillExitProof` starts a helper in `ignore`
mode, expecting a live process that refuses `SIGTERM` so the controller must
escalate. The helper blocks on `select {}`, which Go's runtime reports as
`all goroutines are asleep - deadlock!` and kills with exit 2. The process the
test believes it is force-killing is already dead of its own accord.

The parent compounds it: `cmd.Wait` is started in a goroutine before the
readiness line is read from the stdout pipe, so the helper's immediate exit can
close that pipe first and surface as `read |0: file already closed`.

Both were found by Spec 0064's QA gate, which recorded the decisive probe: run
the parent and helper together fifty times and the parent passes fifty times
while every child emits the fatal deadlock. A green test, a defeated purpose.

## Project Constraints

- Identifier strategy: not applicable — no project-owned Internal Identifier is
  created; the helper modes, test names, and process states keep their existing
  contracts. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the governing clause prohibits
  reading, printing, committing, or generating secrets and forbids inventing
  authentication, authorization, transport, or deployment policy; this Spec
  changes in-process test coordination and opens no transport. Source:
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0089 requires code under test to
  take its environment explicitly, which is what the helper's liveness and the
  parent's readiness handshake must do instead of depending on scheduling.
  ADR-0080 owns QA verdict semantics and ADR-0091 owns the authored QA gate as
  a typed Task node, under which this Spec's own graph is authored. ADR-0093
  surfaces as a relation candidate because it cites ADR-0080; it does not apply
  — it governs the Spec Consistency Check's detection boundary, and this Spec
  changes no detector. Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — no protected tooling mutation is proposed
  or authorized; the surface is Go test code under `internal/store`. Source:
  `docs/agents/agent-instructions.md`.

## Goals

- A helper started in `ignore` mode stays alive until something kills it, so
  the controller's escalation is what ends it.
- Readiness is consumed without racing process exit, so the parent never fails
  on a closed pipe.
- A helper that exits before the controller acts fails the test instead of
  passing it.

## Core Features

1. `ignore` mode blocks in a way the Go runtime does not treat as a deadlock,
   and the process survives `SIGTERM`.
2. Readiness and process-exit observation are ordered so neither can close the
   pipe the other still needs.
3. A regression rejects premature helper exit, naming it, so the proof cannot
   silently return to passing for the wrong reason.
4. The force-kill proof asserts that the controller's escalation caused the
   exit, not merely that the process is gone.

## Non-Goals / Out of Scope

- Changing Force Stop's product behavior, its owner-identity rules, or the
  `--owner-identity-unreadable` last-resort flag. The contract is correct; only
  its proof is broken.
- Rewriting unrelated `internal/store` tests or the process controller itself.
- Adding timing tolerances, retries, or sleeps to make the existing shape pass.

## Success Metrics

- Run in isolation, the `ignore` helper emits readiness and stays alive until
  killed; its exit status reflects the kill, not a runtime deadlock.
- Repeated fifty times, parent and helper both pass, with no child emitting a
  fatal error.
- A deliberately premature helper exit fails the regression, proving the new
  assertion can fail.

## Decisions

- The one-second-scale coordination is fixed by ordering, not by widening a
  timeout. A budget that hides a race is the defect one step later.
- The proof asserts causation — the controller's escalation ended the process —
  because "the process is gone" is satisfied by a process that crashed on its
  own, which is exactly how this defect survived.

## Open Questions

None.
