---
spec: 0354-an-external-skill-is-checked-only-for-presence
prd: _prd.md
created: 2026-11-05
---

# An external skill is checked only for presence — Technical Spec

## Executive Summary

`skills.CheckRepositoryWithExternal` (`skills/repository.go`) hashes every
required external skill folder and compares it with its `skills-lock.json`
`computedHash`, reporting a mismatch as `OutdatedExternal`. Doctor
(`repositorySkillsCheck`, `internal/cli/doctor.go`) and the skills stage of
`baseline update` (`internal/cli/baseline_update.go`) then compare each
lock-matching external skill with its Setup Snapshot
(`baseline.SetupSnapshotDifferences`) and print
`DR-SKILL-DIFFERS-FROM-SNAPSHOT` or list `skills.differing`. A lock mismatch
fails Doctor as `outdated` and makes `baseline update --yes` restore the copy
from the snapshot.

The fix checks an external skill only for presence: a lock entry that names
it and a directory at `.agents/skills/<name>`. The hash comparison, the
`OutdatedExternal` field, the snapshot comparison, the `skills.differing`
field and its text lines are removed, and the managed refresh restores only a
missing skill. The owned-skill rule already compares with the binary's copy
offline and does not change. This repository keeps its own
stricter lock and snapshot rule in its contract tests (ADR-0224, ADR-0360).
ADR-0364 records the decision.

The trade-off accepted: a hand-edited or stale external skill is no longer
reported by any command. The maintainer chose this: the copy in the
repository is the one its own tool installed, and only the deliberate
per-skill restore replaces it.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential is added. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0364 (new) records the rule:
  "A required external skill is checked only for presence". It amends four
  accepted records, whose notices task_03 writes: ADR-0221, whose Decision
  compares each lock-matching external skill with its Setup Snapshot and
  keeps the lock as "the authority for whether a skill is installed and
  intact"; ADR-0264, whose Decision names a Setup Snapshot Difference and
  lists it under `skills.differing`; ADR-0224, which ties this repository's
  tests to "the same comparison the Doctor Command uses"; and ADR-0241,
  whose first paragraph names "the snapshot comparison of ADR-0221" among
  the readers of a repository profile's skill contracts. ADR-0360 (the
  Catalog Skill Bump) and ADR-0322 (the Python profile) cite two of them and
  keep their rules. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle AZ batch 2 answers of 2026-11-05 quoted in the authorization record;
  the Governed Paths are `docs/references/coverage-record.json`, which
  task_01 re-records, and the Roundfix skill files task_03 changes; no
  Makefile, workflow, module or lint configuration changes. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0354-an-external-skill-is-checked-only-for-presence/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/baseline.md`, `.agents/skills/roundfix/references/runtime.md`, `.agents/skills/roundfix/references/setup.md`, `docs/references/coverage-record.json`, `skills/roundfix/SKILL.md`.

## Current behavior

Read on 2026-11-05 at `ad77a797`; `_prd.md` → Validation records the runs.

1. `checkRepositoryWithReadiness` (`skills/repository.go`) reads
   `skills-lock.json` when any external skill is required
   (`readLocalSkillsLock`), requires an entry with a 64-hex `computedHash`
   for each one (`requiredExternalHashes`, which reports missing entries
   through `RepositoryReadinessError.MissingExternal`), records an absent
   folder in `MissingExternal`, refuses a symbolically linked or non-directory
   skill root, then hashes the folder (`skillFolderHash`, which also refuses
   a nested symbolic link) and records a mismatch in `OutdatedExternal`.
2. `doctorSkillReadinessResult` (`internal/cli/doctor.go`) prints
   `outdated: <names>` and a `bunx skills add marcioaltoe/skills@<name>`
   action for each outdated skill (`doctorExternalSkillNextActions`).
3. `repositorySkillsCheck` calls `dependencies.snapshotDifferences`
   (`manifestSnapshotDifferences` → `baseline.SetupSnapshotDifferences`,
   `internal/baseline/skills_trailing.go`) over the names
   `snapshotComparableSkills` keeps, and appends
   `; DR-SKILL-DIFFERS-FROM-SNAPSHOT: differs from the Setup Snapshot, newer
   or older unknown: <names>` or `; snapshot comparison unavailable: <error>`.
4. The `baseline update` preview (`runBaselineUpdate`) fills
   `result.Skills.Missing` from `missingRequiredSkills`
   (`internal/cli/required_skills.go`), which adds outdated skills with the
   reason "required by the Baseline and no longer matches its skills-lock.json
   entry", prints `... required skill(s) are missing or no longer match
   skills-lock.json: <names>`, and fills `result.Skills.Differing` through
   `baselineUpdateSnapshotDifferences`. The skills stage
   (`runBaselineUpdateSkillsStageWith`) does the same comparison and restores
   every name `baselineUpdateExternalDrift` returns: missing and outdated
   ones. `writeBaselineUpdateText` prints `Skills differing: <n>` and one
   `- differs <name>: <reason>` line each; the JSON always carries
   `"differing"`.
5. The owned rule: `inspectOwnedSkillReadiness` compares the declared version
   with `ownedSkillMinimumVersions`, read from the embedded `SKILL.md` of each
   owned skill (`embeddedOwnedSkillVersions`); `Readiness` passes a version
   equal or higher. No network read happens anywhere in the check.
6. This repository's own rules: `TestThisRepositoryHoldsEveryRequiredExternalSkill`
   (`internal/cli/live_documents_docscontract_test.go`) holds its required
   external skills to their lock through `OutdatedExternal`;
   `TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot`
   (`internal/baseline/repository_skills_snapshot_test.go`) and
   `TestEveryVendoredCatalogSkillIsTheTreeItsCatalogCommitPublishes`
   (`internal/baseline/vendored_catalog_skills_docscontract_test.go`) both
   hold them to their Setup Snapshot tree digest.

## System Architecture

No package, command or flag is added. Changes by component:

- `internal/cli/doctor.go`, `internal/cli/baseline_update.go`: no snapshot
  comparison and no `differing` field (task_01); no `outdated:` detail,
  reason or restore for an external skill, with
  `internal/cli/required_skills.go` (task_02).
- `internal/baseline/skills_trailing.go`: deleted with its tests (task_01).
- `skills/repository.go`: the external loop checks presence only;
  `OutdatedExternal` and the `computedHash` reading are removed (task_02).
- Tests: `internal/cli/external_skill_presence_test.go` (task_01),
  `skills/external_presence_test.go` and
  `internal/cli/external_skill_lock_test.go` (task_02), all new, and the
  changes of Declared breaks.
- Documentation (task_03): the Doctor and Baseline guides, the README, the
  Roundfix `baseline`, `setup` and `runtime` references, two `CONTEXT.md`
  entries, and the amendment notices of ADR-0221, ADR-0224, ADR-0241 and
  ADR-0264.

## Implementation Design

### Interfaces

```go
// skills/repository.go

type RepositoryReadiness struct {
	OwnedRequired    int
	ExternalRequired int
	MissingOwned     []string
	Owned            []OwnedSkillReadiness
	MissingExternal  []string // no folder; OutdatedExternal is removed
}

// requiredExternalLockEntries proves that the lock names every required
// external skill. It reads no computedHash (ADR-0364).
func requiredExternalLockEntries(lockPath string, lock localSkillsLock, external []string) error
```

### Data Models

No schema change to any file Roundfix reads. The `baseline update` result
loses its `skills.differing` member (API Contract 3).

### API Contracts

1. API Contract: presence. For each required external skill,
   `CheckRepositoryWithExternal` requires a `skills-lock.json` entry naming
   it, with any members, and an `.agents/skills/<name>` entry that is a
   directory and not a symbolic link. A missing entry is reported through
   `RepositoryReadinessError.MissingExternal`, as today; a missing folder
   through `RepositoryReadiness.MissingExternal`, as today. It opens no file
   inside the folder and reads no `computedHash`. Every lock refusal of
   `readLocalSkillsLock` (absent, symbolic link, not a regular file,
   undecodable, trailing value, version other than `1`) and every
   Roundfix-owned check keep their behavior and messages.
2. API Contract: Doctor's `skills:` line. It never contains `outdated:`,
   `DR-SKILL-DIFFERS-FROM-SNAPSHOT` or `snapshot comparison unavailable`. For
   a present external skill it is `skills: ok (<n> required: <owned>
   Roundfix-owned, <external> external)` when nothing else fails; for one
   without its folder it is `skills: failed (missing: <name>; next: bunx
   skills add marcioaltoe/skills@<name>)` when nothing else fails. The
   `DR-SKILL-LOCK-STALE` suffix, the owned `below minimum` and `unversioned`
   details and their next actions keep their bytes.
3. API Contract: the managed refresh.
   - The JSON result has no `skills.differing` member and the text result no
     `Skills differing:` or `- differs ` line, with or without `--no-skills`.
   - The preview message for missing external skills is `guidance matches
     the current Baseline catalog; <n> required skill(s) are missing:
     <names>`; its next action keeps its bytes. Each `skills.missing` reason
     is `required by the Baseline and not installed; restore with roundfix
     baseline skills restore --profile <id> --skill <name>`.
   - The skills stage restores exactly the names in
     `RepositoryReadiness.MissingExternal` and `RepositoryReadinessError.MissingExternal`,
     and never calls the restore for a present skill.
4. API Contract: the deliberate restore. `roundfix baseline skills restore`
   is unchanged: it previews, exits `3` with a Plan Digest, and with
   `--confirm-plan <digest>` replaces a present external copy with its Setup
   Snapshot copy.
5. API Contract: the owned rule is unchanged: an owned copy passes when its
   declared version is equal to or higher than the embedded one, and fails
   with `below minimum: skill "<name>" requires <bundle>, found <installed>`
   when lower, offline.
6. API Contract: this repository's rule.
   `TestThisRepositoryHoldsEveryRequiredExternalSkill` compares each required
   external skill's `skills.SkillFolderHash` with its lock entry's
   `computedHash` and fails naming the skill and `make
   skills-third-party-bump`; `TestEveryVendoredCatalogSkillIsTheTreeItsCatalogCommitPublishes`
   keeps the snapshot rule.
7. API Contract: the rule in writing.
   - The rule sentence, the same Markdown in the Doctor and Baseline guides
     and in the Roundfix `baseline` and `setup` references:

     ```markdown
     Roundfix checks a required external skill only for presence: an entry in `skills-lock.json` and a folder under `.agents/skills`. It reads neither the folder's files nor the entry's `computedHash`, never compares the skill with its Setup Snapshot, and never replaces a present copy unasked; `roundfix baseline skills restore --profile <id> --skill <name>` replaces one deliberately.
     ```

   - `docs/user-guide/commands/doctor.md`: in the `skills:` bullet, "Each
     required external skill must hash to its `computedHash` in
     `skills-lock.json`." becomes "Each required external skill must be
     present: named by `skills-lock.json` and installed under
     `.agents/skills`."; everything from the paragraph starting "Doctor also
     compares each installed required external skill" through the paragraph
     starting "The comparison accepts repository-owned Baseline Profiles",
     with the example block between them, becomes the rule sentence; "A missing or outdated required skill, or an
     invalid required lock declaration," becomes "A missing required skill,
     or an invalid required lock declaration,"; the example line becomes
     `skills: failed (missing: agentic-cli-design; next: bunx skills add marcioaltoe/skills@agentic-cli-design)`.
   - `docs/user-guide/commands/baseline.md` and the canonical
     `.agents/skills/roundfix/references/baseline.md`: the paragraph from
     "Doctor reports `DR-SKILL-DIFFERS-FROM-SNAPSHOT`" to "`skills.drifted`."
     becomes the rule sentence followed by "The managed refresh restores only
     a missing skill, under `skills.drifted`."; "absent from `.agents/skills`
     or `skills-lock.json`, or no longer matches its lock entry: a Missing
     Required Skill." becomes "absent from `.agents/skills` or
     `skills-lock.json`: a Missing Required Skill."
   - `.agents/skills/roundfix/references/setup.md`: the paragraph "When a
     required external skill matches its lock ... nothing is restored."
     becomes the rule sentence.
   - `.agents/skills/roundfix/references/runtime.md`: "and checks each
     remaining skill only for presence: its `skills-lock.json` entry and its
     folder."
     becomes "and checks each remaining skill only for presence: its
     `skills-lock.json` entry and its folder."; "Each named missing or
     outdated external skill prints" becomes "Each named missing external
     skill prints".
   - `README.md`: "Each required external skill must match its
     `computedHash` in the repository's `skills-lock.json`. Missing,
     outdated, or invalid required skill state" becomes "Each required
     external skill must be present, named by the repository's
     `skills-lock.json` and installed under `.agents/skills`; its content is
     not compared. Missing or invalid required skill state".
   - `CONTEXT.md`, the entry **Missing Required Skill**: its first sentence
     becomes the first block below, and the rest of the entry keeps its
     bytes; the entry **Setup Snapshot Difference**: its body becomes the
     second block, and its `_Avoid_` line keeps its bytes. In `CONTEXT.md`
     each `ADR-NNNN` below is the Markdown link
     `docs/adr/<file>.md` of that record, as every glossary entry writes it.

     ```markdown
     An external skill the adopted modules require and dispatch that is absent from `.agents/skills` or `skills-lock.json`; a present copy is never one, whatever its content (ADR-0364).
     ```

     ```markdown
     A required external skill whose installed tree differs from the `treeDigest` its Setup Snapshot pins. Nothing offline tells whether the copy is newer or older, so Roundfix no longer looks for one: Doctor and the managed refresh check an external skill only for presence, and only a deliberate per-skill restore replaces it (ADR-0264, ADR-0364).
     ```

   - The amendment notices, each after the Decision paragraph named, with
     `ADR-0364` in the bold opening a Markdown link to
     `0364-a-required-external-skill-is-checked-only-for-presence.md`, and
     `updated_at: 2026-11-05T12:00:00Z`:
     - ADR-0221, after the paragraph starting "Doctor now compares each
       required external skill":

       ```markdown
       **Amended by ADR-0364 (2026-11-05).** Doctor and the managed refresh no longer compare a required external skill with its lock's `computedHash` or with its Setup Snapshot: an external skill is checked only for presence, a lock entry and a folder. This repository still holds its own Vendored Skills to the lock and the snapshot in its contract tests (ADR-0224).
       ```

     - ADR-0264, after the paragraph starting "**A difference is not
       trailing.**":

       ```markdown
       **Amended by ADR-0364 (2026-11-05).** Doctor no longer names a Setup Snapshot Difference and the managed refresh no longer lists `skills.differing`: neither compares an external skill with its Setup Snapshot, and a present copy is never restored unasked. `roundfix baseline skills restore` still replaces one deliberately, and the divergence order below still holds.
       ```

     - ADR-0241, after its first paragraph, which starts "A built-in
       Baseline Profile names one Setup Snapshot":

       ```markdown
       **Amended by ADR-0364 (2026-11-05).** The snapshot comparison of ADR-0221 no longer runs, so Doctor never reports it unavailable for a repository profile; skill restore and lock reconciliation still read a repository profile's contracts as this record decides.
       ```

     - ADR-0224, after the paragraph starting "The tests now read the two
       records":

       ```markdown
       **Amended by ADR-0364 (2026-11-05).** The Doctor Command no longer runs this comparison for any repository. This repository's contract tests keep both rules for its own skills: the vendored catalog test holds each tree to its Setup Snapshot, and the live Repository Skill Set test holds each required skill to its lock entry.
       ```

   Line wrapping may differ; the words may not.

### Surface Transcripts

None. Doctor prints a dozen readiness lines whose content depends on the
host's tools, so the changed `skills:` line and the `baseline update` lines
are fixed by API Contracts 2 and 3 and asserted by name in the tests of
Testing Approach.

### Invariants

```text
1. No command reads a file inside a required external skill's folder or a computedHash of an external lock entry.
2. Every result Doctor, skills check and baseline update give for an owned skill keeps its bytes at ad77a797.
3. A repository whose every required external skill has its folder and its lock entry gets the same Doctor skills status and the same baseline update state as at ad77a797 when no copy differs from its lock or its snapshot.
4. Every lock refusal and its message keeps its bytes.
```

### Declared breaks

Measured with a prototype of each Task on scratch copies at `ad77a797`: the
`skills`, `internal/baseline` and `internal/cli` suites, with and without
`docscontract`, failed only where listed here.

- Behavior, task_01: Doctor never appends `DR-SKILL-DIFFERS-FROM-SNAPSHOT`
  or `snapshot comparison unavailable`; `baseline update` drops the JSON
  member `skills.differing` and the text lines `Skills differing: <n>` and
  `- differs <name>: <reason>`, also under `--no-skills`.
- Behavior, task_02:
  - an external copy whose folder hash differs from its `computedHash` is
    ready: Doctor no longer prints `outdated: <name>` or its install action,
    `skills check` and the preview no longer list it, and `baseline update
    --yes` no longer restores it;
  - a lock entry whose `computedHash` is not 64 lowercase hexadecimal digits
    is no longer refused;
  - a symbolic link or a special file inside an external skill folder is no
    longer refused, because no file of the folder is read;
  - the preview message `<n> required skill(s) are missing or no longer
    match skills-lock.json: <names>` becomes `<n> required skill(s) are
    missing: <names>`.
- Tests removed by task_01, with `docs/references/coverage-record.json`
  re-recorded: `TestDoctorNamesASetupSnapshotDifferenceWithoutWarning`,
  `TestDoctorSnapshotComparisonPreservesReadinessAndFiltersSkills` and
  `TestBaselineUpdateKeepsASetupSnapshotDifference` (the whole
  `internal/cli/doctor_trailing_skills_test.go`);
  `TestSetupSnapshotDifferencesCompareTheInstalledTreeWithTheSnapshot` (the
  whole `internal/baseline/skills_trailing_test.go`);
  `TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot` (the whole
  `internal/baseline/repository_skills_snapshot_test.go`), whose check
  `TestEveryVendoredCatalogSkillIsTheTreeItsCatalogCommitPublishes` already
  makes for every lock entry through the same `portableRestoreDigest`; and
  `TestSetupSnapshotDifferencesCompareARepositoryProfile`.
- Tests renamed by task_01, in the same re-record:
  `TestDoctorOwnedLockWarningKeepsSnapshotDetailAndNextActions` →
  `TestDoctorOwnedLockWarningKeepsReadinessDetailAndNextActions`;
  `TestBaselineUpdatePreviewListsADifferenceWithARepositoryProfile` →
  `TestBaselineUpdatePreviewIgnoresADifferenceWithARepositoryProfile`;
  `TestDoctorComparesARepositoryProfileWithItsSnapshot` →
  `TestDoctorChecksARepositoryProfileSkillOnlyForPresence`;
  `TestBaselineUpdatePreviewListsASetupSnapshotDifference` →
  `TestBaselineUpdatePreviewIgnoresASetupSnapshotDifference`.
- Tests changed by task_01: `TestSkillSnapshotProfileNamesTheRepositoryPath`
  loses its `SetupSnapshotDifferences` assertion;
  `TestBaselineUpdateSkillsStageRefreshesSkillsWithARepositoryProfile`
  loses its `snapshotDifferences` dependency and `Differing` assertion.
- Tests changed by task_02, with no rename:
  - `skills/repository_test.go`: rows `changed external file` and `added
    external file` of `TestCheckRepositoryClassifiesMissingAndOutdatedSkills`;
    row `invalid required hash` of
    `TestCheckRepositoryRejectsMalformedLockAndUnsafeRequiredNames`, removed;
    the external nested link and special entry of
    `TestCheckRepositoryHandlesNestedLinksSpecialEntriesAndStableOrdering`;
    `TestCheckRepositoryWithExternalAcceptsEmptyRequirement`;
  - `internal/cli/doctor_test.go`: the mixed row of the Doctor readiness
    table that combined a missing owned skill, a missing external skill, an
    owned skill below its minimum and an outdated `testing-boss`;
  - `internal/cli/doctor_skills_lock_test.go`: the `external failure` row of
    `TestDoctorOwnedLockWarningKeepsReadinessDetailAndNextActions`, whose
    failure came from an outdated `testing-boss`, now comes from a second
    required external skill, `agentic-cli-design`, with no folder and no lock
    entry;
  - `internal/cli/baseline_update_test.go`: the fake readiness of
    `TestBaselineUpdateSkillsStageDegradesUnreachableSourcePerSkill` reports
    its two skills as missing instead of outdated;
  - `internal/cli/baseline_update_missing_skills_test.go`: rows of
    `TestBaselineUpdatePreviewNamesAMissingRequiredSkill` (the message, and
    `folder no longer matches its lock entry` now `current`) and
    `TestMissingRequiredSkillsCollectsNamedEvidence`;
  - `internal/cli/live_documents_docscontract_test.go`:
    `TestThisRepositoryHoldsEveryRequiredExternalSkill` keeps the lock rule
    itself (API Contract 6).
- No test is renamed or removed by task_02, so the coverage record is
  re-recorded once, by task_01.

## Coverage Map

- Goal 1, Core Feature 1, Success Metrics 1 and 2 → API Contracts 1 and 2,
  `TestCheckRepositoryChecksAnExternalSkillOnlyForPresence`,
  `TestDoctorChecksAnExternalSkillOnlyForPresence`,
  `TestDoctorKeepsAnExternalSkillThatDiffersFromItsLockReady`, Invariant 1.
- Goal 2, Core Features 2 and 3, Success Metrics 2 and 3 → API Contract 3,
  `TestBaselineUpdateNeverRestoresAPresentExternalSkill`,
  `TestBaselineUpdateNeverRestoresAnExternalSkillThatDiffersFromItsLock` and
  the rewritten repository-profile tests of Declared breaks.
- Goal 3, Success Metric 4 → API Contracts 1 to 3, the `missing directory`
  and `no lock entry` rows, `TestBaselineUpdatePreviewNamesAMissingRequiredSkill`,
  Invariant 4.
- Goal 4, Success Metric 5 → API Contract 4,
  `TestSkillsRestoreOfflinePreviewApplyAndIdempotence` and
  `TestBaselineSkillsRestoreCommand`, unchanged.
- Goal 5, Success Metric 6 → API Contract 5, Invariant 2, the existing
  owned-skill tests.
- Core Feature 4 → API Contract 6.
- Goal 6, Core Feature 5, Success Metric 7 → API Contract 7.

Deviations from the PRD: None.

## Absence Evidence

1. The three new test files — command: `ls skills/external_presence_test.go internal/cli/external_skill_presence_test.go internal/cli/external_skill_lock_test.go`;
   result: no such files.
2. `requiredExternalLockEntries`, a lock check that reads no hash — command:
   `grep -rn 'requiredExternalLockEntries\|computedHash' skills/repository.go`;
   result: no `requiredExternalLockEntries`; `requiredExternalHashes` reads
   and validates `computedHash` for every required external skill.
3. A lock hash check in this repository's live Repository Skill Set test —
   command: `grep -n 'SkillFolderHash\|computedHash' internal/cli/live_documents_docscontract_test.go`;
   result: no match; the test relies on `OutdatedExternal`.

## ADR Alternatives

- ADR-0364: the decision this Spec records; it rejects keeping the lock hash,
  a folder-only presence check, a network comparison of owned skills and an
  always-empty `skills.differing`, and the design takes none of them.
- ADR-0221: this Spec drops the comparison ADR-0221 added and the lock-hash
  rule it kept; ADR-0221 rejected no alternative, so ADR-0364 amends it and
  task_03 writes the notice.
- ADR-0264: ADR-0264 rejected keeping a warning for a difference, recording
  earlier snapshot digests and comparing upstream ancestry; the design takes
  none, and goes further than its informational difference by not comparing
  at all, so ADR-0364 amends it and task_03 writes the notice; its
  divergence order stays.
- ADR-0224: ADR-0224 rejected a hand-pinned digest constant; the design keeps
  that rejection and both of its rules for this repository, so only its
  reference to the Doctor Command is amended.
- ADR-0241: ADR-0241 rejected binding a repository profile to the closest
  built-in profile and a `setup` field in the profile; the design takes
  neither, and only the snapshot comparison it names among its readers goes,
  so ADR-0364 amends it and task_03 writes the notice.

- ADR-0360: this repository's Catalog Skill Bump stays the only writer of
  its Vendored Skills, their lock entries and the Setup Snapshots; no
  alternative it rejected is taken.
- ADR-0322: the Python profile keeps importing its upstream preset; no
  alternative it rejected is taken.

## Integration Points

None external. The checks read local files; no test reaches the network.

## Testing Approach

All tests are hermetic: `t.TempDir()` repositories built from the embedded
owned bundle, fake external folders and lock files, the existing
`runCLI`, `updateCommandDependenciesForTest` and
`runBaselineUpdateSkillsStageWith` seams, and no network. Every new
top-level test calls `t.Parallel()` first.

- `skills/external_presence_test.go` (task_02):
  `TestCheckRepositoryChecksAnExternalSkillOnlyForPresence` over
  `agentic-cli-design` in `writeReadyRepositoryFixture`, rows `changed file`,
  `added file`, `nested symbolic link` and `lock entry without computedHash`
  ready; `missing directory` in `MissingExternal`; `no lock entry` a
  `RepositoryReadinessError` whose `MissingExternal` names it; `skill root is
  a symbolic link` refused.
- `internal/cli/external_skill_presence_test.go` (task_01):
  - `TestDoctorChecksAnExternalSkillOnlyForPresence` over
    `writeDoctorReadyRepositoryFixture` with profile `go-cli-tui` and the one
    external skill `testing-boss`, whose fixture copy matches its lock and
    differs from its Setup Snapshot: rows `differs from its Setup Snapshot`
    (the exact `ok` line of API Contract 2, no
    `DR-SKILL-DIFFERS-FROM-SNAPSHOT`, exit `0`) and `missing` (the exact
    `failed` line, exit `1`); Doctor writes nothing.
  - `TestBaselineUpdateNeverRestoresAPresentExternalSkill` through
    `runBaselineUpdateSkillsStageWith` with the real
    `skills.CheckRepositoryWithExternal`: `present copy that differs from its
    Setup Snapshot` restores nothing, keeps the copy's bytes and encodes no
    `differing`; `missing copy` restores `testing-boss`.
- `internal/cli/external_skill_lock_test.go` (task_02), on the same fixture
  with `testing-boss`'s `SKILL.md` rewritten after its lock entry was
  recorded:
  - `TestDoctorKeepsAnExternalSkillThatDiffersFromItsLockReady`: the exact
    `ok` line, no `outdated:`, exit `0`;
  - `TestBaselineUpdateNeverRestoresAnExternalSkillThatDiffersFromItsLock`:
    the skills stage restores nothing and keeps the copy's bytes.
- Existing suites: `./skills`, `./internal/cli`, `./internal/baseline`, with
  and without `-tags docscontract`, and `-tags docscontract
  ./internal/docscontract`.
- Sabotages, each run by the author on a scratch copy of the prototype at
  `ad77a797`, with the case whose input makes it observable:
  - the folder hash comparison back in the external loop:
    `TestCheckRepositoryChecksAnExternalSkillOnlyForPresence` fails on
    `changed file`, and `TestDoctorKeepsAnExternalSkillThatDiffersFromItsLockReady`
    fails;
  - the lock entry requirement dropped: it fails on `no lock entry`;
  - the snapshot suffix back in Doctor:
    `TestDoctorChecksAnExternalSkillOnlyForPresence` fails on `differs from
    its Setup Snapshot`;
  - `baselineUpdateExternalDrift` returning every required external skill:
    `TestBaselineUpdateNeverRestoresAPresentExternalSkill` fails on `present
    copy that differs from its Setup Snapshot`;
  - the lock hash check removed from this repository's live test, with one
    required skill's `computedHash` replaced by 64 zeros in a temporary
    copy: `TestThisRepositoryHoldsEveryRequiredExternalSkill` passes, and
    with the check in place it fails naming the skill.

## Build Order

1. The snapshot comparison removed from Doctor and the managed refresh,
   `skills.differing` removed, `internal/baseline/skills_trailing.go` and its
   tests deleted, the repository-profile and owned-lock tests rewritten,
   `internal/cli/external_skill_presence_test.go`, and the coverage record
   re-recorded last (task_01).
2. Presence in `skills/repository.go`; the `outdated` detail, reason and
   restore removed from `internal/cli`; this repository's lock rule moved
   into its live test; `skills/external_presence_test.go` and
   `internal/cli/external_skill_lock_test.go` (task_02) (depends on: 1).
3. The guides, the README, the Roundfix references with the mirrors synced
   and the owned-skill versions recorded, the two glossary entries and the
   four amendment notices, with the Behavior Surface Record re-recorded last
   (task_03) (depends on: 2).
4. QA gate (task_04) (depends on: 3).

task_01 goes first because it deletes `snapshotComparableSkills`, the reader
of `OutdatedExternal` that task_02 would otherwise have to keep compiling;
task_02 removes the field and changes the same CLI files. task_03 documents
what both make true. No Verification depends on another Spec.

## Risks & Considerations

- A hand-edited external skill is no longer reported. ADR-0264's divergence
  order still tells a repository how to correct a vendored skill, and this
  repository keeps failing its own tests on one.
- A Fluxus or Vortex lock written by the skills CLI keeps its entries, so
  presence holds there; a repository that copied a skill folder by hand
  without a lock entry still gets the lock refusal it gets today.
- Removing `skills.differing` is a JSON field removal: the release notes name
  it. No skill or Roundfix command reads it (`grep -rn differing` over
  `.agents/skills`, `skills` and `internal`, outside the removed code).
- Specs 0352 and 0353 also raise the Roundfix skill version; the record
  command picks the next free version at task_03's run.

## Glossary

- changes: **Missing Required Skill**
- changes: **Setup Snapshot Difference**

## Decisions

- Presence is a lock entry and a directory, read by the existing check; no
  new function walks the folder (ADR-0364).
- The lock refusals stay, so a repository changes state only through content
  it already had.
- The owned rule is not touched: it already compares offline with the
  binary's copy, and an equal digest implies an equal version.
- `skills.differing` is removed rather than kept empty (ADR-0364).
- `baseline.SetupSnapshotDifferences` is deleted with its tests, and
  `TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot` with it,
  because `TestEveryVendoredCatalogSkillIsTheTreeItsCatalogCommitPublishes`
  already compares every lock entry's tree with its snapshot through the same
  digest; no production code stays for tests only.
- This repository's lock rule moves into its live test through the existing
  `skills.SkillFolderHash`.
- Secondbrain: `wiki/index.md` read; the `qmd query` named in `_prd.md` →
  Validation found the Fluxus record of 2026-10-06 on the same conflict. Exa:
  the skills CLI's lock format and issue 806.
- ADR-0364 records the rule and amends ADR-0221, ADR-0224, ADR-0241 and
  ADR-0264. Three implementation Tasks in series and a QA gate.
