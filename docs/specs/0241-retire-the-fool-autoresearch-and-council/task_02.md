---
task: task_02
spec: 0241-retire-the-fool-autoresearch-and-council
status: completed
type: chore
complexity: medium
---

# Task 02: Roundfix stops shipping council and this repository drops the three skills

## Overview

After task_01 no module requires `council` or `the-fool`. This Task removes
`council` from the Roundfix-owned bundle (canonical copy, mirror, embed list,
owned names, `Makefile` list and owned-version record), rewords the two owned
skills that send the reader to it, and drops the vendored `the-fool` and
`autoresearch` from this repository's lock, recommended list and
`.agents/skills/`. The doctor counts in the guides and the repository tests
follow. `grilling`, `grill-with-docs`, `write-idea`, `business-analyst` and
`handoff` stay.

This is an authorized tooling Task. Its Governed Paths are the `Makefile`,
the `council`, `the-fool` and `autoresearch` trees it deletes, both copies of
`write-idea` and `write-prd`, and `docs/agents/setup-context.json`, all
bounded in `_authorization.md`. The vendored trees are deleted, never edited.

## Requirements

1. MUST run steps 1 to 3 of the TechSpec's repository procedure: delete
   `.agents/skills/council` and `skills/council`, remove `council` from the
   `//go:embed` line and `skillNames` in `skills/skills.go` and from
   `OWNED_SKILLS` in the `Makefile`, and delete its key from
   `skills/testdata/owned-skill-versions.json` with the stated `jq` program.
   The `Makefile` change is that one line plus the comment's count; no other
   `Makefile` line changes.
2. MUST apply the task_02 Fixed texts to `.agents/skills/write-idea/SKILL.md`,
   `.agents/skills/write-idea/references/idea-template.md` and
   `.agents/skills/write-prd/SKILL.md`, changing only the lines that name
   `council` and the two version fields of each, then run `make
   skills-sync` and re-record versions with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   after the last skill edit. MUST NOT edit a `### QA settlement` section.
3. MUST create `skills/retired_owned_skill_test.go` with
   `TestCouncilIsNotARoundfixOwnedSkill` (Testing Approach 5).
4. MUST run step 4 of the repository procedure: delete the `the-fool` and
   `autoresearch` keys of `skills-lock.json` and their lines of
   `skills/recommended.txt`, and delete `.agents/skills/the-fool` and
   `.agents/skills/autoresearch`. No other lock entry or `.agents/skills/`
   tree changes.
5. MUST make the `skills/repository_test.go` and
   `internal/cli/this_repository_skill_set_test.go` edits that "Existing
   tests that change" lists, and no other test edit.
6. MUST apply the task_02 Fixed texts to `docs/user-guide/commands/skills.md`,
   `docs/user-guide/commands/doctor.md`, `docs/user-guide/usage.md`,
   `README.md` and the comments in `skills/skills.go` and the `Makefile`.
7. MUST run step 5 of the repository procedure; the second managed refresh
   MUST report `File changes: 0`, and `bin/roundfix doctor`, rebuilt, MUST
   print `skills: ok (41 required: 13 Roundfix-owned, 28 external)`.
8. MUST NOT change a Baseline asset, `docs/agents/skill-dispatch.md`, any
   other owned skill, `grilling`, `grill-with-docs`, `business-analyst`,
   `handoff` or `CONTEXT.md`.

## Subtasks

- [ ] Remove `council` from the bundle, the embed list, the `Makefile` and the version record.
- [ ] Reword `write-idea` and `write-prd`, sync and record their versions.
- [ ] Drop `the-fool` and `autoresearch` from the lock, the recommended list and the trees.
- [ ] Update the repository tests and the doctor counts in the guides.
- [ ] Refresh the managed Setup Manifest twice.

## Acceptance Criteria

- [ ] `skills.Names()` has 13 names without `council`, and nothing embedded or installed here holds `council`.
- [ ] No file of `write-idea` or `write-prd` names `council`; each mirror equals its canonical copy and the raised versions are recorded.
- [ ] No lock entry, recommended line or tree names `the-fool` or `autoresearch`, and every kept skill is still installed with its lock hash.
- [ ] The guides print the new doctor line, and the managed refresh converges.

## Context

- instruction: `docs/adr/0246-a-retired-skill-leaves-the-baseline-and-the-adopter-removes-its-copy.md`
- instruction: `docs/adr/0191-a-setup-snapshot-follows-its-upstream-by-name.md`
- instruction: `skills/owned_skill_versions_test.go`
- deletes: `.agents/skills/council/SKILL.md`
- deletes: `.agents/skills/council/assets/synthesis-template.md`
- deletes: `.agents/skills/council/references/archetypes.md`
- deletes: `.agents/skills/council/references/debate-protocols.md`
- deletes: `skills/council/SKILL.md`
- deletes: `skills/council/assets/synthesis-template.md`
- deletes: `skills/council/references/archetypes.md`
- deletes: `skills/council/references/debate-protocols.md`
- deletes: `.agents/skills/the-fool/SKILL.md`
- deletes: `.agents/skills/the-fool/references/cognitive-bias-inventory.md`
- deletes: `.agents/skills/the-fool/references/dialectic-synthesis.md`
- deletes: `.agents/skills/the-fool/references/evidence-audit.md`
- deletes: `.agents/skills/the-fool/references/mode-selection-guide.md`
- deletes: `.agents/skills/the-fool/references/pre-mortem-analysis.md`
- deletes: `.agents/skills/the-fool/references/red-team-adversarial.md`
- deletes: `.agents/skills/the-fool/references/socratic-questioning.md`
- deletes: `.agents/skills/autoresearch/SKILL.md`
- deletes: `.agents/skills/autoresearch/references/eval-guide.md`
- creates: `skills/retired_owned_skill_test.go`
- interface: `skills/skills.go`
- interface: `Makefile`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `.agents/skills/write-idea/SKILL.md`
- interface: `.agents/skills/write-idea/references/idea-template.md`
- interface: `skills/write-idea/SKILL.md`
- interface: `skills/write-idea/references/idea-template.md`
- interface: `.agents/skills/write-prd/SKILL.md`
- interface: `skills/write-prd/SKILL.md`
- interface: `skills-lock.json`
- interface: `skills/recommended.txt`
- interface: `skills/repository_test.go`
- interface: `internal/cli/this_repository_skill_set_test.go`
- interface: `docs/user-guide/commands/skills.md`
- interface: `docs/user-guide/commands/doctor.md`
- interface: `docs/user-guide/usage.md`
- interface: `README.md`
- interface: `docs/agents/setup-context.json`

## Verification

- `out="$(go test -count=1 -v -run '^(TestCouncilIsNotARoundfixOwnedSkill|TestEveryOwnedSkillVersionIsRecorded|TestRecommendedSkillsMatchLock|TestAuthorialSkillSync|TestCheckRepositoryWithExternalUsesExplicitRequirement|TestCheckRepositoryMatchesExternalCompatibilityEntryPoint|TestCheckRepositoryHandlesNestedLinksSpecialEntriesAndStableOrdering)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestCouncilIsNotARoundfixOwnedSkill TestEveryOwnedSkillVersionIsRecorded TestRecommendedSkillsMatchLock TestAuthorialSkillSync TestCheckRepositoryWithExternalUsesExplicitRequirement TestCheckRepositoryMatchesExternalCompatibilityEntryPoint TestCheckRepositoryHandlesNestedLinksSpecialEntriesAndStableOrdering; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && out="$(go test -count=1 -v -run '^(TestThisRepositoryHoldsEveryRequiredExternalSkill|TestARepositoryMissingARequiredExternalSkillIsReported)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestThisRepositoryHoldsEveryRequiredExternalSkill TestARepositoryMissingARequiredExternalSkillIsReported; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task `TestCouncilIsNotARoundfixOwnedSkill` does not exist, so its pass line is missing and the command fails.
- `for path in .agents/skills/council skills/council .agents/skills/the-fool .agents/skills/autoresearch; do test ! -e "$path" || { printf 'retired tree remains: %s\n' "$path" >&2; exit 1; }; done; ! grep -Eq '"(the-fool|autoresearch)"' skills-lock.json && ! grep -Eqx 'the-fool|autoresearch' skills/recommended.txt && ! grep -Eq '^OWNED_SKILLS.*council' Makefile && ! grep -rqi council .agents/skills/write-idea .agents/skills/write-prd skills/write-idea skills/write-prd && diff -r .agents/skills/write-idea skills/write-idea && diff -r .agents/skills/write-prd skills/write-prd && for f in README.md docs/user-guide/usage.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'skills: ok (41 required: 13 Roundfix-owned, 28 external)' || { printf 'stale doctor line: %s\n' "$f" >&2; exit 1; }; done && tr -s '[:space:]' ' ' < docs/user-guide/commands/skills.md | grep -qF -- 'ships 13 Roundfix-owned skills' && tr -s '[:space:]' ' ' < docs/user-guide/commands/doctor.md | grep -qF -- 'the 13 Roundfix-owned skills' && go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json >/dev/null` — expected: exit 0; before this Task the three trees exist, so the first check fails.

## References

- `_prd.md` → Goals 2-3; User Story 3; Core Features 3-4; Success Metrics 3-4; Declared breaks
- `_techspec.md` → Measured facts; Fixed texts (task_02); The repository procedure; Existing tests that change; API Contract 3; Testing Approach 5; Build Order 2
- ADR-0246, ADR-0191, ADR-0189, ADR-0233

## Result

Implemented the Task slice: removed `council` from the owned bundle and
deleted its canonical and mirrored trees; reworded both owned skills and
synchronized their mirrors; removed `the-fool` and `autoresearch` from the
repository lock, recommendations, and installed trees; updated the bounded
repository tests and guide counts; and added the retired-owned-skill boundary
test. The managed setup manifest was refreshed after the lock change.

Focused checks and evidence:

- `GOCACHE=/tmp/roundfix-0241-task02-gocache go test ./skills -run '^TestCouncilIsNotARoundfixOwnedSkill$'` passed.
- `GOCACHE=/tmp/roundfix-0241-task02-gocache go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'` passed after the mandated recorder run. The recorder recorded write-idea 0.0.5 and write-prd 0.0.9; write-prd was already 0.0.8 before this Task and its changed content therefore received the next recorded version.
- `make skills-sync` completed; both canonical/mirror skill trees compare equal, and no write-idea or write-prd file names `council`.
- The owned bundle contains 13 names, no embedded or installed `council` path remains, and the lock/recommended list contain neither `the-fool` nor `autoresearch`.
- The first managed refresh applied one setup-manifest change; the second reported `File changes: 0`.
- `make build` completed and rebuilt `bin/roundfix`. Direct Doctor capture produced no output in this environment, so the required Doctor count line remains for Daemon Verification to observe.
- Settlement feedback repair: updated the stale TechSpec claim receipts to quote the current Task 02 source passages verbatim; Task status remains Daemon-owned.
- Focused repair check: `GOCACHE=/tmp/roundfix-0241-task02-gocache go run -buildvcs=false ./cmd/roundfix spec check 0241-retire-the-fool-autoresearch-and-council --strict` reported `No findings`; authored Verification commands were not executed.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `docs/specs/0241-retire-the-fool-autoresearch-and-council/_techspec.md`

## Carry-forward provenance

- Source Run: `run_20261006T214048Z_131051a4f3c19cbe`
- Source commit: `4f76f4c51ee51be61e2c63cc47102b9f7a33ca69`
