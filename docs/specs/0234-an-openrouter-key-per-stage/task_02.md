---
task: task_02
spec: 0234-an-openrouter-key-per-stage
status: completed
type: backend
complexity: low
---

# Task 02: The Doctor Command lists each stage's keys, and the guides and skills name the stage keys

## Overview

`roundfix doctor`'s `environment:` line lists, by name and in preference
order, the judge's variables and the implementation stage's variables, each
set or not set, from the stage key list task_01 added. The doctor reference,
the configuration guide, the Roundfix Skill's runtime reference and the
write-prd and write-techspec skills name the stage keys, their order, the
shared-key fallback and `key_variable`.

## Requirements

1. MUST make the environment detail end with the text of API Contract 3,
   built from `Questions.KeyVariables` and
   `openrouterkey.Variables(openrouterkey.StageImplement)` per Invariant 11,
   with presence decided by the existing last-entry-wins reading. Its status,
   findings and `next:` suffix are unchanged, and no key value is retained or
   rendered. If the tree already reports the implementation key under
   another label, this line replaces that report.
2. MUST add `TestEnvironmentReadinessListsEachStageKeyByName`: with the judge
   key and the implementation key set to sentinel values and the shared key
   unset, the detail lists each of the five entries of API Contract 3 with
   its state in order, the doctor output contains no sentinel, and the
   generic `OPENROUTER_API_KEY` set alone leaves every entry `not set`. MUST
   update the golden `environment:` phrase of
   `TestDoctorPrintsTheFiveReadinessLines` to API Contract 3.
3. MUST describe in `docs/user-guide/commands/doctor.md` the new line,
   updating its sample output, and add to `docs/user-guide/configuration.md`
   a section on OpenRouter keys per stage that names both stages' variables
   in order, the shared-key fallback, the never-read generic key,
   `key_variable` in the Judge Log and in the implementation spend record,
   and that each stage keeps its own ceiling (`jev.monthly_ceiling_usd`,
   `openrouter.implement_monthly_ceiling_usd`).
4. MUST add to `.agents/skills/roundfix/references/runtime.md`, in the Doctor
   Command paragraph, that the `environment:` line lists `spec judge keys`
   and `implementation keys` by name, and in the OpenCode OpenRouter text
   that an open model reads `ROUNDFIX_OPENROUTER_IMPLEMENT_API_KEY` first and
   the shared key second.
5. MUST replace "either Jev key" in the Advisory judgment sections of
   `.agents/skills/write-prd/SKILL.md` and
   `.agents/skills/write-techspec/SKILL.md` with a sentence that forbids
   printing, storing or asking for any Jev key and names
   `ROUNDFIX_OPENROUTER_JUDGE_API_KEY`, `ROUNDFIX_OPENROUTER_API_KEY` and
   `ROUNDFIX_TYPESAFE_API_KEY`.
6. MUST run `make skills-sync`, re-record the changed owned skills' versions
   with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   after the last skill edit, then `make baseline-digests`, as implementation
   steps outside Verification; never inside a `### QA settlement` section,
   and never a digest or version by hand.
7. MUST NOT change any other doctor line, the doctor's exit codes,
   `CONTEXT.md` or `CHANGELOG.md`.

## Subtasks

- [ ] List both stages' variables in the environment line.
- [ ] Add the doctor test and update the golden phrase.
- [ ] Describe the line and the stage keys in the doctor reference and the configuration guide.
- [ ] Name the stage keys in the runtime reference and the two authoring skills.
- [ ] Sync and record the versions.

## Acceptance Criteria

- [ ] The environment line matches API Contract 3 and renders no key value.
- [ ] The doctor reference, the configuration guide and the runtime
      reference name both stages' variables and the fallback.
- [ ] write-prd and write-techspec forbid printing all three Jev key names.
- [ ] Each mirror equals its canonical file and the raised versions are
      recorded.

## Context

- instruction: `docs/adr/0239-each-openrouter-stage-reads-its-own-roundfix-key-first.md`
- interface: `internal/cli/readiness_toolchain.go`
- interface: `internal/cli/readiness_toolchain_test.go`
- interface: `docs/user-guide/commands/doctor.md`
- interface: `docs/user-guide/configuration.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/runtime.md`
- interface: `.agents/skills/write-prd/SKILL.md`
- interface: `.agents/skills/write-techspec/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/roundfix/references/runtime.md`
- interface: `skills/write-prd/SKILL.md`
- interface: `skills/write-techspec/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`

## Verification

- `out="$(go test -count=1 -v -run '^(TestEnvironmentReadinessListsEachStageKeyByName|TestEnvironmentReadinessNamesAMissingPreloadAndNeverAKeyValue|TestDoctorPrintsTheFiveReadinessLines)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestEnvironmentReadinessListsEachStageKeyByName TestEnvironmentReadinessNamesAMissingPreloadAndNeverAKeyValue TestDoctorPrintsTheFiveReadinessLines; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the new test does not exist, so its pass line is missing and the command fails.
- `for phrase in 'implementation keys' 'ROUNDFIX_OPENROUTER_IMPLEMENT_API_KEY'; do tr -s '[:space:]' ' ' < docs/user-guide/commands/doctor.md | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/doctor.md "$phrase" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' docs/user-guide/configuration.md "$phrase" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/runtime.md | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/runtime.md "$phrase" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- 'openrouter.implement_monthly_ceiling_usd' || { printf 'missing ceiling in configuration guide\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/write-prd/SKILL.md | grep -qF -- 'ROUNDFIX_OPENROUTER_JUDGE_API_KEY' || { printf 'missing judge key in write-prd\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/write-techspec/SKILL.md | grep -qF -- 'ROUNDFIX_OPENROUTER_JUDGE_API_KEY' || { printf 'missing judge key in write-techspec\n' >&2; exit 1; }; cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && cmp .agents/skills/roundfix/references/runtime.md skills/roundfix/references/runtime.md && cmp .agents/skills/write-prd/SKILL.md skills/write-prd/SKILL.md && cmp .agents/skills/write-techspec/SKILL.md skills/write-techspec/SKILL.md && out="$(go test -count=1 -v -run '^(TestEveryOwnedSkillVersionIsRecorded)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- "--- PASS: TestEveryOwnedSkillVersionIsRecorded" || { printf 'missing pass: TestEveryOwnedSkillVersionIsRecorded\n' >&2; exit 1; }` — expected: exit 0; before this Task no guide names the implementation keys and neither authoring skill names the judge key, so the command fails; after it the mirrors equal their canonical files and the raised versions are recorded.

## References

- `_prd.md` → Goals; User Story 3; Core Features 3, 5, 6; Success Metric 4; Success Metric 5
- `_techspec.md` → Invariant 11; API Contract 3; Vocabulary Contract; Testing Approach; Build Order 2
- ADR-0239; ADR-0189

## Result

Implemented the Task 02 slice for Daemon Verification. Task status and the
authored Verification commands remain Daemon-owned; no commit, push or Pull
Request was made.

Acceptance-criterion evidence:

- **Environment line and key secrecy:** `environmentReadiness` now renders
  `Questions.KeyVariables()` followed by
  `openrouterkey.Variables(openrouterkey.StageImplement)`, with the two
  labels and five ordered entries from API Contract 3. The existing
  presence-only, last-entry-wins reader is unchanged, as are findings,
  status and remediation. `TestEnvironmentReadinessListsEachStageKeyByName`
  covers stage-key sentinels, generic-key exclusion, duplicate entries and
  Doctor output without key values. The warning transcript now includes
  both lists before its existing `next:` suffix.
- **Guides and runtime reference:** the Doctor reference and sample output,
  configuration guide and Roundfix runtime reference name each stage's
  variables in preference order and the shared-key fallback. The
  configuration guide also names the unread generic key, `key_variable` in
  both records and both independent ceilings. The existing light-tier text
  was updated to match the stage-key rule.
- **Authoring skill secrecy:** both Advisory judgment sections now forbid
  printing, storing or asking for any Jev key and explicitly name
  `ROUNDFIX_OPENROUTER_JUDGE_API_KEY`, `ROUNDFIX_OPENROUTER_API_KEY` and
  `ROUNDFIX_TYPESAFE_API_KEY`.
- **Mirrors and versions:** `make skills-sync` regenerated the mirrors.
  The required version recorder raised and recorded Roundfix `0.1.37`,
  write-prd `0.0.7` and write-techspec `0.0.7` in both version fields and
  `skills/testdata/owned-skill-versions.json`. A byte comparison of every
  file in these three skill trees confirmed all mirrors match.

Focused checks and implementation steps:

- `GOCACHE=/private/tmp/roundfix-task02-gocache go test ./internal/cli -run
  '^TestEnvironmentReadinessListsEachStageKeyByName$' -count=1`: before the
  production change, all three subtests failed because Doctor omitted the
  stage variables and implementation list, establishing the red signal.
- `gofmt -w internal/cli/readiness_toolchain.go
  internal/cli/readiness_toolchain_test.go`: exit 0.
- `GOCACHE=/private/tmp/roundfix-task02-gocache go test ./internal/cli -run
  'TestEnvironmentReadiness|TestDoctorPrintsTheFiveReadinessLines' -count=1`:
  exit 0 after the implementation, including the existing preload and
  warning-output checks.
- `make skills-sync`: exit 0.
- `GOCACHE=/private/tmp/roundfix-task02-gocache go test ./skills -run
  '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`: the first
  attempt was blocked by the sandbox's read-only `.agents` directory; the
  same required generator with approved filesystem permission exited 0.
- `GOCACHE=/private/tmp/roundfix-task02-gocache make baseline-digests`:
  exit 0; `ok: true`, `changed: false`, derived artifacts already match.
- Python byte comparison of canonical and mirrored skill trees: exit 0;
  Roundfix 20 files, write-prd 2 files, write-techspec 3 files match.
- `git -c core.fsmonitor=false diff --check`: exit 0.

The declared Verification commands were not run. No follow-up work was
identified within this slice.
