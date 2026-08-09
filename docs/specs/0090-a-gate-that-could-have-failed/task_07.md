---
task: task_07
spec: 0090-a-gate-that-could-have-failed
status: pending
type: qa
complexity: high
---

# Task 07: Run the final QA gate

## Overview

The authored terminal gate. It walks every PRD goal against the built binary and
carries one acceptance no fixture substitutes for: a real Roundfix Run whose Task
Graph contains a deliberately vacuous gate, refused at dispatch with no Agent
turn spent. This node also carries the glossary check.

## Requirements

1. MUST execute the `qa-gate` skill as this Spec's authored terminal gate and
   write its report to the Spec's QA directory with a machine-readable verdict.
2. MUST walk every PRD goal against the built binary rather than a fixture.
3. MUST include one row in which a real Roundfix Run executes a Task Graph
   containing a Task whose Verification passes against the unchanged tree,
   recording the Run identifier, the refusal, the offending command named in the
   event stream, and evidence that no Agent turn was spent.
4. MUST record that row as blocked with its reason rather than dropping it when
   the runtime, the network, or the session's own permissions make the Run
   impossible.
5. MUST include one row that proves the authoritative gate returns the same
   verdict twice on one unchanged tree, run cold both times.
6. MUST include one row that exercises the unknown outcome end to end and shows
   its terminal reason is distinguishable from a command that ran and failed.
7. MUST record the origin of the external evidence used, per ADR-0104. The PRD's
   `## External evidence` section names it, and the row must state which of its
   claims the Spec relied on and which came from this repository's own artifacts.
8. MUST run the glossary check against `CONTEXT.md` and record whether this Spec
   coined, changed, or dropped a term. `Negative Control` is a coined term.
9. MUST prove the authoritative gate on a genuinely cold cache — the one the
   Makefile exports, not the user-level one.
10. MUST report `rows_blocked_environment`, `rows_blocked_finding`, and
    `rows_blocked_declared` counts in the report frontmatter.
11. MUST NOT mutate any repository outside this one.

## Subtasks

- [ ] Run the gate and write the report to the Spec's QA directory.
- [ ] Walk every PRD goal against the built binary.
- [ ] Execute a real Run whose graph carries a vacuous gate.
- [ ] Prove the authoritative gate twice on one tree.
- [ ] Run the glossary check and record its result.

## Acceptance Criteria

- [ ] A QA report exists with a machine-readable verdict and all three
      blocked-row counts in frontmatter.
- [ ] Every PRD goal has a row.
- [ ] One row records a Run in which a vacuous gate was refused, with the
      offending command named and no Agent turn spent, or a recorded reason it
      was blocked.
- [ ] One row records two cold runs of the authoritative gate agreeing.
- [ ] One row distinguishes an unknown outcome from a failed command.
- [ ] The report separates external evidence from this repository's own.
- [ ] The report records the glossary check result for `Negative Control`.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`
- instruction: `docs/agents/domain.md`

## Verification

- `test -d docs/specs/0090-a-gate-that-could-have-failed/qa` — expected: exits 0.
- `grep -lq 'rows_blocked_declared' docs/specs/0090-a-gate-that-could-have-failed/qa/qa-report-*.md` — expected: exits 0.
- `grep -lq 'Negative Control' docs/specs/0090-a-gate-that-could-have-failed/qa/qa-report-*.md` — expected: exits 0, proving the glossary check was recorded.
- `grep -lqE 'verdict: (pass|fail|partial)' docs/specs/0090-a-gate-that-could-have-failed/qa/qa-report-*.md` — expected: exits 0, proving a machine-readable verdict was written rather than left pending.
- `GOCACHE="$PWD/.gocache" go clean -testcache` — expected: exits 0. A bare `go clean -testcache` clears the user-level cache, not the one the Makefile exports.
- `make verify` — expected: exits 0 on a genuinely cold cache, with no Go package line reporting `(cached)`.

## References

- `_prd.md` → every Goal; `## External evidence`.
- `_techspec.md` → Testing Approach; Build Order 7.
- ADR-0091, ADR-0096, ADR-0104, ADR-0109, ADR-0110, ADR-0111.
