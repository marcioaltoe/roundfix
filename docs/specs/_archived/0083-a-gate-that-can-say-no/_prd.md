---
spec: 0083-a-gate-that-can-say-no
status: active
created: 2026-08-07
surfaces: [infra, docs]
---

# A gate that can say no

The repository's own rule says the local gate is the only gate. On 2026-08-07
that gate was observed exiting `0` on a working tree whose Go suite exits `1`,
with the failing package omitted from its summary. A gate that cannot say no is
not a gate, and every completion claim routed through it is unverified.

Two further gates cost more than they return. One pins a global finding count
that moves whenever anyone adds an ADR or a Spec; the other asserts wall-clock
under one second on a shared developer machine. Both fired the same day with no
signal — one because two ADRs were authored, one because the machine was busy.
Two timing-sensitive tests flaked in the same session, once in CI and once
locally.

Separately, the live coverage contract sits inside an archived Spec, so the
sanctioned repair for a renamed test is the edit the archival rule forbids.

This is engineering maintenance, not a product change: no Roundfix command
gains, loses, or alters behavior.

## Project Constraints

- Identifier strategy: not applicable — the `go-cli-tui` Baseline Profile does
  not select `identifier.strategy`, and this work persists no record needing a
  generated identifier. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no `docs/agents/backend.md` exists
  because the Profile selects no backend module, and this work opens no HTTP
  surface. Source: `docs/agents/spec-routing.md`.
- Active ADR obligations: applicable — ADR-0081 makes pins rewritten by the
  sanctioned regeneration command deterministic fallout of an authorized edit,
  which this Spec relies on rather than hand-editing; ADR-0089 (code under test
  takes its environment explicitly) governs the flaky-test repair, since a test
  that waits on elapsed time is taking its environment implicitly; and ADR-0090
  (repository facts are read in batches, never cached across mutations) bounds
  how the corpus sweep may be re-expressed when its assertion unit changes. No
  accepted ADR governs the verification gate's composition itself; those rules
  live in `docs/agents/specific-repository.md` as Normative Clauses, and this
  Spec strengthens rather than weakens them. Binding on the authored QA gate:
  ADR-0080, ADR-0088, ADR-0091, ADR-0096, and ADR-0097. ADR-0093 (Spec
  consistency is checked by citation, never by inference) governs how this row
  is validated rather than what the work must do, and is discharged by citing
  every ADR the Spec's artifacts reference. Source:
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — this work edits the Makefile, test sources,
  test data, and moves a file out of an archived Spec. Express maintainer
  authorization: recorded 2026-08-07 in
  `docs/workflow/authorizations/2026-08-07-make-the-gate-honest.md`, granted as
  "Conceder o conjunto completo" and "Sair para um dono semântico". Bounded
  files: `Makefile`; `internal/spec/gate_test.go` (new);
  `internal/speccheck/constraints_characterization_test.go`;
  `internal/speccheck/testdata/corpus-golden.json`;
  `internal/spec/coverage_test.go`; `internal/agent/acpx_runner_test.go`;
  `internal/cli/implement_test.go`; and the authorized move of
  `docs/specs/_archived/0071-verification-cost/coverage-record.json` to
  `docs/references/coverage-record.json`. Source:
  `docs/agents/agent-instructions.md`.

## Goals

- The gate's exit status reflects the suite's exit status, provably, including
  when a failing test emits a large volume of output.
- No gate asserts a fact about the machine it happens to be running on.
- No gate fails because unrelated authoring changed a global count.
- The coverage contract is repairable without editing an archived Spec.
- Every gate that produced signal keeps its teeth.

## Core Features

1. **The authoritative verification invocation does not pass through an
   output-filtering wrapper.** The gate invokes the toolchain directly. Wrapper
   use remains available for human-facing convenience targets, where a lost exit
   status changes nothing.

2. **A regression test proves the gate can fail.** A deliberately failing test,
   emitting output at the volume that triggered the masking, must make the
   authoritative invocation exit non-zero. Without this the repair is a belief.

3. **The global corpus counter stops gating unrelated authoring.** The finding
   count over archived Specs is derived and reported rather than pinned, or the
   gate is retired — either way, adding an ADR or a Spec no longer requires
   editing a recorded number.

4. **The wall-clock budget stops asserting machine speed.** The corpus sweep's
   performance signal becomes advisory, or is expressed in a unit that does not
   depend on competing load. It must not fail a build because another process
   was busy.

5. **The two timing-sensitive tests stop flaking under load.** They wait on the
   condition they mean rather than on elapsed time, or declare the resources
   they need.

6. **The coverage contract lives outside `docs/specs/`.** Its record moves to
   `docs/references/`, and the test reads it there, so a renamed test is
   repaired without touching an archived Spec.

## Non-Goals / Out of Scope

- Diagnosing or fixing the wrapper itself. Its masking is reproduced and
  reported; the tool belongs to the maintainer and is repaired or reported
  upstream separately. This Spec only stops the gate from depending on it.
- Retiring `TestCoverageEquivalence`. It caught a real regression on the day
  this Spec was written and keeps its teeth.
- Weakening `spec check`, the published-example contract, or the authored QA
  gate. Each produced signal on 2026-08-07.
- Auditing every other test for flakiness. Only the two observed failing under
  load are in scope.
- Changing any Roundfix command's behavior.
