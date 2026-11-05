---
spec: 0358-a-large-authored-data-file-is-reviewed-apart-not-omitted
status: active
created: 2026-11-05
surfaces: [backend, docs]
type: feat
---

# A large authored data file is reviewed apart, not omitted

Reported by fiscus on 2026-11-05, point 5 of its report: an authored catalog
of fiscal rules, `icms-es-v2.json` at about 550 KB, passed the pre-PR review
bound on its own, and the only way through was to declare it in
`delivery.derived_paths`, so the reviewer stopped reading the fiscal rules.
After this Spec a repository lists such files in a new Project Config key,
`pre_pr_review.data_paths`; the review takes them out of the code diff and
reviews their diff apart, in up to four Data Review Segments of their own in
the same review session, so the code review fits its bound and the data is
still read.

## Sources

One report, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-05-fiscus-the-delivery-queue-takes-days-and-a-large-data-catalog-is-omitted-from-review.md`](references/2026-11-05-fiscus-the-delivery-queue-takes-days-and-a-large-data-catalog-is-omitted-from-review.md),
   the Fiscus ai-memory message `01a1ab64` of 2026-11-05 08:07 BRT, reported
   by fiscus. This Spec answers its point 5; Specs 0355, 0356 and 0357 answer
   points 1 to 4 from their own copies.

At the authoring base `docs/backlog/` and `docs/findings/` hold no open entry
on the pre-PR review scope, so no other source shares this context.

## Prerequisites

None in `requires`; no Verification of this Spec depends on another Spec.
Authored against `ad77a797`. The planned queue order is 0354, 0356, 0355,
0357, 0358, so Spec 0356 will have merged when this Spec runs. Files and
functions this Spec shares, each edited narrowly so it rebases cleanly:

- `internal/cli/review.go` and `internal/cli/review_segments.go`, which Spec
  0356 changes first (`planReviewSegments`, the lineage in `runReviewCommand`
  and the round-2 delta in `runConfiguredReviewSession`). This Spec adds the
  `data-path` reason in `reviewScopedDiff`, one line in
  `appendReviewOmissions`, the Data Review Segment step in
  `runConfiguredReviewSession` before the readiness proof, the answer headers
  in its prompt loop and the `dataSegment` label in
  `classifyReviewSegmentAnswers`; the data diff starts from that function's
  `measuredBase`, whatever Spec 0356 makes it.
- `internal/config/config.go` and `internal/config/pre_pr_review.go`, next to
  `pre_pr_review.omit_paths`.
- The owned Roundfix skill version (`.agents/skills/roundfix/SKILL.md`, its
  mirror and `skills/testdata/owned-skill-versions.json`), the Roundfix
  `review` reference, `docs/user-guide/commands/review.md`, `CONTEXT.md`,
  `docs/references/behavior-surfaces.json` and ADR-0282, which Spec 0356 also
  changes.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or schema version
  is minted; one Project Config key and two optional record fields are added.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential; the change reads local Git objects and Project Config. Source:
  `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0366 (a conflict-free merge of the default branch is a review delta, a sibling Spec in the same queue) touches the same review code but carries no obligation for this Spec; ADR-0368 (new) records the decision:
  a declared data file is reviewed apart in Data Review Segments. It amends
  ADR-0282, which rejected splitting the diff by file or by bytes, for
  declared data files only, and refines ADR-0352, whose `omit_paths` keeps its
  meaning; task_04 writes the amendment notice into ADR-0282.
  Accounted for and unchanged, each one citation edge from a listed ADR on a part this Spec leaves alone: ADR-0092, ADR-0120, ADR-0165, ADR-0169, ADR-0192, ADR-0196, ADR-0197, ADR-0227, ADR-0233, ADR-0247, ADR-0248, ADR-0250, ADR-0251, ADR-0254, ADR-0256, ADR-0258, ADR-0260, ADR-0261, ADR-0262, ADR-0269, ADR-0272, ADR-0277, ADR-0278, ADR-0279, ADR-0281, ADR-0286, ADR-0290, ADR-0295, ADR-0299, ADR-0306, ADR-0308, ADR-0317, ADR-0345, ADR-0348, ADR-0349. They keep the Claude budget this Spec applies to each Data Review Segment (ADR-0227), the queue's conflict regeneration (ADR-0192), and the history, archive, record, module and binary decisions the design never touches. Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AZ batch 2 answers of 2026-11-05 quoted in the authorization record; the
  governed files are the Roundfix skill files task_04 changes. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0358-a-large-authored-data-file-is-reviewed-apart-not-omitted/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/review.md`, `skills/roundfix/SKILL.md`. No Makefile, workflow, Go module, lint or formatter change.

## Goals

- A repository declares its large authored data files in
  `pre_pr_review.data_paths`, at its default branch, with the same path forms
  and validation as `pre_pr_review.omit_paths`.
- The review takes those files out of every code diff and reviews their diff
  apart, after the code prompts, in the same session, in every round.
- A data file larger than the bound is still reviewed: it is cut at hunks and
  lines into at most four Data Review Segments, each within the provider's
  bound, and every finding anchors to a line of the file.
- `delivery.derived_paths`, `pre_pr_review.omit_paths` and a review without the
  key behave exactly as before.
- The configuration and `review` guides, the Roundfix `review` reference, the
  glossary and the ADRs say so.

## Core Features

1. The Project Config key `pre_pr_review.data_paths` and its matcher.
2. The `data-path` reason: a matched file leaves the code diff and is named in
   the code prompt as reviewed apart; an omission reason wins.
3. The Data Review Segment plan: packing by file, cutting by hunk and line,
   the bound per provider and the four-segment limit with its refusal.
4. The review of the segments in the session, their findings, record fields
   and answer file.
5. The rule in writing: guides, reference, glossary, ADR-0368 and the
   amendment notice in ADR-0282.

## User Stories

- As a maintainer of a repository with a large authored rules catalog, I
  declare it in `pre_pr_review.data_paths` so the pre-PR review reads its
  values instead of skipping it, and the item is not parked `review-blocked`
  for its size.
- As the operator of the Delivery Queue, I read in the review record and
  answer file which findings came from the data files, and dispose them like
  any other finding.

## Non-Goals / Out of Scope

- A summary or parser of any data format; the reviewer reads the diff.
- Changing the meaning of `delivery.derived_paths` or
  `pre_pr_review.omit_paths`, the review bound, the Claude budget, the Review
  Segment limit or any existing refusal text.
- Fiscus points 1 to 4 (Specs 0355, 0356, 0357).
- The Makefile, the CI workflows, `go.mod`, `go.sum`, any Baseline module, the
  `init` templates.

## Success Metrics

1. Success Metric: with `pre_pr_review.data_paths: ['rules/*.json']` at the
   base commit and a candidate adding `code.go` and a 9,000-entry
   `rules/icms.json`, `roundfix review` with `codex` and with `claude` sends
   two prompts: a code prompt without the catalog that names it as reviewed
   apart, and `Data Review Segment 1 of 1` with it; the catalog's finding
   carries `dataSegment: 1`; at `ad77a797` the same Project Config is refused
   with `field data_paths not found in type config.prePRReviewOverlay`.
2. Success Metric: a 24,000-entry data file above the `codex` bound becomes
   two Data Review Segments, each within 917,504 bytes, a finding on a line of
   the second stands its anchor validation, and every added line of every
   piece sits at the line number where the head's file holds it.
3. Success Metric: five data files that need five Data Review Segments block
   the review before any session with `review data too large: <bytes> bytes of
   pre_pr_review.data_paths needs 5 data review segments; a review allows at
   most 4` and the size suffix.
4. Success Metric: a file matched by `pre_pr_review.omit_paths` and
   `pre_pr_review.data_paths` stays omitted as `omit-path`, and every existing
   review, scope and segment test passes unchanged.
5. Success Metric: the configuration guide documents the key, both copies of
   the Roundfix `review` reference and the `review` guide describe Data Review
   Segments, `CONTEXT.md` defines **Review Data Path** and **Data Review
   Segment**, and ADR-0282 carries the amendment naming ADR-0368.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The Fiscus report quoted in the adopted reference, written before this Spec
  existed: point 5, the catalog `icms-es-v2.json` of about 550 KB and Fiscus
  pull request #263, which declared it in `delivery.derived_paths`.
- Git's own content, through the `git` binary the tests already use: every
  `+` line of a piece cut from a data file sits, by its recomputed hunk
  header, at the line number where `git show <head>:<path>` holds that text,
  so the oracle of a valid cut is the file Git stores, not this Spec's
  reading of the format. The git-diff(1) documentation
  (https://git-scm.com/docs/git-diff, read through Exa on 2026-11-05) on the
  input those pieces come from: running git-diff "with the `-p` option
  produces patch text".

## Unreachable Acceptance

None. Every criterion is checked with this repository's files, temporary
homes and temporary repositories, and a fake reviewer.

## Validation

The decisions were checked on 2026-11-05 against these sources.

- Characterization at `ad77a797` with the built binary, a temporary `HOME` and
  a temporary repository whose Project Config declares
  `pre_pr_review.data_paths: [rules/]`: `roundfix review` exits `2` with
  `roundfix: review blocked: parse config "<repo>/.roundfixrc.yml": yaml:
  unmarshal errors:` and `line 3: field data_paths not found in type
  config.prePRReviewOverlay`.
- Reading at `ad77a797`: `reviewScopedDiff` omits a `delivery.derived_paths`
  match as `generated` and a `pre_pr_review.omit_paths` match as `omit-path`,
  both left out of the reviewer's diff; `planReviewSegments` cuts only by
  first-parent commits, so one commit that adds a file above the bound blocks.
- A prototype of the design on a scratch worktree of `ad77a797`: the four
  command tests of `_techspec.md` → Testing Approach passed for `codex` and
  `claude`, and `./internal/cli`, `./internal/config` and `-tags docscontract
  ./internal/docscontract` passed.
- ADRs read: ADR-0192, ADR-0227, ADR-0282 and ADR-0352.
- Secondbrain, read only: `wiki/index.md`, then `qmd query "roundfix pre-PR
  review lineage merge base merge main restart round 1 segment first-parent
  revert" --all --files --min-score 0.3`, which returned the Roundfix mirrors
  of the review guides and ADRs and no page on reviewing data files.
- Exa: git-diff(1) and git-merge-tree(1) at git-scm.com.

## Glossary

- adds: **Review Data Path**
- adds: **Data Review Segment**

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AZ batch 2 brief for Specs 0356 and 0358 of 2026-11-05 and the maintainer's answers below.

The maintainer's answers (AskUserQuestion, 2026-11-05 ~08:20 BRT), verbatim:

- Question: "Relato do Fiscus (fila de 21 Specs há 2 dias). Quais pontos
  viram Spec agora? Os que ficarem de fora viram Backlog Entry." — all four
  options selected; the one this Spec answers, in part:
  - "3+5. Migrations e catálogos" — option text: "Um detector de número de
    migration repetido entre Specs da fila, rodado antes da Run. Uma
    categoria nova para arquivos grandes de dados escritos pelo projeto,
    revisados à parte em vez de omitidos. São desenhos maiores." This Spec
    answers the second half; Spec 0357 answers the first.
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

- "Revisados à parte" is a separate prompt sequence in the same session, not a
  summary: the reviewer must read the values.
- The key lives under `pre_pr_review`, next to `omit_paths`, with the same
  path forms, the same base-commit reading and the same Project-over-User
  replacement.
- Data Review Segments run in every round, from the start of that round's code
  diff, so a round-2 change to the catalog is also read apart.
- Each segment fits the provider's existing bound; the four-segment limit
  mirrors Review Segments; the refusal reuses the actionable size suffix.
- ADR-0368 records the decision, amends ADR-0282 for data files only and
  refines ADR-0352. Four implementation Tasks in series and a QA gate.

## Open Questions

None.
