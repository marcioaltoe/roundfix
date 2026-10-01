---
task: task_03
spec: 0209-sources-that-share-a-context-share-a-spec
status: pending
type: backend
complexity: medium
---

# Task 03: `roundfix spec judge` prints and counts grouping suggestions

## Overview

task_02 makes the judge ask the grouping question, and the command does not
show its outcome. This Task prints each suggested pair, counts suggestions in
the summary line, lists grouping judgments in the JSON document, names the
question in the help, and describes it in the command reference and the
Roundfix Skill. It is verifiable on its own: with Spec 0205's fake transport
and a temporary Roundfix Home, the command reproduces Surface Transcripts 1
to 3 of this Spec and still reproduces Spec 0205's transcripts with
`0 suggested` in their summary.

## Requirements

1. MUST print, in `internal/cli/spec_judge.go`, one `suggested` line per suggested pair in the form of `_techspec.md` → Surface Transcript 1, a `skipped source-grouping <anchor> → <candidate>: <reason>` line per skipped grouping judgment, nothing for a clear one, and the summary `Judge: <a> advisory, <s> suggested, <c> clear, <k> skipped; …` of API Contract 1 in every form that counts judgments. Surface Transcripts 1 to 3 MUST be reproduced byte for byte.
2. MUST list each grouping judgment in the `--format json` document as API Contract 2 states, with no new top-level field.
3. MUST add to `specJudgeUsage` the sentence that the judge also asks whether each open Finding or Backlog Entry belongs with a source the Spec adopted, and that a suggestion never gates, as API Contract 5 states. The flags, exit codes and synopsis MUST NOT change.
4. MUST update the existing tests in `internal/cli/spec_judge_test.go` only where `_techspec.md` → Existing tests that change says: the summary of Spec 0205's Transcripts 1, 4 and 5 gains `0 suggested`. MUST add `TestSpecJudgePrintsAGroupingSuggestion`, `TestSpecJudgeSkipsANonEnglishSource`, `TestSpecJudgeCountsGroupingJudgmentsNotAsked`, `TestSpecJudgePrintsGroupingJSON` and `TestSpecJudgeHelpNamesTheGroupingQuestion` there, each with a temporary repository, a temporary Roundfix Home, a fake transport, a fixed clock and the key in the command's own environment.
5. MUST describe the grouping question in the `### spec judge` section of `docs/user-guide/commands/spec.md`: which pairs it asks, at every stage, the `suggested` line and its threshold of 0.3, that only Findings and Backlog Entries are sent, and that its recall is low.
6. MUST add to the `### Advisory judge` heading of `.agents/skills/roundfix/references/spec.md` that a `suggested` line is answered by adopting the open source within the grouping bound or by stating why it stays apart, and that a suggestion never gates. It MUST add no text inside the `### QA settlement` section of any skill.
7. MUST raise the Roundfix Skill's version in both front-matter fields of `.agents/skills/roundfix/SKILL.md` by one patch step from the value on this Task's starting commit, then run `make skills-sync`, record the version with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`, run `make baseline-digests`, and name in the Result every file those commands rewrote.
8. MUST prove each new gate can fail. The Result MUST record one sabotage of the summary (for example omitting the suggested count) and one of the suggestion line (for example printing the probability to one decimal), each with the test that failed, and that the code was restored.

Spec 0205 creates `internal/cli/spec_judge.go`, its test and `docs/user-guide/commands/spec.md` (`_prd.md` → Prerequisites). They are declared under `creates:` because they do not exist when this Spec is authored; this Task changes them and MUST NOT create one that Spec 0205 has not created. If one of them, or the `### Advisory judge` heading, does not exist, the Task stops and reports that Spec 0205 has not landed.

## Subtasks

- [ ] Print and count suggestions in the text form.
- [ ] List grouping judgments in the JSON form.
- [ ] Extend the help.
- [ ] Describe the question in the command reference and the skill, raise and record the version.
- [ ] Record one sabotage per gate in the Result.

## Acceptance Criteria

- [ ] Surface Transcripts 1, 2 and 3 are reproduced byte for byte.
- [ ] Spec 0205's transcripts still pass, with `0 suggested` in their summary.
- [ ] A grouping judgment in the JSON document has `kind` `source-grouping`, `line` null and the candidate as `target`.
- [ ] The command reference and the skill reference describe the grouping question, and the skill's raised version is recorded.

## Context

- creates: `internal/cli/spec_judge.go`
- creates: `internal/cli/spec_judge_test.go`
- creates: `docs/user-guide/commands/spec.md`
- interface: `.agents/skills/roundfix/references/spec.md`
- interface: `skills/roundfix/references/spec.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`

## Verification

- `out="$(go test -count=1 -v -run "^(TestSpecJudgePrintsAGroupingSuggestion|TestSpecJudgeSkipsANonEnglishSource|TestSpecJudgeCountsGroupingJudgmentsNotAsked|TestSpecJudgePrintsGroupingJSON|TestSpecJudgeHelpNamesTheGroupingQuestion|TestSpecJudgeReportsAdvisoryJudgments|TestSpecJudgeSkipsANonEnglishSpec|TestSpecJudgeStopsWhenTheServiceFails|TestSpecJudgePrintsJSON|TestSpecJudgeHelp|TestEveryOwnedSkillVersionIsRecorded)$" ./internal/cli ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestSpecJudgePrintsAGroupingSuggestion TestSpecJudgeSkipsANonEnglishSource TestSpecJudgeCountsGroupingJudgmentsNotAsked TestSpecJudgePrintsGroupingJSON TestSpecJudgeHelpNamesTheGroupingQuestion TestSpecJudgeReportsAdvisoryJudgments TestSpecJudgeSkipsANonEnglishSpec TestSpecJudgeStopsWhenTheServiceFails TestSpecJudgePrintsJSON TestSpecJudgeHelp TestEveryOwnedSkillVersionIsRecorded; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the five grouping tests do not exist, so the command fails.
- `for pair in "docs/user-guide/commands/spec.md|source-grouping" "docs/user-guide/commands/spec.md|P(same Spec)" "docs/user-guide/commands/spec.md|only Findings and Backlog Entries are sent" ".agents/skills/roundfix/references/spec.md|suggested" ".agents/skills/roundfix/references/spec.md|adopting the open source within the grouping bound or by stating why it stays apart"; do file="${pair%%|*}"; phrase="${pair#*|}"; test -f "$file" || { printf 'missing file: %s\n' "$file" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done; diff -r .agents/skills/roundfix skills/roundfix >/dev/null || { printf 'mirror differs: skills/roundfix\n' >&2; exit 1; }; make skills-sync-check` — expected: exit 0; before this Task neither file describes the grouping question, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 5; User Story 4; Core Feature 7; Success Metric 3; Declared breaks
- [_techspec.md](_techspec.md) — Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; API Contract 1; API Contract 2; API Contract 5; Existing tests that change; Testing Approach 4; Build Order 3
- ADR-0209; ADR-0200; ADR-0184; ADR-0187; ADR-0189
