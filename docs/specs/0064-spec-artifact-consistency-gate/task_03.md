---
task: task_03
spec: 0064-spec-artifact-consistency-gate
status: completed
type: backend
complexity: medium
---

# Task 03: Report undocumented emitted vocabulary

## Overview

Deliver the vocabulary-contract detector: when a Spec declares that a source
emits user-facing tokens and names the document that documents them, the check
proves every emitted token appears in that document. This is the shape of Spec
0058's QA-004 — a workflow emitting five failure prefixes while its runbook
documented four — reduced to a declaration the check can settle in
milliseconds.

The declaration is a new authoring convention and this Spec ships no skill to
teach it, so the finding's fix text is the only teacher: it must carry the
block's exact shape.

## Requirements

1. MUST parse an optional vocabulary contract declared in the TechSpec, each
   entry naming an emitting path, an RE2 pattern selecting the tokens, and a
   documenting path.
2. MUST report `SC-VOCABULARY-UNDOCUMENTED` at severity `error` for each
   distinct token the pattern selects from the emitting file that does not
   appear in the documenting file, locating the emitting line and the
   documenting file.
3. MUST report an unreadable emitting path, unreadable documenting path, or
   invalid pattern as a finding naming the declaration line, never as a panic
   or a silent pass.
4. MUST skip the detector and record the skip when the Spec declares no
   vocabulary contract.
5. MUST carry the declaration block's exact shape in the finding's fix text, so
   an author who has never seen the convention can write one from the
   diagnostic alone.
6. SHOULD deduplicate repeated tokens so one undocumented token reports once
   however many times it is emitted.

## Subtasks

- [ ] Parse the vocabulary contract block into its entries.
- [ ] Extract distinct tokens from each emitting file by its pattern.
- [ ] Implement the detector, including the unreadable-path and
      invalid-pattern findings.
- [ ] Write the fix text carrying the block's shape.
- [ ] Add fixtures: a satisfied contract, a contract missing one token, an
      invalid pattern, a missing emitting file, and a Spec with no contract.

## Acceptance Criteria

- [ ] A fixture whose emitting file yields five tokens and whose documenting
      file carries four reports `SC-VOCABULARY-UNDOCUMENTED` naming exactly the
      fifth token.
- [ ] A fixture whose documenting file carries every token produces no finding.
- [ ] A token emitted three times and documented nowhere reports once.
- [ ] An invalid RE2 pattern reports a finding located at the declaration line
      and the check still returns a result.
- [ ] A Spec declaring no vocabulary contract produces no finding and lists the
      detector as skipped.
- [ ] The reported fix text contains the declaration block's three field names,
      asserted by a test.

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/speccheck -count=1 -run 'Vocabulary' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the vocabulary tests ran and passed.
- `go test ./internal/speccheck -count=1` — expected: exit 0.

## References

- `_prd.md` → Core Feature 6; Goals.
- `_techspec.md` → Data Models; API Contracts; Build Order 4; Risks &
  Considerations.
- ADR-0094.

## Result

Implemented the optional `## Vocabulary Contract` parser and
`SC-VOCABULARY-UNDOCUMENTED` detector. Each list entry declares `emits`,
`pattern`, and `documented-in`; the detector compiles the pattern with Go's
RE2 engine, extracts full matches with their first emitting line, deduplicates
repeated tokens, compares them literally with the documenting file, and emits
located `error` findings for undocumented tokens, invalid patterns, and
unreadable paths. Specs without the section record the detector as skipped.
Every finding's fix includes the copyable declaration block.

Focused checks:

- Pre-change signal: `GOCACHE="$PWD/.gocache" GOFLAGS=-buildvcs=false rtk go
  test ./internal/speccheck -run '^TestCheckVocabulary'` exited 1 because the
  diagnostic constant and detector did not exist.
- Formatting: `rtk gofmt -w` over the changed Go implementation, test, and
  fixture source files exited 0.
- Post-change: `GOCACHE="$PWD/.gocache" GOFLAGS=-buildvcs=false rtk go test
  ./internal/speccheck -run '^(TestCheckVocabulary|TestCheckCleanFixture)'`
  exited 0 and reported 10 passing tests.
- `rtk git diff --check` exited 0; a trailing-whitespace scan over the new
  implementation, test, and fixture trees found no matches.
- The commands under `## Verification` were not run; the Daemon owns them.

Acceptance evidence:

1. `TestCheckVocabularyUndocumented` exercises five emitted tokens against
   four documented tokens and asserts one finding naming `publish:` plus both
   file locations.
2. `TestCheckVocabularySatisfied` asserts that documentation containing all
   five selected tokens produces no vocabulary finding.
3. `TestCheckVocabularyDeduplicatesRepeatedToken` emits `publish:` three times
   and asserts one finding.
4. `TestCheckVocabularyInvalidPatternReturnsFinding` proves an invalid RE2
   returns a result with a finding at the pattern declaration line.
5. `TestCheckVocabularyAbsentContractIsSkipped` proves an absent contract
   produces no finding and records the detector skip.
6. `TestCheckVocabularyFixTeachesDeclarationShape` asserts the exact copyable
   block and its `emits`, `pattern`, and `documented-in` field names.

Additional requirement evidence:

- `TestCheckVocabularyUnreadablePathsReturnFindings` covers both missing
  emitting and documenting files as located findings.
- The missing-token fixture carries two entries, so the result also exercises
  repeated contract parsing rather than only the first entry.
