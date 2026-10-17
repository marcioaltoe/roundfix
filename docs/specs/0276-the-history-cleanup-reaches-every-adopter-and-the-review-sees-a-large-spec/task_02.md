---
task: task_02
spec: 0276-the-history-cleanup-reaches-every-adopter-and-the-review-sees-a-large-spec
status: completed
type: backend
complexity: high
---

# Task 02: The pre-PR review omits retired history and generated files, reviews a large candidate in segments, and exempts the QA Task from C2

## Overview

Today the review refuses a `codex` candidate above the bound before any
reviewer reads it, even when most of its bytes are deleted history and
regenerated goldens, and Delivery Convention C2 asks a QA Task for settlement
sections (`_techspec.md` → Current behavior). This Task adds the
`retired-history` and `generated` omissions, reviews a candidate still above
the bound in up to four Review Segments in one session with one merged record,
rewords C2, and adds **Review Segment** to the glossary. It is verifiable on
its own through review command tests with the fake reviewer runner, which make
no provider call.

## Requirements

1. MUST answer the adopted Backlog Entry of 2026-10-17,
   [the pre-PR review refuses a Spec that deletes retired history](references/2026-10-17-the-pre-pr-review-refuses-a-spec-that-deletes-retired-history.md),
   by implementing API Contract 4 of `_techspec.md` in `reviewScopedDiff` in
   `internal/cli/review.go`: the two reasons, their precedence, the new summary
   line, and derived declarations read from the base commit through
   `config.DeliveryConfigAtCommit`, never from the candidate head, without
   directory entries (ADR-0282).
2. MUST implement API Contract 5 in `runConfiguredReviewSession` and the review
   command: segments only for `codex` round 1 above `reviewDiffBound`, packed
   greedily over first-parent commits, at most `reviewSegmentLimit` (4), each
   prompt sent in order with `RunPrepared` in the one prepared session, the
   blocking reasons of API Contract 5 before any provider call, and one record
   with `segments`, per-finding `segment`, findings numbered across segments
   and each anchor validated against its own segment's diff. A single commit
   above the bound MUST keep today's reason, and round 2 and `claude` MUST keep
   today's behavior. Verdict reuse, dispositions, the convention validator and
   round accounting stay per candidate head (ADR-0197).
3. MUST answer the adopted Backlog Entry of 2026-10-16,
   [review convention C2 asks a QA Task for settlement sections](references/2026-10-16-review-convention-c2-asks-a-qa-task-for-settlement-sections.md),
   by implementing API Contract 6 in `internal/cli/review_conventions.go`:
   version `roundfix/delivery-conventions/v3` and the exact C2 text; the other
   conventions keep their bytes.
4. MUST add the tests named in Verification, each calling `t.Parallel()`
   first: in `internal/cli/review_scope_test.go`,
   `TestReviewOmitsRetiredHistoryAndGeneratedFiles` (a deleted
   `docs/history/` file and a file a base-commit glob declaration names are
   listed with their reasons and left out of the diff, and the summary line
   counts each reason) and
   `TestReviewReadsDerivedPathsFromTheBaseCommit` (a declaration the candidate
   adds, and a directory entry, omit nothing); in the new
   `internal/cli/review_segment_test.go`,
   `TestReviewSegmentsACandidateAboveTheBound` (a three-commit candidate above
   the bound sends three prompts with their segment lines in one prepared
   session and writes one record whose `segments`, `diffBytes`, finding ids and
   `segment` fields follow API Contract 5) and
   `TestReviewBlocksACandidateNeedingMoreThanFourSegments` (five segments block
   with the reason of API Contract 5 and reach no provider); in the new
   `internal/cli/review_convention_qa_task_test.go`,
   `TestDeliveryConventionC2ExemptsTheQATask` (the prompt for a candidate
   holding a completed QA Task file shaped like Spec 0274's `task_05.md`,
   without settlement sections, carries v3 and the exact C2 text, and no
   provider is called).
5. MUST move the summary line asserted by
   `TestReviewOmitsQAEvidenceAndUpstreamSkills` to API Contract 4's form; this
   is the declared break of its old expected text. `TestReviewBlocksAboveTheBoundWithoutAProviderCall`,
   `TestReviewCommandKeepsTheCodexDiffByteBound` and
   `TestReviewCommandSendsAClaudePromptBeyondTheCodexByteBound` MUST pass
   unchanged. MUST NOT rename or remove any test.
6. MUST document the two omission reasons, the summary line, the Review
   Segments with their blocking reasons and record fields, the v3 version and
   the new C2 text in `docs/user-guide/commands/review.md` and
   `.agents/skills/roundfix/references/review.md`, replacing the old C2
   sentence and the v2 name, searching the guides and skill copies for both.
   No text inside any `### QA settlement` section may change.
7. MUST, through `domain-modeling`, add **Review Segment** to `CONTEXT.md` near
   **Reviewer Lineage**, linked to ADR-0282: one of at most four consecutive
   first-parent commit ranges of a pre-PR review candidate whose diff exceeds
   the review bound, each within the bound and reviewed in order in one review
   session, whose findings form one Pre-PR Review Record. The entry MUST pass
   the glossary contract.
8. MUST run `make skills-sync`, raise and record the Roundfix skill version
   with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   and re-record the Behavior Surface Record with
   `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
   because this Task edits a user guide page and command sources. No record or
   mirror is hand-edited.

## Subtasks

- [ ] Add the two omissions, their base-commit declarations and the summary line.
- [ ] Plan, send and merge Review Segments, with their refusals.
- [ ] Reword C2 under v3.
- [ ] Add the five tests and move the summary assertion.
- [ ] Document the review changes and add the glossary entry.
- [ ] Sync the mirrors, record the version and the Behavior Surface Record.

## Acceptance Criteria

- [ ] Deleted history and base-declared generated files are omitted with their reasons; the candidate's own declarations and directory entries omit nothing.
- [ ] A `codex` round-1 candidate above the bound is reviewed in up to four segments in one session with one record; one oversized commit, five segments and an oversized round-2 delta block.
- [ ] The prompt carries v3 and a C2 that exempts the QA Task.
- [ ] The guide, the reference and the glossary describe the change.

## Consumers

- The Pre-PR Review Record and the review prompt — consumers: `./internal/cli`, `./internal/docscontract`

## Context

- instruction: `docs/adr/0282-the-history-cleanup-names-a-superseded-adr-in-its-successor-and-a-large-candidate-is-reviewed-in-segments.md`
- instruction: `docs/adr/0197-a-pre-pr-reviewer-lineage-spans-at-most-two-rounds.md`
- instruction: `docs/adr/0196-a-pre-pr-review-finding-parks-only-after-validation.md`
- instruction: `docs/adr/0227-a-review-through-claude-keeps-its-verdict-and-fits-its-window.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `internal/config/delivery.go`
- interface: `internal/cli/review.go`
- interface: `internal/cli/review_validation.go`
- interface: `internal/cli/review_conventions.go`
- interface: `internal/cli/review_scope_test.go`
- creates: `internal/cli/review_segment_test.go`
- creates: `internal/cli/review_convention_qa_task_test.go`
- interface: `docs/user-guide/commands/review.md`
- interface: `.agents/skills/roundfix/references/review.md`
- interface: `skills/roundfix/references/review.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestReviewOmitsRetiredHistoryAndGeneratedFiles|TestReviewReadsDerivedPathsFromTheBaseCommit|TestReviewSegmentsACandidateAboveTheBound|TestReviewBlocksACandidateNeedingMoreThanFourSegments|TestDeliveryConventionC2ExemptsTheQATask)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestReviewOmitsRetiredHistoryAndGeneratedFiles TestReviewReadsDerivedPathsFromTheBaseCommit TestReviewSegmentsACandidateAboveTheBound TestReviewBlocksACandidateNeedingMoreThanFourSegments TestDeliveryConventionC2ExemptsTheQATask; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; grep -qF -- 'roundfix/delivery-conventions/v3' internal/cli/review_conventions.go || { printf 'the conventions version is not v3\n' >&2; exit 1; }; for f in docs/user-guide/commands/review.md .agents/skills/roundfix/references/review.md; do grep -qF -- 'retired-history' "$f" || { printf 'missing retired-history in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'roundfix/delivery-conventions/v3' || { printf 'missing v3 in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Review Segment' || { printf 'missing Review Segment in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Review Segment**:' || { printf 'CONTEXT.md does not define **Review Segment**\n' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -run 'Review' ./internal/cli || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the five tests do not exist, so the command fails at its first pass check; afterwards the review, documentation, glossary contract, skill version and Behavior Surface Record checks pass.

## References

- `_prd.md` → Goals; User Story 3; User Story 4; Core Feature 3; Core Feature 4; Core Feature 5; Core Feature 9; Success Metric 2; Success Metric 3; Glossary
- `_techspec.md` → API Contract 4; API Contract 5; API Contract 6; API Contract 10; Invariants 4, 5 and 6; Current behavior; Testing Approach; Build Order 2
- ADR-0282; ADR-0196; ADR-0197; ADR-0227


## Result

Implemented this Task's review slice; status and terminal Verification remain
Daemon-owned.

### Acceptance evidence

- Omissions: `reviewScopedDiff` reads derived declarations through
  `DeliveryConfigAtCommit` at the immutable candidate base, with no User
  Config input. Exact/glob file declarations omit generated files; directory
  entries and declarations added by the candidate omit nothing. Deleted
  history is omitted except Archive Records. The five-reason precedence and
  exact summary are implemented. Focused checks passed for
  `TestReviewOmitsRetiredHistoryAndGeneratedFiles`,
  `TestReviewReadsDerivedPathsFromTheBaseCommit`,
  `TestReviewKeepsDeletedArchiveRecordsAndOmissionPrecedence`,
  `TestReviewBlocksOnUnreadableBaseDerivedDeclarations`, and the existing
  QA-evidence, upstream-skill and round-2 scope checks.
- Review Segments: Codex round 1 greedily packs first-parent commits into at
  most four bounded ranges. Every individual commit must fit, including a
  rewrite whose combined range hides its size. Sequential `RunPrepared`
  calls share one prepared session; Spec context occurs only in the first
  prompt. Findings receive global IDs and their segment number, and their
  anchors are validated against their own segment diff. One candidate-head
  record holds the ranges and summed bytes; the answer preserves each
  segment's messages under the specified separators. Classification failures
  stop subsequent prompts. The convention validator, reuse and lineage remain
  per candidate head. Focused checks passed for
  `TestReviewSegmentsACandidateAboveTheBound`,
  `TestReviewBlocksACandidateNeedingMoreThanFourSegments`,
  `TestReviewSegmentsStopAtTheFirstBlockedAnswer`,
  `TestReviewGreedilyPacksFourSegments`, and
  `TestReviewBlocksAnOversizedCommitHiddenByItsRange`.
  `TestReviewBlocksAboveTheBoundWithoutAProviderCall`,
  `TestReviewCommandKeepsTheCodexDiffByteBound` and
  `TestReviewCommandSendsAClaudePromptBeyondTheCodexByteBound` passed without
  edits. Segment fixtures also prove that a candidate-added declaration
  cannot hide subsequent segments, and that same-head reuse sends no further
  prompts. The supporting production helper is
  `internal/cli/review_segments.go`, an ordinary path for this Task's slice.
- C2: the prompt uses `roundfix/delivery-conventions/v3` and API Contract 6's
  exact text. C1 and C3–C5 keep their bytes.
  `TestDeliveryConventionC2ExemptsTheQATask` passed using a committed completed
  QA Task without settlement sections, constructing the prompt without a
  provider call.
- Documentation and vocabulary: the review guide and canonical skill
  reference describe both new reasons, precedence, the summary, segments,
  refusals, record fields and the v3 C2 text. `CONTEXT.md` defines Review
  Segment beside Reviewer Lineage with ADR-0282. Focused glossary contracts
  passed. A direct document audit confirmed exact normalized C2 text in the
  guide and both skill references, no old v2 name there, matching reference
  mirrors and byte-identical QA settlement sections against HEAD.

### Focused commands and outcomes

- Initial regression signal:
  `rtk proxy go test ./internal/cli -run 'TestReviewOmitsRetired|TestDeliveryConventionC2' -count=1`
  exited 1 before implementation: missing omissions and the old v2 C2 prompt.
- Final focused implementation check:
  `rtk proxy go test ./internal/cli -run '^(TestReview(Omits|Reads|Scopes|Blocks|KeepsDeleted|Segments|Greedily|CommandKeepsTheCodexDiffByteBound|CommandSendsAClaudePromptBeyondTheCodexByteBound)|TestDeliveryConvention)' -count=1`
  exited 0 after the final source edit. All reviewer execution used local
  fakes and disposable repositories.
- `rtk proxy go test -tags docscontract ./internal/docscontract -run '^TestGlossary(PhrasesAreDefinedTerms|LinksNameAcceptedADRs|AvoidLinesArePlainLists|StatesEachRuleOnce)$' -count=1`
  exited 0.
- Required generation: `rtk proxy make skills-sync` exited 0;
  `rtk proxy go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
  exited 0 and recorded Roundfix version `0.1.87` in both skill sources and
  the owned-version record;
  `rtk proxy go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`
  exited 0, including a re-record after the final command-source edit.
  Mirrors and generated records were not hand-edited.
- `rtk proxy make baseline-digests` exited 0 and reported no derived changes.
- `rtk proxy git diff --check` exited 0. Final changed-file inspection is
  limited to this review slice, its documentation, the generated skill and
  surface records, and this Task's Result; the status edit predated this turn.

The declared Verification command, repository Verification and full test
suite were not run. No other Task or Task Graph was edited, and no commit,
push or Pull Request was made. The Daemon owns settlement and any required
Recorded paths entry for the supporting helper.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/cli/review_segments.go`
