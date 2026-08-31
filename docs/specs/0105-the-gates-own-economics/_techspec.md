---
spec: 0105-the-gates-own-economics
prd: _prd.md
created: 2026-08-31
---

# The gate's own economics — Technical Spec

## Executive Summary

Five changes, each removing a cost the gate charges for something other than
finding a defect. A Spec crossing an external boundary characterizes the real
thing before its premise reaches the gate. The Pull Request row carries the
equivalent-evidence path by default instead of every Spec rediscovering it. The
QA Task's Verification is derived by Roundfix rather than hand-authored over a
verdict the Daemon already computed. A governance finding blocks the rows it
concerns instead of the whole matrix. And the citation parser accepts the forms
Specs are actually written in.

The primary trade-off is that Core Feature 3 removes an author's control over
one Task's Verification. That is deliberate: the measured failure was a
hand-written predicate that passed a gate which had failed, and an author who
can write that predicate can write it wrong in a way no reviewer catches,
because it reads as ordinary Task authoring. The cost is a Task whose contract
is no longer visible in its own file; the TechSpec answers that by having the
derived command rendered into the Task file rather than hidden, so the file
still shows what will run while no longer being where it is decided.

The second trade-off concerns Core Feature 4 and is a narrowing, not a
widening. ADR-0096 withholds the Agent Session when a blocking machine fact is
present, and that stays. This Spec changes only what happens once a matrix
exists: the measured round built nineteen rows and blocked fifteen of them on
one governance finding, which is cascade rather than withholding. A finding
blocks the rows it names and no others.

## Project Constraints

- Identifier strategy: applicable — QA Report, verdict, blocked-row causes,
  Unreachable Acceptance, and Characterization are glossary terms whose
  obligations change here, and a message that invents a synonym for one of them
  is a defect. This Spec coins no term; the closing node checks whether the work
  introduced, changed, or retired one the glossary should carry. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no network boundary, credential,
  request, or transport is created or read. A characterization Task may reach a
  real boundary, but it does so through the repository's existing command
  surface and declares no new transport policy. Source:
  `docs/agents/specific-repository.md`.
- Active ADR obligations: applicable — ADR-0080 supplies the typed
  blocked-cause distinction the Pull Request row must use rather than each Spec
  rediscovering it, and bounds what an equivalent-evidence path may accept.
  ADR-0091 keeps the gate one terminal Task node of type `qa`, which is the node
  whose Verification becomes derived. ADR-0088 authors the gate into the graph
  rather than requesting it per run, which is why the Pull Request row is
  unreachable by construction. ADR-0097 lets a row carry forward only on
  declared, unmoved evidence, which is the shape the equivalent-evidence path
  follows. ADR-0096 withholds the Agent Session when a blocking machine fact is
  present, and this Spec narrows nothing there: Core Feature 4 changes only how
  a finding attributes blame across a matrix that already exists. ADR-0117
  places a check with the stage that can produce its defect, which is why
  characterization moves to authoring. ADR-0093 checks consistency by citation
  rather than inference, which bounds the parser change to recognising written
  forms and never to inferring intent. ADR-0104 makes a Spec accept on evidence
  it did not author, which is the rule the characterization change serves.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — this Spec edits three Roundfix-owned skills,
  which are protected tooling. Express maintainer authorization: granted
  2026-08-12, recorded at
  `docs/workflow/authorizations/2026-08-12-the-authoring-and-baseline-corrections.md`
  and re-bounded for this Spec at
  `docs/workflow/authorizations/2026-08-31-the-gate-stops-charging-for-its-own-shape.md`.
  Bounded files: `.agents/skills/qa-gate/SKILL.md`,
  `.agents/skills/write-tasks/SKILL.md`,
  `.agents/skills/implement-task/SKILL.md`, `skills/qa-gate/SKILL.md`,
  `skills/write-tasks/SKILL.md`, `skills/implement-task/SKILL.md`. The three
  copies under `skills/` are rewritten by the declared `make skills-sync` and
  stay sanctioned fallout under ADR-0081; they are named because the
  changed-path audit reads `paths:` frontmatter and `skills/` carries no
  ownership declaration that would resolve them, which is the omission that
  refused Spec 0118 on 2026-08-27. The record lands as its own commit before any
  skill edit. Source: `docs/agents/agent-instructions.md`.

## System Architecture

Five components change. No package, layer, or directory is added.

**The task-authoring skill** gains the characterization obligation: a Spec whose
work crosses an external surface authors a Task that records what the real
boundary does, and that Task precedes the work depending on it. This is
guidance, and it is the only change here that acts before any code runs.

**The QA gate skill** gains the Pull Request row's default. The row is
unreachable by construction under ADR-0088, so the gate applies ADR-0080's
environment-blocked path to it and requires the equivalent evidence to be
recorded rather than leaving each Spec to discover the arrangement.

**The QA Task's Verification becomes derived.** `internal/spec` already parses a
Task's `## Verification` verbatim. For a Task of type `qa`, Roundfix supplies
the command instead of reading the author's, and renders what it supplies into
the Task file so the contract stays visible where a reader looks for it. An
authored command for a `qa` Task becomes a Spec check finding rather than
silently losing.

**Blocked-row attribution** stops cascading. A mechanical finding already names
the row it blocks; the gate treats a blocking finding as blocking the matrix.
It blocks the rows it names, and the rest are measured.

**The citation parser** accepts a conjunction as a list separator and a decision
number without its prefix when the context is an obligations line. Where it
still cannot recognise a form, the failure names the form it does recognise,
which today it does not.

## Implementation Design

### Interfaces

```go
// DerivedQAVerification is the Verification Roundfix supplies for a Task of
// type qa. An author no longer writes it: the measured failure was a
// hand-written predicate that accepted a verdict outside the domain, and it
// read as ordinary Task authoring to every reviewer.
func DerivedQAVerification(slug string) []string

// AuthoredQAVerification reports a qa Task that carries its own command, so the
// Spec check can refuse it rather than let it lose silently.
func AuthoredQAVerification(task Task) bool
```

```go
// BlockedRows returns the rows one mechanical finding blocks. A finding names
// the row it concerns; a matrix row it does not name is measured, not blocked.
func (finding MechanicalFinding) BlockedRows() []string
```

### Data Models

No schema change and no new persisted entity. The QA Report's row contract,
verdict rules, and typed blocked-cause counts are unchanged: ADR-0080 requires a
blocked row to carry a typed cause and never credits a journey without evidence,
and this Spec's Non-Goals exclude touching any of it.

### API Contracts

`roundfix spec check` — gains one finding for a `qa` Task that authors its own
Verification, and accepts two citation forms it refuses today. No flag changes.

The QA gate — unchanged verdict rules. The Pull Request row is recorded as
environment-blocked with its equivalent evidence rather than passing silently or
capping the verdict. A governance finding blocks the rows it names.

## Coverage Map

| PRD item | Component |
| --- | --- |
| Goal 1, Story 1, Core Feature 1 | The task-authoring skill's characterization obligation |
| Goal 2, Story 2, Core Feature 2 | The QA gate skill's Pull Request row default |
| Goal 3, Story 4, Core Feature 3 | `DerivedQAVerification`, `AuthoredQAVerification` |
| Goal 4, Story 3, Core Feature 4 | `BlockedRows` and the gate's attribution |
| Story 5, Core Feature 5 | The citation parser's accepted forms |

## Integration Points

- `.agents/skills/{write-tasks,qa-gate,implement-task}/SKILL.md` and their
  generated copies, under the recorded authorization.
- `internal/spec/task.go` and `internal/spec/task_type.go` — the derived
  Verification for type `qa`.
- `internal/speccheck/mechanical.go` — blocked-row attribution.
- `internal/speccheck/citations.go` — the accepted citation forms.

## Testing Approach

Every seam exists; none is added. `internal/spec`, `internal/speccheck`, and
`internal/daemon` already cover the surfaces this Spec changes.

- The derived Verification is exercised against the measured failure: a
  hand-authored predicate that accepted a verdict outside the domain no longer
  reaches the Daemon, and a `qa` Task carrying its own command is refused by
  name.
- Blocked-row attribution is table-driven over a finding that names one row in a
  matrix of several, asserting the unnamed rows are measured rather than
  blocked.
- The two measured citation forms parse: a conjunction separator in both
  repository languages, and a decision number without its prefix on an
  obligations line. A form still unrecognised produces a message naming the
  recognised one, asserted on the message rather than on the absence of a
  finding.
- The two skill changes are read against the delivered commands rather than
  against this TechSpec.

Repository Verification is `rtk make verify`; `make verify-docs` covers the
markdown contracts and is required before the pull request opens.

## Build Order

1. Derive the QA Task's Verification, render it into the Task file, and refuse
   an authored one by name.
2. Attribute a blocking mechanical finding to the rows it names rather than to
   the matrix (depends on: none; serialized after 1 by edit locality in the
   checker packages).
3. Accept the two measured citation forms, and name the recognised form when a
   citation still fails (depends on: 2 — both edit `internal/speccheck`).
4. The task-authoring skill carries the characterization obligation, and the QA
   gate skill carries the Pull Request row's default (depends on: 1, 2, 3 —
   documentation describes delivered behavior).
5. QA gate (depends on: 4).

The bounded authorization is not a Build Order step. It lands as its own commit
during authoring, before any commit that edits a skill, because a record the
graph must name has to exist before the graph names it — and authoring it as a
Task would make its Verification vacuous the moment it existed.

**Every documentation step follows every behavior step.** Spec 0118 ordered its
documentation before a corrective Task changed the rule it described and shipped
two contradicting surfaces; Spec 0116 repeated the lesson from the other
direction, applying a rule uniformly to a skill whose stage could not satisfy
it. Both are why steps 5 and 6 sit last.

## Risks & Considerations

**A derived Verification removes author control from one Task.** Accepted
deliberately, and bounded to type `qa`. The mitigation is visibility, not
choice: the derived command is rendered into the Task file, so a reader still
sees what will run. An authored command is refused by name rather than
overwritten, because silently replacing an author's text is how a contract
becomes invisible.

**The Pull Request row could become a silent pass.** ADR-0080 already refuses to credit a
row without evidence, and this Spec requires the equivalent evidence to be
recorded. A row with no recorded evidence stays blocked exactly as today.

**Blocked-row attribution could under-block.** A governance finding that
genuinely invalidates the whole matrix must still be able to say so. The
mechanism is the finding naming its rows, so a finding that names all of them
blocks all of them — the change removes an implicit cascade, not the ability to
block widely.

**The parser could start inferring.** ADR-0093 bounds it to recognising written
forms. A number without its prefix is accepted only on an obligations line,
where the context makes the form unambiguous; nowhere else.

## Decisions

- The Pull Request row's equivalent-evidence path is a gate default rather than
  a template clause, because a template is per-Spec discovery again. Recorded in
  the PRD.
- The derived Verification is rendered into the Task file rather than left
  implicit, so removing the author's control does not also remove the reader's
  view.
- An authored `qa` Verification is refused rather than overwritten.
- This Spec's authorization is recorded fresh rather than by citing the
  2026-08-12 grant alone, because that record names canonical skills only and
  the changed-path audit reads `paths:`.
- No new ADR is minted. Every decision applies an accepted record; none reverses
  one.
