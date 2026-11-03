---
spec: 0342-the-claude-review-segments-and-omits-declared-paths
status: active
created: 2026-11-03
surfaces: [backend, cli, docs]
type: feat
---

# The Claude review segments a large candidate and omits declared paths

A repository whose pre-PR review runs on Claude loses a delivery whenever one
large generated file lands in a candidate. The review refuses the prompt by
size before any reviewer reads it, the Delivery Queue parks the item as
`review-blocked`, and the refusal names neither the file nor a remedy; the
only way to drop the file today is a derived-path declaration with a
regeneration command that does nothing. After this Spec a Project Config key
lists paths the review leaves out, recorded with their own omission reason;
Claude reviews a large round-1 candidate in up to four Review Segments, each
within its token budget, as Codex already does by bytes; and a review that
still does not fit names the three largest files of its diff and that key, in
the record, on the command's standard error and in the queue's park blocker.

## Sources

One adopter report, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-03-the-claude-pre-pr-review-blocks-on-size-instead-of-segmenting.md`](references/2026-11-03-the-claude-pre-pr-review-blocks-on-size-instead-of-segmenting.md),
   `type: feat`, reported by vortex, ai-memory message `01a1a0ea`.

At the authoring base `docs/backlog/` and `docs/findings/` hold no entry, so
no other source shares this context.

## Prerequisites

None in `requires`. Authored against `4e61bf6d` (v0.129.1), where no other
Spec is active. Specs 0343 and 0344 are authored in parallel and deliver after
this one; `_techspec.md` → Build Order names the files they may share.

## Project Constraints

- Identifier strategy: not applicable — no identifier or slug is minted; the
  Project Config gains one optional list key, the Pre-PR Review Record one
  omission reason value and one optional per-segment field. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the review reads local Git
  objects and sends its prompt through the existing Agent Session; no request
  path, endpoint or credential changes. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0227 bounds the Claude prompt in
  tokens, ADR-0227: "The Claude review prompt is bounded in tokens, not
  bytes"; this Spec keeps that budget and applies it to each Review Segment.
  ADR-0282 splits a large `codex` candidate into Review Segments and reads
  derived declarations only from the base commit, ADR-0282: "the candidate's
  own Project Config is not used, because a candidate must not choose what its
  review skips"; the new key follows the same rule, and Claude joins the
  segmentation. ADR-0352,
  written by this Spec, records the key, the Claude segments and the
  actionable refusal. The accepted ADRs one citation edge from those (each
  cites a listed ADR) are accounted for without change: the review lineage
  and record decisions ADR-0165, ADR-0169, ADR-0196, ADR-0197, ADR-0233, ADR-0272, ADR-0279, ADR-0286 and ADR-0317, whose two-round lineage, unsegmented
  round 2 and Correction Review and validated findings this Spec keeps; the
  history and archive decisions ADR-0092, ADR-0120, ADR-0247, ADR-0248, ADR-0251, ADR-0254, ADR-0258, ADR-0260, ADR-0261, ADR-0262, ADR-0269, ADR-0277, ADR-0278, ADR-0295, ADR-0299, ADR-0306, ADR-0308 and ADR-0348, which this Spec does not
  touch; and ADR-0250, ADR-0256, ADR-0281, ADR-0290, ADR-0345 and ADR-0349, which keep the module record, the release
  wait for lagging skills, executable rebuilds, item binary reclamation,
  the suite guard's sanctioned regeneration and the Release Ceiling this
  Spec's record uses. Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AW answers of 2026-11-03 quoted in the authorization record, which keep the
  cycle AV rules and the cycle AJ grant for the Governed Paths each Spec
  declares; the governed files are the Roundfix skill's `SKILL.md` pair and
  its canonical review and deliver references. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0342-the-claude-review-segments-and-omits-declared-paths/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/deliver.md`, `.agents/skills/roundfix/references/review.md`, `skills/roundfix/SKILL.md`. No Makefile, workflow, Go module, lint or formatter change.

## Goals

- A repository leaves a generated file out of the pre-PR review with one
  Project Config key, without declaring it derived and without a regeneration
  command, and the record says the file was left out and why.
- A Claude round-1 candidate whose prompt is above the budget only because
  several commits add up is reviewed in segments instead of refused.
- A review that still does not fit says which files are largest and which key
  would leave them out, where the operator first reads it.
- A Codex review, a round-2 review and a Correction Review behave as today,
  and a review without the key sends the same prompt bytes as today.

## Core Features

1. Review Omit Paths. Project Config `pre_pr_review.omit_paths` lists exact
   paths, directories ending in `/` and whole-path glob patterns. Roundfix
   refuses an unsafe or malformed entry with the key and the entry, as it
   refuses a derived-path entry. The review reads the key from the Project
   Config at the candidate's base commit, leaves every changed file it
   matches out of the diff of every round, and lists the file in the record's
   `omittedPaths` with reason `omit-path`, after the existing reasons. The
   prompt's omission summary counts these files only when there are any.
   Conflict regeneration keeps reading only the derived-path declarations.
2. Claude Review Segments. A `claude` round-1 prompt estimated above the
   budget is split into at most four consecutive first-parent commit ranges,
   each reviewed in order in the same session, each prompt as sent within the
   budget; the first carries the Spec contexts. One record merges the
   segments' findings as today. A commit too large alone, a fifth segment,
   and a round-2 or Correction Review prompt above the budget still block.
3. Actionable size refusal. Every size refusal — the byte bound, the token
   budget, or the segment limit, for either provider and every round — keeps
   its existing text and appends the three largest files of the diff it
   measured, with each file's bytes, and names `pre_pr_review.omit_paths` as
   the key that leaves them out.
4. The park says why. The Delivery Queue parks a blocked review as
   `review-blocked: <first line of the record's reason>`, so the queue's
   status carries the refusal; its Park Class stays `review`.
5. The rules in writing. The review, configuration and deliver guides, the
   Roundfix skill's review and deliver references, `CONTEXT.md`, ADR-0352 and
   dated amendment lines on ADR-0227 and ADR-0282 state them.

## Non-Goals / Out of Scope

- Omitting large files automatically when the prompt is over the budget
  (the report's third suggestion): a reviewer must never skip what nobody
  chose to skip.
- Detecting generated files from `.gitattributes` or file names.
- Reading the key from User Config, from the candidate's head or from an
  intermediate segment head.
- Changing conflict regeneration, the derived-path declarations, the Codex
  byte bound, the Claude budget, the four-segment limit or the
  two-round Reviewer Lineage.
- Segmenting round 2 or the Correction Review.
- The Makefile, the CI workflows, `go.mod`, `go.sum`, lint or formatter
  configuration, and any `### QA settlement` section.

## Success Metrics

1. Success Metric: with `pre_pr_review.omit_paths` at the base commit listing
   a glob, a candidate that changes a matching file and an authored file
   reviews only the authored file, and the record lists the matching file
   with reason `omit-path`; the same key only in the candidate's head omits
   nothing.
2. Success Metric: an unsafe entry, a malformed pattern and a non-list value
   are each refused with an error naming `pre_pr_review.omit_paths` and the
   entry; a review whose base has no key sends a prompt whose omission
   summary is byte-identical to today's.
3. Success Metric: a `claude` round-1 candidate of three commits whose whole
   prompt is above the budget and whose first two commits fit one prompt is
   reviewed in exactly two segments in one session, each segment's recorded
   estimate within the budget; five such segments block with the segment
   limit named in tokens.
4. Success Metric: a `codex` candidate's segments, records and prompts are
   unchanged, and a `claude` round-2 prompt above the budget still blocks.
5. Success Metric: every size refusal ends with the three largest files of
   its diff in descending bytes, ties by path, including a path with a space
   and non-ASCII letters written as Git lists it, followed by the key; fewer
   than three files lists what exists.
6. Success Metric: a blocked review parks its item as
   `review-blocked: <reason line>`, and a blocked review without a reason
   parks as `review-blocked`, both classified `review`.
7. Success Metric: the review, configuration and deliver guides, both copies
   of the Roundfix review and deliver references and `CONTEXT.md` describe
   the key, the Claude segments and the refusal, and the mirrors equal their
   canonical copies.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The Vortex report quoted verbatim in the adopted reference, written by a
  repository this Spec did not build before this Spec existed: a 541,455-token
  Claude refusal of a 1,049,480-byte diff in which 669 KB were one generated
  Drizzle snapshot, and a `derived_paths` workaround with a regeneration
  command that exits 1.
- GitHub's documentation "Customizing how changed files appear on GitHub"
  (https://docs.github.com/en/repositories/working-with-files/managing-files/customizing-how-changed-files-appear-on-github,
  read through Exa on 2026-11-03): "To keep certain files from displaying in
  diffs by default, or counting toward the repository language, you can mark
  them with the linguist-generated attribute in a .gitattributes file." A
  repository-declared list of generated paths kept out of a review diff is
  established practice.
- Anthropic's Claude Agent SDK documentation "How the agent loop works"
  (https://code.claude.com/docs/en/agent-sdk/agent-loop, read through Exa on
  2026-11-03): "When the context window approaches its limit, the SDK
  automatically compacts the conversation: it summarizes older history to
  free space, keeping your most recent exchanges and key decisions intact."
  Sending Claude segments in order to one session relies on this, as Codex
  segments rely on their runtime.

## Unreachable Acceptance

None. Every criterion is checked with this repository's own files,
temporary repositories and fake Agent runners.

## Validation

The decisions were checked on 2026-11-03 against these sources.

- Reading at `4e61bf6d`: the review segments a round-1 candidate only for
  `codex` (or a Correction Review, which round 2 never segments) above the
  917,504-byte bound; for `claude` it refuses a prompt estimated, at two bytes
  per token, above 500,000 tokens after the Spec contexts are appended, with
  no segmentation; omissions read the derived-path declarations of the base
  commit's Project Config; the omission summary always prints five counts;
  the queue parks every blocked review as the bare `review-blocked`.
- Characterization at `4e61bf6d`: the 57 review prompt-bound, segment and
  scope tests pass, and four of them pin today's refusal text exactly;
  `roundfix profiles show` in a repository whose Project Config sets
  `pre_pr_review.omit_paths` exits 2 with `field omit_paths not found in type
  config.prePRReviewOverlay`, so no adopter can carry the key yet.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd query "roundfix pre-PR review claude prompt too large derived_paths omit review segments" --all --files --min-score 0.3`;
  it returned this repository's own mirrored ADRs and review guide and an
  adopter's mirrored review authorization page, which was not read further
  (the brief forbids reading Vortex content beyond its message); none decides
  review omissions or Claude segmentation.
- Exa: the GitHub and Claude Agent SDK passages quoted above.

## Glossary

- adds: **Review Omit Path** — a Project Config path the pre-PR review leaves out
- changes: **Review Segment** — Claude candidates are segmented too, bounded by tokens
- not a term: **Delivery Queue** — existing glossary term, unchanged
- not a term: **Park Class** — existing glossary term, unchanged
- not a term: **Pre-PR Review Record** — existing glossary term, unchanged

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AW brief for Spec 0342 of 2026-11-03 and the maintainer's answers below.

The maintainer's answers (AskUserQuestion, 2026-11-03), verbatim:

- 0342: "Segmentar + omit_paths (Recommended)" — option text: "Claude também
  divide em até 4 Review Segments; nova chave pre_pr_review.omit_paths
  (globs) deixa arquivos gerados fora da revisão, sem usar derived_paths;
  quando ainda assim bloquear, a mensagem nomeia os 3 maiores arquivos e a
  chave que os tiraria."
- Rules (all three): "Mesmas do AV (Recommended)" — option text: "Push e PR
  dos planos, entrega e merge pela Delivery Queue, consertar e seguir em
  falha técnica (pular só se o conserto contradisser o PRD). Makefile,
  .github/workflows, go.mod/go.sum e config de lint/formatter liberados
  quando o desenho pedir."
- Release (all three): "Uma v0.130.0 no fim (Recommended)" — option text:
  "Minor único depois das três Specs, porque a 0342 cria a chave
  pre_pr_review.omit_paths e a 0343 cria o --all. Notas aos 8 adotantes,
  resposta ao Fluxus com os campos waiting_on/since, aviso ao Vortex da chave
  nova e aviso ao Jarvis de que a versão saiu."

Taken from those answers, ADR-0227, ADR-0282 and the reading above, with no
further interview:

- The key takes the derived-path entry grammar, including a directory ending
  in `/`: the report asks for globs, and a review omission has none of the
  authored-file risk that keeps directory entries out of derived review
  omissions.
- The key is read from the base commit's Project Config only, so a candidate
  cannot choose what its own review skips; an adopter lands it on the default
  branch first.
- A Claude segment is bounded by the estimated tokens of its whole prompt as
  sent, the measure the budget already uses, not by bytes; the planner keeps
  the Codex rule of consecutive first-parent ranges, so a finding still
  anchors in a tree its author made.
- Files are ranked by the bytes of their section of the measured diff: both
  bounds derive from bytes, so the ranking is the same in tokens, and bytes
  are what an operator can check with Git.
- The refusal's suffix is appended to every size refusal, Codex included;
  the existing text before it keeps its bytes. This is the one change of
  Codex behavior, declared in `_techspec.md` → Declared breaks.
- The park carries the first line of the record's reason for every blocked
  review, because the queue only learns the outcome, and the Park Class reads
  only the text before the colon.
- Four implementation Tasks in series (omit paths; Claude segments; the
  actionable refusal and park; documentation) and a QA gate.

## Open Questions

None.
