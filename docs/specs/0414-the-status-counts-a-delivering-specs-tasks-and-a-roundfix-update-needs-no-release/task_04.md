---
task: task_04
spec: 0414-the-status-counts-a-delivering-specs-tasks-and-a-roundfix-update-needs-no-release
status: pending
type: backend
complexity: low
---

# Task 04: The worktree exclusion test gives the queue item its worktree

## Overview

Corrective Task 1 of 2. The QA gate of 2026-11-19 (Run `run_20261119T151201Z_05140c8b41c958d4`, report commit `550d0ef2`, finding F-01, row S03, task_03 Requirement 5.3) found that `TestADeliveryItemCountsTasksOnItsItemBranch/the_worktree_is_not_read` in `internal/cli/follow_up_queue_test.go` creates the Item Branch's linked worktree and leaves `task_02` completed there without committing, but never records that worktree as the queue item's `Worktree`. With `item.Worktree` empty, a reader that loads the item's worktree with `spec.Load` before Git still counts `running 1/2 Tasks`, so the subtest cannot catch the regression it names. task_01 Requirement 5 described the item only as naming its branch. The status view itself is right; only the fixture is incomplete.

## Requirements

1. MUST answer finding F-01 of the QA Report of 2026-11-19: in `TestADeliveryItemCountsTasksOnItsItemBranch`, set the queue item's `Worktree` to the linked worktree the subtest creates for the Item Branch, so `the_worktree_is_not_read` reads an item that names a worktree holding the uncommitted `task_02` completed status.
2. MUST keep that uncommitted completed status in the worktree and keep `the_worktree_is_not_read` expecting `running 1/2 Tasks`.
3. MUST keep the eight subtests, their names and their expected counts unchanged, and keep every one of them passing.
4. MUST NOT change `internal/cli/follow_up_queue.go`, `internal/cli/follow_up_status.go`, `internal/spec/spec.go`, any other test, any guide, the glossary, any record, any skill, the Makefile, a workflow, `go.mod` or a lint configuration.

## Subtasks

- [ ] Record the linked worktree as the queue item's `Worktree` in the fixture.
- [ ] Run the subtests and the worktree-first mutation.

## Acceptance Criteria

- [ ] Every subtest of `TestADeliveryItemCountsTasksOnItsItemBranch` passes, and a reader that loads the item's worktree with `spec.Load` before Git makes `the_worktree_is_not_read` fail with `running 2/2 Tasks`.

## Context

- interface: `internal/cli/follow_up_queue_test.go`
- instruction: `internal/cli/follow_up_queue.go`

## Consumers

- The test `TestADeliveryItemCountsTasksOnItsItemBranch` — consumers: `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^TestADeliveryItemCountsTasksOnItsItemBranch$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in a_running_item_counts_its_branch an_archived_item_counts_the_source_revision an_item_without_a_branch_counts_the_spec_root a_missing_branch_counts_the_spec_root a_branch_without_the_spec_counts_the_spec_root the_worktree_is_not_read the_pending_specs_line_counts_the_same the_rows_write_nothing; do printf '%s\n' "$out" | grep -q -- "--- PASS: TestADeliveryItemCountsTasksOnItsItemBranch/$name " || { echo "missing pass: $name" >&2; exit 1; }; done; root="$(pwd)"; dir="$(mktemp -d)"; trap 'rm -rf "$dir"' EXIT; awk '{ print } /item store[.]DeliveryQueueItem[)] [(][*]spec[.]Graph, error[)] [{]$/ { print "\tif graph, err := spec.Load(filepath.Join(item.Worktree, \"docs\", \"specs\"), item.SpecSlug); err == nil {\n\t\treturn graph, nil\n\t}"; n++ } END { if (n != 1) exit 1 }' internal/cli/follow_up_queue.go > "$dir/follow_up_queue.go" || { echo 'the worktree-first mutation did not apply' >&2; exit 1; }; printf '{"Replace":{"%s":"%s"}}\n' "$root/internal/cli/follow_up_queue.go" "$dir/follow_up_queue.go" > "$dir/overlay.json"; if out="$(go test -count=1 -timeout 20m -v -overlay "$dir/overlay.json" -run '^TestADeliveryItemCountsTasksOnItsItemBranch$/^the_worktree_is_not_read$' ./internal/cli 2>&1)"; then printf '%s\n' "$out"; echo 'the worktree-first mutation survives the_worktree_is_not_read' >&2; exit 1; fi; printf '%s\n' "$out" | grep -q -- '--- FAIL: TestADeliveryItemCountsTasksOnItsItemBranch/the_worktree_is_not_read ' || { printf '%s\n' "$out"; echo 'the_worktree_is_not_read did not fail under the mutation' >&2; exit 1; }; printf '%s\n' "$out" | grep -qF -- 'running 2/2 Tasks' || { printf '%s\n' "$out"; echo 'the mutation failed for another reason' >&2; exit 1; }` — expected: exit 0. Before this Task every subtest passes but the worktree-first mutation, applied through a `go test -overlay` copy of `internal/cli/follow_up_queue.go` that writes nothing to the checkout, also passes `the_worktree_is_not_read`, so the command exits 1 with `the worktree-first mutation survives the_worktree_is_not_read`; afterwards the eight subtests pass and the mutation makes `the_worktree_is_not_read` fail with `running 2/2 Tasks`.

## References

- QA Report of 2026-11-19 (Run `run_20261119T151201Z_05140c8b41c958d4`, commit `550d0ef2`) → F-01, row S03
- task_01 → Requirement 5
- task_03 → Requirement 5.3
