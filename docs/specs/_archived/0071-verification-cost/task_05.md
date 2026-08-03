---
task: task_05
spec: 0071-verification-cost
status: completed
type: docs
complexity: medium
---

# Task 05: Stop charging every Task for the whole suite

## Overview

Spec 0057's fourteen Tasks each carried a whole-package suite command as the
last line of their Verification — about 28 minutes per pass, repaid on every
retry, proving something the Run-level gate already proves. This Task removes
those commands from every active Spec and records the authoring rule that keeps
them from coming back.

## Requirements

1. MUST remove whole-package suite commands from the Verification of every
   active Spec's Task files, leaving each Task's focused checks that prove its
   own effect.
2. MUST NOT weaken any focused check, scope assertion, or build command; only
   the redundant whole-suite line goes.
3. MUST record the authoring rule in the Task-authoring skill: a Task proves
   its own effect, and the Run-level gate proves nothing else regressed.
4. MUST record alongside it the rule that a Verification command must be able
   to fail when no work was done, so a command naming a missing test cannot
   pass.
5. MUST leave archived Specs byte-identical.
6. MUST change only the owned skill pair among protected tooling, per this
   Spec's Tooling authority row; any other tooling path is out of scope.

## Subtasks

- [ ] Remove whole-package suite lines from active Specs' Task Verification.
- [ ] Confirm each Task keeps a focused check proving its own effect.
- [ ] Record both authoring rules in the Task-authoring skill.
- [ ] Confirm archived Specs are untouched.

## Acceptance Criteria

- [ ] No active Spec's Task file carries a whole-package suite command in its
      Verification.
- [ ] Every active Task file still carries at least one focused check naming a
      specific test or asserting a specific effect.
- [ ] The Task-authoring skill states that a Task proves its own effect and the
      Run-level gate owns regression.
- [ ] The Task-authoring skill states that a Verification command must be able
      to fail when no work was done.
- [ ] Archived Spec artifacts are byte-identical.
- [ ] `git status --porcelain` shows no path outside `docs/specs/`,
      `.agents/skills/`, `skills/`, and this task file.

## Context

- instruction: `docs/agents/autonomous-work.md`
- interface: `.agents/skills/write-tasks/SKILL.md`

## Verification

- `grep -rl -- "go test ./internal/" docs/specs --include='task_*.md' | grep -v _archived | xargs -r grep -l -- '-count=1$' | grep -q . && exit 1 || exit 0`
  — expected: exit 0; no active Task ends a Verification line with a bare
  whole-package suite command.
- `grep -q 'prove its own effect' .agents/skills/write-tasks/SKILL.md` —
  expected: exit 0; the authoring rule is recorded.
- `grep -q 'able to fail when no work was done' .agents/skills/write-tasks/SKILL.md`
  — expected: exit 0; the vacuous-Verification rule is recorded.
- `diff -r .agents/skills/write-tasks skills/write-tasks` — expected: exit 0;
  the embedded mirror does not drift.
- `git diff --name-only HEAD -- docs/specs/_archived | grep -q . && exit 1 || exit 0`
  — expected: exit 0; no archived artifact changed.

## References

- `_prd.md` → Core Features 3 and 4; Goals (verification proportional to what
  changed).
- `_techspec.md` → Build Order 5; Project Constraints: Tooling authority.

## Result

### Implementation

- Removed only the redundant single-run package suite from the Verification of
  active Tasks 01 through 04. Their build, coverage-equivalence, repeated-run,
  race, vet, and scope checks remain unchanged.
- Updated both copies of the Task-authoring skill with the rules that a Task
  must prove its own effect, the Run-level gate proves nothing else regressed,
  and every Verification command must be able to fail when no work was done.
- Ran the ADR-0081-sanctioned digest regeneration after changing the
  Roundfix-owned skill. It updated seven derived Baseline digest artifacts
  under `internal/baseline/`.

### Focused-check evidence

- A focused `rtk rg` scan of active Task files found no remaining bare
  `go test ./internal/<package> -count=1` Verification line (no matches).
- A Verification-block scan across every active Task file exited 0; each block
  retains a command naming a test or asserting a specific effect. Tasks 01
  through 04 retain `TestCoverageEquivalence`, repeated-run, race, grep, or vet
  checks as applicable.
- `rtk rg` found `prove its own effect`, `Run-level gate proves nothing else
  regressed`, and `able to fail when no work was done` in both skill copies.
- `rtk go test ./skills -run '^TestAuthorialSkillSync$' -count=1 -v` exited 0;
  RTK reported 18 passing cases.
- `rtk git diff --no-index -- .agents/skills/write-tasks skills/write-tasks`
  exited 0 with no output; the owned skill pair is identical.
- `rtk git -c core.fsmonitor=false status --short docs/specs/_archived` exited
  0 with no output; archived Specs have no worktree changes.
- `rtk make baseline-digests` exited 0 and regenerated the seven derived files
  named above. The authored tooling change remains limited to the owned skill
  pair; the additional paths are deterministic ADR-0081 fallout.

### Acceptance criteria evidence

- No active Task Verification carries the redundant single-run package suite:
  supported by the no-match active-tree scan.
- Every active Task retains an effect-specific check: supported by the
  Verification-block scan and the retained Task 01 through 04 checks.
- Both Task-authoring rules are present: supported by the exact-text scan of
  both skill copies.
- Archived Specs are byte-untouched in the worktree: supported by the scoped
  status inspection.
- Authored paths are confined to active Spec Task files, this Result, and the
  owned skill pair. `internal/baseline/` also contains only the deterministic
  digest fallout produced by the repository-sanctioned regeneration command.

The Task's declared Verification commands were not rerun in this Agent turn.
