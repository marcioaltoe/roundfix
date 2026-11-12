---
spec: 0387-the-techspec-template-states-the-adr-alternatives-instruction-once
status: active
created: 2026-11-12
surfaces: [docs]
type: fix
---

# The TechSpec template states the ADR Alternatives instruction once

The `write-techspec` skill's TechSpec template repeats "Begin the relation
with" in the comment of its `## ADR Alternatives` section: one line ends with
the phrase and the next line begins with it again, so the instruction that
tells an author how to mark an ADR that does not apply reads as a stutter. Spec
0379 introduced the repetition, every adopter received it through `roundfix
baseline update` with `write-techspec` 0.0.19, and Argus's pre-PR review
reported it. After this Spec both copies of the template state the
instruction once, the skill's version rises so adopters receive the
correction through their next `roundfix baseline update`, and a contract test
fails whenever a line of any Roundfix-owned skill ends with a phrase of two or
more words that the next line begins with.

## Sources

One cross-project message, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-12-argus-the-techspec-template-repeats-begin-the-relation-with.md`](references/2026-11-12-argus-the-techspec-template-repeats-begin-the-relation-with.md),
   ai-memory message `01a1d159` of 2026-11-12, copied verbatim; reported by
   argus. It names lines 87 and 88 of
   `skills/write-techspec/references/techspec-template.md` at `726505ca`, the
   canonical copy under `.agents/skills/`, commit `c11c5883` (#833, Spec 0379)
   as the origin, and asks to remove the repetition at the source and publish
   a version of the skill. It asks for no reply.

At the authoring base `docs/backlog/` and `docs/findings/` hold no entry about
skill or template wording, so no other source shares this context.

## Prerequisites

None in `requires`. Authored against `805488c2` (v0.142.0 plus the plans of
Specs 0383 to 0386). The queue delivers 0383, 0385, 0384 and 0386 before this
Spec, which is last. The only file this Spec shares with them is
`skills/testdata/owned-skill-versions.json`: Specs 0383, 0385 and 0386 append
versions of the `roundfix` skill to it, and this Spec appends one version of
`write-techspec`, a different key. task_01 records the version with the
repository's record command on the tree it runs on and never writes a number.
No queued Spec declares a `write-techspec` file. The archived Spec 0379 wrote
the template line this Spec corrects; its Archive Record stays unchanged.

## Project Constraints

- Identifier strategy: not applicable — no identifier, record or schema is
  minted; the skill keeps its name and gains a version. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or
  network call is added; the change is skill text and a test that reads the
  repository's own files. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0389 applies unchanged: the
  template's comment carries its rule that an author begins the relation with
  `not applicable`, or with `não se aplica` in a Spec written in Portuguese,
  and the corrected comment keeps both tokens. ADR-0377 applies unchanged:
  only an ADR a Spec applies or changes seeds `SC-ADR-RELATED`, read from the
  TechSpec's own `## ADR Alternatives` lines and never from the template.
  ADR-0093 and ADR-0168, which cite ADR-0377, do not apply: they decide how
  Spec consistency is checked and when a related-ADR gap opens, and this Spec
  changes no check. Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  answer of 2026-11-12 quoted in the authorization record, on the cycle BD
  terms that keep the grant of every owned skill; the Governed Paths are the
  `write-techspec` skill's `SKILL.md` pair and both copies of its TechSpec
  template. No Makefile, workflow, Go module, lint or formatter change.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0387-the-techspec-template-states-the-adr-alternatives-instruction-once/_authorization.md`;
  bounded files: `.agents/skills/write-techspec/SKILL.md`, `.agents/skills/write-techspec/references/techspec-template.md`, `skills/write-techspec/SKILL.md`, `skills/write-techspec/references/techspec-template.md`.

## Goals

- The `## ADR Alternatives` comment of both copies of the TechSpec template
  states "Begin the relation with" once, followed by `not applicable` and the
  Portuguese token, with every other byte of the template unchanged.
- The `write-techspec` skill's version rises and is recorded, so every
  adopter receives the corrected template through its next `roundfix baseline
  update`.
- A contract test fails when a line of any Roundfix-owned skill, canonical or
  shipped, ends with a phrase of two or more words that the next line begins
  with, and passes on today's owned skills once the template is corrected.
- `SC-ADR-RELATED` keeps reading the not-applying token from the TechSpec's
  own lines, with no dependence on the template's wording.

## Core Features

1. The template comment. In `.agents/skills/write-techspec/references/techspec-template.md`
   the second line of the comment loses its leading "Begin the relation with
   ", so the comment reads "… Begin the relation with `not applicable`, or
   with `não se aplica` in a Spec written in Portuguese, for an ADR the Spec
   lists only to record why it does not apply."; `make skills-sync` writes the
   shipped mirror.
2. The skill version. The record command raises `write-techspec` to its next
   version in both `SKILL.md` copies and records it in
   `skills/testdata/owned-skill-versions.json`.
3. The guard. `skills/repeated_phrase_test.go` holds the detector of a phrase
   repeated across a line break, its unit cases, and the contract test that
   runs it over every file of every owned skill, canonical and shipped,
   naming the file, the line and the phrase.

## Non-Goals / Out of Scope

- Repeated words inside one line, and a single word repeated across a line
  break: "that that" and "had had" are grammatical English, and a one-word
  rule would need a list of exceptions; the guard checks phrases of two or
  more words.
- Vendored skills, guides, `CONTEXT.md`, ADRs and Baseline modules: the guard
  reads only the Roundfix-owned skills Argus received the defect in. A scan
  at the authoring base finds no repeated phrase in `docs/agents/`,
  `docs/user-guide/`, `CONTEXT.md` or `docs/adr/` either.
- Any change to `SC-ADR-RELATED`, its closed token list or its fix text, the
  `write-techspec` skill's `Rejected alternatives` rule, or another template
  line.
- Editing the Archive Record of Spec 0379 or Argus's installed copy; Argus
  receives the corrected copy from its next managed refresh.

## Success Metrics

1. Success Metric: in both copies of the template, the whitespace-normalized
   comment of `## ADR Alternatives` contains "rejected. Begin the relation
   with `not applicable`, or with `não se aplica` in a Spec written in
   Portuguese, for an ADR the Spec lists only to record why it does not
   apply." and "Begin the relation with Begin the relation with" appears
   nowhere; the two copies are byte-identical, and the diff against the
   Delivery Base changes only that comment line.
2. Success Metric: `TestNoOwnedSkillRepeatsAPhraseAcrossALineBreak` passes for
   the `canonical` and `shipped` copies of all twelve owned skills, and
   `TestRepeatedLinePhrasesFindsAPhraseRepeatedAcrossALineBreak` passes its six
   cases; at the authoring base the same scan reports exactly one repeated
   phrase, the one Argus reported, in each copy.
3. Success Metric: `write-techspec` carries its next version in both
   `SKILL.md` copies, recorded in `skills/testdata/owned-skill-versions.json`,
   and `go test ./skills` passes.
4. Success Metric: the `SC-ADR-RELATED` marker and seed tests in
   `internal/speccheck` pass unchanged, and no Go file reads the template's
   comment.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- Argus's report, under `references/`: a measured finding of `roundfix
  review` in another repository, after its own `roundfix baseline update`
  installed `write-techspec` 0.0.19: "A linha 87 termina com ela e a linha 88
  começa de novo com ela". The defect the guard must catch is the one an
  adopter found, not one this Spec designed.
- Matt Might, "3 shell scripts: Kill weasel words, avoid the passive,
  eliminate duplicates"
  (https://matt.might.net/articles/shell-scripts-for-passive-voice-weasel-words-duplicates/),
  read through Exa on 2026-11-12, section "A script to find lexical
  illusions": "Many readers are not aware that the / the brain will
  automatically ignore / a second instance of the word "the" / when it starts
  a new line." and "Duplicating words is a phenomenon of electronic
  composition. They seem to happen as cut and paste accidents". A repetition
  that starts a new line escapes the reader, so a mechanical check across the
  line break is what catches it; the article's script also compares words
  without regard to case.

## Unreachable Acceptance

- criterion: after the release that carries this Spec, an adopter's `roundfix
  baseline update` installs `write-techspec` at its new version, and its
  TechSpec template states "Begin the relation with" once.
  reason: it needs a released binary and another repository's managed
  refresh, which no Verification or QA row may run or read.
  satisfied-by: the operator's or Argus's first managed refresh after the
  release, recorded session evidence.

## Validation

The decisions were checked on 2026-11-12 against these sources.

- Characterization at `805488c2`: lines 86 to 88 of both template copies hold
  the comment, line 87 ends "Begin the relation with" and line 88 begins
  "Begin the relation with `not applicable`". `git show c11c5883` shows Spec
  0379 rewrote line 88 by prefixing the old token line with the phrase that
  line 87 already ended with. `write-techspec` is at 0.0.19 in both
  `SKILL.md` copies.
- Repeated phrases across the twelve owned skills at `805488c2`: a scan of
  every Markdown file under `.agents/skills/<name>/` and `skills/<name>/`
  (the prototype guard later read every file, the agent manifest included,
  with the same result), comparing
  the trailing words of each line with the leading words of the next without
  regard to case, finds, for phrases of two or more words, one hit per copy,
  the template line above. It also finds 14 lines per copy that a fenced
  command listing extends on the next line (`roundfix gc compact` followed by
  `roundfix gc compact --apply`), each line made only of the phrase, and
  Markdown table rows whose borders repeat; the guard exempts both. With
  single words included, the scan finds no other hit outside those two
  shapes. The same scan over `docs/agents/`, `docs/user-guide/`, `CONTEXT.md`
  and `docs/adr/` finds none, and over the vendored skills only box-drawing
  characters in one diagram.
- `SC-ADR-RELATED` reads the token through `notApplicableRelationPattern` in
  `internal/speccheck/authoring_proof.go` and `appliedADRNumbers` in
  `internal/speccheck/citations.go`, over the Spec's own TechSpec; the only Go
  file that names the template is `internal/speccheck/governed.go`, which
  lists it as a Governed Path, and the tests that read it
  (`TestAuthoringTemplatesUseTheBoundedFilesLabel`,
  `TestWriteTechSpecProjectConstraints`, the template contract in
  `internal/docscontract`) assert other sections.
- A prototype on a worktree at `805488c2`: the line edit, `make skills-sync`
  and the record command raised `write-techspec` to 0.0.20 and changed only
  the five files task_01 declares besides the new test; `go test ./skills`,
  `GOOS=windows go vet ./skills` and the documentation contracts
  `TestTheSkillCoverageMapIsCurrent` and `TestGlossaryPhrasesAreDefinedTerms`
  passed, so the Behavior Surface Record does not move. Seven sabotages each
  failed their named test (see `_techspec.md` → Testing Approach).
- Secondbrain, read only: `wiki/index.md`, `wiki/entities/argus.md` (Argus is
  a TypeScript desktop project that adopts the Roundfix Baseline) and a `qmd
  query "duplicated phrase write-techspec template ADR Alternatives not
  applicable" --all --files --min-score 0.3`, which returned no result.
  Nothing there changes the decision.
- Exa: the article quoted under Acceptance evidence. It flags any repeated
  word; this Spec narrows the rule to phrases of two or more words for the
  reason under Non-Goals.

## Glossary

None.

## Decisions

The maintainer's answers (AskUserQuestion, 2026-11-12 ~20:12 BRT), verbatim:

- Question: "A nova Spec 0387 corrige a frase repetida \"Begin the relation
  with\" no template da skill write-techspec (relato do Argus). Que
  autorização ela tem?" — answer: "Mesmas do BD, na v0.143.0 (Recommended)";
  option text: "Push, merge e release; entra na fila depois da 0386 e sai
  junto com ela na minor v0.143.0, chegando aos adotantes pelo baseline
  update."
- Question: "Até o fim da fila (0383 → 0385 → 0384 → 0386 → 0387), quando um
  item parar por falha técnica (QA reprovada, achado de review, teste
  instável no gate), o que eu faço?" — answer: "Consertar e seguir
  (Recommended)" (operator policy for the queue).

Autonomous route: no person present; decisions taken from the operator's cycle BD batch 3 brief for Spec 0387 of 2026-11-12, the answers above and the characterization above.

Taken from those sources, with no further interview:

- Delete the leading "Begin the relation with " of the second comment line,
  so the sentence that line 87 opens continues with the tokens; no other
  wording changes.
- The guard is general over every owned skill, not a phrase pin on one
  template: the same accident can recur in any owned file, and the scan
  shows it stays precise with two exemptions.
- It lives in the `skills` package, beside the other owned-skill contracts,
  and reads both the canonical and the shipped copy, because adopters
  receive the shipped one.
- `type: fix`; the Release Ceiling is `patch`, and the operator releases it
  with Spec 0386 in v0.143.0. One implementation Task and a QA gate; no ADR.

## Open Questions

None.
