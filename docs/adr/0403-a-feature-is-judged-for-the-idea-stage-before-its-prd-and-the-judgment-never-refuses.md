---
status: accepted
created_at: 2026-11-14T15:00:00Z
updated_at: 2026-11-14T16:00:00Z
deprecated_at: null
superseded_by: null
---

# A feature is judged for the Idea stage before its PRD, and the judgment never refuses

Before the PRD of a feature is written, three criteria about the request are
judged, by the judge when it answers and by the agent when it does not. An
indicated Idea stage starts without being asked and the user may dismiss it
with a recorded reason; nothing is ever refused, and every feature PRD records
what was judged and what followed.

## Context

The Spec routing rule sent only a "large or fuzzy product initiative" through
`write-idea` and told the agent to "start with the smaller sufficient route
when two routes fit". The agent decided alone, nothing recorded why the stage
was skipped, and `write-prd` never sent a request back. In practice the stage
ran when the user started it.

On 2026-11-14 the request behind Spec 0392 named its own shape, and the Idea
stage, started by hand, measured that the shape could not work and changed it
before any PRD existed. On 2026-11-13 a study asked for by hand led the
maintainer to decline a Homebrew tap. The same day the maintainer asked for "um
julgamento (usando o JEV) sobre a necessidade do uso do write-idea para
features/specs mesmo sem o acionamento do usuário", chose a new Spec in the
product, and settled the rule in the Idea stage of that Spec: the waste to
avoid first is "Feature na forma errada"; an indicated stage "Roda, com saída
registrada"; without the judge "O agente aplica os critérios"; and success is
"Menos correções depois".

Published evidence pulls both ways. A stage that is always required becomes
ceremony, and a recommendation from an automated aid tends to be accepted
without being checked (Parasuraman and Manzey, 2010), less so when the aid
shows its confidence and informs instead of recommending (Goddard and
colleagues, 2012).

## Decision

**Three criteria, judged before the PRD.** For a feature that has no Idea and
that the user did not send to `write-idea`: the right shape of the solution is
still open; the request rests on an assumption about its users, the business
or the market that it gives no evidence for; deciding to build it depends on a
market or competitive question. They are the criteria the maintainer named,
and they live in the Spec routing rule, so they hold in a repository without a
judge key. A fix or a refactor is never judged.

**The judge answers when it can, the agent when it cannot.** For a request
recorded as an item that is not yet a Spec, and for a request of the session
handed in as a file, `roundfix spec judge --stage idea` asks the judge
([ADR-0404](0404-the-idea-stage-judge-reads-one-item-that-is-not-yet-a-spec.md)).
When the judge does not answer, the agent answers the same three criteria as
yes, no or unsure, and says that the judge did not answer. The agent's answers are recorded as the agent's and never
shown as the judge's.

**Three verdicts.** Indicated when any criterion is yes, not indicated when
every criterion is no, undecided otherwise.

**The verdict starts the stage; it refuses nothing.** Indicated: the agent
starts `write-idea` without asking and shows the verdict with each answer, its
probability and who judged; the user may dismiss the stage at its first
question. Undecided, or in doubt: the user decides, seeing the same answers.
Not indicated: the PRD follows. No command fails for the judgment, no command
checks that the stage ran, and a dismissal needs only its recorded reason, as
[ADR-0200](0200-spec-authoring-gets-an-advisory-judge-that-never-gates.md)
requires of every advisory judgment.

**The record is in the PRD.** Every feature PRD carries a fielded
`## Idea-stage judgment` section: the judged item, who judged, each answer,
the verdict, whether the stage ran or was dismissed and why, and what the Idea
did to the shape of the request (changed, split or unchanged). An Idea carries
the same section, with `shelved` when the stage ended the request. The Archive
Record names the Spec's `source_revision`, so a later reader finds the record
in Git without a database.

**The stage is the user's.** No agent answers the Idea stage's questions or an
undecided judgment for the user. When the stage is indicated or undecided and
no user is present, no PRD is written: the feature stays a Backlog Entry that
records the judgment and waits, with `waiting_on` and `since`.

## Considered Options

- **A command that refuses a feature PRD without an Idea.** Rejected: an
  advisory judgment never gates, and the thresholds are unmeasured on enough
  cases to carry a refusal.
- **Ask the user every time, with the judgment as information only.** Rejected
  by the maintainer: more interruptions, and the recorded dismissal is already
  the check on the judge.
- **The rule alone, with the agent answering the criteria and no judge.** Kept
  as the path without a key, rejected as the only path: an agent asked to judge
  alone tends to proceed, and the maintainer asked for the judge.
- **A structural check that every feature PRD carries the record.** Rejected
  for now: it would make `roundfix spec check --strict` refuse a Spec over this
  judgment, which the Idea rules out. The report of the follow-up Spec counts
  the feature Specs without a record instead.
- **The record in the Archive Record or in the PRD's front matter.** Rejected:
  the archive and the front matter readers would change for a field no command
  reads yet, and the PRD at `source_revision` already outlives the Spec folder.

## Consequences

- Adopters receive the reworded route clauses at their next
  `roundfix baseline update`; no clause is added or removed.
- A feature costs one question to the user more often than before: on this
  repository's history about a third of the feature items read indicated and
  about half undecided at the starting thresholds.
- An unattended cycle with the user away authors only the features whose
  judgment reads not indicated; the others wait.
- A default is strong: an indicated stage that is rarely dismissed says little
  about whether it was needed. The record of what the Idea did to the shape is
  what shows it.
- The calibration report and the judgment shown in the follow-up report are a
  later Spec, which reads these records.
