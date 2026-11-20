---
status: accepted
knowledge: true
created_at: 2026-11-17T18:45:00Z
updated_at: 2026-11-17T18:45:00Z
deprecated_at: null
superseded_by: null
---

# A phase entry holds the model and reasoning effort that ran an Authoring Phase

A phase entry of an Authoring Phase records the model that ran the phase, its
reasoning effort, how Roundfix learned each of them and whether the phase's
Author Subagent or the session itself ran it. The session reports the model,
because nothing a command can read names it, and the Corrections Report reads
these entries. This amends ADR-0407, under which an entry held no environment
value and only the status view read the record.

## Context

The delivery side records the runtime, the model and the reasoning effort of
every Run. The authoring side records neither. On 2026-11-17 the Phase Record
of the Mac mini held 15 phase entries for Specs 0399 to 0404, all with the
same session and none with a model, so the study of Fable and Opus as
authoring models read harness transcripts to learn which model wrote each
plan, and found none for the Specs authored on another machine.

Measured on Claude Code 2.1.325 on 2026-11-17, in the Author Subagents of the
Idea and of the PRD of Spec 0404: a command the subagent runs reads
`CLAUDE_CODE_SESSION_ID`, which holds its parent's session, and
`CLAUDE_EFFORT`, which held `xhigh`, the effort this repository declared for
both phases; no variable's value named a model. The session itself is told
its model.

The Authoring Model Declaration says which model a phase should run on, not
which one ran: a fix is authored in the session on the session's model, and
in a repository that publishes an authoring skill a typed skill runs on the
session's model too (ADR-0408).

## Decision

**What a phase entry gains.** An entry of `idea`, `prd`, `techspec` or
`tasks` records the model as the harness names it, which the session that
runs the phase reports; the reasoning effort, read from the session's
environment where the harness exposes it and otherwise as the session reports
it; for each of the two, how it was learned, `reported`, `environment` or
`unknown`; and whether the phase's Author Subagent or the session ran the
phase, as reported. The keys are added to the schema
`roundfix/phase-record/v1`, and an entry written before them reads as
`unknown`.

**Only for an identified session.** The values are recorded when the
environment names a Claude Code session, the one harness the command
identifies. Any other entry holds `unknown`.

**Never a refusal.** A value that is missing, too long or malformed is
recorded as `unknown` with one warning, and the entry is still written. The
command's refusals and exit codes stay as they are, it still loads no
configuration, and it still reads no earlier entry (ADR-0414).

**The repository an entry belongs to.** An entry of an Authoring Phase also
records, when it is written, an identity of the repository that outlives the
working tree it was written in, because an author works in a clone that is
removed once the plan merges. It never records a credential that a remote
address may carry.

**The declaration is not the source.** The command does not read the
Authoring Model Declaration, and the declaration does not change.

**Which entry gives a phase its model.** For a delivered Spec, the model of a
phase is that of the last entry of that phase recorded before the Spec's
first Implement Run began, or before its delivery merge when the machine
holds no Run of the Spec.

**Who reads the entries.** The Corrections Report reads phase entries, as the
status view of the Follow-up Command does. No command reads a prompt entry or
a harness transcript, prices a token or computes a limit percentage.

**Where it is stated.** The four authoring skills tell the session and the
Author Subagent to report the model and who runs the phase when they record
it, and the Roundfix skill's phase reference and the user guide describe the
entry.

## Considered Options

- **The Authoring Model Declaration as the source.** Rejected: the declared
  model is not always the one that ran, and the command loads no
  configuration.
- **Reading the harness transcript.** Rejected: ADR-0407 keeps every command
  away from the transcripts, and a transcript on another machine is not
  there to read.
- **The model without the effort.** Rejected: this repository declared the
  same model at two efforts for different phases, and the Run side records
  both.
- **The effort only as the session reports it.** Rejected: the session is not
  told its effort, and the environment held it in both sessions measured.
- **A new schema version.** Rejected: added keys keep every reader of
  `roundfix/phase-record/v1`, the reason ADR-0414 kept the `boundary` key.
- **A reported value from any harness.** Not taken: the maintainer approved
  the record for Claude Code sessions, where the measurement was made.
- **Matching an entry to a repository by asking the working tree when the
  report runs**, as ADR-0410 does within seven days. Rejected for a report
  over the Run Retention window: the clone is gone by then.

## Consequences

- A reported model can be wrong, and a session may report none. The entry
  says how each value was learned, and the report counts the Specs whose
  model is `unknown`.
- Two spellings of one model are two groups in the report, so the skills ask
  for the exact identifier the harness states and never an alias.
- A Spec authored on another machine has no entry on this one and reads
  `unknown`.
- ADR-0407's sentence that an entry never holds an environment value now has
  one exception, the reasoning effort.
- The entries stay in the Phase Record's day files, which the Run Retention
  Sweep removes past its cutoff.
