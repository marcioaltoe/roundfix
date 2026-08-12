---
spec: 0083-a-gate-that-can-say-no
prd: _prd.md
created: 2026-08-07
---

# A gate that can say no — Technical Spec

## Executive Summary

The Makefile defines `GO := $(RTK) go`, so every gate target — including the
authoritative `test` — runs the toolchain through an output-filtering wrapper.
On the full package set that wrapper was observed returning `0` for a run the
toolchain returned `1` for, and omitting the failing package from its summary.
The repair is to split the variable: an unfiltered invocation for anything whose
exit status is load-bearing, and the wrapper for convenience targets where it is
not.

The trade-off this design accepts: the gate becomes louder. Filtered output is
what made a 26-package run readable, and removing the filter from `test` gives
back the raw volume. That is the correct direction — a gate is read by a
machine and only skimmed by a human, and the compact summary is exactly what
concealed the failure. Targets a human reads keep the wrapper.

The second theme is unrelated: two gates assert facts about the machine and the
repository's authoring history rather than about the code, and they are repaired
by changing what they assert, not by loosening thresholds.

## Project Constraints

- Identifier strategy: not applicable — the `go-cli-tui` Profile does not select
  `identifier.strategy`; nothing here persists a record needing one. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no `docs/agents/backend.md` exists
  and no HTTP surface is touched. Source: `docs/agents/spec-routing.md`.
- Active ADR obligations: applicable — ADR-0081 governs the derived pins that
  `make baseline-digests` rewrites, which this Spec must regenerate rather than
  transcribe if any authorized edit moves them. ADR-0089 (code under test takes
  its environment explicitly) directly informs the flaky-test repair, and
  ADR-0090 (repository facts read in batches, never cached across mutations)
  bounds how the corpus sweep may be re-expressed. Source:
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — the Makefile, four test sources, one test
  fixture, and one move out of an archived Spec. Express maintainer
  authorization: recorded 2026-08-07 in
  `docs/workflow/authorizations/2026-08-07-make-the-gate-honest.md` as "Conceder
  o conjunto completo" and "Sair para um dono semântico". Bounded files:
  `Makefile`; `internal/spec/gate_test.go` (new);
  `internal/speccheck/constraints_characterization_test.go`;
  `internal/speccheck/testdata/corpus-golden.json`;
  `internal/spec/coverage_test.go`; `internal/agent/acpx_runner_test.go`;
  `internal/cli/implement_test.go`; and the authorized move of
  `docs/specs/_archived/0071-verification-cost/coverage-record.json` to
  `docs/references/coverage-record.json`. Pins rewritten by
  `make baseline-digests` are sanctioned fallout under ADR-0081. Source:
  `docs/agents/agent-instructions.md`.

## System Architecture

No package is added and no seam is created. Three existing surfaces change.

**The Makefile's tool variables.** Today one variable serves every target:

```make
GO := $(RTK) go
```

It becomes two, split by whether the caller reads the exit status or the text:

```make
GO      := go          # authoritative: exit status is the result
GO_HUMAN := $(RTK) go  # convenience: output is the result
```

`test`, `spec-budget`, `skills-check`, and `build` — everything `verify`
composes — use `GO`. Convenience targets a human invokes directly keep
`GO_HUMAN`. The wrapper is not removed from the repository; it is removed from
the path where losing an exit status is fatal.

**`internal/speccheck` characterization tests.** `TestCheckCorpusGolden`
compares a pinned per-code count against the live corpus; `TestCheckCorpusBudget`
compares wall-clock against a one-second constant. Both keep their sweep and
change their assertion.

**`internal/spec` coverage contract.** `coverage_test.go` resolves its record
through two path constants, the archived one first. Both are replaced by one
path under `docs/references/`, and the record is moved there with `git mv` so
its history follows.

## Implementation Design

### Interfaces

The gate's honesty needs a test that can only pass if the gate can fail. It is
the load-bearing artifact of this Spec:

```go
// TestAuthoritativeGateReportsFailure runs the repository's authoritative test
// invocation against a throwaway package whose test fails while emitting output
// at the volume that concealed a real failure on 2026-08-07, and requires a
// non-zero exit. It asserts the gate's contract, not the wrapper's behavior.
func TestAuthoritativeGateReportsFailure(t *testing.T)
```

It lives at `internal/spec/gate_test.go`, joining the package that already owns
the repository's verification contracts. It must run the same command shape the
Makefile runs, in a temporary module, so it proves the composed gate rather than
a hand-built approximation. It must not shell out to `make` itself, which would
make the test depend on the developer's `make` and on the repository's own
state.

The corpus counter changes from a pinned equality to a reported derivation:

```go
// The archived corpus is historical and its finding counts move whenever
// authoring lands. Report them; assert only what a change must not do —
// that no active Spec regresses.
```

The budget changes from a wall-clock assertion to a load-independent one. The
sweep's cost is dominated by work per Spec, so the unit becomes work performed
(file reads, or sweep operations) rather than elapsed time, and wall-clock is
logged for humans without gating.

### Data Models

No new persisted entity. `docs/references/coverage-record.json` is the existing
record at a new path; its schema is unchanged.

### API Contracts

None. No Roundfix command's arguments, output, or exit codes change. `make
verify` keeps its name, its composition, and its meaning — it gains only the
ability to fail.

## Coverage Map

- Core Feature 1 (wrapper off the authoritative path) → Makefile `GO` /
  `GO_HUMAN` split.
- Core Feature 2 (regression test proving failure propagates) →
  `TestAuthoritativeGateReportsFailure`.
- Core Feature 3 (global counter stops gating) → `TestCheckCorpusGolden`
  assertion change; `corpus-golden.json`.
- Core Feature 4 (wall-clock stops asserting machine speed) →
  `TestCheckCorpusBudget` unit change.
- Core Feature 5 (timing tests stop flaking) →
  `TestACPXRunCancellationCommandFailuresWarnAndContinue`,
  `TestRunImplementVerificationCapacityAndDaemonStatusIntegratedFlow`.
- Core Feature 6 (coverage contract outside `docs/specs/`) → `coverage_test.go`
  path constants; `git mv` of the record.
- Goal "gate's exit status reflects the suite's" → Core Features 1 and 2.
- Goal "no gate asserts a fact about its machine" → Core Features 4 and 5.
- Goal "no gate fails on unrelated authoring" → Core Feature 3.
- Goal "coverage contract repairable without editing an archived Spec" → Core
  Feature 6.
- Goal "gates that produced signal keep their teeth" → the Non-Goals in
  `_prd.md`, enforced by leaving `spec check`, the published-example contract,
  and the QA gate untouched.

## Integration Points

- **The output-filtering wrapper** — still used, deliberately, by convenience
  targets. This Spec does not diagnose it. Its observed masking is recorded in
  `docs/findings/2026-08-07-the-only-gate-reports-green-on-a-red-suite.md` for
  an upstream report.
- **`make baseline-digests`** — run if any authorized edit moves a derived pin,
  per ADR-0081.

## Testing Approach

Existing seams throughout; one new test.

- `TestAuthoritativeGateReportsFailure` is the only genuinely new seam, and it
  earns its place: without it, the Makefile change is unproven and the defect can
  return silently the next time a target's variable is edited. It builds a
  temporary module so it never depends on the repository being red.
- The two characterization tests keep their existing sweep helpers and change
  only their assertions, so their coverage of `speccheck` behavior is unchanged.
- The two flaky tests are repaired by waiting on their milestone condition rather
  than on elapsed time. Each must be run repeatedly under induced load to show
  the flake is gone; a single green run proves nothing about a timing flake.
- The coverage record move is proven by `TestCoverageEquivalence` passing from
  the new path and by the archived Spec being byte-identical apart from the
  removal the grant authorizes.

Because the gate is currently dishonest, every task in this Spec reads its
result from the unwrapped invocation until Build Order step 1 lands.

## Build Order

1. **Split the Makefile tool variables** so the authoritative gate invokes the
   toolchain directly, and confirm `make verify` now fails on the currently red
   tree.
2. **Add `TestAuthoritativeGateReportsFailure`**, proving a failing test with
   high output volume propagates a non-zero exit through the authoritative
   invocation. (depends on: 1)
3. **Move the coverage record to `docs/references/`** with `git mv` and point
   `coverage_test.go` at the single new path. (depends on: 1)
4. **Repair the coverage regression** so `TestCoverageEquivalence` passes at the
   new path, making the tree green for the first time in this Spec.
   (depends on: 3)
5. **Retire the global corpus counter's gating role**, reporting archived counts
   and asserting only active-corpus regressions. (depends on: 1)
6. **Re-express the corpus budget** in a load-independent unit, logging
   wall-clock without gating on it. (depends on: 1)
7. **Stabilize the two timing-sensitive tests**, each proven by repeated runs
   under induced load. (depends on: 1)
8. **QA gate** — the authored terminal task, proving the gate can say no and
   that no retired gate took a signal-producing gate with it.
   (depends on: 2, 4, 5, 6, 7)

## Risks & Considerations

- **The gate gets louder.** Removing the filter restores full `go test` output
  for 26 packages in every verification run. That is the accepted cost; the
  compact summary is what hid the failure. If the volume proves unworkable, the
  answer is a filter applied to a *copy* of the output while the exit status
  flows from the unfiltered command — never the reverse.
- **Step 1 makes the tree visibly red.** That is not a regression; it is the
  first honest reading. Steps 3 and 4 make it green again. A reviewer seeing red
  between steps 1 and 4 is seeing the Spec work.
- **A flaky test can look fixed.** Both repairs must be shown under induced
  load, repeatedly. A single pass is the exact evidence that misled this session
  once already.
- **The corpus counter may be load-bearing in a way not yet seen.** It is
  retired from gating, not deleted; if it ever catches something real, the
  reported derivation still surfaces it.
- **Scope creep toward auditing every test.** Only the two observed flaking are
  in scope; a broader audit is a different Spec with a different cost.

## Decisions

- The wrapper stays in the repository and leaves only the authoritative path;
  removing it entirely was rejected as unrelated to the defect and costly to the
  human targets that benefit from it.
- The gate's honesty is proven by a test rather than by inspection, because the
  defect is precisely that inspection reported success.
- The corpus counter is retired from gating rather than re-pinned, because
  re-pinning reproduces the same tax on the next ADR.
- The budget changes unit rather than threshold, because any wall-clock
  threshold on a shared machine is a future false alarm.
- The coverage record moves to `docs/references/` rather than gaining an
  exception, following the repository's own rule that durable knowledge moves to
  a semantic owner before archive.
