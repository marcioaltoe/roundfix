---
task: task_02
spec: 0119-spec-contained-authorization
status: pending
type: backend
complexity: high
---

# Task 02: Read a typed grant by role and keep legacy records resolving

## Overview

Add the authorization reader the rest of the Spec consumes: one parser that
answers whether a record is an operative grant, for whom, over which exact
paths, and under which limits. It is verifiable on its own because a grant, a
proposal, and a malformed record produce three distinguishable answers from the
same input path.

## Requirements

1. MUST parse an authorization record by its declared role into approval state,
   grant date, action, consuming Spec, exact bounded paths, and any sanctioned
   regeneration the record declares.
2. MUST treat a proposal as inspectable and non-granting. An unknown state, an
   absent or unparseable grant date, an empty path list, a consuming field that
   does not name the asking Spec, and a contradictory record are refusals with
   distinguishable reasons, never a grant.
3. MUST refuse a path that escapes the repository, is absolute, carries a glob,
   traverses upward, resolves through a symlink, or repeats another entry.
4. MUST NOT treat a Spec name appearing anywhere in prose as a consuming
   relation; only the declared consuming field creates one.
5. MUST resolve the existing explicit legacy record formats through a bounded
   adapter, preserving a multi-Spec historical grant's actual consuming list and
   scope, without rewriting any historical record.
6. MUST report an unreadable record or an unavailable revision as unresolved
   rather than as either granted or refused, so failure is closed on evidence
   and never on doubt.
7. MUST NOT change the behavior of any existing caller in this Task; the reader
   is introduced and proven, and its consumers are wired in later Tasks.

## Subtasks

- [ ] Define the parsed grant shape and its refusal reasons.
- [ ] Parse the typed record and classify approval state.
- [ ] Validate bounded paths against escape, glob, symlink, and duplication.
- [ ] Adapt the explicit legacy formats, preserving multi-Spec consuming lists.
- [ ] Prove the reader against the preserved historical corpus.

## Acceptance Criteria

- [ ] An approved record for the asking Spec resolves to a grant carrying its
      exact bounded paths; the same record with `status` proposed or a null
      grant date resolves to a refusal naming which field withheld it.
- [ ] A record naming a different consuming Spec refuses even when the asking
      Spec's slug appears in its prose.
- [ ] Each of an absolute path, an upward-traversing path, a glob, a symlinked
      path, and a duplicate entry refuses with a reason naming that path.
- [ ] Every authorization record preserved at Git revision
      `6b8ea48725cbca13974eee0b400b3482202874f6` under
      `docs/workflow/authorizations/` is read and classified without error, and
      none of them silently becomes an approval for this Spec. These records
      were written between 2026-07-30 and 2026-09-08 by the pre-Spec workflow
      that this Spec replaces; they are outside evidence because neither their
      content nor their shape was authored by this Spec, so a reader that only
      fits the new schema fails against them.
- [ ] An unreadable record and an unavailable revision each report unresolved,
      distinguishable from both granted and refused.

## Context

- interface: `internal/spec/spec.go`
- interface: `internal/suiteguardcontract/regeneration.go`

## Verification

- `grep -q 'func TestAuthorizationReaderClassifiesGrantState' internal/spec/authorization_test.go && go test -count=1 ./internal/spec -run '^TestAuthorizationReaderClassifiesGrantState$'` — approved, proposed, null-dated, wrong-consuming and contradictory records resolve to distinguishable answers.
- `grep -q 'func TestAuthorizationReaderRefusesEscapingPaths' internal/spec/authorization_test.go && go test -count=1 ./internal/spec -run '^TestAuthorizationReaderRefusesEscapingPaths$'` — absolute, traversing, globbed, symlinked and duplicate paths each refuse.
- `grep -q 'func TestAuthorizationReaderResolvesPreservedHistoricalRecords' internal/spec/authorization_test.go && go test -count=1 ./internal/spec -run '^TestAuthorizationReaderResolvesPreservedHistoricalRecords$'` — the preserved historical corpus is read and classified, with an unavailable revision reported explicitly rather than skipped.
- `test -f internal/spec/authorization.go && go build -buildvcs=false ./...` — the reader exists and the tree compiles with it.

## References

- `_prd.md` → User Stories 1; Core Features 1, 3, 5; Goals 1, 2, 4.
- `_techspec.md` → Implementation Design: Operative record and source identity; Data Models; Build Order 1.
- ADR-0057, ADR-0130, ADR-0149.
