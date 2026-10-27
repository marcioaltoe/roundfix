---
task: task_02
spec: 0324-the-group-probe-waits-for-a-zombie-and-the-hard-link-rule-sees-all-forms
status: completed
type: test
complexity: medium
---

# Task 02: The Hard Link rule reports every reference in every module test file

## Overview

ADR-0324 keeps every executable a test may start a **Single-Name Test
Binary**, and `TestNoTestCreatesAHardLink` enforces it, but only for calls
through a named import in a test under `internal/`. A dot import of `os` and
`os.Link` taken as a function value pass it. This Task widens the rule with
`go/ast` to every reference, dot imports included, in every Go test file the
`go` command builds, and seeds each new form.

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-10-26-the-hard-link-rule-misses-a-dot-import-and-a-function-value.md`](references/2026-10-26-the-hard-link-rule-misses-a-dot-import-and-a-function-value.md)
   of 2026-10-26 (the Spec 0311 review finding F1 of 2026-10-24) by
   implementing API Contract 3 of `_techspec.md` in
   `internal/testfixture/hard_link_rule_test.go`.
2. MUST walk the repository root given to `hardLinkViolations`, skipping a
   directory named `testdata` or whose name starts with `.` or `_` (the root
   itself is never skipped), and keep parsing only `_test.go` files.
3. MUST report, in source order per file and with today's message format
   (`<path>:<line>: <pkg>.<Name> creates a hard link; copy or use a symbolic link`):
   - every selector `X.Link`, and `X.Linkat` when `X` is not `os`, whose `X`
     is an identifier with no resolved object that names a named or aliased
     import of `os`, `syscall` or `golang.org/x/sys/unix`, whether the
     selector is called or used as a value;
   - in a file with a `.` import of one of those packages, every identifier
     `Link`, or `Linkat` for `syscall` or `unix`, with no resolved object,
     that is not the selected member of a selector, the name of a function or
     method declaration, a field or parameter name, the key of a struct
     literal, or a label, attributed to the first such dot-imported package
     in import order.
4. MUST keep the `repository` subtest and the ten existing seeded subtests
   of `TestNoTestCreatesAHardLink` byte-identical, and add these seeded
   subtests to the same table, each writing
   `internal/seeded/fixture_test.go` (line numbers count from 1 at
   `package seeded`):
   - `dot import of os`: `import . "os"` and a call `Link("source", "target")`
     on line 3 → `internal/seeded/fixture_test.go:3: os.Link`;
   - `dot import of unix Linkat`: `import . "golang.org/x/sys/unix"` and a
     call to `Linkat` on line 3 → `...:3: unix.Linkat`;
   - `os Link as a function value`: `var link = os.Link` on line 3 and
     `func fixture() { apply(os.Link) }` on line 4 → `...:3: os.Link` and
     `...:4: os.Link`;
   - `dot-imported syscall Link as a value`: `import . "syscall"` and
     `var link = Link` on line 3 → `...:3: syscall.Link`;
   - `dot import shadowed by a local Link`: a local variable `Link` declared
     and called inside the function → no finding;
   - `dot import beside Link fields methods keys and labels`: a struct field
     `Link`, a method `Link` on another type and a call of it, a struct
     literal key `Link`, a label `Link` with `break Link`, and `_ = Symlink`
     → no finding;
   - `dot-imported os has no Linkat`: `_ = Linkat` with only `os`
     dot-imported → no finding.
5. MUST add a second table to `TestNoTestCreatesAHardLink` that writes the
   seed `os.Link("source", "target")` on line 3 of `fixture_test.go` in one
   directory per subtest: `a test outside internal` (`cmd/seeded`) →
   `cmd/seeded/fixture_test.go:3: os.Link`; `a test under skills` (`skills`)
   → `skills/fixture_test.go:3: os.Link`; `a testdata directory`
   (`internal/testdata/seeded`), `a hidden directory` (`.agents/seeded`) and
   `an underscore directory` (`internal/_seeded`) → no finding.
6. MUST update the file header's Boundary IN and Boundary OUT lines to the
   widened scope, and MUST NOT change any other file, including
   `docs/agents/specific-repository.md`, whose rule bullet stays true.

## Subtasks

- [ ] Walk every module test file the `go` command builds.
- [ ] Report selector references and dot-imported identifiers.
- [ ] Seed each new form and each location.

## Acceptance Criteria

- [ ] The rule finds no hard link in the repository and reports each new seeded form and location with its file and line.
- [ ] Each seeded non-reference stays unreported, and the ten earlier seeds pass unchanged.

## Consumers

None. Only the `internal/testfixture` rule test changes; no production contract or shared helper moves.

## Context

- instruction: `docs/adr/0324-a-test-never-gives-a-running-executable-a-second-name.md`
- instruction: `docs/agents/specific-repository.md`
- interface: `internal/testfixture/hard_link_rule_test.go`

## Verification

- `out="$(go test -count=1 -timeout 15m -v -run '^TestNoTestCreatesAHardLink$' ./internal/testfixture 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestNoTestCreatesAHardLink TestNoTestCreatesAHardLink/repository TestNoTestCreatesAHardLink/os_Link_names_its_site TestNoTestCreatesAHardLink/shadowed_os_import TestNoTestCreatesAHardLink/dot_import_of_os TestNoTestCreatesAHardLink/dot_import_of_unix_Linkat TestNoTestCreatesAHardLink/os_Link_as_a_function_value TestNoTestCreatesAHardLink/dot-imported_syscall_Link_as_a_value TestNoTestCreatesAHardLink/dot_import_shadowed_by_a_local_Link TestNoTestCreatesAHardLink/dot_import_beside_Link_fields_methods_keys_and_labels TestNoTestCreatesAHardLink/dot-imported_os_has_no_Linkat TestNoTestCreatesAHardLink/a_test_outside_internal TestNoTestCreatesAHardLink/a_test_under_skills TestNoTestCreatesAHardLink/a_testdata_directory TestNoTestCreatesAHardLink/a_hidden_directory TestNoTestCreatesAHardLink/an_underscore_directory; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; GOOS=windows go vet ./internal/testfixture` — expected: exit 0. Before this Task the new seeded subtests do not exist, so the pass check fails on the first of them. Afterwards the rule passes on the repository, every old and new seed passes, and the package compiles for Windows.

## References

- `_prd.md` → Goal 4; Core Feature 3; Success Metrics 4 and 5
- `_techspec.md` → Current behavior 6; Interfaces; API Contract 3; Invariant 4; Testing Approach; Build Order 2
- ADR-0324


## Result

Implemented the Task slice in `internal/testfixture/hard_link_rule_test.go`.
The scan now walks the supplied root, skips descendant `testdata`, hidden
and underscore directories, and parses only `_test.go` files. It reports
named and aliased selector references, including function values, and bare
unresolved dot-imported references in source order. Dot imports retain import
order and `os` never supplies `Linkat`. A syntax-name pass excludes selected
members, declaration names, field and parameter names, struct literal keys,
and labels; resolved local names remain excluded. Updated the boundary header.

Acceptance evidence:

- Repository and new forms/locations: the focused command
  `go test ./internal/testfixture -run '^TestNoTestCreatesAHardLink/' -count=1 -v`
  exited 0. The repository subtest reported no hard links. All seven new
  form seeds and five new location seeds passed, including both ordered
  `os.Link` function-value findings and the `cmd/` and `skills/` findings
  with the required paths and lines.
- Non-references and original seeds: the same focused check passed the
  local-shadow, fields/methods/keys/labels, missing `os.Linkat`, and skipped
  directory cases, plus all ten original seeded cases. A Python comparison
  against `git show HEAD:internal/testfixture/hard_link_rule_test.go`
  confirmed the repository subtest and the contiguous ten original seed
  entries remain byte-identical.
- Red starting point: after adding the seeds but before changing the
  detector, `go test ./internal/testfixture -run
  '^TestNoTestCreatesAHardLink/(dot|os_Link_as|a_test)' -count=1` exited 1:
  dot imports and function values were unreported, outside-internal paths
  errored, and the testdata seed was incorrectly reported.
- `gofmt` applied; `git diff --check` exited 0.

The declared Verification (including Windows vet) and repository Verification
were not run; the Daemon owns those checks and Task settlement. Status and
checkboxes were left unchanged. No follow-up outside this slice was needed.
