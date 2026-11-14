---
spec: 0391-the-grouping-bound-becomes-a-criterion
status: active
created: 2026-11-14
surfaces: [docs]
type: fix
---

# The grouping bound becomes a criterion

The Baseline tells every adopter to group sources into one Spec "only while
the Spec fits four implementation Tasks plus its QA gate", and the authoring
skills, the Roundfix skill's `spec` reference, the `spec` command guide and
the **Grouping Bound** glossary entry repeat the number without saying what it
protects. On 2026-11-14 the maintainer recalled its purpose, giant Specs that
held several contexts and Specs of small Tasks that spent Runs on long
verification, and chose to replace the number with the criterion. After this
Spec a Spec is sized by one context and by Tasks that each earn their
verification cycle; four implementation Tasks plus the QA gate stays as a
reference, not a cap, and every guarantee the old clause gave stays. Every
adopter receives the reworded clause through its next `roundfix baseline
update`, and the `write-tasks` skill gains the same criterion for Task sizing.

## Sources

One maintainer request, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-14-maintainer-the-grouping-bound-becomes-a-criterion.md`](references/2026-11-14-maintainer-the-grouping-bound-becomes-a-criterion.md),
   the maintainer's observation and his answers of 2026-11-14, about
   06:45–06:58 BRT, through AskUserQuestion in `/roundfix --follow-up`,
   recorded verbatim by the operator; the maintainer, 2026-11-14. No
   ai-memory message is adopted, and nobody asked for a completion notice.

At the authoring base `docs/backlog/` holds the Conexus status report and the
Homebrew tap entries and `docs/findings/` holds none; none of them touches how
a Spec is sized, so no other source shares this context.

## Prerequisites

None in `requires`. Authored against `251ed81e` (v0.143.3). Specs 0392 and
0393 are authored beside this one against the same main and deliver after it,
in the order 0391 → 0392 → 0393. Files this Spec shares with them, each edited
only in this Spec's own sentences:

- `internal/baseline/assets/modules/spec-workflow.json`: this Spec owns the
  grouping clause `clause.spec.sources-02-bound-the-group`; 0393 owns the
  route clauses of the same rule `rule.spec.routing`. Both raise that rule's
  `version` by one from the value on the tree the Task runs on, and the
  module's own version comes from the record command, so neither writes a
  number.
- `internal/baseline/module-versions.json`, the Standard TypeScript profile,
  its formatter golden of `docs/agents/spec-routing.md`, the catalog files and
  plan goldens under `internal/baseline/testdata/`, and this repository's
  `docs/agents/spec-routing.md` and `docs/agents/setup-context.json`: task_01
  regenerates them from the tree it runs on, through the record command,
  `make baseline-digests` and the managed refresh.
- The `write-idea`, `write-prd` and `write-techspec` skills: this Spec owns
  their `## Sources that share a context` grouping sentence; 0392 adds
  research-subagent text and 0393 route text elsewhere in them. The
  `write-tasks` skill: this Spec owns its Task sizing guidance.
- The Roundfix skill's `references/spec.md` and `docs/user-guide/commands/spec.md`:
  this Spec owns only the sentence that states the grouping bound; 0393 owns
  the judge question and its command surface.
- The version fields of each changed skill's `SKILL.md` pair and
  `skills/testdata/owned-skill-versions.json` (all three Specs), `CONTEXT.md`
  (this Spec changes only **Grouping Bound**) and
  `docs/references/behavior-surfaces.json`: task_02 raises the versions and
  re-records the fingerprint with the repository's commands on the tree it
  runs on and never writes a number.

## Project Constraints

- Identifier strategy: not applicable — no identifier, record or schema is
  minted; the clause keeps its id. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network
  call is added; the change is Baseline guidance, skill text, the glossary and
  an ADR. Source: `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0400 (new) records the decision: "A
  Spec holds one context, and each of its Tasks earns the verification cycle
  it pays for; four implementation Tasks plus the QA gate stays as a
  reference, not a cap." ADR-0208 is amended: its bullet "**The bound.**
  Grouping stops where the Spec would outgrow four implementation Tasks plus
  its QA gate" no longer holds, while its grouping and extend-before-minting
  clauses and "each source keeps exactly one owning Spec" stay. ADR-0209
  applies unchanged: "A suggested pair is answered by adopting the open source
  or by stating why it stays apart. It never gates"; the judge's grouping
  question is not changed. ADR-0200 and ADR-0201, which ADR-0209 cites, do not
  apply: this Spec changes neither the judge's force nor what it sends.
  ADR-0319, which cites ADR-0209, does not apply: it orders the Follow-up
  Command's items, and this Spec changes no follow-up order. Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  answers of 2026-11-14 quoted in the authorization record, on the cycle BE
  terms, which keep the grants of the Baseline source with its derived files,
  the guides under `docs/agents/` and the owned skills; the Governed Paths are
  the spec-workflow module, the Standard TypeScript profile and its formatter
  golden, this repository's spec routing guide and setup manifest that task_01
  changes through their commands, and the skill files task_02 changes. No
  Makefile, workflow, Go module, lint or formatter change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0391-the-grouping-bound-becomes-a-criterion/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/spec.md`, `.agents/skills/write-idea/SKILL.md`, `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-tasks/SKILL.md`, `.agents/skills/write-techspec/SKILL.md`, `docs/agents/setup-context.json`, `docs/agents/spec-routing.md`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`, `internal/baseline/assets/modules/spec-workflow.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `skills/roundfix/SKILL.md`, `skills/write-idea/SKILL.md`, `skills/write-prd/SKILL.md`, `skills/write-tasks/SKILL.md`, `skills/write-techspec/SKILL.md`.

## Goals

- A Spec is sized by one context and by Tasks that each earn their
  verification cycle, never refused by a count of its Tasks.
- A change too small for its own verification cycle joins a neighbouring Task
  of the same context, and a small source of the same context joins a Spec
  that is open or being planned instead of becoming a one-Task Spec.
- The rule states its two reasons, so an author can apply it to a case the
  reference number does not cover.
- Every guarantee of the old clause stays: each source keeps exactly one
  owning Spec, a split gives each source one owner, a source that shares no
  context is never added only to save a Spec, and `requires` names an unmerged
  Spec that a Verification or premise depends on.
- Every adopter receives the reworded clause through its next
  `roundfix baseline update`, recorded `retained`.

## Core Features

1. The Baseline clause `clause.spec.sources-02-bound-the-group` keeps its id
   and force and states the two criteria, the two reasons, the reference and
   every guarantee it carried; its rendered guide and derived files follow,
   and each kept guarantee and each new criterion is pinned by its own phrase.
2. The `write-idea`, `write-prd` and `write-techspec` skills replace "Group
   only while the scope fits four implementation Tasks plus its QA gate" with
   the criterion; the `write-tasks` skill gains an "Earns its verification
   cycle" rule beside "Sized for one fresh session"; the Roundfix skill's
   `spec` reference and the `spec` command guide stop calling the bound a
   count. Each phrase is pinned in both skill copies.
3. The glossary keeps the term **Grouping Bound** and changes its definition
   to the criterion; ADR-0400 records the decision and ADR-0208 carries its
   amendment notice.

## Non-Goals / Out of Scope

- A check, judge question or command that counts Tasks or measures a Task's
  size. ADR-0400: "whether a change earns its cycle and whether two sources
  share a context are judgments about meaning, not structure".
- The judge's `source-grouping` question and its threshold (ADR-0209), and
  every route clause of `rule.spec.routing`, which Spec 0393 owns.
- The Source Baseline corpus `baseline.standard-typescript-monorepo-0.0.1`: a
  Source Baseline is "the immutable Setup Manifest and governed instruction
  corpus" (`CONTEXT.md`), and a reworded clause keeps its identity, so its row
  stays as recorded.
- The pre-PR review's Review Segments and Data Review Segments, whose "at most
  four" is a different bound, the corrective-Task ceiling of `write-tasks`, and
  every Archive Record, `CHANGELOG.md` entry and dated reference under
  `docs/references/` that quotes the old number.
- Any Spec already planned or delivered: the rule governs authoring from the
  release on.

## Success Metrics

1. Success Metric: a managed refresh of an adopted temporary repository writes
   into `docs/agents/spec-routing.md`, once each, every phrase that
   `_techspec.md` → Testing Approach lists for the clause, and no longer
   "only while the Spec fits four implementation Tasks plus its QA gate"; a
   second refresh reports it current, and the clause keeps its id, force and
   Source Baseline row, recorded `retained`. At `251ed81e` none of the new
   phrases is in the guide.
2. Success Metric: in both copies of `write-idea`, `write-prd`,
   `write-techspec`, `write-tasks` and the Roundfix `spec` reference, each
   phrase `_techspec.md` → Testing Approach lists appears once, and none still
   bounds grouping by a count; the `spec` command guide says the same.
3. Success Metric: `CONTEXT.md` defines the changed **Grouping Bound**,
   ADR-0208 carries one "Amended by ADR-0400" notice dated 2026-11-14, and the
   documentation, glossary and skill contracts pass.
4. Success Metric: each sabotage `_techspec.md` → Testing Approach lists makes
   its named test fail.

## Acceptance evidence

The outside-evidence rows rest on published literature this Spec did not
write, read through Exa on 2026-11-14; the QA gate compares the quoted
passages with the rule and looks nothing up.

- Donald G. Reinertsen, *The Principles of Product Development Flow* (2009),
  Chapter 1, published by the author
  (http://lpd2.com/wp-content/uploads/2013/06/ReinertsenFLOWChap1.pdf):
  "batch sizes have their own economic tradeoffs, which lead to U-curve
  optimization problems." and "Large batches seem attractive because they
  appear to have scale economies that increase efficiency." The book's
  principle list (https://lpd2.com/sample-page/the-principles-of-flow/):
  "B11: The Principle of Batch Size Economics: Economic batch size is a
  U-curve optimization. (p.121)". A Spec and a Task are batches: each pays a
  fixed transaction cost (a Run, its QA gate, review and release; an Agent
  turn and two Verifications), and a large one holds work back. An optimum
  set by two opposing costs is a criterion, not a fixed count, and both
  failures the maintainer named are the two arms of the U-curve.
- Donald G. Reinertsen, "The science of batch size", FlowCon France 2019
  (https://flowcon.io/archives/2019/the-science-of-batch-size-donald-g-reinerstsen-at-flowcon-france-2019-e12/):
  "the correct answer is never on one extreme or the other" and "when you're
  near the bottom of the U-curve, you don't need a very exact answer." This is
  why four implementation Tasks plus the QA gate stays as a reference, not a
  cap, and why the rule refuses neither a fifth Task nor a joined small one by
  count.

## Unreachable Acceptance

None. Every criterion is proved with this repository's files, temporary
copies of the audited head and temporary adopted repositories; no criterion
needs the network, an open Pull Request or another repository.

## Validation

The decisions were checked on 2026-11-14 against these sources.

- Characterization at `251ed81e`: `rule.spec.routing` is at version 10 in the
  `spec-workflow` module at version 24; `clause.spec.sources-02-bound-the-group`
  begins "Group sources into one Spec only while the Spec fits four
  implementation Tasks plus its QA gate." and renders verbatim at line 19 of
  this repository's `docs/agents/spec-routing.md` and of the Standard
  TypeScript formatter golden of that guide.
  `internal/baseline/grouped_sources_clauses_test.go` pins the whole text as a
  literal. The same number is written in the `write-idea`, `write-prd` and
  `write-techspec` skills ("Group only while the scope fits four
  implementation Tasks plus its QA gate"), in the Roundfix `spec` reference
  ("The grouping bound is four implementation Tasks plus its QA gate."), in
  `docs/user-guide/commands/spec.md` ("within the bound of four implementation
  Tasks plus its QA gate"), in `CONTEXT.md` **Grouping Bound** and in ADR-0208.
  No Go code, strict-check rule or judge question counts a Spec's Tasks
  (`internal/spec`, `internal/speccheck`, `internal/judge/questions.json`).
- A prototype of tasks 01 and 02 on the clone at `251ed81e`: the clause edit,
  the module record command, `make baseline-digests` and the managed refresh
  rewrote exactly the derived files task_01 declares, a second refresh
  reported `current`, and `./internal/baseline` and the CLI's Baseline tests
  passed with the two clause tests; the skill edits, `make skills-sync`, the
  version record and the fingerprint re-record passed `./skills` and
  `-tags docscontract ./internal/docscontract`, glossary contracts included;
  a second `make baseline-digests` changed nothing. Every phrase of both new
  tests made its test fail when removed.
- Secondbrain, read only: `wiki/index.md`, then `qmd query` for Spec size,
  Task limits and verification cost, which returned nothing at 0.3 or more; a
  search of `wiki/` found
  `wiki/sources/branas-formacao-ia-processo-de-desenvolvimento.md`, a course
  on spec-driven development whose author proposed about eight Tasks, joined
  them into three to five at the user's request, "recomenda ~3–4 tasks" and
  notes that "muitas tasks dentro do mesmo PRD têm dependência lógica". It
  supports four as a reference and the joining of small Tasks; it gives no
  reason to keep a cap.
- Exa: the Reinertsen sources above. The FlowCon talk also says that teams
  usually overweight transaction cost and should try smaller batches; this
  Spec does not take that as a reason to push Tasks smaller, because here the
  transaction cost, a full verification cycle, is measured and large.

## Glossary

- changes: **Grouping Bound**

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle BF brief for Spec 0391 of 2026-11-14, the maintainer's answers below and the characterization above.

The maintainer's observation, added as free text to his answer about the
scope of Spec 0392 (AskUserQuestion via `/roundfix --follow-up`, 2026-11-14
~06:45–06:58 BRT), verbatim:

- "E também uma observação sobre a regra de 4 tasks que talvez tenha que
  virar uma nova redação: O objetivo da spec com limite de 4 tasks era porque
  tinhamos specs gigantes no inicio com varios contextos na mesma spec e
  tambem specs com tasks pequenas que consumiam runs com grandes tempos de
  verificação (rotina normal) e faziam a spec demorar muito mais do que o
  suficiente."

The rule question and answer, verbatim:

- Question: "Sobre a sua observação: hoje a regra diz "agrupe fontes numa Spec
  só enquanto ela couber em quatro Tasks de implementação mais o QA" (Baseline
  `spec-routing`, skills write-prd e write-techspec, termo Grouping Bound do
  CONTEXT.md, ADR-0208) e não diz o motivo. Você lembrou os dois motivos:
  evitar Spec gigante com vários contextos e evitar Task pequena que gasta um
  ciclo inteiro de verificação. Como a regra deve ficar?"
- Answer: "Trocar o número pelo critério" — option text: "Sai o limite fixo. A
  Spec passa a ser medida por um contexto só e por Tasks que justifiquem cada
  uma o seu ciclo de verificação; o 4 vira referência, não teto. Muda o
  ADR-0208 e o termo do glossário. Vira Spec neste ciclo."

The authorization question and answer, verbatim:

- Question: "Se a nova redação da regra das 4 Tasks virar Spec neste ciclo
  (primeira pergunta), qual autorização vale para ela? Ela muda um módulo do
  Baseline e duas skills de autoria, que os adotantes recebem no próximo
  `baseline update`."
- Answer: "Mesmas do BE (Recommended)" — option text: "Push, PR e merge;
  release patch depois dela; até 3 correções; Makefile, CI, go.mod e lint só
  quando uma correção exigir, descrito no relatório."

Taken from those answers and the characterization, with no further interview:

- The clause keeps its id and force and is reworded, so the Source Baseline
  records it `retained` and its corpus row stays as recorded.
- The term keeps its name, **Grouping Bound**, with a new definition: it is
  still the limit on grouping, and a new name would orphan every earlier
  Archive Record and ADR that uses it.
- The criteria stay prose an author applies; no check counts Tasks, because
  both criteria are judgments about meaning.
- The reasons are written in general terms an adopter can apply ("a Spec that
  holds several contexts grows giant, and a Spec of small Tasks spends Runs on
  long verification"), not as this repository's history.
- `type: fix`: the Spec corrects the wording of an existing rule to the purpose
  the maintainer states and adds no command, flag or binary behavior; the
  Release Ceiling is `patch`. Applied to its own decomposition, the criterion
  gives two implementation Tasks in series, one per carrier family (the
  Baseline clause with its derived files; the skills, the guide, the glossary
  and the ADR notice), and the QA gate: a separate Task for the glossary and
  the ADR notice would not earn its verification cycle.

## Open Questions

None.
