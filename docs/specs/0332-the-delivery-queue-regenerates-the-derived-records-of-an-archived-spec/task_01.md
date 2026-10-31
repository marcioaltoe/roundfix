---
task: task_01
spec: 0332-the-delivery-queue-regenerates-the-derived-records-of-an-archived-spec
status: completed
type: backend
complexity: medium
---

# Task 01: The coverage record writer names its command, and the regeneration contract models an archived item

## Overview

`TestCoverageEquivalence` rewrites `docs/references/coverage-record.json`
under `-update-coverage-record` but never names its command to the suite
guard. The guard therefore refuses the write even when an Archive Record or
an authorization record sanctions it. This Task makes the writer name its
command. It then turns `TestRegenerationIsDeclared` into the class contract
for every derived record the Delivery Queue regenerates:

- the contract's repository copy is sanctioned only by an archived item's
  Archive Record;
- its coverage record is stale;
- it no longer injects a declaration;
- a second copy proves that an unsanctioned write is still refused.

It also records the rule and the glossary revision. This answers the
adopted Backlog Entry of 2026-10-27,
[the queue cannot regenerate the coverage record of an archived Spec](references/2026-10-27-the-queue-cannot-regenerate-the-coverage-record-of-an-archived-spec.md).

## Requirements

1. MUST add to `internal/spec/coverage_test.go` the package constant
   `coverageRecordCommand` with the value
   `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`.
   The value must equal, byte for byte, the `regenerate` string that
   `.roundfixrc.yml` declares for `docs/references/coverage-record.json`.
2. When `-update-coverage-record` is set, `TestCoverageEquivalence` MUST call
   `suiteguard.DeclareSanctionedRegeneration(coverageRecordCommand)` before it
   writes the record. Without the flag it MUST declare nothing
   (`_techspec.md` → API Contract 1). The record's format and the comparison
   MUST keep their behavior.
3. MUST change the fixture of `TestRegenerationIsDeclared`
   (`internal/config/regeneration_declared_test.go`) as API Contract 2
   states.
   - In the copy, remove every `docs/history/specs/*.md` record, the
     `docs/specs/` tree and `docs/workflow/authorizations/`. The copy of a Run
     worktree holds this Spec's own `_authorization.md`, which must not
     sanction the coverage command there.
   - Write the archived item's record `docs/history/specs/0001-the-archived-item.md`
     with `schema: roundfix/archive-record/v1` and the `regeneration` list of
     `_techspec.md` → Data Models.
   - Replace one entry of the copied coverage record with
     `TestRegenerationFixtureStaleEntry`, keeping it valid JSON.
   - Keep the module edit, then commit and run every declaration as today.
   - Assert three things: every command exits 0, the changed paths include
     `docs/references/coverage-record.json`, and the regenerated record no
     longer contains `TestRegenerationFixtureStaleEntry`.
   - Remove `addFixtureRegenerationAuthorization`, the legacy authorization
     it wrote and the `internal/spec/regeneration_fixture_test.go` it injected.
     The file MUST no longer contain `DeclareSanctionedRegeneration` or
     `regeneration_fixture_test.go`.
   - The three ablation subtests keep their behavior.
4. MUST add the subtest `an_unsanctioned_coverage_write_is_refused` to
   `TestRegenerationIsDeclared` (API Contract 3). It builds a second copy the
   same way, but the record lists the module record command and
   `make baseline-digests` and not the coverage command. It runs only the
   coverage declaration and asserts two things: the command fails, and its
   output contains `suiteguard: modified: docs/references/coverage-record.json`.
   Copy the repository in-process with the file's existing `copyRepository`,
   so the Go test cache records the copied sources. Add no new top-level test.
5. MUST add one bullet to the repository rules of
   `docs/agents/specific-repository.md`. The bullet says that a test which
   rewrites a derived record under its record flag calls
   `suiteguard.DeclareSanctionedRegeneration` with the exact command that runs
   it, only while that flag is set and before it writes. It says that the
   suite guard then exempts only the outputs that an authorization record or
   an Archive Record in the tree lists for that command. It cites ADR-0345.
   Keep every existing line of that file.
6. MUST, through `domain-modeling`, revise **Archive Record** in `CONTEXT.md`.
   The revision says that the record carries the `regeneration` commands its
   authorization record sanctioned, with their outputs. It says that the
   suite guard keeps honoring them after the Spec folder leaves the tree, so
   a Delivery Queue regeneration after archive is sanctioned by the item's own
   record. It cites ADR-0345 and no Spec, and it keeps the entry's existing
   sentences.
7. MUST record the Coverage Review for the `config keys` Behavior Surface
   that the PRD's `## Skills` section declares unchanged. Replace that
   entry's `review` in `docs/references/skill-coverage.json` with a dated
   explanation beginning `2026-10-30 — `. The explanation says that only the
   regeneration contract test in `internal/config` changed and that no
   configuration key, default, validation or message changed. Keep every
   other byte of the map.
8. MUST NOT change `internal/suiteguard/`, `internal/suiteguardcontract/`,
   `.roundfixrc.yml`, `internal/cli/deliver_workflow.go` or the Makefile.

## Subtasks

- [ ] Name the command in the coverage writer.
- [ ] Rewrite the regeneration contract's fixture and add the unsanctioned subtest.
- [ ] Add the repository rule, revise the glossary term and record the Coverage Review.

## Acceptance Criteria

- [ ] A coverage re-record exits 0 when a record in the tree lists the command with the record as output, and exits non-zero naming the record when none does.
- [ ] Every `.roundfixrc.yml` regeneration passes in a copy sanctioned only by the archived item's record, with no injected declaration.
- [ ] `docs/agents/specific-repository.md` and **Archive Record** cite ADR-0345.

## Context

- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `docs/adr/0345-a-derived-record-writer-names-its-command-to-the-suite-guard.md`
- instruction: `.roundfixrc.yml`
- instruction: `internal/suiteguard/suiteguard.go`
- instruction: `internal/suiteguardcontract/regeneration.go`
- interface: `internal/spec/coverage_test.go`
- interface: `internal/config/regeneration_declared_test.go`
- interface: `docs/agents/specific-repository.md`
- interface: `CONTEXT.md`
- interface: `docs/references/skill-coverage.json`

## Consumers

- The coverage record writer's command declaration — consumers: `./internal/spec`, `./internal/config`
- The Skill Coverage Map's `config keys` review — consumers: `./internal/docscontract`

## Verification

- `grep -q 'DeclareSanctionedRegeneration(coverageRecordCommand)' internal/spec/coverage_test.go || { echo 'the coverage writer does not name its command' >&2; exit 1; }; grep -q 'an_unsanctioned_coverage_write_is_refused' internal/config/regeneration_declared_test.go || { echo 'the unsanctioned subtest is missing' >&2; exit 1; }; if grep -q 'DeclareSanctionedRegeneration' internal/config/regeneration_declared_test.go; then echo 'the fixture still injects a declaration' >&2; exit 1; fi; if grep -q 'regeneration_fixture_test.go' internal/config/regeneration_declared_test.go; then echo 'the fixture still writes the injected file' >&2; exit 1; fi; tr -s '[:space:]' ' ' < docs/agents/specific-repository.md | grep -qF 'suiteguard.DeclareSanctionedRegeneration' || { echo 'the repository rule is missing' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/agents/specific-repository.md | grep -qF 'ADR-0345' || { echo 'the repository rule does not cite ADR-0345' >&2; exit 1; }; awk 'index($0, "**Archive Record**:") == 1 { f=1; next } f && index($0, "_Avoid_") == 1 { exit } f { print }' CONTEXT.md | tr -s '[:space:]' ' ' | grep -qF 'ADR-0345' || { echo 'CONTEXT.md **Archive Record** does not cite ADR-0345' >&2; exit 1; }; grep -qF '"review": "2026-10-30 ' docs/references/skill-coverage.json || { echo 'the config keys Coverage Review is not recorded' >&2; exit 1; }; go test -count=1 -timeout 20m ./internal/spec -run '^(TestCoverageEquivalence|TestCompareCoverageRecordsReportsMissingTest|TestCompareCoverageRecordsReportsAddedTestWithoutRegression)$' && go test -count=1 -timeout 30m -tags docscontract ./internal/config -run '^TestRegenerationIsDeclared$' -v && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$'` — expected: exit 0. Before this Task the coverage writer names no command, so the first check fails. Afterwards the writer declares its command, and the fixture neither injects a declaration nor writes the injected file. The rule and **Archive Record** cite ADR-0345, the `config keys` Coverage Review is dated, and the Skill Coverage Map contract still passes. The coverage comparison tests pass, and so does `TestRegenerationIsDeclared` with its subtest `an_unsanctioned_coverage_write_is_refused`, in which every regeneration runs under the archived item's record alone.

## References

- `_prd.md` → Goals; Core Features 1, 2, 4; Skills; Success Metric 1; Success Metric 2; Success Metric 4; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-3; Data Models; Invariants 1-3, 5; Build Order 1
- ADR-0345, ADR-0149


## Result

Implemented this Task's slice. Task status remains Daemon-owned.

- Acceptance criterion 1: the coverage writer now declares
  `coverageRecordCommand` only inside the update-flag branch and before writing.
  A focused run of the regeneration contract passed both the sanctioned parent
  fixture and `an_unsanctioned_coverage_write_is_refused`. The latter requires a
  non-zero coverage-command exit and the exact refused-record diagnostic.
- Acceptance criterion 2: both fixtures use the existing in-process repository
  copy, remove active Specs, legacy authorizations and existing Archive Record
  files, and write the archived item's record with the required regeneration
  list. Both retain the module edit and commit a stale, valid coverage record.
  The positive fixture runs all seven declarations, requires every command to
  exit zero, requires the coverage record among changed paths, and rejects a
  remaining stale entry. No declaration is injected; the three ablation bodies
  retain their bytes.
- Acceptance criterion 3: the added repository-rule bullet and the expanded
  **Archive Record** definition both link ADR-0345. Existing repository-rule
  lines and glossary sentences were retained. The `config keys` Coverage Review
  is dated `2026-10-30`; every other byte of the map was retained.

Focused checks:

- Before changing the writer, the revised fixture reproduced the nested coverage
  command's exit 1 and `suiteguard: modified: docs/references/coverage-record.json`.
  The outer guard also reported documentation edited during that initial run;
  the subsequent passing run used the stable, final source and documentation.
- `go test -tags docscontract ./internal/config -run '^TestRegenerationIsDeclared$/^an_unsanctioned_coverage_write_is_refused$' -count=1 -timeout 15m -v`
  exited 0 (`roundfix/internal/config`, 56.666 s). The parent exercised all
  regenerations; the selected unsanctioned subtest passed. The three ablation
  subtests were outside this focused selection.
- `go test ./internal/spec -run '^(TestMarshalCoverageRecordIsDeterministic|TestCoverageRecordCountsNoPackageUnderDocs)$' -count=1 -timeout 5m -v`
  exited 0 (`roundfix/internal/spec`, 1.155 s).
- `go test -tags docscontract ./internal/docscontract -run '^TestGlossary(PhrasesAreDefinedTerms|LinksNameAcceptedADRs|AvoidLinesArePlainLists|StatesEachRuleOnce)$' -count=1 -timeout 5m -v`
  exited 0 (`roundfix/internal/docscontract`, 0.418 s).
- Python inspections confirmed byte equality between the command constant and
  the Project Config declaration, flag/declaration/write ordering, unchanged
  ablation bodies, no injected declaration or source file, exactly one top-level
  config contract test, the bounded documentation edits, and untouched protected
  paths. `gofmt` was applied to both Go files; `git diff --check` exited 0.

The Task's declared Verification and both repository Verification commands were
not run; they remain Daemon-owned. The coverage record re-record belongs to
`task_02` and was not run in this worktree. No commit, push, Pull Request, other
Task-file edit or Task Graph edit was made.
