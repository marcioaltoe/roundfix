---
spec: 0195-owned-skills-and-a-release-step-that-follow-the-bundle
prd: _prd.md
created: 2026-09-30
---

# Owned skills and a release step that follow the bundle — Technical Spec

## Executive Summary

Four changes, each with its own check. The owned-skill minimum is built from
the embedded bundle instead of a literal list, and the `baseline update`
preview reads the installed owned skills. A record file and a test make a
version name one content. Two setup snapshots and one module give the Roundfix
skill membership and a dispatch trigger, and the asset sync keeps an owned
entry the upstream list omits. The release runbook gains a mandatory step, the
release clause gains one sentence, and a test refuses a step that names a
missing check. The trade-off accepted is that Doctor becomes stricter: a
repository fails `skills:` after a binary upgrade until it refreshes its owned
skills. That costs one command per upgrade and buys a comparison that means
something.

## Project Constraints

- Identifier strategy: applicable — one new trigger identifier,
  `trigger.autonomous-work.roundfix`, in the existing form; two optional
  result fields in `roundfix/baseline-update-result/v1`. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — embedded files, local files and
  Git only; no credential and no network call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0189 governs the minimum, the
  version record, the setup minimum, owned membership and the release step.
  ADR-0186 governs the sentence added to the release clause. ADR-0062 and
  ADR-0143 stay in force. ADR-0073, ADR-0081, ADR-0103 and ADR-0149 govern the
  regeneration and this repository's refresh. ADR-0080, ADR-0088, ADR-0091,
  ADR-0093, ADR-0094, ADR-0096, ADR-0104, ADR-0117, ADR-0155 and ADR-0156 bind
  the gate and the checker. These ADRs cite a listed ADR and do not apply,
  because this Spec touches none of their behaviors: ADR-0168 (the
  related-ADR gap) and ADR-0176 (citation checks read authored text) cite
  ADR-0093; ADR-0173 (citations a History Relocation breaks) cites ADR-0073.
  Reached only through those citations, and equally untouched: ADR-0097, ADR-0167, ADR-0177.
  All hold. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30, recorded in [_authorization.md](_authorization.md); bounded
  files: `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `docs/agents/specific-repository.md`,
  `internal/baseline/assets/setups/go-cli.json`,
  `internal/baseline/assets/setups/rust-cli.json`,
  `internal/baseline/assets/modules/autonomous-work.json`,
  `internal/baseline/assets/modules/core.json`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `docs/agents/skill-dispatch.md`, `docs/agents/agent-instructions.md`,
  `docs/agents/setup-context.json`.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## System Architecture

| Concern | Source of truth | Derived by |
| --- | --- | --- |
| Owned skill versions and content | `.agents/skills/<skill>/` | `make skills-sync` writes the mirror under `skills/` |
| Minimum version of an owned skill | the embedded `SKILL.md` | read at start-up in `skills/skills.go` |
| Shipped versions and their digests | `skills/testdata/owned-skill-versions.json` | the record test, with its own flag |
| Setup membership | `internal/baseline/assets/setups/*.json` | hand edit under the grant; `roundfix baseline assets sync` for upstream entries |
| Dispatch triggers and clause text | `internal/baseline/assets/modules/*.json` | hand edit under the grant |
| Formatter goldens, the digest pin, catalog snapshots, plan goldens, the asset-sync parity fixture and its manifest | the modules and setups | `make baseline-digests` |
| This repository's guides and Setup Manifest | the modules | `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text` |
| The release step | `docs/user-guide/release-runbook.md` | hand edit |

The regeneration outputs were measured on 2026-09-30 in a scratch clone of
`9e439dbb`:

- **Membership and trigger (task_03).** Adding the Roundfix entry to the two
  setups and the parity fixture, and the trigger to `autonomous-work.json`,
  then running `go test -buildvcs=false ./skills -run '^TestAuthorialSkillSync$' -update -count=1`
  and `make baseline-digests`, rewrote:
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/baseline/testdata/catalog.diagnostics.golden.json`,
  `internal/baseline/testdata/catalog.digest`,
  `internal/baseline/testdata/catalog.normalized.json`,
  `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`,
  `internal/baseline/testdata/parity-corpus/v1/manifest.json` and the four
  plan goldens under `internal/baseline/testdata/plan-characterization/`.
  The managed refresh then rewrote `docs/agents/skill-dispatch.md` and
  `docs/agents/setup-context.json`, and a second refresh reported
  `File changes: 0`. `./skills`, `./internal/baseline` and `./internal/cli`
  passed.
- **Without the fixture rows** the same edit fails
  `TestAssetsSyncCompatibilityMatchesMaintainedPythonContract`. The two rows
  are therefore added to the fixture by hand, and the regeneration recomputes
  its digests.
- **Release clause (task_04).** An edit to `core.json` rewrites
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`, the digest pin, the three catalog snapshots,
  the four plan goldens, `docs/agents/agent-instructions.md` and
  `docs/agents/setup-context.json`.
- **Minimum and preview (task_01).** Building the minimum from the bundle and
  reading installed owned skills in the preview left `./skills`,
  `./internal/baseline` and `./internal/cli` green. Counting a missing owned
  skill as outdated failed `TestBaselineUpdateFleetSweep`,
  `TestBaselineUpdateNoUnrecordedManagedRegionOmitsOutputs` and
  `TestBaselineUpdateAdoptsADeclaredIncrementalSuggestion`, so a missing
  skill is not counted.
- **A version raise against the literal minimum.** Raising the Roundfix
  skill's version with the literal list in place failed three tests in
  `./skills`. With the minimum built from the bundle, one subtest still
  fails: `TestOwnedSkillContractRejectsSetAndVersionDisagreement/owned_version_below_minimum`
  reads the literal constant. task_01 makes it read the same map the code
  reads.

A module or setup edit must keep the file's existing formatting: replace
strings and numbers in place, and do not pass the file through a JSON encoder.
The parity fixture is the exception; it is stored in the encoder's own format
with sorted keys.

## Implementation Design

### Interfaces

No exported Go signature is removed or changed.

```go
// skills/skills.go
// ownedSkillMinimumVersions is built once from the embedded bundle: each
// owned skill maps to the version its embedded SKILL.md declares. A skill
// whose embedded file declares no valid version keeps the base constant.
var ownedSkillMinimumVersions = embeddedOwnedSkillVersions()

// internal/cli/baseline_update.go
type baselineUpdateSkillOutdated struct {
	Skill    string `json:"skill"`
	Found    string `json:"found"`
	Required string `json:"required"`
}
// baselineUpdateSkillsResult gains:
//   Outdated []baselineUpdateSkillOutdated `json:"outdated,omitempty"`
// and the status value "outdated".

// internal/baseline/assets_sync.go
// buildAssetsSyncSnapshot appends, after the upstream entries and in the
// snapshot's recorded order, each entry of the current snapshot whose
// source type is "repo" and whose name the upstream list did not yield.
```

### Data Models

`skills/testdata/owned-skill-versions.json`:

```json
{
  "schemaVersion": "roundfix/owned-skill-versions/v1",
  "skills": {
    "qa-gate": [
      { "version": "0.0.3", "digest": "<SkillFolderHash of the embedded folder>" }
    ]
  }
}
```

Each skill's list is ordered by version, lowest first. The digest is the
package's existing skill folder hash, computed over the embedded folder. The
file is not embedded in the binary and is not under a regenerated path.

### The record rule

For each owned skill, with `v` the version its embedded `SKILL.md` declares
and `d` the digest of its embedded folder:

1. `v` recorded with digest `d`: pass.
2. `v` recorded with another digest: fail, "content changed under version
   `v`; raise the version".
3. `v` not recorded: fail, and name the recording command.
4. A skill in the record that the bundle no longer ships, or a list that is
   not in ascending version order: fail.

Recording runs the same test with `-record-skill-versions`. It adds `v` with
`d` only in case 3 and only when `v` is higher than every recorded version of
that skill. In case 2 it fails exactly as the check does. The flag is separate
from the package's `-update` flag, and the test's name does not start with
`TestAuthorialSkillSync`, so `make baseline-digests` never runs it.

### The preview rule

In `roundfix baseline update`, when neither `--yes` nor `--confirm-plan` is
given and `--no-skills` is absent, after the Plan is built:

1. Read the repository's owned skills with the package's existing repository
   check, asking for no external skill.
2. Collect each owned skill whose state is below the minimum, sorted by name,
   as `skill`, `found`, `required`.
3. When the list is empty, or the read returns an error, the result is what
   it is today.
4. Otherwise set `skills.status` to `outdated` and `skills.outdated` to the
   list. When the Plan has no file change and no history move, the state is
   `plan_ready`, the category is `approval`, the exit code is `3`, and:
   - message: `guidance matches the current Baseline catalog; <n> Roundfix-owned skill(s) are older than the ones this binary carries`;
   - next action: `rerun with --yes to refresh the Repository Skill Set, or run roundfix skills install --target project`.
   When the Plan has changes, its state, message and next action stay, and
   the list is added.
5. Text output prints `Skills outdated: <n>` and one line per entry,
   `- outdated <skill>: found <found>, requires <required>`, after the
   `Skills drifted` lines, and only when the list is not empty.

### Fixed texts

**Roundfix skill, readiness paragraph.** Replace "For each Roundfix-owned
skill, the running binary declares a minimum version and compares it with the
version declared by the installed `SKILL.md`." with:

> For each Roundfix-owned skill, the minimum version is the version of that
> skill the running binary carries, and Doctor compares it with the version
> declared by the installed `SKILL.md`.

**Roundfix skill, managed refresh paragraph.** After "Without confirmation, a
changed Plan is presented and nothing is written." add:

> The same preview lists, under `skills.outdated`, each installed
> Roundfix-owned skill older than the version the binary carries. With
> unchanged guidance it then reports `plan_ready` instead of `current`, and
> `--yes` refreshes the Repository Skill Set.

Both version fields of the Roundfix skill rise by one patch step from their
value on the Task's base. The `### QA settlement` section is not touched. The
mirror is regenerated with `make skills-sync`.

**`docs/user-guide/commands.md`, the `skills:` bullet.** Add: "The minimum
version of an owned skill is the version of that skill the running binary
carries, so a copy installed by an older binary fails this line until it is
refreshed."

**`docs/user-guide/context-driven-development.md`.** After "Without
confirmation, a changed managed-refresh Plan is presented and the repository
remains unchanged." add: "The preview also lists each installed Roundfix-owned
skill older than the version the binary carries, and then reports
`plan_ready` even when the guidance is unchanged."

**`docs/agents/specific-repository.md`.** Add one rule: "An owned skill's
content changes only together with its version. Raise both version fields,
then record the version with
`go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
and declare `skills/testdata/owned-skill-versions.json` in the Task."

**Setup entry.** Appended as the last entry of `skills` in `go-cli.json` and
`rust-cli.json`, and of the `go-cli` and `rust-cli` setups in the parity
fixture:

```json
{
  "name": "roundfix",
  "path": "skills/06-review-repair/roundfix",
  "source": { "type": "repo", "name": "roundfix" },
  "minimumVersion": "0.0.2"
}
```

**`autonomous-work.json`.** `requiredSkills` becomes `["roundfix"]`, and
`skillDispatch` gains:

```json
{
  "skill": "roundfix",
  "triggers": [
    {
      "id": "trigger.autonomous-work.roundfix",
      "when": "Running, inspecting, recovering or delivering work with the `roundfix` command line, or reading its output."
    }
  ]
}
```

The module version rises by one. No rule, clause or guide version changes.

**`core.json` — `clause.core.plan-the-release-first`.** Append: "Before the
release Pull Request, confirm that the skills and guides the repository ships
describe the behavior being released; the repository's release runbook owns
that check." The clause stays `mandatory`. The versions of
`rule.core.git-delivery`, `guide.agent-instructions` and the `core` module
rise by one.

**`docs/user-guide/release-runbook.md`.** A new section, "Checking skills and
guides before the release", placed before "Cutting a release". It states that
the step is mandatory and runs before the release Pull Request, and lists:

1. `roundfix doctor`: the `skills:` line is `ok`.
2. `go test -count=1 ./skills -run '^(TestEveryOwnedSkillVersionIsRecorded|TestTheOwnedSkillMinimumIsTheEmbeddedVersion)$'`.
3. `go test -count=1 ./internal/baseline -run '^(TestNoTwoBaselineClausesShareText|TestBaselineClauseForceIsCharacterized|TestShippedGuidanceCitesNoRepositoryRecord)$'`
   and `go test -count=1 ./internal/delivery -run '^TestTheLoopClauseOrderMatchesTheDeliveryQueue$'`.
4. `go test -count=1 -tags docscontract ./internal/docscontract -run '^(TestEveryCommandIsNamedInTheRoundfixSkill|TestEveryCommandIsNamedInTheUserGuide)$'`.
5. `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format text`
   reports `current`.
6. A reading pass: for each user-visible change in the release range, confirm
   that the Roundfix skill, the user guide and the Baseline clause that
   describe it say what the release does. Fix a mismatch before the release
   Pull Request, or record a Backlog Entry when the fix needs a Spec.

"Cutting a release" gains a step, after the release plan step and before the
tag step, that points at the new section. The checks in items 3 and 4 belong
to Specs 0193 and 0192; the runbook names them and does not restate what they
assert.

### API Contracts

1. API Contract: `roundfix doctor` — for an installed owned skill below the
   bundle's version the line is
   `skills: failed (below minimum: skill "<skill>" requires <carried>, found <installed>; next: roundfix skills install --target project)`.
2. API Contract: `roundfix baseline update --repo <path> --format json`
   without confirmation — the result follows "The preview rule": `state`
   `plan_ready`, exit `3`, `skills.status` `outdated` and `skills.outdated`
   when an installed owned skill is older; byte-identical output otherwise.
3. API Contract: `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'`
   — fails for a changed digest under a recorded version and for an
   unrecorded version; `-record-skill-versions` adds a new, higher version
   and nothing else.
4. API Contract: rendered guide `docs/agents/skill-dispatch.md` — carries
   `trigger.autonomous-work.roundfix` for every profile.
5. API Contract: `roundfix baseline assets sync` — a refreshed snapshot keeps
   a `repo` entry the upstream list omits, after the upstream entries.
6. API Contract: rendered guide `docs/agents/agent-instructions.md` — the
   release clause carries the appended sentence.

## Coverage Map

- Goal 1 → Interfaces; The preview rule; Testing Approach 1.
- Goal 2 → Data Models; The record rule; Testing Approach 2.
- Goal 3 → Fixed texts (setup entry, `autonomous-work.json`); Testing
  Approach 3.
- Goal 4 → Fixed texts (runbook, release clause); Testing Approach 4.
- Core Feature 1 → Interfaces; The preview rule; API Contracts 1 and 2;
  Testing Approach 1.
- Core Feature 2 → The record rule; API Contract 3; Testing Approach 2.
- Core Feature 3 → Fixed texts; API Contracts 4 and 5; Testing Approach 3.
- Core Feature 4 → Fixed texts; API Contract 6; Testing Approach 4.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 1.
- Success Metric 3 → Testing Approach 2.
- Success Metric 4 → Testing Approach 3.
- Success Metric 5 → Testing Approach 4.
- Success Metric 6 → Testing Approach 5.
- API Contract 1 → Interfaces.
- API Contract 2 → The preview rule.
- API Contract 3 → The record rule.
- API Contract 4 → Fixed texts (`autonomous-work.json`).
- API Contract 5 → Interfaces.
- API Contract 6 → Fixed texts (`core.json`).

## Integration Points

- **Spec 0192 raises owned skill versions first.** It moves the Roundfix
  skill to `0.0.3` and eleven other skills one step. With the literal minimum
  that raise fails three tests in `./skills`; however Spec 0192 settles them,
  task_01 replaces the literal list and makes the tests read the map, so they
  hold for any version.
- **Spec 0193 edits the same modules first.** It rewrites clauses in
  `core.json` and `autonomous-work.json`. This Spec appends one sentence to
  another clause of `core.json` and adds no clause to `autonomous-work.json`.
  Spec 0193's force table lists clauses and levels, and neither changes here.
- **The Spec that refreshes the upstream skill snapshot.** Its asset sync
  reads upstream lists that omit the Roundfix skill. With this Spec's sync
  rule the entry survives; without it the catalog fails closed with
  `catalog.profile.skill.outside-setup`, because `autonomous-work` requires
  the skill.
- **Later Specs that edit an owned skill.** They raise the version and record
  it, and declare the record file in the Task.
- **The owned-skill edit suite.** `TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical`
  appends to a skill without raising its version and runs only
  `make baseline-digests`. The record test is not one of that command's
  steps, so the suite stays green.

## Testing Approach

All tests read embedded files, repository files or temporary repositories
they create. None needs a network. Each negative case is its own test.

1. **Minimum and reports.**
   - New `skills/owned_skill_minimum_test.go`:
     `TestTheOwnedSkillMinimumIsTheEmbeddedVersion` compares the map with
     every embedded version; `TestAMinimumThatDiffersFromTheEmbeddedVersionIsReported`
     gives the comparison a lowered entry and expects it back;
     `TestAnInstalledOwnedSkillOlderThanTheBundleIsBelow` installs the bundle
     into a temporary repository, lowers one skill by a patch step and
     expects the below state with both versions.
   - New `internal/cli/baseline_update_outdated_skills_test.go`:
     `TestBaselineUpdatePreviewReportsAnOlderOwnedSkill`,
     `TestBaselineUpdatePreviewStaysCurrentWhenOwnedSkillsMatch`,
     `TestBaselineUpdatePreviewWithNoSkillsSkipsTheOwnedSkillCheck` and
     `TestDoctorFailsForAnOwnedSkillOlderThanTheBundle`.
   - Existing `TestOwnedSkillContractRejectsSetAndVersionDisagreement`,
     `TestOwnedSkillBundleReadinessKeepsStatesDistinct`,
     `TestCheckRepositoryClassifiesMissingAndOutdatedSkills`,
     `TestBaselineUpdateFleetSweep` and `TestAuthorialSkillSync` stay green.
2. **Version record.** New `skills/owned_skill_versions_test.go`:
   `TestEveryOwnedSkillVersionIsRecorded`,
   `TestAChangedOwnedSkillUnderARecordedVersionIsRefused`,
   `TestAnUnrecordedOwnedSkillVersionIsRefused` and
   `TestRecordingNeverReplacesARecordedVersion`. The three negative tests run
   the rule on an in-memory record.
3. **Membership and dispatch.**
   - New `internal/baseline/roundfix_skill_membership_test.go`:
     `TestEverySetupListsTheRoundfixSkill`,
     `TestASetupWithoutTheRoundfixSkillIsReported` and
     `TestEveryProfileDispatchesTheRoundfixSkill`.
   - New `internal/baseline/assets_sync_owned_membership_test.go`:
     `TestAssetSyncKeepsAnOwnedSkillTheUpstreamListOmits` and
     `TestAssetSyncStillDropsAnExternalSkillTheUpstreamListOmits`.
   - Existing `TestAssetsSyncCompatibilityMatchesMaintainedPythonContract`
     stays green.
4. **Release step.**
   - New `internal/docscontract/release_step_test.go`, under the
     `docscontract` build tag:
     `TestTheReleaseRunbookRequiresTheSkillsAndGuidesCheck` expects the
     section, its six items and the pointer step before the tag step;
     `TestEveryCheckTheReleaseStepNamesExists` reads each `Test…` name in the
     section and expects a matching top-level test function in the
     repository; `TestAReleaseStepThatNamesAMissingCheckIsReported` feeds a
     section with an unknown name and expects it back.
   - New `internal/baseline/release_clause_test.go`:
     `TestTheReleaseClauseNamesTheSkillsAndGuidesCheck`.
   - Existing `TestReleasePlanDocumentationContract` stays green.
5. **Regeneration.** task_03 and task_04 run the existing
   `TestFormatterComposition`, `TestCatalogCompatibility`,
   `TestBaselinePlanCharacterization` and `TestBaselineCompatibilityCorpus`,
   then the read-only managed refresh, which exits non-zero while a plan is
   pending.

## Build Order

1. The minimum built from the bundle, Doctor's and the preview's report, and
   the two documentation edits, task_01 (depends on: none).
2. The version record, its test and the repository rule, task_02 (depends on:
   1).
3. The Roundfix skill's setup membership, its trigger and the sync rule,
   task_03 (depends on: 2).
4. The release step, the release clause sentence and their checks, task_04
   (depends on: 3).
5. Terminal QA, task_05 (depends on: 1, 2, 3, 4).

task_02 follows task_01 because task_01 raises the Roundfix skill's version
and the record must start from the versions that ship. task_03 and task_04
rewrite the same digest pin, catalog snapshots and plan goldens. task_04 names
task_01's and task_02's tests in the runbook.

## Risks & Considerations

- **Doctor fails after an upgrade.** This is the intended result, and the
  line prints its remedy. `roundfix baseline update --yes` also refreshes the
  owned skills.
- **Another Spec restructures the Roundfix skill.** The two skill edits are
  sentence replacements. If the sentences have moved, they are edited where
  they are; if one is gone, the Task states the behavior in the section that
  now describes Doctor's `skills:` line or the managed refresh.
- **The record can be bypassed by recording.** Recording refuses to replace a
  digest, so the only way past a changed skill is a new version, which is the
  rule.
- **The fixture edit is by hand.** The parity fixture's two setups gain the
  entry by hand; the regeneration recomputes only digests. The existing
  compatibility test fails when the fixture and the snapshots disagree.
- **Hazard: generated files.** A Task must never hand-edit a golden, a pin, a
  catalog snapshot, the parity manifest or a rendered guide.

## Decisions

- The minimum is read from the embedded skill; the literal list is removed.
  See ADR-0189.
- The preview counts an installed, older skill and nothing else.
- The record lives beside the test that reads it, outside the embedded bundle
  and outside the regenerated paths.
- The Roundfix skill's trigger lives in `autonomous-work`, the one module all
  three profiles select that is about operating Roundfix.
- The sync rule is general for `repo` entries, not a special case for one
  skill.
