---
task: task_01
spec: 0357-a-numbered-file-collision-between-queued-specs-parks-before-the-run
status: completed
type: backend
complexity: medium
---

# Task 01: Project Config declares the numbered files under delivery.numbered_files

## Overview

Today no configuration names the files a repository numbers in sequence
(`_techspec.md` → Current behavior 5). After this Task User and Project Config
accept `delivery.numbered_files`, a list of path patterns with the Derived
Path Declaration grammar, empty by default, and `(Delivery).IsNumberedFile`
tells a Numbered File from any other matching file. This answers point 3 of
the report reported by fiscus on 2026-11-05
([references](references/2026-11-05-fiscus-delivery-queue-report.md)): the
declaration is how a repository tells the queue which numbers to compare.

## Requirements

1. MUST, before changing code, run
   `go test -count=1 -timeout 20m -run '^TestDerivedLineDeclarationsAreReadAndValidated$' ./internal/config`
   and record in the Result that it passes, and that
   `ResolveConfigProposal(nil, []byte("delivery:\n  numbered_files: [db/*.sql]\n"))`
   is refused today as an unknown key.
2. MUST add `NumberedFiles []string` to `Delivery`, the overlay field
   `NumberedFiles *[]string` with the YAML tag `numbered_files`, its merge in
   `applyConfigContent` beside the `derived_paths` merge (a project value
   replaces the user value, an absent key inherits it, an explicit `[]`
   clears it), and its validation in `Validate` with
   `validateDerivedPathPatterns("delivery.numbered_files", config.Delivery.NumberedFiles)`,
   exactly as `_techspec.md` → API Contract 1 states.
3. MUST add `func (delivery Delivery) IsNumberedFile(name string) bool` as
   API Contract 2 states: a `matchesDerivedPath` match over `NumberedFiles`
   whose `path.Base` starts with a byte `0` through `9`.
4. MUST create `internal/config/numbered_files_test.go`, each top-level test
   calling `t.Parallel()` first, with:
   - `TestNumberedFileDeclarationsAreReadAndValidated`, through
     `ResolveConfigProposal`, with subtests `absent` (nil list),
     `project_replaces_user` (user `[a/*.sql]`, project `[drizzle/*.sql]`
     gives `[drizzle/*.sql]`), `user_inherited` (user list, project without
     the key), `explicit_empty` (user list, project `[]` gives an empty list),
     `unsafe` (`[../x]`, `[/abs]` and `['a\b']` refused with
     `delivery.numbered_files has unsafe path`) and `invalid` (`['[']`
     refused with `delivery.numbered_files has invalid pattern`);
   - `TestIsNumberedFileNeedsAMatchAndALeadingDigit`, with `[drizzle/]`:
     `drizzle/0052_a.sql` and `drizzle/meta/0052_snapshot.json` true,
     `drizzle/meta/_journal.json` false; with `[drizzle/*.sql]`:
     `drizzle/0052_a.sql` true, `drizzle/x0052.sql`, `db/0052_a.sql` and
     `drizzle/meta/0052_snapshot.json` false; with no patterns: false.
5. MUST run sabotages 1 and 2 of `_techspec.md` → Testing Approach on a
   temporary copy and record which named test fails.
6. MUST, as the LAST step, re-record the Behavior Surface Record, whose
   `config keys` fingerprint reads the configuration's YAML tags, with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
7. MUST keep every existing test of `./internal/config` passing unchanged.
8. MUST NOT change `delivery.derived_paths`, `DeliveryConfigAtCommit`, Git
   root detection, `.roundfixrc.yml`, any guide, skill or reference, the
   Makefile, CI workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] Characterize the refusal of the unknown key and the derived paths test.
- [ ] Add the field, overlay, merge and validation.
- [ ] Add `IsNumberedFile`.
- [ ] Add the two tests and run the two sabotages.
- [ ] Re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] `delivery.numbered_files` is read from User and Project Config with project precedence and an explicit empty list, and an unsafe path or invalid pattern is refused naming the key.
- [ ] `IsNumberedFile` is true only for a matched file whose name starts with a digit.
- [ ] The Behavior Surface Record is current.

## Context

- interface: `internal/config/config.go`
- interface: `internal/config/delivery.go`
- creates: `internal/config/numbered_files_test.go`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- `(config.Delivery).IsNumberedFile` and the `delivery.numbered_files` key — consumers: `./internal/config`, `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestNumberedFileDeclarationsAreReadAndValidated|TestIsNumberedFileNeedsAMatchAndALeadingDigit|TestDerivedLineDeclarationsAreReadAndValidated)$' ./internal/config 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestNumberedFileDeclarationsAreReadAndValidated TestIsNumberedFileNeedsAMatchAndALeadingDigit TestNumberedFileDeclarationsAreReadAndValidated/project_replaces_user TestNumberedFileDeclarationsAreReadAndValidated/explicit_empty TestNumberedFileDeclarationsAreReadAndValidated/unsafe; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$'` — expected: exit 0. Before this Task neither new test exists, so the first pass line is missing; afterwards the key is read and validated, the predicate holds, the existing derived paths test passes and the Behavior Surface Record carries the new key.

## References

- `_prd.md` → Goal 1; User Story 1; User Story 5; Core Feature 1; Core Feature 2; Success Metric 1
- `_techspec.md` → Current behavior 5; API Contract 1; API Contract 2; Testing Approach; Build Order 1


## Result

Skills loaded: implement-task, coding-guidelines, golang-structs-interfaces, golang-testing, testing-boss, tech-writer, evidence-gate.

Added `Delivery.NumberedFiles`, its optional `numbered_files` YAML overlay,
and shared path-pattern validation. The merge sits beside `DerivedPaths` in
`applyOverlay`, which `applyConfigContent` calls in this checkout. Project
values replace User Config, absence inherits, and `[]` clears the list.
`IsNumberedFile` reuses `matchesDerivedPath` and requires the basename's first
byte to be an ASCII digit. Existing derived-path behavior and tests are unchanged.

Before changing code:

- `go test -count=1 -timeout 20m -run '^TestDerivedLineDeclarationsAreReadAndValidated$' ./internal/config`
  exited 0 (`ok`, 0.562s).
- A temporary `go test -overlay` probe called
  `ResolveConfigProposal(nil, []byte("delivery:\n  numbered_files: [db/*.sql]\n"))`.
  `TestNumberedFilesUnknownKeyCharacterization` passed by observing
  `field numbered_files not found in type config.deliveryOverlay`.
  The probe created no repository file.

Focused acceptance evidence:

| Criterion | Evidence |
| --- | --- |
| Config precedence, clearing and validation | `TestNumberedFileDeclarationsAreReadAndValidated` covers `absent` (nil), `project_replaces_user`, `user_inherited`, `explicit_empty` (non-nil empty list), all three unsafe patterns and the invalid pattern, including the required key-specific diagnostics. |
| Matching and leading digit | `TestIsNumberedFileNeedsAMatchAndALeadingDigit` covers every required directory/glob/no-pattern case, ASCII boundary digits, punctuation, a non-ASCII digit and an empty name. Both new top-level tests call `t.Parallel()` first. |
| Existing config behavior | `go test -count=1 -timeout 20m ./internal/config` exited 0 (`ok`, 1.943s), exercising all existing config tests unchanged and both new tests. |
| Behavior Surface Record | As the final implementation step, `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` exited 0 (`ok`, 0.439s), regenerating and checking `docs/references/behavior-surfaces.json` from the YAML tags. |

Required sabotages ran on a temporary copy of the working tree with isolated
Git metadata, without a commit; the copy was removed afterward:

1. Dropping the leading-digit check: `go test -count=1 -timeout 20m -v -run '^TestIsNumberedFileNeedsAMatchAndALeadingDigit$' ./internal/config`
   exited 1. `TestIsNumberedFileNeedsAMatchAndALeadingDigit/directory_journal`
   failed because `drizzle/meta/_journal.json` returned true, expected false.
   Other non-digit cases also failed.
2. Removing the numbered-files overlay merge: `go test -count=1 -timeout 20m -v -run '^TestNumberedFileDeclarationsAreReadAndValidated$/^project_replaces_user$' ./internal/config`
   exited 1. `TestNumberedFileDeclarationsAreReadAndValidated/project_replaces_user`
   failed with a nil list instead of `[]string{"drizzle/*.sql"}`.

The prescribed failures were behavioral assertion failures; neither was a
suite-guard refusal. No sabotage changed the working checkout.

Daemon-owned declared Verification and repository Verification were not run.
The pre-existing `status: in_progress` remains unchanged. No Task Graph,
other Task, guide, skill, Project Config, dependency or delivery change was
made. No follow-up was found within this slice.

## Carry-forward provenance

- Source Run: `run_20261105T211401Z_cae788fdc7c4d6e1`
- Source commit: `fa728f5395b9f394ac10d90431162b51c1495683`
