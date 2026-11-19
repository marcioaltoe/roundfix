---
task: task_01
spec: 0412-a-project-declares-whether-an-open-authoring-question-waits-or-follows-the-author
status: completed
type: backend
complexity: high
---

# Task 01: A project declares its Authoring Question Preference, and the status view reports it

## Overview

Project Config gains the key `authoring.open_questions`, `wait` or `follow`, that only a project sets: an absent key means `wait`, a value in User Config is ignored with one warning, and any other value is refused by name. The status view of `roundfix follow-up --status` reads the key at the default branch head, as the local refs hold it, prints the **Authoring Question Preference** in force with its source in one new line before `Sources:`, and carries it in the new JSON field `authoring_questions`. The configuration guide, the `follow-up` guide and the Roundfix skill's `follow-up` reference describe the key and the line.

## Requirements

1. MUST add the key `authoring.open_questions` to Project Config as `_techspec.md` → Preference key, Interfaces and Data Models say: the constants `OpenQuestionsWait` (`wait`) and `OpenQuestionsFollow` (`follow`), the type `Authoring` with `OpenQuestions` and `OpenQuestionsDeclared`, the field `Authoring` of the resolved `Config`, and the built-in value `wait` with `OpenQuestionsDeclared` false. The `authoring` section holds no other key, and an unknown key under it stays refused by the strict decode.
2. MUST validate and filter the key in `prepareAuthoringConfig`, beside `delivery.require_defect_class` (`_techspec.md` → API Contract 1):
   - in User Config the key is removed before validation, even when its value is invalid, with this line on standard error once per load: `config: authoring.open_questions in User Config is ignored; set authoring.open_questions in Project Config`;
   - in Project Config a value that is not the string `wait` or `follow` is refused when the configuration loads with `authoring.open_questions must be wait or follow`, after the `parse config "<label>": ` prefix every configuration error carries;
   - the default configuration that `roundfix init` writes, `DefaultProjectConfigYAML` included, does not hold the key, and no command writes it (`_techspec.md` → Invariant 3).
3. MUST create `internal/cli/follow_up_status_authoring.go` with the reader of `_techspec.md` → Preference reader: it resolves the default branch head with `preflight.ResolveDefaultBranchHead` and reads the configuration at that commit with `roundconfig.DeliveryConfigAtCommit`, through the Git runner the status view already receives; it never fetches, never writes and never fails the command. It fills `followUpAuthoringQuestions` exactly as `_techspec.md` → Data Models, "The status fields", says: `preference`, `source` (`project-config` or `default`), `ref`, `checkout` (the checkout's value only when it differs from `preference`) and `detail` (empty, `no default branch head was found`, or `the default branch's Project Config was not read: <error>`; a head that does not resolve or a configuration that cannot be read gives `wait` and `default`).
4. MUST add to the status document, in `internal/cli/follow_up_status.go`, the field `AuthoringQuestions followUpAuthoringQuestions` with the JSON name `authoring_questions`, immediately after `Maintainer`, and print, immediately after the `Waiting on the maintainer` block and before `Sources:`, the first line of `_techspec.md` → API Contract 2: `Open authoring questions: <preference> (<Project Config|default>)<checkout part><detail part>`, where the checkout part is empty or `; this checkout's Project Config gives <value>, which counts only once the default branch holds it` and the detail part is empty or `; ` and the `detail`. The `Delegated decisions` block is task_02's. The list view, every flag, the help text and every exit code stay as they are (`_techspec.md` → Invariant 6).
5. MUST create `internal/config/authoring_questions_test.go` with `TestTheAuthoringQuestionPreferenceIsAProjectConfigValue`, which calls `t.Parallel()` first, uses `ResolveConfigProposal` and `Load` with a temporary home and repository as `internal/config/require_defect_class_test.go` does, and runs these subtests:
   - `builtin_waits`: the built-in configuration gives `wait`, not declared;
   - `absent_key_waits`: a Project Config without the key gives `wait`, not declared;
   - `project_config_wait` and `project_config_follow`: each value in Project Config is read and declared;
   - `user_config_ignored_once`: `follow` in User Config, and in a second load `maybe` in User Config, each gives `wait`, not declared, and exactly one warning line, the text of requirement 2;
   - `refuses_maybe`, `refuses_a_number`, `refuses_a_list` and `refuses_a_mapping`: `maybe`, `1`, `[follow]` and `{value: follow}` in Project Config are each refused with an error that ends with `authoring.open_questions must be wait or follow`;
   - `refuses_an_unknown_key`: `authoring.other: true` in Project Config is refused;
   - `default_configuration_omits_the_key`: neither default configuration holds `open_questions`.
6. MUST create `internal/cli/follow_up_status_authoring_test.go` with `TestTheStatusViewReportsTheAuthoringQuestionPreference`, which uses `newStatusQueueFixture` of `internal/cli/follow_up_status_queue_test.go` and the Git runner parameter of `runFollowUpStatusWithHistoryGit`, starts no remote and no network, and asserts in each subtest the whole `Open authoring questions` line of the text report and every field of `authoring_questions` in the JSON report:
   - `default_without_the_key`: `main` holds no key: `wait (default)`, `ref` naming `main`;
   - `follow_committed_on_main`: `main` commits `follow`: `follow (Project Config)`;
   - `follow_only_in_the_working_tree`: `follow` is only in the working tree of `main`: `wait (default)` with the checkout part naming `follow`, and `checkout` `follow`;
   - `follow_only_on_a_feature_branch`: `follow` is committed on a branch the checkout is on while `main` holds no key: the same line as the previous subtest;
   - `remote_tracking_ref_wins`: `refs/remotes/origin/main`, written with `git update-ref`, commits `follow` while the local `main` holds no key: `follow (Project Config)` and `ref` `refs/remotes/origin/main`;
   - `no_default_branch_head`: the checkout is on a branch named `work`, with no `main` and no remote: `wait (default); no default branch head was found`;
   - `unreadable_project_config`: `main` commits `authoring.open_questions: maybe` while the checkout is on another branch whose Project Config is valid: exit 0, `wait (default)` and a detail that begins `the default branch's Project Config was not read: `.
7. MUST add the `Open authoring questions` line, as the status view prints it for each fixture, before `Sources:` in every test that pins the whole status text on the default branch when the Task runs; at the authoring base these are `TestTheStatusViewShowsTheSpecQueueTable` and `TestTheStatusViewSurfaceTranscripts` with its constants `statusQueueSurface1` and `statusQueueSurface2`, in `internal/cli/follow_up_status_queue_test.go`, and `TestFollowUpStatusNamesTheManualClassification`, in `internal/cli/follow_up_status_test.go`. No other assertion of those tests changes, and `TestFollowUpStatusWritesNothing` passes unchanged.
8. MUST write the documentation of the line, word for word from `_techspec.md` → API Contracts 11 and 13:
   - `docs/user-guide/configuration.md` gains, at its end, the section `## Open authoring questions` of API Contract 13;
   - `docs/user-guide/commands/follow-up.md` gains, immediately before `## Text and JSON reports`, the section `## Open authoring questions and Delegated Decisions` of API Contract 13 without the paragraph that begins "Each Delegated Decision is read from" and without the line that begins "Each delegated decision carries these JSON fields:", both task_02's; its sentence that lists the text sections names `Open authoring questions` after `Waiting on the maintainer`; its sentence "The status report carries these JSON fields:" ends with `` `maintainer`, `authoring_questions`. ``; and its JSON example of the status report ends with `"maintainer":[],"authoring_questions":{"preference":"wait","source":"default","ref":"refs/heads/main","checkout":"","detail":""}}`;
   - the canonical `.agents/skills/roundfix/references/follow-up.md` gains, immediately before `## Spec queue`, the lines of API Contract 11 without the line that begins "Each delegated decision carries these JSON fields:", and its sentence "The status report carries these JSON fields:" ends with `` `maintainer`, `authoring_questions`. ``.
9. MUST add to `TestTheFollowUpGuidesNameTheJSONFieldsAndMessageTypes`, in `internal/cli/live_documents_docscontract_test.go`, the row `{"The authoring questions report carries these JSON fields:", reflect.TypeOf(followUpAuthoringQuestions{})}`, and change no other row.
10. MUST regenerate, never by hand: the mirrors with `make skills-sync`, then the owned-skill versions with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`, which raises both version fields of `roundfix` on the tree the Task runs on, then `make skills-sync` again; and, as the Task's last edit, after the last change to a user-guide page, the Behavior Surface Record with `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
11. MUST NOT change `.roundfixrc.yml`, any other skill or skill file, the `### QA settlement` section of any skill, any Baseline module, template, profile or golden, any guide under `docs/agents/`, `CONTEXT.md`, any ADR, `internal/speccheck/coherence.go`, the corpus pins, `docs/references/coverage-record.json`, `docs/references/skill-coverage.json`, the Makefile, a workflow, `go.mod`, `go.sum` or a lint configuration. Of the Governed Paths it changes only `.agents/skills/roundfix/references/follow-up.md` and both copies of the Roundfix `SKILL.md`, all three listed in `_authorization.md`.

## Subtasks

- [ ] Add the key with its validation, its User Config warning and its tests.
- [ ] Add the reader at the default branch head, the JSON field and the status line, with their tests.
- [ ] Move the tests that pin the status text to the new line.
- [ ] Write the configuration guide section, the `follow-up` guide and reference texts and the docs contract row.
- [ ] Sync the mirrors, record the owned-skill versions and re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] A Project Config with `wait` or `follow` is read and declared, a project without the key waits, User Config never sets the key and warns once, and any other value is refused by name.
- [ ] The status view reports the value at the default branch head, the remote-tracking ref first, and never the working tree's or the current branch's; a differing checkout value is named on the same line; an unresolved head or an unreadable configuration gives `wait (default)` with its reason and exit 0; text and JSON agree.
- [ ] The guides and both copies of the reference describe the key and the line, and every test that pins the status text passes with the new line.

## Context

- instruction: `docs/adr/0442-an-open-authoring-question-waits-unless-the-project-declares-that-the-session-follows-the-author.md`
- instruction: `docs/adr/0423-an-undeclared-correction-sign-parks-the-item-at-its-merge-as-a-finding.md`
- instruction: `internal/config/require_defect_class_test.go`
- instruction: `internal/cli/correction_signs.go`
- instruction: `internal/config/delivery.go`
- instruction: `internal/preflight/default_branch_head.go`
- interface: `internal/config/config.go`
- interface: `internal/config/authoring.go`
- creates: `internal/config/authoring_questions_test.go`
- interface: `internal/cli/follow_up_status.go`
- creates: `internal/cli/follow_up_status_authoring.go`
- creates: `internal/cli/follow_up_status_authoring_test.go`
- interface: `internal/cli/follow_up_status_queue_test.go`
- interface: `internal/cli/follow_up_status_test.go`
- interface: `internal/cli/live_documents_docscontract_test.go`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/user-guide/commands/follow-up.md`
- interface: `.agents/skills/roundfix/references/follow-up.md`
- interface: `skills/roundfix/references/follow-up.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The resolved configuration and its `authoring.open_questions` value — consumers: `./internal/config`
- The status document of `roundfix follow-up --status`, its text and its JSON field `authoring_questions` — consumers: `./internal/cli`
- The `follow-up` guide, both copies of the `follow-up` reference and the Behavior Surface Record — consumers: `./internal/docscontract`, `./skills`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheAuthoringQuestionPreferenceIsAProjectConfigValue|TestRequireDefectClassIsAProjectConfigBoolean)$' ./internal/config 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheAuthoringQuestionPreferenceIsAProjectConfigValue/builtin_waits TestTheAuthoringQuestionPreferenceIsAProjectConfigValue/absent_key_waits TestTheAuthoringQuestionPreferenceIsAProjectConfigValue/project_config_wait TestTheAuthoringQuestionPreferenceIsAProjectConfigValue/project_config_follow TestTheAuthoringQuestionPreferenceIsAProjectConfigValue/user_config_ignored_once TestTheAuthoringQuestionPreferenceIsAProjectConfigValue/refuses_maybe TestTheAuthoringQuestionPreferenceIsAProjectConfigValue/refuses_a_number TestTheAuthoringQuestionPreferenceIsAProjectConfigValue/refuses_a_list TestTheAuthoringQuestionPreferenceIsAProjectConfigValue/refuses_a_mapping TestTheAuthoringQuestionPreferenceIsAProjectConfigValue/refuses_an_unknown_key TestTheAuthoringQuestionPreferenceIsAProjectConfigValue/default_configuration_omits_the_key TestRequireDefectClassIsAProjectConfigBoolean; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./internal/config` — expected: exit 0. Before this Task the test does not exist, so its pass lines are missing; afterwards each case of the key passes and the whole configuration suite passes.
- `out="$(go test -count=1 -timeout 30m -v -tags docscontract -run '^(TestTheStatusViewReportsTheAuthoringQuestionPreference|TestTheStatusViewShowsTheSpecQueueTable|TestTheStatusViewSurfaceTranscripts|TestFollowUpStatusNamesTheManualClassification|TestFollowUpStatusWritesNothing|TestFollowUpHelpTranscript|TestTheFollowUpGuidesNameTheJSONFieldsAndMessageTypes)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheStatusViewReportsTheAuthoringQuestionPreference/default_without_the_key TestTheStatusViewReportsTheAuthoringQuestionPreference/follow_committed_on_main TestTheStatusViewReportsTheAuthoringQuestionPreference/follow_only_in_the_working_tree TestTheStatusViewReportsTheAuthoringQuestionPreference/follow_only_on_a_feature_branch TestTheStatusViewReportsTheAuthoringQuestionPreference/remote_tracking_ref_wins TestTheStatusViewReportsTheAuthoringQuestionPreference/no_default_branch_head TestTheStatusViewReportsTheAuthoringQuestionPreference/unreadable_project_config TestTheStatusViewShowsTheSpecQueueTable TestTheStatusViewSurfaceTranscripts TestFollowUpStatusNamesTheManualClassification TestFollowUpStatusWritesNothing TestFollowUpHelpTranscript TestTheFollowUpGuidesNameTheJSONFieldsAndMessageTypes; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 40m -run 'FollowUp|Status' ./internal/cli` — expected: exit 0. Before this Task the status test does not exist, so its pass lines are missing; afterwards every case of the reader passes, the tests that pin the status text pass with the new line, the guides name the new JSON fields, and the follow-up and status tests of the command package pass.
- `grep -q 'authoring.open_questions in User Config is ignored; set authoring.open_questions in Project Config' internal/config/authoring.go || { echo 'the User Config warning is missing' >&2; exit 1; }; for phrase in 'A Roundfix older than the release that introduces this key refuses a Project Config that holds it.' 'reads the key at the default branch head, preferring its remote-tracking ref and falling back to its local branch'; do tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- "$phrase" || { printf 'missing phrase in docs/user-guide/configuration.md: %s\n' "$phrase" >&2; exit 1; }; done; for file in docs/user-guide/commands/follow-up.md .agents/skills/roundfix/references/follow-up.md skills/roundfix/references/follow-up.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'The authoring questions report carries these JSON fields:' || { printf 'missing phrase in %s\n' "$file" >&2; exit 1; }; done; make skills-sync-check && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task neither the warning nor the guide texts exist; afterwards the configuration guide and the three `follow-up` documents describe the key and the line, the mirrors match, the skills suite passes with the recorded versions, and the documentation contracts, the current Behavior Surface Record included, pass.

## References

- `_prd.md` → Goals (first and second bullets); User Stories 1 and 9; Core Features 4 and 5; Success Metric 1; Success Metric 2; Success Metric 8
- `_techspec.md` → Preference key; Preference reader; Interfaces; Data Models, "The preference" and "The status fields"; Invariants 1, 2, 3 and 6; API Contract 1; API Contract 2; API Contract 11; API Contract 13; Testing Approach; Build Order 1
- [`references/2026-11-19-prototype-and-governed-paths.md`](references/2026-11-19-prototype-and-governed-paths.md), section 1
- ADR-0442, ADR-0423, ADR-0417, ADR-0392, ADR-0256

## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, tech-writer, evidence-gate, knowledge-workspace, systematic-debugging, no-workarounds.

Implemented only the preference and status-report slice. The pre-existing change was the Daemon's `status: in_progress`; it is preserved. No Task Graph, other Task, repository Project Config, Baseline source, ADR, dependency or verification configuration was changed.

Acceptance evidence:

- Project Config: added `Authoring`, the two preference constants, the strict overlay and the built-in undeclared `wait`. User Config is filtered before value validation with the exact once-per-load warning; malformed Project Config values retain their named diagnostic and parse prefix. Both init defaults omit the key. `go test ./internal/config -run '^TestTheAuthoringQuestionPreferenceIsAProjectConfigValue$' -count=1` exited 0, exercising all eleven required subtests through `Builtin`, `ResolveConfigProposal` and isolated `Load` calls. Before implementation, the focused `builtin_waits` invocation failed to compile on the missing types and field.
- Default-branch status: added the local-object reader using the supplied Git runner, `ResolveDefaultBranchHead` and `DeliveryConfigAtCommit`; added `authoring_questions` immediately after `maintainer` and its text line immediately before `Sources:`. The seven required scenarios assert the entire preference line and all five JSON fields, including exit 0 on unreadable default config, checkout differences and remote-tracking precedence. `go test ./internal/cli -run '^TestTheStatusViewReportsTheAuthoringQuestionPreference$' -count=1` exited 0. The fixtures remove their seeded remote, use only local refs, reject unexpected commands through the injected runner, and fake forge access.
- Documentation and existing transcripts: transcribed API Contracts 11 and 13 with exactly the task_01 exclusions; added only the requested docs-contract row. Updated the three existing pinned status expectations with the new line. `go test ./internal/cli -run '^(TestTheStatusViewShowsTheSpecQueueTable|TestTheStatusViewSurfaceTranscripts|TestFollowUpStatusNamesTheManualClassification|TestFollowUpStatusWritesNothing)$' -count=1` exited 0; the writes-nothing test is unchanged. After the last status-test edit and mirror sync, `go test -tags docscontract ./internal/cli -run '^(TestTheFollowUpGuidesNameTheJSONFieldsAndMessageTypes|TestTheStatusViewReportsTheAuthoringQuestionPreference)$' -count=1` exited 0.

Sanctioned regeneration:

- `make skills-sync` exited 0.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` exited 0; it raised both Roundfix version fields from `0.1.278` to `0.1.279` and recorded the digest.
- `make skills-sync` exited 0 again; byte comparisons confirmed that both modified Roundfix mirrors equal their canonical files.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` exited 0. It is invoked again after this Result as the final repository edit, preserving the required ordering after all user-guide changes. This is focused regeneration evidence, not a Daemon Verification verdict.

A read-only scope comparison found every changed path declared by this Task; `git diff --check` exited 0. The final postflight also includes the generated Behavior Surface Record. No mirror, version field or derived record was hand-edited.

Declared Verification, repository Verification and Task settlement remain with the Daemon. No commit, push, Pull Request or ai-memory message was made. Delegated Decision parsing, fields and rendering remain task_02's work; its future-facing guide text is retained exactly as requirement 8 directs.


### Verification Feedback — attempt 1

Inspected the Daemon diagnostic artifact at `/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261119T203521Z_1b24f6b82481f3c8/verification/batch-001-attempt-1.log` and the related JSON contract test. The failure was a stale assertion in `TestTheStatusJSONCarriesTheQueueAndTheBasis`: it required twelve top-level fields after this Task intentionally added the thirteenth, `authoring_questions`.

- Reproduced with `go test ./internal/cli -run '^TestTheStatusJSONCarriesTheQueueAndTheBasis$' -count=1`, which exited 1 before the repair.
- Updated that test's exact field count to thirteen and extended its field-order assertion to require `authoring_questions` after `maintainer`. Its existing queue, forecast, source and report assertions remain intact; no production behavior changed in this repair.
- `go test ./internal/cli -run '^(TestTheStatusJSONCarriesTheQueueAndTheBasis|TestTheStatusViewReportsTheAuthoringQuestionPreference)$' -count=1` exited 0 after the repair. This covers the original failure and all seven preference scenarios together.
- `git diff --check` exited 0. The sanctioned Behavior Surface recorder is invoked after this evidence update as the final repository edit, retaining requirement 10's ordering. No user-guide or skill source changed during the repair.

The Daemon's declared Verification commands were not rerun, and its Task status remains unchanged. Full Verification and settlement remain with the Daemon.

## Carry-forward provenance

- Source Run: `run_20261119T203521Z_1b24f6b82481f3c8`
- Source commit: `7a715b5ecc8b03a819b4e3b34d9ba42b2d7b5c99`
