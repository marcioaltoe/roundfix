---
spec: 0193-baseline-guides-that-describe-the-product-as-it-is
prd: _prd.md
created: 2026-09-30
---

# Baseline guides that describe the product as it is — Technical Spec

## Executive Summary

Eleven clauses in five Baseline modules change their guidance, one duplicate
entry is removed, and four tests guard what the audit found. The module JSON
is the only source: the formatter goldens, the digest pins, the catalog
snapshots and this repository's guides are regenerated from it by the
sanctioned commands. No clause changes its identifier or its enforcement
level. One small code change makes `deferred` a terminal Backlog status. The
trade-off accepted is that the exact clause texts are fixed here, in the
TechSpec, so an implementing Agent copies sentences instead of composing them.
That costs a longer document and buys guidance whose wording was reviewed once.

## Project Constraints

- Identifier strategy: not applicable — no new identifier; clause, rule, guide
  and module identifiers are kept. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — Baseline assets, local files and
  Git only; no credential and no network call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0186 governs the clause wording.
  ADR-0014, ADR-0092, ADR-0104, ADR-0130, ADR-0163, ADR-0166, ADR-0167 and
  ADR-0179 are the decisions the clauses restate, and none changes. ADR-0073,
  ADR-0081, ADR-0103 and ADR-0149 govern the regeneration and this
  repository's refresh. ADR-0080, ADR-0088, ADR-0091, ADR-0093, ADR-0094,
  ADR-0096, ADR-0117, ADR-0155 and ADR-0156 bind the gate and the checker.
  ADR-0178 has the audit read the grant a Task commit ran under, and the
  unchanged second sentence of the bounded-execution clause stays true under
  it. Nine ADRs cite a listed ADR and do not apply, because this Spec touches
  none of their behaviors: ADR-0020 (the parsed prompt result), ADR-0038 (the
  one Verification repair), ADR-0057 (the Daemon as only writer of Task
  status), ADR-0127 (process residue) and ADR-0160 (the red repository gate)
  cite ADR-0014; ADR-0097 (carrying a QA row forward) cites ADR-0080; ADR-0168
  (the related-ADR gap) and ADR-0176 (citation checks read authored text) cite
  ADR-0093; ADR-0173 (citations a History Relocation breaks) cites ADR-0073.
  Reached only through those citations, and equally untouched: ADR-0056, ADR-0159, ADR-0170, ADR-0177.
  All hold. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30, recorded in [_authorization.md](_authorization.md); bounded
  files: `internal/baseline/assets/modules/backend.json`,
  `internal/baseline/assets/modules/autonomous-work.json`,
  `internal/baseline/assets/modules/core.json`,
  `internal/baseline/assets/modules/spec-workflow.json`,
  `internal/baseline/assets/modules/context-workflow.json`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/backend.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `docs/agents/autonomous-work.md`, `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`, `docs/agents/docs-layout.md`,
  `docs/agents/setup-context.json`, `docs/agents/specific-repository.md`,
  `internal/speccheck/backlog.go`, `internal/baseline/plan_test.go` (added for task_06). Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## System Architecture

| Concern | Source of truth | Derived by |
| --- | --- | --- |
| Clause text and versions | `internal/baseline/assets/modules/*.json` | hand edit under the grant |
| Formatter goldens, the Standard TypeScript Monorepo digest pin, catalog snapshots, plan characterization goldens | the modules | `make baseline-digests` |
| This repository's guides and Setup Manifest | the modules | `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text` |
| Backlog terminal statuses | `internal/speccheck/backlog.go`, `internal/spec/retirement.go` | — |
| Repository-owned guide | `docs/agents/specific-repository.md` | hand edit under the grant |

Every Task edits one or more modules and then runs the two regeneration
commands, in that order. A second managed refresh must report
`File changes: 0`. No pin, golden or generated guide is edited by hand.

The regeneration outputs were measured on 2026-09-30 in a scratch clone of
`9e439dbb`, by editing each module and running both commands:

- every module edit rewrites
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/baseline/testdata/catalog.diagnostics.golden.json`,
  `internal/baseline/testdata/catalog.digest`,
  `internal/baseline/testdata/catalog.normalized.json`, the four plan goldens
  under `internal/baseline/testdata/plan-characterization/`
  (`advisory-only-divergences`, `clean-adoption`,
  `idempotent-replan-after-verified-apply`,
  `same-baseline-changed-profile-and-catalog-digests`), and
  `docs/agents/setup-context.json`;
- each module also rewrites its own formatter golden and this repository's
  matching guide: `backend.md` (golden only; this repository has no backend
  guide), `autonomous-work.md`, `agent-instructions.md`, `spec-routing.md`
  and `docs-layout.md`.

The Source Baseline corpus, the parity corpus and
`internal/baseline/assets/setups/` stay byte-identical.

A module edit must preserve the file's existing formatting. A rewrite through
a JSON encoder reformats `core.json`, and the catalog diagnostics test then
fails on an asset anchor it no longer finds. Each edit therefore replaces the
guidance string and the version numbers in place.

## Implementation Design

### Interfaces

No exported Go signature changes. The tests use the package's existing
embedded-catalog loader.

```go
// internal/speccheck/backlog.go
func terminalBacklogStatus(status string) bool // gains "deferred"

// internal/spec/retirement.go: "deferred" retires like "declined" when the
// entry carries a reason, and like the other terminal statuses otherwise.
```

### Data Models

None. No schema, record or payload changes. Module, rule and guide versions
each rise by one where their content changes.

### Version changes

For each changed clause, raise by one the version of its rule, of every guide
that lists that rule, and of the module. A module changed by two Tasks rises
once per Task.

### Clause texts

Each replacement below is exact. Text outside the quoted sentences stays
byte-identical, and every clause keeps its `enforcement`.

**`backend.json` — `rule.backend.boundary-contracts`.** Remove the last entry
of the rule's `clauses`, whose `id` is `rule.backend.boundary-contracts`. The
entry `clause.backend.boundary-contracts` keeps the same guidance.

**Clause replacement for adopters (task_06).** An adopter whose Setup
Manifest declares the Standard TypeScript Source Baseline is classified clause
by clause against that Source Baseline, which records both
`rule.backend.boundary-contracts` and `clause.backend.boundary-contracts`.
The first delivery Run showed that removing one entry makes the classifier
report it `unaccounted`, and a plan with an unaccounted clause is refused, so
every such adopter's next Baseline update would stop. The surviving clause
therefore declares `"replaces": ["rule.backend.boundary-contracts"]`. The
classifier reports a Source Baseline clause that is absent from the selected
catalog as `replaced`, with the replacing clause as its target, when exactly
one selected clause lists it in `replaces` and carries the same enforcement.
Any other absence stays `unaccounted`. Catalog validation refuses a
`replaces` value that is not a list of unique non-empty clause IDs, an ID that
is still a clause in the catalog, and an ID that two clauses claim. The
declaration is never rendered, so no guide changes. The `replaced`
disposition already exists in the plan contract; only its producer is new.

**`autonomous-work.json` — `clause.autonomous.hook-strictness`.** Replace
"ADR-0014 makes the Daemon the verification authority, and it commits only
after that Verification passes" with "The Daemon is the verification
authority, and it commits only after that Verification passes".

**`autonomous-work.json` — `clause.autonomous.loop-01-qa-once`.** The new
guidance is, in full:

> Follow one order per Spec: implement the graph including its authored gate,
> apply the configured pre-PR review policy, archive and commit the candidate,
> run the repository gate, push the candidate, open the Pull Request, verify
> required checks for the current head, and merge. The Delivery Queue (`roundfix deliver`) runs
> this order for each queued Spec, and a Supervisor driving the steps by hand
> follows the same order. For `codex`, `claude`, or `coderabbit`, require the
> selected provider's current-candidate review. For explicit `none`, skip
> review and provider calls, record the configured omission and continue
> through the other gates. A legacy PR-feedback Watch loop is not a mandatory
> step for every policy, and a missing local adapter never becomes a
> successful review. The authored QA gate is the graph's terminal Task, so the
> gate runs before any Pull Request exists. An environment-blocked row can
> still reach `pass` when the report carries equivalent evidence, and the row
> blocked only because no Pull Request is open never decides a qualifying
> `partial`. Run the Implement Command while implementation work is pending or
> a corrective Task is planned; normally the authored QA gate must also
> settle. When a corrective Task is added after the gate settled, reopen the
> gate with `roundfix reopen --spec <slug>`; never edit the QA Task by hand.
> When a Run ends with Tasks already settled, hand them back with Task
> Carry-Forward before amending the Spec, so proved work is not run again: use
> `roundfix reconcile <run-id> --carry-forward`, or `roundfix deliver retry
> <slug>` for a parked Delivery Queue item. An explicit user-authorized QA
> Archive Override may waive only the covered QA gate's archive prerequisite,
> with the approval and original QA/Task state preserved. Consume an
> applicable prior approval without asking again. That archive does not
> establish Implement Clean or satisfy QA, review, publication, merge or
> required-check obligations; evaluate those separately before continuing
> delivery. The Daemon gates the authored QA Task on every dependency reaching
> completed; re-running the gate after each corrective Task turns discovery
> into a serial chain of full cycles. Close several findings together; more
> than two corrective Tasks generated by QA findings means the decomposition
> needs re-examining. Rebuild any binary the gate exercises before the gate
> becomes runnable, or the gate reports defects the running artifact does not
> contain.

The guidance is one JSON string on one line, with single spaces where the
quote above wraps. It removes the three ADR citations and the sentence about
Spec 0078, and keeps every other sentence of the old clause.

**`core.json` — `clause.core.verification-two-tiers`.** Replace the sentence
that begins "An untrusted source requires an `execution_approvals` entry" and
ends "the approval expires when any of those fields changes." with:

> A Run satisfies this by construction, because its Run Worktree is created
> from a commit. No other command checks it: `roundfix spec check
> --run-verification` executes the commands of the working-tree Spec in a
> checkout of `HEAD`, and `roundfix settle` reads the Task file of the
> directory it settles, so whoever runs them on an uncommitted or modified
> Spec first reads the commands they will execute. No approval record makes
> an untrusted source executable; commit the artifacts or do not execute
> them.

**`core.json` — `clause.core.tooling-commit-choreography`.** Replace "a
proposed, absent, contradictory, or withdrawn record grants no governed
mutation and does not by itself refuse implementation, commit, or push, and
an executor cannot widen the grant it depends on." with:

> a proposed, absent, contradictory, or withdrawn record grants no governed
> mutation and no delivery operation. `roundfix implement` still runs a Spec
> that changes no Governed Path, and `roundfix deliver start` refuses a Spec
> whose record does not grant `implement`, `commit`, `push`, `pull_request`,
> and `merge`. An executor cannot widen the grant it depends on.

**`spec-workflow.json` —
`clause.spec.project-constraints-02-tooling-authorization`.** Append after
"An absent `operations` list grants no operation.":

> A Spec that changes no Governed Path records `paths: []`, which grants its
> listed operations and bounds no path; a record without a `paths` key is
> refused.

**`spec-workflow.json` —
`clause.spec.project-constraints-03-bounded-execution`.** Replace the first
sentence, up to "if changed-file postflight finds another path.", with:

> An authorized tooling Task may change a Governed Path, that is a protected
> tooling path, only when its grant bounds that path; stop before any other
> governed mutation, because the Task fails when its commit changes a Governed
> Path outside the grant. An ordinary path a Task changes without declaring it
> is not refused: the Daemon records it in the Task file under `## Recorded
> paths`, and the QA gate discloses it.

The second sentence, about a new or widened grant, is unchanged.

**`spec-workflow.json` —
`clause.spec.project-constraints-06-outside-evidence`.** Replace the last
sentence, which begins "When the outside source cannot be obtained", with:

> When the outside source cannot be obtained during authoring, record the row
> as blocked with that reason and continue: decomposition never stalls and
> never asks a person. The QA gate then holds Pull Request preparation until
> the row is satisfied or carried forward on declared unmoved evidence.

**`context-workflow.json` — `clause.context.adr-02-active-status`.** Replace
"Only `accepted` is active." with "For an ADR that carries lifecycle
frontmatter, only `accepted` is active."

**`context-workflow.json` — `clause.context.backlog-01-operational-contract`.**
Three replacements:

- in the frontmatter comment, `# open | promoted | declined | done | deprecated
  | superseded | closed | cancelled` gains `deferred` after `declined`;
- "Set `status: declined` only with a non-null `reason`." becomes "Set
  `status: declined` or `status: deferred` only with a non-null `reason`.";
- the terminal list "(`declined`, `done`, `deprecated`, `superseded`, `closed`,
  or `cancelled`)" gains `deferred` after `declined`.

**`context-workflow.json` — `clause.context.docs-one-job-per-directory`.**
Replace "Under ADR-0163, a finished orphan Review Artifact retires" with "A
finished orphan Review Artifact retires".

**`context-workflow.json` — `clause.context.inbox-01-triage`.** Replace
"Preserve the ADR-0092 boundary: evidence never becomes intent without a
human choice." with "Preserve the boundary between evidence and intent:
evidence never becomes intent without a human choice."

**`docs/agents/specific-repository.md`.** Replace "Nothing under a `_archived`
tree is ever validated." with "Nothing under `docs/history/` is ever validated
as live work." This file is repository-owned and is edited directly.

### Clauses left as they are, deliberately

Roundfix delivers the Baseline and runs the loop, so these clauses name it on
purpose and do not change:

- `clause.core.ask-user-answerable-decisions` names the structured-question
  tool of each supported Agent runtime.
- `clause.core.conventional-commit-titles` names the `roundfix/run-` branch
  namespace the tool owns.
- `clause.core.request-review-explicitly` names the review policies Roundfix
  configures, including CodeRabbit as an optional one.
- `clause.spec.keep-artifacts-in-spec-folder` names
  `roundfix archive --qa-override`.
- `clause.autonomous.loop-06-branch-hygiene` explains how Roundfix integrates
  a Run Branch.

The Bun module's command clause and `clause.core.follow-dependency-workflow`
are not touched.

### API Contracts

1. API Contract: rendered guide `docs/agents/autonomous-work.md` — the loop
   clause declares review before archive and names `roundfix deliver`,
   `roundfix reopen --spec <slug>` and Task Carry-Forward.
2. API Contract: rendered guides `docs/agents/agent-instructions.md` and
   `docs/agents/spec-routing.md` — the authorization, verification and
   outside-evidence clauses carry the texts above and none of the four removed
   phrases.
3. API Contract: rendered guide `docs/agents/docs-layout.md` — the Backlog
   contract lists `deferred`, and the ADR and triage clauses cite no ADR
   number.
4. API Contract: `roundfix spec check` — an active Backlog Entry with
   `status: deferred` is reported as a terminal entry that has not left
   `docs/backlog/`, not as an unknown status.

## Coverage Map

- Goal 1 → Clause texts (loop clause); Testing Approach 2.
- Goal 2 → Clause texts (core and spec-workflow clauses); Testing Approach 3.
- Goal 3 → Clause texts (backend, context-workflow); Testing Approach 1 and 4.
- Goal 4 → Testing Approach 1–4.
- Core Feature 1 → Clause texts (backend); Testing Approach 1.
- Core Feature 2 → Clause texts (loop, hook); API Contract 1; Testing
  Approach 2.
- Core Feature 3 → Clause texts (core, spec-workflow); API Contract 2;
  Testing Approach 3.
- Core Feature 4 → Clause texts (context-workflow, repository guide); API
  Contracts 3 and 4; Testing Approach 4.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → Testing Approach 3.
- Success Metric 4 → Testing Approach 4.
- Success Metric 5 → Testing Approach 1.
- Success Metric 6 → Testing Approach 5.
- API Contract 1 → Clause texts (loop clause).
- API Contract 2 → Clause texts (core and spec-workflow clauses).
- API Contract 3 → Clause texts (context-workflow).
- API Contract 4 → Interfaces.

## Integration Points

- **`SC-LOOP-ORDER-DIVERGENT`.** The detector compares the shipped golden,
  this repository's guide and the module. All three are regenerated from the
  module in one Task, so they stay equal. Its marker, "Follow one order per
  Spec:", and the full stop that ends the order sentence are kept.
- **Skills that repeat the loop order.** The Roundfix skill and the
  write-tasks skill carry their own copy of the order sentence. Spec 0192
  owns them; this Spec does not edit a skill.
- **Adopters.** A managed refresh reports each changed guide as an update of
  its managed entry. No retention mapping names a changed clause.

## Testing Approach

All tests read the embedded catalog or repository files and need no network.
Each negative case is its own test.

1. **Uniqueness and force.** New `internal/baseline/clause_characterization_test.go`:
   - `TestNoTwoBaselineClausesShareText` normalizes every clause-level and
     rule-level guidance string (lower case, single spaces) and expects no
     repeat.
   - `TestADuplicatedClauseTextIsReported` runs the same check on a copy of
     the clauses with one guidance repeated and expects that pair back.
   - `TestBaselineClauseForceIsCharacterized` compares every clause
     identifier and enforcement level with a table written in the test. The
     table holds every clause on `9e439dbb` except the duplicate backend
     entry. It fails for a missing clause, a new clause and a changed level.
   - `TestTheBackendBoundaryParagraphRendersOnce` counts the boundary
     sentence in the backend formatter golden.
2. **Loop order.** New `internal/baseline/loop_clause_test.go`:
   `TestLoopClauseNamesTheDeliveryQueueAndItsRecoveryActs` and
   `TestLoopClauseCitesNoSpecOrDecisionNumber`. New
   `internal/delivery/loop_clause_order_test.go`:
   - `TestTheLoopClauseOrderMatchesTheDeliveryQueue` reads the order sentence
     from the module file, maps its phrases to the Delivery Queue's actions
     (`run`, `review`, `archive`, `gate`, `push`, `pull-request`, `checks`,
     `merge`) and compares them with the actions the engine records for a
     reviewed item in the package's existing fake workflow. Every action the
     engine records must have a phrase in the sentence: an unmapped engine
     action fails the test, so a declared loop cannot be shorter than the
     real one.
   - `TestALoopClauseThatArchivesBeforeReviewIsRefused` gives the comparison
     the old order and expects a mismatch.
3. **Authorization and evidence.** New
   `internal/baseline/authorization_clauses_test.go`:
   - `TestAuthorizationClausesStateTodaysRefusals` expects each new sentence
     in its embedded clause and in its formatter golden.
   - `TestAuthorizationClausesDropTheRemovedPhrases` expects none of
     `execution_approvals`, `does not by itself refuse implementation`,
     `may mutate only its bounded` and `never blocks the Spec` in any module
     or formatter golden.
4. **Lifecycle and citations.** New
   `internal/baseline/adopter_neutral_clauses_test.go`:
   - `TestShippedGuidanceCitesNoRepositoryRecord` expects no guidance string
     to match `ADR-` followed by digits or `Spec ` followed by four digits.
   - `TestARepositoryRecordCitationIsReported` runs the check on a guidance
     string that cites one and expects it back.
   - `TestLifecycleClausesCarryTheScopedWording` expects the ADR scope
     sentence and `deferred` in the Backlog contract.
   New `internal/speccheck/backlog_deferred_test.go`:
   `TestADeferredBacklogEntryIsTerminal` and
   `TestADeferredBacklogEntryLeftActiveIsReported`. New
   `internal/spec/retirement_deferred_test.go`:
   `TestADeferredBacklogEntryWithAReasonIsRetired`.
5. **Regeneration.** Each Task's Verification runs the existing
   `TestFormatterComposition`, `TestCatalogCompatibility`,
   `TestBaselinePlanCharacterization` and `TestBaselineCompatibilityCorpus`,
   then the read-only managed refresh, which exits non-zero while a plan is
   pending.

## Build Order

1. Uniqueness and force checks, and the backend duplicate removal, task_01
   (depends on: none).
2. The loop and hook clauses, with the order check against the Delivery
   Queue, task_02 (depends on: 1).
3. The authorization, verification and outside-evidence clauses, task_03
   (depends on: 2).
4. Lifecycle wording, the `deferred` status, the citation check and the
   repository guide, task_04 (depends on: 3).
5. The declared replacement of the removed backend entry, and the two plan
   contracts this Spec moves, task_06 (depends on: 4).
6. Terminal QA, task_05 (depends on: 1, 2, 3, 4, 6).

The chain is serial because every module edit rewrites the same digest pin,
catalog snapshots and plan goldens.

## Risks & Considerations

- **Two Specs edit `core.json`.** Spec 0195 adds one sentence to the release
  clause of the same module. It must be delivered after this Spec, and it then
  regenerates from the state this Spec leaves.
- **Another Spec rewords clauses in `backend.json`.** The stack-clause Spec
  edits clause texts in the `bun`, `typescript`, `backend` and `frontend`
  modules. task_01 removes one entry from `backend.json` and changes no clause
  text, so the two edits do not overlap; they still regenerate the same
  derived files and must be delivered one after the other.
- **The grant a Task ran under.** ADR-0178 refines the grant the audit
  accepts, and adds no Baseline text. The unchanged second sentence of the bounded-execution clause
  stays true under it.
- **The force table is long.** It lists about 137 clauses. It is written once
  and changes only when a clause is added, removed or re-levelled, which is
  the event it exists to make visible.
- **A guide sentence can still be wrong.** The checks cover order, duplicates,
  force and citations. The verification clause now admits a missing control,
  and the open Backlog Entry carries the fix.
- **Hazard: generated files.** A Task must never hand-edit a golden, a pin or
  a rendered guide. The regeneration is the only writer, and the second
  refresh proves convergence.

## Decisions

- Clause texts are fixed in this document. See ADR-0186.
- The order check lives in the delivery package, next to the fake workflow
  that records the engine's actions, and reads the module as a file.
- The citation check forbids the number forms and nothing else, so a clause
  may still name a concept such as the Delivery Queue.
- `deferred` reuses the rule of `declined`: terminal, with a reason.
