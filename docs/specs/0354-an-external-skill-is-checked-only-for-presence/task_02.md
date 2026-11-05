---
task: task_02
spec: 0354-an-external-skill-is-checked-only-for-presence
status: pending
type: backend
complexity: medium
---

# Task 02: A required external skill is checked only for presence

## Overview

The Repository Skill Set check hashes each required external skill and
reports a copy whose hash differs from its `skills-lock.json` entry as
`outdated`: Doctor fails on it, `skills check` rejects it, and `baseline
update --yes` restores it from the Setup Snapshot, the downgrade Fluxus asked
Roundfix never to make. This Task checks an external skill only for presence,
a lock entry and a folder, as ADR-0364 records: the hash comparison, the
`OutdatedExternal` field and every reader of it go, and this repository keeps
its own lock rule in its live Repository Skill Set test. It answers the
fluxus report of 2026-11-05, adopted as
[references/2026-11-05-only-roundfix-owned-skills-are-checked-by-digest-or-version.md](references/2026-11-05-only-roundfix-owned-skills-are-checked-by-digest-or-version.md).

## Requirements

1. MUST change `checkRepositoryWithReadiness` in `skills/repository.go` as API
   Contract 1 states: replace `requiredExternalHashes` with
   `requiredExternalLockEntries`, with the doc comment of `_techspec.md` →
   Interfaces, which reports an absent entry exactly as today and reads no
   `computedHash`; remove the folder hash and the `OutdatedExternal`
   comparison from the external loop, keeping its absent-folder, symbolic
   link and non-directory handling; remove `OutdatedExternal` from
   `RepositoryReadiness` and `Ready`, and the hash members and helpers that
   become unused. Every lock refusal and every owned-skill result keeps its
   bytes.
2. MUST remove every reader of `OutdatedExternal` in `internal/cli`:
   - `doctorSkillReadinessResult` and `doctorExternalSkillNextActions`
     (`internal/cli/doctor.go`) print no `outdated:` detail and render one
     `bunx skills add marcioaltoe/skills@<name>` per missing external skill;
   - `missingRequiredSkills` (`internal/cli/required_skills.go`) no longer
     adds the reason "required by the Baseline and no longer matches its
     skills-lock.json entry";
   - `baselineUpdateExternalDrift` (`internal/cli/baseline_update.go`)
     returns only missing skills, and the preview message becomes
     `guidance matches the current Baseline catalog; %d required skill(s)
     are missing: %s`, as API Contract 3 states.
3. MUST change these tests to the new contract, with no other assertion
   changed:
   - `skills/repository_test.go`: the `changed external file` and `added
     external file` rows of `TestCheckRepositoryClassifiesMissingAndOutdatedSkills`
     expect no classification and the `wantOutdatedExt` member goes;
     the `invalid required hash` row of
     `TestCheckRepositoryRejectsMalformedLockAndUnsafeRequiredNames` goes;
     `TestCheckRepositoryHandlesNestedLinksSpecialEntriesAndStableOrdering`
     expects the external nested symbolic link and the special entry to be
     ready instead of refused; `TestCheckRepositoryWithExternalAcceptsEmptyRequirement`
     drops its `OutdatedExternal` check;
   - `internal/cli/doctor_test.go`: the mixed row of the Doctor readiness
     table drops its outdated skill, its `outdated: testing-boss` detail and
     its `testing-boss` install action;
   - `internal/cli/baseline_update_test.go`: the fake readiness that reported
     `OutdatedExternal` reports the same names as `MissingExternal`;
   - `internal/cli/baseline_update_missing_skills_test.go`: the `folder no
     longer matches its lock entry` row of
     `TestBaselineUpdatePreviewNamesAMissingRequiredSkill` expects `current`
     with exit `0` and no `"missing"`, the other rows expect the message of
     requirement 2, and `TestMissingRequiredSkillsCollectsNamedEvidence`
     builds its evidence from `MissingExternal` only.
4. MUST make `TestThisRepositoryHoldsEveryRequiredExternalSkill`
   (`internal/cli/live_documents_docscontract_test.go`) hold each required
   external skill of this repository to its lock as API Contract 6 states,
   through `skills.SkillFolderHash` and the entry's `computedHash`, failing
   with `external skill <name> no longer matches its skills-lock.json
   computedHash; run make skills-third-party-bump`, only when the test wants
   a ready set; drop its `OutdatedExternal` logging.
5. MUST add `skills/external_presence_test.go` with
   `TestCheckRepositoryChecksAnExternalSkillOnlyForPresence`, calling
   `t.Parallel()` first, over `agentic-cli-design` in
   `writeReadyRepositoryFixture`: rows `changed file`, `added file`, `nested
   symbolic link` and `lock entry without computedHash` are ready with no
   error; `missing directory` reports `MissingExternal` `[agentic-cli-design]`;
   `no lock entry` returns a `RepositoryReadinessError` whose
   `MissingExternal` is `[agentic-cli-design]`; `skill root is a symbolic
   link` returns a `RepositoryReadinessError` with ownership `external` and
   the skill root as its path.
6. MUST add `internal/cli/external_skill_lock_test.go` with two tests, each
   calling `t.Parallel()` first, over `writeDoctorReadyRepositoryFixture`
   with `testing-boss`'s `SKILL.md` rewritten after its lock entry was
   recorded and the helper `setDoctorTestManifestProfile` of task_01:
   - `TestDoctorKeepsAnExternalSkillThatDiffersFromItsLockReady`: `doctor`
     through `runCLI` with `skills.CheckRepositoryWithExternal` and the one
     external skill `testing-boss` prints the line `skills: ok (<n>
     required: <owned> Roundfix-owned, 1 external)`, with `<owned>` read from
     `len(skills.Names())`, no `outdated:`, and exits `0`;
   - `TestBaselineUpdateNeverRestoresAnExternalSkillThatDiffersFromItsLock`:
     `runBaselineUpdateSkillsStageWith` with the real
     `skills.CheckRepositoryWithExternal` restores nothing, reports no
     drifted skill and keeps the copy's bytes.
7. MUST change the `external failure` row of
   `TestDoctorOwnedLockWarningKeepsReadinessDetailAndNextActions`
   (`internal/cli/doctor_skills_lock_test.go`) so its external failure comes
   from a second required external skill, `agentic-cli-design`, that has no
   folder and no lock entry, instead of an outdated `testing-boss`, keeping
   its `failed` status and its assertions.

## Subtasks

- [ ] Check an external skill only for its lock entry and its folder.
- [ ] Remove every reader of `OutdatedExternal` and change the preview message.
- [ ] Move this repository's lock rule into its live test.
- [ ] Add the presence test and the two lock tests, and change the tests of requirements 3 and 7.

## Acceptance Criteria

- [ ] A required external skill whose content differs from its lock is ready for the check, Doctor and the preview, and is never restored.
- [ ] A skill without its folder or its lock entry is still reported, and every lock refusal is unchanged.
- [ ] This repository still fails its live test when a required skill no longer matches its lock.
- [ ] The `skills` and CLI suites, with and without `docscontract`, pass.

## Context

- instruction: `docs/adr/0364-a-required-external-skill-is-checked-only-for-presence.md`
- interface: `skills/repository.go`
- interface: `skills/repository_test.go`
- creates: `skills/external_presence_test.go`
- interface: `internal/cli/doctor.go`
- interface: `internal/cli/required_skills.go`
- interface: `internal/cli/baseline_update.go`
- interface: `internal/cli/doctor_test.go`
- interface: `internal/cli/baseline_update_test.go`
- interface: `internal/cli/baseline_update_missing_skills_test.go`
- interface: `internal/cli/live_documents_docscontract_test.go`
- interface: `internal/cli/doctor_skills_lock_test.go`
- creates: `internal/cli/external_skill_lock_test.go`

## Consumers

- `skills.CheckRepositoryWithExternal` and `RepositoryReadiness` without `OutdatedExternal` — consumers: `./skills`, `./internal/cli`, `./internal/baseline`
- Doctor's `skills:` line and the `baseline update` missing list — consumers: `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^TestCheckRepositoryChecksAnExternalSkillOnlyForPresence$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for row in changed_file added_file nested_symbolic_link lock_entry_without_computedHash missing_directory no_lock_entry skill_root_is_a_symbolic_link; do printf '%s\n' "$out" | grep -q -- "--- PASS: TestCheckRepositoryChecksAnExternalSkillOnlyForPresence/$row " || { echo "missing row: $row" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./skills ./internal/baseline` — expected: exit 0. Before this Task the test does not exist, so the run finds no row; afterwards every row passes, a changed or extended copy is ready, an absent folder or lock entry is still reported, and the `skills` and Baseline suites pass.

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestDoctorKeepsAnExternalSkillThatDiffersFromItsLockReady|TestBaselineUpdateNeverRestoresAnExternalSkillThatDiffersFromItsLock)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDoctorKeepsAnExternalSkillThatDiffersFromItsLockReady TestBaselineUpdateNeverRestoresAnExternalSkillThatDiffersFromItsLock; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m ./internal/cli && go test -count=1 -timeout 30m -tags docscontract ./internal/cli ./skills` — expected: exit 0. Before this Task neither test exists; afterwards a copy that differs from its lock passes Doctor and is not restored, and the CLI suite with and without `docscontract`, this repository's lock rule among them, passes.

## References

- `_prd.md` → Goals; Core Features 1, 3 and 4; Success Metric 2; Success Metric 4; Success Metric 6
- `_techspec.md` → Current behavior 1, 2, 4, 5 and 6; API Contracts 1, 2, 3, 5 and 6; Invariants 1 to 4; Declared breaks; Testing Approach; Build Order 2
