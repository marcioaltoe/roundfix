---
spec: 0349-the-baseline-update-leaves-owned-skill-lock-entries-unpinned
prd: _prd.md
created: 2026-11-04
---

# The baseline update leaves owned skill lock entries unpinned — Technical Spec

## Executive Summary

`RecordInstalledOwnedSkills` in `internal/baseline/owned_skills_lock.go` is
the one writer of owned `skills-lock.json` entries. It already edits the
document by byte offsets so that only the selected value changes. The fix
adds one edit per recorded entry: it removes the `ref` member together with
one separator. The function returns which entries got a new hash and which
lost their `ref`, so `roundfix skills install` and `roundfix baseline update`
can name both.

The same file gains `StaleOwnedLockEntries`, a read-only scan of the same
entries. `repositorySkillsCheck` in `internal/cli/doctor.go` calls it and
appends `DR-SKILL-LOCK-STALE` to the `skills:` line as a warning. Removing
`ref` rather than writing one follows the measurement in `_prd.md` →
Acceptance evidence and ADR-0359: the skills CLI treats an entry without
`ref` as a default-branch install, and the binary knows no catalog commit.

The trade-off accepted: an unpinned entry follows the catalog's default
branch. In the catalog's daily sync window right after a release, that branch
can still hold the previous copy. Doctor's version check reports that case as
it does today.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential; `baseline update` and Doctor stay offline. Source:
  `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0189's addendum keeps every field
  but `computedHash`; ADR-0359 narrows that rule to remove `ref` and adds the
  Doctor warning. ADR-0233 and ADR-0256 govern the Roundfix skill's version
  raise and its coverage in task_03; ADR-0349 is the `minor` Release Ceiling
  the authorization records; ADR-0250, ADR-0197, ADR-0272, ADR-0279 and
  ADR-0286 are not applicable, because no Baseline module and no review
  behavior changes. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle AX answers of 2026-11-04 quoted in the authorization record; the
  Governed Paths are the canonical Roundfix `baseline` and `setup` references
  and the Roundfix `SKILL.md` pair; no Makefile, workflow, module or lint
  configuration changes. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0349-the-baseline-update-leaves-owned-skill-lock-entries-unpinned/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/baseline.md`, `.agents/skills/roundfix/references/setup.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Read on 2026-11-04 at `6a0ed904`; `_prd.md` → Validation records the run.

1. `RecordInstalledOwnedSkills(ctx, root, names) ([]string, error)` loads the
   lock with `loadSkillsLock`. For each sorted, deduplicated name that has an
   object entry and an installed `.agents/skills/<name>` tree, it replaces
   the `computedHash` value with `externalSkillsLockDigest` of the tree, or
   inserts the member when it is absent. It writes through a temporary file
   and rename, keeps the file mode, and returns the names whose hash changed.
   `ref` and every other byte stay.
2. `runBaselineUpdateSkillsStageWith` (`internal/cli/baseline_update.go`)
   stores that list in `skills.lockUpdated` (`omitempty`), and
   `writeBaselineUpdateResult` prints `Skills lock updated: <n>` and one
   `- lock updated <name>` line per name. The stage runs only after an
   approved apply; a `current` repository never reaches it.
3. `runSkillsInstallCommand` (`internal/cli/cli.go`, the project target
   without `--dir`) prints `Updated skills-lock.json computedHash: <name>`
   for each returned name.
4. `repositorySkillsCheck` (`internal/cli/doctor.go`) reports owned skills by
   their version only (`doctorSkillReadinessResult`). With a Setup Manifest it
   appends the Setup Snapshot comparison suffix
   (`DR-SKILL-DIFFERS-FROM-SNAPSHOT`) and keeps the status. A `warn` status
   never fails Doctor (`runDoctorCommand` fails only on `failed`).
   `collectReleasePlanChecks` reuses the same result for the release plan's
   `skills` check.
5. Tests pin today's bytes. `TestRecordInstalledOwnedSkillsRewritesOnlyTheHashOfAnInstalledOwnedSkill`
   expects `"ref":"keep"` to survive, and `TestSkillsInstallRewritesTheOwnedSkillLockHash`
   expects `"ref":"upstream"` to survive. Two stubs return `([]string, error)`
   for `recordLock`: `baseline_update_test.go` (two) and
   `doctor_trailing_skills_test.go` (one).

## System Architecture

No package, command or flag is added. Changes by component:

- `internal/baseline/owned_skills_lock.go`: member offsets, the `ref` removal,
  the result type, `StaleOwnedLockEntries`.
- `internal/cli/baseline_update.go`, `internal/cli/cli.go`: the reports.
- `internal/cli/doctor.go`: the warning.
- Tests in `internal/baseline` and `internal/cli`; the doctor, baseline and
  skills command guides; the Roundfix skill; `CONTEXT.md`.

## Implementation Design

### Interfaces

```go
// internal/baseline/owned_skills_lock.go
type OwnedSkillsLockUpdate struct {
	Hashed   []string // computedHash rewritten or added, sorted
	Unpinned []string // ref member removed, sorted
}
func (update OwnedSkillsLockUpdate) Names() []string // sorted union
func RecordInstalledOwnedSkills(ctx context.Context, root string, names []string) (OwnedSkillsLockUpdate, error)

type StaleOwnedLockEntry struct {
	Skill string
	Ref   string // the ref member's value as text; "" when absent
	Hash  bool   // computedHash absent or not the installed tree's
}
func StaleOwnedLockEntries(ctx context.Context, root string, names []string) ([]StaleOwnedLockEntry, error)
```

### Data Models

`skills-lock.json` keeps its schema. In an owned entry the bundle install
rewrites, the `ref` member is absent; the skills CLI already reads `ref` as
optional. `baseline update` JSON gains
`skills.lockUnpinned: [<name>...]`, with `omitempty`, after `lockUpdated`.

### API Contracts

1. API Contract: the writer. `RecordInstalledOwnedSkills` keeps its entry
   selection (a sorted, deduplicated name with an object entry and an
   installed tree) and its hash rule. In each selected entry it also removes
   the `ref` member, whatever its value, with this span:
   - the previous member's value end to the member's value end, when `ref` is
     not first;
   - otherwise the `ref` key's opening quote to the next key's opening quote;
   - or the key's opening quote to the value end, when `ref` is alone.
   An entry that loses `ref` is listed in `Unpinned`. An entry whose hash is
   written or inserted is listed in `Hashed`. The inserted `computedHash`
   member takes a leading comma only when another member remains. Every byte
   outside these spans and the hash value is unchanged. With no edit the file
   is not written.
2. API Contract: the reports. `runSkillsInstallCommand` prints each `Hashed`
   name as `Updated skills-lock.json computedHash: <name>`, then each
   `Unpinned` name as `Removed skills-lock.json ref: <name>`.
   `runBaselineUpdateSkillsStageWith` sets `LockUpdated` to `Names()` and
   `LockUnpinned` to `Unpinned`. Text prints, after the `- lock updated`
   lines, one `- lock ref removed <name>` line per `LockUnpinned` name.
   `recordLock` in `baselineUpdateSkillsDependencies` takes the new return
   type.
3. API Contract: the scan. `StaleOwnedLockEntries` never writes. An absent
   lock returns nothing. An unsafe or malformed lock returns the
   `loadSkillsLock` error. For each sorted, deduplicated name with an object
   entry and an installed tree, it reports the entry when the entry has a
   `ref` member or its `computedHash` is not the tree's
   `externalSkillsLockDigest`. Entries for names not in `names` are never
   read.
4. API Contract: the warning. With a Setup Manifest,
   `repositorySkillsCheck` calls `StaleOwnedLockEntries(ctx, root,
   skills.Names())` after the snapshot suffix. For one or more entries it
   appends
   `; DR-SKILL-LOCK-STALE: skills update would replace these Roundfix-owned skills with the copy their skills-lock.json entry pins: <sorted names, ", ">`.
   It turns an `ok` status into `warn`, keeps any other status, and appends
   `roundfix skills install --target project` (`ownedSkillsNextAction`) to the
   next action unless it is already there, joined by ` && `. On a scan error
   it appends `; skills-lock.json comparison unavailable: <error>` and keeps
   status and next action.
5. API Contract: the rules in writing. These places state the rule, each in
   its own words:
   - `docs/user-guide/commands/baseline.md` and
     `docs/user-guide/commands/skills.md`: the lock sentence names the removed
     `ref` and the new lines and field.
   - `docs/user-guide/commands/doctor.md`: the `skills:` bullet and a
     paragraph on `DR-SKILL-LOCK-STALE`, including that the release plan's
     `skills` check reads the same line.
   - The canonical Roundfix references `baseline.md` and `setup.md`, then
     their mirrors through `make skills-sync`, with the version raised by the
     record command.
   - `CONTEXT.md`: **Stale Owned Lock Entry**.

### Surface Transcripts

1. Surface Transcript: `roundfix skills install` in a Git repository whose
   lock holds `archive-spec`, `roundfix` and `write-prd` in the skills CLI
   format, pinned to `8d71fade68c54693bedbd917371f0514311af157` with a zero
   hash, beside the external `testing-boss`; `<dir>` is the repository's
   `.agents/skills` and `<n>` the bundle's file count.

   ```transcript
   $ roundfix skills install
   stdout:
   Installed Roundfix skill for project: <dir> (<n> file(s))
   Updated skills-lock.json computedHash: archive-spec
   Updated skills-lock.json computedHash: roundfix
   Updated skills-lock.json computedHash: write-prd
   Removed skills-lock.json ref: archive-spec
   Removed skills-lock.json ref: roundfix
   Removed skills-lock.json ref: write-prd
   stderr:
   exit: 0
   ```

Doctor prints one line per check, and several checks depend on the host
(adapters, profiles, the forge). So no whole-command Doctor transcript is
declared; `TestDoctorNamesAStaleOwnedSkillLockEntry` asserts the whole
printed `skills:` line.

### Invariants

```text
1. Every entry the writer does not select (external, unlisted, or an owned skill not installed) keeps its bytes, its ref included.
2. A selected entry keeps every member except ref, in order, with its bytes; only computedHash's value changes or the member is inserted.
3. The rewritten document is valid JSON that loadSkillsLock accepts.
4. A second bundle install over the result writes nothing and reports nothing.
5. Doctor never names an external entry, never changes its exit code for DR-SKILL-LOCK-STALE, and reads ok after one bundle install over a stale lock.
```

### Declared breaks

- An owned entry's `ref` is removed by the first bundle install after the
  upgrade (ADR-0359); adopters see the change in their lock diff. ADR-0189's
  addendum kept it.
- `TestRecordInstalledOwnedSkillsRewritesOnlyTheHashOfAnInstalledOwnedSkill`
  and `TestSkillsInstallRewritesTheOwnedSkillLockHash` keep their names and
  meaning (only the hash changes when no `ref` exists): their fixtures replace
  the `ref` member with `sourceType`. The new tests own the `ref` behavior.
- Doctor's `skills:` line, and so the release plan's `skills` check, can read
  `warn` for a repository whose lock pins owned entries. This repository's
  root lock holds no owned entry, so its own line stays `ok`.
- `RecordInstalledOwnedSkills` changes its return type. Its only callers are
  the two commands and the three test stubs named in Current behavior 5.

## Coverage Map

- Goal 1, Core Feature 1, Success Metrics 1 and 2 → API Contract 1,
  Invariants 1–4, `TestRecordInstalledOwnedSkillsRemovesTheRefOfAnInstalledOwnedSkill`.
- Goal 2, Core Feature 2, Success Metric 3 → API Contract 2, Surface
  Transcript 1, `TestSkillsInstallRemovesTheRefOfOwnedSkillLockEntries`,
  `TestBaselineUpdateSkillsStageReportsTheUnpinnedLockEntries`.
- Goal 3, Core Feature 3, Success Metric 4 → API Contracts 3 and 4,
  Invariant 5, `TestDoctorNamesAStaleOwnedSkillLockEntry`,
  `TestDoctorKeepsAFailedSkillsLineWithAStaleOwnedSkillLockEntry`.
- Goal 4, Core Feature 4, Success Metric 5 → API Contract 5,
  `TestTheGuidesStateTheUnpinnedOwnedLockEntry`.

Deviations from the PRD: None.

## Absence Evidence

1. `OwnedSkillsLockUpdate`, `StaleOwnedLockEntries` and `StaleOwnedLockEntry` — command: `grep -rn 'OwnedSkillsLockUpdate\|StaleOwnedLock' internal skills`; result: no match.
2. A `ref` removal in the owned writer — command: `grep -n '"ref"' internal/baseline/owned_skills_lock.go`; result: no match.
3. A Doctor reading of owned lock entries — command: `grep -n 'skills-lock\|LOCK' internal/cli/doctor.go`; result: no match.
4. The new report lines and JSON field — command: `grep -rn 'Removed skills-lock.json ref\|lock ref removed\|lockUnpinned' internal docs .agents`; result: no match.
5. `internal/cli/doctor_skills_lock_test.go`, `internal/baseline/owned_skills_lock_unpin_test.go` and `skills/owned_lock_guidance_test.go` — command: `ls internal/cli/doctor_skills_lock_test.go internal/baseline/owned_skills_lock_unpin_test.go skills/owned_lock_guidance_test.go`; result: no such files.

## ADR Alternatives

- ADR-0189: its addendum rejected recording Roundfix as the source because the
  lock format names a provider that can fetch the skill again. The design keeps
  that: `source`, `sourceType` and `skillPath` stay. It also keeps the hash as
  the folder's content and drops only the pin, which ADR-0359 amends.
- ADR-0233: kept as recorded; it rejected a hand-chosen skill version that
  collides at merge, and task_03 raises the Roundfix skill's version only with
  the record command.
- ADR-0256: kept as recorded; it rejected releasing a surface whose covering
  skill was not reviewed, and task_03 edits the covering `baseline` and
  `setup` references in the same Spec.
- ADR-0250: not applicable; it rejected choosing a module version at
  authoring time, and no Baseline module changes here.
- ADR-0349: kept as recorded; it rejected a second approval for a plan at or
  below the recorded ceiling, and the authorization records `minor`.
- ADR-0197: not applicable; its rejected unbounded reviewer lineage concerns
  the pre-PR review, which this Spec does not change.
- ADR-0272: not applicable; its rejected review-record overwrite concerns the
  pre-PR review record, which this Spec does not change.
- ADR-0279: not applicable; its rejected dismissal handling concerns queue
  status and review dismissals, which this Spec does not change.
- ADR-0286: not applicable; its rejected repeated Correction Reviews concern
  the pre-PR review, which this Spec does not change.
- ADR-0359: chosen here. It rejects writing a catalog commit, changing the
  provider or dropping `skillPath`, and failing Doctor. The design takes none
  of those.

## Integration Points

None at run time. The skills CLI behavior the design relies on is the quoted
1.7.6 code in `_prd.md` → Acceptance evidence; no test runs it.

## Testing Approach

All tests are hermetic: temporary directories, the embedded bundle, no
network, no skills CLI. Every new top-level test calls `t.Parallel()` first.

- `internal/baseline/owned_skills_lock_unpin_test.go` (new),
  `TestRecordInstalledOwnedSkillsRemovesTheRefOfAnInstalledOwnedSkill`. It is
  table-driven over the existing `ownedLockFixture` (installed `roundfix`
  only), names `roundfix` and `write-prd`, and checks whole-file bytes and
  `Hashed`/`Unpinned`:
  - `skills tool format`: two-space JSON with `roundfix` pinned between
    `source` and `sourceType`, an external `testing-boss` pinned, and
    `write-prd` pinned but not installed; only the `roundfix` ref line goes
    and only its hash changes.
  - `ref first`: `{ "ref" : "x" , "computedHash":"<stale>"}` becomes
    `{ "computedHash":"<hash>"}`.
  - `ref last with a current hash`: `ref` set to `""`, only `Unpinned`.
  - `ref alone`: `{"ref":null}` becomes `{"computedHash":"<hash>"}`.
- `internal/baseline/owned_skills_lock_test.go`: the declared fixture change;
  results read through `Names()`; a second run reports nothing.
- `internal/cli/owned_skills_lock_test.go`:
  - `TestSkillsInstallRemovesTheRefOfOwnedSkillLockEntries`: Surface
    Transcript 1 through `runCLI` in `withCLIWorkspace`. It asserts whole
    stdout, empty stderr, exit 0 and whole-file lock bytes, with the external
    entry byte-identical.
  - `TestBaselineUpdateSkillsStageReportsTheUnpinnedLockEntries`: the stage
    with `ownedLockStageDependencies`. It asserts lock bytes, an external
    `"ref":"keep"` kept, the text
    `Skills lock updated: 1\n- lock updated roundfix\n- lock ref removed roundfix\nSkills restored: `,
    and the JSON `"lockUpdated":["roundfix"],"lockUnpinned":["roundfix"]`.
  - The declared fixture change of `TestSkillsInstallRewritesTheOwnedSkillLockHash`.
- `internal/cli/doctor_skills_lock_test.go` (new),
  `TestDoctorNamesAStaleOwnedSkillLockEntry`. It runs over a root with the
  installed bundle and an installed external `testing-boss` folder, with
  `doctorDependencies{resolveExternal: nil skills, ok; checkSkills:
  skills.CheckRepositoryWithExternal}`. It prints `repositorySkillsCheck`
  through `printDoctorResult` and compares the whole line:
  - `pinned after a bundle install`: `archive-spec` and `roundfix` carry the
    installed hash plus a `ref`, `write-prd` has a zero hash and no `ref`, and
    `testing-boss` is pinned. The line is
    `skills: warn (<counts>; DR-SKILL-LOCK-STALE: skills update would replace these Roundfix-owned skills with the copy their skills-lock.json entry pins: archive-spec, roundfix, write-prd; next: roundfix skills install --target project)`.
  - `unpinned with the installed hash`: `skills: ok (<counts>)`.
  - After `RecordInstalledOwnedSkills` over each case, the line is
    `skills: ok (<counts>)`.
- `internal/cli/doctor_skills_lock_test.go`,
  `TestDoctorKeepsAFailedSkillsLineWithAStaleOwnedSkillLockEntry`: the
  installed `roundfix` copy's `version:` lines lowered to `0.0.1` and its entry
  pinned with a zero hash. The line starts `skills: failed (`, carries the
  `DR-SKILL-LOCK-STALE` suffix naming `roundfix` before `; next: `, and names
  `roundfix skills install --target project` exactly once.
- `skills/owned_lock_guidance_test.go` (new),
  `TestTheGuidesStateTheUnpinnedOwnedLockEntry`. On whitespace-normalized
  text, both copies of the Roundfix `baseline` reference contain
  `- lock ref removed <name>` and `skills.lockUnpinned`, and both copies of
  the `setup` reference contain `Removed skills-lock.json ref: <name>` and
  `DR-SKILL-LOCK-STALE`.
- Existing suites: `./internal/baseline`, `./internal/cli`, `./skills`,
  `-tags docscontract ./internal/docscontract` (Behavior Surface Record,
  glossary contract, command guides).
- Sabotages, each run by the author on a prototype copy at `6a0ed904`, with
  the case whose input makes it observable:
  - `ref` removal dropped: the unpin test fails on all four rows;
  - first-member span taken as key to own value end: it fails on `ref first`;
  - `ref` removed before the installed check: it fails on `skills tool format`
    (`write-prd` loses its `ref`);
  - the scan ignores the hash: the Doctor test fails on `pinned after a bundle
    install` (`write-prd` missing);
  - the scan reads every lock entry instead of `names`: the Doctor test fails
    on both rows (`testing-boss` named);
  - the status stays `ok`: the Doctor test fails on `pinned after a bundle
    install`;
  - the next action appended without the "already there" check: the failed
    line test fails (the install named twice).

## Build Order

1. The writer, its result type and both reports, their tests and the declared
   fixture changes, the baseline and skills command guides, with the Behavior
   Surface Record re-recorded last (task_01).
2. The scan and the Doctor warning with its tests, the doctor command guide,
   with the Behavior Surface Record re-recorded last (task_02) (depends on:
   1).
3. The Roundfix skill references and their test, `CONTEXT.md`, the mirrors
   synced, the owned-skill version recorded and the Behavior Surface Record
   re-recorded last (task_03) (depends on: 2).
4. QA gate (task_04) (depends on: 3).

Each code Task updates the command guide of the surface it changes, so the
guide lands with the behavior. task_02 follows task_01 because both change
`internal/baseline/owned_skills_lock.go` and
`docs/references/behavior-surfaces.json`. task_03 follows task_02 because the
skill references describe both behaviors and it re-records the same record.
No Verification depends on another Spec. Spec 0347, which delivers first,
also changes `internal/cli/cli.go` (`loadRunReaderCommandConfig` and the
`runs` help); task_01 edits only `runSkillsInstallCommand`, so the two
rebase cleanly.

## Risks & Considerations

- During the catalog's daily sync window after a Roundfix release, `skills
  update` over an unpinned entry can install the catalog's previous copy.
  Doctor's version check fails then, as before this Spec, and one bundle
  install restores the copy. The guides say so.
- Doctor's existing external next action runs `bunx skills update -p -y`. Over
  a stale owned entry it would downgrade; the new warning names the entries
  first, and the bundle install clears them.
- An adopter who pins an owned skill on purpose loses that pin at the next
  bundle install. The bundle already replaces that copy, so the pin never
  held.

## Glossary

- adds: **Stale Owned Lock Entry** — an owned skill's lock entry that pins another copy

## Decisions

- Remove `ref` rather than write one (ADR-0359), after measuring
  `skills@1.7.6`.
- One writer and one scan in `internal/baseline/owned_skills_lock.go`; Doctor
  calls the scan directly, with no new dependency field.
- The warning uses the existing `warn` status and the existing owned next
  action; `skills check` is unchanged.
- Secondbrain: `wiki/index.md` read; the `qmd query` named in `_prd.md` →
  Validation found no deciding page. Exa: the lock write-up and Skillsaw rule
  quoted there.
- Three implementation Tasks in series and a QA gate.
