---
spec: 0138-a-matrix-the-spec-declares
status: active
created: 2026-09-14
surfaces: [backend, docs]
---

# A matrix the Spec declares

The QA gate decides its own matrix today. The qa-gate skill and the QA contract
the Daemon puts in the gate's prompt both tell the executor to cover every
promise, story and criterion. So two executors of the same Spec build different
row sets and set different inputs for them. Spec 0119 needed six QA Reports.
Only two of its five failed reports found a defect in the Spec's own work. The
rest carried noise the matrix produced:

- six rows blocked in cascade by one finding their checks did not depend on;
- an aggregate row that re-counted other rows' findings, failing in one report
  and passing in the next on the same kind of finding;
- a changed-path audit the gate rebuilt by hand, choosing its own revision
  rule, while nothing told it which commits the mechanical stage had audited;
- row inputs widened after the row ran.

Specs 0134 to 0137 wrote narrow matrices into their `qa` Tasks, and those gates
still found real defects. This Spec makes that practice part of the gate
contract. The Spec declares the matrix, a failed check blocks only the rows that
depend on it, and the gate stops redoing work the mechanical stage already does.

## Project Constraints

- Identifier strategy: applicable — the QA Report frontmatter keys, the three typed blocked causes, the Mechanical Refusal Codes and matrix row identifiers keep their spelling and meaning; this Spec introduces no identifier, refusal code or report key. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential, network request or HTTP surface is created or read; the change is guidance for the gate executor and the text of its prompt. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the gate's matrix, verdicts and evidence are governed by accepted decisions this Spec must preserve. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0080 applies: QA verdicts distinguish environment-blocked rows, and every row still ends in one of the typed blocked causes or a result.
  ADR-0088 applies: the QA gate is authored into the graph, which is where the declared matrix is read.
  ADR-0091 applies: the QA gate is a Task node of its own type, so the matrix is declared by the `qa` Task.
  ADR-0096 applies: the QA gate proves machine facts before it spends an agent turn, so the gate stops re-deriving what the mechanical stage computes, and no stage makes a verdict more permissive.
  ADR-0097 applies: a QA row carries forward only on declared unmoved evidence, so row inputs are fixed when the row is planned.
  ADR-0104 applies: a Spec accepts on evidence it did not author, so the outside-evidence row stays in every matrix.
  ADR-0117 applies: a defect is checked by the stage that can produce it, and the commit-dependent audit stays a gate fact executed as a command.
  ADR-0155 applies: the `qa` Task declares the gate's matrix, which this Spec writes into the qa-gate skill and the QA contract.
  ADR-0093 is not applicable to this change: Spec consistency is checked by citation, never by inference, and this Spec changes no Spec Consistency Check rule or citation reader.
- Tooling authority: applicable — express maintainer authorization: "Aprovar como proposto", 2026-09-14, recorded in `docs/specs/0138-a-matrix-the-spec-declares/_authorization.md`; bounded files: `.agents/skills/qa-gate/SKILL.md`, `skills/qa-gate/SKILL.md`. Sanctioned regeneration: `make skills-sync` regenerates the mirror, and `make baseline-digests` rewrites any derived pin the approved skill edit moves. The QA contract in the gate prompt is ordinary source that no authorization has bounded. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## Goals

- Two executors of the same Spec plan the same matrix rows, because every row
  traces to a source the Spec declared before the gate ran.
- A failed check blocks only the rows whose entry point, observable or evidence
  depends on it.
- The gate spends no Agent effort re-deriving a fact the mechanical stage
  computes.
- No defect the current matrix found becomes unreachable.

## User Stories

1. As a Supervisor authoring a Spec, I want the Requirements of my `qa` Task to
   be the gate's complete matrix. A narrow declared matrix can then pass instead
   of reading as a partial gate.
2. As a Supervisor reading a failed QA Report, I want every finding-blocked row
   to name the check it waits on. One finding then cannot hide the result of
   rows that never depended on it.
3. As a maintainer, I want the gate to take the changed-path audit of every
   commit the mechanical stage already audited from that stage's findings. The
   gate then spends no turn repeating an audit the Daemon computed.
4. As a Supervisor rerunning a gate, I want each row's inputs fixed when the row
   is planned. Carry-forward then depends on what the row set out to read, not
   on what it found.
5. As a Supervisor whose `qa` Task names an analyzer beyond the repository
   Verification, I want that analyzer scoped to what the Spec changed. A
   pre-existing diagnostic elsewhere then cannot fail the Spec.

## Core Features

1. **The `qa` Task declares the matrix.** When the Spec's `qa` Task Requirements
   name what the gate verifies, those Requirements are the complete matrix: the
   gate plans a row for each and adds none of its own. A Requirement that only
   tells the gate to run declares nothing. A gate run on a declared matrix is a
   full gate and can reach `pass`. `partial` stays reserved for a `qa` Task that
   explicitly marks its run partial, and for a planned row the gate did not run.
2. **Some obligations no declaration waives.** A declared matrix still carries:
   - the outside-evidence row;
   - the Pull Request row, on its default equivalent-evidence path;
   - the repository Verification;
   - the frontend sweep, when the PRD declares the `frontend` surface.
3. **An undeclared matrix keeps a bounded default.** When the `qa` Task does not
   declare its matrix, the gate derives rows from the PRD user stories and
   Goals, each non-QA Task's Acceptance Criteria, each declared intentional
   break, and the obligations in Core Feature 2. A Non-Goal becomes a row only
   when it names an observable the Spec could ship. That row reads only the
   paths the Spec changed and the surfaces the Non-Goal names.
4. **Every row has its own observable.** No row's result is computed from other
   rows' results. The same holds for a row that re-checks what a Mechanical
   Refusal Code already decides, such as the report's shape or evidence paths.
   The report's coverage and verdict sections carry those totals.
5. **A finding blocks only its dependents.** A failed check blocks the rows whose
   entry point, observable or evidence depends on that check, each naming the
   check it waits on. Every other row runs and records its own result. No
   instruction may require a governance failure to stop the flow rows.
6. **The gate does not repeat the mechanical audit.** For each Task commit the
   mechanical stage audited, the gate reads the stage's changed-path,
   authorization and consequent-fix findings and does not re-derive them. The
   gate audits by command only the Spec's Task commits the stage did not
   receive, and the report names those commits.
7. **Row inputs are fixed at planning.** Each row declares its inputs when it is
   written as `pending`, bounded to what the row reads. A row whose inputs grow
   after it ran cannot be carried forward.
8. **Timing failures are not dismissed.** A timeout or intermittent failure of
   the repository Verification is code-caused unless the unchanged delivery
   target reproduces it under the same command. "Contention" or
   "non-reproducible" is not a proved environmental cause.
9. **A named analyzer reads what the Spec changed.** An analyzer the `qa` Task
   names beyond the repository Verification runs over the packages the Spec
   changed. A diagnostic identical on the delivery target is recorded as
   observed, with the Spec that owns it named, and does not fail the Spec. A
   repository Verification failure stays blocking whether or not it is
   pre-existing.
10. **The prompt states the same contract.** The QA contract the Daemon places in
    the gate prompt describes the declared matrix and the bounded default,
    replacing its instruction to validate every story and criterion.

## User Experience

A Supervisor opens a QA Report and can trace every row to a `qa` Task
Requirement or to one default source. Each finding-blocked row names the check
it waits on. The audit of changed paths against the grant appears once per Task
commit: in the mechanical findings for the commits the stage audited, and as a
gate row only for the commits it did not receive. A rerun shows the same
row set as the round before, with carried rows marked as such.

## Non-Goals / Out of Scope

- Changing verdict semantics, the typed blocked-cause counts, QA Report keys,
  report naming, or which report is newest. The report-sequence disagreement
  captured from the fleet on 2026-09-11 is separate work.
- Adding a mechanical detector, including one that proves row inputs were
  declared before execution.
- Changing what the mechanical stage computes, including which Task commits it
  receives and which revision authorizes them. It audits only the current Run's
  Task commits, so a Spec delivered across several Runs leaves earlier commits
  to the gate. That gap is recorded for follow-up.
- Changing the authoring skills that write `qa` Tasks. Spec 0129 owns
  authoring.
- Widening a grant to cover a Spec's own corrective Tasks, which cost Spec 0119
  more QA rounds than any matrix rule did. That is Spec 0129's class.
- Loosening the repository Verification, the Pull Request row, the outside-
  evidence row or the frontend sweep.
- Adding repository linter or analyzer configuration.

## Declared intentional breaks

- A `qa` Task that declares its matrix no longer gets rows for PRD stories, Core
  Features or Non-Goals its Requirements do not name.
- The gate no longer plans rows that re-count other rows, re-check report
  shape, or repeat the changed-path audit for commits the mechanical stage
  audited.

## Regression locks

- Every verdict rule reaches the same verdict for the same row outcomes; no rule
  becomes more permissive.
- A repository Verification failure is still a `fail`.
- The outside-evidence row and the Pull Request row appear in every matrix,
  declared or default.

## Acceptance evidence

The outside-evidence row replays QA Reports this Spec did not write. The
replay source is the six reports of Spec 0119, one failed and one passing report
from Spec 0136, and Spec 0134's recorded analyzer failure. Under the new
contract:

- **Spec 0119's real defects stay reachable.** Each one still reaches a planned
  row or a mechanical finding: the missing operations list, the dropped template
  clauses, the rewritten setup context, the Verification timeout, and the
  out-of-grant changed paths.
- **Spec 0119's noise does not recur.** The six cascaded blocks and the
  aggregate row do not occur.
- **Spec 0136's real defect stays in the matrix.** Its governed-rename push
  finding came from a declared `qa` Task Requirement, so it remains a planned
  row. Its report-closure row is no longer a matrix row.
- **Spec 0134's analyzer failure passes as pre-existing.** It came from
  diagnostics identical on the delivery target, so its analyzer row passes with
  those diagnostics attributed to their owner.

## Success Metrics

- The replay above holds for every row it names.
- The QA Reports of the next two Specs delivered after this one show:
  - no finding-blocked row that names no dependency;
  - no aggregate row;
  - no changed-path audit repeated for a commit the mechanical stage audited;
  - no row whose inputs grew after it ran.

## Decisions

- QA first, remodelling later: the maintainer chose on 2026-09-14 to narrow the
  gate before Spec 0129. The accepted risk is that weight moves onto Acceptance
  Criteria that Spec 0129 has not yet strengthened.
- The declared matrix is read from the `qa` Task's Requirements, the form Specs
  0134 to 0137 already use, rather than a new section. The authoring skills stay
  unchanged until Spec 0129. See ADR-0155.
- Analyzer rows are scoped to changed packages and attribute diagnostics
  identical on the delivery target to their owner, as Spec 0137 authored.

## Open Questions

- Should a `qa` Task's analyzer row stay scoped to changed packages or cover the
  whole repository? Default: changed packages, as in Core Feature 9, until the
  maintainer decides otherwise.
