---
task: task_03
spec: 0064-spec-artifact-consistency-gate
status: pending
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
