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
3. MUST carry the clause that puts each added path in the set, matching the
   existing entries' shape rather than adding an unexplained exception list.
4. MUST NOT report an ordinary source or documentation path as governed as a
   side effect of the widened predicate.
5. MUST update the recorded governed-membership characterization rows this Task
   intentionally changes and leave the rest untouched.

## Subtasks

- [ ] Identify every owned shipped authoring template the predicate misses.
- [ ] Add them under a clause consistent with the existing set entries.
- [ ] Assert monotonicity against the recorded present membership.
- [ ] Assert that a representative ordinary source path stays ungoverned.

## Acceptance Criteria

- [ ] `skills/write-prd/references/prd-template.md` and
      `skills/write-techspec/references/techspec-template.md` report governed,
      where they reported ungoverned before this Task.
- [ ] Every path recorded as governed by the characterization still reports
      governed; a test proves the set did not shrink.
- [ ] A representative ordinary source path and a representative Spec
      documentation path still report ungoverned.
- [ ] Every added path carries the clause naming why it is in the set.

## Context

- interface: `internal/speccheck/governed.go`

## Verification

- `grep -q 'func TestGovernedSetCoversOwnedShippedTemplates' internal/speccheck/governed_repocontract_test.go && go test -count=1 -tags repocontract ./internal/speccheck -run '^TestGovernedSetCoversOwnedShippedTemplates$'` — the shipped authoring templates report governed.
- `grep -q 'func TestGovernedSetOnlyGrows' internal/speccheck/governed_repocontract_test.go && go test -count=1 -tags repocontract ./internal/speccheck -run '^TestGovernedSetOnlyGrows$'` — no path recorded as governed became ungoverned, and named ordinary paths stayed ungoverned.
- `grep -q 'func TestGovernedSetOnlyGrows' internal/speccheck/governed_repocontract_test.go && go test -count=1 ./internal/speccheck` — the untagged package suite runs with the monotonicity assertion present.

## References

- `_prd.md` → User Stories 3; Goals 2; Decisions: Declared intentional breaks 2.
- `_techspec.md` → Implementation Design: Audit and compatibility; Build Order 3.
- ADR-0130.
