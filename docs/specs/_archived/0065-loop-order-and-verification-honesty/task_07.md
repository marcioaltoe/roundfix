---
task: task_07
spec: 0065-loop-order-and-verification-honesty
status: completed
type: chore
complexity: low
---

# Task 07: Correct the blocked-row count the loop clause cites

## Overview

Corrective Task from the QA gate's F-002 (`Trust-Damage`). The loop-order
rationale cites Spec 0078's gate as passing with **nine of eighteen** rows
blocked. Its archived report records **eleven**: nine of them blocked on the
absent Pull Request, plus one on host-load timing variance and one on a sandbox
that denied Agent state.

The conclusion the citation supports is true — ADR-0080's blocked-row typing
absorbs a Spec whose acceptance observes its own Pull Request, which is why
this Spec keeps the ADR-0091 order. What is false is the number offered as its
proof, and a rationale carrying a wrong measurement teaches the reader to
distrust the correct ones beside it.

## Requirements

1. MUST correct the count in every carrier to eleven of eighteen, naming that
   nine of those eleven were blocked on the absent Pull Request.
2. MUST correct it in all five carriers the gate found:
   `docs/agents/autonomous-work.md`, its Baseline formatter golden,
   `internal/baseline/assets/modules/autonomous-work.json`, this Spec's
   `_techspec.md`, and the Skills that restate it.
3. MUST verify the corrected number against the archived Spec 0078 report
   rather than against another carrier, so the correction cannot inherit the
   defect it repairs.
4. MUST keep the surrounding argument unchanged: the conclusion was never the
   error.
5. MUST run `make skills-sync` if any Skill text changed, then
   `make baseline-digests`, then re-record the two characterization corpora:

   ```bash
   go test ./internal/baseline -count=1 \
     -run TestBaselinePlanCharacterization -update-baseline-plan-characterization
   go test ./internal/baseline -count=1 \
     -run TestCatalogDiagnosticCharacterization -update-catalog-diagnostics
   ```

6. MUST leave `SC-LOOP-ORDER-DIVERGENT` passing: the three order statements
   must still agree after the edit.
7. MUST NOT change Go source.

## Subtasks

- [x] Read the archived Spec 0078 report and record the exact counts.
- [x] Correct every carrier and run the regeneration chain.

## Acceptance Criteria

- [x] No carrier states `nine of eighteen`.
- [x] Every carrier states eleven of eighteen, with nine attributed to the
      absent Pull Request.
- [x] The number matches the archived Spec 0078 report, cited in the Result.
- [x] `SC-LOOP-ORDER-DIVERGENT` passes.
- [x] No Go source changed.

## Context

- instruction: `docs/agents/autonomous-work.md`
- interface: `docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/qa-report-2026-08-05.md`

## Verification

- `grep -rn "nine of eighteen" docs/agents .agents skills internal/baseline/assets | grep -q . && exit 1 || exit 0`
  — expected: exit 0; no shipped carrier states the wrong count. Scoped to the
  shipped surfaces on purpose: this Spec's own artifacts quote the wrong string
  while explaining the correction, and an absence check that cannot tell a
  claim from a citation of that claim is the over-strict shape this Spec exists
  to refuse.
- `grep -q "eleven of eighteen" docs/agents/autonomous-work.md` — expected:
  exit 0; the guide carries the corrected count.
- `grep -q "eleven of eighteen" docs/specs/0065-loop-order-and-verification-honesty/task_01.md`
  — expected: exit 0; the Task evidence carries it too, which the first pass of
  this Task failed to prove because its search omitted Task files.
- `go run -buildvcs=false ./cmd/roundfix spec check > /dev/null` — expected:
  exit 0; the order statements still agree.
- `make verify` — expected: exit 0.
- `git diff --name-only HEAD | grep -E "\.go$" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; no Go source changed.

## References

- `qa/qa-report-2026-08-05.md` → F-002; row R10.
- `_techspec.md` → The order, decided.
- `docs/workflow/authorizations/2026-08-04-queue-tail-tooling.md`.
- ADR-0080, ADR-0081.

## Result

Corrected the quantitative proof without changing the loop order or its
ADR-0080/ADR-0091 conclusion. The guide, shipped formatter golden, Baseline
module, TechSpec, and both order-restating Skills now say that eleven of
eighteen rows were blocked and that nine of those eleven were blocked on the
absent Pull Request. `make skills-sync` updated the two shipped Skill mirrors;
`make baseline-digests` and the two required characterization update commands
re-recorded their deterministic fallout.

The direct source of truth was
`docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/qa-report-2026-08-05.md`:
its frontmatter records `rows_blocked_environment: 11`; its blocked-row detail
assigns R06–R13 and R18 (nine rows) to `no open Pull Request`, R03 to host-load
timing variance, and R04 to sandbox denial of Agent state; its Coverage section
records 7 pass and 11 environment-blocked rows out of 18.

Focused implementation evidence:

- `rtk make skills-sync` — exit 0.
- `rtk make baseline-digests` — exit 0; regenerated the sanctioned Baseline
  digest and characterization fallout.
- `rtk go test ./internal/baseline -count=1 -run TestBaselinePlanCharacterization -update-baseline-plan-characterization`
  — the sandboxed attempt could not open the configured user-level Go build
  cache; the approved full-access rerun exited 0 with 7 passing cases.
- `rtk go test ./internal/baseline -count=1 -run TestCatalogDiagnosticCharacterization -update-catalog-diagnostics`
  — exit 0 with 2 passing cases.
- `rtk rg -n 'nine of eighteen' docs/agents .agents skills internal/baseline/assets docs/specs/0065-loop-order-and-verification-honesty/_techspec.md`
  — exit 1 with no output, proving the stale phrase is absent from the bounded
  carrier scope.
- `rtk rg -n 'eleven of eighteen|nine of those eleven' docs/agents .agents skills internal/baseline/assets docs/specs/0065-loop-order-and-verification-honesty/_techspec.md`
  — exit 0 and listed the guide, module, formatter golden, TechSpec, both
  canonical Skills, and both mirrors with the corrected total and attribution.
- `rtk cmp -s .agents/skills/write-tasks/SKILL.md skills/write-tasks/SKILL.md`
  and `rtk cmp -s .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md`
  — both exited 0; each canonical Skill is byte-identical to its shipped mirror.
- `rtk go test ./internal/speccheck -count=1 -run '^TestCheckLoopOrder'` — exit
  0 with 5 passing cases; the corrected repository and all three seeded
  carrier-divergence cases retain `SC-LOOP-ORDER-DIVERGENT` behavior.
- `rtk go test ./internal/baseline -count=1 -run '^TestFormatterComposition$'`
  — exit 0 with 1 passing case for the formatter golden.
- `rtk git diff --check` — exit 0.
- `rtk git diff --name-only HEAD` — exit 0; the changed paths contain no `.go`
  file. Changes outside the authored carriers and assigned Task are generated
  Skill mirrors or sanctioned Baseline digest/characterization fallout.

Acceptance evidence:

- No stale count: the bounded absence search returned no match.
- Correct total and attribution: every carrier states eleven of eighteen and
  identifies nine of those eleven as absent-Pull-Request rows.
- Archived-report agreement: the report's frontmatter, row detail, and Coverage
  independently establish 11 total environment blocks and the 9/1/1 cause
  split recorded above.
- Order consistency: all 5 focused `TestCheckLoopOrder` cases passed.
- Go-source boundary: the post-edit changed-path inventory contains no `.go`
  file.

The Task's declared `## Verification` commands were not run; the Daemon owns
that verification and terminal settlement.
