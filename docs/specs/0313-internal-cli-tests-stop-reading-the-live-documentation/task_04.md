---
task: task_04
spec: 0313-internal-cli-tests-stop-reading-the-live-documentation
status: completed
type: test
complexity: low
---

# Task 04: The upgrade characterization test runs in a fixture workspace

## Overview

The QA gate of Run `run_20261025T122115Z_fe88d4e7bd75bcb2` found that a
guide edit still re-runs `internal/cli` after task_01 and task_02 (F-01). A
fourteenth untagged test reads the live documentation:
`TestUpgradeStdoutAndExitCodesAreCharacterized` runs `roundfix upgrade` with
the package directory as its working directory, so the upgrade notice
inventories this repository's history and specs. It is the only test in its
file that does not call `withCLIWorkspace`. This Task gives it a fixture
workspace, so that no untagged `internal/cli` test reads the documentation
trees. That is Goal 1.

## Requirements

1. MUST answer F-01 of the QA report of Run
   `run_20261025T122115Z_fe88d4e7bd75bcb2` (`_prd.md` → Sources 2) by
   implementing Invariant 1 of `_techspec.md` for the reader that
   `_techspec.md` → Current behavior 6 measured.
2. MUST add one line, `withCLIWorkspace(t)`, directly after `t.Parallel()`
   at the top of `TestUpgradeStdoutAndExitCodesAreCharacterized` in
   `internal/cli/upgrade_characterization_test.go`. The test's name, its
   table of cases and every assertion keep their bytes, and every other test
   and helper keeps its bytes.
3. MUST NOT change `internal/suiteguard`, a production file, the Makefile,
   the CI workflows, `internal/docscontract` or any file task_01 or task_02
   changed.

## Subtasks

- [ ] Give the upgrade characterization test a fixture home and working directory.

## Acceptance Criteria

- [ ] Go's test input log of the upgrade characterization test names no path under the repository's documentation trees, `.agents`, `skills`, `CONTEXT.md` or `.roundfixrc.yml`, and the test passes.

## Consumers

None. One test gains a fixture workspace; no production contract changes.

## Context

- instruction: `docs/adr/0326-a-code-test-never-reads-the-live-documentation-trees.md`
- instruction: `internal/cli/upgrade.go`
- instruction: `internal/cli/upgrade_test.go`
- interface: `internal/cli/upgrade_characterization_test.go`

## Verification

- `d="$(mktemp -d)" && mkdir "$d/home" && go test -buildvcs=false -c -o "$d/cli.test" ./internal/cli && (cd internal/cli && HOME="$d/home" "$d/cli.test" -test.run '^TestUpgradeStdoutAndExitCodesAreCharacterized$' -test.count 1 -test.timeout 10m -test.testlogfile "$d/inputs.log" >"$d/out.txt" 2>&1) || { cat "$d/out.txt"; exit 1; }; test -s "$d/inputs.log" || { echo 'empty test input log' >&2; exit 1; }; for r in "$(pwd)" "$(pwd -P)"; do if grep -E '^(open|stat) ' "$d/inputs.log" | grep -F -e "$r/docs" -e "$r/.agents" -e "$r/skills/" -e "$r/CONTEXT.md" -e "$r/.roundfixrc.yml"; then echo 'the upgrade characterization test read the live repository' >&2; exit 1; fi; done` — expected: exit 0. Before this Task the test's input log names `docs/history`, `docs/adr`, `docs/specs`, `docs/user-guide` and `.roundfixrc.yml` under the checkout, so the command fails. With the new line removed it fails the same way. Afterwards the test passes and its input log names none of them.

## References

- `_prd.md` → Sources; Goals; Success Metric 3; Success Metric 5
- `_techspec.md` → Current behavior; Invariant 1; Testing Approach; Build Order 3

## Result

Implementation: added `withCLIWorkspace(t)` directly after `t.Parallel()` at
the top of `TestUpgradeStdoutAndExitCodesAreCharacterized` in
`internal/cli/upgrade_characterization_test.go`. The test file diff is exactly
one added line; the test name, its case table and every assertion keep their
bytes, and no other test or helper changed. No production file, Makefile, CI
workflow, `internal/suiteguard` or `internal/docscontract` file was touched.

Focused checks (not the Daemon-owned `## Verification` command):
- Confirmed the source diff: `git diff internal/cli/upgrade_characterization_test.go`
  shows a single `+	withCLIWorkspace(t)` line and nothing else.
- Ran the task's verification command locally as a focused check:
  compiled `./internal/cli` to a temp binary, ran the single test with
  `-test.testlogfile`, then grepped the input log for repository paths under
  `docs`, `.agents`, `skills/`, `CONTEXT.md` and `.roundfixrc.yml`. Result:
  `FOCUSED CHECK PASSED` (exit 0) — the test passes and its input log names
  none of the forbidden live-documentation paths.

Acceptance criteria evidence:
- [x] Go's test input log of the upgrade characterization test names no path
      under the repository's documentation trees, `.agents`, `skills`,
      `CONTEXT.md` or `.roundfixrc.yml`, and the test passes — verified by the
      focused check above (non-empty input log, no forbidden matches, exit 0).

## Carry-forward provenance

- Source Run: `run_20261025T193403Z_ade2bddbacb762cf`
- Source commit: `5bf0be6119df16da2700702bacbcda7a4022f1b6`
