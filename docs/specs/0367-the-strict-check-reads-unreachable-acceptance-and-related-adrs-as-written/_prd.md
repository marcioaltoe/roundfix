---
spec: 0367-the-strict-check-reads-unreachable-acceptance-and-related-adrs-as-written
status: active
created: 2026-11-09
surfaces: [cli, docs]
type: fix
---

# The strict check reads Unreachable Acceptance and related ADRs as written

Two defects of `roundfix spec check --strict` cost Specs their delivery for reasons the check could have named while the Spec was still being written. A PRD whose `## Unreachable Acceptance` section holds an entry the declaration reader does not read, such as Spec 0364's numbered prose paragraph, passes the strict check and fails only at the end of the Run, when `roundfix qa-report accept` counts zero declarations. And `SC-ADR-RELATED` treats an ADR a Spec lists only to record that it does not apply exactly like one the Spec applies, so each such listing draws in the ADRs that cite it; Fiscus's Spec 0081 had to list 72 ADRs over seven rounds, all of them not applying. After this Spec the strict check reports, at its PRD line and with the fielded form in its fix, any text in the Unreachable Acceptance section that no declaration reads, and only an ADR the Spec applies or changes draws in the ADRs that cite it: an ADR whose TechSpec ADR Alternatives line begins with `not applicable` stays listed but draws in nothing.

## Sources

Two sources, adopted (see `references/_index.md`):

1. [`references/2026-11-07-the-strict-check-names-an-unreachable-acceptance-it-cannot-read.md`](references/2026-11-07-the-strict-check-names-an-unreachable-acceptance-it-cannot-read.md),
   the open Backlog Entry requested by the maintainer on 2026-11-07, relayed
   through jarvis, `type: fix`; this Spec adopts it.
2. [`references/2026-11-07-fiscus-sc-adr-related-grows-through-not-applicable-adrs.md`](references/2026-11-07-fiscus-sc-adr-related-grows-through-not-applicable-adrs.md),
   ai-memory message 01a1b6b1, reported by fiscus on 2026-11-07, a finding.

Both change the same detector family in `internal/speccheck` and the same
command; the maintainer chose one Spec for them. At the authoring base
`docs/backlog/` holds no other open entry and `docs/findings/` no unresolved
Finding that shares this context.

## Prerequisites

None in `requires`. Authored against `2e59aebe` (v0.137.1). The queue
delivers Specs 0365 and 0366 first. They may also change `CONTEXT.md` (other
entries), `docs/references/behavior-surfaces.json`,
`skills/testdata/owned-skill-versions.json` and, for 0366, the Roundfix skill
and its references other than `references/spec.md`. This Spec edits only the
**Unreachable Acceptance** and **Spec Consistency Check** entries of
`CONTEXT.md`, the `spec` reference of the Roundfix skill, and the
`write-techspec` and `write-tasks` skills; it re-records the version record and
the Behavior Surface Record by their commands, so a rebase reruns the commands
instead of merging bytes. It does not touch `internal/agent/**`,
`skills/qa-gate/**`, `.agents/skills/qa-gate/**` or `internal/spec/qa*.go`.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the new finding code joins the existing `SC-*` vocabulary of the
  `roundfix-speccheck/v1` report. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the check reads local files and
  Git history only and makes no network call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0377 (new) records the rule,
  ADR-0377: "Only an ADR the Spec applies or changes seeds the related-ADR
  candidates." It amends ADR-0093, whose candidate rule this Spec narrows,
  ADR-0093: "an ADR that cites an ADR the Spec already lists is a candidate
  the Spec should account for", and ADR-0168, whose first paragraph restates
  it, ADR-0168: "reports an accepted ADR that cites an ADR a Spec lists";
  task_02 writes both amendment notices. ADR-0265 defines the ADR Alternatives
  line the new marker lives in, and the check stays structural, ADR-0265:
  "The check reads structure, never meaning." ADR-0094 keeps a detector
  quiet when its artifact is missing, ADR-0094: "is skipped, not failed, when
  one of them is absent", and the new detector skips a missing PRD.
  ADR-0343 matches each declared-blocked QA row with the declaration that
  names its criterion, ADR-0343: "a declared row matches the criterion it
  names", and the new detector guards the declarations that matching reads.
  ADR-0208 lets one Spec adopt sources that share a context, ADR-0208: "One
  Spec may adopt several Inbox Entries, Backlog Entries and Findings whose
  context is similar or complementary", which is why both sources are here.
  Not applicable, and marked `not applicable` in the TechSpec: ADR-0210 and
  ADR-0211, which the adopted Fiscus reference cites by Fiscus's own
  numbering, so this repository's records of those numbers do not bind this
  Spec; and, each within the related-ADR candidates of a listed ADR on a part
  this Spec leaves alone, ADR-0156, ADR-0176, ADR-0183, ADR-0184, ADR-0209,
  ADR-0220, ADR-0274, ADR-0276, ADR-0280, ADR-0295, ADR-0319, ADR-0335 and
  ADR-0353. Source: `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle BA answers of 2026-11-09 quoted in the authorization record grant the
  owned skills this Spec edits and the two Spec Check files an earlier
  authorization bounded; no Makefile, workflow, Go module, lint or formatter
  change. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0367-the-strict-check-reads-unreachable-acceptance-and-related-adrs-as-written/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/write-tasks/SKILL.md`,
  `.agents/skills/write-techspec/SKILL.md`,
  `.agents/skills/write-techspec/references/techspec-template.md`,
  `internal/speccheck/coherence.go`, `internal/speccheck/constraints.go`,
  `skills/roundfix/SKILL.md`, `skills/write-tasks/SKILL.md`,
  `skills/write-techspec/SKILL.md`,
  `skills/write-techspec/references/techspec-template.md`.

## Goals

- The strict check names every piece of text in a PRD's Unreachable
  Acceptance section that the QA gate will not read, at its line, before a Run
  starts, and shows the fielded form to use.
- A section that the reader reads fully, an absent or empty section, and a
  `None.` paragraph keep passing, and the declarations the QA gate and
  `roundfix archive` read do not change.
- An ADR a Spec lists only as not applying stops drawing in the ADRs that
  cite it, while an ADR the Spec applies keeps drawing them in as today.
- ADR-0093, ADR-0168, the glossary, the user guides and the authoring skills
  say how to mark an ADR as not applying and what the new finding means.

## User Stories

1. As a Spec author, I want the strict check to name an Unreachable
   Acceptance entry the QA gate cannot read, so that I rewrite it before the
   Run instead of losing the QA gate at the end.
2. As a Spec author whose PRD lists ADRs that do not apply, I want them to
   stop drawing in the ADRs that cite them, so that my Active ADR obligations
   row lists what the Spec touches and not a cascade of unrelated decisions.
3. As a maintainer reading the ADRs and the guides, I want the related-ADR
   rule and the new finding written where the check is described, so that
   the rule I read is the rule the check enforces.

## Core Features

1. `SC-UNREACHABLE-UNREAD`, an error from the PRD stage: one finding per run
   of consecutive lines in `## Unreachable Acceptance` that no declaration
   reads, at its first line, and one finding for a declaration without its
   `criterion:`, `reason:` or `satisfied-by:` field. The reader in
   `internal/spec` that the QA gate and `roundfix archive` use names those
   lines, so the check and the gate cannot disagree; the declarations it
   returns are unchanged.
2. `SC-ADR-RELATED` seeds its candidates only with the listed ADRs whose
   TechSpec ADR Alternatives line does not begin with `not applicable`; its
   fix says so. Without a TechSpec every listed ADR seeds them, as today.
3. Documentation: ADR-0377 (written by this plan) and the amendment notices
   in ADR-0093 and ADR-0168; the **Unreachable Acceptance** and **Spec
   Consistency Check** entries of `CONTEXT.md`; the `spec` command guide and
   the context-driven development guide; the Roundfix skill's `spec`
   reference and the `write-techspec` and `write-tasks` skills.

## User Experience

On Spec 0364's planned PRD, whose one entry was the paragraph
`1. A real rerun of a failed release job on GitHub Actions against the npm`
and three indented lines, `roundfix spec check <slug> --strict` prints:

```text
[error] SC-UNREACHABLE-UNREAD: docs/specs/<slug>/_prd.md line 165 under ## Unreachable Acceptance is not part of a fielded declaration, so the QA gate does not read it: "1. A real rerun of a failed release job on GitHub Actions against the npm"
  at docs/specs/<slug>/_prd.md:165
  fix: Write each entry in the fielded form: a top-level "- criterion: <text the QA row repeats>" line, then indented "reason: <why no Verification or QA row can reach it>" and "satisfied-by: <evidence or follow-up that satisfies it>" lines; start the section with None. when no criterion is unreachable.
```

A declaration that lacks a field prints
`<prd> line <N> starts an Unreachable Acceptance declaration without its <field>: field`
with the same fix. The `SC-ADR-RELATED` fix gains the sentence
`when it does not apply, begin its ADR Alternatives line in the TechSpec with "not applicable", so the ADRs that cite only it are not drawn in.`

## Non-Goals / Out of Scope

- Making `roundfix qa-report accept`, `roundfix settle` or `roundfix archive`
  refuse text the reader does not read; they keep reading the same
  declarations, and the check names the text before a Run.
- Changing the ADR horizon of ADR-0168, the one-edge depth of the candidates,
  `SC-ADR-UNLISTED`, `SC-ADR-ALTERNATIVES-UNRELATED` or the gap severity of
  `SC-ADR-RELATED`.
- A marker in another language, or a marker inside the PRD's Active ADR
  obligations row; ADR-0377 records why.
- Rewriting any archived Spec, or the `write-prd` template.
- Any change to the Makefile, the workflows, `go.mod`, `go.sum`, lint or
  formatter configuration, the Baseline, or the `### QA settlement` section
  of any skill.

## Success Metrics

1. Success Metric: on Spec 0364's planned PRD, the PRD of commit `4dce106d`,
   the strict check reports exactly one `SC-UNREACHABLE-UNREAD`, at line 165;
   with that entry rewritten in the fielded form, it reports none.
2. Success Metric: a fielded section, an absent section, an empty section, a
   `None.` paragraph and the `write-prd` template comment left in place report
   no `SC-UNREACHABLE-UNREAD`, and every existing test of `Unreachable` in
   `internal/spec` passes unchanged.
3. Success Metric: in a corpus of three accepted fixture decisions, where
   ADR-<B> cites ADR-<A> and ADR-<C> cites ADR-<B>, a Spec that applies ADR-<A> and marks ADR-<B>
   `not applicable` gets no `SC-ADR-RELATED` for ADR-<C>; with ADR-<B>
   applied, or without a TechSpec, it gets one, as at the authoring base.
4. Success Metric: this Spec's own PRD and TechSpec report no finding under
   the delivered strict check.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The Fiscus report quoted in the adopted reference, written by the Fiscus
  project before this Spec existed: "O fecho levou 7 rodadas e 72 ADRs".
- Spec 0364's planned PRD in this repository's history, commit `4dce106d`,
  written for another Spec, whose numbered entry cost its first QA gate.
- The CommonMark specification, section 5.2 List items
  (https://spec.commonmark.org/0.31.2/, read through Exa on 2026-11-09): "A
  bullet list marker is a `-`, `+`, or `*` character." and "An ordered list
  marker is a sequence of 1–9 arabic digits (`0-9`), followed by either a `.`
  character or a `)` character." A `1.` entry is an ordered list item, which
  the declaration reader, built on `- ` bullets, never reads.

## Unreachable Acceptance

- criterion: Fiscus Spec 0081: after the release that carries this Spec,
  Fiscus's Spec 0081 passes the strict check with only the ADRs it applies
  and the ADRs one citation edge from them listed, once its TechSpec marks
  the others `not applicable`.
  reason: it needs another repository, which no Verification or QA row may
  read.
  satisfied-by: the fiscus operator runs `roundfix spec check
  0081-a-exclusao-retirada-sempre-reenvia --strict` after the patch release
  and reports the result through ai-memory, recorded session evidence.

## Validation

The decisions were checked on 2026-11-09 against these sources.

- Characterization with a binary built from `2e59aebe`, in a temporary
  repository whose ADR-<B> cites ADR-<A> and ADR-<C> cites ADR-<B>, and
  whose PRD lists ADR-<A> as applying and ADR-<B> as not applying with a
  `- ADR-<B>: not applicable; …` TechSpec line: `spec check --strict`
  reports `SC-ADR-RELATED` for ADR-<C>, which "cites listed ADR-<B> but is
  not itself listed by the Spec". The Fiscus cascade reproduces. The same
  PRD's `1. A real release run …` Unreachable Acceptance entry reports
  nothing.
- Replay of Spec 0364's planned PRD (`git show
  4dce106d:docs/specs/0364-a-rerun-finishes-a-partial-npm-publication/_prd.md`)
  in a temporary repository: the binary from `2e59aebe` reports no finding
  about the section; a prototype of this design reports one
  `SC-UNREACHABLE-UNREAD` at line 165.
- Survey of the 81 archived PRDs that hold `## Unreachable Acceptance`, each
  read at its record's `source_revision`: the prototype reports no
  `SC-UNREACHABLE-UNREAD` on any of them; 27 of them open the section with a
  `None.` paragraph.
- This Spec's own ADR list: at `2e59aebe` the ADRs it cites, with the three
  numbers the Fiscus reference cites, grow over three rounds of related-ADR
  candidates to 22 listed ADRs, the row above; under the new rule the
  fifteen it marks `not applicable` draw in nothing, and a prototype of this
  design reports no finding on this Spec.
- Reading at `2e59aebe`: `obligationCitationNumbers` reads every zero-padded
  number in the Active ADR obligations row as listed, and
  `detectADRConsistency` seeds `firstListedCitation` with that set;
  `parseUnreachableDeclarations` starts a declaration only at a top-level `- `
  line and ignores other text; `internal/speccheck` never calls
  `spec.Unreachable`.
- Published sources (Exa): the CommonMark specification, section 5.2.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd query "spec check SC-ADR-RELATED related ADR closure not applicable cascade" --all --files --min-score 0.3`
  and
  `qmd query "unreachable acceptance fielded declaration parser QA rows_blocked_declared" --all --files --min-score 0.3`;
  they returned this repository's mirrored ADRs, guides and the adopted
  Backlog Entry, and one archived Fiscus Spec that was not opened, because
  this Spec reads no other repository; no page decides otherwise.
- Pending ai-memory messages: none read by this author; the operator relayed
  the Fiscus message.

## Glossary

- changes: **Unreachable Acceptance** — the Spec Consistency Check reports text in the section that no fielded declaration reads.
- changes: **Spec Consistency Check** — only an ADR the Spec applies or changes draws in the ADRs that cite it.

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle BA brief for Spec 0367 of 2026-11-09 and the maintainer's answers below.

The maintainer's answers (AskUserQuestion, 2026-11-09 ~08:35–08:45 BRT),
verbatim:

- Question "Finding do Fiscus: a regra SC-ADR-RELATED (ADRs 0093/0168
  aceitas) faz uma ADR marcada “não se aplica” puxar todas as ADRs que a
  citam (72 ADRs numa Spec). Emendo a ADR?" → "Emendar (Recommended)" —
  option text: "Só ADRs aplicadas ou alteradas pela Spec entram no fecho; uma
  ADR “não se aplica” deixa de propagar. Task com “Amended by” na ADR
  aceita."
- Question "Os dois itens da checagem estrita (backlog do Unreachable
  Acceptance ilegível + cascata do SC-ADR-RELATED do Fiscus) vão numa Spec só
  ou em duas?" → "Uma Spec (Recommended)" — option text: "Mesmo pacote
  (internal/speccheck), ~2 Tasks + QA, um release patch. Jev: coerente 0,85;
  dividir 0,44."
- Question "Que autorizações você concede de saída para eu levar o ciclo BA
  até o fim sem parar (implementações, merges e releases inclusos como no
  AZ)?" → selected "Mesmas do AZ (Recommended)" — "Release patch assim que
  cada Spec ficar pronta, teto de 2 corretivas por Spec, merge pela fila,
  avisos ao Jarvis, sem notas de release às caixas."; "Scripts/Makefile/CI" —
  "Mexer em ~/.roundfix-operator, Makefile e .github/workflows se a migração
  ou a Spec B precisarem."; "go.mod/lint" — "Mexer em go.mod/go.sum e na
  configuração de lint/formatter se uma Task precisar."; "Avisar
  Quaestor/Fiscus" — "Mensagem ai-memory curta a cada um dizendo qual Spec
  adotou o finding (não é nota de release)."

Taken from those and the evidence, with no further interview:

- The marker is the leading token `not applicable` of the ADR's line under
  the TechSpec's `## ADR Alternatives`, the per-ADR line every listed ADR
  already needs; ADR-0377 records it and the options it rejected. ADR-0377
  amends ADR-0093 and ADR-0168, and task_02 writes both notices.
- The parser does not refuse unread text: `qa-report accept`, `settle` and
  `archive` read the same declarations as before, so a Spec already in a Run
  keeps its gate, and the strict check names the text before the next Run.
  One reader in `internal/spec` serves both, so they cannot drift.
- `SC-UNREACHABLE-UNREAD` is an error, because the check locates the text
  and the reader that drops it, and a `None.` paragraph opening a section
  with no declaration stays accepted, as 27 archived PRDs wrote it.
- Three implementation Tasks and the QA gate, in series: the unread
  detector, the related-ADR seeds, then the guides and skills.
- The release is the operator's patch right after the merge; no Task
  releases.

## Open Questions

None.
