---
task: task_04
spec: 0119-spec-contained-authorization
status: pending
type: backend
complexity: low
---

# Task 04: Extend the governed set to the owned shipped templates

## Overview

The historical Governed Path set must stay monotonic, and the owned shipped
authoring templates that an authorization has bounded are currently missed by
the narrow path predicate. Add them, so a change to a shipped template is
audited against a grant instead of passing as ordinary source. The slice is
verifiable on its own and touches nothing else.

This is an authorized tooling Task. It may change only
`internal/speccheck/governed.go`,
`internal/speccheck/governed_repocontract_test.go`, and this Task file. Stop
before any other mutation. The bounded set comes from the approved grant in
[_authorization.md](_authorization.md).

## Requirements

1. MUST report the owned shipped authoring templates as governed, including
   `skills/write-prd/references/prd-template.md` and
   `skills/write-techspec/references/techspec-template.md`.
2. MUST keep every path the predicate reports as governed today still governed,
   so the set only grows.
3. MUST report as governed every path bounded by a record the widened discovery
   reaches, including the archived Spec 0130 grant's paths, so the contract and
   the predicate cannot disagree about a record that already exists.
3. MUST carry the clause that puts each added path in the set, matching the
   existing entries' shape rather than adding an unexplained exception list.
4. MUST NOT report an ordinary source or documentation path as governed as a
   side effect of the widened predicate.
5. MUST make the bounded-path contract discover operative Spec-contained grants
   rather than only the legacy directory, so it audits real records instead of
   skipping. The legacy directory was removed, so the contract skips today and
   audits nothing, including the grant this Spec consumes.
6. MUST keep the contract's skip meaningful: it may skip only when no operative
   record exists anywhere, never when an approved Spec-contained grant is
   present.
7. MUST update the recorded governed-membership characterization rows this Task
   intentionally changes and leave the rest untouched.

## Subtasks

- [ ] Identify every owned shipped authoring template the predicate misses.
- [ ] Add them under a clause consistent with the existing set entries.
- [ ] Assert monotonicity against the recorded present membership.
- [ ] Assert that a representative ordinary source path stays ungoverned.
- [ ] Make the bounded-path contract read operative Spec-contained grants.

## Acceptance Criteria

- [ ] `skills/write-prd/references/prd-template.md` and
      `skills/write-techspec/references/techspec-template.md` report governed,
      where they reported ungoverned before this Task.
- [ ] Every path recorded as governed by the characterization still reports
      governed; a test proves the set did not shrink.
- [ ] A representative ordinary source path and a representative Spec
      documentation path still report ungoverned.
- [ ] Every added path carries the clause naming why it is in the set.
- [ ] The bounded-path contract's repository-records case runs and passes
      against the approved Spec-contained grant, where it skips today because
      it reads only the removed legacy directory.
- [ ] Every path bounded by this Spec's approved grant reports governed, so the
      contract that consumes the predicate agrees with it.
- [ ] Every path bounded by every record the widened discovery reaches reports
      governed, including the archived Spec 0130 grant's paths; the acceptance
      test enumerates the discovered records rather than a hand-picked sample.

## Context

- interface: `internal/speccheck/governed.go`

## Verification

- `grep -q 'func TestGovernedSetCoversOwnedShippedTemplates' internal/speccheck/governed_repocontract_test.go && go test -count=1 -tags repocontract ./internal/speccheck -run '^TestGovernedSetCoversOwnedShippedTemplates$'` — the shipped authoring templates report governed.
- `grep -q 'func TestGovernedSetOnlyGrows' internal/speccheck/governed_repocontract_test.go && go test -count=1 -tags repocontract ./internal/speccheck -run '^TestGovernedSetOnlyGrows$'` — no path recorded as governed became ungoverned, and named ordinary paths stayed ungoverned.
- `grep -q 'func TestGovernedSetOnlyGrows' internal/speccheck/governed_repocontract_test.go || exit 1; go test -count=1 -tags repocontract ./internal/speccheck -run '^(TestCleanupHistoricalGrantEvidence|TestEveryBoundedPathIsGoverned)$'` — the pre-existing governed contracts still run and pass under the widened predicate.
- `log="$(mktemp)"; go test -count=1 -v -tags repocontract ./internal/speccheck -run '^TestEveryBoundedPathIsGoverned$' > "$log" 2>&1 || { cat "$log"; exit 1; }; grep -q -- '--- PASS: TestEveryBoundedPathIsGoverned/repository_records_are_governed' "$log" || { cat "$log"; exit 1; }` — the bounded-path contract audits real records instead of skipping, which it does today.

## References

- `_prd.md` → User Stories 3; Goals 2; Decisions: Declared intentional breaks 2.
- `_techspec.md` → Implementation Design: Audit and compatibility; Build Order 3.
- ADR-0130.
