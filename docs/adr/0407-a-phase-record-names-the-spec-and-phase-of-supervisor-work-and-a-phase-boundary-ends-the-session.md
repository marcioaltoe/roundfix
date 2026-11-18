---
status: accepted
created_at: 2026-11-16T11:00:00Z
updated_at: 2026-11-18T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A Phase Record names the Spec and phase of Supervisor work, and a phase boundary ends the session

Roundfix appends to one record under Roundfix Home which Spec and phase each
Supervisor session works on and which Work Item each prompt of a Run served.
The skills and the autonomous-work guide tell the Supervisor to end its
session at each phase boundary, and Roundfix says so, refusing nothing, when a
session goes on. The record exists because the usage limit of an Agent harness
is spent in sessions nothing describes, and the rule is an instruction because
Roundfix cannot end a session.

## Context

On 2026-11-14 jarvis reported that the maintainer's weekly Claude limit ran
out before its reset, placed the weight in long orchestration sessions, and
asked Roundfix to write the Run and the phase in a fixed place so that a usage
script could attribute the limit by phase. Measured on the Mac mini from
2026-11-08 to 2026-11-15: 94.2% of the Claude cost ran in sessions opened in a
repository checkout and 0.5% in a Roundfix worktree; 75% of what a main
session cost was cache read, and 59.2% of its calls carried 200 thousand
tokens or more; every one of the 23 ACP session identifiers Roundfix had
recorded for pre-PR reviews equalled the session identifier of a Codex
transcript. A command run by a Claude Code 2.1.325 session reads the session's
identity from `CLAUDE_CODE_SESSION_ID`, and a subagent's command reads its
parent's. The maintainer refused automatic compaction and decided that the
record covers Runs and Supervisor sessions, that Roundfix only exposes it, and
that the boundary is an instruction with a notice.

## Decision

**The Phase Record.** Roundfix appends one JSON line per entry to a day file,
named by the UTC date, in the directory `phase-record` under Roundfix Home,
with the schema `roundfix/phase-record/v1`. The directory is the user's alone
and a file is never rewritten. An entry holds identifiers and times: a Spec
slug, a repository root, a working directory, a harness name, a session
identifier, a Run identifier, a Work Item scope. It never holds prompt text, a
token count, a price or an environment value. The Run Retention Sweep removes
a day file older than its cutoff.

**A phase entry.** `roundfix phase start <phase>` records that the session
that runs it starts a Supervisor Phase: `follow-up`, `idea`, `prd`,
`techspec`, `tasks`, `delivery` or `release`, with the Spec it names. The
Roundfix skill and the authoring skills run it; no read-only command writes
the record. The harness session is read from the environment; without one the
entry is written with none.

**A prompt entry.** When a prompt of an Agent Session ends, the Daemon records
its Run, its scope (a Task, the QA gate, a review Batch), its working
directory, its runtime, the ACP session identifier the adapter reported and
when it began and ended. `roundfix review` records the same for each prompt of
a pre-PR review, which runs outside a Run.

**The Phase Boundary and its notice.** A Phase Boundary follows each Authoring
Phase and the delivery and release of a Spec. There the Supervisor writes the
Session Handoff, says what stays armed and ends the session, and the next
session starts with the follow-up Skill Mode. When `roundfix phase start`
finds that the same session already recorded another phase past a boundary,
it prints the Boundary Notice on standard error, marks the entry and exits as
it would have: the notice never refuses and never changes an exit code.

**Roundfix reports nothing from it.** No command reads a harness transcript,
prices a token or computes a limit percentage. A write to the record that
fails costs one warning and never changes a Run, a review or a command's
result.

**Amended by [ADR-0410](0410-the-status-view-lists-specs-being-authored-and-open-backlog-entries.md) (2026-11-16).** The status view of the Follow-up Command reads the phase entries of the last seven days and names the latest Authoring Phase of a Spec being authored and when it was recorded. It reads no prompt entry and reports no usage: no command reads a harness transcript, prices a token or computes a limit percentage, and a failed write to the record still changes no command's result.

**Amended by [ADR-0414](0414-a-supervisor-session-ends-only-when-the-maintainer-asks.md) (2026-11-17).** No phase of a Spec ends the Supervisor session, which writes the Session Handoff and ends only when the maintainer asks or the work stops. `roundfix phase start` reads no earlier entry, writes `boundary` as `false` without `previous_phase` or `previous_spec`, and prints no Boundary Notice. The Phase Record, its entries and the refusal of nothing stay.

**Amended by ADR-0418 (2026-11-17).** An entry of an Authoring Phase also
records the model and reasoning effort that ran it, how each was learned,
whether its Author Subagent or the session ran it, and its Repository
Identity. The effort is the one environment value an entry may hold; a value
not learned is unknown and never costs the entry. The Corrections Report
reads phase entries, while no command reads prompt entries or harness
transcripts, prices tokens or computes a limit percentage. See
[ADR-0418](0418-a-phase-entry-holds-the-model-and-reasoning-effort-that-ran-an-authoring-phase.md).

**Amended by ADR-0428 (2026-11-17).** The Phase Record also holds
fallback entries: the phase, Spec, repository and Repository Identity,
harness session, declared model and effort, fallback model and effort, limit
name and bounded limit line that caused an Authoring Fallback. A Spec keeps
that fallback for later phases whose declared model matches the limited
family, while the next Spec starts on its declared model. Phase-entry readers
skip fallback entries without counting them as invalid; `phase start` stays
unchanged. See [ADR-0428](0428-a-spec-whose-phase-fell-back-keeps-the-fallback-for-the-limited-model-until-the-spec-ends.md).

## Considered Options

- **A record written only by a Run**, as the request asked. Rejected: it
  would describe half a percent of the Claude cost measured.
- **The Run Database with a read command.** Rejected: an outside reader would
  need the binary or bind to a table, and the request was a fixed place.
- **The Artifact Directory.** Rejected: a repository can move it, so the
  place would depend on each repository's configuration.
- **Recording from the commands the Supervisor already runs.** Rejected: the
  Follow-up Command and its status view write nothing, an Idea stage runs no
  Roundfix command, and a phase would be named after it began.
- **A refusal at a boundary.** Rejected: a small fix is rightly authored and
  delivered in one session, and Roundfix cannot tell.
- **Roundfix opening the next session**, or a command the repository declares
  for it. Not taken now: the maintainer chose instruction and notice.
- **Labels set where the session starts**, as the harness's telemetry does.
  Rejected: a label cannot change inside a session, and one session crosses
  phases today.

## Consequences

- An unattended cycle stops at a Phase Boundary: the Delivery Queue goes on
  delivering what is armed, and the next phase waits for a session a person
  opens.
- A subagent shows its parent's identity, so the record alone cannot tell a
  phase the session wrote from one it handed to an Author Subagent; the
  transcript's agent type does.
- A harness that gives a command no session identity is recorded without one
  and receives no notice.
- A prompt that never returns leaves no prompt entry; its usage stays
  attributable to its Run by the worktree's name, as before.
- The report by phase is the reader's work. It exists once the usage script's
  owner reads the record.
