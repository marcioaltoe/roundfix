---
task: task_02
spec: 0311-the-suite-audit-and-no-test-hard-links-its-binary
status: completed
type: test
complexity: low
---

# Task 02: A repository rule refuses a hard link in any test

## Overview

ADR-0324 decides that a test never gives an executable a second hard link,
because macOS kills a process whose executable path it resolves to a name
that has just disappeared. task_01 removes the only hard link the suite
creates. This Task adds the rule that keeps the class out of every package:
a test in `internal/testfixture` that parses every `_test.go` file under
`internal/` and fails on any call that creates a hard link, writes the rule
into the repository guide, and adds the **Single-Name Test Binary** term to
the glossary.

## Requirements

1. MUST answer the Run failures of 2026-10-23 recorded in
   [`docs/references/2026-10-24-test-suite-audit.md`](../../references/2026-10-24-test-suite-audit.md)
   (`_prd.md` → Sources 1) by implementing API Contract 2 of `_techspec.md`.
2. MUST add the new file `internal/testfixture/hard_link_rule_test.go`
   with a file header in the shape of `written_executable_test.go` (Suite,
   Invariant, Boundary IN, Boundary OUT) and
   `TestNoTestCreatesAHardLink`, which:
   - walks `internal/` from the repository root (`../..`), parses every
     `_test.go` file with `go/parser`, and reports `file:line` for each call
     `X.Link(...)` where `X` resolves, through the file's imports and their
     aliases, to `os`, `syscall` or `golang.org/x/sys/unix`, and for each
     call `X.Linkat(...)` on `syscall` or `golang.org/x/sys/unix`;
   - has a `repository` subtest that requires no finding in this
     repository;
   - has seeded subtests on a temporary tree under `t.TempDir()`: a seeded
     `os.Link` call is reported with its `internal/seeded/fixture_test.go:3`
     location; an aliased import (`filesystem "os"`) is reported; a
     `unix.Linkat` call is reported; a local variable or method named `Link`
     on another receiver, `os.Symlink` and `os.Readlink` are not reported.
3. MUST add one bullet to `## Repository rules` in
   `docs/agents/specific-repository.md`, next to the Signal rule bullet:
   "A test never creates a hard link: a test that needs an executable at
   another path copies it or links it symbolically, so every executable a
   test may start stays a **Single-Name Test Binary**; the Hard Link rule in
   `internal/testfixture` enforces it." The rest of the file keeps its
   bytes.
4. MUST add the term **Single-Name Test Binary** to `CONTEXT.md`, in the
   section that holds **Real-Time Budget** and **Parallel Test Package**,
   following the `domain-modeling` skill, with this definition: "An
   executable a test may start that keeps one directory entry while any
   process may start it; a test that needs it at another path copies it or
   links it symbolically, because macOS kills a new process whose
   executable path names an entry that has disappeared (ADR-0324)." with
   the ADR number written as a Markdown link to
   `docs/adr/0324-a-test-never-gives-a-running-executable-a-second-name.md`,
   in the shape of the neighbouring entries, and an `_Avoid_` line naming "hard-linked fixture, shared inode,
   binary alias".
5. MUST NOT change any other test, production file, skill, the Makefile or
   the CI workflows.

## Subtasks

- [ ] Add the Hard Link rule test with its repository and seeded subtests.
- [ ] Write the rule into the repository guide and the term into the glossary.

## Acceptance Criteria

- [ ] The rule finds no hard link in the repository and reports each seeded one with its location.
- [ ] `docs/agents/specific-repository.md` and `CONTEXT.md` carry the rule and the **Single-Name Test Binary** term.

## Consumers

None. The Task changes test code and documents only; no production contract or shared test helper outside its package changes.

## Context

- instruction: `internal/testfixture/written_executable_test.go`
- instruction: `internal/testfixture/signal_rule_test.go`
- instruction: `docs/adr/0324-a-test-never-gives-a-running-executable-a-second-name.md`
- instruction: `docs/references/2026-10-24-test-suite-audit.md`
- creates: `internal/testfixture/hard_link_rule_test.go`
- interface: `docs/agents/specific-repository.md`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -timeout 15m -v -run '^TestNoTestCreatesAHardLink$' ./internal/testfixture 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestNoTestCreatesAHardLink ' || { echo 'missing pass: TestNoTestCreatesAHardLink' >&2; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestNoTestCreatesAHardLink/repository ' || { echo 'missing repository subtest' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Single-Name Test Binary**:' && tr -s '[:space:]' ' ' < docs/agents/specific-repository.md | grep -qF -- 'the Hard Link rule in' && GOOS=windows go vet ./internal/testfixture` — expected: exit 0. Before this Task the test does not exist and neither document carries the rule or the term, so the command fails at its first pass check; with an `os.Link` call added to any `internal` test file the repository subtest fails naming its location. Afterwards the rule passes on the repository, the guide carries the rule and the glossary the term, and the package compiles for Windows.

## References

- `_prd.md` → Goals; Core Feature 2; Success Metric 3; Glossary
- `_techspec.md` → API Contract 2; Invariant 3; Testing Approach; Glossary; Build Order 2

## Result

Implemented API Contract 2 of the TechSpec and ADR-0324 in this Task's slice.

- Acceptance criterion 1: added `internal/testfixture/hard_link_rule_test.go`
  with the Suite/Invariant/Boundary header and `TestNoTestCreatesAHardLink`.
  The AST scan walks every `_test.go` under `internal/`, resolves imports and
  aliases for `os`, `syscall` and `golang.org/x/sys/unix`, and reports each
  forbidden `Link`/`Linkat` call with its repository-relative file and line.
  The repository subtest found no violation. Ten seeded subtests passed,
  covering all five package/API combinations, aliases, multiple findings,
  the exact `internal/seeded/fixture_test.go:3` location, other receivers,
  a local `Link`, a shadowed `os` import, `Symlink` and `Readlink`.
- Acceptance criterion 2: inserted the exact requested repository-rule
  bullet immediately after the Signal rule and the Single-Name Test Binary
  glossary entry between Parallel Test Package and Real-Time Budget, with
  the ADR-0324 link and requested `_Avoid_` terms. A byte-level comparison
  against `HEAD` confirmed that removing those exact insertions reproduces
  each original document; the ADR link target exists.

Focused checks:

- Before editing, a Python assertion confirmed the guard file and both
  Single-Name Test Binary document entries were absent.
- `go test -race -count=1 -timeout 3m -v -run '^(TestNoTestCreatesAHardLink|TestNoTestWritesAnExecutableOutsideTheResidue)$' ./internal/testfixture`
  exited 0: both rules and every subtest passed, including the new
  repository scan and ten seeded cases.
- `GOOS=windows go test -c -o "$artifact_dir/testfixture.test.exe" ./internal/testfixture`
  exited 0, with `artifact_dir` created by `mktemp` under `/private/tmp`;
  the compiled artifact and its temporary directory were removed.
- `gofmt -l internal/testfixture/hard_link_rule_test.go` returned no paths;
  `git diff --check` exited 0.

The checkout's symlinked writable root prevented sandbox process startup;
checks and edits used its resolved `/Volumes/Brain-SSD/...` checkout with
approved escalation. No host configuration changed.

Declared Task Verification and repository Verification were not run; the
Daemon owns those checks and settlement. Task status was left untouched.
No other Task, Task Graph, production file, existing test, skill, Makefile
or CI workflow was changed. No follow-up was needed within this slice.

## Carry-forward provenance

- Source Run: `run_20261024T115610Z_1c57415017fc7072`
- Source commit: `27a1d6181937cffd75e1a05a234e14be5f77ffdf`
