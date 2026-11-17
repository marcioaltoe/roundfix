---
task: task_02
spec: 0403-the-update-check-names-each-profile-that-differs-from-a-refreshed-recommendation
status: pending
type: backend
complexity: high
---

# Task 02: The Recommended Profile becomes the 2026-11-17 snapshot, and this repository's profiles follow it

## Overview

The Recommended Profile the binary ships becomes the table the maintainer
approved on 2026-11-17, and the built-in `review` profile follows it. The
model selection reference states the new snapshot with the evidence and
sources it rests on, the guides and the Roundfix skill's `profiles` reference
list the new built-ins, and this repository's `.roundfixrc.yml` takes the new
`review`, `docs` and `chore` profiles with comments that describe them. A new
test keeps this repository's profiles equal to the recommendation.

## Requirements

1. MUST set `ModelRecommendationSnapshotVersion` to `2026-11-17` and the
   Recommended Profile in `internal/config/recommendations.go` to the table
   of `references/2026-11-17-maintainer-answers.md`: `review`
   `codex / gpt-6-luna / max` then `codex / gpt-6.1-sol / high`; `docs`
   `codex / gpt-6.1-terra / high` then `claude / sonnet / high`; `chore`
   `codex / gpt-6-luna / max` then `claude / sonnet / high`; every other
   category unchanged (`_techspec.md` → Recommended Profile, invariant 1).
   Each selection keeps a one-sentence rationale without a benchmark figure
   (ADR-0180); the `docs` and `chore` entries stop sharing one helper when
   their selections differ.
2. MUST move "built-in review default" in `internal/agent/catalog.go` from
   `gpt-5.6-luna` to `gpt-6-luna`, and update the matching description in
   `TestModelCatalogsExposeOrderedPickerData`.
3. MUST update `TestRecommendedProfileCoversEveryWorkCategory` to the
   approved table, and add `TestRecommendedProfileUsesNoUltraEffort` to
   `internal/config/recommended_profile_test.go`, calling `t.Parallel()`
   first, asserting that no selection of any category's Recommended Profile
   has the effort `ultra`.
4. MUST update the `docs` rows of
   `TestProfilesShowJSONIsSchemaV2WithoutBenchmarkFields` in
   `internal/cli/profiles_show_recommended_test.go` to
   `codex / gpt-6.1-terra / high` then `claude / sonnet / high`. This is a
   declared change of that test's expectation.
5. MUST add `internal/config/repository_profiles_test.go` with
   `TestRepositoryProfilesFollowTheRecommendedProfile`, calling
   `t.Parallel()` first, which reads this repository's `.roundfixrc.yml`
   with the profile parser Project Config uses, so a key it rejects fails the
   test, with a temporary home directory and no network, and asserts that every profile it defines equals `RecommendedProfile` for its
   category and carries no deviation, one subtest per category.
6. MUST change, in `.roundfixrc.yml`, only the `review`, `docs` and `chore`
   profiles, to the table above, and the comment block above `profiles:`:
   it dates the revision 2026-11-17, describes each group of categories with
   its selection and reason, drops the `gpt-5.6-sol` session measurement and
   every mention of `gpt-5.6-luna`, names Fable 5.2 where it named Fable, and
   keeps the subscription rule (no `gpt-*` or Claude model through
   OpenRouter, ADR-0235) and the pointer to
   `docs/references/model-selection.md`. No other key changes.
7. MUST rewrite the `## Current shipped snapshot` section of
   `docs/references/model-selection.md` for the 2026-11-17 snapshot: the line
   `Updated 2026-11-17.`, the table with every category's selections and
   rationales, and the published evidence of `gpt-6.1-terra`, `gpt-6-luna`,
   Haiku 5.5 and Fable 5.2 with their release dates and sources from
   `references/2026-11-17-model-evidence.md`. The 2026-09-30 table moves under
   `## Historical sections — not current selection advice` with its date.
8. MUST change the built-in lists and the `review` blocks of the profile
   examples in `docs/user-guide/usage.md`, `docs/user-guide/configuration.md`
   and the Roundfix skill's `profiles` reference to
   `codex / gpt-6-luna / max`, and the date "dated `2026-09-30`" of that
   reference to `2026-11-17`.
9. MUST regenerate, never hand-edit, with the commands of
   `_authorization.md` → Sanctioned regeneration, in this order:
   `make skills-sync`; the `-record-skill-versions` step
   (`go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`);
   and the `-record-skill-coverage` step
   (`go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`).
10. MUST change no file outside the ones this Task declares; the Governed
    Paths it changes are in `_authorization.md` `paths`.

## Subtasks

- [ ] Set the new snapshot and its rationales.
- [ ] Update the recommendation, catalog and `profiles show` tests and add the two new tests.
- [ ] Rewrite the model selection reference's current snapshot.
- [ ] Update the guides, the skill reference and `.roundfixrc.yml`.
- [ ] Run the sanctioned regeneration.

## Acceptance Criteria

- [ ] `RecommendedProfile` returns the approved table for every category, dated `2026-11-17`, with no `ultra` effort, and the built-in `review` profile follows it.
- [ ] The model selection reference states the snapshot, every selection and its sources.
- [ ] The guides and both copies of the `profiles` reference list the new built-ins.
- [ ] `.roundfixrc.yml` defines every profile as the Recommended Profile, and its comments name no `gpt-5.6-luna` and no 2026-09-30 measurement.

## Context

- instruction: `docs/specs/0403-the-update-check-names-each-profile-that-differs-from-a-refreshed-recommendation/references/2026-11-17-maintainer-answers.md`
- instruction: `docs/specs/0403-the-update-check-names-each-profile-that-differs-from-a-refreshed-recommendation/references/2026-11-17-model-evidence.md`
- instruction: `docs/adr/0180-built-in-selections-derive-from-one-dated-recommended-profile.md`
- interface: `internal/config/recommendations.go`
- interface: `internal/config/recommended_profile_test.go`
- creates: `internal/config/repository_profiles_test.go`
- interface: `internal/agent/catalog.go`
- interface: `internal/agent/agent_test.go`
- interface: `internal/cli/profiles_show_recommended_test.go`
- interface: `.roundfixrc.yml`
- interface: `docs/references/model-selection.md`
- interface: `docs/user-guide/usage.md`
- interface: `docs/user-guide/configuration.md`
- interface: `.agents/skills/roundfix/references/profiles.md`
- interface: `skills/roundfix/references/profiles.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- `config.RecommendedProfile`, `config.ModelRecommendationSnapshotVersion` and the built-in profiles — consumers: `./internal/config`, `./internal/agent`, `./internal/cli`, `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestRecommendedProfileCoversEveryWorkCategory|TestRecommendedProfileUsesNoUltraEffort|TestRepositoryProfilesFollowTheRecommendedProfile|TestRecommendedModelsAreInTheModelCatalog)$' ./internal/config 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestRecommendedProfileCoversEveryWorkCategory TestRecommendedProfileUsesNoUltraEffort TestRepositoryProfilesFollowTheRecommendedProfile TestRepositoryProfilesFollowTheRecommendedProfile/review TestRepositoryProfilesFollowTheRecommendedProfile/docs TestRepositoryProfilesFollowTheRecommendedProfile/chore; do printf '%s\n' "$out" | grep -qF -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; grep -qF -- 'ModelRecommendationSnapshotVersion = "2026-11-17"' internal/config/recommendations.go && go test -count=1 -timeout 30m ./internal/config ./internal/agent ./internal/cli && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the new tests do not exist and the snapshot is dated 2026-09-30, so the command fails; afterwards the snapshot, this repository's profiles and every package that reads the recommendation pass, the reference and guide contracts included.
- `! grep -qF -- 'gpt-5.6-luna' .roundfixrc.yml && ! grep -qF -- '18.8 minutes' .roundfixrc.yml && grep -qF -- 'Fable 5.2' .roundfixrc.yml && { tr -s '[:space:]' ' ' < docs/references/model-selection.md | grep -qF -- "Updated 2026-11-17." || { printf 'missing phrase in %s: %s\n' docs/references/model-selection.md "Updated 2026-11-17." >&2; exit 1; }; } && go test -count=1 -timeout 20m ./skills && make skills-sync-check` — expected: exit 0. Before this Task `.roundfixrc.yml` names `gpt-5.6-luna` and the 2026-09-30 measurement; afterwards its comments and the reference describe the new snapshot, and the mirror and the records pass.

## References

- `_prd.md` → Goals (first and fifth bullets); Core Features 1, 4 and 5; Success Metric 1
- `_techspec.md` → Recommended Profile; Interfaces (invariant 1); API Contract 5; Testing Approach (Task 2); Build Order 2
- [`references/2026-11-17-maintainer-answers.md`](references/2026-11-17-maintainer-answers.md)
- [`references/2026-11-17-model-evidence.md`](references/2026-11-17-model-evidence.md)
- ADR-0180, ADR-0235, ADR-0256
