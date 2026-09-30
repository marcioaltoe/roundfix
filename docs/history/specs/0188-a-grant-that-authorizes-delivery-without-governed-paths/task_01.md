---
task: task_01
spec: 0188-a-grant-that-authorizes-delivery-without-governed-paths
status: completed
type: backend
complexity: medium
---

# Task 01: An explicitly empty paths list grants operations and bounds nothing

## Overview

`classifyAuthorizationRecord` in `internal/authorization/authorization.go` refuses an approved record whose `paths` has no entry. ADR-0130's contract refuses any listed path that is not governed. A Spec that changes no Governed Path therefore cannot record a delivery grant, and `deliver plan` blocked Spec 0186 on 2026-09-30 with `authorization refused: paths`. This Task makes an explicit empty YAML sequence, `paths: []`, in a Spec-contained record an operative grant of its operations that bounds no Governed Path (ADR-0179). It proves that a governed change under such a grant is still refused, and it states the rule in the delivery guides.

## Requirements

1. MUST record, in `parseAuthorizationFrontmatter`, whether the `paths` node is a YAML sequence with no items, in an unexported presence flag.
2. MUST skip the empty-`paths` refusal in `classifyAuthorizationRecord` only when the record's role is `AuthorizationRoleSpec` and that flag is set. A null `paths:`, an absent `paths` and a legacy-role record MUST keep the refusal with code `paths` and its current message. Every earlier check (status, granted, action, consuming, operations) MUST stay unchanged.
3. MUST NOT change any consumer of the record's `Paths`. An empty set bounds no Governed Path, so the Spec checker's tooling detectors, Implement preflight and the QA mechanical authorization audit refuse a governed change under it unchanged.
4. MUST keep `TestAuthorizationReaderRefusesEmptyPaths` and `TestDeliverStartRefusesASpecWithoutDeliveryAuthority` green without editing them, rename or remove no top-level test, change no exported function signature, and edit no governed file other than the two Roundfix skill files.
5. MUST put the new tests in `internal/authorization/empty_paths_test.go`, `internal/cli/deliver_plan_empty_paths_test.go` and `internal/speccheck/mechanical_empty_grant_test.go`, over temporary repositories and homes only. The Delivery Plan tests MUST drive the existing Delivery Plan workspace fixture. The audit test MUST change `Makefile` as the Governed Path.
6. MUST state, in `docs/user-guide/commands.md` and the Delivery queue section of `.agents/skills/roundfix/SKILL.md`, that a Spec with no Governed Path records an explicit empty list, in a sentence that begins with the exact words `A Spec that changes no Governed Path records` and names `paths: []`, and that this list grants the listed operations and bounds no Governed Path. It then runs `make skills-sync`.
7. MUST NOT name any decision by its `ADR-` identifier in this Task's `## Result`.

## Subtasks

- [ ] Track an explicit empty `paths` sequence and accept it for Spec-contained records.
- [ ] Prove the reader, the Delivery Plan and the mechanical audit with separate tests.
- [ ] State the rule in the guide and the skill, then sync the mirror.

## Acceptance Criteria

- [ ] An approved Spec record with `paths: []` and operations is granted with an empty `Paths`.
- [ ] A null `paths:`, an absent `paths` and a legacy record with an empty sequence are each refused on `paths`.
- [ ] `deliver plan` approves a Spec whose record declares `paths: []` and every delivery operation, `deliver start` accepts it, and the same Spec with `paths:` null is refused on `paths`.
- [ ] Under a `paths: []` grant, a Task commit changing `Makefile` reports `QA-AUTH-PATHS`, and one changing only an ordinary file reports nothing.
- [ ] The guide, the skill and its mirror carry the sentence, and `make skills-sync-check` passes.

## Context

- interface: `internal/authorization/authorization.go`
- creates: `internal/authorization/empty_paths_test.go`
- creates: `internal/cli/deliver_plan_empty_paths_test.go`
- creates: `internal/speccheck/mechanical_empty_grant_test.go`
- instruction: `internal/cli/deliver_plan_test.go`
- instruction: `docs/adr/0179-an-explicitly-empty-paths-list-grants-operations-only.md`
- interface: `docs/user-guide/commands.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestAnExplicitEmptyPathsListGrantsOperations|TestANullPathsFieldStillRefuses|TestAnAbsentPathsFieldStillRefuses|TestALegacyRecordWithEmptyPathsStillRefuses|TestDeliverPlanApprovesAnExplicitEmptyPathsGrant|TestDeliverPlanRefusesANullPathsGrant|TestDeliverStartAcceptsAnExplicitEmptyPathsGrant|TestDeliverStartRefusesASpecWithoutDeliveryAuthority|TestMechanicalAuthPathsRefusesAGovernedChangeUnderAnEmptyGrant|TestMechanicalAuthPathsAcceptsAnOrdinaryChangeUnderAnEmptyGrant|TestAuthorizationReaderRefusesEmptyPaths)$" ./internal/authorization ./internal/cli ./internal/speccheck ./internal/spec 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestAnExplicitEmptyPathsListGrantsOperations TestANullPathsFieldStillRefuses TestAnAbsentPathsFieldStillRefuses TestALegacyRecordWithEmptyPathsStillRefuses TestDeliverPlanApprovesAnExplicitEmptyPathsGrant TestDeliverPlanRefusesANullPathsGrant TestDeliverStartAcceptsAnExplicitEmptyPathsGrant TestDeliverStartRefusesASpecWithoutDeliveryAuthority TestMechanicalAuthPathsRefusesAGovernedChangeUnderAnEmptyGrant TestMechanicalAuthPathsAcceptsAnOrdinaryChangeUnderAnEmptyGrant TestAuthorizationReaderRefusesEmptyPaths; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && for file in docs/user-guide/commands.md .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'A Spec that changes no Governed Path records' || { printf 'missing phrase in %s\n' "$file" >&2; exit 1; }; done && make skills-sync-check` — expected: exit 0; before this Task none of the nine new named tests exists and the phrase is absent, so the command fails.

## References

- [_prd.md](_prd.md) — Goals 1-2; Core Features 1-3; Success Metrics 1-3
- [_techspec.md](_techspec.md) — The explicit empty list; Guides; API Contracts 1-3; Testing Approach 1-4; Build Order 1
- [references/2026-09-30-a-spec-with-no-governed-path-cannot-record-a-grant.md](references/2026-09-30-a-spec-with-no-governed-path-cannot-record-a-grant.md)
- ADR-0179; ADR-0130

## Result

Implemented presence-aware authorization parsing for an explicit empty YAML
sequence. The reader records that syntax separately from null and omission,
returns a non-nil empty `Paths`, and grants it only for a Spec-contained
record. The existing status, grant date, action, consumer and operations
validation remains ahead of the narrowed path refusal.

Removed the mechanical audit's obsolete unconditional empty-bounds finding.
Its existing changed-path loop now treats the empty bounded map normally:
every Governed Path is outside the grant, while ordinary files do not create
an authorization finding. No other bounded-path consumer changed.

Added isolated reader, Delivery Plan/start and mechanical-audit regression
suites over temporary repositories and homes. Added the exact empty-list rule
to the user guide and canonical Roundfix skill, then ran `make skills-sync` to
regenerate the mirror.

Focused-check evidence:

- Reader cases for explicit empty, null, absent and legacy records passed with
  `go test -count=1 -run` against `./internal/authorization`.
- Delivery Plan approval/refusal and Delivery start acceptance passed with
  `go test -count=1 -run` against `./internal/cli`.
- Mechanical governed-path refusal and ordinary-file acceptance passed with
  `go test -count=1 -run` against `./internal/speccheck`.
- `make skills-sync` exited zero, the required sentence is present in the
  guide, canonical skill and mirror, and `cmp` confirms both skill copies are
  byte-identical.
- `make verify-incremental` first reached the test suite but its force-stop
  integration tests could not read the sandboxed process table. The same
  command rerun with process-table permission exited zero, including format,
  vet, all Go tests, skill checks and the build.

Acceptance evidence:

- The reader regression asserts a granted explicit empty sequence, a non-nil
  empty `Paths`, and the listed `implement` operation.
- Separate negative tests preserve the exact `paths` refusal code, field and
  message for null, absent and legacy empty-sequence records.
- The Delivery Plan workspace fixture reports the explicit empty grant as
  approved, reports the null form as `authorization refused: paths`, and
  reaches the injected Delivery start owner exactly once for the grant.
- Real temporary Git histories show a `Makefile` Task commit produces one
  `QA-AUTH-PATHS` Task finding and an ordinary-file-only commit produces none.
- The guide and synchronized skills carry the required sentence verbatim.

The Daemon-owned Verification command was not rerun in this Agent turn.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/speccheck/mechanical.go`
