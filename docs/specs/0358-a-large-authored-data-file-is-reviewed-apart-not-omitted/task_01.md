---
task: task_01
spec: 0358-a-large-authored-data-file-is-reviewed-apart-not-omitted
status: completed
type: backend
complexity: medium
---

# Task 01: Project Config declares the Review Data Paths

## Overview

Project Config can omit a file from the pre-PR review but cannot say that a
file is authored data the reviewer must still read; `pre_pr_review` refuses
any key but `provider` and `omit_paths`. This Task adds
`pre_pr_review.data_paths`, its validation and its matcher, and documents it
in the configuration guide, as ADR-0368 records. It answers point 5 of the
Fiscus report adopted as
[references/2026-11-05-fiscus-the-delivery-queue-takes-days-and-a-large-data-catalog-is-omitted-from-review.md](references/2026-11-05-fiscus-the-delivery-queue-takes-days-and-a-large-data-catalog-is-omitted-from-review.md).

## Requirements

1. MUST add `DataPaths []string` to `PrePRReview` and `DataPaths *[]string`
   with YAML name `data_paths` to `prePRReviewOverlay`
   (`internal/config/config.go`), applied like `OmitPaths`: a Project
   Config list replaces the User Config list and `[]` clears it.
2. MUST refuse a non-list value and a non-string entry with the texts of
   `_techspec.md` → API Contract 1, in the same parser check that guards
   `omit_paths`, and validate the entries with `validateReviewDataPaths`
   (`internal/config/pre_pr_review.go`), which calls
   `validateDerivedPathPatterns("pre_pr_review.data_paths", paths)`, both when
   a config file is parsed and in `Validate`.
3. MUST add `PrePRReview.ReviewsApart(name string) bool`, matching an exact
   file, a directory prefix ending in `/` or a whole-path Go `path.Match`
   pattern exactly as `Omits` does, with the two sharing one matcher.
4. MUST keep `omit_paths`, `provider`, `delivery.derived_paths` and every
   existing message byte for byte, and MUST NOT change `internal/cli`, the
   `init` templates or any existing test.
5. MUST add `internal/config/pre_pr_review_data_paths_test.go`, opening with
   the `Suite`, `Invariant`, `Boundary IN` and `Boundary OUT` comment
   lines, with `TestPrePRReviewDataPathsLoadsAndValidates` of
   `_techspec.md` → Testing Approach, through `ResolveConfigProposal` only:
   rows `absent`, `inherits user`, `project replaces user`, `empty clears
   user`; refusal rows `a mapping`, `a number entry`, `an absolute path`,
   `a parent segment` and `a bad pattern`, each asserting its API Contract 1
   text; and `ReviewsApart` on an exact file, a directory and a pattern. The
   test calls `t.Parallel()` first.
6. MUST add to `docs/user-guide/configuration.md` the row and the sentence of
   `_techspec.md` → API Contract 6 for that guide, word for word, and change
   nothing else in it.
7. MUST re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
   as the Task's last step, after the guide edit, and never hand-edit it.

## Subtasks

- [ ] Add the key, its validation and the matcher.
- [ ] Add the config test.
- [ ] Document the key in the configuration guide, then re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] `pre_pr_review.data_paths` loads, replaces and clears like `omit_paths`, and every malformed value is refused with its own text.
- [ ] `ReviewsApart` matches files, directories and patterns as `Omits` does.
- [ ] The configuration guide documents the key and where an authored catalog belongs.

## Context

- instruction: `docs/adr/0368-a-large-authored-data-file-is-reviewed-apart-in-data-review-segments.md`
- interface: `internal/config/config.go`
- interface: `internal/config/pre_pr_review.go`
- creates: `internal/config/pre_pr_review_data_paths_test.go`
- instruction: `internal/config/pre_pr_review_omit_paths_test.go`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The `pre_pr_review` Project Config schema and `PrePRReview.ReviewsApart` — consumers: `./internal/config`, `./internal/cli`
- The configuration guide fingerprint in the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^TestPrePRReviewDataPathsLoadsAndValidates$' ./internal/config 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestPrePRReviewDataPathsLoadsAndValidates ' || { echo 'missing pass: TestPrePRReviewDataPathsLoadsAndValidates' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -q -- 'belongs in .pre_pr_review.data_paths., not in' || { echo 'missing guide sentence' >&2; exit 1; }; go test -count=1 -timeout 20m ./internal/config ./internal/cli && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the test does not exist, so its pass line is missing; afterwards it passes, the guide states where an authored catalog belongs, and the config and CLI suites and the documentation contracts pass.

## References

- `_prd.md` → Goals; Core Feature 1; User Story 1
- `_techspec.md` → Current behavior 1; API Contract 1; API Contract 6; Invariant 5; Testing Approach; Build Order 1

## Result

Skills loaded: implement-task, coding-guidelines, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, tech-writer, evidence-gate.

- Added `DataPaths` to resolved and overlay config, with Project-over-User
  replacement and explicit empty-list clearing. The shared YAML-node check
  refuses non-lists and non-string entries; `validateReviewDataPaths` delegates
  to `validateDerivedPathPatterns` during parsing and in `Validate`.
- Added `ReviewsApart` and extracted the existing omission matcher so both
  methods retain exact-file, directory-prefix and whole-path pattern behavior,
  including exact names that contain pattern metacharacters.
- Added `TestPrePRReviewDataPathsLoadsAndValidates` through
  `ResolveConfigProposal` only. It covers `absent`, `inherits user`, `project
  replaces user`, `empty clears user`, all five required refusal rows with
  exact errors in both config scopes, positive and negative matching, and that
  declaring data alone does not omit it. Every test calls `t.Parallel()` first.
- Added precisely the API Contract 6 row and sentence to the configuration
  guide. Re-recorded the Behavior Surface Record after the guide edit as the
  last implementation step; only `config keys` and the configuration guide
  fingerprint changed. No generated record was hand-edited.

Focused evidence (commands executed through `rtk proxy`):

1. Before implementation,
   `go test -count=1 -timeout 2m ./internal/config -run 'TestPrePRReview(Data|Omit)PathsLoadsAndValidates'`
   exited 1 because `DataPaths` and `ReviewsApart` were undefined, establishing
   the missing API before adding it.
2. After implementation,
   `go test -count=1 -timeout 2m -v ./internal/config -run 'TestPrePRReview(Data|Omit)PathsLoadsAndValidates|TestPrePRReviewPolicyResolution|TestPrePRReviewProviderRefusesNullValue|TestProjectConfigReadsDerivedPathDeclarations|TestDerivedPathDeclarationsRefuseUnsafeEntries|TestDerivedLineDeclarationsAreReadAndValidated'`
   exited 0. The new test passed every precedence, refusal and matching row;
   the unchanged omission, provider and derived-path tests also passed.
3. `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
   exited 0. The exact guide strings were checked against the TechSpec before
   writing, and the resulting guide and generated-record diffs were inspected.
4. `git diff --check` exited 0. Scope inspection found only the two declared
   config sources, the new test, the configuration guide, its generated record
   and this Task file. Existing tests, CLI code and init templates are unchanged.

Acceptance evidence: item 1 is covered by the precedence and exact refusal
rows plus inspection of both validation call sites; item 2 by the whole-path
matching row and unchanged omission tests; item 3 by the exact guide additions
and successful sanctioned recorder.

Declared Verification and repository Verification were not run; the Daemon
owns them and Task settlement. The initial Task change was its Daemon-written
`in_progress` status, which this turn preserves. No other Task or Task Graph
was edited. No follow-up work was added to this slice.
