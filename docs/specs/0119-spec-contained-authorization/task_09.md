---
task: task_09
spec: 0119-spec-contained-authorization
status: pending
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
