---
task: task_08
spec: 0041-agent-selection-runtime-readiness
status: pending
type: backend
complexity: high
---

# Task 08: Diagnose effective Agent Selection Profiles in Doctor

## Overview

Replace Doctor's legacy configured-runtime model probe with the shared
profile-readiness coordinator. Doctor must prove all effective required
profiles and fallbacks, report one deterministic aggregate result, and continue
unrelated readiness checks when profile proof fails.

## Requirements

1. MUST resolve `general`, `backend`, `frontend`, `qa`, and `review` from the
   effective config rather than substituting `defaults.agent`.
2. MUST prove their distinct Preferred Selections and Fallback Chains through
   the same coordinator used by `profiles validate` and operational preflight.
3. MUST report one deterministic Agent Selection Profile readiness result with
   distinct-tuple and category-reference counts on success.
4. MUST report the first failed tuple, every affected category, classification,
   adapter evidence, and one next action on failure.
5. MUST continue Node, ACPX, adapter, Codex, and other independent readiness
   checks after profile failure, then return the existing failed-Doctor exit.
6. MUST remove the legacy single-model readiness authority without adding a
   second profile prover or a Doctor JSON surface.
7. MUST expose the dependency seam and output order required for Spec 0036 to
   append Repository Skill Set readiness independently.

## Subtasks

- [ ] Inject the shared profile-readiness coordinator into Doctor.
- [ ] Resolve and deduplicate all effective required profiles.
- [ ] Replace legacy single-model readiness output.
- [ ] Render deterministic success and first-failure evidence.
- [ ] Preserve independent checks and failed-Doctor exit behavior.
- [ ] Add the Spec 0036 extension seam and output-order contract.

## Acceptance Criteria

- [ ] Doctor success proves every distinct required profile tuple and reports
      the derived tuple/reference counts.
- [ ] Doctor and `profiles validate` identify the same tuple, affected
      categories, classification, and next action for the same config.
- [ ] A legacy adapter override is reported through effective adapter/profile
      evidence instead of an unrelated legacy model default.
- [ ] Profile failure does not suppress Node, ACPX, adapter, Codex, or injected
      independent checks and makes Doctor exit one.
- [ ] No legacy `model: ok` result can claim readiness for a tuple that the
      effective profiles do not use.
- [ ] Doctor remains text-only and read-only, with no Run, Session prompt,
      config, worktree, artifact, or repository mutation.
- [ ] The Repository Skill Set result can be appended after profile readiness
      without duplicating selection proof.

## Context

- instruction: `docs/specs/0036-doctor-skill-readiness/_techspec.md`
- interface: `internal/cli/doctor.go`
- interface: `internal/cli/health.go`
- interface: `internal/cli/profiles_validate.go`
- interface: `internal/cli/cli_test.go`

## Verification

- `rtk go test ./internal/cli -run 'TestRunDoctor.*ProfileReadiness' -count=1` — expected: effective categories, tuple deduplication, success counts, failure evidence, and no-mutation assertions pass.
- `rtk go test ./internal/cli -run 'TestRunDoctor.*ContinuesChecks' -count=1` — expected: every independent readiness check runs after profile failure and Doctor exits one.
- `rtk go test -race ./internal/cli -run 'TestRunDoctor.*(ProfileReadiness|ContinuesChecks)' -count=1` — expected: injected profile readiness and future skills extension seams are race-free.

## References

- `_prd.md` → User Stories 4, 7, and 9; Core Features 1, 4, 8, and 10; User
  Experience; Success Metrics.
- `_techspec.md` → Shared Readiness Coordinator; Profile-Aware Doctor
  Integration; Error Taxonomy and Diagnostics; Build Order 8.
- `../0036-doctor-skill-readiness/_techspec.md` → Repository Skill Set output
  ordering and ownership boundary.

