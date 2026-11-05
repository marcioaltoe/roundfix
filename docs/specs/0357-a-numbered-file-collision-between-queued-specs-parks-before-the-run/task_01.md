---
task: task_01
spec: 0357-a-numbered-file-collision-between-queued-specs-parks-before-the-run
status: pending
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
