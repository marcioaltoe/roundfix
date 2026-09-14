---
task: task_11
spec: 0119-spec-contained-authorization
status: completed
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

## Result

The typed authorization reader now lives in `internal/authorization`, below
both packages that consume it. `internal/spec` preserves its existing reader
API through type aliases and forwarding functions, while
`internal/suiteguardcontract` imports the shared implementation directly. The
authorization frontmatter decoder and grant classifier exist only in the new
reader package.

Legacy regeneration discovery now consumes the shared reader's parsed
regeneration projection without requiring the typed status, grant date,
action, consuming Spec, or bounded paths that Spec-contained grants require.
Unreadable records still return an error. Active and archived Spec records
still contribute only when the shared reader classifies them as granted.

Acceptance evidence:

1. Before the change, `rtk go test ./internal/spec -run '^$'` failed with the
   import chain `internal/spec` test → `internal/suiteguard` →
   `internal/suiteguardcontract` → `internal/spec`. After the move,
   `rtk go test -count=1 ./internal/spec` passed. A fresh
   `rtk rg -n "roundfix/internal/spec" internal/suiteguardcontract --glob '*.go'`
   returned no matches.
2. The new
   `TestSanctionedRegenerationResolvesLegacyRecordsWithoutFrontmatter` writes a
   legacy record containing only its heading and YAML declaration. Before the
   production change, the focused test failed because discovery returned no
   declarations. After the change, `rtk go test -count=1
   ./internal/suiteguardcontract` passed and preserved both outputs in sorted
   order.
3. That suiteguardcontract run also exercised proposed, null-dated, malformed,
   and unrelated active Spec records as non-operative, plus the approved
   archived multi-Spec grant as operative. `rtk go test -count=1
   ./internal/spec -run 'TestAuthorizationReader'` separately passed all typed
   reader grant, refusal, operation, path, and historical-record cases.
4. `TestSanctionedRegenerationSetOnlyGrows` now creates its own minimal legacy
   record instead of depending on repository contents. As a mutation check, I
   temporarily removed the legacy append and ran `rtk go test -count=1
   ./internal/suiteguardcontract -run 'TestSanctionedRegenerationSetOnlyGrows'`;
   it failed at `regeneration_test.go:302` for the missing recorded legacy
   declaration. I restored the implementation, and the full focused package
   run passed.
5. A fresh search for `parseAuthorizationRecord`,
   `splitAuthorizationFrontmatter`, and `authorizationFrontmatterMapping`
   found all three only in `internal/authorization/authorization.go`. The
   `internal/spec` facade contains no YAML decoder, and the suiteguardcontract
   search above confirms it no longer imports `internal/spec`.

Focused checks:

- `rtk go test -count=1 ./internal/suiteguardcontract` — passed.
- `rtk go test -count=1 ./internal/spec` — passed.
- `rtk go test -count=1 ./internal/suiteguard -run
  'TestSanctionedRegeneration(IsDeclaredInProcess|IsNotAViolationWrongCommandIsRefused|IsNotAViolationUndeclaredCommandIsRefused)'`
  — passed all three in-process declaration cases.
- `rtk go test ./internal/speccheck -run '^$'` — package and shared-reader
  imports compiled without running tests.
- `rtk go test ./internal/authorization` — package compiled; it has no direct
  test files because the preserved `internal/spec` compatibility suite owns
  the reader contract.
- `rtk git diff --check` — passed.

The Task's declared `## Verification` commands were not run; Daemon
Verification remains pending.
