---
spec: 0258-authoring-that-proves-what-it-claims
status: active
created: 2026-10-12
surfaces: [backend, docs]
---

# Authoring that proves what it claims

A Spec can pass `write-techspec`, `write-tasks` and `roundfix spec check` while
it claims things nobody proved. [The adopted Backlog Entry](references/2026-10-12-spec-authoring-does-not-prove-what-it-claims.md)
records four such gaps, measured in another repository that uses these skills.
Two Tasks with green Verification left the repository red, because each moved a
route that a suite in another directory consumed. A Spec ordered a quarantine
that already existed, with eight tests, in the file its TechSpec cited, and the
implementing Agent lowered its threshold and deleted three tests to reconcile
the contradiction. A Spec raised a timeout to the exact value an accepted ADR
had rejected. A TechSpec mapped a PRD promise of per-request records in its
Coverage Map while aggregating them on purpose. After this Spec, an author
records the proof of each of those claims in a named section or line, the
authoring skills say how, and `roundfix spec check --strict` reports a missing
section or line structurally on every active Spec, without judging its content.
The maintainer chose "Skills + gap no spec check" on 2026-10-11, and on
2026-10-12 decided that the new gaps hold every active Spec with no exemption by
date.

## Prerequisites

The new gaps hold every active Spec, and the active-corpus contract tests keep
every finding count at zero. The operator therefore queues this Spec after the
other active Specs of the same cycle, or completes the new sections in any Spec
that is still active when this Spec's first Task runs. This Spec never edits
another Spec's folder. Its own TechSpec and pending Tasks already carry every
section the new detectors read, so the Spec holds itself to its rules from its
first Task. No Verification of this Spec reads another Spec's artifacts.

## Project Constraints

- Identifier strategy: not applicable. No identifier scheme changes; the four
  new finding codes follow the existing `SC-` naming. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable. No request, credential or network
  call is added. Tests use temporary repositories and the repository's own
  fixtures. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0265 (this Spec) decides the four
  proofs, the structural check and the absence of a horizon. ADR-0093: "it
  never judges whether a decision is correct", so every new detector reads
  structure only. ADR-0094: "is skipped, not failed, when one of them is
  absent", so the TechSpec detectors skip a Spec without a TechSpec and the
  Consumer detector skips a Spec without a Task Graph. ADR-0156: "with the
  reason none applies", the declaration form each new section reuses. ADR-0168
  and ADR-0183 gave the related-ADR and receipt gaps a horizon; this Spec
  deliberately takes none, by the maintainer's decision. ADR-0176: "The
  citation walk now reads only the authored projection", which a Consumer
  Declaration belongs to. ADR-0148: "a command that exits zero before the work
  exists proves nothing about the work", so a consuming suite shares its
  command with a condition only the work satisfies. ADR-0252: "the selective
  gate runs the contracts of the affected packages plus the cheap global ones",
  which is why a consumer in another directory goes unrun until the QA gate;
  ADR-0253 runs every contract only on main and before a release, which is
  later still. ADR-0184 asks a TechSpec for its Surface Transcripts; this Spec
  declares none, because no command output format changes. ADR-0259: "A test
  runs in parallel unless it names why it cannot", so every new test calls
  `t.Parallel()` first.
  ADR-0257: "the Daemon runs the Task's declared Verification and the
  repository Verification at settlement", so the Task's own Verification is
  where its consumers must run. ADR-0258: "Each rerun class becomes an
  authoring rule", the pattern this Spec extends with four more classes.
  ADR-0244: "A Spec now carries a Glossary Declaration", which this Spec makes
  for six terms. ADR-0256: "Spec authoring requires a skills Task for a change
  a skill describes", so the skills Task updates the command reference with the
  user guide. Source: `docs/agents/domain.md`.
- Tooling authority: applicable. Express maintainer authorization covers the
  owned skills and the Governed Paths each cycle Spec declares: "considere
  autorizado a ajustar todas as skills se necessário", "Autorizar os dois",
  and the 2026-10-11 cycle grant "Concedo". No Makefile, CI workflow, lint,
  formatter, test-runner configuration or `go.mod` change. Source:
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

## Goals

- A TechSpec records the search that proved each component it declares new is
  absent, as a command and its result.
- A TechSpec states, for each ADR it lists, how the Spec relates to the
  alternatives that ADR rejected.
- A TechSpec's Coverage Map records the cross-read of the PRD and names every
  promise it delivers in another shape.
- A Task that changes a contract other code consumes runs the consuming suites
  in its Verification, or the graph runs them in a checkpoint before the QA
  gate.
- `roundfix spec check --strict` refuses an active Spec that lacks any of
  those sections or lines, and its fix names the section and the skill rule.

## User Stories

1. As a Task agent, I want the TechSpec to show the search that proved a
   component absent, so that I never weaken working code to make a Spec that
   orders it again come true.
2. As a Spec author, I want one line per listed ADR on its rejected
   alternatives, so that I notice when my design is an option a decision
   already refused.
3. As the operator, I want the PRD-to-TechSpec cross-read recorded before
   decomposition, so that a promise the TechSpec delivers in another shape is
   visible before any Task runs.
4. As the operator, I want a contract-changing Task to run the suites that
   consume the contract, so that a green Task never leaves the repository red
   for the QA gate to find.
5. As an adopter, I want each finding to name the missing section and the rule
   that requires it, so that I can complete an active Spec before its next Run
   without reading the skill first.

## Core Features

1. **Absence evidence.** `write-techspec` requires an **Absence Evidence**
   section: one numbered item per component the TechSpec declares new, with
   the search command and its result, or `None.` with the reason no component
   is new. `spec check` reports `SC-ABSENCE-UNPROVEN` for a missing
   declaration or an item without its command or result (ADR-0265).
2. **ADR alternatives.** `write-techspec` requires an **ADR Alternatives**
   section with one line for each ADR the Active ADR obligations rows list.
   `spec check` reports `SC-ADR-ALTERNATIVES-UNRELATED` for each listed ADR
   without its line (ADR-0265).
3. **Coverage cross-read.** `write-techspec` requires the **Coverage
   Cross-Read** in the Coverage Map, a line that says `None.` or lists each PRD
   unit the TechSpec delivers in another shape. `spec check` reports
   `SC-COVERAGE-UNREAD` for a missing line or a deviation that names no PRD
   unit. The existing check that every PRD unit has a Coverage Map entry stays
   as it is (ADR-0265).
4. **Consumer declaration.** `write-tasks` requires a **Consumer
   Declaration** in every non-QA Task: the suites that consume each **Shared
   Contract** the Task changes, run by its own Verification or by an
   **Integration Checkpoint** Task between it and the QA gate, or `None.` with
   the reason. `spec check` reports `SC-CONSUMERS-UNPROVEN` for a missing
   declaration, a malformed entry, a suite no Verification runs or a checkpoint
   outside that position. Completed Tasks are not read (ADR-0265).
5. **Every active Spec.** The four gaps hold every active Spec with no horizon.
   Each finding's fix names the missing section or line and the skill rule.
   The command reference and the user guide list the codes and the steps an
   adopter takes before the next Run (ADR-0265).
6. **Glossary.** The glossary gains the six terms of the Glossary section.

## Non-Goals / Out of Scope

- Judging whether a search was sufficient, a relation true, a deviation
  justified or a consumer list complete. Those stay with review and QA.
- Running Absence Evidence commands at check time, or any new network access.
- Part 1 of the source entry, "the exit status is the verdict". It is already
  covered: `SC-VERIFY-INVERTED-EXIT` refuses a status hidden by a pipe, and
  `write-tasks` already says "A Verification command passes only by exiting
  zero".
- The presence of a Coverage Map entry for each PRD unit, which
  `SC-COVERAGE-UNMAPPED` already reports (ADR-0156).
- Editing another Spec's folder, archived Specs, Archive Records or anything
  under `docs/history/`.
- Baseline modules, guides under `docs/agents/`, the Makefile, CI workflows or
  any dependency.
- The advisory judge; this Spec adds no judge question.

## Success Metrics

1. Success Metric: on fixture Specs in temporary repositories, each of
   `SC-ABSENCE-UNPROVEN`, `SC-ADR-ALTERNATIVES-UNRELATED`, `SC-COVERAGE-UNREAD`
   and `SC-CONSUMERS-UNPROVEN` is reported as a gap when its section or line
   is missing and is absent when the section or line is present, and
   `spec check --strict` exits 1 only in the first case.
2. Success Metric: `write-techspec`, its TechSpec template, `write-tasks` and
   its Task template state the four rules, each owned skill records a raised
   version, and every mirror matches its canonical copy.
3. Success Metric: every new finding's fix names its section or line and the
   skill rule, and the command reference and the user guide list the four codes
   and the adopter steps.
4. Success Metric: at the audited head the active Spec corpus, this Spec
   included, reports zero findings for the four codes, and every archived Spec
   stays byte-identical.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- The adopted Backlog Entry, committed in this Spec's references. It is the
  operator's measurement of four authoring gaps in another repository, which
  the skills maintainer confirmed on 2026-10-09 as uncovered by the synced
  skills.
- Ian Robinson, "Consumer-Driven Contracts: A Service Evolution Pattern"
  (martinfowler.com, 2006): "where a proposed change would in fact break one or
  more consumers, the provider will have immediate insight into the issue".
- Elmar Juergens, "Why and How to Control Cloning in Software Artifacts"
  (TU München, 2011): "Every second change to cloned code that was unaware of
  cloning was faulty."
- Olaf Zimmermann, "The Markdown ADR (MADR) Template Explained and Distilled"
  (2022), on the Considered Options section: "This section lists the
  alternatives (or choices, options, candidate solutions) investigated".
- Gotel and Finkelstein, "An Analysis of the Requirements Traceability
  Problem" (ICRE 1994): "Requirements traceability refers to the ability to
  describe and follow the life of a requirement, in both a forwards and
  backwards direction".

## Unreachable Acceptance

- criterion: fewer QA failures caused by an unrun consumer, a rebuilt existing
  component, a rejected alternative or an unread Coverage Map in Specs
  authored after release
  reason: the effect shows only in Specs authored with the released skills,
  which start after this Spec merges
  satisfied-by: the operator reviews the QA reports of the first five Specs
  delivered after this release and records each failure of those four classes

## Adopter notes

The release notes for this change tell adopters what an active Spec must add
before its next Run:

- In `_techspec.md`, a `## Absence Evidence` section: one numbered item per
  component the TechSpec declares new, written as a component, the search
  command and its result, or `None.` with the reason no component is new.
- In `_techspec.md`, a `## ADR Alternatives` section with one
  `- ADR-NNNN: <relation>` line for each ADR the Active ADR obligations rows
  list.
- In the TechSpec's `## Coverage Map`, a `Deviations from the PRD:` line,
  followed by `None.` or by one bullet per PRD unit delivered in another shape.
- In every non-QA Task that has not completed, a `## Consumers` section that
  names the suites consuming each Shared Contract the Task changes, or `None.`
  with the reason.

Run `roundfix spec check <slug> --strict` to list what remains; each finding
names the section or line and the skill rule.

## Glossary

- adds: **Absence Evidence**
- adds: **ADR Alternatives**
- adds: **Coverage Cross-Read**
- adds: **Consumer Declaration**
- adds: **Shared Contract**
- adds: **Integration Checkpoint**

## Decisions

- Each claim's proof is a named section or line, and the check reads only its
  structure; see ADR-0265.
- The gaps hold every active Spec with no horizon, by the maintainer's
  decision of 2026-10-12; see ADR-0265.
- The parts of the source already covered are closed with their evidence under
  Non-Goals instead of rebuilt.

## Research record

- Secondbrain: `wiki/index.md` read first; `qmd query` over the four gaps
  returned the triaged inbox entries the adopted Backlog Entry cites
  (`inbox/roundfix/_triaged/2026-08-08-write-tasks-nao-governa-exit-status-nem-consumidores.md`,
  `inbox/roundfix/_triaged/2026-08-13-autoria-de-spec-nao-verifica-se-o-que-ela-manda-construir-ja-existe.md`,
  `inbox/roundfix/_triaged/2026-08-18-a-autoria-de-spec-nao-presta-contas-de-findings-nem-verifica-o-coverage-map.md`)
  and the mirrored ADRs on citation-only checks and declared promises. They
  confirmed the four gaps and the maintainer's suggested rule wording, and they
  shaped the `None.` with a reason form. The other repository was read
  read-only during authoring to confirm the cited Specs and the existing
  quarantine constant; nothing from it is copied here.
- Exa: the four published sources above. Robinson supports naming consumers in
  the provider's tests; Juergens supports proving absence before building;
  Zimmermann supports recording a decision's rejected options as a first-class
  section; Gotel and Finkelstein support tracing each requirement forward into
  the design. None of them supports judging content mechanically, which agrees
  with ADR-0093.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the section grammar, the detectors,
the measured fixture breaks, the skill texts and the build order.
