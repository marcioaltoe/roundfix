---
task: task_02
spec: 0397-each-phase-of-a-spec-is-recorded-and-starts-in-a-fresh-session
status: completed
type: backend
complexity: high
---

# Task 02: Each prompt of a Run and of a pre-PR review is recorded, and the Run Retention Sweep removes old day files

## Overview

The Daemon appends a prompt entry to the Phase Record where it already
records a prompt's token usage, `roundfix review` appends one after each
prompt of a pre-PR review, and the Run Retention Sweep removes the day files
older than its cutoff and names the count in the GC report. The slice is
whole on its own: a fixture Run with a fake runner leaves one line per prompt
with its Run, scope, working directory, runtime, session and both times; a
record that cannot be written costs one warning and changes no outcome; and a
dry run of the GC Command counts the old day files without removing one. The
GC and Run Database guides, the usage guide's Phase Record section and the
Roundfix skill's `storage` and `phase` references describe it.

## Requirements

1. MUST add to `internal/phaserecord` the prompt half of `_techspec.md` →
   Interfaces and Data Models: `PromptEntry`, `RecordPrompt` and
   `RemoveBefore`, with those signatures. `RecordPrompt` appends one line
   with the keys of a prompt entry in the order Data Models gives, under the
   rules task 01 gave a phase entry: the same directory, day file, modes,
   single write, time format, empty string for an absent value and 1,024-byte
   cut. `RemoveBefore` removes, or with `dryRun` only counts, each file named
   `YYYY-MM-DD.jsonl` whose date is before the cutoff's UTC date, leaves a
   file whose name is not a date alone, and returns `0` and no error when the
   directory does not exist.
2. MUST give the Daemon the collaborator of `_techspec.md` → Interfaces
   (`PromptRecorder`, and `PhaseRecord` in `Dependencies`, where nil records
   nothing) and write the entry of `_techspec.md` → API Contract 4 in
   `recordPromptUsage`, which receives the time the prompt began: one entry
   for every call that records, or would record, a usage row, a prompt that
   failed or was stopped after it returned included, whether or not the usage
   row itself could be written. A prompt that never returns leaves none. A
   failed write prints
   `roundfix: warning: phase record not written for <scope_kind> <scope_id>: <error>`
   on the Run's progress output, once for that prompt, and changes no result,
   error, Task status, Run state, usage row or event.
3. MUST wire the recorder in every place where the CLI builds the Daemon's
   dependencies, from the command environment's home directory: three places
   on the planning branch, two in `internal/cli/implement.go` and one in
   `internal/cli/cli.go`, and any other the default branch holds when the
   Task runs.
4. MUST make `roundfix review` write the entry of `_techspec.md` → API
   Contract 5 after each prompt its review Agent Session answered, every
   segment of a segmented review included, through the same store and the
   command environment's home directory. A prompt that failed before it was
   sent leaves none, and the sealed convention check, which opens no review
   session and reports no session identifier, records nothing. A failed
   write prints
   `roundfix: warning: phase record not written for pre-pr-review <scope_id>: <error>`
   on standard error and changes neither the review's record, its output nor
   its exit code.
5. MUST apply `_techspec.md` → API Contract 6. `sweepRunDatabase` asks the
   store, after its candidate loop, to remove the day files before its
   cutoff, or to count them in a dry run, and its report carries the count.
   `printRunRetentionReport` prints, right after the line that counts the
   Runs and indented as it is, `Phase Record day files removed: <n>`, or
   `Phase Record day files removable: <n>` in a dry run. The GC Command
   MUST count or remove the day files with the same cutoff when no Run
   Database exists as well: `runGC` calls the sweep only when the database
   exists, and Surface Transcript 8 has none. A failure to remove a day file
   is returned as the sweep returns its other failures, after the Runs were
   handled, so it removes no Run less. The `gc` usage text gains, after the
   sentence that ends "and compacts the Run Database.", the sentence "It
   also removes Phase Record day files older than the Run Retention
   cutoff." The line the Run start prints for a sweep does not change.
6. MUST add these tests, each top-level test calling `t.Parallel()` first,
   with fake runners, temporary homes and fixed clocks through the packages'
   existing seams, without editing `internal/cli/cli_test.go` and without
   starting an Agent harness:
   - `internal/phaserecord`: `TestAPromptEntryIsOneLineOfItsDeclaredKeys`,
     with the subtests `keys_in_order`, `absent_values_are_empty_strings`
     and `oversized_field`; and
     `TestRemoveBeforeRemovesOnlyDayFilesBeforeTheCutoff`, with the subtests
     `live`, `dry_run`, `cutoff_date_is_kept`, `not_a_date` and
     `missing_directory`.
   - `internal/daemon`: `TestEachPromptLeavesOnePromptEntryWithItsScope`,
     with one subtest per scope kind (`task`, `qa`, `review`, `session`),
     each asserting the whole entry the recorder received: the Run, the
     scope, the working directory, the runtime identifier, the session
     identifier the fake runner reported and a start no later than the end;
     `TestAFailedPhaseRecordWriteLeavesTheOutcomeUnchanged`, whose recorder
     returns an error and which asserts the unchanged result, usage rows and
     events and the warning line, word for word, once per prompt; and
     `TestANilPromptRecorderRecordsNothing`.
   - `internal/cli`: `TestEveryEngineTheCLIBuildsRecordsItsPrompts`, which
     parses the package's non-test sources and fails for a
     `daemon.Dependencies` literal that sets no `PhaseRecord`;
     `TestAnImplementRunLeavesOnePromptEntryPerPrompt`, which drives a Run
     through the existing fake-runner harness with a temporary home and
     reads the day file; `TestTheReviewRecordsEachPromptInThePhaseRecord`,
     with the subtests `one_prompt` (an empty `run_id`, the scope kind
     `pre-pr-review`, the review session's name, the reviewed worktree as
     `repository` and `work_dir`, the runtime and the reported session),
     `one_spec_named`, `several_specs_named` (an empty `spec`) and
     `unwritable_record` (a regular file where the record's directory
     belongs: the same record, output and exit code, and the warning once);
     `TestTheRunRetentionSweepRemovesPhaseRecordDayFiles`, through
     `sweepRunDatabase`, with the subtests `live`, `dry_run` and
     `not_a_date`; and `TestGCCountsPhaseRecordDayFilesWithoutARunDatabase`,
     which asserts Surface Transcript 8 with its fixture: an empty standard
     error, exit `0`, the line `  Phase Record day files removable: 2`
     right after the line `  Runs removable: 0`, and the three day files
     still present.
   - `skills`: `TestTheRoundfixSkillNamesThePromptEntriesAndTheirRetention`,
     which reads both copies through `roundfixLayoutSources` and requires
     each phrase requirement 8 lists, exactly once, whitespace-normalized,
     in one named subtest each.
   - Every test that pins the Run Retention section or the `gc` usage on
     the default branch when the Task runs moves to the new text. No
     recorded test is removed or renamed, and every test file compiles for
     Windows.
7. MUST document the behavior in the guides:
   - `docs/user-guide/commands/gc.md`, under `## Run Retention`: the
     paragraph "Run Retention also removes each Phase Record day file whose
     date is before the cutoff's UTC date, with or without a Run Database,
     and leaves a file whose name is not a date alone.", and the two report
     lines in a `text` block, `Phase Record day files removed: <n>` and
     `Phase Record day files removable: <n>`, saying that the second is the
     dry run's and that the line follows the one that counts the Runs;
   - `docs/user-guide/run-database-lifecycle.md`, after the paragraph that
     begins "Run Retention never touches": the sentence "Run Retention also
     removes each Phase Record day file dated before its cutoff; the record
     is files under Roundfix Home and no table of the Run Database.";
   - `docs/user-guide/usage.md`, in the `## Phase Record` section task 01
     wrote: that the Daemon writes a prompt entry when a prompt of an Agent
     Session ends and `roundfix review` after each prompt of a pre-PR
     review, the sentence "A prompt that never returns leaves no entry.",
     the sentence "A failed write of a prompt entry prints one warning and
     changes no Task, Run, review or exit code.", and the sentence "Run
     Retention removes each day file dated before its cutoff, and a file
     whose name is not a date is left alone."
8. MUST add to the canonical Roundfix skill, word for word:
   - at the end of `.agents/skills/roundfix/references/phase.md`, the
     paragraph

     > The Daemon appends a prompt entry when a prompt of an Agent Session ends, and `roundfix review` after each prompt of a pre-PR review: the Run, the Spec, the scope (`task`, `qa`, `review`, `session` or `pre-pr-review`), the repository, the working directory, the runtime, the session the adapter reported and the times the prompt began and ended. A prompt that never returns leaves no entry. When the record cannot be written, one warning is printed and the Run or the review ends as it would have.

     pinned as `the_writers` ("The Daemon appends a prompt entry when a
     prompt of an Agent Session ends, and `roundfix review` after each
     prompt of a pre-PR review"), `a_prompt_that_never_returns` ("A prompt
     that never returns leaves no entry.") and `a_failed_write_warns` ("one
     warning is printed and the Run or the review ends as it would have");
   - in `.agents/skills/roundfix/references/storage.md`, as a paragraph of
     its own right after the paragraph that begins
     "`store.run_retention_days` is User Config only",

     > Run Retention also removes each Phase Record day file dated before its cutoff, in the same sweep and with or without a Run Database. The GC report's Run Retention section names the count as `Phase Record day files removed: <n>`, or `Phase Record day files removable: <n>` in a dry run.

     pinned as `retention_removes_day_files` ("Run Retention also removes
     each Phase Record day file dated before its cutoff") and
     `the_report_names_the_count` ("`Phase Record day files removed: <n>`,
     or `Phase Record day files removable: <n>` in a dry run").

   Every other sentence of the two references stays. MUST then run
   `make skills-sync` and the version record
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   which raises the version of `roundfix` from the tree the Task runs on;
   never write a number or a mirror by hand. MUST NOT add or change text
   inside any `### QA settlement` section.
9. MUST record in `docs/references/skill-coverage.json` the Coverage Review
   of the three surfaces `_prd.md` → Skills declares unchanged, `command
   review`, `command review dispose` and `command implement`: each entry's
   `review` becomes `<the day the Task runs, YYYY-MM-DD> — Spec 0397
   task_02: <one sentence saying why the surface needs no skill text>`,
   followed by one space and any review text the entry already holds. MUST
   then re-record the Behavior Surface Record as the Task's last step, after
   the last guide edit, with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
10. Of the Governed Paths the Task may change only
    `.agents/skills/roundfix/SKILL.md`,
    `.agents/skills/roundfix/references/phase.md`,
    `.agents/skills/roundfix/references/storage.md` and
    `skills/roundfix/SKILL.md`, the first and the last only through the two
    commands of requirement 8. MUST NOT change the `phase start` command,
    `CONTEXT.md`, any guide of the `review` or `implement` commands, any
    other skill, any Baseline file, `.roundfixrc.yml`,
    `docs/references/coverage-record.json`, the Makefile, the CI workflows,
    `go.mod`, `go.sum` or the lint configuration. MUST NOT run a `roundfix`
    found on `PATH` for `skills install` or `baseline update`.

## Subtasks

- [ ] Add the prompt entry and the removal by date to the store, with their tests.
- [ ] Record each prompt in the Daemon and wire the recorder wherever the CLI builds an engine.
- [ ] Record each prompt of `roundfix review`.
- [ ] Remove or count the day files in the Run Retention Sweep and in the GC Command without a Run Database, and print the count.
- [ ] Write the three guides and the two skill references, sync the mirror and record the skill version.
- [ ] Record the three Coverage Reviews and re-record the Behavior Surface Record after the last guide edit.

## Acceptance Criteria

- [ ] A fixture Run with a fake runner leaves one prompt entry for each prompt that returned, for a Task, the QA gate, a review Batch and a session, with its Run, scope, working directory, runtime, reported session identifier and both times; a `roundfix review` with a fake runner leaves one for each review prompt with no Run.
- [ ] With a recorder that fails, or a record that cannot be written, the same Run and the same review end with their usual outcome, events, record and exit code and print the warning once per prompt; a nil recorder writes nothing.
- [ ] Every engine the CLI builds holds the recorder.
- [ ] The Run Retention Sweep removes the day files dated before its cutoff, keeps the others and a file whose name is not a date, and names the count; a dry run removes none; the GC Command does the same without a Run Database.
- [ ] The GC and Run Database guides, the usage guide, both copies of the `storage` and `phase` references, the Skill Coverage Map and the Behavior Surface Record describe the behavior.

## Context

- instruction: `docs/adr/0407-a-phase-record-names-the-spec-and-phase-of-supervisor-work-and-a-phase-boundary-ends-the-session.md`
- instruction: `docs/adr/0255-run-retention-removes-terminal-runs-whole-and-compacts-incrementally.md`
- instruction: `internal/daemon/token_usage_test.go`
- instruction: `internal/agent/agent.go`
- creates: `internal/phaserecord/record.go`
- creates: `internal/phaserecord/prompt.go`
- creates: `internal/phaserecord/prompt_test.go`
- interface: `internal/daemon/engine.go`
- interface: `internal/daemon/agent_session_owner.go`
- creates: `internal/daemon/phase_record_test.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/implement.go`
- interface: `internal/cli/review.go`
- interface: `internal/cli/run_retention.go`
- interface: `internal/cli/gc.go`
- interface: `internal/cli/gc_run_retention_test.go`
- creates: `internal/cli/phase_record_test.go`
- interface: `docs/user-guide/commands/gc.md`
- interface: `docs/user-guide/run-database-lifecycle.md`
- interface: `docs/user-guide/usage.md`
- creates: `.agents/skills/roundfix/references/phase.md`
- creates: `skills/roundfix/references/phase.md`
- interface: `.agents/skills/roundfix/references/storage.md`
- interface: `skills/roundfix/references/storage.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- creates: `skills/phase_record_writers_test.go`
- interface: `docs/references/skill-coverage.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The Daemon's dependencies and the prompt recorder — consumers: `./internal/daemon`, `./internal/cli`
- The Run Retention section of the GC report — consumers: `./internal/cli`
- The Roundfix skill's references and the owned-skill version record — consumers: `./skills`
- The GC, Run Database and usage guides, the Skill Coverage Map and the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 15m -v -run '^(TestAPromptEntryIsOneLineOfItsDeclaredKeys|TestRemoveBeforeRemovesOnlyDayFilesBeforeTheCutoff)$' ./internal/phaserecord 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestAPromptEntryIsOneLineOfItsDeclaredKeys/keys_in_order TestAPromptEntryIsOneLineOfItsDeclaredKeys/absent_values_are_empty_strings TestRemoveBeforeRemovesOnlyDayFilesBeforeTheCutoff/live TestRemoveBeforeRemovesOnlyDayFilesBeforeTheCutoff/dry_run TestRemoveBeforeRemovesOnlyDayFilesBeforeTheCutoff/cutoff_date_is_kept TestRemoveBeforeRemovesOnlyDayFilesBeforeTheCutoff/not_a_date; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; out="$(go test -count=1 -timeout 20m -v -run '^(TestEachPromptLeavesOnePromptEntryWithItsScope|TestAFailedPhaseRecordWriteLeavesTheOutcomeUnchanged|TestANilPromptRecorderRecordsNothing|TestEachPromptRecordsItsUsageWithTheOwnerScope|TestAFailedUsageWriteLeavesTheOutcomeUnchanged)$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestEachPromptLeavesOnePromptEntryWithItsScope/task TestEachPromptLeavesOnePromptEntryWithItsScope/qa TestEachPromptLeavesOnePromptEntryWithItsScope/review TestEachPromptLeavesOnePromptEntryWithItsScope/session TestAFailedPhaseRecordWriteLeavesTheOutcomeUnchanged TestANilPromptRecorderRecordsNothing; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 15m ./internal/phaserecord && GOOS=windows go build -buildvcs=false ./internal/phaserecord` — expected: exit 0. Before this Task the prompt entry, the removal by date and the Daemon's recorder do not exist, so the pass lines are missing; afterwards a prompt entry is one line of its declared keys, only day files before the cutoff are removed, every scope leaves its entry, and a failed write or a nil recorder changes nothing.

- `out="$(go test -count=1 -timeout 40m -v -run '^(TestEveryEngineTheCLIBuildsRecordsItsPrompts|TestAnImplementRunLeavesOnePromptEntryPerPrompt|TestTheReviewRecordsEachPromptInThePhaseRecord|TestTheRunRetentionSweepRemovesPhaseRecordDayFiles|TestGCCountsPhaseRecordDayFilesWithoutARunDatabase)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestEveryEngineTheCLIBuildsRecordsItsPrompts TestAnImplementRunLeavesOnePromptEntryPerPrompt TestTheReviewRecordsEachPromptInThePhaseRecord/one_prompt TestTheReviewRecordsEachPromptInThePhaseRecord/several_specs_named TestTheReviewRecordsEachPromptInThePhaseRecord/unwritable_record TestTheRunRetentionSweepRemovesPhaseRecordDayFiles/live TestTheRunRetentionSweepRemovesPhaseRecordDayFiles/dry_run TestTheRunRetentionSweepRemovesPhaseRecordDayFiles/not_a_date TestGCCountsPhaseRecordDayFilesWithoutARunDatabase; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 40m -run 'Retention' ./internal/cli` — expected: exit 0. Before this Task none of the five tests exists, so the pass lines are missing; afterwards every engine the CLI builds records its prompts, a Run and a review leave one entry per prompt, an unwritable record changes no outcome, Surface Transcript 8 holds, and the CLI's retention tests pass.

- `tr -s '[:space:]' ' ' < docs/user-guide/commands/gc.md | grep -qF -- 'Phase Record day files removable: <n>' || { echo 'the gc guide lacks the day file count' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/run-database-lifecycle.md | grep -qF -- 'Run Retention also removes each Phase Record day file dated before its cutoff' || { echo 'the Run Database guide lacks the day files' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/usage.md | grep -qF -- 'A failed write of a prompt entry prints one warning and changes no Task, Run, review or exit code.' || { echo 'the usage guide lacks the prompt entry writers' >&2; exit 1; }; test "$(grep -c -- 'Spec 0397 task_02' docs/references/skill-coverage.json)" -eq 3 || { echo 'the three Coverage Reviews are not recorded' >&2; exit 1; }; diff -r .agents/skills/roundfix skills/roundfix >/dev/null || { echo 'the Roundfix skill differs from its mirror' >&2; exit 1; }; out="$(go test -count=1 -timeout 20m -v -run '^(TestTheRoundfixSkillNamesThePromptEntriesAndTheirRetention|TestThePhaseReferenceDocumentsTheCommandAndTheRecord|TestEveryOwnedSkillVersionIsRecorded)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheRoundfixSkillNamesThePromptEntriesAndTheirRetention TestThePhaseReferenceDocumentsTheCommandAndTheRecord TestEveryOwnedSkillVersionIsRecorded; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the guides and the references do not name the day files or the writers of a prompt entry; afterwards they do in both copies, the three Coverage Reviews and the skill version are recorded, and the skills suite and the documentation contracts pass.

## References

- `_prd.md` → User Story 3; Core Features 4, 5, 6, 7 and 13; Success Metrics 4 and 5; Skills
- `_techspec.md` → Current behavior; Interfaces; Data Models; API Contract 4; API Contract 5; API Contract 6; Surface Transcript 8; Invariants 1, 2 and 3; Declared breaks; Testing Approach; Build Order 2
- ADR-0407, ADR-0255


## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-context, golang-error-handling, golang-safety, golang-structs-interfaces, golang-testing, testing-boss, tech-writer, systematic-debugging, evidence-gate.

The initial worktree held only the Daemon's `pending` → `in_progress` change in this Task file. The prompt store, recorder collaborator and retention count were absent. Task status, the Task Graph and other Task files were left untouched by this implementation.

### Implementation and acceptance evidence

1. Added `PromptEntry`, `RecordPrompt` and `RemoveBefore` to `internal/phaserecord`. Prompt entries keep the declared key order, UTC seconds, present empty string fields, UTF-8-safe 1,024-byte cuts, private creation modes and one appending write. Daemon entries use the recorded Run repository, the request's working directory/runtime, the owning scope and the adapter-reported session. The scope tests assert each whole entry for `task`, `qa`, `review` and `session`, through both runner dispatch paths. The CLI fixture leaves two entries for its two returned prompts. Review fixtures cover no Spec, one Spec, several Specs, segmented review and the plain runner; pre-prompt failures and the sealed convention check leave no extra entry.
2. Recording runs independently of usage support or a rejected usage write, including a failed or stopped prompt and a canceled caller. The warning tests assert the exact warning once per prompt, unchanged returned results/errors, usage rows/events and Run state; the Task-cycle fixture retains its normal settlement. The review fixture compares the normal and unwritable-record results and persisted/stdout records, normalizing only the newly generated session name, and asserts the exact stderr warning. A nil recorder preserves usage and emits no warning. Usage retains the prompt-end timestamp rather than taking another clock reading after the advisory write.
3. Wired all three production `daemon.Dependencies` literals with a recorder constructed from the command environment's home. The AST test scans every non-test CLI source and rejects a missing or explicit nil `PhaseRecord`.
4. Retention counts/removes only date-named day files before the cutoff's UTC date, after the Run candidate loop. The count follows the Runs count in the report. The exact no-database dry-run transcript passes, keeps all three fixture files and creates no database; the same fixture also checks live removal. A record-directory failure is returned after eligible Runs have been removed, while protected Runs remain. Updated the existing retention assertions and Worktree Build Cache transcript for the new report line; the Run-start sweep line is unchanged.
5. Added the specified GC, Run Database and usage guide text and the exact canonical `phase`/`storage` paragraphs. Inspection confirmed that every earlier reference sentence remains, both skill trees are byte-identical and exactly three Coverage Reviews carry this Task's prefix. `make skills-sync` and the sanctioned skill-version recorder exited 0; both Roundfix version fields in both copies moved from `0.1.241` to `0.1.242`. The final Behavior Surface recorder exited 0 after the source and guide edits; its generated changes cover GC and the three changed guides.

### Focused checks

- `go test ./internal/daemon ./internal/cli ./internal/phaserecord ./skills -run 'Test(EachPrompt|AFailedPhaseRecord|AFailedUsageWrite|ANilPrompt|APromptIsRecordedWithout|APromptEntry|RemoveBefore|EveryEngineTheCLI|AnImplementRunLeaves|TheReviewRecords|TheSealedConventionCheckLeaves|TheRunRetentionSweepRemovesPhase|GCCountsPhase|APhaseRecordRetentionFailure|GCDryRunReportsRunRetention|GCRemovesRunsPastRunRetention|RunDatabaseSweepStopsAtBudgetBoundaries|GCDryRunListsWorktreeBuildCaches|GCRemovesWorktreeBuildCaches|TheRoundfixSkillNamesThePrompt)' -count=1` — exit 0 in all four packages.
- `go test ./internal/cli -run 'Test(TheReviewRecordsEachPromptInThePhaseRecord|TheSealedConventionCheckLeavesNoPromptEntry|ReviewRetries|ReviewNeverRetriesAfterSessionPromptWasSent|ReviewActivatesTheFallback)' -count=1` — exit 0 after the final review-writer change. A preceding targeted regression reproduced an unwanted entry for a plain runner's pre-prompt preparation failure; the recorder now uses the adapter's broader explicit `PromptSent` signal while the retry policy remains unchanged.
- The focused usage/GC/help/skill-reference check also exercised `TestRunGCHelp`, `TestThePhaseReferenceDocumentsTheCommandAndTheRecord` and the existing usage-warning checks — exit 0.
- `GOOS=windows go test -c -o /private/tmp/0397-phaserecord-windows.test.exe ./internal/phaserecord`, and equivalent commands for `daemon`, `cli` and `skills` with their matching output basenames — exit 0. CLI compilation was repeated after the final review change. No Windows binary was executed.
- `make skills-sync` — exit 0.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` — exit 0.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` — exit 0; used as the final regeneration after all source/guide edits and the Result update.
- `git diff --check` and exact-text/mirror/Coverage Review inspection — clean.

The declared Verification commands and repository Verification remain for the Daemon. No Agent harness was started; no commit, push or pull request was made.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/cli/worktree_build_cache_test.go`
