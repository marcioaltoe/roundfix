---
spec: 0170-authoring-that-fails-before-dispatch
status: archived
created: 2026-09-25
surfaces: [backend, cli, docs]
archived: "2026-09-25"
source_slug: 0170-authoring-that-fails-before-dispatch
---


# Authoring that fails before dispatch

Sixteen Runs of Specs 0155 to 0169 ended Unresolved. Seven of them failed for
an authoring mistake that a file read could have caught before any Agent turn
was spent:

- Five Tasks edited a Governed Path that their `_authorization.md` `paths:` and
  the Tooling authority rows never named, so the Daemon's changed-path audit
  refused the commit after the work was done. `roundfix spec check` read the
  same Task files and passed them.
- Two Tasks changed CLI behavior while the skill or guide that describes it
  stayed behind, and a corrective Task followed the pre-PR review.

Five smaller defects in the same authoring path are recorded beside them:

- The PRD and TechSpec templates ask for `bounded paths:`, a label the checker
  refuses with `SC-TOOLING-UNBOUNDED`.
- A Verification that pipes a tool into `grep` inside a command substitution
  and tests the result for emptiness passes when the tool itself fails, and
  `SC-VERIFY-INVERTED-EXIT` approves it.
- A `None.` declaration is accepted only on a single line, so a normally
  wrapped reason is reported as `SC-METRIC-UNDECLARED` or
  `SC-CONTRACT-UNDECLARED`.
- An adopted Finding or Backlog Entry whose original still sits in
  `docs/findings/` or `docs/backlog/` passes `--strict`, although adoption is
  one move.
- No test pins the one-field and blank-field closure cases of an archived
  Finding, so a regression in that check would pass the suite.

This Spec moves each of those refusals to authoring and makes the templates and
authoring skills stop producing the mistakes.

## Project Constraints

- Identifier strategy: applicable — two new stable finding codes,
  `SC-TOOLING-UNDECLARED` and `SC-CLI-UNDOCUMENTED`, join the Spec Consistency
  Check vocabulary; no other identifier changes. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files only; no credential
  and no network call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0093 checks Spec consistency by
  citation and never by inference, so the new checks compare declared paths
  and never guess what a Task will edit; ADR-0094 skips a detector whose
  artifact is absent; ADR-0117 places a check with the stage that can produce
  the defect and keeps the commit audit as the gate row; ADR-0130 makes the
  audit judge Governed Paths; ADR-0131 keeps the tooling row about
  applicability; ADR-0149 resolves sanctioned regeneration outputs from the
  tree; ADR-0083 makes adoption one move; ADR-0092 types the Backlog; ADR-0116
  checks a citation against what it cites and is unchanged here. This Spec's gate is bound by ADR-0080,
  ADR-0091, ADR-0096, ADR-0097, ADR-0104, ADR-0155 and ADR-0156. All hold.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization recorded in
  [_authorization.md](_authorization.md); bounded files:
  `internal/speccheck/constraints.go`, `internal/speccheck/coherence.go`,
  `internal/docscontract/testdata/corpus-golden.json`,
  `internal/spec/archive_layout_characterization_test.go`,
  `.agents/skills/qa-gate/SKILL.md`, `skills/qa-gate/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `.agents/skills/write-prd/references/prd-template.md`,
  `skills/write-prd/references/prd-template.md`,
  `.agents/skills/write-techspec/references/techspec-template.md`,
  `skills/write-techspec/references/techspec-template.md`,
  `.agents/skills/write-tasks/references/task-template.md`,
  `.agents/skills/write-tasks/SKILL.md`, `skills/write-tasks/SKILL.md`,
  `skills/baseline_skill_contract_test.go`. Sanctioned regeneration:
  `make skills-sync`. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`.

## Goals

- A Task that declares a Governed Path its Spec does not authorize is refused
  before dispatch, not after the Agent's work.
- A Task that changes the CLI names the guide that describes the change.
- The authoring templates and skills produce Specs the checker accepts.
- Three honesty gaps in the checker close, and one correct check gains the
  tests that keep it correct.

## Core Features

1. **Undeclared Governed Paths are refused at authoring.** `spec check`
   resolves every path a pending non-QA Task declares under `interface:` or
   `creates:` in its Context, and every repository file its Verification
   names, through `GovernedPath`. A Governed Path missing from the Spec's
   authorization record, or from a Tooling authority `bounded files:` row, is
   refused with `SC-TOOLING-UNDECLARED`; so is a row that names a Governed Path
   the record does not grant.
2. **A CLI change names its guide.** A pending Task whose Context names a CLI
   surface, and which neither names a skill or user guide itself nor depends on
   a Task that does, is reported as the gap `SC-CLI-UNDOCUMENTED`.
3. **Honest Verification and declarations.** `SC-VERIFY-INVERTED-EXIT`
   recognises a tool piped into `grep` inside a command substitution whose
   result is tested for emptiness; a wrapped `None.` paragraph is accepted
   while missing, empty and mixed sections are still refused; an indexed
   adopted source whose original is still in `docs/findings/` or
   `docs/backlog/` is refused; the one-field and blank-field closure cases are
   pinned by tests.
4. **Templates and skills produce accepted Specs.** The PRD and TechSpec
   templates use `bounded files:`; the Task template and the write-tasks skill
   state the status-preserving Verification form, which Context entries are
   audited, and that a CLI change names its guide.

## Non-Goals / Out of Scope

- Predicting paths a Task edits without declaring them. An undeclared edit is
  still caught only by the Daemon's changed-path audit at commit; this Spec
  makes declaring every edited path the authored rule and refuses what is
  declared.
- Changing the governed set, the changed-path audit, or `GovernedPath`.
- Accepting `bounded paths:` as a synonym; the templates change instead.

## Success Metrics

1. A pending Task declaring `Makefile` under `interface:` in a Spec whose
   authorization record omits it produces exactly one `SC-TOOLING-UNDECLARED`
   error naming `Makefile`, and none once the record and both rows name it.
2. Replaying archived Spec 0155 with its Tasks pending and its `paths:` emptied
   reports `SC-TOOLING-UNDECLARED` for each of `Makefile`, `.roundfixrc.yml` and
   `.github/workflows/ci-verify.yml`.
3. A pending Task naming `internal/cli/cli_test.go` with no guide in itself or
   its dependencies produces one `SC-CLI-UNDOCUMENTED` gap, promoted to an error
   by `--strict`.
4. The active corpus golden records `0` for both new codes, and this Spec adds
   no finding to `roundfix spec check --strict`.
5. `x="$(go vet ./... | grep pat)"; test -z "$x"` is refused with
   `SC-VERIFY-INVERTED-EXIT`, and a wrapped `None.` paragraph produces no
   `SC-METRIC-UNDECLARED`.
6. Neither the PRD template nor the TechSpec template contains
   `bounded paths:`.

## Recorded limits

- A Verification path is read only when it resolves to an existing repository
  file or to a path the same Task declares under `creates:`; a directory
  operand such as `./internal/docscontract` is not audited.

## Decisions

- **Refuse what is declared, not what is guessed.** ADR-0093 forbids
  inference; declared Context entries and Verification operands are the only
  inputs, and `instruction:` entries, which name guidance the Agent reads, are
  not audited.
- **One reader of Verification paths.** The check reads Verification operands
  through the same function `spec.Collisions` uses, so the wave-collision check
  and this one cannot disagree about which files a Task names.
- **Change the templates, not the checker.** `bounded paths:` stays refused;
  the templates are the defect.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in
the Task Graph. The negative cases carry the weight: a declared Governed Path
that passes without authorization, or a pipe-to-`grep` form still approved,
would pass the happy path. Success Metric 2 rests on evidence this Spec did not
author: the Tasks and grant of archived Spec 0155, whose Runs the Daemon refused
for exactly those three paths.

## Research basis

The seven Backlog Entries adopted under [references/](references/_index.md)
record the defects and their evidence. The Governed Path refusals are recorded
in the archived authorization records of Specs 0155, 0159, 0162, 0163 and 0167;
the efficiency diagnosis of 2026-09-25 counted the sixteen Unresolved Runs.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
