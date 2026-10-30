---
spec: 0333-runs-load-the-dispatched-skills-and-the-baseline-follows-the-catalog
prd: _prd.md
created: 2026-10-30
---

# Runs load the dispatched skills, and the Baseline follows the curated catalog — Technical Spec

## Executive Summary

Four changes, in order. First, the asset sync moves every upstream Setup
Snapshot to catalog commit `86df8608726979064a8fe46003dfaf403ffe0087`, the
seven skills the catalog removed join the Retired Skills, every module,
trigger and bundle that names them is cleaned, the TypeScript snapshot's
activation bundles are set from the catalog file before the sync, and this
repository restores its three changed upstream skills. Second, five triggers
are reworded so a Run never matches a person-started or polish skill, and the
product UI and public pages bundles split. Third, the Daemon lists
`docs/agents/skill-dispatch.md` in the Spec Context Bundle, the Task prompt,
the QA contract, `implement-task` and `qa-gate` state the load rule, and the
dispatch guide joins the carry-forward inputs (ADR-0346). Fourth,
`implement-spec` becomes a Retired Skill and leaves the binary, the module,
the `Makefile`'s owned list and the skills that name it; its two unique rules
move into the Roundfix skill's delivery reference.

The trade-off accepted: the load rule spends one guide read and the matched
skills' bodies per Task, and trusts the Agent's own Loaded Skill List instead
of a Daemon check; the only proof that it works is the recorded sessions
after the release.

## Project Constraints

- Identifier strategy: not applicable — no project-owned identifier or record
  schema is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path or credential.
  Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — Build Order 3 implements ADR-0346 and
  follows ADR-0284 and ADR-0170; Build Orders 1 and 4 follow ADR-0246 and ADR-0206. All
  are named with their titles in `_prd.md` → Project Constraints. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AS answers of 2026-10-30 and the Governed Paths this Spec declares; the
  maintainer's express grant of 2026-10-30 for the `Makefile`'s
  `OWNED_SKILLS` line only, quoted in the record's Makefile grant section. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0333-runs-load-the-dispatched-skills-and-the-baseline-follows-the-catalog/_authorization.md`;
  bounded files:
  `.agents/skills/archive-spec/SKILL.md`,
  `.agents/skills/context7-cli/references/docs.md`,
  `.agents/skills/context7-cli/SKILL.md`,
  `.agents/skills/grilling/SKILL.md`,
  `.agents/skills/implement-spec/SKILL.md`,
  `.agents/skills/implement-task/SKILL.md`,
  `.agents/skills/qa-gate/SKILL.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `.agents/skills/roundfix/references/reconcile.md`,
  `.agents/skills/roundfix/references/setup.md`,
  `.agents/skills/roundfix/references/spec-delivery.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/tui-design/references/evidence.md`,
  `.agents/skills/tui-design/SKILL.md`,
  `.agents/skills/write-tasks/SKILL.md`,
  `Makefile`,
  `docs/agents/setup-context.json`,
  `docs/agents/skill-dispatch.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`,
  `internal/baseline/assets/modules/backend.json`,
  `internal/baseline/assets/modules/context-workflow.json`,
  `internal/baseline/assets/modules/frontend.json`,
  `internal/baseline/assets/modules/typescript.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/baseline/assets/setups/go-cli-typescript-bun.json`,
  `internal/baseline/assets/setups/go.json`,
  `internal/baseline/assets/setups/python.json`,
  `internal/baseline/assets/setups/rust.json`,
  `internal/baseline/assets/setups/typescript.json`,
  `internal/baseline/assets/skill-activations.json`,
  `skills/archive-spec/SKILL.md`,
  `skills/implement-spec/SKILL.md`,
  `skills/implement-task/SKILL.md`,
  `skills/qa-gate/SKILL.md`,
  `skills/roundfix/SKILL.md`,
  `skills/write-tasks/SKILL.md`.

## Current behavior

Characterized at `35de7f74` (v0.126.0).

- `internal/daemon/task_context.go` `assembleTaskContextBundle` reserves the
  PRD, TechSpec and Task Graph paths, then `AGENTS.md` and
  `.agents/skills/implement-task/SKILL.md` under Instructions when present,
  then the Task's Context paths, then prior changed files, within a 200-path
  cap. `buildQAPromptContext` reuses the same bundle for the QA gate.
- `internal/agent/spec_prompt.go` `taskExecutionInvariants` holds nine
  invariants; none names a skill or the dispatch guide. `qaGateContract`
  holds ten rules; none names a skill other than `qa-gate`.
- `internal/spec/spec.go` `CarryForwardInputs` returns the Task file, the
  three planning files, `AGENTS.md`, `CONTEXT.md`,
  `.agents/skills/implement-task/SKILL.md` and the Task's Context paths;
  `inspectCarryForwardCandidate` (`internal/cli/carryforward.go`) refuses a
  candidate whose existing input changed since its settlement parent.
- `internal/baseline/retired_skills.go` holds `council` and `the-fool`;
  `SyncAssets` drops a Retired Skill from every snapshot it writes, even a
  repository-owned entry; `InstalledRetiredSkills` reports a held copy under
  `Skills retired` and never deletes it (ADR-0246: "Roundfix never deletes an
  installed skill tree").
- The five upstream snapshots name `b56fc5ff`. `typescript.json` and the
  composed `go-cli-typescript-bun.json` carry ten activation bundles that the
  sync copies from the recorded snapshot, never from
  `skill-activations.json`; `go`, `python` and `rust` carry two.
- `modules/typescript.json` requires and dispatches `diagnosing-bugs`,
  `game-changing-features`, `refactoring-analysis`, `tdd` and `to-prompt`;
  `modules/backend.json` `design-patterns`; `modules/frontend.json`
  `fixing-metadata`; `skill-activations.json` puts `diagnosing-bugs` in
  `bundle.debugging`, `tdd` in `bundle.testing`, and six skills in
  `bundle.frontend-ui-quality`.
- `modules/context-workflow.json` requires and dispatches `implement-spec`
  with two triggers; `implement-spec` is owned (`skills/skills.go`,
  `assetsSyncRepoOwnedSkills`, the `Makefile`'s `OWNED_SKILLS`), carries
  `disable-model-invocation: true`, and is named by `archive-spec` (lines 3
  and 18) and `write-tasks` (line 16).
- This repository's `docs/agents/skill-dispatch.md` dispatches
  `implement-spec` and `grill-with-docs` ("Stress-testing a plan or decision
  against repository documents.").

## System Architecture

No new package. The changes extend the Baseline catalog and asset sync
(`internal/baseline`), the Task prompt builder (`internal/agent`), the
Daemon's bundle assembly (`internal/daemon`), the carry-forward input set
(`internal/spec`), the owned skill bundle (`skills`, `.agents/skills`) and
the guides.

## Implementation Design

### Interfaces

```go
// internal/daemon/task_context.go
const skillDispatchGuidePath = "docs/agents/skill-dispatch.md"
// assembleTaskContextBundle adds it to Instructions right after the
// implement-task skill path, only when the file exists in the work dir.

// internal/baseline/retired_skills.go
var retiredSkills = map[string]struct{}{ /* council, design-patterns,
	diagnosing-bugs, fixing-metadata, game-changing-features, implement-spec,
	refactoring-analysis, tdd, the-fool, to-prompt */ }
```

`CarryForwardInputs` appends `docs/agents/skill-dispatch.md` after
`.agents/skills/implement-task/SKILL.md`. No exported signature changes.

### Data Models

1. Data Model: the Retired Skills after task_01 are `council`,
   `design-patterns`, `diagnosing-bugs`, `fixing-metadata`,
   `game-changing-features`, `refactoring-analysis`, `tdd`, `the-fool` and
   `to-prompt`; task_04 adds `implement-spec`.
2. Data Model: `skill-activations.json` after task_02 (version 3):
   `bundle.debugging` = `systematic-debugging`, `no-workarounds`;
   `bundle.testing` = `testing-boss`, `vitest`;
   `bundle.frontend-public-pages` = `frontend-design`,
   `fixing-accessibility`, activated by `trigger.frontend.public-pages`
   ("Designing or reviewing marketing or other public pages.", owner
   `frontend`); `bundle.frontend-ui-quality` = `interaction-design`,
   `interface-design`, `fixing-accessibility`, its trigger reading
   "Designing or reviewing product UI: application screens, forms, and
   dashboards."; every bundle and activation stays sorted by `id`; the other
   bundles are unchanged. task_01 makes only the debugging and testing edits
   and raises the version to 3; task_02 makes the frontend edits.
3. Data Model: the reworded skill triggers (task_02):
   `trigger.typescript.architectural-analysis` "When the person asks for an
   analysis of the TypeScript repository architecture and dependency
   boundaries."; `trigger.typescript.humanizer` "When the person asks to make
   TypeScript repository prose sound natural and direct.";
   `trigger.typescript.writing-clearly-and-concisely` "When the person asks
   to tighten TypeScript repository prose for clarity and brevity.";
   `trigger.context-workflow.grill-with-docs` "When the person asks to
   stress-test a plan or decision against repository documents.";
   `trigger.frontend.web-design-guidelines` "When the person asks for a
   review of the frontend against web interface guidelines.";
   `trigger.frontend.wcag-audit-patterns` "Running a formal WCAG conformance
   audit of the frontend."; `trigger.frontend.frontend-design` "Designing or
   implementing the visual design of marketing or other public pages.";
   `trigger.frontend.interface-design` "Designing product UI hierarchy and
   controls: application screens, forms, and dashboards." `ui-ux-pro-max`
   keeps its trigger.
4. Data Model: the two prompt invariants added before "Implement only this
   Task's slice" (task_03), verbatim:
   `- Before editing, when the Spec Context Bundle lists docs/agents/skill-dispatch.md, read it, select only the triggers and activation bundles that match this Task's type, its declared paths and the work it asks for, and load each selected skill; never load a skill that matches nothing, and do not load a skill only a person can start.`
   and
   `- Name the skills you loaded in ## Result on one line, Skills loaded: <names>, or Skills loaded: none matched.`
   The QA contract rule added before "Run the qa-gate process for this
   Spec.", verbatim:
   `- Before building the matrix, when the Spec Context Bundle lists docs/agents/skill-dispatch.md, load the skills of its QA activation and of each trigger the Spec's declared surfaces match, and only those.`

### API Contracts

1. API Contract: the catalog import (task_01), from the repository root,
   each command verbatim, stopping on the first failure. The module,
   activation and retired-set edits of Data Models 1 and 2 are applied as
   text keeping each file's layout, and the module record command run,
   before the sync:

   ```bash
   sha=86df8608726979064a8fe46003dfaf403ffe0087
   src="$HOME/dev/skills"
   git -C "$src" cat-file -e "$sha^{commit}" || { echo "stop: $src lacks $sha" >&2; exit 1; }
   tmp="$(mktemp -d)"
   git init --quiet "$tmp/skills"
   git -C "$tmp/skills" fetch --quiet "$src" "$sha"
   git -C "$tmp/skills" checkout --quiet --detach "$sha"
   git -C "$tmp/skills" remote add origin https://github.com/marcioaltoe/skills.git
   go build -buildvcs=false -o "$tmp/roundfix" ./cmd/roundfix
   f=internal/baseline/assets/setups/typescript.json
   jq --indent 2 --slurpfile a internal/baseline/assets/skill-activations.json '.activationBundles = $a[0].bundles' "$f" > "$tmp/ts.json" && mv "$tmp/ts.json" "$f"
   h="$(mktemp -d)"; HOME="$h" "$tmp/roundfix" baseline assets sync --source-dir "$tmp/skills/setups" --format text
   HOME="$h" "$tmp/roundfix" baseline assets sync --source-dir "$tmp/skills/setups" --check --format text
   ```

   The sync reports `skills.setup-snapshot.updated` for the five snapshots
   and the check exits `0`; the sync recomputes each digest and the composed
   snapshot's bundles. Then the parity transform of Spec 0208 with the
   composed-snapshot filter of Spec 0317, run as `jq --indent 2 --slurpfile s
   <(jq -s 'sort_by(.id)[]' internal/baseline/assets/setups/*.json) -f
   <program> internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`
   through a temporary copy, with this program:

   ```jq
   .input.sourceRevision as $rev
   | def row($s): {activationBundles: $s.activationBundles, digest: "", id: $s.id,
       schemaVersion: $s.schemaVersion, version: $s.version,
       source: {path: $s.source.path, ref: $rev, repository: "example/skills", type: "github"},
       skills: [$s.skills[] | if .source.type == "github"
         then {name, path, source: {path: .source.path, ref: $rev, repository: "example/skills", type: "github"}, treeDigest: ""}
         else {minimumVersion, name, path, source: {name: .source.name, type: "repo"}} end]}
       | if .activationBundles == null then del(.activationBundles) else . end;
     .manifest.setups = [$s[] | select(.source.type != "composed") | row(.)]
   | .managedEntryLedger = [.manifest.setups[] | {digest: "", id}]
   | .normalizedOutput.result.findings = [.manifest.setups[] | {
       action: "Review the snapshot diff and run asset validation.",
       code: "skills.setup-snapshot.updated", managedId: ("setup." + .id),
       message: "Setup snapshot was synchronized from the canonical source.",
       path: ("assets/setups/" + .id + ".json"), severity: "info"}]
   | .normalizedOutput.result.summary.info = (.manifest.setups | length)
   ```

   Then `make baseline-digests`; then, with a fresh `HOME` each time,
   `"$tmp/roundfix" baseline update --repo . --no-skills --no-history --yes
   --format text` twice (the second reports `File changes: 0`), with the
   binary rebuilt after every source edit; then `"$tmp/roundfix" baseline
   skills restore --profile go-cli-tui --skill context7-cli --skill grilling
   --skill tui-design --source-dir "$tmp/skills" --repo . --format text`, its
   `Plan Digest:` line read, and the same command with `--confirm-plan
   <digest>`. Remove the temporary directories afterwards.
2. API Contract: the re-sync of task_02 and task_04. After their module and
   activation edits and the module record command, each runs the
   `typescript.json` bundle assignment, the sync and its check, the parity
   transform, `make baseline-digests` and the two managed refreshes of API
   Contract 1, from a catalog fetched the same way; neither restores a
   skill.
3. API Contract: the snapshots. Every snapshot whose source type is `github`,
   and each of its `github` skills, names
   `86df8608726979064a8fe46003dfaf403ffe0087`; no snapshot names a Retired
   Skill; the `typescript` and `go-cli-typescript-bun` snapshots carry every
   bundle of `skill-activations.json` with the same skills; their
   `ui-ux-pro-max` tree is the one the catalog commit publishes
   (`34f2a8a3be83f0fff1e9af200c09ad33821d0380298ef2f4952da75652d80807`) and
   their `onioncry` tree
   `66b054d73b38954303e64c4dc2c57916aab971dc8b320f76b930e5c3ef99f3d0`.
4. API Contract: the Task prompt. `BuildTaskPrompt` writes the two invariants
   of Data Model 4 inside the execution invariants, before the Spec Context
   Bundle; `BuildQAPrompt` writes the QA rule inside the QA contract. Every
   other line keeps its bytes.
5. API Contract: the Spec Context Bundle. When
   `docs/agents/skill-dispatch.md` exists in the Run Worktree, the bundle's
   Instructions list it right after `.agents/skills/implement-task/SKILL.md`
   (or after `AGENTS.md` when that skill is absent); otherwise the bundle is
   unchanged. It counts against the 200-path cap like every reserved path.
6. API Contract: carry-forward. `CarryForwardInputs` lists
   `docs/agents/skill-dispatch.md` after the `implement-task` skill; a Task
   whose settlement parent held the guide is refused as moved when the guide
   changed, as for `AGENTS.md`; a parent without it ignores it.
7. API Contract: `implement-task` gains the section "## Skills for the Task"
   before "## 2. Plan", and `qa-gate` the section "### Skills the gate
   loads" before "## 1. Resolve scope and preconditions"; neither touches
   `### QA settlement`. The texts state Data Model 4's rules, `implement-task`
   with the `Skills loaded: <names>` and `Skills loaded: none matched` forms.
8. API Contract: the retirement of `implement-spec` (task_04). It leaves
   `skills/skills.go` (embed line, `skillNames`, the comment's count),
   `assetsSyncRepoOwnedSkills`, the `Makefile`'s `OWNED_SKILLS` (and its
   comment's count), `modules/context-workflow.json`, both trees, and
   `skills/testdata/owned-skill-versions.json`, whose key is deleted with
   `jq --indent 2 'del(.skills["implement-spec"])'` because the record
   command cannot drop a skill. `archive-spec` says "a Roundfix Spec Run" in
   place of the `implement-spec` loop at both places; `write-tasks` says
   "the manifest that `implement-task` and the Roundfix Daemon schedule
   from".
9. API Contract: the Roundfix skill's `references/spec-delivery.md` gains
   the section "### Pending Questions and the planning files" before
   "### Recovering without breaking the loop": ask the maintainer only the
   queue's Pending Question, in its own words and with the action it names,
   adding no confirmation, approval, recommendation or second question, and
   continue only after that action; the Supervisor never edits `_prd.md`,
   `_techspec.md` or `_tasks.md` of a Spec it delivers.

### Surface Transcripts

None. No command changes its output format: `baseline update` prints more
names under its existing `Skills retired` lines, and `skills list` prints one
owned skill fewer in the same shape.

### Invariants

- No module, trigger, activation bundle or snapshot names a Retired Skill
  (the catalog tests read the snapshots' bundles too after task_01).
- `baseline update` never deletes an installed copy and never changes its
  state, category or exit code because of one.
- The Daemon never parses the Loaded Skill List; a Task's verdict depends
  only on its declared Verification.
- The bundle stays path-only: no skill body enters the prompt.

## Coverage Map

- Goal 1 → API Contracts 4, 5 and 7 with Data Model 4: the prompt rules, the bundle entry and the two skill sections (task_03)
- Goal 2 → API Contracts 1, 3 and 8 with Data Model 1: the import, the snapshots and the Retired Skills (task_01, task_04)
- Goal 3 → API Contract 2 with Data Models 2 and 3: the reworded triggers and the split frontend bundles (task_02)
- Goal 4 → API Contracts 8 and 9: the owned bundle without `implement-spec` and the Roundfix delivery reference (task_04)
- Story 1 → task_03 with task_02's bundles
- Story 2 → task_03 (only matching skills)
- Story 3 → task_01, task_04 (Retired Skills, `Skills retired`)
- Story 4 → task_02, task_04
- Story 5 → task_04
- Core Feature 1 → task_03; Core Feature 2 → task_01; Core Feature 3 →
  task_02; Core Feature 4 → task_04; Core Feature 5 → tasks 01, 03 and 04
- Success Metric 1 → task_01
- Success Metric 2 → task_02
- Success Metric 3 → task_03
- Success Metric 4 → task_04
- Success Metric 5 → task_05
- Success Metric 6 → task_03's load rule, proven after the release (Unreachable Acceptance)
- Success Metric 7 → `retiredSkills` of tasks 01 and 04, proven after the release (Unreachable Acceptance)

Deviations from the PRD: None.

## Absence Evidence

1. `skillDispatchGuidePath` and the bundle entry — command: `rg -n
   'skill-dispatch' internal/daemon internal/agent internal/spec`; result: no
   match at `35de7f74`.
2. The tests `TestPersonStartedSkillsAreDispatchedOnlyWhenThePersonAsks`,
   `TestTheFrontendBundlesSeparateProductUIFromPublicPages`,
   `TestTheTypeScriptSnapshotsCarryTheCatalogActivationBundles`,
   `TestTheTypeScriptSetupsCarryTheCompletedUiUxProMaxTree`,
   `TestBuildTaskPromptTellsTheAgentToLoadOnlyTheDispatchedSkills`,
   `TestBuildQAPromptTellsTheGateToLoadItsDispatchedSkills`,
   `TestAssembleTaskContextBundleListsTheSkillDispatchGuideWhenPresent` and
   `TestImplementSpecIsNotARoundfixOwnedSkill` — command: `rg -n
   'PersonStarted|SeparateProductUI|CatalogActivationBundles|UiUxProMax|LoadOnlyTheDispatched|LoadItsDispatched|SkillDispatchGuide|ImplementSpecIsNot'
   -g '*_test.go'`; result: no match.
3. `bundle.frontend-public-pages` and `trigger.frontend.public-pages` —
   command: `rg -n 'public-pages' internal/baseline/assets`; result: no
   match.
4. The glossary term Loaded Skill List and the sections "## Skills for the
   Task", "### Skills the gate loads" and "### Pending Questions and the
   planning files" — command: `rg -n 'Loaded Skill List|Skills for the
   Task|Skills the gate loads|Pending Questions and the planning' CONTEXT.md
   .agents/skills`; result: no match.

## ADR Alternatives

- ADR-0346: this design is the decision; it injects no skill body, loads no
  unmatched skill and adds no Daemon check, the three alternatives the ADR
  rejects.
- ADR-0284: the dispatch guide joins the inputs that refuse on any change,
  like `AGENTS.md`; no binding-text tolerance is extended to it.
- ADR-0170: unchanged; a Task already completed on its target is still
  nothing to carry, whatever its inputs.
- ADR-0246: the eight new Retired Skills follow it exactly; the rejected
  alternatives (keep the names as unrequired snapshot entries, delete the
  adopter's copy) stay rejected.
- ADR-0206: the seven catalog removals would leave the snapshots under this
  rule alone; listing them as Retired Skills adds the adopter report it
  lacks, without editing the upstream list.

## Integration Points

- The local checkout `~/dev/skills`, read only by task_01, task_02 and
  task_04's import commands at `86df860`. No Verification command or QA row
  reads it.
- ACP runtimes receive the prompt text; nothing in the runtime adapter
  changes.

## Testing Approach

- `internal/baseline` (task_01): `TestNoCatalogEntryNamesARetiredSkill` and
  `TestARetiredSkillNameInTheCatalogIsReported` follow the new set, and the
  first also reads each snapshot's activation bundles;
  `TestTheTypeScriptSnapshotsCarryTheCatalogActivationBundles` (new, in
  `retired_skills_test.go`) compares the `typescript` and
  `go-cli-typescript-bun` bundles with `skill-activations.json`;
  `upstreamSetupRevision` and `onioncryTreeDigest` move to the new commit;
  `TestTheTypeScriptSetupsCarryTheCompletedUiUxProMaxTree` (new, in
  `upstream_setup_revision_test.go`) pins the tree of API Contract 3.
- `internal/baseline` (task_02): a new
  `internal/baseline/skill_dispatch_review_test.go` with
  `TestPersonStartedSkillsAreDispatchedOnlyWhenThePersonAsks` (the five
  triggers of Data Model 3 start with "When the person asks", and none of
  those five skills is in an activation bundle) and
  `TestTheFrontendBundlesSeparateProductUIFromPublicPages` (Data Model 2's
  frontend bundles exactly; `wcag-audit-patterns` and
  `web-design-guidelines` in no bundle; no bundle holds both
  `frontend-design` and `interface-design`).
- `internal/agent` (task_03): two new tests in `spec_prompt_test.go` assert
  the Data Model 4 lines inside the invariants and the QA contract, each
  exactly once.
- `internal/daemon` (task_03):
  `TestAssembleTaskContextBundleListsTheSkillDispatchGuideWhenPresent` in
  `task_context_test.go`, a present guide listed right after the
  `implement-task` skill, an absent one not listed.
- `internal/spec` and `internal/cli` (task_03): the expected list of
  `TestCarryForwardInputs…` in `spec_test.go` gains the guide; the table of
  `TestCarryForwardRefusesATaskWhoseBindingTextChanged` in
  `carryforward_amendment_test.go` gains a `skill dispatch guide` case and
  its fixture writes the guide.
- `skills` (task_04): `TestImplementSpecIsNotARoundfixOwnedSkill` (new, in
  `retired_owned_skill_test.go`); `TestCouncilIsNotARoundfixOwnedSkill`
  counts twelve; the `person-only-skills-keep-the-lock` subtest of
  `TestAnAgentMayStartTheSetupSkill` is removed with the only person-only
  owned skill; `internal/baseline/python_setup_test.go` drops
  `implement-spec` from the expected Python setup. No top-level test is
  removed or renamed, so the coverage record does not move.

## Build Order

1. The catalog follows `86df860` and seven skills retire (task_01).
2. The dispatch review (depends on: 1; task_02).
3. Runs load the dispatched skills (depends on: 2; task_03).
4. `implement-spec` retires (depends on: 3; task_04).
5. The QA gate (depends on: 2, 4; task_05).

No other Spec is required. The four Tasks run in series: task_02 edits the
derived files and this repository's dispatch guide task_01 edits; task_03
names the dispatch guide task_02 renders and edits the Roundfix skill, the
owned-skill version record and the Behavior Surface record task_01 edits;
task_04 edits all of them again.

## Risks & Considerations

- The `Makefile` grant covers only the `OWNED_SKILLS` line and the count in
  its comment; task_04 changes no other `Makefile` line.
- A Task carried across this Spec's merge sees its `implement-task` input
  moved and runs again; parked items of other Specs may need a fresh Run.
- The load rule costs context in every Task; the rule's "only what matches"
  and the reworded triggers keep it bounded. Over-matching is visible in the
  Loaded Skill List and in the session records.
- An adopter whose dispatch guide is stale still gets the rule; it loads the
  skills its own guide names.
- The catalog commit may be force-moved upstream; task_01 stops when the
  local checkout lacks it.

## Glossary

- adds: **Loaded Skill List**
- changes: **Spec Context Bundle**

## Decisions

- QA follows the load rule through its contract and `qa-gate`, without a
  list. See ADR-0346.
- No Daemon check of the list. See ADR-0346.
- The dispatch guide is a carry-forward input. See ADR-0346.
- The TypeScript snapshot's bundles are set from `skill-activations.json`
  before each sync, because the sync keeps recorded bundles and no other
  command writes them.
- This repository restores `context7-cli`, `grilling` and `tui-design`, the
  upstream skills its profile holds that changed at `86df860`.
- The owned-version key of `implement-spec` is deleted with `jq`, the only
  hand edit of the record.
