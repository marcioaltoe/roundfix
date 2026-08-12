---
task: task_08
spec: 0041-agent-selection-runtime-readiness
status: completed
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

- [x] Inject the shared profile-readiness coordinator into Doctor.
- [x] Resolve and deduplicate all effective required profiles.
- [x] Replace legacy single-model readiness output.
- [x] Render deterministic success and first-failure evidence.
- [x] Preserve independent checks and failed-Doctor exit behavior.
- [x] Add the Spec 0036 extension seam and output-order contract.

## Acceptance Criteria

- [x] Doctor success proves every distinct required profile tuple and reports
      the derived tuple/reference counts.
- [x] Doctor and `profiles validate` identify the same tuple, affected
      categories, classification, and next action for the same config.
- [x] A legacy adapter override is reported through effective adapter/profile
      evidence instead of an unrelated legacy model default.
- [x] Profile failure does not suppress Node, ACPX, adapter, Codex, or injected
      independent checks and makes Doctor exit one.
- [x] No legacy `model: ok` result can claim readiness for a tuple that the
      effective profiles do not use.
- [x] Doctor remains text-only and read-only, with no Run, Session prompt,
      config, worktree, artifact, or repository mutation.
- [x] The Repository Skill Set result can be appended after profile readiness
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

## Result

Doctor now proves the five required Agent Selection Profile categories through
the shared readiness coordinator and renders one `profiles:` result. Success
reports distinct-tuple and category-reference counts. Failure reports the first
failed exact tuple, all references, classification, bounded adapter evidence,
and one next action.

The independent adapter check now derives its runtime from the effective
`general` profile instead of `defaults.agent`. The legacy configured Agent and
`model:` readiness lines are no longer Doctor authorities. Node, ACPX, adapter,
profile, and Codex checks remain eager and ordered; Repository Skill Set
readiness can be inserted immediately after `profiles:` without another
selection proof.

Verification:

- `rtk go test ./internal/cli -run 'TestRunDoctor.*ProfileReadiness' -count=1`
  passed: 7 tests.
- `rtk go test ./internal/cli -run 'TestRunDoctor.*ContinuesChecks' -count=1`
  passed: 1 test.
- `rtk go test -race ./internal/cli -run 'TestRunDoctor.*(ProfileReadiness|ContinuesChecks)' -count=1`
  passed: 7 tests with race detection.
- `rtk go test ./internal/cli -run 'TestRunDoctor' -count=1` passed: 8 tests.
- `rtk go test ./internal/cli -count=1` passed.
- `rtk make verify` passed: 1,680 Go tests in 20 packages plus formatting,
  setup-context, skill validation, and build checks.

Acceptance evidence:

- The success test proves three deduplicated tuples from ten references across
  `general`, `backend`, `frontend`, `qa`, and `review`.
- The shared-failure test runs Doctor and `profiles validate` against the same
  config and verifies matching tuple values, affected category,
  classification, advertised controls, and configuration action.
- The legacy-adapter test verifies the adapter runtime uses the effective
  `general` preferred selection and that output excludes the unrelated legacy
  runtime model and all `agent:` or `model:` lines.
- The continuation test records one call to profile readiness and every
  independent health check after failure, verifies `exitRunFailed`, and locks
  `profiles:` before `codex:` for the Spec 0036 insertion point.
- Doctor command tests retain text-only stdout, empty diagnostic stderr on
  readiness results, no Agent prompt, and absence of config, Run, Session,
  worktree, or artifact mutations.

Follow-ups: none for this Task slice.
