---
spec: 0138-a-matrix-the-spec-declares
status: active
created: 2026-09-14
surfaces: [backend, docs]
---

# A matrix the Spec declares

Today the QA gate decides its own matrix. The qa-gate skill and the QA contract
the Daemon puts in the gate's prompt both tell the executor to cover every
promise, story and criterion. As a result, two executors of the same Spec build
different row sets with different inputs. Spec 0119 needed six QA Reports. Only
two of its five failed reports found a defect in the Spec's own work. The rest
carried noise the matrix produced:

- six rows blocked in cascade by one finding their checks did not depend on;
- an aggregate row that re-counted other rows' findings, failing in one report
  and passing in the next on the same kind of finding;
- a changed-path audit the gate rebuilt by hand, choosing its own revision rule,
  because nothing told it which commits the mechanical stage had audited;
- row inputs widened after the row ran.

Specs 0134 to 0137 wrote narrow matrices into their `qa` Tasks, and those gates
still found real defects. This Spec makes that practice the gate contract:

- the Spec declares the matrix, by a rule two executors read the same way;
- a failed check blocks only the planned rows that depend on it;
- no row re-checks what a Mechanical Refusal Code already decides.

## Project Constraints

- Identifier strategy: applicable — the QA Report frontmatter keys, the three typed blocked causes, the Mechanical Refusal Codes and matrix row identifiers keep their spelling and meaning; this Spec introduces no identifier, refusal code or report key. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential, network request or HTTP surface is created or read; the change is guidance for the gate executor and the text of its prompt. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the gate's matrix, verdicts and evidence are governed by accepted decisions this Spec must preserve. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0080 applies: QA verdicts distinguish environment-blocked rows, and every row still ends in one of the typed blocked causes or a result.
  ADR-0088 applies: the QA gate is authored into the graph, which is where the declared matrix is read.
  ADR-0091 applies: the QA gate is a Task node of its own type, so the matrix is declared by the `qa` Task.
  ADR-0096 applies: the QA gate proves machine facts before it spends an agent turn, so the stage still withholds the Agent Session on a blocking fact, no row re-checks a fact the stage decides, and no stage makes a verdict more permissive.
  ADR-0097 applies: a QA row carries forward only on declared unmoved evidence, so row inputs are fixed when the row is planned.
  ADR-0104 applies: a Spec accepts on evidence it did not author, so the outside-evidence row stays in every matrix.
  ADR-0117 applies: a defect is checked by the stage that can produce it, and the commit-dependent audit stays a gate fact executed as a command.
  ADR-0155 applies: the `qa` Task declares the gate's matrix, which this Spec writes into the qa-gate skill and the QA contract.
  ADR-0093 is not applicable to this change: Spec consistency is checked by citation, never by inference, and this Spec changes no Spec Consistency Check rule or citation reader.
- Tooling authority: applicable — express maintainer authorization: "Aprovar como proposto", 2026-09-14, recorded in `docs/specs/0138-a-matrix-the-spec-declares/_authorization.md`; bounded files: `.agents/skills/qa-gate/SKILL.md`, `skills/qa-gate/SKILL.md`. Sanctioned regeneration: `make skills-sync` regenerates the mirror, and `make baseline-digests` rewrites any derived pin the approved skill edit moves. The QA contract in the gate prompt is ordinary source that no authorization has bounded. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## Goals

- Two executors of the same Spec plan the same matrix rows, because every row
  traces to a source the Spec declared under a textual rule.
- Once the matrix exists, a failed check blocks only the rows whose entry point,
  observable or evidence depends on it.
- No matrix row re-checks a fact a Mechanical Refusal Code already decides.
- No defect the current matrix found becomes unreachable.

## User Stories

1. As a Supervisor authoring a Spec, I want the verification Requirements of my
   `qa` Task to be the gate's complete matrix, one row each. A narrow declared
   matrix can then pass instead of reading as a partial gate.
2. As a Supervisor reading a failed QA Report, I want every finding-blocked row
   to name the check it waits on. One finding then cannot hide the results of
   rows that never depended on it.
3. As a Supervisor rerunning a gate, I want each row's inputs fixed when the row
   is planned. Carry-forward then depends on what the row set out to read, not
   on what it found.
4. As a Supervisor whose `qa` Task names an analyzer beyond the repository
   Verification, I want that analyzer scoped to what the Spec changed. A
   pre-existing diagnostic elsewhere then cannot fail the Spec.

## Core Features

1. **The `qa` Task declares the matrix, by a textual rule.**
   - A numbered `qa` Task Requirement that starts with `MUST verify` or
     `MUST run` is exactly one matrix row.
   - A `qa` Task with at least one such Requirement has declared its matrix.
     Those Requirements are the complete matrix, and the gate adds no rows of
     its own.
   - Every other Requirement, such as `MUST exercise` or `MUST NOT`, constrains
     how the gate runs and is not a row.
   - A gate run on a declared matrix is a full gate and can reach `pass`.
     `partial` stays reserved for a `qa` Task that explicitly marks its run
     partial, and for a planned row the gate did not run.
2. **Some obligations no declaration waives.** A declared matrix still carries:
   - the outside-evidence row;
   - the Pull Request row, on its default equivalent-evidence path;
   - the repository Verification;
   - the frontend sweep, when the PRD declares the `frontend` surface.
3. **An undeclared matrix keeps a bounded default.** When the `qa` Task declares
   no matrix, the gate derives rows from:
   - the PRD user stories and Goals;
   - each non-QA Task's Acceptance Criteria;
   - each declared intentional break;
   - the obligations in Core Feature 2.

   A Non-Goal becomes a row only when it names an observable the Spec could
   ship. That row reads only the paths the Spec changed and the surfaces the
   Non-Goal names.
4. **Every row has its own observable.**
   - No row's result is computed from other rows' results.
   - No row re-checks what a Mechanical Refusal Code already decides, such as
     the report's shape or its evidence paths.
   - The report's coverage and verdict sections carry those totals.
5. **Once the matrix exists, a finding blocks only its dependents.**
   - A failed check blocks the planned rows whose entry point, observable or
     evidence depends on that check. Each blocked row names the check it waits
     on.
   - Every other row runs and records its own result.
   - No gate instruction may require a governance failure to stop the flow
     rows.
   - The mechanical stage still withholds the Agent Session when a blocking
     machine fact is present before any matrix exists. This feature does not
     change that.
6. **Row inputs are fixed at planning.** Each row declares its inputs when it is
   written as `pending`, bounded to what the row reads. A row whose inputs grow
   after it ran cannot be carried forward.
7. **Timing failures are not dismissed.** A timeout or intermittent failure of
   the repository Verification is code-caused unless the unchanged delivery
   target reproduces it under the same command. "Contention" or
   "non-reproducible" is not a proved environmental cause.
8. **A named analyzer reads what the Spec changed.** An analyzer the `qa` Task
   names beyond the repository Verification runs over the packages the Spec
   changed.
   - A diagnostic identical on the delivery target is recorded as observed,
     names the Spec that owns it, and does not fail the Spec.
   - A repository Verification failure stays blocking whether or not it is
     pre-existing.
9. **The prompt states the same contract.** The QA contract the Daemon places in
   the gate prompt replaces its instruction to validate every story and
   criterion with:
   - the declaration rule;
   - the bounded default;
   - the non-waivable obligations, including the frontend sweep;
   - finding-dependent blocking.

## User Experience

A Supervisor opens a QA Report and can trace every row to a `qa` Task
Requirement that starts with `MUST verify` or `MUST run`, or to one default
source. Each finding-blocked row names the check it waits on. A rerun shows the
same row set as the round before, with carried rows marked as such.

## Non-Goals / Out of Scope

- Changing verdict semantics, the typed blocked-cause counts, QA Report keys,
  report naming, or which report is newest. The report-sequence disagreement
  captured from the fleet on 2026-09-11 is separate work.
- Adding a mechanical detector, including one that proves row inputs were
  declared before execution.
- Changing when the mechanical stage withholds the Agent Session.
- Changing what the mechanical stage computes: which Task commits it receives,
  which revision authorizes them, or whether its report names the commits it
  audited.
  - The stage audits only the current Run's Task commits, and its report does
    not name them.
  - The gate therefore keeps auditing every Task commit by command.
  - Removing that repeated audit is recorded for follow-up.
- Changing the authoring skills that write `qa` Tasks. Spec 0129 owns
  authoring.
- Widening a grant to cover a Spec's own corrective Tasks. That cost Spec 0119
  more QA rounds than any matrix rule did, and it is Spec 0129's class.
- Loosening the repository Verification, the Pull Request row, the
  outside-evidence row or the frontend sweep.
- Adding repository linter or analyzer configuration.

## Declared intentional breaks

- A `qa` Task that declares its matrix no longer gets rows for PRD stories, Core
  Features or Non-Goals that its verification Requirements do not name.
- A `qa` Task Requirement that does not start with `MUST verify` or `MUST run`
  no longer becomes a row.
- The gate no longer plans rows that re-count other rows or re-check report
  shape or evidence paths.

## Regression locks

- Every verdict rule reaches the same verdict for the same row outcomes; no rule
  becomes more permissive.
- A repository Verification failure is still a `fail`.
- The outside-evidence row and the Pull Request row appear in every matrix,
  declared or default.
- The gate still audits Task commits against the grant by command, with the
  clauses the skill's contract tests pin.

## Acceptance evidence

The outside-evidence row replays QA Reports this Spec did not write: the six
reports of Spec 0119, one failed and one passing report from Spec 0136, and
Spec 0134's recorded analyzer failure.

**Real defects stay reachable.** Under the new contract, each of Spec 0119's
real defects still reaches a planned row or a mechanical finding:

- the missing operations list;
- the dropped template clauses;
- the rewritten setup context;
- the Verification timeout;
- the out-of-grant changed paths.

**Noise does not recur.** The six cascaded blocks and the aggregate row do not
occur.

**Spec 0136.** Its governed-rename push finding came from a `qa` Task
Requirement that starts with `MUST verify`, so it remains a planned row. Its
report-closure row is no longer a matrix row.

**Spec 0134.** Its analyzer diagnostics were identical on the delivery target,
so its analyzer row passes with those diagnostics attributed to their owner.

**The declaration rule reproduces what Specs 0134 to 0137 authored.** Applying
it to their `qa` Tasks yields the rows those Specs planned, and treats their
`MUST exercise` and `MUST NOT` Requirements as constraints.

## Success Metrics

- The replay above holds for every row it names.
- The QA Reports of the next two Specs delivered after this one show:
  - no finding-blocked row that names no dependency;
  - no aggregate row;
  - no row that does not trace to a declared source;
  - no row whose inputs grew after it ran.

## Decisions

- **Order of work.** On 2026-09-14 the maintainer chose to narrow the gate
  before remodelling authoring in Spec 0129. The accepted risk is that weight
  moves onto Acceptance Criteria that Spec 0129 has not yet strengthened.
- **Where the matrix is declared.** The declared matrix is read from `qa` Task
  Requirements that start with `MUST verify` or `MUST run`. That is the form
  Specs 0134 to 0137 already use, so no new section is needed and the authoring
  skills stay unchanged until Spec 0129. See ADR-0155.
- **Analyzer scope.** Analyzer rows are scoped to changed packages, and
  diagnostics identical on the delivery target are attributed to their owner,
  as Spec 0137 authored.
- **Repeated audit stays for now.** Removing it needs the mechanical stage to
  report the commits it audited, so it is follow-up work.

## Open Questions

- Should a `qa` Task's analyzer row stay scoped to changed packages or cover the
  whole repository? Default: changed packages, as in Core Feature 8, until the
  maintainer decides otherwise.
