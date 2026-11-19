---
status: accepted
created_at: 2026-11-19T16:00:00Z
updated_at: 2026-11-19T16:00:00Z
deprecated_at: null
superseded_by: null
---

# An open authoring question waits unless the project declares that the session follows the author

When the author of a PRD, a TechSpec or a Task Graph ends its phase with a
question that no source settles and nobody is present, the Spec waits. A
project may declare in its Project Config that the session follows the option
the author recommends instead. The adopted option is a Delegated Decision: it
is recorded in the PRD in a form of its own, it is never an answer of the
user, and the user's later answer wins.

## Context

An Author Subagent that lacks a decision returns its questions, "each with its
options and the one you recommend", and the session that started it must "ask
the user and continue it with the answers". With nobody present, the PRD
author may settle only what "the available sources" settle, and the session
may not take the recommended option, because "A recommendation, silence, a
dismissal or a timeout is never an answer." A question that no source settles
therefore stopped its Spec, and nothing recorded that it stopped or why.

Authoring with nobody present is the usual case in this repository: 41 of the
43 PRDs of Specs 0364 to 0408 name the autonomous route, and 29 hold no
answer of a person. In the three answer records of 2026-11-18 and 2026-11-19
the maintainer took the option labeled recommended in 29 of 33 answers. Both
counts were made on 2026-11-19, the first from the Git history.

On 2026-11-19 the operator asked the maintainer what to do with such a
question while he is away. For the two Specs then being authored he chose
"Seguir o padrão do autor (Recommended)", with the decision recorded "como
decisão do operador por delegação, com pergunta e opções" and his limits,
"dados para fora, escopo novo, gasto", still stopping the Spec. For the
product he asked: "Antes de iniciar o fluxo, perguntar se pode seguir o
padrão de autorização decidido pelo autor ou se deve esperar; seria bacana no
projeto ter uma configuração no .roundfixrc.yml com a preferencia por projeto
e o default ser a spec para e espera o autor."

He settled the shape in the Idea stage of the same day, each answer the
option he chose: the phases, "PRD, TechSpec e Tasks (Recommended)"; the
question and the configuration, "Config responde; sem ela, pergunta
(Recommended)"; how far a Spec goes, "Até o fim, com relato (Recommended)";
what tells a session that nobody is present, "O que a sessão foi informada
(Recommended)"; the number of adopted options, "Sem teto, com contagem
(Recommended)"; and the status view, "Enquanto a Spec estiver ativa
(Recommended)".

Published practice is divided. Tools that pause for a person wait without
limit and read silence as a refusal: AG2
(https://docs.ag2.ai/docs/user-guide/tools/approval_required/) says "an
unanswered approval is not an approval". Appelo's sixth level of authority
(https://www.informit.com/articles/article.aspx?p=1675546&seqNum=3) is the
one this decision takes when a project opts in: "You let the team decide. And
afterwards you inquire about their motivations, or you ask that they actively
keep you informed." Cuadron and Gupta (https://arxiv.org/abs/2611.02207)
measured what an unnoticed mistake costs in chains of code requests, where
agents "complete 95% of requests taken one at a time, but when requests are
chained together, they complete 51%".

## Decision

**An Open Authoring Question is one no source settles.** It is a question the
author of a PRD, a TechSpec or a Task Graph needs answered to finish its
phase, that neither the request, the brief, the approved Idea, a recorded
answer of the user nor another source of the repository settles. The author
states it with two or three options and the one it recommends, with the
reason.

**With the user present, the session asks.** Nothing changes there: a
question the user can see is asked and waits for him, and no recommendation,
silence, dismissal or timeout adopts anything. A session treats the user as
absent only when it was told so: by the brief that started it, by the user's
own words, or because it runs as a continuation that the `unattended` Skill
Mode armed.

**A project that declares nothing waits.** The Authoring Question Preference
is one value of Project Config, wait or follow, and an absent value means
wait. Only Project Config holds it, only as the repository's default branch
holds it, and no session, Agent or command of Roundfix writes it. In a
project that declares nothing, the `unattended` Skill Mode asks once before
the flow, with wait recommended, and the answer holds for the Specs of that
flow.

**Following adopts the author's recommendation as a Delegated Decision.**
With nobody present and follow in force, the session adopts the recommended
option of a question that is not reserved, and the phase goes on. There is no
cap on their number; the session's report counts them.

**A Reserved Question always waits.** Whatever the preference, nobody answers
for the user: the Idea stage and the choice an Idea Stage Judgment leaves to
him; a Rule Challenge; every authorization; a recommended option that adds
scope, spends, sends data outside or contradicts a recorded answer of the
user; a question without two or three options and a recommended one; and any
question the author or the session cannot classify. The author says for each
question whether it is reserved and why, and the session checks that
statement before it adopts anything.

**The record names itself.** A Delegated Decision is an entry of the PRD's
Clarification Record with its date, the phase that asked, the question, the
options, the adopted option with the author's reason, the grant it rests on
and the user's later answer. It is never written as a dated answer of a
person or under the autonomous-route line. The status view of the Follow-up
Command lists the entries of every active Spec.

**The Spec goes to the end, and a later answer wins.** A Spec with Delegated
Decisions is authored and delivered inside the authorization it already
holds. When the user answers a delegated question, his answer joins the
entry; a different answer amends the plan before the Spec's delivery starts
and follows the existing routes after it.

**A Spec that waits, waits alone.** No artifact of that phase is written, an
open Backlog Entry that names the Spec holds the whole question under a
Waiting Declaration, and other work goes on.

## Considered Options

- **Keep the silence.** Rejected: the Spec stopped with no record, and each
  cycle needed a one-time answer that lived in one session's notes.
- **Follow by default.** Rejected by the maintainer: "o default ser a spec
  para e espera o autor". A repository that never chose is never built on an
  answer nobody gave.
- **A timeout that adopts the recommended option while the user is present.**
  Rejected: a question the user can see waits for him. The Apache projects'
  lazy consensus
  (https://community.apache.org/committers/decisionMaking.html), where
  "Silence indicates consent", rests on an announced intent and a waiting
  time among peers; here the absent person is the one who decides.
- **A preference in User Config.** Rejected: the request asks for a
  preference of the project, and a machine-wide switch would follow a person
  into repositories that never chose it.
- **Read the preference from the working tree or from the branch a session
  is on.** Rejected: an Agent could switch it with an edit or a commit of its
  own. The key is read where an Item Branch cannot change it, as
  `delivery.require_defect_class` already is.
- **A hold before delivery.** Rejected by the maintainer, "Até o fim, com
  relato (Recommended)": a hold turns every following Spec into a waiting
  Spec one step later.
- **Record the adopted option as a dated answer.** Rejected: the glossary
  tells authors to avoid assumed approval, and nobody could tell the user's
  answers from the session's.
- **A cap on Delegated Decisions per phase.** Rejected by the maintainer,
  "Sem teto, com contagem (Recommended)".
- **Delegate the Idea stage too.** Rejected by the maintainer: the Idea is
  where scope is decided.

## Consequences

- A project that declares nothing behaves as before, except that a wait now
  has a record: a Backlog Entry that names the Spec and holds the question.
- In a project that follows, a wrong option can reach the TechSpec, the Tasks
  and the code before the user reads it. The record, the report and the
  status view are what make it cheap to overturn.
- The author that classifies a question is the Agent that goes on when the
  question passes. The session is a second reader and also an Agent; the
  person's check is reading the records.
- A teammate's declared preference binds the unattended sessions of everyone
  who works in the repository. A user who is present is always asked.
- The sentence that a recommendation "is never an answer" stays true: a
  Delegated Decision is not an answer and is never recorded as one.
- ADR-0280 is amended: the Clarification Record gains the Delegated Decision
  as a third entry, and its two existing forms and their check are unchanged.
