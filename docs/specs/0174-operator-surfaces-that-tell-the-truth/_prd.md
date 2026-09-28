---
spec: 0174-operator-surfaces-that-tell-the-truth
status: active
created: 2026-09-28
surfaces: [backend, cli, docs]
---

# Operator surfaces that tell the truth

The commands and checks an operator reads must say what is true. Four defects
make them say something else:

- A command that only reads the Run Database (`runs list`, `events`,
  `attach`, `doctor`, `gc`, `storage report`, `reconcile`, `deliver status`,
  `spec audit`, and the Run Database checks of `settle` and `reopen`) refuses
  a Run Database whose schema version differs from the binary's and tells the
  operator to run an operational command to migrate it. No migrate command exists, an
  operational command needs a Pull Request or a Spec, and the same advice is
  printed when the database is newer than the binary, where no migration
  exists. On 2026-09-24 that stopped an operator's monitoring for 13 hours; on
  2026-09-25 the Run Database was at schema 18 while the binary supported 20.
- A bare `help` anywhere in the arguments prints usage and exits `0`:
  `roundfix archive --bogus help` exits `0` where the same line without `help`
  exits `2`, and `roundfix implement --spec help` prints usage, so no flag can
  take the value `help`.
- A QA Report whose front matter is empty settles the QA Task as passed. Spec
  0170's second report opened with two consecutive `---` lines, so the derived
  QA Verification's reader refused it, yet the Daemon's settlement read
  `verdict: pass` from the lines below and committed the report as `(pass)`.
- `GovernedPath` does not match `skills/_ownership.yml`, which archived Spec
  0163 bounds, and the contract test that says so runs in no repository gate,
  so the failure is silent. The same gap hides a second failing contract: two
  packages spawn processes without installing the suite guard.

## Project Constraints

- Identifier strategy: not applicable — no new identifier; the Run Database
  schema version stays the integer SQLite `user_version`. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files, Git and the Run
  Database only; no credential and no network call. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0004 keeps one central Run
  Database; ADR-0009 lets a reader hold a read-only connection beside the
  writer, which this Spec keeps read-only; ADR-0133 makes a diagnostic name the
  literal it requires, so each refusal names the command or binary that
  resolves it; ADR-0015 reads the verdict from the QA Report front matter and
  counts an unreadable one as a failure; ADR-0080 owns verdict semantics;
  ADR-0091 keeps the gate a Task node; ADR-0096 says no stage may make a
  verdict more permissive; ADR-0014 and ADR-0057 keep Verification and Task
  status Daemon-owned; ADR-0130 holds the governed set to every path an
  authorization has bounded; ADR-0126 makes a repository-contract test
  enumerate the packages that install the suite guard; ADR-0081 and ADR-0149
  own the sanctioned regeneration this Spec declares. ADR-0022 cites ADR-0004
  but governs Stop Requests, and ADR-0020, ADR-0038, ADR-0127 and ADR-0160 cite
  ADR-0014 but govern the Agent prompt result, the Verification repair turn,
  process residue and entry to a red repository gate, none of which this Spec
  touches, so they do not apply; nor do ADR-0056 and ADR-0159, which cite
  ADR-0038 for Verification capacity and independent Verification. ADR-0093
  checks Spec
  consistency by citation, ADR-0104 accepts on evidence a Spec did not author,
  ADR-0155 makes the `qa` Task declare the matrix and ADR-0156 makes a declared
  promise name a consuming Task. This Spec's gate is bound by ADR-0080,
  ADR-0091, ADR-0096, ADR-0097 and ADR-0117. All hold. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — the maintainer approved Onda 2 of the
  efficiency sequence in chat on 2026-09-28; the Go sources and tests ride the
  standing grant of 2026-09-21 for governed source, and the skill files ride
  the standing grant of 2026-09-18 for keeping the shipped skills true to the
  CLI, recorded in [_authorization.md](_authorization.md); bounded files:
  `internal/cli/cli_test.go`, `.agents/skills/roundfix/SKILL.md`,
  `skills/roundfix/SKILL.md`, `.agents/skills/qa-gate/SKILL.md`,
  `skills/qa-gate/SKILL.md`, `internal/speccheck/governed.go`,
  `internal/speccheck/governed_repocontract_test.go`, `Makefile`. Sanctioned
  regeneration: `make skills-sync`, `make baseline-digests`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- An operator whose binary and Run Database disagree on the schema learns the
  one action that resolves it, and can take it without opening a Pull Request.
- Read-only commands never write the Run Database.
- A help token means help only where it stands as a command's argument.
- A QA Report the derived QA Verification cannot read never settles a gate.
- Every path an authorization bounds is governed, and the repository-contract
  tests that prove it run in a repository gate.

## Core Features

1. **A direction-aware refusal and a migrate command.** A read-only command
   that meets an older Run Database refuses with an instruction to run
   `roundfix migrate`; one that meets a newer Run Database refuses with an
   instruction to use the newer binary or `roundfix upgrade`. The new
   `roundfix migrate` command upgrades an older Run Database, reports an
   up-to-date or absent one without writing, and refuses a newer one without
   writing. The migration holds the machine-wide write lock from reading the
   version to recording the new one, so two processes migrating at once
   upgrade the database exactly once.
2. **Help only where it is an argument.** A bare `help` requests usage only as
   the first argument a command receives; `-h` and `--help` request usage
   wherever they stand before a `--` terminator; nothing after `--` is ever a
   help token.
3. **Settlement reads the front matter the derived Verification reads.** The QA
   Report reader bounds its front matter exactly as the derived QA
   Verification's reader does, so a report whose front matter is empty, is not
   closed, or carries a duplicated `verdict:` line is unreadable to QA
   settlement, `archive`, `settle` and `qa-report accept` alike, and the QA
   Task's reason names why.
4. **The governed set covers every bounded file and its contract runs in a
   gate.** `GovernedPath` matches `skills/_ownership.yml`; the two packages that
   spawn processes install the suite guard; `make verify-docs` runs every
   `repocontract` test in the repository, and a test run by `make verify`
   fails when a `repocontract` test is left out of that gate.

## Non-Goals / Out of Scope

- Reading an older or newer Run Database schema in place, or migrating from a
  read-only command.
- Downgrading a Run Database, or supporting a schema version the binary does
  not know.
- Guarding `roundfix migrate` against an Active Run started by an older binary;
  every operational command already migrates in that condition, and this Spec
  neither adds nor removes that risk.
- Letting flag parsing own `-h` and `--help` per command; the help rule stays
  in the one shared predicate.
- Changing the derived QA Verification command the `qa` Task renders, or the
  verdicts a readable report may carry.
- Recognizing an `unreadable` or `missing` QA commit as a QA-report-only
  commit.
- Changing the CI workflow; `make verify-docs` already runs in it.

## Success Metrics

1. Against an older Run Database, `roundfix runs list` exits non-zero with a
   stderr diagnostic containing `run 'roundfix migrate'`; `roundfix migrate`
   then exits `0` and reports the from and to schema versions; `roundfix runs
   list` then exits `0`. Against a newer Run Database, `runs list` names
   `roundfix upgrade`, `roundfix migrate` exits `2`, and the database file is
   byte-identical afterwards.
2. `roundfix archive --bogus help` exits `2` and `roundfix implement --spec
   help` prints no usage, while `roundfix archive help` and `roundfix archive
   <slug> --help` print usage and exit `0`.
3. A QA Report opening with two consecutive `---` lines settles the QA Task
   `failed` with verdict `unreadable` and a reason naming the empty front
   matter, and is refused by `qa-report accept` and `archive`; every archived QA
   Report in `docs/history/specs/` reads the same way through the Go reader and
   the derived Verification's reader.
4. `go test -tags repocontract -run TestEveryBoundedPathIsGoverned
   ./internal/speccheck` and `TestEverySpawningPackageInstallsTheSuiteGuard` in
   `./internal/suiteguard` pass, `make repo-test` runs every `repocontract`
   test in the repository, and removing one of them from the gate's list fails
   `go test ./internal/suiteguardcontract`.

## Recorded limits

- The QA Report reader stays stricter than the derived Verification's reader
  where it already was: a readable front matter whose verdict is unsupported,
  or whose blocked-row counts contradict the verdict, is still refused. The
  two readers agree on which front matter is readable, not on every verdict.
- A second front-matter-shaped block in the report body is body text to both
  readers; only the first block counts.

## Decisions

- **Refuse with an exact instruction; migrate only on request.** A read-only
  command keeps opening a read-only connection. Migrating from `runs list`
  would change the Run Database under an older binary's Active Run without
  anyone asking, and the storage reader behind `gc` and `storage report` must
  leave the database byte-identical. Reading older schemas in place would make
  every query schema-aware for a condition an explicit command resolves in one
  step. The refusal names that step, and `roundfix migrate` performs it.
- **One migration path, serialized.** `roundfix migrate` and every operational
  command's writer open share one migration, and it reads the version, derives
  its statements and records the new version under one hold of the
  machine-wide write lock. Today the version and the statements are read
  before the lock, so a second process can replay a migration the first one
  already applied.
- **A newer database names the binary, not a migration.** No binary can
  migrate a schema it does not know, so the refusal names the newer binary and
  `roundfix upgrade`; the writer open gives the same typed refusal instead of
  "schema version N is not supported".
- **Fix the predicate, not its 47 callers.** `commandWantsHelp` is the one
  place every command asks; a flag's value is never the first argument, so a
  leading-only bare `help` cannot collide with a value, and a value spelled
  `-h` or `--help` is never a valid slug, ID or path.
- **One reader, the Verification's bounds.** The QA Report reader is shared by
  settlement, `archive`, `settle` and `qa-report accept`, so bounding its front
  matter as the derived Verification does closes every door at once. A test
  runs both readers on the same corpus so the two cannot drift again.
- **Name the gate's tests, and guard the list.** Running `go test -tags
  repocontract ./...` in full re-runs the whole suite, measured at 3m44s on
  2026-09-28; running the named `repocontract` tests across `./...` took 8s. A
  plain test fails when a `repocontract` test is missing from the Makefile's
  list, so a new contract cannot be silently left out of the gate.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in the
Task Graph. The negative cases carry the weight: a read-only command that
writes, a newer database told to migrate, a migration replayed by a second
process, a flag value read as help, an empty front matter that settles, or a
contract test outside every gate would each pass a happy-path test.

The outside-evidence rows rest on artifacts this Spec did not write:

- Every archived QA Report in `docs/history/specs/` (265 on 2026-09-28) and
  every newest QA Report in the Fiscus mirror at
  `/Users/marcio/dev/secondbrain/projects/fiscus/mirror/docs/history/specs/`
  reads the same way through the Go reader and the derived Verification's
  reader, recorded as evidence this Spec did not author, or the Fiscus half as
  blocked with its reason when the mirror is absent.
- A Run Database created by the installed roundfix 0.15.0 binary at
  `/Users/marcio/go/bin/roundfix` in a disposable home is refused by the built
  `runs list` with the `roundfix migrate` instruction and upgraded by the built
  `roundfix migrate`, or the row is blocked with its reason when that binary is
  absent. The live Run Database under `~/.roundfix` is never opened.
- Every `_authorization.md` under `docs/history/specs/`, written by earlier
  Specs, bounds only governed paths.

## Research basis

The schema refusal was observed in the 2026-09-24 and 2026-09-25 sessions: `runs
list` answered `has schema version 13, but this binary supports schema version
18`, and the installed 0.15.0 binary refused a v18 database. The help-token
defect was captured in the secondbrain inbox
(`inbox/roundfix/2026-09-19-help-token-accepted-anywhere.md`). The empty front
matter is Spec 0170's `qa/qa-report-2026-09-25-01.md` at commit `e3042b40` in Run
`run_20260925T195747Z_0ccff24444c45b55`. The governed-path failure reproduces on
`0160f70a`: `go test -tags repocontract ./...` fails
`TestEveryBoundedPathIsGoverned` on `skills/_ownership.yml` and
`TestEverySpawningPackageInstallsTheSuiteGuard` on `internal/authorization` and
`internal/verifyselect`. The adopted sources are indexed in
[references/_index.md](references/_index.md).

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
