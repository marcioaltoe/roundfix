---
spec: 0258-authoring-that-proves-what-it-claims
prd: _prd.md
created: 2026-10-12
---

# Authoring that proves what it claims — Technical Spec

## Executive Summary

Two new detector files in the Spec Consistency Check read four structural
proofs, and the two authoring skills say how to write them. The TechSpec-stage
detector reads the Absence Evidence section, the ADR Alternatives section and
the Coverage Cross-Read line. The Tasks-stage detector reads each pending
non-QA Task's Consumer Declaration and compares the suites it names with the
Verification commands of that Task or of its Integration Checkpoint. Every
finding is a gap whose fix names the section or line and the skill rule. A
docs Task updates the command reference, the user guide, the Behavior Surface
Record and the glossary.

The trade-off accepted is ADR-0265's decision to hold every active Spec with no
horizon. An authoring prototype at `7a8b7f70` measured the cost: six fixture
TechSpecs, two fixture Task files and two command test files change, and
nothing else in `internal/speccheck`, `internal/cli`, `internal/daemon`,
`internal/delivery`, `internal/spec`, `internal/judge`, `internal/agent` or
`skills`, with and without the `docscontract` and `repocontract` tags. In
exchange, a Spec that has not run yet cannot skip the proofs, and the check
stays structural (ADR-0093: "it never judges whether a decision is correct").

## Project Constraints

- Identifier strategy: not applicable. No identifier scheme changes; the four
  codes follow the existing `SC-` naming. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable. No request, credential or network
  call is added, and the detectors never execute an authored command. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0265 governs every decision below.
  ADR-0093: "it never judges whether a decision is correct". ADR-0094: "is
  skipped, not failed, when one of them is absent". ADR-0156: "with the reason
  none applies". ADR-0168 and ADR-0183 chose horizons this Spec does not take.
  ADR-0176: "The citation walk now reads only the authored projection".
  ADR-0148: "a command that exits zero before the work exists proves nothing
  about the work". ADR-0252: "the selective gate runs the contracts of the
  affected packages plus the cheap global ones", and ADR-0253 runs every
  contract only on main and before a release. ADR-0184 asks for Surface
  Transcripts, declared `None.` below. ADR-0259: "A test runs in parallel
  unless it names why it cannot". ADR-0257: "the Daemon runs the
  Task's declared Verification and the repository Verification at settlement".
  ADR-0258: "Each rerun class becomes an authoring rule". ADR-0244: "A Spec now
  carries a Glossary Declaration". ADR-0256: "Spec authoring requires a skills
  Task for a change a skill describes". Source: `docs/agents/domain.md`.
- Tooling authority: applicable. Express maintainer authorization: "considere
  autorizado a ajustar todas as skills se necessário", "Autorizar os dois" and
  the 2026-10-11 cycle grant "Concedo". Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.
  Spec-contained authorization record:
  `docs/specs/0258-authoring-that-proves-what-it-claims/_authorization.md`.
  Bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/write-tasks/SKILL.md`,
  `.agents/skills/write-tasks/references/task-template.md`,
  `.agents/skills/write-techspec/SKILL.md`,
  `.agents/skills/write-techspec/references/techspec-template.md`,
  `internal/docscontract/testdata/corpus-golden.json`,
  `internal/spec/archive_layout_characterization_test.go`,
  `internal/speccheck/coherence.go`, `skills/baseline_skill_contract_test.go`,
  `skills/roundfix/SKILL.md`, `skills/write-tasks/SKILL.md`,
  `skills/write-techspec/SKILL.md`,
  `skills/write-techspec/references/techspec-template.md`.

## Current behavior

- `write-techspec` tells the author to explore and to "Prefer existing seams to
  new ones", but asks for no evidence that a new component is absent.
- `SC-ADR-UNLISTED` and `SC-ADR-RELATED` hold the Active ADR obligations row to
  the ADRs a Spec cites. Nothing asks how the design relates to what a listed
  ADR rejected; the anti-pattern list names only contradicting a decision.
- `SC-COVERAGE-UNMAPPED` reports a PRD user story, Core Feature or Success
  Metric with no Coverage Map entry. Nothing records that the entry was read
  against the PRD, so a mapped promise delivered in another shape passes.
- `write-tasks` says "Do not add a whole-package suite command to each Task for
  regression coverage", and the Daemon's settlement gate is selective, so a
  consumer in another directory runs first at the QA gate.
- Part 1 of the source, the exit status as the verdict, is covered:
  `SC-VERIFY-INVERTED-EXIT` exists and `write-tasks` says "A Verification
  command passes only by exiting zero".

## System Architecture

Everything lives in `internal/speccheck`, beside the detectors it resembles.
The TechSpec-stage detector runs from two call sites that already read the PRD
and the TechSpec: the full sweep in `detectCitationCoverageAndReferences` and
the TechSpec stage in `detectTechSpecCoverage`. The Tasks-stage detector runs
from the full sweep where the Task Graph is loaded. `stagedDetectors` records
the stage of each code, so `addStageSkips` lists them as skipped at earlier
stages. No new package, command, flag or output format.

## Implementation Design

### Interfaces

```go
const (
	CodeAbsenceUnproven          = "SC-ABSENCE-UNPROVEN"
	CodeADRAlternativesUnrelated = "SC-ADR-ALTERNATIVES-UNRELATED"
	CodeCoverageUnread           = "SC-COVERAGE-UNREAD"
	CodeConsumersUnproven        = "SC-CONSUMERS-UNPROVEN"
)

// detectTechSpecProof reads the Absence Evidence and ADR Alternatives
// sections and the Coverage Cross-Read line of one TechSpec.
func detectTechSpecProof(result *Result, prdContent, techSpecContent []byte, prdDisplay, techSpecDisplay string)

// detectConsumerDeclarations reads the Consumer Declaration of every non-QA
// Task whose status is not completed.
func detectConsumerDeclarations(result *Result, repoRoot string, graph *spec.Graph) error
```

### Data Models

The grammar each detector reads, outside fenced blocks:

- Absence Evidence: a `## Absence Evidence` section whose first non-blank line
  is `None.` with a reason in the same paragraph, or one or more numbered items
  `N. <component> — command: <inline code>; result: <text>`. An item counts
  when `command:` is followed by a non-empty inline code span and `result:` by
  non-blank text.
- ADR Alternatives: the listed set is every decision number in the Active ADR
  obligations row of the PRD and of the TechSpec, read with
  `obligationCitationNumbers`. A `## ADR Alternatives` section holds bullets
  `- ADR-NNNN: <text>` with non-blank text; an empty listed set needs no
  section.
- Coverage Cross-Read: inside `## Coverage Map`, a line starting
  `Deviations from the PRD:`. Its remainder is `None.`, or it is followed by
  bullets, each naming at least one PRD unit (User Story, Core Feature or
  Success Metric with its number) as `addDeclaredReferences` reads it.
- Consumer Declaration: a `## Consumers` section in the Task file. Its first
  non-blank line is `None.` with a reason, or it holds bullets of two kinds:
  `- <contract> — consumers: <inline code>[, <inline code>]` and
  `- checkpoint: task_NN`.

### API Contracts

1. API Contract: `SC-ABSENCE-UNPROVEN` is a gap from the TechSpec stage. It is
   reported once when the TechSpec has no Absence Evidence declaration, and
   once for each item without its command or its result. Summary forms:
   `<techspec> has no Absence Evidence declaration` and
   `<techspec> Absence Evidence item <N> has no command or no result`. The fix
   names `## Absence Evidence` and the write-techspec rule "Absence evidence".
2. API Contract: `SC-ADR-ALTERNATIVES-UNRELATED` is a gap from the TechSpec
   stage, reported once for each listed ADR without its line:
   `<artifact> lists ADR-<NNNN> under Active ADR obligations, but <techspec>
   has no ADR Alternatives line for it`. It locates the row and the section.
   The fix names `## ADR Alternatives` and the write-techspec rule "Rejected
   alternatives".
3. API Contract: `SC-COVERAGE-UNREAD` is a gap from the TechSpec stage,
   reported when the Coverage Map has no `Deviations from the PRD:` line, when
   the line has neither `None.` nor a following bullet, and once for each
   deviation bullet that names no PRD unit. The fix names the line and the
   write-techspec rule "Coverage cross-read".
4. API Contract: `SC-CONSUMERS-UNPROVEN` is a gap from the Tasks stage,
   reported for a pending non-QA Task with no Consumer Declaration, for a
   malformed entry, for each suite that occurs in no Verification command of
   the Task or of a valid checkpoint, and for a checkpoint that is missing, is
   the QA Task, does not transitively need the Task, or is not transitively
   needed by the QA Task. The fix names `## Consumers` and the write-tasks rule
   "Consumers of a shared contract".
5. API Contract: the four codes join the active-corpus code set at 0 in the
   corpus golden and in the characterization that pins it, and the command
   reference and the user guide list them. The text and JSON report formats,
   exit codes and flags do not change.

### Surface Transcripts

None. No command, flag, output line format or exit code changes; the new codes
render through the existing finding lines, and the detector tests assert them.

### Skill and document texts

`write-techspec` gains a section after `## Concrete contracts`:

```markdown
## Proof of what the Spec claims

A TechSpec records three proofs in named places, and `roundfix spec check`
reports a missing place as a gap on every active Spec. The check reads only
structure; whether a proof is right stays with the reader.

- Absence evidence. For each component the TechSpec declares new, record the
  search that proved it absent under `## Absence Evidence`, as a numbered item:
  the component, `command:` with the search in backticks, and `result:` with
  what it found. Search the files the TechSpec cites for the behavior, not only
  for the name. When nothing is new, write `None.` followed by the reason.
  `SC-ABSENCE-UNPROVEN` reports a missing declaration or an item without its
  command or result.
- Rejected alternatives. Under `## ADR Alternatives`, write one line
  `- ADR-NNNN: <relation>` for each ADR the Active ADR obligations rows list:
  how the design relates to the alternatives that ADR rejected. When the design
  takes a rejected alternative, say so and supersede or amend the ADR.
  `SC-ADR-ALTERNATIVES-UNRELATED` reports a listed ADR without its line.
- Coverage cross-read. Before `write-tasks`, read each PRD user story, Core
  Feature and Success Metric against the component the Coverage Map names. Add
  the line `Deviations from the PRD:` to the Coverage Map, followed by `None.`
  or by one bullet per PRD unit the TechSpec delivers in another shape, naming
  the unit and the reason. `SC-COVERAGE-UNREAD` reports a missing line or a
  deviation that names no PRD unit.
```

The TechSpec template gains `## Absence Evidence` and `## ADR Alternatives`
sections after `## Coverage Map`, and its Coverage Map comment adds the line
`Deviations from the PRD: None.`.

`write-tasks` gains a section after `## Surface Transcripts in Tasks`:

```markdown
## Consumers of a shared contract

A Shared Contract is a route, a public signature, a persisted format or a
report that code outside the Task's slice consumes. Every non-QA Task carries
a `## Consumers` section, which `roundfix spec check` reads in every Task that
has not completed:

- `- <contract> — consumers: <suite>, <suite>` names the suites that consume
  each Shared Contract the Task changes, each in backticks exactly as its
  Verification command writes it, such as a package path. Each suite occurs in
  one of the Task's Verification commands, in the same command as a check only
  the Task's work satisfies, so the command still fails before the work.
- `- checkpoint: task_NN` names an Integration Checkpoint instead: a non-QA
  Task that needs this Task, that the QA gate needs, and whose Verification
  runs the suites. Use it when a later Task adapts the consumers.
- `None.` followed by the reason states that the Task changes no Shared
  Contract.

`SC-CONSUMERS-UNPROVEN` reports a missing section, a malformed entry, a suite
no Verification runs, or a checkpoint outside that position.
```

The `write-tasks` bullet that forbids a whole-package suite for regression
coverage gains: "A suite that consumes a Shared Contract the Task changes is
the exception; see Consumers of a shared contract." Process step 1 gains: "Read
the TechSpec's Coverage Cross-Read first; a Coverage Map without its
`Deviations from the PRD:` line goes back to `write-techspec`." The Task
template gains a `## Consumers` section before `## Verification`.

The command reference and the user guide list the four codes and add this
paragraph: "These four gaps hold every active Spec, whatever its age. Before
an active Spec's next Run, add the sections and lines each finding names; the
fix of every finding names the section or line and the skill rule that
requires it."

### Invariants

```text
1. A detector reads structure only: presence, shape, a named unit, a literal
   suite string, and graph reachability. It never executes a command.
2. Each finding is a gap, so it fails only under --strict and the gates that
   promote gaps.
3. A detector whose artifact is absent is skipped and listed as skipped.
4. A completed Task and the QA Task are never read for a Consumer Declaration.
5. Every active Spec is held; no commit horizon applies.
6. Every fix names its section or line and its skill rule.
7. Fenced blocks are not declarations.
```

## Coverage Map

Deviations from the PRD: None.

- User Story 1 → Absence Evidence detector, write-techspec "Absence evidence"
- User Story 2 → ADR Alternatives detector, write-techspec "Rejected alternatives"
- User Story 3 → Coverage Cross-Read detector, write-techspec "Coverage cross-read", write-tasks Process step 1
- User Story 4 → Consumer Declaration detector, write-tasks "Consumers of a shared contract"
- User Story 5 → finding fixes, command reference, user guide
- Core Feature 1 → Absence Evidence detector and skill rule
- Core Feature 2 → ADR Alternatives detector and skill rule
- Core Feature 3 → Coverage Cross-Read detector and skill rule
- Core Feature 4 → Consumer Declaration detector and skill rule
- Core Feature 5 → no horizon, fixes, command reference, user guide
- Core Feature 6 → glossary entries
- Success Metric 1 → detector tests and the QA fixture run
- Success Metric 2 → skill texts, version record, mirrors
- Success Metric 3 → fix texts, command reference, user guide
- Success Metric 4 → corpus golden and characterization at 0

## Absence Evidence

1. TechSpec-stage proof detector — command: `grep -rnE 'Absence Evidence|ADR Alternatives|Deviations from the PRD' internal .agents/skills`; result: no match at `7a8b7f70`; the nearest detector, `SC-COVERAGE-UNMAPPED`, reads only Coverage Map entries.
2. Consumer Declaration detector — command: `grep -rnE '## Consumers|SC-CONSUMERS|Consumer Declaration' internal .agents/skills`; result: no match at `7a8b7f70`; no detector reads a Task section beside Verification, Requirements and Context.
3. Authoring skill rules — command: `grep -rniE 'absence evidence|rejected alternative|consuming suite|integration checkpoint' .agents/skills/write-techspec .agents/skills/write-tasks`; result: no match at `7a8b7f70`; the closest text is the anti-pattern against contradicting an ADR and the rule against whole-package suites.
4. Glossary terms — command: `grep -nE '^[*][*](Absence Evidence|ADR Alternatives|Coverage Cross-Read|Consumer Declaration|Shared Contract|Integration Checkpoint)[*][*]' CONTEXT.md`; result: no match at `7a8b7f70`.

## ADR Alternatives

- ADR-0093: adopts none of its rejected alternatives; every detector reads written structure and never infers from prose.
- ADR-0094: its rejected alternative was demanding the full artifact set; the detectors still skip an absent TechSpec or Task Graph.
- ADR-0148: no alternative it rejected is taken; a consuming suite shares its command with a work-dependent check, so no second prober is needed.
- ADR-0156: of its rejected options, inferring coverage from prose stays rejected; the sections reuse its `None.` with a reason form.
- ADR-0168: its horizon is not repealed for `SC-ADR-RELATED`; the new gaps take no horizon by the maintainer's decision, recorded in ADR-0265.
- ADR-0176: no rejected alternative is taken; a Consumer Declaration is authored text, inside the projection it reads.
- ADR-0183: its horizon stays for receipts; ADR-0265 records why the new gaps hold every active Spec instead.
- ADR-0184: no rejected alternative is taken; the Surface Transcripts section declares `None.` with its reason.
- ADR-0244: no rejected alternative is taken; the six terms follow its declaration contract.
- ADR-0252: no rejected alternative is taken; Task Verification names consumers the selective gate does not select.
- ADR-0253: no rejected alternative is taken; the Full Contract Run stays on main and before a release, and Task Verification names consumers earlier.
- ADR-0259: no rejected alternative is taken; every new test is parallel and none needs a sequential reason.
- ADR-0271: no rejected alternative is taken; the detector tests read files and give no child, goroutine or Git a real-time allowance.
- ADR-0256: its rejected model-judging alternative stays rejected; the skills Task updates the command reference with the user guide.
- ADR-0257: no rejected alternative is taken; agents still run focused tests and the Daemon runs the declared Verification.
- ADR-0258: no rejected alternative is taken; this Spec adds skill rules and checks, and no Baseline clause.
- ADR-0265: this Spec's own decision; it rejects a horizon, content judging, re-running searches and relying on the selective gate alone.

## Integration Points

None outside the repository. The Delivery Plan, queue revalidation and the QA
gate already promote gaps through `spec check --strict`, so they refuse an
active Spec the new detectors report.

## Testing Approach

New tests, each calling `t.Parallel()` first:

- `internal/speccheck/authoring_proof_test.go`:
  `TestAbsenceEvidenceGapFollowsTheSection`,
  `TestADRAlternativesGapNamesEachListedADR`,
  `TestCoverageCrossReadGapFollowsTheDeviationsLine`,
  `TestTechSpecProofRunsAtTheTechSpecStageAndInTheFullSweep` and
  `TestTechSpecProofFixesNameTheSectionAndTheSkillRule`. Each builds a Spec in
  a temporary directory, proves the gap with the section or line missing, and
  proves no finding of that code once it is present. The stage test proves the
  skip without a TechSpec and at the PRD stage, and that `--strict` turns the
  gap into a failure.
- `internal/speccheck/consumers_test.go`:
  `TestConsumersGapFollowsTheSection`,
  `TestConsumersGapNamesASuiteNoVerificationRuns`,
  `TestConsumersCheckpointSitsBetweenTheTaskAndTheGate`,
  `TestConsumersSkipCompletedAndQATasks` and
  `TestConsumersRunAtTheTasksStage`.
- `skills/baseline_skill_contract_test.go`:
  `TestAuthoringSkillsStateTheProofRules` reads both skills and both
  templates.

Declared breaks, measured by the prototype: the techspecs of the `clean`,
`vocabulary-missing`, `citation-dirty`, `coverage-unmapped`,
`coverage-untasked` and `reference-unresolved` fixtures under
`internal/speccheck/testdata/repo/docs/specs/` gain the three proofs, which
keeps `TestCheckCleanFixture`, the three mechanical-stage tests,
`TestCheckCitationCoverageErrorLocations` and the CLI spec check tests at
their current expectations. The `task_01.md` of `coverage-untasked` and
`reference-unresolved`, the Task the spec check Verification workspace writes
in `internal/cli/spec_check_test.go`, and the Task
`newDeliverPlanWorkspace` writes in `internal/cli/deliver_plan_test.go` gain a
`None.` Consumer Declaration. The corpus golden, the characterization that
pins it and `corpusFindingCodes` gain the codes at 0 and one sentence of the
update text each. The Behavior Surface Record is re-recorded with
`go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
after the user guide changes.

## Build Order

1. TechSpec-stage proof detector, its tests, the six fixture techspecs and the
   three corpus pins.
2. Consumer detector, its tests, the two fixture Tasks, the two CLI test
   writers and the three corpus pins (depends on: 1).
3. Authoring skill texts, both templates, mirrors, version record and the skill
   contract test.
4. Command reference, user guide, Behavior Surface Record and glossary
   (depends on: 2, 3).
5. QA gate (depends on: 4).

This Spec rebases after any sibling Spec of the cycle that changes
`write-techspec`, `write-tasks` or the owned skill version record, and the
operator queues it after the other active Specs (see `_prd.md` →
Prerequisites).

## Risks & Considerations

- An active Spec queued before this one merges is refused under `--strict`
  until its author adds the sections. The adopter notes and every fix name
  them.
- An Integration Checkpoint whose Verification only reruns consumers passes on
  the unchanged tree, so the vacuity probe refuses it; the skill tells authors
  to use one only where a later Task adapts the consumers.
- A literal suite match is coarse: `./internal/cli` also matches
  `./internal/client`. The check proves presence, and review judges the rest.
- Sibling Specs of the cycle also edit `write-tasks`, `write-techspec` and the
  owned skill version record; the record command reruns after a rebase.

## Glossary

None.

## Decisions

- Two detector files beside the existing detectors; no new package. See
  ADR-0265.
- One code per proof, each a gap with a fix naming its section and rule. See
  ADR-0265.
- The Coverage Cross-Read is one line in the existing Coverage Map, so the
  presence check `SC-COVERAGE-UNMAPPED` stays unchanged.
- The Consumer Declaration lives in each Task, because the Task owns its
  Verification.
