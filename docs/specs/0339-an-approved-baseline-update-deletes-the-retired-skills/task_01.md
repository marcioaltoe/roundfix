---
task: task_01
spec: 0339-an-approved-baseline-update-deletes-the-retired-skills
status: completed
type: backend
complexity: high
---

# Task 01: An approved Baseline update deletes the Retired Skills it proves unedited

## Overview

`baseline.InstalledRetiredSkills` gives each held Retired Skill a disposition,
`delete` when its copies match the `computedHash` of its `skills-lock.json`
entry and its `.claude/skills` link points at `.agents/skills`, otherwise
`keep` with a reason. A new `baseline.DeleteRetiredSkills` deletes exactly the
planned `delete` entries. `roundfix baseline update` binds the deletions to
the Plan Digest, reports `plan_ready` when a deletion is the only pending
change, deletes after the skills stage of an approved apply and reports
`Skills deleted`. This answers the adopted Backlog Entry of 2026-11-01, "An
approved Baseline update deletes the Retired Skills"
([references](references/2026-11-01-an-approved-baseline-update-deletes-the-retired-skills.md)),
the Argus report of 2026-11-01.

## Requirements

1. MUST characterize first: at the starting tree an applied update over
   `adoptedRepositoryWithRetiredSkills` leaves the copies and the lock
   unchanged (`TestBaselineUpdateAppliesAndListsRetiredSkills`).
2. MUST extend `InstalledRetiredSkill` in
   `internal/baseline/retired_skills_installed.go` with `Action`, `Reason`
   and `LockHash` and add the constants `RetiredSkillDelete`,
   `RetiredSkillKeep`, `RetiredSkillEdited`, `RetiredSkillUnrecorded` and
   `RetiredSkillUnsafe`, exactly as `_techspec.md` → Interfaces states, and
   make `InstalledRetiredSkills` decide each disposition as API Contract 1
   states, hashing a folder with `inspectRestoreTarget` and
   `externalSkillsLockDigest`. It still only reads.
3. MUST add `DeleteRetiredSkills(ctx, repoRoot, planned)` to the same file as
   API Contract 2 states: re-inspect, refuse a name outside the Retired
   Skills set and an entry whose disposition changed, both before deleting
   anything, then remove the paths `.claude` first with `os.RemoveAll` and
   write `skills-lock.json` once through `loadSkillsLock`,
   `removeSkillsEntry` and `marshalIndent`, keeping its mode.
4. MUST change `runBaselineUpdateCommandWithSkillsStage` in
   `internal/cli/baseline_update.go` as API Contracts 3 to 5 state: inspect
   before the confirmation, unless `--no-skills`; fold the `delete` entries
   into the Plan Digest with `retiredSkillsPlanDigest`, which returns its
   input unchanged when there is none; report `plan_ready` with the message
   "guidance matches the current Baseline catalog; the approved update
   deletes <n> Retired Skill(s)" when a deletion is the only pending change;
   after the skills stage of an approved apply call `DeleteRetiredSkills`
   with the entries the digest bound, fill `skills.deleted` (field
   `Deleted`, tag `deleted,omitempty`) and re-inspect into
   `skills.retired`.
5. MUST print the text of API Contract 6 in `writeBaselineUpdateResult`:
   `Skills deleted: <n>` and `- deleted <skill>: <targets>` before
   `Skills retired`, the `delete` line "no longer required by the Baseline;
   the approved update deletes <targets>", and the `keep` line "kept,
   <reason text>; delete <targets> by hand if the repository no longer uses
   it" with the three reason texts verbatim.
6. MUST add to `internal/baseline/retired_skills_installed_test.go`
   `TestRetiredSkillDispositionsProveACopyUnedited` and
   `TestDeleteRetiredSkillsDeletesOnlyThePlannedCopies` with the cases and
   subtests `_techspec.md` → Testing Approach lists, and give the cases of
   `TestInstalledRetiredSkillsListsTreesLinksAndLockEntries` their action and
   reason.
7. MUST add `TestContext7AndExaAreRequiredByEverySnapshotAndNeverRetired` to
   `internal/baseline/retired_skills_test.go`.
8. MUST add to `internal/cli/baseline_update_retired_skills_test.go` the
   fixture `deletableRetiredSkillsRepository`,
   `TestBaselineUpdateApprovedApplyDeletesTheRetiredSkills` (subtests `text`
   and `json`) and `TestBaselineUpdateRefusesARetiredSkillEditedAfterThePlan`
   as Testing Approach lists, each top-level test calling `t.Parallel()`
   first, and update `retiredSkillsTranscript`,
   `expectedInstalledRetiredSkills` and
   `TestBaselineUpdateRetiredSkillDeletionTargets` as Declared breaks states,
   keeping every top-level test name.
9. MUST update `wantRetired` in the combined preview subtest of
   `TestBaselineUpdatePreviewListsASetupSnapshotDifference`
   (`internal/cli/repository_profile_docscontract_test.go`) to the council
   keep line.
10. MUST rewrite the `### Retired skills` section of
    `docs/user-guide/commands/baseline.md` as `_techspec.md` → API Contract 7
    states: the sentence "An approved update deletes each Retired Skill it
    proves unedited" appears verbatim; the section names the
    `.agents/skills/<skill>` folder, the `.claude/skills/<skill>` entry and
    the lock entry, the `plan_ready` preview, the Plan Digest,
    `Skills deleted` (`skills.deleted` in JSON), the three keep reasons and
    that nothing outside the Retired Skills set is deleted, and no longer
    says the update never deletes an installed copy. It keeps the list of
    Retired Skills and its manual commands, introduced as the way to remove
    a kept copy. Change no other section.
11. MUST re-record the Behavior Surface Record as the LAST step, after the
    last edit of the guide, with `go test -count=1 -tags docscontract
    ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$'
    -record-skill-coverage`, and commit `docs/references/behavior-surfaces.json`.
12. MUST NOT delete or change any skill outside `RetiredSkills()`, change the
    Retired Skills set, the asset sync, the Setup Snapshots, Doctor,
    `baseline skills reconcile`, `skills-lock.json` of this repository, any
    Baseline module, `docs/agents/`, the Makefile, CI workflows, `go.mod` or
    `go.sum`. Every new test file compiles for Windows
    (`GOOS=windows go vet ./internal/baseline ./internal/cli`).

## Subtasks

- [ ] Characterize the current apply over the retired fixture.
- [ ] Add the dispositions and the deletion with their baseline tests.
- [ ] Bind the deletions to the Plan Digest, add the preview state, the deletion step and the text, with the CLI tests.
- [ ] Update the tests the Declared breaks name.
- [ ] Rewrite the guide section and re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] An approved update deletes the unedited `tdd` folder, link and lock entry and reports it in text and JSON; the preview lists it and the Plan Digest binds it.
- [ ] `context7-cli`, `exa-web-search`, a repository-installed skill and an edited retired copy survive byte for byte, the edited copy listed as kept.
- [ ] A confirmation after the planned copy changed deletes nothing.
- [ ] The baseline guide describes the deletion, and the Behavior Surface Record passes.

## Context

- instruction: `docs/adr/0350-an-approved-baseline-update-deletes-the-retired-skills-it-proves-unedited.md`
- interface: `internal/baseline/retired_skills_installed.go`
- interface: `internal/baseline/retired_skills_installed_test.go`
- interface: `internal/baseline/retired_skills_test.go`
- interface: `internal/cli/baseline_update.go`
- interface: `internal/cli/baseline_update_retired_skills_test.go`
- interface: `internal/cli/repository_profile_docscontract_test.go`
- interface: `docs/user-guide/commands/baseline.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- `baseline.InstalledRetiredSkill`, `InstalledRetiredSkills` and `DeleteRetiredSkills` — consumers: `./internal/baseline`, `./internal/cli`
- The `skills.retired` and `skills.deleted` JSON and text of `roundfix baseline update` — consumers: `./internal/cli`
- The baseline guide and the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestRetiredSkillDispositionsProveACopyUnedited|TestDeleteRetiredSkillsDeletesOnlyThePlannedCopies|TestContext7AndExaAreRequiredByEverySnapshotAndNeverRetired)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for t in TestRetiredSkillDispositionsProveACopyUnedited TestDeleteRetiredSkillsDeletesOnlyThePlannedCopies TestContext7AndExaAreRequiredByEverySnapshotAndNeverRetired; do printf '%s\n' "$out" | grep -q -- "--- PASS: $t " || { echo "missing pass: $t" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./internal/baseline` — expected: exit 0. Before this Task the three tests do not exist, so a pass line is missing; afterwards they pass with the whole baseline suite.
- `out="$(go test -count=1 -timeout 20m -v -run '^(TestBaselineUpdateApprovedApplyDeletesTheRetiredSkills|TestBaselineUpdateRefusesARetiredSkillEditedAfterThePlan)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for t in TestBaselineUpdateApprovedApplyDeletesTheRetiredSkills TestBaselineUpdateRefusesARetiredSkillEditedAfterThePlan; do printf '%s\n' "$out" | grep -q -- "--- PASS: $t " || { echo "missing pass: $t" >&2; exit 1; }; done; go test -count=1 -timeout 20m -run 'BaselineUpdate|Retired' ./internal/cli && go test -count=1 -timeout 20m -tags docscontract -run '^TestBaselineUpdatePreviewListsASetupSnapshotDifference$' ./internal/cli && GOOS=windows go vet ./internal/baseline ./internal/cli && tr -s '[:space:]' ' ' < docs/user-guide/commands/baseline.md | grep -qF -- 'An approved update deletes each Retired Skill it proves unedited' && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the two tests do not exist; afterwards they pass with every Baseline update test of the CLI, the repository-profile contract test and the Windows build of both packages, the guide states the rule, and the documentation contracts, the Behavior Surface fingerprint among them, pass.

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5
- `_techspec.md` → Current behavior 1-6; API Contracts 1-7; Invariants 1-5; Declared breaks; Testing Approach; Build Order 1

## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-context, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, systematic-debugging, no-workarounds, tech-writer, evidence-gate.

Implemented the assigned retirement slice. Inspection now proves each disposition
with the existing safe-tree inspector and lock digest. Deletion validates every
planned delete entry before removing any path, removes `.claude` before `.agents`,
and rewrites only the deleted lock entries through the ordered lock helpers,
preserving the lock mode. Unsafe parent paths are kept as well as unsafe entries.
The update inspects before confirmation, binds only deletions to the Plan Digest,
reports deletion-only previews as `plan_ready`, deletes after the skills stage,
and reports exact deleted and kept dispositions in text and JSON.

Acceptance evidence:

- Approved deletion and preview: `TestBaselineUpdateApprovedApplyDeletesTheRetiredSkills`
  passes its `text`, `json`, and `yes` subtests. They exercise the adopted
  temporary repository through the existing command seam, check the preview's
  state/category/message and its digest against `--no-skills`, check deletion
  reporting, and confirm the `tdd` folder, link, and lock entry leave.
  `TestBaselineUpdateRetiredSkillDeletionTargets` also checks the exact JSON
  digest encoding and an unchanged digest for keep-only entries.
- Preservation: the same text/JSON apply subtests compare survivor trees with
  `baselinePlanTestTree`, check every survivor link, and compare the whole lock
  after removing only `tdd` from its expected data. `design-patterns` is kept
  as `edited`; `context7-cli`, `exa-web-search`, and `local-only` survive.
  `TestContext7AndExaAreRequiredByEverySnapshotAndNeverRetired` checks every
  embedded Setup Snapshot. The disposition test compares the complete six-skill
  slice including hashes and additionally covers unsafe trees, links, parents,
  link-only copies, and an edited `.claude` tree.
- Stale approval: `TestBaselineUpdateRefusesARetiredSkillEditedAfterThePlan`
  edits the planned copy without making another commit, confirms the old digest,
  and observes exit 3, `action_required`/`approval`, no deletion report, and no
  further tree change. The baseline deletion test checks refusal before any
  removal for an edited copy, an out-of-set name, and a later invalid entry.
  The apply regression also changes `tdd` during the skills stage and observes
  category `skills` with the required next action, unchanged lock bytes, and a
  surviving link. Cancellation and selective/orphaned deletion have focused cases.
- Guide and record: only `### Retired skills` was rewritten, preserving the
  retired-name list and manual commands for kept copies. It describes the proof,
  approval, Plan Digest, `plan_ready`, deletion reporting, all three keep reasons,
  and the Retired Skills boundary. The sanctioned recorder ran after the final
  guide edit and exited 0; its only generated delta is this guide's fingerprint
  in `docs/references/behavior-surfaces.json`.

Focused checks (not Daemon Verification):

- Before implementation: `go test -count=1 -timeout 5m -v ./internal/cli -run '^TestBaselineUpdateAppliesAndListsRetiredSkills$'`
  exited 0, characterizing that the original fixture's copies and lock stay.
  The new `TestBaselineUpdateApprovedApplyDeletesTheRetiredSkills/text` then
  failed on the starting implementation: preview exit 0 and `current` instead
  of the deletion-only approval state.
- After the final test edits: `go test -count=1 -timeout 5m ./internal/baseline ./internal/cli -run '^Test(InstalledRetiredSkills|RetiredSkill|DeleteRetiredSkills|Context7AndExa|NoCatalogEntryNamesARetiredSkill|BaselineUpdate.*Retired|BaselineUpdateWithoutSkills)'`
  exited 0 for both packages.
- `go test -count=1 -timeout 5m -tags docscontract ./internal/cli -run '^TestBaselineUpdatePreviewListsASetupSnapshotDifference$/difference_with_outdated_and_retired_skills$'`
  exited 0 for the changed combined-preview assertion.
- Final Windows test cross-compilation: `GOOS=windows go test -c -o <temporary-output> ./internal/baseline`
  and the equivalent `./internal/cli` command each exited 0. Binaries were
  written under `/private/tmp`; the declared Windows vet command was not run.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
  exited 0 as the final generated-artifact step after the guide edit.
- Five temporary `go test -overlay` probes each exited 1 at the intended
  behavioral assertion: skipping deletion leaves `tdd`; ignoring hash differences
  marks `design-patterns` deletable; omitting the digest fold allows stale
  confirmation and produces the `--no-skills` digest; omitting the name guard
  loses the required refusal; accepting arbitrary link targets loses `unsafe`.
  The transcripts were inspected, not classified from exit codes alone.
  Probe overlays and full logs are under
  `/private/tmp/roundfix-retired-probes-zmy59uxk/`; repository source was not
  changed for these probes.

The original `status: in_progress` is preserved. No other Task, Task Graph,
skill, Setup Snapshot, module, repository lock, or tooling configuration was
edited. Declared Verification and terminal settlement remain Daemon-owned;
no commit, push, or pull request was made. No follow-up outside this Task's
slice was implemented.
