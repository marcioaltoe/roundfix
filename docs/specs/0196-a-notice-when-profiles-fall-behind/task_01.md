---
task: task_01
spec: 0196-a-notice-when-profiles-fall-behind
status: completed
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

Implemented this Task's slice only: a configured profile accepts a dated,
reasoned Profile Deviation, resolution carries a copied record, and configure
fragments write, preview and return it in JSON. Whole-profile replacement
without a deviation removes the old record. Built-in and legacy-derived
profiles carry none, and invocation overrides clear it. Recommendation
comparison remains for task_02.

Pre-change inspection found that `decodeProfile` admitted only `preferred` and
`fallbacks`, and neither new test file existed. No existing top-level test was
renamed, removed or edited, and no exported function signature changed.

### Focused evidence by acceptance criterion

The focused command was
`GOCACHE=/tmp/roundfix-0196-task01-cache rtk proxy go test ./internal/config ./internal/cli -run 'Deviation|TestProfilesConfigure' -count=1`.
It exited 0 after the final implementation edits, including all seven new
named tests and the existing configure tests.

| Acceptance criterion | Implementation and focused-check evidence |
| --- | --- |
| User and Project Config load a valid deviation and resolution carries it | `TestProfileDeviationLoadsFromUserAndProjectConfig` covers both scopes, source attribution, optional-category inheritance, copied resolution records, invocation clearing, and Project Config replacement of a User Config deviation. |
| Every malformed deviation names its path | `TestProfileDeviationRejectsMalformedValues` separately covers an unknown key, missing `from`, missing `reason`, whitespace-only reason, impossible calendar date and non-mapping value. Additional subtests cover duplicate keys, non-string reason, non-scalar date and a timestamp with a time. Each exercises both config loading and fragment parsing and asserts the offending path. |
| Profiles without a deviation retain existing behavior | `TestProfileWithoutDeviationLoadsUnchanged` checks selection tuples and absence on configured, built-in and legacy profiles. `TestProfilesConfigureOutputIsUnchangedWithoutADeviation` compares the complete text and JSON bytes against the existing public output contract. |
| Project configure writes a deviation and later removes it | `TestProfilesFragmentPersistsTheDeviation` checks CategoryChange and proposal copies, persists and reloads the record. `TestReplacingAProfileDropsItsDeviation` checks whole-profile replacement. `TestProfilesConfigurePreviewsAndWritesADeviation` uses the existing fake runner and temporary home/repository to exercise `--scope project --file <fragment> --yes`, preview text, JSON, YAML ordering after `fallbacks`, and removal by a later fragment. No real adapter is used. |
| Existing configure tests pass unchanged | The focused command includes every `TestProfilesConfigure*` test, including `TestProfilesConfigureChangeSummary` and `TestProfilesConfigureExitCodes`. A byte comparison against HEAD confirmed `internal/cli/profiles_configure_test.go` is unchanged. |
| Guides, glossary, skill and mirror describe Profile Deviation; sync check passes | The configuration guide includes the example, two validation rules and the older-release refusal sentence; the command guide describes fragment persistence; CONTEXT defines the term. Both Roundfix skills contain the short paragraph under `### Recommendation check`, beside Agent selection. A byte comparison confirms the mirror matches and `### QA settlement` is unchanged. The incremental check ran `skills-sync-check` successfully. |

### Sabotage evidence

- Refusal gate: temporarily replaced the unknown-key refusal in
  `decodeProfileDeviation` with `continue`. With the task-local cache,
  `go test ./internal/config -run '^TestProfileDeviationRejectsMalformedValues$/unknown_deviation_key$' -count=1`
  exited 1. `TestProfileDeviationRejectsMalformedValues/unknown_deviation_key`
  failed with `expected path refusal, got <nil>`. Restored the decoder.
- Write gate: temporarily disabled the deviation branch of
  `profileWithDeviationYAMLNode`. With the same cache,
  `go test ./internal/config -run '^TestProfilesFragmentPersistsTheDeviation$' -count=1`
  exited 1. `TestProfilesFragmentPersistsTheDeviation` failed with
  `written deviation = <nil>, want {From:2026-09-30 Reason:Keep the validated model}`.
  Restored the writer. The focused command above then exited 0, and was run
  again successfully after the final edits.

### Skill regeneration and additional path

Raised both Roundfix skill version fields from `0.0.4` to `0.0.5` as required by
`docs/agents/specific-repository.md`. Additional Task path declared here:
`skills/testdata/owned-skill-versions.json`, the generated owned-skill version
record. The command
`GOCACHE=/tmp/roundfix-0196-task01-cache rtk proxy go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
exited 0 and rewrote that file to record the final skill content.

`rtk make skills-sync` exited 0. It recreates the shipped bundle, so the files
it copied are named below. Only `skills/roundfix/SKILL.md` gained changed bytes;
the other copied files remain byte-identical to HEAD.

```text
skills/archive-spec/SKILL.md
skills/brainstorming/SKILL.md
skills/business-analyst/SKILL.md
skills/council/SKILL.md
skills/council/assets/synthesis-template.md
skills/council/references/archetypes.md
skills/council/references/debate-protocols.md
skills/evidence-gate/SKILL.md
skills/implement-spec/SKILL.md
skills/implement-task/SKILL.md
skills/qa-gate/SKILL.md
skills/roundfix/SKILL.md
skills/roundfix/agents/openai.yaml
skills/setup-context-driven/SKILL.md
skills/write-idea/SKILL.md
skills/write-idea/references/idea-template.md
skills/write-idea/references/opportunity-scan.md
skills/write-prd/SKILL.md
skills/write-prd/references/prd-template.md
skills/write-tasks/SKILL.md
skills/write-tasks/references/task-template.md
skills/write-techspec/SKILL.md
skills/write-techspec/references/techspec-template.md
```

`GOCACHE=/tmp/roundfix-0196-task01-cache rtk make baseline-digests` exited 0
after the final skill edit and synchronization. It reported `changed:false`
and “derived artifacts already match their canonical sources”; it produced no
changed derived files.

### Incremental check and handoff boundary

- `GOCACHE=/tmp/roundfix-0196-task01-cache rtk proxy make verify-incremental`
  exited 0 with host permissions. Formatting, vet, repository tests,
  `skills-sync-check`, skill validation and build passed. The output is in
  `/tmp/roundfix-0196-task01-incremental.log` for this session.
- The first direct Go check could not access the default Go cache; focused
  checks used the task-local cache thereafter. The first sandboxed incremental
  attempt was interrupted by a network restriction naming `cafe.github.com`
  and yielded no usable completion result. The host-permission rerun above
  supplied the incremental evidence; no repository check was weakened.
- `git diff --check` exited 0. The only pre-existing Task-file difference was
  the Daemon's `status: in_progress`; that value remains untouched.

The authored Verification commands were not run. Task status, subtasks and
acceptance checkboxes remain unchanged for Daemon settlement. No other Task
file or Task Graph was edited, and no commit, push or Pull Request was made.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `skills/testdata/owned-skill-versions.json`

## Carry-forward provenance

- Source Run: `run_20260930T224533Z_22a572af678fb8c4`
- Source commit: `2a9207735fe4a0da788c89fcacae44f19da1ecbc`
