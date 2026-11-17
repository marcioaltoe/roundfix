---
status: accepted
created_at: 2026-11-17T18:45:00Z
updated_at: 2026-11-17T20:30:00Z
deprecated_at: null
superseded_by: null
---

# A Correction carries a declared Defect Class, and a repository requires it only by its own setting

Roundfix records why a delivery was corrected as a declaration by whoever
made the correction: one Correction per cause, with a Defect Class from a
closed list of five and a one-line reason. The class is an opinion, kept
apart from the Cause Class that signatures assign, and a repository makes it
required with a Project Config setting that is off by default. It is declared
because no record and no signature can tell a defect of the plan from a
defect of the implementation, and it is optional by default because no
adopter's delivery may change by surprise.

## Context

On 2026-11-17 the maintainer asked whether Fable is effective for Ideas, PRDs
and TechSpecs compared with Opus. The operator answered with a study of Specs
0370 to 0403 that rebuilt the cause of every correction by hand from harness
transcripts, commit subjects, review artifacts and memory notes, and left
three of 23 delivered Specs with a question mark. Roundfix records a
correction only as a Review Finding Disposition, which says that a finding
was fixed or dismissed and never why it existed. `roundfix runs causes` lists
failed Verification attempts and Corrective Tasks inside Runs: for the Runs
created since 2026-11-09 it listed 10 items in 31 Runs while the study
counted 24 corrections, and its signatures left 100 of 132 items without a
class in the measurement of 2026-09-17 to 2026-10-02.

The maintainer decided the shape on 2026-11-17, each answer verbatim: the
unit, "Uma por causa (Recommended)"; who declares, "Quem corrige
(Recommended)"; whether it is required, "Só onde é ligado e por default é
desligado, nos projetos da gesttione vai ser configurado para ligar"; the
artifact of a plan defect, "Sim (Recommended)"; this repository and the fifth
class, "Ligar e manter unclear (Recommended)". After the TechSpec, at about
17:10 BRT, he answered how a commit by hand that applies his own decision and
corrects no defect is declared: "Sexta resposta (Recommended)", a sixth answer
`not-a-correction` with a reason, which keeps the commit out of the defect
rates instead of forcing it into `operator`.

The pre-PR review of Spec 0398 returned one defect of the plan and one of the
implementation in the same round, and the delivery of Spec 0399 had three
corrections and five causes. This repository squash-merges each delivery: the
nine commits of pull request #892 reached the default branch as one.

## Decision

**A Correction is one cause.** From the moment a Spec's plan is on the
default branch until its delivery merge, each cause that needed a repair is
one Correction: a pre-PR review finding that was fixed, a QA finding that was
repaired, or another defect that needed a change to the delivery candidate or
to the plan. A commit that repairs two causes holds two Corrections, and two
commits for one cause are one. A dismissed finding, a failed Verification
attempt the Agent repaired inside its Run, a stop that needed no repair and a
plan change that records a new decision of the maintainer are not
Corrections.

**The Defect Class.** A Correction carries one class and a reason of one
line. `plan`: the plan said something wrong, missing or impossible.
`implementation`: the plan was right and the delivered change was not.
`environment`: the cause lies outside the Spec, such as the machine, a tool
or a failure already on the default branch. `operator`: an action of the
Supervisor or of a person during the delivery caused it. `unclear`: the
declarer cannot tell, and the reason says what is known. A `plan` Correction
also names the artifact that was wrong, the PRD, the TechSpec or the Tasks; a
defect that comes from the Idea names the PRD, which adopts it.

**A sixth answer: `not-a-correction`.** A declaration may answer
`not-a-correction` instead of a Defect Class, with a reason and no artifact,
for a Correction Sign whose change repaired no defect, such as a commit by
hand that applies a new decision of the maintainer. The declaration is
recorded like any other and covers the signs it names, so they are not
`undeclared`, and the report counts it apart: it enters no count of
Corrections, of plan defects or of a comparison. It is an answer of the
declaration, not a sixth Defect Class; the list of Defect Classes keeps its
five.

**Whoever corrects declares, and the last declaration counts.** The
Supervisor or person who makes the correction declares it: together with the
disposition of a fixed review finding, or in a declaration of its own. The
record keeps who declared and when. Anyone may declare the same Correction
again, the maintainer above all, or withdraw one declared by mistake; the
last declaration counts and the earlier ones stay in the record.

**Correction Signs.** Roundfix knows without a declaration that a delivery of
the Delivery Queue was corrected from three recorded facts: a review finding
with a fixed disposition, a commit on the Item Branch that neither a Run nor
the Delivery Queue made and that is not a merge of the default branch, and a
Corrective Task. A declaration names the signs its Correction covers. A sign
that no standing Correction names counts as one `undeclared` Correction.

**Required only where a repository turns it on.** A Project Config setting,
off by default and ignored in User Config, makes the declaration required.
With it on, the Delivery Queue does not merge an item while one of its
Correction Signs is named by no Correction: the item parks, and its Park line
names each sign, the five classes and the command that declares. With it off
nothing is required, a class that is given is recorded, and nothing a
repository printed before changes. A merge made outside the Delivery Queue is
never refused. This repository turns the setting on once the release that
knows the key is its installed binary.

**Apart from Cause Class and Park Class.** A Cause Class is still assigned by
signature to a failed Verification attempt or a Corrective Task, and a Park
Class still says why an item stopped. Neither changes, no report merges them
with a Defect Class, and `roundfix runs causes` prints what it printed.

**On this machine, within Run Retention.** The declarations and the signs are
kept under Roundfix Home for the Run Retention window and outlive the Delivery
Queue item and the Item Branch they came from. Nothing is added to the
Archive Record, and no text leaves the machine.

## Considered Options

- **A class by signature, as an extension of Cause Class.** Rejected: a
  signature can be shown and repeated, and it left 100 of 132 items without a
  class; no diagnostic text says whether the plan or the change was wrong.
- **A class from the judge.** Rejected for this version: the text of a
  correction would leave the machine, and the Secondbrain's page on evaluation
  records a judge model agreeing with people at kappa 0.32 to 0.53.
- **One Correction per commit, or per review round.** Rejected by the
  maintainer: only one per cause keeps a defect of the plan and a defect of
  the implementation of the same round apart.
- **A `Defect-Class:` trailer on the fix commit as the record.** Rejected as
  the only carrier: a squash merge keeps one commit, and the trailer never
  reaches the default branch.
- **Required in every repository.** Rejected by the maintainer.
- **A refusal at the moment of the fix.** Rejected: a correction commit runs
  no Roundfix command that could refuse, so the merge is the one point every
  correction passes.
- **Declaring a commit that applies a decision of the maintainer as
  `operator`.** Rejected by the maintainer ("Sexta resposta (Recommended)"):
  `operator` would count a decision as a defect of the delivery and raise the
  defect rates; `not-a-correction` records the commit and keeps it out.
- **Four classes, without `unclear`.** Rejected by the maintainer, so that
  nobody is made to guess. Orthogonal Defect Classification warns that a label
  of where a defect was injected is an opinion (Chillarege et al., 1992).
- **The maintainer, or a second reader, declares.** Rejected: whoever
  corrects holds the context at that moment and costs the maintainer no time.
  The later declaration keeps a second reading possible.
- **Facts that travel with the Spec in Git.** Not taken in this version: the
  Archive Record stays as it is.

## Consequences

- The declarer is often the session that wrote the plan, and it is asked
  whether the plan was at fault. The reason, the declarer and every later
  declaration stay visible, and the report counts the Corrections whose class
  a later declaration changed.
- A repository that leaves the setting off collects a thin record, and the
  report says how many Corrections are `undeclared`.
- The signs are read only for a delivery the Delivery Queue made. A Spec
  delivered by hand shows what was declared and nothing else.
- A declarer who answers `not-a-correction` for a repair hides a defect from
  the rates. The reason stays visible, the report counts these answers per
  Spec, and a later declaration may give the sign a Defect Class.
- A merge of the default branch into the Item Branch is not a sign. A repair
  made inside one is a Correction only when someone declares it.
- A Roundfix older than the release that carries the setting refuses a
  Project Config that holds a key it does not know, so a repository turns the
  setting on after every machine that runs its Roundfix is upgraded.
- The history reaches as far back as Run Retention, 30 days by default.
