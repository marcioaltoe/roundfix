---
spec: 0140-a-spec-traces-the-promises-it-makes
status: active
created: 2026-09-17
surfaces: [backend, docs]
---

# A Spec traces the promises it makes

A Spec makes two kinds of promise nothing traces. It promises a measurable
outcome after shipping, in its Success Metrics, and it promises a public
interface, in its TechSpec's API Contracts. The Spec Consistency Check reads
neither: its coverage units come from the PRD's user stories and Core Features
alone. So a Spec can declare a metric no Task will ever settle, or drop the
section entirely, and every stage stays silent.

Measured over the twenty most recent Specs: twelve PRDs carry no Success
Metrics section at all; of the eight that declare metrics, three are named by no
Task; twelve TechSpecs declare API Contracts that no rule traces.

That gap used to be absorbed downstream, because the QA gate rebuilt its own
matrix from every promise it could find. Spec 0138 ended that: the gate now
covers what the `qa` Task declares, and the weight moved onto what the Spec
writes down. This Spec makes each declared promise name the Task that will
settle it, checked where it is written, by citation.

## Project Constraints

- Identifier strategy: applicable — the `SC-*` consistency codes are stable and never renumbered once shipped, and the two severities keep their meaning; this Spec adds coverage unit kinds and reuses or adds codes without renaming or renumbering an existing one. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential, network request, or HTTP surface is created or read; the change is authoring guidance plus a read-only local check over Spec artifacts. Source: `docs/agents/agent-instructions.md` and `docs/agents/cli.md`.
- Active ADR obligations: applicable — the authoring checks, their placement and their evidence are governed by accepted decisions this Spec must preserve. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0093 applies: Spec consistency is checked by citation and never by inference, so a promise is covered only where an artifact names it.
  ADR-0117 applies: a defect is checked by the stage that can produce it, so the declaration is reported at the stage that writes it and the coverage at the stage that can settle it.
  ADR-0155 applies: the `qa` Task declares the gate's matrix, which is why an untraced promise is now uncovered rather than absorbed by a gate-derived matrix.
  ADR-0104 applies: acceptance rests on evidence the Spec did not author, which here is the corpus of Specs delivered before it.
  ADR-0156 applies: a declared promise names a consuming Task, or the Spec declares explicitly that none applies and why.
- Tooling authority: applicable — the authoring rules live in Roundfix-owned Skills, so this Spec carries an operative grant. Express maintainer authorization: "Aprovar como proposto" on 2026-09-17, widened the same day to the Task authoring skill after pre-PR review, recorded in [_authorization.md](_authorization.md); bounded files: `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-prd/references/prd-template.md`, `.agents/skills/write-techspec/SKILL.md`, `.agents/skills/write-techspec/references/techspec-template.md`, `.agents/skills/write-tasks/SKILL.md`, `.agents/skills/write-tasks/references/task-template.md`, `skills/write-prd/SKILL.md`, `skills/write-prd/references/prd-template.md`, `skills/write-techspec/SKILL.md`, `skills/write-techspec/references/techspec-template.md`, `skills/write-tasks/SKILL.md`, `skills/write-tasks/references/task-template.md`. The Spec Consistency Check itself is ordinary source that no authorization has bounded, so it needs no grant. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## Goals

- A Spec that declares a Success Metric or an API Contract names the Task that
  will settle it, and the check reports the promise that no Task names.
- A Spec with nothing to measure or no public interface says so explicitly, with
  a reason, instead of dropping the section.
- Every new finding names both the promise and the artifact that omits it; none
  is produced by resemblance.
- Replaying the Specs delivered before this one shows exactly which promises
  they left untraced, and raises nothing against the ones they traced.

## User Stories

1. As a Supervisor authoring a PRD, I want the check to tell me a Success Metric
   has no consuming Task before the Run starts, so that an unowned promise is
   caught while it is still cheap to assign.
2. As a Supervisor authoring a TechSpec, I want each declared API Contract
   traced the way a Core Feature already is, so that a public interface cannot
   reach delivery without a Task that implements it.
3. As a Supervisor whose Spec has nothing to measure, I want to record that
   explicitly with a reason, so that an honest absence is not confused with an
   omission.
4. As a Supervisor reading a check result, I want each finding to name the
   promise, the artifact that should have carried it, and where I declared it,
   so that I can act on the finding without re-reading the whole Spec.

## Core Features

1. **Promises are declared as units.**
   - The PRD's Success Metrics section is a numbered list, or the single entry
     `None.` with the reason none applies.
   - The TechSpec's API Contracts section is a numbered list, or the single
     entry `None.` with the reason none applies.
   - A section that is absent, empty, or written as unnumbered prose is not a
     declaration.
2. **A declared promise is traced like any other unit.**
   - Each numbered Success Metric appears in the TechSpec Coverage Map.
   - Each numbered Success Metric and each numbered API Contract is named in
     some Task's References.
   - A promise nothing names is reported against the artifact that should have
     named it.
3. **Findings are located, never inferred.** A finding names both sides it
   compared: the declared promise and the artifact that omits it. Where the
   check cannot settle a candidate, it reports the weaker severity rather than
   asserting a contradiction.
4. **Each check runs at the stage that can establish it.** The declaration is
   reported at the stage that writes the artifact carrying it. Coverage in
   References is reported at the Task Graph stage, which is the first stage
   where a Task exists to name the promise.
5. **Existing units are untouched.** User stories and Core Features keep their
   current codes, severities and messages, and no existing code changes meaning
   or number.

## User Experience

A Supervisor runs the Spec Consistency Check while authoring. A PRD that
declares three metrics and a Task Graph that names two reports one finding, and
the finding names the third metric and the artifact missing it. A Spec that
declares `None.` with a reason reports nothing. Nothing about the command
changes: it stays read-only, emits no verdict, and edits no artifact.

## Non-Goals / Out of Scope

- Tracing the obligations a Spec inherits rather than the promises it makes.
  Accepted-ADR obligations already have their own listing rules and cascade, and
  giving each obligation a Task and evidence owner is the next slice carved from
  Spec 0129.
- Property-shaped acceptance, narrowest-seam guidance and CI feasibility of new
  test classes. That is a separate slice of Spec 0129.
- Judging whether the named Task actually settles the promise. That is semantic
  review, which Spec 0126 makes configurable and which may be explicitly
  declined.
- Amendment of a falsified premise, stale-QA recovery, and temporal
  prerequisites. Those remain with Spec 0129.
- Changing the QA gate's matrix, its rows, its verdicts or the QA Report. Spec
  0138 owns the matrix.
- Changing Baseline modules or any guide delivered inside setup-context
  markers.
- Retrofitting Specs already archived. They stay byte-identical, and the replay
  only measures them.
- Requiring a metric where none applies, or inventing a numeric target to
  satisfy the rule.

## Declared intentional breaks

- A PRD with no Success Metrics section now produces a finding where the check
  used to be silent. Twelve of the twenty most recent Specs are in that state.
- A TechSpec whose API Contracts section is unnumbered prose no longer counts as
  a declaration, even when it reads as complete.
- A declared promise that no Task names now produces a finding at the Task Graph
  stage, where nothing was reported before.

## Regression locks

- The codes that cover user stories and Core Features keep their identity,
  their severity and the artifacts they compare.
- The Spec Consistency Check stays read-only: it writes no artifact, emits no QA
  verdict and creates no Run.
- No finding is produced without both compared sides being located.
- A Spec whose promises are all traced produces no new finding.

## Acceptance evidence

At least one acceptance row rests on evidence this Spec did not author: the ten
Specs that sit in authoring today, written before this one and under a different
intent. Archived Specs are never swept, so they are measured in the problem
statement and are not an acceptance input. The replay states its prediction
before it runs, and must show three things.

1. **The predicted gaps are the reported gaps.** Seven of those Specs carry no
   Success Metrics section, three carry one written as bullets or a table, and
   all ten carry an API Contracts section written as prose. The sweep reports
   exactly ten metric findings and ten contract findings, and no other new
   finding.
2. **No new error appears.** None of the ten declares a numbered promise today,
   so the coverage obligation produces no `error` against them. A rule that
   turned the queue red on arrival would be the wrong rule.
3. **No archived or active artifact changes.** The replay reads the corpus and
   modifies nothing in it.

Where an archived Spec cannot be replayed, the row is recorded as blocked with
that reason; it never requires human interaction.

## Success Metrics

1. The replay described above holds for every count it names, including the
   prediction recorded before it ran.
2. The next two Specs authored after this one declare both sections and reach a
   clean check with no promise finding, without either Spec adding an invented
   metric.
3. No existing consistency code changes its number, severity or meaning, proved
   by the checker's own corpus expectations.

## Decisions

- **The slice, and its size.** On 2026-09-17 the maintainer chose to carve this
  Spec out of Spec 0129 rather than author 0129 whole, and chose traced promises
  as the first slice. Spec 0129 remains the portfolio record in authoring. The
  precedent is Spec 0130, written as the bounded repair carved from Spec 0120.
- **Promises, not inherited obligations.** A promise the Spec makes and an
  obligation it inherits are different objects with different existing rules, so
  they are separated into different slices.
- **An honest `None.` is a declaration.** Half of the recent Specs are internal
  repairs with nothing to measure after shipping; requiring an invented metric
  would buy ceremony instead of traceability. See ADR-0156.
- **Checked by citation.** Coverage is a named relationship between artifacts,
  never a resemblance between sentences. See ADR-0093.

## Research basis

**Secondbrain.** Consulted before authoring. The session read `wiki/index.md`
and ran
`qmd query "traceability between requirements, acceptance criteria and tests; coverage of success metrics" --all --files --min-score 0.3`
and
`qmd query "spec-driven development: every requirement must have an owner task and evidence; definition of done" --all --files --min-score 0.35`.
The results were dominated by mirrors of this repository's own Specs, which are
references rather than independent knowledge. The one non-mirror source,
`raw/web/2026-08-10-branas-ia-formacao-ia-spec-driven-na-pratica.md`, describes
spec-driven practice end to end but does not address promise-level traceability,
so it neither supported nor challenged this design. No Secondbrain source
changed the decision.

The repository's pending Inbox Entries were read before authoring. None is a
source for this Spec:
`inbox/roundfix/2026-09-08-adocao-indexada-aceita-origem-duplicada.md` is the
nearest, an authoring-stage check that accepts a duplicate adopted source, and
it is a different rule about a different section. The remaining entries route to
Run reconciliation, QA report numbering and worktree cleanup.

**Exa MCP.** Consultation was attempted, and no Exa MCP tool was available in
this session. No external source was read for this Spec, so no external
validation is claimed.

**Local measurement.** The counts in this PRD were measured over the twenty most
recent Specs by reading their authored artifacts. They are observations of this
repository, not published evidence.

## Open Questions

- May a declared promise name a later Spec as its owner instead of a Task in
  this Spec? Default until the maintainer decides otherwise: no. A promise
  either names a Task here or is removed from the declaration with its reason,
  so that a Spec cannot promise what it has not scheduled.
