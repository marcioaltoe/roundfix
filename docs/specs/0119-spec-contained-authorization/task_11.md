---
task: task_11
spec: 0119-spec-contained-authorization
status: pending
type: backend
complexity: medium
---

# Task 11: Restore legacy declarations and break the reader's import cycle

## Overview

Task 09 routed the suite guard's grant resolution through the reader in
`internal/spec`, which created an import cycle and silently dropped legacy
records that carry no frontmatter. Both defects reached the repository suite
rather than the Task's own gate. This slice repairs them together, because they
share one cause: the reader was placed where its consumers cannot reach it, and
its widened validation was applied to records that never carried the fields it
now demands.

The governed portion of this slice is exactly
`internal/suiteguardcontract/regeneration.go` and
`internal/suiteguardcontract/regeneration_test.go`, both already bounded by the
approved grant. Ordinary source changes outside the governed set stay within
this Spec's declared behavior.

## Requirements

1. MUST remove the import cycle. `internal/suiteguardcontract` must not import
   `internal/spec`, because `internal/spec`'s own test imports
   `internal/suiteguard`, which imports `internal/suiteguardcontract`. The
   typed reader belongs in a package both consumers can import without either
   depending on the other.
2. MUST keep one parsing implementation. Repairing the cycle by duplicating the
   grant parser reintroduces the disagreement Task 09 existed to remove; both
   consumers read the same code.
3. MUST restore every legacy record that resolved before Task 09. A record
   under the preserved legacy location that carries a `Sanctioned regeneration`
   block and no frontmatter contributes its declared command and outputs
   exactly as it did, because the typed grant fields were never required of
   that form.
4. MUST keep the Spec-contained path validating as Task 09 delivered it: a
   proposed, null-dated, malformed or unrelated Spec record still contributes
   nothing, and an archived Spec's approved grant still contributes.
5. MUST cover the legacy path in the regression gate. The existing set-only-
   grows assertion measures the real repository, whose legacy directory no
   longer exists, so that path had no coverage and the regression passed the
   Task that caused it. The gate must exercise a legacy record directly.
6. MUST leave the repository suite passing, including the suite guard's
   in-process sanctioned-declaration fixtures, which are the only place the
   legacy declaration path is exercised end to end.

## Subtasks

- [ ] Move the shared reader where both consumers import it without a cycle.
- [ ] Point `internal/spec` and `internal/suiteguardcontract` at that one reader.
- [ ] Restore frontmatter-free legacy record resolution.
- [ ] Extend the regression gate to exercise a legacy record directly.

## Acceptance Criteria

- [ ] `go vet ./internal/spec` reports no import cycle, and no package under
      `internal/suiteguardcontract` imports `internal/spec`.
- [ ] A legacy record with a `Sanctioned regeneration` block and no frontmatter
      contributes its command and outputs; the same record contributed nothing
      before this Task.
- [ ] A proposed, null-dated, malformed or unrelated Spec-contained record still
      contributes nothing, and an archived Spec's approved grant still
      contributes.
- [ ] The set-only-grows assertion fails when legacy resolution is removed,
      proven by exercising a legacy record rather than by scanning a repository
      that no longer contains one.
- [ ] Exactly one grant parsing implementation exists; a search for a second
      frontmatter decoder over authorization records finds none.

## Context

- interface: `internal/suiteguardcontract/regeneration.go`
- interface: `internal/suiteguard/suiteguard_test.go`

## Verification

- `go vet ./internal/spec ./internal/suiteguardcontract ./internal/suiteguard` — the import cycle is gone; this fails today.
- `matches="$(grep -n 'roundfix/internal/spec' internal/suiteguardcontract/*.go || true)"; test -z "$matches" || { printf '%s\n' "$matches"; exit 1; }` — the low-level contract package no longer depends on the package whose test depends on it.
- `grep -q 'func TestSanctionedRegenerationResolvesLegacyRecordsWithoutFrontmatter' internal/suiteguardcontract/regeneration_test.go && go test -count=1 ./internal/suiteguardcontract -run '^TestSanctionedRegenerationResolvesLegacyRecordsWithoutFrontmatter$'` — a frontmatter-free legacy record contributes its declaration again.
- `go test -count=1 ./internal/suiteguard -run '^TestSanctionedRegeneration'` — the in-process declaration fixtures pass, including the one this regression broke.
- `go test -count=1 ./internal/spec ./internal/speccheck ./internal/suiteguardcontract` — the three packages that share the reader all build and pass together.

## References

- `_prd.md` → Core Features 3; Goals 4; Decisions: Declared intentional breaks.
- `_techspec.md` → Implementation Design: Audit and compatibility; Risks & Considerations.
- `_authorization.md` → the 2026-09-09 amendment bounding the two suiteguardcontract paths.
- ADR-0149.
