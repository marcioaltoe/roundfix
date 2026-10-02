---
task: task_03
spec: 0212-gates-and-run-storage-that-let-a-correct-delivery-finish
status: completed
type: backend
complexity: high
---

# Task 03: A Task declares a path it deletes

## Overview

A Task's `## Context` knows `instruction`, `interface` and `creates`, so
Spec 0200's task_04 declared a deleted skill file as `interface`, the Spec
Consistency Check refused it once the file was gone, and the operator
relabelled it `creates`. This Task adds the `deletes` kind: the
unresolved-path check skips it, the scope audit, the Governed Path
declaration audit and the Wave collision rule count it, the Daemon will not
settle a Task while the path still exists, and the write-tasks skill and
template teach it.

## Requirements

1. MUST add `ContextKindDeletes` and accept `- deletes: <path>` with the path
   rules of the other kinds, refusing a repeated `deletes` path and a path
   declared under two of `interface`, `creates` and `deletes`, with an error
   naming both.
2. MUST make `detectTaskContextReferences` skip `deletes`, and MUST make
   `UndeclaredTaskPaths`, `declaredGovernedTouches`, the Wave collision rule
   and `firstCLISurface` count it, as "The `deletes` Context kind" lists; the
   Daemon's Task context MUST list it with the interface paths.
3. MUST add `undeletedTaskPaths` in the new file
   `internal/daemon/task_deletes.go` and make `verifyTask` fail the attempt,
   before running the Task's Verification commands, with
   `deletes: <path> still exists` per remaining path, through the existing
   failed-Verification route (API Contract 4).
4. MUST NOT add a Spec Consistency Check finding code or change
   `internal/speccheck/coherence.go`.
5. MUST add the tests named in Verification to the new files
   `internal/spec/context_deletes_test.go`,
   `internal/speccheck/context_deletes_test.go` and
   `internal/daemon/task_deletes_test.go`, and MUST leave the existing
   reference, declaration and recorded-path tests unedited and passing.
6. MUST teach the `deletes` kind and its settlement check in the write-tasks
   skill's Context rules and in its task template, raise the skill's version
   by one patch level in both front-matter fields, run `make skills-sync`,
   and re-record the version.

## Subtasks

- [ ] Parse the `deletes` kind.
- [ ] Count it in the reference check, the audits, the collision rule and the Task context.
- [ ] Refuse settlement while a `deletes` path exists.
- [ ] Add the tests.
- [ ] Teach the kind in the write-tasks skill and template, and raise its version.

## Acceptance Criteria

- [ ] `- deletes: old/file.txt` parses; the same path under `interface` and
      `deletes` is refused naming both kinds.
- [ ] A Task whose `deletes` path is gone raises no `SC-REF-UNRESOLVED`.
- [ ] A deleted path is not recorded under `## Recorded paths`; deleting a
      Governed Path counts as a governed touch; two Tasks that edit and
      delete one path cannot share a Wave.
- [ ] Over the task-cycle fixture, a Task whose Agent leaves its `deletes`
      path fails its attempt with `deletes: old/file.txt still exists`, and
      settles `completed` once the path is removed.
- [ ] The write-tasks skill and template teach the kind, the mirrors equal
      their canonical files, and the raised version is recorded.

## Context

- interface: `internal/spec/spec.go`
- interface: `internal/spec/task.go`
- interface: `internal/spec/recorded_paths.go`
- interface: `internal/spec/collision.go`
- interface: `internal/speccheck/citations.go`
- interface: `internal/speccheck/undeclared.go`
- interface: `internal/speccheck/surface.go`
- interface: `internal/daemon/task_context.go`
- interface: `internal/daemon/task_engine.go`
- creates: `internal/daemon/task_deletes.go`
- creates: `internal/daemon/task_deletes_test.go`
- creates: `internal/spec/context_deletes_test.go`
- creates: `internal/speccheck/context_deletes_test.go`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `.agents/skills/write-tasks/references/task-template.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `skills/write-tasks/references/task-template.md`
- interface: `skills/testdata/owned-skill-versions.json`
- instruction: `internal/speccheck/coherence.go`
- instruction: `docs/specs/0212-gates-and-run-storage-that-let-a-correct-delivery-finish/references/2026-10-01-task-context-has-no-kind-for-a-deleted-path.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestContextAcceptsADeletesEntry|TestContextRefusesAPathUnderTwoKinds|TestUndeclaredTaskPathsCountsADeletesPath|TestWaveCollisionCountsADeletesPath)$" ./internal/spec 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestContextAcceptsADeletesEntry TestContextRefusesAPathUnderTwoKinds TestUndeclaredTaskPathsCountsADeletesPath TestWaveCollisionCountsADeletesPath; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the four tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestDeletesPathIsNotAnUnresolvedReference|TestDeletingAGovernedPathNeedsItsGrant)$" ./internal/speccheck 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestDeletesPathIsNotAnUnresolvedReference TestDeletingAGovernedPathNeedsItsGrant; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the two tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestATaskThatLeavesItsDeletesPathDoesNotSettle|TestATaskThatRemovesItsDeletesPathSettles)$" ./internal/daemon 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestATaskThatLeavesItsDeletesPathDoesNotSettle TestATaskThatRemovesItsDeletesPathSettles; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the two tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestCheckReferenceUnresolved|TestContextDeclaredOutput|TestInstructionContextPathIsNotAudited|TestDeletesPathIsNotAnUnresolvedReference)$" ./internal/speccheck 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestCheckReferenceUnresolved TestContextDeclaredOutput TestInstructionContextPathIsNotAudited TestDeletesPathIsNotAnUnresolvedReference; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; the existing reference and declaration tests run unedited beside the new one, which does not exist before this Task.
- `tr -s '[:space:]' ' ' < .agents/skills/write-tasks/SKILL.md | grep -qF -- "- deletes: <path>" || { printf 'missing phrase in %s: %s\n' .agents/skills/write-tasks/SKILL.md "- deletes: <path>" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/write-tasks/SKILL.md | grep -qF -- "still exists" || { printf 'missing phrase in %s: %s\n' .agents/skills/write-tasks/SKILL.md "still exists" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/write-tasks/references/task-template.md | grep -qF -- "- deletes: <path>" || { printf 'missing phrase in %s: %s\n' .agents/skills/write-tasks/references/task-template.md "- deletes: <path>" >&2; exit 1; }; cmp .agents/skills/write-tasks/SKILL.md skills/write-tasks/SKILL.md && cmp .agents/skills/write-tasks/references/task-template.md skills/write-tasks/references/task-template.md && ! grep -q '^version: 0.0.7$' .agents/skills/write-tasks/SKILL.md && out="$(go test -count=1 -v -run "^(TestEveryOwnedSkillVersionIsRecorded|TestSettlementGuidanceIsOneTable|TestTaskAuthoringGuidanceNamesDeclarations)$" ./skills 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestEveryOwnedSkillVersionIsRecorded TestSettlementGuidanceIsOneTable TestTaskAuthoringGuidanceNamesDeclarations; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the write-tasks skill does not teach `- deletes: <path>`, so the command fails.

## References

- `_prd.md` → User Story 4; Core Feature 6; Success Metric 4; Acceptance evidence
- `_techspec.md` → The `deletes` Context kind; Interfaces; API Contract 4; Testing Approach 3; Build Order 3
- ADR-0166; ADR-0178


## Result

Implemented `ContextKindDeletes`, its parser and consumer support, the
pre-Verification deletion settlement check, and write-tasks guidance. Task
status and authored Verification remain Daemon-owned and unchanged.

### Acceptance evidence

1. `TestContextAcceptsADeletesEntry` accepts `old/file.txt` and rejects invalid
   paths. `TestContextRefusesAPathUnderTwoKinds` covers interface/deletes and
   creates/deletes in both orders, plus repeated deletes; diagnostics name
   both kinds. Before implementation, the focused parser test failed with
   `expected "instruction", "interface", or "creates" label`.
2. `TestDeletesPathIsNotAnUnresolvedReference` checks an absent deletion
   through the public Spec Consistency Check without `SC-REF-UNRESOLVED`.
   No finding code was added and `internal/speccheck/coherence.go` is untouched.
3. `TestUndeclaredTaskPathsCountsADeletesPath` excludes the declared deletion
   from paths eligible for `## Recorded paths`.
   `TestDeletingAGovernedPathNeedsItsGrant` refuses an ungranted deletion and
   accepts the matching grant and bounded rows.
   `TestWaveCollisionCountsADeletesPath` reports an edit/delete collision and
   accepts an explicit dependency. `deletes` entries join the Daemon bundle's
   interfaces. `TestDeletesCLISurfaceNeedsAGuideAndDeletionIsNotAGuide` proves
   a CLI deletion needs a guide and a deleted guide cannot satisfy that need.
4. `TestATaskThatLeavesItsDeletesPathDoesNotSettle` exercises the task-cycle
   fixture: both attempts fail with `deletes: old/file.txt still exists`,
   no authored Verification command runs, and no settlement commit occurs.
   `TestATaskThatRemovesItsDeletesPathSettles` removes the path during the
   bounded repair in the same Agent Session, then runs Verification and
   settles the fixture Task completed. The new `undeletedTaskPaths` uses
   `Lstat`; `TestUndeletedTaskPathsCountsDanglingSymlinksAndDirectories`
   covers directories, dangling symlinks and absent paths.
5. Canonical write-tasks skill and template teach `- deletes: <path>` and
   the `still exists` settlement diagnostic. Both skill version fields are
   `0.0.8`; `make skills-sync` regenerated the mirrors and the owned version
   record was regenerated. A Python byte comparison confirmed both canonical
   files equal their mirrors. Sanctioned digest regeneration reported no
   derived changes.

### Focused checks

- `GOCACHE=/tmp/roundfix-task03-go-cache rtk proxy go test ./internal/spec ./internal/speccheck ./internal/daemon ./skills -run 'Deletes|DeletingAGoverned|Undeleted|Context|ReferenceUnresolved|InstructionContextPath|RecordedPaths|Collisions|EveryOwnedSkillVersion|SettlementGuidance|TaskAuthoringGuidance' -count=1`
  — exit 0 in all four packages, including existing reference, declaration
  and recorded-path tests without edits.
- `rtk make skills-sync` — exit 0.
- `GOCACHE=/tmp/roundfix-task03-go-cache rtk proxy go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
  — exit 0; recorded write-tasks `0.0.8`.
- `GOCACHE=/tmp/roundfix-task03-go-cache rtk make baseline-digests`
  — exit 0, `changed: false`.
- `GOCACHE=/tmp/roundfix-task03-go-cache rtk make verify-incremental`
  — exit 0 with process-table permission. The sandboxed attempt passed the
  changed packages but failed two existing CLI force-stop tests because
  process-table reads were denied; the permission-enabled rerun passed.
- `git diff --check` — exit 0.

The Task's declared Verification commands were not run. No other Task or
Task Graph file was edited, and no commit, push or Pull Request was made.
No follow-up outside this slice was identified.

### Verification Feedback — attempt 1

Inspected the Daemon diagnostic artifact at
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261002T100707Z_cbb39973a4b07b9b/verification/batch-003-attempt-1.log`
and the related `internal/cli/implement_detach_teardown_test.go` test. The
Daemon-recorded `make verify-changed` failure was an environment limitation:
`TestImplementDetachChildEndsWhenItsTestBinaryDies` could not probe or kill
its disposable fixture's process group. The Daemon diagnostic recorded the
Task's spec, speccheck and daemon packages passing; those are recorded
Daemon observations, not an Agent rerun of configured Verification.

Focused check:
`GOCACHE=/tmp/roundfix-task03-go-cache rtk proxy go test ./internal/cli -run '^TestImplementDetachChildEndsWhenItsTestBinaryDies$' -count=1 -v`
— exit 0 with process-control permission; the named test passed and the
suite guard reported no repository mutation.

No implementation repair was indicated by this feedback, so Task code and
existing CLI tests remain unchanged. The Daemon retry needs permission to
probe and stop test-owned process groups. Changing CLI teardown behavior,
Verification configuration or sandbox policy is outside task_03's slice.
The Agent did not rerun `make verify-changed` or any declared Verification
command. Task status is unchanged; no commit, push or Pull Request was made.

## Carry-forward provenance

- Source Run: `run_20261002T100707Z_cbb39973a4b07b9b`
- Source commit: `d43995132a7c73feabc8b9076fcf710ddffee035`
