---
task: task_03
spec: 0403-the-update-check-names-each-profile-that-differs-from-a-refreshed-recommendation
status: pending
type: docs
complexity: medium
---

# Task 03: The guides and the --update Skill Mode describe the profiles the Update Check names

## Overview

The update command guide, the Roundfix skill's `setup` reference and the
glossary describe the `Profiles:` lines and the `recommendations` member that
task 04 adds to the Update Check. The `--update` Skill Mode in the `modes`
reference shows each reported category, runs the check again after an
upgrade, and asks once which fix to run. ADR-0181 is marked amended by
ADR-0416. A test pins the skill text in both copies.

## Requirements

1. MUST add to `docs/user-guide/commands/update.md`, after the
   `Authoring models:` text and before the `Next:` text, the two `Profiles:`
   line forms of `_techspec.md` → API Contract 1 with when each appears,
   that a current category or one pinned to the shipped snapshot adds
   nothing, that the comparison uses the running binary's snapshot, writes
   nothing and changes neither `next` nor the exit code, and that no
   comparison runs when a newer Roundfix wrote the Baseline; and add to its
   JSON table the `recommendations` rows of API Contract 2: `snapshot`,
   `categories[].category`, `state` (`differs` or
   `pinned_to_older_snapshot`), `source`, `configured`, `recommended`,
   `deviation`, `repin` and `adopt`. It MUST say that `--apply` adopts every
   differing category of its scope, a category pinned to an older snapshot
   included.
2. MUST add to the Update Check paragraph of the Roundfix skill's `setup`
   reference one or two sentences naming the `Profiles:` lines, the
   `recommendations` member, the re-pin key `profiles.<category>.deviation.from`
   and the adopt command `roundfix profiles check --apply --scope <user|project>`.
3. MUST change the `## --update` section of the Roundfix skill's `modes`
   reference as API Contract 4 says: step 1 shows each `recommendations`
   category with its commands; after step 2 runs `roundfix upgrade`, the mode
   runs `roundfix update check --format json` again with the installed
   binary; and before the report a step asks once, in one question with the
   recommended option first, among adopting every reported category with
   `roundfix profiles check --apply --scope <scope> --yes` for each scope,
   re-pinning each category pinned to an older snapshot by setting
   `profiles.<category>.deviation.from` to the snapshot in the configuration
   file of its source and adopting the rest, and keeping the configuration.
   It runs only what the user chooses, re-pins before it adopts, and when the
   mode runs at session start because `update.mode` is `auto` it only reports
   the categories. MUST leave the other modes' sections byte-identical.
4. MUST change the glossary entry **Update Check** in `CONTEXT.md`, through
   `domain-modeling`, to add that it also names each configured Agent
   Selection Profile that differs from the Recommended Profile, with the
   command that adopts it, citing ADR-0416 as a link in the entry's form; no
   other glossary entry changes.
5. MUST append to `docs/adr/0181-a-configuration-is-compared-with-the-recommendation-only-where-a-person-asked.md`
   one paragraph in the form of the existing amendment lines:
   `**Amended by [ADR-0416](0416-the-update-check-names-each-configured-profile-that-differs-from-the-recommended-profile.md) (2026-11-17).**`
   followed by one or two sentences: the Update Check compares too, naming a
   deviation declared against an older snapshot as pinned to an older
   snapshot with its re-pin key and adopt command; the comparison stays
   read-only, offline and decisionless, and never runs when a Run starts.
   The rest of the ADR stays byte-identical.
6. MUST add `skills/roundfix_update_mode_recommendations_test.go` with
   `TestTheUpdateModeOffersTheProfileCommands`, calling `t.Parallel()`
   first, one subtest per copy that `roundfixLayoutSources` returns,
   asserting with `roundfixWhitespace` that the `--update` section of
   `references/modes.md` names `recommendations`,
   `roundfix profiles check --apply --scope`, `--yes` and `deviation.from`,
   and that `references/setup.md` names `Profiles:` and `recommendations`.
7. MUST regenerate, never hand-edit, with the commands of
   `_authorization.md` → Sanctioned regeneration, in this order:
   `make skills-sync`; the `-record-skill-versions` step
   (`go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`);
   and the `-record-skill-coverage` step
   (`go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`).
8. MUST change no file outside the ones this Task declares; the Governed
   Paths it changes are in `_authorization.md` `paths`.

## Subtasks

- [ ] Describe the lines and the member in the update command guide.
- [ ] Update the `setup` and `modes` references.
- [ ] Change the glossary entry and add ADR-0181's amendment line.
- [ ] Add the skill test and run the sanctioned regeneration.

## Acceptance Criteria

- [ ] The update command guide shows both `Profiles:` line forms and the `recommendations` JSON rows.
- [ ] Both copies of the `--update` Skill Mode re-check after an upgrade and ask once among adopt, re-pin and keep, never running a fix unasked.
- [ ] `CONTEXT.md` defines **Update Check** with the comparison and ADR-0416; ADR-0181 carries its amendment line.
- [ ] The skill mirror, the version record and the Behavior Surface Record are current.

## Context

- instruction: `docs/adr/0416-the-update-check-names-each-configured-profile-that-differs-from-the-recommended-profile.md`
- interface: `docs/user-guide/commands/update.md`
- interface: `.agents/skills/roundfix/references/setup.md`
- interface: `skills/roundfix/references/setup.md`
- interface: `.agents/skills/roundfix/references/modes.md`
- interface: `skills/roundfix/references/modes.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`
- interface: `docs/adr/0181-a-configuration-is-compared-with-the-recommendation-only-where-a-person-asked.md`
- creates: `skills/roundfix_update_mode_recommendations_test.go`

## Consumers

- The Roundfix skill's `modes` and `setup` references and the Behavior Surface Record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^TestTheUpdateModeOffersTheProfileCommands$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -qF -- "--- PASS: TestTheUpdateModeOffersTheProfileCommands " || { echo 'missing pass: TestTheUpdateModeOffersTheProfileCommands' >&2; exit 1; }; go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract && make skills-sync-check` — expected: exit 0. Before this Task the test does not exist, so the command fails; afterwards both copies carry the text, and the skill, mirror, version and surface records pass.
- `tr -s '[:space:]' ' ' < docs/user-guide/commands/update.md | grep -qF -- "is pinned to an older snapshot" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/update.md "is pinned to an older snapshot" >&2; exit 1; }; grep -qF -- 'pinned_to_older_snapshot' docs/user-guide/commands/update.md && { tr -s '[:space:]' ' ' < docs/adr/0181-a-configuration-is-compared-with-the-recommendation-only-where-a-person-asked.md | grep -qF -- 'Amended by [ADR-0416]' || { printf 'missing phrase in %s: %s\n' docs/adr/0181-a-configuration-is-compared-with-the-recommendation-only-where-a-person-asked.md 'Amended by [ADR-0416]' >&2; exit 1; }; }` — expected: exit 0. Before this Task the guide has no `Profiles:` lines and ADR-0181 has no amendment line.
- `grep -F -- '**Update Check**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Update Check\*\*(:| —( |$))' && awk 'found { print; exit } /^\*\*Update Check\*\*:/ { found = 1 }' CONTEXT.md | grep -qF -- 'ADR-0416'` — expected: exit 0. Before this Task the entry's definition does not cite ADR-0416; afterwards the term is defined with the comparison.

## References

- `_prd.md` → Goals (third and fourth bullets); Core Feature 4; Glossary
- `_techspec.md` → Update Check report; API Contracts 1, 2 and 4; Testing Approach (Task 3); Build Order 3; Glossary
- ADR-0416, ADR-0181, ADR-0371, ADR-0256
