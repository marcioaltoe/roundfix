---
spec: 0138-a-matrix-the-spec-declares
status: archived
created: 2026-09-14
surfaces: [backend, docs]
archived: "2026-09-16"
source_slug: 0138-a-matrix-the-spec-declares
---


# A matrix the Spec declares

Today the QA gate decides its own matrix. The qa-gate skill and the QA contract
the Daemon puts in the gate's prompt both tell the executor to cover every
promise, story and criterion. As a result, two executors of the same Spec cover
different sources with different inputs. Spec 0119 needed six QA Reports. Only
two of its five failed reports found a defect in the Spec's own work. The rest
carried noise the matrix produced:

- six rows blocked in cascade by one finding their checks did not depend on;
- an aggregate row that re-counted other rows' findings, failing in one report
  and passing in the next on the same kind of finding;
- a changed-path audit the gate rebuilt by hand, choosing its own revision rule,
  because nothing told it which commits the mechanical stage had audited;
- row inputs widened after the row ran.

Specs 0134 to 0137 wrote narrow verification Requirements into their `qa` Tasks,
and those gates still found real defects. This Spec makes that practice the gate
contract:

- the Spec declares what the matrix must cover, by a rule two executors read the
  same way;
- every row names the sources it covers;
- a failed check blocks only the planned rows that depend on it;
- no row re-checks what a Mechanical Refusal Code already decides.

## Project Constraints

- Identifier strategy: applicable — the QA Report frontmatter keys, the three typed blocked causes, the Mechanical Refusal Codes and matrix row identifiers keep their spelling and meaning; this Spec introduces no identifier, refusal code or report key. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential, network request or HTTP surface is created or read; the change is guidance for the gate executor and the text of its prompt. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the gate's matrix, verdicts and evidence are governed by accepted decisions this Spec must preserve. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0080 applies: QA verdicts distinguish environment-blocked rows, and every row still ends in one of the typed blocked causes or a result.
  ADR-0088 applies: the QA gate is authored into the graph, which is where the declared coverage is read.
  ADR-0091 applies: the QA gate is a Task node of its own type, so the coverage is declared by the `qa` Task.
  ADR-0096 applies: the QA gate proves machine facts before it spends an agent turn, so the stage still withholds the Agent Session on a blocking fact, no row re-checks a fact the stage decides, and no stage makes a verdict more permissive.
  ADR-0097 applies: a QA row carries forward only on declared unmoved evidence, so row inputs are fixed when the row is planned.
  ADR-0104 applies: a Spec accepts on evidence it did not author, so the outside-evidence row stays a source every matrix covers.
  ADR-0117 applies: a defect is checked by the stage that can produce it, and the commit-dependent audit stays a gate fact executed as a command.
  ADR-0155 applies: the `qa` Task declares the gate's matrix, which this Spec writes into the qa-gate skill and the QA contract.
  ADR-0093 is not applicable to this change: Spec consistency is checked by citation, never by inference, and this Spec changes no Spec Consistency Check rule or citation reader.
- Tooling authority: applicable — express maintainer authorization: "Aprovar como proposto", 2026-09-14, recorded in `docs/specs/0138-a-matrix-the-spec-declares/_authorization.md`; bounded files: `.agents/skills/qa-gate/SKILL.md`, `skills/qa-gate/SKILL.md`. Sanctioned regeneration: `make skills-sync` regenerates the mirror, and `make baseline-digests` rewrites any derived pin the approved skill edit moves. The QA contract in the gate prompt is ordinary source that no authorization has bounded. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## Goals

- Two executors of the same Spec cover the same set of sources, because that set
  follows from what the Spec declared under a textual rule.
- Once the matrix exists, a failed check blocks only the rows whose entry point,
  observable or evidence depends on it.
- No matrix row re-checks a fact a Mechanical Refusal Code already decides.
- No defect the current matrix found becomes unreachable.

## User Stories

1. As a Supervisor authoring a Spec, I want the verification Requirements of my
   `qa` Task to define exactly what the gate covers, so that a narrow declared
   matrix can pass instead of reading as a partial gate.
2. As a Supervisor reading a QA Report, I want every row to name the sources it
   covers, so that I can check the matrix is complete and adds nothing else.
3. As a Supervisor reading a failed QA Report, I want every finding-blocked row
   to name the check it waits on, so that one finding cannot hide the results of
   rows that never depended on it.
4. As a Supervisor rerunning a gate, I want each row's inputs fixed when the row
   is planned, so that carry-forward depends on what the row set out to read,
   not on what it found.
5. As a Supervisor whose `qa` Task names an analyzer beyond the repository
   Verification, I want that analyzer scoped to what the Spec changed, so that a
   pre-existing diagnostic elsewhere cannot fail the Spec.

## Core Features

1. **The `qa` Task declares its coverage by a textual rule.**
   - A numbered `qa` Task Requirement that starts with `MUST verify` or
     `MUST run` is a declared verification Requirement.
   - Every other Requirement, such as `MUST exercise` or `MUST NOT`, constrains
     how the gate runs.
   - A `qa` Task with at least one declared verification Requirement has
     declared its matrix. Its coverage sources are those Requirements plus the
     sources in Core Feature 2.
   - When a declared Requirement verifies every Acceptance Criterion of named
     Tasks, each of those criteria is a coverage source of its own.
   - A gate run on a declared matrix is a full gate and can reach `pass`.
     `partial` stays reserved for a `qa` Task that explicitly marks its run
     partial, and for a planned row the gate did not run.
2. **Some coverage sources no declaration waives.** Every matrix covers:
   - the outside-evidence row;
   - the Pull Request row, on its default equivalent-evidence path;
   - the repository Verification;
   - each declaration under the PRD's Unreachable Acceptance section, read as
     the author's claim to test;
   - the frontend sweep, when the PRD declares the `frontend` surface.
3. **An undeclared matrix keeps a bounded default.** When the `qa` Task declares
   no matrix, the coverage sources are:
   - the PRD user stories and Goals;
   - each non-QA Task's Acceptance Criteria;
   - each declared intentional break;
   - the sources in Core Feature 2.

   A Non-Goal becomes a source only when it names an observable the Spec could
   ship. Its row reads only the paths the Spec changed and the surfaces the
   Non-Goal names.
4. **Coverage is complete, closed and named.**
   - Every row names, in its provenance, each coverage source it covers.
   - Every coverage source appears in at least one row's provenance.
   - No row's provenance names a source outside the coverage sources.
   - Rows may group several sources that one observation settles.
   - No row's result is computed from other rows' results.
   - No row re-checks what a Mechanical Refusal Code already decides, such as
     the report's shape or its evidence paths.
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
   the repository Verification is recorded as a failure, never as an environment
   block.
   - "Contention", "non-reproducible" and the same failure on the unchanged
     delivery target are not proved environmental causes.
   - Only a cause that stops the command from running at all is environmental.
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
   - the non-waivable sources, including the frontend sweep;
   - named, complete and closed coverage;
   - scoped blocking.

## User Experience

A Supervisor opens a QA Report and reads each row's provenance. Every source the
Spec declared, and every non-waivable source, appears in at least one row, and
nothing else does. Each finding-blocked row names the check it waits on. A rerun
covers the same sources as the round before, with carried rows marked as such.

## Non-Goals / Out of Scope

- Changing verdict semantics, the typed blocked-cause counts, QA Report keys,
  report naming, or which report is newest. The report-sequence disagreement
  captured from the fleet on 2026-09-11 is separate work.
- Fixing how many rows a matrix has or how rows group sources. Only the covered
  set is fixed.
- Adding a mechanical detector, including one that checks coverage or proves
  row inputs were declared before execution.
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

- A `qa` Task that declares its matrix no longer has rows for PRD stories, Core
  Features or Non-Goals that its verification Requirements do not name.
- A `qa` Task Requirement that does not start with `MUST verify` or `MUST run` no
  longer adds a coverage source.
- A row whose provenance names no coverage source is no longer part of the
  matrix, and no row re-counts other rows or re-checks report shape or evidence
  paths.

## Regression locks

- Every verdict rule reaches the same verdict for the same row outcomes; no rule
  becomes more permissive.
- A repository Verification failure is still a `fail`.
- Every matrix, declared or default, covers the outside-evidence row, the Pull
  Request row and every PRD Unreachable Acceptance declaration.
- The gate still audits Task commits against the grant by command, with the
  clauses the skill's contract tests pin.

## Acceptance evidence

The outside-evidence row replays evidence this Spec did not write:

- the six QA Reports of Spec 0119;
- one failed and one passing QA Report from Spec 0136;
- Spec 0134's recorded analyzer failure;
- the QA Reports of Specs 0134, 0135 and 0137.

Under the new contract, the replay must show five things:

1. **Spec 0119's real defects stay reachable.** Each still reaches a covered
   source or a mechanical finding: the missing operations list, the dropped
   template clauses, the rewritten setup context, the Verification timeout and
   the out-of-grant changed paths.
2. **Spec 0119's noise does not recur.** The six cascaded blocks and the
   aggregate row do not occur.
3. **Spec 0136's real defect stays covered.** Its governed-rename push finding
   came from a `qa` Task Requirement that starts with `MUST verify`, so its row
   covers a declared source. Its report-closure row covers no source.
4. **Spec 0134's analyzer failure resolves.** Its diagnostics were identical on
   the delivery target, so its analyzer row passes with those diagnostics
   attributed to their owner.
5. **Earlier reports show the provenance shape.** The QA Reports of Specs 0134,
   0135 and 0137 already group several criteria in one row and name them in the
   row. The coverage rule adopts that shape and does not claim to reproduce
   their row sets. The replay records, for each of those reports, which of its
   `qa` Task's declared sources no row named.

## Success Metrics

- The replay above holds for every source it names.
- The QA Reports of the next two Specs delivered after this one show:
  - every coverage source named in some row's provenance;
  - no row naming a source outside the set;
  - no finding-blocked row that names no dependency;
  - no aggregate row;
  - no row whose inputs grew after it ran.

## Decisions

- **Order of work.** On 2026-09-14 the maintainer chose to narrow the gate
  before remodelling authoring in Spec 0129. The accepted risk is that weight
  moves onto Acceptance Criteria that Spec 0129 has not yet strengthened.
- **Declaration rule.** A `qa` Task Requirement that starts with `MUST verify`
  or `MUST run` declares a coverage source, the form Specs 0134 to 0137 already
  use. See ADR-0155.
- **Coverage, not row shape.** On 2026-09-15 the maintainer chose a coverage
  rule over one row per Requirement. Three review rounds had each found a new
  place where a row-shape rule needed interpretation: whether a Requirement
  names what the gate verifies, merging with obligations, and plural
  Requirements.
- **Analyzer scope.** Analyzer rows are scoped to changed packages, and
  diagnostics identical on the delivery target are attributed to their owner,
  as Spec 0137 authored.
- **Repeated audit stays for now.** Removing it needs the mechanical stage to
  report the commits it audited, so it is follow-up work.

## Research basis

**Secondbrain.** Consulted before authoring. The session read `wiki/index.md`
and ran
`qmd query "QA gate matrix scope cascade blocked rows cost rerun roundfix" --all --files --min-score 0.3`.
These sources informed the decision:

- `wiki/concepts/verificacao-adversarial-e-oraculos-de-agentes.md`. A gate is
  trustworthy only when it observed the right property and can fail a known
  negative. This supports keeping the outside-evidence row non-waivable. It
  also supports rejecting aggregate rows, whose pass observes nothing of its
  own.
- `inbox/roundfix/_triaged/2026-08-08-atrito-medido-do-loop-autonomo-em-tres-specs.md`.
  Vortex needed fourteen gate runs for three Specs. Eleven did not pass, all for
  contract reasons rather than business logic, yet the gate still caught four
  real defects. This is why the Spec narrows noise without loosening any
  verdict.
- `inbox/roundfix/2026-09-11-o-qa-novo-reutiliza-um-numero-anterior-ao-relatorio-vigente.md`.
  A report-sequence disagreement, kept out of scope.
- `inbox/roundfix/2026-09-10-events-aborta-e-carry-forward-quebra-no-worktree.md`.
  It concerns Run carry-forward, not row carry-forward, and did not change the
  design.

**Exa MCP.** Consultation was attempted, but no Exa MCP tool was available in
this session. No external source was read for this Spec, so no external
validation is claimed. The published work cited inside the Secondbrain concept
page was not re-read.

**Independent review.** Three rounds of Codex pre-PR review
(`gpt-5.6-luna`, effort `max`) raised twelve findings, each confirmed in the
repository before it was accepted. The coverage rule, the scoped blocking, the
timing rule and the Unreachable Acceptance source came from those rounds.

## Open Questions

- Should a `qa` Task's analyzer row stay scoped to changed packages or cover the
  whole repository? Default: changed packages, as in Core Feature 8, until the
  maintainer decides otherwise.
