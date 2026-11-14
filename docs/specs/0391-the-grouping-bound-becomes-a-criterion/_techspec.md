---
spec: 0391-the-grouping-bound-becomes-a-criterion
status: active
created: 2026-11-14
---

# The grouping bound becomes a criterion — Technical Spec

## Executive Summary

The Baseline clause `clause.spec.sources-02-bound-the-group` and five carriers
around it bound a Spec at four implementation Tasks plus its QA gate and never
say why. This Spec rewords the clause, keeping its id and force, so a Spec is
sized by one context and by Tasks that each earn their verification cycle,
with the two reasons stated and four Tasks kept as a reference, not a cap;
every guarantee the old text gave is kept and pinned by its own phrase. The
`write-idea`, `write-prd` and `write-techspec` skills, the Roundfix `spec`
reference and the `spec` command guide follow, the `write-tasks` skill gains
the matching Task sizing rule, the glossary redefines **Grouping Bound**, and
ADR-0400 records the decision and amends ADR-0208.

The trade-off accepted: the criteria are judgments an author makes, not
structure a check can count, so a Spec that drifts into a second context or a
Task too small for its cycle is caught by the author, the judge's grouping
suggestion and review, never by `spec check`.

## Project Constraints

- Identifier strategy: not applicable — no identifier, record or schema is
  minted; the clause keeps its id. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network
  call is added. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0400 (new) records the decision.
  ADR-0208 is amended by task_02: its "**The bound.**" bullet no longer holds;
  the grouping and extend-before-minting clauses and the one owning Spec per
  source stay. ADR-0209 applies unchanged: the judge's grouping suggestion
  stays advisory and its question is not changed. ADR-0200 and ADR-0201, which
  ADR-0209 cites, do not apply: this Spec changes neither the judge's force
  nor what it sends. ADR-0319, which cites ADR-0209, does not apply: it orders
  the Follow-up Command's items, and this Spec changes no follow-up order. Source:
  `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — the maintainer's answers of 2026-11-14
  quoted in `_authorization.md`, on the cycle BE terms; the Governed Paths are
  the spec-workflow module, the Standard TypeScript profile and its formatter
  golden, and this repository's spec routing guide and setup manifest that
  task_01 changes through their commands, and the skill files task_02 changes.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0391-the-grouping-bound-becomes-a-criterion/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/spec.md`, `.agents/skills/write-idea/SKILL.md`, `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-tasks/SKILL.md`, `.agents/skills/write-techspec/SKILL.md`, `docs/agents/setup-context.json`, `docs/agents/spec-routing.md`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`, `internal/baseline/assets/modules/spec-workflow.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `skills/roundfix/SKILL.md`, `skills/write-idea/SKILL.md`, `skills/write-prd/SKILL.md`, `skills/write-tasks/SKILL.md`, `skills/write-techspec/SKILL.md`.

## Current behavior

Read on 2026-11-14 at `251ed81e`; `_prd.md` → Validation records the runs.

1. `internal/baseline/assets/modules/spec-workflow.json` (module version 24):
   `rule.spec.routing` is at version 10, and
   `clause.spec.sources-02-bound-the-group` (`enforcement: mandatory`, no
   `replaces`) reads "Group sources into one Spec only while the Spec fits four
   implementation Tasks plus its QA gate. When the grouped scope needs more,
   split it into Specs that each fit and give each source exactly one owning
   Spec; never add a source that shares no context with the Spec only to save a
   Spec. When a Spec's Verification or premise depends on another Spec that has
   not merged, name that Spec in the Task Graph manifest's `requires` and state
   the dependency in the TechSpec Build Order." It renders verbatim at line 19
   of `docs/agents/spec-routing.md` and of the Standard TypeScript formatter
   golden of that guide.
2. `internal/baseline/grouped_sources_clauses_test.go` pins that text as a
   literal in `groupingClauses`, used by
   `TestTheGroupingClausesCarryTheirForceAndText`,
   `TestTheGroupingClausesRenderInTheGuides`,
   `TestTheGroupingClausesHaveSourceBaselineRows` and
   `TestAnAdopterRetainsTheGroupingClauses`. No other test pins it.
3. The Source Baseline `baseline.standard-typescript-monorepo-0.0.1` holds a row
   for the clause with its own wording ("MUST group sources into one Spec only
   while the Spec fits four implementation Tasks plus its QA gate, and split a
   larger grouped scope into Specs that each fit."), a byte range and a digest.
   `classifySourceClauseTransition` (`internal/baseline/plan.go`) retains a row
   by id and enforcement and never compares text.
4. The three authoring skills' `## Sources that share a context` sections say
   "Group only while the scope fits four implementation Tasks plus its QA gate;
   past that bound, split the scope and give each source one owning Spec." The
   Roundfix `references/spec.md` says "The grouping bound is four
   implementation Tasks plus its QA gate."; `docs/user-guide/commands/spec.md`
   says "within the bound of four implementation Tasks plus its QA gate".
   `write-tasks` sizes a Task only by "**Sized for one fresh session.**" and
   names no Spec bound. No other guide, skill or user-guide page states it: the
   other "at most four" in the review reference and guide bound Review
   Segments.
5. `CONTEXT.md` **Grouping Bound** reads "The size limit on grouping sources
   into one Spec: four implementation Tasks plus its QA gate." ADR-0208's
   "**The bound.**" bullet says the same.
6. Owned skill versions at the base: `write-idea` 0.0.7, `write-prd` 0.0.18,
   `write-techspec` 0.0.20, `write-tasks` 0.0.32, `roundfix` 0.1.233.

## System Architecture

- `internal/baseline/assets/modules/spec-workflow.json`: the reworded clause,
  its pinning tests, the derived files and this repository's managed refresh
  (task_01).
- The authoring and `write-tasks` skills, the Roundfix `spec` reference, the
  `spec` command guide, `CONTEXT.md` and ADR-0208's notice (task_02).

## Implementation Design

### Interfaces

No Go interface, function or type changes. One test is added to each of
`internal/baseline/grouped_sources_clauses_test.go`
(`TestTheGroupingBoundIsACriterionThatKeepsEveryGuarantee`) and
`skills/acceptance_rules_test.go`
(`TestTheAuthoringSkillsSizeASpecByContextAndVerificationCycle`).

### Data Models

None.

### API Contracts

1. **Baseline clause.** In `rule.spec.routing`, whose `version` rises by one
   from the value the Task finds (10 at `251ed81e`), the clause
   `clause.spec.sources-02-bound-the-group` keeps its id, `enforcement:
   mandatory` and absence of `replaces`; every other clause, rule and module is
   unchanged, and its `guidance` becomes, word for word:

   ```text
   Size a Spec by two criteria, never by a count of its Tasks. One context: the sources a Spec adopts and the work it holds change the same component or contract, share a root cause, or complete one another. Tasks that each earn their verification cycle: every Task pays an Agent turn, its declared Verification and the repository Verification at settlement, so a change too small to justify that cost joins a neighbouring Task of the same context instead of becoming its own Task, and a small source of the same context joins a Spec that is open or being planned instead of becoming a one-Task Spec that pays a whole Run, QA gate, review and release. The criteria answer two costs: a Spec that holds several contexts grows giant, and a Spec of small Tasks spends Runs on long verification and takes much longer than its work needs. Four implementation Tasks plus the QA gate is a reference, not a cap: a Spec that needs more is a signal to check whether it holds more than one context, never a refusal by count. When the scope holds more than one context, split it into Specs that each hold one and give each source exactly one owning Spec; never add a source that shares no context with the Spec only to save a Spec. When a Spec's Verification or premise depends on another Spec that has not merged, name that Spec in the Task Graph manifest's `requires` and state the dependency in the TechSpec Build Order.
   ```

2. **Authoring skills.** In each canonical `.agents/skills/<skill>/SKILL.md` of
   `write-idea`, `write-prd` and `write-techspec`, under `## Sources that share
   a context`, the three lines that begin "unresolved Findings for sources
   that share the context. Group only while the" and end "the scope and give
   each source one owning Spec." become, word for word and wrapped as shown:

   ```text
   unresolved Findings for sources that share the context. The bound on grouping
   is a criterion, not a count: one context, where the sources and the work change
   the same component or contract, share a root cause, or complete one another,
   and Tasks that each earn their verification cycle. A change too small to
   justify a Task's Agent turn, its declared Verification and the repository
   Verification at settlement joins a neighbouring Task of the same context, and a
   small source of the same context joins a Spec that is open or being planned
   rather than becoming a one-Task Spec. Four implementation Tasks plus the QA
   gate is a reference, not a cap: needing more is a signal to check whether the
   Spec holds more than one context. When it holds more than one, split the scope
   and give each source one owning Spec.
   ```

   The judge paragraph that follows ("by adopting the open source within the
   bound") is unchanged: "the bound" now names the criterion above.
3. **Task sizing.** In the canonical `.agents/skills/write-tasks/SKILL.md`,
   under `## Decomposition rules`, right after the line that begins
   "- **Sized for one fresh session.**", one new line, word for word:

   ```text
   - **Earns its verification cycle.** Every Task pays an Agent turn, its declared Verification and the repository Verification at settlement. A change too small to justify that cost, such as one sentence, a version raise or a regenerated record, joins a neighbouring Task of the same context instead of becoming its own Task. It joins only a Task whose slice it belongs to, so the vertical-slice rule holds, and Tasks that declare the same file still run in series under the Wave collision rule. Four implementation Tasks plus the QA gate is a reference, not a cap: a graph that needs more is a signal to check whether the Spec holds more than one context, never a refusal by count.
   ```

   The vertical-slice rule, "Sized for one fresh session", the Wave collision
   text and the Corrective-Task ceiling are unchanged.
4. **Roundfix reference and command guide.** In the canonical
   `.agents/skills/roundfix/references/spec.md`, the two lines "The grouping
   bound is four implementation Tasks plus its QA gate. A" / "suggestion never
   gates." become, word for word:

   ```text
   The grouping bound is a criterion, not a count: one context, and Tasks that
   each earn their verification cycle; four implementation Tasks plus its QA gate
   is a reference, not a cap. A suggestion never gates.
   ```

   In `docs/user-guide/commands/spec.md`, the two lines "source within the bound
   of four implementation Tasks plus its QA gate, or by" / "stating why it stays
   apart." become, word for word:

   ```text
   source within the grouping bound, one context and Tasks that each earn their
   verification cycle, or by stating why it stays apart. Four implementation Tasks
   plus its QA gate is a reference, not a cap.
   ```

   `make skills-sync` copies every canonical skill file to `skills/`.
5. **Glossary.** In `CONTEXT.md`, through `domain-modeling`, the
   **Grouping Bound** entry's definition and `_Avoid_` lines become, word for
   word, with the two links as written:

   ```text
   **Grouping Bound**:
   The limit on grouping sources into one Spec, set by two criteria and not by a count of Tasks: one context, where the sources and the work change the same component or contract, share a root cause, or complete one another, and Tasks that each earn their verification cycle, so a change too small for its own Task joins a neighbouring Task of the same context and a small source joins a Spec of its context that is open or being planned. Four implementation Tasks plus the QA gate is a reference, not a cap. A scope that holds more than one context is split into Specs that each hold one, and each source keeps exactly one owning Spec ([ADR-0208](docs/adr/0208-work-items-that-share-a-context-share-a-spec.md), [ADR-0400](docs/adr/0400-a-spec-is-sized-by-one-context-and-by-tasks-that-earn-their-verification-cycle.md)).
   _Avoid_: Spec size rule, Task cap, four-Task limit
   ```

   The texts passed `TestGlossaryPhrasesAreDefinedTerms` and the other glossary
   contracts on the prototype.
6. **Amendment notice.** In
   `docs/adr/0208-work-items-that-share-a-context-share-a-spec.md`, after the
   paragraph that ends "so the clause writes down a limit in use rather than
   inventing one.", one blank line and this paragraph, word for word, and
   `updated_at` set to `2026-11-14T10:00:00Z`; nothing else changes:

   ```text
   **Amended by [ADR-0400](0400-a-spec-is-sized-by-one-context-and-by-tasks-that-earn-their-verification-cycle.md) (2026-11-14).** The fixed bound of four implementation Tasks plus the QA gate is gone. Grouping stops at one context, and each Task earns its verification cycle: a change too small for its own Task joins a neighbouring Task of the same context, and a small source of the same context joins a Spec that is open or being planned. Four implementation Tasks plus the gate stays as a reference, not a cap. Each source still keeps exactly one owning Spec, a split still gives each source one owner, and a source that shares no context is still never added only to save a Spec.
   ```

### Surface Transcripts

None. No command's standard output, standard error or exit code changes: the
Spec changes Baseline guidance, skill text, one user-guide sentence, the
glossary and ADRs.

### Invariants

1. The clause keeps its id, its `mandatory` force and its absence of
   `replaces`, so an adopter's refresh records it `retained`; the Source
   Baseline corpus, manifest and index are byte-identical to the Delivery Base.
2. Every guarantee of the old clause is in the new text and pinned by its own
   phrase: one owning Spec per source, a split by context, no unrelated source
   added to save a Spec, and `requires` plus the Build Order for an unmerged
   prerequisite Spec.
3. No code, check or judge question counts a Spec's Tasks after this Spec, as
   before it.
4. `clause.spec.sources-01-group-by-shared-context` and
   `clause.context.inbox-03-extend-before-minting` are byte-identical.

### Declared breaks

Measured with a prototype of tasks 01 and 02 on the clone at `251ed81e`.

- Behavior, task_01: every adopter's next managed refresh rewrites the clause
  in `docs/agents/spec-routing.md`; this repository's guide and setup manifest
  change through the same refresh. The derived files that change are the module
  version record, the Standard TypeScript profile and its formatter golden, the
  three catalog files and four plan-characterization goldens task_01 declares.
- Test changed, task_01: the `clause.spec.sources-02-bound-the-group` literal in
  `groupingClauses` of `internal/baseline/grouped_sources_clauses_test.go`
  moves to API Contract 1's text; the four tests that read it keep their names
  and assertions. No recorded test is removed or renamed, so
  `docs/references/coverage-record.json` is not re-recorded.
- Behavior, task_02: the authoring skills, `write-tasks`, the Roundfix `spec`
  reference and the `spec` command guide state the criterion; the versions of
  `write-idea`, `write-prd`, `write-techspec`, `write-tasks` and `roundfix`
  rise, and the guide edit moves the `spec` command guide's fingerprint in
  `docs/references/behavior-surfaces.json`, which task_02 re-records.

## Coverage Map

- Goals 1, 2, 3 and 5, Core Feature 1, Success Metric 1 → API Contract 1,
  `TestTheGroupingBoundIsACriterionThatKeepsEveryGuarantee`, Invariant 1, the
  managed refresh (task_01).
- Goal 4, Core Feature 1, Success Metrics 1 and 4 → API Contract 1, the same
  test, Invariants 2 and 4, `TestTheGroupingClausesCarryTheirForceAndText`,
  `TestAnAdopterRetainsTheGroupingClauses` (task_01).
- Goals 1 to 3, Core Feature 2, Success Metrics 2 and 4 → API Contracts 2 to 4,
  `TestTheAuthoringSkillsSizeASpecByContextAndVerificationCycle` (task_02).
- Core Feature 3, Success Metric 3 → API Contracts 5 and 6, the documentation,
  glossary and skill contracts (task_02).

Deviations from the PRD: None.

## Absence Evidence

1. The new clause wording — command:
   `grep -rn 'never by a count of its Tasks' internal/baseline/assets docs/agents .agents/skills skills CONTEXT.md`;
   result: no match.
2. The new tests — command:
   `grep -rn 'TestTheGroupingBoundIsACriterionThatKeepsEveryGuarantee\|TestTheAuthoringSkillsSizeASpecByContextAndVerificationCycle' internal skills`;
   result: no match.
3. The Task sizing rule and the reworded skill and guide sentences — command:
   `grep -rn 'Earns its verification cycle\|earn their verification cycle\|a reference, not a cap' .agents/skills skills docs/user-guide CONTEXT.md`;
   result: no match.
4. The amendment notice — command: `grep -rln 'ADR-0400' docs/adr CONTEXT.md`;
   result: only
   `docs/adr/0400-a-spec-is-sized-by-one-context-and-by-tasks-that-earn-their-verification-cycle.md`.
5. A Task count enforced in code — command:
   `grep -rniE 'four implementation|grouping bound' internal --include=*.go --include=*.json`;
   result: only the clause in `spec-workflow.json` and its literal in
   `grouped_sources_clauses_test.go`, plus derived catalog files under
   `internal/baseline/testdata/`.

## ADR Alternatives

- ADR-0400: the decision this Spec records; it rejects keeping four with the
  reasons stated, a Backlog Entry that keeps the rule, a strict-check rule that
  counts or measures Tasks, and renaming the term, and the design takes none of
  them.
- ADR-0208: amended; it set aside one Spec per source and a new file for an
  observation an open Backlog Entry already holds, and the design takes
  neither; its "**The bound.**" decision, "The bound is a limit on grouping, not
  a size rule for every Spec", is replaced by the criterion, while its grouping
  and extend-before-minting clauses stay.
- ADR-0209: applies unchanged; it set aside a lexical baseline nearly as good
  overall, and the design takes it neither: the grouping suggestion stays a
  Jev advisory, and its question and threshold are unchanged.
- ADR-0200: not applicable; it decides that the judge never gates, and this
  Spec changes no judgment and adds none, so the design takes none of the
  judgments it set aside: no judgment measures a Task or a Spec.
- ADR-0201: not applicable; it decides what the judge sends and what it
  spends, and this Spec sends nothing new, so the design takes neither a
  boundary kept by convention nor the generic OpenRouter key.
- ADR-0319: not applicable; it decides the follow-up order over four Jev
  answers, and this Spec changes no follow-up order, so the design takes none
  of the orders it set aside.

## Integration Points

None outside the repository. The managed refresh is this repository's own
`roundfix baseline update`; no test or QA row reaches the network, the judge
or another repository.

## Testing Approach

Every Go test reads the embedded catalog, the canonical and embedded skill
files, or a temporary repository; none reaches the network or the real
`~/.roundfix`. Each new test calls `t.Parallel()` first.

- `internal/baseline/grouped_sources_clauses_test.go` gains a
  `groupingBoundPhrases` list of `{name, phrase}` pairs and
  `TestTheGroupingBoundIsACriterionThatKeepsEveryGuarantee`, which reads the
  embedded asset
  `formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`,
  joins its whitespace-separated fields with single spaces, fails when it
  contains "only while the Spec fits four implementation Tasks plus its QA
  gate", and runs one subtest per pair, named by the pair's name, that requires
  the phrase exactly once:
  - `no-count`: "Size a Spec by two criteria, never by a count of its Tasks."
  - `one-context`: "One context: the sources a Spec adopts and the work it
    holds change the same component or contract, share a root cause, or
    complete one another."
  - `earned-cycle`: "Tasks that each earn their verification cycle: every Task
    pays an Agent turn, its declared Verification and the repository
    Verification at settlement,"
  - `small-change-joins-a-task`: "a change too small to justify that cost joins
    a neighbouring Task of the same context instead of becoming its own Task,"
  - `small-source-joins-a-spec`: "a small source of the same context joins a
    Spec that is open or being planned instead of becoming a one-Task Spec that
    pays a whole Run, QA gate, review and release."
  - `two-reasons`: "a Spec that holds several contexts grows giant, and a Spec
    of small Tasks spends Runs on long verification and takes much longer than
    its work needs."
  - `reference-not-cap`: "Four implementation Tasks plus the QA gate is a
    reference, not a cap: a Spec that needs more is a signal to check whether
    it holds more than one context, never a refusal by count."
  - `split-by-context`: "When the scope holds more than one context, split it
    into Specs that each hold one"
  - `one-owner`: "give each source exactly one owning Spec;"
  - `no-unrelated-source`: "never add a source that shares no context with the
    Spec only to save a Spec."
  - `requires`: the last sentence of API Contract 1 from "name that Spec in the
    Task Graph manifest's" to "in the TechSpec Build Order.", backticks
    included.
- `skills/acceptance_rules_test.go` gains
  `TestTheAuthoringSkillsSizeASpecByContextAndVerificationCycle`, which, over
  both copies through `roundfixLayoutSources`, reads each file below, compares
  through `normalizeAcceptanceText`, fails when the text contains "fits four
  implementation Tasks plus its QA gate" or "The grouping bound is four
  implementation Tasks", and requires each phrase exactly once, in a subtest
  named `<copy>/<path>`:
  - the three authoring skills, `write-idea`, `write-prd` and
    `write-techspec`, each read at `<skill>/SKILL.md`, hold these seven
    phrases:
    - "The bound on grouping is a criterion, not a count:"
    - "one context, where the sources and the work change the same component
      or contract, share a root cause, or complete one another,"
    - "and Tasks that each earn their verification cycle."
    - "joins a neighbouring Task of the same context,"
    - "small source of the same context joins a Spec that is open or being
      planned rather than becoming a one-Task Spec."
    - "Four implementation Tasks plus the QA gate is a reference, not a cap:
      needing more is a signal to check whether the Spec holds more than one
      context."
    - "When it holds more than one, split the scope and give each source one
      owning Spec."
  - the `write-tasks` skill, read at `write-tasks/SKILL.md`, holds these five
    phrases:
    - "**Sized for one fresh session.**"
    - "**Earns its verification cycle.** Every Task pays an Agent turn, its
      declared Verification and the repository Verification at settlement."
    - "joins a neighbouring Task of the same context instead of becoming its
      own Task."
    - "It joins only a Task whose slice it belongs to, so the vertical-slice
      rule holds, and Tasks that declare the same file still run in series
      under the Wave collision rule."
    - "a graph that needs more is a signal to check whether the Spec holds
      more than one context, never a refusal by count."
  - the Roundfix `spec` reference, read at `roundfix/references/spec.md`,
    holds this phrase:
    - "The grouping bound is a criterion, not a count: one context, and Tasks
      that each earn their verification cycle; four implementation Tasks plus
      its QA gate is a reference, not a cap."

Sabotages the QA runs, each on its own temporary copy of the audited head; on
the prototype, removing each phrase of both lists from the rendered golden or
the skill file made its named subtest fail:

1. The clause's guidance restored to its `251ed81e` text, with the module record
   command and `make baseline-digests` run in that copy:
   `TestTheGroupingBoundIsACriterionThatKeepsEveryGuarantee` fails.
2. In the clause, "give each source exactly one owning Spec; " removed (same
   regeneration): the same test fails in its `one-owner` subtest.
3. In the clause, the sentence that begins "Four implementation Tasks plus the
   QA gate is a reference" removed (same regeneration): the same test fails in
   its `reference-not-cap` subtest.
4. In the clause, "never add a source that shares no context with the Spec only
   to save a Spec" replaced by "add any source" (same regeneration): the same
   test fails in its `no-unrelated-source` subtest.
5. The "**Earns its verification cycle.**" line removed from the canonical
   `write-tasks` skill: `TestTheAuthoringSkillsSizeASpecByContextAndVerificationCycle`
   fails in its `canonical/write-tasks/SKILL.md` subtest.
6. The embedded `skills/write-prd/SKILL.md` grouping sentence restored to its
   `251ed81e` text: the same test fails in its `embedded/write-prd/SKILL.md`
   subtest.

## Build Order

1. task_01: the Baseline clause, its tests, the derived files and this
   repository's managed refresh.
2. task_02: the authoring and `write-tasks` skills, the Roundfix `spec`
   reference, the `spec` command guide, the glossary and ADR-0208's notice;
   depends on step 1, because the skills restate the clause task_01 rewrites
   and both Tasks run the documentation contracts.
3. task_03: the QA gate; depends on step 2.

No other Spec is required. Specs 0392 and 0393 deliver after this one and
share the files `_prd.md` → Prerequisites names: the `rule.spec.routing`
version, the module record and every derived Baseline file (0393), the
authoring skills, `write-tasks`, the owned-skill versions and their record
(0392 and 0393), the Roundfix `spec` reference and the `spec` command guide
(0393), `CONTEXT.md` and `docs/references/behavior-surfaces.json`. Each Task
regenerates derived files and records from the tree it runs on and never
writes a version number.

## Risks & Considerations

- An author may read "a reference, not a cap" as license for large Specs; the
  clause ties any Spec past the reference to a check for a second context, and
  the judge's grouping suggestion and review still apply.
- An author may fold unrelated changes into a neighbouring Task to save a
  cycle; the rule allows it only inside that Task's slice and context, and the
  Wave collision rule still serializes shared files.
- Adopters running an older Roundfix keep the old clause until their next
  Baseline update.

## Glossary

- changes: **Grouping Bound**

## Decisions

Autonomous route: no person present; decisions taken from `_prd.md` → Decisions, the operator's cycle BF brief for Spec 0391 and the characterization above.

- The clause keeps its id; the Source Baseline row is retained by id and
  enforcement and its corpus stays immutable, so no `replaces` and no Source
  Baseline change are needed.
- `rule.spec.routing`'s version rises by one from the value the Task finds, so
  Spec 0393, which changes the same rule, never collides on a number.
- The pinning tests read the rendered guide and the skill copies, so a later
  reword that drops a guarantee fails by name even when the literal is updated
  with it.
- The user-guide sentence and the Roundfix reference change because they state
  the bound; the Review Segment bounds and dated reference documents do not.
