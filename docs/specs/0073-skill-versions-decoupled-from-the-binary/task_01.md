---
task: task_01
spec: 0073-skill-versions-decoupled-from-the-binary
status: completed
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

## Result

### Implementation

- Every canonical member read from `OWNED_SKILLS` now declares top-level
  `version: 0.0.2`, matching its existing authored metadata version.
- `skills-version-check` iterates `OWNED_SKILLS` across `.agents/skills` and
  `skills`, rejects a missing top-level declaration, and is a prerequisite of
  `skills-sync-check`. It never scans unrelated third-party skills.
- `make skills-sync` regenerated the embedded mirrors. The ADR-0081 chain
  regenerated the derived Baseline catalog, parity, plan-characterization,
  and diagnostic artifacts through their sanctioned commands.

### Focused checks

- Pre-change Makefile-variable probe: exit 1, naming all 14 owned canonical
  skills as missing a top-level version.
- `rtk make skills-version-check`: exit 0.
- `rtk make skills-version-check OWNED_SKILLS='roundfix coding-guidelines'`:
  expected exit 2 because the injected member lacks a top-level version. This
  proves that adding a member to the declaring variable cannot skip the check.
- Makefile-derived `diff -r` loop over canonical and mirror roots: exit 0,
  `all owned mirrors match`.
- `rtk make skills-sync`: exit 0.
- `rtk make baseline-digests`: exit 0; reported `ok: true` and regenerated the
  derived artifacts.
- The first direct characterization attempt did not start its test because the
  inherited Go cache pointed outside the Task Worktree. With `GOCACHE` set to
  the repository-local `.gocache`, both required commands passed:
  `TestBaselinePlanCharacterization` (`ok`, 1.486s) and
  `TestCatalogDiagnosticCharacterization` (`ok`, 0.497s).
- `go test ./skills -count=1 -run '^TestAuthorialSkillSync$'` with the local
  Go cache: exit 0 (`ok roundfix/skills`, 0.244s).
- `rtk git -c core.fsmonitor=false diff --name-only -- .agents/skills skills`
  lists only the canonical and mirrored `OWNED_SKILLS` members; no third-party
  skill changed.
- `rtk git -c core.fsmonitor=false diff --name-only -- '*.go'`: no output, so
  this slice adds no Go reader or version-comparison behavior.

### Verification Feedback repair

- Daemon attempt 1 reached the existing owned-skill disagreement test and
  failed because its fixture replaces the first textual `version: 0.0.2` to
  exercise `metadata.version`. The initial field order made that first match
  the new top-level declaration instead, so the intended metadata mismatch was
  no longer produced.
- Moved each new top-level `version` to the end of its existing frontmatter.
  The declaration remains top-level and inert, while the pre-existing fixture
  continues to exercise the metadata contract it owns. No Go or test file was
  changed.
- Regenerated mirrors and the full ADR-0081 chain from the repaired canonical
  skills. `TestBaselinePlanCharacterization` passed (`ok`, 5.761s) and
  `TestCatalogDiagnosticCharacterization` passed (`ok`, 0.729s) with the
  repository-local Go cache.
- Focused reproduction
  `go test ./skills -count=1 -run '^TestOwnedSkillContractRejectsSetAndVersionDisagreement$'`:
  exit 0 (`ok roundfix/skills`, 5.174s).
- Post-repair `rtk make skills-version-check`: exit 0; its injected unversioned
  member negative control still exits 2 as expected.
- Post-repair byte-exact mirror loop: exit 0, `all owned mirrors match`;
  focused `TestAuthorialSkillSync`: exit 0 (`ok roundfix/skills`, 0.206s).
- The Daemon-owned declared Verification commands were not rerun in this Agent
  turn.

### Acceptance evidence

- Every owned canonical and mirror declares the version: supported by the
  positive version check and byte-exact mirror comparison.
- The member set comes from `OWNED_SKILLS`: supported by the Make recipe and
  its expected negative control with an injected unversioned member.
- Third-party skills remain unversioned and outside the check: supported by the
  normal check passing and the skill-path diff containing owned members only.
- Every mirror is byte-identical: supported by the Makefile-derived `diff -r`
  loop and focused `TestAuthorialSkillSync` run.
- `make verify`: not run; the Daemon owns the declared Verification commands.
- No Go behavior reads the field: supported by the empty Go-file diff.
