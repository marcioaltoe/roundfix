---
task: task_06
spec: 0132-a-grant-read-exactly-where-it-lives
status: pending
type: test
complexity: medium
---

# Task 06: Read a record wherever the archive left it

## Overview

Archiving a Spec moves it from the active root to the archive root, and it is
the last step before a Pull Request. Two expectations pin the record to its
active path and break exactly there, after every other gate has passed. This
slice makes them discover the record and makes the corpus invariant hold over
the whole corpus.

This is an authorized tooling Task. It may change only
`internal/speccheck/governed_repocontract_test.go`, this Task file, and the
ungoverned test files `internal/authorization/authorization_test.go` and
`internal/spec/spec_test.go`. Stop before any other mutation. The bounded set
comes from [_authorization.md](_authorization.md).

## Requirements

1. MUST discover a Spec's authorization record wherever it currently lives,
   trying the active root and then the archive root, rather than pinning the
   active path.
2. MUST assert the corpus invariant over active and archived Specs together,
   never requiring at least one active Spec, because the last Spec in a queue
   always empties that set.
3. MUST keep both expectations meaningful: a record that exists in neither root
   still fails, and a corpus that cannot be read still fails.
4. MUST pass against both layouts, proven by exercising a fixture Spec in the
   active root and then in the archive root.
5. MUST NOT weaken what either expectation asserts about the record's contents
   or the corpus's loadability.

## Subtasks

- [ ] Discover the record across active and archive roots.
- [ ] Assert the corpus invariant over the whole corpus.
- [ ] Exercise a fixture Spec in both layouts.
- [ ] Prove an absent record and an unreadable corpus still fail.

## Acceptance Criteria

- [ ] The record expectation passes with the Spec active and with the same Spec
      archived, exercised by moving a fixture between the two roots.
- [ ] The corpus invariant passes when no Spec is active, which fails today.
- [ ] A record present in neither root still fails the expectation.
- [ ] Neither expectation asserts less about contents or loadability than it
      does today.

## Context

- interface: `internal/spec/spec.go`
- instruction: `docs/agents/docs-layout.md`

## Verification

- `grep -q 'func TestRecordDiscoveredInEitherRoot' internal/authorization/authorization_test.go && go test -count=1 ./internal/authorization -run '^TestRecordDiscoveredInEitherRoot$'` — the record resolves in both layouts, exercised by moving a fixture.
- `grep -q 'func TestRecordDiscoveredInEitherRoot' internal/authorization/authorization_test.go || exit 1; go test -count=1 ./internal/spec -run '^TestRepositorySpecCorpusStillLoads$'` — the corpus invariant holds with no active Spec, which fails today.
- `grep -q 'func TestRecordDiscoveredInEitherRoot' internal/authorization/authorization_test.go || exit 1; go test -count=1 -tags repocontract ./internal/speccheck -run '^(TestEveryBoundedPathIsGoverned|TestCleanupHistoricalGrantEvidence)$'` — the governed contract reads the record wherever it lives.
- `grep -q 'func TestRecordDiscoveredInEitherRoot' internal/authorization/authorization_test.go || exit 1; go test -count=1 ./internal/authorization ./internal/spec` — both packages pass with the expectations moved.

## References

- `_prd.md` → Goals 3; User Stories 3; Core Features 4; Decisions: Declared intentional breaks 2.
- `_techspec.md` → Implementation Design: Expectations that survive the archive; Build Order 6.
