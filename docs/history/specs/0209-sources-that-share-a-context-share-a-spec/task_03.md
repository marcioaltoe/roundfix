---
task: task_03
spec: 0209-sources-that-share-a-context-share-a-spec
status: completed
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

## Result

Implemented the command output, help, tests, and documentation for this
Task's slice. Task status and authored Verification remain Daemon-owned;
neither declared Verification command was run. No commit, push, or Pull
Request was made.

### Implementation and acceptance evidence

- Surface Transcripts 1–3: `TestSpecJudgePrintsAGroupingSuggestion/prd`,
  `TestSpecJudgeSkipsANonEnglishSource`, and
  `TestSpecJudgeCountsGroupingJudgmentsNotAsked` compare stdout, stderr,
  and exit code byte for byte against the authored transcripts. The
  suggestion test also checks `techspec` and the default both-artifacts
  stage. Each fixture uses a temporary repository, temporary Roundfix
  Home, fake HTTP transport, fixed UTC clock, and synthetic credentials
  exclusively in the command environment. Missing command credentials
  produce zero requests.
- Spec 0205 compatibility: the existing advisory, non-English Spec, and
  service-stop transcripts pass with `0 suggested` in their summaries.
  The prerequisite also shipped `TestSpecJudgePrintsIndividualSkip`, which
  has a counting summary beyond the three named transcripts; its expected
  summary gained only `0 suggested` to satisfy Requirement 1's contract for
  every counting form. No other existing assertions changed.
- JSON: `TestSpecJudgePrintsGroupingJSON` exercises public command dispatch
  and checks both suggested and clear pairs, the unchanged schema and exact
  top-level field set, `kind: source-grouping`, anchor artifact, candidate
  target, null line/text/section/choice fields, probabilities, and versioned
  model. Task 02's existing `Judgment.MarshalJSON` already provides these
  fields, so no judge-package change was needed.
- Skipped grouping judgments: `TestSpecJudgePrintsSkippedGroupingPair`
  injects an unpinned model and checks both arrow-separated pair lines,
  reasons, and the zero-suggestion summary. Clear pairs print no text line.
- Help: `TestSpecJudgeHelpNamesTheGroupingQuestion` checks the actual help
  output names the open-source question and says a suggestion never gates.
  Existing help tests still pass; flags, synopsis, and exit-code declarations
  are unchanged.
- Documentation and skill version: the command reference now describes
  pairs at every stage, the `source-grouping` line, `P(same Spec)` threshold
  of 0.3, the Findings/Backlog request boundary, and low recall. The Roundfix
  Advisory judge reference describes adopting the open source within the
  grouping bound or stating why it stays apart, and that suggestions never
  gate. Both skill version fields rose from starting-commit `0.1.8` to
  `0.1.9`, recorded in the owned-version history. A focused Python inspection
  confirmed both changed mirror files match their canonical counterparts
  and both `### QA settlement` sections remain byte-identical to `HEAD`.

### Focused checks and regeneration

Go and Make checks used `GOCACHE=/private/tmp/roundfix-task03-gocache`.

- Before implementation, `go test ./internal/cli -run
  '^TestSpecJudge(PrintsAGroupingSuggestion|SkipsANonEnglishSource|HelpNamesTheGroupingQuestion|PrintsSkippedGroupingPair)$'
  -count=1` exited 1: missing suggestion/count/help and incorrect skipped
  pair formatting were exposed.
- After implementation and restoring both sabotages,
  `go test ./internal/cli -run '^TestSpecJudge' -count=1 -v` exited 0;
  all existing and new Spec Judge tests passed.
- `make skills-sync` exited 0 and rewrote
  `skills/roundfix/SKILL.md` and `skills/roundfix/references/spec.md`.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
  -record-skill-versions` exited 0 and rewrote only
  `skills/testdata/owned-skill-versions.json`, adding Roundfix `0.1.9`.
- `make baseline-digests` exited 0 with `changed: false`; it rewrote no file.
- `make verify-incremental` first exited 2 because the sandbox denied
  process-table access in
  `TestRunForceStopLegacyRunWithoutOwnerIdentityStillStopsOwner` and
  `TestRunForceStopOwnerProcessIntegrationProvesExitBeforeStoreCompletion`.
  The same command rerun with host process-table permission exited 0:
  formatting, vet, repository tests, skill synchronization/checks, and build
  passed. No test or implementation was altered to bypass that restriction.
- `git -c core.fsmonitor=false diff --check` exited 0. The changed-file
  inspection contains only this Task's CLI, test, documentation, canonical
  Roundfix skill, required mirrors/version record, and this Task's Result.

### Sabotage evidence

Each mutation ran `go test ./internal/cli -run
'^TestSpecJudgePrintsAGroupingSuggestion/prd$' -count=1` against the mutated
renderer, then restored the original source in a `finally` block.

- Summary sabotage: forced the rendered suggestion count to zero. The
  command exited 1 and
  `TestSpecJudgePrintsAGroupingSuggestion/prd` failed on `0 suggested`
  versus the expected `1 suggested`. The code was restored.
- Suggestion-line sabotage: changed `P(same Spec) %.2f` to `%.1f`. The
  command exited 1 and the same test failed on `P(same Spec) 0.8` versus
  `P(same Spec) 0.81`. The code was restored.

The subsequent focused Spec Judge suite and incremental check passed on the
restored implementation. Final Task Verification and settlement are pending
the Daemon.

## Carry-forward provenance

- Source Run: `run_20261001T233647Z_36941fe094803f4f`
- Source commit: `5fcc743d854201ca3cc001b42720dc98d7ad60b8`
