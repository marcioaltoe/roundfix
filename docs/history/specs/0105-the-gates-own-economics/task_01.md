---
status: completed
type: backend
---

# Task: Roundfix owns the QA Task's Verification

The gate Task's Verification is authored by hand over a verdict the Daemon
already derived. In one measured case that produced a gate which passed itself
having failed — a hand-written predicate accepted a verdict outside the domain,
and it read as ordinary Task authoring to every reviewer.

## Work

- For a Task of type `qa`, Roundfix supplies the Verification instead of reading
  the author's. The derived command requires the newest QA Report to record a
  passing verdict, so a predicate can no longer accept a verdict outside the
  domain or select an older report.
- Render the derived command into the Task file rather than leaving it implicit.
  Removing the author's control must not also remove the reader's view of what
  will run.
- A `qa` Task that carries its own authored command is refused by name through a
  Spec check finding. Refuse rather than overwrite: silently replacing an
  author's text is how a contract becomes invisible.
- Change nothing for any other Task Type. This is bounded to the one node
  ADR-0091 makes terminal.

## References

- `_prd.md` → Goal 3, User Story 4, Core Feature 3
- `_techspec.md` → Build Order 1; Interfaces: `DerivedQAVerification`,
  `AuthoredQAVerification`
- ADR-0091 makes the gate one terminal Task node of type `qa`

## Verification
- `grep -q "DerivedQAVerification" internal/spec/task.go || exit 1; grep -q "AuthoredQAVerification" internal/spec/task.go || exit 1; grep -q "TestDerivedQAVerification" internal/spec/task_test.go || exit 1; go test -count=1 ./internal/spec ./internal/speccheck`

## Result

Roundfix now derives the effective Verification for `qa` Tasks from the Spec
slug. The rendered command orders QA Reports by valid report date and numeric
rerun sequence, reads only the newest report's frontmatter, and accepts exactly
one `verdict: pass`. Loading or reloading a `qa` Task no longer gives an authored
command control over the effective Verification.

The loader retains whether the Task file's rendered command differs from the
derived command. Spec check reports that mismatch as
`SC-QA-VERIFICATION-AUTHORED`, naming the offending Task instead of overwriting
its visible text. A non-`qa` Task keeps its authored Verification unchanged.

Acceptance evidence:

- QA ownership and newest-report semantics: `TestDerivedQAVerificationRequiresTheNewestReportToPass` executes the derived shell command against passing, failed, partial, missing, duplicate-verdict, body-spoofed, invalid-date, `-00`, `-02`, and `-10` report sets.
- Visible derived contract: `TestReloadTaskDerivesOnlyQAVerification/rendered_derived_qa_command_is_accepted` loads the canonical command from the Task file and preserves it as the effective Verification.
- Authored-command refusal: `TestAuthoredQAVerificationIsRefusedByTaskName` runs the Tasks-stage Spec check and observes one `SC-QA-VERIFICATION-AUTHORED` finding whose summary names `task_01`.
- Non-QA boundary: `TestReloadTaskDerivesOnlyQAVerification/non-qa_command_remains_authored` proves a `backend` Task's command is unchanged and is not classified as authored QA Verification.

Focused checks:

- Before implementation, `env GOCACHE=/private/tmp/roundfix-task-0105-01-gocache go test ./internal/spec -run 'TestDerivedQAVerification|TestReloadTaskDerivesOnlyQAVerification' -count=1` failed to compile because `DerivedQAVerification` and `AuthoredQAVerification` did not exist.
- `env GOCACHE=/private/tmp/roundfix-task-0105-01-gocache go test ./internal/spec -run 'Test(DerivedQAVerificationRequiresTheNewestReportToPass|ReloadTaskPicksUpAgentEdits|ReloadTaskDerivesOnlyQAVerification|ReloadTaskNormalizesStatusValues|ReloadTaskReportsBrokenAgentEdits|LoadQAGateContract|LoadQADeclinedContract|LoadRejectsInvalidQAGateShape|LoadInvalidatesSettledQAGateAfterTaskAppend|QAGateLegacyArchivedManifestsLoadUnchanged)$' -count=1` — passed.
- `env GOCACHE=/private/tmp/roundfix-task-0105-01-gocache go test ./internal/speccheck -run 'Test(AuthoredQAVerificationIsRefusedByTaskName|StageScopeRunsOnlyDetectorsTheStageCanDecide|StageScopeNamesTheDetectorsItSkipped|VerifyNonHermeticRegistersAtTasksStage)$' -count=1` — passed.
- After the final source simplification, `env GOCACHE=/private/tmp/roundfix-task-0105-01-gocache go test ./internal/spec -run 'Test(DerivedQAVerificationRequiresTheNewestReportToPass|ReloadTaskDerivesOnlyQAVerification|LoadQAGateContract)$' -count=1` and the corresponding `internal/speccheck` focused run for `TestAuthoredQAVerificationIsRefusedByTaskName` and `TestStageScopeNamesTheDetectorsItSkipped` both passed.
- `git diff --check` — passed.

The daemon-owned command under `## Verification` was not run.

## Carry-forward provenance

- Source Run: `run_20260831T171256Z_df3674a688059467`
- Source commit: `9020e4a48339b43161f0c041e70fb74b22eebb49`
