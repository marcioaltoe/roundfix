---
task: task_01
spec: 0379-an-adr-marked-as-not-applying-in-portuguese-draws-in-no-citing-adr
status: completed
type: docs
complexity: low
---

# Task 01: The decision record, glossary, skills and guides name the Portuguese not-applicable token

## Overview

On 2026-11-10 the operator recorded that an ADR marked as not applying in
Portuguese still seeds the related ADRs
([`references/2026-11-10-an-adr-marked-as-not-applying-in-portuguese-still-seeds-the-related-adrs.md`](references/2026-11-10-an-adr-marked-as-not-applying-in-portuguese-still-seeds-the-related-adrs.md),
recorded by the operator from Fiscus's message). ADR-0389 decides the closed
list of tokens and amends ADR-0377. This Task writes the amendment notice into
ADR-0377 and makes the glossary, the `write-techspec` and `write-tasks`
skills, the Roundfix `spec` reference and the two guides name the Portuguese
token, so the code Task that follows changes behavior they already describe.

## Requirements

1. MUST, in `docs/adr/0377-only-an-adr-a-spec-applies-or-changes-draws-in-the-adrs-that-cite-it.md`:
   - set `updated_at: 2026-11-11T00:00:00Z`;
   - insert, as its own paragraph immediately after the Decision paragraph
     that ends `never judges whether an ADR is relevant (ADR-0093).`, the
     notice
     `**Amended by [ADR-0389](0389-an-adr-alternatives-line-marks-an-adr-as-not-applying-in-english-or-portuguese.md) (2026-11-11).** The leading token is read from a closed list: `not applicable`, `não se aplica` with a precomposed or decomposed `ã`, and `nao se aplica`, matched without regard to case and ending before a character that is not an ASCII letter, digit or underscore. Everything else in this decision still holds.`;
   - append to the Considered Option that begins `A marker phrase inside the
     PRD's Active ADR obligations row` the sentence
     `ADR-0389 keeps this option rejected: it reads the Portuguese token only in the TechSpec's ADR Alternatives line, from a closed list every checkout shares.`
   Nothing else in the ADR changes.
2. MUST, through `domain-modeling`, revise **Spec Consistency Check** in
   `CONTEXT.md`: in its sentence that says an ADR whose `## ADR Alternatives`
   line in the TechSpec begins with `not applicable` stays listed, insert
   `, or with `não se aplica` in a Spec written in Portuguese,` after
   `not applicable`, and add
   `[ADR-0389](docs/adr/0389-an-adr-alternatives-line-marks-an-adr-as-not-applying-in-english-or-portuguese.md)`
   beside the ADR-0377 link in its closing parenthesis, separated by `, `.
   The definition line and every other term keep their bytes.
3. MUST, in `.agents/skills/write-techspec/SKILL.md` (rule "Rejected
   alternatives"), replace the sentence that begins `For an ADR the Spec lists
   only to record why it does not apply` with:
   `For an ADR the Spec lists only to record why it does not apply, begin the relation with `not applicable`, or with `não se aplica` in a Spec written in Portuguese: the ADR stays listed, and `SC-ADR-RELATED` does not ask for the ADRs that cite only it. The check reads these tokens without regard to case and also accepts `nao se aplica` without the accent; any other wording, such as `não se aplicam` or `n/a`, counts as an applied ADR.`
4. MUST, in `.agents/skills/write-techspec/references/techspec-template.md`,
   change the ADR Alternatives comment's last sentence to
   `Begin the relation with `not applicable`, or with `não se aplica` in a Spec written in Portuguese, for an ADR the Spec lists only to record why it does not apply.`
5. MUST, in `.agents/skills/write-tasks/SKILL.md`, change the clause
   `begin a not-applying ADR's line under the TechSpec's `## ADR Alternatives` with `not applicable` so it draws in no ADR that cites only it.`
   to
   `begin a not-applying ADR's line under the TechSpec's `## ADR Alternatives` with `not applicable`, or `não se aplica` in a Spec written in Portuguese, so it draws in no ADR that cites only it.`
6. MUST, in the `SC-ADR-ALTERNATIVES-UNRELATED` entry of
   `.agents/skills/roundfix/references/spec.md` and of
   `docs/user-guide/commands/spec.md`, change
   `A line whose relation begins with `not applicable` keeps the ADR listed`
   to
   `A line whose relation begins with `not applicable` or `não se aplica`, matched without regard to case and also without the accent, keeps the ADR listed`,
   keeping the rest of each entry.
7. MUST, in `docs/user-guide/context-driven-development.md`, change
   `begins with `not applicable` stays listed, and an ADR that cites only such ADRs`
   to
   `begins with `not applicable`, or with `não se aplica` in a Spec written in Portuguese, stays listed, and an ADR that cites only such ADRs`.
8. MUST run `make skills-sync`, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   which raises the version fields of `roundfix`, `write-tasks` and
   `write-techspec` and records them; never write a version number by hand.
9. MUST, as the Task's last step, after every guide and skill edit,
   re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
   (re-record after the last guide edit).
10. MUST NOT touch the `### QA settlement` section of any skill, ADR-0093,
    ADR-0168, ADR-0389, or any other glossary term.

## Subtasks

- [x] Write the amendment notice and the Considered Option note into ADR-0377.
- [x] Revise **Spec Consistency Check** through `domain-modeling`.
- [x] Align the `write-techspec` skill and template, the `write-tasks` skill and the `spec` reference, then sync and record the skill versions.
- [x] Align the two guides and re-record the Behavior Surface Record last.

## Acceptance Criteria

- [x] Every prescribed sentence is present in each file and its mirror, and the mirrors equal their canonical copies.
- [x] The three owned skills carry raised, recorded versions, and the Behavior Surface Record is current.
- [x] ADR-0377 carries the dated notice and its `updated_at`; the glossary keeps its definition line.

## Context

- interface: `docs/adr/0377-only-an-adr-a-spec-applies-or-changes-draws-in-the-adrs-that-cite-it.md`
- interface: `CONTEXT.md`
- interface: `.agents/skills/write-techspec/SKILL.md`
- interface: `skills/write-techspec/SKILL.md`
- interface: `.agents/skills/write-techspec/references/techspec-template.md`
- interface: `skills/write-techspec/references/techspec-template.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `.agents/skills/roundfix/references/spec.md`
- interface: `skills/roundfix/references/spec.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/spec.md`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `docs/references/behavior-surfaces.json`
- instruction: `docs/adr/0389-an-adr-alternatives-line-marks-an-adr-as-not-applying-in-english-or-portuguese.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`

## Consumers

- The owned skill copies and their recorded versions — consumers: `./skills`
- The `spec` command guide and the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `for f in .agents/skills/write-techspec/SKILL.md skills/write-techspec/SKILL.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'any other wording, such as' && tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'counts as an applied ADR' || { printf 'missing phrase in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/write-techspec/references/techspec-template.md skills/write-techspec/references/techspec-template.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'in a Spec written in Portuguese, for an ADR the Spec lists only' || { printf 'missing phrase in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/write-tasks/SKILL.md skills/write-tasks/SKILL.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'in a Spec written in Portuguese, so it draws in no ADR that cites only it' || { printf 'missing phrase in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/spec.md skills/roundfix/references/spec.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'matched without regard to case and also without the accent, keeps the ADR listed' || { printf 'missing phrase in %s\n' "$f" >&2; exit 1; }; done` — expected: every skill copy and mirror carries its prescribed sentence; fails before the work.
- `GOMAXPROCS=2 go test -count=1 -timeout 10m ./skills -run '^(TestEveryOwnedSkillVersionIsRecorded|TestAuthorialSkillSync)$' && tr -s '[:space:]' ' ' < skills/write-techspec/SKILL.md | grep -qF -- 'counts as an applied ADR'` — expected: the mirrors equal their canonical copies and the raised versions are recorded; fails before the work because the sentence is absent.
- `grep -F -- '**Spec Consistency Check**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Spec Consistency Check\*\*(:| —( |$))' && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'in a Spec written in Portuguese, stays listed but opens no' && tr -s '[:space:]' ' ' < docs/adr/0377-only-an-adr-a-spec-applies-or-changes-draws-in-the-adrs-that-cite-it.md | grep -qF -- 'Amended by [ADR-0389](0389-an-adr-alternatives-line-marks-an-adr-as-not-applying-in-english-or-portuguese.md) (2026-11-11).' && tr -s '[:space:]' ' ' < docs/adr/0377-only-an-adr-a-spec-applies-or-changes-draws-in-the-adrs-that-cite-it.md | grep -qF -- 'ADR-0389 keeps this option rejected' && tr -s '[:space:]' ' ' < docs/adr/0377-only-an-adr-a-spec-applies-or-changes-draws-in-the-adrs-that-cite-it.md | grep -qF -- 'updated_at: 2026-11-11T00:00:00Z'` — expected: the glossary keeps **Spec Consistency Check** and names the Portuguese token, and ADR-0377 carries the dated notice; fails before the work.
- `GOMAXPROCS=2 go test -count=1 -timeout 10m -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' && tr -s '[:space:]' ' ' < docs/user-guide/commands/spec.md | grep -qF -- 'matched without regard to case and also without the accent, keeps the ADR listed' && tr -s '[:space:]' ' ' < docs/user-guide/context-driven-development.md | grep -qF -- 'in a Spec written in Portuguese, stays listed, and an ADR that cites only such ADRs'` — expected: both guides name the token and the Behavior Surface Record is current; fails before the work.

## References

- `_prd.md` → Goals; User Stories 2-3; Core Feature 4; Success Metric 3; Glossary
- `_techspec.md` → ADR Alternatives; Build Order 1
- ADR-0389; ADR-0377

## Result

Skills loaded: implement-task, domain-modeling, tech-writer, evidence-gate

Implemented the ADR-0377 amendment, the Portuguese token wording in the
`Spec Consistency Check` glossary entry, both owned authoring skills and their
template/reference text, both user guides, and the generated skill/version
records. The Behavior Surface Record was re-recorded after the final skill and
guide edit.

Focused checks:

- `make skills-sync`: passed; the `skills/` copies were regenerated.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`: passed (`1 passed`); versions recorded as `roundfix 0.1.218`, `write-tasks 0.0.32`, and `write-techspec 0.0.19`.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`: passed (`1 passed`); executed last after all skill and guide edits.
- `git diff --check`: passed.
- `diff -r .agents/skills/<owned-skill> skills/<owned-skill>` for `roundfix`, `write-tasks`, and `write-techspec`: all three mirrors equal their canonical copies.
- Focused phrase and version inspections found the prescribed ADR notice, glossary links, Portuguese token wording, guide entries, raised versions, and coverage record change.

Acceptance evidence:

- Prescribed sentences and mirrors: the focused phrase checks found each required sentence, and all three owned-skill mirror comparisons were equal.
- Raised skill versions and Behavior Surface Record: both recorder tests passed; the version record contains the new entries and `docs/references/behavior-surfaces.json` changed.
- ADR-0377 and glossary: focused inspection found `updated_at: 2026-11-11T00:00:00Z`, the dated ADR-0389 notice, the Considered Option note, and the unchanged `**Spec Consistency Check**` definition heading.

The Task's declared `## Verification` commands were not run; the Daemon owns
that verification and Task status. No commit, push, or pull request was made.
