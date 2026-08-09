---
task: task_06
spec: 0093-a-spec-that-validates-itself
status: completed
type: infra
complexity: high
---

# Task 06: Let the QA gate read product instead of paperwork

## Overview

Measured on Spec 0090, eight of the gate's sixteen rows audited artifacts rather
than product, and the finding that failed the Run came from one of those eight —
reached after 461 tool calls and two context compactions. Every rule those rows
apply is now decided by the checker during authoring. This Task removes them from
the gate's matrix and keeps what a file read cannot settle.

## Requirements

1. MUST remove from the `qa-gate` matrix contract every row whose rule the
   checker now decides, naming each removed rule and the checker rule that runs
   it.
2. MUST keep every row that needs judgement or a running surface: the Spec's
   goals exercised through the surfaces a user reaches, with captured evidence.
3. MUST keep the post-commit rows that no authoring check can answer — which
   paths a Task actually touched against its bounded list — and state that they
   run as commands rather than judgements.
4. MUST NOT reduce what the loop detects. A rule may move; a rule may not
   disappear. Any rule with no checker equivalent stays in the gate.
5. MUST state in the skill that a clean authoring check is a precondition of the
   gate, not a substitute for it.
6. MUST record in the commit message which half of the standing grant it serves
   and how reliability was preserved, per that grant's obligations.

## Subtasks

- [ ] Map each governance row to the checker rule that replaces it.
- [ ] Remove the mapped rows; keep the unmapped ones.
- [ ] State the precondition relationship.

## Acceptance Criteria

- [ ] Every removed row has a named checker rule running its check.
- [ ] No rule is removed without an equivalent.
- [ ] The post-commit rows remain, described as commands.
- [ ] The goal and surface rows are unchanged.

## Rehearsal Cases

- Case: a Spec whose PRD omits a Project Constraints row; Observation: the
  authoring check reports it, and the gate has no row for it.
- Case: a Task whose commit touched a path outside its bounded list;
  Observation: the gate still reports it, because no authoring check can see a
  commit that does not exist yet.
- Case: a Spec whose goals do not work through the CLI; Observation: the gate
  reports it, unchanged from today.

## Bounded scope

Covered by the standing grant at
`docs/workflow/authorizations/2026-08-09-standing-tooling-authority-for-loop-performance.md`.
This Task may create or modify only:

- `.agents/skills/qa-gate/SKILL.md`
- `skills/qa-gate/SKILL.md`, which `make skills-sync` rewrites, following the authorized edit on ADR-0081's principle
- `docs/specs/0093-a-spec-that-validates-itself/task_06.md`

Any other path is out of scope; stop and fail the Task rather than widen it.

## Verification

- `grep -q 'spec check' .agents/skills/qa-gate/SKILL.md` — expected: exits 0, proving the gate names the authoring check it now depends on. This string does not exist in the file before this Task.
- `grep -q 'bounded' .agents/skills/qa-gate/SKILL.md` — expected: exits 0, proving the post-commit path audit was kept rather than removed with the rest.
- `go run -buildvcs=false ./cmd/roundfix skills check` — expected: exits 0.
- `test -z "$(git diff --name-only -- . ':!.agents/skills/qa-gate/SKILL.md' ':!skills/qa-gate/SKILL.md' ':!docs/specs/0093-a-spec-that-validates-itself/task_06.md')"` — expected: exits 0, proving no path outside the bounded list moved.

## Context

- instruction: `docs/workflow/authorizations/2026-08-09-standing-tooling-authority-for-loop-performance.md`

## References

- `_prd.md` → Goal 3.
- `_techspec.md` → Build Order 6.
- ADR-0081, ADR-0091, ADR-0096, ADR-0117.

## Result

The QA skill now requires a clean, full `roundfix spec check <slug> --strict`
authoring result before it builds the matrix and states that this precondition
does not replace the gate. A mapping table names the Task Graph load and every
`SC-*` rule that replaces an authoring-time governance row. If a named detector
is skipped, the skill treats it as no equivalent and retains the corresponding
QA row.

The gate still creates command-backed rows for commit-dependent tooling scope.
Those rows resolve each Task commit's actual paths, current worktree delta,
authorization chronology, prerequisite and consequent fixes, regenerated pins,
and exact bounded files. The skill explicitly retains rules without a checker
equivalent, including missing tooling authorization, outside evidence,
Non-Goals, the report contract, current Task status, and live control or
chronology behavior. Its existing user-story, acceptance, outside-evidence,
Non-Goal, and mandatory-surface row contract is unchanged.

Focused checks:

- `rtk make skills-sync` — exited `0`; it regenerated the embedded skill copy
  from the canonical source.
- `rtk cmp -s .agents/skills/qa-gate/SKILL.md
  skills/qa-gate/SKILL.md` — exited `0` after synchronization.
- `rtk rg -n "precondition of the gate|SC-CONSTRAINT-MISSING|SC-CITATION-UNSUPPORTED|SC-VOCABULARY-UNDOCUMENTED|skipped list|commands, not as a judgement|outside evidence, Non-Goals"
  .agents/skills/qa-gate/SKILL.md skills/qa-gate/SKILL.md` — exited `0` and
  located the precondition, representative mappings, fallback, retained rules,
  and command-backed post-commit audit in both copies.
- Exact `HEAD`/worktree block comparisons with `rtk cmp -s` exited `0` for the
  `Add a row for:` matrix block and the `Surface protocol` block. This proves
  the goal and surface row contracts remained byte-identical.
- `rtk rg -n "SC-(CONSTRAINT|TOOLING|ADR|CITATION|COVERAGE|REF|REQUIREMENT|REHEARSAL|VERIFY|VOCABULARY|LOOP|FINDING|ROLLUP|ARCHIVE|BACKLOG)-"
  internal/speccheck --glob '*.go'` — exited `0` and located every named
  checker family used by the mapping.
- `rtk go run -buildvcs=false ./cmd/roundfix spec check
  0093-a-spec-that-validates-itself --strict` — exited `1` as the new
  precondition requires when it found existing out-of-scope Spec defects:
  `SC-VOCABULARY-UNDOCUMENTED` for `SC-CITATION-UNSUPPORTED` and
  `SC-CITATION-UNSUPPORTED` for the PRD's ADR-0081 claim. It also reported
  `SC-REF-UNRESOLVED` skipped because this Spec has no `references/_index.md`;
  the new fallback therefore keeps that rule in QA rather than treating the
  skip as an equivalent.
- `rtk git diff --check -- .agents/skills/qa-gate/SKILL.md
  skills/qa-gate/SKILL.md
  docs/specs/0093-a-spec-that-validates-itself/task_06.md` — exited `0` after
  the final Result update.

Acceptance evidence:

- Every removed row has a named checker rule running its check: the mapping
  names the Task Graph loader and the exact `SC-*` codes; source lookup found
  every mapped detector family, and the strict check demonstrated that
  checker findings block the gate precondition.
- No rule is removed without an equivalent: a skipped named detector retains
  its row, and the skill lists the non-mechanical and later-created evidence
  classes that remain in QA.
- The post-commit rows remain, described as commands: the skill requires
  command-derived commit paths, worktree delta, chronology, regenerated-pin
  evidence, and bounded-path comparison instead of prose judgement.
- The goal and surface rows are unchanged: both exact `HEAD`/worktree block
  comparisons exited `0`.

Standing-grant commit-message obligation for the Daemon-owned commit:

```text
docs: focus the QA gate on product evidence

Performance: move file-decidable governance rows to roundfix spec check during
authoring.
Reliability: retain named checker mappings, skipped-detector fallbacks,
commit-dependent scope commands, and every goal and surface row.

spec: 0093-a-spec-that-validates-itself / task_06
```

Follow-up outside this Task's slice: the strict authoring precondition exposed
the two existing Spec 0093 findings named above. Their source artifacts are not
in this Task's bounded allowlist and were left unchanged.

The Daemon-owned `## Verification` commands were not run in this Agent turn.
