---
task: task_01
spec: 0191-claims-with-receipts-and-contracts-as-they-ship
status: completed
type: backend
complexity: high
---

# Task 01: A Claim Receipt is proved, and a held attribution without one is a gap

## Overview

`SC-CITATION-UNSUPPORTED` compares a claim's words with the words of the decision record it names. It cannot tell whether the record makes the claim. This Task adds the Claim Receipt of ADR-0183: a source, a colon and a verbatim quote, which the Spec Consistency Check proves against the source file. It also adds the contract horizon that decides which Specs are held to the missing-receipt gap. The Task starts by recording what the parser returns today, so that the change can be shown to move nothing else.

## Requirements

1. MUST first record the characterization of `_techspec.md` → Testing Approach 1, before changing any parser: byte-identical fixture copies of the four archived artifacts named in Context, and a golden file of every claim `CitationClaims` returns for them (artifact, line, target and subject). The characterization test MUST pass on the unchanged parser and after this Task, and the Result MUST state that it was recorded before the parser changed.
2. MUST add `internal/speccheck/authoring_horizon.go` with `ConcreteContractGuidePath` and `newContractHorizon`, deciding in the order `_techspec.md` → The contract horizon states. The PRD's adding commit MUST be resolved through one helper shared with `newADRHorizon`, and `newADRHorizon`'s behavior MUST NOT change.
3. MUST add `internal/speccheck/receipts.go` with `Receipt`, `ReceiptedClaim`, `Receipts`, `ReceiptedClaims` and `ProveReceipt`, as `_techspec.md` → Interfaces and Claim Receipts state: the two source forms, the quote that may wrap across one line break, the shared paragraph walker, the resolution of a decision record and of a path (refusing a symbolic link at any path component, checked with `lstat` from `repoRoot`), the whitespace normalization and the three-word minimum. `TestAReceiptThroughASymbolicLinkDoesNotResolve` MUST prove that a receipt whose path crosses a symbolic link, to a file inside or outside the repository, is `SC-RECEIPT-UNPROVEN` with the summary ending `does not resolve to a file`. It MUST uphold `_techspec.md` → Invariants 1 to 7.
4. MUST declare `SC-RECEIPT-UNPROVEN` and `SC-RECEIPT-MISSING` in the code block of `internal/speccheck/citations.go`, add both to `citationCoverageDetectorCodes`, register both at the `prd` stage in `internal/speccheck/coherence.go`, and call the detector from `detectCitationCoverageAndReferences` and from the stage-scoped check. `internal/speccheck/constraints.go` MUST NOT change.
5. MUST report `SC-RECEIPT-UNPROVEN` as an error for every written receipt that is not proven, whether or not the Spec is held, and `SC-RECEIPT-MISSING` as a gap only for a held Spec, with the summaries, locations and fix texts of `_techspec.md` → Surface Transcripts 1 and 2. A Spec that is not held MUST list `SC-RECEIPT-MISSING` as skipped with the reason the horizon gives.
6. MUST keep `CitationClaims`, `SC-CITATION-UNSUPPORTED` and every existing test unchanged, rename or remove no top-level test, and change no exported function signature.
7. MUST add both codes to `corpusFindingCodes` in `internal/docscontract/corpus_test.go`, to `internal/docscontract/testdata/corpus-golden.json` at `0` with an `update` sentence that names this change, and to the golden pinned in `internal/spec/archive_layout_characterization_test.go`. Every earlier count MUST stay unchanged.
8. MUST add both codes, each with one line saying when it is raised, to the identifier list of the Spec Consistency Check section of `.agents/skills/roundfix/SKILL.md`, then run `make skills-sync`.
9. MUST put the new tests in `internal/speccheck/receipt_characterization_test.go`, `internal/speccheck/authoring_horizon_test.go` and `internal/speccheck/receipts_test.go`. The horizon tests MUST use `gittest` repositories or plain temporary directories, never this repository's own history. Each code MUST have one test where the defect is planted and the finding is reported, and one where it is absent and nothing is reported.

## Subtasks

- [ ] Record the claim characterization before changing the parser.
- [ ] Add the contract horizon and share the PRD-commit helper.
- [ ] Parse receipts, prove them and pair them with claims.
- [ ] Report the two findings and the skip at their stage.
- [ ] Add the codes to the corpus golden, its pin and the Roundfix skill.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] The claims recorded for the four fixture artifacts are identical before and after the change.
- [ ] A verbatim quote is proven, including one that wraps across lines in the artifact and in the source. One changed word, an unresolved source and a two-word quote each report `SC-RECEIPT-UNPROVEN`.
- [ ] A receipt whose source is a backticked path is proved against that file, and a receipt inside a fenced block is not read.
- [ ] In a held Spec, an attribution with no receipt for its record in the same paragraph reports `SC-RECEIPT-MISSING` as a gap, and a receipt for another record does not cover it.
- [ ] A PRD committed before the guide, a repository without the guide and one whose guide is uncommitted report no `SC-RECEIPT-MISSING` and list it as skipped. A PRD committed in or after the guide's commit, an uncommitted PRD and an unreadable history are held.
- [ ] A Spec that is not held still reports an unproven receipt.
- [ ] The rendered text of the two findings holds the finding lines of Surface Transcripts 1 and 2.
- [ ] The corpus golden holds both codes at `0`, every earlier count is unchanged, and the Roundfix skill and its mirror name both codes.

## Context

- instruction: `docs/adr/0183-an-attributed-claim-carries-a-receipt-the-check-proves.md`
- instruction: `docs/adr/0116-a-citation-is-checked-against-what-it-cites.md`
- instruction: `docs/adr/0168-a-related-adr-gap-opens-only-for-adrs-that-predate-the-spec.md`
- instruction: `docs/history/specs/0181-gates-that-refuse-only-what-someone-can-act-on/_prd.md`
- instruction: `docs/history/specs/0181-gates-that-refuse-only-what-someone-can-act-on/_techspec.md`
- instruction: `docs/history/specs/0182-delivery-that-reviews-and-retries-from-where-the-item-stands/_prd.md`
- instruction: `docs/history/specs/0182-delivery-that-reviews-and-retries-from-where-the-item-stands/_techspec.md`
- interface: `internal/speccheck/citations.go`
- interface: `internal/speccheck/coherence.go`
- interface: `internal/speccheck/adr_horizon.go`
- creates: `internal/speccheck/authoring_horizon.go`
- creates: `internal/speccheck/authoring_horizon_test.go`
- creates: `internal/speccheck/receipts.go`
- creates: `internal/speccheck/receipts_test.go`
- creates: `internal/speccheck/receipt_characterization_test.go`
- creates: `internal/speccheck/testdata/receipt-characterization/0181-prd.md`
- creates: `internal/speccheck/testdata/receipt-characterization/0181-techspec.md`
- creates: `internal/speccheck/testdata/receipt-characterization/0182-prd.md`
- creates: `internal/speccheck/testdata/receipt-characterization/0182-techspec.md`
- creates: `internal/speccheck/testdata/receipt-characterization/claims-golden.json`
- interface: `internal/docscontract/corpus_test.go`
- interface: `internal/docscontract/testdata/corpus-golden.json`
- interface: `internal/spec/archive_layout_characterization_test.go`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestReceiptCharacterizationKeepsEveryParsedClaim|TestASpecCommittedBeforeTheGuideIsNotHeld|TestASpecCommittedWithOrAfterTheGuideIsHeld|TestAnUncommittedSpecIsHeldOnceTheGuideIsCommitted|TestARepositoryWithoutTheGuideHoldsNoSpec|TestAnUncommittedGuideHoldsNoSpec|TestAGuideWithUnreadableHistoryHoldsEverySpec|TestAVerbatimReceiptIsProven|TestAReceiptThatWrapsAcrossLinesIsProven|TestAReceiptWithOneChangedWordIsUnproven|TestAReceiptWithAnUnresolvedSourceIsUnproven|TestAReceiptThroughASymbolicLinkDoesNotResolve|TestAReceiptShorterThanThreeWordsIsUnproven|TestAFileReceiptIsProvedAgainstTheFile|TestAReceiptInsideAFencedBlockIsNotRead|TestAHeldAttributionWithoutAReceiptIsAGap|TestAReceiptForAnotherRecordDoesNotCoverTheClaim|TestASpecThatIsNotHeldStillProvesItsReceipts|TestReceiptFindingsRenderSurfaceTranscriptsOneAndTwo|TestCitationCharacterization)$" ./internal/speccheck 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReceiptCharacterizationKeepsEveryParsedClaim TestASpecCommittedBeforeTheGuideIsNotHeld TestASpecCommittedWithOrAfterTheGuideIsHeld TestAnUncommittedSpecIsHeldOnceTheGuideIsCommitted TestARepositoryWithoutTheGuideHoldsNoSpec TestAnUncommittedGuideHoldsNoSpec TestAGuideWithUnreadableHistoryHoldsEverySpec TestAVerbatimReceiptIsProven TestAReceiptThatWrapsAcrossLinesIsProven TestAReceiptWithOneChangedWordIsUnproven TestAReceiptWithAnUnresolvedSourceIsUnproven TestAReceiptThroughASymbolicLinkDoesNotResolve TestAReceiptShorterThanThreeWordsIsUnproven TestAFileReceiptIsProvedAgainstTheFile TestAReceiptInsideAFencedBlockIsNotRead TestAHeldAttributionWithoutAReceiptIsAGap TestAReceiptForAnotherRecordDoesNotCoverTheClaim TestASpecThatIsNotHeldStillProvesItsReceipts TestReceiptFindingsRenderSurfaceTranscriptsOneAndTwo TestCitationCharacterization; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the eighteen new named tests do not exist, so the command fails.
- `out="$(go test -count=1 -tags docscontract -v -run "^(TestCheckCorpusGolden|TestCheckActiveCorpusHasNoErrors)$" ./internal/docscontract 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; pin="$(go test -count=1 -v -run "^TestArchiveLayoutCharacterizationPinsCorpusGoldenAfterSpec0095$" ./internal/spec 2>&1)" || { printf "%s\\n" "$pin"; exit 1; }; for code in SC-RECEIPT-UNPROVEN SC-RECEIPT-MISSING; do grep -qF -- "\"$code\": 0" internal/docscontract/testdata/corpus-golden.json || { printf 'missing golden code: %s\n' "$code" >&2; exit 1; }; for file in .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$code" || { printf 'missing %s in %s\n' "$code" "$file" >&2; exit 1; }; done; done && make skills-sync-check` — expected: exit 0; before this Task the corpus golden and the Roundfix skill do not name either code, so the command fails.

## References

- [_prd.md](_prd.md) — Goals 1 and 3; Core Features 1, 2 and 5; Success Metrics 1 and 3
- [_techspec.md](_techspec.md) — Interfaces; Invariants; The contract horizon; Claim Receipts; API Contracts 1-3; Surface Transcripts 1-2; Testing Approach 1, 2, 3, 5 and 6; Build Order 1
- ADR-0183; ADR-0116; ADR-0168; ADR-0093; ADR-0094

## Result

Implemented the Task 01 slice: Claim Receipt parsing and proof, paragraph-local
claim pairing, the guide commit horizon, both receipt findings and their stage
registration, corpus characterization, and the Roundfix skill identifiers.
The Daemon retains status and declared Verification ownership. No commit,
push, Pull Request, Task Graph edit, other Task edit, or constraints.go edit
was performed.

### Characterization recorded before parser changes

The four archived artifacts were copied byte-identically, and the unchanged
`CitationClaims` parser recorded 13 claims (artifact, line, target and subject)
in `claims-golden.json`. The fixed-golden characterization test passed before
extracting the shared paragraph walker. The initial recording run wrote the
golden and therefore triggered suiteguard's repository-write refusal; the
recording branch was removed, and the read-only comparison exited 0 on the
unchanged parser. The final test contains no regeneration path and continues
to compare both fixture bytes and every parsed claim.

### Acceptance evidence

| Acceptance criterion | Implementation and focused evidence |
| --- | --- |
| Claims remain identical | `TestReceiptCharacterizationKeepsEveryParsedClaim` passes before and after the shared-walker extraction; the four fixture files still match their archived sources byte for byte. Existing citation characterization and the unchanged existing speccheck tests pass. |
| Exact and wrapped quotes; changed words, missing sources, short quotes | `TestAVerbatimReceiptIsProven` and `TestAReceiptThatWrapsAcrossLinesIsProven` prove normalized contiguous presence. The separate changed-word, unresolved-source and shorter-than-three-words tests assert the error finding and its reason. Additional tests prove case sensitivity, literal markup and refusal before reading a short receipt's source. |
| File receipts and fenced examples | `TestAFileReceiptIsProvedAgainstTheFile` proves the backticked path form. `TestAReceiptInsideAFencedBlockIsNotRead` asserts that neither receipts nor claims are read from a fenced block. Source-token tests reject backticked fields, backticked ADR identifiers and whitespace inside paths; Unicode whitespace after the colon is accepted. |
| Held claims require same-record, same-paragraph receipts | `TestAHeldAttributionWithoutAReceiptIsAGap`, `TestAReceiptForAnotherRecordDoesNotCoverTheClaim`, and `TestAReceiptInAnotherParagraphDoesNotCoverTheClaim` report the gap. `TestAHeldAttributionWithItsReceiptHasNoGap` reports neither new finding. Only claims resolving to accepted active records are held; unresolved and inactive claims are separately tested. |
| Horizon and skips | The seven required horizon tests cover older PRDs, same-commit and later PRDs, uncommitted PRDs, absent guides, uncommitted guides and unreadable history, including the rendered skip reason. Tests use only gittest repositories or temporary directories. Additional cases cover shallow history, nested roots, external PRDs, revised old PRDs, divergent ancestry, readoption of the guide, non-file guides and filesystem inspection failure. The PRD adding-commit helper is shared with newADRHorizon; its earlier behavior is preserved. |
| Unheld Specs still prove receipts | `TestASpecThatIsNotHeldStillProvesItsReceipts` reports the unproven error while the missing-receipt detector is skipped. Symlinks to files inside and outside the repository, at the leaf and at a directory component, report an unproven error ending `does not resolve to a file`. ADR directory links, ambiguous records, traversal, absolute paths and non-regular sources are also tested. |
| Finding text matches Surface Transcripts 1 and 2 | `TestReceiptFindingsRenderSurfaceTranscriptsOneAndTwo` asserts each complete finding block, including severity, summary, both locations and fix text. Stage tests prove PRD-only input at prd and PRD plus TechSpec input at later stages, and both receipt skips when the PRD is absent. |
| Corpus and skill mirrors | Both codes join corpusFindingCodes and the corpus golden at 0; every earlier count is unchanged. The archive-layout pin matches the updated golden. The canonical Roundfix skill and mirror name both codes and explain the horizon, and their bytes are identical. Roundfix skill version 0.0.5 is recorded. |

### Focused checks

All Go checks below used `GOCACHE=/tmp/roundfix-0191-gocache`.

- `go test ./internal/speccheck -run '^TestReceiptCharacterizationKeepsEveryParsedClaim$' -count=1` — exit 0 before parser changes.
- `go test ./internal/speccheck -count=1` — exit 0 after the final receipt grammar changes; existing tests remain unchanged.
- `go test ./internal/speccheck ./internal/spec ./skills -count=1` — exit 0; includes the archive-layout corpus pin and owned skill version checks.
- `go test -tags docscontract ./internal/docscontract -run 'TestCheck' -count=1` — exit 0 after the final changes; includes actual active corpus counts and absence of active corpus errors.
- Temporary Go overlay disabling detectReceipts, with the changed-word and missing-receipt tests — expected exit 1; both planted-defect tests failed because findings were absent. The overlay remained under /tmp and did not modify repository sources.
- The three final grammar regression tests first failed on the prior implementation, then passed with the strict path-token and Unicode-whitespace fixes.
- `make skills-sync` — exit 0. `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` — exit 0, generated only the new Roundfix version entry.
- `make baseline-digests` — exit 0, `changed: false`; no Baseline pin changes.
- Read-only corpus/scope audit — only the two new codes are added, both at 0; every previous count is retained, the skill mirror is byte-identical, and excluded task paths are unchanged.
- `make verify-incremental` — first sandboxed run exited 2 because two existing force-stop integration tests could not read the host process table. Both host-permission reruns exited 0, including the final rerun after the grammar fixes; formatting, vet, tests, skill checks and build passed.
- `git -c core.fsmonitor=false diff --check` — exit 0.

### Additional declared generated path

- `skills/testdata/owned-skill-versions.json` — generated by the owned-skill
  recording command required by docs/agents/specific-repository.md after
  raising both Roundfix skill version fields. Only the 0.0.5 version/digest
  entry was added; no digest was edited by hand.

The Task's two declared Verification commands were not run. This Result
records implementation and focused evidence for Daemon Verification and
settlement, without a terminal Task verdict.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `skills/testdata/owned-skill-versions.json`

## Carry-forward provenance

- Source Run: `run_20261001T000219Z_0e911d72b0d72d2f`
- Source commit: `5d2dd2bc5927347a6469ca616c4d44050c2ce11c`
