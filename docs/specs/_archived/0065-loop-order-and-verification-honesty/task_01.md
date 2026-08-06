---
task: task_01
spec: 0065-loop-order-and-verification-honesty
status: completed
type: chore
complexity: medium
---

# Task 01: State the loop order once

## Overview

`docs/agents/autonomous-work.md` states the corrected order in one place and
the original order in another, so the contradiction is live and both are being
followed. The same order is carried to adopting repositories by
`internal/baseline/assets/modules/autonomous-work.json`, which is the third
place it can disagree.

The order itself was settled on 2026-08-05 and does not change here: ADR-0091
stands, so the QA gate is the graph's terminal Task and runs before any Pull
Request exists.

## Requirements

1. MUST state one order in every place it appears: implement the graph
   including its authored gate, archive, open the Pull Request, watch until
   Clean, merge.
2. MUST remove the contradicting restatement rather than annotate it.
3. MUST carry the same order in the Baseline module asset, so adopting
   repositories receive the corrected clause.
4. MUST record why the order is this one: a Spec whose acceptance observes its
   own Pull Request reaches `pass` through ADR-0080's environment-blocked rows
   with equivalent evidence, proven by Spec 0078's gate on 2026-08-05.
5. MUST change only `docs/agents/autonomous-work.md`, the Baseline module asset
   `internal/baseline/assets/modules/autonomous-work.json`, this Task file, and
   the ADR-0081 digest fallout under `DERIVED_DIGEST_PATHS`.
6. MUST run `make baseline-digests` after the module asset edit, then re-record
   the two characterization corpora that command does not reach:

   ```bash
   go test ./internal/baseline -count=1 \
     -run TestBaselinePlanCharacterization -update-baseline-plan-characterization
   go test ./internal/baseline -count=1 \
     -run TestCatalogDiagnosticCharacterization -update-catalog-diagnostics
   ```

7. MUST NOT change Go source, and MUST NOT add the divergence check — that is
   task_04, which depends on this Task precisely so the check never lands while
   the sources disagree.

## Subtasks

- [x] Remove the contradicting restatement and state the order once.
- [x] Carry the same order into the Baseline module asset.
- [x] Run the regeneration chain and both re-records.

## Acceptance Criteria

- [x] `docs/agents/autonomous-work.md` states the order exactly once.
- [x] The Baseline module asset states the same order.
- [x] The recorded reason names ADR-0091, ADR-0080, and the Spec 0078 evidence.
- [x] No Go source changed.
- [ ] `make verify` exits 0.

## Context

- instruction: `docs/agents/autonomous-work.md`
- instruction: `docs/adr/0091-the-qa-gate-is-a-task-node-of-its-own-type.md`

## Verification

- `make verify` — expected: exit 0.
- `git diff --name-only HEAD | grep -E "\.go$" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; no Go source changed.
- `git diff --name-only HEAD | grep -vE "^(docs/agents/autonomous-work\.md$|docs/specs/0065-loop-order-and-verification-honesty/task_01\.md$|internal/baseline/(assets/(modules|setups|source-baselines|formatter-fixtures|profiles)|testdata)/)" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; only the bounded paths changed.

## References

- `_prd.md` → Core Features 1 and 2; Success Metric 1.
- `_techspec.md` → The order, decided; Build Order 1.
- `docs/workflow/authorizations/2026-08-04-queue-tail-tooling.md`.
- ADR-0091, ADR-0080, ADR-0081.

## Result

Implementation evidence:

- Removed the PR-before-gate restatement from the QA-gate authoring section and
  left one canonical sequence in the implementation loop: implement the graph
  including its authored gate, archive, open the Pull Request, watch until
  Clean, and merge.
- Updated `rule.autonomous.loop` in the Baseline module to the same sequence and
  bumped the module and supporting-guide versions from 7 to 8. The generated
  formatter fixture now carries that sequence to adopting repositories.
- Both carriers explain that ADR-0091 keeps the gate terminal, ADR-0080 admits
  equivalent evidence for environment-blocked rows, and Spec 0078's
  2026-08-05 gate proved the path with eleven of eighteen rows blocked, nine of
  them on no open Pull Request. Corrected from a stale `nine of eighteen`
  against the archived Spec 0078 report; task_07 records the same repair for
  the shipped carriers.

Focused-check evidence:

- `rtk make baseline-digests` — exit 0; regenerated the authorized ADR-0081
  fallout under `DERIVED_DIGEST_PATHS`.
- `rtk go test ./internal/baseline -count=1 -run TestBaselinePlanCharacterization -update-baseline-plan-characterization`
  — the sandboxed attempt could not access the shared Go build cache; the
  authorized rerun exited 0 with 7 tests passed.
- `rtk go test ./internal/baseline -count=1 -run TestCatalogDiagnosticCharacterization -update-catalog-diagnostics`
  — exit 0 with 2 tests passed.
- `rtk jq empty internal/baseline/assets/modules/autonomous-work.json` — exit 0.
- Focused searches found one canonical-order occurrence in
  `docs/agents/autonomous-work.md`, one in the module asset, and none of the
  removed PR-before-gate wording in either the guide or regenerated carriers.
- `rtk git -c core.fsmonitor=false diff --name-only HEAD -- '*.go'` — exit 0
  with no output; no Go source changed.
- Changed-path inspection found only the guide, module asset, this Task file,
  and generated files under the Task-authorized `DERIVED_DIGEST_PATHS`.

Daemon-owned verification:

- `make verify` and the other commands under `## Verification` were not run in
  this Agent turn; the Daemon owns them and the terminal Task verdict.

Verification Feedback evidence, attempt 1:

- The Daemon's `make verify` attempt reached the dedicated
  `TestCheckCorpusBudget` step after 3,457 Go tests passed. That step measured
  the Spec corpus sweep at 1.264 seconds against its 1-second wall-clock
  budget.
- A focused rerun of
  `rtk go test -count=1 -parallel=1 ./internal/speccheck -run '^TestCheckCorpusBudget$'`
  exited 0 with the unchanged Task implementation.
- Five verbose focused repetitions measured the sweep at 358.9 ms, 355.9 ms,
  367.1 ms, 357.5 ms, and 370.4 ms; all passed below the 1-second budget. The
  Daemon failure is not reproducible under the same serial selector.
- No Go source, test, or timing threshold was changed: those mutations are
  outside this Task's explicit scope and would hide a transient wall-clock
  signal instead of repairing this documentation and Baseline contract.
