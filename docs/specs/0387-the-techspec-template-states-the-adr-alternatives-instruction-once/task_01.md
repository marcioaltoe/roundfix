---
task: task_01
spec: 0387-the-techspec-template-states-the-adr-alternatives-instruction-once
status: completed
type: docs
complexity: low
---

# Task 01: The TechSpec template states the ADR Alternatives instruction once, and owned skills refuse a phrase repeated across a line break

## Overview

The `## ADR Alternatives` comment of the `write-techspec` TechSpec template
ends one line with "Begin the relation with" and begins the next with the same
phrase. This Task deletes the repeated prefix in the canonical template,
syncs the shipped mirror, raises and records the skill's version with the
repository's record command, and adds a contract test that fails whenever a
line of any Roundfix-owned skill, canonical or shipped, ends with a phrase of
two or more words that the next line begins with.

## Requirements

1. MUST answer Argus's report of 2026-11-12 (ai-memory message `01a1d159`,
   [`references/2026-11-12-argus-the-techspec-template-repeats-begin-the-relation-with.md`](references/2026-11-12-argus-the-techspec-template-repeats-begin-the-relation-with.md),
   `_prd.md` → Sources 1; reported by argus) by implementing API Contracts 1
   to 4 of `_techspec.md`.
2. MUST, in `.agents/skills/write-techspec/references/techspec-template.md`,
   replace the comment line that begins "     Begin the relation with
   `not applicable`" with the line of `_techspec.md` → Interfaces, byte for
   byte: the same line without its leading "Begin the relation with ". Every
   other byte of the template stays.
3. MUST run `make skills-sync`, so `skills/write-techspec/` equals the
   canonical copy, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   which raises `metadata.version` and `version` in both `SKILL.md` copies
   and appends the version to `skills/testdata/owned-skill-versions.json`.
   Never type a version or a digest by hand.
4. MUST add `skills/repeated_phrase_test.go` in `package skills`, with the
   header lines `// Suite: owned skill repeated phrase contract`, `//
   Invariant: no line of an owned skill ends with a phrase of two or more
   words that the next line begins with` and `// Boundary: canonical and
   shipped copies of every owned skill file; a line that is only the phrase
   and Markdown table rows stay exempt`, holding:
   - the type `repeatedLinePhrase` and the function `repeatedLinePhrases` of
     `_techspec.md` → Interfaces, behaving as API Contract 3 (fields split on
     whitespace, compared after `strings.ToLower`, the longest equal run, a
     minimum of two words, no report when the run is the whole line or when
     either line's first word starts with `|`), with `Line` the 1-based
     number of the line that ends with the run and `Phrase` the next line's
     original words of the run joined by one space;
   - `TestRepeatedLinePhrasesFindsAPhraseRepeatedAcrossALineBreak` with the
     six cases of `_techspec.md` → Testing Approach, by those names and
     inputs;
   - `TestNoOwnedSkillRepeatsAPhraseAcrossALineBreak` with subtests
     `canonical` (`../.agents/skills/<name>`) and `shipped` (`<name>` under
     the package directory), walking every file of every name `Names()`
     returns and reporting each hit with `t.Errorf` as
     `<path>:<line> ends with "<phrase>" and the next line begins with it;
     write the phrase once`, the path written with forward slashes.

   Each top-level test and subtest calls `t.Parallel()` first. Helper names
   must not collide with names the package already declares.
5. MUST NOT change any Go production file, any file under `internal/`,
   `docs/agents/`, `docs/user-guide/` or `docs/adr/`, `CONTEXT.md`, any other
   skill or template line, the Makefile, the CI workflows or the `### QA
   settlement` section of any skill.

## Subtasks

- [x] Delete the repeated prefix in the canonical template.
- [x] Sync the mirror and record the raised version.
- [x] Add the detector, its unit cases and the owned-skill guard.

## Acceptance Criteria

- [x] Both template copies state "Begin the relation with" once in the ADR Alternatives comment and keep both tokens.
- [x] `write-techspec` carries its next version, recorded, with the mirror equal to the canonical copy.
- [x] The guard passes on every owned skill, canonical and shipped, and each detector case passes.

## Context

- instruction: `skills/skills.go`
- instruction: `skills/owned_skill_versions_test.go`
- instruction: `skills/baseline_skill_contract_test.go`
- instruction: `internal/speccheck/citations.go`
- interface: `.agents/skills/write-techspec/references/techspec-template.md`
- interface: `skills/write-techspec/references/techspec-template.md`
- interface: `.agents/skills/write-techspec/SKILL.md`
- interface: `skills/write-techspec/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- creates: `skills/repeated_phrase_test.go`

## Consumers

- the shipped `write-techspec` skill folder and the owned-skill version record — consumers: `./skills`

## Verification

- `out="$(go test -timeout 10m -count=1 -v ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestRepeatedLinePhrasesFindsAPhraseRepeatedAcrossALineBreak/doubled-instruction TestRepeatedLinePhrasesFindsAPhraseRepeatedAcrossALineBreak/two-words-in-another-case TestRepeatedLinePhrasesFindsAPhraseRepeatedAcrossALineBreak/one-repeated-word TestRepeatedLinePhrasesFindsAPhraseRepeatedAcrossALineBreak/command-listing TestRepeatedLinePhrasesFindsAPhraseRepeatedAcrossALineBreak/table-rows TestRepeatedLinePhrasesFindsAPhraseRepeatedAcrossALineBreak/single-instruction TestNoOwnedSkillRepeatsAPhraseAcrossALineBreak/canonical TestNoOwnedSkillRepeatsAPhraseAcrossALineBreak/shipped TestEveryOwnedSkillVersionIsRecorded TestAuthorialSkillSync/write-techspec; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; go test -timeout 10m -count=1 -run '^(TestAPortugueseNotApplicableADRDrawsInNoCitingADR|TestAMarkerOutsideTheClosedListStillDrawsIn|TestANotApplicableADRDrawsInNoCitingADR|TestAnAppliedADRStillDrawsInItsCitingADR)$' ./internal/speccheck` — expected: exit 0. Before this Task the detector tests do not exist, so the command fails at its first pass check. Afterwards every detector case and both owned-skill subtests pass, the version record and the `write-techspec` mirror pass, the whole `./skills` package passes, and the `SC-ADR-RELATED` marker and seed tests pass unchanged.
- `tr -s '[:space:]' ' ' < .agents/skills/write-techspec/references/techspec-template.md | grep -qE -- 'rejected[.] Begin the relation with .not applicable., or with .{1,4}o se aplica. in a Spec written in Portuguese, for an ADR the Spec lists only to record why it does not apply[.]' || { printf 'missing corrected comment in the canonical template\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < skills/write-techspec/references/techspec-template.md | grep -qE -- 'rejected[.] Begin the relation with .not applicable., or with .{1,4}o se aplica. in a Spec written in Portuguese' || { printf 'missing corrected comment in the shipped template\n' >&2; exit 1; }; ! tr -s '[:space:]' ' ' < .agents/skills/write-techspec/references/techspec-template.md | grep -qiF -- 'Begin the relation with Begin' && ! tr -s '[:space:]' ' ' < skills/write-techspec/references/techspec-template.md | grep -qiF -- 'Begin the relation with Begin' && cmp -s .agents/skills/write-techspec/references/techspec-template.md skills/write-techspec/references/techspec-template.md` — expected: exit 0. Before this Task the sentence that "rejected." opens continues with a second "Begin the relation with", so both pattern checks fail. Afterwards both copies read the instruction once with both tokens, and they are identical.

## References

- `_prd.md` → Sources; Goals; Core Features 1-3; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4
- `_techspec.md` → Current behavior; Interfaces; API Contracts 1-4; Invariants 1-3; Sabotages; Testing Approach; Build Order 1

## Result

Skills loaded: implement-task, coding-guidelines, golang-testing, testing-boss, tech-writer, evidence-gate

- Removed only the repeated `Begin the relation with` prefix from the canonical TechSpec template comment, synchronized the shipped mirror with `make skills-sync`, and confirmed both copies are byte-identical with one occurrence of the phrase while retaining both `not applicable` and `não se aplica` tokens.
- Ran `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` successfully. The recorder raised both `write-techspec` copies to `0.0.20` and appended the generated digest record.
- Added `repeatedLinePhrase`, `repeatedLinePhrases`, the six named detector cases, and the canonical/shipped owned-skill guard. Focused post-change check `go test ./skills -run '^(TestRepeatedLinePhrasesFindsAPhraseRepeatedAcrossALineBreak|TestNoOwnedSkillRepeatsAPhraseAcrossALineBreak|TestEveryOwnedSkillVersionIsRecorded|TestAuthorialSkillSync/write-techspec)$' -count=1 -v` passed with 11 tests.
- The Task's declared `## Verification` commands were not run; the Daemon owns that verification and settlement.
