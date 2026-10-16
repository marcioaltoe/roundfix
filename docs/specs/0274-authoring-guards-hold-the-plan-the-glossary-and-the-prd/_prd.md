---
spec: 0274-authoring-guards-hold-the-plan-the-glossary-and-the-prd
status: active
created: 2026-10-16
surfaces: [backend, cli, docs]
---

# Authoring guards hold the plan, the glossary and the PRD to their rules

Four authoring rules had no check, and each cost work. On 2026-10-14 a Task of
Spec 0268 edited three user guides while its authorization's Limits forbade the
Behavior Surface Record; its Verification failed twice, the Run ended
Unresolved, and `roundfix spec check --strict` had been clean on the plan
([the first Backlog Entry](references/2026-10-14-spec-check-does-not-flag-a-guide-task-without-the-coverage-record.md)).
The `roundfix` skill and the `implement` user guide state the QA secret refusal
in a form the Daemon does not print
([the second Backlog Entry](references/2026-10-14-the-implement-reference-states-the-qa-secret-refusal-in-another-form.md)).
Every Spec is held to the glossary, but nothing holds `CONTEXT.md` to itself,
and a congruence audit on 2026-10-14 found sixteen undefined phrases and three
closed lists that had drifted for months
([the third Backlog Entry](references/2026-10-16-the-glossary-is-not-held-to-the-rule-it-applies-to-specs.md)).
And PRDs of unattended cycles are written without any interview by
construction, while nothing in them says who settled their questions. After
this Spec each of the four rules fails a check the moment it is broken.

## Prerequisites

Spec 0272, authored beside this one, changes the docs-layout clauses and the
glossary entries History Root, History Sanitize Command, Reduced History Entry,
Rollup, Backlog Entry and Legacy Link Rebase, and adds Git Locator. This Spec
runs after 0272 merges (queue order B 0273, A 0272, C 0274, D 0275) and touches
none of those clauses or entries. No Verification here depends on an artifact
of another unmerged Spec, so the Task Graph declares no `requires`: the
glossary contract reads whatever `CONTEXT.md` the Run starts from, and the Task
that adds it fixes any gap it finds there, including in entries 0272 wrote.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, numbering or
  record schema is added; two Spec Consistency Check codes join the existing
  code set. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request, credential or network
  call is added; the checks read local files and Git history, and tests run on
  temporary repositories. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0280 (this Spec) decides the three
  guards: "A guide Task records its own fingerprint.", "A PRD records who
  decided." and "The glossary is held to its own rule.". ADR-0256: "A release
  waits for the skills that describe a changed surface", whose record the
  first guard protects. ADR-0244: "A Spec declares the domain terms it
  introduces", whose rule the glossary now meets itself, and whose horizon
  the second guard copies: "The section is required only for a Spec whose PRD
  was committed at or after". ADR-0273: "A Secret Match reports
  the repository-relative path, the line and the kind", the form the docs now
  state.
  ADR-0091: "The QA gate is a Task node of its own type", whose `qa` value the
  Task Type mapping holds. ADR-0217: "A Cursor selection names the advertised
  model value"; the ACP Runtime list it changed is not mapped, because no
  exported list of runtimes exists. ADR-0182: "A Task settles on the facts its
  gate will check", which is why the guide Task itself must record its
  fingerprint. ADR-0117: "A defect is checked by the stage that can produce
  it", so the Clarification Record is checked from the PRD stage and the
  surface record from the Tasks stage. ADR-0104: "A Spec accepts on evidence
  it did not author", which the Acceptance evidence row meets. ADR-0096: "The
  QA gate proves machine facts before it spends an agent turn"; the gate's
  mechanical stage is unchanged. ADR-0167: "The pre-PR Pull Request row never
  decides a qualifying partial" and ADR-0240: "One QA partial policy, and rows
  a Run sandbox cannot reach" govern the gate's blocked rows. ADR-0276: "A row
  the authorization keeps offline qualifies, and a byte invariant compares
  with the Delivery Base", which the QA corpus row follows. ADR-0229: "An
  operator archive resumes any park from the Run start" and ADR-0237: "A
  Delivery Retry records a merge made outside the queue" set no obligation
  here: they cite ADR-0091 for the gate's position and this Spec changes
  neither the archive nor the Delivery Queue. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AH approval of 2026-10-16 through AskUserQuestion, option "Governed Paths
  declarados" ("\"Concedo\" como em AE/AF/AG para os Governed Paths que cada
  Spec declarar: módulos da Baseline, registros de versão, make
  baseline-digests, skills e testes governados. Makefile, workflows de CI,
  go.mod e configuração de lint ficam fora."), and the 2026-09-30 skills grant
  ("considere autorizado a ajustar todas as skills se necessário"). Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0274-authoring-guards-hold-the-plan-the-glossary-and-the-prd/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/implement.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/write-prd/SKILL.md`,
  `.agents/skills/write-prd/references/clarification-record.md`,
  `.agents/skills/write-prd/references/prd-template.md`,
  `.agents/skills/write-tasks/SKILL.md`,
  `.agents/skills/write-techspec/SKILL.md`,
  `internal/docscontract/testdata/corpus-golden.json`,
  `internal/spec/archive_layout_characterization_test.go`,
  `internal/speccheck/coherence.go`, `internal/speccheck/constraints.go`,
  `skills/roundfix/SKILL.md`, `skills/write-prd/SKILL.md`,
  `skills/write-prd/references/prd-template.md`,
  `skills/write-tasks/SKILL.md`, `skills/write-techspec/SKILL.md`. No
  Makefile, workflow, Go module, lint, formatter or Baseline module change.

## Goals

- A plan whose guide Task cannot re-record the Behavior Surface Record fails
  `roundfix spec check --strict` before a Run starts.
- Every document that states the QA secret refusal states it as the Daemon
  prints it, and a contract test fails when either side moves alone.
- `CONTEXT.md` fails a contract test when it uses an undefined capitalized
  phrase, links a file that is not an accepted ADR, lets a closed list drift
  from its constants, writes prose in an `_Avoid_` line or repeats a rule.
- A new PRD says who settled its questions: a person, with dates, or nobody,
  with the source the author decided from.

## User Stories

1. As a Spec author, I want `spec check` to name a guide Task that does not
   re-record the Behavior Surface Record, so that the Run does not end
   Unresolved on a fingerprint the Agent may not regenerate.
2. As an operator reading the `roundfix` skill or the `implement` guide, I want
   the QA secret refusal written exactly as the Daemon prints it, so that I can
   match the message in a Run's events.
3. As a maintainer, I want the glossary held to the rule it applies to Specs,
   so that drift is caught by the change that causes it instead of a later
   audit.
4. As a reviewer of a PRD, I want to read whether a person answered its
   questions, so that I know whether its decisions were interviewed or
   decided by the author from a recorded source.

## Core Features

1. Surface record gap. The Spec Consistency Check reports, as the gap
   `SC-SURFACE-RECORD-UNTASKED`, each pending non-QA Task that declares a user
   guide page without declaring the Behavior Surface Record and naming its
   `-record-skill-coverage` step itself, and each such Task of a Spec whose
   authorization Limits name the record. A repository without the record
   skips the detector.
2. QA secret refusal in step. The `roundfix` skill's `implement` reference and
   the `implement` user guide state the refusal as
   `QA evidence carries a secret: <path>:<line>: <kind>`, and a Repository
   Contract Test derives that form from the Daemon's message and the Secret
   Match format and fails when a document states another.
3. Glossary contract. A Repository Contract Test holds `CONTEXT.md` to five
   rules: capitalized phrases in definitions are defined terms, plurals or
   allowlisted names; every link is an accepted ADR; mapped closed lists equal
   their constants; `_Avoid_` lines are plain lists; and no sentence longer
   than sixty characters appears in two entries. It starts green on the tree
   it is added to, and the allowlist and mapping live in one place.
4. Clarification Record. A PRD's `## Decisions` section opens its
   **Clarification Record**: one dated line for each question a person
   answered, or the autonomous-route line naming the source of its decisions.
   The check reports a PRD with neither as the gap `SC-DECISIONS-UNRECORDED`,
   from the PRD stage, for PRDs committed at or after the guide that defines
   it. The `write-prd` and `write-techspec` skills write it.
5. Documentation. The `spec` and `implement` user guides, the Roundfix skill's
   `spec` and `implement` references, the `write-tasks` skill and the glossary
   describe the guards.

## Non-Goals / Out of Scope

- Inferring a command-help change from the files a Task declares; the Behavior
  Surface Record's own check keeps that case.
- Holding a Task that edits an owned skill to its version record in
  `spec check`.
- Judging whether a PRD's interview really happened, or requiring `grilling`
  before a PRD.
- Applying the Clarification Record to PRDs committed before its guide, or
  rewriting any existing PRD.
- Holding an adopter's glossary to Roundfix's constant mapping.
- Changing the Daemon's refusal text, the Secret Match format, the Baseline
  modules, the `grilling` or `grill-with-docs` skills, or anything 0272 owns.

## Success Metrics

1. Success Metric: a Spec Consistency Check of the 0268 plan as committed in
   `e1e376dd`, run with the built binary, reports `SC-SURFACE-RECORD-UNTASKED`
   for its task_01 and fails `--strict`; the same plan with the fix the
   operator applied on 2026-10-14 (the record and its step declared in
   task_01, the record dropped from Limits) reports none.
2. Success Metric: reverting either the skill reference or the user guide to
   the parenthesized form makes the QA secret contract fail, and the delivered
   tree passes.
3. Success Metric: a planted violation of each of the five glossary rules makes
   the glossary contract fail, and the delivered `CONTEXT.md` passes all five.
4. Success Metric: a PRD committed after the horizon without the record yields
   `SC-DECISIONS-UNRECORDED` at `--stage prd`; one with a dated answer or the
   autonomous-route line yields none; a PRD committed before the horizon is
   skipped.
5. Success Metric: every active Spec's count of each new code at the audited
   head equals the corpus golden, and no existing count moves.

## Acceptance evidence

The outside-evidence row rests on published sources this Spec did not write,
read through Exa on 2026-10-16:

- The Go Blog, "Testable Examples in Go" (https://go.dev/blog/examples):
  "Having executable documentation for a package guarantees that the
  information will not go out of date as the API changes." The QA secret
  contract applies that idea to a refusal message: the documented form is
  computed from the producer, not copied.
- Vale documentation, "Vocabularies" (https://docs.vale.sh/keys/vocabularies):
  "Entries in `reject.txt` are automatically added to an existence rule
  (`Vale.Avoid`) that will flag all occurrences as errors." A published prose
  linter enforces an accepted-term list mechanically, which is the shape the
  glossary contract gives `CONTEXT.md`.
- MADR template (https://github.com/adr/madr/blob/main/template/adr-template.md):
  "decision-makers: {list everyone involved in the decision}". A published
  decision-record format records who decided, which is what the Clarification
  Record adds to a PRD.

## Validation

The decisions were checked on 2026-10-16 against two sources, and the change
each prompted is named.

- Secondbrain. `wiki/index.md`, then three `qmd query` runs from the
  Secondbrain root: "glossary held to its own rule contract test undefined
  terms closed lists" returned the mirror of ADR-0244
  (`projects/roundfix/mirror/docs/adr/0244-a-spec-declares-the-domain-terms-it-introduces-and-the-check-holds-them-to-the-glossary.md`)
  and the `write-prd` glossary reference mirror; ADR-0244's horizon reasoning
  ("no Spec authored earlier is blocked by a rule it could not know") is why
  the Clarification Record has a horizon. "PRD grilling decisions record
  autonomous route no person present" returned the `write-prd` skill mirror
  and the finding mirror
  `projects/roundfix/mirror/docs/history/findings/2026-08-04-what-still-needs-a-supervisor-between-a-prd-and-a-merge.md`;
  neither defines a record of who answered. "spec check guide task behavior
  surface record coverage fingerprint stale" returned the mirror of ADR-0256,
  the adopted Backlog Entry and the mirror of the ADR that holds structural
  proof gaps without a horizon, which confirmed that the record's
  check runs at every settlement and that structural proof gaps hold every
  active Spec. The inbox entry
  `inbox/roundfix/_triaged/2026-10-14-hold-the-glossary-to-the-rule-adr-0244-applies-to-specs.md`
  gave the five rules and the audit's numbers.
- Exa. Supporting: the Go Blog on testable examples, Vale's vocabularies, and
  the MADR template above. Challenging: a terminology-linting guide
  (https://weesholapara.com/blog/terminology-consistency-in-documentation-one-term-per-concept-enforced-with-vale)
  warns that a generic word "may be a rejected synonym in one sentence and a
  perfectly legitimate word in the next", and that a noisy rule teaches people
  to ignore it. That is why the phrase rule reads runs of capitalized words
  rather than single words, splits a run into multi-word terms, keeps
  `_Avoid_` items out of it, and starts green with an allowlist whose every
  entry carries a reason. These are standard references, not new research, so
  no Inbox Entry was captured.

## Glossary

- adds: **Clarification Record**

## Decisions

- 2026-10-14: the maintainer noted that "a skill `grilling` não é usada
  corretamente na criação de PRDs"; the minimal contract is a record of the
  route, not a mandatory interview. See ADR-0280.
- 2026-10-16: the maintainer approved the scope of item C through
  AskUserQuestion: "spec check acusa Task de guia sem o record de cobertura; a
  mensagem de segredo da QA no reference igual à do Daemon, com contrato; o
  glossário cobrado de si mesmo (teste de contrato do inbox); registro do
  grilling no PRD."
- 2026-10-16: the maintainer approved the queue order "B → A → C → D
  (Recomendado)" and the grant "Governed Paths declarados".

Taken by the author from those answers and the adopted Backlog Entries, with
no further interview:

- The guide Task itself records the fingerprint; another Task cannot cover it.
  See ADR-0280.
- Command help changes stay with the record's own check. See ADR-0280.
- The Clarification Record has a horizon, unlike the proof gaps. See ADR-0280.
- All five glossary rules fit one Task, so none is dropped; `_Avoid_` items may
  name defined terms. See ADR-0280.

## Open Questions

None.
