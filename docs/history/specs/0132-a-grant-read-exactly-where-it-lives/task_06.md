---
task: task_06
spec: 0132-a-grant-read-exactly-where-it-lives
status: completed
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

## Result

### Implementation

- The tracked Spec 0119 expectation now discovers `_authorization.md` under
  `docs/specs/<slug>/` first and `docs/history/specs/<slug>/` second before
  passing the resolved repository-relative path to the authorization reader.
- `TestRecordDiscoveredInEitherRoot` writes one exact grant, reads and checks
  it while active, moves the same fixture Spec to the History Root, and repeats
  the content and permission checks. Its negative subtest requires a lookup
  error when neither root carries the record.
- The repository corpus expectation now loads every Task Graph found across
  the active and archived roots and applies its nonempty assertion to the
  combined count. An archived-only fixture proves that an empty active root is
  valid, while an unreadable archive root remains an error.
- The governed-path contract still audits operative record contents and
  bounded paths, but locates each required Spec record under the active root
  first and the archived root second. A missing required record now fails
  instead of reaching the former zero-record skip.

### Focused checks

- Pre-change archive reproduction in a disposable repository copy:
  `go test -count=1 ./internal/authorization -run '^TestCurrentRecordPermitsItsDeclaredOperations$'`
  failed with `unreadable_record` for the pinned active path after Spec 0119
  moved to `docs/history/specs/`. The first attempt hit the sandboxed default
  Go cache; rerunning with a cache under `/private/tmp` exposed the expected
  assertion failure.
- `go test -run '^$' ./internal/authorization ./internal/spec` with a private
  Go cache passed compilation for both packages.
- `go test -tags repocontract -run '^$' ./internal/speccheck` with a private Go
  cache passed compilation for the governed contract package.
- `go test -count=1 ./internal/authorization -run '^TestCurrentRecordPermitsItsDeclaredOperations$'`
  passed against the active repository layout and passed again in the
  disposable copy after Spec 0119 moved to the History Root.
- Narrow fixture checks passed for
  `TestRecordDiscoveredInEitherRoot/moves_with_the_Spec_from_active_to_archive`
  and `TestRecordDiscoveredInEitherRoot/missing_from_both_roots`.
- Narrow corpus checks passed for
  `TestRepositorySpecCorpusStillLoads/archived-only_corpus` and
  `TestRepositorySpecCorpusStillLoads/unreadable_corpus`.
- The narrow governed-contract check
  `TestEveryBoundedPathIsGoverned/repository_records_are_governed` passed.
- The required Go 1.26.7 `rtk make verify-incremental` run first reached a
  sandbox refusal at the read-only `api.github.com` skill-sync check. The
  unchanged network-approved retry reached the repository test stage and
  exited 2 because
  `TestTaskCycleQAReportExternalProceedsWithoutStaging`,
  `TestTaskCycleSettlesCompletedWithoutCommitWhenOnlyExternalTaskFileChanged`,
  and
  `TestTaskCommitDropsSymlinkCrossingTaskFileAndCommitsRepositoryPaths`
  supply non-Git external Spec Roots that authorization resolution refuses.
  `internal/authorization`, `internal/spec`, and `internal/speccheck` passed in
  that run.

### Acceptance evidence

1. The fixture move check passed with the same exact record active and then
   archived; the tracked Spec 0119 expectation also passed in both layouts.
2. The archived-only corpus check loaded one Task Graph while the active root
   was empty.
3. The missing-record check passed only when discovery returned an error that
   named the absent Spec and both searched roots.
4. The tracked-record permission loop is unchanged. The fixture additionally
   asserts exact status, action, consuming Spec, bounded paths, operation list,
   and every `Permits` answer. Corpus loading still returns errors from active
   discovery, archive discovery, manifest and Task reads, and `Load`; the
   unreadable-root check exercised that failure path.

The commands under `## Verification` were not run; Daemon Verification owns
those commands and the terminal Task status.

### Follow-up

- The three failing `internal/daemon` external-Spec-root fixtures belong to a
  sibling Task slice and remain unchanged because `internal/daemon` is outside
  this Task's exact mutation allowlist.
