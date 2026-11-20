---
status: accepted
knowledge: true
created_at: 2026-11-19T15:00:00Z
updated_at: 2026-11-19T15:00:00Z
deprecated_at: null
superseded_by: null
---

# A Lesson Learned is written only when a test holds, and leaves when a rule states it

A Lesson Learned is one file under `docs/lessons/` in a fixed form. An Agent
writes it only when a five-part test holds, at two moments of the work, in the
change it is already making, so a person sees it in the Pull Request. The file
is deleted once an ADR, a guide or a skill states what it teaches.

## Context

The only rule about lessons was one sentence of the cross-project clause:
"When a Finding closes, state whether its change teaches a lesson that its
ADR, Archive Record or reference does not already state. Only then send one
message with its source pinned". The lesson had no form and no file, and it
reached the knowledge base only when one session remembered to send it and
another filed it. The archive of a Spec asked nothing. The Finding that held
the proven cause is deleted when it closes. The one count ever made, on
2026-10-14, found that 7 of 383 retired files held a lesson no ADR states,
and a session rescued them by hand.

The maintainer asked on 2026-11-18 for lessons that are worth sending to the
knowledge base, and gave the measure: "Só seria útil se fosse algo como:
tentamos usar x lib e gerou problemas tal e resolvemos assim." He chose the
moments, "Finding fechado e archive c/ tropeço (Recommended)", and who
decides, "Agente grava, você vê no PR (Recommended)".

Published evidence says a store of lessons fails in two ways. It decays when
nobody curates it: an audit of NASA's lessons system found that managers
"rarely consult or contribute" to it, and Nick Milton concludes that "lessons
need to be embedded into design, into principles and into process"
(http://www.nickmilton.com/2017/07/leaving-lessons-in-lessons-database.html).
And lessons without a filter make a coding agent worse: a memory of lessons
distilled from past tasks "does not outperform AutoCodeRover in any
configuration" (Lindenbauer, Groh and Schütze,
https://aclanthology.org/2025.realm-1.30.pdf).

## Decision

**The record.** One file per lesson, `docs/lessons/YYYY-MM-DD-<kebab-slug>.md`,
with a `created` date in its header, a title of one line and eight labelled
lines in this order: Context, Tried, Went wrong, Cause, Solved by, Applies
when, Does not apply when, Links. The labels stay as the Baseline writes them
and the rest is in the language of the repository's documentation.

**The test.** A Lesson Learned is written only when every part holds:
something was tried, it went wrong, the cause is proven, something solved it,
and no accepted ADR, reference document or glossary entry of the repository
already states it. When any part fails, nothing is written.

**Two moments ask the question.** When a Finding closes, in the change that
answers it or retires it, and when a Spec that needed a correction, a failed
QA or a retry is archived. For a Finding a Spec adopted, the change that
answers it is that Spec's delivery, so the question is asked at its archive.

**The Agent writes, and a person sees the Pull Request.** The Agent that
applies the test writes the file in the change it is making. Whoever reviews
may delete it. No judge filters it and no approval is asked before the file
exists.

**The file replaces the message.** A closing Finding no longer sends its
lesson to another project: the mirror carries the file where the Secondbrain
is turned on, and the repository keeps it either way.

**A lesson leaves when a rule takes it in.** Once an accepted ADR, a guide or
a skill states what a Lesson Learned teaches, the file is deleted under the
retirement rule, and Git keeps its text.

**The Delivery Queue asks through the Supervisor.** The queue archives with no
Agent, so the session that supervises a delivery applies the test when it
declares a Correction or retries an item after a repair, and writes the lesson
in the commit that corrects. After the merge it writes the lesson in its next
change.

**Nothing fails for a lesson.** The Doctor Command reports a file under
`docs/lessons/` that is off the form, as a warning. No Run, merge or check
stops for a missing or a malformed lesson, and the Archive Record gains no
field.

## Considered Options

- **Keep the message.** Rejected: it depends on two sessions and leaves
  nothing in the repository that learned.
- **A lesson for every closed Finding and every archived Spec.** Rejected by
  the maintainer and by the evidence above: most work teaches no lesson, and
  an empty record hides the ones that have one.
- **A person approves each lesson before it is written.** Rejected by the
  maintainer: the Pull Request already shows it.
- **A judge filters lessons.** Rejected by the maintainer, and it would send
  the text of a repository's failures to an outside service.
- **Keep every lesson for good.** Rejected: Milton's reading is that a stored
  lesson nobody embeds is not learned. The exit keeps the store small.
- **A section in the Archive Record.** Rejected: a lesson must be found
  without knowing which Spec learned it, and a Finding closed without a Spec
  has no record.
- **The QA gate writes the lesson.** Rejected: its commit holds QA evidence
  only, and a Correction can come after it.

## Consequences

- Lessons are rare by design: the earlier count is 1.8 for each 100 retired
  files.
- A hurried review lets a weak lesson through, and the published measurements
  say weak lessons hurt. The test, the form's two lines on when a lesson
  applies and the exit are the guards.
- A correction the queue makes without a session, such as a Corrective Task,
  is asked about only after the merge.
- The Release Plan counts `docs/lessons/` as planning evidence, so a commit
  that only adds a lesson asks for no release.
- This decision changes the paragraph "Lessons leave before the text" of
  [ADR-0278](0278-a-retirement-deletes-and-the-history-root-holds-only-archive-records.md)
  and adds the Lesson Learned to what its retirement rule deletes.
