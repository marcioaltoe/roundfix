---
spec: 0174-operator-surfaces-that-tell-the-truth
status: active
created: 2026-09-28
surfaces: [backend, cli, docs]
---

# Operator surfaces that tell the truth

## Executive Summary

Make a schema-version refusal name the one action that resolves it and ship
`roundfix migrate` to take it, serialize the migration under the machine-wide
write lock, let a bare `help` mean help only as a command's first argument,
bound the QA Report front matter exactly as the derived QA Verification does,
and govern `skills/_ownership.yml` while running every `repocontract` test in
`make verify-docs`.

## Project Constraints

- Identifier strategy: not applicable — no new identifier. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files, Git and the Run
  Database only; no credential and no network call. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0004, ADR-0009, ADR-0014, ADR-0015,
  ADR-0057, ADR-0080, ADR-0081, ADR-0091, ADR-0093, ADR-0096, ADR-0097,
  ADR-0104, ADR-0117, ADR-0126, ADR-0130, ADR-0133, ADR-0149, ADR-0155 and
  ADR-0156 hold; ADR-0020, ADR-0022, ADR-0038, ADR-0056, ADR-0127, ADR-0159
  and ADR-0160 do not apply. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization recorded in
  [_authorization.md](_authorization.md); bounded files:
  `internal/cli/cli_test.go`, `.agents/skills/roundfix/SKILL.md`,
  `skills/roundfix/SKILL.md`, `.agents/skills/qa-gate/SKILL.md`,
  `skills/qa-gate/SKILL.md`, `internal/speccheck/governed.go`,
  `internal/speccheck/governed_repocontract_test.go`, `Makefile`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## The schema refusal and the migrate command

**Decision (recorded for Backlog Entry B04): commands that only read the Run
Database refuse with an exact, direction-aware instruction; they neither
migrate nor read another schema.** `openReader` in `internal/store/store.go`
keeps opening a read-only connection. Migrating from a reader would change the
Run Database under an older binary's Active Run without a request, and
`OpenStorageReader` must leave the database byte-identical. Reading older
schemas would make every query schema-aware.

`SchemaVersionError` keeps its `Path`, `Found` and `Supported` fields and its
type; only `Error()` changes, by direction:

- older: `Run Database "<path>" has schema version <found>, older than the
  schema version <supported> this binary supports; run 'roundfix migrate' to
  upgrade it`
- newer: `Run Database "<path>" has schema version <found>, newer than the
  schema version <supported> this binary supports; a newer roundfix wrote it,
  so use that binary or run 'roundfix upgrade'`

`store.Open` returns the same typed error for a newer database instead of
`migrate Run Database: schema version N is not supported`, and keeps its
signature; `OpenReader` and `OpenStorageReader` keep theirs, so the journal
consumer corpus harnesses in `internal/store/testdata/` compile unchanged.

`migrate` today reads `MigrationVersion` and derives the conditional delivery
statements (`deliveryOwnerMigrationStatements` through
`deliveryWorktreeProvisioningMigrationStatements`) before `withWriteTx` takes
the machine-wide write lock, so a second process that read the old version
replays statements the first already applied. The migration now reads the
version, derives its statements, applies them and records the new version under
one hold of the write lock; a process that finds the database already at
`schemaVersion` once it holds the lock writes nothing. `store.Open` and the new
command share this one path.

A new exported `store.Migrate(ctx, homeDir)` returns the database path and the
schema version before and after. It never creates a missing Run Database or its
directory: an absent database is reported as absent. It refuses a newer
database with `SchemaVersionError` without writing.

`internal/cli/migrate.go` adds `roundfix migrate`, dispatched from `run` in
`internal/cli/cli.go` and listed in the root usage and `commandUsage`. It takes
no flags or arguments beyond help, needs no Git repository and loads config
only for the home directory, as `runs list` does. Output:

- older database migrated: stdout `Run Database migrated from schema version
  <from> to <to>: <path>`, exit `0`;
- already current: stdout `Run Database is already at schema version <n>:
  <path>`, exit `0`, nothing written;
- absent: stdout `No Run Database at <path>; nothing to migrate`, exit `0`,
  nothing created;
- newer: the newer refusal on stderr, exit `2`, nothing written;
- an unexpected argument: Preflight failure, exit `2`;
- a migration error: stderr, exit `1`.

`docs/user-guide/commands.md` gains a `### migrate` section and the refusal in
the Global contract; the Run discovery section of the Roundfix skill names the
refusal and `roundfix migrate`; `CONTEXT.md` gains a **Migrate Command** entry.

## The help predicate

`commandWantsHelp` in `internal/cli/cli.go` changes once and every caller
inherits it:

- `help` counts only when it is `args[0]`;
- `-h` and `--help` count at any position before the first `--`;
- no token after `--` counts.

A subcommand dispatcher passes `args[1:]` to its subcommand, so `roundfix runs
list help` and `roundfix skills check help` still print usage. The top-level
dispatch in `run` is unchanged. A value spelled `-h` or `--help` keeps
requesting usage, because no Roundfix flag accepts one as a valid slug, ID,
path or reason. `docs/user-guide/commands.md` states the rule in its Global
contract.

## The QA Report front matter

`readQAReport` in `internal/spec/qa.go` stops using the shared
`splitFrontmatter` and bounds the front matter as the awk reader in
`DerivedQAVerification` does:

- the first line must be exactly `---`;
- the front matter closes at the first later line that is exactly `---`;
- inside it, exactly one line begins with `verdict:` in column one, and its
  value trimmed of whitespace is non-empty.

Lines split on `\n` and compare byte-exactly, so a `----` line, an indented
`verdict:` or a `\r` before the newline reads as the awk reader reads it. A
report that breaks a bound is a `QAReportError` naming it: `QA Report front
matter is empty`, `QA Report front matter must open with a "---" first line`,
`QA Report front matter has no closing "---" line`, `QA Report front matter has
<n> "verdict:" lines; expected exactly 1`; a missing `verdict:` line keeps
`frontmatter has no verdict field`. The bounded lines are then YAML-decoded
exactly as today, and the body starts after the closing line. Everything after
that, including the hollow-report check, the verdict switch and the blocked-row
rules, is unchanged. `splitFrontmatter` keeps serving Task files and every other
caller.

`settleQAVerdict` in `internal/daemon/task_engine.go` already settles an
unreadable report as verdict `unreadable` and the QA Task `failed`; it now
returns the read error, and the QA Task's reason becomes `QA verdict
unreadable: <cause>`. `fail`, `pending` and `missing` keep `QA verdict
<verdict>`, and a refused `pass` or `partial` keeps `QA verdict <verdict> not
accepted: <cause>`. The Daemon's seed already writes one front matter; the
qa-gate skill tells the QA Agent to keep it the report's only one.

## The governed set and the repository gate

`skills/_ownership.yml` joins the historical bounded-path literal list in
`internal/speccheck/governed.go`, under the clause ADR-0130 names, and
`TestGovernedSetOnlyGrows` in `internal/speccheck/governed_repocontract_test.go`
lists it as newly governed.

`internal/authorization` and `internal/verifyselect` spawn processes in their
tests but install no suite guard. Each gains a `main_test.go` whose `TestMain`
calls `suiteguard.Main` and a `suiteguard_repocontract_test.go` that calls
`suiteguardcontract.CheckCurrentPackage`, and both join
`guardedSpawningPackages` in `internal/suiteguardcontract/contract.go`.

The Makefile's `REPO_CONTRACT_TESTS` lists every top-level test declared in a
`repocontract` file — adding `TestEverySpawningPackageInstallsTheSuiteGuard`,
`TestGovernedSetCoversOwnedShippedTemplates`, `TestGovernedSetOnlyGrows`,
`TestCleanupHistoricalGrantEvidence` and `TestEveryBoundedPathIsGoverned` — and
`repo-test` runs them across `./...`. A plain test in
`internal/suiteguardcontract/repository_gate_test.go` walks the module as Go
does (skipping directories named `testdata` or starting with `.` or `_`),
collects every top-level `Test` function in a `_test.go` file whose
`//go:build` line requires `repocontract`, and fails naming each one missing
from `REPO_CONTRACT_TESTS`; it also fails when `repo-test` stops passing
`-tags repocontract` or `./...`, or `verify-docs` stops depending on
`repo-test`. `make verify` runs it; CI already runs `make verify-docs`.

## API Contracts

1. `SchemaVersionError.Error()` names `roundfix migrate` for an older Run
   Database and `roundfix upgrade` for a newer one; its fields and type are
   unchanged.
2. `store.Open` returns `SchemaVersionError` for a newer Run Database.
3. `store.Migrate(ctx, homeDir)` reports the path and the schema version before
   and after, creates nothing for an absent database and writes nothing for a
   current or newer one.
4. `roundfix migrate` prints one stdout line and exits `0` for a migrated,
   current or absent Run Database, and exits `2` for a newer one.
5. `commandWantsHelp` returns true for `help` only as the first argument, and
   for `-h` or `--help` only before `--`.
6. A QA Report whose front matter the derived QA Verification's reader refuses
   is a `QAReportError`, and the QA Task it settles carries reason `QA verdict
   unreadable: <cause>`.
7. `GovernedPath("skills/_ownership.yml")` is true.
8. `make repo-test` runs every top-level test declared in a `repocontract` file.

## Coverage Map

- Goal 1 → The schema refusal and the migrate command; API Contracts 1-4.
- Goal 2 → The schema refusal and the migrate command.
- Goal 3 → The help predicate; API Contract 5.
- Goal 4 → The QA Report front matter; API Contract 6.
- Goal 5 → The governed set and the repository gate; API Contracts 7-8.
- Core Feature 1 → The schema refusal and the migrate command.
- Core Feature 2 → The help predicate.
- Core Feature 3 → The QA Report front matter.
- Core Feature 4 → The governed set and the repository gate.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → Testing Approach 3.
- Success Metric 4 → Testing Approach 4.
- API Contracts 1-4 → The schema refusal and the migrate command.
- API Contract 5 → The help predicate.
- API Contract 6 → The QA Report front matter.
- API Contracts 7-8 → The governed set and the repository gate.

## Integration Points

- **Spec 0172.** Made the seeded report `pending` and refused a hollow pass in
  `QAReportEligibility`; this Spec changes only how the report's front matter is
  bounded, below that decision.
- **Spec 0163.** Bounded `skills/_ownership.yml`, the path the governed set now
  covers.
- **Spec 0170.** Its second QA Report is the empty front matter this Spec
  refuses.

## Testing Approach

1. **Schema refusal and migrate.** Store tests build an older fixture (a fresh
   database downgraded as `seedOutdatedV9RunDatabase` does) and a newer one
   (`user_version` advanced), and assert each refusal's text; concurrent
   `store.Open` calls on an older fixture all succeed and leave it at
   `schemaVersion`; `store.Migrate` on an absent home creates nothing. CLI tests
   drive `runs list`, `migrate`, `runs list` on an older database, and `runs
   list` plus `migrate` on a newer one with the file's bytes compared before and
   after. `TestJournalConsumerCorpusReplaysEveryConsumer` stays green.
2. **Help predicate.** CLI tests assert exit codes and stdout for a flag value
   `help`, a trailing positional `help`, a leading `help`, a trailing `--help`
   and a `--help` after `--`, each as its own test.
3. **Front matter.** A table of fixtures and every archived QA Report are read
   by `readQAReport` and by the rendered derived Verification command in a
   temporary repository, and the two readability answers are asserted equal;
   named tests pin each refusal message; a Daemon test settles a report with an
   empty front matter; `qa-report accept` and `archive` refuse one;
   `TestArchivedPassCorpusRemainsArchiveEligible` stays green unchanged.
4. **Governed set and gate.** `TestGovernedPath` value cases for
   `skills/_ownership.yml` and an ordinary sibling; the `repocontract` tests
   named above pass; the gate audit fails on a fixture module whose
   `repocontract` test is missing from the list and passes on the repository.
5. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. The schema refusal and the migrate command (depends on: none).
2. The help predicate (depends on: 1).
3. The QA Report front matter (depends on: 2).
4. The governed set and the repository gate (depends on: none).
5. Terminal QA (depends on: 1, 2, 3, 4).

## Risks & Considerations

- **Shared CLI and guide files.** Tasks 1 and 2 edit `internal/cli/cli.go` and
  `internal/cli/cli_test.go`, and Tasks 1, 2 and 3 edit
  `docs/user-guide/commands.md`, so they form a chain; Task 4 shares no file
  and runs beside them.
- **An older binary's Active Run.** Migrating while a Run started by an older
  binary is Active is the condition every operational command already creates;
  the `roundfix migrate` guide says so rather than hiding it.
- **Guarding two more packages.** The suite guard fails a package whose tests
  write into the repository; if `internal/authorization` or
  `internal/verifyselect` does, the Task fixes the test, never the guard.
- **Top-level test names.** No Task renames or removes an existing top-level
  test, so `docs/references/coverage-record.json` stays unchanged.
