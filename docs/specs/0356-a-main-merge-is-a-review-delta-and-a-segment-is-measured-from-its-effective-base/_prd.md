---
spec: 0356-a-main-merge-is-a-review-delta-and-a-segment-is-measured-from-its-effective-base
status: active
created: 2026-11-05
surfaces: [backend, docs]
type: fix
---

# A main merge is a review delta, and a segment is measured from its Effective Base

Minimal PRD for a bug fix, minted by `write-techspec`. Reported by fiscus on
2026-11-05, points 2 and 4 of its report: merging the default branch into a
Delivery Queue item after `review-findings` moved the merge base, so the next
pre-PR review opened a new Reviewer Lineage at round 1, reread the whole
candidate and raised new findings, with no warning before it ran; and a
segmented review was blocked by size because an early item commit, cut from an
old default branch, was measured against the new merge base and showed about
1.5 to 2 MB of the default branch's work as reverted. After this Spec the
review measures only what the candidate itself changed: a conflict-free merge
of the default branch keeps the lineage and its delta, each Review Segment is
measured from its Effective Base, and `deliver status` and `deliver retry`
name a review that will restart at round 1 before it runs.

## Sources

One report, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-05-fiscus-the-delivery-queue-takes-days-and-a-main-merge-restarts-the-review.md`](references/2026-11-05-fiscus-the-delivery-queue-takes-days-and-a-main-merge-restarts-the-review.md),
   the Fiscus ai-memory message `01a1ab64` of 2026-11-05 08:07 BRT, reported
   by fiscus. This Spec answers its points 2 and 4; Specs 0355 (point 1),
   0357 (point 3) and 0358 (point 5) answer the others from their own copies.

At the authoring base `docs/backlog/` and `docs/findings/` hold no open entry
on the pre-PR review, so no other source shares this context.

## Prerequisites

None in `requires`. Authored against `ad77a797` (Specs 0352 and 0353 planned
and being delivered). The planned queue order is 0354, 0356, 0355, 0357, 0358.
This Spec owns the pre-PR review lineage and segmentation code
(`internal/cli/review_lineage.go`, `internal/cli/review_segments.go`, the
round-2 delta in `runConfiguredReviewSession` and the lineage planning in
`runReviewCommand` of `internal/cli/review.go`) and the Roundfix `review`
reference. Files it shares with the other Specs of the cycle, each edited
narrowly so the later one rebases cleanly:

- `internal/cli/deliver.go`: this Spec adds one call in `runDeliverStatus`
  after the `Warning:` loop and one in `runDeliverRetry` after each
  `printDeliverRetryResult`; Spec 0355 changes queue scheduling and status
  lines in the same file.
- `docs/user-guide/commands/deliver.md` and both copies of the Roundfix
  `deliver` reference, which Spec 0355 owns: one paragraph after the `Owner:`
  status paragraph and one after the retry paragraph, nothing else.
- `internal/cli/review.go` and `internal/cli/review_segments.go`, which Spec
  0358 extends later, in other functions, to review large authored data
  files apart.
- The owned Roundfix skill version (`.agents/skills/roundfix/SKILL.md`, its
  mirror and `skills/testdata/owned-skill-versions.json`), `CONTEXT.md`,
  `docs/references/behavior-surfaces.json` and `docs/user-guide/commands/review.md`,
  which several Specs of the cycle raise or re-record.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or schema version
  is minted; two optional record fields are added to the existing Pre-PR
  Review Record. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential; the change reads local Git objects. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0368 (a large authored data file reviewed apart, a sibling Spec in the same queue) touches the same review code but carries no obligation for this Spec; ADR-0366 (new) records the decision: a
  conflict-free merge of the default branch keeps the Reviewer Lineage, and
  every diff inside the candidate is measured from its Effective Base. It
  refines ADR-0169 and amends ADR-0197 (a moved merge base no longer always
  restarts the lineage), ADR-0282 and ADR-0352 (how a Review Segment is
  measured) and ADR-0286 (where a Correction Review's delta starts); task_04
  writes the four amendment notices.
  Accounted for and unchanged, each one citation edge from a listed ADR on a part this Spec leaves alone: ADR-0092, ADR-0120, ADR-0165, ADR-0192, ADR-0196, ADR-0227, ADR-0233, ADR-0247, ADR-0248, ADR-0250, ADR-0251, ADR-0254, ADR-0256, ADR-0258, ADR-0260, ADR-0261, ADR-0262, ADR-0269, ADR-0272, ADR-0277, ADR-0278, ADR-0279, ADR-0281, ADR-0290, ADR-0295, ADR-0299, ADR-0306, ADR-0308, ADR-0317, ADR-0345, ADR-0348, ADR-0349. They keep the queue's own conflict merge without a repeated review (ADR-0192), finding validation (ADR-0196), the Claude budget (ADR-0227), and the history, archive, record, module and binary decisions the design never touches. Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AZ batch 2 answers of 2026-11-05 quoted in the authorization record; the
  governed files are the Roundfix skill files task_04 changes. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0356-a-main-merge-is-a-review-delta-and-a-segment-is-measured-from-its-effective-base/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/deliver.md`, `.agents/skills/roundfix/references/review.md`, `skills/roundfix/SKILL.md`. No Makefile, workflow, Go module, lint or formatter change.

## Goals

- A conflict-free merge of the default branch into an item after a review
  keeps its Reviewer Lineage: the next review is round 2 or the Correction
  Review, and it reads only the item's own change since the reviewed head.
- A merge that brings in another branch, a conflicting merge, a rebase, a
  rewritten head or a changed provider still starts round 1, and Roundfix
  names which of these happened.
- Before the queue reviews such a head, `deliver status` and `deliver retry`
  say that the review will open a new Reviewer Lineage at round 1 and why.
- Each Review Segment, and each commit inside it, is measured from its
  Effective Base, so a commit cut from an old default branch or a merge of the
  default branch never shows the default branch's work, and does not block
  the review by size.
- The `review` and `deliver` guides, the Roundfix `review` reference, the
  glossary and the amended ADRs say so.

## Core Features

1. The Effective Base measurement: one function that returns what a diff
   from one candidate commit to a later one starts from, carrying the earlier
   commit's own change onto the later one's Effective Base.
2. Segments measured from their Effective Base, for `codex` and `claude`
   alike, recording `measuredFrom` and saying so in the segment prompt.
3. The lineage across a conflict-free merge of the default branch, its
   `lineage.deltaBase` and prompt sentence, for round 2 and the Correction
   Review, and the restart reason when the lineage does not continue.
4. The `Review:` notice of `deliver status` and `deliver retry`.
5. The rule in writing: guides, reference, glossary, ADR-0366 and four
   amendment notices.

## Non-Goals / Out of Scope

- Fiscus points 1, 3 and 5: queue turn and review concurrency (Spec 0355),
  numbered-file collisions (Spec 0357), large authored data files (Spec 0358).
- The review bound, the Claude token budget, the four-segment limit and the
  text of every size refusal.
- Round 2 and the Correction Review stay unsegmented; a delta above the bound
  still blocks.
- The Delivery Retry's ancestry rules, the queue's own merge of the default
  branch on a conflicting Pull Request (ADR-0192) and the review record's
  reuse key.
- `roundfix review`'s own stdout and stderr, apart from the two optional
  record fields.
- The Makefile, the CI workflows, `go.mod`, `go.sum`, any Baseline module,
  and every other line of the Roundfix `deliver` reference.

## Success Metrics

1. Success Metric: in Fiscus's point-4 shape (a commit cut from the old
   default branch, then a merge of a default branch that added more than the
   bound, then more commits), `roundfix review --base main` with `codex` and
   with `claude` exits `0` with outcome `reviewed` and three segments, the
   first two recording `measuredFrom` and its prompt line, and no segment
   prompt holds the default branch's added file; at `ad77a797` the point-4
   shape of Validation exits `2` with `review diff too large: 1223758 bytes after omitting 0
   path(s) exceeds the review bound of 917504 bytes` for `codex`.
2. Success Metric: in Fiscus's point-2 shape (round-1 findings, a fix, then a
   conflict-free merge of the advanced default branch), the next review is
   round 2 whose prompt holds the fix and the `lineage.deltaBase` sentence and
   not the default branch's change; a round-2 fix merged the same way gets its
   Correction Review on the same terms; at `ad77a797` the first shape opens
   round 1 with the full candidate diff.
3. Success Metric: each of the six restart reasons of `_techspec.md` → API
   Contract 4 is returned, word for word, for its case, a conflict-free merge
   returns none, and every row of the existing lineage decision table keeps
   its round.
4. Success Metric: for an item parked after `review-findings` whose head was
   rewritten, `deliver status` and `deliver retry` print the `Review:` line of
   Surface Transcripts 1 and 2; after a conflict-free merge of the default
   branch they print none.
5. Success Metric: both copies of the Roundfix `review` reference and the
   `review` guide state the Effective Base rule and the lineage rule, the
   `deliver` guide and both copies of the `deliver` reference state the
   notice, `CONTEXT.md` defines **Effective Base**
   and its changed entries cite ADR-0366, and ADR-0197, ADR-0282, ADR-0286 and
   ADR-0352 carry the amendment notice naming ADR-0366.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The Fiscus report quoted in the adopted reference, written before this Spec
  existed: points 2 and 4, the Specs named there (0066, 0074, 0070) and the
  measured 1.5 to 2 MB range.
- Git's own merge, through the `git` binary the tests already use: a carry
  that Git reports as conflicting (exit status 1 of `git merge-tree
  --write-tree`) keeps today's measurement, so the oracle of a clean carry is
  Git's merge, not this Spec's reading of it. The git-merge-tree(1)
  documentation (https://git-scm.com/docs/git-merge-tree, read through Exa on
  2026-11-05): "For a successful, non-conflicted merge, the exit status is 0.
  When the merge has conflicts, the exit status is 1."

## Unreachable Acceptance

None. Every criterion is checked with this repository's files, temporary
homes and temporary repositories, and a fake reviewer.

## Validation

The decisions were checked on 2026-11-05 against these sources.

- Characterization at `ad77a797` with two temporary tests through the public
  `review` command and the fake reviewer of `internal/cli`:
  - the point-2 shape: round 1 raised one finding at base `B`; after a fix and
    `git merge --no-edit main` of an advanced default branch, the next review
    recorded `round: 1`, base the new default branch tip, and its prompt held
    `--- BEGIN CANDIDATE DIFF ---`, not `ROUND 2 DELTA`;
  - the point-4 shape: a default branch commit of 1,101,005 bytes merged after
    the fixture's own commit, then two commits of 611,670 bytes each: exit
    `2`, outcome `blocked`, no segments, `review diff too large: 1223758
    bytes after omitting 0 path(s) exceeds the review bound of 917504 bytes;
    largest files in the diff: part1.txt (611794 bytes), part2.txt (611794
    bytes), review.txt (170 bytes); ...`.
- A prototype of the design on a scratch worktree of `ad77a797`: both shapes
  pass (point 4 as two segments, point 2 as round 2 without the default
  branch's change), every test of `./internal/cli` passed unchanged
  (`go test -count=1 ./internal/cli`, 323 s), and each sabotage of
  `_techspec.md` → Testing Approach made its named test fail.
- ADRs read: ADR-0169, ADR-0192, ADR-0197, ADR-0227, ADR-0282, ADR-0286 and
  ADR-0352.
- Secondbrain, read only: `wiki/index.md`, then `qmd query "roundfix pre-PR
  review lineage merge base merge main restart round 1 segment first-parent
  revert" --all --files --min-score 0.3`. It returned the Roundfix mirrors of
  `CONTEXT.md`, the `usage`, `commands` and `deliver` guides, ADR-0197 and
  ADR-0286, and a capture on comparing with the merge base (`git diff A...B`);
  none decides how a merge of the default branch is reviewed.
- Exa: git-merge-tree(1) at git-scm.com, on the `--write-tree` output and exit
  status.

## Glossary

- adds: **Effective Base**
- changes: **Reviewer Lineage**
- changes: **Review Segment**

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AZ batch 2 brief for Specs 0356 and 0358 of 2026-11-05 and the maintainer's answers below.

The maintainer's answers (AskUserQuestion, 2026-11-05 ~08:20 BRT), verbatim:

- Question: "Relato do Fiscus (fila de 21 Specs há 2 dias). Quais pontos
  viram Spec agora? Os que ficarem de fora viram Backlog Entry." — all four
  options selected; the two this Spec answers:
  - "2. Merge do main (Recommended)" — option text: "Um merge do main sem
    conflito, em que o único pai novo é o main, conta como delta e não
    reinicia o review na rodada 1. Quando o review for reiniciar, o retry
    avisa antes."
  - "4. Trecho do review (Recommended)" — option text: "Cada trecho do review
    passa a ser medido contra o merge-base com o main atual. Um commit antigo
    que só parece reverter o main deixa de bloquear o review por tamanho."
- Rules (all): "Mesmas do AZ (Recommended)" — option text: "Sem supervisão,
  com merge pela fila e conserto quando a entrega parar. Makefile, CI, go.mod
  e lint liberados, assim como mudanças de Baseline e skills que chegam aos
  adotantes. Release assim que cada Spec ficar pronta, com notas aos
  adotantes, aviso ao Fluxus e ao Fiscus."
- Standing instruction of 2026-11-05: "Quero que faça os releases hoje assim
  que estiverem prontos. E se existe no baseline regra para um release por dia
  por algum motivo deve ser removida." There is no such rule in the
  repository; this Spec changes nothing for it.

Taken from those answers, the operator's brief and the characterization, with
no further interview:

- "Only new parent is the default branch" is checked on every merge of the
  head's first-parent path since the reviewed head: each parent after the
  first must be contained in the new merge base.
- "Without conflict" is checked by carrying the reviewed head onto the new
  merge base with `git merge-tree --write-tree`; the delta starts from that
  tree, so the fix commits before the merge and any resolution in it are
  reviewed, and the default branch's work is not.
- "The merge-base of the commit with the current default branch" equals its
  merge base with the review's merge base, so no second ref is read.
- The notice is printed by `deliver status` for parked or reviewing items and
  by `deliver retry` when the item returns to `reviewing`; it never blocks.
- ADR-0366 records the decision, refines ADR-0169 and amends ADR-0197,
  ADR-0282, ADR-0286 and ADR-0352. Four implementation Tasks in series and a
  QA gate.

## Open Questions

None.
