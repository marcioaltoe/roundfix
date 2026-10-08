---
task: task_03
spec: 0250-a-run-database-that-keeps-only-recent-runs
status: completed
type: backend
complexity: high
---

# Task 03: The Run Retention Sweep runs once a day at the start of Runs and Delivery Queues

## Overview

Makes the Run Retention Sweep automatic. `implement`, `resolve`, `watch` and
`deliver start` run it when it is due, at most once a day, under a two-second
budget, and print one stderr line only when it removed something, paused or
failed. `roundfix runs show` and `roundfix events` name Run Retention when an
unknown Run ID predates the cutoff. The Roundfix skill and the `runs` and
`events` guides describe both, so the skill matches the shipped CLI. It is
verifiable on its own through the existing command test harnesses on
fixture homes.

## Requirements

1. MUST add `runRetentionAtStart` in the new `internal/cli/run_retention_start.go`,
   implementing `_techspec.md` → API Contract 4: the due rule over
   `LastRunRetentionSweep`, the two-second budget passed to task_02's
   `sweepRunDatabase`, the three exact stderr lines, and no change to the
   command's exit code on any outcome. Its clock is the GC Command's existing
   `now` dependency.
2. MUST call it in `implement`, `resolve` and `watch` directly after their
   existing Journal Retention prune, and in `deliver start` after the Delivery
   Queue is recorded and before that command's store closes. The Journal
   Retention prune and its stderr lines are unchanged (`_techspec.md` →
   Invariant 10).
3. MUST implement `_techspec.md` → API Contract 6 in `roundfix runs show` and
   `roundfix events`, with the exact wording of Surface Transcript 4 and
   Surface Transcript 5, the window read from the loaded configuration, and
   the existing message byte for byte for any other ID, `run_missing`
   included.
4. MUST NOT change `internal/cli/cli_test.go`, a Governed Path. When a test
   of `./internal/cli` changes because the sweep now removes a Run it seeded,
   the Task records the test and the reason in its Result; the TechSpec's
   Risks section expects none.
5. MUST describe in `.agents/skills/roundfix/references/storage.md`
   `store.run_retention_days`, the kept reasons, the once-a-day rule, the
   two-second budget with its paused line, the removal and warning lines, the
   Run Retention section of `roundfix gc --dry-run` and `roundfix gc`, and the
   conversion by `roundfix gc compact --apply`, replacing the sentences that
   say retention never deletes `runs` rows and that nothing compacts
   automatically. MUST add the unknown Run hint to
   `.agents/skills/roundfix/references/runs.md`,
   `.agents/skills/roundfix/references/events.md`,
   `docs/user-guide/commands/runs.md` and `docs/user-guide/commands/events.md`.
   MUST then run `make skills-sync` and
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   which raises both version fields of the Roundfix skill to the next version
   at record time and records it in `skills/testdata/owned-skill-versions.json`.
   No `### QA settlement` section changes.
6. MUST add these tests to `internal/cli/run_retention_start_test.go`, each on
   a fixture home:
   - `TestRunRetentionRunsOnceADayAtRunStart`: a due sweep removes an old
     Run and prints the removal line; a second call within 24 hours removes
     nothing, even with a newly aged Run, and prints nothing; a call 24 hours later, or after
     the window changes from 30 to 7, sweeps again;
   - `TestRunRetentionStopsOnItsBudgetAndResumes`: with a stepping clock that
     spends the budget after the first removal, the call removes one of three
     old Runs, prints the paused line and records no completion; the next call
     removes the rest and records completion;
   - `TestRunRetentionAtStartWarnsAndNeverBlocks`: a store error prints the
     warning line, and the start command still reaches its next step;
   - `TestImplementStartRunsRunRetention`: through the existing implement
     test harness, `roundfix implement` on a fixture home with an old
     terminal Run removes it before the Run is created;
   - `TestDeliverStartRunsRunRetention`: through the existing deliver test
     harness, `deliver start` removes an old Run that no queue references and
     keeps the Runs the recorded queue references;
   - `TestDeliverStatusIsUnchangedByRunRetention`: `deliver status` prints the
     same lines and token totals before and after a sweep that removes every
     other old Run;
   - `TestRunRetentionKeepsTheRunReconcileReads`: an old terminal Run whose
     Run Worktree exists stays after the sweep, and `roundfix reconcile`
     still lists it;
   - `TestUnknownRunNamesRunRetention`: Surface Transcript 4 and Surface
     Transcript 5 exactly, and `run_missing` with the existing message.

## Subtasks

- [ ] Add the due rule, the budget and the stderr lines.
- [ ] Call the sweep from the four start paths.
- [ ] Add the unknown Run hint to `runs show` and `events`.
- [ ] Update the skill references, the two guides, the mirror and the skill version.
- [ ] Write the eight tests and run the affected `./internal/cli` and `./skills` tests.

## Acceptance Criteria

- [ ] A Run or Delivery Queue start sweeps at most once a day and never waits more than the budget plus one step.
- [ ] A queue-referenced Run and a Run whose Run Worktree exists survive every automatic sweep.
- [ ] An unknown old Run ID names Run Retention, and every other refusal is unchanged.
- [ ] The Roundfix skill describes Run Retention and its mirror is byte-identical.

## Context

- creates: `internal/cli/run_retention_start.go`
- interface: `internal/cli/implement.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/deliver.go`
- interface: `internal/cli/runs_show.go`
- interface: `internal/cli/events.go`
- interface: `.agents/skills/roundfix/references/storage.md`
- interface: `.agents/skills/roundfix/references/runs.md`
- interface: `.agents/skills/roundfix/references/events.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/references/storage.md`
- interface: `skills/roundfix/references/runs.md`
- interface: `skills/roundfix/references/events.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/runs.md`
- interface: `docs/user-guide/commands/events.md`
- creates: `internal/cli/run_retention_start_test.go`
- instruction: `docs/adr/0255-run-retention-removes-terminal-runs-whole-and-compacts-incrementally.md`
- instruction: `docs/adr/0225-a-queue-item-runs-the-roundfix-binary-its-branch-builds.md`
- instruction: `docs/adr/0053-terminal-run-worktree-reconciliation-is-proof-based.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestRunRetentionRunsOnceADayAtRunStart|TestRunRetentionStopsOnItsBudgetAndResumes|TestRunRetentionAtStartWarnsAndNeverBlocks|TestImplementStartRunsRunRetention|TestDeliverStartRunsRunRetention|TestDeliverStatusIsUnchangedByRunRetention|TestRunRetentionKeepsTheRunReconcileReads|TestUnknownRunNamesRunRetention)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for t in TestRunRetentionRunsOnceADayAtRunStart TestRunRetentionStopsOnItsBudgetAndResumes TestRunRetentionAtStartWarnsAndNeverBlocks TestImplementStartRunsRunRetention TestDeliverStartRunsRunRetention TestDeliverStatusIsUnchangedByRunRetention TestRunRetentionKeepsTheRunReconcileReads TestUnknownRunNamesRunRetention; do printf '%s\n' "$out" | grep -q -- "--- PASS: $t " || { printf 'missing passing test %s\n' "$t" >&2; exit 1; }; done; s="$(tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/storage.md)"; for p in 'store.run_retention_days' 'Run Retention paused after its 2s budget' 'Run Retention removed runs='; do printf '%s\n' "$s" | grep -qF -- "$p" || { printf 'missing phrase in the storage reference: %s\n' "$p" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/runs.md .agents/skills/roundfix/references/events.md docs/user-guide/commands/runs.md docs/user-guide/commands/events.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Run Retention may have removed it' || { printf 'missing the unknown Run hint in %s\n' "$f" >&2; exit 1; }; done; diff -r .agents/skills/roundfix skills/roundfix >/dev/null || { printf 'the Roundfix skill mirror differs from its canonical copy\n' >&2; exit 1; }; go test -count=1 ./skills && go test -count=1 -run 'Retention|RunsShow|Events|DeliverStatus|DeliverStart|Reconcile|GC' ./internal/cli` — expected: exit 0. Before this Task the eight tests do not exist and the storage reference names no `store.run_retention_days`, so the first check fails. After it, the new tests pass, the skill references and guides carry the new text, the mirror matches, the skill version record holds, and the affected `./internal/cli` tests pass.

## References

- `_prd.md` → Core Feature 2; Core Feature 3; Core Feature 5; Core Feature 6; User Story 1; User Story 3; User Story 5; Success Metric 3; Success Metric 6
- `_techspec.md` → API Contract 4; API Contract 6; Surface Transcript 4; Surface Transcript 5; Invariants 2, 3, 8, 9, 10; Risks & Considerations; Build Order 3
- ADR-0255; ADR-0225; ADR-0053


## Result

Implemented this Task's slice for Daemon Verification. The inherited
`status: in_progress` is unchanged; no declared Verification command was run,
no other Task or Task Graph was edited, and no commit, push or Pull Request
was made. Before implementation, the start helper and its eight tests were
absent, and the storage reference had no `store.run_retention_days` description.

### Implementation and acceptance evidence

| Acceptance criterion | Implementation and focused-check evidence |
| --- | --- |
| A Run or Delivery Queue start sweeps at most once a day and stays within the budget plus one step | `runRetentionAtStart` checks the durable completion and configured window, then passes the existing sweep a two-second budget. Its calls immediately follow Journal Retention in `implement`, `resolve` and `watch`; `deliver start` calls it after recording the queue and before closing its store. `TestRunRetentionRunsOnceADayAtRunStart` checks the exact removal line, silence at 23 hours despite a newly aged Run, and another sweep at exactly 24 hours or after changing 30 to 7 days. `TestRunRetentionStopsOnItsBudgetAndResumes` spends the stepping-clock budget after one of three removals, checks the exact paused line and absent completion, then resumes the other two and records completion. `TestRunRetentionAtStartWarnsAndNeverBlocks` injects a real store error and checks the warning while the implement harness still executes its Agent and exits 0. `TestImplementStartRunsRunRetention` uses a fixture SQLite trigger that refuses Run creation if the old terminal Run still exists, proving the removal precedes creation. |
| Queue-referenced Runs and Runs whose Run Worktree exists survive automatic sweeps | `TestDeliverStartRunsRunRetention` seeds queue references as the command records its new queue, so an incorrectly ordered earlier sweep would remove them. It proves both the item-referenced Run and a separately linked Run survive while the unreferenced Run leaves before owner launch. `TestDeliverStatusIsUnchangedByRunRetention` compares complete status output, including the 123-token total, before and after removing every other old Run. `TestRunRetentionKeepsTheRunReconcileReads` keeps an old terminal Run with a real Run Worktree and proves public reconcile still lists its Run ID and Worktree. The daily test also preserves its queue and Worktree candidates across successive sweeps. |
| An unknown old Run names Run Retention, while other refusals stay unchanged | `unknownRunMessage` validates the timestamp and hexadecimal suffix and compares creation time with the loaded configuration's cutoff through the existing GC clock. `TestUnknownRunNamesRunRetention` asserts stdout, stderr and exit 2 byte for byte for Surface Transcripts 4 and 5, windows 30/15/7, `run_missing`, an ID exactly at the cutoff, a recent ID, invalid dates, non-hex suffixes and empty suffixes. |
| The Roundfix skill describes Run Retention and its mirror is byte-identical | Updated the canonical storage reference with the User Config window, kept reasons, daily schedule, budget and all three diagnostics, GC sections, automatic incremental compaction and explicit conversion. Added the unknown Run hint to both skill references and both user guides, and clarified the distinction between Journal Retention and whole-Run usage removal. Required synchronization and version recording generated both Roundfix version fields at `0.1.56` and the version record. A Python byte comparison confirms all 20 canonical and mirrored files match. Digest regeneration reports no changes. No `### QA settlement` section was changed. |

### Checks run

- `GOCACHE=/private/tmp/roundfix-task03-gocache rtk proxy go test ./internal/cli -run '^(TestRunRetention|TestImplementStartRuns|TestDeliverStartRuns|TestDeliverStatusIs|TestUnknownRun)' -count=1 -v`: exit 0; all eight required tests passed on fixture homes.
- `GOCACHE=/private/tmp/roundfix-task03-gocache rtk proxy go test ./internal/cli -count=1 -v`: exit 0 with required process-table access; the entire CLI package passed, including its repository boundary guard. Log: `/private/tmp/roundfix-task03-cli-check.log`.
- `rtk make skills-sync`: exit 0; mirrored canonical references through the sanctioned generator.
- `GOCACHE=/private/tmp/roundfix-task03-gocache rtk proxy go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`: exit 0; ran after the source edits and again after the final documentation correction, generating the final version and record.
- `GOCACHE=/private/tmp/roundfix-task03-gocache rtk proxy go test ./skills -run 'TestEveryOwnedSkillVersionIsRecorded|TestRecordingRewritesBothVersionFieldsOfASkillAndItsMirror|TestOwnedSkill' -count=1`: exit 0 against the final skill content.
- `GOCACHE=/private/tmp/roundfix-task03-gocache rtk make baseline-digests`: exit 0; `changed: false`, no derived pins rewritten.
- `rtk proxy git -c core.fsmonitor=false diff --check`: exit 0. Python comparison of the canonical and mirrored Roundfix skill: 20 files byte-identical.

The first full CLI attempt used a status fixture that lacked its durable Run
link, lacked sandbox permission to enumerate owned process trees in two
existing force-stop tests, and overlapped implementation edits detected by
suiteguard. The fixture now records its link through the queue API; the final
full run had process access and no concurrent repository edits and passed.
No existing test changed because Run Retention removed a seeded Run;
`internal/cli/cli_test.go` remains untouched. No follow-up implementation was
added outside this Task.
