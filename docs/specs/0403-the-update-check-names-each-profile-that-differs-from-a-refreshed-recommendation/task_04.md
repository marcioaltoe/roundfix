---
task: task_04
spec: 0403-the-update-check-names-each-profile-that-differs-from-a-refreshed-recommendation
status: completed
type: backend
complexity: high
---

# Task 04: The Update Check names each configured profile that differs from the Recommended Profile

## Overview

`roundfix update check` runs the existing read-only comparison with the
Recommended Profile and reports each configured category that differs: in
text, one `Profiles:` line each, naming a deviation declared against an older
snapshot as pinned to an older snapshot with its re-pin key, and the adopt
command; in JSON, the optional `recommendations` member. With nothing to
report, the bytes stay as before. Fixture-repository tests prove the Surface
Transcript, the JSON member and every case that reports nothing.

## Requirements

1. MUST add `internal/cli/update_check_recommendations.go` with the types
   and functions of `_techspec.md` → Interfaces, built from
   `config.CheckRecommendations(loaded.Config)` and formatted with the
   existing `formatRecommendationProfile` and `profilesCheckProfile`.
2. MUST report a category exactly when its status is `differs`, in the
   order `CheckRecommendations` returns, with `state`
   `pinned_to_older_snapshot`, `deviation` and `repin`
   (`profiles.<category>.deviation.from: <snapshot>`) when it carries a
   deviation, and `state` `differs` without them otherwise; `adopt` is
   `roundfix profiles check --apply --scope <source>` (invariants 3 and 4).
3. MUST return no report, so the JSON member is omitted and no line is
   printed, when `loaded.RepositoryVersionErr` is set, when the comparison
   returns an error, or when no category is reported (invariant 5).
4. MUST add `Recommendations *updateCheckRecommendations` with tag
   `json:"recommendations,omitempty"` to the report after
   `AuthoringModels` and after any member Spec 0402 added, before `Next`;
   set it in `runUpdateCommand`; and print the lines in
   `printUpdateCheckReport` after the `Authoring models:` output and any
   line Spec 0402 added there, before `Next:`, exactly as API Contract 1
   writes them, the reason with whitespace collapsed and limited to 300
   characters. `next` and the exit code MUST NOT change (invariant 6).
5. MUST add to the `update` usage text in `internal/cli/cli.go` the sentence
   of API Contract 3 and update the help text of `TestUpdateCheckUsage` in
   `internal/cli/update_check_test.go` to it. MUST change no other existing
   test expectation.
6. MUST add `internal/cli/update_check_recommendations_test.go` with
   `TestUpdateCheckNamesEachProfileThatDiffers`, calling `t.Parallel()`
   first and in each subtest, with the existing update check fixture helpers
   and a fake release lister, one subtest per case, each asserting the whole
   standard output, the standard error and the exit code, never a
   substring:
   - `transcript`: Surface Transcript 1 of `_techspec.md`, in a committed
     temporary Git repository whose `.roundfixrc.yml` holds the transcript's
     `docs` and `review` profiles; every `Profiles:` line in order, the
     expected text built from `config.ModelRecommendationSnapshotVersion` and
     `config.RecommendedProfile`, and no `Profiles:` line for another
     category;
   - `json`: the same fixture with `--format json`; the `recommendations`
     member equal field by field to API Contract 2, after `authoringModels`
     when present and before `next`, and `next` equal to the run without
     profiles;
   - `pinned_and_current_add_nothing`: a `review` deviation whose `from` is
     `config.ModelRecommendationSnapshotVersion` and a `docs` profile equal
     to the recommendation: no `Profiles:` line and no member;
   - `user_config_source`: the `docs` difference in User Config: source
     `user`, `in User Config` never printed, adopt `--scope user`;
   - `no_profiles_keep_the_report_bytes`: no profile configuration: the
     bytes equal those of the same run with the comparison removed, which
     the test states as the literal text of the existing transcripts;
   - `newer_baseline_reports_no_member`: a repository whose Baseline a newer
     Roundfix wrote, with a differing `docs` profile in Project Config: no
     `Profiles:` line and no member.
   No test reads the clock to decide a snapshot.
7. MUST keep `roundfix profiles check`, its JSON, `--apply`, the Doctor line
   and the upgrade notice byte-identical (invariant 7): no change to
   `profiles_check.go`, `profiles_check_apply.go`, `doctor.go` or
   `upgrade.go`.
7a. MUST re-record `docs/references/behavior-surfaces.json` after the last
   edit of the Update Check's sources, with the `-record-skill-coverage` step
   (`go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`),
   so the `command update check` fingerprint is current and
   `make verify-changed` passes.
8. MUST change no Governed Path and no file outside the ones this Task
   declares.

## Subtasks

- [ ] Add the report types, the builder and the printer.
- [ ] Wire the member and the lines into the Update Check.
- [ ] Add the usage sentence and update the usage test.
- [ ] Add the fixture-repository tests, the transcript included.

## Acceptance Criteria

- [ ] Each configured category that differs is named once, with its adopt command, and an older pin with its re-pin key.
- [ ] A current category, a pin to the shipped snapshot, no configuration and an unread Project Config leave the report's bytes unchanged.
- [ ] The report's `next` and exit code never change, and nothing is written.
- [ ] Surface Transcript 1 passes in one named test.

## Context

- instruction: `docs/adr/0416-the-update-check-names-each-configured-profile-that-differs-from-the-recommended-profile.md`
- instruction: `docs/user-guide/commands/update.md`
- instruction: `internal/config/recommendation_check.go`
- instruction: `internal/cli/profiles_check.go`
- creates: `internal/cli/update_check_recommendations.go`
- creates: `internal/cli/update_check_recommendations_test.go`
- interface: `internal/cli/update_check.go`
- interface: `internal/cli/update_check_test.go`
- interface: `internal/cli/cli.go`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The `roundfix/update-check/v1` report and the `update` usage text — consumers: `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestUpdateCheckNamesEachProfileThatDiffers|TestUpdateCheckUsage|TestUpdateCheckTranscripts)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestUpdateCheckNamesEachProfileThatDiffers/transcript TestUpdateCheckNamesEachProfileThatDiffers/json TestUpdateCheckNamesEachProfileThatDiffers/pinned_and_current_add_nothing TestUpdateCheckNamesEachProfileThatDiffers/user_config_source TestUpdateCheckNamesEachProfileThatDiffers/no_profiles_keep_the_report_bytes TestUpdateCheckNamesEachProfileThatDiffers/newer_baseline_reports_no_member TestUpdateCheckUsage TestUpdateCheckTranscripts; do printf '%s\n' "$out" | grep -qF -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; grep -qF -- 'names each configured Agent Selection Profile that differs' internal/cli/cli.go && go test -count=1 -timeout 30m ./internal/cli` — expected: exit 0. Before this Task the new test and its subtests do not exist and the usage text lacks the sentence, so the command fails; afterwards each case passes by name and the whole command package passes.

## References

- `_prd.md` → Goals (third bullet); Core Feature 3; Success Metrics 3 and 4; Non-Goals
- `_techspec.md` → Update Check report; Interfaces (invariants 3 to 7); API Contracts 1, 2 and 3; Surface Transcript 1; Testing Approach (Task 4); Build Order 4
- ADR-0416, ADR-0181, ADR-0371


## Result

Skills loaded: implement-task, agentic-cli-design, coding-guidelines, golang-cli, golang-error-handling, golang-safety, golang-structs-interfaces, golang-testing, testing-boss, tech-writer, systematic-debugging, evidence-gate.

Added the optional Update Check profile comparison, preserving the comparison's
category order and reporting only `differs` rows. Older deviations receive
`pinned_to_older_snapshot`, their deviation, a re-pin key and the source-specific
adopt command. Text uses the existing profile formatter and collapses and bounds
reasons to 300 characters; JSON keeps the original deviation reason. The member
and lines follow authoring models and precede `next`/`Next:`. The comparison is
omitted for an unread Project Config, a comparison error or no differing rows.
The usage sentence and only its existing help expectation were updated.

Focused evidence:

- Before implementation, `go test -count=1 -timeout 5m ./internal/cli -run '^TestUpdateCheckNamesEachProfileThatDiffers$/transcript$'` exited 1 at the exact stdout assertion: both required `Profiles:` lines were absent, with exit 0 and empty stderr.
- `go test -count=1 -timeout 5m -v ./internal/cli -run '^(TestUpdateCheckNamesEachProfileThatDiffers|TestUpdateCheckOmitsAnUnavailableProfileComparison)$'` exited 0. All ten command subtests and both unavailable-comparison subtests passed by name.
- After the final fixture edits, `go test -count=1 -timeout 5m ./internal/cli -run '^(TestUpdateCheck|TestTheUpdateCheckNamesARepository)'` exited 0. This focused regression selection includes existing transcripts, help, release-cache behavior, authoring output and newer-Baseline reporting.
- Three temporary Go overlays changed only the compiled builder, leaving repository source bytes untouched. Reporting current pins failed `pinned_and_current_add_nothing`; dropping deviation classification failed `transcript`; retaining an empty report failed `no_profiles_keep_the_report_bytes`. Each probe exited 1 at its intended full-output assertion, rather than a build or environment failure.
- After the final Update Check edits, `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` exited 0. Its generated diff changes only the `command update check` fingerprint.

Acceptance evidence:

- Each difference and older pin: `transcript` asserts the entire stdout, stderr and exit code for the committed Project Config fixture, with docs before review and no other category. `json` asserts every serialized field and its order, including placement after a present `authoringModels` member. Expectations read the shipped snapshot and recommended profiles from config constants/functions.
- Quiet cases preserve bytes: `pinned_and_current_add_nothing`, `no_profiles_keep_the_report_bytes` and `newer_baseline_reports_no_member` assert the entire text and JSON reports, including omission of `recommendations`. The newer-Baseline fixture also carries a differing User Config profile to exercise the refusal guard. The unavailable-comparison test covers the explicit error branch.
- `next`, exit and read-only behavior: every command case asserts exit 0 and its full report, `json` compares `next` with a run before profiles are configured, and `release_update_preserves_next` retains `apply` and its final text line. File/mode snapshots of the committed repository and isolated home remain unchanged across text, JSON and cached runs after priming the existing release cache. User Config scope and older User Config pins are covered separately; whitespace and a 301-character Unicode reason are covered too.
- Surface Transcript 1: `TestUpdateCheckNamesEachProfileThatDiffers/transcript` passed with the complete controlled release and Baseline prefix and both required profile lines.

The diff stays within this Task's declared paths. The existing profiles-check,
apply, Doctor and upgrade implementations are unchanged. No follow-up work was
identified. Task status remains Daemon-owned; declared Task Verification and
repository Verification were not run in this Agent turn.

## Carry-forward provenance

- Source Run: `run_20261118T033001Z_a32a9a0e9d248cbd`
- Source commit: `61d5d3b607d3ba7d6fb180833d8f42b9758c509f`
