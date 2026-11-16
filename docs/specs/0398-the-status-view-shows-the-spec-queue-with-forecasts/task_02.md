---
task: task_02
spec: 0398-the-status-view-shows-the-spec-queue-with-forecasts
status: pending
type: backend
complexity: high
---

# Task 02: The delivery history, the phase entries and the queue rows are read for the status view

## Overview

Three readers give the status view what the **Spec Queue Table** rests on,
and none of them writes. The Phase Record package of Spec 0397 gains a
reader of phase entries. The command layer gains the delivery history, read
from the Runs it already lists and from one `git log`, and the queue rows:
Delivery Queue items, planned Specs, Specs being authored and open Backlog
Entries, each with its status, its note and the wait the command can name.
The slice is whole on its own: its tests call the readers with temporary
repositories and homes and assert each row, each note and each bound. The
command prints exactly what it printed; task 03 calls the readers and adds
the table.

## Requirements

1. MUST add to the package `internal/phaserecord`, in a new file, the
   reader of `_techspec.md` → Interfaces: the type `RecordedPhase` and the
   function `Phases(homeDir, since)`, with those fields and that signature.
   - It reads the day files of the directory `phase-record` under
     `<homeDir>/.roundfix`, named `YYYY-MM-DD.jsonl`, whose date is not
     before the UTC date of `since`, and returns the entries whose kind is
     `phase`, whose schema is the package's `Schema` and whose `time` is at
     or after `since`, the oldest first.
   - It skips a prompt entry without counting it. It skips and counts a
     line that is not a JSON object, a line of another schema name, and a
     phase entry whose time or phase does not parse.
   - A missing directory is an empty record with no error. An entry of the
     directory that is not a regular file with a date for its name is
     ignored. A directory that cannot be listed or a day file that cannot be
     read is an error.
   - It creates, changes and removes nothing, not even the directory.
   - It uses the package's own schema name, phase type and phase parser, so
     the format keeps one owner; it changes no existing declaration of the
     package and no existing test.
2. MUST add the delivery history of `_techspec.md` → API Contract 4 as
   `readDeliveryHistory`, in a file of its own beside the command's other
   files. Its input carries the Runs the view already listed, the checkout,
   the Spec Root and the archive root, the Run Retention days, the instant
   `now` and the Git runner, a `preflight.GitRunner`; it reads no clock and
   no dependency of the command. The Git read, the default branch head and
   the `git log`, stays apart from the reading of the Runs inside the file,
   because task 03 runs it beside the forge call.
   - A Spec's start is the earliest `CreatedAt` among the Runs of kind
     `implement` with its slug. A Run of another kind never gives a start.
   - Its plan and its end come from one `git log` on the default branch
     head that `preflight.ResolveDefaultBranchHead` names, with exactly the
     arguments of API Contract 4, where the instant of `--since` is `now`
     minus the Run Retention days plus one day, written as RFC 3339 in UTC.
     The plan is the committer time of the commit that added the Spec's
     `_prd.md` under the Spec Root, and the end the committer time of the
     commit that added `<slug>.md` under the archive root. When the window
     holds more than one such commit for a path, the one `git log` prints
     first counts.
   - An end is kept only when the Archive Record at that path in the
     working tree parses and its disposition is not `cancelled`.
   - Every Spec with a start enters the history, and no other; one without
     an end is not delivered, and one whose PRD was added before the window
     has a zero plan.
   - The source is `delivery-history`, `checked` with the number of
     delivered Specs. When no default branch is found, the `git log` fails
     or its output does not parse, the source is `unavailable` with the
     first line of the error and the history is empty. The function never
     returns an error and runs `git log` at most once.
3. MUST add the rows of `_techspec.md` → API Contracts 1, 2, 3 and 7 as
   `followUpQueueRows`, in a file of its own, with the types
   `followUpQueueRow` and `followUpQueueCell` of `_techspec.md` → Data
   Models, field for field and tag for tag. Its input carries what the view
   already read (the progress, the active Specs and the folders the Spec
   listing skipped), the history, the instant `now`, the zone the table
   prints times in and the Git runner; it reads no clock.
   - The rows come in the order of API Contract 1, and `row` counts from 1.
     `spec` is the slug and `ref` is empty for a Spec; for a Backlog Entry
     `spec` is empty, `ref` is its repository-relative path and `title` is
     its title. `title` is empty for a Spec.
   - The status is that of API Contract 2. A Spec whose Task Graph cannot be
     loaded shows `0/0`.
   - The note is that of API Contract 3: `<park class>: <next>` for a parked
     item; the first line of another blocker; `<phase> since <MM-DD HH:MM>`
     for a Spec being authored with a phase entry, with the time of its
     latest counted entry in the given zone; and for a Backlog Entry its
     title followed, with a valid Waiting Declaration, by
     `; waits on <waiting_on> since <since>`. A malformed Waiting
     Declaration adds nothing to the title. The note and the title are one
     line each: line breaks become spaces, and a text longer than 120
     characters is cut to 119 and ends with `…`. The command never
     interprets either.
   - A Spec is being authored by the rule of API Contract 7: the Spec Root
     holds no `_prd.md` for it, the archive root holds no Archive Record for
     it, and either a folder with its slug exists in the Spec Root or the
     Phase Record holds an entry of the phase `idea`, `prd`, `techspec` or
     `tasks` for it, written in the seven days before `now`, whose
     repository counts. A repository counts when it is one of the
     checkout's `roundconfig.RepositoryRoots`, or an existing directory
     whose `git remote get-url origin` names the same host, owner and
     repository as the checkout's by `parseReadinessRemote`. Paths are
     compared clean and absolute with symbolic links resolved, so an entry
     recorded through a linked path of the checkout still counts. The
     function asks the checkout for its remote once. It asks at most eight other
     directories, each once, the most recently recorded first; a directory
     that no longer exists is not asked and does not count among the eight;
     and it asks none when the checkout has no `origin`. It runs no other
     command in another directory.
   - It returns one `queueforecast.Row` for each row that is a Spec, in the
     same order, and none for a Backlog Entry. `Label` is the four digits of
     the slug. `Planned` is true for a Delivery Queue item and a planned
     Spec and false for a Spec being authored. `Queued` is true for a
     Delivery Queue item. `Started` is the Spec's start in the history, or
     zero. `Parked` is true for a parked item.
   - `Wait` is set only for a Delivery Queue item that has not started:
     `after a Run slot frees` when its blocker starts with
     `delivery.WaitRunSlot`; `after prerequisite <NNNN> is merged`, or
     `after prerequisites <NNNN>, <NNNN> are merged`, with the four digits
     of each slug the blocker lists, when it starts with
     `delivery.BlockerPrerequisiteUnmerged`; and
     `after its park is answered` for any other park.
   - It reads the Backlog Entries whose status is `open`, each through the
     reader the follow-up list uses for a Markdown record, and nothing under
     the findings directory or the repository inbox: an unreadable Finding
     or inbox note changes no row and no source. An entry that is not open
     and a Finding are never a row.
   - It returns the sources `backlog`, `checked` with the number of open
     entries, and `phase-record`, `checked` with the number of phase entries
     read. An unreadable Backlog Entry makes `backlog` `unavailable` with
     the first line of the error, and the rows of Specs remain. An
     unreadable Phase Record makes `phase-record` `unavailable`, and the
     Specs being authored come from the folders alone.
4. MUST make `readFollowUpProgress` ask the Run Database once for the
   repository's Runs of every state, as `roundfix runs list` does. The
   active ones fill `in_progress.runs` and the `runs` source exactly as
   before, and the whole list is kept for the history. MUST add the clock
   `followUpNow` to the command dependencies, `time.Now` by default, and no
   other dependency.
5. MUST leave what the command prints unchanged: the text and the JSON of
   the follow-up list and of the status view are byte for byte what they
   were for every fixture of the existing tests, with terminal Runs in the
   Run Database or without. No reader is called by the command in this
   Task.
6. MUST add these tests, each top-level test calling `t.Parallel()` first,
   with temporary repositories and homes through the packages' existing
   seams, never the real home, the real Run Database or the real Phase
   Record, and without editing `internal/cli/cli_test.go`. No recorded test
   is removed or renamed.
   - `internal/phaserecord`:
     `TestPhasesReadsBackThePhaseEntriesOfTheWindow`, with the subtests
     `oldest_first`, `prompt_entries_are_skipped`,
     `a_malformed_line_is_counted`, `another_schema_is_counted`,
     `before_since_is_left_out`, `a_missing_directory_is_an_empty_record`,
     `a_file_that_is_no_day_file_is_ignored`,
     `an_unreadable_record_is_an_error`, `nothing_is_written` and
     `written_by_start_phase`, which reads back, with its phase, Spec and
     repository, an entry the package's own recorder wrote under a clock
     the test supplies.
   - `internal/cli`, the history:
     `TestTheDeliveryHistoryIsReadFromRunsAndOneGitLog`, with the subtests
     `start_is_the_first_implement_run`,
     `end_is_the_archive_records_commit_time`,
     `plan_is_the_prds_commit_time`, `a_cancelled_record_is_not_delivered`,
     `an_unparseable_record_is_not_delivered`,
     `a_spec_without_a_run_is_left_out`, `before_the_window_is_not_read`,
     `one_git_log` (a counting runner sees one `log` with the arguments of
     API Contract 4 and the `--since` instant derived from `now`),
     `no_default_branch_is_unavailable` and `a_failed_log_is_unavailable`.
   - `internal/cli`, the rows: `TestTheQueueRowsComeInTableOrder` (one
     Delivery Queue item, one planned Spec, one folder without a PRD, one
     open Backlog Entry and one waiting Backlog Entry give five rows in that
     order, each with the status and the note of its kind, and four forecast
     rows); `TestAQueueItemNamesItsWait`, with the subtests `run_slot`,
     `one_prerequisite`, `two_prerequisites`, `another_park` and
     `a_started_item_names_none`; `TestASpecIsBeingAuthored`, with the
     subtests `a_folder_without_a_prd`, `an_entry_recorded_in_the_checkout`,
     `an_entry_recorded_in_a_linked_worktree`,
     `an_entry_recorded_in_a_clone_of_the_same_origin`, `another_remote`,
     `a_checkout_recorded_through_a_symbolic_link`,
     `older_than_seven_days`, `a_phase_that_is_no_authoring_phase`,
     `an_archive_record`, `with_a_prd_it_is_planned`,
     `a_directory_that_no_longer_exists`, `the_latest_phase_and_its_time`,
     `at_most_eight_other_directories` and `no_origin_asks_none`;
     `TestTheNoteIsOneLineOfRepositoryText`, with the subtests
     `a_long_waiting_declaration_is_cut` (a Waiting Declaration of 300
     characters with a line that reads as an instruction gives one note
     line of 120 characters ending with `…`), `line_breaks_become_spaces`,
     `a_parked_item`, `another_blocker`, `an_entry_that_is_not_open`,
     `a_malformed_waiting_declaration` and `a_finding_is_never_a_row`;
     `TestAnUnreadableSourceLeavesTheOtherRows`, with the subtests
     `backlog`, `phase_record` and `an_unreadable_finding_changes_nothing`;
     and `TestTheFollowUpProgressKeepsEveryRunForTheHistory`, which seeds one
     active and two terminal Runs and asserts that both views print what
     they print with the active Run alone and that the progress holds the
     three.
   - Time in these tests: each names one instant as its `now` and hands it
     to the reader. Every other time is derived from it: the `created_at`
     of a seeded Run, set by SQL as `runs_causes_test.go` does; the
     committer date of a fixture commit; the time of a phase entry and the
     name of its day file; and the `since` of `Phases`. No test compares a
     seeded time with the real clock, starts a Run, a Daemon or a Delivery
     Owner, sleeps or waits on a timer.
7. MUST add to `docs/user-guide/commands/follow-up.md`, right before the
   section `## Text and JSON reports`, a section `## Spec queue` with these
   three paragraphs, word for word, and change nothing else in the guide:

   > The Spec queue of a repository is its Specs and Backlog Entries in the order they will be worked: each Delivery Queue item that is not merged, by position; then each planned Spec that is not queued, by number; then each Spec being authored, by number; then each open Backlog Entry, by its `created` date and then its path. A Finding and an inbox note are never part of it.

   > A Spec is being authored when the checkout holds no PRD for it, no Archive Record names it, and either its folder exists under the Spec Root or the Phase Record holds an entry of the phase `idea`, `prd`, `techspec` or `tasks` for it written in the last seven days. An entry counts when it was recorded in this checkout or one of its linked worktrees, or in an existing directory whose `origin` remote names the same host, owner and repository as the checkout's. Roundfix asks at most eight other directories for that remote, the most recently recorded first, runs nothing else there, and asks none when the checkout has no `origin`.

   > The delivery history of a repository is the Specs that have an Implement Run in its Run Database. A Spec's start is the start of its first Implement Run. Its end is the commit time of the commit that added its Archive Record to the default branch, unless that record is a cancellation; a Spec with no such end is not delivered. The history is read from the Runs and from one `git log` over the default branch, bounded by the Run Retention days plus one day, so it is one machine's and as long as Run Retention keeps it.

   MUST then re-record the Behavior Surface Record as the Task's last step,
   after the guide edit, with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
8. MUST record in `docs/references/skill-coverage.json` the Coverage Review
   of the surface `command phase start`, whose sources hold the package
   directory this Task adds a file to: its entry carries the `review`
   `<YYYY-MM-DD> — Spec 0398 task_02 adds a reader of phase entries to the package; the phase start command, its output and the format of the record do not change`,
   with the date the Task runs. MUST change nothing else in the map.
9. MUST change no Governed Path and none of these: `internal/cli/cli_test.go`,
   `docs/references/coverage-record.json`, the existing files and tests of
   `internal/phaserecord`, any skill, any Baseline file, `.roundfixrc.yml`,
   the Makefile, the CI workflows, `go.mod`, `go.sum` or the lint
   configuration.

## Subtasks

- [ ] Add the reader of phase entries and its tests.
- [ ] Add the delivery history from the Runs and one `git log`, and its tests.
- [ ] Add the queue rows, their notes, their waits and the rule of a Spec being authored, and their tests.
- [ ] Read the Runs of every state once and add the clock.
- [ ] Write the three paragraphs of the guide, record the Coverage Review and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] A Spec is a row as being authored for a folder without a PRD, for a phase entry recorded in the checkout or one of its linked worktrees, and for one recorded in a clone with the same `origin` remote; it is not a row for an entry of another remote, for an entry older than seven days, for a phase that is no Authoring Phase or when it has an Archive Record; with a PRD it is a planned row.
- [ ] An open Backlog Entry whose Waiting Declaration holds 300 characters and a line that reads as an instruction gives one note line of at most 120 characters; an entry that is not open is not a row; a Finding is never a row and an unreadable one changes nothing.
- [ ] The delivery history takes a Spec's start from its first Implement Run and its end from the commit that added its Archive Record, leaves a cancelled record undelivered, and runs `git log` once whatever the number of Specs; a failed read is an unavailable source and an empty history.
- [ ] The readers ask the checkout for its remote once and at most eight other directories, run nothing else there, and write nothing under the repository or the home.
- [ ] The follow-up list and the status view print exactly what they printed, with terminal Runs in the Run Database or without.
- [ ] The Follow-up guide defines the Spec queue, a Spec being authored and the delivery history, and the documentation contracts pass.

## Context

- instruction: `docs/adr/0410-the-status-view-lists-specs-being-authored-and-open-backlog-entries.md`
- instruction: `docs/adr/0409-a-queue-forecast-is-a-likely-time-and-a-limit-from-the-repository-own-deliveries-or-a-named-condition.md`
- instruction: `docs/adr/0407-a-phase-record-names-the-spec-and-phase-of-supervisor-work-and-a-phase-boundary-ends-the-session.md`
- instruction: `internal/cli/runs.go`
- instruction: `internal/cli/runs_causes_test.go`
- instruction: `internal/cli/follow_up_test.go`
- instruction: `internal/cli/follow_up_status_test.go`
- instruction: `internal/cli/deliver_plan.go`
- instruction: `internal/delivery/engine.go`
- instruction: `internal/delivery/park_class.go`
- instruction: `internal/preflight/default_branch_head.go`
- creates: `internal/phaserecord/read.go`
- creates: `internal/phaserecord/read_test.go`
- creates: `internal/cli/follow_up_history.go`
- creates: `internal/cli/follow_up_history_test.go`
- creates: `internal/cli/follow_up_queue.go`
- creates: `internal/cli/follow_up_queue_test.go`
- interface: `internal/cli/follow_up_inventory.go`
- interface: `internal/cli/cli.go`
- interface: `docs/user-guide/commands/follow-up.md`
- interface: `docs/references/skill-coverage.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The follow-up progress that the list and the status view read, and the command dependencies — consumers: `./internal/cli`
- The Phase Record package — consumers: `./internal/phaserecord`
- The Follow-up guide, the Skill Coverage Map and the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 15m -v -run '^TestPhasesReadsBackThePhaseEntriesOfTheWindow$' ./internal/phaserecord 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestPhasesReadsBackThePhaseEntriesOfTheWindow/oldest_first TestPhasesReadsBackThePhaseEntriesOfTheWindow/prompt_entries_are_skipped TestPhasesReadsBackThePhaseEntriesOfTheWindow/a_malformed_line_is_counted TestPhasesReadsBackThePhaseEntriesOfTheWindow/another_schema_is_counted TestPhasesReadsBackThePhaseEntriesOfTheWindow/before_since_is_left_out TestPhasesReadsBackThePhaseEntriesOfTheWindow/a_missing_directory_is_an_empty_record TestPhasesReadsBackThePhaseEntriesOfTheWindow/an_unreadable_record_is_an_error TestPhasesReadsBackThePhaseEntriesOfTheWindow/nothing_is_written TestPhasesReadsBackThePhaseEntriesOfTheWindow/written_by_start_phase; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 15m ./internal/phaserecord && GOOS=windows go build -buildvcs=false ./internal/phaserecord` — expected: exit 0. Before this Task the reader and its test do not exist, so the pass lines are missing; afterwards phase entries are read back oldest first within their window, prompt entries are skipped, unreadable lines are counted, nothing is written, and the whole package passes and builds for Windows.

- `out="$(go test -count=1 -timeout 30m -v -run '^(TestTheDeliveryHistoryIsReadFromRunsAndOneGitLog|TestTheQueueRowsComeInTableOrder|TestAQueueItemNamesItsWait|TestASpecIsBeingAuthored|TestTheNoteIsOneLineOfRepositoryText|TestAnUnreadableSourceLeavesTheOtherRows|TestTheFollowUpProgressKeepsEveryRunForTheHistory)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheDeliveryHistoryIsReadFromRunsAndOneGitLog/start_is_the_first_implement_run TestTheDeliveryHistoryIsReadFromRunsAndOneGitLog/end_is_the_archive_records_commit_time TestTheDeliveryHistoryIsReadFromRunsAndOneGitLog/plan_is_the_prds_commit_time TestTheDeliveryHistoryIsReadFromRunsAndOneGitLog/a_cancelled_record_is_not_delivered TestTheDeliveryHistoryIsReadFromRunsAndOneGitLog/a_spec_without_a_run_is_left_out TestTheDeliveryHistoryIsReadFromRunsAndOneGitLog/before_the_window_is_not_read TestTheDeliveryHistoryIsReadFromRunsAndOneGitLog/one_git_log TestTheDeliveryHistoryIsReadFromRunsAndOneGitLog/a_failed_log_is_unavailable TestTheQueueRowsComeInTableOrder TestAQueueItemNamesItsWait/run_slot TestAQueueItemNamesItsWait/two_prerequisites TestAQueueItemNamesItsWait/another_park TestAQueueItemNamesItsWait/a_started_item_names_none TestASpecIsBeingAuthored/a_folder_without_a_prd TestASpecIsBeingAuthored/an_entry_recorded_in_the_checkout TestASpecIsBeingAuthored/an_entry_recorded_in_a_clone_of_the_same_origin TestASpecIsBeingAuthored/another_remote TestASpecIsBeingAuthored/older_than_seven_days TestASpecIsBeingAuthored/an_archive_record TestASpecIsBeingAuthored/with_a_prd_it_is_planned TestASpecIsBeingAuthored/at_most_eight_other_directories TestASpecIsBeingAuthored/no_origin_asks_none TestTheNoteIsOneLineOfRepositoryText/a_long_waiting_declaration_is_cut TestTheNoteIsOneLineOfRepositoryText/an_entry_that_is_not_open TestTheNoteIsOneLineOfRepositoryText/a_finding_is_never_a_row TestAnUnreadableSourceLeavesTheOtherRows/backlog TestAnUnreadableSourceLeavesTheOtherRows/phase_record TestAnUnreadableSourceLeavesTheOtherRows/an_unreadable_finding_changes_nothing TestTheFollowUpProgressKeepsEveryRunForTheHistory; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m -run '^TestFollowUp' ./internal/cli` — expected: exit 0. Before this Task the readers and their tests do not exist, so the pass lines are missing; afterwards the history, the rows, the waits, the rule of a Spec being authored and the notes hold by name, and every existing test of the follow-up list and of the status view still passes unchanged.

- `tr -s '[:space:]' ' ' < docs/user-guide/commands/follow-up.md | grep -qF -- 'then each Spec being authored, by number; then each open Backlog Entry' || { echo 'the Follow-up guide lacks the Spec queue' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/follow-up.md | grep -qF -- 'Roundfix asks at most eight other directories for that remote, the most recently recorded first, runs nothing else there' || { echo 'the Follow-up guide lacks the rule of a Spec being authored' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/follow-up.md | grep -qF -- 'bounded by the Run Retention days plus one day' || { echo 'the Follow-up guide lacks the delivery history' >&2; exit 1; }; test "$(grep -c -- 'Spec 0398 task_02' docs/references/skill-coverage.json)" -eq 1 || { echo 'the Coverage Review is not recorded' >&2; exit 1; }; go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the guide holds none of the three paragraphs; afterwards it defines the Spec queue, a Spec being authored and the delivery history, the Coverage Review is recorded once, and the documentation contracts, the Skill Coverage Map and the Behavior Surface Record included, pass.

## References

- `_prd.md` → User Stories 1, 5 and 9; Core Features 1, 2, 3, 5, 9, 10, 13 and 14; Success Metric 7; Success Metric 8; Decisions
- `_techspec.md` → System Architecture; History and row sources; Interfaces; Data Models; API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 7; Integration Points; Testing Approach; Risks & Considerations; Build Order 2
- `_tasks.md` → What the Tasks stage settled beyond the TechSpec; Time in tests; Spec 0397 merged while this graph was written
- ADR-0409, ADR-0410, ADR-0407, ADR-0353, ADR-0255
