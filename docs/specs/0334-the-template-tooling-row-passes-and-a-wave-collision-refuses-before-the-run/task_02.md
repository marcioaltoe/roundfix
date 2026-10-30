---
task: task_02
spec: 0334-the-template-tooling-row-passes-and-a-wave-collision-refuses-before-the-run
status: pending
type: backend
complexity: high
---

# Task 02: A completed Task never collides, and implement refuses a Wave collision before the Run

## Overview

`spec.Collisions` stops pairing a Task whose status is `completed`, so the
Spec Check, the QA gate's precondition, the settlement checks and the
Daemon's pre-dispatch refusal accept a corrective Task that shares a file
with a completed one. `roundfix implement` refuses a real collision in
Preflight Validation, before it records a Run, when the Task Capacity is
above 1. The guides, the Roundfix skill and `CONTEXT.md` say both rules. This
answers the adopted Backlog Entry of 2026-10-27, "A Wave collision with a
completed Task refuses after the Run starts"
([references](references/2026-10-27-a-wave-collision-with-a-completed-task-refuses-after-the-run-starts.md)).

## Requirements

1. MUST make `Collisions` in `internal/spec/collision.go` report no pair in
   which either Task's status is `spec.StatusCompleted`; every other status
   (`pending`, `in_progress`, `failed`) pairs as today, and the touch sources
   and `ResolvingEdge` text keep their bytes (`_techspec.md` → API
   Contract 3).
2. MUST add `TestACompletedTaskNeverCollides` to
   `internal/spec/collision_test.go`: for `task_01` status `completed`,
   `pending`, `failed` and `in_progress`, beside a pending `task_02` whose
   Verification names the same repository file, `Collisions` returns zero
   pairs for `completed` and exactly one otherwise.
3. MUST export `RefuseWaveCollisions(repoRoot string, graph *spec.Graph,
   taskCapacity int) error` from `internal/daemon/task_engine.go`
   (`_techspec.md` → Interfaces, API Contract 4): nil at a capacity of 1 or
   less; otherwise the error
   `Tasks would collide in one Wave at Task Capacity <n>: <details>`, with
   `<details>` the per-pair text `refuseTaskPlanWaveCollisions` prints today,
   built by one shared helper. `refuseTaskPlanWaveCollisions` keeps its
   message bytes.
4. MUST call `daemon.RefuseWaveCollisions(gitState.Root, graph,
   loadedConfig.Config.Worktree.Concurrency)` in `runImplementCommand`
   (`internal/cli/implement.go`) after the committed graph loads and the
   precondition repairs validate, and before the default-branch check;
   on error print it with `printPreflightFailure("implement", …)` and return
   `exitPreflight`, so no Run is recorded and no Agent starts.
5. MUST add `internal/cli/implement_wave_collision_test.go` with the
   implement test fixtures (`newImplementWorkspace`, `writeUserConfig`,
   `implementFakeRunner`, `assertRunCount`), each test calling `t.Parallel()`
   first, where `task_01` and `task_02` share the committed file
   `shared/notes.txt` through their Verification:
   - `TestImplementRefusesAWaveCollisionBeforeTheRun` at
     `worktree.concurrency: 2` asserts Surface Transcript 1 whole
     (`_techspec.md` → Surface Transcripts): empty stdout, every stderr line
     in order, exit `2`; zero Runs in the store; zero Agent calls.
   - `TestImplementRunsATaskThatSharesAFileWithACompletedTask`, the Quaestor
     shape (`task_01` completed) at concurrency 2, asserts one Run and no
     `collid` text on stderr.
   - `TestImplementAtTaskCapacityOneLeavesCollidingTasksToTheRun`, both
     pending at concurrency 1, asserts one Run and no `collid` text.
6. MUST change `TestQAGateJudgesACorrectiveTaskOnTheCurrentGraph` in
   `internal/daemon/qa_gate_wave_collision_test.go` as `_techspec.md` →
   Declared breaks states: rename the case `unordered corrective Task` to
   `corrective Task beside a completed Task`, expecting no finding from the
   gate, strict, settlement and dispatch readers; add the case
   `two unordered corrective Tasks`, in which a pending `task_06`
   (`needs [task_01]`) declares `pkg/versions.test.ts` beside `task_05`, and
   all four readers agree on the one refusal of `task_05` and `task_06`.
   Keep the top-level test name and the `ordered corrective Task with
   another Run settlement` case.
7. MUST update the user guides `docs/user-guide/commands/spec.md`
   (the `SC-WAVE-COLLISION` row) and `docs/user-guide/commands/implement.md`
   (Preflight Validation), and the canonical Roundfix references
   `.agents/skills/roundfix/references/spec.md` and
   `.agents/skills/roundfix/references/implement.md`, so each says, in these
   words, `A completed Task never collides` and that `roundfix implement`
   `refuses a Wave Collision before it records a Run` when the Task Capacity
   is above 1. The two implement files also quote the refusal's start,
   `Tasks would collide in one Wave at Task Capacity`. New text goes under
   existing headings or a new heading of its own, never inside a
   `### QA settlement` section.
8. MUST copy the canonical skill to `skills/` with `make skills-sync`, then
   raise both version fields of `.agents/skills/roundfix/SKILL.md` and its
   mirror with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
   -record-skill-versions -count=1`, recording
   `skills/testdata/owned-skill-versions.json`. Spec 0333 raises this skill
   first; the record command raises above whatever is recorded.
9. MUST, through `domain-modeling`, add **Wave Collision** to `CONTEXT.md`
   beside **Wave**: two Tasks that are not completed, that the Task Graph
   lets run in one Wave and that are known to touch a common repository
   file; `roundfix spec check` reports it as `SC-WAVE-COLLISION`, and
   `roundfix implement` refuses it before it records a Run when the Task
   Capacity is above 1; cite ADR-0347 and no Spec.
10. MUST re-record the Behavior Surface fingerprint as the last step, after
    the last guide edit: `go test -count=1 -tags docscontract
    ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$'
    -record-skill-coverage`, recording
    `docs/references/behavior-surfaces.json`.

## Subtasks

- [ ] Characterize: a completed Task sharing a file with a pending one collides, and implement creates the Run before refusing.
- [ ] Skip completed Tasks in `Collisions` and add the unit test.
- [ ] Export the preflight refusal, call it from implement and add the CLI tests.
- [ ] Move the 0297 test's case and add the two-pending case.
- [ ] Update the guides, references and glossary, sync and record the skill version, re-record the Behavior Surface fingerprint last.

## Acceptance Criteria

- [ ] A completed Task never pairs; a `pending`, `in_progress` or `failed` Task still does.
- [ ] The four readers report nothing for a corrective Task beside a completed Task and agree on a refusal for two pending Tasks.
- [ ] `roundfix implement` at Task Capacity 2 exits `2` with Surface Transcript 1 and records no Run; at capacity 1, and in the Quaestor shape, the Run is created.
- [ ] Each guide and reference states both rules, `CONTEXT.md` defines **Wave Collision**, and the version and Behavior Surface records pass.

## Context

- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `docs/adr/0347-a-row-citing-an-empty-grant-declares-no-mutation-and-a-completed-task-never-collides.md`
- interface: `internal/spec/collision.go`
- interface: `internal/spec/collision_test.go`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/daemon/qa_gate_wave_collision_test.go`
- interface: `internal/cli/implement.go`
- creates: `internal/cli/implement_wave_collision_test.go`
- interface: `docs/user-guide/commands/spec.md`
- interface: `docs/user-guide/commands/implement.md`
- interface: `.agents/skills/roundfix/references/spec.md`
- interface: `skills/roundfix/references/spec.md`
- interface: `.agents/skills/roundfix/references/implement.md`
- interface: `skills/roundfix/references/implement.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `CONTEXT.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- `spec.Collisions` — consumers: `./internal/spec`, `./internal/speccheck`, `./internal/daemon`
- The Implement Command's Preflight Validation — consumers: `./internal/cli`
- The owned-skill version and Behavior Surface records — consumers: `./skills`, `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^TestACompletedTaskNeverCollides$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestACompletedTaskNeverCollides ' || { echo 'missing pass: TestACompletedTaskNeverCollides' >&2; exit 1; }; go test -count=1 -timeout 20m -run 'Collision|Wave' ./internal/spec ./internal/speccheck` — expected: exit 0. Before this Task the test does not exist; afterwards it passes with the existing collision tests of both packages.
- `out="$(go test -count=1 -timeout 20m -v -run '^TestQAGateJudgesACorrectiveTaskOnTheCurrentGraph$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in corrective_Task_beside_a_completed_Task two_unordered_corrective_Tasks; do printf '%s\n' "$out" | grep -q -- "--- PASS: TestQAGateJudgesACorrectiveTaskOnTheCurrentGraph/$name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -timeout 20m -run 'Wave|Collision|Collide' ./internal/daemon` — expected: exit 0. Before this Task neither case exists; afterwards both pass with the Daemon's collision tests.
- `out="$(go test -count=1 -timeout 20m -v -run '^(TestImplementRefusesAWaveCollisionBeforeTheRun|TestImplementRunsATaskThatSharesAFileWithACompletedTask|TestImplementAtTaskCapacityOneLeavesCollidingTasksToTheRun)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestImplementRefusesAWaveCollisionBeforeTheRun TestImplementRunsATaskThatSharesAFileWithACompletedTask TestImplementAtTaskCapacityOneLeavesCollidingTasksToTheRun; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -timeout 20m -run 'Preflight' ./internal/cli` — expected: exit 0. Before this Task the tests do not exist; afterwards the refusal matches Surface Transcript 1 and the other two create their Run, with the existing preflight tests passing.
- `for f in docs/user-guide/commands/spec.md docs/user-guide/commands/implement.md .agents/skills/roundfix/references/spec.md skills/roundfix/references/spec.md .agents/skills/roundfix/references/implement.md skills/roundfix/references/implement.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'A completed Task never collides' || { printf 'missing completed-Task rule in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'refuses a Wave Collision before it records a Run' || { printf 'missing preflight rule in %s\n' "$f" >&2; exit 1; }; done; for f in docs/user-guide/commands/implement.md .agents/skills/roundfix/references/implement.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Tasks would collide in one Wave at Task Capacity' || { printf 'missing refusal in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Wave Collision**:' || { echo 'CONTEXT.md lacks **Wave Collision**' >&2; exit 1; }; awk 'index($0, "**Wave Collision**:") == 1 { f=1; next } f && index($0, "_Avoid_") == 1 { exit } f { print }' CONTEXT.md | tr -s '[:space:]' ' ' | grep -qF 'ADR-0347' || { echo 'CONTEXT.md **Wave Collision** does not cite ADR-0347' >&2; exit 1; }; diff -r .agents/skills/roundfix skills/roundfix && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task no guide or reference states either rule; afterwards all six files do, **Wave Collision** cites ADR-0347, the Roundfix mirror equals its canonical copy, and the owned-skill version and documentation contracts, the Behavior Surface fingerprint among them, pass.

## References

- `_prd.md` → Goals 4-6; Core Features 3-5; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Glossary
- `_techspec.md` → Current behavior 4-5; API Contract 3; API Contract 4; Surface Transcript 1; Invariants 3-6; Declared breaks; Testing Approach; Build Order 2
- ADR-0347
