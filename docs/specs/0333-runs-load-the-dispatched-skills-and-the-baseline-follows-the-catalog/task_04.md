---
task: task_04
spec: 0333-runs-load-the-dispatched-skills-and-the-baseline-follows-the-catalog
status: pending
type: backend
complexity: high
---

# Task 04: implement-spec retires and its two rules move into the Roundfix skill

## Overview

The maintainer decided on 2026-10-30 to remove the person-started
`implement-spec` skill from the Baseline. It becomes a Retired Skill under
ADR-0246: it leaves the binary's bundle, its canonical copy, its mirror, the
owned-version record, the `Makefile`'s owned list, the asset sync's owned set
and the `context-workflow` module; `archive-spec` and `write-tasks` stop
naming it; every Setup Snapshot drops it although the catalog presets still
list it; and `baseline update` lists an adopter's copy to delete. Its two
rules no other skill states, asking the maintainer only the queue's Pending
Question and never editing a delivered Spec's planning files, move into the
Roundfix skill's delivery reference. The guides follow.

## Requirements

1. MUST change only the `Makefile`'s `OWNED_SKILLS` line and the owned-skill
   count in the comment directly above it, as the authorization record's
   Makefile grant allows; every other `Makefile` line stays unchanged.
2. MUST add `implement-spec` to `retiredSkills`
   (`internal/baseline/retired_skills.go`) and to the expected list of
   `TestNoCatalogEntryNamesARetiredSkill`; remove it from
   `assetsSyncRepoOwnedSkills` (`internal/baseline/assets_sync.go`), from
   the `//go:embed` line, `skillNames` and the comment's count (twelve to
   eleven authorial skills) in `skills/skills.go`, and from the `Makefile`'s
   `OWNED_SKILLS` line and its comment's count, with no other `Makefile`
   change (`_techspec.md` → API Contract 8).
3. MUST delete `.agents/skills/implement-spec/SKILL.md` and
   `skills/implement-spec/SKILL.md`, and delete the `implement-spec` key of
   `skills/testdata/owned-skill-versions.json` with `jq --indent 2
   'del(.skills["implement-spec"])'` through a temporary copy, the only hand
   edit of the record.
4. MUST remove, as text that keeps the file's layout, the required-skill
   line and the whole dispatch entry of `implement-spec` (both triggers) from
   `internal/baseline/assets/modules/context-workflow.json`; run the module
   record command with its outputs, then `_techspec.md` → API Contract 2 in
   order: the `typescript.json` bundle assignment, the asset sync and its
   `--check`, the parity transform, `make baseline-digests` and the two
   managed refreshes (the second reports `File changes: 0`).
5. MUST change the tests: `internal/baseline/python_setup_test.go` drops
   `implement-spec` from the expected Python setup;
   `TestCouncilIsNotARoundfixOwnedSkill` (`skills/retired_owned_skill_test.go`)
   counts twelve owned skills; add `TestImplementSpecIsNotARoundfixOwnedSkill`
   in the same file, asserting the name is absent from `Names()`, the embedded
   bundle and `.agents/skills`; remove the
   `person-only-skills-keep-the-lock` subtest of
   `TestAnAgentMayStartTheSetupSkill` (`skills/setup_skill_invocation_test.go`),
   whose only subject was `implement-spec`. No top-level test is removed or
   renamed.
6. MUST reword the canonical `.agents/skills/archive-spec/SKILL.md` at its
   description and at "The trigger is spec completion" to "a Roundfix Spec
   Run" in place of the `implement-spec` loop, and the canonical
   `.agents/skills/write-tasks/SKILL.md` sentence to "the manifest that
   `implement-task` and the Roundfix Daemon schedule from"; add the section
   "### Pending Questions and the planning files" of `_techspec.md` → API
   Contract 9 to the canonical
   `.agents/skills/roundfix/references/spec-delivery.md` before "###
   Recovering without breaking the loop", containing `Pending Question` and
   `_prd.md`, `_techspec.md` or `_tasks.md`. No `### QA settlement` text
   changes.
7. MUST update the guides: `docs/user-guide/commands/skills.md` says the
   binary ships 12 Roundfix-owned skills and no longer names
   `implement-spec`; `docs/user-guide/context-driven-development.md` names
   `implement-task` alone where it named `implement-spec`/`implement-task`;
   `docs/user-guide/commands/baseline.md` section "### Retired skills" adds
   the phrase "the Roundfix-owned `implement-spec` since 0.127.0" to the
   list and its removal commands, and the canonical
   `.agents/skills/roundfix/references/baseline.md` adds `implement-spec` to
   the Retired Skills paragraph task_01 wrote; the canonical
   `.agents/skills/roundfix/references/setup.md`, after the paragraph on
   `roundfix skills check`, says the binary ships twelve Roundfix-owned
   skills, containing `twelve Roundfix-owned skills`. No live guide or skill
   outside `docs/history/`, `CHANGELOG.md` and the Spec folders names
   `implement-spec` except those two Retired Skills lists and their removal
   commands.
8. MUST copy the canonical skills to `skills/` with `make skills-sync`, then
   raise both version fields of `archive-spec`, `write-tasks` and `roundfix`
   and their mirrors with `go test ./skills -run
   '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
   -count=1`.
9. MUST re-record the Behavior Surface fingerprint as the last step, after
   the last guide edit: `go test -tags docscontract ./internal/docscontract
   -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
   -count=1`, recording `docs/references/behavior-surfaces.json`.

## Subtasks

- [ ] Confirm the `Makefile` grant; retire the skill from the code, the bundle, the record and the `Makefile`.
- [ ] Remove it from the module; record the versions, re-sync the snapshots and regenerate the derived files and this repository's guides.
- [ ] Update the tests, the owned skills, the Roundfix delivery reference and the guides; sync and record the skills.
- [ ] Re-record the Behavior Surface fingerprint last.

## Acceptance Criteria

- [ ] The binary ships twelve owned skills, and `implement-spec` is a Retired Skill no module, trigger, bundle or snapshot names.
- [ ] `make skills-sync-check` passes without the skill.
- [ ] The Roundfix delivery reference states the Pending Question rule and the planning-file rule.
- [ ] No live guide or skill points at `implement-spec` outside the Retired Skills list.

## Context

- instruction: `docs/adr/0246-a-retired-skill-leaves-the-baseline-and-the-adopter-removes-its-copy.md`
- interface: `Makefile`
- interface: `internal/baseline/retired_skills.go`
- interface: `internal/baseline/retired_skills_test.go`
- interface: `internal/baseline/assets_sync.go`
- interface: `internal/baseline/python_setup_test.go`
- interface: `internal/baseline/assets/modules/context-workflow.json`
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
- interface: `skills/skills.go`
- interface: `skills/retired_owned_skill_test.go`
- interface: `skills/setup_skill_invocation_test.go`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `.agents/skills/implement-spec/SKILL.md`
- interface: `skills/implement-spec/SKILL.md`
- interface: `.agents/skills/archive-spec/SKILL.md`
- interface: `skills/archive-spec/SKILL.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `.agents/skills/roundfix/references/setup.md`
- interface: `skills/roundfix/references/setup.md`
- interface: `.agents/skills/roundfix/references/spec-delivery.md`
- interface: `skills/roundfix/references/spec-delivery.md`
- interface: `docs/user-guide/commands/skills.md`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `docs/user-guide/commands/baseline.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The owned skill bundle and the Retired Skills — consumers: `./skills`, `./internal/baseline`

## Verification

- `out="$(go test -count=1 -timeout 20m -v ./internal/baseline ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestImplementSpecIsNotARoundfixOwnedSkill TestCouncilIsNotARoundfixOwnedSkill TestNoCatalogEntryNamesARetiredSkill TestThePythonSetupIsThePublishedUpstreamPreset TestAnAgentMayStartTheSetupSkill TestEveryOwnedSkillVersionIsRecorded; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; test ! -e .agents/skills/implement-spec || { echo 'implement-spec tree remains' >&2; exit 1; }; if grep -rlF --include='*.md' --include='*.json' --exclude=baseline.md -- 'implement-spec' docs/agents docs/user-guide/commands/skills.md docs/user-guide/context-driven-development.md .agents/skills skills internal/baseline/assets/modules || grep -qF -- 'implement-spec' Makefile; then echo 'a live file still names implement-spec' >&2; exit 1; fi; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/spec-delivery.md | grep -qF -- '### Pending Questions and the planning files' || { echo 'delivery reference lacks the rules' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/baseline.md | grep -q -- 'Roundfix-owned .implement-spec. since 0[.]127[.]0' || { echo 'baseline guide lacks the retired implement-spec' >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/setup.md | grep -qF -- 'twelve Roundfix-owned skills' || { echo 'setup reference lacks the owned count' >&2; exit 1; }; grep -qF -- 'implement-spec' .agents/skills/roundfix/references/baseline.md || { echo 'baseline reference lacks the retired implement-spec' >&2; exit 1; }; cmp .agents/skills/roundfix/references/setup.md skills/roundfix/references/setup.md && cmp .agents/skills/roundfix/references/baseline.md skills/roundfix/references/baseline.md && cmp .agents/skills/archive-spec/SKILL.md skills/archive-spec/SKILL.md && cmp .agents/skills/write-tasks/SKILL.md skills/write-tasks/SKILL.md && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && cmp .agents/skills/roundfix/references/spec-delivery.md skills/roundfix/references/spec-delivery.md && make skills-sync-check && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task `TestImplementSpecIsNotARoundfixOwnedSkill` does not exist, so it is missing from the passes. Afterwards the Baseline and skills packages pass with twelve owned skills, the tree is gone, no live guide, skill, module or the `Makefile` names the skill except the Retired Skills list of the Baseline guide, the delivery reference carries the two rules, the mirrors equal their canonical copies, `make skills-sync-check` passes, and the documentation contracts pass.

## References

- `_prd.md` → Goals; Core Features 4 and 5; Stories 3, 4 and 5; Success Metric 4; Success Metric 7; Decisions; Open Questions
- `_techspec.md` → Data Model 1; API Contracts 2, 8 and 9; Testing Approach; Build Order 4
- ADR-0246
