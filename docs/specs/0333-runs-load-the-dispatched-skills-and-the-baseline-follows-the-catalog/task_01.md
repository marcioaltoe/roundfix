---
task: task_01
spec: 0333-runs-load-the-dispatched-skills-and-the-baseline-follows-the-catalog
status: completed
type: backend
complexity: high
---

# Task 01: The Baseline follows catalog 86df860 and retires the seven skills it removed

## Overview

Every upstream Setup Snapshot moves to catalog commit `86df860`, which
removed seven skills and completed `ui-ux-pro-max`. The seven become Retired
Skills, so no module, trigger, bundle or snapshot names them and
`baseline update` lists an adopter's copies to delete. The TypeScript
snapshot's activation bundles follow the catalog file, this repository
restores its three upstream skills the commit changed, and the Baseline guide
lists the new Retired Skills and the Python 3 need of `ui-ux-pro-max`. It
answers the two adopted Backlog Entries of 2026-10-30
(`references/2026-10-30-the-baseline-follows-the-curated-skills-catalog.md`
and `references/2026-10-30-the-completed-ui-ux-pro-max-updates-automatically.md`).

## Requirements

1. MUST stop before any write when `~/dev/skills` lacks commit
   `86df8608726979064a8fe46003dfaf403ffe0087` (`_techspec.md` → API Contract
   1, first three lines), and report it.
2. MUST add `design-patterns`, `diagnosing-bugs`, `fixing-metadata`,
   `game-changing-features`, `refactoring-analysis`, `tdd` and `to-prompt` to
   `retiredSkills` in `internal/baseline/retired_skills.go` (Data Model 1).
3. MUST remove, as text that keeps each file's layout, the required-skill
   line and the whole dispatch entry of `diagnosing-bugs`,
   `game-changing-features`, `refactoring-analysis`, `tdd` and `to-prompt`
   from `internal/baseline/assets/modules/typescript.json`, of
   `design-patterns` from `modules/backend.json` and of `fixing-metadata`
   from `modules/frontend.json`; remove `diagnosing-bugs` from
   `bundle.debugging` and `tdd` from `bundle.testing` in
   `internal/baseline/assets/skill-activations.json` and raise its `version`
   from 2 to 3. No other module, trigger or bundle changes in this Task.
4. MUST run `go test ./internal/baseline -run
   '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions
   -count=1`, then the rest of `_techspec.md` → API Contract 1 in order: the
   `typescript.json` bundle assignment, the asset sync and its `--check`
   (exit `0`), the parity transform of
   `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`,
   `make baseline-digests`, the two managed refreshes (the second reports
   `File changes: 0`) and the confirmed restore of `context7-cli`,
   `grilling` and `tui-design`. No snapshot, golden, digest, guide under
   `docs/agents/`, upstream skill file or lock entry is hand-edited.
5. MUST change `internal/baseline/retired_skills_test.go`: the expected list
   of `TestNoCatalogEntryNamesARetiredSkill` is the nine names of Data Model
   1, and that test also passes each snapshot's `activationBundles` to
   `namedSkillFindings` as bundles; add
   `TestTheTypeScriptSnapshotsCarryTheCatalogActivationBundles`, which
   requires the `typescript` and `go-cli-typescript-bun` snapshots to carry
   every bundle of `skill-activations.json` with the same skills in the same
   order.
6. MUST change `internal/baseline/upstream_setup_revision_test.go`:
   `upstreamSetupRevision` becomes
   `86df8608726979064a8fe46003dfaf403ffe0087` and `onioncryTreeDigest` the
   tree of `_techspec.md` → API Contract 3; add
   `TestTheTypeScriptSetupsCarryTheCompletedUiUxProMaxTree`, which requires
   the `typescript` and `go-cli-typescript-bun` snapshots' `ui-ux-pro-max`
   to name that revision and the tree of API Contract 3, through a new
   constant.
7. MUST update `docs/user-guide/commands/baseline.md`, section "### Retired
   skills": after `council` and `the-fool` since 0.52.0, name
   `design-patterns`, `diagnosing-bugs`, `fixing-metadata`,
   `game-changing-features`, `refactoring-analysis`, `tdd` and `to-prompt`
   since 0.127.0, which the catalog removed, and give the removal commands
   for them in the existing shape. Add, before the paragraph on
   `baseline skills reconcile`, a section "### Skills that run scripts":
   the TypeScript setup's `ui-ux-pro-max` runs its searches through Python
   scripts in its own folder, so they need Python 3 (the standard library is
   enough); without Python the skill still applies its `SKILL.md` rules and
   its quick reference; no Baseline command checks for Python. The page MUST
   contain the phrases `which the catalog removed` and `need Python 3`.
8. MUST add to the canonical `.agents/skills/roundfix/references/baseline.md`,
   after "it never deletes the copy and never changes the update's state.",
   one paragraph naming the Retired Skills of requirement 7 the same way and
   saying that `ui-ux-pro-max` searches need Python 3, containing `which the
   catalog removed` and `need Python 3`; copy the canonical skills to
   `skills/` with `make skills-sync`, then raise both version fields of the
   `roundfix` skill and its mirror with `go test ./skills -run
   '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
   -count=1`, recording `skills/testdata/owned-skill-versions.json`.
9. MUST re-record the Behavior Surface fingerprint as the last step, after
   the last guide edit: `go test -tags docscontract ./internal/docscontract
   -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
   -count=1`, recording `docs/references/behavior-surfaces.json`.

## Subtasks

- [ ] Check the catalog commit, then edit the retired set, the three modules and the activation file.
- [ ] Record the module versions, sync the snapshots, transform the parity fixture and regenerate the derived files.
- [ ] Refresh this repository's managed files and restore its three changed upstream skills.
- [ ] Update the catalog tests, the Baseline guide and the Roundfix baseline reference; record the skill versions and re-record the Behavior Surface fingerprint last.

## Acceptance Criteria

- [ ] Every upstream snapshot names `86df860`, carries the completed `ui-ux-pro-max`, and no snapshot, module, trigger or bundle names a Retired Skill.
- [ ] The TypeScript snapshots' bundles equal the catalog file's.
- [ ] This repository's Baseline preview reports no differing skill, and a second managed refresh changes no file.
- [ ] The Baseline guide names the new Retired Skills and the Python 3 need.

## Context

- instruction: `docs/adr/0246-a-retired-skill-leaves-the-baseline-and-the-adopter-removes-its-copy.md`
- instruction: `docs/adr/0206-the-baseline-takes-upstream-setup-names-and-drops-a-removed-skill.md`
- interface: `internal/baseline/retired_skills.go`
- interface: `internal/baseline/retired_skills_test.go`
- interface: `internal/baseline/upstream_setup_revision_test.go`
- interface: `internal/baseline/assets/modules/typescript.json`
- interface: `internal/baseline/assets/modules/backend.json`
- interface: `internal/baseline/assets/modules/frontend.json`
- interface: `internal/baseline/assets/skill-activations.json`
- interface: `internal/baseline/module-versions.json`
- interface: `internal/baseline/assets/setups/go.json`
- interface: `internal/baseline/assets/setups/python.json`
- interface: `internal/baseline/assets/setups/rust.json`
- interface: `internal/baseline/assets/setups/typescript.json`
- interface: `internal/baseline/assets/setups/go-cli-typescript-bun.json`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`
- interface: `internal/baseline/testdata/parity-corpus/v1/manifest.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/setup-context.json`
- interface: `docs/agents/skill-dispatch.md`
- interface: `.agents/skills/context7-cli/SKILL.md`
- interface: `.agents/skills/context7-cli/references/docs.md`
- interface: `.agents/skills/grilling/SKILL.md`
- interface: `.agents/skills/tui-design/SKILL.md`
- creates: `.agents/skills/tui-design/references/evidence.md`
- interface: `skills-lock.json`
- interface: `docs/user-guide/commands/baseline.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The Setup Snapshots and the Retired Skills — consumers: `./internal/baseline`

## Verification

- `out="$(go test -count=1 -timeout 20m -v ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestNoCatalogEntryNamesARetiredSkill TestARetiredSkillNameInTheCatalogIsReported TestTheTypeScriptSnapshotsCarryTheCatalogActivationBundles TestTheTypeScriptSetupsCarryTheCompletedUiUxProMaxTree TestTheTypeScriptSetupsCarryTheRepublishedOnioncryTree TestEveryUpstreamSetupFollowsTheAiMemoryRevision TestEveryUpstreamSetupFollowsThePythonPresetRevision; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/user-guide/commands/baseline.md | grep -qF -- 'which the catalog removed' || { echo 'baseline guide lacks the removed skills' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/baseline.md | grep -qF -- 'need Python 3' || { echo 'baseline guide lacks the Python note' >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/baseline.md | grep -qF -- 'which the catalog removed' || { echo 'baseline reference lacks the removed skills' >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/baseline.md | grep -qF -- 'need Python 3' || { echo 'baseline reference lacks the Python note' >&2; exit 1; }; cmp .agents/skills/roundfix/references/baseline.md skills/roundfix/references/baseline.md && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && go test -count=1 -timeout 20m ./skills || exit 1; h="$(mktemp -d)"; trap 'chmod -R u+w "$h"; rm -rf "$h"' EXIT; go build -buildvcs=false -o "$h/roundfix" ./cmd/roundfix; preview="$(HOME="$h" "$h/roundfix" baseline update --repo . --no-history --format text 2>&1)"; printf '%s\n' "$preview" | grep -qF 'Skills differing: 0' || { printf '%s\n' "$preview" >&2; exit 1; }; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task `TestTheTypeScriptSnapshotsCarryTheCatalogActivationBundles` does not exist, so it is missing from the passes. Afterwards the whole Baseline package passes with the new snapshots, the revision and tree tests name `86df860`, the guide names the removed skills and the Python need, this repository's preview reports no differing skill, and the documentation contracts, the Behavior Surface fingerprint among them, pass.

## References

- `_prd.md` → Goals; Core Features 2 and 5; Story 3; Success Metric 1; Sources
- `_techspec.md` → Data Models 1 and 2; API Contracts 1 and 3; Testing Approach; Build Order 1
- ADR-0246; ADR-0206


## Result

Implemented the catalog-following slice of task_01. The seven catalog removals
join the two existing Retired Skills; only their required-skill lines and
whole dispatch entries were removed from the three modules. The debugging and
testing bundles drop `diagnosing-bugs` and `tdd`, and the activation version
is 3. The module recorder raised backend and frontend to 7 and TypeScript to
9. Person-started trigger changes and `implement-spec` retirement remain
outside this Task.

All five Setup Snapshots were synchronized from the local catalog at
`86df8608726979064a8fe46003dfaf403ffe0087`. The TypeScript snapshots carry
every catalog activation bundle with identical skill order. The new tests
cover that contract, the completed `ui-ux-pro-max` tree, and retired names in
snapshot activation bundles. The existing onioncry assertion now follows the
republished tree.

The sanctioned restore refreshed `context7-cli`, `grilling` and `tui-design`,
including the new TUI evidence reference, and updated only their lock entries.
The user guide gives removal commands for the seven skills since 0.127.0 and
the Python 3 requirement, standard-library sufficiency, fallback guidance and
absence of a Baseline Python check. The canonical Roundfix reference carries
the same facts; its mirror was synchronized and both skill version fields
were recorded at 0.1.159.

### Acceptance evidence

| Acceptance criterion | Focused evidence |
| --- | --- |
| Current upstream snapshots, completed `ui-ux-pro-max`, no retired catalog references | The catalog commit prerequisite exited 0 before any write. Asset sync reported five `skills.setup-snapshot.updated` findings; its `--check` exited 0. Focused retired-name, snapshot-revision, onioncry and `ui-ux-pro-max` tree tests passed. A read-only provenance audit confirmed every GitHub snapshot and skill ref names the new commit. |
| TypeScript snapshot bundles equal the catalog | `TestTheTypeScriptSnapshotsCarryTheCatalogActivationBundles` passed for both snapshots; the audit compared bundles by ID and retained skill order. Composed bundle ordering is preserved. |
| Repository preview has no differing skills; second refresh is idempotent | The second isolated-HOME managed refresh exited 0 with `File changes: 0` and `Idempotence: verified`. After the confirmed restore and final skill edit, a rebuilt binary's isolated-HOME offline preview exited 0 with `Baseline update: current`, `File changes: 0` and `Skills differing: 0`. |
| Both guides name the removed skills and Python need | Read-only assertions found all seven names and both required phrases in each guide. Canonical and mirrored Roundfix guide and SKILL.md bytes match. Owned-skill record and mirror checks passed; the final behavior fingerprint recorder exited 0 after the last guide edit. |

### Commands and outcomes

- Prerequisite: `git -C "$HOME/dev/skills" cat-file -e
  86df8608726979064a8fe46003dfaf403ffe0087^{commit}` — exit 0 before edits.
- Initial regression check: `go test ./internal/baseline -run
  '^(TestNoCatalogEntryNamesARetiredSkill|TestTheTypeScriptSnapshotsCarryTheCatalogActivationBundles|TestTheTypeScriptSetupsCarryTheCompletedUiUxProMaxTree)$'
  -count=1` — exit 1 as expected: the retired set still contained two names
  and both TypeScript snapshots carried the old UI skill revision and tree.
- `go test ./internal/baseline -run
  '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions
  -count=1` — exit 0.
- API Contract 1's offline fetch, TypeScript bundle assignment, asset sync,
  asset `--check`, and parity jq transform ran in order. Both sync commands
  exited 0. `make baseline-digests` then exited 0 and regenerated the declared
  formatter/profile/catalog/parity/plan artifacts.
- After rebuilding the binary, `baseline update --repo . --no-skills
  --no-history --yes --format text` ran twice with fresh HOME directories —
  both exited 0; the first changed only `docs/agents/setup-context.json`, and
  the second changed no file. Existing nested-carrier warnings named the
  formatter fixture and retained Source Baseline fixture; both remained
  untouched.
- `baseline skills restore --profile go-cli-tui --skill context7-cli
  --skill grilling --skill tui-design --source-dir <fetched> --repo .
  --format text` — preview exited 3 with exactly the declared skill and lock
  paths. The same command with `--confirm-plan
  8d0860c77c04fe4f42f4ba5c4a26f5d08543abf16a2f0efd803c22de3b2892fb`
  exited 0 and reported `restore.completed`.
- `make skills-sync`, then `go test ./skills -run
  '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1`
  — both exited 0. A subsequent `make baseline-digests` exited 0 with
  `changed: false` after the owned guide/version edit.
- Final catalog-focused command: `go test ./internal/baseline -run
  '^(TestNoCatalogEntryNamesARetiredSkill|TestARetiredSkillNameInTheCatalogIsReported|TestTheTypeScriptSnapshotsCarryTheCatalogActivationBundles|TestTheTypeScriptSetupsCarryTheCompletedUiUxProMaxTree|TestTheTypeScriptSetupsCarryTheRepublishedOnioncryTree|TestEveryUpstreamSetupFollowsTheAiMemoryRevision|TestEveryUpstreamSetupFollowsThePythonPresetRevision|TestEverySetupCarriesTheKnowledgeWorkspaceThatRetiredTheSecondbrainInbox)$'
  -count=1` — exit 0.
- `go test ./skills -run
  '^(TestEveryOwnedSkillVersionIsRecorded|TestAuthorialSkillSync)$' -count=1`
  — exit 0.
- Rebuilt-binary preview: `baseline update --repo . --no-history
  --skills-source-dir <fetched> --format text`, with a fresh HOME — exit 0,
  current Baseline, no file changes or differing skills.
- Final regeneration: `go test -tags docscontract ./internal/docscontract
  -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`
  — exit 0; recorded `docs/references/behavior-surfaces.json` after the last
  guide edit.
- `git diff --check` and a read-only changed-path audit found no whitespace
  errors or paths outside this Task's declared interfaces. No other Task,
  Task Graph, or Task status was edited by the Agent.

Skills loaded: implement-task, coding-guidelines, golang-testing, testing-boss, tech-writer, evidence-gate, setup-context-driven, knowledge-workspace, roundfix.

The declared Verification, repository Verification and terminal settlement
were not run; they remain Daemon-owned. No commit, push or pull request was
created. No follow-up beyond the already-authored later Tasks was found.
