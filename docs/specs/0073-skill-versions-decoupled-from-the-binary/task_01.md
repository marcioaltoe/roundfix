---
task: task_01
spec: 0073-skill-versions-decoupled-from-the-binary
status: pending
type: chore
complexity: medium
---

# Task 01: Let every owned skill declare its version

## Overview

Owned skill frontmatter carries `name`, `description`, and `metadata` and
nothing else. Because they are authored in this repository, adding a version is
this Spec's own work rather than a dependency on anyone.

This slice is deliberately inert: skills gain a field nothing reads yet. The
repository stays green while the identity exists for task_02 to compare
against. Comparing before declaring would leave every owned skill
`unversioned` at once, and a transition indistinguishable from a regression is
worse than no transition.

## Requirements

1. MUST add a declared version to the frontmatter of every member of
   `OWNED_SKILLS`, canonical and mirror.
2. MUST take the member set from the `Makefile` variable rather than from any
   list copied into a Spec artifact. This Spec exists because a pin drifted
   from what it pinned; a copied member list would repeat that in a new place.
3. MUST add a check that every `OWNED_SKILLS` member declares a version, so a
   skill added later cannot skip it.
4. MUST NOT add a version to any third-party skill, and MUST NOT fail one for
   lacking it.
5. MUST NOT change any behaviour that reads or compares versions — nothing
   reads this field until task_02.
6. MUST regenerate both mirrors with `make skills-sync` and run the ADR-0081
   regeneration chain, including the two characterization corpora that
   `make baseline-digests` does not reach:

   ```bash
   go test ./internal/baseline -count=1 \
     -run TestBaselinePlanCharacterization -update-baseline-plan-characterization
   go test ./internal/baseline -count=1 \
     -run TestCatalogDiagnosticCharacterization -update-catalog-diagnostics
   ```

## Subtasks

- [ ] Add the version field to every owned skill and its mirror.
- [ ] Add the every-member-declares-one check.
- [ ] Run `make skills-sync` and the regeneration chain.

## Acceptance Criteria

- [ ] Every `OWNED_SKILLS` member declares a version, canonical and mirror.
- [ ] The member set is read from the `Makefile` variable, asserted by the
      check failing when a member is added without a version.
- [ ] No third-party skill gained a version and none is failed for lacking one.
- [ ] Every mirror is byte-identical to its canonical skill.
- [ ] `make verify` exits 0.
- [ ] No Go behaviour reads the new field yet.

## Context

- instruction: `.agents/skills/roundfix/SKILL.md`
- interface: `Makefile`

## Verification

- `make skills-sync-check` — expected: exit 0; every mirror matches.
- `go run -buildvcs=false ./cmd/roundfix skills check` — expected: exit 0.
- `make verify` — expected: exit 0.
- `output="$(go test ./skills -count=1 -v 2>&1)"; status=$?; printf '%s\n' "$output"; [ "$status" -eq 0 ] || exit "$status"; printf '%s\n' "$output" | grep -q -- "--- PASS"`
  — expected: exit 0; the skill contract tests ran and passed, with the exit
  status preserved rather than hidden by the pipe.

## References

- `_prd.md` → Core Feature 1; Which skills the contract covers.
- `_techspec.md` → Version identity; Build Order 1.
- `docs/workflow/authorizations/2026-08-04-queue-tail-tooling.md`.
- ADR-0081.
