---
task: task_01
spec: 0234-an-openrouter-key-per-stage
status: pending
type: backend
complexity: medium
---

# Task 01: The judge reads its own OpenRouter key first and names the variable it used

## Overview

Adds the stage key list of ADR-0239 as a leaf package and puts the judge on
it: `roundfix spec judge` reads `ROUNDFIX_OPENROUTER_JUDGE_API_KEY` before the
shared `ROUNDFIX_OPENROUTER_API_KEY` and TypeSafe direct, names the variable
it used in its summary, its JSON report and every Judge Log line, and names
the judge key first when it skips. The `spec judge` reference and the
Roundfix Skill's spec reference say so. This serves the "Separate OpenRouter
keys per stage" bullet of the Backlog Entry "A judge-assigned model tier per
Task" of 2026-10-05, which Spec 0233 owns.

## Requirements

1. MUST add package `internal/openrouterkey` with the constants, `Stage`
   values, `Variables` and `Select` of `_techspec.md` → Interfaces, obeying
   Invariants 1 and 2. It imports nothing outside the standard library and
   never returns or stores a value read from the environment.
2. MUST add `Transport.KeyVariables` and `Questions.KeyVariables` per
   Invariants 3 and 4, without changing `questions.json`, and make
   `selectTransport` return the selected variable per Invariant 5.
3. MUST add `KeyVariable` to `Report` (`key_variable`, null without a key) and
   to the Judge Log line (`key_variable`) per Invariant 6, build the no-key
   skip reason over `Questions.KeyVariables` per Invariant 7, and snapshot
   every variable's value for redaction per Invariant 8. Older Judge Log
   lines without the field still parse and still count toward the ceiling.
4. MUST make `runSpecJudgeCommand` collect exactly the variables of
   `Questions.KeyVariables` from its command environment (Invariant 9),
   end the text summary with `on <variable>` per Invariant 10, and replace
   the help's key paragraph with the four lines of Surface Transcript 2,
   byte for byte, keeping every other help line.
5. MUST add the tests named in Verification: `Variables` order for both
   stages; `Select` with the stage key, the shared key only, empty values,
   a repeated name whose last entry wins, and `OPENROUTER_API_KEY` and
   `TYPESAFE_API_KEY` alone giving no key; a judge run with both OpenRouter
   variables set to different fake values sends only the judge key, through
   a fake transport that asserts the bearer value, and reports and logs
   `ROUNDFIX_OPENROUTER_JUDGE_API_KEY`; with only the shared key it sends and
   names the shared key; no Judge Log line contains either fake value; and
   `spec judge` prints `via openrouter on ROUNDFIX_OPENROUTER_JUDGE_API_KEY`
   and, without keys, the skip line of Surface Transcript 1.
6. MUST update the existing tests whose expected output this Task changes,
   and only those: the skip reason pinned in `internal/judge/judge_test.go`,
   the Judge Log field set in `internal/judge/log_test.go`, the summary, skip
   and JSON field-set goldens in `internal/cli/spec_judge_test.go`, and the
   `selectTransport` call in `internal/judge/task_acceptance_measure_test.go`,
   which only adapts to the new return value.
7. MUST describe the key order, the fallback, the never-read generic key and
   `key_variable` in `docs/user-guide/commands/spec.md` and in the key
   paragraph of `.agents/skills/roundfix/references/spec.md`, never inside a
   `### QA settlement` section; run `make skills-sync`; and re-record the
   Roundfix Skill's version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   after the last skill edit, then `make baseline-digests`. Running those
   commands is an implementation step, never part of Verification, and no
   digest or version is written by hand.
8. MUST NOT change the judge's recipients, endpoints, request models, the
   ceiling or its sum, `questions.json`, `CONTEXT.md` or `CHANGELOG.md`, and
   MUST NOT read the process environment in a test.

## Subtasks

- [ ] Add the stage key list package and its tests.
- [ ] Map the judge's OpenRouter transport onto the judge stage.
- [ ] Record and print the variable; rewrite the skip reason and help.
- [ ] Add the judge and `spec judge` tests and update the pinned goldens.
- [ ] Describe it in the `spec judge` reference and the Roundfix Skill, sync and record the version.

## Acceptance Criteria

- [ ] With both OpenRouter variables set, only the judge key is sent and
      every record names `ROUNDFIX_OPENROUTER_JUDGE_API_KEY`.
- [ ] With only the shared key, the judge behaves as before and names the
      shared key.
- [ ] Without a key, the skip line of Surface Transcript 1 is printed and
      `key_variable` is null in JSON.
- [ ] No Judge Log line, report or output contains a key value.
- [ ] The guide and the skill reference name the order, each mirror equals
      its canonical file, and the raised version is recorded.

## Context

- instruction: `docs/adr/0239-each-openrouter-stage-reads-its-own-roundfix-key-first.md`
- instruction: `docs/adr/0201-the-judge-sends-only-spec-artifacts-and-spends-under-a-monthly-ceiling.md`
- instruction: `internal/judge/questions.json`
- creates: `internal/openrouterkey/openrouterkey.go`
- creates: `internal/openrouterkey/openrouterkey_test.go`
- creates: `internal/judge/stage_key_test.go`
- interface: `internal/judge/questions.go`
- interface: `internal/judge/client.go`
- interface: `internal/judge/judge.go`
- interface: `internal/judge/log.go`
- interface: `internal/judge/judge_test.go`
- interface: `internal/judge/log_test.go`
- interface: `internal/judge/task_acceptance_measure_test.go`
- interface: `internal/cli/spec_judge.go`
- interface: `internal/cli/spec_judge_test.go`
- interface: `docs/user-guide/commands/spec.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/spec.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/roundfix/references/spec.md`
- interface: `skills/testdata/owned-skill-versions.json`

## Verification

- `out="$(go test -count=1 -v -run '^(TestVariablesListEachStageKeyBeforeTheSharedKey|TestSelectReturnsTheFirstSetVariableByName|TestRunPrefersTheJudgeStageKey|TestRunFallsBackToTheSharedKeyAndNamesIt|TestJudgeLogNamesTheKeyVariableAndNoKey|TestSpecJudgeNamesTheKeyVariableItUsed|TestSpecJudgeSkipNamesTheJudgeKeyFirst)$' ./internal/openrouterkey ./internal/judge ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestVariablesListEachStageKeyBeforeTheSharedKey TestSelectReturnsTheFirstSetVariableByName TestRunPrefersTheJudgeStageKey TestRunFallsBackToTheSharedKeyAndNamesIt TestJudgeLogNamesTheKeyVariableAndNoKey TestSpecJudgeNamesTheKeyVariableItUsed TestSpecJudgeSkipNamesTheJudgeKeyFirst; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the package and the seven tests do not exist, so the command fails.
- `for phrase in 'ROUNDFIX_OPENROUTER_JUDGE_API_KEY' 'key_variable'; do tr -s '[:space:]' ' ' < docs/user-guide/commands/spec.md | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/spec.md "$phrase" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/spec.md | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/spec.md "$phrase" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < internal/cli/spec_judge.go | grep -qF -- 'Keys, first set wins: ROUNDFIX_OPENROUTER_JUDGE_API_KEY' || { printf 'missing help paragraph in internal/cli/spec_judge.go\n' >&2; exit 1; }; cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && cmp .agents/skills/roundfix/references/spec.md skills/roundfix/references/spec.md && out="$(go test -count=1 -v -run '^(TestEveryOwnedSkillVersionIsRecorded)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- "--- PASS: TestEveryOwnedSkillVersionIsRecorded" || { printf 'missing pass: TestEveryOwnedSkillVersionIsRecorded\n' >&2; exit 1; }` — expected: exit 0; before this Task neither the guide nor the skill reference names the judge key and the help has no key paragraph, so the command fails; after it the mirrors equal their canonical files and the raised version is recorded.

## References

- `_prd.md` → Goals; User Stories 1, 2, 4; Core Features 1, 3, 4, 6; Success Metric 1; Success Metric 2; Success Metric 5; Acceptance evidence
- `_techspec.md` → Interfaces; Invariants 1 to 10; Data Models; API Contract 1; API Contract 2; Surface Transcript 1; Surface Transcript 2; Vocabulary Contract; Testing Approach; Build Order 1
- ADR-0239; ADR-0201; ADR-0231; ADR-0189
