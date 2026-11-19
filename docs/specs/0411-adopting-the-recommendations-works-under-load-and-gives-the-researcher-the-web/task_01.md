---
task: task_01
spec: 0411-adopting-the-recommendations-works-under-load-and-gives-the-researcher-the-web
status: pending
type: backend
complexity: medium
---

# Task 01: A profile write proves under the configured setup deadline

## Overview

`profiles configure` and `profiles check --apply` share `writeProfilesConfiguration`, which proves the tuples it writes with `roundconfig.Config{Profiles: …}` and so under the built-in 30-second setup deadline, whatever `defaults.profile_proof_setup_timeout` says. This Task adds a reader of that one key to `internal/config`, passes its value to the proof, refuses an invalid value before any proof, and makes the configuration guide, the `profiles` guide and the Roundfix skill's `profiles` reference name both commands. It answers the field reports of fiscus and fluxus of 2026-11-18, recorded in [`references/2026-11-18-maintainer-answers.md`](references/2026-11-18-maintainer-answers.md).

## Requirements

1. MUST add `internal/config/proof_setup_deadline.go` with `ProfileProofSetupTimeout(homeDir string, workDir string) (time.Duration, error)` as `_techspec.md` → Interfaces gives it and invariants 1 and 2 bound it:
   - it reads `$HOME/.roundfix/config.yml` (through `resolveHomeDir` and `userConfigRelPath`) and then `.roundfixrc.yml` at `FindGitRoot` of the resolved working directory, skipping a missing file;
   - it decodes only `defaults.profile_proof_setup_timeout` from each, with the existing `durationValue` type, ignoring every other key, so a file that another key would stop `Load` from accepting still yields its deadline;
   - a value from Project Config wins over one from User Config, and with neither the result is `DefaultProfileProofSetupTimeout`;
   - an unreadable file or a YAML or duration error returns an error naming the file; a file that does not parse as YAML returns `parse config "<path>": <parser error>` with the parser's error unchanged, so the message names the file and the line, as `_techspec.md` → Surface Transcript 2 shows; a value not greater than 0 or above `MaxProfileProofSetupTimeout` returns exactly `defaults.profile_proof_setup_timeout must be greater than 0 and at most 10m`.
2. MUST change `writeProfilesConfiguration` in `internal/cli/profiles_configure.go` as `_techspec.md` → API Contract 1 says: after `PrepareProfilesConfig` and the working-directory resolution, call `roundconfig.ProfileProofSetupTimeout(loadOptions.HomeDir, loadOptions.WorkDir)`; return its error through `printProfilesConfigureError` unchanged, before any proof; pass the value as `Defaults: roundconfig.Defaults{ProfileProofSetupTimeout: …}` of the proof configuration, whose `Profiles` and categories stay those of `profilesConfigureProofScope`. No other configuration value reaches the proof, and `internal/cli/profiles_check_apply.go` does not change.
3. MUST add `internal/config/proof_setup_deadline_test.go` with `TestTheProofSetupDeadlineReaderAgreesWithLoad`, calling `t.Parallel()` first, on a temporary home and a temporary repository marked with `gittest.MarkRepository`, with the subtests `built_in`, `user` (75s), `project` (90s), `project_wins` (user 75s, project 90s, result 90s), `another_key_load_refuses` (User Config with 75s and a backend profile whose runtime is `stale-adapter`: the reader returns 75s while `Load` fails), `zero` (0s), `above_ten_minutes` (11m) and `not_a_duration` (`soon`). Each subtest asserts the reader's value or failure and that `Load` with `WorktreeLocationOptional` fails exactly when expected and otherwise returns the reader's value.
4. MUST add to `internal/cli/profile_proof_setup_deadline_test.go`, which keeps `TestProfileProofPassesTheConfiguredSetupDeadline` and `TestProfilesShowRefusesAnInvalidSetupDeadlineTranscript` unchanged:
   - `TestProfileAdoptionPassesTheConfiguredSetupDeadline`, calling `t.Parallel()` first, with the subtests `configure/built_in` (project scope, nothing declared, 30s), `configure/project_scope_reads_user` (User Config 75s), `configure/user_scope_reads_project` (Project Config 90s), `configure/project_wins` (user scope, User Config 75s, Project Config 90s, result 90s), `check_apply/project` (`writeProfilesCheckFixture` plus Project Config 90s) and `check_apply/user` (`writeProfilesCheckFixture` plus User Config 75s). The configure subtests write a one-category fragment with `--yes --json`; the check subtests run `profiles check --apply --scope <scope> --yes --json`. Each uses `withSuccessfulProfilesConfigureProof`, requires exit 0 and at least one recorded request, and asserts that every recorded `ProbeRequest` carries the expected `SetupTimeout`.
   - `TestProfilesConfigureRefusesAnInvalidSetupDeadlineTranscript`, asserting `_techspec.md` → Surface Transcript 1 whole: a Project Config holding `defaults:` with `profile_proof_setup_timeout: 0s`, a valid one-category fragment, `profiles configure --scope user --file <fragment> --yes`, standard output empty, standard error exactly `roundfix: profiles failed: defaults.profile_proof_setup_timeout must be greater than 0 and at most 10m` and `Run 'roundfix profiles --help' for usage.` each on its own line, exit 2, a proof runner that fails the test if called, and no `.roundfix/config.yml` written under the home.
   - `TestProfilesConfigureNamesTheFileAndLineOfABrokenConfig`, calling `t.Parallel()` first, asserting `_techspec.md` → Surface Transcript 2 whole: a Project Config holding the three lines `defaults:`, `  profile_proof_setup_timeout: 90s` and `broken: a: b`, the fragment of the test above, `profiles configure --scope user --file <fragment> --yes`, standard output empty, standard error exactly `roundfix: profiles failed: parse config ` followed by the Project Config path quoted with `strconv.Quote` and `: yaml: line 3: mapping values are not allowed in this context`, then `Run 'roundfix profiles --help' for usage.`, each on its own line, exit 2, a proof runner that fails the test if called, and no `.roundfix/config.yml` written under the home.
   - The file's header comment names `profiles configure` and `profiles check --apply` in its `Boundary IN` line.
   `TestProfilesConfigureProofScope` and `TestProfilesConfigureChangeSummary` stay unchanged and keep passing; so does every other test that pins the proof of a profile write on the default branch when this Task runs.
5. MUST make the changes of `_techspec.md` → API Contract 2 word for word, in `docs/user-guide/configuration.md` (section `### Profile proof setup deadline`), `docs/user-guide/commands/profiles.md` and `.agents/skills/roundfix/references/profiles.md` (section `## Profile proof setup deadline`), carrying forward every other sentence of those paragraphs.
6. MUST regenerate, in this order and with no hand edit of a mirror, version field or record:
   - the mirror `skills/roundfix/references/profiles.md`, with `make skills-sync`;
   - the Roundfix skill's version, with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`, which raises both version fields of `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` to the next version and records it in `skills/testdata/owned-skill-versions.json`;
   - `make skills-sync` again, then `make skills-sync-check`;
   - as the last step, after the last guide edit, the Behavior Surface Record, with `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`, which rewrites `docs/references/behavior-surfaces.json`.
   Specs 0407 and 0408 deliver before this Spec and may have raised the same version fields and changed the same records; each command runs on the tree this Task starts from.
7. MUST NOT change `internal/agent/fallback.go`, `internal/baselineacp/`, `internal/cli/profiles_validate.go`, `internal/cli/profiles_check_apply.go`, any other skill or skill reference, any Governed Path outside `.agents/skills/roundfix/references/profiles.md`, `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md`, or the `### QA settlement` section of any skill.

## Subtasks

- [ ] Add the deadline reader and its agreement test.
- [ ] Pass the deadline in `writeProfilesConfiguration` and add the adoption and refusal tests.
- [ ] Edit the configuration guide, the `profiles` guide and the canonical `profiles` reference.
- [ ] Synchronize the mirror, record the skill version, and re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] Every proof of `profiles configure` and `profiles check --apply` carries the configured deadline, Project Config winning, for a write to either scope.
- [ ] An invalid deadline refuses `profiles configure` before any proof with the transcript of Surface Transcript 1, a file that does not parse refuses it with the file and the line of Surface Transcript 2, and an untouched invalid category still does not stop a rewrite.
- [ ] The reader agrees with `Load` wherever `Load` succeeds.
- [ ] The three documents name both commands, the mirror equals the canonical reference, the Roundfix skill's version is recorded, and the Behavior Surface Record is current.

## Context

- instruction: `internal/config/config.go`
- instruction: `internal/cli/profiles_check_apply.go`
- creates: `internal/config/proof_setup_deadline.go`
- creates: `internal/config/proof_setup_deadline_test.go`
- interface: `internal/cli/profiles_configure.go`
- interface: `internal/cli/profile_proof_setup_deadline_test.go`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/user-guide/commands/profiles.md`
- interface: `.agents/skills/roundfix/references/profiles.md`
- interface: `skills/roundfix/references/profiles.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The proof configuration of a profile write, read by `proveProfileSelections` — consumers: `./internal/cli`, `./internal/config`
- The Roundfix skill's `profiles` reference, its version and the Behavior Surface Record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `grep -q '^func ProfileProofSetupTimeout(homeDir string, workDir string) (time.Duration, error)' internal/config/proof_setup_deadline.go && grep -q '^func TestTheProofSetupDeadlineReaderAgreesWithLoad(' internal/config/proof_setup_deadline_test.go && grep -q '^func TestProfileAdoptionPassesTheConfiguredSetupDeadline(' internal/cli/profile_proof_setup_deadline_test.go && grep -q '^func TestProfilesConfigureRefusesAnInvalidSetupDeadlineTranscript(' internal/cli/profile_proof_setup_deadline_test.go && grep -q '^func TestProfilesConfigureNamesTheFileAndLineOfABrokenConfig(' internal/cli/profile_proof_setup_deadline_test.go && go test -count=1 -timeout 15m ./internal/config && go test -count=1 -timeout 20m -run '^(TestProfileAdoptionPassesTheConfiguredSetupDeadline|TestProfilesConfigureRefusesAnInvalidSetupDeadlineTranscript|TestProfilesConfigureNamesTheFileAndLineOfABrokenConfig|TestProfileProofPassesTheConfiguredSetupDeadline|TestProfilesShowRefusesAnInvalidSetupDeadlineTranscript|TestProfilesConfigure.*|TestProfilesCheck.*|TestProfileProof.*)$' ./internal/cli` — expected: exit 0. Before this Task the reader and the three tests do not exist, so the command fails at the first `grep`; afterwards the whole `internal/config` package and every profile-write and proof-deadline test of `internal/cli` pass.
- `for file in docs/user-guide/configuration.md .agents/skills/roundfix/references/profiles.md skills/roundfix/references/profiles.md; do text="$(tr -s '[:space:]' ' ' < "$file")"; for phrase in 'Setup, .profiles configure. and .profiles check --apply.' 'and .profiles check --apply. read the key from User Config and Project Config before they prove, with Project Config winning,' 'even when another key of those files would stop loading, and an invalid value refuses the command before any proof.'; do printf '%s\n' "$text" | grep -qE -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done; done; tr -s '[:space:]' ' ' < docs/user-guide/commands/profiles.md | grep -qE -- 'profiles configure. and .profiles check --apply. prove under the same key, read from User Config and Project Config with Project Config winning.' || { echo 'profiles guide lacks the sentence' >&2; exit 1; }; make skills-sync-check && go test -count=1 -timeout 30m ./skills && go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task none of the documents holds the sentences, so the command fails; afterwards the documents, the mirror, the skill version record and the Behavior Surface Record agree, and the skill and documentation contract suites pass.

## References

- `_prd.md` → Goals (first and fourth bullets); Core Features 1 and 4; Success Metric 1; Success Metric 2; Non-Goals
- `_techspec.md` → System Architecture 1 and 2; Interfaces; invariants 1-3; API Contract 1; API Contract 2; Surface Transcript 1; Surface Transcript 2; Testing Approach; Build Order 1
- [`references/2026-11-18-characterization.md`](references/2026-11-18-characterization.md) sections 1 to 4
- ADR-0298
