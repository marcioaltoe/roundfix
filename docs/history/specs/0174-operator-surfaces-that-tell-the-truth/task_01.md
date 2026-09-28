---
task: task_01
spec: 0174-operator-surfaces-that-tell-the-truth
status: completed
type: backend
complexity: high
---

# Task 01: A schema refusal names its remedy and roundfix migrate performs it

## Overview

`openReader` in `internal/store/store.go` refuses a Run Database whose schema version differs from the binary's with a `SchemaVersionError` that tells the operator to run `resolve`, `watch` or `implement`. Every command that only reads the Run Database (`runs list`, `events`, `attach`, the Run Browser, `doctor`, `gc`, `storage report`, `reconcile`, `deliver status`, `spec audit`, and the Run Database checks of `settle` and `reopen`) prints that advice, which needs a Pull Request or a Spec to follow, and prints it even when the database is newer than the binary. The Run Database is the machine-wide file at `~/.roundfix/roundfix.db`, written by whichever Roundfix binary opened it last; its reader is the operator or Supervisor monitoring Runs. `migrate` also reads the schema version and derives its statements before it takes the machine-wide write lock, so two processes opening an older database at once can replay the same migration.

## Requirements

1. MUST keep every command that only reads the Run Database read-only: `openReader` never migrates, and `OpenReader` and `OpenStorageReader` keep their signatures and behavior for a current database.
2. MUST keep `SchemaVersionError`'s type and its `Path`, `Found` and `Supported` fields, and make `Error()` direction-aware. For an older database it reads `Run Database "<path>" has schema version <found>, older than the schema version <supported> this binary supports; run 'roundfix migrate' to upgrade it`; for a newer one `Run Database "<path>" has schema version <found>, newer than the schema version <supported> this binary supports; a newer roundfix wrote it, so use that binary or run 'roundfix upgrade'`. No message may name `resolve, watch, or implement` any longer.
3. MUST make `store.Open` return `SchemaVersionError` for a newer database, replacing `migrate Run Database: schema version N is not supported`, keeping its signature.
4. MUST make the migration read the schema version, derive its statements (including `deliveryOwnerMigrationStatements` through `deliveryWorktreeProvisioningMigrationStatements`), apply them and record the new version under one hold of the machine-wide write lock, so a process that finds the database already at `schemaVersion` once it holds the lock writes nothing. `store.Open` and `store.Migrate` MUST share this one path.
5. MUST add `store.Migrate(ctx, homeDir)` returning the database path and the schema version before and after. It MUST NOT create a missing Run Database or its directory, reporting it as absent; it MUST write nothing to a current database; it MUST refuse a newer one with `SchemaVersionError` and leave the database file byte-identical.
6. MUST add `roundfix migrate` in `internal/cli/migrate.go`, dispatched from `run` in `internal/cli/cli.go`, listed in the root usage and given a `commandUsage` entry. It takes no flags or arguments beyond help, needs no Git repository and reads config only for the home directory. It prints exactly one stdout line and exits `0` for `Run Database migrated from schema version <from> to <to>: <path>`, `Run Database is already at schema version <n>: <path>` and `No Run Database at <path>; nothing to migrate`; it prints the newer refusal on stderr and exits `2`; an unexpected argument is a Preflight failure with exit `2`; a migration error exits `1`.
7. MUST update the assertion of `TestOpenReaderRejectsMismatchedSchemaVersion` in `internal/store/store_test.go` that expects `resolve, watch, or implement` to expect `run 'roundfix migrate'`, keeping the test's name; no other existing test is renamed or removed.
8. MUST put the new store tests in `internal/store/migrate_test.go` and the new command tests in `internal/cli/migrate_test.go`, and add `TestMigrateHelpDocumentsTheContract` to `internal/cli/cli_test.go`, asserting that `roundfix migrate --help` prints the command's usage and that `roundfix --help` lists `roundfix migrate`. An older fixture is a fresh database downgraded the way `seedOutdatedV9RunDatabase` does; a newer fixture advances `user_version` past `schemaVersion`. Concurrency is proved with goroutines released together, never with a fixed wait.
9. MUST add a `### migrate` section to `docs/user-guide/commands.md` naming the four outcomes and stating that a Run started by an older binary that is still Active meets the same migrated database any operational command would leave, and state both refusals in its Global contract, using the phrase `run 'roundfix migrate'`.
10. MUST state in the Run discovery and Attach section of `.agents/skills/roundfix/SKILL.md` that `runs list` refuses a Run Database of another schema version, that an older one is upgraded with `roundfix migrate` and a newer one needs the newer binary, and regenerate `skills/roundfix/SKILL.md` with `make skills-sync`, then run `make baseline-digests`.
11. MUST add a **Migrate Command** entry to `CONTEXT.md` beside the Doctor and GC Command entries: the support command that upgrades an older Run Database to the binary's schema version under the machine-wide write lock, and writes nothing to a current, absent or newer one.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] An older Run Database is refused by a read-only command with `run 'roundfix migrate'`, and a newer one with `roundfix upgrade`, while the reader writes nothing.
- [ ] `roundfix migrate` upgrades an older Run Database, after which `runs list` exits `0`; it reports a current or absent database without writing, and refuses a newer one with exit `2` and the file byte-identical.
- [ ] Concurrent writer opens of an older Run Database all succeed and leave it at the supported schema version.
- [ ] The journal consumer corpus harnesses compile and replay unchanged.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/store/store.go`
- interface: `internal/store/store_test.go`
- creates: `internal/store/migrate_test.go`
- creates: `internal/cli/migrate.go`
- creates: `internal/cli/migrate_test.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/cli_test.go`
- interface: `docs/user-guide/commands.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestOpenReaderRejectsMismatchedSchemaVersion|TestSchemaReviewSkippedReaderRejectsNewerDatabase|TestOpenReaderRefusesAnOlderDatabaseNamingMigrate|TestOpenReaderRefusesANewerDatabaseNamingUpgrade|TestOpenRefusesANewerDatabaseWithTheTypedError|TestConcurrentOpensMigrateAnOlderDatabaseOnce|TestMigrateUpgradesAnOlderDatabase|TestMigrateLeavesACurrentDatabaseUnwritten|TestMigrateCreatesNothingForAnAbsentDatabase|TestMigrateRefusesANewerDatabaseWithoutWriting|TestJournalConsumerCorpusReplaysEveryConsumer|TestRunsListOnAnOlderDatabaseNamesMigrate|TestRunsListOnANewerDatabaseNamesUpgrade|TestMigrateCommandUpgradesAnOlderDatabaseThenRunsListWorks|TestMigrateCommandReportsACurrentDatabase|TestMigrateCommandReportsAnAbsentDatabaseAndCreatesNothing|TestMigrateCommandRefusesANewerDatabase|TestMigrateCommandRefusesAnArgument|TestMigrateHelpDocumentsTheContract)$" ./internal/store ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestOpenReaderRejectsMismatchedSchemaVersion TestSchemaReviewSkippedReaderRejectsNewerDatabase TestOpenReaderRefusesAnOlderDatabaseNamingMigrate TestOpenReaderRefusesANewerDatabaseNamingUpgrade TestOpenRefusesANewerDatabaseWithTheTypedError TestConcurrentOpensMigrateAnOlderDatabaseOnce TestMigrateUpgradesAnOlderDatabase TestMigrateLeavesACurrentDatabaseUnwritten TestMigrateCreatesNothingForAnAbsentDatabase TestMigrateRefusesANewerDatabaseWithoutWriting TestJournalConsumerCorpusReplaysEveryConsumer TestRunsListOnAnOlderDatabaseNamesMigrate TestRunsListOnANewerDatabaseNamesUpgrade TestMigrateCommandUpgradesAnOlderDatabaseThenRunsListWorks TestMigrateCommandReportsACurrentDatabase TestMigrateCommandReportsAnAbsentDatabaseAndCreatesNothing TestMigrateCommandRefusesANewerDatabase TestMigrateCommandRefusesAnArgument TestMigrateHelpDocumentsTheContract; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && ! grep -q "resolve, watch, or implement" internal/store/store.go && grep -q "### migrate" docs/user-guide/commands.md && grep -q "run 'roundfix migrate'" docs/user-guide/commands.md && grep -q "roundfix migrate" .agents/skills/roundfix/SKILL.md && grep -q "Migrate Command" CONTEXT.md && diff -r .agents/skills/roundfix skills/roundfix >/dev/null` — expected: exit 0; before this Task none of the new named tests exists, the refusal still names `resolve, watch, or implement` and no guide names `roundfix migrate`, so the command fails.

## References

- [_techspec.md](_techspec.md) — The schema refusal and the migrate command
- `_prd.md` → Goals 1-2; Core Feature 1; Success Metric 1
- `_techspec.md` → API Contracts 1-4; Testing Approach 1

## Result

### Implementation

- Read-only Run Database opens remain read-only and now return the preserved
  `SchemaVersionError` with an older/migrate or newer/upgrade remedy.
- Writer opens and `store.Migrate` share one migration path that holds the
  machine-wide write lock while reading the version, deriving conditional
  statements, applying them, and recording the supported version.
- `store.Migrate` reports the database path and before/after versions, creates
  nothing for an absent database, writes nothing for a current database, and
  refuses a newer database without changing its bytes.
- `roundfix migrate` is available without a Git repository, prints one
  deterministic stdout outcome for migrated/current/absent databases, maps a
  newer database or invalid argument to exit 2, and maps migration failures to
  exit 1.
- The command reference, Roundfix skill, shipped skill mirror, and glossary
  describe the explicit migration and direction-aware refusal contracts.

### Focused checks

- `GOCACHE=/tmp/roundfix-task01-gocache go test -count=1 -run
  'Test(OpenReaderRejectsMismatchedSchemaVersion|SchemaReviewSkippedReaderRejectsNewerDatabase|OpenReaderRefusesAnOlderDatabaseNamingMigrate|OpenReaderRefusesANewerDatabaseNamingUpgrade|OpenRefusesANewerDatabaseWithTheTypedError|ConcurrentOpensMigrateAnOlderDatabaseOnce|MigrateUpgradesAnOlderDatabase|MigrateLeavesACurrentDatabaseUnwritten|MigrateCreatesNothingForAnAbsentDatabase|MigrateRefusesANewerDatabaseWithoutWriting|JournalConsumerCorpusReplaysEveryConsumer|WriteTxIsTheOnlyWriterTransaction)$'
  ./internal/store` — passed.
- `GOCACHE=/tmp/roundfix-task01-gocache go test -count=1 -run
  'Test(Migrate|RunsListOnA)' ./internal/cli` — passed.
- `GOCACHE=/tmp/roundfix-task01-gocache go test -count=1 ./internal/store` —
  passed.
- `make skills-sync` — passed; canonical and shipped Roundfix skills were
  synchronized.
- `make baseline-digests` — passed and reported no derived changes.
- `GOCACHE=/tmp/roundfix-task01-gocache go test -count=1 ./internal/cli` — the
  Task-related tests passed, but the broader package run was blocked by the
  sandbox denying process-table access in two pre-existing force-stop
  integration tests. This is recorded as an environment limitation, not Task
  evidence.

### Acceptance evidence

- Older and newer reader fixtures are refused with `roundfix migrate` and
  `roundfix upgrade` respectively, and byte comparisons prove the reader did
  not change either database.
- The migrate command upgrades an older fixture and the following `runs list`
  succeeds; current and absent fixtures remain unwritten, while a newer fixture
  exits 2 and remains byte-identical.
- Eight writer opens released from one barrier all succeed against one older
  database, which ends at the supported schema version.
- The journal consumer corpus compiles and replays every consumer unchanged.

### Follow-ups

- None discovered inside Task 01's slice.
