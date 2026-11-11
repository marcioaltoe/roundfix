---
spec: 0381-a-conflicting-merge-of-the-default-branch-is-a-review-delta
status: active
created: 2026-11-11
surfaces: [cli, docs]
type: fix
---

# A conflicting merge of the default branch is a review delta

On 2026-11-11 the pre-PR review of `roundfix deliver` (Roundfix 0.138.0, provider `codex`) raised two false compile findings in Conexus, each claiming that a `catch` or an import was missing although it sat outside the hunk and the head type-checked, and each right after the operator merged `main` into the item branch to clear a `pull-request-conflict` park. A reproduction proved why that review was round 1 over the whole candidate: the reviewed head conflicts with the default branch, which is what the park means, and ADR-0366 restarts the Reviewer Lineage whenever that carry conflicts; in a segmented review the same rule also shows the default branch's newly landed work in a segment that spans the merge. After this Spec a conflicting merge of the default branch keeps the lineage, the round-2, Correction Review and segment diffs start from the Conflicted Carry so they hold only the candidate's own change and its resolution, and every review prompt asks the reviewer to read the whole file before a finding says that something outside the diff is missing.

## Sources

One source, adopted (see `references/_index.md`):

1. [`references/2026-11-11-conexus-review-false-compile-findings-after-a-main-merge.md`](references/2026-11-11-conexus-review-false-compile-findings-after-a-main-merge.md),
   the ai-memory message `01a1cb37-3583-7b31-ba72-e9972d6e77f0` of
   2026-11-11, reported by conexus, nature `measured`, with the operator's
   reproduction note.

At the authoring base `docs/backlog/` and `docs/findings/` hold no open entry
about the pre-PR review, so no other source shares this context.

## Prerequisites

None in `requires`. Authored against `013af3b4` (v0.139.0). The queue
delivers Specs 0377, 0376, 0378, 0379 and 0380 before this one and Spec 0382
after it; none of their code is needed here. This Spec owns the review path:
`internal/cli/review.go`, `review_lineage.go`, `review_effective_base.go`,
`review_segments.go`, their tests, the `review` reference and guide, and the
`deliver` reference and guide where they count restart reasons. It shares,
and edits only through the repository's record commands or narrowly:

- the owned Roundfix skill version pair (`.agents/skills/roundfix/SKILL.md`,
  `skills/roundfix/SKILL.md`) and `skills/testdata/owned-skill-versions.json`,
  which 0376, 0380 and 0382 also raise; task_01 writes "the next version"
  through the record command, never a number;
- `docs/references/behavior-surfaces.json`, re-recorded by its command;
- `CONTEXT.md`, where task_01 adds one entry after **Effective Base** and
  rewords that entry and **Reviewer Lineage**; 0376, 0380 and 0382 edit other
  entries.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted or changed; the review record keeps its fields and
  `lineage.deltaBase` keeps naming a tree. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the change reads the local Git
  repository only and never fetches. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0366 decides the Effective Base
  measurement and restarts the lineage on a conflicting carry, and ADR-0197
  bounds the lineage at two rounds and records that a conflicting merge
  starts round 1; this Spec records ADR-0391, which amends both, with
  "Amended by" notices task_01 writes. Not applicable, each cited by one of
  those or by ADR-0391, with their rules unchanged: ADR-0169 (the merge base
  moves with the candidate's merge), ADR-0282 (Review Segment ranges and the
  bound), ADR-0286 (Correction Review eligibility), ADR-0352 (the Claude
  budget and omit paths), ADR-0233 (the retry path rule) and ADR-0272 (the
  review record surviving a fix). Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — no protected tooling mutation proposed or
  authorized: no Makefile, workflow, `go.mod`, `go.sum`, lint or formatter
  change. The Governed Paths task_01 touches are the owned Roundfix skill
  files and the canonical `review` and `deliver` references, measured with
  `speccheck.GovernedPath`, which the maintainer's cycle grants cover.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0381-a-conflicting-merge-of-the-default-branch-is-a-review-delta/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/review.md`,
  `skills/roundfix/SKILL.md`.

## Goals

- The review after a merge of the default branch that conflicts, the
  recovery every `pull-request-conflict` park asks for, continues the
  Reviewer Lineage: round 2, or the Correction Review, of the delta.
- That delta, and every Review Segment that spans the merge, holds only the
  candidate's own change and its conflict resolution, never the default
  branch's work.
- A reviewer reads the file before it reports that an import, declaration,
  symbol or the end of a block is missing.

## User Stories

1. As an operator who merged `main` into a parked item to resolve a
   conflict, I want the next review to read only what changed since the
   reviewed head, so that the item is not reviewed again in full.
2. As the maintainer, I want a segmented review never to show the default
   branch's landed Specs as the candidate's, so that its findings are about
   the candidate.
3. As an operator running an unattended queue, I want the reviewer to check
   a missing import or `catch` in the file before reporting it, so that a
   false finding does not park the item as `corrective-spec-required`.

## Core Features

1. When `git merge-tree --write-tree` of a candidate commit onto an Effective
   Base exits 1, the review measures from the tree Git wrote, read with
   `merge.conflictStyle=merge`, and knows the conflicted paths Git listed.
2. `decideReviewLineage` no longer restarts on a conflicting carry; the
   restart reason `the merge base moved from <prior base> to <base> and the
   reviewed head <prior head> conflicts with it` is retired, and seven
   reasons remain.
3. The round-2 and Correction Review prompts, after the line naming the
   delta's tree, and a segment prompt, after its `This segment's diff starts
   from` line, add the Conflicted Carry lines of `_techspec.md` → API
   Contract 2.
4. Every review prompt adds the absence-claim sentence of `_techspec.md` →
   API Contract 3 after its read-only line.
5. The `review` and `deliver` guides, the Roundfix skill references and the
   glossary describe the Conflicted Carry, and task_01 writes the amendment
   notices of ADR-0391.

## User Experience

Before this Spec, in the reported shape, `roundfix deliver retry` of the
resolved item prints `Review: <slug> head <head> opens a new Reviewer Lineage
at round 1, a full review of the candidate: the merge base moved from <prior
base> to <base> and the reviewed head <prior head> conflicts with it`, and
the reviewer receives `--- BEGIN CANDIDATE DIFF ---` with the item's whole
change. After it, the retry prints no `Review:` line, and the reviewer
receives `--- BEGIN ROUND 2 DELTA DIFF ---` with the round-1 findings, the
line `The candidate merged the default branch since then; the diff starts
from tree <deltaBase>, …`, the line `Carrying <prior head> onto the default
branch conflicted; where tree <deltaBase> holds Git's conflict markers, the
diff shows how the candidate resolved the conflict.`, the line `Conflicted
paths: <path>.`, and a diff whose conflicted files show Git's markers
removed and the resolution. The review's exit codes and record fields are
unchanged.

## Non-Goals / Out of Scope

- Running the repository's typecheck or build inside the review, or giving
  its result to the reviewer (the reporter's third suggestion): the review
  runs no repository command and serves every stack, the Repository
  Verification already builds the head before the Pull Request, and both
  false findings were claims the file at the head refutes, which the
  read-only session can open.
- Sending whole files to the reviewer (the first suggestion, in full): it
  spends the review bound of ADR-0282 on unchanged lines; the absence-claim
  sentence asks for the read instead.
- A mechanical check that refuses a finding claiming a missing symbol: it
  would be language-specific and judge text the validation does not parse.
- Any change to the two-round ceiling, the Correction Review's eligibility,
  the other seven restart reasons, the record and ledger formats, `review
  dispose`, or the Delivery Queue's conflict park and derived-path
  resolution.
- Any change to the Makefile, workflows, `go.mod`, `go.sum`, lint or
  formatter configuration or the Baseline.

## Success Metrics

1. Success Metric: in a temporary repository of the reported shape, the
   review after a conflicting merge of the default branch exits 0 at round 2
   with `previousHead` set to the round-1 head and `lineage.deltaBase` equal
   to the Conflicted Carry, and its prompt holds neither the candidate diff
   nor any of the thirty files the default branch landed; on the authoring
   base it is round 1.
2. Success Metric: a Correction Review after a conflicting merge reads only
   the fix and the resolution, and a Review Segment spanning a conflicting
   merge is measured from the Conflicted Carry without the default branch's
   files; on the authoring base the first restarts at round 1 and the second
   includes those files.
3. Success Metric: every review prompt carries the absence-claim sentence
   right after its read-only line; on the authoring base none does.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The adopted Conexus message, measured on a repository this Spec did not
  build (`useconexus/conexus`, 2026-11-11, Roundfix 0.138.0): both findings
  came in "round 1" right after "mesclar a `main` na branch do item, como
  pede o park `pull-request-conflict`", and each "cita exatamente um símbolo
  definido fora do trecho alterado".
- The Git manual, read through Exa and the installed manual on 2026-11-11:
  `git merge-tree` writes a tree that "represents what would be checked out
  in the working tree at the end of git merge. If there were conflicts, then
  files within this tree may have embedded conflict markers", and "When the
  merge has conflicts, the exit status is 1"
  (https://git-scm.com/docs/git-merge-tree); `git log --remerge-diff`
  "Remerge[s] two-parent merge commits to create a temporary tree
  object—potentially containing files with conflict markers and such. A diff
  is then shown between that temporary tree and the actual merge commit."
  (https://git-scm.com/docs/git-log).

## Unreachable Acceptance

None.

## Validation

The decisions were checked on 2026-11-11 against these sources.

- Reading at `013af3b4`: `decideReviewLineage`
  (`internal/cli/review_lineage.go`) restarts with `… conflicts with it` when
  `reviewMeasuredFrom` (`internal/cli/review_effective_base.go`) reports a
  conflicting carry; that function returns the earlier commit on `git
  merge-tree` exit 1; `planReviewSegments` (`internal/cli/review_segments.go`)
  measures each segment through it; `buildReviewPrompt`
  (`internal/cli/review.go`) tells the read-only reviewer it "may open
  repository files for context" and never asks it to confirm an absence; the
  `pull-request-conflict` park (`internal/delivery/park_class.go`) tells the
  operator to merge the default branch, and the retry returns the item to
  `reviewing`.
- A temporary-repository reproduction with the review command's fixtures and
  a fake reviewer on `013af3b4`: after a conflicting merge the review was
  round 1 with the conflict restart and a `BEGIN CANDIDATE DIFF` holding the
  item's change but none of the thirty landed files; the `catch` 30 lines
  below the changed line was outside the hunk; a segment measured from an
  item commit across the merge carried the thirty landed files.
- A prototype of the design in a scratch worktree of `013af3b4`: the new
  tests fail on the authoring base and pass with the change; six sabotages
  each fail their named test; every `Review|Lineage|Dispos|Dismiss|
  Correction|Deliver|Segment` test of `./internal/cli` passes once the two
  pinned expectations this Spec moves are updated, and `GOOS=windows go vet
  ./internal/cli` passes.
- Published sources (Exa and the Git 2.56.0 manual): the passages quoted
  above. Git's own "remerge" diff is the precedent: a diff from Git's re-merge
  with conflict markers to the recorded result shows only the resolution.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd query "roundfix pre-PR review lineage merge default branch conflict round 1 false finding" --all --files --min-score 0.3`,
  which returned no file above the score; no wiki page decides how a review
  treats a conflicting merge, so nothing there supports or contradicts this
  design.

## Glossary

- adds: **Conflicted Carry**
- changes: **Effective Base**
- changes: **Reviewer Lineage**

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle BC brief for Spec 0381 of 2026-11-11, the adopted Conexus message, the reproduction, and the maintainer's answer below.

The maintainer's answer (AskUserQuestion via `/roundfix --follow-up`,
2026-11-11 ~14:45–14:50 BRT), verbatim:

- Question "Specs novas D (mensagens: pergunta não bloqueante recebe
  resposta; precedência de 'bloqueado' limitada) e E (review pré-PR revisa só
  o delta e com contexto depois de merge da main): como autorizar?" → "Após a
  0379 (Recommended)" — option text: "Mesmas autorizações do ciclo BB (push,
  merge, release, teto minor); entregues depois da fila 0377→0379, D antes de
  E."

Taken from those and the evidence, with no further interview:

- The cause is ADR-0366's conflict fallback, not the size of the diff: the
  unsegmented round-1 diff never held the default branch's work, so the
  reporter's inference holds only for segments. Both effects share one fix.
- A conflicting carry starts from the tree `git merge-tree --write-tree`
  writes, which the Git manual defines as the merge result with conflict
  markers; the diff from it is the candidate's resolution plus its own
  change. `merge.conflictStyle=merge` is pinned so the tree, and
  `lineage.deltaBase`, do not depend on a user's Git configuration.
- `preflight.GitRunner` folds standard output into its error on exit 1, so
  the conflicted tree is read by a second, direct `git merge-tree` call, as
  `review_recorded_commit.go` already runs Git directly; the first call
  through the runner still decides clean, conflicted or failed, so every
  injected-failure test keeps its meaning.
- The ceiling of ADR-0197 is unchanged: a conflicting merge after a clean
  round 2 closes at the ceiling without review, as every later change does.
- The reviewer gets one sentence, not whole files: it already runs read-only
  in the checkout at the head, and the two findings were absence claims the
  file refutes; the maintainer approved "com contexto".
- Nothing is persisted for the conflict: the prompt lines are recomputed from
  `lineage.deltaBase` and the commits at every run, so the record schema
  stays as it is.
- ADR-0391 records the decision and amends ADR-0366 and ADR-0197; task_01
  writes their "Amended by" notices.
- Two implementation Tasks and the QA gate. `type: fix`; the release is the
  operator's patch after the merge; no Task releases.

## Open Questions

None.
