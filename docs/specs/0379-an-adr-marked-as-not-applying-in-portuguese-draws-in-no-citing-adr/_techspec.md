---
spec: 0379-an-adr-marked-as-not-applying-in-portuguese-draws-in-no-citing-adr
prd: _prd.md
created: 2026-11-10
---

# An ADR marked as not applying in Portuguese draws in no citing ADR — Technical Spec

## Executive Summary

`notApplicableRelationPattern` in `internal/speccheck/authoring_proof.go`
becomes `(?i)^(?:not applicable|n(?:\x{00E3}|a\x{0303}|a)o se aplica)\b`: the
English token, the Portuguese token with a precomposed or decomposed `ã`, and
its unaccented spelling, case-insensitive, ending at an ASCII word boundary.
`appliedADRNumbers` keeps calling it unchanged, so the related-ADR seeds lose
every ADR whose line begins with one of those tokens. The `SC-ADR-RELATED` fix
text in `internal/speccheck/citations.go` names both tokens. ADR-0389 records
the closed list and amends ADR-0377; the docs Task writes the notice into
ADR-0377 and aligns the skills, the reference, the guides and the glossary.
The trade-off accepted is a list that a third language must amend, instead of
a configuration key that would let two checkouts judge one Spec differently.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted or changed.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local file reads only. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0377 (only an applied ADR seeds the
  related-ADR candidates; the not-applicable token), amended by ADR-0389 (the
  closed list of tokens). ADR-0093 and ADR-0168 are not applicable: they cite
  ADR-0377 for the rule and its horizon, which keep their meaning. Source:
  `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — no protected tooling mutation proposed or
  authorized; the Governed Paths task_01 touches are owned skill files, which
  the maintainer's cycle grants cover. Source:
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

## System Architecture

`detectADRRelated` (`internal/speccheck/citations.go`) builds the related-ADR
seeds with `appliedADRNumbers(listed, <TechSpec path>)`, which reads each line
of the TechSpec's `## ADR Alternatives` section with `adrAlternativePattern`
and deletes the ADR from the seeds when `notApplicableRelationPattern` matches
the relation. Nothing else reads the pattern: the Project Constraints parser in
`constraints.go` has its own `not applicable` applicability prefix, which this
Spec leaves alone. The change is the pattern, the fix text, and their tests;
no new function, type or file in production code.

## Implementation Design

### Interfaces

```go
notApplicableRelationPattern = regexp.MustCompile(
	`(?i)^(?:not applicable|n(?:\x{00E3}|a\x{0303}|a)o se aplica)\b`)
```

`\x{00E3}` is `ã`, `a\x{0303}` is `a` with the combining tilde; Go's `(?i)`
folds `ã` and `Ã`; `\b` is the ASCII word boundary.

The `SC-ADR-RELATED` fix becomes, with only the token clause changed:
`Account for ADR-<n> in the Active ADR obligations row in <prd> by listing it or recording why it does not apply; when it does not apply, begin its ADR Alternatives line in the TechSpec with "not applicable" or "não se aplica", so the ADRs that cite only it are not drawn in.`

### Data Models

No changes. Specs, ADRs and findings keep their formats; finding codes,
severities and locations are unchanged.

### API Contracts

1. API Contract: seed tokens. A listed ADR is not a seed when its ADR
   Alternatives relation begins with `not applicable`, `não se aplica`
   (precomposed or decomposed `ã`) or `nao se aplica`, without regard to
   case, followed by the end of the line or a character outside
   `[A-Za-z0-9_]`; every other relation keeps it a seed.
2. API Contract: fix text. The `SC-ADR-RELATED` fix names `"not applicable"
   or "não se aplica"`; summary, severity, code and locations keep their
   bytes.

### Surface Transcripts

None. The command output changes only in the fix line of `SC-ADR-RELATED`,
named verbatim in Interfaces and API Contract 2; which findings appear is
covered by API Contract 1.

## Coverage Map

- Goal 1, User Story 1, Core Feature 1, Success Metric 1 → API Contract 1;
  `notApplicableRelationPattern`; Build Order 2.
- Goal 2, User Story 2, Core Feature 2, Success Metric 2 → API Contract 1;
  `TestAMarkerOutsideTheClosedListStillDrawsIn`; ADR-0389; Build Order 1, 2.
- User Story 3, Core Feature 3 → API Contract 2; Build Order 2.
- Goal 3, Core Feature 4, Success Metric 3 → the docs Task: ADR-0377 notice,
  `CONTEXT.md`, skills, reference and guides; Build Order 1.

Deviations from the PRD: None.

## Absence Evidence

1. A Portuguese not-applicable token — command: `rg -n 'se aplica' internal/speccheck`; result: no match at `b6cd0f14`; the only token is `(?i)^not applicable\b` in `authoring_proof.go`.
2. The tests of the Portuguese forms and guards — command: `rg -n 'Portuguese|se aplica' internal/speccheck --glob '*_test.go'`; result: no match at `b6cd0f14`; `adr_related_seed_test.go` covers `not applicable; unchanged.` and `Not applicable — unchanged.` only.

## ADR Alternatives

- ADR-0377: amended by ADR-0389; the token list grows from one English token
  to a closed English and Portuguese list in the same TechSpec line, which is
  not the PRD-row marker ADR-0377 rejected, and its other rejected options (a
  fifth Project Constraints row, a round limit, removing the check) stay
  rejected.
- ADR-0389: the decision this Spec implements; it rejects a configuration
  key, a skill-only rule, Unicode normalization, synonyms and the PRD-row
  marker, and the design takes none of them.
- ADR-0093: not applicable; the check still compares written declarations
  and never judges relevance.
- ADR-0168: not applicable; the horizon is unchanged.

## Integration Points

None. The check reads local files; no external system is touched.

## Testing Approach

The seam is `speccheck.Check` over a temporary repository, through the
existing fixture `writeADRRelatedSeedFixture` (ADR 0002 cites 0001, ADR 0003
cites 0002, the Spec lists 0001 and 0002, and the relation under test is
the second fixture ADR's line), which observes the seeds `appliedADRNumbers` builds with
the pattern.

- `TestAPortugueseNotApplicableADRDrawsInNoCitingADR` (new): subtests
  `accented` (`não se aplica; …`), `capitalized` (`Não se aplica — …`),
  `upper case` (`NÃO SE APLICA: …`), `unaccented` (`nao se aplica; …`),
  `decomposed` (`não se aplica; …`) and `line end` (`não se aplica`
  alone); each expects no `SC-ADR-RELATED`, `SC-ADR-UNLISTED` or
  `SC-ADR-ALTERNATIVES-UNRELATED`.
- `TestAMarkerOutsideTheClosedListStillDrawsIn` (new): subtests `plural`
  (`não se aplicam; …`), `not leading` (`mantida; a alternativa rejeitada não
  se aplica.`), `abbreviation` (`n/a; unchanged.`) and `other wording` (`não
  aplicável; …`); each expects the one `SC-ADR-RELATED` gap for the third fixture ADR.
- `TestAnAppliedADRStillDrawsInItsCitingADR` (existing, declared break): its
  `wantFix` gains ` or "não se aplica"`.
- `TestANotApplicableADRDrawsInNoCitingADR` and
  `TestWithoutATechSpecEveryListedADRSeeds` stay unchanged.

Proven during authoring on `git archive` copies: the six Portuguese subtests
fail on the starting tree and pass with the change; dropping the Portuguese
alternative, `(?i)`, `a\x{0303}`, the unaccented `a`, or `\b`, or reverting the
fix text, each fails a named subtest.

## Build Order

1. Decision record notice, glossary, skills, reference and guides (task_01):
   the Amended by notice in ADR-0377, the **Spec Consistency Check**
   sentence, the `write-techspec` skill and template, the `write-tasks` skill,
   the Roundfix `spec` reference with `make skills-sync` and the owned skill
   version record, the `spec` command guide and the Context-Driven
   Development guide, then the Behavior Surface Record re-recorded last. It
   runs after Specs 0373 to 0378 merged, which the queue order ensures, and
   reads no artifact of theirs.
2. The pattern, the fix text and their tests (task_02) (depends on: 1).
3. The QA gate (depends on: 2).

## Risks & Considerations

- A relation the author meant as applied that starts with a listed token
  would stop seeding; the tokens are explicit refusals, so this is unlikely,
  and ADR-0389 states it.
- `\b` is ASCII-only in Go, so a token followed by a non-ASCII letter (for
  example `não se aplicaç`) would match; no Portuguese word continues
  `aplica` with a non-ASCII letter, and the guard tests pin the plural.
- Shared files with the earlier Specs (owned skill version fields and record,
  `CONTEXT.md`, the Behavior Surface Record) are edited narrowly or
  regenerated by their commands at the Task's time.

## Glossary

None.

## Decisions

- Pattern, not a token table: one regular expression keeps the single call
  site in `appliedADRNumbers` unchanged.
- No Unicode normalization module: the one decomposable character is named in
  both forms (ADR-0389).
- Docs first, then code, so the code Task names its guides through its
  dependency.
