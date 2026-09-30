---
task: task_01
spec: 0196-a-notice-when-profiles-fall-behind
status: pending
type: backend
complexity: medium
---

# Task 01: A profile can carry a Profile Deviation

## Overview

A configured Agent Selection Profile accepts exactly two keys, `preferred` and `fallbacks`, and `decodeProfile` in `internal/config/profiles.go` refuses any other. A repository that keeps a model on purpose has nowhere to say so. This Task adds the key `deviation`, with the snapshot date it was declared against and a reason, carries it on the loaded profile, and lets `roundfix profiles configure` write it (ADR-0181). It is verifiable on its own: a configuration with a deviation loads, a malformed one is refused, and a fragment with one is written. Nothing reads the deviation yet; task_02 does.

## Requirements

1. MUST add `ProfileDeviation` and the `Deviation` field on `ProfileEntry` and `ResolvedProfile` as `_techspec.md` → Interfaces states, and decode the key `deviation` under the rules of `_techspec.md` → Data Models. Every refusal MUST name the path of the offending key in the existing message style.
2. MUST carry the deviation from the User Config or Project Config entry that defines the profile to the resolved profile, and copy it wherever a `ProfileEntry` is cloned or overlaid. A built-in entry, a legacy-derived entry and an entry with an invocation override MUST carry none.
3. MUST keep a profile without `deviation` decoding, resolving and rendering exactly as today. It MUST rename or remove no top-level test and change no existing exported function signature.
4. MUST make a `profiles configure` fragment accept a deviation, carry it on `CategoryChange`, and write it under the profile after `fallbacks`. A category replaced by a fragment profile without a deviation MUST lose the one it had.
5. MUST make `roundfix profiles configure` print `Deviation: from <date> — <reason>` under the profile in its preview, and add `deviation` to that profile in its JSON, only when the profile carries one. For a profile without one, its text and JSON MUST stay byte-identical; the existing tests of `profiles configure` prove it and MUST pass unchanged.
6. MUST put the new tests in `internal/config/profile_deviation_test.go` and `internal/cli/profiles_deviation_test.go`. The malformed cases MUST be separate subtests: an unknown deviation key, a missing `from`, a missing `reason`, an empty `reason`, a `from` that is not a calendar date, and a deviation that is not a mapping. The CLI tests MUST use the existing fake runner and temporary home and repository, never a real adapter.
7. MUST prove each new gate can fail. The Result MUST record one sabotage for the refusal test (for example accepting an unknown deviation key) and one for the write test (for example not writing the deviation), each with the test that failed, and that the code was restored.
8. MUST document the key in `docs/user-guide/configuration.md`, with an example, its two rules, and the sentence that a Roundfix older than this release refuses a configuration that uses it. It MUST say in the `profiles` section of `docs/user-guide/commands.md` that `profiles configure` writes a Profile Deviation its fragment carries.
9. MUST add the glossary term **Profile Deviation** to `CONTEXT.md`: the dated, reasoned record on a configured Agent Selection Profile that its difference from the Recommended Profile is deliberate, which holds for the snapshot it names.
10. MUST create the heading `### Recommendation check` in `.agents/skills/roundfix/SKILL.md` with one short paragraph on the Profile Deviation. It MUST put the skill text under the heading `### Recommendation check` of `.agents/skills/roundfix/SKILL.md`, add no text inside the `### QA settlement` section, then run `make skills-sync` and `make baseline-digests`, and name in the Result every file either command rewrote.

## Subtasks

- [ ] Decode, carry and clone the deviation.
- [ ] Write it through the fragment path and show it in the configure preview and JSON.
- [ ] Add the tests, each malformed case separate.
- [ ] Document the key, add the glossary term and the skill heading, then sync the mirror and the digests.
- [ ] Record one sabotage per new gate in the Result.

## Acceptance Criteria

- [ ] A Project Config and a User Config profile with a valid `deviation` load, and the resolved profile carries it.
- [ ] Each malformed deviation is refused with a message that names its path.
- [ ] A profile without `deviation` loads and renders exactly as before.
- [ ] `roundfix profiles configure --scope project --file <fragment> --yes` writes the fragment's deviation, and a later fragment without one removes it.
- [ ] The existing `profiles configure` tests pass without a change.
- [ ] The guides, the glossary, the skill and its mirror name the Profile Deviation, and `make skills-sync-check` passes.

## Context

- interface: `internal/config/profiles.go`
- interface: `internal/config/profile_config.go`
- interface: `internal/cli/profiles_configure.go`
- creates: `internal/config/profile_deviation_test.go`
- creates: `internal/cli/profiles_deviation_test.go`
- instruction: `internal/cli/profiles_configure_test.go`
- instruction: `docs/adr/0181-a-configuration-is-compared-with-the-recommendation-only-where-a-person-asked.md`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/user-guide/commands.md`
- interface: `CONTEXT.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestProfileDeviationLoadsFromUserAndProjectConfig|TestProfileDeviationRejectsMalformedValues|TestProfileWithoutDeviationLoadsUnchanged|TestProfilesFragmentPersistsTheDeviation|TestReplacingAProfileDropsItsDeviation|TestProfilesConfigurePreviewsAndWritesADeviation|TestProfilesConfigureOutputIsUnchangedWithoutADeviation|TestProfilesConfigureChangeSummary|TestProfilesConfigureExitCodes)$" ./internal/config ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestProfileDeviationLoadsFromUserAndProjectConfig TestProfileDeviationRejectsMalformedValues TestProfileWithoutDeviationLoadsUnchanged TestProfilesFragmentPersistsTheDeviation TestReplacingAProfileDropsItsDeviation TestProfilesConfigurePreviewsAndWritesADeviation TestProfilesConfigureOutputIsUnchangedWithoutADeviation TestProfilesConfigureChangeSummary TestProfilesConfigureExitCodes; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the seven new named tests exists, so the command fails.
- `for pair in "docs/user-guide/configuration.md|Profile Deviation" "docs/user-guide/configuration.md|deviation:" "docs/user-guide/commands.md|Profile Deviation" ".agents/skills/roundfix/SKILL.md|### Recommendation check" ".agents/skills/roundfix/SKILL.md|Profile Deviation" "skills/roundfix/SKILL.md|Profile Deviation" "CONTEXT.md|**Profile Deviation**"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done && diff -r .agents/skills/roundfix skills/roundfix >/dev/null && make skills-sync-check` — expected: exit 0; before this Task no document names the Profile Deviation, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 3; User Story 3; Core Feature 1; Declared breaks
- [_techspec.md](_techspec.md) — Interfaces; Data Models; API Contract 3; Testing Approach 1; Build Order 1
- ADR-0181; ADR-0049

## Result
