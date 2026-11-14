---
status: accepted
created_at: 2026-11-14T10:00:00Z
updated_at: 2026-11-14T10:00:00Z
deprecated_at: null
superseded_by: null
---

# A Spec is sized by one context and by Tasks that earn their verification cycle

The grouping bound stops being a count. A Spec holds one context, and each of
its Tasks earns the verification cycle it pays for; four implementation Tasks
plus the QA gate stays as a reference, not a cap. The rule says why, so an
author can apply it to a case no number anticipated.

## Context

ADR-0208 wrote down the limit the repository's authoring already worked to:
"Grouping stops where the Spec would outgrow four implementation Tasks plus
its QA gate." The Baseline clause `clause.spec.sources-02-bound-the-group`,
the `write-idea`, `write-prd` and `write-techspec` skills, the Roundfix skill's
`spec` reference, the `spec` command guide and the **Grouping Bound** glossary
entry all repeat the number, and none of them says what it protects.

On 2026-11-14 the maintainer recalled the purpose: "O objetivo da spec com
limite de 4 tasks era porque tinhamos specs gigantes no inicio com varios
contextos na mesma spec e tambem specs com tasks pequenas que consumiam runs
com grandes tempos de verificação (rotina normal) e faziam a spec demorar
muito mais do que o suficiente." Asked how the rule should read, he chose
"Trocar o número pelo critério": "Sai o limite fixo. A Spec passa a ser
medida por um contexto só e por Tasks que justifiquem cada uma o seu ciclo de
verificação; o 4 vira referência, não teto. Muda o ADR-0208 e o termo do
glossário."

A count answers neither failure. A Spec of four Tasks can still hold two
contexts, and a count of four says nothing against four Tasks so small that
each spends an Agent turn, its declared Verification and the repository
Verification at settlement on one sentence. Nothing in the code enforces the
number: no command, check or judge question counts a Spec's Tasks.

## Decision

- **One context.** The sources a Spec adopts and the work it holds change the
  same component or contract, share a root cause, or complete one another, the
  same test ADR-0208 gives for grouping. A scope that holds more than one
  context is split into Specs that each hold one.
- **Each Task earns its verification cycle.** Every Task pays an Agent turn,
  its declared Verification and the repository Verification at settlement. A
  change too small to justify that cost joins a neighbouring Task of the same
  context, inside that Task's slice, instead of becoming its own Task. A small
  source of the same context joins a Spec that is open or being planned instead
  of becoming a one-Task Spec that pays a whole Run, QA gate, review and
  release.
- **Four is a reference, not a cap.** Four implementation Tasks plus the QA
  gate stays as the usual size. A Spec that needs more is a signal to check
  whether it holds more than one context, never a refusal by count.
- **The reasons are part of the rule.** The clause says what the criteria
  answer: a Spec that holds several contexts grows giant, and a Spec of small
  Tasks spends Runs on long verification and takes much longer than its work
  needs.
- **Every other guarantee stays.** Each source keeps exactly one owning Spec;
  a split gives each source one owner; a source that shares no context is
  never added only to save a Spec; a Spec whose Verification or premise depends
  on an unmerged Spec names it in `requires` and in the TechSpec Build Order.
- **The term keeps its name.** **Grouping Bound** now names the criterion. The
  clause keeps its id and its force, so an adopter's next `roundfix baseline
  update` rewrites its words and records it `retained`.

## Considered Options

- **Keep four and state the two reasons.** The recommended option of the
  maintainer's question. Rejected by the maintainer's answer: a count stays a
  count, and the reasons would describe a rule that does not test them.
- **Record the observation as a Backlog Entry and keep the rule.** Rejected by
  the same answer: the rule would keep refusing by count until a later Spec.
- **Enforce a minimum Task size or a maximum Task count in `roundfix spec
  check`.** Rejected: whether a change earns its cycle and whether two sources
  share a context are judgments about meaning, not structure; the strict check
  stays structural, and the judge's grouping question stays advisory
  (ADR-0209).
- **Rename the term.** Rejected: the Grouping Bound is still the limit on
  grouping, and a new name would leave every earlier Archive Record and ADR
  speaking of a term the glossary no longer holds.

## Consequences

- Every adopter of the Spec workflow receives the reworded clause through its
  next `roundfix baseline update`; the clause id, force and Source Baseline row
  are unchanged, and the Source Baseline corpus keeps the wording it recorded.
- An author may plan a fifth implementation Task when the Spec still holds one
  context, and is expected to fold a too-small Task into its neighbour.
- A grouping or a Task the author judged wrong costs a later split or a
  corrective Task, as before; the judge can still suggest a grouping, and the
  author still decides.
