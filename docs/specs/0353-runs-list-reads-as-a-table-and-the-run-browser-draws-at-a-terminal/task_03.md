---
task: task_03
spec: 0353-runs-list-reads-as-a-table-and-the-run-browser-draws-at-a-terminal
status: completed
type: backend
complexity: medium
---

# Task 03: attach and runs show accept the short Run id, and ADR-0363 records the two formats

## Overview

The table task_02 adds shows an eight-digit short id, but every command looks
a Run up only by its exact id (`_techspec.md` → Current behavior 7). After this
Task the store resolves a reference of eight to sixteen lowercase hexadecimal
digits to the one Run whose random suffix starts with it, `attach` and
`runs show` use that resolution, an ambiguous prefix exits `2` listing every
match, and ADR-0363 records the dual text surface and the short id. This makes
the id shown to the maintainer usable, as the report reported by the
maintainer on 2026-11-05 asked for "id curto"
([references](references/2026-11-05-macbook-runs-list-is-confusing-and-attach-shows-nothing.md)).

## Requirements

1. MUST, before changing code, run
   `go test -count=1 -timeout 20m -v -run '^(TestRunsShow.*|TestAttachUnknownRunFailsBeforeTUIStart|TestAttachReplaysCompletedRunReadOnly)$' ./internal/cli`
   and record in the Result that it passes, with `runs show run_missing`
   refused as an unknown Run.
2. MUST add `AmbiguousRunReferenceError` and
   `(*Store).ResolveRunReference` to `internal/store` exactly as
   `_techspec.md` → API Contract 9 states, with the error text
   `Run id prefix "<reference>" matches <n> Runs: <id>, <id>; pass more digits or the full Run id`
   and the ids in descending id order, which is newest first; `(*Store).Run` keeps its exact-id behavior.
3. MUST make `runAttachCommand` (explicit id) and `runRunsShowCommand` call
   `ResolveRunReference` instead of `Run` (API Contract 10): an ambiguity
   prints `roundfix attach failed: <error>`, or
   `roundfix: runs show failed: <error>` and its usage line, and exits `2`;
   not found keeps each command's message; `runs show` reads token usage
   with the resolved `run.ID`. `events`, `stop`, `reconcile`, the Run
   Browser loop and every other caller of `Run` stay unchanged.
4. MUST add, in `commandUsage`, the line
   `<run-id> is a full Run ID or the eight-digit short id that runs list shows.`
   after the description paragraph of `attach` and after
   `Prints recorded token usage per Work Item scope and in total, read-only.`
   of `runs show`.
5. MUST create `internal/store/run_reference_test.go` with
   `TestResolveRunReferenceNamesOneRunByItsShortID`, calling `t.Parallel()`
   first, over Runs created with `CreateRun` in distinct temporary
   repositories and renamed through a raw database connection opened with
   `_pragma=foreign_keys(0)` to `run_20260706T120000Z_93cee2a36ff33269`,
   `run_20260706T120100Z_93cee2a3aaaaaaaa` and
   `run_20260706T120200Z_0123456789abcdef`, with subtests:
   `exact` (the full id gives that Run), `unique prefix`
   (`0123456789abcdef` and `01234567` give the third Run), `ambiguous`
   (`93cee2a3` gives `*AmbiguousRunReferenceError` with both ids, the newest
   first, and the exact text), `longer prefix` (`93cee2a36` gives the first
   Run), `too short` (`0123456` gives not found), `not hex` (`0123456z` and
   `0123456789ABCDEF` give not found), `no match` (`ffffffff` gives not
   found) and `empty` (the `Run ID is required` error).
6. MUST create `internal/cli/run_short_id_test.go`, each test calling
   `t.Parallel()` first, with:
   - `TestAttachAndRunsShowAcceptTheShortID`: `attach <ShortRunID>` of a
     completed seeded Run (as `TestAttachReplaysCompletedRunReadOnly` seeds
     one) prints the same stdout as `attach <full id>`, and
     `runs show <ShortRunID>` of a usage Run (`seedUsageRun`) prints the same
     stdout as `runs show <full id>`;
   - `TestAnAmbiguousShortIDListsEveryMatch`: two Runs seeded with
     `seedRunsForList` and renamed through the raw Run Database (opened as
     `setListedRunCreatedAt` opens it, but with `_pragma=foreign_keys(0)`) to
     `run_20260706T120000Z_93cee2a3cccccccc` and
     `run_20260706T120100Z_93cee2a3bbbbbbbb`; `attach 93cee2a3` exits `2`
     with stderr exactly
     `roundfix attach failed: Run id prefix "93cee2a3" matches 2 Runs: run_20260706T120100Z_93cee2a3bbbbbbbb, run_20260706T120000Z_93cee2a3cccccccc; pass more digits or the full Run id\n`,
     and `runs show 93cee2a3` exits `2` with that error after
     `roundfix: runs show failed: ` and the usage line.
7. MUST create
   `docs/adr/0363-runs-list-prints-a-table-at-a-terminal-and-a-short-run-id-names-a-run.md`
   with the content of `_techspec.md` → ADR-0363, apart from line wrapping.
8. MUST add to `.agents/skills/roundfix/references/runs.md` after the
   paragraph that ends "pass a run id or run `roundfix attach` to pick
   interactively.", verbatim apart from line wrapping:
   "`attach` and `runs show` also accept the short id that the `runs list`
   table shows, or any longer prefix of a Run id's random suffix up to its
   sixteen digits, when it names exactly one Run. A prefix that matches
   several Runs exits `2` and lists them:
   `Run id prefix "<prefix>" matches <n> Runs: <id>, <id>; pass more digits or the full Run id`.
   `events`, `stop`, `reconcile` and every other command need the full Run
   id." Run `make skills-sync`, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
9. MUST add to `docs/user-guide/commands/attach.md` after its usage block,
   and to `docs/user-guide/commands/runs.md` after the paragraph that starts
   "`roundfix runs show <run-id>` reads the Run Database", verbatim apart from
   line wrapping: "`<run-id>` may also be the eight-digit short id that the
   `runs list` table shows, or a longer prefix of the id's random suffix, when
   it names one Run; a prefix that names several Runs exits `2` and lists
   them ([ADR-0363](../../adr/0363-runs-list-prints-a-table-at-a-terminal-and-a-short-run-id-names-a-run.md))."
   Then, as the LAST step after the last guide edit, re-record the Behavior
   Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
10. MUST run sabotages 9 and 10 of `_techspec.md` → Testing Approach on a
    temporary copy and record which named test fails.
11. MUST keep every existing test of `./internal/store` and every
    `TestRunsShow*` and `TestAttach*` test of `./internal/cli` passing
    unchanged.
12. MUST NOT change the Run Database schema, `newRunID`, `(*Store).Run`,
    `events`, `stop`, `reconcile`, the Run Browser, `runs list`,
    `internal/config`, `internal/cli/cli_test.go`, any other skill or reference, the Makefile, CI
    workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] Characterize the exact-id lookup of attach and runs show.
- [ ] Resolve a short id in the store, with the ambiguity error.
- [ ] Use it in attach and runs show; update their help text.
- [ ] Add the store and command tests; run the sabotages.
- [ ] Write ADR-0363; document the short id; sync and record the skill; re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] `attach` and `runs show` act on the one Run a short id or longer suffix prefix names, and an exact id still works.
- [ ] An ambiguous prefix exits `2` listing every match; a prefix of fewer than eight digits, an uppercase or non-hex reference, or no match keeps today's unknown-Run message.
- [ ] ADR-0363 exists with status `accepted`, and the skill and guides state the short id.

## Context

- interface: `internal/store/store.go`
- creates: `internal/store/run_reference_test.go`
- interface: `internal/cli/attach.go`
- interface: `internal/cli/runs_show.go`
- interface: `internal/cli/cli.go`
- creates: `internal/cli/run_short_id_test.go`
- interface: `internal/cli/runs_show_test.go`
- creates: `docs/adr/0363-runs-list-prints-a-table-at-a-terminal-and-a-short-run-id-names-a-run.md`
- interface: `.agents/skills/roundfix/references/runs.md`
- interface: `skills/roundfix/references/runs.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/attach.md`
- interface: `docs/user-guide/commands/runs.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- `(*store.Store).ResolveRunReference` and `AmbiguousRunReferenceError` — consumers: `./internal/store`, `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestResolveRunReferenceNamesOneRunByItsShortID|TestListRuns.*)$' ./internal/store 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestResolveRunReferenceNamesOneRunByItsShortID ' || { echo 'missing pass: TestResolveRunReferenceNamesOneRunByItsShortID' >&2; exit 1; }; out="$(go test -count=1 -timeout 20m -v -run '^(TestAttachAndRunsShowAcceptTheShortID|TestAnAmbiguousShortIDListsEveryMatch|TestRunsShow.*|TestAttach.*|TestRunCommandHelp)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestAttachAndRunsShowAcceptTheShortID TestAnAmbiguousShortIDListsEveryMatch; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; test -f docs/adr/0363-runs-list-prints-a-table-at-a-terminal-and-a-short-run-id-names-a-run.md && grep -q '^status: accepted$' docs/adr/0363-runs-list-prints-a-table-at-a-terminal-and-a-short-run-id-names-a-run.md || { echo 'missing accepted ADR-0363' >&2; exit 1; }; for f in .agents/skills/roundfix/references/runs.md skills/roundfix/references/runs.md docs/user-guide/commands/runs.md docs/user-guide/commands/attach.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'prefix that matches several Runs' || tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'prefix that names several Runs' || { echo "missing short id text in $f" >&2; exit 1; }; done; cmp -s .agents/skills/roundfix/references/runs.md skills/roundfix/references/runs.md || { echo 'runs reference mirror differs' >&2; exit 1; }; go test -count=1 -timeout 20m ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$'` — expected: exit 0. Before this Task none of the three new tests exists, so the first pass line is missing; afterwards a short id names one Run in the store, `attach` and `runs show`, an ambiguous one lists every match, ADR-0363 is accepted, the skill and guides state it with a current version record and Behavior Surface Record, and the existing store, attach and runs show tests pass.

## References

- `_prd.md` → Goals; User Story 3; Core Feature 4; Core Feature 7; Success Metric 4; Success Metric 7
- `_techspec.md` → Current behavior 7; API Contract 9; API Contract 10; Invariant 5; ADR-0363; Testing Approach; Build Order 3


## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, tech-writer, evidence-gate, domain-modeling.

Implemented this Task's short-reference lookup and documentation slice. The
Daemon retains Task status and declared Verification ownership.

### Acceptance evidence

- Unique and exact references: `ResolveRunReference` trims input, uses the
  unchanged exact-id `Run` lookup first, then resolves only lowercase hexadecimal
  suffix prefixes of eight to sixteen digits with the specified escaped `LIKE`
  query. Explicit-id `attach` and `runs show` use it; token usage uses the
  resolved `run.ID`. `TestResolveRunReferenceNamesOneRunByItsShortID` passes
  against a read-only reader over three Runs created in distinct temporary
  repositories and renamed through a raw connection with foreign keys disabled.
  Its exact, unique-prefix, longer-prefix and trimmed cases exercise resolution;
  its `Run_stays_exact` case proves `Run` still refuses a prefix.
- Ambiguity and unknown references: the store test checks the typed error,
  reference, both full ids in descending id order, and exact error text.
  `TestAnAmbiguousShortIDListsEveryMatch` passes with both prescribed renamed
  Runs, asserting exit `2`, empty stdout and each command's exact stderr,
  including the `runs show` usage line. Its invalid-reference case proves both
  commands retain their unknown-Run messages for too-short, uppercase,
  non-hexadecimal and unmatched input. Store cases additionally cover empty,
  whitespace-only and overlong input.
- Command behavior: `TestAttachAndRunsShowAcceptTheShortID` passes. A completed
  Run seeded through the existing attach fixture replays identical stdout for
  short and full ids without probing an Agent; a `seedUsageRun` fixture prints
  identical token usage in both text and JSON for short and full ids.
- Documentation: ADR-0363 is accepted and byte-identical to the TechSpec's ADR
  block. Both help descriptions contain the required short-id line. The
  canonical `runs` skill reference and both user guides contain the prescribed
  paragraphs, and both guide links resolve to ADR-0363. The skill reference and
  `SKILL.md` mirrors are byte-identical. The sanctioned version recorder raised
  both Roundfix version fields from `0.1.187` to `0.1.188` and recorded the new
  digest. The Behavior Surface Record was regenerated after the final guide edit;
  only the attach/show help and attach/runs guide fingerprints changed.

### Focused checks

- Before changing code:
  `go test -count=1 -timeout 20m -v -run '^(TestRunsShow.*|TestAttachUnknownRunFailsBeforeTUIStart|TestAttachReplaysCompletedRunReadOnly)$' ./internal/cli`
  — exit `0`; all eight selected tests passed. In particular,
  `TestRunsShowRefusesAnUnknownRun` refused `runs show run_missing` with exit
  `2`, empty stdout and
  `roundfix: runs show failed: Run "run_missing" does not exist`, followed by
  `Run 'roundfix runs show --help' for usage.`
- Test-first signal:
  `go test -count=1 -run '^TestResolveRunReferenceNamesOneRunByItsShortID$' ./internal/store`
  — exit `1` before implementation, because `ResolveRunReference` and
  `AmbiguousRunReferenceError` were undefined.
- `go test -count=1 -timeout 10m -v -run '^TestResolveRunReferenceNamesOneRunByItsShortID$' ./internal/store`
  — exit `0`; every resolution subtest passed.
- `go test -count=1 -timeout 10m -v -run '^(TestAttachAndRunsShowAcceptTheShortID|TestAnAmbiguousShortIDListsEveryMatch)$' ./internal/cli`
  — exit `0`; both command tests and all their subtests passed.
- `go test -count=1 -timeout 10m ./internal/store` — exit `0`; the entire
  affected store package passed, including its unchanged existing tests.
- `go test -count=1 -timeout 10m -run '^(TestAttach|TestRunsShow|TestRunCommandHelp)' ./internal/cli`
  — exit `0`; the existing attach/show/help tests passed unchanged.
- `make skills-sync`, followed by
  `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
  — both exited `0`; mirror content and both recorded version fields agree.
- `make baseline-digests` — exit `0`; reported
  `baseline-digests: no changes; derived artifacts already match their canonical sources`
  and `changed: false`; no Baseline artifact changed.
- As the last documentation regeneration step, after the last guide edit:
  `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
  — exit `0`; generated `docs/references/behavior-surfaces.json`.
- Artifact inspection confirmed the exact ADR content, matching skill mirrors,
  valid guide links and a clean `git diff --check`.

### Required sabotage probes

Each probe ran on its own disposable copy of the current checkout under
`/private/tmp/rf-task03-sabotage-t4ad3ukg`, with a separate local Git database
for suiteguard. Neither mutation was applied to the Run Worktree.

- Sabotage 9 changed the accepted prefix length from `{8,16}` to `{7,16}`.
  `go test -count=1 -timeout 10m -v -run '^TestResolveRunReferenceNamesOneRunByItsShortID$/^too_short$' ./internal/store`
  exited `1` with named failure
  `TestResolveRunReferenceNamesOneRunByItsShortID/too_short`:
  `reference="0123456" id="run_20260706T120200Z_0123456789abcdef" found=true err=<nil>; want id=""`.
- Sabotage 10 replaced the ambiguity return with the newest matching Run.
  `go test -count=1 -timeout 10m -v -run '^TestResolveRunReferenceNamesOneRunByItsShortID$/^ambiguous$' ./internal/store`
  exited `1` with named failure
  `TestResolveRunReferenceNamesOneRunByItsShortID/ambiguous`:
  `id="run_20260706T120100Z_93cee2a3aaaaaaaa" found=true err=<nil>; want ambiguous reference`.
- Both failing transcripts were inspected: each failed the expected behavioral
  assertion and suiteguard reported no repository-boundary violation. Logs are
  `sabotage-9.log` and `sabotage-10.log` beside the temporary copies.

### Scope and settlement

Only Task 03's declared paths were changed. The pre-existing Task-file change
was its Daemon-owned `in_progress` status, which was preserved. No existing
test was edited. No Run schema, id generator, exact-id lookup, other command,
Run Browser, configuration, other skill/reference, Makefile, workflow, module
file, Task Graph or other Task file changed. No follow-up work was discovered.

The declared `## Verification` command and repository Verification were not
run; the Daemon runs them after this turn. No terminal Task verdict, commit,
push or pull request is claimed or performed.
