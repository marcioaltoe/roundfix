---
spec: 0318-the-baseline-reports-a-required-skill-that-is-not-installed
prd: _prd.md
created: 2026-10-25
---

# The Baseline reports a required skill that is not installed — Technical Spec

## Executive Summary

The `roundfix baseline update` preview already inspects the external
Repository Skill Set, but only to compute `skills.differing`; it never reads
the names the inspection proves missing or lock-mismatched, so it reports
`current` while `--yes` would restore those skills. One helper in
`internal/cli` turns the inspection's result into a list of Missing Required
Skills with their restore command. The preview lists them under
`skills.missing` and reports `plan_ready`, exit `3`, as it does for an owned
skill older than the binary's. `roundfix skills check` uses the same helper on
the repository that contains its work directory and adds one diagnostic per
skill to its existing failure list. Doctor already fails on the class; a test
pins it. No flag, Baseline module or rendered guide changes.

The trade-off accepted: an adopter whose upstream retired a skill the Baseline
still dispatches now gets exit `3` from the preview until the catalog stops
requiring it (Spec 0317 for `handoff`) or the skill is restored; that is the
point, and the reason names the restore command.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local file inspection only.
  Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: not applicable — no Baseline module, clause,
  trigger or Setup Snapshot changes, as `_prd.md` → Project Constraints
  states. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — protected paths limited to one test file
  and the Roundfix skill. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0318-the-baseline-reports-a-required-skill-that-is-not-installed/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `.agents/skills/roundfix/references/setup.md`,
  `internal/cli/cli_test.go`, `skills/roundfix/SKILL.md`.

## Current behavior

Read on 2026-10-25 at `f3c39140` (v0.119.0).

1. `internal/cli/baseline_update.go`, preview branch (`if !request.yes &&
   request.confirmation == ""`, line 285): with skills not skipped it checks
   the owned skills for `ReadinessBelow` (`skills.outdated`), then calls
   `resolveExternalSkillRequirement(request.repo)` and, when the manifest is
   readable and names external skills,
   `roundskills.CheckRepositoryWithExternal(ctx, request.repo, external)`. Its
   result feeds only `baseline.SetupSnapshotDifferences` through
   `snapshotComparableSkills` (`internal/cli/doctor.go:435`), which drops
   `MissingExternal` and `OutdatedExternal` and returns nothing when the check
   returned an error. With no file change, no history move and no outdated
   owned skill the state is `current`, exit `0` (line 337).
2. `skills/repository.go` `CheckRepositoryWithExternal`: a required external
   skill without a `skills-lock.json` entry makes `requiredExternalHashes`
   return a `*RepositoryReadinessError` whose `MissingExternal` names it
   (line 445), before any folder is inspected; a skill with a lock entry and
   no `.agents/skills/<name>` folder is appended to
   `readiness.MissingExternal` (line 203) with a nil error; a folder whose
   hash differs from its lock entry lands in `readiness.OutdatedExternal`. An
   absent or unreadable lock returns an error that names no skill.
3. The applied path (`--yes` or a matching `--confirm-plan`) runs the skills
   stage, whose `baselineUpdateExternalDrift` (line 600) collects exactly
   those three sources and restores each skill from its Setup Snapshot, or
   lists it under `skills.drifted` when the source is unreachable.
4. The precedent for a skill the preview would act on:
   `TestBaselineUpdatePreviewReportsAnOlderOwnedSkill`
   (`internal/cli/baseline_update_outdated_skills_test.go`, Spec 0195) pins
   `state: plan_ready`, `category: approval`, exit `3`, `skills.status:
   outdated`, the message "guidance matches the current Baseline catalog; 1
   Roundfix-owned skill(s) are older than the ones this binary carries" and a
   `--yes` next action. `TestBaselineUpdatePreviewStaysCurrentWhenOwnedSkillsMatch`
   pins that owned skills absent from the repository leave the preview
   `current`. `TestBaselineUpdateListsRetiredSkillsWithoutChangingItsState`
   pins `current` for a fixture whose lock holds only a retired entry.
5. `roundfix doctor` (`repositorySkillsCheck`, `internal/cli/doctor.go:188`)
   reports `skills: failed (...)` with the missing names and `next: bunx
   skills add marcioaltoe/skills@<name>`, exit `1`, for both removal shapes
   (reproduced; see `_prd.md` → Validation). No test runs it on a real
   adopted repository with a removed required skill.
6. `roundfix skills check` (`internal/cli/cli.go`, `case "check"` in
   `runSkillsCommand`) validates only the binary's embedded owned bundle
   (`checkOwnedSkills`): diagnostics go to stderr after
   `roundfix: Roundfix skill check failed:` as `  <path>: <message>` lines,
   exit `1`; otherwise `Roundfix skill check passed: <names>` or the
   `unversioned` line, exit `0`. This repository's `make verify` runs it
   (`Makefile:262`); this repository's Repository Skill Set is complete
   (`roundfix doctor`: `skills: ok (41 required: 13 Roundfix-owned, 28
   external)`).
7. `TestRunSkillsCheck` (`internal/cli/cli_test.go`) and
   `TestDoctorAndSkillsCheckReportSharedOwnedSkillReadiness`
   (`internal/cli/doctor_test.go`) run `skills check` with the test process's
   working directory, which is inside this repository.
8. The docs: `docs/user-guide/commands/baseline.md` and
   `.agents/skills/roundfix/references/baseline.md` describe
   `skills.differing` and end that paragraph with "A missing skill, or one
   that no longer matches its lock, is still restored under
   `skills.drifted`."; `docs/user-guide/commands/skills.md` describes `skills
   check` as the bundle check; `.agents/skills/roundfix/references/setup.md`
   carries the Doctor `skills` line. `CONTEXT.md` has Repository Skill Set,
   Setup Snapshot Difference and Retired Skill, and no term for a required
   skill that is not installed.
9. A guide edit makes `TestTheSkillCoverageMapIsCurrent` report a stale
   fingerprint until `docs/references/behavior-surfaces.json` is re-recorded;
   in the prototype the `skills` help text and both guides moved four
   entries (`command skills check`, `command skills install`, and the two
   guides).

## System Architecture

- `internal/cli/required_skills.go` (new, task_01): the helper that turns a
  Repository Skill Set inspection into Missing Required Skills, and the Setup
  Manifest profile reader.
- `internal/cli/skills_check_repository.go` (new, task_02): the repository
  section of `skills check`.
- `internal/cli/baseline_update.go`: the preview list, status, state, text.
- `internal/cli/cli.go`: `skills check` appends the repository diagnostics;
  the `skills` usage line.
- Tests: `internal/cli/baseline_update_missing_skills_test.go` (new),
  `internal/cli/skills_check_repository_test.go` (new),
  `internal/cli/baseline_update_retired_skills_test.go`,
  `internal/cli/cli_test.go`, `internal/cli/doctor_test.go`.
- Docs: the two guides, the two Roundfix references and their mirrors, the
  Roundfix skill's version pair, `skills/testdata/owned-skill-versions.json`,
  `CONTEXT.md`, `docs/references/behavior-surfaces.json`.

## Implementation Design

### Interfaces

The helper, in `internal/cli/required_skills.go`:

```go
// missingRequiredSkills names each external skill the adopted modules require
// that the inspection proved absent from .agents/skills or skills-lock.json,
// or no longer matching its lock entry. An error naming no skill adds nothing.
func missingRequiredSkills(profile string, readiness roundskills.RepositoryReadiness, checkErr error) []baselineUpdateSkillDrift
```

Sources, in this order: `readiness.MissingExternal`, then the
`MissingExternal` of a `*roundskills.RepositoryReadinessError` found with
`errors.As`, both with the not-installed reason; then
`readiness.OutdatedExternal` with the mismatch reason unless already named.
The result is sorted by skill name, one entry per skill.

Reasons, where `<id>` is the Setup Manifest's profile and `<name>` the skill:

```text
required by the Baseline and not installed; restore with roundfix baseline skills restore --profile <id> --skill <name>
required by the Baseline and no longer matches its skills-lock.json entry; restore with roundfix baseline skills restore --profile <id> --skill <name>
```

The preview, unchanged guidance, one or more Missing Required Skills:

```text
state:      plan_ready
category:   approval
message:    guidance matches the current Baseline catalog; <n> required skill(s) are missing or no longer match skills-lock.json: <name>, <name>
nextAction: rerun with --yes to restore them from the Setup Snapshot, or restore one with roundfix baseline skills restore --profile <id> --skill <name>
exit:       3
```

The text output adds, after the `Skills outdated` block and before `Skills
retired`, only when the list is not empty:

```text
Skills missing: <n>
- missing <name>: <reason>
```

`skills check`, inside an adopted repository, one diagnostic per skill after
the bundle diagnostics, through the existing failure printer:

```text
roundfix: Roundfix skill check failed:
  .agents/skills/<name>: <reason>
```

The `skills` usage entry for `check` becomes:

```text
  check      Validate shipped Roundfix skill artifacts and, inside an adopted
             repository, the skills its Baseline requires
```

The guide paragraph, appended to `docs/user-guide/commands/baseline.md` and
`.agents/skills/roundfix/references/baseline.md` directly after "A missing
skill, or one that no longer matches its lock, is still restored under
`skills.drifted`.":

```markdown
The preview without `--yes` or `--confirm-plan` lists, under `skills.missing`
(`Skills missing` in text), each external skill the adopted modules require and
dispatch that is absent from `.agents/skills` or `skills-lock.json`, or no
longer matches its lock entry: a Missing Required Skill. Each entry names the
skill and the command that restores it,
`roundfix baseline skills restore --profile <id> --skill <name>`.
`skills.status` reads `missing`; with unchanged guidance the state is
`plan_ready` instead of `current` and the exit code is `3`, as for an owned
skill older than the binary's, and `--yes` restores the listed skills.
A skill the upstream catalog retired while the Baseline still dispatches it
is a Baseline defect to report, not a skill to keep restoring by hand.
`--no-skills` skips the list, and an absent or unreadable `skills-lock.json`
adds nothing to it; `roundfix doctor` reports that state.
```

The paragraph appended to `docs/user-guide/commands/skills.md`:

```markdown
Inside an adopted Git repository, `skills check` also inspects the skills the
repository's Baseline requires and fails with exit `1` for each Missing
Required Skill, printing
`.agents/skills/<name>: required by the Baseline and not installed; restore with roundfix baseline skills restore --profile <id> --skill <name>`
on stderr. Outside a Git repository, or without a Setup Manifest, it checks
only the shipped bundle.
```

The paragraph appended to `.agents/skills/roundfix/references/setup.md` after
the `DR-SKILL-DIFFERS-FROM-SNAPSHOT` paragraph that ends "nothing is
restored.":

```markdown
`roundfix skills check` validates the shipped bundle and, inside an adopted
Git repository, fails with exit `1` for each Missing Required Skill, naming it
with `roundfix baseline skills restore --profile <id> --skill <name>`.
```

The `CONTEXT.md` entry, placed before `**Park Class**:`:

```markdown
**Missing Required Skill**:
An external skill the adopted modules require and dispatch that is absent from `.agents/skills` or `skills-lock.json`, or no longer matches its lock entry. The `roundfix baseline update` preview lists it under `skills.missing` with its restore command and reports `plan_ready` instead of `current`; `roundfix doctor` and `roundfix skills check` fail on it.
_Avoid_: Uninstalled skill, absent skill, dispatched skill gap
```

### Data Models

`baselineUpdateSkillsResult` gains, after `Outdated`:

```go
Missing []baselineUpdateSkillDrift `json:"missing,omitempty"`
```

and the status constant `baselineUpdateSkillsMissing = "missing"`, set when
the list is not empty (it overrides `outdated`). The applied path never fills
`missing`.

### API Contracts

1. API Contract: the preview (no `--yes`, no `--confirm-plan`, no
   `--no-skills`) fills `skills.missing` from `missingRequiredSkills` with the
   manifest's profile, using the same inspection call it already makes, and
   sets `skills.status` to `missing`.
2. API Contract: with no file change and no history move, a non-empty list
   returns the state, category, message, next action and exit `3` of
   Interfaces, checked before the outdated-owned branch; when guidance also
   changed, the existing `plan_ready` message and next action are kept and
   the list still appears.
3. API Contract: an empty list leaves every byte of the preview unchanged,
   including no `missing` key in JSON and no `Skills missing` line in text.
   An inspection error that names no skill adds nothing.
4. API Contract: `--no-skills`, `--yes`, `--confirm-plan` and the applied
   path are unchanged.
5. API Contract: `skills check`'s repository section, in
   `internal/cli/skills_check_repository.go`, resolves the work directory with
   `environment.resolveWorkDir`, the repository root with
   `environment.dependencies.resolveSkillsProjectRoot`, the required set with
   `environment.dependencies.doctor.resolveExternal`, and the inspection with
   `environment.dependencies.doctor.checkSkills`. When any of them fails, the
   manifest is unreadable or names no external skill, it adds nothing, so the
   output and exit code are those of the bundle check. Otherwise it appends
   one diagnostic per Missing Required Skill (path `.agents/skills/<name>`,
   message its reason) to the bundle diagnostics, and any diagnostic makes it
   print the existing failure block and exit `1`.
6. API Contract: Doctor's bytes are unchanged.

### Surface Transcripts

None. The new outputs are fixed in Interfaces and asserted by named tests.

### Invariants

```text
1. The baseline update preview never reports current while a required external skill is absent from .agents/skills or skills-lock.json, or no longer matches its lock entry.
2. Every Missing Required Skill is named with roundfix baseline skills restore --profile <id> --skill <name>.
3. A complete Repository Skill Set, --no-skills, an unreadable lock and a non-adopted directory keep their current bytes and exit codes.
4. Doctor, skills.differing, the applied path and a Retired Skill's report are unchanged.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, Core Feature 1, Success Metrics 1 and 2 → API Contracts 1-4, Invariants 1 and 3, task_01
- Goal 2, Core Feature 1, Success Metric 1 → API Contracts 1-2, Invariant 2, task_01
- Goal 3, Core Feature 2, Success Metric 3 → API Contract 5, Invariants 2 and 3, task_02
- Goal 4, Core Feature 3, Success Metric 4 → API Contract 6, Invariant 4, task_01
- Goal 5, Core Feature 4, Success Metric 6 → Interfaces (guides, references, glossary), task_01 and task_02
- Success Metric 5 → Invariant 4, task_01

## Absence Evidence

1. A preview list of missing required skills — command: `grep -rn -e "Skills.Missing" -e "skills.missing" internal/cli`; result: no match at `f3c39140` (the only `MissingExternal` in the update tests, `internal/cli/baseline_update_test.go:397`, feeds the applied path's skills stage).
2. A repository input to skills check — command: `grep -n -A12 'case "check":' internal/cli/cli.go`; result: only `checkOwnedSkills` at `f3c39140`.

## ADR Alternatives

None. No ADR is listed: no Baseline clause, Setup Snapshot rule or
architecture decision changes, and `skills.differing` keeps never restoring a
lock-matching copy.

## Integration Points

None. No network, provider, forge or other repository is reached; the
inspection reads `docs/agents/setup-context.json`, `skills-lock.json` and
`.agents/skills/` of the repository under test.

## Testing Approach

Every test builds temporary repositories; none reads this repository's
skills or the host's.

- `internal/cli/baseline_update_missing_skills_test.go` (new), with a
  `Suite:`/`Invariant:`/`Boundary` header:
  - helper `installRequiredExternalSkills(t, repository) []string`: reads the
    required external set with `resolveExternalSkillRequirement`, writes a
    stub `SKILL.md` folder for each under `.agents/skills/`, writes a
    `skills-lock.json` entry whose `computedHash` is
    `skills.SkillFolderHash` of that folder, keeps existing lock entries, and
    commits. Tests pick the skill to remove from the returned list, never by
    a literal name, so the fixture survives Spec 0317's catalog change.
  - helper `removeInstalledSkill(t, repository, name, lockEntry bool)`.
  - `TestBaselineUpdatePreviewNamesAMissingRequiredSkill`, on
    `newBaselineUpdateRepository` (profile `go-cli-tui`), subtests
    `folder and lock entry removed`, `folder removed with a stale lock
    entry`, `folder no longer matches its lock entry` (rewrite its
    `SKILL.md`): JSON exit `3`, `ok:false`, `plan_ready`, `approval`,
    `skills.status: missing`, `skills.missing` equal to the one entry with
    its exact reason, the exact message and next action; text holds `Skills:
    missing` and the `Skills missing` block; `--no-skills` stays `current`
    with exit `0`; the repository tree is unchanged. Subtest `every required
    skill installed`: `current`, exit `0`, no `"missing"` key.
  - `TestDoctorNamesAMissingRequiredSkill`, subtests `folder and lock entry
    removed` and `folder removed with a stale lock entry`: Doctor with fake
    health checks (`withDoctorFakeLoadedAndReadiness`, `GitRoot` the
    repository) and the real `resolveExternalSkillRequirement` and
    `skills.CheckRepositoryWithExternal`: exit `1`, a `skills: failed (`
    line naming the skill and `bunx skills add marcioaltoe/skills@<name>`.
- `internal/cli/baseline_update_retired_skills_test.go`:
  `adoptedRepositoryWithRetiredSkills` calls
  `installRequiredExternalSkills` after adoption, keeping the retired
  `the-fool` lock entry, so its preview stays `current` (declared change: the
  fixture, not the assertions).
- `internal/cli/skills_check_repository_test.go` (new):
  `TestSkillsCheckNamesAMissingRequiredSkill`, subtests `adopted repository
  with a missing required skill` (exact stderr, empty stdout, exit `1`),
  `adopted repository with every required skill`, `repository without a
  Setup Manifest` (`newBaselinePlanTestRepository`) and `outside a Git
  repository` (`t.TempDir()`), each exit `0` with the pass line; the work
  directory is set with `setCommandWorkDirForTest`.
- `internal/cli/cli_test.go` `TestRunSkillsCheck` and
  `internal/cli/doctor_test.go`
  `TestDoctorAndSkillsCheckReportSharedOwnedSkillReadiness`: set the work
  directory to `t.TempDir()` so they never read this repository (declared
  change: the setup line, not the assertions).
- Unchanged and run: `TestBaselineUpdatePreviewReportsAnOlderOwnedSkill`,
  `TestBaselineUpdatePreviewStaysCurrentWhenOwnedSkillsMatch`, `go test
  ./skills`, and the documentation contracts.
- The prototype (see `_prd.md` → Validation) passed the whole
  `./internal/cli`, `./skills` and `-tags docscontract
  ./internal/docscontract` packages, and each of these mutations failed the
  named subtest: ignoring the lock-entry error's names (`folder and lock
  entry removed`), ignoring `readiness.MissingExternal` (`folder removed
  with a stale lock entry`), ignoring `readiness.OutdatedExternal` (`folder
  no longer matches its lock entry`), letting the preview fall through to
  `current` (`folder and lock entry removed`), reading the list under
  `--no-skills` (`folder and lock entry removed`, its `--no-skills` check),
  dropping the repository diagnostics from `skills check` (`adopted
  repository with a missing required skill`), and reporting a diagnostic when
  no Setup Manifest exists (`repository without a Setup Manifest`).

## Build Order

1. The helper, the preview list, the Doctor characterization, the retired
   fixture, the `baseline` guide and reference paragraph, the glossary term,
   the Roundfix skill version and the Behavior Surface record re-recorded last
   (task_01; no dependency).
2. The `skills check` repository section, its usage line, its tests, the two
   hermetic test fixes, the `skills` guide and `setup` reference paragraphs,
   the Roundfix skill version and the Behavior Surface record re-recorded last
   (task_02; depends on 1 for the helper and the shared skill, version and
   record files).
3. QA gate (task_03) (depends on: 1, 2).

No Spec is listed in `requires`. Spec 0317 delivers first by the maintainer's
order; nothing here reads its code.

## Risks & Considerations

- An adopter that removed a still-required skill now sees exit `3` from the
  preview, and a CI job that runs `roundfix baseline update` without `--yes`
  and fails on a nonzero exit turns red. That is the intended signal; the
  reason names the restore command and the guide says a skill the upstream
  retired is a Baseline defect to report.
- `skills check` now reads the repository it runs in. This repository's
  `make verify` runs it at the root, where the Repository Skill Set is
  complete; Spec 0317 changes this repository's lock and skills and the
  `core` module together, so they stay consistent.
- `skills check` spawns `git rev-parse` only when the work directory has no
  `.git` above it, as `skills install --target project` already does.

## Glossary

- adds: **Missing Required Skill**

## Decisions

- Severity: `plan_ready`, exit `3`, as for an outdated owned skill (Spec
  0195), because `--yes` acts on it; not a `warnings` entry with exit `0`.
- Scope: external required skills only; absent owned skills and an absent
  lock keep their current preview, each named in `_prd.md` → Non-Goals.
- `skills check`: no new flag; the repository section is implicit, like
  `skills install --target project`, and silent outside an adopted
  repository.
- Doctor: a characterization test, no byte change.
- No ADR.
