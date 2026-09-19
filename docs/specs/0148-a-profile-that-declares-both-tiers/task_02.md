---
task: task_02
spec: 0148-a-profile-that-declares-both-tiers
status: pending
type: chore
complexity: low
---

# Task 02: Regenerate the derived catalog and plan

## Overview

A Profile with one more decision is a different Profile, so the catalog digest,
the normalized catalog and the plan characterizations move with it. This slice
runs the sanctioned regeneration and commits its output.

## Requirements

1. MUST regenerate the derived catalog and plan artifacts through the sanctioned
   command, never by hand.
2. MUST leave the generated guides, the guide template, the template index, this
   repository's decision record and the formatter golden untouched: this Spec
   publishes nothing to them.
3. MUST leave the repository Verification green.
4. MUST NOT change any path outside this Spec's bounded file, the derived
   artifacts the command rewrites, and its own Task file.

## Subtasks

- [ ] Run the sanctioned regeneration.
- [ ] Confirm only the derived catalog and plan artifacts moved.
- [ ] Confirm the gate passes.

## Acceptance Criteria

- [ ] The derived catalog and plan artifacts match what the command produces.
- [ ] No generated guide, template, index, decision record or formatter golden
      changed.
- [ ] `make verify` exits 0.

## Context

- interface: `internal/baseline/testdata/catalog.digest`

## Verification

- `derived="$(git diff --name-only HEAD~1 HEAD -- internal/baseline/testdata)"; test -n "$derived" || exit 1; guides="$(git diff --name-only HEAD~1 HEAD -- docs/agents internal/baseline/assets/templates internal/baseline/assets/formatter-fixtures)"; test -z "$guides"` — expected: exit 0; this Task's own commit moved the derived artifacts and published nothing to the guides, their templates or the formatter goldens. Before this Task no derived artifact has moved, so the command fails.
- `go test -count=1 -run "^TestCatalogCompatibility$|^TestBaselinePlanCharacterization$" ./internal/baseline` — expected: exit 0; the derived artifacts agree with the Profile that moved them. Before this Task they disagree, because Task 01 added a decision without regenerating them.

## References

`_prd.md` → Core Feature 4; Declared intentional breaks; Regression locks;
`_techspec.md` → Implementation Design: The derived artifacts; Build Order 2;
`_authorization.md`.
