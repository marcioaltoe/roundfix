---
task: task_04
spec: 0397-each-phase-of-a-spec-is-recorded-and-starts-in-a-fresh-session
status: pending
type: docs
complexity: medium
---

# Task 04: The skills and the autonomous-work guide state the Phase Boundary and keep a fix in the session

## Overview

The four authoring skills gain one section, `## Phase boundary`, that
records the phase and says what a session and an Author Subagent do when the
phase ends, and `write-techspec` and `write-tasks` keep a fix or a refactor
in the session unless the user asks for the Author Subagent. The Roundfix
skill's `modes` reference records the `follow-up`, `delivery` and `release`
phases and states the boundary rule, and one clause of the Baseline's
autonomous-work rule carries the rule to every adopter's guide. Contract
tests pin each new sentence and each guarantee the changed places already
gave, in both copies of every skill file and in the rendered guide. The Task
writes text only: tasks 01 to 03 built everything it describes.

## Requirements

1. MUST insert, in the canonical copy of each of
   `.agents/skills/write-idea/SKILL.md`, `.agents/skills/write-prd/SKILL.md`,
   `.agents/skills/write-techspec/SKILL.md` and
   `.agents/skills/write-tasks/SKILL.md`, this section, word for word,
   immediately before the line `## Process` and followed by one blank line,
   with `<phase>` replaced by the skill's own phase (`idea`, `prd`,
   `techspec`, `tasks`) and `<slug>` left as written:

   > ## Phase boundary
   >
   > Once the Spec's slug is known, record the phase: run `roundfix phase start <phase> --spec <slug>` in the session or the Author Subagent that runs it. When the installed Roundfix does not know the command, say so in the report and continue.
   >
   > When the phase ends, a Phase Boundary follows:
   >
   > - **You are the phase's Author Subagent:** return your report and end no session; the session that started you applies the rule.
   > - **Otherwise** write the Session Handoff with the Spec, the phase that ended, the next phase, each decision still open and what stays armed, tell the user what stays armed, and start no further phase in this session. The next session starts with the Roundfix skill's `follow-up` Skill Mode, which claims the handoff.

   MUST anchor on the heading, never on a line number, MUST NOT change any
   other sentence of the four files outside requirement 2, and MUST NOT add
   a `model` or `effort` key to any frontmatter.
2. MUST append, in the canonical `write-techspec` and `write-tasks` only, to
   the bullet of `## Declared model` that begins "**You activated this skill
   yourself**", after one space, this sentence, word for word:

   > On the fix or refactor route, which a Spec whose `type` is `fix` or `refactor` takes, run the phase here instead, on the session's model, unless the user's request or one of the user's answers asks for the Author Subagent, and say in the report where the phase ran.

   Every sentence the section already holds stays, so each of the thirteen
   phrases `TestTheAuthoringSkillsDelegateADeclaredPhase` pins still occurs
   exactly once. `write-idea` and `write-prd` do not gain the sentence.
3. MUST append in the canonical
   `.agents/skills/roundfix/references/modes.md`, word for word, each after
   one space at the end of the item it names, and change no existing
   sentence:
   - item 1 of `--follow-up`:

     > Before the inventory, record the phase with `roundfix phase start follow-up`; when the installed Roundfix does not know the command, say so in the report and continue.

   - item 7 of `--follow-up`:

     > When the mode starts to follow a Spec's delivery, run `roundfix phase start delivery --spec <slug>`, and before a release run `roundfix phase start release`, with `--spec <slug>` when the release carries one Spec.

   - item 8 of `--follow-up`:

     > At a Phase Boundary, after each Authoring Phase and after a Spec is delivered and released, write the Session Handoff with the Spec, the phase that ended, the next phase, each decision still open and what stays armed, tell the user what stays armed, and start no further phase in this session; the next session starts with this mode, which claims the handoff. A scheduled wake-up returns to the same session with its whole context, so it keeps a phase going and never crosses a Phase Boundary. A fix or a refactor is authored in the session, as the `Declared model` section of `write-techspec` and `write-tasks` says.

   and add to the `--unattended` section two paragraphs of their own, the
   first right before the paragraph that begins "When run alone, arm the
   continuation" and the second at the end of the section:

   > This mode asks no question about where a fix or a refactor is authored: it is authored in the session, and its Author Subagent starts only when the user's request or one of the user's answers asks for it.

   > With nobody present the session still ends at a Phase Boundary and says what stays armed and what waits. A scheduled wake-up keeps a phase going and never crosses a Phase Boundary, because it returns to the same session.

   The `--status` section does not change: that mode records no phase. MUST
   NOT add or change text inside any `### QA settlement` section.
4. MUST add, in a new test file of the `skills` package, two tests that call
   `t.Parallel()` first and read both copies through `roundfixLayoutSources`,
   comparing whitespace-normalized text:
   - `TestTheAuthoringSkillsStateThePhaseBoundary`, with one subtest per copy
     and skill named `<copy>/<skill>` (`canonical/write-idea` to
     `embedded/write-tasks`). Each asserts that the file holds exactly one
     `## Phase boundary` heading, after `## Declared model` and before
     `## Process`, and, in one named subtest each, that the section holds
     exactly once:
     - `records_its_own_phase`: "run `roundfix phase start <phase> --spec <slug>`", with the skill's own phase
     - `an_unknown_command_continues`: "When the installed Roundfix does not know the command, say so in the report and continue."
     - `the_author_subagent_returns`: "return your report and end no session"
     - `the_session_writes_the_handoff`: "write the Session Handoff with the Spec, the phase that ended, the next phase, each decision still open and what stays armed"
     - `the_user_is_told`: "tell the user what stays armed"
     - `no_further_phase`: "start no further phase in this session"
     - `the_next_session_follows_up`: "The next session starts with the Roundfix skill's `follow-up` Skill Mode, which claims the handoff."
   - `TestTheRoundfixSkillRecordsEachPhaseAndStatesThePhaseBoundary`, with
     one subtest per copy (`canonical`, `embedded`) and, under it, one named
     subtest per phrase, each found exactly once in the named place of
     `modes.md`:
     - `follow_up_records_its_phase`, item 1 of `--follow-up`: "record the phase with `roundfix phase start follow-up`"
     - `delivery_is_recorded`, item 7: "run `roundfix phase start delivery --spec <slug>`"
     - `release_is_recorded`, item 7: "before a release run `roundfix phase start release`"
     - `the_handoff_at_a_boundary`, item 8: "write the Session Handoff with the Spec, the phase that ended, the next phase, each decision still open and what stays armed"
     - `no_further_phase`, item 8: "start no further phase in this session"
     - `the_next_session_claims_the_handoff`, item 8: "the next session starts with this mode, which claims the handoff"
     - `a_wake_up_never_crosses`, item 8: "A scheduled wake-up returns to the same session with its whole context, so it keeps a phase going and never crosses a Phase Boundary."
     - `a_fix_stays_in_the_session`, item 8: "A fix or a refactor is authored in the session"
     - `the_turn_still_arms_or_says`, item 8, a guarantee the item already gave: "arm a detached Run, the Delivery Queue or a scheduled wake-up, or say plainly that work stops and what remains"
     - `unattended_asks_no_question`, the `--unattended` section: "This mode asks no question about where a fix or a refactor is authored"
     - `unattended_ends_at_the_boundary`, the same section: "With nobody present the session still ends at a Phase Boundary"
     - `unattended_wake_up`, the same section: "A scheduled wake-up keeps a phase going and never crosses a Phase Boundary, because it returns to the same session."
     - `status_records_no_phase`: the `--status` section does not contain `roundfix phase start`

   MUST also extend `TestTheAuthoringSkillsDelegateADeclaredPhase`, under
   each `<copy>/<skill>` subtest, with the named subtests
   `a_fix_or_refactor_stays_in_the_session` ("On the fix or refactor route,
   which a Spec whose `type` is `fix` or `refactor` takes, run the phase
   here instead, on the session's model"), `the_user_asks_for_the_author`
   ("unless the user's request or one of the user's answers asks for the
   Author Subagent") and `the_report_says_where` ("say in the report where
   the phase ran"): each phrase occurs exactly once in the `## Declared
   model` section of `write-techspec` and `write-tasks` and never in
   `write-idea` or `write-prd`. Its thirteen phrases and every other
   existing test stay as they are; no recorded test is removed or renamed.
5. MUST then run, in this order and never by hand: `make skills-sync`, and
   the version record
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   which raises the versions of `write-idea`, `write-prd`, `write-techspec`,
   `write-tasks` and `roundfix` from the tree the Task runs on; never write
   a number. The five generated files under `.claude/agents/` do not
   change, and no installation is run.
6. MUST change, in `internal/baseline/assets/modules/autonomous-work.json`,
   only the `guidance` of `clause.autonomous.loop-07-outlive-the-turn`, by
   appending after its last sentence, after one space, the two sentences
   below, word for word, and raise the `version` of `rule.autonomous.loop`
   by one from the value the Task finds. MUST keep the clause's id, its
   `enforcement: mandatory`, its absence of `replaces`, its three existing
   sentences, the file's two-space JSON formatting and every other clause
   (Invariant 8):

   > At a phase boundary of a Spec, after each authoring phase and after a Spec is delivered and released, write the session handoff with the Spec, the phase that ended, the next phase, each open decision and what stays armed, tell the user what stays armed, and start no further phase in that session. A scheduled wake-up returns to the same session with its whole context, so it keeps a phase going and is never the way across a boundary.

7. MUST add `internal/baseline/phase_boundary_clause_test.go` with
   `TestTheOutliveTheTurnClauseStatesThePhaseBoundary`, calling
   `t.Parallel()` first and built like `idea_stage_route_clauses_test.go`:
   one named subtest for each phrase below, each asserting that the phrase
   occurs once in the clause, that the clause is `mandatory` and carries no
   `replaces`, and that the phrase occurs once, whitespace-normalized, in
   the Standard TypeScript golden of the autonomous-work guide and in this
   repository's `docs/agents/autonomous-work.md`. Kept: "must arm something
   that outlives its own turn before that turn ends: a detached Run, a
   scheduled wake-up, or a daemon"; "Never end a turn stating that work will
   continue"; "When no such mechanism is armed, say plainly that work is
   stopping and what remains." New: "At a phase boundary of a Spec, after
   each authoring phase and after a Spec is delivered and released"; "write
   the session handoff with the Spec, the phase that ended, the next phase,
   each open decision and what stays armed"; "tell the user what stays
   armed"; "start no further phase in that session"; "A scheduled wake-up
   returns to the same session with its whole context, so it keeps a phase
   going and is never the way across a boundary."
8. MUST regenerate, in this order and never by hand: the module record with
   `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`
   (it writes the next version of `autonomous-work`; the Spec names no
   number); the derived Baseline files with `make baseline-digests`; and, as
   the Task's last step, this repository's managed guide and setup manifest
   with
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`,
   whose second run reports no change. Every `baseline update` the Task
   runs carries `--no-skills` and is the product built from this tree, never
   a `roundfix` found on `PATH`. The derived files are measured on the tree
   the Task runs on: the files the prototype rewrote are declared below, and
   any other file those commands rewrite is recorded by the Daemon.
9. MUST change in `CONTEXT.md`, through `domain-modeling`, only the
   definition lines of **Author Subagent** and **Session Handoff**, where
   each `(ADR-NNNN)` is the glossary's Markdown link to that ADR's file
   under `docs/adr/` as the entries write their links today. **Author
   Subagent** becomes:

   > The subagent that installation generates for one declared Authoring Phase, `roundfix-<phase>-author`, on the declared model and reasoning effort. An authoring skill the session activated hands its phase to it, passing no model and no effort; it cannot ask the user, so the session gathers the decisions first and it returns its open questions instead of guessing one. On the fix or refactor route, `write-techspec` and `write-tasks` run the phase in the session and start it only when the user's request or one of the user's answers asks for it. It is a subagent of the Agent harness session and never the Agent of a Run, and it ends no session at a Phase Boundary (ADR-0401, ADR-0408).

   **Session Handoff** keeps every word of its definition line and gains,
   after one space at its end:

   > At a Phase Boundary the Supervisor writes one with the Spec, the phase that ended, the next phase, each decision still open and what stays armed (ADR-0407).

10. MUST NOT change any Source Baseline file (the corpus, manifest and index
    of `baseline.standard-typescript-monorepo-0.0.1` stay byte-identical),
    any other module, template or profile source by hand, any existing
    Baseline test, any guide under `docs/user-guide/`, any Go source outside
    the three test files, `.roundfixrc.yml`, any file under
    `.claude/agents/`, or any other skill. Of the Governed Paths the Task
    may change only the canonical `SKILL.md` of `write-idea`, `write-prd`,
    `write-techspec`, `write-tasks` and `roundfix` with their five mirrors,
    `.agents/skills/roundfix/references/modes.md`,
    `internal/baseline/assets/modules/autonomous-work.json`,
    `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
    `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md`,
    `docs/agents/autonomous-work.md` and `docs/agents/setup-context.json`,
    the last four through the commands of requirement 8.

## Subtasks

- [ ] Insert the `## Phase boundary` section into the four authoring skills and the fix-route sentence into two of them.
- [ ] Append the sentences and paragraphs to the `modes` reference.
- [ ] Add the two contract tests and extend the delegation test.
- [ ] Sync the mirrors and record the five skill versions.
- [ ] Extend the Baseline clause, add its phrase test, record the module version and regenerate the derived files.
- [ ] Change the two glossary entries.
- [ ] Refresh this repository's guide and setup manifest, and confirm a second refresh changes nothing.

## Acceptance Criteria

- [ ] Both copies of each authoring skill carry one `## Phase boundary` section between `## Declared model` and `## Process` that records the skill's own phase, sends an Author Subagent back with its report and tells a session to write the Session Handoff and start no further phase; no shipped skill declares a model or an effort.
- [ ] Both copies of `write-techspec` and `write-tasks` keep a fix or a refactor in the session unless the user asks for the Author Subagent, `write-idea` and `write-prd` do not say so, and every phrase the delegation test pinned still occurs once.
- [ ] Both copies of the `modes` reference record `follow-up`, `delivery` and `release`, state the boundary rule and the wake-up rule in item 8 and in `--unattended`, and keep what item 8 guaranteed; the `--status` section records no phase.
- [ ] A resolved profile's autonomous-work guide carries the boundary rule and the wake-up rule beside the clause's three earlier sentences; the clause keeps its id and force; the derived Baseline files and this repository's guide and setup manifest are current, and a second managed refresh reports no change.
- [ ] `CONTEXT.md` carries the changed **Author Subagent** and **Session Handoff**, and the five skill versions are recorded.

## Context

- instruction: `docs/adr/0407-a-phase-record-names-the-spec-and-phase-of-supervisor-work-and-a-phase-boundary-ends-the-session.md`
- instruction: `docs/adr/0408-a-repository-that-publishes-an-authoring-skill-receives-no-header-for-it-and-a-fix-is-authored-in-the-session.md`
- instruction: `docs/adr/0338-the-follow-up-mode-claims-the-handoff-and-leaves-no-handed-in-message-pending.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `skills/roundfix_layout_test.go`
- instruction: `internal/baseline/idea_stage_route_clauses_test.go`
- interface: `.agents/skills/write-idea/SKILL.md`
- interface: `skills/write-idea/SKILL.md`
- interface: `.agents/skills/write-prd/SKILL.md`
- interface: `skills/write-prd/SKILL.md`
- interface: `.agents/skills/write-techspec/SKILL.md`
- interface: `skills/write-techspec/SKILL.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/modes.md`
- interface: `skills/roundfix/references/modes.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `skills/authoring_delegation_test.go`
- creates: `skills/phase_boundary_test.go`
- interface: `internal/baseline/assets/modules/autonomous-work.json`
- interface: `internal/baseline/module-versions.json`
- creates: `internal/baseline/phase_boundary_clause_test.go`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/autonomous-work.md`
- interface: `docs/agents/setup-context.json`
- interface: `CONTEXT.md`

## Consumers

- The authoring skills' and the `modes` reference's text and the recorded skill versions — consumers: `./skills`
- The extended clause and the derived catalog files — consumers: `./internal/baseline`, `./internal/cli`
- This repository's managed guide and setup manifest and the glossary entries — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheAuthoringSkillsStateThePhaseBoundary|TestTheRoundfixSkillRecordsEachPhaseAndStatesThePhaseBoundary|TestTheAuthoringSkillsDelegateADeclaredPhase|TestTheFollowUpModeRunsItsStepsInOrder|TestTheFollowUpModeTakesTheHandoffAndSettlesEveryMessage|TestEveryOwnedSkillVersionIsRecorded|TestNoShippedSkillDeclaresAModelOrEffort)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheAuthoringSkillsStateThePhaseBoundary TestTheAuthoringSkillsStateThePhaseBoundary/canonical/write-idea/records_its_own_phase TestTheAuthoringSkillsStateThePhaseBoundary/embedded/write-tasks/no_further_phase TestTheRoundfixSkillRecordsEachPhaseAndStatesThePhaseBoundary TestTheRoundfixSkillRecordsEachPhaseAndStatesThePhaseBoundary/canonical/a_wake_up_never_crosses TestTheRoundfixSkillRecordsEachPhaseAndStatesThePhaseBoundary/embedded/the_turn_still_arms_or_says TestTheAuthoringSkillsDelegateADeclaredPhase TestTheAuthoringSkillsDelegateADeclaredPhase/canonical/write-techspec/a_fix_or_refactor_stays_in_the_session TestTheAuthoringSkillsDelegateADeclaredPhase/embedded/write-prd/a_fix_or_refactor_stays_in_the_session TestTheFollowUpModeRunsItsStepsInOrder TestTheFollowUpModeTakesTheHandoffAndSettlesEveryMessage TestEveryOwnedSkillVersionIsRecorded TestNoShippedSkillDeclaresAModelOrEffort; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; for skill in roundfix write-idea write-prd write-techspec write-tasks; do diff -r ".agents/skills/$skill" "skills/$skill" >/dev/null || { echo "mirror differs: $skill" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./skills` — expected: exit 0. Before this Task the two contract tests and the three new subtests do not exist, so their pass lines are missing; afterwards both copies of every changed skill file carry each phrase once, the fix-route sentence is in two skills and absent from the other two, the mirrors equal the canonical copies, the five versions are recorded, and the skills suite passes.

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheOutliveTheTurnClauseStatesThePhaseBoundary|TestEveryBaselineModuleVersionIsRecorded)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheOutliveTheTurnClauseStatesThePhaseBoundary TestEveryBaselineModuleVersionIsRecorded; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m ./internal/baseline && go test -count=1 -timeout 40m -run 'Baseline|Guidance|SetupManifest' ./internal/cli` — expected: exit 0. Before this Task the phrase test does not exist, so its pass line is missing; afterwards every kept and new phrase renders once, the clause keeps its force, the module version is recorded, and the Baseline suite and the CLI's Baseline tests pass.

- `tr -s '[:space:]' ' ' < docs/agents/autonomous-work.md | grep -qF -- 'A scheduled wake-up returns to the same session with its whole context, so it keeps a phase going and is never the way across a boundary.' || { echo 'the autonomous-work guide lacks the boundary rule' >&2; exit 1; }; grep -F -- '**Author Subagent**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Author Subagent\*\*(:| —( |$))' || { echo 'missing term: Author Subagent' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'never the Agent of a Run, and it ends no session at a Phase Boundary' || { echo 'Author Subagent unchanged' >&2; exit 1; }; grep -F -- '**Session Handoff**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Session Handoff\*\*(:| —( |$))' || { echo 'missing term: Session Handoff' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'At a Phase Boundary the Supervisor writes one with the Spec, the phase that ended, the next phase, each decision still open and what stays armed' || { echo 'Session Handoff unchanged' >&2; exit 1; }; d="$(mktemp -d)" || exit 1; trap 'rm -rf "$d"' EXIT; go build -buildvcs=false -o "$d/roundfix" ./cmd/roundfix || exit 1; mkdir "$d/home" || exit 1; out="$(HOME="$d/home" "$d/roundfix" baseline update --repo . --no-skills --no-history --format json)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q '"state":"current"' || { printf '%s\n' "$out"; exit 1; }; go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the guide holds no boundary rule; afterwards this repository's guide carries it, the two glossary entries are changed, a managed refresh by the product built from the tree reports `current`, and the documentation contracts, the setup manifest included, pass.

## References

- `_prd.md` → User Stories 4, 5 and 9; Core Features 8, 9, 12 and 13; Success Metric 8; Glossary
- `_techspec.md` → Current behavior; API Contract 9; API Contract 10; API Contract 11; Invariants 5, 7 and 8; Declared breaks; Testing Approach; Build Order 4
- ADR-0407, ADR-0408, ADR-0338
