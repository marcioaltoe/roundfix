---
task: task_03
spec: 0333-runs-load-the-dispatched-skills-and-the-baseline-follows-the-catalog
status: pending
type: backend
complexity: high
---

# Task 03: Run Agents load only the skills the dispatch selects and name them in their Result

## Overview

A Run Agent never reads the skill dispatch guide today, so the skills it
makes mandatory never load: over 60 days, 108 Claude sessions of an
adopter's frontend Tasks loaded no frontend skill. This Task implements
ADR-0346. The Daemon lists `docs/agents/skill-dispatch.md` in the Spec
Context Bundle; the Task prompt tells the Agent to load only the skills whose
triggers match its Task and to write the Loaded Skill List in `## Result`;
the QA contract gives the QA Agent the same rule without a list;
`implement-task` and `qa-gate` say the same; the guide becomes a
carry-forward input. The reconcile guides, the Roundfix skill and the
glossary follow.

## Requirements

1. MUST add `skillDispatchGuidePath` to `internal/daemon/task_context.go`
   and list the guide in `assembleTaskContextBundle` as `_techspec.md` → API
   Contract 5 states; `buildQAPromptContext` keeps reusing that bundle.
2. MUST add the two invariants and the QA rule of `_techspec.md` → Data Model
   4, verbatim, to `taskExecutionInvariants` and `qaGateContract` in
   `internal/agent/spec_prompt.go` (API Contract 4); every other line keeps
   its bytes.
3. MUST append `docs/agents/skill-dispatch.md` to the standard inputs of
   `CarryForwardInputs` in `internal/spec/spec.go` right after
   `.agents/skills/implement-task/SKILL.md` and update its doc comment (API
   Contract 6); update the expected list in `internal/spec/spec_test.go`.
4. MUST add the tests `_techspec.md` → Testing Approach names for task_03:
   `TestBuildTaskPromptTellsTheAgentToLoadOnlyTheDispatchedSkills` and
   `TestBuildQAPromptTellsTheGateToLoadItsDispatchedSkills` in
   `internal/agent/spec_prompt_test.go`, each asserting its lines appear
   exactly once and, for the Task prompt, before `Spec Context Bundle:`;
   `TestAssembleTaskContextBundleListsTheSkillDispatchGuideWhenPresent` in
   `internal/daemon/task_context_test.go`, covering a present guide (listed
   right after the `implement-task` skill) and an absent one (not listed);
   and a `skill dispatch guide` case in the table of
   `TestCarryForwardRefusesATaskWhoseBindingTextChanged` in
   `internal/cli/carryforward_amendment_test.go`, whose fixture
   (`newPlanningAmendmentFixture`) also writes the guide.
5. MUST add, in the canonical skills, the section "## Skills for the Task"
   before "## 2. Plan" of `.agents/skills/implement-task/SKILL.md` and the
   section "### Skills the gate loads" before "## 1. Resolve scope and
   preconditions" of `.agents/skills/qa-gate/SKILL.md` (API Contract 7),
   never touching `### QA settlement`. The `implement-task` section MUST
   contain `docs/agents/skill-dispatch.md`, `Skills loaded: <names>` and
   `Skills loaded: none matched`; the `qa-gate` section MUST contain
   `docs/agents/skill-dispatch.md`.
6. MUST update the canonical `.agents/skills/roundfix/SKILL.md` bullet on the
   Spec Context Bundle ("The bundle includes standard Spec artifact paths,
   root instructions, the canonical implement-task skill path,") so it also
   names the skill dispatch guide when the repository has it and the load
   rule, containing `skill dispatch guide`; and the sentence "The Task file,
   `AGENTS.md`, `CONTEXT.md`, the `implement-task` skill and every explicit
   Context path still refuse on any change." in
   `docs/user-guide/commands/reconcile.md` and the canonical
   `.agents/skills/roundfix/references/reconcile.md` so it names
   `docs/agents/skill-dispatch.md` after the `implement-task` skill.
7. MUST, through `domain-modeling`, add **Loaded Skill List** to
   `CONTEXT.md` (the `Skills loaded:` line an Agent writes in a Task's
   `## Result` after loading the skills the dispatch selects; never checked
   by the Daemon), with an `_Avoid_` line and a link to ADR-0346, and revise
   **Spec Context Bundle** so it names the skill dispatch guide. No entry
   references a Spec.
8. MUST copy the canonical skills to `skills/` with `make skills-sync`, then
   raise both version fields of `implement-task`, `qa-gate` and `roundfix`
   and their mirrors with `go test ./skills -run
   '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
   -count=1`, recording `skills/testdata/owned-skill-versions.json`.
9. MUST re-record the Behavior Surface fingerprint as the last step, after
   the last guide edit: `go test -tags docscontract ./internal/docscontract
   -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
   -count=1`, recording `docs/references/behavior-surfaces.json`.

## Subtasks

- [ ] Add the guide to the bundle and the carry-forward inputs.
- [ ] Add the prompt invariants and the QA rule, with their tests.
- [ ] Update `implement-task`, `qa-gate`, the Roundfix skill and the reconcile guides; sync and record the skills.
- [ ] Add and revise the glossary terms; re-record the Behavior Surface fingerprint last.

## Acceptance Criteria

- [ ] A Task prompt in a repository with the guide lists it and tells the Agent to load only matching skills and write the Loaded Skill List.
- [ ] A QA prompt tells the gate to load the QA activation's and its surfaces' skills.
- [ ] A changed dispatch guide refuses a carry-forward like a changed `AGENTS.md`.
- [ ] `implement-task`, `qa-gate`, the Roundfix skill, the reconcile guides and the glossary state the rule.

## Context

- instruction: `docs/adr/0346-a-run-agent-loads-the-skills-the-dispatch-selects-for-its-task.md`
- instruction: `docs/adr/0284-carry-forward-binds-a-task-to-the-planning-text-it-references.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- interface: `internal/daemon/task_context.go`
- interface: `internal/daemon/task_context_test.go`
- interface: `internal/agent/spec_prompt.go`
- interface: `internal/agent/spec_prompt_test.go`
- interface: `internal/spec/spec.go`
- interface: `internal/spec/spec_test.go`
- interface: `internal/cli/carryforward_amendment_test.go`
- interface: `.agents/skills/implement-task/SKILL.md`
- interface: `skills/implement-task/SKILL.md`
- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `skills/qa-gate/SKILL.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/reconcile.md`
- interface: `skills/roundfix/references/reconcile.md`
- interface: `docs/user-guide/commands/reconcile.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `CONTEXT.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The Task and QA prompts and the Spec Context Bundle — consumers: `./internal/agent`, `./internal/daemon`
- The carry-forward input set — consumers: `./internal/spec`, `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v ./internal/agent ./internal/spec 2>&1 && go test -count=1 -timeout 20m -v -run '^(TestAssembleTaskContextBundle.*|TestTaskCyclePromptContainsBundleWithoutReferencedBodies)$' ./internal/daemon 2>&1 && go test -count=1 -timeout 20m -v -run '^TestCarryForwardRefusesATaskWhoseBindingTextChanged$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestBuildTaskPromptTellsTheAgentToLoadOnlyTheDispatchedSkills TestBuildQAPromptTellsTheGateToLoadItsDispatchedSkills TestAssembleTaskContextBundleListsTheSkillDispatchGuideWhenPresent TestAssembleTaskContextBundleReservesExplicitPathsAndCountsOmittedPriorFiles 'TestCarryForwardRefusesATaskWhoseBindingTextChanged/skill_dispatch_guide'; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < .agents/skills/implement-task/SKILL.md | grep -qF -- '## Skills for the Task' || { echo 'implement-task lacks the skills section' >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/implement-task/SKILL.md | grep -qF -- 'Skills loaded: none matched' || { echo 'implement-task lacks the empty list form' >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/implement-task/SKILL.md | grep -qF -- 'Skills loaded: <names>' || { echo 'implement-task lacks the list form' >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/implement-task/SKILL.md | grep -qF -- 'docs/agents/skill-dispatch.md' || { echo 'implement-task lacks the guide' >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/qa-gate/SKILL.md | grep -qF -- '### Skills the gate loads' || { echo 'qa-gate lacks the skills section' >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/SKILL.md | grep -qF -- 'skill dispatch guide' || { echo 'roundfix skill lacks the guide' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/reconcile.md | grep -qF -- 'docs/agents/skill-dispatch.md' || { echo 'reconcile guide lacks the guide input' >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/reconcile.md | grep -qF -- 'docs/agents/skill-dispatch.md' || { echo 'reconcile reference lacks the guide input' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Loaded Skill List**:' || { echo 'CONTEXT.md lacks Loaded Skill List' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Spec Context Bundle**:' || { echo 'CONTEXT.md lacks Spec Context Bundle' >&2; exit 1; }; awk '/^[*][*]Spec Context Bundle[*][*]:/ { getline; if (index($0, "skill dispatch guide") > 0) found=1 } END { exit(found ? 0 : 1) }' CONTEXT.md || { echo 'Spec Context Bundle is not revised' >&2; exit 1; }; cmp .agents/skills/implement-task/SKILL.md skills/implement-task/SKILL.md && cmp .agents/skills/qa-gate/SKILL.md skills/qa-gate/SKILL.md && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && cmp .agents/skills/roundfix/references/reconcile.md skills/roundfix/references/reconcile.md && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the three new tests and the `skill dispatch guide` case do not exist, so the first is missing from the passes. Afterwards the prompt, bundle and carry-forward tests pass, the existing bundle test still counts its reserved paths, the skills, guides and glossary carry the rule, the mirrors equal their canonical copies, and the owned-skill version and documentation contracts pass.

## References

- `_prd.md` → Goals; Core Features 1 and 5; Stories 1 and 2; Success Metric 3; Success Metric 6; Glossary; Decisions
- `_techspec.md` → Data Model 4; API Contracts 4 to 7; Testing Approach; Build Order 3; Glossary
- ADR-0346; ADR-0284
