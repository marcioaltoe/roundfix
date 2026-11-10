---
spec: 0372-a-review-whose-reviewed-head-is-gone-opens-a-new-lineage
prd: _prd.md
created: 2026-11-10
---

# A review whose reviewed head is gone opens a new lineage — Technical Spec

## Executive Summary

A new helper in `internal/cli`, `reviewRecordedCommitPresence`, proves
whether a commit a review record or the disposition ledger names is gone. It
answers "gone" only for a full object ID in a repository that is neither
shallow nor partial, when `git cat-file -e <id>` exits 1 with no output; a
shallow or partial clone answers "history incomplete"; anything else answers
"not proven gone". The four ancestry checks that read such a commit call it
only after `git merge-base --is-ancestor` fails with a status other than 1:
`decideReviewLineage` restarts the lineage with a reason naming the commit,
`review dispose` and the Delivery Retry's correction proof refuse by name, and
the disposition validity check counts a gone fixing commit as invalid. The
trade-off accepted is that a review in a shallow or partial clone whose
recorded commit is missing stays blocked, now with a reason that says how to
recover, rather than restarting on an absence the clone cannot prove; this
keeps the two-round ceiling of ADR-0197 from being reset by anything but a
proven non-descent.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted or changed. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local Git reads only; a partial
  clone is never probed because the probe would fetch. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0197 (a head that does not descend
  starts a new lineage), ADR-0272 (the correction proof of a Delivery Retry)
  and ADR-0279 (the dismissal anchor proof). ADR-0233, ADR-0286, ADR-0366 and
  ADR-0373 are not applicable: they cite the three, and their rules are
  unchanged. Source: `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — no protected tooling mutation proposed or
  authorized; the Governed Paths task_01 changes are the owned Roundfix skill
  files the maintainer's cycle grants cover. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0372-a-review-whose-reviewed-head-is-gone-opens-a-new-lineage/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/review.md`,
  `skills/roundfix/SKILL.md`.

## System Architecture

- `internal/cli/review_recorded_commit.go` (new, task_02):
  `reviewRecordedCommitPresence`, `reviewHistoryIncompleteError` and
  `reviewedHeadGoneRefusal`. It runs Git directly with `exec.CommandContext`
  (as `spec_folder_status.go` and `implement.go` do), because the proof needs
  the exit status and standard error apart, which `preflight.GitRunner`
  folds into one error string.
- `internal/cli/review_lineage.go` (task_02): `decideReviewLineage` calls the
  helper when the head or base ancestry check fails with a status other than
  1; `reviewDispositionValidAtHead` returns `false, nil` for a gone fixing
  commit. `reviewLineageNotice` (`deliver status`, `deliver retry`) changes no
  code and prints the new reasons.
- `internal/cli/review.go` (task_02): the `--fixed-by` branch of the
  `review dispose` command and `proveDismissalAnchor` call
  `reviewedHeadGoneRefusal` before their existing refusal.
- `internal/cli/deliver_review_correction.go` (task_02):
  `ProveReviewCorrection` refuses a gone parked candidate by name.
- Guides and references (task_01): `docs/user-guide/commands/review.md`,
  `docs/user-guide/commands/deliver.md`, the canonical
  `.agents/skills/roundfix/references/review.md` and `deliver.md`, their
  mirrors, the Roundfix `SKILL.md` version pair and the records.

## Implementation Design

### Interfaces

```go
// internal/cli/review_recorded_commit.go
type reviewCommitPresence int

const (
	reviewCommitNotProvenGone     reviewCommitPresence = iota // present, or Git did not answer
	reviewCommitGone                                          // a complete repository proved it absent
	reviewCommitHistoryIncomplete                             // shallow or partial clone
)

func reviewRecordedCommitPresence(ctx context.Context, root, commit string) reviewCommitPresence
func reviewHistoryIncompleteError(commit string) error
func reviewedHeadGoneRefusal(ctx context.Context, root, reviewed string) error
```

`reviewRecordedCommitPresence`, in order: a value that is not 40 or 64
lowercase hexadecimal digits is not proven gone;
`git rev-parse --is-shallow-repository` must print `true` or `false`;
`git config --get-regexp '^(remote\..+\.promisor|extensions\.partialclone)$'`
exit 0 with an `extensions.partialclone` key or a promisor value `true` is a
partial clone, which returns history incomplete without asking Git about the
commit; exit 1 means none; any other status is not proven gone. Then
`git -C <root> -c core.fsmonitor=false cat-file -e <commit>`: exit 1 with
empty standard output and standard error is gone in a complete repository
and history incomplete in a shallow one; every other answer is not proven
gone.

### Data Models

No changes. Review records, Reviewed Head Records and the disposition ledger
keep their bytes and fields; nothing is rewritten.

### API Contracts

1. API Contract: lineage restart. `decideReviewLineage`, when the
   ancestry check of the recorded head fails with a status other than 1:
   gone → `Restart` = `the reviewed head <prior head> is no longer in the
   repository`, round 1, no prior; history incomplete → error
   `check review lineage ancestry: a shallow or partial clone cannot prove
   that the recorded commit <commit> is gone; fetch the full history and try
   again`; otherwise the error it returns today. The recorded merge base
   behaves the same, with `the recorded merge base <prior base> is no longer
   in the repository`.
2. API Contract: dispose refusal. `review dispose --dismiss` and
   `--fixed-by` print `roundfix: review dispose refused: reviewed head
   "<commit>" is no longer in the repository; run roundfix review to open a
   new Reviewer Lineage` and exit 2 when the reviewed head is gone; the
   incomplete-history reason (without the `check review lineage ancestry: `
   prefix) when the clone cannot prove it; their existing refusals otherwise.
3. API Contract: correction proof. `ProveReviewCorrection` returns the
   refusal `parked candidate <commit> is no longer in the repository` for a
   gone parked candidate and the error `check review correction ancestry: `
   plus the incomplete-history reason when the clone cannot prove it; a gone
   fixing commit makes `reviewDispositionValidAtHead` return `false, nil`, so
   the finding gets the existing `requires exactly one valid disposition`
   refusal.

### Surface Transcripts

None. The changed outputs are single lines named verbatim in API Contracts
1–3; `deliver status` and `deliver retry` keep the `Review:` line format of
ADR-0366 with a new reason.

## Coverage Map

- Goal 1, User Story 1, Core Feature 1, Success Metric 1 → Interfaces; API
  Contract 1; `decideReviewLineage`; `reviewLineageNotice`; Build Order 2.
- Goal 2, User Story 2, Core Feature 2, Success Metric 2 → Interfaces;
  `reviewRecordedCommitPresence`; Testing Approach; Build Order 2.
- Goal 3, User Story 3, Core Feature 3, Success Metric 3 → Interfaces; API
  Contracts 2 and 3; `reviewedHeadGoneRefusal`; `ProveReviewCorrection`;
  `reviewDispositionValidAtHead`; Build Order 2.
- Core Feature 4 → the guides and references; Build Order 1.

Deviations from the PRD: None.

## Absence Evidence

1. A proof that a recorded commit is gone — command: `grep -rn 'cat-file", "-e"\|is-shallow-repository\|promisor' internal/cli --include='*.go'`; result: at `9fda0b65`, three production matches, each treating any failure alike: `internal/cli/review.go` (`cat-file -e <revision>^{commit}`, which skips an archived Spec whose source revision cannot be read), `internal/cli/releaseplan_checks.go` and `internal/cli/deliver_workflow.go` (object checks of the release plan and the queue, which read no review record), plus two test lines in `source_revision_recovery_test.go`; no shallow or promisor check anywhere in `internal/cli`.
2. A restart or refusal reason for a gone commit — command: `grep -rn 'no longer in the repository' internal docs/user-guide .agents/skills/roundfix`; result: no match at `9fda0b65`.

## ADR Alternatives

- ADR-0197: kept; the design restarts only on a proven non-descent, which
  ADR-0197 already restarts on, and takes none of its rejected options (an
  unbounded review loop, a fresh full review every time).
- ADR-0272: kept; the correction proof still requires descent, one valid
  disposition per finding and the path bound, and gains only a named refusal.
- ADR-0279: kept; a dismissal still requires the reviewed head or an
  unchanged anchor at a descendant; a gone reviewed head is refused by name.
- ADR-0233: not applicable; the retry path rule is unchanged.
- ADR-0286: not applicable; the Correction Review is unchanged, and its
  planner already counted an unprovable disposition as invalid.
- ADR-0366: not applicable; the merge-base rules are unchanged, and the
  `Review:` line keeps its format.
- ADR-0373: not applicable; it cites ADR-0279 for owner liveness only.

## Integration Points

None outside the repository. The Delivery Queue reaches the change through
`reviewLineageNotice` and `ProveReviewCorrection`; no fetch, no network, no
provider.

## Testing Approach

New file `internal/cli/review_recorded_commit_test.go` (task_02), real Git
in temporary repositories, the existing review fixtures and fake reviewer:

- `TestAReviewWhoseReviewedHeadIsGoneOpensANewLineage`: a first review on
  `feature/review`; a new branch from `main` with a commit; the old branch
  deleted, `git reflog expire --expire=now --all`, `git gc --prune=now`;
  `git cat-file -e` of the old head exits 1. `decideReviewLineage` returns
  the restart reason, `reviewLineageNotice` the `Review:` line, and the
  second `roundfix review` exits 0 at round 1 with an empty stderr.
- `TestARecordedMergeBaseThatIsGoneOpensANewLineage`: a prior record whose
  base is an absent full ID and whose head is an ancestor.
- `TestARecordedCommitThatIsNotProvenGoneStillBlocks`, subtests: `shallow
  clone` (`--depth 1`, the first commit recorded), `partial clone`
  (`--filter=blob:none`, its remote then pointed at a missing directory),
  `object store error` (a `--no-local` clone whose `objects/info/alternates`
  names a missing directory), `abbreviated commit`, and `Git failure on a
  present commit` (a fake runner failing `merge-base`).
- `TestReviewDisposeRefusesAReviewedHeadThatIsGone`, subtests `dismiss` and
  `fixed-by`.
- `TestAReviewCorrectionNamesACommitThatIsGone`, subtests `parked
  candidate` and `fixing commit`.

The declared break: `TestReviewLineageRestartNamesItsReason/git_failures_stay_errors`
(`internal/cli/review_main_merge_lineage_test.go`) ends by asserting that a
recorded head `strings.Repeat("f", 40)` stays an error; that exact case is the
defect, so its last assertion moves to the restart reason. Its five fake
failures on present commits stay errors unchanged, as does
`TestReviewDisposeReportsAncestryFailureWithoutAppending` and the `Git
failure` case of `TestReviewLineageNoticeIsAdvisory`.

Sabotage proofs, each run on the prototype: treating a shallow clone as
complete fails `shallow clone`; probing a partial clone fails `partial
clone`; ignoring standard error fails `object store error`; restarting on
every failure fails `abbreviated commit` and `Git failure on a present
commit`; dropping the fixing-commit rule fails `fixing commit`. No test
reaches the network, starts an Agent session or writes outside its
temporary directories; new top-level tests call `t.Parallel()` first.

## Build Order

1. The guides and Roundfix skill references, which quote the reasons and
   refusals API Contracts 1–3 fix verbatim; text only.
2. The presence proof, the four call sites and the tests, with the moved
   assertion (depends on: 1, so the code Task names its guides through its
   dependency and the guides already match the reasons it prints).
3. Final QA gate (depends on: 1, 2).

No step needs another Spec's code.

## Risks & Considerations

- A partial clone whose recorded commit is missing stays blocked; the reason
  names the recovery. Roundfix's own Run and item worktrees are complete
  `git worktree` checkouts, so the queue is not affected.
- The extra Git commands run only after an ancestry failure, so a healthy
  review pays nothing.
- An alternate object store that prints an error keeps today's block, so a
  temporarily unavailable alternate cannot pass for an absent commit.

## Glossary

None.

## Decisions

- Prove absence with `git cat-file -e` on the full ID, standard error empty,
  in a complete repository; never the `^{commit}` peel, which exits 128.
- Probe only after the ancestry check fails; a partial clone is never asked.
- No ADR; ADR-0197 already decides the restart.
- Two implementation Tasks and the QA gate.
