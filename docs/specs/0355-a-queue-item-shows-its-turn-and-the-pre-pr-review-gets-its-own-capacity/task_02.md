---
task: task_02
spec: 0355-a-queue-item-shows-its-turn-and-the-pre-pr-review-gets-its-own-capacity
status: completed
type: backend
complexity: medium
---

# Task 02: delivery.review_concurrency sets the Review Capacity, and the Limits line prints it

## Overview

`deliver status` prints a fixed `concurrency 1` that Fiscus read as the only
concurrency the queue has ("Concorrência fixa em 1, sem opção", reported by
fiscus on 2026-11-05,
[references](references/2026-11-05-fiscus-delivery-queue-takes-days.md)).
This Task adds the operator-only User Config key `delivery.review_concurrency`
(default `2`, at least `1`) and prints it on the `Limits:` line of `deliver
start` and `deliver status` as `review concurrency <n>` after `concurrency 1`.
Nothing reads the value yet besides the `Limits:` line; task_03 gives it to the
Delivery Engine.

## Requirements

1. MUST, before changing code, run
   `go test -count=1 -timeout 20m -v -run '^(TestDeliveryCheckTimeoutIsReadFromUserConfig|TestProjectConfigCannotSetTheDeliveryCheckTimeout)$' ./internal/config`
   and
   `go test -count=1 -timeout 20m -run '^TestDeliver(Start|Status)' ./internal/cli`
   and record in the Result that both pass; they pin the operator-only key
   pattern and today's `Limits:` line.
2. MUST add `DeliveryReviewConcurrencyKey`, `DefaultDeliveryReviewConcurrency`,
   `Delivery.ReviewConcurrency`, the overlay field, the Project Config removal
   with its warning and the validation exactly as `_techspec.md` → API
   Contract 1 states, beside the existing `delivery.check_timeout` handling in
   `internal/config/config.go`. The default configuration sets the field to
   `DefaultDeliveryReviewConcurrency`; `DefaultConfigYAML` and
   `DefaultProjectConfigYAML` stay byte for byte unchanged.
3. MUST change `printDeliveryLimits` in `internal/cli/deliver.go` as API
   Contract 2 states and pass `loaded.Config.Delivery.ReviewConcurrency` from
   `runDeliverStart` and `runDeliverStatus`.
4. MUST create `internal/config/delivery_review_concurrency_test.go`, opening
   with the `Suite`, `Invariant`, `Boundary IN` and `Boundary OUT` comment
   lines, each top-level test calling `t.Parallel()` first and using only
   `t.TempDir()` directories, with:
   - `TestDeliveryReviewConcurrencyIsReadFromUserConfig`: no User Config gives
     `2` and `DefaultDeliveryReviewConcurrency == 2`; a User Config of
     `delivery:\n  review_concurrency: 1\n` gives `1`; of `4` gives `4`.
   - `TestProjectConfigCannotSetTheDeliveryReviewConcurrency`: with a User
     Config of `3` and a work directory marked with
     `gittest.MarkRepository(t, work)` (Spec 0352), a Project Config
     `.roundfixrc.yml` of `delivery:\n  review_concurrency: <v>\n` for each
     `<v>` of `0`, `5`, `-1` and `two` loads without error, keeps `3`, and
     writes exactly
     `config: delivery.review_concurrency in Project Config is ignored; set delivery.review_concurrency in User Config\n`
     to the loader's stderr.
   - `TestDeliveryReviewConcurrencyMustBePositive`: a User Config of `0` or
     `-1` fails `Load` with an error containing
     `delivery.review_concurrency must be greater than 0`; a User Config of
     `two` fails `Load`.
5. MUST add `TestDeliverLimitsNameTheReviewConcurrency` to
   `internal/cli/deliver_limits_test.go` (`t.Parallel()` first): with
   `delivery:\n  review_concurrency: 3\n` in the test home's User Config,
   `deliver status` over a recorded queue prints
   `Limits: deadline none, retries per item none, concurrency 1, review concurrency 3, tokens none`
   as its `Limits:` line, and a `deliver start` the way
   `TestDeliverStartRecordsNoneForOmittedLimits` drives it prints the same
   line.
6. MUST change, in the sixteen tests `_techspec.md` → Declared breaks names,
   only the expected `Limits:` text, inserting `review concurrency 2, ` before
   `tokens`; no other expectation, fixture or test name changes. The same
   change applies to `TestDeliverStatusNamesAReviewThatOpensANewLineage` in
   `internal/cli/deliver_review_lineage_notice_test.go`, which Spec 0356
   added on the default branch after this plan was written and which pins the
   whole `deliver status` output, `Limits:` line included.
7. MUST run sabotages 8 and 9 of `_techspec.md` → Testing Approach on a
   temporary copy and record which named test fails.
8. MUST, as the last edit, re-record `docs/references/behavior-surfaces.json`
   with `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`,
   because the new key changes the record's `config keys` fingerprint; the
   repository Verification fails with a stale fingerprint otherwise.
9. MUST NOT change the Delivery Engine, the Run Database, any other
   configuration key, the templates, any document or skill other than the
   Behavior Surface Record of requirement 8, the Makefile, CI
   workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] Characterize the operator-only key pattern and the `Limits:` line.
- [ ] Add the key, its default, the Project Config warning and the validation.
- [ ] Print the review concurrency on the `Limits:` line of start and status.
- [ ] Add the configuration and `Limits:` tests and update the declared expectations.
- [ ] Run the two sabotages.

## Acceptance Criteria

- [ ] `delivery.review_concurrency` defaults to `2`, takes a User Config value of at least `1`, refuses `0`, `-1` and `two`, and a Project Config value is ignored with the operator-only warning.
- [ ] `deliver start` and `deliver status` print `concurrency 1, review concurrency <n>` on their `Limits:` line.
- [ ] Only the declared tests change, and only in their expected `Limits:` text.

## Context

- interface: `internal/config/config.go`
- interface: `internal/config/delivery_check_timeout_test.go`
- creates: `internal/config/delivery_review_concurrency_test.go`
- interface: `internal/cli/deliver.go`
- interface: `internal/cli/deliver_limits_test.go`
- interface: `internal/cli/deliver_lane_skip_visible_test.go`
- interface: `internal/cli/deliver_owner_status_test.go`
- interface: `internal/cli/deliver_park_status_test.go`
- interface: `internal/cli/deliver_revalidate_test.go`
- interface: `internal/cli/deliver_run_slot_wait_test.go`
- interface: `internal/cli/deliver_run_unresolved_park_test.go`
- interface: `internal/cli/deliver_test.go`
- interface: `internal/cli/deliver_token_ceiling_test.go`
- interface: `internal/cli/deliver_usage_test.go`
- interface: `internal/cli/deliver_review_lineage_notice_test.go`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- `Delivery.ReviewConcurrency` and the operator-only key list — consumers: `./internal/config`
- The `Limits:` line of `deliver start` and `deliver status` — consumers: `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestDeliveryReviewConcurrencyIsReadFromUserConfig|TestProjectConfigCannotSetTheDeliveryReviewConcurrency|TestDeliveryReviewConcurrencyMustBePositive|TestDeliveryCheckTimeoutIsReadFromUserConfig|TestProjectConfigCannotSetTheDeliveryCheckTimeout)$' ./internal/config 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDeliveryReviewConcurrencyIsReadFromUserConfig TestProjectConfigCannotSetTheDeliveryReviewConcurrency TestDeliveryReviewConcurrencyMustBePositive; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./internal/config || exit 1; out="$(go test -count=1 -timeout 20m -v -run '^TestDeliver(Start|Status|Limits)' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestDeliverLimitsNameTheReviewConcurrency ' || { echo 'missing pass: TestDeliverLimitsNameTheReviewConcurrency' >&2; exit 1; }` — expected: exit 0. Before this Task the three configuration tests and the `Limits:` test do not exist, so the first pass line is missing; afterwards the key behaves as API Contract 1 states, the whole configuration package passes, and every `deliver start` and `deliver status` test passes with the `Limits:` line of API Contract 2.

## References

- `_prd.md` → Goals; Core Feature 6; Success Metric 5
- `_techspec.md` → Current behavior 4; Current behavior 6; API Contract 1; API Contract 2; Declared breaks; Testing Approach; Build Order 1

## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, tech-writer, evidence-gate.

Added the operator-only `delivery.review_concurrency` key beside the existing
Check Wait handling. The built-in default is `2`; the overlay accepts a User
Config integer and validation refuses values below `1`. Project Config removes
the key before strict decoding and emits the exact operator-only warning, even
for invalid values. `deliver start` and `deliver status` pass the loaded value
to `printDeliveryLimits`, which prints it after `concurrency 1`.

### Characterization before code changes

Both required commands passed, exit `0`:

- `go test -count=1 -timeout 20m -v -run '^(TestDeliveryCheckTimeoutIsReadFromUserConfig|TestProjectConfigCannotSetTheDeliveryCheckTimeout)$' ./internal/config`
- `go test -count=1 -timeout 20m -run '^TestDeliver(Start|Status)' ./internal/cli`

The new tests first exposed the missing behavior: the configuration tests
could not compile without the new constant and field; the new command test
failed because strict decoding rejected `review_concurrency` as unknown.

### Acceptance evidence

- Configuration: `go test -count=1 -timeout 3m -v -run 'ReviewConcurrency' ./internal/config`
  passed, exit `0`. All three new named tests passed: default `2`, explicit `1`
  and `4`, User Config refusal of `0`, `-1` and `two`, and Project Config values
  `0`, `5`, `-1` and `two` ignored while retaining User Config `3` and printing
  exactly the required warning. Each top-level test begins with `t.Parallel()`;
  directory fixtures use only `t.TempDir()` and the Project Config fixture uses
  `gittest.MarkRepository`.
- Command output: `go test -count=1 -timeout 5m -v -run '^TestDeliver(LimitsNameTheReviewConcurrency|StartRecords.*|StartPrintsTokensNone|Status.*)$' ./internal/cli`
  passed, exit `0`. `TestDeliverLimitsNameTheReviewConcurrency` exercises start
  with omitted limits and status over its recorded queue, both with User Config
  `3`, and checks the exact `Limits: deadline none, retries per item none, concurrency 1, review concurrency 3, tokens none`
  line. The affected default-value transcripts, including the Spec 0356 review
  lineage transcript, also passed.
- Scope: a byte comparison against `git show HEAD:<path>`, after removing the
  new command test and reversing only `review concurrency 2, ` in existing
  expectations, reproduced all eleven original command test files exactly.
  `DefaultConfigYAML` and `DefaultProjectConfigYAML` are byte-identical to HEAD.
  No other existing test, fixture or expectation changed. `git diff --check`
  passed. The Delivery Engine, Run Database, other keys, templates, skills,
  guides and tooling files were not changed.

### Required sabotages

Both ran on an isolated temporary copy under `/private/tmp/roundfix-0355-task02-1x3zp6d9/repo`,
materialized from HEAD plus this Task's changes, with disposable Git metadata.
Neither mutation touched the assigned worktree.

- Sabotage 8 removed the review concurrency text and its formatting argument
  from `printDeliveryLimits`. `go test -count=1 -timeout 3m -v -run '^TestDeliverLimitsNameTheReviewConcurrency$' ./internal/cli`
  exited `1`: `TestDeliverLimitsNameTheReviewConcurrency` saw a successful start
  printing the old line without `review concurrency 3` and failed its output
  assertion.
- After restoring the command source, sabotage 9 removed the Project Config
  key stripping and allowed its overlay to apply. `go test -count=1 -timeout 3m -v -run '^TestProjectConfigCannotSetTheDeliveryReviewConcurrency$' ./internal/config`
  exited `1`: `TestProjectConfigCannotSetTheDeliveryReviewConcurrency/5` observed
  `ReviewConcurrency = 5, want 3`; `/0` and `/-1` failed validation, and `/two`
  failed decoding. All four subtests failed as intended.

### Final record and settlement boundary

The final repository edit regenerates `docs/references/behavior-surfaces.json`
with `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
No record bytes are edited by hand. The Agent handback reports the command's
outcome.

Declared Verification and repository Verification were not run; they remain
Daemon-owned. Task status and checkboxes are unchanged. No commit, push or
Pull Request was made. Giving this field to the Delivery Engine remains
`task_03`'s slice; there is no additional follow-up in this diff.
