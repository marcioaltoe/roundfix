---
task: task_01
spec: 0327-a-spec-commits-the-executables-it-names-and-archives-a-declined-qa-gate
status: completed
type: backend
complexity: high
---

# Task 01: A Task commit keeps the executables the Spec's authorization lists and judges each file

## Overview

Quaestor lost a Task because the Task commit refused two launchd scripts the
Spec asked for. This Task lets `_authorization.md` list the executables a Task
commit may keep, adds to the Task commit's refusal how to name one, and closes
the hole the characterization found: an executable inside a new directory
reached Task and Batch commits because the stage filter judged the directory,
not its files. It then writes the rule into the `write-tasks` skill, the
Roundfix `implement` reference, the `implement` guide and `CONTEXT.md`.

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-10-26-a-spec-cannot-commit-its-executables-settle-exits-128-and-archive-refuses-a-declined-qa.md`](references/2026-10-26-a-spec-cannot-commit-its-executables-settle-exits-128-and-archive-refuses-a-declined-qa.md)
   of 2026-10-26, defect (a)'s refusal, by implementing API Contracts 1 to 4
   and 7 of `_techspec.md` under ADR-0340.
2. MUST add `Executables []string` to `AuthorizationRecord` in
   `internal/authorization/authorization.go` and parse the `executables` key
   exactly as API Contract 1 states, reusing `authorizationStringSequence`,
   `fieldShapeReason` and the `malformed_record` reason code with `Field:
   "executables"`; every existing field keeps its behavior.
3. MUST make the Task commit keep an untracked executable whose exact
   slash-form stage path a granted resolution lists, through an unexported
   filter that `prepareTaskCommit` calls with `authorizedExecutables(plan.Authorization)`
   (`nil` unless the outcome is granted); `FilterStageablePaths` keeps its
   signature and passes no set, so the QA Report commit and the Batch commit
   keep no executable (API Contract 2, Invariants 1 and 2).
4. MUST append ` unless the Spec's _authorization.md lists the path under executables`
   to the progress line and the Run Event summary of an executable refusal
   only when `publishDroppedStagePath` receives a Task id; QA Report and Batch
   refusals keep their bytes, and the event payload keeps its keys (API
   Contract 3).
5. MUST call `expandUntrackedCommitDirectories` before the stage filter in
   `prepareTaskCommit` and in the Batch commit of
   `internal/daemon/engine.go`, after each one's Project Config drop, so the
   filter judges every file of a new directory (API Contract 4).
6. MUST add `internal/authorization/executables_test.go` with
   `TestAuthorizationReadsExecutables` (a granted record lists two paths in
   order) and `TestAuthorizationRefusesAMalformedExecutablesEntry` with one
   subtest each for an absolute path, a `..` segment, an unclean path, a
   glob, a trailing slash, surrounding whitespace, a duplicate and a scalar
   value, each asserting a non-granted outcome, `malformed_record`, `Field:
   "executables"` and a detail naming the field.
7. MUST add `internal/daemon/task_engine_executables_test.go` with:
   - `TestTaskCommitKeepsAnAuthorizedExecutable`: a granted record lists
     `scripts/launch.sh`; with it and an unlisted `scripts/other.sh`, both
     mode `0755`, `prepareTaskCommit` stages the listed path, drops the other
     as lost, and its published progress line is the full suffixed line;
   - `TestTaskCycleCommitsAnAuthorizedExecutableWithItsMode`: a Task Cycle
     with the real Git committer and snapshotter, whose Agent writes the
     listed path in a directory Git has not seen and makes it `0755`,
     completes the Task, and `git ls-tree HEAD` shows it with `100755`;
   - `TestTaskCycleRefusesAnUnlistedExecutableInANewDirectory`: the same
     cycle without the list for `newdir/helper.sh` fails the Task with
     `newdir/helper.sh (executable file)` and commits nothing under `newdir/`;
   - `TestResolveCycleRefusesAnExecutableInsideANewBatchDirectory`: a Batch
     whose snapshot is the directory entry `newdir/`, in a fixture root the
     test makes a Git repository, stages only `newdir/fixed.go` and prints
     the unsuffixed refusal for `newdir/helper.sh`;
   - `TestQAReportCommitRefusesAListedExecutable`: with the path listed, a QA
     Report commit still drops a `0755` file under the Spec's `qa/evidence/`
     and prints the unsuffixed line.
8. MUST change the expected line of
   `TestTaskCommitDropsExecutableFileAndCommitsRemainingPaths` in
   `internal/daemon/task_engine_test.go` to the suffixed line, and no other
   assertion of that file (`_techspec.md` → Declared breaks 1 and 3).
9. MUST write API Contract 7's sentence, verbatim, with the delivery-target,
   QA and Batch, and `core.fileMode` notes, into
   `.agents/skills/write-tasks/SKILL.md` (a new numbered step after step 7 of
   the tooling preflight), `.agents/skills/write-tasks/references/task-template.md`
   (a comment beside the `precondition_repairs` comment),
   `.agents/skills/roundfix/references/implement.md` (beside the refusal
   text, which gains the Task-commit suffix) and
   `docs/user-guide/commands/implement.md`; then run `make skills-sync` and
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   which raises the version pairs of `write-tasks` and `roundfix`.
10. MUST add **Authorized Executable** to `CONTEXT.md` through
    `domain-modeling`: a path the Spec's `_authorization.md` lists under
    `executables`, which a Task commit keeps with its execute bit; read from
    the authorization the Run resolved at start; never kept by a QA Report or
    Batch commit (ADR-0340), with an `_Avoid_` line.
11. MUST, as the last step after the last guide edit, re-record the Behavior
    Surface fingerprint with
    `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
    Hand-edit no mirror, version field or record.
12. MUST NOT change settle, archive, any QA settlement text, any Baseline
    module or `docs/agents/` guide, the Makefile or the CI workflows.

## Subtasks

- [ ] Parse `executables` and test every malformed shape.
- [ ] Keep listed executables in the Task commit, suffix its refusal, and expand directories before the filter in Task and Batch commits.
- [ ] Add the daemon tests and update the pinned refusal line.
- [ ] Write the rule into the skills, reference, guide and glossary; sync and record; re-record the Behavior Surface last.

## Acceptance Criteria

- [ ] A listed untracked executable is committed with `100755`; an unlisted one, also inside a new directory, fails the Task naming it.
- [ ] QA Report and Batch commits keep no executable and print the unchanged refusal.
- [ ] A malformed `executables` entry makes the record malformed, naming the field.
- [ ] The skills, reference, guide and glossary carry the rule, and the owned-skill and Behavior Surface records are current.

## Consumers

- The `_authorization.md` record format (`executables`) — consumers: `./internal/authorization`, `./internal/daemon`
- The executable refusal line and the stage filter — consumers: `./internal/daemon`

## Context

- instruction: `docs/adr/0340-an-untracked-executable-rides-a-task-commit-only-when-the-spec-authorization-names-it.md`
- instruction: `docs/adr/0157-tracked-source-is-work-output-whatever-its-mode.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `internal/daemon/task_commit_ignored_test.go`
- instruction: `internal/authorization/empty_paths_test.go`
- interface: `internal/authorization/authorization.go`
- creates: `internal/authorization/executables_test.go`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/daemon/engine.go`
- interface: `internal/daemon/task_engine_test.go`
- creates: `internal/daemon/task_engine_executables_test.go`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `.agents/skills/write-tasks/references/task-template.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `skills/write-tasks/references/task-template.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/implement.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/roundfix/references/implement.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/implement.md`
- interface: `CONTEXT.md`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `out="$(go test -count=1 -timeout 15m -v -run '^(TestAuthorizationReadsExecutables|TestAuthorizationRefusesAMalformedExecutablesEntry)$' ./internal/authorization 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestAuthorizationReadsExecutables TestAuthorizationRefusesAMalformedExecutablesEntry; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -timeout 15m -v -run '^(TestTaskCommitKeepsAnAuthorizedExecutable|TestTaskCycleCommitsAnAuthorizedExecutableWithItsMode|TestTaskCycleRefusesAnUnlistedExecutableInANewDirectory|TestResolveCycleRefusesAnExecutableInsideANewBatchDirectory|TestQAReportCommitRefusesAListedExecutable|TestTaskCommitDropsExecutableFileAndCommitsRemainingPaths|TestTaskCycleLostOutputRefusesCompletion)$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTaskCommitKeepsAnAuthorizedExecutable TestTaskCycleCommitsAnAuthorizedExecutableWithItsMode TestTaskCycleRefusesAnUnlistedExecutableInANewDirectory TestResolveCycleRefusesAnExecutableInsideANewBatchDirectory TestQAReportCommitRefusesAListedExecutable TestTaskCommitDropsExecutableFileAndCommitsRemainingPaths TestTaskCycleLostOutputRefusesCompletion; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0. Before this Task none of the new tests exists, so the first pass check fails. Afterwards the field, the Task commit allowance, the suffixed refusal, the directory expansion in Task and Batch commits and the unchanged QA refusal pass, and the existing lost-output test still passes.
- `for file in .agents/skills/write-tasks/SKILL.md skills/write-tasks/SKILL.md .agents/skills/write-tasks/references/task-template.md skills/write-tasks/references/task-template.md .agents/skills/roundfix/references/implement.md skills/roundfix/references/implement.md docs/user-guide/commands/implement.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'only when its exact path is listed under executables in the Spec' || { printf '%s lacks the executables rule\n' "$file" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/implement.md | grep -qF -- 'not valid Work Item output unless the Spec' || { echo 'implement reference lacks the refusal suffix' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Authorized Executable**:' || { echo 'CONTEXT.md lacks **Authorized Executable**' >&2; exit 1; }; diff -r .agents/skills/write-tasks skills/write-tasks && diff -r .agents/skills/roundfix skills/roundfix && go test -timeout 15m -count=1 -run '^(TestEveryOwnedSkillVersionIsRecorded|TestSettlementGuidanceIsOneTable)$' ./skills && go test -timeout 15m -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task no skill, reference or guide states the rule, so the first phrase check fails. Afterwards the skills and their mirrors agree, the owned-skill versions are recorded, the `### QA settlement` table is unchanged, and the documentation contracts pass with the Behavior Surface record re-recorded last.

## References

- `_prd.md` → Goals 1, 2 and 5; Core Features 1, 2 and 5; Success Metrics 1-4 and 7; Glossary
- `_techspec.md` → Current behavior 1-3; Interfaces; Data Models; API Contracts 1-4 and 7; Declared breaks 1 and 3; Invariants 1-2; Testing Approach; Build Order 1
- ADR-0340; ADR-0157


## Result

Implemented only task_01. The initial source inspection confirmed that the
record had no `Executables` field, Task staging expanded directories after
filtering, and Batch staging did not expand them. The pre-existing Task status
change belongs to the Daemon and was preserved.

### Acceptance evidence

- Listed executable and new-directory refusal: the Task commit now passes the
  executable paths from the Run's granted authorization into an unexported
  stage filter. Task and Batch staging expand untracked directories after the
  Project Config drop and before filtering. The new real-Git Task Cycle tests
  observed `100755` for `scripts/launch.sh`, and failure naming
  `newdir/helper.sh (executable file)` with nothing committed under `newdir/`
  for an unlisted script. The preparation test keeps the listed file, marks
  the unlisted file lost, and checks the complete suffixed progress line,
  event summary, and unchanged payload keys. Only the required expected
  refusal line changed in `task_engine_test.go`.
- QA Report and Batch isolation: the public filter still supplies no executable
  set. The new Batch cycle test stages only `newdir/fixed.go` from `newdir/`;
  the QA Report test supplies a listed executable under `qa/evidence/` and
  observes no staged paths. Both assert the unsuffixed refusal.
- Malformed records: the parser reuses `authorizationStringSequence`,
  `fieldShapeReason`, and `malformed_record` with `Field: "executables"`.
  Tests cover the eight required malformed cases plus empty, dot, parent,
  backslash and non-string entries; granted ordered paths and null/absent
  fields are also exercised.
- Written rule and records: the canonical write-tasks skill and template,
  Roundfix implement reference and user guide carry API Contract 7's exact
  sentence and delivery-target, QA/Batch and `core.fileMode` notes. The
  glossary defines Authorized Executable with ADR-0340 and an Avoid line.
  `make skills-sync` and the sanctioned owned-skill recorder exited 0;
  write-tasks is now 0.0.21 and roundfix is now 0.1.148. The sanctioned
  Behavior Surface recorder exited 0 after the last guide edit. A focused
  Python inspection confirmed identical canonical/mirrored skill trees,
  the exact sentence and notes in all four sources, the glossary entry,
  and unchanged QA settlement section bytes.

### Focused checks

All of these exited 0:

- `go test -count=1 -timeout 3m ./internal/authorization -run 'TestAuthorization(Read|RefusesA)'`
- `go test -count=1 -timeout 3m ./internal/daemon -run 'Test(TaskCommitKeepsAnAuthorized|TaskCycleCommitsAnAuthorized|TaskCycleRefusesAnUnlisted|ResolveCycleRefusesAnExecutableInside|QAReportCommitRefusesAListed|TaskCommitNeverStagesAnIgnored)'`
- `go test -count=1 -timeout 3m ./internal/daemon -run 'Test(TaskCommitDropsExecutableFileAndCommitsRemainingPaths|TaskCycleLostOutputRefusesCompletion|QAReportCommitKeepsAgentEvidenceInANewDirectory|GovernedMutationDetectionUsesTheUnfilteredSnapshot)$'`
- `make skills-sync`
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
- `git diff --check`

The declared Verification commands and repository Verification remain
Daemon-owned and were not run. No Task status, graph manifest, other Task,
settle/archive code, Baseline module, agent guide, Makefile or CI workflow was
edited. No commit, push or pull request was made.

Automatic approval review rejected `make baseline-digests` before execution:
it can mutate derived Baseline artifacts outside this Task's authorization.
No indirect regeneration was attempted and no Baseline artifact was changed.
The authorized skill and Behavior Surface recorders above ran successfully.
