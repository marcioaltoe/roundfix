---
task: task_08
spec: 0082-the-manifest-already-answered-that
status: failed
type: qa
complexity: high
---

# Task 08: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It executes `qa-gate` after every
other task settles, deriving its matrix from the PRD and the task evidence, and
walking each user story through the real CLI against real repositories rather
than fixtures alone. The gate's distinguishing risk here is that the feature's
whole promise is about what it does *not* change, so the QA evidence must prove
preservation on a repository nobody constructed for the test.

## Overview of required coverage

Beyond the standard matrix, this gate must exercise the fleet reality the PRD
names: nine adopted repositories across three Baseline Profiles, six distinct
catalog digests, and one repository that has never adopted.

## Requirements

1. MUST derive the QA matrix from the PRD's user stories and the completed task
   evidence, and execute it through the CLI surface end to end.
2. MUST prove preservation against a copy of a real adopted repository that was
   not built as a test fixture, asserting that its Repository-Specific Normative
   Rules and authored prose are byte-identical after a refresh.
3. MUST exercise a repository on the `standard-typescript-monorepo` profile, not
   only this repository's profile.
4. MUST prove the sweep shape: a scripted loop over several repository copies
   that reports per-repository outcomes and stops only where a decision is
   genuinely new.
5. MUST prove idempotence on a real repository copy: refresh, then refresh again,
   and observe zero file changes on the second run.
6. MUST prove no ACP runtime is spawned by a refresh, by observation rather than
   by assertion.
7. MUST classify every finding by user impact and record auditable evidence.
8. MUST operate only on copies; no repository outside this checkout may be
   mutated by the gate.
9. MUST create every repository copy outside this repository and outside the
   Run worktree, under the operating system's temporary directory, with the
   shared name prefix `roundfix-qa-0082-`. A copy created inside the repository
   would be committed by the Run, and a copy created inside another project's
   directory would put a real repository at risk. The prefix is fixed so that
   cleanup is checkable by a command rather than by inspection.
10. MUST delete every repository copy it created before the gate finishes,
    including after a failed or abandoned case, and MUST leave no copy behind
    even when a case fails. Deletion targets only paths the gate itself created
    and matched by that prefix; the gate never deletes by a broad pattern that
    could match a path it did not create.
11. MUST record in the QA report, for every copy, the path it used and that the
    path was removed, so cleanup is auditable rather than assumed.

## Rehearsal Cases

- Case: refresh a copy of an adopted repository with a stale catalog; Observation: managed artifacts reach the current catalog digest and the command asks nothing.
- Case: refresh the same copy a second time; Observation: zero file changes reported.
- Case: refresh a copy carrying hand-written Repository-Specific Normative Rules; Observation: those rules are byte-identical after the refresh.
- Case: refresh a copy of a `standard-typescript-monorepo` repository; Observation: either it applies cleanly or it names a genuinely new decision and writes nothing.
- Case: refresh a copy with a hand-edited managed marker; Observation: the command blocks and names the offending path.
- Case: refresh a copy that never adopted; Observation: the command writes nothing and directs the maintainer to first adoption.
- Case: refresh with the upstream skill source unreachable; Observation: guidance applies, the drifted skill is named as a warning, and the exit is successful.
- Case: run the interactive command on an adopted copy; Observation: it reaches plan confirmation with zero decision prompts and no ACP runtime process appears.
- Case: run first adoption on a copy that never adopted; Observation: the full interview runs, unchanged.
- Case: sweep several copies in a loop; Observation: per-repository outcomes are branchable from exit categories without parsing prose.

## Subtasks

- [ ] Derive the resumable QA matrix from the PRD and task evidence.
- [ ] Prepare repository copies covering both profiles and the never-adopted case.
- [ ] Execute every rehearsal case and capture its observation.
- [ ] Observe process creation during a refresh to prove no ACP runtime spawns.
- [ ] Classify findings by user impact and write the dated QA report.
- [ ] Delete every repository copy created, and record each path and its removal.

## Acceptance Criteria

- [ ] Every rehearsal case above is executed and its observation recorded.
- [ ] Every PRD user story has a matrix row with a verdict and evidence.
- [ ] Preservation is proven by digest on a real repository copy, not asserted.
- [ ] Any environment-blocked row carries equivalent evidence and a stated reason.
- [ ] The dated QA report is written under the Spec's `qa/` directory.
- [ ] No repository outside this checkout was mutated.
- [ ] No repository copy remains on disk after the gate finishes, including for
      cases that failed.
- [ ] The QA report lists every copy path used and states that each was removed.
- [ ] The repository working tree contains no leftover copy, so the Run commits
      no fixture repository.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`
- instruction: `.agents/skills/evidence-gate/SKILL.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exits 0.
- `ls docs/specs/0082-the-manifest-already-answered-that/qa/ | grep -q '\.md$'` — expected: exits 0, proving the dated QA report exists.
- `grep -q -i 'byte-identical' docs/specs/0082-the-manifest-already-answered-that/qa/*.md` — expected: exits 0, proving the preservation evidence was recorded.
- `ls -d "${TMPDIR:-/tmp}"/roundfix-qa-0082-* 2>/dev/null | grep . ; test $? -eq 1` — expected: exits 0, proving every repository copy the gate created was removed.
- `git status --porcelain | grep 'roundfix-qa-0082' ; test $? -eq 1` — expected: exits 0, proving no copy was left inside the repository for the Run to commit.
- `grep -q 'roundfix-qa-0082' docs/specs/0082-the-manifest-already-answered-that/qa/*.md` — expected: exits 0, proving the report records the copy paths it used rather than omitting them.
- `go test ./... -count=1` — expected: exits 0.

## References

- `_prd.md` → all User Stories; Success Metrics.
- `_techspec.md` → Testing Approach; Risks & Considerations.
- ADR-0080, ADR-0088, ADR-0091, ADR-0096, ADR-0097, ADR-0100.
