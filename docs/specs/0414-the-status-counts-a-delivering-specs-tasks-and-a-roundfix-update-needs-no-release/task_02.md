---
task: task_02
spec: 0414-the-status-counts-a-delivering-specs-tasks-and-a-roundfix-update-needs-no-release
status: completed
type: backend
complexity: high
---

# Task 02: In an adopter, a change to the files Roundfix manages or to a CONTEXT.md needs no release

## Overview

In an adopter, `roundfix release plan` asks for a classification by hand for every Roundfix update, every adoption of Roundfix recommendations and every plan that touches a `CONTEXT.md`: the Fiscus range `v0.9.0..ad1381b` blocked on four commits, three of which changed only such files (the maintainer, 2026-11-19, [`references/2026-11-19-maintainer-answers.md`](references/2026-11-19-maintainer-answers.md); measured by the operator, [`references/2026-11-19-characterization.md`](references/2026-11-19-characterization.md), sections 2 and 3). This Task makes the Release Plan count every Adopter Maintenance Path as maintenance-only when the target's `skills-lock.json` records the `roundfix` skill, as ADR-0447 decides, and states the rule in the `release plan` help, the release runbook, the usage guide, both copies of the Roundfix skill's release reference and the glossary.

## Requirements

1. MUST add `internal/releaseplan/managed.go` with `ManagedPathSource`, `ManagedPaths` and `ManagedPaths.Contains` as `_techspec.md` → Interfaces, invariants 3 and 5 and API Contract 2 state, add the field `Managed ManagedPaths` to `ClassifyRequest`, and make the maintenance-only boundary of `ClassifyChanges` and of `Build`'s release-preparation pre-check ask `Managed.Contains` about each path `MaintenanceOnly` rejects. `ClassifyCommit` keeps its signature and its result. `Build` fills `Managed` from the source when the source implements `ManagedPathSource`, before it classifies, and fails the plan with the source's error when that read fails.
2. MUST add `releasePlanGitSource.ManagedPaths` in a new `internal/cli/releaseplan_managed_paths.go`, following invariant 4: `git ls-tree` and `git show` of `skills-lock.json` and `docs/agents/setup-context.json` at the target and at the base, `Adopter` from the target's lock only, `Skills` from the keys of each lock's `skills` object, `Files` from each manifest's `managedArtifacts[].path`. It reads no working-tree file and never `.roundfixrc.yml`; a failing `git` command returns a `releaseplan.GitSourceError`.
3. MUST add `TestAnAdopterManagedPathIsMaintenanceOnly` to a new `internal/releaseplan/managed_test.go`, calling `t.Parallel()` first, classifying one commit per subtest through `ClassifyChanges` with an adopter whose `Files` is `AGENTS.md` and `docs/agents/backend.md` and whose `Skills` is `roundfix` and `write-tasks`, and asserting `CrossesMaintenanceOnlyBoundary`:
   - false for `the_setup_manifest` (`docs/agents/setup-context.json`), `a_recorded_guide` (`docs/agents/backend.md`), `skills_lock` (`skills-lock.json`), `a_locked_skill_file` (`.agents/skills/write-tasks/references/qa.md`), `a_claude_skill_link` (`.claude/skills/roundfix`), `claude_md` (`CLAUDE.md`), `specific_repository` (`docs/agents/specific-repository.md`), `a_generated_author_subagent` (`.claude/agents/roundfix-tasks-author.md`), `project_config` (`.roundfixrc.yml`) and `a_nested_context_md` (`packages/backend/src/contexts/tax/CONTEXT.md`);
   - true for `an_unlocked_local_skill_crosses` (`.agents/skills/revisao/scripts/run.ts`), `another_claude_agent_crosses` (`.claude/agents/reviewer.md`), `context_map_crosses` (`CONTEXT-MAP.md`), `an_unrecorded_guide_crosses` (`docs/agents/frontend.md`), `a_mixed_commit_crosses` (`AGENTS.md` and `packages/backend/.env.example`), `not_an_adopter_changes_nothing` (`CLAUDE.md` with `Adopter` false) and `an_empty_commit_still_crosses` (no changed path).

   Every existing test of `./internal/releaseplan` stays unchanged and keeps passing.
4. MUST add `TestTheReleasePlanReadsTheManagedRecordsAtBaseAndTarget` to a new `internal/cli/releaseplan_managed_paths_test.go`, calling `t.Parallel()` first. Each subtest builds a temporary Git repository with its own stable tag `v0.1.0` on a commit that holds a Setup Manifest recording `AGENTS.md` and `docs/agents/backend.md` and a `skills-lock.json`, then commits, and runs `roundfix release plan --from v0.1.0 --format json` through `runReleasePlanCommand`:
   - `the_fiscus_shape`: the lock records `roundfix` and `write-tasks`; four `chore:` or `docs:` commits shaped like `references/2026-11-19-characterization.md` section 2: one changing `.agents/skills/write-tasks/SKILL.md`, `AGENTS.md`, `docs/agents/setup-context.json` and `skills-lock.json`; one changing `.claude/agents/roundfix-tasks-author.md`, `.roundfixrc.yml` and `docs/agents/specific-repository.md`; one changing `docs/specs/0097-widget/_prd.md` and `packages/backend/src/CONTEXT.md`; one changing `docs/specs/0099-move/_prd.md`, `docs/user-guide/move.md` and `packages/backend/.env.example`. The plan is `manual_classification_required` and its only blocking commit is the fourth;
   - `roundfix_itself`: the same commits with a lock that records only `write-tasks`; all four commits block;
   - `a_skill_removed_in_the_range`: the lock records `roundfix` and `old-skill` at the tag; one commit removes `.agents/skills/old-skill/SKILL.md` and the lock entry; the plan is `no_release`;
   - `a_malformed_lock_makes_no_adopter`: the target's lock is not JSON; the commit changing `CLAUDE.md` blocks;
   - `project_config_is_not_read`: an adopter whose target holds a `.roundfixrc.yml` that is not valid YAML; the commit changing only that file and a `CONTEXT.md` ends `no_release` and the command exits 0.
5. MUST add to the `release plan` help in `internal/cli/cli.go` the paragraph of API Contract 2, word for word apart from line wrapping, right after the paragraph that begins "Planning evidence (docs/specs", and add `TestTheReleasePlanHelpNamesTheAdopterMaintenancePaths` to `internal/cli/releaseplan_command_test.go`, calling `t.Parallel()` first, with one subtest per phrase on whitespace-normalized help text: `adopter_test` ("a repository whose skills-lock.json at the target records the roundfix skill"), `setup_manifest` ("the Setup Manifest docs/agents/setup-context.json and every path it records"), `skills_lock` ("skills-lock.json and every skill directory it records under .agents/skills/ and .claude/skills/"), `claude_md_and_specific_repository` ("CLAUDE.md, docs/agents/specific-repository.md"), `generated_subagents` ("a generated .claude/agents/roundfix-*.md"), `project_config_and_context` (".roundfixrc.yml and every CONTEXT.md"), `base_and_target` ("the two records are read at the base and at the target") and `roundfix_itself` ("Roundfix's own repository authors the roundfix skill and records none"). `TestReleasePlanHelpDescribesFlagsDefaultsStatesAndReadOnlyBoundary` stays unchanged and keeps passing.
6. MUST append the text of API Contract 2 for the release documents, word for word apart from line wrapping, to the paragraph that begins "Backlog Entries under `docs/backlog/`" in `.agents/skills/roundfix/references/release.md`, `docs/user-guide/release-runbook.md` and `docs/user-guide/usage.md`, after "a commit to classify by hand.", then run `make skills-sync` so `skills/roundfix/references/release.md` matches.
7. MUST add `TestTheReleaseDocumentsNameTheAdopterMaintenancePaths` to a new `skills/release_adopter_maintenance_test.go`, calling `t.Parallel()` first, reading the canonical and embedded `roundfix/references/release.md` through `roundfixLayoutSources` and the two guides from `../docs/user-guide/`, with subtests named `<document>/<check>` for the documents `canonical`, `embedded`, `runbook` and `usage` and the checks `adopter_test`, `setup_manifest`, `skills_lock`, `claude_md_and_specific_repository`, `generated_subagents`, `project_config_and_context`, `base_and_target` and `roundfix_itself`, each asserting its phrase of API Contract 2 on whitespace-normalized text, so removing any one phrase from any one document fails exactly its subtest.
8. MUST raise the Roundfix skill to the next version in both `version` fields of `.agents/skills/roundfix/SKILL.md`, run `make skills-sync`, and re-record `skills/testdata/owned-skill-versions.json` with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` on the tree this Task runs on, never by hand. No other Roundfix skill file changes.
9. MUST, through the `domain-modeling` skill, add the glossary entry **Adopter Maintenance Path** to `CONTEXT.md`, right after the entry for the Release Plan, with the definition, the ADR-0447 link and the `_Avoid_` line of API Contract 2, word for word. No other glossary entry changes in this Task.
10. MUST, as its last edit, after the last change to a guide, re-record the Behavior Surface Record with `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`, never by hand.
11. MUST NOT change the status view, `docs/references/skill-coverage.json`, `docs/references/coverage-record.json`, any other owned skill, any Vendored Skill, any Baseline file, this repository's `.roundfixrc.yml` or `skills-lock.json`, the Makefile, a workflow, `go.mod` or a lint configuration.

## Subtasks

- [ ] Add `ManagedPaths` to the classifier and the `ManagedPathSource` seam to `Build`, with `TestAnAdopterManagedPathIsMaintenanceOnly`.
- [ ] Read the two records at the base and the target in the Git source, with `TestTheReleasePlanReadsTheManagedRecordsAtBaseAndTarget`.
- [ ] State the rule in the help, the release reference and the two guides, pinned by their contract tests.
- [ ] Raise the Roundfix skill version, add the glossary entry and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] In an adopter, a commit confined to Adopter Maintenance Paths and planning evidence crosses no boundary, and a commit with any other path, an empty commit, an unlocked skill, another `.claude/agents` file, `CONTEXT-MAP.md` or an unrecorded guide still crosses it.
- [ ] A repository whose lock records no `roundfix` skill, or whose lock is malformed, classifies as before, and the plan reads no `.roundfixrc.yml`.
- [ ] The help, both copies of the release reference and both guides name every kind of path, the adopter test and Roundfix's own exception, each pinned by its own subtest.
- [ ] The Roundfix skill version is raised and recorded, the glossary holds the new entry, and the Behavior Surface Record is current.

## Context

- instruction: `docs/adr/0447-in-an-adopter-a-change-to-the-files-roundfix-manages-or-to-a-context-md-needs-no-release.md`
- instruction: `docs/adr/0412-planning-references-and-the-release-preparation-commit-set-no-release-impact.md`
- instruction: `internal/cli/releaseplan_git_source.go`
- instruction: `skills/roundfix_layout_test.go`
- creates: `internal/releaseplan/managed.go`
- creates: `internal/releaseplan/managed_test.go`
- interface: `internal/releaseplan/classifier.go`
- interface: `internal/releaseplan/build.go`
- interface: `internal/releaseplan/model.go`
- creates: `internal/cli/releaseplan_managed_paths.go`
- creates: `internal/cli/releaseplan_managed_paths_test.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/releaseplan_command_test.go`
- interface: `docs/user-guide/release-runbook.md`
- interface: `docs/user-guide/usage.md`
- interface: `.agents/skills/roundfix/references/release.md`
- interface: `skills/roundfix/references/release.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- creates: `skills/release_adopter_maintenance_test.go`
- interface: `CONTEXT.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The Release Plan classification (`ClassifyRequest`, `Build` and the Git source) — consumers: `./internal/releaseplan`, `./internal/cli`
- The release documents, the Roundfix skill, its version record, the glossary and the Behavior Surface Record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^TestAnAdopterManagedPathIsMaintenanceOnly$' ./internal/releaseplan 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in the_setup_manifest a_recorded_guide skills_lock a_locked_skill_file a_claude_skill_link claude_md specific_repository a_generated_author_subagent project_config a_nested_context_md an_unlocked_local_skill_crosses another_claude_agent_crosses context_map_crosses an_unrecorded_guide_crosses a_mixed_commit_crosses not_an_adopter_changes_nothing an_empty_commit_still_crosses; do printf '%s\n' "$out" | grep -q -- "--- PASS: TestAnAdopterManagedPathIsMaintenanceOnly/$name " || { echo "missing pass: $name" >&2; exit 1; }; done; out="$(go test -count=1 -timeout 20m -v -run '^(TestTheReleasePlanReadsTheManagedRecordsAtBaseAndTarget|TestTheReleasePlanHelpNamesTheAdopterMaintenancePaths|TestReleasePlanHelpDescribesFlagsDefaultsStatesAndReadOnlyBoundary)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheReleasePlanReadsTheManagedRecordsAtBaseAndTarget/the_fiscus_shape TestTheReleasePlanReadsTheManagedRecordsAtBaseAndTarget/roundfix_itself TestTheReleasePlanReadsTheManagedRecordsAtBaseAndTarget/a_skill_removed_in_the_range TestTheReleasePlanReadsTheManagedRecordsAtBaseAndTarget/a_malformed_lock_makes_no_adopter TestTheReleasePlanReadsTheManagedRecordsAtBaseAndTarget/project_config_is_not_read TestTheReleasePlanHelpNamesTheAdopterMaintenancePaths/adopter_test TestTheReleasePlanHelpNamesTheAdopterMaintenancePaths/roundfix_itself TestReleasePlanHelpDescribesFlagsDefaultsStatesAndReadOnlyBoundary; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; out="$(go test -count=1 -timeout 20m -v -run '^TestTheReleaseDocumentsNameTheAdopterMaintenancePaths$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in canonical/adopter_test embedded/roundfix_itself runbook/skills_lock usage/project_config_and_context; do printf '%s\n' "$out" | grep -q -- "--- PASS: TestTheReleaseDocumentsNameTheAdopterMaintenancePaths/$name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0. Before this Task none of the four new tests exists, so their pass lines are missing; afterwards every named subtest and the existing help test pass.
- `grep -F -- '**Adopter Maintenance Path**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Adopter Maintenance Path\*\*(:| —( |$))' || { echo 'glossary lacks the entry' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "so it has no Adopter Maintenance Path" || { echo 'glossary entry lacks its last sentence' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/release-runbook.md | grep -qF -- "an Adopter Maintenance Path is maintenance-only too" || { echo 'runbook lacks the rule' >&2; exit 1; }; make skills-sync-check && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract && go test -count=1 -timeout 30m ./internal/releaseplan ./internal/cli` — expected: exit 0. Before this Task neither text exists; afterwards the glossary and the runbook hold them, the mirrors and the version record are current, and the skill, documentation, classifier and CLI suites pass.

## References

- `_prd.md` → Goals (second and third bullets); Core Features 3 and 4; Success Metric 2; Success Metric 3; Glossary; Decisions
- `_techspec.md` → Interfaces and invariants 3 to 5; API Contract 2; Coverage Map; Testing Approach; Build Order 2
- ADR-0447; ADR-0412
- [`references/2026-11-19-characterization.md`](references/2026-11-19-characterization.md), sections 2 and 3

## Result

Skills loaded: implement-task, agentic-cli-design, coding-guidelines, golang-cli, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, domain-modeling, tech-writer, evidence-gate, no-workarounds, systematic-debugging.

Implemented this Task's adopter maintenance boundary. `ClassifyChanges` and
`Build`'s release-preparation pre-check use `ManagedPaths.Contains`, while
`ClassifyCommit` keeps its signature and legacy classification. The optional
`ManagedPathSource` reads the Setup Manifest and skill lock through Git objects
at both endpoints, joins their recorded paths and skills, and identifies an
adopter from the target lock alone. Missing or malformed records contribute
nothing; Git failures retain `GitSourceError` and propagate through `Build`.
No project configuration content is read.

Acceptance evidence from focused checks:

- Maintenance boundary: the 17 required subtests cover every positive path,
  unlocked skills, other agents, `CONTEXT-MAP.md`, unrecorded guides, mixed
  commits, a non-adopter and an empty commit. Additional checks cover invalid
  paths and skill names, conventional and breaking subjects confined to
  maintenance, unchanged `ClassifyCommit` behavior, and the Build pre-check
  avoiding grant reads for managed maintenance.
- Git source: the five required CLI scenarios assert exact blocking commit
  SHAs and each commit's crossing flag. Additional tests join newly recorded
  and retired guides and skills, overwrite the working-tree records to prove
  they are ignored, reject absent/malformed/non-object lock records, and
  exercise typed failures of `rev-parse`, `ls-tree` and `show`.
- Public contract: the help and all four release documents hold API Contract
  2 verbatim apart from wrapping. Eight separate phrase checks cover each
  document and the help. The existing flags/defaults/states/read-only help
  test is unchanged. `CONTEXT.md` gains only the Adopter Maintenance Path
  entry immediately after Release Plan, with the exact definition, ADR-0447
  link and `_Avoid_` line.
- Generated records: `make skills-sync` updates the embedded reference and
  skill. The sanctioned version recorder raised both Roundfix skill fields
  from `0.1.276` to `0.1.278`; its append-only record also retains `0.1.277`,
  allocated before the document check caught and corrected a hyphen wrapping
  error. The Behavior Surface Record generator reports success and changes
  only the release-plan help and the two edited guide digests. It is run
  again after this Result as the final repository edit.

Focused commands and outcomes:

- Before the implementation,
  `go test ./internal/cli -run '^TestTheReleasePlanReadsTheManagedRecordsAtBaseAndTarget$/^the_fiscus_shape$' -count=1`
  exited 1: all four commits blocked, instead of only the fourth. Before the
  document addition, the canonical `adopter_test` phrase check exited 1 for
  the missing paragraph.
- `go test ./internal/releaseplan ./internal/cli -run 'Test(AnAdopter|Managed|InvalidManaged|BuildReadsManaged|TheReleasePlanReads|TheReleasePlanHelp|ReleasePlan|BuildLoadsRelease|BuildLooksFor|Classify)' -count=1`
  exited 0 in both packages, including the original Fiscus regression.
- `go test ./internal/releaseplan -run 'Test(MaintenanceOnly|ReleasePreparationCommitSetsNoImpact)$' -count=1`
  exited 0; the existing maintenance and release-preparation tests are unchanged.
- `go test ./skills -run 'Test(TheReleaseDocumentsNameTheAdopterMaintenancePaths|EveryOwnedSkillVersionIsRecorded|RoundfixSkillIndex|RoundfixSkillEntryFileStaysWithinItsBudget)' -count=1`
  exited 0 after correcting wrapping. The first run exposed the broken
  `skills-lock.json` word in all four document copies; no assertion was weakened.
- `make skills-sync`,
  `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
  and a subsequent `make skills-sync` exited 0.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
  exited 0; this same sanctioned command is the final edit after this Result.
- `git diff --check` exited 0. Scope inspection found only this Task's declared
  paths; the Daemon-owned initial status change is preserved. No existing
  releaseplan test file, other Task, Task Graph, status view, Baseline file,
  other owned skill, Vendored Skill or prohibited record/configuration changed.

The sandbox executor could not start because its configured writable root
contains a symlink. Commands ran with approved escalation from the physical
checkout. Declared Verification and repository Verification remain for the
Daemon; Task status is untouched. No commit, push, pull request, ai-memory
message or Spec event was issued. No follow-up work identified in this slice.

## Carry-forward provenance

- Source Run: `run_20261119T151201Z_05140c8b41c958d4`
- Source commit: `6a5547beb2051fe1b48fa85c53ba0309f23d6486`
