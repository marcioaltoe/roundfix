---
task: task_02
spec: 0221-a-delivery-queue-that-runs-the-binary-its-item-builds
status: completed
type: backend
complexity: low
---

# Task 02: roundfix migrate --check says whether this binary would change the Run Database

## Overview

The queue owner needs to ask an item binary, before trusting it with the
shared Run Database, whether it would migrate that database. This Task adds
`roundfix migrate --check`, which reads the database through the read-only
reader and reports, without writing, whether its schema version is the one
this binary supports. It is demoable on its own through the CLI.

## Requirements

1. MUST accept exactly one argument, `--check`, in `runMigrateCommand`, and
   keep the existing `unexpected argument` refusal, with its exit code `2`,
   for every other argument or for `--check` followed by another argument.
2. MUST implement the four outcomes of "The migration check" exactly: an
   absent database prints `No Run Database at <path>; nothing to migrate` on
   stdout and exits `0`, creating neither the database nor its directory; a
   current database prints
   `Run Database is at schema version <n>, the version this binary supports: <path>`
   on stdout and exits `0`; a `store.SchemaVersionError` prints
   `roundfix: migrate check: <error>` on stderr and exits `2`; any other
   failure prints `roundfix: migrate check failed: <error>` on stderr and
   exits `1`.
3. MUST read the database only through `store.OpenReader`, never through a
   writer, `store.Migrate` or the machine-wide write lock, and MUST close the
   reader before exiting.
4. MUST change the root usage line and the `migrate` help text in
   `internal/cli/cli.go` to `roundfix migrate [--check]` and say in the help
   text that `--check` reports without writing.
5. MUST add the tests named in Verification to the new file
   `internal/cli/migrate_check_test.go`, reusing the seeds of
   `internal/cli/migrate_test.go` (`seedOutdatedV9RunDatabase`,
   `seedNewerRunDatabase`) and a current database created by the store in a
   disposable home. Each test MUST assert stdout, stderr and the exit code of
   API Contract 1 and Surface Transcripts 1 to 3, and that the database file
   keeps its bytes and `PRAGMA user_version`, and that no file other than
   SQLite's `-wal` and `-shm` sidecars appears in the Roundfix Home.
6. MUST NOT edit `internal/cli/migrate_test.go`, the store package or the
   Run Database schema.

## Subtasks

- [ ] Parse `--check` and keep the argument refusal.
- [ ] Report the four outcomes through the read-only reader.
- [ ] Update the usage and help text.
- [ ] Add the check tests.

## Acceptance Criteria

- [ ] A current and an absent database exit `0` with the exact stdout lines;
      an older and a newer database exit `2` with the remedy on stderr.
- [ ] No outcome changes the database bytes or version, and an absent
      database leaves no `.roundfix` directory.
- [ ] `roundfix migrate --check extra` and `roundfix migrate extra` are refused
      with `unexpected argument` and exit `2`.

## Context

- interface: `internal/cli/migrate.go`
- interface: `internal/cli/cli.go`
- creates: `internal/cli/migrate_check_test.go`
- instruction: `internal/cli/migrate_test.go`
- instruction: `internal/store/store.go`
- instruction: `docs/user-guide/commands/migrate.md`
- instruction: `.agents/skills/roundfix/references/setup.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestMigrateCheckReportsACurrentDatabase|TestMigrateCheckReportsAnOlderDatabaseWithoutMigrating|TestMigrateCheckReportsANewerDatabaseWithoutWriting|TestMigrateCheckReportsAnAbsentDatabaseAndCreatesNothing|TestMigrateCheckRefusesAnExtraArgument)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestMigrateCheckReportsACurrentDatabase TestMigrateCheckReportsAnOlderDatabaseWithoutMigrating TestMigrateCheckReportsANewerDatabaseWithoutWriting TestMigrateCheckReportsAnAbsentDatabaseAndCreatesNothing TestMigrateCheckRefusesAnExtraArgument; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the five tests do not exist, so the command fails.
- `tmp="$(mktemp -d)" || exit 1; go build -buildvcs=false -o "$tmp/roundfix" ./cmd/roundfix || exit 1; mkdir "$tmp/home" || exit 1; out="$(env -u NODE_OPTIONS HOME="$tmp/home" "$tmp/roundfix" migrate --check 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -qF -- "nothing to migrate" || { printf 'unexpected output: %s\n' "$out" >&2; exit 1; }; test ! -e "$tmp/home/.roundfix" || { printf 'migrate --check created %s\n' "$tmp/home/.roundfix" >&2; exit 1; }` — expected: exit 0; before this Task the built binary refuses the flag with exit 2, so the command fails.

## References

- `_prd.md` → User Story 5; Core Feature 7; Success Metric 4; Goal 4
- `_techspec.md` → The migration check; API Contract 1; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Testing Approach 1; Build Order 2
- ADR-0225

## Result

Implemented `migrate --check` through `store.OpenReader`. The check returns the
specified stdout and exit `0` for a current or absent Run Database, stderr
with the existing migration or upgrade remedy and exit `2` for a schema
mismatch, and `migrate check failed` with exit `1` for other read failures.
It closes a successful reader before reporting, including on a version-read
failure, and handles close errors. It calls neither the writer nor migration
nor the machine-wide write lock. Root usage and migrate help now show
`roundfix migrate [--check]` and explain that the check reports without writing.

Focused evidence for the acceptance criteria:

- Current and absent exact stdout, empty stderr, and exit `0`, plus older and
  newer exact stderr remedies, empty stdout, and exit `2`: the four named
  `TestMigrateCheckReports...` acceptance tests passed.
- Database preservation: the current, older, newer, and invalid-argument
  tests compare database bytes and `PRAGMA user_version` before and after,
  inventory Roundfix Home, permit only new SQLite `-wal` and `-shm` sidecars,
  and check preservation of existing files. The absent test confirms neither
  `.roundfix` nor any other home entry appears.
- Argument refusal: `TestMigrateCheckRefusesAnExtraArgument` passed for
  `--check extra`, `extra`, an unknown flag, duplicate `--check`, and
  `--check --help`, asserting exact streams and exit `2` as well as database
  preservation.
- Additional focused cases assert the exact read-failure diagnostic and exit
  `1` without creating files, and the root and migrate help contract.

Focused checks run:

- Initial `go test -count=1 -run '^TestMigrateCheckReportsAnAbsent'
  ./internal/cli` could not use the sandbox-restricted default Go cache.
- With `GOCACHE=/tmp/roundfix-task02-gocache`, the same focused test reproduced
  the pre-change refusal: exit `2`, `unexpected argument "--check"`.
- After implementation, `GOCACHE=/tmp/roundfix-task02-gocache go test
  -count=1 -run '^TestMigrate' ./internal/cli` exited `0`.
- After the final test additions, `GOCACHE=/tmp/roundfix-task02-gocache go test
  -count=1 -v -run '^TestMigrate' ./internal/cli` exited `0`; all seven new
  check tests and the existing migration command tests passed. Commands were
  invoked through `rtk proxy` to preserve output.

The store package, schema, existing `migrate_test.go`, Task Graph, and other
Task files were not edited. Task status and checkboxes remain unchanged.
Declared Verification was not run; the Daemon owns it and Task settlement.
No follow-up work was identified.

## Carry-forward provenance

- Source Run: `run_20261003T234640Z_3e854d86e9475f64`
- Source commit: `e37eb957135be17ee823f68d4a772f3871284750`
