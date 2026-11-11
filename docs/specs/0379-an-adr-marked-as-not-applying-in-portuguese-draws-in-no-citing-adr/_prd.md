---
spec: 0379-an-adr-marked-as-not-applying-in-portuguese-draws-in-no-citing-adr
status: active
created: 2026-11-10
surfaces: [cli, docs]
type: feat
---

# An ADR marked as not applying in Portuguese draws in no citing ADR

A Spec written in Portuguese marks an ADR it lists only to record why it does not apply with "não se aplica", but the Spec Consistency Check reads only the English leading token `not applicable` in the TechSpec's `## ADR Alternatives` line. Every such line counts as an applied ADR, so it keeps drawing in the ADRs that cite it, the cascade ADR-0377 removed, and nothing tells the author why the closure stays large. After this Spec the check reads a small closed list of leading tokens, the English `not applicable` and the Portuguese `não se aplica` with its unaccented spelling, without regard to case, and the authoring skills, the guides and the glossary say so.

## Sources

One source, adopted (see `references/_index.md`):

1. [`references/2026-11-10-an-adr-marked-as-not-applying-in-portuguese-still-seeds-the-related-adrs.md`](references/2026-11-10-an-adr-marked-as-not-applying-in-portuguese-still-seeds-the-related-adrs.md),
   the Backlog Entry recorded by the operator on 2026-11-10 from Fiscus's
   ai-memory message whose id begins `01a1c323`: the 76 ADR lines behind the
   cascade use "não se aplica", and Fiscus did not ask for the change.

At the authoring base `docs/backlog/` holds one other open entry, the
Catalog Skill Bump's stale catalog digest, which Spec 0378 adopts and which
shares no context with the strict check; `docs/findings/` holds none.

## Prerequisites

None in `requires`. Authored against `b6cd0f14` (v0.137.9 plus the 0374,
0375, 0376 and 0377 plans). The queue delivers Specs 0373, 0374, 0375, 0377, 0376 and 0378
before this one. This Spec owns `internal/speccheck/authoring_proof.go`'s
marker pattern, the `SC-ADR-RELATED` fix text in
`internal/speccheck/citations.go`, the `write-techspec` and `write-tasks`
skills and the `SC-ADR-ALTERNATIVES-UNRELATED` entries of the Roundfix `spec`
reference and guide. It shares with the earlier Specs only files that the
docs Task regenerates or edits narrowly:

- the owned Roundfix skill version fields (`.agents/skills/roundfix/SKILL.md`,
  `skills/roundfix/SKILL.md`) and `skills/testdata/owned-skill-versions.json`,
  which Specs 0373, 0374, 0375 and 0376 also raise, re-recorded only with the
  repository's record command;
- `CONTEXT.md`, where Specs 0373 to 0377 change other terms; this Spec edits
  one sentence of **Spec Consistency Check**;
- `docs/references/behavior-surfaces.json`, which Specs 0373 to 0376 also
  re-record, regenerated as the docs Task's last step.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted or changed; the check reads a declaration Specs already write.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — `roundfix spec check` reads local
  files only. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0377 reads the not-applicable
  exception from one English token, and this Spec amends it with ADR-0389
  (ADR-0377: "begins with `not applicable`, matched without regard to case").
  ADR-0093 and ADR-0168 are not applicable: they cite ADR-0377 for the
  related-ADR rule and its horizon, which keep their meaning. Source:
  `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — no protected tooling mutation proposed or
  authorized: no tooling configuration, script, ignore file, plugin
  declaration or version pin changes. The Governed Paths task_01 touches are
  owned skill files, measured with `speccheck.GovernedPath`, which the
  maintainer's cycle grants cover. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0379-an-adr-marked-as-not-applying-in-portuguese-draws-in-no-citing-adr/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/write-tasks/SKILL.md`,
  `.agents/skills/write-techspec/SKILL.md`,
  `.agents/skills/write-techspec/references/techspec-template.md`,
  `skills/roundfix/SKILL.md`, `skills/write-tasks/SKILL.md`,
  `skills/write-techspec/SKILL.md`,
  `skills/write-techspec/references/techspec-template.md`.

## Goals

- A Spec written in Portuguese whose ADR Alternatives lines begin with "não
  se aplica" draws in no ADR that cites only those ADRs, with no rewrite.
- The tokens stay a closed list that every checkout reads the same way, and
  any other wording still counts as an applied ADR.
- The authoring skills, the `spec` guide and reference, the glossary and the
  decision record name the accepted tokens.

## User Stories

1. As an adopter that writes its Specs in Portuguese, I want "não se aplica"
   at the start of an ADR Alternatives line to mark the ADR as not applying,
   so that the strict check stops asking for the ADRs that cite it.
2. As the maintainer, I want the accepted tokens to be a fixed list in the
   check, so that a Spec passes or fails the same way in every checkout.
3. As a Spec author, I want the skill and the check's fix text to name the
   tokens, so that I know which words the check reads.

## Core Features

1. The `SC-ADR-RELATED` seed rule treats a listed ADR as not applying when
   its ADR Alternatives relation begins with `not applicable`, `não se
   aplica` (precomposed or decomposed `ã`) or `nao se aplica`, without regard
   to case, ending at the end of the line or before a character that is not
   an ASCII letter, digit or underscore.
2. Any other relation, including `não se aplicam`, `não aplicável`, `n/a` or
   a token that is not at the start, still counts as an applied ADR.
3. The `SC-ADR-RELATED` fix text names `not applicable` and `não se aplica`.
4. The `write-techspec` skill and its TechSpec template, the `write-tasks`
   skill, the Roundfix `spec` reference, the `spec` command guide, the
   Context-Driven Development guide and the glossary's **Spec Consistency
   Check** name the Portuguese token; ADR-0377 carries the amendment notice.

## User Experience

Before this Spec, a TechSpec line `- ADR-0093: não se aplica; a regra segue
igual.` leaves ADR-0093 seeding the closure, and `roundfix spec check <slug>
--strict` reports `[gap] SC-ADR-RELATED` for every unlisted accepted ADR that
cites ADR-0093, with the fix "begin its ADR Alternatives line in the TechSpec
with \"not applicable\"". After it, the same line marks ADR-0093 as not
applying, those gaps disappear, and the fix text of the gaps that remain reads
"with \"not applicable\" or \"não se aplica\"". `Não se aplica`, `NÃO SE
APLICA` and `nao se aplica` behave the same; `não se aplicam` does not.

## Non-Goals / Out of Scope

- A configurable token, a third language or synonyms (see ADR-0389).
- Reading a marker in the PRD's Active ADR obligations row, which ADR-0377
  keeps rejected (ADR-0377: "because the row is free prose with no boundary
  between ADRs").
- The Project Constraints rows' `applicable` / `not applicable` parsing of
  `SC-CONSTRAINT-*`, which no adopter reported and which this Spec does not
  change.
- A finding for a line that resembles a token but is not one; such a line
  counts as applied, as any other relation does.
- Rewriting the ADR-0377 notices already written into ADR-0093 and ADR-0168:
  they record that amendment and stay true for the English token.
- Any change to the Makefile, workflows, `go.mod`, `go.sum`, lint or formatter
  configuration or the Baseline.

## Success Metrics

1. Success Metric: in a temporary repository with three accepted ADRs, where
   the Spec lists the first two and marks the second with each accepted
   Portuguese form, `roundfix spec check --strict` reports no `SC-ADR-RELATED`
   for the third, which cites only the second; the binary of the Delivery
   Base reports it.
2. Success Metric: `não se aplicam`, `não aplicável`, `n/a` and a token that
   is not at the start still report the gap, and the English token keeps its
   behavior.
3. Success Metric: the skills, the reference, the two guides, the glossary
   and ADR-0377 name the Portuguese token, and the owned skill versions are
   re-recorded.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The adopted Backlog Entry, recorded by the operator from Fiscus's message
  of 2026-11-10 about a repository this Spec did not build: "the 76 ADR lines
  behind the cascade Spec 0367 fixed use \"não se aplica\"".
- The Unicode Standard Annex #15, read through Exa on 2026-11-10
  (https://unicode.org/reports/tr15/): canonical equivalence is "a
  fundamental equivalency between characters or sequences of characters
  which represent the same abstract character, and which when correctly
  displayed should always have the same visual appearance and behavior",
  with the combining sequence example `Ç ↔ C+◌̧`; the decomposed `ã` is such a
  sequence.
- The Go `regexp/syntax` documentation, read through Exa on 2026-11-10
  (https://pkg.go.dev/regexp/syntax): `\b` is "at ASCII word boundary (\w on
  one side and \W, \A, or \z on the other)", which is why `não se aplicam`
  is not a token.

## Unreachable Acceptance

None.

## Validation

The decisions were checked on 2026-11-10 against these sources.

- Reading at `b6cd0f14`: `appliedADRNumbers` (`internal/speccheck/citations.go`)
  removes a listed ADR from the seeds when `notApplicableRelationPattern`
  (`internal/speccheck/authoring_proof.go`) matches its relation; the pattern
  is the English token alone; the fix text of `SC-ADR-RELATED` names only it,
  and `TestAnAppliedADRStillDrawsInItsCitingADR` pins that text.
- A prototype in a `git archive` copy: two new tests, one with six Portuguese
  forms and one with four guards, fail on the starting tree (every Portuguese
  form draws in the citing ADR) and pass with the change, and the whole
  `internal/speccheck` package passes. Six sabotages each fail a named
  subtest: the English-only pattern, no case folding, no decomposed form, no
  unaccented form, no word boundary, and the old fix text.
- Published sources (Exa): the two passages quoted above.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd query "Portuguese marker não se aplica ADR alternatives strict check language of specs" --all --files --min-score 0.3`,
  which returned the adopted Backlog Entry, ADR-0377 and this repository's
  CHANGELOG, all mirrored from this repository; none adds a decision.
- ai-memory: this session was told by the operator to send no message, and
  it read none; the operator triages the inbox.

## Glossary

- changes: **Spec Consistency Check**

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle BB batch 2 brief for Spec 0379 of 2026-11-10, the adopted Backlog Entry, and the maintainer's answers below.

The maintainer's answers (AskUserQuestion via `/roundfix --follow-up`,
2026-11-10 ~20:55–21:00 BRT), verbatim:

- Question "Quais Backlog Entries viram Spec agora, depois da 0374? As que
  você não marcar ficam abertas no backlog." → selected, among four,
  "Marcador \"não se aplica\"" — option text: "Aceitar o marcador em
  português nas ADR Alternatives. O Fiscus não pediu e não está travado."
- Question "Com que autorizações eu toco sem supervisão as seis Specs novas:
  0374 (capacidade da máquina), 0375 (orçamento de núcleos), 0376 (skills de
  usuário velhas), 0377 (portão do CI), 0378 (digest do bump de skills) e
  0379 (marcador \"não se aplica\")?" → "Mesmas do BB, teto minor
  (Recommended)" — option text: "Escrever, entregar pela fila, merge e
  release de cada Spec assim que pronta (patch nas correções, minor nas
  features; 0374 e 0375 saem juntas), teto de 2 Tasks corretivas, CI na 0377
  e Makefile na 0378 só nos arquivos que elas exigirem, eventos ao Jarvis e
  resposta ao Vortex quando a 0374 e a 0375 saírem."

Taken from those and the evidence, with no further interview:

- A closed list in the check, not a Project Config key: every checkout reads
  a Spec the same way. See ADR-0389.
- The unaccented `nao se aplica` is accepted: it is the same token typed
  without the accent, and no other phrase starts with it.
- The decomposed `ã` is accepted: it is canonically equivalent to the
  precomposed one and looks the same, so a byte match alone would fail a line
  the author cannot tell apart.
- Case is ignored, as for the English token.
- The plural `não se aplicam` is not a token: one line names one ADR, and a
  closed list is only closed while each member is exact.
- The skills keep `not applicable` as the first form and accept the
  Portuguese one for a Spec written in Portuguese.
- One docs Task (decision record, glossary, skills, guides) and one code Task
  that names those guides through its dependency, then the QA gate. The
  release is the operator's minor right after the merge; no Task releases.

## Open Questions

None.
