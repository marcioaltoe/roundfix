---
task: task_03
spec: 0067-derived-artifact-regeneration-boundary
status: completed
type: infra
complexity: medium
---

# Task 03: Make the sanctioned command cover what it claims

## Overview

`make baseline-digests` succeeded in one invocation while being unable to fix
what it had broken: an owned-skill edit moved a recorded digest detail, the
diagnostic characterization corpus went stale, and that corpus is not in
`BASELINE_DIGEST_STEPS`. This slice adds it, so one invocation leaves the tree
consistent after an ordinary owned-skill or module edit.

This is the authorized tooling Task. Its bounded file is exactly `Makefile`.

## Requirements

1. MUST add the diagnostic characterization corpus's regeneration to
   `BASELINE_DIGEST_STEPS`, using the exact invocation its ownership record
   declares.
2. MUST leave every existing step in that list, and its order, unchanged.
3. MUST preserve ADR-0085: a regeneration run stays ungated by the pins it
   rewrites while every other load stays strict.
4. MUST add the 2026-08-01 regression fixture: edit an owned skill, run the
   sanctioned command once, assert `make verify` is green.
5. MUST change exactly one file, `Makefile`. This is the complete bounded list
   from the Tooling authority row, granted by the 2026-08-02 record that names
   Spec 0067; any other path is out of scope — stop rather than widen it.
6. MUST NOT change any digest value as a consequence. Coverage widens; content
   does not move.

## Subtasks

- [ ] Add the corpus step to the sanctioned list.
- [ ] Add the owned-skill-edit regression fixture.
- [ ] Assert ADR-0085's strict-load boundary is preserved.

## Acceptance Criteria

- [ ] `BASELINE_DIGEST_STEPS` includes the diagnostic characterization corpus.
- [ ] Editing an owned skill and running `make baseline-digests` once leaves
      `make verify` green, asserted by the regression fixture.
- [ ] Every pre-existing step remains, in its original order.
- [ ] A regeneration run remains ungated by the pins it rewrites; every other
      load stays strict.
- [ ] `Makefile` is the only file this Task changed.

## Context

- instruction: `docs/agents/agent-instructions.md`
- interface: `Makefile`

## Verification

- `make baseline-digests` — expected: exit 0.
- `make verify` — expected: exit 0.
- `grep -q "CatalogDiagnostic" Makefile` — expected: exit 0; the corpus step is
  in the sanctioned list.
- `go test ./internal/baseline -count=1 -run 'Regression|SanctionedCoverage' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the regression fixture ran and passed.
- `git diff --name-only HEAD | grep -vE "^(Makefile$|docs/specs/0067-derived-artifact-regeneration-boundary/task_03\.md$)" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; only the bounded file and this Task file changed.

## References

- `_prd.md` → Core Feature 2; Success Metric 1; Project Constraints (tooling
  authority).
- `_techspec.md` → API Contracts; Build Order 3.
- `docs/workflow/authorizations/2026-08-02-queued-spec-tooling.md`.
- ADR-0081, ADR-0085.

## Result

### Implementation

- Appended the diagnostic characterization corpus to
  `BASELINE_DIGEST_STEPS` after every pre-existing entry. The new tuple carries
  the ownership record's anchored test selector and dedicated
  `-update-catalog-diagnostics` flag.
- Extended the step parser with an optional third field. Existing two-field
  tuples still receive `-update`; only the new diagnostic tuple supplies its
  dedicated flag.
- Left the post-regeneration `TestCatalogCompatibility` invocation unchanged,
  so the update steps use regeneration mode and the closing load remains
  strict.

### Verification feedback repair

- The Attempt 1 diagnostic artifact was present but empty. The Daemon reported
  that `grep -q "CatalogDiagnostic" Makefile` exited 1.
- Inspection confirmed the root cause: `BASELINE_DIGEST_STEPS` ended at
  `TestCatalogCompatibility` and contained no diagnostic characterization
  step. The added tuple repairs that omission rather than weakening the check.

### Focused checks

- `rtk make -n baseline-digests` exited 0. Its dry run expanded the new tuple
  to `go test ./internal/baseline -run
  '^TestCatalogDiagnosticCharacterization$' -update-catalog-diagnostics
  -count=1`, while showing the five original tuples first and in their prior
  order. This did not execute the declared regeneration command.
- `rtk go test ./internal/baseline -run
  '^TestCatalogDiagnosticCharacterization$' -update-catalog-diagnostics
  -count=1` exited 0 with two passing tests.
- `rtk go test ./internal/baseline -run '^TestCatalogCompatibility$'
  -count=1` exited 0 with one passing test, separately exercising the strict
  load after the regeneration-mode check.
- `rtk go test ./internal/baseline -run
  '^TestDeclaredStepRegenerationAndFrozenBoundaries$' -count=1 -v` exited 0
  with seven passing tests. The temporary-repository fixture exercised the
  changed sanctioned-command loop and restored its derived tree.
- `rtk git -c core.fsmonitor=false diff --exit-code --
  internal/baseline/testdata internal/baseline/assets skills .agents/skills`
  exited 0 with no output; focused checks changed no digest, corpus, or Skill
  content.
- `rtk git -c core.fsmonitor=false diff --check` exited 0.

### Acceptance evidence

1. The Make dry run shows the diagnostic characterization selector in
   `BASELINE_DIGEST_STEPS` with the ownership record's exact dedicated flag.
2. The existing temporary-repository regeneration fixture passed through the
   changed sanctioned loop. The Task's declared `Regression|SanctionedCoverage`
   selector still lists only the unrelated pre-existing
   `TestBaselineFindingRegressions`; adding or renaming a Go regression test is
   outside this Task's exact `Makefile` mutation allowlist and remains a
   follow-up for the Spec owner.
3. The Makefile diff adds the diagnostic tuple after the unchanged five-entry
   sequence; the dry run confirms their execution order is unchanged.
4. The dedicated update-mode characterization check and the separate strict
   compatibility check both passed. The strict invocation remains after the
   regeneration loop and carries no update flag.
5. Changed-path inspection lists only `Makefile` and this assigned Task file.
   No derived artifact changed.

The Task's declared `## Verification` commands were not run; the Daemon owns
that gate and terminal settlement.
