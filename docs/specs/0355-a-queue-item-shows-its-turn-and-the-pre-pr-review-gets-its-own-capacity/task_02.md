---
task: task_02
spec: 0355-a-queue-item-shows-its-turn-and-the-pre-pr-review-gets-its-own-capacity
status: pending
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
   `tokens`; no other expectation, fixture or test name changes.
7. MUST run sabotages 8 and 9 of `_techspec.md` → Testing Approach on a
   temporary copy and record which named test fails.
8. MUST NOT change the Delivery Engine, the Run Database, any other
   configuration key, the templates, any document or skill, the Makefile, CI
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

## Consumers

- `Delivery.ReviewConcurrency` and the operator-only key list — consumers: `./internal/config`
- The `Limits:` line of `deliver start` and `deliver status` — consumers: `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestDeliveryReviewConcurrencyIsReadFromUserConfig|TestProjectConfigCannotSetTheDeliveryReviewConcurrency|TestDeliveryReviewConcurrencyMustBePositive|TestDeliveryCheckTimeoutIsReadFromUserConfig|TestProjectConfigCannotSetTheDeliveryCheckTimeout)$' ./internal/config 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDeliveryReviewConcurrencyIsReadFromUserConfig TestProjectConfigCannotSetTheDeliveryReviewConcurrency TestDeliveryReviewConcurrencyMustBePositive; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./internal/config || exit 1; out="$(go test -count=1 -timeout 20m -v -run '^TestDeliver(Start|Status|Limits)' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestDeliverLimitsNameTheReviewConcurrency ' || { echo 'missing pass: TestDeliverLimitsNameTheReviewConcurrency' >&2; exit 1; }` — expected: exit 0. Before this Task the three configuration tests and the `Limits:` test do not exist, so the first pass line is missing; afterwards the key behaves as API Contract 1 states, the whole configuration package passes, and every `deliver start` and `deliver status` test passes with the `Limits:` line of API Contract 2.

## References

- `_prd.md` → Goals; Core Feature 6; Success Metric 5
- `_techspec.md` → Current behavior 4; Current behavior 6; API Contract 1; API Contract 2; Declared breaks; Testing Approach; Build Order 1
