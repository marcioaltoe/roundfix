---
spec: 0107-the-authoring-rules-the-guides-do-not-carry
status: active
created: 2026-08-12
surfaces: [docs, backend]
---

# The authoring rules the guides do not carry

Three fleet repositories independently derived the same class of rule from
measured failures, and every one of them belongs to a guide Roundfix owns — so it
cannot be written where it was learned without being overwritten by the next
update. A Verification that is not hermetic, a requirement that describes data
volume instead of a property, a test that mixes ports with persistence, a contract
changed without its record updated, a commit scope that drags two other Specs into
a Pull Request: each cost a measured Run or a measured round. Two operational
responsibilities have no owner at all — who prepares the environment, and what to
do when the same Work Item fails identically twice. And the mechanism that decides
whether a Spec delivers or freezes is undocumented: where it lives, its entry
format, and what its satisfied-by field means all had to be inferred.

## Project Constraints

- Identifier strategy: applicable — Unreachable Acceptance, Verification, Wave,
  and Work Item are glossary terms whose obligations these clauses state, and a
  clause that names a concept the glossary does not carry is the defect this Spec
  exists to stop repeating. The closing node checks whether the work introduced or
  changed a term. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no network boundary, credential, or
  request is created or read. The work is guidance content rendered from catalog
  modules. Source: `docs/agents/specific-repository.md`.
- Active ADR obligations: applicable — ADR-0093 checks Spec consistency by
  citation rather than inference, which is the bar each new clause must meet if it
  is to be enforced rather than merely stated; ADR-0117 places a check with the
  stage that can produce its defect, which decides where each clause's enforcement
  belongs. ADR-0094 makes the consistency check artifact-presence aware, so a
  clause enforced by a detector must skip where its artifact is absent. ADR-0104
  makes a Spec accept on evidence it did not author, which is what these clauses'
  own evidence is. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — five Baseline catalog modules and the guides
  they render are edited. Express maintainer authorization:
  `docs/workflow/authorizations/2026-08-12-the-authoring-and-baseline-corrections.md`,
  granted 2026-08-12. Bounded files:
  `internal/baseline/assets/modules/autonomous-work.json`,
  `internal/baseline/assets/modules/context-workflow.json`,
  `internal/baseline/assets/modules/spec-workflow.json`,
  `internal/baseline/assets/modules/secondbrain.json`,
  `internal/baseline/assets/modules/core.json`,
  `docs/agents/autonomous-work.md`, `docs/agents/docs-layout.md`,
  `docs/agents/spec-routing.md`, `docs/agents/agent-instructions.md`,
  `docs/agents/secondbrain.md`. Source: `docs/agents/agent-instructions.md`.

## Goals

1. A rule learned by measurement in one repository governs every repository that
   adopts the Baseline.
2. Every clause added names the check that decides it, or states that none does.
3. The mechanism that decides whether a Spec delivers or freezes is documented.
4. Work observed about another project is captured before the session ends.

## User Stories

1. As a Supervisor authoring in any adopting repository, I want the hermeticity,
   property-not-volume, and test-boundary rules stated in my guides, so that I do
   not rediscover them by losing a Run.
2. As a Supervisor whose Spec touches an already-recorded contract, I want the
   rule that the record updates in the same slice, so that a mid-Spec contract
   change does not recur as a finding.
3. As a Supervisor authoring a Spec that freezes on unreachable acceptance, I want
   the section documented, so that I do not infer its format from a citation.
4. As a Supervisor running autonomously, I want an owner for environment
   preparation and a rule for an identically repeated failure, so that neither is
   improvised.
5. As a session that observed something about another project, I want a clause
   that triggers capture, so that the observation survives the session instead of
   competing with the work and losing.

## Core Features

1. **The measured authoring rules are carried.** A Verification is hermetic; a
   requirement describes the property rather than the magnitude of the data; a use
   case test is born against ports and doubles while a persistence proof is born in
   infrastructure; a Task changing an already-recorded contract updates the record
   in the same slice; commit scope is per Spec.
2. **The tooling-authority chronology is stated.** The commit that authorizes a
   tooling mutation does not contain it — one line that resolves a measured
   two-round gate failure and a history rewrite.
3. **Unreachable acceptance is documented.** Where the section lives, its entry
   format, and what its satisfied-by field means.
4. **Wave independence looks at files.** The task-authoring guidance states that
   Tasks in one wave do not share an edit target, or declare the dependency.
5. **The autonomous guide gains its two missing ends.** An owner for environment
   preparation, and the rule that an identically repeated Work Item failure is
   answered by amending the task file rather than by re-invoking.
6. **Tooling authority is asked by class.** The technical-spec stage walks the
   predictable classes and asks for the set once, naming the classes a Spec will
   not touch.
7. **The knowledge guide gains a production trigger.** A session observing a
   defect whose owner is another project captures it in that project's namespace
   before the session ends, and proposed changes to baseline-owned guides are
   captured rather than edited locally.

## User Experience

Each clause reads in the same register as the guidance around it, and a rendered
guide that carries it says the same thing in every adopting repository.

## Non-Goals / Out of Scope

- Building the detectors that would enforce these clauses mechanically. This Spec
  states each clause and names its enforcement owner; implementing a new detector
  belongs to the Spec that owns that check.
- Changing this repository's own guides beyond what the rendered modules produce.
- Re-litigating any clause the guides already carry.
- Moving the canonical method out of the guides, which is a separate Spec.

## Success Metrics

- Each of the five measured authoring rules appears in a rendered guide of a
  repository adopting the Baseline. Source: three repositories this Spec did not
  build, whose 2026-08-07/08 sessions derived them from measured failures.
- The unreachable-acceptance documentation answers the three questions a
  Supervisor had to infer.
- A session observing a cross-project defect produces a capture, measured against
  the session that delivered three Specs and captured nothing until asked.
- No clause ships without naming the check that decides it or stating that none
  does.

## Decisions

- Every clause names its enforcement owner, because a finding filed on 2026-08-06
  predicted a later Spec's blocking defect exactly and nothing stopped it until
  the rule became a detector. Prose alone has a measured failure rate here.

## Open Questions

- Whether the environment-preparation owner is the Supervisor, the repository's
  own runbook, or a preflight check. The repository runbook is the default until
  answered, since it is where the resolving procedure already lives.
- Whether asking tooling authority by class belongs in the technical-spec stage or
  the product stage. The technical-spec stage is the default, since the classes are
  architectural rather than product.
