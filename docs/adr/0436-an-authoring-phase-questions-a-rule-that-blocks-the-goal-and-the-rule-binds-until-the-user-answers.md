---
status: accepted
knowledge: true
created_at: 2026-11-19T11:00:00Z
updated_at: 2026-11-19T11:00:00Z
deprecated_at: null
superseded_by: null
---

# An Authoring Phase questions a rule that blocks the goal, and the rule binds until the user answers

In brainstorming, the Idea and the PRD, an author that finds a recorded
decision in the way of the goal neither obeys it in silence nor overrides it
in silence: it puts a Rule Challenge to the user, who keeps, changes or
retires the rule. The rule binds until the user answers, the artifact records
the question and the answer, and the Spec that asked carries the change. The
TechSpec author asks the same question before a design takes what an accepted
ADR rejected.

## Context

On 2026-11-18 the maintainer asked for it in these words: "Tenho visto com
frequência esse problema na construção de soluções, porque tem uma regra,
então, nem é proposto a mudança da regra, deprecar, regra, simplesmente não
faz, porque existe uma regra. Em algum momento, talvez depois de uma, duas,
três specs, é falado, não, não foi feito, porque existe essa regra, sim, lá do
início do projeto." He named the cause: "As ADRs, regras do specific
repository e outras referências que possam existir dentro do projeto, elas
são para guiar a implementação, são guard-rails de guiação da implementação.
No momento que a gente está trabalhando no IDEA e no PRD, a gente está
trabalhando com criatividade." And the remedy: "que você traga essa
possibilidade de questionamento usando as ferramentas de Ask a question ou
compatíveis no Codex ou no Open Code, para que o usuário decida manter ou não
com base na análise que vai ser feita. [...] Isso tem que existir também no
baseline."

The instructions an author read until then explain the silence. The
Instruction Hierarchy says a narrower guide "cannot weaken a universal
Normative Clause or confirmed project decision" and makes no difference
between the stage that shapes a solution and the stage that builds it. The
domain guide forbids one silence, "flag conflicts instead of silently
overriding repository decisions", and no clause forbade the other one, which
is giving up a goal without saying so. `write-idea`, `write-prd` and
`brainstorming` received a rule only as an obligation to list, and
`write-techspec` could contest an ADR without asking the user.

A count of the plans of Specs 0364 to 0406, made on 2026-11-18 from the Git
history, found that 15 of 41 PRDs draw a scope exclusion at an ADR or a guide,
that 12 amend or supersede an ADR, and that none of the 34 counted lines
records a question to the maintainer about whether a recorded decision is
kept, changed or retired. The count cannot tell a boundary the author chose
from a goal the author gave up, because no file says which it was.

Published evidence explains why an Agent does not ask on its own, and what a
question costs. Xu and Wu, "Stale Instruction Files Override the Code in
Coding Agents: A One-Sentence Fix" (https://arxiv.org/abs/2610.30571),
measured 320 runs over 40 projects: "Claude Code pointed out the conflict in
79 of 80 runs and still followed the file", and one sentence that gave the
code precedence "raised the share of choices that follow the code from 4.2%
to 91.4%". Zhang and others, "Ask Less, Ask Better"
(https://arxiv.org/abs/2610.24411), found that "limiting agents to one to
three questions reduces user-attention cost by 36–73% while keeping 86–94% of
the task reward they earn with unlimited questions". The Spec Kit
constitution takes the opposite position to this decision: "Violations are
resolved by changing the spec, plan, or tasks — not by diluting a principle".

The maintainer settled the shape in the Idea stage on 2026-11-18, each answer
the option he chose: the size of the first version, "Cláusula, passo e
registro (Recommended)"; the TechSpec, "Entra, só ao contradizer ADR
(Recommended)"; who applies the answer, "A própria Spec (Recommended)";
nobody present, "A Spec espera (Recommended)"; the reach, "Sim, vira proposta
ao Roundfix (Recommended)"; the answers, "Três: manter, mudar, retirar
(Recommended)"; and a proposal to the Baseline's owner that has no answer
yet, "Segue dentro da regra (Recommended)".

## Decision

**Both silences are forbidden.** In brainstorming, the Idea and the PRD, an
author never drops or narrows a solution in silence because a rule forbids
it, and never contradicts a rule in silence. It puts a Rule Challenge to the
user.

**A rule is a decision the repository recorded.** An accepted ADR, a
Normative Clause of a guide under `docs/agents/`, a definition in the domain
context, or a document under `docs/references/` that the repository treats as
binding. A user-level instruction file, an ai-memory page, a harness's own
rules and a default of a Vendored Skill are not rules of the repository.

**The trigger is a blocked goal.** A rule blocks the goal when the solution
the author would recommend, or a part the goal needs, cannot exist without
contradicting the rule, or when the rule stops the author from stating the
problem. Extra work, taste and friction are not a block. No Rule Challenge is
asked about a change the user already decided, because the request asks for
it or it follows from an answer of the user that the artifact records; the
artifact names that request or answer.

**The question has three answers.** One question per rule, through the
question contract of the Baseline. It names the rule and where it lives, what
the rule imposes, up to three things it prevents for this goal, and the
author's reading of why the rule exists and whether it looks out of date or
badly written. The answers are keep, change and retire, each with its
consequence and the recommended one first. An Authoring Phase or a
brainstorming session asks at most three Rule Challenges; when more rules
block, it asks about the three that prevent the most and lists the others,
which keep binding, in the artifact's open questions.

**The rule binds until the user answers.** Asking changes nothing. A
recommendation, silence, a dismissal or a timeout is never an answer, and the
rule's text stays as it is until a Task of the Spec changes it.

**The artifact records it.** The Idea and the PRD carry a `## Rule
challenges` section with one fielded record per question: the rule, what it
imposes, what it prevents, when and in which phase it was asked, the answer
in the user's words, and what the Spec does with it. The section reads
`None.` when no rule blocked the goal. The PRD carries the Idea's records,
and a question the TechSpec author asks is recorded in the PRD.

**The Spec that asked carries the change.** On change or retire, the PRD
states the change of the rule as a requirement, and a Task of the same Spec
changes the rule together with the work that depended on it. An ADR changes
through the lifecycle the docs-layout guide defines; this decision adds no
status. On keep, the rule gains a Kept Line
([ADR-0437](0437-a-rule-the-user-keeps-carries-a-kept-line-and-is-asked-about-again-only-with-a-new-fact.md)).

**A Spec waits when nobody can answer.** With no person present, no artifact
of that phase is written and nobody answers for the user. Only that Spec
waits, as a Backlog Entry that holds the whole question, names the Spec
folder when one exists and carries a Waiting Declaration; other work goes on.
An Author Subagent returns the question to the session that started it.

**A rule the Baseline owns becomes a proposal.** A clause inside a Managed
Region changes only in the Baseline. On change or retire, a repository that
does not own the Baseline sends the user's decision as one ai-memory message
to the Baseline's owner, obeys the rule until the Baseline changes, and
records the message and the scope it left out.

**The TechSpec asks before it contradicts an ADR.** When a design takes what
an accepted ADR rejected, or otherwise contradicts its decision, and the user
has not decided it, the TechSpec author asks the same question before it
supersedes or amends the ADR.

**Implementation keeps its guard-rails.** In a Run, in a Task and in the QA
gate a rule is a guard-rail: an Agent flags a conflict there, as the domain
guide already says, and asks no Rule Challenge.

**The Baseline carries the rule, and nothing checks it mechanically.** Two
Normative Clauses of the Spec workflow state it for every adopter, so it also
governs a session that authors through a Vendored Skill. The owned skills
`brainstorming`, `write-idea`, `write-prd` and `write-techspec` carry the
step, and the Idea and PRD templates the section. `roundfix spec check` reads
none of it.

## Considered Options

- **Keep listing a rule as an obligation.** Rejected: it is the silence the
  maintainer reported.
- **Let the author override a rule and record why.** Rejected: the user who
  confirmed a decision is the one who changes it.
- **A clause alone, with the answer in free text.** Rejected by the
  maintainer's choice of the record: without it nobody can tell "the author
  looked and found nothing" from "the author did not look".
- **A mechanical check or a judgment of exclusions.** Left out of the first
  version by the maintainer's choice. A structural check cannot tell a
  boundary from a goal given up: 19 exclusions in 41 PRDs name a rule, and
  most read as boundaries.
- **Change the rule at once, outside the Spec.** Rejected: with the Spec, the
  rule and the work that depended on it land and are verified together.
- **Go on inside the rule when nobody is present.** Rejected by the
  maintainer: the Spec waits, so nothing is built that the user did not
  choose.
- **Four answers, with deprecate and delete apart.** Rejected: the question
  contract takes two or three options, and under the retirement clause both
  end the same way.
- **Questions during implementation.** Rejected: there a rule is a
  guard-rail.

## Consequences

- Every adopter of the Spec workflow receives the two clauses with its next
  `roundfix baseline update`.
- An Idea and a feature PRD say whether a rule blocked the goal, so a later
  reader can tell a chosen boundary from a goal given up.
- The trigger is the author's judgment. An author that wants its own solution
  can read friction as a block; the cap, the binding rule and the Kept Line
  bound the cost, and the record shows the load.
- A Spec can stop on an Agent's judgment when nobody is there to answer.
- The section may fill with `None.` by habit, because nothing checks it. A
  check and a judgment of exclusions are later work that needs these records
  to read.
- The Baseline's owner receives proposals from adopters and triages each as
  any other message.
