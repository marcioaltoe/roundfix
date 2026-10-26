---
spec: 0323-the-status-the-park-table-and-write-tasks-say-what-to-do-next
prd: _prd.md
created: 2026-10-26
---

# The status, the Park Class table and write-tasks say what to do next — Technical Spec

## Executive Summary

Three small corrections, one Task each, in series because they share the
Roundfix skill's version record and the Behavior Surface record.

1. `roundfix follow-up --status` renders the Release Plan's
   `manual_classification_required` state with a reason and a next command
   instead of an empty version.
2. The Park Class table of the `deliver` guide and both Roundfix reference
   copies lists blockers only; the `secret-found` next action becomes its own
   sentence, and a documentation contract compares that sentence and the
   `archived-item-unproven` one with `delivery.ClassifyPark`.
3. The `write-tasks` skill and `CONTEXT.md` **Recorded Path** tell authors
   to declare a test helper, fake, fixture or golden that two Tasks may each
   change in every such Task, because a Recorded Path's settlement commit is
   still Wave collision evidence and the QA gate refuses the pair after both
   have run.

The trade-off accepted for the third source: guidance prevents the lost Run
only when the author foresees the shared file. No code can do better without
changing the collision rule, because the earliest moment a structural
diagnostic has both settlement commits is after both Tasks ran.

## Project Constraints

- Identifier strategy: not applicable — no identifier or schema is minted.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local reads only. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0143 and ADR-0166, and ADR-0176
  recorded as not applying, as `_prd.md` → Project Constraints states; none
  changes. Source:
  `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — protected paths limited to the `roundfix`
  and `write-tasks` skills. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0323-the-status-the-park-table-and-write-tasks-say-what-to-do-next/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/follow-up.md`,
  `.agents/skills/write-tasks/SKILL.md`, `skills/roundfix/SKILL.md`,
  `skills/write-tasks/SKILL.md`.

## Current behavior

Read on 2026-10-26 at `aa8b3db7`.

1. `internal/cli/follow_up_status.go` `runFollowUpStatus` builds the Release
   Plan with `releaseplan.Request{}` (target `HEAD`), copies its state, base
   tag and `ProposedVersion` into `followUpRelease`, and for
   `manual_classification_required` appends the maintainer item
   `{Kind: "release_classification", Detail: "manual classification required"}`;
   `release.reason` stays `nil`. `releaseplan.CalculateProposal` proposes no
   version in that state, and `plan.Classification.BlockingCommits` lists the
   commits that need the classification.
2. `renderFollowUpStatusText` prints `Next release: %s from %s (%s)` for every
   state but `no_unreleased_changes`, `no_release` and `unavailable`, so the
   manual state prints `Next release:  from <base> (manual_classification_required)`,
   and prints each non-question maintainer item as `  release %s: %s` with
   `ProposedVersion`, so the manual item prints `  release : manual classification required`.
3. `internal/cli/follow_up_status_test.go`
   `TestFollowUpStatusReportsTheNextRelease/manual` pins
   `ProposedVersion == ""`, `Reason == nil`, one maintainer item of kind
   `release_classification`, and a text containing `manual classification`.
4. `docs/user-guide/commands/deliver.md` (`| Class | Blockers |`) and
   `.agents/skills/roundfix/references/deliver.md` with its mirror
   (`| Class | Blockers or action |`) end the `finding` row with
   `` `secret-found`: in `<worktree>`, rewrite the item branch … then run
   `roundfix deliver retry <slug>` ``. Each file then states
   ``An `archived-item-unproven` park's next action is `…`.`` in a sentence of
   its own. `internal/delivery/park_class.go` `ClassifyPark` gives the rewrite
   only to `BlockerSecretFound`. No test reads the table.
5. `internal/spec/collision.go` `Collisions` unions each Task's declared and
   Verification paths with the paths of its newest settlement commit
   (`priorRunTaskTouches`, `git log --all`). `internal/daemon/task_engine.go`
   takes each Task's Spec Consistency baseline when the Task starts, so a
   later settlement does not refuse a pair that collided before it started,
   and `specGatePrecondition` runs the strict check with no baseline. Two
   Tasks of one Wave that each change an undeclared file are therefore first
   refused at the QA gate, after both ran.
6. `.agents/skills/write-tasks/SKILL.md` (version `0.0.19`): "The Daemon
   records a path a Task changed without declaring it under `## Recorded
   paths` at commit. … recording discloses a change and reserves nothing.
   Declaring every foreseeable path stays the rule." `CONTEXT.md`
   **Recorded Path** says recording "reserves nothing". Neither says the
   settlement commit is collision evidence.

## System Architecture

- `internal/cli/follow_up_status.go` and its test (task_01), with the
  follow-up guide, the Roundfix `follow-up` reference, the `roundfix` skill
  version pair and the Behavior Surface record.
- `docs/user-guide/commands/deliver.md`, both Roundfix `deliver` reference
  copies and `internal/docscontract/command_documentation_test.go` (task_02),
  with the `roundfix` skill version pair and the Behavior Surface record.
- `.agents/skills/write-tasks/SKILL.md` and its mirror and `CONTEXT.md`
  (task_03), with the `write-tasks` version pair.
- `skills/testdata/owned-skill-versions.json` in every Task.

## Implementation Design

### Interfaces

In `runFollowUpStatus`, case `releaseplan.StateManualClassificationRequired`:

```go
doc.Release.Reason = followUpString(fmt.Sprintf(
	"%d commit(s) since %s need a manual release classification",
	len(plan.Classification.BlockingCommits), plan.Base.Tag))
doc.Maintainer = append(doc.Maintainer, followUpMaintainer{Kind: "release_classification",
	Detail: fmt.Sprintf("classify them with roundfix release plan --from %s --to %s --impact <none|patch|minor|major> --reason <text>",
		plan.Base.Tag, plan.Target.Name)})
```

The command matches the `Next action` that `roundfix release plan` prints in
the same state. `renderFollowUpStatusText` gains a
`string(releaseplan.StateManualClassificationRequired)` case printing
`Next release: not proposed (manual_classification_required): <reason>`, and
prints a `release_classification` item as `  release classification: <detail>`;
the `release_approval` line keeps `  release <version>: <question>`.

The follow-up guide, after the `Next release: no unreleased changes` block:

```markdown
When the Release Plan needs a manual classification, the status proposes no
version and says why instead of leaving the version blank. Its `Next release`
line reads
`Next release: not proposed (manual_classification_required): <n> commit(s) since <base> need a manual release classification`,
and `Waiting on the maintainer` lists
`release classification: classify them with roundfix release plan --from <base> --to HEAD --impact <none|patch|minor|major> --reason <text>`.
The JSON carries the first sentence in `release.reason` and the second in the
`detail` of the `release_classification` item of `maintainer`;
`release.proposed_version` stays empty, as in every state that proposes no
version.
```

The canonical Roundfix `references/follow-up.md`, after the paragraph ending
"it never calls the judge.":

```markdown
When the Release Plan needs a manual classification, the status `release.state`
is `manual_classification_required` and `release.proposed_version` is empty.
`release.reason` says how many commits since the base tag need a manual release
classification, and the `release_classification` item of `maintainer` names the
`roundfix release plan --from <base> --to HEAD --impact <none|patch|minor|major> --reason <text>`
command that classifies them.
```

In the guide and the canonical `references/deliver.md`, the `finding` row
ends at `` `secret-found` |``, and directly after the
`archived-item-unproven` sentence:

```markdown
A `secret-found` park's next action is `in <worktree>, rewrite the item branch from the earliest commit the blocker names so that each named line holds a value no Secret Pattern matches and nothing else in the final tree changes, then run roundfix deliver retry <slug>`.
```

### Data Models

None. `followUpRelease` and `followUpMaintainer` keep their fields and JSON
tags; only the values of `release.reason` and `detail` in the manual state
change.

### Shared test files in the Task Graph

This section delivers PRD Goal 4: an author who follows it declares a shared
test file in both Tasks, so `spec check --strict` reports the collision and
names the `needs` edge before the Run starts, instead of the QA gate refusing
the Run after both Tasks ran.

In `.agents/skills/write-tasks/SKILL.md`, after the paragraph ending
"Declaring every foreseeable path stays the rule.":

```markdown
A test helper, fake, fixture or golden that two Tasks may each change is a
foreseeable path. Declare it under `interface:` in every Task that may change
it, so `roundfix spec check --strict` reports `SC-WAVE-COLLISION` while the
Spec is still authored and the `needs` edge it names runs those Tasks in
series. Recording reserves nothing, but a Task's settlement commit is
collision evidence: when two Tasks of one Wave each change the same
undeclared file, the QA gate's precondition refuses the Run with
`SC-WAVE-COLLISION` after both Tasks have run.
```

`CONTEXT.md` **Recorded Path** gains, after its last definition sentence and
before its `_Avoid_` line:

```markdown
Its settlement commit is still collision evidence: two Tasks of one Wave that
each change the same Recorded Path make the QA gate's precondition refuse with
`SC-WAVE-COLLISION`, so a file two Tasks may each change is declared by both.
```

### API Contracts

1. API Contract: in `manual_classification_required`, `release.reason` is
   `<n> commit(s) since <base> need a manual release classification` with
   `<n>` the number of blocking commits, and the one `release_classification`
   maintainer item's `detail` is the classification command; text prints both
   and no blank version.
2. API Contract: every other release state's text and JSON bytes are
   unchanged, `release.proposed_version` is `""` in every state without a
   version, and the exit code stays `0`.
3. API Contract: the Park Class `finding` row of each documented copy lists
   only backticked blocker names, and each blocker-specific sentence reads
   ``A[n] `<blocker>` park's next action is `<ClassifyPark Next>`.`` with
   `<worktree>` and `<slug>` as placeholders.

### Surface Transcripts

None. The two new lines are fixed in Interfaces and asserted whole by a named
test.

### Invariants

```text
1. The status view never prints an empty version field.
2. The manual classification count is the Release Plan's blocking commits, not all changes.
3. The Park Class table never attaches a next action to a class row.
4. A documented blocker-specific next action equals the action ClassifyPark prints.
5. The Wave collision rule, its evidence and the QA gate precondition are unchanged.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, Core Feature 1, Success Metric 1 → API Contract 1, Invariants 1 and 2, task_01
- Goal 2, Success Metric 2 → API Contract 2, task_01
- Goal 3, Core Feature 2, Success Metric 3 → API Contract 3, Invariants 3 and 4, task_02
- Goal 4, Core Feature 3, Success Metric 4 → Shared test files in the Task Graph, Invariant 5, task_03

## Absence Evidence

1. A manual classification sentence in the status view — command: `grep -rn "need a manual release classification\|release classification:" internal/cli docs/user-guide .agents/skills`; result: no match at `aa8b3db7`.
2. A documentation contract for park next actions — command: `grep -rn "ClassifyPark\|Park Class" internal/docscontract`; result: no match at `aa8b3db7`.
3. A shared-test-file rule in `write-tasks` — command: `grep -rn -i "fake, fixture\|test helper" .agents/skills/write-tasks`; result: no match at `aa8b3db7`.

## ADR Alternatives

- ADR-0143: unchanged; the status view still reads the plan and mutates
  nothing, and it names the manual classification the plan requires, which
  the ADR asks the plan to name.
- ADR-0166: unchanged; this Spec keeps recording a disclosure that reserves
  nothing and tells authors to declare, the alternative the ADR keeps for
  reserving a path, instead of making a Recorded Path reserve anything.
- ADR-0176: does not apply; no citation check changes, and the Daemon-written
  `## Recorded paths` section it keeps out of those checks is untouched.

## Integration Points

None. No network, provider, forge or other repository is reached.

## Testing Approach

Every test builds temporary repositories or reads this repository's own
guides through the `docscontract` contract files.

- `internal/cli/follow_up_status_test.go`:
  - new `TestFollowUpStatusNamesTheManualClassification` (calls
    `t.Parallel()` first): tag `v0.4.1`, then commits `feat: expose widget
    status`, `chore: change widget behavior`, `chore: change widget storage`,
    each changing its own `internal/widget<N>.go`. JSON: state
    `manual_classification_required`, base `v0.4.1`, `proposed_version`
    `""`, `reason` exactly `2 commit(s) since v0.4.1 need a manual release
    classification`, three changes, and maintainer exactly
    `[{release_classification, "classify them with roundfix release plan
    --from v0.4.1 --to HEAD --impact <none|patch|minor|major> --reason
    <text>"}]`. Text: contains `\nNext release: not proposed
    (manual_classification_required): <reason>\n` and `\nWaiting on the
    maintainer: 1\n  release classification: <detail>\nSources:\n`, and
    neither `Next release:  ` nor `  release : `.
  - `TestFollowUpStatusReportsTheNextRelease` (declared change): its shared
    assertion expects a non-nil `reason` only for `manual`, and the `manual`
    case looks for `manual release classification`; every other case keeps
    its expectations.
- `internal/docscontract/command_documentation_test.go` (a `//verify:always`
  file):
  - new `TestTheParkActionsAreDocumentedAsTheQueueClassifiesThem` reads the
    guide and both reference copies and runs `checkParkActionDocumentation`,
    which requires every comma-separated item of the `finding` row's second
    cell to be one backticked blocker name, and for
    `delivery.BlockerArchivedItemUnproven` and `delivery.BlockerSecretFound`
    the whitespace-normalized sentence ``A[n] `<blocker>` park's next action is
    `<next>`.`` where `<next>` is `delivery.ClassifyPark` for an item with
    `SpecSlug: "<slug>"` and `Worktree: "<worktree>"`.
  - new `TestTheParkActionContractReportsAnActionInTheClassTable` feeds the
    guide with the joined row and with the `secret-found` sentence reworded,
    and asserts the two errors.
- A prototype in a scratch worktree passed both packages' tests,
  `go test ./skills` and the full `-tags docscontract ./internal/docscontract`,
  and each mutation failed the named test, whose input makes it observable:
  - counting `len(plan.Changes)` instead of the blocking commits:
    `TestFollowUpStatusNamesTheManualClassification` (three changes, two
    blocking);
  - dropping the text case for the manual state:
    `TestFollowUpStatusNamesTheManualClassification` (the default branch
    prints `Next release:  from v0.4.1`);
  - leaving `release.reason` nil: `TestFollowUpStatusNamesTheManualClassification`
    and `TestFollowUpStatusReportsTheNextRelease/manual`;
  - restoring the joined `finding` row in the guide:
    `TestTheParkActionsAreDocumentedAsTheQueueClassifiesThem`;
  - rewording the `secret-found` sentence in one reference copy:
    `TestTheParkActionsAreDocumentedAsTheQueueClassifiesThem`.

## Build Order

1. The manual classification status, its tests, the follow-up guide, the
   Roundfix `follow-up` reference, the `roundfix` skill version and the
   Behavior Surface record re-recorded last (task_01; no dependency).
2. The Park Class table split in the guide and both reference copies, the
   contract tests, the `roundfix` skill version and the Behavior Surface
   record re-recorded last (task_02; depends on 1: both raise the `roundfix`
   skill version, edit `skills/testdata/owned-skill-versions.json` and
   re-record `docs/references/behavior-surfaces.json`).
3. The `write-tasks` paragraph, the **Recorded Path** revision and the
   `write-tasks` version (task_03; depends on 2: both edit
   `skills/testdata/owned-skill-versions.json`).
4. QA gate (task_04) (depends on: 1, 2, 3).

No Spec is listed in `requires`. Specs 0319–0322 deliver first; nothing here
reads their code.

## Risks & Considerations

- A consumer that matched the old maintainer detail
  `manual classification required` or the `release : ` text no longer finds
  it. No Go package parses the status text, and the Roundfix skill's status
  mode reads the JSON `kind`, which is unchanged.
- The guidance in `write-tasks` cannot force an author to foresee a shared
  file; an unforeseen one still loses the Run at the gate, with the existing
  `SC-WAVE-COLLISION` summary naming the path and the edge.
- The new contract test pins two sentences; a later change to
  `ClassifyPark`'s wording must update the guide and both references in the
  same change, which is the drift it exists to catch.

## Glossary

- changes: **Recorded Path**

## Decisions

- Manual status: a reason sentence and a classification command, the
  `release.proposed_version` field kept empty and typed as today.
- Park Class table: blockers only per class; `secret-found` gets its own
  next-action sentence, pinned with `archived-item-unproven` to
  `ClassifyPark` by a contract test.
- Third source: option (a), guidance in `write-tasks` plus the glossary,
  because option (b) cannot act before both Tasks have run: a Task's
  settlement check compares with a baseline taken when it starts and runs
  before its own settlement commit exists, so the first point that holds
  both commits is the QA gate's precondition, which already refuses with the
  path, its source and the resolving edge. Both options are structural; only
  (a) acts before the Run.
- No new ADR.
