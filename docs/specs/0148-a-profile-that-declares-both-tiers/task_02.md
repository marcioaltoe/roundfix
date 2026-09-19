---
task: task_02
spec: 0148-a-profile-that-declares-both-tiers
status: completed
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

- [x] Run the sanctioned regeneration.
- [x] Confirm only the derived catalog and plan artifacts moved.
- [ ] Confirm the gate passes.

## Acceptance Criteria

- [ ] The derived catalog and plan artifacts match what the command produces.
- [ ] No generated guide, template, index, decision record or formatter golden
      changed.
- [ ] `make verify` exits 0.

## Context

- interface: `internal/baseline/testdata/catalog.digest`

## Verification

- `derived="$(git diff --name-only -- internal/baseline/testdata)"; test -n "$derived" || exit 1; guides="$(git diff --name-only -- docs/agents internal/baseline/assets/templates internal/baseline/assets/formatter-fixtures)"; test -z "$guides"` — expected: exit 0; the regeneration moved the derived artifacts in the working tree and published nothing to the guides, their templates or the formatter goldens. Verification reads the working tree before the Daemon commits, so a commit-relative range cannot observe this Task at all; before the regeneration the tree is clean and the command fails.
- `go test -count=1 -run "^TestCatalogCompatibility$|^TestBaselinePlanCharacterization$" ./internal/baseline` — expected: exit 0; the derived artifacts agree with the Profile that moved them. Before this Task they disagree, because Task 01 added a decision without regenerating them.

## References

`_prd.md` → Core Feature 4; Declared intentional breaks; Regression locks;
`_techspec.md` → Implementation Design: The derived artifacts; Build Order 2;
`_authorization.md`.

## Result

- Ran the sanctioned `make baseline-digests` command successfully. It regenerated
  `internal/baseline/testdata/catalog.digest`,
  `internal/baseline/testdata/catalog.normalized.json`, and the four plan
  characterization goldens under
  `internal/baseline/testdata/plan-characterization/`.
- The generator's embedded focused checks passed, including the baseline package
  characterization/regeneration checks and strict catalog validation; it ended
  with `{"type":"baseline-digests","ok":true,"changed":true}`.
- The post-regeneration worktree contains only this task file and the six
  derived artifacts above. No generated guide, guide template, template index,
  decision record, or formatter golden changed.
- Acceptance evidence: the derived catalog and plan artifacts were produced by
  the sanctioned command; the publication paths remained untouched. The
  repository `make verify` gate was not run in this Daemon-assigned turn and
  remains for Daemon Verification.
