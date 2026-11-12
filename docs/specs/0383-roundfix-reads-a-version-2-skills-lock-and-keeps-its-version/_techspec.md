---
spec: 0383-roundfix-reads-a-version-2-skills-lock-and-keeps-its-version
status: active
created: 2026-11-12
---

# Roundfix reads a version 2 skills lock and keeps its version — Technical Spec

## Executive Summary

Two readers refuse every `skills-lock.json` but version 1, and four places
compute a `computedHash` with the version 1 digest only. This Spec adds the
skills CLI 1.8 digest and one version parser to `internal/skillhash`, makes
both readers accept versions 1 and 2 through that parser, keeps the parsed
version on the Baseline lock document so each hash is computed in its
version's algorithm, and changes the two next actions that answer a refused
version. The guides, the Roundfix skill references and the glossary follow.
It covers `_prd.md` → Goals, User Stories 1–3, Core Features 1–4 and Success
Metrics 1–3; ADR-0393 records the decision.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the lock's existing `version` field is read. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files only, nothing
  fetches. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0393 (this Spec's decision),
  ADR-0359 and ADR-0350 applied; ADR-0364, ADR-0189, ADR-0191, ADR-0246 and
  ADR-0380 not applicable, see `## ADR Alternatives`. Source:
  `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — no protected tooling mutation proposed or
  authorized: no Makefile, workflow, `go.mod`, `go.sum`, lint or formatter
  change. The Governed Paths task_03 touches are the owned Roundfix skill
  files and the canonical `setup`, `baseline` and `runtime` references,
  measured with `speccheck.GovernedPath`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0383-roundfix-reads-a-version-2-skills-lock-and-keeps-its-version/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `.agents/skills/roundfix/references/runtime.md`,
  `.agents/skills/roundfix/references/setup.md`,
  `skills/roundfix/SKILL.md`.

## Current behavior

At `d79b30b4`:

- `readLocalSkillsLock` (`skills/repository.go`) decodes `version` as an
  `int` and returns the readiness error `validate skills lock version
  "<path>": got <n>, want 1` for any other value. It is read only when an
  external skill is required (`CheckRepositoryWithExternal`), by Doctor, the
  release plan's `skills` check, `roundfix skills check` and the missing
  skill list of `baseline update`.
- `loadSkillsLock` (`internal/baseline/skills_restore_lock.go`) refuses a
  lock whose `version` is not 1 with the finding `lock.invalid`, message
  `skills-lock.json must use version 1 with a skills object.`; an absent lock
  yields an empty version 1 document. Its callers are
  `RecordInstalledOwnedSkills` and `StaleOwnedLockEntries`
  (`owned_skills_lock.go`), `InstalledRetiredSkills` and
  `DeleteRetiredSkills` (`retired_skills_installed.go`), the external skill
  restore (`skills_restore.go`) and `ReconcileSkillsLock`
  (`skills_reconcile.go`).
- Every hash is `externalSkillsLockDigest` → `skillhash.Sum`: paths in en-US
  collation order, then for each file the path and the content, no
  separator. It is written by `RecordInstalledOwnedSkills` and
  `expectedSkillsLockEntry`, and compared by `StaleOwnedLockEntries` and
  `InstalledRetiredSkills`.
- `doctorSkillReadinessResult` (`internal/cli/doctor.go`) names
  `externalSkillsNextAction`, `bunx skills experimental_install && bunx
  skills update -p -y`, for every external failure that names no skill,
  including a lock it cannot read. `baseline update` names `repair
  skills-lock.json and rerun roundfix baseline update` at both
  `InstalledRetiredSkills` failures.
- The review scope (`internal/cli/review.go`) reads only the lock's `skills`
  object and the Makefile's `skills-third-party-bump` reads it with `jq`;
  neither reads the version. This repository's lock is version 1 and the
  Catalog Skill Bump writes it only through `roundfix baseline skills
  restore`.

The measurements of `_prd.md` → Validation reproduce the reported lines.

## System Architecture

No new package or seam. `internal/skillhash`, already imported by `skills`
and `internal/baseline`, gains the lock version parser and the version 2
digest beside `Sum`. The Baseline lock document carries the parsed version,
and its callers ask it for the digest instead of calling
`externalSkillsLockDigest`. Doctor and `baseline update` recognise the
refusal with `errors.Is` and choose their next action.

## Implementation Design

### Interfaces

```go
// internal/skillhash/lock.go
var ErrUnsupportedLockVersion = errors.New("unsupported skills-lock.json version")
type UnsupportedLockVersionError struct{ Version string } // as written in the file
func (UnsupportedLockVersionError) Error() string           // API Contract 1
func (UnsupportedLockVersionError) Is(target error) bool   // ErrUnsupportedLockVersion
func ParseLockVersion(text string) (int, error)            // 1 or 2, else the error above
func SumForLockVersion(version int, files []File) (string, error) // 1: Sum; 2: API Contract 3

// internal/baseline/skills_restore_lock.go
type skillsLockDocument struct { before []byte; root orderedJSONValue; version int }
func (skillsLockDocument) digest(files []restoreFile) (string, error)
```

`externalSkillsLockDigest` stays for the version 1 compatibility fixture
check (`validateLockAdapter`); `skills.SkillFolderHash` stays the version 1
digest of a folder.

### Data Models

None changed. `skills-lock.json` keeps its bytes except the `computedHash`
values and `ref` members the existing writers already change; its `version`
is never written.

### API Contracts

1. **Lock refusals.** `skillhash.UnsupportedLockVersionError.Error()` is
   `skills-lock.json version <N> is not supported; this Roundfix reads versions 1 and 2`,
   where `<N>` is the number as the file writes it. `skills/repository.go`
   returns it as the readiness error of operation `validate skills lock
   version` (ownership external), so it prints `validate skills lock version
   "<path>": skills-lock.json version <N> is not supported; this Roundfix
   reads versions 1 and 2`; a missing `version` there prints `validate skills
   lock version "<path>": skills-lock.json must hold a version number and a
   skills object`. `loadSkillsLock` returns the finding code
   `lock.unsupported-version`, message `<Error()>.` and action `Run roundfix
   upgrade; a newer Roundfix may read this version.`, wrapping the error; a
   lock that is not an object, or lacks a numeric `version` or a `skills`
   object, keeps the code `lock.invalid` with the message `skills-lock.json
   must hold a version number and a skills object.`. Both match
   `errors.Is(err, skillhash.ErrUnsupportedLockVersion)` only for the
   version refusal.
2. **Next actions.** When the Doctor `skills:` readiness error matches
   `ErrUnsupportedLockVersion`, its next action is exactly `roundfix
   upgrade`; every other next action is unchanged. With a version 3 lock and
   one required external skill, the line is:

   ```text
   skills: failed (validate skills lock version "<path>": skills-lock.json version 3 is not supported; this Roundfix reads versions 1 and 2; skills-lock.json comparison unavailable: skills-lock.json version 3 is not supported; this Roundfix reads versions 1 and 2.; next: roundfix upgrade)
   ```

   `roundfix baseline update` reports, at both `InstalledRetiredSkills`
   failures, the next action `run roundfix upgrade, then rerun roundfix
   baseline update` for that error and `repair skills-lock.json and rerun
   roundfix baseline update` for every other.
3. **Version 2 digest.** Over the given relative paths and contents: drop a
   file whose path has a component `metadata.json`, or a directory component
   `.git`, `__pycache__` or `__pypackages__`; turn each CRLF into LF in a
   content without a NUL byte; sort by the UTF-8 bytes of the slash path;
   write into one SHA-256, per file, the path, a NUL byte, the decimal length
   of the normalised content, a NUL byte and the content; return lowercase
   hex. It equals the README test vectors quoted in `_prd.md` → Acceptance
   evidence.
4. **Guide text.** task_03 writes the sentences in its Requirements verbatim;
   the Verification of task_03 checks them.

### Surface Transcripts

None. Doctor prints one changed line among checks whose output depends on
the machine, so no whole command transcript is stable; API Contract 2 fixes
that line verbatim, and `TestDoctorNamesRoundfixUpgradeForAnUnsupportedSkillsLockVersion`
asserts it whole.

## Coverage Map

- User Story 1, Core Feature 1, Success Metric 1 → API Contract 1, the two
  readers (task_01) and `TestTheSkillsCommandsWorkOnAVersionTwoSkillsLock`,
  `TestBaselineUpdateNamesRoundfixUpgradeForAnUnsupportedSkillsLockVersion/version_2`
  (task_02).
- User Story 2, Core Feature 2, Success Metrics 1–2 → `SumForLockVersion`,
  `skillsLockDocument.digest` and the four callers (task_01).
- User Story 3, Core Features 1 and 3, Success Metric 3 → API Contracts 1–2
  (task_01, task_02).
- Core Feature 4 → API Contract 4 (task_03).

Deviations from the PRD: None.

## Absence Evidence

1. The version 2 digest and the lock version parser — command: `grep -rn 'NUL\|__pypackages__\|ParseLockVersion\|SumForLockVersion\|version 2' internal/skillhash skills internal/baseline --include='*.go'`; result: no match at `d79b30b4`; `internal/skillhash` holds only `Sum`, the version 1 digest.
2. A lock version kept on the Baseline lock document — command: `grep -n 'version' internal/baseline/skills_restore_lock.go`; result: at `d79b30b4` only the check `version.isOne()` and the empty document's literal `"1"`; nothing stores the version.
3. A next action naming `roundfix upgrade` for a lock — command: `grep -rn 'roundfix upgrade' internal/cli/doctor.go internal/cli/baseline_update.go`; result: at `d79b30b4` only the newer-Baseline refusal of `baseline update` (`refuseNewerBaselineRepository`); Doctor's `skills:` line never names it.

## ADR Alternatives

- ADR-0393: this Spec's decision; it takes none of the options its
  Considered Options reject (accept any shape-compatible version, migrate a
  lock, write a version 1 hash into a version 2 lock, keep the skills command
  as the next action).
- ADR-0359: applied; the bundle install still removes `ref` and writes the
  installed tree's hash, now in the lock version's algorithm, and its
  rejected options (write a catalog commit as `ref`, change `sourceType`,
  fail Doctor on a stale entry) stay rejected.
- ADR-0350: applied; a Retired Skill is still deleted only when its copy
  hashes to its lock entry, now in the lock version's algorithm; none of its
  rejected options is taken.
- ADR-0364: not applicable; a required external skill is still checked
  only for presence, and only the version the presence reader accepts
  changes.
- ADR-0189: not applicable; the owned skill's version and bundle minimum
  are unchanged.
- ADR-0191: not applicable; the Baseline's skill set rules are unchanged.
- ADR-0246: not applicable; retired and superseded by ADR-0350.
- ADR-0380: not applicable; a skill the catalog replaces is still retired
  and its replacements required as before; only the hash that proves a
  Retired Skill unedited follows the lock version.

## Integration Points

None outside the repository. The skills CLI is never run by Roundfix; its
published digest is reproduced in Go and pinned by test vectors.

## Testing Approach

task_01:

- `internal/skillhash/lock_test.go` (new):
  `TestSumForLockVersionMatchesThePublishedSkillsCLI` — the three README
  vectors and the crafted folder of `_prd.md` → Validation for version 2
  (`e5adbe14…`, with `node_modules/pkg/index.js` in the input), and the same
  folder without `node_modules` for version 1 (`416d74f3…`);
  `TestParseLockVersionReadsOnlyVersionsOneAndTwo` — `1` and `2` parse, `3`,
  `0`, `1.5` and `-1` refuse with API Contract 1's text and
  `ErrUnsupportedLockVersion`, and `SumForLockVersion(3, …)` refuses.
- `skills/repository_lock_version_test.go` (new):
  `TestCheckRepositoryReadsASkillsLockOfVersionTwo` and
  `TestCheckRepositoryRefusesAnUnsupportedSkillsLockVersionByName` (the exact
  readiness error).
- `internal/baseline/skills_lock_version_test.go` (new):
  `TestOwnedSkillLockEntriesUseTheHashOfTheLockVersion` (subtests `version
  1`, `version 2`: the other version's hash is stale, the record writes the
  own version's hash and nothing else, then nothing is stale),
  `TestARetiredSkillIsProvenUneditedByTheHashOfTheLockVersion`,
  `TestSkillsRestoreWritesTheHashOfTheLockVersion` (a version 2 lock, a
  source with CRLF text and `metadata.json`) and
  `TestEveryLockReaderRefusesAnUnsupportedVersionByName`.

task_02, `internal/cli/skills_lock_version_test.go` (new):
`TestTheSkillsCommandsWorkOnAVersionTwoSkillsLock` (the pantheon shape: a
version 2 lock whose `roundfix` entry has a `ref` and a stale hash; Doctor
warns `DR-SKILL-LOCK-STALE`; `skills install --target project` exits 0 and
rewrites the lock to the expected bytes, `"version": 2` kept; Doctor is
`ok`), `TestDoctorNamesRoundfixUpgradeForAnUnsupportedSkillsLockVersion`
(API Contract 2's line) and
`TestBaselineUpdateNamesRoundfixUpgradeForAnUnsupportedSkillsLockVersion`
(subtests `version 2`, exit 0, and `version 3`, exit `exitRunFailed` with the
upgrade next action).

Declared breaks, each a pinned expectation that used version 2 as the
invalid version and now uses version 3:
`TestCheckRepositoryRejectsMalformedLockAndUnsafeRequiredNames/wrong_version`
(`skills/repository_test.go`),
`TestRecordInstalledOwnedSkillsRefusesAnUnsafeLock/schema`
(`internal/baseline/owned_skills_lock_test.go`) and
`TestDoctorKeepsReadinessWhenOwnedLockComparisonIsUnavailable/schema`
(`internal/cli/doctor_skills_lock_test.go`).

Sabotage proofs, each run on the prototype: dropping the CRLF rule fails
`TestSumForLockVersionMatchesThePublishedSkillsCLI/version_2_edge_folder`;
returning `Sum` for version 2 fails its three version 2 vector subtests;
accepting 3 instead of 2 fails `TestParseLockVersionReadsOnlyVersionsOneAndTwo`;
making `skillsLockDocument.digest` always use version 1 fails the `version 2`
subtests of the owned-entry and Retired Skill tests and
`TestSkillsRestoreWritesTheHashOfTheLockVersion`; making `skills/repository.go`
refuse version 2 fails `TestCheckRepositoryReadsASkillsLockOfVersionTwo`;
making `loadSkillsLock` refuse version 2 fails
`TestTheSkillsCommandsWorkOnAVersionTwoSkillsLock`; dropping Doctor's
upgrade next action fails
`TestDoctorNamesRoundfixUpgradeForAnUnsupportedSkillsLockVersion`; and
restoring the repair next action for that error fails
`TestBaselineUpdateNamesRoundfixUpgradeForAnUnsupportedSkillsLockVersion/version_3`.
New top-level tests call `t.Parallel()` first; none reaches the network,
runs the skills CLI or writes outside its temporary directories.

## Build Order

1. The digest and version parser in `internal/skillhash`, both lock readers,
   the lock document's version and the four hash call sites, with their
   tests and two moved expectations.
2. Doctor's and `baseline update`'s next actions and the command-level tests,
   with the third moved expectation (depends on: 1, which supplies the error
   and the version 2 behavior the tests drive).
3. The guides, the Roundfix skill references, the glossary entry, the skill
   version and the Behavior Surface Record (depends on: 2, so the text
   describes the delivered next actions).
4. Final QA gate (depends on: 1, 2, 3).

No step needs another Spec's code. Specs 0384 and 0385 may also raise the
Roundfix skill version and edit `CONTEXT.md`; step 3 uses the record command
and edits only **Stale Owned Lock Entry**.

## Risks & Considerations

- A lock version 2 entry written by the CLI for a skill folder with a
  symbolic link inside has a hash Roundfix cannot reproduce; Roundfix refuses
  that folder as today.
- `inspectRestoreTarget` collects every regular file, `node_modules` and
  `.git` included; the version 2 digest drops `.git` itself, and the version 1
  digest keeps today's behavior for each caller.
- A skills CLI release that changes the version 2 digest without a new
  version number breaks the hash agreement; the README vectors pinned in
  `internal/skillhash` record the contract this Spec relies on.

## Glossary

- changes: **Stale Owned Lock Entry**

## Decisions

- Versions 1 and 2 are read, others refused by name; see ADR-0393.
- The parsed version travels on the Baseline lock document; the readiness
  reader only validates it.
- The refusal text lives once, in `internal/skillhash`.
- Three implementation Tasks and the QA gate.
