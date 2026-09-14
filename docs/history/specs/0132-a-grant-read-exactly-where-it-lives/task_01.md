---
task: task_01
spec: 0132-a-grant-read-exactly-where-it-lives
status: completed
type: test
complexity: medium
---

# Task 01: Characterize today's parse, resolution and discovery answers

## Overview

Record what the authorization reader answers today for a malformed closing
delimiter, for a Spec Root outside the code repository, and for a record whose
filename uses the naming already in use. The slice is verifiable alone: it
states the present contract, including the three answers this Spec intends to
move, so the diff shows the moves rather than hiding them.

This is an authorized tooling Task. It may change only
`internal/speccheck/constraints_characterization_test.go` and this Task file,
plus ungoverned test files under `internal/authorization` and
`internal/suiteguardcontract`. Stop before any other mutation. The bounded set
comes from [_authorization.md](_authorization.md).

## Requirements

1. MUST record that a record whose closing frontmatter line carries extra
   characters resolves to a grant today, which is the answer Task 02 moves.
2. MUST record what the operation resolver and the citation resolver answer for
   a Spec Root outside the code repository today, which is the answer Tasks 03
   and 04 move.
3. MUST record that discovery ignores a record whose filename contains
   `authorized` rather than `authorization`, which is the answer Task 05 moves.
4. MUST record the answers that must not move: a well-formed record under the
   default root grants exactly what it grants today.
5. MUST NOT change reader behavior, and MUST NOT weaken or delete an existing
   assertion.

## Subtasks

- [ ] Record the malformed-delimiter answer.
- [ ] Record the external Spec Root answers at both resolvers.
- [ ] Record the discovery answer for the naming in use.
- [ ] Record the default-root answers that must not move.

## Acceptance Criteria

- [ ] A case asserts that a closing line carrying extra characters resolves to a
      grant today, marked as the behavior Task 02 changes.
- [ ] A case asserts each resolver's current answer for an external Spec Root,
      marked as the behavior Tasks 03 and 04 change.
- [ ] A case asserts that a record named with `authorized` is ignored today.
- [ ] A case asserts a well-formed default-root record's grant, unmarked,
      because it must not move.
- [ ] Every assertion that existed in the touched files before this Task still
      runs and passes.

## Context

- interface: `internal/authorization/authorization.go`
- interface: `internal/spec/authorization.go`

## Verification

- `grep -q 'func TestAuthorizationParseCharacterization' internal/authorization/authorization_test.go && go test -count=1 ./internal/authorization -run '^TestAuthorizationParseCharacterization$'` — the malformed-delimiter and default-root answers are recorded; absence cannot pass.
- `grep -q 'func TestExternalSpecRootResolutionCharacterization' internal/speccheck/constraints_characterization_test.go && go test -count=1 ./internal/speccheck -run '^TestExternalSpecRootResolutionCharacterization$'` — both resolvers' current external-root answers are recorded.
- `grep -q 'func TestDiscoveryNamingCharacterization' internal/suiteguardcontract/regeneration_test.go && go test -count=1 ./internal/suiteguardcontract -run '^TestDiscoveryNamingCharacterization$'` — the ignored-naming answer is recorded.
- `grep -q 'func TestAuthorizationParseCharacterization' internal/authorization/authorization_test.go || exit 1; go test -count=1 ./internal/authorization ./internal/suiteguardcontract` — the pre-existing assertions in both packages still pass beside the new characterization.

## References

- `_prd.md` → Decisions: Regression locks; Success Metrics.
- `_techspec.md` → Testing Approach observation 5; Build Order 1.

## Result

Implementation evidence:

- `TestAuthorizationParseCharacterization` records that `---evil` closes the
  frontmatter and grants today, with a comment naming Task 02 as the intended
  change. Its unmarked well-formed case records the same exact approved record,
  bounded path, and sole permitted operation under the default root.
- `TestExternalSpecRootResolutionCharacterization` places the Spec Root in a
  separate temporary directory. It records the operation resolver's current
  `unresolved` answer at the hard-coded default record path and the citation
  resolver's current exact-one-record refusal, with comments naming Tasks 03
  and 04 as the intended changes.
- `TestDiscoveryNamingCharacterization` records that the existing
  `2026-09-08-authorized-qa-archive-override.md` naming is ignored today, with
  a comment naming Task 05 as the intended change.
- No reader or discovery implementation changed; the diff adds assertions
  only.

Focused-check evidence:

- Pre-change
  `rtk rg -n 'func Test(AuthorizationParseCharacterization|ExternalSpecRootResolutionCharacterization|DiscoveryNamingCharacterization)' internal/authorization/authorization_test.go internal/speccheck/constraints_characterization_test.go internal/suiteguardcontract/regeneration_test.go`
  exited 1 with no matches, proving the named characterization cases were
  absent.
- `rtk go test -count=1 ./internal/authorization ./internal/speccheck ./internal/suiteguardcontract -run 'Characterization'`
  initially reached no test binary because the sandbox denied the standard Go
  build cache. The unchanged retry with cache access passed 10 tests across all
  three packages.
- `rtk go test -count=1 ./internal/authorization ./internal/speccheck ./internal/suiteguardcontract`
  passed 371 tests across all three packages, including every pre-existing
  assertion in the touched test files.
- `rtk gofmt -d` over the three test files and `rtk git diff --check` both
  exited 0 with no output.
- The changed-file postflight lists only this Task file and the three authorized
  characterization test files.
- The Task's declared `## Verification` commands were not run; Daemon
  Verification owns them after this handoff.
