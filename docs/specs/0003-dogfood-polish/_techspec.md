---
spec: 0003-dogfood-polish
prd: _prd.md
created: 2026-07-05
---

# Dogfood Polish — Technical Spec

## Executive Summary

Eleven diagnosed findings, batched. Every change lands inside an existing
component with existing tests as the safety net; the accepted trade-off is
breadth over depth — nine small surfaces move in one Spec, so the discipline
that keeps this safe is per-task Verification against the exact contract
tests each surface already has. No new packages, no schema changes, no ADRs.

## System Architecture

All changes extend existing components:

- `internal/daemon` — commit-subject derivation (`TaskCommitMessage`,
  `QACommitMessage`), infrastructure-error propagation.
- `internal/agent` — `InfrastructureError` gains a bounded stderr tail in its
  message; `LogPath` reused for spec Runs.
- `internal/cli` — implement header lines, `stop --spec`, Interactive Input QA
  field wiring, artifact-dir plumbing into `TaskPlan`, spec-discovery
  diagnostics on the picker path.
- `internal/tui` — QA yes/no field in the implement field set.
- `internal/spec` — `ListActive` gains a companion that also reports skipped
  folders with reasons.
- Test helpers across `internal/{daemon,cli,store,spec}` — hermetic git repo
  construction.
- Docs: repo agent instructions, canonical Roundfix skill (+ `make
  skills-sync`).

## Implementation Design

### Interfaces

```go
// internal/daemon — subject derivation only; trailers unchanged
func TaskCommitMessage(slug string, task spec.Task) string // first rune lowered
func QACommitMessage(slug, verdict string) string          // "docs: qa report for <slug> (<verdict>)"

// internal/agent
type InfrastructureError struct {
    ExitCode int
    Reason   string
    Stderr   string // Error() now includes a bounded tail (last ~10 lines / 1 KiB)
}

// internal/daemon
type TaskPlan struct {
    // existing fields...
    ArtifactDir string // agent logs land here, mirroring review Runs
}

// internal/spec
type SkippedSpec struct{ Dir, Reason string }
func ListActiveDetailed(gitRoot string) ([]Spec, []SkippedSpec, error)
```

```go
// test helpers (per package, no shared util package): temp git repos run with
// GIT_CONFIG_GLOBAL=/dev/null, GIT_CONFIG_SYSTEM=/dev/null and explicit
// -c user.name/-c user.email/-c commit.gpgsign=false on every git call,
// matching the daemon's own explicit-config discipline.
```

### Data Models

None. No store, journal, or config schema changes. (`stop --spec` computes the
existing spec target key; no new columns.)

### API Contracts

- `roundfix stop --spec <slug>` — resolves the Active Run for
  `("spec", "<git_root>#<slug>")` in the current repository; error shapes and
  exit codes mirror `--pr`. Mutually exclusive with the other stop selectors.
- Implement Run header: `Budget:` and `Round:` lines removed for spec Runs;
  everything else byte-identical.
- Interactive Input (implement): new final field `QA gate [y/N]`; `--qa` flag
  presets it; `--no-input` behavior unchanged.
- Spec picker: skipped folders print one stderr diagnostic line each
  (`skipped docs/specs/<dir>: <reason>`); stdout untouched.
- QA Report commit subject: `docs: qa report for <slug> (<verdict>)`.
- Task commit subject: first rune lowercased; trailers unchanged.

## Coverage Map

- Story 1 → TaskCommitMessage/QACommitMessage (findings 1, 9)
- Story 2 → implement header rendering (findings 2, 3)
- Story 3 → hermetic git test helpers (finding 21)
- Story 4 → InfrastructureError stderr tail (finding 27)
- Story 5 → stop --spec (finding 13)
- Story 6 → Interactive Input QA field (finding 7)
- Story 7 → ListActiveDetailed + picker diagnostics (finding 6)
- Story 8 → TaskPlan.ArtifactDir + LogPath reuse (finding 8)
- Feature 9 → docs + skill sync (finding 14)

## Integration Points

None new — git, acpx, and the Run Database are touched only through existing
wrappers.

## Testing Approach

Existing seams only. Each task's tests: subject-derivation table tests
(including non-letter first runes); header golden-string asserts updated
deliberately; a canary test that sets `commit.gpgsign=true` in a scoped config
and proves repo-creating helpers stay green; error-message tail assertions;
stop-target resolution over a real temp store; collector-driven QA field test;
picker diagnostics via a broken fixture folder; spec-Run log path assert under
a configured artifact dir. The full suite is the regression net; stdout/exit
contract tests change only for the two removed header lines.

## Build Order

1. Commit subject normalization (no deps)
2. Implement header cleanup (no deps)
3. Hermetic git test helpers across packages (no deps)
4. Infrastructure-error stderr tail (no deps)
5. `stop --spec` (no deps)
6. Interactive Input QA field (no deps)
7. Spec-Run agent logs under the Artifact Directory (no deps)
8. Spec discovery diagnostics (no deps)
9. Docs and skill sync (depends on: 1, 5, 6 — documents shipped behavior)

## Risks & Considerations

- Nine independent tasks touching one CLI package invite mechanical merge
  friction if executed in parallel — the graph stays sequential-friendly and
  tasks touch disjoint files by construction.
- Header-line removal edits existing golden asserts; each removal must be the
  deliberate diff, not collateral.
- The stderr tail must be bounded and secret-agnostic (no env echoing).

## Decisions

- QA commit subject unscoped; task subjects lowercase at derivation — PRD
  Decisions carry both; the 0001 techspec's commit-contract line is amended by
  this Spec's docs task rather than edited retroactively.
- Spec-Run logs adopt the Artifact Directory (unify, not document).
- No new config keys; no ADRs.

## Prior Spec Amendments

- 2026-07-05: This Spec amends the 0001 Implement Command techspec's commit
  contract without editing that shipped file. Per-Task commits still use the
  same type mapping and `Roundfix-Spec`/`Roundfix-Task` trailers, but the first
  rune of the derived Task subject is lowercased. QA Report commits now use
  `docs: qa report for <slug> (<verdict>)` with the `Roundfix-Spec` trailer.
