---
task: task_07
spec: 0093-a-spec-that-validates-itself
status: pending
type: qa
complexity: medium
---

# Task 07: Run the final QA gate

## Overview

The authored terminal gate, and deliberately the smallest in recent Specs. What
this Spec can prove by reading files, it proves during authoring; the gate keeps
only what a file read cannot settle — that the command works for the operator
who runs it, and that the checker's own corpus verdict did not move.

## Requirements

1. MUST execute the `qa-gate` skill as this Spec's authored terminal gate and
   write its report to the Spec's QA directory with a machine-readable verdict.
2. MUST include one row in which an operator runs `roundfix spec check` on an
   artifact carrying Spec 0090's original false citation and observes the
   finding, its two quoted texts, and a non-zero exit.
3. MUST include one row proving the default sweep's findings across the corpus
   are unchanged from before this Spec, so no rule was lost by being moved.
4. MUST include one row measuring a scoped run's wall-clock time against the
   0.04 second baseline recorded in the PRD, since affordability at every stage
   is the premise the design rests on.
5. MUST record the origin of the external evidence used, per ADR-0104, and
   separate it from this repository's own artifacts.
6. MUST run the glossary check against `CONTEXT.md`. `SC-CITATION-UNSUPPORTED`
   is a coined token.
7. MUST prove the authoritative gate on a genuinely cold cache — the one the
   Makefile exports, not the user-level one.
8. MUST report `rows_blocked_environment`, `rows_blocked_finding`, and
   `rows_blocked_declared` counts in the report frontmatter.
9. MUST include one row proving each authoring skill this Spec wired refuses to
   report while a finding stands, exercised through the skill's own instruction
   rather than by reading it.
10. MUST include one row proving the QA gate's own matrix no longer carries a
    row decidable by reading files, and that every rule it dropped is running in
    the checker — named rule by rule.
11. MUST NOT audit this Spec's own artifacts by reading them. Every such check
   now runs during authoring, and duplicating it here would reintroduce the cost
   this Spec exists to remove.

## Subtasks

- [ ] Run the gate and write the report.
- [ ] Observe the finding through the command, as an operator.
- [ ] Prove the corpus verdict is unchanged and measure a scoped run.
- [ ] Run the glossary check.

## Acceptance Criteria

- [ ] A QA report exists with a machine-readable verdict and all three
      blocked-row counts.
- [ ] One row records the operator seeing the finding with both texts and a
      non-zero exit.
- [ ] One row records the corpus verdict unchanged.
- [ ] One row records a scoped run's measured time.
- [ ] The report records the glossary check for `SC-CITATION-UNSUPPORTED`.
- [ ] One row records each wired skill blocking on a finding.
- [ ] One row maps every rule the gate dropped to the checker rule now running
      it.
- [ ] No row re-reads this Spec's own artifacts for governance.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`
- instruction: `docs/agents/domain.md`

## Verification

- `test -d docs/specs/0093-a-spec-that-validates-itself/qa` — expected: exits 0.
- `grep -lq 'rows_blocked_declared' docs/specs/0093-a-spec-that-validates-itself/qa/qa-report-*.md` — expected: exits 0.
- `grep -lq 'SC-CITATION-UNSUPPORTED' docs/specs/0093-a-spec-that-validates-itself/qa/qa-report-*.md` — expected: exits 0.
- `grep -lqE 'verdict: (pass|fail|partial)' docs/specs/0093-a-spec-that-validates-itself/qa/qa-report-*.md` — expected: exits 0, proving a verdict was written rather than left pending.
- `GOCACHE="$PWD/.gocache" go clean -testcache` — expected: exits 0.
- `make verify` — expected: exits 0 on a genuinely cold cache, with no Go package line reporting `(cached)`.

## References

- `_prd.md` → every Goal.
- `_techspec.md` → Testing Approach; Build Order 7.
- ADR-0091, ADR-0096, ADR-0104, ADR-0116, ADR-0117.
