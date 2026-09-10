---
task: task_09
spec: 0119-spec-contained-authorization
status: completed
type: backend
complexity: medium
---

# Task 09: Consolidate the second grant reader in the suite guard

## Overview

The suite guard owns a separate grant parser that decides which regeneration
commands are sanctioned. It scans only the legacy directory and active Spec
roots, so it cannot see an archived Spec's grant and cannot preserve a
multi-Spec consuming list. Route it through the typed reader so one record
cannot mean two things. The slice is verifiable on its own: an archived
approved grant's sanctioned regeneration becomes discoverable where it is
invisible today.

This is an authorized tooling Task. It may change only
`internal/suiteguardcontract/regeneration.go`,
`internal/suiteguardcontract/regeneration_test.go`, and this Task file. Stop
before any other mutation. The bounded set comes from the 2026-09-09 amendment
in [_authorization.md](_authorization.md).

## Requirements

1. MUST resolve grants through the typed authorization reader rather than a
   second private parser, so one record yields one verdict everywhere.
2. MUST discover operative records in archived Specs as well as in active Specs
   and the preserved legacy location, and MUST preserve a historical
   multi-Spec consuming list instead of requiring a single consuming Spec.
3. MUST keep rejecting a proposed, null-dated, malformed or unrelated record as
   a source of sanctioned regeneration, with the same outcome it reaches today.
4. MUST keep the currently resolved set of sanctioned regeneration commands and
   their outputs unchanged for every record that resolves today, so the widened
   discovery only adds records and never drops one.
5. MUST NOT introduce a second exemption list or a parallel ownership rule;
   command-only declarations keep resolving through the repository-owned
   derived-output declarations.

## Subtasks

- [ ] Route the suite guard's grant resolution through the typed reader.
- [ ] Add archived Spec discovery and multi-Spec consuming lists.
- [ ] Prove proposed and malformed records still grant nothing.
- [ ] Prove the resolved command set only grows.

## Acceptance Criteria

- [ ] An approved grant in an archived Spec contributes its sanctioned
      regeneration, where it contributes nothing today.
- [ ] A historical record naming several consuming Specs resolves with its
      actual list rather than being rejected for lacking a single consuming
      value.
- [ ] A proposed, null-dated, malformed or unrelated record contributes no
      sanctioned regeneration.
- [ ] Every command and output resolved before this Task is still resolved
      after it, proven by comparing the resolved set against the recorded
      present set rather than by inspection.
- [ ] The typed reader is the only grant parser the suite guard uses.

## Context

- interface: `internal/suiteguardcontract/regeneration.go`

## Verification

- `grep -q 'func TestSanctionedRegenerationReadsArchivedSpecGrants' internal/suiteguardcontract/regeneration_test.go && go test -count=1 ./internal/suiteguardcontract -run '^TestSanctionedRegenerationReadsArchivedSpecGrants$'` — an archived approved grant contributes its regeneration and a multi-Spec consuming list resolves.
- `grep -q 'func TestSanctionedRegenerationRejectsNonOperativeRecords' internal/suiteguardcontract/regeneration_test.go && go test -count=1 ./internal/suiteguardcontract -run '^TestSanctionedRegenerationRejectsNonOperativeRecords$'` — proposed, null-dated, malformed and unrelated records contribute nothing.
- `grep -q 'func TestSanctionedRegenerationSetOnlyGrows' internal/suiteguardcontract/regeneration_test.go && go test -count=1 ./internal/suiteguardcontract -run '^TestSanctionedRegenerationSetOnlyGrows$'` — every command and output resolved before this Task still resolves.
- `matches="$(grep -n 'yaml.Unmarshal' internal/suiteguardcontract/regeneration.go | grep -i 'grant' || true)"; test -z "$matches" || { printf '%s\n' "$matches"; exit 1; }` — the private grant parser is gone rather than left beside the typed reader.

## References

- `_prd.md` → Core Features 3; Goals 1, 4.
- `_techspec.md` → Implementation Design: Audit and compatibility; Build Order 1, 3.
- `_authorization.md` → the 2026-09-09 amendment adding these two paths.
- ADR-0149.

## Result

The suite guard now discovers candidate Markdown records in the preserved
legacy directory, active Specs, and archived Specs, then delegates grant
classification and sanctioned-regeneration parsing to the typed authorization
reader. Command-only declarations still resolve their outputs through the
repository-owned ownership tree. Explicit historical output lists retain their
sorted result shape.

Focused acceptance evidence:

1. `TestSanctionedRegenerationReadsArchivedSpecGrants` uses an approved record
   under `docs/history/specs/` whose consuming field names two Specs. Its
   command-only declaration resolves to the ownership tree's generated output.
2. The same archived test proves that the enclosing archived Spec can appear
   after another consumer in the declared list and still contributes its
   sanctioned regeneration.
3. `TestSanctionedRegenerationRejectsNonOperativeRecords` independently checks
   proposed, null-dated, malformed, and wrong-consumer records; every case
   resolves to no sanctioned regeneration.
4. `TestSanctionedRegenerationSetOnlyGrows` compares the widened result with a
   recorded pre-change subset containing an active Spec record and a dated
   prose-era legacy record. Both commands and their outputs remain present,
   including the historical sorted order for an explicit multi-output list.
5. `ReadSanctionedRegenerations` now calls `spec.ReadAuthorization`; the private
   `approvedSpecGrant` parser was removed. `ParseSanctionedRegenerations`
   remains exported only for the existing changed-path and ownership audits;
   the suite guard's grant-discovery path no longer calls it.

Focused checks:

- Before the production edit,
  `rtk env GOCACHE=/private/tmp/roundfix-task09-gocache go test -count=1 ./internal/suiteguardcontract -run 'SanctionedRegeneration'`
  failed because the archived multi-Spec record contributed no declaration.
- After the last implementation edit,
  `rtk env GOCACHE=/private/tmp/roundfix-task09-gocache go test -count=1 ./internal/suiteguardcontract -run '^(TestCleanupRegenerationDiscovery|TestSanctionedRegeneration.*)$'`
  passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task09-gocache go build -buildvcs=false ./...`
  passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task09-gocache go test -count=1 ./internal/spec -run '^TestAuthorizationReader(ClassifiesGrantState|ResolvesPreservedHistoricalRecords)$'`
  failed during package setup with the test-only import cycle described below.
- `rtk env GOCACHE=/private/tmp/roundfix-task09-gocache go test -count=1 ./internal/suiteguard -run '^TestSanctionedRegeneration'`
  failed because the positive consumer fixture still supplies the malformed
  legacy record described below.

Follow-up constraints found outside this Task's mutation allowlist:

- `internal/spec/main_test.go` imports `internal/suiteguard`, whose production
  path reaches `internal/suiteguardcontract`. The new typed-reader dependency
  closes that test-only cycle, so the focused `internal/spec` test command
  fails during setup with `import cycle not allowed in test`. Repairing the
  test package boundary requires authority for a path outside this Task.
- `internal/suiteguard/suiteguard_test.go` still creates
  `docs/workflow/authorizations/fixture.md` without the date, consuming work,
  or bounded paths required by the typed legacy contract. Its positive
  sanctioned-regeneration test now refuses that malformed fixture. Updating
  the fixture to an operative dated legacy record requires authority for that
  out-of-scope file.

The Daemon-owned commands under `## Verification` were not run in this Agent
turn.
