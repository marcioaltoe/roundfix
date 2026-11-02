---
spec: 0339-an-approved-baseline-update-deletes-the-retired-skills
prd: _prd.md
created: 2026-11-01
---

# An approved Baseline update deletes the Retired Skills — Technical Spec

## Executive Summary

`baseline.InstalledRetiredSkills` gains a disposition per held Retired Skill:
`delete` when every folder it finds hashes to the `computedHash` of the
skill's `skills-lock.json` entry and every `.claude/skills` link points at
`../../.agents/skills/<skill>`, otherwise `keep` with a reason. A new
`baseline.DeleteRetiredSkills` re-inspects and deletes exactly the planned
`delete` entries. `runBaselineUpdateCommandWithSkillsStage` inspects the
Retired Skills before it fixes the Plan Digest, folds the deletions into the
digest, reports `plan_ready` when a deletion is the only pending change, and
deletes after the skills stage of an approved apply. ADR-0350 supersedes
ADR-0246; the glossary, the guide and the Roundfix baseline reference follow.

The trade-off accepted: a retired copy with no recorded hash is kept even
when it is unedited, because nothing proves it; the two adopters surveyed
record a hash for every retired copy (`_prd.md` → Validation).

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; `roundfix/baseline-update-result/v1` gains optional fields.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential is added. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0350 governs the design and
  supersedes ADR-0246, which task_02 deletes; ADR-0191 is amended for
  Retired Skills only; ADR-0206, which ADR-0246 cites, is unchanged
  (`_prd.md` → Project Constraints; ADR Alternatives below). Source:
  `docs/agents/domain.md`, `docs/agents/docs-layout.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AU answers of 2026-11-01 quoted in the authorization record; the Governed
  Paths are the Roundfix skill's canonical `baseline` reference and its
  `SKILL.md` pair; no Makefile, workflow, module or lint configuration
  changes. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0339-an-approved-baseline-update-deletes-the-retired-skills/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/baseline.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Read on 2026-11-01 at `dd3b2b62`; `_prd.md` → Validation records the runs.

1. `InstalledRetiredSkills(repoRoot)` loads `skills-lock.json`, and for each
   name of `RetiredSkills()` records whether the lock has an entry and which
   of `.agents/skills/<skill>` and `.claude/skills/<skill>` exist (`Lstat`).
   It reads no content and returns `[]InstalledRetiredSkill{Skill, Paths,
   LockEntry}`.
2. `runBaselineUpdateCommandWithSkillsStage` (`internal/cli/baseline_update.go`)
   sets `result.PlanDigest` from `historyPlan.digest(plan.PlanDigest)`. In the
   preview branch, unless `--no-skills`, it fills `result.Skills.Retired`;
   with no file change and no history move the preview reports `current`
   and exit 0 whatever the Retired Skills.
3. On approval it compares the confirmation with `result.PlanDigest`, calls
   `baseline.ApplyPlan`, applies history, runs the skills stage, then fills
   `result.Skills.Retired` again. Nothing is deleted.
4. `writeBaselineUpdateResult` prints `Skills retired: <n>` and, per skill,
   `- retired <skill>: no longer required by the Baseline; delete <targets>`,
   where the targets join the paths and "its skills-lock.json entry".
5. `skills_restore_lock.go` already holds `loadSkillsLock`,
   `removeSkillsEntry` and `marshalIndent`, used by `baseline skills
   reconcile`; `inspectRestoreTarget` reads a skill folder, refusing links and
   hard-linked files, and `externalSkillsLockDigest` computes the lock's
   `computedHash` over its files.
6. Every Setup Snapshot lists `context7-cli` and `exa-web-search`, and
   `RetiredSkills()` names neither; no test pins both facts together.

## System Architecture

No package, command or flag is added. Changes by component:

- `internal/baseline/retired_skills_installed.go`: dispositions and deletion.
- `internal/cli/baseline_update.go`: inspection before the digest, the digest
  fold, the preview state, the deletion step and the text output.
- Tests in `internal/baseline` and `internal/cli`; ADR-0191, ADR-0246,
  `CONTEXT.md`, the baseline guide and the Roundfix baseline reference.

## Implementation Design

### Interfaces

```go
// internal/baseline/retired_skills_installed.go

const (
	RetiredSkillDelete = "delete"
	RetiredSkillKeep   = "keep"

	RetiredSkillEdited     = "edited"
	RetiredSkillUnrecorded = "unrecorded"
	RetiredSkillUnsafe     = "unsafe"
)

type InstalledRetiredSkill struct {
	Skill     string   `json:"skill"`
	Paths     []string `json:"paths"`
	LockEntry bool     `json:"lockEntry"`
	Action    string   `json:"action"`
	Reason    string   `json:"reason,omitempty"`
	LockHash  string   `json:"lockHash,omitempty"`
}

// InstalledRetiredSkills lists each held Retired Skill with its disposition.
// It only reads.
func InstalledRetiredSkills(repoRoot string) ([]InstalledRetiredSkill, error)

// DeleteRetiredSkills deletes exactly the planned delete entries and returns
// them; it deletes nothing when one changed since the plan.
func DeleteRetiredSkills(ctx context.Context, repoRoot string, planned []InstalledRetiredSkill) ([]InstalledRetiredSkill, error)

// internal/cli/baseline_update.go
type baselineUpdateSkillsResult struct {
	// ...existing fields unchanged...
	Deleted []baseline.InstalledRetiredSkill `json:"deleted,omitempty"`
}

func retiredSkillsPlanDigest(digest string, retired []baseline.InstalledRetiredSkill) string
```

### Data Models

`roundfix/baseline-update-result/v1` gains, additively:
`skills.retired[].action` (`delete` or `keep`), `skills.retired[].reason`
(`edited`, `unrecorded` or `unsafe`, only on `keep`),
`skills.retired[].lockHash` (the lock entry's `computedHash`, when recorded),
and `skills.deleted[]` with the same shape, present only after an approved
apply that deleted something. No schema version changes.

### API Contracts

1. API Contract: the disposition. For each name of `RetiredSkills()` with a
   lock entry or an entry under `.agents/skills` or `.claude/skills`:
   - a folder entry is a tree; a `.claude/skills/<skill>` symbolic link whose
     target is exactly `../../.agents/skills/<skill>` is a link; any other
     entry (a file, an `.agents/skills` link, a link elsewhere) is unsafe;
   - `keep`/`unsafe` when an entry is unsafe or `inspectRestoreTarget`
     refuses a tree; else `keep`/`unrecorded` when a tree exists and the lock
     records no `computedHash`; else `keep`/`edited` when a tree's
     `externalSkillsLockDigest` differs from it; else `delete`. A lock entry
     or link with no tree is `delete`.
   Output order stays lexical by skill and, per skill, `.agents` before
   `.claude`.
2. API Contract: the deletion. `DeleteRetiredSkills` re-runs the
   inspection, refuses an entry whose skill is not in the Retired Skills set
   ("refuse to delete %q: it is not a Retired Skill") and an entry whose
   current disposition differs in any field ("retired skill %q changed after
   the plan; rerun roundfix baseline update"), both before deleting anything.
   It then removes, per skill, its paths in reverse order (`.claude` first)
   with `os.RemoveAll`, which never follows a link, and finally writes
   `skills-lock.json` once without the deleted entries through
   `loadSkillsLock`, `removeSkillsEntry` and `marshalIndent`, keeping the
   file mode. Keep entries are never touched.
3. API Contract: the Plan Digest. After the history digest, unless
   `--no-skills`, the command fills `result.Skills.Retired` (inspection
   failure keeps today's category `execution` and next action) and sets
   `result.PlanDigest = retiredSkillsPlanDigest(result.PlanDigest, retired)`.
   With no `delete` entry the digest is returned unchanged; otherwise it is
   `sha256:` and the hex SHA-256 of the JSON object `{"planDigest": <digest>,
   "retiredSkills": [<delete entries>]}`. The confirmation compares against
   this digest, so a change to a planned copy between preview and confirm is
   refused by the existing mismatch path.
4. API Contract: the preview state. With no file change and no history
   move, after the history case and before the missing-skill case, a
   preview with at least one `delete` entry reports state `plan_ready`,
   category `approval`, exit 3, message "guidance matches the current
   Baseline catalog; the approved update deletes <n> Retired Skill(s)" and
   the existing confirm next action. `keep` entries never change the state.
5. API Contract: the approved apply. After the skills stage, unless
   `--no-skills`, the command calls `DeleteRetiredSkills` with the entries
   the digest bound, stores the result in `result.Skills.Deleted`, then
   re-inspects into `result.Skills.Retired`. A deletion error reports
   category `skills`, next action "delete the remaining Retired Skills by
   hand or rerun roundfix baseline update" and exit 1.
6. API Contract: the text. Before `Skills retired`, when
   `Skills.Deleted` is non-empty: `Skills deleted: <n>` and per skill
   `- deleted <skill>: <targets>`. Under `Skills retired`, a `delete` entry
   prints `- retired <skill>: no longer required by the Baseline; the
   approved update deletes <targets>` and a `keep` entry prints `- retired
   <skill>: kept, <reason text>; delete <targets> by hand if the repository
   no longer uses it`, with reason texts: `edited` → "its copy differs from
   its skills-lock.json computedHash"; `unrecorded` → "skills-lock.json
   records no computedHash to prove its copy unedited"; `unsafe` → "an entry
   is not a skill folder or a link to .agents/skills". `<targets>` keeps
   today's joining.
7. API Contract: the rules in writing.
   - ADR-0350 (written by the planning change) states the rule; task_02
     deletes ADR-0246 and appends to ADR-0191 an amendment beginning
     "Amendment (2026-11-01): ADR-0350 lets an approved Baseline update
     delete a Retired Skill it proves unedited." with `updated_at`
     `2026-11-01T00:00:00Z`.
   - `CONTEXT.md`: **Retired Skill** links ADR-0350 instead of ADR-0246 and
     states that an approved `roundfix baseline update` deletes a copy it
     proves unedited against `skills-lock.json`, keeping an edited or
     unrecorded copy.
   - `docs/user-guide/commands/baseline.md` → `### Retired skills` and
     `.agents/skills/roundfix/references/baseline.md`: the update deletes
     the proven copies on approval, lists kept copies with their reason, the
     preview turns `plan_ready`, `Skills deleted` / `skills.deleted` report
     the deletions, and nothing outside the Retired Skills set is deleted;
     the sentence "An approved update deletes each Retired Skill it proves
     unedited" appears verbatim, and "never deletes" is gone from both
     passages. The guide keeps its manual commands for kept copies.

### Surface Transcripts

None. The changed lines are asserted by the named tests; no whole-command
transcript is declared.

### Invariants

```text
1. Only names of RetiredSkills() are ever deleted, and only their .agents/skills and .claude/skills entries and lock entries.
2. A copy that differs from its recorded hash, has no recorded hash, or is not a plain folder or the .agents link is never deleted.
3. Every deletion is listed in the preview and bound by the Plan Digest; a preview with no deletion keeps its digest bytes.
4. A confirmation whose planned copies changed deletes nothing.
5. --no-skills inspects and deletes nothing.
```

### Declared breaks

- An applied update deletes unedited retired copies. `TestBaselineUpdateAppliesAndListsRetiredSkills`
  keeps its name and its assertion that the copies stay, because its fixture
  (`adoptedRepositoryWithRetiredSkills`) records no hash and both copies are
  now kept as `unrecorded`; its transcript `retiredSkillsTranscript` and
  `expectedInstalledRetiredSkills` change to the keep lines and fields.
- `TestBaselineUpdateRetiredSkillDeletionTargets` asserts the `delete` line,
  whose wording changes; `TestInstalledRetiredSkillsListsTreesLinksAndLockEntries`
  gains the action and reason of each case; the combined preview subtest of
  `TestBaselineUpdatePreviewListsASetupSnapshotDifference`
  (`internal/cli/repository_profile_docscontract_test.go`) asserts the
  council keep line.
- A preview whose only pending change is a deletable Retired Skill reports
  `plan_ready` and exit 3 instead of `current` and 0.
- No top-level test is removed or renamed, so the coverage record is
  unchanged.

## Coverage Map

- Goal 1, Core Features 2 and 3, Success Metric 1 → `DeleteRetiredSkills` and
  the approved apply (API Contracts 2, 5 and 6), Invariants 1.
- Goal 2, Core Feature 3, Success Metrics 1 and 3 → the digest fold and the
  preview state (API Contracts 3 and 4), Invariants 3 and 4.
- Goal 3, Core Feature 1, Success Metrics 2 and 4 → the disposition (API
  Contract 1) and the keep lines (API Contract 6), Invariants 2.
- Goal 4, Success Metrics 2, 4 and 5 → the name guard (API Contract 2), the
  snapshot test, Invariants 1 and 5.
- Goal 5, Core Feature 4, Success Metric 6 → ADR-0350, ADR-0191, ADR-0246,
  `CONTEXT.md`, the guide and the reference (API Contract 7).

Deviations from the PRD: None.

## Absence Evidence

1. A deletion of Retired Skills — command: `grep -rn 'RemoveAll\|removeSkillsEntry' internal/baseline/retired_skills_installed.go internal/cli/baseline_update.go`; result: no match; nothing deletes a retired copy.
2. A disposition on the inspection — command: `grep -n 'Action\|computedHash' internal/baseline/retired_skills_installed.go`; result: no match; the struct holds skill, paths and lock entry only.
3. A Retired Skill term in the Plan Digest — command: `grep -n 'Retired' internal/cli/baseline_update_history.go`; result: no match; the digest covers the Baseline plan and the history units only.
4. Retired Skill wording in the setup skill — command: `grep -rn -i 'retired' .agents/skills/setup-context-driven`; result: no match; the skill needs no change.

## ADR Alternatives

- ADR-0350: followed; its rejected alternatives, deleting every listed copy
  edited or not, and refusing the whole update on an edited copy, stay
  rejected: an edited or unrecorded copy is kept with a reason and the rest
  of the update proceeds.
- ADR-0246: superseded and retired; its rejected alternative "deleting the
  adopter's installed copies, or editing its lock, during baseline update" is
  the decision ADR-0350 now takes for proven copies only, while its other
  rejected alternative, leaving retired names in the snapshots, stays
  rejected.
- ADR-0191: amended; its rule that Roundfix never deletes an installed skill
  tree holds for every skill outside the Retired Skills set, renamed skills
  included, and its rejected alternatives on snapshot membership are
  untouched.
- ADR-0204: unaffected; a composed profile still takes its setup composed
  from upstream setups by name, and its rejected alternatives on composition
  stay rejected, because this Spec changes no setup.
- ADR-0219: unaffected; a profile draft still binds to the closest built-in
  profile, and no rejected alternative of it is touched.
- ADR-0241: unaffected; a repository profile still takes its skill contracts
  from the embedded Setup Snapshots, which this Spec reads but never changes.
- ADR-0322: unaffected; the Python profile's import of the upstream preset
  and its rejected alternatives stay as they are.
- ADR-0206: unaffected; dropping a skill upstream removed and the name-based
  setup sync stay as they are, and this Spec adds no snapshot rule.

## Integration Points

None external. The update reads and writes local files and Git.

## Testing Approach

All tests are hermetic: temporary directories and Git repositories through
the existing `newBaselineApplyTestRepository`, `installRequiredExternalSkills`,
`runBaselineUpdateTestCommand` and `baselinePlanTestTree` helpers, no network.

- `internal/baseline/retired_skills_installed_test.go`:
  - `TestRetiredSkillDispositionsProveACopyUnedited`: one repository with
    `tdd` (unedited folder, link to `../../.agents/skills/tdd`, recorded
    hash) → `delete`; `design-patterns` (folder differs) → `keep`/`edited`;
    `council` (folder, no lock entry) → `keep`/`unrecorded`; `to-prompt`
    (link to `../../elsewhere/to-prompt`) → `keep`/`unsafe`;
    `fixing-metadata` (unedited `.claude/skills` copy only) → `delete`;
    `the-fool` (lock entry only) → `delete`; asserted with
    `reflect.DeepEqual` on the whole slice, `LockHash` included.
  - `TestDeleteRetiredSkillsDeletesOnlyThePlannedCopies`, subtests "deletes
    the unedited retired copy and keeps every other skill" (`tdd` folder,
    link and lock entry gone; `design-patterns`, `context7-cli`,
    `exa-web-search`, `local-only` and their links and lock entries stay),
    "refuses a copy edited after the plan" (nothing deleted, lock bytes
    unchanged), "refuses a name outside the Retired Skills" (`context7-cli`
    planned as `delete` is refused and stays).
  - `TestInstalledRetiredSkillsListsTreesLinksAndLockEntries`: the existing
    cases gain `keep`/`unsafe`, `keep`/`unrecorded` and `delete`.
- `internal/baseline/retired_skills_test.go`:
  `TestContext7AndExaAreRequiredByEverySnapshotAndNeverRetired` reads every
  `assets/setups/*.json` and asserts both names in each, and neither in
  `retiredSkills`.
- `internal/cli/baseline_update_retired_skills_test.go`:
  - `TestBaselineUpdateApprovedApplyDeletesTheRetiredSkills`, subtests
    `text` and `json`, over `deletableRetiredSkillsRepository` (adopted
    repository plus `tdd` unedited with link, `design-patterns` edited,
    `context7-cli` with link, `exa-web-search`, `local-only` with link, all
    with recorded hashes): the preview's state, category and message; the
    preview digest absent from the `--no-skills` output; the two preview
    lines; the approved apply's `Skills deleted: 1` line or JSON
    `skills.deleted`; `tdd` folder, link and lock entry gone; every survivor
    tree byte-identical by `baselinePlanTestTree` and its link and lock entry
    present.
  - `TestBaselineUpdateRefusesARetiredSkillEditedAfterThePlan`: preview,
    edit `tdd`, confirm with the old digest → exit 3, `action_required`,
    `approval`, no `skills.deleted`, tree unchanged.
  - The existing tests updated as Declared breaks state.
- Existing suites: `./internal/baseline`, `./internal/cli`, `-tags
  docscontract ./internal/cli` and `./internal/docscontract`, `./skills`.
- Sabotages, each run on the prototype:
  - `DeleteRetiredSkills` returns before deleting:
    `TestDeleteRetiredSkillsDeletesOnlyThePlannedCopies` (first subtest) and
    `TestBaselineUpdateApprovedApplyDeletesTheRetiredSkills` (both subtests)
    fail;
  - the hash comparison never reports a difference:
    `TestRetiredSkillDispositionsProveACopyUnedited` (the `design-patterns`
    case) and `TestBaselineUpdateApprovedApplyDeletesTheRetiredSkills` fail;
  - `retiredSkillsPlanDigest` returns its input unchanged:
    `TestBaselineUpdateApprovedApplyDeletesTheRetiredSkills` (the
    `--no-skills` digest) and `TestBaselineUpdateRefusesARetiredSkillEditedAfterThePlan`
    fail;
  - the Retired Skills name guard removed:
    `TestDeleteRetiredSkillsDeletesOnlyThePlannedCopies` fails on "refuses a
    name outside the Retired Skills";
  - any `.claude/skills` link target accepted:
    `TestRetiredSkillDispositionsProveACopyUnedited` (the `to-prompt` case)
    and `TestInstalledRetiredSkillsListsTreesLinksAndLockEntries` fail.

## Build Order

1. Dispositions, deletion and their tests in `internal/baseline`, and the
   update: inspection before the digest, the fold, the preview state, the
   deletion step, the text, the CLI tests and the baseline guide section,
   with the Behavior Surface Record re-recorded last (task_01). One Task, because the
   new fields change the CLI's JSON and the existing CLI tests in the same
   change.
2. ADR-0246 deleted, ADR-0191 amended, `CONTEXT.md`, the Roundfix
   baseline reference, the skill version and the Behavior Surface Record
   re-recorded last (task_02) (depends on: 1).
3. QA gate (task_03) (depends on: 2).

No Task needs another Spec. Spec 0338 (delivered before) and Spec 0340
(after) also raise the Roundfix skill version and change `CONTEXT.md` and
`docs/references/behavior-surfaces.json`; task_02 rewrites the version
fields, `skills/testdata/owned-skill-versions.json` and the Behavior Surface
Record only through their record commands, so a rebase regenerates them
instead of merging their bytes. task_01 and task_02 both re-record the Behavior Surface Record, so they run
in series. Neither other Spec changes `internal/baseline/retired_skills*.go`,
`internal/cli/baseline_update*.go` or `docs/user-guide/commands/baseline.md`.

## Risks & Considerations

- A lock entry whose `computedHash` was rewritten by hand to match an edited
  copy defeats the proof; the copy is then deleted, and Git keeps it.
- `marshalIndent` rewrites `skills-lock.json` in its canonical indentation,
  as `baseline skills reconcile` already does; key order is preserved.
- A deletion that fails midway leaves earlier deletions done and reports the
  error with the next action; a rerun plans what remains.
- `docs/references/2026-10-08-baseline-audit.md` names ADR-0246 in a dated
  audit; it is a record of that date and keeps its text.
- The history-sanitize ADR compares a Refused Unit with "a Retired Skill",
  which stays true for a kept copy; it is not amended.

## Glossary

- changes: **Retired Skill** — an approved update deletes a copy it proves unedited

## Decisions

- Dispositions with three keep reasons, proven by the lock's `computedHash`;
  a lock-only or link-only entry is deleted.
- The deletion set joins the Plan Digest only when non-empty.
- The deletion runs after the skills stage and re-inspects first.
- Secondbrain: `wiki/index.md` read; the `qmd query` named in `_prd.md` →
  Validation found no page deciding beyond ADR-0246. Exa: the FreeBSD `pkg`
  FAQ and the GSD installer migration guide, which delete only what matches
  the recorded install and never an unknown file.
- Two implementation Tasks in series and a QA gate.
