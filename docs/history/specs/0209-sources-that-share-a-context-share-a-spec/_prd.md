---
spec: 0209-sources-that-share-a-context-share-a-spec
status: archived
created: 2026-10-01
surfaces: [backend, cli, docs]
archived: "2026-10-01"
source_slug: 0209-sources-that-share-a-context-share-a-spec
---


# Sources that share a context share a Spec

Nothing in the Baseline tells an author to look for related work before
minting a Spec or a record. One Spec per Finding looks like the safe default,
and so does a new Backlog Entry for an intent an open entry already holds. The
result is more Specs than the work needs, each paying its own authoring,
review and QA, and records that say the same thing twice. In the fluxus
repository a Spec was authored without consulting the backlog and left its
own Backlog Entry open after it shipped, which a later Spec had to clean up.

On 2026-10-01 the maintainer asked for the opposite, and for it to be a
Baseline rule:

> "podemos ter em uma mesma spec um ou mais inbox/backlogs ou findings não é
> obrigatório e nem desejavel uma spec para um finding com contexto similar ou
> complementar a outro. O cuidado é não ter specs desnecessáriamente complexas
> e grandes, o uso do jev pode ajudar a decidir pela reunião de inbox, backlogs
> e findings em uma só spec. Isso deve fazer parte do baseline de roundfix e
> parte das regras do roundfix."

> "podemos alterar e adicionar e atualizar um finding ou backlog já existente
> e ainda não implementado para que não seja necessário criar outro arquivo"

This Spec adds three Baseline clauses with force: sources that share a
context may share a Spec, grouping stops at four implementation Tasks plus
the QA gate, and an open record is extended instead of duplicated. The
authoring skills teach both rules. The advisory Spec judge gains one measured
question, whether two sources belong in the same Spec, and raises a pair as a
suggestion the author answers.

## Prerequisites

This Spec is delivered after Spec 0205, which adds the advisory
`roundfix spec judge` command, its question file, its readers, its Judge Log
and the `## Advisory judgment` heading of the `write-prd` and
`write-techspec` skills. Tasks 02 and 03 extend that code, and their
Verification runs tests Spec 0205 creates. The Task Graph names Spec 0205 in
`requires`, so the queue owner waits for it; v0.25 ships both, Spec 0205
first.

Spec 0206 also edits the Spec workflow module and its Source Baseline rows.
Whichever lands second raises each version by one from the value on its
starting commit; the clauses of the two Specs do not overlap.

## Project Constraints

- Identifier strategy: not applicable — no new identifier. The three clauses
  are named by fixed Baseline clause identifiers, a grouping judgment by its
  fixed kind `source-grouping`, and a source by its repository path. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: applicable — the grouping question goes through
  the transport Spec 0205 built, unchanged: OpenRouter's System One API with
  `ROUNDFIX_OPENROUTER_API_KEY`, or else TypeSafe directly with
  `ROUNDFIX_TYPESAFE_API_KEY`, each key read only from the command's
  environment and sent only to its own endpoint. What a request may carry
  grows by one reader: the text of Findings and Backlog Entries of this
  repository, within the maintainer's authorization of 2026-09-30. No new
  endpoint, key, header or field. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0208 (this Spec) decides the three
  clauses and the grouping bound. ADR-0209 (this Spec) decides that the judge
  suggests a grouping and never decides it, at the measured threshold, and
  widens the request boundary to Findings and Backlog Entries. ADR-0200 keeps
  the judge advisory and pinned to Jev 1.13, and the grouping question
  inherits both, ADR-0200: "A raised judgment is a finding the authoring
  model answers". ADR-0201 fixes the recipients, the keys, the Judge Log and
  the monthly ceiling, which stay as they are; its statement that Findings
  and Backlog Entries are not read is the one ADR-0209 changes. ADR-0083
  moves an adopted source into exactly one owning Spec, which is how grouping
  records ownership. The extend clause keeps the boundary ADR-0092 draws,
  ADR-0092: "a finding is never a commitment and a backlog entry is never
  evidence", by giving a Finding only an evidence addendum and a Backlog Entry
  only a revision of its intent. The clauses cite no record, ADR-0186: "It
  never cites a Spec number or an ADR number". Each new clause gains its
  Source Baseline row so the update stays accounted, ADR-0060: "A Source
  Baseline must be project-agnostic and independently account for every
  Normative Clause", and ADR-0058: "preview and apply block while any clause
  is unaccounted". That accounting records each new clause `retained`
  without a model, ADR-0099: "it is a mechanical comparison". ADR-0203 governs promoting an
  adopter's hand-written rule into the Baseline; these clauses are the
  maintainer's own workflow rules, not a promotion, so its recurrence test
  does not apply. ADR-0189 ties an owned skill's version to its content, so
  each edited skill raises its version. This Spec names Spec 0205 as its
  prerequisite, ADR-0193: "A Spec now names the Specs it needs". ADR-0081 and
  ADR-0149 sanction the regeneration that follows the Baseline and skill
  edits. ADR-0089 passes the environment to code under test explicitly. The
  TechSpec states the command's new output as transcripts, ADR-0184: "A
  TechSpec now declares numbered Surface Transcripts". ADR-0187 splits the
  Roundfix Skill and the command reference by command. ADR-0182 runs
  Settlement Checks before each Task commit, ADR-0178 authorizes a Task commit
  by its grant, and ADR-0166 records undeclared paths; every Task here
  declares its paths. This gate aims at `pass`, and its Pull Request row is
  the pre-PR row, ADR-0167: "That row no longer counts against a qualifying".
  This Spec's gate is bound by ADR-0080, ADR-0088, ADR-0091, ADR-0096,
  ADR-0104, ADR-0117, ADR-0155 and ADR-0156, and ADR-0093 and ADR-0094 check
  its consistency by citation and artifact presence. ADR-0097 cites ADR-0080
  but carries a QA row forward, ADR-0116, ADR-0176 and ADR-0183 cite ADR-0083
  or ADR-0093 but decide how the Spec Consistency Check reads citations,
  ADR-0168 cites ADR-0093 but narrows the related-ADR check, ADR-0194 and
  ADR-0195 cite ADR-0097 but decide how the Daemon records and observes QA
  rows again, ADR-0210 cites ADR-0097 but decides how an evidence snapshot
  digests its inputs, and ADR-0192 cites ADR-0149 but decides how the queue owner
  resolves a merge conflict on declared derived paths. This Spec changes none
  of these nine, so none of them applies. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30 for the Baseline source and its guides ("Autorizar os dois") and
  for the skills ("considere autorizado a ajustar todas as skills se
  necessário"), and the maintainer's requests of 2026-10-01 quoted above,
  recorded in [_authorization.md](_authorization.md); bounded files:
  `internal/baseline/assets/modules/spec-workflow.json`,
  `internal/baseline/assets/modules/context-workflow.json`,
  `internal/baseline/assets/source-baselines/index.json`,
  `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/baseline.json`,
  `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/manifest.json`,
  `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/spec-routing.md`,
  `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/docs-layout.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `docs/agents/spec-routing.md`, `docs/agents/docs-layout.md`,
  `docs/agents/setup-context.json`,
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/write-idea/SKILL.md`, `skills/write-idea/SKILL.md`,
  `.agents/skills/write-prd/SKILL.md`, `skills/write-prd/SKILL.md`,
  `.agents/skills/write-techspec/SKILL.md`, `skills/write-techspec/SKILL.md`.
  Sanctioned regeneration: `make baseline-digests`, `make skills-sync` and
  the Managed Refresh. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`.

## Goals

- An author looks for open Backlog Entries and unresolved Findings that share
  a source's context before minting a Spec for it, and adopts them into one
  Spec when they share it.
- Grouping never produces a Spec past four implementation Tasks plus its QA
  gate: a larger grouped scope is split, and each source keeps exactly one
  owning Spec.
- An observation or intent that fits an open Backlog Entry or an unresolved
  Finding extends it instead of becoming a new file, while a Finding stays
  immutable history.
- These rules reach every adopter as Baseline clauses with force, and an
  adopter's Baseline update accounts for each of them.
- The Spec judge suggests which open sources may belong in the Spec, at the
  measured threshold, and the suggestion never gates and never decides.

## User Stories

1. As an authoring model about to mint a Spec for one Finding, I want the
   rules to tell me to look for open sources that share its context, so that
   one Spec resolves them together.
2. As an authoring model grouping sources, I want a stated bound, so that the
   grouped Spec stays small enough to deliver.
3. As a session about to record a new observation, I want to be told to
   extend the open record it fits, so that the repository holds one record
   per concern.
4. As an authoring model running the Spec judge, I want each open source
   that likely belongs in this Spec listed with its probability, so that I
   adopt it or state why it stays apart.
5. As a maintainer, I want the grouping question to send only Findings and
   Backlog Entries of this repository, so that the data authorization holds
   without anyone remembering it.
6. As an adopter, I want my Baseline update to account for the new clauses,
   so that it stays ready instead of refusing.

## Core Features

1. **Sources that share a context share a Spec.** A mandatory Spec-workflow
   clause lets one Spec adopt several Inbox Entries, Backlog Entries and
   Findings whose context is similar or complementary, states that one Spec
   per source is neither required nor preferred, asks the author to look for
   such sources before minting a Spec, and calls a grouping suggestion from
   the judge advisory.
2. **The grouping bound.** A mandatory Spec-workflow clause groups sources
   only while the Spec fits four implementation Tasks plus its QA gate, and
   splits a larger grouped scope into Specs that each fit, each source with
   one owner.
3. **Extend before minting.** A mandatory CONTEXT-workflow clause has every
   new Finding or Backlog Entry, Triage included, first look for an existing
   one not yet implemented that fits: an `open` Backlog Entry is revised in
   place, an unresolved Finding gains a dated addendum.
4. **Accounted for adopters.** Each new clause gains its Source Baseline row,
   so an adopter's Baseline update records it `retained`, and no existing
   clause changes.
5. **The skills teach it.** `write-idea`, `write-prd` and `write-techspec`
   each gain a section on grouping sources and on extending open records,
   and `write-prd` and `write-techspec` tell the authoring model how to
   answer a grouping suggestion.
6. **The grouping question.** `roundfix spec judge` asks, for every pair of
   one Finding or Backlog Entry the Spec adopted and one open Backlog Entry
   or unresolved Finding, the measured question whether both should be
   resolved by the same Spec, at every stage.
7. **A suggestion at 0.3.** A pair whose answer is 0.3 or more is printed as
   `suggested` with its probability; any other pinned answer is `clear`. The
   summary line counts suggestions. The command still exits `0`.
8. **The grouping boundary.** The grouping question reads only a Finding or
   Backlog Entry the Spec adopted, an `open` Backlog Entry or an unresolved
   Finding of this repository, as an English regular file directly in its
   directory, prepared exactly as measured.

## User Experience

The author reads one line per suggested pair: the adopted source, the open
source and the probability. Clear pairs are counted, not listed. The skills
say what to do with a line: adopt the open source within the bound, or say in
the report why it stays apart.

## Declared breaks

- `roundfix spec judge`'s summary line gains `N suggested` after
  `N advisory`, so Spec 0205's text Transcripts 1, 4 and 5 gain
  `0 suggested`. No release has shipped the old line.
- A Judge Log line of a grouping judgment has a null `line` and the outcome
  `suggested`, which Spec 0205's log never wrote.
- The judge's requests may carry Finding and Backlog Entry text, which
  ADR-0201 had kept out.
- `docs/agents/spec-routing.md` and `docs/agents/docs-layout.md` gain three
  mandatory clauses in every adopter.

## Non-Goals / Out of Scope

- Gating on a grouping: no exit code, Task, check or queue reads a
  suggestion.
- Moving, adopting or editing a source automatically.
- Sending an Inbox Entry, a terminal Finding or Backlog Entry, or a source
  another Spec adopted.
- Suggesting that a Spec be split, or asking whether two adopted sources
  belong apart.
- A general size rule for every Spec: the bound limits grouping only.
- Changing an existing Baseline clause, including the Triage clause.
- Re-measuring the question, moving its threshold or the pinned model.
- Any live call to OpenRouter or TypeSafe from a test, a Verification or the
  QA gate.

## Success Metrics

1. Success Metric: the rendered `docs/agents/spec-routing.md` and
   `docs/agents/docs-layout.md` of this repository and of the Standard
   TypeScript Monorepo profile state the three clauses with `mandatory`
   force, and a second Managed Refresh changes no file.
2. Success Metric: a Standard TypeScript Monorepo adopter's Baseline update is
   ready and records each new clause `retained`, with no clause
   `unaccounted`.
3. Success Metric: on a fixture Spec that adopted one Finding, with one open
   Backlog Entry answered 0.81 and one answered 0.04, `roundfix spec judge`
   prints the first as `suggested` with `P(same Spec) 0.81`, counts the second
   `clear`, and exits `0`.
4. Success Metric: every string a grouping request carries comes from an
   accepted source, and an Inbox Entry, a `declined` Backlog Entry, a `done`
   Finding, a symbolic link, a non-English source and a Go file send nothing.
5. Success Metric: the grouping question, threshold and preparation in the
   question file equal the measurement's.
6. Success Metric: `write-idea`, `write-prd` and `write-techspec` each carry
   the grouping section, and each mirror equals its canonical copy.

## Recorded limits

- Recall is low: at 0.3 the question found 27% of the pairs that belonged
  together. No suggestion does not mean no grouping fits.
- The measured AUROC, 0.783, is not significantly above the TF-IDF
  baseline's 0.753. Jev is used for its precision at the operating point,
  0.93 against 0.85, and because the maintainer asked for it.
- The threshold was chosen on the data it is reported on.
- The positives are pairs one archived Spec adopted together, so the measure
  rewards the grouping this repository already did, not an ideal one.
- A run asks one question per pair, so a long open backlog makes the run
  slower and dearer; the monthly ceiling still bounds the spend.
- The same suggestion returns at the TechSpec stage after the author answered
  it at the PRD stage.

## Decisions

- **Three clauses, no change to an existing one.** See ADR-0208.
- **A suggestion, at the measured threshold, from source text only.** See
  ADR-0209.
- **The bound limits grouping, not every Spec.** A general size rule would
  break Specs already authored past it and was not asked for.
- **Grouping at every stage.** A refactor or bug fix enters at the TechSpec
  stage and is the Spec most likely to adopt Findings, so the question cannot
  wait for a PRD stage it never runs.
- **Inbox Entries are grouped by hand only.** The clause lets a Spec adopt
  them, and the judge does not send them, because the data authorization
  names Findings and Backlog Entries.

## Acceptance evidence

Each Core Feature requires positive and negative evidence in the Task Graph.
The outside-evidence row rests on sources this Spec did not produce:

- **A measurement this Spec did not design.** On 2026-10-01 the maintainer's
  session measured the grouping question against this repository's archive:
  107 adopted sources from 43 archived Specs, 94 positive pairs (two sources
  one Spec adopted) and 94 hard negatives (sources different Specs adopted
  within three days), Spec and ADR numbers scrubbed. Through OpenRouter,
  `jev-1.13`, reported as `typesafe/jev-1.13-20260917`, cost US$0.0073. Jev
  AUROC 0.783 against TF-IDF 0.753, difference 95% CI [-0.032, 0.093]; at a
  probability of 0.3 or more, precision 0.93 and recall 0.27, against a
  precision of 0.85 for TF-IDF's top 27 pairs. The script is session
  evidence of 2026-10-01, outside this repository; the TechSpec reproduces
  its question and preparation.
- **Published literature on related-record detection.** Sun, Lo, Khoo and
  Jiang, "Towards More Accurate Retrieval of Duplicate Bug Reports" (ASE
  2011, <https://cs.uwaterloo.ca/~cnsun/public/publication/ase11/ase11.pdf>),
  present duplicate detection as a list of candidates for a triager and note
  that several reports can complement one another. Zhang et al., "Duplicate
  Bug Report Detection: How Far Are We?" (TOSEM 2023,
  <https://dl.acm.org/doi/full/10.1145/3576042>), find a simple retrieval
  technique ahead of recent deep models on most projects, and report that the
  VSCode bot in production suggests at most five candidates. Both read
  2026-10-01. They support a suggestion over a decision and a lexical
  baseline as the honest comparison.
- **The question shape, published.** TypeSafe's Noul page
  (<https://docs.typesafe.ai/primitives/noul.md>, read 2026-10-01) states
  that a Noul question has `type` `noul` and `instructions`, that `criteria`
  is optional, and that the answer is a single probability of yes under
  `noul`. The measured question has no `criteria`, and the question file
  keeps it that way.
- **Another repository's record.** The fluxus repository's archived Spec
  0035 records, in Portuguese, that its Spec 0034 was the only Spec of its
  series authored without consulting the backlog, and that the Backlog Entry
  it delivered was still open afterwards. Read through the Secondbrain mirror
  on 2026-10-01.

## Research basis

The Secondbrain was consulted through `wiki/index.md` and the queries
`qmd query "agrupar findings e backlog na mesma spec evitar duplicação" --all
--files --min-score 0.3` and `qmd query "Jev similaridade pares mesma
especificação agrupamento" --all --files --min-score 0.3`. The first returned
the fluxus Spec 0035 PRD cited above, which is why the grouping clause asks the
author to look at the open backlog before minting. The second returned
`wiki/sources/jev-in-the-agent-loop-n01ennn-2026-09-28.md`, which keeps the
thresholds in code and reads a Noul near 0.5 as a state that does not separate
the cases; that is why the threshold lives in the question file and the
suggestion is advisory. The Secondbrain guide's existing rule to extend a
strong verified match instead of duplicating an Inbox Entry is the model for
the extend-before-minting clause. Exa found the two papers cited above.

## Open Questions

- Whether the judge should also send Inbox Entries once the data
  authorization names them. Default until answered: it does not.

## Technical candidate

The [_techspec.md](_techspec.md) records the clause texts, the grouping
question, the readers, the transcripts, coverage and build order.
