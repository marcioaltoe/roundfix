---
task: task_01
spec: 0383-roundfix-reads-a-version-2-skills-lock-and-keeps-its-version
status: completed
type: backend
complexity: medium
---

# Task 01: Every skills lock reader accepts versions 1 and 2 and hashes in the lock's own version

## Overview

The `skills` CLI 1.8 writes `skills-lock.json` as version 2 with a new folder
hash, and Roundfix refuses any lock but version 1 and hashes only in the
version 1 algorithm (vortex and pantheon, 2026-11-12,
[`references/2026-11-12-vortex-skills-lock-version-2-refused.md`](references/2026-11-12-vortex-skills-lock-version-2-refused.md),
reported by vortex;
[`references/2026-11-12-pantheon-skills-lock-version-2-refused.md`](references/2026-11-12-pantheon-skills-lock-version-2-refused.md),
reported by pantheon). This Task adds the version 2 digest and one version
parser to `internal/skillhash`, makes both lock readers accept versions 1 and
2 and refuse any other by name, and makes every `computedHash` Roundfix writes
or compares follow the version the lock declares, as `_techspec.md` →
Interfaces and API Contracts 1 and 3 define.

## Requirements

1. MUST create `internal/skillhash/lock.go` with the exported API of
   `_techspec.md` → Interfaces: `ErrUnsupportedLockVersion`,
   `UnsupportedLockVersionError` (field `Version string`, its `Error()`
   returning exactly
   `skills-lock.json version <Version> is not supported; this Roundfix reads versions 1 and 2`,
   and an `Is` method matching `ErrUnsupportedLockVersion`),
   `ParseLockVersion(text string) (int, error)`, which returns 1 or 2 for a
   number whose value is 1 or 2 and the error for any other text, and
   `SumForLockVersion(version int, files []File) (string, error)`, which
   returns `Sum(files)` for version 1, the digest of `_techspec.md` → API
   Contract 3 for version 2, and the error for any other version. `Sum` and
   `hash.go` stay unchanged.
2. MUST change `readLocalSkillsLock` (`skills/repository.go`) to decode
   `version` as a `json.Number`; a missing version returns the readiness
   error of operation `validate skills lock version` with the error text
   `skills-lock.json must hold a version number and a skills object`, and
   any version `ParseLockVersion` refuses returns that operation with the
   parser's error, so `errors.Is(err, skillhash.ErrUnsupportedLockVersion)`
   holds. Ownership stays external, and versions 1 and 2 read alike.
3. MUST change `loadSkillsLock` (`internal/baseline/skills_restore_lock.go`):
   a lock that is not an object or lacks a numeric `version` or a `skills`
   object returns code `lock.invalid` with the message
   `skills-lock.json must hold a version number and a skills object.`; a
   version `ParseLockVersion` refuses returns code `lock.unsupported-version`,
   the message `<the parser's Error()>.`, the action
   `Run roundfix upgrade; a newer Roundfix may read this version.` and the
   parser's error as the wrapped error. The returned `skillsLockDocument`
   MUST carry the parsed version (1 for the empty document of an absent
   lock), `clone` MUST keep it, and a new method `digest(files []restoreFile)
   (string, error)` MUST return `skillhash.SumForLockVersion` of that
   version. Remove `isOne`, which nothing else uses.
4. MUST replace the version 1 digest with the document's `digest` in
   `RecordInstalledOwnedSkills` and `StaleOwnedLockEntries`
   (`owned_skills_lock.go`), `InstalledRetiredSkills`
   (`retired_skills_installed.go`) and `expectedSkillsLockEntry`
   (`skills_restore.go`, which receives the lock document as its first
   parameter), returning a digest error wrapped with the operation. No writer
   writes the `version` member; `validateLockAdapter` keeps
   `externalSkillsLockDigest`.
5. MUST create the tests of `_techspec.md` → Testing Approach for this Task,
   each top-level test calling `t.Parallel()` first:
   `internal/skillhash/lock_test.go` with
   `TestSumForLockVersionMatchesThePublishedSkillsCLI` and
   `TestParseLockVersionReadsOnlyVersionsOneAndTwo`;
   `skills/repository_lock_version_test.go` with
   `TestCheckRepositoryReadsASkillsLockOfVersionTwo` and
   `TestCheckRepositoryRefusesAnUnsupportedSkillsLockVersionByName`; and
   `internal/baseline/skills_lock_version_test.go` with
   `TestOwnedSkillLockEntriesUseTheHashOfTheLockVersion`,
   `TestARetiredSkillIsProvenUneditedByTheHashOfTheLockVersion`,
   `TestSkillsRestoreWritesTheHashOfTheLockVersion` and
   `TestEveryLockReaderRefusesAnUnsupportedVersionByName`. The version 2
   subtests of `TestSumForLockVersionMatchesThePublishedSkillsCLI` pin the
   README vectors `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`
   (no file), `93c36952d8843b33e68ae970ebabcc03e1687c9bb0dd1ecaf357060318952f55`
   (`SKILL.md` = `# hi\n`) and
   `b22c4d22aa319d2924312fa18fcbfedc88fa89c74dbd8f298c683e1caf7b7f24`
   (`SKILL.md` = `a\n`, `docs/b.md` = `b\n`), and the crafted folder of
   `_prd.md` → Validation (`SKILL.md` = `---\nname: edge\n---\nbody\r\nline two\r\n`,
   `metadata.json` = `meta\n`, `__pycache__/x.pyc` = `cache\n`,
   `references/B.md` = `upper\n`, `references/a.md` = `lower\n`,
   `_x.md` = `under\n`, `node_modules/pkg/index.js` = `dep\n`) at
   `e5adbe14f1b9b04eb73f28a1fc7e58f083b72328dbb7271c4a2813cbcadb42ac`; its
   version 1 subtest pins the same folder without `node_modules/pkg/index.js`
   at `416d74f34b15134cda0d15e3b9e7dd32f3b0fb6aaf9369c51d9931b37138dc73`.
   Each new file starts with the `// Suite:`, `// Invariant:`,
   `// Boundary IN:` and `// Boundary OUT:` header its package's sibling
   tests carry.
6. MUST move the three pinned expectations that used version 2 as the
   invalid version to version 3, and change nothing else in them:
   `TestCheckRepositoryRejectsMalformedLockAndUnsafeRequiredNames/wrong_version`
   (`skills/repository_test.go`, `lock.Version = 3`),
   `TestRecordInstalledOwnedSkillsRefusesAnUnsafeLock/schema`
   (`internal/baseline/owned_skills_lock_test.go`,
   `{"version":3,"skills":{}}`) and
   `TestDoctorKeepsReadinessWhenOwnedLockComparisonIsUnavailable/schema`
   (`internal/cli/doctor_skills_lock_test.go`, `{"version":3,"skills":{}}`),
   plus every other test that pins a version 2 lock as invalid on the
   default branch when the Task runs.
7. MUST reach no network, run no skills CLI, start no Agent session, and
   write only inside temporary directories in every test.

## Subtasks

- [ ] Add the version parser, the refusal error and the version 2 digest to `internal/skillhash`.
- [ ] Accept versions 1 and 2 in both lock readers and carry the version on the Baseline lock document.
- [ ] Hash in the lock's version at the four call sites.
- [ ] Add the new tests and move the three pinned expectations.

## Acceptance Criteria

- [ ] A version 2 lock reads as a version 1 lock does, and a version 3 lock is refused by name in both readers.
- [ ] Every `computedHash` written into or compared with a lock uses that lock's version, and the version 2 digest equals the published vectors.
- [ ] No writer changes the lock's version, and a version 1 lock keeps today's hash.

## Context

- creates: `internal/skillhash/lock.go`
- creates: `internal/skillhash/lock_test.go`
- creates: `skills/repository_lock_version_test.go`
- creates: `internal/baseline/skills_lock_version_test.go`
- interface: `skills/repository.go`
- interface: `skills/repository_test.go`
- interface: `internal/baseline/skills_restore_lock.go`
- interface: `internal/baseline/owned_skills_lock.go`
- interface: `internal/baseline/retired_skills_installed.go`
- interface: `internal/baseline/skills_restore.go`
- interface: `internal/baseline/owned_skills_lock_test.go`
- interface: `internal/cli/doctor_skills_lock_test.go`
- instruction: `internal/skillhash/hash.go`
- instruction: `internal/skillhash/hash_test.go`
- instruction: `internal/baseline/skills_restore_test.go`
- instruction: `internal/baseline/retired_skills_installed_test.go`
- instruction: `docs/adr/0393-roundfix-reads-skills-lock-versions-1-and-2-and-keeps-the-version-it-found.md`
- instruction: `docs/agents/go.md`

## Consumers

- `skillhash.ParseLockVersion`, `SumForLockVersion` and the lock readers, read by Doctor, the release plan, `roundfix skills check`, `roundfix skills install --target project`, `roundfix baseline update` and `roundfix baseline skills restore` and `reconcile` — consumers: `./internal/skillhash`, `./skills`, `./internal/baseline`, `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestSumForLockVersionMatchesThePublishedSkillsCLI|TestParseLockVersionReadsOnlyVersionsOneAndTwo|TestCheckRepositoryReadsASkillsLockOfVersionTwo|TestCheckRepositoryRefusesAnUnsupportedSkillsLockVersionByName|TestCheckRepositoryRejectsMalformedLockAndUnsafeRequiredNames|TestOwnedSkillLockEntriesUseTheHashOfTheLockVersion|TestARetiredSkillIsProvenUneditedByTheHashOfTheLockVersion|TestSkillsRestoreWritesTheHashOfTheLockVersion|TestEveryLockReaderRefusesAnUnsupportedVersionByName|TestRecordInstalledOwnedSkillsRefusesAnUnsafeLock)$' ./internal/skillhash ./skills ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestSumForLockVersionMatchesThePublishedSkillsCLI/version_2_edge_folder TestSumForLockVersionMatchesThePublishedSkillsCLI/version_1_edge_folder TestParseLockVersionReadsOnlyVersionsOneAndTwo TestCheckRepositoryReadsASkillsLockOfVersionTwo TestCheckRepositoryRefusesAnUnsupportedSkillsLockVersionByName TestCheckRepositoryRejectsMalformedLockAndUnsafeRequiredNames/wrong_version TestOwnedSkillLockEntriesUseTheHashOfTheLockVersion/version_2 TestARetiredSkillIsProvenUneditedByTheHashOfTheLockVersion/version_2 TestSkillsRestoreWritesTheHashOfTheLockVersion TestEveryLockReaderRefusesAnUnsupportedVersionByName TestRecordInstalledOwnedSkillsRefusesAnUnsafeLock/schema; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the new tests do not exist, so the packages do not compile with them and their pass lines are missing; afterwards every named test and subtest passes.
- `grep -qF -- 'func SumForLockVersion' internal/skillhash/lock.go && grep -qF -- '{"version":3,"skills":{}}' internal/cli/doctor_skills_lock_test.go && go test -count=1 -timeout 30m ./internal/skillhash ./skills && go test -count=1 -timeout 30m -run 'Lock|Restore|Retired|Reconcile|OwnedSkill' ./internal/baseline && go test -count=1 -timeout 30m -run 'Doctor|Skills|Lock|Retired|OwnedSkill|MissingRequired|ReleasePlan' ./internal/cli && GOOS=windows go vet ./internal/skillhash ./skills ./internal/baseline ./internal/cli` — expected: exit 0. Before this Task `internal/skillhash/lock.go` does not exist, so the first check fails; afterwards the packages' lock, restore, retired-skill, Doctor, skills and release-plan tests pass and the packages vet for Windows.

## References

- `_prd.md` → Goals; User Stories 1-3; Core Features 1-2; Success Metric 1; Success Metric 2; Success Metric 3
- `_techspec.md` → Interfaces; API Contract 1; API Contract 3; Testing Approach; Build Order 1
- ADR-0393; ADR-0359; ADR-0350


## Result

Skills loaded: implement-task, coding-guidelines, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, no-workarounds, tech-writer, evidence-gate, systematic-debugging, golang-safety, golang-cli.

Implemented the Task 01 slice:

- Both readers accept numeric versions 1 and 2 and preserve unsupported-version errors through `errors.Is`. Baseline distinguishes malformed schema (`lock.invalid`) from unsupported numbers (`lock.unsupported-version`) with the specified message and upgrade action. The readiness reader preserves external ownership and rejects missing, quoted, null and boolean versions as malformed.
- `SumForLockVersion` keeps `Sum` unchanged for version 1 and implements the published version 2 filtering, CRLF normalization, UTF-8 path ordering and NUL/length framing. The parser compares exact numeric values so fractions cannot round into a supported version.
- The Baseline document and its clone retain the parsed version; an absent lock starts at version 1. Owned-entry recording/comparison, retired-skill proof and external restore use the document's digest and wrap digest errors with their operation. The compatibility adapter still uses the legacy digest.
- Added all eight specified top-level tests with parallel entry and suite headers. Moved only the invalid-version literal from 2 to 3 in each of the three existing expectations. Tests use temporary trees and offline Git source fixtures; no network, skills CLI or Agent session is needed.

Focused evidence, all exit 0 on the implementation:

1. `go test ./internal/skillhash ./skills -run '^Test(SumForLockVersionMatchesThePublishedSkillsCLI|ParseLockVersionReadsOnlyVersionsOneAndTwo|CheckRepositoryReadsASkillsLockOfVersionTwo|CheckRepositoryRefusesAnUnsupportedSkillsLockVersionByName|CheckRepositoryRejectsMalformedLockAndUnsafeRequiredNames)$' -count=1 -timeout=5m`
   - Acceptance 1: version 2 readiness equals version 1; version 3 and fractional versions produce the exact named readiness error and typed sentinel; missing/nonnumeric versions remain malformed.
   - Acceptance 2: all three README vectors, the version 2 edge folder, binary CRLF preservation and excluded components are exercised.
   - Acceptance 3: the pinned version 1 edge-folder digest remains unchanged and hashing preserves caller-owned input.
2. `go test ./internal/baseline -run '^Test(OwnedSkillLockEntriesUseTheHashOfTheLockVersion|ARetiredSkillIsProvenUneditedByTheHashOfTheLockVersion|SkillsRestoreWritesTheHashOfTheLockVersion|EveryLockReaderRefusesAnUnsupportedVersionByName|RecordInstalledOwnedSkillsRefusesAnUnsafeLock|RecordInstalledOwnedSkillsRewritesOnlyTheHashOfAnInstalledOwnedSkill|RecordInstalledOwnedSkillsLeavesAMatchingLockUntouched|SkillsRestoreOfflinePreviewApplyAndIdempotence|RetiredSkillDispositionsProveACopyUnedited|DeleteRetiredSkillsDeletesOnlyThePlannedCopies|ReconcileSkillsLockRemovesOnlyAbsentUnrequiredEntries)$' -count=1 -timeout=5m`
   - Acceptance 1: unsupported numbers retain the sentinel and exact Baseline finding code, message and action through the document, recorder, stale scan and retired scan; malformed shapes remain `lock.invalid`.
   - Acceptance 2: the other version's hash is stale and cannot prove a retired copy unedited; recording and offline restore write the lock version's digest, including CRLF text and metadata exclusion.
   - Acceptance 3: owned recording changes only the hash bytes, preserves the original `version` representation, and becomes a no-op; restore preserves `2.0` and the unrelated entry and becomes a no-op. Existing version 1 restore, retired deletion and reconciliation regressions also pass.
3. `go test ./internal/baseline -run '^Test(OwnedSkillLockEntriesUseTheHashOfTheLockVersion|ARetiredSkillIsProvenUneditedByTheHashOfTheLockVersion|SkillsRestoreWritesTheHashOfTheLockVersion|EveryLockReaderRefusesAnUnsupportedVersionByName)$' -count=1 -timeout=5m`
   - Rechecked the four new Baseline tests after adding explicit clone digest/version assertions for both versions.
4. `go test ./internal/cli -run '^TestDoctorKeepsReadinessWhenOwnedLockComparisonIsUnavailable$' -count=1 -timeout=5m`
   - The moved version 3 schema fixture preserves Doctor's existing readiness behavior.
5. `git diff --check` — no whitespace errors. Scope inspection found only this Task's declared paths and its Result; the pre-existing Daemon change to `status: in_progress` was preserved. `internal/skillhash/hash.go` and `hash_test.go` have no diff.

Pre-change evidence: the focused version 2 edge-folder test initially failed to compile because the new hashing/version APIs were absent. Added numeric-boundary tests then exposed rounding and quoted-number acceptance; those tests pass after the parser and decoding changes.

The declared Verification commands, repository Verification and incremental Verification were not run; they remain Daemon-owned. No Task status, Task Graph or other Task file was edited, and no commit, push or pull request was made. Doctor/Baseline next-action changes and guide updates remain in their planned Tasks 02 and 03.

## Carry-forward provenance

- Source Run: `run_20261112T212907Z_a8dd596719159f5e`
- Source commit: `976a1939b385cf1eb073b7cbfd30e4faea6ea5ca`
