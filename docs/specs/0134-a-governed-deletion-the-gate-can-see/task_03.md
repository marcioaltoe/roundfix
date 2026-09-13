---
task: task_03
spec: 0134-a-governed-deletion-the-gate-can-see
status: pending
type: backend
complexity: medium
---

# Task 03: Let an enumerated regeneration list be the whole list

## Overview

A record that enumerates its regeneration `outputs` has narrowed its own grant,
and the audit unions that narrowing with every output the repository's ownership
declarations derive for the same command — restoring exactly what the record
excluded. This slice resolves ownership only for a command-only declaration. It
is verifiable alone: an output outside an enumerated list refuses, and a
command-only record still resolves ownership.

## Requirements

1. MUST treat an enumerated `outputs` list as the allowed set for that command,
   and MUST NOT add repository-owned outputs to it.
2. MUST resolve repository-owned outputs when, and only when, the record
   supplies no list, which is the division ADR-0149 draws between the grant and
   the ownership tree.
3. MUST match the suite guard, which already treats an absent list as the signal
   to resolve ownership and a present list as final, so the two readers agree.
4. MUST keep a command-only declaration resolving exactly the outputs it
   resolves today.
5. MUST keep every changed-path refusal that works today working with its
   existing token.

## Subtasks

- [ ] Resolve ownership only for a command-only declaration.
- [ ] Prove an output outside an enumerated list refuses.
- [ ] Prove a command-only record is unchanged.
- [ ] Prove the audit and the suite guard agree.

## Acceptance Criteria

- [ ] A record enumerating one output refuses a consuming commit that changes a
      different owner-derived output for the same command; it passes today.
- [ ] A record enumerating an output still passes a consuming commit that
      changes exactly that output.
- [ ] A command-only record resolves the same outputs it resolves today,
      compared against the repository's own ownership declarations.
- [ ] The audit and the suite guard return the same allowed set for the same
      record, both for an enumerated list and for a command-only declaration.

## Context

- interface: `internal/speccheck/mechanical.go`
- interface: `internal/suiteguardcontract/regeneration.go`

## Verification

- `grep -q 'func TestEnumeratedOutputsAreAuthoritative' internal/speccheck/mechanical_test.go && go test -count=1 ./internal/speccheck -run '^TestEnumeratedOutputsAreAuthoritative$'` — an owner-derived output outside an enumerated list refuses; this fails today.
- `grep -q 'func TestCommandOnlyDeclarationStillResolvesOwnership' internal/speccheck/mechanical_test.go && go test -count=1 ./internal/speccheck -run '^TestCommandOnlyDeclarationStillResolvesOwnership$'` — a command-only record resolves the outputs it resolves today.
- `grep -q 'func TestAuditAndSuiteGuardAgreeOnAllowedOutputs' internal/speccheck/mechanical_test.go && go test -count=1 ./internal/speccheck -run '^TestAuditAndSuiteGuardAgreeOnAllowedOutputs$'` — both readers return the same allowed set for the same record.
- `grep -q 'func TestEnumeratedOutputsAreAuthoritative' internal/speccheck/mechanical_test.go || exit 1; go test -count=1 ./internal/speccheck ./internal/suiteguardcontract` — both packages pass with the division restored.

## References

- `_prd.md` → Goals 2; Core Features 3-4; Declared intentional breaks 2.
- `_techspec.md` → Implementation Design: An enumerated list is the whole list; Build Order 3.
- ADR-0149.
