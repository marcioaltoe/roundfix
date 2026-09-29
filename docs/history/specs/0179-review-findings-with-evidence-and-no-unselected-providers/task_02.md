---
task: task_02
spec: 0179-review-findings-with-evidence-and-no-unselected-providers
status: completed
type: backend
complexity: high
---

# Task 02: Every finding gets a recorded disposition

## Overview

`classifyReviewCommandResult` in `internal/cli/review.go` stores a findings
verdict as one block of text in `pre-pr-review.json`, so a finding has no
identity, and nothing records whether it was fixed or dismissed. On 2026-09-25
the operator recorded those dispositions only in pull request bodies (#257,
#258, #259). This Task gives each finding an ordinal identity and adds
`roundfix review dispose`, which appends one evidence-backed disposition per
finding to a ledger beside the record. The record is written by `roundfix
review` into the Artifact Directory, outside the reviewed tree. The ledger is
written by the operator's `dispose` call and read by later `roundfix review`
calls (task_03) and by whoever audits a delivery, so a refused call must write
nothing, and the evidence text is copied verbatim and never executed.

## Requirements

1. MUST add `reviewFinding` (`id`, `text`) and the record field
   `findingItems` (JSON `findingItems`, omitted when empty), and set it in the
   findings branch of `classifyReviewCommandResult` from a new
   `splitReviewFindings(findings string) []reviewFinding`:
   - A line with no leading whitespace that starts with `- `, `* `, or ASCII
     digits followed by `.` or `)` and a space opens a new finding.
   - Any other line continues the current finding.
   - Text before the first marker is one finding when it is not blank, and
     text with no marker is one finding.
   - Markers and surrounding whitespace are trimmed, blank findings are
     dropped, and identities are `F1`, `F2`, … in order.
   A record read from disk without `findingItems` MUST derive them from
   `findings` the same way. `validateReviewRecord` MUST keep accepting a
   findings record without items, and `record.Findings` MUST keep its current
   text.
2. MUST make `buildReviewPrompt` ask for each finding as a list item that
   starts with `- ` and names its file and line, keeping every other
   instruction.
3. MUST parse `roundfix review dispose <finding-id>` when the first argument of
   `review` is `dispose`, with exactly one of `--dismiss --evidence <text>` or
   `--fixed-by <commit>`. It reads `pre-pr-review.json` from the Artifact
   Directory `runReviewCommand` resolves, and requires a record whose
   repository is the current Git root and whose outcome is `findings` (task_03
   adds `findings-dismissed`).
4. MUST accept `--dismiss` only with non-blank `--evidence` and only when the
   current `HEAD` equals the record's `headCommit`. MUST accept `--fixed-by`
   only when the commit resolves, differs from the record's head, has that head
   as an ancestor, and is reachable from the current `HEAD`. Use `git
   merge-base --is-ancestor` through the command's Git runner.
5. MUST append, on success, one JSON line to `pre-pr-review-dispositions.jsonl`
   in the Artifact Directory and print the same line on stdout, exit `0`. The
   line carries `repository`, `headCommit`, `finding`, `text` (copied from the
   record), `disposition` (`dismissed` or `fixed`), `evidence` or `fixedBy`,
   and `recordedAt` in RFC 3339 UTC. The ledger is only ever appended; an
   existing line is never rewritten.
6. MUST refuse, with exit `2`, stderr `roundfix: review dispose refused:
   <reason>` and no append:
   - no record, a record of another repository, or another outcome;
   - an unknown finding identity;
   - both forms, or neither;
   - blank evidence;
   - a dismissal at a moved `HEAD`;
   - a `--fixed-by` commit that fails any check of Requirement 4;
   - a finding that already has a disposition for the same repository, head,
     identity and text.
   `--help` prints the usage and exits `0`.
7. MUST list both forms in the top-level usage and in `commandUsage("review")`
   in `internal/cli/cli.go`, keeping the existing `roundfix review [--base
   <ref>]` line byte-identical.
8. MUST describe finding identities, both forms, the ledger and the refusals in
   the review section of `docs/user-guide/commands.md` and the Pre-PR review
   section of `.agents/skills/roundfix/SKILL.md`, using the phrase
   `roundfix review dispose` in both. Regenerate `skills/roundfix/SKILL.md`
   with `make skills-sync`. MUST add a Review Finding Disposition entry to
   `CONTEXT.md` defining a disposition as a finding's recorded fix or
   evidence-backed dismissal, tied to the head it reviewed.
9. MUST add `internal/cli/review_disposition_test.go` with
   `TestReviewRecordListsEachFindingWithAnIdentity` (a marked answer yields
   `F1..F3` with their texts, and an unmarked answer yields one `F1`),
   `TestReviewPromptAsksForOneListItemPerFinding`,
   `TestReviewDisposeDismissesAFindingWithEvidence`,
   `TestReviewDisposeRecordsAFixingCommit`,
   `TestReviewDisposeRefusesWithoutAFindingsRecord`,
   `TestReviewDisposeRefusesAnUnknownFinding`,
   `TestReviewDisposeRefusesBothForms`,
   `TestReviewDisposeRefusesNeitherForm`,
   `TestReviewDisposeRefusesBlankEvidence`,
   `TestReviewDisposeRefusesADismissalAtAMovedHead`,
   `TestReviewDisposeRefusesAFixThatDoesNotDescendFromTheReviewedHead`, and
   `TestReviewDisposeRefusesASecondDisposition`. Each refusal test asserts
   exit `2` and a ledger byte-identical to before (absent when it was absent).
   The two success tests assert every field of the appended line and that
   stdout equals it. `TestReviewDisposeRefusesASecondDisposition` records a
   dismissal and then asserts that a fix of the same finding is refused, the
   mirror of the success path.
10. MUST keep `TestReviewCommandExitsOneAndRecordsFindings` and
    `TestReviewRecordRoundTripsEachOutcome` green and unchanged.

## Subtasks

- [ ] Split findings into identified items and ask for list items.
- [ ] Add `roundfix review dispose` and its ledger.
- [ ] Document the command, the ledger and the glossary term.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A findings record lists every finding with an identity.
- [ ] A dismissal needs evidence at the reviewed head; a fix needs a
      descendant commit reachable from `HEAD`.
- [ ] Each finding takes at most one disposition, and every refusal appends
      nothing.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/review.go`
- interface: `internal/cli/cli.go`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `docs/user-guide/commands.md`
- interface: `CONTEXT.md`
- creates: `internal/cli/review_disposition_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestReviewRecordListsEachFindingWithAnIdentity|TestReviewPromptAsksForOneListItemPerFinding|TestReviewDisposeDismissesAFindingWithEvidence|TestReviewDisposeRecordsAFixingCommit|TestReviewDisposeRefusesWithoutAFindingsRecord|TestReviewDisposeRefusesAnUnknownFinding|TestReviewDisposeRefusesBothForms|TestReviewDisposeRefusesNeitherForm|TestReviewDisposeRefusesBlankEvidence|TestReviewDisposeRefusesADismissalAtAMovedHead|TestReviewDisposeRefusesAFixThatDoesNotDescendFromTheReviewedHead|TestReviewDisposeRefusesASecondDisposition|TestReviewCommandExitsOneAndRecordsFindings|TestReviewRecordRoundTripsEachOutcome|TestReviewCommandAppearsOnPublicHelp)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReviewRecordListsEachFindingWithAnIdentity TestReviewPromptAsksForOneListItemPerFinding TestReviewDisposeDismissesAFindingWithEvidence TestReviewDisposeRecordsAFixingCommit TestReviewDisposeRefusesWithoutAFindingsRecord TestReviewDisposeRefusesAnUnknownFinding TestReviewDisposeRefusesBothForms TestReviewDisposeRefusesNeitherForm TestReviewDisposeRefusesBlankEvidence TestReviewDisposeRefusesADismissalAtAMovedHead TestReviewDisposeRefusesAFixThatDoesNotDescendFromTheReviewedHead TestReviewDisposeRefusesASecondDisposition TestReviewCommandExitsOneAndRecordsFindings TestReviewRecordRoundTripsEachOutcome TestReviewCommandAppearsOnPublicHelp; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && tr -s '[:space:]' ' ' < docs/user-guide/commands.md | grep -qF -- "roundfix review dispose" && tr -s '[:space:]' ' ' < .agents/skills/roundfix/SKILL.md | grep -qF -- "roundfix review dispose" && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "**Review Finding Disposition**" && diff -r .agents/skills/roundfix skills/roundfix >/dev/null` — expected: exit 0; before this Task none of the twelve new named tests exists and no guide names `roundfix review dispose`, so the command fails.

## References

- [_techspec.md](_techspec.md) — Finding identity and dispositions
- `_prd.md` → Goal 2; Core Feature 2; Success Metric 2
- `_techspec.md` → API Contracts 1-3; Testing Approach 2

## Result

### Implementation

- Added `findingItems` with ordinal `F<n>` identities while preserving the
  existing findings text. The splitter recognizes the declared bullet and
  numbered markers, keeps unmarked preambles and continuation lines together,
  drops blank items, and derives items when an older record omits them.
- Added `roundfix review dispose` with the dismissal and fixing-commit forms,
  Git-runner ancestry and reachability checks, exact finding lookup, duplicate
  detection, RFC 3339 UTC timestamps, append-only JSONL persistence, and the
  declared stdout, stderr and exit-code contract. Evidence remains inert text.
- Added the two usage forms to top-level and review help. Documented finding
  identities, the disposition forms, ledger and refusal behavior in the user
  guide and canonical Roundfix skill, and added the Review Finding Disposition
  glossary entry.
- Added the twelve named disposition tests plus focused marker and help
  coverage. Every refusal snapshots the ledger before the call and proves it
  remains byte-identical or absent.
- Regenerated `skills/roundfix` with `rtk make skills-sync`; `rtk make
  baseline-digests` reported that all derived artifacts already matched their
  canonical sources and changed no additional files.

### Focused checks and acceptance evidence

1. Finding identity: `rtk env GOCACHE=/tmp/roundfix-task02-gocache go test
   -count=1 -run
   '^TestReview(RecordListsEachFindingWithAnIdentity|PromptAsksForOneListItemPerFinding|Dispose)'
   ./internal/cli` exited 0. The identity test covers marked and unmarked
   answers plus legacy record derivation; the companion splitter test covers a
   preamble, both numeric separators, an indented pseudo-marker and a blank
   item.
2. Head-bound dismissal and fixing commit: the same focused command exited 0
   for the two success paths and the moved-head, unresolved, same-head,
   non-descendant and unreachable fix refusals. Both success tests assert every
   ledger field and byte equality between stdout and the appended line.
3. One disposition and refusal immutability: the focused command exited 0 for
   absent, cross-repository and non-findings records; unknown identity; both,
   neither and blank forms; and the second-disposition mirror path. Each case
   asserted exit 2, the required stderr prefix, empty stdout and byte-identical
   ledger state.
4. Preserved review behavior: `rtk env
   GOCACHE=/tmp/roundfix-task02-gocache go test -count=1 -run
   '^(TestReviewCommandExitsOneAndRecordsFindings|TestReviewRecordRoundTripsEachOutcome|TestReviewCommandAppearsOnPublicHelp)$'
   ./internal/cli` exited 0 without modifying those tests.
5. Documentation and generated copy: `rtk diff -r
   .agents/skills/roundfix skills/roundfix` and `rtk git diff --check` exited 0;
   repository searches found `roundfix review dispose` in both required guides
   and `**Review Finding Disposition**` in `CONTEXT.md`.
6. Repository incremental check: an initial sandboxed package run reached the
   new review coverage but two unrelated force-stop tests could not read the
   host process table. The required elevated rerun of `rtk env
   GOCACHE=/tmp/roundfix-task02-gocache make verify-incremental` exited 0,
   including `go vet`, all Go packages, owned-skill checks and the CLI build.

The Task's authored `## Verification` command was not run; the Daemon owns
that command and terminal settlement.
