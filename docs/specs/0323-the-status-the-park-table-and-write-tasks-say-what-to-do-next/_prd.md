---
spec: 0323-the-status-the-park-table-and-write-tasks-say-what-to-do-next
status: active
created: 2026-10-26
surfaces: [backend, cli]
---

# The status, the Park Class table and write-tasks say what to do next

Minimal PRD for a bug fix, minted by `write-techspec`. Three texts an operator
or an authoring Agent reads tell them the wrong thing or nothing. `roundfix
follow-up --status` prints a blank version and an unexplained
`manual classification required` when the Release Plan needs a manual
classification. The `deliver` Park Class table attaches the `secret-found`
rewrite instruction to the whole `finding` row. And the `write-tasks` skill
never tells an author to declare a test helper or fake that two Tasks may each
change, so two Tasks of one Wave can both edit it undeclared and the QA gate
refuses the Run with `SC-WAVE-COLLISION` after both have run. After this Spec
the status names how many commits need a classification and the command that
classifies them, the table lists blockers only and states the rewrite for
`secret-found` alone, and the skill and the glossary tell authors to declare a
shared test file in every Task that may change it, so the strict check
reports the collision while the Spec is still authored.

## Sources

Three Backlog Entries, written and adopted at birth by this Spec's planning
change (see `references/_index.md`):

1. [`references/2026-10-26-follow-up-status-leaves-the-next-version-blank-under-manual-classification.md`](references/2026-10-26-follow-up-status-leaves-the-next-version-blank-under-manual-classification.md),
   `type: fix`, from the operator's observation of 2026-10-24 on Roundfix
   0.117.0 (cycle AO triage item 11).
2. [`references/2026-10-26-the-park-class-table-joins-the-secret-rewrite-to-every-finding-park.md`](references/2026-10-26-the-park-class-table-joins-the-secret-rewrite-to-every-finding-park.md),
   `type: fix`, from the pre-PR review finding F1 of Spec 0299, dismissed on
   2026-10-22 at 20:05 BRT as outside that Spec (triage item 4).
3. [`references/2026-10-26-two-tasks-of-one-wave-that-change-an-undeclared-test-fake-lose-the-run-at-the-qa-gate.md`](references/2026-10-26-two-tasks-of-one-wave-that-change-an-undeclared-test-fake-lose-the-run-at-the-qa-gate.md),
   `type: fix`, from Fluxus (ai-memory message `01a16a54`, 2026-10-23, triage
   item 10).

The three share one context: each is text the operator or an authoring Agent
reads to decide the next step, and each fix is a few lines of rendering or
documentation. No open Backlog Entry or unresolved Finding shares it: at the
authoring base `docs/backlog/` and `docs/findings/` hold no entry.

## Prerequisites

None in `requires`. Authored against `aa8b3db7`, with Specs 0319–0322 active
under `docs/specs/`; this Spec delivers after them. Spec 0321's task_03 adds a
paragraph to `docs/user-guide/commands/deliver.md` and both Roundfix
`deliver` references before the Linux Lane section; this Spec edits the Park
Class section of the same files and anchors every edit on text, so the two
changes do not overlap. Specs 0319, 0320 and 0321 raise the `roundfix` skill
version and re-record `docs/references/behavior-surfaces.json`; this Spec's
Tasks raise the next version and re-record whatever the tree holds then.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the follow-up JSON schema `roundfix/follow-up/v1` keeps its
  fields and types. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; the status view reads local Git, Specs
  and the Run Database as today. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0143 (the Release Plan is
  read-only and "names any manual classification or approval still
  required") stays unchanged: the status view only renders the plan it
  reads. ADR-0166 (the Daemon records undeclared paths, and "Recording
  discloses a change after the fact, and it reserves nothing") stays
  unchanged: this Spec documents that a Recorded Path's settlement commit is
  still Wave collision evidence. ADR-0176, which cites ADR-0166, does not
  apply: it governs which Spec text the citation checks read, and this Spec
  changes no citation check. Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — protected paths limited to the `roundfix`
  and `write-tasks` skills. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0323-the-status-the-park-table-and-write-tasks-say-what-to-do-next/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/follow-up.md`,
  `.agents/skills/write-tasks/SKILL.md`, `skills/roundfix/SKILL.md`,
  `skills/write-tasks/SKILL.md`.

## Goals

- In the manual classification state, `roundfix follow-up --status` prints no
  blank field: its `Next release` line says that no version is proposed and
  how many commits since the base tag need a manual release classification,
  and `Waiting on the maintainer` names the `roundfix release plan` command
  that classifies them, in text and JSON.
- Every other release state of the status view keeps its bytes.
- The Park Class table of the `deliver` guide and both Roundfix reference
  copies maps each class to its blockers only, and the `secret-found` next
  action is stated for that blocker alone, as the queue prints it; a
  documentation contract keeps the two from drifting.
- The `write-tasks` skill tells authors to declare a test helper, fake,
  fixture or golden that two Tasks may each change in every such Task, and
  says that a Recorded Path's settlement commit is still collision evidence;
  `CONTEXT.md` **Recorded Path** says the same.

## Core Features

1. Manual classification status. `internal/cli/follow_up_status.go` fills
   `release.reason` with `<n> commit(s) since <base> need a manual release
   classification`, where `<n>` counts the Release Plan's blocking commits,
   and gives the `release_classification` maintainer item the detail
   `classify them with roundfix release plan --from <base> --to <target>
   --impact <none|patch|minor|major> --reason <text>`. The text renders
   `Next release: not proposed (manual_classification_required): <reason>`
   and `  release classification: <detail>`. The follow-up guide and the
   Roundfix `follow-up` reference describe the state.
2. Park actions per blocker. The `finding` row of the Park Class table lists
   its blockers only, and a sentence after the `archived-item-unproven` one
   states the `secret-found` park's next action, in the guide and both
   reference copies. A new documentation contract test compares each such
   sentence with `delivery.ClassifyPark` and refuses a class row that carries
   text after its blocker names.
3. Shared test files in the Task Graph. A `write-tasks` paragraph after the
   Recorded Path paragraph, and a revised **Recorded Path** glossary entry.

## User Experience

The operator runs `roundfix follow-up --status` while the release Pull
Request is open and a commit needs a manual classification. Instead of
`Next release:  from v0.117.0 (manual_classification_required)` the report
reads `Next release: not proposed (manual_classification_required): 1
commit(s) since v0.117.0 need a manual release classification`, and under
`Waiting on the maintainer` it reads `release classification: classify them
with roundfix release plan --from v0.117.0 --to HEAD --impact
<none|patch|minor|major> --reason <text>`. An operator reading a
`run-unresolved` park in the `deliver` guide no longer sees a secret rewrite
next to it. An author decomposing a Spec where two Tasks extend the same test
fake declares it in both and gets `SC-WAVE-COLLISION` from `spec check
--strict` before the Run, with the `needs` edge to add.

## Non-Goals / Out of Scope

- A settlement-time diagnostic for undeclared shared paths. When the second
  Task's settlement commit exists, both Tasks have already run in one Wave
  and the gate's refusal is already decided; a diagnostic there cannot save
  the Run, and the existing `SC-WAVE-COLLISION` summary already names each
  path, its source and the edge. See `_techspec.md` → Decisions.
- Changing the Wave collision rule, its evidence sources, the QA gate's
  precondition or the Spec Consistency settlement baseline.
- Changing the Release Plan, its states, `roundfix release plan` output or
  exit codes, or the follow-up JSON field types: `release.proposed_version`
  stays the empty string in every state that proposes no version.
- Changing any Park Class, blocker or queue-printed next action.

## Success Metrics

1. Success Metric: on a fixture with a base tag and three commits, one
   `feat:` and two ambiguous `chore:`, `follow-up --status` prints
   `Next release: not proposed (manual_classification_required): 2 commit(s)
   since v0.4.1 need a manual release classification` and
   `  release classification: classify them with roundfix release plan
   --from v0.4.1 --to HEAD --impact <none|patch|minor|major> --reason
   <text>`, no `Next release:  ` and no `  release : `; the JSON carries the
   first sentence in `release.reason` and the second in the maintainer
   item's `detail`.
2. Success Metric: the `HEAD at release`, `feature`, `dirty`, `patch`,
   `maintenance` and `no tag` cases of
   `TestFollowUpStatusReportsTheNextRelease` keep their expectations.
3. Success Metric: the Park Class `finding` row of the guide and both
   reference copies lists only blocker names, the `secret-found` sentence
   matches `delivery.ClassifyPark` with `<worktree>` and `<slug>`, and the
   new contract test fails on the starting tree's row.
4. Success Metric: the `write-tasks` skill (canonical and mirror) and
   `CONTEXT.md` **Recorded Path** carry the shared-test-file rule; the
   owned-skill versions and the Behavior Surface record are current.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- ADR-0166, accepted on 2026-09-29: "Declaring a path is still how a Task
  plans its work and reserves the path against parallel Tasks. Recording
  discloses a change after the fact, and it reserves nothing." The
  `write-tasks` rule follows it: only a declaration lets the check order two
  Tasks before they run.
- The Fluxus report of 2026-10-23 (ai-memory message `01a16a54`), recorded
  session evidence quoted in the third adopted Backlog Entry: "QA gate
  recusou SC-WAVE-COLLISION quando task_01 e task_02 da mesma onda editaram
  fakes de teste não declarados (só em Recorded paths); spec check --strict
  não prevê."
- The Command Line Interface Guidelines (https://clig.dev/, read on
  2026-10-26): "Suggest commands the user should run. When several commands
  form a workflow, suggesting to the user commands they can run next helps
  them learn how to use your program and discover new functionality." The
  manual classification status follows it by naming the `roundfix release
  plan` command.

## Unreachable Acceptance

- criterion: after release, the operator's next `roundfix follow-up --status`
  in the manual classification state shows the commit count and the
  classification command and no blank version.
  reason: it needs the released binary and the real repository's release
  range, which no Verification or QA row may read.
  satisfied-by: the operator's follow-up after the release that carries this
  Spec, recorded session evidence.

## Validation

The decisions were checked on 2026-10-26 against these sources.

- Characterization at `aa8b3db7`: `TestFollowUpStatusReportsTheNextRelease`
  pins the manual case with `ProposedVersion == ""`, `Reason == nil` and a
  text that merely contains `manual classification`; the renderer's default
  branch prints `Next release: %s from %s (%s)` and the maintainer line
  `  release %s: %s` with the empty version.
- `internal/delivery/park_class.go` `ClassifyPark` gives the rewrite next
  action only to `secret-found`; `revalidation-failed` and
  `archived-item-unproven` have their own, and the other `finding` blockers
  keep the generic `resolve the blocker, then run roundfix deliver retry
  <slug>`.
- `internal/spec/collision.go` `Collisions` reads each Task's newest
  settlement commit through `git log --all`; the Daemon's Spec Consistency
  settlement check compares with a baseline taken when each Task starts
  (`internal/daemon/task_engine.go`), and `specGatePrecondition` runs the
  strict check without a baseline, so the pair is first refused at the gate.
- A prototype of the three Tasks in a scratch worktree passed
  `./internal/cli -run FollowUp`, `./skills` and the full
  `-tags docscontract ./internal/docscontract`; each sabotage in
  `_techspec.md` → Testing Approach failed its named test.
- Secondbrain, read only: `qmd query "roundfix wave collision undeclared test
  fakes recorded paths QA gate" --all --files --min-score 0.3` returned the
  Roundfix mirror's Archive Records of Spec 0097 ("Tasks the graph declares
  independent edit the same file and die at integration") and Spec 0297, and
  the `qa-gate` skill; none records a decision against declaring shared test
  files. `qmd query "roundfix follow-up status next release manual
  classification"` returned the Roundfix mirror's `docs/user-guide/usage.md`
  and the `roundfix` skill, and `qmd query "roundfix deliver park class
  secret-found next action"` returned the mirror's `CONTEXT.md` and
  `CHANGELOG.md`; none prescribes the blank version or the joined row.
- Exa: the Command Line Interface Guidelines, quoted under Acceptance
  evidence.

## Glossary

- changes: **Recorded Path**

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AQ batch 2 brief for Spec 0323 of 2026-10-26, the maintainer's cycle AQ answers quoted in `_authorization.md`, and the characterization above.

Taken from those sources, with no further interview:

- The manual case keeps `release.proposed_version` empty, as every state
  that proposes no version does, and fills the nullable `release.reason`
  instead of changing a field's type.
- The Park Class table lists blockers only; a blocker-specific next action is
  a sentence of its own, as `archived-item-unproven` already is.
- For the third source, guidance in `write-tasks` (option a), not a
  settlement-time diagnostic (option b).
- No new ADR: the number reserved for this Spec stays unused.

## Open Questions

None.
