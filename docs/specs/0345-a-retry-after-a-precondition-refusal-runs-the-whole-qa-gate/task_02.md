---
task: task_02
spec: 0345-a-retry-after-a-precondition-refusal-runs-the-whole-qa-gate
status: completed
type: docs
complexity: medium
---

# Task 02: The implement reference, the configuration guide and the glossary state the import rules

## Overview

The Roundfix skill's implement reference names "an imported pass" only in
its Format Command paragraph, so Pantheon could not tell what a retry
imports, why its refusal came back, or how to start a pass with no previous
report. After this Task both copies of the reference, the configuration
guide, `CONTEXT.md` and ADR-0355 state the QA Report Import, that it never
imports a Precondition Refusal, what else it never brings in, how to start a
pass with no previous report, and that the Format Command covers the
refusal report. This answers the last request of the report reported by
pantheon on 2026-11-03
([references](references/2026-11-03-pantheon-retry-reimports-a-precondition-refusal.md)).

## Requirements

1. MUST edit the canonical `.agents/skills/roundfix/references/implement.md`
   after the Format Command paragraph of step 5, under no `### QA settlement`
   heading, exactly as `_techspec.md` → API Contract 4 states, with these
   phrases verbatim:
   - "the Precondition Refusal report included";
   - "never imports a Precondition Refusal";
   - "precondition refusal: <path>";
   - "a pass with no previous report".
   The list of what the import never brings in names each item of API
   Contract 4.3, and the path to a pass with no previous report says it is a
   destructive Git operation that needs the maintainer's authority.
2. MUST add to the Format Command section of
   `docs/user-guide/configuration.md` the sentence of API Contract 4.1, with
   the phrase "the Precondition Refusal report included", changing no other
   line.
3. MUST, through `domain-modeling`, add **QA Report Import** to `CONTEXT.md`
   as its own `**QA Report Import**:` entry with an `_Avoid_:` line, defined
   as API Contract 4.2 and 4.3 state and linking ADR-0194 and ADR-0355, and
   revise **Precondition Refusal** with the sentence "A later Run never
   imports it as its previous report."; the phrase "a later Run never imports
   it" must appear.
4. MUST create
   `docs/adr/0355-a-precondition-refusal-is-never-imported-as-a-prior-qa-pass.md`
   from `_techspec.md` → ADR-0355, with the lifecycle overlay of
   `docs/agents/docs-layout.md`.
5. MUST copy the canonical skill with `make skills-sync`, then raise both
   version fields of `.agents/skills/roundfix/SKILL.md` and its mirror with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
   -record-skill-versions -count=1`. Never hand-edit a mirror or a version
   field.
6. MUST re-record the Behavior Surface Record as the LAST step, after the
   last edit of any guide or reference, with `go test -count=1 -tags
   docscontract ./internal/docscontract -run
   '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`, and commit
   `docs/references/behavior-surfaces.json`.
7. MUST NOT touch any Go file, any other reference or skill, any Baseline
   asset, `docs/agents/`, any `### QA settlement` section, the Makefile, CI
   workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] Write the import rules into the canonical implement reference.
- [ ] Add the configuration guide sentence.
- [ ] Add and revise the glossary entries and create ADR-0355.
- [ ] Sync, record the skill version, re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] Both copies of the implement reference say what the import brings in, that it never imports a Precondition Refusal, what else it never brings in, and how to start a pass with no previous report.
- [ ] The configuration guide, `CONTEXT.md` and ADR-0355 state the rules, the mirrors equal their canonical copies, and the skill version and Behavior Surface records pass.

## Context

- instruction: `.agents/skills/domain-modeling/SKILL.md`
- interface: `.agents/skills/roundfix/references/implement.md`
- interface: `skills/roundfix/references/implement.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`
- creates: `docs/adr/0355-a-precondition-refusal-is-never-imported-as-a-prior-qa-pass.md`

## Consumers

- The owned-skill version record and the Behavior Surface Record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `for f in .agents/skills/roundfix/references/implement.md skills/roundfix/references/implement.md; do t="$(tr -s '[:space:]' ' ' < "$f")"; printf '%s' "$t" | grep -qF -- 'the Precondition Refusal report included' && printf '%s' "$t" | grep -qF -- 'never imports a Precondition Refusal' && printf '%s' "$t" | grep -qF -- 'precondition refusal: <path>' && printf '%s' "$t" | grep -qF -- 'a pass with no previous report' || { printf 'missing import rule in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- 'the Precondition Refusal report included' || { echo 'missing format sentence in configuration.md' >&2; exit 1; }; grep -F -- '**QA Report Import**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*QA Report Import\*\*(:| —( |$))' || { echo 'missing QA Report Import in CONTEXT.md' >&2; exit 1; }; grep -F -- '**Precondition Refusal**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Precondition Refusal\*\*(:| —( |$))' || { echo 'missing Precondition Refusal in CONTEXT.md' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'a later Run never imports it' || { echo 'missing Precondition Refusal revision' >&2; exit 1; }; adr=docs/adr/0355-a-precondition-refusal-is-never-imported-as-a-prior-qa-pass.md; test -f "$adr" && grep -q '^status: accepted' "$adr" || { echo 'missing ADR-0355' >&2; exit 1; }; cmp .agents/skills/roundfix/references/implement.md skills/roundfix/references/implement.md && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task no reference, guide or glossary entry names the rules and ADR-0355 does not exist, so the first check fails; afterwards they do, the mirrors equal their canonical copies, and the skill version and documentation contracts pass.

## References

- `_prd.md` → Goals; Core Feature 5; Success Metric 5; Glossary
- `_techspec.md` → Current behavior 6; API Contract 4; API Contract 5; API Contract 6; ADR-0355; Build Order 2

## Result

Skills loaded: implement-task, domain-modeling, tech-writer, evidence-gate

- Added the Format Command coverage sentence, QA Report Import rules, refusal
  exclusion, complete non-import list, and destructive no-previous-report path
  to the canonical implement reference; `make skills-sync` copied the same
  content to the mirror.
- Added the required configuration-guide sentence, the `QA Report Import`
  glossary entry and its `_Avoid_` line, the later-Run Precondition Refusal
  rule, and accepted ADR-0355 with the lifecycle overlay.
- Focused checks:
  - `rtk proxy git diff --check`: passed.
  - `rtk proxy make skills-sync`: passed.
  - `rtk proxy go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1`: passed; both Roundfix skill versions and the owned-skill record were raised by the recorder.
  - SHA-256 inspection of both canonical/mirror `SKILL.md` files and both
    implement references: matching pairs.
  - `rtk proxy go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`: passed; the Behavior Surface Record was re-recorded last.

Acceptance evidence:

- Both implement-reference copies contain the Format Command rule, the
  `QA Report Import` behavior, `never imports a Precondition Refusal`, the
  complete exclusion list, `precondition refusal: <path>`, and the destructive
  no-previous-report procedure.
- `docs/user-guide/configuration.md`, `CONTEXT.md`, and ADR-0355 state the
  import and refusal rules; generated skill versions, mirror hashes, and the
  Behavior Surface Record were updated by their sanctioned commands.
- Verification feedback root cause: the glossary contract rejected the
  compound phrase `Daemon QA Report` because `CONTEXT.md` defines `Daemon`
  and `QA Report` separately and the contract does not split a one-word term
  from a multiword term. The glossary now says `QA Report commit the Daemon
  produced`, preserving the API rule without adding an out-of-scope term.
- Focused feedback repair check:
  - `rtk proxy go test -count=1 -tags docscontract ./internal/docscontract -run '^TestGlossaryPhrasesAreDefinedTerms$'`: passed.
  - `rtk proxy git diff --check`: passed; canonical and mirror hashes still
    match for both the implement reference and `SKILL.md`.

Not run:

- The commands in this Task's `## Verification` section and the repository
  Verification remain for Daemon settlement, per the assigned execution mode.
