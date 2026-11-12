---
spec: 0387-the-techspec-template-states-the-adr-alternatives-instruction-once
prd: _prd.md
created: 2026-11-12
---

# The TechSpec template states the ADR Alternatives instruction once — Technical Spec

## Executive Summary

The defect is one line of skill text. Spec 0379 rewrote the second line of the
`## ADR Alternatives` comment in the `write-techspec` TechSpec template by
prefixing it with "Begin the relation with", the phrase the first line already
ends with. This Spec deletes that prefix in the canonical template, copies it
to the shipped mirror with `make skills-sync`, raises and records the skill's
version with the repository's record command, and adds a guard: a pure
detector of a phrase of two or more words that ends one line and begins the
next, with unit cases, and a contract test that runs it over every file of
the twelve Roundfix-owned skills, canonical and shipped. No Go production
code, Baseline module, guide or glossary changes.

The trade-off accepted: the guard skips a single repeated word and a line made
only of the repeated phrase, so a repeated "the" across a line break or a
whole duplicated line still passes. In exchange it needs no exception list:
at the authoring base it reports the Argus line and nothing else across every
file of the twelve owned skills, canonical and shipped, and it does not flag
the command listings in the Roundfix skill references.

## Project Constraints

- Identifier strategy: not applicable — no identifier, record or schema is
  minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or
  network call is added. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0389 and ADR-0377 apply unchanged,
  as `_prd.md` → Project Constraints states: the corrected comment keeps
  ADR-0389's two tokens, and `SC-ADR-RELATED` keeps reading them from the
  Spec's own TechSpec. ADR-0093 and ADR-0168, which cite ADR-0377, do not
  apply: this Spec changes no check. Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  answer of 2026-11-12 quoted in the authorization record; the Governed Paths
  are the `write-techspec` skill's `SKILL.md` pair and both copies of its
  TechSpec template. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0387-the-techspec-template-states-the-adr-alternatives-instruction-once/_authorization.md`;
  bounded files: `.agents/skills/write-techspec/SKILL.md`, `.agents/skills/write-techspec/references/techspec-template.md`, `skills/write-techspec/SKILL.md`, `skills/write-techspec/references/techspec-template.md`.

## Current behavior

Read on 2026-11-12 at `805488c2`.

1. `.agents/skills/write-techspec/references/techspec-template.md` and its
   mirror `skills/write-techspec/references/techspec-template.md` are
   byte-identical. Their lines 86 to 88 are the comment of `## ADR
   Alternatives`; line 87 ends "rejected. Begin the relation with" and line
   88 begins "     Begin the relation with `not applicable`, or with `não se
   aplica` in a Spec written in Portuguese, for an ADR the Spec lists only to
   record why it does not apply. -->". Before `c11c5883` line 88 read
   "     `not applicable` for an ADR the Spec lists only to record why it does
   not apply. -->".
2. `write-techspec` is at 0.0.19 in `metadata.version` and `version` of both
   `SKILL.md` copies; `skills/testdata/owned-skill-versions.json` records it.
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
   -record-skill-versions` raises both fields to the next patch version and
   records the digest when the folder content changed.
3. `Names()` in `skills/skills.go` lists the twelve owned skills, and the
   binary embeds their shipped copies under `skills/<name>/`;
   `TestAuthorialSkillSync` holds the mirror equal to the canonical copy.
   No test reads a repeated phrase.
4. `SC-ADR-RELATED` reads the token with `notApplicableRelationPattern`
   (`internal/speccheck/authoring_proof.go`) in `appliedADRNumbers`
   (`internal/speccheck/citations.go`) over the checked Spec's TechSpec. The
   template is named in Go only by `internal/speccheck/governed.go` (a
   Governed Path) and read by tests that assert its Project Constraints,
   bounded-files label and section headings, none of them the comment.
5. A scan of the owned skills for a phrase repeated across a line break, with
   the result recorded in `_prd.md` → Validation: one hit per copy, the
   template line; 14 lines per copy made only of the phrase that a command
   listing extends; table rows.

## System Architecture

- `.agents/skills/write-techspec/references/techspec-template.md` and its
  mirror: one comment line.
- `.agents/skills/write-techspec/SKILL.md`, `skills/write-techspec/SKILL.md`
  and `skills/testdata/owned-skill-versions.json`: the version raise and its
  record, written by the record command.
- `skills/repeated_phrase_test.go` (new, `package skills`): the detector, its
  unit test and the contract test over the owned skills.

## Implementation Design

### Interfaces

The corrected comment, lines 86 to 88 of both copies:

```markdown
<!-- One line for each ADR listed under Active ADR obligations. State how the design
     relates to the alternatives that ADR rejected. Begin the relation with
     `not applicable`, or with `não se aplica` in a Spec written in Portuguese, for an ADR the Spec lists only to record why it does not apply. -->
```

The detector, in the new test file:

```go
type repeatedLinePhrase struct {
	Line   int    // 1-based number of the line that ends with the phrase
	Phrase string // the phrase as the next line writes it
}

// repeatedLinePhrases compares the trailing words of each line with the
// leading words of the next, split on whitespace and lowered, and reports the
// longest equal run of two or more words. A line made only of that run, and a
// Markdown table row (first word starts with "|"), are not reported.
func repeatedLinePhrases(text string) []repeatedLinePhrase
```

### Data Models

None. The version record keeps its schema; one `write-techspec` entry is
appended by the record command.

### API Contracts

1. API Contract: the template comment. In both copies of the template, line
   88 is the line of Interfaces, every other byte equals the Delivery Base,
   the two copies are byte-identical, and the whitespace-normalized file
   contains "rejected. Begin the relation with `not applicable`, or with `não
   se aplica` in a Spec written in Portuguese, for an ADR the Spec lists only
   to record why it does not apply." and not "Begin the relation with Begin
   the relation with".
2. API Contract: the skill version. Both `SKILL.md` copies of
   `write-techspec` carry the next version written by the record command, and
   `skills/testdata/owned-skill-versions.json` records it with the folder
   digest; no version or digest is typed by hand.
3. API Contract: the detector. `repeatedLinePhrases` returns, for each line
   that ends with a run of two or more words the next line begins with, the
   line number and the next line's words of that run; words compare without
   regard to case; it returns nothing for one repeated word, for a line made
   only of the run, and when either line is a Markdown table row.
4. API Contract: the owned-skill guard.
   `TestNoOwnedSkillRepeatsAPhraseAcrossALineBreak` walks every file of every
   name in `Names()` under `../.agents/skills/` (subtest `canonical`) and
   under the package directory (subtest `shipped`), and reports each hit as
   `<path>:<line> ends with "<phrase>" and the next line begins with it; write
   the phrase once`.

### Surface Transcripts

None. No command, flag, exit code, message or stream placement changes; one
line of skill text, its version and a test change.

### Invariants

```text
1. No line of a Roundfix-owned skill, canonical or shipped, ends with a phrase of two or more words that the next line begins with.
2. The ADR Alternatives comment names both not-applying tokens of ADR-0389, in that order, once.
3. SC-ADR-RELATED reads the not-applying token only from the checked Spec's own TechSpec.
```

### Sabotages

Each runs on a copy of the audited head and names the test that observes the
mutated unit; the prototype at `805488c2` failed each as listed.

1. Put "Begin the relation with " back at the start of line 88 of the
   canonical template only:
   `TestNoOwnedSkillRepeatsAPhraseAcrossALineBreak/canonical` fails.
2. The same in the shipped mirror only:
   `TestNoOwnedSkillRepeatsAPhraseAcrossALineBreak/shipped` fails.
3. Make the detector report nothing (its reporting condition always skips):
   `TestRepeatedLinePhrasesFindsAPhraseRepeatedAcrossALineBreak/doubled-instruction`
   fails.
4. Drop the exemption of a line made only of the run:
   `TestRepeatedLinePhrasesFindsAPhraseRepeatedAcrossALineBreak/command-listing`
   fails (and both owned-skill subtests, on the command listings).
5. Lower the minimum run to one word:
   `TestRepeatedLinePhrasesFindsAPhraseRepeatedAcrossALineBreak/one-repeated-word`
   fails.
6. Drop the table-row exemption:
   `TestRepeatedLinePhrasesFindsAPhraseRepeatedAcrossALineBreak/table-rows`
   fails.
7. Compare words with their case:
   `TestRepeatedLinePhrasesFindsAPhraseRepeatedAcrossALineBreak/two-words-in-another-case`
   fails.

## Coverage Map

Deviations from the PRD: None.

- Goal 1, Core Feature 1, Success Metric 1 → API Contract 1, Invariant 2, task_01
- Goal 2, Core Feature 2, Success Metric 3 → API Contract 2, task_01
- Goal 3, Core Feature 3, Success Metric 2 → API Contracts 3 and 4, Invariant 1, task_01
- Goal 4, Success Metric 4 → Invariant 3, task_01 (unchanged speccheck tests) and task_02, which also records the PRD's Unreachable Acceptance as one declared row

## Absence Evidence

1. The repeated-phrase detector and its owned-skill guard — command: `grep -rn -i -e "repeated" -e "duplicate word" -e "lexical" skills/*.go internal/docscontract/*.go`; result: at `805488c2` the matches are a variable of the version-raise test (`owned_skill_version_raise_test.go`, a repeated recording) and the glossary contract's fixture of one sentence repeated across two glossary entries (`user_guide_contract_test.go`); none reads skill text, and no test compares the end of one line with the start of the next.

## ADR Alternatives

- ADR-0389: the design keeps its closed list and takes none of its rejected options; the template comment still prescribes `not applicable` first and accepts `não se aplica`, with no synonym and no configuration key.
- ADR-0377: the design keeps reading the token only in the TechSpec's own `## ADR Alternatives` line, never in the PRD row or a template, as the decision requires.
- ADR-0093: not applicable — it decides that Spec consistency is checked by citation, and this Spec changes no check; it is listed only because it cites ADR-0377.
- ADR-0168: not applicable — it decides when a related-ADR gap opens, and this Spec changes no check; it is listed only because it cites ADR-0377.

## Integration Points

None. No external system, network or provider is reached; the test reads the
repository's own skill files.

## Testing Approach

- `skills/repeated_phrase_test.go` (new), `package skills`, with a
  `Suite:`/`Invariant:`/`Boundary` header like the package's other contract
  files, each top-level test and subtest calling `t.Parallel()` first:
  - `TestRepeatedLinePhrasesFindsAPhraseRepeatedAcrossALineBreak`, table
    cases: `doubled-instruction` (the Argus lines, reported at line 1 as
    "Begin the relation with"), `two-words-in-another-case` ("Record the
    result in the Task" then "the task file keeps it.", reported as "the
    task"), `one-repeated-word` ("The guide says that" then "that rule still
    holds.", nothing), `command-listing` ("roundfix gc compact" then
    "roundfix gc compact --apply", nothing), `table-rows` ("| id | title | |"
    then "| | | done |", nothing), `single-instruction` (the corrected
    lines, nothing).
  - `TestNoOwnedSkillRepeatsAPhraseAcrossALineBreak`, subtests `canonical`
    and `shipped`, as API Contract 4.
- Unchanged and run: `go test ./skills` (mirror, version record, version
  raise and owned-skill contracts); the `SC-ADR-RELATED` tests
  `TestAPortugueseNotApplicableADRDrawsInNoCitingADR`,
  `TestAMarkerOutsideTheClosedListStillDrawsIn`,
  `TestANotApplicableADRDrawsInNoCitingADR` and
  `TestAnAppliedADRStillDrawsInItsCitingADR` in `./internal/speccheck`; the
  documentation contracts with `-tags docscontract`.
- The file compiles for Windows (`GOOS=windows go vet ./skills`): it uses no
  platform call.

## Build Order

1. The template line, its sync, the version record and the guard test
   (task_01; no dependency).
2. QA gate (task_02) (depends on: 1).

No Spec is listed in `requires`; Specs 0383, 0385, 0384 and 0386 deliver
first and touch only the `roundfix` entries of the shared version record.

## Risks & Considerations

- A future owned-skill edit that legitimately ends a line with the words the
  next line starts with fails the guard; the author rewraps the line. The
  failure names the file, the line and the phrase.
- Adopters keep the stutter until their next `roundfix baseline update`
  after the release; Argus said it will not edit its copy.
- The record command chooses the next version from the tree it runs on; if
  another Spec raised `write-techspec` first, it raises past that one.

## Glossary

None.

## Decisions

- Delete the duplicated prefix; keep every other word of the comment.
- A general guard over all owned skills, in `package skills`, reading both
  copies; phrases of two or more words, case-insensitive, with the two
  exemptions the characterization needs. No ADR: no architecture or Baseline
  decision changes.
