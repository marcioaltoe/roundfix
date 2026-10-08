---
task: task_03
spec: 0251-skills-keep-up-with-the-behavior-they-describe
status: pending
type: backend
complexity: high
---

# Task 03: The release plan refuses a release with a Lagging Surface

## Overview

Adds the `skill-coverage` check to a range Release Plan. It reads the Skill
Coverage Map and the Behavior Surface Record at the plan's base and target
from Git, lists every Lagging Surface, and, for a plan that proposes a
version, exits 3 with a next action that names the check. The decision state,
the proposed version and the skills and baseline checks do not change. The
release runbook, the usage guide, the Roundfix skill's release reference and
the glossary describe it. The slice is verifiable on its own through fixture
repositories and the command's help.

## Requirements

1. MUST implement `_techspec.md` → API Contract 4 in the release plan's check
   collection: read both files at the base and target commits through the
   existing release plan Git runner with `git cat-file -e`, `git show` and
   `git diff --name-only --no-renames`, parse them with `internal/skillcoverage`,
   and judge them with `skillcoverage.Compare`. It MUST NOT read the working
   tree for this check, so `--to` plans a past revision correctly.
2. MUST implement the text and JSON of `_techspec.md` → API Contract 5 and the
   exit code of API Contract 6 exactly, including the `Release blocked:` line
   and the replaced `Next action:` line when blocking. The `--reset-to` mode,
   the JSON schema version, the `skills` and `baseline` checks, the decision
   state and the proposed version MUST stay unchanged.
3. MUST update the `release plan` help as API Contract 6 states, keeping the
   sentence "A range plan also reports the skills and baseline checks
   read-only; they never change the decision state, the proposed version, or
   the exit code." byte-identical and adding the phrase "a blocking
   skill-coverage check".
4. MUST add `internal/cli/releaseplan_skill_coverage_test.go` with these
   tests, each on a temporary repository built with the existing release plan
   fixtures and map and record files committed at the base and the target:
   - `TestReleasePlanBlocksALaggingSurface`: Surface Transcript 1 exactly,
     its JSON (`blocking` true, one `lagging` item), exit 3 for a `ready`
     plan, and an `approval_required` plan that keeps its state, version and
     approval question while exiting 3 with the blocking next action.
   - `TestReleasePlanSkillCoverageIsCurrentWhenItsSkillChanges`: Surface
     Transcript 2 exactly; a changed Coverage Review gives `reviewed` and exit
     0; an unchanged review and a removed surface with an unchanged skill stay
     lagging; an uncovered surface never lags.
   - `TestReleasePlanSkillCoverageNeverBlocksWithoutAMap`: Surface Transcript
     3 exactly, `introduced` for a base without a map, and a `no_release`
     range with a Lagging Surface that reports `behind` with exit 0, each with
     the decision and exit code the same range has without the files.
   - `TestReleasePlanSkillCoverageFailsClosedOnAnUnreadableMap`: a malformed
     map, and a target map without a record, each report `failed`, block a
     `ready` plan with exit 3, and write nothing to the repository.
5. MUST keep every existing release plan test passing unchanged, including
   `TestReleasePlanReportsTheSkillsAndBaselineChecks` and
   `TestReleasePlanChecksNeverChangeTheDecision`.
6. MUST describe the check in `docs/user-guide/release-runbook.md` (its
   skills-and-guides step and its cutting step), in the release section of
   `docs/user-guide/usage.md`, and in the Roundfix skill's release reference,
   replacing every sentence that says the plan's checks never change its exit
   code. The skill edit MUST be synced to its mirror with `make skills-sync`
   and recorded with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   which chooses the next version.
7. MUST re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
   after the help and guide changes, so the repository contract stays green.
8. MUST write the glossary term **Lagging Surface** into `CONTEXT.md` and
   revise **Release Plan** to say that a Lagging Surface blocks a plan that
   proposes a version, through the `domain-modeling` skill, citing ADR-0256.

## Subtasks

- [ ] Read the map and record at both commits and judge them.
- [ ] Print the check, block through the exit code and the next action, and update the help.
- [ ] Write the four fixture tests.
- [ ] Update the runbook, the usage guide and the skill reference; sync, record and re-record.
- [ ] Add and revise the glossary terms.

## Acceptance Criteria

- [ ] A range that changes a surface without its skill exits 3 with the Lagging Surface named, and its state and version are unchanged.
- [ ] A covering skill edit or a changed Coverage Review clears it.
- [ ] No map, a base without a map, and a `no_release` range never block.
- [ ] An unreadable map blocks and the command writes nothing.
- [ ] The docs and the skill no longer say the checks never change the exit code.

## Context

- interface: `internal/cli/releaseplan_checks.go`
- interface: `internal/cli/releaseplan_command.go`
- interface: `internal/cli/releaseplan_git_source.go`
- interface: `internal/cli/cli.go`
- creates: `internal/cli/releaseplan_skill_coverage_test.go`
- interface: `docs/user-guide/release-runbook.md`
- interface: `docs/user-guide/usage.md`
- interface: `.agents/skills/roundfix/references/release.md`
- interface: `skills/roundfix/references/release.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- creates: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`
- instruction: `internal/cli/releaseplan_checks_test.go`
- instruction: `internal/cli/releaseplan_command_test.go`
- instruction: `internal/docscontract/release_step_test.go`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `docs/adr/0256-a-release-waits-for-the-skills-that-describe-a-changed-surface.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestReleasePlanBlocksALaggingSurface|TestReleasePlanSkillCoverageIsCurrentWhenItsSkillChanges|TestReleasePlanSkillCoverageNeverBlocksWithoutAMap|TestReleasePlanSkillCoverageFailsClosedOnAnUnreadableMap|TestReleasePlanReportsTheSkillsAndBaselineChecks|TestReleasePlanChecksNeverChangeTheDecision)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for t in TestReleasePlanBlocksALaggingSurface TestReleasePlanSkillCoverageIsCurrentWhenItsSkillChanges TestReleasePlanSkillCoverageNeverBlocksWithoutAMap TestReleasePlanSkillCoverageFailsClosedOnAnUnreadableMap TestReleasePlanReportsTheSkillsAndBaselineChecks TestReleasePlanChecksNeverChangeTheDecision; do printf '%s\n' "$out" | grep -q -- "--- PASS: $t " || { printf 'missing passing test %s\n' "$t" >&2; exit 1; }; done; help="$(go run -buildvcs=false ./cmd/roundfix release plan --help 2>&1)" || { printf '%s\n' "$help"; exit 1; }; printf '%s\n' "$help" | tr -s '[:space:]' ' ' | grep -qF -- 'a blocking skill-coverage check' || { printf 'release plan help does not name the blocking check\n%s\n' "$help" >&2; exit 1; }; docs="$(go test -count=1 -tags docscontract -run '^(TestTheSkillCoverageMapIsCurrent|TestReleasePlanDocumentationContract|TestTheReleaseRunbookRequiresTheSkillsAndGuidesCheck|TestEveryCheckTheReleaseStepNamesExists)$' ./internal/docscontract 2>&1)" || { printf '%s\n' "$docs"; exit 1; }; go test -count=1 ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' || exit 1; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/release.md | grep -qF -- 'Lagging Surface' || { printf 'the release reference does not describe the Lagging Surface\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Lagging Surface**:' || { printf 'missing glossary term Lagging Surface\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Release Plan**:' || { printf 'missing glossary term Release Plan\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -o -- '[*][*]Release Plan[*][*]:[^*]*' | grep -qF -- 'Lagging Surface' || { printf 'the Release Plan definition does not name the Lagging Surface\n' >&2; exit 1; }` — expected: exit 0. Before this Task the four new tests do not exist, so the first check fails. After it, the new and existing release plan tests pass, the help names the blocking check, the documentation contracts and the skill version record pass with the record re-recorded, and the reference and glossary describe the Lagging Surface.

## References

- `_prd.md` → Core Feature 4; Core Feature 6; User Story 1; User Story 3; User Story 5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 6; Glossary
- `_techspec.md` → API Contract 4; API Contract 5; API Contract 6; Invariant 3; Invariant 4; Invariant 5; Invariant 6; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Build Order 3
- ADR-0256; ADR-0189
